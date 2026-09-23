package v2

import (
	"crypto/sha256"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/lipgloss"
	xansi "github.com/charmbracelet/x/ansi"
	"github.com/opencode/savepoint/internal/data"
)

// The overlay is where V2's trust boundary becomes inspectable: a badge says
// CLEAR, and this surface has to show which Check said so, who recorded it, and
// on what basis anyone called it current. The assertions below are about what
// the record says, never about what the board concluded.

// openTaskDetail opens the board over root, puts the card cursor on taskID, and
// presses the detail key.
func openTaskDetail(t *testing.T, root, taskID string) Model {
	t.Helper()
	return press(t, focusTask(t, openSizedBoard(t, root, 130, 72), taskID), "enter")
}

// openObjectiveDetail does the same from the sidebar, where v is the detail
// key.
func openObjectiveDetail(t *testing.T, root, objectiveID string) Model {
	t.Helper()
	model := press(t, openSizedBoard(t, root, 130, 72), "left")
	return press(t, focusObjective(t, model, objectiveID), "v")
}

func focusTask(t *testing.T, model Model, taskID string) Model {
	t.Helper()
	for _, column := range columnOrder {
		for i, card := range model.Cards[column] {
			if card.Task.ID == taskID {
				model.FocusedColumn = column
				model.FocusedCard = i
				return model
			}
		}
	}
	t.Fatalf("no card on the board for task %s", taskID)
	return model
}

func focusObjective(t *testing.T, model Model, objectiveID string) Model {
	t.Helper()
	index := model.objectiveIndex(objectiveID)
	if index < 0 {
		t.Fatalf("no sidebar row for objective %s", objectiveID)
	}
	model.ObjectiveCursor = index
	return model
}

// screen is the overlay as a reader sees it, with styling stripped.
func screen(model Model) string {
	return xansi.Strip(model.View())
}

func requireContains(t *testing.T, got string, want ...string) {
	t.Helper()
	for _, part := range want {
		if !strings.Contains(got, part) {
			t.Errorf("rendering is missing %q:\n%s", part, got)
		}
	}
}

func TestTaskDetailNamesTheRecordAndItsLifecycle(t *testing.T) {
	got := screen(openTaskDetail(t, writeEvidenceProject(t), "T-002"))

	requireContains(t, got,
		"TASK DETAIL",
		"ID: T-002",
		"Title: Rechecked after the first run found problems",
		"Status: in_progress",
		"Stage: CHECK",
		"Objective: O-001 — Ship the evidence surface",
		"The board shows the record behind the badge.",
	)
}

// A planned Task has no stage. The row is still rendered, because "no stage
// recorded" is a fact about the record rather than a field the overlay dropped.
func TestTaskDetailReportsAPlannedTaskHasNoStage(t *testing.T) {
	got := screen(openTaskDetail(t, writeEvidenceProject(t), "T-007"))

	requireContains(t, got, "ID: T-007", "Status: planned", "Stage: (none recorded)")
}

func TestObjectiveDetailNamesItsTasksAndItsDependencies(t *testing.T) {
	got := screen(openObjectiveDetail(t, writeEvidenceProject(t), "O-001"))

	requireContains(t, got,
		"OBJECTIVE DETAIL",
		"ID: O-001",
		"Title: Ship the evidence surface",
		"Status: in_progress",
		"OWNED TASKS",
		"T-001 — Cleared and accepted by the owner (done)",
		"T-007 — Nothing recorded against it yet (planned)",
		"OBJECTIVE DEPENDENCIES",
		"O-002 — Groundwork nobody has started (planned)",
		"Objective O-002 is not done yet.",
		"BODY",
		"# Ship the evidence surface",
	)
	if strings.Contains(got, "Stage:") {
		t.Errorf("an Objective detail reports a stage, which Objectives do not carry:\n%s", got)
	}
}

