package board

import (
	"fmt"

	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/styles"
)

// linkedFindingsSection renders the read-only "Linked Findings" block shared by
// the task and epic detail overlays: a header, then one row per finding whose
// audit links point back at the work item, or a clear empty state when none do
// (AC3). cursor highlights the selected row. The section is a pure projection of
// the reverse-lookup data and never mutates a finding.
func linkedFindingsSection(findings []data.AuditFinding, cursor, width int) []string {
	section := []string{"", styles.ColumnTitle.Render("Linked Findings:"), ""}
	if len(findings) == 0 {
		return append(section, styles.CardMeta.Render("(no linked findings)"))
	}
	for i, f := range findings {
		section = append(section, renderLinkedFindingRow(f, i == cursor, width))
	}
	return section
}

// renderLinkedFindingRow renders one linked finding as severity tag, ID, title,
// and status, truncated to width. The selected row carries the cursor marker and
// focused styling; others render as dimmed meta rows, matching the audit overlay.
func renderLinkedFindingRow(f data.AuditFinding, selected bool, width int) string {
	row := fmt.Sprintf("%s %s  %s  [%s]", defectSeverityTag(f.Severity), f.ID, f.Title, f.Status)
	if selected {
		return styles.TaskItemFocused.Render(truncate(releaseActiveMarker+" "+row, width))
	}
	return styles.CardMeta.Render(truncate("  "+row, width))
}

// detailFooterHint returns the detail-overlay footer hint, advertising the
// linked-findings navigation only when findings are present so an unlinked work
// item keeps its plain close hint. base is the overlay-specific tail (e.g.
// "esc:close" or "1:Detail 2:Audit  esc:close").
func detailFooterHint(base string, hasFindings bool) string {
	if hasFindings {
		return "↑↓:findings  enter:open  " + base
	}
	return base
}

// focusedTaskFindings returns the findings linked to the board's focused task, or
// nil when no task is focused or no audit register is loaded. It tolerates an
// absent register so detail rendering never blocks (AC7).
func (m Model) focusedTaskFindings() []data.AuditFinding {
	task, ok := m.focusedTask()
	if !ok {
		return nil
	}
	return data.FindingsForTask(m.Audit.Findings, task.ID)
}

// openEpicFindings returns the findings linked to the epic whose detail overlay
// is open, or nil when none resolve or no register is loaded.
func (m Model) openEpicFindings() []data.AuditFinding {
	return data.FindingsForEpic(m.Audit.Findings, m.epicDetailEpic())
}

// moveLinkedFindingCursor moves the linked-findings selection by delta, clamped
// to [0, count-1]. An empty list pins the cursor at 0.
func (m *Model) moveLinkedFindingCursor(delta, count int) {
	if count == 0 {
		m.LinkedFindingCursor = 0
		return
	}
	m.LinkedFindingCursor += delta
	if m.LinkedFindingCursor < 0 {
		m.LinkedFindingCursor = 0
	}
	if m.LinkedFindingCursor >= count {
		m.LinkedFindingCursor = count - 1
	}
}

// activeFinding returns the finding the finding-detail overlay should render,
// resolved from the overlay it was opened from: a work-item detail's selected
// linked finding, or the audit register's cursor finding. It bounds-checks so a
// stale cursor after a reload renders nothing rather than panicking.
func (m Model) activeFinding() (data.AuditFinding, bool) {
	switch m.FindingDetailOrigin {
	case OverlayDetail:
		return findingAt(m.focusedTaskFindings(), m.LinkedFindingCursor)
	case OverlayEpicDetail:
		return findingAt(m.openEpicFindings(), m.LinkedFindingCursor)
	default:
		return m.selectedFinding()
	}
}

// findingDetailReturnOverlay resolves the overlay to return to when the finding
// detail closes. An unset origin falls back to the audit register, the finding's
// canonical home, so a finding detail opened without a recorded origin never
// closes to a blank board.
func (m Model) findingDetailReturnOverlay() OverlayType {
	if m.FindingDetailOrigin == OverlayNone {
		return OverlayAudit
	}
	return m.FindingDetailOrigin
}

// findingAt returns findings[i] when i is in range.
func findingAt(findings []data.AuditFinding, i int) (data.AuditFinding, bool) {
	if i < 0 || i >= len(findings) {
		return data.AuditFinding{}, false
	}
	return findings[i], true
}
