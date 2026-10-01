package codehealth

import (
	"math"
	"reflect"
	"strings"
	"testing"
	"unicode/utf8"
)

const testDigest = "sha256:0000000000000000000000000000000000000000000000000000000000000000"

var testProvider = map[Capability]ProviderKey{
	CapabilityTests:                   ProviderGoTestJSON,
	CapabilityCoverage:                ProviderGoCoverProfile,
	CapabilityComplexity:              ProviderLizardCSV,
	CapabilityDuplication:             ProviderJscpdJSON,
	CapabilityDependencyVulnerability: ProviderOSVScannerJSON,
}

// result builds a fresh, available result in a fixed series.
func result(c Capability, v float64) CapabilityResult {
	return CapabilityResult{
		Capability:   c,
		Outcome:      OutcomeAvailable,
		Freshness:    FreshnessFresh,
		CollectedAt:  "2026-10-01T00:00:00Z",
		ConfigDigest: testDigest,
		Provenance: Provenance{
			Provider: testProvider[c], ProviderVersion: "1", MeasurementDefinition: "def-1",
		},
		Value: &Value{Number: v, Unit: capabilityUnits[c]},
	}
}

func vulnResult(total float64, details ...Detail) CapabilityResult {
	r := result(CapabilityDependencyVulnerability, total)
	r.Details = details
	return r
}

func withOutcome(r CapabilityResult, o Outcome, f Freshness) CapabilityResult {
	r.Outcome, r.Freshness = o, f
	if !o.Measured() {
		r.Value = nil
	}
	return r
}

func official(c Capability, values ...float64) []HistoryEntry {
	var h []HistoryEntry
	for _, v := range values {
		h = append(h, HistoryEntry{Origin: OriginOfficial, Result: result(c, v)})
	}
	return h
}

