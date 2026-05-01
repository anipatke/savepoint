package board

import (
	"testing"

	"github.com/opencode/savepoint/internal/data"
)

func TestNewModel_emptyTasks(t *testing.T) {
	m := NewModel(nil, "v1", "E03")

	if m.SelectedRelease != "v1" {
		t.Errorf("SelectedRelease = %q, want %q", m.SelectedRelease, "v1")
	}
	if m.SelectedEpic != "E03" {
		t.Errorf("SelectedEpic = %q, want %q", m.SelectedEpic, "E03")
	}
	if m.FocusedColumn != data.ColumnPlanned {
		t.Errorf("FocusedColumn = %q, want %q", m.FocusedColumn, data.ColumnPlanned)
	}
	if m.FocusedTask != 0 {
		t.Errorf("FocusedTask = %d, want 0", m.FocusedTask)
	}
	if m.Overlay != OverlayNone {
		t.Errorf("Overlay = %q, want empty", m.Overlay)
	}
	for _, col := range []data.ColumnType{data.ColumnPlanned, data.ColumnInProgress, data.ColumnDone} {
		if _, ok := m.Tasks[col]; !ok {
			t.Errorf("Tasks missing column %q", col)
		}
	}
}

func TestNewModel_groupsByColumn(t *testing.T) {
	tasks := []data.Task{
		{ID: "T1", Column: data.ColumnPlanned},
		{ID: "T2", Column: data.ColumnInProgress},
		{ID: "T3", Column: data.ColumnDone},
		{ID: "T4", Column: data.ColumnPlanned},
	}
	m := NewModel(tasks, "v1", "E03")

	if got := len(m.Tasks[data.ColumnPlanned]); got != 2 {
		t.Errorf("Planned count = %d, want 2", got)
	}
	if got := len(m.Tasks[data.ColumnInProgress]); got != 1 {
		t.Errorf("InProgress count = %d, want 1", got)
	}
	if got := len(m.Tasks[data.ColumnDone]); got != 1 {
		t.Errorf("Done count = %d, want 1", got)
	}
}

func TestNewModel_emptyColumnDefaultsToPlanned(t *testing.T) {
	tasks := []data.Task{{ID: "T1"}}
	m := NewModel(tasks, "v1", "E03")

	if got := len(m.Tasks[data.ColumnPlanned]); got != 1 {
		t.Errorf("Planned count = %d, want 1 (empty column defaulted)", got)
	}
}

func TestModel_Init(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	// Init must not panic and returns a valid Cmd (nil or batch).
	_ = m.Init()
}
