---
id: E17-defect-workflow-tui/T009-defect-warning-glyph
status: done
objective: Standardize defect indicators on a terminal-safe warning glyph
depends_on: [E17-defect-workflow-tui/T005-defect-detail-and-markers]
complexity_tier: low
complexity_reason: "Glyph standardization is a narrow rendering constant and expectation update"
---

# T009: Defect Warning Glyph

## Problem

Defects need a recognizable visual marker, but emoji bug icons can render with inconsistent width across terminals. The TUI should use a terminal-safe defect glyph that communicates warning/repair work without destabilizing layout.

## Context Files

- `internal/board/defect_overlay.go` - defect row severity and list rendering
- `internal/board/defect_detail.go` - related-defect task-card marker helper
- `internal/board/card.go` - task card marker layout constraints
- `internal/board/defect_overlay_test.go` - defect overlay rendering tests
- `internal/board/card_test.go` - task card marker rendering tests

## Acceptance Criteria

- [x] Defect indicators use `⚠` as the default terminal-safe glyph
- [x] Related task-card markers use a compact warning marker such as `⚠ D003`
- [x] Defect overlay rows remain width-stable and readable
- [x] Emoji bug glyphs are not used in core TUI rendering
- [x] Tests cover the updated defect marker output
- [x] `make build && make test` passes

## Implementation Plan

- [x] Add or reuse a single defect glyph constant in board rendering code
- [x] Update related-defect task-card marker formatting to use `⚠`
- [x] Update defect overlay severity/row formatting if needed to consistently signal defects
- [x] Update focused tests for the new marker text
- [x] Run `make build && make test`

## Context Log

- Files read: `internal/board/card.go`, `internal/board/defect_detail.go`, `internal/board/defect_overlay.go`, `internal/board/defect_detail_test.go`, `internal/board/defect_overlay_test.go`, `internal/board/card_test.go`
- Files edited: `internal/board/card.go`, `internal/board/defect_detail.go`, `internal/board/defect_detail_test.go`
- Token estimate: ~8
- Quality gates: `make build && make test` — pass