func TestAssessCurrentValueThresholds(t *testing.T) {
	sev := func(high, critical float64) []Detail {
		return []Detail{{Key: "high", Number: high}, {Key: "critical", Number: critical}}
	}
	tests := []struct {
		name     string
		result   CapabilityResult
		override *Threshold
		want     Classification
	}{
		{"coverage at good boundary", result(CapabilityCoverage, 80), nil, ClassificationGood},
		{"coverage just below good", result(CapabilityCoverage, 79.9), nil, ClassificationWatch},
		{"coverage at watch boundary", result(CapabilityCoverage, 60), nil, ClassificationWatch},
		{"coverage just below watch", result(CapabilityCoverage, 59.9), nil, ClassificationNeedsAttention},
		{"coverage zero", result(CapabilityCoverage, 0), nil, ClassificationNeedsAttention},
		{"coverage full", result(CapabilityCoverage, 100), nil, ClassificationGood},
		{"complexity at good boundary", result(CapabilityComplexity, 10), nil, ClassificationGood},
		{"complexity just above good", result(CapabilityComplexity, 11), nil, ClassificationWatch},
		{"complexity at watch boundary", result(CapabilityComplexity, 20), nil, ClassificationWatch},
		{"complexity just above watch", result(CapabilityComplexity, 21), nil, ClassificationNeedsAttention},
		{"complexity zero", result(CapabilityComplexity, 0), nil, ClassificationGood},
		{"duplication at good boundary", result(CapabilityDuplication, 3), nil, ClassificationGood},
		{"duplication just above good", result(CapabilityDuplication, 3.1), nil, ClassificationWatch},
		{"duplication at watch boundary", result(CapabilityDuplication, 5), nil, ClassificationWatch},
		{"duplication just above watch", result(CapabilityDuplication, 5.1), nil, ClassificationNeedsAttention},
		{"duplication zero", result(CapabilityDuplication, 0), nil, ClassificationGood},
		{"no failing tests", result(CapabilityTests, 0), nil, ClassificationGood},
		{"one failing test", result(CapabilityTests, 1), nil, ClassificationNeedsAttention},
		{"no vulnerabilities", vulnResult(0), nil, ClassificationGood},
		{"no vulnerabilities, severity absent", vulnResult(0), nil, ClassificationGood},
		{"low-severity only", vulnResult(3, sev(0, 0)...), nil, ClassificationWatch},
		{"one high", vulnResult(3, sev(1, 0)...), nil, ClassificationNeedsAttention},
		{"one critical", vulnResult(1, sev(0, 1)...), nil, ClassificationNeedsAttention},
		{"severity missing entirely", vulnResult(2), nil, ClassificationNeedsAttention},
		{"critical count missing", vulnResult(2, Detail{Key: "high", Number: 0}), nil, ClassificationNeedsAttention},
		{"zero total with high", vulnResult(0, sev(1, 0)...), nil, ClassificationNeedsAttention},
		{"zero total with critical", vulnResult(0, sev(0, 1)...), nil, ClassificationNeedsAttention},
		{"zero total with zero severity", vulnResult(0, sev(0, 0)...), nil, ClassificationGood},
		{"high count missing", vulnResult(2, Detail{Key: "critical", Number: 0}), nil, ClassificationNeedsAttention},
		{"medium and low never block", vulnResult(4, append(sev(0, 0), Detail{Key: "medium", Number: 2}, Detail{Key: "low", Number: 2})...), nil, ClassificationWatch},
		{"one unknown blocks", vulnResult(2, append(sev(0, 0), Detail{Key: "unknown", Number: 1})...), nil, ClassificationNeedsAttention},
		{"zero unknown does not block", vulnResult(2, append(sev(0, 0), Detail{Key: "unknown", Number: 0})...), nil, ClassificationWatch},
		{"unknown blocks despite guidance", vulnResult(2, append(sev(0, 0), Detail{Key: "unknown", Number: 1})...), &Threshold{Good: 10, Watch: 10}, ClassificationNeedsAttention},
		// Configured guidance replaces the defaults...
		{"stricter coverage guidance", result(CapabilityCoverage, 85), &Threshold{Good: 90, Watch: 70}, ClassificationWatch},
		{"looser complexity guidance", result(CapabilityComplexity, 15), &Threshold{Good: 15, Watch: 30}, ClassificationGood},
		// ...but never weakens a hard minimum.
		{"guidance cannot excuse failing tests", result(CapabilityTests, 2), &Threshold{Good: 5, Watch: 5}, ClassificationNeedsAttention},
		{"guidance cannot excuse high vulnerability", vulnResult(2, sev(1, 0)...), &Threshold{Good: 10, Watch: 10}, ClassificationNeedsAttention},
		{"guidance cannot excuse unknown severity", vulnResult(2), &Threshold{Good: 10, Watch: 10}, ClassificationNeedsAttention},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := 0.0
			if tc.result.Value != nil {
				v = tc.result.Value.Number
			}
			// Two earlier equal checks plus this one satisfy the three-check
			// rule without any movement.
			got := Assess(OriginOfficial, tc.result, tc.override, official(tc.result.Capability, v, v))
			if got.Classification != tc.want {
				t.Fatalf("classification = %s, want %s (%s)", got.Classification, tc.want, got.Explanation)
			}
			if got.Explanation == "" || len(got.Explanation) > MaxReasonLen {
				t.Fatalf("explanation %q is empty or unbounded", got.Explanation)
			}
		})
	}
}

