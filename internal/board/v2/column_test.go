package v2

import (
	"fmt"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	xansi "github.com/charmbracelet/x/ansi"
	"github.com/muesli/termenv"
	"github.com/opencode/savepoint/internal/data"
)

// manyCards builds count planned cards, enough to overflow any column height a
// test gives them.
func manyCards(count int) []TaskCard {
	cards := make([]TaskCard, 0, count)
	for i := 1; i <= count; i++ {
		cards = append(cards, fixtureCard(
			fixtureTask(fmt.Sprintf("T%03d", i), fmt.Sprintf("Card number %d", i), data.ColumnPlanned, ""),
			data.Clearance{State: data.ClearanceMissing},
			data.GateDecision{Allowed: true},
		))
	}
	return cards
}

func TestRenderColumnHeadsWithItsLabelAndCount(t *testing.T) {
	got := xansi.Strip(renderColumn("PLANNED", manyCards(2), 34, 20, columnCursor{}))

	if !strings.Contains(got, "PLANNED (2)") {
		t.Errorf("column header missing its label and count:\n%s", got)
	}
	if !strings.Contains(got, "Card number 1") || !strings.Contains(got, "Card number 2") {
		t.Errorf("column does not carry its cards:\n%s", got)
	}
}

func TestRenderColumnEmptySaysSo(t *testing.T) {
	got := xansi.Strip(renderColumn("DONE", nil, 34, 20, columnCursor{}))

	if !strings.Contains(got, "DONE (0)") || !strings.Contains(got, "(empty)") {
		t.Errorf("empty column does not read as empty:\n%s", got)
	}
}

// TestRenderColumnScrollsWithIndicators proves a column taller than its budget
// shows what is out of view rather than silently dropping it.
func TestRenderColumnScrollsWithIndicators(t *testing.T) {
	cards := manyCards(10)

	top := xansi.Strip(renderColumn("PLANNED", cards, 34, 16, columnCursor{Holds: true, Focused: true}))
	if !strings.Contains(top, "↓") || !strings.Contains(top, "more") {
		t.Errorf("a column with cards below the fold shows no indicator:\n%s", top)
	}
	if strings.Contains(top, "above") {
		t.Errorf("a column focused on its first card claims cards above it:\n%s", top)
	}

	bottom := xansi.Strip(renderColumn("PLANNED", cards, 34, 16, columnCursor{Card: len(cards) - 1, Holds: true, Focused: true}))
	if !strings.Contains(bottom, "↑") || !strings.Contains(bottom, "above") {
		t.Errorf("a column scrolled to its last card shows no indicator for what is above:\n%s", bottom)
	}
	if !strings.Contains(bottom, "Card number 10") {
		t.Errorf("the focused card scrolled out of view:\n%s", bottom)
	}
}

// TestRenderColumnKeepsItsBudget proves the window and the frame agree: a
// column never renders past the height it was given, at any focus position.
func TestRenderColumnKeepsItsBudget(t *testing.T) {
	cards := manyCards(12)

	for _, height := range []int{10, 16, 24} {
		for _, focus := range []int{0, 5, len(cards) - 1} {
			rendered := renderColumn("PLANNED", cards, 34, height, columnCursor{Card: focus, Holds: true, Focused: true})
			if got := lipgloss.Height(rendered); got != height {
				t.Errorf("height %d with focus %d rendered %d lines", height, focus, got)
			}
		}
	}
}

func TestRenderColumnNeverExceedsItsWidth(t *testing.T) {
	cards := manyCards(3)

	for _, width := range []int{18, 26, 34, 60} {
		rendered := xansi.Strip(renderColumn("IN PROGRESS", cards, width, 20, columnCursor{}))
		want := lipgloss.Width(rendered)
		for _, line := range strings.Split(rendered, "\n") {
			if lipgloss.Width(line) > want {
				t.Errorf("at width %d a line is %d cells wide: %q", width, lipgloss.Width(line), line)
			}
		}
	}
}

// TestRenderColumnFocusChangesColorNotGeometry is the column half of the
// focus-geometry rule the cards already hold.
func TestRenderColumnFocusChangesColorNotGeometry(t *testing.T) {
	forceColorProfile(t, termenv.TrueColor)
	cards := manyCards(3)

	unfocused := renderColumn("PLANNED", cards, 34, 20, columnCursor{})
	focused := renderColumn("PLANNED", cards, 34, 20, columnCursor{Holds: true, Focused: true})

	if lipgloss.Width(unfocused) != lipgloss.Width(focused) {
		t.Errorf("focused width %d, unfocused width %d", lipgloss.Width(focused), lipgloss.Width(unfocused))
	}
	if lipgloss.Height(unfocused) != lipgloss.Height(focused) {
		t.Errorf("focused height %d, unfocused height %d", lipgloss.Height(focused), lipgloss.Height(unfocused))
	}
	if xansi.Strip(unfocused) != xansi.Strip(focused) {
		t.Error("focus changed the column's text, not only its accent")
	}
	if unfocused == focused {
		t.Error("focus produced an identical column; the accent must change")
	}
}

// TestVisibleWindowKeepsTheFocusedCardVisible exercises the windowing directly,
// including the card taller than the whole budget that must still render.
func TestVisibleWindowKeepsTheFocusedCardVisible(t *testing.T) {
	heights := []int{4, 4, 4, 4, 4}

	start, end := visibleWindow(heights, 12, 0)
	if start != 0 || end <= start {
		t.Errorf("window = [%d,%d) focused on the first card", start, end)
	}

	start, end = visibleWindow(heights, 12, 4)
	if 4 < start || 4 >= end {
		t.Errorf("window = [%d,%d), which does not contain the focused card 4", start, end)
	}

	start, end = visibleWindow([]int{20}, 5, 0)
	if start != 0 || end != 1 {
		t.Errorf("window = [%d,%d), want the one oversized card rendered anyway", start, end)
	}
}

func TestRenderColumnPlannedFocusUsesPlannedAccent(t *testing.T) {
	forceColorProfile(t, termenv.TrueColor)
	cards := manyCards(2)

	focusedPlanned := renderColumn("PLANNED", cards, 34, 20, columnCursor{Holds: true, Focused: true})
	// Focused planned column should not use orange accent (RGB 252, 99, 35)
	if strings.Contains(focusedPlanned, "252;99;35") {
		t.Errorf("focused planned column should not use orange accent:\n%s", focusedPlanned)
	}

	focusedInProgress := renderColumn("IN PROGRESS", cards, 34, 20, columnCursor{Holds: true, Focused: true})
	// Focused in_progress column should use orange accent
	if !strings.Contains(focusedInProgress, "252;99;35") {
		t.Errorf("focused in_progress column should use orange accent:\n%s", focusedInProgress)
	}

	focusedDone := renderColumn("DONE", cards, 34, 20, columnCursor{Holds: true, Focused: true})
	// Focused done column should use green accent (RGB 164, 198, 57) and not orange
	if strings.Contains(focusedDone, "252;99;35") {
		t.Errorf("focused done column should not use orange accent:\n%s", focusedDone)
	}
	if !strings.Contains(focusedDone, "163;198;56") {
		t.Errorf("focused done column should use green accent (163;198;56):\n%s", focusedDone)
	}
}
