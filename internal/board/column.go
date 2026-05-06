package board

import (
	"fmt"
	"strings"

	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/styles"
)

// RenderColumn renders a board column: header with label+count, task viewport, bordered container.
func RenderColumn(tasks []data.Task, col data.ColumnType, width, maxHeight, offset, focusedTask int, focused bool, routerState *data.RouterState) string {
	inner := width - colOverhead
	if inner < minColWidth {
		inner = minColWidth
	}
	offset = clampViewportOffset(offset, len(tasks))

	title := columnTitle(col)
	header := fmt.Sprintf("%s (%d)", title, len(tasks))
	if focused {
		header = styles.ColumnTitleFocused.Render(header)
	} else {
		header = styles.ColumnTitle.Render(header)
	}

	lines := []string{header, strings.Repeat("─", inner)}
	if len(tasks) == 0 {
		lines = append(lines, styles.TaskItem.Render("(empty)"))
	} else {
		contentBudget := maxHeight - 2
		if contentBudget < 1 {
			contentBudget = 1
		}

		type cardEntry struct {
			card  string
			lines int
		}
		cardEntries := make([]cardEntry, 0, len(tasks)-offset)
		for i := offset; i < len(tasks); i++ {
			c := RenderCard(tasks[i], inner, focused && i == focusedTask, routerState)
			cardEntries = append(cardEntries, cardEntry{card: c, lines: strings.Count(c, "\n") + 1})
		}

		// Standard window: fit as many cards as possible from the start of cardEntries.
		reserveAbove := 0
		if offset > 0 {
			reserveAbove = 1
		}
		usedLines := reserveAbove
		endIdx := 0
		for endIdx < len(cardEntries) {
			hasMore := (offset + endIdx + 1) < len(tasks)
			bottomReserve := 0
			if hasMore {
				bottomReserve = 1
			}
			if usedLines+cardEntries[endIdx].lines+bottomReserve > contentBudget {
				break
			}
			usedLines += cardEntries[endIdx].lines
			endIdx++
		}
		if endIdx == 0 && len(cardEntries) > 0 {
			endIdx = 1
		}

		// Determine what portion of cardEntries to render.
		renderStart := 0
		renderEnd := endIdx

		focusedRelIdx := focusedTask - offset
		if focused && focusedRelIdx >= 0 && focusedRelIdx < len(cardEntries) && focusedRelIdx >= endIdx {
			// Focused task is beyond the standard window.
			// Anchor viewport at focused task: fill backward to use remaining budget.
			bottomCost := 0
			if focusedTask+1 < len(tasks) {
				bottomCost = 1
			}
			cardsLines := cardEntries[focusedRelIdx].lines
			newStart := focusedRelIdx
			for newStart > 0 {
				prev := cardEntries[newStart-1]
				topCost := 1
				if offset+newStart-1 == 0 {
					topCost = 0
				}
				if cardsLines+prev.lines+topCost+bottomCost > contentBudget {
					break
				}
				cardsLines += prev.lines
				newStart--
			}
			renderStart = newStart
			renderEnd = focusedRelIdx + 1
		}

		effectiveOffset := offset + renderStart
		if effectiveOffset > 0 {
			lines = append(lines, renderScrollIndicator("↑", effectiveOffset, "above"))
		}
		for i := renderStart; i < renderEnd; i++ {
			lines = append(lines, cardEntries[i].card)
		}
		if offset+renderEnd < len(tasks) {
			remaining := len(tasks) - (offset + renderEnd)
			lines = append(lines, renderScrollIndicator("↓", remaining, "more"))
		}
	}

	content := strings.Join(lines, "\n")
	st := styles.ColumnUnfocused.Width(width)
	if focused {
		st = styles.ColumnFocused.Width(width)
	}
	return st.Render(content)
}

func visibleColumnTaskLimit(maxHeight int) int {
	if maxHeight <= 0 {
		return 999999
	}
	limit := (maxHeight - 2) / 3
	if limit < 1 {
		return 1
	}
	return limit
}

func clampViewportOffset(offset, total int) int {
	if offset < 0 || total <= 0 {
		return 0
	}
	if offset >= total {
		return total - 1
	}
	return offset
}

func renderScrollIndicator(arrow string, count int, suffix string) string {
	return styles.ScrollIndicator.Render(fmt.Sprintf("%s %d %s", arrow, count, suffix))
}

func columnTitle(col data.ColumnType) string {
	switch col {
	case data.ColumnPlanned:
		return "PLANNED"
	case data.ColumnInProgress:
		return "IN PROGRESS"
	case data.ColumnDone:
		return "DONE"
	default:
		return strings.ToUpper(string(col))
	}
}
