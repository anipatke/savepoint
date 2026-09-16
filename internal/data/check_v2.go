package data

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// checkIDPatternV2 anchors V2 global Check identity: C plus at least three
// digits.
var checkIDPatternV2 = regexp.MustCompile(`^C[0-9]{3,}$`)

// issueIDPatternV2 anchors the I### identity references a Check's issues
// list may carry. E43 validates the reference shape only: it never loads,
// looks up, or validates an Issue record, which belongs to E44.
var issueIDPatternV2 = regexp.MustCompile(`^I[0-9]{3,}$`)

// CheckResult is the outcome an evaluation records. There is no partial or
// pending result: a Check is written once, after the evaluation concludes.
type CheckResult string

const (
	CheckResultClear     CheckResult = "CLEAR"
	CheckResultNeedsWork CheckResult = "NEEDS WORK"
)

// CheckScopeKind names the record family a Check evaluates.
type CheckScopeKind string

const (
	CheckScopeTask      CheckScopeKind = "task"
	CheckScopeObjective CheckScopeKind = "objective"
)

// CheckScope identifies the single Task or Objective a Check evaluates.
type CheckScope struct {
	Kind CheckScopeKind
	ID   string
}

// ActorRole is the recorded authority behind a Check or evidence entry.
type ActorRole string

const (
	ActorRolePlanner  ActorRole = "planner"
	ActorRoleExecutor ActorRole = "executor"
	ActorRoleChecker  ActorRole = "checker"
	ActorRoleOwner    ActorRole = "owner"
)

// Actor is the provenance of a recorded outcome: who produced it, and which
// session.
type Actor struct {
	Role    ActorRole
	Session string
}

// ReviewedBasis is the recorded scope an evaluation actually covered. An
// absent block, or an absent field within a present block, is decoded as
// absent rather than healed into an empty-but-passing default: absence is a
// fact about what evidence exists, not a value to paper over.
type ReviewedBasis struct {
	BaseCommit   string
	HeadCommit   string
	Files        []string
	Dependencies []string
}

// CheckV2 is a strict, immutable V2 Check record decoded by DecodeCheckV2.
// It is the shared schema for both Task-scoped and Objective-scoped
// evaluations; nothing here is healed, so a malformed field is a named
// diagnostic rather than a silently defaulted value.
type CheckV2 struct {
	ID        string
	Scope     CheckScope
	Result    CheckResult
	CheckedBy Actor
	CheckedAt time.Time
	Reviewed  *ReviewedBasis
	// Issues are I### references decoded for shape here and resolved against
	// the Issue records at index time, where this list is the authoritative
	// record of which Issues the evaluation opened.
	Issues     []string
	Supersedes string // C### reference to the Check this run replaces, empty if none
	Source     V2SourceDocument
}

type checkScopeFrontmatter struct {
	Kind string `yaml:"kind"`
	ID   string `yaml:"id"`
}

type checkActorFrontmatter struct {
	Role    string `yaml:"role"`
	Session string `yaml:"session"`
}

type reviewedFrontmatter struct {
	BaseCommit   string   `yaml:"base_commit"`
	HeadCommit   string   `yaml:"head_commit"`
	Files        []string `yaml:"files"`
	Dependencies []string `yaml:"dependencies"`
}

type checkV2Frontmatter struct {
	ID         string                `yaml:"id"`
	Scope      checkScopeFrontmatter `yaml:"scope"`
	Result     string                `yaml:"result"`
	CheckedBy  checkActorFrontmatter `yaml:"checked_by"`
	CheckedAt  string                `yaml:"checked_at"`
	Reviewed   *reviewedFrontmatter  `yaml:"reviewed"`
	Issues     []string              `yaml:"issues"`
	Supersedes string                `yaml:"supersedes"`
}

