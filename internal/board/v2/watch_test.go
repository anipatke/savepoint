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
		"objectives/O-001-fixture/Objective.md",
		"objectives/O-001-fixture/tasks/T-001-fixture.md",
		"checks/C-001.md",
		"issues/I-001.md",
		"releases/R-001-first/Release.md",
		"router.md",
		"config.yml",
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
		".migration/op/operation.yml",
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

	testutil.WriteFile(t, filepath.Join(root, "objectives", "O-002-fixture", "Objective.md"),
		"---\nid: O-002\ntitle: \"Second objective\"\nstatus: planned\n---\n")
	awaitV2FileChange(t, watchV2Files(watcher, root))

	// The event above registered the new Objective directory. A later edit to
	// that directory must be visible without rebuilding the watcher.
	testutil.WriteFile(t, filepath.Join(root, "objectives", "O-002-fixture", "Objective.md"),
		"---\nid: O-002\ntitle: \"Renamed objective\"\nstatus: planned\n---\n")
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

	// The burst is spaced 10 ms apart; a quiet interval well above that keeps
	// a stalled CI runner from splitting it into two reloads (I-075).
	const quiet = 500 * time.Millisecond
	result := make(chan tea.Msg, 1)
	go func() { result <- watchV2FilesAfter(watcher, root, quiet)() }()
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
	go func() { second <- watchV2FilesAfter(watcher, root, quiet)() }()
	select {
	case msg := <-second:
		t.Fatalf("rapid write burst produced a second reload message: %T", msg)
	case <-time.After(2 * quiet):
	}
}

