package v2

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	xansi "github.com/charmbracelet/x/ansi"
	"github.com/muesli/termenv"
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/testutil"
)

// sidebarBoard opens a board wide enough to carry the sidebar beside the three
// columns, and tall enough for every fixture Objective to be on screen at once
// below the Next area, which is the size every sidebar assertion below is made
// at.
func sidebarBoard(t *testing.T, root string) Model {
	t.Helper()
	return openSizedBoard(t, root, 130, 56)
}

// sidebarLines returns the rendered board's lines with their ANSI stripped,
// narrowed to the sidebar panel: its own columns, from its title down. An
// assertion about the sidebar can then not be satisfied by a card, by the
// header, or by the selection line above it.
func sidebarLines(t *testing.T, model Model) []string {
	t.Helper()
	var lines []string
	for _, line := range strings.Split(xansi.Strip(model.View()), "\n") {
		if lipgloss.Width(line) < sidebarWidth {
			continue
		}
		cut := xansi.Cut(line, 0, sidebarWidth)
		if len(lines) == 0 && !strings.Contains(cut, sidebarTitle) {
			continue
		}
		lines = append(lines, cut)
	}
	return lines
}

func sidebarText(t *testing.T, model Model) string {
	t.Helper()
	return strings.Join(sidebarLines(t, model), "\n")
}

// TestSidebarListsEveryObjectiveInOrder covers the list itself: every Objective
// in the index, in stable O-### order, each row carrying its ID, its title, and
// its recorded status.
func TestSidebarListsEveryObjectiveInOrder(t *testing.T) {
	got := sidebarText(t, sidebarBoard(t, writeNavigationProject(t)))

	if !strings.Contains(got, sidebarTitle) {
		t.Errorf("no Objective sidebar rendered:\n%s", got)
	}

	previous := -1
	for _, id := range []string{"O-001", "O-002", "O-003", "O-004", "O-005", "O-006"} {
		at := strings.Index(got, id)
		if at < 0 {
			t.Fatalf("sidebar is missing Objective %s:\n%s", id, got)
		}
		if at < previous {
			t.Errorf("Objective %s renders out of ascending ID order:\n%s", id, got)
		}
		previous = at
	}

	for _, want := range []string{"Finished and", "Done", "In Progress", "Planned"} {
		if !strings.Contains(got, want) {
			t.Errorf("sidebar is missing %q, which every row must carry:\n%s", want, got)
		}
	}
}

// TestSidebarShowsEachObjectivesOwnCheckBadge proves the sidebar's Check badge
// is the deliberately simple three-notch signal objectiveCheckBadge defines:
// a current Check reads "[✓] Check"; a never-checked Objective reads grey
// "[ ] Check"; and needs_work, unknown, and stale alike collapse into the
// same "[!] Check (needs work)" rather than the fuller five-state vocabulary
// a Task card's badge carries. The Full Objective Check is never waivable, so
// there is no separate waived notch here.
func TestSidebarShowsEachObjectivesOwnCheckBadge(t *testing.T) {
	got := sidebarText(t, sidebarBoard(t, writeNavigationProject(t)))

	// O-001 current; O-002/O-003 missing; O-004 needs_work, O-005 unknown, O-006
	// stale all fold into the same "checked, not clear" badge.
	if !strings.Contains(got, "[✓] Check") {
		t.Errorf("sidebar does not show a current Objective as checked:\n%s", got)
	}
	if count := strings.Count(got, "[ ] Check"); count != 2 {
		t.Errorf("sidebar shows %d Objectives as never checked, want 2 (O-002, O-003):\n%s", count, got)
	}
	// The label wraps across two lines at this column width ("Check (needs" /
	// "work)"), so this counts the badge's flagged glyph rather than the full
	// contiguous label text.
	if count := strings.Count(got, "[!]"); count != 3 {
		t.Errorf("sidebar shows %d Objectives as checked-not-clear, want 3 (needs_work, unknown, and stale all read the same):\n%s", count, got)
	}
}

