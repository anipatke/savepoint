---
id: C-921
scope: {kind: objective, id: O-019}
result: NEEDS WORK
checked_by: {role: checker, session: o019-objective-recheck-20260925}
executed_session: owner-repair-20260925
checked_at: '2026-09-24T22:05:52Z'
reviewed:
  base_commit: 35ce82dffbcfa02a88121f99496d668c67d69db9
  head_commit: 35ce82dffbcfa02a88121f99496d668c67d69db9
  files:
    - cmd/create_task.go
    - cmd/init_test.go
    - internal/data/task_ids.go
    - internal/data/task_create.go
  dependencies: []
issues: [I-049]
supersedes: C-920
---

# C-921: O-019 Full Objective Recheck

## Closure Map and Frozen Scope

| Prior Issue | Status | Reason |
| --- | --- | --- |
| I-049 | Still open | The command now prints a truthful error, but its nonzero exit after persistence still permits an automatic retry to create a duplicate Task. |

This is an independent Check of the owner's `cmd/create_task.go` and `cmd/init_test.go` repair. It uses C-920's numbered scope lock, workflow inventory, and matrix without adding a new axis. All three O-019 Tasks remain owner-completed with recorded optional Task-check waivers. No Task status is changed.

## Recheck Admission Ledger

| Item | Prior Issue or repair claim | Exact C-920 matrix cell | Allowed result |
| --- | --- | --- | --- |
| Failing writer test and CLI `/dev/full` reproduction | I-049; error now names persisted Task | CLI dispatch and output: redirected stdout failure after Task persistence | Close only if command outcome and retry behavior satisfy O-019 Condition 2 and I-049 Proof Needed |
| Same-draft retry following nonzero exit | I-049's safe retry requirement | CLI dispatch and output: retry after redirected stdout failure | Remains blocking if retry duplicates work |
| Allocation boundary, held lock, malformed input, and rollback tests | Unchanged original behavior | Allocator identity; invalid/concurrent state; lock/failure/recovery; draft/direct creation | Reconfirm or mark unverified |
| Two-Objective creation and ID/record representations | Unchanged original behavior | Cross-Objective and representation | Reconfirm or mark unverified |
| Native Windows focused tests and full gate | Original platform/release gates | Linux/Windows | Reconfirm or mark unverified |
| Skill parity, fresh init, upgrades, README/Design | Unchanged guidance | Planner and shipped guidance | Reconfirm or mark unverified |

## Matrix Results

| Original C-920 cell | Recheck evidence | Result |
| --- | --- | --- |
| Allocator identity and boundary | `make test-full` reran `TestAllocateTaskID_*`; independent CLI fixture again issued T-1000 and T-1001 above active T-999, with `last_issued: 1001`. | Proven |
| Allocator invalid and concurrent state | Full gate reran duplicate-path, malformed-high-water, and concurrency tests. | Proven |
| Lock/failure/recovery | Full gate reran held/leftover-lock and rollback tests; independent leftover-lock CLI probe exited 1, named the lock, and left both existing Tasks unchanged. | Proven |
| Draft parser and direct creation | Full gate reran `TestCreateTaskV2*` and exclusive-file tests; repair did not alter `internal/data`. | Proven |
| CLI dispatch and output | `TestRunCreateTaskOutputFailureNamesPersistedTask` passes: error includes ID, path, write failure, and `do not retry`, with one runner call. Fresh CLI `/dev/full` probe: attempt 1 exits 1 with that error and leaves T-001; a retry on exit 1 creates T-002, also exits 1, and high-water is 2. The command status still reports failure for a completed creation. | **I-049 remains open** |
| Cross-Objective and representation | Full gate reran concurrent `CreateTaskV2` and CLI workflow tests; both owner links and strict index remain valid. | Proven |
| Linux/Windows | `make test-full` passed on Linux with Linux/macOS/Windows builds. Focused `internal/data` allocator, duplicate, and concurrent creation tests passed natively on Windows from a fresh temporary NTFS copy (`go1.26.2 windows/amd64`). | Proven |
| Planner and shipped guidance | Full suite's scaffold, upgrade, and skill parity tests passed; guidance paths are unchanged by this repair. | Proven |

The new error wording improves the human diagnostic, but it does not satisfy the frozen CLI output/retry cell. No other original cell failed or became unverified.

## Workflow and Side Effects

C-920's operations 1–5 still have the same ordering and cleanup. At operation 6, `RunCreateTask` now wraps the output error with the committed ID and path. It still returns that error to `main`, which exits nonzero after the Task and high-water mark have been committed. The independent two-attempt CLI probe verifies semantic state, not just the output text. No network, secret, redirect, or cancellation cell beyond the original lock applies.

## Criteria and Guardrails

- O-019 Success Conditions 1 and 3–6, T-010 Done When 1–4, T-011 Done When 1 and 3–5, and T-012 Done When 1–5 remain proven. O-019 Condition 2 and T-011 Done When 2 remain Issue I-049 because the reported nonzero outcome is unsafe to retry. Condition 7's Full gate passed, but the Full Check remains `NEEDS WORK`.
- The original applicable FS, DATA, TPL, ARCH, CFG, and TEST rules were re-evaluated through the frozen matrix. The new failure-path test improves TEST-02 evidence; it does not clear the unmet acceptance condition. The recorded Task-check waivers remain valid as waivers, not technical `CLEAR`.
- Design section 9 still matches the allocator, lock, and strict-loader behavior. It does not supply a safe command-result contract for the output failure.

## Gates

- `go test ./cmd -run '^TestRunCreateTaskOutputFailureNamesPersistedTask$' -count=1`: passed.
- `git diff --check`: passed.
- `make build`: passed.
- `make test-full`: passed fresh after the repair, including all Go packages and Linux/macOS/Windows builds; `go1.26.2 linux/amd64`, exit 0.
- Focused native Windows `internal/data` tests: passed, `go1.26.2 windows/amd64`, exit 0.
- Independent normal/boundary/lock CLI probes: passed. Independent `/dev/full` then nonzero-exit retry probe: reproduced I-049.

## Materiality

| Issue | Likelihood | Impact | Materiality | Recommendation |
| --- | --- | --- | --- | --- |
| I-049 | Low: stdout failure is uncommon, but script retries on nonzero are normal | Medium: the same draft can become two valid Tasks without the caller intending it | Medium within O-019's failed-creation promise | Make the exit status represent successful persistence, while reporting the output failure separately; then recheck this exact cell |

## Observations (non-blocking)

- The warning `do not retry` is clear for a person reading stderr. It cannot change a caller that retries based solely on exit status.

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
- [ ] STYLE-10 **Small diffs** — the full O-019 work still spans allocator, CLI, tests, and guidance; this advisory observation does not affect the verdict.

## Owner Decision

I-049 remains open and O-019 is not ready for owner closure. The checker does not repair the command or mark the Objective complete. A subsequent targeted repair can be rechecked against C-920's original CLI output/retry cell; C-920 and this record remain immutable.
