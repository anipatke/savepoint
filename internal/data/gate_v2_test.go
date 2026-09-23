package data

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

func mustCheck(index *V2Index, id string, scopeID string, result CheckResult) *CheckV2 {
	check := &CheckV2{
		ID:        id,
		Scope:     CheckScope{Kind: CheckScopeTask, ID: scopeID},
		Result:    result,
		CheckedBy: Actor{Role: ActorRoleChecker, Session: "sess-1"},
		Source:    V2SourceDocument{Path: "checks/" + id + ".md"},
	}
	index.Checks[id] = check
	index.LatestCheck[scopeID] = id
	return check
}

func TestResolveClearance_missingWhenNoCheck(t *testing.T) {
	index := newV2TestIndex()
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001"}

	got := ResolveClearance(index, "T-001")
	if got.State != ClearanceMissing {
		t.Fatalf("State = %q, want missing", got.State)
	}
	if got.Check != "" {
		t.Errorf("Check = %q, want empty", got.Check)
	}
}

func TestResolveClearance_needsWorkFromLatestCheck(t *testing.T) {
	index := newV2TestIndex()
	mustCheck(index, "C-001", "T-001", CheckResultNeedsWork)
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001"}

	got := ResolveClearance(index, "T-001")
	if got.State != ClearanceNeedsWork {
		t.Fatalf("State = %q, want needs_work", got.State)
	}
	if got.Check != "C-001" {
		t.Errorf("Check = %q, want C-001", got.Check)
	}
}

// A CLEAR Check signed by a checker is current on its own; no separate
// freshness assessment is required.
func TestResolveClearance_currentWhenNoFreshness(t *testing.T) {
	index := newV2TestIndex()
	mustCheck(index, "C-001", "T-001", CheckResultClear)
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001"}

	got := ResolveClearance(index, "T-001")
	if got.State != ClearanceCurrent {
		t.Fatalf("State = %q, want current", got.State)
	}
	if got.Check != "C-001" {
		t.Errorf("Check = %q, want C-001", got.Check)
	}
}

func TestResolveClearance_currentWhenFreshnessNamesLatestCheckAsCurrent(t *testing.T) {
	index := newV2TestIndex()
	mustCheck(index, "C-001", "T-001", CheckResultClear)
	index.Tasks["T-001"] = &TaskV2{
		ID: "T-001", Objective: "O-001",
		Evidence: &Evidence{Freshness: &Freshness{State: FreshnessCurrent, Check: "C-001", AssessedBy: Actor{Role: ActorRoleChecker, Session: "sess-1"}, Basis: "re-read the diff"}},
	}

	got := ResolveClearance(index, "T-001")
	if got.State != ClearanceCurrent {
		t.Fatalf("State = %q, want current", got.State)
	}
	if got.Check != "C-001" {
		t.Errorf("Check = %q, want C-001", got.Check)
	}
	if got.Freshness == nil || got.Freshness.Basis != "re-read the diff" {
		t.Errorf("Freshness = %+v, want basis %q carried through", got.Freshness, "re-read the diff")
	}
}

func TestResolveClearance_reviewedBasisDoesNotAffectCurrentClearance(t *testing.T) {
	cases := []struct {
		name     string
		reviewed *ReviewedBasis
	}{
		{name: "absent"},
		{name: "empty", reviewed: &ReviewedBasis{Files: []string{}, Dependencies: []string{}}},
		{name: "substantive", reviewed: &ReviewedBasis{
			BaseCommit: "abc123", HeadCommit: "def456",
			Files: []string{"internal/data/check_v2.go"}, Dependencies: []string{"go.mod"},
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			index := newV2TestIndex()
			check := mustCheck(index, "C-001", "T-001", CheckResultClear)
			check.Reviewed = tc.reviewed
			index.Tasks["T-001"] = &TaskV2{
				ID: "T-001", Objective: "O-001",
				Evidence: &Evidence{Freshness: &Freshness{
					State: FreshnessCurrent, Check: "C-001",
					AssessedBy: Actor{Role: ActorRoleChecker, Session: "sess-2"},
					Basis:      "re-read the diff",
				}},
			}

			got := ResolveClearance(index, "T-001")
			if got.State != ClearanceCurrent {
				t.Fatalf("ResolveClearance() = %+v, want current with %s reviewed basis", got, tc.name)
			}
		})
	}
}

// A freshness assessment naming a superseded Check says nothing about the
// newer run, so a later CLEAR re-check is current even over a stale note.
func TestResolveClearance_newerCheckSupersedesFreshnessNamingAnOlderCheck(t *testing.T) {
	index := newV2TestIndex()
	mustCheck(index, "C-001", "T-001", CheckResultClear)
	mustCheck(index, "C-002", "T-001", CheckResultClear) // supersedes C-001 as the latest
	index.Tasks["T-001"] = &TaskV2{
		ID: "T-001", Objective: "O-001",
		Evidence: &Evidence{Freshness: &Freshness{State: FreshnessStale, Check: "C-001", Basis: "stale basis"}},
	}

	got := ResolveClearance(index, "T-001")
	if got.State != ClearanceCurrent {
		t.Fatalf("State = %q, want current (freshness names a superseded check)", got.State)
	}
	if got.Freshness != nil {
		t.Errorf("Freshness = %+v, want nil (the assessment names an older check)", got.Freshness)
	}
	if got.Check != "C-002" {
		t.Errorf("Check = %q, want C-002 (the actual latest)", got.Check)
	}
}

