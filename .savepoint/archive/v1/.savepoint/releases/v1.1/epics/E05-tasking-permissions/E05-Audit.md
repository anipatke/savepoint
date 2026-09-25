---
type: audit-findings
audited: 2026-05-02
---

# Audit Findings: E05 Tasking Permissions

## Main Findings

E05 verified the tasking-permissions shift from the planned `m` hotkey to the final `p` priority hotkey, along with explicit router priority writes that stay in `task-building` and leave audit handoff as an explicit workflow step.

The audit found documentation drift rather than a core behavior failure: live implementation and tests used `p`, while several planning documents still referred to `m`. The file-specific reconciliation blocks are retained under Proposed Changes for apply/close testing.

## Code Style Review


- [ ] One job per file
- [ ] One-sentence functions
- [ ] Test branches
- [ ] Types are documentation
- [ ] Build, don't speculate
- [ ] Errors at boundaries
- [ ] One source of truth
- [ ] Comments explain WHY
- [ ] Content in data files
- [ ] Small diffs

## Proposed Changes

### Target File

`.savepoint/Design.md`

### Replace

```md
- Router updates are explicit TUI actions: after setting a task to `in_progress`, the agent prompts the user to press `m` in the board to update `.savepoint/router.md` to the focused task. Navigation alone must not change router state.
```

### With

```md
- Router updates are explicit TUI actions: after setting a task to `in_progress`, the agent prompts the user to press `p` in the board to mark the focused task as router priority. Navigation alone must not change router task priority.
```

---

### Target File

`.savepoint/Design.md`

### Replace

```md
**Keybindings:** arrow/vim navigation, enter advances, backspace retreats, r/R refreshes, a/A exits toward audit review when proposals exist, q quits.
```

### With

```md
**Keybindings:** arrow/vim navigation, enter opens focused task detail, space advances, backspace retreats, `p` marks the focused non-done task as router priority, `r`/`R` opens release selection or refreshes where supported, `?` opens help, and `q` quits or closes overlays.
```

---

### Target File

`AGENTS.md`

### Replace

```md
3. After setting `in_progress`, press `m` in the TUI to update the router
```

### With

```md
3. After setting `in_progress`, press `p` in the TUI to mark the focused task as router priority
```

---

### Target File

`.savepoint/router.md`

### Replace

```md
**Next:** When starting work, set task `status: in_progress` and press `m` in TUI to update router task. Execute plan, tick checkboxes, run quality gates, update router to next task or `audit-pending`. Stop.
```

### With

```md
**Next:** When starting work, set task `status: in_progress` and press `p` in TUI to mark the focused task as router priority. Execute plan, tick checkboxes, run quality gates, update router to next task or `audit-pending`. Stop.
```

---

### Target File

`.savepoint/releases/v1.1/epics/E05-tasking-permissions/E05-Detail.md`

### Replace

```yaml
status: planned
```

### With

```yaml
status: audited
```

---

### Target File

`.savepoint/releases/v1.1/epics/E05-tasking-permissions/E05-Detail.md`

### Replace

```md
- TUI `m` hotkey implemented: sets router's `task` field to the currently focused task; detects if it's the last uncompleted task in the epic → sets `audit-pending` state
- Help overlay shows `m` key binding
```

### With

```md
- TUI `p` priority hotkey implemented: sets router's `task` field to the currently focused non-done task and keeps router state as `task-building`
- Help overlay shows `p` key binding
```

---

### Target File

`.savepoint/releases/v1.1/epics/E05-tasking-permissions/E05-Detail.md`

### Replace

```md
- The `m` hotkey reads the currently focused task, not the epic panel cursor
- Router state derivation from task:
  - Task is the last uncompleted task in its epic → `state: audit-pending`, `task: <empty>`
  - Otherwise → `state: task-building`, `task: <full-task-id>`
  - `release` and `epic` are always set from the focused task
- The `m` key does nothing when no task is focused, or when focused on a `done` task
- The `m` key is independent of navigation — user can browse freely without router drift
- No code enforcement against agents writing `done`. This is a documentation-only trust boundary. AGENTS.md is the contract.
```

### With

```md
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
```

---

### Target File

`.savepoint/releases/v1.1/epics/E05-tasking-permissions/tasks/T005-update-help-overlay.md`

### Replace

```yaml
objective: "Add `m` key binding to the TUI help overlay"
```

### With

```yaml
objective: "Add `p` priority key binding to the TUI help overlay"
```

---

### Target File

`.savepoint/releases/v1.1/epics/E05-tasking-permissions/tasks/T005-update-help-overlay.md`

### Replace

```md
# T005: Update Help Overlay

## Acceptance Criteria

- Help overlay (`?` key) shows `m` key binding: "m: update router"
- Key binding appears with consistent styling in the shortcuts list
- No other help entries modified

## Implementation Plan

- [ ] Add `helpRow("m", "update router")` to `RenderHelp` in `help.go`
- [ ] Run `make build && make test` to verify
```

### With

```md
# T005: Update Help Overlay

## Acceptance Criteria

- Help overlay (`?` key) shows `p` key binding for router priority
- Key binding appears with consistent styling in the shortcuts list
- No unrelated help entries modified

## Implementation Plan

- [x] Add the `p` priority shortcut to `RenderHelp` in `help.go`
- [x] Verify with `go build -o savepoint.exe main.go` and `go test ./...`
```

---
