package data

import (
	"go/build"
	"os"
	"strings"
	"testing"
)

// TestResolveSelection_exactMatch proves a router naming both an Objective
// and a Task that resolve, and agree on ownership, returns both records and
// no diagnostic.
func TestResolveSelection_exactMatch(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O001"] = &ObjectiveV2{ID: "O001", Title: "Ship it"}
	index.Tasks["T001"] = &TaskV2{ID: "T001", Title: "Write the code", Objective: "O001"}
	router := &RouterStateV2{State: RouterPhaseTask, Objective: "O001", Task: "T001"}

	selection, diagnostic := ResolveSelection(index, router)
	if diagnostic != nil {
		t.Fatalf("ResolveSelection() diagnostic = %+v, want nil", diagnostic)
	}
	if selection.Objective == nil || selection.Objective.ID != "O001" {
		t.Fatalf("Selection.Objective = %+v, want O001", selection.Objective)
	}
	if selection.Task == nil || selection.Task.ID != "T001" {
		t.Fatalf("Selection.Task = %+v, want T001", selection.Task)
	}
}

// TestResolveSelection_nearMissIDNeverSubstituted proves a project holding
// both T014 and T140 resolves a router naming T014 to exactly T014: no
// numeric-proximity or prefix fallback ever selects the other record.
func TestResolveSelection_nearMissIDNeverSubstituted(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O001"] = &ObjectiveV2{ID: "O001", Title: "Ship it"}
	index.Tasks["T014"] = &TaskV2{ID: "T014", Title: "The real target", Objective: "O001"}
	index.Tasks["T140"] = &TaskV2{ID: "T140", Title: "A similarly numbered task", Objective: "O001"}
	router := &RouterStateV2{State: RouterPhaseTask, Objective: "O001", Task: "T014"}

	selection, diagnostic := ResolveSelection(index, router)
	if diagnostic != nil {
		t.Fatalf("ResolveSelection() diagnostic = %+v, want nil", diagnostic)
	}
	if selection.Task == nil || selection.Task.ID != "T014" {
		t.Fatalf("Selection.Task = %+v, want exactly T014, never T140", selection.Task)
	}
}

// TestResolveSelection_absentTask proves a router-named Task ID missing from
// the live index returns a typed not-found diagnostic naming it, rather than
// a partial selection.
func TestResolveSelection_absentTask(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O001"] = &ObjectiveV2{ID: "O001", Title: "Ship it"}
	router := &RouterStateV2{State: RouterPhaseTask, Objective: "O001", Task: "T999"}

	selection, diagnostic := ResolveSelection(index, router)
	if diagnostic == nil {
		t.Fatal("ResolveSelection() diagnostic = nil, want SelectionNotFound for the absent task")
	}
	if diagnostic.Kind != SelectionNotFound {
		t.Errorf("diagnostic.Kind = %q, want SelectionNotFound", diagnostic.Kind)
	}
	if diagnostic.RecordKind != SelectionRecordTask || diagnostic.ID != "T999" {
		t.Errorf("diagnostic = %+v, want RecordKind task, ID T999", diagnostic)
	}
	if selection.Objective != nil || selection.Task != nil {
		t.Errorf("Selection = %+v, want zero value alongside a diagnostic", selection)
	}
}

// TestResolveSelection_absentObjective proves the same for a router-named
// Objective ID missing from the live index.
func TestResolveSelection_absentObjective(t *testing.T) {
	index := newV2TestIndex()
	router := &RouterStateV2{State: RouterPhaseDesign, Objective: "O999"}

	selection, diagnostic := ResolveSelection(index, router)
	if diagnostic == nil {
		t.Fatal("ResolveSelection() diagnostic = nil, want SelectionNotFound for the absent objective")
	}
	if diagnostic.Kind != SelectionNotFound {
		t.Errorf("diagnostic.Kind = %q, want SelectionNotFound", diagnostic.Kind)
	}
	if diagnostic.RecordKind != SelectionRecordObjective || diagnostic.ID != "O999" {
		t.Errorf("diagnostic = %+v, want RecordKind objective, ID O999", diagnostic)
	}
	if selection.Objective != nil || selection.Task != nil {
		t.Errorf("Selection = %+v, want zero value alongside a diagnostic", selection)
	}
}

