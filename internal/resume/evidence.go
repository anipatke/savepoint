// Package resume renders a resolved data.Next projection as read-only
// narrative text, and owns the evidence vocabulary every surface reporting
// that projection speaks. This file holds every piece of evidence,
// freshness, dependency, owner-wait, exception, and Issue phrasing, kept out
// of the layout code in resume.go so each phrase has exactly one source
// (STYLE-07, STYLE-09). EvidenceLines, ActionPhrase, and SelectionPhrase
// expose it to the V2 board's Next area, which lays the same facts out
// compactly rather than restating them; ClearancePhrase, DependencyPhrase,
// ObjectiveDependencyPhrase, ExceptionPhrase, ReplanPhrase, IssueLine, and
// ActorLabel expose the individual phrases to the board's detail surfaces,
// which report one record's own evidence rather than a whole projection.
//
// Every phrase here reports what a record already says. None of them may
// claim that resume verified, ran, checked, or confirmed anything — resume
// read records; it ran nothing (TestNoVerificationClaims in resume_test.go
// enforces this across every rung's output).
package resume

import (
	"fmt"

	"github.com/opencode/savepoint/internal/data"
)

// ActorLabel renders one recorded actor as "<role> session <session>", the
// same provenance shape Evidence and Check records carry. It is exported for
// the same reason the phrases below are: a surface naming who recorded
// something names them in these words, not in a second set.
func ActorLabel(actor data.Actor) string {
	return fmt.Sprintf("%s session %s", actor.Role, actor.Session)
}

// ReleaseIdentityLines names the Release promise carried by a projection.
// Both the board's Next area and resume use these lines so a selected Release
// is accompanied by the same title, outcome, and recorded lifecycle on every
// surface. It returns no lines for the release-free V2 path.
func ReleaseIdentityLines(next data.Next) []string {
	if next.Release == nil {
		return nil
	}
	return []string{
		fmt.Sprintf("Goal: %s — %s", next.Release.ID, next.Release.Title),
		"Goal outcome: " + next.Release.Outcome,
		"Goal status: " + string(next.Release.Status),
	}
}

// ClearancePhrase reports one resolved data.Clearance in wording distinct
// per state, naming the Check when one exists and the freshness assessor,
// date, and recorded basis when a freshness assessment exists. A nil
// clearance reports that no clearance was resolved for this rung, which
// happens when a Task or Objective is allowed to proceed without one (e.g.
// a planned Task with no evidence to speak of yet). It intentionally reports
// the result and freshness authority only: Check.Reviewed is optional scope
// metadata, not a separate requirement for technical CLEAR.
//
// It is exported so a detail surface reporting one record's clearance says
// exactly what the Next area and `savepoint resume` say about the same
// resolved value (STYLE-07, STYLE-09).
func ClearancePhrase(clearance *data.Clearance) string {
	if clearance == nil {
		return "No clearance applies at this step."
	}
	switch clearance.State {
	case data.ClearanceMissing:
		return "No Check has ever been recorded for this target."
	case data.ClearanceNeedsWork:
		return fmt.Sprintf("Check %s recorded NEEDS WORK.", clearance.Check)
	case data.ClearanceStale:
		phrase := fmt.Sprintf("Check %s is recorded CLEAR, but a freshness assessment marks it stale — clearance is stale.", clearance.Check)
		if clearance.Freshness != nil {
			phrase += " " + freshnessBasisPhrase(clearance.Freshness)
		}
		return phrase
	case data.ClearanceUnknown:
		// ResolveClearance reports two different facts as unknown rather than
		// growing a sixth state: a freshness assessment that marks the CLEAR
		// Check unknown, and a CLEAR Check no checker session signed. They are
		// two different things to act on, so they get two different sentences.
		if clearance.Freshness != nil && clearance.Freshness.State == data.FreshnessUnknown {
			return fmt.Sprintf("Check %s is recorded CLEAR, but a freshness assessment marks it unknown — clearance is unknown. %s",
				clearance.Check, freshnessBasisPhrase(clearance.Freshness))
		}
		phrase := fmt.Sprintf("Check %s is recorded CLEAR, but no independent checker session signed it — clearance is not independently established.", clearance.Check)
		if clearance.Freshness != nil {
			phrase += " " + freshnessBasisPhrase(clearance.Freshness)
		}
		return phrase
	case data.ClearanceCurrent:
		phrase := fmt.Sprintf("Check %s is recorded CLEAR.", clearance.Check)
		if clearance.Freshness != nil {
			phrase += " " + freshnessBasisPhrase(clearance.Freshness)
		}
		return phrase
	default:
		return fmt.Sprintf("Clearance state %q is not one this rendering recognizes.", clearance.State)
	}
}

