---
id: E03-ui-visual-refinement/T005-unify-status-glyphs
status: done
objective: "Unify icon/glyph determination across task cards and epic sidebar via a shared statusGlyph helper with backward compatibility"
depends_on: []
---

# T005: Unify Task Status Glyph Determination

## Acceptance Criteria

- A new `status` field is added to the `data.Task` struct (one of: `planned`, `in_progress`, `done`, `audited`)
- A shared `statusGlyph(status string) (styled string)` helper exists in `internal/board/status.go` that maps status strings to the same glyph+style used by the epic sidebar
- `RenderCard` in `card.go` uses the unified helper when `Task.Status` is set, otherwise falls back to the existing Column+Stage logic with identical behavior
- The epic sidebar's `epicSidebarGlyph` function is refactored to call the shared `statusGlyph` helper instead of its own switch
- All existing tests pass without modification
- New tests cover the explicit `Status` path for each status value
- No changes to existing task files are required
- `make build && make test` pass

## Implementation Plan

- [x] Read `internal/data/task.go` — study Task struct and existing columns/stage types
- [x] Read `internal/board/card.go` — study current RenderCard glyph determination path
- [x] Read `internal/board/epic_panel.go` — study epicSidebarGlyph for extraction opportunity
- [x] Edit `internal/data/task.go` — add `Status string` field to Task struct and Status constants
- [x] Create `internal/board/status.go` — shared `statusGlyph(status string) string` function with the unified glyph mapping
- [x] Edit `internal/board/epic_panel.go` — refactor `epicSidebarGlyph` to call the shared `statusGlyph`
- [x] Edit `internal/board/card.go` — update `RenderCard` to use `statusGlyph` when Status is set, else existing backward-compat path
- [x] Edit `internal/board/card_test.go` — add tests for explicit Status path (planned/in_progress/done/audited)
- [x] Run `make build && make test` to verify no regressions

## Context Log

Files read:
- `.savepoint/router.md`
- `.savepoint/releases/v1.1/epics/E03-ui-visual-refinement/E03-Detail.md`
- `.savepoint/releases/v1.1/epics/E03-ui-visual-refinement/tasks/T005-unify-status-glyphs.md`
- `agent-skills/savepoint-build-task/SKILL.md`
- `agent-skills/ink-tui-design/SKILL.md`
- `.savepoint/visual-identity.md`
- `internal/data/task.go`
- `internal/data/parser.go`
- `internal/data/write.go`
- `internal/data/task_test.go`
- `internal/board/card.go`
- `internal/board/epic_panel.go`
- `internal/styles/styles.go`
- `internal/board/card_test.go`
- `internal/board/epic_panel_test.go`
- `.savepoint/releases/v1.1/epics/E04-epic-navigation/tasks/T003-epic-status-glyphs.md`

Estimated input tokens: 11000

Notes:
- Option A approach: backward compatible when `Status` is unset
- The unified glyph mapping matches epic sidebar: planned→○, in_progress→▶, done→◉, audited→✓
- Old `Column`+`Stage` path reproduces exact current behavior (▣/◇/◆ with column-level overrides)
- `go test ./internal/board ./internal/data` initially hit a sandbox AppData build-cache permission error; rerun outside the sandbox passed.
- Literal `make build && make test` could not run because `make` is not installed in this Windows shell.
- Equivalent quality gates passed: `go run ./internal/buildtool build`; `go test ./...`.

## Drift Notes

- Drift: `internal/board/status.go` added, not yet in Codebase Map.
