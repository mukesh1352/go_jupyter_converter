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

func CheckExtension(filePath string, cfg *utils.Config) float64  {
	if filePath == "" {
		log.Println("Warning: provided filename is empty.")
		return 0
	}

	ext := filepath.Ext(filePath)
	log.Printf("Processing: %s (type: %s)", filePath, ext)

	switch ext {
	case ".py":
		api.RunPythonChecker(filePath)
		return 0 // or -1 if not relevant
	case ".ipynb":
		return parser1.RunJupyterChecker(filePath, cfg)
	default:
		log.Printf("Unsupported file type: %s. Skipping.", ext)
		return 0
	}
}

func processFile(path string, cfg *utils.Config) float64 {
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Printf("Error resolving absolute path for %s: %v", path, err)
		return 0
	}

	info, err := os.Stat(absPath)
	if err != nil {
		log.Printf("Error checking file %s: %v", absPath, err)
		return 0
	}
	if info.IsDir() {
		log.Printf("Expected file but got directory: %s", absPath)
		return 0
	} 

	return CheckExtension(absPath, cfg)
}

func processDirectory(dir string, cfg *utils.Config) float64 {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		log.Fatalf("Error resolving absolute directory path: %v", err)
	}

	files, err := os.ReadDir(absDir)
	if err != nil {
		log.Fatalf("Failed to read directory %s: %v", absDir, err)
	}

	var accuracy float64

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fullPath := filepath.Join(absDir, file.Name())
		switch filepath.Ext(file.Name()) {
		case ".py", ".ipynb":
			accuracy = processFile(fullPath, cfg)
		default:
			log.Printf("Skipping unsupported file: %s", fullPath)
		}
	}
	return accuracy
}

func main() {
	log.Printf("📘 Received file path")

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

	var accuracy float64
	if *fileFlag != "" {
		accuracy = processFile(*fileFlag, cfg)
	} else {
		accuracy = processDirectory(*dirFlag, cfg)
	}

	log.Printf("{\"accuracy\": %.2f}\n", accuracy)

	os.Exit(0)
}
	