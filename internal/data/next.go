package data

import (
	"maps"
	"slices"
)

// SelectionRecordKind names which record family a selection diagnostic is
// about: the router's Release, Objective, Task, or Issue.
type SelectionRecordKind string

const (
	SelectionRecordRelease   SelectionRecordKind = "release"
	SelectionRecordObjective SelectionRecordKind = "objective"
	SelectionRecordTask      SelectionRecordKind = "task"
	SelectionRecordIssue     SelectionRecordKind = "issue"
)

// SelectionDiagnosticKind names why the router's contextual selection is
// unresolved or stale. Release diagnostics are separate from the older Task /
// Objective diagnostics so a renderer can explain an archived Release,
// unassigned Objective, and cross-Release mismatch without inspecting the
// index again.
type SelectionDiagnosticKind string

const (
	// SelectionNotFound means a router-named Objective, Task, or Issue ID does not
	// exist among the project's live records. Archived records are out of
	// scope for lookup, so this is also what a migrated-away ID reports.
	SelectionNotFound SelectionDiagnosticKind = "not_found"
	// SelectionMismatch means the router-named Objective and Task both
	// resolve, but the Task's own objective field names a different
	// Objective. Ownership always comes from the Task record, never the
	// router, so this is reported rather than silently reconciled.
	SelectionMismatch SelectionDiagnosticKind = "mismatch"
	// SelectionDone means the exact selection resolved, but its selected
	// Task is done, its selected Objective is done with no Task selected, or
	// its selected Issue is resolved.
	// The selection remains intact so Next continues to describe that record.
	SelectionDone SelectionDiagnosticKind = "done"
	// SelectionReleaseMissing means the router has no Goal selected. No
	// Objective, Task, or Issue selection is resolved until a Goal is chosen.
	SelectionReleaseMissing SelectionDiagnosticKind = "release_missing"
	// SelectionReleaseNotFound means the router names a Release absent from
	// the live index. No similarly named Release is substituted.
	SelectionReleaseNotFound SelectionDiagnosticKind = "release_not_found"
	// SelectionReleaseArchived means the selected Release is a historical
	// migrated record and is not an actionable live delivery context.
	SelectionReleaseArchived SelectionDiagnosticKind = "release_archived"
	// SelectionObjectiveUnassigned means the selected Objective has no
	// Release reference, so it cannot be interpreted inside a Release.
	SelectionObjectiveUnassigned SelectionDiagnosticKind = "objective_unassigned"
	// SelectionReleaseMismatch means the selected Objective names a different
	// Release than the router context.
	SelectionReleaseMismatch SelectionDiagnosticKind = "release_mismatch"

	// Short aliases keep the diagnostic vocabulary easy to discover for
	// callers that reason about selection state rather than record family.
	SelectionArchived   = SelectionReleaseArchived
	SelectionUnassigned = SelectionObjectiveUnassigned
)

// SelectionDiagnostic is the typed, branchable result of an unresolved or
// stale selection. A caller inspects Kind and the ID fields below it; nothing
// here is a formatted message standing in for the IDs read.
type SelectionDiagnostic struct {
	Kind SelectionDiagnosticKind

	// RecordKind and ID are populated when Kind is SelectionNotFound: which
	// record family the router named, and the ID absent from the live
	// index.
	RecordKind SelectionRecordKind
	ID         string

	// Release, Objective, and ObjectiveRelease are populated for Release
	// context diagnostics. ObjectiveRelease is empty for an unassigned
	// Objective.
	Release          string
	Objective        string
	ObjectiveRelease string

	// RouterObjective, Task, and TaskObjective are populated when Kind is
	// SelectionMismatch: the Objective the router named, the Task that
	// otherwise resolved, and the Objective that Task's own record declares
	// as its owner.
	RouterObjective string
	Task            string
	TaskObjective   string
}

// Selection is one honest resolved outcome: no context selected, a Release
// with no Objective selected, an Objective with no Task selected under it, a
// matched Release/Objective/Task context, or a selected Issue. When a Release
// diagnostic names a live Release, the Release is retained for rendering
// while Objective/Task action remains unresolved.
type Selection struct {
	Release   *ReleaseV2
	Objective *ObjectiveV2
	Task      *TaskV2
	Issue     *IssueV2
}