// freshnessBasisPhrase names who assessed a freshness record, when, and on
// what recorded basis — the report of what was written down, never a claim
// that resume itself assessed anything.
func freshnessBasisPhrase(freshness *data.Freshness) string {
	return fmt.Sprintf("Assessed %s by %s on %s, basis: %s.",
		freshness.State, ActorLabel(freshness.AssessedBy), freshness.AssessedAt.Format("2006-01-02"), freshness.Basis)
}

// OwnerWaitPhrase reports an outstanding owner-acceptance block. It is
// exported because Release detail uses the same sentence as the shared Next
// projection and resume output.
func OwnerWaitPhrase(checkID string) string {
	if checkID == "" {
		return "Owner acceptance is required and has not been recorded."
	}
	return fmt.Sprintf("Owner acceptance is required: the owner has not yet accepted Check %s.", checkID)
}

// ExceptionPhrase reports a recorded exception as itself — a reason and an
// owner, never a clearance result, per Exception's own documented contract.
// A detail surface adds the requirement IDs and the time from the record's own
// fields; the claim about what an exception means is made here alone.
func ExceptionPhrase(exception *data.Exception) string {
	if exception == nil {
		return "Allowed by a recorded exception, but no exception detail was carried on this rung."
	}
	return fmt.Sprintf("Allowed by exception, not by clearance: recorded by owner %s for Check %s — %s", exception.Owner, exception.Check, exception.Reason)
}

// ReplanPhrase reports a recorded replan flag by its own reason.
func ReplanPhrase(reason string) string {
	if reason == "" {
		return "A replan has been flagged, with no reason recorded."
	}
	return fmt.Sprintf("A replan has been flagged: %s", reason)
}

// DependencyPhrase reports why one Task dependency is unsatisfied, in
// wording naming the target and the specific unmet requirement rather than
// a bare "blocked".
func DependencyPhrase(block *data.DependencyBlock) string {
	if block == nil {
		return "A dependency is unsatisfied, with no detail carried on this rung."
	}
	switch block.Kind {
	case data.DependencyBlockNotDone:
		return fmt.Sprintf("Waiting on Task %s, which is not done yet.", block.Target)
	case data.DependencyBlockNotCleared:
		return fmt.Sprintf("Waiting on Task %s, which is done but its clearance is %s, not current.", block.Target, block.Clearance)
	case data.DependencyBlockNotAccepted:
		return fmt.Sprintf("Waiting on Task %s, which is done and currently cleared but not yet accepted by the owner.", block.Target)
	default:
		return fmt.Sprintf("Waiting on Task %s for an unrecognized reason %q.", block.Target, block.Kind)
	}
}

