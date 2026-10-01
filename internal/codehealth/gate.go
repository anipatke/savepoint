package codehealth

import (
	"fmt"
	"slices"
	"strings"
)

// optInBlockingCapabilities are the capabilities whose Needs Attention result
// blocks only when the instance sets Blocking. Tests and vulnerabilities
// already block by default.
var optInBlockingCapabilities = []Capability{
	CapabilityCoverage,
	CapabilityComplexity,
	CapabilityDuplication,
}

// Disposition says what a result means for clearance.
type Disposition string

const (
	DispositionBlocks        Disposition = "blocks"
	DispositionReported      Disposition = "reported"
	DispositionNotConfigured Disposition = "not_configured"
)

// Kind names which statement a verdict line makes. A failed tool, thin or old
// evidence, and unhealthy code are always different kinds.
type Kind string

const (
	KindCollectionFailure Kind = "collection_failure"
	KindIncomplete        Kind = "incomplete"
	KindStale             Kind = "stale"
	KindUnhealthy         Kind = "unhealthy_measurement"
	KindNotConfigured     Kind = "not_configured"
	KindNoFinding         Kind = "no_finding"
)

// ResultVerdict is the verdict on one configured instance.
type ResultVerdict struct {
	Capability  Capability
	Provider    ProviderKey
	Name        string
	Required    bool
	Disposition Disposition
	Kind        Kind
	Reason      string
}

// Verdict is the policy outcome for one official snapshot. It says only
// whether health results block clearance; it never says an Objective is clear.
type Verdict struct {
	SnapshotID string
	Results    []ResultVerdict
}

// Blocks reports whether any result blocks clearance.
func (v Verdict) Blocks() bool {
	return slices.ContainsFunc(v.Results, func(r ResultVerdict) bool { return r.Disposition == DispositionBlocks })
}

// failureStatements describes outcomes that produced no measurement.
var failureStatements = map[Outcome]string{
	OutcomeAbsent:      "no report was found",
	OutcomeUnsupported: "the provider is not supported",
	OutcomeUnavailable: "the tool was unavailable",
	OutcomeFailed:      "collection failed",
	OutcomeTimedOut:    "collection timed out",
	OutcomeCancelled:   "collection was cancelled",
}

// Plain statements the verdict is built from, kept together so wording changes
// in one place.
const (
	textNoResult      = "no result was recorded for this configured instance"
	textStale         = "the evidence is stale"
	textIncomplete    = "the evidence is incomplete"
	textNoFinding     = "no blocking finding"
	textNotConfigured = "not configured"
	textFailingTests  = "%s failing"
	textSevereVulns   = "%s known critical or high severity"
	textNeedsReview   = "%s of unknown severity need review"
	textNeedsAttn     = "needs attention and this instance is set to block"
)

// Evaluate applies the blocking policy to an official snapshot and the project
// configuration. It does no IO and trusts the snapshot to be validated. A manual
// snapshot has no standing as Check evidence and is refused.
func Evaluate(s Snapshot, cfg Config) (Verdict, error) {
	if s.Origin != OriginOfficial {
		return Verdict{}, fmt.Errorf("%w: origin is %q; only an official snapshot can be evaluated", ErrManualSnapshot, s.Origin)
	}
	if err := cfg.Validate(); err != nil {
		return Verdict{}, err
	}
	configs := make(map[instanceKey]CapabilityConfig, len(cfg.Capabilities))
	for _, cc := range cfg.Capabilities {
		configs[instanceKey{cc.Capability, cc.Provider, cc.Name}] = cc
	}
	classes := make(map[instanceKey]Classification, len(s.Summary.Capabilities))
	for _, cs := range s.Summary.Capabilities {
		classes[instanceKey{cs.Capability, cs.Provider, cs.Name}] = cs.Classification
	}

	v := Verdict{SnapshotID: s.ID}
	seen := make(map[instanceKey]bool, len(s.Results))
	for _, r := range s.Results {
		key := r.key()
		seen[key] = true
		if r.Outcome == OutcomeNotConfigured {
			v.Results = append(v.Results, ResultVerdict{
				Capability: r.Capability, Disposition: DispositionNotConfigured,
				Kind: KindNotConfigured, Reason: textNotConfigured,
			})
			continue
		}
		v.Results = append(v.Results, judge(r, configs[key], classes[key]))
	}
	for key, cc := range configs {
		if !seen[key] {
			v.Results = append(v.Results, missing(cc))
		}
	}
	slices.SortFunc(v.Results, func(a, b ResultVerdict) int {
		return strings.Compare(sortKey(instanceKey{a.Capability, a.Provider, a.Name}), sortKey(instanceKey{b.Capability, b.Provider, b.Name}))
	})
	return v, nil
}

// missing is the verdict for a configured instance the snapshot has no result
// for: a collection failure, which blocks only when the instance is required.
func missing(cc CapabilityConfig) ResultVerdict {
	return ResultVerdict{
		Capability: cc.Capability, Provider: cc.Provider, Name: cc.Name, Required: cc.Required,
		Disposition: failureDisposition(cc.Required), Kind: KindCollectionFailure, Reason: textNoResult,
	}
}

