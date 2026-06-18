package doctor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/data"
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
