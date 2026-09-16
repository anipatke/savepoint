package data

import (
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strings"
	"time"
)

// IssueType classifies what kind of durable follow-up an Issue records. It is
// descriptive only: no gate decision reads it, so classifying an Issue never
// blocks or unblocks work by itself. Material blocking is expressed by a Check
// recording NEEDS WORK, not by an Issue existing.
type IssueType string

const (
	IssueTypeDefect       IssueType = "defect"
	IssueTypeDrift        IssueType = "drift"
	IssueTypeGuardrail    IssueType = "guardrail"
	IssueTypeVerification IssueType = "verification"
	IssueTypeOther        IssueType = "other"
)

// IssueStatus is the Issue lifecycle vocabulary, and it is deliberately
// disjoint from the Task lifecycle in task.go: an Issue is never planned or
// done, and it never carries a stage. A Task lifecycle value in an Issue's
// status is a named diagnostic rather than a healed value, so the two
// vocabularies cannot quietly merge into one.
type IssueStatus string

const (
	IssueStatusOpen       IssueStatus = "open"
	IssueStatusInProgress IssueStatus = "in_progress"
	IssueStatusResolved   IssueStatus = "resolved"
)

// issueStatuses and issueTypes are the vocabularies themselves, held once so
// decoding and the derived counts below cannot drift apart. Their order is the
// order counts and listings report in.
var (
	issueStatuses = []IssueStatus{IssueStatusOpen, IssueStatusInProgress, IssueStatusResolved}
	issueTypes    = []IssueType{IssueTypeDefect, IssueTypeDrift, IssueTypeGuardrail, IssueTypeVerification, IssueTypeOther}
)

// IssueOriginKind names what produced an Issue: an evaluation that recorded
// it, a direct report, or a migration from an earlier record family.
type IssueOriginKind string

const (
	IssueOriginCheck     IssueOriginKind = "check"
	IssueOriginReport    IssueOriginKind = "report"
	IssueOriginMigration IssueOriginKind = "migration"
)

// IssueOrigin is the decoded source block: what produced the Issue, who
// recorded it, and when. It is required, so an Issue always names its origin
// rather than appearing from nowhere.
type IssueOrigin struct {
	Kind IssueOriginKind
	// Check is the C### that produced the Issue. It is required when Kind is
	// check and rejected otherwise, so a reported or migrated Issue cannot
	// borrow an evaluation's authority.
	Check string
	Actor Actor
	At    time.Time
}

// IssueDisposition is how a resolved Issue was closed. The three dispositions
// carry different obligations — verified is repair with proof, accepted is an
// owner decision that is not repair, and duplicate points at a canonical
// Issue and proves nothing — which is why closure cannot be claimed by
// setting a status alone.
type IssueDisposition string

const (
	IssueDispositionVerified  IssueDisposition = "verified"
	IssueDispositionAccepted  IssueDisposition = "accepted"
	IssueDispositionDuplicate IssueDisposition = "duplicate"
)

// IssueResolution is the decoded resolution block. Decoding validates its
// shape only: the per-disposition proof obligations, and whether a resolution
// may be present at all for a given status, are cross-record rules resolved
// against the full index rather than here.
type IssueResolution struct {
	Disposition IssueDisposition
	Check       string // optional C### proof reference
	Actor       Actor
	At          time.Time
	Reason      string
}

// IssueHistoryKind names what one history entry records.
type IssueHistoryKind string

const (
	IssueHistoryObserved        IssueHistoryKind = "observed"
	IssueHistoryRepairAttempted IssueHistoryKind = "repair_attempted"
	IssueHistoryRechecked       IssueHistoryKind = "rechecked"
	IssueHistoryDeferred        IssueHistoryKind = "deferred"
	IssueHistoryReopened        IssueHistoryKind = "reopened"
	IssueHistoryOwnerDecision   IssueHistoryKind = "owner_decision"
)

// IssueHistoryEntry is one dated, attributed entry in an Issue's append-only
// history. Deferral and reopening are entries here rather than extra
// lifecycle states, so an Issue that waits or recurs keeps the same identity.
type IssueHistoryEntry struct {
	At    time.Time
	Actor Actor
	Kind  IssueHistoryKind
	Note  string
	Check string // optional C### this entry refers to
}

