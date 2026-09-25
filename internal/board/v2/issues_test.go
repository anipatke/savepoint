package v2

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	xansi "github.com/charmbracelet/x/ansi"
	"github.com/muesli/termenv"
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/styles"
	"github.com/opencode/savepoint/internal/testutil"
)

// writeIssuesProject covers all five descriptive types, all three statuses,
// each resolved disposition, a reopened history, and links through both a
// Task and a Check. Every assertion still runs against a temporary project.
func writeIssuesProject(t *testing.T) string {
	t.Helper()
	root := savepointRoot(t)
	writeConfig(t, root)
	writeFixtureRouter(t, root, "task", "O-001", "T-001")
	writeFixtureObjective(t, root, "O-001", "Issue surface", "in_progress", "")
	writeTask(t, root, "O-001", "T-001", "Task carrying follow-ups", "status: in_progress\nstage: audit\n")
	writeCheck(t, root, "C-001", "task", "T-001", "CLEAR")

	writeBoardIssue(t, root, "I-001", "A repair is needed", "defect", "open", `source: {kind: check, check: C-001, actor: {role: checker, session: issue-fixture}, at: '2026-01-01T00:00:00Z'}
tasks: [T-001]
checks: [C-001]
guardrail_ids: [FS-01]
severity: high
history:
  - at: '2026-01-01T00:00:00Z'
    actor: {role: checker, session: issue-fixture}
    kind: observed
    note: first observation
    check: C-001
  - at: '2026-01-02T00:00:00Z'
    actor: {role: executor, session: repair-fixture}
    kind: reopened
    note: the same symptom returned
`)
	writeBoardIssue(t, root, "I-002", "A recorded drift", "drift", "in_progress", `source: {kind: report, actor: {role: owner, session: issue-fixture}, at: '2026-01-03T00:00:00Z'}
`)
	writeBoardIssue(t, root, "I-003", "An accepted risk", "guardrail", "resolved", `source: {kind: report, actor: {role: owner, session: issue-fixture}, at: '2026-01-04T00:00:00Z'}
resolution: {disposition: accepted, actor: {role: owner, session: owner-fixture}, at: '2026-01-05T00:00:00Z', reason: 'the risk is deliberately accepted'}
`)
	writeBoardIssue(t, root, "I-004", "A duplicate report", "verification", "resolved", `source: {kind: report, actor: {role: checker, session: issue-fixture}, at: '2026-01-06T00:00:00Z'}
duplicate_of: I-001
resolution: {disposition: duplicate, actor: {role: checker, session: checker-fixture}, at: '2026-01-06T01:00:00Z'}
`)
	writeBoardIssue(t, root, "I-005", "A verified repair", "other", "resolved", `source: {kind: report, actor: {role: checker, session: issue-fixture}, at: '2026-01-07T00:00:00Z'}
tasks: [T-001]
checks: [C-001]
resolution: {disposition: verified, check: C-001, actor: {role: checker, session: checker-fixture}, at: '2026-01-07T01:00:00Z'}
`)
	return root
}

func writeBoardIssue(t *testing.T, root, id, title, issueType, status, fields string) {
	t.Helper()
	content := fmt.Sprintf("---\nid: %s\ntitle: %q\ntype: %s\nstatus: %s\n%s---\n\n## Summary\n\nThe recorded summary for %s.\n", id, title, issueType, status, fields, id)
	testutil.WriteFile(t, filepath.Join(root, "issues", id+"-fixture.md"), content)
}

func issueBoard(t *testing.T, root string) Model {
	t.Helper()
	return openSizedBoard(t, root, 120, 70)
}

func issueScreen(model Model) string {
	return screen(model)
}

func TestIssuesListUsesStableIdentityOrderAndShowsRecordedSeverity(t *testing.T) {
	model := press(t, issueBoard(t, writeIssuesProject(t)), "i")
	got := issueScreen(model)

	previous := -1
	for _, id := range []string{"I-001", "I-002", "I-003", "I-004", "I-005"} {
		at := strings.Index(got, id)
		if at < 0 {
			t.Fatalf("Issues overlay is missing %s:\n%s", id, got)
		}
		if at < previous {
			t.Errorf("Issue %s is out of stable ID order:\n%s", id, got)
		}
		previous = at
	}
	for _, want := range []string{"Filter: ALL", "✗ DEFECT", "OPEN (1)", "▲ HIGH", "A repair is needed"} {
		if !strings.Contains(got, want) {
			t.Errorf("Issues overlay is missing %q:\n%s", want, got)
		}
	}
}

