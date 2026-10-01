package codehealth

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// fakeTools is a ToolRunner whose behavior is chosen per executable name. It
// records the order of calls and fails if two ever overlap.
type fakeTools struct {
	t        *testing.T
	behavior map[string]func(ctx context.Context, spec ToolSpec) (ToolResult, error)
	calls    []ToolSpec
	inFlight atomic.Int32
}

func (f *fakeTools) Run(ctx context.Context, spec ToolSpec) (ToolResult, error) {
	if f.inFlight.Add(1) != 1 {
		f.t.Error("tools ran concurrently")
	}
	defer f.inFlight.Add(-1)
	f.calls = append(f.calls, spec)
	b, ok := f.behavior[spec.Executable]
	if !ok {
		return ToolResult{}, ErrToolUnavailable
	}
	return b(ctx, spec)
}

func stdout(s string) func(context.Context, ToolSpec) (ToolResult, error) {
	return func(context.Context, ToolSpec) (ToolResult, error) { return ToolResult{Stdout: []byte(s)}, nil }
}

type readerFunc func(ctx context.Context, in ReportInput) (Reading, error)

func (f readerFunc) Read(ctx context.Context, in ReportInput) (Reading, error) { return f(ctx, in) }

// fixedReader reports the report text as the number of bytes, in CCN.
func goodReading(n float64, unit Unit) Reading {
	return Reading{Value: &Value{Number: n, Unit: unit}, Provenance: Provenance{ProviderVersion: "1.0", MeasurementDefinition: "fake"}}
}

func okReader(n float64, unit Unit) Reader {
	return readerFunc(func(context.Context, ReportInput) (Reading, error) { return goodReading(n, unit), nil })
}

var testClock = func() time.Time { return time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC) }

func project(t *testing.T) string {
	t.Helper()
	dir := newRepo(t)
	if err := os.MkdirAll(filepath.Join(dir, ".savepoint"), 0o755); err != nil {
		t.Fatal(err)
	}
	return dir
}

func lizardInstance(name, exe string, scope ...string) CapabilityConfig {
	return CapabilityConfig{Capability: CapabilityComplexity, Provider: ProviderLizardCSV, Name: name, Executable: exe, Scope: scope}
}

func collect(t *testing.T, dir string, cfg Config, readers Readers, runner ToolRunner) Collection {
	t.Helper()
	got, err := Collect(context.Background(), CollectRequest{Root: dir, Origin: OriginManual, Config: cfg, Readers: readers, Runner: runner, Clock: testClock})
	if err != nil {
		t.Fatalf("Collect: %v", err)
	}
	return got
}

func cfgOf(items ...CapabilityConfig) Config {
	return Config{Version: ConfigVersion, Capabilities: items}
}

func collectedFor(t *testing.T, c Collection, capability Capability, name string) Collected {
	t.Helper()
	for _, r := range c.Results {
		if r.Result.Capability == capability && r.Result.Name == name {
			return r
		}
	}
	t.Fatalf("no result for %s %q", capability, name)
	return Collected{}
}

