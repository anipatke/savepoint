---
id: E25-durable-data-writes/T004-frontmatter-helper-consolidation
status: planned
objective: Delete the dead defect disk parser and consolidate frontmatter stripping and CRLF handling on internal/data helpers.
depends_on: []
complexity_tier: medium
complexity_reason: Touches three packages to remove duplicates and requires exporting helpers without changing render output.
---

# T004: Frontmatter Helper Consolidation

## Problem

Three remnants of duplicated parsing logic: (1) `ParseDefectFileFromDisk` in `internal/data/parser.go:224-243` is dead code whose second-truncated mtime would spuriously trigger `ErrMtimeConflict` if ever wired to `WriteDefectStatus`; (2) `stripFrontmatter` in `internal/board/epic_panel.go:43-55` re-implements `---` delimiter scanning with different rules than the canonical parser; (3) `hasAcceptanceCriteria` in `internal/doctor/checks.go:276` does its own CRLF replacement that misses bare `\r`. Audit findings M1 and M5.

## Context Files

- `internal/data/parser.go`
- `internal/data/parser_test.go`
- `internal/board/epic_panel.go`
- `internal/board/epic_panel_test.go`
- `internal/doctor/checks.go`
- `internal/doctor/checks_test.go`

## Acceptance Criteria

- [ ] `ParseDefectFileFromDisk` is deleted and the repo-wide search finds no references.
- [ ] `internal/data` exports `NormalizeLineEndings` and a body-extraction helper (reusing `SplitFrontmatterBody`) for non-parsing consumers.
- [ ] `internal/board/epic_panel.go` strips frontmatter via the exported data helper; epic detail overlay rendering output is unchanged for existing test fixtures.
- [ ] `internal/doctor/checks.go` uses the exported normalization helper; a content string using bare `\r` line endings is handled identically to `\r\n`.
- [ ] `go test ./internal/data ./internal/board ./internal/doctor` passes.

## Implementation Plan

- [ ] Delete `ParseDefectFileFromDisk` from `internal/data/parser.go`.
- [ ] Export `NormalizeLineEndings` (keep the unexported name as a delegating alias or update internal call sites).
- [ ] Add a `data` helper that returns the body lines/content with frontmatter removed, tolerating missing frontmatter by returning the input unchanged.
- [ ] Replace `stripFrontmatter` usage in `internal/board/epic_panel.go` with the new helper and update its tests.
- [ ] Replace the local CRLF replace in `internal/doctor/checks.go:276` with the exported helper and extend the doctor test with a bare-`\r` fixture.
- [ ] Run the quality gates listed in the epic.

## Context Log

Pending.