// ResolveSelection turns the router's Objective/Task/Issue selections into
// exactly matched records against index, or a typed SelectionDiagnostic.
//
// Matching is by exact global ID only: a map lookup against the matching
// record map, with no prefix, suffix, numeric-proximity, title, or
// path-based fallback anywhere in this function. A router-named ID that does
// not resolve is reported as itself; a similarly numbered record is never
// substituted for it.
//
// ResolveSelection reads only its two arguments: no filesystem access, no
// discovery, and no write. It performs no lookup against an archive — the
// live index is the only thing consulted, so a migrated-away ID reports as
// not found rather than as a guess at where it went.
func ResolveSelection(index *V2Index, router *RouterStateV2) (Selection, *SelectionDiagnostic) {
	if index == nil {
		index = &V2Index{}
	}
	if router == nil {
		return Selection{}, nil
	}
	if router.Release == "" {
		return Selection{}, &SelectionDiagnostic{
			Kind: SelectionReleaseMissing, RecordKind: SelectionRecordRelease,
		}
	}

	var issue *IssueV2
	var issueDiagnostic *SelectionDiagnostic
	if router.Issue != "" {
		var ok bool
		issue, ok = index.Issues[router.Issue]
		if !ok {
			issueDiagnostic = &SelectionDiagnostic{
				Kind: SelectionNotFound, RecordKind: SelectionRecordIssue,
				ID: router.Issue,
			}
		} else if issue.Status == IssueStatusResolved {
			issueDiagnostic = &SelectionDiagnostic{
				Kind: SelectionDone, RecordKind: SelectionRecordIssue,
				ID: issue.ID,
			}
		}
	}

	var release *ReleaseV2
	if router.Release != "" {
		candidate, ok := index.Releases[router.Release]
		if !ok {
			return Selection{Issue: issue}, &SelectionDiagnostic{
				Kind: SelectionReleaseNotFound, RecordKind: SelectionRecordRelease,
				ID: router.Release, Release: router.Release,
			}
		}
		release = candidate
		if release.LegacyCompletion != nil {
			return Selection{Release: release, Issue: issue}, &SelectionDiagnostic{
				Kind: SelectionReleaseArchived, RecordKind: SelectionRecordRelease,
				ID: router.Release, Release: router.Release,
			}
		}
	}

	if router.Objective == "" {
		// ReadStateV2 already refuses a Task selected with no Objective
		// (ErrV2InvalidOwnership), so an empty Objective here means no
		// selection at all, regardless of router phase.
		return Selection{Release: release, Issue: issue}, issueDiagnostic
	}

	objective, ok := index.Objectives[router.Objective]
	if !ok {
		return Selection{Release: release, Issue: issue}, &SelectionDiagnostic{
			Kind: SelectionNotFound, RecordKind: SelectionRecordObjective,
			ID: router.Objective, Release: router.Release, Objective: router.Objective,
		}
	}

	if release != nil {
		switch {
		case objective.Release == "":
			return Selection{Release: release}, &SelectionDiagnostic{
				Kind: SelectionObjectiveUnassigned, RecordKind: SelectionRecordObjective,
				ID: objective.ID, Release: release.ID, Objective: objective.ID,
			}
		case objective.Release != release.ID:
			return Selection{Release: release}, &SelectionDiagnostic{
				Kind: SelectionReleaseMismatch, RecordKind: SelectionRecordObjective,
				ID: objective.ID, Release: release.ID, Objective: objective.ID,
				ObjectiveRelease: string(objective.Release),
			}
		}
	}

	if router.Task == "" {
		selection := Selection{Release: release, Objective: objective, Issue: issue}
		if objective.Status == ColumnDone {
			return selection, &SelectionDiagnostic{
				Kind: SelectionDone, RecordKind: SelectionRecordObjective,
				ID: objective.ID, Release: router.Release, Objective: objective.ID,
			}
		}
		return selection, issueDiagnostic
	}

	task, ok := index.Tasks[router.Task]
	if !ok {
		return Selection{Release: release}, &SelectionDiagnostic{
			Kind: SelectionNotFound, RecordKind: SelectionRecordTask,
			ID: router.Task, Release: router.Release, Objective: router.Objective,
		}
	}

	if task.Objective != router.Objective {
		return Selection{Release: release}, &SelectionDiagnostic{
			Kind:            SelectionMismatch,
			Release:         router.Release,
			Objective:       router.Objective,
			RouterObjective: router.Objective,
			Task:            router.Task,
			TaskObjective:   task.Objective,
		}
	}

	selection := Selection{Release: release, Objective: objective, Task: task, Issue: issue}
	if task.Status == ColumnDone {
		return selection, &SelectionDiagnostic{
			Kind: SelectionDone, RecordKind: SelectionRecordTask,
			ID: task.ID, Release: router.Release, Objective: objective.ID,
		}
	}
	return selection, issueDiagnostic
}

