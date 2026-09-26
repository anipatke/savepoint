---
id: I-061
title: Task ID allocator refuses live contention as a leftover lock on Windows
type: defect
status: open
source:
  kind: report
  actor: {role: planner, session: i061-windows-allocator-20260926}
  at: '2026-09-25T21:00:16Z'
tasks: [T-010]
severity: high
history:
  - at: '2026-09-25T21:00:16Z'
    actor: {role: planner, session: i061-windows-allocator-20260926}
    kind: observed
    note: >-
      Owner asked for the Windows CI failure to be investigated before merging
      v2 into master. TestAllocateTaskID_serializesConcurrentCallers fails in
      the windows-tests job in 3 of the 6 runs since that job settled; the
      cause and fix options are recorded in the Analysis section.
  - at: '2026-09-25T21:06:25Z'
    actor: {role: executor, session: i061-windows-allocator-20260926}
    kind: repair_attempted
    note: >-
      Applied fix 1: acquireTaskIDLock restarts its taskIDLockWait deadline
      whenever the lock observed has a different modification time, and retries
      at once when the lock vanishes. New regression test
      TestAllocateTaskID_waitsWhileTheLockChangesHands fails on the previous
      allocator and passes with the fix; the existing held and leftover lock
      refusals still pass. make build and make test-full pass on Linux.
      Windows CI proof is still needed.
  - at: "2026-09-25T21:06:40Z"
    actor:
      role: owner
      session: board-owner
    kind: owner_decision
    note: Moved from open to in_progress by the owner from the board.
  - at: "2026-09-25T21:14:14Z"
    actor:
      role: owner
      session: board-owner
    kind: owner_decision
    note: Resolved by the owner from the board.
  - at: '2026-09-25T23:05:42Z'
    actor: {role: owner, session: user}
    kind: reopened
    note: >-
      Reopened on the owner's instruction to fix I-061 properly for 2.0.3.
      The per-holder deadline repair does not hold on Windows:
      TestAllocateTaskID_serializesConcurrentCallers failed in windows-tests
      on PR #11 (run 36198435217) with 4 of 16 reservations, and again on
      rerun with 1 of 16. The owner-accepted resolution is withdrawn.
  - at: '2026-09-25T23:33:01Z'
    actor: {role: executor, session: i061-repair-20260926}
    kind: repair_attempted
    note: >-
      Windows diagnostic (draft PR #12, run 36199665922, closed unmerged):
      successive lock holders get distinct modification times once they are
      2 ms or more apart; one uncontended allocation takes a median 47 ms
      (22-72 ms) against about 2 ms on Linux; 16 concurrent callers finished
      at up to 828 ms and 2 were refused at about 828 ms. The 500 ms
      per-holder staleness guess is too tight for Windows filesystem
      latency. Commit 845b1be removes the modification-time heuristic:
      acquireTaskIDLock now waits up to taskIDLockWait, raised to 10 s, and
      still refuses a lock present after that as leftover.
      TestAllocateTaskID_waitsForALiveHolder holds a live lock for 1.5 s and
      requires the allocator to wait it out; it fails on the previous logic
      with the I-061 refusal. The refusal test uses a 300 ms wait; the old
      modification-time handoff test was removed with the mechanism. make
      build, make test-fast, and make test-full passed on Linux. windows-tests
      passed 5 of 5 runs on 845b1be (run 36199992156, attempts 1-5,
      23:11-23:32 UTC), including TestAllocateTaskID_serializesConcurrentCallers.
      Issue remains open for independent verification or owner acceptance.
  - at: '2026-09-25T23:45:21Z'
    actor: {role: owner, session: user}
    kind: owner_decision
    note: >-
      Owner instructed the executor to mark this Issue resolved once 2.0.3 was
      live. The fix shipped in savepoint 2.0.3 on npm (tag v2.0.3, merge
      c62c366). No technical CLEAR is implied.
  - at: '2026-09-26T04:44:12Z'
    actor: {role: owner, session: user}
    kind: reopened
    note: >-
      Reopened on the owner's instruction before merging v2 into master. A
      different Windows failure mode: the master CI run for PR #17 (run
      36215377170, 2c567a6, first attempt) failed
      TestAllocateTaskID_serializesConcurrentCallers with one caller getting
      "create lock ...task-ids.lock: Access is denied" and 15 of 16
      reservations. The rerun passed. The owner-accepted resolution is
      withdrawn.
  - at: '2026-09-26T04:49:37Z'
    actor: {role: executor, session: v2-main-flaky}
    kind: repair_attempted
    note: >-
      acquireTaskIDLock now asks taskIDLockBusy whether a failed exclusive
      create means the lock is busy: "exists" on every platform, and on
      Windows also access denied, the delete-pending state. Busy failures
      retry within the existing taskIDLockWait; a Windows permission error
      that lasts the whole wait reports the real error. Added
      TestTaskIDLockBusy_retriesWindowsDeletePendingOnly. GOOS=windows go vet,
      make build, and make test-full passed on Linux, and the concurrency test
      passed 30 of 30 there. Windows proof awaits the windows-tests CI job.
