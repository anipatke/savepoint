package v2

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/styles"
)

// This file is the board's only translation from a typed `data` value into
// something on screen. Every other surface in this package — cards, the
// sidebar, detail, and the plain renderer — reuses these functions rather than
// writing a second copy of the vocabulary (STYLE-07, STYLE-09).
//
// Nothing here decides anything. `ClearanceStale` arrives from
// ResolveClearance and leaves as a glyph, a label, and an accent; the board
// does not work out that clearance is stale, and it introduces no lifecycle
// word of its own (DATA-02).

// Badge is one state marker: a glyph, a short label, and the accent that
// colors both. Glyph and label together carry the state, so a badge stays
// readable with color stripped — which the monochrome and non-TTY paths
// depend on.
type Badge struct {
	Glyph string
	Label string
	Style lipgloss.Style
}

// Text is the badge with no styling: the glyph and label a reader sees when
// color is absent.
func (b Badge) Text() string {
	return b.Glyph + " " + b.Label
}

// Render is the badge as it appears on a card.
func (b Badge) Render() string {
	return b.Style.Render(b.Text())
}

// Glyphs are deliberately narrow, single-cell, and drawn from the set the
// board already uses, so a badge never changes a card's geometry and never
// depends on emoji width.
const (
	glyphBuild     = "▣"
	glyphTest      = "◇"
	glyphAudit     = "◆"
	glyphClear     = "✓"
	glyphNeedsWork = "✗"
	glyphStale     = "◐"
	glyphUnknown   = "?"
	glyphNone      = "○"
	glyphAttention = "⚠"
	glyphWaiting   = "→"
	glyphOwner     = "!"
)

// stageBadge names the implementation stage of a Task under way. It reports
// ok=false for any Task that is not in progress: a planned Task has no stage
// and a done Task's stage is history, so neither carries one.
func stageBadge(status data.ColumnType, stage data.ProgressStage) (Badge, bool) {
	if status != data.ColumnInProgress {
		return Badge{}, false
	}
	switch stage {
	case data.StageBuild:
		return Badge{Glyph: glyphBuild, Label: "BUILD", Style: styles.BadgeAttention}, true
	case data.StageTest:
		return Badge{Glyph: glyphTest, Label: "TEST", Style: styles.BadgeClear}, true
	case data.StageAudit:
		return Badge{Glyph: glyphAudit, Label: "AUDIT", Style: styles.BadgeWaiting}, true
	default:
		return Badge{}, false
	}
}

// clearanceBadge names a resolved clearance state. All five states render with
// their own glyph and label, so "stale" is never mistaken for "clear" on a
// terminal without color.
func clearanceBadge(state data.ClearanceState) Badge {
	switch state {
	case data.ClearanceCurrent:
		return Badge{Glyph: glyphClear, Label: "CLEAR", Style: styles.BadgeClear}
	case data.ClearanceNeedsWork:
		return Badge{Glyph: glyphNeedsWork, Label: "NEEDS WORK", Style: styles.BadgeAttention}
	case data.ClearanceStale:
		return Badge{Glyph: glyphStale, Label: "STALE", Style: styles.BadgeAttention}
	case data.ClearanceUnknown:
		return Badge{Glyph: glyphUnknown, Label: "UNVERIFIED", Style: styles.BadgeAttention}
	default:
		// ClearanceMissing, and any state ResolveClearance could add later:
		// reported as nothing recorded rather than as a met requirement.
		return Badge{Glyph: glyphNone, Label: "NO CHECK", Style: styles.BadgeNeutral}
	}
}

// completionBadge names how a done Task reached done. A recorded exception
// reads differently from an ordinary close, and a close whose clearance is no
// longer current reads as needing attention rather than as finished — both
// distinctions the release design requires to be visible.
func completionBadge(clearance data.ClearanceState, byException bool) Badge {
	switch {
	case byException:
		return Badge{Glyph: glyphOwner, Label: "BY EXCEPTION", Style: styles.BadgeAttention}
	case clearance == data.ClearanceCurrent:
		return Badge{Glyph: glyphClear, Label: "DONE", Style: styles.BadgeClear}
	default:
		return Badge{Glyph: glyphAttention, Label: "DONE", Style: styles.BadgeAttention}
	}
}

