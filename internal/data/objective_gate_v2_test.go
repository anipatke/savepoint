package data

import (
	"strings"
	"testing"
)

func mustObjectiveCheck(index *V2Index, id, objectiveID string, result CheckResult) *CheckV2 {
	check := &CheckV2{
		ID:        id,
		Scope:     CheckScope{Kind: CheckScopeObjective, ID: objectiveID},
		Result:    result,
		CheckedBy: Actor{Role: ActorRoleChecker, Session: "sess-1"},
		Source:    V2SourceDocument{Path: "checks/" + id + ".md"},
	}
	index.Checks[id] = check
	index.LatestCheck[objectiveID] = id
	return check
}

func mustCurrentObjective(id, checkID string, extra *Evidence) *ObjectiveV2 {
	evidence := extra
	if evidence == nil {
		evidence = &Evidence{}
	}
	evidence.Freshness = &Freshness{State: FreshnessCurrent, Check: checkID, AssessedBy: Actor{Role: ActorRoleChecker, Session: "sess-1"}, Basis: "reviewed"}
	return &ObjectiveV2{ID: id, Status: ColumnInProgress, Evidence: evidence}
}

func TestResolveObjectiveCompletion_blocksOnEveryUnfinishedTask(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Status: ColumnInProgress}
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnDone}
	index.Tasks["T-002"] = &TaskV2{ID: "T-002", Objective: "O-001", Status: ColumnInProgress}
	index.Tasks["T-003"] = &TaskV2{ID: "T-003", Objective: "O-001", Status: ColumnPlanned}
	index.ObjectiveTasks["O-001"] = []string{"T-001", "T-002", "T-003"}

	got := ResolveObjectiveCompletion(index, "O-001")
	if got.Allowed {
		t.Fatalf("Allowed = true, want false (two unfinished tasks)")
	}
	if len(got.Blockers) != 2 {
		t.Fatalf("Blockers = %+v, want 2 (one per unfinished task)", got.Blockers)
	}
	for _, b := range got.Blockers {
		if b.Kind != GateBlockInvalidState {
			t.Errorf("Blocker kind = %q, want invalid_state", b.Kind)
		}
	}
	if got.Blockers[0].Detail == got.Blockers[1].Detail {
		t.Errorf("expected each unfinished task named distinctly, got duplicate detail %q", got.Blockers[0].Detail)
	}
}

func TestResolveObjectiveCompletion_allowedWhenTasksDoneAndClearanceCurrent(t *testing.T) {
	index := newV2TestIndex()
	mustObjectiveCheck(index, "C-001", "O-001", CheckResultClear)
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnDone}
	index.ObjectiveTasks["O-001"] = []string{"T-001"}
	index.Objectives["O-001"] = mustCurrentObjective("O-001", "C-001", nil)

	got := ResolveObjectiveCompletion(index, "O-001")
	if !got.Allowed {
		t.Fatalf("Allowed = false, want true, blockers = %+v", got.Blockers)
	}
	if got.Actor != ActorRoleChecker {
		t.Errorf("Actor = %q, want checker", got.Actor)
	}
	if got.AllowedByException {
		t.Errorf("AllowedByException = true, want false (allowed under checker authority)")
	}
}

func TestResolveObjectiveCompletion_emptyObjectiveBlockedAsMissingRatherThanAllowed(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Status: ColumnPlanned}

	got := ResolveObjectiveCompletion(index, "O-001")
	if got.Allowed {
		t.Fatalf("Allowed = true, want false (no owned tasks and no check)")
	}
	if len(got.Blockers) != 1 || got.Blockers[0].Kind != GateBlockClearanceMissing {
		t.Fatalf("Blockers = %+v, want one GateBlockClearanceMissing", got.Blockers)
	}
}

