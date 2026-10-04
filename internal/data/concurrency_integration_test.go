package data

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/opencode/savepoint/internal/testutil"
)

// These tests exercise the lane-planning projection over temporary projects
// loaded from real files, so records, strict decoding, the index and the
// projection are proven together rather than through in-memory fixtures.

const integrationLanes = "lanes:\n  - {key: core, title: Core}\n  - {key: board, title: Board}\n  - {key: docs, title: Docs}\n"

// writeIntegrationProject writes O-001 with objectiveExtra in its
// frontmatter and one Task per entry of tasks (extra frontmatter lines).
func writeIntegrationProject(t *testing.T, objectiveExtra string, tasks map[string]string) string {
	t.Helper()
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, v2ObjectivesDirName, "O-001-x", v2ObjectiveFileName),
		"---\nid: O-001\ntitle: O\nstatus: in_progress\n"+objectiveExtra+"---\n")
	for id, extra := range tasks {
		testutil.WriteFile(t, filepath.Join(root, v2ObjectivesDirName, "O-001-x", v2TasksDirName, id+"-t.md"), integrationTask(id, extra))
	}
	return root
}

// integrationTask writes a planned Task unless extra names its own status.
func integrationTask(id, extra string) string {
	head := strings.Replace(planTaskHead, "status: planned\n", "", 1)
	if !strings.Contains(extra, "status:") {
		head += "status: planned\n"
	}
	return strings.Replace(head, "%s", id, 1) + extra + "---\n\n# body\n"
}

func integrationConcurrency(t *testing.T, root string, enabled bool) (*V2Index, *ConcurrencyV2) {
	t.Helper()
	index := mustLoadV2Index(t, root)
	return index, ResolveConcurrencyV2(index, "O-001", ConcurrencyOptionsV2{Enabled: enabled})
}

// TestIntegration_lanesAcrossOldMixedAndNewRecords proves membership and
// advice for a project with no planning metadata, a mixed one and a fully
// planned one, including a record that ignores its lane entirely.
func TestIntegration_lanesAcrossOldMixedAndNewRecords(t *testing.T) {
	t.Run("old records have no lanes and no advice", func(t *testing.T) {
		root := writeIntegrationProject(t, "", map[string]string{"T-001": "", "T-002": ""})
		index, c := integrationConcurrency(t, root, true)
		if len(index.PlanDiagnostics) != 0 || len(c.Lanes) != 0 || len(c.Candidates) != 0 || len(c.Groups) != 0 {
			t.Fatalf("old project: diagnostics %v, lanes %v, candidates %v, groups %v; want none", index.PlanDiagnostics, c.Lanes, c.Candidates, c.Groups)
		}
		if n := noteFor(c, "T-001", ConcurrencyNoLane); n == nil {
			t.Errorf("notes = %+v, want T-001 reported as having no lane", c.Notes)
		}
	})

	t.Run("mixed records advise only what is planned", func(t *testing.T) {
		root := writeIntegrationProject(t, integrationLanes, map[string]string{
			"T-001": "",
			"T-002": "lane: core\nplanned_reads: [a.go]\nplanned_writes: [b.go]\n",
			"T-003": "lane: board\nplanned_reads: [a.go]\nplanned_writes: [c.go]\n",
			"T-004": "lane: docs\n", // lane without manifests: unknown scope
		})
		_, c := integrationConcurrency(t, root, true)
		wantGroups(t, c, [][]string{{"T-002", "T-003"}})
		if noteFor(c, "T-001", ConcurrencyNoLane) == nil {
			t.Errorf("notes = %+v, want no_lane for T-001", c.Notes)
		}
		if noteFor(c, "T-004", ConcurrencyUnknownScope) == nil {
			t.Errorf("notes = %+v, want unknown_scope for T-004", c.Notes)
		}
	})

	t.Run("new records keep membership with advice off", func(t *testing.T) {
		root := writeIntegrationProject(t, integrationLanes, map[string]string{
			"T-001": "lane: core\nplanned_reads: [a.go]\nplanned_writes: [b.go]\n",
			"T-002": "lane: board\nplanned_reads: [a.go]\nplanned_writes: [c.go]\n",
		})
		_, off := integrationConcurrency(t, root, false)
		_, on := integrationConcurrency(t, root, true)
		if !reflect.DeepEqual(off.Lanes, on.Lanes) {
			t.Errorf("lanes differ with the preference: off %+v, on %+v", off.Lanes, on.Lanes)
		}
		if off.Groups != nil || off.Candidates != nil || off.Notes != nil {
			t.Errorf("off projection carries advice: %+v", off)
		}
		if len(on.Lanes) != 3 || !on.Lanes[2].Empty() {
			t.Errorf("lanes = %+v, want three declared lanes with docs empty", on.Lanes)
		}
	})
}

