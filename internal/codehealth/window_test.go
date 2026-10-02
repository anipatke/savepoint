package codehealth

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

// windowHistory saves size snapshots of the T-094 single-instance fixture: every
// third one manual, so size 14 holds exactly windowOfficial official snapshots
// and size 15 holds one more.
func windowHistory(t *testing.T, size int) (Store, string, Config) {
	t.Helper()
	store, root := newProject(t)
	p := historyProfiles[0]
	cfg := historyConfig(p)
	if _, err := store.SaveConfig(cfg); err != nil {
		t.Fatal(err)
	}
	for n := 0; n < size; n++ {
		mustSave(t, store, historySnapshot(cfg, p, n))
	}
	return store, root, cfg
}

func snapshotIDs(snaps []Snapshot) []string {
	var ids []string
	for _, s := range snaps {
		ids = append(ids, s.ID)
	}
	return ids
}

func TestLoadWindowKeepsTheNewestOfficialAndEverythingNewer(t *testing.T) {
	for _, tc := range []struct {
		size    int
		wantCut bool
	}{{0, false}, {1, false}, {5, false}, {15, false}, {16, true}, {40, true}} {
		store, _, _ := windowHistory(t, tc.size)
		all, err := store.LoadSnapshots()
		if err != nil {
			t.Fatal(err)
		}
		got, cut, err := store.LoadWindow()
		if err != nil {
			t.Fatalf("size %d: %v", tc.size, err)
		}
		if cut != tc.wantCut {
			t.Errorf("size %d: cut = %v, want %v", tc.size, cut, tc.wantCut)
		}
		want := all
		if cut {
			officials, start := 0, 0
			for i := len(all) - 1; i >= 0; i-- {
				if all[i].Origin == OriginOfficial {
					if officials++; officials == windowOfficial {
						start = i
						break
					}
				}
			}
			want = all[start:]
		}
		if !reflect.DeepEqual(snapshotIDs(got), snapshotIDs(want)) {
			t.Errorf("size %d: window = %d snapshots, want the %d newest from the oldest of the %d newest official", tc.size, len(got), len(want), windowOfficial)
		}
	}
}

func TestLoadWindowBreaksTimeTiesByIdentityLikeLoadSnapshots(t *testing.T) {
	store, _ := newProject(t)
	for n := 0; n < 4; n++ {
		s := snapAt(t, OriginOfficial, 0)
		s.Repository = dashRepo(n)
		s.ID = s.ComputeID()
		mustSave(t, store, s)
	}
	all, _ := store.LoadSnapshots()
	got, cut, err := store.LoadWindow()
	if err != nil || cut {
		t.Fatalf("LoadWindow() = %v, cut %v", err, cut)
	}
	if !reflect.DeepEqual(snapshotIDs(got), snapshotIDs(all)) {
		t.Errorf("tied snapshots ordered differently from LoadSnapshots")
	}
}

func TestLoadWindowWithoutHealthStorageIsEmpty(t *testing.T) {
	store, _ := newProject(t)
	snaps, cut, err := store.LoadWindow()
	if err != nil || len(snaps) != 0 || cut {
		t.Errorf("LoadWindow() = %v, %v, %v; want none", snaps, cut, err)
	}
}

// snapshotFileOf finds the file holding snapshot number n of windowHistory.
func snapshotFileOf(t *testing.T, root string, s Snapshot) string {
	t.Helper()
	return filepath.Join(healthPath(root, snapshotsDir), snapshotFileName(s.ID))
}