// IssueV2 is a strict V2 Issue record decoded by DecodeIssueV2: one durable
// follow-up with a global identity that survives repair, recheck, and
// reopening. Unlike a Check it is mutable, but nothing here is healed, so a
// malformed field is a named diagnostic rather than a silently defaulted
// value.
type IssueV2 struct {
	ID     string
	Title  string
	Type   IssueType
	Status IssueStatus
	// Origin decodes the source block. It is named apart from Source so every
	// V2 record family keeps Source as its parsed document.
	Origin IssueOrigin
	// Tasks and Checks are identity references validated for shape here and
	// resolved against the index later.
	Tasks  []string
	Checks []string
	// GuardrailIDs are opaque policy strings. Savepoint never parses a
	// Guardrails file to confirm them: policy owns its own IDs.
	GuardrailIDs []string
	// Severity is policy-owned and therefore opaque here: it is recorded as
	// given, with no vocabulary enforced and no default supplied. Empty means
	// the record declared none.
	Severity    string
	Resolution  *IssueResolution
	DuplicateOf string // optional I### reference to the canonical Issue
	History     []IssueHistoryEntry
	Source      V2SourceDocument
}

type issueOriginFrontmatter struct {
	Kind  string                   `yaml:"kind"`
	Check string                   `yaml:"check"`
	Actor evidenceActorFrontmatter `yaml:"actor"`
	At    string                   `yaml:"at"`
}

type issueResolutionFrontmatter struct {
	Disposition string                   `yaml:"disposition"`
	Check       string                   `yaml:"check"`
	Actor       evidenceActorFrontmatter `yaml:"actor"`
	At          string                   `yaml:"at"`
	Reason      string                   `yaml:"reason"`
}

type issueHistoryFrontmatter struct {
	At    string                   `yaml:"at"`
	Actor evidenceActorFrontmatter `yaml:"actor"`
	Kind  string                   `yaml:"kind"`
	Note  string                   `yaml:"note"`
	Check string                   `yaml:"check"`
}

type issueV2Frontmatter struct {
	ID     string `yaml:"id"`
	Title  string `yaml:"title"`
	Type   string `yaml:"type"`
	Status string `yaml:"status"`
	// Stage is decoded only so it can be rejected: it is a Task lifecycle
	// field that an Issue must never carry.
	Stage        string                      `yaml:"stage"`
	Source       *issueOriginFrontmatter     `yaml:"source"`
	Tasks        []string                    `yaml:"tasks"`
	Checks       []string                    `yaml:"checks"`
	GuardrailIDs []string                    `yaml:"guardrail_ids"`
	Severity     string                      `yaml:"severity"`
	Resolution   *issueResolutionFrontmatter `yaml:"resolution"`
	DuplicateOf  string                      `yaml:"duplicate_of"`
	History      []issueHistoryFrontmatter   `yaml:"history"`
}

