package data

import (
	"errors"
	"testing"
)

func newV2TestIndex() *V2Index {
	return &V2Index{
		Objectives:     map[string]*ObjectiveV2{},
		Tasks:          map[string]*TaskV2{},
		Checks:         map[string]*CheckV2{},
		ObjectiveTasks: map[string][]string{},
		ScopeChecks:    map[string][]string{},
		LatestCheck:    map[string]string{},
	}
}

func TestValidateV2ReferenceGraphs_valid(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Title: "First", Status: ColumnPlanned, Source: V2SourceDocument{Path: "objectives/O-001-a/Objective.md"}}
	index.Objectives["O-002"] = &ObjectiveV2{ID: "O-002", Title: "Second", Status: ColumnPlanned, DependsOn: []string{"O-001"}, Source: V2SourceDocument{Path: "objectives/O-002-b/Objective.md"}}
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Title: "First task", Objective: "O-001", Status: ColumnPlanned, Source: V2SourceDocument{Path: "objectives/O-001-a/tasks/T-001-first.md"}}
	index.Tasks["T-002"] = &TaskV2{ID: "T-002", Title: "Second task", Objective: "O-001", Status: ColumnPlanned, DependsOn: []TaskDependencyV2{{Task: "T-001", Requires: TaskDependencyClear}}, Source: V2SourceDocument{Path: "objectives/O-001-a/tasks/T-002-second.md"}}

	if err := ValidateV2ReferenceGraphs(index); err != nil {
		t.Fatalf("ValidateV2ReferenceGraphs() error = %v, want nil", err)
	}
}

func TestValidateV2ReferenceGraphs_taskSelfDependency(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Title: "First", Status: ColumnPlanned, Source: V2SourceDocument{Path: "objectives/O-001-a/Objective.md"}}
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Title: "Self", Objective: "O-001", Status: ColumnPlanned, DependsOn: []TaskDependencyV2{{Task: "T-001", Requires: TaskDependencyClear}}, Source: V2SourceDocument{Path: "objectives/O-001-a/tasks/T-001-self.md"}}

	err := ValidateV2ReferenceGraphs(index)
	if !errors.Is(err, ErrV2SelfDependency) {
		t.Fatalf("ValidateV2ReferenceGraphs() error = %v, want ErrV2SelfDependency", err)
	}
}

func TestValidateV2ReferenceGraphs_taskMissingTarget(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Title: "First", Status: ColumnPlanned, Source: V2SourceDocument{Path: "objectives/O-001-a/Objective.md"}}
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Title: "Orphan dep", Objective: "O-001", Status: ColumnPlanned, DependsOn: []TaskDependencyV2{{Task: "T-999", Requires: TaskDependencyClear}}, Source: V2SourceDocument{Path: "objectives/O-001-a/tasks/T-001-orphan.md"}}

	err := ValidateV2ReferenceGraphs(index)
	if !errors.Is(err, ErrV2MissingDependencyTarget) {
		t.Fatalf("ValidateV2ReferenceGraphs() error = %v, want ErrV2MissingDependencyTarget", err)
	}
}

func TestValidateV2ReferenceGraphs_taskCycle(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Title: "First", Status: ColumnPlanned, Source: V2SourceDocument{Path: "objectives/O-001-a/Objective.md"}}
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Title: "A", Objective: "O-001", Status: ColumnPlanned, DependsOn: []TaskDependencyV2{{Task: "T-002", Requires: TaskDependencyClear}}, Source: V2SourceDocument{Path: "objectives/O-001-a/tasks/T-001-a.md"}}
	index.Tasks["T-002"] = &TaskV2{ID: "T-002", Title: "B", Objective: "O-001", Status: ColumnPlanned, DependsOn: []TaskDependencyV2{{Task: "T-003", Requires: TaskDependencyClear}}, Source: V2SourceDocument{Path: "objectives/O-001-a/tasks/T-002-b.md"}}
	index.Tasks["T-003"] = &TaskV2{ID: "T-003", Title: "C", Objective: "O-001", Status: ColumnPlanned, DependsOn: []TaskDependencyV2{{Task: "T-001", Requires: TaskDependencyClear}}, Source: V2SourceDocument{Path: "objectives/O-001-a/tasks/T-003-c.md"}}

	err := ValidateV2ReferenceGraphs(index)
	if !errors.Is(err, ErrV2DependencyCycle) {
		t.Fatalf("ValidateV2ReferenceGraphs() error = %v, want ErrV2DependencyCycle", err)
	}
}

