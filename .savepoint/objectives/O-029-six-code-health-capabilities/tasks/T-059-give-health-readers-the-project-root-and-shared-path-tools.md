---
id: T-059
title: Give health readers the project root and shared path tools
objective: O-029
status: in_progress
stage: audit
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

Commands run: `make test-focused TEST='TestRelPath|TestReportInputPath|TestGoModules|TestWorstEvidence|TestCollectHandsReadersTheRoot' PKGS=./internal/codehealth` (pass); `make build && make test-fast` (both pass, no failures).

Per-criterion outcomes:
- `ReportInput` carries the root, and `Collect` passes it: `ReportInput.Root` set in `collector.measure`; `TestCollectHandsReadersTheRoot`.
- Report path helper: `RelPath` and `ReportInput.Path` (also applies instance scope/exclusions). `TestRelPath` covers relative, `./`, backslash, absolute inside, Windows drive (backslash and forward), outside root, sibling-prefix root, traversal, URL, sensitive, empty, and control-character inputs, which are rejected.
- Go import helper: `LoadGoModules` / `GoModules.Resolve` use the longest matching module, nested modules included, scope-limited, and report unresolved import paths as `false`. `TestGoModulesResolve`, `TestGoModulesResolveKeepsToScope`.
- Worst-N helper: `WorstEvidence` caps at `MaxEvidence`, orders worst first then by path/line/note, bounds notes via `sanitizeLine`, drops items that fail evidence validation. `TestWorstEvidence`, `TestWorstEvidenceBoundsNotesAndDropsUnsafeItems`.
- Tests cover Windows-style and absolute inputs, outside-root paths, nested module, unresolvable import: yes, as above.

Files changed: `internal/codehealth/collect.go`; new `reader_paths.go`, `reader_paths_test.go`. Also status/stage edits to this Task and O-029 `Objective.md`.
Files read: only the Context Files, plus `internal/codehealth/runner.go` (extra read: to reuse `sanitizeLine`), `repository.go` (extra read: scope matching via `InputScope.relevant`), `model.go` (extra read: bound constants), `collect_test.go`/`storage_test.go` (extra read: existing test helpers).

Limitations: path mapping is lexical (no symlink resolution); a Windows drive-letter path is matched case-sensitively against the root; module discovery skips hidden, vendor, node_modules, and testdata directories and stops after 20000 directories; `.git`-style sensitive paths are dropped rather than reported.

## Drift Notes

None expected.
