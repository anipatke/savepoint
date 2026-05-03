package doctor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// --- CheckConfig ---

func TestCheckConfigMissing(t *testing.T) {
	root := t.TempDir()
	err := CheckConfig(root)
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Fatalf("CheckConfig() = %v, want not found error", err)
	}
}

func TestCheckConfigInvalidYAML(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "config.yml"), "theme: [broken")
	err := CheckConfig(root)
	if err == nil || !strings.Contains(err.Error(), "invalid YAML") {
		t.Fatalf("CheckConfig() = %v, want invalid YAML error", err)
	}
}

func TestCheckConfigMissingQualityGates(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "config.yml"), "theme:\n  bg: \"#000\"\n")
	err := CheckConfig(root)
	if err == nil || !strings.Contains(err.Error(), "quality_gates") {
		t.Fatalf("CheckConfig() = %v, want quality_gates error", err)
	}
}

func TestCheckConfigMissingTheme(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "config.yml"), "quality_gates:\n  block_on_failure: true\n")
	err := CheckConfig(root)
	if err == nil || !strings.Contains(err.Error(), "theme") {
		t.Fatalf("CheckConfig() = %v, want theme error", err)
	}
}

func TestCheckConfigValid(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "config.yml"), "quality_gates:\n  block_on_failure: true\ntheme:\n  bg: \"#000\"\n")
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
	writeFile(t, filepath.Join(root, "router.md"), "# no state block")
	err := CheckRouter(root, "")
	if err == nil || !strings.Contains(err.Error(), "invalid state block") {
		t.Fatalf("CheckRouter() = %v, want invalid state block error", err)
	}
}

func TestCheckRouterPreImplementation(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "router.md"), routerContent("pre-implementation", "none", "none"))
	if err := CheckRouter(root, ""); err != nil {
		t.Fatalf("CheckRouter() = %v, want nil", err)
	}
}

func TestCheckRouterMissingReleaseDir(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "router.md"), routerContent("task-building", "v1", "none"))
	err := CheckRouter(root, "")
	if err == nil || !strings.Contains(err.Error(), "release") {
		t.Fatalf("CheckRouter() = %v, want release directory error", err)
	}
}

func TestCheckRouterMissingEpicDir(t *testing.T) {
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "releases", "v1"), 0755)
	writeFile(t, filepath.Join(root, "router.md"), routerContent("task-building", "v1", "E03-foo"))
	err := CheckRouter(root, "")
	if err == nil || !strings.Contains(err.Error(), "epic") {
		t.Fatalf("CheckRouter() = %v, want epic directory error", err)
	}
}

func TestCheckRouterValidWithDirs(t *testing.T) {
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "releases", "v1", "epics", "E03-foo"), 0755)
	writeFile(t, filepath.Join(root, "router.md"), routerContent("task-building", "v1", "E03-foo"))
	if err := CheckRouter(root, ""); err != nil {
		t.Fatalf("CheckRouter() = %v, want nil", err)
	}
}

func TestCheckRouterEpicFilterSkip(t *testing.T) {
	root := t.TempDir()
	// release dir missing — would fail without filter
	writeFile(t, filepath.Join(root, "router.md"), routerContent("task-building", "v1", "E03-foo"))
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
	os.MkdirAll(filepath.Join(root, "releases"), 0755)
	problems := CheckStructure(root, "")
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "no release directories found") {
		t.Fatalf("CheckStructure() = %v, want no releases error", problems)
	}
}

func TestCheckStructure_MissingReleasePRD(t *testing.T) {
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "releases", "v1", "epics"), 0755)
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
	os.MkdirAll(filepath.Join(root, "releases", "v1", "epics"), 0755)
	writeReleasePRD(t, filepath.Join(root, "releases", "v1"))
	problems := CheckStructure(root, "")
	for _, p := range problems {
		if strings.Contains(p.File, "v1-PRD.md") {
			t.Fatalf("CheckStructure() unexpected PRD problem: %v", p)
		}
	}
}

