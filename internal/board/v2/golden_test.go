package v2

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/opencode/savepoint/internal/testutil"
)

// The golden files under testdata/golden are the board as it rendered before
// O-046 added anything. They are captured from a temporary project (TEST-04)
// and compared byte for byte, so the opt-in agent feature cannot change what an
// owner who never turns it on sees. Regenerate only on purpose:
//
//	UPDATE_GOLDEN=1 go test ./internal/board/v2 -run TestGolden

const (
	goldenDir       = "testdata/golden"
	goldenWidth     = 130
	goldenHeight    = 40
	updateGoldenEnv = "UPDATE_GOLDEN"
)

// goldenSurfaces renders every surface the opt-in feature could touch for the
// project at root: the TUI board with a Task and with an Objective focused, the
// help overlay for each, the footer hints those views carry, the views after the
// keys O-046 will bind, and the plain (non-TTY) output.
func goldenSurfaces(t *testing.T, root string) map[string]string {
	t.Helper()
	open := func() Model {
		model := openBoard(t, root, "")
		model.Width, model.Height = goldenWidth, goldenHeight
		return model
	}

	task := open()
	objective := press(t, open(), "left")
	if !objective.SidebarFocused {
		t.Fatal("golden fixture: left did not move focus to the Objective sidebar")
	}

	surfaces := map[string]string{
		"tui-task-focus":          task.View(),
		"tui-objective-focus":     objective.View(),
		"help-task-focus":         press(t, task, "?").View(),
		"help-objective-focus":    press(t, objective, "?").View(),
		"keys-unbound-task-focus": press(t, task, "s", "A", "c").View(),
		"keys-unbound-objective":  press(t, objective, "s", "A", "c").View(),
		"keys-unbound-help-task":  press(t, task, "?", "s", "A", "c").View(),
		"plain":                   renderPlain(task.State, selectedObjectiveForRelease(task.State, "", selectedRelease(task.State))),
	}
	return surfaces
}

// writeGoldenProject is the fixed fixture: the navigation project with the
// given agent_launcher lines appended to its config.
func writeGoldenProject(t *testing.T, launcher string) string {
	t.Helper()
	root := writeNavigationProject(t)
	if launcher != "" {
		testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n"+launcher)
	}
	return root
}

func checkGolden(t *testing.T, surfaces map[string]string) {
	t.Helper()
	names := make([]string, 0, len(surfaces))
	for name := range surfaces {
		names = append(names, name)
	}
	sort.Strings(names)

	update := os.Getenv(updateGoldenEnv) != ""
	for _, name := range names {
		path := filepath.Join(goldenDir, name+".golden")
		if update {
			if err := os.MkdirAll(goldenDir, 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte(surfaces[name]), 0o644); err != nil {
				t.Fatal(err)
			}
			continue
		}
		want, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read golden %s: %v (regenerate with %s=1)", path, err, updateGoldenEnv)
		}
		if surfaces[name] != string(want) {
			t.Errorf("%s differs from %s:\n--- got ---\n%s\n--- want ---\n%s", name, path, surfaces[name], want)
		}
	}
}

// TestGoldenBoardBaseline pins the board as it is with no agent_launcher
// config at all. It is the capture taken before any other O-046 change.
func TestGoldenBoardBaseline(t *testing.T) {
	checkGolden(t, goldenSurfaces(t, writeGoldenProject(t, "")))
}
