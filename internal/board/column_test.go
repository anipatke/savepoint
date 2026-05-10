package board

import (
	"fmt"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/opencode/savepoint/internal/data"
)

func TestRenderColumn_headerContainsLabel(t *testing.T) {
	got := RenderColumn(nil, data.ColumnPlanned, 30, 0, 0, 0, false, nil, nil)
	if !strings.Contains(got, "PLANNED") {
		t.Error("RenderColumn missing PLANNED label")
	}
}

func TestRenderColumn_headerContainsCount(t *testing.T) {
	tasks := []data.Task{{ID: "T1", Title: "Task one", Column: data.ColumnPlanned}}
	got := RenderColumn(tasks, data.ColumnPlanned, 30, 0, 0, 0, false, nil, nil)
	if !strings.Contains(got, "(1)") {
		t.Error("RenderColumn missing task count")
	}
}

func TestRenderColumn_emptyShowsPlaceholder(t *testing.T) {
	got := RenderColumn(nil, data.ColumnDone, 30, 0, 0, 0, false, nil, nil)
	if !strings.Contains(got, "(empty)") {
		t.Error("RenderColumn missing (empty) for empty column")
	}
}

func TestRenderColumn_focusedDoesNotPanic(t *testing.T) {
	tasks := []data.Task{{ID: "T1", Column: data.ColumnInProgress}}
	got := RenderColumn(tasks, data.ColumnInProgress, 30, 0, 0, 0, true, nil, nil)
	if got == "" {
		t.Error("RenderColumn returned empty string for focused column")
	}
}

func TestRenderColumn_allColumnTitles(t *testing.T) {
	cases := []struct {
		col   data.ColumnType
		label string
	}{
		{data.ColumnPlanned, "PLANNED"},
		{data.ColumnInProgress, "IN PROGRESS"},
		{data.ColumnDone, "DONE"},
	}
	for _, tc := range cases {
		got := RenderColumn(nil, tc.col, 30, 0, 0, 0, false, nil, nil)
		if !strings.Contains(got, tc.label) {
			t.Errorf("RenderColumn missing label %q for col %q", tc.label, tc.col)
		}
	}
}

func TestRenderColumn_taskTitleRendered(t *testing.T) {
	tasks := []data.Task{{ID: "T2", Title: "Build it", Column: data.ColumnPlanned}}
	got := RenderColumn(tasks, data.ColumnPlanned, 30, 0, 0, 0, false, nil, nil)
	if !strings.Contains(got, "Build it") {
		t.Error("RenderColumn missing task title")
	}
}

func TestRenderColumn_rendersTaskCards(t *testing.T) {
	tasks := []data.Task{{ID: "T2", Title: "Build it", Column: data.ColumnPlanned, Stage: data.StageAudit}}
	got := RenderColumn(tasks, data.ColumnPlanned, 30, 0, 0, 0, true, nil, nil)
	if !strings.Contains(got, glyphAudit) {
		t.Error("RenderColumn should render task phase glyph from card")
	}
	if !strings.Contains(got, "┌") {
		t.Error("RenderColumn should render focused card border")
	}
}

func TestRenderColumn_focusStatesUseStableBorderDimensions(t *testing.T) {
	tasks := []data.Task{{ID: "T2", Title: "Build it", Column: data.ColumnPlanned, Stage: data.StageAudit}}
	unfocused := RenderColumn(tasks, data.ColumnPlanned, 30, 0, 0, -1, false, nil, nil)
	focused := RenderColumn(tasks, data.ColumnPlanned, 30, 0, 0, -1, true, nil, nil)

	if !strings.Contains(unfocused, "┌") {
		t.Error("unfocused column should render a single-line border")
	}
	if !strings.Contains(focused, "┌") {
		t.Error("focused column should render a single-line border")
	}

	unfocusedLines := strings.Split(unfocused, "\n")
	focusedLines := strings.Split(focused, "\n")
	if len(unfocusedLines) != len(focusedLines) {
		t.Fatalf("line count changed between focus states: unfocused=%d focused=%d", len(unfocusedLines), len(focusedLines))
	}
	for i := range unfocusedLines {
		if lipgloss.Width(unfocusedLines[i]) != lipgloss.Width(focusedLines[i]) {
			t.Fatalf("line %d width changed between focus states: unfocused=%d focused=%d", i, lipgloss.Width(unfocusedLines[i]), lipgloss.Width(focusedLines[i]))
		}
	}
}

func TestRenderColumn_emptyCountIsZero(t *testing.T) {
	got := RenderColumn(nil, data.ColumnPlanned, 30, 0, 0, 0, false, nil, nil)
	if !strings.Contains(got, "(0)") {
		t.Error("RenderColumn missing (0) count for empty column")
	}
}

