package codehealth

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var dashProviders = map[Capability]ProviderKey{
	CapabilityTests:                   ProviderGoTestJSON,
	CapabilityCoverage:                ProviderGoCoverProfile,
	CapabilityComplexity:              ProviderLizardCSV,
	CapabilityDuplication:             ProviderJscpdJSON,
	CapabilityDependencyVulnerability: ProviderOSVScannerJSON,
}

// bannedClaims are phrases the dashboard must never use: it reports
// measurements, not proof.
var bannedClaims = []string{
	"proves", "proven", "guarantee", "safe", "secure", "correct", "bug-free",
	"maintainable", "production ready", "production-ready", "healthy",
}

func dashProject(t *testing.T) (Store, string) {
	t.Helper()
	return newProject(t)
}

func dashConfig(t *testing.T, store Store, caps ...Capability) Config {
	t.Helper()
	cfg := Config{Version: ConfigVersion}
	for _, c := range caps {
		cfg.Capabilities = append(cfg.Capabilities, CapabilityConfig{Capability: c, Provider: dashProviders[c], Scope: []string{"**/*.go"}})
	}
	if _, err := store.SaveConfig(cfg); err != nil {
		t.Fatalf("SaveConfig: %v", err)
	}
	return cfg
}

func allCapabilities() []Capability { return Capabilities() }

// dashResult is a complete, fresh, measured result for the instance cc.
func dashResult(cc CapabilityConfig, n float64) CapabilityResult {
	return CapabilityResult{
		Capability: cc.Capability, Name: cc.Name, Outcome: OutcomeAvailable, Freshness: FreshnessFresh,
		CollectedAt:  "2026-10-01T00:00:00Z",
		Provenance:   Provenance{Provider: cc.Provider, ProviderVersion: "1.0", MeasurementDefinition: "def/1"},
		ConfigDigest: cc.Digest(), Scope: cc.Scope,
		Value:    &Value{Number: n, Unit: capabilityUnits[cc.Capability]},
		Evidence: []EvidenceRef{{Path: "report.out", Note: "report"}},
	}
}

func unmeasured(r CapabilityResult, o Outcome) CapabilityResult {
	r.Outcome, r.Freshness, r.Value, r.Evidence = o, FreshnessUnknown, nil, nil
	return r
}

func notConfiguredResult(c Capability) CapabilityResult {
	return CapabilityResult{Capability: c, Outcome: OutcomeNotConfigured, Freshness: FreshnessUnknown, CollectedAt: "2026-10-01T00:00:00Z"}
}

// goodValues are complete good measurements for each capability.
var goodValues = map[Capability]float64{
	CapabilityTests: 0, CapabilityCoverage: 90, CapabilityComplexity: 5,
	CapabilityDuplication: 1, CapabilityDependencyVulnerability: 0,
}

// saveDash stores a snapshot of results, classified against the history
// already stored, exactly as Collect would.
func saveDash(t *testing.T, store Store, cfg Config, origin Origin, n int, results []CapabilityResult, repo RepositoryIdentity) Snapshot {
	t.Helper()
	snaps, err := store.LoadSnapshots()
	if err != nil {
		t.Fatal(err)
	}
	history := dashboardHistory(snaps)
	thresholds := map[instanceKey]*Threshold{}
	for _, cc := range cfg.Capabilities {
		thresholds[instanceKey{cc.Capability, cc.Provider, cc.Name}] = cc.Thresholds
	}
	retention := RetentionPrunable
	if origin == OriginOfficial {
		retention = RetentionPermanent
	}
	s := Snapshot{
		Version: SnapshotVersion, Origin: origin, Retention: retention, Repository: repo,
		CreatedAt: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(n) * time.Minute).Format(timestampLayout),
	}
	var assessed []Assessment
	for _, r := range results {
		a := Assess(origin, r, thresholds[r.key()], ownHistory(r, history[r.historyKey()]))
		assessed = append(assessed, a)
		s.Results = append(s.Results, r)
		s.Summary.Capabilities = append(s.Summary.Capabilities, a.Summary())
	}
	s.Summary.Overall = Overall(assessed)
	s.ID = s.ComputeID()
	mustSave(t, store, s)
	return s
}

func dashRepo(n int) RepositoryIdentity {
	return RepositoryIdentity{Commit: fmt.Sprintf("%040x", n), InputFingerprint: "sha256:" + strings.Repeat("1", 64)}
}

