package board

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/opencode/savepoint/internal/styles"
)

// RenderEpicDetail renders an overlay showing the content of an E##-Detail.md file.
func RenderEpicDetail(epicSlug, content string, overlayW, maxHeight, offset int) string {
	inner := overlayW - detailBorderPad
	if inner < 4 {
		inner = 4
	}

	lines := []string{
		styles.EpicTitleFocused.Render("EPIC DETAIL"),
		strings.Repeat("─", inner),
	}

	body := epicDetailBody(content, inner)
	body = append(body, "", styles.CardMeta.Render("esc:close"))
	lines = append(lines, visibleDetailLines(body, maxHeight-detailVerticalOverhead, offset)...)

	return styles.EpicDetailOverlay.Width(overlayW).Render(strings.Join(lines, "\n"))
}

// epicDetailBody parses markdown content into display lines, stripping frontmatter.
func epicDetailBody(content string, width int) []string {
	if strings.TrimSpace(content) == "" || content == "(no detail available)" {
		return []string{styles.CardMeta.Render("(no detail available)")}
	}

	lines := strings.Split(content, "\n")

	// Strip YAML frontmatter between leading --- markers.
	start := 0
	if len(lines) > 0 && strings.TrimSpace(lines[0]) == "---" {
		for i := 1; i < len(lines); i++ {
			if strings.TrimSpace(lines[i]) == "---" {
				start = i + 1
				break
			}
		}
	}
	lines = lines[start:]

	var body []string
	skip := false
	for _, line := range lines {
		trimmed := strings.TrimRight(line, " \t")
		if strings.HasPrefix(trimmed, "## ") {
			heading := strings.ToLower(strings.TrimPrefix(trimmed, "## "))
			skip = strings.Contains(heading, "component") || strings.Contains(heading, "files")
		}
		if skip {
			continue
		}
		switch {
		case strings.HasPrefix(trimmed, "# "):
			body = append(body, styles.EpicTitleFocused.Render(strings.TrimPrefix(trimmed, "# ")))
		case strings.HasPrefix(trimmed, "## "):
			body = append(body, "", styles.EpicItemFocused.Render(strings.TrimPrefix(trimmed, "## ")))
		case strings.HasPrefix(trimmed, "### "):
			body = append(body, styles.EpicItemFocused.Render(strings.TrimPrefix(trimmed, "### ")))
		case strings.HasPrefix(trimmed, "|"):
			body = append(body, styles.CardMeta.Render(trimmed))
		case trimmed == "":
			body = append(body, "")
		default:
			for _, wrapped := range WrapText(trimmed, width) {
				body = append(body, styles.CardMeta.Render(wrapped))
			}
		}
	}
	return body
}

const epicActiveMarker = "►"

// RenderEpicSidebar renders the fixed left sidebar listing epics with active indicator.
// If epics is empty and selected is non-empty, selected is shown as the sole entry.
func RenderEpicSidebar(epics []string, selected string, width int, focus bool, cursor int, status map[string]string) string {
	inner := width - epicPanelOverhead
	if inner < 2 {
		inner = 2
	}
	list := epics
	if len(list) == 0 && selected != "" {
		list = []string{selected}
	}

	title := styles.ColumnTitle.Render("EPICS")
	if focus {
		title = styles.EpicTitleFocused.Render("EPICS")
	}
	lines := []string{title, strings.Repeat("─", inner)}
	for i, e := range list {
		g := epicSidebarGlyph(status, e)
		gw := lipgloss.Width(g)
		if gw < 1 {
			gw = 1
		}
		label := truncate(e, inner-2-gw)
		if focus && len(epics) > 0 && i == cursor {
			lines = append(lines, styles.EpicItemFocused.Render(epicActiveMarker+" "+g+" "+label))
		} else if !focus && e == selected {
			lines = append(lines, styles.EpicItemFocused.Render(epicActiveMarker+" "+g+" "+label))
		} else {
			lines = append(lines, styles.TaskItem.Render("  "+g+" "+label))
		}
	}
	if len(list) == 0 {
		lines = append(lines, styles.TaskItem.Render("(none)"))
	}
	style := styles.EpicPanel.Width(width)
	if focus && len(epics) > 0 {
		style = styles.EpicPanelFocused.Width(width)
	}
	return style.Render(strings.Join(lines, "\n"))
}

func epicSidebarGlyph(status map[string]string, epicID string) string {
	if status == nil {
		return statusGlyphDefault
	}
	s, ok := status[epicID]
	if !ok {
		return statusGlyphDefault
	}
	return statusGlyph(s)
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
