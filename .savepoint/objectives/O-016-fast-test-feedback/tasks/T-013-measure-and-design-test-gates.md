---
id: T-013
title: Find the slow test work
objective: O-016
status: done
depends_on: []
owner_validation:
    required: true
    accepted_check: ""
planned_by: {role: planner, session: o016-design-20260923}
complexity_tier: spike
complexity_reason: Gate selection and safe reuse depend on measured test cost and shared-state analysis.
check_waiver:
    task: T-013
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-09-22T22:44:25Z"
---

# T-013: Find the slow test work

## Outcome

A measured baseline and a written gate-selection decision identify the costly tests, the exact fast/full test sets, and the isolation limits for optimization.

## User Check

Review the proposed gate membership and confirm that the full gate still runs every existing test. Confirm the recorded machine and timing method are useful for tracking the stated targets.

## Done When

- Cold and warm package and individual-test timings identify the main cost within internal/migrate, with command, Go version, OS, cache state, and repeated sample results recorded.
- An inventory maps every test to full coverage and names the proposed fast subset. No test is omitted from full coverage.
- The decision names repeated fixture work that can share read-only preparation, mutating scenarios that need fresh directories, and process-global hooks that prevent parallel execution.
- Fast/full command design, timing output, failure behavior, and the baseline against which improvements will be compared are documented for T-014-T-016.
- The current configured quality gates pass, and each criterion has evidence.

## Context Files

`Makefile`, `.github/workflows/ci.yml`, `internal/migrate/apply_test.go`, `internal/migrate/end_to_end_test.go`, `internal/migrate/fixture_test.go`, `internal/migrate/operation_test.go`, `internal/migrate/cutover_test.go`, `.savepoint/Guardrails.md`, `.savepoint/objectives/O-016-fast-test-feedback/Objective.md`.

## Design References

Design sections 1, 12, and 13; O-016 Confirmed Verification Policy and Architectural Considerations.

## Guardrails

TEST-01..04, TEST-07..08, CFG-02, STYLE-05, STYLE-07.

## Implementation Plan

1. Record the host/toolchain and run repeatable cold and warm timing samples without editing production code.
2. Attribute slow package time to named tests and fixture operations.
3. Inspect shared hooks and mutation boundaries for safe reuse or parallelism.
4. Produce the gate-selection and optimization decision as Task evidence; flag any target that cannot be met without changing O-016's boundaries.
5. Run the currently configured quality gates before handoff.

## Boundaries

Measurement and a decision deliverable only. Do not remove tests or change the verification policy in this Task.

## Technical Verification

Repeat timing commands; verify that every existing test is included in the proposed full selection; `make build && make test`.

## Technical Evidence

### Start and scope

- Start allowed: depends_on is empty, and O-016 records owner confirmation of the T-013-T-016 sequence and verification policy. No prerequisite or design blocker was found.
- T-013 remained measurement and gate-selection research. No production or test source was changed. Only this Task record was edited.

### Machine and timing method

- Recorded 2026-09-23 on Windows 11 Home build 26200, using the repository on WSL2 Ubuntu 24.04.4 LTS, kernel 6.18.33.2-microsoft-standard-WSL2. Test runtime: Go 1.26.2 linux/amd64; AMD Ryzen 7 7800X3D, 8 cores / 16 logical CPUs; WSL reports 16 CPUs.
- Cold means an isolated, initially empty Go build cache. Every timing command used -count=1 to bypass Go test-result caching. Warm samples reused that build cache and still executed all selected tests. The OS page cache was not flushed.
- /usr/bin/time measured go-command wall time. go test -json Package/Elapsed and Test/Elapsed events supplied test-only package and individual-test durations. JSON logs and cache were kept under /tmp, outside the repository.

| Sample | Cache state | internal/migrate test time | go-command wall time | TestEndToEnd_temporaryRepositoryCopyMigratesWithReleaseAccountability | TestApply_recoversAtEveryPublishBoundaryWithoutOverwritingUserEdits | TestEndToEnd_goldenIsReproducible |
|---|---|---:|---:|---:|---:|---:|
| Cold | fresh isolated GOCACHE | 114.450 s | 129.12 s | 89.080 s | 14.740 s | 1.080 s |
| Warm 1 | reused GOCACHE, uncached tests | 120.869 s | 135.82 s | 95.630 s | 14.560 s | 1.070 s |
| Warm 2 | reused GOCACHE, uncached tests | 116.895 s | 131.82 s | 92.320 s | 14.400 s | 1.010 s |

The repository-copy integration and publish-boundary recovery tests account for about 91% of the two warm package durations. The first copies the repository (excluding .git), restores the archived V1 boundary, then exercises preview, apply, and a no-op second apply. The second recovers at every planned publish boundary and separately proves a user edit is preserved.

