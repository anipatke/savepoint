package codehealth

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func coverageFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "readers", "coverage", name))
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func twoPackageRoot(t *testing.T) string {
	root := t.TempDir()
	writeIn(t, root, "go.mod", "module example.com/m\n")
	return root
}

const vitestRoot = "/work/proj"

type coverageCase struct {
	name     string
	reader   Reader
	input    func(t *testing.T) ReportInput
	percent  float64
	version  string
	details  map[string]float64
	partial  bool
	reason   string
	evidence []EvidenceRef
}

func TestCoverageReadersValue(t *testing.T) {
	goIn := func(file string) func(*testing.T) ReportInput {
		return func(t *testing.T) ReportInput {
			return ReportInput{Provider: ProviderGoCoverProfile, Root: twoPackageRoot(t), Data: coverageFixture(t, file)}
		}
	}
	vitestIn := func(file string) func(*testing.T) ReportInput {
		return func(t *testing.T) ReportInput {
			return ReportInput{Provider: ProviderVitestV8, Root: vitestRoot, Data: coverageFixture(t, file)}
		}
	}
	pyIn := func(file string) func(*testing.T) ReportInput {
		return func(t *testing.T) ReportInput {
			return ReportInput{Provider: ProviderCoveragePyJSON, Root: "/r", Data: coverageFixture(t, file)}
		}
	}
	tests := []coverageCase{
		{name: "go populated", reader: GoCoverReader{}, input: goIn("go-populated.out"), percent: 70, version: "unknown",
			details:  map[string]float64{"covered_statements": 7, "total_statements": 10},
			evidence: []EvidenceRef{{Path: "a/a.go", Note: "40.0% of 5 statements covered"}}},
		{name: "go merges a block seen by several runs", reader: GoCoverReader{}, input: goIn("go-merged.out"), percent: 100,
			details: map[string]float64{"covered_statements": 5, "total_statements": 5}},
		{name: "go fully covered", reader: GoCoverReader{}, input: goIn("go-full.out"), percent: 100,
			details: map[string]float64{"covered_statements": 7, "total_statements": 7}},
		{name: "go unresolvable import path", reader: GoCoverReader{}, input: goIn("go-unresolvable.out"), percent: 100,
			partial: true, reason: "1 profile file(s) could not be mapped",
			details: map[string]float64{"covered_statements": 2, "total_statements": 2}},
		{name: "vitest populated", reader: VitestCoverageReader{}, input: vitestIn("vitest-populated.json"), percent: 400.0 / 6, version: "unknown",
			details: map[string]float64{"covered_statements": 4, "total_statements": 6, "covered_functions": 2, "total_functions": 3,
				"covered_branches": 1, "total_branches": 4},
			evidence: []EvidenceRef{{Path: "src/a.ts", Note: "50.0% of 4 statements covered"}}},
		{name: "vitest fully covered", reader: VitestCoverageReader{}, input: vitestIn("vitest-full.json"), percent: 100,
			details: map[string]float64{"covered_statements": 2, "total_statements": 2, "covered_branches": 2, "total_branches": 2}},
		{name: "vitest file outside the repository", reader: VitestCoverageReader{}, input: vitestIn("vitest-outside.json"), percent: 50,
			partial: true, reason: "1 report file(s) could not be mapped",
			evidence: []EvidenceRef{{Path: "src/a.ts", Note: "50.0% of 2 statements covered"}}},
		{name: "coverage.py populated", reader: CoveragePyReader{}, input: pyIn("coveragepy-populated.json"), percent: 70, version: "7.6.1",
			details:  map[string]float64{"covered_statements": 7, "total_statements": 10},
			evidence: []EvidenceRef{{Path: "app/a.py", Note: "50.0% of 6 statements covered"}}},
		{name: "coverage.py with branches", reader: CoveragePyReader{}, input: pyIn("coveragepy-branches.json"), percent: 50, version: "7.6.1",
			details:  map[string]float64{"covered_statements": 3, "total_statements": 6, "covered_branches": 1, "total_branches": 4},
			evidence: []EvidenceRef{{Path: "app/a.py", Note: "50.0% of 6 statements covered"}}},
		{name: "coverage.py fully covered", reader: CoveragePyReader{}, input: pyIn("coveragepy-full.json"), percent: 100, version: "7.6.1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rd, err := tt.reader.Read(context.Background(), tt.input(t))
			if err != nil {
				t.Fatal(err)
			}
			if rd.Value == nil || rd.Value.Number != tt.percent || rd.Value.Unit != UnitPercent {
				t.Errorf("value = %+v, want %v percent", rd.Value, tt.percent)
			}
			got := detailMap(rd)
			for k, v := range tt.details {
				if got[k] != v {
					t.Errorf("detail %s = %v, want %v", k, got[k], v)
				}
			}
			if _, ok := got["covered_statements"]; !ok {
				t.Error("missing covered_statements")
			}
			if tt.version != "" && rd.Provenance.ProviderVersion != tt.version {
				t.Errorf("version = %q, want %q", rd.Provenance.ProviderVersion, tt.version)
			}
			if rd.Partial != tt.partial || !strings.Contains(rd.Reason, tt.reason) {
				t.Errorf("partial = %v reason = %q, want %v %q", rd.Partial, rd.Reason, tt.partial, tt.reason)
			}
			if tt.evidence == nil && len(rd.Evidence) != 0 || tt.evidence != nil && !equalRefs(rd.Evidence, tt.evidence) {
				t.Errorf("evidence = %v, want %v", rd.Evidence, tt.evidence)
			}
			res := CapabilityResult{Capability: CapabilityCoverage, Value: rd.Value, Details: rd.Details, Evidence: rd.Evidence}
			if err := res.validateBounded("result"); err != nil {
				t.Errorf("reading is not storable: %v", err)
			}
		})
	}
}

