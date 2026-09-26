package migrate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/data"
)

func documentByPath(p *ConversionPlan, path string) (PlannedDocument, bool) {
	for _, d := range p.Documents {
		if d.SourcePath == path {
			return d, true
		}
	}
	return PlannedDocument{}, false
}

// --- PRD -> Idea: byte-preserving relocation -------------------------------

// TestConvertIdea_relocatesPRDByteForByte proves the PRD.md -> Idea.md move
// is a relocation, never a re-render: the rendered content is exactly the
// source bytes, and the original is separately archived.
func TestConvertIdea_relocatesPRDByteForByte(t *testing.T) {
	for _, fixture := range []string{"v1-basic", "v1-history"} {
		t.Run(fixture, func(t *testing.T) {
			root := fixtureProjectRoot(fixture)
			const prdPath = ".savepoint/PRD.md"

			p := mustPlan(t, root)

			doc, ok := documentByPath(p, prdPath)
			if !ok {
				t.Fatalf("no PlannedDocument for %s", prdPath)
			}
			if doc.Kind != DocumentIdea {
				t.Errorf("Kind = %v, want DocumentIdea", doc.Kind)
			}
			if doc.TargetPath != "Idea.md" {
				t.Errorf("TargetPath = %q, want Idea.md", doc.TargetPath)
			}

			archived, ok := archiveByPath(p, prdPath)
			if !ok {
				t.Fatalf("PRD.md was not archived alongside its Idea.md relocation")
			}
			if archived.Role != RoleProductPRD {
				t.Errorf("archived Role = %v, want RoleProductPRD", archived.Role)
			}
			if archived.ArchivePath != ".savepoint/archive/v1/"+prdPath {
				t.Errorf("archived ArchivePath = %q, want .savepoint/archive/v1/%s", archived.ArchivePath, prdPath)
			}

			sourceBytes, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(prdPath)))
			if err != nil {
				t.Fatalf("read source PRD.md: %v", err)
			}

			content, err := ConvertIdea(root, doc)
			if err != nil {
				t.Fatalf("ConvertIdea error = %v", err)
			}
			if content != string(sourceBytes) {
				t.Errorf("ConvertIdea content does not match source bytes exactly:\ngot=%q\nwant=%q", content, string(sourceBytes))
			}
		})
	}
}

func TestConvertIdea_rejectsNonIdeaDocument(t *testing.T) {
	if _, err := ConvertIdea(fixtureProjectRoot("v1-basic"), PlannedDocument{Kind: DocumentRouter, SourcePath: ".savepoint/router.md"}); err == nil {
		t.Error("ConvertIdea on a non-idea document: error = nil, want an error")
	}
}

// --- router: state vocabulary + active selection ---------------------------