func TestResolveObjectiveCompletion_eachClearanceStateBlocksWithADistinctReason(t *testing.T) {
	cases := []struct {
		name     string
		build    func(index *V2Index)
		evidence *Evidence
		wantKind GateBlockKind
	}{
		{
			name:     "missing",
			build:    func(index *V2Index) {},
			wantKind: GateBlockClearanceMissing,
		},
		{
			name:     "needs_work",
			build:    func(index *V2Index) { mustObjectiveCheck(index, "C-001", "O-001", CheckResultNeedsWork) },
			wantKind: GateBlockClearanceNeedsWork,
		},
		{
			name:     "unknown",
			build:    func(index *V2Index) { mustObjectiveCheck(index, "C-001", "O-001", CheckResultClear) },
			wantKind: GateBlockClearanceUnknown,
		},
		{
			name: "stale",
			build: func(index *V2Index) {
				mustObjectiveCheck(index, "C-001", "O-001", CheckResultClear)
			},
			evidence: &Evidence{Freshness: &Freshness{State: FreshnessStale, Check: "C-001", Basis: "flagged stale"}},
			wantKind: GateBlockClearanceStale,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			index := newV2TestIndex()
			index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnDone}
			index.ObjectiveTasks["O-001"] = []string{"T-001"}
			tc.build(index)
			index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Status: ColumnInProgress, Evidence: tc.evidence}

			got := ResolveObjectiveCompletion(index, "O-001")
			if got.Allowed {
				t.Fatalf("Allowed = true, want false")
			}
			if len(got.Blockers) != 1 || got.Blockers[0].Kind != tc.wantKind {
				t.Fatalf("Blockers = %+v, want one %s", got.Blockers, tc.wantKind)
			}
		})
	}
}

func TestResolveObjectiveCompletion_ownerValidationRequiredBlocksOnClearanceAlone(t *testing.T) {
	index := newV2TestIndex()
	mustObjectiveCheck(index, "C-001", "O-001", CheckResultClear)
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnDone}
	index.ObjectiveTasks["O-001"] = []string{"T-001"}
	index.Objectives["O-001"] = mustCurrentObjective("O-001", "C-001", &Evidence{OwnerValidation: &OwnerValidation{Required: true}})

	got := ResolveObjectiveCompletion(index, "O-001")
	if got.Allowed {
		t.Fatalf("Allowed = true, want false (owner has not accepted)")
	}
	if len(got.Blockers) != 1 || got.Blockers[0].Kind != GateBlockOwnerAcceptance {
		t.Fatalf("Blockers = %+v, want one GateBlockOwnerAcceptance", got.Blockers)
	}
}

func TestResolveObjectiveCompletion_ownerValidationSatisfiedWhenAcceptedCurrentCheck(t *testing.T) {
	index := newV2TestIndex()
	mustObjectiveCheck(index, "C-001", "O-001", CheckResultClear)
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnDone}
	index.ObjectiveTasks["O-001"] = []string{"T-001"}
	index.Objectives["O-001"] = mustCurrentObjective("O-001", "C-001", &Evidence{
		OwnerValidation: &OwnerValidation{Required: true, AcceptedCheck: "C-001", AcceptedBy: Actor{Role: ActorRoleOwner, Session: "owner-1"}},
	})

	got := ResolveObjectiveCompletion(index, "O-001")
	if !got.Allowed {
		t.Fatalf("Allowed = false, want true (owner accepted current check), blockers = %+v", got.Blockers)
	}
}

func TestResolveObjectiveCompletion_acceptanceOfSupersededCheckDoesNotClose(t *testing.T) {
	index := newV2TestIndex()
	mustObjectiveCheck(index, "C-001", "O-001", CheckResultClear)
	mustObjectiveCheck(index, "C-002", "O-001", CheckResultClear) // supersedes C-001 as the latest
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnDone}
	index.ObjectiveTasks["O-001"] = []string{"T-001"}
	index.Objectives["O-001"] = mustCurrentObjective("O-001", "C-002", &Evidence{
		OwnerValidation: &OwnerValidation{Required: true, AcceptedCheck: "C-001", AcceptedBy: Actor{Role: ActorRoleOwner, Session: "owner-1"}},
	})

	got := ResolveObjectiveCompletion(index, "O-001")
	if got.Allowed {
		t.Fatalf("Allowed = true, want false (acceptance bound to a superseded check)")
	}
	if len(got.Blockers) != 1 || got.Blockers[0].Kind != GateBlockOwnerAcceptance {
		t.Fatalf("Blockers = %+v, want one GateBlockOwnerAcceptance", got.Blockers)
	}
}

