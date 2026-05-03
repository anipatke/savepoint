package board

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	xansi "github.com/charmbracelet/x/ansi"
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/styles"
)

const defaultTermH = 24

const defaultTermW = 80

func (m Model) View() string {
	w := m.Width
	if w == 0 {
		w = defaultTermW
	}
	m.Width = w
	h := m.Height
	if h == 0 {
		h = defaultTermH
	}
	m.Height = h

	header := m.renderHeader(w)
	nextActivity := m.renderNextActivityLine(w)
	layout := CalculateLayoutWithChrome(w, h, extraHeaderLines(nextActivity))
	topDivider := dividerLine(w)
	board := m.renderBoard(layout)
	bottomDivider := dividerLine(w)
	footer := m.renderFooter(w)
	sections := []string{header}
	if nextActivity != "" {
		sections = append(sections, nextActivity)
	}
	sections = append(sections, topDivider, board, bottomDivider, footer)
	base := lipgloss.JoinVertical(lipgloss.Left, sections...)

	if m.Overlay == OverlayEpic {
		overlay := RenderEpicDropdown(m.Epics, m.EpicCursor, min(40, w))
		return overlayOnBase(dimLines(base), overlay, w, h)
	}

	if m.Overlay == OverlayRelease {
		overlay := RenderReleaseDropdown(m.Releases, m.ReleaseCursor, min(40, w))
		return overlayOnBase(dimLines(base), overlay, w, h)
	}

	if m.Overlay == OverlayHelp {
		help := RenderHelp(overlayWidth(w))
		return overlayOnBase(dimLines(base), help, w, h)
	}

	if m.Overlay == OverlayDetail {
		task, ok := m.focusedTask()
		if !ok {
			return base
		}
		ow := overlayWidth(w)
		detail := RenderDetail(task, ow, m.RouterState, detailMaxHeight(h), m.DetailOffset)
		return overlayOnBase(dimLines(base), detail, w, h)
	}

	if m.Overlay == OverlayEpicDetail {
		ow := overlayWidth(w)
		epicSlug := m.epicDetailEpic()
		var detail string
		if m.EpicDetailTab == 1 {
			detail = RenderEpicAuditTab(epicSlug, m.EpicAuditContent, ow, detailMaxHeight(h), m.EpicDetailOffset, m.EpicDetailTab)
		} else {
			detail = RenderEpicDetail(epicSlug, m.EpicDetailContent, ow, detailMaxHeight(h), m.EpicDetailOffset, m.EpicDetailTab)
		}
		return overlayOnBase(dimLines(base), detail, w, h)
	}

	return base
}

func (m Model) renderHeader(w int) string {
	icon := styles.HeaderIcon.Render("▣")
	text := styles.HeaderText.Render("S A V E P O I N T")
	left := icon + "  " + text
	return styles.HeaderFrame.Width(w).Render(left)
}

func extraHeaderLines(line string) int {
	if line == "" {
		return 0
	}
	return 1
}

func (m Model) renderNextActivityLine(w int) string {
	if w <= 0 {
		w = defaultTermW
	}
	return renderNextActivityLine(m.RouterState, w)
}

func renderNextActivityLine(state *data.RouterState, w int) string {
	tag, style, ok := nextActivityPhase(state)
	if !ok || strings.TrimSpace(state.NextAction) == "" {
		return ""
	}

	content := style.Render(tag+":") + " " + state.NextAction
	if lipgloss.Width(content) > w {
		content = xansi.Truncate(content, w, "…")
	}
	return styles.RootLine.Width(w).Render(content)
}

func nextActivityPhase(state *data.RouterState) (string, lipgloss.Style, bool) {
	if state == nil {
		return "", lipgloss.Style{}, false
	}
	switch state.State {
	case "task-building":
		return "BUILD", styles.FooterPhaseBuild, true
	case "audit-pending":
		return "AUDIT", styles.FooterPhaseAudit, true
	case "pre-implementation", "epic-design", "epic-task-breakdown":
		return "PLAN", styles.FooterPhasePlan, true
	default:
		return "", lipgloss.Style{}, false
	}
}

// FormatNextActivity formats a compact activity string from router state.
// Returns empty string when state is nil. Result is capped at 20 visible chars.
func FormatNextActivity(state *data.RouterState) string {
	if state == nil {
		return ""
	}
	var s string
	switch state.State {
	case "task-building":
		s = fmt.Sprintf("Build %s %s/%s", state.Release, shortID(state.Epic), shortID(state.Task))
	case "audit-pending":
		s = fmt.Sprintf("Audit %s", shortID(state.Epic))
	case "epic-design":
		s = fmt.Sprintf("Design %s", shortID(state.Epic))
	case "epic-task-breakdown":
		s = fmt.Sprintf("Plan %s", shortID(state.Epic))
	case "pre-implementation":
		s = fmt.Sprintf("Planning %s", state.Release)
	default:
		s = state.State
	}
	return xansi.Truncate(s, 20, "…")
}



