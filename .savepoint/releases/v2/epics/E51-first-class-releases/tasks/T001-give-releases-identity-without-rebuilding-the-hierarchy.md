---
id: E51-first-class-releases/T001-give-releases-identity-without-rebuilding-the-hierarchy
status: planned
objective: Load optional Release records with stable identity and derive their Objective membership from one authoritative reference.
depends_on: []
complexity_tier: high
complexity_reason: Adds an indexed record family and cross-record ownership validation while preserving Objective and Task authority.
---

# T001: Give Releases identity without rebuilding the hierarchy

## Problem

V2 currently stores Release as optional text on Objectives. That can filter a list, but it cannot preserve identity, authored release intent, or an accountable destination for migration. The project index needs a real Release record while keeping Objective → Task as the work hierarchy and keeping projects with no Releases valid.

## Context Files

- `.savepoint/releases/v2/epics/E51-first-class-releases/E51-Detail.md`
- `.savepoint/releases/v2/v2-Design.md`
- `internal/data/objective_v2.go`
- `internal/data/objective_v2_test.go`
- `internal/data/project.go`
- `internal/data/project_test.go`
- `internal/data/migration_history_test.go`
- `internal/data/release_v2.go`
- `internal/data/release_v2_test.go`

## Acceptance Criteria

- [ ] V2 decodes a Release at `.savepoint/releases/R###-slug/Release.md` with explicit `id`, `title`, `status`, and the required Outcome, Why, Success Conditions, and Boundaries body sections.
- [ ] Release IDs use global `R###` identity, remain stable when title/path changes, and are never inferred or reassigned from a path.
- [ ] A project with no `releases/` directory or no Release records loads exactly as it did before E51.
- [ ] An Objective may omit `release`; when present, `release` is an `R###` reference to an indexed Release rather than packaging text.
- [ ] `ReleaseObjectives` and Objective-to-Release links are derived from Objective records; Release records contain no manually maintained membership list.
- [ ] Tasks continue to belong only to Objectives and acquire Release context only through the indexed Objective relationship.
- [ ] Duplicate Release identity, malformed lifecycle, unsafe paths, filename/identity mismatch, and dangling Objective references fail closed with the source path and ID named.
- [ ] Unassigned Objectives remain valid and independently indexed beside Objectives assigned to two or more Releases.
- [ ] Existing Task, Objective, Check, and Issue indexing behavior is unchanged for projects without Releases.

## Implementation Plan

- [ ] Add the typed Release model, strict decoder, body-section validation, lifecycle constants, and path/identity checks.
- [ ] Extend V2 discovery to load the optional Release record family through the existing confined-path boundary.
- [ ] Change Objective `release` from free packaging text to an optional typed Release reference.
- [ ] Add Release maps and derived membership links to `V2Index` without storing a second membership list.
- [ ] Include Release IDs in allocation/reservation checks so archived or migrated identities are not reused.
- [ ] Cover zero, one, and multiple Releases; unassigned Objectives; title/path changes; duplicates; dangling references; and confinement failures.
- [ ] Prove legacy V2 fixtures without Release records still load byte-for-byte equivalently at the consumer boundary.

## Context Log

Pending.
