package v2

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	xansi "github.com/charmbracelet/x/ansi"
	"github.com/muesli/termenv"
	"github.com/opencode/savepoint/internal/data"
)

// fixtureCard builds a card from values a resolver would have returned, with no
// project behind it. Every rendering test below uses it, which is what proves
// the renderer reads the resolved values and nothing else: there is no index,
// no Check, and no dependency here to inspect.
func fixtureCard(task *data.TaskV2, clearance data.Clearance, decision data.GateDecision) TaskCard {
	return TaskCard{
		Task:        task,
		Clearance:   clearance,
		Decision:    decision,
		ByException: decision.AllowedByException,
	}
}

func fixtureTask(id, title string, status data.ColumnType, stage data.ProgressStage) *data.TaskV2 {
	return &data.TaskV2{ID: id, Title: title, Objective: "O001", Status: status, Stage: stage}
}

func renderedText(card TaskCard, width int, focused bool) string {
	return xansi.Strip(renderCard(card, width, focused))
}

func TestRenderCardLabelsWithTheTitleAndCarriesTheIdentity(t *testing.T) {
	card := fixtureCard(
		fixtureTask("T001", "Open the board a V2 project already has", data.ColumnPlanned, ""),
		data.Clearance{State: data.ClearanceMissing},
		data.GateDecision{Allowed: true},
	)

	got := renderedText(card, 40, false)

	if !strings.Contains(got, "Open the board a V2 project") {
		t.Errorf("card does not carry its title:\n%s", got)
	}
	if !strings.Contains(got, "T001") {
		t.Errorf("card does not carry its T### identity:\n%s", got)
	}
	if strings.Contains(got, "O001") {
		t.Errorf("card shows its objective reference as display language:\n%s", got)
	}
}

// TestRenderCardReadsOnlyResolvedValues is the derive-nothing proof. The
// clearance says stale while naming no Check and carrying no freshness
// assessment, and the decision reports a dependency wait over a target that
// does not exist in any project. A renderer that re-derived any of it would
// have nothing to derive from and could not produce these badges.
func TestRenderCardReadsOnlyResolvedValues(t *testing.T) {
	card := fixtureCard(
		fixtureTask("T010", "Carry what the resolvers said", data.ColumnInProgress, data.StageAudit),
		data.Clearance{State: data.ClearanceStale},
		data.GateDecision{Blockers: []data.GateBlocker{
			{Kind: data.GateBlockDependency, Dependency: &data.DependencyBlock{Target: "T999"}},
			{Kind: data.GateBlockOwnerAcceptance},
		}},
	)

	got := renderedText(card, 44, false)

	for _, want := range []string{"[!] REVIEW", "WAITS T999", "OWNER"} {
		if !strings.Contains(got, want) {
			t.Errorf("card missing %q:\n%s", want, got)
		}
	}
}

// TestRenderCardOmitsBlockersTheClearanceBadgeAlreadyStates keeps the card from
// saying the same fact twice in two wordings.
func TestRenderCardOmitsBlockersTheClearanceBadgeAlreadyStates(t *testing.T) {
	card := fixtureCard(
		fixtureTask("T011", "One statement per fact", data.ColumnInProgress, data.StageAudit),
		data.Clearance{State: data.ClearanceNeedsWork},
		data.GateDecision{Blockers: []data.GateBlocker{{Kind: data.GateBlockClearanceNeedsWork}}},
	)

	got := renderedText(card, 44, false)

	if count := strings.Count(got, "NEEDS WORK"); count != 1 {
		t.Errorf("card states \"NEEDS WORK\" %d times, want exactly one:\n%s", count, got)
	}
}

