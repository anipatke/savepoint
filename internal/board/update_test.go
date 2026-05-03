package board

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

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

func processCmd(t *testing.T, m Model, cmd tea.Cmd) Model {
	t.Helper()
	msg := cmd()
	got, _ := m.Update(msg)
	return requireModel(t, got)
}

func TestUpdate_pSetsRouterToFocusedTask(t *testing.T) {
	root := writeRouterFixture(t)
	tasks := []data.Task{
		{ID: "E05-tasking-permissions/T004-implement-m-hotkey", Release: "v1.1", Epic: "E05-tasking-permissions", Column: data.ColumnPlanned},
		{ID: "E05-tasking-permissions/T005-update-help-overlay", Release: "v1.1", Epic: "E05-tasking-permissions", Column: data.ColumnPlanned},
	}
	m := NewModel(tasks, "v1.1", "E05-tasking-permissions")
	m.Root = root

	got, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("p")})
	first := requireModel(t, got)
	updated := processCmd(t, first, cmd)

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

func TestUpdate_pSetsRouterToFocusedTaskWhenItIsLastUncompleted(t *testing.T) {
	root := writeRouterFixture(t)
	tasks := []data.Task{
		{ID: "E05-tasking-permissions/T004-implement-m-hotkey", Release: "v1.1", Epic: "E05-tasking-permissions", Column: data.ColumnPlanned},
		{ID: "E05-tasking-permissions/T003-update-design-md", Release: "v1.1", Epic: "E05-tasking-permissions", Column: data.ColumnDone},
	}
	m := NewModel(tasks, "v1.1", "E05-tasking-permissions")
	m.Root = root

	got, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("p")})
	first := requireModel(t, got)
	updated := processCmd(t, first, cmd)

	if !strings.Contains(updated.StatusMessage, "Router set to v1.1 E05-tasking-permissions/T004") {
		t.Fatalf("StatusMessage = %q", updated.StatusMessage)
	}
	state := readRouterFixture(t, root)
	if state.State != "task-building" {
		t.Errorf("router state = %q, want task-building", state.State)
	}
	if state.Task != "E05-tasking-permissions/T004-implement-m-hotkey" {
		t.Errorf("router task = %q, want focused task", state.Task)
	}
}

func TestUpdate_pDoesNothingWhenOverlayOpen(t *testing.T) {
	root := writeRouterFixture(t)
	tasks := []data.Task{{ID: "E05-tasking-permissions/T004-implement-m-hotkey", Release: "v1.1", Epic: "E05-tasking-permissions", Column: data.ColumnPlanned}}
	m := NewModel(tasks, "v1.1", "E05-tasking-permissions")
	m.Root = root
	m.Overlay = OverlayHelp

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("p")})
	updated := requireModel(t, got)

	if updated.StatusMessage != "" {
		t.Fatalf("StatusMessage = %q, want empty", updated.StatusMessage)
	}
	state := readRouterFixture(t, root)
	if state.Task != "T001" {
		t.Errorf("router task = %q, want unchanged T001", state.Task)
	}
}

func TestUpdate_mDoesNotSetRouterTask(t *testing.T) {
	root := writeRouterFixture(t)
	tasks := []data.Task{{ID: "E05-tasking-permissions/T004-implement-m-hotkey", Release: "v1.1", Epic: "E05-tasking-permissions", Column: data.ColumnPlanned}}
	m := NewModel(tasks, "v1.1", "E05-tasking-permissions")
	m.Root = root

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

func TestUpdate_pDoesNotSetDoneTask(t *testing.T) {
	root := writeRouterFixture(t)
	tasks := []data.Task{{ID: "E05-tasking-permissions/T004-implement-m-hotkey", Release: "v1.1", Epic: "E05-tasking-permissions", Column: data.ColumnDone}}
	m := NewModel(tasks, "v1.1", "E05-tasking-permissions")
	m.Root = root
	m.FocusedColumn = data.ColumnDone

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("p")})
	updated := requireModel(t, got)

	if updated.StatusMessage != "Router not updated: focused task is done" {
		t.Fatalf("StatusMessage = %q", updated.StatusMessage)
	}
	state := readRouterFixture(t, root)
	if state.Task != "T001" {
		t.Errorf("router task = %q, want unchanged T001", state.Task)
	}
}

func TestUpdate_spaceShowsPhaseTransitionMessage(t *testing.T) {
	tasks := []data.Task{{ID: "E05/T004", Column: data.ColumnInProgress, Stage: data.StageBuild}}
	m := NewModel(tasks, "v1.1", "E05")
	m.FocusedColumn = data.ColumnInProgress

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeySpace})
	updated := requireModel(t, got)

	if updated.StatusMessage != "Moved T004 to test" {
		t.Fatalf("StatusMessage = %q", updated.StatusMessage)
	}
	if updated.AllTasks[0].Stage != data.StageTest {
		t.Errorf("Stage = %q, want test", updated.AllTasks[0].Stage)
	}
}

