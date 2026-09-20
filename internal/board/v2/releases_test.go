package v2

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	xansi "github.com/charmbracelet/x/ansi"
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/resume"
	"github.com/opencode/savepoint/internal/testutil"
)

func writeReleaseBoardProject(t *testing.T) string {
	t.Helper()
	root := savepointRoot(t)
	writeConfig(t, root)
	writeReleaseRecord(t, root, "R001-first", "R001", "First release")
	writeReleaseRecord(t, root, "R002-second", "R002", "Second release")
	writeRouterWithRelease(t, root, "R001", "O001", "T001")
	writeObjectiveExtra(t, root, "O001", "First objective", "planned", "release: R001\n")
	writeObjectiveExtra(t, root, "O002", "Second objective", "in_progress", "release: R002\n")
	writeTask(t, root, "O001", "T001", "First task", "status: planned\n")
	writeTask(t, root, "O002", "T002", "Second task", "status: in_progress\nstage: build\n")
	return root
}

func writeReleaseRecord(t *testing.T, root, directory, id, title string) {
	t.Helper()
	content := "---\nid: " + id + "\ntitle: \"" + title + "\"\nstatus: planned\n---\n" +
		"# Release\n\n## Outcome\n\nShip the promised delivery.\n\n## Why\n\nThe delivery needs a stable boundary.\n\n## Success Conditions\n\nEvery member Objective is complete.\n\n## Boundaries\n\nRelease does not own Tasks.\n"
	testutil.WriteFile(t, filepath.Join(root, "releases", directory, "Release.md"), content)
}

func writeRouterWithRelease(t *testing.T, root, release, objective, task string) {
	t.Helper()
	content := "# Router\n\n## Current state\n\n```yaml\nstate: task\nrelease: " + release +
		"\nobjective: " + objective + "\ntask: " + task + "\nnext_action: \"do the next thing\"\n```\n"
	testutil.WriteFile(t, filepath.Join(root, "router.md"), content)
}

func releaseBoard(t *testing.T) Model {
	t.Helper()
	return openSizedBoard(t, writeReleaseBoardProject(t), 130, 48)
}

func applyBoardCommands(t *testing.T, model Model, cmd tea.Cmd) Model {
	t.Helper()
	current := model
	for steps := 0; cmd != nil && steps < 3; steps++ {
		updated, next := current.Update(cmd())
		var ok bool
		current, ok = updated.(Model)
		if !ok {
			t.Fatalf("Update() returned %T, want Model", updated)
		}
		cmd = next
	}
	if cmd != nil {
		t.Fatal("board command chain did not settle")
	}
	return current
}

func TestReleaseSelectorOpensOverTheBoardAndStartsOnCurrentRelease(t *testing.T) {
	model := releaseBoard(t)
	before := snapshotProject(t, model.Root)

	updated, cmd := model.Update(keyMsg("r"))
	if cmd != nil {
		t.Fatal("opening the selector returned a command")
	}
	opened := updated.(Model)
	if !opened.ReleaseOverlay {
		t.Fatal("r did not open the Release selector")
	}
	if opened.ReleaseCursor != 0 || opened.SelectedRelease != "R001" {
		t.Errorf("selector state = cursor %d, selected %q; want current R001 at cursor 0", opened.ReleaseCursor, opened.SelectedRelease)
	}
	view := xansi.Strip(opened.View())
	for _, want := range []string{"SELECT RELEASE", "R001", "First release", "S A V E P O I N T"} {
		if !strings.Contains(view, want) {
			t.Errorf("selector view missing %q:\n%s", want, view)
		}
	}

	cancelled, _ := opened.Update(keyMsg("esc"))
	final := cancelled.(Model)
	if final.ReleaseOverlay || final.SelectedRelease != "R001" {
		t.Errorf("Esc changed selector state: open=%t selected=%q", final.ReleaseOverlay, final.SelectedRelease)
	}
	if after := snapshotProject(t, model.Root); !equalSnapshots(before, after) {
		t.Errorf("opening and cancelling changed project files:\n before %v\n  after %v", before, after)
	}
}

