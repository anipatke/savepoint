package data

import (
	"reflect"
	"strings"
	"testing"
)

type concurrencyTask struct {
	id, lane string
	status   ColumnType
	reads    []string // nil is an omitted, unknown scope
	writes   []string
	deps     []string
}

func scope(paths []string) PlanScopeV2 {
	if paths == nil {
		return PlanScopeV2{}
	}
	return PlanScopeV2{declared: true, paths: paths}
}

func newConcurrencyIndex(tasks ...concurrencyTask) *V2Index {
	index := newV2TestIndex()
	index.UnusablePlans = map[string]bool{}
	index.PlanIndependence = map[string][]PlanIndependenceV2{}
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Status: ColumnInProgress, Plan: ObjectivePlanV2{Lanes: []ObjectiveLaneV2{
		{Key: "a", Title: "Lane A"}, {Key: "b", Title: "Lane B"}, {Key: "c", Title: "Lane C"},
	}}}
	for _, t := range tasks {
		task := &TaskV2{ID: t.id, Objective: "O-001", Status: t.status, Plan: TaskPlanV2{Lane: t.lane, Reads: scope(t.reads), Writes: scope(t.writes)}}
		if t.status == ColumnInProgress {
			task.Stage = StageBuild
		}
		for _, dep := range t.deps {
			task.DependsOn = append(task.DependsOn, TaskDependencyV2{Task: dep, Requires: TaskDependencyClear})
		}
		index.Tasks[t.id] = task
		index.ObjectiveTasks["O-001"] = append(index.ObjectiveTasks["O-001"], t.id)
	}
	return index
}

func planned(id, lane string, reads, writes []string, deps ...string) concurrencyTask {
	return concurrencyTask{id, lane, ColumnPlanned, reads, writes, deps}
}

func concurrencyFor(index *V2Index) *ConcurrencyV2 {
	return ResolveConcurrencyV2(index, "O-001", ConcurrencyOptionsV2{Enabled: true})
}

func candidateIDs(c *ConcurrencyV2) []string {
	ids := []string{}
	for _, candidate := range c.Candidates {
		ids = append(ids, candidate.Task)
	}
	return ids
}

func noteFor(c *ConcurrencyV2, id string, kind ConcurrencyKindV2) *ConcurrencyReasonV2 {
	for i := range c.Notes {
		if c.Notes[i].Tasks[0] == id && c.Notes[i].Kind == kind {
			return &c.Notes[i]
		}
	}
	return nil
}

func pairFor(c *ConcurrencyV2, a, b string) *ConcurrencyPairV2 {
	for i := range c.Pairs {
		if c.Pairs[i].Tasks == [2]string{a, b} {
			return &c.Pairs[i]
		}
	}
	return nil
}

func wantGroups(t *testing.T, c *ConcurrencyV2, want [][]string) {
	t.Helper()
	if !reflect.DeepEqual(c.Groups, want) {
		t.Errorf("Groups = %v, want %v (pairs %+v, notes %+v)", c.Groups, want, c.Pairs, c.Notes)
	}
}

func TestResolveConcurrencyV2_unknownObjectiveAndNilIndex(t *testing.T) {
	if ResolveConcurrencyV2(nil, "O-001", ConcurrencyOptionsV2{Enabled: true}) != nil {
		t.Error("nil index projected something")
	}
	if ResolveConcurrencyV2(newConcurrencyIndex(), "O-009", ConcurrencyOptionsV2{Enabled: true}) != nil {
		t.Error("unknown Objective projected something")
	}
}

func TestResolveConcurrencyV2_independentLanesStartTogether(t *testing.T) {
	index := newConcurrencyIndex(
		planned("T-001", "a", []string{"shared.go"}, []string{"a.go"}),
		planned("T-002", "b", []string{"shared.go"}, []string{"b.go"}), // shared reads alone are allowed
	)
	c := concurrencyFor(index)
	wantGroups(t, c, [][]string{{"T-001", "T-002"}})
	if got := candidateIDs(c); !reflect.DeepEqual(got, []string{"T-001", "T-002"}) {
		t.Errorf("candidates = %v", got)
	}
	if pair := pairFor(c, "T-001", "T-002"); pair == nil || !pair.Together || pair.Reason != nil {
		t.Errorf("pair = %+v, want together with no reason", pair)
	}
}

