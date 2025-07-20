package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"go_compiler/api"
	parser "go_compiler/parser1" 
)


func CheckExtension(filePath string) {
	if len(filePath) == 0 {
		log.Fatal("Filename is empty")
	}

	ext := filepath.Ext(filePath)
	fmt.Println("File extension is:", ext)

	switch ext {
	case ".py":
		fmt.Println("This is a Python file.")
		api.RunPythonChecker(filePath)
	case ".ipynb":
		fmt.Println("This is a Jupyter Notebook file.")
		parser.RunJupyterChecker(filePath)
	default:
		fmt.Println("Unknown file type.")
	}
}

func main() {
	directory := "./files"
	filename := "1.ipynb"

	fullPath := filepath.Join(directory, filename)

	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("The file does not exist at the specified path: %s", fullPath)
		}
		log.Fatalf("Error checking file: %v", err)
	}

	if info.IsDir() {
		log.Fatalf("The specified path is a directory, not a file: %s", fullPath)
	}

	CheckExtension(fullPath)
}
