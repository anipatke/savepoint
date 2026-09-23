package migrate

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/opencode/savepoint/internal/data"
)

func cutoverTestOptions() CutoverPreflightOptions {
	return CutoverPreflightOptions{
		Now: func() time.Time {
			return time.Date(2026, 9, 20, 0, 0, 0, 0, time.UTC)
		},
		NewOperationID: func() string { return "op-cutover-test" },
	}
}

func TestPreflightCutover_refusesEveryOperationalCondition(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		setup      func(*testing.T) string
		wantKinds  []CutoverBlockKind
		wantDetail string
	}{
		{
			name: "target is not a project",
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			wantKinds: []CutoverBlockKind{CutoverBlockTarget},
		},
		{
			name: "unsupported schema",
			setup: func(t *testing.T) string {
				root := t.TempDir()
				writeFile(t, filepath.Join(root, ".savepoint", "config.yml"), "schema_version: 3\n")
				return root
			},
			wantKinds:  []CutoverBlockKind{CutoverBlockUnsupportedSchema},
			wantDetail: "schema_version 3",
		},
		{
			name: "malformed schema",
			setup: func(t *testing.T) string {
				root := t.TempDir()
				writeFile(t, filepath.Join(root, ".savepoint", "config.yml"), "schema_version: not-a-number\n")
				return root
			},
			wantKinds:  []CutoverBlockKind{CutoverBlockMalformedSchema},
			wantDetail: "schema_version",
		},
		{
			name: "invalid staged V2",
			setup: func(t *testing.T) string {
				root := t.TempDir()
				writeFile(t, filepath.Join(root, ".savepoint", "config.yml"), "schema_version: 2\n")
				writeFile(t, filepath.Join(root, ".savepoint", "objectives", "O-001-broken", "Objective.md"), "---\nid: O-001\nstatus: done\n---\n\n# Broken\n")
				return root
			},
			wantKinds:  []CutoverBlockKind{CutoverBlockInvalidV2},
			wantDetail: "V2 project cannot be loaded safely",
		},
		{
			name: "ambiguous migration decision",
			setup: func(t *testing.T) string {
				root := t.TempDir()
				writeFile(t, filepath.Join(root, ".savepoint", "config.yml"), "quality_gates: {}\n")
				writeFile(t, filepath.Join(root, ".savepoint", "releases", "v9", "v9-PRD.md"), "---\nname: Ambiguous\nstatus: audited\n---\n\n# Ambiguous\n")
				return root
			},
			wantKinds:  []CutoverBlockKind{CutoverBlockMigrationRequired, CutoverBlockMigrationAmbiguity},
			wantDetail: "unresolved migration decision",
		},
		{
			name: "migration source conflict",
			setup: func(t *testing.T) string {
				root := t.TempDir()
				writeFile(t, filepath.Join(root, ".savepoint", "config.yml"), "quality_gates: {}\n")
				original := []byte("source before operation\n")
				writeFile(t, filepath.Join(root, "README.md"), string(original))
				_, err := CreateOperation(root, "op-source-conflict", []ManifestSource{{Path: "README.md", SHA256: hashCutoverBytes(original)}}, nil, time.Now().UTC())
				if err != nil {
					t.Fatalf("CreateOperation() error = %v", err)
				}
				writeFile(t, filepath.Join(root, "README.md"), "source after operation\n")
				return root
			},
			wantKinds:  []CutoverBlockKind{CutoverBlockPendingOperation, CutoverBlockRecoveryConflict},
			wantDetail: "not safe to resume",
		},
		{
			name: "incomplete operation has no recovery plan",
			setup: func(t *testing.T) string {
				root := t.TempDir()
				writeFile(t, filepath.Join(root, ".savepoint", "config.yml"), "quality_gates: {}\n")
				_, err := CreateOperation(root, "op-no-plan", nil, []JournalEntry{{
					Path: ".savepoint/generated.md", Action: ActionCreate, PlannedHash: strings.Repeat("a", 64),
				}}, time.Now().UTC())
				if err != nil {
					t.Fatalf("CreateOperation() error = %v", err)
				}
				return root
			},
			wantKinds:  []CutoverBlockKind{CutoverBlockPendingOperation, CutoverBlockUnrecoverablePlan},
			wantDetail: "no verified recovery path",
		},
		{
			name: "multiple pending operations",
			setup: func(t *testing.T) string {
				root := t.TempDir()
				writeFile(t, filepath.Join(root, ".savepoint", "config.yml"), "quality_gates: {}\n")
				for _, id := range []string{"op-one", "op-two"} {
					if _, err := CreateOperation(root, id, nil, nil, time.Now().UTC()); err != nil {
						t.Fatalf("CreateOperation(%s) error = %v", id, err)
					}
				}
				return root
			},
			wantKinds:  []CutoverBlockKind{CutoverBlockPendingOperation},
			wantDetail: "more than one incomplete migration operation",
		},
		{
			name: "migration destination conflict",
			setup: func(t *testing.T) string {
				root := t.TempDir()
				writeFile(t, filepath.Join(root, ".savepoint", "config.yml"), "quality_gates: {}\n")
				writeFile(t, filepath.Join(root, ".savepoint", "migrations", "v1-to-v2.yml"), "existing: true\n")
				return root
			},
			wantKinds:  []CutoverBlockKind{CutoverBlockMigrationRequired, CutoverBlockMigrationConflict},
			wantDetail: "migration conflict",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := tt.setup(t)
			got := PreflightCutover(root, cutoverTestOptions())
			if got.Allowed {
				t.Fatalf("PreflightCutover() allowed unsafe candidate: %+v", got)
			}
			for _, want := range tt.wantKinds {
				if !hasCutoverBlock(got, want) {
					t.Errorf("blockers = %+v, want category %q", got.Blockers, want)
				}
			}
			if tt.wantDetail != "" && !containsCutoverDetail(got, tt.wantDetail) {
				t.Errorf("blockers = %+v, want detail containing %q", got.Blockers, tt.wantDetail)
			}
		})
	}
}