func TestResolveConcurrencyV2_offOrAbsentKeepsMembershipOnly(t *testing.T) {
	index := newConcurrencyIndex(
		planned("T-001", "a", []string{}, []string{"a.go"}),
		planned("T-002", "b", []string{}, []string{"b.go"}),
		concurrencyTask{"T-003", "a", ColumnInProgress, []string{}, []string{"x.go"}, nil},
		concurrencyTask{"T-004", "a", ColumnDone, nil, nil, nil},
	)
	index.PlanDiagnostics = []PlanDiagnosticV2{{Record: "T-001", Field: "planned_reads", Message: "bad"}}
	c := ResolveConcurrencyV2(index, "O-001", ConcurrencyOptionsV2{})
	if c.Enabled || c.Withheld != nil || c.Candidates != nil || c.Pairs != nil || c.Groups != nil || c.Notes != nil || c.Active != nil || c.Diagnostics != nil {
		t.Errorf("disabled projection carries advice: %+v", c)
	}
	lane := c.Lanes[0]
	if lane.Ref() != "O-001/a" || lane.Title != "Lane A" ||
		!reflect.DeepEqual(lane.Planned, []string{"T-001"}) || !reflect.DeepEqual(lane.InProgress, []string{"T-003"}) || !reflect.DeepEqual(lane.Done, []string{"T-004"}) {
		t.Errorf("lane a = %+v, want one Task in every column", lane)
	}
	if !c.Lanes[2].Empty() || c.Lanes[0].Empty() {
		t.Error("Empty() does not match membership")
	}

	enabled := concurrencyFor(index)
	if !reflect.DeepEqual(enabled.Lanes, c.Lanes) {
		t.Errorf("membership differs with the preference on: %+v vs %+v", enabled.Lanes, c.Lanes)
	}
	if len(enabled.Diagnostics) != 1 {
		t.Errorf("enabled Diagnostics = %v, want the Task diagnostic", enabled.Diagnostics)
	}
}

func TestResolveConcurrencyV2_dependenciesNeverTogether(t *testing.T) {
	// T-003 depends on T-001 only through T-002, which is not a candidate.
	index := newConcurrencyIndex(
		planned("T-001", "a", []string{}, []string{"a.go"}),
		planned("T-002", "c", []string{}, []string{"c.go"}, "T-001"),
		planned("T-003", "b", []string{}, []string{"b.go"}, "T-002"),
	)
	c := concurrencyFor(index)
	if got := candidateIDs(c); !reflect.DeepEqual(got, []string{"T-001"}) {
		t.Fatalf("candidates = %v, want only T-001", got)
	}
	wantGroups(t, c, nil)
	for _, id := range []string{"T-002", "T-003"} {
		if noteFor(c, id, ConcurrencyStartBlocked) == nil {
			t.Errorf("%s has no start_blocked note: %+v", id, c.Notes)
		}
	}

	// Both start-ready, linked only through a done Task in between.
	index = newConcurrencyIndex(
		planned("T-001", "a", []string{}, []string{"a.go"}),
		concurrencyTask{"T-002", "c", ColumnDone, []string{}, []string{}, []string{"T-001"}},
		planned("T-003", "b", []string{}, []string{"b.go"}, "T-002"),
	)
	index.Tasks["T-002"].Evidence = &Evidence{CheckWaiver: &CheckWaiver{Task: "T-002"}}
	c = concurrencyFor(index)
	pair := pairFor(c, "T-001", "T-003")
	if pair == nil || pair.Together || pair.Reason.Kind != ConcurrencyDependencyPath || !strings.Contains(pair.Reason.Detail, "T-003 depends on T-001") {
		t.Errorf("pair = %+v, want a dependency path reason", pair)
	}
	wantGroups(t, c, nil)
}

