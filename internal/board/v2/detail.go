package v2

import (
	"maps"
	"slices"

	"github.com/opencode/savepoint/internal/data"
)

// This file resolves what an open detail shows, and it is the only place the
// detail reaches the index. Everything a RecordDetail carries is a record field
// or a resolver's own value, settled once when the overlay opens and again on
// every reload, so detail_view.go formats a finished value and decides nothing.
//
// The record body is carried as the author wrote it. It is displayed and
// scrolled; nothing here reads it for meaning, and nothing in this surface
// writes.

// DetailKind names which record family an open detail is about. Its value is
// also the overlay's heading, so the two cannot disagree.
type DetailKind string

const (
	DetailTask      DetailKind = "TASK"
	DetailObjective DetailKind = "OBJECTIVE"
)

// RecordRef is one record named from a detail — an owning Objective, an owned
// Task, or a dependency target — carried by identity, title, and recorded
// status so the overlay never looks one up while rendering.
type RecordRef struct {
	ID     string
	Title  string
	Status data.ColumnType
}

// DependencyEntry is one declared Task dependency: the level the declaration
// requires, and the decision ResolveTaskDependencyV2 reached about it. The
// satisfaction is the resolver's answer — the board never reads the dependency
// Task's own status or evidence to work one out.
type DependencyEntry struct {
	Target   RecordRef
	Requires data.TaskDependencyRequirement
	Decision data.DependencyDecision
}

// ObjectiveDependencyEntry is the same for one declared Objective dependency,
// resolved through ResolveObjectiveDependency.
type ObjectiveDependencyEntry struct {
	Target   RecordRef
	Decision data.ObjectiveDependencyDecision
}

// RecordDetail is everything one open overlay shows about one record: its
// identity and recorded lifecycle, what it depends on and what depends on it,
// the resolved clearance and the recorded evidence behind it, the Check chain,
// the follow-ups the index links to it, and its own body.
//
// Task-only and Objective-only members are left empty for the other family
// rather than split into two types: the evidence half is identical for both,
// and one type keeps it rendered identically.
type RecordDetail struct {
	Kind   DetailKind
	ID     string
	Title  string
	Status data.ColumnType
	// Stage is a Task's recorded implementation stage, empty for a planned or
	// done Task and for every Objective.
	Stage data.ProgressStage
	// Body is the record's author-owned markdown, carried verbatim.
	Body string

	// Owner is the Objective a Task belongs to, nil on an Objective detail.
	Owner *RecordRef
	// OwnedTasks are the Tasks an Objective owns, read from
	// index.ObjectiveTasks and empty on a Task detail.
	OwnedTasks []RecordRef

	// Dependencies are a Task's own declared dependencies;
	// ObjectiveDependencies are an Objective's. A record carries one list or
	// the other, in the order its record declares them.
	Dependencies          []DependencyEntry
	ObjectiveDependencies []ObjectiveDependencyEntry

	Clearance data.Clearance
	Evidence  *data.Evidence
	Checks    []CheckEntry
	Issues    []*data.IssueV2
}

// newTaskDetail resolves the detail for one Task. It returns ok=false for an ID
// the index does not hold, which is what a reload that removed the open record
// reports.
func newTaskDetail(index *data.V2Index, taskID string) (RecordDetail, bool) {
	task, ok := index.Tasks[taskID]
	if !ok {
		return RecordDetail{}, false
	}

	owner := objectiveRef(index, task.Objective)
	detail := RecordDetail{
		Kind:      DetailTask,
		ID:        task.ID,
		Title:     task.Title,
		Status:    task.Status,
		Stage:     task.Stage,
		Body:      task.Source.Body,
		Owner:     &owner,
		Clearance: data.ResolveClearance(index, task.ID),
		Evidence:  task.Evidence,
		Checks:    checkHistory(index, task.ID),
	}

	for _, dependency := range task.DependsOn {
		detail.Dependencies = append(detail.Dependencies, DependencyEntry{
			Target:   taskRef(index, dependency.Task),
			Requires: dependency.Requires,
			Decision: data.ResolveTaskDependencyV2(index, dependency),
		})
	}
	detail.Issues = linkedIssues(index, index.TaskIssues[task.ID], detail.Checks)

	return detail, true
}

