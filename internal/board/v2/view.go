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
	sidebarWidth = 34
	// sidebarBreakpoint is the narrowest terminal that carries the sidebar
	// alongside three columns at their existing 30-cell width.
	sidebarBreakpoint = 124

	// minColumnHeight is the shortest column the board will draw rather than
	// give the columns no room at all.
	minColumnHeight = 6

	// boardMarginX and boardMarginY are the outer breathing room between the
	// terminal's own edge and every surface the board draws — columns, header,
	// footer, overlays alike. They are reserved out of the terminal's reported
	// size before anything is laid out, then restored as padding around the
	// finished frame, so no surface's own geometry has to know about them.
	boardMarginX = 2
	boardMarginY = 1
)

// boardMargin is the padding applied once, around the whole assembled frame.
var boardMargin = lipgloss.NewStyle().Padding(boardMarginY, boardMarginX)

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

	var content string
	switch {
	case m.Diagnostic != "":
		content = m.renderDiagnostic(w)
	case w < narrowNoticeBreakpoint:
		content = m.renderNarrowNotice(w)
	default:
		content = m.renderBoard(w, h)
		if m.ReleaseOverlay {
			content = m.renderReleaseOverlay(content, w, h)
		}
	}
	return boardMargin.Render(content)
}

// terminalWidth and terminalHeight are the size every surface is laid out
// against, falling back to a conventional 80×24 before the first
// WindowSizeMsg arrives. Both are the *content* size: boardMarginX and
// boardMarginY are reserved out of the terminal's own reported size here, so
// every surface downstream lays out against room that already excludes the
// margin View restores around the finished frame.
func (m Model) terminalWidth() int {
	w := defaultTermW
	if m.Width > 0 {
		w = m.Width
	}
	return marginedDimension(w, boardMarginX)
}

func (m Model) terminalHeight() int {
	h := defaultTermH
	if m.Height > 0 {
		h = m.Height
	}
	return marginedDimension(h, boardMarginY)
}