### Test inventory and gate membership

The package counts below combine uncached go test JSON output for the fast selection with the full internal/migrate JSON inventory. A targeted go test -list query confirmed the three full-only names occur only in internal/migrate; make test enumerated every package in the full set. Counts are top-level test entries; nested cases follow their parent. Fast excludes only the three named migration scenarios below.

| Go package | Full entries | Fast entries |
|---|---:|---:|
| github.com/opencode/savepoint | 50 | 50 |
| github.com/opencode/savepoint/cmd | 47 | 47 |
| github.com/opencode/savepoint/internal/board | 452 | 452 |
| github.com/opencode/savepoint/internal/board/v2 | 191 | 191 |
| github.com/opencode/savepoint/internal/buildtool | 16 | 16 |
| github.com/opencode/savepoint/internal/data | 595 | 595 |
| github.com/opencode/savepoint/internal/doctor | 176 | 176 |
| github.com/opencode/savepoint/internal/init | 269 | 269 |
| github.com/opencode/savepoint/internal/migrate | 169 | 166 |
| github.com/opencode/savepoint/internal/resume | 27 | 27 |
| github.com/opencode/savepoint/internal/styles | 8 | 8 |
| github.com/opencode/savepoint/internal/testutil | 0 (no test files) | 0 |
| Total | 2000 | 1997 |

Proposed fast selection command, measured twice with -count=1 and JSON timing enabled:

    go test -json -count=1 -skip '^(TestEndToEnd_temporaryRepositoryCopyMigratesWithReleaseAccountability|TestApply_recoversAtEveryPublishBoundaryWithoutOverwritingUserEdits|TestEndToEnd_goldenIsReproducible)$' ./...

The full-only tests are TestEndToEnd_temporaryRepositoryCopyMigratesWithReleaseAccountability, TestApply_recoversAtEveryPublishBoundaryWithoutOverwritingUserEdits, and TestEndToEnd_goldenIsReproducible. The last full-only test repeats two independently prepared migrations to verify golden reproducibility; TestEndToEnd_goldenConvertedOutput remains in fast. The exact anchored skip pattern omits no other top-level test entry. The full command remains make test, which runs go test ./... with no exclusions, so every existing test remains in full coverage.

| Fast sample | Exit | Wall time | internal/migrate test time |
|---|---:|---:|---:|
| Initial candidate (only repository-copy and recovery excluded) | 0 | 15.67 s | 11.413 s |
| Final selection sample 1 (three exclusions) | 0 | 10.98 s | 10.342 s |
| Final selection sample 2 (three exclusions) | 0 | 10.91 s | 10.267 s |

The initial candidate missed the 15-second target. Adding the 1.01-second golden-reproducibility scenario to the full-only set produced two passing warm-cache samples below target. Keep the budget as measured guidance with attribution, not an in-test timeout.

### Fixture, mutation, and parallelism decision

- Frozen migration fixtures are source evidence. apply_test.go copies them to a fresh t.TempDir through copyFixtureProject; end_to_end_test.go uses copyFixture/migrateFixture; operation_test.go creates each operation under a fresh t.TempDir. Mutating Apply, recovery, and cutover scenarios need separate writable directories so order and previous mutations cannot affect a later case.
- TestApply_recoversAtEveryPublishBoundaryWithoutOverwritingUserEdits reuses a baseline snapshot and one probe plan to enumerate boundaries, then makes a new v1-history copy and plan for each injected failure. The planned optimization may share read-only fixture manifests, source hashes, and boundary metadata, but must preserve a fresh mutable project for each injected boundary.
- TestEndToEnd_temporaryRepositoryCopyMigratesWithReleaseAccountability is intentionally repository-scale and keeps preview, apply, and no-op apply ordered against one private copy. It is full-only in the proposed fast gate. Within fixture-based end-to-end coverage, read-only assertions about one converted result (references, clearance, accountable destinations, and archive hashes) are candidates to consume one independently prepared result with separate named assertions; never share a directory between mutating cases.
- The scoped migration test files contain no t.Parallel calls. The recovery test changes the package-global afterPublishWriteHook and afterPublishRemovalHook and restores them with cleanup. Keep hook-using tests serial until hooks are per-run dependencies. Go's existing package-level concurrency remains safe because test packages run in separate processes and test mutations use temporary directories.

### Gate design, timing output, and failure behavior