// Dependency state is the resolver's answer. T-001 is done, currently cleared,
// and accepted by its owner, so an accepted-level dependency on it is
// satisfied; T-004 is not done, so a clear-level dependency on it is not.
func TestTaskDetailDependenciesCarryTheirLevelAndResolvedState(t *testing.T) {
	root := writeEvidenceProject(t)

	satisfied := screen(openTaskDetail(t, root, "T-002"))
	requireContains(t, satisfied,
		"DEPENDENCIES",
		"T-001 — Cleared and accepted by the owner (done) — requires accepted",
		"Satisfied.",
	)

	blocked := screen(openTaskDetail(t, root, "T-003"))
	requireContains(t, blocked,
		"T-004 — Flagged for replan (in_progress) — requires clear",
		"Waiting on Task T-004, which is not done yet.",
	)
}

// The chain is the evidence: both runs are listed, in recorded order, with the
// rerun marked as the one that counts and the run it replaced marked as
// replaced. Neither is dropped.
func TestCheckHistoryShowsTheWholeChainInRecordedOrder(t *testing.T) {
	got := screen(openTaskDetail(t, writeEvidenceProject(t), "T-002"))

	requireContains(t, got,
		"C-002  NEEDS WORK  [superseded]",
		"C-003  CLEAR  [latest]",
		"recorded by checker session checker-fixture on 2026-01-02T00:00:00Z",
	)
	// Scoped to the section, because the clearance sentence above it names the
	// latest Check before the chain lists either of them.
	history := got[strings.Index(got, "CHECKS"):]
	if strings.Index(history, "C-002") > strings.Index(history, "C-003") {
		t.Errorf("the Check chain is not in recorded order:\n%s", got)
	}
}

func TestObjectiveCheckHistoryShowsItsOwnIntegrationChain(t *testing.T) {
	got := screen(openObjectiveDetail(t, writeEvidenceProject(t), "O-001"))

	requireContains(t, got, "C-010  CLEAR  [superseded]", "C-011  CLEAR  [latest]")
}

func TestARecordWithNoCheckSaysSoRatherThanShowingNothing(t *testing.T) {
	got := screen(openTaskDetail(t, writeEvidenceProject(t), "T-007"))

	requireContains(t, got, "CHECKS", noCheckRecorded)
}

// Each clearance state gets its own sentence, including the second one
// ResolveClearance reports as unknown: a CLEAR Check no checker session signed.
// That shape is unreachable from a project on disk — the V2 Check decoder
// refuses it — so it is proven here, over the resolved value the board would
// be handed.
func TestDetailClearanceStatesReadDistinctly(t *testing.T) {
	assessed := &data.Freshness{
		State:      data.FreshnessCurrent,
		Check:      "C-001",
		AssessedBy: data.Actor{Role: data.ActorRoleChecker, Session: "checker-fixture"},
		AssessedAt: time.Date(2026, 1, 2, 1, 0, 0, 0, time.UTC),
		Basis:      "reran the suite",
	}
	markedUnknown := *assessed
	markedUnknown.State = data.FreshnessUnknown

	cases := []struct {
		name      string
		clearance data.Clearance
		wantParts []string
	}{
		{"missing", data.Clearance{State: data.ClearanceMissing}, []string{"[ ] Check", "No Check has ever been recorded"}},
		{"needs_work", data.Clearance{State: data.ClearanceNeedsWork, Check: "C-001"}, []string{"Check (needs work)", "C-001 recorded NEEDS WORK"}},
		{"stale", data.Clearance{State: data.ClearanceStale, Check: "C-001", Freshness: assessed}, []string{"Check (stale)", "clearance is stale"}},
		{"unknown", data.Clearance{State: data.ClearanceUnknown, Check: "C-001", Freshness: &markedUnknown}, []string{"Check (unverified)", "clearance is unknown"}},
		{"unknown without checker provenance", data.Clearance{State: data.ClearanceUnknown, Check: "C-001"},
			[]string{"no independent checker session", "not independently established"}},
		{"current", data.Clearance{State: data.ClearanceCurrent, Check: "C-001", Freshness: assessed},
			[]string{"[✓] Check", "C-001 is recorded CLEAR.", "Assessed current by checker session checker-fixture on 2026-01-02", "basis: reran the suite"}},
	}

	rendered := map[string]string{}
	for _, c := range cases {
		got := strings.Join(clearanceLines(DetailTask, c.clearance, false, false), " ")
		requireContains(t, got, c.wantParts...)
		for name, other := range rendered {
			if other == got {
				t.Errorf("%s reads identically to %s: %q", c.name, name, got)
			}
		}
		rendered[c.name] = got
	}
}

