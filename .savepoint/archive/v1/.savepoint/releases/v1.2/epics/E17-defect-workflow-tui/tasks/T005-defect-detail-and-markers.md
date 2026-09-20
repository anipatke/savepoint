---
id: E17-defect-workflow-tui/T005-defect-detail-and-markers
status: done
objective: Add defect detail rendering and lightweight related-defect markers on task cards
depends_on: [E17-defect-workflow-tui/T001-defect-data-model, E17-defect-workflow-tui/T004-defects-overlay]
complexity_tier: medium
complexity_reason: "New detail overlay plus width-sensitive card markers across two rendering contexts"
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

- [ ] Add a defect detail overlay renderer modeled on task detail
- [ ] Parse and render defect markdown sections by heading
- [ ] Add optional related-defect marker data to card rendering inputs
- [ ] Keep card rendering stable at narrow widths by omitting markers first
- [ ] Add detail and card tests for defect-specific rendering
- [ ] Run `make build && make test`

## Context Log

- Files read:
- Files edited:
- Token estimate:
- Quality gates:

