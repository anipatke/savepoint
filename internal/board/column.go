package board

import (
	"fmt"
	"strings"

	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/styles"
)

// RenderColumn renders a board column: header with label+count, task list, bordered container.
func RenderColumn(tasks []data.Task, col data.ColumnType, width, focusedTask int, focused bool) string {
	inner := width - colOverhead
	if inner < minColWidth {
		inner = minColWidth
	}

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
		for i, t := range tasks {
			lines = append(lines, RenderCard(t, inner, focused && i == focusedTask))
		}
	}

	content := strings.Join(lines, "\n")
	st := styles.Column.Width(width)
	if focused {
		st = styles.ColumnFocused.Width(width)
	}
	return st.Render(content)
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
