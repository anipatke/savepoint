package v2

import (
	"maps"
	"slices"
	"strings"

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
	DetailRelease   DetailKind = "RELEASE"
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

// ReleaseObjectiveProgress is the resolved progress of one Objective derived
// into a Release. The Objective still owns its Tasks; this value only carries
// the counts and completion decision needed to explain Release readiness.
type ReleaseObjectiveProgress struct {
	Objective  RecordRef
	TasksDone  int
	TasksTotal int
	Decision   data.GateDecision
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
	// MemberObjectives and the fields below are populated only for a Release
	// detail. Membership is derived from each Objective's release reference;
	// the detail carries the resolved values so the renderer never consults the
	// index.
	MemberObjectives  []ReleaseObjectiveProgress
	Outcome           string
	Why               string
	SuccessConditions string
	Boundaries        string
	ReleaseDecision   *data.GateDecision
	LegacyCompletion  *data.LegacyCompletionReference

	Clearance data.Clearance
	// ByException is true only for an Objective whose latest completion
	// decision is allowed despite clearance not being current — the same
	// resolver call ObjectiveRow.ByException uses, so the badge this detail's
	// CLEARANCE section shows (via objectiveCheckBadge) always agrees with the
	// sidebar row for the same Objective. Always false for a Task or Release
	// detail: a Task's exception is its own separate badge, not folded into
	// its Check badge, so this field carries no meaning there.
	ByException bool
	Evidence    *data.Evidence
	Checks      []CheckEntry
	StyleReview *StyleReview
	Issues      []*data.IssueV2
}

// StyleReview carries only the latest Check's authored style section. Present
// distinguishes a missing section from one whose content is empty.
type StyleReview struct {
	CheckID string
	Lines   []string
	Present bool
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
		Kind:        DetailTask,
		ID:          task.ID,
		Title:       task.Title,
		Status:      task.Status,
		Stage:       task.Stage,
		Body:        task.Source.Body,
		Owner:       &owner,
		Clearance:   data.ResolveClearance(index, task.ID),
		Evidence:    task.Evidence,
		Checks:      checkHistory(index, task.ID),
		StyleReview: latestStyleReview(index, task.ID),
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
		Kind:        DetailObjective,
		ID:          objective.ID,
		Title:       objective.Title,
		Status:      objective.Status,
		Body:        objective.Source.Body,
		Clearance:   data.ResolveClearance(index, objective.ID),
		ByException: data.ResolveObjectiveCompletion(index, objectiveID).AllowedByException,
		Evidence:    objective.Evidence,
		Checks:      checkHistory(index, objective.ID),
		StyleReview: latestStyleReview(index, objective.ID),
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

// newReleaseDetail resolves one Release and all of the evidence the overlay
// promises: its authored delivery promise, derived member Objective progress,
// canonical completion decision, Check history, and linked Issues. The
// renderer receives this finished value and does not reinterpret the index.
func newReleaseDetail(index *data.V2Index, releaseID string) (RecordDetail, bool) {
	if index == nil {
		return RecordDetail{}, false
	}
	release, ok := index.Releases[releaseID]
	if !ok {
		return RecordDetail{}, false
	}

	decision := data.ResolveReleaseCompletion(index, release.ID)
	detail := RecordDetail{
		Kind:              DetailRelease,
		ID:                release.ID,
		Title:             release.Title,
		Status:            release.Status,
		Body:              release.Source.Body,
		Outcome:           release.Outcome,
		Why:               release.Why,
		SuccessConditions: release.SuccessConditions,
		Boundaries:        release.Boundaries,
		ReleaseDecision:   &decision,
		LegacyCompletion:  release.LegacyCompletion,
		Clearance:         data.ResolveClearance(index, release.ID),
		Evidence:          release.Evidence,
		Checks:            checkHistory(index, release.ID),
		StyleReview:       latestStyleReview(index, release.ID),
	}

	objectiveIDs := slices.Clone(index.ReleaseObjectives[release.ID])
	slices.Sort(objectiveIDs)
	for _, objectiveID := range objectiveIDs {
		if _, ok := index.Objectives[objectiveID]; !ok {
			continue
		}
		progress := ReleaseObjectiveProgress{
			Objective:  objectiveRef(index, objectiveID),
			TasksTotal: len(index.ObjectiveTasks[objectiveID]),
			Decision:   data.ResolveObjectiveCompletion(index, objectiveID),
		}
		for _, taskID := range index.ObjectiveTasks[objectiveID] {
			if task, ok := index.Tasks[taskID]; ok && task.Status == data.ColumnDone {
				progress.TasksDone++
			}
		}
		detail.MemberObjectives = append(detail.MemberObjectives, progress)
	}
	detail.Issues = linkedIssues(index, nil, detail.Checks)

	return detail, true
}

func latestStyleReview(index *data.V2Index, targetID string) *StyleReview {
	id := index.LatestCheck[targetID]
	if id == "" {
		return nil
	}
	check := index.Checks[id]
	lines, present := codeStyleReviewLines(check.Source.Body)
	return &StyleReview{CheckID: id, Lines: lines, Present: present}
}

// codeStyleReviewLines keeps the authored lines between the exact section
// heading and the next level-two heading, trimming only blank edge lines.
func codeStyleReviewLines(body string) ([]string, bool) {
	body = strings.ReplaceAll(body, "\r\n", "\n")
	lines := strings.Split(body, "\n")
	start := -1
	for i, line := range lines {
		if line == "## Code Style Review" {
			start = i + 1
			break
		}
	}
	if start < 0 {
		return nil, false
	}
	end := len(lines)
	for i := start; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "## ") {
			end = i
			break
		}
	}
	for start < end && strings.TrimSpace(lines[start]) == "" {
		start++
	}
	for end > start && strings.TrimSpace(lines[end-1]) == "" {
		end--
	}
	return lines[start:end], true
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
	if detail.Kind == DetailRelease {
		return newReleaseDetail(index, detail.ID)
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
