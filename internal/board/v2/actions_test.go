package v2

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/testutil"
)

func TestActionsOnlyExposeOwnerAuthorityAndSelection(t *testing.T) {
	root := writeEvidenceProject(t)
	loaded := loadProject(root)
	if loaded.Failed() {
		t.Fatalf("loadProject() diagnostic = %q", loaded.Diagnostic)
	}

	if actions := actionsForRecord(loaded.State.Index, actionTarget{Kind: DetailTask, ID: "T002"}); len(actions) != 0 {
		t.Fatalf("checker-allowed Task actions = %+v, want none", actions)
	}
	if actions := actionsForRecord(loaded.State.Index, actionTarget{Kind: DetailTask, ID: "T006"}); !hasAction(actions, ActionAcceptCheck) {
		t.Fatalf("owner-wait Task actions = %+v, want acceptance", actions)
	}

	// Turn the recorded exception into a currently open completion decision.
	writeTask(t, root, "O001", "T005", "Open exception completion", "status: in_progress\nstage: audit\nlast_check: C004\nexception:\n  requirements: [TEST-02]\n  reason: \"ship the known gap\"\n  owner: \"owner-fixture\"\n  recorded_at: 2026-01-05T00:00:00Z\n  check: C004\n")
	loaded = loadProject(root)
	if loaded.Failed() {
		t.Fatalf("reloaded exception fixture diagnostic = %q", loaded.Diagnostic)
	}
	actions := actionsForRecord(loaded.State.Index, actionTarget{Kind: DetailTask, ID: "T005"})
	action, ok := actionForKey(actions, exceptionCloseKey)
	if !ok || !strings.Contains(action.Label, "ship the known gap") {
		t.Fatalf("exception actions = %+v, want an owner action naming the exception", actions)
	}

	if actions := actionsForRecord(loaded.State.Index, actionTarget{Kind: DetailTask, ID: "T001"}); len(actions) != 0 {
		t.Fatalf("executor-allowed Task actions = %+v, want none", actions)
	}
}

func TestHelpListsOnlyFocusedOwnerCapabilities(t *testing.T) {
	model := openSizedBoard(t, writeValidProject(t), 120, 40)
	model = press(t, model, "?")
	got := screen(model)
	if !strings.Contains(got, "p: record selection") {
		t.Errorf("help omits the available selection key:\n%s", got)
	}
	if strings.Contains(got, "a: accept") || strings.Contains(got, "x: complete") {
		t.Errorf("help offers owner completion capabilities not granted by the focused gate:\n%s", got)
	}
	if !strings.Contains(got, "executor session") {
		t.Errorf("help does not name the executor authority for the focused planned Task:\n%s", got)
	}
}

func TestEveryGateBlockerHasDistinctRefusalWithAuthority(t *testing.T) {
	tests := []struct {
		kind  data.GateBlockKind
		want  string
		other string
	}{
		{data.GateBlockReplan, "planner session", "dependency"},
		{data.GateBlockDependency, "executor session", "replan"},
		{data.GateBlockObjectiveDependency, "executor session", "replan"},
		{data.GateBlockClearanceMissing, "checker session", "stale"},
		{data.GateBlockClearanceNeedsWork, "executor session", "unknown"},
		{data.GateBlockClearanceStale, "checker session", "missing"},
		{data.GateBlockClearanceUnknown, "checker session", "stale"},
		{data.GateBlockOwnerAcceptance, "owner session", "checker authority"},
		{data.GateBlockCheckerAuthority, "checker session", "owner session"},
		{data.GateBlockInvalidState, "planner session", "dependency"},
	}
	seen := map[string]bool{}
	for _, test := range tests {
		got := refusalStatement(string(test.kind), "requirement detail", "", "")
		if !strings.Contains(got, test.want) || !strings.Contains(got, "requirement detail") {
			t.Errorf("%s refusal = %q, want detail and %q", test.kind, got, test.want)
		}
		if strings.Contains(got, test.other) {
			t.Errorf("%s refusal = %q, unexpectedly names %q", test.kind, got, test.other)
		}
		seen[got] = true
	}
	if len(seen) != len(tests) {
		t.Errorf("refusal statements = %d distinct, want %d", len(seen), len(tests))
	}
}

func TestOwnerAcceptanceCommandRereadsAndReloads(t *testing.T) {
	root := writeEvidenceProject(t)
	path := taskPath(root, "O001", "T006")
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	msg := writeOwnerAcceptanceCmd(root, actionTarget{Kind: DetailTask, ID: "T006"})()
	result, ok := msg.(actionMsg)
	if !ok || result.err != nil || !result.reload {
		t.Fatalf("owner acceptance result = %#v, want successful reload message", msg)
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) == string(before) || !strings.Contains(string(after), "accepted_check: C005") {
		t.Fatalf("accepted evidence = %q, want the current Check recorded", after)
	}

	reloaded := loadProject(root)
	if reloaded.Failed() {
		t.Fatalf("reload diagnostic = %q", reloaded.Diagnostic)
	}
	if got := reloaded.State.Index.Tasks["T006"].Evidence.OwnerValidation.AcceptedCheck; got != "C005" {
		t.Errorf("accepted Check after reload = %q, want C005", got)
	}

	secondBefore, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	second := writeOwnerAcceptanceCmd(root, actionTarget{Kind: DetailTask, ID: "T006"})()
	secondResult := second.(actionMsg)
	if secondResult.err == nil {
		t.Fatal("second acceptance = nil error, want the now-unavailable action refused without a write")
	}
	secondAfter, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if !secondAfter.ModTime().Equal(secondBefore.ModTime()) {
		t.Error("repeated acceptance changed the task mtime")
	}
}

