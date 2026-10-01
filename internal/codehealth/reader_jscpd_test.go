package codehealth

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func duplicationFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "readers", "duplication", name))
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestJscpdReaderValue(t *testing.T) {
	tests := []struct {
		name     string
		fixture  string
		scope    []string
		value    float64
		version  string
		partial  bool
		reason   string
		details  map[string]float64
		evidence []EvidenceRef
	}{
		{name: "go, typescript, and python clones", fixture: "mixed.json", value: 10, version: "unknown",
			details: map[string]float64{"duplicated_lines": 100, "total_lines": 1000, "clones": 3, "sources": 12},
			evidence: []EvidenceRef{
				{Path: "web/src/x.ts", Line: 5, Note: "50 lines repeated at web/src/y.ts:7"},
				{Path: "app/p.py", Line: 1, Note: "30 lines repeated at app/q.py:3"},
				{Path: "pkg/a.go", Line: 10, Note: "20 lines repeated at pkg/b.go:40"},
			}},
		{name: "scope keeps one language's clones as items", fixture: "mixed.json", scope: []string{"app/**"}, value: 10, version: "unknown",
			evidence: []EvidenceRef{{Path: "app/p.py", Line: 1, Note: "30 lines repeated at app/q.py:3"}}},
		{name: "no clones", fixture: "none.json", value: 0, version: "unknown",
			details: map[string]float64{"duplicated_lines": 0, "total_lines": 500, "clones": 0, "sources": 4}},
		{name: "absolute paths; clone outside the root is partial", fixture: "absolute.json", value: 20, version: "unknown",
			partial: true, reason: "1 clone(s) in files outside the project could not be mapped",
			evidence: []EvidenceRef{{Path: "src/a.go", Line: 3, Note: "25 lines repeated at src/b.go:9"}}},
		{name: "version is recorded when stated", fixture: "versioned.json", value: 0, version: "4.0.5"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := ReportInput{Provider: ProviderJscpdJSON, Root: vitestRoot, Data: duplicationFixture(t, tt.fixture), Scope: tt.scope}
			rd, err := JscpdReader{}.Read(context.Background(), in)
			if err != nil {
				t.Fatal(err)
			}
			if rd.Value == nil || rd.Value.Number != tt.value || rd.Value.Unit != UnitPercent {
				t.Errorf("value = %+v, want %v percent", rd.Value, tt.value)
			}
			if rd.Partial != tt.partial || !strings.Contains(rd.Reason, tt.reason) {
				t.Errorf("partial = %v reason = %q, want %v %q", rd.Partial, rd.Reason, tt.partial, tt.reason)
			}
			if rd.Provenance.ProviderVersion != tt.version {
				t.Errorf("version = %q, want %q", rd.Provenance.ProviderVersion, tt.version)
			}
			got := detailMap(rd)
			for k, v := range tt.details {
				if got[k] != v {
					t.Errorf("detail %s = %v, want %v", k, got[k], v)
				}
			}
			if tt.evidence == nil && len(rd.Evidence) != 0 || tt.evidence != nil && !equalRefs(rd.Evidence, tt.evidence) {
				t.Errorf("evidence = %v, want %v", rd.Evidence, tt.evidence)
			}
			res := CapabilityResult{Capability: CapabilityDuplication, Value: rd.Value, Details: rd.Details, Evidence: rd.Evidence}
			if err := res.validateBounded("result"); err != nil {
				t.Errorf("reading is not storable: %v", err)
			}
		})
	}
}

func TestJscpdReaderRejectsUnusableReports(t *testing.T) {
	for _, tt := range []struct{ fixture, want string }{
		{"malformed.json", "not a jscpd report"},
		{"no-statistics.json", "no statistics total"},
		{"zero-lines.json", "scanned no lines"},
		{"inconsistent.json", "does not match"},
	} {
		t.Run(tt.fixture, func(t *testing.T) {
			rd, err := JscpdReader{}.Read(context.Background(), ReportInput{Root: vitestRoot, Data: duplicationFixture(t, tt.fixture)})
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("err = %v, want one containing %q", err, tt.want)
			}
			if rd.Value != nil {
				t.Errorf("an unusable report produced value %+v", rd.Value)
			}
		})
	}
}

func TestJscpdReaderStopsWhenCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := (JscpdReader{}).Read(ctx, ReportInput{Root: vitestRoot, Data: duplicationFixture(t, "mixed.json")}); err == nil {
		t.Error("a cancelled read should fail")
	}
}

func TestJscpdEvidenceIsBounded(t *testing.T) {
	var b strings.Builder
	b.WriteString(`{"statistics":{"total":{"lines":1000,"sources":2,"clones":30,"duplicatedLines":300,"percentage":30}},"duplicates":[`)
	for i := 0; i < MaxEvidence+10; i++ {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString(`{"lines":10,"firstFile":{"name":"a.go","start":1},"secondFile":{"name":"b.go","start":1}}`)
	}
	b.WriteString("]}")
	rd, err := JscpdReader{}.Read(context.Background(), ReportInput{Root: vitestRoot, Data: []byte(b.String())})
	if err != nil {
		t.Fatal(err)
	}
	if len(rd.Evidence) != MaxEvidence {
		t.Errorf("evidence has %d items, want %d", len(rd.Evidence), MaxEvidence)
	}
}
