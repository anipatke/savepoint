package v2

import (
	"github.com/charmbracelet/lipgloss"
	"strings"
	"testing"

	xansi "github.com/charmbracelet/x/ansi"
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/resume"
	"github.com/opencode/savepoint/internal/styles"
)

// nextPanelText is the Next area's content as one string, which is what every
// wording assertion below is made against.
func nextPanelText(next data.Next) string {
	return strings.Join(nextLines(next), "\n")
}

// TestNextPanelNamesTheTasksStageIdentityAndTitle proves the compact panel's
// whole content for a selected Task: the Objective and Task lifecycle words,
// identities, and title, in that order, on one line. "Audit" never appears;
// a Task at that stage reads "Check".
func TestNextPanelNamesTheTasksStageIdentityAndTitle(t *testing.T) {
	cases := []struct {
		name      string
		objective *data.ObjectiveV2
		task      *data.TaskV2
		want      string
	}{
		{"planned", &data.ObjectiveV2{ID: "O-001", Title: "Planned objective", Status: data.ColumnPlanned}, &data.TaskV2{ID: "T-001", Objective: "O-001", Title: "Do the thing", Status: data.ColumnPlanned}, "Start T-001 — Do the thing (O-001)"},
		{"build", &data.ObjectiveV2{ID: "O-002", Title: "Active objective", Status: data.ColumnInProgress}, &data.TaskV2{ID: "T-002", Objective: "O-002", Title: "Build it", Status: data.ColumnInProgress, Stage: data.StageBuild}, "Build T-002 — Build it (O-002)"},
		{"test", &data.ObjectiveV2{ID: "O-003", Title: "Active objective", Status: data.ColumnInProgress}, &data.TaskV2{ID: "T-003", Objective: "O-003", Title: "Test it", Status: data.ColumnInProgress, Stage: data.StageTest}, "Test T-003 — Test it (O-003)"},
		{"audit stage reads as Check", &data.ObjectiveV2{ID: "O-004", Title: "Active objective", Status: data.ColumnInProgress}, &data.TaskV2{ID: "T-004", Objective: "O-004", Title: "Prove it", Status: data.ColumnInProgress, Stage: data.StageAudit}, "Close T-004 — Prove it (O-004)"},
		{"done", &data.ObjectiveV2{ID: "O-005", Title: "Completed objective", Status: data.ColumnDone}, &data.TaskV2{ID: "T-005", Objective: "O-005", Title: "Shipped it", Status: data.ColumnDone}, "Done T-005 — Shipped it (O-005)"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := nextPanelText(data.Next{Kind: data.NextExecute, Objective: tc.objective, Task: tc.task})
			if got != tc.want {
				t.Errorf("nextPanelText() = %q, want %q", got, tc.want)
			}
			if strings.Contains(got, "AUDIT") || strings.Contains(got, "Audit") {
				t.Errorf("panel names the removed audit word: %q", got)
			}
		})
	}
}

// TestNextPanelFallsBackToTheObjectiveWithNoTaskSelected covers Objective
// integration, where every owned Task is done and Next names the Check.
func TestNextPanelFallsBackToTheObjectiveWithNoTaskSelected(t *testing.T) {
	objective := &data.ObjectiveV2{ID: "O-001", Title: "Ship the board", Status: data.ColumnInProgress}
	got := nextPanelText(data.Next{Kind: data.NextObjectiveIntegration, Objective: objective})
	if want := "Check O-001 — Ship the board"; got != want {
		t.Errorf("nextPanelText() = %q, want %q", got, want)
	}
}

// TestNextPanelNothingSelected covers a fresh project with no Objective
// and no Task at all: the panel says so plainly rather than rendering blank
// or a raw rung identifier.
func TestNextPanelNothingSelected(t *testing.T) {
	got := nextPanelText(data.Next{Kind: data.NextNothingSelected})
	if got != "Nothing selected" {
		t.Errorf("nextPanelText() = %q, want Nothing selected", got)
	}
}

