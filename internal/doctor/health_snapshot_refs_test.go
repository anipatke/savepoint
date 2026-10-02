package doctor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/codehealth"
	"github.com/opencode/savepoint/internal/testutil"
)

func saveTestSnapshot(t *testing.T, root string, origin codehealth.Origin) string {
	t.Helper()
	retention := codehealth.RetentionPermanent
	if origin == codehealth.OriginManual {
		retention = codehealth.RetentionPrunable
	}
	snapshot := codehealth.Snapshot{
		Version:    codehealth.SnapshotVersion,
		Origin:     origin,
		Retention:  retention,
		CreatedAt:  "2026-10-01T10:00:00Z",
		Repository: codehealth.RepositoryIdentity{Dirty: true, InputFingerprint: "sha256:" + strings.Repeat("a", 64)},
		Summary:    codehealth.Summary{Overall: codehealth.ClassificationUnknown},
	}
	snapshot.ID = snapshot.ComputeID()
	if _, err := codehealth.NewStore(filepath.Dir(root)).SaveSnapshot(snapshot); err != nil {
		t.Fatalf("SaveSnapshot() error = %v", err)
	}
	return snapshot.ID
}

func writeHealthRefCheck(t *testing.T, root, id, ref string) string {
	t.Helper()
	path := filepath.Join(root, "checks", id+".md")
	testutil.WriteFile(t, path, "---\nid: "+id+"\nscope: {kind: task, id: T-001}\nresult: CLEAR\n"+
		"checked_by: {role: checker, session: sess-1}\nexecuted_session: build-1\nchecked_at: '2026-09-14T00:00:00Z'\n"+
		"health_snapshot: "+ref+"\n---\n\n# Check\n")
	return path
}

// healthProject returns a complete V2 project whose root is a .savepoint
// directory, which is where the health store expects to find its parent.
func healthProject(t *testing.T) string {
	t.Helper()
	root := filepath.Join(t.TempDir(), ".savepoint")
	writeCompleteV2Project(t, root)
	return root
}

func healthRefMessages(report *DiagnosticReport) []string {
	var out []string
	for _, p := range report.Project {
		if strings.Contains(p.Message, "health-snapshot") {
			out = append(out, p.Message+" | "+p.File+" | "+p.Repair)
		}
	}
	return out
}

func TestHealthSnapshotRefs_officialReferenceIsClean(t *testing.T) {
	root := healthProject(t)
	id := saveTestSnapshot(t, root, codehealth.OriginOfficial)
	writeHealthRefCheck(t, root, "C-001", id)

	report := RunV2Checks(root)
	if got := healthRefMessages(report); len(got) != 0 {
		t.Fatalf("official reference reported %v, want none", got)
	}
}

func TestHealthSnapshotRefs_missingSnapshot(t *testing.T) {
	root := healthProject(t)
	saveTestSnapshot(t, root, codehealth.OriginOfficial)
	path := writeHealthRefCheck(t, root, "C-001", "sha256:"+strings.Repeat("b", 64))

	got := healthRefMessages(RunV2Checks(root))
	if len(got) != 1 || !strings.Contains(got[0], "[health-snapshot-missing]") || !strings.Contains(got[0], path) || !strings.Contains(got[0], "Set health_snapshot") {
		t.Fatalf("messages = %v, want one health-snapshot-missing naming %s with a repair hint", got, path)
	}
}

func TestHealthSnapshotRefs_missingHealthDirectoryIsMissingSnapshot(t *testing.T) {
	root := healthProject(t)
	writeHealthRefCheck(t, root, "C-001", "sha256:"+strings.Repeat("b", 64))

	got := healthRefMessages(RunV2Checks(root))
	if len(got) != 1 || !strings.Contains(got[0], "[health-snapshot-missing]") {
		t.Fatalf("messages = %v, want one health-snapshot-missing", got)
	}
}

func TestHealthSnapshotRefs_manualSnapshot(t *testing.T) {
	root := healthProject(t)
	id := saveTestSnapshot(t, root, codehealth.OriginManual)
	path := writeHealthRefCheck(t, root, "C-001", id)

	got := healthRefMessages(RunV2Checks(root))
	if len(got) != 1 || !strings.Contains(got[0], "[health-snapshot-manual]") || !strings.Contains(got[0], path) || !strings.Contains(got[0], "official snapshot") {
		t.Fatalf("messages = %v, want one health-snapshot-manual naming %s with a repair hint", got, path)
	}
}

func TestHealthSnapshotRefs_checksWithoutFieldNeverTouchHealthDirectory(t *testing.T) {
	root := healthProject(t)
	writeV2Check(t, root, "C-001", "{kind: task, id: T-001}", "CLEAR", "")
	// An unreadable health directory would be a finding only if it were read.
	testutil.WriteFile(t, filepath.Join(root, "health", "snapshots", "junk.txt"), "not a snapshot")

	report := RunV2Checks(root)
	if got := healthRefMessages(report); len(got) != 0 {
		t.Fatalf("Check without health_snapshot reported %v, want none", got)
	}
}

func TestHealthSnapshotRefs_unreadableStoreIsReported(t *testing.T) {
	root := healthProject(t)
	writeHealthRefCheck(t, root, "C-001", "sha256:"+strings.Repeat("b", 64))
	testutil.WriteFile(t, filepath.Join(root, "health", "snapshots", "junk.txt"), "not a snapshot")

	got := healthRefMessages(RunV2Checks(root))
	if len(got) != 1 || !strings.Contains(got[0], "[health-snapshot-unreadable]") {
		t.Fatalf("messages = %v, want one health-snapshot-unreadable", got)
	}
}

func TestHealthSnapshotRefs_doctorStaysReadOnly(t *testing.T) {
	root := healthProject(t)
	id := saveTestSnapshot(t, root, codehealth.OriginManual)
	writeHealthRefCheck(t, root, "C-001", id)
	writeHealthRefCheck(t, root, "C-002", "sha256:"+strings.Repeat("b", 64))

	before := snapshotTree(t, root)
	RunV2Checks(root)
	after := snapshotTree(t, root)
	if before != after {
		t.Fatalf("doctor changed the project:\nbefore:\n%s\nafter:\n%s", before, after)
	}
}

func snapshotTree(t *testing.T, root string) string {
	t.Helper()
	var b strings.Builder
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, path)
		b.WriteString(rel)
		if info.Mode().IsRegular() {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			b.WriteString(" " + string(content))
		}
		b.WriteString("\n")
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return b.String()
}