func TestReleaseSelectorOpensReadOnlyReleaseDetail(t *testing.T) {
	model := releaseBoard(t)
	before := snapshotProject(t, model.Root)

	detail := press(t, model, "r", "v")
	if detail.ReleaseOverlay || detail.Detail == nil || detail.Detail.Kind != DetailRelease {
		t.Fatalf("release detail state = overlay %t, detail %+v; want a Release detail", detail.ReleaseOverlay, detail.Detail)
	}
	view := xansi.Strip(detail.View())
	for _, want := range []string{
		"RELEASE DETAIL",
		"ID: R001",
		"Title: First release",
		"RELEASE PROMISE",
		"Outcome: Ship the promised delivery.",
		"MEMBER OBJECTIVES",
		"O001 — First objective (planned)",
		"Tasks 0/1 done",
		"RELEASE READINESS",
		"CHECKS",
	} {
		if !strings.Contains(view, want) {
			t.Errorf("Release detail is missing %q:\n%s", want, view)
		}
	}

	closed := press(t, detail, "down", "esc")
	if closed.Detail != nil || closed.ReleaseOverlay {
		t.Fatalf("closing Release detail left an overlay open: detail=%+v release=%t", closed.Detail, closed.ReleaseOverlay)
	}
	if after := snapshotProject(t, model.Root); !equalSnapshots(before, after) {
		t.Errorf("opening, scrolling, and closing Release detail changed project files:\n before %v\n  after %v", before, after)
	}
}

func TestReleaseDetailShowsEvidenceHistoryIssuesAndOwnerBoundary(t *testing.T) {
	root := savepointRoot(t)
	writeConfig(t, root)
	testutil.WriteFile(t, filepath.Join(root, "releases", "R001-first", "Release.md"), `---
id: R001
title: "界 First delivery 🙂"
status: in_progress
last_check: C003
freshness:
  state: current
  check: C003
  assessed_by: {role: checker, session: release-checker}
  assessed_at: '2026-09-14T01:00:00Z'
  basis: "reran the release suite"
owner_validation:
  required: true
  accepted_check: C003
  accepted_by: {role: owner, session: release-owner}
exception:
  requirements: [I001]
  reason: "accepted the documented gap"
  owner: release-owner
  recorded_at: '2026-09-14T02:00:00Z'
  check: C003
---
# Release

## Outcome

Ship the 界 promised delivery 🙂 without losing the source promise.

## Why

The delivery needs a stable boundary.

## Success Conditions

Every member Objective is complete.

## Boundaries

Release does not own Tasks.
`)
	writeObjectiveExtra(t, root, "O001", "Finished Objective", "done", "release: R001\nlast_check: C001\n"+currentFreshness("C001"))
	writeTask(t, root, "O001", "T001", "Finished Task", "status: done\n")
	writeCheck(t, root, "C001", "objective", "O001", "CLEAR")
	writeCheck(t, root, "C002", "release", "R001", "NEEDS WORK")
	writeCheckExtra(t, root, "C003", "release", "R001", "CLEAR", "supersedes: C002\n")
	testutil.WriteFile(t, filepath.Join(root, "issues", "I001.md"), `---
id: I001
title: "Documented release gap"
type: defect
status: open
source: {kind: check, check: C003, actor: {role: checker, session: release-checker}, at: '2026-09-14T00:00:00Z'}
tasks: [T001]
checks: [C002, C003]
---
# Documented release gap
`)
	writeRouterWithRelease(t, root, "R001", "none", "none")

	model := openSizedBoard(t, root, 80, 48)
	detail := press(t, model, "r", "v")
	view := strings.Join(detailLines(*detail.Detail, columnTextWidth(80)), "\n")
	for _, want := range []string{
		"界 First delivery",
		"Outcome: Ship the 界 promised delivery",
		"MEMBER OBJECTIVES",
		"O001 — Finished Objective (done)",
		"Tasks 1/1 done",
		"CLEAR",
		"C002  NEEDS WORK  [superseded]",
		"C003  CLEAR  [latest]",
		"ISSUES",
		"I001 (defect, open): Documented release gap",
		"RELEASE READINESS",
		"Allowed by exception, not by clearance",
		"OWNER VALIDATION",
		"Accepted: Check C003, by owner session release-owner",
	} {
		if !strings.Contains(view, want) {
			t.Errorf("Release detail is missing %q:\n%s", want, view)
		}
	}
	for _, width := range []int{48, 80, 120} {
		sized := press(t, openSizedBoard(t, root, width, 48), "r", "v")
		for lineNo, line := range strings.Split(xansi.Strip(sized.View()), "\n") {
			if got := lipgloss.Width(line); got > width {
				t.Errorf("Release detail line %d is %d cells wide at width %d: %q", lineNo+1, got, width, line)
			}
		}
	}
}

