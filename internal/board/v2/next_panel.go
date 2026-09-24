package v2

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/resume"
	"github.com/opencode/savepoint/internal/styles"
)

// This file is the Next area: a one-line glance at the Objective and Task
// named by the project's resolved projection (data.ResolveNext). The line's
// plain-text wording is shared with `savepoint resume`; the board adds its
// NEXT: label and existing word accents. The rung evidence, action sentence,
// and Issues summary remain in resume and the record's detail overlay.
//
// It reads only the resolved projection and nothing else — no index, no
// record, no gate — so a reader looking at the wrong part of the board still
// gets the project's real next answer rather than whatever the cursor
// happens to be on.

// nextLines is the Next area's whole content as plain lines, shared by the
// terminal panel and the non-TTY rendering so a piped board and a drawn one
// state the same thing.
func nextLines(next data.Next) []string {
	return []string{resume.NextLine(next)}
}

// renderNext draws the Next area: a bold orange "NEXT:" lead-in, then the
// shared line. Objective and Task words use the existing phase accents so
// the glance and the router phase row agree on color.
//
// It is given only the resolved projection, so the sidebar's selection cannot
// move it: the projection answers for the whole project, and a user looking at
// the wrong part of it is exactly who needs that answer.
func (m Model) renderNext(w int) string {
	next := m.State.Next
	lines := nextLines(next)
	rendered := make([]string, 0, len(lines)+1)
	rendered = append(rendered, "")
	for _, line := range lines {
		styledLine := styleNextLine(next, line)
		rendered = append(rendered, wrapTo(w, styles.NextLabel.Render("NEXT:")+" "+styledLine))
	}
	return lipgloss.JoinVertical(lipgloss.Left, rendered...)
}

func styleNextLine(next data.Next, line string) string {
	if next.Task != nil {
		taskWord := resume.TaskStageWord(next.Task)
		if next.Objective == nil {
			return stageWordStyle(next.Task).Render(taskWord) + styles.HeaderWhiteBold.Render(line[len(taskWord):])
		}

		objectiveWord := resume.ObjectiveWord(next.Objective)
		separator := strings.Index(line, " · ")
		if objectiveWord == "" || taskWord == "" || separator < 0 {
			return styles.HeaderWhiteBold.Render(line)
		}
		taskWordStart := separator + len(" · ")
		taskWordEnd := taskWordStart + len(taskWord)
		return objectiveWordStyle(next.Objective).Render(objectiveWord) +
			styles.HeaderWhiteBold.Render(line[len(objectiveWord):taskWordStart]) +
			stageWordStyle(next.Task).Render(taskWord) +
			styles.HeaderWhiteBold.Render(line[taskWordEnd:])
	}

	if next.Objective != nil {
		objectiveWord := resume.ObjectiveWord(next.Objective)
		if objectiveWord == "" {
			return styles.HeaderWhiteBold.Render(line)
		}
		if next.Kind == data.NextObjectiveIntegration || next.Kind == data.NextObjectiveReady {
			separator := strings.Index(line, " · Check — ")
			if separator >= 0 {
				checkStart := separator + len(" · ")
				checkEnd := checkStart + len("Check")
				return objectiveWordStyle(next.Objective).Render(objectiveWord) +
					styles.HeaderWhiteBold.Render(line[len(objectiveWord):checkStart]) +
					styles.FooterPhaseCheck.Render("Check") +
					styles.HeaderWhiteBold.Render(line[checkEnd:])
			}
		}
		return objectiveWordStyle(next.Objective).Render(objectiveWord) +
			styles.HeaderWhiteBold.Render(line[len(objectiveWord):])
	}

	return styles.HeaderWhiteBold.Render(line)
}

func objectiveWordStyle(objective *data.ObjectiveV2) lipgloss.Style {
	switch objective.Status {
	case data.ColumnInProgress:
		return styles.FooterPhaseTask
	case data.ColumnDone:
		return styles.FooterPhaseCheck
	default:
		return styles.HeaderWhiteBold
	}
}

// stageWordStyle is the accent a Task's stage word borrows from the router
// phase row: Build and Test read as the router's Task phase (orange), Check
// reads as the router's Check phase (green). Planned and Done are not a
// router phase, so they stay the panel's plain bold white.
func stageWordStyle(task *data.TaskV2) lipgloss.Style {
	if task.Status == data.ColumnInProgress {
		switch task.Stage {
		case data.StageBuild, data.StageTest:
			return styles.FooterPhaseTask
		case data.StageAudit:
			return styles.FooterPhaseCheck
		}
	}
	return styles.HeaderWhiteBold
}

// wrapTo folds one already-styled line into the terminal's width. A sentence
// longer than the terminal wraps rather than truncating.
func wrapTo(w int, line string) string {
	return lipgloss.NewStyle().Width(w).Render(line)
}
