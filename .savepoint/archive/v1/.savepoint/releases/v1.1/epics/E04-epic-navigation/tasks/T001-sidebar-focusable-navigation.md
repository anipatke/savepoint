---
id: E04-epic-navigation/T001-sidebar-focusable-navigation
status: done
objective: "Make the epic sidebar focusable with up/down navigation to select epics"
depends_on: []
---

# T001: Sidebar Focusable Navigation

## Acceptance Criteria

- Left arrow (`←`/`h`) from the `Planned` column moves keyboard focus to the epic panel (when panel is visible at ≥120 width)
- In epic panel focus mode, `↑`/`k` moves the cursor up through the epic list (clamped to bounds)
- In epic panel focus mode, `↓`/`j` moves the cursor down through the epic list (clamped to bounds)
- The focused epic entry in the sidebar is highlighted with the same focused style used in task columns
- `Enter` in epic panel focus mode selects the focused epic: tasks are filtered to that epic, and focus remains in the panel
- Right arrow (`→`/`l`) from epic panel focus exits back to the `Planned` column with `FocusedTask=0`
- The existing `E` key dropdown overlay still works unchanged
- When no epics exist in the current release, the panel shows no navigable cursor
- All navigation is blocked when any overlay is open (existing behavior)

## Implementation Plan

- [x] Add `EpicPanelFocus bool` and `EpicPanelCursor int` to `Model` in `model.go`
- [x] Add `EpicPanelFocus` tracking to `refreshEpicsForRelease()` — clamp cursor and clear focus if epics list becomes empty
- [x] In `update.go` — before the `switch msg.String()` block, check `m.Overlay == OverlayNone && m.EpicPanelFocus` and route to a new `updateEpicPanel` method
- [x] In `updateEpicPanel` — handle `up`/`k` (decrement cursor), `down`/`j` (increment cursor), `enter` (call `m.selectEpicPanelEpic()`), `right`/`l` (exit to Planned column), `left`/`h` (no-op or wrap)
- [x] In `update.go` main switch — handle `left`/`h` when `FocusedColumn == Planned` by setting `EpicPanelFocus = true` (only when epic panel is visible, i.e. `CalculateLayout(m.Width, m.Height).EpicPanelVisible`)

- [x] Add `m.selectEpicPanelEpic()` method: set `m.SelectedEpic = m.Epics[m.EpicPanelCursor]`, reset `m.FocusedTask = 0`, `m.DetailOffset = 0`, call `m.refreshTasks()`, call `m.ensureFocusedTaskVisible()`, write router state
- [x] Update `epic_panel.go` — update `RenderEpicSidebar` signature to accept `focus bool, cursor int` params; when focused, render the item at cursor with `styles.TaskItemFocused` and `epicActiveMarker`; when not focused, keep current rendering with `SelectedEpic` marker
- [x] Update `view.go` — pass `m.EpicPanelFocus` and `m.EpicPanelCursor` to `RenderEpicSidebar` in `renderEpicPanel`
- [x] Add `m.epicPanelPageSize()` helper returning `m.Height / 2` (or similar) for eventual PgUp/PgDown support
- [x] Run `make build && make test` to verify no regressions

## Context Log

Files read:
- `.savepoint/router.md`
- `agent-skills/savepoint-build-task/SKILL.md`
- `agent-skills/ink-tui-design/SKILL.md`
- `.savepoint/visual-identity.md`
- `.savepoint/releases/v1.1/epics/E04-epic-navigation/E04-Detail.md`
- `.savepoint/releases/v1.1/epics/E04-epic-navigation/tasks/T001-sidebar-focusable-navigation.md`
- `internal/board/model.go`
- `internal/board/update.go`
- `internal/board/view.go`
- `internal/board/epic_panel.go`
- `internal/board/layout.go`
- `internal/board/board.go`
- `internal/board/epic_panel_test.go`
- `internal/board/update_test.go`
- `internal/board/model_test.go`
- `internal/board/view_test.go`

Estimated input tokens: 12000

Notes:
- Model's `EpicCursor` field (used in dropdown) is distinct from new `EpicPanelCursor`
- The epic panel visibility is determined by `layout.EpicPanelVisible` in `CalculateLayout`
- Focused board tests: `go test ./internal/board` passed.
- Required quality gate `make build && make test`: failed to start because `make` is not installed in this shell.
- Equivalent quality gates: `go build -o savepoint main.go` passed; `go test ./...` passed.
- Follow-up fix after manual review: epic panel focus now renders with the focused panel border and focused title so entering the sidebar is visually apparent even when the cursor starts on the already selected epic. Re-ran `go test ./internal/board`, `go build -o savepoint main.go`, and `go test ./...`; all passed.
- Follow-up input-routing fix: global `q`, `e`, `r`, and `?` keys now run before epic-panel-specific key handling, so quit and dropdown overlays still work when the epic panel is focused. Added regression tests for `q`, `e`, and `r` from epic panel focus. Re-ran `go test ./internal/board`, `go build -o savepoint main.go`, and `go test ./...`; all passed.
