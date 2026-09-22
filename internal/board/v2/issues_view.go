package v2

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	xansi "github.com/charmbracelet/x/ansi"
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/resume"
	"github.com/opencode/savepoint/internal/styles"
)

// issuesHeaderLines is the one line the Issues overlay spends above its
// columns: title, optional Task scope, and the active type filter. It is
// deliberately shared across all three columns rather than repeated inside
// each, since it describes the whole overlay's state, not any one status.
const issuesHeaderLines = 1

func renderIssues(model Model, width, height int) string {
	overlay := *model.Issues
	header := fitLine(issuesHeaderLine(overlay), width)

	bodyHeight := height - issuesHeaderLines
	if bodyHeight < minColumnHeight {
		bodyHeight = minColumnHeight
	}

	rows := model.filteredIssueRows()
	if len(rows) == 0 {
		textW := columnTextWidth(width)
		msg := "(no issues recorded)"
		if overlay.Filter != "" || overlay.ScopedTask != "" {
			msg = "(no issues match this filter)"
		}
		body := frameColumn([]string{styles.CardMeta.Render(msg)}, textW, columnBodyHeight(bodyHeight), true)
		return lipgloss.JoinVertical(lipgloss.Left, header, body)
	}

	grouped := model.groupedIssueRows()

	if width < compactBoardBreakpoint {
		idx := issueColumnIndex(overlay.FocusedStatus)
		if idx < 0 {
			idx = 0
		}
		column := issueColumnLabels[idx]
		body := renderIssueColumn(
			column.Label, column.Status, grouped[column.Status], width, bodyHeight,
			columnCursor{Card: overlay.Cursor, Holds: true, Focused: true},
		)
		return lipgloss.JoinVertical(lipgloss.Left, header, body)
	}

	colWidth := issueColumnWidth(width)
	rendered := make([]string, 0, len(issueColumnLabels))
	for _, column := range issueColumnLabels {
		rendered = append(rendered, renderIssueColumn(
			column.Label, column.Status, grouped[column.Status], colWidth, bodyHeight,
			columnCursor{
				Card:    overlay.Cursor,
				Holds:   overlay.FocusedStatus == column.Status,
				Focused: true,
			},
		))
	}
	return lipgloss.JoinVertical(lipgloss.Left, header, lipgloss.JoinHorizontal(lipgloss.Top, rendered...))
}

// issuesHeaderLine names the overlay's own state — Task scope and type
// filter — since that context now sits above the three status columns
// instead of repeating inside each of them.
func issuesHeaderLine(overlay IssueOverlay) string {
	title := "ISSUES"
	if overlay.ScopedTask != "" {
		title += " · " + overlay.ScopedTask
	}
	return styles.ColumnTitleFocused.Render(title) + "   " + styles.CardMeta.Render("Filter: "+issueFilterLabel(overlay.Filter))
}

// issueColumnWidth splits the terminal three ways, the same rule
// columnWidth applies to the Task board's own three columns, without
// reserving space for the Objective sidebar: the Issues overlay always
// replaces the sidebar's region along with the columns'.
func issueColumnWidth(termW int) int {
	termW = terminalWidthOrOne(termW)
	if termW < compactBoardBreakpoint {
		return termW
	}
	width := termW / len(issueColumnLabels)
	if width < minColumnContent+columnChrome {
		return minColumnContent + columnChrome
	}
	return width
}

// renderIssueColumn draws one status column: its label and row count, a
// rule, and as many rows as the height budget fits — the same shape
// renderColumn draws for a Task column, over IssueRow instead of TaskCard.
func renderIssueColumn(label string, status data.IssueStatus, rows []IssueRow, width, height int, cursor columnCursor) string {
	textW := columnTextWidth(width)
	bodyH := columnBodyHeight(height)

	header := fmt.Sprintf("%s (%d)", label, len(rows))
	lines := []string{
		issueColumnTitleStyle(status, cursor.accented()).Render(header),
		styles.Divider.Render(strings.Repeat("─", textW)),
	}

	if len(rows) == 0 {
		lines = append(lines, styles.CardMeta.Render("(empty)"))
		return frameIssueColumn(lines, textW, bodyH, status, cursor.accented())
	}

	// Each item carries its own trailing blank line, so a tile sits apart
	// from its neighbor by construction; visibleWindow's budget math (which
	// works entirely off heights) stays correct because the blank line is
	// counted in the same item's height rather than inserted separately.
	items := make([]string, len(rows))
	heights := make([]int, len(rows))
	for i, row := range rows {
		items[i] = renderIssueRow(row, status, textW, cursor.highlights(i)) + "\n"
		heights[i] = strings.Count(items[i], "\n") + 1
	}

	budget := bodyH - columnHeaderLines
	if budget < 1 {
		budget = 1
	}
	start, end := visibleWindow(heights, budget, focusedIndex(cursor.Holds, cursor.Card, len(rows)))
	if start > 0 {
		lines = append(lines, scrollIndicator("↑", start, "above"))
	}
	lines = append(lines, items[start:end]...)
	if end < len(rows) {
		lines = append(lines, scrollIndicator("↓", len(rows)-end, "more"))
	}
	return frameIssueColumn(lines, textW, bodyH, status, cursor.accented())
}

