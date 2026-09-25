package resume

import (
	"bytes"
	"errors"
	"go/build"
	"strings"
	"testing"
	"time"

	"github.com/opencode/savepoint/internal/data"
)

func TestNextLineFormatsEverySelectionShape(t *testing.T) {
	currentOwnerWait := data.Next{
		Kind:      data.NextObjectiveIntegration,
		Objective: &data.ObjectiveV2{ID: "O-009", Title: "Await owner acceptance", Status: data.ColumnInProgress},
		Clearance: &data.Clearance{State: data.ClearanceCurrent, Check: "C-009"},
		GateDecision: &data.GateDecision{Blockers: []data.GateBlocker{
			{Kind: data.GateBlockOwnerAcceptance},
		}},
	}
	cases := []struct {
		name string
		next data.Next
		want string
	}{
		{
			name: "planned task",
			next: data.Next{Kind: data.NextExecute, Objective: &data.ObjectiveV2{ID: "O-001", Status: data.ColumnPlanned}, Task: &data.TaskV2{ID: "T-001", Objective: "O-001", Title: "Plan it", Status: data.ColumnPlanned}},
			want: "Start T-001 — Plan it (O-001)",
		},
		{
			name: "build task",
			next: data.Next{Kind: data.NextExecute, Objective: &data.ObjectiveV2{ID: "O-002", Status: data.ColumnInProgress}, Task: &data.TaskV2{ID: "T-002", Objective: "O-002", Title: "Build it", Status: data.ColumnInProgress, Stage: data.StageBuild}},
			want: "Build T-002 — Build it (O-002)",
		},
		{
			name: "test task",
			next: data.Next{Kind: data.NextExecute, Objective: &data.ObjectiveV2{ID: "O-003", Status: data.ColumnInProgress}, Task: &data.TaskV2{ID: "T-003", Objective: "O-003", Title: "Test it", Status: data.ColumnInProgress, Stage: data.StageTest}},
			want: "Test T-003 — Test it (O-003)",
		},
		{
			name: "audit task says check",
			next: data.Next{Kind: data.NextCheckNeeded, Objective: &data.ObjectiveV2{ID: "O-004", Status: data.ColumnInProgress}, Task: &data.TaskV2{ID: "T-004", Objective: "O-004", Title: "Review it", Status: data.ColumnInProgress, Stage: data.StageAudit}},
			want: "Check T-004 — Review it (O-004)",
		},
		{
			name: "done task",
			next: data.Next{Kind: data.NextExecute, Objective: &data.ObjectiveV2{ID: "O-005", Status: data.ColumnDone}, Task: &data.TaskV2{ID: "T-005", Objective: "O-005", Title: "Ship it", Status: data.ColumnDone}},
			want: "Done T-005 — Ship it (O-005)",
		},
		{
			name: "task with missing objective record",
			next: data.Next{Kind: data.NextExecute, Task: &data.TaskV2{ID: "T-006", Objective: "O-999", Title: "Keep visible", Status: data.ColumnInProgress, Stage: data.StageBuild}},
			want: "Build T-006 — Keep visible",
		},
		{
			name: "objective with unfinished tasks",
			next: data.Next{Kind: data.NextSelectTask, Objective: &data.ObjectiveV2{ID: "O-007", Title: "Plan tasks", Status: data.ColumnPlanned}},
			want: "Pick a Task in O-007 — Plan tasks",
		},
		{
			name: "objective with no tasks",
			next: data.Next{Kind: data.NextPlanObjective, Objective: &data.ObjectiveV2{ID: "O-007A", Title: "Break this down", Status: data.ColumnPlanned}},
			want: "Plan O-007A — Break this down",
		},
		{
			name: "objective check needed",
			next: data.Next{Kind: data.NextObjectiveIntegration, Objective: &data.ObjectiveV2{ID: "O-008", Title: "Integrate", Status: data.ColumnInProgress}, Clearance: &data.Clearance{State: data.ClearanceMissing}},
			want: "Check O-008 — Integrate",
		},
		{name: "current check awaiting owner", next: currentOwnerWait, want: "Accept O-009 — Await owner acceptance"},
		{
			name: "objective ready for owner closure",
			next: data.Next{Kind: data.NextObjectiveReady, Objective: &data.ObjectiveV2{ID: "O-010", Title: "Ready to close", Status: data.ColumnInProgress}},
			want: "Close O-010 — Ready to close",
		},
		{name: "nothing selected", next: data.Next{Kind: data.NextNothingSelected}, want: "Nothing selected"},
		{name: "open issue", next: data.Next{Kind: data.NextIssue, Issue: &data.IssueV2{ID: "I-042", Title: "Repair the parser", Status: data.IssueStatusOpen}}, want: "Fix I-042 — Repair the parser"},
		{name: "in progress issue", next: data.Next{Kind: data.NextIssue, Issue: &data.IssueV2{ID: "I-043", Title: "Continue the repair", Status: data.IssueStatusInProgress}}, want: "Fix I-043 — Continue the repair"},
		{name: "resolved issue", next: data.Next{Kind: data.NextIssue, Issue: &data.IssueV2{ID: "I-044", Title: "Finished repair", Status: data.IssueStatusResolved}}, want: "Resolved I-044 — Finished repair"},
		{
			name: "planned task blocked",
			next: data.Next{Kind: data.NextDependency, Objective: &data.ObjectiveV2{ID: "O-011"}, Task: &data.TaskV2{ID: "T-011", Title: "Wait", Status: data.ColumnPlanned}},
			want: "Blocked T-011 — Wait (O-011)",
		},
		{
			name: "task replan",
			next: data.Next{Kind: data.NextReplan, Objective: &data.ObjectiveV2{ID: "O-011"}, Task: &data.TaskV2{ID: "T-012", Title: "Rethink", Status: data.ColumnPlanned}},
			want: "Replan T-012 — Rethink (O-011)",
		},
		{
			name: "task awaiting owner acceptance",
			next: data.Next{Kind: data.NextOwnerValidationRequired, Objective: &data.ObjectiveV2{ID: "O-011"}, Task: &data.TaskV2{ID: "T-013", Title: "Accept me", Status: data.ColumnInProgress, Stage: data.StageAudit}},
			want: "Accept T-013 — Accept me (O-011)",
		},
		{
			name: "task without its objective record",
			next: data.Next{Kind: data.NextExecute, Task: &data.TaskV2{ID: "T-014", Title: "Orphan", Status: data.ColumnInProgress, Stage: data.StageBuild}},
			want: "Build T-014 — Orphan",
		},
		{
			name: "current check blocked by more than owner acceptance",
			next: data.Next{
				Kind:         data.NextObjectiveIntegration,
				Objective:    &data.ObjectiveV2{ID: "O-012", Title: "Open issue", Status: data.ColumnInProgress},
				Clearance:    &data.Clearance{State: data.ClearanceCurrent},
				GateDecision: &data.GateDecision{Blockers: []data.GateBlocker{{Kind: data.GateBlockOwnerAcceptance}, {Kind: data.GateBlockObjectiveIssueUnresolved}}},
			},
			want: "Check O-012 — Open issue",
		},
		{
			name: "objective already done",
			next: data.Next{Kind: data.NextObjectiveReady, Objective: &data.ObjectiveV2{ID: "O-013", Title: "Finished", Status: data.ColumnDone}},
			want: "Done O-013 — Finished",
		},
		{name: "goal check needed", next: data.Next{Kind: data.NextReleaseCheckNeeded, Release: &data.ReleaseV2{ID: "R-001", Title: "First delivery"}}, want: "Check R-001 — First delivery"},
		{name: "goal awaiting owner", next: data.Next{Kind: data.NextReleaseOwnerValidationRequired, Release: &data.ReleaseV2{ID: "R-001", Title: "First delivery"}}, want: "Accept R-001 — First delivery"},
		{name: "goal ready", next: data.Next{Kind: data.NextReleaseReady, Release: &data.ReleaseV2{ID: "R-001", Title: "First delivery"}}, want: "Close R-001 — First delivery"},
		{name: "goal with nothing selected", next: data.Next{Kind: data.NextNothingSelected, Release: &data.ReleaseV2{ID: "R-001", Title: "First delivery"}}, want: "Nothing selected"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := NextLine(tc.next); got != tc.want {
				t.Errorf("NextLine() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestNextLine_missingGoalAndObjectiveFacts(t *testing.T) {
	next := data.Next{
		Kind: data.NextNothingSelected,
		SelectionDiagnostic: &data.SelectionDiagnostic{
			Kind: data.SelectionReleaseMissing, RecordKind: data.SelectionRecordRelease,
		},
		ObjectivesWithoutGoal: []string{"O-002", "O-010"},
	}
	if got, want := NextLine(next), "Choose a Goal — press g on the board"; got != want {
		t.Fatalf("NextLine() = %q, want %q", got, want)
	}
	if got := NextVerb(next); got != "Choose" {
		t.Fatalf("NextVerb() = %q, want Choose", got)
	}
	text := renderText(next)
	for _, want := range []string{
		"Selection: The router has no Goal selected.",
		"Objectives without a Goal: O-002, O-010",
		"Next action: Choose a Goal with g on the board.",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("renderText() = %q, want it to contain %q", text, want)
		}
	}
}

func TestRender_validNextHasNoMissingGoalOutput(t *testing.T) {
	text := renderText(data.Next{Kind: data.NextNothingSelected})
	if strings.Contains(text, "Objectives without a Goal") || strings.Contains(text, "Choose a Goal") {
		t.Fatalf("renderText() = %q, want no missing-Goal output for a valid Next", text)
	}
}

func TestRender_flagsUnassignedObjectivesWithGoalSelected(t *testing.T) {
	next := data.Next{
		Kind:    data.NextNothingSelected,
		Release: &data.ReleaseV2{ID: "R-001", Title: "Selected Goal"},
		SelectionDiagnostic: &data.SelectionDiagnostic{
			Kind: data.SelectionObjectiveUnassigned, Release: "R-001", Objective: "O-004",
		},
		ObjectivesWithoutGoal: []string{"O-004"},
	}
	text := renderText(next)
	if !strings.Contains(text, "Objectives without a Goal: O-004") {
		t.Fatalf("renderText() = %q, want Objective O-004 flagged", text)
	}
	if strings.Contains(NextLine(next), "Choose a Goal") {
		t.Errorf("NextLine() = %q, want the selected Goal retained while Objective repair is flagged", NextLine(next))
	}
}

// TestRender_replan is the golden rendering for rung two: a recorded replan
// flag, reported by its own reason.
func TestRender_replan(t *testing.T) {
	next := data.Next{
		Kind: data.NextReplan,
		Task: &data.TaskV2{ID: "T010", Title: "Some task", Status: data.ColumnInProgress, Stage: data.StageBuild},
		GateDecision: &data.GateDecision{Blockers: []data.GateBlocker{
			{Kind: data.GateBlockReplan, Detail: "scope changed"},
		}},
	}
	want := "Task: T010 — Some task\n" +
		"Implementation: Status in_progress, stage build.\n" +
		"\n" +
		"Replan: A replan has been flagged: scope changed\n" +
		"\n" +
		"Next action: Resolve the recorded replan before resuming build, test, or audit.\n"
	assertRenderEquals(t, next, want)
}

// TestRender_dependency is the golden rendering for rung three: an
// unsatisfied Task dependency, naming the target and the unmet requirement.
func TestRender_dependency(t *testing.T) {
	next := data.Next{
		Kind: data.NextDependency,
		Task: &data.TaskV2{ID: "T020", Title: "Blocked task", Status: data.ColumnPlanned},
		GateDecision: &data.GateDecision{Blockers: []data.GateBlocker{
			{Kind: data.GateBlockDependency, Dependency: &data.DependencyBlock{Target: "T019", Kind: data.DependencyBlockNotDone}},
		}},
	}
	want := "Task: T020 — Blocked task\n" +
		"Implementation: Status planned.\n" +
		"\n" +
		"Blocked: Waiting on Task T019, which is not done yet.\n" +
		"\n" +
		"Next action: Wait on the named dependency before starting or advancing this Task.\n"
	assertRenderEquals(t, next, want)
}

// TestRender_objectiveDependency proves an Objective-level dependency block
// is distinguishable from a Task-level one, naming that the wait is at the
// Objective.
func TestRender_objectiveDependency(t *testing.T) {
	next := data.Next{
		Kind: data.NextDependency,
		Task: &data.TaskV2{ID: "T021", Title: "Waits on its objective", Status: data.ColumnPlanned},
		GateDecision: &data.GateDecision{Blockers: []data.GateBlocker{
			{Kind: data.GateBlockObjectiveDependency, ObjectiveDependency: &data.ObjectiveDependencyBlock{Target: "O005", Kind: data.ObjectiveDependencyBlockNotDone}},
		}},
	}
	text := renderText(next)
	if !strings.Contains(text, "Blocked: The owning Objective is waiting: Objective O005 is not done yet.") {
		t.Fatalf("renderText() = %q, want it to name the Objective-level wait", text)
	}
}

// TestRender_execute is the golden rendering for rung four: a Task allowed
// to start.
func TestRender_execute(t *testing.T) {
	next := data.Next{
		Kind:         data.NextExecute,
		Task:         &data.TaskV2{ID: "T030", Title: "Ready task", Status: data.ColumnPlanned},
		GateDecision: &data.GateDecision{Allowed: true, Actor: data.ActorRoleExecutor},
	}
	want := "Task: T030 — Ready task\n" +
		"Implementation: Status planned.\n" +
		"\n" +
		"Ready: Its dependencies are satisfied; it may start.\n" +
		"\n" +
		"Next action: Start Task T030.\n"
	assertRenderEquals(t, next, want)
}

// TestRender_executeAllowedByException proves a completion allowed only by a
// recorded exception is reported as exactly that — its reason and owner —
// and never as clearance.
func TestRender_executeAllowedByException(t *testing.T) {
	next := data.Next{
		Kind: data.NextExecute,
		Task: &data.TaskV2{ID: "T031", Title: "Excepted task", Status: data.ColumnInProgress, Stage: data.StageAudit},
		GateDecision: &data.GateDecision{
			Allowed: true, Actor: data.ActorRoleOwner, AllowedByException: true,
			Exception: &data.Exception{Owner: "alice", Check: "C005", Reason: "accepted known risk"},
		},
	}
	text := renderText(next)
	if !strings.Contains(text, "Completion: Allowed by exception, not by clearance: recorded by owner alice for Check C005 — accepted known risk") {
		t.Fatalf("renderText() = %q, want the exception reported as exception, not clearance", text)
	}
	if strings.Contains(text, "clearance is current") || strings.Contains(text, "Technical clearance:") {
		t.Errorf("renderText() = %q, must not present an exception-allowed completion as clearance", text)
	}
}

// TestRender_checkNeeded is the golden rendering for rung five: missing
// clearance, one of the four non-current states.
func TestRender_checkNeeded(t *testing.T) {
	next := data.Next{
		Kind:      data.NextCheckNeeded,
		Task:      &data.TaskV2{ID: "T040", Title: "Needs check", Status: data.ColumnInProgress, Stage: data.StageAudit},
		Clearance: &data.Clearance{State: data.ClearanceMissing},
	}
	want := "Task: T040 — Needs check\n" +
		"Implementation: Status in_progress, stage audit.\n" +
		"\n" +
		"Technical clearance: No Check has ever been recorded for this target.\n" +
		"\n" +
		"Next action: Owner: request an optional Task Check or record an explicit owner waiver; the Full Objective Check remains mandatory before Objective closure.\n"
	assertRenderEquals(t, next, want)
}

// TestRender_clearanceStatesAreDistinct proves missing, needs_work, stale,
// unknown, and current each render in their own words — plus the sixth
// sentence, for a CLEAR Check no checker session signed — naming the Check
// and the recorded freshness basis when one exists.
func TestRender_clearanceStatesAreDistinct(t *testing.T) {
	assessedAt := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	cases := []struct {
		name      string
		clearance *data.Clearance
		want      string
	}{
		{"missing", &data.Clearance{State: data.ClearanceMissing}, "No Check has ever been recorded for this target."},
		{"needs_work", &data.Clearance{State: data.ClearanceNeedsWork, Check: "C001"}, "Check C001 recorded NEEDS WORK."},
		{"stale", &data.Clearance{State: data.ClearanceStale, Check: "C002", Freshness: &data.Freshness{
			State: data.FreshnessStale, Check: "C002", AssessedBy: data.Actor{Role: data.ActorRoleChecker, Session: "S1"}, AssessedAt: assessedAt, Basis: "diff review",
		}}, "Check C002 is recorded CLEAR, but a freshness assessment marks it stale — clearance is stale. Assessed stale by checker session S1 on 2026-09-01, basis: diff review."},
		{"unknown", &data.Clearance{State: data.ClearanceUnknown, Check: "C003", Freshness: &data.Freshness{
			State: data.FreshnessUnknown, Check: "C003", AssessedBy: data.Actor{Role: data.ActorRoleChecker, Session: "S4"}, AssessedAt: assessedAt, Basis: "not reassessed",
		}}, "Check C003 is recorded CLEAR, but a freshness assessment marks it unknown — clearance is unknown. Assessed unknown by checker session S4 on 2026-09-01, basis: not reassessed."},
		// ResolveClearance also reports a CLEAR Check no checker session signed
		// as unknown. It is a different fact to act on, so it gets its own
		// sentence. The V2 Check decoder refuses that shape, so it is
		// unreachable from a project on disk and proven here instead.
		{"unknown without checker provenance", &data.Clearance{State: data.ClearanceUnknown, Check: "C005"},
			"Check C005 is recorded CLEAR, but no independent checker session signed it — clearance is not independently established."},
		{"current", &data.Clearance{State: data.ClearanceCurrent, Check: "C004", Freshness: &data.Freshness{
			State: data.FreshnessCurrent, Check: "C004", AssessedBy: data.Actor{Role: data.ActorRoleChecker, Session: "S2"}, AssessedAt: assessedAt, Basis: "reran the suite",
		}}, "Check C004 is recorded CLEAR. Assessed current by checker session S2 on 2026-09-01, basis: reran the suite."},
	}

	seen := make(map[string]bool, len(cases))
	for _, c := range cases {
		got := ClearancePhrase(c.clearance)
		if got != c.want {
			t.Errorf("%s: ClearancePhrase() = %q, want %q", c.name, got, c.want)
		}
		if seen[got] {
			t.Errorf("%s: ClearancePhrase() = %q duplicates another state's wording", c.name, got)
		}
		seen[got] = true
	}
}

// TestRender_ownerValidationRequired is the golden rendering for rung six:
// clearance is current, but the owner has not accepted it.
func TestRender_ownerValidationRequired(t *testing.T) {
	next := data.Next{
		Kind: data.NextOwnerValidationRequired,
		Task: &data.TaskV2{ID: "T050", Title: "Needs owner", Status: data.ColumnInProgress, Stage: data.StageAudit},
		Clearance: &data.Clearance{State: data.ClearanceCurrent, Check: "C010", Freshness: &data.Freshness{
			State: data.FreshnessCurrent, Check: "C010", AssessedBy: data.Actor{Role: data.ActorRoleChecker, Session: "S1"},
			AssessedAt: time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC), Basis: "reviewed diff",
		}},
		GateDecision: &data.GateDecision{Blockers: []data.GateBlocker{
			{Kind: data.GateBlockOwnerAcceptance, Detail: "owner has not accepted current check C010"},
		}},
	}
	want := "Task: T050 — Needs owner\n" +
		"Implementation: Status in_progress, stage audit.\n" +
		"\n" +
		"Technical clearance: Check C010 is recorded CLEAR. Assessed current by checker session S1 on 2026-09-01, basis: reviewed diff.\n" +
		"Owner wait: Owner acceptance is required: the owner has not yet accepted Check C010.\n" +
		"\n" +
		"Next action: Ask the owner to accept the current Check.\n"
	assertRenderEquals(t, next, want)
}

// TestRender_ownerWaitDistinctFromTechnicalBlock proves an owner wait always
// carries its own "Owner wait:" line, distinct from the "Technical
// clearance:" line describing the underlying Check.
func TestRender_ownerWaitDistinctFromTechnicalBlock(t *testing.T) {
	next := data.Next{
		Kind:      data.NextOwnerValidationRequired,
		Task:      &data.TaskV2{ID: "T051", Title: "Needs owner too", Status: data.ColumnInProgress, Stage: data.StageAudit},
		Clearance: &data.Clearance{State: data.ClearanceCurrent, Check: "C011"},
		GateDecision: &data.GateDecision{Blockers: []data.GateBlocker{
			{Kind: data.GateBlockOwnerAcceptance},
		}},
	}
	text := renderText(next)
	if !strings.Contains(text, "Technical clearance:") || !strings.Contains(text, "Owner wait:") {
		t.Fatalf("renderText() = %q, want both a Technical clearance line and a distinct Owner wait line", text)
	}
}

// TestRender_objectiveIntegration is the golden rendering for rung seven: an
// Objective whose owned Tasks are all done but whose own integration
// clearance is not current.
func TestRender_objectiveIntegration(t *testing.T) {
	next := data.Next{
		Kind:      data.NextObjectiveIntegration,
		Objective: &data.ObjectiveV2{ID: "O001", Title: "Ship it", Status: data.ColumnInProgress},
		Clearance: &data.Clearance{State: data.ClearanceMissing},
		GateDecision: &data.GateDecision{Blockers: []data.GateBlocker{
			{Kind: data.GateBlockClearanceMissing, Detail: "no recorded check"},
		}},
	}
	want := "Objective: O001 — Ship it\n" +
		"Implementation: Status in_progress.\n" +
		"\n" +
		"Technical clearance: No Check has ever been recorded for this target.\n" +
		"\n" +
		"Next action: Record the Objective O001 integration Check.\n"
	assertRenderEquals(t, next, want)
}

// TestRender_objectiveIntegrationOwnerWait proves the same rung reports an
// owner wait, distinctly, when clearance is current but owner acceptance is
// outstanding.
func TestRender_objectiveIntegrationOwnerWait(t *testing.T) {
	next := data.Next{
		Kind:      data.NextObjectiveIntegration,
		Objective: &data.ObjectiveV2{ID: "O002", Title: "Ship it too", Status: data.ColumnInProgress},
		Clearance: &data.Clearance{State: data.ClearanceCurrent, Check: "C020"},
		GateDecision: &data.GateDecision{Blockers: []data.GateBlocker{
			{Kind: data.GateBlockOwnerAcceptance, Detail: "owner has not accepted current check C020"},
		}},
	}
	text := renderText(next)
	if !strings.Contains(text, "Owner wait: Owner acceptance is required: the owner has not yet accepted Check C020.") {
		t.Fatalf("renderText() = %q, want the objective integration owner wait named", text)
	}
	if !strings.Contains(text, "Next action: Ask the owner to accept the Objective's current integration Check.") {
		t.Fatalf("renderText() = %q, want the owner-specific next action", text)
	}
}

func TestRender_releaseRungsKeepPromiseEvidenceAndActionDistinct(t *testing.T) {
	release := &data.ReleaseV2{
		ID: "R001", Title: "First delivery", Status: data.ColumnInProgress,
		Outcome: "Ship the promised outcome.",
		Evidence: &data.Evidence{OwnerValidation: &data.OwnerValidation{
			AcceptedCheck: "C003", AcceptedBy: data.Actor{Role: data.ActorRoleOwner, Session: "owner-1"},
		}},
	}
	current := &data.Clearance{State: data.ClearanceCurrent, Check: "C003", Freshness: &data.Freshness{
		State: data.FreshnessCurrent, Check: "C003", AssessedBy: data.Actor{Role: data.ActorRoleChecker, Session: "checker-1"},
		AssessedAt: time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC), Basis: "release suite",
	}}

	cases := []struct {
		name string
		next data.Next
		want []string
	}{
		{
			name: "integration",
			next: data.Next{Kind: data.NextReleaseIntegration, Release: release, GateDecision: &data.GateDecision{Blockers: []data.GateBlocker{{Kind: data.GateBlockReleaseObjectiveIncomplete, Objective: "O001", Detail: "member objective is not complete"}}}},
			want: []string{"Goal: R001 — First delivery", "Goal outcome: Ship the promised outcome.", "Goal readiness: Member Objective O001 is not complete", "Complete the member Objectives of Goal R001"},
		},
		{
			name: "check needed",
			next: data.Next{Kind: data.NextReleaseCheckNeeded, Release: release, Clearance: &data.Clearance{State: data.ClearanceStale, Check: "C003", Freshness: current.Freshness}},
			want: []string{"Technical clearance:", "clearance is stale", "Record a fresh Goal Check for R001"},
		},
		{
			name: "owner wait",
			next: data.Next{Kind: data.NextReleaseOwnerValidationRequired, Release: release, Clearance: current, GateDecision: &data.GateDecision{Blockers: []data.GateBlocker{{Kind: data.GateBlockOwnerAcceptance}}}},
			want: []string{"Technical clearance:", "Owner wait: Owner acceptance is required", "Ask the owner to accept the current Goal Check"},
		},
		{
			name: "ready",
			next: data.Next{Kind: data.NextReleaseReady, Release: release, Clearance: current, GateDecision: &data.GateDecision{Allowed: true}},
			want: []string{"Goal readiness: Check C003 is current and accepted by owner session owner-1", "Record Goal R001 as done"},
		},
		{
			name: "done by exception",
			next: data.Next{Kind: data.NextReleaseReady, Release: release, GateDecision: &data.GateDecision{Allowed: true, AllowedByException: true, Exception: &data.Exception{Owner: "owner-1", Check: "C003", Reason: "accepted risk"}}},
			want: []string{"Completion: Allowed by exception, not by clearance", "accepted risk", "Record Goal R001 as done under the recorded exception"},
		},
		{
			name: "historical completion",
			next: data.Next{Kind: data.NextReleaseReady, Release: &data.ReleaseV2{ID: "R001", Title: "First delivery", Status: data.ColumnDone}, GateDecision: &data.GateDecision{Allowed: true, AllowedByLegacyCompletion: true, LegacyCompletion: &data.LegacyCompletionReference{ArchivePath: ".savepoint/archive/release.md"}}},
			want: []string{"Historical completion: Goal R001", "not a new V2 CLEAR Check", "Review the archived historical completion"},
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			text := renderText(test.next)
			for _, want := range test.want {
				if !strings.Contains(text, want) {
					t.Errorf("renderText() = %q, want %q", text, want)
				}
			}
		})
	}
}

