package v2

import (
	"maps"
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"
	xansi "github.com/charmbracelet/x/ansi"
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/styles"
)

// The sidebar is where an Objective stops being a folder and becomes a record
// with its own state. V1 navigated a release, then an epic — two selectors over
// directory structure. V2 has one level and it is not a directory: an Objective
// owns a Task because that Task's own objective field says so, which is what
// index.ObjectiveTasks records. Membership is read from that map and from
// nothing else; no path, directory name, or ID prefix is consulted anywhere in
// this file.

const (
	sidebarTitle = "OBJECTIVES"
	// sidebarHeaderLines is the title plus the rule under it, the same shape a
	// column's header has.
	sidebarHeaderLines = 2
	// sidebarEmpty is what an Objective-less project shows. A project with no
	// Objectives yet is a normal state — the wording matches an empty column so
	// the board reads as one surface.
	sidebarEmpty = "(empty)"
)

// glyphSelected marks the Objective the columns are filtered to; glyphCursor
// marks the row the sidebar cursor is on. They are separate cells because a row
// can be both, and because the visual identity requires selection and focus to
// be legible without color.
const (
	glyphSelected = "●"
	glyphCursor   = "▸"
	// rowMarkerCells is the width those two glyphs and their trailing space
	// occupy on every row, marked or not, so a row's text starts in the same
	// column whatever its state.
	rowMarkerCells = 3
)

// ObjectiveRow is one Objective together with the resolved values its sidebar
// row shows. Like TaskCard it exists so rendering is pure: renderSidebar
// formats rows and calls no resolver, so what is on screen is exactly what the
// resolvers returned.
type ObjectiveRow struct {
	Objective *data.ObjectiveV2
	Clearance data.Clearance
	// TasksComplete is true when the Objective owns at least one Task and every
	// one of them is done. It is the first half of the completion question; the
	// Objective's own integration clearance is the other, and an Objective that
	// has one without the other is exactly the state the sidebar exists to make
	// visible.
	TasksComplete bool
	// ByException is true when a recorded owner exception against the
	// Objective's latest Check grants ResolveObjectiveCompletion's Allowed
	// result despite clearance not being current. It is never true merely
	// because clearance is current — that is a passed Check, not an
	// exception.
	ByException bool
	// Waits are the Objective's own declared dependencies that
	// ResolveObjectiveDependency reports unsatisfied, in declared order.
	Waits []data.ObjectiveDependencyBlock
}

// ID is the Objective's global identity, the value selection is keyed by.
func (r ObjectiveRow) ID() string {
	return r.Objective.ID
}

// objectiveRows resolves every Objective in index into a sidebar row, in
// ascending O### order so a project renders the same way twice. Every value on
// a row comes from a resolver or from the record itself.
func objectiveRows(index *data.V2Index) []ObjectiveRow {
	return objectiveRowsForRelease(index, "")
}

// objectiveRowsForRelease resolves the Objectives visible in one Release
// context. When releaseID is empty it keeps the release-free V2 behavior. A
// non-empty context reads the derived reverse link and never guesses from a
// title, path, or identifier prefix.
func objectiveRowsForRelease(index *data.V2Index, releaseID string) []ObjectiveRow {
	if index == nil {
		return nil
	}

	ids := objectiveIDsForRelease(index, releaseID)
	rows := make([]ObjectiveRow, 0, len(ids))
	for _, id := range ids {
		objective := index.Objectives[id]
		rows = append(rows, ObjectiveRow{
			Objective:     objective,
			Clearance:     data.ResolveClearance(index, id),
			TasksComplete: ownedTasksComplete(index, id),
			ByException:   data.ResolveObjectiveCompletion(index, id).AllowedByException,
			Waits:         unsatisfiedObjectiveWaits(index, objective),
		})
	}
	return rows
}

func objectiveIDsForRelease(index *data.V2Index, releaseID string) []string {
	if index == nil {
		return nil
	}
	if releaseID != "" {
		ids := slices.Clone(index.ReleaseObjectives[releaseID])
		slices.Sort(ids)
		return ids
	}
	return slices.Sorted(maps.Keys(index.Objectives))
}

// ownedTasksComplete reports whether every Task the Objective owns is done.
// An Objective owning no Task is not complete: there is no finished work to
// report, and calling it complete would read as an achievement.
func ownedTasksComplete(index *data.V2Index, objectiveID string) bool {
	owned := index.ObjectiveTasks[objectiveID]
	if len(owned) == 0 {
		return false
	}
	for _, taskID := range owned {
		if index.Tasks[taskID].Status != data.ColumnDone {
			return false
		}
	}
	return true
}