// NextKind names the selected record's next action. ResolveNext never
// searches for a different Objective or Task when the router selection is
// absent or cannot be resolved.
// ResolveNext: owner validation sits below check-needed because asking the
// owner to accept work with no current technical clearance would invert the
// authority model E43 established.
type NextKind string

const (
	// NextReplan means the selected Task carries a recorded replan flag.
	NextReplan NextKind = "replan"
	// NextDependency means the selected Task is blocked on an unsatisfied
	// Task dependency or an unsatisfied Objective dependency of its owning
	// Objective.
	NextDependency NextKind = "dependency"
	// NextExecute means the selected Task may start, advance, or complete
	// now — including completion allowed only by a recorded exception,
	// which GateDecision.AllowedByException and .Exception name rather than
	// this rung claiming a CLEAR result.
	NextExecute NextKind = "execute"
	// NextCheckNeeded means the selected Task or its owning Objective needs
	// a fresh Check: its clearance is missing, needs_work, stale, unknown,
	// or current without independent checker provenance.
	NextCheckNeeded NextKind = "check_needed"
	// NextOwnerValidationRequired means clearance is current but the owner
	// has not accepted it, and owner_validation.required is set.
	NextOwnerValidationRequired NextKind = "owner_validation_required"
	// NextObjectiveIntegration means every Task the Objective owns is done,
	// but the Objective's own integration clearance is not current (or is
	// current without owner acceptance where required).
	NextObjectiveIntegration NextKind = "objective_integration"
	// NextObjectiveReady means all owned Tasks are done and the Objective's
	// existing completion resolver allows the owner to record it as done.
	NextObjectiveReady NextKind = "objective_ready"
	// NextReleaseReady means the Release completion decision is allowed and
	// the next action is for the owner to record the Release as done.
	NextReleaseReady NextKind = "release_ready"
	// NextSelectTask means the selected Objective owns unfinished Tasks but
	// the router has not selected which one to work on.
	NextSelectTask NextKind = "select_task"
	// NextNothingSelected means there is no resolved Objective or Task
	// selection. It never implies that the project has no available work.
	NextNothingSelected NextKind = "nothing_selected"
	// NextPlanObjective means the selected Objective owns no Tasks yet, so
	// its next action is to plan the work under that Objective.
	NextPlanObjective NextKind = "plan_objective"
	// NextIssue means the router selected an Issue without an Objective or
	// Task, so the Issue itself is the next record to work on.
	NextIssue NextKind = "issue"
)

