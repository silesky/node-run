package fuzzsearch

import (
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
	quitting bool
	styles   Styles
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
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		// return the selected value and quit
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
		line := fmt.Sprintf("%s [%s] %s - %s", cursor, cmd.PackageName, cmd.Name, cmd.Command)
		line = strings.Replace(line, cursor, "", 1)
		if m.cursor == i {
			lines.WriteString(cursor + m.styles.selected.Render(line))
		} else {
			lines.WriteString(cursor + line)
		}

		lines.WriteString("\n")
	}

	// if no results, show no results found. otherwise, only show results message if filter is active.
	if len(m.filtered) == 0 {
		lines.WriteString("\n" + m.styles.helpText.Render("No results found.") + "\n")
	} else if len(m.filtered) != len(m.commands) {
		filterCommand := fmt.Sprintf("Displaying %d of %d results", len(m.filtered), len(m.commands))
		lines.WriteString("\n" + m.styles.helpText.Render(filterCommand))
	}

	quitHelp := m.styles.helpText.Render("Press q to quit.")
	lines.WriteString("\n\n" + quitHelp + "\n")
	return m.styles.container.Render(lines.String())
}

// filterCommands filters commands based on user input
func filterCommands(commands []Command, query string) []Command {
	var result []Command
	query = strings.TrimSpace(query)
	for _, cmd := range commands {
		if strings.Contains(strings.ToLower(cmd.PackageName), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(cmd.Name), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(cmd.Command), strings.ToLower(query)) {
			result = append(result, cmd)
		}
	}
	return result
}

type Styles struct {
	selected    lipgloss.Style
	filterInput lipgloss.Style
	container   lipgloss.Style
	helpText    lipgloss.Style
}

func newStyles() Styles {
	hotPink := lipgloss.Color("205")
	darkGray := lipgloss.Color("#A9A9A9")
	return Styles{
		selected:    lipgloss.NewStyle().Foreground(hotPink),
		filterInput: lipgloss.NewStyle().Foreground(hotPink),
		container:   lipgloss.NewStyle(),
		helpText:    lipgloss.NewStyle().Foreground(darkGray),
	}
}

// DisplayCommandSelector displays the command selector UI
func DisplayCommandSelector(commands []Command) (*Command, error) {
	ti := textinput.New()
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
		return nil, err
	}

	if m.cursor >= 0 && m.cursor < len(m.filtered) && !m.quitting {
		return &m.filtered[m.cursor], nil
	}

	return nil, fmt.Errorf("no command selected")
}
