package data

import (
	"errors"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/testutil"
	"gopkg.in/yaml.v3"
)

func TestLoadProject(t *testing.T) {
	cases := []struct {
		name          string
		configContent string
		omitConfig    bool
		wantVersion   SchemaVersion
		wantErr       error
	}{
		{
			name:        "absent config selects transitional V1",
			omitConfig:  true,
			wantVersion: SchemaVersionV1,
		},
		{
			name:          "absent schema_version selects transitional V1",
			configContent: "theme:\n  bg: \"#000000\"\n",
			wantVersion:   SchemaVersionV1,
		},
		{
			name:          "explicit version 2 dispatches to the V2 index loader",
			configContent: "schema_version: 2\n",
			wantVersion:   SchemaVersionV2,
		},
		{
			name:          "malformed explicit version fails named, does not fall back to V1",
			configContent: "schema_version: nope\n",
			wantErr:       ErrMalformedSchemaVersion,
		},
		{
			name:          "unsupported explicit version fails named, does not fall back to V1",
			configContent: "schema_version: 99\n",
			wantErr:       ErrUnsupportedSchemaVersion,
		},
		{
			name:          "unrelated version-shaped fields do not affect dispatch",
			configContent: "theme:\n  bg: \"#000000\"\n",
			wantVersion:   SchemaVersionV1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			if !tc.omitConfig {
				testutil.WriteFile(t, filepath.Join(root, "config.yml"), tc.configContent)
			}
			// Package/release/upgrade-manifest content must never influence
			// schema selection, even when it carries its own version fields.
			testutil.WriteFile(t, filepath.Join(root, ".upgrade-manifest.yml"), "schema_version: 2\nversion: 1.3.1\n")
			testutil.MkdirAll(t, filepath.Join(root, "releases", "v2", "epics"))

			project, err := LoadProject(root)

			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("LoadProject() error = %v, want wrapping %v", err, tc.wantErr)
				}
				if project != nil {
					t.Errorf("LoadProject() project = %+v, want nil on error", project)
				}
				return
			}

			if err != nil {
				t.Fatalf("LoadProject() unexpected error = %v", err)
			}
			if project.SchemaVersion != tc.wantVersion {
				t.Errorf("LoadProject() SchemaVersion = %v, want %v", project.SchemaVersion, tc.wantVersion)
			}
			switch tc.wantVersion {
			case SchemaVersionV2:
				if project.V2 == nil {
					t.Fatal("LoadProject() V2 index = nil, want non-nil for V2 dispatch")
				}
			default:
				if project.V1 == nil {
					t.Fatal("LoadProject() V1 discover adapter = nil, want non-nil for transitional V1 dispatch")
				}
			}
		})
	}
}

// TestLoadProjectV1DiscoveryUnchanged proves the schema dispatch boundary
// does not alter existing V1 Discover behavior: a project reachable via
// LoadProject's V1 adapter still lists releases/epics/tasks exactly as
// calling Discover directly would.
func TestLoadProjectV1DiscoveryUnchanged(t *testing.T) {
	root := t.TempDir()
	savepointRoot := filepath.Join(root, ".savepoint")
	testutil.SetupMinimalProject(t, savepointRoot, "v1", "E01-example")

	project, err := LoadProject(savepointRoot)
	if err != nil {
		t.Fatalf("LoadProject() error = %v", err)
	}
	if project.SchemaVersion != SchemaVersionV1 {
		t.Fatalf("LoadProject() SchemaVersion = %v, want SchemaVersionV1", project.SchemaVersion)
	}

	viaProject, err := project.V1.ListEpics(savepointRoot, "v1")
	if err != nil {
		t.Fatalf("project.V1.ListEpics() error = %v", err)
	}

	direct, err := NewDiscover().ListEpics(savepointRoot, "v1")
	if err != nil {
		t.Fatalf("NewDiscover().ListEpics() error = %v", err)
	}

	if len(viaProject) != len(direct) || len(viaProject) != 1 || viaProject[0].ID != direct[0].ID {
		t.Errorf("ListEpics() via project = %+v, direct = %+v, want matching single-epic results", viaProject, direct)
	}
}

func TestLoadV2Index_valid(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O001-first", "O001", "First objective")
	writeV2ObjectiveFixture(t, root, "O002-second", "O002", "Second objective", "O001")
	writeV2TaskFixture(t, root, "O001-first", "T001-alpha.md", "T001", "Alpha", "O001")
	writeV2TaskFixture(t, root, "O001-first", "T002-beta.md", "T002", "Beta", "O001")
	writeV2TaskFixture(t, root, "O002-second", "T003-gamma.md", "T003", "Gamma", "O002")

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}

	if len(index.Objectives) != 2 || len(index.Tasks) != 3 {
		t.Fatalf("LoadV2Index() Objectives = %d, Tasks = %d, want 2, 3", len(index.Objectives), len(index.Tasks))
	}

	owned := index.ObjectiveTasks["O001"]
	if len(owned) != 2 || owned[0] != "T001" || owned[1] != "T002" {
		t.Errorf("ObjectiveTasks[O001] = %v, want [T001 T002] in sorted order", owned)
	}
	if got := index.ObjectiveTasks["O002"]; len(got) != 1 || got[0] != "T003" {
		t.Errorf("ObjectiveTasks[O002] = %v, want [T003]", got)
	}
}

func TestLoadV2Index_emptyProject(t *testing.T) {
	root := t.TempDir()

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v, want nil for a project with no Objectives yet", err)
	}
	if len(index.Objectives) != 0 || len(index.Tasks) != 0 || len(index.ObjectiveTasks) != 0 {
		t.Errorf("LoadV2Index() = %+v, want an empty index", index)
	}
}

// TestLoadV2Index_movedTaskRetainsOwnership proves the assembled index keeps
// explicit Objective ownership from the Task record's own field even when
// the Task file is filed under a different Objective's tasks/ directory.
func TestLoadV2Index_movedTaskRetainsOwnership(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O001-first", "O001", "First objective")
	writeV2ObjectiveFixture(t, root, "O002-second", "O002", "Second objective")
	writeV2TaskFixture(t, root, "O002-second", "T001-alpha.md", "T001", "Alpha", "O001")

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}

	if got := index.ObjectiveTasks["O001"]; len(got) != 1 || got[0] != "T001" {
		t.Errorf("ObjectiveTasks[O001] = %v, want [T001] despite the file living under O002-second/tasks", got)
	}
	if len(index.ObjectiveTasks["O002"]) != 0 {
		t.Errorf("ObjectiveTasks[O002] = %v, want empty: the containing directory does not grant ownership", index.ObjectiveTasks["O002"])
	}
}

func TestLoadV2Index_missingOwner(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O001-first", "O001", "First objective")
	writeV2TaskFixture(t, root, "O001-first", "T001-alpha.md", "T001", "Alpha", "O999")

	_, err := LoadV2Index(root)
	if !errors.Is(err, ErrV2MissingOwner) {
		t.Fatalf("LoadV2Index() error = %v, want ErrV2MissingOwner", err)
	}
}

