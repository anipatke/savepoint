package resume

import (
	"fmt"
	"io"
	"strings"

	"github.com/opencode/savepoint/internal/data"
)

// Render writes next as deterministic, plain-text narrative to w. It takes
// no project root and performs no discovery: every fact it writes is read
// straight off next, the value data.ResolveNext already resolved (E48-Detail
// §"Resume reads; it does not act."). Render carries no ANSI styling and
// produces byte-identical output for the same next on every call.
func Render(w io.Writer, next data.Next) error {
	_, err := io.WriteString(w, renderText(next))
	return err
}

// renderText builds the full narrative as a string, so exact-text and
// determinism tests never have to go through an io.Writer.
func renderText(next data.Next) string {
	var lines []string

	if next.SelectionDiagnostic != nil {
		lines = append(lines, "Selection: "+SelectionPhrase(next.SelectionDiagnostic), "")
	}

	lines = append(lines, identityLines(next)...)
	lines = append(lines, rungLines(next)...)
	lines = append(lines, issueLines(next.Issues)...)
	lines = append(lines, "Next action: "+ActionPhrase(next))

	return strings.Join(lines, "\n") + "\n"
}

// identityLines names the selected Objective and/or Task and its recorded
// implementation state. It renders nothing for a rung with no selection —
// pending migration and plan-next-Objective both carry neither.
func identityLines(next data.Next) []string {
	var lines []string
	if next.Objective != nil {
		lines = append(lines, fmt.Sprintf("Objective: %s — %s", next.Objective.ID, next.Objective.Title))
		if next.Task == nil {
			lines = append(lines, "Implementation: "+objectiveStatusPhrase(next.Objective))
		}
	}
	if next.Task != nil {
		lines = append(lines, fmt.Sprintf("Task: %s — %s", next.Task.ID, next.Task.Title))
		lines = append(lines, "Implementation: "+implementationPhrase(next.Task))
	}
	if len(lines) > 0 {
		lines = append(lines, "")
	}
	return lines
}

// rungLines is EvidenceLines laid out for the narrative: the same lines,
// followed by the blank line that separates sections here. A rung carrying no
// evidence contributes nothing rather than a stray blank.
func rungLines(next data.Next) []string {
	lines := EvidenceLines(next)
	if len(lines) == 0 {
		return nil
	}
	return append(lines, "")
}

// EvidenceLines renders the rung-specific evidence: the technical clearance,
// owner-wait, exception, dependency, or replan state that explains why this
// rung — and no other — was reached. Each rung's rendering is distinct, and
// none of it claims anything was verified, run, checked, or confirmed — these
// lines report what a record already says.
//
// It is exported because the V2 board's Next area reports the same facts from
// the same data.Next in a compact layout of its own. The layouts differ; the
// wording must not, so both surfaces read it from here rather than keeping a
// second copy of the evidence vocabulary (STYLE-07, STYLE-09).
func EvidenceLines(next data.Next) []string {
	switch next.Kind {
	case data.NextPendingMigration:
		return []string{
			fmt.Sprintf("Migration: an operation is in progress (%s). Every other action is on hold until it resolves.", next.Migration.OperationID),
		}
	case data.NextReplan:
		return []string{"Replan: " + replanBlockerPhrase(next.GateDecision)}
	case data.NextDependency:
		return dependencyBlockerLines(next.GateDecision)
	case data.NextExecute:
		return executeLines(next.GateDecision, next.Task)
	case data.NextCheckNeeded:
		return []string{"Technical clearance: " + ClearancePhrase(next.Clearance)}
	case data.NextOwnerValidationRequired:
		return []string{
			"Technical clearance: " + ClearancePhrase(next.Clearance),
			"Owner wait: " + ownerWaitPhrase(clearanceCheckID(next.Clearance)),
		}
	case data.NextObjectiveIntegration:
		lines := []string{"Technical clearance: " + ClearancePhrase(next.Clearance)}
		if hasBlockerKind(next.GateDecision, data.GateBlockOwnerAcceptance) {
			lines = append(lines, "Owner wait: "+ownerWaitPhrase(clearanceCheckID(next.Clearance)))
		}
		return lines
	case data.NextReady, data.NextPlanObjective:
		return nil
	default:
		return nil
	}
}

// replanBlockerPhrase reads the recorded replan reason straight off the
// GateDecision's own GateBlockReplan blocker, never restating the flag
// itself.
func replanBlockerPhrase(decision *data.GateDecision) string {
	if decision != nil {
		for _, blocker := range decision.Blockers {
			if blocker.Kind == data.GateBlockReplan {
				return ReplanPhrase(blocker.Detail)
			}
		}
	}
	return ReplanPhrase("")
}

