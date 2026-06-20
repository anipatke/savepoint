package board

import (
	"strings"
	"testing"
)

// maxLineWidth returns the widest rune-count among lines.
func maxLineWidth(lines []string) int {
	w := 0
	for _, l := range lines {
		if n := len([]rune(l)); n > w {
			w = n
		}
	}
	return w
}

func TestWrapText_wrapsAtWordBoundaries(t *testing.T) {
	got := WrapText("the quick brown fox jumps", 10)
	want := []string{"the quick", "brown fox", "jumps"}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Errorf("WrapText = %q, want %q", got, want)
	}
}

func TestWrapText_neverExceedsWidth(t *testing.T) {
	got := WrapText("Resolve the defect v1.2/D017-epic-view-line-wrapping before relying on internal/board/epic_panel_test.go coverage.", 30)
	if w := maxLineWidth(got); w > 30 {
		t.Errorf("WrapText produced line of width %d > 30: %q", w, got)
	}
}

// TestSplitLongWord_breaksAtNaturalSeparators is the D017 regression: long path
// and identifier tokens must break after structural separators, not mid-token.
func TestSplitLongWord_breaksAtNaturalSeparators(t *testing.T) {
	cases := []struct {
		name  string
		word  string
		width int
		want  []string
	}{
		{
			name:  "task/defect id breaks after hyphen",
			word:  "v1.2/D017-epic-view-line-wrapping",
			width: 26,
			want:  []string{"v1.2/D017-epic-view-line-", "wrapping"},
		},
		{
			name:  "file path breaks after underscore/slash",
			word:  "internal/board/epic_panel_test.go",
			width: 26,
			want:  []string{"internal/board/epic_panel_", "test.go"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := SplitLongWord(tc.word, tc.width)
			if strings.Join(got, "|") != strings.Join(tc.want, "|") {
				t.Errorf("SplitLongWord(%q, %d) = %q, want %q", tc.word, tc.width, got, tc.want)
			}
			if w := maxLineWidth(got); w > tc.width {
				t.Errorf("SplitLongWord(%q, %d) width %d exceeds limit", tc.word, tc.width, w)
			}
		})
	}
}

func TestSplitLongWord_hardCutsUnbreakableRun(t *testing.T) {
	// No natural separators: must still fall back to a rune-count cut within width.
	got := SplitLongWord("aaaaaaaaaaaaaaaaaa", 6)
	for _, line := range got {
		if len([]rune(line)) > 6 {
			t.Errorf("unbreakable run not cut to width: %q", got)
		}
	}
	if strings.Join(got, "") != "aaaaaaaaaaaaaaaaaa" {
		t.Errorf("SplitLongWord lost characters: %q", got)
	}
}

func TestSplitLongWord_packsSeparatorSegments(t *testing.T) {
	// Segments that fit together stay on one line up to the width limit.
	got := SplitLongWord("a/b/c/d/e/f/g/h", 6)
	if w := maxLineWidth(got); w > 6 {
		t.Errorf("packed line exceeds width: %q", got)
	}
	if strings.Join(got, "") != "a/b/c/d/e/f/g/h" {
		t.Errorf("SplitLongWord lost characters: %q", got)
	}
}
