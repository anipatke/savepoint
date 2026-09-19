package v2

import (
	"fmt"
	"strings"

	xansi "github.com/charmbracelet/x/ansi"
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/resume"
	"github.com/opencode/savepoint/internal/styles"
)

const issueHeaderLines = 3

func renderIssues(model Model, width, height int) string {
	overlay := *model.Issues
	textW := columnTextWidth(width)
	bodyH := columnBodyHeight(height)
	rows := issueListRows(model)

	title := "ISSUES"
	if overlay.ScopedTask != "" {
		title += " · " + overlay.ScopedTask
	}
	lines := []string{
		styles.ColumnTitleFocused.Render(title),
		styles.CardMeta.Render("Filter: " + issueFilterLabel(overlay.Filter)),
		styles.Divider.Render(strings.Repeat("─", textW)),
	}

	budget := bodyH - issueHeaderLines
	if budget < 1 {
		budget = 1
	}
	if len(rows) == 0 {
		if overlay.Filter == "" && overlay.ScopedTask == "" {
			lines = append(lines, styles.CardMeta.Render("(no issues recorded)"))
		} else {
			lines = append(lines, styles.CardMeta.Render("(no issues match this filter)"))
		}
		return frameColumn(lines, textW, bodyH, true)
	}

	items := make([]string, len(rows))
	heights := make([]int, len(rows))
	for i, row := range rows {
		items[i] = renderIssueRow(row, textW, i == overlay.Cursor)
		heights[i] = 1
	}
	start, end := visibleWindow(heights, budget, focusedIndex(true, overlay.Cursor, len(rows)))
	if start > 0 {
		lines = append(lines, scrollIndicator("↑", start, "above"))
	}
	lines = append(lines, items[start:end]...)
	if end < len(items) {
		lines = append(lines, scrollIndicator("↓", len(items)-end, "more"))
	}
	return frameColumn(lines, textW, bodyH, true)
}

func renderIssueRow(row IssueRow, width int, selected bool) string {
	issue := row.Issue
	severity := ""
	if issue.Severity != "" {
		severity = "  severity:" + issue.Severity
	}
	text := fmt.Sprintf("%s  %s  %s  %s%s", issue.ID, issue.Title, issue.Type, issue.Status, severity)
	text = xansi.Truncate(text, width-2, "…")
	marker := "  "
	style := styles.CardMeta
	if selected {
		marker = "▸ "
		style = styles.TaskItemFocused
	}
	return style.Render(xansi.Truncate(marker+text, width, "…"))
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

// issueListRows is the renderer's model-specific adapter. It keeps all data
// selection outside the formatting function while allowing tests to render a
// standalone overlay value.
func issueListRows(model Model) []IssueRow {
	return model.filteredIssueRows()
}
