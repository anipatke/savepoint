package codehealth

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultReadersCoverTheWholeCatalogue(t *testing.T) {
	readers := DefaultReaders()
	for key := range providerCapability {
		if r, ok := readers[key]; !ok || r == nil {
			t.Errorf("catalogue provider %s has no reader", key)
		}
	}
	for key := range readers {
		if _, ok := providerCapability[key]; !ok {
			t.Errorf("reader registered for %s, which is not in the catalogue", key)
		}
	}
}

func fixtureBytes(t *testing.T, dir, name string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "readers", dir, name))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestCollectWithRealReadersKeepsEveryMeasureTruthful(t *testing.T) {
	root := project(t)
	write(t, root, "go.mod", "module example.com/m\n")
	report := func(rel, dir, fixture string) string {
		// Fixtures that name absolute paths are pointed at this project.
		write(t, root, rel, strings.ReplaceAll(fixtureBytes(t, dir, fixture), "/work/proj", filepath.ToSlash(root)))
		return rel
	}

	type want struct {
		capability Capability
		provider   ProviderKey
		name       string
		outcome    Outcome // zero means any measured outcome
	}
	cfg := cfgOf(
		// Report-only providers, each with a valid report.
		CapabilityConfig{Capability: CapabilityTests, Provider: ProviderGoTestJSON, Report: report("reports/go-test.jsonl", "tests", "go-fail.jsonl")},
		CapabilityConfig{Capability: CapabilityTests, Provider: ProviderVitestJUnit, Name: "ok", Scope: []string{"web/**"}, Report: report("reports/vitest.xml", "tests", "vitest-fail.xml")},
		CapabilityConfig{Capability: CapabilityTests, Provider: ProviderVitestJUnit, Name: "broken", Scope: []string{"broken/**"}, Report: report("reports/broken.xml", "tests", "junit-malformed.xml")},
		CapabilityConfig{Capability: CapabilityTests, Provider: ProviderPytestJUnit, Report: report("reports/pytest.xml", "tests", "pytest-fail.xml")},
		CapabilityConfig{Capability: CapabilityCoverage, Provider: ProviderGoCoverProfile, Name: "ok", Scope: []string{"a/**"}, Report: report("reports/go.out", "coverage", "go-populated.out")},
		CapabilityConfig{Capability: CapabilityCoverage, Provider: ProviderGoCoverProfile, Name: "absent", Scope: []string{"absent/**"}, Report: "reports/never-written.out"},
		CapabilityConfig{Capability: CapabilityCoverage, Provider: ProviderVitestV8, Report: report("reports/vitest-cov.json", "coverage", "vitest-populated.json")},
		CapabilityConfig{Capability: CapabilityCoverage, Provider: ProviderCoveragePyJSON, Report: report("reports/coverage.json", "coverage", "coveragepy-populated.json")},
		// Executed providers answered by a fake runner.
		CapabilityConfig{Capability: CapabilityComplexity, Provider: ProviderLizardCSV, Name: "all", Executable: "lizard-all"},
		CapabilityConfig{Capability: CapabilityComplexity, Provider: ProviderLizardCSV, Name: "app", Executable: "lizard-app", Scope: []string{"app/**"}},
		CapabilityConfig{Capability: CapabilityComplexity, Provider: ProviderLizardCSV, Name: "slow", Executable: "lizard-slow", Scope: []string{"slow/**"}},
		CapabilityConfig{Capability: CapabilityDuplication, Provider: ProviderJscpdJSON, Name: "ok", Executable: "jscpd", Scope: []string{"web/**"}},
		CapabilityConfig{Capability: CapabilityDuplication, Provider: ProviderJscpdJSON, Name: "missing", Executable: "jscpd-missing", Scope: []string{"missing/**"}},
		CapabilityConfig{Capability: CapabilityDependencyVulnerability, Provider: ProviderOSVScannerJSON, Executable: "osv-scanner"},
	)
	wants := []want{
		{CapabilityTests, ProviderGoTestJSON, "", ""},
		{CapabilityTests, ProviderVitestJUnit, "ok", ""},
		{CapabilityTests, ProviderVitestJUnit, "broken", OutcomeFailed},
		{CapabilityTests, ProviderPytestJUnit, "", ""},
		{CapabilityCoverage, ProviderGoCoverProfile, "ok", ""},
		{CapabilityCoverage, ProviderGoCoverProfile, "absent", OutcomeAbsent},
		{CapabilityCoverage, ProviderVitestV8, "", ""},
		{CapabilityCoverage, ProviderCoveragePyJSON, "", ""},
		{CapabilityComplexity, ProviderLizardCSV, "all", ""},
		{CapabilityComplexity, ProviderLizardCSV, "app", ""},
		{CapabilityComplexity, ProviderLizardCSV, "slow", OutcomeTimedOut},
		{CapabilityDuplication, ProviderJscpdJSON, "ok", ""},
		{CapabilityDuplication, ProviderJscpdJSON, "missing", OutcomeUnavailable},
		{CapabilityDependencyVulnerability, ProviderOSVScannerJSON, "", ""},
	}

	lizardOut := stdout(fixtureBytes(t, "complexity", "mixed.csv"))
	tools := &fakeTools{t: t, behavior: map[string]func(context.Context, ToolSpec) (ToolResult, error){
		"lizard-all": lizardOut,
		"lizard-app": lizardOut,
		"lizard-slow": func(context.Context, ToolSpec) (ToolResult, error) {
			return ToolResult{}, context.DeadlineExceeded
		},
		"jscpd": stdout(fixtureBytes(t, "duplication", "mixed.json")),
		// osv-scanner exits 1 when it finds vulnerabilities.
		"osv-scanner": func(context.Context, ToolSpec) (ToolResult, error) {
			return ToolResult{Stdout: []byte(fixtureBytes(t, "vulnerabilities", "all-severities.json")), ExitCode: 1}, nil
		},
	}}

	got := collect(t, root, cfg, DefaultReaders(), tools)

	for _, w := range wants {
		var r CapabilityResult
		for _, c := range got.Results {
			if c.Result.Provenance.Provider == w.provider && c.Result.Name == w.name {
				r = c.Result
			}
		}
		if r.Capability != w.capability {
			t.Errorf("%s %q: no result", w.provider, w.name)
			continue
		}
		if w.outcome != "" {
			if r.Outcome != w.outcome || r.Value != nil {
				t.Errorf("%s %q = %s with value %v (%s), want %s without a value", w.provider, w.name, r.Outcome, r.Value, r.Reason, w.outcome)
			}
			continue
		}
		if !r.Outcome.Measured() || r.Value == nil || len(r.Details) == 0 {
			t.Errorf("%s %q = %s value %v details %d (%s), want a measured value with details", w.provider, w.name, r.Outcome, r.Value, len(r.Details), r.Reason)
		}
	}

	snaps, err := NewStore(root).LoadSnapshots()
	if err != nil || len(snaps) != 1 {
		t.Fatalf("LoadSnapshots = %d snapshots, %v", len(snaps), err)
	}
	snap := snaps[0]
	if snap.ID != got.SnapshotID {
		t.Errorf("stored snapshot %s, collected %s", snap.ID, got.SnapshotID)
	}
	if err := snap.Validate(); err != nil {
		t.Errorf("stored snapshot is invalid: %v", err)
	}
	for _, w := range wants {
		if w.outcome == "" {
			continue
		}
		for _, s := range snap.Summary.Capabilities {
			if s.Provider == w.provider && s.Name == w.name && s.Classification != ClassificationUnknown {
				t.Errorf("%s %q failed to measure but was classified %s", w.provider, w.name, s.Classification)
			}
		}
	}
}