// DecodeCheckV2 strictly decodes a V2 Check record from content. It requires
// a valid global C### ID, a scope naming a T### or O### target consistent
// with its kind, a CLEAR or NEEDS WORK result, checked_by actor provenance,
// and a parseable checked_at timestamp. A CLEAR Check must be recorded by a
// checker; other roles may record NEEDS WORK evidence, but cannot author a
// clearance-capable result. reviewed, issues, and supersedes are optional;
// when present they are validated for shape only — reviewed's fields are
// recorded as given, issues are I### references with no Issue lookup, and
// supersedes is a C### reference resolved later against the full index, not
// here.
func DecodeCheckV2(path, content string) (*CheckV2, error) {
	doc, err := ParseV2Document(path, content)
	if err != nil {
		return nil, err
	}

	var fields checkV2Frontmatter
	if err := doc.Frontmatter.Decode(&fields); err != nil {
		return nil, fmt.Errorf("%w: %s: %v", ErrV2Malformed, path, err)
	}

	if !checkIDPatternV2.MatchString(fields.ID) {
		return nil, fmt.Errorf("%w: %s: check id %q must match C plus at least three digits", ErrV2InvalidID, path, fields.ID)
	}

	scope, err := decodeCheckScope(path, fields.ID, fields.Scope)
	if err != nil {
		return nil, err
	}

	result := CheckResult(fields.Result)
	switch result {
	case "":
		return nil, fmt.Errorf("%w: %s: check %s missing required field result", ErrV2MissingField, path, fields.ID)
	case CheckResultClear, CheckResultNeedsWork:
	default:
		return nil, fmt.Errorf("%w: %s: check %s result %q; use CLEAR or NEEDS WORK", ErrV2CheckMalformed, path, fields.ID, fields.Result)
	}

	checkedBy, err := decodeCheckActor(path, fields.ID, fields.CheckedBy)
	if err != nil {
		return nil, err
	}
	if result == CheckResultClear && !isCheckerActor(checkedBy) {
		return nil, fmt.Errorf("%w: %s: check %s CLEAR result requires checked_by.role checker", ErrV2CheckMalformed, path, fields.ID)
	}

	if strings.TrimSpace(fields.CheckedAt) == "" {
		return nil, fmt.Errorf("%w: %s: check %s missing required field checked_at", ErrV2MissingField, path, fields.ID)
	}
	checkedAt, err := time.Parse(time.RFC3339, fields.CheckedAt)
	if err != nil {
		return nil, fmt.Errorf("%w: %s: check %s checked_at %q is not a parseable RFC 3339 timestamp: %v", ErrV2CheckMalformed, path, fields.ID, fields.CheckedAt, err)
	}

	var reviewed *ReviewedBasis
	if fields.Reviewed != nil {
		reviewed = &ReviewedBasis{
			BaseCommit:   fields.Reviewed.BaseCommit,
			HeadCommit:   fields.Reviewed.HeadCommit,
			Files:        fields.Reviewed.Files,
			Dependencies: fields.Reviewed.Dependencies,
		}
	}

	issues := make([]string, 0, len(fields.Issues))
	for _, ref := range fields.Issues {
		if !issueIDPatternV2.MatchString(ref) {
			return nil, fmt.Errorf("%w: %s: check %s issues %q must match I plus at least three digits", ErrV2InvalidID, path, fields.ID, ref)
		}
		issues = append(issues, ref)
	}

	if fields.Supersedes != "" && !checkIDPatternV2.MatchString(fields.Supersedes) {
		return nil, fmt.Errorf("%w: %s: check %s supersedes %q must match C plus at least three digits", ErrV2InvalidID, path, fields.ID, fields.Supersedes)
	}

	return &CheckV2{
		ID:         fields.ID,
		Scope:      scope,
		Result:     result,
		CheckedBy:  checkedBy,
		CheckedAt:  checkedAt,
		Reviewed:   reviewed,
		Issues:     issues,
		Supersedes: fields.Supersedes,
		Source:     doc,
	}, nil
}

func decodeCheckScope(path, checkID string, raw checkScopeFrontmatter) (CheckScope, error) {
	if raw.Kind == "" {
		return CheckScope{}, fmt.Errorf("%w: %s: check %s missing required field scope.kind", ErrV2MissingField, path, checkID)
	}
	kind := CheckScopeKind(raw.Kind)
	if kind != CheckScopeTask && kind != CheckScopeObjective {
		return CheckScope{}, fmt.Errorf("%w: %s: check %s scope.kind %q; use task or objective", ErrV2CheckMalformed, path, checkID, raw.Kind)
	}
	if raw.ID == "" {
		return CheckScope{}, fmt.Errorf("%w: %s: check %s missing required field scope.id", ErrV2MissingField, path, checkID)
	}

	pattern := taskIDPatternV2
	if kind == CheckScopeObjective {
		pattern = objectiveIDPattern
	}
	if !pattern.MatchString(raw.ID) {
		return CheckScope{}, fmt.Errorf("%w: %s: check %s scope.id %q does not match scope.kind %s", ErrV2InvalidID, path, checkID, raw.ID, kind)
	}

	return CheckScope{Kind: kind, ID: raw.ID}, nil
}

func decodeCheckActor(path, checkID string, raw checkActorFrontmatter) (Actor, error) {
	if raw.Role == "" {
		return Actor{}, fmt.Errorf("%w: %s: check %s missing required field checked_by.role", ErrV2MissingField, path, checkID)
	}
	role := ActorRole(raw.Role)
	switch role {
	case ActorRolePlanner, ActorRoleExecutor, ActorRoleChecker, ActorRoleOwner:
	default:
		return Actor{}, fmt.Errorf("%w: %s: check %s checked_by.role %q; use planner, executor, checker, or owner", ErrV2CheckMalformed, path, checkID, raw.Role)
	}
	if raw.Session == "" {
		return Actor{}, fmt.Errorf("%w: %s: check %s missing required field checked_by.session", ErrV2MissingField, path, checkID)
	}
	return Actor{Role: role, Session: raw.Session}, nil
}

// isCheckerActor identifies the only actor provenance that may contribute to
// a clearance decision. The session check keeps an in-memory decision from
// treating an incomplete actor value as independent evidence; decoding also
// enforces the same field as required input.
func isCheckerActor(actor Actor) bool {
	return actor.Role == ActorRoleChecker && strings.TrimSpace(actor.Session) != ""
}
