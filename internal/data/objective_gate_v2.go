package data

import "fmt"

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
// never as a CLEAR result or current clearance. An unfinished owned Task is
// not excusable by exception: cross-Task repair goes back through Tasks, and
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
		blockers = append(blockers, GateBlocker{Kind: GateBlockClearanceStale, Detail: fmt.Sprintf("freshness assessment does not name latest check %s as current", clearance.Check)})
	case ClearanceUnknown:
		if untrustedCurrentClearance(index, objectiveID, clearance.Check) {
			blockers = append(blockers, GateBlocker{Kind: GateBlockCheckerAuthority, Detail: fmt.Sprintf("latest check %s lacks independent checker provenance", clearance.Check)})
		} else {
			blockers = append(blockers, GateBlocker{Kind: GateBlockClearanceUnknown, Detail: fmt.Sprintf("no freshness assessment recorded for latest check %s", clearance.Check)})
		}
	case ClearanceCurrent:
		if ownerValidationRequired(objective.Evidence) && !ownerAcceptedCheck(objective.Evidence, clearance.Check) {
			blockers = append(blockers, GateBlocker{Kind: GateBlockOwnerAcceptance, Detail: fmt.Sprintf("owner has not accepted current check %s", clearance.Check)})
		}
	}

	if len(blockers) == 0 {
		return GateDecision{Allowed: true, Actor: ActorRoleChecker}
	}

	if exception := applicableException(objective.Evidence, index.LatestCheck[objectiveID]); exception != nil {
		return GateDecision{Allowed: true, Actor: ActorRoleOwner, AllowedByException: true, Exception: exception}
	}

	return GateDecision{Blockers: blockers}
}
