---
type: epic-design
status: planned
---

# E28: Parallel Worktree Launch

## Purpose

Create isolated Git worktrees for an eligible task group and reuse the interactive launcher to start one scoped builder in each worktree.

## What this epic adds

- A testable Git worktree service with clean-checkout and collision preflight.
- Ignored coordinator and per-worktree runtime manifests that do not replace durable Savepoint state.
- Board actions for eligible group launch, active worktree visibility, and duplicate-dispatch prevention.
- Explicit, conservative cleanup that never removes dirty or mismatched worktrees.

## Components and files

| Module | Purpose |
|--------|---------|
| `internal/worktree` | Own Git command execution, branch/path naming, manifests, and cleanup validation |
| `internal/launcher` | Launch the configured builder in each worktree with a scoped prompt |
| `internal/board` | Expose eligible parallel actions and active runtime state |
| `.savepoint/runtime` | Hold ignored operational coordination files |

## Architectural delta

Savepoint gains an optional Git process boundary beside the launcher process boundary. The primary checkout remains the coordinator and router owner; each worker receives an ignored launch-scope file and edits only its own branch, task file, and declared source files.

## Boundaries

**In scope:**
- Clean-checkout preflight
- Deterministic worktree and branch creation
- Runtime manifests and scoped prompts
- Interactive batch launch and explicit cleanup

**Out of scope:**
- Agent process monitoring or cancellation
- Automatic commit, merge, rebase, conflict resolution, or branch deletion
- Parallel defects, audits, or subtasks
- OS-level sandbox enforcement

## Quality gates

- Worktree tests use temporary Git repositories and no paid agent calls.
- Board tests use fake worktree and launcher services.
- Dirty or mismatched worktrees are never removed automatically.
- `go test ./internal/worktree ./internal/launcher ./internal/board` passes.
- `make build && make test` passes.

## Open decisions

None. Runtime manifests are advisory operational state and are rebuilt or diagnosed against Git reality.
