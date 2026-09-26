package init

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/doctor"
	"github.com/opencode/savepoint/internal/resume"
)

// v2ScaffoldSavepointFiles are the files at the root of
// templates/project-v2/.savepoint that must be present.
var v2ScaffoldSavepointFiles = []string{"Idea.md", "Design.md", "Guardrails.md", "config.yml", "router.md"}

// v2ScaffoldForbiddenFiles are V1-lifecycle files that must never reach the
// V2 scaffold: a fresh project has no epics, no release helper document, no
// Concept, no Health-Check, and no audit register.
var v2ScaffoldForbiddenFiles = []string{"PRD.md", "Concept.md", "Health-Check.md"}
var v2ScaffoldForbiddenDirs = []string{"audit"}

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
	goalPath := filepath.Join(root, "releases", "G-001-first-goal", "Release.md")
	if _, err := os.Stat(goalPath); err != nil {
		t.Errorf("templates/project-v2/.savepoint/releases/G-001-first-goal/Release.md missing: %v", err)
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

func TestV2ScaffoldCreatesProjectGoalWithInterpolatedName(t *testing.T) {
	target := t.TempDir()
	templates := os.DirFS(filepath.Join("..", "..", "templates", "project-v2"))
	if err := Scaffold(templates, target, "myapp", false); err != nil {
		t.Fatalf("Scaffold() from templates/project-v2 error = %v", err)
	}

	goalPath := filepath.Join(target, ".savepoint", "releases", "G-001-first-goal", "Release.md")
	content, err := os.ReadFile(goalPath)
	if err != nil {
		t.Fatalf("read fresh Goal record: %v", err)
	}
	for _, want := range []string{
		"id: G-001",
		"title: myapp",
		"status: in_progress",
		"## Outcome",
		"## Why",
		"## Success Conditions",
		"## Boundaries",
	} {
		if !strings.Contains(string(content), want) {
			t.Errorf("fresh Goal record missing %q", want)
		}
	}
	if strings.Contains(string(content), "{{PROJECT_NAME}}") {
		t.Error("fresh Goal record retained the project-name placeholder")
	}
	index, err := data.LoadV2Index(filepath.Join(target, ".savepoint"))
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	routerBytes, err := os.ReadFile(filepath.Join(target, ".savepoint", "router.md"))
	if err != nil {
		t.Fatalf("read fresh router: %v", err)
	}
	router, err := data.NewRouterReader().ReadStateV2(string(routerBytes))
	if err != nil {
		t.Fatalf("ReadStateV2() error = %v", err)
	}
	selection, diagnostic := data.ResolveSelection(index, router)
	if diagnostic != nil || selection.Release == nil || selection.Release.ID != "G-001" {
		t.Fatalf("fresh Goal selection = %+v, diagnostic = %+v; want resolved G-001", selection, diagnostic)
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
	assertContains(t, content, "release: G-001")
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

func TestDesignDocumentsTaskCreationAndDurableAllocation(t *testing.T) {
	content := readTemplate(t, filepath.Join("..", ".."), ".savepoint", "Design.md")
	for _, phrase := range []string{
		"`savepoint create-task --objective O-### --draft <path> [dir]`",
		"`.savepoint/task-ids.yml` stores `last_issued: N`",
		"`.savepoint/task-ids.lock`",
		"reserved ID stays retired",
		"After creating or renaming any other identity-bearing record",
	} {
		if !strings.Contains(content, phrase) {
			t.Errorf(".savepoint/Design.md is missing Task allocator detail %q", phrase)
		}
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
	assertContains(t, content, "Exception: agents may run `savepoint create-task --objective O-### --draft <path> [dir]` only to create a new Task from an ID-free draft.")
	assertContains(t, content, "After creating or renaming any other identity-bearing V2 record, run `savepoint resume` to require strict loading of the full V2 index.")
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

	// Optional files degrade gracefully; Goal context is required.
	for _, optional := range []string{"Concept", "Health-Check", "procedures file"} {
		if !strings.Contains(body, optional) {
			t.Errorf("adoption section does not name optional file %q as never required", optional)
		}
	}
	if !strings.Contains(body, "Their absence is normal, not a finding") {
		t.Error("adoption section does not state absence of optional files is normal, not a finding")
	}
	for _, phrase := range []string{"A Goal is required", "Choose a Goal", "savepoint doctor", "G-001", "savepoint init", "migration keeps or selects an existing live Goal"} {
		if !strings.Contains(body, phrase) {
			t.Errorf("adoption section does not explain required Goal context: %q", phrase)
		}
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
	if len(index.Releases) != 1 || len(index.Objectives) != 0 || len(index.Tasks) != 0 || len(index.Checks) != 0 || len(index.Issues) != 0 {
		t.Errorf("fresh V2 scaffold index has unexpected records: %+v", index)
	}
	report := doctor.RunV2Checks(filepath.Join(target, ".savepoint"))
	if report.HasProblems() {
		t.Fatalf("doctor reports problems for a fresh V2 scaffold: %s", report.Format())
	}
}

func TestV2ScaffoldResumesWithItsProjectGoalSelected(t *testing.T) {
	target := t.TempDir()
	templates := os.DirFS(filepath.Join("..", "..", "templates", "project-v2"))
	if err := Scaffold(templates, target, "myapp", false); err != nil {
		t.Fatalf("Scaffold() from templates/project-v2 error = %v", err)
	}

	savepointRoot := filepath.Join(target, ".savepoint")
	index, err := data.LoadV2Index(savepointRoot)
	if err != nil {
		t.Fatalf("data.LoadV2Index() on fresh V2 scaffold error = %v", err)
	}
	routerContent, err := os.ReadFile(filepath.Join(savepointRoot, "router.md"))
	if err != nil {
		t.Fatalf("read fresh V2 router: %v", err)
	}
	router, err := data.NewRouterReader().ReadStateV2(string(routerContent))
	if err != nil {
		t.Fatalf("read fresh V2 router state: %v", err)
	}
	next := data.ResolveNext(data.NextInput{Index: index, Router: router})
	var output bytes.Buffer
	if err := resume.Render(&output, next); err != nil {
		t.Fatalf("resume.Render() on fresh V2 scaffold error = %v", err)
	}
	if !strings.Contains(output.String(), "Next action:") {
		t.Fatalf("fresh V2 scaffold resume output = %q, want a Next action", output.String())
	}
	if strings.Contains(output.String(), "Choose a Goal") {
		t.Fatalf("fresh V2 scaffold resume unexpectedly asks for a Goal: %q", output.String())
	}
}
