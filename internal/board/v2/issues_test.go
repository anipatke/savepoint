package v2

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/testutil"
)

// writeIssuesProject covers all five descriptive types, all three statuses,
// each resolved disposition, a reopened history, and links through both a
// Task and a Check. Every assertion still runs against a temporary project.
func writeIssuesProject(t *testing.T) string {
	t.Helper()
	root := savepointRoot(t)
	writeConfig(t, root)
	writeRouter(t, root, "task", "O001", "T001")
	writeObjective(t, root, "O001", "Issue surface", "in_progress")
	writeTask(t, root, "O001", "T001", "Task carrying follow-ups", "status: in_progress\nstage: audit\n")
	writeCheck(t, root, "C001", "task", "T001", "CLEAR")

	writeBoardIssue(t, root, "I001", "A repair is needed", "defect", "open", `source: {kind: check, check: C001, actor: {role: checker, session: issue-fixture}, at: '2026-01-01T00:00:00Z'}
tasks: [T001]
checks: [C001]
guardrail_ids: [FS-01]
severity: high
history:
  - at: '2026-01-01T00:00:00Z'
    actor: {role: checker, session: issue-fixture}
    kind: observed
    note: first observation
    check: C001
  - at: '2026-01-02T00:00:00Z'
    actor: {role: executor, session: repair-fixture}
    kind: reopened
    note: the same symptom returned
`)
	writeBoardIssue(t, root, "I002", "A recorded drift", "drift", "in_progress", `source: {kind: report, actor: {role: owner, session: issue-fixture}, at: '2026-01-03T00:00:00Z'}
`)
	writeBoardIssue(t, root, "I003", "An accepted risk", "guardrail", "resolved", `source: {kind: report, actor: {role: owner, session: issue-fixture}, at: '2026-01-04T00:00:00Z'}
resolution: {disposition: accepted, actor: {role: owner, session: owner-fixture}, at: '2026-01-05T00:00:00Z', reason: 'the risk is deliberately accepted'}
`)
	writeBoardIssue(t, root, "I004", "A duplicate report", "verification", "resolved", `source: {kind: report, actor: {role: checker, session: issue-fixture}, at: '2026-01-06T00:00:00Z'}
duplicate_of: I001
resolution: {disposition: duplicate, actor: {role: checker, session: checker-fixture}, at: '2026-01-06T01:00:00Z'}
`)
	writeBoardIssue(t, root, "I005", "A verified repair", "other", "resolved", `source: {kind: report, actor: {role: checker, session: issue-fixture}, at: '2026-01-07T00:00:00Z'}
tasks: [T001]
checks: [C001]
resolution: {disposition: verified, check: C001, actor: {role: checker, session: checker-fixture}, at: '2026-01-07T01:00:00Z'}
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
	for _, id := range []string{"I001", "I002", "I003", "I004", "I005"} {
		at := strings.Index(got, id)
		if at < 0 {
			t.Fatalf("Issues overlay is missing %s:\n%s", id, got)
		}
		if at < previous {
			t.Errorf("Issue %s is out of stable ID order:\n%s", id, got)
		}
		previous = at
	}
	for _, want := range []string{"Filter: ALL", "defect", "open", "severity:high", "A repair is needed"} {
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
		if strings.Contains(got, "Filter: "+want+"\n") && want == "DRIFT" && strings.Contains(got, "I001") {
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
		"ISSUE DETAIL", "ID: I001", "Type: defect", "Status: open",
		"SUMMARY", "The recorded summary for I001.", "Kind: check", "Check: C001",
		"Actor: checker session issue-fixture", "LINKED TASKS", "T001 — Task carrying follow-ups",
		"LINKED CHECKS", "C001 — CLEAR", "GUARDRAILS", "FS-01", "HISTORY",
		"first observation", "the same symptom returned",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("Issue detail is missing %q:\n%s", want, got)
		}
	}

	// I001 is Open; I003 (accepted) is the first of three Resolved rows, so
	// reaching it crosses into the Resolved column rather than scrolling down
	// a single shared list.
	accepted := press(t, model, "esc", "right", "right", "enter")
	acceptedText := issueScreen(accepted)
	if !strings.Contains(acceptedText, "Disposition: accepted") || !strings.Contains(acceptedText, "Not proof of repair") {
		t.Errorf("accepted disposition is not distinct from repair proof:\n%s", acceptedText)
	}

	duplicate := press(t, accepted, "esc", "down", "enter")
	duplicateText := issueScreen(duplicate)
	if !strings.Contains(duplicateText, "Disposition: duplicate") || !strings.Contains(duplicateText, "Canonical: I001") {
		t.Errorf("duplicate detail does not name its canonical target:\n%s", duplicateText)
	}
	toCanonical := press(t, duplicate, "enter")
	if !strings.Contains(issueScreen(toCanonical), "ID: I001") {
		t.Errorf("duplicate navigation did not open the canonical Issue:\n%s", issueScreen(toCanonical))
	}
	verified := press(t, issueBoard(t, root), "i", "right", "right", "down", "down", "enter")
	verifiedText := issueScreen(verified)
	if !strings.Contains(verifiedText, "Disposition: verified") || !strings.Contains(verifiedText, "Proof: Check C001") {
		t.Errorf("verified disposition does not show its proof Check:\n%s", verifiedText)
	}
}

// TestIssuesSplitIntoThreeStatusColumns covers the Open/In Progress/Resolved
// column layout: all three headers with their counts render together, and
// left/right moves focus between columns while up/down stays inside one.
func TestIssuesSplitIntoThreeStatusColumns(t *testing.T) {
	root := writeIssuesProject(t)

	got := issueScreen(press(t, issueBoard(t, root), "i"))
	for _, want := range []string{"OPEN (1)", "IN PROGRESS (1)", "RESOLVED (3)"} {
		if !strings.Contains(got, want) {
			t.Errorf("Issues overlay is missing column header %q:\n%s", want, got)
		}
	}

	toInProgress := press(t, issueBoard(t, root), "i", "right", "enter")
	if !strings.Contains(issueScreen(toInProgress), "ID: I002") {
		t.Errorf("right did not focus the In Progress column:\n%s", issueScreen(toInProgress))
	}

	toResolved := press(t, issueBoard(t, root), "i", "right", "right", "down", "enter")
	if !strings.Contains(issueScreen(toResolved), "ID: I004") {
		t.Errorf("right right down did not reach the second Resolved row:\n%s", issueScreen(toResolved))
	}

	back := press(t, toResolved, "esc", "left", "left", "enter")
	if !strings.Contains(issueScreen(back), "ID: I001") {
		t.Errorf("left left did not return focus to the Open column:\n%s", issueScreen(back))
	}
}

func TestTaskScopedIssuesUseIndexedLinksAndRestoreBoardFocus(t *testing.T) {
	root := writeIssuesProject(t)
	model := focusTask(t, issueBoard(t, root), "T001")
	originalColumn, originalCard := model.FocusedColumn, model.FocusedCard
	model = press(t, model, "I")
	got := issueScreen(model)
	if !strings.Contains(got, "ISSUES · T001") || !strings.Contains(got, "I001") || !strings.Contains(got, "I005") {
		t.Errorf("task-scoped overlay is missing direct linked Issues:\n%s", got)
	}
	if strings.Contains(got, "I002") {
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