func TestPreflightCutover_allowsReleaseFreeV2(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".savepoint", "config.yml"), "schema_version: 2\n")
	before := snapshotCutoverTree(t, root)

	got := PreflightCutover(root, cutoverTestOptions())
	if !got.Allowed || len(got.Blockers) != 0 {
		t.Fatalf("PreflightCutover() = %+v, want allowed Release-free V2 candidate", got)
	}
	if got.Index == nil || len(got.Index.Releases) != 0 {
		t.Fatalf("Index = %+v, want a loaded Release-free V2 index", got.Index)
	}
	assertCutoverTreeEqual(t, before, snapshotCutoverTree(t, root))
}

func TestPreflightCutover_allowsAcceptedMultiReleaseCandidate(t *testing.T) {
	t.Parallel()
	root := writeAcceptedMultiReleaseCandidate(t, 2)
	before := snapshotCutoverTree(t, root)

	got := PreflightCutover(root, cutoverTestOptions())
	if !got.Allowed || len(got.Blockers) != 0 {
		t.Fatalf("PreflightCutover() = %+v, want accepted multi-Release candidate", got)
	}
	if got.Index == nil || len(got.Index.Releases) != 2 {
		t.Fatalf("Index.Releases = %d, want two Releases", len(got.Index.Releases))
	}
	assertCutoverTreeEqual(t, before, snapshotCutoverTree(t, root))
}

