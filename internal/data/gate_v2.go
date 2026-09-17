package data

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

// ClearanceState is the resolved clearance for a Task or Objective's
// recorded evidence. It is derived only from Checks and freshness
// assessments already present in the loaded index — never from file
// scanning, hashing, timestamp comparison, or Git inspection.
type ClearanceState string

const (
	// ClearanceCurrent means the target's latest Check is CLEAR and a
	// freshness assessment names that same Check with state: current.
	ClearanceCurrent ClearanceState = "current"
	// ClearanceNeedsWork means the target's latest Check recorded NEEDS
	// WORK.
	ClearanceNeedsWork ClearanceState = "needs_work"
	// ClearanceStale means the latest Check is CLEAR but the freshness
	// assessment names a different Check, or does not assert current for
	// the latest one.
	ClearanceStale ClearanceState = "stale"
	// ClearanceUnknown means the latest Check is CLEAR but the target
	// carries no freshness assessment at all.
	ClearanceUnknown ClearanceState = "unknown"
	// ClearanceMissing means the target has no recorded Check.
	ClearanceMissing ClearanceState = "missing"
)

// Clearance is the resolved clearance for one Task or Objective: the state
// it landed in, the latest Check it was derived from, and the freshness
// assessment consulted, if any (which carries the recorded basis).
type Clearance struct {
	State     ClearanceState
	Check     string // C### of the target's latest Check; empty when State is missing
	Freshness *Freshness
}

// ResolveClearance resolves targetID's clearance from index, re-reading the
// latest Check and evidence already loaded there. targetID may name either a
// Task or an Objective; both share the same evidence shape. It returns a
// decision for every reachable state rather than an error: a target with no
// evidence and no Check resolves to missing, never a failure.
func ResolveClearance(index *V2Index, targetID string) Clearance {
	latestCheckID := index.LatestCheck[targetID]
	if latestCheckID == "" {
		return Clearance{State: ClearanceMissing}
	}

	latestCheck := index.Checks[latestCheckID]
	if latestCheck.Result == CheckResultNeedsWork {
		return Clearance{State: ClearanceNeedsWork, Check: latestCheckID}
	}

	freshness := evidenceFreshness(index, targetID)
	if freshness == nil {
		return Clearance{State: ClearanceUnknown, Check: latestCheckID}
	}

	if freshness.Check == latestCheckID && freshness.State == FreshnessCurrent {
		if independentCheckerProvenance(latestCheck, freshness) {
			return Clearance{State: ClearanceCurrent, Check: latestCheckID, Freshness: freshness}
		}
		// A current claim without checker provenance is not a stale claim: it
		// is unknown whether the recorded CLEAR result was independently
		// checked. Keep the five-state clearance vocabulary and fail closed.
		return Clearance{State: ClearanceUnknown, Check: latestCheckID, Freshness: freshness}
	}

	return Clearance{State: ClearanceStale, Check: latestCheckID, Freshness: freshness}
}

// independentCheckerProvenance is the authority boundary for a CLEAR Check.
// Both the recorded Check and the freshness assessment must identify a
// checker session. Savepoint deliberately does not authenticate either
// session; it only refuses to treat executor, planner, or owner self-report
// as clearance evidence.
func independentCheckerProvenance(check *CheckV2, freshness *Freshness) bool {
	return check != nil && check.Result == CheckResultClear && isCheckerActor(check.CheckedBy) &&
		freshness != nil && freshness.State == FreshnessCurrent && isCheckerActor(freshness.AssessedBy)
}

func untrustedCurrentClearance(index *V2Index, targetID, checkID string) bool {
	check := index.Checks[checkID]
	freshness := evidenceFreshness(index, targetID)
	return check != nil && check.Result == CheckResultClear && freshness != nil &&
		freshness.Check == checkID && freshness.State == FreshnessCurrent &&
		!independentCheckerProvenance(check, freshness)
}