// marginedDimension reserves margin cells off both ends of outer, clamped so
// a terminal too small to carry the full margin still gets at least one cell
// of content rather than a negative or zero size.
func marginedDimension(outer, margin int) int {
	inner := outer - margin*2
	if inner < 1 {
		return 1
	}
	return inner
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
// constant means a surface above the board — the reload diagnostic, or a Next
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
// reload diagnostic, or a Next area whose clearance phrase wrapped — cannot
// silently push the body past the bottom of the terminal.
func (m Model) boardChrome(w int) (above, below []string) {
	above = []string{m.renderHeader(w), m.renderSelection(w)}
	if notice := unassignedGoalNotice(m.State.Index); notice != "" {
		above = append(above, lipgloss.NewStyle().Width(w).Render(notice))
	}
	if m.ReloadDiagnostic != "" {
		above = append(above, m.renderReloadDiagnostic(w))
	}
	above = append(above, m.renderNext(w), styles.Divider.Render(strings.Repeat("─", w)))
	below = []string{styles.Divider.Render(strings.Repeat("─", w)), m.renderPhaseRow(w), "", m.renderStatusBar(w)}
	return above, below
}

// renderPhaseRow draws the four V2 router phases together, in the same
// colors the public site (getsavepoint.dev) uses for them: white Idea,
// purple Design, orange Task, green Check. It is a fixed reference row, not
// a status indicator — it does not read the router's own state — the same
// shape V1's board drew for its PLAN/BUILD/AUDIT phases.
func (m Model) renderPhaseRow(w int) string {
	row := styles.FooterPhaseIdea.Render("IDEA") +
		styles.FooterDivider.Render(" │ ") +
		styles.FooterPhaseDesign.Render("DESIGN") +
		styles.FooterDivider.Render(" │ ") +
		styles.FooterPhaseTask.Render("TASK") +
		styles.FooterDivider.Render(" │ ") +
		styles.FooterPhaseCheck.Render("CHECK")
	return styles.RootLine.Width(w).Align(lipgloss.Center).Render(row)
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

// renderSelection states the selected Goal context — a bold capitalized
// "GOAL:" label with the record's own ID and title in plain white after it —
// and, only when no Objective filter is in effect, the Goal-scoped
// allObjectivesLabel. The Objective the columns are filtered to is otherwise
// the sidebar's own purple-accented selection marker (glyphSelected), not
// restated here. Nothing here is truncated by fitLine: a styled line carries
// ANSI codes fitLine's rune count would miscount, so overflow is left to the
// terminal to wrap. With no valid Goal selected, this line is blank.
func (m Model) renderSelection(w int) string {
	var parts []string
	if m.SelectedRelease != "" {
		releaseText := m.SelectedRelease
		if r := m.selectedReleaseRecord(); r != nil {
			releaseText += " — " + r.Title
		}
		parts = append(parts, styles.HeaderWhiteBold.Render(strings.ToUpper(goalLabel)+":")+" "+styles.HeaderWhite.Render(releaseText))
	}
	if m.SelectedRelease != "" && m.SelectedObjective == "" && m.State.Index != nil {
		label := allObjectivesLabel + " IN GOAL " + m.SelectedRelease
		parts = append(parts, styles.HeaderWhiteBold.Render(label))
	}
	return styles.RootLine.Width(w).Render(strings.Join(parts, "  ·  "))
}

// allObjectivesLabel marks the unfiltered view within the selected Goal. No
// sidebar row is marked selected then, so without it the view could read as
// one Objective's Tasks with a broken filter.
const allObjectivesLabel = "ALL OBJECTIVES"

func (m Model) selectedReleaseRecord() *data.ReleaseV2 {
	if m.State.Index == nil || m.SelectedRelease == "" {
		return nil
	}
	return m.State.Index.Releases[m.SelectedRelease]
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
		return styles.RootLine.Width(w).Align(lipgloss.Center).Render(styles.StatusBar.Render(fitLine(m.StatusMessage, w)))
	}
	return styles.RootLine.Width(w).Align(lipgloss.Center).Render(styles.FooterHints.Render(fitLine(m.hints(), w)))
}

func columnLabel(status data.ColumnType) string {
	for _, column := range columnLabels {
		if column.Status == status {
			return column.Label
		}
	}
	return columnLabels[0].Label
}

// hints name the keys that exist on the surface holding focus, and only the
// ones actually available right now: a terminal too narrow to draw the
// sidebar is not offered the key that would focus it, a Release selector is
// not offered when no Release exists, opening a record is not offered on an
// empty surface, and clearing the Objective filter is not offered when
// nothing is selected. The full key map, including keys omitted here,
// remains in Help.
func (m Model) hints() string {
	switch {
	case m.Help:
		// esc and q both only close Help here; neither reaches the global quit.
		return "esc/q:close"
	case m.ReleaseOverlay:
		return "↑↓ / j k:Goal  enter:select  v:detail  esc/q:cancel"
	case m.Issues != nil && m.Issues.Detail != nil:
		if m.Issues.Detail.DuplicateTarget != nil {
			return "↑↓:scroll  enter:canonical  esc:back  q:quit"
		}
		return "↑↓:scroll  esc:back  q:quit"
	case m.Issues != nil:
		return "↑↓:issues  f:filter  enter:open  esc:close  q:quit"
	case m.Detail != nil:
		return joinHints("↑↓:scroll  esc:close", m.focusedActionText(), "?:help  q:quit")
	case !m.sidebarVisible():
		return joinHints("↑↓←→:card  space:advance  backspace:retreat  i:issues", m.releaseHint(), m.detailHint("enter:detail"), m.focusedActionText(), "?:help  q:quit")
	case m.SidebarFocused:
		return joinHints("↑↓:objective  →:cards", m.releaseHint(), m.detailHint("v:detail"), "i:issues", m.clearObjectiveHint(), m.focusedActionText(), "?:help  q:quit")
	default:
		return joinHints("↑↓←→:card  space:advance  backspace:retreat  i:issues", m.releaseHint(), m.detailHint("enter:detail"), m.focusedActionText(), "?:help  q:quit")
	}
}

// releaseHint offers the Goal selector only when a Goal actually
// exists to select.
func (m Model) releaseHint() string {
	if len(m.Releases) == 0 {
		return ""
	}
	return goalSelectorKey + ":" + strings.ToLower(goalsLabel)
}

// detailHint offers opening a record only when the focused surface actually
// has one under its cursor.
func (m Model) detailHint(label string) string {
	if !m.hasDetailTarget() {
		return ""
	}
	return label
}

// clearObjectiveHint offers clearing the Objective filter only when one is
// currently selected.
func (m Model) clearObjectiveHint() string {
	if m.SelectedObjective == "" {
		return ""
	}
	return "esc:clear"
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
