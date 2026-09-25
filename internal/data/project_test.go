package data

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/testutil"
	"gopkg.in/yaml.v3"
)

func TestCheckRuntimeSchema(t *testing.T) {
	cases := []struct {
		name          string
		configContent string
		omitConfig    bool
		wantErr       error
	}{
		{
			name:       "absent config requires migration",
			omitConfig: true,
			wantErr:    ErrSchemaMigrationRequired,
		},
		{
			name:          "absent schema_version requires migration",
			configContent: "theme:\n  bg: \"#000000\"\n",
			wantErr:       ErrSchemaMigrationRequired,
		},
		{
			name:          "explicit schema version 2 passes the runtime gate",
			configContent: "schema_version: 2\n",
		},
		{
			name:          "malformed explicit version fails with its named diagnostic",
			configContent: "schema_version: nope\n",
			wantErr:       ErrMalformedSchemaVersion,
		},
		{
			name:          "unsupported explicit version fails with its named diagnostic",
			configContent: "schema_version: 99\n",
			wantErr:       ErrUnsupportedSchemaVersion,
		},
		{
			name:          "unrelated version-shaped fields do not pass the schema gate",
			configContent: "package_version: 2\n",
			wantErr:       ErrSchemaMigrationRequired,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			savepointRoot := filepath.Join(root, ".savepoint")
			testutil.MkdirAll(t, savepointRoot)
			if !tc.omitConfig {
				testutil.WriteFile(t, filepath.Join(savepointRoot, "config.yml"), tc.configContent)
			}
			// Package/release/upgrade-manifest content must never influence
			// the schema gate, even when it carries its own version fields.
			testutil.WriteFile(t, filepath.Join(savepointRoot, ".upgrade-manifest.yml"), "schema_version: 2\nversion: 1.3.1\n")

			err := CheckRuntimeSchema(root)
			if tc.wantErr == nil {
				if err != nil {
					t.Fatalf("CheckRuntimeSchema() error = %v, want nil", err)
				}
				return
			}
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("CheckRuntimeSchema() error = %v, want wrapping %v", err, tc.wantErr)
			}
		})
	}
}

func TestLoadV2Index_valid(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O-001-first", "O-001", "First objective")
	writeV2ObjectiveFixture(t, root, "O-002-second", "O-002", "Second objective", "O-001")
	writeV2TaskFixture(t, root, "O-001-first", "T-001-alpha.md", "T-001", "Alpha", "O-001")
	writeV2TaskFixture(t, root, "O-001-first", "T-002-beta.md", "T-002", "Beta", "O-001")
	writeV2TaskFixture(t, root, "O-002-second", "T-003-gamma.md", "T-003", "Gamma", "O-002")

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}

	if len(index.Objectives) != 2 || len(index.Tasks) != 3 {
		t.Fatalf("LoadV2Index() Objectives = %d, Tasks = %d, want 2, 3", len(index.Objectives), len(index.Tasks))
	}

	owned := index.ObjectiveTasks["O-001"]
	if len(owned) != 2 || owned[0] != "T-001" || owned[1] != "T-002" {
		t.Errorf("ObjectiveTasks[O-001] = %v, want [T-001 T-002] in sorted order", owned)
	}
	if got := index.ObjectiveTasks["O-002"]; len(got) != 1 || got[0] != "T-003" {
		t.Errorf("ObjectiveTasks[O-002] = %v, want [T-003]", got)
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

func TestLoadV2Index_recordsObjectivesWithoutGoalInIDOrder(t *testing.T) {
	tests := []struct {
		name string
		ids  []string
		want []string
	}{
		{name: "one Objective", ids: []string{"O-010"}, want: []string{"O-010"}},
		{name: "several Objectives", ids: []string{"O-010", "O-002"}, want: []string{"O-002", "O-010"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			for _, id := range tc.ids {
				writeV2ObjectiveFixture(t, root, id+"-fixture", id, id+" objective")
			}
			index, err := LoadV2Index(root)
			if err != nil {
				t.Fatalf("LoadV2Index() error = %v, want non-fatal missing Goal facts", err)
			}
			if !reflect.DeepEqual(index.ObjectivesWithoutGoal, tc.want) {
				t.Fatalf("ObjectivesWithoutGoal = %v, want sorted live Objective IDs %v", index.ObjectivesWithoutGoal, tc.want)
			}
		})
	}
}

func TestOrderedObjectiveIDsForGoalUsesPriorityRankAndID(t *testing.T) {
	index := &V2Index{
		Objectives: map[string]*ObjectiveV2{
			"O-010": {ID: "O-010", Priority: ObjectivePriorityCritical, Rank: 2},
			"O-002": {ID: "O-002", Priority: ObjectivePriorityCritical, Rank: 2},
			"O-005": {ID: "O-005", Priority: ObjectivePriorityCritical, Rank: 1},
			"O-006": {ID: "O-006", Priority: ObjectivePriorityCritical},
			"O-004": {ID: "O-004", Priority: ObjectivePriorityHigh},
			"O-003": {ID: "O-003", Priority: ObjectivePriorityHigh, Rank: 1},
			"O-001": {ID: "O-001", Priority: ObjectivePriorityMedium},
			"O-009": {ID: "O-009", Priority: ObjectivePriorityMedium, Rank: 1},
			"O-007": {ID: "O-007", Priority: ObjectivePriorityMedium, Rank: 3},
			"O-008": {ID: "O-008", Priority: ObjectivePriorityMedium, Rank: 3},
			"O-011": {ID: "O-011", Priority: ObjectivePriorityLow, Rank: 1},
			"O-012": {ID: "O-012", Status: ColumnDone, Priority: ObjectivePriorityCritical, Rank: 1},
			"O-013": {ID: "O-013", Status: ColumnDone, Priority: ObjectivePriorityLow, Rank: 2},
		},
		ReleaseObjectives: map[string][]string{
			"R-001": {"O-010", "O-002", "O-005", "O-006", "O-004", "O-003", "O-001", "O-009", "O-007", "O-008", "O-011", "O-013", "O-012"},
		},
	}

	want := []string{"O-005", "O-002", "O-010", "O-006", "O-003", "O-004", "O-009", "O-007", "O-008", "O-001", "O-011", "O-012", "O-013"}
	if got := OrderedObjectiveIDsForGoal(index, "R-001"); !reflect.DeepEqual(got, want) {
		t.Fatalf("OrderedObjectiveIDsForGoal() = %v, want %v", got, want)
	}
	if got := OrderedObjectiveIDsForGoal(index, "R-999"); len(got) != 0 {
		t.Errorf("unknown Goal order = %v, want empty", got)
	}
}

func TestLoadV2IndexIgnoresDuplicateRanksForDoneObjectives(t *testing.T) {
	root := t.TempDir()
	writeV2ReleaseFixture(t, root, "R-001-current", "R-001", "Current Goal")
	writeRankedObjectiveForGoal(t, root, "O-001", "R-001", "critical", 2, "done")
	writeRankedObjectiveForGoal(t, root, "O-002", "R-001", "critical", 2, "done")

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	if len(index.DuplicateObjectiveRanks) != 0 {
		t.Fatalf("DuplicateObjectiveRanks = %+v, want done Objectives excluded", index.DuplicateObjectiveRanks)
	}
}

func TestReopenedDoneObjectiveReturnsToRankedGroupAndHealsCollision(t *testing.T) {
	root := t.TempDir()
	writeV2ReleaseFixture(t, root, "R-001-current", "R-001", "Current Goal")
	writeRankedObjectiveForGoal(t, root, "O-001", "R-001", "critical", 2, "done")
	writeRankedObjectiveForGoal(t, root, "O-002", "R-001", "critical", 2)

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() with done Objective error = %v", err)
	}
	if got, want := OrderedObjectiveIDsForGoal(index, "R-001"), []string{"O-002", "O-001"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("order with done Objective = %v, want open Objective before done %v", got, want)
	}

	path := index.Objectives["O-001"].Source.Path
	if !filepath.IsAbs(path) {
		path = filepath.Join(root, path)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(done Objective) error = %v", err)
	}
	reopened := strings.Replace(string(content), "status: done\n", "status: planned\n", 1)
	if reopened == string(content) {
		t.Fatal("done Objective status was not present to reopen")
	}
	if err := os.WriteFile(path, []byte(reopened), 0644); err != nil {
		t.Fatalf("WriteFile(reopened Objective) error = %v", err)
	}

	index, err = LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() after reopen error = %v", err)
	}
	if got, want := OrderedObjectiveIDsForGoal(index, "R-001"), []string{"O-001", "O-002"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("reopened order = %v, want the Objective back at its recorded rank %v", got, want)
	}
	if len(index.DuplicateObjectiveRanks) != 1 {
		t.Fatalf("duplicate ranks after reopen = %+v, want the recorded rank collision", index.DuplicateObjectiveRanks)
	}

	if err := WriteObjectiveGroupOrderV2(index, "R-001", ObjectivePriorityCritical, []string{"O-001", "O-002"}); err != nil {
		t.Fatalf("WriteObjectiveGroupOrderV2() to heal reopened collision error = %v", err)
	}
	index, err = LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() after healing collision error = %v", err)
	}
	if len(index.DuplicateObjectiveRanks) != 0 {
		t.Fatalf("duplicate ranks after move = %+v, want collision healed", index.DuplicateObjectiveRanks)
	}
	if got, want := OrderedObjectiveIDsForGoal(index, "R-001"), []string{"O-001", "O-002"}; !reflect.DeepEqual(got, want) {
		t.Errorf("order after healing collision = %v, want %v", got, want)
	}
}

