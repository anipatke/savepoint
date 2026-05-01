---
id: E06-atari-noir-layout/T009-router-priority-marker
status: planned
objective: "Highlight the router priority task with a green marker on the card and optional text highlight"
depends_on: []
---

# T009: Highlight Router Priority Task on Board

## Acceptance Criteria

- The model reads `RouterState.Task` (the `task:` field from `router.md`) and stores it as `Model.RouterTask`
- In the board view, the task card whose ID matches `RouterTask` shows a green `▣` glyph (using `TagDone`/green style) instead of the normal phase-colored glyph
- Non-priority tasks keep their existing phase-colored glyphs unchanged
- The task detail overlay shows a `"(router priority)"` label or marker when the displayed task matches `RouterTask`
- All existing keyboard interactions, column rendering, and card behavior remain intact

## Implementation Plan

- [ ] Edit `internal/board/model.go` — add `RouterTask string` field to `Model`.
- [ ] Edit `internal/board/board.go` — in `newProjectModel()`, read `routerState.Task` into `model.RouterTask`.
- [ ] Edit `internal/board/card.go` — update `RenderCard` to accept a `routerTaskID string` parameter; when `t.ID == routerTaskID`, replace the phase glyph with a green `▣` using `TagDone` style.
- [ ] Edit `internal/board/card.go` — update `RenderCard` signature and callsites (`column.go`, `view.go`).
- [ ] Edit `internal/board/detail.go` — update `RenderDetail` to accept a `routerTaskID string`; when matching, append a `"(router priority)"` green label line.
- [ ] Update all callsites of `RenderCard` and `RenderDetail` to pass `m.RouterTask`.
- [ ] Run `make build && make test` to verify no regressions.

## Context Log

Files read:
- `internal/board/model.go`
- `internal/board/board.go`
- `internal/board/card.go`
- `internal/board/column.go`
- `internal/board/detail.go`
- `internal/data/router.go`
- `internal/board/view.go`

Estimated input tokens: 1000

Notes:
