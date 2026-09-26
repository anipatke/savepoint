---
id: I-069
title: The migrate preview says ready when apply will refuse the working tree
type: defect
status: open
source:
  kind: report
  actor: {role: owner, session: user}
  at: '2026-09-26T03:34:20Z'
severity: medium
history:
  - at: '2026-09-26T03:34:20Z'
    actor: {role: owner, session: user}
    kind: observed
    note: >-
      Migrating galaxy with 2.0.5: the preview printed "Status: ready", then
      --apply refused twice, first for untracked and modified audit files,
      then for git-ignored screenshot folders under planned paths.
---

# I-069: The migrate preview says ready when apply will refuse the working tree

## Summary

`RunCommand` checks `ensureCleanGitTree` only on the apply path. The preview
never runs it, so it reports "Status: ready" even when apply will refuse
modified, untracked, or ignored files at planned write or removal paths. The
owner finds out one refusal at a time.

## Evidence

- `internal/migrate/command.go`: `ensureCleanGitTree(root, plan, opts.RunGit)`
  runs only after `!opts.Write` returns.
- galaxy: the preview said ready; apply then refused `!! .savepoint/releases/v1/epics/E10-natural-rendering-exploration/evidence/`
  and `!! .savepoint/releases/v1/epics/E12-luminous-galaxy-renderer/`
  (285 MB of git-ignored PNGs), after an earlier refusal for uncommitted
  audit records.

## Proof Needed

- The preview runs the same working-tree check as apply and, when it would
  fail, reports "Status: blocked" listing every offending path grouped as
  modified, untracked, or ignored, with the same repair advice apply gives.
- A clean tree still previews as ready; a project outside Git is reported
  as apply does.
- Tests cover modified, untracked, and ignored paths in the preview;
  `make build && make test-fast` and `make test-full` pass.
