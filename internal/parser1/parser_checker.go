package parser1

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"go_compiler/internal/utils"
)

type Notebook struct {
	Cells []Cell `json:"cells"`
}

type Cell struct {
	CellType string   `json:"cell_type"`
	Source   []string `json:"source"`
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

	// Start a persistent Python subprocess
	cmd := exec.Command("python3", "../../internal/executor/run_code.py")

	// Set OUTPUT_DIR from config
	if cfg.OutputDir != "" {
		cmd.Env = append(os.Environ(), "OUTPUT_DIR="+cfg.OutputDir)
	} else {
		log.Println("⚠️ Warning: Output directory not set in config. Using default current directory.")
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("Failed to get stdin: %v", err)
	}
	defer stdin.Close()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to get stdout: %v", err)
	}

	cmd.Stderr = os.Stderr // Optional: forward stderr to console

	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start Python process: %v", err)
	}

	reader := bufio.NewScanner(stdout)

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

			// Display results
			fmt.Printf("\n--- Cell %d Output ---\n", i)
			fmt.Println("Success:", result["success"])
			fmt.Println("Stdout:\n", result["stdout"])
			fmt.Println("Stderr:\n", result["stderr"])
			fmt.Println("----------------------")
		} else {
			log.Printf("No response from Python for cell %d", i)
		}
	}

	// Send exit command to Python subprocess
	exitPayload := map[string]string{"code": "__EXIT__"}
	exitJSON, _ := json.Marshal(exitPayload)
	stdin.Write(append(exitJSON, '\n'))

	cmd.Wait()
}