// TestRenderCardDoneCardsShowOneReviewOutcomeAndNoCompletionBadge proves O012's
// retired vocabulary — "DONE", "BY EXCEPTION", "BY WAIVER", "Check (stale)" —
// never appears on a Done card, and that a Done card's own review outcome
// still tells an ordinary close, a stale one, and an owner-accepted one apart.
func TestRenderCardDoneCardsShowOneReviewOutcomeAndNoCompletionBadge(t *testing.T) {
	ordinary := fixtureCard(
		fixtureTask("T020", "Closed the ordinary way", data.ColumnDone, ""),
		data.Clearance{State: data.ClearanceCurrent},
		data.GateDecision{},
	)
	exception := TaskCard{
		Task:        fixtureTask("T021", "Closed by a recorded exception", data.ColumnDone, ""),
		Clearance:   data.Clearance{State: data.ClearanceNeedsWork},
		ByException: true,
	}
	stale := fixtureCard(
		fixtureTask("T022", "Closed, and the check went stale", data.ColumnDone, ""),
		data.Clearance{State: data.ClearanceStale},
		data.GateDecision{},
	)

	ordinaryText := renderedText(ordinary, 44, false)
	exceptionText := renderedText(exception, 44, false)
	staleText := renderedText(stale, 44, false)

	if !strings.Contains(ordinaryText, "[✓] CHECK") {
		t.Errorf("an ordinary done card does not read as checked:\n%s", ordinaryText)
	}
	if !strings.Contains(exceptionText, "[✓] OWNER ACCEPTED") {
		t.Errorf("a done-by-exception card does not name the owner's acceptance:\n%s", exceptionText)
	}
	if !strings.Contains(staleText, "[!] REVIEW") {
		t.Errorf("a stale done card does not read as needing review:\n%s", staleText)
	}
	for _, text := range []struct {
		name, got string
	}{{"ordinary", ordinaryText}, {"exception", exceptionText}, {"stale", staleText}} {
		for _, retired := range []string{"✓ DONE", "⚠ DONE", "BY EXCEPTION", "BY WAIVER", "Check (stale)"} {
			if strings.Contains(text.got, retired) {
				t.Errorf("%s done card still carries the retired badge %q:\n%s", text.name, retired, text.got)
			}
		}
	}
}

// TestRenderCardOpenExceptionOmitsTheOwnerBlockerItResolved proves an open
// Task allowed by a recorded exception shows "[✓] OWNER ACCEPTED" alone, with
// no separate "OWNER" blocker badge alongside it — matching
// data.ResolveTaskCompletion, which reports an exception-allowed decision
// with no Blockers at all, never the owner-acceptance blocker it overrode.
func TestRenderCardOpenExceptionOmitsTheOwnerBlockerItResolved(t *testing.T) {
	card := TaskCard{
		Task:        fixtureTask("T023", "Owner accepted the risk while still open", data.ColumnInProgress, data.StageAudit),
		Clearance:   data.Clearance{State: data.ClearanceNeedsWork},
		Decision:    data.GateDecision{Allowed: true, Actor: data.ActorRoleOwner, AllowedByException: true},
		ByException: true,
	}

	got := renderedText(card, 44, false)

	if !strings.Contains(got, "[✓] OWNER ACCEPTED") {
		t.Errorf("open exception card does not show OWNER ACCEPTED:\n%s", got)
	}
	if strings.Contains(got, "! OWNER") {
		t.Errorf("open exception card still shows the owner blocker its exception resolved:\n%s", got)
	}
}

func TestRenderCardStageAbsentOffAnInProgressTask(t *testing.T) {
	planned := renderedText(fixtureCard(
		fixtureTask("T030", "Not started", data.ColumnPlanned, ""),
		data.Clearance{State: data.ClearanceMissing},
		data.GateDecision{Allowed: true},
	), 40, false)

	for _, stage := range []string{"BUILD", "TEST", "CHECK"} {
		if strings.Contains(planned, stage) {
			t.Errorf("a planned card carries stage %q:\n%s", stage, planned)
		}
	}
}

