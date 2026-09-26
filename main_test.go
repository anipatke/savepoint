package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/opencode/savepoint/internal/data"
	savepointinit "github.com/opencode/savepoint/internal/init"
)

func TestMainVersionFlagPrintsVersion(t *testing.T) {
	result := runMainForTest(t, []string{"--version"}, "v9.8.7-test")

	if result.err != nil {
		t.Fatalf("savepoint --version failed: %v\nstderr: %s", result.err, result.stderr)
	}
	if strings.TrimSpace(result.stdout) != "v9.8.7-test" {
		t.Fatalf("stdout = %q, want version only", result.stdout)
	}
	if result.stderr != "" {
		t.Fatalf("stderr = %q, want empty", result.stderr)
	}
}

func TestMainHelpPrintsV2CommandContract(t *testing.T) {
	result := runMainForTest(t, []string{"--help"}, "")

	if result.err != nil {
		t.Fatalf("savepoint --help failed: %v\nstderr: %s", result.err, result.stderr)
	}
	for _, want := range []string{
		"Usage: savepoint <command> [options]",
		"create-task --objective <O-###> --draft <path> [dir]",
		"board [--objective <objective>]",
		"migrate [dir] [--apply]",
		"upgrade-assets [dir]",
		"Migration and asset upgrades are separate",
	} {
		if !strings.Contains(result.stdout, want) {
			t.Errorf("stdout = %q, want %q", result.stdout, want)
		}
	}
}