// TestLoadV2Index_integratedProjectScenario combines the behaviors E42's
// tasks proved in isolation into one project, as the completed-epic
// boundary: distinct human titles vs O###/T### identity, valid Objectives
// and Tasks, a Task moved to a different Objective's tasks/ directory that
// still resolves ownership from its own field, unknown frontmatter content
// surviving the full discover-then-decode path, and deterministic
// (repeat-run-stable) indexing.
func TestLoadV2Index_integratedProjectScenario(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O001-ship", "O001", "Ship the thing")
	writeV2ObjectiveFixture(t, root, "O002-polish", "O002", "Polish the thing")
	writeV2TaskFixture(t, root, "O001-ship", "T001-write-code.md", "T001", "Write the code", "O001")

	// T002 carries unknown top-level and nested-mapping fields; the source
	// document must retain them even though DecodeTaskV2 projects only the
	// typed fields it knows about.
	testutil.WriteFile(t, filepath.Join(root, v2ObjectivesDirName, "O001-ship", v2TasksDirName, "T002-review-code.md"),
		"---\nid: T002\ntitle: \"Review the code\"\nobjective: O001\nstatus: planned\nnotes: kept for reviewers\nmetadata:\n  reviewer:\n    name: sam\n---\n\n# Review the code\n")

	// T003 is owned by O002 but filed under O001's tasks/ directory: a moved
	// Task whose ownership must come from its own objective field, not its
	// containing directory.
	writeV2TaskFixture(t, root, "O001-ship", "T003-ship-it.md", "T003", "Ship it", "O002")

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}

	if len(index.Objectives) != 2 || len(index.Tasks) != 3 {
		t.Fatalf("LoadV2Index() Objectives = %d, Tasks = %d, want 2, 3", len(index.Objectives), len(index.Tasks))
	}

	for _, task := range index.Tasks {
		if task.Title == task.Objective || task.Title == task.ID {
			t.Errorf("task %s Title = %q, want a human display title distinct from its ID/objective identity", task.ID, task.Title)
		}
	}

	if got := index.ObjectiveTasks["O001"]; len(got) != 2 || got[0] != "T001" || got[1] != "T002" {
		t.Errorf("ObjectiveTasks[O001] = %v, want [T001 T002]", got)
	}
	if got := index.ObjectiveTasks["O002"]; len(got) != 1 || got[0] != "T003" {
		t.Errorf("ObjectiveTasks[O002] = %v, want [T003] despite T003's file living under O001-ship/tasks", got)
	}

	t2 := index.Tasks["T002"]
	t2Mapping := t2.Source.Frontmatter.Content[0]
	if notes, ok := mappingFieldValue(t2Mapping, "notes"); !ok || notes != "kept for reviewers" {
		t.Errorf("T002 preserved notes field = %q, ok = %v, want \"kept for reviewers\", true", notes, ok)
	}
	metadataNode := findMappingChild(t2Mapping, "metadata")
	if metadataNode == nil {
		t.Fatal("T002 lost its unknown nested metadata field")
	}
	reviewerNode := findMappingChild(metadataNode, "reviewer")
	if reviewerNode == nil {
		t.Fatal("T002 lost its unknown nested metadata.reviewer field")
	}
	if name, ok := mappingFieldValue(reviewerNode, "name"); !ok || name != "sam" {
		t.Errorf("T002 preserved metadata.reviewer.name = %q, ok = %v, want \"sam\", true", name, ok)
	}

	reindex, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() second run error = %v", err)
	}
	if !reflect.DeepEqual(index.ObjectiveTasks, reindex.ObjectiveTasks) {
		t.Errorf("LoadV2Index() ObjectiveTasks not deterministic across runs: %v vs %v", index.ObjectiveTasks, reindex.ObjectiveTasks)
	}
}

// findMappingChild returns the value node for key in mapping, or nil if not
// present. mappingFieldValue (write.go) only returns scalar values; this
// walks the same Content pairs to reach a nested mapping node instead.
func findMappingChild(mapping *yaml.Node, key string) *yaml.Node {
	for i := 0; i < len(mapping.Content)-1; i += 2 {
		if mapping.Content[i].Value == key {
			return mapping.Content[i+1]
		}
	}
	return nil
}

// TestLoadV2Index_taskCycleFromDisk proves the full discover-then-validate
// path surfaces a Task reference cycle decoded from real files, not just the
// in-memory graph algorithm covered directly in dependency_test.go.
func TestLoadV2Index_taskCycleFromDisk(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O001-first", "O001", "First objective")
	testutil.WriteFile(t, filepath.Join(root, v2ObjectivesDirName, "O001-first", v2TasksDirName, "T001-alpha.md"),
		"---\nid: T001\ntitle: \"Alpha\"\nobjective: O001\nstatus: planned\ndepends_on: [{task: T002}]\n---\n\n# Alpha\n")
	testutil.WriteFile(t, filepath.Join(root, v2ObjectivesDirName, "O001-first", v2TasksDirName, "T002-beta.md"),
		"---\nid: T002\ntitle: \"Beta\"\nobjective: O001\nstatus: planned\ndepends_on: [{task: T001}]\n---\n\n# Beta\n")

	_, err := LoadV2Index(root)
	if !errors.Is(err, ErrV2DependencyCycle) {
		t.Fatalf("LoadV2Index() error = %v, want ErrV2DependencyCycle", err)
	}
}

func TestLoadV2Index_checksValid(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O001-first", "O001", "First objective")
	writeV2TaskFixture(t, root, "O001-first", "T001-alpha.md", "T001", "Alpha", "O001")
	writeV2CheckFixture(t, root, "C001-first.md", "C001", "task", "T001")

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}

	if len(index.Checks) != 1 || index.Checks["C001"] == nil {
		t.Fatalf("Checks = %+v, want C001", index.Checks)
	}
	if got := index.ScopeChecks["T001"]; len(got) != 1 || got[0] != "C001" {
		t.Errorf("ScopeChecks[T001] = %v, want [C001]", got)
	}
	if got := index.LatestCheck["T001"]; got != "C001" {
		t.Errorf("LatestCheck[T001] = %q, want C001", got)
	}
}

func TestLoadV2Index_checksAbsentDirLoadsEmpty(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O001-first", "O001", "First objective")
	writeV2TaskFixture(t, root, "O001-first", "T001-alpha.md", "T001", "Alpha", "O001")

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v, want nil for a project with no checks/ yet", err)
	}
	if len(index.Checks) != 0 || len(index.ScopeChecks) != 0 || len(index.LatestCheck) != 0 {
		t.Errorf("Checks = %+v, ScopeChecks = %+v, LatestCheck = %+v, want all empty", index.Checks, index.ScopeChecks, index.LatestCheck)
	}
}

func TestLoadV2Index_checkMissingScopeTarget(t *testing.T) {
	root := t.TempDir()
	writeV2CheckFixture(t, root, "C001-first.md", "C001", "task", "T999")

	_, err := LoadV2Index(root)
	if !errors.Is(err, ErrV2CheckMissingScopeTarget) {
		t.Fatalf("LoadV2Index() error = %v, want ErrV2CheckMissingScopeTarget", err)
	}
}

// TestLoadV2Index_checkSupersedesChain proves a valid rerun chain resolves
// to the newest Check as the target's latest, while both Checks remain
// individually addressable by ID.
func TestLoadV2Index_checkSupersedesChain(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O001-first", "O001", "First objective")
	writeV2TaskFixture(t, root, "O001-first", "T001-alpha.md", "T001", "Alpha", "O001")
	writeV2CheckFixture(t, root, "C001-first.md", "C001", "task", "T001")
	testutil.WriteFile(t, filepath.Join(root, v2ChecksDirName, "C002-rerun.md"),
		"---\nid: C002\nscope: {kind: task, id: T001}\nresult: CLEAR\nchecked_by: {role: checker, session: sess-2}\nchecked_at: '2026-09-15T00:00:00Z'\nsupersedes: C001\n---\n\n# Check\n")

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}

	if got := index.ScopeChecks["T001"]; len(got) != 2 || got[0] != "C001" || got[1] != "C002" {
		t.Errorf("ScopeChecks[T001] = %v, want [C001 C002]", got)
	}
	if got := index.LatestCheck["T001"]; got != "C002" {
		t.Errorf("LatestCheck[T001] = %q, want C002", got)
	}
}

