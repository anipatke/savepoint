package data

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/testutil"
)

func TestFindSavepointRoot(t *testing.T) {
	d := NewDiscover()
	savepointRoot := createDiscoveryFixture(t)
	start := filepath.Join(filepath.Dir(savepointRoot), "nested", "child")
	testutil.MkdirAll(t, start)

	root, err := d.FindSavepointRoot(start)
	if err != nil {
		t.Fatalf("FindSavepointRoot() error = %v", err)
	}
	if root != savepointRoot {
		t.Errorf("FindSavepointRoot() = %v, want %v", root, savepointRoot)
	}
}

func TestListReleases(t *testing.T) {
	d := NewDiscover()
	root := createDiscoveryFixture(t)

	releases, err := d.ListReleases(root)
	if err != nil {
		t.Fatalf("ListReleases() error = %v", err)
	}

	if len(releases) != 2 {
		t.Fatalf("ListReleases() returned %d releases, want 2", len(releases))
	}
	if releases[0].ID != "v1" || releases[1].ID != "v2" {
		t.Errorf("ListReleases() IDs = %v, want [v1 v2]", []string{releases[0].ID, releases[1].ID})
	}
}

func TestListRootDirs(t *testing.T) {
	d := NewDiscover()
	root := t.TempDir()
	testutil.MkdirAll(t, filepath.Join(root, "beta"))
	testutil.MkdirAll(t, filepath.Join(root, "alpha"))
	testutil.WriteFile(t, filepath.Join(root, "notes.txt"), "test")

	dirs, err := d.ListRootDirs(root)
	if err != nil {
		t.Fatalf("ListRootDirs() error = %v", err)
	}

	if len(dirs) != 2 || dirs[0] != "alpha" || dirs[1] != "beta" {
		t.Fatalf("ListRootDirs() = %v, want [alpha beta]", dirs)
	}
}

func TestListRootDirsRejectsFile(t *testing.T) {
	d := NewDiscover()
	root := t.TempDir()
	path := filepath.Join(root, "not-dir")
	testutil.WriteFile(t, path, "test")

	_, err := d.ListRootDirs(path)
	if err == nil {
		t.Fatal("ListRootDirs() error = nil, want not directory error")
	}
}

func TestListEpics(t *testing.T) {
	d := NewDiscover()
	root := createDiscoveryFixture(t)

	epics, err := d.ListEpics(root, "v1")
	if err != nil {
		t.Fatalf("ListEpics() error = %v", err)
	}

	if len(epics) != 2 {
		t.Fatalf("ListEpics() returned %d epics, want 2", len(epics))
	}
	if epics[0].ID != "E01-go-setup" || epics[1].ID != "E02-data-readers" {
		t.Errorf("ListEpics() IDs = %v, want [E01-go-setup E02-data-readers]", []string{epics[0].ID, epics[1].ID})
	}
}

func TestListTasks(t *testing.T) {
	d := NewDiscover()
	root := createDiscoveryFixture(t)

	tasks, err := d.ListTasks(root, "v1", "E02-data-readers")
	if err != nil {
		t.Fatalf("ListTasks() error = %v", err)
	}

	if len(tasks) != 2 {
		t.Fatalf("ListTasks() returned %d tasks, want 2", len(tasks))
	}
	if tasks[0].ID != "T001-task-struct" || tasks[1].ID != "T002-frontmatter-parser" {
		t.Errorf("ListTasks() IDs = %v, want [T001-task-struct T002-frontmatter-parser]", []string{tasks[0].ID, tasks[1].ID})
	}
}

