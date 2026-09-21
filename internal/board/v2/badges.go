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
	glyphUnknown   = "?"
	glyphAttention = "⚠"
	glyphWaiting   = "→"
	glyphOwner     = "!"

	// Check-badge glyphs: an explicit checkbox rather than the narrower
	// clearance glyphs above, since this badge is read on its own as the
	// friendly "has this been checked" answer — grey empty, green ticked,
	// amber flagged — not as one entry in the fuller clearance vocabulary.
	glyphCheckPending = "[ ]"
	glyphCheckClear   = "[✓]"
	glyphCheckFlagged = "[!]"
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
		return Badge{Glyph: glyphBuild, Label: stageLabel(stage), Style: styles.BadgeAttention}, true
	case data.StageTest:
		return Badge{Glyph: glyphTest, Label: stageLabel(stage), Style: styles.BadgeClear}, true
	case data.StageAudit:
		return Badge{Glyph: glyphAudit, Label: stageLabel(stage), Style: styles.BadgeWaiting}, true
	default:
		return Badge{}, false
	}
}

// stageLabel is the one place a Task's stage becomes owner-facing text.
// data.StageAudit's stored value stays "audit" — every existing Task record,
// including ones already on disk, keeps parsing unchanged — but nothing the
// board displays says that word: this project's stages read as build, test,
// check, matching the rest of V2's Check vocabulary. A stage outside the
// three canonical values reports itself rather than a guess.
func stageLabel(stage data.ProgressStage) string {
	switch stage {
	case data.StageBuild:
		return "BUILD"
	case data.StageTest:
		return "TEST"
	case data.StageAudit:
		return "CHECK"
	default:
		return string(stage)
	}
}

// taskCheckBadge is the friendly three-notch Check badge a Task card shows:
// grey "[ ] Check" before any Check has ever been requested, green "[✓]
// Check" once an independent Check recorded it current, and amber "[!]
// Check" when the owner waived the local Check instead (never a fourth,
// unmarked state — a waiver is not technical CLEAR, and this badge does not
// pretend otherwise). A Check that was recorded but came back NEEDS WORK,
// or whose freshness is stale or unknown, is a real, distinct problem — it
// keeps the amber "[!]" accent so it still reads as needing attention, but
// its own label names which of the three it is rather than collapsing them
// into the same word as a waiver.
func taskCheckBadge(clearance data.ClearanceState, waived bool) Badge {
	if waived {
		return Badge{Glyph: glyphCheckFlagged, Label: "Check (waived)", Style: styles.BadgeAttention}
	}
	switch clearance {
	case data.ClearanceCurrent:
		return Badge{Glyph: glyphCheckClear, Label: "Check", Style: styles.BadgeClear}
	case data.ClearanceNeedsWork:
		return Badge{Glyph: glyphCheckFlagged, Label: "Check (needs work)", Style: styles.BadgeAttention}
	case data.ClearanceStale:
		return Badge{Glyph: glyphCheckFlagged, Label: "Check (stale)", Style: styles.BadgeAttention}
	case data.ClearanceUnknown:
		return Badge{Glyph: glyphCheckFlagged, Label: "Check (unverified)", Style: styles.BadgeAttention}
	default:
		// ClearanceMissing: no Check has been requested yet, and no waiver is
		// recorded either — still mid-build, nothing to flag.
		return Badge{Glyph: glyphCheckPending, Label: "Check", Style: styles.BadgeNeutral}
	}
}

// objectiveCheckBadge is the Objective sidebar row's one completion-state
// badge — deliberately the only one. An earlier version of this row carried
// a second badge ("INTEGRATED" / "NEEDS INTEGRATION") that restated this same
// clearance fact in different words once every owned Task was done; that was
// confusing and was removed rather than reconciled; see I012. Do not
// reintroduce a second badge for Task-completeness — clearance is the one
// fact this row states about whether an Objective is ready to close, and it
// states it once.
//
// The row deliberately shows only three notches, not the five-state
// vocabulary taskCheckBadge carries. Savepoint's audience does not need
// needs_work, stale, and unverified told apart at a glance — all three are
// "this Objective was checked and is not clear," so they read identically
// here. Anyone who needs the distinction reads it in the detail overlay's
// CLEARANCE section (resume.ClearancePhrase), which still states the exact
// recorded reason. Do not re-expand this badge to five states; that was
// tried and explicitly walked back for being over the reader's head.
//
// byException folds an owner-recorded exception into the same green tick a
// current Check gets, rather than into its own "BY EXCEPTION" notch. The
// reasoning an owner accepted, and the fact that a re-check was skipped for
// it, are durable and already recorded — Evidence.Exception, read by the
// detail overlay's EXCEPTION section — but they are not information the
// compact row needs to flag. An owner who chose to accept the findings and
// move on gets the same clean "checked" signal a clear Check does; anyone
// who wants to know why opens the detail. Do not add a fourth notch for this
// without a fresh product decision — it was deliberately cut once already.
func objectiveCheckBadge(clearance data.ClearanceState, byException bool) Badge {
	if byException || clearance == data.ClearanceCurrent {
		return Badge{Glyph: glyphCheckClear, Label: "Check", Style: styles.BadgeClear}
	}
	if clearance == data.ClearanceMissing {
		return Badge{Glyph: glyphCheckPending, Label: "Check", Style: styles.BadgeNeutral}
	}
	// needs_work, stale, and unverified: checked, not clear, collapsed to one
	// wording on purpose (see doc comment above).
	return Badge{Glyph: glyphCheckFlagged, Label: "Check (needs work)", Style: styles.BadgeAttention}
}

// completionBadge names how a done Task reached done. A recorded exception or
// waiver reads differently from an ordinary close, and a close whose
// clearance is no longer current reads as needing attention rather than as
// finished — every distinction the release design requires to be visible.
func completionBadge(clearance data.ClearanceState, byException, byWaiver bool) Badge {
	switch {
	case byException:
		return Badge{Glyph: glyphOwner, Label: "BY EXCEPTION", Style: styles.BadgeAttention}
	case byWaiver:
		return Badge{Glyph: glyphCheckFlagged, Label: "BY WAIVER", Style: styles.BadgeAttention}
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

// clearanceIsMissing reports whether clearance means no Check has ever been
// recorded at all — the one state a Task-check waiver may stand in for, at
// Space's completion write in io.go.
func clearanceIsMissing(clearance data.Clearance) bool {
	return clearance.State == data.ClearanceMissing
}

func hasOwnerAcceptanceBlock(decision data.GateDecision) bool {
	for _, blocker := range decision.Blockers {
		if blocker.Kind == data.GateBlockOwnerAcceptance {
			return true
		}
	}
	return false
}
