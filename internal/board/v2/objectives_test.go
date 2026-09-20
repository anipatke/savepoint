package v2

import (
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
)

// sidebarBoard opens a board wide enough to carry the sidebar beside the three
// columns, and tall enough for every fixture Objective to be on screen at once
// below the Next area, which is the size every sidebar assertion below is made
// at.
func sidebarBoard(t *testing.T, root string) Model {
	t.Helper()
	return openSizedBoard(t, root, 130, 48)
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
// in the index, in stable O### order, each row carrying its ID, its title, and
// its recorded status.
func TestSidebarListsEveryObjectiveInOrder(t *testing.T) {
	got := sidebarText(t, sidebarBoard(t, writeNavigationProject(t)))

	if !strings.Contains(got, sidebarTitle) {
		t.Errorf("no Objective sidebar rendered:\n%s", got)
	}

	previous := -1
	for _, id := range []string{"O001", "O002", "O003", "O004", "O005", "O006"} {
		at := strings.Index(got, id)
		if at < 0 {
			t.Fatalf("sidebar is missing Objective %s:\n%s", id, got)
		}
		if at < previous {
			t.Errorf("Objective %s renders out of ascending ID order:\n%s", id, got)
		}
		previous = at
	}

	for _, want := range []string{"Finished and in", "done", "in_progress", "planned"} {
		if !strings.Contains(got, want) {
			t.Errorf("sidebar is missing %q, which every row must carry:\n%s", want, got)
		}
	}
}

// TestSidebarShowsEachObjectivesOwnCheckBadge proves the sidebar's Check badge
// is the deliberately simple two-notch signal objectiveCheckBadge defines: a
// current Check reads "[✓] Check", and every other clearance state — missing,
// needs_work, unknown, and stale alike — reads the same grey "[ ] Check"
// rather than the fuller five-state vocabulary a Task card's badge carries.
// The Full Objective Check is never waivable, so there is no third notch here.
func TestSidebarShowsEachObjectivesOwnCheckBadge(t *testing.T) {
	got := sidebarText(t, sidebarBoard(t, writeNavigationProject(t)))

	// O001 current; O002/O003 missing, O004 needs work, O005 unknown, O006
	// stale all fold into the same "not yet" badge.
	if !strings.Contains(got, "[✓] Check") {
		t.Errorf("sidebar does not show a current Objective as checked:\n%s", got)
	}
	if count := strings.Count(got, "[ ] Check"); count != 5 {
		t.Errorf("sidebar shows %d Objectives as not-yet-checked, want 5 (missing, needs_work, unknown, and stale all read the same):\n%s", count, got)
	}
}

// TestSidebarSeparatesFinishedTasksFromAFinishedObjective is the state the
// sidebar exists to make visible: an Objective whose every Task is done but
// whose own integration Check is not current looks finished in the columns and
// must not look finished here.
func TestSidebarSeparatesFinishedTasksFromAFinishedObjective(t *testing.T) {
	rows := sidebarLines(t, sidebarBoard(t, writeNavigationProject(t)))

	integrated := rowsFor(rows, "O001")
	waiting := rowsFor(rows, "O002")

	if !strings.Contains(integrated, "✓ INTEGRATED") {
		t.Errorf("an Objective whose Tasks are done and whose integration is current does not read as integrated:\n%s", integrated)
	}
	if !strings.Contains(waiting, "⚠ NEEDS INTEGRATION") {
		t.Errorf("an Objective whose Tasks are done but whose integration is not current does not say so:\n%s", waiting)
	}
	if strings.Contains(waiting, "✓ INTEGRATED") {
		t.Errorf("an uncleared Objective reads as integrated:\n%s", waiting)
	}
}

// TestSidebarNamesTheObjectiveAWaitIsOn proves the dependency wait comes from
// ResolveObjectiveDependency and names its target.
func TestSidebarNamesTheObjectiveAWaitIsOn(t *testing.T) {
	rows := sidebarLines(t, sidebarBoard(t, writeNavigationProject(t)))

	waiting := rowsFor(rows, "O003")
	if !strings.Contains(waiting, "WAITS O002") {
		t.Errorf("O003 does not name the Objective it waits on:\n%s", waiting)
	}
	if independent := rowsFor(rows, "O004"); strings.Contains(independent, "WAITS") {
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
	if got := loaded.State.Index.Tasks["T004"].Objective; got != "O003" {
		t.Fatalf("fixture T004 is owned by %q, want O003 declared on the record itself", got)
	}
	if path := loaded.State.Index.Tasks["T004"].Source.Path; !strings.Contains(path, "O001") {
		t.Fatalf("fixture T004 lives at %q, want a path under another Objective's directory", path)
	}

	// The router selects O003, so the board opens filtered to it.
	model := sidebarBoard(t, root)
	if got := cardIDsInView(model); !equalIDs(got, []string{"T003", "T004"}) {
		t.Errorf("columns show %v, want exactly the Tasks O003 owns", got)
	}

	// Selecting O001 from the sidebar filters to its Tasks — and T004, whose
	// file sits in O001's directory, is not one of them.
	model = press(t, model, "tab", "up", "up", "enter")
	if model.SelectedObjective != "O001" {
		t.Fatalf("SelectedObjective = %q after selecting the first row, want O001", model.SelectedObjective)
	}
	if got := cardIDsInView(model); !equalIDs(got, []string{"T001"}) {
		t.Errorf("columns show %v, want only the Task O001's own records claim", got)
	}
}

func TestNoSelectionShowsEveryTaskInTheProject(t *testing.T) {
	model := press(t, sidebarBoard(t, writeNavigationProject(t)), "tab", "esc")

	if model.SelectedObjective != "" {
		t.Fatalf("SelectedObjective = %q after clearing, want nothing selected", model.SelectedObjective)
	}
	if got := cardIDsInView(model); !equalIDs(got, []string{"T001", "T002", "T003", "T004"}) {
		t.Errorf("columns show %v, want every Task in the project", got)
	}
}

// TestInitialSelectionPrecedence covers the three ways a board opens: the flag
// wins, the router is next, and nothing selected is the honest default.
func TestInitialSelectionPrecedence(t *testing.T) {
	root := writeNavigationProject(t)

	fromRouter := sidebarBoard(t, root)
	if fromRouter.SelectedObjective != "O003" {
		t.Errorf("SelectedObjective = %q, want the router's O003", fromRouter.SelectedObjective)
	}

	fromFlag := openBoard(t, root, "O001")
	if fromFlag.SelectedObjective != "O001" {
		t.Errorf("SelectedObjective = %q, want --objective to win over the router", fromFlag.SelectedObjective)
	}
	if got := cardIDsInView(fromFlag); !equalIDs(got, []string{"T001"}) {
		t.Errorf("columns show %v, want the flagged Objective's Tasks", got)
	}

	writeRouter(t, root, "design", "none", "none")
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
	writeRouter(t, root, "task", "O009", "none")

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
	if !strings.Contains(model.State.Next.SelectionDiagnostic.ID, "O009") {
		t.Errorf("selection diagnostic = %+v, want it to name O009", model.State.Next.SelectionDiagnostic)
	}
	if !strings.Contains(got, "OBJECTIVES") || !strings.Contains(got, "PLANNED (") {
		t.Errorf("board did not open over a stale router hint:\n%s", got)
	}
	if strings.Contains(got, diagnosticHeading) {
		t.Errorf("a stale router hint was reported as invalid project data:\n%s", got)
	}
}

// TestSidebarNavigationClampsAndIsIdempotent covers the cursor: it moves, it
// stops at both ends, and a repeated press at an end changes nothing at all.
func TestSidebarNavigationClampsAndIsIdempotent(t *testing.T) {
	model := press(t, sidebarBoard(t, writeNavigationProject(t)), "tab")

	if !model.SidebarFocused {
		t.Fatal("tab did not move focus to the sidebar")
	}
	if model.ObjectiveCursor != 2 {
		t.Fatalf("ObjectiveCursor = %d, want the row holding the selected O003", model.ObjectiveCursor)
	}

	atTop := press(t, model, "up", "up", "up", "up")
	if atTop.ObjectiveCursor != 0 {
		t.Errorf("ObjectiveCursor = %d after walking off the top, want 0", atTop.ObjectiveCursor)
	}
	if press(t, atTop, "up").View() != atTop.View() {
		t.Error("pressing up at the first Objective changed the board")
	}

	atBottom := press(t, model, "down", "down", "down", "down", "down", "down")
	if atBottom.ObjectiveCursor != 5 {
		t.Errorf("ObjectiveCursor = %d after walking off the bottom, want the last row", atBottom.ObjectiveCursor)
	}
	if press(t, atBottom, "down").View() != atBottom.View() {
		t.Error("pressing down at the last Objective changed the board")
	}

	// Moving the cursor selects nothing by itself.
	if atBottom.SelectedObjective != "O003" {
		t.Errorf("SelectedObjective = %q after moving the cursor, want the selection unchanged", atBottom.SelectedObjective)
	}
}

// TestSidebarSurvivesAReloadThatShortensTheList proves a cursor past the end of
// a reloaded list lands somewhere that renders.
func TestSidebarSurvivesAReloadThatShortensTheList(t *testing.T) {
	root := writeNavigationProject(t)
	model := press(t, sidebarBoard(t, root), "tab", "down", "down", "down")
	if model.ObjectiveCursor != 5 {
		t.Fatalf("ObjectiveCursor = %d, want the last row", model.ObjectiveCursor)
	}

	// The three Objectives go, and so do the Checks that name them: a Check
	// pointing at a record that no longer exists is a load refusal, not a
	// shorter list.
	for _, id := range []string{"O004", "O005", "O006"} {
		removeObjective(t, root, id)
	}
	for _, id := range []string{"C003", "C004", "C005"} {
		removeCheck(t, root, id)
	}
	reloaded, _ := model.Update(loadCmd(root)().(projectLoadedMsg))
	after := reloaded.(Model)

	if after.ObjectiveCursor >= len(after.Objectives) {
		t.Errorf("ObjectiveCursor = %d over %d rows, want a cursor inside the list", after.ObjectiveCursor, len(after.Objectives))
	}
	if got := sidebarText(t, after); strings.Contains(got, "O006") {
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
	// A title wider than the sidebar costs one line, truncated, not three.
	full := strings.Join(sidebarLines(t, sidebarBoard(t, writeNavigationProject(t))), "\n")
	if !strings.Contains(full, "O006") {
		t.Fatalf("the last Objective is not on screen, so nothing here proves how its title renders:\n%s", full)
	}
	if strings.Contains(full, "Integration has gone stale") {
		t.Errorf("a title longer than the sidebar wrapped instead of truncating:\n%s", full)
	}
	if !strings.Contains(full, "…") {
		t.Errorf("no title was truncated, so the width rule is untested:\n%s", full)
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
	// Moving and selecting over an empty list is a no-op, not a panic.
	empty := press(t, model, "tab", "down", "enter", "up", "enter")
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
	sidebar := press(t, model, "tab").View()

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

	narrow := press(t, model, "tab")
	if narrow.SidebarFocused {
		t.Error("tab focused a sidebar the terminal is too narrow to draw")
	}

	// A terminal that shrinks while the sidebar has focus hands it back.
	wide := press(t, sidebarBoard(t, writeNavigationProject(t)), "tab")
	shrunk, _ := wide.Update(tea.WindowSizeMsg{Width: sidebarBreakpoint - 1, Height: 40})
	if shrunk.(Model).SidebarFocused {
		t.Error("focus stayed on the sidebar after the terminal shrank past the breakpoint")
	}
	if narrow.View() != model.View() {
		t.Error("tab changed a board with no sidebar on it")
	}
	if strings.Contains(xansi.Strip(narrow.View()), "tab:") {
		t.Error("a board with no sidebar offers the key that would focus it")
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
		"tab", "down", "down", "enter", "up", "up", "up", "enter",
		"esc", "tab", "right", "down", "left", "tab", "down", "enter")
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
// which is the only line carrying an O### identity.
func objectiveRowStart(line string) bool {
	trimmed := strings.TrimLeft(strings.TrimPrefix(strings.TrimSpace(line), "│"), " "+glyphCursor+glyphSelected)
	return strings.HasPrefix(trimmed, "O") && len(trimmed) > 4 && trimmed[1] >= '0' && trimmed[1] <= '9'
}