func TestPreflightCutover_translatesCanonicalReleaseBlockers(t *testing.T) {
	t.Parallel()
	root := writeAcceptedMultiReleaseCandidate(t, 2)
	releasePath := filepath.Join(root, ".savepoint", "releases", "R-002-release-2", "Release.md")
	content, err := os.ReadFile(releasePath)
	if err != nil {
		t.Fatalf("ReadFile(%s) error = %v", releasePath, err)
	}
	content = []byte(strings.Replace(string(content), "owner_validation:\n  required: true\n  accepted_check: C-004\n  accepted_by: {role: owner, session: owner-2}\n", "", 1))
	if err := os.WriteFile(releasePath, content, 0644); err != nil {
		t.Fatalf("WriteFile(%s) error = %v", releasePath, err)
	}

	got := PreflightCutover(root, cutoverTestOptions())
	if got.Allowed {
		t.Fatalf("PreflightCutover() allowed an unaccepted Release: %+v", got)
	}
	var found bool
	for _, blocker := range got.Blockers {
		if blocker.Kind == CutoverBlockRelease && blocker.ReleaseID == "R-002" && blocker.Gate.Kind == data.GateBlockOwnerAcceptance {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("Blockers = %+v, want canonical R-002 owner-acceptance blocker", got.Blockers)
	}
}

func TestPreflightCutover_isReadOnlyAcrossV1V2AndRecovery(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		setup func(*testing.T) string
	}{
		{
			name: "V1 plan",
			setup: func(t *testing.T) string {
				root := t.TempDir()
				writeFile(t, filepath.Join(root, ".savepoint", "config.yml"), "quality_gates: {}\n")
				writeFile(t, filepath.Join(root, ".savepoint", "releases", "v9", "v9-PRD.md"), "---\nname: Ambiguous\nstatus: audited\n---\n\n# Ambiguous\n")
				return root
			},
		},
		{
			name: "V2 index",
			setup: func(t *testing.T) string {
				return writeAcceptedMultiReleaseCandidate(t, 1)
			},
		},
		{
			name: "pending operation",
			setup: func(t *testing.T) string {
				root := t.TempDir()
				writeFile(t, filepath.Join(root, ".savepoint", "config.yml"), "quality_gates: {}\n")
				if _, err := CreateOperation(root, "op-read-only", nil, nil, time.Now().UTC()); err != nil {
					t.Fatalf("CreateOperation() error = %v", err)
				}
				return root
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := tt.setup(t)
			before := snapshotCutoverTree(t, root)
			_ = PreflightCutover(root, cutoverTestOptions())
			assertCutoverTreeEqual(t, before, snapshotCutoverTree(t, root))
		})
	}
}

func hasCutoverBlock(result CutoverPreflightResult, kind CutoverBlockKind) bool {
	for _, blocker := range result.Blockers {
		if blocker.Kind == kind {
			return true
		}
	}
	return false
}

func containsCutoverDetail(result CutoverPreflightResult, text string) bool {
	for _, blocker := range result.Blockers {
		if strings.Contains(blocker.Detail, text) {
			return true
		}
	}
	return false
}