func TestResolveConcurrencyV2_prerequisitesUseOrdinaryStartGate(t *testing.T) {
	cases := map[string]struct {
		prepare   func(*V2Index)
		wantReady bool
	}{
		"clear": {func(i *V2Index) { mustCheck(i, "C-001", "T-001", CheckResultClear) }, true},
		"waived": {func(i *V2Index) {
			i.Tasks["T-001"].Evidence = &Evidence{CheckWaiver: &CheckWaiver{Task: "T-001"}}
		}, true},
		"no clearance": {func(*V2Index) {}, false},
		"needs work":   {func(i *V2Index) { mustCheck(i, "C-001", "T-001", CheckResultNeedsWork) }, false},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			index := newConcurrencyIndex(
				concurrencyTask{"T-001", "c", ColumnDone, []string{}, []string{}, nil},
				planned("T-002", "a", []string{}, []string{"a.go"}, "T-001"),
				planned("T-003", "b", []string{}, []string{"b.go"}),
			)
			tc.prepare(index)
			c := concurrencyFor(index)
			if ready := len(c.Groups) == 1; ready != tc.wantReady {
				t.Errorf("groups = %v, want ready = %t", c.Groups, tc.wantReady)
			}
			if !tc.wantReady && noteFor(c, "T-002", ConcurrencyStartBlocked) == nil {
				t.Errorf("blocked Task is not explained: %+v", c.Notes)
			}
		})
	}

	// A waiver never satisfies requires: accepted.
	index := newConcurrencyIndex(
		concurrencyTask{"T-001", "c", ColumnDone, []string{}, []string{}, nil},
		planned("T-002", "a", []string{}, []string{"a.go"}),
	)
	index.Tasks["T-001"].Evidence = &Evidence{CheckWaiver: &CheckWaiver{Task: "T-001"}}
	index.Tasks["T-002"].DependsOn = []TaskDependencyV2{{Task: "T-001", Requires: TaskDependencyAccepted}}
	if c := concurrencyFor(index); len(c.Candidates) != 0 || noteFor(c, "T-002", ConcurrencyStartBlocked) == nil {
		t.Errorf("waived prerequisite satisfied requires: accepted: %+v", c)
	}
}

func TestResolveConcurrencyV2_objectiveDependencyBlocksEveryTask(t *testing.T) {
	index := newConcurrencyIndex(
		planned("T-001", "a", []string{}, []string{"a.go"}),
		planned("T-002", "b", []string{}, []string{"b.go"}),
	)
	index.Objectives["O-000"] = &ObjectiveV2{ID: "O-000", Status: ColumnPlanned}
	index.Objectives["O-001"].DependsOn = []string{"O-000"}
	c := concurrencyFor(index)
	if len(c.Candidates) != 0 || c.Groups != nil {
		t.Errorf("candidates = %v, groups = %v, want none", c.Candidates, c.Groups)
	}
	for _, id := range []string{"T-001", "T-002"} {
		if n := noteFor(c, id, ConcurrencyStartBlocked); n == nil || !strings.Contains(n.Detail, "O-000") {
			t.Errorf("%s note = %+v, want the Objective dependency", id, n)
		}
	}
}

