package data

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadStateV2_decodesSelectedTask(t *testing.T) {
	content := "## Current state\n\n```yaml\nstate: task\nobjective: O-001\ntask: T-001\nnext_action: \"Build T-001\"\n```\n"

	state, err := NewRouterReader().ReadStateV2(content)
	if err != nil {
		t.Fatalf("ReadStateV2() error = %v", err)
	}

	want := RouterStateV2{State: RouterPhaseTask, Objective: "O-001", Task: "T-001", HasRetiredNextAction: true}
	if *state != want {
		t.Errorf("ReadStateV2() = %+v, want %+v", *state, want)
	}
}

func TestReadStateV2_decodesSelectedRelease(t *testing.T) {
	content := "## Current state\n\n```yaml\nstate: task\nrelease: R-001\nobjective: O-001\ntask: T-001\nnext_action: \"Build T-001\"\n```\n"

	state, err := NewRouterReader().ReadStateV2(content)
	if err != nil {
		t.Fatalf("ReadStateV2() error = %v", err)
	}
	if state.Release != "R-001" || state.Objective != "O-001" || state.Task != "T-001" {
		t.Fatalf("ReadStateV2() selections = release %q objective %q task %q, want R-001/O-001/T-001", state.Release, state.Objective, state.Task)
	}
}

func TestReadStateV2_decodesIssueAloneWithoutObjective(t *testing.T) {
	content := "## Current state\n\n```yaml\nstate: task\nissue: I-042\n```\n"

	state, err := NewRouterReader().ReadStateV2(content)
	if err != nil {
		t.Fatalf("ReadStateV2() error = %v", err)
	}
	if state.Issue != "I-042" || state.Objective != "" || state.Task != "" {
		t.Fatalf("ReadStateV2() selections = issue %q objective %q task %q, want I-042 with no Objective or Task", state.Issue, state.Objective, state.Task)
	}
	if state.HasRetiredNextAction {
		t.Error("HasRetiredNextAction = true without a next_action key")
	}
}

func TestReadStateV2_noneSelectionsDecodeAsEmpty(t *testing.T) {
	content := "## Current state\n\n```yaml\nstate: idea\nrelease: none\nobjective: none\ntask: none\nissue: none\nnext_action: \"\"\n```\n"

	state, err := NewRouterReader().ReadStateV2(content)
	if err != nil {
		t.Fatalf("ReadStateV2() error = %v", err)
	}
	if state.Objective != "" {
		t.Errorf("Objective = %q, want empty for the none sentinel", state.Objective)
	}
	if state.Release != "" {
		t.Errorf("Release = %q, want empty for the none sentinel", state.Release)
	}
	if state.Task != "" {
		t.Errorf("Task = %q, want empty for the none sentinel", state.Task)
	}
	if state.Issue != "" {
		t.Errorf("Issue = %q, want empty for the none sentinel", state.Issue)
	}
	if state.State != RouterPhaseIdea {
		t.Errorf("State = %q, want idea", state.State)
	}
	if !state.HasRetiredNextAction {
		t.Error("HasRetiredNextAction = false, want true for a present next_action key")
	}
}

func TestReadStateV2_missingBlankAndNoneGoalSelectNothing(t *testing.T) {
	cases := []struct {
		name string
		goal string
	}{
		{name: "missing"},
		{name: "blank", goal: "release: \"\"\n"},
		{name: "none", goal: "release: none\n"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			content := "## Current state\n\n```yaml\nstate: task\n" + tc.goal + "objective: O-001\ntask: none\nissue: I-001\n```\n"
			router, err := NewRouterReader().ReadStateV2(content)
			if err != nil {
				t.Fatalf("ReadStateV2() error = %v", err)
			}
			index := &V2Index{
				Releases: map[string]*ReleaseV2{"R-001": {ID: "R-001"}},
				Objectives: map[string]*ObjectiveV2{
					"O-001": {ID: "O-001", Release: "R-001"},
				},
				Issues: map[string]*IssueV2{
					"I-001": {ID: "I-001", Status: IssueStatusOpen},
				},
			}
			selection, diagnostic := ResolveSelection(index, router)
			if diagnostic == nil || diagnostic.Kind != SelectionReleaseMissing {
				t.Fatalf("ResolveSelection() diagnostic = %+v, want missing Goal", diagnostic)
			}
			if selection != (Selection{}) {
				t.Errorf("ResolveSelection() = %+v, want no selection", selection)
			}
		})
	}
}