// The board never says it verified anything. Every sentence on the overlay
// reports what a record already says.
func TestDetailClaimsNoVerification(t *testing.T) {
	root := writeEvidenceProject(t)
	for _, taskID := range []string{"T-001", "T-002", "T-003", "T-004", "T-005", "T-006", "T-007"} {
		got := strings.ToLower(screen(openTaskDetail(t, root, taskID)))
		for _, claim := range []string{"verified that", "we checked", "confirmed that", "board verified"} {
			if strings.Contains(got, claim) {
				t.Errorf("%s detail claims %q, which the board did not do:\n%s", taskID, claim, got)
			}
		}
	}
}

func TestTaskDetailRendersAReplanByItsRecordedReason(t *testing.T) {
	got := screen(openTaskDetail(t, writeEvidenceProject(t), "T-004"))

	requireContains(t, got,
		"REPLAN",
		"A replan has been flagged: the approach no longer fits",
		"Recorded by planner session planner-fixture on 2026-01-04T00:00:00Z",
	)
}

func TestTaskDetailRendersAnExceptionWithEveryRecordedFact(t *testing.T) {
	got := screen(openTaskDetail(t, writeEvidenceProject(t), "T-005"))

	requireContains(t, got,
		"EXCEPTION",
		"Allowed by exception, not by clearance",
		"shipped with a known gap",
		"Requirements: TEST-02, TEST-05",
		"Applies to: Check C-004",
		"Recorded by owner the owner on 2026-01-05T00:00:00Z",
	)
}

func TestTaskDetailReportsOwnerValidationAndWhatWasAccepted(t *testing.T) {
	root := writeEvidenceProject(t)

	accepted := screen(openTaskDetail(t, root, "T-001"))
	requireContains(t, accepted, "OWNER VALIDATION", "Required: yes", "Accepted: Check C-001, by owner session owner-fixture")

	waiting := screen(openTaskDetail(t, root, "T-006"))
	requireContains(t, waiting, "Required: yes", "Accepted: "+notRecorded)

	notRequired := screen(openTaskDetail(t, root, "T-007"))
	requireContains(t, notRequired, "Required: no")
}

// The Issue link maps are the index's own: I-001 names T-002 as carrying its
// repair and C-002 as the run that observed it, and it is listed once rather
// than once per link.
func TestDetailListsTheIssuesLinkedToTheRecord(t *testing.T) {
	root := writeEvidenceProject(t)

	got := screen(openTaskDetail(t, root, "T-002"))
	requireContains(t, got, "ISSUES", "- I-001 (defect, open): The retry loop drops the last attempt")
	if count := strings.Count(got, "I-001 (defect, open)"); count != 1 {
		t.Errorf("I-001 is listed %d times, want once:\n%s", count, got)
	}

	if unlinked := screen(openTaskDetail(t, root, "T-007")); strings.Contains(unlinked, "ISSUES") {
		t.Errorf("a record with no linked Issue renders an Issues section:\n%s", unlinked)
	}
}

// The overlay scrolls, and both ends hold: a press past either end changes
// nothing at all.
func TestDetailScrollsAndClampsAtBothEnds(t *testing.T) {
	model := press(t, focusTask(t, openSizedBoard(t, writeEvidenceProject(t), 100, 22), "T-002"), "enter")
	if model.Detail == nil {
		t.Fatal("enter did not open the detail overlay")
	}
	top := screen(model)

	scrolled := press(t, model, "down", "down")
	if scrolled.DetailOffset != 2 {
		t.Fatalf("DetailOffset = %d after two downs, want 2", scrolled.DetailOffset)
	}
	if screen(scrolled) == top {
		t.Error("scrolling the overlay changed nothing on screen")
	}

	atTop := press(t, model, "up", "up")
	if atTop.DetailOffset != 0 || screen(atTop) != top {
		t.Errorf("DetailOffset = %d after scrolling up at the top, want 0 and an unchanged screen", atTop.DetailOffset)
	}

	atEnd := press(t, model, strings.Split(strings.Repeat("down ", 60), " ")[:60]...)
	if press(t, atEnd, "down").DetailOffset != atEnd.DetailOffset {
		t.Error("scrolling down past the end of the overlay moved it")
	}
	requireContains(t, screen(atEnd), "above")
}

