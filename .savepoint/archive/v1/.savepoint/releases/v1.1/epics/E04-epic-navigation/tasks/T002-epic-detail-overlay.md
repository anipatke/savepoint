---
id: E04-epic-navigation/T002-epic-detail-overlay
status: done
objective: "Add Epic Detail overlay showing E##-Detail.md content, triggered by Enter on a focused epic"
depends_on:
    - E04-epic-navigation/T001-sidebar-focusable-navigation
---

# T002: Epic Detail Overlay

## Acceptance Criteria

- Pressing `Enter` on a focused epic in the sidebar opens a centered overlay titled "EPIC DETAIL"
- The overlay content is read from the epic's `E##-Detail.md` file on disk
- If the file exists, it shows: epic title, purpose, definition of done, and the components table (rendered as text)
- If the file does not exist, the overlay shows an informative "(no detail available)" message
- The overlay is height-capped at ~70% of terminal height with scroll indicators (↑/↓) for overflow
- `↑`/`k` and `↓`/`j` scroll the detail content (reuses existing scroll indicator pattern from `detail.go`)
- `PgUp`/`PgDown` scroll in page increments
- `Esc` or `q` closes the overlay and returns to epic panel focus
- The overlay uses `styles.DetailOverlay` for consistent visual styling with task detail

## Implementation Plan

- [x] Add `OverlayEpicDetail OverlayType = "detail-epic"` to `model.go`
- [x] Add `EpicDetailOffset int` to `Model` in `model.go`
- [x] Add `EpicDetailContent string` to `Model` in `model.go` (cached content of the detail file)

- [x] In `epic_panel.go` — add `RenderEpicDetail(epicID, content string, overlayW, maxHeight, offset int) string`:
  - Parse frontmatter lines from the markdown content (lines between `---` markers) for metadata display
  - Render header: epic ID and title from markdown heading
  - Render body: the rest of the markdown content as plain text (no markdown rendering)
  - Use `visibleDetailLines` from `detail.go` for height-capped scrolling
  - Reuse `styles.DetailOverlay`, `styles.ColumnTitleFocused`, `styles.CardMeta`, `styles.ColumnTitle`

- [x] In `update.go` — when `m.EpicPanelFocus && enter`:
  - Derive detail file path: `{m.Root}/releases/{m.SelectedRelease}/epics/{epicSlug}/{shortID(epicSlug)}-Detail.md`
  - Read the file content with `os.ReadFile`
  - If error (file missing, etc.), set `m.EpicDetailContent = "(no detail available)"`
  - Set `m.Overlay = OverlayEpicDetail`, `m.EpicDetailOffset = 0`
  - Cache the content in `m.EpicDetailContent`

- [x] In `updateOverlay` in `update.go` — match `OverlayEpicDetail` and handle:
  - `esc`/`q` → `OverlayNone` (exit back to epic panel focus)
  - `up`/`k` → decrement `EpicDetailOffset` if > 0
  - `down`/`j` → increment `EpicDetailOffset`
  - `pgup`/`pgdown` → page scroll using existing `detailPageSize()`

- [x] In `view.go` — match `OverlayEpicDetail` and call `RenderEpicDetail(...)` with the cached content
- [x] Reuse `visibleDetailLines`, `clampDetailOffset`, `detailRow`, `WrapText` from `detail.go` — no need to copy
- [x] Run `make build && make test` to verify no regressions

## Context Log

Files read:
- `internal/board/model.go`
- `internal/board/update.go`
- `internal/board/view.go`
- `internal/board/epic_panel.go`
- `internal/board/detail.go`
- `internal/board/card.go` — `shortID` helper
- `.savepoint/releases/v1.1/epics/E01-tui-optimisation/E01-Detail.md` — example detail file

Estimated input tokens: ~2800

Quality gate: `go build` ✓, `go test ./...` ✓ (all 3 packages pass)

Notes:
- `detail.go` functions `visibleDetailLines`, `clampDetailOffset`, `WrapText` are unexported — `epic_panel.go` is in the same `board` package so they're accessible
- Content loading is synchronous (simple file read) — no async command needed
- `shortID` extracts "E04" from "E04-epic-navigation" via first `-` split
- Updated `TestUpdate_epicPanelEnterSelectsFocusedEpic` → `TestUpdate_epicPanelEnterOpensDetailOverlay` to reflect T002 behavior change
- Added 7 new tests covering overlay open, scroll, pgup/pgdown, esc close, view rendering, no-content fallback
