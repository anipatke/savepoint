---
id: T-081
title: Add a history list behind h in the popover
objective: O-035
status: done
depends_on: [{task: T-080, requires: clear}]
owner_validation:
    required: true
    accepted_check: ""
planned_by: {role: planner, session: planning-o035-20261002}
complexity_tier: low
complexity_reason: One more fixed-height view over a list the dashboard already bounds to ten entries.
check_waiver:
    task: T-081
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-01T22:23:43Z"
---

# Add a history list behind h in the popover

## Outcome

Pressing `h` in the popover shows the last ten checks, newest first, with date, official or manual, and overall label; `h` or Esc returns to the signals.

## User Check

In a project with several snapshots, open the popover and press `h`: ten rows at most, newest first, manual refreshes dimmed. Press `h` or Esc to go back to the signals; press Esc again to close.

## Done When

- `h` switches the popover body to the history list; `h` or Esc returns to the signals, and a second Esc closes the popover.
- Each row shows `WhenText`, `OriginText` and the overall label with its colour, at most 10 rows, newest first. Manual refreshes are dimmed with a note that they do not feed the sparklines.
- With no history it says so in one line. The view uses the same fixed frame and never scrolls; a long list truncates at ten.
- Help and the hints name `h`. No other key changes.
- Tests cover more than ten snapshots, exactly one, none, manual versus official rendering, `h` and Esc return paths, and fixed width and height at 80x20, 80x24 and 80x40.

## Context Files

`.savepoint/objectives/O-035-glanceable-code-health/Objective.md`; `internal/board/v2/health.go`, `internal/board/v2/health_view.go`, `internal/board/v2/health_test.go`, `internal/board/v2/help.go`, `internal/board/v2/view.go`; `internal/codehealth/dashboard.go` (read only: history entries).

## Design References

O-035 Success Conditions (history behind `h`); O-031 Confirmed Design Decisions (short history).

## Guardrails

ARCH-02, STYLE-09, TEST-01, TEST-02, TEST-04.

## Implementation Plan

1. Confirm `Dashboard.History` is bounded to ten with `WhenText`; return REPLAN REQUIRED otherwise.
2. Add a history mode to the Health overlay state and its renderer; wire `h` and Esc.
3. Update Help and hints and add the tests in Done When.

## Boundaries

No snapshot deletion, filtering or comparison view, and no change to what is stored or how many snapshots are kept.

## Technical Verification

Focused tests during iteration (`make test-focused`); `make build && make test-fast` for ordinary handoff. Later evaluated per `agent-skills/references/check-method.md`.

## Technical Evidence

Commands: `make build && make test-fast` (exit 0, 2026-10-02); focused `go test ./internal/board/...` during iteration.

Per criterion:
- `h` toggles history; `h` or Esc returns to the signals; a second Esc closes (`TestHealthHistoryReturnPaths`). Esc still cancels a refresh first.
- Rows show `WhenText`, `OriginText`, glyph and overall label (coloured for official), at most 10, newest first; manual rows render in the dim `CardMeta` style with a note that they do not feed sparklines (`TestHealthHistoryShowsAtMostTenNewestFirst`, `...RendersManualRefreshesDimmed`).
- No history prints "No checks yet."; the frame stays 17 lines and truncates at ten (`...OneAndNoSnapshots`, 14-snapshot case at 80x20, 80x24, 80x40).
- Help and hints name `h` (`TestHealthHistoryNamedInHelpAndHints`); no other key changed; up/down are inert while history shows.

Step 1 confirmed: `Dashboard.History` is bounded by `MaxDashboardHistory = 10` with `WhenText`.

Files read: Context Files plus `internal/board/v2/update.go` (key dispatch, extra read to confirm `h` reaches the Health handler before the board's column keys) and `health_popover_test.go`/`health_test.go` fixtures. Changed: `health.go`, `health_view.go`, `help.go`, `view.go`, new `health_history_test.go`.

Limitations: dimming is by style only; tests run without a colour profile so they assert text and layout, not ANSI dimming. Not manually tried in a real terminal; owner validation is required.

## Drift Notes

New module or architecture delta beyond the documented Codebase Map, reconciled through the planner before Check.
