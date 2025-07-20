package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"go_compiler/api"
	parser "go_compiler/parser1"
)

func CheckExtension(filePath string) {
	if len(filePath) == 0 {
		log.Println("Filename is empty")
		return
	}

	ext := filepath.Ext(filePath)
	fmt.Printf("Processing: %s (type: %s)\n", filePath, ext)

	switch ext {
	case ".py":
		fmt.Println("This is a Python file.")
		api.RunPythonChecker(filePath)
	case ".ipynb":
		fmt.Println("This is a Jupyter Notebook file.")
		parser.RunJupyterChecker(filePath)
	default:
		fmt.Println("Unknown file type. Skipping:", filePath)
	}
}

func processFile(path string) {
	info, err := os.Stat(path)
	if err != nil {
		log.Printf("Error checking file %s: %v", path, err)
		return
	}
	if info.IsDir() {
		log.Printf("Expected a file, but got a directory: %s", path)
		return
	}
	CheckExtension(path)
}

func processDirectory(dir string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("Failed to read directory %s: %v", dir, err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fullPath := filepath.Join(dir, file.Name())
		ext := filepath.Ext(file.Name())
		if ext == ".py" || ext == ".ipynb" {
			processFile(fullPath)
		}
	}
}

func main() {
	dirFlag := flag.String("dir", "", "Directory containing .py or .ipynb files")
	fileFlag := flag.String("file", "", "Single .py or .ipynb file to process")

	flag.Parse()

	if *dirFlag == "" && *fileFlag == "" {
		log.Fatal("Please provide either --dir <directory> or --file <filepath>")
	}
	if *dirFlag != "" && *fileFlag != "" {
		log.Fatal("Please provide only one of --dir or --file, not both.")
	}

	if *fileFlag != "" {
		processFile(*fileFlag)
	} else {
		processDirectory(*dirFlag)
	}
}
