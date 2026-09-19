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

// TestRender_pendingMigration is the golden rendering for rung one: pending
// migration outranks every other rung and carries no Task or Objective.
func TestRender_pendingMigration(t *testing.T) {
	next := data.Next{Kind: data.NextPendingMigration, Migration: data.MigrationState{Pending: true, OperationID: "OP001"}}
	want := "Migration: an operation is in progress (OP001). Every other action is on hold until it resolves.\n" +
		"\n" +
		"Next action: Wait for the pending migration to finish; nothing else is actionable until it resolves.\n"
	assertRenderEquals(t, next, want)
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
		"Next action: Record a fresh Check against this target.\n"
	assertRenderEquals(t, next, want)
}

// TestRender_clearanceStatesAreDistinct proves missing, needs_work, stale,
// unknown, and current each render in their own words — plus the sixth
// sentence, for the CLEAR Check whose current assessment lacks independent
// checker provenance — naming the Check and the recorded freshness basis when
// one exists.
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
		}}, "Check C002 is recorded CLEAR, but its freshness assessment does not name it current — clearance is stale. Assessed stale by checker session S1 on 2026-09-01, basis: diff review."},
		{"unknown", &data.Clearance{State: data.ClearanceUnknown, Check: "C003"}, "Check C003 is recorded CLEAR, but no freshness assessment has ever been recorded for it — clearance is unknown."},
		// ResolveClearance also reports a CLEAR Check whose current assessment
		// carries no independent checker session as unknown, but with that
		// assessment attached. It is a different fact to act on, so it gets its
		// own sentence. Both V2 decoders refuse the shapes that produce it, so
		// it is unreachable from a project on disk and proven here instead.
		{"unknown without checker provenance", &data.Clearance{State: data.ClearanceUnknown, Check: "C005", Freshness: &data.Freshness{
			State: data.FreshnessCurrent, Check: "C005", AssessedBy: data.Actor{Role: data.ActorRoleExecutor, Session: "S3"}, AssessedAt: assessedAt, Basis: "self-reported",
		}}, "Check C005 is recorded CLEAR and its freshness assessment names it current, but that evidence carries no independent checker session — clearance is not independently established. Assessed current by executor session S3 on 2026-09-01, basis: self-reported."},
		{"current", &data.Clearance{State: data.ClearanceCurrent, Check: "C004", Freshness: &data.Freshness{
			State: data.FreshnessCurrent, Check: "C004", AssessedBy: data.Actor{Role: data.ActorRoleChecker, Session: "S2"}, AssessedAt: assessedAt, Basis: "reran the suite",
		}}, "Check C004 is recorded CLEAR and its freshness assessment names it current. Assessed current by checker session S2 on 2026-09-01, basis: reran the suite."},
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
		"Technical clearance: Check C010 is recorded CLEAR and its freshness assessment names it current. Assessed current by checker session S1 on 2026-09-01, basis: reviewed diff.\n" +
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

// TestRender_readyTask is the golden rendering for rung eight when the ready
// record is a Task.
func TestRender_readyTask(t *testing.T) {
	next := data.Next{
		Kind:         data.NextReady,
		Task:         &data.TaskV2{ID: "T060", Title: "Ready elsewhere", Status: data.ColumnPlanned},
		GateDecision: &data.GateDecision{Allowed: true, Actor: data.ActorRoleExecutor},
	}
	want := "Task: T060 — Ready elsewhere\n" +
		"Implementation: Status planned.\n" +
		"\n" +
		"Next action: Start Task T060.\n"
	assertRenderEquals(t, next, want)
}

// TestRender_readyObjective proves rung eight's Objective-only fallback:
// ready to plan Tasks under an Objective with none yet.
func TestRender_readyObjective(t *testing.T) {
	next := data.Next{
		Kind:      data.NextReady,
		Objective: &data.ObjectiveV2{ID: "O003", Title: "New objective", Status: data.ColumnPlanned},
	}
	want := "Objective: O003 — New objective\n" +
		"Implementation: Status planned.\n" +
		"\n" +
		"Next action: Plan Tasks under Objective O003.\n"
	assertRenderEquals(t, next, want)
}

// TestRender_planObjective is the golden rendering for rung nine: nothing
// selected, nothing ready, plan the next Objective. A fresh project with no
// Objectives lands here with the same wording.
func TestRender_planObjective(t *testing.T) {
	next := data.Next{Kind: data.NextPlanObjective}
	want := "Next action: Plan the next Objective — for a project with nothing underway yet, start with Idea/Design.\n"
	assertRenderEquals(t, next, want)
}

// TestRender_selectionDiagnosticAndNextActionTogether proves a selection
// diagnostic renders alongside the next action the project's own records
// still support, rather than in place of it.
func TestRender_selectionDiagnosticAndNextActionTogether(t *testing.T) {
	next := data.Next{
		Kind:                data.NextPlanObjective,
		SelectionDiagnostic: &data.SelectionDiagnostic{Kind: data.SelectionNotFound, RecordKind: data.SelectionRecordTask, ID: "T999"},
	}
	want := "Selection: The router names task T999, which does not exist among the project's live records.\n" +
		"\n" +
		"Next action: Plan the next Objective — for a project with nothing underway yet, start with Idea/Design.\n"
	assertRenderEquals(t, next, want)
}

// TestRender_selectionMismatch proves the mismatch diagnostic names both IDs
// and treats the Task record's ownership as authoritative in its wording.
func TestRender_selectionMismatch(t *testing.T) {
	next := data.Next{
		Kind: data.NextPlanObjective,
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
	next := data.Next{Kind: data.NextPlanObjective}
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
	next := data.Next{Kind: data.NextPlanObjective}

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
		{Kind: data.NextPendingMigration, Migration: data.MigrationState{Pending: true, OperationID: "OP001"}},
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
			Kind:         data.NextReady,
			Task:         &data.TaskV2{ID: "T060", Title: "Ready elsewhere", Status: data.ColumnPlanned},
			GateDecision: &data.GateDecision{Allowed: true, Actor: data.ActorRoleExecutor},
		},
		{
			Kind:      data.NextReady,
			Objective: &data.ObjectiveV2{ID: "O003", Title: "New objective", Status: data.ColumnPlanned},
		},
		{Kind: data.NextPlanObjective},
		{
			Kind:                data.NextPlanObjective,
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