// unsatisfiedObjectiveWaits resolves each Objective dependency through
// ResolveObjectiveDependency and keeps the blocks it reports. The sidebar says
// an Objective is waiting because the resolver said so, not because the
// dependency's status field looked unfinished.
func unsatisfiedObjectiveWaits(index *data.V2Index, objective *data.ObjectiveV2) []data.ObjectiveDependencyBlock {
	var waits []data.ObjectiveDependencyBlock
	for _, dependencyID := range objective.DependsOn {
		if decision := data.ResolveObjectiveDependency(index, dependencyID); decision.Block != nil {
			waits = append(waits, *decision.Block)
		}
	}
	return waits
}

// taskIDsInView is the ownership filter, and the only place membership is
// decided: the Tasks index.ObjectiveTasks records for the selected Objective,
// or every Task when nothing is selected. Both orders are ascending by ID, so
// the columns do not reorder when a selection changes.
func taskIDsInView(index *data.V2Index, objectiveID string) []string {
	return taskIDsInReleaseView(index, "", objectiveID)
}

// taskIDsInReleaseView is the only release-aware Task membership filter. A
// selected Objective narrows through ObjectiveTasks; otherwise the selected
// Release expands through ReleaseObjectives and then ObjectiveTasks. Both
// maps are index links built from authored ownership fields.
func taskIDsInReleaseView(index *data.V2Index, releaseID, objectiveID string) []string {
	if index == nil {
		return nil
	}
	if objectiveID == "" {
		if releaseID == "" {
			return slices.Sorted(maps.Keys(index.Tasks))
		}
		var taskIDs []string
		for _, objectiveID := range objectiveIDsForRelease(index, releaseID) {
			taskIDs = append(taskIDs, index.ObjectiveTasks[objectiveID]...)
		}
		slices.Sort(taskIDs)
		return taskIDs
	}
	if releaseID != "" {
		objective, ok := index.Objectives[objectiveID]
		if !ok || string(objective.Release) != releaseID {
			return nil
		}
	}
	return index.ObjectiveTasks[objectiveID]
}

// badges is the row's state line: the Objective's own Check state, and each
// Objective it waits on. Every badge comes from the one mapping in badges.go.
//
// This is deliberately one completion-state badge, not two. An earlier
// version also carried objectiveIntegrationBadge ("INTEGRATED" / "NEEDS
// INTEGRATION"), which — once every owned Task was done — restated the exact
// same clearance.State the Check badge already showed, in different and
// more alarming-sounding words. That duplication was flagged (I012) and
// removed rather than patched: objectiveCheckBadge alone, with byException
// folded in, is now this row's whole answer to "is this Objective ready to
// close." TasksComplete stays on ObjectiveRow because callers besides this
// method may still need "are the owned Tasks done" as its own fact; it is no
// longer read here.
func (r ObjectiveRow) badges() []Badge {
	badges := []Badge{objectiveCheckBadge(r.Clearance.State, r.ByException)}
	for _, wait := range r.Waits {
		badges = append(badges, objectiveWaitBadge(wait))
	}
	return badges
}

// renderSidebar draws the Objective list at the given outer size: the title, a
// rule, and as many rows as the height budget fits, with an indicator for
// whatever is scrolled out of view. It reuses the columns' windowing, so a long
// sidebar scrolls exactly as a long column does.
//
// focused changes the accent and adds the cursor glyph. It changes no width and
// no line count: the frame, the marker cells, and the text width are identical
// in both states.
func renderSidebar(rows []ObjectiveRow, selected string, cursor int, focused bool, width, height int) string {
	textW := columnTextWidth(width)
	bodyH := columnBodyHeight(height)

	title := styles.ColumnTitle.Render(sidebarTitle)
	if focused {
		title = styles.SidebarTitleFocused.Render(sidebarTitle)
	}
	lines := []string{title, styles.Divider.Render(strings.Repeat("─", textW))}

	if len(rows) == 0 {
		lines = append(lines, styles.CardMeta.Render(sidebarEmpty))
		return frameSidebar(lines, textW, bodyH, focused)
	}

	rendered := make([]string, len(rows))
	heights := make([]int, len(rows))
	for i, row := range rows {
		// Each row ends in a blank line so the list reads as blocks rather than
		// one run of text. The spacer belongs to the row, so the window budget
		// counts it like any other line.
		rendered[i] = renderObjectiveRow(row, textW, row.ID() == selected, focused && i == cursor) + "\n"
		heights[i] = strings.Count(rendered[i], "\n") + 1
	}

	budget := bodyH - sidebarHeaderLines
	if budget < 1 {
		budget = 1
	}
	// The sidebar always keeps its own cursor row in view, focused or not, so
	// the selected Objective stays on screen while the keys are in the columns.
	start, end := visibleWindow(heights, budget, focusedIndex(true, cursor, len(rows)))

	if start > 0 {
		lines = append(lines, scrollIndicator("↑", start, "above"))
	}
	lines = append(lines, rendered[start:end]...)
	if end < len(rows) {
		lines = append(lines, scrollIndicator("↓", len(rows)-end, "more"))
	}

	return frameSidebar(lines, textW, bodyH, focused)
}

