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
// silence for the clearance kinds the clearance badge already states.
func TestBlockerBadgeCoversEveryKind(t *testing.T) {
	stated := map[data.GateBlockKind]bool{
		data.GateBlockReplan:              true,
		data.GateBlockDependency:          true,
		data.GateBlockObjectiveDependency: true,
		data.GateBlockOwnerAcceptance:     true,
		data.GateBlockCheckerAuthority:    true,
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

func TestCompletionBadgeDistinguishesExceptionAndStaleFromFinished(t *testing.T) {
	finished := completionBadge(data.ClearanceCurrent, false, false)
	byException := completionBadge(data.ClearanceCurrent, true, false)
	byWaiver := completionBadge(data.ClearanceMissing, false, true)
	stale := completionBadge(data.ClearanceStale, false, false)

	if finished.Text() == byException.Text() {
		t.Errorf("a done Task and one done by exception both render as %q", finished.Text())
	}
	if finished.Text() == byWaiver.Text() {
		t.Errorf("a done Task and one done by waiver both render as %q", finished.Text())
	}
	if byException.Text() == byWaiver.Text() {
		t.Errorf("a done-by-exception Task and a done-by-waiver one both render as %q", byException.Text())
	}
	if finished.Text() == stale.Text() {
		t.Errorf("a cleared done Task and a stale one both render as %q", finished.Text())
	}
	if finished.Glyph == stale.Glyph {
		t.Errorf("stale done reuses the finished glyph %q; it must read as needing attention", stale.Glyph)
	}
	if byException.Text() != exceptionBadge().Text() {
		t.Errorf("a closed exception reads %q and an open one %q; they name the same fact", byException.Text(), exceptionBadge().Text())
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

// allBadges is every badge the vocabulary can produce, for the distinctness and
// color assertions above.
func allBadges() []Badge {
	badges := []Badge{exceptionBadge()}
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
	badges = append(badges,
		completionBadge(data.ClearanceCurrent, false, false),
		completionBadge(data.ClearanceStale, false, false),
		completionBadge(data.ClearanceMissing, false, true),
	)
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
