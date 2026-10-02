package codehealth

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"testing"
)

func TestOlderAndNewerRecordVersionsAreRefused(t *testing.T) {
	for _, version := range []string{"0", "2", "-1"} {
		t.Run("config "+version, func(t *testing.T) {
			st, root := newProject(t)
			content := `{"version":` + version + `,"capabilities":[]}`
			writeFile(t, healthPath(root, configFile), content)
			if _, err := st.LoadConfig(); !errors.Is(err, ErrUnsupportedVersion) {
				t.Fatalf("LoadConfig() = %v, want ErrUnsupportedVersion", err)
			}
		})
		t.Run("snapshot "+version, func(t *testing.T) {
			st, root := newProject(t)
			s := snapAt(t, OriginManual, 1)
			writeFile(t, healthPath(root, snapshotsDir, snapshotFileName(s.ID)), `{"version":`+version+`}`)
			if _, err := st.LoadSnapshots(); !errors.Is(err, ErrUnsupportedVersion) {
				t.Fatalf("LoadSnapshots() = %v, want ErrUnsupportedVersion", err)
			}
		})
	}
}

func TestOversizedRecordsAreRefusedUnread(t *testing.T) {
	st, root := newProject(t)
	big := strings.Repeat(" ", maxRecordBytes+1)
	writeFile(t, healthPath(root, configFile), big)
	if _, err := st.LoadConfig(); !errors.Is(err, ErrUnboundedDetail) {
		t.Fatalf("LoadConfig() = %v, want ErrUnboundedDetail", err)
	}
	s := snapAt(t, OriginManual, 1)
	writeFile(t, healthPath(root, snapshotsDir, snapshotFileName(s.ID)), big)
	if _, err := st.LoadSnapshots(); !errors.Is(err, ErrUnboundedDetail) {
		t.Fatalf("LoadSnapshots() = %v, want ErrUnboundedDetail", err)
	}
	if _, err := st.SaveSnapshot(s); !errors.Is(err, ErrSnapshotConflict) && !errors.Is(err, ErrUnboundedDetail) {
		t.Fatalf("SaveSnapshot() over oversized file = %v", err)
	}
}

// A destination the rename cannot replace fails after the temporary file was
// created: the destination must be untouched and no temporary file may remain.
func TestFailedReplacementAfterTempCreationLeavesNoTemp(t *testing.T) {
	t.Run("config", func(t *testing.T) {
		st, root := newProject(t)
		blocker := healthPath(root, configFile, "keep")
		writeFile(t, blocker, "owner bytes")
		if _, err := st.SaveConfig(validConfig(t)); err == nil {
			t.Fatal("SaveConfig() over a directory succeeded")
		}
		assertNoTemps(t, healthPath(root))
		if got, _ := os.ReadFile(blocker); string(got) != "owner bytes" {
			t.Fatalf("blocker changed to %q", got)
		}
	})
	t.Run("report", func(t *testing.T) {
		st, root := newProject(t)
		blocker := healthPath(root, reportFile, "keep")
		writeFile(t, blocker, "owner bytes")
		if _, err := st.WriteReport("text\n"); err == nil {
			t.Fatal("WriteReport() over a directory succeeded")
		}
		assertNoTemps(t, healthPath(root))
		if got, _ := os.ReadFile(blocker); string(got) != "owner bytes" {
			t.Fatalf("blocker changed to %q", got)
		}
	})
}

func assertNoTemps(t *testing.T, dir string) {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if isTempName(e.Name()) {
			t.Errorf("leftover temporary file %s in %s", e.Name(), dir)
		}
	}
}

// Failed replacements keep the owner's bytes and the prior valid report, and a
// repeat after the obstruction is removed succeeds.
func TestFailedWritesPreserveOwnerBytesAndRepeatSucceeds(t *testing.T) {
	st, root := newProject(t)
	owned := "# mine\nsnapshots/\n"
	writeFile(t, healthPath(root, gitignoreFile), owned)
	if _, err := st.WriteReport("good report\n"); err != nil {
		t.Fatal(err)
	}
	cfgBytes, _ := encodeRecord(validConfig(t))
	writeFile(t, healthPath(root, configFile), string(cfgBytes))

	readOnlyOrSkip(t, healthPath(root))
	cfg := validConfig(t)
	cfg.Capabilities = cfg.Capabilities[:1]
	if _, err := st.SaveConfig(cfg); err == nil {
		t.Fatal("SaveConfig() into a read-only directory succeeded")
	}
	if _, err := st.WriteReport("new report\n"); err == nil {
		t.Fatal("WriteReport() into a read-only directory succeeded")
	}
	if got, _ := os.ReadFile(healthPath(root, gitignoreFile)); string(got) != owned {
		t.Errorf(".gitignore = %q, want the owner's bytes", got)
	}
	if got, _ := os.ReadFile(healthPath(root, configFile)); string(got) != string(cfgBytes) {
		t.Errorf("config bytes changed by a failed save")
	}
	if got := readReport(t, root); got != "good report\n" {
		t.Errorf("report = %q, want the prior valid report", got)
	}
	assertNoTemps(t, healthPath(root))

	if err := os.Chmod(healthPath(root), 0o755); err != nil {
		t.Fatal(err)
	}
	if changed, err := st.WriteReport("new report\n"); err != nil || !changed {
		t.Fatalf("repeat WriteReport() = %v, %v", changed, err)
	}
	if changed, err := st.SaveConfig(cfg); err != nil || !changed {
		t.Fatalf("repeat SaveConfig() = %v, %v", changed, err)
	}
}

