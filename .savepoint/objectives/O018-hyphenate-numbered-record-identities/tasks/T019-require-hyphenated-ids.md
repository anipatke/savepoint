---
id: T019
title: Require and generate hyphenated IDs
objective: O018
status: planned
owner_validation: {required: false}
planned_by: {role: planner, session: o018-replan-20260923}
---

# T019: Require and generate hyphenated IDs

## Outcome

V2 validation accepts only hyphenated identities, and every code path that
mints an identity emits the hyphenated form.

## User Check

No separate owner interaction is needed. The optional Task Check may be
requested; if skipped, record the owner's explicit waiver in Task evidence.

## Done When

One shared rule replaces the five per-kind regexes; `nextV2CheckID` and the
Issue generator emit `C-###`/`I-###`; the V1→V2 allocator emits `X-###`;
skill, template, and fixture examples are hyphenated with canonical/scaffold
copies byte-identical. Tests cover acceptance, rejection of old and malformed
IDs, and generated output.

## Context Files

`internal/data/release_v2.go`, `internal/data/objective_v2.go`,
`internal/data/task_v2.go`, `internal/data/check_v2.go`,
`internal/data/write.go`, `internal/migrate/plan.go`, `agent-skills/`.

## Guardrails

DATA-01..04, TPL-01..04, TEST-01..04.

## Implementation Plan

1. Add a shared identity pattern/helper in `internal/data`; point each kind's
   validation at it.
2. Update `write.go` Check/Issue allocation (format and prefix parsing) and
   the `migrate/plan.go` allocator.
3. Update skill/template examples and test fixtures; run the skill-copy parity
   check.

## Boundaries

Lands in the same commit as T020; the repository does not load in between.

## Technical Verification

Focused tests while iterating; `make build && make test-fast` before T020.

## Technical Evidence

Pending execution.