func TestResolveObjectiveCompletion_allowedByExceptionWhenOtherwiseBlocked(t *testing.T) {
	index := newV2TestIndex()
	mustObjectiveCheck(index, "C-001", "O-001", CheckResultNeedsWork)
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnDone}
	index.ObjectiveTasks["O-001"] = []string{"T-001"}
	index.Objectives["O-001"] = &ObjectiveV2{
		ID:     "O-001",
		Status: ColumnInProgress,
		Evidence: &Evidence{
			Exception: &Exception{Requirements: []string{"integration-check"}, Reason: "shipping deadline", Owner: "owner-1", Check: "C-001"},
		},
	}

	got := ResolveObjectiveCompletion(index, "O-001")
	if !got.Allowed {
		t.Fatalf("Allowed = false, want true (exception names the latest check), blockers = %+v", got.Blockers)
	}
	if !got.AllowedByException || got.Exception == nil {
		t.Fatalf("AllowedByException = %v, Exception = %v, want exception attached", got.AllowedByException, got.Exception)
	}
	if got.Actor != ActorRoleOwner {
		t.Errorf("Actor = %q, want owner", got.Actor)
	}
}

func TestResolveObjectiveCompletion_exceptionNamingOtherCheckDoesNotApply(t *testing.T) {
	index := newV2TestIndex()
	mustObjectiveCheck(index, "C-001", "O-001", CheckResultNeedsWork)
	mustObjectiveCheck(index, "C-002", "O-001", CheckResultNeedsWork) // supersedes C-001 as the latest
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnDone}
	index.ObjectiveTasks["O-001"] = []string{"T-001"}
	index.Objectives["O-001"] = &ObjectiveV2{
		ID:     "O-001",
		Status: ColumnInProgress,
		Evidence: &Evidence{
			Exception: &Exception{Requirements: []string{"integration-check"}, Reason: "shipping deadline", Owner: "owner-1", Check: "C-001"},
		},
	}

	got := ResolveObjectiveCompletion(index, "O-001")
	if got.Allowed {
		t.Fatalf("Allowed = true, want false (exception names a superseded check, not the latest)")
	}
	if len(got.Blockers) != 1 || got.Blockers[0].Kind != GateBlockClearanceNeedsWork {
		t.Fatalf("Blockers = %+v, want one GateBlockClearanceNeedsWork", got.Blockers)
	}
}

func TestResolveObjectiveCompletion_unfinishedTaskIsNotExcusableByException(t *testing.T) {
	index := newV2TestIndex()
	mustObjectiveCheck(index, "C-001", "O-001", CheckResultClear)
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnInProgress}
	index.ObjectiveTasks["O-001"] = []string{"T-001"}
	index.Objectives["O-001"] = mustCurrentObjective("O-001", "C-001", &Evidence{
		Exception: &Exception{Requirements: []string{"tasks-done"}, Reason: "shipping deadline", Owner: "owner-1", Check: "C-001"},
	})

	got := ResolveObjectiveCompletion(index, "O-001")
	if got.Allowed {
		t.Fatalf("Allowed = true, want false (an unfinished task is not excusable by exception)")
	}
	if len(got.Blockers) != 1 || got.Blockers[0].Kind != GateBlockInvalidState {
		t.Fatalf("Blockers = %+v, want one GateBlockInvalidState", got.Blockers)
	}
}

func TestResolveObjectiveCompletion_unknownObjectiveReturnsZeroDecision(t *testing.T) {
	index := newV2TestIndex()

	got := ResolveObjectiveCompletion(index, "O-404")
	if got.Allowed || len(got.Blockers) != 0 {
		t.Fatalf("got = %+v, want zero-value decision for an unknown objective", got)
	}
}

func TestResolveObjectiveDependency_missingObjectiveBlocksAsNotDone(t *testing.T) {
	index := newV2TestIndex()

	got := ResolveObjectiveDependency(index, "O-404")
	if got.Satisfied {
		t.Fatalf("Satisfied = true, want false (no such objective)")
	}
	if got.Block == nil || got.Block.Kind != ObjectiveDependencyBlockNotDone || got.Block.Target != "O-404" {
		t.Fatalf("Block = %+v, want NotDone naming O-404", got.Block)
	}
}

func TestResolveObjectiveDependency_notDoneBlocksReadiness(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-002"] = &ObjectiveV2{ID: "O-002", Status: ColumnInProgress}

	got := ResolveObjectiveDependency(index, "O-002")
	if got.Satisfied {
		t.Fatalf("Satisfied = true, want false (dependency objective not done)")
	}
	if got.Block == nil || got.Block.Kind != ObjectiveDependencyBlockNotDone || got.Block.Target != "O-002" {
		t.Fatalf("Block = %+v, want NotDone naming O-002", got.Block)
	}
}

