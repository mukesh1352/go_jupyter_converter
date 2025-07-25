package api

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func RunPythonChecker(filePath string) {
	startTime := time.Now()

	cmd := exec.Command("C:\\Users\\Albert Nedumudy\\AppData\\Local\\Programs\\Python\\Python311\\python.exe", filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error running Python script:", err)
	}

	duration := time.Since(startTime)
	fmt.Printf("\nExecution Time: %v\n", duration)
}