// evidenceFreshness looks up targetID's recorded freshness assessment,
// checking Tasks then Objectives. It returns nil when the target carries no
// evidence or no freshness block, never a healed value.
func evidenceFreshness(index *V2Index, targetID string) *Freshness {
	if task, ok := index.Tasks[targetID]; ok {
		if task.Evidence == nil {
			return nil
		}
		return task.Evidence.Freshness
	}
	if objective, ok := index.Objectives[targetID]; ok {
		if objective.Evidence == nil {
			return nil
		}
		return objective.Evidence.Freshness
	}
	return nil
}

// GateBlockKind names why a Task start, advance, or completion decision is
// not allowed.
type GateBlockKind string

const (
	// GateBlockReplan means the Task carries a recorded replan flag, which
	// blocks start and advance until it is cleared by an explicit write.
	GateBlockReplan GateBlockKind = "replan"
	// GateBlockDependency means one of the Task's dependencies is
	// unsatisfied; Dependency names which one and why.
	GateBlockDependency GateBlockKind = "dependency"
	// GateBlockClearanceMissing means the Task has no recorded Check.
	GateBlockClearanceMissing GateBlockKind = "clearance_missing"
	// GateBlockClearanceNeedsWork means the Task's latest Check recorded
	// NEEDS WORK.
	GateBlockClearanceNeedsWork GateBlockKind = "clearance_needs_work"
	// GateBlockClearanceStale means the Task's freshness assessment does not
	// name its latest Check as current.
	GateBlockClearanceStale GateBlockKind = "clearance_stale"
	// GateBlockClearanceUnknown means the Task's latest Check is CLEAR but
	// carries no freshness assessment.
	GateBlockClearanceUnknown GateBlockKind = "clearance_unknown"
	// GateBlockOwnerAcceptance means the Task declares owner_validation.
	// required and the owner has not accepted the current Check.
	GateBlockOwnerAcceptance GateBlockKind = "owner_acceptance_required"
	// GateBlockInvalidState means the Task's recorded status/stage is not a
	// state AdvanceTaskLifecycleState recognizes, or, for an Objective
	// completion decision, that one of its owned Tasks is not done.
	GateBlockInvalidState GateBlockKind = "invalid_state"
	// GateBlockCheckerAuthority means a CLEAR Check or its current freshness
	// assessment was not recorded by an identified checker session.
	GateBlockCheckerAuthority GateBlockKind = "checker_authority"
	// GateBlockObjectiveDependency means the Task's owning Objective has an
	// unsatisfied Objective dependency; ObjectiveDependency names which one
	// and why, so a consumer can explain that the wait is at the Objective
	// level rather than the Task's own dependencies.
	GateBlockObjectiveDependency GateBlockKind = "objective_dependency"
)

// GateBlocker names one unmet requirement blocking a start, advance, or
// completion decision. Dependency is set only for GateBlockDependency, so
// callers get the same typed DependencyBlock ResolveTaskDependencyV2 already
// reports rather than a second dependency vocabulary. ObjectiveDependency is
// set only for GateBlockObjectiveDependency, mirroring the same pattern for
// ResolveObjectiveDependency's typed block.
type GateBlocker struct {
	Kind                GateBlockKind
	Detail              string
	Dependency          *DependencyBlock
	ObjectiveDependency *ObjectiveDependencyBlock
}

// GateDecision is the resolved outcome for one Task start, advance, or
// completion action: whether it is allowed, which actor authority may take
// it, and every unmet requirement blocking it when it is not. A decision
// allowed only through a recorded exception sets AllowedByException and
// Exception instead of reporting a CLEAR result or current clearance.
type GateDecision struct {
	Allowed            bool
	Actor              ActorRole // meaningful only when Allowed is true
	Blockers           []GateBlocker
	AllowedByException bool
	Exception          *Exception // set only when AllowedByException
}

