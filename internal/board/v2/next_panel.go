package v2

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/styles"
)

// This file is the Next area: a one-line glance at the Task the project's
// resolved projection (data.ResolveNext) is pointing at right now — its
// lifecycle word, its T-### identity, and its title. Nothing more: the rung
// label, the clearance/owner-wait/dependency evidence, the action sentence,
// and the Issues summary that used to live here are deliberately not
// repeated in this compact panel. That fuller narrative remains
// `savepoint resume`'s job (internal/resume) and the record's own detail
// overlay (detail_view.go) — both unchanged, both still read straight off
// the same projection and resolvers. This panel and those two surfaces are
// allowed to diverge in wording on purpose now: the Next area is a glance,
// not a report.
//
// It reads only the resolved projection and nothing else — no index, no
// record, no gate — so a reader looking at the wrong part of the board still
// gets the project's real next answer rather than whatever the cursor
// happens to be on.

// nextLines is the Next area's whole content as plain lines, shared by the
// terminal panel and the non-TTY rendering so a piped board and a drawn one
// state the same thing.
func nextLines(next data.Next) []string {
	if next.Task != nil {
		return []string{taskStageWord(next.Task) + " " + next.Task.ID + " — " + next.Task.Title}
	}
	if next.Objective != nil {
		return []string{next.Objective.ID + " — " + next.Objective.Title}
	}
	if next.Kind == data.NextPendingMigration {
		return []string{"Migration in progress (" + next.Migration.OperationID + ")"}
	}
	return []string{"Nothing selected yet"}
}

// taskStageWord is the one word this panel leads a Task with: its
// implementation stage while in progress (Build, Test, or Check — never
// "Audit"; see stageLabel), or its own status otherwise.
func taskStageWord(task *data.TaskV2) string {
	if task.Status == data.ColumnInProgress {
		switch task.Stage {
		case data.StageBuild:
			return "Build"
		case data.StageTest:
			return "Test"
		case data.StageAudit:
			return "Check"
		}
	}
	switch task.Status {
	case data.ColumnPlanned:
		return "Planned"
	case data.ColumnDone:
		return "Done"
	default:
		return string(task.Status)
	}
}

// renderNext draws the Next area: a bold orange "NEXT:" lead-in, then the
// word, the identity, and the title read as a single glance. The word itself
// — Build, Test, or Check — carries the same accent the router row uses for
// that phase (Task's orange for Build/Test, Check's green for Check), so the
// glance and the phase row always agree on color.
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
		if i == 0 && next.Task != nil {
			word := taskStageWord(next.Task)
			styledLine = stageWordStyle(next.Task).Render(word) + styles.HeaderWhiteBold.Render(line[len(word):])
		}
		rendered = append(rendered, wrapTo(w, styles.NextLabel.Render("NEXT:")+" "+styledLine))
	}
	return lipgloss.JoinVertical(lipgloss.Left, rendered...)
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