func TestRender_releaseSelectionDiagnosticsRequireNewSelection(t *testing.T) {
	cases := []struct {
		name       string
		diagnostic *data.SelectionDiagnostic
		want       []string
	}{
		{"missing", &data.SelectionDiagnostic{Kind: data.SelectionReleaseNotFound, Release: "R999", ID: "R999"}, []string{"Goal R999", "does not exist", "Next action:"}},
		{"archived", &data.SelectionDiagnostic{Kind: data.SelectionReleaseArchived, Release: "R001", ID: "R001"}, []string{"Goal R001", "historical", "Next action:"}},
		{"unassigned", &data.SelectionDiagnostic{Kind: data.SelectionObjectiveUnassigned, Release: "R001", Objective: "O001"}, []string{"Objective O001", "inside Goal R001", "no Goal reference", "Next action:"}},
		{"mismatch", &data.SelectionDiagnostic{Kind: data.SelectionReleaseMismatch, Release: "R001", Objective: "O001", ObjectiveRelease: "R002"}, []string{"Objective O001", "Goal R001", "Goal R002", "Next action:"}},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			text := renderText(data.Next{Kind: data.NextNothingSelected, SelectionDiagnostic: test.diagnostic})
			for _, want := range test.want {
				if !strings.Contains(text, want) {
					t.Errorf("renderText() = %q, want %q", text, want)
				}
			}
		})
	}
}

