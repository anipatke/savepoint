package board

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/opencode/savepoint/internal/data"
)

func requireModel(t *testing.T, got tea.Model) Model {
	t.Helper()
	m, ok := got.(Model)
	if !ok {
		t.Fatalf("updated model type = %T, want board.Model", got)
	}
	return m
}

func TestUpdate_qQuits(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if cmd == nil {
		t.Fatal("expected tea.Quit cmd, got nil")
	}
}

func TestUpdate_ctrlCQuits(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Fatal("expected tea.Quit cmd, got nil")
	}
}

func TestUpdate_windowSizeUpdatesModel(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	got, cmd := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	updated := requireModel(t, got)
	if cmd != nil {
		t.Errorf("expected nil cmd for window resize, got %v", cmd)
	}
	if updated.Width != 120 {
		t.Errorf("Width = %d, want 120", updated.Width)
	}
	if updated.Height != 40 {
		t.Errorf("Height = %d, want 40", updated.Height)
	}
}

func TestUpdate_rightMovesColumn(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	// FocusedColumn starts at Planned
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")})
	updated := requireModel(t, got)
	if updated.FocusedColumn != data.ColumnInProgress {
		t.Errorf("FocusedColumn = %q, want %q", updated.FocusedColumn, data.ColumnInProgress)
	}
}

func TestUpdate_leftWrapsAround(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	// Planned -> left -> Done (wrap)
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("h")})
	updated := requireModel(t, got)
	if updated.FocusedColumn != data.ColumnDone {
		t.Errorf("FocusedColumn = %q, want %q", updated.FocusedColumn, data.ColumnDone)
	}
}

func TestUpdate_rightWrapsAround(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.FocusedColumn = data.ColumnDone
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")})
	updated := requireModel(t, got)
	if updated.FocusedColumn != data.ColumnPlanned {
		t.Errorf("FocusedColumn = %q, want %q", updated.FocusedColumn, data.ColumnPlanned)
	}
}

func TestUpdate_downMovesTaskFocus(t *testing.T) {
	tasks := []data.Task{
		{ID: "T1", Column: data.ColumnPlanned},
		{ID: "T2", Column: data.ColumnPlanned},
	}
	m := NewModel(tasks, "v1", "E03")
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	updated := requireModel(t, got)
	if updated.FocusedTask != 1 {
		t.Errorf("FocusedTask = %d, want 1", updated.FocusedTask)
	}
}

func TestUpdate_upMovesTaskFocus(t *testing.T) {
	tasks := []data.Task{
		{ID: "T1", Column: data.ColumnPlanned},
		{ID: "T2", Column: data.ColumnPlanned},
	}
	m := NewModel(tasks, "v1", "E03")
	m.FocusedTask = 1
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	updated := requireModel(t, got)
	if updated.FocusedTask != 0 {
		t.Errorf("FocusedTask = %d, want 0", updated.FocusedTask)
	}
}

func TestUpdate_taskFocusClampedAtEnd(t *testing.T) {
	tasks := []data.Task{{ID: "T1", Column: data.ColumnPlanned}}
	m := NewModel(tasks, "v1", "E03")
	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	updated := requireModel(t, got)
	if updated.FocusedTask != 0 {
		t.Errorf("FocusedTask = %d, want 0", updated.FocusedTask)
	}
}

func TestUpdate_unknownMsgNoOp(t *testing.T) {
	m := NewModel(nil, "v1", "E03")
	m.Width = 80
	got, cmd := m.Update(nil)
	updated := requireModel(t, got)
	if cmd != nil {
		t.Errorf("expected nil cmd, got %v", cmd)
	}
	if updated.Width != 80 {
		t.Errorf("Width changed unexpectedly: %d", updated.Width)
	}
}