// TestRenderCardFocusChangesColorNotGeometry holds the visual identity's rule
// that a layout must not move under focus.
func TestRenderCardFocusChangesColorNotGeometry(t *testing.T) {
	forceColorProfile(t, termenv.TrueColor)

	card := fixtureCard(
		fixtureTask("T040", "A card that must not move when focused", data.ColumnInProgress, data.StageBuild),
		data.Clearance{State: data.ClearanceCurrent},
		data.GateDecision{Allowed: true},
	)

	unfocused := renderCard(card, 40, false)
	focused := renderCard(card, 40, true)

	if lipgloss.Width(unfocused) != lipgloss.Width(focused) {
		t.Errorf("focused width %d, unfocused width %d", lipgloss.Width(focused), lipgloss.Width(unfocused))
	}
	if lipgloss.Height(unfocused) != lipgloss.Height(focused) {
		t.Errorf("focused height %d, unfocused height %d", lipgloss.Height(focused), lipgloss.Height(unfocused))
	}
	if xansi.Strip(unfocused) != xansi.Strip(focused) {
		t.Errorf("focus changed the card's text, not only its accent:\n%s\n%s", xansi.Strip(unfocused), xansi.Strip(focused))
	}
	if unfocused == focused {
		t.Error("focus produced an identical rendering; the accent must change")
	}
}

func TestRenderCardNeverExceedsItsWidth(t *testing.T) {
	card := fixtureCard(
		fixtureTask("T050", "A title long enough to need more than one line at any sensible width", data.ColumnInProgress, data.StageAudit),
		data.Clearance{State: data.ClearanceUnknown},
		data.GateDecision{Blockers: []data.GateBlocker{
			{Kind: data.GateBlockObjectiveDependency, ObjectiveDependency: &data.ObjectiveDependencyBlock{Target: "O009"}},
		}},
	)

	for _, width := range []int{20, 28, 40, 72} {
		for _, line := range strings.Split(renderedText(card, width, true), "\n") {
			if lipgloss.Width(line) > width {
				t.Errorf("at width %d a line is %d cells wide: %q", width, lipgloss.Width(line), line)
			}
		}
	}
}

// TestGroupTaskCardsGroupsByRecordedStatus proves the columns come from
// TaskV2.Status and that resolution reaches every Task.
func TestGroupTaskCardsGroupsByRecordedStatus(t *testing.T) {
	root := writeBadgeProject(t)
	loaded := loadProject(root)
	if loaded.Failed() {
		t.Fatalf("fixture project did not load: %s", loaded.Diagnostic)
	}

	grouped := groupTaskCards(loaded.State.Index)

	if len(grouped) != 3 {
		t.Errorf("grouped into %d columns, want exactly three", len(grouped))
	}
	counts := map[data.ColumnType]int{data.ColumnPlanned: 2, data.ColumnInProgress: 4, data.ColumnDone: 3}
	for column, want := range counts {
		if got := len(grouped[column]); got != want {
			t.Errorf("column %q holds %d cards, want %d", column, got, want)
		}
	}
	for column, cards := range grouped {
		for _, card := range cards {
			if card.Task.Status != column {
				t.Errorf("card %s with status %q sits in column %q", card.Task.ID, card.Task.Status, column)
			}
		}
	}

	planned := grouped[data.ColumnPlanned]
	if len(planned) < 2 || planned[0].Task.ID != "T001" || planned[1].Task.ID != "T002" {
		t.Errorf("planned cards = %v, want ascending Task ID order", cardIDs(planned))
	}
}

// TestGroupTaskCardsResolvesTheSameDecisionTheProjectionDoes proves a card and
// the Next area cannot disagree: for the Task the projection selected, the
// card's decision is the decision the projection carried.
func TestGroupTaskCardsResolvesTheSameDecisionTheProjectionDoes(t *testing.T) {
	root := writeBadgeProject(t)
	loaded := loadProject(root)
	if loaded.Failed() {
		t.Fatalf("fixture project did not load: %s", loaded.Diagnostic)
	}
	next := loaded.State.Next
	if next.Task == nil || next.GateDecision == nil {
		t.Fatalf("fixture projection selected no Task with a decision: %+v", next)
	}

	card, ok := findCard(groupTaskCards(loaded.State.Index), next.Task.ID)
	if !ok {
		t.Fatalf("no card for the projection's selected Task %s", next.Task.ID)
	}

	if card.Decision.Allowed != next.GateDecision.Allowed || len(card.Decision.Blockers) != len(next.GateDecision.Blockers) {
		t.Errorf("card decision %+v differs from the projection's %+v", card.Decision, *next.GateDecision)
	}
}

