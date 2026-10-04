package data

import (
	"fmt"
	"slices"
	"strings"
)

// This file projects the advisory parallel-planning metadata of one Objective
// into stable lane membership and conservative suggestions about which lanes
// could start together. It is advice only: no gate, Next selection, start,
// advance or completion decision consumes it, and the owner can ignore every
// suggestion. The projection is pure. It reads the loaded index and the
// saved preference, probes no filesystem or worktree, and never infers that
// a done Task has merged.

// ConcurrencyKindV2 names why a Task or pair of Tasks is not suggested to run
// together, or why all advice for an Objective is withheld.
type ConcurrencyKindV2 string

const (
	// ConcurrencyDependencyPath: one Task depends on the other, directly or
	// through other Tasks.
	ConcurrencyDependencyPath ConcurrencyKindV2 = "dependency_path"
	// ConcurrencySharedWrite: both Tasks plan to write the same file.
	ConcurrencySharedWrite ConcurrencyKindV2 = "shared_write"
	// ConcurrencyWriteReadOverlap: one Task writes a file the other reads and
	// no valid independence explanation covers it.
	ConcurrencyWriteReadOverlap ConcurrencyKindV2 = "write_read_overlap"
	// ConcurrencyUnknownScope: a Task has no usable reviewed read/write scope.
	ConcurrencyUnknownScope ConcurrencyKindV2 = "unknown_scope"
	// ConcurrencyPathAlias: paths differing only by case may be one file, and
	// lexical comparison cannot certify them as distinct.
	ConcurrencyPathAlias ConcurrencyKindV2 = "path_alias"
	// ConcurrencyStartBlocked: the ordinary start gate does not allow the
	// Task yet; Detail repeats its blockers.
	ConcurrencyStartBlocked ConcurrencyKindV2 = "start_blocked"
	// ConcurrencyNoLane: the Task has no saved lane.
	ConcurrencyNoLane ConcurrencyKindV2 = "no_lane"
	// ConcurrencyUnusablePlan: the Task's planning metadata is unusable.
	ConcurrencyUnusablePlan ConcurrencyKindV2 = "unusable_plan"
	// ConcurrencySameLane: another Task is suggested first in the lane.
	ConcurrencySameLane ConcurrencyKindV2 = "same_lane"
	// ConcurrencyLaneBusy: a Task in the lane is already in progress.
	ConcurrencyLaneBusy ConcurrencyKindV2 = "lane_busy"
	// ConcurrencyReplanRecorded: a Task of the Objective records a replan, so
	// the plan is under revision and no launch is suggested.
	ConcurrencyReplanRecorded ConcurrencyKindV2 = "replan_recorded"
	// ConcurrencySelectionStale: the router selection is missing or out of
	// date, so no launch is suggested.
	ConcurrencySelectionStale ConcurrencyKindV2 = "selection_stale"
)

// ConcurrencyReasonV2 explains one withheld suggestion in plain words.
// Tasks names the Tasks the reason is about, the first being the Task that is
// not suggested.
type ConcurrencyReasonV2 struct {
	Kind   ConcurrencyKindV2
	Tasks  []string
	Detail string
}

// ConcurrencyLaneV2 is one declared lane's stable membership, derived from
// each Task's lane reference. Membership never depends on the preference or on
// whether the plan is usable for advice, so a board can keep headings stable
// as Tasks move between columns.
type ConcurrencyLaneV2 struct {
	Objective  string
	Key        string
	Title      string
	Planned    []string
	InProgress []string
	Done       []string
}

// Ref returns the Objective-namespaced lane identity, such as O-033/core.
func (l ConcurrencyLaneV2) Ref() string { return l.Objective + "/" + l.Key }

// Empty reports whether no Task belongs to the lane.
func (l ConcurrencyLaneV2) Empty() bool {
	return len(l.Planned)+len(l.InProgress)+len(l.Done) == 0
}

// ConcurrencyCandidateV2 is the one Task suggested next in a lane.
type ConcurrencyCandidateV2 struct {
	Task string
	Lane string // lane key
}

