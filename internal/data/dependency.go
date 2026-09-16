package data

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

type DependencyKind string

const (
	DependencyMissing DependencyKind = ""
	DependencyTask    DependencyKind = "task"
	DependencyEpic    DependencyKind = "epic"
)

type DependencyResolution struct {
	Kind       DependencyKind
	ID         string
	TaskStatus ColumnType
	EpicStatus string
}

func ResolveDependency(ref string, dependent Task, tasks []Task, epicStatuses map[string]string) DependencyResolution {
	for _, task := range tasks {
		if task.ID == ref && sameRelease(dependent.Release, task.Release) {
			return DependencyResolution{Kind: DependencyTask, ID: task.ID, TaskStatus: task.Column}
		}
	}

	if isShortTaskRef(ref) {
		shortRef := taskShortID(ref)
		for _, task := range tasks {
			if taskShortID(task.ID) == shortRef && sameRelease(dependent.Release, task.Release) && sameEpic(dependent.Epic, task.Epic) {
				return DependencyResolution{Kind: DependencyTask, ID: task.ID, TaskStatus: task.Column}
			}
		}
	}

	if status, ok := epicStatuses[ref]; ok {
		return DependencyResolution{Kind: DependencyEpic, ID: ref, EpicStatus: status}
	}

	if isShortEpicRef(ref) {
		for epicID, status := range epicStatuses {
			if epicShortID(epicID) == ref {
				return DependencyResolution{Kind: DependencyEpic, ID: epicID, EpicStatus: status}
			}
		}
	}

	return DependencyResolution{}
}

func sameRelease(a, b string) bool {
	return a == "" || b == "" || a == b
}

func sameEpic(a, b string) bool {
	return a == "" || b == "" || a == b
}

func isShortTaskRef(ref string) bool {
	if strings.Contains(ref, "/") {
		return false
	}
	short := taskShortID(ref)
	return len(short) == 4 && short[0] == 'T' && allDigits(short[1:])
}

func isShortEpicRef(ref string) bool {
	return len(ref) == 3 && ref[0] == 'E' && allDigits(ref[1:])
}

func allDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return s != ""
}

func taskShortID(id string) string {
	if slash := strings.LastIndexByte(id, '/'); slash >= 0 {
		id = id[slash+1:]
	}
	if dash := strings.IndexByte(id, '-'); dash >= 0 {
		return id[:dash]
	}
	return id
}

func epicShortID(id string) string {
	if dash := strings.IndexByte(id, '-'); dash >= 0 {
		return id[:dash]
	}
	return id
}

// DependencyBlockKind names why a V2 Task dependency is unsatisfied.
type DependencyBlockKind string

const (
	// DependencyBlockNotDone means the dependency Task's status is not done.
	DependencyBlockNotDone DependencyBlockKind = "not_done"
	// DependencyBlockNotCleared means the dependency Task is done but its
	// clearance is not current; Clearance names the observed state.
	DependencyBlockNotCleared DependencyBlockKind = "not_cleared"
	// DependencyBlockNotAccepted means the dependency Task is done and
	// currently cleared but requires accepted, and the owner has not
	// accepted that current Check.
	DependencyBlockNotAccepted DependencyBlockKind = "not_accepted"
)

// DependencyBlock is the typed reason one V2 Task dependency is unsatisfied.
// It names the target and the unmet requirement rather than a bare boolean,
// reusing ClearanceState so a done-but-not-cleared block names the same
// clearance vocabulary ResolveClearance reports.
type DependencyBlock struct {
	Target    string
	Requires  TaskDependencyRequirement
	Kind      DependencyBlockKind
	Clearance ClearanceState // set for NotCleared and NotAccepted; empty for NotDone
}

// DependencyDecision is the resolved satisfaction for one V2 Task
// dependency.
type DependencyDecision struct {
	Satisfied bool
	Block     *DependencyBlock // nil when Satisfied is true
}

