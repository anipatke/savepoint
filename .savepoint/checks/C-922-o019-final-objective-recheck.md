---
id: C-922
scope: {kind: objective, id: O-019}
result: CLEAR
checked_by: {role: checker, session: o019-final-recheck-20260925}
executed_session: unrecorded-external-executor
checked_at: '2026-09-24T23:46:20Z'
reviewed:
  base_commit: 35ce82dffbcfa02a88121f99496d668c67d69db9
  head_commit: 35ce82dffbcfa02a88121f99496d668c67d69db9
  files:
    - cmd/create_task.go
    - cmd/init_test.go
    - main.go
    - internal/data/project_test.go
    - internal/data/task_create.go
    - internal/data/task_ids.go
  dependencies: []
issues: []
supersedes: C-921
---

# C-922: O-019 Final Full Objective Recheck

## Closure Map and Independence

| Prior Issue | Result | Evidence |
| --- | --- | --- |
| I-049 | Verified by this CLEAR Check | A real `create-task` run with stdout redirected to `/dev/full` exits 0 after persisting T-001, warns on stderr with its ID and path, and leaves one Task and `last_issued: 1`. |

This checker session did not implement the repair. C-920's numbered scope lock and matrix remain the perimeter; C-921's admission ledger binds this final targeted recheck to the original CLI output/retry cell. The owner asked for every post-write path in that workflow to be probed. This record reports the secondary cleanup limitation without treating a new synthetic test setup as a product failure.

## Admission Ledger and Frozen Matrix

| Recheck item | Prior cell | Probe and result |
| --- | --- | --- |
| Stdout failure after persistence and retry signal | C-920 CLI dispatch/output; I-049 | `TestRunCreateTaskOutputFailureSucceedsWithStderrWarning` passes. Independent real binary with `/dev/full`: exit 0, warning names T-001 and path, one Task remains. **Proven.** |
| Strict-load failure after Task write | C-920 draft/direct creation and workflow operation 5 | Independent CLI draft with missing T-999 dependency: exit 1, T-002 removed, T-001 untouched, next successful creation issues T-003 and `last_issued: 3`. Existing strict-validation rollback test passes. **Proven.** |
| Lock-release failure after write | C-920 lock/failure/recovery and workflow operation 5 | `TestWithTaskIDReservationLockReleaseFailureRollsBackAndRetiresID` passes on Linux: synthetic release error removes its created file and retains the ID reservation. Source trace shows the callback rollback is invoked when release fails. **Proven for the scoped Linux path; Windows test limitation below.** |
| Rollback itself fails | C-920 workflow cleanup/secondary failure | `TestWithTaskIDReservationRollbackFailureIsReported` passes: primary and rollback errors are joined. The test injects a failing callback; it does not create a real Task, so it proves diagnostic propagation, not that a real Task remains. **Observed limitation below.** |
| Allocator identity, invalid input, concurrency, high-water boundary and lock timeout | C-920 allocator and lock rows | Fresh full suite passes the named `internal/data` tests. The earlier independent T-999→T-1000/T-1001 and leftover-lock probes remain applicable; production allocator code is unchanged by this repair. **Proven.** |
| Cross-Objective creation, frontmatter/path/index, and guidance | C-920 cross-Objective and planner rows | Fresh full suite reruns concurrent CLI/data, byte preservation, scaffold, upgrade, and guidance tests. Changes are confined to command output and post-write tests. **Proven.** |
| Linux/Windows release gate | C-920 platform row | Fresh `make test-full` passes on Linux with Linux/macOS/Windows builds. The exact focused native Windows CI set (`TestAllocateTaskID_*`, duplicate diagnostic, concurrent creation) passes from a fresh NTFS copy. **Proven for the required platform cell.** |

## Workflow and Side Effects

The production sequence remains C-920's parse → strict-load → lock → reserve/sync → exclusive Task write → strict-load → unlock → report. Before and during the Task write, malformed input and occupied paths fail without overwriting existing content. After the write, strict-load failure rolls back the new Task and retires the ID. A lock-release error runs the rollback callback before returning failure. After successful persistence and lock release, stdout failure now returns success and emits a warning on stderr. The point of externally reported success therefore agrees with the committed Task. No server, network response, secret, colour, or redirect behavior beyond the documented CLI output path applies.

## Acceptance and Guardrails

- O-019 Success Conditions 1–7: Proven against the frozen matrix and fresh Full gate. T-010 Done When 1–4, T-011 Done When 1–5, and T-012 Done When 1–5: Proven. The three Tasks remain `done` with owner Task-check waivers; the waivers were inspected directly rather than treated as technical `CLEAR`.
- FS-01/04..06, DATA-01/03/05, TPL-01/02/04, ARCH-01/03/04, CFG-01/02, and TEST-01..05/08/09: satisfied for the supported paths and required gates. The CLI failure-path regression now covers the I-049 reproduction. No user-authored file is overwritten by the new creation path.
- Design sections 1, 2, 5, 6, and 9 remain reconciled with the allocator, lock, ID-free command, and strict-load behavior. This targeted repair did not change those contracts.

## Gates and Independent Probes

- `git diff --check`: passed.
- `make build`: passed.
- `make test-full`: passed fresh after the repair, `go1.26.2 linux/amd64`, exit 0; includes the complete Go suite and Linux/macOS/Windows builds.
- Focused post-write tests for strict-validation rollback, lock-release rollback, rollback-error reporting, and the CLI stderr warning: passed on Linux.
- Real binary `/dev/full`/strict-load/next-ID sequence: stdout failure exited 0 with one committed Task; invalid dependency then exited 1 and removed its new Task; next success used T-003. High-water finished at 3.
- Native Windows focused CI set: passed, `go1.26.2 windows/amd64`, exit 0.

## Materiality and Observations

I-049 is verified; no materiality action remains inside the frozen Issue. No new Issue is admitted from the completed matrix.

- A broader native Windows run that included `TestWithTaskIDReservationLockReleaseFailureRollsBackAndRetiresID` failed twice at `internal/data/project_test.go:1600`. Its test callback tries to remove the lock file while the allocator still holds it; Windows refuses that setup operation, so the callback returns before installing the rollback it intends to test. The required Windows CI test set passes, and this synthetic failure does not demonstrate a production Task-creation failure. The test needs a portable failure-injection method if the full Windows suite is later brought into scope (see the already open I-030 for that deferred full-suite work).
- If an actual rollback cannot remove a newly written Task after another failure, `withTaskIDReservation` returns both the primary and rollback errors. The injected callback test proves error reporting only. A secondary filesystem failure or concurrent outside edit could still leave a Task while returning an error; the Check does not claim a cleanup guarantee in that state. This rare compound failure is recorded as a residual observation, and no supported real-path reproduction was established in this run.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches**
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [x] STYLE-07 **One source of truth**
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [ ] STYLE-10 **Small diffs** — the integrated Objective still spans allocator, CLI, tests, and guidance; this is advisory only.

## Owner Handoff

Technical clearance is current with C-922. I-049 is resolved as verified by this Check. O-019 is ready for the owner's completion decision; the checker does not set Objective status to `done` or accept the R-006 Goal Check for the owner.
