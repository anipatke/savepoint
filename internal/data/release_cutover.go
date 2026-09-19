package data

import (
	"fmt"
	"maps"
	"slices"
)

// ReleaseCutoverBlocker identifies the Release whose canonical completion
// decision is preventing a project-level cutover. Gate is deliberately the
// original GateDecision's blocker rather than a second cutover vocabulary.
type ReleaseCutoverBlocker struct {
	ReleaseID string
	Gate      GateBlocker
}

// ReleaseCutoverDecision is the project-level composition E50 can consume.
// A project with no Releases is allowed because Releases are optional. For
// every Release that exists, ResolveReleaseCompletion remains the only source
// of readiness rules; this function only orders and qualifies those decisions
// with their Release identity.
type ReleaseCutoverDecision struct {
	Allowed  bool
	Blockers []ReleaseCutoverBlocker
}

// ResolveReleaseCutover decides whether all declared Releases permit a V2
// cutover. It does not decide migration-plan ambiguity, schema validity, or
// filesystem recovery; callers must complete those preconditions separately.
// It also does not add a new Release rule: every returned blocker is emitted
// by the canonical ResolveReleaseCompletion resolver.
func ResolveReleaseCutover(index *V2Index) ReleaseCutoverDecision {
	if index == nil {
		return ReleaseCutoverDecision{Blockers: []ReleaseCutoverBlocker{{
			Gate: GateBlocker{
				Kind:   GateBlockInvalidState,
				Detail: "V2 project index is unavailable",
			},
		}}}
	}

	var blockers []ReleaseCutoverBlocker
	for _, releaseID := range slices.Sorted(maps.Keys(index.Releases)) {
		decision := ResolveReleaseCompletion(index, releaseID)
		for _, blocker := range decision.Blockers {
			if blocker.Detail == "" {
				blocker.Detail = fmt.Sprintf("release %s is not ready", releaseID)
			}
			blockers = append(blockers, ReleaseCutoverBlocker{ReleaseID: releaseID, Gate: blocker})
		}
	}

	return ReleaseCutoverDecision{Allowed: len(blockers) == 0, Blockers: blockers}
}
