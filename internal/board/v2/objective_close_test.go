package v2

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/testutil"
)

func closeableObjectiveProject(t *testing.T) (string, string, string) {
	t.Helper()
	root := writeEmptyProjectFromTemplate(t)
	addTestRelease(t, root, "R-006", "Test Goal")
	writeObjectiveExtra(t, root, "O-001", "Ready Objective", "in_progress",
		"release: R-006\nlast_check: C-001\n"+currentFreshness("C-001"))
	writeTask(t, root, "O-001", "T-001", "Done Task", "status: done\n")
	writeCheck(t, root, "C-001", "objective", "O-001", "CLEAR")
	// This Issue is router context only; it is unrelated to the Objective Check.
	testutil.WriteFile(t, filepath.Join(root, "issues", "I-002.md"),
		"---\nid: I-002\ntitle: Other Issue\ntype: defect\nstatus: open\nsource: {kind: report, actor: {role: owner, session: fixture}, at: '2026-01-02T00:00:00Z'}\n---\n\n# Other Issue\n")
	router := routerContent("task", "O-001", "T-001")
	router = strings.Replace(router, "objective: O-001\n", "release: R-006\nobjective: O-001\n", 1)
	router = strings.Replace(router, "task: T-001\n", "task: T-001\nissue: I-002\n", 1)
	path := filepath.Join(root, "router.md")
	testutil.WriteFile(t, path, router)
	return root, path, router
}

func TestSpaceClosesClearedObjectiveAndPreservesRouterContext(t *testing.T) {
	root, path, before := closeableObjectiveProject(t)
	model := focusObjective(t, openSizedBoard(t, root, 120, 40), "O-001")
	model.SidebarFocused = true
	if got := renderHelp(model, 120, 40); !strings.Contains(got, "space: close Objective") {
		t.Fatalf("help does not list permitted Objective closure: %s", got)
	}
	_, cmd := model.Update(keyMsg(" "))
	if cmd == nil {
		t.Fatal("Space returned no completion command")
	}
	message := cmd().(actionMsg)
	if message.err != nil || !message.reload {
		t.Fatalf("Objective completion = %+v", message)
	}
	loaded := loadProject(root)
	if loaded.Failed() || loaded.State.Index.Objectives["O-001"].Status != data.ColumnDone {
		t.Fatalf("Objective not done after reload: %s", loaded.Diagnostic)
	}
	want := strings.Replace(before, "objective: O-001\n", "objective: none\n", 1)
	want = strings.Replace(want, "task: T-001\n", "task: none\n", 1)
	if got := readActionFile(t, path); got != want {
		t.Fatalf("router bytes changed outside Objective and Task: got %q, want %q", got, want)
	}
	if got := writeObjectiveCompletionCmd(root, "O-001")().(actionMsg); got.err == nil || !strings.Contains(got.err.Error(), "already done") {
		t.Fatalf("repeat completion = %+v", got)
	}
}

