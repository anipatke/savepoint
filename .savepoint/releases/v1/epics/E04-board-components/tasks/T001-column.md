---
id: E04-board-components/T001-column
status: done
objective: "Render column with header, count, and task list"
depends_on: [E03-board-tui-core/T005-layout]
---

# T001: Column Component

## Acceptance Criteria

- Column renders header with status label (PLANNED / IN PROGRESS / DONE).
- Header shows task count.
- Column has bordered container.
- Focused column has accent border.
- Empty columns show "(empty)".

## Implementation Plan

- [x] Create `internal/board/column.go`.
- [x] Implement `RenderColumn(tasks, col, width, focusedTask, focused)`.
- [x] Style header and border.
- [x] Write tests.
- [x] Run `go test`.

## Context Log

**Files read:** `internal/board/view.go`, `internal/board/model.go`, `internal/styles/styles.go`, `internal/styles/palette.go`, `internal/data/task.go`, `internal/board/view_test.go`

**Estimated input tokens:** ~4 000

**Notes:** Extracted `columnTitle` and `taskLabel` from `view.go` into `column.go`. `renderColumn` method in `view.go` now delegates to `RenderColumn`. Added task count to header format `LABEL (N)` and `(empty)` placeholder for empty columns. All 7 column tests + prior view/board tests pass. `go build` and `go vet` clean.

**Quality gates:** `go build ./...` ✅ | `go vet ./...` ✅ | `go test ./...` ✅
