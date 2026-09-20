package data

import (
	"fmt"
	"strings"
	"time"
)

// FreshnessState is the recorded clearance freshness a Task or Objective
// evidence block asserts. There is no default: an absent assessment is a nil
// Freshness, never a defaulted "unknown" value.
type FreshnessState string

const (
	FreshnessCurrent FreshnessState = "current"
	FreshnessStale   FreshnessState = "stale"
	FreshnessUnknown FreshnessState = "unknown"
)

// Freshness is a recorded assessment naming the Check it evaluates, who
// assessed it and when, and the basis for the assessment. It is decoded only
// when the freshness block is present, and only when every required field is
// present and valid — a partially filled block is a diagnostic, never a
// healed value.
type Freshness struct {
	State      FreshnessState
	Check      string // C### this assessment names
	AssessedBy Actor
	AssessedAt time.Time
	Basis      string
}

// OwnerValidation records whether a Task or Objective declares that owner
// acceptance is required for completion, and which Check the owner has
// accepted, if any. An absent OwnerValidation block means owner validation
// is not required.
type OwnerValidation struct {
	Required      bool
	AcceptedCheck string // optional C### reference
	AcceptedBy    Actor  // required when AcceptedCheck is present
}

// Exception records an owner's decision to accept completion despite unmet
// requirements. It is reported as completion by exception, never as a CLEAR
// clearance result.
type Exception struct {
	Requirements []string // non-empty list of requirement IDs
	Reason       string
	Owner        string
	RecordedAt   time.Time
	Check        string // C### this exception applies to
}

// Replan flags that a Task or Objective's plan needs revisiting. It never
// changes the record's status or stage on its own; it is a signal consumers
// read to block start/advance decisions.
type Replan struct {
	Reason     string
	RecordedBy Actor
	RecordedAt time.Time
}

// CheckWaiver records an owner's explicit decision to skip the optional Task
// Check for one Task. It waives only that local Check: it is not technical
// CLEAR, never satisfies a dependency that explicitly requires clear, and
// never substitutes for the mandatory Full Objective Check or a Release
// Check. It applies only to a Task's own evidence — decodeCheckWaiverV2
// rejects one recorded on an Objective or Release.
type CheckWaiver struct {
	Task       string // T### this waiver names; must equal the owning Task's own ID
	Reason     string
	Actor      Actor // required role: owner
	RecordedAt time.Time
}

// Evidence is the shared freshness, acceptance, exception, replan, and
// check-waiver evidence block decoded identically for Task and Objective
// records by decodeEvidenceV2. Every field is independently optional; a
// record carrying none of them decodes to a nil Evidence rather than a
// defaulted one.
type Evidence struct {
	LastCheck       string // optional C### reference to the most recent Check
	Freshness       *Freshness
	OwnerValidation *OwnerValidation
	Exception       *Exception
	Replan          *Replan
	CheckWaiver     *CheckWaiver
}

type evidenceActorFrontmatter struct {
	Role    string `yaml:"role"`
	Session string `yaml:"session"`
}

type freshnessV2Frontmatter struct {
	State      string                   `yaml:"state"`
	Check      string                   `yaml:"check"`
	AssessedBy evidenceActorFrontmatter `yaml:"assessed_by"`
	AssessedAt string                   `yaml:"assessed_at"`
	Basis      string                   `yaml:"basis"`
}

type ownerValidationV2Frontmatter struct {
	Required      bool                      `yaml:"required"`
	AcceptedCheck string                    `yaml:"accepted_check"`
	AcceptedBy    *evidenceActorFrontmatter `yaml:"accepted_by,omitempty"`
}

type exceptionV2Frontmatter struct {
	Requirements []string `yaml:"requirements"`
	Reason       string   `yaml:"reason"`
	Owner        string   `yaml:"owner"`
	RecordedAt   string   `yaml:"recorded_at"`
	Check        string   `yaml:"check"`
}

type replanV2Frontmatter struct {
	Reason     string                   `yaml:"reason"`
	RecordedBy evidenceActorFrontmatter `yaml:"recorded_by"`
	RecordedAt string                   `yaml:"recorded_at"`
}

type checkWaiverV2Frontmatter struct {
	Task       string                   `yaml:"task"`
	Reason     string                   `yaml:"reason"`
	Actor      evidenceActorFrontmatter `yaml:"actor"`
	RecordedAt string                   `yaml:"recorded_at"`
}

