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

// SplitLongWord splits a long word into chunks of at most width runes.
func SplitLongWord(word string, width int) []string {
	runes := []rune(word)
	lines := []string{}
	for len(runes) > width {
		lines = append(lines, string(runes[:width]))
		runes = runes[width:]
	}
	if len(runes) > 0 {
		lines = append(lines, string(runes))
	}
	return lines
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