func TestMainCreateTaskAllocatesIDAndPrintsCreatedPath(t *testing.T) {
	dir := t.TempDir()
	objectiveDir := filepath.Join(dir, ".savepoint", "objectives", "O-001-first")
	if err := os.MkdirAll(objectiveDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".savepoint", "config.yml"), []byte("schema_version: 2\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	objective := "---\nid: O-001\ntitle: First objective\nstatus: planned\n---\n\n# First objective\n"
	if err := os.WriteFile(filepath.Join(objectiveDir, "Objective.md"), []byte(objective), 0o644); err != nil {
		t.Fatal(err)
	}
	draftPath := filepath.Join(dir, "draft.md")
	draft := "---\ntitle: Review the command\nobjective: O-001\nplanned_by: {role: planner, session: cli-test}\nstatus: planned\n---\n\n# Review the command\n"
	if err := os.WriteFile(draftPath, []byte(draft), 0o644); err != nil {
		t.Fatal(err)
	}

	result := runMainForTest(t, []string{"create-task", "--objective", "O-001", "--draft", draftPath, dir}, "")
	if result.err != nil {
		t.Fatalf("savepoint create-task failed: %v\nstderr: %s", result.err, result.stderr)
	}
	if want := "Created T-001 at .savepoint/objectives/O-001-first/tasks/T-001-review-the-command.md\n"; result.stdout != want {
		t.Fatalf("stdout = %q, want %q", result.stdout, want)
	}
	index, err := data.LoadV2Index(filepath.Join(dir, ".savepoint"))
	if err != nil {
		t.Fatalf("LoadV2Index() after command: %v", err)
	}
	if task := index.Tasks["T-001"]; task == nil || task.Objective != "O-001" {
		t.Fatalf("created Task = %+v, want T-001 owned by O-001", task)
	}
}

func TestMainCreateTaskConcurrentAcrossObjectivesStrictLoadsIndex(t *testing.T) {
	dir := t.TempDir()
	savepointDir := filepath.Join(dir, ".savepoint")
	if err := os.MkdirAll(savepointDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(savepointDir, "config.yml"), []byte("schema_version: 2\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	objectives := []struct {
		id, slug, title string
	}{
		{id: "O-001", slug: "O-001-first", title: "First objective"},
		{id: "O-002", slug: "O-002-second", title: "Second objective"},
	}
	for _, objective := range objectives {
		objectiveDir := filepath.Join(savepointDir, "objectives", objective.slug)
		if err := os.MkdirAll(objectiveDir, 0o755); err != nil {
			t.Fatal(err)
		}
		content := "---\nid: " + objective.id + "\ntitle: " + objective.title + "\nstatus: planned\n---\n\n# " + objective.title + "\n"
		if err := os.WriteFile(filepath.Join(objectiveDir, "Objective.md"), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	type result struct {
		objective string
		main      mainResult
	}
	results := make(chan result, len(objectives))
	for _, objective := range objectives {
		draftPath := filepath.Join(dir, "draft-"+objective.id+".md")
		draft := "---\ntitle: Concurrent review " + objective.id + "\nplanned_by: {role: planner, session: scenario}\nstatus: planned\n---\n\n# Concurrent review\n"
		if err := os.WriteFile(draftPath, []byte(draft), 0o644); err != nil {
			t.Fatal(err)
		}
		objective := objective
		go func() {
			mainResult := runMainForTest(t, []string{"create-task", "--objective", objective.id, "--draft", draftPath, dir}, "")
			results <- result{objective: objective.id, main: mainResult}
		}()
	}

	ids := make(map[string]string, len(objectives))
	for range len(objectives) {
		got := <-results
		if got.main.err != nil {
			t.Fatalf("savepoint create-task for %s failed: %v\nstderr: %s", got.objective, got.main.err, got.main.stderr)
		}
		fields := strings.Fields(got.main.stdout)
		if len(fields) != 4 || fields[0] != "Created" || fields[2] != "at" {
			t.Fatalf("create-task output for %s = %q, want assigned ID and path", got.objective, got.main.stdout)
		}
		id := fields[1]
		if previous, exists := ids[id]; exists {
			t.Fatalf("Objectives %s and %s received duplicate Task ID %s", previous, got.objective, id)
		}
		ids[id] = got.objective
	}

	index, err := data.LoadV2Index(savepointDir)
	if err != nil {
		t.Fatalf("strict V2 load after concurrent create-task calls: %v", err)
	}
	if len(index.Tasks) != len(objectives) || len(ids) != len(objectives) {
		t.Fatalf("created IDs=%v, indexed Tasks=%v; want one distinct Task per Objective", ids, index.Tasks)
	}
	for id, objective := range ids {
		if task := index.Tasks[id]; task == nil || task.Objective != objective {
			t.Errorf("indexed Task %s = %+v, want owner %s", id, task, objective)
		}
	}
}

func TestMainInitHelpStillUsesNormalDispatch(t *testing.T) {
	result := runMainForTest(t, []string{"init", "--help"}, "")

	if result.err != nil {
		t.Fatalf("savepoint init --help failed: %v\nstderr: %s", result.err, result.stderr)
	}
	if !strings.Contains(result.stdout, "Usage: init [dir] [--force] [--install]") {
		t.Fatalf("stdout = %q, want init usage", result.stdout)
	}
}

func TestMainUpgradeAssetsPrintsPartialWorkOnFailure(t *testing.T) {
	// A write failure part-way through must not hide what was already applied:
	// the user needs the report to know which files changed and where any
	// backup went.
	if runtime.GOOS == "windows" {
		t.Skip("directory permissions are not enforced the same way on Windows")
	}
	if os.Geteuid() == 0 {
		t.Skip("running as root bypasses directory permissions")
	}

	dir := t.TempDir()
	mkdirAll(t, filepath.Join(dir, ".savepoint"))
	if err := os.WriteFile(filepath.Join(dir, ".savepoint", "config.yml"), []byte("schema_version: 2\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// A stale skill in a directory that cannot be written: the walk reaches it
	// after it has already installed earlier skills.
	blocked := filepath.Join(dir, "agent-skills", "savepoint-check")
	mkdirAll(t, blocked)
	if err := os.WriteFile(filepath.Join(blocked, "SKILL.md"), []byte("# Stale\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(blocked, 0500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(blocked, 0755) })

	result := runMainForTest(t, []string{"upgrade-assets", dir}, "")

	if result.err == nil {
		t.Fatal("upgrade-assets succeeded, want the blocked write to fail")
	}
	if !strings.Contains(result.stdout, "Upgrade Report:") {
		t.Errorf("stdout = %q, want the partial report", result.stdout)
	}
	if !strings.Contains(result.stdout, "failed  agent-skills/savepoint-check/SKILL.md") {
		t.Errorf("stdout = %q, want the failed path named", result.stdout)
	}
	if !strings.Contains(result.stdout, "agent-skills/references/check-method.md") {
		t.Errorf("stdout = %q, want the already-applied work named", result.stdout)
	}
}

// v1OnlyPaths are V1-specific paths that must never appear on a fresh V2
// scaffold. The V2 Goal record also uses .savepoint/releases/.
var v1OnlyPaths = []string{
	filepath.Join(".savepoint", "releases", "v2", "epics"),
	filepath.Join(".savepoint", "releases", "v2", "v2-PRD.md"),
	filepath.Join(".savepoint", "epics"),
	filepath.Join(".savepoint", "PRD.md"),
	filepath.Join(".savepoint", "Concept.md"),
	filepath.Join(".savepoint", "Health-Check.md"),
	filepath.Join(".savepoint", "audit"),
}

// v1OnlySkills are the nine V1 skills that must never reach a fresh V2
// project; savepoint init writes only the four V2 skills.
var v1OnlySkills = []string{
	"savepoint-draft-prd", "savepoint-system-design", "savepoint-create-task",
	"savepoint-build-task", "savepoint-audit-epic", "savepoint-audit-task",
	"savepoint-audit-register", "savepoint-create-defect", "savepoint-create-plan",
}

func TestMainInitScaffoldsV2ProjectWithProjectGoal(t *testing.T) {
	dir := t.TempDir()

	result := runMainForTest(t, []string{"init", dir}, "")
	if result.err != nil {
		t.Fatalf("savepoint init failed: %v\nstderr: %s", result.err, result.stderr)
	}

	configPath := filepath.Join(dir, ".savepoint", "config.yml")
	version, err := data.ReadSchemaVersion(configPath)
	if err != nil {
		t.Fatalf("ReadSchemaVersion() error = %v", err)
	}
	if version != data.SchemaVersionV2 {
		t.Fatalf("SchemaVersion = %v, want SchemaVersionV2", version)
	}

	if err := data.CheckRuntimeSchema(dir); err != nil {
		t.Fatalf("CheckRuntimeSchema() on fresh init error = %v", err)
	}

	index, err := data.LoadV2Index(filepath.Join(dir, ".savepoint"))
	if err != nil {
		t.Fatalf("LoadV2Index() on fresh init error = %v", err)
	}
	if len(index.Releases) != 1 || len(index.Objectives) != 0 || len(index.Tasks) != 0 || len(index.Checks) != 0 || len(index.Issues) != 0 {
		t.Errorf("fresh init V2 index has unexpected records: %+v", index)
	}
	goalPath := filepath.Join(dir, ".savepoint", "releases", "G-001-first-goal", "Release.md")
	goal, err := os.ReadFile(goalPath)
	if err != nil {
		t.Fatalf("read fresh-init G-001 Goal: %v", err)
	}
	if !strings.Contains(string(goal), "title: "+filepath.Base(dir)) {
		t.Errorf("fresh-init Goal title = %q, want project name %q", goal, filepath.Base(dir))
	}
	router, err := os.ReadFile(filepath.Join(dir, ".savepoint", "router.md"))
	if err != nil {
		t.Fatalf("read fresh-init router: %v", err)
	}
	if !strings.Contains(string(router), "release: G-001") {
		t.Errorf("fresh-init router does not select G-001: %q", router)
	}
	agents, err := os.ReadFile(filepath.Join(dir, "AGENTS.md"))
	if err != nil {
		t.Fatalf("read fresh-init AGENTS.md: %v", err)
	}
	if !strings.Contains(string(agents), "Exception: agents may run `savepoint create-task --objective O-### --draft <path> [dir]` only to create a new Task from an ID-free draft.") {
		t.Error("fresh init AGENTS.md does not allow the narrow Task creation command")
	}
}

func TestMainInitWritesNoV1OnlyPathOrSkill(t *testing.T) {
	dir := t.TempDir()

	result := runMainForTest(t, []string{"init", dir}, "")
	if result.err != nil {
		t.Fatalf("savepoint init failed: %v\nstderr: %s", result.err, result.stderr)
	}

	for _, rel := range v1OnlyPaths {
		if _, err := os.Stat(filepath.Join(dir, rel)); !os.IsNotExist(err) {
			t.Errorf("fresh V2 init has V1-only path %s (stat err = %v)", rel, err)
		}
	}
	for _, skill := range v1OnlySkills {
		path := filepath.Join(dir, "agent-skills", skill)
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Errorf("fresh V2 init has V1 skill %s (stat err = %v)", skill, err)
		}
	}
}

func TestMainInitManifestRecordsExactlyTheFourV2Skills(t *testing.T) {
	dir := t.TempDir()

	result := runMainForTest(t, []string{"init", dir}, "")
	if result.err != nil {
		t.Fatalf("savepoint init failed: %v\nstderr: %s", result.err, result.stderr)
	}

	manifest, err := savepointinit.LoadManifest(dir)
	if err != nil {
		t.Fatalf("LoadManifest() error = %v", err)
	}

	want := []string{
		"agent-skills/savepoint-idea/SKILL.md",
		"agent-skills/savepoint-design/SKILL.md",
		"agent-skills/savepoint-task/SKILL.md",
		"agent-skills/savepoint-check/SKILL.md",
	}
	if len(manifest.Skills) != len(want) {
		t.Fatalf("manifest.Skills = %v, want exactly %v", manifest.Skills, want)
	}
	for _, key := range want {
		if _, ok := manifest.Hash(key); !ok {
			t.Errorf("manifest missing provenance entry for %s", key)
		}
	}
}

// fileHash returns the hex SHA-256 of path's content, for tests that must
// prove a file's bytes did not change rather than merely that it still
// exists.
func fileHash(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:])
}

func TestMainInitOnExistingCodebasePreservesEveryFileByteIdentical(t *testing.T) {
	dir := t.TempDir()

	seed := map[string]string{
		"package.json": `{"name":"demo"}` + "\n",
		"README.md":    "# Demo\n",
		"src/index.js": "console.log('hi')\n",
		".gitignore":   "node_modules\n",
	}
	for rel, content := range seed {
		path := filepath.Join(dir, filepath.FromSlash(rel))
		mkdirAll(t, filepath.Dir(path))
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	before := map[string]string{}
	for rel := range seed {
		before[rel] = fileHash(t, filepath.Join(dir, filepath.FromSlash(rel)))
	}

	result := runMainForTest(t, []string{"init", dir}, "")
	if result.err != nil {
		t.Fatalf("savepoint init failed: %v\nstderr: %s", result.err, result.stderr)
	}

	for rel, want := range before {
		got := fileHash(t, filepath.Join(dir, filepath.FromSlash(rel)))
		if got != want {
			t.Errorf("%s content changed by init", rel)
		}
	}
}

const migrateFixtureProject = "internal/data/testdata/migration/v1-basic/project"

// copyMigrateFixture copies the frozen v1-basic fixture into a fresh
// temporary directory so command tests can preview or apply against it
// without ever mutating the fixture itself.
func copyMigrateFixture(t *testing.T) string {
	t.Helper()
	dst := t.TempDir()
	err := filepath.WalkDir(migrateFixtureProject, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(migrateFixtureProject, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, content, 0644)
	})
	if err != nil {
		t.Fatalf("copy migrate fixture: %v", err)
	}
	return dst
}

type fileSnapshot struct {
	hash    string
	modTime time.Time
}

// snapshotDir records every regular file's content hash and mtime under
// root, so a test can prove a command wrote nothing at all rather than
// merely "no new top-level file".
func snapshotDir(t *testing.T, root string) map[string]fileSnapshot {
	t.Helper()
	snap := map[string]fileSnapshot{}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		sum := sha256.Sum256(content)
		snap[rel] = fileSnapshot{hash: hex.EncodeToString(sum[:]), modTime: info.ModTime()}
		return nil
	})
	if err != nil {
		t.Fatalf("snapshot %s: %v", root, err)
	}
	return snap
}

func assertSameSnapshot(t *testing.T, before, after map[string]fileSnapshot) {
	t.Helper()
	if len(before) != len(after) {
		t.Fatalf("file count changed: before %d, after %d", len(before), len(after))
	}
	for path, want := range before {
		got, ok := after[path]
		if !ok {
			t.Fatalf("%s was removed", path)
		}
		if got.hash != want.hash {
			t.Fatalf("%s content changed", path)
		}
		if !got.modTime.Equal(want.modTime) {
			t.Fatalf("%s mtime changed: before %v, after %v", path, want.modTime, got.modTime)
		}
	}
}

func writeMigrateMinimalProject(t *testing.T, root string) {
	t.Helper()
	mkdirAll(t, filepath.Join(root, ".savepoint"))
	if err := os.WriteFile(filepath.Join(root, ".savepoint", "config.yml"), []byte("quality_gates: {}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".savepoint", "router.md"), []byte("# Router\n"), 0644); err != nil {
		t.Fatal(err)
	}
	releaseDir := filepath.Join(root, ".savepoint", "releases", "v1")
	mkdirAll(t, releaseDir)
	if err := os.WriteFile(filepath.Join(releaseDir, "v1-PRD.md"), []byte("---\nname: V1\nstatus: in_progress\n---\n\n# V1\n"), 0644); err != nil {
		t.Fatal(err)
	}
}

// writeMigrateAmbiguousProject writes a minimal V1 project whose single task
// carries an unrecognized status, so Plan raises exactly one blocking
// AmbiguityUnrecognizedLifecycle ambiguity. It returns that ambiguity's
// stable ID (kind:path, matching addAmbiguity in internal/migrate).
func writeMigrateAmbiguousProject(t *testing.T, root string) string {
	t.Helper()
	writeMigrateMinimalProject(t, root)

	epicDir := filepath.Join(root, ".savepoint", "releases", "v1", "epics", "E01-x")
	mkdirAll(t, epicDir)
	if err := os.WriteFile(filepath.Join(epicDir, "E01-Detail.md"), []byte("---\nstatus: in_progress\n---\n\n# E01\n"), 0644); err != nil {
		t.Fatal(err)
	}

	taskDir := filepath.Join(epicDir, "tasks")
	mkdirAll(t, taskDir)
	taskPath := ".savepoint/releases/v1/epics/E01-x/tasks/T001-weird.md"
	if err := os.WriteFile(filepath.Join(taskDir, "T001-weird.md"),
		[]byte("---\nid: E01-x/T001-weird\nstatus: escalated\ndepends_on: []\n---\n\n# T001\n"), 0644); err != nil {
		t.Fatal(err)
	}

	return "unrecognized_lifecycle:" + taskPath
}

func TestMainMigrateHelpStillUsesNormalDispatch(t *testing.T) {
	result := runMainForTest(t, []string{"migrate", "--help"}, "")

	if result.err != nil {
		t.Fatalf("savepoint migrate --help failed: %v\nstderr: %s", result.err, result.stderr)
	}
	if !strings.Contains(result.stdout, "Usage: migrate [dir]") {
		t.Fatalf("stdout = %q, want migrate usage", result.stdout)
	}
	if !strings.Contains(result.stdout, "Preview is the default") {
		t.Fatalf("stdout = %q, want it to state preview is the default", result.stdout)
	}
}

func TestMainMigrateRejectsUnknownFlag(t *testing.T) {
	dir := copyMigrateFixture(t)
	before := snapshotDir(t, dir)

	result := runMainForTest(t, []string{"migrate", dir, "--bogus"}, "")

	if result.err == nil {
		t.Fatal("savepoint migrate --bogus succeeded, want a nonzero exit")
	}
	if !strings.Contains(result.stderr, "unknown migrate flag") {
		t.Fatalf("stderr = %q, want a named unknown-flag error", result.stderr)
	}
	if !strings.Contains(result.stdout, "Usage: migrate [dir]") {
		t.Fatalf("stdout = %q, want usage text rather than the flag being silently ignored", result.stdout)
	}
	assertSameSnapshot(t, before, snapshotDir(t, dir))
}

func TestMainMigratePreviewDefaultWritesNothing(t *testing.T) {
	dir := copyMigrateFixture(t)
	initMigrateGitRepo(t, dir)
	before := snapshotDir(t, dir)

	result := runMainForTest(t, []string{"migrate", dir}, "")

	if result.err != nil {
		t.Fatalf("savepoint migrate failed: %v\nstderr: %s", result.err, result.stderr)
	}
	if !strings.Contains(result.stdout, "Migration preview") {
		t.Fatalf("stdout = %q, want a preview report", result.stdout)
	}
	if !strings.Contains(result.stdout, "Will create") || strings.Contains(result.stdout, "Planned records") {
		t.Fatalf("stdout = %q, want the summary preview by default", result.stdout)
	}
	assertSameSnapshot(t, before, snapshotDir(t, dir))
}

func TestMainMigrateVerboseListsEveryPlannedRecord(t *testing.T) {
	dir := copyMigrateFixture(t)
	initMigrateGitRepo(t, dir)
	before := snapshotDir(t, dir)

	result := runMainForTest(t, []string{"migrate", dir, "--verbose"}, "")

	if result.err != nil {
		t.Fatalf("savepoint migrate --verbose failed: %v\nstderr: %s", result.err, result.stderr)
	}
	for _, want := range []string{"Planned records", "Archives"} {
		if !strings.Contains(result.stdout, want) {
			t.Fatalf("stdout = %q, want the full listing section %q", result.stdout, want)
		}
	}
	assertSameSnapshot(t, before, snapshotDir(t, dir))
}

func TestMainMigrateDryRunIsSynonymOfDefault(t *testing.T) {
	dir := copyMigrateFixture(t)
	initMigrateGitRepo(t, dir)
	before := snapshotDir(t, dir)

	result := runMainForTest(t, []string{"migrate", dir, "--dry-run"}, "")

	if result.err != nil {
		t.Fatalf("savepoint migrate --dry-run failed: %v\nstderr: %s", result.err, result.stderr)
	}
	assertSameSnapshot(t, before, snapshotDir(t, dir))
}

func TestMainMigrateApplyAndDryRunTogetherPreviews(t *testing.T) {
	dir := copyMigrateFixture(t)
	initMigrateGitRepo(t, dir)
	before := snapshotDir(t, dir)

	result := runMainForTest(t, []string{"migrate", dir, "--apply", "--dry-run"}, "")

	if result.err != nil {
		t.Fatalf("savepoint migrate --apply --dry-run failed: %v\nstderr: %s", result.err, result.stderr)
	}
	assertSameSnapshot(t, before, snapshotDir(t, dir))
}

func TestMainMigratePreviewIsDeterministic(t *testing.T) {
	dir := copyMigrateFixture(t)
	initMigrateGitRepo(t, dir)

	first := runMainForTest(t, []string{"migrate", dir}, "")
	if first.err != nil {
		t.Fatalf("first savepoint migrate failed: %v\nstderr: %s", first.err, first.stderr)
	}
	second := runMainForTest(t, []string{"migrate", dir}, "")
	if second.err != nil {
		t.Fatalf("second savepoint migrate failed: %v\nstderr: %s", second.err, second.stderr)
	}

	if normalizeMigratePreview(first.stdout) != normalizeMigratePreview(second.stdout) {
		t.Fatalf("preview output is not deterministic across runs:\n--- first ---\n%s\n--- second ---\n%s",
			first.stdout, second.stdout)
	}
}

// normalizeMigratePreview strips the operation id and generated-at lines,
// which legitimately vary run to run against a real clock, so the test can
// assert everything else about the report's structure and ordering is
// stable.
func normalizeMigratePreview(output string) string {
	var kept []string
	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(line, "operation:") || strings.HasPrefix(line, "generated:") {
			continue
		}
		kept = append(kept, line)
	}
	return strings.Join(kept, "\n")
}

func TestMainMigrateApplyWritesAndActivatesSchema(t *testing.T) {
	dir := copyMigrateFixture(t)
	initMigrateGitRepo(t, dir)

	result := runMainForTest(t, []string{"migrate", dir, "--apply"}, "")

	if result.err != nil {
		t.Fatalf("savepoint migrate --apply failed: %v\nstderr: %s", result.err, result.stderr)
	}
	if !strings.Contains(result.stdout, "migration complete") {
		t.Fatalf("stdout = %q, want a completion message", result.stdout)
	}

	config, err := os.ReadFile(filepath.Join(dir, ".savepoint", "config.yml"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(config), "schema_version: 2") {
		t.Fatalf("config.yml = %q, want schema_version activated to 2", config)
	}
}

func TestMainMigrateAmbiguousPlanBlocksAndNamesTheID(t *testing.T) {
	dir := t.TempDir()
	wantID := writeMigrateAmbiguousProject(t, dir)
	before := snapshotDir(t, dir)

	result := runMainForTest(t, []string{"migrate", dir}, "")

	if result.err == nil {
		t.Fatal("savepoint migrate over an unresolved blocking ambiguity succeeded, want a nonzero exit")
	}
	if !strings.Contains(result.stderr, wantID) {
		t.Fatalf("stderr = %q, want it to name the unresolved ambiguity id %s", result.stderr, wantID)
	}
	assertSameSnapshot(t, before, snapshotDir(t, dir))
}

func TestMainMigrateDecisionsFileResolvesAmbiguity(t *testing.T) {
	dir := t.TempDir()
	ambiguityID := writeMigrateAmbiguousProject(t, dir)

	decisionsPath := filepath.Join(dir, "decisions.yml")
	decisionsContent := "decisions:\n  - id: " + ambiguityID + "\n    value: planned\n"
	if err := os.WriteFile(decisionsPath, []byte(decisionsContent), 0644); err != nil {
		t.Fatal(err)
	}
	initMigrateGitRepo(t, dir)

	result := runMainForTest(t, []string{"migrate", dir, "--decisions", decisionsPath}, "")

	if result.err != nil {
		t.Fatalf("savepoint migrate --decisions failed: %v\nstderr: %s", result.err, result.stderr)
	}
	if strings.Contains(result.stdout, "BLOCKED") {
		t.Fatalf("stdout = %q, want the resolved ambiguity to no longer block", result.stdout)
	}
}

func TestMainMigrateMissingDirectory(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "does-not-exist")

	result := runMainForTest(t, []string{"migrate", missing}, "")

	if result.err == nil {
		t.Fatal("savepoint migrate over a missing directory succeeded, want a nonzero exit")
	}
	if !strings.Contains(result.stderr, "target directory does not exist") {
		t.Fatalf("stderr = %q, want a named missing-directory error", result.stderr)
	}
}

func TestMainMigrateNotASavepointProject(t *testing.T) {
	dir := t.TempDir()

	result := runMainForTest(t, []string{"migrate", dir}, "")

	if result.err == nil {
		t.Fatal("savepoint migrate over a non-Savepoint directory succeeded, want a nonzero exit")
	}
	if !strings.Contains(result.stderr, "is not a Savepoint project") {
		t.Fatalf("stderr = %q, want a named not-a-project error", result.stderr)
	}
}

func TestMainMigrateUnwritableDirectory(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("directory permissions are not enforced the same way on Windows")
	}
	if os.Geteuid() == 0 {
		t.Skip("running as root bypasses directory permissions")
	}

	dir := copyMigrateFixture(t)
	initMigrateGitRepo(t, dir)
	if err := os.Chmod(dir, 0500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(dir, 0755) })

	result := runMainForTest(t, []string{"migrate", dir, "--apply"}, "")

	if result.err == nil {
		t.Fatal("savepoint migrate over an unwritable directory succeeded, want a nonzero exit")
	}
	if !strings.Contains(result.stderr, "target directory is not writable") {
		t.Fatalf("stderr = %q, want a named unwritable-directory error", result.stderr)
	}
}

func initMigrateGitRepo(t *testing.T, dir string) {
	t.Helper()
	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git is not available; migration apply requires git")
	}
	run := func(args ...string) {
		t.Helper()
		command := exec.Command(gitPath, args...)
		command.Dir = dir
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, output)
		}
	}
	run("init", "-q")
	run("config", "user.name", "Savepoint Test")
	run("config", "user.email", "savepoint-test@example.invalid")
	// A commit can start detached auto-maintenance that creates and removes
	// .git lock files while a test snapshots the directory.
	run("config", "maintenance.auto", "false")
	run("config", "gc.auto", "0")
	run("add", "-A")
	run("commit", "-qm", "fixture baseline")
}

