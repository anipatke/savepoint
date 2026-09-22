package v2

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/testutil"
)

func TestV2WatchSetIncludesLiveFilesAndExcludesHistoricalTrees(t *testing.T) {
	root := savepointRoot(t)
	watched := []string{
		"objectives/O001-fixture/Objective.md",
		"objectives/O001-fixture/tasks/T001-fixture.md",
		"checks/C001.md",
		"issues/I001.md",
		"releases/R001-first/Release.md",
		"router.md",
		"config.yml",
		".migration/op/operation.yml",
	}
	for _, rel := range watched {
		if !isV2WatchedPath(root, filepath.Join(root, filepath.FromSlash(rel))) {
			t.Errorf("%s is not in the V2 watch set", rel)
		}
	}
	for _, rel := range []string{
		"archive/v1/router.md",
		"defects/D001.md",
		"audit/register.md",
		"migrations/v1-to-v2.yml",
	} {
		if isV2WatchedPath(root, filepath.Join(root, filepath.FromSlash(rel))) {
			t.Errorf("%s is excluded from the V2 watch set but was accepted", rel)
		}
	}
}

func TestV2WatcherReloadsLiveFilesAndLearnsNewObjectiveDirectories(t *testing.T) {
	root := writeValidProject(t)
	watcher, err := newV2Watcher(root)
	if err != nil {
		t.Fatalf("newV2Watcher() error = %v", err)
	}
	t.Cleanup(func() { _ = watcher.Close() })

	testutil.WriteFile(t, filepath.Join(root, "objectives", "O002-fixture", "Objective.md"),
		"---\nid: O002\ntitle: \"Second objective\"\nstatus: planned\n---\n")
	awaitV2FileChange(t, watchV2Files(watcher, root))

	// The event above registered the new Objective directory. A later edit to
	// that directory must be visible without rebuilding the watcher.
	testutil.WriteFile(t, filepath.Join(root, "objectives", "O002-fixture", "Objective.md"),
		"---\nid: O002\ntitle: \"Renamed objective\"\nstatus: planned\n---\n")
	awaitV2FileChange(t, watchV2Files(watcher, root))
}

func TestV2WatcherIgnoresExcludedTrees(t *testing.T) {
	root := writeValidProject(t)
	watcher, err := newV2Watcher(root)
	if err != nil {
		t.Fatalf("newV2Watcher() error = %v", err)
	}
	t.Cleanup(func() { _ = watcher.Close() })

	result := make(chan tea.Msg, 1)
	go func() { result <- watchV2Files(watcher, root)() }()
	for _, directory := range []string{"archive", "defects", "audit"} {
		testutil.WriteFile(t, filepath.Join(root, directory, "changed.md"), "historical\n")
	}
	select {
	case msg := <-result:
		t.Fatalf("excluded tree produced reload message %T", msg)
	case <-time.After(250 * time.Millisecond):
	}
}

