---
id: E06-audit-command/T013-handle-tab-keys
status: done
objective: Handle 1/2 keys to switch tabs in Epic Detail overlay
depends_on:
    - E06-audit-command/T011-model-tab-state
    - E06-audit-command/T012-epic-audit-render
---

# T013: Handle Tab Keys

## Acceptance Criteria

- [x] Pressing `1` switches to Detail tab (EpicDetailTab = 0)
- [x] Pressing `2` switches to Audit tab (EpicDetailTab = 1)
- [x] When switching to Audit tab, loads E##-Audit.md content
- [x] Audit content cached in EpicAuditContent
- [x] Scrolling works on both tabs
- [x] Tests pass

## Implementation Plan

- [ ] In `internal/board/update.go`:
  - When `m.Overlay == OverlayEpicDetail`:
    - If key is `1` → set `m.EpicDetailTab = 0`
    - If key is `2` →:
      - set `m.EpicDetailTab = 1`
      - if `m.EpicAuditContent == ""`:
        - derive audit file path
        - read file with `os.ReadFile`
        - if error → set `m.EpicAuditContent = "(no audit available)"`
        - else → cache content
- [ ] Add test for tab switching
- [ ] Add test for audit file loading on first access