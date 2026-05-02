package board

import (
	"os"
	"path/filepath"
	"strings"
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

func TestUpdate_downAutoScrollsFocusedTaskIntoViewport(t *testing.T) {
	tasks := []data.Task{
		{ID: "T1", Column: data.ColumnPlanned},
		{ID: "T2", Column: data.ColumnPlanned},
		{ID: "T3", Column: data.ColumnPlanned},
		{ID: "T4", Column: data.ColumnPlanned},
		{ID: "T5", Column: data.ColumnPlanned},
	}
	m := NewModel(tasks, "v1", "E03")
	m.Width = 100
	m.Height = 24
	m.FocusedTask = 3

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	updated := requireModel(t, got)

	if updated.FocusedTask != 4 {
		t.Errorf("FocusedTask = %d, want 4", updated.FocusedTask)
	}
	if updated.ColumnOffsets[data.ColumnPlanned] != 1 {
		t.Errorf("ColumnOffsets[planned] = %d, want 1", updated.ColumnOffsets[data.ColumnPlanned])
	}
}

func TestUpdate_pageDownScrollsFocusedColumnByPage(t *testing.T) {
	tasks := []data.Task{
		{ID: "T1", Column: data.ColumnPlanned},
		{ID: "T2", Column: data.ColumnPlanned},
		{ID: "T3", Column: data.ColumnPlanned},
		{ID: "T4", Column: data.ColumnPlanned},
		{ID: "T5", Column: data.ColumnPlanned},
		{ID: "T6", Column: data.ColumnPlanned},
	}
	m := NewModel(tasks, "v1", "E03")
	m.Width = 100
	m.Height = 24

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyPgDown})
	updated := requireModel(t, got)

	if updated.ColumnOffsets[data.ColumnPlanned] != 2 {
		t.Errorf("ColumnOffsets[planned] = %d, want 2", updated.ColumnOffsets[data.ColumnPlanned])
	}
	if updated.FocusedTask != 2 {
		t.Errorf("FocusedTask = %d, want 2", updated.FocusedTask)
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

func TestUpdate_mSetsRouterToFocusedTask(t *testing.T) {
	root := writeRouterFixture(t)
	tasks := []data.Task{
		{ID: "E05-tasking-permissions/T004-implement-m-hotkey", Release: "v1.1", Epic: "E05-tasking-permissions", Column: data.ColumnPlanned},
		{ID: "E05-tasking-permissions/T005-update-help-overlay", Release: "v1.1", Epic: "E05-tasking-permissions", Column: data.ColumnPlanned},
	}
	m := NewModel(tasks, "v1.1", "E05-tasking-permissions")
	m.Root = root

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("m")})
	updated := requireModel(t, got)

	if !strings.Contains(updated.StatusMessage, "Router set to v1.1 E05-tasking-permissions/T004") {
		t.Fatalf("StatusMessage = %q", updated.StatusMessage)
	}
	state := readRouterFixture(t, root)
	if state.State != "task-building" {
		t.Errorf("router state = %q, want task-building", state.State)
	}
	if state.Release != "v1.1" {
		t.Errorf("router release = %q, want v1.1", state.Release)
	}
	if state.Epic != "E05-tasking-permissions" {
		t.Errorf("router epic = %q, want E05-tasking-permissions", state.Epic)
	}
	if state.Task != "E05-tasking-permissions/T004-implement-m-hotkey" {
		t.Errorf("router task = %q, want focused task", state.Task)
	}
}

func TestUpdate_mSetsAuditPendingForLastUncompletedTask(t *testing.T) {
	root := writeRouterFixture(t)
	tasks := []data.Task{
		{ID: "E05-tasking-permissions/T004-implement-m-hotkey", Release: "v1.1", Epic: "E05-tasking-permissions", Column: data.ColumnPlanned},
		{ID: "E05-tasking-permissions/T003-update-design-md", Release: "v1.1", Epic: "E05-tasking-permissions", Column: data.ColumnDone},
	}
	m := NewModel(tasks, "v1.1", "E05-tasking-permissions")
	m.Root = root

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("m")})
	updated := requireModel(t, got)

	if updated.StatusMessage != "Audit pending for E05-tasking-permissions" {
		t.Fatalf("StatusMessage = %q", updated.StatusMessage)
	}
	state := readRouterFixture(t, root)
	if state.State != "audit-pending" {
		t.Errorf("router state = %q, want audit-pending", state.State)
	}
	if state.Task != "" {
		t.Errorf("router task = %q, want empty", state.Task)
	}
}

func TestUpdate_mDoesNothingWhenOverlayOpen(t *testing.T) {
	root := writeRouterFixture(t)
	tasks := []data.Task{{ID: "E05-tasking-permissions/T004-implement-m-hotkey", Release: "v1.1", Epic: "E05-tasking-permissions", Column: data.ColumnPlanned}}
	m := NewModel(tasks, "v1.1", "E05-tasking-permissions")
	m.Root = root
	m.Overlay = OverlayHelp

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("m")})
	updated := requireModel(t, got)

	if updated.StatusMessage != "" {
		t.Fatalf("StatusMessage = %q, want empty", updated.StatusMessage)
	}
	state := readRouterFixture(t, root)
	if state.Task != "T001" {
		t.Errorf("router task = %q, want unchanged T001", state.Task)
	}
}

func TestUpdate_mDoesNotSetDoneTask(t *testing.T) {
	root := writeRouterFixture(t)
	tasks := []data.Task{{ID: "E05-tasking-permissions/T004-implement-m-hotkey", Release: "v1.1", Epic: "E05-tasking-permissions", Column: data.ColumnDone}}
	m := NewModel(tasks, "v1.1", "E05-tasking-permissions")
	m.Root = root
	m.FocusedColumn = data.ColumnDone

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("m")})
	updated := requireModel(t, got)

	if updated.StatusMessage != "Router not updated: focused task is done" {
		t.Fatalf("StatusMessage = %q", updated.StatusMessage)
	}
	state := readRouterFixture(t, root)
	if state.Task != "T001" {
		t.Errorf("router task = %q, want unchanged T001", state.Task)
	}
}

func writeRouterFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	content := "# Agent State Machine\n\n## Current state\n\n```yaml\nstate: task-building\nrelease: v1.1\nepic: E05\ntask: T001\nnext_action: Build T001.\n```\n"
	if err := os.WriteFile(filepath.Join(root, "router.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return root
}

func readRouterFixture(t *testing.T, root string) *data.RouterState {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(root, "router.md"))
	if err != nil {
		t.Fatal(err)
	}
	state, err := data.NewRouterReader().ReadState(string(content))
	if err != nil {
		t.Fatal(err)
	}
	return state
}