func TestLoadV2Index_checkOrderingAcrossC999C1000Boundary(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O001-first", "O001", "First objective")
	writeV2TaskFixture(t, root, "O001-first", "T001-alpha.md", "T001", "Alpha", "O001")
	writeV2TaskFixture(t, root, "O001-first", "T002-beta.md", "T002", "Beta", "O001")
	writeV2CheckFixture(t, root, "C999-older.md", "C999", "task", "T001")
	writeV2CheckFixture(t, root, "C1000-newer.md", "C1000", "task", "T001")

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	if got := index.ScopeChecks["T001"]; len(got) != 2 || got[0] != "C999" || got[1] != "C1000" {
		t.Fatalf("ScopeChecks[T001] = %v, want [C999 C1000] in numeric order", got)
	}
	if got := index.LatestCheck["T001"]; got != "C1000" {
		t.Fatalf("LatestCheck[T001] = %q, want C1000", got)
	}

	t001 := index.Tasks["T001"]
	t001.Evidence = &Evidence{Freshness: &Freshness{
		State: FreshnessCurrent, Check: "C1000",
		AssessedBy: Actor{Role: ActorRoleChecker, Session: "checker-1"},
		Basis:      "reviewed the newer check",
	}}
	t001.Status = ColumnDone
	if got := ResolveClearance(index, "T001"); got.State != ClearanceCurrent || got.Check != "C1000" {
		t.Fatalf("ResolveClearance(T001) = %+v, want current/C1000", got)
	}
	index.Tasks["T002"].Status = ColumnPlanned
	index.Tasks["T002"].DependsOn = []TaskDependencyV2{{Task: "T001", Requires: TaskDependencyClear}}
	if dependency := ResolveTaskDependencyV2(index, index.Tasks["T002"].DependsOn[0]); !dependency.Satisfied {
		t.Fatalf("ResolveTaskDependencyV2(T002 -> T001) = %+v, want satisfied", dependency)
	}

	t001.Status = ColumnInProgress
	t001.Stage = StageAudit
	if decision := ResolveTaskCompletion(index, "T001"); !decision.Allowed {
		t.Fatalf("ResolveTaskCompletion(T001) = %+v, want allowed from latest C1000", decision)
	}
}

func TestLoadV2Index_checkSupersedesMissingTarget(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O001-first", "O001", "First objective")
	writeV2TaskFixture(t, root, "O001-first", "T001-alpha.md", "T001", "Alpha", "O001")
	testutil.WriteFile(t, filepath.Join(root, v2ChecksDirName, "C002-rerun.md"),
		"---\nid: C002\nscope: {kind: task, id: T001}\nresult: CLEAR\nchecked_by: {role: checker, session: sess-2}\nchecked_at: '2026-09-15T00:00:00Z'\nsupersedes: C999\n---\n\n# Check\n")

	_, err := LoadV2Index(root)
	if !errors.Is(err, ErrV2CheckMissingReference) {
		t.Fatalf("LoadV2Index() error = %v, want ErrV2CheckMissingReference", err)
	}
}

func TestLoadV2Index_checkSupersedesScopeMismatch(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O001-first", "O001", "First objective")
	writeV2TaskFixture(t, root, "O001-first", "T001-alpha.md", "T001", "Alpha", "O001")
	writeV2TaskFixture(t, root, "O001-first", "T002-beta.md", "T002", "Beta", "O001")
	writeV2CheckFixture(t, root, "C001-first.md", "C001", "task", "T001")
	testutil.WriteFile(t, filepath.Join(root, v2ChecksDirName, "C002-rerun.md"),
		"---\nid: C002\nscope: {kind: task, id: T002}\nresult: CLEAR\nchecked_by: {role: checker, session: sess-2}\nchecked_at: '2026-09-15T00:00:00Z'\nsupersedes: C001\n---\n\n# Check\n")

	_, err := LoadV2Index(root)
	if !errors.Is(err, ErrV2CheckSupersedesConflict) {
		t.Fatalf("LoadV2Index() error = %v, want ErrV2CheckSupersedesConflict", err)
	}
}

// TestLoadV2Index_checkSupersedesFork proves two Checks superseding the same
// target Check is rejected, not silently treated as two valid chain heads.
func TestLoadV2Index_checkSupersedesFork(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O001-first", "O001", "First objective")
	writeV2TaskFixture(t, root, "O001-first", "T001-alpha.md", "T001", "Alpha", "O001")
	writeV2CheckFixture(t, root, "C001-first.md", "C001", "task", "T001")
	testutil.WriteFile(t, filepath.Join(root, v2ChecksDirName, "C002-rerun.md"),
		"---\nid: C002\nscope: {kind: task, id: T001}\nresult: CLEAR\nchecked_by: {role: checker, session: sess-2}\nchecked_at: '2026-09-15T00:00:00Z'\nsupersedes: C001\n---\n\n# Check\n")
	testutil.WriteFile(t, filepath.Join(root, v2ChecksDirName, "C003-also-rerun.md"),
		"---\nid: C003\nscope: {kind: task, id: T001}\nresult: NEEDS WORK\nchecked_by: {role: checker, session: sess-3}\nchecked_at: '2026-09-15T01:00:00Z'\nsupersedes: C001\n---\n\n# Check\n")

	_, err := LoadV2Index(root)
	if !errors.Is(err, ErrV2CheckSupersedesConflict) {
		t.Fatalf("LoadV2Index() error = %v, want ErrV2CheckSupersedesConflict", err)
	}
}

func TestLoadV2Index_checkSupersedesCycle(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O001-first", "O001", "First objective")
	writeV2TaskFixture(t, root, "O001-first", "T001-alpha.md", "T001", "Alpha", "O001")
	testutil.WriteFile(t, filepath.Join(root, v2ChecksDirName, "C001-first.md"),
		"---\nid: C001\nscope: {kind: task, id: T001}\nresult: CLEAR\nchecked_by: {role: checker, session: sess-1}\nchecked_at: '2026-09-14T00:00:00Z'\nsupersedes: C002\n---\n\n# Check\n")
	testutil.WriteFile(t, filepath.Join(root, v2ChecksDirName, "C002-second.md"),
		"---\nid: C002\nscope: {kind: task, id: T001}\nresult: CLEAR\nchecked_by: {role: checker, session: sess-2}\nchecked_at: '2026-09-15T00:00:00Z'\nsupersedes: C001\n---\n\n# Check\n")

	_, err := LoadV2Index(root)
	if !errors.Is(err, ErrV2CheckSupersedesConflict) {
		t.Fatalf("LoadV2Index() error = %v, want ErrV2CheckSupersedesConflict", err)
	}
}

func TestLoadV2Index_checkDeterministicDiagnosticOrder(t *testing.T) {
	root := t.TempDir()
	writeV2CheckFixture(t, root, "C001-first.md", "C001", "task", "T998")
	writeV2CheckFixture(t, root, "C002-second.md", "C002", "task", "T999")

	_, err1 := LoadV2Index(root)
	_, err2 := LoadV2Index(root)
	if err1 == nil || err2 == nil || err1.Error() != err2.Error() {
		t.Fatalf("LoadV2Index() errors not deterministic across runs: %v vs %v", err1, err2)
	}
	if !errors.Is(err1, ErrV2CheckMissingScopeTarget) {
		t.Fatalf("LoadV2Index() error = %v, want ErrV2CheckMissingScopeTarget", err1)
	}
}

