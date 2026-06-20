---
id: E33-release-docs-view/T004-release-docs-verification
status: done
objective: Add regression coverage and quality-gate verification for Release Docs behavior
depends_on:
    - E33-release-docs-view/T003-release-docs-renderer
complexity_tier: low
complexity_reason: "Focuses on tests and documented verification after the feature modules are in place."
---

# T004: Release Docs Verification

## Problem

The Release Docs feature touches board navigation and rendering, so it needs
focused regression tests that prove existing Epic detail behavior remains intact
while the new read-only view works at narrow and normal widths.

## Context Files

- `internal/board/epic_panel_test.go`
- `internal/board/update_test.go`
- `internal/board/io_test.go`
- `internal/data/release_doc_test.go`
- `Makefile`

## Acceptance Criteria

- [x] Tests prove the Epic detail overlay defaults to existing Epic content.
- [x] Tests prove Release Docs can render PRD and Design labels and selected
      body content.
- [x] Tests prove missing documents render without panics or fatal overlay
      errors.
- [x] Tests prove long document lines do not exceed the rendered content width.
- [x] Tests prove document selection and scroll state behave predictably.
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Add or extend board render tests for default Epic detail behavior.
- [x] Add Release Docs render tests at normal and narrow widths.
- [x] Add update tests for subview switching, document switching, and scroll
      behavior.
- [x] Add IO command/message tests for successful and missing release docs.
- [x] Run `make build && make test` and record the result in this task's
      Context Log when implemented.

## Context Log

- Read: `internal/board/epic_panel.go`, `update.go`, `io.go`,
  `internal/data/release_doc.go`, and the existing test files
  (`epic_panel_test.go`, `update_test.go`, `io_test.go`,
  `release_doc_test.go`, `util.go`) to map current coverage against the ACs.
- Found T003 already landed the bulk of the render/state/loader coverage:
  selector labels, selected/missing/empty/no-docs bodies, width-bound body
  wrapping at width 48, frontmatter stripping (`epic_panel_test.go`); tab
  switching, `[`/`]` doc switching, per-doc scroll, index clamping, and overlay
  state reset (`update_test.go`); and loader success/missing/read-error
  (`internal/data/release_doc_test.go`). These satisfy the
  selection/scroll/missing-doc/normal-width ACs and were verified, not
  duplicated.
- Closed the remaining gaps:
  - `internal/board/io_test.go`: added `loadReleaseDocsCmd` command tests —
    success returns `releaseDocsMsg` with both bodies, a missing doc yields an
    unavailable entry (not fatal), and a read error (directory in place of
    `PRD.md`) surfaces as `errorMsg` with path context. Added a package-local
    `docByID` helper.
  - `internal/board/epic_panel_test.go`: added
    `TestView_epicDetailDefaultsToEpicContent` (overlay defaults to Epic
    content even with docs loaded), `TestRenderEpicReleaseDocs_narrowWidth`
    (header/selector/body render at overlay width 36), and
    `TestRenderReleaseDocBody_wrapsWithinNarrowWidth` (long unbreakable token
    hard-cut so no line exceeds width 16).
- Quality gates: `make build && make test` pass; the six new tests pass under
  `-run`; `go vet ./internal/board/ ./internal/data/` clean.
