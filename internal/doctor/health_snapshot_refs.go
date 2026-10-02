package doctor

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/opencode/savepoint/internal/codehealth"
	"github.com/opencode/savepoint/internal/data"
)

const (
	healthSnapshotMissing = "health-snapshot-missing"
	healthSnapshotManual  = "health-snapshot-manual"
	healthSnapshotUnread  = "health-snapshot-unreadable"
)

// healthSnapshotRefProblems reports Checks whose health_snapshot names no
// stored snapshot, or names a manual one that a Check may not rely on. It only
// reads: nothing is created, repaired, or pruned. Snapshots load once, and
// only when some Check carries the field, so a project without a health
// directory is never flagged for Checks that do not use it.
func healthSnapshotRefProblems(root string, index *data.V2Index) []Problem {
	var refs []*data.CheckV2
	for _, check := range index.Checks {
		if check.HealthSnapshot != "" {
			refs = append(refs, check)
		}
	}
	if len(refs) == 0 {
		return nil
	}
	sort.Slice(refs, func(i, j int) bool { return refs[i].ID < refs[j].ID })

	snapshots, err := codehealth.NewStore(filepath.Dir(root)).LoadSnapshots()
	if err != nil {
		return []Problem{{
			File:     filepath.Join(root, "health", "snapshots"),
			Message:  fmt.Sprintf("[%s] cannot read stored snapshots to verify %d Check reference(s): %v", healthSnapshotUnread, len(refs), err),
			Repair:   V2ProblemRepair(healthSnapshotUnread),
			Category: HealthMalformedData,
		}}
	}
	origins := make(map[string]codehealth.Origin, len(snapshots))
	for _, snapshot := range snapshots {
		origins[snapshot.ID] = snapshot.Origin
	}

	var problems []Problem
	for _, check := range refs {
		origin, found := origins[check.HealthSnapshot]
		name := ""
		switch {
		case !found:
			name = healthSnapshotMissing
		case origin == codehealth.OriginManual:
			name = healthSnapshotManual
		default:
			continue
		}
		problems = append(problems, Problem{
			File:     filepath.Join(root, check.Source.Path),
			Message:  fmt.Sprintf("[%s] check %s health_snapshot %q", name, check.ID, check.HealthSnapshot),
			Repair:   V2ProblemRepair(name),
			Category: HealthMissingEvidence,
		})
	}
	return problems
}
