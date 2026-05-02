---
id: E06-atari-noir-layout/T010-auto-refresh-watcher
status: done
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

- [x] Add `github.com/fsnotify/fsnotify` to `go.mod`
- [x] Create `internal/board/watch.go`:
  - Define `fileChangeMsg struct{}` and `reloadMsg struct{ tasks []data.Task }`
  - `watchFiles(w *fsnotify.Watcher) tea.Cmd` — blocks on watcher.Events, drains channel for 100ms, emits single `fileChangeMsg`
  - `reloadTasks(root string) tea.Cmd` — calls `loadAllTasks(root)` and emits `reloadMsg{tasks}`
  - `loadAllTasks(root string) ([]data.Task, error)` — extracted from `board.go`, reuses `Discover` + `Parser`
- [x] Edit `internal/board/model.go`:
  - Add `Watcher *fsnotify.Watcher` field to `Model`
  - `Init()` — returns `watchFiles(m.Watcher)` if watcher non-nil, else nil
- [x] Edit `internal/board/board.go`:
  - Extract task-discovery loop from `newProjectModel` into `loadAllTasks(root string) ([]data.Task, error)`
  - `newProjectModel` calls `loadAllTasks` + `newWatcher`, sets `model.Watcher`
  - `loadEpicTasks` now sets `task.Path` and `task.Mtime` from stat
- [x] Edit `internal/board/update.go`:
  - Handle `fileChangeMsg` → return `reloadTasks(m.Root)` cmd
  - Handle `reloadMsg` → swap `m.AllTasks = msg.tasks`, call `m.refreshTasks()`, return `watchFiles(m.Watcher)` to re-subscribe
  - On quit (`q`/`ctrl+c`), close `m.Watcher` before `tea.Quit`
  - `space`/`backspace` now call `data.WriteTaskStatus` to persist to disk
- [x] `internal/data/task.go`: added `Path string` and `Mtime time.Time` (yaml:"-") to Task
- [x] Run `go build ./...` and `go test ./...` to verify no regressions

## Context Log

Files read:
- `internal/board/model.go`
- `internal/board/board.go`
- `internal/board/update.go`
- `internal/board/transitions.go`
- `internal/board/watch.go` (created)
- `internal/data/write.go`
- `internal/data/task.go`
- `internal/data/parser.go`
- `.savepoint/AGENTS.md`

Estimated input tokens: 2200

Quality gates: `go build ./...` PASS, `go test ./...` PASS (board: 0.317s, data: 0.343s)

Notes:
- fsnotify v1.10.0 has no recursive AddWith option; used `filepath.WalkDir` + `w.Add` per subdir instead.
- `data.Task` gained `Path` and `Mtime` (yaml:"-") fields populated in `loadEpicTasks`.
- Advance/Retreat disk writes live in `update.go`, not `transitions.go` (keeps transitions pure).
- Watcher created in `newProjectModel`; Init() subscribes if non-nil; nil watcher = test-safe.
- Follow-up fix (2026-05-02): `newWatcher` receives `root` as the `.savepoint` directory, so it watches `root/releases` directly rather than `.savepoint/.savepoint/releases`. Watch startup errors now fail model creation instead of being silently ignored.
- Follow-up fix (2026-05-02): file-change reloads refresh task data plus release/epic indexes, and create events add watches for newly-created directories before the debounced reload.
