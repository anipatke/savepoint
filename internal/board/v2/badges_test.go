package v2

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	xansi "github.com/charmbracelet/x/ansi"
	"github.com/muesli/termenv"
	"github.com/opencode/savepoint/internal/data"
)

// clearanceStates and blockerKinds are the whole typed vocabularies the badge
// mapping covers. Listing them here is what makes the exhaustiveness tests
// below fail when data grows a value the board has no wording for.
var clearanceStates = []data.ClearanceState{
	data.ClearanceMissing,
	data.ClearanceNeedsWork,
	data.ClearanceStale,
	data.ClearanceUnknown,
	data.ClearanceCurrent,
}

var blockerKinds = []data.GateBlockKind{
	data.GateBlockReplan,
	data.GateBlockDependency,
	data.GateBlockObjectiveDependency,
	data.GateBlockOwnerAcceptance,
	data.GateBlockCheckerAuthority,
	data.GateBlockClearanceMissing,
	data.GateBlockClearanceNeedsWork,
	data.GateBlockClearanceStale,
	data.GateBlockClearanceUnknown,
	data.GateBlockInvalidState,
}

// taskReviewOutcomeCases is every path taskReviewOutcomeBadge's precedence
// covers, paired with the exact card text O012 fixes: owner risk acceptance
// (exception) first, an owner's recorded Check waiver second, and the
// resolved clearance state otherwise, with stale and unknown folded into one
// "REVIEW" wording.
var taskReviewOutcomeCases = []struct {
	name      string
	clearance data.ClearanceState
	waived    bool
	exception bool
	want      string
}{
	{"missing", data.ClearanceMissing, false, false, "[ ] CHECK"},
	{"current", data.ClearanceCurrent, false, false, "[✓] CHECK"},
	{"needs_work", data.ClearanceNeedsWork, false, false, "[!] NEEDS WORK"},
	{"stale", data.ClearanceStale, false, false, "[!] REVIEW"},
	{"unknown", data.ClearanceUnknown, false, false, "[!] REVIEW"},
	{"waived", data.ClearanceMissing, true, false, "[✓] WAIVED"},
	{"owner accepted", data.ClearanceNeedsWork, false, true, "[✓] OWNER ACCEPTED"},
	{"owner accepted outranks waived", data.ClearanceMissing, true, true, "[✓] OWNER ACCEPTED"},
}

func TestTaskCheckBadgeRendersEveryStateDistinctly(t *testing.T) {
	seen := map[string]data.ClearanceState{}
	for _, state := range clearanceStates {
		badge := taskCheckBadge(state, false)
		if badge.Glyph == "" || badge.Label == "" {
			t.Errorf("clearance %q has no glyph or label: %+v", state, badge)
		}
		if other, ok := seen[badge.Text()]; ok {
			t.Errorf("clearance %q and %q both render as %q", state, other, badge.Text())
		}
		seen[badge.Text()] = state
	}
}

// TestTaskCheckBadgeWaivedOutranksClearance proves a recorded waiver always
// reads as "waived", whatever the underlying clearance state happens to be —
// waived is only reachable from ClearanceMissing in practice, but the badge
// itself does not depend on that to stay honest.
func TestTaskCheckBadgeWaivedOutranksClearance(t *testing.T) {
	waived := taskCheckBadge(data.ClearanceMissing, true)
	notWaived := taskCheckBadge(data.ClearanceMissing, false)
	if waived.Text() == notWaived.Text() {
		t.Errorf("a waived Task and one simply not yet checked both render as %q", waived.Text())
	}
	if waived.Glyph != glyphCheckFlagged {
		t.Errorf("waived badge glyph = %q, want the flagged glyph", waived.Glyph)
	}
}

// TestObjectiveCheckBadgeIsThreeNotchOnly proves the Objective badge collapses
// needs_work, stale, and unverified into one "checked, not clear" wording —
// deliberately simpler than taskCheckBadge's five states, per the owner's
// call that the extra distinctions are not worth the reader's attention here.
func TestObjectiveCheckBadgeIsThreeNotchOnly(t *testing.T) {
	seen := map[string]bool{}
	for _, state := range clearanceStates {
		badge := objectiveCheckBadge(state, false)
		seen[badge.Text()] = true
	}
	if len(seen) != 3 {
		t.Errorf("objectiveCheckBadge produced %d distinct renderings across every clearance state, want exactly 3 (missing, needs work, current)", len(seen))
	}
}

