package v2

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	xansi "github.com/charmbracelet/x/ansi"
)

func TestTruncateCellsUsesTerminalWidthAndWholeGraphemes(t *testing.T) {
	text := "A界e\u0301🙂B"
	for width := 1; width <= lipgloss.Width(text); width++ {
		got := truncateCells(text, width)
		if actual := lipgloss.Width(got); actual > width {
			t.Errorf("width %d produced %d cells: %q", width, actual, got)
		}
		if lipgloss.Width(text) > width && !strings.HasSuffix(got, "…") {
			t.Errorf("width %d truncated without a visible ellipsis: %q", width, got)
		}
	}

	// The combining mark remains attached to its base when the string is
	// clipped at the grapheme boundary. A byte or rune slice would leave either
	// a dangling mark or half of the emoji here.
	got := truncateCells(text, lipgloss.Width("A界e\u0301")+1)
	if got != "A界e\u0301…" {
		t.Errorf("grapheme-aware truncation = %q, want the complete accented grapheme", got)
	}
	if xansi.StringWidth(got) > lipgloss.Width("A界e\u0301")+1 {
		t.Errorf("truncated text exceeds its cell budget: %q", got)
	}
}

func TestNarrowBoardCollapsesInDefinedOrder(t *testing.T) {
	root := writeBadgeProject(t)

	withoutSidebar := xansi.Strip(openSizedBoard(t, root, sidebarBreakpoint-1, 32).View())
	if strings.Contains(withoutSidebar, sidebarTitle) {
		t.Fatalf("sidebar remains visible below its breakpoint:\n%s", withoutSidebar)
	}
	for _, label := range []string{"PLANNED (", "IN PROGRESS (", "DONE ("} {
		if !strings.Contains(withoutSidebar, label) {
			t.Errorf("three-column layout lost %q before compact collapse:\n%s", label, withoutSidebar)
		}
	}

	compact := xansi.Strip(openSizedBoard(t, root, compactBoardBreakpoint-1, 32).View())
	focused := "PLANNED ("
	if !strings.Contains(compact, focused) {
		t.Fatalf("compact layout does not retain the focused column %q:\n%s", focused, compact)
	}
	for _, label := range []string{"IN PROGRESS (", "DONE ("} {
		if strings.Contains(compact, label) {
			t.Errorf("compact layout still renders collapsed column %q:\n%s", label, compact)
		}
	}
}

func TestSidebarBreakpointKeepsTaskColumnsReadable(t *testing.T) {
	root := writeBadgeProject(t)
	outerAtBreakpoint := sidebarBreakpoint + boardMarginX*2

	visible := openSizedBoard(t, root, outerAtBreakpoint, 32)
	visibleText := xansi.Strip(visible.View())
	if !strings.Contains(visibleText, sidebarTitle) {
		t.Fatalf("sidebar is hidden at its content-width breakpoint %d:\n%s", sidebarBreakpoint, visibleText)
	}
	if got := columnWidth(sidebarBreakpoint); got != 30 {
		t.Errorf("Task column width at the sidebar breakpoint = %d, want 30", got)
	}
	assertSurfaceFits(t, "visible sidebar", visible.View(), outerAtBreakpoint)

	narrow := openSizedBoard(t, root, outerAtBreakpoint-1, 32)
	narrowText := xansi.Strip(narrow.View())
	if strings.Contains(narrowText, sidebarTitle) {
		t.Errorf("sidebar remains visible one cell below its content-width breakpoint:\n%s", narrowText)
	}
	assertSurfaceFits(t, "collapsed sidebar", narrow.View(), outerAtBreakpoint-1)
}

func TestEveryV2SurfaceFitsNarrowWidths(t *testing.T) {
	root := writeEvidenceProject(t)
	for _, width := range []int{20, 24, 32, compactBoardBreakpoint - 1, compactBoardBreakpoint, 80, 120} {
		t.Run(string(rune(width)), func(t *testing.T) {
			model := openSizedBoard(t, root, width, 36)
			assertSurfaceFits(t, "board", model.View(), width)
			assertSurfaceFits(t, "detail", press(t, model, "enter").View(), width)
			assertSurfaceFits(t, "issues", press(t, model, "i").View(), width)
			assertSurfaceFits(t, "issue detail", press(t, model, "i", "enter").View(), width)
			assertSurfaceFits(t, "help", press(t, model, "?").View(), width)
		})
	}
}

func assertSurfaceFits(t *testing.T, name, rendered string, width int) {
	t.Helper()
	for lineNo, line := range strings.Split(xansi.Strip(rendered), "\n") {
		if actual := lipgloss.Width(line); actual > width {
			t.Errorf("%s line %d is %d cells wide at width %d: %q", name, lineNo+1, actual, width, line)
		}
	}
}