// DecodeIssueV2 strictly decodes a V2 Issue record from content. It requires a
// valid global I### ID, a non-empty title, a known type, an Issue-vocabulary
// status, and a source naming what produced the Issue with its actor and
// time. tasks, checks, guardrail_ids, severity, resolution, duplicate_of, and
// history are optional; when present they are validated for shape only —
// tasks and checks are T### and C### references with no record lookup,
// guardrail_ids are opaque policy strings, and the proof obligations a
// resolution carries are resolved later against the full index, not here.
func DecodeIssueV2(path, content string) (*IssueV2, error) {
	doc, err := ParseV2Document(path, content)
	if err != nil {
		return nil, err
	}

	var fields issueV2Frontmatter
	if err := doc.Frontmatter.Decode(&fields); err != nil {
		return nil, fmt.Errorf("%w: %s: %v", ErrV2Malformed, path, err)
	}

	if !issueIDPatternV2.MatchString(fields.ID) {
		return nil, fmt.Errorf("%w: %s: issue id %q must match I plus at least three digits", ErrV2InvalidID, path, fields.ID)
	}

	if strings.TrimSpace(fields.Title) == "" {
		return nil, fmt.Errorf("%w: %s: issue %s missing required field title", ErrV2MissingField, path, fields.ID)
	}

	issueType, err := decodeIssueType(path, fields.ID, fields.Type)
	if err != nil {
		return nil, err
	}

	status, err := decodeIssueStatus(path, fields.ID, fields.Status, fields.Stage)
	if err != nil {
		return nil, err
	}

	origin, err := decodeIssueOrigin(path, fields.ID, fields.Source)
	if err != nil {
		return nil, err
	}

	tasks, err := decodeIssueReferences(path, fields.ID, "tasks", taskIDPatternV2, "T plus at least three digits", fields.Tasks)
	if err != nil {
		return nil, err
	}

	checks, err := decodeIssueReferences(path, fields.ID, "checks", checkIDPatternV2, "C plus at least three digits", fields.Checks)
	if err != nil {
		return nil, err
	}

	guardrailIDs, err := decodeIssueGuardrailIDs(path, fields.ID, fields.GuardrailIDs)
	if err != nil {
		return nil, err
	}

	if fields.Severity != "" && strings.TrimSpace(fields.Severity) == "" {
		return nil, fmt.Errorf("%w: %s: issue %s severity is blank; omit it rather than recording an empty value", ErrV2IssueMalformed, path, fields.ID)
	}

	resolution, err := decodeIssueResolution(path, fields.ID, fields.Resolution)
	if err != nil {
		return nil, err
	}

	if fields.DuplicateOf != "" && !issueIDPatternV2.MatchString(fields.DuplicateOf) {
		return nil, fmt.Errorf("%w: %s: issue %s duplicate_of %q must match I plus at least three digits", ErrV2InvalidID, path, fields.ID, fields.DuplicateOf)
	}

	history, err := decodeIssueHistory(path, fields.ID, fields.History)
	if err != nil {
		return nil, err
	}

	return &IssueV2{
		ID:           fields.ID,
		Title:        fields.Title,
		Type:         issueType,
		Status:       status,
		Origin:       origin,
		Tasks:        tasks,
		Checks:       checks,
		GuardrailIDs: guardrailIDs,
		Severity:     fields.Severity,
		Resolution:   resolution,
		DuplicateOf:  fields.DuplicateOf,
		History:      history,
		Source:       doc,
	}, nil
}

func decodeIssueType(path, id, raw string) (IssueType, error) {
	if raw == "" {
		return "", fmt.Errorf("%w: %s: issue %s missing required field type", ErrV2MissingField, path, id)
	}
	issueType := IssueType(raw)
	if !slices.Contains(issueTypes, issueType) {
		return "", fmt.Errorf("%w: %s: issue %s type %q; use defect, drift, guardrail, verification, or other", ErrV2IssueMalformed, path, id, raw)
	}
	return issueType, nil
}

// decodeIssueStatus decodes the Issue lifecycle vocabulary and refuses the
// Task one. A stage is rejected outright so an Issue's in_progress can never
// be read as the Task lifecycle's staged state, and planned, done, and every
// other Task value fall through to a named diagnostic rather than being
// healed into a valid-looking Issue status.
func decodeIssueStatus(path, id, raw, stage string) (IssueStatus, error) {
	if raw == "" {
		return "", fmt.Errorf("%w: %s: issue %s missing required field status", ErrV2MissingField, path, id)
	}
	if stage != "" {
		return "", fmt.Errorf("%w: %s: issue %s declares stage %q; issue status carries no stage", ErrV2InvalidLifecycle, path, id, stage)
	}
	status := IssueStatus(raw)
	if !slices.Contains(issueStatuses, status) {
		return "", fmt.Errorf("%w: %s: issue %s status %q; use open, in_progress, or resolved", ErrV2InvalidLifecycle, path, id, raw)
	}
	return status, nil
}

