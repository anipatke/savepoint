package codehealth

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// signalWords is the plain-language reading of one dashboard row. It is derived
// from a result and the configured thresholds only; the stored classification
// text, thresholds and snapshots are never touched.
type signalWords struct {
	Question string
	Value    string
	Aim      string
	Meaning  string
	NextStep string
	Where    string
}

// copyEntry is what a label means for one signal and the one thing to do next.
type copyEntry struct {
	Meaning  string
	NextStep string
}

// Detail keys the tests reader writes beside its failing count.
const (
	detailTotalTests  = "total_tests"
	detailPassedTests = "passed_tests"
)

// Every phrase of the plain-language wording lives in these tables (STYLE-09).
var (
	questionText = map[Capability]string{
		CapabilityTests:                   "Do the tests pass?",
		CapabilityCoverage:                "How much code do the tests actually run?",
		CapabilityComplexity:              "How tangled is the hardest code?",
		CapabilityDuplication:             "How much is copy-pasted?",
		CapabilityDependencyVulnerability: "Known security problems in libraries we use",
	}
	// zeroAimText is the aim when the Good threshold is zero.
	zeroAimText = map[Capability]string{
		CapabilityTests:                   "aim for none failing",
		CapabilityDependencyVulnerability: "aim for none",
	}
	meaningText = map[Capability]map[Classification]copyEntry{
		CapabilityTests: {
			ClassificationGood:           {"Every test passed in the last check.", "Nothing to do; keep the tests passing."},
			ClassificationWatch:          {"The test result is incomplete or out of date.", "Run the health check again for a full result."},
			ClassificationNeedsAttention: {"Some tests fail, so the code does not do everything it is expected to.", "Ask your agent to fix the failing tests."},
			ClassificationUnknown:        unknownCopy,
		},
		CapabilityCoverage: {
			ClassificationGood:           {"The tests run most of the code, so a break is likely to be caught.", "Nothing to do; keep adding tests with new code."},
			ClassificationWatch:          {"The tests skip a fair amount of the code, so some breaks could go unnoticed.", "Ask your agent to add tests for the least-tested files."},
			ClassificationNeedsAttention: {"Much of the code is never run by the tests, so breaks can slip through unnoticed.", "Ask your agent to add tests for the least-tested files."},
			ClassificationUnknown:        unknownCopy,
		},
		CapabilityComplexity: {
			ClassificationGood:           {"Even the most tangled function is still easy to follow.", "Nothing to do; keep functions small."},
			ClassificationWatch:          {"The most tangled function is getting hard to follow.", "Ask your agent to split the most tangled function."},
			ClassificationNeedsAttention: {"The most tangled function is very hard to follow and easy to break.", "Ask your agent to split the most tangled function."},
			ClassificationUnknown:        unknownCopy,
		},
		CapabilityDuplication: {
			ClassificationGood:           {"Very little code is copy-pasted.", "Nothing to do; reuse code instead of copying it."},
			ClassificationWatch:          {"A fair amount of code is copy-pasted, so a fix may need repeating.", "Ask your agent to merge the copied code into shared code."},
			ClassificationNeedsAttention: {"A lot of code is copy-pasted, so a fix may need repeating in many places.", "Ask your agent to merge the copied code into shared code."},
			ClassificationUnknown:        unknownCopy,
		},
		CapabilityDependencyVulnerability: {
			ClassificationGood:           {"No known security problems in the libraries you use.", "Nothing to do; update libraries regularly."},
			ClassificationWatch:          {"Some libraries have known low-risk problems.", "Ask your agent to update the affected libraries."},
			ClassificationNeedsAttention: {"A library has a serious or unrated security problem.", "Ask your agent to update or replace the affected library now."},
			ClassificationUnknown:        unknownCopy,
		},
	}
	// notConfiguredCopy is the reading of a signal nobody has set up.
	notConfiguredCopy = copyEntry{
		"This signal is not set up, so nothing is judged.",
		"Run savepoint health setup to add a tool for it.",
	}
	unknownCopy = copyEntry{
		"The last check gave no usable result, so this signal is not judged.",
		"Run the health check again; if it still fails, ask your agent to fix the tool setup.",
	}
)

// Why a label differs from what the value alone says (STYLE-09). Each cause
// follows a lead sentence saying where the value stands.
const (
	textLeadGood     = "Meets the aim"
	textLeadWatch    = "Within the watch range"
	textCauseWorse   = ", but it is clearly worse than recent checks."
	textCausePartial = ", but only part was measured."
	textCauseStale   = ", but the result may be out of date."
	textCauseThin    = "; it stays Watch until three comparable checks."
	textNextThin     = "Nothing to fix; it settles after a few more official checks."
	textNextAgain    = "Run the health check again for a full result."
)

