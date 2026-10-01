package codehealth

import (
	"errors"
	"fmt"
	"slices"
)

// DashboardState says which of the three screens the board should show.
type DashboardState string

const (
	DashboardNotConfigured DashboardState = "not_configured"
	DashboardFirstRun      DashboardState = "first_run"
	DashboardMeasured      DashboardState = "measured"
)

// MaxDashboardHistory bounds the history list.
const MaxDashboardHistory = 10

// Dashboard is the read-only description of current health the board renders.
// It holds plain strings and the closed vocabularies only, so the board never
// handles snapshots, results, or provider types and never decides a label.
type Dashboard struct {
	State DashboardState

	// Set only when State is DashboardMeasured.
	Overall     Classification
	OverallText string
	Origin      Origin
	OriginText  string
	SnapshotID  string
	MeasuredAt  string
	// Recorded is the repository state the snapshot was measured at. Pass it
	// to DashboardFreshness; it is not the snapshot.
	Recorded RepositoryIdentity
	Rows     []DashboardRow
	History  []DashboardHistoryEntry
}

// DashboardRow describes one configured instance, or the placeholder for a
// capability that has none.
type DashboardRow struct {
	Capability      Capability
	CapabilityText  string
	Name            string
	NotConfigured   bool
	Label           Classification
	LabelText       string
	Explanation     string
	Trend           string
	Basis           string
	Outcome         Outcome
	OutcomeText     string
	FreshnessText   string
	Required        bool
	RequiredText    string
	Provider        ProviderKey
	ProviderVersion string
	Scope           []string
	CollectedAt     string
	Evidence        []EvidenceRef
}

// DashboardHistoryEntry is one earlier snapshot, newest first.
type DashboardHistoryEntry struct {
	CreatedAt   string
	Origin      Origin
	OriginText  string
	Overall     Classification
	OverallText string
}

// Every phrase the dashboard says lives here (STYLE-09). Classification
// wording comes from the persisted summary, never from this table.
var (
	classificationText = map[Classification]string{
		ClassificationGood:           "Good",
		ClassificationWatch:          "Watch",
		ClassificationNeedsAttention: "Needs Attention",
		ClassificationUnknown:        "Unknown",
	}
	originText = map[Origin]string{
		OriginOfficial: "Official check",
		OriginManual:   "Manual refresh",
	}
	outcomeText = map[Outcome]string{
		OutcomeAvailable:     "measured",
		OutcomePartial:       "partial",
		OutcomeAbsent:        "absent: no report was found",
		OutcomeUnsupported:   "unsupported",
		OutcomeUnavailable:   "unavailable: the tool could not be used",
		OutcomeNotConfigured: "not configured",
		OutcomeFailed:        "failed",
		OutcomeTimedOut:      "timed out",
		OutcomeCancelled:     "cancelled",
	}
	freshnessText = map[Freshness]string{
		FreshnessFresh:   "fresh",
		FreshnessStale:   "stale",
		FreshnessUnknown: "freshness unknown",
	}
	capabilityText = map[Capability]string{
		CapabilityTests:                   "Tests",
		CapabilityCoverage:                "Coverage",
		CapabilityComplexity:              "Complexity",
		CapabilityDuplication:             "Duplication",
		CapabilityDependencyVulnerability: "Dependency vulnerabilities",
	}
	trendText = map[TrendDirection]string{
		TrendNone:      "No trend yet",
		TrendImproving: "Improving",
		TrendSteady:    "Steady",
		TrendDeclining: "Declining",
	}
)

const (
	textRequired         = "required"
	textOptional         = "optional"
	textNotConfiguredRow = "Not configured, so health for this signal is unknown."
	textNoResultRow      = "This snapshot has no result for this instance; refresh to measure it."
	textNoBasis          = "No comparable official history yet."
	textBasisFormat      = "Compared with %d earlier comparable official %s."
	textTrendRange       = "%s over %d official checks, %s to %s"
)

// CapabilityText is the plain name of a signal, for screens that only hold the
// closed Capability value (a refresh in progress).
func CapabilityText(c Capability) string {
	if text, ok := capabilityText[c]; ok {
		return text
	}
	return string(c)
}

// LoadDashboard reads the saved configuration and snapshots and describes the
// newest snapshot of either origin. It runs no subprocess and no tool. A
// damaged configuration or history is returned as the storage error; nothing is
// repaired.
func LoadDashboard(root string) (Dashboard, error) {
	store := NewStore(root)
	cfg, err := store.LoadConfig()
	if errors.Is(err, ErrConfigNotFound) {
		return Dashboard{State: DashboardNotConfigured}, nil
	}
	if err != nil {
		return Dashboard{}, err
	}
	snaps, err := store.LoadSnapshots()
	if err != nil {
		return Dashboard{}, err
	}
	if len(snaps) == 0 {
		return Dashboard{State: DashboardFirstRun}, nil
	}
	newest := snaps[len(snaps)-1]
	d := Dashboard{
		State:       DashboardMeasured,
		Overall:     newest.Summary.Overall,
		OverallText: classificationText[newest.Summary.Overall],
		Origin:      newest.Origin,
		OriginText:  originText[newest.Origin],
		SnapshotID:  newest.ID,
		MeasuredAt:  newest.CreatedAt,
		Recorded:    newest.Repository,
		Rows:        dashboardRows(cfg, newest, dashboardHistory(snaps[:len(snaps)-1])),
	}
	for i := len(snaps) - 1; i >= 0 && len(d.History) < MaxDashboardHistory; i-- {
		s := snaps[i]
		d.History = append(d.History, DashboardHistoryEntry{
			CreatedAt: s.CreatedAt, Origin: s.Origin, OriginText: originText[s.Origin],
			Overall: s.Summary.Overall, OverallText: classificationText[s.Summary.Overall],
		})
	}
	return d, nil
}