func TestResolveObjectiveDependency_satisfiedWhenDoneAndClearanceCurrent(t *testing.T) {
	index := newV2TestIndex()
	mustObjectiveCheck(index, "C-001", "O-002", CheckResultClear)
	dependency := mustCurrentObjective("O-002", "C-001", nil)
	dependency.Status = ColumnDone
	index.Objectives["O-002"] = dependency

	got := ResolveObjectiveDependency(index, "O-002")
	if !got.Satisfied {
		t.Fatalf("Satisfied = false, want true, block = %+v", got.Block)
	}
	if got.Block != nil {
		t.Errorf("Block = %+v, want nil when satisfied", got.Block)
	}
}

func TestResolveObjectiveDependency_eachClearanceStateBlocksWithItsOwnReason(t *testing.T) {
	cases := []struct {
		name     string
		build    func(index *V2Index)
		evidence *Evidence
	}{
		{
			name:  "missing",
			build: func(index *V2Index) {},
		},
		{
			name:  "needs_work",
			build: func(index *V2Index) { mustObjectiveCheck(index, "C-001", "O-002", CheckResultNeedsWork) },
		},
		{
			name:  "unknown",
			build: func(index *V2Index) { mustObjectiveCheck(index, "C-001", "O-002", CheckResultClear) },
		},
		{
			name:     "stale",
			build:    func(index *V2Index) { mustObjectiveCheck(index, "C-001", "O-002", CheckResultClear) },
			evidence: &Evidence{Freshness: &Freshness{State: FreshnessStale, Check: "C-001", Basis: "flagged stale"}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			index := newV2TestIndex()
			tc.build(index)
			index.Objectives["O-002"] = &ObjectiveV2{ID: "O-002", Status: ColumnDone, Evidence: tc.evidence}

			got := ResolveObjectiveDependency(index, "O-002")
			if got.Satisfied {
				t.Fatalf("Satisfied = true, want false")
			}
			wantClearance := ClearanceState(tc.name)
			if got.Block == nil || got.Block.Kind != ObjectiveDependencyBlockNotCleared || got.Block.Clearance != wantClearance {
				t.Fatalf("Block = %+v, want NotCleared with clearance %q", got.Block, wantClearance)
			}
		})
	}
}

func TestResolveObjectiveDependency_clearedByExceptionReportedDistinctlyNotAsCurrent(t *testing.T) {
	index := newV2TestIndex()
	mustObjectiveCheck(index, "C-001", "O-002", CheckResultNeedsWork)
	index.Objectives["O-002"] = &ObjectiveV2{
		ID:     "O-002",
		Status: ColumnDone,
		Evidence: &Evidence{
			Exception: &Exception{Requirements: []string{"integration-check"}, Reason: "shipping deadline", Owner: "owner-1", Check: "C-001"},
		},
	}

	got := ResolveObjectiveDependency(index, "O-002")
	if got.Satisfied {
		t.Fatalf("Satisfied = true, want false (an exception is not current clearance)")
	}
	if got.Block == nil || got.Block.Kind != ObjectiveDependencyBlockClearedByException {
		t.Fatalf("Block = %+v, want ClearedByException, reported distinctly from NotCleared", got.Block)
	}
}

func TestResolveObjectiveDependency_exceptionNamingOtherCheckStillBlocksAsNotCleared(t *testing.T) {
	index := newV2TestIndex()
	mustObjectiveCheck(index, "C-001", "O-002", CheckResultNeedsWork)
	mustObjectiveCheck(index, "C-002", "O-002", CheckResultNeedsWork) // supersedes C-001 as the latest
	index.Objectives["O-002"] = &ObjectiveV2{
		ID:     "O-002",
		Status: ColumnDone,
		Evidence: &Evidence{
			Exception: &Exception{Requirements: []string{"integration-check"}, Reason: "shipping deadline", Owner: "owner-1", Check: "C-001"},
		},
	}

	got := ResolveObjectiveDependency(index, "O-002")
	if got.Satisfied {
		t.Fatalf("Satisfied = true, want false")
	}
	if got.Block == nil || got.Block.Kind != ObjectiveDependencyBlockNotCleared || got.Block.Clearance != ClearanceNeedsWork {
		t.Fatalf("Block = %+v, want NotCleared/needs_work (exception names a superseded check, not the latest)", got.Block)
	}
}

