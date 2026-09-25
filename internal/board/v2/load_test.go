package v2

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/testutil"
)

func TestLoadProjectLoadsIndexRouterAndNext(t *testing.T) {
	root := writeValidProject(t)

	loaded := loadProject(root)

	if loaded.Failed() {
		t.Fatalf("loadProject() diagnostic = %q, want a loaded project", loaded.Diagnostic)
	}
	if got := loaded.State.objectiveCount(); got != 1 {
		t.Errorf("objectiveCount() = %d, want 1", got)
	}
	if got := loaded.State.taskCount(); got != 2 {
		t.Errorf("taskCount() = %d, want 2", got)
	}
	grouped := groupTaskCards(loaded.State.Index)
	if got := len(grouped[data.ColumnPlanned]); got != 1 {
		t.Errorf("planned tasks = %d, want 1", got)
	}
	if got := len(grouped[data.ColumnDone]); got != 1 {
		t.Errorf("done tasks = %d, want 1", got)
	}
	if loaded.State.Router == nil || loaded.State.Router.State != data.RouterPhaseTask {
		t.Fatalf("Router = %+v, want the decoded task phase", loaded.State.Router)
	}
	if loaded.State.Next.Kind != data.NextExecute {
		t.Errorf("Next.Kind = %q, want %q", loaded.State.Next.Kind, data.NextExecute)
	}
	if loaded.State.Next.Task == nil || loaded.State.Next.Task.ID != "T-001" {
		t.Errorf("Next.Task = %+v, want T-001", loaded.State.Next.Task)
	}
}

func TestLoadProjectEmptyProjectFromTemplateIsNormal(t *testing.T) {
	root := writeEmptyProjectFromTemplate(t)

	loaded := loadProject(root)

	if loaded.Failed() {
		t.Fatalf("loadProject() diagnostic = %q, want a fresh project to load cleanly", loaded.Diagnostic)
	}
	if loaded.State.objectiveCount() != 0 || loaded.State.taskCount() != 0 {
		t.Errorf("counts = %d objectives / %d tasks, want an empty project", loaded.State.objectiveCount(), loaded.State.taskCount())
	}
	if loaded.State.Next.Kind != data.NextNothingSelected {
		t.Errorf("Next.Kind = %q, want %q for a project with nothing selected", loaded.State.Next.Kind, data.NextNothingSelected)
	}
	if loaded.State.Router == nil || loaded.State.Router.Objective != "" {
		t.Errorf("Router objective = %+v, want no selection", loaded.State.Router)
	}
}

// TestLoadProjectRefusedIndexReportsNamedDiagnostic exercises each structural
// refusal LoadV2Index makes that a board user can cause, and proves every one
// of them arrives as a diagnostic naming the offending file and the problem —
// never as a loaded project and never as a V1 discovery attempt.
func TestLoadProjectRefusedIndexReportsNamedDiagnostic(t *testing.T) {
	for _, test := range invalidProjectCases() {
		t.Run(test.Name, func(t *testing.T) {
			root := writeValidProject(t)
			test.Build(t, root)

			loaded := loadProject(root)

			if !loaded.Failed() {
				t.Fatalf("loadProject() loaded a project the index must refuse")
			}
			if !strings.Contains(loaded.Diagnostic, test.WantPath) {
				t.Errorf("diagnostic = %q, want the offending file path %q named", loaded.Diagnostic, test.WantPath)
			}
			for _, part := range test.WantParts {
				if !strings.Contains(loaded.Diagnostic, part) {
					t.Errorf("diagnostic = %q, want it to name %q", loaded.Diagnostic, part)
				}
			}
			if strings.Contains(loaded.Diagnostic, "releases directory not found") {
				t.Errorf("diagnostic = %q, want no V1 discovery fallback", loaded.Diagnostic)
			}
			if loaded.State.Index != nil || loaded.State.Router != nil {
				t.Errorf("State = %+v, want nothing loaded alongside a diagnostic", loaded.State)
			}
		})
	}
}

