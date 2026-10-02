---
id: T-091
title: Verify provider processes stop safely
objective: O-032
status: done
depends_on: [{task: T-085, requires: clear}]
complexity_tier: high
complexity_reason: Unix and Windows process-tree cleanup and bounded reports are platform-sensitive.
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: planning-o032-20261002-owner-confirmed}
check_waiver:
    task: T-091
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-02T04:16:51Z"
---

# Verify provider processes stop safely

## Outcome

Provider execution has repeatable safety evidence for cancellation, timeout, cleanup, argument/path handling and bounded sanitized failures on supported platforms.

## User Check

Verify a cancelled helper process tree ends and later collection can run normally.

## Done When

- Use controlled local helper processes to prove child/descendant termination, timeout/cancel ordering, pipe release and repeat after failure; retain distinct unavailable/failed/timed_out outcomes.
- Probe stdout report limit and stderr cap at below/exact/above boundaries, malformed/oversized file reports, no-shell argument handling and safe diagnostics without source/secrets leakage.
- Verify independent instances remain independent and cancellation before save writes no snapshot. Report failure after save remains separate from collection failure.
- Use scanner decision as evidence, not permission for an unknown setup/install change; only repair demonstrated execution-contract defects. Provider-specific unknown changes return to design.
- Record Linux full evidence and the existing native Windows CI verification path/results; no Windows skip for a supported real-world case without justified policy. Normal tests run without provider installations or network.

## Context Files

`.savepoint/objectives/O-032-code-health-release-hardening/tasks/T-085-find-why-the-vulnerability-scan-fails.md` (dependency evidence); `internal/codehealth/runner.go`; `internal/codehealth/runner_unix.go`; `internal/codehealth/runner_windows.go`; `internal/codehealth/runner_test.go`; `internal/codehealth/runner_alive_unix_test.go`; `internal/codehealth/runner_alive_windows_test.go`; `internal/codehealth/collect.go`; `internal/codehealth/collect_test.go`; `internal/codehealth/collect_repair_test.go`.

## Design References

O-032 Confirmed Release-Hardening Design and Owner-Requested Health Analysis; Design sections 1, 6, 7, 9, 12 and 13 as applicable.

## Guardrails

FS-01, FS-04, FS-05, FS-06, ARCH-01, ARCH-02, ARCH-03, CFG-01, CFG-02, TEST-01, TEST-02, TEST-03, TEST-04, TEST-05, TEST-08, TEST-09, STYLE-01, STYLE-02, STYLE-07, CFG-03, REL-01, REL-02, REL-03.

## Implementation Plan

1. Reconcile scanner research with execution invariants.
2. Inventory existing runner and collector failure tests.
3. Add boundary and subprocess-tree scenarios; repair only bounded reproduced safety defects.
4. Run relevant local tests and fresh full gate; record native platform limits honestly.

## Boundaries

No provider installation, shell execution, network-dependent normal tests, new capability, thresholds or provider sandbox promise.

## Technical Verification

Fresh `make test-full` for this platform/migration/distribution-sensitive Task; focused `make test-focused TEST=... PKGS=...` only during iteration. Record named per-criterion evidence and limitations. The later independent Check applies `agent-skills/references/check-method.md`; the mandatory Full Objective Check runs its fresh full gate and official Code Health collection. This executor never writes a Check or runs savepoint health setup/report/check.

## Technical Evidence

Run 2026-10-02 04:15-04:16Z (UTC), go1.26.2 linux/amd64 (WSL2). No provider installed, no network, no savepoint health command run, no scanner/setup change.

### Inventory (plan step 2)

Existing coverage kept as is: stdout truncation 2 MiB over the cap; stderr sanitized/bounded; missing and not-runnable executable = unavailable; no-shell `$(echo hi)` / `a;b`; deadline, cancel, already-cancelled, grandchild killed on cancel; distinct unavailable/failed/timed_out/partial outcomes; cancellation before/mid/after tools writes no snapshot; report-write failure keeps the snapshot; oversized report file +10 bytes; missing/stale/directory/escaping-symlink report files; reader error contained; real-process temp-file cleanup.

### Added

`internal/codehealth/runner_safety_test.go` (new), helper modes in `runner_test.go` (`stdout-size`, `stderr-tail`, `args`, `orphan`).

### Production repair (demonstrated defect)

