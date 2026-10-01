package codehealth

import (
	"context"
	"strings"
	"testing"
)

// Regression tests for the O-029 Full Check issues: instance scope (I-097),
// unusable test reports (I-098), invalid coverage measurements (I-099), and
// repeated vulnerability groups (I-100).

func goLines(lines ...string) []byte { return []byte(strings.Join(lines, "\n") + "\n") }

func scopedGoRoot(t *testing.T) string {
	root := t.TempDir()
	writeIn(t, root, "go.mod", "module example.com/m\n")
	return root
}

func TestGoTestReaderCountsOnlyTheInstanceScope(t *testing.T) {
	data := goLines(
		`{"Action":"run","Package":"example.com/m/a","Test":"TestA"}`,
		`{"Action":"pass","Package":"example.com/m/a","Test":"TestA"}`,
		`{"Action":"pass","Package":"example.com/m/a"}`,
		`{"Action":"run","Package":"example.com/m/b","Test":"TestB"}`,
		`{"Action":"fail","Package":"example.com/m/b","Test":"TestB"}`,
		`{"Action":"fail","Package":"example.com/m/b"}`,
		`{"Action":"fail","Package":"example.com/other/c"}`,
	)
	root := scopedGoRoot(t)
	for _, tt := range []struct {
		name            string
		scope, exclude  []string
		failed, total   float64
		wantPartialNone bool
	}{
		{name: "include a", scope: []string{"a/**"}, failed: 1, total: 1}, // example.com/other/c is unowned, so still counted
		{name: "exclude b", exclude: []string{"b/**"}, failed: 1, total: 1},
		{name: "include b", scope: []string{"b/**"}, failed: 2, total: 1},
		{name: "no scope", failed: 2, total: 2},
	} {
		t.Run(tt.name, func(t *testing.T) {
			rd, err := GoTestReader{}.Read(context.Background(), ReportInput{Root: root, Data: data, Scope: tt.scope, Exclusions: tt.exclude})
			if err != nil {
				t.Fatal(err)
			}
			if rd.Value.Number != tt.failed || detailMap(rd)["total_tests"] != tt.total {
				t.Errorf("failed %v total %v, want %v %v", rd.Value.Number, detailMap(rd)["total_tests"], tt.failed, tt.total)
			}
		})
	}
}

func TestGoTestReaderRejectsReportsWithoutEvents(t *testing.T) {
	for name, data := range map[string][]byte{
		"null":         []byte("null\n"),
		"wrong schema": []byte(`{"results":[]}` + "\n"),
		"array":        []byte("[]\n"),
	} {
		if _, err := (GoTestReader{}).Read(context.Background(), ReportInput{Root: scopedGoRoot(t), Data: data}); err == nil {
			t.Errorf("%s: want an error, not a measurement", name)
		}
	}
}

func TestGoTestReaderIsPartialWithoutAFinalPackageResult(t *testing.T) {
	for name, data := range map[string][]byte{
		"finished tests, no package result": goLines(
			`{"Action":"run","Package":"example.com/m/a","Test":"TestA"}`,
			`{"Action":"pass","Package":"example.com/m/a","Test":"TestA"}`),
		"start only": goLines(`{"Action":"start","Package":"example.com/m/a"}`),
	} {
		t.Run(name, func(t *testing.T) {
			rd, err := GoTestReader{}.Read(context.Background(), ReportInput{Root: scopedGoRoot(t), Data: data})
			if err != nil {
				t.Fatal(err)
			}
			if !rd.Partial || !strings.Contains(rd.Reason, "final result") {
				t.Errorf("partial = %v reason = %q, want a partial reading naming the missing result", rd.Partial, rd.Reason)
			}
		})
	}
}

func TestGoTestReaderIgnoresAnExcludedUnfinishedPackage(t *testing.T) {
	data := goLines(
		`{"Action":"run","Package":"example.com/m/a","Test":"TestA"}`,
		`{"Action":"pass","Package":"example.com/m/a","Test":"TestA"}`,
		`{"Action":"pass","Package":"example.com/m/a"}`,
		`{"Action":"start","Package":"example.com/m/b"}`,
	)
	rd, err := GoTestReader{}.Read(context.Background(), ReportInput{Root: scopedGoRoot(t), Data: data, Scope: []string{"a/**"}})
	if err != nil {
		t.Fatal(err)
	}
	if rd.Partial {
		t.Errorf("an excluded package left the reading partial: %q", rd.Reason)
	}
}

func junitSuiteXML(attrs string, cases ...string) []byte {
	return []byte(`<testsuites><testsuite ` + attrs + `>` + strings.Join(cases, "") + `</testsuite></testsuites>`)
}