func TestConcurrentSnapshotSavesAndReads(t *testing.T) {
	st, root := newProject(t)
	same := snapAt(t, OriginOfficial, 1)
	var others []Snapshot
	for i := 2; i < 10; i++ {
		others = append(others, snapAt(t, OriginManual, i))
	}
	var wg sync.WaitGroup
	errs := make(chan error, 64)
	created := make(chan bool, 16)
	for range 8 {
		wg.Go(func() {
			ok, err := st.SaveSnapshot(same)
			if err != nil {
				errs <- err
			}
			created <- ok
		})
	}
	for _, s := range others {
		wg.Go(func() {
			if _, err := st.SaveSnapshot(s); err != nil {
				errs <- err
			}
		})
	}
	for range 8 {
		wg.Go(func() {
			if _, err := st.LoadSnapshots(); err != nil {
				errs <- err
			}
		})
	}
	wg.Wait()
	close(errs)
	close(created)
	for err := range errs {
		t.Errorf("concurrent operation failed: %v", err)
	}
	wins := 0
	for ok := range created {
		if ok {
			wins++
		}
	}
	if wins != 1 {
		t.Errorf("%d writers reported creating the same snapshot, want exactly 1", wins)
	}
	got, err := st.LoadSnapshots()
	if err != nil || len(got) != 1+len(others) {
		t.Fatalf("LoadSnapshots() = %d snapshots, %v; want %d", len(got), err, 1+len(others))
	}
	assertNoTemps(t, healthPath(root, snapshotsDir))
}

// platformRefusal reports whether err is Windows refusing to replace or open a
// file another handle holds: access denied or a sharing violation. Elsewhere
// nothing is a platform refusal.
func platformRefusal(err error) bool {
	if runtime.GOOS != "windows" {
		return false
	}
	var errno syscall.Errno
	return errors.Is(err, fs.ErrPermission) || (errors.As(err, &errno) && (errno == 5 || errno == 32))
}

// Readers during replacement see a complete old or new file, never a mixture.
// Windows may refuse a replacement while a reader holds the file; that refusal
// is accepted, but the owner's bytes must stay complete and no temp file may
// remain.
func TestReadersSeeCompleteConfigAndReportDuringReplacement(t *testing.T) {
	st, root := newProject(t)
	full := validConfig(t)
	short := full
	short.Capabilities = full.Capabilities[:1]
	if _, err := st.SaveConfig(full); err != nil {
		t.Fatal(err)
	}
	oldText, newText := strings.Repeat("old\n", 4000), strings.Repeat("new\n", 4000)
	if _, err := st.WriteReport(oldText); err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	errs := make(chan error, 8)
	wg.Go(func() {
		for i := range 100 {
			cfg, text := full, newText
			if i%2 == 1 {
				cfg, text = short, oldText
			}
			if _, err := st.SaveConfig(cfg); err != nil && !platformRefusal(err) {
				errs <- err
				return
			}
			if _, err := st.WriteReport(text); err != nil && !platformRefusal(err) {
				errs <- err
				return
			}
		}
	})
	for range 4 {
		wg.Go(func() {
			for range 200 {
				cfg, err := st.LoadConfig()
				if err != nil {
					if platformRefusal(err) {
						continue
					}
					errs <- err
					return
				}
				if n := len(cfg.Capabilities); n != len(full.Capabilities) && n != 1 {
					errs <- errors.New("config read was neither the old nor the new content")
					return
				}
				data, err := os.ReadFile(filepath.Join(healthPath(root), reportFile))
				if err != nil {
					if platformRefusal(err) {
						continue
					}
					errs <- err
					return
				}
				if s := string(data); s != oldText && s != newText {
					errs <- errors.New("report read was neither the old nor the new content")
					return
				}
			}
		})
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Error(err)
	}
	if cfg, err := st.LoadConfig(); err != nil || (len(cfg.Capabilities) != len(full.Capabilities) && len(cfg.Capabilities) != 1) {
		t.Errorf("config after replacement = %d capabilities, %v; want the old or new content", len(cfg.Capabilities), err)
	}
	if data, err := os.ReadFile(filepath.Join(healthPath(root), reportFile)); err != nil || (string(data) != oldText && string(data) != newText) {
		t.Errorf("report after replacement is not the old or new content: %v", err)
	}
	assertNoTemps(t, healthPath(root))
}

func TestPruneToleratesAFileRemovedByAnotherPrune(t *testing.T) {
	st, root := newProject(t)
	ids := saveN(t, st, OriginManual, 1, 12)
	var wg sync.WaitGroup
	errs := make(chan error, 2)
	for range 2 {
		wg.Go(func() {
			if _, err := st.Prune(); err != nil && !errors.Is(err, ErrMalformedRecord) {
				errs <- err
			}
		})
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Errorf("concurrent Prune() = %v", err)
	}
	if got := len(snapshotFiles(t, root)); got != ManualRetention {
		t.Errorf("%d files remain, want %d", got, ManualRetention)
	}
	for _, id := range ids[2:] {
		if _, err := os.Stat(healthPath(root, snapshotsDir, snapshotFileName(id))); err != nil {
			t.Errorf("retained snapshot missing: %v", err)
		}
	}
}
