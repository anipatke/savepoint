---
id: T-132
title: Start an agent run from the focused Task with s
objective: O-046
status: planned
depends_on: [{task: T-131, requires: clear}]
owner_validation: {required: true}
planned_by: {role: planner, session: planning-o046-20261011}
---

# Start an agent run from the focused Task with s

## Outcome

With the feature on, pressing `s` on a Task whose Next verb is Start, Build, or Test and whose gate allows it starts the run subcommand detached and shows a status line; the board holds no child handle, and closing it does not end the run. A Blocked Task, a non-Task focus, or a second start while another run is `running` is refused with a plain reason.

## User Check

Owner validation required. With `agent_launcher.enabled: true` locally, press `s` on a planned Task, see the status line and a new `.savepoint/runs/<id>/run.yml` within a second or two, quit the board with `q`, confirm the run continues (events grow), and reopen the board. Press `s` again while it runs and read the refusal.

## Done When

1. A new `internal/run` helper (platform files beside `proc_unix.go`/`proc_windows.go`) starts `os.Executable()` with `run --task T-### <project root>` detached: `Setsid` on Unix, `DETACHED_PROCESS|CREATE_NEW_PROCESS_GROUP` on Windows, stdin from the null device, stdout/stderr to `.savepoint/runs/launch.log` (replaced per start), then releases the process. No shell (CFG-02, FS-05).
2. `s` is bound, hinted in the footer, and listed in help only when the flag is on and the focus is a Task whose Next verb is Start/Build/Test and whose gate (the one `actionsForRecord` uses) allows it. T-131's flag-off golden tests still pass byte for byte.
3. The board refuses with a status-line reason when any run reads `running` and is not inferred lost by `run.RunnerLost`.
4. The spawn is a `tea.Cmd`; `Update` performs no filesystem work (ARCH-02). After a spawn, if no new run directory appears within a few seconds, the board shows the last line of `launch.log` as the reason; a run that appears clears the pending state.
5. `runs` joins the watched directories only when the flag is on; the watcher creates nothing when `.savepoint/runs/` is absent.
6. Tests: gating matrix (flag off/on, Objective focus, Blocked Task, Check verb, running run, lost run), a spawn test with a fake executable proving the parent can exit while the child continues on Unix and Windows, and the launch-log refusal path (TEST-02).
7. `make build && make test-fast` and fresh `make test-full` pass, with native Windows CI evidence named, because this Task touches process creation. Per-criterion evidence follows AGENTS.md's Verification Policy; owner validation is recorded after the User Check.

## Context Files

`internal/run/launch.go`, `internal/run/run.go`, `internal/run/stop.go`, `internal/run/proc_unix.go`, `internal/run/proc_windows.go`, `internal/run/lost.go`, `internal/run/record.go`, `internal/board/v2/actions.go`, `internal/board/v2/update.go`, `internal/board/v2/help.go`, `internal/board/v2/view.go`, `internal/board/v2/io.go`, `internal/board/v2/watch.go`, `internal/board/v2/watch_test.go`, `internal/board/v2/model.go`, `.savepoint/objectives/O-046-start-and-watch-a-run-from-the-board-opt-in/Objective.md`.

## Design References

Design §8 (keybindings, board persistence and refresh), §12 (Windows support). O-046 Confirmed Design (detached start, board-side gating).

## Guardrails

ARCH-02, ARCH-03, ARCH-04, CFG-02, CFG-03, FS-05, TEST-01, TEST-02, TEST-04, TEST-08.

## Implementation Plan

1. Confirm `internal/run` still exposes `Record`, `RunnerLost`, and the run directory layout; return REPLAN REQUIRED if not.
2. Add the detached spawn helper with Unix and Windows files and tests.
3. Add a read-only "is a run active" helper over `.savepoint/runs/` in `internal/run`.
4. Wire `s` through the board's action and hint paths, gated on the flag.
5. Add the pending-start status and launch-log refusal.
6. Add `runs` to the watcher when on; gates; push for Windows CI.

## Boundaries

No Runs overlay, activity view, or cancel (later Tasks). No change to `Launch`'s own checks, `resume`, gates, or Task records.

## Technical Verification

Focused tests during iteration; handoff and Full Check gates per AGENTS.md's Verification Policy, with the native Windows CI job as the platform evidence the Full Check needs, produced by the executor's push. Evaluated later under `agent-skills/references/check-method.md`.

## Technical Evidence

Pending execution.

## Drift Notes

Record architecture deltas and reconcile through the planner before Check.
