package board

import (
	"strings"

	"github.com/opencode/savepoint/internal/styles"
)

const releaseActiveMarker = "►"

// RenderReleaseDropdown renders the release selection dropdown overlay.
func RenderReleaseDropdown(releases []string, cursor int, width int) string {
	inner := width - epicPanelOverhead
	if inner < 2 {
		inner = 2
	}

	lines := []string{styles.ColumnTitleFocused.Render("SELECT RELEASE"), strings.Repeat("─", inner)}
	for i, r := range releases {
		label := truncate(r, inner-2)
		if i == cursor {
			lines = append(lines, styles.TaskItemFocused.Render(releaseActiveMarker+" "+label))
		} else {
			lines = append(lines, styles.TaskItem.Render("  "+label))
		}
	}
	if len(releases) == 0 {
		lines = append(lines, styles.TaskItem.Render("(none)"))
	}
	lines = append(lines, "", styles.CardMeta.Render("↑↓:nav  enter:select  esc:cancel"))
	return styles.EpicPanel.Width(width).Render(strings.Join(lines, "\n"))
}


