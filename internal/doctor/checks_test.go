package doctor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/migrate"
	"github.com/opencode/savepoint/internal/testutil"
)

// --- CheckConfig ---

func TestCheckConfigMissing(t *testing.T) {
	root := t.TempDir()
	err := CheckConfig(root)
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Fatalf("CheckConfig() = %v, want not found error", err)
	}
}

func TestCheckStructureReportsAgentCompleteStatusAlias(t *testing.T) {
	root := t.TempDir()
	testutil.SetupMinimalProject(t, root, "v1", "E01-foo")
	testutil.WriteTask(t, root, "v1", "E01-foo", testutil.TaskFixture{
		Slug:      "T001-agent-complete",
		Status:    "complete",
		Objective: "Agent wrote complete",
	})

	problems := CheckStructure(root, "")
	found := false
	for _, p := range problems {
		if strings.Contains(p.Message, `non-canonical status "complete"`) && strings.Contains(p.Message, `"done"`) {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("CheckStructure() = %v, want non-canonical complete status problem", problems)
	}
}

func TestCheckStructureReportsOverlongComplexityReasonWords(t *testing.T) {
	root := t.TempDir()
	testutil.SetupMinimalProject(t, root, "v1", "E01-foo")
	testutil.WriteTask(t, root, "v1", "E01-foo", testutil.TaskFixture{
		Slug:      "T001-long-complexity",
		Status:    "planned",
		Objective: "Long complexity",
		Extra: map[string]string{
			"complexity_tier":   "high",
			"complexity_reason": `"` + strings.TrimSpace(strings.Repeat("word ", data.MaxComplexityReasonWords+1)) + `"`,
		},
	})

	problems := CheckStructure(root, "")
	found := false
	for _, p := range problems {
		if strings.Contains(p.Message, "task complexity invalid") && strings.Contains(p.Message, "maximum is") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("CheckStructure() = %v, want overlong complexity reason problem", problems)
	}
}

func TestCheckConfigInvalidYAML(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "theme: [broken")
	err := CheckConfig(root)
	if err == nil || !strings.Contains(err.Error(), "invalid YAML") {
		t.Fatalf("CheckConfig() = %v, want invalid YAML error", err)
	}
}

func TestCheckConfigMissingQualityGates(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "theme:\n  bg: \"#000\"\n")
	err := CheckConfig(root)
	if err == nil || !strings.Contains(err.Error(), "quality_gates") {
		t.Fatalf("CheckConfig() = %v, want quality_gates error", err)
	}
}

func TestCheckConfigMissingTheme(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "quality_gates:\n  block_on_failure: true\n")
	err := CheckConfig(root)
	if err == nil || !strings.Contains(err.Error(), "theme") {
		t.Fatalf("CheckConfig() = %v, want theme error", err)
	}
}

func TestCheckConfigValid(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "quality_gates:\n  block_on_failure: true\ntheme:\n  bg: \"#000\"\n")
	if err := CheckConfig(root); err != nil {
		t.Fatalf("CheckConfig() = %v, want nil", err)
	}
}

// --- CheckMigration ---

func TestCheckMigration_noneExists(t *testing.T) {
	root := t.TempDir()
	if problems := CheckMigration(root); len(problems) != 0 {
		t.Fatalf("CheckMigration() = %v, want no problems when nothing is pending", problems)
	}
}

func TestCheckMigration_incompleteOperationIsNamedDiagnostic(t *testing.T) {
	projectDir := t.TempDir()
	root := filepath.Join(projectDir, ".savepoint")
	if err := os.MkdirAll(root, 0755); err != nil {
		t.Fatal(err)
	}
	entries := []migrate.JournalEntry{{Path: "objectives/O001.md", Action: migrate.ActionCreate}}
	if _, err := migrate.CreateOperation(projectDir, "op-1", nil, entries, time.Now()); err != nil {
		t.Fatalf("CreateOperation() error = %v", err)
	}
	problems := CheckMigration(root)
	if len(problems) != 1 {
		t.Fatalf("CheckMigration() = %v, want exactly one problem", problems)
	}
	p := problems[0]
	if !strings.Contains(p.Message, "migrate-operation-incomplete") {
		t.Errorf("Message = %q, want the migrate-operation-incomplete diagnostic name", p.Message)
	}
	if !strings.Contains(p.Message, "op-1") {
		t.Errorf("Message = %q, want it to name the operation", p.Message)
	}
	if !strings.Contains(p.Message, "--recover") {
		t.Errorf("Message = %q, want the recovery command", p.Message)
	}
	if !strings.Contains(p.Repair, "migrate --recover") {
		t.Errorf("Repair = %q, want the recovery command and no write of its own", p.Repair)
	}
}

func TestCheckMigration_multipleOperationsIsNamedDiagnostic(t *testing.T) {
	projectDir := t.TempDir()
	root := filepath.Join(projectDir, ".savepoint")
	if err := os.MkdirAll(root, 0755); err != nil {
		t.Fatal(err)
	}
	if _, err := migrate.CreateOperation(projectDir, "op-a", nil, nil, time.Now()); err != nil {
		t.Fatalf("CreateOperation(op-a) error = %v", err)
	}
	if _, err := migrate.CreateOperation(projectDir, "op-b", nil, nil, time.Now()); err != nil {
		t.Fatalf("CreateOperation(op-b) error = %v", err)
	}
	problems := CheckMigration(root)
	if len(problems) != 1 {
		t.Fatalf("CheckMigration() = %v, want exactly one problem", problems)
	}
	if !strings.Contains(problems[0].Message, "migrate-multiple-operations") {
		t.Errorf("Message = %q, want the migrate-multiple-operations diagnostic name", problems[0].Message)
	}
}

// --- CheckRouter ---

func TestCheckRouterMissing(t *testing.T) {
	root := t.TempDir()
	err := CheckRouter(root, "")
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Fatalf("CheckRouter() = %v, want not found error", err)
	}
}

func TestCheckRouterInvalidStateBlock(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "router.md"), "# no state block")
	err := CheckRouter(root, "")
	if err == nil || !strings.Contains(err.Error(), "invalid state block") {
		t.Fatalf("CheckRouter() = %v, want invalid state block error", err)
	}
}

func TestCheckRouterPreImplementation(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "router.md"), routerContent("pre-implementation", "none", "none"))
	if err := CheckRouter(root, ""); err != nil {
		t.Fatalf("CheckRouter() = %v, want nil", err)
	}
}

func TestCheckRouterMissingReleaseDir(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "router.md"), routerContent("task-building", "v1", "none"))
	err := CheckRouter(root, "")
	if err == nil || !strings.Contains(err.Error(), "release") {
		t.Fatalf("CheckRouter() = %v, want release directory error", err)
	}
}

func TestCheckRouterMissingEpicDir(t *testing.T) {
	root := t.TempDir()
	testutil.MkdirAll(t, filepath.Join(root, "releases", "v1"))
	testutil.WriteFile(t, filepath.Join(root, "router.md"), routerContent("task-building", "v1", "E03-foo"))
	err := CheckRouter(root, "")
	if err == nil || !strings.Contains(err.Error(), "epic") {
		t.Fatalf("CheckRouter() = %v, want epic directory error", err)
	}
}

func TestCheckRouterValidWithDirs(t *testing.T) {
	root := t.TempDir()
	testutil.MkdirAll(t, filepath.Join(root, "releases", "v1", "epics", "E03-foo"))
	testutil.WriteFile(t, filepath.Join(root, "router.md"), routerContent("task-building", "v1", "E03-foo"))
	if err := CheckRouter(root, ""); err != nil {
		t.Fatalf("CheckRouter() = %v, want nil", err)
	}
}

func TestCheckRouterEpicFilterSkip(t *testing.T) {
	root := t.TempDir()
	// release dir missing — would fail without filter
	testutil.WriteFile(t, filepath.Join(root, "router.md"), routerContent("task-building", "v1", "E03-foo"))
	// filter doesn't match router epic → skip dir checks
	if err := CheckRouter(root, "E99-other"); err != nil {
		t.Fatalf("CheckRouter() = %v, want nil (filter skip)", err)
	}
}

// --- CheckStructure ---

func TestCheckStructure_MissingReleasesDir(t *testing.T) {
	root := t.TempDir()
	problems := CheckStructure(root, "")
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "releases directory not found") {
		t.Fatalf("CheckStructure() = %v, want releases directory error", problems)
	}
}

func TestCheckStructure_EmptyReleases(t *testing.T) {
	root := t.TempDir()
	testutil.MkdirAll(t, filepath.Join(root, "releases"))
	problems := CheckStructure(root, "")
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "no release directories found") {
		t.Fatalf("CheckStructure() = %v, want no releases error", problems)
	}
}

func TestCheckStructure_MissingReleasePRD(t *testing.T) {
	root := t.TempDir()
	testutil.MkdirAll(t, filepath.Join(root, "releases", "v1", "epics"))
	problems := CheckStructure(root, "")
	found := false
	for _, p := range problems {
		if strings.Contains(p.Message, "release PRD file not found") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("CheckStructure() = %v, want release PRD file not found problem", problems)
	}
}

func TestCheckStructure_ReleasePRDValid(t *testing.T) {
	root := t.TempDir()
	testutil.MkdirAll(t, filepath.Join(root, "releases", "v1", "epics"))
	testutil.WriteReleasePRD(t, filepath.Join(root, "releases", "v1"))
	problems := CheckStructure(root, "")
	for _, p := range problems {
		if strings.Contains(p.File, "v1-PRD.md") {
			t.Fatalf("CheckStructure() unexpected PRD problem: %v", p)
		}
	}
}

