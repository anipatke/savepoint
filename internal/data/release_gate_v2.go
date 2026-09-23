package data

import (
	"fmt"
	"slices"
	"strings"
)

// ResolveReleaseCompletion decides whether a Release may be recorded as done.
// Release membership comes only from the Objective release references. The
// decision composes the existing Objective completion decision, the shared
// Check freshness resolver, the Check-to-Issue link map, and the shared owner
// exception vocabulary; it does not create a second version of any of those
// rules.
func ResolveReleaseCompletion(index *V2Index, releaseID string) GateDecision {
	if index == nil {
		return GateDecision{}
	}
	release, ok := index.Releases[releaseID]
	if !ok {
		return GateDecision{}
	}
	return resolveReleaseCompletionForRecord(index, release)
}

func resolveReleaseCompletionForRecord(index *V2Index, release *ReleaseV2) GateDecision {
	objectiveIDs := index.ReleaseObjectives[release.ID]
	// A migrated settled Release may have no live Objective members because
	// its completed V1 work was archived. Its typed legacy completion is an
	// explicit historical outcome, not a current CLEAR Check, and is therefore
	// enough to preserve that settled disposition. If live members do exist,
	// they still have to be historically done below.
	if release.Status == ColumnDone && release.LegacyCompletion != nil && len(objectiveIDs) == 0 {
		return GateDecision{
			Allowed:                   true,
			Actor:                     ActorRoleOwner,
			AllowedByLegacyCompletion: true,
			LegacyCompletion:          release.LegacyCompletion,
		}
	}
	if len(objectiveIDs) == 0 {
		return GateDecision{Blockers: []GateBlocker{{
			Kind:   GateBlockReleaseNoObjectives,
			Detail: fmt.Sprintf("release %s has no member objectives", release.ID),
		}}}
	}

	// Migration may explicitly preserve a settled historical disposition. It
	// is a separate allowed outcome: no historical reference is consulted by
	// ResolveClearance and it never becomes a new CLEAR result.
	if release.Status == ColumnDone && release.LegacyCompletion != nil {
		var blockers []GateBlocker
		for _, objectiveID := range objectiveIDs {
			objective, ok := index.Objectives[objectiveID]
			if !ok || objective.Status != ColumnDone {
				blockers = append(blockers, GateBlocker{
					Kind:      GateBlockReleaseObjectiveIncomplete,
					Objective: objectiveID,
					Detail:    fmt.Sprintf("member objective %s is not historically done", objectiveID),
				})
			}
		}
		if len(blockers) == 0 {
			return GateDecision{
				Allowed:                   true,
				Actor:                     ActorRoleOwner,
				AllowedByLegacyCompletion: true,
				LegacyCompletion:          release.LegacyCompletion,
			}
		}
		return GateDecision{Blockers: blockers}
	}

	var blockers []GateBlocker
	for _, objectiveID := range objectiveIDs {
		decision := ResolveObjectiveCompletion(index, objectiveID)
		if decision.Allowed {
			continue
		}
		if len(decision.Blockers) == 0 {
			blockers = append(blockers, GateBlocker{
				Kind:      GateBlockReleaseObjectiveIncomplete,
				Objective: objectiveID,
				Detail:    fmt.Sprintf("member objective %s is not complete", objectiveID),
			})
			continue
		}
		for _, blocker := range decision.Blockers {
			blockers = append(blockers, GateBlocker{
				Kind:      GateBlockReleaseObjectiveIncomplete,
				Objective: objectiveID,
				Detail:    fmt.Sprintf("member objective %s: %s", objectiveID, blocker.Detail),
			})
		}
	}
	if len(blockers) > 0 {
		return GateDecision{Blockers: blockers}
	}

	clearance := ResolveClearance(index, release.ID)
	switch clearance.State {
	case ClearanceMissing:
		return GateDecision{Blockers: []GateBlocker{{Kind: GateBlockClearanceMissing, Detail: "no recorded Release Check"}}}
	case ClearanceNeedsWork:
		return GateDecision{Blockers: []GateBlocker{{Kind: GateBlockClearanceNeedsWork, Detail: fmt.Sprintf("latest Release Check %s recorded NEEDS WORK", clearance.Check)}}}
	case ClearanceStale:
		return GateDecision{Blockers: []GateBlocker{{Kind: GateBlockClearanceStale, Detail: fmt.Sprintf("Release freshness assessment marks latest Check %s stale", clearance.Check)}}}
	case ClearanceUnknown:
		if untrustedCurrentClearance(index, release.ID, clearance.Check) {
			return GateDecision{Blockers: []GateBlocker{{Kind: GateBlockCheckerAuthority, Detail: fmt.Sprintf("latest Release Check %s lacks independent checker provenance", clearance.Check)}}}
		}
		return GateDecision{Blockers: []GateBlocker{{Kind: GateBlockClearanceUnknown, Detail: fmt.Sprintf("Release freshness assessment marks latest Check %s unknown", clearance.Check)}}}
	case ClearanceCurrent:
		// Continue with material Issue and owner-acceptance composition below.
	}

	exception := applicableException(release.Evidence, clearance.Check)
	exceptionUsed := false
	for _, issueID := range checkIssueIDs(index, clearance.Check) {
		issue := index.Issues[issueID]
		if issue == nil || issue.Status == IssueStatusResolved {
			continue
		}
		if exception != nil && releaseExceptionCoversIssue(exception, issueID) {
			exceptionUsed = true
			continue
		}
		blockers = append(blockers, GateBlocker{
			Kind:   GateBlockReleaseIssueUnresolved,
			Issue:  issueID,
			Detail: fmt.Sprintf("material Issue %s linked to Release Check %s remains %s", issueID, clearance.Check, issue.Status),
		})
	}

	// Release acceptance is mandatory even when every material Issue is
	// covered by an owner exception. A Release promise is not complete merely
	// because a checker recorded CLEAR.
	if !ownerAcceptedCheck(release.Evidence, clearance.Check) {
		blockers = append(blockers, GateBlocker{
			Kind:   GateBlockOwnerAcceptance,
			Detail: fmt.Sprintf("owner has not accepted current Release Check %s", clearance.Check),
		})
	}

	if len(blockers) > 0 {
		return GateDecision{Blockers: blockers}
	}
	if exceptionUsed {
		return GateDecision{Allowed: true, Actor: ActorRoleOwner, AllowedByException: true, Exception: exception}
	}
	return GateDecision{Allowed: true, Actor: ActorRoleChecker}
}

// checkIssueIDs reads the indexed Check-to-Issue links and falls back
// to the immutable Check's own list for hand-built unit indexes. A loaded
// project always uses the index map, whose links have already passed the
// pairing and target validation gates.
func checkIssueIDs(index *V2Index, checkID string) []string {
	ids := make(map[string]struct{})
	for _, issueID := range index.CheckIssues[checkID] {
		ids[issueID] = struct{}{}
	}
	if check := index.Checks[checkID]; check != nil {
		for _, issueID := range check.Issues {
			ids[issueID] = struct{}{}
		}
	}
	keys := mapsKeysString(ids)
	slices.SortFunc(keys, strings.Compare)
	return keys
}

func mapsKeysString(values map[string]struct{}) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	return keys
}

func releaseExceptionCoversIssue(exception *Exception, issueID string) bool {
	if exception == nil {
		return false
	}
	for _, requirement := range exception.Requirements {
		if requirement == issueID || requirement == "issue:"+issueID || requirement == "release.issue:"+issueID {
			return true
		}
	}
	return false
}
