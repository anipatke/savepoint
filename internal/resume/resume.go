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
// plan-next-Objective carries neither.
func identityLines(next data.Next) []string {
	lines := ReleaseIdentityLines(next)
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
			"Owner wait: " + OwnerWaitPhrase(clearanceCheckID(next.Clearance)),
		}
	case data.NextObjectiveIntegration:
		lines := []string{"Technical clearance: " + ClearancePhrase(next.Clearance)}
		if hasBlockerKind(next.GateDecision, data.GateBlockOwnerAcceptance) {
			lines = append(lines, "Owner wait: "+OwnerWaitPhrase(clearanceCheckID(next.Clearance)))
		}
		return lines
	case data.NextReleaseIntegration:
		return releaseIntegrationLines(next)
	case data.NextReleaseCheckNeeded:
		return []string{"Technical clearance: " + ClearancePhrase(next.Clearance)}
	case data.NextReleaseOwnerValidationRequired:
		return []string{
			"Technical clearance: " + ClearancePhrase(next.Clearance),
			"Owner wait: " + OwnerWaitPhrase(clearanceCheckID(next.Clearance)),
		}
	case data.NextReleaseReady:
		return releaseReadyLines(next)
	case data.NextReady, data.NextPlanObjective:
		return nil
	default:
		return nil
	}
}

// releaseIntegrationLines reports the typed blockers that keep a selected
// Release from reaching its Check or completion rung. Release membership and
// gate decisions were resolved before this function runs; it only words those
// values and never reaches back to an index.
func releaseIntegrationLines(next data.Next) []string {
	if next.GateDecision == nil || len(next.GateDecision.Blockers) == 0 {
		return []string{"Goal readiness: integration is not complete."}
	}

	lines := make([]string, 0, len(next.GateDecision.Blockers))
	for _, blocker := range next.GateDecision.Blockers {
		if blocker.Kind == data.GateBlockOwnerAcceptance {
			lines = append(lines, "Owner wait: "+OwnerWaitPhrase(clearanceCheckID(next.Clearance)))
			continue
		}
		lines = append(lines, "Goal readiness: "+ReleaseBlockerPhrase(blocker))
	}
	return lines
}

// releaseReadyLines keeps three allowed outcomes visibly separate: a current
// V2 Check accepted by the owner, an owner exception, and migrated historical
// completion. None of those states is collapsed into a generic "ready" line.
func releaseReadyLines(next data.Next) []string {
	if next.GateDecision != nil {
		switch {
		case next.GateDecision.AllowedByLegacyCompletion:
			return []string{HistoricalCompletionPhrase(releaseID(next), next.GateDecision.LegacyCompletion)}
		case next.GateDecision.AllowedByException:
			return []string{"Completion: " + ExceptionPhrase(next.GateDecision.Exception)}
		}
	}
	return []string{
		"Technical clearance: " + ClearancePhrase(next.Clearance),
		ReleaseAcceptancePhrase(next.Release, next.Clearance),
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
	case data.NextReleaseIntegration:
		return releaseIntegrationNextActionPhrase(next)
	case data.NextReleaseCheckNeeded:
		if id := releaseID(next); id != "" {
			return fmt.Sprintf("Record a fresh Goal Check for %s.", id)
		}
		return "Record a fresh Goal Check."
	case data.NextReleaseOwnerValidationRequired:
		return "Ask the owner to accept the current Goal Check."
	case data.NextReleaseReady:
		return releaseReadyNextActionPhrase(next)
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

func releaseIntegrationNextActionPhrase(next data.Next) string {
	id := releaseID(next)
	var hasObjectiveBlock, hasIssueBlock bool
	for _, blocker := range blockers(next.GateDecision) {
		switch blocker.Kind {
		case data.GateBlockReleaseNoObjectives, data.GateBlockReleaseObjectiveIncomplete:
			hasObjectiveBlock = true
		case data.GateBlockReleaseIssueUnresolved:
			hasIssueBlock = true
		}
	}

	switch {
	case hasObjectiveBlock:
		if id != "" {
			return fmt.Sprintf("Complete the member Objectives of Goal %s before recording its Goal Check.", id)
		}
		return "Complete the member Objectives before recording the Goal Check."
	case hasIssueBlock:
		if id != "" {
			return fmt.Sprintf("Resolve the material Issues blocking Goal %s before completing it.", id)
		}
		return "Resolve the material Issues blocking the Goal before completing it."
	default:
		if id != "" {
			return fmt.Sprintf("Resolve the integration blockers for Goal %s.", id)
		}
		return "Resolve the Goal integration blockers."
	}
}

func releaseReadyNextActionPhrase(next data.Next) string {
	id := releaseID(next)
	if next.GateDecision != nil {
		if next.GateDecision.AllowedByLegacyCompletion {
			if id != "" {
				return fmt.Sprintf("Review the archived historical completion for Goal %s; it is already marked done.", id)
			}
			return "Review the archived historical Goal completion; it is already marked done."
		}
		if next.GateDecision.AllowedByException {
			if id != "" {
				return fmt.Sprintf("Record Goal %s as done under the recorded exception.", id)
			}
			return "Record the Goal as done under the recorded exception."
		}
	}
	if next.Release != nil && next.Release.Status == data.ColumnDone {
		if id != "" {
			return fmt.Sprintf("Goal %s is already done; no further Goal action is required.", id)
		}
		return "The Goal is already done; no further Goal action is required."
	}
	if id != "" {
		return fmt.Sprintf("Record Goal %s as done.", id)
	}
	return "Record the Goal as done."
}

func blockers(decision *data.GateDecision) []data.GateBlocker {
	if decision == nil {
		return nil
	}
	return decision.Blockers
}

func releaseID(next data.Next) string {
	if next.Release == nil {
		return ""
	}
	return next.Release.ID
}

func readyNextActionPhrase(next data.Next) string {
	if next.Task != nil {
		return fmt.Sprintf("Start Task %s.", next.Task.ID)
	}
	return fmt.Sprintf("Plan Tasks under Objective %s.", next.Objective.ID)
}
