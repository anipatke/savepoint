package v2

import (
	"fmt"
	"strings"
	"testing"

	xansi "github.com/charmbracelet/x/ansi"
	"github.com/opencode/savepoint/internal/codehealth"
)

func historyDashboard(n int) codehealth.Dashboard {
	d := measuredDashboard(fiveSignals())
	d.History = nil
	for i := 0; i < n; i++ {
		e := codehealth.DashboardHistoryEntry{
			WhenText: fmt.Sprintf("%d Oct 09:00 UTC", 30-i), Origin: codehealth.OriginOfficial, OriginText: "Official check",
			Overall: codehealth.ClassificationGood, OverallText: "Good",
		}
		if i == 1 {
			e.Origin, e.OriginText = codehealth.OriginManual, "Manual refresh"
		}
		d.History = append(d.History, e)
	}
	return d
}

func TestHealthHistoryShowsAtMostTenNewestFirst(t *testing.T) {
	for _, size := range popoverSizes {
		m := press(t, openHealthScreenAt(t, &fakeHealth{dashboard: historyDashboard(14)}, size[0], size[1]), "h")
		view := xansi.Strip(screen(m))
		popoverBounds(t, view, size[0], size[1])
		requireContains(t, view, "last checks, newest first", "30 Oct 09:00 UTC  Official check  ✓ Good", "21 Oct 09:00 UTC", "Manual refresh", "do not feed the sparklines", "Re-run: savepoint health check")
		if strings.Contains(view, "20 Oct") {
			t.Errorf("%v: eleventh check shown:\n%s", size, view)
		}
		if strings.Index(view, "30 Oct") > strings.Index(view, "29 Oct") {
			t.Errorf("%v: not newest first", size)
		}
		if got := strings.Count(m.renderHealthPopover(m.terminalWidth()), "\n") + 1; got != healthPopoverHeight {
			t.Errorf("%v: popover is %d lines, want %d", size, got, healthPopoverHeight)
		}
	}
}

func TestHealthHistoryOneAndNoSnapshots(t *testing.T) {
	m := press(t, openHealthScreen(t, &fakeHealth{dashboard: historyDashboard(1)}), "h")
	requireContains(t, xansi.Strip(screen(m)), "30 Oct 09:00 UTC  Official check  ✓ Good")
	m = press(t, openHealthScreen(t, &fakeHealth{dashboard: historyDashboard(0)}), "h")
	view := xansi.Strip(screen(m))
	requireContains(t, view, "No checks yet.")
	popoverBounds(t, view, 120, 40)
}

func TestHealthHistoryRendersManualRefreshesDimmed(t *testing.T) {
	m := press(t, openHealthScreen(t, &fakeHealth{dashboard: historyDashboard(3)}), "h")
	view := xansi.Strip(screen(m))
	requireContains(t, view, "29 Oct 09:00 UTC  Manual refresh  ✓ Good", "30 Oct 09:00 UTC  Official check  ✓ Good")
	d := historyDashboard(3)
	if d.History[1].Origin != codehealth.OriginManual {
		t.Fatalf("fixture lost its manual entry")
	}
}

func TestHealthHistoryReturnPaths(t *testing.T) {
	open := func() Model {
		return openHealthScreen(t, &fakeHealth{dashboard: historyDashboard(3)})
	}
	for _, back := range []string{"h", "esc"} {
		m := press(t, open(), "h")
		if !m.Health.History {
			t.Fatalf("h did not open history")
		}
		m = press(t, m, back)
		if m.Health == nil || m.Health.History {
			t.Fatalf("%q did not return to the signals", back)
		}
		requireContains(t, xansi.Strip(screen(m)), "✓ Tests")
		if m = press(t, m, "esc"); m.Health != nil {
			t.Errorf("second esc did not close the popover")
		}
	}
	m := press(t, open(), "h", "down", "j")
	if m.Health.Cursor != 0 {
		t.Errorf("signal cursor moved inside history")
	}
}

func TestHealthHistoryNamedInHelpAndHints(t *testing.T) {
	m := openHealthScreen(t, &fakeHealth{dashboard: historyDashboard(3)})
	requireContains(t, m.hints(), "h:history")
	requireContains(t, renderHelp(m, 100, 60), "show the last ten checks")
	requireContains(t, press(t, m, "h").hints(), "h/esc:signals")
}
