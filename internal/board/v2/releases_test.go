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
	writeReleaseRecord(t, root, "R-001-first", "R-001", "First release")
	writeReleaseRecord(t, root, "R-002-second", "R-002", "Second release")
	writeRouterWithRelease(t, root, "R-001", "O-001", "T-001")
	writeObjectiveExtra(t, root, "O-001", "First objective", "planned", "release: R-001\n")
	writeObjectiveExtra(t, root, "O-002", "Second objective", "in_progress", "release: R-002\n")
	writeTask(t, root, "O-001", "T-001", "First task", "status: planned\n")
	writeTask(t, root, "O-002", "T-002", "Second task", "status: in_progress\nstage: build\n")
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

	updated, cmd := model.Update(keyMsg(goalSelectorKey))
	if cmd != nil {
		t.Fatal("opening the selector returned a command")
	}
	opened := updated.(Model)
	if !opened.ReleaseOverlay {
		t.Fatal("g did not open the Goal selector")
	}
	if opened.ReleaseCursor != 0 || opened.SelectedRelease != "R-001" {
		t.Errorf("selector state = cursor %d, selected %q; want current R-001 at cursor 0", opened.ReleaseCursor, opened.SelectedRelease)
	}
	view := xansi.Strip(opened.View())
	for _, want := range []string{"SELECT GOAL", "R-001", "First release", "S A V E P O I N T"} {
		if !strings.Contains(view, want) {
			t.Errorf("selector view missing %q:\n%s", want, view)
		}
	}

	cancelled, _ := opened.Update(keyMsg("esc"))
	final := cancelled.(Model)
	if final.ReleaseOverlay || final.SelectedRelease != "R-001" {
		t.Errorf("Esc changed selector state: open=%t selected=%q", final.ReleaseOverlay, final.SelectedRelease)
	}
	if after := snapshotProject(t, model.Root); !equalSnapshots(before, after) {
		t.Errorf("opening and cancelling changed project files:\n before %v\n  after %v", before, after)
	}
}

func TestGoalSelectorCannotClearTheCurrentGoal(t *testing.T) {
	model := releaseBoard(t)
	opened, cmd := model.Update(keyMsg(goalSelectorKey))
	selector := opened.(Model)
	if cmd != nil || !selector.ReleaseOverlay || selector.SelectedRelease != "R-001" {
		t.Fatalf("selector state = open %t, Goal %q, command %t; want open/R-001/no command", selector.ReleaseOverlay, selector.SelectedRelease, cmd != nil)
	}
	view := strings.ToLower(xansi.Strip(selector.View()))
	if strings.Contains(view, "clear goal") || strings.Contains(view, "(none)") || strings.Contains(view, "(no goal)") {
		t.Errorf("Goal selector offers an empty selection:\n%s", view)
	}

	for _, key := range []string{"backspace", "delete"} {
		updated, cmd := selector.Update(keyMsg(key))
		selector = updated.(Model)
		if cmd != nil || !selector.ReleaseOverlay || selector.SelectedRelease != "R-001" || selector.State.Router.Release != "R-001" {
			t.Errorf("%q cleared or changed the selected Goal: overlay=%t view=%q router=%+v command=%t", key, selector.ReleaseOverlay, selector.SelectedRelease, selector.State.Router, cmd != nil)
		}
	}
}

func TestLegacyReleaseKeyRemainsAnUndisclosedGoalSelectorAlias(t *testing.T) {
	model := releaseBoard(t)
	opened, cmd := model.Update(keyMsg(goalSelectorAlias))
	selector := opened.(Model)
	if cmd != nil || !selector.ReleaseOverlay {
		t.Fatalf("legacy r alias opened %t and returned command %t; want the Goal selector with no command", selector.ReleaseOverlay, cmd != nil)
	}
	selectorView := xansi.Strip(selector.View())
	if !strings.Contains(selectorView, "SELECT GOAL") {
		t.Errorf("legacy alias opened the wrong selector:\n%s", selectorView)
	}
	if strings.Contains(selectorView, "r:open") || strings.Contains(selectorView, "r:releases") {
		t.Errorf("legacy alias is exposed in selector UI:\n%s", selectorView)
	}
}