func TestJUnitReaderCountsOnlyTheInstanceScope(t *testing.T) {
	data := junitSuiteXML(`tests="2" failures="1"`,
		`<testcase name="ok" classname="a/ok.test.ts"/>`,
		`<testcase name="bad" classname="b/bad.test.ts"><failure/></testcase>`)
	for _, tt := range []struct {
		name           string
		scope, exclude []string
		failed, total  float64
		evidence       int
	}{
		{name: "include a", scope: []string{"a/**"}, failed: 0, total: 1},
		{name: "exclude b", exclude: []string{"b/**"}, failed: 0, total: 1},
		{name: "include b", scope: []string{"b/**"}, failed: 1, total: 1, evidence: 1},
		{name: "no scope", failed: 1, total: 2, evidence: 1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			rd, err := JUnitReader{}.Read(context.Background(), ReportInput{Provider: ProviderVitestJUnit, Root: "/r", Data: data, Scope: tt.scope, Exclusions: tt.exclude})
			if err != nil {
				t.Fatal(err)
			}
			if rd.Value.Number != tt.failed || detailMap(rd)["total_tests"] != tt.total || len(rd.Evidence) != tt.evidence || rd.Partial {
				t.Errorf("failed %v total %v evidence %v partial %v, want %v %v %d false", rd.Value.Number, detailMap(rd)["total_tests"], rd.Evidence, rd.Partial, tt.failed, tt.total, tt.evidence)
			}
		})
	}
}

func TestJUnitReaderScopesPytestClassnames(t *testing.T) {
	data := junitSuiteXML(`tests="2"`,
		`<testcase name="test_a" classname="a.test_a"/>`,
		`<testcase name="test_b" classname="b.test_b"><failure/></testcase>`)
	rd, err := JUnitReader{}.Read(context.Background(), ReportInput{Provider: ProviderPytestJUnit, Root: "/r", Data: data, Scope: []string{"a/**"}})
	if err != nil {
		t.Fatal(err)
	}
	if rd.Value.Number != 0 || detailMap(rd)["total_tests"] != 1 {
		t.Errorf("failed %v total %v, want the b failure left out", rd.Value.Number, detailMap(rd)["total_tests"])
	}
}

func TestJUnitReaderKeepsCasesItCannotPlace(t *testing.T) {
	data := junitSuiteXML(`tests="1"`, `<testcase name="anon"><failure/></testcase>`)
	rd, err := JUnitReader{}.Read(context.Background(), ReportInput{Provider: ProviderVitestJUnit, Root: "/r", Data: data, Scope: []string{"a/**"}})
	if err != nil {
		t.Fatal(err)
	}
	if rd.Value.Number != 1 {
		t.Errorf("failed = %v, want a case without a path still counted", rd.Value.Number)
	}
}

func TestJUnitReaderDiagnosesUnusableStructure(t *testing.T) {
	if _, err := (JUnitReader{}).Read(context.Background(), ReportInput{Root: "/r", Data: []byte(`<testsuite/><testsuite/>`)}); err == nil {
		t.Error("two document roots: want an error")
	}
	if _, err := (JUnitReader{}).Read(context.Background(), ReportInput{Root: "/r", Data: []byte(`<testsuite tests="many"/>`)}); err == nil {
		t.Error("a non-numeric declared count: want an error")
	}
	for name, data := range map[string][]byte{
		"declares an omitted failure": junitSuiteXML(`tests="2" failures="1"`, `<testcase name="ok"/>`),
		"declares fewer than present": junitSuiteXML(`tests="0"`, `<testcase name="ok"/>`),
	} {
		rd, err := JUnitReader{}.Read(context.Background(), ReportInput{Root: "/r", Data: data})
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if !rd.Partial || !strings.Contains(rd.Reason, "declare a test count") {
			t.Errorf("%s: partial = %v reason = %q", name, rd.Partial, rd.Reason)
		}
	}
	rd, err := JUnitReader{}.Read(context.Background(), ReportInput{Root: "/r", Data: []byte(`<testsuites tests="0"><testsuite tests="0"/></testsuites>`)})
	if err != nil || rd.Partial || rd.Reason != noTestsRanReason {
		t.Errorf("a genuine empty suite: err %v partial %v reason %q", err, rd.Partial, rd.Reason)
	}
}

