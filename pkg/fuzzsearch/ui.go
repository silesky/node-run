package fuzzsearch

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TeaCommandModel represents the bubbletea Model API
type TeaCommandModel struct {
	commands []Command
	filtered []Command
	cursor   int
	input    string
	quitting bool
	styles   Styles
}

// Init is required implementation
func (m *TeaCommandModel) Init() tea.Cmd {
	return nil
}

// Init is required implementation
func (m *TeaCommandModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		// return the selected value and quit
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		// return the selected varlua and quit
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

		case "left", "right":
			// do nothing

		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
				m.filtered = filterCommands(m.commands, m.input)
				if m.cursor >= len(m.filtered) {
					m.cursor = len(m.filtered) - 1
				}
			}

		default:
			// if filtering using the search box
			m.input += msg.String()
			m.filtered = filterCommands(m.commands, m.input)
		}
	}

	return m, nil
}

// View is required implementation for tea
func (m TeaCommandModel) View() string {
	var lines strings.Builder

	lines.WriteString("\n Filter: ")
	lines.WriteString(m.styles.filterInput.Render(m.input))
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

	if len(m.filtered) == 0 {
		lines.WriteString("No results found.\n")
	}

	lines.WriteString("\nPress q to quit.\n")
	return lines.String()
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
}

func newStyles() Styles {
	hotPink := lipgloss.Color("205")
	return Styles{
		selected:    lipgloss.NewStyle().Foreground(hotPink),
		filterInput: lipgloss.NewStyle().Foreground(hotPink),
	}
}

// commandSelector displays the command selector UI
func DisplayCommandSelector(commands []Command) (*Command, error) {
	m := &TeaCommandModel{
		commands: commands,
		filtered: commands, // Start with all commands
		styles:   newStyles(),
	}

	program := tea.NewProgram(m)
	if _, err := program.Run(); err != nil {
		return nil, err
	}

	if m.cursor >= 0 && m.cursor < len(m.filtered) {
		return &m.filtered[m.cursor], nil
	}

	return nil, fmt.Errorf("no command selected")
}
