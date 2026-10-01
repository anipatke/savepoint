package codehealth

import (
	"math"
	"slices"
)

// Snapshot is one immutable observation of project health (schema version 1).
// Its ID is derived from its canonical content, so any edit is detectable.
type Snapshot struct {
	Version    int                `json:"version"`
	ID         string             `json:"id"`
	Origin     Origin             `json:"origin"`
	Retention  Retention          `json:"retention"`
	CreatedAt  string             `json:"created_at"`
	Repository RepositoryIdentity `json:"repository"`
	Results    []CapabilityResult `json:"results"`
	Summary    Summary            `json:"summary"`
}

// RepositoryIdentity describes the observed repository state without storing
// any file content. Commit is empty for a repository without a useful commit.
type RepositoryIdentity struct {
	Commit           string `json:"commit,omitempty"`
	Dirty            bool   `json:"dirty"`
	InputFingerprint string `json:"input_fingerprint"`
}

// Provenance records which tool and definition produced a result.
type Provenance struct {
	Provider              ProviderKey `json:"provider,omitempty"`
	ProviderVersion       string      `json:"provider_version,omitempty"`
	ReportSchema          string      `json:"report_schema,omitempty"`
	MeasurementDefinition string      `json:"measurement_definition,omitempty"`
}

// Value is a measured number with its unit.
type Value struct {
	Number float64 `json:"number"`
	Unit   Unit    `json:"unit"`
}

// UnmarshalJSON requires a present number; see requiredNumber.
func (v *Value) UnmarshalJSON(data []byte) error {
	var raw struct {
		Number *float64 `json:"number"`
		Unit   Unit     `json:"unit"`
	}
	if err := decodeRecord(data, &raw); err != nil {
		return err
	}
	n, err := requiredNumber("value.number", raw.Number)
	if err != nil {
		return err
	}
	*v = Value{Number: n, Unit: raw.Unit}
	return nil
}

// Detail is one bounded, normalized supporting number such as total_tests.
type Detail struct {
	Key    string  `json:"key"`
	Number float64 `json:"number"`
}

// UnmarshalJSON requires a present number; see requiredNumber.
func (d *Detail) UnmarshalJSON(data []byte) error {
	var raw struct {
		Key    string   `json:"key"`
		Number *float64 `json:"number"`
	}
	if err := decodeRecord(data, &raw); err != nil {
		return err
	}
	n, err := requiredNumber("details.number", raw.Number)
	if err != nil {
		return err
	}
	*d = Detail{Key: raw.Key, Number: n}
	return nil
}

// EvidenceRef points at a repository-relative file that supports a result.
type EvidenceRef struct {
	Path string `json:"path"`
	Line int    `json:"line,omitempty"`
	Note string `json:"note,omitempty"`
}

// CapabilityResult is the scoped result for one configured instance.
type CapabilityResult struct {
	Capability   Capability    `json:"capability"`
	Name         string        `json:"name,omitempty"`
	Outcome      Outcome       `json:"outcome"`
	Freshness    Freshness     `json:"freshness"`
	CollectedAt  string        `json:"collected_at"`
	Provenance   Provenance    `json:"provenance"`
	ConfigDigest string        `json:"config_digest,omitempty"`
	Scope        []string      `json:"scope,omitempty"`
	Exclusions   []string      `json:"exclusions,omitempty"`
	Value        *Value        `json:"value,omitempty"`
	Reason       string        `json:"reason,omitempty"`
	Details      []Detail      `json:"details,omitempty"`
	Evidence     []EvidenceRef `json:"evidence,omitempty"`
}

// Summary holds deterministic labels. There is no composite score.
type Summary struct {
	Overall      Classification      `json:"overall"`
	Capabilities []CapabilitySummary `json:"capabilities,omitempty"`
}

// CapabilitySummary labels one result.
type CapabilitySummary struct {
	Capability     Capability     `json:"capability"`
	Provider       ProviderKey    `json:"provider,omitempty"`
	Name           string         `json:"name,omitempty"`
	Classification Classification `json:"classification"`
	Explanation    string         `json:"explanation,omitempty"`
}