func TestResolveClearance_staleWhenFreshnessNamesLatestButNotCurrentState(t *testing.T) {
	index := newV2TestIndex()
	mustCheck(index, "C-001", "T-001", CheckResultClear)
	index.Tasks["T-001"] = &TaskV2{
		ID: "T-001", Objective: "O-001",
		Evidence: &Evidence{Freshness: &Freshness{State: FreshnessStale, Check: "C-001", Basis: "flagged stale"}},
	}

	got := ResolveClearance(index, "T-001")
	if got.State != ClearanceStale {
		t.Fatalf("State = %q, want stale", got.State)
	}
}

func TestResolveClearance_resolvesObjectiveTargetsToo(t *testing.T) {
	index := newV2TestIndex()
	index.Checks["C-001"] = &CheckV2{ID: "C-001", Scope: CheckScope{Kind: CheckScopeObjective, ID: "O-001"}, Result: CheckResultClear, CheckedBy: Actor{Role: ActorRoleChecker, Session: "sess-1"}}
	index.LatestCheck["O-001"] = "C-001"
	index.Objectives["O-001"] = &ObjectiveV2{
		ID:       "O-001",
		Evidence: &Evidence{Freshness: &Freshness{State: FreshnessCurrent, Check: "C-001", AssessedBy: Actor{Role: ActorRoleChecker, Session: "sess-1"}, Basis: "owner reviewed"}},
	}

	got := ResolveClearance(index, "O-001")
	if got.State != ClearanceCurrent {
		t.Fatalf("State = %q, want current", got.State)
	}
}

func mustCurrentTask(id string, checkID string, extra *Evidence) *TaskV2 {
	evidence := extra
	if evidence == nil {
		evidence = &Evidence{}
	}
	evidence.Freshness = &Freshness{State: FreshnessCurrent, Check: checkID, AssessedBy: Actor{Role: ActorRoleChecker, Session: "sess-1"}, Basis: "reviewed"}
	return &TaskV2{ID: id, Objective: "O-001", Status: ColumnInProgress, Stage: StageAudit, Evidence: evidence}
}

func TestResolveTaskStart_allowedWhenNoReplanAndDependenciesSatisfied(t *testing.T) {
	index := newV2TestIndex()
	mustCheck(index, "C-001", "T-002", CheckResultClear)
	index.Tasks["T-002"] = mustCurrentTask("T-002", "C-001", nil)
	index.Tasks["T-002"].Status = ColumnDone
	index.Tasks["T-001"] = &TaskV2{
		ID: "T-001", Objective: "O-001", Status: ColumnPlanned,
		DependsOn: []TaskDependencyV2{{Task: "T-002", Requires: TaskDependencyClear}},
	}

	got := ResolveTaskStart(index, "T-001")
	if !got.Allowed || len(got.Blockers) != 0 {
		t.Fatalf("ResolveTaskStart() = %+v, want allowed with no blockers", got)
	}
	if got.Actor != ActorRoleExecutor {
		t.Errorf("Actor = %q, want executor", got.Actor)
	}
}

func TestResolveTaskStart_blockedByReplan(t *testing.T) {
	index := newV2TestIndex()
	index.Tasks["T-001"] = &TaskV2{
		ID: "T-001", Objective: "O-001", Status: ColumnPlanned,
		Evidence: &Evidence{Replan: &Replan{Reason: "scope changed", RecordedBy: Actor{Role: ActorRoleOwner, Session: "sess-1"}}},
	}

	got := ResolveTaskStart(index, "T-001")
	if got.Allowed {
		t.Fatalf("ResolveTaskStart() Allowed = true, want false (replan flag set)")
	}
	if len(got.Blockers) != 1 || got.Blockers[0].Kind != GateBlockReplan {
		t.Fatalf("Blockers = %+v, want one GateBlockReplan", got.Blockers)
	}
}

func TestResolveTaskStart_namesEveryUnsatisfiedDependency(t *testing.T) {
	index := newV2TestIndex()
	index.Tasks["T-002"] = &TaskV2{ID: "T-002", Objective: "O-001", Status: ColumnInProgress}
	index.Tasks["T-003"] = &TaskV2{ID: "T-003", Objective: "O-001", Status: ColumnPlanned}
	index.Tasks["T-001"] = &TaskV2{
		ID: "T-001", Objective: "O-001", Status: ColumnPlanned,
		DependsOn: []TaskDependencyV2{
			{Task: "T-002", Requires: TaskDependencyClear},
			{Task: "T-003", Requires: TaskDependencyClear},
		},
	}

	got := ResolveTaskStart(index, "T-001")
	if got.Allowed {
		t.Fatalf("ResolveTaskStart() Allowed = true, want false")
	}
	if len(got.Blockers) != 2 {
		t.Fatalf("Blockers = %+v, want one blocker per unsatisfied dependency", got.Blockers)
	}
	targets := map[string]bool{}
	for _, blocker := range got.Blockers {
		if blocker.Kind != GateBlockDependency || blocker.Dependency == nil {
			t.Fatalf("Blocker = %+v, want GateBlockDependency with Dependency set", blocker)
		}
		targets[blocker.Dependency.Target] = true
	}
	if !targets["T-002"] || !targets["T-003"] {
		t.Fatalf("Blocker targets = %+v, want T-002 and T-003 both named", targets)
	}
}

