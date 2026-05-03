package board

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// applyColorProfile detects terminal color support and configures lipgloss.
// Priority: NO_COLOR → COLORTERM env hint → termenv auto-detect.
func applyColorProfile() {
	lipgloss.SetColorProfile(detectProfile())
}

func detectProfile() termenv.Profile {
	if os.Getenv("NO_COLOR") != "" {
		return termenv.Ascii
	}
	if c := os.Getenv("COLORTERM"); c == "truecolor" || c == "24bit" {
		return termenv.TrueColor
	}
	return termenv.ColorProfile()
}
