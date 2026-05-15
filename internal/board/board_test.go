package board

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/testutil"
)

func TestNewModelUsesBoardCore(t *testing.T) {
	m := NewModel(nil, "v1", "E03-board-tui-core")
	m.Width = 100

	got := m.View()
	if strings.Contains(got, "Welcome to Savepoint") {
		t.Fatal("program model still renders placeholder welcome screen")
	}
	for _, title := range []string{"PLANNED", "IN PROGRESS", "DONE"} {
		if !strings.Contains(got, title) {
			t.Fatalf("program model view missing board column %q", title)
		}
	}
}

func TestNewProjectModelLoadsReleasesEpicsAndTasks(t *testing.T) {
	projectRoot := t.TempDir()
	savepointRoot := filepath.Join(projectRoot, ".savepoint")
	testutil.WriteRouter(t, savepointRoot, "task-building", "v2", "E03-live", "", "test")
	writeTask(t, savepointRoot, "v1", "E01-old", "T001-old", data.ColumnPlanned)
	writeTask(t, savepointRoot, "v2", "E03-live", "T001-live", data.ColumnInProgress)

	model, err := newProjectModel(projectRoot, "", "")
	if err != nil {
		t.Fatalf("newProjectModel() error = %v", err)
	}

	if model.Root != savepointRoot {
		t.Errorf("Root = %q, want %q", model.Root, savepointRoot)
	}
	if model.SelectedRelease != "v2" {
		t.Errorf("SelectedRelease = %q, want v2", model.SelectedRelease)
	}
	if model.SelectedEpic != "E03-live" {
		t.Errorf("SelectedEpic = %q, want E03-live", model.SelectedEpic)
	}
	if len(model.Releases) != 2 {
		t.Errorf("Releases = %v, want two releases", model.Releases)
	}
	if len(model.Epics) != 1 || model.Epics[0] != "E03-live" {
		t.Errorf("Epics = %v, want [E03-live]", model.Epics)
	}
	tasks := model.Tasks[data.ColumnInProgress]
	if len(tasks) != 1 || tasks[0].ID != "E03-live/T001-live" {
		t.Errorf("visible in-progress tasks = %v, want E03-live/T001-live", tasks)
	}
	if model.Watcher == nil {
		t.Fatal("Watcher is nil, want auto-refresh watcher")
	}
	t.Cleanup(func() { model.Watcher.Close() })
}

func TestNewProjectModelUsesPathReleaseForTaskWithoutReleaseFrontmatter(t *testing.T) {
	projectRoot := t.TempDir()
	savepointRoot := filepath.Join(projectRoot, ".savepoint")
	testutil.WriteRouter(t, savepointRoot, "task-building", "v1.1", "E01-tui-optimisation", "E01-tui-optimisation/T001-border-resize-fix", "test")
	writeTaskWithoutRelease(t, savepointRoot, "v1.1", "E01-tui-optimisation", "T001-border-resize-fix", data.ColumnInProgress)

	model, err := newProjectModel(projectRoot, "", "")
	if err != nil {
		t.Fatalf("newProjectModel() error = %v", err)
	}

	if model.SelectedRelease != "v1.1" {
		t.Errorf("SelectedRelease = %q, want v1.1", model.SelectedRelease)
	}
	tasks := model.Tasks[data.ColumnInProgress]
	if len(tasks) != 1 {
		t.Fatalf("visible in-progress tasks = %v, want one v1.1 task", tasks)
	}
	if tasks[0].Release != "v1.1" {
		t.Errorf("Task.Release = %q, want v1.1", tasks[0].Release)
	}
	if model.Watcher != nil {
		t.Cleanup(func() { model.Watcher.Close() })
	}
}

func TestNewProjectModelResolvesShortRouterEpicToFullEpicID(t *testing.T) {
	projectRoot := t.TempDir()
	savepointRoot := filepath.Join(projectRoot, ".savepoint")
	testutil.WriteRouter(t, savepointRoot, "task-building", "v1.1", "E03", "T001", "Build v1.1 E03/T001")
	writeTask(t, savepointRoot, "v1.1", "E01-tui-optimisation", "T007-column-focus-border-stability", data.ColumnInProgress)
	writeTask(t, savepointRoot, "v1.1", "E03-ui-visual-refinement", "T001-border-resize-fix", data.ColumnInProgress)

	model, err := newProjectModel(projectRoot, "", "")
	if err != nil {
		t.Fatalf("newProjectModel() error = %v", err)
	}

	if model.SelectedEpic != "E03-ui-visual-refinement" {
		t.Errorf("SelectedEpic = %q, want E03-ui-visual-refinement", model.SelectedEpic)
	}
	tasks := model.Tasks[data.ColumnInProgress]
	if len(tasks) != 1 || tasks[0].ID != "E03-ui-visual-refinement/T001-border-resize-fix" {
		t.Errorf("visible in-progress tasks = %v, want E03-ui-visual-refinement/T001-border-resize-fix", tasks)
	}
	if model.Watcher != nil {
		t.Cleanup(func() { model.Watcher.Close() })
	}
}

