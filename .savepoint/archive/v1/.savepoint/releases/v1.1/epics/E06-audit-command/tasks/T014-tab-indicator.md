---
id: E06-audit-command/T014-tab-indicator
status: done
objective: Update Epic Detail overlay header to show tab indicator and switch between tabs in view.go
depends_on:
    - E06-audit-command/T011-model-tab-state
    - E06-audit-command/T012-epic-audit-render
    - E06-audit-command/T013-handle-tab-keys
---

# T014: Tab Indicator

## Acceptance Criteria

- [ ] Tab indicator shown in overlay header: "EPIC DETAIL [1] │ AUDIT [2]"
- [ ] Active tab highlighted/bold
- [ ] Switches between RenderEpicDetail and RenderEpicAuditTab based on EpicDetailTab state
- [ ] Tests pass and verify indicator rendering

## Implementation Plan

- [ ] In `internal/board/view.go`:
  - When `m.Overlay == OverlayEpicDetail`:
    - Check `m.EpicDetailTab` state
    - If 0 → call existing `RenderEpicDetail()`
    - If 1 → call `RenderEpicAuditTab()`
  - Pass `EpicDetailTab` to renderer for header rendering
- [ ] Update `RenderEpicDetail()` in `epic_panel.go`:
  - Accept `tab int` parameter
  - Render tab indicator line before content
  - Highlight active tab number (bold or different color)
- [ ] Update `RenderEpicAuditTab()` similarly
- [ ] Add test for tab indicator rendering with both tab states
