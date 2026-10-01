package codehealth

import (
	"strings"
	"testing"
)

func seq(from float64, step float64, n int) []float64 {
	out := make([]float64, n)
	for i := range out {
		out[i] = from + step*float64(i)
	}
	return out
}

func TestSparklinePointCounts(t *testing.T) {
	tests := []struct {
		name     string
		values   []float64
		wantLen  int
		wantNote string
	}{
		{"one point", seq(50, 10, 1), 0, "not enough history yet (1 of 3)"},
		{"two points", seq(50, 10, 2), 0, "not enough history yet (2 of 3)"},
		{"three points", seq(50, 10, 3), 3, "early"},
		{"four points", seq(50, 10, 4), 4, "early"},
		{"five points", seq(50, 10, 5), 5, ""},
		{"twelve points", seq(10, 5, 12)[2:], 10, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spark, _, note := sparkline(CapabilityCoverage, tt.values, TrendImproving, false)
			if got := len([]rune(spark)); got != tt.wantLen {
				t.Fatalf("spark %q has %d blocks, want %d", spark, got, tt.wantLen)
			}
			if note != tt.wantNote {
				t.Fatalf("note = %q, want %q", note, tt.wantNote)
			}
		})
	}
}

func TestSparkValuesKeepsTheNewestTen(t *testing.T) {
	s := series{values: seq(1, 1, 11)}
	cur := dashResult(CapabilityConfig{Capability: CapabilityTests, Provider: ProviderGoTestJSON}, 99)
	got := sparkValues(cur, OriginOfficial, s)
	if len(got) != MaxSparkPoints || got[0] != 3 || got[len(got)-1] != 99 {
		t.Fatalf("values = %v, want the newest %d ending at the current value", got, MaxSparkPoints)
	}
	if got := sparkValues(cur, OriginManual, s); len(got) != MaxSparkPoints || got[len(got)-1] != 11 {
		t.Fatalf("manual newest must not contribute: %v", got)
	}
}

func TestSparklineEarlyBoundary(t *testing.T) {
	_, _, note := sparkline(CapabilityCoverage, seq(50, 10, 4), TrendImproving, false)
	if note != "early" {
		t.Fatalf("4 points note = %q, want early", note)
	}
	if _, _, note := sparkline(CapabilityCoverage, seq(50, 10, 5), TrendImproving, false); note != "" {
		t.Fatalf("5 points note = %q, want none", note)
	}
}

func TestSparklineFlatAndBelowMaterialDrawTheMiddleBlock(t *testing.T) {
	for name, values := range map[string][]float64{
		"flat":   {80, 80, 80, 80},
		"wobble": {80, 82, 81, 83}, // spread 3 < coverage's 5
	} {
		spark, _, _ := sparkline(CapabilityCoverage, values, TrendSteady, false)
		if spark != strings.Repeat("▄", 4) {
			t.Errorf("%s: spark = %q, want all middle blocks", name, spark)
		}
	}
	// At the material size the shape is drawn.
	if spark, _, _ := sparkline(CapabilityCoverage, []float64{80, 85, 80}, TrendSteady, false); spark != "▁█▁" {
		t.Errorf("material spread spark = %q, want ▁█▁", spark)
	}
}

func TestSparklineTallerIsAlwaysBiggerNumber(t *testing.T) {
	// The same shape for a higher-is-better and a lower-is-better signal;
	// only the word differs.
	rising := []float64{10, 20, 30, 40, 50}
	falling := []float64{50, 40, 30, 20, 10}
	for _, c := range []Capability{CapabilityCoverage, CapabilityComplexity} {
		if spark, _, _ := sparkline(c, rising, TrendNone, false); spark != "▁▃▅▆█" {
			t.Errorf("%s rising spark = %q", c, spark)
		}
		if spark, _, _ := sparkline(c, falling, TrendNone, false); spark != "█▆▅▃▁" {
			t.Errorf("%s falling spark = %q", c, spark)
		}
	}
}

func TestSparklineWordFollowsTrendDirection(t *testing.T) {
	for dir, want := range map[TrendDirection]string{
		TrendImproving: "better", TrendDeclining: "worse", TrendSteady: "steady", TrendNone: "",
	} {
		if _, word, _ := sparkline(CapabilityCoverage, seq(1, 10, 5), dir, false); word != want {
			t.Errorf("direction %q word = %q, want %q", dir, word, want)
		}
	}
}

func TestSparklineRestartedNote(t *testing.T) {
	_, _, note := sparkline(CapabilityCoverage, seq(50, 10, 8), TrendSteady, true)
	if note != "restarted because settings changed" {
		t.Fatalf("note = %q", note)
	}
	_, _, note = sparkline(CapabilityCoverage, seq(50, 10, 1), TrendNone, true)
	if note != "not enough history yet (1 of 3) restarted because settings changed" {
		t.Fatalf("thin and restarted note = %q", note)
	}
}

