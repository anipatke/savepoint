---
id: E17-defect-workflow-tui/T011-task-defect-count-marker
status: done
objective: Show related defect counts on task cards instead of a single defect number
depends_on: [E17-defect-workflow-tui/T005-defect-detail-and-markers, E17-defect-workflow-tui/T009-defect-warning-glyph]
---

# T011: Task Defect Count Marker

## Problem

Task cards currently show the first matching related defect ID only. That hides how many defects are attached to the task and gives no quick sense of whether they are still open or already resolved.

## Context Files

- `internal/board/defect_detail.go` - related-defect matching and marker formatting
- `internal/board/view.go` - task marker aggregation and board rendering
- `internal/board/card.go` - task card layout and width constraints
- `internal/board/defect_detail_test.go` - related-defect marker tests
- `internal/board/card_test.go` - card width and marker rendering tests

## Acceptance Criteria

- [x] Task cards with related defects show a compact count marker instead of a single defect ID
- [x] The marker reports open and resolved counts for all defects referencing the task
- [x] `open` includes defects with `status: open` and `status: in_progress`
- [x] `resolved` includes only defects with `status: resolved`
- [x] Matching still works for both full task IDs and short task IDs
- [x] Cards keep their width behavior and omit the marker first when space is tight
- [x] Tasks with no related defects do not show a marker
- [x] Existing defect detail rendering remains unchanged
- [x] Tests cover aggregated counts, status grouping, ID matching, and narrow-width omission
- [x] `make build && make test` passes

## Implementation Plan

- [x] Replace first-match related-defect marker logic with aggregation over all matching defects
- [x] Count unresolved defects as open and resolved defects separately
- [x] Update task marker formatting to render compact open/resolved counts
- [x] Preserve narrow-card behavior by omitting the marker before truncating required task content
- [x] Add or update tests for multiple matching defects, short-ID matching, and width limits
- [x] Run `make build && make test`

## Context Log

- Files read: defect_detail.go, view.go, card.go, defect_detail_test.go, card_test.go, data/defect.go
- Files edited: internal/board/defect_detail.go, internal/board/defect_detail_test.go
- Token estimate: ~6k
- Quality gates: make build && make test — all pass