func TestCheckStructure_ReleasePRDCorruptYAML(t *testing.T) {
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "releases", "v1", "epics"), 0755)
	writeFile(t, filepath.Join(root, "releases", "v1", "v1-PRD.md"), "---\ntype: [broken\n---\n")
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
	os.MkdirAll(epicPath, 0755)
	writeReleasePRD(t, releasePath)
	writeFile(t, filepath.Join(epicPath, "E01-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E01: Foo\n")
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
	os.MkdirAll(epicPath, 0755)
	writeReleasePRD(t, releasePath)
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

func TestCheckStructure_ValidTask(t *testing.T) {
	root := t.TempDir()
	releasePath := filepath.Join(root, "releases", "v1")
	epicPath := filepath.Join(releasePath, "epics", "E01-foo")
	tasksPath := filepath.Join(epicPath, "tasks")
	os.MkdirAll(tasksPath, 0755)
	writeReleasePRD(t, releasePath)
	writeFile(t, filepath.Join(epicPath, "E01-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E01: Foo\n")
	writeFile(t, filepath.Join(tasksPath, "T001-task.md"), "---\nid: E01-foo/T001-task\nstatus: planned\nobjective: \"Do the thing\"\ndepends_on: []\n---\n\n# T001: Task\n\n## Acceptance Criteria\n\n- It works\n")
	problems := CheckStructure(root, "")
	if len(problems) > 0 {
		t.Fatalf("CheckStructure() = %v, want no problems", problems)
	}
}

func TestCheckStructure_TaskMissingRequiredField(t *testing.T) {
	root := t.TempDir()
	releasePath := filepath.Join(root, "releases", "v1")
	epicPath := filepath.Join(releasePath, "epics", "E01-foo")
	tasksPath := filepath.Join(epicPath, "tasks")
	os.MkdirAll(tasksPath, 0755)
	writeReleasePRD(t, releasePath)
	writeFile(t, filepath.Join(epicPath, "E01-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E01: Foo\n")
	writeFile(t, filepath.Join(tasksPath, "T001-task.md"), "---\nid: E01-foo/T001-task\nobjective: \"Do the thing\"\n---\n\n# T001: Task\n\n## Acceptance Criteria\n\n- It works\n")
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
	os.MkdirAll(tasksPath, 0755)
	writeReleasePRD(t, releasePath)
	writeFile(t, filepath.Join(epicPath, "E01-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E01: Foo\n")
	writeFile(t, filepath.Join(tasksPath, "T001-task.md"), "---\nid: E01-foo/T001-task\nstatus: planned\nobjective: \"Do the thing\"\n---\n\n# T001: Task\n")
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
	os.MkdirAll(tasksPath, 0755)
	writeReleasePRD(t, releasePath)
	writeFile(t, filepath.Join(epicPath, "E01-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E01: Foo\n")
	writeFile(t, filepath.Join(tasksPath, "T001-task.md"), "---\nid: \"unclosed\nstatus: planned\n---\n")
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
	os.MkdirAll(epic1Path, 0755)
	os.MkdirAll(epic2Path, 0755)
	writeReleasePRD(t, releasePath)
	writeFile(t, filepath.Join(epic1Path, "E01-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E01: Foo\n")
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
	os.MkdirAll(epicPath, 0755)
	writeReleasePRD(t, releasePath)
	writeFile(t, filepath.Join(epicPath, "E01-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E01: Foo\n")
	os.MkdirAll(filepath.Join(epicPath, "tasks"), 0755)
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
	os.MkdirAll(tasksPath, 0755)
	writeReleasePRD(t, releasePath)
	writeFile(t, filepath.Join(epicPath, "E01-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E01: Foo\n")

	// Two epics, same task ID
	epic2Path := filepath.Join(releasePath, "epics", "E02-bar")
	tasks2Path := filepath.Join(epic2Path, "tasks")
	os.MkdirAll(tasks2Path, 0755)
	writeFile(t, filepath.Join(epic2Path, "E02-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E02: Bar\n")

	writeFile(t, filepath.Join(tasksPath, "T001-task.md"), "---\nid: E01-foo/T001-task\nstatus: planned\nobjective: \"A\"\ndepends_on: []\n---\n\n# T001\n\n## Acceptance Criteria\n\n- it works\n")
	writeFile(t, filepath.Join(tasks2Path, "T001-task.md"), "---\nid: E01-foo/T001-task\nstatus: planned\nobjective: \"A\"\ndepends_on: []\n---\n\n# T001\n\n## Acceptance Criteria\n\n- it works\n")

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
	os.MkdirAll(filepath.Join(releasePath, "epics", "E01-foo", "tasks"), 0755)
	os.MkdirAll(filepath.Join(releasePath, "epics", "E02-bar", "tasks"), 0755)
	writeReleasePRD(t, releasePath)
	writeFile(t, filepath.Join(releasePath, "epics", "E01-foo", "E01-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E01: Foo\n")
	writeFile(t, filepath.Join(releasePath, "epics", "E02-bar", "E02-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E02: Bar\n")

	// E02 has a missing dep — should be invisible with filter
	taskE2 := `---\nid: E02-bar/T001-task\nstatus: planned\nobjective: \"B\"\ndepends_on: [\"E02-bar/T999-nonexistent\"]\n---\n\n# T001\n\n## Acceptance Criteria\n\n- it works\n`
	writeFile(t, filepath.Join(releasePath, "epics", "E01-foo", "tasks", "T001-task.md"), "---\nid: E01-foo/T001-task\nstatus: planned\nobjective: \"A\"\ndepends_on: []\n---\n\n# T001\n\n## Acceptance Criteria\n\n- it works\n")
	writeFile(t, filepath.Join(releasePath, "epics", "E02-bar", "tasks", "T001-task.md"), strings.ReplaceAll(taskE2, "\\n", "\n"))

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
	os.MkdirAll(filepath.Join(root, "releases", "v1", "epics", "E01-foo"), 0755)
	writeFile(t, filepath.Join(root, "router.md"), routerContent("task-building", "v1", "E01-foo"))
	problems := CheckAuditState(root)
	if len(problems) > 0 {
		t.Fatalf("CheckAuditState() = %v, want no problems", problems)
	}
}

func TestCheckAuditState_MatchesRouter(t *testing.T) {
	root := t.TempDir()
	epicPath := filepath.Join(root, "releases", "v1", "epics", "E01-foo")
	os.MkdirAll(epicPath, 0755)
	writeFile(t, filepath.Join(epicPath, "E01-Audit.md"), "---\ntype: audit-findings\n---\n\n# Audit\n")
	writeFile(t, filepath.Join(root, "router.md"), routerContent("audit-pending", "v1", "E01-foo"))
	problems := CheckAuditState(root)
	if len(problems) > 0 {
		t.Fatalf("CheckAuditState() = %v, want no problems when router matches", problems)
	}
}

func TestCheckAuditState_ProposalWithoutPending(t *testing.T) {
	root := t.TempDir()
	epicPath := filepath.Join(root, "releases", "v1", "epics", "E01-foo")
	os.MkdirAll(epicPath, 0755)
	writeFile(t, filepath.Join(epicPath, "E01-Audit.md"), "---\ntype: audit-findings\n---\n\n# Audit\n")
	writeFile(t, filepath.Join(root, "router.md"), routerContent("task-building", "v1", "E01-foo"))
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
	os.MkdirAll(epic1Path, 0755)
	os.MkdirAll(epic2Path, 0755)
	writeFile(t, filepath.Join(epic1Path, "E01-Audit.md"), "---\ntype: audit-findings\n---\n\n# Audit\n")
	writeFile(t, filepath.Join(root, "router.md"), routerContent("audit-pending", "v1", "E02-bar"))
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
	os.MkdirAll(epic1Path, 0755)
	os.MkdirAll(epic2Path, 0755)
	writeFile(t, filepath.Join(epic1Path, "E01-Audit.md"), "---\ntype: audit-findings\n---\n\n# Audit\n")
	writeFile(t, filepath.Join(epic2Path, "E02-Audit.md"), "---\ntype: audit-findings\n---\n\n# Audit\n")
	writeFile(t, filepath.Join(root, "router.md"), routerContent("task-building", "v1", "E03-baz"))
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
	os.MkdirAll(tasksPath, 0755)
	writeReleasePRD(t, releasePath)
	writeFile(t, filepath.Join(epicPath, "E01-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E01: Foo\n")
	writeFile(t, filepath.Join(tasksPath, "T001-task.md"), "---\nid: E02-bar/T001-task\nstatus: planned\nobjective: \"Task\"\ndepends_on: []\n---\n\n# T001\n\n## Acceptance Criteria\n\n- it works\n")
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
	os.MkdirAll(filepath.Join(epic1Path, "tasks"), 0755)
	os.MkdirAll(filepath.Join(epic2Path, "tasks"), 0755)
	writeReleasePRD(t, releasePath)
	writeFile(t, filepath.Join(epic1Path, "E01-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E01: Foo\n")
	writeFile(t, filepath.Join(epic2Path, "E02-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# E02: Bar\n")
	writeFile(t, filepath.Join(epic1Path, "tasks", "T001-task.md"), "---\nid: E02-bar/T001-task\nstatus: planned\nobjective: \"Task\"\ndepends_on: []\n---\n\n# T001\n\n## Acceptance Criteria\n\n- it works\n")
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
	writeFile(t, filepath.Join(tasksPath, "T002-bad.md"), "---\nstatus: planned\nobjective: \"No ID\"\ndepends_on: []\n---\n\n# T002\n\n## Acceptance Criteria\n\n- it works\n")
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
	os.MkdirAll(tasksPath, 0755)

	prefix := epicID
	if idx := strings.IndexByte(epicID, '-'); idx != -1 {
		prefix = epicID[:idx]
	}

	writeReleasePRD(t, releasePath)
	writeFile(t, filepath.Join(epicPath, prefix+"-Detail.md"), "---\ntype: epic-design\nstatus: planned\n---\n\n# Epic\n")

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
		writeFile(t, filepath.Join(tasksPath, fmt.Sprintf("T%03d-task.md", i+1)), content)
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func writeReleasePRD(t *testing.T, releasePath string) {
	t.Helper()
	writeFile(t, filepath.Join(releasePath, filepath.Base(releasePath)+"-PRD.md"), "---\ntype: project-prd\nstatus: active\n---\n\n# Release\n")
}

func routerContent(state, release, epic string) string {
	return "## Current state\n\n```yaml\nstate: " + state + "\nrelease: " + release + "\nepic: " + epic + "\ntask: none\nnext_action: \"\"\n```\n"
}
