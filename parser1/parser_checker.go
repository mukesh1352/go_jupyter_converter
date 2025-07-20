package parser1

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
)

type Notebook struct {
	Cells []Cell `json:"cells"`
}

type Cell struct {
	CellType string   `json:"cell_type"`
	Source   []string `json:"source"`
}

func RunJupyterChecker(filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading the file: %v", err)
	}

	var nb Notebook
	if err := json.Unmarshal(data, &nb); err != nil {
		log.Fatalf("Error parsing notebook JSON: %v", err)
	}

	// Start Python subprocess (persistent)
	cmd := exec.Command("python3", "executor/run_code.py")

	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("Failed to get stdin: %v", err)
	}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to get stdout: %v", err)
	}

	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start Python process: %v", err)
	}

	reader := bufio.NewScanner(stdoutPipe)

	for i, cell := range nb.Cells {
		if cell.CellType != "code" {
			continue
		}

		code := ""
		for _, line := range cell.Source {
			code += line
		}

		// Send code as JSON to Python
		input := map[string]string{"code": code}
		jsonData, _ := json.Marshal(input)
		_, err := stdinPipe.Write(append(jsonData, '\n')) 
		if err != nil {
			log.Printf("Failed to send code for cell %d: %v", i, err)
			continue
		}

		// Read JSON output line from Python
		if reader.Scan() {
			line := reader.Text()
			var result map[string]interface{}
			json.Unmarshal([]byte(line), &result)

			fmt.Printf(" Cell %d Output:\n", i)
			fmt.Println(" Success:", result["success"])
			fmt.Println(" Stdout:\n", result["stdout"])
			fmt.Println(" Stderr:\n", result["stderr"])
			fmt.Println("------")
		} else {
			log.Printf("No response from Python for cell %d", i)
		}
	}

	exitInput := map[string]string{"code": "__EXIT__"}
	exitJSON, _ := json.Marshal(exitInput)
	stdinPipe.Write(append(exitJSON, '\n'))
	stdinPipe.Close()
	cmd.Wait()
}