func TestResolveTaskStart_blockedWhenOwningObjectiveDependencyUnsatisfied(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-002"] = &ObjectiveV2{ID: "O-002", Status: ColumnInProgress}
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Status: ColumnInProgress, DependsOn: []string{"O-002"}}
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnPlanned}

	got := ResolveTaskStart(index, "T-001")
	if got.Allowed {
		t.Fatalf("Allowed = true, want false (owning objective waits on O-002)")
	}
	if len(got.Blockers) != 1 {
		t.Fatalf("Blockers = %+v, want exactly one", got.Blockers)
	}
	blocker := got.Blockers[0]
	if blocker.Kind != GateBlockObjectiveDependency || blocker.ObjectiveDependency == nil {
		t.Fatalf("Blocker = %+v, want GateBlockObjectiveDependency with ObjectiveDependency set", blocker)
	}
	if blocker.ObjectiveDependency.Target != "O-002" {
		t.Errorf("ObjectiveDependency.Target = %q, want O-002", blocker.ObjectiveDependency.Target)
	}
	if !strings.Contains(blocker.Detail, "O-001") || !strings.Contains(blocker.Detail, "O-002") {
		t.Errorf("Detail = %q, want both the waiting objective O-001 and the dependency O-002 named", blocker.Detail)
	}
}

func TestResolveTaskStart_namesEveryUnsatisfiedObjectiveDependency(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-002"] = &ObjectiveV2{ID: "O-002", Status: ColumnInProgress}
	index.Objectives["O-003"] = &ObjectiveV2{ID: "O-003", Status: ColumnInProgress}
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Status: ColumnInProgress, DependsOn: []string{"O-002", "O-003"}}
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnPlanned}

	got := ResolveTaskStart(index, "T-001")
	if got.Allowed {
		t.Fatalf("Allowed = true, want false")
	}
	if len(got.Blockers) != 2 {
		t.Fatalf("Blockers = %+v, want one per unsatisfied objective dependency", got.Blockers)
	}
	targets := map[string]bool{}
	for _, blocker := range got.Blockers {
		if blocker.Kind != GateBlockObjectiveDependency || blocker.ObjectiveDependency == nil {
			t.Fatalf("Blocker = %+v, want GateBlockObjectiveDependency", blocker)
		}
		targets[blocker.ObjectiveDependency.Target] = true
	}
	if !targets["O-002"] || !targets["O-003"] {
		t.Fatalf("Blocker targets = %+v, want O-002 and O-003 both named", targets)
	}
}

func TestResolveTaskStart_allowedWhenOwningObjectiveDependencySatisfied(t *testing.T) {
	index := newV2TestIndex()
	mustObjectiveCheck(index, "C-001", "O-002", CheckResultClear)
	dependency := mustCurrentObjective("O-002", "C-001", nil)
	dependency.Status = ColumnDone
	index.Objectives["O-002"] = dependency
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Status: ColumnInProgress, DependsOn: []string{"O-002"}}
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnPlanned}

	got := ResolveTaskStart(index, "T-001")
	if !got.Allowed || len(got.Blockers) != 0 {
		t.Fatalf("ResolveTaskStart() = %+v, want allowed with no blockers", got)
	}
}

func TestResolveTaskStart_noObjectiveDependenciesStartsAsBefore(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Status: ColumnInProgress}
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnPlanned}

	got := ResolveTaskStart(index, "T-001")
	if !got.Allowed || len(got.Blockers) != 0 {
		t.Fatalf("ResolveTaskStart() = %+v, want allowed with no blockers (objective declares no dependencies)", got)
	}
}

func TestResolveTaskStart_unrelatedObjectiveUnaffectedByAnotherBlockedObjective(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-004"] = &ObjectiveV2{ID: "O-004", Status: ColumnInProgress}
	index.Objectives["O-003"] = &ObjectiveV2{ID: "O-003", Status: ColumnInProgress, DependsOn: []string{"O-004"}}
	index.Objectives["O-005"] = &ObjectiveV2{ID: "O-005", Status: ColumnInProgress}
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-005", Status: ColumnPlanned}

	got := ResolveTaskStart(index, "T-001")
	if !got.Allowed || len(got.Blockers) != 0 {
		t.Fatalf("ResolveTaskStart() = %+v, want allowed (an unrelated objective's blocked dependency does not spill over)", got)
	}
}

// TestResolveTaskAdvanceAndCompletion_unaffectedByObjectiveDependencyGate is
// the regression case required alongside the Objective dependency addition
// to ResolveTaskStart: only start consults Objective readiness.
func TestResolveTaskAdvanceAndCompletion_unaffectedByObjectiveDependencyGate(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-002"] = &ObjectiveV2{ID: "O-002", Status: ColumnInProgress} // unsatisfied dependency
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Status: ColumnInProgress, DependsOn: []string{"O-002"}}
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnInProgress, Stage: StageBuild}

	advance := ResolveTaskAdvance(index, "T-001")
	if !advance.Allowed {
		t.Fatalf("ResolveTaskAdvance() = %+v, want allowed (objective dependency readiness only gates start)", advance)
	}

	mustCheck(index, "C-001", "T-001", CheckResultClear)
	index.Tasks["T-001"] = mustCurrentTask("T-001", "C-001", nil)

	completion := ResolveTaskCompletion(index, "T-001")
	if !completion.Allowed {
		t.Fatalf("ResolveTaskCompletion() = %+v, want allowed (objective dependency readiness only gates start)", completion)
	}
}

func TestResolveTaskAdvance_blockedByReplan(t *testing.T) {
	index := newV2TestIndex()
	index.Tasks["T-001"] = &TaskV2{
		ID: "T-001", Objective: "O-001", Status: ColumnInProgress, Stage: StageBuild,
		Evidence: &Evidence{Replan: &Replan{Reason: "needs rework", RecordedBy: Actor{Role: ActorRoleOwner, Session: "sess-1"}}},
	}

	got := ResolveTaskAdvance(index, "T-001")
	if got.Allowed {
		t.Fatalf("ResolveTaskAdvance() Allowed = true, want false (replan flag set)")
	}
	if len(got.Blockers) != 1 || got.Blockers[0].Kind != GateBlockReplan {
		t.Fatalf("Blockers = %+v, want one GateBlockReplan", got.Blockers)
	}
}

