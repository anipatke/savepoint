package init

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/data"
)

// v2ScaffoldSavepointFiles are the files templates/project-v2/.savepoint must
// contain, and no others besides objectives/.gitkeep.
var v2ScaffoldSavepointFiles = []string{"Idea.md", "Design.md", "Guardrails.md", "config.yml", "router.md"}

// v2ScaffoldForbiddenFiles are V1-lifecycle files that must never reach the
// V2 scaffold: a fresh project has no epics, no release helper document, no
// Concept, no Health-Check, and no audit register.
var v2ScaffoldForbiddenFiles = []string{"PRD.md", "Concept.md", "Health-Check.md"}
var v2ScaffoldForbiddenDirs = []string{"releases", "audit"}

func TestV2ScaffoldSavepointFileSet(t *testing.T) {
	root := filepath.Join("..", "..", "templates", "project-v2", ".savepoint")

	for _, name := range v2ScaffoldSavepointFiles {
		if _, err := os.Stat(filepath.Join(root, name)); err != nil {
			t.Errorf("templates/project-v2/.savepoint/%s missing: %v", name, err)
		}
	}
	if info, err := os.Stat(filepath.Join(root, "objectives")); err != nil || !info.IsDir() {
		t.Errorf("templates/project-v2/.savepoint/objectives missing or not a directory: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "objectives", ".gitkeep")); err != nil {
		t.Errorf("templates/project-v2/.savepoint/objectives/.gitkeep missing: %v", err)
	}

	for _, name := range v2ScaffoldForbiddenFiles {
		if _, err := os.Stat(filepath.Join(root, name)); err == nil {
			t.Errorf("templates/project-v2/.savepoint/%s must not exist on the V2 scaffold", name)
		}
	}
	for _, name := range v2ScaffoldForbiddenDirs {
		if _, err := os.Stat(filepath.Join(root, name)); err == nil {
			t.Errorf("templates/project-v2/.savepoint/%s must not exist on the V2 scaffold", name)
		}
	}
}

func TestV2ScaffoldDoesNotCreateReleaseRecordOrPromise(t *testing.T) {
	target := t.TempDir()
	templates := os.DirFS(filepath.Join("..", "..", "templates", "project-v2"))
	if err := Scaffold(templates, target, "myapp", false); err != nil {
		t.Fatalf("Scaffold() from templates/project-v2 error = %v", err)
	}

	if _, err := os.Stat(filepath.Join(target, ".savepoint", "releases")); !os.IsNotExist(err) {
		t.Fatalf("fresh V2 scaffold has a releases directory, stat err = %v", err)
	}

	err := filepath.WalkDir(target, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.Name() == "Release.md" {
			t.Errorf("fresh V2 scaffold contains a Release record at %s", path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk fresh V2 scaffold: %v", err)
	}
}

func TestV2ScaffoldConfigDeclaresSchemaTwoAndQualityGates(t *testing.T) {
	content := readTemplate(t, filepath.Join("..", ".."), "templates", "project-v2", ".savepoint", "config.yml")

	assertContains(t, content, "schema_version: 2")
	for _, key := range []string{"quality_gates:", "lint:", "typecheck:", "build:", "test:", "block_on_failure:", "gate_timeout:", "theme:"} {
		assertContains(t, content, key)
	}
	assertNotContains(t, content, "audit:")
}

func TestV2ScaffoldRouterOpensAtIdeaWithObjectiveField(t *testing.T) {
	content := readTemplate(t, filepath.Join("..", ".."), "templates", "project-v2", ".savepoint", "router.md")

	assertContains(t, content, "state: idea")
	assertContains(t, content, "objective:")
	assertContains(t, content, "| idea | savepoint-idea |")
	assertContains(t, content, "| design | savepoint-design |")
	assertContains(t, content, "| task | savepoint-task |")
	assertContains(t, content, "| check | savepoint-check |")

	for _, v1State := range []string{"pre-implementation", "epic-design", "epic-task-breakdown", "task-building", "audit-pending", "defect-building"} {
		assertNotContains(t, content, "state: "+v1State)
	}
}

// v2DesignSections are the six sections the V2 Design template must carry,
// matching the contract savepoint-design's SKILL.md states in full.
var v2DesignSections = []string{
	"## Architecture",
	"## Components/Codebase Map",
	"## Interfaces and Data Flow",
	"## Boundaries",
	"## Decisions",
	"## Current Technical State",
}

func TestV2ScaffoldDesignCarriesSixSections(t *testing.T) {
	content := readTemplate(t, filepath.Join("..", ".."), "templates", "project-v2", ".savepoint", "Design.md")

	for _, heading := range v2DesignSections {
		assertContains(t, content, heading)
	}
}

// v2AgentsGuideForbiddenVocabulary is V1-only vocabulary the V2 scaffold's
// AGENTS.md must never carry: V2 renamed Epic/PRD to Objective/Idea, folded
// the audit register into Check records, and has no defect-building state or
// "phase" as a lifecycle term.
var v2AgentsGuideForbiddenVocabulary = []string{
	"epic", "PRD", "audit register", "defect-building",
	"inactive until cutover", "phase",
}

func TestV2ScaffoldAgentsGuideIsLiveAndUsesV2Vocabulary(t *testing.T) {
	content := readTemplate(t, filepath.Join("..", ".."), "templates", "project-v2", "AGENTS.md")

	assertContains(t, content, "| idea | savepoint-idea |")
	assertContains(t, content, "| design | savepoint-design |")
	assertContains(t, content, "| task | savepoint-task |")
	assertContains(t, content, "| check | savepoint-check |")

	for _, stale := range v2AgentsGuideForbiddenVocabulary {
		assertNotContains(t, content, stale)
	}
}

// v2AdoptionLoadBearingPhrases are the statements the existing-codebase
// adoption section must carry: reconstruction is targeted reads, intent comes
// from the owner, and the managed guide block is the only authored-file
// exception.
var v2AdoptionLoadBearingPhrases = []string{
	"reconstructed from the code through targeted reads",
	"recorded in `.savepoint/Idea.md` through `savepoint-idea`, never inferred from source",
	"may add or refresh the Savepoint-managed block in an existing agent guide",
	"preserving every byte outside that block",
}

func TestV2AgentsGuideCarriesExistingCodebaseAdoptionSection(t *testing.T) {
	content := readTemplate(t, filepath.Join("..", ".."), "templates", "project-v2", "AGENTS.md")

	body, found := sectionBody(content, "## Existing Codebase Adoption")
	if !found || body == "" {
		t.Fatal("templates/project-v2/AGENTS.md missing or empty ## Existing Codebase Adoption section")
	}

	for _, phrase := range v2AdoptionLoadBearingPhrases {
		if !strings.Contains(body, phrase) {
			t.Errorf("adoption section missing load-bearing phrase %q", phrase)
		}
	}

	// No whole-repository scan, no automatic analysis, no model service call;
	// unread areas are recorded as unknown rather than inferred.
	for _, phrase := range []string{
		"whole-repository scan",
		"automatic analysis pass",
		"call to a model service",
		"recorded as unknown",
	} {
		if !strings.Contains(body, phrase) {
			t.Errorf("adoption section missing scope phrase %q", phrase)
		}
	}

	// Which record receives which finding.
	if !strings.Contains(body, "What exists goes to `.savepoint/Design.md`") {
		t.Error("adoption section does not name Design.md as the destination for what exists")
	}
	if !strings.Contains(body, "What the code is for goes to `.savepoint/Idea.md`") {
		t.Error("adoption section does not name Idea.md as the destination for what the code is for")
	}

	// Optional files degrade gracefully; absence is normal, not a finding.
	for _, optional := range []string{"Concept", "Health-Check", "procedures file", "Release record"} {
		if !strings.Contains(body, optional) {
			t.Errorf("adoption section does not name optional file %q as never required", optional)
		}
	}
	if !strings.Contains(body, "Their absence is normal, not a finding") {
		t.Error("adoption section does not state absence of optional files is normal, not a finding")
	}
}

func TestV2DesignCurrentTechnicalStateInstructsVerifiedOnly(t *testing.T) {
	content := readTemplate(t, filepath.Join("..", ".."), "templates", "project-v2", ".savepoint", "Design.md")

	body, found := sectionBody(content, "## Current Technical State")
	if !found || body == "" {
		t.Fatal("templates/project-v2/.savepoint/Design.md missing or empty ## Current Technical State section")
	}

	if !strings.Contains(body, "Record only what was verified by reading") {
		t.Error("Current Technical State does not instruct recording only what was verified by reading")
	}
	if !strings.Contains(body, "name what has not yet been examined as unknown") {
		t.Error("Current Technical State does not instruct naming what has not yet been examined")
	}
}

// TestV2ScaffoldIntoPopulatedDirectoryPreservesExistingFiles proves adoption
// into an existing codebase touches nothing it did not ship: every
// pre-existing file, including one whose name collides with nothing the V2
// scaffold writes, stays byte-identical after Scaffold runs.
func TestV2ScaffoldIntoPopulatedDirectoryPreservesExistingFiles(t *testing.T) {
	dir := t.TempDir()

	preexisting := map[string]string{
		"main.go":          "package main\n\nfunc main() {}\n",
		"src/app.py":       "print('hello')\n",
		"README.md":        "# My Existing Project\n",
		"requirements.txt": "flask==2.0.0\n", // collides with nothing templates/project-v2 writes
	}
	for relPath, content := range preexisting {
		full := filepath.Join(dir, filepath.FromSlash(relPath))
		if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	templates := os.DirFS(filepath.Join("..", "..", "templates", "project-v2"))
	if err := Scaffold(templates, dir, "myapp", false); err != nil {
		t.Fatalf("Scaffold() into populated directory error = %v", err)
	}

	for relPath, want := range preexisting {
		got, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(relPath)))
		if err != nil {
			t.Fatalf("read %s after scaffold: %v", relPath, err)
		}
		if string(got) != want {
			t.Errorf("%s changed by scaffold: got %q, want %q", relPath, string(got), want)
		}
	}

	// The V2 scaffold itself still lands as expected alongside the untouched files.
	for _, name := range v2ScaffoldSavepointFiles {
		if _, err := os.Stat(filepath.Join(dir, ".savepoint", name)); err != nil {
			t.Errorf("V2 scaffold file .savepoint/%s not installed: %v", name, err)
		}
	}
}

// TestV2ScaffoldLoadsCleanThroughRuntimeSchemaAndIndex proves the V2 tree is
// structurally sound: its schema passes the runtime gate and its empty
// objectives/ tree loads without a diagnostic.
func TestV2ScaffoldLoadsCleanThroughRuntimeSchemaAndIndex(t *testing.T) {
	target := t.TempDir()
	templates := os.DirFS(filepath.Join("..", "..", "templates", "project-v2"))
	if err := Scaffold(templates, target, "myapp", false); err != nil {
		t.Fatalf("Scaffold() from templates/project-v2 error = %v", err)
	}

	if err := data.CheckRuntimeSchema(target); err != nil {
		t.Fatalf("data.CheckRuntimeSchema() on fresh V2 scaffold error = %v", err)
	}
	index, err := data.LoadV2Index(filepath.Join(target, ".savepoint"))
	if err != nil {
		t.Fatalf("data.LoadV2Index() on fresh V2 scaffold error = %v", err)
	}
	if len(index.Objectives) != 0 || len(index.Tasks) != 0 || len(index.Checks) != 0 || len(index.Issues) != 0 {
		t.Errorf("fresh V2 scaffold index not empty: %+v", index)
	}
}