func TestResolveConcurrencyV2_sharedWriteAndOverlapAreWithheld(t *testing.T) {
	cases := map[string]struct {
		a, b []string // reads, writes pairs flattened below
		kind ConcurrencyKindV2
		a1   concurrencyTask
		b1   concurrencyTask
	}{
		"shared write": {kind: ConcurrencySharedWrite,
			a1: planned("T-001", "a", []string{}, []string{"x.go"}), b1: planned("T-002", "b", []string{}, []string{"x.go"})},
		"write then read": {kind: ConcurrencyWriteReadOverlap,
			a1: planned("T-001", "a", []string{}, []string{"x.go"}), b1: planned("T-002", "b", []string{"x.go"}, []string{"b.go"})},
		"read then write": {kind: ConcurrencyWriteReadOverlap,
			a1: planned("T-001", "a", []string{"x.go"}, []string{"a.go"}), b1: planned("T-002", "b", []string{}, []string{"x.go"})},
		"case alias": {kind: ConcurrencyPathAlias,
			a1: planned("T-001", "a", []string{}, []string{"Dir/X.go"}), b1: planned("T-002", "b", []string{"dir/x.go"}, []string{"b.go"})},
		"unknown reads": {kind: ConcurrencyUnknownScope,
			a1: planned("T-001", "a", []string{}, []string{"a.go"}), b1: planned("T-002", "b", nil, []string{"b.go"})},
		"unknown writes": {kind: ConcurrencyUnknownScope,
			a1: planned("T-001", "a", []string{}, []string{"a.go"}), b1: planned("T-002", "b", []string{}, nil)},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			c := concurrencyFor(newConcurrencyIndex(tc.a1, tc.b1))
			wantGroups(t, c, nil)
			if name[:7] == "unknown" {
				// A Task with no scope is not a candidate at all.
				if n := noteFor(c, "T-002", ConcurrencyUnknownScope); n == nil {
					t.Errorf("notes = %+v, want unknown_scope for T-002", c.Notes)
				}
				return
			}
			pair := pairFor(c, "T-001", "T-002")
			if pair == nil || pair.Together || pair.Reason == nil || pair.Reason.Kind != tc.kind {
				t.Errorf("pair = %+v, want withheld as %s", pair, tc.kind)
			}
		})
	}
}

func TestResolveConcurrencyV2_independenceExplanationsAffectOnlyTheirOverlap(t *testing.T) {
	index := newConcurrencyIndex(
		planned("T-001", "a", []string{"two.go"}, []string{"one.go", "two2.go"}),
		planned("T-002", "b", []string{"one.go", "two2.go"}, []string{"b.go"}),
	)
	index.PlanIndependence["O-001"] = []PlanIndependenceV2{{
		Tasks:    [2]string{"T-001", "T-002"},
		Overlaps: []PlanOverlapV2{{Writer: "T-001", Reader: "T-002", Path: "one.go"}},
		Reason:   "stable copy",
	}}
	c := concurrencyFor(index)
	pair := pairFor(c, "T-001", "T-002")
	if pair == nil || pair.Together || !strings.Contains(pair.Reason.Detail, "two2.go") {
		t.Fatalf("pair = %+v, want the unexplained two2.go overlap withheld", pair)
	}

	index.Tasks["T-002"].Plan.Reads = scope([]string{"one.go"})
	wantGroups(t, concurrencyFor(index), [][]string{{"T-001", "T-002"}})

	// An explanation never covers a shared write.
	index.Tasks["T-002"].Plan.Writes = scope([]string{"one.go"})
	c = concurrencyFor(index)
	if pair := pairFor(c, "T-001", "T-002"); pair == nil || pair.Together || pair.Reason.Kind != ConcurrencySharedWrite {
		t.Errorf("pair = %+v, want shared_write", pair)
	}
}

func TestResolveConcurrencyV2_staleExplanationWithholdsAndSaysSo(t *testing.T) {
	tasks := map[string]string{
		"T-001": "lane: core\nplanned_reads: []\nplanned_writes: [shared.go]\n",
		"T-002": "lane: core\nplanned_reads: [shared.go]\nplanned_writes: []\n",
	}
	overlap := "{writer: T-001, reader: T-002, path: shared.go}"
	probe, err := LoadV2Index(writePlanProject(t, "", tasks))
	if err != nil {
		t.Fatal(err)
	}
	digest := PlanReviewDigest("T-001", probe.Tasks["T-001"].Plan, "T-002", probe.Tasks["T-002"].Plan)

	// Two Tasks in one lane are never suggested together, so declare a
	// second lane by extending the Objective's lanes list.
	laneTwo := map[string]string{"T-001": tasks["T-001"], "T-002": strings.Replace(tasks["T-002"], "lane: core", "lane: other", 1)}
	for name, reviewed := range map[string]string{"current": digest, "stale": "sha256:0"} {
		extra := "  - {key: other, title: Other}\n" + independenceYAML(reviewed, overlap, "reads a stable copy")
		root := writePlanProject(t, extra, laneTwo)
		index, err := LoadV2Index(root)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		pair := pairFor(ResolveConcurrencyV2(index, "O-001", ConcurrencyOptionsV2{Enabled: true}), "T-001", "T-002")
		switch name {
		case "current":
			if pair == nil || !pair.Together {
				t.Errorf("current explanation: pair = %+v, want together", pair)
			}
		case "stale":
			if pair == nil || pair.Together || pair.Reason.Kind != ConcurrencyWriteReadOverlap || !strings.Contains(pair.Reason.Detail, "stale or unusable") {
				t.Errorf("stale explanation: pair = %+v, want a withheld overlap that mentions it", pair)
			}
		}
	}
}

