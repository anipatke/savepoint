package data

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/testutil"
)

const planTaskHead = "---\nid: %s\ntitle: \"T\"\nobjective: O-001\nplanned_by: {role: planner, session: s}\nstatus: planned\n"

func planTask(id, extra string) string {
	return strings.Replace(planTaskHead, "%s", id, 1) + extra + "---\n\n# body\n"
}

func decodePlanTask(t *testing.T, extra string) TaskPlanV2 {
	t.Helper()
	task, err := DecodeTaskV2("T-001.md", planTask("T-001", extra))
	if err != nil {
		t.Fatalf("DecodeTaskV2() error = %v; advisory metadata must not fail a load", err)
	}
	return task.Plan
}

func TestDecodeTaskV2_planOmittedIsUnknownAndEmptyIsReviewed(t *testing.T) {
	none := decodePlanTask(t, "")
	if none.Lane != "" || none.Reads.Declared() || none.Writes.Declared() || len(none.Diagnostics()) != 0 {
		t.Errorf("old task plan = %+v, want entirely unknown and clean", none)
	}

	empty := decodePlanTask(t, "lane: core\nplanned_reads: []\nplanned_writes: [internal/new/file.go]\n")
	if !empty.Reads.Declared() || len(empty.Reads.Paths()) != 0 {
		t.Errorf("Reads = %+v, want declared and empty", empty.Reads)
	}
	if !empty.Writes.Contains("internal/new/file.go") || empty.Lane != "core" || len(empty.Diagnostics()) != 0 {
		t.Errorf("plan = %+v, want new path accepted without probing the filesystem", empty)
	}
}

func TestDecodeTaskV2_planReturnsCopies(t *testing.T) {
	plan := decodePlanTask(t, "planned_writes: [a.go]\n")
	plan.Writes.Paths()[0] = "changed.go"
	if !plan.Writes.Contains("a.go") {
		t.Error("mutating the returned paths changed the plan")
	}
}

func TestDecodeTaskV2_malformedPlanIsNonfatalAndNamed(t *testing.T) {
	cases := map[string]string{
		"glob":          "planned_writes: [internal/*.go]\n",
		"directory":     "planned_writes: [internal/data/]\n",
		"absolute":      "planned_reads: [/etc/passwd]\n",
		"drive":         "planned_reads: ['C:/x/y.go']\n",
		"unc":           "planned_reads: ['//host/share/f.go']\n",
		"backslash":     "planned_reads: ['a\\b.go']\n",
		"traversal":     "planned_writes: [../outside.go]\n",
		"dot segment":   "planned_writes: [./a.go]\n",
		"empty segment": "planned_writes: [a//b.go]\n",
		"trailing dot":  "planned_writes: ['a/b.']\n",
		"case alias":    "planned_writes: [A.go, a.go]\n",
		"alias across":  "planned_reads: [A.go]\nplanned_writes: [a.go]\n",
		"not a list":    "planned_writes: a.go\n",
		"null list":     "planned_writes:\n",
		"non string":    "planned_writes: [12]\n",
		"bad lane":      "lane: Not A Lane\n",
		"lane list":     "lane: [core]\n",
	}
	for name, extra := range cases {
		t.Run(name, func(t *testing.T) {
			plan := decodePlanTask(t, extra)
			diagnostics := plan.Diagnostics()
			if len(diagnostics) == 0 {
				t.Fatalf("plan %+v has no diagnostic", plan)
			}
			if got := diagnostics[0].String(); !strings.Contains(got, "T-001") || !strings.Contains(got, "T-001.md") {
				t.Errorf("diagnostic %q must name the record and its file", got)
			}
		})
	}
}