func TestCoverageOmitsBranchDetailsWhenReportHasNone(t *testing.T) {
	rd, err := CoveragePyReader{}.Read(context.Background(), ReportInput{Root: "/r", Data: coverageFixture(t, "coveragepy-populated.json")})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := detailMap(rd)["total_branches"]; ok {
		t.Error("branch counts reported although branch coverage was off")
	}
}

func TestGoCoverReaderMonorepoResolvesEachModule(t *testing.T) {
	root := t.TempDir()
	writeIn(t, root, "one/go.mod", "module example.com/one\n")
	writeIn(t, root, "two/go.mod", "module example.com/two\n")
	rd, err := GoCoverReader{}.Read(context.Background(), ReportInput{Root: root, Data: coverageFixture(t, "go-monorepo.out")})
	if err != nil {
		t.Fatal(err)
	}
	want := []EvidenceRef{{Path: "two/svc/b.go", Note: "0.0% of 2 statements covered"}}
	if rd.Value.Number != 50 || rd.Partial || !equalRefs(rd.Evidence, want) {
		t.Errorf("value %v partial %v evidence %v, want 50 false %v", rd.Value.Number, rd.Partial, rd.Evidence, want)
	}
}

func TestGoCoverReaderDropsFilesOutsideScopeWithoutGoingPartial(t *testing.T) {
	rd, err := GoCoverReader{}.Read(context.Background(), ReportInput{
		Root: twoPackageRoot(t), Data: coverageFixture(t, "go-populated.out"), Scope: []string{"a"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if rd.Value.Number != 40 || rd.Partial {
		t.Errorf("value %v partial %v, want 40 false", rd.Value.Number, rd.Partial)
	}
}

func TestCoverageReadersRejectUnusableReports(t *testing.T) {
	tests := []struct {
		name   string
		reader Reader
		data   []byte
	}{
		{"go zero statements", GoCoverReader{}, coverageFixture(t, "go-zero.out")},
		{"go missing mode", GoCoverReader{}, coverageFixture(t, "go-malformed.out")},
		{"go empty", GoCoverReader{}, nil},
		{"go bad block", GoCoverReader{}, []byte("mode: set\nnot a block\n")},
		{"vitest zero statements", VitestCoverageReader{}, coverageFixture(t, "vitest-zero.json")},
		{"vitest malformed", VitestCoverageReader{}, coverageFixture(t, "vitest-malformed.json")},
		{"vitest empty report", VitestCoverageReader{}, []byte("{}")},
		{"coverage.py zero statements", CoveragePyReader{}, coverageFixture(t, "coveragepy-zero.json")},
		{"coverage.py no totals", CoveragePyReader{}, coverageFixture(t, "coveragepy-no-totals.json")},
		{"coverage.py malformed", CoveragePyReader{}, coverageFixture(t, "coveragepy-malformed.json")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rd, err := tt.reader.Read(context.Background(), ReportInput{Root: twoPackageRoot(t), Data: tt.data})
			if err == nil {
				t.Fatalf("want an error, got %+v", rd)
			}
		})
	}
}

func TestCoverageReadersStopWhenCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := (GoCoverReader{}).Read(ctx, ReportInput{Root: twoPackageRoot(t), Data: coverageFixture(t, "go-populated.out")}); err == nil {
		t.Error("go: want a cancellation error")
	}
	if _, err := (VitestCoverageReader{}).Read(ctx, ReportInput{Root: vitestRoot, Data: coverageFixture(t, "vitest-populated.json")}); err == nil {
		t.Error("vitest: want a cancellation error")
	}
}
