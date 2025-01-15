package commandselector

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TeaCommandModel represents the bubbletea Model API
type TeaCommandModel struct {
	commands []Command
	filtered []Command
	cursor   int
	input    textinput.Model
	styles   Styles
	quitting bool
	runner   bool
}

// Init is required implementation
func (m *TeaCommandModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update is required implementation
func (m *TeaCommandModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		// return the selected value and quit
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit

		// run the command in interactive mode
		case "ctrl+r":
			m.runner = true
			return m, tea.Quit

		// run the command
		case "enter":
			return m, tea.Quit

		// move the cursor up
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}

		// move the cursor down
		case "down":
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
			}

		case "backspace":
			m.input, cmd = m.input.Update(msg)
			m.filtered = filterCommands(m.commands, m.input.Value())
			m.cursor = 0 // reset cursor to the first command

		default:
			m.input, cmd = m.input.Update(msg)
			m.filtered = filterCommands(m.commands, m.input.Value())
			m.cursor = 0 // reset cursor to the first command
		}
	}

	return m, cmd
}

func renderQuit(styles Styles, lines []string) string {
	var message string
	for _, m := range lines {
		message += styles.gray.Render(m) + "\n"
	}
	return message
}

// View is required implementation for tea
func (m TeaCommandModel) View() string {
	var lines strings.Builder

	lines.WriteString("\nFilter: ")
	lines.WriteString(m.input.View())
	lines.WriteString("\n\n")

	for i, cmd := range m.filtered {
		cursor := "  "
		if m.cursor == i {
			cursor = "> "
		}

		// Format the line for the current command
		line := fmt.Sprintf("%s [%s] %s (%s)",
			cursor, cmd.PackageName, m.styles.purple.Render(cmd.CommandName), m.styles.gray.Render(cmd.CommandValue))
		line = strings.Replace(line, cursor, "", 1)
		if m.cursor == i {
			lines.WriteString(cursor + m.styles.magenta.Render(line))
		} else {
			lines.WriteString(cursor + line)
		}

		lines.WriteString("\n")
	}

	// if no results, show no results found. otherwise, only show results message if filter is active.
	if len(m.filtered) == 0 {
		lines.WriteString("\n" + m.styles.gray.Render("No results found.") + "\n")
	} else if len(m.filtered) != len(m.commands) {
		filterCommand := fmt.Sprintf("Displaying %d of %d results", len(m.filtered), len(m.commands))
		lines.WriteString("\n" + m.styles.gray.Render(filterCommand))
	}

	quitHelp := renderQuit(m.styles, []string{
		"Press enter to run", "Press ctrl+r to run interactively.", "Press ctrl+c to quit."})

	lines.WriteString("\n\n" + quitHelp + "\n")
	return m.styles.container.Render(lines.String())
}

// filterCommands filters commands based on user input
func filterCommands(commands []Command, query string) []Command {
	// create an _empty_ slice, not a nil slice (var []Command{} creates a nil slice)
	// https://stackoverflow.com/questions/49104157/what-is-the-point-of-having-nil-slice-and-empty-slice-in-golang
	result := []Command{}
	query = strings.TrimSpace(query)
	tokens := strings.Fields(strings.ToLower(query)) // Split query into tokens

	for _, cmd := range commands {
		combinedFields := strings.ToLower(cmd.PackageName) + " " +
			strings.ToLower(cmd.CommandName) + " " +
			strings.ToLower(cmd.CommandValue)

		matches := true
		for _, token := range tokens {
			if !strings.Contains(combinedFields, token) {
				matches = false
				break
			}
		}

		if matches {
			result = append(result, cmd)
		}
	}

	return result
}

type Styles struct {
	magenta   lipgloss.Style
	container lipgloss.Style
	gray      lipgloss.Style
	purple    lipgloss.Style
}

func newStyles() Styles {
	// https://hexdocs.pm/color_palette/ansi_color_codes.html
	magenta := lipgloss.Color(Colors.magenta)
	charcoal := lipgloss.Color(Colors.charcoal)
	purple := lipgloss.Color(Colors.purple)

	return Styles{
		magenta:   lipgloss.NewStyle().Foreground(magenta),
		container: lipgloss.NewStyle(),
		gray:      lipgloss.NewStyle().Foreground(charcoal),
		purple:    lipgloss.NewStyle().Foreground(purple),
	}
}

// Define custom errors
var (
	ErrUserAbort = errors.New("user aborted command selection")
)

// DisplayCommandSelector displays the command selector UI
func DisplayCommandSelector(commands []Command, initialInputValue string) (Command, error) {
	ti := textinput.New()
	ti.SetValue(initialInputValue)
	ti.Placeholder = "Type to filter"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	m := &TeaCommandModel{
		commands: commands,
		filtered: commands, // Start with all commands
		input:    ti,
		styles:   newStyles(),
	}

	program := tea.NewProgram(m)

	if _, err := program.Run(); err != nil {
		return Command{}, err
	}

	if m.cursor >= 0 && m.cursor < len(m.filtered) && !m.quitting {
		cmd := m.filtered[m.cursor]
		cmd.ExecOptions.WithRunner = m.runner
		return cmd, nil
	}

	if m.quitting {
		return Command{}, ErrUserAbort
	}

	return Command{}, fmt.Errorf("no command selected")
}