func TestResolveTaskAdvance_movesThroughBuildAndTest(t *testing.T) {
	index := newV2TestIndex()
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnInProgress, Stage: StageBuild}
	if got := ResolveTaskAdvance(index, "T-001"); !got.Allowed || got.Actor != ActorRoleExecutor {
		t.Fatalf("ResolveTaskAdvance() from build = %+v, want allowed under executor authority", got)
	}

	index.Tasks["T-001"].Stage = StageTest
	if got := ResolveTaskAdvance(index, "T-001"); !got.Allowed {
		t.Fatalf("ResolveTaskAdvance() from test = %+v, want allowed", got)
	}
}

func TestResolveTaskAdvance_fromAuditDefersToCompletionDecision(t *testing.T) {
	index := newV2TestIndex()
	mustCheck(index, "C-001", "T-001", CheckResultClear)
	index.Tasks["T-001"] = mustCurrentTask("T-001", "C-001", nil)

	advance := ResolveTaskAdvance(index, "T-001")
	completion := ResolveTaskCompletion(index, "T-001")
	if !reflect.DeepEqual(advance, completion) {
		t.Fatalf("ResolveTaskAdvance() from audit = %+v, want it to equal ResolveTaskCompletion() = %+v", advance, completion)
	}
}

func TestResolveTaskCompletion_technicalTaskAllowedUnderCheckerAuthorityWhenCurrent(t *testing.T) {
	index := newV2TestIndex()
	mustCheck(index, "C-001", "T-001", CheckResultClear)
	index.Tasks["T-001"] = mustCurrentTask("T-001", "C-001", nil)

	got := ResolveTaskCompletion(index, "T-001")
	if !got.Allowed || got.AllowedByException {
		t.Fatalf("ResolveTaskCompletion() = %+v, want plain allowed", got)
	}
	if got.Actor != ActorRoleChecker {
		t.Errorf("Actor = %q, want checker", got.Actor)
	}
}

func TestResolveTaskCompletion_ownerValidationRequiredBlocksOnClearanceAlone(t *testing.T) {
	index := newV2TestIndex()
	mustCheck(index, "C-001", "T-001", CheckResultClear)
	index.Tasks["T-001"] = mustCurrentTask("T-001", "C-001", &Evidence{OwnerValidation: &OwnerValidation{Required: true}})

	got := ResolveTaskCompletion(index, "T-001")
	if got.Allowed {
		t.Fatalf("ResolveTaskCompletion() Allowed = true, want false (owner has not accepted)")
	}
	if len(got.Blockers) != 1 || got.Blockers[0].Kind != GateBlockOwnerAcceptance {
		t.Fatalf("Blockers = %+v, want one GateBlockOwnerAcceptance", got.Blockers)
	}
}

func TestResolveTaskCompletion_rejectsUnattributedOwnerAcceptance(t *testing.T) {
	index := newV2TestIndex()
	mustCheck(index, "C-001", "T-001", CheckResultClear)
	index.Tasks["T-001"] = mustCurrentTask("T-001", "C-001", &Evidence{
		OwnerValidation: &OwnerValidation{Required: true, AcceptedCheck: "C-001"},
	})

	decision := ResolveTaskCompletion(index, "T-001")
	if decision.Allowed || len(decision.Blockers) != 1 || decision.Blockers[0].Kind != GateBlockOwnerAcceptance {
		t.Fatalf("ResolveTaskCompletion() = %+v, want owner-acceptance blocker", decision)
	}
}

func TestResolveTaskCompletion_ownerValidationSatisfiedWhenAcceptedCurrentCheck(t *testing.T) {
	index := newV2TestIndex()
	mustCheck(index, "C-001", "T-001", CheckResultClear)
	index.Tasks["T-001"] = mustCurrentTask("T-001", "C-001", &Evidence{OwnerValidation: &OwnerValidation{Required: true, AcceptedCheck: "C-001", AcceptedBy: Actor{Role: ActorRoleOwner, Session: "owner-1"}}})

	got := ResolveTaskCompletion(index, "T-001")
	if !got.Allowed {
		t.Fatalf("ResolveTaskCompletion() Allowed = false, want true (owner accepted current check)")
	}
	if got.Actor != ActorRoleChecker {
		t.Errorf("Actor = %q, want checker (checker finalizes an acceptance that still applies)", got.Actor)
	}
}

func TestResolveTaskCompletion_acceptanceOfSupersededCheckDoesNotClose(t *testing.T) {
	index := newV2TestIndex()
	mustCheck(index, "C-001", "T-001", CheckResultClear)
	mustCheck(index, "C-002", "T-001", CheckResultClear) // supersedes C-001 as the latest
	index.Tasks["T-001"] = mustCurrentTask("T-001", "C-002", &Evidence{OwnerValidation: &OwnerValidation{Required: true, AcceptedCheck: "C-001", AcceptedBy: Actor{Role: ActorRoleOwner, Session: "owner-1"}}})

	got := ResolveTaskCompletion(index, "T-001")
	if got.Allowed {
		t.Fatalf("ResolveTaskCompletion() Allowed = true, want false (acceptance bound to a superseded check)")
	}
	if len(got.Blockers) != 1 || got.Blockers[0].Kind != GateBlockOwnerAcceptance {
		t.Fatalf("Blockers = %+v, want one GateBlockOwnerAcceptance", got.Blockers)
	}
}

