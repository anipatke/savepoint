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
	extra := extraHeaderLines(nextActivity)
	statusLines := wrappedStatusLines(m.StatusMessage, w)
	footerExtra := len(statusLines) - 1
	layout := CalculateLayoutWithChrome(w, h, extra)
	topDivider := dividerLine(w)
	board := m.renderBoard(layout)
	boardBudget := h - 8 - extra - footerExtra
	if boardBudget < 0 {
		boardBudget = 0
	}
	boardLines := strings.Split(board, "\n")
	if len(boardLines) > boardBudget {
		board = strings.Join(boardLines[:boardBudget], "\n")
	}
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

	if m.Overlay == OverlayReleaseDocs {
		ow := overlayWidth(w)
		overlay := RenderReleaseDocs(m.ReleaseDocs, m.ReleaseDocIndex, ow, detailMaxHeight(h), m.selectedReleaseDocOffset())
		return overlayOnBase(dimLines(base), overlay, w, h)
	}

	if m.Overlay == OverlayAudit {
		ow := overlayWidth(w)
		overlay := RenderAuditOverlay(m.Audit, m.AuditTab, m.FindingCursor, ow, detailMaxHeight(h), m.selectedAuditOffset())
		return overlayOnBase(dimLines(base), overlay, w, h)
	}

	if m.Overlay == OverlayFindingDetail {
		if finding, ok := m.activeFinding(); ok {
			ow := overlayWidth(w)
			detail := RenderFindingDetail(finding, ow, detailMaxHeight(h), m.FindingDetailOffset)
			return overlayOnBase(dimLines(base), detail, w, h)
		}
		return base
	}

	if m.Overlay == OverlayDefect {
		defects := defectsForOverlay(m.AllDefects, m.SelectedRelease)
		overlay := RenderDefectsOverlay(defects, m.DefectCursor, min(overlayWidth(w), 60))
		return overlayOnBase(dimLines(base), overlay, w, h)
	}

	if m.Overlay == OverlayDefectDetail {
		defects := defectsForOverlay(m.AllDefects, m.SelectedRelease)
		if m.DefectCursor >= 0 && m.DefectCursor < len(defects) {
			ow := overlayWidth(w)
			detail := RenderDefectDetail(defects[m.DefectCursor], ow, detailMaxHeight(h), m.DefectDetailOffset)
			return overlayOnBase(dimLines(base), detail, w, h)
		}
		return base
	}

	if m.Overlay == OverlayDetail {
		task, ok := m.focusedTask()
		if !ok {
			return base
		}
		ow := overlayWidth(w)
		detail := RenderDetail(task, ow, m.RouterState, detailMaxHeight(h), m.DetailOffset, m.focusedTaskFindings(), m.LinkedFindingCursor)
		return overlayOnBase(dimLines(base), detail, w, h)
	}

	if m.Overlay == OverlayEpicDetail {
		ow := overlayWidth(w)
		epicSlug := m.epicDetailEpic()
		var detail string
		switch m.EpicDetailTab {
		case 1:
			detail = RenderEpicAuditTab(epicSlug, m.EpicAuditContent, ow, detailMaxHeight(h), m.EpicDetailOffset, m.EpicDetailTab)
		default:
			detail = RenderEpicDetail(epicSlug, m.EpicDetailContent, ow, detailMaxHeight(h), m.EpicDetailOffset, m.EpicDetailTab, m.openEpicFindings(), m.LinkedFindingCursor)
		}
		return overlayOnBase(dimLines(base), detail, w, h)
	}

	return base
}

func (m Model) renderHeader(w int) string {
	icon := styles.HeaderIcon.Render("▣")
	text := styles.HeaderText.Render("S A V E P O I N T")
	left := icon + "  " + text

	var parts []string
	if m.SelectedRelease != "" {
		parts = append(parts, styles.HeaderRelease.Render(m.SelectedRelease))
	}
	if count := m.openDefectCount(); count > 0 {
		parts = append(parts, styles.HeaderRight.Render(fmt.Sprintf("⚠  %d open", count)))
	}

	if len(parts) > 0 {
		right := strings.Join(parts, styles.HeaderRight.Render(" │ "))
		inner := w - 2 // HeaderFrame padding(1,1)
		gap := inner - lipgloss.Width(left) - lipgloss.Width(right)
		if gap > 0 {
			return styles.HeaderFrame.Width(w).Render(left + strings.Repeat(" ", gap) + right)
		}
	}
	return styles.HeaderFrame.Width(w).Render(left)
}

