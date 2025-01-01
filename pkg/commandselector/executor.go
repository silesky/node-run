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

// InteractiveRunner manages the interactive command execution loop.
type InteractiveRunner struct {
	lastCommand string
}

// NewInteractiveRunner initializes a new InteractiveRunner instance.
func NewInteractiveRunner() *InteractiveRunner {
	return &InteractiveRunner{}
}

// Start begins the interactive command execution loop.
func (ir *InteractiveRunner) Start() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to the Interactive Runner!")
	fmt.Println("Commands:")
	fmt.Println("  r - Rerun the last command")
	fmt.Println("  exit - Exit the runner")

	for {
		fmt.Print("runner> ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		// Trim whitespace
		input = strings.TrimSpace(input)

		// Handle exit
		if input == "exit" {
			fmt.Println("Exiting Interactive Runner.")
			break
		}

		// Handle rerun last command
		if input == "r" {
			if ir.lastCommand == "" {
				fmt.Println("No command to rerun.")
				continue
			}
			input = ir.lastCommand
		} else {
			ir.lastCommand = input
		}

		// Execute the command
		if err := ir.runCommand(input); err != nil {
			fmt.Println("Error:", err)
		}
	}
}

// runCommand executes a shell command.
func (ir *InteractiveRunner) runCommand(command string) error {
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
	NewInteractiveRunner().runCommand(cmd)
}
