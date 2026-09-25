---
type: epic-design
status: planned
---

# E26: Board Watch Resilience

## Purpose

Keep the TUI board's live file watching alive through transient failures, closing audit findings H2 and L6 from `project-audit/audit_report_fable_5.md`.

## What this epic adds

- File watching that survives a failed reload: a malformed file no longer permanently stops live updates.
- Board reloads that skip unparseable files with a visible status note instead of aborting the entire reload.
- Watcher errors surfaced through the existing debug log instead of being silently discarded.

## Components and files

| Module | Purpose |
|--------|---------|
| `internal/board/watch.go` | Distinguish reload failure from generic errors; log watcher errors |
| `internal/board/update.go` | Re-arm the watch command on the reload-failure path |
| `internal/board/board.go` | Skip-and-report malformed files during `loadBoardData` |

## Architectural delta

The watch loop becomes a guaranteed cycle: every message produced by a watch-triggered command path must end by re-arming `watchFiles`. Reload changes from all-or-nothing to best-effort with explicit problem reporting, matching how doctor already treats malformed files as diagnosable problems rather than fatal errors.

## Boundaries

**In scope:**
- Reload failure message type and watcher re-arm
- Per-file parse-error tolerance in board data loading with status reporting
- Watcher error channel logging

**Out of scope:**
- Write durability (E25)
- Incremental/partial reload optimization — full reload on change stays
- Watcher debounce timing and directory-watch strategy
- Doctor behavior

## Quality gates

- `go test ./internal/board` passes.
- A regression test proves a failed reload still re-arms file watching.
- A regression test proves one malformed task file does not hide the remaining tasks.
- `make build && make test` passes.

## Open decisions

None.
