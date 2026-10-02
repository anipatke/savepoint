package codehealth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

// Cross-stack reader matrix (T-093). Every scenario is offline: report-only
// providers read files written into a temporary project and executed
// providers are answered by a fake runner.
//
//	Reader (provider)       Go           JavaScript        Python
//	go-test-json            measured     -                 -
//	vitest-junit            -            malformed=failed  -
//	pytest-junit            -            -                 measured
//	go-cover-profile        measured     -                 -
//	vitest-v8               -            absent (required) -
//	coveragepy-json         -            -                 malformed=failed
//	lizard-csv              measured     unavailable       measured
//	jscpd-json              -            report total only report total only
//	osv-scanner-json        scanner error reported, no value
//
// Each reader's own good, partial, malformed and boundary cases live in the
// reader_*_test.go files; this file proves they stay independent when one
// Collect runs them together, and that the verdict, dashboard and history
// read the same snapshot the same way.

const (
	stackGoTests = "reports/go-test.jsonl"
	stackJSTests = "reports/vitest.xml"
	stackPyTests = "reports/pytest.xml"
	stackGoCover = "reports/go.out"
	stackJSCover = "reports/vitest-cov.json" // never written: the report is absent
	stackPyCover = "reports/coverage.json"
)

// stackProject is a project whose report files are the fixtures of a Go,
// JavaScript and Python stack. The JavaScript tests and Python coverage
// reports are malformed on purpose, and the JavaScript coverage report is
// missing.
func stackProject(t *testing.T) string {
	t.Helper()
	root := project(t)
	write(t, root, "go.mod", "module example.com/m\n")
	write(t, root, "a/x.go", "package a\n")
	put := func(rel, dir, fixture string) {
		write(t, root, rel, strings.ReplaceAll(fixtureBytes(t, dir, fixture), "/work/proj", filepath.ToSlash(root)))
	}
	put(stackGoTests, "tests", "go-fail.jsonl")
	put(stackJSTests, "tests", "junit-malformed.xml")
	put(stackPyTests, "tests", "pytest-fail.xml")
	put(stackGoCover, "coverage", "go-populated.out")
	put(stackPyCover, "coverage", "coveragepy-malformed.json")
	// A report is stale when an input it covers is newer, so date the reports
	// ahead of every input except the Go coverage profile, which predates the
	// Go sources it covers.
	later, earlier := time.Now().Add(time.Hour), time.Now().Add(-time.Hour)
	for _, rel := range []string{stackGoTests, stackJSTests, stackPyTests, stackPyCover} {
		if err := os.Chtimes(filepath.Join(root, rel), later, later); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Chtimes(filepath.Join(root, stackGoCover), earlier, earlier); err != nil {
		t.Fatal(err)
	}
	return root
}

// stackConfig has one instance per reader. jsCoverageRequired makes the
// missing JavaScript coverage report a required instance.
func stackConfig(jsCoverageRequired bool) Config {
	jsCover := CapabilityConfig{Capability: CapabilityCoverage, Provider: ProviderVitestV8, Name: "js", Report: stackJSCover, Required: jsCoverageRequired}
	return cfgOf(
		CapabilityConfig{Capability: CapabilityTests, Provider: ProviderGoTestJSON, Name: "go", Report: stackGoTests},
		CapabilityConfig{Capability: CapabilityTests, Provider: ProviderVitestJUnit, Name: "js", Scope: []string{"web/**"}, Report: stackJSTests},
		CapabilityConfig{Capability: CapabilityTests, Provider: ProviderPytestJUnit, Name: "py", Scope: []string{"tests/**"}, Report: stackPyTests},
		CapabilityConfig{Capability: CapabilityCoverage, Provider: ProviderGoCoverProfile, Name: "go", Scope: []string{"a/**"}, Report: stackGoCover},
		jsCover,
		CapabilityConfig{Capability: CapabilityCoverage, Provider: ProviderCoveragePyJSON, Name: "py", Report: stackPyCover},
		CapabilityConfig{Capability: CapabilityComplexity, Provider: ProviderLizardCSV, Name: "go", Executable: "lizard-go", Scope: []string{"cmd/**", "pkg/**"}},
		CapabilityConfig{Capability: CapabilityComplexity, Provider: ProviderLizardCSV, Name: "js", Executable: "lizard-js", Scope: []string{"web/**"}},
		CapabilityConfig{Capability: CapabilityComplexity, Provider: ProviderLizardCSV, Name: "py", Executable: "lizard-py", Scope: []string{"app/**"}},
		CapabilityConfig{Capability: CapabilityDuplication, Provider: ProviderJscpdJSON, Name: "js", Executable: "jscpd", Scope: []string{"web/**"}},
		CapabilityConfig{Capability: CapabilityDuplication, Provider: ProviderJscpdJSON, Name: "py", Executable: "jscpd", Scope: []string{"app/**"}},
		CapabilityConfig{Capability: CapabilityDependencyVulnerability, Provider: ProviderOSVScannerJSON, Executable: "osv-scanner"},
	)
}

// stackTools answers the executed providers. lizard-js is not installed, and
// osv-scanner exits 1 with a report whose only content is a scanner error.
func stackTools(t *testing.T, root string) *fakeTools {
	t.Helper()
	csv := fixtureBytes(t, "complexity", "mixed.csv")
	return &fakeTools{t: t, behavior: map[string]func(context.Context, ToolSpec) (ToolResult, error){
		// The Go run sees only files inside the project; the Python run also
		// lists one outside it, which the reader cannot place.
		"lizard-go": stdout(strings.ReplaceAll(csv, "/work/proj", filepath.ToSlash(root))),
		"lizard-py": stdout(csv),
		// jscpd 5 writes no per-file counts, so a scoped run has only the
		// report's own totals to offer.
		"jscpd": stdout(fixtureBytes(t, "duplication", "jscpd5.json")),
		"osv-scanner": func(context.Context, ToolSpec) (ToolResult, error) {
			return ToolResult{Stdout: []byte(fixtureBytes(t, "vulnerabilities", "scanner-error.json")), ExitCode: 1}, nil
		},
	}}
}

// claimsHealth matches the word "healthy" but not "unhealthy".
var claimsHealth = regexp.MustCompile(`(?i)\bhealthy\b`)

func collectAt(t *testing.T, root string, origin Origin, cfg Config, tools ToolRunner, at time.Time) Collection {
	t.Helper()
	got, err := Collect(context.Background(), CollectRequest{
		Root: root, Origin: origin, Config: cfg, Readers: DefaultReaders(), Runner: tools,
		Clock: func() time.Time { return at },
	})
	if err != nil {
		t.Fatalf("Collect: %v", err)
	}
	return got
}

func stackResult(t *testing.T, c Collection, p ProviderKey, name string) CapabilityResult {
	t.Helper()
	for _, r := range c.Results {
		if r.Result.Provenance.Provider == p && r.Result.Name == name {
			return r.Result
		}
	}
	t.Fatalf("no result for %s %q", p, name)
	return CapabilityResult{}
}

func verdictFor(t *testing.T, v Verdict, p ProviderKey, name string) ResultVerdict {
	t.Helper()
	for _, r := range v.Results {
		if r.Provider == p && r.Name == name {
			return r
		}
	}
	t.Fatalf("no verdict for %s %q", p, name)
	return ResultVerdict{}
}

func TestPolyglotReadersKeepEachStackIndependent(t *testing.T) {
	cfg := stackConfig(true)
	root := stackProject(t)
	got := collectAt(t, root, OriginOfficial, cfg, stackTools(t, root), testClock())

	assertStacksReadAloneAsBesideOthers(t, cfg, got)
	assertStackOutcomes(t, got)
	assertScopedJscpdKeepsReportTotals(t, got)

	official := loadOfficialSnapshot(t, root)
	assertSnapshotClassifications(t, official)
	assertStackVerdicts(t, official, cfg)
	assertOnlyOfficialSnapshotsHaveStanding(t, root, cfg)
}

// stackWantOutcome is the outcome each configured instance must end with.
var stackWantOutcome = map[string]Outcome{
	string(ProviderGoTestJSON) + "/go":     OutcomeAvailable,
	string(ProviderVitestJUnit) + "/js":    OutcomeFailed,
	string(ProviderPytestJUnit) + "/py":    OutcomeAvailable,
	string(ProviderGoCoverProfile) + "/go": OutcomeAvailable,
	string(ProviderVitestV8) + "/js":       OutcomeAbsent,
	string(ProviderCoveragePyJSON) + "/py": OutcomeFailed,
	string(ProviderLizardCSV) + "/go":      OutcomeAvailable,
	string(ProviderLizardCSV) + "/js":      OutcomeUnavailable,
	string(ProviderLizardCSV) + "/py":      OutcomePartial,
	string(ProviderJscpdJSON) + "/js":      OutcomeAvailable,
	string(ProviderJscpdJSON) + "/py":      OutcomeAvailable,
	string(ProviderOSVScannerJSON) + "/":   OutcomePartial,
}

// assertStacksReadAloneAsBesideOthers checks that every instance, run alone in
// a fresh project, reads exactly what it read beside the others: a failed tool
// changes nobody else's result.
func assertStacksReadAloneAsBesideOthers(t *testing.T, cfg Config, got Collection) {
	t.Helper()
	for _, cc := range cfg.Capabilities {
		soloRoot := stackProject(t)
		solo := collectAt(t, soloRoot, OriginManual, cfgOf(cc), stackTools(t, soloRoot), testClock())
		want, have := stackResult(t, solo, cc.Provider, cc.Name), stackResult(t, got, cc.Provider, cc.Name)
		if have.Outcome != want.Outcome || fmt.Sprint(have.Value) != fmt.Sprint(want.Value) || fmt.Sprint(have.Details) != fmt.Sprint(want.Details) {
			t.Errorf("%s %q beside the others = %s %v %v, alone = %s %v %v",
				cc.Provider, cc.Name, have.Outcome, have.Value, have.Details, want.Outcome, want.Value, want.Details)
		}
	}
}

func assertStackOutcomes(t *testing.T, got Collection) {
	t.Helper()
	for _, col := range got.Results {
		r := col.Result
		if r.Outcome == OutcomeNotConfigured {
			continue
		}
		key := string(r.Provenance.Provider) + "/" + r.Name
		want, ok := stackWantOutcome[key]
		if !ok {
			t.Errorf("unexpected result %s", key)
			continue
		}
		if r.Outcome != want {
			t.Errorf("%s = %s (%s), want %s", key, r.Outcome, r.Reason, want)
		}
		assertValueMatchesOutcome(t, key, r)
	}
}

// assertValueMatchesOutcome checks that whatever did not measure carries no
// number, so it can never read as a healthy zero.
func assertValueMatchesOutcome(t *testing.T, key string, r CapabilityResult) {
	t.Helper()
	if !r.Outcome.Measured() && (r.Value != nil || len(r.Details) > 0 || len(r.Evidence) > 0) {
		t.Errorf("%s is %s but carries value %v details %v evidence %v", key, r.Outcome, r.Value, r.Details, r.Evidence)
	}
	if r.Outcome.Measured() && r.Value == nil {
		t.Errorf("%s is %s without a value", key, r.Outcome)
	}
}

// assertScopedJscpdKeepsReportTotals checks duplication: jscpd 5 has no
// per-file counts, so each scoped instance keeps the report's own whole-run
// totals and invents no scoped figure.
func assertScopedJscpdKeepsReportTotals(t *testing.T, got Collection) {
	t.Helper()
	for _, name := range []string{"js", "py"} {
		r := stackResult(t, got, ProviderJscpdJSON, name)
		lines := 0.0
		for _, d := range r.Details {
			if d.Key == DetailTotalLines {
				lines = d.Number
			}
		}
		if r.Value == nil || r.Value.Number != 3.33 || lines != 600 {
			t.Errorf("jscpd %s = %v with %v total lines, want the report's 3.33%% of 600", name, r.Value, lines)
		}
	}
}

func loadOfficialSnapshot(t *testing.T, root string) Snapshot {
	t.Helper()
	snaps, err := NewStore(root).LoadSnapshots()
	if err != nil || len(snaps) != 1 {
		t.Fatalf("LoadSnapshots = %d snapshots, %v", len(snaps), err)
	}
	if err := snaps[0].Validate(); err != nil {
		t.Fatalf("snapshot is invalid: %v", err)
	}
	return snaps[0]
}

func assertSnapshotClassifications(t *testing.T, official Snapshot) {
	t.Helper()
	for _, s := range official.Summary.Capabilities {
		switch stackWantOutcome[string(s.Provider)+"/"+s.Name] {
		case OutcomeAvailable:
		case OutcomePartial:
			if s.Classification == ClassificationGood {
				t.Errorf("%s %q is partial but was classified good", s.Provider, s.Name)
			}
		default:
			if s.Classification != ClassificationUnknown {
				t.Errorf("%s %q did not measure but was classified %s", s.Provider, s.Name, s.Classification)
			}
		}
	}
}

// assertStackVerdicts checks that the verdict keeps the stacks apart too, and
// that required versus optional only changes the absent coverage report.
func assertStackVerdicts(t *testing.T, official Snapshot, cfg Config) {
	t.Helper()
	verdict, err := Evaluate(official, cfg)
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	wantVerdict := []struct {
		provider ProviderKey
		name     string
		blocks   bool
		kind     Kind
	}{
		{ProviderGoTestJSON, "go", true, KindUnhealthy},
		{ProviderPytestJUnit, "py", true, KindUnhealthy},
		{ProviderVitestJUnit, "js", false, KindCollectionFailure},
		{ProviderVitestV8, "js", true, KindCollectionFailure}, // required and absent
		{ProviderCoveragePyJSON, "py", false, KindCollectionFailure},
		{ProviderLizardCSV, "js", false, KindCollectionFailure},
		{ProviderOSVScannerJSON, "", false, KindIncomplete},
		{ProviderLizardCSV, "py", false, KindIncomplete},
		{ProviderLizardCSV, "go", false, KindNoFinding},
		{ProviderGoCoverProfile, "go", false, KindStale},
		{ProviderJscpdJSON, "js", false, KindNoFinding},
	}
	for _, w := range wantVerdict {
		rv := verdictFor(t, verdict, w.provider, w.name)
		if (rv.Disposition == DispositionBlocks) != w.blocks || rv.Kind != w.kind {
			t.Errorf("%s %q = %s/%s (%s), want blocks=%v kind %s", w.provider, w.name, rv.Disposition, rv.Kind, rv.Reason, w.blocks, w.kind)
		}
	}
	if claimsHealth.MatchString(verdict.Render()) {
		t.Errorf("verdict claims health:\n%s", verdict.Render())
	}

	// The passing Go and Python tests' failures block either way, so Blocks()
	// alone cannot show the difference.
	optional, err := Evaluate(official, stackConfig(false))
	if err != nil {
		t.Fatalf("Evaluate optional: %v", err)
	}
	if rv := verdictFor(t, optional, ProviderVitestV8, "js"); rv.Disposition != DispositionReported || rv.Kind != KindCollectionFailure {
		t.Errorf("optional absent coverage = %s/%s, want reported collection failure", rv.Disposition, rv.Kind)
	}
}

// assertOnlyOfficialSnapshotsHaveStanding checks that a manual snapshot has no
// standing as Check evidence.
func assertOnlyOfficialSnapshotsHaveStanding(t *testing.T, root string, cfg Config) {
	t.Helper()
	manual := collectAt(t, root, OriginManual, cfg, stackTools(t, root), testClock().Add(time.Minute))
	all, err := NewStore(root).LoadSnapshots()
	if err != nil || len(all) != 2 {
		t.Fatalf("LoadSnapshots = %d, %v", len(all), err)
	}
	for _, s := range all {
		if s.ID != manual.SnapshotID {
			continue
		}
		if _, err := Evaluate(s, cfg); !errors.Is(err, ErrManualSnapshot) {
			t.Errorf("Evaluate(manual) = %v, want ErrManualSnapshot", err)
		}
	}
}

func TestPolyglotDashboardAgreesWithTheVerdict(t *testing.T) {
	cfg := stackConfig(true)
	root := stackProject(t)
	if _, err := NewStore(root).SaveConfig(cfg); err != nil {
		t.Fatal(err)
	}
	got := collectAt(t, root, OriginOfficial, cfg, stackTools(t, root), testClock())
	snaps, err := NewStore(root).LoadSnapshots()
	if err != nil || len(snaps) != 1 || snaps[0].ID != got.SnapshotID {
		t.Fatalf("LoadSnapshots = %d snapshots, %v", len(snaps), err)
	}
	verdict, err := Evaluate(snaps[0], cfg)
	if err != nil {
		t.Fatal(err)
	}
	d := mustLoad(t, root)

	if !verdict.Blocks() || d.SignOff != textSignOffBlocks {
		t.Errorf("verdict blocks = %v, dashboard sign-off %q", verdict.Blocks(), d.SignOff)
	}
	if claimsHealth.MatchString(d.Headline + d.SignOff + d.OverallText) {
		t.Errorf("dashboard claims health: %q %q %q", d.Headline, d.SignOff, d.OverallText)
	}
	for _, rv := range verdict.Results {
		if rv.Disposition == DispositionNotConfigured {
			continue
		}
		var row DashboardRow
		for _, r := range d.Rows {
			if r.Capability == rv.Capability && r.Provider == rv.Provider && r.Name == rv.Name {
				row = r
			}
		}
		if row.Capability == "" {
			t.Errorf("no dashboard row for %s %s %q", rv.Capability, rv.Provider, rv.Name)
			continue
		}
		if row.BlocksSignOff() != (rv.Disposition == DispositionBlocks) {
			t.Errorf("%s %q: dashboard blocks=%v, verdict %s", rv.Provider, rv.Name, row.BlocksSignOff(), rv.Disposition)
		}
		if !row.Outcome.Measured() && (row.Label != ClassificationUnknown || row.ValueLabel != ClassificationUnknown) {
			t.Errorf("%s %q is %s but shows %s/%s", rv.Provider, rv.Name, row.Outcome, row.Label, row.ValueLabel)
		}
	}
}

// historyTools lets a scenario change the jscpd report between collections.
func historyTools(t *testing.T, root string, jscpd *string) *fakeTools {
	t.Helper()
	tools := stackTools(t, root)
	tools.behavior["jscpd"] = func(context.Context, ToolSpec) (ToolResult, error) {
		return ToolResult{Stdout: []byte(*jscpd)}, nil
	}
	return tools
}

func TestPolyglotHistoryComparabilityAcrossChanges(t *testing.T) {
	root := stackProject(t)
	report := fixtureBytes(t, "duplication", "jscpd5.json")
	tools := historyTools(t, root, &report)
	cfgAt := func(api string, webScope ...string) Config {
		return cfgOf(
			CapabilityConfig{Capability: CapabilityComplexity, Provider: ProviderLizardCSV, Name: api, Executable: "lizard-go", Scope: []string{"app/**"}},
			CapabilityConfig{Capability: CapabilityComplexity, Provider: ProviderLizardCSV, Name: "web", Executable: "lizard-go", Scope: webScope},
			CapabilityConfig{Capability: CapabilityDuplication, Provider: ProviderJscpdJSON, Name: "dup", Executable: "jscpd"},
		)
	}
	cfg := cfgAt("api", "web/**")
	at := testClock()
	next := func() time.Time { at = at.Add(time.Minute); return at }
	rowOf := func(d Dashboard, c Capability, name string) DashboardRow {
		for _, r := range d.Rows {
			if r.Capability == c && r.Name == name {
				return r
			}
		}
		t.Fatalf("no row for %s %q", c, name)
		return DashboardRow{}
	}
	save := func(c Config) {
		if _, err := NewStore(root).SaveConfig(c); err != nil {
			t.Fatal(err)
		}
	}

	for i := 0; i < 3; i++ {
		save(cfg)
		collectAt(t, root, OriginOfficial, cfg, tools, next())
	}
	d := mustLoad(t, root)
	if got := rowOf(d, CapabilityComplexity, "api"); got.Spark == "" || got.SparkNote != textSparkEarly {
		t.Fatalf("three official checks: spark %q note %q, want a drawing with the early note", got.Spark, got.SparkNote)
	}

	// Renaming an instance keeps its series: the new name inherits the trend.
	renamed := cfgAt("backend", "web/**")
	save(renamed)
	collectAt(t, root, OriginOfficial, renamed, tools, next())
	d = mustLoad(t, root)
	if got := rowOf(d, CapabilityComplexity, "backend"); got.Spark == "" || strings.Contains(got.SparkNote, textSparkRestarted) || strings.Contains(got.Basis, "not compared") {
		t.Errorf("renamed instance lost its history: spark %q note %q basis %q", got.Spark, got.SparkNote, got.Basis)
	}

	// Changing the scope starts a new series for that instance only.
	rescoped := cfgAt("backend", "web/src/**")
	save(rescoped)
	collectAt(t, root, OriginOfficial, rescoped, tools, next())
	d = mustLoad(t, root)
	web, backend := rowOf(d, CapabilityComplexity, "web"), rowOf(d, CapabilityComplexity, "backend")
	if !strings.Contains(web.SparkNote, textSparkRestarted) || !strings.Contains(web.Basis, "not compared") || web.Spark != "" {
		t.Errorf("rescoped instance: spark %q note %q basis %q, want a restart", web.Spark, web.SparkNote, web.Basis)
	}
	if strings.Contains(backend.SparkNote, textSparkRestarted) || backend.Spark == "" {
		t.Errorf("sibling instance was restarted by another's scope: spark %q note %q", backend.Spark, backend.SparkNote)
	}

	// A new provider version restarts the series too, and only for the
	// instance whose provider reported it.
	report = fixtureBytes(t, "duplication", "versioned.json")
	collectAt(t, root, OriginOfficial, rescoped, tools, next())
	d = mustLoad(t, root)
	if dup := rowOf(d, CapabilityDuplication, "dup"); !strings.Contains(dup.SparkNote, textSparkRestarted) || dup.ProviderVersion != "4.0.5" {
		t.Errorf("provider version change: note %q version %q, want a restart on 4.0.5", dup.SparkNote, dup.ProviderVersion)
	}
	if backend := rowOf(d, CapabilityComplexity, "backend"); strings.Contains(backend.SparkNote, textSparkRestarted) {
		t.Errorf("a jscpd version change restarted lizard history: %q", backend.SparkNote)
	}

	// The newest snapshot is manual: it is shown but never counted, and sign-off
	// still refers to the official check.
	before := rowOf(d, CapabilityComplexity, "backend").Spark
	collectAt(t, root, OriginManual, rescoped, tools, next())
	d = mustLoad(t, root)
	if d.Origin != OriginManual || !strings.Contains(d.SignOff, "manual refresh does not affect sign-off") {
		t.Errorf("manual newest: origin %s sign-off %q", d.Origin, d.SignOff)
	}
	if got := rowOf(d, CapabilityComplexity, "backend"); got.Spark != before {
		t.Errorf("manual result changed the trend: spark %q (was %q) basis %q", got.Spark, before, got.Basis)
	}

	// A newer official snapshot takes sign-off back.
	collectAt(t, root, OriginOfficial, rescoped, tools, next())
	d = mustLoad(t, root)
	if d.Origin != OriginOfficial || strings.Contains(d.SignOff, "manual refresh") {
		t.Errorf("official newest: origin %s sign-off %q", d.Origin, d.SignOff)
	}
}