// issueRowIndent is how far a row's title and badge lines sit from the
// column edge — lined up under the ID text, which follows a two-cell marker
// exactly as wide.
const issueRowIndent = "  "

// renderIssueRow spends several lines per Issue rather than truncating
// everything onto one: a column at a third of the board's width has too
// little of it for ID, title, type, and severity to survive on a single row
// the way the unsplit overlay could afford. It mirrors renderCard's own
// shape — an ID line, a wrapped title, then a badge line — over an IssueRow
// instead of a TaskCard. Status stays off it entirely: the column it sits in
// already says that.
func renderIssueRow(row IssueRow, status data.IssueStatus, width int, selected bool) string {
	issue := row.Issue
	marker := "  "
	idStyle := styles.CardMeta
	titleStyle := issueRowTitleStyle(status, selected)
	if selected {
		marker = "▸ "
	}

	lines := []string{idStyle.Render(xansi.Truncate(marker+issue.ID, width, "…"))}

	innerWidth := width - len(issueRowIndent)
	if innerWidth < 1 {
		innerWidth = 1
	}
	for _, titleLine := range wrapTitleLines(issue.Title, innerWidth, 2) {
		lines = append(lines, titleStyle.Render(issueRowIndent+titleLine))
	}

	badges := []Badge{issueTypeBadge(issue.Type)}
	if severity, ok := issueSeverityBadge(issue.Severity); ok {
		badges = append(badges, severity)
	}
	for _, badgeLine := range renderBadgeLines(badges, innerWidth) {
		lines = append(lines, issueRowIndent+badgeLine)
	}

	return strings.Join(lines, "\n")
}

// frameIssueColumn draws the column's frame around its lines at exactly the
// height it was budgeted, mirroring frameColumn's Height/MaxHeight split
// between filling a short column out and clipping a long one.
func frameIssueColumn(lines []string, textW, bodyH int, status data.IssueStatus, focused bool) string {
	return issueColumnStyle(status, focused).
		Width(textW + paddingCells).
		Height(bodyH).
		MaxHeight(bodyH + borderCells).
		Render(strings.Join(lines, "\n"))
}

// issueRowTitleStyle takes the column's own accent when selected — orange
// for In Progress, green for Resolved, matching what a selected Task card's
// title does. Open cannot reuse that pattern outright: a Task card's
// selected Planned title stays plain white only because the card still gets
// its own bordered box, which is what actually carries the selection there.
// An Issue row has no per-row box, so Open's selected title instead takes
// the same purple the Objective sidebar's cursor already uses for exactly
// this situation — a plain, unboxed list row that needs its own selection
// accent distinct from every status color. The two surfaces are never shown
// at once, so reusing it here does not collide with its sidebar meaning.
func issueRowTitleStyle(status data.IssueStatus, selected bool) lipgloss.Style {
	if !selected {
		return styles.TaskItem
	}
	switch status {
	case data.IssueStatusOpen:
		return styles.ObjectiveItemFocused
	case data.IssueStatusResolved:
		return styles.TaskItemFocusedDone
	default:
		return styles.TaskItemFocused
	}
}

// issueColumnStyle and issueColumnTitleStyle give each status column one
// accent that its border, its heading, and its own selected row's title
// (issueRowTitleStyle) all agree on: Resolved wears Done's green, In
// Progress wears the plain focused orange every other column wears, and
// Open wears the same purple the Objective sidebar's panel and cursor wear
// — not the Planned column's grey a first pass gave it, which left the
// heading and border out of step with Open's own purple-selected row.
func issueColumnStyle(status data.IssueStatus, focused bool) lipgloss.Style {
	if !focused {
		return styles.ColumnUnfocused
	}
	switch status {
	case data.IssueStatusOpen:
		return styles.SidebarPanelFocused
	case data.IssueStatusResolved:
		return styles.ColumnFocusedDone
	default:
		return styles.ColumnFocused
	}
}

func issueColumnTitleStyle(status data.IssueStatus, focused bool) lipgloss.Style {
	if !focused {
		return styles.ColumnTitle
	}
	switch status {
	case data.IssueStatusOpen:
		return styles.SidebarTitleFocused
	case data.IssueStatusResolved:
		return styles.ColumnTitleFocusedDone
	default:
		return styles.ColumnTitleFocused
	}
}

func renderIssueDetail(detail IssueDetail, width, height, offset int) string {
	textW := columnTextWidth(width)
	bodyH := columnBodyHeight(height)
	lines := []string{
		styles.ColumnTitleFocused.Render("ISSUE DETAIL"),
		styles.Divider.Render(strings.Repeat("─", textW)),
	}
	content := issueDetailLines(detail, textW)
	lines = append(lines, detailWindow(content, detailBudget(bodyH), offset)...)
	return frameColumn(lines, textW, bodyH, true)
}

