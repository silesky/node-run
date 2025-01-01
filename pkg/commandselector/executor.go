package commandselector

import (
	"bufio"
	"fmt"
	"log"
	"node-task-runner/pkg/logger"
	"os"
	"os/exec"
	"strings"
)

// InteractiveRunner represents the interactive command runner.
type InteractiveRunner struct {
}

func NewInteractiveRunner() *InteractiveRunner {
	return &InteractiveRunner{}
}

// Start begins the interactive command execution loop.
func (ir *InteractiveRunner) Start(cmd string) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("\n-------------------------\nRerun the last command")
	fmt.Println("Commands:")
	fmt.Println("  r - Rerun the last command")
	fmt.Println("  q - Quit the runner")

	fmt.Println("Enter text (Ctrl+D to end):")
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}

	for scanner.Scan() {
		input := scanner.Text()
		// Execute the command
		if input == "r" {
			if err := runCommand(cmd); err != nil {
				fmt.Println("Error:", err)
			}
		}

		// Trim whitespace
		input = strings.TrimSpace(input)

		// Handle exit
		if input == "q" {
			fmt.Println("Exiting Interactive Runner.")
			break
		}

	}
}

// runCommand executes a shell command.
func runCommand(command string) error {
	logger.Debugf("Running CLI command: %s", command)
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("no command provided")
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func createCLICommand(proj Project, command Command) string {
	switch proj.Manager {
	case Npm:
		return fmt.Sprintf("npm --prefix %s run %s", command.PackageDir, command.CommandName)
	case Pnpm:
		return fmt.Sprintf("pnpm --prefix %s run %s", command.PackageDir, command.CommandName)
	case Yarn:
		return fmt.Sprintf("yarn --cwd %s %s", command.PackageDir, command.CommandName)
	default:
		log.Fatalf("invariant")
		return ""
	}
}
func Executor(command Command, project Project) {
	cmd := createCLICommand(project, command)
	if err := runCommand(cmd); err != nil {
		fmt.Println("Error:", err)
	}
	ir := NewInteractiveRunner()
	ir.Start(cmd)
}
