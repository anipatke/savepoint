package v2

import (
	"strings"
	"testing"

	xansi "github.com/charmbracelet/x/ansi"
	"github.com/opencode/savepoint/internal/data"
)

// nextPanelText is the Next area's content as one string, which is what every
// wording assertion below is made against.
func nextPanelText(next data.Next) string {
	return strings.Join(nextLines(next), "\n")
}

// TestNextPanelNamesTheTasksStageIdentityAndTitle proves the compact panel's
// whole content for a selected Task: its lifecycle word, its T### identity,
// and its title, in that order, on one line — nothing else. "Audit" never
// appears; a Task at that stage reads "Check" (see taskStageWord/stageLabel).
func TestNextPanelNamesTheTasksStageIdentityAndTitle(t *testing.T) {
	cases := []struct {
		name string
		task *data.TaskV2
		want string
	}{
		{"planned", &data.TaskV2{ID: "T001", Title: "Do the thing", Status: data.ColumnPlanned}, "Planned T001 — Do the thing"},
		{"build", &data.TaskV2{ID: "T002", Title: "Build it", Status: data.ColumnInProgress, Stage: data.StageBuild}, "Build T002 — Build it"},
		{"test", &data.TaskV2{ID: "T003", Title: "Test it", Status: data.ColumnInProgress, Stage: data.StageTest}, "Test T003 — Test it"},
		{"audit stage reads as Check", &data.TaskV2{ID: "T004", Title: "Prove it", Status: data.ColumnInProgress, Stage: data.StageAudit}, "Check T004 — Prove it"},
		{"done", &data.TaskV2{ID: "T005", Title: "Shipped it", Status: data.ColumnDone}, "Done T005 — Shipped it"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := nextPanelText(data.Next{Kind: data.NextExecute, Task: tc.task})
			if got != tc.want {
				t.Errorf("nextPanelText() = %q, want %q", got, tc.want)
			}
			if strings.Contains(got, "AUDIT") || strings.Contains(got, "Audit") {
				t.Errorf("panel names the removed audit word: %q", got)
			}
		})
	}
}

// TestNextPanelFallsBackToTheObjectiveWithNoTaskSelected covers the rungs
// with no Task of their own — Objective integration, planning a first
// Objective — which name the Objective by ID and title instead.
func TestNextPanelFallsBackToTheObjectiveWithNoTaskSelected(t *testing.T) {
	objective := &data.ObjectiveV2{ID: "O001", Title: "Ship the board"}
	got := nextPanelText(data.Next{Kind: data.NextObjectiveIntegration, Objective: objective})
	if want := "O001 — Ship the board"; got != want {
		t.Errorf("nextPanelText() = %q, want %q", got, want)
	}
}

// TestNextPanelPendingMigrationNamesTheOperation covers rung one: no Task or
// Objective is named while a migration operation is in progress, and the
// panel still says which one.
func TestNextPanelPendingMigrationNamesTheOperation(t *testing.T) {
	got := nextPanelText(data.Next{Kind: data.NextPendingMigration, Migration: data.MigrationState{Pending: true, OperationID: "OP001"}})
	if !strings.Contains(got, "OP001") {
		t.Errorf("panel does not name the pending operation:\n%s", got)
	}
}

// TestNextPanelNothingSelectedYet covers a fresh project with no Objective
// and no Task at all: the panel says so plainly rather than rendering blank
// or a raw rung identifier.
func TestNextPanelNothingSelectedYet(t *testing.T) {
	got := nextPanelText(data.Next{Kind: data.NextPlanObjective})
	if got == "" || strings.Contains(got, "NextPlanObjective") {
		t.Errorf("nextPanelText() = %q, want a plain fallback, not a raw rung name", got)
	}
}

// TestNextAreaIgnoresTheSidebarSelection is the panel's independence from
// navigation: the projection answers for the project, and the user needs that
// answer most when they are looking at the wrong part of it.
func TestNextAreaIgnoresTheSidebarSelection(t *testing.T) {
	model := openSizedBoard(t, writeNavigationProject(t), 120, 48)
	before := xansi.Strip(model.renderNext(120))

	for _, objective := range []string{"O001", "O005", ""} {
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

	if !strings.Contains(before, "T001 — Do the thing") {
		t.Fatalf("the Next area does not name the projection's own Task:\n%s", before)
	}

	model.Cards = groupTaskCards(nil)
	model.Objectives = nil
	model.SelectedObjective = ""
	model.State.Index = nil
	if got := xansi.Strip(model.renderNext(100)); got != before {
		t.Errorf("emptying the index and the cards moved the Next area:\nbefore:\n%s\nafter:\n%s", before, got)
	}

	model.State.Next = data.Next{Kind: data.NextPlanObjective}
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