// Closing returns the keys to the surface the overlay was opened from, with the
// cursor where the reader left it — including a sidebar cursor deliberately
// moved off the selected Objective.
func TestClosingTheDetailRestoresTheSurfaceItWasOpenedFrom(t *testing.T) {
	root := writeEvidenceProject(t)

	fromColumns := focusTask(t, openSizedBoard(t, root, 130, 72), "T-005")
	reopened := press(t, fromColumns, "enter", "down", "esc")
	if reopened.Detail != nil {
		t.Fatal("esc left the overlay open")
	}
	if reopened.SidebarFocused || reopened.FocusedColumn != fromColumns.FocusedColumn || reopened.FocusedCard != fromColumns.FocusedCard {
		t.Errorf("closing left focus at %q card %d (sidebar %v), want %q card %d on the columns",
			reopened.FocusedColumn, reopened.FocusedCard, reopened.SidebarFocused,
			fromColumns.FocusedColumn, fromColumns.FocusedCard)
	}
	if screen(reopened) != screen(fromColumns) {
		t.Error("closing the overlay did not return the board to the surface it was opened over")
	}

	fromSidebar := focusObjective(t, press(t, openSizedBoard(t, root, 130, 72), "left"), "O-002")
	closed := press(t, fromSidebar, "v", "esc")
	if !closed.SidebarFocused || closed.ObjectiveCursor != fromSidebar.ObjectiveCursor {
		t.Errorf("closing left the sidebar cursor at %d (focused %v), want %d on the sidebar",
			closed.ObjectiveCursor, closed.SidebarFocused, fromSidebar.ObjectiveCursor)
	}
	if closed.SelectedObjective != fromSidebar.SelectedObjective {
		t.Errorf("opening a detail changed the selection to %q, want %q unchanged",
			closed.SelectedObjective, fromSidebar.SelectedObjective)
	}
}

// While an overlay is open the keys belong to it: esc closes it rather than
// clearing the sidebar's selection, and a surface-crossing arrow key does not
// move focus behind it.
func TestTheOpenOverlayHoldsTheKeys(t *testing.T) {
	model := press(t, openSizedBoard(t, writeEvidenceProject(t), 120, 44), "left", "v")
	if model.Detail == nil {
		t.Fatal("v did not open the Objective detail")
	}
	selected := model.SelectedObjective

	crossed := press(t, model, "right")
	if crossed.Detail == nil || crossed.SidebarFocused != model.SidebarFocused {
		t.Error("right moved focus behind the open overlay")
	}

	closed := press(t, model, "esc")
	if closed.Detail != nil {
		t.Error("esc did not close the overlay")
	}
	if closed.SelectedObjective != selected {
		t.Errorf("esc cleared the selection to %q while closing the overlay, want %q", closed.SelectedObjective, selected)
	}
}

// The body is displayed, not read. A body full of things that look like
// frontmatter, checkboxes, and headings reaches the screen verbatim and changes
// nothing the overlay says about the record.
func TestTheRecordBodyIsDisplayedAndNotParsed(t *testing.T) {
	got := screen(openTaskDetail(t, writeEvidenceProject(t), "T-002"))

	requireContains(t, got, "BODY", "# Outcome", "- [x] status: done", "- [ ] still open")
	// The body says "status: done"; the record says in_progress, and the record wins.
	requireContains(t, got, "Status: in_progress")
}

// Long and wide content costs extra lines inside the overlay and never a wider
// board: the frame is measured, not pushed.
func TestLongAndWideBodyContentDoesNotWidenTheFrame(t *testing.T) {
	root := savepointRoot(t)
	writeConfig(t, root)
	writeRouter(t, root, "task", "O-001", "T-001")
	writeObjective(t, root, "O-001", "Wide content", "in_progress")
	writeTaskBody(t, root, "O-001", "T-001", "Carries a body nothing can wrap", "status: planned\n",
		strings.Repeat("x", 500)+"\n\n"+strings.Repeat("日本語のテキスト", 40)+"\n")

	for _, width := range []int{80, 100, 120, 160} {
		got := screen(press(t, openSizedBoard(t, root, width, 30), "enter"))
		for _, line := range strings.Split(got, "\n") {
			if lipgloss.Width(line) > width {
				t.Fatalf("at width %d the overlay drew a %d-cell line: %q", width, lipgloss.Width(line), line)
			}
		}
	}
}