// ConcurrencyPairV2 compares two candidates. Reason is nil when they are
// suggested together.
type ConcurrencyPairV2 struct {
	Tasks    [2]string
	Together bool
	Reason   *ConcurrencyReasonV2
}

// ConcurrencyV2 is the projection for one Objective.
type ConcurrencyV2 struct {
	Objective string
	// Enabled is the saved parallel-planning preference the projection was
	// computed under. When false only Lanes is filled.
	Enabled bool
	Lanes   []ConcurrencyLaneV2
	// Withheld is set when all advice for the Objective is withheld.
	Withheld *ConcurrencyReasonV2
	// Diagnostics are the Objective's unusable planning metadata.
	Diagnostics []PlanDiagnosticV2
	// Active lists the Objective's in-progress Tasks, which count as running
	// work. Every Candidate is compatible with all of them.
	Active []string
	// Candidates holds at most one start-ready Task per lane, in lane order.
	// A candidate alone claims parallelism only alongside Active Tasks.
	Candidates []ConcurrencyCandidateV2
	// Pairs compares every two candidates.
	Pairs []ConcurrencyPairV2
	// Groups are sets of two or more candidates in which every pair is
	// together. They come from a deterministic first-fit pass, so a candidate
	// joins a group only when compatible with every member; they are not an
	// optimal schedule.
	Groups [][]string
	// Notes explain each planned Task that is not a candidate.
	Notes []ConcurrencyReasonV2
}

// ConcurrencyOptionsV2 carries the inputs of ResolveConcurrencyV2 besides the
// index.
type ConcurrencyOptionsV2 struct {
	// Enabled is the saved parallel-planning preference.
	Enabled bool
	// Selection is the router selection diagnostic when the caller has one.
	// Any diagnostic withholds launch suggestions; a nil value means no
	// selection problem, as for a board's focused Objective.
	Selection *SelectionDiagnostic
}

// ResolveConcurrencyV2 projects objectiveID's lanes and, when enabled,
// suggestions for which lanes may start together. It returns nil for an
// unknown Objective. It reuses ResolveTaskStart for ordinary start
// eligibility and adds no lane condition to any gate.
func ResolveConcurrencyV2(index *V2Index, objectiveID string, opts ConcurrencyOptionsV2) *ConcurrencyV2 {
	if index == nil {
		return nil
	}
	objective, ok := index.Objectives[objectiveID]
	if !ok {
		return nil
	}
	var tasks []*TaskV2
	for _, id := range index.ObjectiveTasks[objectiveID] {
		if task, ok := index.Tasks[id]; ok {
			tasks = append(tasks, task)
		}
	}

	c := &ConcurrencyV2{Objective: objectiveID, Enabled: opts.Enabled, Lanes: concurrencyLanes(objective, tasks)}
	if !opts.Enabled {
		return c
	}

	owned := map[string]bool{objectiveID: true}
	for _, task := range tasks {
		owned[task.ID] = true
	}
	for _, d := range index.PlanDiagnostics {
		if owned[d.Record] {
			c.Diagnostics = append(c.Diagnostics, d)
		}
	}

	if opts.Selection != nil {
		c.Withheld = &ConcurrencyReasonV2{Kind: ConcurrencySelectionStale, Detail: fmt.Sprintf("the router selection is not current (%s), so no launch is suggested", opts.Selection.Kind)}
		return c
	}
	for _, task := range tasks {
		if task.Status != ColumnDone && task.Evidence != nil && task.Evidence.Replan != nil {
			c.Withheld = &ConcurrencyReasonV2{Kind: ConcurrencyReplanRecorded, Tasks: []string{task.ID}, Detail: fmt.Sprintf("%s records a replan, so the plan is being revised and no launch is suggested", task.ID)}
			return c
		}
	}

	p := &concurrencyPass{index: index, objective: objective, c: c, reach: map[string]map[string]bool{}}
	p.run(tasks)
	return c
}