// TestLoadV2Index_evidenceReferencesResolve proves last_check,
// freshness.check, owner_validation.accepted_check, and exception.check on
// a Task record all resolve against Checks discovered in the same index.
func TestLoadV2Index_evidenceReferencesResolve(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O001-first", "O001", "First objective")
	writeV2CheckFixture(t, root, "C001-first.md", "C001", "task", "T001")
	testutil.WriteFile(t, filepath.Join(root, v2ObjectivesDirName, "O001-first", v2TasksDirName, "T001-alpha.md"),
		"---\nid: T001\ntitle: \"Alpha\"\nobjective: O001\nstatus: planned\n"+
			"last_check: C001\n"+
			"freshness: {state: current, check: C001, assessed_by: {role: checker, session: sess-1}, assessed_at: '2026-09-15T00:00:00Z', basis: reviewed}\n"+
			"owner_validation: {required: true, accepted_check: C001, accepted_by: {role: owner, session: owner-1}}\n"+
			"exception: {requirements: [TEST-08], reason: waived, owner: ani, recorded_at: '2026-09-15T00:00:00Z', check: C001}\n"+
			"---\n\n# Alpha\n")

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}

	task := index.Tasks["T001"]
	if task.Evidence == nil {
		t.Fatal("Tasks[T001].Evidence = nil, want decoded evidence")
	}
	if task.Evidence.LastCheck != "C001" {
		t.Errorf("Evidence.LastCheck = %q, want C001", task.Evidence.LastCheck)
	}
}

// TestLoadV2Index_evidenceMissingReference proves each of the four evidence
// Check reference fields fails the load closed with ErrV2EvidenceMissingReference
// when the named Check does not exist in the index.
func TestLoadV2Index_evidenceMissingReference(t *testing.T) {
	tests := []struct {
		name  string
		field string
	}{
		{"last_check", "last_check: C999\n"},
		{"freshness.check", "freshness: {state: current, check: C999, assessed_by: {role: checker, session: sess-1}, assessed_at: '2026-09-15T00:00:00Z', basis: reviewed}\n"},
		{"owner_validation.accepted_check", "owner_validation: {required: true, accepted_check: C999, accepted_by: {role: owner, session: owner-1}}\n"},
		{"exception.check", "exception: {requirements: [TEST-08], reason: waived, owner: ani, recorded_at: '2026-09-15T00:00:00Z', check: C999}\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			writeV2ObjectiveFixture(t, root, "O001-first", "O001", "First objective")
			testutil.WriteFile(t, filepath.Join(root, v2ObjectivesDirName, "O001-first", v2TasksDirName, "T001-alpha.md"),
				"---\nid: T001\ntitle: \"Alpha\"\nobjective: O001\nstatus: planned\n"+tt.field+"---\n\n# Alpha\n")

			_, err := LoadV2Index(root)
			if !errors.Is(err, ErrV2EvidenceMissingReference) {
				t.Fatalf("LoadV2Index() error = %v, want ErrV2EvidenceMissingReference", err)
			}
		})
	}
}

// TestLoadV2Index_evidenceMissingReferenceOnObjective proves the same
// reference resolution runs for Objective evidence, not just Task evidence.
func TestLoadV2Index_evidenceMissingReferenceOnObjective(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, v2ObjectivesDirName, "O001-first", v2ObjectiveFileName),
		"---\nid: O001\ntitle: \"First objective\"\nstatus: planned\nlast_check: C999\n---\n\n# First objective\n")

	_, err := LoadV2Index(root)
	if !errors.Is(err, ErrV2EvidenceMissingReference) {
		t.Fatalf("LoadV2Index() error = %v, want ErrV2EvidenceMissingReference", err)
	}
}

func TestLoadV2Index_issuesValid(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O001-first", "O001", "First objective")
	writeV2TaskFixture(t, root, "O001-first", "T001-alpha.md", "T001", "Alpha", "O001")
	writeV2IssueFixture(t, root, "I001-alpha.md", "I001", "open")
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{
		fileName: "I002-beta.md", id: "I002", status: "resolved",
		resolution: "{disposition: accepted, actor: {role: owner, session: owner-1}, at: '2026-09-16T00:00:00Z', reason: accepted risk}",
	})

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}

	if len(index.Issues) != 2 {
		t.Fatalf("Issues = %+v, want 2 keyed by I###", index.Issues)
	}
	if index.Issues["I001"] == nil || index.Issues["I001"].Status != IssueStatusOpen {
		t.Errorf("Issues[I001] = %+v, want the open issue", index.Issues["I001"])
	}
	if index.Issues["I002"] == nil || index.Issues["I002"].Status != IssueStatusResolved {
		t.Errorf("Issues[I002] = %+v, want the resolved issue", index.Issues["I002"])
	}
}

func TestLoadV2Index_issuesAbsentDirLoadsEmpty(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O001-first", "O001", "First objective")
	writeV2TaskFixture(t, root, "O001-first", "T001-alpha.md", "T001", "Alpha", "O001")

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v, want nil for a project with no issues/ yet", err)
	}
	if len(index.Issues) != 0 {
		t.Errorf("Issues = %+v, want empty", index.Issues)
	}
}

// writeV2DuplicateIssueFixture writes an Issue that names another Issue as
// its canonical record.
func writeV2DuplicateIssueFixture(t *testing.T, root, fileName, id, duplicateOf string) {
	t.Helper()
	content := "---\nid: " + id + "\ntitle: \"Follow-up\"\ntype: defect\nstatus: open\n" +
		"source: {kind: report, actor: {role: owner, session: owner-1}, at: '2026-09-15T00:00:00Z'}\n" +
		"duplicate_of: " + duplicateOf + "\n---\n\n# Issue\n"
	testutil.WriteFile(t, filepath.Join(root, v2IssuesDirName, fileName), content)
}

func TestLoadV2Index_issueDuplicateOfResolves(t *testing.T) {
	root := t.TempDir()
	writeV2IssueFixture(t, root, "I001-canonical.md", "I001", "open")
	writeV2DuplicateIssueFixture(t, root, "I002-duplicate.md", "I002", "I001")

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	if index.Issues["I002"].DuplicateOf != "I001" {
		t.Errorf("Issues[I002].DuplicateOf = %q, want I001", index.Issues["I002"].DuplicateOf)
	}
}

func TestLoadV2Index_issueDuplicateOfMissingTarget(t *testing.T) {
	root := t.TempDir()
	writeV2DuplicateIssueFixture(t, root, "I002-duplicate.md", "I002", "I999")

	_, err := LoadV2Index(root)
	if !errors.Is(err, ErrV2IssueMissingDuplicateTarget) {
		t.Fatalf("LoadV2Index() error = %v, want ErrV2IssueMissingDuplicateTarget", err)
	}
	if !strings.Contains(err.Error(), "I002") || !strings.Contains(err.Error(), "I999") {
		t.Errorf("LoadV2Index() error = %v, want both the issue and its missing target named", err)
	}
}

