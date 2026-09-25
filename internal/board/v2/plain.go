package v2

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

// plainNonTTYNotice tells a reader piping the board why they are looking at
// text instead of the board.
const plainNonTTYNotice = "[non-interactive mode — run in a TTY to launch the board UI]"

// renderPlain renders one loaded ProjectState as plain text for a writer that
// is not a terminal. It leads with the Next area — the same lines the drawn
// board shows, so piping the board and looking at it answer alike — because a
// reader who pipes the board is asking what is going on. It is built with fmt
// alone, so the output carries no escape sequence, and it reads only state, so
// two runs over the same project produce identical bytes.
//
// selected is the Objective in view, resolved exactly as the terminal board
// resolves it, so the same invocation over the same project counts the same
// Tasks on both surfaces.
func renderPlain(state ProjectState, selected string) string {
	var b strings.Builder

	fmt.Fprintln(&b, plainNonTTYNotice)
	fmt.Fprintln(&b)
	for _, line := range boardNextLines(state) {
		// The interactive Next panel keeps its selection-scoped Issue summary
		// beside the answer. Plain output places the complete project summary
		// after the columns, so a piped board has one deterministic Issues line
		// in the order promised to readers.
		if strings.HasPrefix(line, "Issues: ") {
			continue
		}
		fmt.Fprintln(&b, line)
	}
	fmt.Fprintln(&b)
	fmt.Fprintf(&b, "Objectives: %d  Tasks: %d\n", state.objectiveCount(), state.taskCount())
	release := selectedRelease(state)
	if release != "" {
		fmt.Fprintf(&b, "Selected %s: %s\n", goalLabel, release)
	}
	// "Selected", not "Objective": the Next area's own Objective line names
	// what the projection chose, which the sidebar's filter never moves.
	selection := "no Goal selected"
	if release != "" {
		if selected != "" {
			selection = selected
		} else {
			selection = "all Objectives in Goal " + release
		}
	}
	fmt.Fprintf(&b, "Selected: %s\n", selection)
	fmt.Fprintln(&b, "Objectives:")
	lastPriority := ""
	for _, row := range objectiveRowsForRelease(state.Index, release) {
		priority := string(row.Objective.Priority)
		if priority != lastPriority {
			fmt.Fprintln(&b, objectivePriorityHeading(row.Objective.Priority))
			lastPriority = priority
		}
		fmt.Fprintf(&b, "  %s — %s", row.ID(), row.Objective.Title)
		if badges := row.badges(); len(badges) > 0 {
			labels := make([]string, 0, len(badges))
			for _, badge := range badges {
				labels = append(labels, badge.Text())
			}
			fmt.Fprintf(&b, "  %s", strings.Join(labels, "  "))
		}
		b.WriteByte('\n')
	}
	if notice := unassignedGoalNotice(state.Index); notice != "" {
		fmt.Fprintln(&b, notice)
	}
	fmt.Fprintln(&b)

	cards := groupTaskCardsForRelease(state.Index, release, selected)
	for _, column := range columnLabels {
		fmt.Fprintf(&b, "%-12s %d\n", column.Label, len(cards[column.Status]))
		for _, card := range cards[column.Status] {
			fmt.Fprintf(&b, "  %s — %s", card.Task.ID, card.Task.Title)
			if badges := card.badges(); len(badges) > 0 {
				labels := make([]string, 0, len(badges))
				for _, badge := range badges {
					labels = append(labels, badge.Text())
				}
				fmt.Fprintf(&b, "  [%s]", strings.Join(labels, "  "))
			}
			b.WriteByte('\n')
		}
	}
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, plainIssuesSummary(state.Issues.Rows))

	return stripTerminalControls(b.String())
}

func plainIssuesSummary(rows []IssueRow) string {
	if len(rows) == 0 {
		return "Issues: 0"
	}
	counts := make(map[string]int, len(rows))
	for _, row := range rows {
		counts[string(row.Issue.Type)]++
	}
	parts := make([]string, 0, len(counts))
	for _, issueType := range slices.Sorted(maps.Keys(counts)) {
		parts = append(parts, fmt.Sprintf("%d %s", counts[issueType], issueType))
	}
	return fmt.Sprintf("Issues: %d — %s", len(rows), strings.Join(parts, ", "))
}