// TestConvertRouter_activeSelectionResolvesGlobalIDs proves that when the V1
// router's selected epic/task both converted, the V2 router names their
// allocated O-###/T-###, drops retired fields, and adds no migration note.
func TestConvertRouter_activeSelectionResolvesGlobalIDs(t *testing.T) {
	cases := []struct {
		fixture     string
		epicPath    string
		taskPath    string
		wantV2State string
	}{
		{
			fixture:     "v1-basic",
			epicPath:    ".savepoint/releases/v1/epics/E01-example/E01-Detail.md",
			taskPath:    ".savepoint/releases/v1/epics/E01-example/tasks/T002-follow-up.md",
			wantV2State: "task",
		},
		{
			fixture:     "v1-history",
			epicPath:    ".savepoint/releases/v1.1/epics/E01-example/E01-Detail.md",
			taskPath:    ".savepoint/releases/v1.1/epics/E01-example/tasks/T001-shared.md",
			wantV2State: "task",
		},
	}

	for _, tc := range cases {
		t.Run(tc.fixture, func(t *testing.T) {
			root := fixtureProjectRoot(tc.fixture)
			p := mustPlan(t, root)

			epicTarget, ok := targetByPath(p, tc.epicPath)
			if !ok {
				t.Fatalf("epic %s was not planned as an Objective target", tc.epicPath)
			}
			taskTarget, ok := targetByPath(p, tc.taskPath)
			if !ok {
				t.Fatalf("task %s was not planned as an active target", tc.taskPath)
			}

			doc, ok := documentByPath(p, ".savepoint/router.md")
			if !ok {
				t.Fatalf("no PlannedDocument for router.md")
			}
			if doc.Kind != DocumentRouter || doc.TargetPath != "router.md" {
				t.Errorf("router PlannedDocument = %+v, want Kind DocumentRouter, TargetPath router.md", doc)
			}

			content, err := ConvertRouter(root, p, doc)
			if err != nil {
				t.Fatalf("ConvertRouter error = %v", err)
			}

			if !strings.Contains(content, "state: "+tc.wantV2State) {
				t.Errorf("content missing %q:\n%s", "state: "+tc.wantV2State, content)
			}
			if !strings.Contains(content, "objective: "+epicTarget.GlobalID) {
				t.Errorf("content missing objective %q:\n%s", epicTarget.GlobalID, content)
			}
			if !strings.Contains(content, "task: "+taskTarget.GlobalID) {
				t.Errorf("content missing task %q:\n%s", taskTarget.GlobalID, content)
			}
			if !strings.Contains(content, "release: "+p.GoalSelection.GoalID) {
				t.Errorf("content missing live Goal %q:\n%s", p.GoalSelection.GoalID, content)
			}
			for _, retiredKey := range []string{"next_action:", "epic:", "defect:"} {
				if strings.Contains(content, retiredKey) {
					t.Errorf("content retained retired router key %q:\n%s", retiredKey, content)
				}
			}
			if strings.Contains(content, "## Migration Note") {
				t.Errorf("content has a Migration Note for a fully-resolved selection:\n%s", content)
			}
			if !strings.Contains(content, "## Read order") {
				t.Errorf("content lost surrounding V1 prose (## Read order):\n%s", content)
			}

			readBack, err := data.NewRouterReader().ReadState(content)
			if err != nil {
				t.Fatalf("rendered router content did not parse back: %v", err)
			}
			if readBack.State != tc.wantV2State {
				t.Errorf("readBack.State = %q, want %q", readBack.State, tc.wantV2State)
			}
			if readBack.Task != taskTarget.GlobalID {
				t.Errorf("readBack.Task = %q, want %q", readBack.Task, taskTarget.GlobalID)
			}
			if readBack.Release != p.GoalSelection.GoalID {
				t.Errorf("readBack.Release = %q, want selected live Goal %q", readBack.Release, p.GoalSelection.GoalID)
			}
		})
	}
}

