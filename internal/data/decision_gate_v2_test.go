package data

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/testutil"
)

func checkerCarry(check string, applies bool, change string) DecisionCarry {
	return DecisionCarry{Check: check, Applies: applies, AssessedBy: Actor{Role: ActorRoleChecker, Session: "check-1"}, Reason: "assessed", MaterialChange: change}
}

func ownerCarry(check string) DecisionCarry {
	return DecisionCarry{Check: check, Applies: true, AssessedBy: Actor{Role: ActorRoleOwner, Session: "owner-2"}, Reason: "renewed"}
}

func TestAssessDecision(t *testing.T) {
	tests := []struct {
		name       string
		origin     string
		carried    []DecisionCarry
		latest     string
		want       DecisionApplicability
		wantChange string
	}{
		{"recorded against the latest check", "C-005", nil, "C-005", DecisionApplies, ""},
		{"no assessment at the latest check", "C-005", nil, "C-006", DecisionUnassessed, ""},
		{"assessment only at an earlier check", "C-005", []DecisionCarry{checkerCarry("C-006", true, "")}, "C-007", DecisionUnassessed, ""},
		{"checker says it applies", "C-005", []DecisionCarry{checkerCarry("C-006", true, "")}, "C-006", DecisionApplies, ""},
		{"checker says it changed", "C-005", []DecisionCarry{checkerCarry("C-006", false, "AC-2 changed")}, "C-006", DecisionChanged, "AC-2 changed"},
		{"owner renewal applies", "C-005", []DecisionCarry{ownerCarry("C-006")}, "C-006", DecisionApplies, ""},
		{"owner renewal outranks a checker change", "C-005", []DecisionCarry{checkerCarry("C-006", false, "AC-2 changed"), ownerCarry("C-006")}, "C-006", DecisionApplies, ""},
		{"no latest check", "C-005", []DecisionCarry{ownerCarry("C-006")}, "", DecisionUnassessed, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := assessDecision(tt.origin, tt.carried, tt.latest)
			if got.State != tt.want || got.MaterialChange != tt.wantChange {
				t.Errorf("assessDecision() = %+v, want state %s change %q", got, tt.want, tt.wantChange)
			}
		})
	}
}

func theShedExceptionEvidence(carried ...DecisionCarry) *Evidence {
	return &Evidence{
		OwnerValidation: &OwnerValidation{
			Required: true, AcceptedCheck: "C-005", AcceptedBy: Actor{Role: ActorRoleOwner, Session: "owner-1"},
			CarriedForward: carried,
		},
		Exception: &Exception{
			Requirements: []string{"TEST-08"}, Reason: "owner accepts unproven evidence", Owner: "ani", Check: "C-005",
			CarriedForward: carried,
		},
	}
}

// objectiveNeedingWork builds an Objective whose latest Check is a NEEDS WORK
// that lists unmet, with one done Task.
func objectiveNeedingWork(evidence *Evidence, unmet ...string) *V2Index {
	index := newV2TestIndex()
	mustObjectiveCheck(index, "C-005", "O-001", CheckResultNeedsWork)
	mustObjectiveCheck(index, "C-006", "O-001", CheckResultNeedsWork).Unmet = unmet
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnDone}
	index.ObjectiveTasks["O-001"] = []string{"T-001"}
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Status: ColumnInProgress, Evidence: evidence}
	return index
}

func blockerKinds(blockers []GateBlocker) []GateBlockKind {
	kinds := make([]GateBlockKind, 0, len(blockers))
	for _, blocker := range blockers {
		kinds = append(kinds, blocker.Kind)
	}
	return kinds
}

func TestResolveObjectiveCompletion_exceptionCarriedForwardGrantsCompletion(t *testing.T) {
	index := objectiveNeedingWork(theShedExceptionEvidence(checkerCarry("C-006", true, "")), "TEST-08")

	got := ResolveObjectiveCompletion(index, "O-001")
	if !got.Allowed || !got.AllowedByException || got.Actor != ActorRoleOwner {
		t.Fatalf("ResolveObjectiveCompletion() = %+v, want allowed by exception under owner authority", got)
	}
	if got.ExceptionCarry == nil || got.ExceptionCarry.Check != "C-006" {
		t.Errorf("ExceptionCarry = %+v, want the C-006 entry", got.ExceptionCarry)
	}
	if clearance := ResolveClearance(index, "O-001"); clearance.State != ClearanceNeedsWork {
		t.Errorf("clearance = %s, want needs_work while evidence is waived", clearance.State)
	}
}