func TestValidateV2ReferenceGraphs_objectiveSelfDependency(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Title: "Self", Status: ColumnPlanned, DependsOn: []string{"O-001"}, Source: V2SourceDocument{Path: "objectives/O-001-a/Objective.md"}}

	err := ValidateV2ReferenceGraphs(index)
	if !errors.Is(err, ErrV2SelfDependency) {
		t.Fatalf("ValidateV2ReferenceGraphs() error = %v, want ErrV2SelfDependency", err)
	}
}

func TestValidateV2ReferenceGraphs_objectiveMissingTarget(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Title: "Orphan dep", Status: ColumnPlanned, DependsOn: []string{"O-999"}, Source: V2SourceDocument{Path: "objectives/O-001-a/Objective.md"}}

	err := ValidateV2ReferenceGraphs(index)
	if !errors.Is(err, ErrV2MissingDependencyTarget) {
		t.Fatalf("ValidateV2ReferenceGraphs() error = %v, want ErrV2MissingDependencyTarget", err)
	}
}

func TestValidateV2ReferenceGraphs_objectiveCycle(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Title: "A", Status: ColumnPlanned, DependsOn: []string{"O-002"}, Source: V2SourceDocument{Path: "objectives/O-001-a/Objective.md"}}
	index.Objectives["O-002"] = &ObjectiveV2{ID: "O-002", Title: "B", Status: ColumnPlanned, DependsOn: []string{"O-001"}, Source: V2SourceDocument{Path: "objectives/O-002-b/Objective.md"}}

	err := ValidateV2ReferenceGraphs(index)
	if !errors.Is(err, ErrV2DependencyCycle) {
		t.Fatalf("ValidateV2ReferenceGraphs() error = %v, want ErrV2DependencyCycle", err)
	}
}

// TestValidateV2ReferenceGraphs_taskAndObjectiveCyclesNotConflated proves a
// Task graph cycle is reported even when the Objective graph is
// simultaneously acyclic, and vice versa: the two graphs are validated
// independently rather than sharing one combined graph.
func TestValidateV2ReferenceGraphs_taskAndObjectiveCyclesNotConflated(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Title: "A", Status: ColumnPlanned, Source: V2SourceDocument{Path: "objectives/O-001-a/Objective.md"}}
	index.Objectives["O-002"] = &ObjectiveV2{ID: "O-002", Title: "B", Status: ColumnPlanned, Source: V2SourceDocument{Path: "objectives/O-002-b/Objective.md"}}
	// Task IDs collide with Objective IDs in shape only; the graphs must stay
	// separate keyspaces (T-### vs O-###), so this is not a real risk in
	// practice, but the independent-cycle guarantee still needs a task cycle
	// under acyclic objectives to be provable.
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Title: "A", Objective: "O-001", Status: ColumnPlanned, DependsOn: []TaskDependencyV2{{Task: "T-002", Requires: TaskDependencyClear}}, Source: V2SourceDocument{Path: "objectives/O-001-a/tasks/T-001-a.md"}}
	index.Tasks["T-002"] = &TaskV2{ID: "T-002", Title: "B", Objective: "O-001", Status: ColumnPlanned, DependsOn: []TaskDependencyV2{{Task: "T-001", Requires: TaskDependencyClear}}, Source: V2SourceDocument{Path: "objectives/O-001-a/tasks/T-002-b.md"}}

	err := ValidateV2ReferenceGraphs(index)
	if !errors.Is(err, ErrV2DependencyCycle) {
		t.Fatalf("ValidateV2ReferenceGraphs() error = %v, want ErrV2DependencyCycle from the task graph", err)
	}
}

