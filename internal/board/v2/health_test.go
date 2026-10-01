package v2

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	xansi "github.com/charmbracelet/x/ansi"
	"github.com/opencode/savepoint/internal/codehealth"
	"github.com/opencode/savepoint/internal/styles"
	"github.com/opencode/savepoint/internal/testutil"
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
		Question: text + "?", Value: "42", Aim: "aim: 80 or more", Meaning: reason, NextStep: "Do the next thing.",
		Spark: "▁▃▅▇", SparkWord: "better", SignOff: "Advisory only", Where: "internal/foo/foo.go",
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
		MeasuredAt: "2026-10-02T09:00:00Z", MeasuredText: "2 Oct 2026", Headline: "2 signals need a look.", SignOff: "Doesn't block sign-off", Rows: rows,
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

func TestHealthLoadErrorKeepsThePreviousView(t *testing.T) {
	f := &fakeHealth{dashboard: measuredDashboard(fiveSignals())}
	m := openHealthScreen(t, f)
	f.loadErr = errors.New("history is damaged")
	m = settle(t, m, healthLoadCmd(f.funcs(), m.Root))
	got := screen(m)

	requireContains(t, got, "history is damaged", "Tests", "Coverage")
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

	for _, want := range []string{"Refreshing 1 of 2: Tests", "Refreshing 2 of 2: Coverage"} {
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
	requireContains(t, got, "Refresh cancelled; nothing was saved.", "Coverage")
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

	requireContains(t, screen(m), "Refresh failed; nothing was saved: disk is full", "Coverage")
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

func TestHealthCtrlCFromHelpDuringRefreshCancelsBeforeQuitting(t *testing.T) {
	cancelled := make(chan struct{})
	f := &fakeHealth{dashboard: measuredDashboard(fiveSignals())}
	f.refresh = func(ctx context.Context, _ func(codehealth.Progress)) error {
		<-ctx.Done()
		close(cancelled)
		return codehealth.ErrCollectionCancelled
	}
	m := openHealthScreen(t, f)
	m, _ = startRefresh(t, m)
	m = press(t, m, "?")
	if !m.Help {
		t.Fatal("? did not open Help during the refresh")
	}
	_, cmd := sendKey(t, m, "ctrl+c")

	if cmd == nil {
		t.Fatalf("ctrl+c from Help did not quit")
	}
	select {
	case <-cancelled:
	case <-time.After(2 * time.Second):
		t.Fatal("quitting from Help did not cancel the refresh")
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

// headerLine is the header's content row, the one carrying the title.
func headerLine(m Model) string {
	for _, line := range strings.Split(screen(m), "\n") {
		if strings.Contains(line, "S A V E P O I N T") {
			return line
		}
	}
	return ""
}

func withChip(t *testing.T, width int, chip codehealth.Chip) Model {
	t.Helper()
	m := openSizedBoard(t, writeValidProject(t), width, 24)
	m.State.HealthChip = chip
	return m
}

var measuredChip = codehealth.Chip{State: codehealth.ChipMeasured, Overall: codehealth.ClassificationWatch, Label: "Watch", Good: 3, Signals: 5}

func TestHeaderChipStatesAndWording(t *testing.T) {
	tests := []struct {
		name string
		chip codehealth.Chip
		want string
	}{
		{"measured", measuredChip, "♥ Health 3/5"},
		{"no check", codehealth.Chip{State: codehealth.ChipNoCheck}, "♥ Health: no check yet"},
		{"not set up", codehealth.Chip{State: codehealth.ChipNotSetUp}, "♥ Health: not set up"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := headerLine(withChip(t, 120, tc.chip))
			requireContains(t, got, tc.want, "Objectives", "Tasks", "Issues")
		})
	}
}

func TestHeaderChipColourFollowsOverallLabel(t *testing.T) {
	bad := codehealth.ClassificationNeedsAttention
	tests := []struct {
		chip codehealth.Chip
		want lipgloss.TerminalColor
	}{
		{codehealth.Chip{State: codehealth.ChipMeasured, Overall: codehealth.ClassificationGood}, styles.HealthGood.GetForeground()},
		{codehealth.Chip{State: codehealth.ChipMeasured, Overall: codehealth.ClassificationWatch}, styles.HealthWatch.GetForeground()},
		{codehealth.Chip{State: codehealth.ChipMeasured, Overall: bad}, styles.HealthNeedsAttention.GetForeground()},
		{codehealth.Chip{State: codehealth.ChipMeasured, Overall: codehealth.ClassificationUnknown}, styles.HealthUnknown.GetForeground()},
		{codehealth.Chip{State: codehealth.ChipNoCheck}, styles.HealthUnknown.GetForeground()},
		{codehealth.Chip{State: codehealth.ChipNotSetUp}, styles.HealthUnknown.GetForeground()},
	}
	for _, tc := range tests {
		if got := chipStyle(tc.chip).GetForeground(); got != tc.want {
			t.Errorf("chipStyle(%+v) foreground = %v, want %v", tc.chip, got, tc.want)
		}
	}
	if styles.HealthGood.GetForeground() == styles.HealthWatch.GetForeground() || styles.HealthWatch.GetForeground() == styles.HealthNeedsAttention.GetForeground() {
		t.Fatal("overall colours are not distinct")
	}
}

func TestHeaderChipStaysOneRowAt80AndDegradesBeforeCounts(t *testing.T) {
	for _, width := range []int{80, 100, 70, 60, 50} {
		m := withChip(t, width, measuredChip)
		inner := m.terminalWidth()
		header := m.renderHeader(inner)
		bare := withChip(t, width, codehealth.Chip{})
		if got, want := lipgloss.Height(header), lipgloss.Height(bare.renderHeader(inner)); got != want {
			t.Fatalf("width %d: header is %d rows with the chip, %d without:\n%s", width, got, want, header)
		}
		line := xansi.Strip(header)
		if lipgloss.Width(line) > inner {
			t.Errorf("width %d: header is %d cells, terminal content is %d", width, lipgloss.Width(line), inner)
		}
		if strings.Contains(line, "♥") && !strings.Contains(line, "Issues") {
			t.Errorf("width %d: chip kept while a count was dropped: %q", width, line)
		}
	}
	wide := xansi.Strip(withChip(t, 100, measuredChip).renderHeader(96))
	requireContains(t, wide, "♥ Health 3/5")

	// Words go first, then the chip; the counts keep their text throughout.
	var sawShort, sawGone bool
	for width := 100; width >= 50; width-- {
		m := withChip(t, width, measuredChip)
		line := xansi.Strip(m.renderHeader(m.terminalWidth()))
		switch {
		case strings.Contains(line, "♥ Health 3/5"):
			if sawShort || sawGone {
				t.Fatalf("width %d: full chip returned after it shortened", width)
			}
		case strings.Contains(line, "♥ 3/5"):
			sawShort = true
			if sawGone {
				t.Fatalf("width %d: short chip returned after it vanished", width)
			}
		default:
			sawGone = true
		}
	}
	if !sawShort || !sawGone {
		t.Errorf("chip never passed through short form and disappearance (short=%v gone=%v)", sawShort, sawGone)
	}
}

func TestHeaderChipUnreadableStorageIsNotABoardError(t *testing.T) {
	root := writeValidProject(t)
	testutil.WriteFile(t, filepath.Join(root, "health", "snapshots", "junk.json"), "{")
	testutil.WriteFile(t, filepath.Join(root, "health", "config.json"), "not json")
	msg := loadProject(root)
	if msg.Failed() {
		t.Fatalf("loadProject failed on damaged health storage: %s", msg.Diagnostic)
	}
	if got := msg.State.HealthChip.State; got != codehealth.ChipNotSetUp {
		t.Errorf("chip state = %v, want not set up", got)
	}
}

func TestHeaderChipRefreshesAfterRefreshResult(t *testing.T) {
	f := &fakeHealth{dashboard: codehealth.Dashboard{State: codehealth.DashboardFirstRun}}
	m := openHealthScreen(t, f)
	m.State.HealthChip = codehealth.Chip{State: codehealth.ChipNoCheck}

	measured := measuredDashboard(fiveSignals())
	f.mu.Lock()
	f.dashboard = measured
	f.mu.Unlock()
	f.refresh = func(context.Context, func(codehealth.Progress)) error { return nil }

	m, wait := startRefresh(t, m)
	for done := false; !done; {
		msg := wait()
		next, cmd := m.Update(msg)
		m = next.(Model)
		if _, done = msg.(healthRefreshDoneMsg); done {
			m = settle(t, m, cmd)
		}
	}
	if got, want := m.State.HealthChip, measured.Chip(); got != want {
		t.Fatalf("chip after refresh = %+v, want %+v", got, want)
	}
	if got := m.State.HealthChip.State; got != codehealth.ChipMeasured {
		t.Fatalf("chip state = %v, want measured", got)
	}
}
