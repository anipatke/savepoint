package v2

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	xansi "github.com/charmbracelet/x/ansi"
	"github.com/opencode/savepoint/internal/data"
)

// openSizedBoard opens the board over root at a given terminal size, through
// the same load the program runs.
func openSizedBoard(t *testing.T, root string, width, height int) Model {
	t.Helper()
	model := openBoard(t, root, "")
	model.Width = width
	model.Height = height
	return model
}

// TestBoardShowsTasksInTheColumnTheirStatusNames is the whole surface end to
// end: three columns, cards labelled by title, badges from the resolvers. The
// terminal is tall enough for every fixture card to be on screen at once below
// the Next area, so a badge missing here is a badge the board does not render
// rather than one that scrolled out of view.
func TestBoardShowsTasksInTheColumnTheirStatusNames(t *testing.T) {
	got := xansi.Strip(openSizedBoard(t, writeBadgeProject(t), 120, 48).View())

	for _, want := range []string{
		"PLANNED (2)", "IN PROGRESS (4)", "DONE (3)",
		"Planned with nothing recorded",
		"Being built right now",
		"Done and cleared",
		"▣ BUILD", "◇ TEST", "◆ CHECK",
		"[!] NEEDS WORK", "[✓] CHECK",
		"→ WAITS T001", "⚠ REPLAN", "! AWAITS OWNER",
		"[!] REVIEW", "[✓] OWNER ACCEPTED",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("board is missing %q:\n%s", want, got)
		}
	}

	// This fixture's build- and test-stage Tasks (T003, T004) have never been
	// checked, so their review outcome is suppressed rather than shown as an
	// unearned "[ ] CHECK" — see TaskCard.showsReviewOutcome. The pending
	// "[ ] CHECK" rendering itself is covered where it still belongs, at
	// audit stage: TestRenderCardAuditAlwaysShowsCheckBadgeEvenWhenMissing.
	if strings.Contains(got, "[ ] CHECK") {
		t.Errorf("board shows an unearned pending check badge on a build/test-stage card:\n%s", got)
	}

	// The retired completion vocabulary O012 removes from Task cards must
	// never reappear: the Done column carries completion, and one review
	// outcome badge carries the rest.
	for _, retired := range []string{"✓ DONE", "⚠ DONE", "BY EXCEPTION", "BY WAIVER", "Check (stale)", "Check (unverified)"} {
		if strings.Contains(got, retired) {
			t.Errorf("board still carries the retired badge %q:\n%s", retired, got)
		}
	}
}

func TestBoardDrawsExactlyThreeColumns(t *testing.T) {
	got := xansi.Strip(openSizedBoard(t, writeBadgeProject(t), 120, 40).View())

	for _, label := range []string{"PLANNED", "IN PROGRESS", "DONE"} {
		if count := strings.Count(got, label+" ("); count != 1 {
			t.Errorf("column %q appears %d times, want exactly one", label, count)
		}
	}
	// Nothing a badge reports is ever promoted to a column of its own.
	for _, notAColumn := range []string{"REPLAN (", "OWNER (", "AUDIT (", "BLOCKED"} {
		if strings.Contains(got, notAColumn) {
			t.Errorf("board grew a column for %q:\n%s", notAColumn, got)
		}
	}
}

func TestBoardLinesNeverExceedTheTerminalWidth(t *testing.T) {
	for _, width := range []int{80, 100, 120, 160} {
		got := xansi.Strip(openSizedBoard(t, writeBadgeProject(t), width, 30).View())
		for _, line := range strings.Split(got, "\n") {
			if lipgloss.Width(line) > width {
				t.Errorf("at width %d a line is %d cells wide: %q", width, lipgloss.Width(line), line)
			}
		}
	}
}

func TestBoardEmptyProjectStillDrawsThreeEmptyColumns(t *testing.T) {
	got := xansi.Strip(openSizedBoard(t, writeEmptyProjectFromTemplate(t), 100, 30).View())

	for _, label := range []string{"PLANNED (0)", "IN PROGRESS (0)", "DONE (0)"} {
		if !strings.Contains(got, label) {
			t.Errorf("empty project is missing column %q:\n%s", label, got)
		}
	}
	if strings.Count(got, "(empty)") != 3 {
		t.Errorf("empty project does not report all three columns as empty:\n%s", got)
	}
}

