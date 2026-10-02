package codehealth

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func vulnerabilityFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "readers", "vulnerabilities", name))
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestOSVScannerReaderValue(t *testing.T) {
	tests := []struct {
		name     string
		fixture  string
		scope    []string
		total    float64
		severity map[string]float64
		version  string
		partial  bool
		reason   string
		evidence []EvidenceRef
	}{
		{name: "go, npm, and python lockfiles", fixture: "multi-ecosystem.json", total: 3, version: "unknown",
			severity: map[string]float64{"critical": 1, "high": 1, "medium": 1, "low": 0, "unknown": 0},
			reason:   "snapshot is unknown",
			evidence: []EvidenceRef{
				{Path: "web/package-lock.json", Note: "lodash 4.17.0: GHSA-a"},
				{Path: "go.mod", Note: "golang.org/x/net 0.1.0: GO-2023-1, CVE-2023-1"},
				{Path: "app/poetry.lock", Note: "flask 2.0.0: PYSEC-1"},
			}},
		{name: "scope keeps one lockfile", fixture: "multi-ecosystem.json", scope: []string{"web/**"}, total: 1, version: "unknown",
			severity: map[string]float64{"critical": 1},
			evidence: []EvidenceRef{{Path: "web/package-lock.json", Note: "lodash 4.17.0: GHSA-a"}}},
		{name: "no findings", fixture: "none.json", total: 0, version: "unknown",
			severity: map[string]float64{"critical": 0, "high": 0, "medium": 0, "low": 0, "unknown": 0},
			reason:   "queried at collection time"},
		{name: "every severity bucket and its boundaries", fixture: "all-severities.json", total: 10, version: "2.2.0",
			severity: map[string]float64{"critical": 2, "high": 2, "medium": 2, "low": 2, "unknown": 2},
			evidence: nil},
		{name: "missing max_severity is unknown", fixture: "missing-severity.json", total: 1, version: "unknown",
			severity: map[string]float64{"unknown": 1, "critical": 0, "high": 0, "medium": 0, "low": 0},
			evidence: []EvidenceRef{{Path: "go.mod", Note: "a 1: A-1, A-2, A-3 +2 more"}}},
		{name: "manifest without resolved versions is partial", fixture: "unresolved.json", total: 1, version: "unknown",
			partial: true, reason: "1 package(s) scanned without a resolved version",
			severity: map[string]float64{"medium": 1},
			evidence: []EvidenceRef{{Path: "requirements.txt", Note: "requests unknown: PYSEC-9"}}},
		{name: "scanner error is partial", fixture: "scanner-error.json", total: 0, version: "unknown",
			partial: true, reason: "1 scanner error(s) in the report"},
		{name: "vulnerabilities without groups are partial", fixture: "ungrouped.json", total: 0, version: "unknown",
			partial: true, reason: "1 package(s) with vulnerabilities but no groups"},
		{name: "absolute paths; source outside the root is partial", fixture: "absolute.json", total: 2, version: "unknown",
			partial: true, reason: "1 source(s) outside the project have no affected items",
			severity: map[string]float64{"high": 1, "medium": 1},
			evidence: []EvidenceRef{{Path: "go.mod", Note: "a 1: A-1"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := ReportInput{Provider: ProviderOSVScannerJSON, Root: vitestRoot, Data: vulnerabilityFixture(t, tt.fixture), Scope: tt.scope}
			rd, err := OSVScannerReader{}.Read(context.Background(), in)
			if err != nil {
				t.Fatal(err)
			}
			if rd.Value == nil || rd.Value.Number != tt.total || rd.Value.Unit != UnitCount {
				t.Errorf("value = %+v, want %v count", rd.Value, tt.total)
			}
			if rd.Partial != tt.partial || !strings.Contains(rd.Reason, tt.reason) {
				t.Errorf("partial = %v reason = %q, want %v %q", rd.Partial, rd.Reason, tt.partial, tt.reason)
			}
			if rd.Provenance.ProviderVersion != tt.version {
				t.Errorf("version = %q, want %q", rd.Provenance.ProviderVersion, tt.version)
			}
			got := detailMap(rd)
			var sum float64
			for _, key := range vulnerabilitySeverityKeys {
				if _, ok := got[key]; !ok {
					t.Errorf("detail %s is missing", key)
				}
				sum += got[key]
			}
			if sum != tt.total {
				t.Errorf("severity counts sum to %v, want the total %v", sum, tt.total)
			}
			for k, v := range tt.severity {
				if got[k] != v {
					t.Errorf("detail %s = %v, want %v", k, got[k], v)
				}
			}
			if tt.fixture != "all-severities.json" && !equalRefs(rd.Evidence, tt.evidence) {
				t.Errorf("evidence = %v, want %v", rd.Evidence, tt.evidence)
			}
			res := CapabilityResult{Capability: CapabilityDependencyVulnerability, Value: rd.Value, Details: rd.Details, Evidence: rd.Evidence, Reason: rd.Reason}
			if err := res.validateBounded("result"); err != nil {
				t.Errorf("reading is not storable: %v", err)
			}
		})
	}
}

func TestOSVScannerEvidenceIsWorstFirstAndBounded(t *testing.T) {
	rd, err := OSVScannerReader{}.Read(context.Background(), ReportInput{Root: vitestRoot, Data: vulnerabilityFixture(t, "all-severities.json")})
	if err != nil {
		t.Fatal(err)
	}
	var notes []string
	for _, e := range rd.Evidence {
		notes = append(notes, strings.Fields(e.Note)[0])
	}
	// a 10.0, b 9.0, c 8.9, then d 7.0 and the unknown j, i which rank with high.
	if got, want := strings.Join(notes, ""), "abcdijefgh"; got != want {
		t.Errorf("evidence order = %q, want %q", got, want)
	}

	var b strings.Builder
	b.WriteString(`{"results":[{"source":{"path":"go.mod"},"packages":[`)
	for i := 0; i < MaxEvidence+10; i++ {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString(`{"package":{"name":"p","version":"1"},"groups":[{"ids":["X"],"max_severity":"5"}]}`)
	}
	b.WriteString("]}]}")
	rd, err = OSVScannerReader{}.Read(context.Background(), ReportInput{Root: vitestRoot, Data: []byte(b.String())})
	if err != nil {
		t.Fatal(err)
	}
	if len(rd.Evidence) != MaxEvidence || rd.Value.Number != MaxEvidence+10 {
		t.Errorf("evidence has %d items for value %v, want %d items for %d", len(rd.Evidence), rd.Value.Number, MaxEvidence, MaxEvidence+10)
	}
}

func TestOSVScannerUnknownSeverityBlocks(t *testing.T) {
	rd, err := OSVScannerReader{}.Read(context.Background(), ReportInput{Root: vitestRoot, Data: vulnerabilityFixture(t, "missing-severity.json")})
	if err != nil {
		t.Fatal(err)
	}
	res := CapabilityResult{
		Capability: CapabilityDependencyVulnerability, Outcome: OutcomeAvailable, Freshness: FreshnessFresh,
		Value: rd.Value, Details: rd.Details,
	}
	if got := Assess(OriginOfficial, res, nil, nil).Classification; got != ClassificationNeedsAttention {
		t.Errorf("classification = %q, want needs_attention", got)
	}
}

func TestOSVScannerReaderRejectsUnusableReports(t *testing.T) {
	for _, tt := range []struct{ fixture, want string }{
		{"malformed.json", "not an OSV-Scanner report"},
		{"no-results.json", "no results list"},
	} {
		t.Run(tt.fixture, func(t *testing.T) {
			rd, err := OSVScannerReader{}.Read(context.Background(), ReportInput{Root: vitestRoot, Data: vulnerabilityFixture(t, tt.fixture)})
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("err = %v, want one containing %q", err, tt.want)
			}
			if rd.Value != nil {
				t.Errorf("an unusable report produced value %+v", rd.Value)
			}
		})
	}
}

func TestOSVScannerReaderStopsWhenCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := (OSVScannerReader{}).Read(ctx, ReportInput{Root: vitestRoot, Data: vulnerabilityFixture(t, "multi-ecosystem.json")}); err == nil {
		t.Error("a cancelled read should fail")
	}
}