func TestResolveObjectiveCompletion_exceptionAtItsOwnCheckHasNoCarry(t *testing.T) {
	index := objectiveNeedingWork(theShedExceptionEvidence(), "TEST-08")
	index.Objectives["O-001"].Evidence.Exception.Check = "C-006"

	got := ResolveObjectiveCompletion(index, "O-001")
	if !got.Allowed || !got.AllowedByException || got.ExceptionCarry != nil {
		t.Fatalf("ResolveObjectiveCompletion() = %+v, want allowed by exception with no carry entry (legacy shape)", got)
	}
}

func TestResolveObjectiveCompletion_unassessedExceptionNamesLatestCheck(t *testing.T) {
	index := objectiveNeedingWork(theShedExceptionEvidence(), "TEST-08")

	got := ResolveObjectiveCompletion(index, "O-001")
	if got.Allowed {
		t.Fatalf("ResolveObjectiveCompletion() = %+v, want blocked without a carry entry", got)
	}
	want := []GateBlockKind{GateBlockClearanceNeedsWork, GateBlockDecisionUnassessed}
	if !slices.Equal(blockerKinds(got.Blockers), want) {
		t.Fatalf("blocker kinds = %v, want %v", blockerKinds(got.Blockers), want)
	}
	if blocker := got.Blockers[1]; blocker.Decision != DecisionKindException || !strings.Contains(blocker.Detail, "C-006") {
		t.Errorf("blocker = %+v, want an exception blocker naming C-006", blocker)
	}
}

func TestResolveObjectiveCompletion_changedExceptionBlocksOnlyTheException(t *testing.T) {
	evidence := theShedExceptionEvidence()
	evidence.OwnerValidation.CarriedForward = []DecisionCarry{checkerCarry("C-006", true, "")}
	evidence.Exception.CarriedForward = []DecisionCarry{checkerCarry("C-006", false, "waiver now also covers STYLE-01")}
	index := objectiveNeedingWork(evidence, "TEST-08")

	got := ResolveObjectiveCompletion(index, "O-001")
	if got.Allowed {
		t.Fatalf("ResolveObjectiveCompletion() = %+v, want blocked", got)
	}
	last := got.Blockers[len(got.Blockers)-1]
	if last.Kind != GateBlockDecisionChanged || last.Decision != DecisionKindException || last.Change != "waiver now also covers STYLE-01" {
		t.Errorf("last blocker = %+v, want the exception's change named", last)
	}
	if !acceptanceApplies(index.Objectives["O-001"].Evidence, "C-006") {
		t.Error("acceptance stopped applying; a changed exception must not affect it")
	}
}

func TestResolveObjectiveCompletion_ownerRenewalOutranksCheckerChange(t *testing.T) {
	evidence := theShedExceptionEvidence()
	evidence.Exception.CarriedForward = []DecisionCarry{checkerCarry("C-006", false, "scope changed"), ownerCarry("C-006")}
	index := objectiveNeedingWork(evidence, "TEST-08")

	got := ResolveObjectiveCompletion(index, "O-001")
	if !got.Allowed || !got.AllowedByException || got.ExceptionCarry.AssessedBy.Role != ActorRoleOwner {
		t.Fatalf("ResolveObjectiveCompletion() = %+v, want allowed by the owner's renewal", got)
	}
}

func TestResolveObjectiveCompletion_exceptionMustCoverEveryUnmetRequirement(t *testing.T) {
	index := objectiveNeedingWork(theShedExceptionEvidence(checkerCarry("C-006", true, "")), "TEST-08", "DESIGN-01", "STYLE-04")

	got := ResolveObjectiveCompletion(index, "O-001")
	if got.Allowed {
		t.Fatalf("ResolveObjectiveCompletion() = %+v, want blocked by uncovered requirements", got)
	}
	last := got.Blockers[len(got.Blockers)-1]
	if last.Kind != GateBlockExceptionScope || !slices.Equal(last.Requirements, []string{"DESIGN-01", "STYLE-04"}) {
		t.Errorf("last blocker = %+v, want exception_scope naming DESIGN-01 and STYLE-04", last)
	}
}