// ResolveTaskStart decides whether taskID may move from planned to
// in_progress. Starting is valid only for a planned Task with no stage.
// Starting is blocked by a recorded replan flag and by any unsatisfied
// dependency; every unmet dependency is named, not just the first. It grants
// executor authority: the agent picking up the Task starts it.
func ResolveTaskStart(index *V2Index, taskID string) GateDecision {
	task, ok := index.Tasks[taskID]
	if !ok {
		return GateDecision{}
	}

	var blockers []GateBlocker
	if task.Evidence != nil && task.Evidence.Replan != nil {
		blockers = append(blockers, GateBlocker{
			Kind:   GateBlockReplan,
			Detail: task.Evidence.Replan.Reason,
		})
	}

	if task.Status != ColumnPlanned || task.Stage != "" {
		blockers = append(blockers, GateBlocker{
			Kind:   GateBlockInvalidState,
			Detail: fmt.Sprintf("task start requires status planned with no stage (got status %q, stage %q)", task.Status, task.Stage),
		})
	}
	if len(blockers) > 0 {
		return GateDecision{Blockers: blockers}
	}

	for _, dep := range task.DependsOn {
		decision := ResolveTaskDependencyV2(index, dep)
		if decision.Satisfied {
			continue
		}
		blockers = append(blockers, GateBlocker{
			Kind:       GateBlockDependency,
			Detail:     fmt.Sprintf("dependency %s requires %s: %s", decision.Block.Target, decision.Block.Requires, decision.Block.Kind),
			Dependency: decision.Block,
		})
	}

	if len(blockers) > 0 {
		return GateDecision{Blockers: blockers}
	}

	if objective, ok := index.Objectives[task.Objective]; ok {
		for _, dep := range objective.DependsOn {
			decision := ResolveObjectiveDependency(index, dep)
			if decision.Satisfied {
				continue
			}
			blockers = append(blockers, GateBlocker{
				Kind:                GateBlockObjectiveDependency,
				Detail:              fmt.Sprintf("owning objective %s waits on objective %s: %s", task.Objective, decision.Block.Target, decision.Block.Kind),
				ObjectiveDependency: decision.Block,
			})
		}
	}

	if len(blockers) > 0 {
		return GateDecision{Blockers: blockers}
	}
	return GateDecision{Allowed: true, Actor: ActorRoleExecutor}
}

// ResolveTaskAdvance decides whether an in_progress Task may move to its
// next stage. From the audit stage the next move is completion, so it
// defers to ResolveTaskCompletion rather than restating that decision here.
// Otherwise it reuses AdvanceTaskLifecycleState to confirm the next stage is
// a recognized canonical move — build to test, or test to audit — so the
// stage sequence stays defined once in lifecycle.go. A replan flag is checked
// before every route, including the audit-to-completion route.
func ResolveTaskAdvance(index *V2Index, taskID string) GateDecision {
	task, ok := index.Tasks[taskID]
	if !ok {
		return GateDecision{}
	}

	if task.Evidence != nil && task.Evidence.Replan != nil {
		return GateDecision{Blockers: []GateBlocker{{
			Kind:   GateBlockReplan,
			Detail: task.Evidence.Replan.Reason,
		}}}
	}

	if task.Status != ColumnInProgress {
		return GateDecision{Blockers: []GateBlocker{{
			Kind:   GateBlockInvalidState,
			Detail: fmt.Sprintf("task advance requires status in_progress (got %q)", task.Status),
		}}}
	}

	if task.Stage == StageAudit {
		return ResolveTaskCompletion(index, taskID)
	}

	if _, err := AdvanceTaskLifecycleState(TaskLifecycleState{Status: task.Status, Stage: task.Stage}); err != nil {
		return GateDecision{Blockers: []GateBlocker{{
			Kind:   GateBlockInvalidState,
			Detail: err.Error(),
		}}}
	}

	return GateDecision{Allowed: true, Actor: ActorRoleExecutor}
}

