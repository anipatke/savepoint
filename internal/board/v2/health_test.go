package v2

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/opencode/savepoint/internal/codehealth"
)

// fakeHealth stands in for codehealth: every call is counted and every result
// is set by the test, so the screen is driven by messages alone.
type fakeHealth struct {
	mu        sync.Mutex
	dashboard codehealth.Dashboard
	loadErr   error
	fresh     codehealth.CodeFreshness
	loads     int
	saved     int
	// refresh runs inside the refresh goroutine.
	refresh func(ctx context.Context, progress func(codehealth.Progress)) error
}

func (f *fakeHealth) funcs() HealthFuncs {
	return HealthFuncs{
		Load: func(string) (codehealth.Dashboard, error) {
			f.mu.Lock()
			defer f.mu.Unlock()
			f.loads++
			return f.dashboard, f.loadErr
		},
		Freshness: func(context.Context, string, codehealth.RepositoryIdentity) codehealth.CodeFreshness {
			return f.fresh
		},
		Refresh: func(ctx context.Context, _ string, progress func(codehealth.Progress)) error {
			return f.refresh(ctx, progress)
		},
	}
}

func healthRow(c codehealth.Capability, text string, label codehealth.Classification, labelText, reason, trend string) codehealth.DashboardRow {
	return codehealth.DashboardRow{
		Capability: c, CapabilityText: text, Label: label, LabelText: labelText,
		Explanation: reason, Trend: trend, Basis: "Compared with 3 earlier comparable official checks.",
		Outcome: codehealth.OutcomeAvailable, OutcomeText: "measured",
		Required: true, RequiredText: "required",
		Provider: "go_test_json", ProviderVersion: "1.2", Scope: []string{"./..."},
		CollectedAt: "2026-10-02T09:00:00Z",
		Evidence:    []codehealth.EvidenceRef{{Path: "internal/foo/foo.go", Line: 12, Note: "complex function"}},
	}
}

func fiveSignals() []codehealth.DashboardRow {
	good, watch, bad, unknown := codehealth.ClassificationGood, codehealth.ClassificationWatch, codehealth.ClassificationNeedsAttention, codehealth.ClassificationUnknown
	return []codehealth.DashboardRow{
		healthRow(codehealth.CapabilityTests, "Tests", good, "Good", "All 120 tests pass.", "Steady over 4 official checks"),
		healthRow(codehealth.CapabilityCoverage, "Coverage", watch, "Watch", "Coverage is 61%, down from 68%.", "Declining over 4 official checks"),
		healthRow(codehealth.CapabilityComplexity, "Complexity", bad, "Needs Attention", "One function is far too complex.", "No trend yet"),
		healthRow(codehealth.CapabilityDuplication, "Duplication", unknown, "Unknown", "The tool could not be used.", "No trend yet"),
		healthRow(codehealth.CapabilityDependencyVulnerability, "Dependency vulnerabilities", good, "Good", "No known vulnerabilities.", "Steady over 4 official checks"),
	}
}

func measuredDashboard(rows []codehealth.DashboardRow) codehealth.Dashboard {
	return codehealth.Dashboard{
		State: codehealth.DashboardMeasured, Overall: codehealth.ClassificationWatch, OverallText: "Watch",
		Origin: codehealth.OriginOfficial, OriginText: "Official check", SnapshotID: "snap-1",
		MeasuredAt: "2026-10-02T09:00:00Z", Rows: rows,
		History: []codehealth.DashboardHistoryEntry{
			{CreatedAt: "2026-10-02T09:00:00Z", OriginText: "Official check", OverallText: "Watch"},
			{CreatedAt: "2026-10-01T09:00:00Z", OriginText: "Manual refresh", OverallText: "Good"},
		},
	}
}

// settle runs a command and feeds its message to the model, then does the same
// for the follow-up command, until none is left. Only for commands that finish
// on their own.
func settle(t *testing.T, m Model, cmd tea.Cmd) Model {
	t.Helper()
	for cmd != nil {
		msg := cmd()
		if msg == nil {
			return m
		}
		if _, batch := msg.(tea.BatchMsg); batch {
			t.Fatalf("settle reached a batch; start refreshes with startRefresh")
		}
		next, c := m.Update(msg)
		m, cmd = next.(Model), c
	}
	return m
}

func healthBoard(t *testing.T, f *fakeHealth, width, height int) Model {
	t.Helper()
	m := openSizedBoard(t, writeValidProject(t), width, height)
	m.HealthFuncs = f.funcs()
	return m
}

