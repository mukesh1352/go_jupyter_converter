package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"go_compiler/internal/api"
	"go_compiler/internal/parser1"
	"go_compiler/internal/utils"
)

func CheckExtension(filePath string, cfg *utils.Config) {
	if filePath == "" {
		log.Println("Warning: provided filename is empty.")
		return
	}

	ext := filepath.Ext(filePath)
	log.Printf("Processing: %s (type: %s)", filePath, ext)

	switch ext {
	case ".py":
		log.Println("Detected Python file.")
		api.RunPythonChecker(filePath)
	case ".ipynb":
		log.Println("Detected Jupyter Notebook file.")
		parser1.RunJupyterChecker(filePath, cfg) // ✅ Pass full config
	default:
		log.Printf("Unsupported file type: %s. Skipping.", ext)
	}
}

func processFile(path string, cfg *utils.Config) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Printf("Error resolving absolute path for %s: %v", path, err)
		return
	}

	info, err := os.Stat(absPath)
	if err != nil {
		log.Printf("Error checking file %s: %v", absPath, err)
		return
	}
	if info.IsDir() {
		log.Printf("Expected file but got directory: %s", absPath)
		return
	}

	CheckExtension(absPath, cfg)
}

func processDirectory(dir string, cfg *utils.Config) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		log.Fatalf("Error resolving absolute directory path: %v", err)
	}

	files, err := os.ReadDir(absDir)
	if err != nil {
		log.Fatalf("Failed to read directory %s: %v", absDir, err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fullPath := filepath.Join(absDir, file.Name())
		switch filepath.Ext(file.Name()) {
		case ".py", ".ipynb":
			processFile(fullPath, cfg)
		default:
			log.Printf("Skipping unsupported file: %s", fullPath)
		}
	}
}

func main() {
	dirFlag := flag.String("dir", "", "Directory containing .py or .ipynb files")
	fileFlag := flag.String("file", "", "Single .py or .ipynb file to process")
	flag.Parse()

	if *dirFlag == "" && *fileFlag == "" {
		log.Fatal("Error: Please provide either --dir <directory> or --file <filepath>")
	}
	if *dirFlag != "" && *fileFlag != "" {
		log.Fatal("Error: Please provide only one of --dir or --file at a time")
	}

	cfg, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if *fileFlag != "" {
		processFile(*fileFlag, cfg)
	} else {
		processDirectory(*dirFlag, cfg)
	}
}
	