func TestCoverageReadersRejectInvalidMeasurements(t *testing.T) {
	root := scopedGoRoot(t)
	for name, data := range map[string]string{
		"go unsupported mode": "mode: garbage\nexample.com/m/a/a.go:1.1,2.1 2 1\n",
		"go bad coordinates":  "mode: set\nexample.com/m/a/a.go:bad,bad 2 1\n",
		"go partial span":     "mode: set\nexample.com/m/a/a.go:1.1,2 2 1\n",
		"go negative count":   "mode: count\nexample.com/m/a/a.go:1.1,2.1 2 -1\n",
	} {
		if _, err := (GoCoverReader{}).Read(context.Background(), ReportInput{Root: root, Data: []byte(data)}); err == nil {
			t.Errorf("%s: want an error", name)
		}
	}
	for _, mode := range []string{"set", "count", "atomic"} {
		in := ReportInput{Root: root, Data: []byte("mode: " + mode + "\nexample.com/m/a/a.go:1.1,2.1 2 1\n")}
		if _, err := (GoCoverReader{}).Read(context.Background(), in); err != nil {
			t.Errorf("mode %s: %v", mode, err)
		}
	}
	for name, data := range map[string]string{
		"vitest null statement":     `{"a/a.ts":{"s":{"0":null}}}`,
		"vitest negative statement": `{"a/a.ts":{"s":{"0":-1}}}`,
		"vitest null function":      `{"a/a.ts":{"s":{"0":1},"f":{"0":null}}}`,
		"vitest negative branch":    `{"a/a.ts":{"s":{"0":1},"b":{"0":[1,-1]}}}`,
		"vitest null branch":        `{"a/a.ts":{"s":{"0":1},"b":{"0":[null]}}}`,
	} {
		if _, err := (VitestCoverageReader{}).Read(context.Background(), ReportInput{Root: vitestRoot, Data: []byte(data)}); err == nil {
			t.Errorf("%s: want an error, not a measurement", name)
		}
	}
	rd, err := VitestCoverageReader{}.Read(context.Background(), ReportInput{Root: vitestRoot, Data: []byte(`{"a/a.ts":{"s":{"0":0,"1":2}}}`)})
	if err != nil || rd.Value.Number != 50 {
		t.Errorf("a zero counter is a genuine uncovered statement: %v %+v", err, rd.Value)
	}
}

const scopedCoveragePyReport = `{"meta":{"version":"7.4","branch_coverage":true},
 "files":{
  "a/a.py":{"summary":{"covered_lines":1,"num_statements":1,"num_branches":2,"covered_branches":2}},
  "b/b.py":{"summary":{"covered_lines":0,"num_statements":1,"num_branches":2,"covered_branches":0}}},
 "totals":{"covered_lines":1,"num_statements":2,"num_branches":4,"covered_branches":2}}`

func TestCoveragePyReaderScopesItsHeadline(t *testing.T) {
	in := ReportInput{Provider: ProviderCoveragePyJSON, Root: "/r", Data: []byte(scopedCoveragePyReport), Scope: []string{"a/**"}}
	rd, err := CoveragePyReader{}.Read(context.Background(), in)
	if err != nil {
		t.Fatal(err)
	}
	got := detailMap(rd)
	if rd.Value.Number != 100 || got[DetailTotalStatements] != 1 || got[DetailTotalBranches] != 2 || got[DetailCoveredBranches] != 2 || rd.Partial {
		t.Errorf("value %v details %v partial %v, want the a/ file alone", rd.Value.Number, got, rd.Partial)
	}
	if rd.Provenance.ProviderVersion != "7.4" || len(rd.Evidence) != 0 {
		t.Errorf("version %q evidence %v", rd.Provenance.ProviderVersion, rd.Evidence)
	}

	in.Scope, in.Exclusions = nil, []string{"a/**"}
	rd, err = CoveragePyReader{}.Read(context.Background(), in)
	if err != nil || rd.Value.Number != 0 || len(rd.Evidence) != 1 {
		t.Errorf("exclude a: %v %+v %v", err, rd.Value, rd.Evidence)
	}

	in.Exclusions = nil
	rd, err = CoveragePyReader{}.Read(context.Background(), in)
	if err != nil || rd.Value.Number != 50 {
		t.Errorf("no scope keeps the report's totals: %v %+v", err, rd.Value)
	}

	in.Scope = []string{"nowhere/**"}
	if _, err := (CoveragePyReader{}).Read(context.Background(), in); err == nil {
		t.Error("a scope holding no statements has nothing to measure")
	}
}

