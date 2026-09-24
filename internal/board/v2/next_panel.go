package v2

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/resume"
	"github.com/opencode/savepoint/internal/styles"
)

// This file is the Next area: a one-line glance at the selected Objective,
// Task, or Issue named by the project's resolved projection (data.ResolveNext). The line's
// plain-text wording is shared with `savepoint resume`; the board adds its
// NEXT: label and the verb's accent. The rung evidence, action sentence,
// and Issues summary remain in resume and the record's detail overlay.
//
// It reads only the resolved projection and nothing else — no index, no
// record, no gate — so a reader looking at the wrong part of the board still
// gets the project's real next answer rather than whatever the cursor
// happens to be on.

// nextLines is the Next area's whole content as plain lines, shared by the
// terminal panel and the non-TTY rendering so a piped board and a drawn one
// state the same thing. A selection diagnostic follows the Next line; the
// diagnostic never replaces the action the selected records still support.
func nextLines(next data.Next) []string {
	lines := []string{resume.NextLine(next)}
	if next.Issue != nil && (next.Task != nil || next.Objective != nil) {
		lines = append(lines, resume.IssueContextLine(next.Issue))
	}
	if next.SelectionDiagnostic != nil {
		lines = append(lines, resume.SelectionPhrase(next.SelectionDiagnostic))
	}
	return lines
}

// renderNext draws the Next area: a bold orange "NEXT:" lead-in, then the
// shared line. Only the verb is accented (verbStyle), in the router phase
// row's colours, so the glance and the phase row agree on color.
//
// It is given only the resolved projection, so the sidebar's selection cannot
// move it: the projection answers for the whole project, and a user looking at
// the wrong part of it is exactly who needs that answer.
func (m Model) renderNext(w int) string {
	next := m.State.Next
	lines := nextLines(next)
	rendered := make([]string, 0, len(lines)+1)
	rendered = append(rendered, "")
	for i, line := range lines {
		styledLine := styles.HeaderWhiteBold.Render(line)
		prefix := "      "
		if i == 0 {
			styledLine = styleNextLine(next, line)
			prefix = styles.NextLabel.Render("NEXT:") + " "
		}
		rendered = append(rendered, wrapTo(w, prefix+styledLine))
	}
	return lipgloss.JoinVertical(lipgloss.Left, rendered...)
}

func styleNextLine(next data.Next, line string) string {
	verb := resume.NextVerb(next)
	if verb == "" || !strings.HasPrefix(line, verb) {
		return styles.HeaderWhiteBold.Render(line)
	}
	return verbStyle(verb).Render(verb) + styles.HeaderWhiteBold.Render(line[len(verb):])
}

// verbStyle is the accent the Next line's verb borrows from the router phase
// row: work in the Task phase reads orange, a Check or the owner's closing
// decision reads green, and every other verb stays the panel's bold white.
func verbStyle(verb string) lipgloss.Style {
	switch verb {
	case "Build", "Test", "Fix":
		return styles.FooterPhaseTask
	case "Check", "Accept", "Close":
		return styles.FooterPhaseCheck
	default:
		return styles.HeaderWhiteBold
	}
}

// wrapTo folds one already-styled line into the terminal's width. A sentence
// longer than the terminal wraps rather than truncating.
func wrapTo(w int, line string) string {
	return lipgloss.NewStyle().Width(w).Render(line)
}
