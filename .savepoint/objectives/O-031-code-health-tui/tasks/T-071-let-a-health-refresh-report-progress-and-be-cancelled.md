---
id: T-071
title: Let a health refresh report progress and be cancelled
objective: O-031
status: done
depends_on: []
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: planning-o031-20261002}
complexity_tier: medium
complexity_reason: Small collector API change, but cancellation crosses process-group handling on Unix and Windows.
check_waiver:
    task: T-071
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-01T20:27:49Z"
---

# Let a health refresh report progress and be cancelled

## Outcome

A caller of Code Health collection can see which instance is running ("2 of 5: coverage") and can cancel it; a cancelled collection saves nothing, so a half-finished run never replaces the last real result.

## User Check

None beyond the Full Objective Check.

## Done When

- `CollectRequest` gains an optional progress callback called once before each configured instance with its position, total, capability, provider, and name. Nil means no reporting; existing callers compile and behave as before.
- When the context is cancelled before the snapshot is saved, `Collect` saves nothing and returns a named `ErrCollectionCancelled`; the store is unchanged (test compares snapshot directory listings before and after).
- `savepoint health check` reports a cancelled run as "cancelled; no snapshot was saved" and exits non-zero, consistent with its existing no-snapshot rule.
- A running tool is stopped on cancel on Unix and Windows using the existing runner; no orphaned process remains (existing runner alive tests extended or reused by name).
- Tests cover progress order and counts, cancel before the first instance, cancel mid-tool, cancel after the last instance but before save, and an uncancelled run unchanged.

## Context Files

`.savepoint/objectives/O-031-code-health-tui/Objective.md`; `internal/codehealth/collect.go`, `internal/codehealth/collect_test.go`, `internal/codehealth/errors.go`, `internal/codehealth/runner.go`, `internal/codehealth/runner_test.go`, `internal/codehealth/runner_alive_unix_test.go`, `internal/codehealth/runner_alive_windows_test.go`; `internal/healthcheck/healthcheck.go`, `internal/healthcheck/healthcheck_test.go`.

## Design References

Design sections 1 (Code Health readers), 6 (`health check`), 9, and 11; O-031 Confirmed Design Decisions (Refresh, Cancel).

## Guardrails

CFG-02, CFG-03, FS-06, ARCH-03, STYLE-06, TEST-01, TEST-02, TEST-04, TEST-08.

## Implementation Plan

1. Confirm `Collect` saves after all instances and that the runner already kills the process group on context cancel; return REPLAN REQUIRED if cancellation cannot reach a running tool.
2. Add the progress callback and call it in configuration order.
3. Check the context before saving; return the named error and save nothing.
4. Map the error in `healthcheck.Run`.
5. Add the tests above with fake runners; reuse the platform alive tests for real process stop.

## Boundaries

No board code, no new origin, no change to classification or snapshot schema.

## Technical Verification

Cancellation is platform-sensitive: fresh `make test-full` before handoff. Later evaluated per `agent-skills/references/check-method.md`.

## Technical Evidence

Commands: `make build && make test-full` (fresh, linux; cross-builds for darwin/windows compiled) passed; `go vet` clean on both packages.

Per criterion:
- Progress callback: `CollectRequest.Progress func(Progress)` called before each configured instance in order with position, total, capability, provider, name; nil is a no-op. `TestCollectReportsProgressInOrder`.
- Cancel saves nothing: `Collect` returns `ErrCollectionCancelled` when the context is done after observation or before save. `TestCollectCancelledSavesNothing` compares store file listings (before first instance, mid tool, after last instance before save) and checks an uncancelled run still saves. The old test expecting a snapshot of cancelled outcomes was replaced.
- `health check`: `healthcheck.Run` returns "health check: cancelled; no snapshot was saved: ..." (non-zero via returned error); `TestRun_cancelledContextSavesNothing` tightened. Interrupt already cancels via `signal.NotifyContext` in main.go.
- Tool stop on cancel: existing `ExecRunner` process-group kill (unix `Setpgid`, windows `Cancel`) reused unchanged; `TestExecRunnerStopsOnDeadlineAndCancel` (children die) and `TestCollectCancelledRealProcessLeavesNoTemporaryFiles` (now expects `ErrCollectionCancelled`).
- Tests listed all present.

Extra reads: `main.go` (health check wiring/signal handling) to confirm interrupt reaches the context. Limitation: Windows kill path not executed here (compiled only).

## Drift Notes

Changes official `health check` behavior on interrupt from "saves a snapshot with cancelled outcomes" to "saves nothing"; reconcile Design section 6 wording at the Objective Check.
