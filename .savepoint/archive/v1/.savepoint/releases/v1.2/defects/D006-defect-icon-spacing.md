---
id: v1.2/D006-defect-icon-spacing
release: v1.2
status: resolved
severity: low
title: "Defect icon spacing is too tight in the board header and task cards"
reference: E17-defect-workflow-tui/T009-defect-warning-glyph
---

# D006: Defect Icon Spacing

## Symptom

The defect icon in the top-right of the board header sits too close to the adjacent text, and the same tight spacing also appears on task cards.

## Expected Behavior

The icon should have clear spacing from nearby text in both the board header and task cards so the UI reads cleanly and the icon does not appear visually crowded.

## Reproduction

1. Open the board for release v1.2.
2. View the top-right board header area on E17.
3. Observe the defect icon placement next to the text.
4. Open a task card and observe the same icon spacing next to its text.

## Impact

This is a visual polish issue, but it reduces readability in a high-visibility part of the board and makes both the header and task cards feel tighter than intended.

## Fix Plan

Increase the spacing around the defect icon in the board header and task cards, preserving the current layout and alignment.

## Acceptance Criteria

- [x] The defect icon no longer appears too close to the adjacent text in the board header.
- [x] The defect icon no longer appears too close to the adjacent text on task cards.
- [x] The spacing is visually consistent with the surrounding elements in both contexts.

## Resolution Notes

Increased spacing after `⚠` from one space to two in both the board header (`view.go`) and the task-card defect marker (`defect_detail.go`). Updated all related tests. All tests pass.