func TestLoadV2IndexReportsDuplicateObjectiveRanksByGoalAndPriority(t *testing.T) {
	root := t.TempDir()
	writeV2ReleaseFixture(t, root, "R-001-current", "R-001", "Current Goal")
	writeV2ReleaseFixture(t, root, "R-002-other", "R-002", "Other Goal")
	writeRankedObjectiveForGoal(t, root, "O-001", "R-001", "critical", 2)
	writeRankedObjectiveForGoal(t, root, "O-002", "R-001", "critical", 2)
	writeRankedObjectiveForGoal(t, root, "O-003", "R-001", "high", 2)
	writeRankedObjectiveForGoal(t, root, "O-004", "R-002", "critical", 2)
	writeRankedObjectiveForGoal(t, root, "O-005", "R-001", "medium", 0)

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v, want duplicate ranks to remain non-fatal", err)
	}
	want := []DuplicateObjectiveRankV2{{
		GoalID:       "R-001",
		Priority:     ObjectivePriorityCritical,
		Rank:         2,
		ObjectiveIDs: []string{"O-001", "O-002"},
	}}
	if !reflect.DeepEqual(index.DuplicateObjectiveRanks, want) {
		t.Fatalf("DuplicateObjectiveRanks = %+v, want %+v", index.DuplicateObjectiveRanks, want)
	}
}

func writeRankedObjectiveForGoal(t *testing.T, root, id, goalID, priority string, rank int, statuses ...string) {
	t.Helper()
	status := "planned"
	if len(statuses) > 0 {
		status = statuses[0]
	}
	dir := filepath.Join(root, v2ObjectivesDirName, id+"-ordering")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("MkdirAll(%s) error = %v", dir, err)
	}
	content := "---\nid: " + id + "\ntitle: \"" + id + " ordering\"\nstatus: " + status + "\nrelease: " + goalID + "\npriority: " + priority + "\n"
	if rank > 0 {
		content += "rank: " + fmt.Sprint(rank) + "\n"
	}
	content += "---\n\n# " + id + "\n"
	path := filepath.Join(dir, v2ObjectiveFileName)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("WriteFile(%s) error = %v", path, err)
	}
}

func TestLoadV2Index_validGoalReferenceHasNoMissingGoalFact(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "releases", "R-001", "Release.md"),
		"---\nid: R-001\ntitle: \"Test Goal\"\nstatus: planned\n---\n\n"+
			"## Outcome\n\nDeliver the recorded promise.\n\n"+
			"## Why\n\nThe delivery boundary needs a stable identity.\n\n"+
			"## Success Conditions\n\n- All member Objectives are complete.\n\n"+
			"## Boundaries\n\nMembership is derived from Objective records.\n")
	testutil.WriteFile(t, filepath.Join(root, "objectives", "O-001-valid", "Objective.md"),
		"---\nid: O-001\ntitle: \"Valid Objective\"\nstatus: planned\nrelease: R-001\n---\n\n# Valid Objective\n")

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	if len(index.ObjectivesWithoutGoal) != 0 {
		t.Fatalf("ObjectivesWithoutGoal = %v, want no fact for an Objective with a valid Goal reference", index.ObjectivesWithoutGoal)
	}
	if got := index.ReleaseObjectives["R-001"]; len(got) != 1 || got[0] != "O-001" {
		t.Fatalf("ReleaseObjectives[R-001] = %v, want [O-001]", got)
	}
}

func TestLoadV2Index_invalidGoalReferencesRemainFatal(t *testing.T) {
	tests := []struct {
		name      string
		reference string
		want      error
	}{
		{name: "malformed", reference: "R-1", want: ErrV2InvalidReleaseReference},
		{name: "unknown", reference: "R-999", want: ErrV2MissingRelease},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			writeV2ObjectiveFixture(t, root, "O-001-invalid", "O-001", "Invalid reference")
			path := filepath.Join(root, v2ObjectivesDirName, "O-001-invalid", "Objective.md")
			raw, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("ReadFile(Objective.md) error = %v", err)
			}
			updated := strings.Replace(string(raw), "status: planned\n", "status: planned\nrelease: "+tc.reference+"\n", 1)
			if updated == string(raw) {
				t.Fatalf("fixture Objective has no planned status to extend")
			}
			if err := os.WriteFile(path, []byte(updated), 0644); err != nil {
				t.Fatalf("WriteFile(Objective.md) error = %v", err)
			}

			_, err = LoadV2Index(root)
			if !errors.Is(err, tc.want) {
				t.Fatalf("LoadV2Index() error = %v, want %v", err, tc.want)
			}
		})
	}
}

// TestLoadV2Index_movedTaskRetainsOwnership proves the assembled index keeps
// explicit Objective ownership from the Task record's own field even when
// the Task file is filed under a different Objective's tasks/ directory.
func TestLoadV2Index_movedTaskRetainsOwnership(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O-001-first", "O-001", "First objective")
	writeV2ObjectiveFixture(t, root, "O-002-second", "O-002", "Second objective")
	writeV2TaskFixture(t, root, "O-002-second", "T-001-alpha.md", "T-001", "Alpha", "O-001")

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}

	if got := index.ObjectiveTasks["O-001"]; len(got) != 1 || got[0] != "T-001" {
		t.Errorf("ObjectiveTasks[O-001] = %v, want [T-001] despite the file living under O-002-second/tasks", got)
	}
	if len(index.ObjectiveTasks["O-002"]) != 0 {
		t.Errorf("ObjectiveTasks[O-002] = %v, want empty: the containing directory does not grant ownership", index.ObjectiveTasks["O-002"])
	}
}

func TestLoadV2Index_missingOwner(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O-001-first", "O-001", "First objective")
	writeV2TaskFixture(t, root, "O-001-first", "T-001-alpha.md", "T-001", "Alpha", "O-999")

	_, err := LoadV2Index(root)
	if !errors.Is(err, ErrV2MissingOwner) {
		t.Fatalf("LoadV2Index() error = %v, want ErrV2MissingOwner", err)
	}
}

