package v2

import (
	"strings"

	xansi "github.com/charmbracelet/x/ansi"
)

const (
	// compactBoardBreakpoint is the last width at which three columns can be
	// shown without making each one narrower than a readable card. Below it,
	// the sidebar has already collapsed and the board shows only the focused
	// column. This is an intentional degradation, not terminal overflow.
	compactBoardBreakpoint = 48

	// narrowNoticeBreakpoint is the smallest width at which a framed board can
	// carry its border and padding. The notice below it is itself fitted to the
	// terminal, so even an unusually small reported width stays safe.
	narrowNoticeBreakpoint = borderCells + paddingCells + 2
)

func terminalWidthOrOne(width int) int {
	if width < 1 {
		return 1
	}
	return width
}

// truncateCells is the one-cell-width truncation boundary used by board
// labels. x/ansi measures terminal cells and grapheme clusters, so CJK,
// emoji, and combining marks are never split in the middle.
func truncateCells(text string, width int) string {
	if width <= 0 {
		return ""
	}
	return xansi.Truncate(text, width, "…")
}

// fitLine keeps a non-wrapping surface to one terminal row. Content surfaces
// such as Next and detail deliberately wrap; headers, hints, and status lines
// use this helper because a focus change must not change their geometry.
func fitLine(text string, width int) string {
	return truncateCells(strings.ReplaceAll(text, "\n", " "), terminalWidthOrOne(width))
}
