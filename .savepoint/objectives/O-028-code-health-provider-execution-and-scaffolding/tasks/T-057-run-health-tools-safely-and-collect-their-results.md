---
id: T-057
title: Run health tools safely and collect their results
objective: O-028
status: done
depends_on: [{task: T-055, requires: clear}]
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: planning-o028-20261001}
check_waiver:
    task: T-057
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-01T07:56:40Z"
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

Executor: Claude Sonnet 5.5, 2026-10-01. Status `in_progress`, stage `audit` (ready for a Check; not a pass).

**Per criterion**

- `Collect` sequencing, config order, `not_configured`, saved ID, origin/retention: `TestCollectRunsInstancesSequentiallyInConfigOrderAndSaves`, `TestCollectOfficialOriginIsPermanent` (fake runner fails the test if two tools overlap).
- Executed tools without a shell, `{report}` temp file outside the repository, stdout as report, 32 MiB cap, bounded sanitized stderr, effective timeout, cancel and process group, no temp files left: `TestExecRunner*`, `TestCollectWithRealProcessesAndReportFiles`, `TestCollectCancelledRealProcessLeavesNoTemporaryFiles` (re-exec of the test binary as the tool: success, exit code, hang, huge output, child process, `{report}` writer). Windows path: `GOOS=windows go vet` and `make build-all` pass; the runtime cancel path there is covered by CI only.
- Report-only providers: `TestCollectReportOnlyNeverRunsATool` (absent, fresh, stale, directory, symlink escape; no tool run), `TestCollectReportFreshnessIgnoresTheReportItself`, `TestCollectOversizedReportFileIsPartial`.
- Distinct outcomes (unavailable, failed, OSV exit 1 as findings, timed_out, cancelled with remaining instances cancelled, reader error, partial on over-cap, unsupported with no tool run): `TestCollectExecutedOutcomesStayDistinct`, `TestCollectReaderProblemsAreContained`, `TestCollectCancellationStopsRemainingInstances`.
- Reader seam: `Reader`/`Readers`/`Reading`/`ReportInput`; fakes only, no production readers.
- Snapshot: config digest, scope, exclusions, name, repository identity, `Assess`/`Overall`, `SaveSnapshot`; one failure never alters a neighbour; nothing pruned (`TestCollectNeverPrunesAndAccumulatesHistory`). `Required` is carried on each `Collected`.

**Commands**: `make build && make test-fast` passed; fresh `make test-full` passed (exit 0, go1.26.2 linux/amd64, includes `build-all`). `GOOS=windows go vet ./internal/codehealth/` passed.

**Files read**: all Context Files that exist, plus `AGENTS.md`, `agent-skills/savepoint-task/SKILL.md`, `.savepoint/router.md`, `internal/codehealth/{history.go,identity.go,primitives.go}` and `repository_test.go` helpers (extra reads: needed for `boundText`, `Digest`, `validate*`, and the test helpers `newRepo`/`write`).

**Files changed**: new `internal/codehealth/{runner.go,runner_unix.go,runner_windows.go,collect.go,runner_test.go,collect_test.go,runner_alive_unix_test.go,runner_alive_windows_test.go}`; edited `internal/codehealth/repository.go` (extra change, below); this Task file.

**Deviations and limitations**

- `repository.go` was changed (in Context Files): `Observation` gains `Newest` (latest mtime of relevant inputs) via a new `fingerprintNewest`; `fingerprintInputs` keeps its signature. Needed because freshness compares report mtime to the newest relevant input and `ObserveRepository` did not expose it.
- `CollectRequest` has an extra optional `Git CommandRunner` (default `GitRunner{}`) for the repository observation; the plan's field list did not name one.
- The existing `cappedBuffer` was not reused: it fails the command on overflow, but the plan needs over-cap output to become `partial`. A new `limitedBuffer` drops the excess and records truncation.
- Over-cap partial results record `unknown` for provider version and measurement definition, since no reader ran.
- Damaged snapshot history fails `Collect` before any tool runs; it is not skipped.
- Freshness is `unknown` when the per-instance repository observation fails for a reason other than cancellation.
- Windows process-tree kill is not available in the standard library; there `cmd.Cancel` kills the tool and `WaitDelay` bounds surviving children. Not run on Windows locally.
- The 32-result snapshot bound can be exceeded when many instances plus `not_configured` entries are configured; `SaveSnapshot` then refuses the whole snapshot.
- No Task Check requested or waived; none recorded.

## Drift Notes

AGENTS.md Codebase Map for `internal/codehealth` must now say it runs configured tools and saves snapshots; reconcile at the Full Objective Check.