func issueDetailLines(detail IssueDetail, width int) []string {
	issue := detail.Issue
	lines := []string{
		issueField(width, "ID", issue.ID),
		issueField(width, "Title", issue.Title),
		issueField(width, "Type", string(issue.Type)),
		issueField(width, "Status", string(issue.Status)),
	}
	if issue.Severity != "" {
		lines = append(lines, issueField(width, "Severity", issue.Severity))
	}

	lines = append(lines, "", styles.ColumnTitle.Render("SUMMARY"))
	lines = append(lines, issueBodyLines(issue.Source.Body, width)...)

	lines = append(lines, "", styles.ColumnTitle.Render("ORIGIN"))
	lines = append(lines, issueField(width, "Kind", string(issue.Origin.Kind)))
	if issue.Origin.Check != "" {
		lines = append(lines, issueField(width, "Check", issue.Origin.Check))
	}
	lines = append(lines, issueField(width, "Actor", resume.ActorLabel(issue.Origin.Actor)))
	lines = append(lines, issueField(width, "Time", issue.Origin.At.Format(detailTimeFormat)))

	lines = append(lines, "", styles.ColumnTitle.Render("LINKED TASKS"))
	lines = append(lines, issueLinkLines(detail.Tasks, "(none recorded)", width)...)
	lines = append(lines, "", styles.ColumnTitle.Render("LINKED CHECKS"))
	lines = append(lines, issueLinkLines(detail.Checks, "(none recorded)", width)...)

	if len(detail.GuardrailIDs) > 0 {
		lines = append(lines, "", styles.ColumnTitle.Render("GUARDRAILS"))
		for _, id := range detail.GuardrailIDs {
			lines = append(lines, issueLine(width, id))
		}
	}

	if issue.Resolution != nil {
		lines = append(lines, "", styles.ColumnTitle.Render("RESOLUTION"))
		resolution := issue.Resolution
		lines = append(lines, issueField(width, "Disposition", string(resolution.Disposition)))
		switch resolution.Disposition {
		case data.IssueDispositionVerified:
			if resolution.Check != "" {
				lines = append(lines, issueField(width, "Proof", "Check "+resolution.Check))
			}
		case data.IssueDispositionAccepted:
			lines = append(lines, issueLine(width, "Not proof of repair; this is an owner decision."))
		case data.IssueDispositionDuplicate:
			lines = append(lines, issueLine(width, "Not proof of repair; this points to a canonical Issue."))
			if detail.DuplicateTarget != nil {
				lines = append(lines, issueField(width, "Canonical", detail.DuplicateTarget.ID+" — "+detail.DuplicateTarget.Label))
			}
		case data.IssueDispositionEscalated:
			lines = append(lines, issueLine(width, "Not proof of repair; this points to the Objective the repair was promoted into."))
			if detail.EscalationTarget != nil {
				lines = append(lines, issueField(width, "Objective", detail.EscalationTarget.ID+" — "+detail.EscalationTarget.Label))
			}
		}
		lines = append(lines, issueField(width, "Actor", resume.ActorLabel(resolution.Actor)))
		lines = append(lines, issueField(width, "Time", resolution.At.Format(detailTimeFormat)))
		if strings.TrimSpace(resolution.Reason) != "" {
			lines = append(lines, issueField(width, "Reason", resolution.Reason))
		}
	}

	if len(issue.History) > 0 {
		lines = append(lines, "", styles.ColumnTitle.Render("HISTORY"))
		for _, entry := range issue.History {
			lines = append(lines, issueLine(width, issueHistoryLine(entry)))
		}
	}
	return lines
}

func issueField(width int, label, value string) string {
	return issueLine(width, label+": "+value)
}

func issueLine(width int, text string) string {
	return xansi.Truncate("  "+text, width, "…")
}

func issueBodyLines(body string, width int) []string {
	if strings.TrimSpace(body) == "" {
		return []string{issueLine(width, "(empty)")}
	}
	lines := strings.Split(body, "\n")
	for i, line := range lines {
		lines[i] = issueLine(width, line)
	}
	return lines
}

func issueLinkLines(links []IssueLink, empty string, width int) []string {
	if len(links) == 0 {
		return []string{issueLine(width, empty)}
	}
	lines := make([]string, 0, len(links))
	for _, link := range links {
		text := link.ID
		if link.Label != "" {
			text += " — " + link.Label
		}
		if link.Detail != "" {
			text += " (" + link.Detail + ")"
		}
		lines = append(lines, issueLine(width, text))
	}
	return lines
}

func issueDetailScrollLimit(detail IssueDetail, width, height int) int {
	lines := issueDetailLines(detail, columnTextWidth(width))
	return clampDetailOffset(len(lines), len(lines), detailBudget(columnBodyHeight(height)))
}
