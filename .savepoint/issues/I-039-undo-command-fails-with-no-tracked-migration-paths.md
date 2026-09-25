---
id: I-039
title: Printed migrate undo command fails when no planned path is tracked
type: defect
status: resolved
source:
  kind: check
  check: C-912
  actor: {role: checker, session: t025-task-check-20260924}
  at: '2026-09-23T21:00:00Z'
tasks: [T-025]
checks: [C-912, C-913]
severity: low
resolution:
  disposition: verified
  check: C-913
  actor: {role: checker, session: t025-task-check-recheck-20260924}
  at: '2026-09-23T21:18:45Z'
history:
  - at: '2026-09-23T21:00:00Z'
    actor: {role: checker, session: t025-task-check-20260924}
    kind: observed
    check: C-912
    note: >-
      Reproduced with a built binary against a committed minimal V1 project
      during the T-025 Task Check.
  - at: '2026-09-23T21:00:42Z'
    actor: {role: executor, session: t025-task-check-20260924}
    kind: repair_attempted
    note: >-
      Owner-directed repair: gitUndoCommand omits the restore clause when no
      planned path is tracked. The new test runs the exact printed string. A
      fresh make test-full passed. Awaiting a re-check that supersedes C-912.
  - at: '2026-09-23T21:18:45Z'
    actor: {role: checker, session: t025-task-check-recheck-20260924}
    kind: rechecked
    note: >-
      C-913 CLEAR: the printed clean-only command succeeds for an empty
      tracked-path group, removes the generated config and manifest, and
      leaves the original committed content and Git status clean.
    check: C-913
---

# I-039: Printed migrate undo command fails when no planned path is tracked

## Summary

`gitUndoCommand` (`internal/migrate/command.go:287`) always emits
`git restore --source=HEAD --staged --worktree -- <tracked>` followed by
`&& git clean -fdx -- <untracked>`. `gitUndoPathGroups`
(`internal/migrate/command.go:293`) puts a path in the tracked group only for
an existing `config.yml`, the router, or an archived source. When a V1
project has none of those, the tracked group is empty, so the restore runs
with no pathspec, git exits 128 with `fatal: you must specify path(s) to
restore`, and `&&` skips the clean. The generated `config.yml` and
`migrations/v1-to-v2.yml` are left behind.

This breaks T-025's outcome ("tells the user how to undo with git") for the
success message and, through `applyFailure`, for the error report too.

## Evidence

Minimal reproduction (git on PATH):

```sh
mkdir -p p/.savepoint && cd p && echo '# D' > .savepoint/Design.md
git init -q && git add -A && git commit -qm base
savepoint migrate --apply
# migration complete. Undo from the project root with:
#   git --literal-pathspecs restore --source=HEAD --staged --worktree --  && git --literal-pathspecs clean -fdx -- '.savepoint/config.yml' '.savepoint/migrations/v1-to-v2.yml'
sh -c "<printed command>"   # fatal: you must specify path(s) to restore; exit 128
git status --porcelain       # ?? .savepoint/config.yml, ?? .savepoint/migrations/v1-to-v2.yml
```

Expected: running the printed command returns the tree to the commit.

The real-git integration test
(`TestRunCommand_realGitApplyIsReversibleAndRejectsDirtyPaths`) uses only the
`v1-history` fixture, which always has tracked paths, and it runs the two
git steps separately instead of the printed command string. That is why it
does not catch this.

## Proof Needed

When the tracked group is empty, `gitUndoCommand` leaves out the restore
clause (or otherwise emits a command that works). A test runs the exact
printed string in a real temp repo holding a project with no tracked planned
paths, then shows `git status` is clean afterwards.
