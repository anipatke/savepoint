package codehealth

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
)

// Window sizes. A recent range and trend need minObservations; the range looks
// at no more than recentWindow; the decline baseline is the median of the last
// baselineSize earlier observations.
const (
	minObservations = 3
	recentWindow    = 5
	baselineSize    = 3
)

// TrendDirection says how the recent window moved, from the point of view of
// health rather than of the raw number.
type TrendDirection string

const (
	TrendNone      TrendDirection = ""
	TrendImproving TrendDirection = "improving"
	TrendSteady    TrendDirection = "steady"
	TrendDeclining TrendDirection = "declining"
)

// Trend is the recent range and direction. It is zero until at least three
// comparable official observations exist.
type Trend struct {
	Observations int
	HasRange     bool
	Min, Max     float64
	Direction    TrendDirection
}

// series is the comparable official history for one result, plus what was left
// out of it so the explanation can say so.
type series struct {
	values       []float64 // oldest first
	manual       int       // manual observations in the same series
	incompatible int       // official observations in a different series
}

// selectSeries keeps only official, available, finite observations that share
// current's series identity. Manual and incompatible observations stay
// visible through the counts but never enter the values.
func selectSeries(current CapabilityResult, history []HistoryEntry) series {
	var s series
	id := current.SeriesID()
	for _, o := range history {
		r := o.Result
		if r.Capability != current.Capability {
			continue
		}
		sameSeries := r.Outcome.Measured() && r.SeriesID() == id
		switch {
		case o.Origin == OriginManual:
			if sameSeries {
				s.manual++
			}
		case o.Origin != OriginOfficial:
		case !sameSeries:
			if r.Outcome.Measured() {
				s.incompatible++
			}
		case usableValue(r):
			s.values = append(s.values, r.Value.Number)
		}
	}
	return s
}

// usableValue is true for a complete, finite measurement. Partial results
// measure only part of the scope, so they never feed comparisons.
func usableValue(r CapabilityResult) bool {
	return r.Outcome == OutcomeAvailable && r.Value != nil &&
		!math.IsNaN(r.Value.Number) && !math.IsInf(r.Value.Number, 0) && r.Value.Number >= 0
}

// buildTrend computes the recent range and direction over the newest official
// observations, including current when it is itself a complete official one.
func buildTrend(current CapabilityResult, origin Origin, s series) Trend {
	values := slices.Clone(s.values)
	if origin == OriginOfficial && usableValue(current) {
		values = append(values, current.Value.Number)
	}
	if len(values) < minObservations {
		return Trend{}
	}
	if len(values) > recentWindow {
		values = values[len(values)-recentWindow:]
	}
	t := Trend{Observations: len(values), HasRange: true, Min: slices.Min(values), Max: slices.Max(values)}
	worseBy := worseningOf(current.Capability, values[0], values[len(values)-1])
	switch move := materialMovement[current.Capability]; {
	case worseBy >= move:
		t.Direction = TrendDeclining
	case -worseBy >= move:
		t.Direction = TrendImproving
	default:
		t.Direction = TrendSteady
	}
	return t
}

// worseningOf is how much worse `to` is than `from`; negative means better.
// Rounding keeps decimal boundaries such as a 5.0-point move from missing on
// floating-point noise.
func worseningOf(c Capability, from, to float64) float64 {
	d := to - from
	if higherIsBetter(c) {
		d = from - to
	}
	return math.Round(d*1e6) / 1e6
}

type movement struct {
	median             float64
	declined, improved bool
}

// baselineMovement compares current with the median of the last baselineSize
// earlier comparable official observations. It reports nothing until that many
// exist.
func baselineMovement(current CapabilityResult, s series) (movement, bool) {
	if len(s.values) < baselineSize {
		return movement{}, false
	}
	last := slices.Clone(s.values[len(s.values)-baselineSize:])
	slices.Sort(last)
	m := movement{median: last[baselineSize/2]}
	worseBy := worseningOf(current.Capability, m.median, current.Value.Number)
	move := materialMovement[current.Capability]
	m.declined = worseBy >= move
	m.improved = -worseBy >= move
	return m, true
}

// describe returns the history sentences: the recent range and direction once
// enough observations exist, otherwise how much history is missing, plus what
// was left out of comparison.
func (t Trend) describe(current CapabilityResult, s series) []string {
	var out []string
	switch {
	case t.HasRange && t.Min == t.Max:
		out = append(out, fmt.Sprintf("Unchanged at %s over %d official checks.", formatNumber(t.Min, current), t.Observations))
	case t.HasRange:
		out = append(out, fmt.Sprintf("Recent range %s to %s over %d official checks, %s.",
			formatNumber(t.Min, current), formatNumber(t.Max, current), t.Observations, t.Direction))
	default:
		out = append(out, "Not enough comparable official history for a recent range yet.")
	}
	if s.incompatible > 0 {
		out = append(out, fmt.Sprintf("%d earlier official %s not compared.", s.incompatible, pluralize(s.incompatible, "result", "results")))
	}
	if s.manual > 0 {
		out = append(out, fmt.Sprintf("%d manual %s shown, not counted.", s.manual, pluralize(s.manual, "result", "results")))
	}
	return out
}

func pluralize(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}

func join(parts []string) string { return strings.Join(parts, " ") }

// boundText keeps an explanation within the persisted bound without splitting a
// UTF-8 character.
func boundText(s string) string {
	if len(s) <= MaxReasonLen {
		return s
	}
	cut := MaxReasonLen - len("…")
	for cut > 0 && !isRuneStart(s[cut]) {
		cut--
	}
	return s[:cut] + "…"
}

func isRuneStart(b byte) bool { return b&0xC0 != 0x80 }

func formatValue(r CapabilityResult) string { return formatNumber(r.Value.Number, r) }

// formatNumber rounds to one decimal so a reading never implies more precision
// than the measurement has.
func formatNumber(v float64, r CapabilityResult) string {
	s := formatPlain(v)
	if capabilityUnits[r.Capability] == UnitPercent {
		return s + "%"
	}
	return s
}

func formatPlain(v float64) string {
	return strconv.FormatFloat(math.Round(v*10)/10, 'f', -1, 64)
}