func TestCheckStructure_ReleasePRDCorruptYAML(t *testing.T) {
	root := t.TempDir()
	testutil.MkdirAll(t, filepath.Join(root, "releases", "v1", "epics"))
	testutil.WriteFile(t, filepath.Join(root, "releases", "v1", "v1-PRD.md"), "---\ntype: [broken\n---\n")
	problems := CheckStructure(root, "")
	found := false
	for _, p := range problems {
		if strings.Contains(p.File, "v1-PRD.md") && p.Line > 0 {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("CheckStructure() = %v, want corrupt YAML with line in v1-PRD.md", problems)
	}
}

func TestCheckStructure_ValidEpicDetail(t *testing.T) {
	root := t.TempDir()
	releasePath := filepath.Join(root, "releases", "v1")
	epicPath := filepath.Join(releasePath, "epics", "E01-foo")
	testutil.MkdirAll(t, epicPath)
	testutil.WriteReleasePRD(t, releasePath)
	testutil.WriteFile(t, filepath.Join(epicPath, "E01-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E01: Foo\n")
	problems := CheckStructure(root, "")
	for _, p := range problems {
		if strings.Contains(p.File, "Detail.md") {
			t.Fatalf("CheckStructure() unexpected Detail.md problem: %v", p)
		}
	}
}

func TestCheckStructure_MissingEpicDetail(t *testing.T) {
	root := t.TempDir()
	releasePath := filepath.Join(root, "releases", "v1")
	epicPath := filepath.Join(releasePath, "epics", "E01-foo")
	testutil.MkdirAll(t, epicPath)
	testutil.WriteReleasePRD(t, releasePath)
	problems := CheckStructure(root, "")
	found := false
	for _, p := range problems {
		if strings.Contains(p.Message, "epic detail file not found") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("CheckStructure() = %v, want epic detail file not found problem", problems)
	}
}

func TestCheckStructure_EpicCanonicalStatusNoProblem(t *testing.T) {
	root := t.TempDir()
	releasePath := filepath.Join(root, "releases", "v1")
	epicPath := filepath.Join(releasePath, "epics", "E01-foo")
	testutil.MkdirAll(t, epicPath)
	testutil.WriteReleasePRD(t, releasePath)
	testutil.WriteFile(t, filepath.Join(epicPath, "E01-Detail.md"), "---\ntype: epic-design\nstatus: audited\n---\n\n# E01: Foo\n")
	problems := CheckStructure(root, "")
	for _, p := range problems {
		if strings.Contains(p.Message, "epic") && strings.Contains(p.Message, "status") {
			t.Fatalf("CheckStructure() unexpected epic status problem: %v", p)
		}
	}
}

func TestCheckStructure_EpicNonCanonicalStatusReported(t *testing.T) {
	root := t.TempDir()
	releasePath := filepath.Join(root, "releases", "v1")
	epicPath := filepath.Join(releasePath, "epics", "E01-foo")
	testutil.MkdirAll(t, epicPath)
	testutil.WriteReleasePRD(t, releasePath)
	testutil.WriteFile(t, filepath.Join(epicPath, "E01-Detail.md"), "---\ntype: epic-design\nstatus: epic-design\n---\n\n# E01: Foo\n")
	problems := CheckStructure(root, "")
	found := false
	for _, p := range problems {
		if strings.Contains(p.Message, `epic uses non-canonical status "epic-design"`) {
			if SuggestRepair(p) != "Set the epic status to planned, in_progress, done, or audited" {
				t.Fatalf("SuggestRepair(%q) = %q, want epic status hint", p.Message, SuggestRepair(p))
			}
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("CheckStructure() = %v, want non-canonical epic status problem", problems)
	}
}

func TestCheckStructure_ValidTask(t *testing.T) {
	root := t.TempDir()
	releasePath := filepath.Join(root, "releases", "v1")
	epicPath := filepath.Join(releasePath, "epics", "E01-foo")
	tasksPath := filepath.Join(epicPath, "tasks")
	testutil.MkdirAll(t, tasksPath)
	testutil.WriteReleasePRD(t, releasePath)
	testutil.WriteFile(t, filepath.Join(epicPath, "E01-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E01: Foo\n")
	testutil.WriteFile(t, filepath.Join(tasksPath, "T001-task.md"), "---\nid: E01-foo/T001-task\nstatus: planned\nobjective: \"Do the thing\"\ndepends_on: []\n---\n\n# T001: Task\n\n## Acceptance Criteria\n\n- It works\n")
	problems := CheckStructure(root, "")
	if len(problems) > 0 {
		t.Fatalf("CheckStructure() = %v, want no problems", problems)
	}
}

func TestCheckStructure_TaskInProgressMissingStage(t *testing.T) {
	root := t.TempDir()
	releasePath := filepath.Join(root, "releases", "v1")
	epicPath := filepath.Join(releasePath, "epics", "E01-foo")
	tasksPath := filepath.Join(epicPath, "tasks")
	testutil.MkdirAll(t, tasksPath)
	testutil.WriteReleasePRD(t, releasePath)
	testutil.WriteFile(t, filepath.Join(epicPath, "E01-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E01: Foo\n")
	testutil.WriteFile(t, filepath.Join(tasksPath, "T001-task.md"), "---\nid: E01-foo/T001-task\nstatus: in_progress\nobjective: \"Do the thing\"\ndepends_on: []\n---\n\n# T001: Task\n\n## Acceptance Criteria\n\n- It works\n")

	problems := CheckStructure(root, "")
	found := false
	for _, p := range problems {
		if strings.Contains(p.Message, "task stage is required") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("CheckStructure() = %v, want missing task stage problem", problems)
	}
}

func TestCheckStructure_TaskInvalidStage(t *testing.T) {
	root := t.TempDir()
	releasePath := filepath.Join(root, "releases", "v1")
	epicPath := filepath.Join(releasePath, "epics", "E01-foo")
	tasksPath := filepath.Join(epicPath, "tasks")
	testutil.MkdirAll(t, tasksPath)
	testutil.WriteReleasePRD(t, releasePath)
	testutil.WriteFile(t, filepath.Join(epicPath, "E01-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E01: Foo\n")
	testutil.WriteFile(t, filepath.Join(tasksPath, "T001-task.md"), "---\nid: E01-foo/T001-task\nstatus: in_progress\nstage: done\nobjective: \"Do the thing\"\ndepends_on: []\n---\n\n# T001: Task\n\n## Acceptance Criteria\n\n- It works\n")

	problems := CheckStructure(root, "")
	found := false
	for _, p := range problems {
		if strings.Contains(p.Message, "task stage invalid") && strings.Contains(p.Message, "done") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("CheckStructure() = %v, want invalid task stage problem", problems)
	}
}

func TestCheckStructure_TaskStageOutsideInProgress(t *testing.T) {
	root := t.TempDir()
	releasePath := filepath.Join(root, "releases", "v1")
	epicPath := filepath.Join(releasePath, "epics", "E01-foo")
	tasksPath := filepath.Join(epicPath, "tasks")
	testutil.MkdirAll(t, tasksPath)
	testutil.WriteReleasePRD(t, releasePath)
	testutil.WriteFile(t, filepath.Join(epicPath, "E01-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E01: Foo\n")
	testutil.WriteFile(t, filepath.Join(tasksPath, "T001-task.md"), "---\nid: E01-foo/T001-task\nstatus: planned\nstage: build\nobjective: \"Do the thing\"\ndepends_on: []\n---\n\n# T001: Task\n\n## Acceptance Criteria\n\n- It works\n")

	problems := CheckStructure(root, "")
	found := false
	for _, p := range problems {
		if strings.Contains(p.Message, "task stage field") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("CheckStructure() = %v, want stage outside in_progress problem", problems)
	}
}

func TestCheckStructure_TaskImplementationStageOutsideInProgress(t *testing.T) {
	root := t.TempDir()
	releasePath := filepath.Join(root, "releases", "v1")
	epicPath := filepath.Join(releasePath, "epics", "E01-foo")
	tasksPath := filepath.Join(epicPath, "tasks")
	testutil.MkdirAll(t, tasksPath)
	testutil.WriteReleasePRD(t, releasePath)
	testutil.WriteFile(t, filepath.Join(epicPath, "E01-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E01: Foo\n")
	testutil.WriteFile(t, filepath.Join(tasksPath, "T001-task.md"), "---\nid: E01-foo/T001-task\nstatus: planned\nstage: implementation\nobjective: \"Do the thing\"\ndepends_on: []\n---\n\n# T001: Task\n\n## Acceptance Criteria\n\n- It works\n")

	problems := CheckStructure(root, "")
	found := false
	for _, p := range problems {
		if strings.Contains(p.Message, "task stage field") && strings.Contains(p.Message, "implementation") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("CheckStructure() = %v, want implementation stage outside in_progress problem", problems)
	}
}

func TestCheckStructure_TaskImplementationPhaseOutsideInProgress(t *testing.T) {
	root := t.TempDir()
	releasePath := filepath.Join(root, "releases", "v1")
	epicPath := filepath.Join(releasePath, "epics", "E01-foo")
	tasksPath := filepath.Join(epicPath, "tasks")
	testutil.MkdirAll(t, tasksPath)
	testutil.WriteReleasePRD(t, releasePath)
	testutil.WriteFile(t, filepath.Join(epicPath, "E01-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E01: Foo\n")
	testutil.WriteFile(t, filepath.Join(tasksPath, "T001-task.md"), "---\nid: E01-foo/T001-task\nstatus: done\nphase: implementation\nobjective: \"Do the thing\"\ndepends_on: []\n---\n\n# T001: Task\n\n## Acceptance Criteria\n\n- It works\n")

	problems := CheckStructure(root, "")
	found := false
	for _, p := range problems {
		if strings.Contains(p.Message, "legacy frontmatter field phase") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("CheckStructure() = %v, want legacy implementation phase problem", problems)
	}
}

func TestCheckStructure_TaskLegacyPhaseField(t *testing.T) {
	root := t.TempDir()
	releasePath := filepath.Join(root, "releases", "v1")
	epicPath := filepath.Join(releasePath, "epics", "E01-foo")
	tasksPath := filepath.Join(epicPath, "tasks")
	testutil.MkdirAll(t, tasksPath)
	testutil.WriteReleasePRD(t, releasePath)
	testutil.WriteFile(t, filepath.Join(epicPath, "E01-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E01: Foo\n")
	testutil.WriteFile(t, filepath.Join(tasksPath, "T001-task.md"), "---\nid: E01-foo/T001-task\nstatus: in_progress\nphase: build\nobjective: \"Do the thing\"\ndepends_on: []\n---\n\n# T001: Task\n\n## Acceptance Criteria\n\n- It works\n")

	problems := CheckStructure(root, "")
	found := false
	for _, p := range problems {
		if strings.Contains(p.Message, "legacy frontmatter field phase") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("CheckStructure() = %v, want legacy phase problem", problems)
	}
}

func TestCheckStructure_TaskLegacyPhaseDoneReportsInvalidStage(t *testing.T) {
	root := t.TempDir()
	releasePath := filepath.Join(root, "releases", "v1")
	epicPath := filepath.Join(releasePath, "epics", "E01-foo")
	tasksPath := filepath.Join(epicPath, "tasks")
	testutil.MkdirAll(t, tasksPath)
	testutil.WriteReleasePRD(t, releasePath)
	testutil.WriteFile(t, filepath.Join(epicPath, "E01-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E01: Foo\n")
	testutil.WriteFile(t, filepath.Join(tasksPath, "T001-task.md"), "---\nid: E01-foo/T001-task\nstatus: in_progress\nphase: done\nobjective: \"Do the thing\"\ndepends_on: []\n---\n\n# T001: Task\n\n## Acceptance Criteria\n\n- It works\n")

	problems := CheckStructure(root, "")
	found := false
	for _, p := range problems {
		if strings.Contains(p.Message, "task stage invalid") && strings.Contains(p.Message, "done") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("CheckStructure() = %v, want invalid legacy phase stage problem", problems)
	}
}

func TestCheckStructure_TaskLegacyImplementationPhaseReportsInvalidStage(t *testing.T) {
	root := t.TempDir()
	releasePath := filepath.Join(root, "releases", "v1")
	epicPath := filepath.Join(releasePath, "epics", "E01-foo")
	tasksPath := filepath.Join(epicPath, "tasks")
	testutil.MkdirAll(t, tasksPath)
	testutil.WriteReleasePRD(t, releasePath)
	testutil.WriteFile(t, filepath.Join(epicPath, "E01-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E01: Foo\n")
	testutil.WriteFile(t, filepath.Join(tasksPath, "T001-task.md"), "---\nid: E01-foo/T001-task\nstatus: in_progress\nphase: implementation\nobjective: \"Do the thing\"\ndepends_on: []\n---\n\n# T001: Task\n\n## Acceptance Criteria\n\n- It works\n")

	problems := CheckStructure(root, "")
	found := false
	for _, p := range problems {
		if strings.Contains(p.Message, "task stage invalid") && strings.Contains(p.Message, "implementation") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("CheckStructure() = %v, want invalid legacy implementation phase problem", problems)
	}
}

func TestCheckStructure_TaskLifecycleDiagnosticsAreReadOnly(t *testing.T) {
	root := t.TempDir()
	releasePath := filepath.Join(root, "releases", "v1")
	epicPath := filepath.Join(releasePath, "epics", "E01-foo")
	tasksPath := filepath.Join(epicPath, "tasks")
	taskPath := filepath.Join(tasksPath, "T001-task.md")
	testutil.MkdirAll(t, tasksPath)
	testutil.WriteReleasePRD(t, releasePath)
	testutil.WriteFile(t, filepath.Join(epicPath, "E01-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E01: Foo\n")
	content := "---\nid: E01-foo/T001-task\nstatus: done\nstage: implementation\nphase: implementation\nobjective: \"Do the thing\"\ndepends_on: []\n---\n\n# T001: Task\n\n## Acceptance Criteria\n\n- It works\n"
	testutil.WriteFile(t, taskPath, content)

	problems := CheckStructure(root, "")
	if len(problems) == 0 {
		t.Fatal("CheckStructure() expected lifecycle problems")
	}

	after, err := os.ReadFile(taskPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != content {
		t.Fatalf("CheckStructure() changed task content:\n%s", string(after))
	}
}

func TestCheckStructure_TaskMissingRequiredField(t *testing.T) {
	root := t.TempDir()
	releasePath := filepath.Join(root, "releases", "v1")
	epicPath := filepath.Join(releasePath, "epics", "E01-foo")
	tasksPath := filepath.Join(epicPath, "tasks")
	testutil.MkdirAll(t, tasksPath)
	testutil.WriteReleasePRD(t, releasePath)
	testutil.WriteFile(t, filepath.Join(epicPath, "E01-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E01: Foo\n")
	testutil.WriteFile(t, filepath.Join(tasksPath, "T001-task.md"), "---\nid: E01-foo/T001-task\nobjective: \"Do the thing\"\n---\n\n# T001: Task\n\n## Acceptance Criteria\n\n- It works\n")
	problems := CheckStructure(root, "")
	found := false
	for _, p := range problems {
		if strings.Contains(p.Message, "status") && strings.Contains(p.Message, "missing") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("CheckStructure() = %v, want missing status field problem", problems)
	}
}

func TestCheckStructure_TaskMissingAcceptanceCriteria(t *testing.T) {
	root := t.TempDir()
	releasePath := filepath.Join(root, "releases", "v1")
	epicPath := filepath.Join(releasePath, "epics", "E01-foo")
	tasksPath := filepath.Join(epicPath, "tasks")
	testutil.MkdirAll(t, tasksPath)
	testutil.WriteReleasePRD(t, releasePath)
	testutil.WriteFile(t, filepath.Join(epicPath, "E01-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E01: Foo\n")
	testutil.WriteFile(t, filepath.Join(tasksPath, "T001-task.md"), "---\nid: E01-foo/T001-task\nstatus: planned\nobjective: \"Do the thing\"\n---\n\n# T001: Task\n")
	problems := CheckStructure(root, "")
	found := false
	for _, p := range problems {
		if strings.Contains(p.Message, "Acceptance Criteria") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("CheckStructure() = %v, want missing acceptance criteria problem", problems)
	}
}

func TestCheckStructure_TaskCorruptYAML(t *testing.T) {
	root := t.TempDir()
	releasePath := filepath.Join(root, "releases", "v1")
	epicPath := filepath.Join(releasePath, "epics", "E01-foo")
	tasksPath := filepath.Join(epicPath, "tasks")
	testutil.MkdirAll(t, tasksPath)
	testutil.WriteReleasePRD(t, releasePath)
	testutil.WriteFile(t, filepath.Join(epicPath, "E01-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E01: Foo\n")
	testutil.WriteFile(t, filepath.Join(tasksPath, "T001-task.md"), "---\nid: \"unclosed\nstatus: planned\n---\n")
	problems := CheckStructure(root, "")
	foundLine := false
	for _, p := range problems {
		if strings.Contains(p.File, "T001-task.md") && p.Line > 0 {
			foundLine = true
			break
		}
	}
	if !foundLine {
		t.Fatalf("CheckStructure() = %v, want corrupt YAML with line number in task", problems)
	}
}

func TestCheckStructure_EpicFilter(t *testing.T) {
	root := t.TempDir()
	releasePath := filepath.Join(root, "releases", "v1")
	epic1Path := filepath.Join(releasePath, "epics", "E01-foo")
	epic2Path := filepath.Join(releasePath, "epics", "E02-bar")
	testutil.MkdirAll(t, epic1Path)
	testutil.MkdirAll(t, epic2Path)
	testutil.WriteReleasePRD(t, releasePath)
	testutil.WriteFile(t, filepath.Join(epic1Path, "E01-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E01: Foo\n")
	// E02 has no detail file — should not appear when filtering to E01
	problems := CheckStructure(root, "E01-foo")
	for _, p := range problems {
		if strings.Contains(p.Message, "E02") {
			t.Fatalf("CheckStructure() with epicFilter=E01-foo should skip E02, got: %v", p)
		}
	}
}

func TestCheckStructure_EpicFilterByPrefix(t *testing.T) {
	root := t.TempDir()
	releasePath := filepath.Join(root, "releases", "v1")
	epicPath := filepath.Join(releasePath, "epics", "E01-foo")
	testutil.MkdirAll(t, epicPath)
	testutil.WriteReleasePRD(t, releasePath)
	testutil.WriteFile(t, filepath.Join(epicPath, "E01-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E01: Foo\n")
	testutil.MkdirAll(t, filepath.Join(epicPath, "tasks"))
	problems := CheckStructure(root, "E01")
	if len(problems) > 0 {
		t.Fatalf("CheckStructure() with epicFilter=E01 prefix = %v, want no problems", problems)
	}
}

// --- CheckDependencies ---

func TestCheckDependencies_NoReleases(t *testing.T) {
	root := t.TempDir()
	problems := CheckDependencies(root, "")
	if len(problems) == 0 {
		t.Fatal("CheckDependencies() = no problems, want error about releases")
	}
}

func TestCheckDependencies_NoDeps(t *testing.T) {
	root := t.TempDir()
	setupMinimalProject(t, root, "v1", "E01-foo", nil)
	problems := CheckDependencies(root, "")
	if len(problems) > 0 {
		t.Fatalf("CheckDependencies() = %v, want no problems", problems)
	}
}

func TestCheckDependencies_ValidDeps(t *testing.T) {
	root := t.TempDir()
	setupMinimalProject(t, root, "v1", "E01-foo", []taskSpec{
		{id: "E01-foo/T001-task", deps: []string{}},
		{id: "E01-foo/T002-task", deps: []string{"E01-foo/T001-task"}},
		{id: "E01-foo/T003-task", deps: []string{"E01-foo/T002-task"}},
	})
	problems := CheckDependencies(root, "")
	if len(problems) > 0 {
		t.Fatalf("CheckDependencies() = %v, want no problems", problems)
	}
}

func TestCheckDependencies_ShortSameEpicTaskDepAccepted(t *testing.T) {
	root := t.TempDir()
	setupMinimalProject(t, root, "v1", "E01-foo", []taskSpec{
		{id: "E01-foo/T001-task", deps: []string{}},
		{id: "E01-foo/T002-task", deps: []string{"T001"}},
	})
	problems := CheckDependencies(root, "")
	if len(problems) > 0 {
		t.Fatalf("CheckDependencies() = %v, want no problems for short same-epic task dep", problems)
	}
}

func TestCheckDependencies_FilenameStyleSameEpicTaskDepAccepted(t *testing.T) {
	root := t.TempDir()
	setupMinimalProject(t, root, "v1", "E01-foo", []taskSpec{
		{id: "E01-foo/T001-task", deps: []string{}},
		{id: "E01-foo/T002-task", deps: []string{"T001-task"}},
	})
	problems := CheckDependencies(root, "")
	if len(problems) > 0 {
		t.Fatalf("CheckDependencies() = %v, want no problems for filename-style same-epic task dep", problems)
	}
}

func TestCheckDependencies_ShortTaskDepMissingOutsideSameEpic(t *testing.T) {
	root := t.TempDir()
	setupMinimalProject(t, root, "v1", "E01-foo", []taskSpec{
		{id: "E01-foo/T001-task", deps: []string{"T009"}},
	})
	setupMinimalProject(t, root, "v1", "E02-bar", []taskSpec{
		{id: "E02-bar/T009-task", deps: []string{}},
	})
	problems := CheckDependencies(root, "")
	found := false
	for _, p := range problems {
		if strings.Contains(p.Message, "T009") && strings.Contains(p.Message, "non-existent") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("CheckDependencies() = %v, want missing same-epic short task dependency problem", problems)
	}
}

func TestCheckDependencies_ShortEpicDepAccepted(t *testing.T) {
	root := t.TempDir()
	setupMinimalProject(t, root, "v1", "E03-canvas-baseline", nil)
	setupMinimalProject(t, root, "v1", "E06-canvas-polish", []taskSpec{
		{id: "E06-canvas-polish/T001-task", deps: []string{"E03"}},
	})
	problems := CheckDependencies(root, "")
	if len(problems) > 0 {
		t.Fatalf("CheckDependencies() = %v, want no problems for short epic dep", problems)
	}
}

func TestCheckDependencies_FullEpicDepAccepted(t *testing.T) {
	root := t.TempDir()
	setupMinimalProject(t, root, "v1", "E03-canvas-baseline", nil)
	setupMinimalProject(t, root, "v1", "E06-canvas-polish", []taskSpec{
		{id: "E06-canvas-polish/T001-task", deps: []string{"E03-canvas-baseline"}},
	})
	problems := CheckDependencies(root, "")
	if len(problems) > 0 {
		t.Fatalf("CheckDependencies() = %v, want no problems for full epic dep", problems)
	}
}

func TestCheckDependencies_EpicDepMissing(t *testing.T) {
	root := t.TempDir()
	setupMinimalProject(t, root, "v1", "E06-canvas-polish", []taskSpec{
		{id: "E06-canvas-polish/T001-task", deps: []string{"E03"}},
	})
	problems := CheckDependencies(root, "")
	found := false
	for _, p := range problems {
		if strings.Contains(p.Message, "E03") && strings.Contains(p.Message, "non-existent") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("CheckDependencies() = %v, want missing epic dependency problem", problems)
	}
}

func TestCheckDependencies_MissingDep(t *testing.T) {
	root := t.TempDir()
	setupMinimalProject(t, root, "v1", "E01-foo", []taskSpec{
		{id: "E01-foo/T001-task", deps: []string{}},
		{id: "E01-foo/T002-task", deps: []string{"E01-foo/T999-nonexistent"}},
	})
	problems := CheckDependencies(root, "")
	found := false
	for _, p := range problems {
		if strings.Contains(p.Message, "non-existent") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("CheckDependencies() = %v, want missing dependency problem", problems)
	}
}

func TestCheckDependencies_DuplicateIDs(t *testing.T) {
	root := t.TempDir()
	releasePath := filepath.Join(root, "releases", "v1")
	epicPath := filepath.Join(releasePath, "epics", "E01-foo")
	tasksPath := filepath.Join(epicPath, "tasks")
	testutil.MkdirAll(t, tasksPath)
	testutil.WriteReleasePRD(t, releasePath)
	testutil.WriteFile(t, filepath.Join(epicPath, "E01-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E01: Foo\n")

	// Two epics, same task ID
	epic2Path := filepath.Join(releasePath, "epics", "E02-bar")
	tasks2Path := filepath.Join(epic2Path, "tasks")
	testutil.MkdirAll(t, tasks2Path)
	testutil.WriteFile(t, filepath.Join(epic2Path, "E02-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E02: Bar\n")

	testutil.WriteFile(t, filepath.Join(tasksPath, "T001-task.md"), "---\nid: E01-foo/T001-task\nstatus: planned\nobjective: \"A\"\ndepends_on: []\n---\n\n# T001\n\n## Acceptance Criteria\n\n- it works\n")
	testutil.WriteFile(t, filepath.Join(tasks2Path, "T001-task.md"), "---\nid: E01-foo/T001-task\nstatus: planned\nobjective: \"A\"\ndepends_on: []\n---\n\n# T001\n\n## Acceptance Criteria\n\n- it works\n")

	problems := CheckDependencies(root, "")
	found := false
	for _, p := range problems {
		if strings.Contains(p.Message, "duplicate task ID") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("CheckDependencies() = %v, want duplicate task ID problem", problems)
	}
}

func TestCheckDependencies_Cycle(t *testing.T) {
	root := t.TempDir()
	setupMinimalProject(t, root, "v1", "E01-foo", []taskSpec{
		{id: "E01-foo/T001-task", deps: []string{"E01-foo/T003-task"}},
		{id: "E01-foo/T002-task", deps: []string{"E01-foo/T001-task"}},
		{id: "E01-foo/T003-task", deps: []string{"E01-foo/T002-task"}},
	})
	problems := CheckDependencies(root, "")
	found := false
	for _, p := range problems {
		if strings.Contains(p.Message, "cycle") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("CheckDependencies() = %v, want cycle problem", problems)
	}
}

func TestCheckDependencies_CycleAccuratePath(t *testing.T) {
	root := t.TempDir()
	setupMinimalProject(t, root, "v1", "E01-foo", []taskSpec{
		{id: "E01-foo/T001-task", deps: []string{"E01-foo/T002-task"}},
		{id: "E01-foo/T002-task", deps: []string{"E01-foo/T003-task"}},
		{id: "E01-foo/T003-task", deps: []string{"E01-foo/T001-task"}},
	})
	problems := CheckDependencies(root, "")
	var cycleMsg string
	for _, p := range problems {
		if strings.Contains(p.Message, "cycle") {
			cycleMsg = p.Message
			break
		}
	}
	if cycleMsg == "" {
		t.Fatal("CheckDependencies() = no cycle problem, want one")
	}
	// The cycle path should contain T001, T002, T003 in the correct order
	if !strings.Contains(cycleMsg, "T001") || !strings.Contains(cycleMsg, "T002") || !strings.Contains(cycleMsg, "T003") {
		t.Fatalf("CheckDependencies() cycle path = %q, should contain all three tasks", cycleMsg)
	}
	// Each arrow should separate consecutive nodes in the cycle
	if !strings.Contains(cycleMsg, "T001-task") || !strings.Contains(cycleMsg, "T002-task") || !strings.Contains(cycleMsg, "T003-task") {
		t.Fatalf("CheckDependencies() cycle path = %q, should reference task files", cycleMsg)
	}
}

func TestCheckDependencies_SelfReference(t *testing.T) {
	root := t.TempDir()
	setupMinimalProject(t, root, "v1", "E01-foo", []taskSpec{
		{id: "E01-foo/T001-task", deps: []string{"E01-foo/T001-task"}},
	})
	problems := CheckDependencies(root, "")
	found := false
	for _, p := range problems {
		if strings.Contains(p.Message, "cycle") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("CheckDependencies() = %v, want cycle problem (self-reference)", problems)
	}
}

func TestCheckDependencies_EpicFilter(t *testing.T) {
	root := t.TempDir()
	releasePath := filepath.Join(root, "releases", "v1")
	testutil.MkdirAll(t, filepath.Join(releasePath, "epics", "E01-foo", "tasks"))
	testutil.MkdirAll(t, filepath.Join(releasePath, "epics", "E02-bar", "tasks"))
	testutil.WriteReleasePRD(t, releasePath)
	testutil.WriteFile(t, filepath.Join(releasePath, "epics", "E01-foo", "E01-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E01: Foo\n")
	testutil.WriteFile(t, filepath.Join(releasePath, "epics", "E02-bar", "E02-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E02: Bar\n")

	// E02 has a missing dep — should be invisible with filter
	taskE2 := `---\nid: E02-bar/T001-task\nstatus: planned\nobjective: \"B\"\ndepends_on: [\"E02-bar/T999-nonexistent\"]\n---\n\n# T001\n\n## Acceptance Criteria\n\n- it works\n`
	testutil.WriteFile(t, filepath.Join(releasePath, "epics", "E01-foo", "tasks", "T001-task.md"), "---\nid: E01-foo/T001-task\nstatus: planned\nobjective: \"A\"\ndepends_on: []\n---\n\n# T001\n\n## Acceptance Criteria\n\n- it works\n")
	testutil.WriteFile(t, filepath.Join(releasePath, "epics", "E02-bar", "tasks", "T001-task.md"), strings.ReplaceAll(taskE2, "\\n", "\n"))

	problems := CheckDependencies(root, "E01-foo")
	for _, p := range problems {
		if strings.Contains(p.Message, "E02-bar") {
			t.Fatalf("CheckDependencies() with epicFilter=E01-foo should skip E02, got: %v", p)
		}
	}
}

// --- CheckAuditState ---

func TestCheckAuditState_NoAuditFiles(t *testing.T) {
	root := t.TempDir()
	testutil.MkdirAll(t, filepath.Join(root, "releases", "v1", "epics", "E01-foo"))
	testutil.WriteFile(t, filepath.Join(root, "router.md"), routerContent("task-building", "v1", "E01-foo"))
	problems := CheckAuditState(root)
	if len(problems) > 0 {
		t.Fatalf("CheckAuditState() = %v, want no problems", problems)
	}
}

func TestCheckAuditState_MatchesRouter(t *testing.T) {
	root := t.TempDir()
	epicPath := filepath.Join(root, "releases", "v1", "epics", "E01-foo")
	testutil.MkdirAll(t, epicPath)
	testutil.WriteFile(t, filepath.Join(epicPath, "E01-Audit.md"), "---\ntype: audit-findings\n---\n\n# Audit\n")
	testutil.WriteFile(t, filepath.Join(root, "router.md"), routerContent("audit-pending", "v1", "E01-foo"))
	problems := CheckAuditState(root)
	if len(problems) > 0 {
		t.Fatalf("CheckAuditState() = %v, want no problems when router matches", problems)
	}
}

func TestCheckAuditState_ProposalWithoutPending(t *testing.T) {
	root := t.TempDir()
	epicPath := filepath.Join(root, "releases", "v1", "epics", "E01-foo")
	testutil.MkdirAll(t, epicPath)
	testutil.WriteFile(t, filepath.Join(epicPath, "E01-Audit.md"), "---\ntype: audit-findings\n---\n\n# Audit\n")
	testutil.WriteFile(t, filepath.Join(root, "router.md"), routerContent("task-building", "v1", "E01-foo"))
	problems := CheckAuditState(root)
	if len(problems) != 1 {
		t.Fatalf("CheckAuditState() = %v, want 1 problem (audit file without audit-pending)", problems)
	}
	if !strings.Contains(problems[0].Message, "audit proposal exists") {
		t.Fatalf("CheckAuditState() = %v, want 'audit proposal exists' message", problems)
	}
}

func TestCheckAuditState_DifferentEpicInRouter(t *testing.T) {
	root := t.TempDir()
	epic1Path := filepath.Join(root, "releases", "v1", "epics", "E01-foo")
	epic2Path := filepath.Join(root, "releases", "v1", "epics", "E02-bar")
	testutil.MkdirAll(t, epic1Path)
	testutil.MkdirAll(t, epic2Path)
	testutil.WriteFile(t, filepath.Join(epic1Path, "E01-Audit.md"), "---\ntype: audit-findings\n---\n\n# Audit\n")
	testutil.WriteFile(t, filepath.Join(root, "router.md"), routerContent("audit-pending", "v1", "E02-bar"))
	problems := CheckAuditState(root)
	if len(problems) != 1 {
		t.Fatalf("CheckAuditState() = %v, want 1 problem (E01 audit but E02 in router)", problems)
	}
	if !strings.Contains(problems[0].Message, "E01") {
		t.Fatalf("CheckAuditState() = %v, want problem mentioning E01", problems)
	}
}

func TestCheckAuditState_MultipleStale(t *testing.T) {
	root := t.TempDir()
	epic1Path := filepath.Join(root, "releases", "v1", "epics", "E01-foo")
	epic2Path := filepath.Join(root, "releases", "v1", "epics", "E02-bar")
	testutil.MkdirAll(t, epic1Path)
	testutil.MkdirAll(t, epic2Path)
	testutil.WriteFile(t, filepath.Join(epic1Path, "E01-Audit.md"), "---\ntype: audit-findings\n---\n\n# Audit\n")
	testutil.WriteFile(t, filepath.Join(epic2Path, "E02-Audit.md"), "---\ntype: audit-findings\n---\n\n# Audit\n")
	testutil.WriteFile(t, filepath.Join(root, "router.md"), routerContent("task-building", "v1", "E03-baz"))
	problems := CheckAuditState(root)
	if len(problems) != 2 {
		t.Fatalf("CheckAuditState() = %v, want 2 problems (both audit files stale)", problems)
	}
}

// --- CheckOrphans ---

func TestCheckOrphans_NoOrphans(t *testing.T) {
	root := t.TempDir()
	setupMinimalProject(t, root, "v1", "E01-foo", []taskSpec{
		{id: "E01-foo/T001-task", deps: []string{}},
		{id: "E01-foo/T002-task", deps: []string{}},
	})
	problems := CheckOrphans(root)
	if len(problems) > 0 {
		t.Fatalf("CheckOrphans() = %v, want no problems", problems)
	}
}

func TestCheckOrphans_TaskRefersNonexistentEpic(t *testing.T) {
	root := t.TempDir()
	setupMinimalProject(t, root, "v1", "E01-foo", []taskSpec{
		{id: "E99-ghost/T001-task", deps: []string{}},
	})
	problems := CheckOrphans(root)
	found := false
	for _, p := range problems {
		if strings.Contains(p.Message, "orphaned") && strings.Contains(p.Message, "E99-ghost") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("CheckOrphans() = %v, want orphaned task problem for E99-ghost", problems)
	}
}

func TestCheckOrphans_CrossReleaseEpicRef(t *testing.T) {
	root := t.TempDir()
	releasePath := filepath.Join(root, "releases", "v1")
	epicPath := filepath.Join(releasePath, "epics", "E01-foo")
	tasksPath := filepath.Join(epicPath, "tasks")
	testutil.MkdirAll(t, tasksPath)
	testutil.WriteReleasePRD(t, releasePath)
	testutil.WriteFile(t, filepath.Join(epicPath, "E01-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E01: Foo\n")
	testutil.WriteFile(t, filepath.Join(tasksPath, "T001-task.md"), "---\nid: E02-bar/T001-task\nstatus: planned\nobjective: \"Task\"\ndepends_on: []\n---\n\n# T001\n\n## Acceptance Criteria\n\n- it works\n")
	// E02-bar does not exist in any release
	problems := CheckOrphans(root)
	found := false
	for _, p := range problems {
		if strings.Contains(p.Message, "orphaned") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("CheckOrphans() = %v, want orphaned task problem for E02-bar", problems)
	}
}

func TestCheckOrphans_ValidCrossReleaseRef(t *testing.T) {
	root := t.TempDir()
	releasePath := filepath.Join(root, "releases", "v1")
	epic1Path := filepath.Join(releasePath, "epics", "E01-foo")
	epic2Path := filepath.Join(releasePath, "epics", "E02-bar")
	testutil.MkdirAll(t, filepath.Join(epic1Path, "tasks"))
	testutil.MkdirAll(t, filepath.Join(epic2Path, "tasks"))
	testutil.WriteReleasePRD(t, releasePath)
	testutil.WriteFile(t, filepath.Join(epic1Path, "E01-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E01: Foo\n")
	testutil.WriteFile(t, filepath.Join(epic2Path, "E02-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E02: Bar\n")
	testutil.WriteFile(t, filepath.Join(epic1Path, "tasks", "T001-task.md"), "---\nid: E02-bar/T001-task\nstatus: planned\nobjective: \"Task\"\ndepends_on: []\n---\n\n# T001\n\n## Acceptance Criteria\n\n- it works\n")
	// E02-bar exists
	problems := CheckOrphans(root)
	for _, p := range problems {
		if strings.Contains(p.Message, "orphaned") {
			t.Fatalf("CheckOrphans() = %v, want no orphan problems for cross-epic ref that exists", problems)
		}
	}
}

func TestCheckOrphans_EmptyID(t *testing.T) {
	root := t.TempDir()
	setupMinimalProject(t, root, "v1", "E01-foo", []taskSpec{
		{id: "E01-foo/T001-task", deps: []string{}},
	})
	// Write a task with empty ID
	tasksPath := filepath.Join(root, "releases", "v1", "epics", "E01-foo", "tasks")
	testutil.WriteFile(t, filepath.Join(tasksPath, "T002-bad.md"), "---\nstatus: planned\nobjective: \"No ID\"\ndepends_on: []\n---\n\n# T002\n\n## Acceptance Criteria\n\n- it works\n")
	problems := CheckOrphans(root)
	// Should not crash, should handle missing ID gracefully
	if len(problems) > 0 {
		// Only allow non-orphan problems (e.g. missing ID)
		for _, p := range problems {
			if strings.Contains(p.Message, "orphaned") {
				t.Fatalf("CheckOrphans() = %v, want no orphan problems for task with missing ID", problems)
			}
		}
	}
}

func TestCheckOrphans_NoReleasesDir(t *testing.T) {
	root := t.TempDir()
	problems := CheckOrphans(root)
	// Should report releases dir problem, not crash
	if len(problems) == 0 {
		t.Fatal("CheckOrphans() = no problems, want error about missing releases")
	}
}

// helpers

type taskSpec struct {
	id   string
	deps []string
}

func setupMinimalProject(t *testing.T, root, releaseID, epicID string, tasks []taskSpec) {
	t.Helper()
	releasePath := filepath.Join(root, "releases", releaseID)
	epicPath := filepath.Join(releasePath, "epics", epicID)
	tasksPath := filepath.Join(epicPath, "tasks")
	testutil.MkdirAll(t, tasksPath)

	prefix := epicID
	if idx := strings.IndexByte(epicID, '-'); idx != -1 {
		prefix = epicID[:idx]
	}

	testutil.WriteReleasePRD(t, releasePath)
	testutil.WriteFile(t, filepath.Join(epicPath, prefix+"-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# Epic\n")

	for i, ts := range tasks {
		depsYAML := "[]"
		if len(ts.deps) > 0 {
			quoted := make([]string, len(ts.deps))
			for j, d := range ts.deps {
				quoted[j] = fmt.Sprintf("%q", d)
			}
			depsYAML = "[" + strings.Join(quoted, ", ") + "]"
		}
		content := fmt.Sprintf("---\nid: %s\nstatus: planned\nobjective: \"Task %d\"\ndepends_on: %s\n---\n\n# T%03d\n\n## Acceptance Criteria\n\n- it works\n", ts.id, i, depsYAML, i+1)
		testutil.WriteFile(t, filepath.Join(tasksPath, fmt.Sprintf("T%03d-task.md", i+1)), content)
	}
}

func routerContent(state, release, epic string) string {
	return "## Current state\n\n```yaml\nstate: " + state + "\nrelease: " + release + "\nepic: " + epic + "\ntask: none\nnext_action: \"\"\n```\n"
}

// --- CheckStructure complexity ---

func TestCheckStructure_TaskInvalidComplexityTier(t *testing.T) {
	root := t.TempDir()
	releasePath := filepath.Join(root, "releases", "v1")
	epicPath := filepath.Join(releasePath, "epics", "E01-foo")
	tasksPath := filepath.Join(epicPath, "tasks")
	testutil.MkdirAll(t, tasksPath)
	testutil.WriteReleasePRD(t, releasePath)
	testutil.WriteFile(t, filepath.Join(epicPath, "E01-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E01: Foo\n")
	testutil.WriteFile(t, filepath.Join(tasksPath, "T001-task.md"),
		"---\nid: E01-foo/T001-task\nstatus: planned\nobjective: \"Do the thing\"\ncomplexity_tier: extreme\ncomplexity_reason: \"Some reason.\"\ndepends_on: []\n---\n\n# T001\n\n## Acceptance Criteria\n\n- It works\n")
	problems := CheckStructure(root, "")
	found := false
	for _, p := range problems {
		if strings.Contains(p.Message, "complexity") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("CheckStructure() = %v, want complexity invalid problem", problems)
	}
}

func TestCheckStructure_TaskValidComplexity(t *testing.T) {
	root := t.TempDir()
	releasePath := filepath.Join(root, "releases", "v1")
	epicPath := filepath.Join(releasePath, "epics", "E01-foo")
	tasksPath := filepath.Join(epicPath, "tasks")
	testutil.MkdirAll(t, tasksPath)
	testutil.WriteReleasePRD(t, releasePath)
	testutil.WriteFile(t, filepath.Join(epicPath, "E01-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E01: Foo\n")
	testutil.WriteFile(t, filepath.Join(tasksPath, "T001-task.md"),
		"---\nid: E01-foo/T001-task\nstatus: planned\nobjective: \"Do the thing\"\ncomplexity_tier: medium\ncomplexity_reason: \"Touches parser and write paths.\"\ndepends_on: []\n---\n\n# T001\n\n## Acceptance Criteria\n\n- It works\n")
	problems := CheckStructure(root, "")
	for _, p := range problems {
		if strings.Contains(p.Message, "complexity") {
			t.Fatalf("CheckStructure() unexpected complexity problem: %v", p)
		}
	}
}

func TestCheckStructure_TaskComplexityAbsentNoProblems(t *testing.T) {
	root := t.TempDir()
	releasePath := filepath.Join(root, "releases", "v1")
	epicPath := filepath.Join(releasePath, "epics", "E01-foo")
	tasksPath := filepath.Join(epicPath, "tasks")
	testutil.MkdirAll(t, tasksPath)
	testutil.WriteReleasePRD(t, releasePath)
	testutil.WriteFile(t, filepath.Join(epicPath, "E01-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E01: Foo\n")
	testutil.WriteFile(t, filepath.Join(tasksPath, "T001-task.md"),
		"---\nid: E01-foo/T001-task\nstatus: planned\nobjective: \"Do the thing\"\ndepends_on: []\n---\n\n# T001\n\n## Acceptance Criteria\n\n- It works\n")
	problems := CheckStructure(root, "")
	for _, p := range problems {
		if strings.Contains(p.Message, "complexity") {
			t.Fatalf("CheckStructure() unexpected complexity problem for absent fields: %v", p)
		}
	}
}

// --- CheckDefects ---

func writeDefect(t *testing.T, defectsDir, filename, content string) {
	t.Helper()
	testutil.MkdirAll(t, defectsDir)
	testutil.WriteFile(t, filepath.Join(defectsDir, filename), content)
}

func TestCheckDefects_NoDefectsDir(t *testing.T) {
	root := t.TempDir()
	testutil.MkdirAll(t, filepath.Join(root, "releases", "v1"))
	testutil.WriteReleasePRD(t, filepath.Join(root, "releases", "v1"))
	problems := CheckDefects(root)
	if len(problems) > 0 {
		t.Fatalf("CheckDefects() = %v, want no problems when defects dir absent", problems)
	}
}

func TestCheckDefects_ValidDefect(t *testing.T) {
	root := t.TempDir()
	defectsDir := filepath.Join(root, "releases", "v1", "defects")
	writeDefect(t, defectsDir, "D001-crash.md",
		"---\nid: v1/D001-crash\nstatus: open\nseverity: high\ntitle: Auth crash\n---\n\n## Problem\n\nIt crashes.\n")
	problems := CheckDefects(root)
	if len(problems) > 0 {
		t.Fatalf("CheckDefects() = %v, want no problems for valid defect", problems)
	}
}

func TestCheckDefects_MalformedYAML(t *testing.T) {
	root := t.TempDir()
	defectsDir := filepath.Join(root, "releases", "v1", "defects")
	writeDefect(t, defectsDir, "D001-bad.md", "---\n: bad: [yaml\n---\n")
	problems := CheckDefects(root)
	found := false
	for _, p := range problems {
		if strings.Contains(p.Message, "parse error") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("CheckDefects() = %v, want parse error problem", problems)
	}
}

func TestCheckDefects_InvalidStatus(t *testing.T) {
	root := t.TempDir()
	defectsDir := filepath.Join(root, "releases", "v1", "defects")
	writeDefect(t, defectsDir, "D001-bad-status.md",
		"---\nid: v1/D001-bad-status\nstatus: blocked\nseverity: low\ntitle: Bad\n---\n")
	problems := CheckDefects(root)
	found := false
	for _, p := range problems {
		if strings.Contains(p.Message, "defect status invalid") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("CheckDefects() = %v, want invalid status problem", problems)
	}
}

func TestCheckDefects_InProgressMissingStage(t *testing.T) {
	root := t.TempDir()
	defectsDir := filepath.Join(root, "releases", "v1", "defects")
	writeDefect(t, defectsDir, "D001-in-progress.md",
		"---\nid: v1/D001-in-progress\nstatus: in_progress\nseverity: medium\ntitle: In progress\n---\n")
	problems := CheckDefects(root)
	found := false
	for _, p := range problems {
		if strings.Contains(p.Message, "stage is required") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("CheckDefects() = %v, want missing stage problem", problems)
	}
}

func TestCheckDefects_InvalidStageWhenInProgress(t *testing.T) {
	root := t.TempDir()
	defectsDir := filepath.Join(root, "releases", "v1", "defects")
	writeDefect(t, defectsDir, "D001-bad-stage.md",
		"---\nid: v1/D001-bad-stage\nstatus: in_progress\nstage: review\nseverity: high\ntitle: Bad stage\n---\n")
	problems := CheckDefects(root)
	found := false
	for _, p := range problems {
		if strings.Contains(p.Message, "defect stage invalid") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("CheckDefects() = %v, want invalid stage problem", problems)
	}
}

func TestCheckDefects_TaskStyleStatusAlias(t *testing.T) {
	root := t.TempDir()
	defectsDir := filepath.Join(root, "releases", "v1", "defects")
	writeDefect(t, defectsDir, "D001-alias.md",
		"---\nid: v1/D001-alias\nstatus: done\nseverity: low\ntitle: Alias\n---\n")
	problems := CheckDefects(root)
	found := false
	for _, p := range problems {
		if strings.Contains(p.Message, "non-canonical status") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("CheckDefects() = %v, want non-canonical status problem", problems)
	}
}

func TestCheckDefects_StaleStageOutsideInProgress(t *testing.T) {
	root := t.TempDir()
	defectsDir := filepath.Join(root, "releases", "v1", "defects")
	writeDefect(t, defectsDir, "D001-stale.md",
		"---\nid: v1/D001-stale\nstatus: resolved\nstage: audit\nseverity: low\ntitle: Stale\n---\n")
	problems := CheckDefects(root)
	found := false
	for _, p := range problems {
		if strings.Contains(p.Message, "only valid when status is in_progress") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("CheckDefects() = %v, want stale stage problem", problems)
	}
}

func TestCheckDefects_MissingID(t *testing.T) {
	root := t.TempDir()
	defectsDir := filepath.Join(root, "releases", "v1", "defects")
	writeDefect(t, defectsDir, "D001-no-id.md",
		"---\nstatus: open\nseverity: low\ntitle: No ID\n---\n")
	problems := CheckDefects(root)
	found := false
	for _, p := range problems {
		if strings.Contains(p.Message, "missing required frontmatter field: id") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("CheckDefects() = %v, want missing id problem", problems)
	}
}

func TestCheckDefects_MissingSeverity(t *testing.T) {
	root := t.TempDir()
	defectsDir := filepath.Join(root, "releases", "v1", "defects")
	writeDefect(t, defectsDir, "D001-no-severity.md",
		"---\nid: v1/D001-no-severity\nstatus: open\ntitle: No Severity\n---\n")
	problems := CheckDefects(root)
	found := false
	for _, p := range problems {
		if strings.Contains(p.Message, "missing required frontmatter field: severity") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("CheckDefects() = %v, want missing severity problem", problems)
	}
}

func TestCheckDefects_BrokenReferenceEmptyComponent(t *testing.T) {
	root := t.TempDir()
	defectsDir := filepath.Join(root, "releases", "v1", "defects")
	writeDefect(t, defectsDir, "D001-bad-ref.md",
		"---\nid: v1/D001-bad-ref\nstatus: open\nseverity: low\ntitle: Bad ref\nreference: \"/T003\"\n---\n")
	problems := CheckDefects(root)
	found := false
	for _, p := range problems {
		if strings.Contains(p.Message, "empty epic or task component") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("CheckDefects() = %v, want broken reference problem", problems)
	}
}

func TestCheckDefects_ReferenceNoSlashIsAccepted(t *testing.T) {
	root := t.TempDir()
	defectsDir := filepath.Join(root, "releases", "v1", "defects")
	writeDefect(t, defectsDir, "D001-plain-ref.md",
		"---\nid: v1/D001-plain-ref\nstatus: open\nseverity: low\ntitle: Plain ref\nreference: \"v1.0.5\"\n---\n")
	problems := CheckDefects(root)
	if len(problems) > 0 {
		t.Fatalf("CheckDefects() = %v, want no problems for non-slash reference", problems)
	}
}

func TestCheckDefects_ReferenceMatchesExistingTask(t *testing.T) {
	root := t.TempDir()
	// Set up a real task
	setupMinimalProject(t, root, "v1", "E01-foo", []taskSpec{
		{id: "E01-foo/T001-task", deps: []string{}},
	})
	defectsDir := filepath.Join(root, "releases", "v1", "defects")
	writeDefect(t, defectsDir, "D001-valid-ref.md",
		"---\nid: v1/D001-valid-ref\nstatus: open\nseverity: low\ntitle: Valid ref\nreference: \"E01-foo/T001-task\"\n---\n")
	problems := CheckDefects(root)
	if len(problems) > 0 {
		t.Fatalf("CheckDefects() = %v, want no problems when reference matches task", problems)
	}
}

func TestCheckDefects_ReferenceDoesNotMatchTask(t *testing.T) {
	root := t.TempDir()
	setupMinimalProject(t, root, "v1", "E01-foo", []taskSpec{
		{id: "E01-foo/T001-task", deps: []string{}},
	})
	defectsDir := filepath.Join(root, "releases", "v1", "defects")
	writeDefect(t, defectsDir, "D001-bad-ref.md",
		"---\nid: v1/D001-bad-ref\nstatus: open\nseverity: low\ntitle: Bad ref\nreference: \"E01-foo/T999-ghost\"\n---\n")
	problems := CheckDefects(root)
	found := false
	for _, p := range problems {
		if strings.Contains(p.Message, "does not match any known task ID") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("CheckDefects() = %v, want missing task reference problem", problems)
	}
}

// --- CheckAuditRegister ---

// writeAuditFinding writes a finding record under root/audit/findings with the
// given frontmatter block (without --- delimiters).
func writeAuditFinding(t *testing.T, root, filename, frontmatter string) string {
	t.Helper()
	path := filepath.Join(root, "audit", "findings", filename)
	testutil.WriteFile(t, path, "---\n"+frontmatter+"---\n\n# Finding\n\n## Proof\n\nPending.\n")
	return path
}

// cleanFindingFrontmatter is a fully valid F001 finding referencing the v1
// release, the E01-foo epic, and the E01-foo/T001-task task.
const cleanFindingFrontmatter = `id: F001
title: "Sample finding"
status: open
severity: medium
confidence: high
proof_needed: "regression test in checks_test.go"
first_seen: "2026-07-01"
last_seen: "2026-07-01"
releases: ["v1"]
epics: ["E01-foo"]
tasks: ["E01-foo/T001-task"]
`

func TestCheckAuditRegister_AbsentAuditDir(t *testing.T) {
	root := t.TempDir()
	setupMinimalProject(t, root, "v1", "E01-foo", []taskSpec{
		{id: "E01-foo/T001-task", deps: []string{}},
	})

	problems := CheckAuditRegister(root)
	if len(problems) > 0 {
		t.Fatalf("CheckAuditRegister() = %v, want no problems when audit/ is absent", problems)
	}
}

func TestCheckAuditRegister_PartialAdoptionWithoutFindings(t *testing.T) {
	root := t.TempDir()
	setupMinimalProject(t, root, "v1", "E01-foo", []taskSpec{
		{id: "E01-foo/T001-task", deps: []string{}},
	})
	testutil.WriteFile(t, filepath.Join(root, "audit", "register.md"), "# Audit Register\n")
	testutil.WriteFile(t, filepath.Join(root, "audit", "findings", "README.md"), "# Findings guidance\n")
	testutil.WriteFile(t, filepath.Join(root, "audit", "runs", "README.md"), "# Runs guidance\n")

	problems := CheckAuditRegister(root)
	if len(problems) > 0 {
		t.Fatalf("CheckAuditRegister() = %v, want no problems for partially adopted register", problems)
	}
}

func TestCheckAuditRegister_CleanFinding(t *testing.T) {
	root := t.TempDir()
	setupMinimalProject(t, root, "v1", "E01-foo", []taskSpec{
		{id: "E01-foo/T001-task", deps: []string{}},
	})
	writeAuditFinding(t, root, "F001-sample.md", cleanFindingFrontmatter)

	problems := CheckAuditRegister(root)
	if len(problems) > 0 {
		t.Fatalf("CheckAuditRegister() = %v, want no problems for clean finding", problems)
	}
}

func TestCheckAuditRegister_ReportsHealedFields(t *testing.T) {
	root := t.TempDir()
	setupMinimalProject(t, root, "v1", "E01-foo", []taskSpec{
		{id: "E01-foo/T001-task", deps: []string{}},
	})
	writeAuditFinding(t, root, "F002-bad-fields.md", `id: F001
title: "Bad fields"
status: reviewing
severity: blocker
confidence: certain
proof_needed: "test"
first_seen: "2026-07-01"
last_seen: "2026-07-01"
`)

	problems := CheckAuditRegister(root)
	wants := []string{
		`finding id "F001" does not match filename id "F002"`,
		`finding status invalid "reviewing"`,
		`finding severity invalid "blocker"`,
		`finding confidence invalid "certain"`,
	}
	for _, want := range wants {
		if !problemsContain(problems, want) {
			t.Errorf("CheckAuditRegister() missing problem %q in %v", want, problems)
		}
	}
}

func TestCheckAuditRegister_ReportsMissingFields(t *testing.T) {
	root := t.TempDir()
	setupMinimalProject(t, root, "v1", "E01-foo", []taskSpec{
		{id: "E01-foo/T001-task", deps: []string{}},
	})
	writeAuditFinding(t, root, "F003-missing-fields.md", "id: F003\nstatus: open\n")

	problems := CheckAuditRegister(root)
	for _, field := range []string{"title", "severity", "confidence", "proof_needed", "first_seen", "last_seen"} {
		want := "finding missing required field: " + field
		if !problemsContain(problems, want) {
			t.Errorf("CheckAuditRegister() missing problem %q in %v", want, problems)
		}
	}
}

func TestCheckAuditRegister_ReportsVerifiedMissingProof(t *testing.T) {
	root := t.TempDir()
	setupMinimalProject(t, root, "v1", "E01-foo", []taskSpec{
		{id: "E01-foo/T001-task", deps: []string{}},
	})
	writeAuditFinding(t, root, "F001-verified.md", `id: F001
title: "Verified without proof"
status: verified
severity: high
confidence: high
proof_needed: "regression test"
first_seen: "2026-07-01"
last_seen: "2026-07-01"
`)

	problems := CheckAuditRegister(root)
	if !problemsContain(problems, "verified finding has no named proof") {
		t.Fatalf("CheckAuditRegister() = %v, want verified-missing-proof problem", problems)
	}
}

func TestCheckAuditRegister_ReportsDuplicateProblems(t *testing.T) {
	root := t.TempDir()
	setupMinimalProject(t, root, "v1", "E01-foo", []taskSpec{
		{id: "E01-foo/T001-task", deps: []string{}},
	})
	writeAuditFinding(t, root, "F001-dup-no-target.md", `id: F001
title: "Duplicate without target"
status: duplicate
severity: low
confidence: medium
proof_needed: "n/a"
first_seen: "2026-07-01"
last_seen: "2026-07-01"
`)
	writeAuditFinding(t, root, "F002-dup-bad-target.md", `id: F002
title: "Duplicate of missing finding"
status: duplicate
severity: low
confidence: medium
proof_needed: "n/a"
first_seen: "2026-07-01"
last_seen: "2026-07-01"
duplicate_of: F999
`)

	problems := CheckAuditRegister(root)
	if !problemsContain(problems, "duplicate finding has no duplicate_of") {
		t.Errorf("CheckAuditRegister() = %v, want duplicate-missing-target problem", problems)
	}
	if !problemsContain(problems, `duplicate_of references unknown finding "F999"`) {
		t.Errorf("CheckAuditRegister() = %v, want unknown duplicate_of problem", problems)
	}
}

func TestCheckAuditRegister_ReportsBrokenWorkItemLinks(t *testing.T) {
	root := t.TempDir()
	setupMinimalProject(t, root, "v1", "E01-foo", []taskSpec{
		{id: "E01-foo/T001-task", deps: []string{}},
	})
	writeDefect(t, filepath.Join(root, "releases", "v1", "defects"), "D001-real.md",
		"---\nid: v1/D001-real\nstatus: open\nseverity: low\ntitle: Real\n---\n")
	writeAuditFinding(t, root, "F001-broken-links.md", `id: F001
title: "Broken links"
status: open
severity: medium
confidence: medium
proof_needed: "test"
first_seen: "2026-07-01"
last_seen: "2026-07-01"
releases: ["v9"]
epics: ["E99-ghost"]
tasks: ["E01-foo/T999-ghost"]
defects: ["D001-real", "D404-nope"]
`)

	problems := CheckAuditRegister(root)
	wants := []string{
		`references unknown release "v9"`,
		`references unknown epic "E99-ghost"`,
		`references unknown task "E01-foo/T999-ghost"`,
		`references unknown defect "D404-nope"`,
	}
	for _, want := range wants {
		if !problemsContain(problems, want) {
			t.Errorf("CheckAuditRegister() missing problem %q in %v", want, problems)
		}
	}
	if problemsContain(problems, `unknown defect "D001-real"`) {
		t.Errorf("CheckAuditRegister() = %v, flagged existing defect D001-real", problems)
	}
}

func TestCheckAuditRegister_MalformedFindingFile(t *testing.T) {
	root := t.TempDir()
	setupMinimalProject(t, root, "v1", "E01-foo", []taskSpec{
		{id: "E01-foo/T001-task", deps: []string{}},
	})
	testutil.WriteFile(t, filepath.Join(root, "audit", "findings", "F001-broken.md"),
		"---\nid: [broken\n---\n\n# Finding\n")

	problems := CheckAuditRegister(root)
	if !problemsContain(problems, "audit finding unloadable") {
		t.Fatalf("CheckAuditRegister() = %v, want structural finding problem", problems)
	}
}

func TestCheckAuditRegister_MalformedRunFile(t *testing.T) {
	root := t.TempDir()
	setupMinimalProject(t, root, "v1", "E01-foo", []taskSpec{
		{id: "E01-foo/T001-task", deps: []string{}},
	})
	writeAuditFinding(t, root, "F001-sample.md", cleanFindingFrontmatter)
	testutil.WriteFile(t, filepath.Join(root, "audit", "runs", "2026-07-01-full.md"),
		"---\ndate: [broken\n---\n\n# Run\n")

	problems := CheckAuditRegister(root)
	if !problemsContain(problems, "audit run unloadable") {
		t.Fatalf("CheckAuditRegister() = %v, want structural run problem", problems)
	}
}

func problemsContain(problems []Problem, want string) bool {
	for _, p := range problems {
		if strings.Contains(p.Message, want) {
			return true
		}
	}
	return false
}

// --- CheckProject ---

func writeV2Objective(t *testing.T, root, dirName, id, title string) string {
	t.Helper()
	path := filepath.Join(root, "objectives", dirName, "Objective.md")
	testutil.WriteFile(t, path, "---\nid: "+id+"\ntitle: \""+title+"\"\nstatus: planned\n---\n\n# "+title+"\n")
	return path
}

func writeV2Task(t *testing.T, root, objDirName, fileName, id, title, objective string) string {
	t.Helper()
	path := filepath.Join(root, "objectives", objDirName, "tasks", fileName)
	testutil.WriteFile(t, path, "---\nid: "+id+"\ntitle: \""+title+"\"\nobjective: "+objective+"\nstatus: planned\n---\n\n# "+title+"\n")
	return path
}

func TestCheckProject_v1ProjectNoProblems(t *testing.T) {
	root := t.TempDir()
	testutil.SetupMinimalProject(t, root, "v1", "E01-foo")

	if problems := CheckProject(root); len(problems) != 0 {
		t.Fatalf("CheckProject() = %v, want no problems for a V1 project", problems)
	}
}

func TestCheckProject_v2ValidNoProblems(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Objective(t, root, "O001-ship", "O001", "Ship it")
	writeV2Task(t, root, "O001-ship", "T001-write.md", "T001", "Write it", "O001")

	if problems := CheckProject(root); len(problems) != 0 {
		t.Fatalf("CheckProject() = %v, want no problems for a valid V2 project", problems)
	}
}

func TestCheckProject_SchemaVersionMalformed(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: nope\n")

	problems := CheckProject(root)
	if len(problems) != 1 {
		t.Fatalf("CheckProject() = %v, want exactly 1 problem", problems)
	}
	if !strings.Contains(problems[0].Message, "[schema-version-malformed]") {
		t.Errorf("CheckProject()[0].Message = %q, want the schema-version-malformed diagnostic name", problems[0].Message)
	}
	if problems[0].File != root {
		t.Errorf("CheckProject()[0].File = %q, want project root %q", problems[0].File, root)
	}
	if problems[0].Repair == "" {
		t.Error("CheckProject()[0].Repair is empty, want a typed repair suggestion")
	}
}

func TestCheckProject_SchemaVersionUnsupported(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 99\n")

	problems := CheckProject(root)
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[schema-version-unsupported]") {
		t.Fatalf("CheckProject() = %v, want 1 problem naming schema-version-unsupported", problems)
	}
}

// TestCheckProject_MissingTaskTitleNotBackfilledFromObjective proves doctor
// reports a missing V2 Task title as an actionable schema error rather than
// silently accepting the Task's objective owner reference as a display
// title substitute.
func TestCheckProject_MissingTaskTitleNotBackfilledFromObjective(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Objective(t, root, "O001-ship", "O001", "Ship it")
	testutil.WriteFile(t, filepath.Join(root, "objectives", "O001-ship", "tasks", "T001-write.md"),
		"---\nid: T001\nobjective: O001\nstatus: planned\n---\n\n# Write it\n")

	problems := CheckProject(root)
	if len(problems) != 1 {
		t.Fatalf("CheckProject() = %v, want exactly 1 problem", problems)
	}
	if !strings.Contains(problems[0].Message, "[v2-missing-field]") {
		t.Errorf("CheckProject()[0].Message = %q, want the v2-missing-field diagnostic name", problems[0].Message)
	}
	if !strings.Contains(problems[0].Message, "missing required field title") {
		t.Errorf("CheckProject()[0].Message = %q, want it to name the missing title field, not accept objective O001 as a substitute", problems[0].Message)
	}
}

func TestCheckProject_MissingOwner(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Objective(t, root, "O001-ship", "O001", "Ship it")
	writeV2Task(t, root, "O001-ship", "T001-write.md", "T001", "Write it", "O999")

	problems := CheckProject(root)
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-missing-owner]") {
		t.Fatalf("CheckProject() = %v, want 1 problem naming v2-missing-owner", problems)
	}
}

func TestCheckProject_DependencyCycle(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Objective(t, root, "O001-ship", "O001", "Ship it")
	testutil.WriteFile(t, filepath.Join(root, "objectives", "O001-ship", "tasks", "T001-alpha.md"),
		"---\nid: T001\ntitle: \"Alpha\"\nobjective: O001\nstatus: planned\ndepends_on: [{task: T002}]\n---\n\n# Alpha\n")
	testutil.WriteFile(t, filepath.Join(root, "objectives", "O001-ship", "tasks", "T002-beta.md"),
		"---\nid: T002\ntitle: \"Beta\"\nobjective: O001\nstatus: planned\ndepends_on: [{task: T001}]\n---\n\n# Beta\n")

	problems := CheckProject(root)
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-dependency-cycle]") {
		t.Fatalf("CheckProject() = %v, want 1 problem naming v2-dependency-cycle", problems)
	}
}

// TestCheckProject_ReadOnly proves CheckProject never writes to the project
// it inspects, for a valid V1 project with quality gates configured, a valid
// V2 project, and an invalid V2 project.
func TestCheckProject_ReadOnly(t *testing.T) {
	t.Run("valid V1 project with quality gates configured", func(t *testing.T) {
		root := t.TempDir()
		testutil.SetupMinimalProject(t, root, "v1", "E01-foo")
		configPath := filepath.Join(root, "config.yml")
		before, err := os.ReadFile(configPath)
		if err != nil {
			t.Fatal(err)
		}

		RunAllChecks(root, "")

		after, err := os.ReadFile(configPath)
		if err != nil {
			t.Fatal(err)
		}
		if string(before) != string(after) {
			t.Fatalf("RunAllChecks() changed config.yml:\nbefore: %s\nafter: %s", before, after)
		}
	})

	t.Run("valid V2 project", func(t *testing.T) {
		root := t.TempDir()
		testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
		writeV2Objective(t, root, "O001-ship", "O001", "Ship it")
		taskPath := writeV2Task(t, root, "O001-ship", "T001-write.md", "T001", "Write it", "O001")
		before, err := os.ReadFile(taskPath)
		if err != nil {
			t.Fatal(err)
		}

		if problems := CheckProject(root); len(problems) != 0 {
			t.Fatalf("CheckProject() = %v, want no problems", problems)
		}

		after, err := os.ReadFile(taskPath)
		if err != nil {
			t.Fatal(err)
		}
		if string(before) != string(after) {
			t.Fatalf("CheckProject() changed %s:\nbefore: %s\nafter: %s", taskPath, before, after)
		}
	})

	t.Run("invalid V2 project", func(t *testing.T) {
		root := t.TempDir()
		testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
		writeV2Objective(t, root, "O001-ship", "O001", "Ship it")
		taskPath := writeV2Task(t, root, "O001-ship", "T001-write.md", "T001", "Write it", "O999")
		before, err := os.ReadFile(taskPath)
		if err != nil {
			t.Fatal(err)
		}

		if problems := CheckProject(root); len(problems) == 0 {
			t.Fatal("CheckProject() = no problems, want the missing-owner diagnostic")
		}

		after, err := os.ReadFile(taskPath)
		if err != nil {
			t.Fatal(err)
		}
		if string(before) != string(after) {
			t.Fatalf("CheckProject() changed %s:\nbefore: %s\nafter: %s", taskPath, before, after)
		}
	})
}

// --- CheckProject: Check and evidence diagnostics ---

func writeV2Check(t *testing.T, root, id, scope, result, supersedes string) string {
	t.Helper()
	path := filepath.Join(root, "checks", id+".md")
	body := "---\nid: " + id + "\nscope: " + scope + "\nresult: " + result + "\n" +
		"checked_by: {role: checker, session: sess-1}\nchecked_at: '2026-09-14T00:00:00Z'\n"
	if supersedes != "" {
		body += "supersedes: " + supersedes + "\n"
	}
	body += "---\n\n# Check\n"
	testutil.WriteFile(t, path, body)
	return path
}

func TestCheckProject_CheckMalformed(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Objective(t, root, "O001-ship", "O001", "Ship it")
	writeV2Task(t, root, "O001-ship", "T001-write.md", "T001", "Write it", "O001")
	writeV2Check(t, root, "C001", "{kind: task, id: T001}", "BOGUS", "")

	problems := CheckProject(root)
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-check-malformed]") {
		t.Fatalf("CheckProject() = %v, want 1 problem naming v2-check-malformed", problems)
	}
	if problems[0].Repair == "" {
		t.Error("CheckProject()[0].Repair is empty, want a typed repair suggestion")
	}
}

func TestCheckProject_CheckMissingScopeTarget(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Objective(t, root, "O001-ship", "O001", "Ship it")
	writeV2Task(t, root, "O001-ship", "T001-write.md", "T001", "Write it", "O001")
	writeV2Check(t, root, "C001", "{kind: task, id: T999}", "CLEAR", "")

	problems := CheckProject(root)
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-check-missing-scope-target]") {
		t.Fatalf("CheckProject() = %v, want 1 problem naming v2-check-missing-scope-target", problems)
	}
}

func TestCheckProject_CheckMissingReference(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Objective(t, root, "O001-ship", "O001", "Ship it")
	writeV2Task(t, root, "O001-ship", "T001-write.md", "T001", "Write it", "O001")
	writeV2Check(t, root, "C002", "{kind: task, id: T001}", "CLEAR", "C001")

	problems := CheckProject(root)
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-check-missing-reference]") {
		t.Fatalf("CheckProject() = %v, want 1 problem naming v2-check-missing-reference", problems)
	}
}

func TestCheckProject_CheckSupersedesConflict(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Objective(t, root, "O001-ship", "O001", "Ship it")
	writeV2Task(t, root, "O001-ship", "T001-write.md", "T001", "Write it", "O001")
	writeV2Check(t, root, "C001", "{kind: task, id: T001}", "CLEAR", "")
	writeV2Check(t, root, "C002", "{kind: task, id: T001}", "CLEAR", "C001")
	writeV2Check(t, root, "C003", "{kind: task, id: T001}", "CLEAR", "C001")

	problems := CheckProject(root)
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-check-supersedes-conflict]") {
		t.Fatalf("CheckProject() = %v, want 1 problem naming v2-check-supersedes-conflict", problems)
	}
}

func TestCheckProject_EvidenceMalformed(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Objective(t, root, "O001-ship", "O001", "Ship it")
	writeV2Check(t, root, "C001", "{kind: task, id: T001}", "CLEAR", "")
	testutil.WriteFile(t, filepath.Join(root, "objectives", "O001-ship", "tasks", "T001-write.md"),
		"---\nid: T001\ntitle: \"Write it\"\nobjective: O001\nstatus: planned\n"+
			"freshness: {state: bogus, check: C001, assessed_by: {role: checker, session: s}, assessed_at: '2026-09-14T00:00:00Z', basis: reviewed}\n"+
			"---\n\n# Write it\n")

	problems := CheckProject(root)
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-evidence-malformed]") {
		t.Fatalf("CheckProject() = %v, want 1 problem naming v2-evidence-malformed", problems)
	}
}

func TestCheckProject_EvidenceMissingReference(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Objective(t, root, "O001-ship", "O001", "Ship it")
	testutil.WriteFile(t, filepath.Join(root, "objectives", "O001-ship", "tasks", "T001-write.md"),
		"---\nid: T001\ntitle: \"Write it\"\nobjective: O001\nstatus: planned\nlast_check: C999\n---\n\n# Write it\n")

	problems := CheckProject(root)
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-evidence-missing-reference]") {
		t.Fatalf("CheckProject() = %v, want 1 problem naming v2-evidence-missing-reference", problems)
	}
}

func writeV2Issue(t *testing.T, root, fileName, id, status, issueType, extra string) string {
	t.Helper()
	path := filepath.Join(root, "issues", fileName)
	body := "---\nid: " + id + "\ntitle: \"Follow-up\"\ntype: " + issueType + "\nstatus: " + status +
		"\nsource: {kind: report, actor: {role: owner, session: owner-1}, at: '2026-09-15T00:00:00Z'}\n" + extra +
		"---\n\n# Issue\n"
	testutil.WriteFile(t, path, body)
	return path
}

func TestCheckProject_IssueMalformed(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Issue(t, root, "I001-bad.md", "I001", "open", "bogus", "")

	problems := CheckProject(root)
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-issue-malformed]") {
		t.Fatalf("CheckProject() = %v, want 1 problem naming v2-issue-malformed", problems)
	}
}

func TestCheckProject_IssueMissingDuplicateTarget(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Issue(t, root, "I001-x.md", "I001", "open", "defect", "duplicate_of: I999\n")

	problems := CheckProject(root)
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-issue-missing-duplicate-target]") {
		t.Fatalf("CheckProject() = %v, want 1 problem naming v2-issue-missing-duplicate-target", problems)
	}
}

func TestCheckProject_IssueSelfDuplicate(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Issue(t, root, "I001-x.md", "I001", "open", "defect", "duplicate_of: I001\n")

	problems := CheckProject(root)
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-issue-self-duplicate]") {
		t.Fatalf("CheckProject() = %v, want 1 problem naming v2-issue-self-duplicate", problems)
	}
}

func TestCheckProject_IssueDuplicateCycle(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Issue(t, root, "I001-x.md", "I001", "open", "defect", "duplicate_of: I002\n")
	writeV2Issue(t, root, "I002-y.md", "I002", "open", "defect", "duplicate_of: I001\n")

	problems := CheckProject(root)
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-issue-duplicate-cycle]") {
		t.Fatalf("CheckProject() = %v, want 1 problem naming v2-issue-duplicate-cycle", problems)
	}
}

func TestCheckProject_IssueMissingLinkTarget(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Issue(t, root, "I001-x.md", "I001", "open", "defect", "tasks: [T999]\n")

	problems := CheckProject(root)
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-issue-missing-link-target]") {
		t.Fatalf("CheckProject() = %v, want 1 problem naming v2-issue-missing-link-target", problems)
	}
}

func TestCheckProject_IssueUnpairedCheckLink(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Objective(t, root, "O001-ship", "O001", "Ship it")
	writeV2Task(t, root, "O001-ship", "T001-write.md", "T001", "Write it", "O001")
	testutil.WriteFile(t, filepath.Join(root, "checks", "C001.md"),
		"---\nid: C001\nscope: {kind: task, id: T001}\nresult: CLEAR\n"+
			"checked_by: {role: checker, session: sess-1}\nchecked_at: '2026-09-14T00:00:00Z'\nissues: [I001]\n---\n\n# Check\n")
	writeV2Issue(t, root, "I001-x.md", "I001", "open", "defect", "")

	problems := CheckProject(root)
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-issue-unpaired-check-link]") {
		t.Fatalf("CheckProject() = %v, want 1 problem naming v2-issue-unpaired-check-link", problems)
	}
}

func TestCheckProject_IssueResolutionRequired(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Issue(t, root, "I001-x.md", "I001", "resolved", "defect", "")

	problems := CheckProject(root)
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-issue-resolution-required]") {
		t.Fatalf("CheckProject() = %v, want 1 problem naming v2-issue-resolution-required", problems)
	}
}