func TestCollectRunsInstancesSequentiallyInConfigOrderAndSaves(t *testing.T) {
	dir := project(t)
	tools := &fakeTools{t: t, behavior: map[string]func(context.Context, ToolSpec) (ToolResult, error){
		"lizard-web": stdout("w"), "lizard-api": stdout("a"),
	}}
	web, api := lizardInstance("web", "lizard-web", "web/**"), lizardInstance("api", "lizard-api", "api/**")
	web.Required = true
	got := collect(t, dir, cfgOf(web, api), Readers{ProviderLizardCSV: okReader(7, UnitCCN)}, tools)

	if len(tools.calls) != 2 || tools.calls[0].Executable != "lizard-web" || tools.calls[1].Executable != "lizard-api" {
		t.Fatalf("calls not in config order: %+v", tools.calls)
	}
	if tools.calls[0].Dir != dir {
		t.Fatalf("tool ran in %q, want project root %q", tools.calls[0].Dir, dir)
	}
	if got.Results[0].Result.Name != "web" || !got.Results[0].Required || got.Results[1].Required {
		t.Fatalf("required flags or order wrong: %+v", got.Results[:2])
	}
	for _, c := range []Capability{CapabilityTests, CapabilityCoverage, CapabilityDuplication, CapabilityDependencyVulnerability} {
		r := collectedFor(t, got, c, "").Result
		if r.Outcome != OutcomeNotConfigured || r.Freshness != FreshnessUnknown {
			t.Errorf("%s = %s/%s, want not_configured/unknown", c, r.Outcome, r.Freshness)
		}
	}
	snaps, err := NewStore(dir).LoadSnapshots()
	if err != nil || len(snaps) != 1 {
		t.Fatalf("snapshots = %d, err %v", len(snaps), err)
	}
	s := snaps[0]
	if s.ID != got.SnapshotID || !got.Created || s.Origin != OriginManual || s.Retention != RetentionPrunable {
		t.Fatalf("saved %+v, collection %+v", s, got)
	}
	r := collectedFor(t, got, CapabilityComplexity, "web").Result
	if r.Outcome != OutcomeAvailable || r.Freshness != FreshnessFresh || r.Value.Number != 7 || r.Provenance.Provider != ProviderLizardCSV ||
		r.ConfigDigest != web.Digest() || r.CollectedAt != "2026-10-01T12:00:00Z" || r.Scope[0] != "web/**" {
		t.Fatalf("web result = %+v", r)
	}
	if s.Summary.Overall == ClassificationGood {
		t.Fatal("unconfigured capabilities must keep the overall label from being good")
	}
}

func TestCollectOfficialOriginIsPermanent(t *testing.T) {
	dir := project(t)
	got, err := Collect(context.Background(), CollectRequest{Root: dir, Origin: OriginOfficial, Config: cfgOf(), Clock: testClock})
	if err != nil {
		t.Fatal(err)
	}
	snaps, _ := NewStore(dir).LoadSnapshots()
	if len(snaps) != 1 || snaps[0].ID != got.SnapshotID || snaps[0].Retention != RetentionPermanent {
		t.Fatalf("got %+v", snaps)
	}
}

func TestCollectExecutedOutcomesStayDistinct(t *testing.T) {
	dir := project(t)
	exit := func(code int, stderr string) func(context.Context, ToolSpec) (ToolResult, error) {
		return func(context.Context, ToolSpec) (ToolResult, error) {
			return ToolResult{ExitCode: code, Stderr: stderr}, nil
		}
	}
	tools := &fakeTools{t: t, behavior: map[string]func(context.Context, ToolSpec) (ToolResult, error){
		"fails":    exit(2, "boom"),
		"timeout":  func(context.Context, ToolSpec) (ToolResult, error) { return ToolResult{}, context.DeadlineExceeded },
		"findings": exit(1, ""),
		"huge": func(context.Context, ToolSpec) (ToolResult, error) {
			return ToolResult{Stdout: []byte("x"), Truncated: true}, nil
		},
		"good": stdout("ok"),
	}}
	cfg := cfgOf(
		lizardInstance("missing", "absent-tool", "a/**"),
		lizardInstance("fails", "fails", "b/**"),
		lizardInstance("timeout", "timeout", "c/**"),
		lizardInstance("huge", "huge", "d/**"),
		lizardInstance("good", "good", "e/**"),
		CapabilityConfig{Capability: CapabilityDependencyVulnerability, Provider: ProviderOSVScannerJSON, Executable: "findings"},
		CapabilityConfig{Capability: CapabilityDuplication, Provider: ProviderJscpdJSON, Executable: "nope"},
	)
	readers := Readers{
		ProviderLizardCSV:      okReader(3, UnitCCN),
		ProviderOSVScannerJSON: okReader(2, UnitCount),
		// jscpd deliberately has no reader.
	}
	got := collect(t, dir, cfg, readers, tools)

	want := map[string]Outcome{"missing": OutcomeUnavailable, "fails": OutcomeFailed, "timeout": OutcomeTimedOut, "huge": OutcomePartial, "good": OutcomeAvailable}
	for name, outcome := range want {
		r := collectedFor(t, got, CapabilityComplexity, name).Result
		if r.Outcome != outcome {
			t.Errorf("%s = %s (%s), want %s", name, r.Outcome, r.Reason, outcome)
		}
		if outcome != OutcomeAvailable && outcome != OutcomePartial && r.Value != nil {
			t.Errorf("%s carries a value", name)
		}
	}
	if r := collectedFor(t, got, CapabilityComplexity, "fails").Result; !strings.Contains(r.Reason, "status 2") || !strings.Contains(r.Reason, "boom") {
		t.Errorf("failure reason = %q", r.Reason)
	}
	if r := collectedFor(t, got, CapabilityComplexity, "huge").Result; r.Value != nil || !strings.Contains(r.Reason, "32 MiB") {
		t.Errorf("truncated result = %+v", r)
	}
	if r := collectedFor(t, got, CapabilityDependencyVulnerability, "").Result; r.Outcome != OutcomeAvailable || r.Value.Number != 2 {
		t.Errorf("exit 1 means findings for OSV-Scanner, got %s %q", r.Outcome, r.Reason)
	}
	if r := collectedFor(t, got, CapabilityDuplication, "").Result; r.Outcome != OutcomeUnsupported {
		t.Errorf("no reader = %s", r.Outcome)
	}
	for _, call := range tools.calls {
		if call.Executable == "nope" {
			t.Error("a tool whose report cannot be read was still run")
		}
	}
	// The good instance is untouched by its neighbours' failures.
	if r := collectedFor(t, got, CapabilityComplexity, "good").Result; r.Outcome != OutcomeAvailable || r.Value.Number != 3 || r.Freshness != FreshnessFresh {
		t.Errorf("good = %+v", r)
	}
}