func TestReadStateV2_taskAbsentObjectivePresentDecodesCleanly(t *testing.T) {
	content := "## Current state\n\n```yaml\nstate: design\nobjective: O-002\ntask: none\nnext_action: \"Plan O-002\"\n```\n"

	state, err := NewRouterReader().ReadStateV2(content)
	if err != nil {
		t.Fatalf("ReadStateV2() error = %v", err)
	}
	if state.Objective != "O-002" {
		t.Errorf("Objective = %q, want O-002", state.Objective)
	}
	if state.Task != "" {
		t.Errorf("Task = %q, want empty", state.Task)
	}
}

func TestReadStateV2_recordsEmptyRetiredNextActionPresence(t *testing.T) {
	content := "## Current state\n\n```yaml\nstate: check\nobjective: O-003\ntask: none\nnext_action:\n```\n"

	state, err := NewRouterReader().ReadStateV2(content)
	if err != nil {
		t.Fatalf("ReadStateV2() error = %v", err)
	}
	if !state.HasRetiredNextAction {
		t.Error("HasRetiredNextAction = false, want true for an empty next_action key")
	}
}

func TestReadStateV2_unknownStateIsNamedDiagnostic(t *testing.T) {
	content := "## Current state\n\n```yaml\nstate: building\nobjective: none\ntask: none\nnext_action: \"\"\n```\n"

	_, err := NewRouterReader().ReadStateV2(content)
	if !errors.Is(err, ErrV2InvalidLifecycle) {
		t.Fatalf("ReadStateV2() error = %v, want ErrV2InvalidLifecycle", err)
	}
	if !strings.Contains(err.Error(), "building") {
		t.Errorf("error = %q, want it to name the unrecognized state", err.Error())
	}
}

func TestReadStateV2_emptyStateIsNamedDiagnostic(t *testing.T) {
	content := "## Current state\n\n```yaml\nstate: \"\"\nobjective: none\ntask: none\nnext_action: \"\"\n```\n"

	_, err := NewRouterReader().ReadStateV2(content)
	if !errors.Is(err, ErrV2InvalidLifecycle) {
		t.Fatalf("ReadStateV2() error = %v, want ErrV2InvalidLifecycle", err)
	}
}

func TestReadStateV2_malformedObjectiveID(t *testing.T) {
	content := "## Current state\n\n```yaml\nstate: design\nobjective: abc\ntask: none\nnext_action: \"\"\n```\n"

	_, err := NewRouterReader().ReadStateV2(content)
	if !errors.Is(err, ErrV2InvalidID) {
		t.Fatalf("ReadStateV2() error = %v, want ErrV2InvalidID", err)
	}
}

func TestReadStateV2_malformedReleaseID(t *testing.T) {
	content := "## Current state\n\n```yaml\nstate: design\nrelease: release-1\nobjective: none\ntask: none\nnext_action: \"\"\n```\n"

	_, err := NewRouterReader().ReadStateV2(content)
	if !errors.Is(err, ErrV2InvalidID) {
		t.Fatalf("ReadStateV2() error = %v, want ErrV2InvalidID", err)
	}
}

func TestReadStateV2_malformedTaskID(t *testing.T) {
	content := "## Current state\n\n```yaml\nstate: task\nobjective: O-001\ntask: xyz\nnext_action: \"\"\n```\n"

	_, err := NewRouterReader().ReadStateV2(content)
	if !errors.Is(err, ErrV2InvalidID) {
		t.Fatalf("ReadStateV2() error = %v, want ErrV2InvalidID", err)
	}
}

func TestReadStateV2_malformedIssueID(t *testing.T) {
	content := "## Current state\n\n```yaml\nstate: task\nissue: I042\n```\n"

	_, err := NewRouterReader().ReadStateV2(content)
	if !errors.Is(err, ErrV2InvalidID) {
		t.Fatalf("ReadStateV2() error = %v, want ErrV2InvalidID", err)
	}
}

