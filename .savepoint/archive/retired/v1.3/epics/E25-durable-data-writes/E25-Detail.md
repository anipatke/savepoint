---
type: epic-design
status: planned
---

# E25: Durable Data Writes

## Purpose

Make every write to `.savepoint/` files crash-safe and consolidate file parsing/writing on one source of truth, closing audit findings H1, M1, M5, and L5 from `project-audit/audit_report_fable_5.md`.

## What this epic adds

- Atomic temp-file + rename writes for all task, defect, router, and proposal updates, so a crash or full disk can never truncate a `.savepoint/` file.
- One shared atomic write helper used by both `internal/data` and `internal/init` instead of two write strategies.
- A router state write that cannot corrupt files containing non-ASCII text.
- Removal of the dead `ParseDefectFileFromDisk` and its latent mtime-truncation bug.
- A single canonical frontmatter-stripping and line-ending-normalization path shared by board rendering and doctor checks.

## Components and files

| Module | Purpose |
|--------|---------|
| `internal/data/atomic.go` (new) | Own the shared atomic write helper |
| `internal/data/write.go` | Route all five write paths through atomic writes; fix router block search |
| `internal/data/parser.go` | Remove dead defect parser; export canonical frontmatter/normalization helpers |
| `internal/init/write.go` | Delegate to the shared atomic helper instead of owning a copy |
| `internal/board/epic_panel.go` | Replace ad-hoc frontmatter stripping with the canonical helper |
| `internal/doctor/checks.go` | Replace local CRLF handling with the canonical helper |

## Architectural delta

Write durability moves from an `init`-only concern to a `data`-package guarantee. `internal/data` becomes the single owner of frontmatter parsing, line-ending normalization, and file replacement; `init`, `board`, and `doctor` consume those helpers instead of carrying private copies.

## Boundaries

**In scope:**
- Atomic write helper relocation and adoption at every `internal/data` write site
- `WriteRouterState` search hardening
- Dead code removal in `internal/data/parser.go`
- Frontmatter/CRLF helper consolidation in board and doctor

**Out of scope:**
- File watcher behavior (E26)
- Any change to lifecycle rules, mtime conflict semantics, or file formats
- fsync policy changes beyond what the existing helper already does

## Quality gates

- `go test ./internal/data ./internal/init ./internal/board ./internal/doctor` passes.
- `make build && make test` passes.
- No remaining direct `os.WriteFile` calls in `internal/data/write.go`.

## Open decisions

None.
