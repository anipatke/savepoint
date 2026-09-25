---
id: E04-board-components/T004-detail-overlay
status: done
objective: "Render task detail as centered overlay with board behind"
depends_on: [E04-board-components/T003-epic-panel]
---

# T004: Detail Overlay

## Acceptance Criteria

- `Enter` opens detail overlay for focused task.
- `Esc` closes overlay.
- Overlay shows: ID, title, epic, release, status, phase, description, acceptance criteria.
- Board visible and dimmed behind overlay.
- Overlay width = min(terminal width - 4, 80).

## Implementation Plan

- [x] Create `internal/board/detail.go`.
- [x] Implement `RenderDetail(task, overlayW)`.
- [x] Add board dimming (`dimLines` + `overlayOnBase` in `view.go`).
- [x] Wire Enter/Esc into Update loop.
- [x] Write tests.
- [x] Run `go test`.

## Context Log

**Files read:** model.go, update.go, view.go, epic_panel.go, card.go, layout.go, update_test.go, view_test.go, epic_panel_test.go, internal/data/task.go, internal/styles/styles.go, internal/styles/palette.go, go.mod

**Estimated input tokens:** ~8 000

**Notes:**
- `overlayOnBase` uses `charmbracelet/x/ansi.Truncate` for ANSI-aware left truncation; each base line is independently faint-wrapped so truncation doesn't inherit prior escape states.
- `overlayWidth` clamps to `min(termW-4, 80)` with a floor of 20.
- `phaseLabel` lives in `detail.go`; `focusedTask`, `overlayWidth`, `dimLines`, `overlayOnBase` live in `view.go`.
- `OverlayDetail` constant added to `model.go`; `styles.DetailOverlay` added to `styles.go`.
- No drift: all new files/exports match the Codebase Map entry for `internal/board/detail.go`.

**Quality gates:** `go build ./...` ✓ · `go vet ./...` ✓ · `go test ./...` ✓ (board: 32 tests pass)
