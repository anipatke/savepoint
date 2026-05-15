package board

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/styles"
)

const (
	glyphBuild   = "▣"
	glyphTest    = "◇"
	glyphAudit   = "◆"
	glyphWarning = "⚠"

	cardOverhead = 4 // border (2) + padding (2×1)
)

// RenderCard renders a task card with stage glyph, truncated ID+title, and focus styling.
// When router state matches t's release/epic/task, a green priority glyph replaces the stage glyph.
// defectMarker is an optional compact string (e.g. "! D003") appended to the id line when width allows.
func RenderCard(t data.Task, width int, focused bool, routerState *data.RouterState, defectMarker string) string {
	inner := width - cardOverhead
	if inner < 2 {
		inner = 2
	}

	glyph := taskGlyph(t, routerState)
	stage := taskStageText(t)
	idWidth := inner - 2
	if stage != "" {
		idWidth -= lipgloss.Width(stage) + 1
	}
	if idWidth < 1 {
		idWidth = 1
	}

	idLine := fmt.Sprintf("%s %s", glyph, truncate(shortID(t.ID), idWidth))
	if stage != "" && lipgloss.Width(idLine)+1+lipgloss.Width(stage) <= inner {
		idLine += " " + stage
	}
	if defectMarker != "" {
		markerStyled := styles.CardMeta.Render(defectMarker)
		if lipgloss.Width(idLine)+1+lipgloss.Width(markerStyled) <= inner {
			idLine += " " + markerStyled
		}
	}
	titleLine := styles.CardMeta.Render(strings.Join(WrapText(t.Title, inner), "\n"))

	content := idLine + "\n" + titleLine
	if t.ComplexityTier != "" {
		content += "\n" + styles.CardMeta.Render(complexityCardLabel(t.ComplexityTier))
	}

	if focused {
		return styles.CardFocused.Width(width).Render(content)
	}
	return styles.Card.Width(width).Render(content)
}

func taskGlyph(t data.Task, routerState *data.RouterState) string {
	if t.Column == data.ColumnInProgress {
		return stageGlyphStyled(t.Stage)
	}
	if t.Column == data.ColumnDone {
		return styles.GlyphBuild.Render(glyphBuild)
	}
	if isRouterPriority(t, routerState) {
		return styles.TagDone.Render(glyphBuild)
	}
	return stageGlyphStyled(t.Stage)
}

func taskStageText(t data.Task) string {
	switch t.Column {
	case data.ColumnInProgress:
		return styles.CardMeta.Render(strings.ToUpper(stageLabel(t.Stage)))
	case data.ColumnDone:
		return styles.CardMeta.Render("DONE")
	default:
		return ""
	}
}

func stageGlyphStyled(stage data.ProgressStage) string {
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
	if state.State == "defect-building" {
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

func complexityCardLabel(tier data.ComplexityTier) string {
	switch tier {
	case data.ComplexityLow:
		return "· low"
	case data.ComplexityMedium:
		return "· med"
	case data.ComplexityHigh:
		return "· high"
	case data.ComplexitySpike:
		return "· spike"
	default:
		return "· " + string(tier)
	}
}