// TestRender_selectedPlannedTask renders the explicitly selected planned
// Task through the existing start resolver.
func TestRender_selectedPlannedTask(t *testing.T) {
	next := data.Next{
		Kind:         data.NextExecute,
		Objective:    &data.ObjectiveV2{ID: "O060", Title: "Selected Objective", Status: data.ColumnInProgress},
		Task:         &data.TaskV2{ID: "T060", Title: "Selected Task", Objective: "O060", Status: data.ColumnPlanned},
		GateDecision: &data.GateDecision{Allowed: true, Actor: data.ActorRoleExecutor},
	}
	want := "Objective: O060 — Selected Objective\n" +
		"Task: T060 — Selected Task\n" +
		"Implementation: Status planned.\n" +
		"\n" +
		"Ready: Its dependencies are satisfied; it may start.\n" +
		"\n" +
		"Next action: Start Task T060.\n"
	assertRenderEquals(t, next, want)
}

// TestRender_planSelectedObjective proves a Task-less selected Objective
// directs planning under that Objective.
func TestRender_planSelectedObjective(t *testing.T) {
	next := data.Next{
		Kind:      data.NextPlanObjective,
		Objective: &data.ObjectiveV2{ID: "O003", Title: "New objective", Status: data.ColumnPlanned},
	}
	want := "Objective: O003 — New objective\n" +
		"Implementation: Status planned.\n" +
		"\n" +
		"Next action: Plan Tasks under Objective O003.\n"
	assertRenderEquals(t, next, want)
}