func decodeIssueOrigin(path, id string, raw *issueOriginFrontmatter) (IssueOrigin, error) {
	if raw == nil {
		return IssueOrigin{}, fmt.Errorf("%w: %s: issue %s missing required field source", ErrV2MissingField, path, id)
	}

	if raw.Kind == "" {
		return IssueOrigin{}, fmt.Errorf("%w: %s: issue %s missing required field source.kind", ErrV2MissingField, path, id)
	}
	kind := IssueOriginKind(raw.Kind)
	switch kind {
	case IssueOriginCheck, IssueOriginReport, IssueOriginMigration:
	default:
		return IssueOrigin{}, fmt.Errorf("%w: %s: issue %s source.kind %q; use check, report, or migration", ErrV2IssueMalformed, path, id, raw.Kind)
	}

	if kind == IssueOriginCheck {
		if raw.Check == "" {
			return IssueOrigin{}, fmt.Errorf("%w: %s: issue %s missing required field source.check", ErrV2MissingField, path, id)
		}
		if !checkIDPatternV2.MatchString(raw.Check) {
			return IssueOrigin{}, fmt.Errorf("%w: %s: issue %s source.check %q must match C plus at least three digits", ErrV2InvalidID, path, id, raw.Check)
		}
	} else if raw.Check != "" {
		return IssueOrigin{}, fmt.Errorf("%w: %s: issue %s source.check is recorded only when source.kind is check", ErrV2IssueMalformed, path, id)
	}

	actor, err := decodeV2Actor(ErrV2IssueMalformed, path, "issue", id, "source.actor", raw.Actor)
	if err != nil {
		return IssueOrigin{}, err
	}

	at, err := decodeV2Timestamp(ErrV2IssueMalformed, path, "issue", id, "source.at", raw.At)
	if err != nil {
		return IssueOrigin{}, err
	}

	return IssueOrigin{Kind: kind, Check: raw.Check, Actor: actor, At: at}, nil
}

// decodeIssueReferences validates one list of identity references for shape
// only. Whether the named records exist is resolved against the full index,
// so a reference that is well formed but dangling is not a decoding error.
func decodeIssueReferences(path, id, field string, pattern *regexp.Regexp, shape string, refs []string) ([]string, error) {
	out := make([]string, 0, len(refs))
	for _, ref := range refs {
		if !pattern.MatchString(ref) {
			return nil, fmt.Errorf("%w: %s: issue %s %s %q must match %s", ErrV2InvalidID, path, id, field, ref, shape)
		}
		out = append(out, ref)
	}
	return out, nil
}

// decodeIssueGuardrailIDs validates that each policy ID is non-empty and
// records it verbatim. The IDs are never resolved against a Guardrails file:
// Savepoint does not parse policy, so it cannot confirm that a rule exists.
func decodeIssueGuardrailIDs(path, id string, raw []string) ([]string, error) {
	out := make([]string, 0, len(raw))
	for _, guardrailID := range raw {
		if strings.TrimSpace(guardrailID) == "" {
			return nil, fmt.Errorf("%w: %s: issue %s guardrail_ids contains an empty entry", ErrV2IssueMalformed, path, id)
		}
		out = append(out, guardrailID)
	}
	return out, nil
}

// decodeIssueResolution decodes the resolution block's shape: a known
// disposition, a well-formed optional proof reference, and the actor and time
// behind the decision. Which disposition may name a proof Check, and whether
// a resolution belongs on this Issue's status at all, are cross-record
// obligations resolved against the index instead.
func decodeIssueResolution(path, id string, raw *issueResolutionFrontmatter) (*IssueResolution, error) {
	if raw == nil {
		return nil, nil
	}

	if raw.Disposition == "" {
		return nil, fmt.Errorf("%w: %s: issue %s missing required field resolution.disposition", ErrV2MissingField, path, id)
	}
	disposition := IssueDisposition(raw.Disposition)
	switch disposition {
	case IssueDispositionVerified, IssueDispositionAccepted, IssueDispositionDuplicate:
	default:
		return nil, fmt.Errorf("%w: %s: issue %s resolution.disposition %q; use verified, accepted, or duplicate", ErrV2IssueMalformed, path, id, raw.Disposition)
	}

	if raw.Check != "" && !checkIDPatternV2.MatchString(raw.Check) {
		return nil, fmt.Errorf("%w: %s: issue %s resolution.check %q must match C plus at least three digits", ErrV2InvalidID, path, id, raw.Check)
	}

	actor, err := decodeV2Actor(ErrV2IssueMalformed, path, "issue", id, "resolution.actor", raw.Actor)
	if err != nil {
		return nil, err
	}

	at, err := decodeV2Timestamp(ErrV2IssueMalformed, path, "issue", id, "resolution.at", raw.At)
	if err != nil {
		return nil, err
	}

	return &IssueResolution{
		Disposition: disposition,
		Check:       raw.Check,
		Actor:       actor,
		At:          at,
		Reason:      raw.Reason,
	}, nil
}

