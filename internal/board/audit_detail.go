package board

import (
	"strings"

	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/styles"
)

// findingLink pairs a detail label with the linked IDs it renders, so the link
// rows share one definition rather than repeating the "skip when empty" branch.
type findingLink struct {
	label string
	ids   []string
}

// RenderFindingDetail renders the read-only finding detail overlay: the finding's
// fixed metadata, its links and locations, the proof it still needs, its
// first/last-seen dates, and its markdown body. It never mutates the finding.
func RenderFindingDetail(f data.AuditFinding, overlayW, maxHeight, offset int) string {
	inner := overlayW - detailBorderPad
	if inner < 4 {
		inner = 4
	}

	lines := []string{
		styles.ColumnTitleFocused.Render("FINDING DETAIL"),
		strings.Repeat("─", inner),
	}

	body := []string{
		detailRow("ID", f.ID, inner),
		detailRow("Title", f.Title, inner),
		detailRow("Status", string(f.Status), inner),
		detailRow("Severity", string(f.Severity), inner),
		detailRow("Confidence", string(f.Confidence), inner),
	}

	body = append(body, findingLinkRows(f, inner)...)

	if len(f.Locations) > 0 {
		body = append(body, detailRow("Locations", strings.Join(f.Locations, ", "), inner))
	}
	if f.ProofNeeded != "" {
		body = append(body, detailRow("Proof Needed", f.ProofNeeded, inner))
	}
	body = append(body,
		detailRow("First Seen", f.FirstSeen, inner),
		detailRow("Last Seen", f.LastSeen, inner),
	)

	if sections := renderFindingBody(f.Body, inner); len(sections) > 0 {
		body = append(body, "")
		body = append(body, sections...)
	}

	body = append(body, "", styles.CardMeta.Render("esc:close"))
	lines = append(lines, visibleDetailLines(body, maxHeight-detailVerticalOverhead, offset)...)
	return styles.DetailOverlay.Width(overlayW).Render(strings.Join(lines, "\n"))
}

// findingLinkRows renders the finding's work-item link and its related-record
// links (releases, epics, tasks, defects, guardrails), skipping any that are
// absent so an unlinked finding shows no empty rows.
func findingLinkRows(f data.AuditFinding, width int) []string {
	var rows []string
	if f.WorkItem != "" {
		rows = append(rows, detailRow("Work Item", f.WorkItem, width))
	}
	links := []findingLink{
		{"Releases", f.Releases},
		{"Epics", f.Epics},
		{"Tasks", f.Tasks},
		{"Defects", f.Defects},
		{"Guardrails", f.GuardrailIDs},
	}
	for _, link := range links {
		if len(link.ids) > 0 {
			rows = append(rows, detailRow(link.label, strings.Join(link.ids, ", "), width))
		}
	}
	return rows
}

// renderFindingBody renders the finding's markdown body sections as wrapped,
// styled display lines, preserving the body's own heading and list structure.
// A blank body yields no lines so the detail view omits the section entirely.
func renderFindingBody(body string, width int) []string {
	if strings.TrimSpace(body) == "" {
		return nil
	}
	return renderReleaseDocBody(body, width)
}
