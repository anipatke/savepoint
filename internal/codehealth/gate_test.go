package codehealth

import (
	"errors"
	"strings"
	"testing"
)

func gateResult(c Capability, p ProviderKey, o Outcome, f Freshness, n float64) CapabilityResult {
	r := CapabilityResult{Capability: c, Outcome: o, Freshness: f, Provenance: Provenance{Provider: p}}
	if o.Measured() {
		r.Value = &Value{Number: n, Unit: capabilityUnits[c]}
	}
	if o == OutcomePartial {
		r.Reason = "two packages missing"
	}
	return r
}

func gateSnapshot(class Classification, results ...CapabilityResult) Snapshot {
	s := Snapshot{ID: "sha256:test", Origin: OriginOfficial}
	for _, r := range results {
		s.Results = append(s.Results, r)
		s.Summary.Capabilities = append(s.Summary.Capabilities, CapabilitySummary{
			Capability: r.Capability, Provider: r.Provenance.Provider, Name: r.Name, Classification: class,
		})
	}
	return s
}

func gateConfig(ccs ...CapabilityConfig) Config {
	return Config{Version: ConfigVersion, Capabilities: ccs}
}

func TestEvaluateRules(t *testing.T) {
	testCfg := CapabilityConfig{Capability: CapabilityTests, Provider: ProviderGoTestJSON}
	covCfg := CapabilityConfig{Capability: CapabilityCoverage, Provider: ProviderGoCoverProfile}
	vulnCfg := CapabilityConfig{Capability: CapabilityDependencyVulnerability, Provider: ProviderOSVScannerJSON}
	req := func(c CapabilityConfig) CapabilityConfig { c.Required = true; return c }
	block := func(c CapabilityConfig) CapabilityConfig { c.Blocking = true; return c }
	vuln := func(key string, n float64) CapabilityResult {
		r := gateResult(CapabilityDependencyVulnerability, ProviderOSVScannerJSON, OutcomeAvailable, FreshnessFresh, n)
		r.Details = []Detail{{Key: key, Number: n}}
		return r
	}
	tests := []struct {
		name  string
		cfg   CapabilityConfig
		res   CapabilityResult
		class Classification
		disp  Disposition
		kind  Kind
	}{
		{"failing optional tests block", testCfg, gateResult(CapabilityTests, ProviderGoTestJSON, OutcomeAvailable, FreshnessFresh, 2), ClassificationNeedsAttention, DispositionBlocks, KindUnhealthy},
		{"passing tests report no finding", testCfg, gateResult(CapabilityTests, ProviderGoTestJSON, OutcomeAvailable, FreshnessFresh, 0), ClassificationGood, DispositionReported, KindNoFinding},
		{"stale required tests block", req(testCfg), gateResult(CapabilityTests, ProviderGoTestJSON, OutcomeAvailable, FreshnessStale, 0), ClassificationUnknown, DispositionBlocks, KindStale},
		{"stale optional tests are reported", testCfg, gateResult(CapabilityTests, ProviderGoTestJSON, OutcomeAvailable, FreshnessStale, 0), ClassificationUnknown, DispositionReported, KindStale},
		{"failed optional coverage is reported", covCfg, gateResult(CapabilityCoverage, ProviderGoCoverProfile, OutcomeFailed, FreshnessUnknown, 0), ClassificationUnknown, DispositionReported, KindCollectionFailure},
		{"failed required coverage blocks", req(covCfg), gateResult(CapabilityCoverage, ProviderGoCoverProfile, OutcomeFailed, FreshnessUnknown, 0), ClassificationUnknown, DispositionBlocks, KindCollectionFailure},
		{"timed out required blocks", req(covCfg), gateResult(CapabilityCoverage, ProviderGoCoverProfile, OutcomeTimedOut, FreshnessUnknown, 0), ClassificationUnknown, DispositionBlocks, KindCollectionFailure},
		{"unavailable required blocks", req(covCfg), gateResult(CapabilityCoverage, ProviderGoCoverProfile, OutcomeUnavailable, FreshnessUnknown, 0), ClassificationUnknown, DispositionBlocks, KindCollectionFailure},
		{"absent required blocks", req(covCfg), gateResult(CapabilityCoverage, ProviderGoCoverProfile, OutcomeAbsent, FreshnessUnknown, 0), ClassificationUnknown, DispositionBlocks, KindCollectionFailure},
		{"unsupported required blocks", req(covCfg), gateResult(CapabilityCoverage, ProviderGoCoverProfile, OutcomeUnsupported, FreshnessUnknown, 0), ClassificationUnknown, DispositionBlocks, KindCollectionFailure},
		{"cancelled required blocks", req(covCfg), gateResult(CapabilityCoverage, ProviderGoCoverProfile, OutcomeCancelled, FreshnessUnknown, 0), ClassificationUnknown, DispositionBlocks, KindCollectionFailure},
		{"cancelled optional is reported", covCfg, gateResult(CapabilityCoverage, ProviderGoCoverProfile, OutcomeCancelled, FreshnessUnknown, 0), ClassificationUnknown, DispositionReported, KindCollectionFailure},
		{"needs attention without flag is reported", covCfg, gateResult(CapabilityCoverage, ProviderGoCoverProfile, OutcomeAvailable, FreshnessFresh, 40), ClassificationNeedsAttention, DispositionReported, KindNoFinding},
		{"needs attention with flag blocks", block(covCfg), gateResult(CapabilityCoverage, ProviderGoCoverProfile, OutcomeAvailable, FreshnessFresh, 40), ClassificationNeedsAttention, DispositionBlocks, KindUnhealthy},
		{"watch with flag does not block", block(covCfg), gateResult(CapabilityCoverage, ProviderGoCoverProfile, OutcomeAvailable, FreshnessFresh, 70), ClassificationWatch, DispositionReported, KindNoFinding},
		{"partial required alone does not block", req(covCfg), gateResult(CapabilityCoverage, ProviderGoCoverProfile, OutcomePartial, FreshnessFresh, 90), ClassificationWatch, DispositionReported, KindIncomplete},
		{"partial with opt-in needs attention blocks", block(covCfg), gateResult(CapabilityCoverage, ProviderGoCoverProfile, OutcomePartial, FreshnessFresh, 30), ClassificationNeedsAttention, DispositionBlocks, KindUnhealthy},
		{"partial tests with failures block", testCfg, gateResult(CapabilityTests, ProviderGoTestJSON, OutcomePartial, FreshnessFresh, 1), ClassificationNeedsAttention, DispositionBlocks, KindUnhealthy},
		{"critical vulnerability blocks", vulnCfg, vuln(DetailCriticalVulnerabilities, 1), ClassificationNeedsAttention, DispositionBlocks, KindUnhealthy},
		{"high vulnerability blocks", vulnCfg, vuln(DetailHighVulnerabilities, 1), ClassificationNeedsAttention, DispositionBlocks, KindUnhealthy},
		{"medium vulnerability is reported", vulnCfg, vuln(DetailMediumVulnerabilities, 1), ClassificationWatch, DispositionReported, KindNoFinding},
		{"unknown severity needs review without blocking", vulnCfg, vuln(DetailUnknownVulnerabilities, 2), ClassificationWatch, DispositionReported, KindUnhealthy},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v, err := Evaluate(gateSnapshot(tc.class, tc.res), gateConfig(tc.cfg))
			if err != nil {
				t.Fatal(err)
			}
			if len(v.Results) != 1 {
				t.Fatalf("got %d verdicts, want 1", len(v.Results))
			}
			got := v.Results[0]
			if got.Disposition != tc.disp || got.Kind != tc.kind {
				t.Errorf("got %s/%s (%s), want %s/%s", got.Disposition, got.Kind, got.Reason, tc.disp, tc.kind)
			}
			if v.Blocks() != (tc.disp == DispositionBlocks) {
				t.Errorf("Blocks() = %v", v.Blocks())
			}
		})
	}
}