func openHealthScreen(t *testing.T, f *fakeHealth) Model {
	t.Helper()
	return openHealthScreenAt(t, f, 120, 40)
}

func openHealthScreenAt(t *testing.T, f *fakeHealth, width, height int) Model {
	t.Helper()
	m := healthBoard(t, f, width, height)
	next, cmd := m.Update(keyMsg("H"))
	return settle(t, next.(Model), cmd)
}

func sendKey(t *testing.T, m Model, key string) (Model, tea.Cmd) {
	t.Helper()
	next, cmd := m.Update(keyMsg(key))
	return next.(Model), cmd
}

// startRefresh presses R and runs the refresh goroutine, returning the model
// and a function that waits for the next message the refresh sent.
func startRefresh(t *testing.T, m Model) (Model, func() tea.Msg) {
	t.Helper()
	m, cmd := sendKey(t, m, "R")
	if cmd == nil || m.Health.Refresh == nil {
		t.Fatalf("R did not start a refresh")
	}
	parts, ok := cmd().(tea.BatchMsg)
	if !ok || len(parts) != 2 {
		t.Fatalf("R returned %T, want the refresh and its message reader", cmd())
	}
	go parts[0]()
	events := m.Health.Refresh.Events
	return m, func() tea.Msg {
		select {
		case msg := <-events:
			return msg
		case <-time.After(2 * time.Second):
			t.Fatal("no refresh message arrived")
			return nil
		}
	}
}

func feed(t *testing.T, m Model, msg tea.Msg) Model {
	t.Helper()
	next, _ := m.Update(msg)
	return next.(Model)
}

func TestHealthNotConfiguredExplainsAndRefreshDoesNothing(t *testing.T) {
	f := &fakeHealth{dashboard: codehealth.Dashboard{State: codehealth.DashboardNotConfigured}}
	m := openHealthScreen(t, f)

	requireContains(t, screen(m), "CODE HEALTH", "not set up", "savepoint health setup")
	m, cmd := sendKey(t, m, "R")
	if cmd != nil || m.Health.Refresh != nil {
		t.Fatalf("R started a refresh with health not configured")
	}
	if strings.Contains(screen(m), "R:refresh") {
		t.Errorf("footer offers R when it does nothing:\n%s", screen(m))
	}
}

func TestHealthFirstRunOffersRefresh(t *testing.T) {
	f := &fakeHealth{dashboard: codehealth.Dashboard{State: codehealth.DashboardFirstRun}}
	m := openHealthScreen(t, f)

	requireContains(t, screen(m), "Nothing has been measured yet", "Press R", "R:refresh")
}

func TestHealthShowsFiveSignalsWithLabelsReasonsAndTrends(t *testing.T) {
	f := &fakeHealth{dashboard: measuredDashboard(fiveSignals()), fresh: codehealth.CodeFreshness{State: codehealth.CodeMatches, Text: "Your code matches what was measured."}}
	m := openHealthScreen(t, f)
	got := screen(m)

	requireContains(t, got,
		"Watch  overall", "Official check", "Your code matches what was measured.",
		"Tests", "Coverage", "Complexity", "Duplication", "Dependency vulnerabilities",
		"Coverage is 61%, down from 68%.", "Declining over 4 official checks",
		"Needs Attention", "Unknown", "RECENT CHECKS", "Manual refresh",
	)
	if strings.Contains(got, "Overall score") || strings.Contains(got, "/100") {
		t.Errorf("screen shows an overall number:\n%s", got)
	}
}

func TestHealthNonGoodRowsAreNeverShownAsGood(t *testing.T) {
	rows := fiveSignals()
	rows[0] = healthRow(codehealth.CapabilityTests, "Tests", codehealth.ClassificationUnknown, "Unknown", "Only part of the tests ran.", "No trend yet")
	rows[0].Outcome, rows[0].OutcomeText = codehealth.OutcomePartial, "partial"
	rows[1].Label, rows[1].LabelText, rows[1].Explanation = codehealth.ClassificationNeedsAttention, "Needs Attention", "A required tool failed."
	rows[1].Outcome, rows[1].OutcomeText = codehealth.OutcomeFailed, "failed"
	rows[3].Required, rows[3].RequiredText, rows[3].Outcome, rows[3].OutcomeText = false, "optional", codehealth.OutcomeUnavailable, "unavailable: the tool could not be used"
	rows[4] = codehealth.DashboardRow{
		Capability: codehealth.CapabilityDependencyVulnerability, CapabilityText: "Dependency vulnerabilities", NotConfigured: true,
		Label: codehealth.ClassificationUnknown, LabelText: "Unknown", Explanation: "Not configured, so health for this signal is unknown.",
		Outcome: codehealth.OutcomeNotConfigured, OutcomeText: "not configured", RequiredText: "optional",
	}
	d := measuredDashboard(rows)
	d.Overall, d.OverallText = codehealth.ClassificationUnknown, "Unknown"
	d.History[1].OverallText = "Watch"
	m := openHealthScreen(t, &fakeHealth{dashboard: d})
	got := screen(m)

	if strings.Contains(got, "Good") {
		t.Errorf("a partial, failed, unavailable, or not-configured row reads as Good:\n%s", got)
	}
	requireContains(t, got, "Only part of the tests ran.", "A required tool failed.", "Not configured, so health")

	m = press(t, m, "down", "down", "down", "enter")
	requireContains(t, screen(m), "Outcome: unavailable", "Required: optional")
}