func TestResolveTaskCompletion_eachClearanceStateBlocksWithADistinctReason(t *testing.T) {
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
			build:    func(index *V2Index) { mustCheck(index, "C-001", "T-001", CheckResultNeedsWork) },
			wantKind: GateBlockClearanceNeedsWork,
		},
		{
			name:     "unknown",
			build:    func(index *V2Index) { mustCheck(index, "C-001", "T-001", CheckResultClear) },
			evidence: &Evidence{Freshness: &Freshness{State: FreshnessUnknown, Check: "C-001", Basis: "not reassessed"}},
			wantKind: GateBlockClearanceUnknown,
		},
		{
			name:     "stale",
			build:    func(index *V2Index) { mustCheck(index, "C-001", "T-001", CheckResultClear) },
			evidence: &Evidence{Freshness: &Freshness{State: FreshnessStale, Check: "C-001", Basis: "code changed after the check"}},
			wantKind: GateBlockClearanceStale,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			index := newV2TestIndex()
			tc.build(index)
			index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnInProgress, Stage: StageAudit, Evidence: tc.evidence}

			got := ResolveTaskCompletion(index, "T-001")
			if got.Allowed {
				t.Fatalf("ResolveTaskCompletion() Allowed = true, want false")
			}
			if len(got.Blockers) != 1 || got.Blockers[0].Kind != tc.wantKind {
				t.Fatalf("Blockers = %+v, want one %q", got.Blockers, tc.wantKind)
			}
		})
	}
}

func TestResolveClearance_rejectsExecutorCheckAndFreshnessAsCurrent(t *testing.T) {
	index := newV2TestIndex()
	index.Checks["C-001"] = &CheckV2{
		ID:        "C-001",
		Scope:     CheckScope{Kind: CheckScopeTask, ID: "T-001"},
		Result:    CheckResultClear,
		CheckedBy: Actor{Role: ActorRoleExecutor, Session: "executor-1"},
	}
	index.LatestCheck["T-001"] = "C-001"
	index.Tasks["T-001"] = &TaskV2{
		ID: "T-001", Objective: "O-001", Status: ColumnInProgress, Stage: StageAudit,
		Evidence: &Evidence{Freshness: &Freshness{
			State: FreshnessCurrent, Check: "C-001",
			AssessedBy: Actor{Role: ActorRoleExecutor, Session: "executor-1"},
			Basis:      "executor self-report",
		}},
	}

	clearance := ResolveClearance(index, "T-001")
	if clearance.State == ClearanceCurrent {
		t.Fatalf("ResolveClearance() = %+v, must not treat executor evidence as current", clearance)
	}
	decision := ResolveTaskCompletion(index, "T-001")
	if decision.Allowed {
		t.Fatalf("ResolveTaskCompletion() = %+v, must block executor self-report", decision)
	}
	if len(decision.Blockers) != 1 || decision.Blockers[0].Kind != GateBlockCheckerAuthority {
		t.Fatalf("Blockers = %+v, want one checker-authority blocker", decision.Blockers)
	}
}

// A checker's CLEAR Check is sufficient on its own, so a "current" note from
// someone else adds nothing and takes nothing away.
func TestResolveTaskCompletion_checkerClearStandsWithoutCheckerFreshness(t *testing.T) {
	index := newV2TestIndex()
	check := mustCheck(index, "C-001", "T-001", CheckResultClear)
	check.CheckedBy = Actor{Role: ActorRoleChecker, Session: "checker-1"}
	index.Tasks["T-001"] = &TaskV2{
		ID: "T-001", Objective: "O-001", Status: ColumnInProgress, Stage: StageAudit,
		Evidence: &Evidence{Freshness: &Freshness{
			State: FreshnessCurrent, Check: "C-001",
			AssessedBy: Actor{Role: ActorRoleExecutor, Session: "executor-1"},
			Basis:      "executor self-report",
		}},
	}

	decision := ResolveTaskCompletion(index, "T-001")
	if !decision.Allowed {
		t.Fatalf("ResolveTaskCompletion() = %+v, want allowed on the checker's CLEAR Check", decision)
	}
}

func TestResolveTaskStart_requiresPlannedState(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status ColumnType
		stage  ProgressStage
	}{
		{name: "in progress", status: ColumnInProgress, stage: StageBuild},
		{name: "done", status: ColumnDone},
	} {
		t.Run(tc.name, func(t *testing.T) {
			index := newV2TestIndex()
			index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: tc.status, Stage: tc.stage}
			decision := ResolveTaskStart(index, "T-001")
			if decision.Allowed || len(decision.Blockers) != 1 || decision.Blockers[0].Kind != GateBlockInvalidState {
				t.Fatalf("ResolveTaskStart() = %+v, want invalid-state blocker", decision)
			}
		})
	}
}

func TestResolveTaskCompletion_requiresAuditStage(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status ColumnType
		stage  ProgressStage
	}{
		{name: "planned", status: ColumnPlanned},
		{name: "build", status: ColumnInProgress, stage: StageBuild},
		{name: "test", status: ColumnInProgress, stage: StageTest},
		{name: "done", status: ColumnDone},
	} {
		t.Run(tc.name, func(t *testing.T) {
			index := newV2TestIndex()
			index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: tc.status, Stage: tc.stage}
			decision := ResolveTaskCompletion(index, "T-001")
			if decision.Allowed || len(decision.Blockers) != 1 || decision.Blockers[0].Kind != GateBlockInvalidState {
				t.Fatalf("ResolveTaskCompletion() = %+v, want invalid-state blocker", decision)
			}
		})
	}
}