// TestResolveTaskCompletion_unchangedByObjectiveGate is the regression case
// required alongside ResolveObjectiveCompletion: adding Objective completion
// must not alter a single Task completion decision.
func TestResolveTaskCompletion_unchangedByObjectiveGate(t *testing.T) {
	index := newV2TestIndex()
	mustCheck(index, "C-001", "T-001", CheckResultClear)
	index.Tasks["T-001"] = mustCurrentTask("T-001", "C-001", &Evidence{OwnerValidation: &OwnerValidation{Required: true}})

	got := ResolveTaskCompletion(index, "T-001")
	if got.Allowed {
		t.Fatalf("Allowed = true, want false (owner has not accepted)")
	}
	if len(got.Blockers) != 1 || got.Blockers[0].Kind != GateBlockOwnerAcceptance {
		t.Fatalf("Blockers = %+v, want one GateBlockOwnerAcceptance", got.Blockers)
	}

	index.Tasks["T-001"].Evidence.OwnerValidation.AcceptedCheck = "C-001"
	index.Tasks["T-001"].Evidence.OwnerValidation.AcceptedBy = Actor{Role: ActorRoleOwner, Session: "owner-1"}

	got = ResolveTaskCompletion(index, "T-001")
	if !got.Allowed || got.Actor != ActorRoleChecker {
		t.Fatalf("got = %+v, want allowed under checker authority once owner accepts", got)
	}
}

func TestResolveObjectiveCompletion_repairAndRecheckDoesNotRequireRetreatingDoneTasks(t *testing.T) {
	// I-019: an Objective Check's NEEDS WORK must be repairable and rechecked
	// without retreating a Task that is already done. ResolveObjectiveCompletion
	// only ever reads Task status; it never requires flipping it back to
	// in_progress to unblock on a fresh, superseding Check.
	index := newV2TestIndex()
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnDone}
	index.Tasks["T-002"] = &TaskV2{ID: "T-002", Objective: "O-001", Status: ColumnDone}
	index.ObjectiveTasks["O-001"] = []string{"T-001", "T-002"}

	mustObjectiveCheck(index, "C-001", "O-001", CheckResultNeedsWork)
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Status: ColumnInProgress}

	got := ResolveObjectiveCompletion(index, "O-001")
	if got.Allowed {
		t.Fatalf("Allowed = true, want false (objective check recorded NEEDS WORK)")
	}
	if len(got.Blockers) != 1 || got.Blockers[0].Kind != GateBlockClearanceNeedsWork {
		t.Fatalf("Blockers = %+v, want one GateBlockClearanceNeedsWork", got.Blockers)
	}
	if index.Tasks["T-001"].Status != ColumnDone || index.Tasks["T-002"].Status != ColumnDone {
		t.Fatalf("owned tasks changed status while only the objective check was NEEDS WORK")
	}

	// The finding becomes new work under the same Objective. The two
	// completed Tasks are historical work and must not be reopened.
	index.Tasks["T-003"] = &TaskV2{ID: "T-003", Objective: "O-001", Status: ColumnPlanned}
	index.ObjectiveTasks["O-001"] = append(index.ObjectiveTasks["O-001"], "T-003")
	got = ResolveObjectiveCompletion(index, "O-001")
	if got.Allowed || len(got.Blockers) != 1 || got.Blockers[0].Kind != GateBlockInvalidState {
		t.Fatalf("new remediation Task must block Objective completion, got %+v", got)
	}
	if index.Tasks["T-001"].Status != ColumnDone || index.Tasks["T-002"].Status != ColumnDone {
		t.Fatal("adding remediation work reopened a completed Task")
	}

	// Finishing remediation alone cannot override the failed Check. A fresh,
	// superseding CLEAR Check is still required.
	index.Tasks["T-003"].Status = ColumnDone
	got = ResolveObjectiveCompletion(index, "O-001")
	if got.Allowed || len(got.Blockers) != 1 || got.Blockers[0].Kind != GateBlockClearanceNeedsWork {
		t.Fatalf("repair without recheck must remain blocked, got %+v", got)
	}
	mustObjectiveCheck(index, "C-002", "O-001", CheckResultClear)
	index.Objectives["O-001"] = mustCurrentObjective("O-001", "C-002", nil)
	index.Issues = map[string]*IssueV2{"I-001": {ID: "I-001", Status: IssueStatusOpen}}
	index.CheckIssues = map[string][]string{"C-002": {"I-001"}}

	got = ResolveObjectiveCompletion(index, "O-001")
	if got.Allowed || len(got.Blockers) != 1 || got.Blockers[0].Kind != GateBlockObjectiveIssueUnresolved || got.Blockers[0].Issue != "I-001" {
		t.Fatalf("current CLEAR Check with unresolved material Issue must block, got %+v", got)
	}
	index.Objectives["O-001"].Evidence.Exception = &Exception{
		Requirements: []string{"other-requirement"}, Reason: "accepted separate risk",
		Owner: "owner-1", Check: "C-002",
	}
	got = ResolveObjectiveCompletion(index, "O-001")
	if got.Allowed || got.Blockers[0].Kind != GateBlockObjectiveIssueUnresolved {
		t.Fatalf("unrelated Objective exception must not accept an open Issue, got %+v", got)
	}
	index.Objectives["O-001"].Evidence.Exception = nil
	index.Issues["I-001"].Status = IssueStatusResolved
	got = ResolveObjectiveCompletion(index, "O-001")
	if !got.Allowed {
		t.Fatalf("Allowed = false, want true after the recheck cleared, blockers = %+v", got.Blockers)
	}
	if index.Tasks["T-001"].Status != ColumnDone || index.Tasks["T-002"].Status != ColumnDone {
		t.Fatalf("owned tasks changed status during repair and recheck, want both to remain done throughout")
	}
}