// TestSidebarSeparatesFinishedTasksFromAFinishedObjective is the state the
// sidebar exists to make visible: an Objective whose every Task is done but
// which has never had its own Check reads differently from one whose Check is
// current, and finished Tasks alone are never enough to read as checked. This
// used to be a second badge ("INTEGRATED" / "NEEDS INTEGRATION") that just
// restated the Check badge in scarier words (I-012); it is gone, and the Check
// badge alone carries this distinction now.
func TestSidebarSeparatesFinishedTasksFromAFinishedObjective(t *testing.T) {
	rows := sidebarLines(t, sidebarBoard(t, writeNavigationProject(t)))

	checked := rowsFor(rows, "O-001")
	neverChecked := rowsFor(rows, "O-002")

	if !strings.Contains(checked, "[✓] Check") {
		t.Errorf("an Objective whose Tasks are done and whose Check is current does not read as checked:\n%s", checked)
	}
	if !strings.Contains(neverChecked, "[ ] Check") {
		t.Errorf("an Objective whose Tasks are done but which has never been checked does not say so:\n%s", neverChecked)
	}
	if strings.Contains(neverChecked, "[✓] Check") {
		t.Errorf("an unchecked Objective reads as checked:\n%s", neverChecked)
	}
	for _, retired := range []string{"INTEGRATED", "NEEDS INTEGRATION"} {
		if strings.Contains(checked+neverChecked, retired) {
			t.Errorf("sidebar still renders the retired %q wording:\n%s\n%s", retired, checked, neverChecked)
		}
	}
}

// TestObjectiveRowBadgesFoldExceptionIntoCheck proves ObjectiveRow.badges()
// wires ByException into objectiveCheckBadge rather than leaving it unread: a
// row whose owner accepted an exception against a not-current Check shows the
// same "[✓] Check" a clear Check gets, not the not-clear badge its raw
// clearance state alone would produce.
func TestObjectiveRowBadgesFoldExceptionIntoCheck(t *testing.T) {
	row := ObjectiveRow{
		Objective:   &data.ObjectiveV2{ID: "O-001", Title: "Accepted with a known gap", Status: "in_progress"},
		Clearance:   data.Clearance{State: data.ClearanceNeedsWork},
		ByException: true,
	}
	badges := row.badges()
	if len(badges) != 1 {
		t.Fatalf("badges() = %d entries, want exactly 1 (the Check badge, no waits)", len(badges))
	}
	if got := badges[0].Text(); got != "[✓] Check" {
		t.Errorf("badge = %q, want an exception-accepted row to read exactly like a current Check", got)
	}
}

// TestSidebarNamesTheObjectiveAWaitIsOn proves the dependency wait comes from
// ResolveObjectiveDependency and names its target.
func TestSidebarNamesTheObjectiveAWaitIsOn(t *testing.T) {
	rows := sidebarLines(t, sidebarBoard(t, writeNavigationProject(t)))

	waiting := rowsFor(rows, "O-003")
	if !strings.Contains(waiting, "WAITS O-002") {
		t.Errorf("O-003 does not name the Objective it waits on:\n%s", waiting)
	}
	if independent := rowsFor(rows, "O-004"); strings.Contains(independent, "WAITS") {
		t.Errorf("an Objective with no declared dependency reports a wait:\n%s", independent)
	}
}

// TestSelectionFiltersColumnsByRecordedOwnership covers the filter and its
// source: the Tasks a selected Objective owns come from index.ObjectiveTasks,
// so a Task filed under another Objective's directory still follows the owner
// its own record names.
func TestSelectionFiltersColumnsByRecordedOwnership(t *testing.T) {
	root := writeNavigationProject(t)
	loaded := loadProject(root)
	if loaded.Failed() {
		t.Fatalf("fixture project did not load: %s", loaded.Diagnostic)
	}
	if got := loaded.State.Index.Tasks["T-004"].Objective; got != "O-003" {
		t.Fatalf("fixture T-004 is owned by %q, want O-003 declared on the record itself", got)
	}
	if path := loaded.State.Index.Tasks["T-004"].Source.Path; !strings.Contains(path, "O-001") {
		t.Fatalf("fixture T-004 lives at %q, want a path under another Objective's directory", path)
	}

	// The router selects O-003, so the board opens filtered to it.
	model := sidebarBoard(t, root)
	if got := cardIDsInView(model); !equalIDs(got, []string{"T-003", "T-004"}) {
		t.Errorf("columns show %v, want exactly the Tasks O-003 owns", got)
	}

	// Selecting O-001 from the sidebar filters to its Tasks — and T-004, whose
	// file sits in O-001's directory, is not one of them.
	model = press(t, model, "left", "up", "up")
	if model.SelectedObjective != "O-001" {
		t.Fatalf("SelectedObjective = %q after selecting the first row, want O-001", model.SelectedObjective)
	}
	if got := cardIDsInView(model); !equalIDs(got, []string{"T-001"}) {
		t.Errorf("columns show %v, want only the Task O-001's own records claim", got)
	}
}