// goodResults is one good measured result per configured instance, with a
// not-configured placeholder for each capability that has none.
func goodResults(cfg Config, tweak func(*CapabilityResult)) []CapabilityResult {
	var out []CapabilityResult
	have := map[Capability]bool{}
	for _, cc := range cfg.Capabilities {
		have[cc.Capability] = true
		r := dashResult(cc, goodValues[cc.Capability])
		if tweak != nil {
			tweak(&r)
		}
		out = append(out, r)
	}
	for _, c := range Capabilities() {
		if !have[c] {
			out = append(out, notConfiguredResult(c))
		}
	}
	return out
}

func rowFor(t *testing.T, d Dashboard, c Capability) DashboardRow {
	t.Helper()
	for _, r := range d.Rows {
		if r.Capability == c {
			return r
		}
	}
	t.Fatalf("no row for %s", c)
	return DashboardRow{}
}

func mustLoad(t *testing.T, root string) Dashboard {
	t.Helper()
	d, err := LoadDashboard(root)
	if err != nil {
		t.Fatalf("LoadDashboard: %v", err)
	}
	return d
}

// dashboardText gathers every phrase a dashboard produced.
func dashboardText(d Dashboard) string {
	parts := []string{d.OverallText, d.OriginText}
	for _, r := range d.Rows {
		parts = append(parts, r.CapabilityText, r.LabelText, r.Explanation, r.Trend, r.Basis, r.OutcomeText, r.FreshnessText, r.RequiredText)
	}
	for _, h := range d.History {
		parts = append(parts, h.OriginText, h.OverallText)
	}
	return strings.ToLower(strings.Join(parts, "\n"))
}

func assertNoBannedClaims(t *testing.T, name, text string) {
	t.Helper()
	for _, banned := range bannedClaims {
		if strings.Contains(text, banned) {
			t.Errorf("%s: text contains banned claim %q:\n%s", name, banned, text)
		}
	}
}

func TestLoadDashboardNotConfigured(t *testing.T) {
	_, root := dashProject(t)
	d := mustLoad(t, root)
	if d.State != DashboardNotConfigured || len(d.Rows) != 0 {
		t.Fatalf("dashboard = %+v, want not configured and empty", d)
	}
}

func TestLoadDashboardFirstRun(t *testing.T) {
	store, root := dashProject(t)
	dashConfig(t, store, CapabilityTests)
	d := mustLoad(t, root)
	if d.State != DashboardFirstRun || len(d.Rows) != 0 {
		t.Fatalf("dashboard = %+v, want first run and empty", d)
	}
}

func TestLoadDashboardNeverRunsAnythingOrWrites(t *testing.T) {
	store, root := dashProject(t)
	cfg := dashConfig(t, store, allCapabilities()...)
	saveDash(t, store, cfg, OriginOfficial, 1, goodResults(cfg, nil), dashRepo(1))
	before := snapshotFiles(t, root)
	mustLoad(t, root)
	if after := snapshotFiles(t, root); strings.Join(after, ",") != strings.Join(before, ",") {
		t.Fatalf("loading changed the snapshot files: %v -> %v", before, after)
	}
}

func TestLoadDashboardReportsDamagedStorage(t *testing.T) {
	store, root := dashProject(t)
	dashConfig(t, store, CapabilityTests)
	writeFile(t, filepath.Join(healthPath(root, snapshotsDir), strings.Repeat("a", 64)+".json"), "{")
	if _, err := LoadDashboard(root); err == nil {
		t.Fatal("LoadDashboard on a damaged snapshot = nil error, want a named diagnostic")
	}
	writeFile(t, healthPath(root, configFile), "{")
	if _, err := LoadDashboard(root); err == nil {
		t.Fatal("LoadDashboard on a damaged configuration = nil error")
	}
}

func TestDashboardAllFiveSignalsMeasured(t *testing.T) {
	store, root := dashProject(t)
	cfg := dashConfig(t, store, allCapabilities()...)
	for n := 1; n <= 3; n++ {
		saveDash(t, store, cfg, OriginOfficial, n, goodResults(cfg, nil), dashRepo(n))
	}
	d := mustLoad(t, root)
	if d.State != DashboardMeasured || d.Origin != OriginOfficial || d.Overall != ClassificationGood || d.OverallText != "Good" {
		t.Fatalf("dashboard header = %+v", d)
	}
	if len(d.Rows) != 5 {
		t.Fatalf("rows = %d, want 5", len(d.Rows))
	}
	for i, c := range Capabilities() {
		r := d.Rows[i]
		if r.Capability != c || r.Label != ClassificationGood || r.LabelText != "Good" {
			t.Errorf("row %d = %s %s, want %s Good", i, r.Capability, r.Label, c)
		}
		if r.Outcome != OutcomeAvailable || r.OutcomeText != "measured" || r.FreshnessText != "fresh" {
			t.Errorf("%s outcome = %q, %q", c, r.OutcomeText, r.FreshnessText)
		}
		if r.Provider != dashProviders[c] || r.ProviderVersion != "1.0" || r.CollectedAt == "" || len(r.Scope) != 1 || len(r.Evidence) != 1 {
			t.Errorf("%s provenance/scope/evidence = %+v", c, r)
		}
		if r.RequiredText != "optional" || r.Required {
			t.Errorf("%s required = %v %q", c, r.Required, r.RequiredText)
		}
		if !strings.HasPrefix(r.Trend, "Steady over 3 official checks") || !strings.Contains(r.Basis, "2 earlier comparable official checks") {
			t.Errorf("%s trend/basis = %q / %q", c, r.Trend, r.Basis)
		}
	}
	assertNoBannedClaims(t, "all measured", dashboardText(d))
}

