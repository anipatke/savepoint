package data

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// --- Migration fixture helpers ----------------------------------------
//
// internal/data/testdata/migration/ holds frozen V1 projects kept as
// migration source evidence for the V2 design. Each fixture directory
// carries a hand-authored manifest.yml recording every file's exact-byte
// SHA-256; these helpers read that manifest as the single source of truth
// and never regenerate a hash to compare it against itself. See
// testdata/migration/README.md.
//
// TestMigrationSourceHistory* (a dependent fixture task) reuses these
// helpers from within this same test package.

type migrationManifest struct {
	Fixture        string                    `yaml:"fixture"`
	Description    string                    `yaml:"description"`
	Files          []migrationManifestFile   `yaml:"files"`
	AbsentByDesign []migrationManifestAbsent `yaml:"absent_by_design"`
}

type migrationManifestFile struct {
	Path                    string `yaml:"path"`
	SHA256                  string `yaml:"sha256"`
	Role                    string `yaml:"role"`
	ScopedID                string `yaml:"scoped_id,omitempty"`
	RawStatus               string `yaml:"raw_status,omitempty"`
	RawPhase                string `yaml:"raw_phase,omitempty"`
	ExpectedStageAfterParse string `yaml:"expected_stage_after_parse,omitempty"`
	DependsOn               string `yaml:"depends_on,omitempty"`
	Encoding                string `yaml:"encoding,omitempty"`
	ExpectedClassification  string `yaml:"expected_classification,omitempty"`
}

type migrationManifestAbsent struct {
	Path   string `yaml:"path"`
	Reason string `yaml:"reason"`
}

func migrationFixtureDir(fixture string) string {
	return filepath.Join("testdata", "migration", fixture)
}

func loadMigrationManifest(t *testing.T, fixture string) migrationManifest {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(migrationFixtureDir(fixture), "manifest.yml"))
	if err != nil {
		t.Fatalf("read manifest for fixture %s: %v", fixture, err)
	}
	var manifest migrationManifest
	if err := yaml.Unmarshal(raw, &manifest); err != nil {
		t.Fatalf("parse manifest for fixture %s: %v", fixture, err)
	}
	return manifest
}

func readMigrationFile(t *testing.T, fixture, relPath string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(migrationFixtureDir(fixture), relPath))
	if err != nil {
		t.Fatalf("read migration fixture %s/%s: %v", fixture, relPath, err)
	}
	return data
}

func sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// assertFixtureBytesMatchManifest recomputes every recorded file's hash and
// compares it against manifest.yml's independently authored value. A test
// that only compares a hash to itself proves nothing; this is the check
// that actually catches fixture drift.
func assertFixtureBytesMatchManifest(t *testing.T, fixture string) {
	t.Helper()
	manifest := loadMigrationManifest(t, fixture)
	for _, f := range manifest.Files {
		got := sha256Hex(readMigrationFile(t, fixture, f.Path))
		if got != f.SHA256 {
			t.Errorf("fixture %s/%s sha256 = %s, want %s (recorded in manifest.yml)", fixture, f.Path, got, f.SHA256)
		}
	}
}

