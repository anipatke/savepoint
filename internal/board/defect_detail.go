package board

import (
	"fmt"
	"strings"

	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/styles"
)

// defectSections lists the markdown section headings rendered in the defect detail overlay,
// in presentation order. Sections absent from the body are silently skipped.
var defectSections = []string{
	"Symptom",
	"Expected Behavior",
	"Reproduction",
	"Impact",
	"Fix Plan",
	"Acceptance Criteria",
	"Resolution Notes",
}

// RenderDefectDetail renders a defect detail overlay at the given display width.
func RenderDefectDetail(d data.Defect, overlayW, maxHeight, offset int) string {
	inner := overlayW - detailBorderPad
	if inner < 4 {
		inner = 4
	}

	lines := []string{
		styles.ColumnTitleFocused.Render("DEFECT DETAIL"),
		strings.Repeat("─", inner),
	}

	body := []string{
		detailRow("ID", shortDefectID(d.ID), inner),
		detailRow("Title", d.Title, inner),
		detailRow("Status", string(d.Status), inner),
		detailRow("Severity", string(d.Severity), inner),
		detailRow("Release", d.Release, inner),
	}
	if d.Introduced != "" {
		body = append(body, detailRow("Introduced", d.Introduced, inner))
	}
	if d.Reference != "" {
		body = append(body, detailRow("Reference", d.Reference, inner))
	}

	sections := parseDefectSections(d.Body)
	for _, heading := range defectSections {
		content, ok := sections[heading]
		if !ok || strings.TrimSpace(content) == "" {
			continue
		}
		body = append(body, "", styles.ColumnTitle.Render(heading+":"))
		for _, line := range WrapText(strings.TrimSpace(content), inner) {
			body = append(body, styles.CardMeta.Render(line))
		}
	}

	body = append(body, "", styles.CardMeta.Render("esc:close"))
	lines = append(lines, visibleDetailLines(body, maxHeight-detailVerticalOverhead, offset)...)
	return styles.DetailOverlay.Width(overlayW).Render(strings.Join(lines, "\n"))
}

// parseDefectSections parses a markdown body into a map of heading → body text.
// Only `## Level` headings are recognised. Content before the first heading is ignored.
func parseDefectSections(body string) map[string]string {
	sections := map[string]string{}
	var currentHeading string
	var buf strings.Builder

	for _, line := range strings.Split(body, "\n") {
		if heading, ok := extractH2Heading(line); ok {
			if currentHeading != "" {
				sections[currentHeading] = buf.String()
			}
			currentHeading = heading
			buf.Reset()
			continue
		}
		if currentHeading != "" {
			buf.WriteString(line)
			buf.WriteByte('\n')
		}
	}
	if currentHeading != "" {
		sections[currentHeading] = buf.String()
	}
	return sections
}

// extractH2Heading returns the heading text when line is a `## Heading` line.
func extractH2Heading(line string) (string, bool) {
	if !strings.HasPrefix(line, "## ") {
		return "", false
	}
	heading := strings.TrimSpace(line[3:])
	if heading == "" {
		return "", false
	}
	return heading, true
}

// shortDefectID strips any release-prefix path component and slug from a defect ID.
// "v1.1/D003-auth-crash" → "D003"
func shortDefectID(id string) string {
	if idx := strings.LastIndex(id, "/"); idx >= 0 {
		id = id[idx+1:]
	}
	if idx := strings.Index(id, "-"); idx >= 0 {
		id = id[:idx]
	}
	return id
}

// defectMarkerForTask returns a compact count marker (e.g. "⚠ 2/1") summarising
// how many defects reference the given task: open count (open+in_progress) / resolved count.
// Returns "" when no defects match.
// Matching uses the defect Reference field against both the full task ID and the short task ID.
func defectMarkerForTask(taskID string, defects []data.Defect) string {
	short := shortID(taskID)
	var open, resolved int
	for _, d := range defects {
		if d.Reference == "" {
			continue
		}
		ref := d.Reference
		if ref != taskID && ref != short {
			continue
		}
		if d.Status == data.DefectResolved {
			resolved++
		} else {
			open++
		}
	}
	if open+resolved == 0 {
		return ""
	}
	return fmt.Sprintf("%s  %d/%d", glyphWarning, open, resolved)
}
