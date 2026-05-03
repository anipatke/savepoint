package board

import (
	"strings"

	"github.com/opencode/savepoint/internal/styles"
)

const helpBorderPad = 4 // rounded border (2) + padding (2x1)

// RenderHelp renders the keyboard shortcut reference overlay.
func RenderHelp(width int) string {
	inner := width - helpBorderPad
	if inner < 4 {
		inner = 4
	}

	lines := []string{
		styles.ColumnTitleFocused.Render("KEYBOARD SHORTCUTS"),
		strings.Repeat("─", inner),
		helpRow("h / left", "previous column"),
		helpRow("l / right", "next column"),
		helpRow("enter", "open task detail / select item"),
		helpRow("e", "open epic selector on narrow screens"),
		helpRow("r", "open release selector"),
		helpRow("p", "mark focused task as priority"),
		helpRow("up / k", "move selector up"),
		helpRow("down / j", "move selector down"),
		helpRow("?", "open help"),
		helpRow("esc / q", "close overlay"),
		helpRow("q / ctrl+c", "quit from board"),
		"",
		styles.CardMeta.Render("esc/q:close"),
	}

	return styles.DetailOverlay.Width(width).Render(strings.Join(lines, "\n"))
}

func helpRow(key, action string) string {
	return styles.ColumnTitle.Render(key+": ") + styles.CardMeta.Render(action)
}
