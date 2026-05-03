package board

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/data"
)

func TestNewProgramModelUsesBoardCore(t *testing.T) {
	m := newProgramModel()
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
	writeFile(t, filepath.Join(savepointRoot, "router.md"), `# Agent State Machine

## Current state

`+"```"+`yaml
state: task-building
release: v2
epic: E03-live
task: ""
next_action: "test"
`+"```"+`
`)
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
	writeFile(t, filepath.Join(savepointRoot, "router.md"), `# Agent State Machine

## Current state

`+"```"+`yaml
state: task-building
release: v1.1
epic: E01-tui-optimisation
task: E01-tui-optimisation/T001-border-resize-fix
next_action: "test"
`+"```"+`
`)
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
	writeFile(t, filepath.Join(savepointRoot, "router.md"), `# Agent State Machine

## Current state

`+"```"+`yaml
state: task-building
release: v1.1
epic: E03
task: T001
next_action: "Build v1.1 E03/T001"
`+"```"+`
`)
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

func writeTask(t *testing.T, root, release, epic, task string, column data.ColumnType) {
	t.Helper()
	path := filepath.Join(root, "releases", release, "epics", epic, "tasks", task+".md")
	phase := ""
	if column == data.ColumnInProgress {
		phase = "phase: build\n"
	}
	content := `---
id: ` + epic + `/` + task + `
release: ` + release + `
status: ` + string(column) + `
` + phase + `objective: "Test task"
depends_on: []
---

# Test task
`
	writeFile(t, path, content)
}

func writeTaskWithoutRelease(t *testing.T, root, release, epic, task string, column data.ColumnType) {
	t.Helper()
	path := filepath.Join(root, "releases", release, "epics", epic, "tasks", task+".md")
	phase := ""
	if column == data.ColumnInProgress {
		phase = "phase: build\n"
	}
	content := `---
id: ` + epic + `/` + task + `
status: ` + string(column) + `
` + phase + `objective: "Test task"
depends_on: []
---

# Test task
`
	writeFile(t, path, content)
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}
