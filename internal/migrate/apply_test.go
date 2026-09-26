package migrate

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/data"
)

func copyFixtureProject(t *testing.T, fixture string) string {
	t.Helper()
	src := fixtureProjectRoot(fixture)
	dst := t.TempDir()
	if err := filepath.WalkDir(src, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if entry.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, content, 0644)
	}); err != nil {
		t.Fatalf("copy fixture %s: %v", fixture, err)
	}
	return dst
}

func mustLoadV2Index(t *testing.T, root string) *data.V2Index {
	t.Helper()
	if err := data.CheckRuntimeSchema(root); err != nil {
		t.Fatalf("CheckRuntimeSchema(%s) error = %v", root, err)
	}
	index, err := data.LoadV2Index(filepath.Join(root, ".savepoint"))
	if err != nil {
		t.Fatalf("LoadV2Index(%s) error = %v", root, err)
	}
	return index
}

func mustExist(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected %s to exist: %v", path, err)
	}
}

func mustNotExist(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Lstat(path); err == nil {
		t.Fatalf("expected %s to be gone, but it still exists", path)
	} else if !os.IsNotExist(err) {
		t.Fatalf("stat %s: %v", path, err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(content)
}

func TestApply_v1BasicEndToEnd(t *testing.T) {
	root := copyFixtureProject(t, "v1-basic")
	plan := mustPlan(t, root)
	if !plan.Appliable {
		t.Fatalf("plan.Appliable = false, want true; ambiguities = %+v", plan.Ambiguities)
	}

	result, err := Apply(root, plan)
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	if result.AlreadyMigrated {
		t.Fatalf("Apply() result = %+v, want a new migration", result)
	}
	if len(result.WrittenPaths) == 0 || result.WrittenPaths[len(result.WrittenPaths)-1] != ".savepoint/config.yml" {
		t.Fatalf("written paths end at %v, want schema config last", result.WrittenPaths)
	}

	version, err := data.ReadSchemaVersion(filepath.Join(root, ".savepoint", "config.yml"))
	if err != nil || version != data.SchemaVersionV2 {
		t.Fatalf("ReadSchemaVersion() = %v, %v, want SchemaVersionV2", version, err)
	}
	if config := readFile(t, filepath.Join(root, ".savepoint", "config.yml")); !strings.Contains(config, `bg: "#000000"`) {
		t.Errorf("config.yml lost its unrelated theme key:\n%s", config)
	}

	mustNotExist(t, filepath.Join(root, ".savepoint", "releases", "v1", "epics", "E01-example", "E01-Detail.md"))
	mustNotExist(t, filepath.Join(root, ".savepoint", "releases", "v1", "epics", "E01-example", "tasks", "T001-original.md"))
	activeTaskSource := filepath.Join(root, ".savepoint", "releases", "v1", "epics", "E01-example", "tasks", "T002-follow-up.md")
	mustNotExist(t, activeTaskSource)
	mustExist(t, filepath.Join(root, ".savepoint", "archive", "v1", ".savepoint", "releases", "v1", "epics", "E01-example", "tasks", "T002-follow-up.md"))
	mustNotExist(t, filepath.Join(root, ".savepoint", "PRD.md"))
	mustNotExist(t, filepath.Join(root, ".savepoint", ".migration"))

	archived, ok := archiveByPath(plan, ".savepoint/releases/v1/epics/E01-example/tasks/T001-original.md")
	if !ok {
		t.Fatal("T001-original was not archived by the plan")
	}
	archiveContent := readFile(t, filepath.Join(root, filepath.FromSlash(archived.ArchivePath)))
	originalContent := readFile(t, filepath.Join(fixtureProjectRoot("v1-basic"), filepath.FromSlash(archived.SourcePath)))
	if archiveContent != originalContent {
		t.Error("archived T001-original does not match its original source bytes")
	}
	mustExist(t, filepath.Join(root, ".savepoint", "Idea.md"))
	mustExist(t, filepath.Join(root, ".savepoint", "migrations", "v1-to-v2.yml"))

	originalDesign := readFile(t, filepath.Join(fixtureProjectRoot("v1-basic"), ".savepoint", "Design.md"))
	if got := readFile(t, filepath.Join(root, ".savepoint", "Design.md")); got != originalDesign {
		t.Error("Design.md changed; preserved-in-place documents must stay byte-identical")
	}

	index := mustLoadV2Index(t, root)
	if len(index.Objectives) != 1 || len(index.Tasks) != 1 {
		t.Errorf("migrated index = %d Objectives, %d Tasks; want 1, 1", len(index.Objectives), len(index.Tasks))
	}
	mustNotExist(t, filepath.Join(root, ".savepoint", "archive", "v1", "archive"))
}

func TestApply_v1HistoryKeepsReleaseScopedIDs(t *testing.T) {
	root := copyFixtureProject(t, "v1-history")
	plan := mustPlan(t, root)
	if _, err := Apply(root, plan); err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	index := mustLoadV2Index(t, root)
	if len(index.Tasks) != 2 {
		t.Fatalf("migrated index = %d Tasks; want two active Tasks", len(index.Tasks))
	}
	seen := map[string]bool{}
	for id := range index.Tasks {
		if seen[id] {
			t.Errorf("global Task ID %s allocated more than once", id)
		}
		seen[id] = true
	}
	mustNotExist(t, filepath.Join(root, ".savepoint", "releases", "v1", "epics", "E01-example", "tasks", "T001-shared.md"))
	mustNotExist(t, filepath.Join(root, ".savepoint", "releases", "v1.1", "epics", "E01-example", "E01-Detail.md"))
	mustNotExist(t, filepath.Join(root, ".savepoint", "releases", "v1.1", "epics", "E01-example", "tasks", "T001-shared.md"))
	mustExist(t, filepath.Join(root, ".savepoint", "archive", "v1", ".savepoint", "releases", "v1.1", "epics", "E01-example", "tasks", "T001-shared.md"))
	mustExist(t, filepath.Join(root, ".savepoint", "archive", "v1", ".savepoint", "releases", "v1", "epics", "E01-example", "tasks", "T001-shared.md"))
	mustExist(t, filepath.Join(root, ".savepoint", "archive", "v1", ".savepoint", "releases", "v1.1", "epics", "E01-example", "E01-Detail.md"))
}

func TestApply_missingConfigCreatesSchemaConfigAsLastWrite(t *testing.T) {
	root := copyFixtureProject(t, "v1-basic")
	configPath := filepath.Join(root, ".savepoint", "config.yml")
	if err := os.Remove(configPath); err != nil {
		t.Fatalf("remove fixture config: %v", err)
	}
	plan := mustPlan(t, root)
	if planHasSource(plan, ".savepoint/config.yml") {
		t.Fatal("plan records a config source that was removed before planning")
	}

	result, err := Apply(root, plan)
	if err != nil {
		t.Fatalf("Apply() without a source config error = %v", err)
	}
	if got := result.WrittenPaths[len(result.WrittenPaths)-1]; got != ".savepoint/config.yml" {
		t.Fatalf("last written path = %q, want schema config activation", got)
	}
	if got := readFile(t, configPath); !strings.Contains(got, "schema_version: 2") {
		t.Fatalf("created config.yml = %q, want schema_version: 2", got)
	}
	tracked, untracked := gitUndoPathGroups(plan)
	if containsPath(tracked, ".savepoint/config.yml") || !containsPath(untracked, ".savepoint/config.yml") {
		t.Fatalf("undo groups for absent source config = tracked %v, untracked %v", tracked, untracked)
	}
}

func containsPath(paths []string, want string) bool {
	for _, path := range paths {
		if path == want {
			return true
		}
	}
	return false
}

func TestApply_notAppliableWritesNothing(t *testing.T) {
	root := copyFixtureProject(t, "v1-basic")
	taskPath := filepath.Join(root, ".savepoint", "releases", "v1", "epics", "E01-example", "tasks", "T002-follow-up.md")
	content := readFile(t, taskPath)
	if !strings.Contains(content, "status: in_progress") {
		t.Fatal("fixture does not contain the expected active status")
	}
	if err := os.WriteFile(taskPath, []byte(strings.Replace(content, "status: in_progress", "status: sideways", 1)), 0644); err != nil {
		t.Fatalf("corrupt fixture copy: %v", err)
	}
	plan := mustPlan(t, root)
	if plan.Appliable {
		t.Fatal("plan.Appliable = true, want an unrecognized status to block")
	}
	before := snapshotTree(t, root)
	_, err := Apply(root, plan)
	if !errors.Is(err, ErrPlanNotAppliable) {
		t.Fatalf("Apply() error = %v, want ErrPlanNotAppliable", err)
	}
	assertSnapshotsEqual(t, before, snapshotTree(t, root))
}

func TestApply_planConflictWritesNothing(t *testing.T) {
	root := copyFixtureProject(t, "v1-basic")
	manifest := filepath.Join(root, ".savepoint", "migrations", "v1-to-v2.yml")
	if err := os.MkdirAll(filepath.Dir(manifest), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(manifest, []byte("stale: true\n"), 0644); err != nil {
		t.Fatal(err)
	}
	plan := mustPlan(t, root)
	if len(plan.Conflicts) == 0 {
		t.Fatal("plan.Conflicts is empty, want the existing-manifest conflict")
	}
	before := snapshotTree(t, root)
	_, err := Apply(root, plan)
	if !errors.Is(err, ErrPlanConflict) {
		t.Fatalf("Apply() error = %v, want ErrPlanConflict", err)
	}
	assertSnapshotsEqual(t, before, snapshotTree(t, root))
}

func TestApply_lateCreateDestinationsPreserveUserContent(t *testing.T) {
	cases := []struct {
		name string
		path func(*ConversionPlan) string
	}{
		{name: "target", path: func(plan *ConversionPlan) string { return savepointPath(plan.Targets[0].InstallPath()) }},
		{name: "document", path: func(plan *ConversionPlan) string {
			for _, document := range plan.Documents {
				if document.Kind == DocumentIdea {
					return savepointPath(document.TargetPath)
				}
			}
			return ""
		}},
		{name: "archive", path: func(plan *ConversionPlan) string { return plan.Archives[0].ArchivePath }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := copyFixtureProject(t, "v1-basic")
			plan := mustPlan(t, root)
			relPath := tc.path(plan)
			if relPath == "" {
				t.Fatal("fixture did not produce the planned destination")
			}
			path := filepath.Join(root, filepath.FromSlash(relPath))
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				t.Fatal(err)
			}
			const userBytes = "user content created after planning\n"
			if err := os.WriteFile(path, []byte(userBytes), 0644); err != nil {
				t.Fatal(err)
			}
			_, err := Apply(root, plan)
			if !errors.Is(err, ErrCreateDestinationExists) {
				t.Fatalf("Apply() error = %v, want ErrCreateDestinationExists", err)
			}
			if got := readFile(t, path); got != userBytes {
				t.Fatalf("destination %s changed to %q; want %q", relPath, got, userBytes)
			}
			if !strings.Contains(err.Error(), "undo from the project root") {
				t.Errorf("partial-apply error omitted Git undo guidance: %v", err)
			}
		})
	}
}

func TestApply_writesRemovalsManifestAndSchemaInOrder(t *testing.T) {
	root := copyFixtureProject(t, "v1-history")
	plan := mustPlan(t, root)
	result, err := Apply(root, plan)
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	if len(result.WrittenPaths) == 0 || result.WrittenPaths[len(result.WrittenPaths)-1] != ".savepoint/config.yml" {
		t.Fatalf("written paths end at %v, want schema activation last", result.WrittenPaths)
	}
	manifestIndex := indexOf(result.WrittenPaths, manifestRelPath())
	if manifestIndex < 0 || manifestIndex != len(result.WrittenPaths)-2 {
		t.Fatalf("manifest index = %d in %v, want immediately before config activation", manifestIndex, result.WrittenPaths)
	}
	firstRemoval := len(result.WrittenPaths)
	for i, path := range result.WrittenPaths {
		for _, archive := range plan.Archives {
			if path == archive.SourcePath && i < firstRemoval {
				firstRemoval = i
			}
		}
	}
	if firstRemoval == len(result.WrittenPaths) || firstRemoval >= manifestIndex {
		t.Fatalf("removals at %d, manifest at %d; want removals before manifest: %v", firstRemoval, manifestIndex, result.WrittenPaths)
	}
	for i := 0; i < firstRemoval; i++ {
		if result.WrittenPaths[i] == manifestRelPath() {
			t.Fatalf("manifest written before removal at path index %d: %v", i, result.WrittenPaths)
		}
	}
}

func TestApply_failureReportsTouchedPathsAndGitUndo(t *testing.T) {
	root := copyFixtureProject(t, "v1-basic")
	plan := mustPlan(t, root)
	fault := errors.New("synthetic filesystem interruption")
	oldHook := afterApplyMutationHook
	calls := 0
	afterApplyMutationHook = func(string) error {
		calls++
		if calls == 1 {
			return fault
		}
		return nil
	}
	t.Cleanup(func() { afterApplyMutationHook = oldHook })

	_, err := Apply(root, plan)
	if !errors.Is(err, fault) {
		t.Fatalf("Apply() error = %v, want injected failure", err)
	}
	if calls != 1 {
		t.Fatalf("mutation hook ran %d times, want stop after first failure", calls)
	}
	if !strings.Contains(err.Error(), "written paths:") || !strings.Contains(err.Error(), "git --literal-pathspecs restore") || !strings.Contains(err.Error(), "git --literal-pathspecs clean") {
		t.Fatalf("Apply() error omitted touched paths or undo command:\n%v", err)
	}
	if version, readErr := data.ReadSchemaVersion(filepath.Join(root, ".savepoint", "config.yml")); readErr != nil || version != data.SchemaVersionV1 {
		t.Fatalf("schema after interrupted apply = %v, %v; want V1", version, readErr)
	}
	mustNotExist(t, filepath.Join(root, ".savepoint", "migrations", "v1-to-v2.yml"))
}

func TestApply_secondFullRunIsNoOp(t *testing.T) {
	root := copyFixtureProject(t, "v1-basic")
	if _, err := Apply(root, mustPlan(t, root)); err != nil {
		t.Fatalf("first Apply() error = %v", err)
	}
	before := snapshotTree(t, root)
	plan := mustPlan(t, root)
	if !plan.SchemaAlreadyV2 {
		t.Fatal("second Plan() did not report SchemaAlreadyV2")
	}
	result, err := Apply(root, plan)
	if err != nil {
		t.Fatalf("second Apply() error = %v", err)
	}
	if !result.AlreadyMigrated || len(result.WrittenPaths) != 0 {
		t.Fatalf("second Apply() result = %+v, want no-op", result)
	}
	assertSnapshotsEqual(t, before, snapshotTree(t, root))
}

func TestApply_preservesExistingMigrationsDirectoryContent(t *testing.T) {
	root := copyFixtureProject(t, "v1-basic")
	if err := os.MkdirAll(filepath.Join(root, ".savepoint", "migrations"), 0755); err != nil {
		t.Fatal(err)
	}
	readme := filepath.Join(root, ".savepoint", "migrations", "README.md")
	skill := filepath.Join(root, ".savepoint", "migrations", "savepoint-audit-SKILL.md")
	if err := os.WriteFile(readme, []byte("# preexisting readme\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(skill, []byte("# archived legacy skill\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := Apply(root, mustPlan(t, root)); err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	if got := readFile(t, readme); got != "# preexisting readme\n" {
		t.Errorf("migrations/README.md changed: %q", got)
	}
	if got := readFile(t, skill); got != "# archived legacy skill\n" {
		t.Errorf("migrations/savepoint-audit-SKILL.md changed: %q", got)
	}
	mustExist(t, filepath.Join(root, ".savepoint", "migrations", "v1-to-v2.yml"))
}

func TestWriteTempAndRenameReplacesAtomicallyAndCleansTemporaryFiles(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "record.md")
	if err := os.WriteFile(path, []byte("original\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := writeTempAndRename(path, []byte("replacement\n"), 0600, false); err != nil {
		t.Fatalf("writeTempAndRename() error = %v", err)
	}
	if got := readFile(t, path); got != "replacement\n" {
		t.Fatalf("replaced content = %q", got)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	// Windows keeps only a read-only flag, so a writable file always reports
	// 666 there; the requested mode is checked where it can be represented.
	if runtime.GOOS != "windows" && info.Mode().Perm() != 0600 {
		t.Errorf("replacement mode = %o, want 600", info.Mode().Perm())
	}
	matches, err := filepath.Glob(filepath.Join(dir, ".savepoint-migrate-apply-*"))
	if err != nil || len(matches) != 0 {
		t.Fatalf("temporary files = %v, %v; want none", matches, err)
	}
}

func replaceOnce(t *testing.T, content, old, new string) string {
	t.Helper()
	if !strings.Contains(content, old) {
		t.Fatalf("content does not contain %q", old)
	}
	return strings.Replace(content, old, new, 1)
}

func indexOf(paths []string, want string) int {
	for i, path := range paths {
		if path == want {
			return i
		}
	}
	return -1
}