// TestObjectiveCheckBadgeFoldsExceptionIntoCurrent proves a recorded owner
// exception reads exactly like a current Check on the compact badge — the
// row does not carry a fourth "BY EXCEPTION" notch. That fact still lives in
// the detail overlay's EXCEPTION section; the row only needs to say the
// Objective is not blocked on its Check anymore.
func TestObjectiveCheckBadgeFoldsExceptionIntoCurrent(t *testing.T) {
	current := objectiveCheckBadge(data.ClearanceCurrent, false)
	byException := objectiveCheckBadge(data.ClearanceNeedsWork, true)

	if byException.Text() != current.Text() {
		t.Errorf("an Objective accepted by exception reads %q, want the same wording as a current Check %q", byException.Text(), current.Text())
	}

	needsWork := objectiveCheckBadge(data.ClearanceNeedsWork, false)
	if byException.Text() == needsWork.Text() {
		t.Errorf("exception acceptance did not change the badge from plain needs-work: both render as %q", byException.Text())
	}
}

func TestStageBadgeIsPresentOnlyWhileInProgress(t *testing.T) {
	stages := []struct {
		stage data.ProgressStage
		label string
	}{
		{data.StageBuild, "BUILD"},
		{data.StageTest, "TEST"},
		{data.StageAudit, "CHECK"},
	}

	seen := map[string]bool{}
	for _, test := range stages {
		badge, ok := stageBadge(data.ColumnInProgress, test.stage)
		if !ok {
			t.Fatalf("stage %q has no badge while in progress", test.stage)
		}
		if badge.Label != test.label {
			t.Errorf("stage %q label = %q, want %q", test.stage, badge.Label, test.label)
		}
		if seen[badge.Text()] {
			t.Errorf("stage %q repeats another stage's rendering %q", test.stage, badge.Text())
		}
		seen[badge.Text()] = true
	}

	for _, status := range []data.ColumnType{data.ColumnPlanned, data.ColumnDone} {
		if _, ok := stageBadge(status, data.StageBuild); ok {
			t.Errorf("status %q carries a stage badge; stage belongs to a Task under way", status)
		}
	}
}

// TestBlockerBadgeCoversEveryKind proves the mapping has an answer for every
// blocker kind: wording for the ones only it can state, and a deliberate
// silence for the clearance kinds — checker_authority included — the review
// outcome badge already states (see taskReviewOutcomeBadge).
func TestBlockerBadgeCoversEveryKind(t *testing.T) {
	stated := map[data.GateBlockKind]bool{
		data.GateBlockReplan:              true,
		data.GateBlockDependency:          true,
		data.GateBlockObjectiveDependency: true,
		data.GateBlockOwnerAcceptance:     true,
	}

	seen := map[string]data.GateBlockKind{}
	for _, kind := range blockerKinds {
		badge, ok := blockerBadge(data.GateBlocker{Kind: kind})
		if ok != stated[kind] {
			t.Errorf("blocker %q: badge present = %v, want %v", kind, ok, stated[kind])
		}
		if !ok {
			continue
		}
		if badge.Glyph == "" || badge.Label == "" {
			t.Errorf("blocker %q has no glyph or label: %+v", kind, badge)
		}
		if other, exists := seen[badge.Text()]; exists {
			t.Errorf("blockers %q and %q both render as %q", kind, other, badge.Text())
		}
		seen[badge.Text()] = kind
	}
}

func TestBlockerBadgeNamesTheWaitTarget(t *testing.T) {
	task, _ := blockerBadge(data.GateBlocker{
		Kind:       data.GateBlockDependency,
		Dependency: &data.DependencyBlock{Target: "T042"},
	})
	if !strings.Contains(task.Label, "T042") {
		t.Errorf("dependency badge = %q, want the Task it waits on named", task.Label)
	}

	objective, _ := blockerBadge(data.GateBlocker{
		Kind:                data.GateBlockObjectiveDependency,
		ObjectiveDependency: &data.ObjectiveDependencyBlock{Target: "O007"},
	})
	if !strings.Contains(objective.Label, "O007") {
		t.Errorf("objective dependency badge = %q, want the Objective it waits on named", objective.Label)
	}
	if objective.Text() == task.Text() {
		t.Error("a Task wait and an Objective wait render identically")
	}

	// A blocker whose typed block is absent still says the Task is waiting,
	// rather than trailing an empty target.
	bare, _ := blockerBadge(data.GateBlocker{Kind: data.GateBlockDependency})
	if strings.HasSuffix(bare.Label, " ") {
		t.Errorf("bare dependency badge = %q, want no trailing separator", bare.Label)
	}
}

