package codehealth

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"slices"
	"strings"
)

// splitPath splits a validated forward-slash path into segments.
func splitPath(p string) []string { return strings.Split(p, "/") }

// hashCanonical returns "sha256:" plus the digest of v's canonical JSON. All
// hashed types use structs and sorted slices only, so encoding/json output is
// byte-stable across runs and platforms.
func hashCanonical(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		// Only finite numbers and plain strings reach here after validation.
		panic("codehealth: canonical encoding failed: " + err.Error())
	}
	sum := sha256.Sum256(b)
	return digestPrefix + hex.EncodeToString(sum[:])
}

// Canonical returns a deep copy with order-insensitive lists sorted, so equal
// observations serialize identically.
func (s Snapshot) Canonical() Snapshot {
	c := s
	c.Results = make([]CapabilityResult, len(s.Results))
	for i, r := range s.Results {
		c.Results[i] = r.canonical()
	}
	slices.SortFunc(c.Results, func(a, b CapabilityResult) int {
		return strings.Compare(sortKey(a.key()), sortKey(b.key()))
	})
	c.Summary.Capabilities = slices.Clone(s.Summary.Capabilities)
	slices.SortFunc(c.Summary.Capabilities, func(a, b CapabilitySummary) int {
		return strings.Compare(sortKey(instanceKey{a.Capability, a.Provider, a.Name}), sortKey(instanceKey{b.Capability, b.Provider, b.Name}))
	})
	return c
}

// sortKey orders instances by capability, provider, then name. A NUL separator
// cannot appear in a validated token, so keys never run together.
func sortKey(k instanceKey) string {
	return string(k.capability) + "/" + string(k.provider) + "\x00" + k.name
}

func (r CapabilityResult) canonical() CapabilityResult {
	r.Scope = sortedCopy(r.Scope)
	r.Exclusions = sortedCopy(r.Exclusions)
	r.Details = slices.Clone(r.Details)
	slices.SortFunc(r.Details, func(a, b Detail) int { return strings.Compare(a.Key, b.Key) })
	r.Evidence = slices.Clone(r.Evidence)
	slices.SortFunc(r.Evidence, func(a, b EvidenceRef) int {
		if c := strings.Compare(a.Path, b.Path); c != 0 {
			return c
		}
		if a.Line != b.Line {
			return a.Line - b.Line
		}
		return strings.Compare(a.Note, b.Note)
	})
	if r.Value != nil {
		v := *r.Value
		r.Value = &v
	}
	return r
}

func sortedCopy(in []string) []string {
	if len(in) == 0 {
		return nil
	}
	out := slices.Clone(in)
	slices.Sort(out)
	return out
}

// ComputeID derives the snapshot identity from canonical content with the ID
// field blanked. It identifies one observation, so it includes the commit and
// the time.
func (s Snapshot) ComputeID() string {
	c := s.Canonical()
	c.ID = ""
	return hashCanonical(c)
}

// seriesKey is everything that decides whether two results are comparable. The
// observed commit, time, and value are deliberately absent.
type seriesKey struct {
	Capability            Capability  `json:"capability"`
	Provider              ProviderKey `json:"provider"`
	ProviderVersion       string      `json:"provider_version"`
	ReportSchema          string      `json:"report_schema"`
	MeasurementDefinition string      `json:"measurement_definition"`
	ConfigDigest          string      `json:"config_digest"`
	Scope                 []string    `json:"scope"`
	Exclusions            []string    `json:"exclusions"`
}

// SeriesID identifies the comparison series a result belongs to. A change to
// provider, schema, measurement definition, effective configuration,
// exclusions, or scope begins a new series; a new commit does not.
func (r CapabilityResult) SeriesID() string {
	return hashCanonical(seriesKey{
		Capability:            r.Capability,
		Provider:              r.Provenance.Provider,
		ProviderVersion:       r.Provenance.ProviderVersion,
		ReportSchema:          r.Provenance.ReportSchema,
		MeasurementDefinition: r.Provenance.MeasurementDefinition,
		ConfigDigest:          r.ConfigDigest,
		Scope:                 nonNil(sortedCopy(r.Scope)),
		Exclusions:            nonNil(sortedCopy(r.Exclusions)),
	})
}

func nonNil(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}

// digestKey is the part of a configured instance that changes what is
// measured. Timeouts, thresholds, executable location, and report path do not.
type digestKey struct {
	Capability Capability  `json:"capability"`
	Provider   ProviderKey `json:"provider"`
	Args       []string    `json:"args"`
}

// Digest fingerprints the measurement-relevant configuration of an instance.
// Scope and exclusions are hashed separately, in the series identity.
func (cc CapabilityConfig) Digest() string {
	args := cc.Args
	if args == nil {
		args = []string{}
	}
	return hashCanonical(digestKey{cc.Capability, cc.Provider, args})
}
