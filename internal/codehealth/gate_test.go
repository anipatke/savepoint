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
