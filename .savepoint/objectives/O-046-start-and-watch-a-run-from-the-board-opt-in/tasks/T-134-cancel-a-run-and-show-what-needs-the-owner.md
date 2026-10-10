---
id: T-134
title: Cancel a run and show what needs the owner
objective: O-046
status: planned
depends_on: [{task: T-133, requires: clear}]
owner_validation: {required: true}
planned_by: {role: planner, session: planning-o046-20261011}
---

# Cancel a run and show what needs the owner

## Outcome

With the feature on, `c` in the Runs overlay or activity view asks a running run to stop by writing its cancel file, and the row reaches `cancelled` from the runner's own record. A finished run shows the agent's final message under an "agent reports" label with no board verdict; a run with denials reads "needs you" in red and names each refused action; a lost runner reads `runner_lost`.

## User Check

Owner validation required. Start a run, press `c`, and confirm it reads `cancelled` within 15 seconds with no `claude` process left. Run once with an empty allow-list and read the "needs you" block. Read a finished run's "agent reports" text and confirm the board gives no verdict.

## Done When

1. `c` writes `<run dir>/cancel` through a `tea.Cmd` and touches no process; it is offered only for a run that reads `running` and is not lost, and refused with a reason otherwise. The board never rewrites `run.yml`.
2. A finished run shows `agent_reported` under an "agent reports" label in purple; nothing on the board calls the work passed, failed, or done.
3. Any denial makes the row and view read "needs you" in red, with one line per denial naming tool and target; the overview header counts runs that need the owner (review §4.1, §4.3), without the re-run, takeover, or dismiss keys.
4. States `failed`, `cancelled`, `timed_out`, and inferred `runner_lost` render distinctly and plainly; text from records has control sequences removed (ARCH-05).
5. Tests: cancel-file write and refusal paths, a fake runner reaching `cancelled`, denial and agent-reports rendering, lost-runner rendering. T-131's flag-off golden tests still pass.
6. `make build && make test-fast` and fresh `make test-full` pass, with native Windows CI evidence named for cancel. Per-criterion evidence follows AGENTS.md's Verification Policy; owner validation is recorded after the User Check.

## Context Files

`internal/run/stop.go`, `internal/run/record.go`, `internal/run/lost.go`, `internal/board/v2/io.go`, `internal/board/v2/update.go`, `internal/board/v2/help.go`, `internal/board/v2/view.go`, `.savepoint/visual-identity.md`, `docs/reviews/agent-orchestration-pass-1.md`, `.savepoint/objectives/O-046-start-and-watch-a-run-from-the-board-opt-in/Objective.md`, plus the Runs overlay files T-133 created (named in its Technical Evidence).

## Design References

Design §8 (TUI, keybindings). Review §4.1, §4.3. O-046 Confirmed Design (cancel, presentation).

## Guardrails

ARCH-02, ARCH-05, CFG-02, CFG-03, FS-05, STYLE-07, TEST-01, TEST-02, TEST-04, TEST-08.

## Implementation Plan

1. Confirm the cancel file name in `internal/run` and expose it rather than duplicating the string (STYLE-07).
2. Add the cancel command and key handling.
3. Add the needs-you, agent-reports, and end-state presentation.
4. Tests, golden re-run, gates, push for Windows CI.

## Boundaries

No re-run with a wider allow-list, takeover, dismiss, or permission editing. No Task, Objective, Check, or router write.

## Technical Verification

Focused tests during iteration; handoff and Full Check gates per AGENTS.md's Verification Policy, with the native Windows CI job as platform evidence for cancel, produced by the executor's push. Evaluated later under `agent-skills/references/check-method.md`.

## Technical Evidence

Pending execution.

## Drift Notes

Record architecture deltas and reconcile through the planner before Check.