// Validate checks the snapshot, including that ID matches its content.
func (s Snapshot) Validate() error {
	if s.Version != SnapshotVersion {
		return fieldError(ErrUnsupportedVersion, "version", "got %d, want %d", s.Version, SnapshotVersion)
	}
	if err := s.validateOrigin(); err != nil {
		return err
	}
	if err := validateTimestamp("created_at", s.CreatedAt); err != nil {
		return err
	}
	if err := s.Repository.validate("repository"); err != nil {
		return err
	}
	if err := s.validateResults(); err != nil {
		return err
	}
	if err := s.Summary.validate("summary", s.Results); err != nil {
		return err
	}
	if err := validateDigest("id", s.ID); err != nil {
		return err
	}
	if want := s.ComputeID(); s.ID != want {
		return fieldError(ErrIdentityMismatch, "id", "got %s, content hashes to %s", s.ID, want)
	}
	return nil
}

// validateOrigin enforces the fixed origin-to-retention pairing: official
// snapshots are permanent, manual snapshots may be pruned.
func (s Snapshot) validateOrigin() error {
	want, ok := map[Origin]Retention{OriginOfficial: RetentionPermanent, OriginManual: RetentionPrunable}[s.Origin]
	if !ok {
		return fieldError(ErrInvalidOrigin, "origin", "%q", s.Origin)
	}
	if s.Retention != want {
		return fieldError(ErrInvalidOrigin, "retention", "%s snapshots must be %s, got %q", s.Origin, want, s.Retention)
	}
	return nil
}

func (r RepositoryIdentity) validate(field string) error {
	if r.Commit != "" {
		if err := validateCommit(field+".commit", r.Commit); err != nil {
			return err
		}
	}
	if err := validateDigest(field+".input_fingerprint", r.InputFingerprint); err != nil {
		return err
	}
	if r.Commit == "" && !r.Dirty {
		return fieldError(ErrInvalidRepository, field, "a repository without a commit has no clean state; mark it dirty")
	}
	return nil
}

func (s Snapshot) validateResults() error {
	if len(s.Results) > MaxResults {
		return fieldError(ErrUnboundedDetail, "results", "has %d entries; at most %d", len(s.Results), MaxResults)
	}
	seen := make(map[instanceKey]bool, len(s.Results))
	for i, r := range s.Results {
		field := indexed("results", i)
		if err := r.validate(field); err != nil {
			return err
		}
		key := r.key()
		if seen[key] {
			return fieldError(ErrDuplicateInstance, field, "%s via %q named %q appears more than once", r.Capability, r.Provenance.Provider, r.Name)
		}
		seen[key] = true
	}
	return nil
}

// key is the instance this result belongs to.
func (r CapabilityResult) key() instanceKey {
	return instanceKey{r.Capability, r.Provenance.Provider, r.Name}
}

func (r CapabilityResult) validate(field string) error {
	if _, ok := capabilityUnits[r.Capability]; !ok {
		return fieldError(ErrUnknownCapability, field+".capability", "%q", r.Capability)
	}
	if !slices.Contains(Outcomes(), r.Outcome) {
		return fieldError(ErrInvalidOutcome, field+".outcome", "%q", r.Outcome)
	}
	if err := validateToken(field+".name", r.Name, false); err != nil {
		return err
	}
	if err := r.validateFreshness(field); err != nil {
		return err
	}
	if err := validateTimestamp(field+".collected_at", r.CollectedAt); err != nil {
		return err
	}
	if err := r.validateProvenance(field); err != nil {
		return err
	}
	if err := r.validateValue(field); err != nil {
		return err
	}
	if err := r.validateBounded(field); err != nil {
		return err
	}
	if err := validatePatternList(field+".scope", r.Scope); err != nil {
		return err
	}
	return validatePatternList(field+".exclusions", r.Exclusions)
}