// TestTaskReviewOutcomeBadgeRendersEveryPathAtExactText proves every path
// O012's Done When names renders exactly the fixed text, whether the Task is
// still open or already done — completion never changes the wording, since
// the Done column carries that fact instead (see card.go's badges()).
func TestTaskReviewOutcomeBadgeRendersEveryPathAtExactText(t *testing.T) {
	for _, test := range taskReviewOutcomeCases {
		t.Run(test.name, func(t *testing.T) {
			got := taskReviewOutcomeBadge(test.clearance, test.waived, test.exception).Text()
			if got != test.want {
				t.Errorf("taskReviewOutcomeBadge(%q, waived=%v, exception=%v) = %q, want %q",
					test.clearance, test.waived, test.exception, got, test.want)
			}
		})
	}
}

// TestTaskReviewOutcomeBadgeStaysDistinctAcrossThePrecedenceLadder proves the
// eight cases above never collapse two different facts into the same text,
// except the two collisions the design deliberately intends: "stale" and
// "unknown" both fold into "REVIEW" (see taskReviewOutcomeBadge), and "owner
// accepted outranks waived" is the same fact as "owner accepted" under a
// different clearance.
func TestTaskReviewOutcomeBadgeStaysDistinctAcrossThePrecedenceLadder(t *testing.T) {
	deliberateDuplicate := map[string]bool{
		"unknown":                        true,
		"owner accepted outranks waived": true,
	}
	seen := map[string]string{}
	for _, test := range taskReviewOutcomeCases {
		got := taskReviewOutcomeBadge(test.clearance, test.waived, test.exception).Text()
		if deliberateDuplicate[test.name] {
			continue
		}
		if other, ok := seen[got]; ok {
			t.Errorf("cases %q and %q both render as %q", test.name, other, got)
		}
		seen[got] = test.name
	}
}

// TestBadgesStayDistinctWithColorDisabled is the monochrome guarantee: with
// every accent stripped, no two badges in the whole vocabulary collapse into
// the same text.
func TestBadgesStayDistinctWithColorDisabled(t *testing.T) {
	forceColorProfile(t, termenv.Ascii)

	seen := map[string]string{}
	for _, badge := range allBadges() {
		rendered := badge.Render()
		if xansi.Strip(rendered) != rendered {
			t.Errorf("badge %q still carries escapes with color disabled", badge.Text())
		}
		if other, ok := seen[rendered]; ok {
			t.Errorf("badges %q and %q are indistinguishable without color", badge.Text(), other)
		}
		seen[rendered] = badge.Text()
	}
}

// TestBadgesCarryGlyphAndLabelUnderColor proves color is reinforcement: with a
// full-color profile, the text under the escapes is the same text the
// monochrome terminal shows.
func TestBadgesCarryGlyphAndLabelUnderColor(t *testing.T) {
	forceColorProfile(t, termenv.TrueColor)

	for _, badge := range allBadges() {
		if stripped := xansi.Strip(badge.Render()); stripped != badge.Text() {
			t.Errorf("badge renders as %q under color, want %q", stripped, badge.Text())
		}
	}
}

// issueTypesForBadges is every data.IssueType value the decoder admits, plus
// one word outside that vocabulary — issueTypeBadge switches on the type's
// string value rather than the data package's named constants (see its own
// doc comment), so this proves the fallback still covers a value the decoder
// itself would reject.
var issueTypesForBadges = []data.IssueType{"defect", "drift", "guardrail", "verification", "other", "not-a-real-type"}

func TestIssueTypeBadgeCoversEveryTypeDistinctly(t *testing.T) {
	seen := map[string]data.IssueType{}
	for _, issueType := range issueTypesForBadges {
		badge := issueTypeBadge(issueType)
		if badge.Glyph == "" || badge.Label == "" {
			t.Errorf("type %q has no glyph or label: %+v", issueType, badge)
		}
		if other, ok := seen[badge.Text()]; ok {
			t.Errorf("types %q and %q both render as %q", issueType, other, badge.Text())
		}
		seen[badge.Text()] = issueType
	}
}

// issueSeveritiesForBadges is the five-word scale this board reads out of the
// opaque, policy-owned severity field, plus a word outside that scale and the
// empty string severity itself allows (internal/data/issue_v2.go).
var issueSeveritiesForBadges = []string{"blocker", "high", "medium", "low", "cosmetic", "urgent"}