// TestNavigationMovesFocusAndClamps covers the column cursor: it moves, it
// clamps at both ends, and a repeated press at an end changes nothing.
func TestNavigationMovesFocusAndClamps(t *testing.T) {
	model := openSizedBoard(t, writeBadgeProject(t), 120, 40)

	if model.FocusedColumn != data.ColumnPlanned || model.FocusedCard != 0 {
		t.Fatalf("board opened focused on %q card %d, want the first planned card", model.FocusedColumn, model.FocusedCard)
	}

	model = press(t, model, "right")
	if model.FocusedColumn != data.ColumnInProgress {
		t.Errorf("right moved focus to %q, want in_progress", model.FocusedColumn)
	}
	model = press(t, model, "down", "down")
	if model.FocusedCard != 2 {
		t.Errorf("two downs left the cursor at %d, want 2", model.FocusedCard)
	}

	// Moving to a shorter column pulls the cursor back into range.
	model = press(t, model, "left")
	if model.FocusedColumn != data.ColumnPlanned || model.FocusedCard != 1 {
		t.Errorf("focus = %q card %d, want the planned column clamped to its last card", model.FocusedColumn, model.FocusedCard)
	}

	// Both ends clamp, and a repeated press at an end is a no-op.
	atStart := press(t, model, "left", "left", "up", "up", "up")
	if atStart.FocusedColumn != data.ColumnPlanned || atStart.FocusedCard != 0 {
		t.Errorf("focus = %q card %d, want it held at the first card of the first column", atStart.FocusedColumn, atStart.FocusedCard)
	}
	if press(t, atStart, "left").View() != atStart.View() {
		t.Error("pressing left at the first column changed the board")
	}

	atEnd := press(t, model, "right", "right", "right", "down", "down", "down", "down")
	if atEnd.FocusedColumn != data.ColumnDone || atEnd.FocusedCard != 2 {
		t.Errorf("focus = %q card %d, want it held at the last card of the last column", atEnd.FocusedColumn, atEnd.FocusedCard)
	}
	if press(t, atEnd, "down").View() != atEnd.View() {
		t.Error("pressing down at the last card changed the board")
	}
}

func TestNavigationOverAnEmptyColumnHoldsAValidCursor(t *testing.T) {
	model := press(t, openSizedBoard(t, writeEmptyProjectFromTemplate(t), 100, 30), "down", "right", "down")

	if model.FocusedCard != 0 {
		t.Errorf("FocusedCard = %d over a project with no Tasks, want 0", model.FocusedCard)
	}
}

// TestReloadClampsFocusIntoTheCardsThatRemain proves a reload that removed the
// focused Task leaves a focus something renders.
func TestReloadClampsFocusIntoTheCardsThatRemain(t *testing.T) {
	root := writeBadgeProject(t)
	model := press(t, openSizedBoard(t, root, 120, 40), "right", "down", "down", "down")
	if model.FocusedCard != 3 {
		t.Fatalf("FocusedCard = %d, want the fourth in-progress card", model.FocusedCard)
	}

	removeTask(t, root, "O001", "T004")
	removeTask(t, root, "O001", "T005")
	removeTask(t, root, "O001", "T006")
	removeCheck(t, root, "C001")
	removeCheck(t, root, "C002")
	reloaded, _ := model.Update(loadCmd(root)().(projectLoadedMsg))

	after := reloaded.(Model)
	if after.FocusedCard != 0 {
		t.Errorf("FocusedCard = %d after the column shrank to one card, want 0", after.FocusedCard)
	}
	if strings.Contains(xansi.Strip(after.View()), "Being tested") {
		t.Error("the board still shows a Task the reload removed")
	}
}

// press sends a sequence of keys and returns the resulting model.
func press(t *testing.T, model Model, keys ...string) Model {
	t.Helper()
	current := tea.Model(model)
	for _, key := range keys {
		updated, _ := current.Update(keyMsg(key))
		current = updated
	}
	next, ok := current.(Model)
	if !ok {
		t.Fatalf("Update() returned %T, want Model", current)
	}
	return next
}

func keyMsg(key string) tea.KeyMsg {
	switch key {
	case "left":
		return tea.KeyMsg{Type: tea.KeyLeft}
	case "right":
		return tea.KeyMsg{Type: tea.KeyRight}
	case "up":
		return tea.KeyMsg{Type: tea.KeyUp}
	case "down":
		return tea.KeyMsg{Type: tea.KeyDown}
	case "tab":
		return tea.KeyMsg{Type: tea.KeyTab}
	case "enter":
		return tea.KeyMsg{Type: tea.KeyEnter}
	case "esc":
		return tea.KeyMsg{Type: tea.KeyEsc}
	default:
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
	}
}
