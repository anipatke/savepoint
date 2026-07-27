package init

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/opencode/savepoint/internal/testutil"
)

// An upgrade writes several files in sequence, so the honest question is not
// whether a write can fail — it is what the project looks like afterwards. The
// tests below inject a failure at each point in that sequence and assert the
// state that remains: nothing beyond the failure is applied, whatever was
// applied is reported, the user's content is still recoverable, and provenance
// matches what is actually on disk.

var errInjected = errors.New("injected write failure")

// countingWriter fails the failAt-th write (1-based) and performs every other
// one for real.
type countingWriter struct {
	failAt  int
	calls   int
	written []string
}

func (w *countingWriter) write(path string, content []byte) error {
	w.calls++
	if w.calls == w.failAt {
		return errInjected
	}
	if err := AtomicWrite(path, content); err != nil {
		return err
	}
	w.written = append(w.written, path)
	return nil
}

const (
	skillA = "agent-skills/skill-a/SKILL.md"
	skillB = "agent-skills/skill-b/SKILL.md"
)

// failureProject builds a project whose upgrade writes a known sequence under
// --force: skill-a's backup, skill-a, skill-b, then the manifest.
func failureProject(t *testing.T) (dir string, templates fstest.MapFS) {
	t.Helper()
	dir = t.TempDir()
	testutil.MkdirAll(t, filepath.Join(dir, ".savepoint"))

	testutil.WriteFile(t, filepath.Join(dir, filepath.FromSlash(skillA)), "# My edits")
	testutil.WriteFile(t, filepath.Join(dir, filepath.FromSlash(skillB)), "# Ours v1")

	manifest := NewManifest()
	manifest.Record(skillA, []byte("# Ours v1"))
	manifest.Record(skillB, []byte("# Ours v1"))
	if err := manifest.Save(dir); err != nil {
		t.Fatal(err)
	}

	return dir, fstest.MapFS{
		skillA: &fstest.MapFile{Data: []byte("# Ours v2")},
		skillB: &fstest.MapFile{Data: []byte("# Ours v2")},
	}
}

func projectFile(t *testing.T, dir, path string) (string, bool) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(path)))
	if os.IsNotExist(err) {
		return "", false
	}
	if err != nil {
		t.Fatal(err)
	}
	return string(data), true
}

func assertFileIs(t *testing.T, dir, path, want string) {
	t.Helper()
	got, ok := projectFile(t, dir, path)
	if !ok {
		t.Fatalf("%s missing, want %q", path, want)
	}
	if got != want {
		t.Errorf("%s = %q, want %q", path, got, want)
	}
}

func TestUpgrade_writeFailureLeavesRecoverableState(t *testing.T) {
	// The write sequence under --force, one case per failure point.
	cases := []struct {
		name   string
		failAt int
	}{
		{"before any write", 1},
		{"after the backup, before the replacement", 2},
		{"after a live asset replacement", 3},
		{"while committing the manifest", 4},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir, templates := failureProject(t)
			manifestBefore, _ := projectFile(t, dir, ".savepoint/.upgrade-manifest.yml")

			w := &countingWriter{failAt: tc.failAt}
			report, err := upgradeProjectAssets(templates, dir, false, true, w.write)

			if !errors.Is(err, errInjected) {
				t.Fatalf("error = %v, want the injected failure to surface", err)
			}
			if report == nil {
				t.Fatal("no report returned; the applied work is invisible to the user")
			}

			backup, hasBackup := projectFile(t, dir, skillA+backupSuffix)
			liveA, _ := projectFile(t, dir, skillA)
			liveB, _ := projectFile(t, dir, skillB)
			manifestAfter, _ := projectFile(t, dir, ".savepoint/.upgrade-manifest.yml")

			switch tc.failAt {
			case 1:
				// Nothing was applied, so nothing may have changed.
				if hasBackup {
					t.Error("failed-before-any-write upgrade left a backup")
				}
				assertFileIs(t, dir, skillA, "# My edits")
				assertFileIs(t, dir, skillB, "# Ours v1")
				if manifestAfter != manifestBefore {
					t.Error("failed-before-any-write upgrade rewrote the manifest")
				}
			case 2:
				// The replacement failed, but the user's content survives twice
				// over: untouched in place, and copied to the backup.
				if !hasBackup || backup != "# My edits" {
					t.Errorf("backup = %q (exists %v), want the user's content", backup, hasBackup)
				}
				assertFileIs(t, dir, skillA, "# My edits")
				assertFileIs(t, dir, skillB, "# Ours v1")

				// Writing a .bak and not saying so is the same invisible
				// partial state the report exists to prevent.
				entry, found := entryFor(report, skillA)
				if !found {
					t.Fatalf("report does not mention %s at all: %+v", skillA, report.Actions)
				}
				if entry.Action != ActionFailed {
					t.Errorf("%s action = %v, want failed", skillA, entry.Action)
				}
				if entry.Note != noteBackup {
					t.Errorf("%s note = %q, want it to name the backup", skillA, entry.Note)
				}
				if !strings.Contains(report.Format(), backupSuffix) {
					t.Errorf("report does not name the backup:\n%s", report.Format())
				}
			case 3:
				// skill-a is applied and skill-b is untouched: no half-written
				// file, and the manifest now matches what is on disk.
				if !hasBackup || backup != "# My edits" {
					t.Errorf("backup = %q (exists %v), want the user's content", backup, hasBackup)
				}
				assertFileIs(t, dir, skillA, "# Ours v2")
				assertFileIs(t, dir, skillB, "# Ours v1")
				assertRecordedHash(t, dir, skillA, "# Ours v2")
				assertRecordedHash(t, dir, skillB, "# Ours v1")
				if action, found := actionFor(report, skillA); !found || action != ActionUpdated {
					t.Errorf("report does not name the applied %s: %v (found %v)", skillA, action, found)
				}
			case 4:
				// Both assets are applied; only the provenance record failed.
				assertFileIs(t, dir, skillA, "# Ours v2")
				assertFileIs(t, dir, skillB, "# Ours v2")
				if !strings.Contains(err.Error(), "manifest") {
					t.Errorf("error = %v, want it to name the manifest", err)
				}
				for _, path := range []string{skillA, skillB} {
					if action, found := actionFor(report, path); !found || action != ActionUpdated {
						t.Errorf("report does not name the applied %s: %v (found %v)", path, action, found)
					}
				}
			}

			if liveA == "" || liveB == "" {
				t.Error("a skill file was left empty by the failed upgrade")
			}
		})
	}
}