func TestV2WatcherDebouncesRapidWrites(t *testing.T) {
	root := writeValidProject(t)
	watcher, err := newV2Watcher(root)
	if err != nil {
		t.Fatalf("newV2Watcher() error = %v", err)
	}
	t.Cleanup(func() { _ = watcher.Close() })

	result := make(chan tea.Msg, 1)
	go func() { result <- watchV2Files(watcher, root)() }()
	config := filepath.Join(root, "config.yml")
	for i := 0; i < 5; i++ {
		testutil.WriteFile(t, config, "schema_version: 2\n# edit "+string(rune('a'+i))+"\n")
		time.Sleep(10 * time.Millisecond)
	}
	select {
	case msg := <-result:
		if _, ok := msg.(v2FileChangeMsg); !ok {
			t.Fatalf("watch command returned %T, want v2FileChangeMsg", msg)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("watch command did not emit a debounced reload")
	}

	second := make(chan tea.Msg, 1)
	go func() { second <- watchV2Files(watcher, root)() }()
	select {
	case msg := <-second:
		t.Fatalf("rapid write burst produced a second reload message: %T", msg)
	case <-time.After(200 * time.Millisecond):
	}
}

func TestV2WatcherReportsMigrationAppearingAndCompleting(t *testing.T) {
	root := writeValidProject(t)
	watcher, err := newV2Watcher(root)
	if err != nil {
		t.Fatalf("newV2Watcher() error = %v", err)
	}
	t.Cleanup(func() { _ = watcher.Close() })

	operationID := createPendingOperation(t, root)
	awaitV2FileChange(t, watchV2Files(watcher, root))
	if loaded := loadProject(root); loaded.State.Next.Kind != data.NextPendingMigration {
		t.Fatalf("Next.Kind while migration %s is present = %q, want pending migration", operationID, loaded.State.Next.Kind)
	}

	if err := os.RemoveAll(filepath.Join(filepath.Dir(root), ".savepoint", ".migration", operationID)); err != nil {
		t.Fatalf("remove migration operation: %v", err)
	}
	awaitV2FileChange(t, watchV2Files(watcher, root))
	if loaded := loadProject(root); loaded.Failed() || loaded.State.Next.Kind == data.NextPendingMigration {
		t.Fatalf("load after migration completion = diagnostic %q, Next %q; want ordinary project state", loaded.Diagnostic, loaded.State.Next.Kind)
	}
}

func TestReloadWritesNothing(t *testing.T) {
	root := writeValidProject(t)
	model := openBoard(t, root, "")
	path := taskPath(root, "O001", "T001")
	testutil.WriteFile(t, path,
		"---\nid: T001\ntitle: \"Edited outside the board\"\nobjective: O001\nplanned_by: {role: planner, session: board-fixture}\nstatus: planned\n---\n")
	contentBefore, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	infoBefore, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	updated, _ := model.Update(loadCmd(root)().(projectLoadedMsg))
	_ = updated
	contentAfter, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	infoAfter, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(contentAfter) != string(contentBefore) || !infoAfter.ModTime().Equal(infoBefore.ModTime()) {
		t.Errorf("reload changed task bytes or mtime: bytes equal=%t, mtime equal=%t", string(contentAfter) == string(contentBefore), infoAfter.ModTime().Equal(infoBefore.ModTime()))
	}
}

func TestReloadFollowsTaskIdentityAcrossColumnsAndPreservesOverlay(t *testing.T) {
	root := writeValidProject(t)
	model := openBoard(t, root, "")
	model.Width = 120
	model.Height = 30
	detail, _ := newTaskDetail(model.State.Index, "T001")
	model.Detail = &detail
	model.DetailOffset = 1

	// Keep the same identity while moving the Task from Planned to In Progress.
	testutil.WriteFile(t, taskPath(root, "O001", "T001"),
		"---\nid: T001\ntitle: \"Do the thing now\"\nobjective: O001\nplanned_by: {role: planner, session: board-fixture}\nstatus: in_progress\nstage: build\n---\n\n# changed\n")
	updated, _ := model.Update(loadCmd(root)().(projectLoadedMsg))
	after := updated.(Model)
	if after.FocusedColumn != data.ColumnInProgress || after.FocusedCard != 0 {
		t.Errorf("focus = %s/%d, want T001's new In Progress position", after.FocusedColumn, after.FocusedCard)
	}
	if after.Detail == nil || after.Detail.ID != "T001" || after.Detail.Title != "Do the thing now" {
		t.Errorf("detail = %+v, want the surviving edited Task", after.Detail)
	}
	if after.DetailOffset != 1 {
		t.Errorf("DetailOffset = %d, want the prior scroll offset preserved", after.DetailOffset)
	}
	if after.State.Next.Task == nil || after.State.Next.Task.Title != "Do the thing now" {
		t.Errorf("Next.Task = %+v, want the edited Task", after.State.Next.Task)
	}
}

func TestFailedReloadKeepsLastGoodBoardAndRecovers(t *testing.T) {
	root := writeValidProject(t)
	model := openBoard(t, root, "")
	goodTitle := model.State.Index.Tasks["T001"].Title

	testutil.WriteFile(t, taskPath(root, "O001", "T001"),
		"---\nid: T001\nobjective: O001\nstatus: planned\n---\n")
	updated, _ := model.Update(loadCmd(root)().(projectLoadedMsg))
	after := updated.(Model)
	if after.State.Index.Tasks["T001"].Title != goodTitle {
		t.Errorf("last good title = %q, want %q", after.State.Index.Tasks["T001"].Title, goodTitle)
	}
	view := after.View()
	if after.ReloadDiagnostic == "" || !strings.Contains(view, "missing required field title") {
		t.Errorf("View() does not show the failed reload diagnostic:\n%s", view)
	}
	if got := strings.Count(view, "missing required field title"); got != 1 {
		t.Errorf("failed reload diagnostic appears %d times in View(); want exactly once:\n%s", got, view)
	}
	if strings.Contains(after.StatusMessage, "missing required field title") {
		t.Errorf("StatusMessage repeats the failed reload diagnostic: %q", after.StatusMessage)
	}

	writeTask(t, root, "O001", "T001", "Recovered title", "status: planned\n")
	recovered, _ := after.Update(loadCmd(root)().(projectLoadedMsg))
	final := recovered.(Model)
	if final.ReloadDiagnostic != "" || final.State.Index.Tasks["T001"].Title != "Recovered title" {
		t.Errorf("recovery state = diagnostic %q, title %q; want a clean successful reload", final.ReloadDiagnostic, final.State.Index.Tasks["T001"].Title)
	}
}

func TestReloadRouterEditUpdatesSelectionDiagnostic(t *testing.T) {
	root := writeValidProject(t)
	model := openBoard(t, root, "")
	writeRouter(t, root, "task", "O999", "T999")
	updated, _ := model.Update(loadCmd(root)().(projectLoadedMsg))
	after := updated.(Model)
	if after.SelectedObjective != "" {
		t.Errorf("SelectedObjective = %q, want no substitute for missing router Objective", after.SelectedObjective)
	}
	if after.State.Next.SelectionDiagnostic == nil {
		t.Fatal("Next.SelectionDiagnostic = nil, want the router edit reported")
	}
}

func TestReloadRefreshesCheckAndIssueSurfaces(t *testing.T) {
	checkRoot := writeEvidenceProject(t)
	checkModel := openTaskDetail(t, checkRoot, "T002")
	writeCheckExtra(t, checkRoot, "C006", "task", "T002", "NEEDS WORK", "supersedes: C003\n")
	updated, _ := checkModel.Update(loadCmd(checkRoot)().(projectLoadedMsg))
	refreshed := updated.(Model)
	if refreshed.Detail == nil || len(refreshed.Detail.Checks) != 3 || refreshed.Detail.Checks[2].Check.ID != "C006" {
		t.Errorf("reloaded Check history = %+v, want the new C006 entry", refreshed.Detail)
	}
	if refreshed.State.Next.Task == nil || refreshed.State.Next.Task.ID != "T002" {
		t.Errorf("reloaded Next.Task = %+v, want T002 after its Check changed", refreshed.State.Next.Task)
	}

	issueRoot := writeIssuesProject(t)
	issueModel := press(t, issueBoard(t, issueRoot), "i")
	writeBoardIssue(t, issueRoot, "I006", "A newly opened follow-up", "other", "open",
		"source: {kind: report, actor: {role: owner, session: issue-fixture}, at: '2026-01-08T00:00:00Z'}\n")
	updated, _ = issueModel.Update(loadCmd(issueRoot)().(projectLoadedMsg))
	refreshed = updated.(Model)
	if _, ok := refreshed.State.Issues.ByID["I006"]; !ok {
		t.Fatal("reloaded Issues catalog is missing newly written I006")
	}

	writeBoardIssue(t, issueRoot, "I001", "A repair is needed", "defect", "resolved",
		"source: {kind: check, check: C001, actor: {role: checker, session: issue-fixture}, at: '2026-01-01T00:00:00Z'}\n"+
			"tasks: [T001]\nchecks: [C001]\n"+
			"resolution: {disposition: accepted, actor: {role: owner, session: issue-fixture}, at: '2026-01-08T00:00:00Z', reason: 'accepted for the fixture'}\n")
	if loaded := loadProject(issueRoot); loaded.Failed() {
		t.Fatalf("resolved Issue edit refused: %s", loaded.Diagnostic)
	}
	updated, _ = refreshed.Update(loadCmd(issueRoot)().(projectLoadedMsg))
	refreshed = updated.(Model)
	if got := refreshed.State.Issues.ByID["I001"].Issue.Status; got != data.IssueStatusResolved {
		t.Errorf("reloaded I001 status = %q, want resolved", got)
	}
}

func awaitV2FileChange(t *testing.T, cmd tea.Cmd) tea.Msg {
	t.Helper()
	if cmd == nil {
		return nil
	}
	result := make(chan tea.Msg, 1)
	go func() { result <- cmd() }()
	select {
	case msg := <-result:
		return msg
	case <-time.After(2 * time.Second):
		t.Fatal("watch command did not emit a file-change message")
		return nil
	}
}