func TestLoadV2Index_issueDuplicateOfSelf(t *testing.T) {
	root := t.TempDir()
	writeV2DuplicateIssueFixture(t, root, "I001-self.md", "I001", "I001")

	_, err := LoadV2Index(root)
	if !errors.Is(err, ErrV2IssueSelfDuplicate) {
		t.Fatalf("LoadV2Index() error = %v, want ErrV2IssueSelfDuplicate", err)
	}
}

// TestLoadV2Index_issueDuplicateOfCycle proves a ring of duplicates, where no
// Issue is canonical, is refused rather than silently followed forever.
func TestLoadV2Index_issueDuplicateOfCycle(t *testing.T) {
	root := t.TempDir()
	writeV2DuplicateIssueFixture(t, root, "I001-first.md", "I001", "I002")
	writeV2DuplicateIssueFixture(t, root, "I002-second.md", "I002", "I003")
	writeV2DuplicateIssueFixture(t, root, "I003-third.md", "I003", "I001")

	_, err := LoadV2Index(root)
	if !errors.Is(err, ErrV2IssueDuplicateCycle) {
		t.Fatalf("LoadV2Index() error = %v, want ErrV2IssueDuplicateCycle", err)
	}
}

// TestLoadV2Index_issueDuplicateDiagnosticsAreDistinct proves the three
// duplicate-graph failures are separately identifiable, so a caller can tell
// a typo'd reference from a self-reference from a cycle.
func TestLoadV2Index_issueDuplicateDiagnosticsAreDistinct(t *testing.T) {
	sentinels := []error{ErrV2IssueMissingDuplicateTarget, ErrV2IssueSelfDuplicate, ErrV2IssueDuplicateCycle}
	for i, outer := range sentinels {
		for j, inner := range sentinels {
			if i != j && errors.Is(outer, inner) {
				t.Errorf("sentinel %v matches %v, want three distinct diagnostics", outer, inner)
			}
		}
	}
}

func TestLoadV2Index_issueDeterministicDiagnosticOrder(t *testing.T) {
	root := t.TempDir()
	writeV2DuplicateIssueFixture(t, root, "I001-first.md", "I001", "I997")
	writeV2DuplicateIssueFixture(t, root, "I002-second.md", "I002", "I998")
	writeV2DuplicateIssueFixture(t, root, "I003-third.md", "I003", "I999")

	_, err1 := LoadV2Index(root)
	_, err2 := LoadV2Index(root)
	if err1 == nil || err2 == nil || err1.Error() != err2.Error() {
		t.Fatalf("LoadV2Index() errors not deterministic across runs: %v vs %v", err1, err2)
	}
	if !errors.Is(err1, ErrV2IssueMissingDuplicateTarget) {
		t.Fatalf("LoadV2Index() error = %v, want ErrV2IssueMissingDuplicateTarget", err1)
	}
	if !strings.Contains(err1.Error(), "I001") {
		t.Errorf("LoadV2Index() error = %v, want the lowest-ID issue reported first", err1)
	}
}

// TestLoadV2Index_malformedIssueFailsClosed proves a malformed Issue stops
// the whole load the way E42 and E43 structural diagnostics do, rather than
// loading a project that is missing one record.
func TestLoadV2Index_malformedIssueFailsClosed(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O001-first", "O001", "First objective")
	writeV2TaskFixture(t, root, "O001-first", "T001-alpha.md", "T001", "Alpha", "O001")
	testutil.WriteFile(t, filepath.Join(root, v2IssuesDirName, "I001-broken.md"),
		"---\nid: I001\ntitle: \"Broken\"\ntype: defect\nstatus: open\n---\n\n# Issue\n")

	index, err := LoadV2Index(root)
	if !errors.Is(err, ErrV2MissingField) {
		t.Fatalf("LoadV2Index() error = %v, want ErrV2MissingField", err)
	}
	if index != nil {
		t.Errorf("LoadV2Index() = %+v, want no partially loaded index alongside the diagnostic", index)
	}
}

// v2IssueFixture describes one Issue file for the link tests, which care
// about what an Issue points at rather than the rest of its schema.
type v2IssueFixture struct {
	fileName    string
	id          string
	issueType   string
	status      string
	tasks       []string
	checks      []string
	resolution  string // optional raw resolution frontmatter block
	duplicateOf string
	history     string // optional raw history frontmatter list
}

func writeV2LinkedIssueFixture(t *testing.T, root string, fixture v2IssueFixture) {
	t.Helper()
	if fixture.issueType == "" {
		fixture.issueType = "defect"
	}
	if fixture.status == "" {
		fixture.status = "open"
	}
	content := "---\nid: " + fixture.id + "\ntitle: \"Follow-up\"\ntype: " + fixture.issueType +
		"\nstatus: " + fixture.status +
		"\nsource: {kind: report, actor: {role: owner, session: owner-1}, at: '2026-09-15T00:00:00Z'}\n"
	if len(fixture.tasks) > 0 {
		content += "tasks: [" + joinIDs(fixture.tasks) + "]\n"
	}
	if len(fixture.checks) > 0 {
		content += "checks: [" + joinIDs(fixture.checks) + "]\n"
	}
	if fixture.resolution != "" {
		content += "resolution: " + fixture.resolution + "\n"
	}
	if fixture.duplicateOf != "" {
		content += "duplicate_of: " + fixture.duplicateOf + "\n"
	}
	if fixture.history != "" {
		content += "history: " + fixture.history + "\n"
	}
	content += "---\n\n# Issue\n"
	testutil.WriteFile(t, filepath.Join(root, v2IssuesDirName, fixture.fileName), content)
}

// writeV2CheckWithIssuesFixture writes a Check that records the Issues its
// evaluation opened.
func writeV2CheckWithIssuesFixture(t *testing.T, root, fileName, id, scopeKind, scopeID string, issues []string) {
	t.Helper()
	content := "---\nid: " + id + "\nscope: {kind: " + scopeKind + ", id: " + scopeID +
		"}\nresult: CLEAR\nchecked_by: {role: checker, session: sess-1}\nchecked_at: '2026-09-14T00:00:00Z'\n" +
		"issues: [" + joinIDs(issues) + "]\n---\n\n# Check\n"
	testutil.WriteFile(t, filepath.Join(root, v2ChecksDirName, fileName), content)
}

// writeV2LinkedProject writes the shared shape the link tests vary from: one
// Objective, one Task, and one Check scoped to that Task.
func writeV2LinkedProject(t *testing.T, root string) {
	t.Helper()
	writeV2ObjectiveFixture(t, root, "O001-first", "O001", "First objective")
	writeV2TaskFixture(t, root, "O001-first", "T001-alpha.md", "T001", "Alpha", "O001")
}

// TestLoadV2Index_issueLinksResolve proves a fully paired Check-to-Issue link,
// and an Issue naming the Task carrying its repair, load into both link maps.
func TestLoadV2Index_issueLinksResolve(t *testing.T) {
	root := t.TempDir()
	writeV2LinkedProject(t, root)
	writeV2CheckWithIssuesFixture(t, root, "C001-alpha.md", "C001", "task", "T001", []string{"I001"})
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{
		fileName: "I001-alpha.md", id: "I001", tasks: []string{"T001"}, checks: []string{"C001"},
	})

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	if !reflect.DeepEqual(index.TaskIssues["T001"], []string{"I001"}) {
		t.Errorf("TaskIssues[T001] = %v, want [I001]", index.TaskIssues["T001"])
	}
	if !reflect.DeepEqual(index.CheckIssues["C001"], []string{"I001"}) {
		t.Errorf("CheckIssues[C001] = %v, want [I001]", index.CheckIssues["C001"])
	}
}