func TestResolveTaskAdvance_replanBlocksAuditToCompletion(t *testing.T) {
	index := newV2TestIndex()
	index.Tasks["T-001"] = &TaskV2{
		ID: "T-001", Objective: "O-001", Status: ColumnInProgress, Stage: StageAudit,
		Evidence: &Evidence{Replan: &Replan{Reason: "scope changed", RecordedBy: Actor{Role: ActorRolePlanner, Session: "planner-1"}}},
	}

	decision := ResolveTaskAdvance(index, "T-001")
	if decision.Allowed || len(decision.Blockers) != 1 || decision.Blockers[0].Kind != GateBlockReplan {
		t.Fatalf("ResolveTaskAdvance() = %+v, want replan blocker before completion", decision)
	}
}

func TestResolveTaskCompletion_allowedByExceptionWhenOtherwiseBlocked(t *testing.T) {
	index := newV2TestIndex()
	mustCheck(index, "C-001", "T-001", CheckResultNeedsWork)
	index.Tasks["T-001"] = &TaskV2{
		ID: "T-001", Objective: "O-001", Status: ColumnInProgress, Stage: StageAudit,
		Evidence: &Evidence{Exception: &Exception{
			Requirements: []string{"AC-3"},
			Reason:       "owner accepted the risk",
			Owner:        "ani",
			RecordedAt:   mustParseTime(t, "2026-09-01T00:00:00Z"),
			Check:        "C-001",
		}},
	}

	got := ResolveTaskCompletion(index, "T-001")
	if !got.Allowed || !got.AllowedByException {
		t.Fatalf("ResolveTaskCompletion() = %+v, want allowed-by-exception", got)
	}
	if got.Exception == nil || got.Exception.Reason != "owner accepted the risk" {
		t.Fatalf("Exception = %+v, want the recorded exception attached", got.Exception)
	}
	if got.Actor != ActorRoleOwner {
		t.Errorf("Actor = %q, want owner", got.Actor)
	}
}

func TestResolveTaskCompletion_exceptionDoesNotCarryToASupersedingCheck(t *testing.T) {
	index := newV2TestIndex()
	mustCheck(index, "C-001", "T-001", CheckResultNeedsWork)
	index.Tasks["T-001"] = &TaskV2{
		ID: "T-001", Objective: "O-001", Status: ColumnInProgress, Stage: StageAudit,
		Evidence: &Evidence{Exception: &Exception{
			Requirements: []string{"AC-3"},
			Reason:       "owner accepted the risk for C-001",
			Owner:        "ani",
			RecordedAt:   mustParseTime(t, "2026-09-01T00:00:00Z"),
			Check:        "C-001",
		}},
	}
	// A rerun records a new, superseding Check. The exception named C-001 only.
	mustCheck(index, "C-002", "T-001", CheckResultNeedsWork)

	got := ResolveTaskCompletion(index, "T-001")
	if got.Allowed {
		t.Fatalf("ResolveTaskCompletion() Allowed = true, want false (exception named the superseded check only)")
	}
	if got.AllowedByException {
		t.Fatalf("ResolveTaskCompletion() AllowedByException = true, want false")
	}
	if len(got.Blockers) != 1 || got.Blockers[0].Kind != GateBlockClearanceNeedsWork {
		t.Fatalf("Blockers = %+v, want one GateBlockClearanceNeedsWork", got.Blockers)
	}
}

func validTaskCheckWaiver(taskID string) *CheckWaiver {
	return &CheckWaiver{
		Task:       taskID,
		Reason:     "owner waived the local Check; the Full Objective Check will cover it",
		Actor:      Actor{Role: ActorRoleOwner, Session: "owner-1"},
		RecordedAt: time.Date(2026, 9, 20, 0, 0, 0, 0, time.UTC),
	}
}

func TestResolveTaskCompletion_allowedByWaiverWhenNoCheckWasEverRequested(t *testing.T) {
	index := newV2TestIndex()
	index.Tasks["T-001"] = &TaskV2{
		ID: "T-001", Objective: "O-001", Status: ColumnInProgress, Stage: StageAudit,
		Evidence: &Evidence{CheckWaiver: validTaskCheckWaiver("T-001")},
	}

	got := ResolveTaskCompletion(index, "T-001")
	if !got.Allowed || !got.AllowedByWaiver {
		t.Fatalf("ResolveTaskCompletion() = %+v, want allowed-by-waiver", got)
	}
	if got.Waiver == nil || got.Waiver.Reason == "" {
		t.Fatalf("Waiver = %+v, want the recorded waiver attached", got.Waiver)
	}
	if got.Actor != ActorRoleOwner {
		t.Errorf("Actor = %q, want owner", got.Actor)
	}
	if got.AllowedByException {
		t.Errorf("AllowedByException = true, want false (a waiver is not an exception)")
	}
}

