---
id: E06-audit-command/T012-epic-audit-render
status: done
objective: Add RenderEpicAuditTab function to render audit findings + code style checklist
depends_on:
    - E06-audit-command/T011-model-tab-state
---

# T012: Render Epic Audit Tab

## Acceptance Criteria

- [ ] `RenderEpicAuditTab(epicSlug, content string, overlayW, maxHeight, offset int) string` function created
- [ ] Renders "EPIC AUDIT" header
- [ ] Parses E##-Audit.md Markdown content (Main Findings + Code Style Review sections)
- [ ] Renders code style checkboxes with appropriate styling
- [ ] Uses existing scroll helpers (visibleDetailLines)
- [ ] Tests pass

## Implementation Plan

- [x] Add `RenderEpicAuditTab()` function to `internal/board/epic_panel.go`:
  - Similar structure to `RenderEpicDetail()`
  - Parse markdown content: strip frontmatter, extract sections
  - Render header with audit title style
  - Render body: findings content + checklist
  - Reuse `visibleDetailLines`, `WrapText` from detail.go
- [x] Add unit tests for parsing and rendering
- [x] Ensure scroll behavior matches detail overlay

## Fixed

- Section filtering: `epicAuditBody` now only renders "Quality Review" and "Code Style Review" sections, excluding all other sections per acceptance criteria