package data

import (
	"fmt"
	"slices"
	"strings"
)

// DecisionApplicability is whether an owner decision still holds at a target's
// latest Check.
type DecisionApplicability string

const (
	// DecisionApplies means the decision was recorded against the latest
	// Check, a checker assessed it as still applying there, or the owner
	// renewed it there.
	DecisionApplies DecisionApplicability = "applies"
	// DecisionUnassessed means no one has assessed the decision at the latest
	// Check.
	DecisionUnassessed DecisionApplicability = "unassessed"
	// DecisionChanged means a checker assessed that a material change ended
	// the decision at the latest Check and the owner has not renewed it.
	DecisionChanged DecisionApplicability = "changed"
)

// DecisionKind names which owner decision a blocker is about, so a change to
// one decision never reads as a change to the other.
type DecisionKind string

const (
	DecisionKindAcceptance DecisionKind = "acceptance"
	DecisionKindException  DecisionKind = "exception"
)

// DecisionAssessment is the resolved applicability of one decision at a
// target's latest Check.
type DecisionAssessment struct {
	State          DecisionApplicability
	Check          string         // the latest Check the decision was assessed at
	Carry          *DecisionCarry // deciding entry; nil when recorded against Check itself
	MaterialChange string         // set only for DecisionChanged
}

// assessDecision is the one rule for whether a decision recorded against
// origin still applies at latest. An owner renewal at latest outranks a
// checker entry for the same Check, so an owner can always reinstate a
// decision a checker found changed.
func assessDecision(origin string, carried []DecisionCarry, latest string) DecisionAssessment {
	assessment := DecisionAssessment{State: DecisionUnassessed, Check: latest}
	if latest == "" {
		return assessment
	}
	if origin == latest {
		assessment.State = DecisionApplies
		return assessment
	}

	var checker *DecisionCarry
	for i := range carried {
		carry := &carried[i]
		if carry.Check != latest {
			continue
		}
		if carry.AssessedBy.Role == ActorRoleOwner {
			assessment.State = DecisionApplies
			assessment.Carry = carry
			return assessment
		}
		checker = carry
	}
	if checker == nil {
		return assessment
	}

	assessment.Carry = checker
	if checker.Applies {
		assessment.State = DecisionApplies
	} else {
		assessment.State = DecisionChanged
		assessment.MaterialChange = checker.MaterialChange
	}
	return assessment
}

// recordedAcceptance returns evidence's owner acceptance when the owner has
// recorded one against a Check.
func recordedAcceptance(evidence *Evidence) *OwnerValidation {
	if evidence == nil || evidence.OwnerValidation == nil {
		return nil
	}
	accepted := evidence.OwnerValidation
	if accepted.AcceptedCheck == "" || accepted.AcceptedBy.Role != ActorRoleOwner || strings.TrimSpace(accepted.AcceptedBy.Session) == "" {
		return nil
	}
	return accepted
}

// acceptanceApplies reports whether the owner's recorded acceptance holds at
// latest.
func acceptanceApplies(evidence *Evidence, latest string) bool {
	accepted := recordedAcceptance(evidence)
	return accepted != nil && assessDecision(accepted.AcceptedCheck, accepted.CarriedForward, latest).State == DecisionApplies
}

// acceptanceBlockers explains why a required owner acceptance does not hold at
// latest, or returns nil when it does.
func acceptanceBlockers(evidence *Evidence, latest string) []GateBlocker {
	accepted := recordedAcceptance(evidence)
	if accepted == nil {
		return []GateBlocker{{Kind: GateBlockOwnerAcceptance, Detail: fmt.Sprintf("owner has not accepted current check %s", latest)}}
	}
	return decisionBlockers(DecisionKindAcceptance, accepted.AcceptedCheck, assessDecision(accepted.AcceptedCheck, accepted.CarriedForward, latest))
}

func decisionBlockers(kind DecisionKind, origin string, assessment DecisionAssessment) []GateBlocker {
	switch assessment.State {
	case DecisionUnassessed:
		return []GateBlocker{{
			Kind:     GateBlockDecisionUnassessed,
			Decision: kind,
			Detail:   fmt.Sprintf("owner %s recorded at %s has not been assessed at latest check %s", kind, origin, assessment.Check),
		}}
	case DecisionChanged:
		return []GateBlocker{{
			Kind:     GateBlockDecisionChanged,
			Decision: kind,
			Change:   assessment.MaterialChange,
			Detail:   fmt.Sprintf("owner %s recorded at %s no longer applies at latest check %s: %s", kind, origin, assessment.Check, assessment.MaterialChange),
		}}
	}
	return nil
}

// exceptionAssessment is an exception's applicability at its target's latest
// Check together with the unmet requirements it does not cover.
type exceptionAssessment struct {
	Exception *Exception
	DecisionAssessment
	// Uncovered lists requirement IDs the latest Check found unmet that the
	// exception does not name. Empty when the Check lists none.
	Uncovered []string
}

// assessException resolves evidence's exception at targetID's latest Check. It
// returns nil when no exception is recorded or the target has no Check.
func assessException(index *V2Index, evidence *Evidence, targetID string) *exceptionAssessment {
	latest := index.LatestCheck[targetID]
	if evidence == nil || evidence.Exception == nil || latest == "" {
		return nil
	}
	exception := evidence.Exception
	assessed := &exceptionAssessment{
		Exception:          exception,
		DecisionAssessment: assessDecision(exception.Check, exception.CarriedForward, latest),
	}
	if check := index.Checks[latest]; check != nil {
		for _, id := range check.Unmet {
			if !slices.Contains(exception.Requirements, id) {
				assessed.Uncovered = append(assessed.Uncovered, id)
			}
		}
	}
	return assessed
}

// grants reports whether the exception allows completion: it applies at the
// latest Check and covers every requirement that Check lists as unmet.
func (a *exceptionAssessment) grants() bool {
	return a != nil && a.State == DecisionApplies && len(a.Uncovered) == 0
}

// blockers explains why the exception does not grant completion.
func (a *exceptionAssessment) blockers() []GateBlocker {
	blockers := decisionBlockers(DecisionKindException, a.Exception.Check, a.DecisionAssessment)
	if a.State == DecisionApplies && len(a.Uncovered) > 0 {
		blockers = append(blockers, GateBlocker{
			Kind:         GateBlockExceptionScope,
			Decision:     DecisionKindException,
			Requirements: a.Uncovered,
			Detail:       fmt.Sprintf("latest check %s lists unmet requirements the exception does not cover: %s", a.Check, strings.Join(a.Uncovered, ", ")),
		})
	}
	return blockers
}

// exceptionGrants reports whether targetID's recorded exception allows
// completion at its latest Check.
func exceptionGrants(index *V2Index, evidence *Evidence, targetID string) bool {
	return assessException(index, evidence, targetID).grants()
}

// exceptionDecision returns the GateDecision for an exception that grants
// completion.
func (a *exceptionAssessment) decision() GateDecision {
	return GateDecision{Allowed: true, Actor: ActorRoleOwner, AllowedByException: true, Exception: a.Exception, ExceptionCarry: a.Carry}
}