func TestCheckProject_IssueResolutionNotAllowed(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Issue(t, root, "I001-x.md", "I001", "open", "defect",
		"resolution: {disposition: accepted, actor: {role: owner, session: owner-1}, at: '2026-09-15T00:00:00Z', reason: \"risk accepted\"}\n")

	problems := CheckProject(root)
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-issue-resolution-not-allowed]") {
		t.Fatalf("CheckProject() = %v, want 1 problem naming v2-issue-resolution-not-allowed", problems)
	}
}

func TestCheckProject_IssueResolutionMissingProof(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Issue(t, root, "I001-x.md", "I001", "resolved", "defect",
		"resolution: {disposition: verified, actor: {role: checker, session: sess-1}, at: '2026-09-15T00:00:00Z'}\n")

	problems := CheckProject(root)
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-issue-resolution-missing-proof]") {
		t.Fatalf("CheckProject() = %v, want 1 problem naming v2-issue-resolution-missing-proof", problems)
	}
}

func TestCheckProject_IssueResolutionUnusableProof(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Objective(t, root, "O001-ship", "O001", "Ship it")
	writeV2Task(t, root, "O001-ship", "T001-write.md", "T001", "Write it", "O001")
	writeV2Check(t, root, "C001", "{kind: task, id: T001}", "CLEAR", "")
	writeV2Issue(t, root, "I001-x.md", "I001", "resolved", "defect",
		"resolution: {disposition: verified, check: C001, actor: {role: checker, session: sess-1}, at: '2026-09-15T00:00:00Z'}\n")

	problems := CheckProject(root)
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-issue-resolution-unusable-proof]") {
		t.Fatalf("CheckProject() = %v, want 1 problem naming v2-issue-resolution-unusable-proof", problems)
	}
}