func TestNoObjectiveFilterShowsEveryTaskInTheSelectedGoal(t *testing.T) {
	model := press(t, sidebarBoard(t, writeNavigationProject(t)), "left", "esc")

	if model.SelectedObjective != "" {
		t.Fatalf("SelectedObjective = %q after clearing, want nothing selected", model.SelectedObjective)
	}
	if got := cardIDsInView(model); !equalIDs(got, []string{"T-001", "T-002", "T-003", "T-004"}) {
		t.Errorf("columns show %v, want every Task in the selected Goal", got)
	}
}

// TestInitialSelectionPrecedence covers the three ways a board opens: the flag
// wins, the router is next, and nothing selected is the honest default.
func TestInitialSelectionPrecedence(t *testing.T) {
	root := writeNavigationProject(t)

	fromRouter := sidebarBoard(t, root)
	if fromRouter.SelectedObjective != "O-003" {
		t.Errorf("SelectedObjective = %q, want the router's O-003", fromRouter.SelectedObjective)
	}

	fromFlag := openBoard(t, root, "O-001")
	if fromFlag.SelectedObjective != "O-001" {
		t.Errorf("SelectedObjective = %q, want --objective to win over the router", fromFlag.SelectedObjective)
	}
	if got := cardIDsInView(fromFlag); !equalIDs(got, []string{"T-001"}) {
		t.Errorf("columns show %v, want the flagged Objective's Tasks", got)
	}

	writeRouterWithRelease(t, root, "R-001", "none", "none")
	unselected := sidebarBoard(t, root)
	if unselected.SelectedObjective != "" {
		t.Errorf("SelectedObjective = %q over a router selecting nothing, want nothing", unselected.SelectedObjective)
	}
	if got := len(cardIDsInView(unselected)); got != 4 {
		t.Errorf("columns show %d Tasks with nothing selected, want all 4", got)
	}
}

// TestRouterNamingAMissingObjectiveOpensTheBoardAnyway proves a stale router
// hint costs nothing: no selection, every Task rendered, and no refusal to
// open. The selection-diagnostic sentence itself (data.SelectionDiagnostic,
// still resolved onto State.Next and still what `savepoint resume` reports)
// is no longer echoed in the board's own one-line Next area — that panel is
// a glance at the resolved Task now, not a diagnostics surface.
func TestRouterNamingAMissingObjectiveOpensTheBoardAnyway(t *testing.T) {
	root := writeNavigationProject(t)
	writeFixtureRouter(t, root, "task", "O-009", "none")

	model := sidebarBoard(t, root)
	got := xansi.Strip(model.View())

	if model.SelectedObjective != "" {
		t.Errorf("SelectedObjective = %q, want nothing selected for a record the index does not have", model.SelectedObjective)
	}
	if ids := cardIDsInView(model); len(ids) != 4 {
		t.Errorf("columns show %v, want every Task while nothing is selected", ids)
	}
	if model.State.Next.SelectionDiagnostic == nil {
		t.Fatalf("State.Next carries no selection diagnostic for a missing router Objective")
	}
	if !strings.Contains(model.State.Next.SelectionDiagnostic.ID, "O-009") {
		t.Errorf("selection diagnostic = %+v, want it to name O-009", model.State.Next.SelectionDiagnostic)
	}
	if !strings.Contains(got, "OBJECTIVES") || !strings.Contains(got, "PLANNED (") {
		t.Errorf("board did not open over a stale router hint:\n%s", got)
	}
	if strings.Contains(got, diagnosticHeading) {
		t.Errorf("a stale router hint was reported as invalid project data:\n%s", got)
	}
}