func (m Model) openDefectCount() int {
	count := 0
	for _, d := range m.AllDefects {
		if m.SelectedRelease != "" && d.Release != m.SelectedRelease {
			continue
		}
		if d.Status != data.DefectResolved {
			count++
		}
	}
	return count
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
	case "defect-building":
		return "DEFECT", styles.FooterPhaseDefect, true
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
	case "defect-building":
		s = fmt.Sprintf("Defect %s", shortID(state.Defect))
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
		epic := m.renderEpicPanel(layout.EpicPanelWidth, layout.ContentHeight)
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

func (m Model) renderEpicPanel(w int, maxHeight int) string {
	return RenderEpicSidebar(m.Epics, m.SelectedEpic, w, m.EpicPanelFocus, m.EpicPanelCursor, m.EpicStatus, maxHeight)
}

func (m Model) renderColumn(col data.ColumnType, colW, maxHeight int) string {
	focused := !m.EpicPanelFocus && m.FocusedColumn == col
	markers := buildTaskMarkers(m.Tasks[col], m.AllDefects)
	return RenderColumn(m.Tasks[col], col, colW, maxHeight, m.ColumnOffsets[col], m.FocusedTask, focused, m.RouterState, markers)
}

// buildTaskMarkers returns a map from task ID to defect marker string for any
// task in the slice that is referenced by a defect.
func buildTaskMarkers(tasks []data.Task, defects []data.Defect) map[string]string {
	if len(defects) == 0 {
		return nil
	}
	markers := map[string]string{}
	for _, t := range tasks {
		if m := defectMarkerForTask(t.ID, defects); m != "" {
			markers[t.ID] = m
		}
	}
	return markers
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

// footerHints is the board's keyboard hint bar. Single-space separators keep the
// full hint set — including A:audits — on one line at the 80-column minimum
// supported width, where two-space separators would overflow and truncate hints.
const footerHints = "←→:nav p:prio ctrl+r:refresh R:release d:defects D:docs A:audits ?:help q:quit"

func (m Model) renderFooter(termW int) string {
	phase := footerLine(termW,
		styles.FooterPhasePlan.Render("PLAN")+
			styles.FooterDivider.Render(" │ ")+
			styles.FooterPhaseBuild.Render("BUILD")+
			styles.FooterDivider.Render(" │ ")+
			styles.FooterPhaseAudit.Render("AUDIT"),
	)
	hints := footerLine(termW, styles.FooterHints.Render(footerHints))
	statusLines := wrappedStatusLines(m.StatusMessage, termW)
	renderedStatus := make([]string, len(statusLines))
	for i, line := range statusLines {
		status := ""
		if line != "" {
			status = styles.StatusBar.Render(line)
		}
		renderedStatus[i] = footerLine(termW, status)
	}
	sections := append([]string{phase}, renderedStatus...)
	sections = append(sections, hints)
	return lipgloss.JoinVertical(lipgloss.Center, sections...)
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

const maxStatusFooterLines = 3

func wrappedStatusLines(message string, termW int) []string {
	if strings.TrimSpace(message) == "" {
		return []string{""}
	}
	if termW <= 0 {
		termW = defaultTermW
	}

	words := strings.Fields(message)
	lines := make([]string, 0, maxStatusFooterLines)
	current := ""
	for _, word := range words {
		for lipgloss.Width(word) > termW {
			part := xansi.Truncate(word, termW, "")
			lines = append(lines, part)
			if len(lines) == maxStatusFooterLines {
				lines[len(lines)-1] = xansi.Truncate(lines[len(lines)-1], max(termW-1, 1), "…")
				return lines
			}
			word = strings.TrimPrefix(word, part)
		}

		candidate := word
		if current != "" {
			candidate = current + " " + word
		}
		if lipgloss.Width(candidate) <= termW {
			current = candidate
			continue
		}
		lines = append(lines, current)
		if len(lines) == maxStatusFooterLines {
			lines[len(lines)-1] = xansi.Truncate(lines[len(lines)-1], max(termW-1, 1), "…")
			return lines
		}
		current = word
	}

	if current != "" {
		lines = append(lines, current)
	}
	if len(lines) > maxStatusFooterLines {
		lines = lines[:maxStatusFooterLines]
		lines[len(lines)-1] = xansi.Truncate(lines[len(lines)-1], max(termW-1, 1), "…")
	}
	return lines
}
