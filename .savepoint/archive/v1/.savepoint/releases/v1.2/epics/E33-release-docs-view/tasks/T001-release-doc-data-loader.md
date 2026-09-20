---
id: E33-release-docs-view/T001-release-doc-data-loader
status: done
objective: Add a data-layer reader for the known supporting release documents
depends_on: []
complexity_tier: medium
complexity_reason: "Adds a small data API with file IO and missing-file behavior that board code will consume."
---

# T001: Release Doc Data Loader

## Problem

The board needs a bounded, testable way to load supporting documents without
letting UI code know file paths or parse filesystem errors directly.

## Context Files

- `internal/data/discover.go`
- `internal/data/errors.go`
- `internal/data/release_doc.go`
- `internal/data/release_doc_test.go`

## Acceptance Criteria

- [x] `internal/data` exposes a typed release-doc model with stable IDs for PRD
      and Design.
- [x] The loader reads `.savepoint/PRD.md` and `.savepoint/Design.md` from the
      provided `.savepoint` root.
- [x] Missing docs are represented as unavailable document entries instead of a
      fatal load failure.
- [x] Unexpected read errors are returned with path context.
- [x] Tests cover present docs, missing docs, and read-error behavior where
      practical.

## Implementation Plan

- [x] Add `ReleaseDocID` constants for PRD and Design in
      `internal/data/release_doc.go`.
- [x] Add a `ReleaseDoc` struct containing ID, label, relative path, body, and
      availability/error state.
- [x] Implement `LoadReleaseDocs(root string) ([]ReleaseDoc, error)` using the
      `.savepoint` root path passed to board data readers.
- [x] Treat `os.IsNotExist` as a non-fatal unavailable document.
- [x] Add table-driven unit tests in `internal/data/release_doc_test.go`.

## Context Log

- Read: `internal/data/discover.go`, `internal/data/errors.go`,
  `internal/data/defect.go`, `internal/data/task.go`,
  `internal/data/discover_test.go`, `internal/testutil` helpers for conventions.
- Added: `internal/data/release_doc.go` — `ReleaseDocID` constants (`prd`,
  `design`), `ReleaseDoc` model (ID, Label, relative Path, Body, Available), a
  bounded `releaseDocSpecs` list, and `LoadReleaseDocs(root)`. Missing files map
  to `Available: false`; non-`IsNotExist` errors abort with path context via
  `fmt.Errorf("read release doc %s: %w", ...)`.
- Added: `internal/data/release_doc_test.go` — table-driven coverage for both
  present, both missing, mixed presence; label/path assertions; and a portable
  read-error case using a directory in place of `PRD.md` (non-`IsNotExist`).
- Quality gates: `make build && make test` pass; `go test ./internal/data/`
  green.

## Drift Notes

`internal/data/release_doc.go` is a new file but was already listed in the E33
epic Components-and-files map, so no architecture drift beyond the documented
plan. The struct exposes `Available` (no separate `Err` field); unexpected read
errors are surfaced through the function return per the acceptance criteria
rather than stored per-doc, keeping the model minimal.