func TestIssuesFilterCyclesThroughEveryTypeAndBackToAll(t *testing.T) {
	model := press(t, issueBoard(t, writeIssuesProject(t)), "i")
	for _, want := range []string{"DEFECT", "DRIFT", "GUARDRAIL", "VERIFICATION", "OTHER"} {
		model = press(t, model, "f")
		got := issueScreen(model)
		if !strings.Contains(got, "Filter: "+want) {
			t.Fatalf("filter is missing %q:\n%s", want, got)
		}
		if strings.Contains(got, "Filter: "+want+"\n") && want == "DRIFT" && strings.Contains(got, "I-001") {
			t.Errorf("drift filter included the defect row:\n%s", got)
		}
	}
	model = press(t, model, "f")
	if !strings.Contains(issueScreen(model), "Filter: ALL") {
		t.Errorf("filter did not cycle back to ALL:\n%s", issueScreen(model))
	}
}

func TestIssueDetailShowsOriginLinksGuardrailsResolutionAndHistory(t *testing.T) {
	root := writeIssuesProject(t)
	model := press(t, issueBoard(t, root), "i", "enter")
	got := issueScreen(model)
	for _, want := range []string{
		"ISSUE DETAIL", "ID: I-001", "Type: defect", "Status: open",
		"SUMMARY", "The recorded summary for I-001.", "Kind: check", "Check: C-001",
		"Actor: checker session issue-fixture", "LINKED TASKS", "T-001 — Task carrying follow-ups",
		"LINKED CHECKS", "C-001 — CLEAR", "GUARDRAILS", "FS-01", "HISTORY",
		"first observation", "the same symptom returned",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("Issue detail is missing %q:\n%s", want, got)
		}
	}

	// I-001 is Open; I-003 (accepted) is the first of three Resolved rows, so
	// reaching it crosses into the Resolved column rather than scrolling down
	// a single shared list.
	accepted := press(t, model, "esc", "right", "right", "enter")
	acceptedText := issueScreen(accepted)
	if !strings.Contains(acceptedText, "Disposition: accepted") || !strings.Contains(acceptedText, "Not proof of repair") {
		t.Errorf("accepted disposition is not distinct from repair proof:\n%s", acceptedText)
	}

	duplicate := press(t, accepted, "esc", "down", "enter")
	duplicateText := issueScreen(duplicate)
	if !strings.Contains(duplicateText, "Disposition: duplicate") || !strings.Contains(duplicateText, "Canonical: I-001") {
		t.Errorf("duplicate detail does not name its canonical target:\n%s", duplicateText)
	}
	toCanonical := press(t, duplicate, "enter")
	if !strings.Contains(issueScreen(toCanonical), "ID: I-001") {
		t.Errorf("duplicate navigation did not open the canonical Issue:\n%s", issueScreen(toCanonical))
	}
	verified := press(t, issueBoard(t, root), "i", "right", "right", "down", "down", "enter")
	verifiedText := issueScreen(verified)
	if !strings.Contains(verifiedText, "Disposition: verified") || !strings.Contains(verifiedText, "Proof: Check C-001") {
		t.Errorf("verified disposition does not show its proof Check:\n%s", verifiedText)
	}
}

