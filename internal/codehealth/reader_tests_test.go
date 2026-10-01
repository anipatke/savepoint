package codehealth

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func readerFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "readers", "tests", name))
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func detailMap(rd Reading) map[string]float64 {
	m := map[string]float64{}
	for _, d := range rd.Details {
		m[d.Key] = d.Number
	}
	return m
}

func goRoot(t *testing.T) string {
	root := t.TempDir()
	writeIn(t, root, "go.mod", "module example.com/m\n")
	return root
}

func TestGoTestReader(t *testing.T) {
	tests := []struct {
		name, file string
		failed     float64
		details    map[string]float64
		partial    bool
		reason     string
		evidence   []EvidenceRef
	}{
		{name: "passing", file: "go-pass.jsonl", details: map[string]float64{"total_tests": 1, "passed_tests": 1, "skipped_tests": 0, "build_failures": 0}},
		{name: "failing with a skip", file: "go-fail.jsonl", failed: 1,
			details:  map[string]float64{"total_tests": 3, "passed_tests": 1, "skipped_tests": 1, "build_failures": 0},
			evidence: []EvidenceRef{{Path: "a", Note: "test TestBad failed"}}},
		{name: "build failure counts once", file: "go-build-failure.jsonl", failed: 1,
			details:  map[string]float64{"total_tests": 1, "passed_tests": 1, "skipped_tests": 0, "build_failures": 1},
			evidence: []EvidenceRef{{Path: "broken", Note: "package example.com/m/broken failed to build"}}},
		{name: "empty suite", file: "go-empty.jsonl", details: map[string]float64{"total_tests": 0}, reason: "no tests ran"},
		{name: "truncated stream", file: "go-truncated.jsonl", partial: true, reason: "ended before 1 started",
			details: map[string]float64{"total_tests": 2, "passed_tests": 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rd, err := GoTestReader{}.Read(context.Background(), ReportInput{Provider: ProviderGoTestJSON, Root: goRoot(t), Data: readerFixture(t, tt.file)})
			if err != nil {
				t.Fatal(err)
			}
			if rd.Value == nil || rd.Value.Number != tt.failed || rd.Value.Unit != UnitCount {
				t.Errorf("value = %+v, want %v count", rd.Value, tt.failed)
			}
			got := detailMap(rd)
			for k, v := range tt.details {
				if got[k] != v {
					t.Errorf("detail %s = %v, want %v", k, got[k], v)
				}
			}
			if rd.Partial != tt.partial || !strings.Contains(rd.Reason, tt.reason) {
				t.Errorf("partial = %v reason = %q, want %v %q", rd.Partial, rd.Reason, tt.partial, tt.reason)
			}
			if tt.evidence != nil && !equalRefs(rd.Evidence, tt.evidence) {
				t.Errorf("evidence = %v, want %v", rd.Evidence, tt.evidence)
			}
			if rd.Provenance.ProviderVersion != "unknown" {
				t.Errorf("version = %q", rd.Provenance.ProviderVersion)
			}
		})
	}
}

func equalRefs(a, b []EvidenceRef) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestGoTestReaderMonorepoResolvesEachModule(t *testing.T) {
	root := t.TempDir()
	writeIn(t, root, "one/go.mod", "module example.com/one\n")
	writeIn(t, root, "two/go.mod", "module example.com/two\n")
	rd, err := GoTestReader{}.Read(context.Background(), ReportInput{Root: root, Data: readerFixture(t, "go-monorepo.jsonl")})
	if err != nil {
		t.Fatal(err)
	}
	want := []EvidenceRef{{Path: "one/a", Note: "test TestA failed"}, {Path: "two/svc/b", Note: "test TestB failed"}}
	if rd.Value.Number != 2 || !equalRefs(rd.Evidence, want) {
		t.Errorf("value %v evidence %v, want 2 %v", rd.Value.Number, rd.Evidence, want)
	}
}