// TestLoadV2Index_checkNamesMissingIssue proves a Check's issues list stops
// being an identity-only string: a reference to an Issue that does not exist
// fails the load closed and names the Check and the reference.
func TestLoadV2Index_checkNamesMissingIssue(t *testing.T) {
	root := t.TempDir()
	writeV2LinkedProject(t, root)
	writeV2CheckWithIssuesFixture(t, root, "C001-alpha.md", "C001", "task", "T001", []string{"I999"})

	_, err := LoadV2Index(root)
	if !errors.Is(err, ErrV2IssueMissingLinkTarget) {
		t.Fatalf("LoadV2Index() error = %v, want ErrV2IssueMissingLinkTarget", err)
	}
	if !strings.Contains(err.Error(), "C001") || !strings.Contains(err.Error(), "I999") {
		t.Errorf("LoadV2Index() error = %v, want the check and the missing reference named", err)
	}
}

// TestLoadV2Index_issueNamesMissingRecord proves an Issue's own outbound
// references resolve too, in both directions it can point.
func TestLoadV2Index_issueNamesMissingRecord(t *testing.T) {
	tests := []struct {
		name    string
		fixture v2IssueFixture
		missing string
	}{
		{"missing task", v2IssueFixture{fileName: "I001-alpha.md", id: "I001", tasks: []string{"T999"}}, "T999"},
		{"missing check", v2IssueFixture{fileName: "I001-alpha.md", id: "I001", checks: []string{"C999"}}, "C999"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			writeV2LinkedProject(t, root)
			writeV2LinkedIssueFixture(t, root, tt.fixture)

			_, err := LoadV2Index(root)
			if !errors.Is(err, ErrV2IssueMissingLinkTarget) {
				t.Fatalf("LoadV2Index() error = %v, want ErrV2IssueMissingLinkTarget", err)
			}
			if !strings.Contains(err.Error(), "I001") || !strings.Contains(err.Error(), tt.missing) {
				t.Errorf("LoadV2Index() error = %v, want the issue and %s named", err, tt.missing)
			}
		})
	}
}

// TestLoadV2Index_checkIssueLinkUnpaired proves the immutable Check is the
// authoritative direction: an Issue that does not mirror a Check naming it
// fails the load closed rather than being reconciled onto either side. The two
// cases differ only in which record sorts first, so the diagnostic does not
// depend on discovery order.
func TestLoadV2Index_checkIssueLinkUnpaired(t *testing.T) {
	tests := []struct {
		name    string
		checkID string
		issueID string
	}{
		{"check named first", "C001", "I002"},
		{"issue named first", "C002", "I001"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			writeV2LinkedProject(t, root)
			writeV2CheckWithIssuesFixture(t, root, tt.checkID+"-alpha.md", tt.checkID, "task", "T001", []string{tt.issueID})
			writeV2LinkedIssueFixture(t, root, v2IssueFixture{fileName: tt.issueID + "-alpha.md", id: tt.issueID})

			index, err := LoadV2Index(root)
			if !errors.Is(err, ErrV2IssueUnpairedCheckLink) {
				t.Fatalf("LoadV2Index() error = %v, want ErrV2IssueUnpairedCheckLink", err)
			}
			if !strings.Contains(err.Error(), tt.checkID) || !strings.Contains(err.Error(), tt.issueID) {
				t.Errorf("LoadV2Index() error = %v, want both records named", err)
			}
			if index != nil {
				t.Errorf("LoadV2Index() = %+v, want no index alongside the diagnostic", index)
			}
		})
	}
}

// TestLoadV2Index_issueMayNameCheckThatDoesNotNameIt proves the asymmetry is
// legal in the other direction: a recheck or proof Check can observe an Issue
// without having recorded it, so only the Check-to-Issue direction is paired.
func TestLoadV2Index_issueMayNameCheckThatDoesNotNameIt(t *testing.T) {
	root := t.TempDir()
	writeV2LinkedProject(t, root)
	writeV2CheckFixture(t, root, "C001-recheck.md", "C001", "task", "T001")
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{
		fileName: "I001-alpha.md", id: "I001", checks: []string{"C001"},
	})

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v, want an issue naming an unrecording check to be legal", err)
	}
	if !reflect.DeepEqual(index.CheckIssues["C001"], []string{"I001"}) {
		t.Errorf("CheckIssues[C001] = %v, want [I001]", index.CheckIssues["C001"])
	}
}

// TestLoadV2Index_issueLinkMapsAreOrdered proves both maps list their Issues
// in ID order regardless of the order the referencing records are discovered
// in, and that a repeated reference is reported once.
func TestLoadV2Index_issueLinkMapsAreOrdered(t *testing.T) {
	root := t.TempDir()
	writeV2LinkedProject(t, root)
	writeV2CheckFixture(t, root, "C001-alpha.md", "C001", "task", "T001")
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{
		fileName: "I003-third.md", id: "I003", tasks: []string{"T001"}, checks: []string{"C001"},
	})
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{
		fileName: "I001-first.md", id: "I001", tasks: []string{"T001", "T001"}, checks: []string{"C001"},
	})
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{
		fileName: "I002-second.md", id: "I002", tasks: []string{"T001"}, checks: []string{"C001"},
	})

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	want := []string{"I001", "I002", "I003"}
	if !reflect.DeepEqual(index.TaskIssues["T001"], want) {
		t.Errorf("TaskIssues[T001] = %v, want %v", index.TaskIssues["T001"], want)
	}
	if !reflect.DeepEqual(index.CheckIssues["C001"], want) {
		t.Errorf("CheckIssues[C001] = %v, want %v", index.CheckIssues["C001"], want)
	}
}

// TestLoadV2Index_issueLinkMapsEmptyWithoutReferences proves a project whose
// Issues point at nothing loads with empty maps rather than an error.
func TestLoadV2Index_issueLinkMapsEmptyWithoutReferences(t *testing.T) {
	root := t.TempDir()
	writeV2LinkedProject(t, root)
	writeV2IssueFixture(t, root, "I001-alpha.md", "I001", "open")

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	if len(index.TaskIssues) != 0 || len(index.CheckIssues) != 0 {
		t.Errorf("TaskIssues = %v, CheckIssues = %v, want both empty", index.TaskIssues, index.CheckIssues)
	}
}

// TestLoadV2Index_issueLinkDiagnosticOrder proves a project with several
// dangling links reports the same one first on every run.
func TestLoadV2Index_issueLinkDiagnosticOrder(t *testing.T) {
	root := t.TempDir()
	writeV2LinkedProject(t, root)
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{fileName: "I001-first.md", id: "I001", tasks: []string{"T997"}})
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{fileName: "I002-second.md", id: "I002", tasks: []string{"T998"}})
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{fileName: "I003-third.md", id: "I003", tasks: []string{"T999"}})

	_, err1 := LoadV2Index(root)
	_, err2 := LoadV2Index(root)
	if err1 == nil || err2 == nil || err1.Error() != err2.Error() {
		t.Fatalf("LoadV2Index() errors not deterministic across runs: %v vs %v", err1, err2)
	}
	if !strings.Contains(err1.Error(), "I001") {
		t.Errorf("LoadV2Index() error = %v, want the lowest-ID issue reported first", err1)
	}
}

