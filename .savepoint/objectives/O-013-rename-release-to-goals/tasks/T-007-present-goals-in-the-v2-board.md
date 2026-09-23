---
id: T-007
title: Present Goals in the V2 board
objective: O-013
planned_by: {role: planner, session: goals-terminology-20260921}
status: planned
complexity_tier: medium
complexity_reason: "The V2 board exposes the context in interactive, overlay, reload, and non-TTY paths, all of which must share one vocabulary without changing filtering behavior."
depends_on: [{task: T-006, requires: clear}]
owner_validation:
    required: true
---

# T-007: Present Goals in the V2 board

## Outcome

The V2 interactive and plain board presents the selected delivery context as a
Goal while retaining the existing Objective filtering, selection persistence,
and Release-backed data behavior.

## User Check

Open the V2 board with an existing R### project. Confirm that the selector,
selected-context line, detail view, help, reload/status messages, and plain
output say Goal/Goals; confirm `g` opens the selector and that any retained `r`
alias behaves identically. Confirm Objective and Task membership has not
changed.

## Done When

- Selector, selected-context header, detail overlay, help, and status/error
  messages use the Goal vocabulary.
- `g` is documented and dispatched as the canonical selector shortcut; a
  retained `r` alias has no separate behavior or public wording.
- Interactive and non-TTY output agree on the Goal label and selected record.
- Reload, missing-selection, narrow-terminal, and empty-Goal cases remain
  understandable and deterministic.
- Existing R### data continues to drive the same Objective and Task filtering;
  no rendering path performs new IO or gate evaluation.
- V1 board behavior and historical Release-specific compatibility surfaces are
  left unchanged.

## Context Files

`internal/board/v2/releases.go`, `internal/board/v2/releases_test.go`, `internal/board/v2/model.go`, `internal/board/v2/view.go`, `internal/board/v2/view_test.go`, `internal/board/v2/help.go`, `internal/board/v2/update.go`, `internal/board/v2/objectives.go`, `internal/board/v2/objectives_test.go`, `internal/board/v2/plain.go`, `internal/board/v2/run.go`, `internal/board/v2/run_test.go`, `internal/board/v2/load_test.go`, `internal/board/v2/watch_test.go`.

## Design References

Design sections 2, 3, 6, 8, 11, and 13.

## Guardrails

DATA-02, ARCH-02, ARCH-03, TPL-02, TEST-01, TEST-02, TEST-04, TEST-06,
TEST-08, STYLE-07, STYLE-09, and STYLE-10.

## Implementation Plan

1. Centralize the user-facing Goal labels and shortcut vocabulary at the V2
   board presentation/dispatch boundary rather than scattering new strings.
2. Update selector, header, detail, help, status, reload, and plain-output
   paths while keeping the existing Release-backed model and index calls.
3. Cover the canonical `g` key, optional `r` compatibility alias, selected
   context, empty data, reload, narrow layout, and non-TTY output.
4. Run focused V2 board tests, `git diff --check`, `make build`, and `make test`.

## Boundaries

No data model rename, persisted-file migration, new Goal membership map, gate
policy change, V1 board rewrite, or UI redesign beyond terminology and the
selector shortcut.

## Technical Verification

Focused `go test ./internal/board/v2` with interactive, overlay, reload,
narrow-width, and non-TTY assertions; existing R### fixture load;
`git diff --check`; `make build`; and `make test`.

## Technical Evidence

Pending execution: named UI paths and test cases, files read/changed, command
results, and limitations.

## Drift Notes

If the board cannot expose Goals without changing persisted identity or gate
ownership, stop and return REPLAN REQUIRED rather than introducing a second
domain model.