func TestEvaluateNotConfiguredAndMissing(t *testing.T) {
	cov := CapabilityConfig{Capability: CapabilityCoverage, Provider: ProviderGoCoverProfile, Required: true}
	s := gateSnapshot(ClassificationUnknown, CapabilityResult{Capability: CapabilityTests, Outcome: OutcomeNotConfigured, Freshness: FreshnessUnknown})
	v, err := Evaluate(s, gateConfig(cov))
	if err != nil {
		t.Fatal(err)
	}
	if len(v.Results) != 2 {
		t.Fatalf("got %d verdicts, want 2", len(v.Results))
	}
	if v.Results[0].Disposition != DispositionBlocks || v.Results[0].Kind != KindCollectionFailure {
		t.Errorf("missing required instance = %+v", v.Results[0])
	}
	if v.Results[1].Capability != CapabilityTests || v.Results[1].Disposition != DispositionNotConfigured {
		t.Errorf("not_configured verdict = %+v", v.Results[1])
	}
}

func TestEvaluateRefusesManualSnapshot(t *testing.T) {
	s := gateSnapshot(ClassificationGood)
	s.Origin = OriginManual
	if _, err := Evaluate(s, gateConfig()); !errors.Is(err, ErrManualSnapshot) {
		t.Fatalf("err = %v, want ErrManualSnapshot", err)
	}
}