// TestLoadV2Index_issueResolutionRequiredWhenResolved proves a resolved Issue
// with no resolution block fails closed rather than being read as silently
// closed by status alone.
func TestLoadV2Index_issueResolutionRequiredWhenResolved(t *testing.T) {
	root := t.TempDir()
	writeV2LinkedProject(t, root)
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{fileName: "I001-alpha.md", id: "I001", status: "resolved"})

	_, err := LoadV2Index(root)
	if !errors.Is(err, ErrV2IssueResolutionRequired) {
		t.Fatalf("LoadV2Index() error = %v, want ErrV2IssueResolutionRequired", err)
	}
}

// TestLoadV2Index_issueResolutionNotAllowedWhenNotResolved proves an Issue
// returned to open or in_progress with a stale resolution still present is a
// named diagnostic, so reopening must clear it.
func TestLoadV2Index_issueResolutionNotAllowedWhenNotResolved(t *testing.T) {
	for _, status := range []string{"open", "in_progress"} {
		t.Run(status, func(t *testing.T) {
			root := t.TempDir()
			writeV2LinkedProject(t, root)
			writeV2LinkedIssueFixture(t, root, v2IssueFixture{
				fileName: "I001-alpha.md", id: "I001", status: status,
				resolution: "{disposition: accepted, actor: {role: owner, session: owner-1}, at: '2026-09-16T00:00:00Z', reason: accepted risk}",
			})

			_, err := LoadV2Index(root)
			if !errors.Is(err, ErrV2IssueResolutionNotAllowed) {
				t.Fatalf("LoadV2Index() error = %v, want ErrV2IssueResolutionNotAllowed", err)
			}
		})
	}
}

// TestLoadV2Index_issueResolutionDiagnosticsAreDistinct proves the missing-
// and not-allowed obligations are separately identifiable, as required for a
// status-versus-resolution mismatch in either direction.
func TestLoadV2Index_issueResolutionDiagnosticsAreDistinct(t *testing.T) {
	if errors.Is(ErrV2IssueResolutionRequired, ErrV2IssueResolutionNotAllowed) || errors.Is(ErrV2IssueResolutionNotAllowed, ErrV2IssueResolutionRequired) {
		t.Fatal("ErrV2IssueResolutionRequired and ErrV2IssueResolutionNotAllowed must be distinct diagnostics")
	}
}

// TestLoadV2Index_issueVerifiedResolutionRequiresProof proves verified
// closure requires a proof Check that exists, recorded CLEAR, and appears in
// the Issue's own checks list.
func TestLoadV2Index_issueVerifiedResolutionRequiresProof(t *testing.T) {
	verifiedResolution := func(check string) string {
		return "{disposition: verified, check: " + check + ", actor: {role: checker, session: sess-1}, at: '2026-09-16T00:00:00Z', reason: repaired}"
	}

	tests := []struct {
		name       string
		setup      func(t *testing.T, root string)
		resolution string
		wantErr    error
	}{
		{
			name:       "missing proof check",
			resolution: "{disposition: verified, actor: {role: checker, session: sess-1}, at: '2026-09-16T00:00:00Z', reason: repaired}",
			wantErr:    ErrV2IssueResolutionMissingProof,
		},
		{
			name:       "proof check does not exist",
			resolution: verifiedResolution("C999"),
			wantErr:    ErrV2IssueResolutionMissingProof,
		},
		{
			name: "proof check recorded needs work",
			setup: func(t *testing.T, root string) {
				testutil.WriteFile(t, filepath.Join(root, v2ChecksDirName, "C001-alpha.md"),
					"---\nid: C001\nscope: {kind: task, id: T001}\nresult: NEEDS WORK\nchecked_by: {role: checker, session: sess-1}\nchecked_at: '2026-09-14T00:00:00Z'\n---\n\n# Check\n")
			},
			resolution: verifiedResolution("C001"),
			wantErr:    ErrV2IssueResolutionUnusableProof,
		},
		{
			name: "proof check not listed on the issue",
			setup: func(t *testing.T, root string) {
				writeV2CheckFixture(t, root, "C001-alpha.md", "C001", "task", "T001")
			},
			resolution: verifiedResolution("C001"),
			wantErr:    ErrV2IssueResolutionUnusableProof,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			writeV2LinkedProject(t, root)
			if tt.setup != nil {
				tt.setup(t, root)
			}
			writeV2LinkedIssueFixture(t, root, v2IssueFixture{
				fileName: "I001-alpha.md", id: "I001", status: "resolved", resolution: tt.resolution,
			})

			_, err := LoadV2Index(root)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("LoadV2Index() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

// TestLoadV2Index_issueVerifiedResolutionValid proves a proof Check that
// exists, recorded CLEAR, and appears in the Issue's checks list closes the
// Issue cleanly.
func TestLoadV2Index_issueVerifiedResolutionValid(t *testing.T) {
	root := t.TempDir()
	writeV2LinkedProject(t, root)
	writeV2CheckFixture(t, root, "C001-alpha.md", "C001", "task", "T001")
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{
		fileName: "I001-alpha.md", id: "I001", status: "resolved", checks: []string{"C001"},
		resolution: "{disposition: verified, check: C001, actor: {role: checker, session: sess-1}, at: '2026-09-16T00:00:00Z', reason: repaired and rechecked}",
	})

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	if index.Issues["I001"].Resolution == nil || index.Issues["I001"].Resolution.Disposition != IssueDispositionVerified {
		t.Errorf("Resolution = %+v, want verified", index.Issues["I001"].Resolution)
	}
}

// TestLoadV2Index_issueAcceptedResolutionObligations proves accepted closure
// requires an owner actor and a non-empty reason and must not name a proof
// Check, because accepting risk is an owner decision, not a repair.
func TestLoadV2Index_issueAcceptedResolutionObligations(t *testing.T) {
	tests := []struct {
		name       string
		resolution string
	}{
		{
			name:       "missing owner actor",
			resolution: "{disposition: accepted, actor: {role: checker, session: sess-1}, at: '2026-09-16T00:00:00Z', reason: accepted risk}",
		},
		{
			name:       "missing reason",
			resolution: "{disposition: accepted, actor: {role: owner, session: owner-1}, at: '2026-09-16T00:00:00Z'}",
		},
		{
			name:       "names a proof check",
			resolution: "{disposition: accepted, check: C001, actor: {role: owner, session: owner-1}, at: '2026-09-16T00:00:00Z', reason: accepted risk}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			writeV2LinkedProject(t, root)
			writeV2CheckFixture(t, root, "C001-alpha.md", "C001", "task", "T001")
			writeV2LinkedIssueFixture(t, root, v2IssueFixture{
				fileName: "I001-alpha.md", id: "I001", status: "resolved", resolution: tt.resolution,
			})

			_, err := LoadV2Index(root)
			if !errors.Is(err, ErrV2IssueResolutionFieldMismatch) {
				t.Fatalf("LoadV2Index() error = %v, want ErrV2IssueResolutionFieldMismatch", err)
			}
		})
	}
}

// TestLoadV2Index_issueAcceptedResolutionValid proves an owner actor with a
// non-empty reason and no named proof Check closes the Issue as accepted.
func TestLoadV2Index_issueAcceptedResolutionValid(t *testing.T) {
	root := t.TempDir()
	writeV2LinkedProject(t, root)
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{
		fileName: "I001-alpha.md", id: "I001", status: "resolved",
		resolution: "{disposition: accepted, actor: {role: owner, session: owner-1}, at: '2026-09-16T00:00:00Z', reason: accepted risk}",
	})

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	if index.Issues["I001"].Resolution == nil || index.Issues["I001"].Resolution.Disposition != IssueDispositionAccepted {
		t.Errorf("Resolution = %+v, want accepted", index.Issues["I001"].Resolution)
	}
}

// TestLoadV2Index_issueDuplicateResolutionObligations proves duplicate
// closure requires duplicate_of naming the canonical Issue and must not name
// a proof Check, because a duplicate proves nothing itself.
func TestLoadV2Index_issueDuplicateResolutionObligations(t *testing.T) {
	tests := []struct {
		name        string
		resolution  string
		duplicateOf string
	}{
		{
			name:        "names a proof check",
			resolution:  "{disposition: duplicate, check: C001, actor: {role: owner, session: owner-1}, at: '2026-09-16T00:00:00Z', reason: same as I002}",
			duplicateOf: "I002",
		},
		{
			name:       "missing duplicate_of",
			resolution: "{disposition: duplicate, actor: {role: owner, session: owner-1}, at: '2026-09-16T00:00:00Z', reason: same as another}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			writeV2LinkedProject(t, root)
			writeV2CheckFixture(t, root, "C001-alpha.md", "C001", "task", "T001")
			writeV2IssueFixture(t, root, "I002-canonical.md", "I002", "open")
			writeV2LinkedIssueFixture(t, root, v2IssueFixture{
				fileName: "I001-alpha.md", id: "I001", status: "resolved",
				resolution: tt.resolution, duplicateOf: tt.duplicateOf,
			})

			_, err := LoadV2Index(root)
			if !errors.Is(err, ErrV2IssueResolutionFieldMismatch) {
				t.Fatalf("LoadV2Index() error = %v, want ErrV2IssueResolutionFieldMismatch", err)
			}
		})
	}
}

