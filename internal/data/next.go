package data

import (
	"maps"
	"slices"
)

// SelectionRecordKind names which record family a selection diagnostic is
// about: the router's Release, Objective, or Task.
type SelectionRecordKind string

const (
	SelectionRecordRelease   SelectionRecordKind = "release"
	SelectionRecordObjective SelectionRecordKind = "objective"
	SelectionRecordTask      SelectionRecordKind = "task"
)

// SelectionDiagnosticKind names why the router's contextual selection could
// not be honored. Release diagnostics are separate from the older Task /
// Objective diagnostics so a renderer can explain an archived Release,
// unassigned Objective, and cross-Release mismatch without inspecting the
// index again.
type SelectionDiagnosticKind string

const (
	// SelectionNotFound means a router-named Objective or Task ID does not
	// exist among the project's live records. Archived records are out of
	// scope for lookup, so this is also what a migrated-away ID reports.
	SelectionNotFound SelectionDiagnosticKind = "not_found"
	// SelectionMismatch means the router-named Objective and Task both
	// resolve, but the Task's own objective field names a different
	// Objective. Ownership always comes from the Task record, never the
	// router, so this is reported rather than silently reconciled.
	SelectionMismatch SelectionDiagnosticKind = "mismatch"
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

// SelectionDiagnostic is the typed, branchable result of a selection that
// could not be honored. A caller inspects Kind and the ID fields below it;
// nothing here is a formatted message standing in for the IDs read.
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
// with no Objective selected, an Objective with no Task selected under it, or
// a matched Release/Objective/Task context. When a Release diagnostic names a
// live Release, the Release is retained for rendering while Objective/Task
// action remains unresolved.
type Selection struct {
	Release   *ReleaseV2
	Objective *ObjectiveV2
	Task      *TaskV2
}

// ResolveSelection turns router's Objective/Task hint into an exactly
// matched Selection against index, or a typed SelectionDiagnostic.
//
// Matching is by exact global ID only: a map lookup against index.Objectives
// and index.Tasks, with no prefix, suffix, numeric-proximity, title, or
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

	var release *ReleaseV2
	if router.Release != "" {
		candidate, ok := index.Releases[router.Release]
		if !ok {
			return Selection{}, &SelectionDiagnostic{
				Kind: SelectionReleaseNotFound, RecordKind: SelectionRecordRelease,
				ID: router.Release, Release: router.Release,
			}
		}
		release = candidate
		if release.LegacyCompletion != nil {
			return Selection{Release: release}, &SelectionDiagnostic{
				Kind: SelectionReleaseArchived, RecordKind: SelectionRecordRelease,
				ID: router.Release, Release: router.Release,
			}
		}
	}

	if router.Objective == "" {
		// ReadStateV2 already refuses a Task selected with no Objective
		// (ErrV2InvalidOwnership), so an empty Objective here means no
		// selection at all, regardless of router phase.
		return Selection{Release: release}, nil
	}

	objective, ok := index.Objectives[router.Objective]
	if !ok {
		return Selection{Release: release}, &SelectionDiagnostic{
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
		return Selection{Release: release, Objective: objective}, nil
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

	return Selection{Release: release, Objective: objective, Task: task}, nil
}

// NextKind names the rung of the precedence ladder a Next value landed on.
// The global values are evaluated in their existing order; Release-specific
// values are reached only after a valid Release context has narrowed the
// candidate records.
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
	// NextReleaseIntegration means the selected Release still has member
	// Objective work or a material Release-level integration blocker.
	NextReleaseIntegration NextKind = "release_integration"
	// NextReleaseCheckNeeded means member Objective work is complete but the
	// selected Release lacks usable current technical clearance.
	NextReleaseCheckNeeded NextKind = "release_check_needed"
	// NextReleaseOwnerValidationRequired means the Release Check is current
	// and the owner must accept that exact Check.
	NextReleaseOwnerValidationRequired NextKind = "release_owner_validation_required"
	// NextReleaseReady means the Release completion decision is allowed and
	// the next action is to record the Release as done.
	NextReleaseReady NextKind = "release_ready"
	// NextReady means no record is under active work, but a specific Task
	// or Objective elsewhere in the project is ready to pick up.
	NextReady NextKind = "ready"
	// NextPlanObjective means nothing in the project is ready: there is no
	// in-progress work and no ready Task or Objective, so the next action
	// is planning — Idea/Design work on a new or existing Objective. A
	// fresh project with no Objectives lands here.
	NextPlanObjective NextKind = "plan_objective"
)

// Next is the one derived answer ResolveNext returns for a project: which
// rung of the ladder it landed on, the Objective and/or Task selected there
// (when the rung names one), and the resolver-derived value the rung was
// read from. Every field beyond Kind is populated only when the rung makes
// it meaningful — GateDecision only for a rung read from ResolveTaskStart,
// ResolveTaskAdvance, ResolveTaskCompletion, or ResolveObjectiveCompletion,
// and Clearance only for a rung read from ResolveClearance. Nothing here is
// computed by this package: every readiness claim is a value obtained from
// the existing E43/E44 resolvers and carried forward whole (STYLE-07, DATA-02).
type Next struct {
	Kind      NextKind
	Release   *ReleaseV2
	Objective *ObjectiveV2
	Task      *TaskV2

	// GateDecision is the decision the rung was derived from. For a Release
	// rung it is ResolveReleaseCompletion; for older rungs it remains the
	// Task- or Objective-level decision described below.
	GateDecision *GateDecision
	// Clearance is the resolved clearance the rung reports on, for
	// NextCheckNeeded, NextOwnerValidationRequired, NextObjectiveIntegration,
	// and the Release-specific evidence rungs.
	Clearance *Clearance

	// SelectionDiagnostic is set whenever the router's contextual hint did not
	// resolve, regardless of which rung was ultimately reached: an unresolved
	// selection still yields whatever next action the project's records
	// themselves support.
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

// ResolveNext computes the one next action for a V2 project: the precedence
// ladder documented on NextKind, evaluated top to bottom. It returns a value
// for every reachable project state, including a fresh project with no
// Objectives or Tasks, and never returns an error. A nil Index or Router is
// read as an empty project and an empty selection rather than panicking; see
// the boundary note in the body.
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
	next.SelectionDiagnostic = diagnostic
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

// resolveLadder keeps the pre-E51 ladder intact when no valid Release context
// is selected. A Release diagnostic deliberately falls through to the same
// global search as any other unresolved router hint, so a bad context never
// hides unrelated work.
func resolveLadder(index *V2Index, selection Selection, diagnostic *SelectionDiagnostic) Next {
	if selection.Release != nil && diagnostic == nil {
		return resolveReleaseLadder(index, selection)
	}

	if task := selection.Task; task != nil && task.Status != ColumnDone {
		return resolveTaskRung(index, task)
	}

	if objectiveID := objectiveInView(selection); objectiveID != "" {
		if next, ok := resolveObjectiveIntegrationRung(index, objectiveID); ok {
			return next
		}
	}

	if next, ok := resolveProjectWideIntegrationRung(index); ok {
		return next
	}

	if next, ok := resolveReadyRung(index); ok {
		return next
	}

	return Next{Kind: NextPlanObjective}
}

// resolveReleaseLadder applies the existing Task and Objective resolvers to
// only the selected Release's derived members. Once those records have no
// actionable work left, the Release completion resolver supplies the Check,
// owner-acceptance, or ready outcome. No Release membership is inferred from
// paths, titles, or Task IDs.
func resolveReleaseLadder(index *V2Index, selection Selection) Next {
	release := selection.Release
	if task := selection.Task; task != nil && task.Status != ColumnDone {
		return withRelease(resolveTaskRung(index, task), release)
	}

	if next, ok := resolveReleaseActiveTaskRung(index, release.ID); ok {
		return withRelease(next, release)
	}

	if objectiveID := objectiveInView(selection); objectiveID != "" {
		if next, ok := resolveObjectiveIntegrationRung(index, objectiveID); ok {
			return withRelease(next, release)
		}
	}

	if next, ok := resolveReleaseProjectWideIntegrationRung(index, release.ID); ok {
		return withRelease(next, release)
	}

	if next, ok := resolveReleaseReadyRung(index, release.ID); ok {
		return withRelease(next, release)
	}

	return resolveReleaseCompletionRung(index, release)
}

// resolveReleaseActiveTaskRung keeps an in-progress Task actionable after a
// Release switch clears an Objective/Task hint. The global no-Release path
// intentionally retains its pre-E51 behavior; this search is only the
// contextual Release projection and is deterministic by Task ID.
func resolveReleaseActiveTaskRung(index *V2Index, releaseID string) (Next, bool) {
	members := make(map[string]struct{}, len(index.ReleaseObjectives[releaseID]))
	for _, objectiveID := range sortedReleaseObjectiveIDs(index, releaseID) {
		members[objectiveID] = struct{}{}
	}

	for _, id := range slices.Sorted(maps.Keys(index.Tasks)) {
		task := index.Tasks[id]
		if task.Status != ColumnInProgress {
			continue
		}
		if _, ok := members[task.Objective]; !ok {
			continue
		}
		return resolveTaskRung(index, task), true
	}
	return Next{}, false
}

func withRelease(next Next, release *ReleaseV2) Next {
	next.Release = release
	return next
}

func sortedReleaseObjectiveIDs(index *V2Index, releaseID string) []string {
	ids := slices.Clone(index.ReleaseObjectives[releaseID])
	slices.Sort(ids)
	return ids
}

// resolveReleaseProjectWideIntegrationRung mirrors the existing project-wide
// Objective integration search, restricted to the selected Release's reverse
// membership map and kept in deterministic Objective ID order.
func resolveReleaseProjectWideIntegrationRung(index *V2Index, releaseID string) (Next, bool) {
	for _, objectiveID := range sortedReleaseObjectiveIDs(index, releaseID) {
		objective := index.Objectives[objectiveID]
		if objective == nil || objective.Status == ColumnDone || len(index.ObjectiveTasks[objectiveID]) == 0 {
			continue
		}
		if next, ok := resolveObjectiveIntegrationRung(index, objectiveID); ok {
			return next, true
		}
	}
	return Next{}, false
}

// resolveReleaseReadyRung is the release-scoped version of the global ready
// search. It only admits Tasks whose owning Objective names the selected
// Release, then Objectives in that same derived member set.
func resolveReleaseReadyRung(index *V2Index, releaseID string) (Next, bool) {
	members := make(map[string]struct{}, len(index.ReleaseObjectives[releaseID]))
	for _, objectiveID := range sortedReleaseObjectiveIDs(index, releaseID) {
		members[objectiveID] = struct{}{}
	}

	for _, id := range slices.Sorted(maps.Keys(index.Tasks)) {
		task := index.Tasks[id]
		if task.Status != ColumnPlanned {
			continue
		}
		if _, ok := members[task.Objective]; !ok {
			continue
		}
		decision := ResolveTaskStart(index, id)
		if decision.Allowed {
			return Next{Kind: NextReady, Task: task, GateDecision: &decision}, true
		}
	}

	for _, objectiveID := range slices.Sorted(maps.Keys(members)) {
		objective := index.Objectives[objectiveID]
		if objective == nil || objective.Status == ColumnDone || len(index.ObjectiveTasks[objectiveID]) > 0 {
			continue
		}
		if objectiveDependenciesSatisfied(index, objective) {
			return Next{Kind: NextReady, Objective: objective}, true
		}
	}

	return Next{}, false
}

// resolveReleaseCompletionRung projects the one Release gate decision only
// after member Task/Objective work has had its chance to win. Clearance and
// owner acceptance remain the same typed evidence values used by the gate
// resolver; this function only maps them to contextual action rungs.
func resolveReleaseCompletionRung(index *V2Index, release *ReleaseV2) Next {
	decision := ResolveReleaseCompletion(index, release.ID)
	clearance := ResolveClearance(index, release.ID)
	next := Next{Release: release, GateDecision: &decision, Clearance: &clearance}

	if decision.Allowed {
		next.Kind = NextReleaseReady
		return next
	}

	if len(decision.Blockers) == 1 && decision.Blockers[0].Kind == GateBlockOwnerAcceptance {
		next.Kind = NextReleaseOwnerValidationRequired
		return next
	}

	for _, blocker := range decision.Blockers {
		switch blocker.Kind {
		case GateBlockClearanceMissing, GateBlockClearanceNeedsWork,
			GateBlockClearanceStale, GateBlockClearanceUnknown, GateBlockCheckerAuthority:
			next.Kind = NextReleaseCheckNeeded
			return next
		}
	}

	next.Kind = NextReleaseIntegration
	return next
}

// resolveProjectWideIntegrationRung answers rung seven for an Objective the
// router does not name: the lowest-ID Objective that is not yet done, owns at
// least one Task, and whose own integration clearance is what remains unmet.
// Without it, an outstanding integration Check is reported only while the
// router happens to select its Objective, so the same project state would
// answer differently depending on a hint that decides nothing. An Objective
// already recorded done is skipped: thin evidence behind a closed Objective is
// doctor's missing-evidence report, not the project's next action.
func resolveProjectWideIntegrationRung(index *V2Index) (Next, bool) {
	for _, id := range slices.Sorted(maps.Keys(index.Objectives)) {
		objective := index.Objectives[id]
		if objective.Status == ColumnDone || len(index.ObjectiveTasks[id]) == 0 {
			continue
		}
		if next, ok := resolveObjectiveIntegrationRung(index, id); ok {
			return next, true
		}
	}
	return Next{}, false
}

// objectiveInView names the Objective rung seven should evaluate: the
// selected Task's own owner once that Task is done, or the directly
// selected Objective. It returns "" when neither is in view, so the caller
// skips straight to the project-wide search.
func objectiveInView(selection Selection) string {
	switch {
	case selection.Task != nil:
		return selection.Task.Objective
	case selection.Objective != nil:
		return selection.Objective.ID
	default:
		return ""
	}
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

// resolveObjectiveIntegrationRung reads ResolveObjectiveCompletion for
// objectiveID and reports rung seven only when every owned Task is already
// done and the Objective's own integration clearance (or owner acceptance)
// is what remains unmet. When ResolveObjectiveCompletion is blocked by an
// incomplete owned Task instead, the real next action is that Task, so this
// returns ok=false and lets the caller fall through to the project-wide
// search, which walks every Task including that one. When completion is
// already allowed — current clearance, or allowed by exception — there is
// nothing to report at this rung either, for the same reason.
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
	if decision.Allowed {
		return Next{}, false
	}
	for _, blocker := range decision.Blockers {
		if blocker.Kind == GateBlockInvalidState {
			return Next{}, false
		}
	}

	clearance := ResolveClearance(index, objectiveID)
	return Next{
		Kind:         NextObjectiveIntegration,
		Objective:    objective,
		GateDecision: &decision,
		Clearance:    &clearance,
	}, true
}

// resolveReadyRung is the project-wide search for rung eight: the first
// planned Task (in sorted ID order, for a deterministic answer) that
// ResolveTaskStart allows, or, failing that, the first Objective (also
// sorted) that owns no Task yet and whose own dependencies are satisfied.
// It reads every candidate through the existing resolvers rather than
// inferring readiness from status alone.
func resolveReadyRung(index *V2Index) (Next, bool) {
	for _, id := range slices.Sorted(maps.Keys(index.Tasks)) {
		task := index.Tasks[id]
		if task.Status != ColumnPlanned {
			continue
		}
		decision := ResolveTaskStart(index, id)
		if decision.Allowed {
			return Next{Kind: NextReady, Task: task, GateDecision: &decision}, true
		}
	}

	for _, id := range slices.Sorted(maps.Keys(index.Objectives)) {
		objective := index.Objectives[id]
		if objective.Status == ColumnDone || len(index.ObjectiveTasks[id]) > 0 {
			continue
		}
		if objectiveDependenciesSatisfied(index, objective) {
			return Next{Kind: NextReady, Objective: objective}, true
		}
	}

	return Next{}, false
}

// objectiveDependenciesSatisfied re-resolves each of objective's own
// dependencies through ResolveObjectiveDependency rather than reading any
// cached judgement.
func objectiveDependenciesSatisfied(index *V2Index, objective *ObjectiveV2) bool {
	for _, dep := range objective.DependsOn {
		if !ResolveObjectiveDependency(index, dep).Satisfied {
			return false
		}
	}
	return true
}
