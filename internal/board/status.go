package board

import (
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/styles"
)

const statusGlyphDefault = " "

func statusGlyph(status string) string {
	switch status {
	case string(data.StatusPlanned):
		return styles.CardMeta.Render("○")
	case string(data.StatusInProgress):
		return styles.GlyphBuild.Render("▶")
	case string(data.StatusDone):
		return styles.TagDone.Render("◉")
	case string(data.StatusAudited):
		return styles.TagDone.Render("✓")
	default:
		return statusGlyphDefault
	}
}