func TestMainUpgradeAssetsStillWorksAfterMigrateAdded(t *testing.T) {
	// A regression guard for the new migrate dispatch case: the existing
	// upgrade-assets command must behave exactly as before.
	dir := t.TempDir()
	mkdirAll(t, filepath.Join(dir, ".savepoint"))

	result := runMainForTest(t, []string{"upgrade-assets", dir, "--dry-run"}, "")

	if result.err != nil {
		t.Fatalf("savepoint upgrade-assets --dry-run failed: %v\nstderr: %s", result.err, result.stderr)
	}
	if !strings.Contains(result.stdout, "Upgrade Report:") {
		t.Fatalf("stdout = %q, want the upgrade report", result.stdout)
	}
}

func TestMainUpgradeAssetsV1ProjectRefusesMutationAndNamesMigrateRoute(t *testing.T) {
	dir := t.TempDir()
	mkdirAll(t, filepath.Join(dir, ".savepoint"))
	before := snapshotDir(t, dir)

	result := runMainForTest(t, []string{"upgrade-assets", dir}, "")
	if result.err != nil {
		t.Fatalf("savepoint upgrade-assets failed: %v\nstderr: %s", result.err, result.stderr)
	}

	if !strings.Contains(result.stdout, "savepoint migrate") {
		t.Errorf("stdout = %q, want the migrate-route note", result.stdout)
	}
	assertSameSnapshot(t, before, snapshotDir(t, dir))
}

