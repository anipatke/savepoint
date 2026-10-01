package v2

import (
	"context"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	xansi "github.com/charmbracelet/x/ansi"
	"github.com/opencode/savepoint/internal/codehealth"
)

func popoverBounds(t *testing.T, view string, width, height int) {
	t.Helper()
	lines := strings.Split(view, "\n")
	if len(lines) > height {
		t.Errorf("%d lines overflow a %d-line terminal", len(lines), height)
	}
	for _, line := range lines {
		if lipgloss.Width(line) > width {
			t.Errorf("line is %d cells, terminal is %d: %q", lipgloss.Width(line), width, xansi.Strip(line))
		}
	}
}

func TestHealthPopoverFitsEverySupportedSize(t *testing.T) {
	rows := fiveSignals()
	rows[2].Meaning = strings.Repeat("A very long explanation that must be cut rather than wrap. ", 4)
	f := &fakeHealth{dashboard: measuredDashboard(rows), fresh: codehealth.CodeFreshness{Text: strings.Repeat("Your code has moved on. ", 6)}}
	for _, size := range [][2]int{{80, 20}, {80, 24}, {80, 40}} {
		m := openHealthScreenAt(t, f, size[0], size[1])
		for step := 0; step < 5; step++ {
			view := xansi.Strip(screen(m))
			popoverBounds(t, view, size[0], size[1])
			requireContains(t, view, "CODE HEALTH", "2 Oct 2026 · Official check · 2 signals need a look.", "Re-run: savepoint health check")
			if got := strings.Count(m.renderHealthPopover(m.terminalWidth()), "\n") + 1; got != healthPopoverHeight {
				t.Errorf("%v step %d: popover is %d lines, want %d", size, step, got, healthPopoverHeight)
			}
			m = press(t, m, "down")
		}
	}
}

func TestHealthPopoverRowsShowMarkValueSparkAndAim(t *testing.T) {
	m := openHealthScreen(t, &fakeHealth{dashboard: measuredDashboard(fiveSignals())})
	got := xansi.Strip(screen(m))

	requireContains(t, got, "✓ Tests", "~ Coverage", "✗ Complexity", "? Duplication", "▁▃▅▇ better", "aim: 80 or more", "42")
	for _, banned := range []string{"go_test_json", "snap-1", "2026-10-02T09:00:00Z", "RECENT CHECKS", "SIGNAL DETAIL"} {
		if strings.Contains(got, banned) {
			t.Errorf("popover shows %q:\n%s", banned, got)
		}
	}
}

func TestHealthPopoverSelectionChangesTheExplanation(t *testing.T) {
	m := openHealthScreen(t, &fakeHealth{dashboard: measuredDashboard(fiveSignals())})
	requireContains(t, xansi.Strip(screen(m)), "Tests?", "All 120 tests pass.", "Sign-off: Advisory only", "Where: internal/foo/foo.go", "Next: Do the next thing.", "Doesn't block sign-off")
	m = press(t, m, "down")
	got := xansi.Strip(screen(m))
	requireContains(t, got, "Coverage?", "Coverage is 61%, down from 68%.")
	if strings.Contains(got, "Tests?") {
		t.Errorf("explanation did not follow the selection:\n%s", got)
	}
	if m = press(t, m, "k"); m.Health.Cursor != 0 {
		t.Errorf("k left the cursor at %d", m.Health.Cursor)
	}
	if m = press(t, m, "up"); m.Health.Cursor != 0 {
		t.Errorf("up moved above the first signal")
	}
	if m = press(t, m, "down", "down", "down", "down", "down", "down"); m.Health.Cursor != 4 {
		t.Errorf("down went past the last signal: %d", m.Health.Cursor)
	}
}

