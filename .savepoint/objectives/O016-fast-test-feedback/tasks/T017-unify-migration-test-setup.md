---
id: T017
title: Unify migration test setup
objective: O016
status: done
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: o016-remediation-design-20260923}
check_waiver:
    task: T017
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-09-23T04:21:24Z"
---

# T017: Unify migration test setup

## Outcome

The migration package has one package setup path on every platform, while
retaining shared read-only fixture preparation and the Windows helper-process
tests.

## User Check

Review the focused Windows check for shared fixture preparation and both
environment-driven helper paths. The full Windows runtime suite is outside
this Task's scope.

## Done When

- The Windows test build has exactly one effective TestMain for
  internal/migrate.
- Shared fixture preparation and Windows environment-driven helper modes both
  run through the package setup without bypassing their existing cleanup and
  isolation rules.
- The focused Windows tests pass for shared read-only fixture results,
  destination-held-open helper processes, and interruption helper processes.
- The Windows CI job runs those focused migration tests rather than the full
  Windows repository suite.
- A fresh make test-full passes for Task handoff. Record toolchain and result.
- I029 remains open for independent Check evidence; this Task records repair
  evidence but does not close the Issue.

## Context Files

internal/migrate/end_to_end_test.go,
internal/migrate/replace_windows_test.go,
.github/workflows/ci.yml.

## Design References

.savepoint/Design.md sections 1 and 13; O016 Confirmed Verification Policy.

## Guardrails

CFG-02, TEST-05, TEST-08.

## Implementation Plan

1. Reconcile the unconditional shared-fixture TestMain with the Windows-only
   helper-process setup so only one package TestMain is compiled on Windows.
2. Dispatch Windows helper-process modes before shared fixture preparation so
   child invocations preserve their existing behavior without preparing unused
   fixtures.
3. Preserve shared read-only fixture protection, Windows private copies, and
   fresh copies for mutating and recovery scenarios.
4. Run the focused Windows setup/helper tests and the fresh Linux
   make test-full gate. Record the exact Windows Go version and named tests.

## Boundaries

No migration scenario is removed, skipped, or made order-dependent to resolve
the package setup conflict. The Windows CI job runs only the focused setup
check; the full Windows runtime suite and unrelated platform repairs are
deferred. Do not change T014 or T016 status.

## Technical Verification

Windows:
`go run ./internal/buildtool test -json -count=1 ./internal/migrate -run '^(TestEndToEnd_sharedMigratedResultsAreReadOnly|TestProbeDestinationHeldOpen|TestProbeInterruptionBeforeReplace)$'`;
fresh `make test-full`.

## Technical Evidence

**Owner scope update (2026-09-23 03:05 UTC):** The owner narrowed Windows runtime coverage to this focused setup check; the full Windows suite and unrelated platform repairs are deferred. The Linux `make test-full` gate remains required.

- Focused Windows check passed with Go `go1.26.2 windows/amd64`: `go run ./internal/buildtool test -json -count=1 ./internal/migrate -run '^(TestEndToEnd_sharedMigratedResultsAreReadOnly|TestProbeDestinationHeldOpen|TestProbeInterruptionBeforeReplace)$'`. It ran from a local NTFS copy (the WSL UNC path cannot run Go directly) and reported the package at 1.912s. This covers shared fixture preparation/private copies and the held-open and crash helper-process paths.
- Owner scope confirmation (2026-09-23 03:55 UTC): remove T018 and T019 from O016 because their only basis was the full Windows runtime suite, which is deferred. The focused T017 Windows check and Linux `make test-full` remain in scope; the Objective records this narrowing.