func createDiscoveryFixture(t *testing.T) string {
	t.Helper()

	root := t.TempDir()
	savepointRoot := filepath.Join(root, ".savepoint")
	paths := []string{
		filepath.Join(savepointRoot, "releases", "v1", "epics", "E02-data-readers", "tasks"),
		filepath.Join(savepointRoot, "releases", "v1", "epics", "E01-go-setup", "tasks"),
		filepath.Join(savepointRoot, "releases", "v1", "epics", "_archived"),
		filepath.Join(savepointRoot, "releases", "v2", "epics"),
	}
	for _, path := range paths {
		testutil.MkdirAll(t, path)
	}

	files := []string{
		filepath.Join(savepointRoot, "releases", "v1", "epics", "E02-data-readers", "tasks", "T002-frontmatter-parser.md"),
		filepath.Join(savepointRoot, "releases", "v1", "epics", "E02-data-readers", "tasks", "T001-task-struct.md"),
		filepath.Join(savepointRoot, "releases", "v1", "epics", "E02-data-readers", "tasks", "notes.txt"),
	}
	for _, file := range files {
		testutil.WriteFile(t, file, "test")
	}

	return savepointRoot
}

func writeV2ObjectiveFixture(t *testing.T, root, dirName, id, title string, dependsOn ...string) {
	t.Helper()
	content := "---\nid: " + id + "\ntitle: \"" + title + "\"\nstatus: planned\n"
	if len(dependsOn) > 0 {
		content += "depends_on: [" + joinIDs(dependsOn) + "]\n"
	}
	content += "---\n\n# " + title + "\n"
	testutil.WriteFile(t, filepath.Join(root, v2ObjectivesDirName, dirName, v2ObjectiveFileName), content)
}

func writeV2TaskFixture(t *testing.T, root, objDirName, fileName, id, title, objective string) {
	t.Helper()
	content := "---\nid: " + id + "\ntitle: \"" + title + "\"\nobjective: " + objective + "\nstatus: planned\n---\n\n# " + title + "\n"
	testutil.WriteFile(t, filepath.Join(root, v2ObjectivesDirName, objDirName, v2TasksDirName, fileName), content)
}

func joinIDs(ids []string) string {
	out := ""
	for i, id := range ids {
		if i > 0 {
			out += ", "
		}
		out += id
	}
	return out
}

func TestDiscoverV2Records_valid(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O001-first", "O001", "First objective")
	writeV2ObjectiveFixture(t, root, "O002-second", "O002", "Second objective", "O001")
	writeV2TaskFixture(t, root, "O001-first", "T001-alpha.md", "T001", "Alpha", "O001")
	writeV2TaskFixture(t, root, "O001-first", "T002-beta.md", "T002", "Beta", "O001")

	objectives, tasks, err := DiscoverV2Records(root)
	if err != nil {
		t.Fatalf("DiscoverV2Records() error = %v", err)
	}

	if len(objectives) != 2 {
		t.Fatalf("DiscoverV2Records() objectives = %d, want 2", len(objectives))
	}
	if objectives["O001"] == nil || objectives["O001"].Title != "First objective" {
		t.Errorf("objectives[O001] = %+v, want title %q", objectives["O001"], "First objective")
	}
	if objectives["O002"] == nil || len(objectives["O002"].DependsOn) != 1 || objectives["O002"].DependsOn[0] != "O001" {
		t.Errorf("objectives[O002] = %+v, want DependsOn [O001]", objectives["O002"])
	}

	if len(tasks) != 2 {
		t.Fatalf("DiscoverV2Records() tasks = %d, want 2", len(tasks))
	}
	if tasks["T001"] == nil || tasks["T001"].Objective != "O001" {
		t.Errorf("tasks[T001] = %+v, want objective O001", tasks["T001"])
	}
}

func TestDiscoverV2Records_emptyProjectHasNoObjectivesDir(t *testing.T) {
	root := t.TempDir()

	objectives, tasks, err := DiscoverV2Records(root)
	if err != nil {
		t.Fatalf("DiscoverV2Records() error = %v, want nil for a project with no Objectives yet", err)
	}
	if len(objectives) != 0 || len(tasks) != 0 {
		t.Errorf("DiscoverV2Records() = %v, %v, want both empty", objectives, tasks)
	}
}

