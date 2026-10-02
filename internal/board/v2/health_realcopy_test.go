package v2

import (
	"context"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	xansi "github.com/charmbracelet/x/ansi"
	"github.com/muesli/termenv"
	"github.com/opencode/savepoint/internal/codehealth"
)

// realCopyRows are the five signals with the lengths the real projection
// writes: the longest signal name, long values and the real default aims.
func realCopyRows() []codehealth.DashboardRow {
	rows := fiveSignals()
	words := []struct{ value, figure, aim string }{
		{"all 3,045 pass", "0", "aim for none failing"},
		{"86%", "86%", "aim for 80% or more"},
		{"hardest function scores 46", "46", "aim for 15 or less"},
		{"8% copy-pasted", "8%", "aim for 5% or less"},
		{"2 unrated", "2", "aim for none"},
	}
	for i, w := range words {
		rows[i].Value, rows[i].Figure, rows[i].Aim = w.value, w.figure, w.aim
		rows[i].Spark, rows[i].SparkWord, rows[i].SparkNote = "▁▄█", "better", ""
	}
	return rows
}

func TestHealthPopoverKeepsEveryYardstickWithRealCopy(t *testing.T) {
	rows := realCopyRows()
	for _, size := range [][2]int{{80, 20}, {80, 24}, {80, 40}, {200, 50}} {
		m := openHealthScreenAt(t, &fakeHealth{dashboard: measuredDashboard(rows)}, size[0], size[1])
		view := xansi.Strip(screen(m))
		popoverBounds(t, view, size[0], size[1])
		for _, row := range rows {
			requireContains(t, view, row.Aim)
		}
		requireContains(t, view, "▁▄█ better")
		requireContains(t, strings.Join(strings.Fields(view), " "), "Tests failing 0", "Coverage 86%", "Complexity 46", "Duplication 8%", "Dependencies 2")
		if got := lipgloss.Width(m.renderHealthPopover(size[0])); got > 72 {
			t.Errorf("%v: popover is %d cells wide, the cap is 72", size, got)
		}
		if strings.Contains(view, "wor…") || strings.Count(view, "better") < 5 {
			t.Errorf("%v: direction word cut:\n%s", size, view)
		}
	}
}

func TestHealthPopoverShowsSparkNotesWithAndWithoutBlocks(t *testing.T) {
	const early, restarted = "early", "restarted because settings changed"
	thin := "not enough history yet (2 of 3)"
	for name, c := range map[string]struct{ spark, word, note string }{
		"early":             {"▁▄█", "better", early},
		"restarted":         {"▁▄█", "better", restarted},
		"early and restart": {"▁▄█", "better", early + " " + restarted},
		"thin":              {"", "", thin},
	} {
		t.Run(name, func(t *testing.T) {
			rows := realCopyRows()
			for i := range rows {
				rows[i].Spark, rows[i].SparkWord, rows[i].SparkNote = c.spark, c.word, c.note
			}
			for _, size := range [][2]int{{80, 20}, {80, 24}, {120, 40}} {
				m := openHealthScreenAt(t, &fakeHealth{dashboard: measuredDashboard(rows)}, size[0], size[1])
				for range rows {
					view := xansi.Strip(screen(m))
					popoverBounds(t, view, size[0], size[1])
					requireContains(t, view, c.note, rows[m.Health.Cursor].Meaning, "Next: ", "Sign-off: ", "Where: ")
					if c.spark != "" {
						requireContains(t, view, c.spark+" "+c.word)
					}
					m = press(t, m, "down")
				}
			}
		})
	}
}

func TestHealthRefreshNeverNamesAProvider(t *testing.T) {
	providers := map[codehealth.Capability]codehealth.ProviderKey{
		codehealth.CapabilityTests:                   "go_test_json",
		codehealth.CapabilityCoverage:                "go_cover",
		codehealth.CapabilityComplexity:              "gocyclo_cli",
		codehealth.CapabilityDuplication:             "jscpd",
		codehealth.CapabilityDependencyVulnerability: "govulncheck",
	}
	for i, c := range []codehealth.Capability{
		codehealth.CapabilityTests, codehealth.CapabilityCoverage, codehealth.CapabilityComplexity,
		codehealth.CapabilityDuplication, codehealth.CapabilityDependencyVulnerability,
	} {
		for _, history := range []bool{false, true} {
			m, _ := startRefresh(t, openHealthScreenAt(t, &fakeHealth{dashboard: measuredDashboard(realCopyRows()),
				refresh: func(ctx context.Context, _ func(codehealth.Progress)) error { <-ctx.Done(); return ctx.Err() }}, 80, 24))
			m.Health.History = history
			m = feed(t, m, healthProgressMsg{Progress: codehealth.Progress{Position: i + 1, Total: 5, Capability: c, Provider: providers[c]}})
			view := xansi.Strip(screen(m))
			requireContains(t, view, "Refreshing "+string(rune('1'+i))+" of 5: "+codehealth.CapabilityText(c), "Esc to cancel")
			if strings.Contains(view, string(providers[c])) {
				t.Errorf("%s history=%v shows a provider:\n%s", c, history, view)
			}
		}
	}
}