// ResolveTaskCompletion decides whether taskID may close. Completion is valid
// only for a Task at status in_progress and stage audit. A technical Task —
// one with no owner_validation.required — closes under checker authority
// once its clearance is current. A Task declaring owner_validation.required
// additionally needs the owner's acceptance of that same current Check;
// acceptance bound to a Check a later Check has superseded does not count,
// so the Task stays blocked until the owner accepts again. Missing,
// needs_work, stale, and unknown clearance each block completion with a
// distinct reason. When completion would otherwise be blocked, a recorded
// exception naming the Task's current latest Check grants completion by
// exception under owner authority instead — never as a CLEAR result — and an
// exception naming any other Check does not apply.
func ResolveTaskCompletion(index *V2Index, taskID string) GateDecision {
	task, ok := index.Tasks[taskID]
	if !ok {
		return GateDecision{}
	}

	if task.Evidence != nil && task.Evidence.Replan != nil {
		return GateDecision{Blockers: []GateBlocker{{
			Kind:   GateBlockReplan,
			Detail: task.Evidence.Replan.Reason,
		}}}
	}

	if task.Status != ColumnInProgress || task.Stage != StageAudit {
		return GateDecision{Blockers: []GateBlocker{{
			Kind:   GateBlockInvalidState,
			Detail: fmt.Sprintf("task completion requires status in_progress and stage audit (got status %q, stage %q)", task.Status, task.Stage),
		}}}
	}

	clearance := ResolveClearance(index, taskID)
	var blockers []GateBlocker

	switch clearance.State {
	case ClearanceMissing:
		blockers = append(blockers, GateBlocker{Kind: GateBlockClearanceMissing, Detail: "no recorded check"})
	case ClearanceNeedsWork:
		blockers = append(blockers, GateBlocker{Kind: GateBlockClearanceNeedsWork, Detail: fmt.Sprintf("latest check %s recorded NEEDS WORK", clearance.Check)})
	case ClearanceStale:
		blockers = append(blockers, GateBlocker{Kind: GateBlockClearanceStale, Detail: fmt.Sprintf("freshness assessment does not name latest check %s as current", clearance.Check)})
	case ClearanceUnknown:
		if untrustedCurrentClearance(index, taskID, clearance.Check) {
			blockers = append(blockers, GateBlocker{Kind: GateBlockCheckerAuthority, Detail: fmt.Sprintf("latest check %s lacks independent checker provenance", clearance.Check)})
		} else {
			blockers = append(blockers, GateBlocker{Kind: GateBlockClearanceUnknown, Detail: fmt.Sprintf("no freshness assessment recorded for latest check %s", clearance.Check)})
		}
	case ClearanceCurrent:
		if ownerValidationRequired(task.Evidence) && !ownerAcceptedCheck(task.Evidence, clearance.Check) {
			blockers = append(blockers, GateBlocker{Kind: GateBlockOwnerAcceptance, Detail: fmt.Sprintf("owner has not accepted current check %s", clearance.Check)})
		}
	}

	if len(blockers) == 0 {
		return GateDecision{Allowed: true, Actor: ActorRoleChecker}
	}

	if exception := applicableException(task.Evidence, index.LatestCheck[taskID]); exception != nil {
		return GateDecision{Allowed: true, Actor: ActorRoleOwner, AllowedByException: true, Exception: exception}
	}

	return GateDecision{Blockers: blockers}
}

// ownerValidationRequired and ownerAcceptedCheck read the shared Evidence
// block so Task and Objective completion decisions apply the same owner-
// acceptance rule without restating it.
func ownerValidationRequired(evidence *Evidence) bool {
	return evidence != nil && evidence.OwnerValidation != nil && evidence.OwnerValidation.Required
}

func ownerAcceptedCheck(evidence *Evidence, checkID string) bool {
	if evidence == nil || evidence.OwnerValidation == nil || checkID == "" {
		return false
	}
	accepted := evidence.OwnerValidation
	return accepted.AcceptedCheck == checkID && accepted.AcceptedBy.Role == ActorRoleOwner && strings.TrimSpace(accepted.AcceptedBy.Session) != ""
}