func TestWideAndCombiningTaskTitleIsMeasuredInCells(t *testing.T) {
	root := writeValidProject(t)
	writeTaskBody(t, root, "O-001", "T-001", "界界 e\u0301 🙂 a title that must truncate", "status: planned\n", "")

	model := openSizedBoard(t, root, 20, 30)
	got := xansi.Strip(model.View())
	if !strings.Contains(got, "…") {
		t.Fatalf("wide title was not visibly truncated:\n%s", got)
	}
	assertSurfaceFits(t, "wide title board", model.View(), 20)
}

func TestFullBrowseAndReloadLeaveProjectUntouched(t *testing.T) {
	root := writeEvidenceProject(t)
	before := snapshotProject(t, root)
	model := openSizedBoard(t, root, 120, 40)

	// Visit the sidebar, every Issues filter, both detail surfaces, help, and a
	// reload. These are all read operations; only the explicit owner action keys
	// are allowed to reach the write commands.
	model = press(t, model, "left", "down", "v", "down", "esc", "right")
	model = press(t, model, "i")
	for range issueFilterOrder {
		model = press(t, model, "f")
	}
	model = press(t, model, "enter", "down", "esc", "esc", "?", "esc")
	updated, _ := model.Update(loadCmd(root)())
	model, ok := updated.(Model)
	if !ok {
		t.Fatalf("reload returned %T, want Model", updated)
	}
	_ = model.View()

	after := snapshotProject(t, root)
	if len(before) != len(after) {
		t.Fatalf("browse changed file count from %d to %d", len(before), len(after))
	}
	for path, state := range before {
		if got := after[path]; got != state {
			t.Errorf("%s changed during browse/reload:\n before %s\n  after %s", path, state, got)
		}
	}
}

func TestWrapTitleLines(t *testing.T) {
	// Zero / negative limits
	if wrapTitleLines("hello", 10, 0) != nil {
		t.Error("maxLines 0 should return nil")
	}
	if wrapTitleLines("hello", 0, 2) != nil {
		t.Error("width 0 should return nil")
	}
	if got := wrapTitleLines("", 10, 2); len(got) != 1 || got[0] != "" {
		t.Errorf("empty text should return [\"\"], got %v", got)
	}

	// Fits on one line
	got := wrapTitleLines("Short text", 20, 2)
	if len(got) != 1 || got[0] != "Short text" {
		t.Errorf("single-line text wrapped unexpectedly: %v", got)
	}

	// Wraps across two lines cleanly at word boundary
	got = wrapTitleLines("The quick brown fox jumps", 15, 2)
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d: %v", len(got), got)
	}
	if got[0] != "The quick brown" {
		t.Errorf("line 1 = %q, want %q", got[0], "The quick brown")
	}
	if got[1] != "fox jumps" {
		t.Errorf("line 2 = %q, want %q", got[1], "fox jumps")
	}

	// Exceeds two lines, truncates line two with ellipsis
	got = wrapTitleLines("The quick brown fox jumps over the lazy dog repeatedly", 15, 2)
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d: %v", len(got), got)
	}
	if got[0] != "The quick brown" {
		t.Errorf("line 1 = %q, want %q", got[0], "The quick brown")
	}
	if !strings.HasSuffix(got[1], "…") {
		t.Errorf("line 2 = %q, want suffix '…'", got[1])
	}

	// Very long word exceeding width across two lines
	got = wrapTitleLines("Supercalifragilisticexpialidocious", 10, 2)
	if len(got) != 2 {
		t.Fatalf("expected 2 lines for long word, got %d: %v", len(got), got)
	}
	if !strings.HasSuffix(got[1], "…") {
		t.Errorf("line 2 = %q, want suffix '…'", got[1])
	}
}

// TestSidebarGrowsWithTerminalWidth proves the sidebar takes a share of a wide
// terminal instead of staying at its breakpoint width, never at the cost of the
// Task columns' breakpoint width, and stops growing at its cap.
func TestSidebarGrowsWithTerminalWidth(t *testing.T) {
	root := writeBadgeProject(t)
	for _, tc := range []struct {
		name  string
		width int
		want  int
	}{
		{"breakpoint", sidebarBreakpoint, sidebarMinWidth},
		{"limited by column room", sidebarBreakpoint + 2, sidebarMinWidth + 2},
		{"intermediate", 160, 48},
		{"cap", 200, sidebarMaxWidth},
		{"past cap", 300, sidebarMaxWidth},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := sidebarWidth(tc.width); got != tc.want {
				t.Errorf("sidebarWidth(%d) = %d, want %d", tc.width, got, tc.want)
			}
			if got := columnWidth(tc.width); got < 30 {
				t.Errorf("Task column width at %d = %d, below its 30-cell breakpoint width", tc.width, got)
			}
			outer := tc.width + boardMarginX*2
			model := openSizedBoard(t, root, outer, 32)
			if !strings.Contains(xansi.Strip(model.View()), sidebarTitle) {
				t.Fatalf("sidebar is hidden at content width %d", tc.width)
			}
			assertSurfaceFits(t, tc.name, model.View(), outer)
		})
	}
}
