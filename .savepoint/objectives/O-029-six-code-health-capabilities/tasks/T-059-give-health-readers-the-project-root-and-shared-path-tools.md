---
id: T-059
title: Give health readers the project root and shared path tools
objective: O-029
status: planned
depends_on: []
owner_validation: {required: false}
planned_by: {role: planner, session: planning-o029-20261001}
complexity_tier: low
complexity_reason: Small seam change plus pure helpers; no provider formats.
---

# Give health readers the project root and shared path tools

## Outcome

Every report reader can turn the paths a tool writes (absolute, `./`-prefixed, backslashed, or Go import paths) into safe repository-relative paths and keep only a bounded, ordered list of affected items, using one shared set of helpers instead of five copies.

## User Check

None beyond the Full Objective Check.

## Done When

- `ReportInput` carries the project root, and `Collect` passes it to every reader.
- A helper maps a report path to a repository-relative slash path: absolute paths inside the root and relative paths are cleaned; paths outside the root or that fail `validateRelativePath` are rejected so the reader can drop them.
- A helper resolves a Go import path to a repository directory from the nearest `go.mod` at or under the root, within the instance scope (monorepo modules included), and reports an unresolvable import path as unresolved instead of guessing.
- A helper keeps the worst N affected items (at most `MaxEvidence`) in a deterministic order and trims notes to `MaxNoteLen`.
- Tests cover Windows-style and absolute inputs, paths outside the root, a nested module, and an unresolvable import path.

## Context Files

`.savepoint/objectives/O-029-six-code-health-capabilities/Objective.md`; `internal/codehealth/collect.go`; `internal/codehealth/primitives.go`; `internal/codehealth/snapshot.go`; new `internal/codehealth/reader_paths.go`, `internal/codehealth/reader_paths_test.go`.

## Design References

Design sections 1 and 11; O-026 "Shared normalized record and bounds"; O-029 Confirmed Design Decisions.

## Guardrails

FS-05, ARCH-03, ARCH-04, CFG-02, CFG-03, STYLE-07, TEST-01, TEST-02, TEST-04, TEST-08.

## Implementation Plan

1. Add `Root` to `ReportInput` and set it in `collector.measure`.
2. Add the path, Go module, and evidence-bounding helpers in `reader_paths.go`, reusing `validateRelativePath` and `sanitizeLine`.
3. Table tests for each helper, including the Windows and outside-root cases.

## Boundaries

No provider-specific parsing and no reader registration.

## Technical Verification

Focused tests while iterating; `make build && make test-fast` at handoff.

## Technical Evidence

Pending execution.

## Drift Notes

None expected.