// TestResolveSelection_mismatch proves a Task that resolves but is owned by
// a different Objective than the router names returns a typed mismatch
// diagnostic naming both, and does not substitute the Task's own Objective
// for the router's.
func TestResolveSelection_mismatch(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O001"] = &ObjectiveV2{ID: "O001", Title: "Router's guess"}
	index.Objectives["O002"] = &ObjectiveV2{ID: "O002", Title: "Task's real owner"}
	index.Tasks["T001"] = &TaskV2{ID: "T001", Title: "Moved task", Objective: "O002"}
	router := &RouterStateV2{State: RouterPhaseTask, Objective: "O001", Task: "T001"}

	selection, diagnostic := ResolveSelection(index, router)
	if diagnostic == nil {
		t.Fatal("ResolveSelection() diagnostic = nil, want SelectionMismatch")
	}
	if diagnostic.Kind != SelectionMismatch {
		t.Errorf("diagnostic.Kind = %q, want SelectionMismatch", diagnostic.Kind)
	}
	if diagnostic.RouterObjective != "O001" || diagnostic.Task != "T001" || diagnostic.TaskObjective != "O002" {
		t.Errorf("diagnostic = %+v, want RouterObjective O001, Task T001, TaskObjective O002", diagnostic)
	}
	if selection.Objective != nil || selection.Task != nil {
		t.Errorf("Selection = %+v, want zero value alongside a diagnostic", selection)
	}
}

// TestResolveSelection_objectiveOnly proves a router naming an Objective
// with no Task resolves cleanly to that Objective alone, and is not a
// diagnostic.
func TestResolveSelection_objectiveOnly(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O001"] = &ObjectiveV2{ID: "O001", Title: "Plan me"}
	router := &RouterStateV2{State: RouterPhaseDesign, Objective: "O001"}

	selection, diagnostic := ResolveSelection(index, router)
	if diagnostic != nil {
		t.Fatalf("ResolveSelection() diagnostic = %+v, want nil", diagnostic)
	}
	if selection.Objective == nil || selection.Objective.ID != "O001" {
		t.Fatalf("Selection.Objective = %+v, want O001", selection.Objective)
	}
	if selection.Task != nil {
		t.Errorf("Selection.Task = %+v, want nil", selection.Task)
	}
}

// TestResolveSelection_noSelection proves a router in idea or design with no
// Objective and no Task resolves cleanly to no selection, not a diagnostic.
func TestResolveSelection_noSelection(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O001"] = &ObjectiveV2{ID: "O001", Title: "Unrelated"}
	router := &RouterStateV2{State: RouterPhaseIdea}

	selection, diagnostic := ResolveSelection(index, router)
	if diagnostic != nil {
		t.Fatalf("ResolveSelection() diagnostic = %+v, want nil", diagnostic)
	}
	if selection.Objective != nil || selection.Task != nil {
		t.Errorf("Selection = %+v, want zero value for no selection", selection)
	}
}

// TestResolveSelection_readsNoFilesystem builds its V2Index directly in
// memory, with no discovery or file read anywhere in the call, proving
// ResolveSelection consults only the two values it is handed.
func TestResolveSelection_readsNoFilesystem(t *testing.T) {
	index := &V2Index{
		Objectives: map[string]*ObjectiveV2{
			"O001": {ID: "O001", Title: "In-memory only"},
		},
		Tasks: map[string]*TaskV2{
			"T001": {ID: "T001", Title: "In-memory only", Objective: "O001"},
		},
	}
	router := &RouterStateV2{State: RouterPhaseTask, Objective: "O001", Task: "T001"}

	selection, diagnostic := ResolveSelection(index, router)
	if diagnostic != nil {
		t.Fatalf("ResolveSelection() diagnostic = %+v, want nil", diagnostic)
	}
	if selection.Objective == nil || selection.Task == nil {
		t.Fatalf("Selection = %+v, want both records resolved from the in-memory index", selection)
	}
}

// TestResolveNext_pendingMigrationOutranksEverything proves rung one wins
// even over a project that would otherwise land on NextExecute.
func TestResolveNext_pendingMigrationOutranksEverything(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O001"] = &ObjectiveV2{ID: "O001", Status: ColumnPlanned}
	index.Tasks["T001"] = &TaskV2{ID: "T001", Objective: "O001", Status: ColumnPlanned}
	index.ObjectiveTasks["O001"] = []string{"T001"}
	router := &RouterStateV2{State: RouterPhaseIdea}

	next := ResolveNext(NextInput{Index: index, Router: router, Migration: MigrationState{Pending: true, OperationID: "op-1"}})
	if next.Kind != NextPendingMigration {
		t.Fatalf("Kind = %q, want pending_migration even though a ready task exists", next.Kind)
	}
	if next.Migration.OperationID != "op-1" {
		t.Errorf("Migration = %+v, want OperationID op-1 carried through", next.Migration)
	}
}