func TestReadStateV2_issueMustBeScalar(t *testing.T) {
	content := "## Current state\n\n```yaml\nstate: task\nissue: [I-042]\n```\n"

	_, err := NewRouterReader().ReadStateV2(content)
	if !errors.Is(err, ErrV2Malformed) {
		t.Fatalf("ReadStateV2() error = %v, want ErrV2Malformed", err)
	}
}

func TestReadStateV2_rejectsUnhyphenatedSelections(t *testing.T) {
	tests := []struct {
		name   string
		fields string
	}{
		{"release", "release: R001\nobjective: none\ntask: none"},
		{"objective", "release: none\nobjective: O001\ntask: none"},
		{"task", "release: none\nobjective: O-001\ntask: T001"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := "## Current state\n\n```yaml\nstate: task\n" + tt.fields + "\nnext_action: \"\"\n```\n"
			_, err := NewRouterReader().ReadStateV2(content)
			if !errors.Is(err, ErrV2InvalidID) {
				t.Fatalf("ReadStateV2() error = %v, want ErrV2InvalidID", err)
			}
		})
	}
}

func TestReadStateV2_taskWithoutObjectiveIsNamedDiagnostic(t *testing.T) {
	content := "## Current state\n\n```yaml\nstate: task\nobjective: none\ntask: T-001\nnext_action: \"\"\n```\n"

	_, err := NewRouterReader().ReadStateV2(content)
	if !errors.Is(err, ErrV2InvalidOwnership) {
		t.Fatalf("ReadStateV2() error = %v, want ErrV2InvalidOwnership", err)
	}
}

func TestReadStateV2_missingBlock(t *testing.T) {
	_, err := NewRouterReader().ReadStateV2("# No state block here")
	if err == nil {
		t.Fatal("ReadStateV2() expected error for missing state block")
	}
	if !strings.Contains(err.Error(), "block") {
		t.Errorf("error = %q, want it to name the missing block", err.Error())
	}
}

func TestReadStateV2_malformedYAML(t *testing.T) {
	content := "## Current state\n\n```yaml\nstate: [task\nobjective: none\n```\n"

	_, err := NewRouterReader().ReadStateV2(content)
	if !errors.Is(err, ErrV2Malformed) {
		t.Fatalf("ReadStateV2() error = %v, want ErrV2Malformed", err)
	}
}

func TestReadStateV2_unknownKeyIsNamedDiagnostic(t *testing.T) {
	// Unlike the V1 reader (TestRouterReader_ignoresUnknownKeys), the V2
	// reader decodes strictly: an unrecognized key is a diagnostic, not
	// silently ignored.
	content := "## Current state\n\n```yaml\nstate: task\nobjective: O-001\ntask: T-001\nnext_action: \"\"\nowner: someone\n```\n"

	_, err := NewRouterReader().ReadStateV2(content)
	if !errors.Is(err, ErrV2Malformed) {
		t.Fatalf("ReadStateV2() error = %v, want ErrV2Malformed for an unknown key", err)
	}
}

func TestReadStateV2_v1RouterStateDecodeUnaffected(t *testing.T) {
	// ReadState (V1) must keep decoding unknown keys tolerantly after
	// ReadState and ReadStateV2 started sharing extractStateBlock.
	content := "## Current state\n\n```yaml\nstate: task-building\nrelease: v1\nepic: E01\ntask: E01/T-001\nnext_action: \"Build T-001\"\nowner: someone\n```\n"

	state, err := NewRouterReader().ReadState(content)
	if err != nil {
		t.Fatalf("ReadState() error = %v, want unknown keys ignored", err)
	}
	if state.State != "task-building" {
		t.Errorf("State = %q, want task-building", state.State)
	}
}

func TestReadStateV2_shippedTemplateDecodesCleanly(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "templates", "project-v2", ".savepoint", "router.md"))
	if err != nil {
		t.Fatalf("read shipped V2 router template: %v", err)
	}

	state, err := NewRouterReader().ReadStateV2(string(data))
	if err != nil {
		t.Fatalf("ReadStateV2() on the shipped V2 router template error = %v", err)
	}
	if state.State != RouterPhaseIdea {
		t.Errorf("State = %q, want idea", state.State)
	}
	if state.Objective != "" {
		t.Errorf("Objective = %q, want empty for a fresh project", state.Objective)
	}
	if state.Task != "" {
		t.Errorf("Task = %q, want empty for a fresh project", state.Task)
	}
}