// decodeIssueHistory decodes history entries in recorded order, keeping the
// sequence exactly as written: the order is the record, so nothing here sorts
// or reorders it. Every entry must name when it happened, who recorded it,
// and a recognized kind; note and check stay optional, because a recheck that
// names its Check needs no prose to be meaningful.
func decodeIssueHistory(path, id string, raw []issueHistoryFrontmatter) ([]IssueHistoryEntry, error) {
	entries := make([]IssueHistoryEntry, 0, len(raw))
	for i, entry := range raw {
		field := fmt.Sprintf("history[%d]", i)

		if entry.Kind == "" {
			return nil, fmt.Errorf("%w: %s: issue %s missing required field %s.kind", ErrV2MissingField, path, id, field)
		}
		kind := IssueHistoryKind(entry.Kind)
		switch kind {
		case IssueHistoryObserved, IssueHistoryRepairAttempted, IssueHistoryRechecked,
			IssueHistoryDeferred, IssueHistoryReopened, IssueHistoryOwnerDecision:
		default:
			return nil, fmt.Errorf("%w: %s: issue %s %s.kind %q; use observed, repair_attempted, rechecked, deferred, reopened, or owner_decision", ErrV2IssueMalformed, path, id, field, entry.Kind)
		}

		at, err := decodeV2Timestamp(ErrV2IssueMalformed, path, "issue", id, field+".at", entry.At)
		if err != nil {
			return nil, err
		}

		actor, err := decodeV2Actor(ErrV2IssueMalformed, path, "issue", id, field+".actor", entry.Actor)
		if err != nil {
			return nil, err
		}

		if entry.Check != "" && !checkIDPatternV2.MatchString(entry.Check) {
			return nil, fmt.Errorf("%w: %s: issue %s %s.check %q must match C plus at least three digits", ErrV2InvalidID, path, id, field, entry.Check)
		}

		entries = append(entries, IssueHistoryEntry{
			At:    at,
			Actor: actor,
			Kind:  kind,
			Note:  entry.Note,
			Check: entry.Check,
		})
	}
	return entries, nil
}

// validateIssueResolutionObligations enforces the structural rules that make
// Issue closure honest: a resolution is present exactly when status is
// resolved — never on an open or in_progress Issue, so reopening must clear
// it — and each disposition's own proof obligations hold. It runs at index
// time, after Check-Issue pairing has resolved every Issue's checks list and
// every Check's recorded result, and walks Issue IDs in sorted order so a
// project with more than one violation reports the same one first on every
// run.
func validateIssueResolutionObligations(index *V2Index) error {
	for _, id := range slices.Sorted(maps.Keys(index.Issues)) {
		issue := index.Issues[id]

		if issue.Resolution == nil {
			if issue.Status == IssueStatusResolved {
				return fmt.Errorf("%w: %s: issue %s status resolved has no resolution", ErrV2IssueResolutionRequired, issue.Source.Path, id)
			}
			continue
		}

		if issue.Status != IssueStatusResolved {
			return fmt.Errorf("%w: %s: issue %s status %s carries a resolution; reopening must clear it", ErrV2IssueResolutionNotAllowed, issue.Source.Path, id, issue.Status)
		}

		if err := validateIssueDispositionObligations(index, issue); err != nil {
			return err
		}
	}
	return nil
}

// validateIssueDispositionObligations enforces the proof obligation specific
// to issue's disposition. Decoding already guarantees the disposition is one
// of the three known values.
func validateIssueDispositionObligations(index *V2Index, issue *IssueV2) error {
	switch issue.Resolution.Disposition {
	case IssueDispositionVerified:
		return validateVerifiedResolution(index, issue)
	case IssueDispositionAccepted:
		return validateAcceptedResolution(issue)
	default: // IssueDispositionDuplicate
		return validateDuplicateResolution(issue)
	}
}

