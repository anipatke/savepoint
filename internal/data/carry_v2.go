package data

import (
	"fmt"
	"strings"
	"time"
)

// DecisionCarry records, at one later Check, whether an owner decision (an
// acceptance or an exception) still applies. A checker entry assesses
// applicability; an owner entry renews the decision and always applies. The
// originating Check stays on the decision as provenance, and entries are only
// ever appended.
type DecisionCarry struct {
	Check          string // C-### the assessment was made at; never the originating Check
	Applies        bool
	AssessedBy     Actor // role checker or owner
	AssessedAt     time.Time
	Reason         string
	MaterialChange string // required when Applies is false
}

type decisionCarryV2Frontmatter struct {
	Check          string                   `yaml:"check"`
	Applies        *bool                    `yaml:"applies"`
	AssessedBy     evidenceActorFrontmatter `yaml:"assessed_by"`
	AssessedAt     string                   `yaml:"assessed_at"`
	Reason         string                   `yaml:"reason"`
	MaterialChange string                   `yaml:"material_change,omitempty"`
}

// decodeDecisionCarriesV2 decodes one decision's carried_forward list. field
// names the list in diagnostics ("owner_validation.carried_forward") and
// originCheck is the decision's own Check, which an entry may not name.
func decodeDecisionCarriesV2(path, recordKind, id, field, originCheck string, raw []decisionCarryV2Frontmatter) ([]DecisionCarry, error) {
	if len(raw) == 0 {
		return nil, nil
	}

	type roleKey struct {
		check string
		role  ActorRole
	}
	seen := make(map[roleKey]bool, len(raw))
	carries := make([]DecisionCarry, 0, len(raw))
	for i, entry := range raw {
		name := fmt.Sprintf("%s[%d]", field, i)

		if entry.Check == "" {
			return nil, fmt.Errorf("%w: %s: %s %s missing required field %s.check", ErrV2MissingField, path, recordKind, id, name)
		}
		if !matchesV2Identity(entry.Check, 'C') {
			return nil, fmt.Errorf("%w: %s: %s %s %s.check %q must match C- plus at least three digits", ErrV2InvalidID, path, recordKind, id, name, entry.Check)
		}
		if entry.Check == originCheck {
			return nil, fmt.Errorf("%w: %s: %s %s %s.check %s is the decision's own check; name a later check", ErrV2EvidenceMalformed, path, recordKind, id, name, entry.Check)
		}

		if entry.Applies == nil {
			return nil, fmt.Errorf("%w: %s: %s %s missing required field %s.applies", ErrV2MissingField, path, recordKind, id, name)
		}
		applies := *entry.Applies

		assessedBy, err := decodeV2Actor(ErrV2EvidenceMalformed, path, recordKind, id, name+".assessed_by", entry.AssessedBy)
		if err != nil {
			return nil, err
		}
		if assessedBy.Role != ActorRoleChecker && assessedBy.Role != ActorRoleOwner {
			return nil, fmt.Errorf("%w: %s: %s %s %s.assessed_by.role %q; use checker or owner", ErrV2EvidenceMalformed, path, recordKind, id, name, assessedBy.Role)
		}
		if assessedBy.Role == ActorRoleOwner && !applies {
			return nil, fmt.Errorf("%w: %s: %s %s %s is an owner entry and must apply; an owner renewal cannot end a decision", ErrV2EvidenceMalformed, path, recordKind, id, name)
		}

		assessedAt, err := decodeV2Timestamp(ErrV2EvidenceMalformed, path, recordKind, id, name+".assessed_at", entry.AssessedAt)
		if err != nil {
			return nil, err
		}

		if strings.TrimSpace(entry.Reason) == "" {
			return nil, fmt.Errorf("%w: %s: %s %s missing required field %s.reason", ErrV2MissingField, path, recordKind, id, name)
		}
		if !applies && strings.TrimSpace(entry.MaterialChange) == "" {
			return nil, fmt.Errorf("%w: %s: %s %s missing required field %s.material_change when applies is false", ErrV2MissingField, path, recordKind, id, name)
		}

		if i > 0 && compareV2CheckIDs(carries[i-1].Check, entry.Check) > 0 {
			return nil, fmt.Errorf("%w: %s: %s %s %s.check %s is out of order; entries run in check ID order", ErrV2EvidenceMalformed, path, recordKind, id, name, entry.Check)
		}
		key := roleKey{entry.Check, assessedBy.Role}
		if seen[key] {
			return nil, fmt.Errorf("%w: %s: %s %s %s repeats the %s entry for check %s", ErrV2EvidenceMalformed, path, recordKind, id, name, assessedBy.Role, entry.Check)
		}
		seen[key] = true

		carries = append(carries, DecisionCarry{
			Check:          entry.Check,
			Applies:        applies,
			AssessedBy:     assessedBy,
			AssessedAt:     assessedAt,
			Reason:         entry.Reason,
			MaterialChange: entry.MaterialChange,
		})
	}
	return carries, nil
}

