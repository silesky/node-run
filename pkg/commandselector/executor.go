package commandselector

import (
	"fmt"
	"log"
	"node-task-runner/pkg/logger"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// InteractiveRunner represents the interactive command runner.
type InteractiveRunner struct {
	command string
	input   string
}

// NewInteractiveRunner creates a new InteractiveRunner.
func NewInteractiveRunner(command string) *InteractiveRunner {
	return &InteractiveRunner{
		command: command,
	}
}

// Init initializes the Bubble Tea program.
func (ir *InteractiveRunner) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model.
func (ir *InteractiveRunner) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return ir, tea.Quit
		case "r", "enter":
			ir.runCommand()
		}
	}
	return ir, nil
}

// View renders the UI.
func (ir *InteractiveRunner) View() string {
	const template = `
---------------------
Welcome to Interactive Mode.
Input command: %s

Commands:
  r - Rerun the command
  q - Quit the runner

%s`
	return fmt.Sprintf(template, ir.command, ir.input)
}

// runCommand executes a shell command.
func (ir *InteractiveRunner) runCommand() {
	command := ir.command
	logger.Debugf("Running CLI command: %s", command)
	parts := strings.Fields(command)
	if len(parts) == 0 {
		fmt.Println("No command provided")
		return
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		fmt.Println("Error:", err)
	}
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
	ir := NewInteractiveRunner(cmd)
	ir.command = cmd
	ir.runCommand()
	program := tea.NewProgram(ir)
	if _, err := program.Run(); err != nil {
		logger.Fatalf("%v", err)
	}

}
