---
type: epic-design
status: audited
---

# E31: Audit Register Data Model

## Purpose

Add typed parsing, discovery, and validation for audit prompt, register, finding, and run files so later board and doctor work can rely on structured data.

## What this epic adds

- Audit finding model with stable IDs, lifecycle status, proof, locations, and work-item links.
- Audit run model with prompt version, commit SHA, mode, coverage, and headline counts.
- Repo-wide discovery helpers for `.savepoint/audit/`.
- Validation for malformed statuses, missing proof, duplicate links, and broken work-item references.
- Canonical read helpers that tolerate absent audit-register files.

## Components and files

| Module | Purpose |
|--------|---------|
| `internal/data` | Own audit-register types, parsing, discovery, validation, and errors |
| `.savepoint/audit/` | Source layout consumed by the data model |

## Architectural delta

The data package gains an audit-register domain beside existing task, defect, router, and release-doc models. The model is read-oriented in v1.4 except for reusable canonical helpers needed by future workflows.

## Boundaries

**In scope:**
- Typed parsing and validation
- Absent-register tolerance
- Work-item link checks against discovered tasks and defects
- Unit tests for valid and invalid audit files

**Out of scope:**
- Board rendering
- Doctor report formatting
- Automatic finding matching
- Mutating finding status from the TUI

## Quality gates

- `go test ./internal/data` passes.
- Data tests cover absent, empty, valid, and invalid audit-register states.
- Existing task, defect, router, and release-doc parsing remains unchanged.

## Open decisions

None.