func TestCheckProject_IssueResolutionFieldMismatch(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Objective(t, root, "O001-ship", "O001", "Ship it")
	writeV2Task(t, root, "O001-ship", "T001-write.md", "T001", "Write it", "O001")
	writeV2Check(t, root, "C001", "{kind: task, id: T001}", "CLEAR", "")
	writeV2Issue(t, root, "I001-x.md", "I001", "resolved", "defect",
		"checks: [C001]\nresolution: {disposition: accepted, check: C001, actor: {role: owner, session: owner-1}, at: '2026-09-15T00:00:00Z', reason: \"n/a\"}\n")

	problems := CheckProject(root)
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-issue-resolution-field-mismatch]") {
		t.Fatalf("CheckProject() = %v, want 1 problem naming v2-issue-resolution-field-mismatch", problems)
	}
}

// TestV2DiagnosticName_noNewIssueSentinelFallsThrough proves every new E44
// data sentinel — including the write-time-only ones CheckProject can never
// actually surface, since doctor never writes — maps to a name other than the
// generic fallback, so v2DiagnosticName stays a complete registry of every
// named diagnostic.
func TestV2DiagnosticName_noNewIssueSentinelFallsThrough(t *testing.T) {
	sentinels := []error{
		data.ErrV2IssueMalformed,
		data.ErrV2IssueMissingDuplicateTarget,
		data.ErrV2IssueSelfDuplicate,
		data.ErrV2IssueDuplicateCycle,
		data.ErrV2IssueMissingLinkTarget,
		data.ErrV2IssueUnpairedCheckLink,
		data.ErrV2IssueResolutionRequired,
		data.ErrV2IssueResolutionNotAllowed,
		data.ErrV2IssueResolutionMissingProof,
		data.ErrV2IssueResolutionUnusableProof,
		data.ErrV2IssueResolutionFieldMismatch,
		data.ErrV2IssueAlreadyExists,
		data.ErrV2IssueHistoryNotAppendOnly,
	}
	for _, sentinel := range sentinels {
		name := v2DiagnosticName(fmt.Errorf("wrap: %w", sentinel))
		if name == "v2-project-error" {
			t.Errorf("v2DiagnosticName(%v) fell through to the generic fallback", sentinel)
		}
		if V2ProblemRepair(name) == "Review the V2 project diagnostic and fix the reported record" {
			t.Errorf("V2ProblemRepair(%q) fell through to the generic repair suggestion", name)
		}
	}
}

