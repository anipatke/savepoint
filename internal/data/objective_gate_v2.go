package data

import (
	"fmt"
	"maps"
	"slices"
)

// ResolveObjectiveCompletion decides whether objectiveID may be treated as
// complete. Completion requires every Task the Objective owns — membership
// read from index.ObjectiveTasks, never a list maintained on the Objective
// itself — to be done, and the Objective's own integration clearance,
// resolved by ResolveClearance over its scope.kind: objective Checks, to be
// current. An Objective declaring owner_validation.required additionally
// needs owner acceptance of that same current Check; acceptance bound to a
// superseded Check does not satisfy it. When clearance would otherwise block
// completion, a recorded exception naming the Objective's current latest
// Check grants completion by exception under owner authority instead —
// never as a CLEAR result or current clearance. Material Issues linked to
// the current Objective Check must be resolved; an Objective exception does
// not silently accept an open Issue. An unfinished owned Task is not
// excusable by exception: cross-Task repair goes back through Tasks, and
// no Objective Check ever closes a Task.
func ResolveObjectiveCompletion(index *V2Index, objectiveID string) GateDecision {
	objective, ok := index.Objectives[objectiveID]
	if !ok {
		return GateDecision{}
	}

	var incomplete []GateBlocker
	for _, taskID := range index.ObjectiveTasks[objectiveID] {
		task := index.Tasks[taskID]
		if task.Status != ColumnDone {
			incomplete = append(incomplete, GateBlocker{
				Kind:   GateBlockInvalidState,
				Detail: fmt.Sprintf("task %s is not done (status %q)", taskID, task.Status),
			})
		}
	}
	if len(incomplete) > 0 {
		return GateDecision{Blockers: incomplete}
	}

	clearance := ResolveClearance(index, objectiveID)
	var blockers []GateBlocker

	switch clearance.State {
	case ClearanceMissing:
		blockers = append(blockers, GateBlocker{Kind: GateBlockClearanceMissing, Detail: "no recorded check"})
	case ClearanceNeedsWork:
		blockers = append(blockers, GateBlocker{Kind: GateBlockClearanceNeedsWork, Detail: fmt.Sprintf("latest check %s recorded NEEDS WORK", clearance.Check)})
	case ClearanceStale:
		blockers = append(blockers, GateBlocker{Kind: GateBlockClearanceStale, Detail: fmt.Sprintf("freshness assessment marks latest check %s stale", clearance.Check)})
	case ClearanceUnknown:
		if untrustedCurrentClearance(index, objectiveID, clearance.Check) {
			blockers = append(blockers, GateBlocker{Kind: GateBlockCheckerAuthority, Detail: fmt.Sprintf("latest check %s lacks independent checker provenance", clearance.Check)})
		} else {
			blockers = append(blockers, GateBlocker{Kind: GateBlockClearanceUnknown, Detail: fmt.Sprintf("freshness assessment marks latest check %s unknown", clearance.Check)})
		}
	case ClearanceCurrent:
		if ownerValidationRequired(objective.Evidence) && !ownerAcceptedCheck(objective.Evidence, clearance.Check) {
			blockers = append(blockers, GateBlocker{Kind: GateBlockOwnerAcceptance, Detail: fmt.Sprintf("owner has not accepted current check %s", clearance.Check)})
		}
	}

	unresolvedIssue := false
	if clearance.State == ClearanceCurrent {
		for _, issueID := range checkIssueIDs(index, clearance.Check) {
			issue := index.Issues[issueID]
			if issue == nil || issue.Status == IssueStatusResolved {
				continue
			}
			unresolvedIssue = true
			blockers = append(blockers, GateBlocker{
				Kind:   GateBlockObjectiveIssueUnresolved,
				Issue:  issueID,
				Detail: fmt.Sprintf("material Issue %s linked to Objective Check %s remains %s", issueID, clearance.Check, issue.Status),
			})
		}
	}

	if len(blockers) == 0 {
		return GateDecision{Allowed: true, Actor: ActorRoleChecker}
	}

	// An Objective exception cannot silently accept a still-open material Issue.
	// Owner acceptance of that Issue is recorded as its resolution.
	if unresolvedIssue {
		return GateDecision{Blockers: blockers}
	}
	if exception := applicableException(objective.Evidence, index.LatestCheck[objectiveID]); exception != nil {
		return GateDecision{Allowed: true, Actor: ActorRoleOwner, AllowedByException: true, Exception: exception}
	}

	return GateDecision{Blockers: blockers}
}