// TestIssuesSplitIntoThreeStatusColumns covers the Open/In Progress/Resolved
// column layout: all three headers with their counts render together, and
// left/right moves focus between columns while up/down stays inside one.
// TestIssuesColumnSortsMostSevereFirst covers the severity sort (blocker at
// the top of a column, cosmetic-or-unrecorded at the bottom) independent of
// the ID order the Issues were written in — I-901 is written last but
// declares the most severe word, so a stable-ID-only ordering would leave it
// at the bottom instead of the top.
func TestIssuesColumnSortsMostSevereFirst(t *testing.T) {
	root := savepointRoot(t)
	writeConfig(t, root)
	writeRouter(t, root, "task", "O-001", "T-001")
	writeObjective(t, root, "O-001", "Severity ordering", "in_progress")
	writeTask(t, root, "O-001", "T-001", "Task carrying follow-ups", "status: planned\n")

	writeBoardIssue(t, root, "I-801", "A low-severity item", "other", "open", `source: {kind: report, actor: {role: owner, session: severity-fixture}, at: '2026-01-01T00:00:00Z'}
severity: low
`)
	writeBoardIssue(t, root, "I-802", "An item with no severity", "other", "open", `source: {kind: report, actor: {role: owner, session: severity-fixture}, at: '2026-01-01T00:00:00Z'}
`)
	writeBoardIssue(t, root, "I-803", "A medium-severity item", "other", "open", `source: {kind: report, actor: {role: owner, session: severity-fixture}, at: '2026-01-01T00:00:00Z'}
severity: medium
`)
	writeBoardIssue(t, root, "I-901", "A blocking item written last", "other", "open", `source: {kind: report, actor: {role: owner, session: severity-fixture}, at: '2026-01-01T00:00:00Z'}
severity: blocker
`)

	got := issueScreen(press(t, issueBoard(t, root), "i"))
	previous := -1
	for _, id := range []string{"I-901", "I-803", "I-801", "I-802"} {
		at := strings.Index(got, id)
		if at < 0 {
			t.Fatalf("Issues overlay is missing %s:\n%s", id, got)
		}
		if at < previous {
			t.Errorf("Issue %s is out of severity order (want blocker, medium, low, unrecorded):\n%s", id, got)
		}
		previous = at
	}
}

// TestIssueRowTitleStyleGivesEveryColumnItsOwnDistinctSelectionAccent proves
// a selected row's title always changes color from the plain unselected
// style — In Progress orange, Resolved green, and Open the Issues
// red, since Open cannot reuse a Task card's plain-selected-Planned
// choice: that one stays legible because the card still gets its own
// bordered box, which an Issue row does not have.
func TestIssueRowTitleStyleGivesEveryColumnItsOwnDistinctSelectionAccent(t *testing.T) {
	forceColorProfile(t, termenv.TrueColor)
	sample := "x"

	plain := issueRowTitleStyle(data.IssueStatusOpen, false).Render(sample)
	for _, status := range []data.IssueStatus{data.IssueStatusOpen, data.IssueStatusInProgress, data.IssueStatusResolved} {
		if got := issueRowTitleStyle(status, false).Render(sample); got != plain {
			t.Errorf("unselected %s title = %q, want every column's unselected title to render the same plain style (%q)", status, got, plain)
		}
	}

	selected := map[data.IssueStatus]string{
		data.IssueStatusOpen:       issueRowTitleStyle(data.IssueStatusOpen, true).Render(sample),
		data.IssueStatusInProgress: issueRowTitleStyle(data.IssueStatusInProgress, true).Render(sample),
		data.IssueStatusResolved:   issueRowTitleStyle(data.IssueStatusResolved, true).Render(sample),
	}
	seen := map[string]data.IssueStatus{plain: ""}
	for status, rendered := range selected {
		if rendered == plain {
			t.Errorf("selected %s title did not change color from the plain unselected style %q", status, plain)
		}
		if other, ok := seen[rendered]; ok {
			t.Errorf("selected %s and %s titles render identically: %q", status, other, rendered)
		}
		seen[rendered] = status
	}
}

// TestIssueColumnAccentMatchesItsOwnSelectedRow proves a focused column's
// heading always uses the same color its own selected row's title does —
// Open's heading must not stay the Planned column's grey while its selected
// row wears the Issues red; every column wears exactly one accent
// color, head to row. The heading is bold and the row title is not, so this
// compares foreground color alone rather than the full rendered style.
func TestIssueColumnAccentMatchesItsOwnSelectedRow(t *testing.T) {
	forceColorProfile(t, termenv.TrueColor)

	for _, status := range []data.IssueStatus{data.IssueStatusOpen, data.IssueStatusInProgress, data.IssueStatusResolved} {
		headingColor := issueColumnTitleStyle(status, true).GetForeground()
		rowColor := issueRowTitleStyle(status, true).GetForeground()
		if headingColor != rowColor {
			t.Errorf("%s column heading color = %v, selected row color = %v; want the same accent", status, headingColor, rowColor)
		}
	}

	// An unfocused column never carries a status accent at all — only the
	// column holding board focus does.
	if got := issueColumnTitleStyle(data.IssueStatusOpen, false); got.Render("x") != styles.ColumnTitle.Render("x") {
		t.Errorf("unfocused Open heading = %q, want the plain unaccented column title style", got.Render("x"))
	}
}