// TestSidebarNavigationClampsAndIsIdempotent covers the cursor: it moves, it
// stops at both ends, it applies the row it lands on as the selection
// immediately — the V1 interaction, with no separate enter:select step — and
// a repeated press at an end changes nothing at all.
func TestSidebarNavigationClampsAndIsIdempotent(t *testing.T) {
	model := press(t, sidebarBoard(t, writeNavigationProject(t)), "left")

	if !model.SidebarFocused {
		t.Fatal("left did not move focus to the sidebar")
	}
	if model.ObjectiveCursor != 2 {
		t.Fatalf("ObjectiveCursor = %d, want the row holding the selected O-003", model.ObjectiveCursor)
	}

	atTop := press(t, model, "up", "up", "up", "up")
	if atTop.ObjectiveCursor != 0 {
		t.Errorf("ObjectiveCursor = %d after walking off the top, want 0", atTop.ObjectiveCursor)
	}
	if press(t, atTop, "up").View() != atTop.View() {
		t.Error("pressing up at the first Objective changed the board")
	}
	if atTop.SelectedObjective != "O-001" {
		t.Errorf("SelectedObjective = %q at the top row, want up to have applied it immediately", atTop.SelectedObjective)
	}

	atBottom := press(t, model, "down", "down", "down", "down", "down", "down")
	if atBottom.ObjectiveCursor != 5 {
		t.Errorf("ObjectiveCursor = %d after walking off the bottom, want the last row", atBottom.ObjectiveCursor)
	}
	if press(t, atBottom, "down").View() != atBottom.View() {
		t.Error("pressing down at the last Objective changed the board")
	}

	// Moving the cursor applies the row it lands on, immediately.
	if atBottom.SelectedObjective != "O-006" {
		t.Errorf("SelectedObjective = %q after moving the cursor, want it to follow the cursor to the last row", atBottom.SelectedObjective)
	}
}

// TestLeftArrowAtPlannedColumnEntersSidebar proves arrow keys alone can reach
// every column end to end: left from the Planned column, the board's
// leftmost surface, crosses into the sidebar.
func TestLeftArrowAtPlannedColumnEntersSidebar(t *testing.T) {
	viaLeft := press(t, sidebarBoard(t, writeNavigationProject(t)), "left")

	if !viaLeft.SidebarFocused {
		t.Fatal("left at the Planned column did not move focus to the sidebar")
	}
	if viaLeft.ObjectiveCursor != 2 {
		t.Errorf("ObjectiveCursor = %d via left, want the row holding the router's selected O-003", viaLeft.ObjectiveCursor)
	}

	// h is the same key as left.
	viaH := press(t, sidebarBoard(t, writeNavigationProject(t)), "h")
	if !viaH.SidebarFocused {
		t.Error("h at the Planned column did not move focus to the sidebar")
	}

	// A column right of Planned still clamps rather than crossing.
	inProgress := press(t, sidebarBoard(t, writeNavigationProject(t)), "right", "left")
	if inProgress.SidebarFocused {
		t.Error("left from a column other than Planned crossed into the sidebar")
	}
}

// TestRightArrowInSidebarReturnsToColumns is the mirror crossing: right from
// the sidebar hands focus back to the Planned column, matching the edge left
// crossed in from.
func TestRightArrowInSidebarReturnsToColumns(t *testing.T) {
	inSidebar := press(t, sidebarBoard(t, writeNavigationProject(t)), "left")

	back := press(t, inSidebar, "right")
	if back.SidebarFocused {
		t.Fatal("right in the sidebar did not return focus to the columns")
	}
	if back.FocusedColumn != data.ColumnPlanned {
		t.Errorf("FocusedColumn = %q after right from the sidebar, want Planned", back.FocusedColumn)
	}

	// l is the same key as right.
	viaL := press(t, inSidebar, "l")
	if viaL.SidebarFocused {
		t.Error("l in the sidebar did not return focus to the columns")
	}
}