func TestReleaseSelectorNavigationClampsAndQCancels(t *testing.T) {
	model := releaseBoard(t)
	model.SidebarFocused = true
	model.ObjectiveCursor = 0
	opened, _ := model.Update(keyMsg("r"))
	selector := opened.(Model)

	selector = press(t, selector, "j", "j", "k")
	if selector.ReleaseCursor != 0 {
		t.Errorf("ReleaseCursor = %d after j/j/k, want 0", selector.ReleaseCursor)
	}
	selector = press(t, selector, "down")
	if selector.ReleaseCursor != 1 {
		t.Errorf("ReleaseCursor = %d after down, want 1", selector.ReleaseCursor)
	}
	selector = press(t, selector, "down")
	if selector.ReleaseCursor != 1 {
		t.Errorf("ReleaseCursor = %d past the end, want 1", selector.ReleaseCursor)
	}
	closed, _ := selector.Update(keyMsg("q"))
	final := closed.(Model)
	if final.ReleaseOverlay || !final.SidebarFocused || final.ObjectiveCursor != 0 {
		t.Errorf("q did not restore the prior focus: overlay=%t sidebar=%t cursor=%d", final.ReleaseOverlay, final.SidebarFocused, final.ObjectiveCursor)
	}
}

func TestReleaseSelectionFiltersIndexedObjectivesAndPersistsOnlyRouterContext(t *testing.T) {
	root := writeReleaseBoardProject(t)
	model := openSizedBoard(t, root, 130, 48)
	beforeRouter, err := os.ReadFile(filepath.Join(root, "router.md"))
	if err != nil {
		t.Fatal(err)
	}

	opened, _ := model.Update(keyMsg("r"))
	selector := opened.(Model)
	selector = press(t, selector, "j")
	updated, cmd := selector.Update(keyMsg("enter"))
	switched := updated.(Model)
	if cmd == nil {
		t.Fatal("selecting a Release returned no persistence command")
	}
	if switched.ReleaseOverlay || switched.SelectedRelease != "R002" {
		t.Errorf("selection state = open %t, release %q; want closed/R002", switched.ReleaseOverlay, switched.SelectedRelease)
	}
	if got := cardIDsInView(switched); !equalIDs(got, []string{"T002"}) {
		t.Errorf("selected Release shows Tasks %v, want only indexed member T002", got)
	}
	if len(switched.Objectives) != 1 || switched.Objectives[0].ID() != "O002" {
		t.Errorf("selected Release Objectives = %v, want only O002", switched.Objectives)
	}

	final := applyBoardCommands(t, switched, cmd)
	if final.State.Router == nil || final.State.Router.Release != "R002" {
		t.Fatalf("router Release = %+v, want R002 after the canonical write", final.State.Router)
	}
	if final.State.Router.Objective != "" || final.State.Router.Task != "" {
		t.Fatalf("router selection = %+v, want stale Objective/Task cleared for R002", final.State.Router)
	}
	if final.State.Next.SelectionDiagnostic != nil {
		t.Fatalf("reloaded Release selection diagnostic = %+v, want nil", final.State.Next.SelectionDiagnostic)
	}
	if final.State.Next.Release == nil || final.State.Next.Release.ID != "R002" {
		t.Fatalf("reloaded Next.Release = %+v, want R002-scoped Next", final.State.Next.Release)
	}
	if got := cardIDsInView(final); !equalIDs(got, []string{"T002"}) {
		t.Errorf("reloaded selected Release shows Tasks %v, want T002", got)
	}
	if len(final.Objectives) != 1 || final.Objectives[0].ID() != "O002" {
		t.Errorf("reloaded visible Objectives = %v, want only O002", final.Objectives)
	}
	// The board's own Next area and `savepoint resume` intentionally carry
	// different amounts of detail now (the board is a one-line glance; resume
	// stays the full narrative) — so each is checked against its own wording
	// rather than for a shared substring between them.
	view := xansi.Strip(final.View())
	for _, want := range []string{
		"RELEASE: R002 — Second release",
		"Build T002 — Second task",
	} {
		if !strings.Contains(view, want) {
			t.Errorf("board view missing %q:\n%s", want, view)
		}
	}
	var resumeOutput bytes.Buffer
	if err := resume.Render(&resumeOutput, final.State.Next); err != nil {
		t.Fatalf("resume.Render() error = %v", err)
	}
	for _, want := range []string{
		"Task: T002 — Second task",
		"Next action: " + resume.ActionPhrase(final.State.Next),
	} {
		if !strings.Contains(resumeOutput.String(), want) {
			t.Errorf("resume output missing %q:\n%s", want, resumeOutput.String())
		}
	}
	content, err := os.ReadFile(filepath.Join(root, "router.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "release: R002") || !strings.Contains(string(content), "next_action: \"do the next thing\"") {
		t.Errorf("router write lost the Release or next_action:\n%s", content)
	}
	if !strings.Contains(string(content), "objective: none") || !strings.Contains(string(content), "task: none") {
		t.Errorf("Release switch did not clear stale router ownership:\n%s", content)
	}
	if string(beforeRouter) == string(content) {
		t.Error("canonical Release selection did not change router context")
	}
}

func TestSelectionWriteRejectsCrossReleaseObjective(t *testing.T) {
	root := writeReleaseBoardProject(t)
	before, err := os.ReadFile(filepath.Join(root, "router.md"))
	if err != nil {
		t.Fatal(err)
	}

	message, ok := writeSelectionCmd(root, data.RouterSelectionV2{
		Release: "R002", Objective: "O001", Task: "T001",
	})().(actionMsg)
	if !ok || message.err == nil || !strings.Contains(message.err.Error(), "belongs to Release R001, not R002") {
		t.Fatalf("cross-Release selection result = %#v, want a named ownership refusal", message)
	}
	after, err := os.ReadFile(filepath.Join(root, "router.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(before) {
		t.Error("cross-Release selection changed router.md")
	}
}

func TestReleaseSelectionConflictRollsBackAndReloads(t *testing.T) {
	root := writeReleaseBoardProject(t)
	model := openSizedBoard(t, root, 130, 48)
	opened, _ := model.Update(keyMsg("r"))
	selector := opened.(Model)
	selector = press(t, selector, "j")
	updated, cmd := selector.Update(keyMsg("enter"))
	switched := updated.(Model)

	time.Sleep(10 * time.Millisecond)
	writeRouterWithRelease(t, root, "R001", "O001", "T001")
	message := cmd()
	action, ok := message.(actionMsg)
	if !ok || action.err == nil || !strings.Contains(action.err.Error(), "changed on disk") {
		t.Fatalf("conflict message = %#v, want a named mtime refusal", message)
	}
	reloaded, reloadCmd := switched.Update(action)
	rolledBack := reloaded.(Model)
	if rolledBack.SelectedRelease != "R001" {
		t.Errorf("conflict rollback selected %q, want prior R001", rolledBack.SelectedRelease)
	}
	final := applyBoardCommands(t, rolledBack, reloadCmd)
	if final.SelectedRelease != "R001" || final.State.Router.Release != "R001" {
		t.Errorf("post-conflict selection = %q, router = %+v; want truthful R001", final.SelectedRelease, final.State.Router)
	}
	if !strings.Contains(final.StatusMessage, "changed on disk") {
		t.Errorf("post-conflict status = %q, want the refusal reason retained", final.StatusMessage)
	}
}

func TestReleaseSelectionPendingMigrationRefusesWithoutPartialContext(t *testing.T) {
	root := writeReleaseBoardProject(t)
	model := openSizedBoard(t, root, 130, 48)
	operationID := createPendingOperation(t, root)
	opened, _ := model.Update(keyMsg("r"))
	selector := opened.(Model)
	selector = press(t, selector, "j")
	updated, cmd := selector.Update(keyMsg("enter"))
	switched := updated.(Model)
	action, ok := cmd().(actionMsg)
	if !ok || action.err == nil || !strings.Contains(action.err.Error(), operationID) {
		t.Fatalf("pending migration message = %#v, want operation %s named", action, operationID)
	}
	reloaded, reloadCmd := switched.Update(action)
	rolledBack := reloaded.(Model)
	if rolledBack.SelectedRelease != "R001" {
		t.Errorf("pending migration rollback selected %q, want prior R001", rolledBack.SelectedRelease)
	}
	final := applyBoardCommands(t, rolledBack, reloadCmd)
	if !final.State.Migration.Pending || final.SelectedRelease != "" {
		t.Errorf("pending migration state = %+v, selected %q; want migration held and no partial Release", final.State.Migration, final.SelectedRelease)
	}
	if !strings.Contains(final.StatusMessage, operationID) {
		t.Errorf("pending migration status = %q, want refusal reason retained", final.StatusMessage)
	}
}

func TestReleaseReloadPreservesFocusAndDiagnosesRemovedSelection(t *testing.T) {
	root := writeReleaseBoardProject(t)
	model := openSizedBoard(t, root, 130, 48)
	opened, _ := model.Update(keyMsg("r"))
	selector := opened.(Model)
	selector = press(t, selector, "j")
	updated, cmd := selector.Update(keyMsg("enter"))
	model = applyBoardCommands(t, updated.(Model), cmd)
	if model.SelectedRelease != "R002" {
		t.Fatalf("selected Release = %q, want R002", model.SelectedRelease)
	}
	focused := model.focusedTaskID()
	if focused != "T002" {
		t.Fatalf("focused Task = %q, want T002 in R002", focused)
	}

	writeReleaseRecord(t, root, "R002-second", "R002", "Renamed second release")
	reloaded, _ := model.Update(loadCmd(root)().(projectLoadedMsg))
	model = reloaded.(Model)
	if model.SelectedRelease != "R002" || model.focusedTaskID() != focused {
		t.Errorf("reload lost Release/focus: release=%q task=%q", model.SelectedRelease, model.focusedTaskID())
	}

	if err := os.RemoveAll(filepath.Join(root, "releases", "R002-second")); err != nil {
		t.Fatal(err)
	}
	writeObjective(t, root, "O002", "Second objective", "in_progress")
	reloaded, _ = model.Update(loadCmd(root)().(projectLoadedMsg))
	model = reloaded.(Model)
	if model.SelectedRelease != "" {
		t.Errorf("removed Release selection = %q, want no substitute", model.SelectedRelease)
	}
	if model.State.Next.SelectionDiagnostic == nil || model.State.Next.SelectionDiagnostic.Kind != data.SelectionReleaseNotFound {
		t.Errorf("selection diagnostic = %+v, want canonical release_not_found", model.State.Next.SelectionDiagnostic)
	}
}

func TestReleaseSelectorWithoutReleasesIsExplicitAndSafe(t *testing.T) {
	model := openSizedBoard(t, writeEmptyProjectFromTemplate(t), 130, 30)
	opened, cmd := model.Update(keyMsg("r"))
	selector := opened.(Model)
	if cmd != nil || !selector.ReleaseOverlay || !strings.Contains(xansi.Strip(selector.View()), "(none)") {
		t.Errorf("empty Release selector = open %t, cmd %t, view:\n%s", selector.ReleaseOverlay, cmd != nil, xansi.Strip(selector.View()))
	}
	selected, cmd := selector.Update(keyMsg("enter"))
	if cmd != nil || selected.(Model).SelectedRelease != "" {
		t.Errorf("Enter over no Releases changed context or returned a command: %#v", selected)
	}
}

func TestV2HelpAdvertisesTheReleaseSelector(t *testing.T) {
	model := releaseBoard(t)
	opened, _ := model.Update(keyMsg("?"))
	if !strings.Contains(xansi.Strip(opened.(Model).View()), "r: open the Release selector") {
		t.Errorf("help does not advertise r:\n%s", xansi.Strip(opened.(Model).View()))
	}
}

func equalSnapshots(left, right map[string]string) bool {
	if len(left) != len(right) {
		return false
	}
	for path, content := range left {
		if right[path] != content {
			return false
		}
	}
	return true
}