func TestNextPanelNamesIssueAndShowsItAsContext(t *testing.T) {
	issue := &data.IssueV2{ID: "I-042", Title: "Repair the parser", Status: data.IssueStatusOpen}
	issueOnly := data.Next{Kind: data.NextIssue, Issue: issue}
	if got, want := nextPanelText(issueOnly), "Fix I-042 — Repair the parser"; got != want {
		t.Errorf("nextPanelText() = %q, want Issue line %q", got, want)
	}

	taskContext := data.Next{
		Kind:      data.NextExecute,
		Objective: &data.ObjectiveV2{ID: "O-014", Status: data.ColumnInProgress},
		Task:      &data.TaskV2{ID: "T-028", Objective: "O-014", Title: "Copy the line", Status: data.ColumnInProgress, Stage: data.StageBuild},
		Issue:     issue,
	}
	want := "Build T-028 — Copy the line (O-014)\nIssue: Fix I-042 — Repair the parser"
	if got := nextPanelText(taskContext); got != want {
		t.Errorf("nextPanelText() = %q, want Task line plus selected Issue context %q", got, want)
	}
	model := Model{}
	model.State.Next = taskContext
	if got := xansi.Strip(model.renderNext(120)); !strings.Contains(got, "Issue: Fix I-042 — Repair the parser") {
		t.Errorf("renderNext() = %q, want selected Issue context", got)
	}
}

func TestNextPanelAddsSelectionDiagnosticUnderTheExistingNextLine(t *testing.T) {
	cases := []struct {
		name       string
		next       data.Next
		wantPhrase string
	}{
		{
			name: "existing not-found wording",
			next: data.Next{
				Kind: data.NextNothingSelected,
				SelectionDiagnostic: &data.SelectionDiagnostic{
					Kind: data.SelectionNotFound, RecordKind: data.SelectionRecordTask, ID: "T-999",
				},
			},
			wantPhrase: "The router names task T-999, which does not exist among the project's live records.",
		},
		{
			name: "finished Task wording",
			next: data.Next{
				Kind:      data.NextSelectTask,
				Objective: &data.ObjectiveV2{ID: "O-001", Title: "Ship it", Status: data.ColumnInProgress},
				SelectionDiagnostic: &data.SelectionDiagnostic{
					Kind: data.SelectionDone, RecordKind: data.SelectionRecordTask, ID: "T-001",
				},
			},
			wantPhrase: "Warning: router still selects finished Task T-001.",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			want := resume.NextLine(tc.next) + "\n" + tc.wantPhrase
			if got := nextPanelText(tc.next); got != want {
				t.Errorf("nextPanelText() = %q, want Next plus diagnostic %q", got, want)
			}
		})
	}
}

func TestRenderNextIncludesTheSelectionDiagnostic(t *testing.T) {
	model := openBoard(t, writeValidProject(t), "")
	model.State.Next = data.Next{
		Kind:      data.NextSelectTask,
		Objective: &data.ObjectiveV2{ID: "O-001", Title: "Ship it", Status: data.ColumnInProgress},
		SelectionDiagnostic: &data.SelectionDiagnostic{
			Kind: data.SelectionDone, RecordKind: data.SelectionRecordTask, ID: "T-001",
		},
	}
	got := xansi.Strip(model.renderNext(120))
	if want := "Warning: router still selects finished Task T-001."; !strings.Contains(got, want) {
		t.Errorf("renderNext() = %q, want the shared selection diagnostic %q", got, want)
	}
}

func TestRenderNextAccentsOnlyTheVerb(t *testing.T) {
	tests := []struct {
		name  string
		next  data.Next
		style lipgloss.Style
		verb  string
		rest  string
	}{
		{
			name: "build",
			next: data.Next{
				Kind:      data.NextExecute,
				Objective: &data.ObjectiveV2{ID: "O-014", Status: data.ColumnInProgress},
				Task:      &data.TaskV2{ID: "T-028", Title: "Copy the line", Objective: "O-014", Status: data.ColumnInProgress, Stage: data.StageBuild},
			},
			style: styles.FooterPhaseTask, verb: "Build", rest: " T-028 — Copy the line (O-014)",
		},
		{
			name:  "objective check",
			next:  data.Next{Kind: data.NextObjectiveIntegration, Objective: &data.ObjectiveV2{ID: "O-015", Title: "Finish the Objective", Status: data.ColumnInProgress}},
			style: styles.FooterPhaseCheck, verb: "Check", rest: " O-015 — Finish the Objective",
		},
		{
			name:  "objective close",
			next:  data.Next{Kind: data.NextObjectiveReady, Objective: &data.ObjectiveV2{ID: "O-016", Title: "Ready to close", Status: data.ColumnInProgress}},
			style: styles.FooterPhaseCheck, verb: "Close", rest: " O-016 — Ready to close",
		},
		{
			name:  "plan",
			next:  data.Next{Kind: data.NextPlanObjective, Objective: &data.ObjectiveV2{ID: "O-017", Title: "Break it down", Status: data.ColumnPlanned}},
			style: styles.HeaderWhiteBold, verb: "Plan", rest: " O-017 — Break it down",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			want := test.style.Render(test.verb) + styles.HeaderWhiteBold.Render(test.rest)
			if got := styleNextLine(test.next, resume.NextLine(test.next)); got != want {
				t.Errorf("styled line = %q, want %q", got, want)
			}
		})
	}
}