// TestNarrowTerminalLeftArrowStaysOnColumns proves the arrow-key crossing
// follows the same screen guard the sidebar itself does: a terminal too
// narrow for the sidebar leaves left/right clamping at the columns rather
// than crossing into a surface nothing draws.
func TestNarrowTerminalLeftArrowStaysOnColumns(t *testing.T) {
	model := openSizedBoard(t, writeNavigationProject(t), sidebarBreakpoint-1, 40)

	narrow := press(t, model, "left")
	if narrow.SidebarFocused {
		t.Error("left focused a sidebar the terminal is too narrow to draw")
	}
}

// TestSidebarSurvivesAReloadThatShortensTheList proves a cursor past the end of
// a reloaded list lands somewhere that renders.
func TestSidebarSurvivesAReloadThatShortensTheList(t *testing.T) {
	root := writeNavigationProject(t)
	model := press(t, sidebarBoard(t, root), "left", "down", "down", "down")
	if model.ObjectiveCursor != 5 {
		t.Fatalf("ObjectiveCursor = %d, want the last row", model.ObjectiveCursor)
	}

	// The three Objectives go, and so do the Checks that name them: a Check
	// pointing at a record that no longer exists is a load refusal, not a
	// shorter list.
	for _, id := range []string{"O-004", "O-005", "O-006"} {
		removeObjective(t, root, id)
	}
	for _, id := range []string{"C-003", "C-004", "C-005"} {
		removeCheck(t, root, id)
	}
	reloaded, _ := model.Update(loadCmd(root)().(projectLoadedMsg))
	after := reloaded.(Model)

	if after.ObjectiveCursor >= len(after.Objectives) {
		t.Errorf("ObjectiveCursor = %d over %d rows, want a cursor inside the list", after.ObjectiveCursor, len(after.Objectives))
	}
	if got := sidebarText(t, after); strings.Contains(got, "O-006") {
		t.Errorf("sidebar still lists an Objective the reload removed:\n%s", got)
	}
}

// TestSidebarScrollsRatherThanWrapping proves a list longer than the viewport
// reports what is out of view and keeps every line inside the sidebar's width.
func TestSidebarScrollsRatherThanWrapping(t *testing.T) {
	model := openSizedBoard(t, writeNavigationProject(t), 130, 24)

	lines := sidebarLines(t, model)
	got := strings.Join(lines, "\n")
	if !strings.Contains(got, "more") {
		t.Errorf("a sidebar taller than its viewport reports nothing out of view:\n%s", got)
	}
	for _, line := range lines {
		if lipgloss.Width(line) > sidebarWidth {
			t.Errorf("sidebar line is %d cells wide, past the %d it has: %q", lipgloss.Width(line), sidebarWidth, line)
		}
	}
	// A title wider than the sidebar wraps across up to two lines so the whole title can be read.
	full := strings.Join(sidebarLines(t, sidebarBoard(t, writeNavigationProject(t))), "\n")
	if !strings.Contains(full, "O-006") {
		t.Fatalf("the last Objective is not on screen, so nothing here proves how its title renders:\n%s", full)
	}
	if !strings.Contains(full, "Integration has gone") || !strings.Contains(full, "stale") {
		t.Errorf("a title longer than one line should wrap across two lines:\n%s", full)
	}
}

