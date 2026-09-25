---
id: E01-tui-optimisation/T002-rename-epic-design-files
status: done
objective: "Rename all epic Design.md files to E##-Detail.md convention"
depends_on: []
---

# T002: Rename Epic Design Files

## Acceptance Criteria

- All 8 active epic `Design.md` files renamed to `E##-Detail.md` pattern
- Archived epics under `_archived/` left unchanged
- Root `.savepoint/Design.md` left unchanged
- No content changes in renamed files
- All renamed files remain valid markdown with intact frontmatter

## Implementation Plan

- [x] Rename `.savepoint/releases/v1/epics/E01-go-setup/Design.md` → `E01-Detail.md`
- [x] Rename `.savepoint/releases/v1/epics/E02-data-readers/Design.md` → `E02-Detail.md`
- [x] Rename `.savepoint/releases/v1/epics/E03-board-tui-core/Design.md` → `E03-Detail.md`
- [x] Rename `.savepoint/releases/v1/epics/E04-board-components/Design.md` → `E04-Detail.md`
- [x] Rename `.savepoint/releases/v1/epics/E05-phase-transitions/Design.md` → `E05-Detail.md`
- [x] Rename `.savepoint/releases/v1/epics/E06-atari-noir-layout/Design.md` → `E06-Detail.md`
- [x] Rename `.savepoint/releases/v1.1/epics/E01-tui-optimisation/Design.md` → `E01-Detail.md`
- [x] Rename `.savepoint/releases/v1.1/epics/E02-cross-platform-compatibility/Design.md` → `E02-Detail.md`

## Context Log

Files read:
- None (file renames only)

Estimated input tokens: 300

Notes:
- Content is not modified, only file paths change
- Archived epics explicitly excluded per user request
