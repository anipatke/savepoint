package codehealth

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"testing"
	"time"
)

func newProject(t *testing.T) (Store, string) {
	t.Helper()
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, savepointDir), 0o755); err != nil {
		t.Fatal(err)
	}
	return NewStore(root), root
}

func healthPath(root string, elems ...string) string {
	return filepath.Join(append([]string{root, savepointDir, healthDir}, elems...)...)
}

// snapAt builds a valid snapshot that differs by origin, time, and commit.
func snapAt(t *testing.T, origin Origin, n int) Snapshot {
	t.Helper()
	s := validSnapshot(t)
	s.Origin, s.Retention = OriginManual, RetentionPrunable
	if origin == OriginOfficial {
		s.Origin, s.Retention = OriginOfficial, RetentionPermanent
	}
	s.CreatedAt = time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(n) * time.Minute).Format(timestampLayout)
	s.Repository.Commit = fmt.Sprintf("%040x", n)
	return reseal(s)
}

func mustSave(t *testing.T, st Store, s Snapshot) {
	t.Helper()
	if _, err := st.SaveSnapshot(s); err != nil {
		t.Fatalf("SaveSnapshot() = %v", err)
	}
}

func snapshotFiles(t *testing.T, root string) []string {
	t.Helper()
	entries, err := os.ReadDir(healthPath(root, snapshotsDir))
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	return names
}

func writeFile(t *testing.T, path string, data string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
}

func symlinkOrSkip(t *testing.T, target, link string) {
	t.Helper()
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlinks unavailable on %s: %v", runtime.GOOS, err)
	}
}

// readOnlyOrSkip makes dir unwritable, which root and Windows ignore.
func readOnlyOrSkip(t *testing.T, dir string) {
	t.Helper()
	if runtime.GOOS == "windows" || os.Geteuid() == 0 {
		t.Skip("directory permissions are not enforced here")
	}
	if err := os.Chmod(dir, 0o555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(dir, 0o755) })
}

func TestConfigAbsentIsNotFound(t *testing.T) {
	st, root := newProject(t)
	if _, err := st.LoadConfig(); !errors.Is(err, ErrConfigNotFound) {
		t.Fatalf("no health dir: LoadConfig() = %v, want ErrConfigNotFound", err)
	}
	if err := os.Mkdir(healthPath(root), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := st.LoadConfig(); !errors.Is(err, ErrConfigNotFound) {
		t.Fatalf("empty health dir: LoadConfig() = %v, want ErrConfigNotFound", err)
	}
}

func TestConfigRoundTripAndRepeat(t *testing.T) {
	st, root := newProject(t)
	cfg := validConfig(t)
	if changed, err := st.SaveConfig(cfg); err != nil || !changed {
		t.Fatalf("first SaveConfig() = %v, %v; want changed", changed, err)
	}
	got, err := st.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if got.Validate() != nil || len(got.Capabilities) != len(cfg.Capabilities) {
		t.Fatalf("round trip lost data: %+v", got)
	}
	before, _ := os.Stat(healthPath(root, configFile))
	if changed, err := st.SaveConfig(cfg); err != nil || changed {
		t.Fatalf("repeat SaveConfig() = %v, %v; want unchanged", changed, err)
	}
	after, _ := os.Stat(healthPath(root, configFile))
	if !after.ModTime().Equal(before.ModTime()) {
		t.Fatal("identical save rewrote the file")
	}
}

func TestSaveConfigReplacesAtomicallyWithoutLeftovers(t *testing.T) {
	st, root := newProject(t)
	cfg := validConfig(t)
	if _, err := st.SaveConfig(cfg); err != nil {
		t.Fatal(err)
	}
	cfg.Capabilities = cfg.Capabilities[:1]
	if changed, err := st.SaveConfig(cfg); err != nil || !changed {
		t.Fatalf("SaveConfig(changed) = %v, %v", changed, err)
	}
	got, err := st.LoadConfig()
	if err != nil || len(got.Capabilities) != 1 {
		t.Fatalf("LoadConfig() = %+v, %v", got, err)
	}
	entries, _ := os.ReadDir(healthPath(root))
	if len(entries) != 1 {
		t.Fatalf("health dir has %d entries, want only the config", len(entries))
	}
}

func TestConfigFailedReplaceKeepsOldFile(t *testing.T) {
	st, root := newProject(t)
	cfg := validConfig(t)
	if _, err := st.SaveConfig(cfg); err != nil {
		t.Fatal(err)
	}
	readOnlyOrSkip(t, healthPath(root))
	cfg.Capabilities = cfg.Capabilities[:1]
	if _, err := st.SaveConfig(cfg); err == nil {
		t.Fatal("SaveConfig() into a read-only directory succeeded")
	}
	got, err := st.LoadConfig()
	if err != nil || len(got.Capabilities) != len(validConfig(t).Capabilities) {
		t.Fatalf("old config damaged: %+v, %v", got, err)
	}
}

func TestLoadConfigRejectsBadFiles(t *testing.T) {
	tests := map[string]struct {
		content string
		want    error
	}{
		"malformed json":      {`{"version": 1,`, ErrMalformedRecord},
		"unknown field":       {`{"version":1,"capabilities":[],"extra":true}`, ErrMalformedRecord},
		"unsupported version": {`{"version":2,"capabilities":[]}`, ErrUnsupportedVersion},
		"trailing data":       {`{"version":1,"capabilities":[]} {}`, ErrMalformedRecord},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			st, root := newProject(t)
			path := healthPath(root, configFile)
			writeFile(t, path, tc.content)
			if _, err := st.LoadConfig(); !errors.Is(err, tc.want) {
				t.Fatalf("LoadConfig() = %v, want %v", err, tc.want)
			}
			if data, _ := os.ReadFile(path); string(data) != tc.content {
				t.Fatal("loading rewrote a malformed file")
			}
		})
	}
}