func TestMainUpgradeAssetsV2ProjectInstallsOnlyV2Skills(t *testing.T) {
	dir := t.TempDir()
	mkdirAll(t, filepath.Join(dir, ".savepoint"))
	if err := os.WriteFile(filepath.Join(dir, ".savepoint", "config.yml"), []byte("schema_version: 2\n"), 0644); err != nil {
		t.Fatal(err)
	}

	result := runMainForTest(t, []string{"upgrade-assets", dir}, "")
	if result.err != nil {
		t.Fatalf("savepoint upgrade-assets failed: %v\nstderr: %s", result.err, result.stderr)
	}

	for _, skill := range []string{"savepoint-idea", "savepoint-design", "savepoint-task", "savepoint-check"} {
		if _, err := os.Stat(filepath.Join(dir, "agent-skills", skill, "SKILL.md")); err != nil {
			t.Errorf("V2 skill %s not installed: %v", skill, err)
		}
	}
	for _, reference := range []string{"check-method.md", "commands-and-procedures.md", "issue-capture.md"} {
		if _, err := os.Stat(filepath.Join(dir, "agent-skills", "references", reference)); err != nil {
			t.Errorf("V2 shared reference %s not installed: %v", reference, err)
		}
	}
	for _, skill := range v1OnlySkills {
		if _, err := os.Stat(filepath.Join(dir, "agent-skills", skill)); !os.IsNotExist(err) {
			t.Errorf("V1 skill %s installed on a V2 project, stat err = %v", skill, err)
		}
	}

	if strings.Contains(result.stdout, "savepoint migrate") {
		t.Errorf("stdout = %q, a V2 project should carry no migrate-route note", result.stdout)
	}

	config, err := os.ReadFile(filepath.Join(dir, ".savepoint", "config.yml"))
	if err != nil {
		t.Fatal(err)
	}
	if string(config) != "schema_version: 2\n" {
		t.Errorf("config.yml = %q, want it untouched by the upgrade", string(config))
	}
}

