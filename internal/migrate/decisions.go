// decisions.go gives every ambiguity Plan raises a stable identity, a
// blocking-versus-advisory classification, and a place to answer: the
// --decisions FILE input this file parses. Plan (plan.go) still detects each
// ambiguity as it walks the project; this file owns what an ambiguity *is*
// and how an owner's answer to it is read, validated, and recorded — never
// how it is guessed.
package migrate

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"time"

	"gopkg.in/yaml.v3"
)

// AmbiguityKind names one of the ways Plan refuses to guess. Blocking kinds
// change what the migrated project would mean and refuse apply until an
// owner decision resolves them; advisory kinds are reported for visibility
// and never block. See isBlockingAmbiguityKind for the classification.
type AmbiguityKind string

const (
	// AmbiguityUnrecognizedLifecycle is an active record (epic or task)
	// carrying a status value that is neither canonical nor a known healed
	// alias. Choices are the record's canonical status vocabulary.
	AmbiguityUnrecognizedLifecycle AmbiguityKind = "unrecognized_lifecycle"
	// AmbiguityMissingDependencyTarget is an active task's depends_on
	// reference that resolves to no known task, planned or archived.
	// Migration never invents a replacement target; the only choice offered
	// is to drop the reference and let the owner fix the V1 source and rerun
	// if a real target was intended.
	AmbiguityMissingDependencyTarget AmbiguityKind = "missing_dependency_target"
	// AmbiguityDuplicateSourceIdentity is two source files declaring the same
	// original ID within the same identity scope (a release and epic, for a
	// task), so allocation cannot tell which one the ID actually names.
	AmbiguityDuplicateSourceIdentity AmbiguityKind = "duplicate_source_identity"
	// AmbiguityUnresolvedNarrativeFind is a `duplicate` audit finding whose
	// duplicate_of names no finding this project actually has — a narrative
	// cross-reference to an identity that does not exist, rather than a
	// finding with no identity of its own.
	AmbiguityUnresolvedNarrativeFind AmbiguityKind = "unresolved_narrative_finding"

	// AmbiguityUnclassifiedFile is a source file matching no known Role. It
	// is archived intact regardless, so it never blocks; it is named so the
	// owner can see what migration could not place.
	AmbiguityUnclassifiedFile AmbiguityKind = "unclassified_file"
)

// isBlockingAmbiguityKind reports whether kind refuses apply until resolved.
// This is the single source of the blocking/advisory split: Plan and the
// manifest both consult it rather than each keeping their own list.
func isBlockingAmbiguityKind(kind AmbiguityKind) bool {
	switch kind {
	case AmbiguityUnrecognizedLifecycle, AmbiguityMissingDependencyTarget,
		AmbiguityDuplicateSourceIdentity, AmbiguityUnresolvedNarrativeFind:
		return true
	default:
		return false
	}
}

// Ambiguity is one named, stably identified condition Plan found that either
// blocks apply until an owner decision resolves it, or is purely advisory.
// ID is derived from Kind and Path (plus a discriminator for a kind that can
// recur more than once against the same path), so the same unchanged project
// previews the same ambiguity IDs every run.
type Ambiguity struct {
	ID       string
	Kind     AmbiguityKind
	Path     string
	Detail   string
	Blocking bool
	// Choices lists the concrete decision values a decisions file may supply
	// for this ambiguity. Empty for a purely advisory ambiguity, which offers
	// nothing to decide.
	Choices []string

	Resolved bool
	Decision string // the owner-supplied resolution, once Decisions names this ID

	// DecisionSourceFile and DecisionAt are the provenance of Decision, once
	// resolved: the --decisions FILE path it was read from and when. Both are
	// zero when Resolved is false.
	DecisionSourceFile string
	DecisionAt         time.Time
}

// DecisionValue is one owner-supplied resolution for a single ambiguity ID,
// carrying the provenance the manifest records against it.
type DecisionValue struct {
	Value      string
	SourceFile string
	DecidedAt  time.Time
}

// Decisions maps an Ambiguity ID to the owner-supplied resolution recorded
// against it. Plan itself only consumes this map; ReadDecisionsFile builds
// one from the --decisions FILE input. Nothing in this package ever produces
// a Decisions entry except an explicit, owner-authored one — no default, no
// heuristic, no similarity match.
type Decisions map[string]DecisionValue

