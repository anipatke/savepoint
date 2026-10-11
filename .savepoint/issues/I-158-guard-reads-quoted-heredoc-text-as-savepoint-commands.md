---
id: I-158
title: Router guard reads quoted heredoc text as savepoint commands
type: defect
status: open
source:
  kind: report
  actor: {role: owner, session: owner-20261011}
  at: '2026-10-11T06:05:00Z'
guardrail_ids: []
severity: low
history:
  - at: '2026-10-11T06:05:00Z'
    actor: {role: executor, session: check-o048-recheck-20261011}
    kind: observed
    note: 'Hit while updating PR #25 after the O-048 recheck; the owner asked for it to be recorded. Not repaired.'
---

# I-158: Router guard reads quoted heredoc text as savepoint commands

## Summary

The Claude Code guard hook (`.claude/hooks/guard.js`, shipped copy
`templates/project-v2/.claude/hooks/guard.js`) splits a Bash command into
simple commands in `simpleCommands` (lines 93-128) without understanding
heredocs. The body of a heredoc is treated as more shell: every newline and
backtick ends a "command". When a quoted heredoc (`<<'EOF'`, whose body the
shell never runs) contains a line starting with `savepoint` or a backticked
`` `savepoint run` ``, `checkBash` returns `MSG_CLI` and blocks the whole
Bash call, including writing a harmless file.

Unquoted heredocs (`<<EOF`) do run backtick substitutions in a real shell, so
blocking those is correct; only quoted-delimiter heredoc bodies are false
positives.

## Evidence

Reproduced with `decide()` from `.claude/hooks/guard.js` at 6a5dc6c:

| Command | Guard | Correct |
|---|---|---|
| ``cat > f.md <<'EOF'`` / ``- `savepoint run` is gone.`` / `EOF` | blocked | allow (literal text) |
| `cat <<'EOF'` / `savepoint board` / `EOF` | blocked | allow |
| `cat <<EOF` / `savepoint board` / `EOF` | blocked | allow (line is text, not a command) |
| ``git commit -m "Remove `savepoint run`"`` | allowed | allow |
| `gh api repos/anipatke/savepoint/pulls/25` | allowed | allow |

Real occurrence: writing the v2.3.1 PR description (which mentions
`` `savepoint run` ``) with a quoted heredoc was refused with "Blocked: only
the owner runs this savepoint command…". The owner had to run the PR update
by hand.

Missing test: `internal/init/claude_guard_test.go` `TestGuard_savepointCommands`
has no heredoc case.

Separate, not this Issue: T-138 notes the guard also blocks a bare
`./savepoint` (no arguments). That opens the board, an owner command, so
blocking it is by design.

## Proof Needed

Quoted-heredoc bodies are skipped when splitting commands; unquoted heredoc
bodies still block a backticked or `$(...)` savepoint command. Add guard test
cases for both, keep the repository and scaffold copies identical
(`TestGuard_repositoryCopyMatchesTheScaffold`), and pass the configured gates.
