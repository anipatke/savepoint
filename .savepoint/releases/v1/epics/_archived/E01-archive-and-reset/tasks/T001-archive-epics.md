---
id: E01-archive-and-reset/T001-archive-epics
status: planned
objective: "Move all 11 original v1 epic directories into _archived/"
depends_on: []
---

# T001: Archive Epics

## Acceptance Criteria

- All 11 original epic directories (E01-scaffolding through E11-release-validation) exist under `releases/v1/epics/_archived/`.
- No original epic directories remain under `releases/v1/epics/` except `_archived/`.
- All task files and Design.md files within archived epics are preserved.

## Implementation Plan

- [ ] Create `releases/v1/epics/_archived/` directory.
- [ ] Move each of the 11 original epic directories into `_archived/`.
- [ ] Verify no original epic dirs remain at top level.