func TestResolveObjectiveCompletion_exceptionWithoutUnmetListRelyOnCheckerAssessment(t *testing.T) {
	index := objectiveNeedingWork(theShedExceptionEvidence(checkerCarry("C-006", true, "")))
	if got := ResolveObjectiveCompletion(index, "O-001"); !got.Allowed || !got.AllowedByException {
		t.Fatalf("ResolveObjectiveCompletion() = %+v, want allowed when the Check lists no unmet requirements", got)
	}
}

func TestResolveObjectiveCompletion_carriedExceptionStillRefusesAnOpenIssueAndUnfinishedTask(t *testing.T) {
	index := objectiveNeedingWork(theShedExceptionEvidence(checkerCarry("C-006", true, "")), "TEST-08")
	index.Tasks["T-001"].Status = ColumnInProgress

	if got := ResolveObjectiveCompletion(index, "O-001"); got.Allowed {
		t.Fatalf("ResolveObjectiveCompletion() = %+v, want an unfinished Task to stay blocking", got)
	}
}

func TestResolveTaskCompletion_acceptanceCarriedForwardClosesUnderCheckerAuthority(t *testing.T) {
	index := newV2TestIndex()
	mustCheck(index, "C-001", "T-001", CheckResultClear)
	mustCheck(index, "C-002", "T-001", CheckResultClear)
	index.Tasks["T-001"] = mustCurrentTask("T-001", "C-002", &Evidence{OwnerValidation: &OwnerValidation{
		Required: true, AcceptedCheck: "C-001", AcceptedBy: Actor{Role: ActorRoleOwner, Session: "owner-1"},
		CarriedForward: []DecisionCarry{checkerCarry("C-002", true, "")},
	}})

	got := ResolveTaskCompletion(index, "T-001")
	if !got.Allowed || got.Actor != ActorRoleChecker || got.AllowedByException {
		t.Fatalf("ResolveTaskCompletion() = %+v, want allowed under checker authority", got)
	}
}

func TestResolveTaskCompletion_changedAcceptanceNamesTheChange(t *testing.T) {
	index := newV2TestIndex()
	mustCheck(index, "C-001", "T-001", CheckResultClear)
	mustCheck(index, "C-002", "T-001", CheckResultClear)
	index.Tasks["T-001"] = mustCurrentTask("T-001", "C-002", &Evidence{OwnerValidation: &OwnerValidation{
		Required: true, AcceptedCheck: "C-001", AcceptedBy: Actor{Role: ActorRoleOwner, Session: "owner-1"},
		CarriedForward: []DecisionCarry{checkerCarry("C-002", false, "AC-2 now rejects empty input")},
	}})

	got := ResolveTaskCompletion(index, "T-001")
	if got.Allowed || len(got.Blockers) != 1 {
		t.Fatalf("ResolveTaskCompletion() = %+v, want one blocker", got)
	}
	if blocker := got.Blockers[0]; blocker.Kind != GateBlockDecisionChanged || blocker.Decision != DecisionKindAcceptance || blocker.Change != "AC-2 now rejects empty input" {
		t.Errorf("blocker = %+v, want the acceptance change named", blocker)
	}

	index.Tasks["T-001"].Evidence.OwnerValidation.CarriedForward = append(index.Tasks["T-001"].Evidence.OwnerValidation.CarriedForward, ownerCarry("C-002"))
	if got := ResolveTaskCompletion(index, "T-001"); !got.Allowed {
		t.Errorf("ResolveTaskCompletion() = %+v, want the owner's renewal to restore the acceptance", got)
	}
}