func TestAssessUnmeasuredAndIncompleteEvidence(t *testing.T) {
	good := result(CapabilityCoverage, 90)
	bad := result(CapabilityCoverage, 40)
	tests := []struct {
		name string
		in   CapabilityResult
		want Classification
		text string
	}{
		{"absent", withOutcome(good, OutcomeAbsent, FreshnessUnknown), ClassificationUnknown, "No report"},
		{"unsupported", withOutcome(good, OutcomeUnsupported, FreshnessUnknown), ClassificationUnknown, "not supported"},
		{"unavailable", withOutcome(good, OutcomeUnavailable, FreshnessUnknown), ClassificationUnknown, "could not be used"},
		{"not configured", withOutcome(good, OutcomeNotConfigured, FreshnessUnknown), ClassificationUnknown, "Not configured"},
		{"failed", withOutcome(good, OutcomeFailed, FreshnessUnknown), ClassificationUnknown, "failed"},
		{"timed out", withOutcome(good, OutcomeTimedOut, FreshnessUnknown), ClassificationUnknown, "timed out"},
		{"cancelled", withOutcome(good, OutcomeCancelled, FreshnessUnknown), ClassificationUnknown, "cancelled"},
		{"unrecognised outcome", withOutcome(good, Outcome("bogus"), FreshnessUnknown), ClassificationUnknown, "not recognised"},
		{"available without a value", func() CapabilityResult { r := good; r.Value = nil; return r }(), ClassificationUnknown, "no value"},
		{"NaN value", result(CapabilityCoverage, math.NaN()), ClassificationUnknown, "not a finite"},
		{"infinite value", result(CapabilityCoverage, math.Inf(1)), ClassificationUnknown, "not a finite"},
		{"negative value", result(CapabilityCoverage, -1), ClassificationUnknown, "not a finite"},
		{"partial good value is only watch", withOutcome(good, OutcomePartial, FreshnessFresh), ClassificationWatch, "part of the evidence"},
		{"partial poor value still needs attention", withOutcome(bad, OutcomePartial, FreshnessFresh), ClassificationNeedsAttention, "part of the evidence"},
		{"stale good value is only watch", withOutcome(good, OutcomeAvailable, FreshnessStale), ClassificationWatch, "stale"},
		{"stale poor value still needs attention", withOutcome(bad, OutcomeAvailable, FreshnessStale), ClassificationNeedsAttention, "stale"},
		{"unknown-freshness good value is only watch", withOutcome(good, OutcomeAvailable, FreshnessUnknown), ClassificationWatch, "Freshness is unknown"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Assess(OriginOfficial, tc.in, nil, nil)
			if got.Classification != tc.want {
				t.Fatalf("classification = %s, want %s (%s)", got.Classification, tc.want, got.Explanation)
			}
			if !strings.Contains(got.Explanation, tc.text) {
				t.Fatalf("explanation %q does not mention %q", got.Explanation, tc.text)
			}
			if got.Classification == ClassificationGood {
				t.Fatalf("incomplete evidence must never be good")
			}
		})
	}
	t.Run("partial without a value is unknown", func(t *testing.T) {
		r := withOutcome(good, OutcomePartial, FreshnessFresh)
		r.Value = nil
		if got := Assess(OriginOfficial, r, nil, nil); got.Classification != ClassificationUnknown {
			t.Fatalf("got %s", got.Classification)
		}
	})
}