// TestCardReportsAnObjectiveLevelWait proves the Objective-dependency badge
// reaches a card from a real project: the Task itself declares no dependency,
// so the wait can only have come from the gate decision over its owning
// Objective.
func TestCardReportsAnObjectiveLevelWait(t *testing.T) {
	root := writeObjectiveDependencyProject(t)
	loaded := loadProject(root)
	if loaded.Failed() {
		t.Fatalf("fixture project did not load: %s", loaded.Diagnostic)
	}

	card, ok := findCard(groupTaskCards(loaded.State.Index), "T002")
	if !ok {
		t.Fatal("no card for T002")
	}
	if len(card.Task.DependsOn) != 0 {
		t.Fatalf("fixture Task declares its own dependencies: %+v", card.Task.DependsOn)
	}

	got := renderedText(card, 44, false)
	if !strings.Contains(got, "OBJECTIVE WAITS O001") {
		t.Errorf("card does not report the Objective-level wait:\n%s", got)
	}
}

func TestGroupTaskCardsHandlesAProjectWithNoIndex(t *testing.T) {
	grouped := groupTaskCards(nil)

	if len(grouped) != 3 {
		t.Fatalf("grouped into %d columns, want three empty ones", len(grouped))
	}
	for column, cards := range grouped {
		if len(cards) != 0 {
			t.Errorf("column %q holds %d cards over no index", column, len(cards))
		}
	}
}

func cardIDs(cards []TaskCard) []string {
	ids := make([]string, 0, len(cards))
	for _, card := range cards {
		ids = append(ids, card.Task.ID)
	}
	return ids
}

func findCard(grouped map[data.ColumnType][]TaskCard, id string) (TaskCard, bool) {
	for _, cards := range grouped {
		for _, card := range cards {
			if card.Task.ID == id {
				return card, true
			}
		}
	}
	return TaskCard{}, false
}

func TestRenderCardPlannedOmitsCheckBadge(t *testing.T) {
	card := fixtureCard(
		fixtureTask("T001", "Planned task", data.ColumnPlanned, ""),
		data.Clearance{State: data.ClearanceMissing},
		data.GateDecision{Allowed: true},
	)

	got := renderedText(card, 40, false)
	if strings.Contains(got, "CHECK") {
		t.Errorf("planned card should not render check badge:\n%s", got)
	}
}

// TestRenderCardInProgressOmitsCheckBadgeWhenNotActionable proves an open
// Task — at build, test, or audit stage alike — carries no "[ ] CHECK" or
// "[✓] CHECK" badge: "not checked yet" and "checked and clear" are
// completion-outcome vocabulary that belongs to the Done column (see
// TaskCard.showsReviewOutcome), not an open card, where it would either
// restate the stage badge or, worse, read as "all clear" beside a blocker
// that says otherwise.
func TestRenderCardInProgressOmitsCheckBadgeWhenNotActionable(t *testing.T) {
	for _, test := range []struct {
		name      string
		stage     data.ProgressStage
		clearance data.ClearanceState
	}{
		{"build, nothing recorded", data.StageBuild, data.ClearanceMissing},
		{"test, nothing recorded", data.StageTest, data.ClearanceMissing},
		{"audit, nothing recorded", data.StageAudit, data.ClearanceMissing},
		{"audit, checked and clear", data.StageAudit, data.ClearanceCurrent},
	} {
		t.Run(test.name, func(t *testing.T) {
			card := fixtureCard(
				fixtureTask("T002", "In progress task", data.ColumnInProgress, test.stage),
				data.Clearance{State: test.clearance},
				data.GateDecision{Allowed: true},
			)
			got := renderedText(card, 40, false)
			// "◆ CHECK" is the audit stage badge's own label, present
			// whenever stage is audit regardless of outcome; only the
			// bracketed review-outcome forms are what this test forbids.
			for _, retired := range []string{"[ ] CHECK", "[✓] CHECK"} {
				if strings.Contains(got, retired) {
					t.Errorf("%s: should not render %q:\n%s", test.name, retired, got)
				}
			}
		})
	}
}

