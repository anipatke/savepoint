---
type: epic-design
status: planned
---

# E30: Epic Status Self-Heal

## Purpose

Stop epics from rendering with a blank status glyph when an agent writes a
non-canonical `status:` into an `E##-Detail.md`. Epic status is currently loaded
as a raw free string with no normalization, so any value outside the four the
glyph map knows about disappears from the board. This epic gives epic status the
same load-time self-heal + doctor-surfacing contract that task and defect
lifecycles already have.

## What this epic adds

- A canonical epic-status vocabulary defined in one place in `internal/data`.
- Load-time normalization so an unknown or alias status heals to a safe canonical
  value and always renders a glyph.
- A visible fallback glyph so a truly unmapped status is never invisible.
- A `savepoint doctor` diagnostic + repair hint that reports the heal so the
  underlying file gets corrected.

## Components and files

| Module | Purpose |
|--------|---------|
| `internal/data/lifecycle.go` | Canonical epic statuses, `NormalizeEpicStatusForLoad`, `ResolveEpicStatusAlias`, `DiagnoseEpicStatus` |
| `internal/data/lifecycle_test.go` | Normalization + diagnostic unit tests |
| `internal/board/board.go` | Apply normalization when the `epicStatuses` map is built at load |
| `internal/board/status.go` | Visible fallback glyph for an unmapped status |
| `internal/board/status_test.go` | Glyph fallback test |
| `internal/doctor/checks.go` | Surface `DiagnoseEpicStatus` from `checkEpicDetail` |
| `internal/doctor/repairs.go` | Repair guidance string for non-canonical epic status |

## Architectural delta

Before: epic status is read verbatim (`internal/board/board.go:134`) and mapped
to a glyph by `statusGlyph` (`internal/board/status.go:9`), whose `default`
branch returns a blank space. Any status not in
`{planned, in_progress, done, audited}` — including agent typos and stray
router-state leaks such as `epic-design` — renders invisibly.

After: epic status has a single canonical vocabulary owned by `internal/data`,
is normalized on load exactly like task `stage` and defect `status`, surfaces the
heal through `savepoint doctor`, and never renders blank.

## Boundaries

**In scope:**
- Canonical epic-status constants, load normalization, glyph fallback, doctor
  diagnostic + repair hint, and tests.

**Out of scope:**
- Any TUI control that *changes* epic status (the audited shortcut is E31).
- Router state vocabulary (`epic-design`, `audit-pending`, …) — those remain
  router states; this epic only treats them as non-canonical *epic statuses*.
- Rewriting existing epic files on disk; healing is at load + doctor-reported.

## Quality gates

- `NormalizeEpicStatusForLoad` maps every non-canonical value to a canonical one
  and is the single source of the epic-status vocabulary.
- The board never renders a blank epic glyph for any input.
- `savepoint doctor` reports a non-canonical epic status with an actionable hint.
- `make build && make test` passes.

## Open decisions

None. Canonical set is `{planned, in_progress, done, audited}`; unknown heals to
`planned` (lowest lifecycle), mirroring defect unknown → `open`.