func TestAssessMaterialDeclineAndImprovement(t *testing.T) {
	tests := []struct {
		name    string
		cap     Capability
		history []float64
		current float64
		want    Classification
		text    string // required in the explanation; empty means none
	}{
		// Good base, worsened one level at exactly the minimum movement.
		{"complexity +5 from good", CapabilityComplexity, []float64{5, 5, 5}, 10, ClassificationWatch, "Materially worse"},
		{"complexity +4 is not material", CapabilityComplexity, []float64{5, 5, 5}, 9, ClassificationGood, ""},
		{"duplication +2 from good", CapabilityDuplication, []float64{1, 1, 1}, 3, ClassificationWatch, "Materially worse"},
		{"duplication +1.9 is not material", CapabilityDuplication, []float64{1, 1, 1}, 2.9, ClassificationGood, ""},
		{"coverage -5 from good", CapabilityCoverage, []float64{95, 95, 95}, 90, ClassificationWatch, "Materially worse"},
		{"coverage -4.9 is not material", CapabilityCoverage, []float64{95, 95, 95}, 90.1, ClassificationGood, ""},
		// Watch base worsens to needs attention.
		{"coverage -5 from watch", CapabilityCoverage, []float64{70, 70, 70}, 65, ClassificationNeedsAttention, "Materially worse"},
		{"complexity +5 from watch", CapabilityComplexity, []float64{12, 12, 12}, 17, ClassificationNeedsAttention, "Materially worse"},
		// Needs attention cannot fall further.
		{"needs attention stays", CapabilityCoverage, []float64{58, 58, 58}, 40, ClassificationNeedsAttention, "Materially worse"},
		// Only one level, however large the drop.
		{"one level only", CapabilityCoverage, []float64{100, 100, 100}, 81, ClassificationWatch, "Materially worse"},
		{"one level only, big drop", CapabilityCoverage, []float64{100, 100, 100}, 80, ClassificationWatch, "Materially worse"},
		// Improvement explains but never upgrades a poor current value.
		{"improvement leaves watch as watch", CapabilityCoverage, []float64{50, 50, 50}, 70, ClassificationWatch, "Improved from"},
		{"improvement leaves needs attention", CapabilityCoverage, []float64{30, 30, 30}, 40, ClassificationNeedsAttention, "Improved from"},
		{"small improvement is silent", CapabilityCoverage, []float64{50, 50, 50}, 54, ClassificationNeedsAttention, ""},
		{"improved complexity stays good", CapabilityComplexity, []float64{9, 9, 9}, 2, ClassificationGood, "Improved from"},
		// The baseline is the median of the last three, not the mean or the latest.
		{"median ignores an outlier", CapabilityCoverage, []float64{95, 10, 95}, 91, ClassificationGood, ""},
		{"only the last three count", CapabilityCoverage, []float64{10, 10, 95, 95, 95}, 90, ClassificationWatch, "Materially worse"},
		// Fewer than three earlier observations give no baseline.
		{"two observations give no baseline", CapabilityComplexity, []float64{5, 5}, 10, ClassificationGood, ""},
		{"no history gives no baseline", CapabilityComplexity, nil, 10, ClassificationWatch, "Good needs"},
		// Counts: any increase is material, and hard blockers already need attention.
		{"vulnerability total rises", CapabilityDependencyVulnerability, []float64{0, 0, 0}, 1, ClassificationNeedsAttention, "Materially worse"},
		{"vulnerability total flat", CapabilityDependencyVulnerability, []float64{2, 2, 2}, 2, ClassificationWatch, ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cur := result(tc.cap, tc.current)
			if tc.cap == CapabilityDependencyVulnerability {
				cur.Details = []Detail{{Key: "high", Number: 0}, {Key: "critical", Number: 0}}
			}
			got := Assess(OriginOfficial, cur, nil, official(tc.cap, tc.history...))
			if got.Classification != tc.want {
				t.Fatalf("classification = %s, want %s (%s)", got.Classification, tc.want, got.Explanation)
			}
			has := strings.Contains(got.Explanation, "Materially worse") || strings.Contains(got.Explanation, "Improved from")
			if tc.text == "" && has {
				t.Fatalf("unexpected movement wording: %s", got.Explanation)
			}
			if tc.text != "" && !strings.Contains(got.Explanation, tc.text) {
				t.Fatalf("explanation %q does not mention %q", got.Explanation, tc.text)
			}
		})
	}
}

func TestAssessTestsAndHardBlockersIgnoreDecline(t *testing.T) {
	// A hard blocker is already at the floor; history cannot change or soften it.
	history := official(CapabilityTests, 5, 5, 5)
	got := Assess(OriginOfficial, result(CapabilityTests, 1), nil, history)
	if got.Classification != ClassificationNeedsAttention {
		t.Fatalf("got %s", got.Classification)
	}
	if strings.Contains(got.Explanation, "Improved") {
		t.Fatalf("a hard blocker must not read as improving: %s", got.Explanation)
	}
}