func TestActionPhraseSelectionGuidance(t *testing.T) {
	cases := []struct {
		name string
		next data.Next
		want string
	}{
		{
			name: "plan under selected Objective",
			next: data.Next{Kind: data.NextPlanObjective, Objective: &data.ObjectiveV2{ID: "O-003"}},
			want: "Plan Tasks under Objective O-003.",
		},
		{
			name: "select a Task under selected Objective",
			next: data.Next{Kind: data.NextSelectTask, Objective: &data.ObjectiveV2{ID: "O-003"}},
			want: "Select a Task under Objective O-003: press p on the board, or ask the agent to \"set router to O-### T-###\".",
		},
		{
			name: "select an Objective when none is selected",
			next: data.Next{Kind: data.NextNothingSelected},
			want: "Select an Objective: press p on the board, or ask the agent to \"set router to O-### T-###\".",
		},
		{
			name: "choose optional Task Check or owner waiver",
			next: data.Next{Kind: data.NextCheckNeeded, Task: &data.TaskV2{ID: "T-004", Status: data.ColumnInProgress, Stage: data.StageAudit}},
			want: "Owner: request an optional Task Check or record an explicit owner waiver; the Full Objective Check remains mandatory before Objective closure.",
		},
		{
			name: "work on selected Issue",
			next: data.Next{Kind: data.NextIssue, Issue: &data.IssueV2{ID: "I-042", Status: data.IssueStatusOpen}},
			want: "Work on Issue I-042.",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := ActionPhrase(tc.next); got != tc.want {
				t.Errorf("ActionPhrase() = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestRender_nothingSelected gives an explicit way to choose the next
// Objective when the router has no current selection.
func TestRender_nothingSelected(t *testing.T) {
	next := data.Next{Kind: data.NextNothingSelected}
	want := "Next action: Select an Objective: press p on the board, or ask the agent to \"set router to O-### T-###\".\n"
	assertRenderEquals(t, next, want)
}

// TestRender_selectionDiagnosticAndNextActionTogether proves a selection
// diagnostic renders alongside the next action the project's own records
// still support, rather than in place of it.
func TestRender_selectionDiagnosticAndNextActionTogether(t *testing.T) {
	next := data.Next{
		Kind:                data.NextNothingSelected,
		SelectionDiagnostic: &data.SelectionDiagnostic{Kind: data.SelectionNotFound, RecordKind: data.SelectionRecordTask, ID: "T999"},
	}
	want := "\nSelection: The router names task T999, which does not exist among the project's live records.\n" +
		"\n" +
		"Next action: Select an Objective: press p on the board, or ask the agent to \"set router to O-### T-###\".\n"
	assertRenderEquals(t, next, want)
}

func TestSelectionPhrase_doneSelectionsNameFinishedRecord(t *testing.T) {
	cases := []struct {
		name       string
		recordKind data.SelectionRecordKind
		id         string
		want       string
	}{
		{"Task", data.SelectionRecordTask, "T-001", "Warning: router still selects finished Task T-001."},
		{"Objective", data.SelectionRecordObjective, "O-001", "Warning: router still selects finished Objective O-001."},
		{"Issue", data.SelectionRecordIssue, "I-042", "Warning: router still selects resolved Issue I-042."},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := SelectionPhrase(&data.SelectionDiagnostic{
				Kind: data.SelectionDone, RecordKind: tc.recordKind, ID: tc.id,
			})
			if got != tc.want {
				t.Errorf("SelectionPhrase() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestRender_selectedIssueIsContextWhenTaskWinsNext(t *testing.T) {
	issue := &data.IssueV2{ID: "I-042", Title: "Repair the parser", Status: data.IssueStatusInProgress}
	next := data.Next{
		Kind:      data.NextExecute,
		Objective: &data.ObjectiveV2{ID: "O-014", Title: "Router Issue target", Status: data.ColumnInProgress},
		Task:      &data.TaskV2{ID: "T-028", Title: "Copy the line", Objective: "O-014", Status: data.ColumnInProgress, Stage: data.StageBuild},
		Issue:     issue,
	}

	if got, want := NextLine(next), "Build T-028 — Copy the line (O-014)"; got != want {
		t.Fatalf("NextLine() = %q, want Task line %q", got, want)
	}
	if got, want := IssueContextLine(issue), "Issue: Fix I-042 — Repair the parser"; got != want {
		t.Fatalf("IssueContextLine() = %q, want %q", got, want)
	}
	if text := renderText(next); !strings.Contains(text, IssueContextLine(issue)) {
		t.Errorf("renderText() = %q, want selected Issue context %q", text, IssueContextLine(issue))
	}
}

func TestRender_doneSelectionWarningAlongsideCurrentNext(t *testing.T) {
	diagnostic := &data.SelectionDiagnostic{
		Kind: data.SelectionDone, RecordKind: data.SelectionRecordTask, ID: "T-001",
	}
	next := data.Next{
		Kind:                data.NextSelectTask,
		Objective:           &data.ObjectiveV2{ID: "O-001", Title: "Ship it", Status: data.ColumnInProgress},
		SelectionDiagnostic: diagnostic,
	}

	text := renderText(next)
	if want := "Selection: " + SelectionPhrase(diagnostic); !strings.Contains(text, want) {
		t.Errorf("renderText() = %q, want stale-selection warning %q", text, want)
	}
	if !strings.Contains(text, "Objective: O-001 — Ship it") || !strings.Contains(text, "Next action: ") {
		t.Errorf("renderText() = %q, want the selected Objective and its existing Next action alongside the warning", text)
	}
}

// TestRender_selectionMismatch proves the mismatch diagnostic names both IDs
// and treats the Task record's ownership as authoritative in its wording.
func TestRender_selectionMismatch(t *testing.T) {
	next := data.Next{
		Kind: data.NextNothingSelected,
		SelectionDiagnostic: &data.SelectionDiagnostic{
			Kind: data.SelectionMismatch, RouterObjective: "O001", Task: "T001", TaskObjective: "O002",
		},
	}
	text := renderText(next)
	want := "The router names Objective O001 and Task T001, but Task T001's own record names O002 as its owner — the Task record wins, so this selection is not honored."
	if !strings.Contains(text, want) {
		t.Fatalf("renderText() = %q, want it to contain %q", text, want)
	}
}

// TestRender_issuesSection proves relevant Issues render with their ID,
// type, and status.
func TestRender_issuesSection(t *testing.T) {
	next := data.Next{
		Kind:      data.NextCheckNeeded,
		Task:      &data.TaskV2{ID: "T070", Title: "Has issues", Status: data.ColumnInProgress, Stage: data.StageAudit},
		Clearance: &data.Clearance{State: data.ClearanceNeedsWork, Check: "C020"},
		Issues: []*data.IssueV2{
			{ID: "I001", Title: "Found a bug", Type: data.IssueTypeDefect, Status: data.IssueStatusOpen},
			{ID: "I002", Title: "Drifted from design", Type: data.IssueTypeDrift, Status: data.IssueStatusInProgress},
		},
	}
	text := renderText(next)
	want := "Issues:\n- I001 (defect, open): Found a bug\n- I002 (drift, in_progress): Drifted from design\n"
	if !strings.Contains(text, want) {
		t.Fatalf("renderText() = %q, want it to contain %q", text, want)
	}
}

// TestRender_noIssuesRendersNoSection proves a Next with no Issues renders
// no "Issues:" section at all, rather than an empty one.
func TestRender_noIssuesRendersNoSection(t *testing.T) {
	next := data.Next{Kind: data.NextNothingSelected}
	text := renderText(next)
	if strings.Contains(text, "Issues:") {
		t.Fatalf("renderText() = %q, want no Issues section when Issues is nil", text)
	}
}

// TestRender_noVerificationClaims scans every rung's rendering — including
// exception, owner-wait, and every clearance state — for language claiming
// resume verified, ran, checked, or confirmed anything. Resume read
// records; it ran nothing.
func TestRender_noVerificationClaims(t *testing.T) {
	forbidden := []string{"verified", "verify", "confirmed", "confirm", "checked ", " ran ", "resume ran"}

	for _, next := range allRungFixtures() {
		text := strings.ToLower(renderText(next))
		for _, word := range forbidden {
			if strings.Contains(text, word) {
				t.Errorf("renderText(%s) = %q, contains forbidden verification claim %q", next.Kind, text, word)
			}
		}
	}
}

// TestRender_determinism proves two renders of the same projection are
// byte-identical.
func TestRender_determinism(t *testing.T) {
	for _, next := range allRungFixtures() {
		first := renderText(next)
		second := renderText(next)
		if first != second {
			t.Errorf("renderText(%s) is not deterministic: %q != %q", next.Kind, first, second)
		}
	}
}

// TestRender_noANSIEscapes proves output never carries a terminal escape
// sequence: resume produces plain text, readable in a pipe, a log, or an
// agent transcript exactly as on a terminal.
func TestRender_noANSIEscapes(t *testing.T) {
	for _, next := range allRungFixtures() {
		if text := renderText(next); strings.ContainsRune(text, '\x1b') {
			t.Errorf("renderText(%s) contains an ANSI escape", next.Kind)
		}
	}
}

// TestRender_narrowWidthReadable proves every rendered line is a short,
// unpadded sentence rather than a fixed-width table column, so the text
// wraps and stays readable at both 80 and 40 columns without truncation or
// misaligned columns.
func TestRender_narrowWidthReadable(t *testing.T) {
	boxDrawing := []string{"─", "│", "┌", "\t"}
	for _, next := range allRungFixtures() {
		text := renderText(next)
		for _, glyph := range boxDrawing {
			if strings.Contains(text, glyph) {
				t.Errorf("renderText(%s) contains fixed-width layout glyph %q, which does not stay readable at 40 columns", next.Kind, glyph)
			}
		}
	}
}

// TestRender_propagatesWriterError proves Render returns a failing writer's
// error rather than swallowing it.
func TestRender_propagatesWriterError(t *testing.T) {
	wantErr := errors.New("disk full")
	next := data.Next{Kind: data.NextNothingSelected}

	err := Render(failingWriter{err: wantErr}, next)
	if !errors.Is(err, wantErr) {
		t.Fatalf("Render() error = %v, want %v", err, wantErr)
	}
}

// TestPackage_noFilesystemNetworkOrSubprocessImports proves resume's own
// source imports nothing that performs filesystem, subprocess, or network
// access — Render's only IO is the io.Writer it is handed.
func TestPackage_noFilesystemNetworkOrSubprocessImports(t *testing.T) {
	pkg, err := build.ImportDir(".", build.IgnoreVendor)
	if err != nil {
		t.Fatalf("scan package imports: %v", err)
	}
	forbidden := map[string]bool{"os": true, "os/exec": true, "net": true, "net/http": true, "syscall": true}
	for _, imported := range pkg.Imports {
		if forbidden[imported] {
			t.Errorf("internal/resume imports %q, which performs IO beyond the handed io.Writer", imported)
		}
	}
}

type failingWriter struct{ err error }

func (f failingWriter) Write([]byte) (int, error) { return 0, f.err }

func assertRenderEquals(t *testing.T, next data.Next, want string) {
	t.Helper()
	var buf bytes.Buffer
	if err := Render(&buf, next); err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	want = NextLine(next) + "\n" + want
	if got := buf.String(); got != want {
		t.Fatalf("Render() = %q, want %q", got, want)
	}
}

// allRungFixtures returns one representative Next per reachable rung, plus
// the exception, owner-wait, and Issues variants, for the cross-cutting
// tests above that must hold across every rendering.
func allRungFixtures() []data.Next {
	assessedAt := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	return []data.Next{
		{
			Kind: data.NextReplan,
			Task: &data.TaskV2{ID: "T010", Title: "Some task", Status: data.ColumnInProgress, Stage: data.StageBuild},
			GateDecision: &data.GateDecision{Blockers: []data.GateBlocker{
				{Kind: data.GateBlockReplan, Detail: "scope changed"},
			}},
		},
		{
			Kind: data.NextDependency,
			Task: &data.TaskV2{ID: "T020", Title: "Blocked task", Status: data.ColumnPlanned},
			GateDecision: &data.GateDecision{Blockers: []data.GateBlocker{
				{Kind: data.GateBlockDependency, Dependency: &data.DependencyBlock{Target: "T019", Kind: data.DependencyBlockNotDone}},
			}},
		},
		{
			Kind:         data.NextExecute,
			Task:         &data.TaskV2{ID: "T030", Title: "Ready task", Status: data.ColumnPlanned},
			GateDecision: &data.GateDecision{Allowed: true, Actor: data.ActorRoleExecutor},
		},
		{
			Kind: data.NextExecute,
			Task: &data.TaskV2{ID: "T031", Title: "Excepted task", Status: data.ColumnInProgress, Stage: data.StageAudit},
			GateDecision: &data.GateDecision{
				Allowed: true, Actor: data.ActorRoleOwner, AllowedByException: true,
				Exception: &data.Exception{Owner: "alice", Check: "C005", Reason: "accepted known risk", RecordedAt: assessedAt},
			},
		},
		{
			Kind:      data.NextCheckNeeded,
			Task:      &data.TaskV2{ID: "T040", Title: "Needs check", Status: data.ColumnInProgress, Stage: data.StageAudit},
			Clearance: &data.Clearance{State: data.ClearanceMissing},
		},
		{
			Kind:      data.NextCheckNeeded,
			Task:      &data.TaskV2{ID: "T041", Title: "Needs rework", Status: data.ColumnInProgress, Stage: data.StageAudit},
			Clearance: &data.Clearance{State: data.ClearanceNeedsWork, Check: "C006"},
		},
		{
			Kind: data.NextCheckNeeded,
			Task: &data.TaskV2{ID: "T042", Title: "Stale check", Status: data.ColumnInProgress, Stage: data.StageAudit},
			Clearance: &data.Clearance{State: data.ClearanceStale, Check: "C007", Freshness: &data.Freshness{
				State: data.FreshnessStale, Check: "C007", AssessedBy: data.Actor{Role: data.ActorRoleChecker, Session: "S3"}, AssessedAt: assessedAt, Basis: "diff review",
			}},
		},
		{
			Kind:      data.NextCheckNeeded,
			Task:      &data.TaskV2{ID: "T043", Title: "Unknown freshness", Status: data.ColumnInProgress, Stage: data.StageAudit},
			Clearance: &data.Clearance{State: data.ClearanceUnknown, Check: "C008"},
		},
		{
			Kind: data.NextOwnerValidationRequired,
			Task: &data.TaskV2{ID: "T050", Title: "Needs owner", Status: data.ColumnInProgress, Stage: data.StageAudit},
			Clearance: &data.Clearance{State: data.ClearanceCurrent, Check: "C010", Freshness: &data.Freshness{
				State: data.FreshnessCurrent, Check: "C010", AssessedBy: data.Actor{Role: data.ActorRoleChecker, Session: "S1"}, AssessedAt: assessedAt, Basis: "reviewed diff",
			}},
			GateDecision: &data.GateDecision{Blockers: []data.GateBlocker{
				{Kind: data.GateBlockOwnerAcceptance, Detail: "owner has not accepted current check C010"},
			}},
		},
		{
			Kind:      data.NextObjectiveIntegration,
			Objective: &data.ObjectiveV2{ID: "O001", Title: "Ship it", Status: data.ColumnInProgress},
			Clearance: &data.Clearance{State: data.ClearanceMissing},
			GateDecision: &data.GateDecision{Blockers: []data.GateBlocker{
				{Kind: data.GateBlockClearanceMissing, Detail: "no recorded check"},
			}},
		},
		{
			Kind:      data.NextObjectiveIntegration,
			Objective: &data.ObjectiveV2{ID: "O002", Title: "Ship it too", Status: data.ColumnInProgress},
			Clearance: &data.Clearance{State: data.ClearanceCurrent, Check: "C020"},
			GateDecision: &data.GateDecision{Blockers: []data.GateBlocker{
				{Kind: data.GateBlockOwnerAcceptance, Detail: "owner has not accepted current check C020"},
			}},
		},
		{
			Kind:         data.NextObjectiveReady,
			Objective:    &data.ObjectiveV2{ID: "O003", Title: "Ready to close", Status: data.ColumnInProgress},
			Clearance:    &data.Clearance{State: data.ClearanceCurrent, Check: "C021"},
			GateDecision: &data.GateDecision{Allowed: true, Actor: data.ActorRoleChecker},
		},
		{
			Kind:      data.NextSelectTask,
			Objective: &data.ObjectiveV2{ID: "O004", Title: "Choose a Task", Status: data.ColumnInProgress},
		},
		{Kind: data.NextNothingSelected},
		{
			Kind:                data.NextNothingSelected,
			SelectionDiagnostic: &data.SelectionDiagnostic{Kind: data.SelectionNotFound, RecordKind: data.SelectionRecordTask, ID: "T999"},
		},
		{
			Kind:      data.NextCheckNeeded,
			Task:      &data.TaskV2{ID: "T070", Title: "Has issues", Status: data.ColumnInProgress, Stage: data.StageAudit},
			Clearance: &data.Clearance{State: data.ClearanceNeedsWork, Check: "C020"},
			Issues: []*data.IssueV2{
				{ID: "I001", Title: "Found a bug", Type: data.IssueTypeDefect, Status: data.IssueStatusOpen},
				{ID: "I002", Title: "Drifted from design", Type: data.IssueTypeDrift, Status: data.IssueStatusInProgress},
			},
		},
	}
}
