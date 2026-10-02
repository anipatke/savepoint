---
id: T-076
title: Show a Health chip in the board header
objective: O-035
status: done
depends_on: []
owner_validation:
    required: true
    accepted_check: ""
planned_by: {role: planner, session: planning-o035-20261002}
complexity_tier: medium
complexity_reason: A new cheap data read on every board load plus a header layout that must degrade at narrow widths.
check_waiver:
    task: T-076
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-01T22:10:42Z"
---

# Show a Health chip in the board header

## Outcome

The board header always shows a small Health chip, coloured by overall health, that costs no extra row and reads saved data only.

## User Check

Open the board in a project with a saved health snapshot: the header shows `Health 3/5` in a colour that matches the worst signal. Shrink the terminal: the chip loses its words and then disappears before the counts do. Open a project with no health setup: the chip says so quietly.

## Done When

- A new `codehealth` lookup returns a small chip value from saved data only: state (not set up, no check yet, measured), the overall label, the number of Good signals, and the number of signals. It reads the saved configuration and the newest snapshot and runs no tool or Git command.
- Reading must stay cheap as official snapshots accumulate: if the store has no way to read just the newest snapshot, add one rather than loading every snapshot on each board load.
- The board loads the chip in `loadProject`, not while rendering. After a health refresh completes, the chip updates without restarting the board.
- The header shows `♥ Health 3/5` coloured Good, Watch, Needs Attention or dim by overall label; `♥ Health: no check yet` and `♥ Health: not set up` are dim. The existing counts, their order and the one-row header are unchanged.
- At narrow widths the chip shortens to `♥ 3/5`, then disappears before any count is truncated. Nothing wraps or adds a row at 80 columns.
- A missing or unreadable health directory is not a board error and shows the not-set-up chip.
- Tests cover the three states, each overall colour, 80 and 100 columns, narrow degradation, the unreadable-storage case, and the chip refreshing after a refresh result.

## Context Files

`.savepoint/objectives/O-035-glanceable-code-health/Objective.md`; `internal/codehealth/chip.go` (new), `internal/codehealth/chip_test.go` (new), `internal/codehealth/storage.go`; `internal/codehealth/dashboard.go` (read only: label wording); `internal/board/v2/view.go`, `internal/board/v2/model.go`, `internal/board/v2/load.go`, `internal/board/v2/health.go`, `internal/board/v2/view_test.go`, `internal/board/v2/load_test.go`, `internal/board/v2/health_test.go`; `internal/styles/styles.go`.

## Design References

O-035 Confirmed Design Decisions (placement); O-031 Check summary decision (read-only lookup during load).

## Guardrails

ARCH-02, ARCH-03, DATA-03, STYLE-09, TEST-01, TEST-02, TEST-04.

## Implementation Plan

1. Confirm `HealthGood`, `HealthWatch` and `HealthUnknown` styles and `renderHeader` exist; return REPLAN REQUIRED otherwise.
2. Add the chip lookup and, if needed, a newest-snapshot read in `storage.go`.
3. Carry the chip on the project state from `loadProject`; update it when a health load or refresh result arrives.
4. Render it in `renderHeader` with the width fallbacks.
5. Add the tests in Done When.

## Boundaries

No popover, no key bindings, no change to the existing header counts or to the full-screen view.

## Technical Verification

Focused tests during iteration (`make test-focused`); `make build && make test-fast` for ordinary handoff. Later evaluated per `agent-skills/references/check-method.md`.

## Technical Evidence

Gate: `make build && make test-fast` exit 0 (2026-10-02, after the final edit).

Per criterion:
- Lookup from saved data only: `codehealth.LoadChip` / `Dashboard.Chip` in `chip.go` (states not set up, no check yet, measured; label, Good count, signal count). Reads config and one snapshot, runs no tool or Git. Tests: `TestLoadChipNotSetUpWithoutConfig`, `TestLoadChipNoCheckYetWithConfigAndNoSnapshots`, `TestLoadChipMeasuredCountsGoodOfSignals`, `TestDashboardChipMatchesStates`.
- Cheap newest read: `Store.LatestSnapshot` picks by file modification time and decodes only that file. `TestLatestSnapshotReturnsNewestWithoutLoadingTheOthers` damages an older file and still passes; `TestLatestSnapshotAbsentStorage`.
- Loaded in `loadProject` (`ProjectState.HealthChip`); updated on every good `healthLoadedMsg`, which a finished refresh triggers. `TestHeaderChipRefreshesAfterRefreshResult`.
- Header wording/colour: `headerChipForms`, `chipStyle` in `view.go`. `TestHeaderChipStatesAndWording`, `TestHeaderChipColourFollowsOverallLabel` (Good, Watch, Needs Attention, Unknown/dim).
- Narrow widths: `TestHeaderChipStaysOneRowAt80AndDegradesBeforeCounts` checks widths 100, 80, 70, 60, 50: header height equals the chip-less header, no overflow, full -> `♥ 3/5` -> gone in order, never kept while a count is cut. At 80 the short form shows.
- Unreadable storage: `TestHeaderChipUnreadableStorageIsNotABoardError`, `TestLoadChipUnreadableStorageIsNotSetUp`.

Files changed: `internal/codehealth/chip.go`, `chip_test.go` (new); `internal/board/v2/view.go`, `load.go`, `health.go`, `health_test.go`, `detail_test.go`.

Extra reads: `internal/codehealth/history.go`, `snapshot.go` fields, `testdata` fixtures, `internal/board/v2/update.go` (message routing), `detail_test.go` (two failing tests) - to find the snapshot summary shape, how health messages reach the model, and why two old tests broke.

Limitations:
- `detail_test.go` is outside the Context Files: two tests searched the whole screen for "Health:", which the chip's own wording now matches; they now ignore the `♥` line.
- "Newest" uses file modification time, not `created_at`; a restored or copied snapshot set could order differently from the dashboard until the dashboard reloads. Signals count equals dashboard rows (configured instances, or one placeholder per unconfigured signal).
- Non-measured chips (`not set up`, `no check yet`) have one form and vanish whole when narrow; only the measured chip has a shortened form.
- No Task Check requested; `make test-full` not run (not a migration/platform Task).

## Drift Notes

New module or architecture delta beyond the documented Codebase Map, reconciled through the planner before Check.
