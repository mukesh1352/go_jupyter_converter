package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"go_compiler/api"
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
		api.RunPythonChecker(filePath)
	default:
		fmt.Println("Unknown file type.")
	}
}

func main() {
	directory := "./files"
	filename := "1.py" 


	fullpath := filepath.Join(directory, filename)

	info, err := os.Stat(fullpath)
	if os.IsNotExist(err) {
		log.Fatal("The file does not exist in the specified directory...")
	}
	if info.IsDir() {
		log.Fatal("The specified path is a directory, not a file...")
	}

	CheckExtension(fullpath)
}
