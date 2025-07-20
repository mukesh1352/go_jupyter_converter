package api

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Structs for parsing Jupyter notebook JSON
type Notebook struct {
	Cells []Cell `json:"cells"`
}

type Cell struct {
	CellType string   `json:"cell_type"`
	Source   []string `json:"source"`
}

// RunJupyterChecker reads and parses a Jupyter notebook (.ipynb)
func RunJupyterChecker(filepath string) {
	// Read the file content
	data, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println("Error reading the file:", err)
		os.Exit(1)
	}

	// Unmarshal the JSON content into the Notebook struct
	var nb Notebook
	if err := json.Unmarshal(data, &nb); err != nil {
		log.Fatalf("Error parsing the JSON file: %v", err)
	}

	// Loop through all cells in the notebook
	for i, cell := range nb.Cells {
		switch cell.CellType {
		case "code":
			fmt.Printf("Code cell %d:\n", i)
			for _, line := range cell.Source {
				fmt.Print(line)
			}
			fmt.Println("\n--n--")

		case "markdown":
			fmt.Printf("Markdown cell %d:\n", i)
			for _, line := range cell.Source {
				fmt.Print(line)
			}
			fmt.Println("\n--n--")

		default:
			fmt.Printf("Unknown cell type %q in cell %d\n", cell.CellType, i)
		}
	}
}