// TestLoadV2Index_missingObjectiveStillDiscoversTasks proves that a task is
// not silently dropped when its containing Objective.md is absent. The index
// must retain the task long enough to return the named missing-owner error.
func TestLoadV2Index_missingObjectiveStillDiscoversTasks(t *testing.T) {
	root := t.TempDir()
	writeV2TaskFixture(t, root, "O999-missing", "T003-orphan.md", "T003", "Orphan task", "O999")

	_, err := LoadV2Index(root)
	if !errors.Is(err, ErrV2MissingOwner) {
		t.Fatalf("LoadV2Index() error = %v, want ErrV2MissingOwner", err)
	}
	if !strings.Contains(err.Error(), "T003") || !strings.Contains(err.Error(), "O999") {
		t.Fatalf("LoadV2Index() error = %v, want task and owner context", err)
	}
}

// TestDiscoverV2Records_movedTaskRetainsPathIndependentOwnership proves that
// a Task filed under a different Objective's tasks/ directory than the one
// it declares ownership to is still indexed by its own ID with its declared
// objective field intact: ownership is never inferred from the containing
// directory.
func TestDiscoverV2Records_movedTaskRetainsPathIndependentOwnership(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O001-first", "O001", "First objective")
	writeV2ObjectiveFixture(t, root, "O002-second", "O002", "Second objective")
	// T001 lives under O002's directory but still declares O001 as owner.
	writeV2TaskFixture(t, root, "O002-second", "T001-alpha.md", "T001", "Alpha", "O001")

	_, tasks, err := DiscoverV2Records(root)
	if err != nil {
		t.Fatalf("DiscoverV2Records() error = %v", err)
	}
	if tasks["T001"] == nil {
		t.Fatal("DiscoverV2Records() did not index the moved task")
	}
	if tasks["T001"].Objective != "O001" {
		t.Errorf("moved task Objective = %q, want O001 (declared owner, not containing directory)", tasks["T001"].Objective)
	}
}

func TestDiscoverV2Records_objectiveDirectoryNameMismatch(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O002-wrong-dir", "O001", "First objective")

	_, _, err := DiscoverV2Records(root)
	if !errors.Is(err, ErrV2PathMismatch) {
		t.Fatalf("DiscoverV2Records() error = %v, want ErrV2PathMismatch", err)
	}
}

func TestDiscoverV2Records_taskFileNameMismatch(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O001-first", "O001", "First objective")
	writeV2TaskFixture(t, root, "O001-first", "T002-wrong-name.md", "T001", "Alpha", "O001")

	_, _, err := DiscoverV2Records(root)
	if !errors.Is(err, ErrV2PathMismatch) {
		t.Fatalf("DiscoverV2Records() error = %v, want ErrV2PathMismatch", err)
	}
}

func TestDiscoverV2Records_duplicateObjectiveID(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O001-first", "O001", "First objective")
	writeV2ObjectiveFixture(t, root, "O001-duplicate", "O001", "Duplicate objective")

	_, _, err := DiscoverV2Records(root)
	if !errors.Is(err, ErrV2DuplicateID) {
		t.Fatalf("DiscoverV2Records() error = %v, want ErrV2DuplicateID", err)
	}
}

func TestDiscoverV2Records_duplicateTaskID(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O001-first", "O001", "First objective")
	writeV2ObjectiveFixture(t, root, "O002-second", "O002", "Second objective")
	writeV2TaskFixture(t, root, "O001-first", "T001-alpha.md", "T001", "Alpha", "O001")
	writeV2TaskFixture(t, root, "O002-second", "T001-alpha-again.md", "T001", "Alpha again", "O002")

	_, _, err := DiscoverV2Records(root)
	if !errors.Is(err, ErrV2DuplicateID) {
		t.Fatalf("DiscoverV2Records() error = %v, want ErrV2DuplicateID", err)
	}
}

