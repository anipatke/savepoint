package v2

import "github.com/opencode/savepoint/internal/data"

// A badge compresses a Check into a glyph. This file resolves what the glyph
// compressed: the whole chain of Checks recorded against one target, in the
// order they were recorded.
//
// Rendering only the newest one would hide the thing a history is for. A target
// that was CLEAR, then NEEDS WORK, then CLEAR again says something no single
// record does, and a Check is immutable — a rerun is a new record naming the
// one it replaces — so the chain is the evidence. Nothing here drops an entry.

// CheckEntry is one recorded Check together with its place in that chain.
// Latest is index.LatestCheck's own answer about which Check counts, and
// Superseded is read from the supersedes link a later Check recorded; neither
// is worked out by comparing timestamps or reading a result.
type CheckEntry struct {
	Check      *data.CheckV2
	Latest     bool
	Superseded bool
}

// checkHistory is every Check naming targetID as its scope, in the recorded
// order index.ScopeChecks holds them in. targetID may name a Task or an
// Objective: both are Check scope targets, and both read the same two maps.
func checkHistory(index *data.V2Index, targetID string) []CheckEntry {
	ids := index.ScopeChecks[targetID]
	if len(ids) == 0 {
		return nil
	}

	// A Check is superseded because a later one says so, not because a newer
	// one exists: an entry with no successor naming it stays unmarked.
	superseded := make(map[string]bool, len(ids))
	for _, id := range ids {
		if replaced := index.Checks[id].Supersedes; replaced != "" {
			superseded[replaced] = true
		}
	}

	latest := index.LatestCheck[targetID]
	entries := make([]CheckEntry, 0, len(ids))
	for _, id := range ids {
		entries = append(entries, CheckEntry{
			Check:      index.Checks[id],
			Latest:     id == latest,
			Superseded: superseded[id],
		})
	}
	return entries
}
