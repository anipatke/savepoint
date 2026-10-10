---
id: G-001
title: Worktree Support
status: done
---

## Outcome

AI-assisted solo builders can understand which planned Tasks are safe to run concurrently with separate agents in Git worktrees, while sequential execution remains the default.

## Why

Savepoint can describe independent Tasks, but it does not persist or present a trustworthy concurrency plan. Owners must currently infer safe parallel work from dependencies and hand-write additional agent instructions.

## Success Conditions

- Objective planning records enough information to distinguish harmless shared reads from conflicting planned writes.
- Savepoint derives safe concurrent Task opportunities from dependencies and planned writes rather than relying on informal lane labels.
- The board and `savepoint resume` explain the same concurrent and sequential Task groupings in plain language.
- Every suggested concurrent Task has a complete, copyable instruction for starting a separate worktree session.
- Existing Task lifecycle states remain unchanged, and projects without concurrency metadata default safely to sequential execution.

## Boundaries

**In scope:** planning data for concurrency decisions, deterministic concurrency projection, board and resume presentation, copyable worktree prompts, workflow guidance, verification, and upgrade compatibility.

**Out of scope:** automatic worktree creation or deletion, automatic merging or conflict resolution, background Git monitoring, multi-user locking, cloud coordination, general multi-agent orchestration, and new lane lifecycle statuses.