- Exploratory full Windows package run (now outside scope): `go1.26.2 windows/amd64`. Running `go test -count=1 ./internal/migrate` from the WSL UNC workspace failed before startup with `RLock ... go.mod: Incorrect function`. A local NTFS copy (excluding `.git`) ran for 138.044s; the duplicate `TestMain` setup error was gone. `TestEndToEnd_sharedMigratedResultsAreReadOnly` and `TestApply_recoversAtEveryPublishBoundaryWithoutOverwritingUserEdits` reported no failure. The inventory suite ran but `TestInventory_performsNoWriteOnFailure/case_collision` and `TestInventory_rejectsCaseCollision` got nil instead of `ErrCaseCollision`. `TestOperation_installFailureLeavesLiveFileIntact` allowed install into a read-only directory. Directory mtime assertions failed in `TestPreflightCutover_isReadOnlyAcrossV1V2AndRecovery` (V1 plan, V2 index, pending operation), `TestPreflightCutover_allowsAcceptedMultiReleaseCandidate`, and `TestPreflightCutover_allowsReleaseFreeV2`. Existing I030 covers the board watcher and guide-casing failures, not these migration-package failures.
- First full-gate attempt: `make test-full`, invoked from the WSL workspace on `go1.26.2 linux/amd64`, failed in `internal/init/TestIntegration_InstallDependencies`. WSL has no Linux Node/npm; `npm` resolves to Windows npm. The test launched it from a WSL `/tmp` UNC path, and npm attempted to write `C:\Windows\package-lock.json` before failing with `EPERM`. A fresh `go test -count=1 ./internal/init` reproduced that environment failure. This is not a valid passing full-gate result.
- Read: `agent-skills/savepoint-design/SKILL.md`, `agent-skills/savepoint-task/SKILL.md`, `agent-skills/references/issue-capture.md`, `.savepoint/router.md`, `.savepoint/Idea.md`, `.savepoint/Design.md`, `.savepoint/Guardrails.md` (CFG-02, TEST-05, TEST-08), the O016 Objective and this Task, `internal/migrate/end_to_end_test.go`, `internal/migrate/replace_windows_test.go`, and `.github/workflows/ci.yml`. Read I029/I030 to check existing Windows findings. Earlier replan investigation read the relevant failing assertions in `internal/migrate/inventory_test.go`, `internal/migrate/operation_test.go`, and `internal/migrate/cutover_test.go` and their matching implementations to classify the out-of-scope full-suite failures.
- Changed: `.savepoint/router.md`, the O016 Objective, this Task's scope/lifecycle/evidence, `.github/workflows/ci.yml`, `internal/migrate/replace_windows_test.go`, I029/I030 incident records, and removed the now-unneeded planned T018/T019 files from O016. Neither Issue was closed.
- Second full-gate attempt: `make test-full` with `TMPDIR=/mnt/c/Users/User/AppData/Local/Temp` let Windows npm use a local path, but the root package and `internal/migrate` were reported failed. The migration package took 8m24.327s. The command stopped at `test` before its `build-all` phase, so this is not passing full-gate evidence; the Windows-mounted temp directory may also affect filesystem-specific behavior. The timing summary showed `TestEndToEnd_temporaryRepositoryCopyMigratesWithReleaseAccountability` at 4m43.03s and `TestApply_recoversAtEveryPublishBoundaryWithoutOverwritingUserEdits` at 3m20.99s. The captured output did not preserve individual failure names.
- Separate `make build-all` passed on `go1.26.2 linux/amd64`, covering the configured Linux, Darwin, and Windows build groups.
- Limitation (at executor handoff, now superseded below): no passing fresh `make test-full` evidence. The focused Windows setup check and standalone cross-build passed; the broad Windows suite is outside the owner-approved scope. I029 remains open; no Issue was closed.

### Full-gate evidence cited from C909

- The independent Full Objective Check C909 ran `make test-full` and it passed with exit 0. It started 2026-09-23T04:14:15Z and finished at 04:16:21Z, taking 2:06.39 wall time. The toolchain was Go 1.26.2 linux/amd64 under WSL2 (kernel 6.18.33.2-microsoft-standard-WSL2), with native Linux Node v24.16.0 and npm 11.13.0 from nvm. The run used no npm bridge and no `TMPDIR` override. `internal/migrate` took 1m48.107s, `TestIntegration_InstallDependencies` passed without being skipped, and all three cross-build phases completed.
- Inputs at that run were HEAD `e7e72bb3d736028d208bc070f012fd6c66e8c930`. `git diff | sha256sum` gave `a1e3a570df1fcc8e25ec1ebdc5490a68afdf1a73faf7b17ee4c8c6aa355bb3a1`. The code, tests, fixtures, dependencies, and gate definitions are unchanged since then. Later edits touched only `.savepoint` metadata: C909, I029, O016 freshness, and this record. Under TEST-08, this counts as metadata-only reuse of that full-gate result.
- C909 also reproduced the focused Windows check independently with Go 1.26.2 windows/amd64. All three named tests passed, and `internal/migrate` took 1.890s. When C909 temporarily added a second `TestMain`, the CI command failed with `[setup failed]`. This proves the focused CI job guards the I029 regression.
- I029 was closed as `verified` by C909. This Task did not close it.
- Per-criterion status: exactly one effective Windows `TestMain` is met. Setup and helper dispatch are met. The focused Windows tests pass. The Windows CI job runs only the focused check. A fresh `make test-full` passed, as cited above. I029 was left for independent Check closure. All six criteria are met.
- Task-check waiver: the owner skipped the optional Task Check for T017 and completed the Task on the board. The frontmatter `check_waiver` records the task, reason, actor (owner, `board-owner`), and time (2026-09-23T04:21:24Z). The owner confirmed that waiver in chat on 2026-09-23. It is not technical CLEAR. C909 is the independent Full Objective Check covering T017.

## Drift Notes

No architecture or verification-policy change is expected. Return to design if
composing the setup requires changing fixture ownership or the migration test
contract.

