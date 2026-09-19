package v2

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/styles"
)

const (
	// defaultTermW and defaultTermH are the size the board lays out against
	// before a terminal reports its own.
	defaultTermW = 80
	defaultTermH = 24

	// sidebarWidth is the outer width the Objective sidebar occupies once it
	// exists. Reserving it here — and only at a width where three columns still
	// fit beside it — keeps column geometry from moving when the sidebar
	// arrives.
	sidebarWidth = 28
	// sidebarBreakpoint is the narrowest terminal that carries the sidebar
	// alongside three columns.
	sidebarBreakpoint = 120

	// minColumnHeight is the shortest column the board will draw rather than
	// give the columns no room at all.
	minColumnHeight = 6
)

// diagnosticHeading is the load-diagnostic screen's title. It names the class
// of problem — the project's data, not the board — so the reader knows which
// thing to go and fix.
const diagnosticHeading = "INVALID PROJECT DATA"

// columnLabels is the three-column vocabulary, keyed by the recorded Task
// status it counts. The columns stay Planned/In Progress/Done: implementation
// stage, clearance, owner wait, and replan are badges, not columns.
var columnLabels = []struct {
	Label  string
	Status data.ColumnType
}{
	{"PLANNED", data.ColumnPlanned},
	{"IN PROGRESS", data.ColumnInProgress},
	{"DONE", data.ColumnDone},
}

func (m Model) View() string {
	w, h := m.terminalWidth(), m.terminalHeight()

	if !m.Loaded {
		return styles.StatusBar.Render("Loading project…")
	}
	if m.Diagnostic != "" {
		return m.renderDiagnostic(w)
	}
	if w < narrowNoticeBreakpoint {
		return m.renderNarrowNotice(w)
	}
	base := m.renderBoard(w, h)
	if m.ReleaseOverlay {
		return m.renderReleaseOverlay(base, w, h)
	}
	return base
}

// terminalWidth and terminalHeight are the size every surface is laid out
// against, falling back to a conventional 80×24 before the first
// WindowSizeMsg arrives.
func (m Model) terminalWidth() int {
	if m.Width <= 0 {
		return defaultTermW
	}
	return m.Width
}

func (m Model) terminalHeight() int {
	if m.Height <= 0 {
		return defaultTermH
	}
	return m.Height
}

// renderDiagnostic is the whole screen for a project that did not load: the
// heading, the diagnostic naming the file and the problem, and why nothing else
// is shown. It draws no columns and no counts, because there is no loaded
// project to count — a board here would be showing something other than the
// project.
func (m Model) renderDiagnostic(w int) string {
	body := lipgloss.NewStyle().Width(w).Render(m.Diagnostic)
	explain := lipgloss.NewStyle().Width(w).Render("No board is drawn: this project's records did not load.")
	return lipgloss.JoinVertical(lipgloss.Left,
		styles.HeaderFrame.Width(w).Render(styles.HeaderIcon.Render("▣")+"  "+styles.HeaderText.Render(diagnosticHeading)),
		body,
		"",
		explain,
		"",
		styles.FooterHints.Render("q:quit"),
	)
}

func (m Model) renderNarrowNotice(w int) string {
	width := terminalWidthOrOne(w)
	return fitLine("Terminal too narrow; widen to continue", width)
}

