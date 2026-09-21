package v2

import (
	"strings"
	"unicode"

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

// wrapTitleLines wraps text into at most maxLines lines, each of terminal
// cell width at most width. If text fits on a single line, exactly one line
// is returned without trailing blank lines. If it exceeds maxLines, the final
// line is truncated with an ellipsis.
func wrapTitleLines(text string, width int, maxLines int) []string {
	if maxLines <= 0 || width <= 0 {
		return nil
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return []string{""}
	}
	if xansi.StringWidth(text) <= width {
		return []string{text}
	}
	if maxLines == 1 {
		return []string{truncateCells(text, width)}
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	wordIdx := 0

	for lineNum := 1; lineNum <= maxLines && wordIdx < len(words); lineNum++ {
		isLastLine := lineNum == maxLines
		if isLastLine {
			remaining := strings.Join(words[wordIdx:], " ")
			if xansi.StringWidth(remaining) <= width {
				lines = append(lines, remaining)
			} else {
				lines = append(lines, truncateCells(remaining, width))
			}
			break
		}

		firstWord := words[wordIdx]
		if xansi.StringWidth(firstWord) > width {
			cut := xansi.Truncate(firstWord, width, "")
			if cut == "" {
				runes := []rune(firstWord)
				cut = string(runes[:1])
			}
			lines = append(lines, cut)
			words[wordIdx] = strings.TrimPrefix(firstWord, cut)
			continue
		}

		current := firstWord
		wordIdx++
		for wordIdx < len(words) {
			next := words[wordIdx]
			if xansi.StringWidth(current+" "+next) <= width {
				current += " " + next
				wordIdx++
			} else {
				break
			}
		}
		lines = append(lines, current)
	}

	return lines
}

// fitLine keeps a non-wrapping surface to one terminal row. Content surfaces
// such as Next and detail deliberately wrap; headers, hints, and status lines
// use this helper because a focus change must not change their geometry.
func fitLine(text string, width int) string {
	return truncateCells(strings.ReplaceAll(text, "\n", " "), terminalWidthOrOne(width))
}

// stripTerminalControls makes redirected output plain even when authored
// record text contains YAML-escaped ANSI or other terminal control bytes.
// Newlines produced by the renderer remain structural; authored tabs become
// spaces and other controls are removed.
func stripTerminalControls(text string) string {
	text = xansi.Strip(text)
	return strings.Map(func(r rune) rune {
		switch {
		case r == '\n':
			return r
		case r == '\t':
			return ' '
		case unicode.IsControl(r):
			return -1
		default:
			return r
		}
	}, text)
}