// frameSidebar draws the sidebar's own frame, mirroring frameColumn's shape
// exactly: the same gray border the Task columns wear when unfocused, and the
// sidebar's own purple accent — distinct from the columns' orange — only once
// the sidebar itself holds focus.
func frameSidebar(lines []string, textW, bodyH int, focused bool) string {
	return sidebarStyle(focused).
		Width(textW + paddingCells).
		Height(bodyH).
		MaxHeight(bodyH + borderCells).
		Render(strings.Join(lines, "\n"))
}

// sidebarStyle is the sidebar frame in two accents, mirroring columnStyle:
// gray unfocused, purple focused. Only the color differs — the border and
// padding are identical in both states.
func sidebarStyle(focused bool) lipgloss.Style {
	if focused {
		return styles.SidebarPanelFocused
	}
	return styles.ColumnUnfocused
}

// renderObjectiveRow draws one row: its markers and O### identity, the human
// title its author wrote, the status its record records, and its badges. The
// title wraps across up to two lines before truncating so the whole title can
// be read, while subsequent lines are indented past the marker column.
func renderObjectiveRow(row ObjectiveRow, width int, selected, cursor bool) string {
	textW := width - rowMarkerCells
	if textW < 4 {
		textW = 4
	}

	style := styles.TaskItem
	switch {
	case cursor:
		style = styles.ObjectiveItemFocused
	case selected:
		style = styles.SidebarSelected
	}

	titleLines := wrapTitleLines(row.ID()+" "+row.Objective.Title, textW, 2)
	var lines []string
	if len(titleLines) > 0 {
		lines = append(lines, style.Render(rowMarkers(selected, cursor)+titleLines[0]))
		for _, line := range titleLines[1:] {
			lines = append(lines, style.Render(indent(line)))
		}
	} else {
		lines = append(lines, style.Render(rowMarkers(selected, cursor)))
	}

	lines = append(lines, styles.CardMeta.Render(indent(xansi.Truncate(objectiveStatusLabel(row.Objective.Status), textW, "…"))))
	for _, line := range renderBadgeLines(row.badges(), textW) {
		lines = append(lines, indent(line))
	}
	return strings.Join(lines, "\n")
}

// objectiveStatusLabel is the one place an Objective's recorded Status
// becomes owner-facing text (I013): the row shows "Planned", "In Progress",
// and "Done" rather than the raw frontmatter values, which stay
// `planned`/`in_progress`/`done` in every record and resolver untouched. A
// status outside those three reports itself rather than a guess.
func objectiveStatusLabel(status data.ColumnType) string {
	switch status {
	case data.ColumnPlanned:
		return "Planned"
	case data.ColumnInProgress:
		return "In Progress"
	case data.ColumnDone:
		return "Done"
	default:
		return string(status)
	}
}

// rowMarkers is the fixed-width cursor and selection column. Both glyphs keep
// their cell whether or not they are drawn, so no row shifts sideways when the
// cursor moves onto it.
func rowMarkers(selected, cursor bool) string {
	marker := func(on bool, glyph string) string {
		if on {
			return glyph
		}
		return " "
	}
	return marker(cursor, glyphCursor) + marker(selected, glyphSelected) + " "
}

// indent aligns a row's continuation lines under its title, past the marker
// cells.
func indent(line string) string {
	return strings.Repeat(" ", rowMarkerCells) + line
}