// validateFreshness keeps outcome and freshness orthogonal while preventing
// evidence that was never measured from claiming to be fresh or stale.
func (r CapabilityResult) validateFreshness(field string) error {
	switch r.Freshness {
	case FreshnessFresh, FreshnessStale:
		if !r.Outcome.Measured() {
			return fieldError(ErrInvalidFreshness, field+".freshness", "%s evidence was never measured, so freshness must be unknown", r.Outcome)
		}
	case FreshnessUnknown:
	default:
		return fieldError(ErrInvalidFreshness, field+".freshness", "%q", r.Freshness)
	}
	return nil
}

func (r CapabilityResult) validateProvenance(field string) error {
	p := r.Provenance
	if r.Outcome == OutcomeNotConfigured {
		if p != (Provenance{}) || r.ConfigDigest != "" || r.Name != "" {
			return fieldError(ErrInvalidOutcome, field, "a not_configured result has no provider, name, or configuration")
		}
		return nil
	}
	if err := validateCapabilityProvider(field+".provenance", r.Capability, p.Provider); err != nil {
		return err
	}
	if err := validateDigest(field+".config_digest", r.ConfigDigest); err != nil {
		return err
	}
	measured := r.Outcome.Measured()
	for _, f := range []struct {
		name, v string
	}{
		{"provider_version", p.ProviderVersion},
		{"report_schema", p.ReportSchema},
		{"measurement_definition", p.MeasurementDefinition},
	} {
		if err := validateToken(field+".provenance."+f.name, f.v, measured && f.name != "report_schema"); err != nil {
			return err
		}
	}
	return nil
}

// validateValue allows a value only on available or partial results and
// checks it against the capability's unit and range.
func (r CapabilityResult) validateValue(field string) error {
	if !r.Outcome.Measured() {
		if r.Value != nil {
			return fieldError(ErrValueNotAllowed, field+".value", "%s results cannot carry a measured value", r.Outcome)
		}
		return nil
	}
	if r.Outcome == OutcomeAvailable && r.Value == nil {
		return fieldError(ErrMissingValue, field+".value", "available results need a value")
	}
	if r.Outcome == OutcomePartial && r.Reason == "" {
		return fieldError(ErrMissingValue, field+".reason", "partial results must say what is missing")
	}
	if r.Value == nil {
		return nil
	}
	return r.Value.validate(field+".value", r.Capability)
}

func (v Value) validate(field string, c Capability) error {
	if want := capabilityUnits[c]; v.Unit != want {
		return fieldError(ErrIncompatibleUnit, field+".unit", "%s is measured in %s, got %q", c, want, v.Unit)
	}
	if math.IsNaN(v.Number) || math.IsInf(v.Number, 0) || v.Number < 0 {
		return fieldError(ErrIncompatibleUnit, field+".number", "%v must be a finite non-negative number", v.Number)
	}
	if v.Unit == UnitPercent && v.Number > 100 {
		return fieldError(ErrIncompatibleUnit, field+".number", "%v exceeds 100 percent", v.Number)
	}
	if v.Unit == UnitCount && v.Number != math.Trunc(v.Number) {
		return fieldError(ErrIncompatibleUnit, field+".number", "%v is not a whole count", v.Number)
	}
	return nil
}

// validateSeverity checks the vulnerability severity counts: each must be a
// whole, non-negative count, and together they cannot exceed the total.
func (r CapabilityResult) validateSeverity(field string) error {
	if r.Capability != CapabilityDependencyVulnerability {
		return nil
	}
	var sum float64
	for i, d := range r.Details {
		if d.Key != DetailHighVulnerabilities && d.Key != DetailCriticalVulnerabilities {
			continue
		}
		if d.Number < 0 || d.Number != math.Trunc(d.Number) {
			return fieldError(ErrIncompatibleUnit, indexed(field+".details", i)+".number", "%v must be a whole, non-negative count", d.Number)
		}
		sum += d.Number
	}
	if r.Value != nil && sum > r.Value.Number {
		return fieldError(ErrIncompatibleUnit, field+".details", "severity counts total %v, more than the %v vulnerabilities reported", sum, r.Value.Number)
	}
	return nil
}