func TestInspectObjectiveConsistency_ignoresObjectivesNotDone(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Status: ColumnInProgress, Source: V2SourceDocument{Path: "objectives/O-001-a/Objective.md"}}

	if got := InspectObjectiveConsistency(index); len(got) != 0 {
		t.Fatalf("InspectObjectiveConsistency() = %+v, want none for a not-done objective", got)
	}
}

func TestInspectObjectiveConsistency_doneWithoutCurrentClearance(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Status: ColumnDone, Source: V2SourceDocument{Path: "objectives/O-001-a/Objective.md"}}

	got := InspectObjectiveConsistency(index)
	if len(got) != 1 || got[0].Kind != ObjectiveConsistencyDoneWithoutClearance || got[0].Objective != "O-001" {
		t.Fatalf("InspectObjectiveConsistency() = %+v, want one ObjectiveConsistencyDoneWithoutClearance naming O-001", got)
	}
}

func TestInspectObjectiveConsistency_doneWithIncompleteTask(t *testing.T) {
	index := newV2TestIndex()
	mustObjectiveCheck(index, "C-001", "O-001", CheckResultClear)
	index.Objectives["O-001"] = mustCurrentObjective("O-001", "C-001", nil)
	index.Objectives["O-001"].Status = ColumnDone
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnInProgress}
	index.ObjectiveTasks["O-001"] = []string{"T-001"}

	got := InspectObjectiveConsistency(index)
	if len(got) != 1 || got[0].Kind != ObjectiveConsistencyIncompleteTask || got[0].Objective != "O-001" {
		t.Fatalf("InspectObjectiveConsistency() = %+v, want one ObjectiveConsistencyIncompleteTask naming O-001", got)
	}
	if !strings.Contains(got[0].Detail, "T-001") {
		t.Errorf("Detail = %q, want the incomplete task named", got[0].Detail)
	}
}

func TestInspectObjectiveConsistency_cleanObjectiveReportsNothing(t *testing.T) {
	index := newV2TestIndex()
	mustObjectiveCheck(index, "C-001", "O-001", CheckResultClear)
	index.Objectives["O-001"] = mustCurrentObjective("O-001", "C-001", nil)
	index.Objectives["O-001"].Status = ColumnDone
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnDone}
	index.ObjectiveTasks["O-001"] = []string{"T-001"}

	if got := InspectObjectiveConsistency(index); len(got) != 0 {
		t.Fatalf("InspectObjectiveConsistency() = %+v, want none for a consistent done objective", got)
	}
}

func TestInspectObjectiveConsistency_sortedOrderReturnsEveryProblem(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-002"] = &ObjectiveV2{ID: "O-002", Status: ColumnDone, Source: V2SourceDocument{Path: "objectives/O-002-b/Objective.md"}}
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Status: ColumnDone, Source: V2SourceDocument{Path: "objectives/O-001-a/Objective.md"}}

	got := InspectObjectiveConsistency(index)
	if len(got) != 2 {
		t.Fatalf("InspectObjectiveConsistency() = %+v, want 2 problems (one per done objective)", got)
	}
	if got[0].Objective != "O-001" || got[1].Objective != "O-002" {
		t.Fatalf("InspectObjectiveConsistency() order = [%s, %s], want sorted [O-001, O-002]", got[0].Objective, got[1].Objective)
	}
}