// TestResolveNext_replanOutranksExecution proves a replan flag reports
// NextReplan rather than NextExecute, even though the Task has no
// dependency and would otherwise be ready to start.
func TestResolveNext_replanOutranksExecution(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O001"] = &ObjectiveV2{ID: "O001", Status: ColumnPlanned}
	index.Tasks["T001"] = &TaskV2{
		ID: "T001", Objective: "O001", Status: ColumnPlanned,
		Evidence: &Evidence{Replan: &Replan{Reason: "scope changed", RecordedBy: Actor{Role: ActorRoleOwner, Session: "sess-1"}}},
	}
	index.ObjectiveTasks["O001"] = []string{"T001"}
	router := &RouterStateV2{State: RouterPhaseTask, Objective: "O001", Task: "T001"}

	next := ResolveNext(NextInput{Index: index, Router: router})
	if next.Kind != NextReplan {
		t.Fatalf("Kind = %q, want replan", next.Kind)
	}
}

// TestResolveNext_replanOutranksCheckNeeded proves the same at the audit
// stage, where the Task would otherwise need a fresh Check.
func TestResolveNext_replanOutranksCheckNeeded(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O001"] = &ObjectiveV2{ID: "O001", Status: ColumnPlanned}
	index.Tasks["T001"] = &TaskV2{
		ID: "T001", Objective: "O001", Status: ColumnInProgress, Stage: StageAudit,
		Evidence: &Evidence{Replan: &Replan{Reason: "scope changed", RecordedBy: Actor{Role: ActorRoleOwner, Session: "sess-1"}}},
	}
	index.ObjectiveTasks["O001"] = []string{"T001"}
	router := &RouterStateV2{State: RouterPhaseCheck, Objective: "O001", Task: "T001"}

	next := ResolveNext(NextInput{Index: index, Router: router})
	if next.Kind != NextReplan {
		t.Fatalf("Kind = %q, want replan (no check is recorded at all, which would otherwise be check_needed)", next.Kind)
	}
}

// TestResolveNext_taskDependencyBlocks proves an unsatisfied Task
// dependency reaches NextDependency carrying ResolveTaskStart's typed
// DependencyBlock.
func TestResolveNext_taskDependencyBlocks(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O001"] = &ObjectiveV2{ID: "O001", Status: ColumnPlanned}
	index.Tasks["T002"] = &TaskV2{ID: "T002", Objective: "O001", Status: ColumnPlanned}
	index.Tasks["T001"] = &TaskV2{
		ID: "T001", Objective: "O001", Status: ColumnPlanned,
		DependsOn: []TaskDependencyV2{{Task: "T002", Requires: TaskDependencyClear}},
	}
	index.ObjectiveTasks["O001"] = []string{"T001", "T002"}
	router := &RouterStateV2{State: RouterPhaseTask, Objective: "O001", Task: "T001"}

	next := ResolveNext(NextInput{Index: index, Router: router})
	if next.Kind != NextDependency {
		t.Fatalf("Kind = %q, want dependency", next.Kind)
	}
	if next.GateDecision == nil || len(next.GateDecision.Blockers) == 0 || next.GateDecision.Blockers[0].Kind != GateBlockDependency {
		t.Fatalf("GateDecision = %+v, want a GateBlockDependency blocker", next.GateDecision)
	}
	if dep := next.GateDecision.Blockers[0].Dependency; dep == nil || dep.Target != "T002" {
		t.Errorf("Dependency block = %+v, want target T002", dep)
	}
}

// TestResolveNext_objectiveDependencyBlocks proves an unsatisfied Objective
// dependency of the Task's owning Objective is distinguishable from a Task
// dependency: same NextDependency rung, but carrying
// ResolveObjectiveDependency's typed ObjectiveDependencyBlock instead.
func TestResolveNext_objectiveDependencyBlocks(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O002"] = &ObjectiveV2{ID: "O002", Status: ColumnPlanned}
	index.Objectives["O001"] = &ObjectiveV2{ID: "O001", Status: ColumnPlanned, DependsOn: []string{"O002"}}
	index.Tasks["T001"] = &TaskV2{ID: "T001", Objective: "O001", Status: ColumnPlanned}
	index.ObjectiveTasks["O001"] = []string{"T001"}
	router := &RouterStateV2{State: RouterPhaseTask, Objective: "O001", Task: "T001"}

	next := ResolveNext(NextInput{Index: index, Router: router})
	if next.Kind != NextDependency {
		t.Fatalf("Kind = %q, want dependency", next.Kind)
	}
	if next.GateDecision == nil || len(next.GateDecision.Blockers) == 0 || next.GateDecision.Blockers[0].Kind != GateBlockObjectiveDependency {
		t.Fatalf("GateDecision = %+v, want a GateBlockObjectiveDependency blocker", next.GateDecision)
	}
	if dep := next.GateDecision.Blockers[0].ObjectiveDependency; dep == nil || dep.Target != "O002" {
		t.Errorf("ObjectiveDependency block = %+v, want target O002", dep)
	}
}