// evidenceV2Frontmatter is the shared block of raw evidence fields, embedded
// inline into both taskV2Frontmatter and objectiveV2Frontmatter so the two
// records expose identical YAML shape and decode through the one
// decodeEvidenceV2 decoder below.
type evidenceV2Frontmatter struct {
	LastCheck       string                        `yaml:"last_check"`
	Freshness       *freshnessV2Frontmatter       `yaml:"freshness"`
	OwnerValidation *ownerValidationV2Frontmatter `yaml:"owner_validation"`
	Exception       *exceptionV2Frontmatter       `yaml:"exception"`
	Replan          *replanV2Frontmatter          `yaml:"replan"`
	CheckWaiver     *checkWaiverV2Frontmatter     `yaml:"check_waiver"`
}

// decodeEvidenceV2 strictly decodes the shared evidence block for a Task or
// Objective record. recordKind ("task" or "objective") and id name the
// owning record in diagnostics. Every sub-block is independently optional; a
// record carrying none of last_check, freshness, owner_validation,
// exception, replan, or check_waiver returns a nil Evidence. A present
// sub-block must be fully and validly filled — an unknown state, a
// malformed timestamp, or a partially filled block is a named diagnostic,
// never a healed value.
func decodeEvidenceV2(path, recordKind, id string, raw evidenceV2Frontmatter) (*Evidence, error) {
	if raw.LastCheck == "" && raw.Freshness == nil && raw.OwnerValidation == nil && raw.Exception == nil && raw.Replan == nil && raw.CheckWaiver == nil {
		return nil, nil
	}

	evidence := &Evidence{}

	if raw.LastCheck != "" {
		if !checkIDPatternV2.MatchString(raw.LastCheck) {
			return nil, fmt.Errorf("%w: %s: %s %s last_check %q must match C plus at least three digits", ErrV2InvalidID, path, recordKind, id, raw.LastCheck)
		}
		evidence.LastCheck = raw.LastCheck
	}

	if raw.Freshness != nil {
		freshness, err := decodeFreshnessV2(path, recordKind, id, *raw.Freshness)
		if err != nil {
			return nil, err
		}
		evidence.Freshness = freshness
	}

	if raw.OwnerValidation != nil {
		ownerValidation, err := decodeOwnerValidationV2(path, recordKind, id, *raw.OwnerValidation)
		if err != nil {
			return nil, err
		}
		evidence.OwnerValidation = ownerValidation
	}

	if raw.Exception != nil {
		exception, err := decodeExceptionV2(path, recordKind, id, *raw.Exception)
		if err != nil {
			return nil, err
		}
		evidence.Exception = exception
	}

	if raw.Replan != nil {
		replan, err := decodeReplanV2(path, recordKind, id, *raw.Replan)
		if err != nil {
			return nil, err
		}
		evidence.Replan = replan
	}

	if raw.CheckWaiver != nil {
		waiver, err := decodeCheckWaiverV2(path, recordKind, id, *raw.CheckWaiver)
		if err != nil {
			return nil, err
		}
		evidence.CheckWaiver = waiver
	}

	return evidence, nil
}

func decodeFreshnessV2(path, recordKind, id string, raw freshnessV2Frontmatter) (*Freshness, error) {
	state := FreshnessState(raw.State)
	switch state {
	case "":
		return nil, fmt.Errorf("%w: %s: %s %s missing required field freshness.state", ErrV2MissingField, path, recordKind, id)
	case FreshnessCurrent, FreshnessStale, FreshnessUnknown:
	default:
		return nil, fmt.Errorf("%w: %s: %s %s freshness.state %q; use current, stale, or unknown", ErrV2EvidenceMalformed, path, recordKind, id, raw.State)
	}

	if raw.Check == "" {
		return nil, fmt.Errorf("%w: %s: %s %s missing required field freshness.check", ErrV2MissingField, path, recordKind, id)
	}
	if !checkIDPatternV2.MatchString(raw.Check) {
		return nil, fmt.Errorf("%w: %s: %s %s freshness.check %q must match C plus at least three digits", ErrV2InvalidID, path, recordKind, id, raw.Check)
	}

	assessedBy, err := decodeV2Actor(ErrV2EvidenceMalformed, path, recordKind, id, "freshness.assessed_by", raw.AssessedBy)
	if err != nil {
		return nil, err
	}
	if state == FreshnessCurrent && !isCheckerActor(Actor{Role: assessedBy.Role, Session: assessedBy.Session}) {
		return nil, fmt.Errorf("%w: %s: %s %s current freshness requires freshness.assessed_by.role checker", ErrV2EvidenceMalformed, path, recordKind, id)
	}

	assessedAt, err := decodeV2Timestamp(ErrV2EvidenceMalformed, path, recordKind, id, "freshness.assessed_at", raw.AssessedAt)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(raw.Basis) == "" {
		return nil, fmt.Errorf("%w: %s: %s %s missing required field freshness.basis", ErrV2MissingField, path, recordKind, id)
	}

	return &Freshness{
		State:      state,
		Check:      raw.Check,
		AssessedBy: assessedBy,
		AssessedAt: assessedAt,
		Basis:      raw.Basis,
	}, nil
}

