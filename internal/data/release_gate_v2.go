package data

import (
	"fmt"
	"slices"
	"strings"
)

// ResolveReleaseCompletion decides whether a Release may be recorded as done.
// Release membership comes only from Objective release references, and every
// member must pass the existing Objective completion decision.
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

	return GateDecision{Allowed: true, Actor: ActorRoleOwner}
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
