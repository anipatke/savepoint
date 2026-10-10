---
id: I-146
title: Session-start hook's Windows .cmd path is untested and breaks on paths with spaces
type: defect
status: resolved
source:
  kind: check
  check: C-974
  actor: {role: checker, session: check-o043-20261010}
  at: '2026-10-10T07:52:00Z'
tasks: [T-117]
checks: [C-974, C-975, C-976]
guardrail_ids: [CFG-02, CFG-03]
severity: medium
resolution:
  disposition: verified
  check: C-976
  actor: {role: checker, session: check-o043-recheck2-20261010}
  at: '2026-10-10T08:10:00Z'
  reason: 'Fix verified in C-975; session-start.js unchanged since, TestSessionStartHook passes in native Windows run 38036497392 at 3596b2c.'
history:
  - at: '2026-10-10T07:52:00Z'
    actor: {role: checker, session: check-o043-20261010}
    kind: observed
    note: Found by the Full Objective Check of O-043.
    check: C-974
  - at: '2026-10-10T07:59:00Z'
    actor: {role: checker, session: check-o043-20261010}
    kind: rechecked
    check: C-975
    note: Native Windows run 38036069019 at 0fb8c37 ran every TestSessionStartHook subtest with a .cmd fake, including the directory-with-a-space case, all PASS; quoted .cmd spawn confirmed by probe.
  - at: '2026-10-10T08:10:00Z'
    actor: {role: checker, session: check-o043-recheck2-20261010}
    kind: rechecked
    check: C-976
    note: 'Fix verified in C-975; session-start.js unchanged since, TestSessionStartHook passes in native Windows run 38036497392 at 3596b2c.'
---

# I-146: Session-start hook's Windows .cmd path is untested and breaks on paths with spaces

## Summary

On Windows the session-start hook runs a found `savepoint.cmd` with `shell: true` (`templates/project-v2/.claude/hooks/session-start.js:20`). With `shell: true`, Node joins the command and arguments into one unquoted command line, so a `savepoint.cmd` whose path contains a space (a project under `C:\Users\First Last\...\node_modules\.bin`, or an npm global prefix in such a profile) fails to start. The hook then adds nothing, silently. The owner gets no Next line and no sign why.

This path is never exercised on any platform. `fakeSavepoint` (`internal/init/claude_hooks_test.go:208`) skips on Windows with the reason "fake shell executable cannot exist on Windows", but a fake `savepoint.cmd` batch file can exist there. So on Windows every `TestSessionStartHook` subtest that needs a fake `savepoint` skips, and only "savepoint missing" runs. `TestFindSavepoint_windowsNames` checks only the lookup, not the spawn.

Violates CFG-03 (a Windows skip is allowed only when the situation cannot exist there), CFG-02 (platform differences explicit and tested), and T-117 Done When "Hook commands ... work on Windows".

## Evidence

- Mechanism reproduced on Linux by forcing the same option: `spawnSync("<scratch>/dir with space/savepoint", ["resume"], {shell: true})` → status 127, `/bin/sh: ... not found`; the same call with `shell: false` → status 0, `Next action: ok`. Node builds the Windows `cmd.exe /d /s /c "..."` line the same unquoted way.
- `grep -n "Skip" internal/init/claude_hooks_test.go` → line 209, the Windows skip above.

Expected: on Windows, a `savepoint.cmd` found on PATH or in `node_modules/.bin`, including under a directory with a space, yields the Next line.

## Proof Needed

- The `.cmd` spawn quotes the path (or avoids `shell: true`) so a directory with spaces works.
- The Windows skip is removed: on Windows the fake is a `savepoint.cmd` batch file, so the Next-line, failure, no-Next, and `node_modules/.bin` subtests run there; one subtest puts the fake under a directory whose name contains a space.
- Native Windows CI runs those subtests (see I-147); `make build && make test-fast` passes; a re-check confirms.
