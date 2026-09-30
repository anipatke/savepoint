// Package codehealth owns the versioned, stack-agnostic records that describe
// project health evidence. Missing, failed, stale, or incomparable evidence is
// always a distinct state and never a healthy value.
package codehealth

// Schema versions of the persisted records. A reader rejects any other value.
const (
	ConfigVersion   = 1
	SnapshotVersion = 1
)

// Capability is one of the five fixed measures. The set is closed: there is no
// plugin or user-defined capability.
type Capability string

const (
	CapabilityTests                   Capability = "tests"
	CapabilityCoverage                Capability = "coverage"
	CapabilityComplexity              Capability = "complexity"
	CapabilityDuplication             Capability = "duplication"
	CapabilityDependencyVulnerability Capability = "dependency_vulnerabilities"
)

// Capabilities lists the fixed set in canonical order.
func Capabilities() []Capability {
	return []Capability{
		CapabilityTests,
		CapabilityCoverage,
		CapabilityComplexity,
		CapabilityDuplication,
		CapabilityDependencyVulnerability,
	}
}

// Outcome says what collection produced. It is independent of Freshness.
type Outcome string

const (
	OutcomeAvailable     Outcome = "available"
	OutcomePartial       Outcome = "partial"
	OutcomeAbsent        Outcome = "absent"
	OutcomeUnsupported   Outcome = "unsupported"
	OutcomeUnavailable   Outcome = "unavailable"
	OutcomeNotConfigured Outcome = "not_configured"
	OutcomeFailed        Outcome = "failed"
	OutcomeTimedOut      Outcome = "timed_out"
	OutcomeCancelled     Outcome = "cancelled"
)

// Outcomes lists every outcome in canonical order.
func Outcomes() []Outcome {
	return []Outcome{
		OutcomeAvailable, OutcomePartial, OutcomeAbsent, OutcomeUnsupported,
		OutcomeUnavailable, OutcomeNotConfigured, OutcomeFailed,
		OutcomeTimedOut, OutcomeCancelled,
	}
}

// Measured reports whether the outcome may carry measured values.
func (o Outcome) Measured() bool {
	return o == OutcomeAvailable || o == OutcomePartial
}

// Freshness says whether measured evidence still matches the repository and
// configuration. Evidence that was never measured is always FreshnessUnknown.
type Freshness string

const (
	FreshnessFresh   Freshness = "fresh"
	FreshnessStale   Freshness = "stale"
	FreshnessUnknown Freshness = "unknown"
)

// Origin records who produced a snapshot. Official snapshots come from a Full
// Objective Check; manual snapshots come from an explicit refresh.
type Origin string

const (
	OriginOfficial Origin = "official"
	OriginManual   Origin = "manual"
)

// Retention is how long a snapshot is kept. It is fixed by Origin.
type Retention string

const (
	RetentionPermanent Retention = "permanent"
	RetentionPrunable  Retention = "prunable"
)

// Classification is a deterministic summary label. There is no numeric score.
type Classification string

const (
	ClassificationGood           Classification = "good"
	ClassificationWatch          Classification = "watch"
	ClassificationNeedsAttention Classification = "needs_attention"
	ClassificationUnknown        Classification = "unknown"
)

// Unit names what a Value measures. Each capability accepts exactly one unit.
type Unit string

const (
	UnitCount   Unit = "count"
	UnitPercent Unit = "percent"
	UnitCCN     Unit = "ccn"
)

// capabilityUnits is the single source of which unit each capability measures.
var capabilityUnits = map[Capability]Unit{
	CapabilityTests:                   UnitCount,   // failed tests
	CapabilityCoverage:                UnitPercent, // covered statements
	CapabilityComplexity:              UnitCCN,     // highest function CCN
	CapabilityDuplication:             UnitPercent, // duplicated lines
	CapabilityDependencyVulnerability: UnitCount,   // known vulnerabilities
}

// higherIsBetter is true only for coverage; every other measure improves as it
// falls.
func higherIsBetter(c Capability) bool { return c == CapabilityCoverage }

// ProviderKey names one supported report path from the approved catalogue.
type ProviderKey string

const (
	ProviderGoTestJSON     ProviderKey = "go-test-json"
	ProviderVitestJUnit    ProviderKey = "vitest-junit"
	ProviderPytestJUnit    ProviderKey = "pytest-junit"
	ProviderGoCoverProfile ProviderKey = "go-cover-profile"
	ProviderVitestV8       ProviderKey = "vitest-v8-coverage"
	ProviderCoveragePyJSON ProviderKey = "coverage-py-json"
	ProviderLizardCSV      ProviderKey = "lizard-csv"
	ProviderJscpdJSON      ProviderKey = "jscpd-json"
	ProviderOSVScannerJSON ProviderKey = "osv-scanner-json"
)

// providerCapability maps each catalogue key to the one capability it serves.
var providerCapability = map[ProviderKey]Capability{
	ProviderGoTestJSON:     CapabilityTests,
	ProviderVitestJUnit:    CapabilityTests,
	ProviderPytestJUnit:    CapabilityTests,
	ProviderGoCoverProfile: CapabilityCoverage,
	ProviderVitestV8:       CapabilityCoverage,
	ProviderCoveragePyJSON: CapabilityCoverage,
	ProviderLizardCSV:      CapabilityComplexity,
	ProviderJscpdJSON:      CapabilityDuplication,
	ProviderOSVScannerJSON: CapabilityDependencyVulnerability,
}

// Bounds on persisted records, so history stays compact and reviewable.
const (
	MaxReasonLen     = 200
	MaxNoteLen       = 200
	MaxTokenLen      = 64
	MaxPathLen       = 256
	MaxDetails       = 32
	MaxEvidence      = 20
	MaxScopeEntries  = 64
	MaxResults       = 32
	MaxArgs          = 64
	MaxArgLen        = 512
	MaxTimeoutSecond = 120
)
