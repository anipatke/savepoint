package v2

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	xansi "github.com/charmbracelet/x/ansi"
	"github.com/opencode/savepoint/internal/data"
)

func liveCard(status data.ColumnType, stage data.ProgressStage, frame int) TaskCard {
	card := fixtureCard(
		fixtureTask("T-001", "Animate the tile", status, stage),
		data.Clearance{State: data.ClearanceMissing},
		data.GateDecision{Allowed: true},
	)
	card.Live, card.Frame = true, frame
	return card
}

func TestLiveBuildTileShowsSpinnerAndBuilding(t *testing.T) {
	got := renderedText(liveCard(data.ColumnInProgress, data.StageBuild, 0), 40, false)
	if !strings.Contains(got, spinnerFrame(0)+" Building…") {
		t.Errorf("build tile has no spinner and Building…:\n%s", got)
	}
	if strings.Contains(got, "BUILD") {
		t.Errorf("build tile still shows the static badge:\n%s", got)
	}

	// Failure path: a card not drawn on the board keeps the static badge.
	card := liveCard(data.ColumnInProgress, data.StageBuild, 0)
	card.Live = false
	if static := renderedText(card, 40, false); !strings.Contains(static, "▣ BUILD") || strings.Contains(static, "Building…") {
		t.Errorf("a non-live build card must keep the static badge:\n%s", static)
	}
}

func TestLiveTestTileShowsSpinnerAndTesting(t *testing.T) {
	got := renderedText(liveCard(data.ColumnInProgress, data.StageTest, 3), 40, false)
	if !strings.Contains(got, spinnerFrame(3)+" Testing…") {
		t.Errorf("test tile has no spinner and Testing…:\n%s", got)
	}
	if strings.Contains(got, "TEST") {
		t.Errorf("test tile still shows the static badge:\n%s", got)
	}

	// Failure path: a stage outside build and test is not given words.
	if _, ok := liveStageBadge(data.ColumnInProgress, data.ProgressStage("deploy"), 0); ok {
		t.Error("an unknown stage must not get a live badge")
	}
}

func TestSpinnerFrameWrapsAndIgnoresNegativeCounters(t *testing.T) {
	n := len(spinnerFrames)
	if spinnerFrame(n) != spinnerFrame(0) {
		t.Error("the frame counter does not wrap")
	}
	if got := spinnerFrame(-1); got != spinnerFrames[n-1] {
		t.Errorf("negative frame = %q, want the last frame", got)
	}
}

func TestOtherTilesAreUnchangedWhenLive(t *testing.T) {
	cases := []struct {
		name   string
		status data.ColumnType
		stage  data.ProgressStage
	}{
		{"planned", data.ColumnPlanned, ""},
		{"audit", data.ColumnInProgress, data.StageAudit},
		{"done", data.ColumnDone, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			live := liveCard(tc.status, tc.stage, 4)
			plain := live
			plain.Live = false
			if a, b := renderCard(live, 40, false), renderCard(plain, 40, false); a != b {
				t.Errorf("%s tile changes when live:\nlive:\n%s\nstatic:\n%s", tc.name, a, b)
			}
			if strings.Contains(renderedText(live, 40, false), "…ing") {
				t.Errorf("%s tile shows a live word", tc.name)
			}
		})
	}
}

func TestLiveLabelKeepsOtherBadgesInOrder(t *testing.T) {
	card := liveCard(data.ColumnInProgress, data.StageBuild, 0)
	card.ByWaiver = true
	card.Decision = data.GateDecision{Blockers: []data.GateBlocker{{Kind: data.GateBlockReplan}}}
	got := renderedText(card, 60, false)
	live, waived, replan := strings.Index(got, "Building…"), strings.Index(got, "WAIVED"), strings.Index(got, "REPLAN")
	if live < 0 || waived < live || replan < waived {
		t.Errorf("badge order is not live, review outcome, blocker:\n%s", got)
	}
}

func TestLiveLabelIsCutNotWrappedOnANarrowTile(t *testing.T) {
	card := liveCard(data.ColumnInProgress, data.StageBuild, 0)
	got := renderCard(card, 11, false)
	for _, line := range strings.Split(xansi.Strip(got), "\n") {
		if w := xansi.StringWidth(line); w > 11 {
			t.Errorf("line %q is %d cells wide on an 11-cell tile", line, w)
		}
	}
	if strings.Contains(xansi.Strip(got), "ing…\n") && strings.Count(xansi.Strip(got), "Building") > 1 {
		t.Errorf("the label wrapped:\n%s", got)
	}
	static := card
	static.Live = false
	if lines := strings.Count(got, "\n"); lines != strings.Count(renderCard(static, 11, false), "\n") {
		t.Errorf("the live tile is a different height from the static one (%d lines)", lines+1)
	}
}

func liveModel(stage data.ProgressStage) Model {
	m := NewModel(Options{})
	m.Cards[data.ColumnInProgress] = []TaskCard{liveCard(data.ColumnInProgress, stage, 0)}
	return m
}

