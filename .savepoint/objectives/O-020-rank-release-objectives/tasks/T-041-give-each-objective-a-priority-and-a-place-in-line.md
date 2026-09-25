---
id: T-041
title: Give each Objective a priority and a place in line
objective: O-020
status: planned
complexity_tier: medium
complexity_reason: Adds two optional validated Objective fields, one ordered projection, one duplicate-rank index fact with a doctor warning, and a multi-record order writer.
depends_on: []
owner_validation: {required: true}
planned_by: {role: planner, session: o020-rank-objectives-20260925}
---

# Give each Objective a priority and a place in line

## Outcome

`internal/data` reads an optional `priority` and `rank` on every Objective,
returns a Goal's Objectives in one canonical display order, warns through
doctor when two Objectives share a rank, and can safely rewrite a group's
order. No view parses these fields itself.

## User Check

In a scratch project, add `priority: high` and `rank: 2` to one Objective and
`priority: high` and `rank: 1` to another, and leave a third without either.
Confirm `savepoint doctor` reports nothing new. Give the two High Objectives
the same rank and confirm doctor names both and the shared rank as a warning.
Set `priority: urgent` and confirm the project refuses to load with a message
naming the file and the allowed values.

## Done When

- `ObjectiveV2` gains `Priority` (typed: critical, high, medium, low; missing
  decodes as medium) and `Rank` (positive integer; zero value means unranked).
  An unknown priority or a non-positive or non-integer rank is a strict decode
  error naming the path, Objective, field, and allowed values.
- One exported function returns a Goal's Objective IDs in display order:
  priority group (Critical, High, Medium, Low), then rank ascending with
  unranked after ranked, then ID. Legacy records with neither field keep ID
  order inside Medium.
- The V2 index records duplicate ranks per Goal and priority group as a
  non-fatal fact. Doctor reports each as a warning naming the Goal, priority,
  rank, and Objectives. It never fails the load.
- One exported writer takes a Goal, a priority, and an ordered list of
  Objective IDs and sets `priority` and `rank` 1..n on each, writing only the
  records whose values change. It checks that every record in the set is fresh
  before writing any; a stale record returns a conflict error and writes
  nothing. Each file is replaced atomically and all other frontmatter and body
  bytes are preserved.
- The Objective Artifact Template in `agent-skills/savepoint-design/SKILL.md`
  and its scaffold copy show the two optional fields, byte-identical (TPL-01).
- Tests cover decode defaults and each invalid value, all four groups,
  ranked-before-unranked, ties by ID, duplicate-rank reporting, a no-op write,
  a renumber that changes only some files, a stale record refusing the whole
  set, and preservation of unrelated content. A failure injected after the
  first file of a multi-file write leaves a loadable project whose order is
  deterministic and whose duplicate is reported.

## Context Files

`internal/data/objective_v2.go`, `internal/data/objective_v2_test.go`,
`internal/data/project.go`, `internal/data/project_test.go`,
`internal/data/write.go`, `internal/data/write_test.go`,
`internal/doctor/v2_runtime.go`, `internal/doctor/v2_runtime_test.go`,
`agent-skills/savepoint-design/SKILL.md`,
`templates/project-v2/agent-skills/savepoint-design/SKILL.md`,
`.savepoint/objectives/O-020-rank-release-objectives/Objective.md`.

## Design References

Design sections 1 (V2 evidence and identity boundary), 9 (concurrency), and
11 (failure modes). O-020 Confirmed Design Decisions.

## Guardrails

FS-01, FS-04, DATA-01, DATA-02, DATA-03, DATA-04, TPL-01, ARCH-04,
TEST-01..04, TEST-08, STYLE-07.

## Implementation Plan

1. Add a typed priority with its four values, a parse helper, and the
   Critical-to-Low sort weight. Decode `priority` and `rank` in
   `DecodeObjectiveV2` with named errors.
2. Add the ordered-IDs function next to `ReleaseObjectives` in the index, and
   the duplicate-rank fact built during index load.
3. Report the fact as a doctor warning in `v2_runtime.go`.
4. Add the order writer on top of `writeV2Record`: collect patches, run the
   freshness check on every record, then write each changed record. If
   `write_test.go` does not exist, create it.
5. Update the Objective template in both skill copies.
6. Write the tests listed above.

## Boundaries

No board changes, key handling, or rendering. No change to status,
dependencies, gates, `ResolveNext`, or Goal completion. No Goal-owned member
list.

## Technical Verification

Focused tests while iterating; `make build && make test-fast` at handoff.

## Technical Evidence

Pending execution.

## Drift Notes

None expected. Design section 1 gains one sentence on priority and rank at
Objective reconciliation.