// TestIssueSurfaceWearsItsOwnRedAccent proves the Issues surface does not
// read as the Task board: every Issue ID, in a row and in the detail, both
// Issues headings, and the focused Open column's border wear the red accent, and that red is none of the
// board's own orange, green, or purple accents.
func TestIssueSurfaceWearsItsOwnRedAccent(t *testing.T) {
	forceColorProfile(t, termenv.TrueColor)

	issue := &data.IssueV2{ID: "I-031", Title: "Reload errors", Type: "defect", Status: data.IssueStatusOpen}
	redID := styles.IssueAccent.Render("I-031")

	// The row's selection marker sits inside the ID's styled span.
	if row := renderIssueRow(IssueRow{Issue: issue}, data.IssueStatusOpen, 40, false); !strings.Contains(row, styles.IssueAccent.Render("  I-031")) {
		t.Errorf("Issue row ID is not red:\n%q", row)
	}
	if lines := issueDetailLines(IssueDetail{Issue: issue}, 40); !strings.Contains(strings.Join(lines, "\n"), redID) {
		t.Errorf("Issue detail ID is not red:\n%q", strings.Join(lines, "\n"))
	}
	if header := issuesHeaderLine(IssueOverlay{}); !strings.Contains(header, styles.IssueAccent.Render("ISSUES")) {
		t.Errorf("Issues heading is not red: %q", header)
	}
	if detail := renderIssueDetail(IssueDetail{Issue: issue}, 60, 20, 0); !strings.Contains(detail, styles.IssueAccent.Render("ISSUE DETAIL")) {
		t.Errorf("Issue detail heading is not red:\n%q", detail)
	}
	if got, want := issueColumnStyle(data.IssueStatusOpen, true).GetBorderTopForeground(), styles.IssueAccent.GetForeground(); got != want {
		t.Errorf("focused Open column border = %v, want the Issues red %v", got, want)
	}

	red := styles.IssueAccent.GetForeground()
	for name, other := range map[string]lipgloss.Style{
		"orange": styles.ColumnTitleFocused,
		"green":  styles.ColumnTitleFocusedDone,
		"purple": styles.SidebarTitleFocused,
	} {
		if other.GetForeground() == red {
			t.Errorf("Issue accent shares the board's %s accent", name)
		}
	}
}

// TestIssueColumnsGroupByTypeUnderTheAllFilter proves the ALL filter orders
// each column by type in filter order, then severity within a type, and
// heads each type group once with its count; a type filter leaves one type,
// so it draws no group headings.
func TestIssueColumnsGroupByTypeUnderTheAllFilter(t *testing.T) {
	issue := func(id string, issueType data.IssueType, severity string) IssueRow {
		return IssueRow{Issue: &data.IssueV2{ID: id, Title: id, Type: issueType, Severity: severity, Status: data.IssueStatusOpen}}
	}
	model := Model{
		Issues: &IssueOverlay{FocusedStatus: data.IssueStatusOpen},
		State: ProjectState{Issues: IssueCatalog{Rows: []IssueRow{
			issue("I-001", "other", ""),
			issue("I-002", "defect", "low"),
			issue("I-003", "drift", ""),
			issue("I-004", "defect", "blocker"),
			issue("I-005", "verification", ""),
		}}},
	}

	var got []string
	for _, row := range model.focusedIssueRows() {
		got = append(got, row.Issue.ID)
	}
	if want := []string{"I-004", "I-002", "I-003", "I-005", "I-001"}; strings.Join(got, " ") != strings.Join(want, " ") {
		t.Errorf("ALL filter order = %v, want %v (type in filter order, then severity)", got, want)
	}

	column := xansi.Strip(renderIssueColumn("OPEN", data.IssueStatusOpen, model.focusedIssueRows(), 40, 60,
		columnCursor{Holds: true, Focused: true}, model.issuesGroupedByType()))
	for _, heading := range []string{"DEFECT (2)", "DRIFT (1)", "VERIFICATION (1)", "OTHER (1)"} {
		if strings.Count(column, heading) != 1 {
			t.Errorf("ALL column should head the %q group exactly once:\n%s", heading, column)
		}
	}

	model.Issues.Filter = "defect"
	if model.issuesGroupedByType() {
		t.Fatal("a type filter must not group by type")
	}
	filtered := xansi.Strip(renderIssueColumn("OPEN", data.IssueStatusOpen, model.focusedIssueRows(), 40, 60,
		columnCursor{Holds: true, Focused: true}, model.issuesGroupedByType()))
	if strings.Contains(filtered, "DEFECT (") {
		t.Errorf("DEFECT filter should draw no group heading:\n%s", filtered)
	}
}