func TestResolveConcurrencyV2_oneTaskPerLane(t *testing.T) {
	index := newConcurrencyIndex(
		planned("T-001", "a", []string{}, []string{"a.go"}),
		planned("T-002", "a", []string{}, []string{"a2.go"}),
		planned("T-003", "b", []string{}, []string{"b.go"}),
	)
	c := concurrencyFor(index)
	if got := candidateIDs(c); !reflect.DeepEqual(got, []string{"T-001", "T-003"}) {
		t.Errorf("candidates = %v, want one per lane", got)
	}
	if n := noteFor(c, "T-002", ConcurrencySameLane); n == nil || n.Tasks[1] != "T-001" {
		t.Errorf("note = %+v, want same_lane naming T-001", n)
	}
	wantGroups(t, c, [][]string{{"T-001", "T-003"}})
}

func TestResolveConcurrencyV2_remainingLaneScopeConflictsWithhold(t *testing.T) {
	// T-001 and T-003 are disjoint, but T-002 later in lane a writes b.go.
	index := newConcurrencyIndex(
		planned("T-001", "a", []string{}, []string{"a.go"}),
		planned("T-002", "a", []string{}, []string{"b.go"}, "T-001"),
		planned("T-003", "b", []string{}, []string{"b.go"}),
	)
	c := concurrencyFor(index)
	pair := pairFor(c, "T-001", "T-003")
	if pair == nil || pair.Together || pair.Reason.Kind != ConcurrencySharedWrite || !strings.Contains(pair.Reason.Detail, "later work") {
		t.Errorf("pair = %+v, want a later-work shared_write reason", pair)
	}
	wantGroups(t, c, nil)

	// A later Task of the other lane without a manifest leaves safety unknown.
	index = newConcurrencyIndex(
		planned("T-001", "a", []string{}, []string{"a.go"}),
		planned("T-002", "b", []string{}, []string{"b.go"}),
		planned("T-003", "b", nil, nil, "T-002"),
	)
	c = concurrencyFor(index)
	if pair := pairFor(c, "T-001", "T-002"); pair == nil || pair.Together || pair.Reason.Kind != ConcurrencyUnknownScope {
		t.Errorf("pair = %+v, want unknown_scope", pair)
	}
}

