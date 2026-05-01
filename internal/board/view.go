package board

import (
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

	layout := CalculateLayout(w, m.Height)
	icon := styles.HeaderIcon.Render("▣")
	text := styles.HeaderText.Render("S A V E P O I N T")
	header := styles.HeaderFrame.Width(w).Render(icon + "  " + text)
	board := m.renderBoard(layout)
	footer := m.renderFooter(w)
	base := lipgloss.JoinVertical(lipgloss.Left, header, board, footer)

	h := m.Height
	if h == 0 {
		h = defaultTermH
	}

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
		detail := RenderDetail(task, ow)
		return overlayOnBase(dimLines(base), detail, w, h)
	}

	return base
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
		return m.renderColumn(m.FocusedColumn, layout.ColWidths[0])
	}
	allCols := []data.ColumnType{data.ColumnPlanned, data.ColumnInProgress, data.ColumnDone}
	rendered := make([]string, len(allCols))
	for i, col := range allCols {
		rendered[i] = m.renderColumn(col, layout.ColWidths[i])
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, rendered...)
}

func (m Model) renderEpicPanel(w int) string {
	return RenderEpicSidebar(m.Epics, m.SelectedEpic, w)
}

func (m Model) renderColumn(col data.ColumnType, colW int) string {
	focused := m.FocusedColumn == col
	return RenderColumn(m.Tasks[col], col, colW, m.FocusedTask, focused)
}

func (m Model) renderFooter(termW int) string {
	phase := footerLine(termW,
		styles.FooterPhasePlan.Render("PLAN")+
			styles.FooterDivider.Render(" │ ")+
			styles.FooterPhaseBuild.Render("BUILD")+
			styles.FooterDivider.Render(" │ ")+
			styles.FooterPhaseAudit.Render("AUDIT"),
	)
	hints := footerLine(termW, styles.FooterHints.Render("←/→:nav  E:epic  R:release  ?:help  q:quit"))
	spacer := footerLine(termW, "")
	return lipgloss.JoinVertical(lipgloss.Center, phase, spacer, hints)
}

func footerLine(termW int, content string) string {
	if termW <= 0 {
		termW = defaultTermW
	}
	if lipgloss.Width(content) > termW {
		content = xansi.Truncate(content, termW, "")
	}
	return lipgloss.NewStyle().Width(termW).Align(lipgloss.Center).Render(content)
}
