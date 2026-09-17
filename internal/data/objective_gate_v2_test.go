package data

import "testing"

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
	index.Objectives["O001"] = &ObjectiveV2{ID: "O001", Status: ColumnInProgress}
	index.Tasks["T001"] = &TaskV2{ID: "T001", Objective: "O001", Status: ColumnDone}
	index.Tasks["T002"] = &TaskV2{ID: "T002", Objective: "O001", Status: ColumnInProgress}
	index.Tasks["T003"] = &TaskV2{ID: "T003", Objective: "O001", Status: ColumnPlanned}
	index.ObjectiveTasks["O001"] = []string{"T001", "T002", "T003"}

	got := ResolveObjectiveCompletion(index, "O001")
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
	mustObjectiveCheck(index, "C001", "O001", CheckResultClear)
	index.Tasks["T001"] = &TaskV2{ID: "T001", Objective: "O001", Status: ColumnDone}
	index.ObjectiveTasks["O001"] = []string{"T001"}
	index.Objectives["O001"] = mustCurrentObjective("O001", "C001", nil)

	got := ResolveObjectiveCompletion(index, "O001")
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
	index.Objectives["O001"] = &ObjectiveV2{ID: "O001", Status: ColumnPlanned}

	got := ResolveObjectiveCompletion(index, "O001")
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
			build:    func(index *V2Index) { mustObjectiveCheck(index, "C001", "O001", CheckResultNeedsWork) },
			wantKind: GateBlockClearanceNeedsWork,
		},
		{
			name:     "unknown",
			build:    func(index *V2Index) { mustObjectiveCheck(index, "C001", "O001", CheckResultClear) },
			wantKind: GateBlockClearanceUnknown,
		},
		{
			name: "stale",
			build: func(index *V2Index) {
				mustObjectiveCheck(index, "C001", "O001", CheckResultClear)
			},
			evidence: &Evidence{Freshness: &Freshness{State: FreshnessStale, Check: "C001", Basis: "flagged stale"}},
			wantKind: GateBlockClearanceStale,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			index := newV2TestIndex()
			index.Tasks["T001"] = &TaskV2{ID: "T001", Objective: "O001", Status: ColumnDone}
			index.ObjectiveTasks["O001"] = []string{"T001"}
			tc.build(index)
			index.Objectives["O001"] = &ObjectiveV2{ID: "O001", Status: ColumnInProgress, Evidence: tc.evidence}

			got := ResolveObjectiveCompletion(index, "O001")
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
	mustObjectiveCheck(index, "C001", "O001", CheckResultClear)
	index.Tasks["T001"] = &TaskV2{ID: "T001", Objective: "O001", Status: ColumnDone}
	index.ObjectiveTasks["O001"] = []string{"T001"}
	index.Objectives["O001"] = mustCurrentObjective("O001", "C001", &Evidence{OwnerValidation: &OwnerValidation{Required: true}})

	got := ResolveObjectiveCompletion(index, "O001")
	if got.Allowed {
		t.Fatalf("Allowed = true, want false (owner has not accepted)")
	}
	if len(got.Blockers) != 1 || got.Blockers[0].Kind != GateBlockOwnerAcceptance {
		t.Fatalf("Blockers = %+v, want one GateBlockOwnerAcceptance", got.Blockers)
	}
}

func TestResolveObjectiveCompletion_ownerValidationSatisfiedWhenAcceptedCurrentCheck(t *testing.T) {
	index := newV2TestIndex()
	mustObjectiveCheck(index, "C001", "O001", CheckResultClear)
	index.Tasks["T001"] = &TaskV2{ID: "T001", Objective: "O001", Status: ColumnDone}
	index.ObjectiveTasks["O001"] = []string{"T001"}
	index.Objectives["O001"] = mustCurrentObjective("O001", "C001", &Evidence{
		OwnerValidation: &OwnerValidation{Required: true, AcceptedCheck: "C001", AcceptedBy: Actor{Role: ActorRoleOwner, Session: "owner-1"}},
	})

	got := ResolveObjectiveCompletion(index, "O001")
	if !got.Allowed {
		t.Fatalf("Allowed = false, want true (owner accepted current check), blockers = %+v", got.Blockers)
	}
}

