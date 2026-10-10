---
id: I-147
title: Supply native Windows evidence for O-043
type: verification
status: open
source:
  kind: check
  check: C-974
  actor: {role: checker, session: check-o043-20261010}
  at: '2026-10-10T07:52:00Z'
tasks: [T-115, T-116, T-117, T-118]
checks: [C-974]
guardrail_ids: [CFG-03]
severity: medium
history:
  - at: '2026-10-10T07:52:00Z'
    actor: {role: checker, session: check-o043-20261010}
    kind: observed
    note: Found by the Full Objective Check of O-043.
    check: C-974
---

# I-147: Supply native Windows evidence for O-043

## Summary

CFG-03 (Blocker) requires the full Go test suite to pass natively on Windows in CI. O-043 adds Node hooks, a git-worktree lane check, a `.cmd` lookup, and a new embed root (`all:templates/project-v2`), all platform-sensitive. All four Task records say native Windows CI was not run. The O-043 work is uncommitted on top of `0e95269`, so no CI run covers it.

## Evidence

- T-115, T-117, T-118 Technical Evidence each list native Windows CI as not run.
- `git status` at Check time: every O-043 code file is modified or untracked; `HEAD` = `0e95269`.
- `.github/workflows/ci.yml` defines a `windows-tests` job (`go run ./internal/buildtool test -json -count=1 ./...` on `windows-latest`) that would supply this evidence once the work is pushed.

## Proof Needed

- A successful `ci` and `windows-tests` run on a commit that contains the O-043 work (and the I-146 fix), with its run ID and head commit recorded.
- A re-check confirms the run matches the reviewed head.
