package v2

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/testutil"
)

// scaffoldConfig reads the config.yml a fresh `savepoint init` would write, so
// the journey starts from the shipped default rather than a hand-made one.
func scaffoldConfig(t *testing.T) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join("..", "..", "..", "templates", "project-v2", ".savepoint", "config.yml"))
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}

// projectFiles maps each file under root to its content, keyed by relative
// path, so two temporary projects can be compared.
func projectFiles(t *testing.T, root string) map[string]string {
	t.Helper()
	files := map[string]string{}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		files[filepath.ToSlash(rel)] = string(content)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	return files
}

// TestOptionsJourneyOnTemporaryProject walks the owner's check end to end with
// no lane data anywhere: save, restart, switch off and on, an external edit, a
// stale save, and close. It also pins that only config.yml is ever written.
func TestOptionsJourneyOnTemporaryProject(t *testing.T) {
	root := writeValidProject(t)
	scaffold := scaffoldConfig(t)
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), scaffold)
	before := projectFiles(t, root)

	model := press(t, openSizedBoard(t, root, 130, 40), optionsKey)
	if !strings.Contains(optionsScreen(model), "OFF") {
		t.Fatalf("fresh scaffold did not open with the option off:\n%s", optionsScreen(model))
	}

	on := toggleOptions(t, model)
	if !strings.Contains(optionsScreen(on), "Saved. Parallel planning is on.") {
		t.Fatalf("enable not confirmed:\n%s", optionsScreen(on))
	}

	restarted := press(t, openSizedBoard(t, root, 130, 40), optionsKey)
	if !strings.Contains(optionsScreen(restarted), "ON") {
		t.Fatalf("restart lost the preference:\n%s", optionsScreen(restarted))
	}

	off := toggleOptions(t, restarted)
	if !strings.Contains(optionsScreen(off), "Saved. Parallel planning is off.") {
		t.Fatalf("switch-off not confirmed:\n%s", optionsScreen(off))
	}
	if got := readOptionsConfig(t, off); got != scaffold {
		t.Errorf("config.yml after switching off = %q", got)
	}

	// An edit made outside the board, while the screen is open, wins: the next
	// save is refused once, then retried against the file as it now stands.
	external := readOptionsConfig(t, off) + "# edited by the owner\n"
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), external)
	refused := toggleOptions(t, off)
	if !strings.Contains(optionsScreen(refused), "Nothing was saved") {
		t.Fatalf("stale save was not refused:\n%s", optionsScreen(refused))
	}
	if got := readOptionsConfig(t, refused); got != external {
		t.Fatalf("stale save overwrote the external edit: %q", got)
	}
	reOn := toggleOptions(t, refused)
	if got := readOptionsConfig(t, reOn); !strings.Contains(got, "# edited by the owner") || !strings.Contains(got, "parallel_planning: true") {
		t.Fatalf("retry lost the edit or the value: %q", got)
	}

	closed := press(t, reOn, "esc")
	if closed.Options != nil {
		t.Error("esc did not close Advanced Options")
	}

	after := projectFiles(t, root)
	for path, content := range before {
		if path != "config.yml" && after[path] != content {
			t.Errorf("%s changed during the journey", path)
		}
	}
	for path := range after {
		if _, existed := before[path]; !existed {
			t.Errorf("journey created %s; only config.yml may be written", path)
		}
	}
}

// TestOptionsLeaveLifecycleAndBoardIdenticalWhenOn runs the same board action
// in two identical projects, one with the option saved on, and requires the
// same board and the same record files: the option changes no lifecycle or
// Code Health behavior.
func TestOptionsLeaveLifecycleAndBoardIdenticalWhenOn(t *testing.T) {
	build := func(on bool) (Model, string) {
		root := writeValidProject(t)
		testutil.WriteFile(t, filepath.Join(root, "config.yml"), scaffoldConfig(t))
		model := openSizedBoard(t, root, 130, 40)
		if on {
			model = press(t, toggleOptions(t, press(t, model, optionsKey)), "esc")
		}
		return model, root
	}
	offModel, offRoot := build(false)
	onModel, onRoot := build(true)

	if on, off := optionsScreen(onModel), optionsScreen(offModel); on != off {
		t.Errorf("board differs with the option on:\noff:\n%s\non:\n%s", off, on)
	}
	if onModel.State.HealthChip != offModel.State.HealthChip {
		t.Errorf("Code Health chip differs: off %+v, on %+v", offModel.State.HealthChip, onModel.State.HealthChip)
	}

	next := func(m Model) Model {
		updated, cmd := pressOptionsKey(t, m, " ")
		return applyBoardCommands(t, updated, cmd)
	}
	offAfter, onAfter := next(offModel), next(onModel)
	if optionsScreen(onAfter) != optionsScreen(offAfter) {
		t.Errorf("advancing a Task renders differently with the option on")
	}
	offFiles, onFiles := projectFiles(t, offRoot), projectFiles(t, onRoot)
	delete(offFiles, "config.yml")
	delete(onFiles, "config.yml")
	if len(offFiles) != len(onFiles) {
		t.Fatalf("file sets differ: off %d, on %d", len(offFiles), len(onFiles))
	}
	for path, content := range offFiles {
		if onFiles[path] != content {
			t.Errorf("%s differs with the option on", path)
		}
	}
}
