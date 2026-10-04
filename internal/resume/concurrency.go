package resume

import (
	"fmt"
	"slices"
	"strings"

	xansi "github.com/charmbracelet/x/ansi"
	"github.com/opencode/savepoint/internal/data"
)

// This file words the optional parallel-planning advice that data.ResolveConcurrencyV2
// projected: which Tasks may start together, why others are not suggested, and
// a copyable instruction for a fresh session per eligible Task. It is the one
// formatter resume, the Objective and Task details and the plain board share,
// so the same projection reads the same everywhere (STYLE-07, STYLE-09). It is
// pure: it decides no safety (the projection did), probes no worktree, creates
// no branch, and touches no clipboard. Every line is advice the owner may ignore.

const parallelHeading = "Parallel planning is optional: these are suggestions you may ignore. Any Task can run on main or in another worktree."

// ParallelLines renders the advice in c for one surface. focus names a Task
// when the surface is about that Task, and is empty for a whole Objective. It
// returns nil when the preference is off or there is nothing to say, so an
// off project loses the section without losing any record.
func ParallelLines(index *data.V2Index, c *data.ConcurrencyV2, focus string) []string {
	if index == nil || c == nil || !c.Enabled {
		return nil
	}
	body := diagnosticLines(c, focus)
	switch {
	case c.Withheld != nil:
		body = append(body, "Withheld: "+c.Withheld.Detail)
	case focus == "":
		body = append(body, objectiveAdvice(index, c)...)
	default:
		body = append(body, taskAdvice(index, c, focus)...)
	}
	if len(body) == 0 {
		return nil
	}
	lines := append([]string{parallelHeading}, body...)
	for i, line := range lines {
		lines[i] = cleanText(line)
	}
	return lines
}

// diagnosticLines names the planning metadata the projection could not use,
// so the owner can see which record, file and field to look at. They are
// advisory: nothing is blocked and nothing asks the owner to repair them. A
// Task view lists only that Task's own diagnostics.
func diagnosticLines(c *data.ConcurrencyV2, focus string) []string {
	var lines []string
	for _, d := range c.Diagnostics {
		if focus != "" && d.Record != focus {
			continue
		}
		lines = append(lines, fmt.Sprintf("Planning metadata not used — %s %s, field %s: %s", d.Record, d.Path, d.Field, d.Message))
	}
	return lines
}

func objectiveAdvice(index *data.V2Index, c *data.ConcurrencyV2) []string {
	var lines []string
	if len(c.Active) > 0 {
		lines = append(lines, "In progress: "+strings.Join(c.Active, ", "))
	}
	grouped := map[string]bool{}
	for _, group := range c.Groups {
		lines = append(lines, "May start together: "+strings.Join(group, ", "))
		for _, id := range group {
			grouped[id] = true
		}
	}
	for _, candidate := range c.Candidates {
		switch {
		case grouped[candidate.Task]:
		case len(c.Active) > 0:
			lines = append(lines, "May start alongside the work in progress: "+candidate.Task)
		default:
			lines = append(lines, "Ready on its own, with no other lane suggested beside it: "+candidate.Task)
		}
	}
	for _, pair := range c.Pairs {
		if !pair.Together && pair.Reason != nil {
			lines = append(lines, fmt.Sprintf("Not together: %s and %s — %s", pair.Tasks[0], pair.Tasks[1], pair.Reason.Detail))
		}
	}
	for _, note := range c.Notes {
		lines = append(lines, "Not suggested now: "+note.Detail)
	}
	for _, candidate := range c.Candidates {
		if eligible(c, candidate.Task) {
			lines = append(lines, instructionLines(index, c.Objective, candidate)...)
		}
	}
	return lines
}

