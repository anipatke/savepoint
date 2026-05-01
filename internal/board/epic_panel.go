package board

import (
	"strings"

	"github.com/opencode/savepoint/internal/styles"
)

const epicActiveMarker = "►"

// RenderEpicSidebar renders the fixed left sidebar listing epics with active indicator.
// If epics is empty and selected is non-empty, selected is shown as the sole entry.
func RenderEpicSidebar(epics []string, selected string, width int) string {
	inner := width - epicPanelOverhead
	if inner < 2 {
		inner = 2
	}
	list := epics
	if len(list) == 0 && selected != "" {
		list = []string{selected}
	}

	lines := []string{styles.ColumnTitle.Render("EPICS"), strings.Repeat("─", inner)}
	for _, e := range list {
		label := truncate(e, inner-2)
		if e == selected {
			lines = append(lines, styles.TaskItemFocused.Render(epicActiveMarker+" "+label))
		} else {
			lines = append(lines, styles.TaskItem.Render("  "+label))
		}
	}
	if len(list) == 0 {
		lines = append(lines, styles.TaskItem.Render("(none)"))
	}
	return styles.EpicPanel.Width(width).Render(strings.Join(lines, "\n"))
}

// RenderEpicDropdown renders the epic selection dropdown overlay.
func RenderEpicDropdown(epics []string, cursor int, width int) string {
	inner := width - epicPanelOverhead
	if inner < 2 {
		inner = 2
	}

	lines := []string{styles.ColumnTitleFocused.Render("SELECT EPIC"), strings.Repeat("─", inner)}
	for i, e := range epics {
		label := truncate(e, inner-2)
		if i == cursor {
			lines = append(lines, styles.TaskItemFocused.Render(epicActiveMarker+" "+label))
		} else {
			lines = append(lines, styles.TaskItem.Render("  "+label))
		}
	}
	if len(epics) == 0 {
		lines = append(lines, styles.TaskItem.Render("(none)"))
	}
	lines = append(lines, "", styles.CardMeta.Render("↑↓:nav  enter:select  esc:cancel"))
	return styles.EpicPanel.Width(width).Render(strings.Join(lines, "\n"))
}

// epicIndex returns the index of selected in epics, or 0 if not found.
func epicIndex(epics []string, selected string) int {
	for i, e := range epics {
		if e == selected {
			return i
		}
	}
	return 0
}
