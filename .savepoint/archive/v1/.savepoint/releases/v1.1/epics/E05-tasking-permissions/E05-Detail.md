---
type: epic-design
status: audited
---

# Epic E05: Tasking Permissions & Router Updates

## Purpose

Restrict agent task status mutation to `in_progress` only (no `done`, no retreat). Add an explicit TUI priority hotkey (`p`) to update the router to the focused non-done task, decoupling router task priority from navigation. The user retains full control over `done` transitions and retreats via the board.

## Definition of Done

- AGENTS.md updated: agents set `status: in_progress` when starting implementation; only user can set `done` or retreat
- Router.md updated: document that agent sets router task when starting work via explicit `p` action
- Design.md Section 4 updated: status model distinguishes agent vs user capabilities
- TUI `p` priority hotkey implemented: sets router's `task` field to the currently focused non-done task and keeps router state as `task-building`
- Help overlay shows `p` key binding
- `make build && make test` passes

## Components and files

| Path | Purpose |
|------|---------|
| `AGENTS.md` | Remove "set status: done" from Task Completion Protocol; add `in_progress` guidance |
| `.savepoint/router.md` | Add agent router-update guidance in task-building state |
| `.savepoint/Design.md` | Revise Section 4 status model: agent `in_progress` only, user `done`/retreat only |
| `.savepoint/releases/v1.1/v1.1-PRD.md` | Add E05 to epic breakdown |
| `.savepoint/router.md` | Update current state for E05/T001 |
| `internal/board/update.go` | Add `p` key handler for explicit router priority update |
| `internal/board/model.go` | Add `writeRouterTask()` method to update router release, epic, task, state, and next action |
| `internal/board/help.go` | Add `p` key to keyboard shortcuts |
| `internal/data/router.go` | No changes needed (WriteRouterState already supports all fields) |
| `internal/data/write.go` | No changes needed |

## Architectural notes

- The `p` hotkey reads the currently focused task, not the epic panel cursor.
- Router priority writes always set `state: task-building`, `task: <full-task-id>`, `release`, `epic`, and `next_action` from the focused task.
- The `p` key does nothing when no task is focused, when an overlay is open, or when focused on a `done` task.
- The `p` key is independent of navigation — user can browse freely without router task drift.
- Audit state is not inferred by the priority key; audit handoff remains an explicit workflow/router closeout action.
- No code enforcement prevents agents from writing `done`. This is a documentation-only trust boundary. AGENTS.md is the contract.

## Implemented as

- `AGENTS.md`, `templates/project/AGENTS.md`, and phase skill guides now describe the three-status canon and agent/user task-status ownership boundary.
- `.savepoint/router.md` was simplified around the current state machine and task-building handoff.
- `.savepoint/Design.md` Section 4 records that agents may only set `in_progress`; user-only operations are `done` and retreat.
- `internal/board/model.go` adds `writeRouterTask(task data.Task)`, which reads the current router, writes release/epic/task/state/next_action via `data.WriteRouterState`, and updates in-memory router fields.
- `internal/board/update.go` handles `p` before normal navigation when no overlay is open, no-ops for done tasks, and writes user-visible status messages.
- `internal/board/help.go` and `internal/board/view.go` advertise `p: Priority` and render priority/status feedback in the footer.
- `internal/board/update_test.go`, `help_test.go`, `view_test.go`, `card_test.go`, and transition tests cover priority routing, done-task no-op, overlay no-op, help text, footer status feedback, and transition messages.
- Implementation deviation: the original design specified `m` and last-task `audit-pending` derivation. The final behavior uses `p` and keeps priority writes in `task-building` so audit handoff remains explicit.