func TestHealthFreshnessLineSaysWhyResultsMayBeOld(t *testing.T) {
	for name, text := range map[string]string{
		"stale":    "Your code has moved on by 3 commits since this was measured.",
		"dirty":    "Files you are working on differ from what was measured.",
		"diverged": "This was measured on a different branch.",
		"unknown":  "Git could not be read, so freshness is unknown.",
	} {
		t.Run(name, func(t *testing.T) {
			f := &fakeHealth{dashboard: measuredDashboard(fiveSignals()), fresh: codehealth.CodeFreshness{State: codehealth.CodeUnknown, Text: text}}
			requireContains(t, screen(openHealthScreen(t, f)), text)
		})
	}
}

func TestHealthShowsResultsBeforeFreshnessArrives(t *testing.T) {
	f := &fakeHealth{dashboard: measuredDashboard(fiveSignals()), fresh: codehealth.CodeFreshness{Text: "Your code matches what was measured."}}
	m := healthBoard(t, f, 120, 40)
	next, cmd := m.Update(keyMsg("H"))
	next, freshness := next.(Model).Update(cmd())
	got := screen(next.(Model))

	requireContains(t, got, "Tests", "Checking whether your code has changed")
	if freshness == nil {
		t.Fatalf("loading a measured dashboard did not start the freshness command")
	}
	requireContains(t, screen(settle(t, next.(Model), freshness)), "Your code matches what was measured.")
}

func TestHealthStaleFreshnessForAnOlderSnapshotIsIgnored(t *testing.T) {
	f := &fakeHealth{dashboard: measuredDashboard(fiveSignals())}
	m := openHealthScreen(t, f)
	m = feed(t, m, healthFreshnessMsg{SnapshotID: "older", Freshness: codehealth.CodeFreshness{Text: "ANSWER FOR ANOTHER SNAPSHOT"}})
	if strings.Contains(screen(m), "ANSWER FOR ANOTHER SNAPSHOT") {
		t.Errorf("freshness for a different snapshot was applied")
	}
}

func TestHealthComparableAndResetTrendsAreShownAsWorded(t *testing.T) {
	rows := fiveSignals()
	rows[1].Trend, rows[1].Basis = "Improving over 4 official checks, 55% to 68%", "Compared with 3 earlier comparable official checks."
	rows[2].Trend, rows[2].Basis = "No trend yet", "No comparable official history yet. 2 earlier official results not compared."
	m := openHealthScreen(t, &fakeHealth{dashboard: measuredDashboard(rows)})

	requireContains(t, screen(m), "Improving over 4 official checks, 55% to 68%")
	m = press(t, m, "down", "down", "enter")
	requireContains(t, screen(m), "No comparable official history yet. 2 earlier official results not compared.")
}

func TestHealthDetailsExplainTheLabelAndEscReturnsToTheList(t *testing.T) {
	m := openHealthScreen(t, &fakeHealth{dashboard: measuredDashboard(fiveSignals())})
	m = press(t, m, "down", "enter")
	got := screen(m)

	requireContains(t, got,
		"SIGNAL DETAIL", "Coverage", "Coverage is 61%, down from 68%.",
		"Trend: Declining over 4 official checks", "Compared with: Compared with 3 earlier",
		"Outcome: measured", "Required: required", "Provider: go_test_json 1.2", "Scope: ./...",
		"Measured: 2026-10-02T09:00:00Z", "AFFECTED AREAS", "internal/foo/foo.go:12 — complex function",
	)
	m = press(t, m, "esc")
	if m.Health == nil || m.Health.Detail {
		t.Fatalf("esc in details did not return to the list")
	}
	requireContains(t, screen(m), "RECENT CHECKS")
	m = press(t, m, "v")
	requireContains(t, screen(m), "SIGNAL DETAIL")
}

