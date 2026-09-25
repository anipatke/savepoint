---
id: I-061
title: Task ID allocator refuses live contention as a leftover lock on Windows
type: defect
status: in_progress
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