// decisionsFileEntry and decisionsFileSchema are the --decisions FILE
// input's YAML shape: an explicit list naming each ambiguity ID and the
// single concrete choice that resolves it. A list rather than a map keyed by
// ambiguity ID keeps IDs that themselves contain ":" or "#" legible without
// YAML key-escaping.
type decisionsFileEntry struct {
	ID    string `yaml:"id"`
	Value string `yaml:"value"`
}

type decisionsFileSchema struct {
	Decisions []decisionsFileEntry `yaml:"decisions"`
}

// ReadDecisionsFile parses the --decisions FILE input at path into a
// Decisions map, stamping every entry with path and at as provenance. It
// only reads the file: nothing here resolves an ambiguity, because doing
// that needs the ambiguities Plan actually raised and their stated Choices,
// which this function never sees. Validating an entry's ID and value against
// a real plan happens inside Plan (see validateDecisions), once those
// ambiguities exist to validate against.
func ReadDecisionsFile(path string, at time.Time) (Decisions, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read decisions file %s: %w", path, err)
	}
	var schema decisionsFileSchema
	if err := yaml.Unmarshal(content, &schema); err != nil {
		return nil, fmt.Errorf("parse decisions file %s: %w", path, err)
	}
	decisions := make(Decisions, len(schema.Decisions))
	for _, e := range schema.Decisions {
		decisions[e.ID] = DecisionValue{Value: e.Value, SourceFile: path, DecidedAt: at}
	}
	return decisions, nil
}

// Distinct decisions-file diagnostics, so a caller (and a test) can tell a
// stale decisions file apart from one with a value the ambiguity never
// offered, rather than string-matching a message.
var (
	ErrUnknownAmbiguityID      = errors.New("migrate: decisions file names an ambiguity id the plan did not raise")
	ErrDecisionValueNotAllowed = errors.New("migrate: decision value is not one of the ambiguity's stated choices")
)

// validateDecisions checks a Decisions map against the ambiguities a build
// actually raised, once every ambiguity (and its Choices) is known. An ID
// naming no raised ambiguity is an error, not a silently ignored entry, so a
// stale decisions file from an earlier preview of a since-changed project is
// never mistaken for a current one. A value outside the ambiguity's stated
// Choices is rejected the same way, rather than accepted as whatever the
// owner happened to type.
func validateDecisions(decisions Decisions, ambiguities []Ambiguity) error {
	byID := make(map[string]Ambiguity, len(ambiguities))
	for _, a := range ambiguities {
		byID[a.ID] = a
	}

	ids := make([]string, 0, len(decisions))
	for id := range decisions {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	for _, id := range ids {
		amb, ok := byID[id]
		if !ok {
			return fmt.Errorf("%w: %s", ErrUnknownAmbiguityID, id)
		}
		value := decisions[id].Value
		if len(amb.Choices) > 0 && !containsString(amb.Choices, value) {
			return fmt.Errorf("%w: ambiguity %s got %q, want one of %v", ErrDecisionValueNotAllowed, id, value, amb.Choices)
		}
	}
	return nil
}

// addAmbiguity records one ambiguity with an ID derived solely from kind and
// path, for a kind that can never recur more than once against the same
// path.
func (b *planBuilder) addAmbiguity(kind AmbiguityKind, path, detail string, choices []string) {
	b.addAmbiguityDiscriminated(kind, path, "", detail, choices)
}

// addAmbiguityDiscriminated is addAmbiguity plus a discriminator folded into
// the ID for a kind that can recur more than once against the same path (for
// example, a task naming more than one missing dependency): the
// discriminator is the specific reference or identity text that makes this
// occurrence distinct, so the ID stays stable across repeated previews
// without colliding with a sibling ambiguity on the same path.
func (b *planBuilder) addAmbiguityDiscriminated(kind AmbiguityKind, path, discriminator, detail string, choices []string) {
	id := string(kind) + ":" + path
	if discriminator != "" {
		id += "#" + discriminator
	}
	dv, resolved := b.decisions[id]
	amb := Ambiguity{
		ID:       id,
		Kind:     kind,
		Path:     path,
		Detail:   detail,
		Blocking: isBlockingAmbiguityKind(kind),
		Choices:  choices,
		Resolved: resolved,
	}
	if resolved {
		amb.Decision = dv.Value
		amb.DecisionSourceFile = dv.SourceFile
		amb.DecisionAt = dv.DecidedAt
	}
	b.ambiguities = append(b.ambiguities, amb)
}