func TestDecodeObjectiveV2_lanesAreNonfatalAndStable(t *testing.T) {
	good, err := DecodeObjectiveV2("O.md", "---\nid: O-001\ntitle: O\nstatus: planned\nlanes:\n  - {key: core, title: Core lane}\n  - {key: board, title: Board lane}\n---\n")
	if err != nil {
		t.Fatalf("DecodeObjectiveV2() error = %v", err)
	}
	if len(good.Plan.Lanes) != 2 || good.Plan.Lanes[0].Key != "core" || good.Plan.Lanes[1].Title != "Board lane" {
		t.Errorf("Lanes = %+v, want authored order kept", good.Plan.Lanes)
	}
	if _, ok := good.Plan.Lane("board"); !ok {
		t.Error("Lane(board) not found")
	}

	bad, err := DecodeObjectiveV2("O.md", "---\nid: O-001\ntitle: O\nstatus: planned\nlanes:\n  - {key: core, title: A}\n  - {key: core, title: B}\n  - {key: Bad Key, title: C}\n  - {key: ok}\n  - nonsense\n---\n")
	if err != nil {
		t.Fatalf("DecodeObjectiveV2() error = %v", err)
	}
	if len(bad.Plan.Lanes) != 0 || len(bad.Plan.Diagnostics()) != 4 {
		t.Errorf("Lanes = %+v, diagnostics = %v; want the ambiguous lane dropped and four diagnostics", bad.Plan.Lanes, bad.Plan.Diagnostics())
	}
}

func TestComparePlanPaths_neverCertifiesCaseAliases(t *testing.T) {
	if ComparePlanPaths("a/B.go", "a/B.go") != PlanPathsSame || ComparePlanPaths("a/B.go", "a/b.go") != PlanPathsMayAlias || ComparePlanPaths("a.go", "b.go") != PlanPathsDistinct {
		t.Error("ComparePlanPaths relations are wrong")
	}
	if err := ValidatePlanPath("internal/data/brand_new.go"); err != nil {
		t.Errorf("ValidatePlanPath(new file) = %v, want valid", err)
	}
}

func writePlanProject(t *testing.T, objectiveExtra string, tasks map[string]string) string {
	t.Helper()
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, v2ObjectivesDirName, "O-001-x", v2ObjectiveFileName),
		"---\nid: O-001\ntitle: O\nstatus: planned\nlanes:\n  - {key: core, title: Core}\n"+objectiveExtra+"---\n")
	for id, extra := range tasks {
		testutil.WriteFile(t, filepath.Join(root, v2ObjectivesDirName, "O-001-x", v2TasksDirName, id+"-t.md"), planTask(id, extra))
	}
	return root
}

func TestLoadV2Index_oldProjectsAreCleanAndUnknown(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O-001-first", "O-001", "First")
	writeV2TaskFixture(t, root, "O-001-first", "T-001-a.md", "T-001", "A", "O-001")
	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	if len(index.PlanDiagnostics) != 0 || len(index.UnusablePlans) != 0 || index.Tasks["T-001"].Plan.Reads.Declared() {
		t.Errorf("old project planning state = %+v, want clean and unknown", index.PlanDiagnostics)
	}
}

func TestLoadV2Index_advisoryProblemsNameRecordsAndKeepLifecycleGates(t *testing.T) {
	root := writePlanProject(t, "", map[string]string{
		"T-001": "lane: core\nplanned_writes: [a.go]\n",
		"T-002": "lane: missing\nplanned_writes: [b.go]\n",
		"T-003": "planned_writes: [c/*.go]\n",
	})
	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v; advisory metadata must not block loading", err)
	}
	if index.UnusablePlans["T-001"] || !index.UnusablePlans["T-002"] || !index.UnusablePlans["T-003"] {
		t.Errorf("UnusablePlans = %v, want only T-002 and T-003", index.UnusablePlans)
	}
	if len(index.PlanDiagnostics) != 2 || index.PlanDiagnostics[0].Record != "T-002" || index.PlanDiagnostics[1].Record != "T-003" {
		t.Errorf("PlanDiagnostics = %v, want T-002 then T-003", index.PlanDiagnostics)
	}
	if got := index.PlanDiagnostics[0].String(); !strings.Contains(got, "T-002-t.md") || !strings.Contains(got, `"missing"`) {
		t.Errorf("diagnostic %q must name the file and unresolved lane", got)
	}
	// Bad advice never changes a gate: an unfinished target still blocks for
	// the ordinary reason, and nothing mentions planning.
	decision := ResolveTaskDependencyV2(index, TaskDependencyV2{Task: "T-002", Requires: TaskDependencyClear})
	if decision.Satisfied || decision.Block == nil || decision.Block.Kind != DependencyBlockNotDone {
		t.Errorf("decision = %+v, want the ordinary not-done block", decision)
	}
}

func independenceYAML(reviewed string, overlap string, reason string) string {
	return "independence:\n  - tasks: [T-001, T-002]\n    overlaps:\n      - " + overlap + "\n    reason: \"" + reason + "\"\n    reviewed: " + reviewed + "\n"
}

