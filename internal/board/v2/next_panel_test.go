package v2

import (
	"strings"
	"testing"
	"time"

	xansi "github.com/charmbracelet/x/ansi"
	"github.com/opencode/savepoint/internal/data"
)

// nextPanelText is the Next area's content as one string, which is what every
// wording assertion below is made against.
func nextPanelText(next data.Next) string {
	return strings.Join(nextLines(next), "\n")
}

// rungFixtures is one data.Next per rung of the ladder, built in memory so a
// wording assertion names the projection value it is about rather than a
// project state that happens to produce it. The gate decisions and clearances
// here are the shapes the resolvers return for each rung.
func rungFixtures() []data.Next {
	objective := &data.ObjectiveV2{ID: "O001", Title: "Ship the board", Status: data.ColumnInProgress}
	planned := &data.TaskV2{ID: "T001", Title: "Do the thing", Objective: "O001", Status: data.ColumnPlanned}
	audit := &data.TaskV2{ID: "T002", Title: "Prove the thing", Objective: "O001", Status: data.ColumnInProgress, Stage: data.StageAudit}

	return []data.Next{
		{Kind: data.NextPendingMigration, Migration: data.MigrationState{Pending: true, OperationID: "OP001"}},
		{
			Kind: data.NextReplan, Objective: objective, Task: planned,
			GateDecision: &data.GateDecision{Blockers: []data.GateBlocker{{Kind: data.GateBlockReplan, Detail: "the plan no longer matches"}}},
		},
		{
			Kind: data.NextDependency, Objective: objective, Task: planned,
			GateDecision: &data.GateDecision{Blockers: []data.GateBlocker{{
				Kind:       data.GateBlockDependency,
				Dependency: &data.DependencyBlock{Kind: data.DependencyBlockNotDone, Target: "T009"},
			}}},
		},
		{
			Kind: data.NextExecute, Objective: objective, Task: planned,
			GateDecision: &data.GateDecision{Allowed: true},
		},
		{
			Kind: data.NextCheckNeeded, Objective: objective, Task: audit,
			Clearance: &data.Clearance{State: data.ClearanceMissing},
		},
		{
			Kind: data.NextOwnerValidationRequired, Objective: objective, Task: audit,
			Clearance:    &data.Clearance{State: data.ClearanceCurrent, Check: "C001", Freshness: freshnessFixture("C001", data.FreshnessCurrent)},
			GateDecision: &data.GateDecision{Blockers: []data.GateBlocker{{Kind: data.GateBlockOwnerAcceptance}}},
		},
		{
			Kind: data.NextObjectiveIntegration, Objective: objective,
			Clearance:    &data.Clearance{State: data.ClearanceMissing},
			GateDecision: &data.GateDecision{Blockers: []data.GateBlocker{{Kind: data.GateBlockClearanceMissing}}},
		},
		{Kind: data.NextReady, Objective: objective, Task: planned},
		{Kind: data.NextPlanObjective},
	}
}

func freshnessFixture(check string, state data.FreshnessState) *data.Freshness {
	return &data.Freshness{
		State:      state,
		Check:      check,
		AssessedBy: data.Actor{Role: data.ActorRoleChecker, Session: "checker-fixture"},
		AssessedAt: time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
		Basis:      "reran the suite",
	}
}

// TestNextPanelRendersEveryRungDistinctly proves each of the nine rungs
// reaches the panel as its own readable statement: a rung line and an action
// line no other rung produces.
func TestNextPanelRendersEveryRungDistinctly(t *testing.T) {
	seen := make(map[string]data.NextKind, 9)

	for _, next := range rungFixtures() {
		lines := nextLines(next)
		rung, action := lines[0], lines[len(lines)-1]

		if !strings.HasPrefix(rung, "NEXT: ") || rung == "NEXT: " {
			t.Errorf("rung %s renders no leading statement: %q", next.Kind, rung)
		}
		if !strings.HasPrefix(action, "Action: ") || action == "Action: " {
			t.Errorf("rung %s renders no action: %q", next.Kind, action)
		}
		if rung == "NEXT: "+string(next.Kind) {
			t.Errorf("rung %s renders its raw kind rather than wording a reader can act on: %q", next.Kind, rung)
		}

		statement := rung + "\n" + action
		if other, ok := seen[statement]; ok {
			t.Errorf("rungs %s and %s render the same statement %q", other, next.Kind, statement)
		}
		seen[statement] = next.Kind
	}

	if len(seen) != 9 {
		t.Errorf("the panel rendered %d distinct rung statements, want one per rung of the ladder", len(seen))
	}
}

