---
id: T-037
objective: O-022
title: Keep the board inside one Goal
status: done
complexity_tier: medium
complexity_reason: Removes the board's Goal-less branches across the TUI and plain output without touching selection rules.
depends_on: [{task: T-036, requires: clear}]
owner_validation:
    required: true
    accepted_check: ""
planned_by: {role: planner, session: o022-goal-required-20260925}
check_waiver:
    task: T-037
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-09-25T01:35:05Z"
---

# Keep the board inside one Goal

## Outcome

The board only ever shows one Goal's Objectives and Tasks. With no valid Goal
selected, it shows the `Choose a Goal` diagnostic and how to pick one, not a
project-wide list of the work.

## User Check

Open the board on a scratch project whose router names no Goal. Confirm the
sidebar and columns are empty, the Next area reads `Choose a Goal`, and `g`
opens the selector. Pick a Goal and confirm the view fills with that Goal's
Objectives. Repeat with a project that declares no Goals and confirm the board
points to `savepoint doctor`. Run `savepoint board` non-interactively and
confirm the plain output says the same thing.

## Done When

- With no valid router Goal, `taskIDsInReleaseView`, the Objective sidebar, and
  plain output list nothing project-wide. The view shows the T-036 diagnostic
  and the `g` hint, or a doctor pointer when zero Goals exist.
- With a Goal selected and no Objective filter, the `ALL OBJECTIVES` label and
  the plain `Selected:` line say the view is scoped to that Goal. The
  unscoped `Selected: all Objectives` text is gone.
- When Objectives without a Goal exist, the board shows a one-line notice with
  their count and a pointer to doctor. The notice is drawn from the T-036 index
  fact, not recomputed.
- The Goal selector still cannot clear the Goal. A test pins that behaviour.
- Tests cover no Goal, zero Goals, a selected Goal with and without an
  Objective filter, the unassigned notice, narrow widths, and matching TTY and
  non-TTY text.

## Context Files

`internal/board/v2/releases.go`, `internal/board/v2/releases_test.go`,
`internal/board/v2/card.go`, `internal/board/v2/view.go`,
`internal/board/v2/view_test.go`, `internal/board/v2/plain.go`,
`internal/board/v2/objectives.go`, `internal/board/v2/objectives_test.go`,
`internal/board/v2/next_panel.go`, `internal/board/v2/next_panel_test.go`,
`internal/board/v2/update.go`, `cmd/board_test.go`.

## Design References

Design section 8 (TUI) and 10; O-022 Success Conditions.

## Guardrails

ARCH-01, ARCH-03, TEST-01, TEST-02, TEST-04, TEST-08, STYLE-07.

## Implementation Plan

1. Confirm T-036's diagnostic and index fact exist; return REPLAN REQUIRED if
   they do not.
2. Remove the empty-Goal branch from the board's card, sidebar, and plain
   grouping. Render the diagnostic and hint instead.
3. Relabel the no-Objective-filter view as Goal-scoped, in the TUI and plain
   output.
4. Add the unassigned-Objectives notice.
5. Run focused tests while iterating, then `make build && make test-fast`.

## Boundaries

No new selection rules, no board action that edits an Objective's Goal, and
no Objective ranking (O-020).

## Technical Verification

Focused tests during iteration; `make build && make test-fast` for handoff.

## Technical Evidence

- Started 2026-09-25. Router selection remains O-022/T-037 with `release: R-006`; no selection keys or Goal ownership were changed.
- Acceptance evidence:
  - No valid Goal yields no Objective rows or Task IDs, shows the canonical missing-Goal diagnostic and `g` hint, and prints no project-wide records: `TestBoardWithoutAValidGoalShowsNoProjectWideRecords`.
  - Zero Goals yields a `savepoint doctor` pointer in TUI and plain output: `TestBoardWithZeroGoalsPointsToDoctor`, `TestViewEmptyTemplateProjectOpensWithThreeEmptyColumns`.
  - The unfiltered Goal label and plain `Selected:` line name the Goal, while the Objective-filtered view omits that label: `TestGoalScopedSelectionWithObjectiveFilter`, `TestReleaseSelectionFiltersIndexedObjectivesAndPersistsOnlyRouterContext`.
  - The one-line unassigned Objective notice uses `V2Index.ObjectivesWithoutGoal` and appears in both renderers: `TestUnassignedGoalNoticeUsesIndexFactInBothRenderers`.
  - The selector has no empty choice and ignores clear keys: `TestGoalSelectorCannotClearTheCurrentGoal`.
  - Narrow chrome, selector width, and matching TTY/non-TTY guidance are covered by `TestGoalScopedChromeFitsNarrowWidths`, `TestGoalSelectorFitsNarrowBoardWidths`, `TestBoardNextAndResumeReportTheSameAnswer`, and `TestBuiltBoardAndResumeReportTheSameAnswer`.
- Verification:
  - `go test ./internal/board/v2` — passed.
  - `go test ./... -run TestGoalScopedSelectionWithObjectiveFilter` — passed.
  - `go test . -run 'Test(BoardNextAndResumeReportTheSameAnswer|BuiltBoardAndResumeReportTheSameAnswer)'` — passed.
  - `make build && make test-fast` — passed on the final run. Two earlier runs exposed the outdated zero-Goal parity expectation; both source and built-binary parity cases were updated and then passed.
  - Two `make test-focused TEST=...` attempts returned package-setup errors from the wrapper; direct equivalent Go test commands above passed.
- Extra reads, logged before inspection: `internal/data/project.go` and `internal/data/next.go` (confirm T-036's `ObjectivesWithoutGoal` index fact and `SelectionReleaseMissing` diagnostic); `internal/board/v2/fixture_test.go` (understand `writeNavigationProject` Goal setup); `internal/board/v2/card_test.go` and `internal/board/v2/issues_test.go` (repair focused-test fixtures that expected cards without a Goal); `main_board_next_parity_test.go` (align the cross-surface zero-Goal expectation with the Task).
- Files read within the Task context: `.savepoint/router.md`, O-022 `Objective.md`, this Task, `.savepoint/Guardrails.md`, `internal/board/v2/releases.go`, `releases_test.go`, `card.go`, `view.go`, `view_test.go`, `plain.go`, `objectives.go`, `objectives_test.go`, `next_panel.go`, `next_panel_test.go`, `update.go`, and `cmd/board_test.go`; plus `agent-skills/savepoint-task/SKILL.md`.
- Files changed: this Task record; `internal/board/v2/card.go`, `card_test.go`, `issues_test.go`, `next_panel.go`, `objectives.go`, `objectives_test.go`, `plain.go`, `releases.go`, `releases_test.go`, `update.go`, `view.go`, `view_test.go`; and `main_board_next_parity_test.go`.
- Limitations: the documented interactive scratch-project User Check was not run in a live TTY; model-rendering and non-TTY behavior were tested. No optional independent Task Check, owner waiver, or technical CLEAR is recorded.

## Drift Notes

None yet.