// Next is the one derived answer ResolveNext returns for a project: the
// router-selected Objective, Task, or Issue, its next action kind, and any
// resolver-derived gate evidence. When a Task is selected, Objective also
// carries its owning record when that record exists. GateDecision and
// Clearance are populated only when the selected action reports them; each
// value comes from the existing E43/E44 resolvers (STYLE-07, DATA-02).
type Next struct {
	Kind      NextKind
	Release   *ReleaseV2
	Objective *ObjectiveV2
	Task      *TaskV2
	Issue     *IssueV2
	// ObjectivesWithoutGoal is the sorted index fact copied into the shared
	// Next projection so renderers can flag unassigned Objectives without
	// consulting project data independently.
	ObjectivesWithoutGoal []string

	// GateDecision is the decision the rung was derived from. For a Release
	// rung it is ResolveReleaseCompletion; for older rungs it remains the
	// Task- or Objective-level decision described below.
	GateDecision *GateDecision
	// Clearance is the resolved clearance the rung reports on, for
	// NextCheckNeeded, NextOwnerValidationRequired, NextObjectiveIntegration,
	// and the Release-specific evidence rungs.
	Clearance *Clearance

	// SelectionDiagnostic is set whenever the router's selection did not
	// resolve or points at finished work. Unresolved selections yield
	// NextNothingSelected; a SelectionDone diagnostic leaves the resolved
	// selection in place so its existing Next rung remains visible.
	SelectionDiagnostic *SelectionDiagnostic

	// Issues are the Issues relevant to Task (when set) or otherwise
	// Objective, resolved here so a rendering surface never has to consult
	// the index on its own to answer what follow-up hangs off the selected
	// record (STYLE-07). Nil when neither is set, or when none are linked.
	Issues []*IssueV2
}

// NextInput carries everything ResolveNext reads: the project's index and the
// decoded router hint. ResolveNext performs no discovery, no filesystem
// access, and no write; it consults only these two values.
type NextInput struct {
	Index  *V2Index
	Router *RouterStateV2
}

// ResolveNext computes the next action for the router's exact Objective,
// Task, or Issue selection, using the existing resolvers for Objective/Task
// gate state. When no record is selected, or a selection cannot be resolved, it does
// not substitute work found elsewhere. A valid selected Release retains its
// completion rungs once all of its member Objectives are done. ResolveNext
// never returns an error. A nil Index or Router is read as an empty project
// and an empty selection rather than panicking; see the boundary note below.
func ResolveNext(input NextInput) Next {
	index, router := input.Index, input.Router
	// A nil index or router is a caller whose load has not completed, or did
	// not succeed. Reading them as an empty project and an empty selection
	// keeps a half-assembled call from panicking in a long-running consumer
	// such as the board, which holds this input across reloads rather than
	// building it once per process the way a command does. It is a boundary
	// check, not an interpretation: reporting that a project failed to load
	// stays the caller's job, because only the caller holds the diagnostic.
	if index == nil {
		index = &V2Index{}
	}
	if router == nil {
		router = &RouterStateV2{}
	}

	selection, diagnostic := ResolveSelection(index, router)

	next := resolveLadder(index, selection, diagnostic)
	next.Issue = selection.Issue
	if next.Task != nil && next.Objective == nil {
		next.Objective = index.Objectives[next.Task.Objective]
	}
	next.SelectionDiagnostic = diagnostic
	next.ObjectivesWithoutGoal = index.ObjectivesWithoutGoal
	next.Issues = relevantIssues(index, next)
	return next
}

// relevantIssues collects the Issues linked to next's selected Task,
// Objective, or Release, walking only the link maps LoadV2Index already built rather
// than re-deriving any relationship. A Task's Issues are read straight from
// TaskIssues. An Objective with no Task selected has no such direct map,
// because IssueV2 carries no Objective reference of its own, so its
// relevant Issues are the union of every owned Task's Issues and any Issue
// linked to a Check scoped to the Objective itself. It returns nil rather
// than an empty slice when nothing applies, so a caller can render no
// section instead of an empty one.
func relevantIssues(index *V2Index, next Next) []*IssueV2 {
	var ids []string
	switch {
	case next.Task != nil:
		ids = index.TaskIssues[next.Task.ID]
	case next.Objective != nil:
		seen := make(map[string]bool)
		for _, taskID := range index.ObjectiveTasks[next.Objective.ID] {
			for _, issueID := range index.TaskIssues[taskID] {
				seen[issueID] = true
			}
		}
		for _, checkID := range index.ScopeChecks[next.Objective.ID] {
			for _, issueID := range index.CheckIssues[checkID] {
				seen[issueID] = true
			}
		}
		ids = slices.Sorted(maps.Keys(seen))
	case next.Release != nil:
		seen := make(map[string]bool)
		for _, checkID := range index.ScopeChecks[next.Release.ID] {
			for _, issueID := range index.CheckIssues[checkID] {
				seen[issueID] = true
			}
		}
		ids = slices.Sorted(maps.Keys(seen))
	default:
		return nil
	}

	if len(ids) == 0 {
		return nil
	}
	issues := make([]*IssueV2, 0, len(ids))
	for _, id := range ids {
		if issue, ok := index.Issues[id]; ok {
			issues = append(issues, issue)
		}
	}
	return issues
}