// TestLoadProjectRouterDiagnostics covers DATA-03 for the router: neither an
// unreadable file nor a state value outside the closed V2 vocabulary is healed
// into a default.
func TestLoadProjectRouterDiagnostics(t *testing.T) {
	tests := []struct {
		name      string
		build     func(t *testing.T, root string)
		wantParts []string
	}{
		{
			name: "unreadable router",
			build: func(t *testing.T, root string) {
				if err := os.Remove(filepath.Join(root, "router.md")); err != nil {
					t.Fatal(err)
				}
			},
			wantParts: []string{"router.md"},
		},
		{
			name: "malformed router body",
			build: func(t *testing.T, root string) {
				testutil.WriteFile(t, filepath.Join(root, "router.md"), "# Router\n\nno state anchor here\n")
			},
			wantParts: []string{"router.md"},
		},
		{
			name: "unknown router state",
			build: func(t *testing.T, root string) {
				testutil.WriteFile(t, filepath.Join(root, "router.md"), routerContent("shipping", "O-001", "T-001"))
			},
			wantParts: []string{"router.md", "router state", "shipping"},
		},
		{
			name: "unknown router key",
			build: func(t *testing.T, root string) {
				testutil.WriteFile(t, filepath.Join(root, "router.md"),
					"# Router\n\n## Current state\n\n```yaml\nstate: task\nobjective: O-001\ntask: T-001\nepic: E01-legacy\n```\n")
			},
			wantParts: []string{"router.md"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := writeValidProject(t)
			test.build(t, root)

			loaded := loadProject(root)

			if !loaded.Failed() {
				t.Fatalf("loadProject() healed a router problem into a default")
			}
			for _, part := range test.wantParts {
				if !strings.Contains(loaded.Diagnostic, part) {
					t.Errorf("diagnostic = %q, want it to name %q", loaded.Diagnostic, part)
				}
			}
		})
	}
}

// TestLoadCmdReturnsTheSameSingleMessage proves startup and every later reload
// share one load path: the command is a wrapper over loadProject and returns
// its one message type.
func TestLoadCmdReturnsTheSameSingleMessage(t *testing.T) {
	root := writeValidProject(t)

	msg := loadCmd(root)()

	loaded, ok := msg.(projectLoadedMsg)
	if !ok {
		t.Fatalf("loadCmd() returned %T, want projectLoadedMsg", msg)
	}
	if loaded.State.taskCount() != loadProject(root).State.taskCount() {
		t.Error("loadCmd() and loadProject() disagree about the same project")
	}
}

// TestHeaderCountsFinishedOutOfTotal proves each header count's first number
// is the finished records only: done Objectives and Tasks, resolved Issues.
func TestHeaderCountsFinishedOutOfTotal(t *testing.T) {
	state := ProjectState{Index: &data.V2Index{
		Objectives: map[string]*data.ObjectiveV2{
			"O-001": {Status: data.ColumnDone},
			"O-002": {Status: data.ColumnInProgress},
		},
		Tasks: map[string]*data.TaskV2{
			"T-001": {Status: data.ColumnDone},
			"T-002": {Status: data.ColumnDone},
			"T-003": {Status: data.ColumnPlanned},
		},
		Issues: map[string]*data.IssueV2{
			"I-001": {Status: data.IssueStatusResolved},
			"I-002": {Status: data.IssueStatusInProgress},
			"I-003": {Status: data.IssueStatusOpen},
			"I-004": {Status: data.IssueStatusOpen},
		},
	}}

	if finished, total := state.objectiveCounts(); finished != 1 || total != 2 {
		t.Errorf("objectiveCounts() = %d/%d, want 1/2", finished, total)
	}
	if finished, total := state.taskCounts(); finished != 2 || total != 3 {
		t.Errorf("taskCounts() = %d/%d, want 2/3", finished, total)
	}
	if finished, total := state.issueCounts(); finished != 1 || total != 4 {
		t.Errorf("issueCounts() = %d/%d, want 1/4", finished, total)
	}
}