const (
	textNotMeasured     = "not measured"
	textNoTestsRan      = "no tests ran"
	textNoneFailing     = "none failing"
	textNone            = "none"
	textAllPass         = "all %s pass"
	textFailing         = "%s failing"
	textHardest         = "hardest function scores %s"
	textCopyPasted      = "%s copy-pasted"
	textFound           = "%s found"
	textUnrated         = "unrated"
	textAimMore         = "aim for %s or more"
	textAimLess         = "aim for %s or less"
	textPastWatch       = " It is past the watch line of %s."
	textWhereMany       = "%d files; savepoint health check lists them"
	textBoundMore       = "%s or more"
	textBoundLess       = "%s or less"
	textThousandsSep    = ","
	vulnerabilityJoiner = ", "
)

// aimText says the Good threshold as a yardstick, from the instance's
// configured thresholds or else the built-in defaults.
func aimText(c Capability, configured *Threshold) string {
	good := thresholdFor(c, configured).Good
	if good == 0 && !higherIsBetter(c) {
		if text, ok := zeroAimText[c]; ok {
			return text
		}
	}
	if higherIsBetter(c) {
		return fmt.Sprintf(textAimMore, boundNumber(c, good))
	}
	return fmt.Sprintf(textAimLess, boundNumber(c, good))
}

func boundNumber(c Capability, v float64) string {
	s := formatPlain(v)
	if capabilityUnits[c] == UnitPercent {
		return s + "%"
	}
	return s
}

// watchLineText names the Watch boundary, only for a value beyond it. Signals
// whose Good and Watch coincide, or whose Watch is unbounded, have no line.
func watchLineText(c Capability, t Threshold, value float64) string {
	if t.Watch == t.Good || t.Watch >= math.MaxFloat64 {
		return ""
	}
	beyond := value > t.Watch
	bound := textBoundLess
	if higherIsBetter(c) {
		beyond = value < t.Watch
		bound = textBoundMore
	}
	if !beyond {
		return ""
	}
	return fmt.Sprintf(textPastWatch, fmt.Sprintf(bound, boundNumber(c, t.Watch)))
}

// measuredWords reads a measured result in plain words.
func measuredWords(configured *Threshold, r CapabilityResult, label Classification) signalWords {
	w := signalWords{
		Question: questionText[r.Capability],
		Aim:      aimText(r.Capability, configured),
		Where:    whereText(r.Evidence),
	}
	if r.Value == nil || !r.Outcome.Measured() {
		w.Value = textNotMeasured
		w.Meaning, w.NextStep = unknownCopy.Meaning, unknownCopy.NextStep
		return w
	}
	w.Value = valueText(r)
	entry := meaningText[r.Capability][label]
	w.Meaning, w.NextStep = entry.Meaning, entry.NextStep
	if label == ClassificationNeedsAttention {
		w.Meaning += watchLineText(r.Capability, thresholdFor(r.Capability, configured), r.Value.Number)
	}
	return w
}

// whyLabel explains a label that is worse than the value alone earns. The
// value's own reading says where it stands; the cause says what moved the
// label: a drop from recent checks, a partial or out-of-date result, or too
// few comparable checks. ok is false when the label is just the value's own.
func whyLabel(configured *Threshold, r CapabilityResult, label Classification, s series, t Trend) (words signalWords, ok bool) {
	if r.Value == nil || !r.Outcome.Measured() || label == ClassificationUnknown {
		return signalWords{}, false
	}
	base, hard, why := classifyValue(r, thresholdFor(r.Capability, configured))
	if why == "" || hard || severity[label] <= severity[base] {
		return signalWords{}, false
	}
	lead := textLeadWatch
	if base == ClassificationGood {
		lead = textLeadGood
	}
	move, hasMove := baselineMovement(r, s)
	switch {
	case hasMove && move.declined && worsen(base) == label:
		return signalWords{Meaning: lead + textCauseWorse, NextStep: meaningText[r.Capability][label].NextStep}, true
	case base != ClassificationGood:
		return signalWords{}, false
	case r.Outcome == OutcomePartial:
		return signalWords{Meaning: lead + textCausePartial, NextStep: textNextAgain}, true
	case confidenceCaveat(r) != "":
		return signalWords{Meaning: lead + textCauseStale, NextStep: textNextAgain}, true
	case !t.HasRange:
		return signalWords{Meaning: lead + textCauseThin, NextStep: textNextThin}, true
	}
	return signalWords{}, false
}