// ObjectiveDependencyPhrase mirrors DependencyPhrase for a blocked Objective
// dependency. It is stated about the dependency rather than about whoever is
// waiting on it: a Task rung waits on its owning Objective's dependencies and
// an Objective's own detail waits on its own, and the reason is the same fact
// in both. The subject belongs to the caller's line; the reason lives here
// once (STYLE-07).
func ObjectiveDependencyPhrase(block *data.ObjectiveDependencyBlock) string {
	if block == nil {
		return "an Objective dependency is unsatisfied, with no detail recorded."
	}
	switch block.Kind {
	case data.ObjectiveDependencyBlockNotDone:
		return fmt.Sprintf("Objective %s is not done yet.", block.Target)
	case data.ObjectiveDependencyBlockNotCleared:
		return fmt.Sprintf("Objective %s is done but its clearance is %s, not current.", block.Target, block.Clearance)
	case data.ObjectiveDependencyBlockClearedByException:
		return fmt.Sprintf("Objective %s reached done only by exception, not current clearance.", block.Target)
	default:
		return fmt.Sprintf("Objective %s is unsatisfied for an unrecognized reason %q.", block.Target, block.Kind)
	}
}

// IssueLine renders one Issue as its ID, type, status, and title — the four
// facts an Issue listing promises and nothing derived beyond them.
func IssueLine(issue *data.IssueV2) string {
	return fmt.Sprintf("- %s (%s, %s): %s", issue.ID, issue.Type, issue.Status, issue.Title)
}

// ReleaseBlockerPhrase turns one typed Release gate blocker into readable
// evidence. The blocker was already resolved by internal/data; this function
// only gives that value one shared sentence for resume, the Next area, and
// Release detail.
func ReleaseBlockerPhrase(blocker data.GateBlocker) string {
	switch blocker.Kind {
	case data.GateBlockReleaseNoObjectives:
		return "No member Objective is recorded."
	case data.GateBlockReleaseObjectiveIncomplete:
		if blocker.Objective != "" {
			return fmt.Sprintf("Member Objective %s is not complete: %s", blocker.Objective, blocker.Detail)
		}
		return "A member Objective is not complete: " + blocker.Detail
	case data.GateBlockReleaseIssueUnresolved:
		if blocker.Issue != "" {
			return fmt.Sprintf("Issue %s remains unresolved: %s", blocker.Issue, blocker.Detail)
		}
		return "A material Issue remains unresolved: " + blocker.Detail
	case data.GateBlockOwnerAcceptance:
		if blocker.Detail != "" {
			return "Owner acceptance is required: " + blocker.Detail
		}
		return OwnerWaitPhrase("")
	default:
		if blocker.Detail != "" {
			return "Goal completion is blocked: " + blocker.Detail
		}
		return fmt.Sprintf("Goal completion is blocked by %s.", blocker.Kind)
	}
}

// HistoricalCompletionPhrase keeps migrated historical proof distinct from a
// current V2 Check. The archive path is evidence's recorded destination, not
// a claim that resume opened or verified it.
func HistoricalCompletionPhrase(releaseID string, reference *data.LegacyCompletionReference) string {
	if reference == nil {
		return fmt.Sprintf("Historical completion: Goal %s is recorded done from archived legacy evidence; it is not a new V2 CLEAR Check.", releaseID)
	}
	return fmt.Sprintf("Historical completion: Goal %s is recorded done from archived legacy evidence at %s; it is not a new V2 CLEAR Check.", releaseID, reference.ArchivePath)
}

// ReleaseAcceptancePhrase reports the owner decision that permits a current
// Release Check to close the Release promise. It is deliberately separate
// from ClearancePhrase: technical currentness and product acceptance are two
// different facts even when the gate has allowed completion.
func ReleaseAcceptancePhrase(release *data.ReleaseV2, clearance *data.Clearance) string {
	if release == nil {
		return ReleaseAcceptancePhraseForEvidence("", nil, clearance)
	}
	return ReleaseAcceptancePhraseForEvidence(release.ID, release.Evidence, clearance)
}

