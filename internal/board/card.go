package board

import (
	"fmt"
	"strings"

	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/styles"
)

const (
	glyphBuild = "▣"
	glyphTest  = "◇"
	glyphAudit = "◆"

	cardOverhead = 4 // border (2) + padding (2×1)
)

// RenderCard renders a task card with phase glyph, truncated ID+title, and focus styling.
// When router state matches t's release/epic/task, a green priority glyph replaces the phase glyph.
func RenderCard(t data.Task, width int, focused bool, routerState *data.RouterState) string {
	inner := width - cardOverhead
	if inner < 2 {
		inner = 2
	}

	var glyph string
	if t.Status != "" {
		glyph = statusGlyph(t.Status)
	} else if t.Column == data.ColumnDone {
		glyph = styles.GlyphBuild.Render(glyphBuild)
	} else if isRouterPriority(t, routerState) {
		glyph = styles.TagDone.Render(glyphBuild)
	} else {
		glyph = phaseGlyphStyled(t.Stage)
	}
	// glyph is 1 rune + 1 space prefix; leave room for "▣ "
	idLine := fmt.Sprintf("%s %s", glyph, truncate(shortID(t.ID), inner-2))
	titleLine := styles.CardMeta.Render(strings.Join(WrapText(t.Title, inner), "\n"))

	content := idLine + "\n" + titleLine

	if focused {
		return styles.CardFocused.Width(width).Render(content)
	}
	return styles.Card.Width(width).Render(content)
}

func phaseGlyphStyled(stage data.ProgressStage) string {
	switch stage {
	case data.StageTest:
		return styles.GlyphTest.Render(glyphTest)
	case data.StageAudit:
		return styles.GlyphAudit.Render(glyphAudit)
	default:
		return styles.GlyphBuild.Render(glyphBuild)
	}
}

func isRouterPriority(t data.Task, state *data.RouterState) bool {
	if state == nil || state.Task == "" {
		return false
	}
	if shortID(t.ID) != shortID(state.Task) {
		return false
	}
	if state.Release != "" && t.Release != "" && t.Release != state.Release {
		return false
	}
	routerEpic := state.Epic
	if routerEpic == "" {
		routerEpic = taskEpic(state.Task)
	}
	if routerEpic != "" && t.Epic != "" && shortID(t.Epic) != shortID(routerEpic) {
		return false
	}
	return true
}

func taskEpic(taskID string) string {
	if idx := strings.LastIndex(taskID, "/"); idx >= 0 {
		return taskID[:idx]
	}
	return ""
}

// shortID strips the epic prefix and slug from a task ID.
// "E06-atari-noir-layout/T004-component-refinement" → "T004"
func shortID(id string) string {
	if idx := strings.LastIndex(id, "/"); idx >= 0 {
		id = id[idx+1:]
	}
	if idx := strings.Index(id, "-"); idx >= 0 {
		id = id[:idx]
	}
	return id
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
