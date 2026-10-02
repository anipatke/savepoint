package codehealth

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Benchmarks for the cost of health history (T-094). They run only with
// -bench, so they never fail a gate on machine speed. Every fixture is built
// from this file's own constants: no provider, network or developer project.
//
//	go test ./internal/codehealth -run '^$' -bench 'HealthHistory' -benchmem -count=5

// historySizes are the snapshot counts the decision covers.
var historySizes = []int{10, 100, 1000}

// historyProfile names how large each snapshot is: how many instances each
// capability has, and how much evidence each result carries.
type historyProfile struct {
	name      string
	instances int
	evidence  int
	noteLen   int
}

var historyProfiles = []historyProfile{
	{name: "single", instances: 1, evidence: 1, noteLen: 20},
	{name: "multi", instances: 4, evidence: 3, noteLen: 60},
	{name: "heavy", instances: 4, evidence: MaxEvidence, noteLen: MaxNoteLen},
}

// historyConfig has p.instances instances of every capability. A single
// instance is unnamed, as in a normal project; several get distinct scopes,
// which configuration validation requires.
func historyConfig(p historyProfile) Config {
	cfg := Config{Version: ConfigVersion}
	for _, c := range Capabilities() {
		for i := 0; i < p.instances; i++ {
			name, scope := "", "**/*.go"
			if p.instances > 1 {
				name, scope = fmt.Sprintf("inst%d", i), fmt.Sprintf("mod%d/**/*.go", i)
			}
			cfg.Capabilities = append(cfg.Capabilities, CapabilityConfig{
				Capability: c, Provider: dashProviders[c], Name: name, Scope: []string{scope},
			})
		}
	}
	return cfg
}

// historySnapshot is snapshot number n of a deterministic history: every
// third one manual, values drifting with n so trends have something to say.
func historySnapshot(cfg Config, p historyProfile, n int) Snapshot {
	origin, retention := OriginOfficial, RetentionPermanent
	if n%3 == 2 {
		origin, retention = OriginManual, RetentionPrunable
	}
	s := Snapshot{
		Version: SnapshotVersion, Origin: origin, Retention: retention,
		CreatedAt:  time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(n) * time.Minute).Format(timestampLayout),
		Repository: dashRepo(n),
	}
	note := strings.Repeat("n", p.noteLen)
	var assessed []Assessment
	for _, cc := range cfg.Capabilities {
		r := dashResult(cc, goodValues[cc.Capability]+float64(n%7))
		r.Evidence = nil
		for e := 0; e < p.evidence; e++ {
			r.Evidence = append(r.Evidence, EvidenceRef{Path: fmt.Sprintf("pkg/dir%d/file%d.go", e%5, e), Line: e + 1, Note: note})
		}
		a := Assess(origin, r, cc.Thresholds, nil)
		assessed = append(assessed, a)
		s.Results = append(s.Results, r)
		s.Summary.Capabilities = append(s.Summary.Capabilities, a.Summary())
	}
	s.Summary.Overall = Overall(assessed)
	s.ID = s.ComputeID()
	return s
}

// historyFixture writes a project with size snapshots below root and returns
// the store plus the bytes the snapshot files take.
func historyFixture(b *testing.B, p historyProfile, size int) (Store, Config, string, int64) {
	b.Helper()
	root := b.TempDir()
	if err := os.Mkdir(filepath.Join(root, savepointDir), 0o755); err != nil {
		b.Fatal(err)
	}
	store := NewStore(root)
	cfg := historyConfig(p)
	if _, err := store.SaveConfig(cfg); err != nil {
		b.Fatal(err)
	}
	for n := 0; n < size; n++ {
		if _, err := store.SaveSnapshot(historySnapshot(cfg, p, n)); err != nil {
			b.Fatalf("snapshot %d: %v", n, err)
		}
	}
	entries, err := os.ReadDir(healthPath(root, snapshotsDir))
	if err != nil {
		b.Fatal(err)
	}
	var total int64
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			b.Fatal(err)
		}
		total += info.Size()
	}
	return store, cfg, root, total
}

// forEachHistory runs fn once per profile and size, building each fixture
// before any timing starts.
func forEachHistory(b *testing.B, fn func(b *testing.B, store Store, cfg Config, root string)) {
	for _, p := range historyProfiles {
		for _, size := range historySizes {
			store, cfg, root, bytes := historyFixture(b, p, size)
			b.Run(fmt.Sprintf("%s/n=%d", p.name, size), func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()
				fn(b, store, cfg, root)
				// After the run: ResetTimer would clear custom metrics.
				b.ReportMetric(float64(bytes), "fixture-bytes")
				b.ReportMetric(float64(bytes)/float64(size), "bytes/snapshot")
			})
		}
	}
}

