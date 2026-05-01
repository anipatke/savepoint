package board

import (
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/data"
)

func TestRenderColumn_headerContainsLabel(t *testing.T) {
	got := RenderColumn(nil, data.ColumnPlanned, 30, 0, false)
	if !strings.Contains(got, "PLANNED") {
		t.Error("RenderColumn missing PLANNED label")
	}
}

func TestRenderColumn_headerContainsCount(t *testing.T) {
	tasks := []data.Task{{ID: "T1", Title: "Task one", Column: data.ColumnPlanned}}
	got := RenderColumn(tasks, data.ColumnPlanned, 30, 0, false)
	if !strings.Contains(got, "(1)") {
		t.Error("RenderColumn missing task count")
	}
}

func TestRenderColumn_emptyShowsPlaceholder(t *testing.T) {
	got := RenderColumn(nil, data.ColumnDone, 30, 0, false)
	if !strings.Contains(got, "(empty)") {
		t.Error("RenderColumn missing (empty) for empty column")
	}
}

func TestRenderColumn_focusedDoesNotPanic(t *testing.T) {
	tasks := []data.Task{{ID: "T1", Column: data.ColumnInProgress}}
	got := RenderColumn(tasks, data.ColumnInProgress, 30, 0, true)
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
		got := RenderColumn(nil, tc.col, 30, 0, false)
		if !strings.Contains(got, tc.label) {
			t.Errorf("RenderColumn missing label %q for col %q", tc.label, tc.col)
		}
	}
}

func TestRenderColumn_taskTitleRendered(t *testing.T) {
	tasks := []data.Task{{ID: "T2", Title: "Build it", Column: data.ColumnPlanned}}
	got := RenderColumn(tasks, data.ColumnPlanned, 30, 0, false)
	if !strings.Contains(got, "Build it") {
		t.Error("RenderColumn missing task title")
	}
}

func TestRenderColumn_rendersTaskCards(t *testing.T) {
	tasks := []data.Task{{ID: "T2", Title: "Build it", Column: data.ColumnPlanned, Stage: data.StageAudit}}
	got := RenderColumn(tasks, data.ColumnPlanned, 30, 0, true)
	if !strings.Contains(got, glyphAudit) {
		t.Error("RenderColumn should render task phase glyph from card")
	}
	if !strings.Contains(got, "╭") {
		t.Error("RenderColumn should render bordered card")
	}
}

func TestRenderColumn_emptyCountIsZero(t *testing.T) {
	got := RenderColumn(nil, data.ColumnPlanned, 30, 0, false)
	if !strings.Contains(got, "(0)") {
		t.Error("RenderColumn missing (0) count for empty column")
	}
}