func TestHealthLoadErrorKeepsThePreviousView(t *testing.T) {
	f := &fakeHealth{dashboard: measuredDashboard(fiveSignals())}
	m := openHealthScreen(t, f)
	f.loadErr = errors.New("history is damaged")
	m = settle(t, m, healthLoadCmd(f.funcs(), m.Root))
	got := screen(m)

	requireContains(t, got, "history is damaged", "Coverage is 61%")
}

func TestHealthLoadErrorOnFirstOpenIsAStatusLineNotACrash(t *testing.T) {
	m := openHealthScreen(t, &fakeHealth{loadErr: errors.New("history is damaged")})
	requireContains(t, screen(m), "CODE HEALTH", "history is damaged")
}

func TestHealthRefreshReportsProgressInOrderAndReloadsOnCompletion(t *testing.T) {
	release := make(chan struct{})
	f := &fakeHealth{dashboard: measuredDashboard(fiveSignals())}
	f.refresh = func(ctx context.Context, progress func(codehealth.Progress)) error {
		progress(codehealth.Progress{Position: 1, Total: 2, Capability: codehealth.CapabilityTests, Provider: "go_test_json"})
		progress(codehealth.Progress{Position: 2, Total: 2, Capability: codehealth.CapabilityCoverage, Provider: "go_cover"})
		<-release
		return nil
	}
	m := openHealthScreen(t, f)
	loadsBefore := f.loads
	m, next := startRefresh(t, m)

	for _, want := range []string{"Refreshing 1 of 2: Tests (go_test_json)", "Refreshing 2 of 2: Coverage (go_cover)"} {
		updated, cmd := m.Update(next())
		m = updated.(Model)
		requireContains(t, screen(m), want, "Esc to cancel")
		if cmd == nil {
			t.Fatalf("progress did not keep reading the refresh")
		}
	}

	// A second R while running is ignored.
	again, cmd := sendKey(t, m, "R")
	if cmd != nil || again.Health.Refresh != m.Health.Refresh {
		t.Fatalf("R during a refresh started another")
	}

	close(release)
	done := next()
	if _, ok := done.(healthRefreshDoneMsg); !ok {
		t.Fatalf("last message = %T, want the result", done)
	}
	updated, cmd := m.Update(done)
	m = updated.(Model)
	if m.Health.Refresh != nil {
		t.Fatalf("a finished refresh is still shown as running")
	}
	m = settle(t, m, cmd)
	if f.loads <= loadsBefore {
		t.Errorf("completion did not reload the dashboard")
	}
}

func TestHealthEscCancelsARefreshAndKeepsThePreviousResult(t *testing.T) {
	f := &fakeHealth{dashboard: measuredDashboard(fiveSignals())}
	f.refresh = func(ctx context.Context, progress func(codehealth.Progress)) error {
		progress(codehealth.Progress{Position: 1, Total: 5, Capability: codehealth.CapabilityTests, Provider: "go_test_json"})
		<-ctx.Done()
		return codehealth.ErrCollectionCancelled
	}
	m := openHealthScreen(t, f)
	loadsBefore := f.loads
	m, next := startRefresh(t, m)
	m = feed(t, m, next())

	m, cmd := sendKey(t, m, "esc")
	if cmd != nil || m.Health == nil || m.Health.Refresh == nil {
		t.Fatalf("esc during a refresh must cancel it and stay on the screen")
	}
	requireContains(t, screen(m), "Cancelling")

	m = feed(t, m, next())
	got := screen(m)
	requireContains(t, got, "Refresh cancelled; nothing was saved.", "Coverage is 61%")
	if m.Health.Refresh != nil || f.loads != loadsBefore || f.saved != 0 {
		t.Errorf("a cancelled refresh left state behind: running=%v loads=%d saved=%d", m.Health.Refresh != nil, f.loads, f.saved)
	}

	m, _ = sendKey(t, m, "esc")
	if m.Health != nil {
		t.Errorf("esc with nothing running did not close the screen")
	}
}

func TestHealthRefreshErrorIsAStatusLine(t *testing.T) {
	f := &fakeHealth{dashboard: measuredDashboard(fiveSignals())}
	f.refresh = func(context.Context, func(codehealth.Progress)) error { return errors.New("disk is full") }
	m := openHealthScreen(t, f)
	m, next := startRefresh(t, m)
	m = feed(t, m, next())

	requireContains(t, screen(m), "Refresh failed; nothing was saved: disk is full", "Coverage is 61%")
}