// ResolveTaskDependencyV2 resolves whether dep is satisfied, re-resolving
// the dependency Task's status, clearance, and owner acceptance from index
// on every call rather than trusting any cached judgement. A requires: clear
// dependency is satisfied only when its target Task is done with current
// clearance; requires: accepted additionally needs recorded owner acceptance
// of that Task's current Check — acceptance of a superseded Check does not
// satisfy it.
func ResolveTaskDependencyV2(index *V2Index, dep TaskDependencyV2) DependencyDecision {
	target, ok := index.Tasks[dep.Task]
	if !ok || target.Status != ColumnDone {
		return DependencyDecision{Block: &DependencyBlock{Target: dep.Task, Requires: dep.Requires, Kind: DependencyBlockNotDone}}
	}

	clearance := ResolveClearance(index, dep.Task)
	if clearance.State != ClearanceCurrent {
		return DependencyDecision{Block: &DependencyBlock{Target: dep.Task, Requires: dep.Requires, Kind: DependencyBlockNotCleared, Clearance: clearance.State}}
	}

	if dep.Requires == TaskDependencyAccepted {
		if !ownerAcceptedCheck(target.Evidence, clearance.Check) {
			return DependencyDecision{Block: &DependencyBlock{Target: dep.Task, Requires: dep.Requires, Kind: DependencyBlockNotAccepted, Clearance: clearance.State}}
		}
	}

	return DependencyDecision{Satisfied: true}
}

// ValidateV2ReferenceGraphs validates the Task and Objective dependency
// graphs of index by global ID. It reports the first self-dependency,
// missing dependency target, or reference cycle found while walking IDs in
// sorted order, so a given project reports the same diagnostic on every run.
// The Task graph and Objective graph are validated independently: a Task
// cycle is never conflated with an Objective cycle.
func ValidateV2ReferenceGraphs(index *V2Index) error {
	if err := validateTaskDependencyGraph(index); err != nil {
		return err
	}
	return validateObjectiveDependencyGraph(index)
}

func validateTaskDependencyGraph(index *V2Index) error {
	ids := slices.Sorted(maps.Keys(index.Tasks))
	edges := make(map[string][]string, len(ids))

	for _, id := range ids {
		task := index.Tasks[id]
		for _, dep := range task.DependsOn {
			if dep.Task == id {
				return fmt.Errorf("%w: %s: task %s depends on itself", ErrV2SelfDependency, task.Source.Path, id)
			}
			if _, ok := index.Tasks[dep.Task]; !ok {
				return fmt.Errorf("%w: %s: task %s depends on missing task %s", ErrV2MissingDependencyTarget, task.Source.Path, id, dep.Task)
			}
			edges[id] = append(edges[id], dep.Task)
		}
	}

	if cycle := findV2Cycle(ids, edges); cycle != nil {
		return fmt.Errorf("%w: task cycle %s", ErrV2DependencyCycle, describeV2Cycle(cycle, func(id string) string { return index.Tasks[id].Source.Path }))
	}
	return nil
}

func validateObjectiveDependencyGraph(index *V2Index) error {
	ids := slices.Sorted(maps.Keys(index.Objectives))
	edges := make(map[string][]string, len(ids))

	for _, id := range ids {
		objective := index.Objectives[id]
		for _, dep := range objective.DependsOn {
			if dep == id {
				return fmt.Errorf("%w: %s: objective %s depends on itself", ErrV2SelfDependency, objective.Source.Path, id)
			}
			if _, ok := index.Objectives[dep]; !ok {
				return fmt.Errorf("%w: %s: objective %s depends on missing objective %s", ErrV2MissingDependencyTarget, objective.Source.Path, id, dep)
			}
			edges[id] = append(edges[id], dep)
		}
	}

	if cycle := findV2Cycle(ids, edges); cycle != nil {
		return fmt.Errorf("%w: objective cycle %s", ErrV2DependencyCycle, describeV2Cycle(cycle, func(id string) string { return index.Objectives[id].Source.Path }))
	}
	return nil
}

// findV2Cycle runs an iterative DFS over edges, visiting order first so two
// runs over the same graph always report the same cycle. It returns the
// cycle as an ID sequence ending back at its own start, or nil if the graph
// is acyclic.
func findV2Cycle(order []string, edges map[string][]string) []string {
	const (
		white = iota
		gray
		black
	)
	color := make(map[string]int, len(order))
	var stack []string

	var visit func(id string) []string
	visit = func(id string) []string {
		color[id] = gray
		stack = append(stack, id)
		for _, next := range edges[id] {
			switch color[next] {
			case white:
				if cycle := visit(next); cycle != nil {
					return cycle
				}
			case gray:
				start := slices.Index(stack, next)
				cycle := append([]string{}, stack[start:]...)
				return append(cycle, next)
			}
		}
		stack = stack[:len(stack)-1]
		color[id] = black
		return nil
	}

	for _, id := range order {
		if color[id] == white {
			if cycle := visit(id); cycle != nil {
				return cycle
			}
		}
	}
	return nil
}

func describeV2Cycle(cycle []string, pathFor func(string) string) string {
	parts := make([]string, len(cycle))
	for i, id := range cycle {
		parts[i] = fmt.Sprintf("%s (%s)", id, pathFor(id))
	}
	return strings.Join(parts, " -> ")
}