func TestHealthPopoverColoursEachLabel(t *testing.T) {
	for label, style := range map[codehealth.Classification]string{
		codehealth.ClassificationGood:           healthLabelStyle(codehealth.ClassificationGood).Render("x"),
		codehealth.ClassificationWatch:          healthLabelStyle(codehealth.ClassificationWatch).Render("x"),
		codehealth.ClassificationNeedsAttention: healthLabelStyle(codehealth.ClassificationNeedsAttention).Render("x"),
		codehealth.ClassificationUnknown:        healthLabelStyle(codehealth.ClassificationUnknown).Render("x"),
	} {
		row := codehealth.DashboardRow{Label: label, Spark: "▁▇", SparkWord: "worse"}
		got := healthSpark(row, healthLabelStyle(label))
		if !strings.HasSuffix(got, " worse") || xansi.Strip(got) != "▁▇ worse" {
			t.Errorf("%s: spark %q", label, got)
		}
		if strings.Contains(got, "▇") && got != "▁"+strings.Replace(style, "x", "▇", 1)+" worse" {
			t.Errorf("%s: last block is not coloured by the label: %q", label, got)
		}
		if healthGlyph[label] == "" {
			t.Errorf("%s has no mark", label)
		}
	}
}

func TestHealthPopoverSignOffLines(t *testing.T) {
	for _, text := range []string{"Blocks sign-off", "A manual refresh does not affect sign-off.", "There is no official check yet, so nothing is judged for sign-off."} {
		d := measuredDashboard(fiveSignals())
		d.SignOff = text
		requireContains(t, xansi.Strip(screen(openHealthScreen(t, &fakeHealth{dashboard: d}))), text)
	}
}

func TestHealthPopoverRerunCommand(t *testing.T) {
	f := &fakeHealth{dashboard: measuredDashboard(fiveSignals())}

	m := openHealthScreen(t, f)
	want := "Re-run: savepoint health check " + m.Objectives[m.ObjectiveCursor].ID()
	requireContains(t, xansi.Strip(screen(m)), want)

	b := healthBoard(t, f, 120, 40)
	b.Objectives = nil
	b.State.Router.Objective = "O-077"
	if got := b.healthRerunCommand(); got != "Re-run: savepoint health check O-077" {
		t.Errorf("router fallback = %q", got)
	}
	b.State.Router = nil
	if got := b.healthRerunCommand(); got != "Re-run: savepoint health check O-###" {
		t.Errorf("no objective = %q", got)
	}

	f.refresh = func(ctx context.Context, _ func(codehealth.Progress)) error { <-ctx.Done(); return ctx.Err() }
	m, _ = startRefresh(t, m)
	defer m.cancelHealthRefresh()
	if strings.Contains(xansi.Strip(screen(m)), "Re-run:") {
		t.Errorf("re-run line shown during a refresh")
	}
}

func TestHealthPopoverRemovedKeysDoNothing(t *testing.T) {
	m := openHealthScreen(t, &fakeHealth{dashboard: measuredDashboard(fiveSignals())})
	before := xansi.Strip(screen(m))
	for _, key := range []string{"enter", "v", "pgdown", "pgup", "end", "home"} {
		next, cmd := sendKey(t, m, key)
		if cmd != nil || next.Health == nil || next.Health.Cursor != 0 || xansi.Strip(screen(next)) != before {
			t.Errorf("%q changed the popover", key)
		}
	}
}

var popoverSizes = [][2]int{{80, 20}, {80, 24}, {80, 40}}

func TestHealthPopoverFirstRunStatesFitEverySize(t *testing.T) {
	for name, tc := range map[string]struct {
		d    codehealth.Dashboard
		want []string
	}{
		"not set up": {codehealth.Dashboard{State: codehealth.DashboardNotConfigured}, []string{"not set up", "savepoint health setup"}},
		"no check":   {codehealth.Dashboard{State: codehealth.DashboardFirstRun}, []string{"Nothing has been measured yet", "Press R", "Re-run: savepoint health check"}},
	} {
		for _, size := range popoverSizes {
			m := openHealthScreenAt(t, &fakeHealth{dashboard: tc.d}, size[0], size[1])
			view := xansi.Strip(screen(m))
			popoverBounds(t, view, size[0], size[1])
			requireContains(t, view, tc.want...)
			if name == "not set up" && strings.Contains(view, "Re-run") {
				t.Errorf("not set up must not offer a re-run line:\n%s", view)
			}
			if got := strings.Count(m.renderHealthPopover(m.terminalWidth()), "\n") + 1; got != healthPopoverHeight {
				t.Errorf("%s %v: popover is %d lines, want %d", name, size, got, healthPopoverHeight)
			}
		}
	}
}

