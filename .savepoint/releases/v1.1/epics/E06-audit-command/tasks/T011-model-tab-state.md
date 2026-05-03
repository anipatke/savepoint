---
id: E06-audit-command/T011-model-tab-state
status: done
objective: Add tab state to board Model for Detail/Audit switching
depends_on: []
---

# T011: Model Tab State

## Acceptance Criteria

- [x] `EpicDetailTab int` added to Model (0=Detail, 1=Audit)
- [x] `EpicAuditContent string` added to Model (cached audit file content)
- [x] Defaults to 0 (Detail tab)
- [x] `make build` passes
- [x] `make test` passes

## Implementation Plan

- [ ] Add to `Model` struct in `internal/board/model.go`:
  ```go
  EpicDetailTab    int    // 0=Detail, 1=Audit
  EpicAuditContent string // cached E##-Audit.md content
  ```
- [ ] Run `make build`
- [ ] Run `make test`