// TestLoadV2Index_issueDuplicateResolutionValid proves duplicate_of naming an
// existing, different Issue with no named proof Check closes the Issue as a
// duplicate.
func TestLoadV2Index_issueDuplicateResolutionValid(t *testing.T) {
	root := t.TempDir()
	writeV2LinkedProject(t, root)
	writeV2IssueFixture(t, root, "I002-canonical.md", "I002", "open")
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{
		fileName: "I001-alpha.md", id: "I001", status: "resolved", duplicateOf: "I002",
		resolution: "{disposition: duplicate, actor: {role: owner, session: owner-1}, at: '2026-09-16T00:00:00Z', reason: same as I002}",
	})

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	if index.Issues["I001"].Resolution == nil || index.Issues["I001"].Resolution.Disposition != IssueDispositionDuplicate {
		t.Errorf("Resolution = %+v, want duplicate", index.Issues["I001"].Resolution)
	}
}

// TestLoadV2Index_issueVerifiedProofSupersededStillSatisfiesAtLoad proves the
// verified obligation reads the resolution's own named proof Check directly,
// not the target's current latest Check: a proof Check that a later rerun
// has since superseded still recorded CLEAR and still appears in the Issue's
// checks, so it still satisfies verified at load. Whether that proof has
// since gone stale is a consistency question for a later inspection, not a
// load-time obligation.
func TestLoadV2Index_issueVerifiedProofSupersededStillSatisfiesAtLoad(t *testing.T) {
	root := t.TempDir()
	writeV2LinkedProject(t, root)
	writeV2CheckFixture(t, root, "C001-first.md", "C001", "task", "T001")
	testutil.WriteFile(t, filepath.Join(root, v2ChecksDirName, "C002-rerun.md"),
		"---\nid: C002\nscope: {kind: task, id: T001}\nresult: CLEAR\nchecked_by: {role: checker, session: sess-2}\nchecked_at: '2026-09-15T00:00:00Z'\nsupersedes: C001\n---\n\n# Check\n")
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{
		fileName: "I001-alpha.md", id: "I001", status: "resolved", checks: []string{"C001"},
		resolution: "{disposition: verified, check: C001, actor: {role: checker, session: sess-1}, at: '2026-09-14T00:00:00Z', reason: repaired}",
	})

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v, want the superseded-but-still-CLEAR proof to satisfy verified at load", err)
	}
	if index.LatestCheck["T001"] != "C002" {
		t.Fatalf("LatestCheck[T001] = %q, want C002 (the reference case for this test)", index.LatestCheck["T001"])
	}
}

// TestLoadV2Index_issueReopeningClearsResolution proves an Issue returned to
// open with a stale resolution still present is refused, while the same I###
// identity with the resolution cleared and a dated reopened history entry
// loads cleanly as the one recurring record — never a second Issue.
func TestLoadV2Index_issueReopeningClearsResolution(t *testing.T) {
	root := t.TempDir()
	writeV2LinkedProject(t, root)
	writeV2CheckFixture(t, root, "C001-alpha.md", "C001", "task", "T001")

	writeV2LinkedIssueFixture(t, root, v2IssueFixture{
		fileName: "I001-alpha.md", id: "I001", status: "open", checks: []string{"C001"},
		resolution: "{disposition: verified, check: C001, actor: {role: checker, session: sess-1}, at: '2026-09-15T00:00:00Z', reason: repaired}",
	})
	if _, err := LoadV2Index(root); !errors.Is(err, ErrV2IssueResolutionNotAllowed) {
		t.Fatalf("LoadV2Index() error = %v, want ErrV2IssueResolutionNotAllowed for the stale resolution", err)
	}

	writeV2LinkedIssueFixture(t, root, v2IssueFixture{
		fileName: "I001-alpha.md", id: "I001", status: "open", checks: []string{"C001"},
		history: "[{at: '2026-09-16T00:00:00Z', actor: {role: owner, session: owner-1}, kind: reopened, note: regressed}]",
	})

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v, want the reopened issue to load once resolution is cleared", err)
	}
	if index.Issues["I001"].Status != IssueStatusOpen || index.Issues["I001"].Resolution != nil {
		t.Errorf("Issues[I001] = %+v, want open with no resolution", index.Issues["I001"])
	}
	if len(index.Issues) != 1 {
		t.Fatalf("Issues = %+v, want the single I001 identity, not a second record for the same recurrence", index.Issues)
	}
}

// TestLoadV2Index_issueDeferralStaysOpen proves a deferral is a dated history
// entry on an Issue that stays open, not a fourth lifecycle status.
func TestLoadV2Index_issueDeferralStaysOpen(t *testing.T) {
	root := t.TempDir()
	writeV2LinkedProject(t, root)
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{
		fileName: "I001-alpha.md", id: "I001", status: "open",
		history: "[{at: '2026-09-16T00:00:00Z', actor: {role: owner, session: owner-1}, kind: deferred, note: revisit next release}]",
	})

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	if index.Issues["I001"].Status != IssueStatusOpen {
		t.Errorf("Status = %q, want open", index.Issues["I001"].Status)
	}
	if len(index.Issues["I001"].History) != 1 || index.Issues["I001"].History[0].Kind != IssueHistoryDeferred {
		t.Errorf("History = %+v, want one deferred entry", index.Issues["I001"].History)
	}
}

func TestLoadProjectRejectsDirectoryReadFailure(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "config.yml")
	testutil.MkdirAll(t, configPath) // a directory where config.yml should be a file

	_, err := LoadProject(root)
	if err == nil {
		t.Fatal("LoadProject() error = nil, want error reading config.yml")
	}
	if errors.Is(err, ErrMalformedSchemaVersion) || errors.Is(err, ErrUnsupportedSchemaVersion) {
		t.Errorf("LoadProject() error = %v, want a read failure, not a schema diagnostic", err)
	}
}