func taskAdvice(index *data.V2Index, c *data.ConcurrencyV2, id string) []string {
	var lines []string
	if slices.Contains(c.Active, id) {
		lines = append(lines, id+" is in progress; the other lanes are compared against it.")
	}
	for _, group := range c.Groups {
		if slices.Contains(group, id) {
			lines = append(lines, "May start together: "+strings.Join(group, ", "))
		}
	}
	for _, candidate := range c.Candidates {
		if candidate.Task == id && eligible(c, id) {
			return append(lines, instructionLines(index, c.Objective, candidate)...)
		}
	}
	for _, pair := range c.Pairs {
		if !pair.Together && pair.Reason != nil && slices.Contains(pair.Tasks[:], id) {
			lines = append(lines, fmt.Sprintf("Not together: %s and %s — %s", pair.Tasks[0], pair.Tasks[1], pair.Reason.Detail))
		}
	}
	for _, note := range c.Notes {
		if slices.Contains(note.Tasks, id) {
			lines = append(lines, "Not suggested now: "+note.Detail)
		}
	}
	return lines
}

// eligible reports whether a candidate is actually suggested in parallel: it
// belongs to a together group, or runs beside work already in progress. A
// lone candidate with nothing beside it is ordinary sequential work.
func eligible(c *data.ConcurrencyV2, id string) bool {
	if len(c.Active) > 0 {
		return true
	}
	return slices.ContainsFunc(c.Groups, func(group []string) bool { return slices.Contains(group, id) })
}

// instructionLines is one self-contained instruction a fresh session can be
// given. It names no path and runs nothing: the owner prepares any worktree.
func instructionLines(index *data.V2Index, objectiveID string, candidate data.ConcurrencyCandidateV2) []string {
	task, ok := index.Tasks[candidate.Task]
	objective := index.Objectives[objectiveID]
	if !ok || objective == nil {
		return nil
	}
	lane := candidate.Lane
	if l, ok := objective.Plan.Lane(candidate.Lane); ok {
		lane = fmt.Sprintf("%s (%s)", l.Key, l.Title)
	}
	prerequisites := "none"
	if len(task.DependsOn) > 0 {
		parts := make([]string, len(task.DependsOn))
		for i, dep := range task.DependsOn {
			parts[i] = fmt.Sprintf("%s (requires %s)", dep.Task, dep.Requires)
		}
		prerequisites = strings.Join(parts, ", ")
	}
	return []string{
		fmt.Sprintf("Instruction for %s — copy the lines between the markers into a fresh agent session:", task.ID),
		fmt.Sprintf("----- begin %s -----", task.ID),
		startLine(task, objective),
		"Use savepoint-task for the Start line above. It selects this Task even if the shared router names another.",
		fmt.Sprintf("Objective: %s; Task: %s; lane: %s.", objective.ID, task.ID, lane),
		"Reads: " + pathList(task.Plan.Reads.Paths()),
		"Writes: " + pathList(task.Plan.Writes.Paths()),
		"Prerequisites: " + prerequisites + ". Their finished work must already be in your checkout.",
		"Setup: you prepare any worktree and branch yourself, starting from a checkout that includes the prerequisite changes; nothing is created for you and no path is assumed.",
		"This lane is only a suggestion. Running the Task on main or in another worktree is fine, and no lane rule is a completion condition.",
		"If you do use a separate worktree: leave router.md and the Goal selection unchanged; do not allocate Task, Check or Issue identities there; record Task evidence and commit locally on the worktree branch; do not push or merge. Merging and Checks happen on main.",
		fmt.Sprintf("----- end %s -----", task.ID),
	}
}

// startLine is the standalone selection a fresh session routes on, worded
// exactly as the Next line for this Task so the workflow reads it the same way.
func startLine(task *data.TaskV2, objective *data.ObjectiveV2) string {
	return NextLine(data.Next{Task: task, Objective: objective})
}

func pathList(paths []string) string {
	if len(paths) == 0 {
		return "none declared"
	}
	return strings.Join(paths, ", ")
}

// cleanText removes terminal control sequences and control characters from a
// line built from project record strings, so a record cannot move the cursor
// or recolor the reader's terminal. Newlines inside one string become spaces:
// each entry stays one line.
func cleanText(text string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case r == '\n' || r == '\r' || r == '\t':
			return ' '
		case r < 0x20 || (r >= 0x7f && r < 0xa0):
			return -1
		}
		return r
	}, xansi.Strip(text))
}

// parallelFocus is the Task the narrative is about, empty when only an
// Objective is selected.
func parallelFocus(next data.Next) string {
	if next.Task == nil {
		return ""
	}
	return next.Task.ID
}
