---
id: E33-release-docs-view/T004-release-docs-verification
status: planned
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

- [ ] Tests prove the Epic detail overlay defaults to existing Epic content.
- [ ] Tests prove Release Docs can render PRD and Design labels and selected
      body content.
- [ ] Tests prove missing documents render without panics or fatal overlay
      errors.
- [ ] Tests prove long document lines do not exceed the rendered content width.
- [ ] Tests prove document selection and scroll state behave predictably.
- [ ] `make build && make test` passes.

## Implementation Plan

- [ ] Add or extend board render tests for default Epic detail behavior.
- [ ] Add Release Docs render tests at normal and narrow widths.
- [ ] Add update tests for subview switching, document switching, and scroll
      behavior.
- [ ] Add IO command/message tests for successful and missing release docs.
- [ ] Run `make build && make test` and record the result in this task's
      Context Log when implemented.

## Context Log

Pending.
