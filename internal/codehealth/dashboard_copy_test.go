package codehealth

import (
	"strings"
	"testing"
)

func wordsResult(c Capability, n float64, details ...Detail) CapabilityResult {
	r := dashResult(CapabilityConfig{Capability: c, Provider: dashProviders[c]}, n)
	r.Details = details
	return r
}

func TestQuestionsAreFixedPerSignal(t *testing.T) {
	want := map[Capability]string{
		CapabilityTests:                   "Do the tests pass?",
		CapabilityCoverage:                "How much code do the tests actually run?",
		CapabilityComplexity:              "How tangled is the hardest code?",
		CapabilityDuplication:             "How much is copy-pasted?",
		CapabilityDependencyVulnerability: "Known security problems in libraries we use",
	}
	for c, q := range want {
		if got := measuredWords(nil, wordsResult(c, 0), ClassificationGood).Question; got != q {
			t.Errorf("%s question = %q, want %q", c, got, q)
		}
	}
}

func TestValueReadsAsWords(t *testing.T) {
	tests := []struct {
		name string
		r    CapabilityResult
		want string
	}{
		{"all tests pass", wordsResult(CapabilityTests, 0, Detail{Key: "total_tests", Number: 3045}), "all 3,045 pass"},
		{"small pass count", wordsResult(CapabilityTests, 0, Detail{Key: "total_tests", Number: 12}), "all 12 pass"},
		{"failing tests", wordsResult(CapabilityTests, 2, Detail{Key: "total_tests", Number: 3045}), "2 failing"},
		{"no tests ran", wordsResult(CapabilityTests, 0, Detail{Key: "total_tests", Number: 0}), "no tests ran"},
		{"no total reported", wordsResult(CapabilityTests, 0), "none failing"},
		{"coverage", wordsResult(CapabilityCoverage, 86), "86%"},
		{"complexity", wordsResult(CapabilityComplexity, 46), "hardest function scores 46"},
		{"duplication", wordsResult(CapabilityDuplication, 8), "8% copy-pasted"},
		{"vulnerabilities none", wordsResult(CapabilityDependencyVulnerability, 0), "none"},
		{"vulnerabilities unrated", wordsResult(CapabilityDependencyVulnerability, 2, Detail{Key: DetailUnknownVulnerabilities, Number: 2}), "2 unrated"},
		{"vulnerabilities by severity", wordsResult(CapabilityDependencyVulnerability, 3,
			Detail{Key: DetailLowVulnerabilities, Number: 2}, Detail{Key: DetailHighVulnerabilities, Number: 1}), "1 high, 2 low"},
		{"vulnerabilities without severity", wordsResult(CapabilityDependencyVulnerability, 4), "4 found"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := measuredWords(nil, tt.r, ClassificationGood).Value; got != tt.want {
				t.Errorf("Value = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestUnmeasuredResultSaysNotMeasured(t *testing.T) {
	r := unmeasured(wordsResult(CapabilityCoverage, 0), OutcomeFailed)
	w := measuredWords(nil, r, ClassificationUnknown)
	if w.Value != "not measured" || w.Meaning != unknownCopy.Meaning || w.NextStep != unknownCopy.NextStep {
		t.Errorf("unmeasured words = %+v", w)
	}
	if w.Aim != "aim for 80% or more" {
		t.Errorf("an unmeasured signal still states its aim, got %q", w.Aim)
	}
}

func TestAimComesFromDefaultsOrConfiguredGood(t *testing.T) {
	tests := []struct {
		c    Capability
		conf *Threshold
		want string
	}{
		{CapabilityTests, nil, "aim for none failing"},
		{CapabilityCoverage, nil, "aim for 80% or more"},
		{CapabilityComplexity, nil, "aim for 10 or less"},
		{CapabilityDuplication, nil, "aim for 3% or less"},
		{CapabilityDependencyVulnerability, nil, "aim for none"},
		{CapabilityCoverage, &Threshold{Good: 90, Watch: 70}, "aim for 90% or more"},
		{CapabilityComplexity, &Threshold{Good: 15, Watch: 25}, "aim for 15 or less"},
		{CapabilityDuplication, &Threshold{Good: 1.5, Watch: 4}, "aim for 1.5% or less"},
		{CapabilityTests, &Threshold{Good: 2, Watch: 5}, "aim for 2 or less"},
	}
	for _, tt := range tests {
		if got := aimText(tt.c, tt.conf); got != tt.want {
			t.Errorf("aimText(%s, %v) = %q, want %q", tt.c, tt.conf, got, tt.want)
		}
	}
}

func TestMeaningAndNextStepExistForEverySignalAndLabel(t *testing.T) {
	labels := []Classification{ClassificationGood, ClassificationWatch, ClassificationNeedsAttention, ClassificationUnknown}
	for _, c := range Capabilities() {
		for _, l := range labels {
			e := meaningText[c][l]
			if e.Meaning == "" || e.NextStep == "" {
				t.Errorf("%s/%s has no meaning or next step: %+v", c, l, e)
			}
			w := measuredWords(nil, wordsResult(c, goodValues[c]), l)
			if l != ClassificationUnknown && (w.Meaning == "" || w.NextStep == "") {
				t.Errorf("%s/%s row words empty: %+v", c, l, w)
			}
		}
	}
}

func TestWatchBoundaryAppearsOnlyBeyondIt(t *testing.T) {
	tests := []struct {
		name  string
		r     CapabilityResult
		conf  *Threshold
		label Classification
		want  string // substring of Meaning; empty means no boundary sentence
	}{
		{"complexity beyond watch", wordsResult(CapabilityComplexity, 46), nil, ClassificationNeedsAttention, "past the watch line of 20 or less"},
		{"complexity within watch", wordsResult(CapabilityComplexity, 15), nil, ClassificationNeedsAttention, ""},
		{"complexity watch label", wordsResult(CapabilityComplexity, 15), nil, ClassificationWatch, ""},
		{"coverage beyond watch", wordsResult(CapabilityCoverage, 40), nil, ClassificationNeedsAttention, "past the watch line of 60% or more"},
		{"duplication beyond configured watch", wordsResult(CapabilityDuplication, 8), &Threshold{Good: 2, Watch: 6}, ClassificationNeedsAttention, "past the watch line of 6% or less"},
		{"failing tests have no watch line", wordsResult(CapabilityTests, 2), nil, ClassificationNeedsAttention, ""},
		{"vulnerabilities have no watch line", wordsResult(CapabilityDependencyVulnerability, 5), nil, ClassificationNeedsAttention, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := measuredWords(tt.conf, tt.r, tt.label).Meaning
			if tt.want == "" && strings.Contains(got, "watch line") {
				t.Errorf("Meaning mentions the watch line: %q", got)
			}
			if tt.want != "" && !strings.Contains(got, tt.want) {
				t.Errorf("Meaning = %q, want it to contain %q", got, tt.want)
			}
		})
	}
}

func TestWhereNamesOneFileOrCountsMany(t *testing.T) {
	ref := func(paths ...string) []EvidenceRef {
		var out []EvidenceRef
		for i, p := range paths {
			out = append(out, EvidenceRef{Path: p, Line: i + 1})
		}
		return out
	}
	tests := []struct {
		name string
		ev   []EvidenceRef
		want string
	}{
		{"no evidence", nil, ""},
		{"one path", ref("internal/doctor/repairs.go"), "internal/doctor/repairs.go"},
		{"many paths", ref("internal/doctor/repairs.go", "a.go", "b.go"), "3 files; savepoint health check lists them"},
		{"two paths", ref("internal/doctor/repairs.go", "a.go"), "2 files; savepoint health check lists them"},
		{"same file twice counts once", ref("a.go", "a.go", "b.go", "b.go"), "2 files; savepoint health check lists them"},
		{"one file at several lines names the file", ref("a.go", "a.go"), "a.go"},
	}
	for _, tt := range tests {
		if got := whereText(tt.ev); got != tt.want {
			t.Errorf("%s: whereText = %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestNotConfiguredRowSaysHowToSetItUp(t *testing.T) {
	w := unmeasuredWords(CapabilityDuplication, nil, true)
	if w.Value != "not measured" || w.Meaning != notConfiguredCopy.Meaning || w.NextStep != notConfiguredCopy.NextStep || w.Where != "" {
		t.Errorf("not-configured words = %+v", w)
	}
	if w.Aim != "aim for 3% or less" || w.Question == "" {
		t.Errorf("not-configured row keeps its question and aim: %+v", w)
	}
}

func TestDashboardRowsCarryPlainWords(t *testing.T) {
	store, root := dashProject(t)
	cfg := dashConfig(t, store, CapabilityTests, CapabilityComplexity, CapabilityDuplication)
	cfg.Capabilities[2].Thresholds = &Threshold{Good: 2, Watch: 6}
	if _, err := store.SaveConfig(cfg); err != nil {
		t.Fatal(err)
	}
	results := goodResults(cfg, func(r *CapabilityResult) {
		switch r.Capability {
		case CapabilityTests:
			r.Details = []Detail{{Key: "total_tests", Number: 3045}}
		case CapabilityComplexity:
			r.Value.Number = 46
			r.Evidence = []EvidenceRef{{Path: "internal/doctor/repairs.go", Line: 10}, {Path: "a.go", Line: 3}}
		case CapabilityDuplication:
			r.Value.Number = 8
			r.Evidence = nil
		}
	})
	saveDash(t, store, cfg, OriginOfficial, 1, results, dashRepo(1))
	d := mustLoad(t, root)

	tests := rowFor(t, d, CapabilityTests)
	if tests.Question != "Do the tests pass?" || tests.Value != "all 3,045 pass" || tests.Aim != "aim for none failing" {
		t.Errorf("tests row = %+v", tests)
	}
	cx := rowFor(t, d, CapabilityComplexity)
	if cx.Value != "hardest function scores 46" || cx.Aim != "aim for 10 or less" ||
		cx.Where != "2 files; savepoint health check lists them" || cx.NextStep != "Ask your agent to split the most tangled function." {
		t.Errorf("complexity row = %+v", cx)
	}
	if !strings.Contains(cx.Meaning, "past the watch line of 20 or less") {
		t.Errorf("complexity meaning lacks the watch line: %q", cx.Meaning)
	}
	dup := rowFor(t, d, CapabilityDuplication)
	if dup.Aim != "aim for 2% or less" || dup.Where != "" || !strings.Contains(dup.Meaning, "6% or less") {
		t.Errorf("a configured threshold must override the default: %+v", dup)
	}
	cov := rowFor(t, d, CapabilityCoverage)
	if !cov.NotConfigured || cov.Value != "not measured" || cov.Aim != "aim for 80% or more" || cov.NextStep != notConfiguredCopy.NextStep {
		t.Errorf("not-configured row = %+v", cov)
	}
}

func TestDashboardInstanceMissingFromSnapshotHasWords(t *testing.T) {
	store, root := dashProject(t)
	cfg := dashConfig(t, store, CapabilityCoverage)
	saveDash(t, store, cfg, OriginOfficial, 1, goodResults(cfg, nil), dashRepo(1))
	cfg.Capabilities = append(cfg.Capabilities, CapabilityConfig{Capability: CapabilityDuplication, Provider: ProviderJscpdJSON, Scope: []string{"**/*.go"}})
	if _, err := store.SaveConfig(cfg); err != nil {
		t.Fatal(err)
	}
	row := rowFor(t, mustLoad(t, root), CapabilityDuplication)
	if row.Value != "not measured" || row.Aim != "aim for 3% or less" || row.Meaning != unknownCopy.Meaning {
		t.Errorf("no-result row = %+v", row)
	}
}

func TestPlainWordingNeverRepeatsBlockingClaimsOrBannedClaims(t *testing.T) {
	var all []string
	for _, q := range questionText {
		all = append(all, q)
	}
	for _, a := range zeroAimText {
		all = append(all, a)
	}
	for _, byLabel := range meaningText {
		for _, e := range byLabel {
			all = append(all, e.Meaning, e.NextStep)
		}
	}
	all = append(all, notConfiguredCopy.Meaning, notConfiguredCopy.NextStep,
		textNotMeasured, textNoTestsRan, textNoneFailing, textAllPass, textFailing, textHardest,
		textCopyPasted, textFound, textUnrated, textAimMore, textAimLess, textPastWatch)
	text := strings.ToLower(strings.Join(all, "\n"))
	assertNoBannedClaims(t, "plain wording", text)
	if strings.Contains(text, "treated as blocking") || strings.Contains(text, "blocks") {
		t.Errorf("plain wording must not make blocking claims:\n%s", text)
	}
	for _, c := range Capabilities() {
		if questionText[c] == "" {
			t.Errorf("no question for %s", c)
		}
	}
}

func TestGroupDigits(t *testing.T) {
	for in, want := range map[float64]string{0: "0", 999: "999", 1000: "1,000", 3045: "3,045", 1234567: "1,234,567"} {
		if got := groupDigits(in); got != want {
			t.Errorf("groupDigits(%v) = %q, want %q", in, got, want)
		}
	}
}

// The label is one word for several causes. The plain reading must name the
// cause that actually produced it, and never tell anyone to repair a reading
// that already meets the aim.
func TestMeaningNamesTheCauseOfTheLabel(t *testing.T) {
	tests := []struct {
		name      string
		c         Capability
		values    []float64
		tweak     func(i int, r *CapabilityResult)
		label     Classification
		want      []string
		notWanted []string
	}{
		{"coverage with one check", CapabilityCoverage, []float64{90}, nil, ClassificationWatch,
			[]string{"three comparable checks"}, []string{"skip", "add tests"}},
		{"coverage with two checks", CapabilityCoverage, []float64{90, 90}, nil, ClassificationWatch,
			[]string{"three comparable checks"}, []string{"skip", "add tests"}},
		{"passing tests with one check", CapabilityTests, []float64{0}, nil, ClassificationWatch,
			[]string{"three comparable checks"}, []string{"incomplete", "out of date"}},
		{"low complexity with two checks", CapabilityComplexity, []float64{5, 5}, nil, ClassificationWatch,
			[]string{"three comparable checks"}, []string{"hard to follow", "split"}},
		{"low duplication with one check", CapabilityDuplication, []float64{1}, nil, ClassificationWatch,
			[]string{"three comparable checks"}, []string{"copy-pasted, so", "merge"}},
		{"coverage with three checks", CapabilityCoverage, []float64{90, 90, 90}, nil, ClassificationGood,
			[]string{"tests run most"}, []string{"three comparable"}},
		{"partial coverage", CapabilityCoverage, []float64{90, 90, 90, 90}, func(i int, r *CapabilityResult) {
			if i == 3 {
				r.Outcome, r.Reason = OutcomePartial, "one package did not report"
			}
		}, ClassificationWatch, []string{"only part was measured"}, []string{"skip", "add tests"}},
		{"stale coverage", CapabilityCoverage, []float64{90, 90, 90, 90}, func(i int, r *CapabilityResult) {
			if i == 3 {
				r.Freshness = FreshnessStale
			}
		}, ClassificationWatch, []string{"out of date"}, []string{"skip", "add tests"}},
		{"coverage worse than recent checks", CapabilityCoverage, []float64{95, 95, 95, 95, 80}, nil, ClassificationWatch,
			[]string{"worse than recent checks"}, []string{"three comparable"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row := sparkRow(t, tt.c, tt.values, tt.tweak, nil)
			if row.Label != tt.label {
				t.Fatalf("setup: label = %s, want %s", row.Label, tt.label)
			}
			text := row.Meaning + " | " + row.NextStep
			if n := len([]rune(row.Meaning)); n > 68 {
				t.Errorf("meaning is %d cells and would be cut at 80 columns: %q", n, row.Meaning)
			}
			for _, want := range tt.want {
				if !strings.Contains(text, want) {
					t.Errorf("reading %q does not say %q", text, want)
				}
			}
			for _, banned := range tt.notWanted {
				if strings.Contains(text, banned) {
					t.Errorf("reading %q wrongly says %q", text, banned)
				}
			}
		})
	}
}

func TestFigureIsTheNumberAlone(t *testing.T) {
	tests := []struct {
		r    CapabilityResult
		want string
	}{
		{wordsResult(CapabilityTests, 0, Detail{Key: "total_tests", Number: 3045}), "0"},
		{wordsResult(CapabilityTests, 2, Detail{Key: "total_tests", Number: 3045}), "2"},
		{wordsResult(CapabilityCoverage, 86.3), "86.3%"},
		{wordsResult(CapabilityComplexity, 46), "46"},
		{wordsResult(CapabilityDuplication, 9), "9%"},
		{wordsResult(CapabilityDependencyVulnerability, 0), "none"},
		{wordsResult(CapabilityDependencyVulnerability, 3, Detail{Key: DetailHighVulnerabilities, Number: 1}, Detail{Key: DetailLowVulnerabilities, Number: 2}), "3"},
		{unmeasured(wordsResult(CapabilityCoverage, 0), OutcomeFailed), "—"},
	}
	for _, tt := range tests {
		if got := measuredWords(nil, tt.r, ClassificationGood).Figure; got != tt.want {
			t.Errorf("%s figure = %q, want %q", tt.r.Capability, got, tt.want)
		}
	}
}