// TestLoadV2Index_integratedProjectScenario combines the behaviors E42's
// tasks proved in isolation into one project, as the completed-epic
// boundary: distinct human titles vs O-###/T-### identity, valid Objectives
// and Tasks, a Task moved to a different Objective's tasks/ directory that
// still resolves ownership from its own field, unknown frontmatter content
// surviving the full discover-then-decode path, and deterministic
// (repeat-run-stable) indexing.
func TestLoadV2Index_integratedProjectScenario(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O-001-ship", "O-001", "Ship the thing")
	writeV2ObjectiveFixture(t, root, "O-002-polish", "O-002", "Polish the thing")
	writeV2TaskFixture(t, root, "O-001-ship", "T-001-write-code.md", "T-001", "Write the code", "O-001")

	// T-002 carries unknown top-level and nested-mapping fields; the source
	// document must retain them even though DecodeTaskV2 projects only the
	// typed fields it knows about.
	testutil.WriteFile(t, filepath.Join(root, v2ObjectivesDirName, "O-001-ship", v2TasksDirName, "T-002-review-code.md"),
		"---\nid: T-002\ntitle: \"Review the code\"\nobjective: O-001\nplanned_by: {role: planner, session: planning-fixture}\nstatus: planned\nnotes: kept for reviewers\nmetadata:\n  reviewer:\n    name: sam\n---\n\n# Review the code\n")

	// T-003 is owned by O-002 but filed under O-001's tasks/ directory: a moved
	// Task whose ownership must come from its own objective field, not its
	// containing directory.
	writeV2TaskFixture(t, root, "O-001-ship", "T-003-ship-it.md", "T-003", "Ship it", "O-002")

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

	if got := index.ObjectiveTasks["O-001"]; len(got) != 2 || got[0] != "T-001" || got[1] != "T-002" {
		t.Errorf("ObjectiveTasks[O-001] = %v, want [T-001 T-002]", got)
	}
	if got := index.ObjectiveTasks["O-002"]; len(got) != 1 || got[0] != "T-003" {
		t.Errorf("ObjectiveTasks[O-002] = %v, want [T-003] despite T-003's file living under O-001-ship/tasks", got)
	}

	t2 := index.Tasks["T-002"]
	t2Mapping := t2.Source.Frontmatter.Content[0]
	if notes, ok := mappingFieldValue(t2Mapping, "notes"); !ok || notes != "kept for reviewers" {
		t.Errorf("T-002 preserved notes field = %q, ok = %v, want \"kept for reviewers\", true", notes, ok)
	}
	metadataNode := findMappingChild(t2Mapping, "metadata")
	if metadataNode == nil {
		t.Fatal("T-002 lost its unknown nested metadata field")
	}
	reviewerNode := findMappingChild(metadataNode, "reviewer")
	if reviewerNode == nil {
		t.Fatal("T-002 lost its unknown nested metadata.reviewer field")
	}
	if name, ok := mappingFieldValue(reviewerNode, "name"); !ok || name != "sam" {
		t.Errorf("T-002 preserved metadata.reviewer.name = %q, ok = %v, want \"sam\", true", name, ok)
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
	writeV2ObjectiveFixture(t, root, "O-001-first", "O-001", "First objective")
	testutil.WriteFile(t, filepath.Join(root, v2ObjectivesDirName, "O-001-first", v2TasksDirName, "T-001-alpha.md"),
		"---\nid: T-001\ntitle: \"Alpha\"\nobjective: O-001\nplanned_by: {role: planner, session: planning-fixture}\nstatus: planned\ndepends_on: [{task: T-002}]\n---\n\n# Alpha\n")
	testutil.WriteFile(t, filepath.Join(root, v2ObjectivesDirName, "O-001-first", v2TasksDirName, "T-002-beta.md"),
		"---\nid: T-002\ntitle: \"Beta\"\nobjective: O-001\nplanned_by: {role: planner, session: planning-fixture}\nstatus: planned\ndepends_on: [{task: T-001}]\n---\n\n# Beta\n")

	_, err := LoadV2Index(root)
	if !errors.Is(err, ErrV2DependencyCycle) {
		t.Fatalf("LoadV2Index() error = %v, want ErrV2DependencyCycle", err)
	}
}

func TestLoadV2Index_checksValid(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O-001-first", "O-001", "First objective")
	writeV2TaskFixture(t, root, "O-001-first", "T-001-alpha.md", "T-001", "Alpha", "O-001")
	writeV2CheckFixture(t, root, "C-001-first.md", "C-001", "task", "T-001")

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}

	if len(index.Checks) != 1 || index.Checks["C-001"] == nil {
		t.Fatalf("Checks = %+v, want C-001", index.Checks)
	}
	if got := index.ScopeChecks["T-001"]; len(got) != 1 || got[0] != "C-001" {
		t.Errorf("ScopeChecks[T-001] = %v, want [C-001]", got)
	}
	if got := index.LatestCheck["T-001"]; got != "C-001" {
		t.Errorf("LatestCheck[T-001] = %q, want C-001", got)
	}
}

func TestLoadV2Index_checksAbsentDirLoadsEmpty(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O-001-first", "O-001", "First objective")
	writeV2TaskFixture(t, root, "O-001-first", "T-001-alpha.md", "T-001", "Alpha", "O-001")

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
	writeV2CheckFixture(t, root, "C-001-first.md", "C-001", "task", "T-999")

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
	writeV2ObjectiveFixture(t, root, "O-001-first", "O-001", "First objective")
	writeV2TaskFixture(t, root, "O-001-first", "T-001-alpha.md", "T-001", "Alpha", "O-001")
	writeV2CheckFixture(t, root, "C-001-first.md", "C-001", "task", "T-001")
	testutil.WriteFile(t, filepath.Join(root, v2ChecksDirName, "C-002-rerun.md"),
		"---\nid: C-002\nscope: {kind: task, id: T-001}\nresult: CLEAR\nchecked_by: {role: checker, session: sess-2}\nexecuted_session: build-fixture\nchecked_at: '2026-09-15T00:00:00Z'\nsupersedes: C-001\n---\n\n# Check\n")

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}

	if got := index.ScopeChecks["T-001"]; len(got) != 2 || got[0] != "C-001" || got[1] != "C-002" {
		t.Errorf("ScopeChecks[T-001] = %v, want [C-001 C-002]", got)
	}
	if got := index.LatestCheck["T-001"]; got != "C-002" {
		t.Errorf("LatestCheck[T-001] = %q, want C-002", got)
	}
}

func TestLoadV2Index_checkOrderingAcrossC999C1000Boundary(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O-001-first", "O-001", "First objective")
	writeV2TaskFixture(t, root, "O-001-first", "T-001-alpha.md", "T-001", "Alpha", "O-001")
	writeV2TaskFixture(t, root, "O-001-first", "T-002-beta.md", "T-002", "Beta", "O-001")
	writeV2CheckFixture(t, root, "C-999-older.md", "C-999", "task", "T-001")
	testutil.WriteFile(t, filepath.Join(root, v2ChecksDirName, "C-1000-newer.md"),
		"---\nid: C-1000\nscope: {kind: task, id: T-001}\nresult: CLEAR\nchecked_by: {role: checker, session: sess-2}\nexecuted_session: build-fixture\nchecked_at: '2026-09-15T00:00:00Z'\nsupersedes: C-999\n---\n\n# Check\n")

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	if got := index.ScopeChecks["T-001"]; len(got) != 2 || got[0] != "C-999" || got[1] != "C-1000" {
		t.Fatalf("ScopeChecks[T-001] = %v, want [C-999 C-1000] in numeric order", got)
	}
	if got := index.LatestCheck["T-001"]; got != "C-1000" {
		t.Fatalf("LatestCheck[T-001] = %q, want C-1000", got)
	}

	t001 := index.Tasks["T-001"]
	t001.Evidence = &Evidence{Freshness: &Freshness{
		State: FreshnessCurrent, Check: "C-1000",
		AssessedBy: Actor{Role: ActorRoleChecker, Session: "checker-1"},
		Basis:      "reviewed the newer check",
	}}
	t001.Status = ColumnDone
	if got := ResolveClearance(index, "T-001"); got.State != ClearanceCurrent || got.Check != "C-1000" {
		t.Fatalf("ResolveClearance(T-001) = %+v, want current/C-1000", got)
	}
	index.Tasks["T-002"].Status = ColumnPlanned
	index.Tasks["T-002"].DependsOn = []TaskDependencyV2{{Task: "T-001", Requires: TaskDependencyClear}}
	if dependency := ResolveTaskDependencyV2(index, index.Tasks["T-002"].DependsOn[0]); !dependency.Satisfied {
		t.Fatalf("ResolveTaskDependencyV2(T-002 -> T-001) = %+v, want satisfied", dependency)
	}

	t001.Status = ColumnInProgress
	t001.Stage = StageAudit
	if decision := ResolveTaskCompletion(index, "T-001"); !decision.Allowed {
		t.Fatalf("ResolveTaskCompletion(T-001) = %+v, want allowed from latest C-1000", decision)
	}
	index.Releases = map[string]*ReleaseV2{"R-001": {ID: "R-001", Title: "Test Goal"}}
	index.Objectives["O-001"].Release = "R-001"
	index.ReleaseObjectives["R-001"] = []string{"O-001"}
	index.ObjectivesWithoutGoal = nil
	next := ResolveNext(NextInput{
		Index:  index,
		Router: &RouterStateV2{State: RouterPhaseTask, Release: "R-001", Objective: "O-001", Task: "T-001"},
	})
	if next.Kind != NextExecute || next.Task == nil || next.Task.ID != "T-001" {
		t.Fatalf("ResolveNext() = %+v, want execute T-001 from latest C-1000", next)
	}
}

func TestLoadV2Index_checkSupersedesMultipleHeads(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O-001-first", "O-001", "First objective")
	writeV2TaskFixture(t, root, "O-001-first", "T-001-alpha.md", "T-001", "Alpha", "O-001")
	writeV2CheckFixture(t, root, "C-001-first.md", "C-001", "task", "T-001")
	testutil.WriteFile(t, filepath.Join(root, v2ChecksDirName, "C-002-rerun.md"),
		"---\nid: C-002\nscope: {kind: task, id: T-001}\nresult: CLEAR\nchecked_by: {role: checker, session: sess-2}\nexecuted_session: build-fixture\nchecked_at: '2026-09-15T00:00:00Z'\nsupersedes: C-001\n---\n\n# Check\n")
	writeV2CheckFixture(t, root, "C-003-stray.md", "C-003", "task", "T-001")

	_, err := LoadV2Index(root)
	if !errors.Is(err, ErrV2CheckSupersedesConflict) {
		t.Fatalf("LoadV2Index() error = %v, want ErrV2CheckSupersedesConflict", err)
	}
	if !strings.Contains(err.Error(), "C-003") || !strings.Contains(err.Error(), "C-002") {
		t.Fatalf("LoadV2Index() error = %v, want the stray head and required predecessor named", err)
	}
}