func TestDashboardLabelsComeFromThePersistedSummary(t *testing.T) {
	store, root := dashProject(t)
	cfg := dashConfig(t, store, allCapabilities()...)
	s := saveDash(t, store, cfg, OriginOfficial, 1, goodResults(cfg, nil), dashRepo(1))
	d := mustLoad(t, root)
	if d.Overall != s.Summary.Overall {
		t.Errorf("overall = %s, want persisted %s", d.Overall, s.Summary.Overall)
	}
	for i, cs := range s.Summary.Capabilities {
		if r := d.Rows[i]; r.Label != cs.Classification || r.Explanation != cs.Explanation {
			t.Errorf("row %d = %s %q, want persisted %s %q", i, r.Label, r.Explanation, cs.Classification, cs.Explanation)
		}
	}
	// One official check is too little history for Good.
	if d.Overall == ClassificationGood {
		t.Error("a first official check must not read as Good")
	}
}

func TestDashboardNotConfiguredCapabilitiesGetPlaceholderRows(t *testing.T) {
	store, root := dashProject(t)
	cfg := dashConfig(t, store, CapabilityTests, CapabilityCoverage)
	saveDash(t, store, cfg, OriginManual, 1, goodResults(cfg, nil), dashRepo(1))
	d := mustLoad(t, root)
	if len(d.Rows) != 5 {
		t.Fatalf("rows = %d, want 5", len(d.Rows))
	}
	for _, c := range []Capability{CapabilityComplexity, CapabilityDuplication, CapabilityDependencyVulnerability} {
		r := rowFor(t, d, c)
		if !r.NotConfigured || r.Label != ClassificationUnknown || r.OutcomeText != "not configured" {
			t.Errorf("%s = %+v, want a not-configured Unknown row", c, r)
		}
	}
	if d.Overall == ClassificationGood {
		t.Error("overall is Good with unconfigured signals")
	}
}

func TestDashboardMultipleInstancesOfOneCapabilityEachGetARow(t *testing.T) {
	store, root := dashProject(t)
	cfg := Config{Version: ConfigVersion, Capabilities: []CapabilityConfig{
		{Capability: CapabilityComplexity, Provider: ProviderLizardCSV, Name: "api", Scope: []string{"api/**"}},
		{Capability: CapabilityComplexity, Provider: ProviderLizardCSV, Name: "web", Scope: []string{"web/**"}},
	}}
	if _, err := store.SaveConfig(cfg); err != nil {
		t.Fatal(err)
	}
	saveDash(t, store, cfg, OriginManual, 1, goodResults(cfg, nil), dashRepo(1))
	var names []string
	for _, r := range mustLoad(t, root).Rows {
		if r.Capability == CapabilityComplexity {
			names = append(names, r.Name)
		}
	}
	if strings.Join(names, ",") != "api,web" {
		t.Fatalf("complexity rows = %v, want api,web in configuration order", names)
	}
}

func TestDashboardInstanceMissingFromSnapshotIsUnknown(t *testing.T) {
	store, root := dashProject(t)
	cfg := dashConfig(t, store, CapabilityTests)
	saveDash(t, store, cfg, OriginManual, 1, goodResults(cfg, nil), dashRepo(1))
	cfg = dashConfig(t, store, CapabilityTests, CapabilityCoverage)
	r := rowFor(t, mustLoad(t, root), CapabilityCoverage)
	if r.NotConfigured || r.Label != ClassificationUnknown || !strings.Contains(r.Explanation, "refresh") {
		t.Fatalf("coverage row = %+v, want Unknown and a refresh pointer", r)
	}
}