func TestAssessTrendWording(t *testing.T) {
	tests := []struct {
		name         string
		cap          Capability
		history      []float64
		current      float64
		observations int
		min, max     float64
		direction    TrendDirection
	}{
		{"no history", CapabilityCoverage, nil, 80, 0, 0, 0, TrendNone},
		{"two observations", CapabilityCoverage, []float64{80}, 81, 0, 0, 0, TrendNone},
		{"three steady", CapabilityCoverage, []float64{80, 82}, 81, 3, 80, 82, TrendSteady},
		{"improving", CapabilityCoverage, []float64{70, 75}, 80, 3, 70, 80, TrendImproving},
		{"declining", CapabilityCoverage, []float64{80, 75}, 70, 3, 70, 80, TrendDeclining},
		{"complexity rising is declining", CapabilityComplexity, []float64{5, 8}, 10, 3, 5, 10, TrendDeclining},
		{"complexity falling is improving", CapabilityComplexity, []float64{12, 9}, 7, 3, 7, 12, TrendImproving},
		{"small net movement is steady", CapabilityDuplication, []float64{2, 3}, 3.5, 3, 2, 3.5, TrendSteady},
		{"exact threshold movement counts", CapabilityDuplication, []float64{1, 2}, 3, 3, 1, 3, TrendDeclining},
		{"all ties", CapabilityCoverage, []float64{80, 80}, 80, 3, 80, 80, TrendSteady},
		{"all zeros", CapabilityTests, []float64{0, 0}, 0, 3, 0, 0, TrendSteady},
		{"zero start", CapabilityTests, []float64{0, 0}, 0, 3, 0, 0, TrendSteady},
		{"window is the newest five", CapabilityCoverage, []float64{10, 60, 70, 80, 90}, 95, 5, 60, 95, TrendImproving},
		{"down then back is steady", CapabilityCoverage, []float64{80, 70}, 80, 3, 70, 80, TrendSteady},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Assess(OriginOfficial, result(tc.cap, tc.current), nil, official(tc.cap, tc.history...)).Trend
			want := Trend{Observations: tc.observations, HasRange: tc.observations > 0, Min: tc.min, Max: tc.max, Direction: tc.direction}
			if got != want {
				t.Fatalf("trend = %+v, want %+v", got, want)
			}
		})
	}
}

func TestAssessTrendExactWording(t *testing.T) {
	tests := []struct {
		name    string
		cap     Capability
		history []float64
		current float64
		want    string
	}{
		{"short history", CapabilityCoverage, []float64{80}, 82,
			"Coverage 82% meets the good threshold. Not enough comparable official history for a recent range yet. Good needs three comparable official checks."},
		{"range", CapabilityCoverage, []float64{80, 84}, 82,
			"Coverage 82% meets the good threshold. Recent range 80% to 84% over 3 official checks, steady."},
		{"unchanged", CapabilityComplexity, []float64{4, 4}, 4,
			"Highest complexity 4 meets the good threshold. Unchanged at 4 over 3 official checks."},
		{"rounded", CapabilityDuplication, []float64{1.04, 1.06}, 1.26,
			"Duplication 1.3% meets the good threshold. Recent range 1% to 1.3% over 3 official checks, steady."},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Assess(OriginOfficial, result(tc.cap, tc.current), nil, official(tc.cap, tc.history...))
			if got.Explanation != tc.want {
				t.Fatalf("explanation =\n%q\nwant\n%q", got.Explanation, tc.want)
			}
		})
	}
}