func TestMigrationSourceBasicInventory(t *testing.T) {
	const fixture = "v1-basic"
	manifest := loadMigrationManifest(t, fixture)

	if len(manifest.Files) == 0 {
		t.Fatal("manifest.yml lists no files")
	}
	assertFixtureBytesMatchManifest(t, fixture)

	recorded := map[string]bool{}
	for _, f := range manifest.Files {
		recorded[filepath.ToSlash(f.Path)] = true
	}

	root := migrationFixtureDir(fixture)
	var onDisk []string
	err := filepath.Walk(filepath.Join(root, "project"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		onDisk = append(onDisk, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		t.Fatalf("walk fixture project tree: %v", err)
	}

	if len(onDisk) != len(recorded) {
		t.Errorf("fixture has %d files on disk, manifest records %d: on disk = %v", len(onDisk), len(recorded), onDisk)
	}
	for _, rel := range onDisk {
		if !recorded[rel] {
			t.Errorf("fixture file %s is on disk but not recorded in manifest.yml", rel)
		}
	}

	for _, a := range manifest.AbsentByDesign {
		if _, err := os.Stat(filepath.Join(root, a.Path)); !os.IsNotExist(err) {
			t.Errorf("expected %s to be absent by design (%s), stat err = %v", a.Path, a.Reason, err)
		}
	}
}

func TestMigrationSourceBasicInterpretation(t *testing.T) {
	const fixture = "v1-basic"
	savepointRoot := filepath.Join(migrationFixtureDir(fixture), "project", ".savepoint")
	discover := NewDiscover()
	parser := NewParser()

	releases, err := discover.ListReleases(savepointRoot)
	if err != nil {
		t.Fatalf("ListReleases() error = %v", err)
	}
	if len(releases) != 1 || releases[0].ID != "v1" {
		t.Fatalf("ListReleases() = %v, want [v1]", releases)
	}

	epics, err := discover.ListEpics(savepointRoot, "v1")
	if err != nil {
		t.Fatalf("ListEpics() error = %v", err)
	}
	if len(epics) != 1 || epics[0].ID != "E01-example" {
		t.Fatalf("ListEpics() = %v, want [E01-example]", epics)
	}

	taskInfos, err := discover.ListTasks(savepointRoot, "v1", "E01-example")
	if err != nil {
		t.Fatalf("ListTasks() error = %v", err)
	}
	if len(taskInfos) != 2 || taskInfos[0].ID != "T001-original" || taskInfos[1].ID != "T002-follow-up" {
		t.Fatalf("ListTasks() = %v, want [T001-original T002-follow-up]", taskInfos)
	}

	tasks := make([]Task, 0, len(taskInfos))
	rawContent := map[string]string{}
	for _, info := range taskInfos {
		raw, err := os.ReadFile(info.Path)
		if err != nil {
			t.Fatalf("read %s: %v", info.Path, err)
		}
		rawContent[info.ID] = string(raw)

		task, err := parser.ParseTaskFile(info.Path, string(raw))
		if err != nil {
			t.Fatalf("ParseTaskFile(%s) error = %v", info.ID, err)
		}
		// Discovery normally supplies release/epic scope explicitly (see
		// internal/board/board.go's loadEpicTasks); mirror that here so
		// ResolveDependency below sees the same context a real load would.
		task.Release = "v1"
		task.Epic = "E01-example"
		tasks = append(tasks, *task)
	}
	t1, t2 := tasks[0], tasks[1]

	if t1.Column != ColumnDone {
		t.Errorf("T001.Column = %v, want done", t1.Column)
	}
	if t2.Column != ColumnInProgress {
		t.Errorf("T002.Column = %v, want in_progress", t2.Column)
	}
	if t2.Stage != StageBuild {
		t.Errorf("T002.Stage = %v, want build (healed from legacy phase: implementation)", t2.Stage)
	}

	// Raw YAML still reports the original legacy phase; the healed Stage
	// above is a parser interpretation, not a rewrite of the source.
	rawFrontmatter := rawFrontmatterMap(t, rawContent["T002-follow-up"])
	if rawFrontmatter["phase"] != "implementation" {
		t.Errorf("raw T002 frontmatter phase = %v, want implementation", rawFrontmatter["phase"])
	}
	if _, hasStage := rawFrontmatter["stage"]; hasStage {
		t.Error("raw T002 frontmatter has a stage key; source should only carry the legacy phase key")
	}

	// An unknown top-level field and an unknown nested field parse without
	// error and are simply absent from the typed Task.
	if rawFrontmatter["legacy_note"] != "kept from v1 authoring; not part of the current schema" {
		t.Errorf("raw T002 frontmatter legacy_note = %v", rawFrontmatter["legacy_note"])
	}
	metadata, ok := rawFrontmatter["metadata"].(map[string]interface{})
	if !ok {
		t.Fatalf("raw T002 frontmatter metadata = %T, want a nested mapping", rawFrontmatter["metadata"])
	}
	reviewer, ok := metadata["reviewer"].(map[string]interface{})
	if !ok || reviewer["name"] != "sam" {
		t.Errorf("raw T002 frontmatter metadata.reviewer = %v, want name sam", metadata["reviewer"])
	}

	// The dependency reference uses the legacy short-with-suffix form
	// ("T001-original", not "T001" or the full scoped ID); ResolveDependency
	// must still land on T001 using T002's release/epic scope.
	if len(t2.DependsOn) != 1 || t2.DependsOn[0] != "T001-original" {
		t.Fatalf("T002.DependsOn = %v, want [T001-original]", t2.DependsOn)
	}
	resolution := ResolveDependency(t2.DependsOn[0], t2, tasks, map[string]string{})
	if resolution.Kind != DependencyTask || resolution.ID != t1.ID || resolution.TaskStatus != ColumnDone {
		t.Errorf("ResolveDependency() = %+v, want task %s done", resolution, t1.ID)
	}

	// Source bytes are CRLF throughout; the parser normalizes internally
	// but never rewrites the file, so the fixture itself must still be CRLF.
	if !strings.Contains(rawContent["T002-follow-up"], "\r\n") {
		t.Error("raw T002 content has no CRLF line endings; fixture must stay CRLF")
	}

	// The body comment survives in raw content but never becomes a
	// checklist item.
	if !strings.Contains(rawContent["T002-follow-up"], "<!-- Author note:") {
		t.Error("raw T002 content lost its author-note comment")
	}
	for _, item := range t2.Checklist {
		if strings.Contains(item.Text, "Author note") {
			t.Errorf("Checklist picked up the body comment: %+v", item)
		}
	}

	// Multiline authored text: raw bytes keep every authored line, but the
	// two checklist extractors normalize it differently. Implementation
	// Plan items join their continuation lines into one string...
	if len(t2.Checklist) != 2 {
		t.Fatalf("T002.Checklist = %+v, want 2 items", t2.Checklist)
	}
	wantJoined := "Extend it with the following steps: first gather inputs, then apply the follow-up transform, then record results"
	if t2.Checklist[1].Text != wantJoined {
		t.Errorf("T002.Checklist[1].Text = %q, want %q", t2.Checklist[1].Text, wantJoined)
	}

	// ...but Acceptance Criteria extraction does not merge continuation
	// lines at all, so only the first authored line survives into the typed
	// model. This is exactly the gap raw-byte preservation exists to cover.
	if len(t2.Acceptance) != 2 {
		t.Fatalf("T002.Acceptance = %+v, want 2 items", t2.Acceptance)
	}
	wantTruncated := "Follow-up documented across more than"
	if t2.Acceptance[1] != wantTruncated {
		t.Errorf("T002.Acceptance[1] = %q, want %q (Acceptance extraction drops continuation lines)", t2.Acceptance[1], wantTruncated)
	}
}

// TestMigrationSourceBasicLoadProjectDispatch proves the E42 schema-dispatch
// boundary (LoadProject) reaches the frozen v1-basic fixture through
// transitional V1 dispatch exactly as calling Discover directly would, and
// that doing so never touches the frozen source bytes.
func TestMigrationSourceBasicLoadProjectDispatch(t *testing.T) {
	const fixture = "v1-basic"
	savepointRoot := filepath.Join(migrationFixtureDir(fixture), "project", ".savepoint")

	project, err := LoadProject(savepointRoot)
	if err != nil {
		t.Fatalf("LoadProject() error = %v", err)
	}
	if project.SchemaVersion != SchemaVersionV1 {
		t.Fatalf("LoadProject() SchemaVersion = %v, want SchemaVersionV1", project.SchemaVersion)
	}
	if project.V1 == nil {
		t.Fatal("LoadProject() V1 discover adapter = nil, want non-nil for transitional V1 dispatch")
	}

	epics, err := project.V1.ListEpics(savepointRoot, "v1")
	if err != nil {
		t.Fatalf("project.V1.ListEpics() error = %v", err)
	}
	if len(epics) != 1 || epics[0].ID != "E01-example" {
		t.Fatalf("project.V1.ListEpics() = %v, want [E01-example]", epics)
	}

	assertFixtureBytesMatchManifest(t, fixture)
}

func TestMigrationSourceBasicFailures(t *testing.T) {
	const fixture = "v1-basic"
	const relT1 = "project/.savepoint/releases/v1/epics/E01-example/tasks/T001-original.md"
	const relT2 = "project/.savepoint/releases/v1/epics/E01-example/tasks/T002-follow-up.md"
	parser := NewParser()

	t.Run("malformed frontmatter produces the existing named parse error", func(t *testing.T) {
		raw := string(readMigrationFile(t, fixture, relT2))
		corrupted := strings.Replace(raw, "status: in_progress", "status: [unterminated", 1)
		if corrupted == raw {
			t.Fatal("test setup: corruption target line not found in fixture")
		}
		dest := filepath.Join(t.TempDir(), "T002-corrupt.md")
		if err := os.WriteFile(dest, []byte(corrupted), 0644); err != nil {
			t.Fatalf("write corrupted copy: %v", err)
		}

		content, err := os.ReadFile(dest)
		if err != nil {
			t.Fatalf("read corrupted copy: %v", err)
		}
		if _, err := parser.ParseTaskFile(dest, string(content)); err == nil {
			t.Fatal("ParseTaskFile() error = nil, want a malformed-YAML parse error")
		} else if !strings.Contains(err.Error(), "parse error for") {
			t.Errorf("ParseTaskFile() error = %v, want the existing named parse-error shape", err)
		}
	})

	t.Run("broken dependency reference resolves to the existing missing result", func(t *testing.T) {
		t1raw, err := parser.ParseTaskFile("T001-original.md", string(readMigrationFile(t, fixture, relT1)))
		if err != nil {
			t.Fatalf("ParseTaskFile(T001) error = %v", err)
		}
		t1raw.Release, t1raw.Epic = "v1", "E01-example"

		raw := string(readMigrationFile(t, fixture, relT2))
		broken := strings.Replace(raw, "depends_on: [T001-original]", "depends_on: [T099-missing]", 1)
		if broken == raw {
			t.Fatal("test setup: depends_on line not found in fixture")
		}
		dest := filepath.Join(t.TempDir(), "T002-broken-dep.md")
		if err := os.WriteFile(dest, []byte(broken), 0644); err != nil {
			t.Fatalf("write broken-dependency copy: %v", err)
		}

		content, err := os.ReadFile(dest)
		if err != nil {
			t.Fatalf("read broken-dependency copy: %v", err)
		}
		t2broken, err := parser.ParseTaskFile(dest, string(content))
		if err != nil {
			t.Fatalf("ParseTaskFile() error = %v, want the broken-dependency copy to still parse", err)
		}
		t2broken.Release, t2broken.Epic = "v1", "E01-example"

		resolution := ResolveDependency(t2broken.DependsOn[0], *t2broken, []Task{*t1raw, *t2broken}, map[string]string{})
		if resolution != (DependencyResolution{}) {
			t.Errorf("ResolveDependency() = %+v, want the existing zero-value missing-reference result", resolution)
		}
	})

	// Both failure variants ran against temporary copies only; the frozen
	// fixture bytes must be exactly what the manifest recorded before them.
	assertFixtureBytesMatchManifest(t, fixture)
}