// TestResolveTaskCompletion_waiverDoesNotApplyOnceACheckExists proves a
// waiver only substitutes for a Task Check that was never requested at all.
// Once a Check was actually run, its outcome — NEEDS WORK here — stands; the
// owner cannot retroactively wave away a result that already came back.
func TestResolveTaskCompletion_waiverDoesNotApplyOnceACheckExists(t *testing.T) {
	index := newV2TestIndex()
	mustCheck(index, "C-001", "T-001", CheckResultNeedsWork)
	index.Tasks["T-001"] = &TaskV2{
		ID: "T-001", Objective: "O-001", Status: ColumnInProgress, Stage: StageAudit,
		Evidence: &Evidence{CheckWaiver: validTaskCheckWaiver("T-001")},
	}

	got := ResolveTaskCompletion(index, "T-001")
	if got.Allowed {
		t.Fatalf("ResolveTaskCompletion() Allowed = true, want false (a waiver does not override a recorded NEEDS WORK)")
	}
	if got.AllowedByWaiver {
		t.Errorf("AllowedByWaiver = true, want false")
	}
	if len(got.Blockers) != 1 || got.Blockers[0].Kind != GateBlockClearanceNeedsWork {
		t.Fatalf("Blockers = %+v, want one GateBlockClearanceNeedsWork", got.Blockers)
	}
}

// TestResolveTaskCompletion_waiverNamingAnotherTaskDoesNotApply proves a
// waiver is bound to the Task that names it, mirroring how an exception is
// bound to the Check it names.
func TestResolveTaskCompletion_waiverNamingAnotherTaskDoesNotApply(t *testing.T) {
	index := newV2TestIndex()
	index.Tasks["T-001"] = &TaskV2{
		ID: "T-001", Objective: "O-001", Status: ColumnInProgress, Stage: StageAudit,
		Evidence: &Evidence{CheckWaiver: validTaskCheckWaiver("T-002")},
	}

	got := ResolveTaskCompletion(index, "T-001")
	if got.Allowed {
		t.Fatalf("ResolveTaskCompletion() Allowed = true, want false (waiver names a different task)")
	}
	if len(got.Blockers) != 1 || got.Blockers[0].Kind != GateBlockClearanceMissing {
		t.Fatalf("Blockers = %+v, want one GateBlockClearanceMissing", got.Blockers)
	}
}

// TestResolveObjectiveCompletion_unaffectedByAnOwnedTaskCheckWaiver proves
// the Full Objective Check stays mandatory: an owned Task closing by waiver
// still leaves the Objective needing its own current clearance, never
// inheriting the Task's waiver as if it were CLEAR.
func TestResolveObjectiveCompletion_unaffectedByAnOwnedTaskCheckWaiver(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001"}
	index.ObjectiveTasks["O-001"] = []string{"T-001"}
	index.Tasks["T-001"] = &TaskV2{
		ID: "T-001", Objective: "O-001", Status: ColumnDone,
		Evidence: &Evidence{CheckWaiver: validTaskCheckWaiver("T-001")},
	}

	got := ResolveObjectiveCompletion(index, "O-001")
	if got.Allowed {
		t.Fatalf("ResolveObjectiveCompletion() Allowed = true, want false (no Objective-scope Check recorded)")
	}
	if len(got.Blockers) != 1 || got.Blockers[0].Kind != GateBlockClearanceMissing {
		t.Fatalf("Blockers = %+v, want one GateBlockClearanceMissing", got.Blockers)
	}
}

// TestInspectTaskConsistency_waiverAllowedCompletionIsNotAContradiction
// proves a Task sitting in_progress/audit with a valid, applicable waiver is
// not reported as evidence-contradicts-status: like an exception, a waiver
// is a deliberate owner decision awaiting the owner's own status: done
// write, not an accidental hand-edit for doctor to flag.
func TestInspectTaskConsistency_waiverAllowedCompletionIsNotAContradiction(t *testing.T) {
	index := newV2TestIndex()
	index.Tasks["T-001"] = &TaskV2{
		ID: "T-001", Objective: "O-001", Status: ColumnInProgress, Stage: StageAudit,
		Evidence: &Evidence{CheckWaiver: validTaskCheckWaiver("T-001")},
	}

	got := InspectTaskConsistency(index)
	if len(got) != 0 {
		t.Fatalf("InspectTaskConsistency() = %+v, want no diagnostics for a waiver-allowed task awaiting owner closure", got)
	}
}

func TestInspectTaskConsistency_reportsEveryProblemNotOnlyTheFirst(t *testing.T) {
	index := newV2TestIndex()

	// T-001: done, but its latest check recorded NEEDS WORK.
	mustCheck(index, "C-001", "T-001", CheckResultNeedsWork)
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnDone}

	// T-002: owner accepted C-002, but a later check C-003 has since superseded it.
	mustCheck(index, "C-002", "T-002", CheckResultClear)
	mustCheck(index, "C-003", "T-002", CheckResultClear)
	index.Tasks["T-002"] = &TaskV2{
		ID: "T-002", Objective: "O-001", Status: ColumnInProgress, Stage: StageAudit,
		Evidence: &Evidence{
			Freshness:       &Freshness{State: FreshnessCurrent, Check: "C-003", AssessedBy: Actor{Role: ActorRoleChecker, Session: "sess-1"}, Basis: "rechecked"},
			OwnerValidation: &OwnerValidation{Required: true, AcceptedCheck: "C-002", AcceptedBy: Actor{Role: ActorRoleOwner, Session: "owner-1"}},
		},
	}

	// T-003: evidence clears completion, but status was never advanced to done.
	mustCheck(index, "C-004", "T-003", CheckResultClear)
	index.Tasks["T-003"] = mustCurrentTask("T-003", "C-004", nil)
	index.Tasks["T-003"].Status = ColumnInProgress

	got := InspectTaskConsistency(index)
	if len(got) != 3 {
		t.Fatalf("InspectTaskConsistency() = %+v, want 3 diagnostics (one per task)", got)
	}

	kinds := map[string]ConsistencyDiagnosticKind{}
	for _, diagnostic := range got {
		kinds[diagnostic.Task] = diagnostic.Kind
	}
	if kinds["T-001"] != ConsistencyDoneWithoutClearance {
		t.Errorf("T-001 kind = %q, want done_without_current_clearance", kinds["T-001"])
	}
	if kinds["T-002"] != ConsistencyAcceptanceSuperseded {
		t.Errorf("T-002 kind = %q, want acceptance_names_superseded_check", kinds["T-002"])
	}
	if kinds["T-003"] != ConsistencyEvidenceContradictsStatus {
		t.Errorf("T-003 kind = %q, want evidence_contradicts_status", kinds["T-003"])
	}
}