// TestNextPanelNamesTheSelectedRecordsByIDAndTitle covers the identity lines,
// including the rungs that name neither record.
func TestNextPanelNamesTheSelectedRecordsByIDAndTitle(t *testing.T) {
	for _, next := range rungFixtures() {
		got := nextPanelText(next)

		if next.Objective != nil {
			if want := "Objective: O001 — Ship the board"; !strings.Contains(got, want) {
				t.Errorf("rung %s does not name its Objective by ID and title (%q):\n%s", next.Kind, want, got)
			}
		} else if strings.Contains(got, "Objective: ") {
			t.Errorf("rung %s names an Objective it did not select:\n%s", next.Kind, got)
		}

		if next.Task != nil {
			if !strings.Contains(got, "Task: "+next.Task.ID+" — "+next.Task.Title) {
				t.Errorf("rung %s does not name its Task by ID and title:\n%s", next.Kind, got)
			}
		} else if strings.Contains(got, "Task: ") {
			t.Errorf("rung %s names a Task it did not select:\n%s", next.Kind, got)
		}
	}
}

// TestNextPanelSelectionDiagnosticAccompaniesTheAction is the specific
// behavior the projection forbids being turned into a substitution: a router
// naming a record the project does not have is reported as itself, beside the
// action the project's own records still support.
func TestNextPanelSelectionDiagnosticAccompaniesTheAction(t *testing.T) {
	next := data.Next{
		Kind: data.NextReady,
		Task: &data.TaskV2{ID: "T001", Title: "Do the thing", Objective: "O001", Status: data.ColumnPlanned},
		SelectionDiagnostic: &data.SelectionDiagnostic{
			Kind: data.SelectionNotFound, RecordKind: data.SelectionRecordTask, ID: "T900",
		},
	}

	got := nextPanelText(next)

	if !strings.Contains(got, "Selection: ") || !strings.Contains(got, "T900") {
		t.Errorf("the panel does not name the record the router named:\n%s", got)
	}
	if !strings.Contains(got, "Action: Start Task T001.") {
		t.Errorf("the panel dropped the available next action for a stale router line:\n%s", got)
	}
	// T001 is named as the action's own target, never as a stand-in for T900:
	// the diagnostic line itself must carry no substitute.
	for _, line := range nextLines(next) {
		if strings.HasPrefix(line, "Selection: ") && strings.Contains(line, "T001") {
			t.Errorf("the selection diagnostic offers a similarly numbered record in place of T900: %q", line)
		}
	}
}

// TestNextPanelSelectionMismatchNamesBothRecords covers the second diagnostic
// kind: a Task that resolves but disagrees with the router about its owner.
func TestNextPanelSelectionMismatchNamesBothRecords(t *testing.T) {
	got := nextPanelText(data.Next{
		Kind: data.NextPlanObjective,
		SelectionDiagnostic: &data.SelectionDiagnostic{
			Kind: data.SelectionMismatch, RouterObjective: "O001", Task: "T001", TaskObjective: "O002",
		},
	})

	for _, want := range []string{"Selection: ", "O001", "T001", "O002"} {
		if !strings.Contains(got, want) {
			t.Errorf("the mismatch diagnostic does not name %q:\n%s", want, got)
		}
	}
	if !strings.Contains(got, "Action: ") {
		t.Errorf("the panel dropped the next action for a mismatched selection:\n%s", got)
	}
}