// dashboardHistory groups earlier results the way Collect does, so trends
// match what the snapshot was classified against.
func dashboardHistory(earlier []Snapshot) map[instanceKey][]HistoryEntry {
	out := make(map[instanceKey][]HistoryEntry)
	for _, s := range earlier {
		for _, r := range s.Results {
			out[r.historyKey()] = append(out[r.historyKey()], HistoryEntry{Origin: s.Origin, Result: r})
		}
	}
	return out
}

// dashboardRows lists the five signals in canonical order: each configured
// instance in configuration order, or one placeholder for a capability with
// none.
func dashboardRows(cfg Config, snap Snapshot, history map[instanceKey][]HistoryEntry) []DashboardRow {
	results := make(map[instanceKey]CapabilityResult, len(snap.Results))
	for _, r := range snap.Results {
		results[r.key()] = r
	}
	summaries := make(map[instanceKey]CapabilitySummary, len(snap.Summary.Capabilities))
	for _, cs := range snap.Summary.Capabilities {
		summaries[instanceKey{cs.Capability, cs.Provider, cs.Name}] = cs
	}
	var rows []DashboardRow
	for _, c := range Capabilities() {
		var instances []CapabilityConfig
		for _, cc := range cfg.Capabilities {
			if cc.Capability == c {
				instances = append(instances, cc)
			}
		}
		if len(instances) == 0 {
			rows = append(rows, notConfiguredRow(c))
			continue
		}
		for _, cc := range instances {
			key := instanceKey{cc.Capability, cc.Provider, cc.Name}
			r, ok := results[key]
			if !ok {
				rows = append(rows, noResultRow(cc))
				continue
			}
			rows = append(rows, measuredRow(snap.Origin, cc, r, summaries[key], history[r.historyKey()]))
		}
	}
	return rows
}

func notConfiguredRow(c Capability) DashboardRow {
	return DashboardRow{
		Capability: c, CapabilityText: capabilityText[c], NotConfigured: true,
		Label: ClassificationUnknown, LabelText: classificationText[ClassificationUnknown],
		Explanation: textNotConfiguredRow,
		Outcome:     OutcomeNotConfigured, OutcomeText: outcomeText[OutcomeNotConfigured],
		RequiredText: textOptional,
	}
}

// noResultRow covers an instance configured after the snapshot was taken.
func noResultRow(cc CapabilityConfig) DashboardRow {
	return DashboardRow{
		Capability: cc.Capability, CapabilityText: capabilityText[cc.Capability], Name: cc.Name,
		Label: ClassificationUnknown, LabelText: classificationText[ClassificationUnknown],
		Explanation: textNoResultRow,
		Outcome:     OutcomeNotConfigured, OutcomeText: outcomeText[OutcomeNotConfigured],
		Required: cc.Required, RequiredText: requiredText(cc.Required),
		Provider: cc.Provider, Scope: slices.Clone(cc.Scope),
	}
}

func measuredRow(origin Origin, cc CapabilityConfig, r CapabilityResult, cs CapabilitySummary, history []HistoryEntry) DashboardRow {
	s := selectSeries(r, ownHistory(r, history))
	trend := buildTrend(r, origin, s)
	row := DashboardRow{
		Capability: r.Capability, CapabilityText: capabilityText[r.Capability], Name: r.Name,
		Label: cs.Classification, LabelText: classificationText[cs.Classification],
		Explanation: cs.Explanation,
		Trend:       trendWords(r, trend), Basis: basisWords(s),
		Outcome: r.Outcome, OutcomeText: outcomeText[r.Outcome],
		Required: cc.Required, RequiredText: requiredText(cc.Required),
		Provider: r.Provenance.Provider, ProviderVersion: r.Provenance.ProviderVersion,
		Scope: slices.Clone(r.Scope), CollectedAt: r.CollectedAt,
		Evidence: slices.Clone(r.Evidence),
	}
	if r.Outcome.Measured() {
		row.FreshnessText = freshnessText[r.Freshness]
	}
	return row
}

func requiredText(required bool) string {
	if required {
		return textRequired
	}
	return textOptional
}

func trendWords(r CapabilityResult, t Trend) string {
	if !t.HasRange {
		return trendText[TrendNone]
	}
	return fmt.Sprintf(textTrendRange, trendText[t.Direction], t.Observations, formatNumber(t.Min, r), formatNumber(t.Max, r))
}

// basisWords says what the trend was compared with, and what was left out.
func basisWords(s series) string {
	out := textNoBasis
	if n := len(s.values); n > 0 {
		out = fmt.Sprintf(textBasisFormat, n, pluralize(n, "check", "checks"))
	}
	if s.incompatible > 0 {
		out += fmt.Sprintf(" %d earlier official %s not compared.", s.incompatible, pluralize(s.incompatible, "result", "results"))
	}
	if s.manual > 0 {
		out += fmt.Sprintf(" %d manual %s shown, not counted.", s.manual, pluralize(s.manual, "result", "results"))
	}
	return out
}