func TestConvertRouter_fallbackSelectionPointsAtLiveGoal(t *testing.T) {
	tests := []struct {
		name          string
		routerRelease string
		releaseStatus string
		wantGoal      string
		wantGenerated bool
	}{
		{name: "missing release", routerRelease: "", releaseStatus: "in_progress", wantGoal: "R-001"},
		{name: "unresolvable release", routerRelease: "v9-missing", releaseStatus: "in_progress", wantGoal: "R-001"},
		{name: "archived release", routerRelease: "v1", releaseStatus: "done", wantGoal: "G-001", wantGenerated: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			writeFile(t, filepath.Join(root, ".savepoint", "config.yml"), "quality_gates: {}\n")
			writeFile(t, filepath.Join(root, ".savepoint", "router.md"), routerFixtureContent(
				"task-building", tc.routerRelease, "", "", `"Continue work."`))
			writeFile(t, filepath.Join(root, ".savepoint", "releases", "v1", "v1-PRD.md"),
				"---\nname: V1\nstatus: "+tc.releaseStatus+"\n---\n\n# V1\n")

			p := mustPlan(t, root)
			if p.GoalSelection.Generated != tc.wantGenerated {
				t.Fatalf("GoalSelection.Generated = %t, want %t: %+v", p.GoalSelection.Generated, tc.wantGenerated, p.GoalSelection)
			}
			if p.GoalSelection.GoalID != tc.wantGoal {
				t.Fatalf("GoalSelection.GoalID = %q, want %q", p.GoalSelection.GoalID, tc.wantGoal)
			}
			doc, ok := documentByPath(p, ".savepoint/router.md")
			if !ok {
				t.Fatal("no PlannedDocument for router.md")
			}
			content, err := ConvertRouter(root, p, doc)
			if err != nil {
				t.Fatalf("ConvertRouter() error = %v", err)
			}
			state, err := data.NewRouterReader().ReadState(content)
			if err != nil {
				t.Fatalf("rendered router did not parse: %v", err)
			}
			if state.Release != p.GoalSelection.GoalID {
				t.Errorf("router Release = %q, want selected live Goal %q", state.Release, p.GoalSelection.GoalID)
			}
			if tc.wantGenerated {
				if !strings.Contains(content, p.GoalSelection.GoalID) || !strings.Contains(content, continuationGoalTitle) {
					t.Errorf("migration note does not explain continuation Goal selection:\n%s", content)
				}
			} else if strings.Contains(content, "## Migration Note") {
				t.Errorf("router has an unnecessary migration note for an existing live Goal:\n%s", content)
			}
			if strings.Contains(content, "router now selects no Release") {
				t.Errorf("router retains obsolete no-Release note:\n%s", content)
			}
		})
	}
}

// --- router: archived selection --------------------------------------------

// TestConvertRouter_archivedTaskSelection_recordsNoteAndNoDanglingReference
// covers a router pointing at a completed (archived) task under an otherwise
// active epic: the Objective still resolves, but the Task does not, and a
// named note explains why instead of a dangling T-### reference.
func TestConvertRouter_archivedTaskSelection_recordsNoteAndNoDanglingReference(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".savepoint", "config.yml"), "quality_gates: {}\n")
	writeFile(t, filepath.Join(root, ".savepoint", "router.md"), routerFixtureContent(
		"task-building", "v1", "E01-x", "E01-x/T001-done", `"Handle E01-x/T001-done."`))
	writeFile(t, filepath.Join(root, ".savepoint", "releases", "v1", "v1-PRD.md"),
		"---\nname: V1\nstatus: in_progress\n---\n\n# V1\n")
	writeFile(t, filepath.Join(root, ".savepoint", "releases", "v1", "epics", "E01-x", "E01-Detail.md"),
		"---\nstatus: in_progress\n---\n\n# E01\n")
	writeFile(t, filepath.Join(root, ".savepoint", "releases", "v1", "epics", "E01-x", "tasks", "T001-done.md"),
		"---\nid: E01-x/T001-done\nstatus: done\n---\n\n# T001\n")

	p := mustPlan(t, root)

	epicTarget, ok := targetByPath(p, ".savepoint/releases/v1/epics/E01-x/E01-Detail.md")
	if !ok {
		t.Fatalf("epic was not planned as an active Objective target")
	}
	taskArchive, ok := archiveByPath(p, ".savepoint/releases/v1/epics/E01-x/tasks/T001-done.md")
	if !ok {
		t.Fatalf("done task was not archived")
	}

	doc, ok := documentByPath(p, ".savepoint/router.md")
	if !ok {
		t.Fatalf("no PlannedDocument for router.md")
	}

	content, err := ConvertRouter(root, p, doc)
	if err != nil {
		t.Fatalf("ConvertRouter error = %v", err)
	}

	if !strings.Contains(content, "objective: "+epicTarget.GlobalID) {
		t.Errorf("content missing resolved objective %q:\n%s", epicTarget.GlobalID, content)
	}
	if !strings.Contains(content, "## Migration Note") {
		t.Errorf("content has no Migration Note for an archived task selection:\n%s", content)
	}
	if !strings.Contains(content, taskArchive.ArchivePath) {
		t.Errorf("Migration Note does not name the archive path %q:\n%s", taskArchive.ArchivePath, content)
	}

	readBack, err := data.NewRouterReader().ReadState(content)
	if err != nil {
		t.Fatalf("rendered router content did not parse back: %v", err)
	}
	if readBack.Task != "" {
		t.Errorf("readBack.Task = %q, want empty: no dangling reference to an archived task", readBack.Task)
	}
}

