---
id: T-057
title: Run health tools safely and collect their results
objective: O-028
status: planned
depends_on: [{task: T-055, requires: clear}]
owner_validation: {required: false}
planned_by: {role: planner, session: planning-o028-20261001}
---

# Run health tools safely and collect their results

## Outcome

A `Collect` service runs the configured analysis tools one at a time without a shell, reads existing test and coverage reports without running tests, keeps every failure mode distinct, hands report bytes to a per-provider reader seam, and saves one immutable snapshot with the given origin.

## User Check

None beyond the Full Objective Check; collection is invoked by O-030 and O-031.

## Done When

- `Collect(ctx, CollectRequest{Root, Origin, Config, Readers, Runner, Clock})` processes configured instances sequentially in config order and returns one `CapabilityResult` per instance plus the saved snapshot ID. A capability with no configured instance yields `not_configured`.
- **Executed tools:** run `Executable` with `Args` in the project root via a `ToolRunner` interface (argument vector, never a shell). A literal `{report}` argument is replaced with a fresh temporary file path outside the repository; without it, stdout is the report. Stdout and stderr are capped (report input at 32 MiB, stderr kept only as a bounded, sanitized reason). The per-instance effective timeout applies. Cancellation stops the process and its children (process group on Unix, `cmd.Cancel` + `WaitDelay` everywhere) and leaves no temporary files.
- **Report-only providers:** read the configured report path (bounded, regular file, inside the project). Missing is `absent`. The report is `stale` when older than the newest relevant input in that instance's scope (from `ObserveRepository`), otherwise `fresh`. No test command is ever run.
- **Distinct outcomes:** executable not found or not runnable is `unavailable`; non-zero exit that the provider does not define as "findings present" is `failed`; deadline is `timed_out`; context cancel is `cancelled` and stops remaining instances as `cancelled`; a reader error is `failed` with a named reason; over-cap input is `partial` with a truncation reason; a provider with no registered reader is `unsupported`. None of these carries a value.
- **Reader seam:** `Reader` interface keyed by `ProviderKey` turns report bytes plus scope into a measured result; O-029 registers the real readers. O-028 ships the interface and fakes only.
- **Snapshot:** builds results with config digest, scope, exclusions, instance name, and repository identity; classifies with the existing `Assess`/`Overall`; saves through `Store.SaveSnapshot`. One instance's failure never removes or alters another instance's result. Collection never prunes.
- Required instances are reported as such in the result so O-030 can apply policy; `Collect` itself blocks nothing.

## Context Files

`.savepoint/objectives/O-028-code-health-provider-execution-and-scaffolding/Objective.md`; `internal/codehealth/model.go`; `internal/codehealth/config.go`; `internal/codehealth/snapshot.go`; `internal/codehealth/classification.go`; `internal/codehealth/storage.go`; `internal/codehealth/repository.go`; new `internal/codehealth/runner.go`, `internal/codehealth/runner_unix.go`, `internal/codehealth/runner_windows.go`, `internal/codehealth/collect.go`, `internal/codehealth/collect_test.go`, `internal/codehealth/runner_test.go`.

## Design References

Design sections 1, 9, 11, and 12; O-026 Module Boundary "Result and error states"; O-028 Confirmed Design Decisions (2026-10-01).

## Guardrails

FS-01, FS-05, FS-06, ARCH-03, ARCH-04, CFG-02, CFG-03, DEP-01, STYLE-06, TEST-01, TEST-02, TEST-04, TEST-08.

## Implementation Plan

1. Confirm the first Task's instance `Name`, `Required`, and `EffectiveTimeout` exist; return REPLAN REQUIRED if not.
2. Add `ToolRunner` and its `exec` implementation with capped buffers, pinned environment additions (no prompts), timeout, and platform process-group cleanup. Reuse the existing `cappedBuffer`.
3. Add the `Reader` interface and registry type; no production readers.
4. Implement `Collect`: per-instance dispatch (executed vs report-only), freshness via `ObserveRepository` scope and file mtimes, outcome mapping, result assembly, `Assess`/`Overall`, and `SaveSnapshot`.
5. Tests use a helper test binary (re-exec of the test executable) to simulate success, non-zero exit, hang, huge output, child process, and a `{report}` writer; fake readers cover success, malformed, and partial. Cover cancellation mid-sequence and verify no temp files remain.

## Boundaries

No real report parsing, no CLI, no discovery, no Check or TUI wiring, no pruning, no network in tests, no running of project test suites.

## Technical Verification

Focused runs while iterating. Process handling is platform-sensitive: fresh `make test-full` at handoff, with the Windows cancellation path covered by CI.

## Technical Evidence

Pending execution.

## Drift Notes

AGENTS.md Codebase Map for `internal/codehealth` must now say it runs configured tools and saves snapshots; reconcile at the Full Objective Check.