func TestUpdateReloadMsgRefreshesReleaseEpicIndex(t *testing.T) {
	m := NewModel(nil, "v1", "E01-old")
	m.Releases = []string{"v1"}
	m.ReleaseEpics = map[string][]string{"v1": []string{"E01-old"}}

	task := data.Task{
		ID:      "E02-new/T001-new",
		Release: "v1",
		Epic:    "E02-new",
		Column:  data.ColumnPlanned,
	}
	got, _ := m.Update(reloadMsg{
		tasks:        []data.Task{task},
		releases:     []string{"v1"},
		releaseEpics: map[string][]string{"v1": []string{"E02-new"}},
	})
	updated := requireModel(t, got)

	if updated.SelectedEpic != "E02-new" {
		t.Errorf("SelectedEpic = %q, want E02-new", updated.SelectedEpic)
	}
	if len(updated.Epics) != 1 || updated.Epics[0] != "E02-new" {
		t.Errorf("Epics = %v, want [E02-new]", updated.Epics)
	}
	if len(updated.Tasks[data.ColumnPlanned]) != 1 {
		t.Errorf("planned tasks = %v, want reloaded task visible", updated.Tasks[data.ColumnPlanned])
	}
}

func TestNewProjectModelFreshInit(t *testing.T) {
	projectRoot := t.TempDir()
	savepointRoot := filepath.Join(projectRoot, ".savepoint")
	testutil.WriteRouter(t, savepointRoot, "pre-implementation", "v1", "none", "", "")
	testutil.MkdirAll(t, filepath.Join(savepointRoot, "releases", "v1", "epics"))

	model, err := newProjectModel(projectRoot, "", "")
	if err != nil {
		t.Fatalf("newProjectModel() error = %v (board must not fail after fresh init)", err)
	}
	if model.Root != savepointRoot {
		t.Errorf("Root = %q, want %q", model.Root, savepointRoot)
	}
	if len(model.Releases) != 1 || model.Releases[0] != "v1" {
		t.Errorf("Releases = %v, want [v1]", model.Releases)
	}
	if len(model.Epics) != 0 {
		t.Errorf("Epics = %v, want empty for zero-epic project", model.Epics)
	}
	if model.Watcher != nil {
		t.Cleanup(func() { model.Watcher.Close() })
	}
}

func writeTask(t *testing.T, root, release, epic, taskSlug string, column data.ColumnType) {
	t.Helper()
	tf := testutil.TaskFixture{
		Slug:      taskSlug,
		Release:   release,
		Status:    string(column),
		Objective: "Test task",
	}
	if column == data.ColumnInProgress {
		tf.Stage = "build"
	}
	testutil.WriteTask(t, root, release, epic, tf)
}

func TestNewProjectModelLoadsDefectsForRelease(t *testing.T) {
	projectRoot := t.TempDir()
	savepointRoot := filepath.Join(projectRoot, ".savepoint")
	testutil.WriteRouter(t, savepointRoot, "task-building", "v2", "E01-alpha", "", "")
	writeTask(t, savepointRoot, "v2", "E01-alpha", "T001-task", data.ColumnPlanned)
	writeDefect(t, savepointRoot, "v2", "D001-crash", data.DefectOpen)
	writeDefect(t, savepointRoot, "v2", "D002-done-bug", data.DefectResolved)

	model, err := newProjectModel(projectRoot, "", "")
	if err != nil {
		t.Fatalf("newProjectModel() error = %v", err)
	}
	if len(model.AllDefects) != 2 {
		t.Fatalf("AllDefects len = %d, want 2", len(model.AllDefects))
	}
	if model.Watcher != nil {
		t.Cleanup(func() { model.Watcher.Close() })
	}
}

func TestNewProjectModelZeroDefectsWhenNoDefectDir(t *testing.T) {
	projectRoot := t.TempDir()
	savepointRoot := filepath.Join(projectRoot, ".savepoint")
	testutil.WriteRouter(t, savepointRoot, "pre-implementation", "v1", "none", "", "")
	testutil.MkdirAll(t, filepath.Join(savepointRoot, "releases", "v1", "epics"))

	model, err := newProjectModel(projectRoot, "", "")
	if err != nil {
		t.Fatalf("newProjectModel() error = %v", err)
	}
	if len(model.AllDefects) != 0 {
		t.Errorf("AllDefects = %v, want empty", model.AllDefects)
	}
	if model.Watcher != nil {
		t.Cleanup(func() { model.Watcher.Close() })
	}
}

func writeDefect(t *testing.T, root, release, slug string, status data.DefectStatus) {
	t.Helper()
	content := "---\nid: " + slug + "\nrelease: " + release + "\nstatus: " + string(status) + "\nseverity: medium\ntitle: " + slug + "\n---\n\n# " + slug + "\n"
	path := filepath.Join(root, "releases", release, "defects", slug+".md")
	testutil.WriteFile(t, path, content)
}

func writeTaskWithoutRelease(t *testing.T, root, release, epic, taskSlug string, column data.ColumnType) {
	t.Helper()
	tf := testutil.TaskFixture{
		Slug:      taskSlug,
		Status:    string(column),
		Objective: "Test task",
	}
	if column == data.ColumnInProgress {
		tf.Stage = "build"
	}
	testutil.WriteTask(t, root, release, epic, tf)
}
