package v2

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunWithoutTTYLeadsWithNextAndReportsCounts(t *testing.T) {
	root := writeValidProject(t)
	var stdout bytes.Buffer

	if err := Run(Options{Root: root, Stdout: &stdout, TTY: false}); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	got := stdout.String()
	if want := "Planned T001 — Do the thing"; !strings.Contains(got, want) {
		t.Errorf("output does not lead with the resolved Task, missing %q:\n%s", want, got)
	}
	if !strings.Contains(got, "Objectives: 1  Tasks: 2") {
		t.Errorf("output does not report the loaded counts:\n%s", got)
	}
	if !strings.Contains(got, "PLANNED      1") || !strings.Contains(got, "DONE         1") {
		t.Errorf("output does not report the three columns:\n%s", got)
	}
	if strings.Contains(got, "\x1b[") {
		t.Errorf("output carries ANSI escapes without a TTY:\n%q", got)
	}
}

func TestRunWithoutTTYStripsAuthoredTerminalControls(t *testing.T) {
	root := writeValidProject(t)
	writeTask(t, root, "O001", "T001", `\u001b[31mred\u001b[0m\tname`, "status: planned\n")

	var stdout bytes.Buffer
	if err := Run(Options{Root: root, Stdout: &stdout, TTY: false}); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if strings.Contains(stdout.String(), "\x1b") {
		t.Fatalf("plain output contains authored ANSI escapes: %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "red name") {
		t.Errorf("plain output did not preserve readable title text: %q", stdout.String())
	}
}

func TestRunWithoutTTYIsDeterministic(t *testing.T) {
	root := writeValidProject(t)

	var first, second bytes.Buffer
	for _, out := range []*bytes.Buffer{&first, &second} {
		if err := Run(Options{Root: root, Stdout: out, TTY: false}); err != nil {
			t.Fatalf("Run() error = %v", err)
		}
	}

	if first.String() != second.String() {
		t.Errorf("two runs over the same project differ:\n%q\n%q", first.String(), second.String())
	}
}

func TestRunWithoutTTYIncludesTitlesBadgesAndIssueSummary(t *testing.T) {
	root := writeEvidenceProject(t)
	var stdout bytes.Buffer

	if err := Run(Options{Root: root, Stdout: &stdout, TTY: false}); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	got := stdout.String()
	for _, want := range []string{
		"Rechecked after the first run found problems",
		"[◆ CHECK  [✓] Check",
		"Issues: 1 — 1 defect",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("plain output missing %q:\n%s", want, got)
		}
	}
}

// TestRunWithoutTTYRefusedIndexFailsAndWritesNoBoard proves every refusal the
// screen reports also fails the piped board: the diagnostic is returned — the
// command prints it to stderr and exits nonzero — and stdout stays empty.
func TestRunWithoutTTYRefusedIndexFailsAndWritesNoBoard(t *testing.T) {
	for _, test := range invalidProjectCases() {
		t.Run(test.Name, func(t *testing.T) {
			root := writeValidProject(t)
			test.Build(t, root)
			var stdout bytes.Buffer

			err := Run(Options{Root: root, Stdout: &stdout, TTY: false})

			if err == nil {
				t.Fatal("Run() error = nil, want the load diagnostic reported as a failure")
			}
			if !strings.Contains(err.Error(), test.WantPath) {
				t.Errorf("error = %q, want the offending file path %q named", err.Error(), test.WantPath)
			}
			for _, part := range test.WantParts {
				if !strings.Contains(err.Error(), part) {
					t.Errorf("error = %q, want it to name %q", err.Error(), part)
				}
			}
			if stdout.Len() != 0 {
				t.Errorf("stdout = %q, want no partial board written", stdout.String())
			}
		})
	}
}