// TestConvertRouter_archivedEpicSelection_recordsBothNotes covers a router
// pointing at a fully-closed epic (and therefore its task too): neither
// resolves, and both are named in the note rather than silently dropped.
func TestConvertRouter_archivedEpicSelection_recordsBothNotes(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".savepoint", "config.yml"), "quality_gates: {}\n")
	writeFile(t, filepath.Join(root, ".savepoint", "router.md"), routerFixtureContent(
		"audit-pending", "v1", "E02-y", "E02-y/T001-y", `"Audit E02-y."`))
	writeFile(t, filepath.Join(root, ".savepoint", "releases", "v1", "epics", "E02-y", "E02-Detail.md"),
		"---\nstatus: done\n---\n\n# E02\n")
	writeFile(t, filepath.Join(root, ".savepoint", "releases", "v1", "epics", "E02-y", "tasks", "T001-y.md"),
		"---\nid: E02-y/T001-y\nstatus: done\n---\n\n# T001\n")

	p := mustPlan(t, root)

	if _, ok := targetByPath(p, ".savepoint/releases/v1/epics/E02-y/E02-Detail.md"); ok {
		t.Fatalf("done epic was planned as an Objective target, want archive-only")
	}
	epicArchive, ok := archiveByPath(p, ".savepoint/releases/v1/epics/E02-y/E02-Detail.md")
	if !ok {
		t.Fatalf("done epic was not archived")
	}
	taskArchive, ok := archiveByPath(p, ".savepoint/releases/v1/epics/E02-y/tasks/T001-y.md")
	if !ok {
		t.Fatalf("task under the done epic was not archived")
	}

	doc, ok := documentByPath(p, ".savepoint/router.md")
	if !ok {
		t.Fatalf("no PlannedDocument for router.md")
	}

	content, err := ConvertRouter(root, p, doc)
	if err != nil {
		t.Fatalf("ConvertRouter error = %v", err)
	}

	if !strings.Contains(content, epicArchive.ArchivePath) {
		t.Errorf("note missing epic archive path %q:\n%s", epicArchive.ArchivePath, content)
	}
	if !strings.Contains(content, taskArchive.ArchivePath) {
		t.Errorf("note missing task archive path %q:\n%s", taskArchive.ArchivePath, content)
	}

	readBack, err := data.NewRouterReader().ReadState(content)
	if err != nil {
		t.Fatalf("rendered router content did not parse back: %v", err)
	}
	if readBack.State != "check" {
		t.Errorf("readBack.State = %q, want check (audit-pending maps to check)", readBack.State)
	}
	if readBack.Task != "" {
		t.Errorf("readBack.Task = %q, want empty", readBack.Task)
	}
}

// routerFixtureContent renders a minimal router.md matching the shape both
// frozen fixtures use, so ad hoc tests exercise the same anchor-finding path.
func routerFixtureContent(state, release, epic, task, nextAction string) string {
	var b strings.Builder
	b.WriteString("# Agent State Machine\n\n## Read order\n\n1. This file (router.md)\n\n## Current state\n\n```yaml\n")
	b.WriteString("state: " + state + "\n")
	b.WriteString("release: " + release + "\n")
	b.WriteString("epic: " + epic + "\n")
	b.WriteString("task: " + task + "\n")
	b.WriteString("next_action: " + nextAction + "\n")
	b.WriteString("```\n")
	return b.String()
}

// --- determinism ------------------------------------------------------------

