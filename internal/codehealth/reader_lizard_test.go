package codehealth

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func complexityFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "readers", "complexity", name))
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestLizardReaderValue(t *testing.T) {
	tests := []struct {
		name     string
		fixture  string
		scope    []string
		value    float64
		partial  bool
		reason   string
		details  map[string]float64
		evidence []EvidenceRef
	}{
		{name: "go, python, and typescript in one report", fixture: "mixed.csv", value: 22,
			details: map[string]float64{"functions": 5, "average_ccn": 8, "functions_over_10": 2, "functions_over_20": 1},
			evidence: []EvidenceRef{
				{Path: "web/src/view.ts", Line: 20, Note: "render has complexity 22"},
				{Path: "pkg/parse.go", Line: 10, Note: "Parse has complexity 12"},
			}},
		{name: "scope keeps one language", fixture: "mixed.csv", scope: []string{"app/**"}, value: 3,
			details: map[string]float64{"functions": 1, "average_ccn": 3, "functions_over_10": 0, "functions_over_20": 0}},
		{name: "empty report has no functions", fixture: "empty.csv", value: 0, reason: "no functions",
			details: map[string]float64{"functions": 0, "average_ccn": 0, "functions_over_10": 0, "functions_over_20": 0}},
		{name: "function outside the root is partial", fixture: "outside.csv", value: 4, partial: true,
			reason:  "1 function(s) in files outside the project could not be mapped",
			details: map[string]float64{"functions": 1, "average_ccn": 4, "functions_over_10": 0, "functions_over_20": 0}},
		{name: "everything outside the root measures nothing", fixture: "all-outside.csv", partial: true,
			reason: "1 function(s) in files outside the project could not be mapped"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := ReportInput{Provider: ProviderLizardCSV, Root: vitestRoot, Data: complexityFixture(t, tt.fixture), Scope: tt.scope}
			rd, err := LizardReader{}.Read(context.Background(), in)
			if err != nil {
				t.Fatal(err)
			}
			switch {
			case tt.partial && tt.value == 0:
				if rd.Value != nil {
					t.Errorf("value = %+v, want none: nothing was measured", rd.Value)
				}
			case rd.Value == nil || rd.Value.Number != tt.value || rd.Value.Unit != UnitCCN:
				t.Errorf("value = %+v, want %v ccn", rd.Value, tt.value)
			}
			if rd.Partial != tt.partial || !strings.Contains(rd.Reason, tt.reason) {
				t.Errorf("partial = %v reason = %q, want %v %q", rd.Partial, rd.Reason, tt.partial, tt.reason)
			}
			if !strings.Contains(rd.Reason, "version") && !tt.partial {
				t.Errorf("reason %q should say Lizard states no version", rd.Reason)
			}
			if rd.Provenance.ProviderVersion != "unknown" {
				t.Errorf("version = %q, want unknown", rd.Provenance.ProviderVersion)
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
			res := CapabilityResult{Capability: CapabilityComplexity, Value: rd.Value, Details: rd.Details, Evidence: rd.Evidence}
			if err := res.validateBounded("result"); err != nil {
				t.Errorf("reading is not storable: %v", err)
			}
		})
	}
}

func TestLizardReaderRejectsMalformedReports(t *testing.T) {
	for _, tt := range []struct{ fixture, want string }{
		{"wrong-columns.csv", "wrong number of fields"},
		{"bad-ccn.csv", `CCN "many"`},
		{"bad-line.csv", `start line "five"`},
	} {
		t.Run(tt.fixture, func(t *testing.T) {
			rd, err := LizardReader{}.Read(context.Background(), ReportInput{Root: vitestRoot, Data: complexityFixture(t, tt.fixture)})
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("err = %v, want one containing %q", err, tt.want)
			}
			if rd.Value != nil {
				t.Errorf("a malformed report produced value %+v", rd.Value)
			}
		})
	}
}

func TestLizardReaderStopsWhenCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := (LizardReader{}).Read(ctx, ReportInput{Root: vitestRoot, Data: complexityFixture(t, "mixed.csv")}); err == nil {
		t.Error("a cancelled read should fail")
	}
}

func TestLizardEvidenceIsBounded(t *testing.T) {
	var b strings.Builder
	for i := 0; i < MaxEvidence+10; i++ {
		b.WriteString("9,11,50,0,9,\"f@1-9@a.go\",\"a.go\",\"f\",\"f()\",1,9\n")
	}
	rd, err := LizardReader{}.Read(context.Background(), ReportInput{Root: vitestRoot, Data: []byte(b.String())})
	if err != nil {
		t.Fatal(err)
	}
	if len(rd.Evidence) != MaxEvidence {
		t.Errorf("evidence has %d items, want %d", len(rd.Evidence), MaxEvidence)
	}
}
