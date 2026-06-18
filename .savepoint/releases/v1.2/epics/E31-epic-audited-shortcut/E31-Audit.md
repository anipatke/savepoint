---
type: audit-findings
audited: 2026-06-18
---

# Audit Findings: E31 Mark Epic Audited Shortcut

## Main Findings

E31 satisfies the task acceptance criteria. Pressing `a` while the Epic Detail overlay is open dispatches an epic-status write command, persists `status: audited` through `data.UpdateEpicStatus`, and updates the in-memory `EpicStatus` map from the resulting message so the sidebar glyph can change without a manual refresh. The already-audited path is a no-op with a clear status message.

The write path follows the existing task/defect optimistic-concurrency pattern used in this codebase: the open epic detail file mtime is captured when the overlay opens, `writeEpicStatusCmd` compares it before writing, and conflict/error paths return `errorMsg` without changing model state. Current tests cover the successful keybinding flow, already-audited no-op, conflict messaging, map refresh, and direct file update.

Verification run during audit: `make build && make test` passed.

Unresolved audit item: `.savepoint/Design.md` has minor documentation drift. The board architecture and keybinding sections do not yet mention the new Epic Detail `a` shortcut or immediate epic glyph refresh after marking an epic audited. Proposed replacement blocks are included below for apply/close.

## Code Style Review

- [x] One job per file
- [x] One job per function
- [x] Test branches
- [x] Types document intent
- [x] Build only what is needed
- [x] Handle errors at boundaries
- [x] One source of truth
- [x] Comments explain WHY
- [x] Content in data files
- [x] Small diffs

## Proposed Changes

### Target File
.savepoint/Design.md

### Replace
```md
Epic status glyphs are cached from each epic's `E##-Detail.md` frontmatter and shown in the wide epic sidebar only.
```

### With
```md
Epic status glyphs are cached from each epic's `E##-Detail.md` frontmatter and shown in the wide epic sidebar only; the Epic Detail overlay can persist `status: audited` through an mtime-guarded `a` shortcut and refresh the glyph cache immediately.
```

### Target File
.savepoint/Design.md

### Replace
```md
**Keybindings:** arrow/vim navigation, enter opens focused task detail, space advances, backspace retreats, `p` marks the focused non-done task as router priority, `r`/`R` opens release selection or refreshes where supported, `?` opens help, and `q` quits or closes overlays.
```

### With
```md
**Keybindings:** arrow/vim navigation, enter opens focused task detail, space advances, backspace retreats, `p` marks the focused non-done task as router priority, `a` marks the open Epic Detail epic audited, `r`/`R` opens release selection or refreshes where supported, `?` opens help, and `q` quits or closes overlays.
```
