package v2

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

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

func TestRecordSelectionPreservesReleaseAndIssueForTaskAndObjective(t *testing.T) {
	tests := []struct {
		name       string
		targetKind DetailKind
		targetID   string
		wantTask   string
	}{
		{name: "Task", targetKind: DetailTask, targetID: "T-007", wantTask: "T-007"},
		{name: "Objective", targetKind: DetailObjective, targetID: "O-001", wantTask: "none"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			root := writeEvidenceProject(t)
			path, before := addTestReleaseAndRouterContext(t, root, "T-002")
			loaded := loadProject(root)
			if loaded.Failed() {
				t.Fatalf("loadProject() diagnostic = %q", loaded.Diagnostic)
			}

			model := Model{Root: root, State: loaded.State}
			message, ok := model.runAction(BoardAction{
				Kind:       ActionRecordSelection,
				TargetKind: tc.targetKind,
				TargetID:   tc.targetID,
			})().(actionMsg)
			if !ok || message.err != nil || !message.reload {
				t.Fatalf("record selection result = %#v, want successful reload", message)
			}

			want := strings.Replace(before, "task: T-002\n", "task: "+tc.wantTask+"\n", 1)
			if got := readActionFile(t, path); got != want {
				t.Fatalf("router selection write changed bytes outside the selected Task:\n got:\n%s\nwant:\n%s", got, want)
			}
		})
	}
}

func TestTaskCompletionMovesRouterToLowestUnfinishedBlockedTask(t *testing.T) {
	root := writeEvidenceProject(t)
	path, before := addTestReleaseAndRouterContext(t, root, "T-002")

	message := writeTaskAdvanceCmd(root, "T-002")().(actionMsg)
	if message.err != nil || !message.reload {
		t.Fatalf("completion result = %#v, want successful reload", message)
	}
	loaded := loadProject(root)
	if loaded.Failed() {
		t.Fatalf("loadProject() diagnostic = %q", loaded.Diagnostic)
	}
	if got := loaded.State.Index.Tasks["T-002"].Status; got != data.ColumnDone {
		t.Fatalf("T-002 status = %q, want done", got)
	}
	if decision := data.ResolveTaskStart(loaded.State.Index, "T-003"); decision.Allowed || len(decision.Blockers) == 0 {
		t.Fatalf("T-003 start decision = %+v, want the next Task blocked", decision)
	}
	if got := loaded.State.Router.Task; got != "T-003" {
		t.Fatalf("router Task = %q, want blocked lowest unfinished T-003", got)
	}
	if line := nextPanelText(loaded.State.Next); !strings.Contains(line, "T-003") {
		t.Fatalf("board Next line = %q, want the newly selected T-003", line)
	}
	want := strings.Replace(before, "task: T-002\n", "task: T-003\n", 1)
	if got := readActionFile(t, path); got != want {
		t.Fatalf("completion router write changed bytes outside task selection:\n got:\n%s\nwant:\n%s", got, want)
	}
}

func TestTaskWaiverCompletionAdvancesRouterAndKeepsReleaseAndIssue(t *testing.T) {
	root := writeEvidenceProject(t)
	path, before := addTestReleaseAndRouterContext(t, root, "T-007")
	writeTask(t, root, "O-001", "T-007", "Never checked", "status: in_progress\nstage: audit\n")

	message := writeTaskAdvanceCmd(root, "T-007")().(actionMsg)
	if message.err != nil || !message.reload || !strings.Contains(message.message, "waiver") {
		t.Fatalf("waiver completion result = %#v, want successful reload with waiver evidence", message)
	}
	loaded := loadProject(root)
	if loaded.Failed() {
		t.Fatalf("loadProject() diagnostic = %q", loaded.Diagnostic)
	}
	task := loaded.State.Index.Tasks["T-007"]
	if task.Status != data.ColumnDone || task.Evidence.CheckWaiver == nil {
		t.Fatalf("T-007 status/evidence = %q/%+v, want done with owner waiver", task.Status, task.Evidence)
	}
	if got := loaded.State.Router.Task; got != "T-002" {
		t.Fatalf("router Task = %q, want lowest unfinished T-002", got)
	}
	if line := nextPanelText(loaded.State.Next); !strings.Contains(line, "T-002") {
		t.Fatalf("board Next line = %q, want the newly selected T-002", line)
	}
	want := strings.Replace(before, "task: T-007\n", "task: T-002\n", 1)
	if got := readActionFile(t, path); got != want {
		t.Fatalf("waiver completion router write changed bytes outside task selection:\n got:\n%s\nwant:\n%s", got, want)
	}
}