// TestIntegration_sequentialDependencyAndSharedPathCases covers same-lane
// sequence, satisfied and blocked dependencies, shared writes and shared reads.
func TestIntegration_sequentialDependencyAndSharedPathCases(t *testing.T) {
	t.Run("same lane is sequential", func(t *testing.T) {
		root := writeIntegrationProject(t, integrationLanes, map[string]string{
			"T-001": "lane: core\nplanned_reads: []\nplanned_writes: [a.go]\n",
			"T-002": "lane: core\nplanned_reads: []\nplanned_writes: [b.go]\n",
		})
		_, c := integrationConcurrency(t, root, true)
		if got := candidateIDs(c); !reflect.DeepEqual(got, []string{"T-001"}) {
			t.Errorf("candidates = %v, want only the first Task of the lane", got)
		}
		if n := noteFor(c, "T-002", ConcurrencySameLane); n == nil || n.Tasks[1] != "T-001" {
			t.Errorf("notes = %+v, want same_lane naming T-001", c.Notes)
		}
	})

	t.Run("unfinished dependency blocks the start", func(t *testing.T) {
		root := writeIntegrationProject(t, integrationLanes, map[string]string{
			"T-001": "lane: core\nplanned_reads: []\nplanned_writes: [a.go]\n",
			"T-002": "lane: board\nplanned_reads: []\nplanned_writes: [b.go]\ndepends_on: [{task: T-001, requires: clear}]\n",
		})
		index, c := integrationConcurrency(t, root, true)
		wantGroups(t, c, nil)
		n := noteFor(c, "T-002", ConcurrencyStartBlocked)
		if n == nil || !strings.Contains(n.Detail, "T-001") {
			t.Errorf("notes = %+v, want start_blocked repeating the dependency blocker", c.Notes)
		}
		if ResolveTaskStart(index, "T-002").Allowed {
			t.Error("ordinary start gate allowed a Task with an unmet dependency")
		}
	})

	t.Run("satisfied dependency is no obstacle", func(t *testing.T) {
		root := writeIntegrationProject(t, integrationLanes, map[string]string{
			"T-001": "status: done\n",
			"T-002": "lane: core\nplanned_reads: []\nplanned_writes: [a.go]\ndepends_on: [{task: T-001, requires: clear}]\n",
			"T-003": "lane: board\nplanned_reads: []\nplanned_writes: [b.go]\n",
		})
		writeScenarioCheck(t, root, "C-001", CheckScopeTask, "T-001", "sess-1", time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC), "")
		index, c := integrationConcurrency(t, root, true)
		if !ResolveTaskStart(index, "T-002").Allowed {
			t.Fatal("dependency fixture is not satisfied")
		}
		wantGroups(t, c, [][]string{{"T-002", "T-003"}})
	})

	t.Run("shared write is withheld, shared read alone is not", func(t *testing.T) {
		root := writeIntegrationProject(t, integrationLanes, map[string]string{
			"T-001": "lane: core\nplanned_reads: [shared.go]\nplanned_writes: [a.go, same.go]\n",
			"T-002": "lane: board\nplanned_reads: [shared.go]\nplanned_writes: [b.go, same.go]\n",
			"T-003": "lane: docs\nplanned_reads: [shared.go]\nplanned_writes: [c.go]\n",
		})
		_, c := integrationConcurrency(t, root, true)
		if pair := pairFor(c, "T-001", "T-002"); pair == nil || pair.Together || pair.Reason.Kind != ConcurrencySharedWrite {
			t.Errorf("T-001/T-002 = %+v, want shared_write", pair)
		}
		for _, other := range []string{"T-002", "T-001"} {
			a, b := "T-003", other
			if a > b {
				a, b = b, a
			}
			if pair := pairFor(c, a, b); pair == nil || !pair.Together {
				t.Errorf("%s/%s = %+v, want together on a shared read alone", a, b, pair)
			}
		}
	})

	t.Run("active work with unknown scope withholds a start", func(t *testing.T) {
		root := writeIntegrationProject(t, integrationLanes, map[string]string{
			"T-001": "status: in_progress\nstage: build\nlane: core\n", // running, no manifests
			"T-002": "lane: board\nplanned_reads: []\nplanned_writes: [b.go]\n",
		})
		_, c := integrationConcurrency(t, root, true)
		if len(c.Active) != 1 || c.Active[0] != "T-001" {
			t.Fatalf("Active = %v, want T-001", c.Active)
		}
		if candidateIDs(c) != nil && len(candidateIDs(c)) != 0 {
			t.Errorf("candidates = %v, want none while the running Task's scope is unknown", candidateIDs(c))
		}
		if n := noteFor(c, "T-002", ConcurrencyUnknownScope); n == nil {
			t.Errorf("notes = %+v, want unknown_scope for T-002", c.Notes)
		}
	})
}