// TestResolveNext_executeWhenTaskPlannedAndReady proves a planned Task with
// satisfied dependencies reaches NextExecute under executor authority.
func TestResolveNext_executeWhenTaskPlannedAndReady(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O001"] = &ObjectiveV2{ID: "O001", Status: ColumnPlanned}
	index.Tasks["T001"] = &TaskV2{ID: "T001", Objective: "O001", Status: ColumnPlanned}
	index.ObjectiveTasks["O001"] = []string{"T001"}
	router := &RouterStateV2{State: RouterPhaseTask, Objective: "O001", Task: "T001"}

	next := ResolveNext(NextInput{Index: index, Router: router})
	if next.Kind != NextExecute {
		t.Fatalf("Kind = %q, want execute", next.Kind)
	}
	if next.GateDecision == nil || !next.GateDecision.Allowed || next.GateDecision.Actor != ActorRoleExecutor {
		t.Errorf("GateDecision = %+v, want allowed under executor authority", next.GateDecision)
	}
}

// TestResolveNext_executeWhenTaskInProgress proves the same for a Task mid
// build/test, read through ResolveTaskAdvance rather than ResolveTaskStart.
func TestResolveNext_executeWhenTaskInProgress(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O001"] = &ObjectiveV2{ID: "O001", Status: ColumnPlanned}
	index.Tasks["T001"] = &TaskV2{ID: "T001", Objective: "O001", Status: ColumnInProgress, Stage: StageBuild}
	index.ObjectiveTasks["O001"] = []string{"T001"}
	router := &RouterStateV2{State: RouterPhaseTask, Objective: "O001", Task: "T001"}

	next := ResolveNext(NextInput{Index: index, Router: router})
	if next.Kind != NextExecute {
		t.Fatalf("Kind = %q, want execute", next.Kind)
	}
}

// TestResolveNext_checkNeeded_missing, _needsWork, _stale, and _unknown
// prove each of the four non-current clearance states reaches the same
// NextCheckNeeded rung while carrying its own distinguishable
// Clearance.State — no restatement of freshness logic in next.go, only the
// value ResolveClearance already computed.
func TestResolveNext_checkNeeded_missing(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O001"] = &ObjectiveV2{ID: "O001", Status: ColumnPlanned}
	index.Tasks["T001"] = &TaskV2{ID: "T001", Objective: "O001", Status: ColumnInProgress, Stage: StageAudit}
	index.ObjectiveTasks["O001"] = []string{"T001"}
	router := &RouterStateV2{State: RouterPhaseCheck, Objective: "O001", Task: "T001"}

	next := ResolveNext(NextInput{Index: index, Router: router})
	if next.Kind != NextCheckNeeded || next.Clearance == nil || next.Clearance.State != ClearanceMissing {
		t.Fatalf("next = %+v, want check_needed/missing", next)
	}
}

func TestResolveNext_checkNeeded_needsWork(t *testing.T) {
	index := newV2TestIndex()
	mustCheck(index, "C001", "T001", CheckResultNeedsWork)
	index.Objectives["O001"] = &ObjectiveV2{ID: "O001", Status: ColumnPlanned}
	index.Tasks["T001"] = &TaskV2{ID: "T001", Objective: "O001", Status: ColumnInProgress, Stage: StageAudit}
	index.ObjectiveTasks["O001"] = []string{"T001"}
	router := &RouterStateV2{State: RouterPhaseCheck, Objective: "O001", Task: "T001"}

	next := ResolveNext(NextInput{Index: index, Router: router})
	if next.Kind != NextCheckNeeded || next.Clearance == nil || next.Clearance.State != ClearanceNeedsWork {
		t.Fatalf("next = %+v, want check_needed/needs_work", next)
	}
}

func TestResolveNext_checkNeeded_stale(t *testing.T) {
	index := newV2TestIndex()
	mustCheck(index, "C001", "T001", CheckResultClear)
	mustCheck(index, "C002", "T001", CheckResultClear) // supersedes C001 as the latest
	index.Objectives["O001"] = &ObjectiveV2{ID: "O001", Status: ColumnPlanned}
	index.Tasks["T001"] = &TaskV2{
		ID: "T001", Objective: "O001", Status: ColumnInProgress, Stage: StageAudit,
		Evidence: &Evidence{Freshness: &Freshness{State: FreshnessCurrent, Check: "C001", AssessedBy: Actor{Role: ActorRoleChecker, Session: "sess-1"}, Basis: "stale basis"}},
	}
	index.ObjectiveTasks["O001"] = []string{"T001"}
	router := &RouterStateV2{State: RouterPhaseCheck, Objective: "O001", Task: "T001"}

	next := ResolveNext(NextInput{Index: index, Router: router})
	if next.Kind != NextCheckNeeded || next.Clearance == nil || next.Clearance.State != ClearanceStale {
		t.Fatalf("next = %+v, want check_needed/stale", next)
	}
}

