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

	model, err := newProjectModel(projectRoot)
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

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}