func TestResolveTaskDependencyV2_clearSatisfiedWhenDoneAndCurrent(t *testing.T) {
	index := newV2TestIndex()
	mustCheck(index, "C-001", "T-001", CheckResultClear)
	index.Tasks["T-001"] = &TaskV2{
		ID: "T-001", Objective: "O-001", Status: ColumnDone,
		Evidence: &Evidence{Freshness: &Freshness{State: FreshnessCurrent, Check: "C-001", AssessedBy: Actor{Role: ActorRoleChecker, Session: "sess-1"}, Basis: "checked"}},
	}

	got := ResolveTaskDependencyV2(index, TaskDependencyV2{Task: "T-001", Requires: TaskDependencyClear})
	if !got.Satisfied || got.Block != nil {
		t.Fatalf("ResolveTaskDependencyV2() = %+v, want satisfied with no block", got)
	}
}

// TestResolveTaskDependencyV2_clearSatisfiedByCheckWaiver proves an owner
// Task-check waiver satisfies a requires: clear dependency exactly the way
// current clearance does — the waiver is the owner's own completion
// decision, so a downstream Task's default dependency does not re-demand an
// independent Check the owner already chose to skip.
func TestResolveTaskDependencyV2_clearSatisfiedByCheckWaiver(t *testing.T) {
	index := newV2TestIndex()
	index.Tasks["T-001"] = &TaskV2{
		ID: "T-001", Objective: "O-001", Status: ColumnDone,
		Evidence: &Evidence{CheckWaiver: validTaskCheckWaiver("T-001")},
	}

	got := ResolveTaskDependencyV2(index, TaskDependencyV2{Task: "T-001", Requires: TaskDependencyClear})
	if !got.Satisfied || got.Block != nil {
		t.Fatalf("ResolveTaskDependencyV2() = %+v, want satisfied by the recorded waiver", got)
	}
}

// TestResolveTaskDependencyV2_acceptedNeverSatisfiedByCheckWaiver proves the
// waiver's stand-in for "clear" does not extend to requires: accepted —
// there is no Check for the owner to have accepted, so that stricter
// dependency level stays blocked.
func TestResolveTaskDependencyV2_acceptedNeverSatisfiedByCheckWaiver(t *testing.T) {
	index := newV2TestIndex()
	index.Tasks["T-001"] = &TaskV2{
		ID: "T-001", Objective: "O-001", Status: ColumnDone,
		Evidence: &Evidence{CheckWaiver: validTaskCheckWaiver("T-001")},
	}

	got := ResolveTaskDependencyV2(index, TaskDependencyV2{Task: "T-001", Requires: TaskDependencyAccepted})
	if got.Satisfied {
		t.Fatalf("ResolveTaskDependencyV2() satisfied = true, want false: a waiver is not an accepted Check")
	}
	if got.Block == nil || got.Block.Kind != DependencyBlockNotAccepted {
		t.Fatalf("Block = %+v, want {Kind: not_accepted}", got.Block)
	}
}

func TestResolveTaskDependencyV2_notDoneUnsatisfied(t *testing.T) {
	index := newV2TestIndex()
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnInProgress}

	got := ResolveTaskDependencyV2(index, TaskDependencyV2{Task: "T-001", Requires: TaskDependencyClear})
	if got.Satisfied {
		t.Fatalf("ResolveTaskDependencyV2() satisfied = true, want false")
	}
	if got.Block == nil || got.Block.Kind != DependencyBlockNotDone || got.Block.Target != "T-001" {
		t.Fatalf("Block = %+v, want {Target: T-001, Kind: not_done}", got.Block)
	}
}