func concurrencyLanes(objective *ObjectiveV2, tasks []*TaskV2) []ConcurrencyLaneV2 {
	lanes := make([]ConcurrencyLaneV2, len(objective.Plan.Lanes))
	for i, lane := range objective.Plan.Lanes {
		lanes[i] = ConcurrencyLaneV2{Objective: objective.ID, Key: lane.Key, Title: lane.Title}
		for _, task := range tasks {
			if task.Plan.Lane != lane.Key {
				continue
			}
			switch task.Status {
			case ColumnPlanned:
				lanes[i].Planned = append(lanes[i].Planned, task.ID)
			case ColumnInProgress:
				lanes[i].InProgress = append(lanes[i].InProgress, task.ID)
			case ColumnDone:
				lanes[i].Done = append(lanes[i].Done, task.ID)
			}
		}
	}
	return lanes
}

// concurrencySide is one candidate or active Task together with the
// unfinished work of its lane, which is treated as scope that will run later.
type concurrencySide struct {
	self  *TaskV2
	tasks []*TaskV2
}

type concurrencyPass struct {
	index     *V2Index
	objective *ObjectiveV2
	c         *ConcurrencyV2
	reach     map[string]map[string]bool
}

func (p *concurrencyPass) run(tasks []*TaskV2) {
	c := p.c
	var active, planned []*TaskV2
	busy := map[string]bool{}
	for _, task := range tasks {
		switch task.Status {
		case ColumnInProgress:
			active = append(active, task)
			c.Active = append(c.Active, task.ID)
			if task.Plan.Lane != "" {
				busy[task.Plan.Lane] = true
			}
		case ColumnPlanned:
			planned = append(planned, task)
		}
	}

	heads := map[string]*TaskV2{} // lane key -> suggested Task
	for _, task := range planned {
		if decision := ResolveTaskStart(p.index, task.ID); !decision.Allowed {
			p.note(ConcurrencyStartBlocked, []string{task.ID}, "%s cannot start yet: %s", task.ID, blockerDetails(decision.Blockers))
			continue
		}
		lane := task.Plan.Lane
		switch {
		case p.index.UnusablePlans[task.ID]:
			p.note(ConcurrencyUnusablePlan, []string{task.ID}, "%s has unusable planning metadata; see the planning diagnostics", task.ID)
		case lane == "":
			p.note(ConcurrencyNoLane, []string{task.ID}, "%s has no saved lane, so no lane suggestion is made", task.ID)
		case busy[lane]:
			p.note(ConcurrencyLaneBusy, []string{task.ID}, "%s is in lane %s, which already has work in progress", task.ID, lane)
		case heads[lane] != nil:
			p.note(ConcurrencySameLane, []string{task.ID, heads[lane].ID}, "%s is suggested first in lane %s; one Task at a time per lane", heads[lane].ID, lane)
		case scopeUnknown(p.index, task):
			p.note(ConcurrencyUnknownScope, []string{task.ID}, "%s has no reviewed read/write scope, so safety is unknown", task.ID)
		default:
			heads[lane] = task
		}
	}

	var candidates []ConcurrencyCandidateV2
	var sides []concurrencySide
	for _, lane := range p.objective.Plan.Lanes {
		head := heads[lane.Key]
		if head == nil {
			continue
		}
		side := p.laneSide(head)
		if reason := p.againstActive(side, active); reason != nil {
			p.c.Notes = append(p.c.Notes, *reason)
			continue
		}
		candidates = append(candidates, ConcurrencyCandidateV2{Task: head.ID, Lane: lane.Key})
		sides = append(sides, side)
	}
	c.Candidates = candidates

	together := map[[2]string]bool{}
	for i := range sides {
		for j := i + 1; j < len(sides); j++ {
			pair := ConcurrencyPairV2{Tasks: [2]string{sides[i].self.ID, sides[j].self.ID}}
			pair.Reason = p.sidesConflict(sides[i], sides[j])
			pair.Together = pair.Reason == nil
			together[pair.Tasks] = pair.Together
			c.Pairs = append(c.Pairs, pair)
		}
	}
	c.Groups = firstFitGroups(candidates, together)

	slices.SortStableFunc(c.Notes, func(a, b ConcurrencyReasonV2) int {
		return strings.Compare(a.Tasks[0], b.Tasks[0])
	})
}

