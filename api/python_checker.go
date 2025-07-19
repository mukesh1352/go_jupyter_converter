package api

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func RunPythonChecker(filePath string) {
	startTime := time.Now()

	cmd := exec.Command("python3", filePath)
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