func TestResolveTaskDependencyV2_doneButNotClearedUnsatisfied(t *testing.T) {
	cases := []struct {
		name          string
		buildEvidence func(index *V2Index)
		evidence      *Evidence
		wantClearance ClearanceState
	}{
		{
			name: "needs_work",
			buildEvidence: func(index *V2Index) {
				mustCheck(index, "C-001", "T-001", CheckResultNeedsWork)
			},
			wantClearance: ClearanceNeedsWork,
		},
		{
			name: "unknown",
			buildEvidence: func(index *V2Index) {
				mustCheck(index, "C-001", "T-001", CheckResultClear)
			},
			wantClearance: ClearanceUnknown,
		},
		{
			name: "stale",
			buildEvidence: func(index *V2Index) {
				mustCheck(index, "C-001", "T-001", CheckResultClear)
				mustCheck(index, "C-002", "T-001", CheckResultClear)
			},
			wantClearance: ClearanceStale,
			evidence:      &Evidence{Freshness: &Freshness{State: FreshnessCurrent, Check: "C-001", Basis: "names the superseded check"}},
		},
		{
			name:          "missing",
			buildEvidence: func(index *V2Index) {},
			wantClearance: ClearanceMissing,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			index := newV2TestIndex()
			tc.buildEvidence(index)
			index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnDone, Evidence: tc.evidence}

			got := ResolveTaskDependencyV2(index, TaskDependencyV2{Task: "T-001", Requires: TaskDependencyClear})
			if got.Satisfied {
				t.Fatalf("ResolveTaskDependencyV2() satisfied = true, want false")
			}
			if got.Block == nil || got.Block.Kind != DependencyBlockNotCleared {
				t.Fatalf("Block = %+v, want kind not_cleared", got.Block)
			}
			if got.Block.Clearance != tc.wantClearance {
				t.Errorf("Block.Clearance = %q, want %q", got.Block.Clearance, tc.wantClearance)
			}
		})
	}
}

func TestResolveTaskDependencyV2_acceptedSatisfiedWhenOwnerAcceptedCurrentCheck(t *testing.T) {
	index := newV2TestIndex()
	mustCheck(index, "C-001", "T-001", CheckResultClear)
	index.Tasks["T-001"] = &TaskV2{
		ID: "T-001", Objective: "O-001", Status: ColumnDone,
		Evidence: &Evidence{
			Freshness:       &Freshness{State: FreshnessCurrent, Check: "C-001", AssessedBy: Actor{Role: ActorRoleChecker, Session: "sess-1"}, Basis: "checked"},
			OwnerValidation: &OwnerValidation{Required: true, AcceptedCheck: "C-001", AcceptedBy: Actor{Role: ActorRoleOwner, Session: "owner-1"}},
		},
	}

	got := ResolveTaskDependencyV2(index, TaskDependencyV2{Task: "T-001", Requires: TaskDependencyAccepted})
	if !got.Satisfied || got.Block != nil {
		t.Fatalf("ResolveTaskDependencyV2() = %+v, want satisfied with no block", got)
	}
}

func TestResolveTaskDependencyV2_acceptedUnsatisfiedWithNoAcceptance(t *testing.T) {
	index := newV2TestIndex()
	mustCheck(index, "C-001", "T-001", CheckResultClear)
	index.Tasks["T-001"] = &TaskV2{
		ID: "T-001", Objective: "O-001", Status: ColumnDone,
		Evidence: &Evidence{Freshness: &Freshness{State: FreshnessCurrent, Check: "C-001", AssessedBy: Actor{Role: ActorRoleChecker, Session: "sess-1"}, Basis: "checked"}},
	}

	got := ResolveTaskDependencyV2(index, TaskDependencyV2{Task: "T-001", Requires: TaskDependencyAccepted})
	if got.Satisfied {
		t.Fatalf("ResolveTaskDependencyV2() satisfied = true, want false")
	}
	if got.Block == nil || got.Block.Kind != DependencyBlockNotAccepted {
		t.Fatalf("Block = %+v, want kind not_accepted", got.Block)
	}
	if got.Block.Clearance != ClearanceCurrent {
		t.Errorf("Block.Clearance = %q, want current (blocked on acceptance, not clearance)", got.Block.Clearance)
	}
}