func TestSaveConfigRefusesInvalidWithoutWriting(t *testing.T) {
	st, root := newProject(t)
	cfg := validConfig(t)
	cfg.Version = 9
	if _, err := st.SaveConfig(cfg); !errors.Is(err, ErrUnsupportedVersion) {
		t.Fatalf("SaveConfig() = %v, want ErrUnsupportedVersion", err)
	}
	if _, err := os.Lstat(healthPath(root)); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("a refused save created the health directory")
	}
}

func TestSaveRefusesNonProjectWithoutScaffolding(t *testing.T) {
	root := t.TempDir()
	st := NewStore(root)
	if _, err := st.SaveSnapshot(snapAt(t, OriginManual, 1)); !errors.Is(err, ErrNotProject) {
		t.Fatalf("SaveSnapshot() = %v, want ErrNotProject", err)
	}
	if _, err := st.SaveConfig(validConfig(t)); !errors.Is(err, ErrNotProject) {
		t.Fatalf("SaveConfig() = %v, want ErrNotProject", err)
	}
	if entries, _ := os.ReadDir(root); len(entries) != 0 {
		t.Fatalf("refused saves left %d entries", len(entries))
	}
	missing := NewStore(filepath.Join(root, "absent"))
	if _, err := missing.SaveSnapshot(snapAt(t, OriginManual, 1)); !errors.Is(err, ErrNotProject) {
		t.Fatalf("missing root: SaveSnapshot() = %v, want ErrNotProject", err)
	}
}

func TestLoadSnapshotsAbsentDirectories(t *testing.T) {
	root := t.TempDir()
	for _, st := range []Store{NewStore(root), NewStore(filepath.Join(root, "absent"))} {
		if got, err := st.LoadSnapshots(); err != nil || got != nil {
			t.Fatalf("LoadSnapshots() = %v, %v; want empty", got, err)
		}
	}
	st, root := newProject(t)
	if got, err := st.LoadSnapshots(); err != nil || got != nil {
		t.Fatalf("no health dir: LoadSnapshots() = %v, %v; want empty", got, err)
	}
	if _, err := os.Lstat(healthPath(root)); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("loading created the health directory")
	}
}