func failureDisposition(required bool) Disposition {
	if required {
		return DispositionBlocks
	}
	return DispositionReported
}

// statement is one thing the verdict says about a result.
type statement struct {
	kind   Kind
	reason string
}

// judge decides one measured or unmeasured result. A default or opt-in
// unhealthy-measurement rule wins as the primary statement, so failing tests are
// never softened by being optional, stale, or partial. Evidence-quality and
// review statements that also apply are kept after it, never dropped.
func judge(r CapabilityResult, cc CapabilityConfig, class Classification) ResultVerdict {
	rv := ResultVerdict{
		Capability: r.Capability, Provider: r.Provenance.Provider, Name: r.Name,
		Required: cc.Required, Disposition: DispositionReported, Kind: KindNoFinding, Reason: textNoFinding,
	}
	if !r.Outcome.Measured() {
		rv.Kind = KindCollectionFailure
		rv.Reason = failureStatements[r.Outcome]
		rv.Disposition = failureDisposition(cc.Required)
		return rv
	}
	var statements []statement
	if reason, blocks := defaultBlocker(r); blocks {
		rv.Disposition = DispositionBlocks
		statements = append(statements, statement{KindUnhealthy, reason})
	} else if cc.Blocking && r.Freshness == FreshnessFresh && class == ClassificationNeedsAttention {
		rv.Disposition = DispositionBlocks
		statements = append(statements, statement{KindUnhealthy, textNeedsAttn})
	}
	if r.Freshness == FreshnessStale {
		if rv.Disposition != DispositionBlocks {
			rv.Disposition = failureDisposition(cc.Required)
		}
		statements = append(statements, statement{KindStale, textStale})
	}
	if r.Outcome == OutcomePartial {
		reason := textIncomplete
		if r.Reason != "" {
			reason += ": " + r.Reason
		}
		statements = append(statements, statement{KindIncomplete, reason})
	}
	if n := detailNumber(r, DetailUnknownVulnerabilities); n > 0 && r.Capability == CapabilityDependencyVulnerability {
		statements = append(statements, statement{KindUnhealthy, fmt.Sprintf(textNeedsReview, count(n))})
	}
	if len(statements) == 0 {
		return rv
	}
	rv.Kind = statements[0].kind
	reasons := make([]string, len(statements))
	for i, st := range statements {
		reasons[i] = st.reason
	}
	rv.Reason = strings.Join(reasons, "; ")
	return rv
}

// defaultBlocker applies the rules that block whether the instance is required
// or optional: any failing test, and any critical or high vulnerability.
func defaultBlocker(r CapabilityResult) (string, bool) {
	switch r.Capability {
	case CapabilityTests:
		if r.Value != nil && r.Value.Number > 0 {
			return fmt.Sprintf(textFailingTests, count(r.Value.Number)), true
		}
	case CapabilityDependencyVulnerability:
		if n := detailNumber(r, DetailCriticalVulnerabilities) + detailNumber(r, DetailHighVulnerabilities); n > 0 {
			return fmt.Sprintf(textSevereVulns, count(n)), true
		}
	}
	return "", false
}

func detailNumber(r CapabilityResult, key string) float64 {
	for _, d := range r.Details {
		if d.Key == key {
			return d.Number
		}
	}
	return 0
}

func count(n float64) string { return fmt.Sprintf("%.0f", n) }

// Headlines are the only statements about the whole report. Neither says clear
// or healthy: clearance belongs to the independent checker.
const (
	headlineBlocks    = "Code Health blocks clearance."
	headlineNoBlock   = "Code Health does not block clearance."
	headlineNoResults = "Code Health has no results."
)

var dispositionLabels = map[Disposition]string{
	DispositionBlocks:        "blocks clearance",
	DispositionReported:      "reported only",
	DispositionNotConfigured: "not configured",
}

var kindLabels = map[Kind]string{
	KindCollectionFailure: "collection failure",
	KindIncomplete:        "incomplete",
	KindStale:             "stale",
	KindUnhealthy:         "unhealthy measurement",
	KindNotConfigured:     "not configured",
	KindNoFinding:         "no finding",
}

// Render writes the verdict as deterministic plain text, one line per result.
func (v Verdict) Render() string {
	var b strings.Builder
	switch {
	case v.Blocks():
		b.WriteString(headlineBlocks)
	case len(v.Results) == 0:
		b.WriteString(headlineNoResults)
	default:
		b.WriteString(headlineNoBlock)
	}
	b.WriteString("\n")
	for _, r := range v.Results {
		b.WriteString(r.line())
		b.WriteString("\n")
	}
	return b.String()
}

func (r ResultVerdict) line() string {
	label := string(r.Capability)
	if r.Provider != "" {
		label += " via " + string(r.Provider)
	}
	if r.Name != "" {
		label += " (" + r.Name + ")"
	}
	need := "optional"
	if r.Required {
		need = "required"
	}
	if r.Disposition == DispositionNotConfigured {
		return fmt.Sprintf("- %s: %s", label, dispositionLabels[r.Disposition])
	}
	return fmt.Sprintf("- %s [%s]: %s; %s: %s.", label, need, dispositionLabels[r.Disposition], kindLabels[r.Kind], r.Reason)
}