func writeAcceptedMultiReleaseCandidate(t *testing.T, releaseCount int) string {
	t.Helper()
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".savepoint", "config.yml"), "schema_version: 2\n")
	checkedAt := "2026-09-15T00:00:00Z"
	releaseBody := "\n# Release\n\n## Outcome\n\nDeliver the promise.\n\n## Why\n\nA Release owns the delivery promise.\n\n## Success Conditions\n\n- Every member Objective is complete.\n\n## Boundaries\n\nRelease does not own Tasks.\n"
	for i := 1; i <= releaseCount; i++ {
		releaseID := fmt.Sprintf("R-%03d", i)
		objectiveID := fmt.Sprintf("O-%03d", i)
		taskID := fmt.Sprintf("T-%03d", i)
		objectiveCheckID := fmt.Sprintf("C-%03d", (i-1)*2+1)
		releaseCheckID := fmt.Sprintf("C-%03d", i*2)

		writeFile(t, filepath.Join(root, ".savepoint", "releases", releaseID+"-release-"+fmt.Sprint(i), "Release.md"), fmt.Sprintf("---\nid: %s\ntitle: Release %d\nstatus: in_progress\nfreshness:\n  state: current\n  check: %s\n  assessed_by: {role: checker, session: release-checker-%d}\n  assessed_at: '%s'\n  basis: release integration reviewed\nowner_validation:\n  required: true\n  accepted_check: %s\n  accepted_by: {role: owner, session: owner-%d}\n---\n%s", releaseID, i, releaseCheckID, i, checkedAt, releaseCheckID, i, releaseBody))
		writeFile(t, filepath.Join(root, ".savepoint", "objectives", objectiveID+"-objective-"+fmt.Sprint(i), "Objective.md"), fmt.Sprintf("---\nid: %s\ntitle: Objective %d\nstatus: done\nrelease: %s\nfreshness:\n  state: current\n  check: %s\n  assessed_by: {role: checker, session: objective-checker-%d}\n  assessed_at: '%s'\n  basis: objective integration reviewed\n---\n\n# Objective\n", objectiveID, i, releaseID, objectiveCheckID, i, checkedAt))
		writeFile(t, filepath.Join(root, ".savepoint", "objectives", objectiveID+"-objective-"+fmt.Sprint(i), "tasks", taskID+"-task.md"), fmt.Sprintf("---\nid: %s\ntitle: Task %d\nobjective: %s\nplanned_by: {role: planner, session: plan-%d}\nstatus: done\n---\n\n# Task\n", taskID, i, objectiveID, i))
		writeFile(t, filepath.Join(root, ".savepoint", "checks", objectiveCheckID+"-objective.md"), fmt.Sprintf("---\nid: %s\nscope: {kind: objective, id: %s}\nresult: CLEAR\nchecked_by: {role: checker, session: objective-checker-%d}\nexecuted_session: objective-build-%d\nchecked_at: '%s'\n---\n\n# Check\n", objectiveCheckID, objectiveID, i, i, checkedAt))
		writeFile(t, filepath.Join(root, ".savepoint", "checks", releaseCheckID+"-release.md"), fmt.Sprintf("---\nid: %s\nscope: {kind: release, id: %s}\nresult: CLEAR\nchecked_by: {role: checker, session: release-checker-%d}\nexecuted_session: release-build-%d\nchecked_at: '%s'\n---\n\n# Check\n", releaseCheckID, releaseID, i, i, checkedAt))
	}
	return root
}

type cutoverFileSnapshot struct {
	directory bool
	content   []byte
	modTime   time.Time
}

func snapshotCutoverTree(t *testing.T, root string) map[string]cutoverFileSnapshot {
	t.Helper()
	snapshot := make(map[string]cutoverFileSnapshot)
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		state := cutoverFileSnapshot{directory: entry.IsDir(), modTime: info.ModTime()}
		if !entry.IsDir() {
			state.content, err = os.ReadFile(path)
			if err != nil {
				return err
			}
		}
		snapshot[rel] = state
		return nil
	})
	if err != nil {
		t.Fatalf("snapshotCutoverTree(%s) error = %v", root, err)
	}
	return snapshot
}

func assertCutoverTreeEqual(t *testing.T, before, after map[string]cutoverFileSnapshot) {
	t.Helper()
	if len(before) != len(after) {
		t.Fatalf("filesystem entry count changed: before %d, after %d", len(before), len(after))
	}
	for path, want := range before {
		got, ok := after[path]
		if !ok {
			t.Errorf("filesystem entry %s was removed", path)
			continue
		}
		if want.directory != got.directory {
			t.Errorf("filesystem entry %s changed directory kind", path)
		}
		if !want.modTime.Equal(got.modTime) {
			t.Errorf("filesystem entry %s modTime changed: before %v, after %v", path, want.modTime, got.modTime)
		}
		if string(want.content) != string(got.content) {
			t.Errorf("filesystem entry %s content changed", path)
		}
	}
}

func hashCutoverBytes(content []byte) string {
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:])
}