// renderBoard assembles the surfaces around the columns, then gives the columns
// whatever height is left. Measuring the chrome rather than subtracting a
// constant means a surface above the board — the migration line, or a Next
// area that grew a line because a clearance phrase wrapped — cannot silently
// push the columns past the bottom of the terminal.
func (m Model) renderBoard(w, h int) string {
	above, below := m.boardChrome(w)
	sections := append(above, m.renderBody(w, boardBodyHeight(h, above, below)))
	sections = append(sections, below...)
	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// boardChrome is everything drawn above and below the board's body. It is built
// rather than measured by a constant so a surface that grew a line — the
// migration line, or a Next area whose clearance phrase wrapped — cannot
// silently push the body past the bottom of the terminal.
func (m Model) boardChrome(w int) (above, below []string) {
	above = []string{m.renderHeader(w), m.renderSelection(w)}
	if m.ReloadDiagnostic != "" {
		above = append(above, m.renderReloadDiagnostic(w))
	}
	if m.State.Migration.Pending {
		above = append(above, m.renderMigration(w))
	}
	above = append(above, m.renderNext(w), styles.Divider.Render(strings.Repeat("─", w)))
	below = []string{styles.Divider.Render(strings.Repeat("─", w)), m.renderStatusBar(w)}
	return above, below
}

// renderReloadDiagnostic keeps the last good board visible while making the
// failed external load explicit. The next successful load removes this line.
func (m Model) renderReloadDiagnostic(w int) string {
	return lipgloss.NewStyle().Width(w).Render(
		styles.HeaderIcon.Render("RELOAD") + "  " + styles.StatusBar.Render(m.ReloadDiagnostic),
	)
}

// boardBodyHeight is whatever height the chrome left for the body.
func boardBodyHeight(h int, above, below []string) int {
	height := h - renderedLines(above) - renderedLines(below)
	if height < minColumnHeight {
		return minColumnHeight
	}
	return height
}

// renderBody is the columns, or the open detail in their place. The overlay
// takes the columns' region rather than the whole screen, so the header, the
// Next answer, and the status bar stay where they were while a record is read.
func (m Model) renderBody(w, height int) string {
	if m.Help {
		return renderHelp(m, w, height)
	}
	if m.Issues != nil {
		if m.Issues.Detail != nil {
			return renderIssueDetail(*m.Issues.Detail, w, height, m.Issues.DetailOffset)
		}
		return renderIssues(m, w, height)
	}
	if m.Detail != nil {
		return renderDetail(*m.Detail, w, height, m.DetailOffset)
	}
	return m.renderColumns(w, height)
}

// detailViewport is the region an open overlay is drawn in. Update sizes its
// scrolling through the same function View draws with, so a scroll key and the
// window it scrolls can never disagree about how much fits.
func (m Model) detailViewport() (width, height int) {
	w, h := m.terminalWidth(), m.terminalHeight()
	above, below := m.boardChrome(w)
	return w, boardBodyHeight(h, above, below)
}

// renderedLines counts the terminal lines a set of rendered sections occupies.
func renderedLines(sections []string) int {
	total := 0
	for _, section := range sections {
		total += strings.Count(section, "\n") + 1
	}
	return total
}

// renderHeader carries the proof that the index reached the model: the counts
// of the Objectives and Tasks this load actually indexed.
func (m Model) renderHeader(w int) string {
	left := styles.HeaderIcon.Render("▣") + "  " + styles.HeaderText.Render("S A V E P O I N T")
	right := styles.HeaderRight.Render(fmt.Sprintf("%d objectives · %d tasks", m.State.objectiveCount(), m.State.taskCount()))

	inner := w - 2 // HeaderFrame padding(1,1)
	if inner < 1 {
		return fitLine(left, w)
	}
	gap := inner - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 1 {
		return styles.HeaderFrame.Width(w).Render(truncateCells(left, inner))
	}
	return styles.HeaderFrame.Width(w).Render(left + strings.Repeat(" ", gap) + right)
}

// renderSelection states the optional Release context and which Objective the
// columns are filtered to. Both are navigation state; the Next area's answer
// remains the load command's shared projection.
func (m Model) renderSelection(w int) string {
	text := "Objective: none selected"
	if m.SelectedRelease != "" {
		text = "Release: " + m.SelectedRelease
		if release := m.selectedReleaseRecord(); release != nil {
			text += " — " + release.Title
		}
		text += " · Objective: none selected"
	}
	if m.SelectedObjective != "" {
		text = "Objective: " + m.SelectedObjective
		if objective := m.selectedObjectiveRecord(); objective != nil {
			text += " — " + objective.Title
		}
		if m.SelectedRelease != "" {
			text = "Release: " + m.SelectedRelease
			if release := m.selectedReleaseRecord(); release != nil {
				text += " — " + release.Title
			}
			text += " · "
			text += "Objective: " + m.SelectedObjective
			if objective := m.selectedObjectiveRecord(); objective != nil {
				text += " — " + objective.Title
			}
		}
		// The header counts the whole project, so a filtered board says how
		// much of it the columns are showing.
		text += fmt.Sprintf(" · %d of %d tasks", m.cardCount(), m.State.taskCount())
	}
	return styles.RootLine.Width(w).Render(styles.CardMeta.Render(fitLine(text, w)))
}

func (m Model) selectedObjectiveRecord() *data.ObjectiveV2 {
	if m.State.Index == nil || m.SelectedObjective == "" {
		return nil
	}
	return m.State.Index.Objectives[m.SelectedObjective]
}

func (m Model) selectedReleaseRecord() *data.ReleaseV2 {
	if m.State.Index == nil || m.SelectedRelease == "" {
		return nil
	}
	return m.State.Index.Releases[m.SelectedRelease]
}

// renderMigration reports an incomplete conversion with the recovery guidance
// migrate.PendingOperation produced, rather than a board's own summary of it.
func (m Model) renderMigration(w int) string {
	return lipgloss.NewStyle().Width(w).Render(styles.HeaderIcon.Render("MIGRATION") + "  " + m.State.MigrationGuidance)
}

// renderColumns draws the three columns side by side, each holding the cards
// whose recorded status it names.
func (m Model) renderColumns(w, height int) string {
	width := columnWidth(w)
	rendered := make([]string, 0, len(columnLabels)+1)
	if w < compactBoardBreakpoint {
		label := columnLabel(m.FocusedColumn)
		return renderColumn(
			label,
			m.Cards[m.FocusedColumn],
			width,
			height,
			columnCursor{Card: m.FocusedCard, Holds: true, Focused: true},
		)
	}
	if m.sidebarVisible() {
		rendered = append(rendered, renderSidebar(
			m.Objectives,
			m.SelectedObjective,
			m.ObjectiveCursor,
			m.SidebarFocused,
			sidebarWidth,
			height,
		))
	}
	for _, column := range columnLabels {
		rendered = append(rendered, renderColumn(
			column.Label,
			m.Cards[column.Status],
			width,
			height,
			columnCursor{
				Card:    m.FocusedCard,
				Holds:   m.FocusedColumn == column.Status,
				Focused: !m.SidebarFocused,
			},
		))
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, rendered...)
}

// columnWidth splits the terminal three ways, after setting aside the width the
// Objective sidebar occupies on a terminal wide enough to carry it.
func columnWidth(termW int) int {
	termW = terminalWidthOrOne(termW)
	if termW < compactBoardBreakpoint {
		return termW
	}
	available := termW
	if termW >= sidebarBreakpoint {
		available -= sidebarWidth
	}
	width := available / len(columnLabels)
	if width < minColumnContent+columnChrome {
		return minColumnContent + columnChrome
	}
	return width
}

func (m Model) renderStatusBar(w int) string {
	if strings.TrimSpace(m.StatusMessage) != "" {
		return styles.RootLine.Width(w).Render(styles.StatusBar.Render(fitLine(m.StatusMessage, w)))
	}
	return styles.RootLine.Width(w).Render(styles.FooterHints.Render(fitLine(m.hints(), w)))
}

func columnLabel(status data.ColumnType) string {
	for _, column := range columnLabels {
		if column.Status == status {
			return column.Label
		}
	}
	return columnLabels[0].Label
}

// hints name the keys that exist on the surface holding focus. A terminal too
// narrow to draw the sidebar is not offered the key that would focus it.
func (m Model) hints() string {
	switch {
	case m.Help:
		return "esc:close help  q:quit"
	case m.ReleaseOverlay:
		return "↑↓ / j k:release  enter:select  esc/q:cancel"
	case m.Issues != nil && m.Issues.Detail != nil:
		return "↑↓:scroll  enter:canonical  esc:back  q:quit"
	case m.Issues != nil:
		return "↑↓:issues  f:filter  enter:open  I:task issues  esc:close  q:quit"
	case m.Detail != nil:
		return joinHints("↑↓:scroll  esc:close", m.focusedActionText(), "?:help  q:quit")
	case !m.sidebarVisible():
		return joinHints("↑↓←→:card  r:releases  enter:detail", m.focusedActionText(), "?:help  q:quit")
	case m.SidebarFocused:
		return joinHints("↑↓:objective  r:releases  enter:select  v:detail  esc:clear  tab:cards", m.focusedActionText(), "?:help  q:quit")
	default:
		return joinHints("↑↓←→:card  r:releases  enter:detail  tab:objectives", m.focusedActionText(), "?:help  q:quit")
	}
}

func joinHints(parts ...string) string {
	var nonEmpty []string
	for _, part := range parts {
		if strings.TrimSpace(part) != "" {
			nonEmpty = append(nonEmpty, part)
		}
	}
	return strings.Join(nonEmpty, "  ")
}