func TestBlockingFlagValidation(t *testing.T) {
	for _, c := range []CapabilityConfig{
		{Capability: CapabilityTests, Provider: ProviderGoTestJSON, Blocking: true},
		{Capability: CapabilityDependencyVulnerability, Provider: ProviderOSVScannerJSON, Blocking: true},
	} {
		if err := gateConfig(c).Validate(); !errors.Is(err, ErrBlockingNotAllowed) {
			t.Errorf("%s: err = %v, want ErrBlockingNotAllowed", c.Capability, err)
		}
	}
	for _, c := range []CapabilityConfig{
		{Capability: CapabilityCoverage, Provider: ProviderGoCoverProfile, Blocking: true},
		{Capability: CapabilityComplexity, Provider: ProviderLizardCSV, Blocking: true},
		{Capability: CapabilityDuplication, Provider: ProviderJscpdJSON, Blocking: true},
	} {
		if err := gateConfig(c).Validate(); err != nil {
			t.Errorf("%s: %v", c.Capability, err)
		}
	}
}

func TestBlockingFlagDoesNotChangeDigest(t *testing.T) {
	a := CapabilityConfig{Capability: CapabilityCoverage, Provider: ProviderGoCoverProfile}
	b := a
	b.Blocking = true
	if a.Digest() != b.Digest() {
		t.Fatal("setting blocking changed the digest, which would reset comparison series")
	}
	if _, err := DecodeConfig([]byte(`{"version":1,"capabilities":[{"capability":"coverage","provider":"go-cover-profile"}]}`)); err != nil {
		t.Fatalf("config without blocking: %v", err)
	}
}

func TestRenderGolden(t *testing.T) {
	cfg := gateConfig(
		CapabilityConfig{Capability: CapabilityTests, Provider: ProviderGoTestJSON, Required: true},
		CapabilityConfig{Capability: CapabilityCoverage, Provider: ProviderGoCoverProfile},
	)
	s := gateSnapshot(ClassificationUnknown,
		gateResult(CapabilityTests, ProviderGoTestJSON, OutcomeAvailable, FreshnessFresh, 2),
		gateResult(CapabilityCoverage, ProviderGoCoverProfile, OutcomeFailed, FreshnessUnknown, 0),
		CapabilityResult{Capability: CapabilityDuplication, Outcome: OutcomeNotConfigured, Freshness: FreshnessUnknown},
	)
	v, err := Evaluate(s, cfg)
	if err != nil {
		t.Fatal(err)
	}
	want := `Code Health blocks clearance.
- coverage via go-cover-profile [optional]: reported only; collection failure: collection failed.
- duplication: not configured
- tests via go-test-json [required]: blocks clearance; unhealthy measurement: 2 failing.
`
	if got := v.Render(); got != want {
		t.Errorf("render:\n%s\nwant:\n%s", got, want)
	}

	ok, err := Evaluate(gateSnapshot(ClassificationGood, gateResult(CapabilityTests, ProviderGoTestJSON, OutcomeAvailable, FreshnessFresh, 0)), cfg)
	if err != nil {
		t.Fatal(err)
	}
	out := ok.Render()
	if !strings.HasPrefix(out, "Code Health does not block clearance.\n") {
		t.Errorf("headline = %q", out)
	}
	for _, banned := range []string{"healthy", "Clear ", "is clear"} {
		if strings.Contains(out, banned) {
			t.Errorf("render contains %q", banned)
		}
	}
	if (Verdict{}).Render() != "Code Health has no results.\n" {
		t.Errorf("empty render = %q", Verdict{}.Render())
	}
}