func TestUpdate_spaceWarnsAfterStaleMtime(t *testing.T) {
	path := filepath.Join(t.TempDir(), "T004-task.md")
	content := "---\nid: E05/T004\nstatus: in_progress\nstage: build\nphase: build\n---\n\n# Task\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	tasks := []data.Task{{
		ID:     "E05/T004",
		Column: data.ColumnInProgress,
		Stage:  data.StageBuild,
		Path:   path,
		Mtime:  fi.ModTime().Add(-time.Hour),
	}}
	m := NewModel(tasks, "v1.1", "E05")
	m.FocusedColumn = data.ColumnInProgress

	got, cmd := m.Update(tea.KeyMsg{Type: tea.KeySpace})
	// cmd should be errorMsg since mtime is stale
	msg := cmd()
	if _, ok := msg.(errorMsg); !ok {
		t.Fatalf("expected errorMsg, got %T", msg)
	}
	updated := requireModel(t, got)
	got2, _ := updated.Update(msg)
	updated2 := requireModel(t, got2)

	if updated2.StatusMessage != "mtime conflict: refresh before retrying" {
		t.Fatalf("StatusMessage = %q", updated2.StatusMessage)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := data.NewParser().ParseTaskFile(path, string(raw))
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Stage != data.StageBuild {
		t.Errorf("persisted Stage = %q, want build", parsed.Stage)
	}
	if !strings.Contains(string(raw), "stage:") {
		t.Error("legacy stage field should remain when write is rejected")
	}
	if updated2.AllTasks[0].Stage != data.StageBuild {
		t.Errorf("model Stage = %q after rejected write, want build", updated2.AllTasks[0].Stage)
	}
}

func TestUpdate_backspaceShowsRetreatMessageAndSyncsStatus(t *testing.T) {
	tasks := []data.Task{{ID: "E05/T004", Status: string(data.StatusDone), Column: data.ColumnDone}}
	m := NewModel(tasks, "v1.1", "E05")
	m.FocusedColumn = data.ColumnDone

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	updated := requireModel(t, got)

	if updated.StatusMessage != "Moved back T004 to audit" {
		t.Fatalf("StatusMessage = %q", updated.StatusMessage)
	}
	if updated.AllTasks[0].Column != data.ColumnInProgress {
		t.Errorf("Column = %q, want in_progress", updated.AllTasks[0].Column)
	}
	if updated.AllTasks[0].Stage != data.StageAudit {
		t.Errorf("Stage = %q, want audit", updated.AllTasks[0].Stage)
	}
	if updated.AllTasks[0].Status != string(data.StatusInProgress) {
		t.Errorf("Status = %q, want in_progress", updated.AllTasks[0].Status)
	}
}

func TestUpdate_key1SwitchesToDetailTab(t *testing.T) {
	m := NewModel(nil, "v1.1", "E06-audit-command")
	m.Epics = []string{"E06-audit-command"}
	m.Overlay = OverlayEpicDetail
	m.EpicDetailTab = 1
	m.EpicDetailOffset = 5

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("1")})
	updated := requireModel(t, got)

	if updated.EpicDetailTab != 0 {
		t.Errorf("EpicDetailTab = %d, want 0", updated.EpicDetailTab)
	}
	if updated.EpicDetailOffset != 0 {
		t.Errorf("EpicDetailOffset = %d, want 0 (reset on tab switch)", updated.EpicDetailOffset)
	}
}

func TestUpdate_key2SwitchesToAuditTabAndLoadsContent(t *testing.T) {
	root := t.TempDir()
	auditDir := filepath.Join(root, "releases", "v1.1", "epics", "E06-audit-command")
	if err := os.MkdirAll(auditDir, 0755); err != nil {
		t.Fatal(err)
	}
	auditContent := "# E06 Audit\n\n## Findings\n\n- [x] All good\n"
	if err := os.WriteFile(filepath.Join(auditDir, "E06-Audit.md"), []byte(auditContent), 0644); err != nil {
		t.Fatal(err)
	}

	m := NewModel(nil, "v1.1", "E06-audit-command")
	m.Root = root
	m.Epics = []string{"E06-audit-command"}
	m.EpicPanelCursor = 0
	m.Overlay = OverlayEpicDetail
	m.EpicDetailTab = 0
	m.EpicDetailOffset = 3

	got, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("2")})
	updated := requireModel(t, got)

	if updated.EpicDetailTab != 1 {
		t.Errorf("EpicDetailTab = %d, want 1", updated.EpicDetailTab)
	}
	if updated.EpicDetailOffset != 0 {
		t.Errorf("EpicDetailOffset = %d, want 0 (reset on tab switch)", updated.EpicDetailOffset)
	}

	msg := cmd()
	got2, _ := updated.Update(msg)
	updated2 := requireModel(t, got2)
	if updated2.EpicAuditContent != auditContent {
		t.Errorf("EpicAuditContent = %q, want %q", updated2.EpicAuditContent, auditContent)
	}
}