---

# I-061: Task ID allocator refuses live contention as a leftover lock on Windows

## Summary

`AllocateTaskID` waits at most 500 ms in total for `task-ids.lock`, and
reports any caller still waiting at that point as facing "a leftover lock
[that] requires owner cleanup". The deadline starts when a caller begins
waiting, not when the current holder took the lock, so a caller queued behind
several live holders is refused even though every holder is healthy. On the
Windows CI runner, 16 concurrent callers cross that deadline and the test
fails; v2 CI is red on every such run, which blocks merging v2 into master.

## Evidence

- windows-tests job, `internal/data`:
  `--- FAIL: TestAllocateTaskID_serializesConcurrentCallers (0.54s)`, with
  three callers reporting
  `allocate Task ID: lock file C:\Users\RUNNER~1\...\.savepoint\task-ids.lock is present; refusing after 500ms (a leftover lock requires owner cleanup)`
  and `unique reservations = 13, want 16`.
- The same failure, at 0.54-0.57 s, in runs 36129088982 (4ae8762),
  36187571259 (3b1e607), and 36188137830 (5ce6b01). Windows runs 36125466183,
  36126277219, and 36129100326 passed. Run 36125057404 (05d161f), which
  introduced the job, failed across many packages and is not counted.
- The `ci` (Linux) job passes on every one of those runs.
- Linux, local: the test passed 30 of 30 runs at 0.15-0.16 s. Sixteen
  sequential allocations take about 30 ms in total (about 1.8 ms each), so
  the concurrent run's time is almost entirely the 10 ms retry sleep
  (`taskIDLockRetry`) between handoffs; even on Linux the test spends about a
  third of the 500 ms budget.

## Analysis

Recorded 2026-09-26 (analysis only; no code changed).

- `acquireTaskIDLock` (`internal/data/task_ids.go`) computes one deadline,
  `time.Now().Add(taskIDLockWait)`, when a caller starts, and never extends
  it. A caller behind N holders must see all N finish within 500 ms.
- Each hold is a strict project load under the lock, a temp-file write and
  rename of `task-ids.yml` (`replaceV2File`), and a close plus remove of the
  lock. Those filesystem operations are slower on the Windows runner, and
  each handoff also costs up to one 10 ms retry sleep per waiter.
- Polling is not first-come-first-served: a waiter can lose the race to
  newer waiters repeatedly, so the worst-placed caller can wait longer than
  its queue position implies.
- The refusal message is wrong in this case: the lock is not left over; it
  has changed hands several times while the caller waited.

Product impact is lower than the test's: real contention is two or three
agents allocating at once, not sixteen. It is still possible on a slow
Windows or network filesystem, and the resulting message tells the owner to
delete a lock that is in active use.

