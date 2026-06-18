package board

import (
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/data"
)

func TestStatusGlyph_canonicalStatusesRenderNonBlank(t *testing.T) {
	for _, status := range []string{
		string(data.ColumnPlanned),
		string(data.ColumnInProgress),
		string(data.ColumnDone),
		string(data.EpicStatusAudited),
	} {
		if strings.TrimSpace(statusGlyph(status)) == "" {
			t.Fatalf("statusGlyph(%q) rendered blank, want a visible glyph", status)
		}
	}
}

func TestStatusGlyph_unknownStatusFallsBackToVisibleGlyph(t *testing.T) {
	got := statusGlyph("epic-design")
	if strings.TrimSpace(got) == "" {
		t.Fatal("statusGlyph(unknown) rendered blank, want a visible fallback glyph")
	}
	if !strings.Contains(got, statusGlyphUnknown) {
		t.Fatalf("statusGlyph(unknown) = %q, want it to contain %q", got, statusGlyphUnknown)
	}
}

func TestStatusGlyph_emptyStatusStaysBlank(t *testing.T) {
	if got := statusGlyph(""); got != statusGlyphDefault {
		t.Fatalf("statusGlyph(\"\") = %q, want the blank default %q", got, statusGlyphDefault)
	}
}