// TestEscalatedIssueRowNamesItsObjective proves an escalated Issue's ID line
// carries the Objective it was promoted into, in the Objective purple, and
// that an Issue without escalated_to shows its ID alone.
func TestEscalatedIssueRowNamesItsObjective(t *testing.T) {
	forceColorProfile(t, termenv.TrueColor)

	escalated := &data.IssueV2{ID: "I-060", Title: "Remove the Goal Check", Type: "drift", Status: data.IssueStatusResolved, EscalatedTo: "O-025"}
	row := renderIssueRow(IssueRow{Issue: escalated}, data.IssueStatusResolved, 40, false)
	if first := xansi.Strip(strings.SplitN(row, "\n", 2)[0]); first != "  I-060 → O-025" {
		t.Errorf("escalated ID line = %q, want %q", first, "  I-060 → O-025")
	}
	if !strings.Contains(row, styles.IssueEscalatedObjective.Render("O-025")) {
		t.Errorf("escalation target is not in the Objective purple:\n%q", row)
	}
	if got, want := styles.IssueEscalatedObjective.GetForeground(), styles.SidebarTitleFocused.GetForeground(); got != want {
		t.Errorf("escalation target color = %v, want the sidebar's Objective purple %v", got, want)
	}

	plain := &data.IssueV2{ID: "I-031", Title: "Reload errors", Type: "defect", Status: data.IssueStatusOpen}
	if first := xansi.Strip(strings.SplitN(renderIssueRow(IssueRow{Issue: plain}, data.IssueStatusOpen, 40, false), "\n", 2)[0]); first != "  I-031" {
		t.Errorf("unescalated ID line = %q, want %q", first, "  I-031")
	}
}

func TestIssuesSplitIntoThreeStatusColumns(t *testing.T) {
	root := writeIssuesProject(t)

	got := issueScreen(press(t, issueBoard(t, root), "i"))
	for _, want := range []string{"OPEN (1)", "IN PROGRESS (1)", "RESOLVED (3)"} {
		if !strings.Contains(got, want) {
			t.Errorf("Issues overlay is missing column header %q:\n%s", want, got)
		}
	}

	toInProgress := press(t, issueBoard(t, root), "i", "right", "enter")
	if !strings.Contains(issueScreen(toInProgress), "ID: I-002") {
		t.Errorf("right did not focus the In Progress column:\n%s", issueScreen(toInProgress))
	}

	toResolved := press(t, issueBoard(t, root), "i", "right", "right", "down", "enter")
	if !strings.Contains(issueScreen(toResolved), "ID: I-004") {
		t.Errorf("right right down did not reach the second Resolved row:\n%s", issueScreen(toResolved))
	}

	back := press(t, toResolved, "esc", "left", "left", "enter")
	if !strings.Contains(issueScreen(back), "ID: I-001") {
		t.Errorf("left left did not return focus to the Open column:\n%s", issueScreen(back))
	}
}

func TestTaskScopedIssuesUseIndexedLinksAndRestoreBoardFocus(t *testing.T) {
	root := writeIssuesProject(t)
	model := focusTask(t, issueBoard(t, root), "T-001")
	originalColumn, originalCard := model.FocusedColumn, model.FocusedCard
	model = press(t, model, "I")
	got := issueScreen(model)
	if !strings.Contains(got, "ISSUES · T-001") || !strings.Contains(got, "I-001") || !strings.Contains(got, "I-005") {
		t.Errorf("task-scoped overlay is missing direct linked Issues:\n%s", got)
	}
	if strings.Contains(got, "I-002") {
		t.Errorf("task-scoped overlay included an unlinked Issue:\n%s", got)
	}

	closed := press(t, model, "esc")
	if closed.Issues != nil || closed.FocusedColumn != originalColumn || closed.FocusedCard != originalCard {
		t.Errorf("closing task-scoped Issues did not restore focus: overlay=%v column=%s card=%d", closed.Issues, closed.FocusedColumn, closed.FocusedCard)
	}
}