func TestDiscoverV2Records_rejectsSymlinkEscape(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink creation requires elevated privileges on windows")
	}

	root := t.TempDir()
	outside := t.TempDir()
	writeV2ObjectiveFixture(t, outside, "O001-first", "O001", "Escaped objective")

	testutil.MkdirAll(t, filepath.Join(root, v2ObjectivesDirName))
	target := filepath.Join(outside, v2ObjectivesDirName, "O001-first")
	link := filepath.Join(root, v2ObjectivesDirName, "O001-first")
	if err := os.Symlink(target, link); err != nil {
		t.Fatalf("os.Symlink() error = %v", err)
	}

	_, _, err := DiscoverV2Records(root)
	if !errors.Is(err, ErrV2UnsafePath) {
		t.Fatalf("DiscoverV2Records() error = %v, want ErrV2UnsafePath", err)
	}
}

func TestDiscoverV2Records_rejectsCaseCollision(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O001-first", "O001", "First objective")
	writeV2ObjectiveFixture(t, root, "o001-First", "O002", "Case-colliding objective")

	_, _, err := DiscoverV2Records(root)
	if !errors.Is(err, ErrV2UnsafePath) {
		t.Fatalf("DiscoverV2Records() error = %v, want ErrV2UnsafePath", err)
	}
}

func writeV2CheckFixture(t *testing.T, root, fileName, id, scopeKind, scopeID string) {
	t.Helper()
	content := "---\nid: " + id + "\nscope: {kind: " + scopeKind + ", id: " + scopeID + "}\nresult: CLEAR\nchecked_by: {role: checker, session: sess-1}\nchecked_at: '2026-09-14T00:00:00Z'\n---\n\n# Check\n"
	testutil.WriteFile(t, filepath.Join(root, v2ChecksDirName, fileName), content)
}

func TestDiscoverV2Checks_valid(t *testing.T) {
	root := t.TempDir()
	writeV2CheckFixture(t, root, "C001-alpha.md", "C001", "task", "T001")
	writeV2CheckFixture(t, root, "C002-beta.md", "C002", "objective", "O001")

	checks, err := DiscoverV2Checks(root)
	if err != nil {
		t.Fatalf("DiscoverV2Checks() error = %v", err)
	}
	if len(checks) != 2 {
		t.Fatalf("DiscoverV2Checks() = %d checks, want 2", len(checks))
	}
	if checks["C001"] == nil || checks["C001"].Scope.ID != "T001" {
		t.Errorf("checks[C001] = %+v, want scope id T001", checks["C001"])
	}
	if checks["C002"] == nil || checks["C002"].Scope.Kind != CheckScopeObjective {
		t.Errorf("checks[C002] = %+v, want objective scope", checks["C002"])
	}
}

func TestDiscoverV2Checks_absentChecksDir(t *testing.T) {
	root := t.TempDir()

	checks, err := DiscoverV2Checks(root)
	if err != nil {
		t.Fatalf("DiscoverV2Checks() error = %v, want nil for a project with no checks/ yet", err)
	}
	if len(checks) != 0 {
		t.Errorf("DiscoverV2Checks() = %v, want empty", checks)
	}
}

func TestDiscoverV2Checks_fileNameMismatch(t *testing.T) {
	root := t.TempDir()
	writeV2CheckFixture(t, root, "C002-wrong-name.md", "C001", "task", "T001")

	_, err := DiscoverV2Checks(root)
	if !errors.Is(err, ErrV2PathMismatch) {
		t.Fatalf("DiscoverV2Checks() error = %v, want ErrV2PathMismatch", err)
	}
}

func TestDiscoverV2Checks_duplicateID(t *testing.T) {
	root := t.TempDir()
	writeV2CheckFixture(t, root, "C001-alpha.md", "C001", "task", "T001")
	writeV2CheckFixture(t, root, "C001-alpha-again.md", "C001", "task", "T001")

	_, err := DiscoverV2Checks(root)
	if !errors.Is(err, ErrV2DuplicateID) {
		t.Fatalf("DiscoverV2Checks() error = %v, want ErrV2DuplicateID", err)
	}
}