func TestRenderObjectiveRow_WrapsUpToTwoLinesAndTruncates(t *testing.T) {
	shortRow := ObjectiveRow{
		Objective: &data.ObjectiveV2{ID: "O-001", Title: "Short", Status: "planned"},
		Clearance: data.Clearance{State: data.ClearanceMissing},
	}
	shortText := xansi.Strip(renderObjectiveRow(shortRow, sidebarWidth, false, false))
	if strings.Contains(shortText, "…") {
		t.Errorf("short objective title should not be truncated:\n%s", shortText)
	}

	twoLineRow := ObjectiveRow{
		Objective: &data.ObjectiveV2{ID: "O-002", Title: "Implement user authentication subsystem", Status: "in_progress"},
		Clearance: data.Clearance{State: data.ClearanceMissing},
	}
	twoLineText := xansi.Strip(renderObjectiveRow(twoLineRow, sidebarWidth, false, false))
	if !strings.Contains(twoLineText, "Implement user") {
		t.Errorf("two-line objective title missing line 1:\n%s", twoLineText)
	}
	if !strings.Contains(twoLineText, "subsystem") {
		t.Errorf("two-line objective title missing line 2:\n%s", twoLineText)
	}

	longRow := ObjectiveRow{
		Objective: &data.ObjectiveV2{ID: "O-003", Title: "Implement user authentication and authorization subsystem with multi factor security and token lifecycle management", Status: "in_progress"},
		Clearance: data.Clearance{State: data.ClearanceMissing},
	}
	longText := xansi.Strip(renderObjectiveRow(longRow, 28, false, false))
	if !strings.Contains(longText, "…") {
		t.Errorf("objective title exceeding two lines should truncate with ellipsis:\n%s", longText)
	}
}

func TestEmptyProjectRendersAnEmptySidebar(t *testing.T) {
	model := sidebarBoard(t, writeEmptyProjectFromTemplate(t))
	got := xansi.Strip(model.View())

	if len(model.Objectives) != 0 {
		t.Errorf("Objectives = %v over a project with none", model.Objectives)
	}
	if !strings.Contains(sidebarText(t, model), sidebarTitle) {
		t.Errorf("no sidebar rendered for an empty project:\n%s", got)
	}
	for _, forbidden := range []string{diagnosticHeading, "error", "not found"} {
		if strings.Contains(got, forbidden) {
			t.Errorf("an empty project's sidebar reads as broken, containing %q:\n%s", forbidden, got)
		}
	}
	// Moving over an empty list is a no-op, not a panic.
	empty := press(t, model, "left", "down", "up")
	if empty.SelectedObjective != "" || empty.ObjectiveCursor != 0 {
		t.Errorf("navigating an empty sidebar selected %q at cursor %d", empty.SelectedObjective, empty.ObjectiveCursor)
	}
}

// TestFocusChangesAccentsNotGeometry is the focus half of the visual rule: the
// board occupies exactly the same cells with the sidebar focused as with the
// columns focused, and differs only in accent and marker glyphs.
func TestFocusChangesAccentsNotGeometry(t *testing.T) {
	forceColorProfile(t, termenv.TrueColor)
	model := sidebarBoard(t, writeNavigationProject(t))

	columns := model.View()
	sidebar := press(t, model, "left").View()

	if lipgloss.Width(columns) != lipgloss.Width(sidebar) {
		t.Errorf("width with the sidebar focused = %d, with the columns focused = %d", lipgloss.Width(sidebar), lipgloss.Width(columns))
	}
	if lipgloss.Height(columns) != lipgloss.Height(sidebar) {
		t.Errorf("height with the sidebar focused = %d, with the columns focused = %d", lipgloss.Height(sidebar), lipgloss.Height(columns))
	}
	if columns == sidebar {
		t.Error("moving focus changed nothing; the accent and the cursor glyph must change")
	}

	columnLines := strings.Split(xansi.Strip(columns), "\n")
	sidebarFocusedLines := strings.Split(xansi.Strip(sidebar), "\n")
	for i := range columnLines {
		if lipgloss.Width(columnLines[i]) != lipgloss.Width(sidebarFocusedLines[i]) {
			t.Errorf("line %d is %d cells with the columns focused and %d with the sidebar focused",
				i, lipgloss.Width(columnLines[i]), lipgloss.Width(sidebarFocusedLines[i]))
		}
	}
	if !strings.Contains(xansi.Strip(sidebar), glyphCursor) {
		t.Error("the focused sidebar shows no cursor glyph; focus must be legible without color")
	}
	if strings.Contains(xansi.Strip(columns), glyphCursor) {
		t.Error("the unfocused sidebar shows a cursor glyph")
	}
}