// BenchmarkHealthHistoryLoadSnapshots is the filesystem part of a dashboard
// load: list, read, decode and validate every snapshot.
func BenchmarkHealthHistoryLoadSnapshots(b *testing.B) {
	forEachHistory(b, func(b *testing.B, store Store, _ Config, _ string) {
		for i := 0; i < b.N; i++ {
			if _, err := store.LoadSnapshots(); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkHealthHistoryLoadDashboard is everything the board's load command
// does: config, every snapshot, rows, sign-off and the capped history list.
func BenchmarkHealthHistoryLoadDashboard(b *testing.B) {
	forEachHistory(b, func(b *testing.B, _ Store, _ Config, root string) {
		for i := 0; i < b.N; i++ {
			if _, err := LoadDashboard(root); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkHealthHistoryProject is the in-memory projection from loaded
// snapshots to dashboard rows, with no file access.
func BenchmarkHealthHistoryProject(b *testing.B) {
	forEachHistory(b, func(b *testing.B, store Store, cfg Config, _ string) {
		b.StopTimer()
		snaps, err := store.LoadSnapshots()
		if err != nil {
			b.Fatal(err)
		}
		newest := snaps[len(snaps)-1]
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			dashboardRows(cfg, newest, dashboardHistory(snaps[:len(snaps)-1]), false)
		}
	})
}

// BenchmarkHealthHistoryRenderReport is the pure report text from a loaded
// dashboard.
func BenchmarkHealthHistoryRenderReport(b *testing.B) {
	forEachHistory(b, func(b *testing.B, _ Store, _ Config, root string) {
		b.StopTimer()
		d, err := LoadDashboard(root)
		if err != nil {
			b.Fatal(err)
		}
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			if _, err := RenderReport(d, "O-001"); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkHealthHistoryFreshnessScripted is the freshness code with Git
// answered from memory: only this package's own cost. It does not depend on
// how many snapshots exist, since freshness reads only the newest identity.
func BenchmarkHealthHistoryFreshnessScripted(b *testing.B) {
	root := b.TempDir()
	git := gitScript{head: strings.Repeat("a", 40), counts: "0\t3", ancestors: map[string]bool{}}
	recorded := RepositoryIdentity{Commit: strings.Repeat("b", 40), InputFingerprint: "sha256:" + strings.Repeat("1", 64)}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		DashboardFreshness(context.Background(), root, git, recorded)
	}
}

// BenchmarkHealthHistoryFreshnessGit runs the real git binary against a
// throwaway repository this benchmark creates, so the number is process
// launch plus Git's own work, not Savepoint's.
func BenchmarkHealthHistoryFreshnessGit(b *testing.B) {
	if _, err := exec.LookPath("git"); err != nil {
		b.Skip("git is not installed")
	}
	root := b.TempDir()
	run := func(args ...string) {
		cmd := exec.Command("git", append([]string{"-c", "user.name=bench", "-c", "user.email=bench@example.invalid", "-c", "commit.gpgsign=false"}, args...)...)
		cmd.Dir = root
		cmd.Env = append(os.Environ(), "GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_SYSTEM=/dev/null")
		if out, err := cmd.CombinedOutput(); err != nil {
			b.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	run("init", "-q")
	for i := 0; i < 20; i++ {
		if err := os.WriteFile(filepath.Join(root, fmt.Sprintf("f%d.go", i)), []byte("package f\n"), 0o644); err != nil {
			b.Fatal(err)
		}
	}
	run("add", ".")
	run("commit", "-q", "-m", "one")
	recorded, err := ObserveRepository(context.Background(), GitRunner{}, root, InputScope{})
	if err != nil {
		b.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "f0.go"), []byte("package f\n// moved\n"), 0o644); err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DashboardFreshness(context.Background(), root, GitRunner{}, recorded.Identity)
	}
}

// TestHistoryFixturesAreValidAndDeterministic keeps the benchmark builders
// honest without timing anything: the same inputs give the same identities,
// and every profile produces snapshots the real loader accepts.
func TestHistoryFixturesAreValidAndDeterministic(t *testing.T) {
	for _, p := range historyProfiles {
		cfg := historyConfig(p)
		a, b := historySnapshot(cfg, p, 3), historySnapshot(cfg, p, 3)
		if a.ID != b.ID {
			t.Errorf("%s: identity differs between builds", p.name)
		}
		if err := a.Validate(); err != nil {
			t.Errorf("%s: %v", p.name, err)
		}
		if got, want := len(a.Results), len(Capabilities())*p.instances; got != want {
			t.Errorf("%s: %d results, want %d", p.name, got, want)
		}
		if got := len(a.Results[0].Evidence); got != p.evidence {
			t.Errorf("%s: %d evidence entries, want %d", p.name, got, p.evidence)
		}
	}
}
