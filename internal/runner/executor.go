package runner

import (
	"fmt"
	"os"
	"os/exec"
)

// Execute runs a shortcut command in the user's shell.
func Execute(command string) error {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}

	cmd := exec.Command(shell, "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("execute command: %w", err)
	}
	return nil
}