func TestCollectReaderProblemsAreContained(t *testing.T) {
	dir := project(t)
	tools := &fakeTools{t: t, behavior: map[string]func(context.Context, ToolSpec) (ToolResult, error){"l": stdout("data")}}
	byName := map[string]func(ReportInput) (Reading, error){
		"malformed": func(ReportInput) (Reading, error) { return Reading{}, errors.New("bad csv\nline 3") },
		"partial": func(ReportInput) (Reading, error) {
			r := goodReading(9, UnitCCN)
			r.Partial, r.Reason = true, "two files skipped"
			return r, nil
		},
		"nonsense": func(ReportInput) (Reading, error) { return goodReading(9, UnitPercent), nil }, // wrong unit
		"ok":       func(in ReportInput) (Reading, error) { return goodReading(float64(len(in.Data)), UnitCCN), nil },
	}
	reader := readerFunc(func(_ context.Context, in ReportInput) (Reading, error) {
		return byName[in.Scope[0]](in)
	})
	var items []CapabilityConfig
	for name := range byName {
		items = append(items, lizardInstance(name, "l", name))
	}
	got := collect(t, dir, cfgOf(items...), Readers{ProviderLizardCSV: reader}, tools)

	if r := collectedFor(t, got, CapabilityComplexity, "malformed").Result; r.Outcome != OutcomeFailed || !strings.Contains(r.Reason, "lizard-csv") || strings.Contains(r.Reason, "\n") {
		t.Errorf("malformed = %s %q", r.Outcome, r.Reason)
	}
	if r := collectedFor(t, got, CapabilityComplexity, "partial").Result; r.Outcome != OutcomePartial || r.Value == nil || r.Reason != "two files skipped" {
		t.Errorf("partial = %+v", r)
	}
	if r := collectedFor(t, got, CapabilityComplexity, "nonsense").Result; r.Outcome != OutcomeFailed || r.Value != nil || !strings.Contains(r.Reason, "unusable") {
		t.Errorf("nonsense = %+v", r)
	}
	if r := collectedFor(t, got, CapabilityComplexity, "ok").Result; r.Outcome != OutcomeAvailable || r.Value.Number != 4 {
		t.Errorf("ok = %+v", r)
	}
	if _, err := NewStore(dir).LoadSnapshots(); err != nil {
		t.Fatalf("snapshot with a bad reader did not save cleanly: %v", err)
	}
}

