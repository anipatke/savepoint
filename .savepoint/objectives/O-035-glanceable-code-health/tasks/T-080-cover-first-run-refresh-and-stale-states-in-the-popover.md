---
id: T-080
title: Cover first run, refresh and stale states in the popover
objective: O-035
status: done
depends_on: [{task: T-079, requires: clear}]
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: planning-o035-20261002}
complexity_tier: medium
complexity_reason: Moves existing refresh, cancel, freshness and first-run behaviour into a fixed-height frame without losing any of it.
check_waiver:
    task: T-080
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-01T22:21:08Z"
---

# Cover first run, refresh and stale states in the popover

## Outcome

The popover handles not set up, no check yet, refresh progress and cancel, stale code, partial and unknown signals, and load errors in the same fixed frame.

## User Check

In a project with no health setup, press `H`: the popover explains Code Health and points to `savepoint health setup`. With setup but no snapshot it offers `R`. Press `R` and watch `Refreshing 2 of 5: coverage`; press Esc to cancel and see the last result kept. Change a file and reopen: a short line says the code has moved on since the check.

## Done When

- Not set up shows two short lines that explain the feature and name `savepoint health setup`; no check yet offers `R` and shows the same copyable `savepoint health check O-###` re-run line. Neither opens tools or scans.
- `R` starts a refresh; one line inside the popover shows `Refreshing n of 5: <signal>`, the footer offers Esc to cancel, and Esc cancels, saves nothing and keeps the last saved result. Ctrl+C and `q` cancel before quitting, including from Help (I-105 behaviour is kept).
- The freshness comparison runs as before as an explicit command; when code has moved on, a single line says so and rows from a different branch or older code are marked, with the existing wording preserved.
- Partial, unknown and not-measured rows keep `?` or `~` marks and their plain `Meaning`; a load error shows one line and keeps the last good result.
- Everything fits the fixed frame at 80x20, 80x24 and 80x40.
- Tests cover each state, cancel keeping the last result, quit during refresh from the popover and from Help, stale and divergent notices, a load error, and every line within width and height.

## Context Files

`.savepoint/objectives/O-035-glanceable-code-health/Objective.md`; `internal/board/v2/health.go`, `internal/board/v2/health_view.go`, `internal/board/v2/health_test.go`, `internal/board/v2/view.go`, `internal/board/v2/update.go`, `internal/board/v2/help.go`; `internal/codehealth/dashboard_freshness.go` (read only: freshness wording).

## Design References

O-031 Confirmed Design Decisions (opening, refresh, cancel); I-105.

## Guardrails

ARCH-02, DATA-03, STYLE-09, TEST-01, TEST-02, TEST-04.

## Implementation Plan

1. Confirm the popover frame and the existing refresh and freshness state still exist after the previous Task; return REPLAN REQUIRED otherwise.
2. Re-home first-run, not-configured, progress, cancel, freshness and error lines in the popover frame.
3. Keep refresh, cancel and quit behaviour identical to today and add the tests in Done When.

## Boundaries

No change to how refresh collects, to freshness computation, or to snapshot creation; no history view.

## Technical Verification

Focused tests during iteration (`make test-focused`); `make build && make test-fast` for ordinary handoff. Later evaluated per `agent-skills/references/check-method.md`.

## Technical Evidence

Commands: `go test ./internal/board/v2` (pass); `make build && make test-fast` (pass).

Per criterion:
- Not set up / no check yet: two short lines each; no check yet offers R and shows the re-run line; not set up shows no re-run line. TestHealthPopoverFirstRunStatesFitEverySize.
- Refresh/cancel/quit: behaviour unchanged; existing tests (Esc cancel keeps last result, quit and Ctrl+C from Help during refresh) pass; TestHealthPopoverRefreshAndNoticesFitEverySize adds progress line, cancelled notice and load error at 80x20/24/40.
- Freshness: wording unchanged; moved-on and other-branch now mark every value with `*` and the selected signal with "(measured on other code)"; matches/unknown do not. TestHealthPopoverMarksRowsMeasuredOnOtherCode.
- Partial/unknown marks and Meaning kept: TestHealthPopoverKeepsUnknownAndPartialMarks.
- Fixed frame: every state at 80x20, 80x24, 80x40 within width/height and exactly healthPopoverHeight lines.

Files changed: internal/board/v2/health_view.go, internal/board/v2/health_popover_test.go. Read beyond Context Files: internal/board/v2/health_popover_test.go, internal/codehealth/dashboard.go (row Freshness fields) to place tests and the older-code mark.

Limitations: the older-code mark follows the whole-code freshness state, not per-row report staleness (`FreshnessText`); not run in a real terminal; make test-full not run.

## Drift Notes

New module or architecture delta beyond the documented Codebase Map, reconciled through the planner before Check.