func TestSaveSnapshotIsIdempotent(t *testing.T) {
	st, root := newProject(t)
	s := snapAt(t, OriginManual, 1)
	if created, err := st.SaveSnapshot(s); err != nil || !created {
		t.Fatalf("first SaveSnapshot() = %v, %v; want created", created, err)
	}
	path := healthPath(root, snapshotsDir, snapshotFileName(s.ID))
	before, _ := os.Stat(path)
	if created, err := st.SaveSnapshot(s); err != nil || created {
		t.Fatalf("repeat SaveSnapshot() = %v, %v; want unchanged", created, err)
	}
	after, _ := os.Stat(path)
	if !after.ModTime().Equal(before.ModTime()) || len(snapshotFiles(t, root)) != 1 {
		t.Fatal("repeat save changed the store")
	}
	got, err := st.LoadSnapshots()
	if err != nil || len(got) != 1 || got[0].ID != s.ID {
		t.Fatalf("LoadSnapshots() = %v, %v", got, err)
	}
}

func TestSaveSnapshotRefusesInvalid(t *testing.T) {
	st, root := newProject(t)
	s := snapAt(t, OriginManual, 1)
	s.Results[0].Value.Number = 1 // content no longer matches ID
	if _, err := st.SaveSnapshot(s); !errors.Is(err, ErrIdentityMismatch) {
		t.Fatalf("SaveSnapshot() = %v, want ErrIdentityMismatch", err)
	}
	s.ID = "sha256:../../escape"
	if _, err := st.SaveSnapshot(s); !errors.Is(err, ErrMalformedIdentity) {
		t.Fatalf("SaveSnapshot(traversing ID) = %v, want ErrMalformedIdentity", err)
	}
	if _, err := os.Lstat(healthPath(root)); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("a refused save created the health directory")
	}
}

func TestSaveSnapshotConflictKeepsExistingFile(t *testing.T) {
	tests := map[string]string{
		"corrupt content":     `{"version":1`,
		"different snapshot":  "",
		"unsupported version": `{"version":2}`,
	}
	for name, content := range tests {
		t.Run(name, func(t *testing.T) {
			st, root := newProject(t)
			s := snapAt(t, OriginManual, 1)
			if content == "" {
				other := snapAt(t, OriginManual, 2)
				data, _ := encodeRecord(other)
				content = string(data)
			}
			path := healthPath(root, snapshotsDir, snapshotFileName(s.ID))
			writeFile(t, path, content)
			if _, err := st.SaveSnapshot(s); !errors.Is(err, ErrSnapshotConflict) {
				t.Fatalf("SaveSnapshot() = %v, want ErrSnapshotConflict", err)
			}
			if data, _ := os.ReadFile(path); string(data) != content {
				t.Fatal("conflicting save replaced the existing file")
			}
		})
	}
}

func TestInterruptedWriteLeftoverIsIgnoredAndRetrySucceeds(t *testing.T) {
	st, root := newProject(t)
	writeFile(t, healthPath(root, snapshotsDir, tempPrefix+"123"+tempSuffix), `{"version":1`)
	if got, err := st.LoadSnapshots(); err != nil || len(got) != 0 {
		t.Fatalf("LoadSnapshots() with leftover = %v, %v", got, err)
	}
	mustSave(t, st, snapAt(t, OriginManual, 1))
	if got, err := st.LoadSnapshots(); err != nil || len(got) != 1 {
		t.Fatalf("LoadSnapshots() after retry = %v, %v", got, err)
	}
}

func TestFailedSnapshotWriteLeavesStoreValid(t *testing.T) {
	st, root := newProject(t)
	first := snapAt(t, OriginManual, 1)
	mustSave(t, st, first)
	readOnlyOrSkip(t, healthPath(root, snapshotsDir))
	if _, err := st.SaveSnapshot(snapAt(t, OriginManual, 2)); err == nil {
		t.Fatal("SaveSnapshot() into a read-only directory succeeded")
	}
	got, err := st.LoadSnapshots()
	if err != nil || len(got) != 1 || got[0].ID != first.ID {
		t.Fatalf("store damaged by failed write: %v, %v", got, err)
	}
}

func TestLoadSnapshotsOrderIsDeterministic(t *testing.T) {
	st, _ := newProject(t)
	a := snapAt(t, OriginManual, 5)
	b := snapAt(t, OriginOfficial, 1)
	// Same time, different commit: only the identity can break the tie.
	c, d := snapAt(t, OriginManual, 3), snapAt(t, OriginManual, 3)
	d.Repository.Commit = fmt.Sprintf("%040x", 99)
	d = reseal(d)
	for _, s := range []Snapshot{a, d, b, c} {
		mustSave(t, st, s)
	}
	got, err := st.LoadSnapshots()
	if err != nil || len(got) != 4 {
		t.Fatalf("LoadSnapshots() = %d snapshots, %v", len(got), err)
	}
	tied := []string{c.ID, d.ID}
	slices.Sort(tied)
	want := []string{b.ID, tied[0], tied[1], a.ID}
	for i, s := range got {
		if s.ID != want[i] {
			t.Fatalf("position %d = %s, want %s", i, s.ID, want[i])
		}
	}
}

