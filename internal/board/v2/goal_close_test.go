package v2

import (
	"path/filepath"
	"strings"
	"testing"

	xansi "github.com/charmbracelet/x/ansi"
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/testutil"
)

func releasePath(root, id string) string {
	return filepath.Join(root, "releases", id+"-test", "Release.md")
}

func goalSelectorOn(t *testing.T, root, id string) Model {
	t.Helper()
	model := openSizedBoard(t, root, 120, 40)
	updated, _ := model.Update(keyMsg(goalSelectorKey))
	model = updated.(Model)
	model.ReleaseCursor = releaseIndex(model.Releases, id)
	return model
}

func TestGoalToggleClosesAReadyGoalAndReopensIt(t *testing.T) {
	root, routerPath, routerBefore := closeableObjectiveProject(t)
	path := releasePath(root, "R-006")
	before := readActionFile(t, path)

	model := goalSelectorOn(t, root, "R-006")
	if got := xansi.Strip(model.View()); !strings.Contains(got, glyphCheckPending+" R-006") {
		t.Fatalf("open Goal row lacks the empty marker:\n%s", got)
	}
	_, cmd := model.Update(keyMsg(goalToggleKey))
	if cmd == nil {
		t.Fatal("C returned no Goal status command")
	}
	if message := cmd().(actionMsg); message.err != nil || !message.reload || message.message != "R-006 closed." {
		t.Fatalf("closing a ready Goal = %+v", message)
	}
	want := strings.Replace(before, "status: planned\n", "status: done\n", 1)
	if got := readActionFile(t, path); got != want {
		t.Fatalf("Goal write changed more than status: got %q, want %q", got, want)
	}
	if got := readActionFile(t, routerPath); got != routerBefore {
		t.Fatalf("closing a Goal changed the router: got %q, want %q", got, routerBefore)
	}

	closed := goalSelectorOn(t, root, "R-006")
	if got := xansi.Strip(closed.View()); !strings.Contains(got, glyphCheckClear+" R-006") {
		t.Fatalf("done Goal row lacks the ticked marker:\n%s", got)
	}
	if message := writeGoalStatusCmd(root, "R-006")().(actionMsg); message.err != nil || message.message != "R-006 reopened." {
		t.Fatalf("reopening a done Goal = %+v", message)
	}
	loaded := loadProject(root)
	if loaded.Failed() || loaded.State.Index.Releases["R-006"].Status != data.ColumnInProgress {
		t.Fatalf("Goal not in_progress after reopen: %s", loaded.Diagnostic)
	}
}

func TestGoalToggleRefusesAGoalWithAnIncompleteObjective(t *testing.T) {
	root, _, _ := closeableObjectiveProject(t)
	writeTask(t, root, "O-001", "T-001", "Open Task", "status: in_progress\nstage: build\n")
	path := releasePath(root, "R-006")
	before := readActionFile(t, path)

	message := writeGoalStatusCmd(root, "R-006")().(actionMsg)
	if message.err == nil || !strings.Contains(message.err.Error(), "cannot close R-006") || !strings.Contains(message.err.Error(), "O-001") {
		t.Fatalf("blocked Goal close = %+v, want a refusal naming O-001", message)
	}
	if got := readActionFile(t, path); got != before {
		t.Fatalf("refused Goal close wrote the record: %q", got)
	}
}

func TestGoalToggleRefusesAnArchivedHistoricalGoal(t *testing.T) {
	root, _, _ := closeableObjectiveProject(t)
	path := filepath.Join(root, "releases", "R-001-history", "Release.md")
	content := "---\nid: R-001\ntitle: \"Historical Goal\"\nstatus: done\n" +
		"legacy_completion:\n  source_path: releases/v1/PRD.md\n  archive_path: archive/v1/PRD.md\n" +
		"  sha256: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\n---\n\n" +
		"## Outcome\n\nPreserve historical completion.\n\n## Why\n\nArchived.\n\n" +
		"## Success Conditions\n\nDone.\n\n## Boundaries\n\nArchived.\n"
	testutil.WriteFile(t, path, content)

	message := writeGoalStatusCmd(root, "R-001")().(actionMsg)
	if message.err == nil || !strings.Contains(message.err.Error(), "archived historical Goal") {
		t.Fatalf("toggling an archived Goal = %+v, want a refusal", message)
	}
	if got := readActionFile(t, path); got != content {
		t.Fatalf("refused archived Goal toggle wrote the record: %q", got)
	}
}

func TestGoalSelectorAdvertisesTheToggleKey(t *testing.T) {
	root, _, _ := closeableObjectiveProject(t)
	model := goalSelectorOn(t, root, "R-006")
	if got := xansi.Strip(model.View()); !strings.Contains(got, goalToggleKey+":close/reopen") {
		t.Fatalf("Goal selector hints omit %s:\n%s", goalToggleKey, got)
	}
	if got := renderHelp(model, 120, 40); !strings.Contains(got, goalToggleKey) {
		t.Fatalf("help omits the Goal toggle key:\n%s", got)
	}
}