func TestResolveNext_checkNeeded_unknown(t *testing.T) {
	index := newV2TestIndex()
	mustCheck(index, "C001", "T001", CheckResultClear)
	index.Objectives["O001"] = &ObjectiveV2{ID: "O001", Status: ColumnPlanned}
	index.Tasks["T001"] = &TaskV2{ID: "T001", Objective: "O001", Status: ColumnInProgress, Stage: StageAudit}
	index.ObjectiveTasks["O001"] = []string{"T001"}
	router := &RouterStateV2{State: RouterPhaseCheck, Objective: "O001", Task: "T001"}

	next := ResolveNext(NextInput{Index: index, Router: router})
	if next.Kind != NextCheckNeeded || next.Clearance == nil || next.Clearance.State != ClearanceUnknown {
		t.Fatalf("next = %+v, want check_needed/unknown", next)
	}
}

// TestResolveNext_checkNeededOutranksOwnerValidation proves a Task needing
// both a fresh Check and owner validation reports the Check rung first: the
// switch inside ResolveTaskCompletion only ever adds an owner-acceptance
// blocker once clearance is already current.
func TestResolveNext_checkNeededOutranksOwnerValidation(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O001"] = &ObjectiveV2{ID: "O001", Status: ColumnPlanned}
	index.Tasks["T001"] = &TaskV2{
		ID: "T001", Objective: "O001", Status: ColumnInProgress, Stage: StageAudit,
		Evidence: &Evidence{OwnerValidation: &OwnerValidation{Required: true}},
	}
	index.ObjectiveTasks["O001"] = []string{"T001"}
	router := &RouterStateV2{State: RouterPhaseCheck, Objective: "O001", Task: "T001"}

	next := ResolveNext(NextInput{Index: index, Router: router})
	if next.Kind != NextCheckNeeded {
		t.Fatalf("Kind = %q, want check_needed (clearance is missing, so owner validation is not yet in view)", next.Kind)
	}
}

// TestResolveNext_ownerValidationRequiredAfterClearanceCurrent proves owner
// validation is reported only once clearance is current.
func TestResolveNext_ownerValidationRequiredAfterClearanceCurrent(t *testing.T) {
	index := newV2TestIndex()
	mustCheck(index, "C001", "T001", CheckResultClear)
	index.Objectives["O001"] = &ObjectiveV2{ID: "O001", Status: ColumnPlanned}
	index.Tasks["T001"] = mustCurrentTask("T001", "C001", &Evidence{OwnerValidation: &OwnerValidation{Required: true}})
	index.ObjectiveTasks["O001"] = []string{"T001"}
	router := &RouterStateV2{State: RouterPhaseCheck, Objective: "O001", Task: "T001"}

	next := ResolveNext(NextInput{Index: index, Router: router})
	if next.Kind != NextOwnerValidationRequired {
		t.Fatalf("Kind = %q, want owner_validation_required", next.Kind)
	}
	if next.Clearance == nil || next.Clearance.State != ClearanceCurrent {
		t.Errorf("Clearance = %+v, want current", next.Clearance)
	}
}

// TestResolveNext_exceptionAllowedCompletionReportsExecuteNotClearance
// proves a completion allowed only by a recorded exception reaches
// NextExecute carrying AllowedByException and the Exception itself, and
// carries no Clearance — it is never presented as a CLEAR result.
func TestResolveNext_exceptionAllowedCompletionReportsExecuteNotClearance(t *testing.T) {
	index := newV2TestIndex()
	mustCheck(index, "C001", "T001", CheckResultNeedsWork)
	index.Objectives["O001"] = &ObjectiveV2{ID: "O001", Status: ColumnPlanned}
	index.Tasks["T001"] = &TaskV2{
		ID: "T001", Objective: "O001", Status: ColumnInProgress, Stage: StageAudit,
		Evidence: &Evidence{Exception: &Exception{
			Requirements: []string{"R1"},
			Reason:       "known risk accepted",
			Owner:        "owner-1",
			Check:        "C001",
		}},
	}
	index.ObjectiveTasks["O001"] = []string{"T001"}
	router := &RouterStateV2{State: RouterPhaseCheck, Objective: "O001", Task: "T001"}

	next := ResolveNext(NextInput{Index: index, Router: router})
	if next.Kind != NextExecute {
		t.Fatalf("Kind = %q, want execute (allowed by exception)", next.Kind)
	}
	if next.GateDecision == nil || !next.GateDecision.AllowedByException || next.GateDecision.Exception == nil {
		t.Fatalf("GateDecision = %+v, want AllowedByException with Exception carried", next.GateDecision)
	}
	if next.Clearance != nil {
		t.Errorf("Clearance = %+v, want nil (an exception-allowed completion is never presented as clearance)", next.Clearance)
	}
}