// sparkRow saves official snapshots of coverage or complexity values and
// returns the row of the newest.
func sparkRow(t *testing.T, c Capability, values []float64, tweak func(i int, r *CapabilityResult), origins func(i int) Origin) DashboardRow {
	t.Helper()
	store, root := dashProject(t)
	cfg := dashConfig(t, store, c)
	cc := cfg.Capabilities[0]
	for i, v := range values {
		r := dashResult(cc, v)
		if tweak != nil {
			tweak(i, &r)
		}
		origin := OriginOfficial
		if origins != nil {
			origin = origins(i)
		}
		saveDash(t, store, cfg, origin, i+1, append([]CapabilityResult{r}, notConfiguredFor(c)...), dashRepo(i+1))
	}
	return rowFor(t, mustLoad(t, root), c)
}

func TestDashboardSparklineThroughTheRow(t *testing.T) {
	// Coverage rising: better. Complexity rising: worse.
	row := sparkRow(t, CapabilityCoverage, []float64{60, 70, 80, 90, 95, 98}, nil, nil)
	if row.Spark == "" || row.SparkWord != "better" || row.SparkNote != "" {
		t.Fatalf("coverage rising = %q / %q / %q", row.Spark, row.SparkWord, row.SparkNote)
	}
	row = sparkRow(t, CapabilityComplexity, []float64{5, 8, 12, 16, 20, 25}, nil, nil)
	if row.SparkWord != "worse" || []rune(row.Spark)[5] != '█' {
		t.Fatalf("complexity rising = %q / %q", row.Spark, row.SparkWord)
	}
	row = sparkRow(t, CapabilityCoverage, []float64{90, 80, 70, 60, 50, 40}, nil, nil)
	if row.SparkWord != "worse" {
		t.Fatalf("coverage falling word = %q", row.SparkWord)
	}
	row = sparkRow(t, CapabilityComplexity, []float64{25, 20, 16, 12, 8, 5}, nil, nil)
	if row.SparkWord != "better" {
		t.Fatalf("complexity falling word = %q", row.SparkWord)
	}
}

func TestDashboardSparklineThinHistory(t *testing.T) {
	row := sparkRow(t, CapabilityCoverage, []float64{90, 88}, nil, nil)
	if row.Spark != "" || row.SparkWord != "" || row.SparkNote != "not enough history yet (2 of 3)" {
		t.Fatalf("two points = %q / %q / %q", row.Spark, row.SparkWord, row.SparkNote)
	}
}

func TestDashboardSparklineIgnoresManualPartialAndKeepsTrend(t *testing.T) {
	origins := func(i int) Origin {
		if i == 2 {
			return OriginManual
		}
		return OriginOfficial
	}
	tweak := func(i int, r *CapabilityResult) {
		if i == 3 {
			r.Outcome, r.Reason = OutcomePartial, "only part of the scope was measured"
		}
	}
	// Official: 60, 70, [manual 99], [partial 10], 80, 90 -> four usable points.
	row := sparkRow(t, CapabilityCoverage, []float64{60, 70, 99, 10, 80, 90}, tweak, origins)
	if got := []rune(row.Spark); len(got) != 4 || got[0] != '▁' || got[3] != '█' {
		t.Fatalf("spark = %q, want 4 points from 60 to 90", row.Spark)
	}
	if row.SparkNote != "early" || row.SparkWord != "better" {
		t.Fatalf("note/word = %q / %q", row.SparkNote, row.SparkWord)
	}
}

func TestDashboardSparklineRestartedSeries(t *testing.T) {
	tweak := func(i int, r *CapabilityResult) {
		if i >= 4 {
			r.Provenance.ProviderVersion = "2.0"
		}
	}
	row := sparkRow(t, CapabilityCoverage, []float64{60, 70, 80, 90, 85, 86, 87}, tweak, nil)
	if got := []rune(row.Spark); len(got) != 3 {
		t.Fatalf("spark = %q, want only the 3 points since the restart", row.Spark)
	}
	if row.SparkNote != "early restarted because settings changed" {
		t.Fatalf("note = %q", row.SparkNote)
	}
}

func TestDashboardSparklineLeavesTrendWordingAlone(t *testing.T) {
	row := sparkRow(t, CapabilityCoverage, []float64{90, 88, 86}, nil, nil)
	if !strings.Contains(row.Trend, "over 3 official checks") || !strings.Contains(row.Basis, "2 earlier comparable official checks") {
		t.Fatalf("trend/basis = %q / %q", row.Trend, row.Basis)
	}
}