// TestResolveTaskDependencyV2_acceptedUnsatisfiedWhenAcceptanceBoundToSupersededCheck
// proves that owner acceptance of an earlier Check does not satisfy
// requires: accepted once a newer Check supersedes it, even though the
// dependency Task's own clearance can still be current via a fresh
// freshness assessment naming the new Check.
func TestResolveTaskDependencyV2_acceptedUnsatisfiedWhenAcceptanceBoundToSupersededCheck(t *testing.T) {
	index := newV2TestIndex()
	mustCheck(index, "C-001", "T-001", CheckResultClear)
	mustCheck(index, "C-002", "T-001", CheckResultClear)
	index.Tasks["T-001"] = &TaskV2{
		ID: "T-001", Objective: "O-001", Status: ColumnDone,
		Evidence: &Evidence{
			Freshness:       &Freshness{State: FreshnessCurrent, Check: "C-002", AssessedBy: Actor{Role: ActorRoleChecker, Session: "sess-1"}, Basis: "rechecked"},
			OwnerValidation: &OwnerValidation{Required: true, AcceptedCheck: "C-001", AcceptedBy: Actor{Role: ActorRoleOwner, Session: "owner-1"}},
		},
	}

	got := ResolveTaskDependencyV2(index, TaskDependencyV2{Task: "T-001", Requires: TaskDependencyAccepted})
	if got.Satisfied {
		t.Fatalf("ResolveTaskDependencyV2() satisfied = true, want false (acceptance bound to superseded C-001, latest is C-002)")
	}
	if got.Block == nil || got.Block.Kind != DependencyBlockNotAccepted {
		t.Fatalf("Block = %+v, want kind not_accepted", got.Block)
	}
}

// TestResolveTaskDependencyV2_chainAcrossSeveralTasks proves each dependency
// in a chain is resolved independently against the live index rather than a
// cached judgement: T-003 depends on T-002, which depends on T-001. T-001 is
// done and current; T-002 is done but only unknown (no freshness assessment),
// so T-002's own dependency on T-001 is satisfied while T-003's dependency on
// T-002 is not.
func TestResolveTaskDependencyV2_chainAcrossSeveralTasks(t *testing.T) {
	index := newV2TestIndex()
	mustCheck(index, "C-001", "T-001", CheckResultClear)
	mustCheck(index, "C-002", "T-002", CheckResultClear)
	index.Tasks["T-001"] = &TaskV2{
		ID: "T-001", Objective: "O-001", Status: ColumnDone,
		Evidence: &Evidence{Freshness: &Freshness{State: FreshnessCurrent, Check: "C-001", AssessedBy: Actor{Role: ActorRoleChecker, Session: "sess-1"}, Basis: "checked"}},
	}
	index.Tasks["T-002"] = &TaskV2{
		ID: "T-002", Objective: "O-001", Status: ColumnDone,
		DependsOn: []TaskDependencyV2{{Task: "T-001", Requires: TaskDependencyClear}},
	}
	index.Tasks["T-003"] = &TaskV2{
		ID: "T-003", Objective: "O-001", Status: ColumnDone,
		DependsOn: []TaskDependencyV2{{Task: "T-002", Requires: TaskDependencyClear}},
	}

	gotT2onT1 := ResolveTaskDependencyV2(index, index.Tasks["T-002"].DependsOn[0])
	if !gotT2onT1.Satisfied {
		t.Fatalf("T-002 -> T-001 = %+v, want satisfied", gotT2onT1)
	}

	gotT3onT2 := ResolveTaskDependencyV2(index, index.Tasks["T-003"].DependsOn[0])
	if gotT3onT2.Satisfied {
		t.Fatalf("T-003 -> T-002 satisfied = true, want false (T-002 has no freshness assessment)")
	}
	if gotT3onT2.Block == nil || gotT3onT2.Block.Kind != DependencyBlockNotCleared || gotT3onT2.Block.Clearance != ClearanceUnknown {
		t.Fatalf("Block = %+v, want {Kind: not_cleared, Clearance: unknown}", gotT3onT2.Block)
	}
}