// TestResolveNext_objectiveIntegrationRung proves an Objective whose owned
// Tasks are all done, but whose own integration clearance is not current,
// reaches NextObjectiveIntegration carrying ResolveObjectiveCompletion's
// decision and the Objective's own Clearance.
func TestResolveNext_objectiveIntegrationRung(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O001"] = &ObjectiveV2{ID: "O001", Status: ColumnInProgress}
	index.Tasks["T001"] = &TaskV2{ID: "T001", Objective: "O001", Status: ColumnDone}
	index.ObjectiveTasks["O001"] = []string{"T001"}
	router := &RouterStateV2{State: RouterPhaseCheck, Objective: "O001", Task: "T001"}

	next := ResolveNext(NextInput{Index: index, Router: router})
	if next.Kind != NextObjectiveIntegration {
		t.Fatalf("Kind = %q, want objective_integration", next.Kind)
	}
	if next.Objective == nil || next.Objective.ID != "O001" {
		t.Fatalf("Objective = %+v, want O001", next.Objective)
	}
	if next.GateDecision == nil || next.GateDecision.Allowed {
		t.Fatalf("GateDecision = %+v, want blocked (no objective check recorded)", next.GateDecision)
	}
	if next.Clearance == nil || next.Clearance.State != ClearanceMissing {
		t.Errorf("Clearance = %+v, want missing", next.Clearance)
	}
}

// TestResolveNext_objectiveIntegrationFallsThroughWhenTaskIncomplete proves
// that an Objective with an incomplete owned Task never reaches the
// integration rung: the real next action is that Task, found by the
// project-wide search.
func TestResolveNext_objectiveIntegrationFallsThroughWhenTaskIncomplete(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O001"] = &ObjectiveV2{ID: "O001", Status: ColumnInProgress}
	index.Tasks["T001"] = &TaskV2{ID: "T001", Objective: "O001", Status: ColumnPlanned}
	index.ObjectiveTasks["O001"] = []string{"T001"}
	router := &RouterStateV2{State: RouterPhaseDesign, Objective: "O001"}

	next := ResolveNext(NextInput{Index: index, Router: router})
	if next.Kind != NextReady {
		t.Fatalf("Kind = %q, want ready (T001 is the real next action, not integration)", next.Kind)
	}
	if next.Task == nil || next.Task.ID != "T001" {
		t.Fatalf("Task = %+v, want T001", next.Task)
	}
}

// TestResolveNext_readyRungPicksLowestSortedTaskID proves the project-wide
// search is deterministic: it always names the same ready Task first.
func TestResolveNext_readyRungPicksLowestSortedTaskID(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O001"] = &ObjectiveV2{ID: "O001", Status: ColumnPlanned}
	index.Tasks["T002"] = &TaskV2{ID: "T002", Objective: "O001", Status: ColumnPlanned}
	index.Tasks["T001"] = &TaskV2{ID: "T001", Objective: "O001", Status: ColumnPlanned}
	index.ObjectiveTasks["O001"] = []string{"T001", "T002"}
	router := &RouterStateV2{State: RouterPhaseIdea}

	next := ResolveNext(NextInput{Index: index, Router: router})
	if next.Kind != NextReady || next.Task == nil || next.Task.ID != "T001" {
		t.Fatalf("next = %+v, want ready/T001", next)
	}
}

// TestResolveNext_readyRungFallsBackToObjectiveWithNoTasks proves the
// project-wide search also finds a ready Objective — one with no Tasks yet
// and satisfied dependencies — when no Task anywhere is ready.
func TestResolveNext_readyRungFallsBackToObjectiveWithNoTasks(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O001"] = &ObjectiveV2{ID: "O001", Status: ColumnPlanned}
	router := &RouterStateV2{State: RouterPhaseIdea}

	next := ResolveNext(NextInput{Index: index, Router: router})
	if next.Kind != NextReady {
		t.Fatalf("Kind = %q, want ready", next.Kind)
	}
	if next.Objective == nil || next.Objective.ID != "O001" || next.Task != nil {
		t.Fatalf("next = %+v, want Objective O001 selected with no Task", next)
	}
}