// ObjectiveDependencyBlockKind names why a V2 Objective dependency is
// unsatisfied.
type ObjectiveDependencyBlockKind string

const (
	// ObjectiveDependencyBlockNotDone means the dependency Objective's status
	// is not done.
	ObjectiveDependencyBlockNotDone ObjectiveDependencyBlockKind = "not_done"
	// ObjectiveDependencyBlockNotCleared means the dependency Objective is
	// done but its integration clearance is not current; Clearance names the
	// observed state.
	ObjectiveDependencyBlockNotCleared ObjectiveDependencyBlockKind = "not_cleared"
	// ObjectiveDependencyBlockClearedByException means the dependency
	// Objective reached done under a recorded exception rather than current
	// clearance, which does not satisfy readiness.
	ObjectiveDependencyBlockClearedByException ObjectiveDependencyBlockKind = "cleared_by_exception"
)

// ObjectiveDependencyBlock is the typed reason one Objective dependency is
// unsatisfied. It mirrors DependencyBlock's shape so a Task's Objective
// dependency block reads the same way as its Task dependency blocks.
type ObjectiveDependencyBlock struct {
	Target    string
	Kind      ObjectiveDependencyBlockKind
	Clearance ClearanceState // set only for ObjectiveDependencyBlockNotCleared
}

// ObjectiveDependencyDecision is the resolved satisfaction for one Objective
// dependency.
type ObjectiveDependencyDecision struct {
	Satisfied bool
	Block     *ObjectiveDependencyBlock // nil when Satisfied is true
}

// ResolveObjectiveDependency resolves whether dependencyID, named in an
// owning Objective's depends_on, is satisfied — re-resolving the dependency
// Objective's status and clearance from index on every call rather than
// trusting any cached judgement. It is satisfied only when the dependency
// Objective is done and its integration clearance, resolved by
// ResolveClearance over its scope.kind: objective Checks, is current;
// Task-only clearance inside that Objective is not enough. A dependency
// Objective that reached done only through a recorded exception is reported
// distinctly rather than silently read as current clearance.
func ResolveObjectiveDependency(index *V2Index, dependencyID string) ObjectiveDependencyDecision {
	target, ok := index.Objectives[dependencyID]
	if !ok || target.Status != ColumnDone {
		return ObjectiveDependencyDecision{Block: &ObjectiveDependencyBlock{Target: dependencyID, Kind: ObjectiveDependencyBlockNotDone}}
	}

	clearance := ResolveClearance(index, dependencyID)
	if clearance.State == ClearanceCurrent {
		return ObjectiveDependencyDecision{Satisfied: true}
	}

	if applicableException(target.Evidence, index.LatestCheck[dependencyID]) != nil {
		return ObjectiveDependencyDecision{Block: &ObjectiveDependencyBlock{Target: dependencyID, Kind: ObjectiveDependencyBlockClearedByException}}
	}

	return ObjectiveDependencyDecision{Block: &ObjectiveDependencyBlock{Target: dependencyID, Kind: ObjectiveDependencyBlockNotCleared, Clearance: clearance.State}}
}

// ObjectiveConsistencyDiagnosticKind names one way an Objective's recorded
// status can contradict its recorded integration evidence or its owned Tasks
// after a hand edit.
type ObjectiveConsistencyDiagnosticKind string