func TestLoadV2Index_checkSupersedesOrderMismatch(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O-001-first", "O-001", "First objective")
	writeV2TaskFixture(t, root, "O-001-first", "T-001-alpha.md", "T-001", "Alpha", "O-001")
	writeV2CheckFixture(t, root, "C-001-first.md", "C-001", "task", "T-001")
	testutil.WriteFile(t, filepath.Join(root, v2ChecksDirName, "C-002-rerun.md"),
		"---\nid: C-002\nscope: {kind: task, id: T-001}\nresult: CLEAR\nchecked_by: {role: checker, session: sess-2}\nexecuted_session: build-fixture\nchecked_at: '2026-09-15T00:00:00Z'\nsupersedes: C-003\n---\n\n# Check\n")
	testutil.WriteFile(t, filepath.Join(root, v2ChecksDirName, "C-003-rerun.md"),
		"---\nid: C-003\nscope: {kind: task, id: T-001}\nresult: CLEAR\nchecked_by: {role: checker, session: sess-3}\nexecuted_session: build-fixture\nchecked_at: '2026-09-16T00:00:00Z'\nsupersedes: C-001\n---\n\n# Check\n")

	_, err := LoadV2Index(root)
	if !errors.Is(err, ErrV2CheckSupersedesConflict) {
		t.Fatalf("LoadV2Index() error = %v, want ErrV2CheckSupersedesConflict", err)
	}
}

func TestLoadV2Index_checkSupersedesMissingTarget(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O-001-first", "O-001", "First objective")
	writeV2TaskFixture(t, root, "O-001-first", "T-001-alpha.md", "T-001", "Alpha", "O-001")
	testutil.WriteFile(t, filepath.Join(root, v2ChecksDirName, "C-002-rerun.md"),
		"---\nid: C-002\nscope: {kind: task, id: T-001}\nresult: CLEAR\nchecked_by: {role: checker, session: sess-2}\nexecuted_session: build-fixture\nchecked_at: '2026-09-15T00:00:00Z'\nsupersedes: C-999\n---\n\n# Check\n")

	_, err := LoadV2Index(root)
	if !errors.Is(err, ErrV2CheckMissingReference) {
		t.Fatalf("LoadV2Index() error = %v, want ErrV2CheckMissingReference", err)
	}
}

func TestLoadV2Index_checkSupersedesScopeMismatch(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O-001-first", "O-001", "First objective")
	writeV2TaskFixture(t, root, "O-001-first", "T-001-alpha.md", "T-001", "Alpha", "O-001")
	writeV2TaskFixture(t, root, "O-001-first", "T-002-beta.md", "T-002", "Beta", "O-001")
	writeV2CheckFixture(t, root, "C-001-first.md", "C-001", "task", "T-001")
	testutil.WriteFile(t, filepath.Join(root, v2ChecksDirName, "C-002-rerun.md"),
		"---\nid: C-002\nscope: {kind: task, id: T-002}\nresult: CLEAR\nchecked_by: {role: checker, session: sess-2}\nexecuted_session: build-fixture\nchecked_at: '2026-09-15T00:00:00Z'\nsupersedes: C-001\n---\n\n# Check\n")

	_, err := LoadV2Index(root)
	if !errors.Is(err, ErrV2CheckSupersedesConflict) {
		t.Fatalf("LoadV2Index() error = %v, want ErrV2CheckSupersedesConflict", err)
	}
}

// TestLoadV2Index_checkSupersedesFork proves two Checks superseding the same
// target Check is rejected, not silently treated as two valid chain heads.
func TestLoadV2Index_checkSupersedesFork(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O-001-first", "O-001", "First objective")
	writeV2TaskFixture(t, root, "O-001-first", "T-001-alpha.md", "T-001", "Alpha", "O-001")
	writeV2CheckFixture(t, root, "C-001-first.md", "C-001", "task", "T-001")
	testutil.WriteFile(t, filepath.Join(root, v2ChecksDirName, "C-002-rerun.md"),
		"---\nid: C-002\nscope: {kind: task, id: T-001}\nresult: CLEAR\nchecked_by: {role: checker, session: sess-2}\nexecuted_session: build-fixture\nchecked_at: '2026-09-15T00:00:00Z'\nsupersedes: C-001\n---\n\n# Check\n")
	testutil.WriteFile(t, filepath.Join(root, v2ChecksDirName, "C-003-also-rerun.md"),
		"---\nid: C-003\nscope: {kind: task, id: T-001}\nresult: NEEDS WORK\nchecked_by: {role: checker, session: sess-3}\nexecuted_session: build-fixture\nchecked_at: '2026-09-15T01:00:00Z'\nsupersedes: C-001\n---\n\n# Check\n")

	_, err := LoadV2Index(root)
	if !errors.Is(err, ErrV2CheckSupersedesConflict) {
		t.Fatalf("LoadV2Index() error = %v, want ErrV2CheckSupersedesConflict", err)
	}
}

func TestLoadV2Index_checkSupersedesCycle(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O-001-first", "O-001", "First objective")
	writeV2TaskFixture(t, root, "O-001-first", "T-001-alpha.md", "T-001", "Alpha", "O-001")
	testutil.WriteFile(t, filepath.Join(root, v2ChecksDirName, "C-001-first.md"),
		"---\nid: C-001\nscope: {kind: task, id: T-001}\nresult: CLEAR\nchecked_by: {role: checker, session: sess-1}\nexecuted_session: build-fixture\nchecked_at: '2026-09-14T00:00:00Z'\nsupersedes: C-002\n---\n\n# Check\n")
	testutil.WriteFile(t, filepath.Join(root, v2ChecksDirName, "C-002-second.md"),
		"---\nid: C-002\nscope: {kind: task, id: T-001}\nresult: CLEAR\nchecked_by: {role: checker, session: sess-2}\nexecuted_session: build-fixture\nchecked_at: '2026-09-15T00:00:00Z'\nsupersedes: C-001\n---\n\n# Check\n")

	_, err := LoadV2Index(root)
	if !errors.Is(err, ErrV2CheckSupersedesConflict) {
		t.Fatalf("LoadV2Index() error = %v, want ErrV2CheckSupersedesConflict", err)
	}
}

func TestLoadV2Index_checkDeterministicDiagnosticOrder(t *testing.T) {
	root := t.TempDir()
	writeV2CheckFixture(t, root, "C-001-first.md", "C-001", "task", "T-998")
	writeV2CheckFixture(t, root, "C-002-second.md", "C-002", "task", "T-999")

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
	writeV2ObjectiveFixture(t, root, "O-001-first", "O-001", "First objective")
	writeV2CheckFixture(t, root, "C-001-first.md", "C-001", "task", "T-001")
	testutil.WriteFile(t, filepath.Join(root, v2ObjectivesDirName, "O-001-first", v2TasksDirName, "T-001-alpha.md"),
		"---\nid: T-001\ntitle: \"Alpha\"\nobjective: O-001\nplanned_by: {role: planner, session: planning-fixture}\nstatus: planned\n"+
			"last_check: C-001\n"+
			"freshness: {state: current, check: C-001, assessed_by: {role: checker, session: sess-1}, assessed_at: '2026-09-15T00:00:00Z', basis: reviewed}\n"+
			"owner_validation: {required: true, accepted_check: C-001, accepted_by: {role: owner, session: owner-1}}\n"+
			"exception: {requirements: [TEST-08], reason: waived, owner: ani, recorded_at: '2026-09-15T00:00:00Z', check: C-001}\n"+
			"---\n\n# Alpha\n")

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}

	task := index.Tasks["T-001"]
	if task.Evidence == nil {
		t.Fatal("Tasks[T-001].Evidence = nil, want decoded evidence")
	}
	if task.Evidence.LastCheck != "C-001" {
		t.Errorf("Evidence.LastCheck = %q, want C-001", task.Evidence.LastCheck)
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
		{"last_check", "last_check: C-999\n"},
		{"freshness.check", "freshness: {state: current, check: C-999, assessed_by: {role: checker, session: sess-1}, assessed_at: '2026-09-15T00:00:00Z', basis: reviewed}\n"},
		{"owner_validation.accepted_check", "owner_validation: {required: true, accepted_check: C-999, accepted_by: {role: owner, session: owner-1}}\n"},
		{"exception.check", "exception: {requirements: [TEST-08], reason: waived, owner: ani, recorded_at: '2026-09-15T00:00:00Z', check: C-999}\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			writeV2ObjectiveFixture(t, root, "O-001-first", "O-001", "First objective")
			testutil.WriteFile(t, filepath.Join(root, v2ObjectivesDirName, "O-001-first", v2TasksDirName, "T-001-alpha.md"),
				"---\nid: T-001\ntitle: \"Alpha\"\nobjective: O-001\nplanned_by: {role: planner, session: planning-fixture}\nstatus: planned\n"+tt.field+"---\n\n# Alpha\n")

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
	testutil.WriteFile(t, filepath.Join(root, v2ObjectivesDirName, "O-001-first", v2ObjectiveFileName),
		"---\nid: O-001\ntitle: \"First objective\"\nstatus: planned\nlast_check: C-999\n---\n\n# First objective\n")

	_, err := LoadV2Index(root)
	if !errors.Is(err, ErrV2EvidenceMissingReference) {
		t.Fatalf("LoadV2Index() error = %v, want ErrV2EvidenceMissingReference", err)
	}
}

func TestLoadV2Index_issuesValid(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O-001-first", "O-001", "First objective")
	writeV2TaskFixture(t, root, "O-001-first", "T-001-alpha.md", "T-001", "Alpha", "O-001")
	writeV2IssueFixture(t, root, "I-001-alpha.md", "I-001", "open")
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{
		fileName: "I-002-beta.md", id: "I-002", status: "resolved",
		resolution: "{disposition: accepted, actor: {role: owner, session: owner-1}, at: '2026-09-16T00:00:00Z', reason: accepted risk}",
	})

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}

	if len(index.Issues) != 2 {
		t.Fatalf("Issues = %+v, want 2 keyed by I-###", index.Issues)
	}
	if index.Issues["I-001"] == nil || index.Issues["I-001"].Status != IssueStatusOpen {
		t.Errorf("Issues[I-001] = %+v, want the open issue", index.Issues["I-001"])
	}
	if index.Issues["I-002"] == nil || index.Issues["I-002"].Status != IssueStatusResolved {
		t.Errorf("Issues[I-002] = %+v, want the resolved issue", index.Issues["I-002"])
	}
}