func (r CapabilityResult) validateBounded(field string) error {
	if err := validateText(field+".reason", r.Reason, MaxReasonLen); err != nil {
		return err
	}
	if len(r.Details) > MaxDetails {
		return fieldError(ErrUnboundedDetail, field+".details", "has %d entries; at most %d", len(r.Details), MaxDetails)
	}
	keys := make(map[string]bool, len(r.Details))
	for i, d := range r.Details {
		f := indexed(field+".details", i)
		if err := validateToken(f+".key", d.Key, true); err != nil {
			return err
		}
		if keys[d.Key] {
			return fieldError(ErrUnboundedDetail, f+".key", "duplicate key %q", d.Key)
		}
		keys[d.Key] = true
		if math.IsNaN(d.Number) || math.IsInf(d.Number, 0) {
			return fieldError(ErrIncompatibleUnit, f+".number", "%v must be finite", d.Number)
		}
	}
	if err := r.validateSeverity(field); err != nil {
		return err
	}
	if len(r.Evidence) > MaxEvidence {
		return fieldError(ErrUnboundedDetail, field+".evidence", "has %d entries; at most %d", len(r.Evidence), MaxEvidence)
	}
	for i, e := range r.Evidence {
		if err := e.validate(indexed(field+".evidence", i)); err != nil {
			return err
		}
	}
	return nil
}

func (e EvidenceRef) validate(field string) error {
	if err := validateRelativePath(ErrUnsafeEvidence, field+".path", e.Path); err != nil {
		return err
	}
	for _, seg := range splitPath(e.Path) {
		if sensitiveSegment(seg) {
			return fieldError(ErrUnsafeEvidence, field+".path", "%q points at a sensitive location", e.Path)
		}
	}
	if e.Line < 0 {
		return fieldError(ErrUnsafeEvidence, field+".line", "%d is negative", e.Line)
	}
	return validateText(field+".note", e.Note, MaxNoteLen)
}

func (s Summary) validate(field string, results []CapabilityResult) error {
	if err := validateClassification(field+".overall", s.Overall); err != nil {
		return err
	}
	if len(s.Capabilities) != len(results) {
		return fieldError(ErrInvalidSummary, field+".capabilities", "has %d entries for %d results", len(s.Capabilities), len(results))
	}
	byKey := make(map[instanceKey]CapabilityResult, len(results))
	for _, r := range results {
		byKey[r.key()] = r
	}
	seen := make(map[instanceKey]bool, len(results))
	for i, cs := range s.Capabilities {
		f := indexed(field+".capabilities", i)
		key := instanceKey{cs.Capability, cs.Provider, cs.Name}
		r, ok := byKey[key]
		if !ok || seen[key] {
			return fieldError(ErrInvalidSummary, f, "%s via %q named %q does not match exactly one result", cs.Capability, cs.Provider, cs.Name)
		}
		seen[key] = true
		if err := validateClassification(f+".classification", cs.Classification); err != nil {
			return err
		}
		if cs.Classification == ClassificationGood && (r.Outcome != OutcomeAvailable || r.Freshness != FreshnessFresh) {
			return fieldError(ErrInvalidSummary, f, "good requires available, fresh evidence; got %s and %s", r.Outcome, r.Freshness)
		}
		if err := validateText(f+".explanation", cs.Explanation, MaxReasonLen); err != nil {
			return err
		}
	}
	if s.Overall == ClassificationGood {
		for _, cs := range s.Capabilities {
			if cs.Classification != ClassificationGood {
				return fieldError(ErrInvalidSummary, field+".overall", "good requires every capability to be good")
			}
		}
	}
	return nil
}

func validateClassification(field string, c Classification) error {
	switch c {
	case ClassificationGood, ClassificationWatch, ClassificationNeedsAttention, ClassificationUnknown:
		return nil
	}
	return fieldError(ErrInvalidSummary, field, "%q", c)
}

// DecodeSnapshot parses and validates a JSON snapshot.
func DecodeSnapshot(data []byte) (Snapshot, error) {
	var s Snapshot
	if err := decodeStrict(data, &s); err != nil {
		return Snapshot{}, err
	}
	return s, s.Validate()
}
