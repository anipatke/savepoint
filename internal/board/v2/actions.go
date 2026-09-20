package v2

import (
	"fmt"
	"strings"

	"github.com/opencode/savepoint/internal/data"
)

// ActionKind is a capability the owner can see on the focused record. The
// board does not invent lifecycle actions: executor and checker decisions are
// reported as work for those sessions, not exposed as owner keys.
type ActionKind string

const (
	ActionRecordSelection     ActionKind = "record_selection"
	ActionAcceptCheck         ActionKind = "accept_check"
	ActionCompleteByException ActionKind = "complete_by_exception"
)

const (
	selectionKey      = "p"
	acceptanceKey     = "a"
	exceptionCloseKey = "x"
	ownerBoardSession = "board-owner"
)

// BoardAction is one key the current owner session may use. Decision is
// retained for tests and status copy; rendering never resolves a gate itself.
type BoardAction struct {
	Key        string
	Kind       ActionKind
	TargetKind DetailKind
	TargetID   string
	Label      string
	Decision   data.GateDecision
}

type actionTarget struct {
	Kind DetailKind
	ID   string
}

// actionsForRecord maps one freshly resolved completion decision to owner
// capabilities. The only admitting gate authority is ActorRoleOwner. Owner
// acceptance is the one intentional exception: its gate blocker names the
// owner action that satisfies it, but the decision is not yet Allowed.
func actionsForRecord(index *data.V2Index, target actionTarget) []BoardAction {
	if index == nil || target.ID == "" {
		return nil
	}
	decision, clearance, ok := recordDecision(index, target)
	if !ok {
		return nil
	}

	var actions []BoardAction
	if decision.Allowed && decision.Actor == data.ActorRoleOwner && decision.AllowedByException {
		actions = append(actions, BoardAction{
			Key:        exceptionCloseKey,
			Kind:       ActionCompleteByException,
			TargetKind: target.Kind,
			TargetID:   target.ID,
			Label:      exceptionLabel(decision.Exception),
			Decision:   decision,
		})
	}

	if clearanceIsCurrent(clearance) && hasOwnerAcceptanceBlock(decision) {
		actions = append(actions, BoardAction{
			Key:        acceptanceKey,
			Kind:       ActionAcceptCheck,
			TargetKind: target.Kind,
			TargetID:   target.ID,
			Label:      fmt.Sprintf("accept current Check %s", clearance.Check),
			Decision:   decision,
		})
	}
	return actions
}

func recordDecision(index *data.V2Index, target actionTarget) (data.GateDecision, data.Clearance, bool) {
	switch target.Kind {
	case DetailTask:
		task, ok := index.Tasks[target.ID]
		if !ok {
			return data.GateDecision{}, data.Clearance{}, false
		}
		if task.Status == data.ColumnDone {
			return data.GateDecision{}, data.Clearance{}, true
		}
		var decision data.GateDecision
		switch {
		case task.Status == data.ColumnPlanned:
			decision = data.ResolveTaskStart(index, target.ID)
		case task.Stage == data.StageAudit:
			decision = data.ResolveTaskCompletion(index, target.ID)
		default:
			decision = data.ResolveTaskAdvance(index, target.ID)
		}
		return decision, data.ResolveClearance(index, target.ID), true
	case DetailObjective:
		if _, ok := index.Objectives[target.ID]; !ok {
			return data.GateDecision{}, data.Clearance{}, false
		}
		return data.ResolveObjectiveCompletion(index, target.ID), data.ResolveClearance(index, target.ID), true
	default:
		return data.GateDecision{}, data.Clearance{}, false
	}
}

func exceptionLabel(exception *data.Exception) string {
	if exception == nil {
		return "complete by recorded exception"
	}
	return "complete by exception: " + exception.Reason
}

// refusalTemplates is the one source for gate refusal copy. Detail is filled
// from the resolver's blocker, so a refusal names both the unmet requirement
// and the role whose session owns the next step.
var refusalTemplates = map[string]string{
	"replan":                    "Replan is required (%s); the planner session must resolve it.",
	"dependency":                "Dependency is unmet (%s); the executor session must resolve it.",
	"objective_dependency":      "Objective dependency is unmet (%s); the executor session must resolve it.",
	"clearance_missing":         "Clearance is missing (%s); the checker session must record a Check.",
	"clearance_needs_work":      "Clearance needs work (%s); the executor session must repair it and the checker session must re-check it.",
	"clearance_stale":           "Clearance is stale (%s); the checker session must reassess it.",
	"clearance_unknown":         "Clearance is unknown (%s); the checker session must assess it.",
	"owner_acceptance_required": "Acceptance by the owner is required (%s); the owner session may accept the current Check.",
	"checker_authority":         "Checker authority is required (%s); the checker session must complete the decision.",
	"invalid_state":             "The recorded state is invalid (%s); the planner session must repair the record.",
}

func refusalStatement(kind, detail, dependencyTarget, objectiveTarget string) string {
	detail = blockerDetail(kind, detail, dependencyTarget, objectiveTarget)
	template, ok := refusalTemplates[kind]
	if !ok {
		return fmt.Sprintf("Requirement %s is unmet (%s); the owning session must resolve it.", kind, detail)
	}
	return fmt.Sprintf(template, detail)
}