func TestLoadV2Index_issuesAbsentDirLoadsEmpty(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O-001-first", "O-001", "First objective")
	writeV2TaskFixture(t, root, "O-001-first", "T-001-alpha.md", "T-001", "Alpha", "O-001")

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
	writeV2IssueFixture(t, root, "I-001-canonical.md", "I-001", "open")
	writeV2DuplicateIssueFixture(t, root, "I-002-duplicate.md", "I-002", "I-001")

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	if index.Issues["I-002"].DuplicateOf != "I-001" {
		t.Errorf("Issues[I-002].DuplicateOf = %q, want I-001", index.Issues["I-002"].DuplicateOf)
	}
}

func TestLoadV2Index_issueDuplicateOfMissingTarget(t *testing.T) {
	root := t.TempDir()
	writeV2DuplicateIssueFixture(t, root, "I-002-duplicate.md", "I-002", "I-999")

	_, err := LoadV2Index(root)
	if !errors.Is(err, ErrV2IssueMissingDuplicateTarget) {
		t.Fatalf("LoadV2Index() error = %v, want ErrV2IssueMissingDuplicateTarget", err)
	}
	if !strings.Contains(err.Error(), "I-002") || !strings.Contains(err.Error(), "I-999") {
		t.Errorf("LoadV2Index() error = %v, want both the issue and its missing target named", err)
	}
}

func TestLoadV2Index_issueDuplicateOfSelf(t *testing.T) {
	root := t.TempDir()
	writeV2DuplicateIssueFixture(t, root, "I-001-self.md", "I-001", "I-001")

	_, err := LoadV2Index(root)
	if !errors.Is(err, ErrV2IssueSelfDuplicate) {
		t.Fatalf("LoadV2Index() error = %v, want ErrV2IssueSelfDuplicate", err)
	}
}

// TestLoadV2Index_issueDuplicateOfCycle proves a ring of duplicates, where no
// Issue is canonical, is refused rather than silently followed forever.
func TestLoadV2Index_issueDuplicateOfCycle(t *testing.T) {
	root := t.TempDir()
	writeV2DuplicateIssueFixture(t, root, "I-001-first.md", "I-001", "I-002")
	writeV2DuplicateIssueFixture(t, root, "I-002-second.md", "I-002", "I-003")
	writeV2DuplicateIssueFixture(t, root, "I-003-third.md", "I-003", "I-001")

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
	writeV2DuplicateIssueFixture(t, root, "I-001-first.md", "I-001", "I-997")
	writeV2DuplicateIssueFixture(t, root, "I-002-second.md", "I-002", "I-998")
	writeV2DuplicateIssueFixture(t, root, "I-003-third.md", "I-003", "I-999")

	_, err1 := LoadV2Index(root)
	_, err2 := LoadV2Index(root)
	if err1 == nil || err2 == nil || err1.Error() != err2.Error() {
		t.Fatalf("LoadV2Index() errors not deterministic across runs: %v vs %v", err1, err2)
	}
	if !errors.Is(err1, ErrV2IssueMissingDuplicateTarget) {
		t.Fatalf("LoadV2Index() error = %v, want ErrV2IssueMissingDuplicateTarget", err1)
	}
	if !strings.Contains(err1.Error(), "I-001") {
		t.Errorf("LoadV2Index() error = %v, want the lowest-ID issue reported first", err1)
	}
}

// TestLoadV2Index_malformedIssueFailsClosed proves a malformed Issue stops
// the whole load the way E42 and E43 structural diagnostics do, rather than
// loading a project that is missing one record.
func TestLoadV2Index_malformedIssueFailsClosed(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O-001-first", "O-001", "First objective")
	writeV2TaskFixture(t, root, "O-001-first", "T-001-alpha.md", "T-001", "Alpha", "O-001")
	testutil.WriteFile(t, filepath.Join(root, v2IssuesDirName, "I-001-broken.md"),
		"---\nid: I-001\ntitle: \"Broken\"\ntype: defect\nstatus: open\n---\n\n# Issue\n")

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
	escalatedTo string
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
	if fixture.escalatedTo != "" {
		content += "escalated_to: " + fixture.escalatedTo + "\n"
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
		"}\nresult: CLEAR\nchecked_by: {role: checker, session: sess-1}\nexecuted_session: build-fixture\nchecked_at: '2026-09-14T00:00:00Z'\n" +
		"issues: [" + joinIDs(issues) + "]\n---\n\n# Check\n"
	testutil.WriteFile(t, filepath.Join(root, v2ChecksDirName, fileName), content)
}

// writeV2LinkedProject writes the shared shape the link tests vary from: one
// Objective, one Task, and one Check scoped to that Task.
func writeV2LinkedProject(t *testing.T, root string) {
	t.Helper()
	writeV2ObjectiveFixture(t, root, "O-001-first", "O-001", "First objective")
	writeV2TaskFixture(t, root, "O-001-first", "T-001-alpha.md", "T-001", "Alpha", "O-001")
}

func TestCreateTaskV2InjectsIdentityAndPreservesDraftBytes(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O-001-first", "O-001", "First objective")
	draft := "---\r\ntitle: Review the Code\r\nplanned_by: {role: planner, session: draft-source}\r\nstatus: planned\r\ndepends_on: []\r\n---\r\n\r\n# Review\r\n\r\nKeep this body byte for byte.\r\n"

	created, err := CreateTaskV2(root, "O-001", draft)
	if err != nil {
		t.Fatalf("CreateTaskV2() error = %v", err)
	}
	if created.ID != "T-001" {
		t.Errorf("created ID = %q, want T-001", created.ID)
	}
	wantPath := filepath.Join("objectives", "O-001-first", "tasks", "T-001-review-the-code.md")
	if created.Path != wantPath {
		t.Errorf("created path = %q, want %q", created.Path, wantPath)
	}
	content, err := os.ReadFile(filepath.Join(root, created.Path))
	if err != nil {
		t.Fatal(err)
	}
	wantPrefix := "---\r\nid: T-001\r\nobjective: O-001\r\n"
	if !strings.HasPrefix(string(content), wantPrefix) || !strings.HasSuffix(string(content), draft[len("---\r\n"):]) {
		t.Fatalf("created content did not preserve draft bytes around injected identity:\n%s", content)
	}
	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() after creation: %v", err)
	}
	task := index.Tasks[created.ID]
	if task == nil || task.Objective != "O-001" || task.Title != "Review the Code" || task.PlannedBy.Session != "draft-source" {
		t.Fatalf("indexed Task = %+v, want decoded fields preserved with allocated owner", task)
	}
}

func TestCreateTaskV2RejectsInvalidDraftAndMissingObjectiveWithoutReservation(t *testing.T) {
	validDraft := "---\ntitle: Review\nplanned_by: {role: planner, session: draft-source}\nstatus: planned\n---\n\n# Review\n"
	tests := []struct {
		name      string
		objective string
		draft     string
	}{
		{name: "authored ID", objective: "O-001", draft: strings.Replace(validDraft, "title: Review", "id: T-777\ntitle: Review", 1)},
		{name: "mismatched owner", objective: "O-001", draft: strings.Replace(validDraft, "title: Review", "objective: O-002\ntitle: Review", 1)},
		{name: "malformed YAML", objective: "O-001", draft: "---\ntitle: [broken\n---\nbody\n"},
		{name: "missing Objective", objective: "O-999", draft: validDraft},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			writeV2ObjectiveFixture(t, root, "O-001-first", "O-001", "First objective")
			if _, err := CreateTaskV2(root, tc.objective, tc.draft); err == nil {
				t.Fatal("CreateTaskV2() error = nil, want rejection")
			}
			for _, name := range []string{taskIDHighWaterFile, taskIDLockFile} {
				if _, err := os.Lstat(filepath.Join(root, name)); !errors.Is(err, os.ErrNotExist) {
					t.Errorf("%s stat error = %v, want no allocation state", name, err)
				}
			}
			index, err := LoadV2Index(root)
			if err != nil || len(index.Tasks) != 0 {
				t.Errorf("project after rejected draft: index=%+v error=%v, want valid and task-free", index, err)
			}
		})
	}
}