func (p *concurrencyPass) note(kind ConcurrencyKindV2, tasks []string, format string, args ...any) {
	p.c.Notes = append(p.c.Notes, ConcurrencyReasonV2{Kind: kind, Tasks: tasks, Detail: fmt.Sprintf(format, args...)})
}

// laneSide gathers the unfinished work of self's lane; a Task without a lane
// stands alone.
func (p *concurrencyPass) laneSide(self *TaskV2) concurrencySide {
	side := concurrencySide{self: self, tasks: []*TaskV2{self}}
	if self.Plan.Lane == "" {
		return side
	}
	side.tasks = nil
	for _, id := range p.index.ObjectiveTasks[self.Objective] {
		task, ok := p.index.Tasks[id]
		if ok && task.Plan.Lane == self.Plan.Lane && task.Status != ColumnDone {
			side.tasks = append(side.tasks, task)
		}
	}
	return side
}

// againstActive returns why candidate cannot start alongside the recorded
// active work, or nil.
func (p *concurrencyPass) againstActive(candidate concurrencySide, active []*TaskV2) *ConcurrencyReasonV2 {
	for _, task := range active {
		if reason := p.sidesConflict(candidate, p.laneSide(task)); reason != nil {
			reason.Tasks = slices.Insert(slices.DeleteFunc(slices.Clone(reason.Tasks), func(id string) bool { return id == candidate.self.ID }), 0, candidate.self.ID)
			reason.Detail = fmt.Sprintf("%s is not suggested alongside %s, which is in progress: %s", candidate.self.ID, task.ID, reason.Detail)
			return reason
		}
	}
	return nil
}

// sidesConflict returns the first reason x and y are not suggested together.
// The two chosen Tasks must be unrelated by any dependency path; for the rest
// of each lane a dependency path only orders the work, so it is no conflict.
func (p *concurrencyPass) sidesConflict(x, y concurrencySide) *ConcurrencyReasonV2 {
	for _, u := range x.tasks {
		for _, v := range y.tasks {
			strict := u == x.self && v == y.self
			reason := p.pairConflict(u, v, strict)
			if reason == nil {
				continue
			}
			if !strict {
				reason.Detail += " (later work in the lanes)"
			}
			return reason
		}
	}
	return nil
}

func (p *concurrencyPass) pairConflict(u, v *TaskV2, strict bool) *ConcurrencyReasonV2 {
	pair := []string{u.ID, v.ID}
	if dependent, dependency, linked := p.linked(u.ID, v.ID); linked {
		if !strict {
			return nil
		}
		return &ConcurrencyReasonV2{Kind: ConcurrencyDependencyPath, Tasks: pair, Detail: fmt.Sprintf("%s depends on %s, directly or through other Tasks", dependent, dependency)}
	}
	for _, task := range []*TaskV2{u, v} {
		if scopeUnknown(p.index, task) {
			return &ConcurrencyReasonV2{Kind: ConcurrencyUnknownScope, Tasks: pair, Detail: fmt.Sprintf("%s has no reviewed read/write scope, so safety is unknown", task.ID)}
		}
	}

	var shared, overlap, alias *ConcurrencyReasonV2
	record := func(slot **ConcurrencyReasonV2, kind ConcurrencyKindV2, format string, args ...any) {
		if *slot == nil {
			*slot = &ConcurrencyReasonV2{Kind: kind, Tasks: pair, Detail: fmt.Sprintf(format, args...)}
		}
	}
	for _, a := range u.Plan.Writes.paths {
		for _, b := range v.Plan.Writes.paths {
			switch ComparePlanPaths(a, b) {
			case PlanPathsSame:
				record(&shared, ConcurrencySharedWrite, "%s and %s both plan to write %s", u.ID, v.ID, a)
			case PlanPathsMayAlias:
				record(&alias, ConcurrencyPathAlias, "%s and %s plan to write %s and %s, which may be one file", u.ID, v.ID, a, b)
			}
		}
	}
	for _, order := range [][2]*TaskV2{{u, v}, {v, u}} {
		writer, reader := order[0], order[1]
		for _, w := range writer.Plan.Writes.paths {
			for _, r := range reader.Plan.Reads.paths {
				switch ComparePlanPaths(w, r) {
				case PlanPathsSame:
					if !p.explained(writer.ID, reader.ID, w) {
						record(&overlap, ConcurrencyWriteReadOverlap, "%s writes %s, which %s reads%s", writer.ID, w, reader.ID, p.staleHint(u.ID, v.ID))
					}
				case PlanPathsMayAlias:
					record(&alias, ConcurrencyPathAlias, "%s writes %s and %s reads %s, which may be one file", writer.ID, w, reader.ID, r)
				}
			}
		}
	}
	for _, reason := range []*ConcurrencyReasonV2{shared, overlap, alias} {
		if reason != nil {
			return reason
		}
	}
	return nil
}