const (
	// ObjectiveConsistencyDoneWithoutClearance means an Objective's status is
	// done but its own integration clearance is not current.
	ObjectiveConsistencyDoneWithoutClearance ObjectiveConsistencyDiagnosticKind = "done_without_current_clearance"
	// ObjectiveConsistencyIncompleteTask means an Objective's status is done
	// while one of its owned Tasks is not done.
	ObjectiveConsistencyIncompleteTask ObjectiveConsistencyDiagnosticKind = "done_with_incomplete_task"
	// ObjectiveConsistencyPlannedWithStartedTask means an Objective's status is
	// still planned while one of its owned Tasks has started or finished, so
	// every surface showing the Objective's status understates its progress.
	ObjectiveConsistencyPlannedWithStartedTask ObjectiveConsistencyDiagnosticKind = "planned_with_started_task"
)

// ObjectiveConsistencyDiagnostic names one inconsistency
// InspectObjectiveConsistency found between an Objective's recorded status
// and its recorded integration evidence or owned Tasks.
type ObjectiveConsistencyDiagnostic struct {
	Objective string
	Kind      ObjectiveConsistencyDiagnosticKind
	Detail    string
}

// InspectObjectiveConsistency reports every inconsistency between an
// Objective's recorded status and its recorded integration evidence or owned
// Tasks, without rewriting any record: an Objective done without current
// integration clearance or an applicable owner exception, an Objective done while an owned Task is not done,
// and an Objective still planned after one of its Tasks started. It walks Objective IDs in sorted order and returns every problem
// found across every Objective, not only the first, mirroring
// InspectTaskConsistency's read-only, sorted, return-everything shape.
func InspectObjectiveConsistency(index *V2Index) []ObjectiveConsistencyDiagnostic {
	var diagnostics []ObjectiveConsistencyDiagnostic

	for _, id := range slices.Sorted(maps.Keys(index.Objectives)) {
		objective := index.Objectives[id]
		if objective.Status == ColumnPlanned {
			if diagnostic, ok := plannedWithStartedTask(index, id); ok {
				diagnostics = append(diagnostics, diagnostic)
			}
			continue
		}
		if objective.Status != ColumnDone {
			continue
		}

		// An owner exception naming the latest Check completes the Objective
		// as ResolveObjectiveCompletion allows; it is not a contradiction.
		clearance := ResolveClearance(index, id)
		if clearance.State != ClearanceCurrent && applicableException(objective.Evidence, index.LatestCheck[id]) == nil {
			diagnostics = append(diagnostics, ObjectiveConsistencyDiagnostic{
				Objective: id,
				Kind:      ObjectiveConsistencyDoneWithoutClearance,
				Detail:    fmt.Sprintf("objective is done but clearance is %s", clearance.State),
			})
		}

		for _, taskID := range index.ObjectiveTasks[id] {
			task := index.Tasks[taskID]
			if task.Status != ColumnDone {
				diagnostics = append(diagnostics, ObjectiveConsistencyDiagnostic{
					Objective: id,
					Kind:      ObjectiveConsistencyIncompleteTask,
					Detail:    fmt.Sprintf("objective is done but task %s is not done (status %q)", taskID, task.Status),
				})
			}
		}
	}

	return diagnostics
}

// plannedWithStartedTask names the first owned Task, in the Objective's sorted
// Task order, that has left planned while its Objective has not. One
// diagnostic per Objective is enough: the repair is the same single status
// change however many Tasks have started.
func plannedWithStartedTask(index *V2Index, objectiveID string) (ObjectiveConsistencyDiagnostic, bool) {
	for _, taskID := range index.ObjectiveTasks[objectiveID] {
		task := index.Tasks[taskID]
		if task != nil && task.Status != ColumnPlanned {
			return ObjectiveConsistencyDiagnostic{
				Objective: objectiveID,
				Kind:      ObjectiveConsistencyPlannedWithStartedTask,
				Detail:    fmt.Sprintf("objective is planned but task %s has started (status %q)", taskID, task.Status),
			}, true
		}
	}
	return ObjectiveConsistencyDiagnostic{}, false
}