func TestRunObjectiveFilterSelectsAndRejects(t *testing.T) {
	root := writeValidProject(t)

	var stdout bytes.Buffer
	if err := Run(Options{Root: root, ObjectiveFilter: "O001", Stdout: &stdout, TTY: false}); err != nil {
		t.Fatalf("Run() with a known objective error = %v", err)
	}

	var rejected bytes.Buffer
	err := Run(Options{Root: root, ObjectiveFilter: "O002", Stdout: &rejected, TTY: false})
	if err == nil {
		t.Fatal("Run() error = nil, want an unknown --objective refused")
	}
	if !strings.Contains(err.Error(), "--objective O002") {
		t.Errorf("error = %q, want the flag and the value it was given named", err.Error())
	}
	if strings.Contains(err.Error(), "O001") {
		t.Errorf("error = %q, want nothing guessed at in place of the named objective", err.Error())
	}
	if rejected.Len() != 0 {
		t.Errorf("stdout = %q, want no board written for a refused filter", rejected.String())
	}
}

// TestRunWithoutTTYCountsOnlyTheSelectedObjectivesTasks proves the piped board
// reads the same selection the terminal board does: the router's Objective, or
// the flag when one was passed, filtering the counts either way.
func TestRunWithoutTTYCountsOnlyTheSelectedObjectivesTasks(t *testing.T) {
	root := writeNavigationProject(t)

	var fromRouter bytes.Buffer
	if err := Run(Options{Root: root, Stdout: &fromRouter, TTY: false}); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	for _, want := range []string{"Tasks: 4", "Selected: O003", "PLANNED      2", "DONE         0"} {
		if !strings.Contains(fromRouter.String(), want) {
			t.Errorf("output is missing %q:\n%s", want, fromRouter.String())
		}
	}

	var fromFlag bytes.Buffer
	if err := Run(Options{Root: root, ObjectiveFilter: "O001", Stdout: &fromFlag, TTY: false}); err != nil {
		t.Fatalf("Run() with --objective error = %v", err)
	}
	for _, want := range []string{"Selected: O001", "PLANNED      0", "DONE         1"} {
		if !strings.Contains(fromFlag.String(), want) {
			t.Errorf("output does not count only the flagged Objective's Tasks, missing %q:\n%s", want, fromFlag.String())
		}
	}
}

// TestUpdateUnknownObjectiveFilterQuitsWithError proves the terminal path
// refuses the same filter the plain path refuses: the model records the error
// and asks to quit, and Run hands that error back to the command.
func TestUpdateUnknownObjectiveFilterQuitsWithError(t *testing.T) {
	root := writeValidProject(t)

	model := openBoard(t, root, "O002")

	if model.FatalErr == nil {
		t.Fatal("FatalErr = nil, want the unknown --objective refused")
	}
	if !strings.Contains(model.FatalErr.Error(), "--objective O002") {
		t.Errorf("FatalErr = %q, want the flag named", model.FatalErr.Error())
	}
}

func TestUpdateObjectiveFilterSelectsThatObjective(t *testing.T) {
	root := writeValidProject(t)

	model := openBoard(t, root, "O001")

	if model.FatalErr != nil {
		t.Fatalf("FatalErr = %v, want a known objective accepted", model.FatalErr)
	}
	if model.SelectedObjective != "O001" {
		t.Errorf("SelectedObjective = %q, want O001", model.SelectedObjective)
	}
}

// TestUpdateRouterSelectionThatNoLongerResolvesSelectsNothing proves the board
// selects no substitute for a router-named Objective that is gone, leaving the
// projection's own selection diagnostic as the report.
func TestUpdateRouterSelectionThatNoLongerResolvesSelectsNothing(t *testing.T) {
	root := writeValidProject(t)
	writeRouter(t, root, "design", "O009", "")

	model := openBoard(t, root, "")

	if model.SelectedObjective != "" {
		t.Errorf("SelectedObjective = %q, want nothing selected for a missing record", model.SelectedObjective)
	}
	if model.State.Next.SelectionDiagnostic == nil {
		t.Error("Next.SelectionDiagnostic = nil, want the unresolved selection reported")
	}
}