// exceptionBadge names a completion a recorded exception allows for a Task
// still open. It is the same wording completionBadge gives a closed one, so a
// reader learns one phrase for one fact.
func exceptionBadge() Badge {
	return Badge{Glyph: glyphOwner, Label: "BY EXCEPTION", Style: styles.BadgeAttention}
}

// objectiveIntegrationBadge names where an Objective stands once every Task it
// owns is done. An Objective reaches done only when its Tasks meet completion
// rules *and* its own integration Check is current, so finished Tasks alone are
// reported as work still to do rather than as a finished Objective — otherwise
// that state is invisible, because every column looks complete.
//
// It reports ok=false for an Objective whose Tasks are not all done: there is
// no integration question to answer yet, and the clearance badge already states
// what has been recorded.
func objectiveIntegrationBadge(tasksComplete bool, clearance data.ClearanceState) (Badge, bool) {
	switch {
	case !tasksComplete:
		return Badge{}, false
	case clearance == data.ClearanceCurrent:
		return Badge{Glyph: glyphClear, Label: "INTEGRATED", Style: styles.BadgeClear}, true
	default:
		return Badge{Glyph: glyphAttention, Label: "NEEDS INTEGRATION", Style: styles.BadgeAttention}, true
	}
}

// objectiveWaitBadge names one unsatisfied Objective dependency, read from the
// typed block ResolveObjectiveDependency returned. It reuses the wait wording a
// Task card already carries, so one phrase means one fact across the board.
func objectiveWaitBadge(block data.ObjectiveDependencyBlock) Badge {
	return Badge{Glyph: glyphWaiting, Label: waitLabel("WAITS", block.Target), Style: styles.BadgeWaiting}
}

// blockerBadge names one unmet requirement from a gate decision.
//
// It reports ok=false for the clearance blockers — missing, needs_work, stale,
// unknown — because clearanceBadge already states exactly those, from the same
// resolved value; showing both would print the same fact twice in two
// wordings. Every other blocker kind is a fact no other badge carries.
func blockerBadge(blocker data.GateBlocker) (Badge, bool) {
	switch blocker.Kind {
	case data.GateBlockReplan:
		return Badge{Glyph: glyphAttention, Label: "REPLAN", Style: styles.BadgeAttention}, true
	case data.GateBlockDependency:
		return Badge{Glyph: glyphWaiting, Label: waitLabel("WAITS", dependencyTarget(blocker)), Style: styles.BadgeWaiting}, true
	case data.GateBlockObjectiveDependency:
		return Badge{Glyph: glyphWaiting, Label: waitLabel("OBJECTIVE WAITS", objectiveDependencyTarget(blocker)), Style: styles.BadgeWaiting}, true
	case data.GateBlockOwnerAcceptance:
		return Badge{Glyph: glyphOwner, Label: "OWNER", Style: styles.BadgeAttention}, true
	case data.GateBlockCheckerAuthority:
		return Badge{Glyph: glyphUnknown, Label: "CHECKER", Style: styles.BadgeAttention}, true
	default:
		// The clearance kinds, and GateBlockInvalidState, which names a
		// recorded status/stage combination the strict V2 decoder cannot
		// admit and a card therefore cannot be showing.
		return Badge{}, false
	}
}

// waitLabel names what a wait is on. A blocker whose typed block carries no
// target still says that the Task is waiting, rather than trailing a space
// where an ID should be.
func waitLabel(prefix, target string) string {
	if target == "" {
		return prefix
	}
	return prefix + " " + target
}

// dependencyTarget and objectiveDependencyTarget read the target off the typed
// block the resolver attached to the blocker. Neither looks the dependency up
// or re-evaluates it; an unnamed target degrades to the label alone rather
// than to a guess.
func dependencyTarget(blocker data.GateBlocker) string {
	if blocker.Dependency == nil {
		return ""
	}
	return blocker.Dependency.Target
}

func objectiveDependencyTarget(blocker data.GateBlocker) string {
	if blocker.ObjectiveDependency == nil {
		return ""
	}
	return blocker.ObjectiveDependency.Target
}

func clearanceIsCurrent(clearance data.Clearance) bool {
	return clearance.State == data.ClearanceCurrent
}

func hasOwnerAcceptanceBlock(decision data.GateDecision) bool {
	for _, blocker := range decision.Blockers {
		if blocker.Kind == data.GateBlockOwnerAcceptance {
			return true
		}
	}
	return false
}