func TestAssessHistoryFiltering(t *testing.T) {
	cap := CapabilityComplexity
	other := func(mutate func(*CapabilityResult), values ...float64) []HistoryEntry {
		h := official(cap, values...)
		for i := range h {
			mutate(&h[i].Result)
		}
		return h
	}
	manual := func(values ...float64) []HistoryEntry {
		h := official(cap, values...)
		for i := range h {
			h[i].Origin = OriginManual
		}
		return h
	}
	concat := func(parts ...[]HistoryEntry) []HistoryEntry {
		var out []HistoryEntry
		for _, p := range parts {
			out = append(out, p...)
		}
		return out
	}
	tests := []struct {
		name         string
		history      []HistoryEntry
		origin       Origin
		current      float64
		want         Classification
		observations int
		text         []string
		notText      []string
	}{
		{"manual history never counts", manual(5, 5, 5), OriginOfficial, 10, ClassificationWatch, 0,
			[]string{"3 manual results shown, not counted"}, []string{"Materially"}},
		{"one manual result is singular", manual(5), OriginOfficial, 10, ClassificationWatch, 0,
			[]string{"1 manual result shown, not counted"}, nil},
		{"provider version change resets", other(func(r *CapabilityResult) { r.Provenance.ProviderVersion = "2" }, 5, 5, 5), OriginOfficial, 10, ClassificationWatch, 0,
			[]string{"3 earlier official results not compared"}, []string{"Materially"}},
		{"measurement definition change resets", other(func(r *CapabilityResult) { r.Provenance.MeasurementDefinition = "def-2" }, 5, 5, 5), OriginOfficial, 10, ClassificationWatch, 0,
			[]string{"not compared"}, []string{"Materially"}},
		{"report schema change resets", other(func(r *CapabilityResult) { r.Provenance.ReportSchema = "v2" }, 5, 5, 5), OriginOfficial, 10, ClassificationWatch, 0,
			[]string{"not compared"}, []string{"Materially"}},
		{"configuration change resets", other(func(r *CapabilityResult) { r.ConfigDigest = "sha256:" + strings.Repeat("1", 64) }, 5, 5, 5), OriginOfficial, 10, ClassificationWatch, 0,
			[]string{"not compared"}, []string{"Materially"}},
		{"scope change resets", other(func(r *CapabilityResult) { r.Scope = []string{"internal/**"} }, 5, 5, 5), OriginOfficial, 10, ClassificationWatch, 0,
			[]string{"not compared"}, []string{"Materially"}},
		{"exclusion change resets", other(func(r *CapabilityResult) { r.Exclusions = []string{"vendor/**"} }, 5, 5, 5), OriginOfficial, 10, ClassificationWatch, 0,
			[]string{"not compared"}, []string{"Materially"}},
		{"provider change resets", other(func(r *CapabilityResult) { r.Provenance.Provider = ProviderLizardCSV + "-x" }, 5, 5, 5), OriginOfficial, 10, ClassificationWatch, 0,
			[]string{"not compared"}, []string{"Materially"}},
		{"old series are skipped, new series counts", concat(other(func(r *CapabilityResult) { r.Provenance.ProviderVersion = "2" }, 1, 1, 1), official(cap, 5, 5, 5)), OriginOfficial, 10, ClassificationWatch, 4,
			[]string{"Materially worse", "3 earlier official results not compared"}, nil},
		{"scope order is irrelevant", other(func(r *CapabilityResult) { r.Scope = []string{"b/**", "a/**"} }, 5, 5, 5), OriginOfficial, 10, ClassificationWatch, 0, nil, nil},
		{"other capability is ignored", official(CapabilityCoverage, 95, 95, 95), OriginOfficial, 10, ClassificationWatch, 0, nil, []string{"not compared", "manual"}},
		{"failed history is dropped", []HistoryEntry{
			{Origin: OriginOfficial, Result: withOutcome(result(cap, 5), OutcomeFailed, FreshnessUnknown)},
			{Origin: OriginOfficial, Result: withOutcome(result(cap, 5), OutcomeAbsent, FreshnessUnknown)},
			{Origin: OriginOfficial, Result: withOutcome(result(cap, 5), OutcomeTimedOut, FreshnessUnknown)},
		}, OriginOfficial, 10, ClassificationWatch, 0, nil, []string{"not compared"}},
		{"partial history is dropped", []HistoryEntry{
			{Origin: OriginOfficial, Result: withOutcome(result(cap, 5), OutcomePartial, FreshnessFresh)},
			{Origin: OriginOfficial, Result: withOutcome(result(cap, 5), OutcomePartial, FreshnessFresh)},
			{Origin: OriginOfficial, Result: withOutcome(result(cap, 5), OutcomePartial, FreshnessFresh)},
		}, OriginOfficial, 10, ClassificationWatch, 0, nil, []string{"Materially"}},
		{"missing history values are dropped", []HistoryEntry{
			{Origin: OriginOfficial, Result: func() CapabilityResult { r := result(cap, 5); r.Value = nil; return r }()},
			{Origin: OriginOfficial, Result: result(cap, math.NaN())},
			{Origin: OriginOfficial, Result: result(cap, math.Inf(1))},
			{Origin: OriginOfficial, Result: result(cap, -3)},
		}, OriginOfficial, 10, ClassificationWatch, 0, nil, []string{"Materially"}},
		{"unknown origin is ignored", func() []HistoryEntry {
			h := official(cap, 5, 5, 5)
			for i := range h {
				h[i].Origin = "imported"
			}
			return h
		}(), OriginOfficial, 10, ClassificationWatch, 0, nil, []string{"Materially", "not compared"}},
		// A manual current result is judged against official history but never
		// joins it.
		{"manual current uses official baseline", official(cap, 5, 5, 5), OriginManual, 10, ClassificationWatch, 3,
			[]string{"Materially worse"}, nil},
		{"manual current is not an observation", official(cap, 5, 5), OriginManual, 10, ClassificationWatch, 0, nil, nil},
		{"partial official current is not an observation", official(cap, 5, 5), OriginOfficial, 10, ClassificationWatch, 0, nil, nil},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cur := result(cap, tc.current)
			if tc.name == "partial official current is not an observation" {
				cur = withOutcome(cur, OutcomePartial, FreshnessFresh)
				tc.want = ClassificationWatch // partial caps good at watch
			}
			got := Assess(tc.origin, cur, nil, tc.history)
			if got.Classification != tc.want {
				t.Fatalf("classification = %s, want %s (%s)", got.Classification, tc.want, got.Explanation)
			}
			if got.Trend.Observations != tc.observations {
				t.Fatalf("observations = %d, want %d (%s)", got.Trend.Observations, tc.observations, got.Explanation)
			}
			for _, s := range tc.text {
				if !strings.Contains(got.Explanation, s) {
					t.Errorf("explanation %q does not mention %q", got.Explanation, s)
				}
			}
			for _, s := range tc.notText {
				if strings.Contains(got.Explanation, s) {
					t.Errorf("explanation %q unexpectedly mentions %q", got.Explanation, s)
				}
			}
		})
	}
}