// A blocking finding must not hide that the evidence is also stale or incomplete.
func TestEvaluateKeepsEvidenceLimitations(t *testing.T) {
	testCfg := CapabilityConfig{Capability: CapabilityTests, Provider: ProviderGoTestJSON, Required: true}
	covCfg := CapabilityConfig{Capability: CapabilityCoverage, Provider: ProviderGoCoverProfile, Blocking: true}
	vulnCfg := CapabilityConfig{Capability: CapabilityDependencyVulnerability, Provider: ProviderOSVScannerJSON}
	vuln := func(o Outcome, f Freshness, details ...Detail) CapabilityResult {
		r := gateResult(CapabilityDependencyVulnerability, ProviderOSVScannerJSON, o, f, 1)
		r.Details = details
		return r
	}
	partialTests := func(f Freshness) CapabilityResult {
		r := gateResult(CapabilityTests, ProviderGoTestJSON, OutcomePartial, f, 1)
		r.Reason = "package example/missing has no terminal event"
		return r
	}
	tests := []struct {
		name  string
		cfg   CapabilityConfig
		res   CapabilityResult
		class Classification
		want  string
		kind  Kind
	}{
		{"fresh partial failing tests", testCfg, partialTests(FreshnessFresh), ClassificationNeedsAttention,
			"1 failing; the evidence is incomplete: package example/missing has no terminal event", KindUnhealthy},
		{"stale partial failing tests", testCfg, partialTests(FreshnessStale), ClassificationNeedsAttention,
			"1 failing; the evidence is stale; the evidence is incomplete: package example/missing has no terminal event", KindUnhealthy},
		{"stale failing tests", testCfg, gateResult(CapabilityTests, ProviderGoTestJSON, OutcomeAvailable, FreshnessStale, 3), ClassificationNeedsAttention,
			"3 failing; the evidence is stale", KindUnhealthy},
		{"partial opt-in coverage needs attention", covCfg, gateResult(CapabilityCoverage, ProviderGoCoverProfile, OutcomePartial, FreshnessFresh, 30), ClassificationNeedsAttention,
			"needs attention and this instance is set to block; the evidence is incomplete: two packages missing", KindUnhealthy},
		{"severe vulnerability on partial evidence", vulnCfg, vuln(OutcomePartial, FreshnessFresh, Detail{Key: DetailCriticalVulnerabilities, Number: 1}), ClassificationNeedsAttention,
			"1 known critical or high severity; the evidence is incomplete: two packages missing", KindUnhealthy},
		{"severe and unknown vulnerabilities, stale", vulnCfg, vuln(OutcomeAvailable, FreshnessStale, Detail{Key: DetailHighVulnerabilities, Number: 1}, Detail{Key: DetailUnknownVulnerabilities, Number: 2}), ClassificationNeedsAttention,
			"1 known critical or high severity; the evidence is stale; 2 of unknown severity need review", KindUnhealthy},
		{"unknown severity on partial evidence keeps review wording", vulnCfg, vuln(OutcomePartial, FreshnessFresh, Detail{Key: DetailUnknownVulnerabilities, Number: 2}), ClassificationWatch,
			"the evidence is incomplete: two packages missing; 2 of unknown severity need review", KindIncomplete},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v, err := Evaluate(gateSnapshot(tc.class, tc.res), gateConfig(tc.cfg))
			if err != nil {
				t.Fatal(err)
			}
			got := v.Results[0]
			if got.Reason != tc.want || got.Kind != tc.kind {
				t.Errorf("got %s: %q, want %s: %q", got.Kind, got.Reason, tc.kind, tc.want)
			}
		})
	}
}

