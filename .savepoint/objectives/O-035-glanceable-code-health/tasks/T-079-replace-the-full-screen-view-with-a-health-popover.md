---
id: T-079
title: Replace the full-screen view with a Health popover
objective: O-035
status: done
depends_on: [{task: T-078, requires: clear}, {task: T-076, requires: clear}]
owner_validation:
    required: true
    accepted_check: ""
planned_by: {role: planner, session: planning-o035-20261002}
complexity_tier: high
complexity_reason: Replaces a shipped screen, its detail view and its scrolling; most of the board's overlay, help, hint and test surface is touched.
check_waiver:
    task: T-079
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-01T22:19:06Z"
---

# Replace the full-screen view with a Health popover

## Outcome

`H` opens a small fixed-height popover over the dimmed board listing the five signals with value, mark, sparkline and aim, and a short explanation of the selected signal; it never scrolls and Esc returns to where the cursor was.

## User Check

On a project with a saved official snapshot, press `H`: a popover shows five rows, each with a mark, value, sparkline and aim, and a `Re-run: savepoint health check O-0xx` line you can copy. Move with up and down: the explanation under the list changes and says what the result means, whether it blocks sign-off, where, and one next step. Try 80x20, 80x24 and 80x40: nothing scrolls and nothing is cut off. Esc closes it with the board cursor where it was.

## Done When

- `H` opens a popover drawn over the dimmed board with the existing overlay helper, at most 72 columns wide and 17 lines tall, fixed regardless of data. The title shows the headline and the measured date and origin.
- Each row is one line: mark (`✓`, `~`, `✗`, `?`), the plain signal name, `Value`, the sparkline with its last block coloured by label and `SparkWord`, and `Aim`. Rows truncate with an ellipsis; they never wrap or scroll.
- The selected row's block under a rule shows `Question`, `Meaning`, `SignOff`, `Where` and `NextStep` as at most five lines, each truncated, and the dashboard-level sign-off line.
- A last line under the list shows the command to re-run the official check, ready to copy: `Re-run: savepoint health check O-035`, using the Objective in view on the board, or the router's Objective, and `O-###` when neither exists. It is built in the board's load path from state the model already holds, not while rendering, and it is not shown while a refresh runs. `R` is described as a manual refresh that does not count as an official check.
- Up and down move the selection; no key scrolls. The full-screen detail view, its Enter and `v` keys, the `Offset` and `Scrolled` state, the overview history block and their tests are removed, and Help and the hints describe the new keys.
- Esc closes the popover and restores the board cursor as it was. Provider names, snapshot hashes and raw timestamps appear nowhere in it.
- Tests cover the measured state at 80x20, 80x24 and 80x40 with every line within the width and the popover within the height, selection changes, each label's mark and colour, the sign-off lines, the re-run command for a selected Objective, for the router's Objective and with neither, truncation of a long explanation, Esc restoration, and the removed keys doing nothing.

## Context Files

`.savepoint/objectives/O-035-glanceable-code-health/Objective.md`; `internal/board/v2/health.go`, `internal/board/v2/health_view.go`, `internal/board/v2/health_test.go`, `internal/board/v2/view.go`, `internal/board/v2/update.go`, `internal/board/v2/help.go`, `internal/board/v2/releases.go` (read only: overlay helper), `internal/board/v2/view_test.go`, `internal/board/v2/footer_test.go`; `internal/styles/styles.go`; `internal/codehealth/dashboard.go` (read only: row fields).

## Design References

O-035 Confirmed Design Decisions (placement, affected files, verdict basis); I-104 (scroll defect that no longer applies).

## Guardrails

ARCH-02, STYLE-01, STYLE-05, STYLE-09, TEST-01, TEST-02, TEST-04, TEST-05.

## Implementation Plan

1. Confirm the dashboard fields from the three data Tasks and the chip exist; return REPLAN REQUIRED otherwise.
2. Add the popover renderer on the existing overlay helper and route `H`, selection and Esc through the current Health overlay state.
3. Remove the detail view, the scrolling state, the overview history block and their tests; update Help and hints.
4. Add the tests in Done When, including a regression that none of the removed keys opens or scrolls anything.

## Boundaries

No first-run, not-configured, refresh or freshness wording yet, and no history view; those arrive in the next Tasks. No new keys other than the popover's own selection.

## Technical Verification

Focused tests during iteration (`make test-focused`); `make build && make test-fast` for ordinary handoff. Later evaluated per `agent-skills/references/check-method.md`.

## Technical Evidence

Executed 2026-10-02.

Per criterion:
- Popover on the overlay helper, at most 72 wide and 17 tall, title with headline, date and origin: `renderHealthPopover` draws over `dimV2Lines(base)`; `TestHealthPopoverFitsEverySupportedSize` (80x20, 80x24, 80x40, all five selections: every line within width, 17-line popover, title text present).
- One-line rows with mark, name, value, spark with last block coloured, word, aim: `TestHealthPopoverRowsShowMarkValueSparkAndAim`, `TestHealthPopoverColoursEachLabel`.
- Selected-row block (Question, Meaning, sign-off, Where, Next, at most five lines) and dashboard sign-off line: `TestHealthPopoverSelectionChangesTheExplanation`, `TestHealthPopoverSignOffLines`; truncation of a long Meaning in the fit test.
- Re-run line built at open (`healthRerunCommand`) for the sidebar Objective, the router's Objective, and `O-###`; hidden during refresh: `TestHealthPopoverRerunCommand`. Help describes `R` as a manual refresh that is not an official check.
- Up/down only; Detail, Offset, Scrolled, history block, Enter/`v` and their tests removed; Help and hints updated: `TestHealthPopoverRemovedKeysDoNothing`.
- Esc restores the cursor: existing `TestHealthEscRestoresTheBoardCursor`. No provider, snapshot hash or raw timestamp shown: asserted in `TestHealthPopoverRowsShow...`.

Commands: `go test ./internal/board/v2` ok; `make build && make test-fast` passed.

Files read: the Context Files plus `internal/codehealth/dashboard_copy.go`, `chip.go`, `spark.go` (extra reads, to confirm row fields) and `internal/board/v2/model.go`, `width.go`, `actions.go`, `objectives.go` (extra reads, for overlay geometry and the Objective in view).
Files changed: `internal/board/v2/health.go`, `health_view.go`, `health_test.go`, `health_popover_test.go` (new), `help.go`, `update.go`, `view.go`.

Limitations: Help takes precedence over the popover (drawn instead of it). Not-configured, first-run, loading and notice states reuse the existing wording inside the popover; their redesign is T-080. The re-run command is chosen when the popover opens, not on later cursor moves. Not viewed in a real terminal; owner validation pending. `make test-full` not run (not required for ordinary handoff).

## Drift Notes

New module or architecture delta beyond the documented Codebase Map, reconciled through the planner before Check.