func TestInspectTaskConsistency_reportsNothingForAConsistentProject(t *testing.T) {
	index := newV2TestIndex()
	mustCheck(index, "C-001", "T-001", CheckResultClear)
	index.Tasks["T-001"] = mustCurrentTask("T-001", "C-001", nil)
	index.Tasks["T-001"].Status = ColumnDone

	got := InspectTaskConsistency(index)
	if len(got) != 0 {
		t.Fatalf("InspectTaskConsistency() = %+v, want no diagnostics for a consistent record", got)
	}
}

// TestGateDecisions_unaffectedByOpenIssuesOfEveryType proves no Issue type,
// severity, or open count changes any Task gate decision: ResolveTaskStart,
// ResolveTaskAdvance, and ResolveTaskCompletion return identical decisions
// with and without open Issues of every type present. Issues never gate
// directly — material blocking is expressed by a Check recording NEEDS WORK
// alone.
func TestGateDecisions_unaffectedByOpenIssuesOfEveryType(t *testing.T) {
	buildIndex := func() *V2Index {
		index := newV2TestIndex()
		mustCheck(index, "C-001", "T-002", CheckResultClear)
		index.Tasks["T-002"] = mustCurrentTask("T-002", "C-001", nil)
		index.Tasks["T-002"].Status = ColumnDone
		index.Tasks["T-001"] = &TaskV2{
			ID: "T-001", Objective: "O-001", Status: ColumnPlanned,
			DependsOn: []TaskDependencyV2{{Task: "T-002", Requires: TaskDependencyClear}},
		}
		return index
	}

	without := buildIndex()
	startWithout := ResolveTaskStart(without, "T-001")

	without.Tasks["T-001"].Status = ColumnInProgress
	without.Tasks["T-001"].Stage = StageBuild
	advanceWithout := ResolveTaskAdvance(without, "T-001")

	without.Tasks["T-001"].Stage = StageAudit
	completionWithout := ResolveTaskCompletion(without, "T-001")

	with := buildIndex()
	with.Issues = map[string]*IssueV2{}
	for i, issueType := range []IssueType{IssueTypeDefect, IssueTypeDrift, IssueTypeGuardrail, IssueTypeVerification, IssueTypeOther} {
		id := fmt.Sprintf("I-%03d", i+1)
		with.Issues[id] = &IssueV2{ID: id, Type: issueType, Status: IssueStatusOpen, Severity: "critical"}
	}
	with.TaskIssues = map[string][]string{"T-001": {"I-001", "I-002", "I-003", "I-004", "I-005"}}

	startWith := ResolveTaskStart(with, "T-001")
	if !reflect.DeepEqual(startWithout, startWith) {
		t.Fatalf("ResolveTaskStart() with open issues = %+v, want identical to without = %+v", startWith, startWithout)
	}

	with.Tasks["T-001"].Status = ColumnInProgress
	with.Tasks["T-001"].Stage = StageBuild
	advanceWith := ResolveTaskAdvance(with, "T-001")
	if !reflect.DeepEqual(advanceWithout, advanceWith) {
		t.Fatalf("ResolveTaskAdvance() with open issues = %+v, want identical to without = %+v", advanceWith, advanceWithout)
	}

	with.Tasks["T-001"].Stage = StageAudit
	completionWith := ResolveTaskCompletion(with, "T-001")
	if !reflect.DeepEqual(completionWithout, completionWith) {
		t.Fatalf("ResolveTaskCompletion() with open issues = %+v, want identical to without = %+v", completionWith, completionWithout)
	}
}

// TestResolveTaskCompletion_needsWorkCheckReferencingIssueBlocksThroughClearanceAlone
// proves a NEEDS WORK Check that names an Issue still blocks solely through
// E43's existing clearance rules, with no Issue-derived blocker kind
// introduced.
func TestResolveTaskCompletion_needsWorkCheckReferencingIssueBlocksThroughClearanceAlone(t *testing.T) {
	index := newV2TestIndex()
	check := mustCheck(index, "C-001", "T-001", CheckResultNeedsWork)
	check.Issues = []string{"I-001"}
	index.Issues = map[string]*IssueV2{
		"I-001": {ID: "I-001", Type: IssueTypeDefect, Status: IssueStatusOpen},
	}
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnInProgress, Stage: StageAudit}

	got := ResolveTaskCompletion(index, "T-001")
	if got.Allowed {
		t.Fatalf("ResolveTaskCompletion() Allowed = true, want false (NEEDS WORK check)")
	}
	if len(got.Blockers) != 1 || got.Blockers[0].Kind != GateBlockClearanceNeedsWork {
		t.Fatalf("Blockers = %+v, want exactly one GateBlockClearanceNeedsWork, no Issue-derived blocker kind", got.Blockers)
	}
}

func mustParseTime(t *testing.T, value string) time.Time {
	t.Helper()
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Fatalf("time.Parse(%q) error = %v", value, err)
	}
	return parsed
}
