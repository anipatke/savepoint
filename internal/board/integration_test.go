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

// writeRouterForIntegration creates a minimal router.md in the given savepoint root.
func writeRouterForIntegration(t *testing.T, savepointRoot, release, epic string) {
	t.Helper()
	content := "# Agent State Machine\n\n## Current state\n\n```yaml\nstate: task-building\nrelease: " + release + "\nepic: " + epic + "\ntask: \"\"\nnext_action: \"test\"\n```\n"
	writeFile(t, filepath.Join(savepointRoot, "router.md"), content)
}

// writeTaskWithBody creates a task file with a body section to verify content preservation.
func writeTaskWithBody(t *testing.T, root, release, epic, task string, column data.ColumnType, body string) string {
	t.Helper()
	path := filepath.Join(root, "releases", release, "epics", epic, "tasks", task+".md")
	phase := ""
	if column == data.ColumnInProgress {
		phase = "phase: build\n"
	}
	content := "---\nid: " + epic + "/" + task + "\nrelease: " + release + "\nstatus: " + string(column) + "\n" + phase + "objective: \"Test task\"\ndepends_on: []\n---\n\n" + body
	writeFile(t, path, content)
	return path
}

// TestBoardPipeline_endToEnd loads a real project from disk and renders the full board view.
func TestBoardPipeline_endToEnd(t *testing.T) {
	projectRoot := t.TempDir()
	savepointRoot := filepath.Join(projectRoot, ".savepoint")
	writeRouterForIntegration(t, savepointRoot, "v1", "E01-init")
	writeTask(t, savepointRoot, "v1", "E01-init", "T001-scaffold", data.ColumnPlanned)
	writeTask(t, savepointRoot, "v1", "E01-init", "T002-validate", data.ColumnInProgress)
	writeTask(t, savepointRoot, "v1", "E01-init", "T003-done-task", data.ColumnDone)

	model, err := newProjectModel(projectRoot, "", "")
	if err != nil {
		t.Fatalf("newProjectModel: %v", err)
	}
	if model.Watcher != nil {
		t.Cleanup(func() { model.Watcher.Close() })
	}

	model.Width = 120
	model.Height = 40
	view := model.View()

	for _, want := range []string{"PLANNED", "IN PROGRESS", "DONE", "T001", "T002", "T003"} {
		if !strings.Contains(view, want) {
			t.Errorf("board view missing %q", want)
		}
	}
}

// TestRunPlainOutput_endToEnd calls runPlainOutput against a real temp project root.
func TestRunPlainOutput_endToEnd(t *testing.T) {
	projectRoot := t.TempDir()
	savepointRoot := filepath.Join(projectRoot, ".savepoint")
	writeRouterForIntegration(t, savepointRoot, "v1", "E01-init")
	writeTask(t, savepointRoot, "v1", "E01-init", "T001-scaffold", data.ColumnPlanned)
	writeTask(t, savepointRoot, "v1", "E01-init", "T002-validate", data.ColumnDone)

	model, err := newProjectModel(projectRoot, "", "")
	if err != nil {
		t.Fatalf("newProjectModel: %v", err)
	}
	if model.Watcher != nil {
		t.Cleanup(func() { model.Watcher.Close() })
	}

	out := RenderPlainTable(model)

	if !strings.Contains(out, plainNonTTYWarning) {
		t.Error("plain output missing non-TTY warning")
	}
	for _, want := range []string{"PLANNED", "DONE", "T001-scaffold", "T002-validate"} {
		if !strings.Contains(out, want) {
			t.Errorf("plain output missing %q", want)
		}
	}
}

// TestStatusWrite_preservesTaskBody advances a task via space key and verifies the body text is unchanged.
func TestStatusWrite_preservesTaskBody(t *testing.T) {
	root := t.TempDir()
	body := "## Acceptance Criteria\n\n- [ ] thing one\n- [ ] thing two\n"
	path := writeTaskWithBody(t, root, "v1", "E01-init", "T001-scaffold", data.ColumnPlanned, body)

	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	task := data.Task{
		ID:     "E01-init/T001-scaffold",
		Column: data.ColumnPlanned,
		Path:   path,
		Mtime:  fi.ModTime(),
	}

	m := NewModel([]data.Task{task}, "v1", "E01-init")
	m.FocusedColumn = data.ColumnPlanned

	got, cmd := m.Update(tea.KeyMsg{Type: tea.KeySpace})
	msg := cmd()
	got2, _ := got.Update(msg)
	updated := requireModel(t, got2)

	if updated.AllTasks[0].Column != data.ColumnInProgress {
		t.Errorf("Column = %q, want in_progress", updated.AllTasks[0].Column)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), body) {
		t.Errorf("task body was altered after status write; got:\n%s", raw)
	}
}

