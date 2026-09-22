---
id: I027
title: Frontmatter parse errors reach the repo before anything catches them
type: guardrail
status: in_progress
source:
  kind: report
  actor: {role: owner, session: user}
  at: '2026-09-22T00:00:00Z'
severity: medium
history:
  - at: '2026-09-22T00:00:00Z'
    actor: {role: owner, session: user}
    kind: observed
    note: >-
      A hand edit to I026 dropped the closing "---" frontmatter delimiter and
      committed cleanly; the board only reported "no closing frontmatter
      delimiter found" later. Nothing in the repo validates frontmatter
      before it lands. Owner reported this as recurring.
  - at: '2026-09-22T00:00:01Z'
    actor: {role: executor, session: user}
    kind: repair_attempted
    note: >-
      Added scripts/git-hooks/pre-commit and check_frontmatter.py: a git
      pre-commit hook that YAML-validates every staged .savepoint/**/*.md
      file's frontmatter (opening delimiter, closing delimiter, valid YAML,
      decodes to a mapping) and refuses the commit otherwise. Installed
      locally at .git/hooks/pre-commit for this session. Verified it allows
      a valid file and rejects a deliberately truncated one.
  - at: '2026-09-22T00:00:02Z'
    actor: {role: executor, session: user}
    kind: repair_attempted
    note: >-
      Added a `make install-hooks` target running
      `git config core.hooksPath scripts/git-hooks`, committable so every
      clone has the one command to opt in. Not run automatically by the
      agent — updating git config is outside agent authority in this
      environment regardless of request; the owner runs it per clone.
      Local hook copy at .git/hooks/pre-commit from the prior repair still
      works standalone until `make install-hooks` is run.
---

# I027: Frontmatter parse errors reach the repo before anything catches them

## Summary

Nothing validates a `.savepoint/**/*.md` file's YAML frontmatter before a
commit lands. A hand edit can drop the closing `---` delimiter or otherwise
malform the block, and the first anyone hears about it is a board or doctor
parse error, sometime later, disconnected from the edit that caused it.

## Evidence

- I026's `deferred` history entry was added without its closing `---`
  delimiter, producing `parse error ... no closing frontmatter delimiter
  found` reported against the board.
- `DATA-01` in Guardrails.md requires parsing to preserve frontmatter on
  program-driven rewrite, but nothing covers hand/agent edits made outside
  the program's own writers.
- No `.githooks` directory, no `core.hooksPath`, and no Makefile target
  existed before this repair.

## Proof Needed

- Confirm the installed pre-commit hook (`scripts/git-hooks/pre-commit` +
  `check_frontmatter.py`) is the accepted shape, or replace it with a Go
  implementation reusing `internal/data`'s own frontmatter split if that is
  preferred over a Python dependency.
- Decide whether every clone needs to install the hook manually (documented
  step) or whether `core.hooksPath` should be committed/scripted so it's
  automatic — this repair only installed it for the current session's local
  clone.
- Consider whether CI should run the same check independently of local
  hooks, since a hook can be skipped (`--no-verify`) or never installed.
- Reconcile with Guardrails.md — this may warrant its own rule (e.g.
  `DATA-02`) rather than staying implicit tooling.