func TestLoadSnapshotsRejectsBadFiles(t *testing.T) {
	valid := snapAt(t, OriginManual, 1)
	validData, _ := encodeRecord(valid)
	tests := map[string]struct {
		name, content string
		want          error
	}{
		"malformed json":      {snapshotFileName(valid.ID), `{"version":1`, ErrMalformedRecord},
		"unsupported version": {snapshotFileName(valid.ID), `{"version":2}`, ErrUnsupportedVersion},
		"name mismatch":       {snapshotFileName(snapAt(t, OriginManual, 2).ID), string(validData), ErrIdentityMismatch},
		"not a snapshot name": {"notes.txt", "hello", ErrMalformedRecord},
		"uppercase digest":    {"ABCDEF" + snapshotFileName(valid.ID)[6:], string(validData), ErrMalformedRecord},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			st, root := newProject(t)
			mustSave(t, st, snapAt(t, OriginOfficial, 7))
			path := healthPath(root, snapshotsDir, tc.name)
			writeFile(t, path, tc.content)
			if _, err := st.LoadSnapshots(); !errors.Is(err, tc.want) {
				t.Fatalf("LoadSnapshots() = %v, want %v", err, tc.want)
			}
			if data, _ := os.ReadFile(path); string(data) != tc.content {
				t.Fatal("loading rewrote a bad file")
			}
		})
	}
}

func TestLoadSnapshotsRejectsDuplicateContentUnderAnotherName(t *testing.T) {
	st, root := newProject(t)
	s := snapAt(t, OriginManual, 1)
	mustSave(t, st, s)
	data, _ := os.ReadFile(healthPath(root, snapshotsDir, snapshotFileName(s.ID)))
	writeFile(t, healthPath(root, snapshotsDir, "copy.json"), string(data))
	if _, err := st.LoadSnapshots(); err == nil {
		t.Fatal("a second copy of one identity loaded")
	}
}