func TestReloadWritesNothing(t *testing.T) {
	root := writeValidProject(t)
	model := openBoard(t, root, "")
	path := taskPath(root, "O-001", "T-001")
	testutil.WriteFile(t, path,
		"---\nid: T-001\ntitle: \"Edited outside the board\"\nobjective: O-001\nplanned_by: {role: planner, session: board-fixture}\nstatus: planned\n---\n")
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
	detail, _ := newTaskDetail(model.State.Index, "T-001")
	model.Detail = &detail
	model.DetailOffset = 1

	// Keep the same identity while moving the Task from Planned to In Progress.
	testutil.WriteFile(t, taskPath(root, "O-001", "T-001"),
		"---\nid: T-001\ntitle: \"Do the thing now\"\nobjective: O-001\nplanned_by: {role: planner, session: board-fixture}\nstatus: in_progress\nstage: build\n---\n\n# changed\n")
	updated, _ := model.Update(loadCmd(root)().(projectLoadedMsg))
	after := updated.(Model)
	if after.FocusedColumn != data.ColumnInProgress || after.FocusedCard != 0 {
		t.Errorf("focus = %s/%d, want T-001's new In Progress position", after.FocusedColumn, after.FocusedCard)
	}
	if after.Detail == nil || after.Detail.ID != "T-001" || after.Detail.Title != "Do the thing now" {
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
	goodTitle := model.State.Index.Tasks["T-001"].Title

	testutil.WriteFile(t, taskPath(root, "O-001", "T-001"),
		"---\nid: T-001\nobjective: O-001\nstatus: planned\n---\n")
	updated, _ := model.Update(loadCmd(root)().(projectLoadedMsg))
	if updated.(Model).ReloadDiagnostic != "" {
		t.Fatal("first failed reload shows RELOAD; want one quiet retry first")
	}
	updated, _ = updated.Update(retriedLoad(root))
	after := updated.(Model)
	if after.State.Index.Tasks["T-001"].Title != goodTitle {
		t.Errorf("last good title = %q, want %q", after.State.Index.Tasks["T-001"].Title, goodTitle)
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

	writeTask(t, root, "O-001", "T-001", "Recovered title", "status: planned\n")
	recovered, _ := after.Update(loadCmd(root)().(projectLoadedMsg))
	final := recovered.(Model)
	if final.ReloadDiagnostic != "" || final.State.Index.Tasks["T-001"].Title != "Recovered title" {
		t.Errorf("recovery state = diagnostic %q, title %q; want a clean successful reload", final.ReloadDiagnostic, final.State.Index.Tasks["T-001"].Title)
	}
}

func TestFailedReloadRetriesBeforeReporting(t *testing.T) {
	root := writeValidProject(t)
	model := openBoard(t, root, "")

	// An agent's multi-step edit is read half-finished, then completed
	// before the retry.
	testutil.WriteFile(t, taskPath(root, "O-001", "T-001"),
		"---\nid: T-001\nobjective: O-001\nstatus: planned\n---\n")
	updated, cmd := model.Update(loadCmd(root)().(projectLoadedMsg))
	if cmd == nil || updated.(Model).ReloadDiagnostic != "" {
		t.Fatalf("failed reload = diagnostic %q, cmd %v; want no RELOAD and a retry", updated.(Model).ReloadDiagnostic, cmd)
	}
	writeTask(t, root, "O-001", "T-001", "Finished edit", "status: planned\n")
	retried, ok := cmd().(projectLoadedMsg)
	if !ok || !retried.Retry {
		t.Fatalf("retry cmd returned %#v, want a Retry projectLoadedMsg", retried)
	}
	updated, _ = updated.Update(retried)
	after := updated.(Model)
	if after.ReloadDiagnostic != "" || after.State.Index.Tasks["T-001"].Title != "Finished edit" {
		t.Errorf("after retry = diagnostic %q, title %q; want the finished edit and no RELOAD", after.ReloadDiagnostic, after.State.Index.Tasks["T-001"].Title)
	}
}

func TestOlderLoadResultDoesNotOverwriteNewer(t *testing.T) {
	root := writeValidProject(t)
	model := openBoard(t, root, "")

	older := loadCmd(root)().(projectLoadedMsg)
	writeTask(t, root, "O-001", "T-001", "Newer title", "status: planned\n")
	newer := loadCmd(root)().(projectLoadedMsg)

	updated, _ := model.Update(newer)
	updated, _ = updated.Update(older)
	if got := updated.(Model).State.Index.Tasks["T-001"].Title; got != "Newer title" {
		t.Errorf("title = %q after a stale load arrived last; want %q", got, "Newer title")
	}
}

func TestReloadRouterEditUpdatesSelectionDiagnostic(t *testing.T) {
	root := writeValidProject(t)
	model := openBoard(t, root, "")
	writeRouter(t, root, "task", "O-999", "T-999")
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
	checkModel := openTaskDetail(t, checkRoot, "T-002")
	writeCheckExtra(t, checkRoot, "C-006", "task", "T-002", "NEEDS WORK", "supersedes: C-003\n")
	updated, _ := checkModel.Update(loadCmd(checkRoot)().(projectLoadedMsg))
	refreshed := updated.(Model)
	if refreshed.Detail == nil || len(refreshed.Detail.Checks) != 3 || refreshed.Detail.Checks[2].Check.ID != "C-006" {
		t.Errorf("reloaded Check history = %+v, want the new C-006 entry", refreshed.Detail)
	}
	if refreshed.State.Next.Task == nil || refreshed.State.Next.Task.ID != "T-002" {
		t.Errorf("reloaded Next.Task = %+v, want T-002 after its Check changed", refreshed.State.Next.Task)
	}

	issueRoot := writeIssuesProject(t)
	issueModel := press(t, issueBoard(t, issueRoot), "i")
	writeBoardIssue(t, issueRoot, "I-006", "A newly opened follow-up", "other", "open",
		"source: {kind: report, actor: {role: owner, session: issue-fixture}, at: '2026-01-08T00:00:00Z'}\n")
	updated, _ = issueModel.Update(loadCmd(issueRoot)().(projectLoadedMsg))
	refreshed = updated.(Model)
	if _, ok := refreshed.State.Issues.ByID["I-006"]; !ok {
		t.Fatal("reloaded Issues catalog is missing newly written I-006")
	}

	writeBoardIssue(t, issueRoot, "I-001", "A repair is needed", "defect", "resolved",
		"source: {kind: check, check: C-001, actor: {role: checker, session: issue-fixture}, at: '2026-01-01T00:00:00Z'}\n"+
			"tasks: [T-001]\nchecks: [C-001]\n"+
			"resolution: {disposition: accepted, actor: {role: owner, session: issue-fixture}, at: '2026-01-08T00:00:00Z', reason: 'accepted for the fixture'}\n")
	if loaded := loadProject(issueRoot); loaded.Failed() {
		t.Fatalf("resolved Issue edit refused: %s", loaded.Diagnostic)
	}
	updated, _ = refreshed.Update(loadCmd(issueRoot)().(projectLoadedMsg))
	refreshed = updated.(Model)
	if got := refreshed.State.Issues.ByID["I-001"].Issue.Status; got != data.IssueStatusResolved {
		t.Errorf("reloaded I-001 status = %q, want resolved", got)
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
