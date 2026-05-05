package board

import (
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/styles"
)

const statusGlyphDefault = " "

func statusGlyph(status string) string {
	switch status {
	case string(data.ColumnPlanned):
		return styles.CardMeta.Render("○")
	case string(data.ColumnInProgress):
		return styles.GlyphBuild.Render("▶")
	case string(data.ColumnDone):
		return styles.TagDone.Render("◉")
	case "audited":
		return styles.TagDone.Render("✓")
	default:
		return statusGlyphDefault
	}
}
