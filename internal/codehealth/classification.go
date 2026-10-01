package codehealth

import (
	"fmt"
	"math"
)

// Detail keys that carry dependency-vulnerability severity counts alongside the
// total value.
const (
	DetailHighVulnerabilities     = "high"
	DetailCriticalVulnerabilities = "critical"
)

// Minimum movement against the recent baseline that counts as material. Smaller
// movement is noise. Counts move by whole units, so any change is material.
var materialMovement = map[Capability]float64{
	CapabilityTests:                   1,
	CapabilityCoverage:                5,
	CapabilityComplexity:              5,
	CapabilityDuplication:             2,
	CapabilityDependencyVulnerability: 1,
}

// defaultThresholds apply when a project configures none. Tests and
// vulnerabilities tolerate no failures or findings to be good; any vulnerability
// that is not a hard blocker is watch-worthy rather than needing attention.
var defaultThresholds = map[Capability]Threshold{
	CapabilityTests:                   {Good: 0, Watch: 0},
	CapabilityCoverage:                {Good: 80, Watch: 60},
	CapabilityComplexity:              {Good: 10, Watch: 20},
	CapabilityDuplication:             {Good: 3, Watch: 5},
	CapabilityDependencyVulnerability: {Good: 0, Watch: math.MaxFloat64},
}

// HistoryEntry is one earlier result of a capability together with the origin of
// the snapshot it came from.
type HistoryEntry struct {
	Origin Origin
	Result CapabilityResult
}

// Assessment is the deterministic label and plain-language explanation for one
// capability result. It carries no numeric score.
type Assessment struct {
	Capability     Capability
	Provider       ProviderKey
	Name           string
	Classification Classification
	Explanation    string
	Trend          Trend
}

// Summary converts the assessment to the persisted summary form.
func (a Assessment) Summary() CapabilitySummary {
	return CapabilitySummary{
		Capability:     a.Capability,
		Provider:       a.Provider,
		Name:           a.Name,
		Classification: a.Classification,
		Explanation:    a.Explanation,
	}
}

// Assess classifies one current result. origin is the snapshot origin of
// current, configured is the project's threshold guidance (nil for defaults),
// and history holds earlier observations of the same instance, oldest first,
// not including current. Only official observations in current's series feed
// baselines, ranges, and trends.
func Assess(origin Origin, current CapabilityResult, configured *Threshold, history []HistoryEntry) Assessment {
	a := Assessment{
		Capability:     current.Capability,
		Provider:       current.Provenance.Provider,
		Name:           current.Name,
		Classification: ClassificationUnknown,
	}
	if reason, ok := unmeasuredReason(current); ok {
		a.Explanation = boundText(reason)
		return a
	}

	series := selectSeries(current, history)
	a.Trend = buildTrend(current, origin, series)

	base, hard, why := classifyValue(current, thresholdFor(current.Capability, configured))
	if why == "" {
		a.Explanation = boundText(fmt.Sprintf("Value %s is not a usable measurement.", formatValue(current)))
		return a
	}
	class := base
	parts := []string{why}

	if hard {
		class = ClassificationNeedsAttention
	} else if move, ok := baselineMovement(current, series); ok {
		switch {
		case move.declined:
			class = worsen(class)
			parts = append(parts, fmt.Sprintf("Materially worse than recent median %s.", formatNumber(move.median, current)))
		case move.improved:
			parts = append(parts, fmt.Sprintf("Improved from recent median %s; current value still decides.", formatNumber(move.median, current)))
		}
	}

	if caveat := confidenceCaveat(current); caveat != "" && class == ClassificationGood {
		class = ClassificationWatch
		parts = append(parts, caveat)
	} else if caveat != "" {
		parts = append(parts, caveat)
	}
	parts = append(parts, a.Trend.describe(current, series)...)
	if class == ClassificationGood && !a.Trend.HasRange {
		class = ClassificationWatch
		parts = append(parts, "Good needs three comparable official checks.")
	}

	a.Classification = class
	a.Explanation = boundText(join(parts))
	return a
}

// Overall combines assessments without scoring: any needs-attention wins, then
// any unknown, then any watch. Good requires at least one result and every
// result good.
func Overall(assessments []Assessment) Classification {
	if len(assessments) == 0 {
		return ClassificationUnknown
	}
	worst := ClassificationGood
	for _, a := range assessments {
		worst = worse(worst, a.Classification)
	}
	return worst
}

// severity orders labels from best to worst. Unknown ranks between watch and
// needs attention: it is never good, but it is not a confirmed problem.
var severity = map[Classification]int{
	ClassificationGood:           0,
	ClassificationWatch:          1,
	ClassificationUnknown:        2,
	ClassificationNeedsAttention: 3,
}

func worse(a, b Classification) Classification {
	if severity[b] > severity[a] {
		return b
	}
	return a
}

