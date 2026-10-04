package v2

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	xansi "github.com/charmbracelet/x/ansi"
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/testutil"
)

const laneObjectiveExtra = "lanes:\n  - key: core\n    title: Lane / Proposed worktree — Core\n  - key: board\n    title: Lane / Proposed worktree — Board\n"

// writeLaneProject writes O-001 with two lanes and Tasks in every column:
// core holds a done, an in-progress and a planned Task, board holds a planned
// Task, and T-005 has no lane.
func writeLaneProject(t *testing.T, enabled bool) string {
	t.Helper()
	root := writeValidProject(t)
	config := "schema_version: 2\n"
	if enabled {
		config += "features:\n  parallel_planning: true\n"
	}
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), config)
	writeFixtureObjective(t, root, "O-001", "First objective", "in_progress", laneObjectiveExtra)
	writeTask(t, root, "O-001", "T-001", "Core done", "status: done\nlane: core\n")
	writeTask(t, root, "O-001", "T-002", "Board planned", "status: planned\nlane: board\n")
	writeTask(t, root, "O-001", "T-003", "Core running", "status: in_progress\nstage: build\nlane: core\n")
	writeTask(t, root, "O-001", "T-004", "Core planned", "status: planned\nlane: core\n")
	writeTask(t, root, "O-001", "T-005", "Sequential work", "status: planned\n")
	return root
}

func TestLaneHeadingsStayOffWhileTheOptionIsOffEvenWithSavedLanes(t *testing.T) {
	model := openSizedBoard(t, writeLaneProject(t, false), 120, 48)

	got := xansi.Strip(model.View())
	if strings.Contains(got, "Proposed worktree") || strings.Contains(got, ungroupedTitle) {
		t.Errorf("headings drawn with the option off:\n%s", got)
	}
	if ids := cardIDs(model.Cards[data.ColumnPlanned]); strings.Join(ids, " ") != "T-002 T-004 T-005" {
		t.Errorf("planned order = %v, want the ordinary ID order", ids)
	}
}

func TestLaneHeadingsFollowMembershipAcrossEveryColumn(t *testing.T) {
	model := openSizedBoard(t, writeLaneProject(t, true), 120, 48)

	if ids := cardIDs(model.Cards[data.ColumnPlanned]); strings.Join(ids, " ") != "T-004 T-002 T-005" {
		t.Errorf("planned order = %v, want core, board, then ungrouped", ids)
	}
	got := xansi.Strip(model.View())
	for _, want := range []string{"Lane / Proposed worktree — Core", "Lane / Proposed worktree — Board", ungroupedTitle} {
		if !strings.Contains(got, want) {
			t.Errorf("board lacks heading %q:\n%s", want, got)
		}
	}
	// core has a card in all three columns, so its heading appears in each.
	if n := strings.Count(got, "Proposed worktree — Core"); n != 3 {
		t.Errorf("core heading drawn %d times, want once per column (3):\n%s", n, got)
	}
	// board has only a planned card, so only one heading: empty ones are hidden.
	if n := strings.Count(got, "Proposed worktree — Board"); n != 1 {
		t.Errorf("board heading drawn %d times, want 1:\n%s", n, got)
	}
	for column, want := range map[data.ColumnType]int{data.ColumnPlanned: 3, data.ColumnInProgress: 1, data.ColumnDone: 1} {
		if got := len(model.Cards[column]); got != want {
			t.Errorf("column %s holds %d cards, want %d: headings must not change counts", column, got, want)
		}
	}
}

func TestLaneHeadingsAreNotSelectableAndFocusFollowsTaskIdentity(t *testing.T) {
	root := writeLaneProject(t, false)
	model := openSizedBoard(t, root, 120, 48)
	model = press(t, model, "down", "down") // T-005 in the ungrouped order
	if id := model.focusedTaskID(); id != "T-005" {
		t.Fatalf("setup focus = %q, want T-005", id)
	}

	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\nfeatures:\n  parallel_planning: true\n")
	updated, _ := model.Update(loadProject(root))
	model = updated.(Model)
	if id := model.focusedTaskID(); id != "T-005" {
		t.Errorf("focus after grouping = %q, want T-005 kept by identity", id)
	}
	model = press(t, model, "up")
	if id := model.focusedTaskID(); id != "T-002" {
		t.Errorf("up from T-005 = %q, want the previous card T-002, never a heading", id)
	}
	model = press(t, model, "up", "up", "up")
	if id := model.focusedTaskID(); id != "T-004" {
		t.Errorf("up past the first card = %q, want it to stop at T-004", id)
	}
}

