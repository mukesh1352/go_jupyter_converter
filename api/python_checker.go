package api

import (
	"fmt"
	"os"
	"os/exec"
)

func RunPythonChecker(filePath string) {
	cmd := exec.Command("python3", filePath)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error running Python script:", err)
	}
}