// newObjectiveDetail resolves the detail for one Objective: the same evidence
// shape a Task carries, plus the Tasks it owns and the Objectives it waits on.
func newObjectiveDetail(index *data.V2Index, objectiveID string) (RecordDetail, bool) {
	objective, ok := index.Objectives[objectiveID]
	if !ok {
		return RecordDetail{}, false
	}

	detail := RecordDetail{
		Kind:      DetailObjective,
		ID:        objective.ID,
		Title:     objective.Title,
		Status:    objective.Status,
		Body:      objective.Source.Body,
		Clearance: data.ResolveClearance(index, objective.ID),
		Evidence:  objective.Evidence,
		Checks:    checkHistory(index, objective.ID),
	}

	// Ownership is index.ObjectiveTasks' answer, built from each Task's own
	// objective field — never the directory a Task file happens to sit in.
	for _, taskID := range index.ObjectiveTasks[objectiveID] {
		detail.OwnedTasks = append(detail.OwnedTasks, taskRef(index, taskID))
	}
	for _, dependencyID := range objective.DependsOn {
		detail.ObjectiveDependencies = append(detail.ObjectiveDependencies, ObjectiveDependencyEntry{
			Target:   objectiveRef(index, dependencyID),
			Decision: data.ResolveObjectiveDependency(index, dependencyID),
		})
	}
	detail.Issues = linkedIssues(index, nil, detail.Checks)

	return detail, true
}

// reopenDetail re-resolves an open detail against a freshly loaded index, so a
// reload under an open overlay shows the records as they now are. A record the
// load no longer holds reports ok=false: there is nothing left to show, and the
// copy already on screen would be older than the project.
func reopenDetail(index *data.V2Index, detail RecordDetail) (RecordDetail, bool) {
	if index == nil {
		return RecordDetail{}, false
	}
	if detail.Kind == DetailObjective {
		return newObjectiveDetail(index, detail.ID)
	}
	return newTaskDetail(index, detail.ID)
}

// objectiveRef and taskRef name a referenced record by identity, title, and
// recorded status. A reference the index cannot resolve degrades to its ID
// rather than to a guess — LoadV2Index refuses a project with a dangling
// reference, so this is a guard rather than a state a loaded board reaches.
func objectiveRef(index *data.V2Index, id string) RecordRef {
	if objective, ok := index.Objectives[id]; ok {
		return RecordRef{ID: id, Title: objective.Title, Status: objective.Status}
	}
	return RecordRef{ID: id}
}

func taskRef(index *data.V2Index, id string) RecordRef {
	if task, ok := index.Tasks[id]; ok {
		return RecordRef{ID: id, Title: task.Title, Status: task.Status}
	}
	return RecordRef{ID: id}
}

// linkedIssues collects the follow-ups the index's own link maps hang off this
// record: the Issues naming the Task as carrying their repair, and the Issues
// naming any Check in the record's history. Both maps were built by
// LoadV2Index; nothing here walks the Issues to work out which apply.
//
// The result is sorted by ID and lists each Issue once, so an Issue named by
// two Checks in the same chain appears once rather than twice.
func linkedIssues(index *data.V2Index, taskIssueIDs []string, checks []CheckEntry) []*data.IssueV2 {
	linked := make(map[string]bool, len(taskIssueIDs))
	for _, id := range taskIssueIDs {
		linked[id] = true
	}
	for _, entry := range checks {
		for _, id := range index.CheckIssues[entry.Check.ID] {
			linked[id] = true
		}
	}
	if len(linked) == 0 {
		return nil
	}

	issues := make([]*data.IssueV2, 0, len(linked))
	for _, id := range slices.Sorted(maps.Keys(linked)) {
		issues = append(issues, index.Issues[id])
	}
	return issues
}
