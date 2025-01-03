package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/creack/pty"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go-run [command] [args...]")
		return
	}

	// Command to execute
	command := os.Args[1]

	// Create a pseudo-terminal
	ptyMaster, tty, err := pty.Open()
	if err != nil {
		fmt.Println("Error creating PTY:", err)
		return
	}
	defer ptyMaster.Close()
	defer tty.Close()

	// Build the shell script logic in a single command
	script := fmt.Sprintf(`printf "%s"; read k; exec %s`, strings.Join(os.Args[1:], " "), command)

	// Start a shell with the script
	cmd := exec.Command("sh", "-c", script)
	cmd.Stdin = tty
	cmd.Stdout = tty
	cmd.Stderr = tty

	// Start the command
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting command:", err)
		return
	}

	// Attach the PTY to the current terminal
	go func() {
		_, _ = os.Stdout.ReadFrom(ptyMaster)
	}()

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		fmt.Println("Error waiting for command:", err)
		return
	}
}