func TestOverallHasNoScoreAndNeverHidesProblems(t *testing.T) {
	c := func(cs ...Classification) []Assessment {
		var out []Assessment
		for _, x := range cs {
			out = append(out, Assessment{Classification: x})
		}
		return out
	}
	const (
		g = ClassificationGood
		w = ClassificationWatch
		n = ClassificationNeedsAttention
		u = ClassificationUnknown
	)
	tests := []struct {
		name string
		in   []Assessment
		want Classification
	}{
		{"none", nil, u},
		{"all good", c(g, g, g, g, g), g},
		{"one watch", c(g, w, g), w},
		{"one unknown", c(g, u, g), u},
		{"watch and unknown", c(w, u), u},
		{"one needs attention", c(g, w, u, n), n},
		{"needs attention first", c(n, g), n},
		{"single good", c(g), g},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := Overall(tc.in); got != tc.want {
				t.Fatalf("got %s, want %s", got, tc.want)
			}
		})
	}
}

func TestAssessmentSummaryMatchesPersistedShape(t *testing.T) {
	r := result(CapabilityCoverage, 90)
	a := Assess(OriginOfficial, r, nil, official(CapabilityCoverage, 90, 90))
	want := CapabilitySummary{
		Capability: CapabilityCoverage, Provider: ProviderGoCoverProfile,
		Classification: ClassificationGood, Explanation: a.Explanation,
	}
	if got := a.Summary(); got != want {
		t.Fatalf("summary = %+v, want %+v", got, want)
	}
}