func TestCoveragePyReaderIsPartialForInScopeFilesWithoutCounts(t *testing.T) {
	data := `{"files":{"a/a.py":{"summary":{"covered_lines":1,"num_statements":1}},"a/b.py":{"summary":{}}},"totals":{"covered_lines":1,"num_statements":1}}`
	rd, err := CoveragePyReader{}.Read(context.Background(), ReportInput{Root: "/r", Data: []byte(data), Scope: []string{"a/**"}})
	if err != nil {
		t.Fatal(err)
	}
	if !rd.Partial || !strings.Contains(rd.Reason, "lack valid statement counts") {
		t.Errorf("partial = %v reason = %q", rd.Partial, rd.Reason)
	}
}

func jscpdReportJSON(sources, clones string) []byte {
	return []byte(`{"statistics":{"total":{"lines":100,"sources":2,"clones":1,"duplicatedLines":20,"percentage":20},` + sources + `},"duplicates":[` + clones + `]}`)
}

const jscpdSources = `"formats":{"typescript":{"sources":{"a/c.ts":{"lines":40,"duplicatedLines":0},"b/b.ts":{"lines":60,"duplicatedLines":20}}}}`

func TestJscpdReaderScopesItsTotalsAndEvidence(t *testing.T) {
	bOnly := `{"lines":20,"firstFile":{"name":"b/b.ts","start":1},"secondFile":{"name":"b/c.ts","start":2}}`
	mixed := `{"lines":20,"firstFile":{"name":"b/b.ts","start":1},"secondFile":{"name":"a/c.ts","start":9}}`

	rd, err := JscpdReader{}.Read(context.Background(), ReportInput{Root: "/r", Scope: []string{"a/**"}, Data: jscpdReportJSON(jscpdSources, bOnly)})
	if err != nil {
		t.Fatal(err)
	}
	got := detailMap(rd)
	if rd.Value.Number != 0 || got[DetailTotalLines] != 40 || got[DetailClones] != 0 || len(rd.Evidence) != 0 {
		t.Errorf("all clones outside scope: value %v details %v evidence %v", rd.Value.Number, got, rd.Evidence)
	}

	rd, err = JscpdReader{}.Read(context.Background(), ReportInput{Root: "/r", Scope: []string{"a/**"}, Data: jscpdReportJSON(jscpdSources, mixed)})
	if err != nil {
		t.Fatal(err)
	}
	want := []EvidenceRef{{Path: "a/c.ts", Line: 9, Note: "20 lines repeated at b/b.ts:1"}}
	if !equalRefs(rd.Evidence, want) || detailMap(rd)[DetailClones] != 1 {
		t.Errorf("mixed clone: evidence %v, want %v", rd.Evidence, want)
	}
	for _, e := range rd.Evidence {
		if !(ReportInput{Scope: []string{"a/**"}}).inputScope().relevant(e.Path) {
			t.Errorf("evidence path %s is outside the scope", e.Path)
		}
	}

	rd, err = JscpdReader{}.Read(context.Background(), ReportInput{Root: "/r", Exclusions: []string{"a/**"}, Data: jscpdReportJSON(jscpdSources, mixed)})
	if err != nil || rd.Value.Number != 100.0*20/60 || detailMap(rd)[DetailSources] != 1 {
		t.Errorf("exclude a: %v value %+v details %v", err, rd.Value, detailMap(rd))
	}
}

func TestJscpdReaderWithoutPerFileCountsKeepsTheReportTotal(t *testing.T) {
	rd, err := JscpdReader{}.Read(context.Background(), ReportInput{Root: "/r", Scope: []string{"a/**"}, Data: jscpdReportJSON(`"x":1`, "")})
	if err != nil {
		t.Fatal(err)
	}
	if rd.Value == nil || rd.Value.Number != 20 || rd.Partial {
		t.Errorf("value %v partial %v, want the report's own total kept", rd.Value, rd.Partial)
	}
	rd, err = JscpdReader{}.Read(context.Background(), ReportInput{Root: "/r", Data: jscpdReportJSON(`"x":1`, "")})
	if err != nil || rd.Value == nil || rd.Value.Number != 20 {
		t.Errorf("an unscoped report keeps its totals: %v %+v", err, rd.Value)
	}
}

func osvGroupsReport(packages string) []byte {
	return []byte(`{"results":[{"source":{"path":"go.mod"},"packages":[` + packages + `]}]}`)
}