func decisionCarriesToFrontmatter(carries []DecisionCarry) []decisionCarryV2Frontmatter {
	if len(carries) == 0 {
		return nil
	}
	raw := make([]decisionCarryV2Frontmatter, 0, len(carries))
	for _, carry := range carries {
		applies := carry.Applies
		raw = append(raw, decisionCarryV2Frontmatter{
			Check:          carry.Check,
			Applies:        &applies,
			AssessedBy:     evidenceActorFrontmatter{Role: string(carry.AssessedBy.Role), Session: carry.AssessedBy.Session},
			AssessedAt:     carry.AssessedAt.Format(time.RFC3339),
			Reason:         carry.Reason,
			MaterialChange: carry.MaterialChange,
		})
	}
	return raw
}

// decodeBlankFreeList copies a list of IDs, refusing a blank entry under the
// named field.
func decodeBlankFreeList(path, recordKind, id, field string, raw []string) ([]string, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	out := make([]string, 0, len(raw))
	for _, entry := range raw {
		if strings.TrimSpace(entry) == "" {
			return nil, fmt.Errorf("%w: %s: %s %s %s contains an empty entry", ErrV2EvidenceMalformed, path, recordKind, id, field)
		}
		out = append(out, entry)
	}
	return out, nil
}

// decisionCarryTargets lists each carried Check of a Task, Objective or
// Release decision in declaration order, with the field that names it.
func decisionCarryTargets(evidence *Evidence) []struct{ field, origin, check string } {
	var targets []struct{ field, origin, check string }
	if evidence == nil {
		return targets
	}
	if evidence.OwnerValidation != nil {
		for i, carry := range evidence.OwnerValidation.CarriedForward {
			targets = append(targets, struct{ field, origin, check string }{
				fmt.Sprintf("owner_validation.carried_forward[%d].check", i), evidence.OwnerValidation.AcceptedCheck, carry.Check,
			})
		}
	}
	if evidence.Exception != nil {
		for i, carry := range evidence.Exception.CarriedForward {
			targets = append(targets, struct{ field, origin, check string }{
				fmt.Sprintf("exception.carried_forward[%d].check", i), evidence.Exception.Check, carry.Check,
			})
		}
	}
	return targets
}

// checkDecisionCarryReferences requires every carried Check to exist and to be
// scoped to the same record as the Check the decision originated at.
func checkDecisionCarryReferences(index *V2Index, path, recordKind, id string, evidence *Evidence) error {
	for _, target := range decisionCarryTargets(evidence) {
		carried, ok := index.Checks[target.check]
		if !ok {
			return fmt.Errorf("%w: %s: %s %s %s names missing check %s", ErrV2EvidenceMissingReference, path, recordKind, id, target.field, target.check)
		}
		origin, ok := index.Checks[target.origin]
		if !ok {
			continue // the originating reference is reported by checkEvidenceReferences
		}
		if carried.Scope != origin.Scope {
			return fmt.Errorf("%w: %s: %s %s %s check %s is scoped to %s %s, not %s %s of the originating check %s",
				ErrV2EvidenceMalformed, path, recordKind, id, target.field, target.check,
				carried.Scope.Kind, carried.Scope.ID, origin.Scope.Kind, origin.Scope.ID, target.origin)
		}
	}
	return nil
}
