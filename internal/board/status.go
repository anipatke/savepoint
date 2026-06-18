package board

import (
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/styles"
)

const statusGlyphDefault = " "
const statusGlyphUnknown = "?"

func statusGlyph(status string) string {
	switch status {
	case string(data.ColumnPlanned):
		return styles.CardMeta.Render("○")
	case string(data.ColumnInProgress):
		return styles.GlyphBuild.Render("▶")
	case string(data.ColumnDone):
		return styles.TagDone.Render("◉")
	case string(data.EpicStatusAudited):
		return styles.TagDone.Render("✓")
	case "":
		return statusGlyphDefault
	default:
		// Defense in depth: a status that escaped load-time normalization
		// still renders a visible glyph rather than a blank cell.
		return styles.CardMeta.Render(statusGlyphUnknown)
	}
}
