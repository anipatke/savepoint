---
id: T-133
title: Show the runs list and each run's live activity
objective: O-046
status: planned
depends_on: [{task: T-132, requires: clear}]
owner_validation: {required: true}
planned_by: {role: planner, session: planning-o046-20261011}
---

# Show the runs list and each run's live activity

## Outcome

With the feature on, `A` opens a Runs overlay listing runs newest first with Task, state, elapsed time, and last event; `enter` opens a live activity view of that run's real events that updates as `events.jsonl` grows, stays responsive on a 50 MB event file, and keeps only a bounded window in memory.

## User Check

Owner validation required. During a real run, open `A`, watch the row's last event change, press `enter` and see reads, edits, and commands appear without a restart; close and reopen the board mid-run and confirm the newer events are shown.

## Done When

1. Read-only helpers in `internal/run` list run directories with their decoded `Record` (and the `runner_lost` inference) and read `events.jsonl` from a byte offset, returning display events decoded with the provider adapter plus the new offset; a trailing partial line is left for the next read. No board code parses provider lines.
2. The Runs overlay follows the existing overlay pattern (nullable struct on `Model`, mutually exclusive with Issues, Health, Options, Detail, and the Goal selector; origin restored on close). `↑↓` move, `enter` opens activity, `esc` returns (from activity to the list, from the list to the board).
3. Reads are `tea.Cmd`s triggered by open and by the watcher's debounced reloads; the reducer only appends to a bounded window of 500 events; rendering does no IO and caps each line at terminal width (ARCH-02, ARCH-05).
4. Layout follows review §4.1 and §4.2 within the Atari-Noir rules: real event kinds only, no percentages; orange running, dim for other non-final states; follow mode at the end, scroll otherwise.
5. Tests: empty `runs/`, malformed `run.yml` named rather than crashing (DATA-03), incremental read across appends and a partial line, burst coalescing through the debounce, a 50 MB file opened in bounded time and memory, and rendering snapshots. T-131's flag-off golden tests still pass.
6. `make build && make test-fast` and fresh `make test-full` pass. Per-criterion evidence follows AGENTS.md's Verification Policy; owner validation is recorded after the User Check.

## Context Files

`internal/run/events.go`, `internal/run/claude.go`, `internal/run/record.go`, `internal/run/lost.go`, `internal/run/sanitize.go`, `internal/board/v2/model.go`, `internal/board/v2/update.go`, `internal/board/v2/view.go`, `internal/board/v2/help.go`, `internal/board/v2/watch.go`, `internal/board/v2/issues.go`, `internal/board/v2/issues_view.go`, `internal/board/v2/options.go`, `.savepoint/visual-identity.md`, `docs/reviews/agent-orchestration-pass-1.md`, `.savepoint/objectives/O-046-start-and-watch-a-run-from-the-board-opt-in/Objective.md`.

## Design References

Design §8 (TUI, visual guardrail, border and color policy). Review §4.1, §4.2, and "Keeping the TUI responsive". O-046 Confirmed Design (reading runs).

## Guardrails

ARCH-02, ARCH-04, ARCH-05, DATA-03, STYLE-01, TEST-01, TEST-02, TEST-04, TEST-08.

## Implementation Plan

1. Add the run listing and offset reader to `internal/run` with tests, including the 50 MB case.
2. Add the Runs overlay state, keys, and view in new board files.
3. Add the activity view with the bounded window and follow mode.
4. Hook reads to open and to debounced reloads.
5. Rendering snapshots, golden re-run, gates.

## Boundaries

No cancel key and no needs-you or agent-reports block (next Task). No writes of any kind.

## Technical Verification

Focused tests during iteration; handoff and Full Check gates per AGENTS.md's Verification Policy. Evaluated later under `agent-skills/references/check-method.md`.

## Technical Evidence

Pending execution.

## Drift Notes

Record architecture deltas and reconcile through the planner before Check.