func TestEmptyIssuesOverlayIsExplicitAndScrollableDetailDoesNotWrap(t *testing.T) {
	empty := press(t, issueBoard(t, writeEmptyProjectFromTemplate(t)), "i")
	if !strings.Contains(issueScreen(empty), "(no issues recorded)") {
		t.Errorf("empty project does not report an explicit Issues state:\n%s", issueScreen(empty))
	}

	model := press(t, openSizedBoard(t, writeIssuesProject(t), 100, 22), "i", "enter")
	for _, line := range strings.Split(issueScreen(model), "\n") {
		if len([]rune(line)) > 120 {
			t.Errorf("Issue detail wrapped past the terminal width: %q", line)
		}
	}
	if model.Issues == nil || model.Issues.Detail == nil {
		t.Fatal("enter did not open Issue detail")
	}
	if scrolled := press(t, model, "down", "down"); scrolled.Issues.DetailOffset != 2 {
		t.Errorf("Issue detail offset = %d after two downs, want 2", scrolled.Issues.DetailOffset)
	}
}

func TestIssuesNavigationAndFilteringAreReadOnly(t *testing.T) {
	root := writeIssuesProject(t)
	before := issueFilesSnapshot(t, root)
	model := press(t, issueBoard(t, root), "i", "f", "down", "enter", "down", "esc", "esc")
	if model.Issues != nil {
		t.Fatal("esc did not close the Issues overlay")
	}
	if after := issueFilesSnapshot(t, root); after != before {
		t.Error("opening, filtering, navigating, and closing Issues changed project files")
	}
}

func TestIssuesSpaceAndBackspaceMoveSelectedIssueAndRetainFocus(t *testing.T) {
	root := writeIssuesProject(t)
	model := press(t, issueBoard(t, root), "i")
	const issueID = "I-001"
	wantStatuses := []data.IssueStatus{
		data.IssueStatusInProgress,
		data.IssueStatusResolved,
		data.IssueStatusInProgress,
		data.IssueStatusOpen,
	}
	keys := []string{" ", " ", "backspace", "backspace"}
	for i, key := range keys {
		var action actionMsg
		model, action = runIssueTransition(t, model, key)
		if action.err != nil {
			t.Fatalf("transition %d (%q) failed: %v", i+1, key, action.err)
		}
		if !action.reload || action.issueSelectedID != issueID {
			t.Fatalf("transition %d action = %#v, want reload and selected Issue %s", i+1, action, issueID)
		}
		if model.Issues == nil || model.Issues.SelectedID != issueID || model.Issues.FocusedStatus != wantStatuses[i] {
			t.Fatalf("transition %d focus = %#v, want %s selected in %s", i+1, model.Issues, issueID, wantStatuses[i])
		}
		rows := model.focusedIssueRows()
		if model.Issues.Cursor < 0 || model.Issues.Cursor >= len(rows) || rows[model.Issues.Cursor].Issue.ID != issueID {
			t.Fatalf("transition %d cursor %d does not select %s in %s", i+1, model.Issues.Cursor, issueID, wantStatuses[i])
		}
	}

	index, err := data.LoadV2Index(root)
	if err != nil {
		t.Fatal(err)
	}
	issue := index.Issues[issueID]
	if issue.Status != data.IssueStatusOpen || issue.Resolution != nil {
		t.Fatalf("Issue after four moves = status %s, resolution %#v; want open with no resolution", issue.Status, issue.Resolution)
	}
	wantKinds := []string{"owner_decision", "owner_decision", "reopened", "owner_decision"}
	if len(issue.History) < len(wantKinds) {
		t.Fatalf("Issue history has %d entries, want at least %d", len(issue.History), len(wantKinds))
	}
	for i, want := range wantKinds {
		entry := issue.History[len(issue.History)-len(wantKinds)+i]
		if string(entry.Kind) != want || entry.Actor.Role != data.ActorRoleOwner || entry.Actor.Session != ownerBoardSession || entry.At.IsZero() {
			t.Errorf("history entry %d = %#v, want %s by board owner at a timestamp", i+1, entry, want)
		}
	}
	if !strings.Contains(issue.History[len(issue.History)-3].Note, "Resolved by the owner from the board.") {
		t.Errorf("resolution history note = %q, want the board reason", issue.History[len(issue.History)-3].Note)
	}
}