func TestLoadSnapshotsRejectsDirectoryEntry(t *testing.T) {
	st, root := newProject(t)
	s := snapAt(t, OriginManual, 1)
	if err := os.MkdirAll(healthPath(root, snapshotsDir, snapshotFileName(s.ID)), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := st.LoadSnapshots(); !errors.Is(err, ErrUnsafePath) {
		t.Fatalf("LoadSnapshots() = %v, want ErrUnsafePath", err)
	}
}

func TestSymlinkEscapesAreRefused(t *testing.T) {
	t.Run("health directory", func(t *testing.T) {
		st, root := newProject(t)
		outside := t.TempDir()
		symlinkOrSkip(t, outside, healthPath(root))
		if _, err := st.SaveSnapshot(snapAt(t, OriginManual, 1)); !errors.Is(err, ErrUnsafePath) {
			t.Fatalf("SaveSnapshot() = %v, want ErrUnsafePath", err)
		}
		if _, err := st.SaveConfig(validConfig(t)); !errors.Is(err, ErrUnsafePath) {
			t.Fatalf("SaveConfig() = %v, want ErrUnsafePath", err)
		}
		if _, err := st.LoadSnapshots(); !errors.Is(err, ErrUnsafePath) {
			t.Fatalf("LoadSnapshots() = %v, want ErrUnsafePath", err)
		}
		if entries, _ := os.ReadDir(outside); len(entries) != 0 {
			t.Fatal("a write escaped through the symlink")
		}
	})
	t.Run("snapshots directory", func(t *testing.T) {
		st, root := newProject(t)
		if err := os.MkdirAll(healthPath(root), 0o755); err != nil {
			t.Fatal(err)
		}
		outside := t.TempDir()
		symlinkOrSkip(t, outside, healthPath(root, snapshotsDir))
		if _, err := st.SaveSnapshot(snapAt(t, OriginManual, 1)); !errors.Is(err, ErrUnsafePath) {
			t.Fatalf("SaveSnapshot() = %v, want ErrUnsafePath", err)
		}
		if _, err := st.Prune(); !errors.Is(err, ErrUnsafePath) {
			t.Fatalf("Prune() = %v, want ErrUnsafePath", err)
		}
	})
	t.Run("config file", func(t *testing.T) {
		st, root := newProject(t)
		outside := filepath.Join(t.TempDir(), "outside.json")
		writeFile(t, outside, `{"version":1,"capabilities":[]}`)
		if err := os.MkdirAll(healthPath(root), 0o755); err != nil {
			t.Fatal(err)
		}
		symlinkOrSkip(t, outside, healthPath(root, configFile))
		if _, err := st.LoadConfig(); !errors.Is(err, ErrUnsafePath) {
			t.Fatalf("LoadConfig() = %v, want ErrUnsafePath", err)
		}
		if _, err := st.SaveConfig(validConfig(t)); !errors.Is(err, ErrUnsafePath) {
			t.Fatalf("SaveConfig() = %v, want ErrUnsafePath", err)
		}
		if data, _ := os.ReadFile(outside); string(data) != `{"version":1,"capabilities":[]}` {
			t.Fatal("SaveConfig wrote through the symlink")
		}
	})
	t.Run("snapshot file", func(t *testing.T) {
		st, root := newProject(t)
		s := snapAt(t, OriginManual, 1)
		data, _ := encodeRecord(s)
		outside := filepath.Join(t.TempDir(), "outside.json")
		writeFile(t, outside, string(data))
		if err := os.MkdirAll(healthPath(root, snapshotsDir), 0o755); err != nil {
			t.Fatal(err)
		}
		symlinkOrSkip(t, outside, healthPath(root, snapshotsDir, snapshotFileName(s.ID)))
		if _, err := st.LoadSnapshots(); !errors.Is(err, ErrUnsafePath) {
			t.Fatalf("LoadSnapshots() = %v, want ErrUnsafePath", err)
		}
	})
}

func TestIncompatibleSeriesStayApart(t *testing.T) {
	st, _ := newProject(t)
	a := snapAt(t, OriginOfficial, 1)
	b := snapAt(t, OriginOfficial, 2)
	b.Results[0].Provenance.ProviderVersion = "99.0"
	b = reseal(b)
	if a.Results[0].SeriesID() == b.Results[0].SeriesID() {
		t.Fatal("fixture did not produce two series")
	}
	mustSave(t, st, a)
	mustSave(t, st, b)
	got, err := st.LoadSnapshots()
	if err != nil || len(got) != 2 {
		t.Fatalf("LoadSnapshots() = %d, %v; want both series kept", len(got), err)
	}
	capability := a.Results[0].Capability
	series := func(s Snapshot) string {
		for _, r := range s.Results {
			if r.Capability == capability {
				return r.SeriesID()
			}
		}
		t.Fatalf("snapshot lost %s result", capability)
		return ""
	}
	if series(got[0]) == series(got[1]) {
		t.Fatal("loading merged two incompatible series")
	}
}

func saveN(t *testing.T, st Store, origin Origin, from, count int) []string {
	t.Helper()
	var ids []string
	for i := from; i < from+count; i++ {
		s := snapAt(t, origin, i)
		mustSave(t, st, s)
		ids = append(ids, s.ID)
	}
	return ids
}

func TestPruneRetentionBoundaries(t *testing.T) {
	for _, n := range []int{0, 9, 10, 11, 14} {
		t.Run(fmt.Sprintf("%d manual", n), func(t *testing.T) {
			st, root := newProject(t)
			ids := saveN(t, st, OriginManual, 1, n)
			wantRemoved := max(0, n-ManualRetention)
			plan, err := st.PlanPrune()
			if err != nil || len(plan.Removed) != wantRemoved || len(plan.Retained) != n-wantRemoved {
				t.Fatalf("PlanPrune() = %+v, %v; want %d removed", plan, err, wantRemoved)
			}
			if n > 0 && len(snapshotFiles(t, root)) != n {
				t.Fatal("planning removed files")
			}
			done, err := st.Prune()
			if err != nil || !slices.Equal(done.Removed, plan.Removed) || !slices.Equal(done.Retained, plan.Retained) {
				t.Fatalf("Prune() = %+v, %v; want %+v", done, err, plan)
			}
			if !slices.Equal(done.Removed, ids[:wantRemoved]) {
				t.Fatalf("removed %v, want the %d oldest %v", done.Removed, wantRemoved, ids[:wantRemoved])
			}
			if n > 0 {
				if got := len(snapshotFiles(t, root)); got != n-wantRemoved {
					t.Fatalf("%d files remain, want %d", got, n-wantRemoved)
				}
			}
		})
	}
}

func TestPruneNeverTouchesOfficialSnapshots(t *testing.T) {
	st, root := newProject(t)
	official := saveN(t, st, OriginOfficial, 100, 12)
	manual := saveN(t, st, OriginManual, 1, 13)
	done, err := st.Prune()
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(done.Removed, manual[:3]) {
		t.Fatalf("Removed = %v, want the three oldest manual %v", done.Removed, manual[:3])
	}
	for _, id := range official {
		if !slices.Contains(done.Retained, id) {
			t.Fatalf("official snapshot %s was not retained", id)
		}
		if _, err := os.Stat(healthPath(root, snapshotsDir, snapshotFileName(id))); err != nil {
			t.Fatalf("official snapshot file gone: %v", err)
		}
	}
	if len(snapshotFiles(t, root)) != 12+ManualRetention {
		t.Fatalf("unexpected file count %d", len(snapshotFiles(t, root)))
	}
}

func TestPruneTiesUseIdentity(t *testing.T) {
	st, _ := newProject(t)
	var same []Snapshot
	for i := range ManualRetention + 2 {
		s := snapAt(t, OriginManual, 1) // identical time for all
		s.Repository.Commit = fmt.Sprintf("%040x", 500+i)
		s = reseal(s)
		mustSave(t, st, s)
		same = append(same, s)
	}
	ids := make([]string, len(same))
	for i, s := range same {
		ids[i] = s.ID
	}
	slices.Sort(ids)
	done, err := st.Prune()
	if err != nil || !slices.Equal(done.Removed, ids[:2]) {
		t.Fatalf("Prune() = %+v, %v; want two lowest identities %v", done, err, ids[:2])
	}
}

func TestPruneIsRepeatable(t *testing.T) {
	st, _ := newProject(t)
	saveN(t, st, OriginManual, 1, 12)
	if _, err := st.Prune(); err != nil {
		t.Fatal(err)
	}
	again, err := st.Prune()
	if err != nil || len(again.Removed) != 0 || len(again.Retained) != ManualRetention {
		t.Fatalf("second Prune() = %+v, %v; want nothing to remove", again, err)
	}
}

func TestSavingNeverPrunes(t *testing.T) {
	st, root := newProject(t)
	saveN(t, st, OriginManual, 1, ManualRetention+5)
	if got := len(snapshotFiles(t, root)); got != ManualRetention+5 {
		t.Fatalf("%d files after saving, want %d", got, ManualRetention+5)
	}
}

func TestPruneRefusesWhenHistoryIsInvalid(t *testing.T) {
	st, root := newProject(t)
	saveN(t, st, OriginManual, 1, 12)
	writeFile(t, healthPath(root, snapshotsDir, "stray.txt"), "x")
	if _, err := st.Prune(); !errors.Is(err, ErrMalformedRecord) {
		t.Fatalf("Prune() = %v, want ErrMalformedRecord", err)
	}
	if got := len(snapshotFiles(t, root)); got != 13 {
		t.Fatalf("refused prune changed the store: %d files", got)
	}
}

func TestPrunePartialFailureKeepsRemainingValid(t *testing.T) {
	st, root := newProject(t)
	ids := saveN(t, st, OriginManual, 1, 12)
	readOnlyOrSkip(t, healthPath(root, snapshotsDir))
	done, err := st.Prune()
	if err == nil {
		t.Fatal("Prune() in a read-only directory succeeded")
	}
	if len(done.Removed) != 0 || len(done.Retained) != 12 {
		t.Fatalf("report after failure = %+v", done)
	}
	got, lerr := st.LoadSnapshots()
	if lerr != nil || len(got) != len(ids) {
		t.Fatalf("store damaged by failed prune: %d, %v", len(got), lerr)
	}
}