func TestDashboardPartialSupportStaysVisibleAndNotGood(t *testing.T) {
	store, root := dashProject(t)
	cfg := dashConfig(t, store, allCapabilities()...)
	results := goodResults(cfg, func(r *CapabilityResult) {
		if r.Capability == CapabilityCoverage {
			r.Outcome, r.Reason = OutcomePartial, "report truncated"
		}
	})
	saveDash(t, store, cfg, OriginOfficial, 1, results, dashRepo(1))
	d := mustLoad(t, root)
	r := rowFor(t, d, CapabilityCoverage)
	if r.OutcomeText != "partial" || r.Label == ClassificationGood || r.Explanation == "" {
		t.Fatalf("coverage = %+v, want a visible partial row that is not Good", r)
	}
	if d.Overall == ClassificationGood {
		t.Error("overall Good with a partial result")
	}
	assertNoBannedClaims(t, "partial", dashboardText(d))
}

func TestDashboardFailuresStayDistinctAndRequiredIsShown(t *testing.T) {
	store, root := dashProject(t)
	cfg := dashConfig(t, store, allCapabilities()...)
	cfg.Capabilities[0].Required = true // tests
	if _, err := store.SaveConfig(cfg); err != nil {
		t.Fatal(err)
	}
	outcomes := map[Capability]Outcome{
		CapabilityTests:                   OutcomeFailed,
		CapabilityCoverage:                OutcomeAbsent,
		CapabilityComplexity:              OutcomeTimedOut,
		CapabilityDuplication:             OutcomeUnavailable,
		CapabilityDependencyVulnerability: OutcomeUnsupported,
	}
	results := goodResults(cfg, nil)
	for i := range results {
		results[i] = unmeasured(results[i], outcomes[results[i].Capability])
	}
	saveDash(t, store, cfg, OriginOfficial, 1, results, dashRepo(1))
	d := mustLoad(t, root)
	want := map[Capability]string{
		CapabilityTests: "failed", CapabilityCoverage: "absent", CapabilityComplexity: "timed out",
		CapabilityDuplication: "unavailable", CapabilityDependencyVulnerability: "unsupported",
	}
	for c, w := range want {
		r := rowFor(t, d, c)
		if !strings.HasPrefix(r.OutcomeText, w) || r.Label != ClassificationUnknown || r.FreshnessText != "" {
			t.Errorf("%s = %q %s, want %q Unknown", c, r.OutcomeText, r.Label, w)
		}
	}
	if r := rowFor(t, d, CapabilityTests); !r.Required || r.RequiredText != "required" {
		t.Errorf("tests required = %v %q", r.Required, r.RequiredText)
	}
	if r := rowFor(t, d, CapabilityCoverage); r.Required || r.RequiredText != "optional" {
		t.Errorf("coverage required = %v %q", r.Required, r.RequiredText)
	}
	if d.Overall == ClassificationGood {
		t.Error("overall Good with failed results")
	}
	assertNoBannedClaims(t, "failures", dashboardText(d))
}

func TestDashboardNeedsAttentionForFailingTests(t *testing.T) {
	store, root := dashProject(t)
	cfg := dashConfig(t, store, allCapabilities()...)
	results := goodResults(cfg, func(r *CapabilityResult) {
		if r.Capability == CapabilityTests {
			r.Value.Number = 3
		}
	})
	saveDash(t, store, cfg, OriginOfficial, 1, results, dashRepo(1))
	d := mustLoad(t, root)
	if r := rowFor(t, d, CapabilityTests); r.Label != ClassificationNeedsAttention || r.LabelText != "Needs Attention" {
		t.Errorf("tests = %s, want Needs Attention", r.Label)
	}
	if d.Overall != ClassificationNeedsAttention {
		t.Errorf("overall = %s, want needs_attention", d.Overall)
	}
}

func TestDashboardManualNewestWinsOverOfficial(t *testing.T) {
	store, root := dashProject(t)
	cfg := dashConfig(t, store, allCapabilities()...)
	saveDash(t, store, cfg, OriginOfficial, 1, goodResults(cfg, nil), dashRepo(1))
	manual := saveDash(t, store, cfg, OriginManual, 2, goodResults(cfg, nil), dashRepo(2))
	d := mustLoad(t, root)
	if d.Origin != OriginManual || d.OriginText != "Manual refresh" || d.SnapshotID != manual.ID {
		t.Fatalf("dashboard shows %s %s, want the newer manual snapshot", d.Origin, d.SnapshotID)
	}
	if d.Recorded != manual.Repository {
		t.Errorf("recorded = %+v, want the manual snapshot's repository", d.Recorded)
	}
}