func TestRenderSaysWhenNothingWasMeasured(t *testing.T) {
	cfg := gateConfig(
		CapabilityConfig{Capability: CapabilityComplexity, Provider: ProviderLizardCSV, Executable: "lizard"},
		CapabilityConfig{Capability: CapabilityTests, Provider: ProviderVitestJUnit, Report: "junit.xml"},
		CapabilityConfig{Capability: CapabilityTests, Provider: ProviderPytestJUnit, Report: "pytest-junit.xml"},
	)
	s := gateSnapshot(ClassificationUnknown,
		gateResult(CapabilityComplexity, ProviderLizardCSV, OutcomeUnavailable, FreshnessUnknown, 0),
		gateResult(CapabilityTests, ProviderVitestJUnit, OutcomeAbsent, FreshnessUnknown, 0),
		gateResult(CapabilityTests, ProviderPytestJUnit, OutcomeAbsent, FreshnessUnknown, 0),
		CapabilityResult{Capability: CapabilityCoverage, Outcome: OutcomeNotConfigured, Freshness: FreshnessUnknown},
	)
	for i := range s.Results {
		s.Results[i].Name = []string{"", "a", "b", ""}[i]
	}
	cfg.Capabilities[1].Name, cfg.Capabilities[2].Name = "a", "b"
	v, err := Evaluate(s, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if !v.NoData() {
		t.Fatal("NoData() = false for a report in which nothing was measured")
	}
	out := v.Render()
	for _, want := range []string{
		"Code Health does not block clearance.\nNo signal produced data, so nothing was judged. Do not read this as a pass.\n",
		"the tool was unavailable. Install lizard (pip install lizard), then run the check again.",
		"no report was found. Run your tests so they write the report first, for example: vitest run --reporter=junit --outputFile=junit.xml.",
		"for example: pytest --junitxml=pytest-junit.xml.",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("render is missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "healthy") {
		t.Errorf("render says healthy:\n%s", out)
	}
}

func TestNoDataLineAbsentOnceAnythingIsMeasured(t *testing.T) {
	cfg := gateConfig(
		CapabilityConfig{Capability: CapabilityTests, Provider: ProviderGoTestJSON},
		CapabilityConfig{Capability: CapabilityCoverage, Provider: ProviderGoCoverProfile},
	)
	s := gateSnapshot(ClassificationGood,
		gateResult(CapabilityTests, ProviderGoTestJSON, OutcomeAvailable, FreshnessFresh, 0),
		gateResult(CapabilityCoverage, ProviderGoCoverProfile, OutcomeAbsent, FreshnessUnknown, 0),
	)
	v, err := Evaluate(s, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if v.NoData() || strings.Contains(v.Render(), "No signal produced data") {
		t.Errorf("one measured result must clear the no-data line:\n%s", v.Render())
	}
	if (Verdict{}).NoData() {
		t.Error("an empty verdict is the no-results case, not no-data")
	}
}

func TestRemedyOmittedWhenNothingSpecificCanBeSaid(t *testing.T) {
	r := ResultVerdict{Provider: ProviderGoTestJSON, Kind: KindCollectionFailure, Outcome: OutcomeFailed}
	if got := r.remedy(); got != "" {
		t.Errorf("remedy for a generic failure = %q", got)
	}
	r = ResultVerdict{Provider: ProviderLizardCSV, Kind: KindNoFinding, Outcome: OutcomeUnavailable}
	if got := r.remedy(); got != "" {
		t.Errorf("remedy for a measured result = %q", got)
	}
}