The stderr boundary probe failed: with stderr above `maxStderrBytes` (64 KiB), `limitedBuffer` kept the head, so the final error that `sanitizeTail` (T-085) is meant to keep was lost (`…pppEN` / no `END`). Fixed in `runner.go` with a small `tailBuffer` used only for stderr; stdout still keeps its start. Tests failed before the repair and pass after. No other production change.

### Per-criterion outcomes

1. Process tree / ordering / pipes / repeat: `TestExecRunnerStopsOnDeadlineAndCancel/children_die_with_the_tool` (existing, cancel), `TestExecRunnerStopsTheWholeTreeOnDeadlineAndRunsAgain` (descendant gone after deadline, then a normal run succeeds), `TestCollectRealProcessTimeoutEndsDescendantsAndLaterInstanceRuns` (timed_out, descendant gone, next instance available), `TestExecRunnerBoundsPipesHeldByAnOrphan` (finished tool whose child holds the pipes returns a result in about killWait=2s), `TestCollectCancelWinsOverTheInstanceDeadline` (cancel beats timed_out, nothing saved). unavailable/failed/timed_out stay distinct (`TestCollectExecutedOutcomesStayDistinct`). Met on Linux.
2. Boundaries: `TestExecRunnerOutputCapsAtTheBoundary` stdout limit-1/limit/limit+1 (Truncated only above) and stderr cap-1/cap/cap+1/8x cap (end kept, <= MaxReasonLen, tool never blocked); `TestCollectReportFileSizeBoundary` file report limit-1/limit available, limit+1 partial and never read; malformed reader error and oversized reasons bounded (`TestCollectFailureReasonsLeakNeitherReportNorEnvironment`); `TestExecRunnerPassesArgumentsAndDirectoryUnchanged` (nine shell-meta/empty/newline args unchanged, directory with spaces); failing-tool stdout (source) and environment secret never in the saved reason. Met.
3. Independence and save ordering: existing `TestCollectCancelledSavesNothing`, `TestCollectReportWriteFailureKeepsTheSnapshot`, `TestCollectRunsInstancesSequentially...`, plus the neighbour-after-timeout check above. Met.
4. Scanner decision used as evidence only: the only repair is the demonstrated stderr-tail defect inside the execution contract; no setup/install/provider change. Met.
5. Evidence: Linux `make test-full` passed (below). Native Windows: the existing `windows-tests` job in `.github/workflows/ci.yml` (`go run ./internal/buildtool test -json -count=1 ./...`) runs these tests natively; `GOOS=windows go vet ./internal/codehealth` passes locally. Not run natively here. No Windows skip added; tests need no provider or network. Met with the Windows limitation below.

### Commands

- `go test ./internal/codehealth -count=5 -run 'TestExecRunner|TestCollect'` -> ok (35.8s)
- `gofmt -l`, `go vet ./internal/codehealth`, `GOOS=windows go vet ./internal/codehealth` -> clean
- `make test-full` (test + build-all) -> exit 0, no FAIL lines (about 21s)

### Files

Read: the Task, T-085 (dependency evidence), router.md, Objective.md, runner.go, runner_unix.go, runner_windows.go, runner_test.go, runner_alive_unix_test.go, runner_alive_windows_test.go, collect.go, collect_test.go, collect_repair_test.go. Extra reads: AGENTS.md, the skill, `.github/workflows/ci.yml` and `Makefile` (to name the gate and the Windows CI path), `history.go` boundText (sanitization bound). Changed: this Task, `internal/codehealth/runner.go`, `runner_test.go`, `runner_safety_test.go`.

### Limitations

- Windows behavior (taskkill tree walk, orphan pipe handling, argument quoting of empty/newline arguments) is verified only by cross-vet here; the native result is the next `windows-tests` CI run, not yet seen.
- A descendant that leaves the Unix process group (setsid) or detaches on Windows is not killed by cancel; only the 2s pipe bound applies. Not reproduced as a product case, left as is.
- The orphan test intentionally leaves a short-lived grandchild until the test kills it.

## Drift Notes

Responsibility change: `ExecRunner` stderr capture now keeps the last 64 KiB (`tailBuffer`) instead of the first, completing T-085's tail-preserving repair. Planner to reconcile.

Record implemented responsibility or interface changes for planner reconciliation before the mandatory Full Objective Check. If evidence requires a material unknown repair, return REPLAN REQUIRED rather than inventing another scope.
