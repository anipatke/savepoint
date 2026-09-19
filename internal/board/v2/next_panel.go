package v2

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/resume"
	"github.com/opencode/savepoint/internal/styles"
)

// This file is the Next area: the one next action, formatted from the one
// data.Next the load command resolved. It reads that value and nothing else —
// no index, no record, no gate — so the board cannot quietly arrive at a
// different answer than `savepoint resume` does from the same project.
//
// The evidence wording is internal/resume's, called rather than copied. Two
// surfaces phrasing one answer differently is precisely the divergence the
// shared projection exists to prevent, and a second copy of the clearance,
// exception, owner-wait, and dependency vocabulary is the way that starts
// (STYLE-07, STYLE-09). The layout here is the board's own: a compact block,
// not resume's narrative.

// nextKindLabels is the presentation wording for each rung data.ResolveNext
// can land on. It is a mapping over a resolved value, not a second opinion
// about it: the board chooses the words for NextCheckNeeded, it does not
// decide that a Check is needed.
var nextKindLabels = map[data.NextKind]string{
	data.NextPendingMigration:               "Finish the pending migration",
	data.NextReplan:                         "Replan the selected Task",
	data.NextDependency:                     "Waiting on a dependency",
	data.NextExecute:                        "Work the selected Task",
	data.NextCheckNeeded:                    "A Check is needed",
	data.NextOwnerValidationRequired:        "Accept the recorded Check",
	data.NextObjectiveIntegration:           "Check the Objective's integration",
	data.NextReleaseIntegration:             "Integrate the selected Release",
	data.NextReleaseCheckNeeded:             "Check the selected Release",
	data.NextReleaseOwnerValidationRequired: "Accept the Release Check",
	data.NextReleaseReady:                   "Complete the selected Release",
	data.NextReady:                          "Ready to pick up",
	data.NextPlanObjective:                  "Plan an Objective",
}

// nextLines is the Next area's whole content as plain lines, shared by the
// terminal panel and the non-TTY rendering so a piped board and a drawn one
// state the same thing. The first line is the rung; what follows is the
// record it named, the evidence behind it, any unresolved router selection,
// the follow-ups hanging off the selection, and the action itself.
//
// The action is always present, including when a selection diagnostic is:
// a router line that rotted still leaves a project with work to do, so the
// diagnostic accompanies the answer rather than replacing it.
func nextLines(next data.Next) []string {
	lines := []string{"NEXT: " + nextKindLabel(next.Kind)}
	lines = append(lines, nextIdentityLines(next)...)
	lines = append(lines, resume.EvidenceLines(next)...)
	if next.SelectionDiagnostic != nil {
		lines = append(lines, "Selection: "+resume.SelectionPhrase(next.SelectionDiagnostic))
	}
	if summary := issuesSummary(next.Issues); summary != "" {
		lines = append(lines, summary)
	}
	return append(lines, "Action: "+resume.ActionPhrase(next))
}

// nextKindLabel names one rung. An unrecognized kind reports itself rather
// than falling back to a familiar-looking label, so a rung added later is
// visibly unlabelled instead of silently mislabelled.
func nextKindLabel(kind data.NextKind) string {
	if label, ok := nextKindLabels[kind]; ok {
		return label
	}
	return string(kind)
}

// nextIdentityLines names the record the rung selected, by ID and by title.
// A rung that selected neither — pending migration, and planning a first
// Objective — contributes nothing here rather than an empty label.
func nextIdentityLines(next data.Next) []string {
	var lines []string
	if next.Release != nil {
		lines = append(lines, fmt.Sprintf("Release: %s — %s", next.Release.ID, next.Release.Title))
	}
	if next.Objective != nil {
		lines = append(lines, fmt.Sprintf("Objective: %s — %s", next.Objective.ID, next.Objective.Title))
	}
	if next.Task != nil {
		lines = append(lines, fmt.Sprintf("Task: %s — %s", next.Task.ID, next.Task.Title))
	}
	return lines
}

// issuesSummary states the follow-ups the projection already attached to its
// own selection: how many there are, and how they break down by the type each
// Issue records. It counts the types present rather than listing a fixed
// vocabulary, so a type added to the record model needs no change here, and
// it sorts them so two renders of one project agree. The overlay lists the
// Issues themselves; this says they exist without one being open.
func issuesSummary(issues []*data.IssueV2) string {
	if len(issues) == 0 {
		return ""
	}
	counts := make(map[string]int, len(issues))
	for _, issue := range issues {
		counts[string(issue.Type)]++
	}
	parts := make([]string, 0, len(counts))
	for _, kind := range slices.Sorted(maps.Keys(counts)) {
		parts = append(parts, fmt.Sprintf("%d %s", counts[kind], kind))
	}
	return fmt.Sprintf("Issues: %d — %s", len(issues), strings.Join(parts, ", "))
}

// renderNext draws the Next area. The rung leads in the accent the board uses
// for the work in hand; the rest is meta text, so the panel reads as one block
// rather than competing with the columns for attention.
//
// It is given only the resolved projection, so the sidebar's selection cannot
// move it: the projection answers for the whole project, and a user looking at
// the wrong part of it is exactly who needs that answer.
func (m Model) renderNext(w int) string {
	lines := nextLines(m.State.Next)
	rendered := make([]string, 0, len(lines)+1)
	// A blank line above, and everything under the rung indented beneath it, so
	// the answer reads as one block rather than as more header.
	rendered = append(rendered, "", wrapTo(w, styles.FooterPhaseBuild.Render(lines[0])))
	for _, line := range lines[1:] {
		rendered = append(rendered, wrapTo(w, nextIndent+styles.CardMeta.Render(line)))
	}
	return lipgloss.JoinVertical(lipgloss.Left, rendered...)
}

// nextIndent is the gutter the panel's continuation lines sit in.
const nextIndent = "  "

// wrapTo folds one already-styled line into the terminal's width. A sentence
// longer than the terminal wraps rather than truncating: the evidence behind a
// next action is the part a reader cannot reconstruct from a glyph.
func wrapTo(w int, line string) string {
	return lipgloss.NewStyle().Width(w).Render(line)
}