func TestSpinnerTickAdvancesTheFrameAndReschedules(t *testing.T) {
	m := liveModel(data.StageBuild)
	m.spinnerRunning = true
	updated, cmd := m.Update(spinnerTickMsg{})
	got := updated.(Model)
	if got.spinnerFrame != 1 || cmd == nil || !got.spinnerRunning {
		t.Errorf("tick: frame %d, rescheduled %v, running %v; want 1, true, true", got.spinnerFrame, cmd != nil, got.spinnerRunning)
	}

	// Failure path: with nothing live the tick neither advances nor reschedules.
	idle := NewModel(Options{})
	idle.spinnerRunning = true
	updated, cmd = idle.Update(spinnerTickMsg{})
	if got := updated.(Model); got.spinnerFrame != 0 || cmd != nil || got.spinnerRunning {
		t.Errorf("idle tick: frame %d, rescheduled %v, running %v; want 0, false, false", got.spinnerFrame, cmd != nil, got.spinnerRunning)
	}
}

func TestSpinnerStopsAfterTheLastLiveTileMovesOn(t *testing.T) {
	m := liveModel(data.StageTest)
	m.spinnerRunning = true
	m.Cards[data.ColumnInProgress][0].Task.Stage = data.StageAudit
	updated, cmd := m.Update(spinnerTickMsg{})
	if cmd != nil || updated.(Model).spinnerRunning {
		t.Error("the tick kept running after the last live tile moved to check")
	}

	// Failure path: a second live tile keeps it going.
	m = liveModel(data.StageTest)
	m.spinnerRunning = true
	m.Cards[data.ColumnInProgress] = append(m.Cards[data.ColumnInProgress], liveCard(data.ColumnInProgress, data.StageBuild, 0))
	m.Cards[data.ColumnInProgress][0].Task.Stage = data.StageAudit
	if _, cmd := m.Update(spinnerTickMsg{}); cmd == nil {
		t.Error("the tick stopped while a live tile remained")
	}
}

func TestIdleBoardStartsNoTimer(t *testing.T) {
	m := NewModel(Options{})
	if cmd := m.spinnerCmd(); cmd != nil || m.spinnerRunning {
		t.Error("an empty board started a timer")
	}
	planned := NewModel(Options{})
	planned.Cards[data.ColumnPlanned] = []TaskCard{liveCard(data.ColumnPlanned, "", 0)}
	if cmd := planned.spinnerCmd(); cmd != nil {
		t.Error("a board with only a planned tile started a timer")
	}

	// Failure path: a live tile starts exactly one timer, not one per call.
	live := liveModel(data.StageBuild)
	if cmd := live.spinnerCmd(); cmd == nil || !live.spinnerRunning {
		t.Fatal("a live tile did not start the timer")
	}
	if cmd := live.spinnerCmd(); cmd != nil {
		t.Error("a second call started a second timer")
	}
}

func reload(t *testing.T, m Model, root string) (Model, tea.Cmd) {
	t.Helper()
	updated, cmd := m.Update(loadCmd(root)().(projectLoadedMsg))
	return updated.(Model), cmd
}

func TestStageChangeBetweenReloadsStartsAndStopsTheTimer(t *testing.T) {
	root := writeBadgeProject(t)
	m := openBoard(t, root, "")
	m.spinnerRunning = false

	// The fixture has build and test Tasks, so a reload starts the timer.
	m, cmd := reload(t, m, root)
	if cmd == nil || !m.spinnerRunning {
		t.Fatal("a reload with live tiles did not start the timer")
	}

	// Both live Tasks move on to check: the next reload starts nothing.
	writeTask(t, root, "O-001", "T-003", "Being built right now", "status: in_progress\nstage: audit\n")
	writeTask(t, root, "O-001", "T-004", "Being tested, and the plan slipped", "status: in_progress\nstage: audit\n")
	m.spinnerRunning = false
	m, cmd = reload(t, m, root)
	if cmd != nil || m.spinnerRunning {
		t.Error("a reload with no live tile started a timer")
	}

	// Failure path: a Task sent back to build brings the timer back.
	writeTask(t, root, "O-001", "T-003", "Being built right now", "status: in_progress\nstage: build\n")
	if _, cmd = reload(t, m, root); cmd == nil {
		t.Error("a Task moved back to build did not restart the timer")
	}
}

func TestPlainOutputCarriesNoLiveLabel(t *testing.T) {
	m := openBoard(t, writeBadgeProject(t), "")
	m.spinnerFrame = 2
	for _, cards := range m.Cards {
		for _, card := range cards {
			for _, badge := range card.badges() {
				if strings.Contains(badge.Text(), "…") {
					t.Errorf("plain badges carry a live label: %q", badge.Text())
				}
			}
		}
	}
}

func TestSpinnerRestartsWhenAnObjectiveSwitchShowsALiveTileAgain(t *testing.T) {
	root := writeBadgeProject(t)
	writeObjective(t, root, "O-002", "Quiet", "planned")
	writeTask(t, root, "O-002", "T-900", "Not started yet", "status: planned\n")
	m := openBoard(t, root, "")
	if !m.spinnerRunning {
		t.Fatal("the opened board with live tiles has no timer")
	}

	// Narrowing to an Objective with no live tile lets the timer end.
	m.selectObjective("O-002")
	updated, _ := m.Update(spinnerTickMsg{})
	m = updated.(Model)
	if m.spinnerRunning {
		t.Fatal("the timer kept running with no live tile shown")
	}

	// Failure path from I-155: going back to the live view must restart it.
	m.SidebarFocused = true
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = updated.(Model)
	if !m.hasLiveTile() {
		t.Fatal("esc did not return to a view with a live tile")
	}
	if cmd == nil || !m.spinnerRunning {
		t.Error("the spinner stayed frozen after switching back to a live tile")
	}
}