func decodeOwnerValidationV2(path, recordKind, id string, raw ownerValidationV2Frontmatter) (*OwnerValidation, error) {
	if raw.AcceptedCheck != "" && !checkIDPatternV2.MatchString(raw.AcceptedCheck) {
		return nil, fmt.Errorf("%w: %s: %s %s owner_validation.accepted_check %q must match C plus at least three digits", ErrV2InvalidID, path, recordKind, id, raw.AcceptedCheck)
	}

	if raw.AcceptedCheck == "" {
		if raw.AcceptedBy != nil {
			return nil, fmt.Errorf("%w: %s: %s %s owner_validation.accepted_by requires accepted_check", ErrV2EvidenceMalformed, path, recordKind, id)
		}
		return &OwnerValidation{Required: raw.Required}, nil
	}
	if raw.AcceptedBy == nil {
		return nil, fmt.Errorf("%w: %s: %s %s missing required field owner_validation.accepted_by", ErrV2MissingField, path, recordKind, id)
	}
	acceptedBy, err := decodeV2Actor(ErrV2EvidenceMalformed, path, recordKind, id, "owner_validation.accepted_by", *raw.AcceptedBy)
	if err != nil {
		return nil, err
	}
	if acceptedBy.Role != ActorRoleOwner {
		return nil, fmt.Errorf("%w: %s: %s %s owner_validation.accepted_by.role %q; use owner", ErrV2EvidenceMalformed, path, recordKind, id, acceptedBy.Role)
	}

	return &OwnerValidation{Required: raw.Required, AcceptedCheck: raw.AcceptedCheck, AcceptedBy: acceptedBy}, nil
}

func decodeExceptionV2(path, recordKind, id string, raw exceptionV2Frontmatter) (*Exception, error) {
	if len(raw.Requirements) == 0 {
		return nil, fmt.Errorf("%w: %s: %s %s missing required non-empty field exception.requirements", ErrV2MissingField, path, recordKind, id)
	}
	requirements := make([]string, 0, len(raw.Requirements))
	for _, requirement := range raw.Requirements {
		if strings.TrimSpace(requirement) == "" {
			return nil, fmt.Errorf("%w: %s: %s %s exception.requirements contains an empty entry", ErrV2EvidenceMalformed, path, recordKind, id)
		}
		requirements = append(requirements, requirement)
	}

	if strings.TrimSpace(raw.Reason) == "" {
		return nil, fmt.Errorf("%w: %s: %s %s missing required field exception.reason", ErrV2MissingField, path, recordKind, id)
	}
	if strings.TrimSpace(raw.Owner) == "" {
		return nil, fmt.Errorf("%w: %s: %s %s missing required field exception.owner", ErrV2MissingField, path, recordKind, id)
	}

	recordedAt, err := decodeV2Timestamp(ErrV2EvidenceMalformed, path, recordKind, id, "exception.recorded_at", raw.RecordedAt)
	if err != nil {
		return nil, err
	}

	if raw.Check == "" {
		return nil, fmt.Errorf("%w: %s: %s %s missing required field exception.check", ErrV2MissingField, path, recordKind, id)
	}
	if !checkIDPatternV2.MatchString(raw.Check) {
		return nil, fmt.Errorf("%w: %s: %s %s exception.check %q must match C plus at least three digits", ErrV2InvalidID, path, recordKind, id, raw.Check)
	}

	return &Exception{
		Requirements: requirements,
		Reason:       raw.Reason,
		Owner:        raw.Owner,
		RecordedAt:   recordedAt,
		Check:        raw.Check,
	}, nil
}

