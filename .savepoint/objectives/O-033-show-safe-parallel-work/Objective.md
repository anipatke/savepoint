---
id: O-033
title: Show which Tasks can safely run in parallel worktrees
status: planned
depends_on: [O-032]
release: G-001
priority: medium
rank: 1
---

# O-033: Show which Tasks can safely run in parallel worktrees

## Outcome

A solo builder can see which planned Tasks may safely be handed to parallel agents in separate worktrees, why those Tasks are independent, and which work must remain sequential.

## Why

The current guidance encourages worktree lanes without persisting their plan or distinguishing shared reads from shared writes. This leaves the owner to infer concurrency safety and compose secondary agent prompts by hand.

## Success Conditions

- Sequential execution remains the default when concurrency cannot be proven.
- Task planning records exact anticipated writes separately from the executor's readable Context Files.
- Concurrent suggestions require satisfied dependencies, no dependency path between the Tasks, and disjoint planned writes.
- Semantic or shared-contract coupling that file overlap cannot express remains an explicit Task dependency.
- The concurrency plan persists outside planning chat and survives later sessions and replanning.
- The board and `savepoint resume` show the same concurrent groups, sequential constraints, and understandable reasons.
- Each suggested parallel Task has a complete, copyable instruction for starting its worktree session.
- Existing `planned`, `in_progress`, and `done` lifecycle states remain unchanged.
- Existing projects and Tasks without the new planning data remain valid and default to sequential execution.

## Architectural Considerations

`internal/data` must own the canonical concurrency decision so the board and resume cannot diverge. Context Files remain the read budget; planned writes are separate structured planning data. Dependencies remain the source of truth for ordering and semantic coupling. Presentation may describe parallel opportunities, but it must not introduce another Task lifecycle or make worktree state part of completion.

The exact persisted shape, prompt format, and board presentation remain design decisions to settle before Task detailing.

## Boundaries

**In scope:** concurrency planning metadata, canonical safety rules, board and resume presentation, copyable worktree instructions, planner/executor/checker guidance, scaffold parity, and existing-project upgrade behavior.

**Out of scope:** automatic worktree creation or deletion, branch management, automatic merges or conflict resolution, claimed/running/merged lane statuses, background repository monitoring, multi-user coordination, cloud services, and general-purpose agent orchestration.