func TestResolveObjectiveCompletion_acceptanceOfSupersededCheckDoesNotClose(t *testing.T) {
	index := newV2TestIndex()
	mustObjectiveCheck(index, "C001", "O001", CheckResultClear)
	mustObjectiveCheck(index, "C002", "O001", CheckResultClear) // supersedes C001 as the latest
	index.Tasks["T001"] = &TaskV2{ID: "T001", Objective: "O001", Status: ColumnDone}
	index.ObjectiveTasks["O001"] = []string{"T001"}
	index.Objectives["O001"] = mustCurrentObjective("O001", "C002", &Evidence{
		OwnerValidation: &OwnerValidation{Required: true, AcceptedCheck: "C001", AcceptedBy: Actor{Role: ActorRoleOwner, Session: "owner-1"}},
	})

	got := ResolveObjectiveCompletion(index, "O001")
	if got.Allowed {
		t.Fatalf("Allowed = true, want false (acceptance bound to a superseded check)")
	}
	if len(got.Blockers) != 1 || got.Blockers[0].Kind != GateBlockOwnerAcceptance {
		t.Fatalf("Blockers = %+v, want one GateBlockOwnerAcceptance", got.Blockers)
	}
}

func TestResolveObjectiveCompletion_allowedByExceptionWhenOtherwiseBlocked(t *testing.T) {
	index := newV2TestIndex()
	mustObjectiveCheck(index, "C001", "O001", CheckResultNeedsWork)
	index.Tasks["T001"] = &TaskV2{ID: "T001", Objective: "O001", Status: ColumnDone}
	index.ObjectiveTasks["O001"] = []string{"T001"}
	index.Objectives["O001"] = &ObjectiveV2{
		ID:     "O001",
		Status: ColumnInProgress,
		Evidence: &Evidence{
			Exception: &Exception{Requirements: []string{"integration-check"}, Reason: "shipping deadline", Owner: "owner-1", Check: "C001"},
		},
	}

	got := ResolveObjectiveCompletion(index, "O001")
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
	mustObjectiveCheck(index, "C001", "O001", CheckResultNeedsWork)
	mustObjectiveCheck(index, "C002", "O001", CheckResultNeedsWork) // supersedes C001 as the latest
	index.Tasks["T001"] = &TaskV2{ID: "T001", Objective: "O001", Status: ColumnDone}
	index.ObjectiveTasks["O001"] = []string{"T001"}
	index.Objectives["O001"] = &ObjectiveV2{
		ID:     "O001",
		Status: ColumnInProgress,
		Evidence: &Evidence{
			Exception: &Exception{Requirements: []string{"integration-check"}, Reason: "shipping deadline", Owner: "owner-1", Check: "C001"},
		},
	}

	got := ResolveObjectiveCompletion(index, "O001")
	if got.Allowed {
		t.Fatalf("Allowed = true, want false (exception names a superseded check, not the latest)")
	}
	if len(got.Blockers) != 1 || got.Blockers[0].Kind != GateBlockClearanceNeedsWork {
		t.Fatalf("Blockers = %+v, want one GateBlockClearanceNeedsWork", got.Blockers)
	}
}

func TestResolveObjectiveCompletion_unfinishedTaskIsNotExcusableByException(t *testing.T) {
	index := newV2TestIndex()
	mustObjectiveCheck(index, "C001", "O001", CheckResultClear)
	index.Tasks["T001"] = &TaskV2{ID: "T001", Objective: "O001", Status: ColumnInProgress}
	index.ObjectiveTasks["O001"] = []string{"T001"}
	index.Objectives["O001"] = mustCurrentObjective("O001", "C001", &Evidence{
		Exception: &Exception{Requirements: []string{"tasks-done"}, Reason: "shipping deadline", Owner: "owner-1", Check: "C001"},
	})

	got := ResolveObjectiveCompletion(index, "O001")
	if got.Allowed {
		t.Fatalf("Allowed = true, want false (an unfinished task is not excusable by exception)")
	}
	if len(got.Blockers) != 1 || got.Blockers[0].Kind != GateBlockInvalidState {
		t.Fatalf("Blockers = %+v, want one GateBlockInvalidState", got.Blockers)
	}
}

func TestResolveObjectiveCompletion_unknownObjectiveReturnsZeroDecision(t *testing.T) {
	index := newV2TestIndex()

	got := ResolveObjectiveCompletion(index, "O404")
	if got.Allowed || len(got.Blockers) != 0 {
		t.Fatalf("got = %+v, want zero-value decision for an unknown objective", got)
	}
}

func TestResolveObjectiveDependency_missingObjectiveBlocksAsNotDone(t *testing.T) {
	index := newV2TestIndex()

	got := ResolveObjectiveDependency(index, "O404")
	if got.Satisfied {
		t.Fatalf("Satisfied = true, want false (no such objective)")
	}
	if got.Block == nil || got.Block.Kind != ObjectiveDependencyBlockNotDone || got.Block.Target != "O404" {
		t.Fatalf("Block = %+v, want NotDone naming O404", got.Block)
	}
}

