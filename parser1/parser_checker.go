package parser1

import (
	"bytes"
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

	for i, cell := range nb.Cells {
		if cell.CellType != "code" {
			continue
		}

		code := ""
		for _, line := range cell.Source {
			code += line
		}

		// Call Python script
		input := map[string]string{"code": code}
		jsonData, _ := json.Marshal(input)

		cmd := exec.Command("python3", "executor/run_code.py")
		cmd.Stdin = bytes.NewBuffer(jsonData)

		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			log.Printf("Execution failed on cell %d: %v", i, err)
			continue
		}

		var result map[string]interface{}
		json.Unmarshal(out.Bytes(), &result)

		fmt.Printf(" Cell %d Output:\n", i)
		fmt.Println(" Success:", result["success"])
		fmt.Println(" Stdout:\n", result["stdout"])
		fmt.Println(" Stderr:\n", result["stderr"])
		fmt.Println("------")
	}
}