func TestTurningTheOptionOffRestoresTheOrdinaryBoard(t *testing.T) {
	root := writeLaneProject(t, true)
	on := openSizedBoard(t, root, 120, 48)
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\nfeatures:\n  parallel_planning: false\n")
	updated, _ := on.Update(loadProject(root))
	off := updated.(Model)

	if ids := cardIDs(off.Cards[data.ColumnPlanned]); strings.Join(ids, " ") != "T-002 T-004 T-005" {
		t.Errorf("planned order after off = %v", ids)
	}
	if strings.Contains(xansi.Strip(off.View()), "Proposed worktree") {
		t.Error("headings remain after turning the option off")
	}
}

func TestLegacyObjectiveWithoutLanesKeepsItsLayoutWhenEnabled(t *testing.T) {
	root := writeValidProject(t)
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\nfeatures:\n  parallel_planning: true\n")
	model := openSizedBoard(t, root, 120, 48)

	for _, cards := range model.Cards {
		for _, card := range cards {
			if card.Heading.Present() {
				t.Errorf("%s has heading %+v in a project with no lanes", card.Task.ID, card.Heading)
			}
		}
	}
}

func TestGoalWideHeadingsNamespaceByObjective(t *testing.T) {
	root := writeLaneProject(t, true)
	writeFixtureObjective(t, root, "O-002", "Second objective", "planned", laneObjectiveExtra)
	writeTask(t, root, "O-002", "T-010", "Other core", "status: planned\nlane: core\n")
	model := openSizedBoard(t, root, 160, 60)
	model.selectObjective("")

	got := xansi.Strip(model.View())
	// Titles wrap inside the column, so look for the namespaced opening.
	for _, want := range []string{"▸ O-001 · Lane", "▸ O-002 · Lane"} {
		if !strings.Contains(got, want) {
			t.Errorf("goal-wide board lacks %q:\n%s", want, got)
		}
	}
	keys := map[string]bool{}
	for _, card := range model.Cards[data.ColumnPlanned] {
		keys[card.Heading.Key] = true
	}
	if !keys["O-001/core"] || !keys["O-002/core"] {
		t.Errorf("heading keys = %v, want O-001/core and O-002/core kept apart", keys)
	}
}

func TestLaneReadinessComesFromTheProjection(t *testing.T) {
	model := openSizedBoard(t, writeLaneProject(t, true), 120, 48)

	byKey := map[string]string{}
	for _, cards := range model.Cards {
		for _, card := range cards {
			byKey[card.Heading.Key] = card.Heading.Readiness
		}
	}
	if byKey["O-001/core"] != "in progress" {
		t.Errorf("core readiness = %q, want in progress", byKey["O-001/core"])
	}
	if byKey["O-001/"] != "" {
		t.Errorf("ungrouped readiness = %q, want none", byKey["O-001/"])
	}
}

func laneCards(lane, title string, count, from int) []TaskCard {
	cards := manyCards(count)
	for i := range cards {
		cards[i].Task.ID = fmt.Sprintf("T%03d", from+i)
		cards[i].Heading = LaneHeading{Key: "O-001/" + lane, Title: title}
	}
	return cards
}

func TestVisibleLaneWindowCountsHeadingHeightAndKeepsFocusVisible(t *testing.T) {
	heights := []int{4, 4, 4, 4}
	headingH := []int{2, 2, 2, 2}
	newGroup := []bool{true, false, true, false}

	// Card 0 costs its heading (2) and itself (4), card 1 costs 4, and the
	// "more" indicator one: 11 of 12 lines. Card 2 would add a heading.
	start, end := visibleLaneWindow(heights, headingH, newGroup, 12, 0)
	if start != 0 || end != 2 {
		t.Errorf("window = [%d,%d), want cards 0 and 1 inside 12 lines", start, end)
	}
	if start, end = visibleLaneWindow(heights, headingH, newGroup, 8, 0); start != 0 || end != 1 {
		t.Errorf("window = [%d,%d), want only card 0 once its heading is counted", start, end)
	}
	for focus := range heights {
		start, end = visibleLaneWindow(heights, headingH, newGroup, 12, focus)
		if focus < start || focus >= end {
			t.Errorf("window [%d,%d) loses focused card %d", start, end, focus)
		}
	}
	start, end = visibleLaneWindow([]int{20}, []int{2}, []bool{true}, 5, 0)
	if start != 0 || end != 1 {
		t.Errorf("window = [%d,%d), want the oversized card kept", start, end)
	}
}