func decodeReplanV2(path, recordKind, id string, raw replanV2Frontmatter) (*Replan, error) {
	if strings.TrimSpace(raw.Reason) == "" {
		return nil, fmt.Errorf("%w: %s: %s %s missing required field replan.reason", ErrV2MissingField, path, recordKind, id)
	}

	recordedBy, err := decodeV2Actor(ErrV2EvidenceMalformed, path, recordKind, id, "replan.recorded_by", raw.RecordedBy)
	if err != nil {
		return nil, err
	}

	recordedAt, err := decodeV2Timestamp(ErrV2EvidenceMalformed, path, recordKind, id, "replan.recorded_at", raw.RecordedAt)
	if err != nil {
		return nil, err
	}

	return &Replan{Reason: raw.Reason, RecordedBy: recordedBy, RecordedAt: recordedAt}, nil
}

// decodeCheckWaiverV2 decodes an explicit owner waiver of the optional Task
// Check. TEST-09 requires the waiver to name the Task, reason, actor, and
// time; this decoder enforces all four plus the policy boundary that a
// waiver applies only to a Task's own evidence and only under owner
// authority, never planner, executor, or checker self-report.
func decodeCheckWaiverV2(path, recordKind, id string, raw checkWaiverV2Frontmatter) (*CheckWaiver, error) {
	if recordKind != "task" {
		return nil, fmt.Errorf("%w: %s: %s %s check_waiver only applies to a task's evidence", ErrV2EvidenceMalformed, path, recordKind, id)
	}

	if strings.TrimSpace(raw.Task) == "" {
		return nil, fmt.Errorf("%w: %s: %s %s missing required field check_waiver.task", ErrV2MissingField, path, recordKind, id)
	}
	if raw.Task != id {
		return nil, fmt.Errorf("%w: %s: %s %s check_waiver.task %q must name its own task %s", ErrV2EvidenceMalformed, path, recordKind, id, raw.Task, id)
	}

	if strings.TrimSpace(raw.Reason) == "" {
		return nil, fmt.Errorf("%w: %s: %s %s missing required field check_waiver.reason", ErrV2MissingField, path, recordKind, id)
	}

	actor, err := decodeV2Actor(ErrV2EvidenceMalformed, path, recordKind, id, "check_waiver.actor", raw.Actor)
	if err != nil {
		return nil, err
	}
	if actor.Role != ActorRoleOwner {
		return nil, fmt.Errorf("%w: %s: %s %s check_waiver.actor.role %q; use owner", ErrV2EvidenceMalformed, path, recordKind, id, actor.Role)
	}

	recordedAt, err := decodeV2Timestamp(ErrV2EvidenceMalformed, path, recordKind, id, "check_waiver.recorded_at", raw.RecordedAt)
	if err != nil {
		return nil, err
	}

	return &CheckWaiver{Task: raw.Task, Reason: raw.Reason, Actor: actor, RecordedAt: recordedAt}, nil
}

// decodeV2Actor decodes one {role, session} provenance block. malformed is
// the calling record family's own diagnostic, so an Issue's bad actor is
// reported as an Issue problem rather than an evidence one.
func decodeV2Actor(malformed error, path, recordKind, id, fieldName string, raw evidenceActorFrontmatter) (Actor, error) {
	if raw.Role == "" {
		return Actor{}, fmt.Errorf("%w: %s: %s %s missing required field %s.role", ErrV2MissingField, path, recordKind, id, fieldName)
	}
	role := ActorRole(raw.Role)
	switch role {
	case ActorRolePlanner, ActorRoleExecutor, ActorRoleChecker, ActorRoleOwner:
	default:
		return Actor{}, fmt.Errorf("%w: %s: %s %s %s.role %q; use planner, executor, checker, or owner", malformed, path, recordKind, id, fieldName, raw.Role)
	}
	if raw.Session == "" {
		return Actor{}, fmt.Errorf("%w: %s: %s %s missing required field %s.session", ErrV2MissingField, path, recordKind, id, fieldName)
	}
	return Actor{Role: role, Session: raw.Session}, nil
}

// decodeV2Timestamp decodes one required RFC 3339 field, reporting an
// unparseable value under the calling record family's malformed diagnostic.
func decodeV2Timestamp(malformed error, path, recordKind, id, fieldName, raw string) (time.Time, error) {
	if strings.TrimSpace(raw) == "" {
		return time.Time{}, fmt.Errorf("%w: %s: %s %s missing required field %s", ErrV2MissingField, path, recordKind, id, fieldName)
	}
	parsed, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %s: %s %s %s %q is not a parseable RFC 3339 timestamp: %v", malformed, path, recordKind, id, fieldName, raw, err)
	}
	return parsed, nil
}