func TestResolveTaskDependencyV2_acceptedFollowsApplicability(t *testing.T) {
	build := func(carried ...DecisionCarry) *V2Index {
		index := newV2TestIndex()
		mustCheck(index, "C-001", "T-002", CheckResultClear)
		mustCheck(index, "C-002", "T-002", CheckResultClear)
		index.Tasks["T-002"] = mustCurrentTask("T-002", "C-002", &Evidence{OwnerValidation: &OwnerValidation{
			Required: true, AcceptedCheck: "C-001", AcceptedBy: Actor{Role: ActorRoleOwner, Session: "owner-1"}, CarriedForward: carried,
		}})
		index.Tasks["T-002"].Status = ColumnDone
		return index
	}
	dep := TaskDependencyV2{Task: "T-002", Requires: TaskDependencyAccepted}

	if got := ResolveTaskDependencyV2(build(checkerCarry("C-002", true, "")), dep); !got.Satisfied {
		t.Errorf("carried acceptance: %+v, want satisfied", got)
	}
	if got := ResolveTaskDependencyV2(build(), dep); got.Satisfied || got.Block.Kind != DependencyBlockNotAccepted {
		t.Errorf("unassessed acceptance: %+v, want not_accepted", got)
	}
	if got := ResolveTaskDependencyV2(build(checkerCarry("C-002", false, "changed")), dep); got.Satisfied {
		t.Errorf("changed acceptance: %+v, want blocked", got)
	}
}

func TestInspectTaskConsistency_carriedAcceptanceIsNotSuperseded(t *testing.T) {
	build := func(carried ...DecisionCarry) *V2Index {
		index := newV2TestIndex()
		mustCheck(index, "C-001", "T-001", CheckResultClear)
		mustCheck(index, "C-002", "T-001", CheckResultClear)
		index.Tasks["T-001"] = mustCurrentTask("T-001", "C-002", &Evidence{OwnerValidation: &OwnerValidation{
			Required: true, AcceptedCheck: "C-001", AcceptedBy: Actor{Role: ActorRoleOwner, Session: "owner-1"}, CarriedForward: carried,
		}})
		index.Tasks["T-001"].Status = ColumnDone
		return index
	}

	if got := InspectTaskConsistency(build(checkerCarry("C-002", true, ""))); len(got) != 0 {
		t.Errorf("carried acceptance: diagnostics = %+v, want none", got)
	}
	if got := InspectTaskConsistency(build()); len(got) != 1 || got[0].Kind != ConsistencyAcceptanceSuperseded {
		t.Errorf("unassessed acceptance: diagnostics = %+v, want one acceptance_names_superseded_check", got)
	}
	got := InspectTaskConsistency(build(checkerCarry("C-002", false, "AC-2 changed")))
	if len(got) != 1 || !strings.Contains(got[0].Detail, "AC-2 changed") {
		t.Errorf("changed acceptance: diagnostics = %+v, want one naming the change", got)
	}
}

func TestObjectiveGates_carriedExceptionIsHonouredByDependencyAndConsistency(t *testing.T) {
	index := objectiveNeedingWork(theShedExceptionEvidence(checkerCarry("C-006", true, "")), "TEST-08")
	index.Objectives["O-001"].Status = ColumnDone

	if got := ResolveObjectiveDependency(index, "O-001"); got.Block == nil || got.Block.Kind != ObjectiveDependencyBlockClearedByException {
		t.Errorf("ResolveObjectiveDependency() = %+v, want cleared_by_exception", got)
	}
	if got := InspectObjectiveConsistency(index); len(got) != 0 {
		t.Errorf("InspectObjectiveConsistency() = %+v, want none for a done Objective closed by a carried exception", got)
	}

	index.Objectives["O-001"].Evidence.Exception.CarriedForward = nil
	if got := ResolveObjectiveDependency(index, "O-001"); got.Block == nil || got.Block.Kind != ObjectiveDependencyBlockNotCleared {
		t.Errorf("ResolveObjectiveDependency() = %+v, want not_cleared without the carry entry", got)
	}
	if got := InspectObjectiveConsistency(index); len(got) != 1 || got[0].Kind != ObjectiveConsistencyDoneWithoutClearance {
		t.Errorf("InspectObjectiveConsistency() = %+v, want done_without_current_clearance", got)
	}
}

