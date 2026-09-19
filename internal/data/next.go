package data

import (
	"maps"
	"slices"
)

// SelectionRecordKind names which record family a SelectionNotFound
// diagnostic is about: the router's objective, or its task.
type SelectionRecordKind string

const (
	SelectionRecordObjective SelectionRecordKind = "objective"
	SelectionRecordTask      SelectionRecordKind = "task"
)

// SelectionDiagnosticKind names the two honest ways selection resolution can
// fail to produce a value: a named ID absent from the live index, or a
// resolvable Task whose own objective disagrees with the router's.
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

	// RouterObjective, Task, and TaskObjective are populated when Kind is
	// SelectionMismatch: the Objective the router named, the Task that
	// otherwise resolved, and the Objective that Task's own record declares
	// as its owner.
	RouterObjective string
	Task            string
	TaskObjective   string
}

// Selection is one honest resolved outcome: no record selected, an
// Objective with no Task selected under it, or a matched Objective/Task
// pair. It is the zero value both when nothing was selected and when
// resolution produced a SelectionDiagnostic instead; callers distinguish the
// two by whether ResolveSelection returned a non-nil diagnostic.
type Selection struct {
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
	if router.Objective == "" {
		// ReadStateV2 already refuses a Task selected with no Objective
		// (ErrV2InvalidOwnership), so an empty Objective here means no
		// selection at all, regardless of router phase.
		return Selection{}, nil
	}

	objective, ok := index.Objectives[router.Objective]
	if !ok {
		return Selection{}, &SelectionDiagnostic{Kind: SelectionNotFound, RecordKind: SelectionRecordObjective, ID: router.Objective}
	}

	if router.Task == "" {
		return Selection{Objective: objective}, nil
	}

	task, ok := index.Tasks[router.Task]
	if !ok {
		return Selection{}, &SelectionDiagnostic{Kind: SelectionNotFound, RecordKind: SelectionRecordTask, ID: router.Task}
	}

	if task.Objective != router.Objective {
		return Selection{}, &SelectionDiagnostic{
			Kind:            SelectionMismatch,
			RouterObjective: router.Objective,
			Task:            router.Task,
			TaskObjective:   task.Objective,
		}
	}

	return Selection{Objective: objective, Task: task}, nil
}

// MigrationState is the smallest injected fact ResolveNext needs about an
// in-flight conversion: whether one is pending, and which operation. It
// exists so the projection can outrank every other rung on an incomplete
// migration without importing internal/migrate itself — that package already
// imports internal/data, so the reverse import would close a cycle. main.go
// fills this from migrate.PendingOperation at the same point it resolves the
// project root.
type MigrationState struct {
	Pending     bool
	OperationID string
}

// NextKind names the rung of the precedence ladder a Next value landed on.
// The nine values are the whole ladder, evaluated in this order by
// ResolveNext: pending migration outranks everything, because a project
// midway through conversion holds records whose meaning is not yet settled;
// owner validation sits below check-needed because asking the owner to
// accept work with no current technical clearance would invert the
// authority model E43 established.
type NextKind string

const (
	// NextPendingMigration means an incomplete migration operation exists;
	// every other rung is unreachable until it resolves.
	NextPendingMigration NextKind = "pending_migration"
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
// Clearance only for a rung read from ResolveClearance, and Migration only
// for NextPendingMigration. Nothing here is computed by this package: every
// readiness claim is a value obtained from the existing E43/E44 resolvers
// and carried forward whole (STYLE-07, DATA-02).
type Next struct {
	Kind      NextKind
	Objective *ObjectiveV2
	Task      *TaskV2

	// GateDecision is the Task- or Objective-level decision the rung was
	// derived from: ResolveTaskStart/Advance/Completion for a Task rung,
	// ResolveObjectiveCompletion for NextObjectiveIntegration.
	GateDecision *GateDecision
	// Clearance is the resolved clearance the rung reports on, for
	// NextCheckNeeded, NextOwnerValidationRequired, and
	// NextObjectiveIntegration.
	Clearance *Clearance
	// Migration is set only when Kind is NextPendingMigration.
	Migration MigrationState

	// SelectionDiagnostic is set whenever the router's Objective/Task hint
	// did not resolve, regardless of which rung was ultimately reached: an
	// unresolved selection still yields whatever next action the project's
	// records themselves support.
	SelectionDiagnostic *SelectionDiagnostic

	// Issues are the Issues relevant to Task (when set) or otherwise
	// Objective, resolved here so a rendering surface never has to consult
	// the index on its own to answer what follow-up hangs off the selected
	// record (STYLE-07). Nil when neither is set, or when none are linked.
	Issues []*IssueV2
}

// NextInput carries everything ResolveNext reads: the project's index, the
// decoded router hint, and injected migration state. ResolveNext performs no
// discovery, no filesystem access, and no write; it consults only these
// three values.
type NextInput struct {
	Index     *V2Index
	Router    *RouterStateV2
	Migration MigrationState
}

// ResolveNext computes the one next action for a V2 project: the precedence
// ladder documented on NextKind, evaluated top to bottom. It returns a value
// for every reachable project state, including a fresh project with no
// Objectives or Tasks, and never returns an error.
func ResolveNext(input NextInput) Next {
	if input.Migration.Pending {
		return Next{Kind: NextPendingMigration, Migration: input.Migration}
	}

	selection, diagnostic := ResolveSelection(input.Index, input.Router)

	next := resolveLadder(input.Index, selection)
	next.SelectionDiagnostic = diagnostic
	next.Issues = relevantIssues(input.Index, next)
	return next
}

// relevantIssues collects the Issues linked to next's selected Task or
// Objective, walking only the link maps LoadV2Index already built rather
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

// resolveLadder runs rungs two through nine over selection. An unresolved
// selection arrives here as a zero Selection, which falls straight through
// to the project-wide search the same way a project with no selection at
// all does — the diagnostic naming why is attached by the caller.
func resolveLadder(index *V2Index, selection Selection) Next {
	if task := selection.Task; task != nil && task.Status != ColumnDone {
		return resolveTaskRung(index, task)
	}

	if objectiveID := objectiveInView(selection); objectiveID != "" {
		if next, ok := resolveObjectiveIntegrationRung(index, objectiveID); ok {
			return next
		}
	}

	if next, ok := resolveReadyRung(index); ok {
		return next
	}

	return Next{Kind: NextPlanObjective}
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