// TestIntegration_explanationInvalidatesWhenEitherManifestChanges proves a
// recorded read/write explanation holds only for the reviewed manifests, on
// real files edited between loads.
func TestIntegration_explanationInvalidatesWhenEitherManifestChanges(t *testing.T) {
	writer := "lane: core\nplanned_reads: []\nplanned_writes: [shared.go]\n"
	reader := "lane: board\nplanned_reads: [shared.go]\nplanned_writes: []\n"
	probe := mustLoadV2Index(t, writeIntegrationProject(t, integrationLanes, map[string]string{"T-001": writer, "T-002": reader}))
	digest := PlanReviewDigest("T-001", probe.Tasks["T-001"].Plan, "T-002", probe.Tasks["T-002"].Plan)
	explained := integrationLanes + independenceYAML(digest, "{writer: T-001, reader: T-002, path: shared.go}", "reads a stable copy")

	root := writeIntegrationProject(t, explained, map[string]string{"T-001": writer, "T-002": reader})
	_, c := integrationConcurrency(t, root, true)
	wantGroups(t, c, [][]string{{"T-001", "T-002"}})

	for name, edit := range map[string]struct{ task, extra string }{
		"writer broadened": {"T-001", "lane: core\nplanned_reads: []\nplanned_writes: [shared.go, more.go]\n"},
		"reader broadened": {"T-002", "lane: board\nplanned_reads: [shared.go, more.go]\nplanned_writes: []\n"},
	} {
		t.Run(name, func(t *testing.T) {
			tasks := map[string]string{"T-001": writer, "T-002": reader}
			tasks[edit.task] = edit.extra
			root := writeIntegrationProject(t, explained, tasks)
			_, c := integrationConcurrency(t, root, true)
			pair := pairFor(c, "T-001", "T-002")
			if pair == nil || pair.Together || pair.Reason.Kind != ConcurrencyWriteReadOverlap || !strings.Contains(pair.Reason.Detail, "stale or unusable") {
				t.Errorf("pair = %+v, want withheld with a stale-explanation note", pair)
			}
		})
	}
}

// TestIntegration_adviceNeverChangesExistingDecisions loads the same Tasks
// with and without lane and manifest metadata and with the preference on and
// off, and requires identical start, advance and completion decisions.
func TestIntegration_adviceNeverChangesExistingDecisions(t *testing.T) {
	planned := map[string]string{
		"T-001": "status: done\n",
		"T-002": "lane: core\nplanned_reads: [a.go]\nplanned_writes: [b.go]\ndepends_on: [{task: T-001, requires: clear}]\n",
		"T-003": "lane: core\nplanned_reads: [a.go]\nplanned_writes: [b.go]\n", // shares a write with T-002
		"T-004": "lane: board\nstatus: in_progress\nstage: audit\nplanned_reads: [a.go]\nplanned_writes: [b.go]\n",
		"T-005": "lane: board\nstatus: in_progress\nstage: build\nplanned_writes: not-a-list\n", // malformed advice
	}
	plain := map[string]string{
		"T-001": "status: done\n",
		"T-002": "depends_on: [{task: T-001, requires: clear}]\n",
		"T-003": "",
		"T-004": "status: in_progress\nstage: audit\n",
		"T-005": "status: in_progress\nstage: build\n",
	}
	build := func(tasks map[string]string, objectiveExtra string) *V2Index {
		root := writeIntegrationProject(t, objectiveExtra, tasks)
		writeScenarioCheck(t, root, "C-001", CheckScopeTask, "T-001", "sess-1", time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC), "")
		return mustLoadV2Index(t, root)
	}
	withLanes := build(planned, integrationLanes)
	without := build(plain, "")

	decisions := func(index *V2Index) map[string]GateDecision {
		out := map[string]GateDecision{}
		for _, id := range []string{"T-001", "T-002", "T-003", "T-004", "T-005"} {
			out["start "+id] = ResolveTaskStart(index, id)
			out["advance "+id] = ResolveTaskAdvance(index, id)
			out["complete "+id] = ResolveTaskCompletion(index, id)
		}
		return out
	}
	before := decisions(withLanes)
	for _, enabled := range []bool{false, true, false} {
		ResolveConcurrencyV2(withLanes, "O-001", ConcurrencyOptionsV2{Enabled: enabled})
		if got := decisions(withLanes); !reflect.DeepEqual(got, before) {
			t.Errorf("decisions changed after projecting with Enabled=%v", enabled)
		}
	}
	if got := decisions(without); !reflect.DeepEqual(got, before) {
		t.Errorf("decisions differ when lanes and manifests are absent:\nwith lanes: %+v\nwithout:    %+v", before, got)
	}
	if len(withLanes.PlanDiagnostics) == 0 {
		t.Error("malformed planned_writes produced no diagnostic; the degradation case is not exercised")
	}
	if !before["start T-002"].Allowed || before["start T-003"].Allowed == false {
		t.Errorf("start decisions = %+v / %+v, want ordinary readiness regardless of a shared write", before["start T-002"], before["start T-003"])
	}
}