func TestOwnerActionKeyReturnsTypedCommandAndReloadsThroughLoad(t *testing.T) {
	root := writeEvidenceProject(t)
	model := focusTask(t, openSizedBoard(t, root, 120, 48), "T006")

	updated, cmd := model.Update(keyMsg("a"))
	if cmd == nil {
		t.Fatal("acceptance key returned no command")
	}
	message, ok := cmd().(actionMsg)
	if !ok || message.err != nil || !message.reload {
		t.Fatalf("acceptance command message = %#v, want typed success", message)
	}

	updated, reload := updated.(Model).Update(message)
	if reload == nil {
		t.Fatal("successful owner action returned no reload command")
	}
	loadedMessage, ok := reload().(projectLoadedMsg)
	if !ok {
		t.Fatalf("reload command returned %T, want projectLoadedMsg", reload())
	}
	final, _ := updated.(Model).Update(loadedMessage)
	finalModel := final.(Model)
	if got := finalModel.State.Index.Tasks["T006"].Evidence.OwnerValidation.AcceptedCheck; got != "C005" {
		t.Errorf("reloaded acceptance = %q, want C005", got)
	}

	// The command is a Bubble Tea command, not a direct filesystem branch in
	// Update; this compile-time assertion keeps that contract visible here.
	var _ tea.Cmd = cmd
}

func TestSelectionCommandOnlyChangesRouterSelectionAndReloads(t *testing.T) {
	root := writeValidProject(t)
	path := filepath.Join(root, "router.md")
	testutil.WriteFile(t, path, "# Router prose\n\n## Current state\n\n```yaml\nstate: task\nobjective: O001\ntask: T001\nnext_action: \\\"human prose survives\\\"\n```\n\nBody survives.\n")
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	msg := writeSelectionCmd(root, data.RouterSelectionV2{Objective: "O001", Task: "T002"})()
	result, ok := msg.(actionMsg)
	if !ok || result.err != nil || !result.reload {
		t.Fatalf("selection result = %#v, want successful reload message", msg)
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(after), "task: T002") || !strings.Contains(string(after), "human prose survives") || !strings.Contains(string(after), "Body survives.") {
		t.Fatalf("selection write = %q, want only the selection changed", after)
	}
	if strings.Contains(string(after), "state: task\nobjective: O001\ntask: T001") {
		t.Errorf("selection write left the old task selected: %q", after)
	}
	if !strings.Contains(string(before), "Body survives.") {
		t.Fatal("selection fixture did not contain its preservation body")
	}
}

func TestWritesRefusePendingMigrationWithoutChangingFiles(t *testing.T) {
	root := writeEvidenceProject(t)
	path := taskPath(root, "O001", "T006")
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if operationID := createPendingOperation(t, root); operationID == "" {
		t.Fatal("pending operation has no ID")
	}

	for _, name := range []string{"acceptance", "selection"} {
		t.Run(name, func(t *testing.T) {
			var msg teaMsg
			switch name {
			case "acceptance":
				msg = writeOwnerAcceptanceCmd(root, actionTarget{Kind: DetailTask, ID: "T006"})()
			case "selection":
				msg = writeSelectionCmd(root, data.RouterSelectionV2{Objective: "O001", Task: "T006"})()
			}
			result := msg.(actionMsg)
			if result.err == nil || !strings.Contains(result.err.Error(), "savepoint migrate --recover") {
				t.Fatalf("result = %#v, want migration recovery guidance", result)
			}
			after, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if string(after) != string(before) {
				t.Error("pending migration write changed task bytes")
			}
		})
	}
}

// teaMsg keeps the pending-migration table compact while still making the
// command result assertion explicit at each call site.
type teaMsg interface{}

func hasAction(actions []BoardAction, kind ActionKind) bool {
	for _, action := range actions {
		if action.Kind == kind {
			return true
		}
	}
	return false
}

func TestExceptionCompletionCommandWritesOnlyWhenOwnerAuthorityRemains(t *testing.T) {
	root := writeEvidenceProject(t)
	writeTask(t, root, "O001", "T005", "Open exception completion", "status: in_progress\nstage: audit\nlast_check: C004\nexception:\n  requirements: [TEST-02]\n  reason: \"ship the known gap\"\n  owner: \"owner-fixture\"\n  recorded_at: 2026-01-05T00:00:00Z\n  check: C004\n")

	msg := writeExceptionCompletionCmd(root, actionTarget{Kind: DetailTask, ID: "T005"})()
	result := msg.(actionMsg)
	if result.err != nil || !result.reload {
		t.Fatalf("exception completion result = %#v, want successful reload message", result)
	}
	loaded := loadProject(root)
	if loaded.Failed() {
		t.Fatalf("reload diagnostic = %q", loaded.Diagnostic)
	}
	if got := loaded.State.Index.Tasks["T005"].Status; got != data.ColumnDone {
		t.Errorf("T005 status = %q, want done", got)
	}
}