func TestDiscoverV2Checks_rejectsSymlinkEscape(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink creation requires elevated privileges on windows")
	}

	root := t.TempDir()
	outside := t.TempDir()
	outsideFile := filepath.Join(outside, "C001-alpha.md")
	testutil.WriteFile(t, outsideFile, "---\nid: C001\nscope: {kind: task, id: T001}\nresult: CLEAR\nchecked_by: {role: checker, session: sess-1}\nchecked_at: '2026-09-14T00:00:00Z'\n---\n\n# Check\n")

	testutil.MkdirAll(t, filepath.Join(root, v2ChecksDirName))
	link := filepath.Join(root, v2ChecksDirName, "C001-alpha.md")
	if err := os.Symlink(outsideFile, link); err != nil {
		t.Fatalf("os.Symlink() error = %v", err)
	}

	_, err := DiscoverV2Checks(root)
	if !errors.Is(err, ErrV2UnsafePath) {
		t.Fatalf("DiscoverV2Checks() error = %v, want ErrV2UnsafePath", err)
	}
}

func TestDiscoverV2Checks_rejectsCaseCollision(t *testing.T) {
	root := t.TempDir()
	writeV2CheckFixture(t, root, "C001-alpha.md", "C001", "task", "T001")
	writeV2CheckFixture(t, root, "c001-Alpha.md", "C002", "task", "T001")

	_, err := DiscoverV2Checks(root)
	if !errors.Is(err, ErrV2UnsafePath) {
		t.Fatalf("DiscoverV2Checks() error = %v, want ErrV2UnsafePath", err)
	}
}

func writeV2IssueFixture(t *testing.T, root, fileName, id, status string) {
	t.Helper()
	content := "---\nid: " + id + "\ntitle: \"Follow-up\"\ntype: defect\nstatus: " + status +
		"\nsource: {kind: report, actor: {role: owner, session: owner-1}, at: '2026-09-15T00:00:00Z'}\n---\n\n# Issue\n"
	testutil.WriteFile(t, filepath.Join(root, v2IssuesDirName, fileName), content)
}

func TestDiscoverV2Issues_valid(t *testing.T) {
	root := t.TempDir()
	writeV2IssueFixture(t, root, "I001-alpha.md", "I001", "open")
	writeV2IssueFixture(t, root, "I002-beta.md", "I002", "in_progress")

	issues, err := DiscoverV2Issues(root)
	if err != nil {
		t.Fatalf("DiscoverV2Issues() error = %v", err)
	}
	if len(issues) != 2 {
		t.Fatalf("DiscoverV2Issues() = %d issues, want 2", len(issues))
	}
	if issues["I001"] == nil || issues["I001"].Status != IssueStatusOpen {
		t.Errorf("issues[I001] = %+v, want status open", issues["I001"])
	}
	if issues["I002"] == nil || issues["I002"].Source.Path != filepath.Join(v2IssuesDirName, "I002-beta.md") {
		t.Errorf("issues[I002] = %+v, want its project-relative source path", issues["I002"])
	}
}

func TestDiscoverV2Issues_absentIssuesDir(t *testing.T) {
	root := t.TempDir()

	issues, err := DiscoverV2Issues(root)
	if err != nil {
		t.Fatalf("DiscoverV2Issues() error = %v, want nil for a project with no issues/ yet", err)
	}
	if len(issues) != 0 {
		t.Errorf("DiscoverV2Issues() = %v, want empty", issues)
	}
}

func TestDiscoverV2Issues_fileNameMismatch(t *testing.T) {
	root := t.TempDir()
	writeV2IssueFixture(t, root, "I002-wrong-name.md", "I001", "open")

	_, err := DiscoverV2Issues(root)
	if !errors.Is(err, ErrV2PathMismatch) {
		t.Fatalf("DiscoverV2Issues() error = %v, want ErrV2PathMismatch", err)
	}
}