// explained reports whether a valid independence explanation covers the exact
// write/read overlap.
func (p *concurrencyPass) explained(writer, reader, path string) bool {
	pair := sortedPair(writer, reader)
	for _, e := range p.index.PlanIndependence[p.objective.ID] {
		if e.Tasks != pair {
			continue
		}
		for _, o := range e.Overlaps {
			if o.Writer == writer && o.Reader == reader && o.Path == path {
				return true
			}
		}
	}
	return false
}

// staleHint notes an authored explanation for the pair that was not usable.
func (p *concurrencyPass) staleHint(a, b string) string {
	pair := sortedPair(a, b)
	for _, e := range p.objective.Plan.Independent {
		if e.Tasks == pair {
			return "; an independence explanation exists but is stale or unusable"
		}
	}
	return ""
}

func sortedPair(a, b string) [2]string {
	if a > b {
		a, b = b, a
	}
	return [2]string{a, b}
}

// linked reports whether one Task depends on the other through any chain of
// dependencies, naming the dependent and the dependency.
func (p *concurrencyPass) linked(a, b string) (dependent, dependency string, ok bool) {
	switch {
	case p.dependsOn(a)[b]:
		return a, b, true
	case p.dependsOn(b)[a]:
		return b, a, true
	}
	return "", "", false
}

// dependsOn returns every Task id reachable from id through depends_on.
func (p *concurrencyPass) dependsOn(id string) map[string]bool {
	if seen, ok := p.reach[id]; ok {
		return seen
	}
	seen := map[string]bool{}
	p.reach[id] = seen
	stack := []string{id}
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		task, ok := p.index.Tasks[current]
		if !ok {
			continue
		}
		for _, dep := range task.DependsOn {
			if !seen[dep.Task] {
				seen[dep.Task] = true
				stack = append(stack, dep.Task)
			}
		}
	}
	return seen
}

func scopeUnknown(index *V2Index, task *TaskV2) bool {
	return index.UnusablePlans[task.ID] || !task.Plan.Reads.declared || !task.Plan.Writes.declared
}

func blockerDetails(blockers []GateBlocker) string {
	details := make([]string, len(blockers))
	for i, blocker := range blockers {
		details[i] = blocker.Detail
	}
	return strings.Join(details, "; ")
}

// firstFitGroups places each candidate, in order, into the first group whose
// members it is together with, and keeps only groups of two or more.
func firstFitGroups(candidates []ConcurrencyCandidateV2, together map[[2]string]bool) [][]string {
	var all [][]string
	for _, candidate := range candidates {
		placed := false
		for i, group := range all {
			if !slices.ContainsFunc(group, func(member string) bool { return !together[[2]string{member, candidate.Task}] }) {
				all[i] = append(group, candidate.Task)
				placed = true
				break
			}
		}
		if !placed {
			all = append(all, []string{candidate.Task})
		}
	}
	var groups [][]string
	for _, group := range all {
		if len(group) > 1 {
			groups = append(groups, group)
		}
	}
	return groups
}
