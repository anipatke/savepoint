package testutil

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

// WriteRouter writes a router.md at root with the given fields.
// If nextAction is empty, it defaults to "".
func WriteRouter(t testing.TB, root, state, release, epic, task, nextAction string) {
	t.Helper()
	if nextAction == "" {
		nextAction = `""`
	} else {
		nextAction = `"` + nextAction + `"`
	}
	if task == "" {
		task = `""`
	} else {
		task = `"` + task + `"`
	}
	content := "# Agent State Machine\n\n## Current state\n\n```yaml\nstate: " + state + "\nrelease: " + release + "\nepic: " + epic + "\ntask: " + task + "\nnext_action: " + nextAction + "\n```\n"
	WriteFile(t, filepath.Join(root, "router.md"), content)
}

// WriteReleasePRD writes a minimal release PRD file.
func WriteReleasePRD(t testing.TB, releasePath string) {
	t.Helper()
	releaseID := filepath.Base(releasePath)
	WriteFile(t, filepath.Join(releasePath, releaseID+"-PRD.md"), "---\ntype: project-prd\nstatus: active\n---\n\n# Release\n")
}

// WriteEpicDetail writes a minimal epic detail file.
func WriteEpicDetail(t testing.TB, epicPath, prefix string) {
	t.Helper()
	WriteFile(t, filepath.Join(epicPath, prefix+"-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# Epic\n")
}

// TaskFixture describes a task file to create.
type TaskFixture struct {
	Slug      string // e.g. "T001-task" — becomes filename
	Release   string // optional; omitted if empty
	Status    string
	Stage     string // optional canonical stage; omitted if empty
	Phase     string // optional legacy phase; omitted if empty
	Objective string
	DependsOn []string          // optional; defaults to empty list
	Body      string            // optional; defaults to minimal body
	Extra     map[string]string // optional; extra frontmatter fields
}

// WriteTask writes a task file in the savepoint directory structure.
func WriteTask(t testing.TB, root, release, epic string, task TaskFixture) {
	t.Helper()
	path := filepath.Join(root, "releases", release, "epics", epic, "tasks", task.Slug+".md")
	id := epic + "/" + task.Slug

	var b strings.Builder
	b.WriteString("---\n")
	b.WriteString("id: " + id + "\n")
	if task.Release != "" {
		b.WriteString("release: " + task.Release + "\n")
	}
	b.WriteString("status: " + task.Status + "\n")
	if task.Stage != "" {
		b.WriteString("stage: " + task.Stage + "\n")
	} else if task.Phase != "" {
		b.WriteString("phase: " + task.Phase + "\n")
	}
	b.WriteString("objective: \"" + task.Objective + "\"\n")

	deps := "[]"
	if len(task.DependsOn) > 0 {
		quoted := make([]string, len(task.DependsOn))
		for i, d := range task.DependsOn {
			quoted[i] = fmt.Sprintf("%q", d)
		}
		deps = "[" + strings.Join(quoted, ", ") + "]"
	}
	b.WriteString("depends_on: " + deps + "\n")
	for k, v := range task.Extra {
		b.WriteString(k + ": " + v + "\n")
	}
	b.WriteString("---\n\n")

	if task.Body != "" {
		b.WriteString(task.Body)
	} else {
		b.WriteString("# " + task.Slug + "\n\n## Acceptance Criteria\n\n- it works\n")
	}

	WriteFile(t, path, b.String())
}

// SetupMinimalProject creates a minimal valid savepoint project structure.
// It creates config.yml, router.md, release PRD, epic detail, and the directory tree.
func SetupMinimalProject(t testing.TB, root, release, epic string) {
	t.Helper()
	releasePath := filepath.Join(root, "releases", release)
	//nolint:gocritic // assignment to same epicPath variable is fine
	epicPath := filepath.Join(releasePath, "epics", epic)
	tasksPath := filepath.Join(epicPath, "tasks")
	MkdirAll(t, tasksPath)

	prefix := epic
	if idx := strings.IndexByte(epic, '-'); idx != -1 {
		prefix = epic[:idx]
	}

	WriteFile(t, filepath.Join(root, "config.yml"), "quality_gates:\n  lint: null\n  typecheck: null\n  test: null\ntheme:\n  bg: \"#000\"\n")
	WriteRouter(t, root, "task-building", release, epic, "", "")
	WriteReleasePRD(t, releasePath)
	WriteEpicDetail(t, epicPath, prefix)
}