func TestExceptionCompletionAdvancesOnlyWhenRouterSelectedClosedTask(t *testing.T) {
	tests := []struct {
		name         string
		selectedTask string
		wantTask     string
	}{
		{name: "selected", selectedTask: "T-005", wantTask: "T-002"},
		{name: "not selected", selectedTask: "T-002", wantTask: "T-002"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			root := writeEvidenceProject(t)
			path, before := addTestReleaseAndRouterContext(t, root, tc.selectedTask)
			writeTask(t, root, "O-001", "T-005", "Open exception completion", "status: in_progress\nstage: audit\nlast_check: C-004\nexception:\n  requirements: [TEST-02]\n  reason: \"ship the known gap\"\n  owner: \"owner-fixture\"\n  recorded_at: 2026-01-05T00:00:00Z\n  check: C-004\n")

			message := writeExceptionCompletionCmd(root, actionTarget{Kind: DetailTask, ID: "T-005"})().(actionMsg)
			if message.err != nil || !message.reload {
				t.Fatalf("exception completion result = %#v, want successful reload", message)
			}
			loaded := loadProject(root)
			if loaded.Failed() {
				t.Fatalf("loadProject() diagnostic = %q", loaded.Diagnostic)
			}
			if got := loaded.State.Index.Tasks["T-005"].Status; got != data.ColumnDone {
				t.Fatalf("T-005 status = %q, want done", got)
			}
			if got := loaded.State.Router.Task; got != tc.wantTask {
				t.Fatalf("router Task = %q, want %q", got, tc.wantTask)
			}
			want := before
			if tc.selectedTask == "T-005" {
				want = strings.Replace(before, "task: T-005\n", "task: T-002\n", 1)
			}
			if got := readActionFile(t, path); got != want {
				t.Fatalf("exception completion router bytes = %q, want %q", got, want)
			}
		})
	}
}

func TestRouterConflictAfterCompletionKeepsRecordAndReportsStaleSelection(t *testing.T) {
	root := writeEvidenceProject(t)
	path, _ := addTestReleaseAndRouterContext(t, root, "T-007")
	writeTask(t, root, "O-001", "T-007", "Already completed", "status: done\n")

	message := completedRecordAction(root, data.RouterSelectionV2{Objective: "O-001", Task: "T-007"}, "T-007 completed.", func(root string, selection data.RouterSelectionV2, expected time.Time) error {
		time.Sleep(10 * time.Millisecond)
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.WriteFile(path, append(content, []byte("\n# concurrent router edit\n")...), 0644); err != nil {
			return err
		}
		return data.WriteRouterStateV2(root, selection, expected)
	}).(actionMsg)
	if message.err == nil || !message.reload || !strings.Contains(message.err.Error(), "router still selects T-007") {
		t.Fatalf("completion conflict result = %#v, want reload warning naming stale selection", message)
	}
	loaded := loadProject(root)
	if loaded.Failed() {
		t.Fatalf("loadProject() diagnostic = %q", loaded.Diagnostic)
	}
	if got := loaded.State.Index.Tasks["T-007"].Status; got != data.ColumnDone {
		t.Fatalf("T-007 status = %q, want completed record retained", got)
	}
	if got := loaded.State.Router.Task; got != "T-007" {
		t.Fatalf("router Task = %q, want unchanged stale selection", got)
	}
	if got := readActionFile(t, path); !strings.Contains(got, "release: R-006\n") || !strings.Contains(got, "issue: I-001\n") {
		t.Fatalf("router context changed on conflict: %q", got)
	}
}