// TestMtimeConflict_directDetection verifies WriteTaskStatus returns ErrMtimeConflict on mtime mismatch.
func TestMtimeConflict_directDetection(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "T001.md")
	content := "---\nid: E01/T001\nstatus: planned\nphase: build\n---\n\n# Task\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	task := &data.Task{
		ID:     "E01/T001",
		Column: data.ColumnInProgress,
		Stage:  data.StageBuild,
	}
	staleTime := time.Now().Add(-time.Hour)
	err := data.WriteTaskStatus(path, task, staleTime)
	if err != data.ErrMtimeConflict {
		t.Errorf("WriteTaskStatus with stale mtime = %v, want ErrMtimeConflict", err)
	}
}

// TestMtimeConflict_boardWarns verifies the board surfaces an mtime conflict instead of overwriting external edits.
func TestMtimeConflict_boardWarns(t *testing.T) {
	path := filepath.Join(t.TempDir(), "T001.md")
	content := "---\nid: E01/T001\nstatus: in_progress\nphase: build\n---\n\n# Task\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	task := data.Task{
		ID:     "E01/T001",
		Column: data.ColumnInProgress,
		Stage:  data.StageBuild,
		Path:   path,
		Mtime:  fi.ModTime().Add(-time.Minute), // intentionally stale
	}
	m := NewModel([]data.Task{task}, "v1", "E01")
	m.FocusedColumn = data.ColumnInProgress

	got, cmd := m.Update(tea.KeyMsg{Type: tea.KeySpace})
	msg := cmd()
	got2, _ := got.Update(msg)
	updated := requireModel(t, got2)

	if !strings.Contains(updated.StatusMessage, "mtime conflict") {
		t.Errorf("StatusMessage = %q, want mtime conflict warning", updated.StatusMessage)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), "phase: build") {
		t.Errorf("task file was overwritten despite mtime conflict:\n%s", raw)
	}
}

// TestReleaseFilter_showsOnlyMatchingRelease verifies the --release flag filters tasks.
func TestReleaseFilter_showsOnlyMatchingRelease(t *testing.T) {
	projectRoot := t.TempDir()
	savepointRoot := filepath.Join(projectRoot, ".savepoint")
	writeRouterForIntegration(t, savepointRoot, "v1", "E01-init")
	writeTask(t, savepointRoot, "v1", "E01-init", "T001-v1-task", data.ColumnPlanned)
	writeTask(t, savepointRoot, "v2", "E01-init", "T001-v2-task", data.ColumnPlanned)

	model, err := newProjectModel(projectRoot, "v2", "")
	if err != nil {
		t.Fatalf("newProjectModel: %v", err)
	}
	if model.Watcher != nil {
		t.Cleanup(func() { model.Watcher.Close() })
	}

	if model.SelectedRelease != "v2" {
		t.Errorf("SelectedRelease = %q, want v2", model.SelectedRelease)
	}
	planned := model.Tasks[data.ColumnPlanned]
	for _, task := range planned {
		if task.Release != "v2" {
			t.Errorf("task %q has release %q, want v2 only", task.ID, task.Release)
		}
	}
}

// TestEpicFilter_showsOnlyMatchingEpic verifies the --epic flag filters tasks.
func TestEpicFilter_showsOnlyMatchingEpic(t *testing.T) {
	projectRoot := t.TempDir()
	savepointRoot := filepath.Join(projectRoot, ".savepoint")
	writeRouterForIntegration(t, savepointRoot, "v1", "E01-init")
	writeTask(t, savepointRoot, "v1", "E01-init", "T001-alpha", data.ColumnPlanned)
	writeTask(t, savepointRoot, "v1", "E02-build", "T001-beta", data.ColumnPlanned)

	model, err := newProjectModel(projectRoot, "v1", "E02-build")
	if err != nil {
		t.Fatalf("newProjectModel: %v", err)
	}
	if model.Watcher != nil {
		t.Cleanup(func() { model.Watcher.Close() })
	}

	if model.SelectedEpic != "E02-build" {
		t.Errorf("SelectedEpic = %q, want E02-build", model.SelectedEpic)
	}
	planned := model.Tasks[data.ColumnPlanned]
	for _, task := range planned {
		if task.Epic != "E02-build" {
			t.Errorf("task %q has epic %q, want E02-build only", task.ID, task.Epic)
		}
	}
}

// TestDetailPane_opensOnEnter verifies enter key opens the detail overlay.
func TestDetailPane_opensOnEnter(t *testing.T) {
	tasks := []data.Task{{ID: "E01/T001", Title: "Scaffold init", Column: data.ColumnPlanned}}
	m := NewModel(tasks, "v1", "E01")
	m.FocusedColumn = data.ColumnPlanned

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	updated := requireModel(t, got)

	if updated.Overlay != OverlayDetail {
		t.Errorf("Overlay = %q, want %q", updated.Overlay, OverlayDetail)
	}
}

// TestDetailPane_escClosesOverlay verifies esc dismisses the detail overlay.
func TestDetailPane_escClosesOverlay(t *testing.T) {
	m := NewModel(nil, "v1", "E01")
	m.Overlay = OverlayDetail

	got, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	updated := requireModel(t, got)

	if updated.Overlay != OverlayNone {
		t.Errorf("Overlay = %q after esc, want none", updated.Overlay)
	}
}
