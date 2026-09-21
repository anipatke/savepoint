---
id: T006
title: Preserve existing project records
objective: O013
planned_by: {role: planner, session: goals-terminology-20260921}
status: planned
complexity_tier: medium
complexity_reason: "The public rename must be proven against V2 parsing, membership, router selection, completion, cutover, and migration compatibility without changing the stored contract."
depends_on: []
owner_validation:
    required: true
---

# T006: Preserve existing project records

## Outcome

The compatibility boundary for the Goal rename is explicit and tested: current
V2 projects continue to use their existing `R###` Release records and
`release:` references while the public presentation vocabulary changes.

## User Check

Review the compatibility statement and representative existing-project
fixtures. Confirm that this Objective changes the user-facing name only and
does not silently introduce a `G###`/`goal:` schema migration. If a full
persisted rename is desired, stop and replan this Objective before
implementation.

## Done When

- Existing `R###` Release records and `.savepoint/releases/` paths still load
  through the V2 index.
- Objective `release:` references and router `release:` selections still
  resolve to the same Objectives and Tasks.
- The canonical Release completion and project cutover resolvers remain the
  only gate decisions; no Goal-specific duplicate is introduced.
- Existing migration conversion and historical-completion behavior remains
  covered by named tests or is explicitly recorded as unchanged evidence.
- No parser, writer, or migration path rewrites stored Release data as part of
  the terminology change.

## Context Files

`internal/data/release_v2.go`, `internal/data/release_v2_test.go`, `internal/data/release_gate_v2.go`, `internal/data/release_gate_v2_test.go`, `internal/data/release_cutover.go`, `internal/data/objective_v2.go`, `internal/data/objective_v2_test.go`, `internal/data/project.go`, `internal/data/project_test.go`, `internal/data/router_v2.go`, `internal/data/router_v2_test.go`, `internal/migrate/convert_releases.go`, `internal/migrate/convert_releases_test.go`, `internal/migrate/convert_test.go`.

## Design References

Design sections 2, 3, 4, 5, 6, 10, and 13.

## Guardrails

DATA-02, DATA-03, FS-01, TPL-02, TEST-01, TEST-02, TEST-04, TEST-06,
TEST-08, STYLE-07, and STYLE-10.

## Implementation Plan

1. Confirm the existing R###, `release:`, router, membership, completion, and
   migration contracts before changing presentation code.
2. Add only the missing regression assertions needed to prove old project data
   remains loadable and gate behavior is unchanged; reuse existing tests when
   they already provide exact evidence.
3. Record the compatibility boundary for the UI and documentation Tasks.
4. Run focused data and migration tests, `git diff --check`, `make build`, and
   `make test`.

## Boundaries

No `G###` identity vocabulary, `goal:` frontmatter, goals directory, migration
rewrite, new completion resolver, or changes to the Release/Objective gate
policy.

## Technical Verification

Focused `internal/data` and `internal/migrate` tests covering existing V2
records, Objective membership, router selection, completion/cutover decisions,
and legacy conversion; `git diff --check`; `make build`; and `make test`.

## Technical Evidence

Pending execution: named compatibility cases, tests reused or added, files
read/changed, command results, and limitations.

## Drift Notes

Any requirement to rename persisted identities, paths, frontmatter, or router
keys is a material scope change and routes back to design before implementation.
