---
id: E13-audit-remediation/T003-centralize-duplication
status: done
objective: Centralize \r\n normalization, extract SplitFrontmatterBody, consolidate shared board utilities
depends_on:
  - E13-audit-remediation/T001-safe-cleanup
---

# T003: Centralize Duplicated Logic — Normalization, Frontmatter, Utilities

## Context Files

- `internal/data/parser.go` — `extractFrontmatter()` and `\r\n` normalization
- `internal/data/write.go` — duplicated `delimLen + len(raw) + delimLen` body offset calculation (2 places)
- `internal/board/epic_panel.go` — duplicated frontmatter stripping in `epicDetailBody` and `epicAuditBody`; `epicIndex()` function
- `internal/board/card.go` — `shortID()` used by 6+ files
- `internal/board/view.go` — `shortRouterID()` near-duplicate of `shortID()`
- `internal/board/detail.go` — `WrapText()` and `SplitLongWord()` general-purpose text utilities
- `internal/board/release.go` — `releaseIndex()` identical to `epicIndex()`

## Acceptance Criteria

- [ ] `normalizeLineEndings(s string) string` added to `internal/data/parser.go`; all `\r\n` replacement call sites updated to use it
- [ ] `SplitFrontmatterBody(content string) (yaml string, body string, err error)` added to `internal/data/write.go`; `updateFrontmatterField()` and `WriteTaskStatus()` refactored to use it
- [ ] `stripFrontmatter(content string) string` helper added in `internal/board/epic_panel.go`; `epicDetailBody` and `epicAuditBody` use it
- [ ] `shortID()` and `shortRouterID()` consolidated into single `shortID()` in `internal/board/card.go`; all call sites updated
- [ ] `epicIndex()` and `releaseIndex()` consolidated into `sliceIndex(items []string, target string) int` in `internal/board/util.go`
- [ ] `WrapText()` and `SplitLongWord()` moved to `internal/board/util.go`
- [ ] `truncate()` moved to `internal/board/util.go`
- [ ] `go test ./...` passes with no regressions

## Implementation Plan

- [x] Add `normalizeLineEndings()` in `data/parser.go`; update call sites in `parser.go`, `write.go` (3 places), and `epic_panel.go`
- [x] Add `SplitFrontmatterBody()` in `data/write.go`; refactor `updateFrontmatterField()` and `WriteTaskStatus()` to use it; eliminate duplicated body offset calculation
- [x] Add `stripFrontmatter()` in `board/epic_panel.go`; refactor `epicDetailBody` and `epicAuditBody`
- [x] Create `internal/board/util.go` with `sliceIndex`, `WrapText`, `SplitLongWord`, `truncate` moved from their current locations
- [x] Consolidate `shortID`/`shortRouterID` into single `shortID` in `card.go`; update `view.go` callers
- [x] Replace `epicIndex`/`releaseIndex` with `sliceIndex` from `util.go`; update `epic_panel.go` and `release.go`
- [x] Run `make build && make test`

## Context Log

- Files modified:
  - `internal/data/parser.go` — added `normalizeLineEndings()`, updated 3 call sites
  - `internal/data/write.go` — added `SplitFrontmatterBody()`, refactored `updateFrontmatterField()` and `WriteTaskStatus()`, updated 3 call sites to use `normalizeLineEndings()`
  - `internal/board/epic_panel.go` — added `stripFrontmatter()`, refactored `epicDetailBody` and `epicAuditBody`, replaced `epicIndex()` with `sliceIndex()`
  - `internal/board/card.go` — removed `truncate()` (moved to util.go)
  - `internal/board/detail.go` — removed `WrapText()` and `SplitLongWord()` (moved to util.go)
  - `internal/board/view.go` — replaced `shortRouterID()` with `shortID()`, removed `shortRouterID()`
  - `internal/board/release.go` — replaced `releaseIndex()` with `sliceIndex()`
  - `internal/board/update.go` — replaced `epicIndex()`/`releaseIndex()` with `sliceIndex()`
  - `internal/board/model.go` — replaced `epicIndex()` with `sliceIndex()`
- Files created: `internal/board/util.go`
- Test files updated: `epic_panel_test.go`, `release_test.go`
- Quality gates: `make build` ✓, `make test` ✓ (all packages pass)

## Drift Notes

No drift — `internal/board/util.go` is a new file within the existing `internal/board/` module, consistent with the Codebase Map. No architecture changes from Design.md.