---
id: I-148
title: Router guard blocks router edits on the main checkout on Windows
type: defect
status: open
source:
  kind: check
  check: C-975
  actor: {role: checker, session: check-o043-20261010}
  at: '2026-10-10T07:59:00Z'
tasks: [T-118]
checks: [C-975]
guardrail_ids: [CFG-02, CFG-03]
severity: medium
history:
  - at: '2026-10-10T07:59:00Z'
    actor: {role: checker, session: check-o043-20261010}
    kind: observed
    note: Found by the O-043 Full Objective re-check from the native Windows CI run.
    check: C-975
---

# I-148: Router guard blocks router edits on the main checkout on Windows

## Summary

On native Windows the PreToolUse guard treats an ordinary (main) checkout as a worktree lane and denies an Edit of `.savepoint/router.md` with "Blocked: do not edit .savepoint/router.md in a worktree lane ...". T-118 Done When requires router edits to be blocked only when `--git-dir` differs from `--git-common-dir` and allowed on the main checkout. O-043 SC4 limits the block to "inside a worktree lane". The native Windows job fails, so CFG-03 is not met.

## Evidence

- GitHub Actions run 38036069019 (https://github.com/anipatke/savepoint/actions/runs/38036069019), head `0fb8c3777884e4bd7e9f0b61db881aa9dfa3991e`, job `windows-tests`: `claude_guard_test.go:151: main checkout: want allow, got "Blocked: do not edit .savepoint/router.md in a worktree lane; ..."`, `--- FAIL: TestGuard_routerInLane`, `FAIL github.com/opencode/savepoint/internal/init`. Every other `internal/init` test in that job passed.
- On Linux the same test and an independent real-worktree probe pass (lane blocked, main allowed).
- Code: `inLane` in `templates/project-v2/.claude/hooks/guard.js` resolves each `git rev-parse` output with `path.resolve(cwd, ...)` and compares `fs.realpathSync` of the two. Likely cause, not verified on Windows: one of the two paths comes back in a different form (for example an 8.3 short name such as `RUNNER~1` in the temp directory, or a relative value resolved against a different base), and the JavaScript `fs.realpathSync` does not normalise it, so the two never compare equal on Windows. `git rev-parse --path-format=absolute` or `fs.realpathSync.native` are candidate fixes.

Expected: on Windows, main checkout allows; lane blocks; outside a repo allows.

## Proof Needed

- `TestGuard_routerInLane` passes in the native `windows-tests` job, with the main-checkout and lane cases unchanged.
- The fix keeps the Linux behaviour; `make build && make test-fast` passes; a re-check confirms with the CI run ID and head commit.