func TestRenderColumn_viewportShowsScrollIndicators(t *testing.T) {
	tasks := []data.Task{
		{ID: "T1", Title: "Task one", Column: data.ColumnPlanned},
		{ID: "T2", Title: "Task two", Column: data.ColumnPlanned},
		{ID: "T3", Title: "Task three", Column: data.ColumnPlanned},
		{ID: "T4", Title: "Task four", Column: data.ColumnPlanned},
	}

	got := RenderColumn(tasks, data.ColumnPlanned, 30, 8, 1, 1, true, nil, nil)

	if !strings.Contains(got, "↑ 1 above") {
		t.Error("RenderColumn missing above indicator")
	}
	if !strings.Contains(got, "↓ 2 more") {
		t.Errorf("RenderColumn missing more indicator, got:\n%s", got)
	}
	if strings.Contains(got, "Task one") {
		t.Error("RenderColumn should not render tasks above viewport")
	}
	if strings.Contains(got, "Task three") {
		t.Error("RenderColumn should not render tasks that don't fit budget")
	}
	if strings.Contains(got, "Task four") {
		t.Error("RenderColumn should not render tasks below viewport")
	}
}

func TestVisibleColumnTaskLimitDefaultsToFourAtStandardHeight(t *testing.T) {
	if got := visibleColumnTaskLimit(CalculateLayout(100, 24).ContentHeight); got != 4 {
		t.Errorf("visibleColumnTaskLimit(standard height) = %d, want 4", got)
	}
}

func BenchmarkRenderColumn_empty(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		RenderColumn(nil, data.ColumnPlanned, 30, 20, 0, 0, false, nil, nil)
	}
}

func BenchmarkRenderColumn_fewTasks(b *testing.B) {
	tasks := []data.Task{
		{ID: "E06/T001", Title: "First task", Column: data.ColumnPlanned, Stage: data.StageBuild},
		{ID: "E06/T002", Title: "Second task", Column: data.ColumnPlanned, Stage: data.StageTest},
		{ID: "E06/T003", Title: "Third task", Column: data.ColumnPlanned, Stage: data.StageAudit},
	}
	b.ReportAllocs()
	for b.Loop() {
		RenderColumn(tasks, data.ColumnPlanned, 30, 20, 0, 0, false, nil, nil)
	}
}

func BenchmarkRenderColumn_manyTasks(b *testing.B) {
	tasks := make([]data.Task, 20)
	stages := []data.ProgressStage{data.StageBuild, data.StageTest, data.StageAudit}
	for i := range tasks {
		tasks[i] = data.Task{
			ID:      fmt.Sprintf("E06/T%03d", i+1),
			Title:   fmt.Sprintf("Task number %d with a reasonable length title", i+1),
			Column:  data.ColumnPlanned,
			Stage:   stages[i%3],
			Release: "v1.1",
		}
	}
	b.ReportAllocs()
	for b.Loop() {
		RenderColumn(tasks, data.ColumnPlanned, 40, 24, 0, 0, false, nil, nil)
	}
}

func BenchmarkRenderColumn_focused(b *testing.B) {
	tasks := []data.Task{
		{ID: "E06/T001", Title: "First task", Column: data.ColumnInProgress, Stage: data.StageBuild},
		{ID: "E06/T002", Title: "Second task", Column: data.ColumnInProgress, Stage: data.StageTest},
		{ID: "E06/T003", Title: "Third task", Column: data.ColumnInProgress, Stage: data.StageAudit},
	}
	router := &data.RouterState{Release: "v1", Epic: "E06", Task: "T001"}
	b.ReportAllocs()
	for b.Loop() {
		RenderColumn(tasks, data.ColumnInProgress, 40, 20, 0, 0, true, router, nil)
	}
}

func TestRenderColumn_focusedLastTaskVisibleWhenScrolled(t *testing.T) {
	// When scrolled (offset > 0) the top scroll indicator steals 1 line.
	// The focused last-in-page task must still appear, not be cut off.
	tasks := []data.Task{
		{ID: "T1", Title: "Task one", Column: data.ColumnPlanned},
		{ID: "T2", Title: "Task two", Column: data.ColumnPlanned},
		{ID: "T3", Title: "Task three", Column: data.ColumnPlanned},
		{ID: "T4", Title: "Task four", Column: data.ColumnPlanned},
		{ID: "T5", Title: "Task five", Column: data.ColumnPlanned},
	}
	// offset=2, focusedTask=4 (last task), maxHeight=14 (24-line terminal)
	got := RenderColumn(tasks, data.ColumnPlanned, 30, 14, 2, 4, true, nil, nil)
	if !strings.Contains(got, "Task five") {
		t.Errorf("focused last task must be visible when scrolled, got:\n%s", got)
	}
	if !strings.Contains(got, "↑") {
		t.Error("scroll indicator above must appear when offset > 0")
	}
}