func TestHealthPopoverShowsOneRowPerSignalWithoutScrolling(t *testing.T) {
	base := realCopyRows()
	named := func(row codehealth.DashboardRow, name string, label codehealth.Classification, labelText string) codehealth.DashboardRow {
		row.Name, row.Question, row.Label, row.LabelText = name, row.Question+"?", label, labelText
		return row
	}
	blocking := base[4]
	blocking.Label, blocking.LabelText, blocking.SignOff, blocking.Meaning = codehealth.ClassificationNeedsAttention, "Needs Attention", "Blocks sign-off", "A serious library problem."
	cases := map[string][]codehealth.DashboardRow{
		"two tests first": {base[0], named(base[0], "integration", codehealth.ClassificationNeedsAttention, "Needs Attention"), base[1], base[2], base[3], blocking},
		"two tests and complexity": {base[0], named(base[0], "integration", codehealth.ClassificationWatch, "Watch"), base[1], base[2],
			named(base[2], "scripts", codehealth.ClassificationNeedsAttention, "Needs Attention"), base[3], blocking},
	}
	for name, rows := range cases {
		t.Run(name, func(t *testing.T) {
			for _, size := range [][2]int{{80, 20}, {80, 24}, {80, 40}} {
				m := openHealthScreenAt(t, &fakeHealth{dashboard: measuredDashboard(rows)}, size[0], size[1])
				first := ""
				for step := 0; step < 5; step++ {
					view := xansi.Strip(screen(m))
					popoverBounds(t, view, size[0], size[1])
					signals := strings.Join(strings.Split(view, "\n")[3:8], "\n")
					signals = strings.NewReplacer("▸", " ").Replace(signals)
					if first == "" {
						first = signals
					} else if signals != first {
						t.Fatalf("%v step %d: rows moved:\n%s\n--- was\n%s", size, step, signals, first)
					}
					requireContains(t, view, "Tests failing ×2", "Dependencies")
					m = press(t, m, "down")
				}
				if m.Health.Cursor != 4 {
					t.Fatalf("%v: cursor at %d, want the fifth signal", size, m.Health.Cursor)
				}
				requireContains(t, xansi.Strip(screen(m)), "Sign-off: Blocks sign-off", "A serious library problem.")
				requireContains(t, strings.Join(strings.Fields(xansi.Strip(screen(m))), " "), "Dependencies 2!")
				m = press(t, m, "k", "k", "k", "k")
				requireContains(t, xansi.Strip(screen(m)), "(worst of 2: integration)")
			}
		})
	}
}

func TestHealthSparkLastBlockFollowsTheValueNotTheConfidence(t *testing.T) {
	forceColorProfile(t, termenv.TrueColor)
	good := healthLabelStyle(codehealth.ClassificationGood).Render("█")
	watch := healthLabelStyle(codehealth.ClassificationWatch).Render("█")
	if good == watch {
		t.Fatal("good and watch render alike, so colour cannot be checked")
	}
	rows := realCopyRows()[:1]
	rows[0].Label, rows[0].ValueLabel = codehealth.ClassificationWatch, codehealth.ClassificationGood
	d := measuredDashboard(rows)
	m := openHealthScreenAt(t, &fakeHealth{dashboard: d}, 80, 24)
	line := m.Health.rows()[0]
	got := strings.Join(healthSignalRows(m.Health, 68), "\n")
	if !strings.Contains(got, good) || strings.Contains(got, watch) {
		t.Errorf("good value on stale evidence: last block should be drawn good:\n%q", got)
	}
	line.ValueLabel = ""
	m.Health.Dashboard.Rows[0] = line
	if got := strings.Join(healthSignalRows(m.Health, 68), "\n"); !strings.Contains(got, watch) {
		t.Errorf("without a value label the block follows the label:\n%q", got)
	}
}

func TestHealthRowMarkFollowsTheValue(t *testing.T) {
	rows := realCopyRows()
	rows[0].Label, rows[0].ValueLabel = codehealth.ClassificationWatch, codehealth.ClassificationGood
	m := openHealthScreenAt(t, &fakeHealth{dashboard: measuredDashboard(rows)}, 80, 24)
	view := strings.Join(strings.Fields(xansi.Strip(screen(m))), " ")
	requireContains(t, view, "✓ Tests failing 0")
	if strings.Contains(view, "~ Tests failing") {
		t.Errorf("a good value on doubtful evidence still shows the Watch mark:\n%s", view)
	}
}

func TestHealthPopoverNextStepPointsAtTheReportAtEverySize(t *testing.T) {
	for _, step := range []string{
		"Ask your agent to investigate .savepoint/health/report.md",
		"Run savepoint health report, then ask your agent to investigate it.",
	} {
		rows := realCopyRows()
		rows[0].Label, rows[0].LabelText, rows[0].ReportStep = codehealth.ClassificationNeedsAttention, "Needs Attention", step
		rows[0].NextStep = "Ask your agent to fix the failing tests."
		rows[0].Where = "13 files"
		for _, size := range [][2]int{{80, 20}, {80, 24}, {80, 40}} {
			m := openHealthScreenAt(t, &fakeHealth{dashboard: measuredDashboard(rows)}, size[0], size[1])
			view := xansi.Strip(screen(m))
			popoverBounds(t, view, size[0], size[1])
			flat := strings.Join(strings.Fields(strings.ReplaceAll(view, "│", " ")), " ")
			requireContains(t, flat, "Next: "+step, "Where: 13 files")
			if strings.Contains(flat, "fix the failing tests") || strings.Contains(flat, "lists them") {
				t.Errorf("%v: stale wording:\n%s", size, view)
			}
		}
	}
}