// Opening, scrolling, and closing a detail is a read. Nothing under the project
// root changes: not its bytes, not its modification times.
func TestOpeningScrollingAndClosingADetailWritesNothing(t *testing.T) {
	root := writeEvidenceProject(t)
	before := snapshotTree(t, root)

	model := focusTask(t, openSizedBoard(t, root, 100, 20), "T-002")
	model = press(t, model, "enter", "down", "down", "up", "esc")
	model.SidebarFocused = true
	model = press(t, model, "v", "down", "esc")

	if model.Detail != nil {
		t.Fatal("the key sequence left an overlay open")
	}
	if after := snapshotTree(t, root); after != before {
		t.Errorf("opening a detail changed the project:\nbefore:\n%s\nafter:\n%s", before, after)
	}
}

// snapshotTree records every file under root by path, size, content hash, and
// modification time, so a comparison catches a rewrite that preserved the bytes
// as well as one that changed them.
func snapshotTree(t *testing.T, root string) string {
	t.Helper()
	var out strings.Builder
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		fmt.Fprintf(&out, "%s %d %x %s\n", path, len(content), sha256.Sum256(content), info.ModTime().UTC().Format(time.RFC3339Nano))
		return nil
	})
	if err != nil {
		t.Fatalf("snapshot %s: %v", root, err)
	}
	return out.String()
}

// TestDetailRenderingReadsOnlyTheResolvedValue is the derive-nothing rule as
// behavior rather than as a grep: with the loaded index taken away entirely,
// an open overlay still renders exactly what it rendered before. Nothing on it
// was being looked up at render time.
func TestDetailRenderingReadsOnlyTheResolvedValue(t *testing.T) {
	model := openTaskDetail(t, writeEvidenceProject(t), "T-002")
	// The board's body alone: the chrome around it does report the loaded
	// counts, and removing the index legitimately changes those.
	before := xansi.Strip(model.renderBody(120, 48))

	model.State.Index = nil
	if after := xansi.Strip(model.renderBody(120, 48)); after != before {
		t.Errorf("the overlay changed once the index was removed, so it was resolving at render time:\nbefore:\n%s\nafter:\n%s", before, after)
	}
}

// A reload under an open overlay re-resolves it, and a record the reload
// removed closes it rather than leaving a copy older than the project.
func TestReloadRefreshesAnOpenDetailAndClosesADeletedOne(t *testing.T) {
	root := writeEvidenceProject(t)
	model := openTaskDetail(t, root, "T-002")
	requireContains(t, screen(model), "C-003  CLEAR  [latest]")

	// A third run against T-002, recorded while its overlay is open.
	writeCheckExtra(t, root, "C-006", "task", "T-002", "NEEDS WORK", "supersedes: C-003\n")
	reloaded, _ := model.Update(loadCmd(root)().(projectLoadedMsg))
	refreshed := reloaded.(Model)
	if refreshed.Detail == nil {
		t.Fatal("a reload closed an overlay whose record still exists")
	}
	requireContains(t, screen(refreshed),
		"C-003  CLEAR  [superseded]",
		"C-006  NEEDS WORK  [latest]",
		"Check C-006 recorded NEEDS WORK.",
	)

	// A record the reload no longer has: the overlay closes rather than keeping
	// a copy older than the project.
	open := openTaskDetail(t, root, "T-007")
	removeTask(t, root, "O-001", "T-007")
	closed, _ := open.Update(loadCmd(root)().(projectLoadedMsg))
	if closed.(Model).Detail != nil {
		t.Error("a reload that removed the open record left its overlay open")
	}
}

// There is nothing to open on an empty project, and pressing the detail key
// there changes nothing.
func TestTheDetailKeyOverAnEmptyProjectOpensNothing(t *testing.T) {
	model := openSizedBoard(t, writeEmptyProjectFromTemplate(t), 100, 30)

	pressed := press(t, model, "enter", "left", "v")
	if pressed.Detail != nil {
		t.Error("the detail key opened an overlay over a project with no records")
	}
}