func TestLastTaskCompletionClearsRouterTaskAndShowsObjectiveCheck(t *testing.T) {
	root := writeEmptyProjectFromTemplate(t)
	addTestRelease(t, root, "R-006", "Test Release")
	writeObjectiveExtra(t, root, "O-001", "One Task Objective", "in_progress", "release: R-006\n")
	writeTask(t, root, "O-001", "T-001", "Last Task", "status: in_progress\nstage: audit\n")
	router := routerContent("task", "O-001", "T-001")
	router = strings.Replace(router, "objective: O-001\n", "release: R-006\nobjective: O-001\n", 1)
	router = strings.Replace(router, "task: T-001\n", "task: T-001\nissue: none\n", 1)
	path := filepath.Join(root, "router.md")
	testutil.WriteFile(t, path, router)

	message := writeTaskAdvanceCmd(root, "T-001")().(actionMsg)
	if message.err != nil || !message.reload {
		t.Fatalf("last Task completion result = %#v, want successful reload", message)
	}
	loaded := loadProject(root)
	if loaded.Failed() {
		t.Fatalf("loadProject() diagnostic = %q", loaded.Diagnostic)
	}
	if got := loaded.State.Router.Task; got != "" {
		t.Fatalf("router Task = %q, want cleared after the last Task", got)
	}
	if !strings.HasPrefix(nextPanelText(loaded.State.Next), "Check O-001 — ") {
		t.Fatalf("board Next line = %q, want the Objective Check rung", nextPanelText(loaded.State.Next))
	}
	want := strings.Replace(router, "task: T-001\n", "task: none\n", 1)
	if got := readActionFile(t, path); got != want {
		t.Fatalf("last Task completion changed bytes outside Task selection:\n got:\n%s\nwant:\n%s", got, want)
	}
}

func TestObjectiveExceptionCompletionClearsOnlyObjectiveAndTaskSelection(t *testing.T) {
	root := writeEmptyProjectFromTemplate(t)
	addTestRelease(t, root, "R-006", "Test Release")
	writeObjectiveExtra(t, root, "O-001", "Objective exception", "in_progress",
		"release: R-006\nlast_check: C-001\nexception:\n"+
			"  requirements: [TEST-02]\n  reason: \"owner accepted the documented integration gap\"\n"+"  owner: \"owner-fixture\"\n  recorded_at: 2026-01-05T00:00:00Z\n  check: C-001\n")
	writeCheck(t, root, "C-001", "objective", "O-001", "NEEDS WORK")
	writeTask(t, root, "O-001", "T-001", "Completed with owner waiver",
		"status: done\ncheck_waiver:\n  task: T-001\n  reason: \"owner completed without requesting a Task Check\"\n"+
			"  actor: {role: owner, session: owner-fixture}\n  recorded_at: 2026-01-05T00:00:00Z\n")
	router := routerContent("task", "O-001", "T-001")
	router = strings.Replace(router, "objective: O-001\n", "release: R-006\nobjective: O-001\n", 1)
	router = strings.Replace(router, "task: T-001\n", "task: T-001\nissue: none\n", 1)
	path := filepath.Join(root, "router.md")
	testutil.WriteFile(t, path, router)
	before := router

	message := writeExceptionCompletionCmd(root, actionTarget{Kind: DetailObjective, ID: "O-001"})().(actionMsg)
	if message.err != nil || !message.reload {
		t.Fatalf("Objective exception completion result = %#v, want successful reload; error: %v", message, message.err)
	}
	loaded := loadProject(root)
	if loaded.Failed() {
		t.Fatalf("loadProject() diagnostic = %q", loaded.Diagnostic)
	}
	if got := loaded.State.Index.Objectives["O-001"].Status; got != data.ColumnDone {
		t.Fatalf("O-001 status = %q, want done", got)
	}
	if got := loaded.State.Router.Objective; got != "" {
		t.Fatalf("router Objective = %q, want cleared", got)
	}
	if got := loaded.State.Router.Task; got != "" {
		t.Fatalf("router Task = %q, want cleared", got)
	}
	want := strings.Replace(before, "objective: O-001\n", "objective: none\n", 1)
	want = strings.Replace(want, "task: T-001\n", "task: none\n", 1)
	if got := readActionFile(t, path); got != want {
		t.Fatalf("Objective closure changed bytes outside objective/task selection:\n got:\n%s\nwant:\n%s", got, want)
	}
}