// TestCheckProject_ObjectiveConsistencyDiagnostics proves doctor reports every
// evaluation-level inconsistency InspectObjectiveConsistency finds — a done
// Objective without current integration clearance, and a done Objective with
// an incomplete owned Task — named against the Objective's own source record.
func TestCheckProject_ObjectiveConsistencyDiagnostics(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	testutil.WriteFile(t, filepath.Join(root, "objectives", "O001-ship", "Objective.md"),
		"---\nid: O001\ntitle: \"Ship it\"\nstatus: done\n---\n\n# Ship it\n")
	testutil.WriteFile(t, filepath.Join(root, "objectives", "O001-ship", "tasks", "T001-write.md"),
		"---\nid: T001\ntitle: \"Write it\"\nobjective: O001\nstatus: in_progress\nstage: build\n---\n\n# Write it\n")

	problems := CheckProject(root)
	if len(problems) != 2 {
		t.Fatalf("CheckProject() = %+v, want 2 objective consistency problems", problems)
	}
	if !strings.Contains(problems[0].Message, "[v2-objective-done-without-clearance]") {
		t.Errorf("problems[0].Message = %q, want v2-objective-done-without-clearance", problems[0].Message)
	}
	if !strings.Contains(problems[1].Message, "[v2-objective-done-with-incomplete-task]") {
		t.Errorf("problems[1].Message = %q, want v2-objective-done-with-incomplete-task", problems[1].Message)
	}
	for _, p := range problems {
		if p.Repair == "" {
			t.Errorf("problem %+v missing repair", p)
		}
		if !strings.HasSuffix(p.File, "Objective.md") {
			t.Errorf("problem File = %q, want the Objective's source record path", p.File)
		}
	}
}