// ReleaseAcceptancePhraseForEvidence is the same Release acceptance wording
// for a resolved detail, which carries the record's evidence fields rather
// than the full Release pointer.
func ReleaseAcceptancePhraseForEvidence(releaseID string, evidence *data.Evidence, clearance *data.Clearance) string {
	if clearance == nil || clearance.Check == "" {
		if releaseID != "" {
			return fmt.Sprintf("Goal readiness: Goal %s completion is allowed by the recorded Goal decision.", releaseID)
		}
		return "Goal readiness: completion is allowed by the recorded Goal decision."
	}
	if evidence != nil && evidence.OwnerValidation != nil {
		accepted := evidence.OwnerValidation
		if accepted.AcceptedCheck == clearance.Check {
			return fmt.Sprintf("Goal readiness: Check %s is current and accepted by %s.", clearance.Check, ActorLabel(accepted.AcceptedBy))
		}
	}
	return fmt.Sprintf("Goal readiness: Check %s is current and the recorded owner decision allows completion.", clearance.Check)
}

// implementationPhrase reports a Task's recorded status and stage as a
// plain fact, never as a judgement about whether that status is correct.
func implementationPhrase(task *data.TaskV2) string {
	if task.Status == data.ColumnInProgress {
		return fmt.Sprintf("Status %s, stage %s.", task.Status, task.Stage)
	}
	return fmt.Sprintf("Status %s.", task.Status)
}

// objectiveStatusPhrase mirrors implementationPhrase for an Objective, which
// carries no stage of its own.
func objectiveStatusPhrase(objective *data.ObjectiveV2) string {
	return fmt.Sprintf("Status %s.", objective.Status)
}

// SelectionPhrase reports an unresolved or stale router selection, naming the
// ID or mismatch it read rather than guessing at a replacement record. It is
// exported alongside EvidenceLines and ActionPhrase so the board's Next area
// names a diagnostic in these same words.
func SelectionPhrase(diagnostic *data.SelectionDiagnostic) string {
	switch diagnostic.Kind {
	case data.SelectionDone:
		switch diagnostic.RecordKind {
		case data.SelectionRecordTask:
			return fmt.Sprintf("Warning: router still selects finished Task %s.", diagnostic.ID)
		case data.SelectionRecordObjective:
			return fmt.Sprintf("Warning: router still selects finished Objective %s.", diagnostic.ID)
		case data.SelectionRecordIssue:
			return fmt.Sprintf("Warning: router still selects resolved Issue %s.", diagnostic.ID)
		default:
			return fmt.Sprintf("Warning: router selection %s is finished.", diagnostic.ID)
		}
	case data.SelectionNotFound:
		return fmt.Sprintf("The router names %s %s, which does not exist among the project's live records.", diagnostic.RecordKind, diagnostic.ID)
	case data.SelectionReleaseMissing:
		return "The router has no Goal selected."
	case data.SelectionMismatch:
		return fmt.Sprintf("The router names Objective %s and Task %s, but Task %s's own record names %s as its owner — the Task record wins, so this selection is not honored.",
			diagnostic.RouterObjective, diagnostic.Task, diagnostic.Task, diagnostic.TaskObjective)
	case data.SelectionReleaseNotFound:
		return fmt.Sprintf("The router names Goal %s, which does not exist among the project's live records.", selectionReleaseID(diagnostic))
	case data.SelectionReleaseArchived:
		return fmt.Sprintf("The router names Goal %s, which is historical and cannot be used as a live delivery context.", selectionReleaseID(diagnostic))
	case data.SelectionObjectiveUnassigned:
		return fmt.Sprintf("The router selects Objective %s inside Goal %s, but the Objective has no Goal reference, so this selection is not honored.", diagnostic.Objective, diagnostic.Release)
	case data.SelectionReleaseMismatch:
		return fmt.Sprintf("The router selects Objective %s inside Goal %s, but its record names Goal %s, so this selection is not honored.", diagnostic.Objective, diagnostic.Release, diagnostic.ObjectiveRelease)
	default:
		return fmt.Sprintf("The router selection did not resolve, for an unrecognized reason %q.", diagnostic.Kind)
	}
}

func selectionReleaseID(diagnostic *data.SelectionDiagnostic) string {
	if diagnostic.Release != "" {
		return diagnostic.Release
	}
	return diagnostic.ID
}
