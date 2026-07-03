package board

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/testutil"
)

// writeTaskWithBody creates a task file with a body section to verify content preservation.
func writeTaskWithBody(t *testing.T, root, release, epic, taskSlug string, column data.ColumnType, body string) string {
	t.Helper()
	tf := testutil.TaskFixture{
		Slug:      taskSlug,
		Release:   release,
		Status:    string(column),
		Objective: "Test task",
		Body:      body,
	}
	if column == data.ColumnInProgress {
		tf.Stage = "build"
	}
	testutil.WriteTask(t, root, release, epic, tf)
	return filepath.Join(root, "releases", release, "epics", epic, "tasks", taskSlug+".md")
}

// TestBoardPipeline_endToEnd loads a real project from disk and renders the full board view.
func TestBoardPipeline_endToEnd(t *testing.T) {
	projectRoot := t.TempDir()
	savepointRoot := filepath.Join(projectRoot, ".savepoint")
	testutil.WriteRouter(t, savepointRoot, "task-building", "v1", "E01-init", "", "test")
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
	testutil.WriteRouter(t, savepointRoot, "task-building", "v1", "E01-init", "", "test")
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
	content := "---\nid: E01/T001\nstatus: planned\nstage: build\n---\n\n# Task\n"
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

// TestMtimeConflict_boardRefreshesChangedTask verifies the board refreshes instead of overwriting external edits.
func TestMtimeConflict_boardRefreshesChangedTask(t *testing.T) {
	path := filepath.Join(t.TempDir(), "T001.md")
	content := "---\nid: E01/T001\nstatus: in_progress\nstage: test\n---\n\n# Task\n"
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

	if !strings.Contains(updated.StatusMessage, "task changed on disk") {
		t.Errorf("StatusMessage = %q, want changed-on-disk warning", updated.StatusMessage)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), "stage: test") {
		t.Errorf("task file was overwritten despite changed-on-disk refresh:\n%s", raw)
	}
}