func TestObjectiveSpaceRefusesBlockedAndExceptionOnlyGates(t *testing.T) {
	tests := []struct {
		name string
		edit func(*testing.T, string)
		want string
	}{
		{"unfinished Task", func(t *testing.T, root string) {
			writeTask(t, root, "O-001", "T-001", "Open", "status: in_progress\nstage: audit\n")
		}, "not done"},
		{"missing Check", func(t *testing.T, root string) {
			if err := os.Remove(filepath.Join(root, "checks", "C-001.md")); err != nil {
				t.Fatal(err)
			}
			writeObjectiveExtra(t, root, "O-001", "Ready Objective", "in_progress", "release: R-006\n")
		}, "missing"},
		{"needs work", func(t *testing.T, root string) { writeCheck(t, root, "C-001", "objective", "O-001", "NEEDS WORK") }, "needs work"},
		{"stale Check", func(t *testing.T, root string) { setObjectiveFreshness(t, root, "stale") }, "stale"},
		{"unknown Check", func(t *testing.T, root string) { setObjectiveFreshness(t, root, "unknown") }, "unknown"},
		{"owner acceptance", func(t *testing.T, root string) { setObjectiveExtra(t, root, "owner_validation:\n  required: true\n") }, "Acceptance"},
		{"open linked Issue", func(t *testing.T, root string) {
			writeCheckExtra(t, root, "C-001", "objective", "O-001", "CLEAR", "issues: [I-001]\n")
			writeIssue(t, root, "I-001", "Open finding", "defect", "C-001", "T-001")
		}, "Issue I-001"},
		{"exception only", func(t *testing.T, root string) {
			writeCheck(t, root, "C-001", "objective", "O-001", "NEEDS WORK")
			setObjectiveExtra(t, root, "exception:\n  requirements: [TEST-02]\n  reason: accepted risk\n  owner: owner-fixture\n  recorded_at: 2026-01-05T00:00:00Z\n  check: C-001\n")
		}, "owner session"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			root, path, before := closeableObjectiveProject(t)
			tc.edit(t, root)
			model := focusObjective(t, openSizedBoard(t, root, 120, 40), "O-001")
			model.SidebarFocused = true
			if got := renderHelp(model, 120, 40); strings.Contains(got, "space: close Objective") {
				t.Fatalf("help offers blocked closure: %s", got)
			}
			_, cmd := model.Update(keyMsg(" "))
			if cmd == nil {
				t.Fatal("Space returned no refusal command")
			}
			message := cmd().(actionMsg)
			if message.err == nil || !strings.Contains(message.err.Error(), tc.want) {
				t.Fatalf("refusal = %+v, want %q", message, tc.want)
			}
			if got := readActionFile(t, path); got != before {
				t.Fatalf("refusal changed router: %q", got)
			}
			index, err := data.LoadV2Index(root)
			if err != nil || index.Objectives["O-001"].Status != data.ColumnInProgress {
				t.Fatalf("refusal changed Objective: %v", err)
			}
		})
	}
}

func TestObjectiveCompletionReresolvesAfterBoardLoad(t *testing.T) {
	root, path, before := closeableObjectiveProject(t)
	model := focusObjective(t, openSizedBoard(t, root, 120, 40), "O-001")
	model.SidebarFocused = true
	writeCheck(t, root, "C-001", "objective", "O-001", "NEEDS WORK")
	_, cmd := model.Update(keyMsg(" "))
	if cmd == nil {
		t.Fatal("stale board returned no completion command")
	}
	message := cmd().(actionMsg)
	if message.err == nil || !strings.Contains(message.err.Error(), "needs work") {
		t.Fatalf("fresh gate refusal = %+v", message)
	}
	if got := readActionFile(t, path); got != before {
		t.Fatalf("fresh gate refusal changed router: %q", got)
	}
	index, err := data.LoadV2Index(root)
	if err != nil || index.Objectives["O-001"].Status != data.ColumnInProgress {
		t.Fatalf("fresh gate refusal changed Objective: %v", err)
	}
}

func setObjectiveExtra(t *testing.T, root, extra string) {
	t.Helper()
	path := objectivePath(root, "O-001")
	content := readActionFile(t, path)
	testutil.WriteFile(t, path, strings.Replace(content, "---\n\n#", extra+"---\n\n#", 1))
}

func setObjectiveFreshness(t *testing.T, root, state string) {
	t.Helper()
	path := objectivePath(root, "O-001")
	content := readActionFile(t, path)
	testutil.WriteFile(t, path, strings.Replace(content, "state: current", "state: "+state, 1))
}

func TestObjectiveRouterConflictReportsStaleSelection(t *testing.T) {
	root, path, _ := closeableObjectiveProject(t)
	objective := objectivePath(root, "O-001")
	content := readActionFile(t, objective)
	testutil.WriteFile(t, objective, strings.Replace(content, "status: in_progress", "status: done", 1))
	message := completedRecordAction(root, data.RouterSelectionV2{Objective: "O-001"}, "O-001 completed.", func(root string, selection data.RouterSelectionV2, expected time.Time) error {
		content := readActionFile(t, path)
		if err := os.WriteFile(path, []byte(content+"\n# concurrent edit\n"), 0644); err != nil {
			return err
		}
		return data.WriteRouterStateV2(root, selection, expected)
	}).(actionMsg)
	if message.err == nil || !message.reload || !strings.Contains(message.err.Error(), "router still selects O-001") {
		t.Fatalf("router conflict = %+v", message)
	}
}
