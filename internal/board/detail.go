package board

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/styles"
)

const detailBorderPad = 4        // rounded border (2) + padding (2×1)
const detailVerticalOverhead = 4 // overlay border (2) + fixed title/header rows (2)

// RenderDetail renders a task detail overlay panel at the given display width.
// When router state matches t's release/epic/task, a "(router priority)" label is shown.
func RenderDetail(t data.Task, overlayW int, routerState *data.RouterState, maxHeight, offset int) string {
	inner := overlayW - detailBorderPad
	if inner < 4 {
		inner = 4
	}

	lines := []string{
		styles.ColumnTitleFocused.Render("TASK DETAIL"),
		strings.Repeat("─", inner),
	}
	body := []string{
		detailRow("ID", t.ID, inner),
		detailRow("Title", t.Title, inner),
		detailRow("Epic", t.Epic, inner),
		detailRow("Release", t.Release, inner),
		detailRow("Status", string(t.Column), inner),
		detailRow("Phase", phaseLabel(t.Stage), inner),
	}

	if t.Description != "" {
		body = append(body,
			"",
			styles.ColumnTitle.Render("Description:"),
		)
		for _, line := range WrapText(t.Description, inner) {
			body = append(body, styles.CardMeta.Render(line))
		}
	}

	if len(t.Checklist) > 0 {
		body = append(body, "", styles.ColumnTitle.Render("Implementation Plan:"), "")
		for _, item := range t.Checklist {
			glyph := "[ ] "
			style := styles.CardMeta
			if item.Done {
				glyph = "[x] "
				style = styles.TagDone
			}
			body = append(body, renderChecklistSentences(item.Text, glyph, inner, style)...)
		}
	}

	if t.Column != data.ColumnDone && isRouterPriority(t, routerState) {
		body = append(body, "", styles.TagDone.Render("(router priority)"))
	}
	body = append(body, "", styles.CardMeta.Render("esc:close"))
	lines = append(lines, visibleDetailLines(body, maxHeight-detailVerticalOverhead, offset)...)

	return styles.DetailOverlay.Width(overlayW).Render(strings.Join(lines, "\n"))
}

func renderChecklistSentences(text, glyph string, width int, style lipgloss.Style) []string {
	textWidth := width - len(glyph)
	if textWidth < 4 {
		textWidth = 4
	}

	lines := []string{}
	continuationIndent := strings.Repeat(" ", len(glyph))
	for _, sentence := range splitChecklistSentences(text) {
		wrapped := WrapText(sentence, textWidth)
		for i, line := range wrapped {
			if i == 0 {
				lines = append(lines, style.Render(glyph+line))
				continue
			}
			lines = append(lines, style.Render(continuationIndent+line))
		}
	}
	return lines
}

func splitChecklistSentences(text string) []string {
	fields := strings.Fields(text)
	if len(fields) == 0 {
		return nil
	}
	normalized := strings.Join(fields, " ")

	sentences := []string{}
	start := 0
	for i, r := range normalized {
		if r != '.' && r != '!' && r != '?' {
			continue
		}
		end := i + len(string(r))
		if end < len(normalized) && normalized[end] != ' ' {
			continue
		}
		sentence := strings.TrimSpace(normalized[start:end])
		if sentence != "" {
			sentences = append(sentences, sentence)
		}
		start = end
	}
	if tail := strings.TrimSpace(normalized[start:]); tail != "" {
		sentences = append(sentences, tail)
	}
	return sentences
}

func visibleDetailLines(lines []string, maxBodyHeight, offset int) []string {
	total := len(lines)
	if maxBodyHeight <= 0 || total <= maxBodyHeight {
		return lines
	}
	offset = clampDetailOffset(offset, total)
	available := maxBodyHeight
	if offset > 0 {
		available--
	}
	if available < 1 {
		available = 1
	}
	end := min(offset+available, total)
	if end < total && available > 1 {
		available--
		end = min(offset+available, total)
	}

	visible := make([]string, 0, available+2)
	if offset > 0 {
		visible = append(visible, renderScrollIndicator("↑", offset, "above"))
	}
	visible = append(visible, lines[offset:end]...)
	if end < total {
		visible = append(visible, renderScrollIndicator("↓", total-end, "more"))
	}
	return visible
}

func clampDetailOffset(offset, total int) int {
	if offset < 0 || total <= 0 {
		return 0
	}
	if offset >= total {
		return total - 1
	}
	return offset
}

func detailRow(label, value string, width int) string {
	prefix := label + ": "
	wrapped := WrapText(value, width-len(prefix))
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