// TestNextPanelClearanceStatesReadDistinctly proves missing, unknown, and
// stale are three different statements rather than three shades of "not
// clear", and that a state resting on a recorded freshness assessment names
// the Check and the basis that assessment recorded.
func TestNextPanelClearanceStatesReadDistinctly(t *testing.T) {
	task := &data.TaskV2{ID: "T002", Title: "Prove the thing", Objective: "O001", Status: data.ColumnInProgress, Stage: data.StageAudit}

	cases := []struct {
		state     data.ClearanceState
		clearance *data.Clearance
		wantParts []string
	}{
		{
			state:     data.ClearanceMissing,
			clearance: &data.Clearance{State: data.ClearanceMissing},
			wantParts: []string{"No Check has ever been recorded"},
		},
		{
			state:     data.ClearanceUnknown,
			clearance: &data.Clearance{State: data.ClearanceUnknown, Check: "C001"},
			wantParts: []string{"C001", "no freshness assessment has ever been recorded", "unknown"},
		},
		{
			state:     data.ClearanceStale,
			clearance: &data.Clearance{State: data.ClearanceStale, Check: "C001", Freshness: freshnessFixture("C001", data.FreshnessStale)},
			wantParts: []string{"C001", "stale", "reran the suite", "2026-01-02"},
		},
		{
			state:     data.ClearanceNeedsWork,
			clearance: &data.Clearance{State: data.ClearanceNeedsWork, Check: "C001"},
			wantParts: []string{"C001", "NEEDS WORK"},
		},
	}

	seen := make(map[string]data.ClearanceState, len(cases))
	for _, test := range cases {
		got := nextPanelText(data.Next{Kind: data.NextCheckNeeded, Task: task, Clearance: test.clearance})

		for _, part := range test.wantParts {
			if !strings.Contains(got, part) {
				t.Errorf("clearance %s does not report %q:\n%s", test.state, part, got)
			}
		}

		clearanceLine := lineWithPrefix(t, nextLines(data.Next{Kind: data.NextCheckNeeded, Task: task, Clearance: test.clearance}), "Technical clearance: ")
		if other, ok := seen[clearanceLine]; ok {
			t.Errorf("clearance states %s and %s render identically: %q", other, test.state, clearanceLine)
		}
		seen[clearanceLine] = test.state
	}
}