// worsen lowers a label by exactly one level. Needs attention is the floor.
func worsen(c Classification) Classification {
	switch c {
	case ClassificationGood:
		return ClassificationWatch
	case ClassificationWatch:
		return ClassificationNeedsAttention
	}
	return c
}

// thresholdFor returns project guidance, or the built-in default.
func thresholdFor(c Capability, configured *Threshold) Threshold {
	if configured != nil {
		return *configured
	}
	return defaultThresholds[c]
}

// unmeasuredReason explains results that carry no usable measurement. Such
// results are never good.
func unmeasuredReason(r CapabilityResult) (string, bool) {
	switch r.Outcome {
	case OutcomeAvailable, OutcomePartial:
		if r.Value == nil {
			return "Measurement is incomplete and carries no value, so health is unknown.", true
		}
		if math.IsNaN(r.Value.Number) || math.IsInf(r.Value.Number, 0) || r.Value.Number < 0 {
			return "Measured value is not a finite number, so health is unknown.", true
		}
		return "", false
	case OutcomeAbsent:
		return "No report was found, so health is unknown.", true
	case OutcomeUnsupported:
		return "This project is not supported by the provider, so health is unknown.", true
	case OutcomeUnavailable:
		return "The provider could not be used, so health is unknown.", true
	case OutcomeNotConfigured:
		return "Not configured, so health is unknown.", true
	case OutcomeFailed:
		return "Collection failed, so health is unknown.", true
	case OutcomeTimedOut:
		return "Collection timed out, so health is unknown.", true
	case OutcomeCancelled:
		return "Collection was cancelled, so health is unknown.", true
	}
	return "Outcome is not recognised, so health is unknown.", true
}

// confidenceCaveat names why a measured result cannot be trusted as current or
// complete. Such a result is never reported as good.
func confidenceCaveat(r CapabilityResult) string {
	switch {
	case r.Outcome == OutcomePartial:
		return "Only part of the evidence was measured, so this is not a full result."
	case r.Freshness == FreshnessStale:
		return "The evidence is stale; re-measure before relying on it."
	case r.Freshness == FreshnessUnknown:
		return "Freshness is unknown; re-measure before relying on it."
	}
	return ""
}

// classifyValue applies hard minimums and then the current-value thresholds. It
// returns the base label, whether a hard blocker forced it, and a sentence. An
// empty sentence means the value could not be classified.
func classifyValue(r CapabilityResult, t Threshold) (Classification, bool, string) {
	if blocked, why := hardBlocker(r); blocked {
		return ClassificationNeedsAttention, true, why
	}
	v := r.Value.Number
	label := fmt.Sprintf("%s %s", r.Capability.label(), formatValue(r))
	within := func(limit float64) bool { return v <= limit }
	if higherIsBetter(r.Capability) {
		within = func(limit float64) bool { return v >= limit }
	}
	switch {
	case within(t.Good):
		return ClassificationGood, false, label + " meets the good threshold."
	case within(t.Watch):
		return ClassificationWatch, false, label + " is within the watch range."
	}
	return ClassificationNeedsAttention, false, label + " is beyond the watch range."
}

// hardBlocker reports the confirmed rules that configuration cannot weaken:
// failing tests, and high or critical vulnerabilities, including findings whose
// severity is unknown.
func hardBlocker(r CapabilityResult) (bool, string) {
	switch r.Capability {
	case CapabilityTests:
		if r.Value.Number > 0 {
			return true, fmt.Sprintf("%s failing; failing tests always need attention.", formatNumber(r.Value.Number, r))
		}
	case CapabilityDependencyVulnerability:
		// Severity is read before the total so contradictory evidence (a zero
		// total beside a positive high or critical count) still blocks.
		high, hasHigh := detail(r, DetailHighVulnerabilities)
		critical, hasCritical := detail(r, DetailCriticalVulnerabilities)
		switch {
		case high > 0 || critical > 0:
			return true, fmt.Sprintf("%s high and %s critical vulnerabilities; these always need attention.", formatPlain(high), formatPlain(critical))
		case r.Value.Number > 0 && (!hasHigh || !hasCritical):
			return true, "Vulnerabilities were found with unknown severity, so they are treated as blocking."
		}
	}
	return false, ""
}

func detail(r CapabilityResult, key string) (float64, bool) {
	for _, d := range r.Details {
		if d.Key == key {
			return d.Number, true
		}
	}
	return 0, false
}

func (c Capability) label() string {
	switch c {
	case CapabilityTests:
		return "Failing tests"
	case CapabilityCoverage:
		return "Coverage"
	case CapabilityComplexity:
		return "Highest complexity"
	case CapabilityDuplication:
		return "Duplication"
	case CapabilityDependencyVulnerability:
		return "Vulnerabilities"
	}
	return string(c)
}
