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
		limit := visibleColumnTaskLimit(maxHeight)
		end := min(offset+limit, len(tasks))
		if offset > 0 {
			lines = append(lines, renderScrollIndicator("↑", offset, "above"))
		}
		for i, t := range tasks[offset:end] {
			taskIndex := offset + i
			lines = append(lines, RenderCard(t, inner, focused && taskIndex == focusedTask, routerState))
		}
		if end < len(tasks) {
			lines = append(lines, renderScrollIndicator("↓", len(tasks)-end, "more"))
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

func taskLabel(t data.Task) string {
	if t.Title == "" {
		return t.ID
	}
	return fmt.Sprintf("%s  %s", t.ID, t.Title)
}
