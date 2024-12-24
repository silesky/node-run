// MAYBE: rename to command_interactive.go
package commandselector

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// InteractivePackageCommandRunner represents the interactive command runner.
type InteractivePackageCommandRunner struct {
	command     string
	escape      bool
	lastCommand CommandOutputMsg
}

// createInteractiveRunnerModel creates a new InteractiveRunner.
func createInteractiveRunnerModel(command string) *InteractivePackageCommandRunner {
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
		case "ctrl+c":
			return ir, tea.Quit
		case "esc":
			ir.escape = true
			return ir, tea.Quit
		case "enter", "r":
			return ir,
				func() tea.Msg {
					return ir.execCommand()
				}
		}
	case CommandOutputMsg:
		ir.lastCommand = msg
		return ir, nil
	}
	return ir, nil
}

func renderOutput(command CommandOutputMsg) string {
	getCurrentTime := func() string {
		now := time.Now()
		dateTimeString := now.Format("15:04:05")
		return dateTimeString
	}

	if command.Output == "" {
		return ""
	}

	timeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(Colors.purple)).
		Bold(true)

	timeContainerStyle := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(0, 1)

	time := timeContainerStyle.Render("Last executed: " + timeStyle.Render(getCurrentTime()))
	return "\n" + command.Output + "\n" + time
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

%s

%s - Re-run the command
%s - Go back

%s`
	return fmt.Sprintf(template,
		titleStyle.Render("Interactive Mode"),
		commandStyle.Render(ir.command),
		renderOutput(ir.lastCommand),
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

// CommandOutputMsg is a message that contains the output of a command.
type CommandOutputMsg struct {
	Output string
	Err    error
}

// execCommand creates a new exec.Cmd from a command string.
func execCommand(command string) *exec.Cmd {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		fmt.Fprintf(os.Stderr, "command string is empty")
	}
	cmd := exec.Command(parts[0], parts[1:]...)
	return cmd
}

func (ir *InteractivePackageCommandRunner) execCommand() CommandOutputMsg {
	cmd := execCommand(ir.command)
	output, err := cmd.CombinedOutput()
	return CommandOutputMsg{
		Output: string(output),
		Err:    err,
	}
}

var (
	ErrEscape = errors.New("user pressed escape")
)

func Exec(command Command, project Project) error {
	cmd := createCLICommand(project, command)
	// initial run command

	if command.ExecOptions.WithRunner {
		ir := createInteractiveRunnerModel(cmd)
		initialCommand := ir.execCommand()
		ir.lastCommand = initialCommand
		program := tea.NewProgram(ir)
		if _, err := program.Run(); err != nil {
			log.Fatalf("%v", err)
		}
		if ir.escape {
			return ErrEscape
		}
	} else {
		cmd := execCommand(cmd)
		// the first part is the binary (e.g. npm), the rest are the args
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			fmt.Println("Error:", err)
		}
	}
	return nil
}