func TestCreateTaskV2RefusesOccupiedDestinationAndRetiresID(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O-001-first", "O-001", "First objective")
	tasksDir := filepath.Join(root, "objectives", "O-001-first", "tasks")
	occupiedPath := filepath.Join(tasksDir, "T-001-collision-task.md")
	markerPath := filepath.Join(occupiedPath, "keep.txt")
	testutil.WriteFile(t, markerPath, "preserve occupied destination\n")
	draft := "---\ntitle: Collision Task\nplanned_by: {role: planner, session: draft-source}\nstatus: planned\n---\n\n# Collision\n"

	if _, err := CreateTaskV2(root, "O-001", draft); err == nil {
		t.Fatal("CreateTaskV2() error = nil, want occupied destination rejection")
	}
	if got, err := os.ReadFile(markerPath); err != nil || string(got) != "preserve occupied destination\n" {
		t.Errorf("occupied marker = %q, error = %v; existing content must remain unchanged", got, err)
	}
	if got, err := os.ReadFile(filepath.Join(root, taskIDHighWaterFile)); err != nil || string(got) != "last_issued: 1\n" {
		t.Errorf("high-water mark = %q, error = %v; failed ID must stay retired", got, err)
	}
	if _, err := LoadV2Index(root); err != nil {
		t.Errorf("LoadV2Index() after occupied path refusal: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tasksDir, "T-001-collision-task.md")); err != nil {
		t.Errorf("occupied destination changed after refusal: %v", err)
	}
}

func TestCreateExclusiveTaskFileDoesNotOverwriteExistingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "T-001-collision.md")
	testutil.WriteFile(t, path, "keep existing Task bytes\n")
	if _, err := createExclusiveTaskFile(path, []byte("new Task bytes\n")); err == nil {
		t.Fatal("createExclusiveTaskFile() error = nil, want occupied-file refusal")
	}
	if got, err := os.ReadFile(path); err != nil || string(got) != "keep existing Task bytes\n" {
		t.Errorf("existing file = %q, error = %v; content must remain unchanged", got, err)
	}
}

func TestCreateTaskV2InvalidProjectDoesNotChangeExistingContent(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O-001-first", "O-001", "First objective")
	invalidPath := filepath.Join(root, "objectives", "O-001-first", "tasks", "broken.md")
	invalidContent := "not a V2 Task record\n"
	testutil.WriteFile(t, invalidPath, invalidContent)
	draft := "---\ntitle: Review\nplanned_by: {role: planner, session: draft-source}\nstatus: planned\n---\n\n# Review\n"

	if _, err := CreateTaskV2(root, "O-001", draft); err == nil {
		t.Fatal("CreateTaskV2() error = nil, want invalid project refusal")
	}
	if got, err := os.ReadFile(invalidPath); err != nil || string(got) != invalidContent {
		t.Errorf("existing invalid content = %q, error = %v; it must remain unchanged", got, err)
	}
	for _, name := range []string{taskIDHighWaterFile, taskIDLockFile} {
		if _, err := os.Lstat(filepath.Join(root, name)); !errors.Is(err, os.ErrNotExist) {
			t.Errorf("%s stat error = %v, want no creation state", name, err)
		}
	}
}

func TestCreateTaskV2StrictValidationFailureRemovesTaskButRetiresID(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O-001-first", "O-001", "First objective")
	draft := "---\ntitle: Broken Link\nplanned_by: {role: planner, session: draft-source}\nstatus: planned\ndepends_on: [{task: T-999}]\n---\n\n# Broken Link\n"

	if _, err := CreateTaskV2(root, "O-001", draft); err == nil {
		t.Fatal("CreateTaskV2() error = nil, want strict index validation failure")
	}
	if _, err := os.Stat(filepath.Join(root, "objectives", "O-001-first", "tasks")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("new tasks directory stat error = %v, want rollback to remove it", err)
	}
	if got, err := os.ReadFile(filepath.Join(root, taskIDHighWaterFile)); err != nil || string(got) != "last_issued: 1\n" {
		t.Errorf("high-water mark = %q, error = %v; failed ID must stay retired", got, err)
	}
	if _, err := os.Lstat(filepath.Join(root, taskIDLockFile)); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("lock stat error = %v, want lock released", err)
	}
	index, err := LoadV2Index(root)
	if err != nil || len(index.Tasks) != 0 {
		t.Errorf("project after rollback: index=%+v error=%v, want valid and task-free", index, err)
	}
}

func TestCreateTaskV2ConcurrentAcrossObjectives(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O-001-first", "O-001", "First objective")
	writeV2ObjectiveFixture(t, root, "O-002-second", "O-002", "Second objective")
	draft := "---\ntitle: Concurrent Task\nplanned_by: {role: planner, session: concurrent-test}\nstatus: planned\n---\n\n# Concurrent\n"
	type outcome struct {
		objective string
		created   TaskCreationResult
		err       error
	}
	results := make(chan outcome, 2)
	for _, objective := range []string{"O-001", "O-002"} {
		objective := objective
		go func() {
			created, err := CreateTaskV2(root, objective, draft)
			results <- outcome{objective: objective, created: created, err: err}
		}()
	}

	ids := make(map[string]string, 2)
	for range 2 {
		result := <-results
		if result.err != nil {
			t.Fatalf("CreateTaskV2(%s) error = %v", result.objective, result.err)
		}
		if previous, exists := ids[result.created.ID]; exists {
			t.Fatalf("Objectives %s and %s received duplicate ID %s", previous, result.objective, result.created.ID)
		}
		ids[result.created.ID] = result.objective
	}
	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() after concurrent creation: %v", err)
	}
	if len(ids) != 2 || len(index.Tasks) != 2 {
		t.Fatalf("created IDs=%v, indexed tasks=%v; want two distinct valid Tasks", ids, index.Tasks)
	}
	for id, objective := range ids {
		if task := index.Tasks[id]; task == nil || task.Objective != objective {
			t.Errorf("indexed Task %s = %+v, want owner %s", id, task, objective)
		}
	}
}

// TestLoadV2Index_issueLinksResolve proves a fully paired Check-to-Issue link,
// and an Issue naming the Task carrying its repair, load into both link maps.
func TestLoadV2Index_issueLinksResolve(t *testing.T) {
	root := t.TempDir()
	writeV2LinkedProject(t, root)
	writeV2CheckWithIssuesFixture(t, root, "C-001-alpha.md", "C-001", "task", "T-001", []string{"I-001"})
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{
		fileName: "I-001-alpha.md", id: "I-001", tasks: []string{"T-001"}, checks: []string{"C-001"},
	})

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	if !reflect.DeepEqual(index.TaskIssues["T-001"], []string{"I-001"}) {
		t.Errorf("TaskIssues[T-001] = %v, want [I-001]", index.TaskIssues["T-001"])
	}
	if !reflect.DeepEqual(index.CheckIssues["C-001"], []string{"I-001"}) {
		t.Errorf("CheckIssues[C-001] = %v, want [I-001]", index.CheckIssues["C-001"])
	}
}

