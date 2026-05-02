---
type: epic-design
status: planned
---

# Epic E05: Tasking Permissions & Router Updates

## Purpose

Restrict agent task status mutation to `in_progress` only (no `done`, no retreat). Add an explicit TUI hotkey (`m`) to update the router to the current task, decoupling router updates from navigation. The user retains full control over `done` transitions and retreats via the board.

## Definition of Done

- AGENTS.md updated: agents set `status: in_progress` when starting implementation; only user can set `done` or retreat
- Router.md updated: document that agent sets router task when starting work via explicit `m` action
- Design.md Section 4 updated: status model distinguishes agent vs user capabilities
- TUI `m` hotkey implemented: sets router's `task` field to the currently focused task; detects if it's the last uncompleted task in the epic → sets `audit-pending` state
- Help overlay shows `m` key binding
- `make build && make test` passes

## Components and files

| Path | Purpose |
|------|---------|
| `AGENTS.md` | Remove "set status: done" from Task Completion Protocol; add `in_progress` guidance |
| `.savepoint/router.md` | Add agent router-update guidance in task-building state |
| `.savepoint/Design.md` | Revise Section 4 status model: agent `in_progress` only, user `done`/retreat only |
| `.savepoint/releases/v1.1/v1.1-PRD.md` | Add E05 to epic breakdown |
| `.savepoint/router.md` | Update current state for E05/T001 |
| `internal/board/update.go` | Add `m` key handler for explicit router update |
| `internal/board/model.go` | Add `writeRouterTask()` method to update router's `task` field and derive state |
| `internal/board/help.go` | Add `m` key to keyboard shortcuts |
| `internal/data/router.go` | No changes needed (WriteRouterState already supports all fields) |
| `internal/data/write.go` | No changes needed |

## Architectural notes

- The `m` hotkey reads the currently focused task, not the epic panel cursor
- Router state derivation from task:
  - Task is the last uncompleted task in its epic → `state: audit-pending`, `task: <empty>`
  - Otherwise → `state: task-building`, `task: <full-task-id>`
  - `release` and `epic` are always set from the focused task
- The `m` key does nothing when no task is focused, or when focused on a `done` task
- The `m` key is independent of navigation — user can browse freely without router drift
- No code enforcement against agents writing `done`. This is a documentation-only trust boundary. AGENTS.md is the contract.