func TestDiscoverV2Issues_duplicateID(t *testing.T) {
	root := t.TempDir()
	writeV2IssueFixture(t, root, "I001-alpha.md", "I001", "open")
	writeV2IssueFixture(t, root, "I001-alpha-again.md", "I001", "open")

	_, err := DiscoverV2Issues(root)
	if !errors.Is(err, ErrV2DuplicateID) {
		t.Fatalf("DiscoverV2Issues() error = %v, want ErrV2DuplicateID", err)
	}
}

// TestDiscoverV2Issues_rejectsTraversalFilename proves a filename that tries
// to climb out of issues/ is never treated as a record: discovery walks real
// directory entries, so the escape attempt is simply not an Issue file.
func TestDiscoverV2Issues_rejectsTraversalFilename(t *testing.T) {
	root := t.TempDir()
	writeV2IssueFixture(t, root, "I001-alpha.md", "I001", "open")
	testutil.WriteFile(t, filepath.Join(root, "I999-escaped.md"),
		"---\nid: I999\ntitle: \"Escaped\"\ntype: defect\nstatus: open\nsource: {kind: report, actor: {role: owner, session: o}, at: '2026-09-15T00:00:00Z'}\n---\n\n# Issue\n")

	issues, err := DiscoverV2Issues(root)
	if err != nil {
		t.Fatalf("DiscoverV2Issues() error = %v", err)
	}
	if _, ok := issues["I999"]; ok {
		t.Error("DiscoverV2Issues() indexed a record outside issues/, want confinement to the issues directory")
	}
}

func TestDiscoverV2Issues_rejectsSymlinkEscape(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink creation requires elevated privileges on windows")
	}

	root := t.TempDir()
	outside := t.TempDir()
	outsideFile := filepath.Join(outside, "I001-alpha.md")
	testutil.WriteFile(t, outsideFile,
		"---\nid: I001\ntitle: \"Escaped\"\ntype: defect\nstatus: open\nsource: {kind: report, actor: {role: owner, session: o}, at: '2026-09-15T00:00:00Z'}\n---\n\n# Issue\n")

	testutil.MkdirAll(t, filepath.Join(root, v2IssuesDirName))
	link := filepath.Join(root, v2IssuesDirName, "I001-alpha.md")
	if err := os.Symlink(outsideFile, link); err != nil {
		t.Fatalf("os.Symlink() error = %v", err)
	}

	_, err := DiscoverV2Issues(root)
	if !errors.Is(err, ErrV2UnsafePath) {
		t.Fatalf("DiscoverV2Issues() error = %v, want ErrV2UnsafePath", err)
	}
}

func TestDiscoverV2Issues_rejectsCaseCollision(t *testing.T) {
	root := t.TempDir()
	writeV2IssueFixture(t, root, "I001-alpha.md", "I001", "open")
	writeV2IssueFixture(t, root, "i001-Alpha.md", "I002", "open")

	_, err := DiscoverV2Issues(root)
	if !errors.Is(err, ErrV2UnsafePath) {
		t.Fatalf("DiscoverV2Issues() error = %v, want ErrV2UnsafePath", err)
	}
}

// TestDiscoverV2Issues_malformedRecordFailsClosed proves a structurally bad
// Issue stops the load rather than being skipped, so a project never reports
// fewer Issues than it actually has.
func TestDiscoverV2Issues_malformedRecordFailsClosed(t *testing.T) {
	root := t.TempDir()
	writeV2IssueFixture(t, root, "I001-alpha.md", "I001", "open")
	writeV2IssueFixture(t, root, "I002-beta.md", "I002", "done")

	_, err := DiscoverV2Issues(root)
	if !errors.Is(err, ErrV2InvalidLifecycle) {
		t.Fatalf("DiscoverV2Issues() error = %v, want ErrV2InvalidLifecycle", err)
	}
}