func assertRecordedHash(t *testing.T, dir, path, want string) {
	t.Helper()
	manifest, err := LoadManifest(dir)
	if err != nil {
		t.Fatal(err)
	}
	hash, ok := manifest.Hash(path)
	if !ok {
		t.Fatalf("no manifest entry for %s", path)
	}
	if hash != hashContent([]byte(want)) {
		t.Errorf("manifest hash for %s does not match %q", path, want)
	}
}

func TestUpgrade_refusesWhenTheManifestCannotBeWritten(t *testing.T) {
	// The real cause of a late manifest failure is a read-only .savepoint/.
	// The preflight check must turn that into a refusal with nothing applied.
	if runtime.GOOS == "windows" {
		t.Skip("directory permissions are not enforced the same way on Windows")
	}
	if os.Geteuid() == 0 {
		t.Skip("running as root bypasses directory permissions")
	}

	dir, templates := failureProject(t)
	savepointDir := filepath.Join(dir, ".savepoint")
	if err := os.Chmod(savepointDir, 0500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(savepointDir, 0755) })

	_, err := UpgradeProjectAssets(templates, dir, false, true)
	if err == nil {
		t.Fatal("expected an error for an unwritable .savepoint directory")
	}
	if !strings.Contains(err.Error(), "upgrade manifest") {
		t.Errorf("error = %v, want it to name the upgrade manifest", err)
	}

	// Nothing may have been applied: the refusal came before the first write.
	assertFileIs(t, dir, skillA, "# My edits")
	assertFileIs(t, dir, skillB, "# Ours v1")
	if _, ok := projectFile(t, dir, skillA+backupSuffix); ok {
		t.Error("refused upgrade left a backup behind")
	}
}

func TestDiscardProbe_reportsCleanupFailure(t *testing.T) {
	// The probe's cleanup result is part of the failure evidence: a probe that
	// cannot be closed or removed means the filesystem is not in the state the
	// check just claimed it was.
	probe, err := os.CreateTemp(t.TempDir(), ".tmp-*.manifest-probe")
	if err != nil {
		t.Fatal(err)
	}
	name := probe.Name()
	if err := probe.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(name); err != nil {
		t.Fatal(err)
	}

	err = discardProbe(probe)
	if err == nil {
		t.Fatal("discardProbe() = nil, want the close and remove failures reported")
	}
	if !strings.Contains(err.Error(), name) {
		t.Errorf("error = %v, want it to name the probe file", err)
	}
	if !errors.Is(err, os.ErrClosed) {
		t.Errorf("error = %v, want the close failure preserved", err)
	}
}

func TestUpgrade_unchangedRunDoesNotProbe(t *testing.T) {
	// The probe is a real file creation inside .savepoint/. Running it on a
	// no-op upgrade would change that directory's modification time, so a run
	// with nothing to write must never reach it.
	dir, templates := failureProject(t)

	// Make the project match the templates so the second run has no writes.
	if _, err := upgradeProjectAssets(templates, dir, false, true, AtomicWrite); err != nil {
		t.Fatalf("first UpgradeProjectAssets() error = %v", err)
	}

	probed := false
	write := beforeFirstWrite(AtomicWrite, func() error {
		probed = true
		return nil
	})
	report, err := upgradeProjectAssets(templates, dir, false, false, write)
	if err != nil {
		t.Fatalf("second UpgradeProjectAssets() error = %v", err)
	}
	for _, e := range report.Actions {
		if e.Action != ActionUnchanged {
			t.Errorf("second run %s = %v, want unchanged", e.Path, e.Action)
		}
	}
	if probed {
		t.Error("an unchanged upgrade ran the writability probe")
	}
}

func TestUpgrade_dryRunNeedsNoWritableManifest(t *testing.T) {
	// A dry run writes nothing at all, so it must not demand write access.
	if runtime.GOOS == "windows" {
		t.Skip("directory permissions are not enforced the same way on Windows")
	}
	if os.Geteuid() == 0 {
		t.Skip("running as root bypasses directory permissions")
	}

	dir, templates := failureProject(t)
	savepointDir := filepath.Join(dir, ".savepoint")
	if err := os.Chmod(savepointDir, 0500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(savepointDir, 0755) })

	report, err := UpgradeProjectAssets(templates, dir, true, false)
	if err != nil {
		t.Fatalf("dry run on a read-only .savepoint error = %v", err)
	}
	if action, found := actionFor(report, skillA); !found || action != ActionConflict {
		t.Errorf("%s dry-run action = %v (found %v), want conflict", skillA, action, found)
	}
}