// TestCheckProject_IssueVerifiedProofSuperseded proves doctor reports a
// verified Issue whose proof Check has been superseded by a later Check
// recorded against the same scope — an inconsistency a successful load
// cannot refuse by itself.
func TestCheckProject_IssueVerifiedProofSuperseded(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Objective(t, root, "O001-ship", "O001", "Ship it")
	writeV2Task(t, root, "O001-ship", "T001-write.md", "T001", "Write it", "O001")
	writeV2Check(t, root, "C001", "{kind: task, id: T001}", "CLEAR", "")
	writeV2Check(t, root, "C002", "{kind: task, id: T001}", "CLEAR", "C001")
	writeV2Issue(t, root, "I001-flaky.md", "I001", "resolved", "defect",
		"checks: [C001]\nresolution: {disposition: verified, check: C001, actor: {role: checker, session: sess-1}, at: '2026-09-14T00:00:00Z'}\n")

	problems := CheckProject(root)
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-issue-verified-proof-superseded]") {
		t.Fatalf("CheckProject() = %+v, want 1 problem naming v2-issue-verified-proof-superseded", problems)
	}
	if !strings.HasSuffix(problems[0].File, "I001-flaky.md") {
		t.Errorf("problems[0].File = %q, want the Issue's source record path", problems[0].File)
	}
}