func TestGoalSelectionPreservesIssueContext(t *testing.T) {
	root := writeEvidenceProject(t)
	path, before := addTestReleaseAndRouterContext(t, root, "T-002")
	addTestRelease(t, root, "R-007", "Second test Release")
	objectiveFile := objectivePath(root, "O-002")
	objectiveBytes, err := os.ReadFile(objectiveFile)
	if err != nil {
		t.Fatal(err)
	}
	objective := strings.Replace(string(objectiveBytes), "status: planned\n", "status: planned\nrelease: R-007\n", 1)
	if objective == string(objectiveBytes) {
		t.Fatal("O-002 fixture has no planned status to attach Release")
	}
	testutil.WriteFile(t, objectiveFile, objective)

	message := writeReleaseSelectionCmd(root, "R-007")().(actionMsg)
	if message.err != nil || !message.reload {
		t.Fatalf("Goal selection result = %#v, want successful reload", message)
	}
	want := strings.Replace(before, "release: R-006\n", "release: R-007\n", 1)
	want = strings.Replace(want, "objective: O-001\n", "objective: none\n", 1)
	want = strings.Replace(want, "task: T-002\n", "task: none\n", 1)
	if got := readActionFile(t, path); got != want {
		t.Fatalf("Goal selection changed bytes outside its Release and matching Objective/Task keys:\n got:\n%s\nwant:\n%s", got, want)
	}
}

func addTestReleaseAndRouterContext(t *testing.T, root, selectedTask string) (string, string) {
	t.Helper()
	addTestRelease(t, root, "R-006", "Test Release")

	objectiveFile := objectivePath(root, "O-001")
	objectiveBytes, err := os.ReadFile(objectiveFile)
	if err != nil {
		t.Fatal(err)
	}
	objective := strings.Replace(string(objectiveBytes), "status: in_progress\n", "status: in_progress\nrelease: R-006\n", 1)
	if objective == string(objectiveBytes) {
		t.Fatal("Objective fixture has no in_progress status to attach Release")
	}
	testutil.WriteFile(t, objectiveFile, objective)

	path := filepath.Join(root, "router.md")
	router := routerContent("task", "O-001", selectedTask)
	router = strings.Replace(router, "objective: O-001\n", "release: R-006\nobjective: O-001\n", 1)
	router = strings.Replace(router, "task: "+selectedTask+"\n", "task: "+selectedTask+"\nissue: I-001\n", 1)
	testutil.WriteFile(t, path, router)
	return path, router
}

func addTestRelease(t *testing.T, root, id, title string) {
	t.Helper()
	releaseDir := filepath.Join(root, "releases", id+"-test")
	testutil.WriteFile(t, filepath.Join(releaseDir, "Release.md"), "---\nid: "+id+"\ntitle: "+title+"\nstatus: planned\n---\n\n## Outcome\n\nA test Release.\n\n## Why\n\nA test Release fixture.\n\n## Success Conditions\n\n- The record loads.\n\n## Boundaries\n\n- Test fixture only.\n")
}

func readActionFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}

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

// writeObjectiveStartProject is one Objective with one planned Task, selected
// by the router, with the Objective at the given status.
func writeObjectiveStartProject(t *testing.T, objectiveStatus string) string {
	t.Helper()
	root := savepointRoot(t)
	writeConfig(t, root)
	writeRouter(t, root, "task", "O-001", "T-001")
	writeObjective(t, root, "O-001", "First objective", objectiveStatus)
	writeTask(t, root, "O-001", "T-001", "First task", "status: planned\n")
	return root
}

