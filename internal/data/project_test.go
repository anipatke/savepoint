package data

import (
	"errors"
	"path/filepath"
	"reflect"
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