// TestResolveNext_planObjectiveForEmptyProject proves a fresh project with
// no Objectives and no Tasks reaches NextPlanObjective with no diagnostic
// and no error.
func TestResolveNext_planObjectiveForEmptyProject(t *testing.T) {
	index := newV2TestIndex()
	router := &RouterStateV2{State: RouterPhaseIdea}

	next := ResolveNext(NextInput{Index: index, Router: router})
	if next.Kind != NextPlanObjective {
		t.Fatalf("Kind = %q, want plan_objective", next.Kind)
	}
	if next.SelectionDiagnostic != nil {
		t.Errorf("SelectionDiagnostic = %+v, want nil for an empty project with no router selection", next.SelectionDiagnostic)
	}
	if next.Objective != nil || next.Task != nil {
		t.Errorf("next = %+v, want no Objective or Task selected", next)
	}
}

// TestResolveNext_planObjectiveWhenNothingReady proves the same rung is
// reached by a project that is not empty but has nothing left ready: every
// Objective and Task is done.
func TestResolveNext_planObjectiveWhenNothingReady(t *testing.T) {
	index := newV2TestIndex()
	mustObjectiveCheck(index, "C001", "O001", CheckResultClear)
	index.Objectives["O001"] = mustCurrentObjective("O001", "C001", nil)
	index.Objectives["O001"].Status = ColumnDone
	index.Tasks["T001"] = &TaskV2{ID: "T001", Objective: "O001", Status: ColumnDone}
	index.ObjectiveTasks["O001"] = []string{"T001"}
	router := &RouterStateV2{State: RouterPhaseIdea}

	next := ResolveNext(NextInput{Index: index, Router: router})
	if next.Kind != NextPlanObjective {
		t.Fatalf("Kind = %q, want plan_objective (everything done, nothing ready)", next.Kind)
	}
}

// TestResolveNext_unresolvedSelectionStillYieldsAnAvailableAction proves an
// unresolved router selection reports its diagnostic and still yields the
// next action the project's own records support.
func TestResolveNext_unresolvedSelectionStillYieldsAnAvailableAction(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O001"] = &ObjectiveV2{ID: "O001", Status: ColumnPlanned}
	index.Tasks["T001"] = &TaskV2{ID: "T001", Objective: "O001", Status: ColumnPlanned}
	index.ObjectiveTasks["O001"] = []string{"T001"}
	router := &RouterStateV2{State: RouterPhaseTask, Objective: "O001", Task: "T999"} // T999 does not exist

	next := ResolveNext(NextInput{Index: index, Router: router})
	if next.SelectionDiagnostic == nil || next.SelectionDiagnostic.Kind != SelectionNotFound {
		t.Fatalf("SelectionDiagnostic = %+v, want SelectionNotFound for T999", next.SelectionDiagnostic)
	}
	if next.Kind != NextReady || next.Task == nil || next.Task.ID != "T001" {
		t.Fatalf("next = %+v, want ready/T001 derived from the records despite the unresolved selection", next)
	}
}

// TestResolveNext_noSelectionCarriesNoDiagnostic proves a router naming no
// selection at all — idea or design with no Objective — is not itself a
// diagnostic.
func TestResolveNext_noSelectionCarriesNoDiagnostic(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O001"] = &ObjectiveV2{ID: "O001", Status: ColumnPlanned}
	index.Tasks["T001"] = &TaskV2{ID: "T001", Objective: "O001", Status: ColumnPlanned}
	index.ObjectiveTasks["O001"] = []string{"T001"}
	router := &RouterStateV2{State: RouterPhaseIdea}

	next := ResolveNext(NextInput{Index: index, Router: router})
	if next.SelectionDiagnostic != nil {
		t.Errorf("SelectionDiagnostic = %+v, want nil when the router names no selection at all", next.SelectionDiagnostic)
	}
}

// TestResolveNext_readsNoFilesystem builds its V2Index directly in memory,
// with no discovery or file read anywhere in the call, proving ResolveNext
// consults only the values it is handed.
func TestResolveNext_readsNoFilesystem(t *testing.T) {
	index := &V2Index{
		Objectives:     map[string]*ObjectiveV2{"O001": {ID: "O001", Status: ColumnPlanned}},
		Tasks:          map[string]*TaskV2{"T001": {ID: "T001", Objective: "O001", Status: ColumnPlanned}},
		ObjectiveTasks: map[string][]string{"O001": {"T001"}},
		Checks:         map[string]*CheckV2{},
		ScopeChecks:    map[string][]string{},
		LatestCheck:    map[string]string{},
	}
	router := &RouterStateV2{State: RouterPhaseTask, Objective: "O001", Task: "T001"}

	next := ResolveNext(NextInput{Index: index, Router: router})
	if next.Kind != NextExecute {
		t.Fatalf("Kind = %q, want execute, from an in-memory-only index", next.Kind)
	}
}