// TestLoadV2Index_checkNamesMissingIssue proves a Check's issues list stops
// being an identity-only string: a reference to an Issue that does not exist
// fails the load closed and names the Check and the reference.
func TestLoadV2Index_checkNamesMissingIssue(t *testing.T) {
	root := t.TempDir()
	writeV2LinkedProject(t, root)
	writeV2CheckWithIssuesFixture(t, root, "C-001-alpha.md", "C-001", "task", "T-001", []string{"I-999"})

	_, err := LoadV2Index(root)
	if !errors.Is(err, ErrV2IssueMissingLinkTarget) {
		t.Fatalf("LoadV2Index() error = %v, want ErrV2IssueMissingLinkTarget", err)
	}
	if !strings.Contains(err.Error(), "C-001") || !strings.Contains(err.Error(), "I-999") {
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
		{"missing task", v2IssueFixture{fileName: "I-001-alpha.md", id: "I-001", tasks: []string{"T-999"}}, "T-999"},
		{"missing check", v2IssueFixture{fileName: "I-001-alpha.md", id: "I-001", checks: []string{"C-999"}}, "C-999"},
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
			if !strings.Contains(err.Error(), "I-001") || !strings.Contains(err.Error(), tt.missing) {
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
		{"check named first", "C-001", "I-002"},
		{"issue named first", "C-002", "I-001"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			writeV2LinkedProject(t, root)
			writeV2CheckWithIssuesFixture(t, root, tt.checkID+"-alpha.md", tt.checkID, "task", "T-001", []string{tt.issueID})
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
	writeV2CheckFixture(t, root, "C-001-recheck.md", "C-001", "task", "T-001")
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{
		fileName: "I-001-alpha.md", id: "I-001", checks: []string{"C-001"},
	})

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v, want an issue naming an unrecording check to be legal", err)
	}
	if !reflect.DeepEqual(index.CheckIssues["C-001"], []string{"I-001"}) {
		t.Errorf("CheckIssues[C-001] = %v, want [I-001]", index.CheckIssues["C-001"])
	}
}

// TestLoadV2Index_issueLinkMapsAreOrdered proves both maps list their Issues
// in ID order regardless of the order the referencing records are discovered
// in, and that a repeated reference is reported once.
func TestLoadV2Index_issueLinkMapsAreOrdered(t *testing.T) {
	root := t.TempDir()
	writeV2LinkedProject(t, root)
	writeV2CheckFixture(t, root, "C-001-alpha.md", "C-001", "task", "T-001")
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{
		fileName: "I-003-third.md", id: "I-003", tasks: []string{"T-001"}, checks: []string{"C-001"},
	})
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{
		fileName: "I-001-first.md", id: "I-001", tasks: []string{"T-001", "T-001"}, checks: []string{"C-001"},
	})
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{
		fileName: "I-002-second.md", id: "I-002", tasks: []string{"T-001"}, checks: []string{"C-001"},
	})

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	want := []string{"I-001", "I-002", "I-003"}
	if !reflect.DeepEqual(index.TaskIssues["T-001"], want) {
		t.Errorf("TaskIssues[T-001] = %v, want %v", index.TaskIssues["T-001"], want)
	}
	if !reflect.DeepEqual(index.CheckIssues["C-001"], want) {
		t.Errorf("CheckIssues[C-001] = %v, want %v", index.CheckIssues["C-001"], want)
	}
}

// TestLoadV2Index_issueLinkMapsEmptyWithoutReferences proves a project whose
// Issues point at nothing loads with empty maps rather than an error.
func TestLoadV2Index_issueLinkMapsEmptyWithoutReferences(t *testing.T) {
	root := t.TempDir()
	writeV2LinkedProject(t, root)
	writeV2IssueFixture(t, root, "I-001-alpha.md", "I-001", "open")

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
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{fileName: "I-001-first.md", id: "I-001", tasks: []string{"T-997"}})
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{fileName: "I-002-second.md", id: "I-002", tasks: []string{"T-998"}})
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{fileName: "I-003-third.md", id: "I-003", tasks: []string{"T-999"}})

	_, err1 := LoadV2Index(root)
	_, err2 := LoadV2Index(root)
	if err1 == nil || err2 == nil || err1.Error() != err2.Error() {
		t.Fatalf("LoadV2Index() errors not deterministic across runs: %v vs %v", err1, err2)
	}
	if !strings.Contains(err1.Error(), "I-001") {
		t.Errorf("LoadV2Index() error = %v, want the lowest-ID issue reported first", err1)
	}
}

// TestLoadV2Index_issueResolutionRequiredWhenResolved proves a resolved Issue
// with no resolution block fails closed rather than being read as silently
// closed by status alone.
func TestLoadV2Index_issueResolutionRequiredWhenResolved(t *testing.T) {
	root := t.TempDir()
	writeV2LinkedProject(t, root)
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{fileName: "I-001-alpha.md", id: "I-001", status: "resolved"})

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
				fileName: "I-001-alpha.md", id: "I-001", status: status,
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
			resolution: verifiedResolution("C-999"),
			wantErr:    ErrV2IssueResolutionMissingProof,
		},
		{
			name: "proof check recorded needs work",
			setup: func(t *testing.T, root string) {
				testutil.WriteFile(t, filepath.Join(root, v2ChecksDirName, "C-001-alpha.md"),
					"---\nid: C-001\nscope: {kind: task, id: T-001}\nresult: NEEDS WORK\nchecked_by: {role: checker, session: sess-1}\nexecuted_session: build-fixture\nchecked_at: '2026-09-14T00:00:00Z'\n---\n\n# Check\n")
			},
			resolution: verifiedResolution("C-001"),
			wantErr:    ErrV2IssueResolutionUnusableProof,
		},
		{
			name: "proof check not listed on the issue",
			setup: func(t *testing.T, root string) {
				writeV2CheckFixture(t, root, "C-001-alpha.md", "C-001", "task", "T-001")
			},
			resolution: verifiedResolution("C-001"),
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
				fileName: "I-001-alpha.md", id: "I-001", status: "resolved", resolution: tt.resolution,
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
	writeV2CheckFixture(t, root, "C-001-alpha.md", "C-001", "task", "T-001")
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{
		fileName: "I-001-alpha.md", id: "I-001", status: "resolved", checks: []string{"C-001"},
		resolution: "{disposition: verified, check: C-001, actor: {role: checker, session: sess-1}, at: '2026-09-16T00:00:00Z', reason: repaired and rechecked}",
	})

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	if index.Issues["I-001"].Resolution == nil || index.Issues["I-001"].Resolution.Disposition != IssueDispositionVerified {
		t.Errorf("Resolution = %+v, want verified", index.Issues["I-001"].Resolution)
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
			resolution: "{disposition: accepted, check: C-001, actor: {role: owner, session: owner-1}, at: '2026-09-16T00:00:00Z', reason: accepted risk}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			writeV2LinkedProject(t, root)
			writeV2CheckFixture(t, root, "C-001-alpha.md", "C-001", "task", "T-001")
			writeV2LinkedIssueFixture(t, root, v2IssueFixture{
				fileName: "I-001-alpha.md", id: "I-001", status: "resolved", resolution: tt.resolution,
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
		fileName: "I-001-alpha.md", id: "I-001", status: "resolved",
		resolution: "{disposition: accepted, actor: {role: owner, session: owner-1}, at: '2026-09-16T00:00:00Z', reason: accepted risk}",
	})

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	if index.Issues["I-001"].Resolution == nil || index.Issues["I-001"].Resolution.Disposition != IssueDispositionAccepted {
		t.Errorf("Resolution = %+v, want accepted", index.Issues["I-001"].Resolution)
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
			resolution:  "{disposition: duplicate, check: C-001, actor: {role: owner, session: owner-1}, at: '2026-09-16T00:00:00Z', reason: same as I-002}",
			duplicateOf: "I-002",
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
			writeV2CheckFixture(t, root, "C-001-alpha.md", "C-001", "task", "T-001")
			writeV2IssueFixture(t, root, "I-002-canonical.md", "I-002", "open")
			writeV2LinkedIssueFixture(t, root, v2IssueFixture{
				fileName: "I-001-alpha.md", id: "I-001", status: "resolved",
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
	writeV2IssueFixture(t, root, "I-002-canonical.md", "I-002", "open")
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{
		fileName: "I-001-alpha.md", id: "I-001", status: "resolved", duplicateOf: "I-002",
		resolution: "{disposition: duplicate, actor: {role: owner, session: owner-1}, at: '2026-09-16T00:00:00Z', reason: same as I-002}",
	})

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	if index.Issues["I-001"].Resolution == nil || index.Issues["I-001"].Resolution.Disposition != IssueDispositionDuplicate {
		t.Errorf("Resolution = %+v, want duplicate", index.Issues["I-001"].Resolution)
	}
}

// TestLoadV2Index_issueEscalatedResolutionObligations proves escalated
// closure requires a planner actor and escalated_to naming the Objective, and
// must not name a proof Check, because escalation proves nothing itself.
func TestLoadV2Index_issueEscalatedResolutionObligations(t *testing.T) {
	tests := []struct {
		name        string
		resolution  string
		escalatedTo string
	}{
		{
			name:        "names a proof check",
			resolution:  "{disposition: escalated, check: C-001, actor: {role: planner, session: planner-1}, at: '2026-09-16T00:00:00Z', reason: promoted to O-001}",
			escalatedTo: "O-001",
		},
		{
			name:        "missing escalated_to",
			resolution:  "{disposition: escalated, actor: {role: planner, session: planner-1}, at: '2026-09-16T00:00:00Z', reason: promoted}",
			escalatedTo: "",
		},
		{
			name:        "non-planner actor",
			resolution:  "{disposition: escalated, actor: {role: owner, session: owner-1}, at: '2026-09-16T00:00:00Z', reason: promoted to O-001}",
			escalatedTo: "O-001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			writeV2LinkedProject(t, root)
			writeV2CheckFixture(t, root, "C-001-alpha.md", "C-001", "task", "T-001")
			writeV2LinkedIssueFixture(t, root, v2IssueFixture{
				fileName: "I-001-alpha.md", id: "I-001", status: "resolved",
				resolution: tt.resolution, escalatedTo: tt.escalatedTo,
			})

			_, err := LoadV2Index(root)
			if !errors.Is(err, ErrV2IssueResolutionFieldMismatch) {
				t.Fatalf("LoadV2Index() error = %v, want ErrV2IssueResolutionFieldMismatch", err)
			}
		})
	}
}

// TestLoadV2Index_issueEscalatedResolutionValid proves a planner actor with
// escalated_to naming an existing Objective and no named proof Check closes
// the Issue as escalated.
func TestLoadV2Index_issueEscalatedResolutionValid(t *testing.T) {
	root := t.TempDir()
	writeV2LinkedProject(t, root)
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{
		fileName: "I-001-alpha.md", id: "I-001", status: "resolved", escalatedTo: "O-001",
		resolution: "{disposition: escalated, actor: {role: planner, session: planner-1}, at: '2026-09-16T00:00:00Z', reason: promoted to O-001}",
	})

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	if index.Issues["I-001"].Resolution == nil || index.Issues["I-001"].Resolution.Disposition != IssueDispositionEscalated {
		t.Errorf("Resolution = %+v, want escalated", index.Issues["I-001"].Resolution)
	}
	if index.Issues["I-001"].EscalatedTo != "O-001" {
		t.Errorf("Issues[I-001].EscalatedTo = %q, want O-001", index.Issues["I-001"].EscalatedTo)
	}
}