func mkdirAll(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
}

func TestMainHelperProcess(t *testing.T) {
	if os.Getenv("SAVEPOINT_TEST_MAIN") != "1" {
		return
	}
	if value := os.Getenv("SAVEPOINT_TEST_VERSION"); value != "" {
		version = value
	}
	os.Args = append([]string{"savepoint"}, helperArgs(os.Args)...)
	main()
}

type mainResult struct {
	stdout string
	stderr string
	err    error
}

func runMainForTest(t *testing.T, args []string, testVersion string) mainResult {
	t.Helper()

	cmdArgs := []string{"-test.run=TestMainHelperProcess", "--"}
	cmdArgs = append(cmdArgs, args...)
	cmd := exec.Command(os.Args[0], cmdArgs...)
	cmd.Env = append(os.Environ(), "SAVEPOINT_TEST_MAIN=1")
	if testVersion != "" {
		cmd.Env = append(cmd.Env, "SAVEPOINT_TEST_VERSION="+testVersion)
	}

	stdout, err := cmd.Output()
	stderr := ""
	if exitErr, ok := err.(*exec.ExitError); ok {
		stderr = string(exitErr.Stderr)
	}
	return mainResult{
		stdout: string(stdout),
		stderr: stderr,
		err:    err,
	}
}

func helperArgs(args []string) []string {
	for i, arg := range args {
		if arg == "--" {
			return args[i+1:]
		}
	}
	return nil
}