// TestConvertDocs_deterministicAcrossRepeatedRuns proves ConvertIdea and
// ConvertRouter render byte-identical content across repeated calls over the
// same plan, matching convert.go's determinism guarantee for records.
func TestConvertDocs_deterministicAcrossRepeatedRuns(t *testing.T) {
	for _, fixture := range []string{"v1-basic", "v1-history"} {
		t.Run(fixture, func(t *testing.T) {
			root := fixtureProjectRoot(fixture)
			p := mustPlan(t, root)

			for _, doc := range p.Documents {
				var first, second string
				var err error
				switch doc.Kind {
				case DocumentIdea:
					first, err = ConvertIdea(root, doc)
					if err == nil {
						second, err = ConvertIdea(root, doc)
					}
				case DocumentRouter:
					first, err = ConvertRouter(root, p, doc)
					if err == nil {
						second, err = ConvertRouter(root, p, doc)
					}
				}
				if err != nil {
					t.Fatalf("convert doc %s error = %v", doc.SourcePath, err)
				}
				if first != second {
					t.Errorf("convert doc %s not deterministic:\nfirst=%q\nsecond=%q", doc.SourcePath, first, second)
				}
			}
		})
	}
}

// --- Health-Check.md: archived, candidate commands never inferred as gates -

func TestCandidateHealthCheckCommands_extractsOnlyFencedLines(t *testing.T) {
	body := "# Health Check\n\n" +
		"Prose mentioning `make test` inline should not count.\n\n" +
		"```sh\nmake build\nmake test\n```\n\n" +
		"More prose.\n\n" +
		"```\nnpm run lint\n```\n"

	got := CandidateHealthCheckCommands(body)
	want := []string{"make build", "make test", "npm run lint"}

	if len(got) != len(want) {
		t.Fatalf("CandidateHealthCheckCommands = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("CandidateHealthCheckCommands[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

// TestPlan_healthCheckArchivedWithCandidateCommandsGatesUnaffected proves
// Health-Check.md is archived with its candidate commands recorded for
// preview, and that a project's config.yml quality_gates are exactly what it
// declared — never anything inferred from Health-Check.md's prose, even when
// that prose is full of command-looking lines.
func TestPlan_healthCheckArchivedWithCandidateCommandsGatesUnaffected(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, ".savepoint", "config.yml")
	writeFile(t, configPath, "quality_gates:\n  test: \"declared test command\"\n")
	writeFile(t, filepath.Join(root, ".savepoint", "router.md"), "# Router\n")
	writeFile(t, filepath.Join(root, ".savepoint", "Health-Check.md"),
		"---\ntype: health-check\n---\n\n# Health Check\n\n```sh\nmake build\nmake test\nnpm run lint\n```\n")

	p := mustPlan(t, root)

	const healthCheckPath = ".savepoint/Health-Check.md"
	archived, ok := archiveByPath(p, healthCheckPath)
	if !ok {
		t.Fatalf("Health-Check.md was not archived")
	}
	if archived.Role != RoleHealthCheck {
		t.Errorf("archived Role = %v, want RoleHealthCheck", archived.Role)
	}
	if archived.ArchivePath != ".savepoint/archive/v1/"+healthCheckPath {
		t.Errorf("archived ArchivePath = %q, want .savepoint/archive/v1/%s", archived.ArchivePath, healthCheckPath)
	}
	wantCommands := []string{"make build", "make test", "npm run lint"}
	if len(archived.CandidateCommands) != len(wantCommands) {
		t.Fatalf("CandidateCommands = %v, want %v", archived.CandidateCommands, wantCommands)
	}
	for i, want := range wantCommands {
		if archived.CandidateCommands[i] != want {
			t.Errorf("CandidateCommands[%d] = %q, want %q", i, archived.CandidateCommands[i], want)
		}
	}

	if _, ok := targetByPath(p, healthCheckPath); ok {
		t.Errorf("Health-Check.md was planned as a target, want archive-only")
	}
	if _, ok := documentByPath(p, healthCheckPath); ok {
		t.Errorf("Health-Check.md was planned as a document, want archive-only")
	}

	cfg, err := data.NewConfigReader().Read(configPath)
	if err != nil {
		t.Fatalf("read config.yml: %v", err)
	}
	if cfg.QualityGates.Test == nil || *cfg.QualityGates.Test != "declared test command" {
		t.Errorf("QualityGates.Test = %v, want \"declared test command\" — never inferred from Health-Check.md prose", cfg.QualityGates.Test)
	}
	if cfg.QualityGates.Lint != nil {
		t.Errorf("QualityGates.Lint = %v, want nil: no gate was declared, none should be inferred from `npm run lint`", cfg.QualityGates.Lint)
	}
}

// TestPlan_v1Basic_absentHealthCheckConvertsWithoutFinding proves a project
// missing the optional Health-Check.md (v1-basic's actual shape) plans
// cleanly: no archive entry, no document, no ambiguity naming it.
func TestPlan_v1Basic_absentHealthCheckConvertsWithoutFinding(t *testing.T) {
	p := mustPlan(t, fixtureProjectRoot("v1-basic"))

	for _, a := range p.Archives {
		if a.Role == RoleHealthCheck {
			t.Errorf("v1-basic has no Health-Check.md, but an archive entry claims one: %+v", a)
		}
	}
	for _, amb := range p.Ambiguities {
		if strings.Contains(amb.Path, "Health-Check") {
			t.Errorf("absent optional Health-Check.md produced an ambiguity: %+v", amb)
		}
	}
}

// --- documents preserved untouched in place ---------------------------------

// TestPlan_configAndInPlaceDocsAreNeverArchivedOrTargeted proves config.yml,
// Design.md, Guardrails.md, and visual-identity.md never appear in
// plan.Targets, plan.Archives, or plan.Documents: omission from all three is
// itself the "preserved in place, untouched" outcome.
func TestPlan_configAndInPlaceDocsAreNeverArchivedOrTargeted(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".savepoint", "config.yml"), "quality_gates: {}\n")
	writeFile(t, filepath.Join(root, ".savepoint", "router.md"), "# Router\n")
	writeFile(t, filepath.Join(root, ".savepoint", "Design.md"), "# Design\n\nArchitecture.\n")
	writeFile(t, filepath.Join(root, ".savepoint", "Guardrails.md"), "# Guardrails\n\nSTYLE-01: ...\n")
	writeFile(t, filepath.Join(root, ".savepoint", "visual-identity.md"), "# Visual Identity\n\nPalette.\n")

	p := mustPlan(t, root)

	untouched := []string{
		".savepoint/config.yml",
		".savepoint/Design.md",
		".savepoint/Guardrails.md",
		".savepoint/visual-identity.md",
	}
	for _, path := range untouched {
		if _, ok := targetByPath(p, path); ok {
			t.Errorf("%s was planned as a target, want untouched", path)
		}
		if _, ok := archiveByPath(p, path); ok {
			t.Errorf("%s was archived, want untouched", path)
		}
		if _, ok := documentByPath(p, path); ok {
			t.Errorf("%s was planned as a document, want untouched", path)
		}
	}

	for _, path := range untouched {
		before, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		after, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
		if err != nil {
			t.Fatalf("re-read %s: %v", path, err)
		}
		if string(before) != string(after) {
			t.Errorf("%s bytes changed", path)
		}
	}
}

// TestPlan_missingOptionalDocsConvertWithoutFinding proves a project missing
// Guardrails.md and visual-identity.md entirely (v1-basic's actual shape,
// since neither fixture carries them) plans with no ambiguity naming either.
func TestPlan_missingOptionalDocsConvertWithoutFinding(t *testing.T) {
	for _, fixture := range []string{"v1-basic", "v1-history"} {
		p := mustPlan(t, fixtureProjectRoot(fixture))
		for _, amb := range p.Ambiguities {
			if strings.Contains(amb.Path, "Guardrails") || strings.Contains(amb.Path, "visual-identity") {
				t.Errorf("%s: absent optional doc produced an ambiguity: %+v", fixture, amb)
			}
		}
	}
}
