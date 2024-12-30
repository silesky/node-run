package fuzzsearch

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the Bubble Tea CommandModel
type CommandModel struct {
	commands    []Command
	filtered    []Command
	cursor      int
	input       string
	quitting    bool
	highlighted lipgloss.Style
}

func (m *CommandModel) Init() tea.Cmd {
	return nil
}

func (m *CommandModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			return m, tea.Quit

		case "up":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down":
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
			}

		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
				m.filtered = filterCommands(m.commands, m.input)
				if m.cursor >= len(m.filtered) {
					m.cursor = len(m.filtered) - 1
				}
			}

		default:
			m.input += msg.String()
			m.filtered = filterCommands(m.commands, m.input)
		}
	}

	return m, nil
}

// View renders the UI
func (m *CommandModel) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	var b strings.Builder
	b.WriteString("Search: " + m.input + "\n\n")

	for i, cmd := range m.filtered {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		line := fmt.Sprintf("%s [%s] %s - %s\n", cursor, cmd.PackageName, cmd.Name, cmd.Command)
		if m.cursor == i {
			b.WriteString(m.highlighted.Render(line))
		} else {
			b.WriteString(line)
		}
	}

	if len(m.filtered) == 0 {
		b.WriteString("No results found.\n")
	}

	return b.String()
}

// filterCommands filters commands based on user input
func filterCommands(commands []Command, query string) []Command {
	var result []Command
	for _, cmd := range commands {
		if strings.Contains(strings.ToLower(cmd.PackageName), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(cmd.Name), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(cmd.Command), strings.ToLower(query)) {
			result = append(result, cmd)
		}
	}
	return result
}

// commandSelector displays the command selector UI
func DisplayCommandSelector(commands []Command) (*Command, error) {
	m := &CommandModel{
		commands:    commands,
		filtered:    commands, // Start with all commands
		highlighted: lipgloss.NewStyle().Foreground(lipgloss.Color("205")),
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