// validateVerifiedResolution requires a proof Check that exists, recorded
// CLEAR, and appears in the Issue's own checks list — a duplicate or an
// accepted risk proves nothing, so only a verified repair may point at
// evidence at all.
func validateVerifiedResolution(index *V2Index, issue *IssueV2) error {
	check := issue.Resolution.Check
	if check == "" {
		return fmt.Errorf("%w: %s: issue %s verified resolution names no proof check", ErrV2IssueResolutionMissingProof, issue.Source.Path, issue.ID)
	}
	proof, ok := index.Checks[check]
	if !ok {
		return fmt.Errorf("%w: %s: issue %s verified resolution names missing proof check %s", ErrV2IssueResolutionMissingProof, issue.Source.Path, issue.ID, check)
	}
	if proof.Result != CheckResultClear {
		return fmt.Errorf("%w: %s: issue %s verified resolution's proof check %s recorded %s, not CLEAR", ErrV2IssueResolutionUnusableProof, issue.Source.Path, issue.ID, check, proof.Result)
	}
	if !slices.Contains(issue.Checks, check) {
		return fmt.Errorf("%w: %s: issue %s verified resolution's proof check %s does not appear in the issue's checks", ErrV2IssueResolutionUnusableProof, issue.Source.Path, issue.ID, check)
	}
	return nil
}

// validateAcceptedResolution requires an owner actor and a non-empty reason,
// and refuses a named proof Check: accepting risk is an owner decision, not a
// repair, so it cannot borrow a Check's authority.
func validateAcceptedResolution(issue *IssueV2) error {
	resolution := issue.Resolution
	if resolution.Check != "" {
		return fmt.Errorf("%w: %s: issue %s accepted resolution names proof check %s; accepting risk is not repair", ErrV2IssueResolutionFieldMismatch, issue.Source.Path, issue.ID, resolution.Check)
	}
	if resolution.Actor.Role != ActorRoleOwner || strings.TrimSpace(resolution.Actor.Session) == "" {
		return fmt.Errorf("%w: %s: issue %s accepted resolution requires an owner actor", ErrV2IssueResolutionFieldMismatch, issue.Source.Path, issue.ID)
	}
	if strings.TrimSpace(resolution.Reason) == "" {
		return fmt.Errorf("%w: %s: issue %s accepted resolution requires a non-empty reason", ErrV2IssueResolutionFieldMismatch, issue.Source.Path, issue.ID)
	}
	return nil
}

// validateDuplicateResolution requires duplicate_of to name the canonical
// Issue and refuses a named proof Check: a duplicate points at another
// record and proves nothing itself. duplicate_of's existence, self-reference,
// and cycle rules are validated once for every Issue by
// validateIssueDuplicateGraph, not repeated here.
func validateDuplicateResolution(issue *IssueV2) error {
	resolution := issue.Resolution
	if resolution.Check != "" {
		return fmt.Errorf("%w: %s: issue %s duplicate resolution names proof check %s; a duplicate proves nothing", ErrV2IssueResolutionFieldMismatch, issue.Source.Path, issue.ID, resolution.Check)
	}
	if issue.DuplicateOf == "" {
		return fmt.Errorf("%w: %s: issue %s duplicate resolution requires duplicate_of naming the canonical issue", ErrV2IssueResolutionFieldMismatch, issue.Source.Path, issue.ID)
	}
	return nil
}

// IssueIDsWithStatus returns the sorted IDs of every Issue currently in
// status.
//
// This and the counts below are derived from the loaded index on every call.
// Savepoint stores no Issue summary, register, or cached total anywhere in a
// project's files, so a listing cannot fall out of step with the records it
// describes, and reading one never writes anything.
func (index *V2Index) IssueIDsWithStatus(status IssueStatus) []string {
	ids := make([]string, 0, len(index.Issues))
	for _, id := range slices.Sorted(maps.Keys(index.Issues)) {
		if index.Issues[id].Status == status {
			ids = append(ids, id)
		}
	}
	return ids
}

// IssueStatusCounts returns how many Issues sit in each status. Every status
// is present even at zero, so a caller never has to read an absent key as
// none.
func (index *V2Index) IssueStatusCounts() map[IssueStatus]int {
	counts := make(map[IssueStatus]int, len(issueStatuses))
	for _, status := range issueStatuses {
		counts[status] = 0
	}
	for _, issue := range index.Issues {
		counts[issue.Status]++
	}
	return counts
}

// IssueTypeCounts returns how many Issues carry each type, every type present
// even at zero. The counts describe the backlog's shape; they decide nothing,
// because no gate reads an Issue's type.
func (index *V2Index) IssueTypeCounts() map[IssueType]int {
	counts := make(map[IssueType]int, len(issueTypes))
	for _, issueType := range issueTypes {
		counts[issueType] = 0
	}
	for _, issue := range index.Issues {
		counts[issue.Type]++
	}
	return counts
}
