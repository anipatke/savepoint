package board

import (
	"strings"

	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/styles"
)

const detailBorderPad = 4 // rounded border (2) + padding (2×1)

// RenderDetail renders a task detail overlay panel at the given display width.
func RenderDetail(t data.Task, overlayW int) string {
	inner := overlayW - detailBorderPad
	if inner < 4 {
		inner = 4
	}

	lines := []string{
		styles.ColumnTitleFocused.Render("TASK DETAIL"),
		strings.Repeat("─", inner),
	}
	lines = append(lines,
		detailRow("ID", t.ID, inner),
		detailRow("Title", t.Title, inner),
		detailRow("Epic", t.Epic, inner),
		detailRow("Release", t.Release, inner),
		detailRow("Status", string(t.Column), inner),
		detailRow("Phase", phaseLabel(t.Stage), inner),
	)

	if t.Description != "" {
		lines = append(lines,
			"",
			styles.ColumnTitle.Render("Description:"),
		)
		for _, line := range wrapText(t.Description, inner) {
			lines = append(lines, styles.CardMeta.Render(line))
		}
	}

	if len(t.Acceptance) > 0 {
		lines = append(lines, "", styles.ColumnTitle.Render("Acceptance Criteria:"))
		for _, a := range t.Acceptance {
			for _, line := range wrapText(a, inner-2) {
				lines = append(lines, styles.CardMeta.Render("  • "+line))
			}
		}
	}

	if len(t.Checklist) > 0 {
		lines = append(lines, "", styles.ColumnTitle.Render("Implementation Plan:"))
		for _, item := range t.Checklist {
			for _, line := range wrapText(item, inner-2) {
				lines = append(lines, styles.CardMeta.Render("  □ "+line))
			}
		}
	}

	lines = append(lines, "", styles.CardMeta.Render("esc:close"))

	return styles.DetailOverlay.Width(overlayW).Render(strings.Join(lines, "\n"))
}

func detailRow(label, value string, width int) string {
	prefix := label + ": "
	wrapped := wrapText(value, width-len(prefix))
	if len(wrapped) == 0 {
		wrapped = []string{""}
	}
	lines := make([]string, 0, len(wrapped))
	for i, line := range wrapped {
		if i == 0 {
			lines = append(lines, styles.ColumnTitle.Render(prefix)+styles.CardMeta.Render(line))
			continue
		}
		lines = append(lines, strings.Repeat(" ", len(prefix))+styles.CardMeta.Render(line))
	}
	return strings.Join(lines, "\n")
}

func phaseLabel(s data.ProgressStage) string {
	switch s {
	case data.StageTest:
		return "test"
	case data.StageAudit:
		return "audit"
	default:
		return "build"
	}
}

func wrapText(s string, width int) []string {
	if width < 4 {
		width = 4
	}
	words := strings.Fields(s)
	if len(words) == 0 {
		return nil
	}
	lines := []string{}
	current := ""
	for _, word := range words {
		if len([]rune(word)) > width {
			if current != "" {
				lines = append(lines, current)
				current = ""
			}
			lines = append(lines, splitLongWord(word, width)...)
			continue
		}
		if current == "" {
			current = word
			continue
		}
		if len([]rune(current))+1+len([]rune(word)) <= width {
			current += " " + word
			continue
		}
		lines = append(lines, current)
		current = word
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

func splitLongWord(word string, width int) []string {
	runes := []rune(word)
	lines := []string{}
	for len(runes) > width {
		lines = append(lines, string(runes[:width]))
		runes = runes[width:]
	}
	if len(runes) > 0 {
		lines = append(lines, string(runes))
	}
	return lines
}
