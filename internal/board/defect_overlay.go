package board

import (
	"strings"

	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/styles"
)

var defectStatusOrder = []data.DefectStatus{
	data.DefectOpen,
	data.DefectInProgress,
	data.DefectResolved,
}

var defectSectionLabel = map[data.DefectStatus]string{
	data.DefectOpen:       "OPEN",
	data.DefectInProgress: "IN PROGRESS",
	data.DefectResolved:   "RESOLVED",
}

// defectsForOverlay returns defects filtered to release and sorted by status
// order (open → in-progress → resolved), preserving discovery order within
// each group.
func defectsForOverlay(all []data.Defect, release string) []data.Defect {
	buckets := map[data.DefectStatus][]data.Defect{}
	for _, d := range all {
		if release != "" && d.Release != release {
			continue
		}
		col := d.Status
		if col == "" {
			col = data.DefectOpen
		}
		buckets[col] = append(buckets[col], d)
	}
	out := make([]data.Defect, 0, len(all))
	for _, col := range defectStatusOrder {
		out = append(out, buckets[col]...)
	}
	return out
}

// RenderDefectsOverlay renders the defect list overlay for a release.
func RenderDefectsOverlay(defects []data.Defect, cursor, width int) string {
	inner := width - epicPanelOverhead
	if inner < 4 {
		inner = 4
	}

	lines := []string{
		styles.ColumnTitleFocused.Render("DEFECTS"),
		strings.Repeat("─", inner),
	}

	if len(defects) == 0 {
		lines = append(lines, styles.TaskItem.Render("(no defects)"))
	} else {
		lastSection := data.DefectStatus("")
		for i, d := range defects {
			col := d.Status
			if col == "" {
				col = data.DefectOpen
			}
			if col != lastSection {
				if lastSection != "" {
					lines = append(lines, "")
				}
				label := defectSectionLabel[col]
				if label == "" {
					label = string(col)
				}
				lines = append(lines, styles.ColumnTitle.Render(label))
				lastSection = col
			}
			row := defectRowText(d, inner-4)
			if i == cursor {
				lines = append(lines, styles.TaskItemFocused.Render(releaseActiveMarker+" "+row))
			} else {
				lines = append(lines, styles.TaskItem.Render("  "+row))
			}
		}
	}

	lines = append(lines, "", styles.CardMeta.Render("↑↓:nav  space:resolve  esc:close"))
	return styles.EpicPanel.Width(width).Render(strings.Join(lines, "\n"))
}

// defectRowText composes the display line for one defect entry.
func defectRowText(d data.Defect, maxWidth int) string {
	sev := defectSeverityTag(d.Severity)
	id := shortID(d.ID)
	title := d.Title
	ref := d.Reference

	row := sev + " " + id
	if title != "" {
		row += "  " + title
	}
	if ref != "" {
		row += "  (" + ref + ")"
	}
	return truncate(row, maxWidth)
}

func defectSeverityTag(s data.DefectSeverity) string {
	switch s {
	case data.SeverityCritical:
		return "[⚠]"
	case data.SeverityHigh:
		return "[h]"
	case data.SeverityMedium:
		return "[m]"
	case data.SeverityLow:
		return "[l]"
	default:
		return "[?]"
	}
}
