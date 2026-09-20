---
type: epic-design
status: audited
---

# E31: Mark Epic Audited Shortcut

## Purpose

Give the user a way to mark an epic `audited` from the board. The audit workflow
ends with "mark the epic audited and advance the router," and
`data.UpdateEpicStatus` already exists to write the field — but nothing in the
TUI ever calls it, so the only way to set an epic audited is to hand-edit
frontmatter. This epic wires a keyboard shortcut in the Epic view to that
existing write helper.

## What this epic adds

- A keyboard shortcut in the epic-detail overlay that sets the focused epic's
  status to `audited`.
- A board write command that persists the status via `data.UpdateEpicStatus` and
  refreshes the in-memory `EpicStatus` map so the glyph updates immediately.
- Help/footer and audit-tab hint text describing the shortcut.

## Components and files

| Module | Purpose |
|--------|---------|
| `internal/board/update.go` | Handle the audited key in `handleEpicDetailOverlay` |
| `internal/board/io.go` | `writeEpicStatusCmd` wrapping `data.UpdateEpicStatus` |
| `internal/board/model.go` | Refresh `EpicStatus` after a successful write |
| `internal/board/help.go` | Help row for the shortcut |
| `internal/board/epic_panel.go` | Audit-tab hint text |
| `internal/data/write.go` | Existing `UpdateEpicStatus` (reused, not changed) |
| `internal/board/update_test.go` / `io_test.go` | Behaviour + write tests |

## Architectural delta

Before: `handleEpicDetailOverlay` (`internal/board/update.go:376`) handles only
tab switching, scrolling, and close. `data.UpdateEpicStatus`
(`internal/board/../data/write.go:31`) is dead code from the TUI's perspective.
After: pressing the shortcut in the Epic view writes `status: audited` through
the same frontmatter-write path the rest of the board uses, with an mtime-guarded
command and an immediate glyph refresh — consistent with how task advancement
writes status.

## Boundaries

**In scope:**
- One epic-detail keybinding, its write command, the `EpicStatus` refresh, hint
  text, and tests.

**Out of scope:**
- Defining/normalizing the epic-status vocabulary (that is E30; this epic writes
  the canonical literal `audited`).
- Advancing the router state or writing the `E##-Audit.md` file (the audit skill
  owns that).
- A general epic status-cycling control; this epic only sets `audited`.

## Quality gates

- The shortcut persists `status: audited` to the correct `E##-Detail.md` and the
  glyph updates without a manual refresh.
- Writes are mtime-guarded and report a clear status message on conflict/error.
- `make build && make test` passes.

## Open decisions

- Whether to gate the shortcut on a prior state (e.g. only from `done`) or allow
  it from any state with a status message. Default: allow from any non-`audited`
  state, no-op + message if already `audited`. Revisit in task breakdown.
