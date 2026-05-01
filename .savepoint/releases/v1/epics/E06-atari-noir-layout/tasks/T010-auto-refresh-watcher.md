---
id: E06-atari-noir-layout/T010-auto-refresh-watcher
status: planned
objective: "Auto-refresh the board when task files change on disk via fsnotify file watcher"
depends_on: []
---

# T010: Auto-Refresh Board on File Changes

## Acceptance Criteria

- The board detects task `.md` file changes (write, create, delete) under `.savepoint/releases/` without requiring a restart
- When a file change is detected, the board reloads all tasks from disk and re-renders the columns
- Rapid successive file saves (e.g. editor auto-save + format) are debounced — only one reload fires within a 100ms window
- The watcher does not prevent graceful shutdown (closes on `q`/`ctrl+c`)
- An in-memory status change via `space`/`backspace` persists to disk via `data.WriteTaskStatus` and triggers the file watcher naturally
- All existing keyboard navigation, overlay, and column behavior remain intact
- `go build ./...` and `go test ./...` pass cleanly

## Implementation Plan

- [ ] Add `github.com/fsnotify/fsnotify` to `go.mod`
- [ ] Create `internal/board/watch.go`:
  - Define `fileChangeMsg struct{}` and `reloadMsg struct{ tasks []data.Task }`
  - `watchFiles(w *fsnotify.Watcher) tea.Cmd` — blocks on watcher.Events, drains channel for 100ms, emits single `fileChangeMsg`
  - `reloadTasks(root string) tea.Cmd` — calls `loadAllTasks(root)` and emits `reloadMsg{tasks}`
  - `loadAllTasks(root string) ([]data.Task, error)` — extracted from `board.go`, reuses `Discover` + `Parser`
- [ ] Edit `internal/board/model.go`:
  - Add `watcher *fsnotify.Watcher` field to `Model`
  - `Init()` — create watcher, `watcher.AddRecursive(root)` on `.savepoint/releases/`, return `watchFiles(watcher)`
  - Close watcher in `Init()` if root is empty (test mode)
- [ ] Edit `internal/board/board.go`:
  - Extract task-discovery loop from `newProjectModel` into `loadAllTasks(root string) ([]data.Task, error)`
  - `newProjectModel` calls `loadAllTasks` instead of inline logic
- [ ] Edit `internal/board/update.go`:
  - Handle `fileChangeMsg` → return `reloadTasks(m.Root)` cmd
  - Handle `reloadMsg` → swap `m.AllTasks = msg.tasks`, call `m.refreshTasks()`, return `watchFiles(m.watcher)` to re-subscribe
  - On quit (`q`/`ctrl+c`), close `m.watcher` before `tea.Quit`
- [ ] Edit `internal/board/transitions.go`:
  - In `Advance` / `Retreat`, call `data.WriteTaskStatus(taskPath, task, expectedMtime)` to persist to disk so the watcher picks it up
- [ ] Run `go build ./...` and `go test ./...` to verify no regressions

## Context Log

Files read:
- `internal/board/model.go`
- `internal/board/board.go`
- `internal/board/update.go`
- `internal/board/transitions.go`
- `internal/board/watch.go` (does not exist yet)
- `internal/data/write.go`
- `.savepoint/AGENTS.md`
- `.savepoint/releases/v1/epics/E06-atari-noir-layout/Design.md`
- `.savepoint/releases/v1/epics/E06-atari-noir-layout/tasks/T009-router-priority-marker.md`

Estimated input tokens: 1800

Notes:
- The watcher watches `.savepoint/releases/` recursively. `fsnotify` v1.7+ supports `AddWith(path, fsnotify.Recursive)` on macOS/Windows; on Linux, fall back to walking and adding each subdirectory.
- The 100ms debounce in `watchFiles` collapses rapid write sequences (e.g. editor save + lsp format) into a single reload.
- `WriteTaskStatus` needs the task file path. Task struct may need a `Path` field populated at load time to make persistence possible.
