---
id: T-137
title: Animate tiles that are building or testing
objective: O-048
status: done
depends_on: []
owner_validation:
    required: true
    accepted_check: ""
planned_by: {role: planner, session: planning-o048-animate-tiles-20261011}
check_waiver:
    task: T-137
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-11T05:34:55Z"
---

# Animate tiles that are building or testing

## Outcome

A Task tile at `status: in_progress` with `stage: build` shows a one-cell spinner and "Building…", and one at `stage: test` shows a spinner and "Testing…", in place of today's stage badge. Every other tile, and a board with no such Task, is exactly as today. The spinner moves only while at least one such tile is shown.

## User Check

Owner validation required. Open the board with a Task at build and confirm its tile shows a moving spinner with "Building…". Move it to test and confirm "Testing…". Move it on to check and confirm the tile shows today's badges and nothing moves. Open a project with no Task at build or test and confirm the board looks as it did before.

## Done When

1. The live label replaces the stage badge only for `in_progress` + `build` or `test`, beside `stageBadge` in `badges.go`. The text is the stage word from `stageLabel` plus "…", and the spinner frame comes first. It keeps the build and test accents. It is one terminal cell plus the words, cut to the tile width and never wrapped. Other badges on the tile (review outcome, blockers) are unchanged and keep their order.
2. A frame counter on `Model` is advanced by a tick `tea.Cmd`. The tick is scheduled after a load or reload only when a shown tile qualifies, and it does not reschedule once none does. An idle board runs no timer. Rendering and `Update` do no file work (ARCH-02). The frames and interval are named constants (STYLE-09).
3. Plain (non-TTY) output, help, the key map, and footer hints are unchanged, and no key is added. The committed golden tests in `internal/board/v2/golden_test.go` pass unchanged.
4. `.savepoint/visual-identity.md` records the spinner as the owner-approved exception to "running animation skip" (TPL-02).
5. Tests: the build label, the test label, the frame advancing on tick, audit/done/planned tiles unchanged, the tick not starting on an idle board and stopping after the last qualifying tile moves on, a narrow tile with no wrap, and a stage change between reloads. Each test covers its happy path and a failure-path case (TEST-02).
6. `make build && make test-fast` pass. Per-criterion evidence follows AGENTS.md's Verification Policy. Owner validation is recorded after the User Check.

## Context Files

`internal/board/v2/badges.go`, `internal/board/v2/card.go`, `internal/board/v2/card_test.go`, `internal/board/v2/model.go`, `internal/board/v2/update.go`, `internal/board/v2/view.go`, `internal/board/v2/golden_test.go`, `internal/styles/styles.go`, `.savepoint/visual-identity.md`, `.savepoint/objectives/O-048-animate-tiles/Objective.md`.

## Design References

Design §8 (TUI, badges, visual guardrail, "Task tile run animation"). O-048 Confirmed Design.

## Guardrails

ARCH-02, TPL-02, STYLE-07, STYLE-09, TEST-01, TEST-02, TEST-04, TEST-08.

## Implementation Plan

1. Add the live stage badge beside `stageBadge`, chosen in `TaskCard.badges`.
2. Add the frame counter, the tick message and command, and the self-stopping schedule in `update.go`.
3. Pass the frame to tile rendering without any IO.
4. Record the visual-identity exception.
5. Tests, golden re-run, gates.

## Boundaries

No run records, liveness checks, or file-age timeouts. No change to Objective rows, plain output, keys, options, or gates. No Task, Objective, Check, or router write.

## Technical Verification

Focused tests during iteration. Handoff gates per AGENTS.md's Verification Policy. No process, path, or file-replace work, so the O-046-style Windows evidence is not needed beyond the O-048 Full Objective Check's CI run. Evaluated later under `agent-skills/references/check-method.md`.

## Technical Evidence

**Commands run:** `make build && make test-fast` (exit 0, run twice, the second to confirm the exit code). Focused runs during iteration: `go test ./internal/board/...` and the new spinner tests. `make test-full` not run: no migration or platform-sensitive work.

**Per-criterion outcomes**

1. Live label: `liveStageBadge` (`badges.go`) applies only to `in_progress` + `build`/`test`. Words "Building…"/"Testing…" come from `liveStageWords`; spinner frame first; keeps the attention (build) and clear (test) accents; `NoWrap` cuts it to the tile width in `renderBadgeLines`. Review outcome and blocker badges keep their order (`TestLiveLabelKeepsOtherBadgesInOrder`).
2. Timer: `spinnerFrame` and `spinnerRunning` on `Model`; `spinnerTickMsg` advances the frame and reschedules only while `hasLiveTile()`; `applyLoad` calls `spinnerCmd()`, which starts one timer only when a live tile is in scope. No file work in Update or rendering. Frames and interval are named in `spinner.go` (STYLE-09).
3. Plain output, help, keys, footer unchanged: live drawing is opt-in (`TaskCard.Live`, set only by the board's `renderColumns`), so `badges()` and plain output are unchanged (`TestPlainOutputCarriesNoLiveLabel`). The committed golden files pass unchanged and none were modified.
4. `.savepoint/visual-identity.md` records the spinner as the owner-approved exception.
5. Tests in `spinner_test.go`: build label, test label, frame wrap, audit/done/planned unchanged, tick advances, tick stops when none live, idle board starts no timer, narrow tile no wrap, stage change between reloads. Each has a failure-path case.
6. `make build && make test-fast` pass.

**Existing tests changed:** `columns_view_test.go` now expects the live labels instead of the static BUILD/TEST badges on the board. `objectives_test.go` and `releases_test.go` helpers stop following a spinner tick, since a live tile's timer never ends by design.

**Files read:** the Task, Objective, `badges.go`, `card.go`, `model.go`, `update.go`, `view.go`, `column.go`, `visual-identity.md`, `Guardrails.md` (rule lines). **Extra reads beyond Context Files:** `column.go` (where tiles are drawn, needed to pass the frame), `fixture_test.go` and `view_test.go` helpers (to reuse fixtures), `Guardrails.md` (rule wording). **Files changed:** `badges.go`, `card.go`, `column.go`, `model.go`, `update.go`, `view.go`, new `spinner.go` and `spinner_test.go`, the three test files above, `visual-identity.md`.

**Limitations:** Not run in a real terminal; owner validation (the User Check) is still to do. In the narrow single-column layout the timer runs while any in-scope tile is live, even if the In Progress column is not the one on screen. A Task left at build keeps animating until moved on, as the Objective decided.

## Drift Notes

Record architecture deltas and reconcile through the planner before Check.