func TestOSVReaderCountsEachDistinctGroupOnce(t *testing.T) {
	pkg := func(name string, groups ...string) string {
		return `{"package":{"name":"` + name + `","version":"1"},"groups":[` + strings.Join(groups, ",") + `]}`
	}
	for _, tt := range []struct {
		name     string
		packages string
		total    float64
		high     float64
		items    int
	}{
		{"identical groups", pkg("x", `{"ids":["X"],"max_severity":"7"}`, `{"ids":["X"],"max_severity":"7"}`), 1, 1, 1},
		{"reordered aliases", pkg("x", `{"ids":["A","B"],"max_severity":"7"}`, `{"ids":["B","A"],"max_severity":"7"}`), 1, 1, 1},
		{"different groups", pkg("x", `{"ids":["A"],"max_severity":"7"}`, `{"ids":["B"],"max_severity":"7"}`), 2, 2, 2},
		{"same group in two packages", pkg("x", `{"ids":["X"],"max_severity":"7"}`) + "," + pkg("y", `{"ids":["X"],"max_severity":"7"}`), 2, 2, 2},
	} {
		t.Run(tt.name, func(t *testing.T) {
			rd, err := OSVScannerReader{}.Read(context.Background(), ReportInput{Root: "/r", Data: osvGroupsReport(tt.packages)})
			if err != nil {
				t.Fatal(err)
			}
			if rd.Value.Number != tt.total || detailMap(rd)[DetailHighVulnerabilities] != tt.high || len(rd.Evidence) != tt.items {
				t.Errorf("value %v high %v items %d, want %v %v %d", rd.Value.Number, detailMap(rd)[DetailHighVulnerabilities], len(rd.Evidence), tt.total, tt.high, tt.items)
			}
		})
	}
}

func TestCollectKeepsUnusableReadingsOutOfMeasuredResults(t *testing.T) {
	root := project(t)
	write(t, root, "go.mod", "module example.com/m\n")
	write(t, root, "reports/go.jsonl", "null\n")
	write(t, root, "reports/vitest-cov.json", `{"a/a.ts":{"s":{"0":null}}}`)
	write(t, root, "reports/go.out", "mode: garbage\nexample.com/m/a/a.go:1.1,2.1 2 1\n")
	write(t, root, "reports/mismatch.xml", string(junitSuiteXML(`tests="2"`, `<testcase name="ok" classname="a/ok.test.ts"/>`)))
	write(t, root, "reports/good.xml", string(junitSuiteXML(`tests="1"`, `<testcase name="ok" classname="a/ok.test.ts"/>`)))
	cfg := cfgOf(
		CapabilityConfig{Capability: CapabilityTests, Provider: ProviderGoTestJSON, Name: "null", Report: "reports/go.jsonl"},
		CapabilityConfig{Capability: CapabilityCoverage, Provider: ProviderVitestV8, Name: "nullcount", Report: "reports/vitest-cov.json"},
		CapabilityConfig{Capability: CapabilityCoverage, Provider: ProviderGoCoverProfile, Name: "mode", Report: "reports/go.out"},
		CapabilityConfig{Capability: CapabilityTests, Provider: ProviderVitestJUnit, Name: "mismatch", Report: "reports/mismatch.xml"},
		CapabilityConfig{Capability: CapabilityTests, Provider: ProviderPytestJUnit, Name: "good", Report: "reports/good.xml"},
	)
	got := collect(t, root, cfg, DefaultReaders(), &fakeTools{t: t})

	find := func(name string) CapabilityResult {
		for _, c := range got.Results {
			if c.Result.Name == name {
				return c.Result
			}
		}
		t.Fatalf("no result named %q", name)
		return CapabilityResult{}
	}
	for _, name := range []string{"null", "nullcount", "mode"} {
		if r := find(name); r.Outcome != OutcomeFailed || r.Value != nil {
			t.Errorf("%s = %s with value %v, want failed without a value", name, r.Outcome, r.Value)
		}
	}
	if r := find("mismatch"); r.Outcome != OutcomePartial {
		t.Errorf("mismatch = %s, want partial", r.Outcome)
	}
	if r := find("good"); !r.Outcome.Measured() || r.Value == nil || r.Value.Number != 0 || r.Outcome == OutcomePartial {
		t.Errorf("sibling = %s %v, want an intact measured value", r.Outcome, r.Value)
	}
	snaps, err := NewStore(root).LoadSnapshots()
	if err != nil || len(snaps) != 1 {
		t.Fatalf("LoadSnapshots = %d, %v", len(snaps), err)
	}
	for _, s := range snaps[0].Summary.Capabilities {
		switch s.Name {
		case "null", "nullcount", "mode":
			if s.Classification != ClassificationUnknown {
				t.Errorf("%s classified %s, want unknown", s.Name, s.Classification)
			}
		case "mismatch":
			// A partial reading is never reported as good.
			if s.Classification == ClassificationGood {
				t.Errorf("mismatch classified %s, want less than good", s.Classification)
			}
		}
	}
}