func TestHealthPopoverMarksRowsMeasuredOnOtherCode(t *testing.T) {
	for name, tc := range map[string]struct {
		state codehealth.CodeState
		marks bool
	}{
		"matches":      {codehealth.CodeMatches, false},
		"unknown":      {codehealth.CodeUnknown, false},
		"moved on":     {codehealth.CodeMovedOn, true},
		"other branch": {codehealth.CodeOtherBranch, true},
	} {
		t.Run(name, func(t *testing.T) {
			text := "Freshness sentence for " + name
			f := &fakeHealth{dashboard: measuredDashboard(fiveSignals()), fresh: codehealth.CodeFreshness{State: tc.state, Text: text}}
			for _, size := range popoverSizes {
				view := xansi.Strip(screen(openHealthScreenAt(t, f, size[0], size[1])))
				popoverBounds(t, view, size[0], size[1])
				requireContains(t, view, text)
				if got := strings.Count(view, "42*") == 5; got != tc.marks {
					t.Errorf("%v: marked rows = %v, want %v:\n%s", size, got, tc.marks, view)
				}
				if got := strings.Contains(view, "measured on other code"); got != tc.marks {
					t.Errorf("%v: selected note = %v, want %v", size, got, tc.marks)
				}
			}
		})
	}
}

func TestHealthPopoverKeepsUnknownAndPartialMarks(t *testing.T) {
	rows := fiveSignals()
	rows[3].Meaning = "Duplication could not be measured."
	view := xansi.Strip(screen(openHealthScreen(t, &fakeHealth{dashboard: measuredDashboard(rows)})))
	requireContains(t, view, "? Duplication", "~ Coverage")
	m := openHealthScreen(t, &fakeHealth{dashboard: measuredDashboard(rows)})
	m = press(t, press(t, press(t, m, "down"), "down"), "down")
	requireContains(t, xansi.Strip(screen(m)), "Duplication could not be measured.")
}

func TestHealthPopoverRefreshAndNoticesFitEverySize(t *testing.T) {
	for _, size := range popoverSizes {
		f := &fakeHealth{dashboard: measuredDashboard(fiveSignals()), refresh: func(ctx context.Context, _ func(codehealth.Progress)) error {
			<-ctx.Done()
			return codehealth.ErrCollectionCancelled
		}}
		m := openHealthScreenAt(t, f, size[0], size[1])
		m, _ = startRefresh(t, m)
		m = feed(t, m, healthProgressMsg{Progress: codehealth.Progress{Position: 2, Total: 5, Capability: codehealth.CapabilityCoverage, Provider: "go_cover"}})
		view := xansi.Strip(screen(m))
		popoverBounds(t, view, size[0], size[1])
		requireContains(t, view, "Refreshing 2 of 5: Coverage", "Esc to cancel")

		m = feed(t, m, healthRefreshDoneMsg{Err: codehealth.ErrCollectionCancelled})
		view = xansi.Strip(screen(m))
		popoverBounds(t, view, size[0], size[1])
		requireContains(t, view, textHealthCancelled, "Tests")

		m = feed(t, m, healthLoadedMsg{Err: context.DeadlineExceeded})
		view = xansi.Strip(screen(m))
		popoverBounds(t, view, size[0], size[1])
		requireContains(t, view, "Could not read history", "Tests")
	}
}
