---
id: T-026
title: Reconcile guidance with the smaller code base
objective: O-021
planned_by: {role: planner, session: o021-design-20260923}
status: planned
complexity_tier: low
complexity_reason: "Documentation-only reconciliation plus a line-count report; no production behavior changes."
depends_on: [{task: T-025, requires: clear}]
owner_validation:
    required: false
    accepted_check: ""
---

# T-026: Reconcile guidance with the smaller code base

## Outcome

AGENTS.md, `Design.md`, and the V2 scaffold guidance describe the code base
as it now is: no V1 board, no V1 templates or skills, a `data`-owned schema
check, and a clean-git migration. The Objective records the measured size
reduction.

## User Check

Read the AGENTS.md Codebase Map: every row is one or two plain sentences,
and no row describes a deleted package, file, or behavior.

## Done When

- AGENTS.md's Codebase Map lists only live modules, each in one or two
  sentences. Rows for deleted V1 board behavior, `release_cutover.go`,
  recoverable apply, and the V1 template tree are gone or corrected. The
  closing "V1 readers ... compatibility code" paragraph states what
  remains: the V1 readers `internal/migrate` uses.
- `Design.md` no longer claims recoverable apply, interruption recovery,
  or the migration journal (command table, `.savepoint/` tree, test
  table, and header notes), and states the clean-git rule and the `data`
  schema check.
- `templates/project-v2/AGENTS.md` and any V2 skill text that mentions
  recovery or pending migrations are aligned; template freshness tests
  pass.
- The Objective's Technical Evidence (or this Task's) records production
  and test line counts per package against the baseline in O-021's Why,
  and a final `deadcode` report.
- `git diff --check`, `make build && make test-fast` pass.

## Context Files

`AGENTS.md`, `.savepoint/Design.md`, `templates/project-v2/AGENTS.md`,
`README.md`, `internal/init/template_freshness_test.go`,
`internal/init/agent_skills_test.go`.

## Design References

Design sections 1, 2, 10, and 11.

## Guardrails

ARCH-04, TPL-01, TPL-02, TEST-06, TEST-08.

## Implementation Plan

1. Search guidance for references to deleted code and behavior.
2. Rewrite the Codebase Map rows concisely and correct Design sections.
3. Align V2 scaffold guidance and run template tests.
4. Measure and record line counts and `deadcode` output.

## Boundaries

No production code changes. No edits to archived, Check, or closed Issue
records. O-014's router guidance changes stay with O-014.

## Technical Verification

`make build && make test-fast`; the Full Objective Check then runs a fresh
`make test-full`.

## Technical Evidence

Pending execution.

## Drift Notes

None expected.