// theShedProject writes the TheShed O-002 shape into a temporary directory:
// done Tasks, an owner acceptance and an evidence exception recorded against
// C-005, C-005 NEEDS WORK for Design drift plus waived evidence, and C-006
// superseding it for the waived evidence only. carry is the carried_forward
// list written on both decisions.
func theShedProject(t *testing.T, carry string) string {
	t.Helper()
	root := t.TempDir()
	objective := "---\nid: O-002\ntitle: \"Shed\"\nstatus: in_progress\n" +
		"owner_validation: {required: true, accepted_check: C-005, accepted_by: {role: owner, session: owner-1}" + carry + "}\n" +
		"exception: {requirements: [TEST-08], reason: owner accepts unproven evidence, owner: ani, recorded_at: '2026-10-09T00:00:00Z', check: C-005" + carry + "}\n" +
		"---\n\n# Shed\n"
	testutil.WriteFile(t, filepath.Join(root, v2ObjectivesDirName, "O-002-shed", v2ObjectiveFileName), objective)
	for _, id := range []string{"T-001", "T-002"} {
		testutil.WriteFile(t, filepath.Join(root, v2ObjectivesDirName, "O-002-shed", v2TasksDirName, id+"-task.md"),
			"---\nid: "+id+"\ntitle: \"Task\"\nobjective: O-002\nplanned_by: {role: planner, session: planning-fixture}\nstatus: done\n---\n\n# Task\n")
	}
	check := func(id, supersedes, unmet string) {
		content := "---\nid: " + id + "\nscope: {kind: objective, id: O-002}\nresult: NEEDS WORK\nchecked_by: {role: checker, session: check-" + id + "}\n" +
			"executed_session: build-fixture\nchecked_at: '2026-10-09T00:00:00Z'\nunmet: " + unmet + "\n" + supersedes + "---\n\n# Check\n"
		testutil.WriteFile(t, filepath.Join(root, v2ChecksDirName, id+"-check.md"), content)
	}
	check("C-005", "", "[DESIGN-01, TEST-08]")
	check("C-006", "supersedes: C-005\n", "[TEST-08]")
	return root
}

const theShedCarry = ", carried_forward: [{check: C-006, applies: true, assessed_by: {role: checker, session: check-C-006}, assessed_at: '2026-10-10T00:00:00Z', reason: only Design text changed}]"

func TestTheShedO002_decisionsCarriedForwardAllowCompletionByException(t *testing.T) {
	index := mustLoadV2Index(t, theShedProject(t, theShedCarry))

	got := ResolveObjectiveCompletion(index, "O-002")
	if !got.Allowed || !got.AllowedByException || got.Actor != ActorRoleOwner {
		t.Fatalf("ResolveObjectiveCompletion() = %+v, want completion allowed by exception with owner authority", got)
	}
	if got.ExceptionCarry == nil || got.ExceptionCarry.Check != "C-006" {
		t.Errorf("ExceptionCarry = %+v, want the C-006 entry", got.ExceptionCarry)
	}
	if clearance := ResolveClearance(index, "O-002"); clearance.State != ClearanceNeedsWork || clearance.Check != "C-006" {
		t.Errorf("clearance = %+v, want NEEDS WORK on C-006, never current", clearance)
	}
	if !acceptanceApplies(index.Objectives["O-002"].Evidence, "C-006") {
		t.Error("owner acceptance does not apply at C-006, want carried forward")
	}
}

// TestTheShedO002_withoutCarryEntriesReportsUnassessed is the failing scenario
// before this change: the same records only produced a plain NEEDS WORK block.
func TestTheShedO002_withoutCarryEntriesReportsUnassessed(t *testing.T) {
	index := mustLoadV2Index(t, theShedProject(t, ""))

	got := ResolveObjectiveCompletion(index, "O-002")
	if got.Allowed {
		t.Fatalf("ResolveObjectiveCompletion() = %+v, want blocked without carry entries", got)
	}
	last := got.Blockers[len(got.Blockers)-1]
	if last.Kind != GateBlockDecisionUnassessed || !strings.Contains(last.Detail, "C-006") {
		t.Errorf("last blocker = %+v, want unassessed naming C-006", last)
	}
}

func TestResolveObjectiveIntegrationRung_uncoveredRequirementAsksTheOwner(t *testing.T) {
	index := objectiveNeedingWork(theShedExceptionEvidence(checkerCarry("C-006", true, "")), "TEST-08", "DESIGN-02")

	got, ok := resolveObjectiveIntegrationRung(index, "O-001")
	if !ok || got.Kind != NextOwnerValidationRequired {
		t.Fatalf("resolveObjectiveIntegrationRung() = %+v, %v, want the owner-facing Accept rung", got, ok)
	}
}
