package v2

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/opencode/savepoint/internal/styles"
)

const (
	// columnChrome is the horizontal space a column's frame takes, the same
	// border and padding a card spends.
	columnChrome = borderCells + paddingCells
	// minColumnContent is the narrowest content a column will render at rather
	// than overflow; a width below it clamps instead of wrapping into its
	// neighbour.
	minColumnContent = 12
	// columnHeaderLines is the header plus the rule under it.
	columnHeaderLines = 2
)

// columnCursor is where the card cursor stands in relation to one column.
// Holds says the cursor is in this column; Focused says the columns are the
// surface the keys are driving. Keeping them apart is what lets focus move to
// the sidebar and back without scrolling a long column away from where the
// reader left it: the window follows Holds, and only the accent follows Focused.
type columnCursor struct {
	Card    int
	Holds   bool
	Focused bool
}

// accented reports whether this column wears the focus accent.
func (c columnCursor) accented() bool {
	return c.Holds && c.Focused
}

// highlights reports whether the card at i is the one under the cursor, on the
// focused surface.
func (c columnCursor) highlights(i int) bool {
	return c.accented() && i == c.Card
}

// renderColumn draws one column: its label and card count, a rule, and as many
// cards as the height budget fits, with an indicator for whatever is scrolled
// out of view above or below. width and height are the column's outer size.
//
// The window is anchored on the card under the cursor, so moving the cursor
// down a long column scrolls it rather than losing it. Nothing else decides
// what is visible: cards arrive already resolved and already grouped.
func renderColumn(label string, cards []TaskCard, width, height int, cursor columnCursor) string {
	textW := columnTextWidth(width)
	bodyH := columnBodyHeight(height)

	header := fmt.Sprintf("%s (%d)", label, len(cards))
	if cursor.accented() {
		if isPlannedColumn(label) {
			header = styles.ColumnTitleFocusedPlanned.Render(header)
		} else if isDoneColumn(label) {
			header = styles.ColumnTitleFocusedDone.Render(header)
		} else {
			header = styles.ColumnTitleFocused.Render(header)
		}
	} else {
		header = styles.ColumnTitle.Render(header)
	}
	lines := []string{header, styles.Divider.Render(strings.Repeat("─", textW))}

	if len(cards) == 0 {
		lines = append(lines, styles.CardMeta.Render("(empty)"))
		return frameColumn(lines, textW, bodyH, cursor.accented(), label)
	}

	rendered := make([]string, len(cards))
	heights := make([]int, len(cards))
	for i, card := range cards {
		rendered[i] = renderCard(card, textW, cursor.highlights(i))
		heights[i] = strings.Count(rendered[i], "\n") + 1
	}

	budget := bodyH - columnHeaderLines
	if budget < 1 {
		budget = 1
	}
	start, end := visibleWindow(heights, budget, focusedIndex(cursor.Holds, cursor.Card, len(cards)))

	if start > 0 {
		lines = append(lines, scrollIndicator("↑", start, "above"))
	}
	lines = append(lines, rendered[start:end]...)
	if end < len(cards) {
		lines = append(lines, scrollIndicator("↓", len(cards)-end, "more"))
	}

	return frameColumn(lines, textW, bodyH, cursor.accented(), label)
}

// frameColumn draws the column's frame around its lines at exactly the height
// it was budgeted: Height fills a short column out and MaxHeight clips one
// whose last card would otherwise push past the bottom of the board.
func frameColumn(lines []string, textW, bodyH int, focused bool, label ...string) string {
	return columnStyle(focused, label...).
		Width(textW + paddingCells).
		Height(bodyH).
		MaxHeight(bodyH + borderCells).
		Render(strings.Join(lines, "\n"))
}

// focusedIndex is the row the window must keep visible: the one under the
// cursor in the column or list that holds it, and the first row everywhere
// else, so a list without the cursor always shows its top rather than a
// position the user cannot see moving.
func focusedIndex(holdsCursor bool, cursor, total int) int {
	if !holdsCursor || cursor < 0 || cursor >= total {
		return 0
	}
	return cursor
}

// visibleWindow picks the run of cards to draw: as many as fit from the top,
// then slid forward until the card that must stay visible is inside it. It
// reserves a line for each scroll indicator it will need, so an indicator never
// pushes a card out of the frame it was measured for. At least one card is
// always returned — a card taller than the whole budget is clipped by the
// frame rather than dropped.
func visibleWindow(heights []int, budget, mustShow int) (int, int) {
	for start := 0; ; start++ {
		end := fitFrom(heights, start, budget)
		if end > mustShow || end == len(heights) {
			return start, end
		}
	}
}

// fitFrom reports how far a window starting at start reaches within budget,
// counting the indicator lines that window itself makes necessary.
func fitFrom(heights []int, start, budget int) int {
	used := 0
	if start > 0 {
		used++ // "n above"
	}
	end := start
	for end < len(heights) {
		below := 0
		if end+1 < len(heights) {
			below = 1 // "n more"
		}
		if used+heights[end]+below > budget {
			break
		}
		used += heights[end]
		end++
	}
	if end == start && start < len(heights) {
		end = start + 1
	}
	return end
}

// columnTextWidth is the width text and cards get inside a column's frame,
// clamped so a very narrow terminal renders a narrow column rather than one
// that overflows into its neighbour.
func columnTextWidth(width int) int {
	textW := width - columnChrome
	if textW < 1 {
		return 1
	}
	return textW
}

// columnBodyHeight is the height inside a column's frame: the budget the header,
// the rule, and the cards share.
func columnBodyHeight(height int) int {
	bodyH := height - borderCells
	if bodyH < columnHeaderLines+1 {
		return columnHeaderLines + 1
	}
	return bodyH
}

// columnStyle is the column frame in two accents. Like the card frame, the two
// differ in color only: same border, same padding, same width.
func columnStyle(focused bool, label ...string) lipgloss.Style {
	if focused {
		if len(label) > 0 {
			if isPlannedColumn(label[0]) {
				return styles.ColumnFocusedPlanned
			}
			if isDoneColumn(label[0]) {
				return styles.ColumnFocusedDone
			}
		}
		return styles.ColumnFocused
	}
	return styles.ColumnUnfocused
}

func isPlannedColumn(label string) bool {
	return strings.HasPrefix(label, "PLANNED")
}

func isDoneColumn(label string) bool {
	return strings.HasPrefix(label, "DONE")
}

func scrollIndicator(arrow string, count int, suffix string) string {
	return styles.ScrollIndicator.Render(fmt.Sprintf("%s %d %s", arrow, count, suffix))
}