func TestUpdate_key2LoadsAuditForOpenedEpicWhenPanelCursorStale(t *testing.T) {
	root := t.TempDir()
	epicA := filepath.Join(root, "releases", "v1.1", "epics", "E02-cross-platform-compatibility")
	epicB := filepath.Join(root, "releases", "v1.1", "epics", "E06-audit-command")
	if err := os.MkdirAll(epicA, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(epicB, 0755); err != nil {
		t.Fatal(err)
	}
	auditContent := "# E06 Audit\n\n## Main Findings\nE06 content\n"
	if err := os.WriteFile(filepath.Join(epicB, "E06-Audit.md"), []byte(auditContent), 0644); err != nil {
		t.Fatal(err)
	}

	m := NewModel(nil, "v1.1", "E06-audit-command")
	m.Root = root
	m.Epics = []string{"E02-cross-platform-compatibility", "E06-audit-command"}
	m.SelectedEpic = "E06-audit-command"
	m.EpicDetailEpic = "E06-audit-command"
	m.EpicPanelCursor = 0
	m.Overlay = OverlayEpicDetail

	got, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("2")})
	updated := requireModel(t, got)

	msg := cmd()
	got2, _ := updated.Update(msg)
	updated2 := requireModel(t, got2)
	if updated2.EpicAuditContent != auditContent {
		t.Errorf("EpicAuditContent = %q, want opened epic audit content", updated2.EpicAuditContent)
	}
}

func TestUpdate_key2FallsBackWhenNoAuditFile(t *testing.T) {
	m := NewModel(nil, "v1.1", "E06-audit-command")
	m.Epics = []string{"E06-audit-command"}
	m.EpicPanelCursor = 0
	m.Overlay = OverlayEpicDetail

	got, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("2")})
	updated := requireModel(t, got)

	if updated.EpicDetailTab != 1 {
		t.Errorf("EpicDetailTab = %d, want 1", updated.EpicDetailTab)
	}

	msg := cmd()
	got2, _ := updated.Update(msg)
	updated2 := requireModel(t, got2)
	if updated2.EpicAuditContent != "(no audit available)" {
		t.Errorf("EpicAuditContent = %q, want \"(no audit available)\"", updated2.EpicAuditContent)
	}
}

func TestUpdate_key2CachesAuditContent(t *testing.T) {
	m := NewModel(nil, "v1.1", "E06-audit-command")
	m.Epics = []string{"E06-audit-command"}
	m.EpicPanelCursor = 0
	m.Overlay = OverlayEpicDetail
	m.EpicAuditContent = "already cached"

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("2")})
	updated := requireModel(t, got)

	if updated.EpicAuditContent != "already cached" {
		t.Errorf("EpicAuditContent = %q, want cached value preserved", updated.EpicAuditContent)
	}
}

func TestUpdate_openEpicDetailOverlayResetsTabState(t *testing.T) {
	m := NewModel(nil, "v1.1", "E06-audit-command")
	m.Epics = []string{"E06-audit-command"}
	m.EpicPanelCursor = 0
	m.EpicDetailTab = 1
	m.EpicAuditContent = "stale content"

	m.openEpicDetailOverlay()

	if m.EpicDetailTab != 0 {
		t.Errorf("EpicDetailTab = %d, want 0 after overlay open", m.EpicDetailTab)
	}
	if m.EpicAuditContent != "" {
		t.Errorf("EpicAuditContent = %q, want empty after overlay open", m.EpicAuditContent)
	}
}

func TestUpdate_tabKeysNoopOutsideEpicDetailOverlay(t *testing.T) {
	m := NewModel(nil, "v1.1", "E06-audit-command")
	m.Overlay = OverlayHelp
	m.EpicDetailTab = 0

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("2")})
	updated := requireModel(t, got)

	if updated.EpicDetailTab != 0 {
		t.Errorf("EpicDetailTab changed outside EpicDetail overlay: got %d", updated.EpicDetailTab)
	}
}

func TestReloadMsgUpdatesRouterState(t *testing.T) {
	m := NewModel(nil, "v1", "E01")
	m.RouterState = &data.RouterState{State: "task-building", Task: "E01/T001", NextAction: "Build E01/T001."}
	m.RouterTask = "E01/T001"

	newState := &data.RouterState{State: "task-building", Task: "E01/T002", NextAction: "Build E01/T002."}
	got, _ := m.Update(reloadMsg{
		tasks:        nil,
		releases:     []string{"v1"},
		releaseEpics: map[string][]string{"v1": {"E01"}},
		routerState:  newState,
	})
	updated := requireModel(t, got)

	if updated.RouterTask != "E01/T002" {
		t.Errorf("RouterTask = %q, want E01/T002", updated.RouterTask)
	}
	if updated.RouterState == nil || updated.RouterState.NextAction != "Build E01/T002." {
		t.Errorf("RouterState.NextAction = %q, want Build E01/T002.", updated.RouterState.NextAction)
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