// TestNextAreaIgnoresTheSidebarSelection is the panel's independence from
// navigation: the projection answers for the project, and the user needs that
// answer most when they are looking at the wrong part of it.
func TestNextAreaIgnoresTheSidebarSelection(t *testing.T) {
	model := openSizedBoard(t, writeNavigationProject(t), 120, 48)
	before := xansi.Strip(model.renderNext(120))

	for _, objective := range []string{"O-001", "O-005", ""} {
		model.selectObjective(objective)
		if got := xansi.Strip(model.renderNext(120)); got != before {
			t.Errorf("selecting %q moved the Next area:\nbefore:\n%s\nafter:\n%s", objective, before, got)
		}
	}
}

// TestNextAreaTracksOnlyTheProjection is the derive-nothing claim as behavior:
// replacing the records the board is holding leaves the Next area unchanged,
// and replacing the projection alone changes it. A panel that consulted the
// index to second-guess the answer would fail the first half; one that ignored
// the projection would fail the second.
func TestNextAreaTracksOnlyTheProjection(t *testing.T) {
	model := openBoard(t, writeValidProject(t), "")
	before := xansi.Strip(model.renderNext(100))

	if want := "NEXT: " + resume.NextLine(model.State.Next); !strings.Contains(before, want) {
		t.Fatalf("the Next area does not name the projection's own Task:\n%s", before)
	}

	model.Cards = groupTaskCards(nil)
	model.Objectives = nil
	model.SelectedObjective = ""
	model.State.Index = nil
	if got := xansi.Strip(model.renderNext(100)); got != before {
		t.Errorf("emptying the index and the cards moved the Next area:\nbefore:\n%s\nafter:\n%s", before, got)
	}

	model.State.Next = data.Next{Kind: data.NextNothingSelected}
	after := xansi.Strip(model.renderNext(100))
	if after == before {
		t.Errorf("replacing the projection left the Next area unchanged:\n%s", after)
	}
}

// TestNextAreaReportsTheTaskTheLoadResolved closes the loop between the
// panel's wording and the projection the load command actually resolved,
// over a real project rather than an in-memory value.
func TestNextAreaReportsTheTaskTheLoadResolved(t *testing.T) {
	model := openBoard(t, writeValidProject(t), "")

	if model.State.Next.Kind != data.NextExecute {
		t.Fatalf("Next.Kind = %q, want the fixture's execute rung", model.State.Next.Kind)
	}
	if model.State.Next.Task == nil {
		t.Fatalf("Next.Task is nil, want the fixture's resolved Task")
	}
	want := nextPanelText(model.State.Next)
	if !strings.Contains(xansi.Strip(model.View()), want) {
		t.Errorf("the board does not report the Task it resolved (%q):\n%s", want, xansi.Strip(model.View()))
	}
}

// TestNextAreaOpensAFreshProjectWithNoError covers the project a user gets from
// `savepoint init`: no Objectives, no Tasks, and a plain fallback rather than
// anything that reads as broken.
func TestNextAreaOpensAFreshProjectWithNoError(t *testing.T) {
	got := xansi.Strip(openBoard(t, writeEmptyProjectFromTemplate(t), "").View())

	for _, forbidden := range []string{diagnosticHeading, "NextPlanObjective"} {
		if strings.Contains(got, forbidden) {
			t.Errorf("a fresh project's Next area contains %q, which reads as a problem:\n%s", forbidden, got)
		}
	}
}
