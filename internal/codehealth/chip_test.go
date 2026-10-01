package codehealth

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func saveConfig(t *testing.T, st Store) Config {
	t.Helper()
	cfg := validConfig(t)
	if _, err := st.SaveConfig(cfg); err != nil {
		t.Fatalf("SaveConfig() = %v", err)
	}
	return cfg
}

func TestLoadChipNotSetUpWithoutConfig(t *testing.T) {
	st, _ := newProject(t)
	if got := LoadChip(st.projectPath); got.State != ChipNotSetUp {
		t.Fatalf("LoadChip() = %+v, want not set up", got)
	}
}

func TestLoadChipNoCheckYetWithConfigAndNoSnapshots(t *testing.T) {
	st, root := newProject(t)
	saveConfig(t, st)
	if got := LoadChip(root); got.State != ChipNoCheck {
		t.Fatalf("LoadChip() = %+v, want no check yet", got)
	}
}

func TestLoadChipMeasuredCountsGoodOfSignals(t *testing.T) {
	st, root := newProject(t)
	cfg := saveConfig(t, st)
	mustSave(t, st, snapAt(t, OriginOfficial, 1))

	got := LoadChip(root)
	want := Chip{State: ChipMeasured, Overall: ClassificationNeedsAttention, Label: "Needs Attention", Good: 1, Signals: signalCount(cfg)}
	if got != want {
		t.Fatalf("LoadChip() = %+v, want %+v", got, want)
	}
	// Three configured capabilities plus a placeholder each for the other two.
	if got.Signals != 5 {
		t.Errorf("Signals = %d, want 5", got.Signals)
	}
}

func TestLoadChipUnreadableStorageIsNotSetUp(t *testing.T) {
	st, root := newProject(t)
	saveConfig(t, st)
	writeFile(t, healthPath(root, snapshotsDir, "not-a-snapshot.json"), "{")
	if got := LoadChip(root); got.State != ChipNotSetUp {
		t.Fatalf("LoadChip() = %+v, want not set up for damaged storage", got)
	}
	if got := LoadChip(filepath.Join(root, "missing")); got.State != ChipNotSetUp {
		t.Fatalf("LoadChip(missing dir) = %+v, want not set up", got)
	}
}

func TestLatestSnapshotReturnsNewestWithoutLoadingTheOthers(t *testing.T) {
	st, root := newProject(t)
	older, newer := snapAt(t, OriginOfficial, 1), snapAt(t, OriginManual, 2)
	mustSave(t, st, older)
	mustSave(t, st, newer)
	base := healthPath(root, snapshotsDir)
	if err := os.Chtimes(filepath.Join(base, snapshotFileName(older.ID)), fixedTime(1), fixedTime(1)); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(filepath.Join(base, snapshotFileName(newer.ID)), fixedTime(2), fixedTime(2)); err != nil {
		t.Fatal(err)
	}
	// A damaged older file must not matter: only the newest is decoded.
	writeFile(t, filepath.Join(base, snapshotFileName(older.ID)), "{")
	if err := os.Chtimes(filepath.Join(base, snapshotFileName(older.ID)), fixedTime(1), fixedTime(1)); err != nil {
		t.Fatal(err)
	}

	got, found, err := st.LatestSnapshot()
	if err != nil || !found || got.ID != newer.ID {
		t.Fatalf("LatestSnapshot() = %s, %v, %v; want %s", got.ID, found, err, newer.ID)
	}
}

func TestLatestSnapshotAbsentStorage(t *testing.T) {
	st, _ := newProject(t)
	if _, found, err := st.LatestSnapshot(); found || err != nil {
		t.Fatalf("LatestSnapshot() found=%v err=%v, want neither", found, err)
	}
}

func TestDashboardChipMatchesStates(t *testing.T) {
	if got := (Dashboard{State: DashboardNotConfigured}).Chip(); got.State != ChipNotSetUp {
		t.Errorf("not configured chip = %+v", got)
	}
	if got := (Dashboard{State: DashboardFirstRun}).Chip(); got.State != ChipNoCheck {
		t.Errorf("first run chip = %+v", got)
	}
	d := Dashboard{State: DashboardMeasured, Overall: ClassificationWatch, OverallText: "Watch", Rows: []DashboardRow{
		{Label: ClassificationGood}, {Label: ClassificationGood}, {Label: ClassificationWatch}, {Label: ClassificationUnknown}, {Label: ClassificationGood},
	}}
	want := Chip{State: ChipMeasured, Overall: ClassificationWatch, Label: "Watch", Good: 3, Signals: 5}
	if got := d.Chip(); got != want {
		t.Errorf("measured chip = %+v, want %+v", got, want)
	}
}

func fixedTime(n int) time.Time {
	return time.Date(2026, 10, 1, 0, n, 0, 0, time.UTC)
}

func setMTimes(t *testing.T, root string, times map[string]time.Time) {
	t.Helper()
	for id, at := range times {
		if err := os.Chtimes(filepath.Join(healthPath(root, snapshotsDir), snapshotFileName(id)), at, at); err != nil {
			t.Fatal(err)
		}
	}
}

// The header and the dashboard must name the same newest snapshot whatever the
// files' modification times say: a restored or copied history reorders them.
func TestLatestSnapshotFollowsStoredRecencyNotFileTime(t *testing.T) {
	tie := func(n int, commit int) Snapshot {
		s := snapAt(t, OriginOfficial, n)
		s.Repository.Commit = fmt.Sprintf("%040x", commit)
		return reseal(s)
	}
	for name, tc := range map[string]struct {
		snaps  []Snapshot
		mtimes func(s []Snapshot) map[string]time.Time
	}{
		"mtime reversed": {
			[]Snapshot{snapAt(t, OriginOfficial, 1), snapAt(t, OriginManual, 2)},
			func(s []Snapshot) map[string]time.Time {
				return map[string]time.Time{s[0].ID: fixedTime(9), s[1].ID: fixedTime(1)}
			},
		},
		"equal mtimes": {
			[]Snapshot{snapAt(t, OriginManual, 2), snapAt(t, OriginOfficial, 1), snapAt(t, OriginOfficial, 3)},
			func(s []Snapshot) map[string]time.Time {
				return map[string]time.Time{s[0].ID: fixedTime(5), s[1].ID: fixedTime(5), s[2].ID: fixedTime(5)}
			},
		},
		"same stored time": {
			[]Snapshot{tie(4, 1), tie(4, 2), tie(4, 3)},
			func(s []Snapshot) map[string]time.Time {
				return map[string]time.Time{s[0].ID: fixedTime(3), s[1].ID: fixedTime(1), s[2].ID: fixedTime(2)}
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			st, root := newProject(t)
			for _, s := range tc.snaps {
				mustSave(t, st, s)
			}
			setMTimes(t, root, tc.mtimes(tc.snaps))
			all, err := st.LoadSnapshots()
			if err != nil {
				t.Fatal(err)
			}
			got, found, err := st.LatestSnapshot()
			if err != nil || !found || got.ID != all[len(all)-1].ID {
				t.Fatalf("LatestSnapshot() = %s, %v, %v; want %s", got.ID, found, err, all[len(all)-1].ID)
			}
			saveConfig(t, st)
			d, err := LoadDashboard(root)
			if err != nil {
				t.Fatal(err)
			}
			if chip, want := LoadChip(root), d.Chip(); chip != want {
				t.Errorf("header chip %+v differs from the dashboard's %+v", chip, want)
			}
		})
	}
}
