package v2

import (
	"strings"

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

	// Issue type glyphs: one per type.IssueType value, agreed with the owner
	// before implementing rather than guessed. The "defect" type reads as a
	// cross since the board's no-emoji rule (above) rules out a literal bug.
	glyphIssueCross        = "✗"
	glyphIssueDrift        = "≈"
	glyphIssueGuardrail    = "‖"
	glyphIssueVerification = "◎"
	glyphIssueOther        = "…"

	// Issue severity glyphs taper from solid to faint — blocker heaviest,
	// cosmetic lightest — so the shape alone carries the five-step scale
	// even with color stripped; color then adds urgency on top of that:
	// orange for the top two, plain text for medium, dim for the bottom two.
	glyphSeverityBlocker  = "●"
	glyphSeverityHigh     = "▲"
	glyphSeverityMedium   = "▪"
	glyphSeverityLow      = "▫"
	glyphSeverityCosmetic = "·"
	// glyphSeverityOther marks a recorded severity outside the five-word
	// vocabulary above: severity is policy-owned and opaque to this board
	// (internal/data/issue_v2.go), so an unrecognized word is still shown,
	// at neutral weight, rather than guessed at or dropped.
	glyphSeverityOther = "◦"
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

// taskReviewOutcomeBadge is the one review-outcome badge a non-planned Task
// card shows. It replaces the separate completion-plus-Check composition an
// earlier version carried (completionBadge, exceptionBadge — see O012):
// completion is now the Done column's own fact, so this badge states review
// standing alone, at fixed precedence — an owner's recorded risk acceptance
// (exception) first, an owner's recorded Check waiver second, and the
// resolved clearance state otherwise — so one card never prints two
// competing review outcomes. Stale and unknown clearance collapse to the
// same "REVIEW" wording a checker-authority failure gets (see blockerBadge):
// both mean the same thing to a reader — this needs an independent look —
// and the exact reason remains available off the card, in the detail
// overlay's CLEARANCE section. Waiver and owner-accepted risk keep the same
// green "clear" accent a passing Check gets, per the palette rule that green
// means an outcome accepted for the Task; only the label tells them apart,
// deliberately, since Savepoint's data layer never treats either as
// technical clearance (internal/data/gate_v2.go, evidence_v2.go).
func taskReviewOutcomeBadge(clearance data.ClearanceState, waived, byException bool) Badge {
	switch {
	case byException:
		return Badge{Glyph: glyphCheckClear, Label: "OWNER ACCEPTED", Style: styles.BadgeClear}
	case waived:
		return Badge{Glyph: glyphCheckClear, Label: "WAIVED", Style: styles.BadgeClear}
	}
	switch clearance {
	case data.ClearanceCurrent:
		return Badge{Glyph: glyphCheckClear, Label: "CHECK", Style: styles.BadgeClear}
	case data.ClearanceNeedsWork:
		return Badge{Glyph: glyphCheckFlagged, Label: "NEEDS WORK", Style: styles.BadgeAttention}
	case data.ClearanceStale, data.ClearanceUnknown:
		return Badge{Glyph: glyphCheckFlagged, Label: "REVIEW", Style: styles.BadgeAttention}
	default:
		// ClearanceMissing: no Check has been requested yet, and no waiver is
		// recorded either — still mid-build, nothing to flag.
		return Badge{Glyph: glyphCheckPending, Label: "CHECK", Style: styles.BadgeNeutral}
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
// It reports ok=false for the clearance blockers — missing, needs_work,
// stale, unknown, and checker_authority — because taskReviewOutcomeBadge
// already states exactly those, from the same resolved clearance value
// (stale, unknown, and a checker-authority failure all read as the card's
// one "REVIEW" outcome); showing both would print the same fact twice in
// two wordings. Every other blocker kind is a fact no other badge carries.
func blockerBadge(blocker data.GateBlocker) (Badge, bool) {
	switch blocker.Kind {
	case data.GateBlockReplan:
		return Badge{Glyph: glyphAttention, Label: "REPLAN", Style: styles.BadgeAttention}, true
	case data.GateBlockDependency:
		return Badge{Glyph: glyphWaiting, Label: waitLabel("WAITS", dependencyTarget(blocker)), Style: styles.BadgeWaiting}, true
	case data.GateBlockObjectiveDependency:
		return Badge{Glyph: glyphWaiting, Label: waitLabel("OBJECTIVE WAITS", objectiveDependencyTarget(blocker)), Style: styles.BadgeWaiting}, true
	case data.GateBlockOwnerAcceptance:
		return Badge{Glyph: glyphOwner, Label: "AWAITS OWNER", Style: styles.BadgeAttention}, true
	default:
		// The clearance kinds, GateBlockCheckerAuthority (folded into the
		// review outcome's "REVIEW" wording), and GateBlockInvalidState,
		// which names a recorded status/stage combination the strict V2
		// decoder cannot admit and a card therefore cannot be showing.
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

// reviewOutcomeIsActionable reports whether a clearance state is itself the
// fact a still-open Task's card needs to state. NEEDS WORK and the
// stale/unknown REVIEW fold are real attention items a reader cannot infer
// from the stage badge alone; ClearanceCurrent and ClearanceMissing are
// completion-outcome vocabulary that belongs to the Done column instead (see
// TaskCard.showsReviewOutcome in card.go).
func reviewOutcomeIsActionable(state data.ClearanceState) bool {
	switch state {
	case data.ClearanceNeedsWork, data.ClearanceStale, data.ClearanceUnknown:
		return true
	default:
		return false
	}
}

func hasOwnerAcceptanceBlock(decision data.GateDecision) bool {
	for _, blocker := range decision.Blockers {
		if blocker.Kind == data.GateBlockOwnerAcceptance {
			return true
		}
	}
	return false
}

// issueTypeBadge names an Issue's descriptive type. It stays neutral weight
// regardless of which type: type is descriptive, never itself a severity or
// a blocker (see agent-skills/references/issue-capture.md), so it never
// borrows the attention accent severity badges use.
//
// It switches on the type's own string value rather than the data package's
// named constants: the V1/V2 boundary this package holds
// (TestPackageCarriesNoReleaseOrEpicSurface) refuses one of those names
// outright, left over from a retired V1 surface with no V2 meaning, and this
// V2 Issue type is an unrelated field that happens to share the English word.
func issueTypeBadge(issueType data.IssueType) Badge {
	label := strings.ToUpper(string(issueType))
	switch string(issueType) {
	case "defect":
		return Badge{Glyph: glyphIssueCross, Label: label, Style: styles.CardMeta}
	case "drift":
		return Badge{Glyph: glyphIssueDrift, Label: label, Style: styles.CardMeta}
	case "guardrail":
		return Badge{Glyph: glyphIssueGuardrail, Label: label, Style: styles.CardMeta}
	case "verification":
		return Badge{Glyph: glyphIssueVerification, Label: label, Style: styles.CardMeta}
	default: // "other", and any value the decoder would otherwise reject
		return Badge{Glyph: glyphIssueOther, Label: label, Style: styles.CardMeta}
	}
}

// issueSeverityBadge names a recorded severity, or reports ok=false for an
// Issue that declared none — severity is optional, and a card says nothing
// about it rather than showing an empty badge. The five-word scale is this
// board's own reading of an opaque, policy-owned field: a word outside it is
// still shown, uppercased, at the same neutral weight unrecognized-but-real
// data gets everywhere else in this package, never dropped or guessed at.
func issueSeverityBadge(severity string) (Badge, bool) {
	trimmed := strings.TrimSpace(severity)
	if trimmed == "" {
		return Badge{}, false
	}
	label := strings.ToUpper(trimmed)
	switch strings.ToLower(trimmed) {
	case "blocker":
		return Badge{Glyph: glyphSeverityBlocker, Label: label, Style: styles.BadgeAttention}, true
	case "high":
		return Badge{Glyph: glyphSeverityHigh, Label: label, Style: styles.BadgeAttention}, true
	case "medium":
		return Badge{Glyph: glyphSeverityMedium, Label: label, Style: styles.TaskItem}, true
	case "low":
		return Badge{Glyph: glyphSeverityLow, Label: label, Style: styles.BadgeNeutral}, true
	case "cosmetic":
		return Badge{Glyph: glyphSeverityCosmetic, Label: label, Style: styles.BadgeNeutral}, true
	default:
		return Badge{Glyph: glyphSeverityOther, Label: label, Style: styles.BadgeNeutral}, true
	}
}

// issueSeverityRank is this board's own read of the same five-word scale,
// for sorting rather than display: blocker first, cosmetic last, a
// recognized-but-listed word after cosmetic, and no recorded severity at all
// last of everything. It never touches the stored field — severity stays
// exactly as recorded (internal/data/issue_v2.go) — this is a display-only
// ordering, not a validated vocabulary.
func issueSeverityRank(severity string) int {
	switch strings.ToLower(strings.TrimSpace(severity)) {
	case "blocker":
		return 0
	case "high":
		return 1
	case "medium":
		return 2
	case "low":
		return 3
	case "cosmetic":
		return 4
	case "":
		return 6
	default:
		return 5
	}
}