func TestResolveConcurrencyV2_recordedActiveWork(t *testing.T) {
	active := func(lane string, reads, writes []string) concurrencyTask {
		return concurrencyTask{"T-009", lane, ColumnInProgress, reads, writes, nil}
	}
	cases := map[string]struct {
		active      concurrencyTask
		wantKind    ConcurrencyKindV2 // of T-001's note; empty when T-001 stays a candidate
		wantActive  bool
		wantNoteFor string
	}{
		"disjoint active work":      {active: active("b", []string{}, []string{"z.go"}), wantActive: true},
		"active shared write":       {active: active("b", []string{}, []string{"a.go"}), wantKind: ConcurrencySharedWrite},
		"active unknown scope":      {active: active("b", nil, nil), wantKind: ConcurrencyUnknownScope},
		"active without lane":       {active: active("", []string{"a.go"}, []string{"z.go"}), wantKind: ConcurrencyWriteReadOverlap},
		"active owns the same lane": {active: active("a", []string{}, []string{"z.go"}), wantKind: ConcurrencyLaneBusy},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			index := newConcurrencyIndex(tc.active, planned("T-001", "a", []string{}, []string{"a.go"}))
			c := concurrencyFor(index)
			if !reflect.DeepEqual(c.Active, []string{"T-009"}) {
				t.Errorf("Active = %v", c.Active)
			}
			if tc.wantKind == "" {
				if got := candidateIDs(c); !reflect.DeepEqual(got, []string{"T-001"}) {
					t.Errorf("candidates = %v, want T-001 alongside active work", got)
				}
				wantGroups(t, c, nil) // alongside active work is not a claim about two candidates
				return
			}
			if len(c.Candidates) != 0 {
				t.Errorf("candidates = %v, want none", c.Candidates)
			}
			if n := noteFor(c, "T-001", tc.wantKind); n == nil {
				t.Errorf("notes = %+v, want %s for T-001", c.Notes, tc.wantKind)
			}
		})
	}
}

func TestResolveConcurrencyV2_incompatibleTriangleNeverGroupsAConflictingPair(t *testing.T) {
	// a-b and b-c are fine, a-c share a write: no group may hold both a and c.
	index := newConcurrencyIndex(
		planned("T-001", "a", []string{}, []string{"a.go", "ac.go"}),
		planned("T-002", "b", []string{}, []string{"b.go"}),
		planned("T-003", "c", []string{}, []string{"c.go", "ac.go"}),
	)
	c := concurrencyFor(index)
	wantGroups(t, c, [][]string{{"T-001", "T-002"}})
	if pair := pairFor(c, "T-001", "T-003"); pair == nil || pair.Together {
		t.Errorf("T-001/T-003 = %+v, want withheld", pair)
	}
	if pair := pairFor(c, "T-002", "T-003"); pair == nil || !pair.Together {
		t.Errorf("T-002/T-003 = %+v, want together as a pair", pair)
	}
	for _, group := range c.Groups {
		if len(group) < 2 {
			t.Errorf("group %v claims parallelism with fewer than two Tasks", group)
		}
	}
}

func TestResolveConcurrencyV2_singletonAndEmptyClaimNothing(t *testing.T) {
	c := concurrencyFor(newConcurrencyIndex(planned("T-001", "a", []string{}, []string{"a.go"})))
	if len(c.Candidates) != 1 || c.Groups != nil || c.Pairs != nil {
		t.Errorf("single candidate = %+v, want no pairs or groups", c)
	}
	c = concurrencyFor(newConcurrencyIndex())
	if c.Candidates != nil || c.Groups != nil || c.Withheld != nil {
		t.Errorf("empty Objective = %+v, want nothing", c)
	}
}

func TestResolveConcurrencyV2_ungroupableTasksAreExplained(t *testing.T) {
	index := newConcurrencyIndex(
		planned("T-001", "a", []string{}, []string{"a.go"}),
		planned("T-002", "", []string{}, []string{"b.go"}),
		planned("T-003", "b", []string{}, []string{"c.go"}),
	)
	index.UnusablePlans["T-003"] = true
	c := concurrencyFor(index)
	if noteFor(c, "T-002", ConcurrencyNoLane) == nil || noteFor(c, "T-003", ConcurrencyUnusablePlan) == nil {
		t.Errorf("notes = %+v, want no_lane and unusable_plan", c.Notes)
	}
	if len(c.Groups) != 0 {
		t.Errorf("groups = %v, want none", c.Groups)
	}
}