func TestCollectCancellationStopsRemainingInstances(t *testing.T) {
	dir := project(t)
	ctx, cancel := context.WithCancel(context.Background())
	tools := &fakeTools{t: t, behavior: map[string]func(context.Context, ToolSpec) (ToolResult, error){
		"first":  stdout("x"),
		"second": func(ctx context.Context, _ ToolSpec) (ToolResult, error) { cancel(); return ToolResult{}, ctx.Err() },
		"third":  stdout("x"),
	}}
	cfg := cfgOf(lizardInstance("a", "first", "a/**"), lizardInstance("b", "second", "b/**"), lizardInstance("c", "third", "c/**"))
	got, err := Collect(ctx, CollectRequest{Root: dir, Origin: OriginManual, Config: cfg, Readers: Readers{ProviderLizardCSV: okReader(1, UnitCCN)}, Runner: tools, Clock: testClock})
	if err != nil {
		t.Fatal(err)
	}
	if len(tools.calls) != 2 {
		t.Fatalf("third instance ran after cancellation: %+v", tools.calls)
	}
	for name, want := range map[string]Outcome{"a": OutcomeAvailable, "b": OutcomeCancelled, "c": OutcomeCancelled} {
		if r := collectedFor(t, got, CapabilityComplexity, name).Result; r.Outcome != want {
			t.Errorf("%s = %s, want %s", name, r.Outcome, want)
		}
	}
	if snaps, _ := NewStore(dir).LoadSnapshots(); len(snaps) != 1 {
		t.Fatalf("a cancelled run still saves its snapshot; got %d", len(snaps))
	}
}