// TestNarrowTerminalOffersNoSidebarFocus proves the keys follow the screen: a
// terminal too narrow to draw the sidebar beside three columns does not move
// focus into a surface nothing is drawing.
func TestNarrowTerminalOffersNoSidebarFocus(t *testing.T) {
	model := openSizedBoard(t, writeNavigationProject(t), sidebarBreakpoint-1, 40)

	if strings.Contains(xansi.Strip(model.View()), sidebarTitle) {
		t.Fatalf("a terminal below the sidebar breakpoint drew one:\n%s", xansi.Strip(model.View()))
	}

	narrow := press(t, model, "left")
	if narrow.SidebarFocused {
		t.Error("left focused a sidebar the terminal is too narrow to draw")
	}

	// A terminal that shrinks while the sidebar has focus hands it back.
	wide := press(t, sidebarBoard(t, writeNavigationProject(t)), "left")
	shrunk, _ := wide.Update(tea.WindowSizeMsg{Width: sidebarBreakpoint - 1, Height: 40})
	if shrunk.(Model).SidebarFocused {
		t.Error("focus stayed on the sidebar after the terminal shrank past the breakpoint")
	}
	if narrow.View() != model.View() {
		t.Error("left changed a board with no sidebar on it")
	}
	if strings.Contains(xansi.Strip(narrow.View()), "tab") {
		t.Error("a board with no sidebar offers a Tab key that no longer exists")
	}
}

// TestNavigationWritesNothingToTheProject holds the task's no-write rule with a
// byte-and-mtime snapshot around a full navigation and selection sequence.
// Selection here is view state; the key that records one in router.md arrives
// with the rest of the write path.
func TestNavigationWritesNothingToTheProject(t *testing.T) {
	root := writeNavigationProject(t)
	model := sidebarBoard(t, root)
	before := snapshotProject(t, root)

	model = press(t, model,
		"left", "down", "down", "up", "up", "up",
		"esc", "right", "down", "left", "down")
	_ = model.View()

	after := snapshotProject(t, root)
	if len(before) != len(after) {
		t.Fatalf("project holds %d files after navigating, want the %d it started with", len(after), len(before))
	}
	for path, state := range before {
		if got, ok := after[path]; !ok || got != state {
			t.Errorf("%s changed under navigation:\n before %s\n  after %s", path, state, got)
		}
	}
}

// snapshotProject records every file under root by content and modification
// time, so a write of identical bytes is caught as well as a changed one.
func snapshotProject(t *testing.T, root string) map[string]string {
	t.Helper()
	snapshot := map[string]string{}
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		snapshot[path] = fmt.Sprintf("modified %s, content %q",
			info.ModTime().UTC().Format("2006-01-02T15:04:05.000000000Z"), content)
		return nil
	})
	if err != nil {
		t.Fatalf("snapshot %s: %v", root, err)
	}
	if len(snapshot) == 0 {
		t.Fatalf("snapshot of %s is empty; the assertion would pass vacuously", root)
	}
	return snapshot
}

// cardIDsInView lists the Task IDs the three columns are showing, sorted, since
// what an ownership filter changes is which Tasks are there rather than which
// column each sits in.
func cardIDsInView(model Model) []string {
	var ids []string
	for _, column := range columnOrder {
		ids = append(ids, cardIDs(model.Cards[column])...)
	}
	slices.Sort(ids)
	return ids
}

func equalIDs(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}

// rowsFor returns the sidebar lines belonging to one Objective: the line naming
// it and every line up to the next Objective's.
func rowsFor(lines []string, id string) string {
	var block []string
	collecting := false
	for _, line := range lines {
		switch {
		case strings.Contains(line, id):
			collecting = true
		case collecting && objectiveRowStart(line):
			return strings.Join(block, "\n")
		}
		if collecting {
			block = append(block, line)
		}
	}
	return strings.Join(block, "\n")
}

// objectiveRowStart reports whether a sidebar line is the first line of a row,
// which is the only line carrying an O-### identity.
func objectiveRowStart(line string) bool {
	trimmed := strings.TrimLeft(strings.TrimPrefix(strings.TrimSpace(line), "│"), " "+glyphCursor+glyphSelected)
	return strings.HasPrefix(trimmed, "O") && len(trimmed) > 4 && trimmed[1] >= '0' && trimmed[1] <= '9'
}