func TestIssueSeverityBadgeCoversTheScaleDistinctlyAndDegradesForUnknown(t *testing.T) {
	seen := map[string]string{}
	for _, severity := range issueSeveritiesForBadges {
		badge, ok := issueSeverityBadge(severity)
		if !ok {
			t.Errorf("severity %q reported ok=false, want a badge", severity)
		}
		if badge.Label != strings.ToUpper(severity) {
			t.Errorf("severity %q badge label = %q, want %q", severity, badge.Label, strings.ToUpper(severity))
		}
		if other, existing := seen[badge.Text()]; existing {
			t.Errorf("severities %q and %q both render as %q", severity, other, badge.Text())
		}
		seen[badge.Text()] = severity
	}
}

func TestIssueSeverityBadgeIsAbsentWhenNoneIsRecorded(t *testing.T) {
	if _, ok := issueSeverityBadge(""); ok {
		t.Error("empty severity produced a badge; a card should show none instead")
	}
	if _, ok := issueSeverityBadge("   "); ok {
		t.Error("whitespace-only severity produced a badge; a card should show none instead")
	}
}

// TestIssueSeverityRankOrdersBlockerFirstAndUnrecordedLast proves the sort
// this board applies within an Issues column: blocker first, cosmetic last
// of the five-word scale, a real word outside that scale ranked after it
// (still shown, never dropped), and no recorded severity at all ranked last
// of everything.
func TestIssueSeverityRankOrdersBlockerFirstAndUnrecordedLast(t *testing.T) {
	ranked := []string{"blocker", "high", "medium", "low", "cosmetic", "urgent", ""}
	for i := 1; i < len(ranked); i++ {
		prev, next := issueSeverityRank(ranked[i-1]), issueSeverityRank(ranked[i])
		if prev >= next {
			t.Errorf("issueSeverityRank(%q)=%d is not before issueSeverityRank(%q)=%d", ranked[i-1], prev, ranked[i], next)
		}
	}
	// Case- and whitespace-insensitive: the same word ranks the same
	// regardless of how the record spelled it.
	if issueSeverityRank("HIGH") != issueSeverityRank(" high ") {
		t.Error("issueSeverityRank is not case/whitespace insensitive")
	}
}

// allBadges is every badge the vocabulary can produce, for the distinctness and
// color assertions above.
func allBadges() []Badge {
	var badges []Badge
	for _, stage := range []data.ProgressStage{data.StageBuild, data.StageTest, data.StageAudit} {
		if badge, ok := stageBadge(data.ColumnInProgress, stage); ok {
			badges = append(badges, badge)
		}
	}
	// objectiveCheckBadge is deliberately excluded here: it is a Task card's
	// sibling vocabulary for a different row context (the sidebar), and its
	// "[✓] Check" for a current Objective is meant to read as the same fact
	// taskCheckBadge's "[✓] Check" does for a current Task — reused
	// presentation for reused meaning, not a collision. Its own two-state
	// distinctness is covered by TestObjectiveCheckBadgeIsTwoNotchOnly.
	for _, state := range clearanceStates {
		badges = append(badges, taskCheckBadge(state, false))
	}
	badges = append(badges, taskCheckBadge(data.ClearanceMissing, true))
	for _, kind := range blockerKinds {
		if badge, ok := blockerBadge(data.GateBlocker{Kind: kind}); ok {
			badges = append(badges, badge)
		}
	}
	// taskReviewOutcomeBadge's "unknown" and "owner accepted outranks waived"
	// cases are deliberately excluded: each renders identical text to another
	// case by design — "unknown" folds into "stale"'s "REVIEW", and "owner
	// accepted outranks waived" repeats "owner accepted"'s text under a
	// different clearance — which
	// TestTaskReviewOutcomeBadgeStaysDistinctAcrossThePrecedenceLadder already
	// covers on its own.
	skipDuplicate := map[string]bool{"unknown": true, "owner accepted outranks waived": true}
	for _, test := range taskReviewOutcomeCases {
		if skipDuplicate[test.name] {
			continue
		}
		badges = append(badges, taskReviewOutcomeBadge(test.clearance, test.waived, test.exception))
	}
	return badges
}

// forceColorProfile pins Lip Gloss's color profile for one test and restores
// whatever the suite was using afterwards.
func forceColorProfile(t *testing.T, profile termenv.Profile) {
	t.Helper()
	previous := lipgloss.ColorProfile()
	lipgloss.SetColorProfile(profile)
	t.Cleanup(func() { lipgloss.SetColorProfile(previous) })
}
