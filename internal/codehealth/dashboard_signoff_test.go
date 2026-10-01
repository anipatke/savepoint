package codehealth

import (
	"strings"
	"testing"
)

// signOffProject saves an official snapshot built from tweak and returns the
// loaded dashboard with the verdict Evaluate gives for the same snapshot.
func signOffProject(t *testing.T, tweak func(*CapabilityResult)) (Dashboard, Verdict) {
	t.Helper()
	store, root := dashProject(t)
	cfg := dashConfig(t, store, allCapabilities()...)
	snap := saveDash(t, store, cfg, OriginOfficial, 1, goodResults(cfg, tweak), dashRepo(1))
	verdict, err := Evaluate(snap, cfg)
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	return mustLoad(t, root), verdict
}

func TestDashboardSignOffMatchesEvaluate(t *testing.T) {
	vulnUnknown := func(r *CapabilityResult) {
		if r.Capability == CapabilityDependencyVulnerability {
			r.Value.Number = 2
			r.Details = []Detail{{Key: DetailUnknownVulnerabilities, Number: 2}}
		}
	}
	failingTests := func(r *CapabilityResult) {
		if r.Capability == CapabilityTests {
			r.Value.Number = 3
		}
	}
	for _, tc := range []struct {
		name       string
		tweak      func(*CapabilityResult)
		wantBlocks bool
	}{
		{"all good", nil, false},
		{"failing tests block", failingTests, true},
		{"unknown severity vulnerabilities are reported only", vulnUnknown, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d, v := signOffProject(t, tc.tweak)
			if v.Blocks() != tc.wantBlocks {
				t.Fatalf("Evaluate Blocks() = %v, want %v", v.Blocks(), tc.wantBlocks)
			}
			want := textSignOffClear
			if v.Blocks() {
				want = textSignOffBlocks
			}
			if d.SignOff != want {
				t.Errorf("Dashboard.SignOff = %q, want %q", d.SignOff, want)
			}
			for _, rv := range v.Results {
				row := rowFor(t, d, rv.Capability)
				want := map[Disposition]string{DispositionBlocks: textSignOffBlocks, DispositionReported: textSignOffAdvisory}[rv.Disposition]
				if row.SignOff != want {
					t.Errorf("%s SignOff = %q, want %q (%s)", rv.Capability, row.SignOff, want, rv.Disposition)
				}
			}
		})
	}
}

func TestDashboardSignOffBlocksAndAdvisoryRows(t *testing.T) {
	d, _ := signOffProject(t, func(r *CapabilityResult) {
		if r.Capability == CapabilityTests {
			r.Value.Number = 1
		}
	})
	if got := rowFor(t, d, CapabilityTests).SignOff; got != "Blocks sign-off" {
		t.Errorf("tests SignOff = %q", got)
	}
	if got := rowFor(t, d, CapabilityCoverage).SignOff; got != "Advisory only" {
		t.Errorf("coverage SignOff = %q", got)
	}
}

func TestDashboardSignOffManualNewest(t *testing.T) {
	store, root := dashProject(t)
	cfg := dashConfig(t, store, allCapabilities()...)
	saveDash(t, store, cfg, OriginOfficial, 1, goodResults(cfg, nil), dashRepo(1))
	saveDash(t, store, cfg, OriginManual, 2, goodResults(cfg, nil), dashRepo(2))
	d := mustLoad(t, root)
	if !strings.Contains(d.SignOff, "manual refresh does not affect sign-off") {
		t.Errorf("SignOff = %q", d.SignOff)
	}
	for _, row := range d.Rows {
		if row.SignOff != "" {
			t.Errorf("%s manual row SignOff = %q, want empty", row.Capability, row.SignOff)
		}
	}
}

func TestDashboardSignOffNoOfficial(t *testing.T) {
	store, root := dashProject(t)
	cfg := dashConfig(t, store, allCapabilities()...)
	saveDash(t, store, cfg, OriginManual, 1, goodResults(cfg, nil), dashRepo(1))
	d := mustLoad(t, root)
	if !strings.Contains(d.SignOff, "no official check yet") {
		t.Errorf("SignOff = %q", d.SignOff)
	}
}

func TestDashboardSignOffEvaluateErrorKeepsDashboard(t *testing.T) {
	store, root := dashProject(t)
	cfg := dashConfig(t, store, allCapabilities()...)
	snap := saveDash(t, store, cfg, OriginOfficial, 1, goodResults(cfg, nil), dashRepo(1))
	d := mustLoad(t, root)

	d.applySignOff(Config{Version: ConfigVersion + 1}, snap, true)
	if d.SignOff != "Sign-off status unavailable" {
		t.Errorf("SignOff = %q", d.SignOff)
	}
	for _, row := range d.Rows {
		if row.SignOff != "Sign-off status unavailable" {
			t.Errorf("%s SignOff = %q", row.Capability, row.SignOff)
		}
	}
	if len(d.Rows) == 0 || d.Headline == "" {
		t.Errorf("rest of the dashboard was dropped: %+v", d)
	}
}

func TestDashboardHeadline(t *testing.T) {
	rows := func(labels ...Classification) []DashboardRow {
		var out []DashboardRow
		for _, l := range labels {
			out = append(out, DashboardRow{Label: l})
		}
		return out
	}
	g, w, n, u := ClassificationGood, ClassificationWatch, ClassificationNeedsAttention, ClassificationUnknown
	for _, tc := range []struct {
		name string
		rows []DashboardRow
		want string
	}{
		{"all good", rows(g, g, g, g, g), "All 5 look fine"},
		{"two need a look", rows(g, w, n, g, g), "2 of 5 need a look"},
		{"one needs a look", rows(g, g, w, g, g), "1 of 5 needs a look"},
		{"unknown counts as needing a look", rows(g, u, g), "1 of 3 needs a look"},
		{"nothing judged", rows(u, u, u), "Not enough to judge yet"},
		{"no rows", nil, "Not enough to judge yet"},
	} {
		if got := headlineText(tc.rows); got != tc.want {
			t.Errorf("%s: headline = %q, want %q", tc.name, got, tc.want)
		}
	}
}

func TestDashboardDatesAreFixedUTC(t *testing.T) {
	store, root := dashProject(t)
	cfg := dashConfig(t, store, allCapabilities()...)
	saveDash(t, store, cfg, OriginOfficial, 1, goodResults(cfg, nil), dashRepo(1))
	saveDash(t, store, cfg, OriginOfficial, 2, goodResults(cfg, nil), dashRepo(2))
	d := mustLoad(t, root)
	if d.MeasuredText != "1 Oct 00:02 UTC" {
		t.Errorf("MeasuredText = %q", d.MeasuredText)
	}
	if len(d.History) != 2 || d.History[0].WhenText != "1 Oct 00:02 UTC" || d.History[1].WhenText != "1 Oct 00:01 UTC" {
		t.Errorf("History WhenText = %+v", d.History)
	}
	if got := whenText("2026-10-01T20:53:09Z"); got != "1 Oct 20:53 UTC" {
		t.Errorf("whenText = %q", got)
	}
	if got := whenText("not a time"); got != "not a time" {
		t.Errorf("whenText fallback = %q", got)
	}
}