// applicableException returns evidence's recorded exception only when it
// names latestCheckID, the target's current latest Check. An exception
// naming any other Check — one a later Check has superseded — does not carry
// forward and returns nil, so the caller falls back to reporting normal
// blockers. Shared by Task and Objective completion.
func applicableException(evidence *Evidence, latestCheckID string) *Exception {
	if evidence == nil || evidence.Exception == nil {
		return nil
	}
	exception := evidence.Exception
	if latestCheckID == "" || exception.Check != latestCheckID {
		return nil
	}
	return exception
}

// ConsistencyDiagnosticKind names one way a Task's recorded status can
// contradict its recorded evidence after a hand edit.
type ConsistencyDiagnosticKind string

const (
	// ConsistencyDoneWithoutClearance means a Task's status is done but its
	// clearance is not current.
	ConsistencyDoneWithoutClearance ConsistencyDiagnosticKind = "done_without_current_clearance"
	// ConsistencyAcceptanceSuperseded means a Task's recorded owner
	// acceptance names a Check that is no longer its latest Check.
	ConsistencyAcceptanceSuperseded ConsistencyDiagnosticKind = "acceptance_names_superseded_check"
	// ConsistencyEvidenceContradictsStatus means a Task's recorded evidence
	// would allow completion but its status was never advanced to done.
	ConsistencyEvidenceContradictsStatus ConsistencyDiagnosticKind = "evidence_contradicts_status"
)

// ConsistencyDiagnostic names one inconsistency InspectTaskConsistency found
// between a Task's recorded status and its recorded evidence.
type ConsistencyDiagnostic struct {
	Task   string
	Kind   ConsistencyDiagnosticKind
	Detail string
}

// InspectTaskConsistency reports every inconsistency between a Task's
// recorded status and its recorded evidence without rewriting any record: a
// done Task without current clearance, an owner acceptance naming a Check a
// later Check has superseded, and evidence that would clear completion while
// status was never advanced to done. It walks Task IDs in sorted order and
// returns every problem found across every Task, not only the first. It
// does not authenticate who wrote the evidence or prevent further external
// edits; it only names mismatches already present in the loaded index.
func InspectTaskConsistency(index *V2Index) []ConsistencyDiagnostic {
	var diagnostics []ConsistencyDiagnostic

	for _, id := range slices.Sorted(maps.Keys(index.Tasks)) {
		task := index.Tasks[id]
		clearance := ResolveClearance(index, id)

		if task.Status == ColumnDone && clearance.State != ClearanceCurrent {
			diagnostics = append(diagnostics, ConsistencyDiagnostic{
				Task:   id,
				Kind:   ConsistencyDoneWithoutClearance,
				Detail: fmt.Sprintf("task is done but clearance is %s", clearance.State),
			})
		}

		if task.Evidence != nil && task.Evidence.OwnerValidation != nil {
			accepted := task.Evidence.OwnerValidation.AcceptedCheck
			latest := index.LatestCheck[id]
			if accepted != "" && accepted != latest {
				diagnostics = append(diagnostics, ConsistencyDiagnostic{
					Task:   id,
					Kind:   ConsistencyAcceptanceSuperseded,
					Detail: fmt.Sprintf("owner accepted %s but the latest check is %s", accepted, latest),
				})
			}
		}

		if task.Status != ColumnDone {
			completion := ResolveTaskCompletion(index, id)
			if completion.Allowed && !completion.AllowedByException {
				diagnostics = append(diagnostics, ConsistencyDiagnostic{
					Task:   id,
					Kind:   ConsistencyEvidenceContradictsStatus,
					Detail: fmt.Sprintf("evidence clears completion but status is %s", task.Status),
				})
			}
		}
	}

	return diagnostics
}