func TestLoadV2Index_independenceBindsToReviewedManifests(t *testing.T) {
	tasks := map[string]string{
		"T-001": "planned_writes: [shared.go]\n",
		"T-002": "planned_reads: [shared.go]\n",
	}
	probe, err := LoadV2Index(writePlanProject(t, "", tasks))
	if err != nil {
		t.Fatal(err)
	}
	digest := PlanReviewDigest("T-002", probe.Tasks["T-002"].Plan, "T-001", probe.Tasks["T-001"].Plan)
	if digest != PlanReviewDigest("T-001", probe.Tasks["T-001"].Plan, "T-002", probe.Tasks["T-002"].Plan) {
		t.Fatal("digest depends on argument order")
	}
	overlap := "{writer: T-001, reader: T-002, path: shared.go}"

	valid, err := LoadV2Index(writePlanProject(t, independenceYAML(digest, overlap, "reads a stable copy"), tasks))
	if err != nil || len(valid.PlanDiagnostics) != 0 || len(valid.PlanIndependence["O-001"]) != 1 {
		t.Fatalf("valid explanation: err = %v, diagnostics = %v, kept = %v", err, valid.PlanDiagnostics, valid.PlanIndependence)
	}

	broadened := map[string]string{"T-001": "planned_writes: [shared.go, other.go]\n", "T-002": tasks["T-002"]}
	cases := map[string]struct {
		extra string
		tasks map[string]string
		want  string
	}{
		"stale after manifest edit": {independenceYAML(digest, overlap, "r"), broadened, "stale"},
		"unreviewed":                {independenceYAML("''", overlap, "r"), tasks, "stale or unreviewed"},
		"empty reason":              {independenceYAML(digest, overlap, " "), tasks, "non-empty reason"},
		"wrong direction":           {independenceYAML(digest, "{writer: T-002, reader: T-001, path: shared.go}", "r"), tasks, "not a write by T-002"},
		"path not in both":          {independenceYAML(digest, "{writer: T-001, reader: T-002, path: nope.go}", "r"), tasks, "not a write by T-001"},
		"third task":                {independenceYAML(digest, "{writer: T-001, reader: T-009, path: shared.go}", "r"), tasks, "one Task as writer"},
		"missing task":              {strings.Replace(independenceYAML(digest, overlap, "r"), "T-002]", "T-009]", 1), tasks, "missing Task T-009"},
		"duplicate declaration":     {independenceYAML(digest, overlap, "r") + "  - tasks: [T-002, T-001]\n    overlaps: [" + overlap + "]\n    reason: again\n    reviewed: " + digest + "\n", tasks, "more than once"},
		"shared write":              {independenceYAML(digest, overlap, "r"), map[string]string{"T-001": tasks["T-001"], "T-002": "planned_reads: [shared.go]\nplanned_writes: [shared.go]\n"}, "shared write"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			index, err := LoadV2Index(writePlanProject(t, tc.extra, tc.tasks))
			if err != nil {
				t.Fatalf("LoadV2Index() error = %v; advisory metadata must not block loading", err)
			}
			if len(index.PlanIndependence["O-001"]) != 0 {
				t.Errorf("explanation kept: %v", index.PlanIndependence)
			}
			found := false
			for _, d := range index.PlanDiagnostics {
				found = found || strings.Contains(d.String(), tc.want)
			}
			if !found {
				t.Errorf("diagnostics %v lack %q", index.PlanDiagnostics, tc.want)
			}
		})
	}
}

func TestLoadV2Index_laneReferencesResolveOnlyInOwner(t *testing.T) {
	root := writePlanProject(t, "", map[string]string{"T-001": "lane: core\n"})
	testutil.WriteFile(t, filepath.Join(root, v2ObjectivesDirName, "O-002-y", v2ObjectiveFileName), "---\nid: O-002\ntitle: Other\nstatus: planned\n---\n")
	other := strings.Replace(planTask("T-002", "lane: core\n"), "objective: O-001", "objective: O-002", 1)
	testutil.WriteFile(t, filepath.Join(root, v2ObjectivesDirName, "O-002-y", v2TasksDirName, "T-002-t.md"), other)
	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatal(err)
	}
	if index.UnusablePlans["T-001"] || !index.UnusablePlans["T-002"] {
		t.Errorf("UnusablePlans = %v, want only the Task naming another Objective's lane", index.UnusablePlans)
	}
}