// TestCheckProject_OpenAdvisoryIssueDoesNotFailLoad proves an open Issue with
// no resolution is purely advisory: structural and gate evidence decide
// health, so a project carrying only an open Issue alongside otherwise-clean
// records reports no problems.
func TestCheckProject_OpenAdvisoryIssueDoesNotFailLoad(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Objective(t, root, "O001-ship", "O001", "Ship it")
	writeV2Task(t, root, "O001-ship", "T001-write.md", "T001", "Write it", "O001")
	writeV2Issue(t, root, "I001-flaky.md", "I001", "open", "defect", "")

	if problems := CheckProject(root); len(problems) != 0 {
		t.Fatalf("CheckProject() = %v, want no problems: an open advisory Issue must not make the project unhealthy", problems)
	}
}

// TestIssuePostureReport_derivedCountsNoStoredSummary proves Issue posture is
// computed fresh from the loaded index every call, with no summary cached
// between them.
func TestIssuePostureReport_derivedCountsNoStoredSummary(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Issue(t, root, "I001-a.md", "I001", "open", "defect", "")

	posture := IssuePostureReport(root)
	if posture == nil {
		t.Fatal("IssuePostureReport() = nil, want a posture for a V2 project")
	}
	if posture.StatusCounts[data.IssueStatusOpen] != 1 {
		t.Errorf("StatusCounts[open] = %d, want 1", posture.StatusCounts[data.IssueStatusOpen])
	}
	if posture.TypeCounts[data.IssueTypeDefect] != 1 {
		t.Errorf("TypeCounts[defect] = %d, want 1", posture.TypeCounts[data.IssueTypeDefect])
	}

	writeV2Issue(t, root, "I002-b.md", "I002", "open", "drift", "")
	posture2 := IssuePostureReport(root)
	if posture2.StatusCounts[data.IssueStatusOpen] != 2 {
		t.Errorf("StatusCounts[open] after a second Issue = %d, want 2 (recomputed, not cached)", posture2.StatusCounts[data.IssueStatusOpen])
	}
}