// TestNextPanelNamesAnExceptionAsAnException proves a completion a recorded
// exception allows is never reported as a CLEAR result or as current
// clearance — the distinction the release design requires to stay visible.
func TestNextPanelNamesAnExceptionAsAnException(t *testing.T) {
	got := nextPanelText(data.Next{
		Kind: data.NextExecute,
		Task: &data.TaskV2{ID: "T001", Title: "Do the thing", Objective: "O001", Status: data.ColumnInProgress, Stage: data.StageAudit},
		GateDecision: &data.GateDecision{
			Allowed:            true,
			AllowedByException: true,
			Exception: &data.Exception{
				Owner: "the owner", Check: "C004", Reason: "shipped with a known gap",
			},
		},
	})

	for _, want := range []string{"exception", "the owner", "C004", "shipped with a known gap"} {
		if !strings.Contains(got, want) {
			t.Errorf("the panel does not name the recorded exception %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "CLEAR") || strings.Contains(got, "clearance is current") {
		t.Errorf("the panel reports an exception as a clearance result:\n%s", got)
	}
}

// TestNextPanelClaimsNoVerification scans every rung for language claiming the
// board verified, ran, checked, or confirmed anything. The board read records;
// it ran nothing. It mirrors internal/resume's own assertion, because the two
// surfaces now state the same facts and must be held to the same standard.
func TestNextPanelClaimsNoVerification(t *testing.T) {
	forbidden := []string{"verified", "verify", "confirmed", "confirm", "checked ", " ran ", "board ran"}

	for _, next := range rungFixtures() {
		text := strings.ToLower(nextPanelText(next))
		for _, word := range forbidden {
			if strings.Contains(text, word) {
				t.Errorf("rung %s contains the verification claim %q:\n%s", next.Kind, word, text)
			}
		}
	}
}

// TestNextPanelSummarisesIssuesByCountAndType proves the follow-ups relevant
// to the selection are visible without the overlay being open, and that a rung
// carrying none renders no Issues line at all rather than an empty one.
func TestNextPanelSummarisesIssuesByCountAndType(t *testing.T) {
	next := data.Next{
		Kind: data.NextReady,
		Task: &data.TaskV2{ID: "T001", Title: "Do the thing", Objective: "O001", Status: data.ColumnPlanned},
		Issues: []*data.IssueV2{
			{ID: "I001", Title: "One", Type: data.IssueTypeDrift, Status: data.IssueStatusOpen},
			{ID: "I002", Title: "Two", Type: data.IssueTypeGuardrail, Status: data.IssueStatusOpen},
			{ID: "I003", Title: "Three", Type: data.IssueTypeDrift, Status: data.IssueStatusInProgress},
		},
	}

	line := lineWithPrefix(t, nextLines(next), "Issues: ")
	if want := "Issues: 3 — 2 drift, 1 guardrail"; line != want {
		t.Errorf("issues summary = %q, want %q", line, want)
	}

	next.Issues = nil
	if strings.Contains(nextPanelText(next), "Issues:") {
		t.Errorf("a selection with no Issues renders an empty Issues line:\n%s", nextPanelText(next))
	}
}

// TestNextPanelPendingMigrationOutranksAndNamesTheOperation covers rung one
// end to end, over a real project held back by an incomplete conversion.
func TestNextPanelPendingMigrationOutranksAndNamesTheOperation(t *testing.T) {
	root := writeValidProject(t)
	operationID := createPendingOperation(t, root)

	got := xansi.Strip(openBoard(t, root, "").View())

	if !strings.Contains(got, "NEXT: Finish the pending migration") {
		t.Errorf("the Next area does not lead with the pending migration:\n%s", got)
	}
	if !strings.Contains(got, operationID) {
		t.Errorf("the Next area does not report the operation %q:\n%s", operationID, got)
	}
	// The project has a router selecting T001, which every lower rung would
	// have named. Rung one outranking them means no record is named at all.
	if strings.Contains(got, "Task: T001") {
		t.Errorf("a lower rung's record reached the Next area past a pending migration:\n%s", got)
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

	if !strings.Contains(before, "Task: T001 — Do the thing") {
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
	if !strings.Contains(after, "NEXT: Plan an Objective") {
		t.Errorf("the Next area does not track the projection it was given:\n%s", after)
	}
}

// TestNextAreaReportsTheRungTheLoadResolved closes the loop between the panel's
// rung wording and the projection the load command actually resolved, over a
// real project rather than an in-memory value.
func TestNextAreaReportsTheRungTheLoadResolved(t *testing.T) {
	model := openBoard(t, writeValidProject(t), "")

	if model.State.Next.Kind != data.NextExecute {
		t.Fatalf("Next.Kind = %q, want the fixture's execute rung", model.State.Next.Kind)
	}
	if want := "NEXT: " + nextKindLabel(model.State.Next.Kind); !strings.Contains(xansi.Strip(model.View()), want) {
		t.Errorf("the board does not report the rung it resolved (%q):\n%s", want, xansi.Strip(model.View()))
	}
}

// TestNextAreaOpensAFreshProjectWithNoError covers the project a user gets from
// `savepoint init`: no Objectives, no Tasks, and a planning next action rather
// than anything that reads as broken.
func TestNextAreaOpensAFreshProjectWithNoError(t *testing.T) {
	got := xansi.Strip(openBoard(t, writeEmptyProjectFromTemplate(t), "").View())

	if !strings.Contains(got, "NEXT: Plan an Objective") {
		t.Errorf("a fresh project does not render the planning next action:\n%s", got)
	}
	if !strings.Contains(got, "Action: Plan the next Objective") {
		t.Errorf("a fresh project does not render the planning action:\n%s", got)
	}
	for _, forbidden := range []string{diagnosticHeading, "Selection: ", "Issues:"} {
		if strings.Contains(got, forbidden) {
			t.Errorf("a fresh project's Next area contains %q, which reads as a problem:\n%s", forbidden, got)
		}
	}
}

// lineWithPrefix returns the one line starting with prefix, failing when none
// does — so an assertion about a line's exact content cannot pass vacuously.
func lineWithPrefix(t *testing.T, lines []string, prefix string) string {
	t.Helper()
	for _, line := range lines {
		if strings.HasPrefix(line, prefix) {
			return line
		}
	}
	t.Fatalf("no line starts with %q:\n%s", prefix, strings.Join(lines, "\n"))
	return ""
}
