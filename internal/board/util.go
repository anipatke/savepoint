package board

import "strings"

// sliceIndex returns the index of target in items, or 0 if not found.
func sliceIndex(items []string, target string) int {
	for i, e := range items {
		if e == target {
			return i
		}
	}
	return 0
}

// WrapText wraps s to fit within width, splitting on word boundaries.
func WrapText(s string, width int) []string {
	if width < 4 {
		width = 4
	}
	words := strings.Fields(s)
	if len(words) == 0 {
		return nil
	}
	lines := []string{}
	current := ""
	for _, word := range words {
		if len([]rune(word)) > width {
			if current != "" {
				lines = append(lines, current)
				current = ""
			}
			lines = append(lines, SplitLongWord(word, width)...)
			continue
		}
		if current == "" {
			current = word
			continue
		}
		if len([]rune(current))+1+len([]rune(word)) <= width {
			current += " " + word
			continue
		}
		lines = append(lines, current)
		current = word
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

// SplitLongWord splits a word too long for width into chunks that fit. It
// prefers breaking after structural separators (path, identifier, and version
// delimiters) so tokens like "internal/board/epic_panel_test.go" or
// "v1.2/D017-epic-view-line-wrapping" break at readable boundaries instead of
// mid-token. A run with no separator that still exceeds width is hard-cut by
// rune count as a last resort.
func SplitLongWord(word string, width int) []string {
	if width < 1 {
		width = 1
	}
	lines := []string{}
	current := ""
	for _, seg := range splitNaturalSegments(word) {
		segRunes := []rune(seg)
		if len([]rune(current))+len(segRunes) <= width {
			current += seg
			continue
		}
		if current != "" {
			lines = append(lines, current)
			current = ""
		}
		for len(segRunes) > width {
			lines = append(lines, string(segRunes[:width]))
			segRunes = segRunes[width:]
		}
		current = string(segRunes)
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

// splitNaturalSegments breaks word into segments that each end immediately after
// a structural separator (the separator stays attached to its left fragment),
// plus any trailing remainder. A word with no separators yields itself.
func splitNaturalSegments(word string) []string {
	runes := []rune(word)
	segments := []string{}
	start := 0
	for i, r := range runes {
		if naturalBreakAfter(r) {
			segments = append(segments, string(runes[start:i+1]))
			start = i + 1
		}
	}
	if start < len(runes) {
		segments = append(segments, string(runes[start:]))
	}
	return segments
}

// naturalBreakAfter reports whether a long word may break immediately after r.
// These are path, identifier, and version-string separators where a wrap keeps
// the fragment readable.
func naturalBreakAfter(r rune) bool {
	switch r {
	case '/', '\\', '-', '_', '.', ':':
		return true
	}
	return false
}

// truncate clips s to max runes, appending "…" if clipped.
func truncate(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	if max <= 1 {
		return "…"
	}
	return string(runes[:max-1]) + "…"
}