func TestIssuesTransitionRefusalsShowInStatusWithoutWriting(t *testing.T) {
	tests := []struct {
		name string
		keys []string
		key  string
		want string
	}{
		{name: "Backspace on Open", keys: []string{"i"}, key: "backspace", want: "cannot retreat Issue I-001"},
		{name: "Space on Resolved", keys: []string{"i", "right", "right"}, key: " ", want: "cannot advance Issue I-003"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := writeIssuesProject(t)
			model := press(t, issueBoard(t, root), tt.keys...)
			before := issueFilesSnapshot(t, root)
			model, action := runIssueTransition(t, model, tt.key)
			if action.err == nil {
				t.Fatal("transition unexpectedly succeeded")
			}
			if !strings.Contains(model.StatusMessage, tt.want) {
				t.Errorf("status line = %q, want %q", model.StatusMessage, tt.want)
			}
			if after := issueFilesSnapshot(t, root); after != before {
				t.Error("refused Issue transition changed project files")
			}
		})
	}
}

func TestIssuesTransitionDoesNothingWithoutASelectionOrInDetail(t *testing.T) {
	root := writeIssuesProject(t)
	empty := press(t, issueBoard(t, root), "i")
	empty.Issues.Filter = "guardrail"
	empty.clampIssueCursorToStatus()
	for _, key := range []string{" ", "backspace"} {
		_, cmd := empty.Update(issueTransitionKeyMsg(key))
		if cmd != nil {
			t.Errorf("empty filtered column scheduled a command for %q", key)
		}
	}

	unselected := press(t, issueBoard(t, root), "i")
	unselected.Issues.SelectedID = ""
	_, cmd := unselected.Update(issueTransitionKeyMsg(" "))
	if cmd != nil {
		t.Error("missing selection scheduled an Issue transition")
	}

	detail := press(t, issueBoard(t, root), "i", "enter")
	for _, key := range []string{" ", "backspace"} {
		_, cmd := detail.Update(issueTransitionKeyMsg(key))
		if cmd != nil {
			t.Errorf("Issue detail scheduled a transition for %q", key)
		}
	}
}

func TestIssuesWriteFailureAppearsInStatusLine(t *testing.T) {
	root := writeIssuesProject(t)
	model := press(t, issueBoard(t, root), "i")
	issuesDir := filepath.Join(root, "issues")
	issuePath := filepath.Join(issuesDir, "I-001-fixture.md")
	dirInfo, err := os.Stat(issuesDir)
	if err != nil {
		t.Fatal(err)
	}
	fileInfo, err := os.Stat(issuePath)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chmod(issuesDir, dirInfo.Mode().Perm())
		_ = os.Chmod(issuePath, fileInfo.Mode().Perm())
	})
	if err := os.Chmod(issuePath, 0o444); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(issuesDir, 0o555); err != nil {
		t.Fatal(err)
	}

	model, action := runIssueTransition(t, model, " ")
	if action.err == nil {
		t.Fatal("transition succeeded while the Issue directory was read-only")
	}
	if !strings.Contains(strings.ToLower(model.StatusMessage), "permission denied") {
		t.Errorf("status line = %q, want the write permission failure", model.StatusMessage)
	}
}

func issueTransitionKeyMsg(key string) tea.KeyMsg {
	if key == "backspace" {
		return tea.KeyMsg{Type: tea.KeyBackspace}
	}
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
}

func runIssueTransition(t *testing.T, model Model, key string) (Model, actionMsg) {
	t.Helper()
	updated, cmd := model.Update(issueTransitionKeyMsg(key))
	current, ok := updated.(Model)
	if !ok {
		t.Fatalf("Update() returned %T, want Model", updated)
	}
	if cmd == nil {
		t.Fatalf("Issue transition key %q returned no command", key)
	}
	message := cmd()
	action, ok := message.(actionMsg)
	if !ok {
		t.Fatalf("transition command returned %T, want actionMsg", message)
	}
	updated, reloadCmd := current.Update(action)
	current, ok = updated.(Model)
	if !ok {
		t.Fatalf("Update(actionMsg) returned %T, want Model", updated)
	}
	if action.err != nil {
		return current, action
	}
	if !action.reload || reloadCmd == nil {
		t.Fatalf("successful transition action = %#v, reload command %v", action, reloadCmd != nil)
	}
	updated, _ = current.Update(reloadCmd())
	current, ok = updated.(Model)
	if !ok {
		t.Fatalf("Update(projectLoadedMsg) returned %T, want Model", updated)
	}
	return current, action
}

func issueFilesSnapshot(t *testing.T, root string) string {
	t.Helper()
	var b strings.Builder
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		fmt.Fprintf(&b, "%s\n%d\n%s\n", path, info.ModTime().UnixNano(), content)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return b.String()
}
