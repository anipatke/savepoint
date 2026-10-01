package codehealth

import (
	"errors"
	"fmt"
	"slices"
	"strings"
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
	// MeasuredText is MeasuredAt in the friendly fixed UTC form.
	MeasuredText string
	// Headline is one sentence on how many signals need a look. SignOff says
	// whether the result blocks sign-off, from the same Evaluate as the
	// health check command.
	Headline string
	SignOff  string
	// Recorded is the repository state the snapshot was measured at. Pass it
	// to DashboardFreshness; it is not the snapshot.
	Recorded RepositoryIdentity
	Rows     []DashboardRow
	History  []DashboardHistoryEntry
}

// DashboardRow describes one configured instance, or the placeholder for a
// capability that has none.
type DashboardRow struct {
	Capability     Capability
	CapabilityText string
	Name           string
	NotConfigured  bool
	Label          Classification
	LabelText      string
	Explanation    string
	Trend          string
	Basis          string
	// Spark draws the last official checks, oldest first, one block each; it
	// is empty below three points. SparkWord is better, worse or steady, empty
	// when there is no trend. SparkNote says when history is thin or restarted.
	Spark           string
	SparkWord       string
	SparkNote       string
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
	// SignOff is empty for a manual snapshot and for a signal nobody set up.
	SignOff string

	// Plain-language reading of the signal. All plain strings; Where is empty
	// when there is no evidence.
	Question string
	Value    string
	// Figure is the number alone, for the popover's value column.
	Figure   string
	Aim      string
	Meaning  string
	NextStep string
	Where    string
}

// DashboardHistoryEntry is one earlier snapshot, newest first.
type DashboardHistoryEntry struct {
	CreatedAt   string
	WhenText    string
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
	d.MeasuredText = whenText(newest.CreatedAt)
	d.Headline = headlineText(d.Rows)
	d.applySignOff(cfg, newest, hasOfficial(snaps))
	for i := len(snaps) - 1; i >= 0 && len(d.History) < MaxDashboardHistory; i-- {
		s := snaps[i]
		d.History = append(d.History, DashboardHistoryEntry{
			CreatedAt: s.CreatedAt, WhenText: whenText(s.CreatedAt), Origin: s.Origin, OriginText: originText[s.Origin],
			Overall: s.Summary.Overall, OverallText: classificationText[s.Summary.Overall],
		})
	}
	return d, nil
}

func hasOfficial(snaps []Snapshot) bool {
	return slices.ContainsFunc(snaps, func(s Snapshot) bool { return s.Origin == OriginOfficial })
}

// applySignOff says what the newest snapshot means for sign-off. Only an
// official snapshot is evaluated, with the gate the health check command
// uses; a failed evaluation is said plainly and the rest still loads.
func (d *Dashboard) applySignOff(cfg Config, newest Snapshot, anyOfficial bool) {
	if newest.Origin != OriginOfficial {
		d.SignOff = textSignOffManual
		if !anyOfficial {
			d.SignOff = textSignOffNoOfficial
		}
		return
	}
	verdict, err := Evaluate(newest, cfg)
	if err != nil {
		d.SignOff = textSignOffUnavailable
		for i := range d.Rows {
			d.Rows[i].SignOff = textSignOffUnavailable
		}
		return
	}
	d.SignOff = textSignOffClear
	if verdict.Blocks() {
		d.SignOff = textSignOffBlocks
	}
	byKey := make(map[instanceKey]Disposition, len(verdict.Results))
	for _, rv := range verdict.Results {
		byKey[instanceKey{rv.Capability, rv.Provider, rv.Name}] = rv.Disposition
	}
	for i, row := range d.Rows {
		d.Rows[i].SignOff = rowSignOff[byKey[instanceKey{row.Capability, row.Provider, row.Name}]]
	}
}

// SnapshotLabel is what a screen needs to point at one stored snapshot: its
// overall label in the dashboard's words and a short form of its identity.
type SnapshotLabel struct {
	Label   string
	ShortID string
}

// shortIDLength is how many digest characters ShortID keeps.
const shortIDLength = 8

// SnapshotLabels maps each stored snapshot ID to its SnapshotLabel. It reads
// files only. Storage that is absent or cannot be read yields no labels and no
// error: a screen that only decorates a Check with a health line must not fail
// because health is unavailable, and doctor reports the damage.
func SnapshotLabels(projectPath string) map[string]SnapshotLabel {
	snaps, err := NewStore(projectPath).LoadSnapshots()
	if err != nil || len(snaps) == 0 {
		return nil
	}
	out := make(map[string]SnapshotLabel, len(snaps))
	for _, s := range snaps {
		short := strings.TrimPrefix(s.ID, digestPrefix)
		if len(short) > shortIDLength {
			short = short[:shortIDLength]
		}
		out[s.ID] = SnapshotLabel{Label: classificationText[s.Summary.Overall], ShortID: short}
	}
	return out
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

// BlocksSignOff is true when the gate says this signal blocks sign-off.
func (row DashboardRow) BlocksSignOff() bool { return row.SignOff == textSignOffBlocks }

// withWords copies a signal's plain-language reading onto its row.
func (row DashboardRow) withWords(w signalWords) DashboardRow {
	row.Question, row.Value, row.Figure, row.Aim = w.Question, w.Value, w.Figure, w.Aim
	row.Meaning, row.NextStep, row.Where = w.Meaning, w.NextStep, w.Where
	return row
}

func notConfiguredRow(c Capability) DashboardRow {
	return DashboardRow{
		Capability: c, CapabilityText: capabilityText[c], NotConfigured: true,
		Label: ClassificationUnknown, LabelText: classificationText[ClassificationUnknown],
		Explanation: textNotConfiguredRow,
		Outcome:     OutcomeNotConfigured, OutcomeText: outcomeText[OutcomeNotConfigured],
		RequiredText: textOptional,
	}.withWords(unmeasuredWords(c, nil, true))
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
	}.withWords(unmeasuredWords(cc.Capability, cc.Thresholds, false))
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
	row.Spark, row.SparkWord, row.SparkNote = sparkline(r.Capability, sparkValues(r, origin, s), trend.Direction, s.incompatible > 0)
	if r.Outcome.Measured() {
		row.FreshnessText = freshnessText[r.Freshness]
	}
	row = row.withWords(measuredWords(cc.Thresholds, r, cs.Classification))
	if why, ok := whyLabel(cc.Thresholds, r, cs.Classification, s, trend); ok {
		row.Meaning, row.NextStep = why.Meaning, why.NextStep
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