// unmeasuredWords reads a signal with no result: nothing configured, or an
// instance the snapshot does not cover.
func unmeasuredWords(c Capability, configured *Threshold, notConfigured bool) signalWords {
	entry := unknownCopy
	if notConfigured {
		entry = notConfiguredCopy
	}
	return signalWords{
		Question: questionText[c], Value: textNotMeasured, Aim: aimText(c, configured),
		Meaning: entry.Meaning, NextStep: entry.NextStep,
	}
}

// valueText says the measured value as words.
func valueText(r CapabilityResult) string {
	n := r.Value.Number
	switch r.Capability {
	case CapabilityTests:
		return testsValue(r)
	case CapabilityCoverage:
		return formatNumber(n, r)
	case CapabilityComplexity:
		return fmt.Sprintf(textHardest, formatPlain(n))
	case CapabilityDuplication:
		return fmt.Sprintf(textCopyPasted, formatNumber(n, r))
	case CapabilityDependencyVulnerability:
		return vulnerabilityValue(r)
	}
	return formatValue(r)
}

func testsValue(r CapabilityResult) string {
	if failed := r.Value.Number; failed > 0 {
		return fmt.Sprintf(textFailing, groupDigits(failed))
	}
	total, hasTotal := detail(r, detailTotalTests)
	switch {
	case !hasTotal:
		return textNoneFailing
	case total == 0:
		return textNoTestsRan
	}
	return fmt.Sprintf(textAllPass, groupDigits(total))
}

// vulnerabilityValue lists the findings by severity, worst first.
func vulnerabilityValue(r CapabilityResult) string {
	var parts []string
	for _, key := range vulnerabilitySeverityKeys {
		count, ok := detail(r, key)
		if !ok || count <= 0 {
			continue
		}
		name := key
		if key == DetailUnknownVulnerabilities {
			name = textUnrated
		}
		parts = append(parts, formatPlain(count)+" "+name)
	}
	switch {
	case len(parts) > 0:
		return strings.Join(parts, vulnerabilityJoiner)
	case r.Value.Number > 0:
		return fmt.Sprintf(textFound, formatPlain(r.Value.Number))
	}
	return textNone
}

// whereText names the affected file when the evidence covers one path, and
// otherwise only how many files. Stored evidence is in path order, so no file
// is ever called the top or worst.
func whereText(evidence []EvidenceRef) string {
	paths := map[string]bool{}
	for _, e := range evidence {
		paths[e.Path] = true
	}
	switch len(paths) {
	case 0:
		return ""
	case 1:
		return evidence[0].Path
	}
	return fmt.Sprintf(textWhereMany, len(paths))
}

// groupDigits writes a whole number with thousands separators.
func groupDigits(v float64) string {
	s := strconv.FormatInt(int64(math.Round(v)), 10)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + textThousandsSep + s[i:]
	}
	return s
}

// Sign-off and headline wording (STYLE-09).
const (
	textSignOffBlocks      = "Blocks sign-off"
	textSignOffAdvisory    = "Advisory only"
	textSignOffClear       = "Doesn't block sign-off"
	textSignOffManual      = "A manual refresh does not affect sign-off."
	textSignOffNoOfficial  = "There is no official check yet, so nothing is judged for sign-off."
	textSignOffUnavailable = "Sign-off status unavailable"
	textHeadlineNone       = "Not enough to judge yet"
	textHeadlineAllFine    = "All %d look fine"
	textHeadlineSomeLook   = "%d of %d %s a look"
	whenLayout             = "2 Jan 15:04 UTC"
)

// rowSignOff maps the gate's disposition for a signal to its row wording. A
// signal nobody set up has no sign-off statement.
var rowSignOff = map[Disposition]string{
	DispositionBlocks:   textSignOffBlocks,
	DispositionReported: textSignOffAdvisory,
}

// headlineText counts the signals that are not Good. With nothing judged at all
// it says so rather than claiming all is fine.
func headlineText(rows []DashboardRow) string {
	look, judged := 0, 0
	for _, row := range rows {
		if row.Label != ClassificationGood {
			look++
		}
		if row.Label != ClassificationUnknown {
			judged++
		}
	}
	switch {
	case judged == 0:
		return textHeadlineNone
	case look == 0:
		return fmt.Sprintf(textHeadlineAllFine, len(rows))
	}
	return fmt.Sprintf(textHeadlineSomeLook, look, len(rows), pluralize(look, "needs", "need"))
}

// whenText writes a stored time in one fixed UTC form, independent of the
// machine's clock and zone. An unreadable time is shown as stored.
func whenText(stored string) string {
	t, err := time.Parse(timestampLayout, stored)
	if err != nil {
		return stored
	}
	return t.UTC().Format(whenLayout)
}