func TestIssuePostureReport_v1ProjectReturnsNil(t *testing.T) {
	root := t.TempDir()
	testutil.SetupMinimalProject(t, root, "v1", "E01-foo")

	if got := IssuePostureReport(root); got != nil {
		t.Errorf("IssuePostureReport() = %+v, want nil for a V1 project", got)
	}
}

// TestCheckProject_v2ValidWithChecksNoProblems proves a clean V2 project
// carrying real Check and evidence records — current clearance on a
// technical Task and an owner-accepted Task — reports no problems.
func TestCheckProject_v2ValidWithChecksNoProblems(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Objective(t, root, "O001-ship", "O001", "Ship it")
	writeV2Check(t, root, "C001", "{kind: task, id: T001}", "CLEAR", "")
	testutil.WriteFile(t, filepath.Join(root, "objectives", "O001-ship", "tasks", "T001-write.md"),
		"---\nid: T001\ntitle: \"Write it\"\nobjective: O001\nstatus: done\n"+
			"freshness: {state: current, check: C001, assessed_by: {role: checker, session: s}, assessed_at: '2026-09-14T00:00:00Z', basis: reviewed}\n"+
			"---\n\n# Write it\n")

	if problems := CheckProject(root); len(problems) != 0 {
		t.Fatalf("CheckProject() = %v, want no problems for a clean V2 project with checks", problems)
	}
}

// TestCheckProject_ConsistencyDiagnostics proves doctor reports every
// evaluation-level inconsistency InspectTaskConsistency finds — done without
// current clearance, owner acceptance naming a superseded check, and
// evidence that clears completion while status lags — alongside structural
// load failures, in deterministic (task ID) order, without failing the load.
func TestCheckProject_ConsistencyDiagnostics(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Objective(t, root, "O001-ship", "O001", "Ship it")

	// T001: marked done but no Check was ever recorded.
	testutil.WriteFile(t, filepath.Join(root, "objectives", "O001-ship", "tasks", "T001-alpha.md"),
		"---\nid: T001\ntitle: \"Alpha\"\nobjective: O001\nstatus: done\n---\n\n# Alpha\n")

	// T002: owner accepted C002, but C003 has since superseded it as latest.
	writeV2Check(t, root, "C002", "{kind: task, id: T002}", "CLEAR", "")
	writeV2Check(t, root, "C003", "{kind: task, id: T002}", "CLEAR", "C002")
	testutil.WriteFile(t, filepath.Join(root, "objectives", "O001-ship", "tasks", "T002-beta.md"),
		"---\nid: T002\ntitle: \"Beta\"\nobjective: O001\nstatus: in_progress\nstage: audit\n"+
			"freshness: {state: current, check: C003, assessed_by: {role: checker, session: s}, assessed_at: '2026-09-14T00:00:00Z', basis: rechecked}\n"+
			"owner_validation: {required: true, accepted_check: C002, accepted_by: {role: owner, session: owner-1}}\n"+
			"---\n\n# Beta\n")

	// T003: evidence clears completion (current clearance, no owner
	// validation required) but status was never advanced to done.
	writeV2Check(t, root, "C004", "{kind: task, id: T003}", "CLEAR", "")
	testutil.WriteFile(t, filepath.Join(root, "objectives", "O001-ship", "tasks", "T003-gamma.md"),
		"---\nid: T003\ntitle: \"Gamma\"\nobjective: O001\nstatus: in_progress\nstage: audit\n"+
			"freshness: {state: current, check: C004, assessed_by: {role: checker, session: s}, assessed_at: '2026-09-14T00:00:00Z', basis: reviewed}\n"+
			"---\n\n# Gamma\n")

	problems := CheckProject(root)
	if len(problems) != 3 {
		t.Fatalf("CheckProject() = %+v, want 3 consistency problems", problems)
	}

	wantOrder := []string{"[v2-done-without-clearance]", "[v2-acceptance-superseded]", "[v2-evidence-contradicts-status]"}
	for i, want := range wantOrder {
		if !strings.Contains(problems[i].Message, want) {
			t.Errorf("problems[%d].Message = %q, want containing %q (deterministic task-ID order)", i, problems[i].Message, want)
		}
		if problems[i].Repair == "" {
			t.Errorf("problems[%d].Repair is empty, want a typed repair suggestion", i)
		}
		if !strings.HasSuffix(problems[i].File, ".md") {
			t.Errorf("problems[%d].File = %q, want the Task's source record path", i, problems[i].File)
		}
	}
}

// TestCheckProject_ConsistencyReadOnly proves CheckProject never writes to a
// V2 project carrying Check and evidence records that trip a consistency
// diagnostic: every file's bytes and modification time stay unchanged.
func TestCheckProject_ConsistencyReadOnly(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Objective(t, root, "O001-ship", "O001", "Ship it")
	taskPath := filepath.Join(root, "objectives", "O001-ship", "tasks", "T001-alpha.md")
	testutil.WriteFile(t, taskPath, "---\nid: T001\ntitle: \"Alpha\"\nobjective: O001\nstatus: done\n---\n\n# Alpha\n")

	type snapshot struct {
		bytes []byte
		mtime time.Time
	}
	before := map[string]snapshot{}
	for _, p := range []string{taskPath} {
		info, err := os.Stat(p)
		if err != nil {
			t.Fatal(err)
		}
		raw, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		before[p] = snapshot{bytes: raw, mtime: info.ModTime()}
	}

	if problems := CheckProject(root); len(problems) == 0 {
		t.Fatal("CheckProject() = no problems, want the done-without-clearance diagnostic")
	}

	for p, want := range before {
		info, err := os.Stat(p)
		if err != nil {
			t.Fatal(err)
		}
		raw, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		if string(raw) != string(want.bytes) {
			t.Fatalf("CheckProject() changed %s bytes:\nbefore: %s\nafter: %s", p, want.bytes, raw)
		}
		if !info.ModTime().Equal(want.mtime) {
			t.Fatalf("CheckProject() changed %s mtime: before=%v after=%v", p, want.mtime, info.ModTime())
		}
	}
}

// TestCheckProject_ObjectiveAndIssueConsistencyReadOnly proves CheckProject
// never writes to a V2 project carrying Objective and Issue records that trip
// the new E44 consistency diagnostics: every file's bytes and modification
// time stay unchanged.
func TestCheckProject_ObjectiveAndIssueConsistencyReadOnly(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	objectivePath := filepath.Join(root, "objectives", "O001-ship", "Objective.md")
	testutil.WriteFile(t, objectivePath, "---\nid: O001\ntitle: \"Ship it\"\nstatus: done\n---\n\n# Ship it\n")
	taskPath := writeV2Task(t, root, "O001-ship", "T001-write.md", "T001", "Write it", "O001")
	checkC001Path := writeV2Check(t, root, "C001", "{kind: task, id: T001}", "CLEAR", "")
	checkC002Path := writeV2Check(t, root, "C002", "{kind: task, id: T001}", "CLEAR", "C001")
	issuePath := writeV2Issue(t, root, "I001-flaky.md", "I001", "resolved", "defect",
		"checks: [C001]\nresolution: {disposition: verified, check: C001, actor: {role: checker, session: sess-1}, at: '2026-09-14T00:00:00Z'}\n")

	type snapshot struct {
		bytes []byte
		mtime time.Time
	}
	before := map[string]snapshot{}
	for _, p := range []string{objectivePath, taskPath, checkC001Path, checkC002Path, issuePath} {
		info, err := os.Stat(p)
		if err != nil {
			t.Fatal(err)
		}
		raw, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		before[p] = snapshot{bytes: raw, mtime: info.ModTime()}
	}

	if problems := CheckProject(root); len(problems) == 0 {
		t.Fatal("CheckProject() = no problems, want the objective and issue consistency diagnostics")
	}

	for p, want := range before {
		info, err := os.Stat(p)
		if err != nil {
			t.Fatal(err)
		}
		raw, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		if string(raw) != string(want.bytes) {
			t.Fatalf("CheckProject() changed %s bytes:\nbefore: %s\nafter: %s", p, want.bytes, raw)
		}
		if !info.ModTime().Equal(want.mtime) {
			t.Fatalf("CheckProject() changed %s mtime: before=%v after=%v", p, want.mtime, info.ModTime())
		}
	}
}
