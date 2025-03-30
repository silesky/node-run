package commandselector

import "github.com/charmbracelet/lipgloss"

type TerminalColors struct {
	magenta  lipgloss.AdaptiveColor
	charcoal lipgloss.AdaptiveColor
	blue     lipgloss.AdaptiveColor
	purple   lipgloss.AdaptiveColor
	white    lipgloss.AdaptiveColor
	yellow   lipgloss.AdaptiveColor
	black    lipgloss.AdaptiveColor
}

var (
	Colors = TerminalColors{
		black: lipgloss.AdaptiveColor{
			Light: "#000000", // Black for light backgrounds
			Dark:  "#ffffff", // White for dark backgrounds
		},
		white: lipgloss.AdaptiveColor{
			Light: "#000000", // Black for light backgrounds
			Dark:  "#ffffff", // White for dark backgrounds
		},
		magenta: lipgloss.AdaptiveColor{
			Light: "#8b005d", // Dark magenta for light backgrounds
			Dark:  "#ff79c6", // Bright magenta for dark backgrounds
		},
		charcoal: lipgloss.AdaptiveColor{
			Light: "#222222", // Dark gray for light backgrounds
			Dark:  "#dddddd", // Light gray for dark backgrounds
		},
		blue: lipgloss.AdaptiveColor{
			Light: "#003366", // Dark blue for light backgrounds
			Dark:  "#8be9fd", // Bright cyan-blue for dark backgrounds
		},
		purple: lipgloss.AdaptiveColor{
			Light: "#4b0082", // Dark purple for light backgrounds
			Dark:  "#bd93f9", // Soft purple for dark backgrounds
		},
		yellow: lipgloss.AdaptiveColor{
			Light: "#8b8000", // Dark yellow for light backgrounds
			Dark:  "#f1fa8c", // Bright yellow for dark backgrounds
		},
	}
)
