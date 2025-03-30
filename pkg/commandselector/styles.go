package commandselector

import "github.com/charmbracelet/lipgloss"

type TerminalColors struct {
	magenta lipgloss.AdaptiveColor
	grey    lipgloss.AdaptiveColor
	blue    lipgloss.AdaptiveColor
	purple  lipgloss.AdaptiveColor
	white   lipgloss.AdaptiveColor
	yellow  lipgloss.AdaptiveColor
	black   lipgloss.AdaptiveColor
}

// Light and Dark correspond to the color that will display based on the terminal background
var (
	Colors = TerminalColors{
		black: lipgloss.AdaptiveColor{
			Light: "#000000",
			Dark:  "#ffffff",
		},
		white: lipgloss.AdaptiveColor{
			Light: "#000000",
			Dark:  "#ffffff",
		},
		magenta: lipgloss.AdaptiveColor{
			Light: "#8b005d",
			Dark:  "#ff79c6",
		},
		grey: lipgloss.AdaptiveColor{
			Light: "#222222",
			Dark:  "#717171",
		},
		blue: lipgloss.AdaptiveColor{
			Light: "#003366",
			Dark:  "#8be9fd",
		},
		purple: lipgloss.AdaptiveColor{
			Light: "#4b0082",
			Dark:  "#bd93f9",
		},
		yellow: lipgloss.AdaptiveColor{
			Light: "#8b8000",
			Dark:  "#f1fa8c",
		},
	}
)