// TestRenderCardInProgressStillShowsAnActionableOutcome proves the
// suppression above is narrow: a Task a NEEDS WORK Check sent back to build
// for repair, or whose clearance has otherwise gone stale, still shows that
// outcome at any stage — only the "nothing to report" cases (not yet
// checked, or checked and clear) are hidden on an open card.
func TestRenderCardInProgressStillShowsAnActionableOutcome(t *testing.T) {
	needsWork := renderedText(fixtureCard(
		fixtureTask("T003", "Sent back to build for repair", data.ColumnInProgress, data.StageBuild),
		data.Clearance{State: data.ClearanceNeedsWork},
		data.GateDecision{Blockers: []data.GateBlocker{{Kind: data.GateBlockClearanceNeedsWork}}},
	), 44, false)
	if !strings.Contains(needsWork, "[!] NEEDS WORK") {
		t.Errorf("build-stage repair card does not show its real outcome:\n%s", needsWork)
	}

	stale := renderedText(fixtureCard(
		fixtureTask("T004", "Testing again after clearance went stale", data.ColumnInProgress, data.StageTest),
		data.Clearance{State: data.ClearanceStale},
		data.GateDecision{},
	), 44, false)
	if !strings.Contains(stale, "[!] REVIEW") {
		t.Errorf("test-stage card with stale clearance does not show its real outcome:\n%s", stale)
	}
}

// TestRenderCardCurrentCheckAwaitingOwnerShowsOnlyTheOwnerBlocker is the
// concrete case that motivated the rule above: a Task at audit with a
// current, clear Check but still requiring owner sign-off used to show
// "[✓] CHECK" right next to the "AWAITS OWNER" blocker, reading as a
// contradiction — checked and clear, yet still blocked. The Check outcome is
// completion vocabulary for the Done column; while owner sign-off is
// outstanding, the blocker alone is the actionable fact.
func TestRenderCardCurrentCheckAwaitingOwnerShowsOnlyTheOwnerBlocker(t *testing.T) {
	card := fixtureCard(
		fixtureTask("T006", "At audit and waiting on the owner", data.ColumnInProgress, data.StageAudit),
		data.Clearance{State: data.ClearanceCurrent},
		data.GateDecision{Blockers: []data.GateBlocker{{Kind: data.GateBlockOwnerAcceptance}}},
	)
	got := renderedText(card, 44, false)

	// "◆ CHECK" is the audit stage badge's own label; only the bracketed
	// review-outcome form is what a current-and-clear Check must not add.
	if strings.Contains(got, "[✓] CHECK") {
		t.Errorf("card should not also show the completed-check outcome while owner sign-off is outstanding:\n%s", got)
	}
	if !strings.Contains(got, "AWAITS OWNER") {
		t.Errorf("card is missing the owner blocker:\n%s", got)
	}
}

func TestRenderCardPlannedFocusUsesMutedStyling(t *testing.T) {
	forceColorProfile(t, termenv.TrueColor)

	card := fixtureCard(
		fixtureTask("T001", "Planned task", data.ColumnPlanned, ""),
		data.Clearance{State: data.ClearanceMissing},
		data.GateDecision{Allowed: true},
	)

	unfocused := renderCard(card, 40, false)
	focused := renderCard(card, 40, true)

	if lipgloss.Width(unfocused) != lipgloss.Width(focused) {
		t.Errorf("focused width %d, unfocused width %d", lipgloss.Width(focused), lipgloss.Width(unfocused))
	}
	if lipgloss.Height(unfocused) != lipgloss.Height(focused) {
		t.Errorf("focused height %d, unfocused height %d", lipgloss.Height(focused), lipgloss.Height(unfocused))
	}
	// Focused planned card border should not contain orange
	if strings.Contains(focused, "252;99;35") { // 252;99;35 is #FC6323 in truecolor ANSI
		t.Errorf("focused planned card should not use orange accent:\n%s", focused)
	}
	// Title text color remains unfocused color
	if !strings.Contains(focused, "Planned task") {
		t.Errorf("focused planned card missing title:\n%s", focused)
	}
}