// writeIssueOnlyRouter selects issue alone, the way an Issue repair is routed.
func writeIssueOnlyRouter(t *testing.T, root, issue string) {
	t.Helper()
	testutil.WriteFile(t, filepath.Join(root, "router.md"),
		"# Router\n\n## Current state\n\n```yaml\nstate: task\nrelease: R-001\nobjective: none\ntask: none\nissue: "+issue+"\n```\n")
}

// TestIssueOnlySelectionOpensOnTheIssuesObjective covers I-045: a router that
// selects an Issue alone opens the board on the one Objective the Issue's
// linked Tasks and Objective-scoped Checks belong to, and on the labelled
// unfiltered view when that is not exactly one Objective.
func TestIssueOnlySelectionOpensOnTheIssuesObjective(t *testing.T) {
	tests := []struct {
		name        string
		task, check string
		want        string
	}{
		// C-002 is scoped to Task T-001, so it names no Objective.
		{name: "task link only", task: "T-003", check: "C-002", want: "O-003"},
		{name: "task and objective check agree", task: "T-002", check: "C-002", want: "O-002"},
		// C-003 is scoped to Objective O-004; T-003 belongs to O-003.
		{name: "two objectives", task: "T-003", check: "C-003", want: ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := writeNavigationProject(t)
			writeIssue(t, root, "I-001", "Repair something", "defect", test.check, test.task)
			writeIssueOnlyRouter(t, root, "I-001")

			model := sidebarBoard(t, root)
			if model.SelectedObjective != test.want {
				t.Fatalf("SelectedObjective = %q, want %q", model.SelectedObjective, test.want)
			}
			labelled := strings.Contains(xansi.Strip(model.View()), allObjectivesLabel)
			if labelled != (test.want == "") {
				t.Errorf("%s shown = %v with SelectedObjective %q", allObjectivesLabel, labelled, test.want)
			}

			var plain bytes.Buffer
			if err := Run(Options{Root: root, Stdout: &plain, TTY: false}); err != nil {
				t.Fatalf("Run() error = %v", err)
			}
			wantLine := "Selected: " + test.want
			if test.want == "" {
				wantLine = "Selected: all Objectives in Goal R-001"
			}
			if !strings.Contains(plain.String(), wantLine+"\n") {
				t.Errorf("plain output is missing %q:\n%s", wantLine, plain.String())
			}
		})
	}
}

// TestIssueOnlySelectionWithUnknownIssueOpensUnfiltered keeps a router naming
// a missing Issue on the labelled unfiltered view instead of guessing one.
func TestIssueOnlySelectionWithUnknownIssueOpensUnfiltered(t *testing.T) {
	root := writeNavigationProject(t)
	writeIssueOnlyRouter(t, root, "I-404")

	model := sidebarBoard(t, root)
	if model.SelectedObjective != "" {
		t.Fatalf("SelectedObjective = %q, want nothing for a missing Issue", model.SelectedObjective)
	}
	if !strings.Contains(xansi.Strip(model.View()), allObjectivesLabel) {
		t.Errorf("unfiltered view is missing %s", allObjectivesLabel)
	}
}

// TestFilteredViewHasNoAllObjectivesLabel keeps the label to the unfiltered
// view: a board filtered to one Objective never claims to show them all.
func TestFilteredViewHasNoAllObjectivesLabel(t *testing.T) {
	model := sidebarBoard(t, writeNavigationProject(t))
	if model.SelectedObjective == "" {
		t.Fatal("fixture router should select an Objective")
	}
	if strings.Contains(xansi.Strip(model.View()), allObjectivesLabel) {
		t.Errorf("filtered view shows %s", allObjectivesLabel)
	}
	model.selectObjective("")
	if !strings.Contains(xansi.Strip(model.View()), allObjectivesLabel) {
		t.Errorf("view after clearing the filter is missing %s", allObjectivesLabel)
	}
}