// dependencyBlockerLines renders one "Blocked:" line per unmet Task or
// Objective dependency the GateDecision carries. ResolveTaskStart never
// mixes dependency blockers with any other kind on the same decision, so
// every blocker here is one of the two dependency kinds.
func dependencyBlockerLines(decision *data.GateDecision) []string {
	if decision == nil {
		return []string{"Blocked: a dependency is unsatisfied, with no detail carried on this rung."}
	}
	var lines []string
	for _, blocker := range decision.Blockers {
		switch blocker.Kind {
		case data.GateBlockDependency:
			lines = append(lines, "Blocked: "+DependencyPhrase(blocker.Dependency))
		case data.GateBlockObjectiveDependency:
			lines = append(lines, "Blocked: The owning Objective is waiting: "+ObjectiveDependencyPhrase(blocker.ObjectiveDependency))
		}
	}
	if len(lines) == 0 {
		lines = append(lines, "Blocked: a dependency is unsatisfied, with no detail carried on this rung.")
	}
	return lines
}

// executeLines renders the NextExecute rung: a completion allowed only by a
// recorded exception is reported as exactly that, never as clearance: a
// Task that may start or advance is reported by what allows it.
func executeLines(decision *data.GateDecision, task *data.TaskV2) []string {
	if decision != nil && decision.AllowedByException {
		return []string{"Completion: " + ExceptionPhrase(decision.Exception)}
	}
	return []string{"Ready: " + executeReadyPhrase(task)}
}

func executeReadyPhrase(task *data.TaskV2) string {
	switch {
	case task.Status == data.ColumnPlanned:
		return "Its dependencies are satisfied; it may start."
	case task.Status == data.ColumnInProgress && task.Stage != data.StageAudit:
		return "It may advance to its next stage."
	default:
		return "Its technical clearance is current; it may complete."
	}
}

// clearanceCheckID reads the Check ID a resolved Clearance names, or "" when
// no clearance was carried on this rung.
func clearanceCheckID(clearance *data.Clearance) string {
	if clearance == nil {
		return ""
	}
	return clearance.Check
}

// hasBlockerKind reports whether decision carries a blocker of kind.
func hasBlockerKind(decision *data.GateDecision, kind data.GateBlockKind) bool {
	if decision == nil {
		return false
	}
	for _, blocker := range decision.Blockers {
		if blocker.Kind == kind {
			return true
		}
	}
	return false
}

// issueLines renders the Issues section, omitted entirely rather than
// printed empty when next carries none.
func issueLines(issues []*data.IssueV2) []string {
	if len(issues) == 0 {
		return nil
	}
	lines := make([]string, 0, len(issues)+2)
	lines = append(lines, "Issues:")
	for _, issue := range issues {
		lines = append(lines, IssueLine(issue))
	}
	return append(lines, "")
}

// ActionPhrase is the one closing instruction every rendering ends with:
// what a user or agent reading this output should actually do next. It is
// exported for the same reason EvidenceLines is — the board's Next area
// states the same action, and two surfaces phrasing one answer differently
// is the divergence the shared projection exists to prevent.
func ActionPhrase(next data.Next) string {
	switch next.Kind {
	case data.NextPendingMigration:
		return "Wait for the pending migration to finish; nothing else is actionable until it resolves."
	case data.NextReplan:
		return "Resolve the recorded replan before resuming build, test, or audit."
	case data.NextDependency:
		return "Wait on the named dependency before starting or advancing this Task."
	case data.NextExecute:
		return executeNextActionPhrase(next)
	case data.NextCheckNeeded:
		return "Record a fresh Check against this target."
	case data.NextOwnerValidationRequired:
		return "Ask the owner to accept the current Check."
	case data.NextObjectiveIntegration:
		return objectiveIntegrationNextActionPhrase(next)
	case data.NextReady:
		return readyNextActionPhrase(next)
	case data.NextPlanObjective:
		return "Plan the next Objective — for a project with nothing underway yet, start with Idea/Design."
	default:
		return fmt.Sprintf("No next action is defined for rung %q.", next.Kind)
	}
}

func executeNextActionPhrase(next data.Next) string {
	if next.GateDecision != nil && next.GateDecision.AllowedByException {
		return "Proceed under the recorded exception; no further clearance is required."
	}
	switch {
	case next.Task.Status == data.ColumnPlanned:
		return fmt.Sprintf("Start Task %s.", next.Task.ID)
	case next.Task.Status == data.ColumnInProgress && next.Task.Stage != data.StageAudit:
		return fmt.Sprintf("Advance Task %s to its next stage.", next.Task.ID)
	default:
		return fmt.Sprintf("Complete Task %s.", next.Task.ID)
	}
}

func objectiveIntegrationNextActionPhrase(next data.Next) string {
	if hasBlockerKind(next.GateDecision, data.GateBlockOwnerAcceptance) {
		return "Ask the owner to accept the Objective's current integration Check."
	}
	return fmt.Sprintf("Record the Objective %s integration Check.", next.Objective.ID)
}

func readyNextActionPhrase(next data.Next) string {
	if next.Task != nil {
		return fmt.Sprintf("Start Task %s.", next.Task.ID)
	}
	return fmt.Sprintf("Plan Tasks under Objective %s.", next.Objective.ID)
}