func TestRenderCardDoneFocusUsesGreenStyling(t *testing.T) {
	forceColorProfile(t, termenv.TrueColor)

	card := fixtureCard(
		fixtureTask("T001", "Done task", data.ColumnDone, ""),
		data.Clearance{State: data.ClearanceCurrent},
		data.GateDecision{},
	)

	unfocused := renderCard(card, 40, false)
	focused := renderCard(card, 40, true)

	if lipgloss.Width(unfocused) != lipgloss.Width(focused) {
		t.Errorf("focused width %d, unfocused width %d", lipgloss.Width(focused), lipgloss.Width(unfocused))
	}
	if lipgloss.Height(unfocused) != lipgloss.Height(focused) {
		t.Errorf("focused height %d, unfocused height %d", lipgloss.Height(focused), lipgloss.Height(unfocused))
	}
	// Focused done card should not contain orange
	if strings.Contains(focused, "252;99;35") {
		t.Errorf("focused done card should not use orange accent:\n%s", focused)
	}
	// Focused done card should contain green (163;198;56) for border and title
	if !strings.Contains(focused, "163;198;56") {
		t.Errorf("focused done card should use green accent:\n%s", focused)
	}
}

func TestRenderCardTitleWrapsUpToTwoLines(t *testing.T) {
	shortCard := fixtureCard(
		fixtureTask("T001", "Short title", data.ColumnPlanned, ""),
		data.Clearance{State: data.ClearanceMissing},
		data.GateDecision{Allowed: true},
	)
	shortGot := renderedText(shortCard, 30, false)
	if strings.Contains(shortGot, "…") {
		t.Errorf("short card title should not truncate:\n%s", shortGot)
	}

	wrapCard := fixtureCard(
		fixtureTask("T002", "Implement user authentication subsystem", data.ColumnPlanned, ""),
		data.Clearance{State: data.ClearanceMissing},
		data.GateDecision{Allowed: true},
	)
	wrapGot := renderedText(wrapCard, 30, false)
	// Outer width 30 -> textW = 30 - 4 = 26
	// "Implement user" (14) + "authentication" (14) -> 29 > 26
	// Line 1: "Implement user", Line 2: "authentication subsystem"
	if !strings.Contains(wrapGot, "Implement user") {
		t.Errorf("card title missing line 1:\n%s", wrapGot)
	}
	if !strings.Contains(wrapGot, "authentication subsystem") {
		t.Errorf("card title missing line 2:\n%s", wrapGot)
	}
	// Height of wrapped card is 1 line taller than short card
	shortLines := strings.Count(shortGot, "\n") + 1
	wrapLines := strings.Count(wrapGot, "\n") + 1
	if wrapLines != shortLines+1 {
		t.Errorf("wrapLines = %d, shortLines = %d (want difference of 1)", wrapLines, shortLines)
	}

	longCard := fixtureCard(
		fixtureTask("T003", "Implement user authentication subsystem with OAuth2 and SAML providers and tokens", data.ColumnPlanned, ""),
		data.Clearance{State: data.ClearanceMissing},
		data.GateDecision{Allowed: true},
	)
	longGot := renderedText(longCard, 30, false)
	if !strings.Contains(longGot, "…") {
		t.Errorf("long card title exceeding two lines should truncate with ellipsis:\n%s", longGot)
	}
	longLines := strings.Count(longGot, "\n") + 1
	if longLines != wrapLines {
		t.Errorf("longLines = %d, want equal to wrapLines %d (capped at 2 title lines)", longLines, wrapLines)
	}
}
