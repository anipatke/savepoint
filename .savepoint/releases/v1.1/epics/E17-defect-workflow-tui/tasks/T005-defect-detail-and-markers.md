---
id: E17-defect-workflow-tui/T005-defect-detail-and-markers
status: in_progress
stage: build
objective: Add defect detail rendering and lightweight related-defect markers on task cards
depends_on: [E17-defect-workflow-tui/T001-defect-data-model, E17-defect-workflow-tui/T004-defects-overlay]
---

# T005: Defect Detail And Related Markers

## Problem

The defect list gives users visibility, but repair work needs evidence: symptom, expected behavior, reproduction, impact, fix plan, acceptance criteria, and resolution notes. Related task cards should also hint when a visible task has an associated defect.

## Context Files

- `internal/board/detail.go` - task detail overlay and checklist rendering
- `internal/board/card.go` - task card line layout and router glyph behavior
- `internal/board/column.go` - task card viewport rendering
- `internal/board/view.go` - detail overlay dispatch
- `internal/board/detail_test.go` - detail rendering tests
- `internal/board/card_test.go` - card rendering tests
- `internal/styles/styles.go` - warning and semantic styles

## Acceptance Criteria

- [ ] Enter on a selected defect opens a defect detail overlay
- [ ] Defect detail renders symptom, expected behavior, reproduction, impact, fix plan, acceptance criteria, and resolution notes when present
- [ ] Long defect content wraps and scrolls with the same behavior as task detail overlays
- [ ] Defects that reference a visible task can produce a compact task-card marker such as `! D003`
- [ ] Card markers only render when width permits and never displace required task id/title content
- [ ] Existing task detail rendering remains unchanged
- [ ] Tests cover full defect detail, missing optional sections, wrapping/scrolling, and card marker width behavior
- [ ] `make build && make test` passes

## Implementation Plan

- [x] Add a defect detail overlay renderer modeled on task detail
- [x] Parse and render defect markdown sections by heading
- [x] Add optional related-defect marker data to card rendering inputs
- [x] Keep card rendering stable at narrow widths by omitting markers first
- [x] Add detail and card tests for defect-specific rendering
- [x] Run `make build && make test`

## Context Log

- Files read: detail.go, card.go, column.go, view.go, update.go, model.go, defect_overlay.go, defect.go, dependency.go, styles.go, detail_test.go, card_test.go, defect_test.go, column_test.go
- Files edited: model.go, card.go, column.go, view.go, update.go, card_test.go, column_test.go, render_policy_test.go, view_test.go, update_test.go
- Files created: defect_detail.go, defect_detail_test.go
- Token estimate: ~12k
- Quality gates: make build && make test — all pass

