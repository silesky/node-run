// MAYBE: rename to command_interactive.go
package commandselector

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// InteractivePackageCommandRunner represents the interactive command runner.
type InteractivePackageCommandRunner struct {
	command string
	escape  bool
}

// createInitialInteractiveRunnerModel creates a new InteractiveRunner.
func createInitialInteractiveRunnerModel(command string) *InteractivePackageCommandRunner {
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
		case "esc":
			ir.escape = true
			return ir, tea.Quit
		case "enter", "r":
			ir.runCommand()
			return ir, nil
		}
	}
	return ir, nil
}

func (ir *InteractivePackageCommandRunner) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(Colors.white)).
		Background(lipgloss.Color(Colors.purple)).
		Padding(1, 1)

	commandStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(Colors.purple)).
		Bold(true)

	helpCommandsStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(Colors.white)).
		Background(lipgloss.Color(Colors.purple)).
		Padding(0, 1)

	FooterCommandStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(Colors.charcoal))
	template := `
%s

%s

%s - Re-run the command
%s - Go back

%s`
	return fmt.Sprintf(template,
		titleStyle.Render("Interactive Mode"),
		commandStyle.Render(ir.command),
		helpCommandsStyle.Render("enter"),
		helpCommandsStyle.Render("esc"),
		FooterCommandStyle.Render("Press ctrl+c to quit."),
	)
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

var (
	ErrEscape = errors.New("user pressed escape")
)

func Exec(command Command, project Project) error {
	cmd := createCLICommand(project, command)
	ir := createInitialInteractiveRunnerModel(cmd)
	// initial run command
	ir.runCommand()

	// command runner
	program := tea.NewProgram(ir)
	if _, err := program.Run(); err != nil {
		log.Fatalf("%v", err)
	}
	if ir.escape {
		return ErrEscape
	}
	return nil
}