func TestDashboardComparableTrendAndReset(t *testing.T) {
	store, root := dashProject(t)
	cfg := dashConfig(t, store, CapabilityCoverage)
	for n, v := range []float64{90, 88, 86} {
		saveDash(t, store, cfg, OriginOfficial, n+1, []CapabilityResult{dashResult(cfg.Capabilities[0], v)}, dashRepo(n+1))
	}
	if r := rowFor(t, mustLoad(t, root), CapabilityCoverage); !strings.Contains(r.Trend, "over 3 official checks") || !strings.Contains(r.Basis, "2 earlier comparable official checks") {
		t.Fatalf("comparable trend/basis = %q / %q", r.Trend, r.Basis)
	}

	// A new tool version starts a new series: earlier results are not compared.
	reset := dashResult(cfg.Capabilities[0], 85)
	reset.Provenance.ProviderVersion = "2.0"
	reset.Outcome = OutcomeAvailable
	saveDash(t, store, cfg, OriginOfficial, 4, append([]CapabilityResult{reset}, notConfiguredFor(CapabilityCoverage)...), dashRepo(4))
	r := rowFor(t, mustLoad(t, root), CapabilityCoverage)
	if r.Trend != "No trend yet" || !strings.Contains(r.Basis, "3 earlier official results not compared") {
		t.Fatalf("reset trend/basis = %q / %q", r.Trend, r.Basis)
	}
}

// notConfiguredFor is a placeholder for every capability except keep.
func notConfiguredFor(keep Capability) []CapabilityResult {
	var out []CapabilityResult
	for _, c := range Capabilities() {
		if c != keep {
			out = append(out, notConfiguredResult(c))
		}
	}
	return out
}

func TestDashboardManualResultsAreShownNotCounted(t *testing.T) {
	store, root := dashProject(t)
	cfg := dashConfig(t, store, CapabilityCoverage)
	cc := cfg.Capabilities[0]
	for n := 1; n <= 2; n++ {
		saveDash(t, store, cfg, OriginManual, n, append([]CapabilityResult{dashResult(cc, 90)}, notConfiguredFor(CapabilityCoverage)...), dashRepo(n))
	}
	r := rowFor(t, mustLoad(t, root), CapabilityCoverage)
	if r.Trend != "No trend yet" || !strings.Contains(r.Basis, "1 manual result shown, not counted") {
		t.Fatalf("trend/basis = %q / %q", r.Trend, r.Basis)
	}
}

func TestDashboardHistoryIsBoundedAndNewestFirst(t *testing.T) {
	store, root := dashProject(t)
	cfg := dashConfig(t, store, allCapabilities()...)
	for n := 1; n <= MaxDashboardHistory+2; n++ {
		origin := OriginManual
		if n%2 == 0 {
			origin = OriginOfficial
		}
		saveDash(t, store, cfg, origin, n, goodResults(cfg, nil), dashRepo(n))
	}
	d := mustLoad(t, root)
	if len(d.History) != MaxDashboardHistory {
		t.Fatalf("history = %d entries, want %d", len(d.History), MaxDashboardHistory)
	}
	for i := 1; i < len(d.History); i++ {
		if d.History[i-1].CreatedAt <= d.History[i].CreatedAt {
			t.Fatalf("history not newest first at %d: %s then %s", i, d.History[i-1].CreatedAt, d.History[i].CreatedAt)
		}
	}
	first := d.History[0]
	if first.Origin != OriginOfficial || first.OriginText != "Official check" || first.OverallText == "" || first.CreatedAt != d.MeasuredAt {
		t.Errorf("newest history entry = %+v", first)
	}
}

func TestDashboardWordingTableHasNoBannedClaims(t *testing.T) {
	var all []string
	for _, v := range classificationText {
		all = append(all, v)
	}
	for _, v := range originText {
		all = append(all, v)
	}
	for _, v := range outcomeText {
		all = append(all, v)
	}
	for _, v := range freshnessText {
		all = append(all, v)
	}
	for _, v := range capabilityText {
		all = append(all, v)
	}
	for _, v := range trendText {
		all = append(all, v)
	}
	all = append(all, textNotConfiguredRow, textNoResultRow, textNoBasis, textBasisFormat, textTrendRange)
	assertNoBannedClaims(t, "wording table", strings.ToLower(strings.Join(all, "\n")))
	for _, c := range Capabilities() {
		if capabilityText[c] == "" {
			t.Errorf("no wording for %s", c)
		}
	}
	for _, o := range Outcomes() {
		if outcomeText[o] == "" {
			t.Errorf("no wording for outcome %s", o)
		}
	}
}