func TestResolveConcurrencyV2_replanAndStaleSelectionWithholdLaunchAdvice(t *testing.T) {
	index := newConcurrencyIndex(
		planned("T-001", "a", []string{}, []string{"a.go"}),
		planned("T-002", "b", []string{}, []string{"b.go"}),
	)
	wantGroups(t, concurrencyFor(index), [][]string{{"T-001", "T-002"}})

	stale := ResolveConcurrencyV2(index, "O-001", ConcurrencyOptionsV2{Enabled: true, Selection: &SelectionDiagnostic{Kind: SelectionDone}})
	if stale.Withheld == nil || stale.Withheld.Kind != ConcurrencySelectionStale || stale.Groups != nil || stale.Candidates != nil {
		t.Errorf("stale selection = %+v, want advice withheld", stale)
	}
	if len(stale.Lanes) != 3 {
		t.Error("withheld advice dropped lane membership")
	}

	index.Tasks["T-002"].Evidence = &Evidence{Replan: &Replan{Reason: "interface changed"}}
	replan := concurrencyFor(index)
	if replan.Withheld == nil || replan.Withheld.Kind != ConcurrencyReplanRecorded || replan.Withheld.Tasks[0] != "T-002" || replan.Groups != nil {
		t.Errorf("replan = %+v, want advice withheld naming T-002", replan.Withheld)
	}
}

func TestResolveNext_concurrencyIsAdviceOnly(t *testing.T) {
	build := func() (*V2Index, *RouterStateV2) {
		index := newConcurrencyIndex(
			planned("T-001", "a", []string{}, []string{"a.go"}),
			planned("T-002", "b", []string{}, []string{"b.go"}),
		)
		router := &RouterStateV2{State: RouterPhaseTask, Objective: "O-001", Task: "T-001"}
		return index, router
	}
	strip := func(n Next) Next { n.Concurrency = nil; return n }

	var baseline Next
	for i, enabled := range []bool{false, true} {
		index, router := build()
		next := resolveNextWithTestGoal(NextInput{Index: index, Router: router, ParallelPlanning: enabled})
		if next.Concurrency == nil || next.Concurrency.Objective != "O-001" || next.Concurrency.Enabled != enabled {
			t.Fatalf("enabled=%t: Concurrency = %+v", enabled, next.Concurrency)
		}
		if enabled && len(next.Concurrency.Groups) != 1 {
			t.Errorf("enabled Groups = %v, want one", next.Concurrency.Groups)
		}
		if !enabled && next.Concurrency.Groups != nil {
			t.Errorf("disabled Groups = %v, want none", next.Concurrency.Groups)
		}
		if i == 0 {
			baseline = strip(next)
		} else if !reflect.DeepEqual(strip(next), baseline) {
			t.Errorf("Next changed with the preference:\n on: %+v\noff: %+v", strip(next), baseline)
		}
	}

	// A lane conflict changes advice, never the selected Task or its gate.
	index, router := build()
	index.Tasks["T-002"].Plan.Writes = scope([]string{"a.go"})
	next := resolveNextWithTestGoal(NextInput{Index: index, Router: router, ParallelPlanning: true})
	if next.Kind != NextExecute || next.Task == nil || next.Task.ID != "T-001" || next.GateDecision == nil || !next.GateDecision.Allowed {
		t.Errorf("Next = %+v, want T-001 still allowed to execute", next)
	}
	if len(next.Concurrency.Groups) != 0 {
		t.Errorf("Groups = %v, want the shared write withheld", next.Concurrency.Groups)
	}

	// No resolved Objective, no projection.
	index, _ = build()
	none := resolveNextWithTestGoal(NextInput{Index: index, Router: &RouterStateV2{}, ParallelPlanning: true})
	if none.Concurrency != nil {
		t.Errorf("Concurrency = %+v, want nil without an Objective", none.Concurrency)
	}
}

func TestResolveConcurrencyV2_isDeterministic(t *testing.T) {
	index := newConcurrencyIndex(
		planned("T-001", "a", []string{"s.go"}, []string{"a.go"}),
		planned("T-002", "b", []string{"s.go"}, []string{"b.go"}),
		planned("T-003", "c", []string{}, []string{"a.go"}),
		planned("T-004", "a", []string{}, []string{"d.go"}),
	)
	first := concurrencyFor(index)
	for range 20 {
		if got := concurrencyFor(index); !reflect.DeepEqual(got, first) {
			t.Fatalf("projection varies between runs:\n%+v\n%+v", got, first)
		}
	}
}