func TestGoalSelectorFitsNarrowBoardWidths(t *testing.T) {
	root := writeReleaseBoardProject(t)
	for _, width := range []int{32, 48} {
		model := openSizedBoard(t, root, width, 40)
		opened, _ := model.Update(keyMsg(goalSelectorKey))
		selector := opened.(Model)
		view := xansi.Strip(selector.View())
		if !selector.ReleaseOverlay || !strings.Contains(view, "SELECT GOAL") {
			t.Errorf("width %d did not render the Goal selector:\n%s", width, view)
		}
		for lineNo, line := range strings.Split(view, "\n") {
			if got := lipgloss.Width(line); got > width {
				t.Errorf("selector line %d is %d cells wide at terminal width %d: %q", lineNo+1, got, width, line)
			}
		}
	}
}

func TestReleaseSelectorOpensReadOnlyReleaseDetail(t *testing.T) {
	model := releaseBoard(t)
	before := snapshotProject(t, model.Root)

	detail := press(t, model, goalSelectorKey, "v")
	if detail.ReleaseOverlay || detail.Detail == nil || detail.Detail.Kind != DetailRelease {
		t.Fatalf("release detail state = overlay %t, detail %+v; want a Release detail", detail.ReleaseOverlay, detail.Detail)
	}
	view := xansi.Strip(detail.View())
	for _, want := range []string{
		"GOAL DETAIL",
		"ID: R-001",
		"Title: First release",
		"GOAL PROMISE",
		"Outcome: Ship the promised delivery.",
		"MEMBER OBJECTIVES",
		"O-001 — First objective (planned)",
		"Tasks 0/1 done",
		"GOAL READINESS",
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
	testutil.WriteFile(t, filepath.Join(root, "releases", "R-001-first", "Release.md"), `---
id: R-001
title: "界 First delivery 🙂"
status: in_progress
last_check: C-003
freshness:
  state: current
  check: C-003
  assessed_by: {role: checker, session: release-checker}
  assessed_at: '2026-09-14T01:00:00Z'
  basis: "reran the release suite"
owner_validation:
  required: true
  accepted_check: C-003
  accepted_by: {role: owner, session: release-owner}
exception:
  requirements: [I-001]
  reason: "accepted the documented gap"
  owner: release-owner
  recorded_at: '2026-09-14T02:00:00Z'
  check: C-003
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
	writeObjectiveExtra(t, root, "O-001", "Finished Objective", "done", "release: R-001\nlast_check: C-001\n"+currentFreshness("C-001"))
	writeTask(t, root, "O-001", "T-001", "Finished Task", "status: done\n")
	writeCheck(t, root, "C-001", "objective", "O-001", "CLEAR")
	writeCheck(t, root, "C-002", "release", "R-001", "NEEDS WORK")
	writeCheckExtra(t, root, "C-003", "release", "R-001", "CLEAR", "supersedes: C-002\n")
	testutil.WriteFile(t, filepath.Join(root, "issues", "I-001.md"), `---
id: I-001
title: "Documented release gap"
type: defect
status: open
source: {kind: check, check: C-003, actor: {role: checker, session: release-checker}, at: '2026-09-14T00:00:00Z'}
tasks: [T-001]
checks: [C-002, C-003]
---
# Documented release gap
`)
	writeRouterWithRelease(t, root, "R-001", "none", "none")

	model := openSizedBoard(t, root, 80, 48)
	detail := press(t, model, goalSelectorKey, "v")
	view := strings.Join(detailLines(*detail.Detail, columnTextWidth(80)), "\n")
	for _, want := range []string{
		"界 First delivery",
		"Outcome: Ship the 界 promised delivery",
		"MEMBER OBJECTIVES",
		"O-001 — Finished Objective (done)",
		"Tasks 1/1 done",
		"CLEAR",
		"C-002  NEEDS WORK  [superseded]",
		"C-003  CLEAR  [latest]",
		"ISSUES",
		"I-001 (defect, open): Documented release gap",
		"GOAL READINESS",
		"Allowed by exception, not by clearance",
		"OWNER VALIDATION",
		"Accepted: Check C-003, by owner session release-owner",
	} {
		if !strings.Contains(view, want) {
			t.Errorf("Release detail is missing %q:\n%s", want, view)
		}
	}
	for _, width := range []int{48, 80, 120} {
		sized := press(t, openSizedBoard(t, root, width, 48), goalSelectorKey, "v")
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
	opened, _ := model.Update(keyMsg(goalSelectorKey))
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

	opened, _ := model.Update(keyMsg(goalSelectorKey))
	selector := opened.(Model)
	selector = press(t, selector, "j")
	updated, cmd := selector.Update(keyMsg("enter"))
	switched := updated.(Model)
	if cmd == nil {
		t.Fatal("selecting a Release returned no persistence command")
	}
	if switched.ReleaseOverlay || switched.SelectedRelease != "R-002" {
		t.Errorf("selection state = open %t, release %q; want closed/R-002", switched.ReleaseOverlay, switched.SelectedRelease)
	}
	if got := cardIDsInView(switched); !equalIDs(got, []string{"T-002"}) {
		t.Errorf("selected Release shows Tasks %v, want only indexed member T-002", got)
	}
	if len(switched.Objectives) != 1 || switched.Objectives[0].ID() != "O-002" {
		t.Errorf("selected Release Objectives = %v, want only O-002", switched.Objectives)
	}

	final := applyBoardCommands(t, switched, cmd)
	if final.State.Router == nil || final.State.Router.Release != "R-002" {
		t.Fatalf("router Release = %+v, want R-002 after the canonical write", final.State.Router)
	}
	if final.State.Router.Objective != "" || final.State.Router.Task != "" {
		t.Fatalf("router selection = %+v, want stale Objective/Task cleared for R-002", final.State.Router)
	}
	if final.State.Next.SelectionDiagnostic != nil {
		t.Fatalf("reloaded Release selection diagnostic = %+v, want nil", final.State.Next.SelectionDiagnostic)
	}
	if final.State.Next.Release == nil || final.State.Next.Release.ID != "R-002" {
		t.Fatalf("reloaded Next.Release = %+v, want R-002-scoped Next", final.State.Next.Release)
	}
	if got := cardIDsInView(final); !equalIDs(got, []string{"T-002"}) {
		t.Errorf("reloaded selected Release shows Tasks %v, want T-002", got)
	}
	if len(final.Objectives) != 1 || final.Objectives[0].ID() != "O-002" {
		t.Errorf("reloaded visible Objectives = %v, want only O-002", final.Objectives)
	}
	// The board's own Next area and `savepoint resume` intentionally carry
	// different amounts of detail now (the board is a one-line glance; resume
	// stays the full narrative) — so each is checked against its own wording
	// rather than for a shared substring between them.
	view := xansi.Strip(final.View())
	for _, want := range []string{
		"GOAL: R-002 — Second release",
		"ALL OBJECTIVES IN GOAL R-002",
		"NEXT: Nothing selected",
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
		"Nothing selected",
		"Next action: " + resume.ActionPhrase(final.State.Next),
	} {
		if !strings.Contains(resumeOutput.String(), want) {
			t.Errorf("resume output missing %q:\n%s", want, resumeOutput.String())
		}
	}
	if strings.Contains(resumeOutput.String(), "Task: T-002 — Second task") {
		t.Errorf("resume output substituted the Release's member Task without an Objective selection:\n%s", resumeOutput.String())
	}
	content, err := os.ReadFile(filepath.Join(root, "router.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "release: R-002") || !strings.Contains(string(content), "next_action: \"do the next thing\"") {
		t.Errorf("router write lost the Release or next_action:\n%s", content)
	}
	if !strings.Contains(string(content), "objective: none") || !strings.Contains(string(content), "task: none") {
		t.Errorf("Release switch did not clear stale router ownership:\n%s", content)
	}
	if string(beforeRouter) == string(content) {
		t.Error("canonical Release selection did not change router context")
	}

	var plain bytes.Buffer
	if err := Run(Options{Root: root, Stdout: &plain, TTY: false}); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if !strings.Contains(plain.String(), "Selected: all Objectives in Goal R-002\n") {
		t.Errorf("plain output does not scope the unfiltered view to R-002:\n%s", plain.String())
	}
	if strings.Contains(plain.String(), "Selected: all Objectives\n") {
		t.Errorf("plain output retains the unscoped selection label:\n%s", plain.String())
	}
}

func TestGoalScopedSelectionWithObjectiveFilter(t *testing.T) {
	root := writeReleaseBoardProject(t)
	writeRouterWithRelease(t, root, "R-002", "none", "none")

	unfiltered := openSizedBoard(t, root, 130, 48)
	if unfiltered.SelectedRelease != "R-002" || unfiltered.SelectedObjective != "" {
		t.Fatalf("unfiltered selection = Goal %q, Objective %q; want R-002 and none", unfiltered.SelectedRelease, unfiltered.SelectedObjective)
	}
	if got := cardIDsInView(unfiltered); !equalIDs(got, []string{"T-002"}) {
		t.Errorf("unfiltered Goal view shows Tasks %v, want only R-002's T-002", got)
	}
	if !strings.Contains(xansi.Strip(unfiltered.View()), "ALL OBJECTIVES IN GOAL R-002") {
		t.Errorf("unfiltered TUI view does not name its Goal scope:\n%s", xansi.Strip(unfiltered.View()))
	}

	filtered := openBoard(t, root, "O-002")
	if filtered.SelectedRelease != "R-002" || filtered.SelectedObjective != "O-002" {
		t.Fatalf("filtered selection = Goal %q, Objective %q; want R-002/O-002", filtered.SelectedRelease, filtered.SelectedObjective)
	}
	if got := cardIDsInView(filtered); !equalIDs(got, []string{"T-002"}) {
		t.Errorf("filtered Goal view shows Tasks %v, want only T-002", got)
	}
	if strings.Contains(xansi.Strip(filtered.View()), allObjectivesLabel) {
		t.Errorf("filtered TUI view carries the unfiltered Goal label:\n%s", xansi.Strip(filtered.View()))
	}

	var plain bytes.Buffer
	if err := Run(Options{Root: root, ObjectiveFilter: "O-002", Stdout: &plain, TTY: false}); err != nil {
		t.Fatalf("Run() with Objective filter error = %v", err)
	}
	if !strings.Contains(plain.String(), "Selected: O-002\n") || strings.Contains(plain.String(), "Selected: all Objectives") {
		t.Errorf("filtered plain output does not retain the Objective filter:\n%s", plain.String())
	}
}

func TestBoardWithoutAValidGoalShowsNoProjectWideRecords(t *testing.T) {
	root := writeReleaseBoardProject(t)
	writeRouterWithRelease(t, root, "", "none", "none")

	model := openSizedBoard(t, root, 130, 48)
	if model.State.Next.SelectionDiagnostic == nil || model.State.Next.SelectionDiagnostic.Kind != data.SelectionReleaseMissing {
		t.Fatalf("Next selection diagnostic = %+v, want the canonical missing-Goal diagnostic", model.State.Next.SelectionDiagnostic)
	}
	if len(taskIDsInReleaseView(model.State.Index, "", "")) != 0 || len(cardIDsInView(model)) != 0 || len(model.Objectives) != 0 {
		t.Errorf("empty-Goal projection shows Tasks %v and Objectives %v", cardIDsInView(model), model.Objectives)
	}
	if model.SelectedRelease != "" || model.SelectedObjective != "" {
		t.Errorf("empty-Goal selection = %q/%q, want no Goal or Objective", model.SelectedRelease, model.SelectedObjective)
	}
	view := xansi.Strip(model.View())
	for _, want := range []string{"Choose a Goal", resume.SelectionPhrase(model.State.Next.SelectionDiagnostic), "g:goals", "(empty)"} {
		if !strings.Contains(view, want) {
			t.Errorf("no-Goal TUI is missing %q:\n%s", want, view)
		}
	}
	for _, forbidden := range []string{"ALL OBJECTIVES", "O-001", "O-002", "T-001", "T-002", "First objective", "Second objective"} {
		if strings.Contains(view, forbidden) {
			t.Errorf("no-Goal TUI shows project-wide record %q:\n%s", forbidden, view)
		}
	}

	filtered := openBoard(t, root, "O-001")
	if filtered.SelectedObjective != "" || len(cardIDsInView(filtered)) != 0 {
		t.Errorf("Objective filter escaped the missing-Goal scope: selected %q, Tasks %v", filtered.SelectedObjective, cardIDsInView(filtered))
	}

	var plain bytes.Buffer
	if err := Run(Options{Root: root, Stdout: &plain, TTY: false}); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	for _, want := range []string{"Choose a Goal", resume.SelectionPhrase(model.State.Next.SelectionDiagnostic), "Selected: no Goal selected"} {
		if !strings.Contains(plain.String(), want) {
			t.Errorf("no-Goal plain output is missing %q:\n%s", want, plain.String())
		}
	}
	for _, forbidden := range []string{"Selected: all Objectives", "T-001 — First task", "T-002 — Second task"} {
		if strings.Contains(plain.String(), forbidden) {
			t.Errorf("no-Goal plain output shows project-wide content %q:\n%s", forbidden, plain.String())
		}
	}
}

func TestBoardWithZeroGoalsPointsToDoctor(t *testing.T) {
	root := writeEmptyProjectFromTemplate(t)
	writeRouterWithRelease(t, root, "R-001", "none", "none")
	if err := os.RemoveAll(filepath.Join(root, "releases")); err != nil {
		t.Fatal(err)
	}

	assertBoardPointsToDoctor(t, root)
}

func TestBoardWithOnlyArchivedGoalsPointsToDoctor(t *testing.T) {
	root := writeEmptyProjectFromTemplate(t)
	writeRouterWithRelease(t, root, "R-001", "none", "none")
	testutil.WriteFile(t, filepath.Join(root, "releases", "R-001-history", "Release.md"),
		"---\nid: R-001\ntitle: \"Historical Goal\"\nstatus: done\n"+
			"legacy_completion:\n  source_path: releases/v1/PRD.md\n  archive_path: archive/v1/PRD.md\n"+
			"  sha256: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\n---\n\n"+
			"## Outcome\n\nPreserve historical completion.\n\n"+
			"## Why\n\nThe original scope is archived.\n\n"+
			"## Success Conditions\n\nThe original work was completed.\n\n"+
			"## Boundaries\n\nThis Goal is archived.\n")

	assertBoardPointsToDoctor(t, root)
}

func assertBoardPointsToDoctor(t *testing.T, root string) {
	t.Helper()
	model := openSizedBoard(t, root, 130, 48)
	view := xansi.Strip(model.View())
	if !strings.Contains(view, "savepoint doctor") || strings.Contains(view, "Choose a Goal") {
		t.Errorf("project with no live Goals does not direct the owner to doctor:\n%s", view)
	}

	var plain bytes.Buffer
	if err := Run(Options{Root: root, Stdout: &plain, TTY: false}); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if !strings.Contains(plain.String(), "savepoint doctor") || strings.Contains(plain.String(), "Choose a Goal") {
		t.Errorf("plain output with no live Goals does not direct the owner to doctor:\n%s", plain.String())
	}
}

func TestUnassignedGoalNoticeUsesIndexFactInBothRenderers(t *testing.T) {
	root := writeReleaseBoardProject(t)
	writeObjectiveExtra(t, root, "O-003", "Unassigned objective", "planned", "")
	model := openSizedBoard(t, root, 130, 48)
	if got := len(model.State.Index.ObjectivesWithoutGoal); got != 1 {
		t.Fatalf("T-036 missing-Goal index fact has %d IDs, want 1", got)
	}
	want := "1 Objective has no Goal; run savepoint doctor."
	if !strings.Contains(xansi.Strip(model.View()), want) {
		t.Errorf("TUI omits the missing-Goal count and repair pointer:\n%s", xansi.Strip(model.View()))
	}
	var plain bytes.Buffer
	if err := Run(Options{Root: root, Stdout: &plain, TTY: false}); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if !strings.Contains(plain.String(), want) {
		t.Errorf("plain output omits the missing-Goal count and repair pointer:\n%s", plain.String())
	}

	// The display count follows the index fact even if the Objective map would
	// imply a different count, proving the board does not recompute membership.
	model.State.Index.ObjectivesWithoutGoal = []string{"O-003", "O-004"}
	want = "2 Objectives have no Goal; run savepoint doctor."
	if !strings.Contains(xansi.Strip(model.View()), want) || !strings.Contains(renderPlain(model.State, model.SelectedObjective), want) {
		t.Errorf("renderers recomputed instead of using ObjectivesWithoutGoal; TUI:\n%s\nplain:\n%s", xansi.Strip(model.View()), renderPlain(model.State, model.SelectedObjective))
	}
}

func TestGoalScopedChromeFitsNarrowWidths(t *testing.T) {
	root := writeReleaseBoardProject(t)
	model := openSizedBoard(t, root, 130, 48)
	for _, width := range []int{64, 80, 112} {
		above, _ := model.boardChrome(width)
		for _, section := range above {
			for lineNo, line := range strings.Split(xansi.Strip(section), "\n") {
				if got := lipgloss.Width(line); got > width {
					t.Errorf("board chrome line %d is %d cells at width %d: %q", lineNo+1, got, width, line)
				}
			}
		}
	}
}

func TestSelectionWriteRejectsCrossReleaseObjective(t *testing.T) {
	root := writeReleaseBoardProject(t)
	before, err := os.ReadFile(filepath.Join(root, "router.md"))
	if err != nil {
		t.Fatal(err)
	}

	message, ok := writeSelectionCmd(root, data.RouterSelectionV2{
		Release: "R-002", Objective: "O-001", Task: "T-001",
	})().(actionMsg)
	if !ok || message.err == nil || !strings.Contains(message.err.Error(), "belongs to Goal R-001, not R-002") {
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
	opened, _ := model.Update(keyMsg(goalSelectorKey))
	selector := opened.(Model)
	selector = press(t, selector, "j")
	updated, cmd := selector.Update(keyMsg("enter"))
	switched := updated.(Model)

	time.Sleep(10 * time.Millisecond)
	writeRouterWithRelease(t, root, "R-001", "O-001", "T-001")
	message := cmd()
	action, ok := message.(actionMsg)
	if !ok || action.err == nil || !strings.Contains(action.err.Error(), "changed on disk") {
		t.Fatalf("conflict message = %#v, want a named mtime refusal", message)
	}
	reloaded, reloadCmd := switched.Update(action)
	rolledBack := reloaded.(Model)
	if rolledBack.SelectedRelease != "R-001" {
		t.Errorf("conflict rollback selected %q, want prior R-001", rolledBack.SelectedRelease)
	}
	final := applyBoardCommands(t, rolledBack, reloadCmd)
	if final.SelectedRelease != "R-001" || final.State.Router.Release != "R-001" {
		t.Errorf("post-conflict selection = %q, router = %+v; want truthful R-001", final.SelectedRelease, final.State.Router)
	}
	if !strings.Contains(final.StatusMessage, "changed on disk") {
		t.Errorf("post-conflict status = %q, want the refusal reason retained", final.StatusMessage)
	}
}

func TestReleaseReloadPreservesFocusAndDiagnosesRemovedSelection(t *testing.T) {
	root := writeReleaseBoardProject(t)
	model := openSizedBoard(t, root, 130, 48)
	opened, _ := model.Update(keyMsg(goalSelectorKey))
	selector := opened.(Model)
	selector = press(t, selector, "j")
	updated, cmd := selector.Update(keyMsg("enter"))
	model = applyBoardCommands(t, updated.(Model), cmd)
	if model.SelectedRelease != "R-002" {
		t.Fatalf("selected Release = %q, want R-002", model.SelectedRelease)
	}
	focused := model.focusedTaskID()
	if focused != "T-002" {
		t.Fatalf("focused Task = %q, want T-002 in R-002", focused)
	}

	writeReleaseRecord(t, root, "R-002-second", "R-002", "Renamed second release")
	reloaded, _ := model.Update(loadCmd(root)().(projectLoadedMsg))
	model = reloaded.(Model)
	if model.SelectedRelease != "R-002" || model.focusedTaskID() != focused {
		t.Errorf("reload lost Release/focus: release=%q task=%q", model.SelectedRelease, model.focusedTaskID())
	}

	if err := os.RemoveAll(filepath.Join(root, "releases", "R-002-second")); err != nil {
		t.Fatal(err)
	}
	writeObjective(t, root, "O-002", "Second objective", "in_progress")
	reloaded, _ = model.Update(loadCmd(root)().(projectLoadedMsg))
	model = reloaded.(Model)
	if model.SelectedRelease != "" {
		t.Errorf("removed Release selection = %q, want no substitute", model.SelectedRelease)
	}
	if !strings.Contains(model.StatusMessage, "Goal R-002 no longer exists; selection cleared.") {
		t.Errorf("removed selection status = %q, want a Goal-specific explanation", model.StatusMessage)
	}
	if model.State.Next.SelectionDiagnostic == nil || model.State.Next.SelectionDiagnostic.Kind != data.SelectionReleaseNotFound {
		t.Errorf("selection diagnostic = %+v, want canonical release_not_found", model.State.Next.SelectionDiagnostic)
	}
}

func TestReloadClosingARemovedGoalDetailNamesItAGoal(t *testing.T) {
	root := writeReleaseBoardProject(t)
	writeReleaseRecord(t, root, "R-003-empty", "R-003", "Empty goal")
	model := press(t, openSizedBoard(t, root, 130, 48), goalSelectorKey, "j", "j", "v")
	if model.Detail == nil || model.Detail.ID != "R-003" {
		t.Fatalf("detail = %+v, want the R-003 Goal detail", model.Detail)
	}

	if err := os.RemoveAll(filepath.Join(root, "releases", "R-003-empty")); err != nil {
		t.Fatal(err)
	}
	reloaded, _ := model.Update(loadCmd(root)().(projectLoadedMsg))
	model = reloaded.(Model)
	if model.Detail != nil {
		t.Errorf("removed Goal detail stayed open: %+v", model.Detail)
	}
	if !strings.Contains(model.StatusMessage, "GOAL R-003 no longer exists; detail closed.") || strings.Contains(model.StatusMessage, "RELEASE") {
		t.Errorf("removed detail status = %q, want Goal wording", model.StatusMessage)
	}
}

func TestReleaseSelectorWithoutReleasesIsExplicitAndSafe(t *testing.T) {
	model := openSizedBoard(t, writeEmptyProjectFromTemplate(t), 130, 30)
	opened, cmd := model.Update(keyMsg(goalSelectorKey))
	selector := opened.(Model)
	if cmd != nil || !selector.ReleaseOverlay || !strings.Contains(xansi.Strip(selector.View()), "(no Goals in this project)") {
		t.Errorf("empty Release selector = open %t, cmd %t, view:\n%s", selector.ReleaseOverlay, cmd != nil, xansi.Strip(selector.View()))
	}
	selected, cmd := selector.Update(keyMsg("enter"))
	if cmd != nil || selected.(Model).SelectedRelease != "" {
		t.Errorf("Enter over no Releases changed context or returned a command: %#v", selected)
	}
}

func TestV2HelpAdvertisesTheCanonicalGoalSelectorKeyOnly(t *testing.T) {
	model := releaseBoard(t)
	opened, _ := model.Update(keyMsg("?"))
	view := xansi.Strip(opened.(Model).View())
	if !strings.Contains(view, "g: open the Goal selector") {
		t.Errorf("help does not advertise g:\n%s", view)
	}
	if strings.Contains(view, "r: open the Release selector") {
		t.Errorf("help exposes the legacy alias:\n%s", view)
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