// resolveLadder resolves only the router's selection. A diagnostic or an
// empty Objective selection never falls through to a project-wide search.
func resolveLadder(index *V2Index, selection Selection, diagnostic *SelectionDiagnostic) Next {
	if diagnostic != nil && diagnostic.Kind != SelectionDone {
		validIssueContext := diagnostic.RecordKind == SelectionRecordIssue && (selection.Task != nil || selection.Objective != nil)
		issueWinsOverReleaseDiagnostic := selection.Issue != nil && selection.Task == nil && selection.Objective == nil && diagnostic.RecordKind == SelectionRecordRelease
		if !validIssueContext && !issueWinsOverReleaseDiagnostic {
			return Next{Kind: NextNothingSelected, Release: selection.Release}
		}
	}

	if selection.Task == nil && selection.Objective == nil && selection.Issue != nil {
		return Next{Kind: NextIssue, Release: selection.Release, Issue: selection.Issue}
	}

	if selection.Release != nil {
		return resolveReleaseLadder(index, selection)
	}

	if task := selection.Task; task != nil && task.Status != ColumnDone {
		return resolveTaskRung(index, task)
	}

	if selection.Objective != nil {
		return resolveSelectedObjective(index, selection.Objective)
	}

	return Next{Kind: NextNothingSelected}
}

// resolveReleaseLadder applies existing gate resolvers to the selected
// Release's explicitly selected Objective/Task. With no Objective selected,
// it reports whether the member Objectives allow the owner to complete the
// Release. No other member is selected by search.
func resolveReleaseLadder(index *V2Index, selection Selection) Next {
	release := selection.Release
	if task := selection.Task; task != nil && task.Status != ColumnDone {
		return withRelease(resolveTaskRung(index, task), release)
	}

	if selection.Objective != nil {
		return withRelease(resolveSelectedObjective(index, selection.Objective), release)
	}

	return resolveReleaseCompletionRung(index, release)
}

func withRelease(next Next, release *ReleaseV2) Next {
	next.Release = release
	return next
}

// resolveReleaseCompletionRung maps the Release completion decision to the
// ready action or an unselected state. The caller reaches it only when no
// Objective is selected.
func resolveReleaseCompletionRung(index *V2Index, release *ReleaseV2) Next {
	decision := ResolveReleaseCompletion(index, release.ID)
	next := Next{Release: release, GateDecision: &decision}

	if decision.Allowed {
		next.Kind = NextReleaseReady
		return next
	}

	next.Kind = NextNothingSelected
	return next
}

// resolveTaskRung reads the one gate decision that governs task's next move
// — ResolveTaskStart for a planned Task, ResolveTaskAdvance for one mid
// build/test, ResolveTaskCompletion for one at audit — and maps it onto a
// rung. It calls exactly one resolver per Task, so a Task audited at the
// same moment its owning Objective would also block never reports both.
func resolveTaskRung(index *V2Index, task *TaskV2) Next {
	switch {
	case task.Status == ColumnPlanned:
		decision := ResolveTaskStart(index, task.ID)
		return taskDecisionRung(task, decision)
	case task.Stage != StageAudit:
		decision := ResolveTaskAdvance(index, task.ID)
		return taskDecisionRung(task, decision)
	default:
		decision := ResolveTaskCompletion(index, task.ID)
		next := taskDecisionRung(task, decision)
		if next.Kind == NextCheckNeeded || next.Kind == NextOwnerValidationRequired {
			clearance := ResolveClearance(index, task.ID)
			next.Clearance = &clearance
		}
		return next
	}
}