## Possible Fixes

1. Recommended: measure staleness per holder, not per wait. While waiting,
   stat the lock file and restart the 500 ms deadline whenever it is a
   different file from the one last seen (`os.SameFile` or a changed
   modification time). A lock held by one holder for 500 ms is still refused
   as leftover, as T-010 intended, but a queue of healthy holders is not. The
   test stays as it is.
2. Test-only: lower the concurrency or raise the budget for the test. This
   makes CI green but leaves the misleading refusal in the product.
3. Shorten `taskIDLockRetry` (for example to 2 ms). This reduces handoff
   overhead but does not remove the fixed-total-deadline flaw.

## Proof Needed

Show `TestAllocateTaskID_serializesConcurrentCallers` passing repeatedly in
the windows-tests job, and add a regression test in which a lock that keeps
changing hands for longer than `taskIDLockWait` does not cause a refusal,
while a single lock held longer than `taskIDLockWait` still does.

## Recurrence After 2.0.3 (2026-09-26)

The 2.0.3 repair (a 10 s wait measured per caller) fixed the deadline
failure; none of the later Windows runs fail that way. A different
Windows-only failure now appears in the same test.

Evidence:

- Run 36215377170 (master, merge of PR #17, 2c567a6), windows-tests, first
  attempt: `discover_test.go:398: AllocateTaskID() concurrent call error =
  allocate Task ID: create lock C:\Users\RUNNER~1\...\.savepoint\task-ids.lock:
  open ...task-ids.lock: Access is denied.` and `unique reservations = 15,
  want 16`. The Linux `ci` job passed; a rerun of the Windows job passed.
- Across the last 60 CI runs, the only failures of this test are this one and
  the three pre-2.0.3 deadline failures already recorded above.

Analysis:

- `acquireTaskIDLock` retries only when `os.OpenFile(..., O_CREATE|O_EXCL)`
  returns `os.ErrExist`; any other error fails the caller at once.
- `releaseTaskIDLock` closes the lock and then calls `os.Remove`. On Windows
  a deleted file stays in a "delete pending" state until every handle to it
  is closed (another process such as an antivirus scanner can hold one
  briefly). Creating a file with that name meanwhile fails with
  `ERROR_ACCESS_DENIED`, which Go reports as a permission error, not
  `os.ErrExist`. So a caller that lands in that brief window is refused
  outright instead of waiting its turn. This is inferred from the error and
  Windows file semantics; it was not reproduced locally (no Windows host).

Possible fix: on Windows, treat a permission error from the exclusive create
as "lock busy" and keep retrying within the existing `taskIDLockWait`
deadline, reporting the last real error if the deadline passes. Linux
behavior is unchanged.

Proof Needed for the recurrence:

- A Windows access-denied error while another allocator releases the lock is
  retried, not reported; a lasting permission problem still fails with the
  real error once the wait ends.
- Unit coverage of the retry decision on every platform, plus the Windows CI
  job passing `TestAllocateTaskID_serializesConcurrentCallers` on repeated
  runs (for example `-count=20` in a diagnostic run).

### Repair Attempt Evidence (2026-09-26)

- `internal/data/task_ids.go`: new `taskIDLockBusy(err, goos)`;
  `acquireTaskIDLock` retries while it returns true and, after the wait,
  reports a lasting non-exists error as
  `create lock ...: still failing after 10s: <error>`.
- `internal/data/discover_test.go`
  `TestTaskIDLockBusy_retriesWindowsDeletePendingOnly` covers "exists" on
  Linux and Windows, access denied on Windows (retried), access denied on
  Linux (not retried), and other Windows errors (not retried).
- Linux: `make build`, `make test-full`, and
  `go test -count=30 -run TestAllocateTaskID_serializesConcurrentCallers
  ./internal/data` passed. `GOOS=windows go vet ./internal/data` passed.
- Not yet proven on Windows; the next windows-tests CI run is the evidence.