func TestAssessIsDeterministicAndDoesNotMutateInput(t *testing.T) {
	history := official(CapabilityCoverage, 95, 90, 92, 91)
	before := make([]HistoryEntry, len(history))
	copy(before, history)
	cur := result(CapabilityCoverage, 85)
	first := Assess(OriginOfficial, cur, nil, history)
	second := Assess(OriginOfficial, cur, nil, history)
	if !reflect.DeepEqual(first, second) {
		t.Fatalf("results differ between runs:\n%+v\n%+v", first, second)
	}
	if !reflect.DeepEqual(history, before) {
		t.Fatal("history was mutated")
	}
}

func TestBoundTextKeepsUTF8AndLimit(t *testing.T) {
	long := strings.Repeat("é", MaxReasonLen)
	got := boundText(long)
	if len(got) > MaxReasonLen || !utf8.ValidString(got) || !strings.HasSuffix(got, "…") {
		t.Fatalf("bad bounded text: len=%d valid=%v", len(got), utf8.ValidString(got))
	}
	short := "short"
	if boundText(short) != short {
		t.Fatal("short text must be unchanged")
	}
	exact := strings.Repeat("a", MaxReasonLen)
	if boundText(exact) != exact {
		t.Fatal("text at the limit must be unchanged")
	}
}

func TestEverySummaryFitsTheSnapshotContract(t *testing.T) {
	// Every outcome and freshness pairing must yield a label and explanation
	// that a persisted summary accepts: bounded, single-line, and never good
	// without available, fresh evidence.
	for _, c := range Capabilities() {
		for _, o := range Outcomes() {
			for _, f := range []Freshness{FreshnessFresh, FreshnessStale, FreshnessUnknown} {
				r := result(c, 1)
				r.Outcome, r.Freshness = o, f
				if !o.Measured() {
					r.Value = nil
				}
				a := Assess(OriginOfficial, r, nil, nil)
				if len(a.Explanation) > MaxReasonLen || hasControl(a.Explanation) || a.Explanation == "" {
					t.Fatalf("%s/%s/%s: bad explanation %q", c, o, f, a.Explanation)
				}
				if a.Classification == ClassificationGood && (o != OutcomeAvailable || f != FreshnessFresh) {
					t.Fatalf("%s/%s/%s: good without available, fresh evidence", c, o, f)
				}
				if err := validateClassification("c", a.Classification); err != nil {
					t.Fatal(err)
				}
			}
		}
	}
}

func TestGoodIsCappedAtWatchUntilThreeComparableChecks(t *testing.T) {
	c := CapabilityCoverage
	manual := official(c, 90, 90, 90)
	for i := range manual {
		manual[i].Origin = OriginManual
	}
	incompatible := official(c, 90, 90, 90)
	for i := range incompatible {
		incompatible[i].Result.ConfigDigest = "sha256:" + strings.Repeat("1", 64)
	}
	tests := []struct {
		name    string
		history []HistoryEntry
		want    Classification
	}{
		{"no history", nil, ClassificationWatch},
		{"one earlier check", official(c, 90), ClassificationWatch},
		{"two earlier checks", official(c, 90, 90), ClassificationGood},
		{"manual history only", manual, ClassificationWatch},
		{"incompatible history only", incompatible, ClassificationWatch},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := Assess(OriginOfficial, result(c, 90), nil, tc.history); got.Classification != tc.want {
				t.Fatalf("got %s, want %s (%s)", got.Classification, tc.want, got.Explanation)
			}
		})
	}
	t.Run("poor values are not softened or double-noted", func(t *testing.T) {
		got := Assess(OriginOfficial, result(c, 50), nil, nil)
		if got.Classification != ClassificationNeedsAttention || strings.Contains(got.Explanation, "Good needs") {
			t.Fatalf("got %s: %s", got.Classification, got.Explanation)
		}
	})
}