// taskDecisionRung turns one GateDecision into a Next: allowed maps to
// NextExecute regardless of whether it was allowed by exception, and a
// blocked decision maps by its blockers' kinds via rungForBlockers.
func taskDecisionRung(task *TaskV2, decision GateDecision) Next {
	next := Next{Task: task, GateDecision: &decision}
	if decision.Allowed {
		next.Kind = NextExecute
		return next
	}
	next.Kind = rungForBlockers(decision.Blockers)
	return next
}

// rungForBlockers ranks a GateDecision's blockers onto the ladder's rungs
// two, three, five, and six, in that priority order. It never re-derives
// the resolver's own judgement — it only reads the typed Kind the resolver
// already assigned. Rung four's invalid-state blocker cannot occur here in
// practice, because the strict V2 Task decoder only ever admits a canonical
// status/stage combination and resolveTaskRung dispatches on that same
// combination; the dependency rung is the closest honest fallback for any
// blocker this function does not otherwise recognize.
func rungForBlockers(blockers []GateBlocker) NextKind {
	for _, blocker := range blockers {
		if blocker.Kind == GateBlockReplan {
			return NextReplan
		}
	}
	for _, blocker := range blockers {
		if blocker.Kind == GateBlockDependency || blocker.Kind == GateBlockObjectiveDependency {
			return NextDependency
		}
	}
	for _, blocker := range blockers {
		switch blocker.Kind {
		case GateBlockClearanceMissing, GateBlockClearanceNeedsWork, GateBlockClearanceStale, GateBlockClearanceUnknown, GateBlockCheckerAuthority:
			return NextCheckNeeded
		}
	}
	for _, blocker := range blockers {
		if blocker.Kind == GateBlockOwnerAcceptance {
			return NextOwnerValidationRequired
		}
	}
	return NextDependency
}

// resolveSelectedObjective returns the next step for the Objective the
// router selected. It does not choose which Task to work on: when owned work
// remains, the Objective stays selected until a Task is named explicitly.
func resolveSelectedObjective(index *V2Index, objective *ObjectiveV2) Next {
	taskIDs := index.ObjectiveTasks[objective.ID]
	if len(taskIDs) == 0 {
		return Next{Kind: NextPlanObjective, Objective: objective}
	}

	for _, taskID := range taskIDs {
		task, ok := index.Tasks[taskID]
		if !ok || task.Status != ColumnDone {
			return Next{Kind: NextSelectTask, Objective: objective}
		}
	}

	if next, ok := resolveObjectiveIntegrationRung(index, objective.ID); ok {
		return next
	}

	return Next{Kind: NextSelectTask, Objective: objective}
}

// resolveObjectiveIntegrationRung reads ResolveObjectiveCompletion for
// objectiveID after the selected Objective's Tasks are done. A blocked
// decision maps to the Objective integration action; an allowed decision
// means the owner may record the Objective as done.
func resolveObjectiveIntegrationRung(index *V2Index, objectiveID string) (Next, bool) {
	objective, ok := index.Objectives[objectiveID]
	if !ok {
		return Next{}, false
	}

	// An Objective that owns no Task has not been broken down yet, so
	// "every owned Task is done" is vacuously true and this rung would
	// ask for an integration Check over work that does not exist. The
	// real next action is planning, which rung eight answers.
	if len(index.ObjectiveTasks[objectiveID]) == 0 {
		return Next{}, false
	}

	decision := ResolveObjectiveCompletion(index, objectiveID)
	for _, blocker := range decision.Blockers {
		if blocker.Kind == GateBlockInvalidState {
			return Next{}, false
		}
	}

	clearance := ResolveClearance(index, objectiveID)
	if decision.Allowed {
		return Next{
			Kind: NextObjectiveReady, Objective: objective,
			GateDecision: &decision, Clearance: &clearance,
		}, true
	}

	return Next{
		Kind:         NextObjectiveIntegration,
		Objective:    objective,
		GateDecision: &decision,
		Clearance:    &clearance,
	}, true
}
