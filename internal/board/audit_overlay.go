package board

import (
	"fmt"
	"strings"

	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/styles"
)

// auditTab identifies a tab in the Audit Register overlay. The zero value is the
// Prompt tab so a freshly opened overlay lands on the prompt.
type auditTab int

const (
	auditTabPrompt auditTab = iota
	auditTabFindings
	auditTabRuns
)

// auditTabLabels are the tab strip labels in tab order; the slice length is the
// single source of truth for how many tabs exist.
var auditTabLabels = []string{"Prompt", "Findings", "Runs"}

// selectAuditTab moves the audit overlay to tab, ignoring an out-of-range index.
// Each tab keeps its own scroll offset, so switching does not disturb another
// tab's position.
func (m *Model) selectAuditTab(tab auditTab) {
	if tab < 0 || int(tab) >= len(auditTabLabels) {
		return
	}
	m.AuditTab = tab
}

// scrollAuditTab adjusts the selected tab's scroll offset by delta, clamped at
// the top. The upper bound is left to the renderer, matching the Detail and
// Release Docs scroll behavior. Offsets are stored per tab so an unselected tab
// retains its position.
func (m *Model) scrollAuditTab(delta int) {
	if m.AuditOffsets == nil {
		m.AuditOffsets = map[auditTab]int{}
	}
	offset := m.AuditOffsets[m.AuditTab] + delta
	if offset < 0 {
		offset = 0
	}
	m.AuditOffsets[m.AuditTab] = offset
}

// selectedAuditOffset returns the stored scroll offset for the selected tab, or
// 0 when that tab has not been scrolled.
func (m *Model) selectedAuditOffset() int {
	return m.AuditOffsets[m.AuditTab]
}

// RenderAuditOverlay renders the top-level Audit Register overlay: a tab strip
// (Prompt, Findings, Runs) above a scrollable, read-only body for the selected
// tab. Missing or empty audit assets render an inline empty state rather than an
// error. offset scrolls the selected tab's body. The tab strip is one fixed
// header row beyond the title, so its row is reserved from the body viewport
// like the Release Docs selector.
func RenderAuditOverlay(set data.AuditRegisterSet, tab auditTab, findingCursor, overlayW, maxHeight, offset int) string {
	inner := overlayW - detailBorderPad
	if inner < 4 {
		inner = 4
	}

	lines := []string{
		styles.GlyphAudit.Render("AUDIT REGISTER"),
		renderAuditTabStrip(tab),
	}

	body := auditTabBody(set, tab, findingCursor, inner)
	body = append(body, "", styles.CardMeta.Render("[/]:tab  esc:close"))
	lines = append(lines, visibleDetailLines(body, maxHeight-detailVerticalOverhead-1, offset)...)

	return styles.EpicDetailOverlay.Width(overlayW).Render(strings.Join(lines, "\n"))
}

// renderAuditTabStrip renders the Prompt/Findings/Runs tab strip, emphasizing the
// selected tab using the same active/inactive styling as the other overlays.
func renderAuditTabStrip(tab auditTab) string {
	parts := make([]string, len(auditTabLabels))
	for i, label := range auditTabLabels {
		parts[i] = tabLabel(label, auditTab(i) == tab)
	}
	return strings.Join(parts, styles.CardMeta.Render("  │  "))
}

// auditTabBody returns the rendered body lines for the selected tab.
func auditTabBody(set data.AuditRegisterSet, tab auditTab, findingCursor, width int) []string {
	switch tab {
	case auditTabFindings:
		return auditFindingsBody(set.Findings, findingCursor, width)
	case auditTabRuns:
		return auditRunsBody(set.Runs, width)
	default:
		return auditPromptBody(set.Prompt, width)
	}
}

// auditPromptBody renders the audit prompt markdown, substituting a read-only
// empty state when the prompt is absent or blank.
func auditPromptBody(prompt data.AuditPrompt, width int) []string {
	switch {
	case !prompt.Available:
		return []string{styles.CardMeta.Render("(no audit prompt at " + prompt.Path + ")")}
	case strings.TrimSpace(prompt.Body) == "":
		return []string{styles.CardMeta.Render("(audit prompt is empty)")}
	default:
		return renderReleaseDocBody(prompt.Body, width)
	}
}

// auditFindingsBody renders the findings tab: findings grouped under a status
// header (findings arrive pre-sorted by status, then severity, then ID from the
// data loader) with the cursor row highlighted, or a read-only empty state when
// none are recorded. It is a thin styled projection of auditFindingLayout so the
// renderer and cursor-visibility math share one grouping.
func auditFindingsBody(findings []data.AuditFinding, cursor, width int) []string {
	if len(findings) == 0 {
		return []string{styles.CardMeta.Render("(no audit findings recorded)")}
	}
	lines, _ := auditFindingLayout(findings, cursor, width)
	return lines
}

// auditFindingLayout builds the findings-tab body lines and, in parallel, the
// finding index each line displays (-1 for status headers and separators). The
// single pass is the one source of grouping truth: the renderer emits the lines
// and ensureFindingCursorVisible reads findingLine to locate the cursor row, so
// the two never disagree about where a finding sits.
func auditFindingLayout(findings []data.AuditFinding, cursor, width int) (lines []string, findingLine []int) {
	lastStatus := data.FindingStatus("")
	for i, f := range findings {
		if f.Status != lastStatus {
			if lastStatus != "" {
				lines, findingLine = append(lines, ""), append(findingLine, -1)
			}
			lines = append(lines, styles.ColumnTitle.Render(findingStatusLabel(f.Status)))
			findingLine = append(findingLine, -1)
			lastStatus = f.Status
		}
		lines = append(lines, renderAuditFindingRow(f, i == cursor, width))
		findingLine = append(findingLine, i)
	}
	return lines, findingLine
}

// renderAuditFindingRow renders one finding row as severity tag, ID, title,
// confidence, and linked work item, truncated to the body width. The selected
// row carries a cursor marker and focused styling; others are dimmed meta rows.
func renderAuditFindingRow(f data.AuditFinding, selected bool, width int) string {
	row := fmt.Sprintf("%s %s  %s  conf:%s", defectSeverityTag(f.Severity), f.ID, f.Title, f.Confidence)
	if f.WorkItem != "" {
		row += "  (" + f.WorkItem + ")"
	}
	if selected {
		return styles.TaskItemFocused.Render(truncate(releaseActiveMarker+" "+row, width))
	}
	return styles.CardMeta.Render(truncate("  "+row, width))
}

// findingStatusLabel renders a finding status as an uppercase section header,
// deriving the label from the status token so the status constants stay the sole
// source of truth (e.g. "in_progress" → "IN PROGRESS").
func findingStatusLabel(s data.FindingStatus) string {
	return strings.ToUpper(strings.ReplaceAll(string(s), "_", " "))
}

// auditRunsBody renders one read-only row per run, or an empty state when none
// are recorded. Runs arrive newest-first from the data loader.
func auditRunsBody(runs []data.AuditRun, width int) []string {
	if len(runs) == 0 {
		return []string{styles.CardMeta.Render("(no audit runs recorded)")}
	}
	body := make([]string, 0, len(runs))
	for _, r := range runs {
		body = append(body, renderAuditRunRow(r, width))
	}
	return body
}

// renderAuditRunRow renders a single run as date, mode, and label, truncated to
// the body width.
func renderAuditRunRow(r data.AuditRun, width int) string {
	label := fmt.Sprintf("%s  %-11s %s", r.Date, r.Mode, r.Label)
	return styles.CardMeta.Render(truncate(label, width))
}
