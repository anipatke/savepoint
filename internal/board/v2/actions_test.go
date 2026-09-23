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

	if actions := actionsForRecord(loaded.State.Index, actionTarget{Kind: DetailTask, ID: "T-002"}); len(actions) != 0 {
		t.Fatalf("checker-allowed Task actions = %+v, want none", actions)
	}
	if actions := actionsForRecord(loaded.State.Index, actionTarget{Kind: DetailTask, ID: "T-006"}); !hasAction(actions, ActionAcceptCheck) {
		t.Fatalf("owner-wait Task actions = %+v, want acceptance", actions)
	}

	// Turn the recorded exception into a currently open completion decision.
	writeTask(t, root, "O-001", "T-005", "Open exception completion", "status: in_progress\nstage: audit\nlast_check: C-004\nexception:\n  requirements: [TEST-02]\n  reason: \"ship the known gap\"\n  owner: \"owner-fixture\"\n  recorded_at: 2026-01-05T00:00:00Z\n  check: C-004\n")
	loaded = loadProject(root)
	if loaded.Failed() {
		t.Fatalf("reloaded exception fixture diagnostic = %q", loaded.Diagnostic)
	}
	actions := actionsForRecord(loaded.State.Index, actionTarget{Kind: DetailTask, ID: "T-005"})
	action, ok := actionForKey(actions, exceptionCloseKey)
	if !ok || action.Label != "complete by exception" {
		t.Fatalf("exception actions = %+v, want a short owner action label, not the exception's own reason text (that belongs in the detail overlay, not the footer)", actions)
	}
	if strings.Contains(action.Label, "ship the known gap") {
		t.Fatalf("exception action label = %q, leaked the recorded reason text into what must stay a short footer hint", action.Label)
	}

	if actions := actionsForRecord(loaded.State.Index, actionTarget{Kind: DetailTask, ID: "T-001"}); len(actions) != 0 {
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
	path := taskPath(root, "O-001", "T-006")
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	msg := writeOwnerAcceptanceCmd(root, actionTarget{Kind: DetailTask, ID: "T-006"})()
	result, ok := msg.(actionMsg)
	if !ok || result.err != nil || !result.reload {
		t.Fatalf("owner acceptance result = %#v, want successful reload message", msg)
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) == string(before) || !strings.Contains(string(after), "accepted_check: C-005") {
		t.Fatalf("accepted evidence = %q, want the current Check recorded", after)
	}

	reloaded := loadProject(root)
	if reloaded.Failed() {
		t.Fatalf("reload diagnostic = %q", reloaded.Diagnostic)
	}
	if got := reloaded.State.Index.Tasks["T-006"].Evidence.OwnerValidation.AcceptedCheck; got != "C-005" {
		t.Errorf("accepted Check after reload = %q, want C-005", got)
	}

	secondBefore, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	second := writeOwnerAcceptanceCmd(root, actionTarget{Kind: DetailTask, ID: "T-006"})()
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
	model := focusTask(t, openSizedBoard(t, root, 120, 48), "T-006")

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
	if got := finalModel.State.Index.Tasks["T-006"].Evidence.OwnerValidation.AcceptedCheck; got != "C-005" {
		t.Errorf("reloaded acceptance = %q, want C-005", got)
	}

	// The command is a Bubble Tea command, not a direct filesystem branch in
	// Update; this compile-time assertion keeps that contract visible here.
	var _ tea.Cmd = cmd
}