func TestCollectReportOnlyNeverRunsATool(t *testing.T) {
	dir := project(t)
	tools := &fakeTools{t: t}
	cover := func(name, report string, scope ...string) CapabilityConfig {
		return CapabilityConfig{Capability: CapabilityCoverage, Provider: ProviderGoCoverProfile, Name: name, Report: report, Scope: scope}
	}
	readers := Readers{ProviderGoCoverProfile: okReader(85, UnitPercent)}

	old := time.Now().Add(-time.Hour)
	write(t, dir, "reports/stale.out", "mode: set\n")
	write(t, dir, "reports/fresh.out", "mode: set\n")
	if err := os.Chtimes(filepath.Join(dir, "reports/stale.out"), old, old); err != nil {
		t.Fatal(err)
	}
	write(t, dir, "reports/dir.out/x", "x")
	outside := t.TempDir()
	write(t, outside, "secret.out", "mode: set\n")
	if err := os.Symlink(filepath.Join(outside, "secret.out"), filepath.Join(dir, "reports/escape.out")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	cfg := cfgOf(
		cover("missing", "reports/none.out", "a/**"),
		cover("stale", "reports/stale.out", "."), // replaced below
		cover("fresh", "reports/fresh.out"),
		cover("dir", "reports/dir.out", "d/**"),
		cover("escape", "reports/escape.out", "e/**"),
	)
	// Instances sharing a provider need distinct scopes; make stale and fresh
	// both whole-repository by giving them one pattern each.
	cfg.Capabilities[1].Scope = []string{"*.go", "sub/**"}
	cfg.Capabilities[2].Scope = []string{"*.go"}
	got := collect(t, dir, cfg, readers, tools)

	if len(tools.calls) != 0 {
		t.Fatalf("report-only providers ran a tool: %+v", tools.calls)
	}
	want := map[string][2]any{
		"missing": {OutcomeAbsent, FreshnessUnknown},
		"stale":   {OutcomeAvailable, FreshnessStale},
		"fresh":   {OutcomeAvailable, FreshnessFresh},
		"dir":     {OutcomeFailed, FreshnessUnknown},
		"escape":  {OutcomeFailed, FreshnessUnknown},
	}
	for name, w := range want {
		r := collectedFor(t, got, CapabilityCoverage, name).Result
		if r.Outcome != w[0] || r.Freshness != w[1] {
			t.Errorf("%s = %s/%s (%s), want %v/%v", name, r.Outcome, r.Freshness, r.Reason, w[0], w[1])
		}
	}
}

func TestCollectReportFreshnessIgnoresTheReportItself(t *testing.T) {
	dir := project(t)
	// The report is newer than every other input only because it was just
	// written; it must not make itself stale or count as an input.
	write(t, dir, "coverage.out", "mode: set\n")
	future := time.Now().Add(time.Hour)
	os.Chtimes(filepath.Join(dir, "coverage.out"), future, future)
	cfg := cfgOf(CapabilityConfig{Capability: CapabilityCoverage, Provider: ProviderGoCoverProfile, Report: "coverage.out"})
	got := collect(t, dir, cfg, Readers{ProviderGoCoverProfile: okReader(90, UnitPercent)}, &fakeTools{t: t})
	if r := collectedFor(t, got, CapabilityCoverage, "").Result; r.Freshness != FreshnessFresh {
		t.Fatalf("freshness = %s", r.Freshness)
	}
	// An input newer than the report makes it stale.
	past := time.Now().Add(-2 * time.Hour)
	os.Chtimes(filepath.Join(dir, "coverage.out"), past, past)
	got = collect(t, dir, cfg, Readers{ProviderGoCoverProfile: okReader(90, UnitPercent)}, &fakeTools{t: t})
	if r := collectedFor(t, got, CapabilityCoverage, "").Result; r.Freshness != FreshnessStale {
		t.Fatalf("freshness = %s, want stale", r.Freshness)
	}
}

func TestCollectOversizedReportFileIsPartial(t *testing.T) {
	dir := project(t)
	f, err := os.Create(filepath.Join(dir, "big.out"))
	if err != nil {
		t.Fatal(err)
	}
	f.Truncate(MaxReportBytes + 10)
	f.Close()
	cfg := cfgOf(CapabilityConfig{Capability: CapabilityCoverage, Provider: ProviderGoCoverProfile, Report: "big.out"})
	called := false
	reader := readerFunc(func(context.Context, ReportInput) (Reading, error) { called = true; return Reading{}, nil })
	got := collect(t, dir, cfg, Readers{ProviderGoCoverProfile: reader}, &fakeTools{t: t})
	r := collectedFor(t, got, CapabilityCoverage, "").Result
	if r.Outcome != OutcomePartial || r.Value != nil || called || !strings.Contains(r.Reason, "32 MiB") {
		t.Fatalf("got %+v, reader called=%v", r, called)
	}
}

func TestCollectWithRealProcessesAndReportFiles(t *testing.T) {
	helperEnv(t)
	tmp := t.TempDir()
	t.Setenv("TMPDIR", tmp)
	t.Setenv("TMP", tmp)
	t.Setenv("TEMP", tmp)
	dir := project(t)
	exe, _ := os.Executable()
	base := []string{"-test.run=^TestHelperProcess$", "--"}
	tool := func(name string, timeout int, args ...string) CapabilityConfig {
		cc := lizardInstance(name, exe, name+"/**")
		cc.Args, cc.TimeoutSeconds = append(append([]string{}, base...), args...), timeout
		return cc
	}
	var seen []string
	reader := readerFunc(func(_ context.Context, in ReportInput) (Reading, error) {
		seen = append(seen, string(in.Data))
		return goodReading(1, UnitCCN), nil
	})
	cfg := cfgOf(
		tool("stdout", 0, "ok"),
		tool("file", 0, "report", "{report}"),
		tool("silent", 0, "silent", "{report}"),
		tool("hang", 1, "hang"),
		tool("bad", 0, "exit", "9"),
	)
	got := collect(t, dir, cfg, Readers{ProviderLizardCSV: reader}, ExecRunner{})

	if strings.Join(seen, ",") != "report-bytes,file-report" {
		t.Errorf("reader saw %q", seen)
	}
	want := map[string]Outcome{"stdout": OutcomeAvailable, "file": OutcomeAvailable, "silent": OutcomeFailed, "hang": OutcomeTimedOut, "bad": OutcomeFailed}
	for name, outcome := range want {
		if r := collectedFor(t, got, CapabilityComplexity, name).Result; r.Outcome != outcome {
			t.Errorf("%s = %s (%s), want %s", name, r.Outcome, r.Reason, outcome)
		}
	}
	if entries, _ := os.ReadDir(tmp); len(entries) != 0 {
		t.Errorf("temporary files left behind: %v", entries)
	}
}

func TestCollectCancelledRealProcessLeavesNoTemporaryFiles(t *testing.T) {
	helperEnv(t)
	tmp := t.TempDir()
	t.Setenv("TMPDIR", tmp)
	t.Setenv("TMP", tmp)
	t.Setenv("TEMP", tmp)
	dir := project(t)
	exe, _ := os.Executable()
	cc := lizardInstance("", exe)
	cc.Args = []string{"-test.run=^TestHelperProcess$", "--", "hang", "{report}"}
	ctx, cancel := context.WithCancel(context.Background())
	go func() { time.Sleep(300 * time.Millisecond); cancel() }()
	start := time.Now()
	got, err := Collect(ctx, CollectRequest{Root: dir, Origin: OriginManual, Config: cfgOf(cc), Readers: Readers{ProviderLizardCSV: okReader(1, UnitCCN)}, Clock: testClock})
	if err != nil {
		t.Fatal(err)
	}
	if r := collectedFor(t, got, CapabilityComplexity, "").Result; r.Outcome != OutcomeCancelled || time.Since(start) > 10*time.Second {
		t.Fatalf("got %s after %s", r.Outcome, time.Since(start))
	}
	if entries, _ := os.ReadDir(tmp); len(entries) != 0 {
		t.Errorf("temporary files left behind: %v", entries)
	}
}

func TestCollectRefusesBeforeRunningAnything(t *testing.T) {
	dir := project(t)
	tools := &fakeTools{t: t, behavior: map[string]func(context.Context, ToolSpec) (ToolResult, error){"l": stdout("x")}}
	req := CollectRequest{Root: dir, Origin: OriginManual, Config: cfgOf(lizardInstance("", "l")), Readers: Readers{ProviderLizardCSV: okReader(1, UnitCCN)}, Runner: tools, Clock: testClock}

	t.Run("bad origin", func(t *testing.T) {
		r := req
		r.Origin = "bogus"
		if _, err := Collect(context.Background(), r); !errors.Is(err, ErrInvalidOrigin) {
			t.Fatalf("got %v", err)
		}
	})
	t.Run("invalid configuration", func(t *testing.T) {
		r := req
		r.Config = Config{Version: 99}
		if _, err := Collect(context.Background(), r); !errors.Is(err, ErrUnsupportedVersion) {
			t.Fatalf("got %v", err)
		}
	})
	t.Run("damaged history", func(t *testing.T) {
		write(t, dir, ".savepoint/health/snapshots/junk.json", "{")
		if _, err := Collect(context.Background(), req); err == nil {
			t.Fatal("collected over a damaged history")
		}
		os.Remove(filepath.Join(dir, ".savepoint/health/snapshots/junk.json"))
	})
	t.Run("not a repository", func(t *testing.T) {
		r := req
		r.Root = t.TempDir()
		if _, err := Collect(context.Background(), r); err == nil {
			t.Fatal("collected outside a repository")
		}
	})
	if len(tools.calls) != 0 {
		t.Fatalf("a tool ran despite the refusal: %+v", tools.calls)
	}
}

func TestCollectNeverPrunesAndAccumulatesHistory(t *testing.T) {
	dir := project(t)
	tools := &fakeTools{t: t, behavior: map[string]func(context.Context, ToolSpec) (ToolResult, error){"l": stdout("x")}}
	cfg := cfgOf(lizardInstance("", "l"))
	n := ManualRetention + 2
	for i := 0; i < n; i++ {
		clock := func() time.Time { return time.Date(2026, 10, 1, 12, 0, i, 0, time.UTC) }
		if _, err := Collect(context.Background(), CollectRequest{Root: dir, Origin: OriginManual, Config: cfg, Readers: Readers{ProviderLizardCSV: okReader(1, UnitCCN)}, Runner: tools, Clock: clock}); err != nil {
			t.Fatal(err)
		}
	}
	snaps, err := NewStore(dir).LoadSnapshots()
	if err != nil || len(snaps) != n {
		t.Fatalf("snapshots = %d, want %d (err %v)", len(snaps), n, err)
	}
}

func TestCollectDefaultsToRealRunnersAndClock(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("covered by the other real-process tests")
	}
	dir := project(t)
	got, err := Collect(context.Background(), CollectRequest{Root: dir, Origin: OriginManual, Config: cfgOf(lizardInstance("", filepath.Join(dir, "missing-tool"))), Readers: Readers{ProviderLizardCSV: okReader(1, UnitCCN)}})
	if err != nil {
		t.Fatal(err)
	}
	if r := collectedFor(t, got, CapabilityComplexity, "").Result; r.Outcome != OutcomeUnavailable {
		t.Fatalf("got %s", r.Outcome)
	}
}