- Keep make test as the full Go test command (go test ./...). Keep make build as the required native build. CI's existing make ci continues to include the full test, build, distribution, and package checks.
- Add a future make test-fast target using the exact three-test skip selection above. Both proposed commands use go test -json for package and individual-test timing; the Make/CI wrapper should summarize the slowest packages/tests while preserving the Go command's nonzero exit status. Any included test failure fails the gate. Timing targets remain reported budgets, not flaky timeout assertions.
- Current .github/workflows/ci.yml has one Ubuntu job running make ci. The Makefile's ci target includes test, build, dist, and package-check, but the workflow does not run the test suite natively on Windows. The confirmed O-016 policy requires a Windows test job in CI; T-014-T-016 must add it while retaining the Linux full gate and cross-build coverage.

### Acceptance evidence and commands

1. Timing criterion: met. Three cold/warm uncached package samples and individual durations identify the two dominant tests; machine, Go version, OS, cache state, and method are recorded above.
2. Inventory criterion: met. Package-by-package full and fast counts map all 2000 top-level entries; the full command has no skip and covers every package/test. The fast command excludes only the three named parents and their nested cases.
3. Isolation criterion: met. Reusable read-only metadata/result assertions, fresh mutable directories, and process-global hook restrictions are recorded above.
4. Gate-design criterion: met. Exact fast/full commands, JSON timing reporting, nonzero failure behavior, and measured baseline for T-014-T-016 are recorded above.
5. Current quality gate: passed. On 2026-09-23, make build && make test exited 0 on the final retry. That run executed internal/data (0.543 s) and internal/migrate (118.868 s); unchanged packages with successful same-input results were reported cached. A fresh full run with the temporary npm bridge also passed internal/init and all other packages except a transient internal/data test-log write error; internal/data then passed alone with -count=1 and passed again in the successful full retry. Since that gate run, only this Task record status/evidence metadata changed; Go code, tests, fixtures, Makefile, and workflow are unchanged. Under O-016's confirmed metadata-only correction policy, the named full-gate result remains current.

Commands and run notes:

- Three samples: go test -json -count=1 ./internal/migrate, with one fresh GOCACHE and two reuses of that cache.
- Fast selection: the exact proposed -skip command above passed twice; the preliminary two-exclusion candidate passed at 15.67 s. go test -list with the three-name regex confirmed all full-only names are in internal/migrate only.
- Required quality gate: make build && make test passed. A direct Windows Go invocation from the WSL UNC path failed before running tests (Go could not lock go.mod). Native WSL Go was used thereafter.
- The npm integration initially failed because WSL interop launched Windows npm with an unsupported UNC temp working directory. A temporary /tmp wrapper used cmd pushd to map that same WSL temp directory to a Windows drive and invoked the real npm install; TestIntegration_InstallDependencies passed. No test source was modified. Pointing TMPDIR at /mnt/c was rejected because it caused expected Linux case-sensitive and permission tests to fail.
- A first uncached full run with the bridge hit the transient internal/data test-log write error; go test -count=1 ./internal/data then passed, and the final full-gate retry passed. The successful retry used Go's cache only for unchanged package inputs.
- go test -count=1 ./internal/init -run '^TestIntegration_InstallDependencies$' passed with the temporary npm bridge.
- go test -count=1 ./internal/data passed after the transient test-log error.
- The temporary bridge, JSON logs, and caches are under /tmp, not in the repository. make ci was inspected but not run; T-013's required local gate is make build && make test.

### Files read, changed, and extra reads

- Task Context Files read: Makefile; .github/workflows/ci.yml; internal/migrate/apply_test.go; internal/migrate/end_to_end_test.go; internal/migrate/fixture_test.go; internal/migrate/operation_test.go; internal/migrate/cutover_test.go; .savepoint/Guardrails.md; .savepoint/objectives/O-016-fast-test-feedback/Objective.md.
- Required routing reads: .savepoint/router.md and agent-skills/savepoint-task/SKILL.md.
- Extra read: .savepoint/Design.md sections 1, 12, and 13, because T-013 names them as Design References and they define build/test/CI boundaries.
- Extra read: internal/init/integration_test.go around TestIntegration_InstallDependencies, to diagnose the make-test WSL/npm failure and choose a real-npm bridge.
- Changed file: this T-013 record only. git status --short showed only this file modified; no production or test source changed.

### Owner decision

- Owner acceptance: the owner said “Happy with this analysis” and directed T-013 closure. The analysis and recorded evidence are accepted for T-013.
- Owner waiver: Task T-013; the optional independent Task Check is waived. Reason: the owner reviewed and accepted the analysis and directed closure without an independent Task Check. Actor: owner. Time: 2026-09-22 22:39:42 UTC.
- This waiver skips the optional Task Check only; it does not establish technical CLEAR. T-014 has a dependency on T-013 requiring CLEAR.
