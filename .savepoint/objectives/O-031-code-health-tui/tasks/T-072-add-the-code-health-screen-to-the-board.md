---
id: T-072
title: Add the Code Health screen to the board
objective: O-031
status: done
depends_on: [{task: T-070, requires: clear}, {task: T-071, requires: clear}]
owner_validation:
    required: true
    accepted_check: ""
planned_by: {role: planner, session: planning-o031-20261002}
complexity_tier: high
complexity_reason: New full-screen board surface with async load, Git comparison, refresh progress, and cancellation, all through explicit commands.
check_waiver:
    task: T-072
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-01T20:34:50Z"
---

# Add the Code Health screen to the board

## Outcome

Pressing `H` on the board opens a Code Health screen that answers "is my code staying healthy?" in plain words, lets the owner open any signal's details, and lets them refresh with visible progress and cancel without losing the last result.

## User Check

In a project with health configured, open the board, press `H`: the screen appears at once from saved results, then marks anything stale. Press Enter on a signal and read why it has its label. Press `R`, watch "n of 5" progress, press Esc mid-run, and confirm the previous result is still shown. Press `R` again and let it finish. In a project without health configured, `H` explains the feature and points to `savepoint health setup`. Esc returns to the board where you were.

## Done When

- `H` opens a full-screen Code Health view over the board, restoring the board cursor on Esc, mutually exclusive with the detail, Issues, Goal, and help overlays (same pattern as the Issues overlay).
- Opening dispatches one command that calls `codehealth.LoadDashboard`; a second command then calls `codehealth.DashboardFreshness`. Rendering does no IO (ARCH-02) and the board imports no snapshot or provider type.
- Not configured: explains in a few lines what Code Health is and to run `savepoint health setup`; `R` does nothing. First run: explains nothing is measured yet and offers `R`.
- Measured: overall label, measured time and origin (official or manual), the freshness line, then one row per signal with label, short reason, and trend word. Partial, failed, unavailable, optional, and not-configured rows are visible and never shown as Good. Labels use Atari-Noir tokens from `.savepoint/visual-identity.md`; no overall number.
- Enter or `v` on a row opens its details: full explanation, trend and comparison basis, affected areas (evidence references), outcome, required or optional, provider and version, scope, and time. A history section lists recent snapshots.
- `R` starts a manual refresh through an explicit command with a cancellable context. The view shows "Refreshing n of N: <signal> (<provider>)" from the progress callback and "Esc to cancel". Esc cancels and shows "Refresh cancelled; nothing was saved." Completion reloads the dashboard and freshness. A second `R` while running is ignored; `q`/ctrl+c during refresh cancels before quitting.
- A collection error (for example damaged history) is shown as a plain status line, not a crash, and leaves the previous view.
- Help (`?`) lists `H` on the board and the view's own keys.
- Model-level tests with injected dashboard, freshness, and refresh functions cover: not configured, first run, five signals, partial support, required and optional failures, stale/dirty/diverged/unknown freshness, comparable and reset trends, details, history, refresh progress messages in order, cancel saves nothing, refresh error, and cursor restore on close. A width test covers the view at the board's minimum width.

## Context Files

`.savepoint/objectives/O-031-code-health-tui/Objective.md`; `.savepoint/visual-identity.md`; `internal/board/v2/model.go`, `internal/board/v2/update.go`, `internal/board/v2/view.go`, `internal/board/v2/help.go`, `internal/board/v2/io.go`, `internal/board/v2/issues.go` (overlay pattern), `internal/board/v2/width.go`, `internal/board/v2/fixture_test.go`; `internal/codehealth/dashboard.go`, `internal/codehealth/dashboard_freshness.go`, `internal/codehealth/collect.go`; new `internal/board/v2/health.go`, `internal/board/v2/health_view.go`, `internal/board/v2/health_test.go`.

## Design References

Design section 8 (layout, keybindings, visual and border policy, persistence and refresh) and section 1 (Code Health readers); O-031 Confirmed Design Decisions.

## Guardrails

ARCH-02, ARCH-03, ARCH-04, DATA-03, CFG-02, DEP-02, STYLE-07, STYLE-09, TEST-01, TEST-02, TEST-04.

## Implementation Plan

