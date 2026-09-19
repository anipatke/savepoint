package v2

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	xansi "github.com/charmbracelet/x/ansi"
	"github.com/opencode/savepoint/internal/data"
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
	return openSizedBoard(t, writeReleaseBoardProject(t), 120, 48)
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
	model := openSizedBoard(t, root, 120, 48)
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
	if got := cardIDsInView(final); !equalIDs(got, []string{"T002"}) {
		t.Errorf("reloaded selected Release shows Tasks %v, want T002", got)
	}
	content, err := os.ReadFile(filepath.Join(root, "router.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "release: R002") || !strings.Contains(string(content), "next_action: \"do the next thing\"") {
		t.Errorf("router write lost the Release or next_action:\n%s", content)
	}
	if strings.Contains(string(content), "objective: none") || strings.Contains(string(content), "task: none") {
		t.Errorf("Release switch cleared existing router ownership:\n%s", content)
	}
	if string(beforeRouter) == string(content) {
		t.Error("canonical Release selection did not change router context")
	}
}

func TestReleaseSelectionConflictRollsBackAndReloads(t *testing.T) {
	root := writeReleaseBoardProject(t)
	model := openSizedBoard(t, root, 120, 48)
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
	model := openSizedBoard(t, root, 120, 48)
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
	model := openSizedBoard(t, root, 120, 48)
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
	model := openSizedBoard(t, writeEmptyProjectFromTemplate(t), 120, 30)
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