func (m Model) focusedTask() (data.Task, bool) {
	tasks := m.Tasks[m.FocusedColumn]
	if len(tasks) == 0 || m.FocusedTask >= len(tasks) {
		return data.Task{}, false
	}
	return tasks[m.FocusedTask], true
}

func overlayWidth(termW int) int {
	ow := termW - 4
	if ow > 80 {
		ow = 80
	}
	if ow < 20 {
		ow = 20
	}
	return ow
}

// dimLines applies faint ANSI styling to each line individually.
func dimLines(s string) string {
	dim := lipgloss.NewStyle().Faint(true)
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		lines[i] = dim.Render(l)
	}
	return strings.Join(lines, "\n")
}

// overlayOnBase places overlay centered on base, preserving base lines outside
// the overlay area and replacing the left portion of intersecting lines.
func overlayOnBase(base, overlay string, termW, termH int) string {
	baseLines := strings.Split(base, "\n")
	overlayLines := strings.Split(overlay, "\n")

	overlayH := len(overlayLines)
	overlayW := 0
	for _, l := range overlayLines {
		if lw := lipgloss.Width(l); lw > overlayW {
			overlayW = lw
		}
	}

	startY := (termH - overlayH) / 2
	if startY < 0 {
		startY = 0
	}
	startX := (termW - overlayW) / 2
	if startX < 0 {
		startX = 0
	}

	for len(baseLines) < termH {
		baseLines = append(baseLines, "")
	}

	result := make([]string, len(baseLines))
	for i, line := range baseLines {
		oi := i - startY
		if oi >= 0 && oi < overlayH {
			left := xansi.Truncate(line, startX, "")
			leftW := lipgloss.Width(left)
			if leftW < startX {
				left += strings.Repeat(" ", startX-leftW)
			}
			result[i] = left + overlayLines[oi]
		} else {
			result[i] = line
		}
	}
	return strings.Join(result, "\n")
}

func (m Model) renderBoard(layout Layout) string {
	cols := m.renderColumns(layout)
	var content string
	if layout.EpicPanelVisible {
		epic := m.renderEpicPanel(layout.EpicPanelWidth)
		content = lipgloss.JoinHorizontal(lipgloss.Top, epic, cols)
	} else {
		content = cols
	}
	return styles.BoardFrame.Width(m.Width).Render(content)
}

func (m Model) renderColumns(layout Layout) string {
	if layout.ColCount == 1 {
		return m.renderColumn(m.FocusedColumn, layout.ColWidths[0], layout.ContentHeight)
	}
	allCols := []data.ColumnType{data.ColumnPlanned, data.ColumnInProgress, data.ColumnDone}
	rendered := make([]string, len(allCols))
	for i, col := range allCols {
		rendered[i] = m.renderColumn(col, layout.ColWidths[i], layout.ContentHeight)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, rendered...)
}

func (m Model) renderEpicPanel(w int) string {
	return RenderEpicSidebar(m.Epics, m.SelectedEpic, w, m.EpicPanelFocus, m.EpicPanelCursor, m.EpicStatus)
}

func (m Model) renderColumn(col data.ColumnType, colW, maxHeight int) string {
	focused := !m.EpicPanelFocus && m.FocusedColumn == col
	return RenderColumn(m.Tasks[col], col, colW, maxHeight, m.ColumnOffsets[col], m.FocusedTask, focused, m.RouterState)
}

func detailMaxHeight(termH int) int {
	if termH <= 0 {
		termH = defaultTermH
	}
	h := termH * 7 / 10
	if h < 6 {
		h = 6
	}
	return h
}

func (m Model) renderFooter(termW int) string {
	phase := footerLine(termW,
		styles.FooterPhasePlan.Render("PLAN")+
			styles.FooterDivider.Render(" │ ")+
			styles.FooterPhaseBuild.Render("BUILD")+
			styles.FooterDivider.Render(" │ ")+
			styles.FooterPhaseAudit.Render("AUDIT"),
	)
	hints := footerLine(termW, styles.FooterHints.Render("←/→:nav  p: Priority  R:release  ?:help  q:quit"))
	status := ""
	if m.StatusMessage != "" {
		status = styles.StatusBar.Render(m.StatusMessage)
	}
	statusLine := footerLine(termW, status)
	return lipgloss.JoinVertical(lipgloss.Center, phase, statusLine, hints)
}

func dividerLine(termW int) string {
	if termW <= 0 {
		termW = defaultTermW
	}
	return styles.Divider.Render(strings.Repeat("─", termW))
}

func footerLine(termW int, content string) string {
	if termW <= 0 {
		termW = defaultTermW
	}
	if lipgloss.Width(content) > termW {
		content = xansi.Truncate(content, termW, "")
	}
	return styles.RootLine.Width(termW).Align(lipgloss.Center).Render(content)
}