1. Confirm T-070's dashboard API and T-071's progress callback and cancel error exist as described; return REPLAN REQUIRED otherwise.
2. Add the health overlay state and open/close handling mirroring the Issues overlay's origin restore.
3. Add load, freshness, and refresh commands in `io.go` behind injectable function fields so tests drive them with messages.
4. Stream progress from the refresh goroutine to the program through messages; hold the cancel function in the model.
5. Render overview, details, and history in `health_view.go`, copy kept in one table.
6. Add help rows and tests.

## Boundaries

No background or timed refresh, no file watcher on health files, no setup from the board, no AI wording, no Issue creation, no change to Code Health policy or `health check`.

## Technical Verification

Focused board tests while iterating; `make build && make test-fast` for handoff. Owner validation per User Check. Later evaluated per `agent-skills/references/check-method.md`.

## Technical Evidence

Commands run: `go test -race -count=3 ./internal/board/v2 -run Health` (pass); `make build && make test-fast` (pass, exit 0, no FAIL lines). Owner validation (User Check in a real terminal) has not been done by the executor.

Per criterion:
- `H` opens a full-screen view, Esc restores the cursor, exclusive with detail/Issues/Goal/help: `openHealth` guards every other overlay; `TestHealthEscRestoresTheBoardCursor`, `TestHealthIsExclusiveWithTheOtherOverlays`.
- Open dispatches `LoadDashboard`, then `DashboardFreshness`; no IO in rendering: `healthLoadCmd`/`healthFreshnessCmd` in `io.go`, injectable via `HealthFuncs`; `TestHealthShowsResultsBeforeFreshnessArrives`, `TestHealthStaleFreshnessForAnOlderSnapshotIsIgnored`. Board uses only `codehealth.Dashboard`, `CodeFreshness`, `Progress`, and the opaque `RepositoryIdentity` it hands back; no snapshot or provider-reader type.
- Not configured / first run: `TestHealthNotConfiguredExplainsAndRefreshDoesNothing`, `TestHealthFirstRunOffersRefresh`.
- Measured: overall, time and origin, freshness, rows with label/reason/trend, glyph plus word, no number, Atari-Noir tokens (`styles.Health*`): `TestHealthShowsFiveSignalsWithLabelsReasonsAndTrends`, `TestHealthNonGoodRowsAreNeverShownAsGood`.
- Details and history: `TestHealthDetailsExplainTheLabelAndEscReturnsToTheList`; history is listed under the rows as "RECENT CHECKS".
- Refresh: progress text in order, Esc cancels with the exact cancelled line, second R ignored, completion reloads, q/ctrl+c cancels first: `TestHealthRefreshReportsProgressInOrderAndReloadsOnCompletion`, `TestHealthEscCancelsARefreshAndKeepsThePreviousResult`, `TestHealthQuitDuringRefreshCancelsBeforeQuitting`.
- Errors are status lines: `TestHealthLoadErrorKeepsThePreviousView`, `TestHealthLoadErrorOnFirstOpenIsAStatusLineNotACrash`, `TestHealthRefreshErrorIsAStatusLine`.
- Help: `TestHealthHelpListsTheKeys`. Width at the board's minimum: `TestHealthFitsTheBoardsMinimumWidth`.

Files changed: `internal/board/v2/{health.go,health_view.go,health_test.go,io.go,model.go,update.go,view.go,help.go}`, `internal/styles/styles.go` (Health* styles), `internal/codehealth/dashboard.go` (exported `CapabilityText`).

Extra reads/edits beyond Context Files: `internal/styles/styles.go` (read and edited, to add the Atari-Noir label styles the plan requires; `visual-identity.md` names the tokens but the board reaches colors only through `styles`), `internal/board/v2/boundary_test.go` (read, to learn the badge-vocabulary rule that forced `styles.Health*`), `internal/healthcheck/healthcheck.go` and `cmd`-side wiring (read, to copy how a collection is invoked), `internal/codehealth/classification.go` (read, to confirm a private capability label was not the display text).

Limitations: terminal behavior (colors, real refresh with real tools, real Git freshness) is unverified until the owner runs the User Check. A cancelled refresh relies on T-071's process-group kill; on `q` the program exits right after cancelling, so a tool process is stopped by that cancellation, not waited for. The "reset" trend wording is shown as `codehealth` words it; the test uses injected text.

## Drift Notes

New board surface and key `H`; Design section 8 keybindings and the AGENTS.md Codebase Map are reconciled at the Objective Check.