func TestHealthQuitDuringRefreshCancelsBeforeQuitting(t *testing.T) {
	cancelled := make(chan struct{})
	f := &fakeHealth{dashboard: measuredDashboard(fiveSignals())}
	f.refresh = func(ctx context.Context, _ func(codehealth.Progress)) error {
		<-ctx.Done()
		close(cancelled)
		return codehealth.ErrCollectionCancelled
	}
	m := openHealthScreen(t, f)
	m, _ = startRefresh(t, m)
	_, cmd := sendKey(t, m, "ctrl+c")

	if cmd == nil {
		t.Fatalf("ctrl+c did not quit")
	}
	select {
	case <-cancelled:
	case <-time.After(2 * time.Second):
		t.Fatal("quitting did not cancel the refresh")
	}
}

func TestHealthEscRestoresTheBoardCursor(t *testing.T) {
	f := &fakeHealth{dashboard: measuredDashboard(fiveSignals())}
	m := healthBoard(t, f, 120, 40)
	m = press(t, m, "right", "down")
	column, card := m.FocusedColumn, m.FocusedCard
	next, cmd := m.Update(keyMsg("H"))
	m = settle(t, next.(Model), cmd)
	m = press(t, m, "down", "esc")

	if m.Health != nil || m.FocusedColumn != column || m.FocusedCard != card {
		t.Errorf("closing Code Health left the cursor at %v/%d, want %v/%d", m.FocusedColumn, m.FocusedCard, column, card)
	}
}

func TestHealthIsExclusiveWithTheOtherOverlays(t *testing.T) {
	f := &fakeHealth{dashboard: measuredDashboard(fiveSignals())}
	m := healthBoard(t, f, 120, 40)
	for _, open := range []string{"i", "?"} {
		opened := press(t, m, open)
		next, cmd := opened.Update(keyMsg("H"))
		if cmd != nil || next.(Model).Health != nil {
			t.Errorf("H opened Code Health over the %q overlay", open)
		}
	}
	m = openHealthScreen(t, f)
	if m = press(t, m, "i"); m.Issues != nil {
		t.Errorf("i opened Issues over Code Health")
	}
}

func TestHealthHelpListsTheKeys(t *testing.T) {
	boardHelp := screen(press(t, healthBoard(t, &fakeHealth{}, 120, 50), "?"))
	requireContains(t, boardHelp, "H:", "open Code Health")

	m := openHealthScreen(t, &fakeHealth{dashboard: measuredDashboard(fiveSignals())})
	m = press(t, m, "?")
	requireContains(t, screen(m), "in Code Health: refresh now", "esc cancels it")
	if m = press(t, m, "esc"); m.Health == nil || m.Help {
		t.Errorf("closing help did not return to Code Health")
	}
}

func TestHealthFitsTheBoardsMinimumWidth(t *testing.T) {
	rows := fiveSignals()
	rows[2].Explanation = strings.Repeat("A very long explanation that must wrap rather than widen the screen. ", 3)
	rows[2].Evidence = []codehealth.EvidenceRef{{Path: "internal/a/very/long/path/that/keeps/going/and/going/file.go", Line: 99, Note: "note"}}
	f := &fakeHealth{dashboard: measuredDashboard(rows), fresh: codehealth.CodeFreshness{Text: strings.Repeat("Your code has moved on. ", 4)}}

	for _, width := range []int{compactBoardBreakpoint - 1, compactBoardBreakpoint, 80} {
		m := openHealthScreenAt(t, f, width, 24)
		for step, view := range []string{screen(m), screen(press(t, m, "down", "down", "enter"))} {
			for _, line := range strings.Split(view, "\n") {
				if lipgloss.Width(line) > width {
					t.Errorf("width %d, view %d: line is %d cells: %q", width, step, lipgloss.Width(line), line)
				}
			}
			if got := strings.Count(view, "\n") + 1; got > 24 {
				t.Errorf("width %d, view %d: %d lines overflow a 24-line terminal", width, step, got)
			}
		}
	}
}

func TestHealthKeepsTheCursorRowOnScreenInAShortTerminal(t *testing.T) {
	m := openHealthScreenAt(t, &fakeHealth{dashboard: measuredDashboard(fiveSignals())}, 100, 20)
	m = press(t, m, "down", "down", "down", "down")

	requireContains(t, screen(m), "▸")
	requireContains(t, screen(m), "Dependency vulnerabilities")
}