// TestLoadV2Index_issueEscalatedToMissingTarget proves escalated_to must name
// an Objective that actually exists.
func TestLoadV2Index_issueEscalatedToMissingTarget(t *testing.T) {
	root := t.TempDir()
	writeV2LinkedProject(t, root)
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{
		fileName: "I-001-alpha.md", id: "I-001", status: "resolved", escalatedTo: "O-999",
		resolution: "{disposition: escalated, actor: {role: planner, session: planner-1}, at: '2026-09-16T00:00:00Z', reason: promoted}",
	})

	_, err := LoadV2Index(root)
	if !errors.Is(err, ErrV2IssueMissingEscalationTarget) {
		t.Fatalf("LoadV2Index() error = %v, want ErrV2IssueMissingEscalationTarget", err)
	}
	if !strings.Contains(err.Error(), "I-001") || !strings.Contains(err.Error(), "O-999") {
		t.Errorf("LoadV2Index() error = %v, want both the issue and its missing target named", err)
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
	writeV2CheckFixture(t, root, "C-001-first.md", "C-001", "task", "T-001")
	testutil.WriteFile(t, filepath.Join(root, v2ChecksDirName, "C-002-rerun.md"),
		"---\nid: C-002\nscope: {kind: task, id: T-001}\nresult: CLEAR\nchecked_by: {role: checker, session: sess-2}\nexecuted_session: build-fixture\nchecked_at: '2026-09-15T00:00:00Z'\nsupersedes: C-001\n---\n\n# Check\n")
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{
		fileName: "I-001-alpha.md", id: "I-001", status: "resolved", checks: []string{"C-001"},
		resolution: "{disposition: verified, check: C-001, actor: {role: checker, session: sess-1}, at: '2026-09-14T00:00:00Z', reason: repaired}",
	})

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v, want the superseded-but-still-CLEAR proof to satisfy verified at load", err)
	}
	if index.LatestCheck["T-001"] != "C-002" {
		t.Fatalf("LatestCheck[T-001] = %q, want C-002 (the reference case for this test)", index.LatestCheck["T-001"])
	}
}

// TestLoadV2Index_issueReopeningClearsResolution proves an Issue returned to
// open with a stale resolution still present is refused, while the same I-###
// identity with the resolution cleared and a dated reopened history entry
// loads cleanly as the one recurring record — never a second Issue.
func TestLoadV2Index_issueReopeningClearsResolution(t *testing.T) {
	root := t.TempDir()
	writeV2LinkedProject(t, root)
	writeV2CheckFixture(t, root, "C-001-alpha.md", "C-001", "task", "T-001")

	writeV2LinkedIssueFixture(t, root, v2IssueFixture{
		fileName: "I-001-alpha.md", id: "I-001", status: "open", checks: []string{"C-001"},
		resolution: "{disposition: verified, check: C-001, actor: {role: checker, session: sess-1}, at: '2026-09-15T00:00:00Z', reason: repaired}",
	})
	if _, err := LoadV2Index(root); !errors.Is(err, ErrV2IssueResolutionNotAllowed) {
		t.Fatalf("LoadV2Index() error = %v, want ErrV2IssueResolutionNotAllowed for the stale resolution", err)
	}

	writeV2LinkedIssueFixture(t, root, v2IssueFixture{
		fileName: "I-001-alpha.md", id: "I-001", status: "open", checks: []string{"C-001"},
		history: "[{at: '2026-09-16T00:00:00Z', actor: {role: owner, session: owner-1}, kind: reopened, note: regressed}]",
	})

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v, want the reopened issue to load once resolution is cleared", err)
	}
	if index.Issues["I-001"].Status != IssueStatusOpen || index.Issues["I-001"].Resolution != nil {
		t.Errorf("Issues[I-001] = %+v, want open with no resolution", index.Issues["I-001"])
	}
	if len(index.Issues) != 1 {
		t.Fatalf("Issues = %+v, want the single I-001 identity, not a second record for the same recurrence", index.Issues)
	}
}

// TestLoadV2Index_issueDeferralStaysOpen proves a deferral is a dated history
// entry on an Issue that stays open, not a fourth lifecycle status.
func TestLoadV2Index_issueDeferralStaysOpen(t *testing.T) {
	root := t.TempDir()
	writeV2LinkedProject(t, root)
	writeV2LinkedIssueFixture(t, root, v2IssueFixture{
		fileName: "I-001-alpha.md", id: "I-001", status: "open",
		history: "[{at: '2026-09-16T00:00:00Z', actor: {role: owner, session: owner-1}, kind: deferred, note: revisit next release}]",
	})

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	if index.Issues["I-001"].Status != IssueStatusOpen {
		t.Errorf("Status = %q, want open", index.Issues["I-001"].Status)
	}
	if len(index.Issues["I-001"].History) != 1 || index.Issues["I-001"].History[0].Kind != IssueHistoryDeferred {
		t.Errorf("History = %+v, want one deferred entry", index.Issues["I-001"].History)
	}
}

func TestCheckRuntimeSchemaRejectsDirectoryReadFailure(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, ".savepoint", "config.yml")
	testutil.MkdirAll(t, configPath) // a directory where config.yml should be a file

	err := CheckRuntimeSchema(root)
	if err == nil {
		t.Fatal("CheckRuntimeSchema() error = nil, want error reading config.yml")
	}
	if errors.Is(err, ErrMalformedSchemaVersion) || errors.Is(err, ErrUnsupportedSchemaVersion) || errors.Is(err, ErrSchemaMigrationRequired) {
		t.Errorf("CheckRuntimeSchema() error = %v, want a read failure, not a schema diagnostic", err)
	}
}

// A failure after the Task is persisted must never leave a Task behind with a
// nonzero result: lock-release failure rolls the new file back and keeps the ID retired.
func TestWithTaskIDReservationLockReleaseFailureRollsBackAndRetiresID(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O-001-first", "O-001", "First objective")
	created := filepath.Join(root, "created.tmp")

	id, err := withTaskIDReservation(root, nil, func(_ string, _ *V2Index) (func() error, error) {
		if err := os.WriteFile(created, []byte("task"), 0o644); err != nil {
			return nil, err
		}
		// Make the later lock release fail.
		if err := os.Remove(filepath.Join(root, taskIDLockFile)); err != nil {
			return nil, err
		}
		return func() error { return os.Remove(created) }, nil
	})
	if err == nil || id != "" {
		t.Fatalf("withTaskIDReservation() = %q, %v; want empty ID and release-lock error", id, err)
	}
	if _, statErr := os.Stat(created); !errors.Is(statErr, os.ErrNotExist) {
		t.Errorf("created file stat error = %v, want rolled back", statErr)
	}
	if got, readErr := os.ReadFile(filepath.Join(root, taskIDHighWaterFile)); readErr != nil || string(got) != "last_issued: 1\n" {
		t.Errorf("high-water mark = %q, error = %v; ID must stay retired", got, readErr)
	}
}

// If rollback itself fails the Task remains, so the error must say so and name the cause.
func TestWithTaskIDReservationRollbackFailureIsReported(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O-001-first", "O-001", "First objective")

	_, err := withTaskIDReservation(root, nil, func(string, *V2Index) (func() error, error) {
		return func() error { return errors.New("stuck file") }, errors.New("strict load failed")
	})
	if err == nil || !strings.Contains(err.Error(), "strict load failed") || !strings.Contains(err.Error(), "rollback created Task: stuck file") {
		t.Fatalf("error = %v, want operation error joined with rollback failure", err)
	}
}