func decisionRefusal(decision data.GateDecision) string {
	if len(decision.Blockers) == 0 {
		if decision.Allowed {
			return fmt.Sprintf("This action belongs to the %s session.", decision.Actor)
		}
		return "No owner action is available for this record."
	}
	parts := make([]string, 0, len(decision.Blockers))
	for _, blocker := range decision.Blockers {
		dependencyTarget := ""
		if blocker.Dependency != nil {
			dependencyTarget = blocker.Dependency.Target
		}
		objectiveTarget := ""
		if blocker.ObjectiveDependency != nil {
			objectiveTarget = blocker.ObjectiveDependency.Target
		}
		parts = append(parts, refusalStatement(string(blocker.Kind), blocker.Detail, dependencyTarget, objectiveTarget))
	}
	return strings.Join(parts, " ")
}

func blockerDetail(kind, detail, dependencyTarget, objectiveTarget string) string {
	if strings.TrimSpace(detail) != "" {
		return detail
	}
	if dependencyTarget != "" {
		return "dependency " + dependencyTarget
	}
	if objectiveTarget != "" {
		return "objective " + objectiveTarget + " is not ready"
	}
	return kind
}

func actionForKey(actions []BoardAction, key string) (BoardAction, bool) {
	for _, action := range actions {
		if action.Key == key {
			return action, true
		}
	}
	return BoardAction{}, false
}

func (m Model) focusedActionTarget() (actionTarget, bool) {
	if m.State.Index == nil {
		return actionTarget{}, false
	}
	if m.Detail != nil {
		return actionTarget{Kind: m.Detail.Kind, ID: m.Detail.ID}, true
	}
	if m.SidebarFocused {
		if m.ObjectiveCursor < 0 || m.ObjectiveCursor >= len(m.Objectives) {
			return actionTarget{}, false
		}
		return actionTarget{Kind: DetailObjective, ID: m.Objectives[m.ObjectiveCursor].ID()}, true
	}
	cards := m.Cards[m.FocusedColumn]
	if m.FocusedCard < 0 || m.FocusedCard >= len(cards) || cards[m.FocusedCard].Task == nil {
		return actionTarget{}, false
	}
	return actionTarget{Kind: DetailTask, ID: cards[m.FocusedCard].Task.ID}, true
}

func (m Model) focusedActions() []BoardAction {
	target, ok := m.focusedActionTarget()
	if !ok || m.Issues != nil {
		return nil
	}
	actions := actionsForRecord(m.State.Index, target)
	if selection, ok := selectionAction(m.State.Index, target); ok {
		actions = append(actions, selection)
	}
	return actions
}

func selectionAction(index *data.V2Index, target actionTarget) (BoardAction, bool) {
	if index == nil {
		return BoardAction{}, false
	}
	switch target.Kind {
	case DetailObjective:
		if _, ok := index.Objectives[target.ID]; !ok {
			return BoardAction{}, false
		}
	case DetailTask:
		if _, ok := index.Tasks[target.ID]; !ok {
			return BoardAction{}, false
		}
	default:
		return BoardAction{}, false
	}
	return BoardAction{
		Key:        selectionKey,
		Kind:       ActionRecordSelection,
		TargetKind: target.Kind,
		TargetID:   target.ID,
		Label:      "record selection",
	}, true
}

func selectionForTarget(index *data.V2Index, target actionTarget) (data.RouterSelectionV2, error) {
	if index == nil {
		return data.RouterSelectionV2{}, fmt.Errorf("selection requires a loaded project")
	}
	switch target.Kind {
	case DetailObjective:
		if _, ok := index.Objectives[target.ID]; !ok {
			return data.RouterSelectionV2{}, fmt.Errorf("selection target %s is no longer present", target.ID)
		}
		return data.RouterSelectionV2{Objective: target.ID}, nil
	case DetailTask:
		task, ok := index.Tasks[target.ID]
		if !ok {
			return data.RouterSelectionV2{}, fmt.Errorf("selection target %s is no longer present", target.ID)
		}
		return data.RouterSelectionV2{Objective: task.Objective, Task: task.ID}, nil
	default:
		return data.RouterSelectionV2{}, fmt.Errorf("selection is not supported for %s", target.Kind)
	}
}

// focusedActionText names the owner keys available on the focused record, and
// nothing else: the footer is keys only. A blocked action is never explained
// here — that prose belongs to the status-bar message a refused attempt
// produces (decisionRefusal, used only inside the write commands), not to a
// hint line shown at every idle moment regardless of whether anything was
// attempted.
func (m Model) focusedActionText() string {
	actions := m.focusedActions()
	if len(actions) == 0 {
		return ""
	}
	parts := make([]string, 0, len(actions))
	for _, action := range actions {
		parts = append(parts, action.Key+":"+action.Label)
	}
	return strings.Join(parts, "  ")
}

func (m Model) actionForKey(key string) (BoardAction, bool) {
	return actionForKey(m.focusedActions(), key)
}