// TestReleaseFilter_showsOnlyMatchingRelease verifies the --release flag filters tasks.
func TestReleaseFilter_showsOnlyMatchingRelease(t *testing.T) {
	projectRoot := t.TempDir()
	savepointRoot := filepath.Join(projectRoot, ".savepoint")
	testutil.WriteRouter(t, savepointRoot, "task-building", "v1", "E01-init", "", "test")
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

func TestReleaseFilter_acceptsNamedRelease(t *testing.T) {
	projectRoot := t.TempDir()
	savepointRoot := filepath.Join(projectRoot, ".savepoint")
	release := "quiz-tales-journals"
	testutil.WriteRouter(t, savepointRoot, "task-building", "v1", "E01-init", "", "test")
	writeTask(t, savepointRoot, "v1", "E01-init", "T001-v1-task", data.ColumnPlanned)
	writeTask(t, savepointRoot, release, "E01-init", "T001-named-release-task", data.ColumnPlanned)

	model, err := newProjectModel(projectRoot, release, "")
	if err != nil {
		t.Fatalf("newProjectModel: %v", err)
	}
	if model.Watcher != nil {
		t.Cleanup(func() { model.Watcher.Close() })
	}

	if model.SelectedRelease != release {
		t.Errorf("SelectedRelease = %q, want %q", model.SelectedRelease, release)
	}
	planned := model.Tasks[data.ColumnPlanned]
	if len(planned) != 1 {
		t.Fatalf("planned tasks = %v, want one named release task", planned)
	}
	if planned[0].Release != release {
		t.Errorf("task release = %q, want %q", planned[0].Release, release)
	}
}

func TestReleaseFilter_allowsTasklessEpicInNamedRelease(t *testing.T) {
	projectRoot := t.TempDir()
	savepointRoot := filepath.Join(projectRoot, ".savepoint")
	release := "quiz-tales-journals"
	testutil.WriteRouter(t, savepointRoot, "task-building", release, "E01-init", "", "test")
	writeTask(t, savepointRoot, release, "E01-init", "T001-named-release-task", data.ColumnPlanned)
	testutil.MkdirAll(t, filepath.Join(savepointRoot, "releases", release, "epics", "E04-THE_AUDIT"))

	model, err := newProjectModel(projectRoot, release, "")
	if err != nil {
		t.Fatalf("newProjectModel: %v", err)
	}
	if model.Watcher != nil {
		t.Cleanup(func() { model.Watcher.Close() })
	}

	if model.SelectedRelease != release {
		t.Errorf("SelectedRelease = %q, want %q", model.SelectedRelease, release)
	}
	if len(model.ReleaseEpics[release]) != 2 {
		t.Fatalf("release epics = %v, want task epic and taskless audit epic", model.ReleaseEpics[release])
	}
	planned := model.Tasks[data.ColumnPlanned]
	if len(planned) != 1 || planned[0].ID != "E01-init/T001-named-release-task" {
		t.Fatalf("planned tasks = %v, want only named release task", planned)
	}
}

func TestReleaseFilter_unknownReleaseReturnsError(t *testing.T) {
	projectRoot := t.TempDir()
	savepointRoot := filepath.Join(projectRoot, ".savepoint")
	testutil.WriteRouter(t, savepointRoot, "task-building", "v1", "E01-init", "", "test")
	writeTask(t, savepointRoot, "v1", "E01-init", "T001-v1-task", data.ColumnPlanned)

	_, err := newProjectModel(projectRoot, "quiz-tales-journals", "")
	if err == nil {
		t.Fatal("newProjectModel error = nil, want unknown release error")
	}
	if !strings.Contains(err.Error(), `release "quiz-tales-journals" not found`) {
		t.Fatalf("newProjectModel error = %v, want unknown release message", err)
	}
}

// TestEpicFilter_showsOnlyMatchingEpic verifies the --epic flag filters tasks.
func TestEpicFilter_showsOnlyMatchingEpic(t *testing.T) {
	projectRoot := t.TempDir()
	savepointRoot := filepath.Join(projectRoot, ".savepoint")
	testutil.WriteRouter(t, savepointRoot, "task-building", "v1", "E01-init", "", "test")
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

// TestAuditOverlay_openCloseLeavesBoardBehaviorUnchanged loads a real project
// with an on-disk audit register, opens and closes the Audit Register overlay,
// and proves the task-detail, defect, and release-docs overlays still behave as
// before (v1.4 release regression).
func TestAuditOverlay_openCloseLeavesBoardBehaviorUnchanged(t *testing.T) {
	projectRoot := t.TempDir()
	savepointRoot := filepath.Join(projectRoot, ".savepoint")
	testutil.WriteRouter(t, savepointRoot, "task-building", "v1", "E01-init", "", "test")
	writeTask(t, savepointRoot, "v1", "E01-init", "T001-scaffold", data.ColumnPlanned)
	testutil.WriteFile(t, filepath.Join(savepointRoot, "audit", "prompt.md"), "# Audit Prompt\nreview the release")
	testutil.WriteFile(t, filepath.Join(savepointRoot, "audit", "findings", "F001-sample.md"), `---
id: F001
title: "Sample finding"
status: open
severity: medium
confidence: high
proof_needed: "regression test"
first_seen: "2026-07-01"
last_seen: "2026-07-01"
---

# Finding
`)

	model, err := newProjectModel(projectRoot, "", "")
	if err != nil {
		t.Fatalf("newProjectModel: %v", err)
	}
	if model.Watcher != nil {
		t.Cleanup(func() { model.Watcher.Close() })
	}
	model.Width = 120
	model.Height = 40

	// Open the audit overlay and apply its async load command.
	got, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("A")})
	m := requireModel(t, got)
	if m.Overlay != OverlayAudit {
		t.Fatalf("Overlay = %q after A, want %q", m.Overlay, OverlayAudit)
	}
	if cmd == nil {
		t.Fatal("expected an audit-register load command")
	}
	got, _ = m.Update(cmd())
	m = requireModel(t, got)
	if view := m.View(); !strings.Contains(view, "AUDIT REGISTER") {
		t.Error("board view missing AUDIT REGISTER overlay after open")
	}

	// Close it.
	got, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = requireModel(t, got)
	if m.Overlay != OverlayNone {
		t.Fatalf("Overlay = %q after esc, want none", m.Overlay)
	}

	// Task detail still opens and closes.
	m.FocusedColumn = data.ColumnPlanned
	got, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = requireModel(t, got)
	if m.Overlay != OverlayDetail {
		t.Errorf("Overlay = %q after enter, want %q", m.Overlay, OverlayDetail)
	}
	got, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = requireModel(t, got)

	// Defect overlay still opens and closes.
	got, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("d")})
	m = requireModel(t, got)
	if m.Overlay != OverlayDefect {
		t.Errorf("Overlay = %q after d, want %q", m.Overlay, OverlayDefect)
	}
	got, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = requireModel(t, got)

	// Release-docs overlay still opens.
	got, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("D")})
	m = requireModel(t, got)
	if m.Overlay != OverlayReleaseDocs {
		t.Errorf("Overlay = %q after D, want %q", m.Overlay, OverlayReleaseDocs)
	}
}