func TestResolveObjectiveDependency_notDoneBlocksReadiness(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O002"] = &ObjectiveV2{ID: "O002", Status: ColumnInProgress}

	got := ResolveObjectiveDependency(index, "O002")
	if got.Satisfied {
		t.Fatalf("Satisfied = true, want false (dependency objective not done)")
	}
	if got.Block == nil || got.Block.Kind != ObjectiveDependencyBlockNotDone || got.Block.Target != "O002" {
		t.Fatalf("Block = %+v, want NotDone naming O002", got.Block)
	}
}

func TestResolveObjectiveDependency_satisfiedWhenDoneAndClearanceCurrent(t *testing.T) {
	index := newV2TestIndex()
	mustObjectiveCheck(index, "C001", "O002", CheckResultClear)
	dependency := mustCurrentObjective("O002", "C001", nil)
	dependency.Status = ColumnDone
	index.Objectives["O002"] = dependency

	got := ResolveObjectiveDependency(index, "O002")
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
			build: func(index *V2Index) { mustObjectiveCheck(index, "C001", "O002", CheckResultNeedsWork) },
		},
		{
			name:  "unknown",
			build: func(index *V2Index) { mustObjectiveCheck(index, "C001", "O002", CheckResultClear) },
		},
		{
			name:     "stale",
			build:    func(index *V2Index) { mustObjectiveCheck(index, "C001", "O002", CheckResultClear) },
			evidence: &Evidence{Freshness: &Freshness{State: FreshnessStale, Check: "C001", Basis: "flagged stale"}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			index := newV2TestIndex()
			tc.build(index)
			index.Objectives["O002"] = &ObjectiveV2{ID: "O002", Status: ColumnDone, Evidence: tc.evidence}

			got := ResolveObjectiveDependency(index, "O002")
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
	mustObjectiveCheck(index, "C001", "O002", CheckResultNeedsWork)
	index.Objectives["O002"] = &ObjectiveV2{
		ID:     "O002",
		Status: ColumnDone,
		Evidence: &Evidence{
			Exception: &Exception{Requirements: []string{"integration-check"}, Reason: "shipping deadline", Owner: "owner-1", Check: "C001"},
		},
	}

	got := ResolveObjectiveDependency(index, "O002")
	if got.Satisfied {
		t.Fatalf("Satisfied = true, want false (an exception is not current clearance)")
	}
	if got.Block == nil || got.Block.Kind != ObjectiveDependencyBlockClearedByException {
		t.Fatalf("Block = %+v, want ClearedByException, reported distinctly from NotCleared", got.Block)
	}
}

func TestResolveObjectiveDependency_exceptionNamingOtherCheckStillBlocksAsNotCleared(t *testing.T) {
	index := newV2TestIndex()
	mustObjectiveCheck(index, "C001", "O002", CheckResultNeedsWork)
	mustObjectiveCheck(index, "C002", "O002", CheckResultNeedsWork) // supersedes C001 as the latest
	index.Objectives["O002"] = &ObjectiveV2{
		ID:     "O002",
		Status: ColumnDone,
		Evidence: &Evidence{
			Exception: &Exception{Requirements: []string{"integration-check"}, Reason: "shipping deadline", Owner: "owner-1", Check: "C001"},
		},
	}

	got := ResolveObjectiveDependency(index, "O002")
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
	mustCheck(index, "C001", "T001", CheckResultClear)
	index.Tasks["T001"] = mustCurrentTask("T001", "C001", &Evidence{OwnerValidation: &OwnerValidation{Required: true}})

	got := ResolveTaskCompletion(index, "T001")
	if got.Allowed {
		t.Fatalf("Allowed = true, want false (owner has not accepted)")
	}
	if len(got.Blockers) != 1 || got.Blockers[0].Kind != GateBlockOwnerAcceptance {
		t.Fatalf("Blockers = %+v, want one GateBlockOwnerAcceptance", got.Blockers)
	}

	index.Tasks["T001"].Evidence.OwnerValidation.AcceptedCheck = "C001"
	index.Tasks["T001"].Evidence.OwnerValidation.AcceptedBy = Actor{Role: ActorRoleOwner, Session: "owner-1"}

	got = ResolveTaskCompletion(index, "T001")
	if !got.Allowed || got.Actor != ActorRoleChecker {
		t.Fatalf("got = %+v, want allowed under checker authority once owner accepts", got)
	}
}
