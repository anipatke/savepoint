---
id: E04-board-components/T002-card
status: done
objective: "Render task card with phase glyphs, truncation, and focus styling"
depends_on: [E04-board-components/T001-column]
---

# T002: Card Component

## Acceptance Criteria

- Card shows: task ID (truncated), title (truncated), phase glyphs (▣/◇/◆).
- Fixed width, never wraps.
- Metadata dimmed.
- Focused card has orange accent border and background tint.
- Phase glyph uses phase accent color.

## Implementation Plan

- [x] Create `internal/board/card.go`.
- [x] Implement `RenderCard(task, width, focused)`.
- [x] Add truncation helper.
- [x] Apply focus and phase styling.
- [x] Write tests.
- [x] Run `go test`.

## Context Log

Files read: `internal/board/column.go`, `internal/board/column_test.go`, `internal/styles/styles.go`, `internal/styles/palette.go`, `internal/data/task.go`, `.savepoint/releases/v1/epics/E04-board-components/Design.md`, `AGENTS.md`, `.savepoint/router.md`

Estimated input tokens: ~4 000

Notes: Added `Dim`/`Dim256`/`Dim16` palette constants and `Card`, `CardFocused`, `CardMeta`, `GlyphBuild`, `GlyphTest`, `GlyphAudit` styles. `truncate` helper is package-private in `board`. All 13 tests pass. `go build ./...`, `go vet ./...`, `go test ./...` green.
