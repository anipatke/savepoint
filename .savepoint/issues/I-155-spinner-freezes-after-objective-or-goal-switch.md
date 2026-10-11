---
id: I-155
title: Spinner freezes after switching the board to another Objective or Goal and back
type: defect
status: resolved
source:
  kind: check
  check: C-980
  actor: {role: checker, session: check-o048-20261011}
  at: '2026-10-11T05:45:00Z'
tasks: [T-137]
checks: [C-980, C-981]
guardrail_ids: [ARCH-02, TEST-02]
severity: medium
resolution:
  disposition: verified
  check: C-981
  actor: {role: checker, session: check-o048-recheck-20261011}
  at: '2026-10-11T05:55:00Z'
  reason: 'Commit 6145746 starts the timer from Update after any non-tick message that returns no command; reload path still starts it. Independent probes: esc back to whole Goal, sidebar up/down walk, and Goal picker switch R-001->R-002->R-001->R-002 all restart or stop the timer correctly; no second timer from an ordinary key; idle board keys start none. New regression test TestSpinnerRestartsWhenAnObjectiveSwitchShowsALiveTileAgain passes.'
history:
  - at: '2026-10-11T05:45:00Z'
    actor: {role: checker, session: check-o048-20261011}
    kind: observed
    check: C-980
    note: Found by the O-048 Full Objective Check.
  - at: '2026-10-11T05:55:00Z'
    actor: {role: checker, session: check-o048-recheck-20261011}
    kind: rechecked
    check: C-981
    note: 'Commit 6145746 starts the timer from Update after any non-tick message that returns no command; reload path still starts it. Independent probes: esc back to whole Goal, sidebar up/down walk, and Goal picker switch R-001->R-002->R-001->R-002 all restart or stop the timer correctly; no second timer from an ordinary key; idle board keys start none. New regression test TestSpinnerRestartsWhenAnObjectiveSwitchShowsALiveTileAgain passes.'
---

# I-155: Spinner freezes after switching the board to another Objective or Goal and back

## Summary

The spinner timer is started only from `applyLoad` (`internal/board/v2/update.go:724`, `m.spinnerCmd()`). Changing the shown cards without a reload, through `selectObjective` (`update.go:594-602`, reached by sidebar up/down and esc) or `applyReleaseSelection` (`update.go:608-624`, the Goal picker), never calls `spinnerCmd`. Once the timer has stopped because the narrowed view had no build or test tile, returning to a view with one shows "Building…"/"Testing…" with a frozen frame until some file change triggers a reload.

Violated requirements: O-048 SC1 ("shows a moving spinner"), SC3 (ticks while at least one shown tile is at build or test); T-137 DW2 and DW5 (tick start/stop tests cover reloads only).

## Evidence

- Smallest reproduction (temporary in-package probe, removed after the run): `writeBadgeProject` plus O-002 with one planned Task; `openBoard` -> `spinnerRunning=true`; `selectObjective("O-002")` then one `spinnerTickMsg` -> timer ends (correct); sidebar focused, `esc` -> `SelectedObjective=""`, `hasLiveTile()=true`, `spinnerRunning=false`, returned cmd nil, frame stays 0.
- Expected: the spinner moves again as soon as a live tile is shown. Actual: it stays still.
- Owner path: simply moving the sidebar cursor across Objectives (each move calls `selectObjective`) passes through Objectives without live tiles.
- Missing test: `spinner_test.go` covers start/stop only across reloads and direct ticks, not across a scope change in Update.

## Proof Needed

Any Update path that changes `m.Cards` starts the timer when a live tile becomes shown (without starting a second one), with a regression test for Objective selection and Goal selection, plus the existing gates.
