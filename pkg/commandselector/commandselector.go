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
	commands     []Command
	filtered     []Command
	cursor       int
	input        textinput.Model
	styles       Styles
	quitting     bool
	runner       bool
	currentPage  int
	itemsPerPage int
}

// Init is required implementation
func (m *TeaCommandModel) Init() tea.Cmd {
	return m.input.Focus()
}

// Update is required implementation
func (m *TeaCommandModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		case "ctrl+r":
			m.runner = true
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
		case "left":
			if m.currentPage > 0 {
				m.currentPage--
				m.cursor = 0
			}
		case "right":
			if (m.currentPage+1)*m.itemsPerPage < len(m.filtered) {
				m.currentPage++
				m.cursor = 0
			}
		default:
			m.filtered = filterCommands(m.commands, m.input.Value())
			m.cursor = 0 // reset cursor to the first command
			m.input, cmd = m.input.Update(msg)
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

func (m TeaCommandModel) View() string {
	var lines strings.Builder

	lines.WriteString(m.styles.purple.Render("\nFilter: "))
	lines.WriteString(m.input.View())
	lines.WriteString("\n\n")

	start := m.currentPage * m.itemsPerPage
	end := start + m.itemsPerPage
	if end > len(m.filtered) {
		end = len(m.filtered)
	}

	for i, cmd := range m.filtered[start:end] {
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

	// Calculate the total number of pages
	totalPages := (len(m.filtered) + m.itemsPerPage - 1) / m.itemsPerPage
	currentPage := m.currentPage + 1

	// Display pagination information
	paginationInfo := fmt.Sprintf("Page %d of %d", currentPage, totalPages)
	lines.WriteString("\n" + m.styles.gray.Render(paginationInfo))

	// if no results, show no results found. otherwise, only show results message if filter is active.
	filterCommand := fmt.Sprintf("Displaying %d of %d results", len(m.filtered), len(m.commands))
	lines.WriteString("\n" + m.styles.gray.Render(filterCommand))

	quitHelp := renderQuit(m.styles, []string{
		"Press left/right to navigate pages.", "Press enter to run.", "Press ctrl+r to run interactively (beta).", "Press ctrl+c to quit.",
	})

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
	yellow    lipgloss.Style
}

func newStyles() Styles {
	// https://hexdocs.pm/color_palette/ansi_color_codes.html
	magenta := lipgloss.Color(Colors.magenta)
	charcoal := lipgloss.Color(Colors.charcoal)
	purple := lipgloss.Color(Colors.purple)
	yellow := lipgloss.Color(Colors.yellow)

	return Styles{
		magenta:   lipgloss.NewStyle().Foreground(magenta),
		container: lipgloss.NewStyle(),
		gray:      lipgloss.NewStyle().Foreground(charcoal),
		purple:    lipgloss.NewStyle().Foreground(purple),
		yellow:    lipgloss.NewStyle().Foreground(yellow),
	}
}

// Define custom errors
var (
	ErrUserAbort = errors.New("user aborted command selection")
)

// DisplayCommandSelector displays the command selector UI
func DisplayCommandSelector(commands []Command, initialInputValue string) (Command, error) {
	ti := textinput.New()
	ti.Placeholder = "Type to filter"
	ti.SetValue(initialInputValue)
	filtered := commands
	if initialInputValue != "" {
		filtered = filterCommands(commands, initialInputValue)
	}

	m := &TeaCommandModel{
		commands:     commands,
		filtered:     filtered,
		input:        ti,
		styles:       newStyles(),
		itemsPerPage: 50, // Set the number of items per page
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