func TestLoadWindowDamageInsideAndOutsideTheWindow(t *testing.T) {
	store, root, _ := windowHistory(t, 40)
	all, err := store.LoadSnapshots()
	if err != nil {
		t.Fatal(err)
	}
	// Break the body after the head: the time and origin stay readable.
	damage := func(s Snapshot) {
		path := snapshotFileOf(t, root, s)
		data, _ := os.ReadFile(path)
		if err := os.WriteFile(path, []byte(strings.Replace(string(data), `"results"`, `"resultz"`, 1)), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	damage(all[0])
	if _, _, err := store.LoadWindow(); err != nil {
		t.Errorf("a damaged body outside the window failed the load: %v", err)
	}
	if _, err := store.LoadSnapshots(); err == nil {
		t.Error("LoadSnapshots accepted the damaged file; full validation must still see it")
	}
	damage(all[len(all)-1])
	if _, _, err := store.LoadWindow(); err == nil {
		t.Error("a damaged snapshot inside the window loaded")
	}
}

func TestLoadWindowStillFailsOnNamesAndUnreadableHeads(t *testing.T) {
	for name, tc := range map[string]func(t *testing.T, dir string, any Snapshot){
		"not a snapshot name": func(t *testing.T, dir string, _ Snapshot) {
			if err := os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("x"), 0o644); err != nil {
				t.Fatal(err)
			}
		},
		"head that is not JSON": func(t *testing.T, dir string, s Snapshot) {
			if err := os.WriteFile(filepath.Join(dir, snapshotFileName(s.ID)), []byte("garbage"), 0o644); err != nil {
				t.Fatal(err)
			}
		},
		"head with a bad time": func(t *testing.T, dir string, s Snapshot) {
			path := filepath.Join(dir, snapshotFileName(s.ID))
			data, _ := os.ReadFile(path)
			if err := os.WriteFile(path, []byte(strings.Replace(string(data), s.CreatedAt, "yesterday", 1)), 0o644); err != nil {
				t.Fatal(err)
			}
		},
		"directory in place of a file": func(t *testing.T, dir string, s Snapshot) {
			path := filepath.Join(dir, snapshotFileName(s.ID))
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
			if err := os.Mkdir(path, 0o755); err != nil {
				t.Fatal(err)
			}
		},
	} {
		t.Run(name, func(t *testing.T) {
			store, root, _ := windowHistory(t, 40)
			all, _ := store.LoadSnapshots()
			tc(t, healthPath(root, snapshotsDir), all[0]) // outside the window
			if _, _, err := store.LoadWindow(); err == nil {
				t.Error("LoadWindow() succeeded")
			}
		})
	}
}

// A head that does not give time and origin first is decoded whole instead.
func TestLoadWindowDecodesAFileWhoseHeadIsInAnotherOrder(t *testing.T) {
	store, root, _ := windowHistory(t, 40)
	all, _ := store.LoadSnapshots()
	oldest := all[0]
	path := snapshotFileOf(t, root, oldest)
	data, _ := os.ReadFile(path)
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		t.Fatal(err)
	}
	var b strings.Builder
	b.WriteString("{")
	for i, k := range []string{"results", "summary", "repository", "version", "id", "origin", "retention", "created_at"} {
		if i > 0 {
			b.WriteString(",")
		}
		key, _ := json.Marshal(k)
		b.Write(key)
		b.WriteString(":")
		b.Write(fields[k])
	}
	b.WriteString("}")
	if err := os.WriteFile(path, []byte(b.String()), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, _, err := store.LoadWindow(); err != nil {
		t.Errorf("a valid file in another field order failed the load: %v", err)
	}
	// And its damage is found, because the fallback validates it.
	if err := os.WriteFile(path, []byte(strings.Replace(b.String(), `"retention"`, `"retentionx"`, 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, _, err := store.LoadWindow(); err == nil {
		t.Error("a damaged file with an unreadable head loaded")
	}
}

// fullDashboard is what LoadDashboard gave before the window: every snapshot.
func fullDashboard(t *testing.T, root string) Dashboard {
	t.Helper()
	store := NewStore(root)
	cfg, err := store.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	snaps, err := store.LoadSnapshots()
	if err != nil {
		t.Fatal(err)
	}
	return buildDashboard(store, cfg, snaps, false)
}

func TestLoadDashboardIsUnchangedWithinTheWindow(t *testing.T) {
	for _, size := range []int{1, 2, 3, 8, 15} {
		_, root, _ := windowHistory(t, size)
		got, err := LoadDashboard(root)
		if err != nil {
			t.Fatal(err)
		}
		if want := fullDashboard(t, root); !reflect.DeepEqual(got, want) {
			t.Errorf("size %d: dashboard differs from the full-history dashboard", size)
		}
	}
}

func TestLoadDashboardBeyondTheWindowKeepsTrendsAndSaysOlderWasNotRead(t *testing.T) {
	for _, size := range []int{16, 30, 100} {
		_, root, _ := windowHistory(t, size)
		got, err := LoadDashboard(root)
		if err != nil {
			t.Fatal(err)
		}
		want := fullDashboard(t, root)
		if len(got.Rows) != len(want.Rows) {
			t.Fatalf("size %d: %d rows, want %d", size, len(got.Rows), len(want.Rows))
		}
		for i := range got.Rows {
			g, w := got.Rows[i], want.Rows[i]
			if !strings.HasSuffix(g.Basis, textOlderNotRead) || !strings.Contains(g.Basis, "most recent") {
				t.Errorf("size %d row %d: basis %q does not say older history was left unread", size, i, g.Basis)
			}
			if strings.Contains(w.Basis, "not read") {
				t.Errorf("full-history basis %q must not mention the window", w.Basis)
			}
			g.Basis, w.Basis = "", ""
			if !reflect.DeepEqual(g, w) {
				t.Errorf("size %d row %d (%s): window row differs from full history beyond Basis:\n got  %+v\n want %+v", size, i, g.CapabilityText, g, w)
			}
		}
		got.Rows, want.Rows = nil, nil
		if !reflect.DeepEqual(got, want) {
			t.Errorf("size %d: headline, sign-off, history list or chip inputs differ from full history", size)
		}
	}
}

func TestBasisWordsNamesTheWindowOnlyWhenHistoryWasCut(t *testing.T) {
	s := series{values: []float64{1, 2, 3}, manual: 2, incompatible: 1}
	if got, want := basisWords(s), "Compared with 3 earlier comparable official checks. 1 earlier official result not compared. 2 manual results shown, not counted."; got != want {
		t.Errorf("basisWords = %q, want %q", got, want)
	}
	s.cut = true
	if got, want := basisWords(s), "Compared with the 3 most recent comparable official checks. 1 earlier official result not compared. 2 manual results shown, not counted. Older history was not read."; got != want {
		t.Errorf("basisWords (cut) = %q, want %q", got, want)
	}
	if got, want := basisWords(series{cut: true}), "No comparable official history yet. Older history was not read."; got != want {
		t.Errorf("basisWords (cut, empty) = %q, want %q", got, want)
	}
}

func TestLoadDashboardReadsNoSnapshotBodyOutsideTheWindow(t *testing.T) {
	store, root, _ := windowHistory(t, 40)
	all, _ := store.LoadSnapshots()
	for _, s := range all[:10] {
		path := snapshotFileOf(t, root, s)
		data, _ := os.ReadFile(path)
		// Same head, same size, unreadable body.
		bad := []byte(strings.Replace(string(data), `"results"`, `"resultz"`, 1))
		if err := os.WriteFile(path, bad, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	got, err := LoadDashboard(root)
	if err != nil {
		t.Fatalf("LoadDashboard() = %v", err)
	}
	if got.State != DashboardMeasured || len(got.History) != MaxDashboardHistory {
		t.Errorf("dashboard = %v with %d history entries", got.State, len(got.History))
	}
}

func TestLoadDashboardWithOnlyManualHistoryStillSaysNoOfficial(t *testing.T) {
	store, root := newProject(t)
	p := historyProfiles[0]
	cfg := historyConfig(p)
	if _, err := store.SaveConfig(cfg); err != nil {
		t.Fatal(err)
	}
	for n := 0; n < 3; n++ {
		s := historySnapshot(cfg, p, 2) // manual
		s.CreatedAt = historySnapshot(cfg, p, 2+3*n).CreatedAt
		s.ID = s.ComputeID()
		mustSave(t, store, s)
	}
	got, err := LoadDashboard(root)
	if err != nil {
		t.Fatal(err)
	}
	if got.SignOff != textSignOffNoOfficial {
		t.Errorf("SignOff = %q, want %q", got.SignOff, textSignOffNoOfficial)
	}
}

// A history of at most windowOfficial official snapshots is never cut, however
// many manual ones lead, trail or sit among them.
func TestLoadDashboardIsUnchangedWithLeadingManualHistory(t *testing.T) {
	p := historyProfiles[0]
	for name, order := range map[string][]int{
		"leading manuals":     {2, 5, 6, 7, 9, 10, 12, 13, 15, 16, 18, 19},
		"manual-only":         {2, 5, 8},
		"trailing manuals":    {0, 1, 3, 4, 6, 7, 9, 10, 12, 13, 14},
		"interleaved":         {0, 2, 3, 5, 6, 8, 9, 11, 12, 14, 15, 17},
		"ten officials alone": {0, 1, 3, 4, 6, 7, 9, 10, 12, 13},
	} {
		store, root := newProject(t)
		cfg := historyConfig(p)
		if _, err := store.SaveConfig(cfg); err != nil {
			t.Fatal(err)
		}
		for _, n := range order {
			mustSave(t, store, historySnapshot(cfg, p, n))
		}
		if _, cut, err := store.LoadWindow(); err != nil || cut {
			t.Errorf("%s: LoadWindow cut = %v, err = %v; want the whole history", name, cut, err)
		}
		got, err := LoadDashboard(root)
		if err != nil {
			t.Fatal(err)
		}
		if want := fullDashboard(t, root); !reflect.DeepEqual(got, want) {
			t.Errorf("%s: dashboard differs from the full-history dashboard", name)
		}
	}
}

// A repeated created_at or origin is read as the full decoder reads it (the last
// value wins), so the window never files a snapshot under the wrong time.
func TestLoadWindowAgreesWithTheDecoderOnRepeatedHeaderFields(t *testing.T) {
	after := regexp.MustCompile(`"origin":\s*"[a-z]+",`)
	for name, edit := range map[string]func(string) string{
		"created_at first": func(d string) string {
			return strings.Replace(d, "{", `{"created_at":"2000-01-01T00:00:00Z",`, 1)
		},
		"origin first": func(d string) string {
			return strings.Replace(d, "{", `{"origin":"manual",`, 1)
		},
		"created_at after the header fields": func(d string) string {
			return after.ReplaceAllStringFunc(d, func(m string) string { return m + `"created_at":"2000-01-01T00:00:00Z",` })
		},
	} {
		store, root, _ := windowHistory(t, 40)
		all, err := store.LoadSnapshots()
		if err != nil {
			t.Fatal(err)
		}
		path := snapshotFileOf(t, root, all[len(all)-1])
		data, _ := os.ReadFile(path)
		edited := edit(string(data))
		if edited == string(data) {
			t.Fatalf("%s: edit changed nothing", name)
		}
		if err := os.WriteFile(path, []byte(edited), 0o644); err != nil {
			t.Fatal(err)
		}
		full, err := store.LoadSnapshots()
		if err != nil {
			t.Fatalf("%s: full decoder rejected the file: %v", name, err)
		}
		got, _, err := store.LoadWindow()
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if g, w := got[len(got)-1].ID, full[len(full)-1].ID; g != w {
			t.Errorf("%s: window newest = %s, decoder newest = %s", name, g, w)
		}
	}
}