// TestTaskStartMovesPlannedObjectiveToInProgress is I-044's regression test:
// every earlier fixture wrote its Objective already in_progress, so nothing
// covered starting work under a planned one.
func TestTaskStartMovesPlannedObjectiveToInProgress(t *testing.T) {
	root := writeObjectiveStartProject(t, "planned")
	before := readActionFile(t, objectivePath(root, "O-001"))

	result := writeTaskAdvanceCmd(root, "T-001")().(actionMsg)
	if result.err != nil || !result.reload || !strings.Contains(result.message, "Objective O-001 is now in progress") {
		t.Fatalf("start result = %#v, want a successful reload naming the Objective", result)
	}
	loaded := loadProject(root)
	if loaded.Failed() {
		t.Fatalf("reload diagnostic = %q", loaded.Diagnostic)
	}
	if got := loaded.State.Index.Objectives["O-001"].Status; got != data.ColumnInProgress {
		t.Fatalf("O-001 status = %q, want in_progress after its first Task started", got)
	}
	want := strings.Replace(before, "status: planned\n", "status: in_progress\n", 1)
	if got := readActionFile(t, objectivePath(root, "O-001")); got != want {
		t.Errorf("Objective write changed bytes beyond status:\n got:\n%s\nwant:\n%s", got, want)
	}
	if got, want := nextPanelText(loaded.State.Next), "Build T-001 — First task (O-001)"; got != want {
		t.Errorf("Next line = %q, want %q", got, want)
	}
}

func TestTaskStartLeavesInProgressAndDoneObjectivesUntouched(t *testing.T) {
	root := writeObjectiveStartProject(t, "in_progress")
	path := objectivePath(root, "O-001")
	before := readActionFile(t, path)
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	result := writeTaskAdvanceCmd(root, "T-001")().(actionMsg)
	if result.err != nil || strings.Contains(result.message, "now in progress") {
		t.Fatalf("start result = %#v, want a plain start with no Objective change", result)
	}
	after, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := readActionFile(t, path); got != before || !after.ModTime().Equal(info.ModTime()) {
		t.Errorf("an in_progress Objective was rewritten on Task start")
	}

	loaded := loadProject(root)
	objective := loaded.State.Index.Objectives["O-001"]
	objective.Status = data.ColumnDone
	msg := startedTaskAction(loaded.State.Index, loaded.State.Index.Tasks["T-001"], "T-001 moved.", func(*data.ObjectiveV2) error {
		t.Fatal("a done Objective must never be written on Task start")
		return nil
	}).(actionMsg)
	if msg.err != nil {
		t.Errorf("done Objective start result = %#v, want no error", msg)
	}
}

func TestStartedTaskActionReportsObjectiveWriteFailure(t *testing.T) {
	root := writeObjectiveStartProject(t, "planned")
	loaded := loadProject(root)

	msg := startedTaskAction(loaded.State.Index, loaded.State.Index.Tasks["T-001"], "T-001 moved.", func(*data.ObjectiveV2) error {
		return errors.New("disk full")
	}).(actionMsg)
	if msg.err == nil || !msg.reload || !strings.Contains(msg.err.Error(), "O-001 is still planned") || !strings.Contains(msg.err.Error(), "disk full") {
		t.Fatalf("result = %#v, want a reloading error naming the still-planned Objective and the cause", msg)
	}
}

func TestTaskRetreatLeavesObjectiveStatusAlone(t *testing.T) {
	root := writeObjectiveStartProject(t, "planned")
	if result := writeTaskAdvanceCmd(root, "T-001")().(actionMsg); result.err != nil {
		t.Fatalf("start error = %v", result.err)
	}
	if result := writeTaskRetreatCmd(root, "T-001")().(actionMsg); result.err != nil {
		t.Fatalf("retreat error = %v", result.err)
	}
	loaded := loadProject(root)
	if got := loaded.State.Index.Tasks["T-001"].Status; got != data.ColumnPlanned {
		t.Fatalf("T-001 status = %q, want planned after retreat", got)
	}
	if got := loaded.State.Index.Objectives["O-001"].Status; got != data.ColumnInProgress {
		t.Errorf("O-001 status = %q, want in_progress kept after the Task retreated", got)
	}
}
