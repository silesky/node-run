// MAYBE: rename to command_interactive.go
package commandselector

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// InteractivePackageCommandRunner represents the interactive command runner.
type InteractivePackageCommandRunner struct {
	command string
	input   string
}

// NewInteractiveRunner creates a new InteractiveRunner.
func NewInteractiveRunner(command string) *InteractivePackageCommandRunner {
	return &InteractivePackageCommandRunner{
		command: command,
	}
}

// Init initializes the Bubble Tea program.
func (ir *InteractivePackageCommandRunner) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model.
func (ir *InteractivePackageCommandRunner) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
func (ir *InteractivePackageCommandRunner) View() string {
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

func createCLICommand(proj Project, command Command) string {
	switch proj.Manager {
	case Npm:
		return fmt.Sprintf("npm --prefix %s run %s", command.PackageDir, command.CommandName)
	case Pnpm:
		return fmt.Sprintf("pnpm --prefix %s run %s", command.PackageDir, command.CommandName)
	case Yarn:
		return fmt.Sprintf("yarn --cwd %s %s", command.PackageDir, command.CommandName)
	default:
		log.Fatalf("invariant -- unknown package manager: %s", proj.Manager)
		return ""
	}
}

func (ir *InteractivePackageCommandRunner) runCommand() {
	command := ir.command
	parts := strings.Fields(command)
	if len(parts) == 0 {
		fmt.Println("No command provided")
		return
	}

	// the first part is the binary (e.g. npm), the rest are the args
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}

func Exec(command Command, project Project) {
	cmd := createCLICommand(project, command)
	ir := NewInteractiveRunner(cmd)
	ir.runCommand()
	program := tea.NewProgram(ir)
	if _, err := program.Run(); err != nil {
		log.Fatalf("%v", err)
	}
}
