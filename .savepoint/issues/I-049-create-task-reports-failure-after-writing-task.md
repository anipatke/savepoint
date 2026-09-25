---
id: I-049
title: create-task reports failure after writing the Task
type: defect
status: resolved
source:
  kind: check
  check: C-920
  actor: {role: checker, session: o019-objective-check-20260925}
  at: '2026-09-24T21:22:30Z'
tasks: [T-011]
checks: [C-920, C-921, C-922]
guardrail_ids: [TEST-02]
severity: medium
resolution:
  disposition: verified
  check: C-922
  actor: {role: checker, session: o019-final-recheck-20260925}
  at: '2026-09-24T23:46:20Z'
  reason: The real CLI now exits 0 and warns on stderr when stdout fails after the Task is committed; a retry on nonzero is no longer prompted.
history:
  - at: '2026-09-24T21:22:30Z'
    actor: {role: checker, session: o019-objective-check-20260925}
    kind: observed
    check: C-920
    note: A redirected stdout write failure makes create-task exit nonzero after its Task file and high-water mark have been committed.
  - at: '2026-09-24T22:05:52Z'
    actor: {role: checker, session: o019-objective-recheck-20260925}
    kind: rechecked
    check: C-921
    note: The error now names the committed Task and warns not to retry, but a nonzero-exit retry still creates another Task from the same draft.
  - at: '2026-09-24T23:46:20Z'
    actor: {role: checker, session: o019-final-recheck-20260925}
    kind: rechecked
    check: C-922
    note: The exact real-binary stdout-failure reproduction now exits 0, retains one valid Task, and warns with its ID and path; the original Issue is verified.
---

# I-049: create-task reports failure after writing the Task

## Summary

O-019 Success Condition 2 and T-011's Done When require a failed creation to leave no new Task file. A failure writing the CLI success message instead leaves a new Task while the command exits nonzero. Retrying the same draft creates another Task with a new ID, so the planner can unknowingly duplicate the work.

## Evidence

- On a valid temporary V2 project, run `./savepoint create-task --objective O-001 --draft draft.md <project>` with stdout sent to `/dev/full`. The command exits 1 and prints `write /dev/stdout: no space left on device` to stderr. The project contains `T-001-output-failure.md` and `.savepoint/task-ids.yml` says `last_issued: 1`.
- `cmd/create_task.go:30-35` calls the creation runner before writing the success message, then returns the output write error as the command result. `internal/data/task_create.go` has already completed the exclusive file write and strict load.
- Existing tests cover successful output and creation failures, but no test covers output failure after persistence. The full Linux gate and focused native Windows tests do not exercise this case.

## Proof Needed

- Make the command's reported outcome accurately reflect the persisted Task when output fails, or ensure a failed command leaves no Task while retaining the reserved ID. Preserve the never-reuse guarantee and avoid deleting any pre-existing file.
- Add a regression scenario with a failing output writer or redirected sink. Verify the command result, Task files, high-water mark, and safe retry behavior.

## Recheck C-921

The new failing-writer test proves that the error names the committed Task and warns a human not to retry. The independent `/dev/full` CLI reproduction still exits 1 after creating T-001. Repeating the command on that nonzero result creates T-002; `last_issued` becomes 2 and both Task files remain. The original safe-retry requirement is therefore still open. The Full Linux gate and focused native Windows allocation tests passed.

## Final Recheck C-922

The real `create-task` command with stdout redirected to `/dev/full` now exits 0. Stderr warns that T-001 was created and names its path. Exactly one Task exists, with `last_issued: 1`. The new `TestRunCreateTaskOutputFailureSucceedsWithStderrWarning` passes, and a fresh `make test-full` passes. C-922 records the full Objective verdict and the secondary cleanup observations separately.