func TestGoTestReaderKeepsFailureCountWhenPackageUnresolved(t *testing.T) {
	rd, err := GoTestReader{}.Read(context.Background(), ReportInput{Root: t.TempDir(), Data: readerFixture(t, "go-fail.jsonl")})
	if err != nil {
		t.Fatal(err)
	}
	if rd.Value.Number != 1 || len(rd.Evidence) != 0 {
		t.Errorf("value %v evidence %v", rd.Value.Number, rd.Evidence)
	}
}

func TestGoTestReaderRejectsUnusableReports(t *testing.T) {
	for name, data := range map[string][]byte{
		"malformed": readerFixture(t, "go-malformed.jsonl"),
		"empty":     nil,
		"blank":     []byte("\n\n"),
	} {
		if _, err := (GoTestReader{}).Read(context.Background(), ReportInput{Root: goRoot(t), Data: data}); err == nil {
			t.Errorf("%s: want an error", name)
		}
	}
}

func TestJUnitReader(t *testing.T) {
	tests := []struct {
		name, file string
		provider   ProviderKey
		failed     float64
		details    map[string]float64
		reason     string
		evidence   []EvidenceRef
	}{
		{name: "vitest failing", file: "vitest-fail.xml", provider: ProviderVitestJUnit, failed: 1,
			details:  map[string]float64{"total_tests": 3, "skipped_tests": 1, "errors": 0},
			evidence: []EvidenceRef{{Path: "src/a.test.ts", Note: "test subtracts failed"}}},
		{name: "pytest failures and errors", file: "pytest-fail.xml", provider: ProviderPytestJUnit, failed: 2,
			details: map[string]float64{"total_tests": 4, "skipped_tests": 1, "errors": 1},
			evidence: []EvidenceRef{
				{Path: "tests/test_a.py", Note: "test test_bad failed"},
				{Path: "tests/test_b.py", Line: 12, Note: "test test_err errored"},
			}},
		{name: "passing", file: "junit-pass.xml", provider: ProviderPytestJUnit, details: map[string]float64{"total_tests": 1}},
		{name: "empty suite", file: "junit-empty.xml", provider: ProviderVitestJUnit, details: map[string]float64{"total_tests": 0}, reason: "no tests ran"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rd, err := JUnitReader{}.Read(context.Background(), ReportInput{Provider: tt.provider, Root: t.TempDir(), Data: readerFixture(t, tt.file)})
			if err != nil {
				t.Fatal(err)
			}
			if rd.Value.Number != tt.failed || rd.Partial || !strings.Contains(rd.Reason, tt.reason) {
				t.Errorf("value %v partial %v reason %q", rd.Value.Number, rd.Partial, rd.Reason)
			}
			got := detailMap(rd)
			for k, v := range tt.details {
				if got[k] != v {
					t.Errorf("detail %s = %v, want %v", k, got[k], v)
				}
			}
			if tt.evidence != nil && !equalRefs(rd.Evidence, tt.evidence) {
				t.Errorf("evidence = %v, want %v", rd.Evidence, tt.evidence)
			}
			for _, e := range rd.Evidence {
				if strings.Contains(e.Note, "SECRET") || strings.Contains(e.Note, "ci-box") {
					t.Errorf("evidence leaks report content: %v", e)
				}
			}
		})
	}
}

func TestJUnitReaderRejectsUnusableReports(t *testing.T) {
	for name, data := range map[string][]byte{
		"malformed": readerFixture(t, "junit-malformed.xml"),
		"empty":     nil,
		"not junit": []byte("<html></html>"),
	} {
		if _, err := (JUnitReader{}).Read(context.Background(), ReportInput{Root: "/r", Data: data}); err == nil {
			t.Errorf("%s: want an error", name)
		}
	}
}

func TestJUnitReaderHonoursScope(t *testing.T) {
	rd, err := JUnitReader{}.Read(context.Background(), ReportInput{Provider: ProviderVitestJUnit, Root: "/r", Scope: []string{"lib/**"}, Data: readerFixture(t, "vitest-fail.xml")})
	if err != nil {
		t.Fatal(err)
	}
	if rd.Value.Number != 1 || len(rd.Evidence) != 0 {
		t.Errorf("value %v evidence %v", rd.Value.Number, rd.Evidence)
	}
}
