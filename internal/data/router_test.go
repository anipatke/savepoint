package data

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRouterReaderReadState(t *testing.T) {
	r := NewRouterReader()
	content := "## Current state\n\n```yaml\nstate: building\nrelease: v1\nepic: E01-go-setup\ntask: E01-go-setup/T002-entrypoint\nnext_action: \"Start T002-entrypoint\"\n```\n"

	state, err := r.ReadState(content)
	if err != nil {
		t.Fatalf("ReadState() error = %v", err)
	}

	if state.State != "building" {
		t.Errorf("State = %v, want building", state.State)
	}
	if state.Epic != "E01-go-setup" {
		t.Errorf("Epic = %v, want E01-go-setup", state.Epic)
	}
	if state.Task != "E01-go-setup/T002-entrypoint" {
		t.Errorf("Task = %v, want E01-go-setup/T002-entrypoint", state.Task)
	}
}

func TestRouterReaderCapitalizedHeading(t *testing.T) {
	r := NewRouterReader()
	content := "## Current State\n\n```yaml\nstate: task-building\nrelease: v1\nepic: E01\ntask: E01/T001\nnext_action: \"Do thing\"\n```\n"

	state, err := r.ReadState(content)
	if err != nil {
		t.Fatalf("ReadState() error = %v", err)
	}
	if state.State != "task-building" {
		t.Errorf("State = %v, want task-building", state.State)
	}
}

func TestRouterReaderMissing(t *testing.T) {
	r := NewRouterReader()
	content := "# No state block here"

	_, err := r.ReadState(content)
	if err == nil {
		t.Error("ReadState() expected error for missing state block")
	}
}

func TestRouterReaderDefectBuilding(t *testing.T) {
	r := NewRouterReader()
	content := "## Current state\n\n```yaml\nstate: defect-building\nrelease: v1.1\nepic: E17-defect-workflow-tui\ndefect: D001-nil-pointer-on-empty-task\nnext_action: \"Fix D001\"\n```\n"

	state, err := r.ReadState(content)
	if err != nil {
		t.Fatalf("ReadState() error = %v", err)
	}
	if state.State != "defect-building" {
		t.Errorf("State = %v, want defect-building", state.State)
	}
	if state.Defect != "D001-nil-pointer-on-empty-task" {
		t.Errorf("Defect = %v, want D001-nil-pointer-on-empty-task", state.Defect)
	}
}

func TestRouterReaderNoDefectField(t *testing.T) {
	r := NewRouterReader()
	content := "## Current state\n\n```yaml\nstate: task-building\nrelease: v1.1\nepic: E01\ntask: E01/T001\nnext_action: \"Build T001\"\n```\n"

	state, err := r.ReadState(content)
	if err != nil {
		t.Fatalf("ReadState() error = %v", err)
	}
	if state.Defect != "" {
		t.Errorf("Defect = %v, want empty for older router file", state.Defect)
	}
}

// --- Reader contract -------------------------------------------------------
//
// Upgrade never rewrites `.savepoint/router.md`, so every router Savepoint has
// ever written stays in the wild in its original shape. Compatibility is
// therefore a property of the reader, and these tests pin it: the format
// anchors ReadState requires, and the tolerance it must keep.

// readFixture loads a frozen legacy fixture. The fixtures live under
// internal/init/testdata/legacy/ so one copy serves both the reader tests here
// and the upgrade tests there; see that directory's README for the freeze rule.
func readFixture(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "init", "testdata", "legacy", name))
	if err != nil {
		t.Fatalf("read legacy fixture %s: %v", name, err)
	}
	return string(data)
}

func TestRouterReader_legacyFixtureParses(t *testing.T) {
	state, err := NewRouterReader().ReadState(readFixture(t, "router.md"))
	if err != nil {
		t.Fatalf("ReadState() error = %v, want a pre-v1.5 router to still parse", err)
	}

	want := RouterState{
		State:      "task-building",
		Release:    "v1",
		Epic:       "E03-board-tui",
		Task:       "E03-board-tui/T002-columns",
		NextAction: "Build E03-board-tui/T002-columns.",
	}
	if *state != want {
		t.Errorf("ReadState() = %+v, want %+v", *state, want)
	}
}

func TestRouterReader_ignoresUnknownKeys(t *testing.T) {
	// Decoding stays non-strict: a router written by a newer or older Savepoint
	// may carry keys this version does not model, and must still load.
	content := "## Current state\n\n```yaml\nstate: task-building\nrelease: v1\nepic: E01\ntask: E01/T001\nnext_action: \"Build T001\"\nupdated: 2024-11-02\nowner: someone\n```\n"

	state, err := NewRouterReader().ReadState(content)
	if err != nil {
		t.Fatalf("ReadState() error = %v, want unknown keys ignored", err)
	}
	if state.State != "task-building" {
		t.Errorf("State = %q, want task-building", state.State)
	}
	if state.Task != "E01/T001" {
		t.Errorf("Task = %q, want E01/T001", state.Task)
	}
}

func TestRouterReader_requiresStructuralAnchors(t *testing.T) {
	// The heading and its fenced YAML block are the parse contract. Losing
	// either must fail loudly rather than silently yielding an empty state.
	cases := map[string]string{
		"no current state heading": "# Agent State Machine\n\n```yaml\nstate: task-building\n```\n",
		"no fenced yaml block":     "## Current state\n\nstate: task-building\n",
		"unfenced yaml block":      "## Current state\n\n```\nstate: task-building\n```\n",
	}

	for name, content := range cases {
		t.Run(name, func(t *testing.T) {
			state, err := NewRouterReader().ReadState(content)
			if err == nil {
				t.Fatalf("ReadState() = %+v, want an error", state)
			}
			if !strings.Contains(err.Error(), "block") {
				t.Errorf("error = %q, want it to name the missing block", err.Error())
			}
		})
	}
}

func TestRouterReader_shippedTemplateCarriesBothAnchors(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "templates", "project", ".savepoint", "router.md"))
	if err != nil {
		t.Fatalf("read shipped router template: %v", err)
	}
	shipped := string(data)

	if !strings.Contains(shipped, stateBlockStart) {
		t.Errorf("shipped router template missing %q heading", stateBlockStart)
	}
	if !strings.Contains(shipped, "```yaml") {
		t.Error("shipped router template missing its fenced yaml block")
	}
	if _, err := NewRouterReader().ReadState(shipped); err != nil {
		t.Errorf("ReadState() on the shipped router template error = %v", err)
	}
}