// TestResolveNext_issuesLinkedToSelectedTask proves a selected Task's own
// linked Issues (TaskIssues) land on Next.Issues, in the same sorted order
// the index already keeps them in.
func TestResolveNext_issuesLinkedToSelectedTask(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O001"] = &ObjectiveV2{ID: "O001", Status: ColumnPlanned}
	index.Tasks["T001"] = &TaskV2{ID: "T001", Objective: "O001", Status: ColumnPlanned}
	index.ObjectiveTasks["O001"] = []string{"T001"}
	index.Issues = map[string]*IssueV2{
		"I001": {ID: "I001", Title: "Found during build", Type: IssueTypeDefect, Status: IssueStatusOpen},
		"I002": {ID: "I002", Title: "Unrelated", Type: IssueTypeDrift, Status: IssueStatusOpen},
	}
	index.TaskIssues = map[string][]string{"T001": {"I001"}}
	router := &RouterStateV2{State: RouterPhaseTask, Objective: "O001", Task: "T001"}

	next := ResolveNext(NextInput{Index: index, Router: router})
	if len(next.Issues) != 1 || next.Issues[0].ID != "I001" {
		t.Fatalf("Issues = %+v, want exactly [I001]", next.Issues)
	}
}

// TestResolveNext_issuesForObjectiveUnionOwnedTasksAndOwnChecks proves an
// Objective selected with no Task reports the union of every owned Task's
// linked Issues and any Issue linked to a Check scoped to the Objective
// itself — since IssueV2 carries no direct Objective reference — sorted and
// deduplicated.
func TestResolveNext_issuesForObjectiveUnionOwnedTasksAndOwnChecks(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O001"] = &ObjectiveV2{ID: "O001", Status: ColumnPlanned}
	index.Tasks["T001"] = &TaskV2{ID: "T001", Objective: "O001", Status: ColumnDone}
	index.ObjectiveTasks["O001"] = []string{"T001"}
	index.Issues = map[string]*IssueV2{
		"I001": {ID: "I001", Title: "From owned task", Type: IssueTypeDefect, Status: IssueStatusOpen},
		"I002": {ID: "I002", Title: "From objective's own check", Type: IssueTypeGuardrail, Status: IssueStatusOpen},
	}
	index.TaskIssues = map[string][]string{"T001": {"I001"}}
	index.ScopeChecks = map[string][]string{"O001": {"C001"}}
	index.CheckIssues = map[string][]string{"C001": {"I002"}}
	router := &RouterStateV2{State: RouterPhaseDesign, Objective: "O001"}

	next := ResolveNext(NextInput{Index: index, Router: router})
	if len(next.Issues) != 2 || next.Issues[0].ID != "I001" || next.Issues[1].ID != "I002" {
		t.Fatalf("Issues = %+v, want [I001 I002] sorted", next.Issues)
	}
}

// TestResolveNext_noIssuesLeavesNilNotEmpty proves a selected Task with no
// linked Issues reports a nil slice, so a rendering surface can tell "no
// Issues" apart from "an empty Issues section".
func TestResolveNext_noIssuesLeavesNilNotEmpty(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O001"] = &ObjectiveV2{ID: "O001", Status: ColumnPlanned}
	index.Tasks["T001"] = &TaskV2{ID: "T001", Objective: "O001", Status: ColumnPlanned}
	index.ObjectiveTasks["O001"] = []string{"T001"}
	router := &RouterStateV2{State: RouterPhaseTask, Objective: "O001", Task: "T001"}

	next := ResolveNext(NextInput{Index: index, Router: router})
	if next.Issues != nil {
		t.Fatalf("Issues = %+v, want nil", next.Issues)
	}
}

// TestNext_packageDoesNotImportMigrate proves internal/data does not import
// internal/migrate. internal/migrate already imports internal/data for the
// V2 records it converts into, so the reverse import would close a cycle —
// which is exactly why MigrationState exists as an injected value instead.
// Mirrors internal/migrate's own TestReplaceFile_doesNotReachAtomicWrite.
func TestNext_packageDoesNotImportMigrate(t *testing.T) {
	if _, err := os.Stat("next.go"); err != nil {
		t.Skipf("package source is not beside the test binary: %v", err)
	}
	pkg, err := build.ImportDir(".", build.IgnoreVendor)
	if err != nil {
		t.Fatalf("scan package imports: %v", err)
	}
	for _, imported := range append(pkg.Imports, pkg.TestImports...) {
		if strings.HasSuffix(imported, "/internal/migrate") {
			t.Errorf("internal/data imports %s, which already imports internal/data and would close a cycle", imported)
		}
	}
}