func TestRenderColumnRepeatsTheHeadingWhenTheWindowStartsMidLane(t *testing.T) {
	cards := laneCards("core", "Lane / Proposed worktree — Core", 6, 1)
	for focus := range cards {
		got := xansi.Strip(renderColumn("PLANNED", cards, 34, 22, columnCursor{Card: focus, Holds: true, Focused: true}))

		if !strings.Contains(got, fmt.Sprintf("T%03d", focus+1)) {
			t.Errorf("focus %d: focused card scrolled out of the column:\n%s", focus, got)
		}
		if !strings.Contains(got, "Proposed worktree") {
			t.Errorf("focus %d: no heading in the window:\n%s", focus, got)
		}
		if lipgloss.Height(got) > 22 {
			t.Errorf("focus %d: column is %d lines, budget 22", focus, lipgloss.Height(got))
		}
		for _, line := range strings.Split(got, "\n") {
			if lipgloss.Width(line) > 34 {
				t.Errorf("focus %d: line %q wider than the column", focus, line)
			}
		}
	}
}

func TestRenderColumnCompactHeightDropsAHeadingBeforeClippingTheFocusedCard(t *testing.T) {
	cards := laneCards("core", "Lane / Proposed worktree — Planning and safety", 3, 1)
	got := xansi.Strip(renderColumn("PLANNED", cards, 30, 11, columnCursor{Card: 1, Holds: true, Focused: true}))

	if !strings.Contains(got, "Card number 2") || !strings.Contains(got, "T002") {
		t.Errorf("focused card clipped by its heading:\n%s", got)
	}
}

func TestRenderColumnWithoutHeadingsIsUnchanged(t *testing.T) {
	cards := manyCards(4)
	got := xansi.Strip(renderColumn("PLANNED", cards, 34, 40, columnCursor{}))
	if strings.Contains(got, "▸") {
		t.Errorf("ungrouped column drew a heading:\n%s", got)
	}
}

func TestLongHeadingTitlesWrapWithinTheColumn(t *testing.T) {
	cards := laneCards("a", "Lane / Proposed worktree — A very long lane title that cannot fit on one line", 1, 1)
	for _, width := range []int{14, 20, 34} {
		got := renderColumn("PLANNED", cards, width, 30, columnCursor{})
		for _, line := range strings.Split(xansi.Strip(got), "\n") {
			if lipgloss.Width(line) > width {
				t.Errorf("width %d: line %q is %d cells wide", width, line, lipgloss.Width(line))
			}
		}
	}
}

func TestHeadingPersistsWhileATaskMovesThroughEveryColumn(t *testing.T) {
	root := writeLaneProject(t, true)
	for _, status := range []string{"planned", "in_progress", "done"} {
		stage := ""
		if status == "in_progress" {
			stage = "stage: build\n"
		}
		writeTask(t, root, "O-001", "T-002", "Board planned", "status: "+status+"\n"+stage+"lane: board\n")
		model := openSizedBoard(t, root, 120, 48)

		card, ok := findCard(model.Cards, "T-002")
		if !ok {
			t.Fatalf("status %s: T-002 missing from the board", status)
		}
		if card.Heading.Key != "O-001/board" || card.Task.Status != data.ColumnType(status) {
			t.Errorf("status %s: heading %q, status %q", status, card.Heading.Key, card.Task.Status)
		}
	}
}

func TestCompactAndNarrowBoardsKeepGroupedColumnsInsideTheTerminal(t *testing.T) {
	root := writeLaneProject(t, true)
	for _, size := range [][2]int{{40, 24}, {30, 20}, {47, 30}, {80, 24}} {
		model := openSizedBoard(t, root, size[0], size[1])
		for _, keys := range [][]string{nil, {"down"}, {"down", "down"}} {
			shown := press(t, model, keys...)
			view := xansi.Strip(shown.View())
			for _, line := range strings.Split(view, "\n") {
				if lipgloss.Width(line) > size[0] {
					t.Errorf("%dx%d: line %q exceeds the terminal width", size[0], size[1], line)
				}
			}
			if id := shown.focusedTaskID(); id != "" && !strings.Contains(view, id) {
				t.Errorf("%dx%d after %v: focused %s is not visible:\n%s", size[0], size[1], keys, id, view)
			}
		}
	}
}