func TestSelectionCommandOnlyChangesRouterSelectionAndReloads(t *testing.T) {
	root := writeValidProject(t)
	path := filepath.Join(root, "router.md")
	testutil.WriteFile(t, path, "# Router prose\n\n## Current state\n\n```yaml\nstate: task\nobjective: O-001\ntask: T-001\nnext_action: \\\"human prose survives\\\"\n```\n\nBody survives.\n")
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	msg := writeSelectionCmd(root, data.RouterSelectionV2{Objective: "O-001", Task: "T-002"})()
	result, ok := msg.(actionMsg)
	if !ok || result.err != nil || !result.reload {
		t.Fatalf("selection result = %#v, want successful reload message", msg)
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(after), "task: T-002") || !strings.Contains(string(after), "human prose survives") || !strings.Contains(string(after), "Body survives.") {
		t.Fatalf("selection write = %q, want only the selection changed", after)
	}
	if strings.Contains(string(after), "state: task\nobjective: O-001\ntask: T-001") {
		t.Errorf("selection write left the old task selected: %q", after)
	}
	if !strings.Contains(string(before), "Body survives.") {
		t.Fatal("selection fixture did not contain its preservation body")
	}
}

func TestWritesRefusePendingMigrationWithoutChangingFiles(t *testing.T) {
	root := writeEvidenceProject(t)
	path := taskPath(root, "O-001", "T-006")
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
				msg = writeOwnerAcceptanceCmd(root, actionTarget{Kind: DetailTask, ID: "T-006"})()
			case "selection":
				msg = writeSelectionCmd(root, data.RouterSelectionV2{Objective: "O-001", Task: "T-006"})()
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

// TestTaskAdvanceCommandStartsAndAdvancesNormally covers Space's ordinary
// path through the lifecycle it shares with ResolveTaskStart/ResolveTaskAdvance:
// planned moves to in_progress/build, and build moves to test.
func TestTaskAdvanceCommandStartsAndAdvancesNormally(t *testing.T) {
	root := writeEvidenceProject(t)
	writeTask(t, root, "O-002", "T-900", "Not started yet", "status: planned\n")

	msg := writeTaskAdvanceCmd(root, "T-900")()
	result := msg.(actionMsg)
	if result.err != nil || !result.reload {
		t.Fatalf("start result = %#v, want successful reload", result)
	}
	loaded := loadProject(root)
	if loaded.Failed() {
		t.Fatalf("reload diagnostic = %q", loaded.Diagnostic)
	}
	task := loaded.State.Index.Tasks["T-900"]
	if task.Status != data.ColumnInProgress || task.Stage != data.StageBuild {
		t.Fatalf("T-900 status/stage = %q/%q, want in_progress/build after starting", task.Status, task.Stage)
	}

	msg = writeTaskAdvanceCmd(root, "T-900")()
	result = msg.(actionMsg)
	if result.err != nil || !result.reload {
		t.Fatalf("advance result = %#v, want successful reload", result)
	}
	loaded = loadProject(root)
	if got := loaded.State.Index.Tasks["T-900"].Stage; got != data.StageTest {
		t.Errorf("T-900 stage = %q, want test after advancing from build", got)
	}
}

// TestTaskAdvanceCommandCompletingWithNoCheckRequestedRecordsAnOwnerWaiver is
// the behavior an owner asked for directly: pressing Space to complete a
// Task at stage check that has no recorded Check at all is itself treated as
// the explicit Task-check waiver TEST-09 requires — recorded automatically,
// not demanded of the owner as a separate manual step first.
func TestTaskAdvanceCommandCompletingWithNoCheckRequestedRecordsAnOwnerWaiver(t *testing.T) {
	root := writeEvidenceProject(t)
	writeTask(t, root, "O-001", "T-900", "Never checked", "status: in_progress\nstage: audit\n")

	msg := writeTaskAdvanceCmd(root, "T-900")()
	result := msg.(actionMsg)
	if result.err != nil || !result.reload {
		t.Fatalf("completion result = %#v, want successful reload", result)
	}
	if !strings.Contains(result.message, "waiver") {
		t.Errorf("completion message = %q, want it to say a waiver was recorded", result.message)
	}

	loaded := loadProject(root)
	if loaded.Failed() {
		t.Fatalf("reload diagnostic = %q", loaded.Diagnostic)
	}
	task := loaded.State.Index.Tasks["T-900"]
	if task.Status != data.ColumnDone {
		t.Fatalf("T-900 status = %q, want done", task.Status)
	}
	waiver := task.Evidence.CheckWaiver
	if waiver == nil || waiver.Task != "T-900" {
		t.Fatalf("CheckWaiver = %+v, want one naming T-900", waiver)
	}
	if waiver.Actor.Role != data.ActorRoleOwner || waiver.Actor.Session != ownerBoardSession {
		t.Errorf("CheckWaiver.Actor = %+v, want owner/%s", waiver.Actor, ownerBoardSession)
	}
}

// TestTaskAdvanceCommandNeverAutoWaivesARealCheckFinding proves the
// auto-waiver applies only when no Check was ever requested. A Task whose
// recorded Check actually found a problem stays blocked — Space never
// silently overrides an independent checker's NEEDS WORK.
func TestTaskAdvanceCommandNeverAutoWaivesARealCheckFinding(t *testing.T) {
	root := writeEvidenceProject(t)
	writeCheck(t, root, "C-900", "task", "T-900", "NEEDS WORK")
	writeTask(t, root, "O-001", "T-900", "Checked and found wanting", "status: in_progress\nstage: audit\nlast_check: C-900\n")

	msg := writeTaskAdvanceCmd(root, "T-900")()
	result := msg.(actionMsg)
	if result.err == nil {
		t.Fatalf("completion result = %#v, want a refusal for a NEEDS WORK check", result)
	}
	if strings.Contains(result.err.Error(), "waiver") {
		t.Errorf("refusal = %q, want no mention of a waiver being an option here", result.err.Error())
	}

	loaded := loadProject(root)
	if loaded.Failed() {
		t.Fatalf("reload diagnostic = %q", loaded.Diagnostic)
	}
	task := loaded.State.Index.Tasks["T-900"]
	if task.Status != data.ColumnInProgress || task.Evidence.CheckWaiver != nil {
		t.Fatalf("T-900 = status %q, CheckWaiver %+v, want unchanged and unwaived", task.Status, task.Evidence.CheckWaiver)
	}
}

// TestTaskRetreatCommandMovesBackOneStepAndIsUngated proves Backspace moves a
// Task backward through the lifecycle without any Check gate — only the
// owner's own keypress reaches it.
func TestTaskRetreatCommandMovesBackOneStepAndIsUngated(t *testing.T) {
	root := writeEvidenceProject(t)
	writeTask(t, root, "O-001", "T-900", "Mid test", "status: in_progress\nstage: test\n")

	msg := writeTaskRetreatCmd(root, "T-900")()
	result := msg.(actionMsg)
	if result.err != nil || !result.reload {
		t.Fatalf("retreat result = %#v, want successful reload", result)
	}
	loaded := loadProject(root)
	if loaded.Failed() {
		t.Fatalf("reload diagnostic = %q", loaded.Diagnostic)
	}
	if got := loaded.State.Index.Tasks["T-900"].Stage; got != data.StageBuild {
		t.Errorf("T-900 stage = %q, want build after retreating from test", got)
	}
}

func TestExceptionCompletionCommandWritesOnlyWhenOwnerAuthorityRemains(t *testing.T) {
	root := writeEvidenceProject(t)
	writeTask(t, root, "O-001", "T-005", "Open exception completion", "status: in_progress\nstage: audit\nlast_check: C-004\nexception:\n  requirements: [TEST-02]\n  reason: \"ship the known gap\"\n  owner: \"owner-fixture\"\n  recorded_at: 2026-01-05T00:00:00Z\n  check: C-004\n")

	msg := writeExceptionCompletionCmd(root, actionTarget{Kind: DetailTask, ID: "T-005"})()
	result := msg.(actionMsg)
	if result.err != nil || !result.reload {
		t.Fatalf("exception completion result = %#v, want successful reload message", result)
	}
	loaded := loadProject(root)
	if loaded.Failed() {
		t.Fatalf("reload diagnostic = %q", loaded.Diagnostic)
	}
	if got := loaded.State.Index.Tasks["T-005"].Status; got != data.ColumnDone {
		t.Errorf("T-005 status = %q, want done", got)
	}
}
