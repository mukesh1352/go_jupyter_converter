package parser1

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"go_compiler/internal/utils"
)

type Notebook struct {
	Cells []Cell `json:"cells"`
}

type Cell struct {
	CellType string   `json:"cell_type"`
	Source   []string `json:"source"`
}

type ExecutionResult struct {
	CellNumber int
	Success    bool
	Stdout     string
	Stderr     string
	Images     []string
}

// RunJupyterChecker processes a Jupyter notebook and sends code cells to Python
func RunJupyterChecker(filePath string, cfg *utils.Config) {
	// Read the Jupyter Notebook file
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading the file %s: %v", filePath, err)
	}

	var nb Notebook
	if err := json.Unmarshal(data, &nb); err != nil {
		log.Fatalf("Failed to parse notebook JSON: %v", err)
	}

	// Initialize path rewriter for dataset normalization
	rewriter := utils.NewPathRewriter(cfg.DatasetRoot)

	// Ensure output directory exists
	if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory %s: %v", cfg.OutputDir, err)
	}

	// Create output file
	notebookName := filepath.Base(filePath)
	outputName := strings.TrimSuffix(notebookName, filepath.Ext(notebookName)) + "_output.txt"
	outputPath := filepath.Join(cfg.OutputDir, outputName)

	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("Failed to create output file %s: %v", outputPath, err)
	}
	defer outputFile.Close()

	// Write header to output file
	header := fmt.Sprintf("Execution Log for: %s\nTime: %s\n\n", notebookName, time.Now().Format(time.RFC1123))
	if _, err := outputFile.WriteString(header); err != nil {
		log.Printf("⚠️ Warning: Failed to write header to output file: %v", err)
	}

	// Start a persistent Python subprocess
	cmd := exec.Command("python3", "../../internal/executor/run_code.py")

	// Set environment variables
	cmd.Env = append(os.Environ(),
		"OUTPUT_DIR="+cfg.OutputDir,
		"DATASET_ROOT="+cfg.DatasetRoot,
	)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("Failed to get stdin: %v", err)
	}
	defer stdin.Close()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to get stdout: %v", err)
	}

	cmd.Stderr = os.Stderr // Forward stderr to console

	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start Python process: %v", err)
	}

	reader := bufio.NewScanner(stdout)
	var results []ExecutionResult

	// Process each code cell
	for i, cell := range nb.Cells {
		if cell.CellType != "code" {
			continue
		}

		// Join all lines into a single code string
		code := strings.Join(cell.Source, "")

		// Rewrite dataset paths
		code = rewriter.Rewrite(code)

		// Prepare JSON payload
		input := map[string]string{"code": code}
		jsonData, err := json.Marshal(input)
		if err != nil {
			log.Printf("Failed to marshal JSON for cell %d: %v", i, err)
			continue
		}

		// Send code to Python
		_, err = stdin.Write(append(jsonData, '\n'))
		if err != nil {
			log.Printf("Failed to send code for cell %d: %v", i, err)
			continue
		}

		// Read Python's response
		if reader.Scan() {
			var result map[string]interface{}
			line := reader.Text()

			if err := json.Unmarshal([]byte(line), &result); err != nil {
				log.Printf("Failed to parse Python output for cell %d: %v\nOutput: %s", i, err, line)
				continue
			}

			// Create execution result
			execResult := ExecutionResult{
				CellNumber: i,
				Success:    result["success"].(bool),
				Stdout:     result["stdout"].(string),
				Stderr:     result["stderr"].(string),
			}

			if images, ok := result["images"].([]interface{}); ok {
				for _, img := range images {
					if imgPath, ok := img.(string); ok {
						execResult.Images = append(execResult.Images, imgPath)
					}
				}
			}

			results = append(results, execResult)

			// Write to console
			printCellResult(execResult)

			// Write to file
			writeCellResult(outputFile, execResult)
		} else {
			log.Printf("No response from Python for cell %d", i)
		}
	}

	// Send exit command to Python subprocess
	exitPayload := map[string]string{"code": "__EXIT__"}
	exitJSON, _ := json.Marshal(exitPayload)
	stdin.Write(append(exitJSON, '\n'))

	cmd.Wait()

	// Write summary to output file
	writeSummary(outputFile, results)
	log.Printf("✅ Execution complete. Results saved to: %s", outputPath)
}

func printCellResult(result ExecutionResult) {
	fmt.Printf("\n--- Cell %d Output ---\n", result.CellNumber)
	fmt.Println("Success:", result.Success)
	
	if result.Stdout != "" {
		fmt.Println("Stdout:")
		fmt.Println(result.Stdout)
	}
	
	if result.Stderr != "" {
		fmt.Println("Stderr:")
		fmt.Println(result.Stderr)
	}
	
	if len(result.Images) > 0 {
		fmt.Println("Generated Images:")
		for _, img := range result.Images {
			fmt.Println("-", img)
		}
	}
	
	fmt.Println("----------------------")
}

func writeCellResult(file *os.File, result ExecutionResult) {
	output := fmt.Sprintf("\n--- Cell %d Output ---\n", result.CellNumber)
	output += fmt.Sprintf("Success: %v\n", result.Success)
	
	if result.Stdout != "" {
		output += fmt.Sprintf("Stdout:\n%s\n", result.Stdout)
	}
	
	if result.Stderr != "" {
		output += fmt.Sprintf("Stderr:\n%s\n", result.Stderr)
	}
	
	if len(result.Images) > 0 {
		output += "Generated Images:\n"
		for _, img := range result.Images {
			output += fmt.Sprintf("- %s\n", img)
		}
	}
	
	output += "----------------------\n"
	
	if _, err := file.WriteString(output); err != nil {
		log.Printf("⚠️ Warning: Failed to write cell %d results to file: %v", result.CellNumber, err)
	}
}

func writeSummary(file *os.File, results []ExecutionResult) {
	successCount := 0
	failureCount := 0
	imageCount := 0

	for _, result := range results {
		if result.Success {
			successCount++
		} else {
			failureCount++
		}
		imageCount += len(result.Images)
	}

	summary := fmt.Sprintf("\n\n=== Execution Summary ===\n")
	summary += fmt.Sprintf("Total Cells Processed: %d\n", len(results))
	summary += fmt.Sprintf("Successful Executions: %d\n", successCount)
	summary += fmt.Sprintf("Failed Executions: %d\n", failureCount)
	summary += fmt.Sprintf("Images Generated: %d\n", imageCount)
	summary += fmt.Sprintf("=======================\n")

	if _, err := file.WriteString(summary); err != nil {
		log.Printf("⚠️ Warning: Failed to write summary to file: %v", err)
	}
}
