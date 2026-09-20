---
id: E01-tui-optimisation/T005-update-cross-references
status: done
objective: "Update all cross-references in task files and audit snapshots"
depends_on:
  - E01-tui-optimisation/T002-rename-epic-design-files
  - E01-tui-optimisation/T004-update-instruction-files
---

# T005: Update Cross-References

## Acceptance Criteria

- All task file "Files read" sections reference `E##-Detail.md` instead of `Design.md`
- All audit `proposals.md` files reference renamed epic design files
- All audit `snapshot.md` files reference renamed epic design files
- No stale `Design.md` references remain in any epic-level file paths
- 244 grep matches reduced to 0 for epic-level `Design.md` references

## Implementation Plan

- [x] Find all task files with "Design.md" in content and update epic-level references to `E##-Detail.md`
- [x] Update `.savepoint/audit/E02-data-readers/proposals.md`
- [x] Update `.savepoint/audit/E03-board-tui-core/proposals.md`
- [x] Update `.savepoint/audit/E03-board-tui-core/snapshot.md`
- [x] Update `.savepoint/audit/E04-board-components/proposals.md`
- [x] Update `.savepoint/audit/E04-board-components/snapshot.md`
- [x] Update `.savepoint/audit/E05-phase-transitions/proposals.md`
- [x] Update `.savepoint/audit/E05-phase-transitions/snapshot.md`
- [x] Update `.savepoint/releases/v1/epics/E03-board-tui-core/tasks/T*.md` files
- [x] Update `.savepoint/releases/v1/epics/E04-board-components/tasks/T*.md` files
- [x] Update `.savepoint/releases/v1/epics/E06-atari-noir-layout/tasks/T*.md` files
- [x] Run grep verification: no `epics/.*Design\.md` references remain

## Context Log

Files read:
- All task and audit files listed above

Estimated input tokens: 1200

Notes:
- Only update epic-level Design.md paths (e.g., `E03-board-tui-core/Design.md` → `E03-Detail.md`)
- Root `.savepoint/Design.md` references must NOT be changed
- Audit must-fix applied: stale implementation checklist items were ticked after verifying the completed cross-reference cleanup task scope.
