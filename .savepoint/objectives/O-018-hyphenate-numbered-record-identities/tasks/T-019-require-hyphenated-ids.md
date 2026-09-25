---
id: T-019
title: Require and generate hyphenated IDs
objective: O-018
status: done
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: o018-replan-20260923}
check_waiver:
    task: T-019
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-09-23T07:57:57Z"
---

# T-019: Require and generate hyphenated IDs

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

Lands in the same commit as T-020; the repository does not load in between.

## Technical Verification

Focused tests while iterating; `make build && make test-fast` before T-020.

## Technical Evidence

Execution began after confirming Objective O-018 is planned, T-019 has no Task
dependencies, and the active Design is `active`. The starting worktree already
contained user changes in Objective, CLI, migration, and data files; those were
preserved. The owner’s `RELOAD` interruption surfaced malformed router YAML;
the `next_action` indentation was repaired and the route remains at
`state: task`, O-018, T-019. The router and current Task intentionally retain the
pre-T-020 identity form per the same-commit boundary.

Acceptance evidence:

- Shared validation rule: `internal/data/identity_v2.go` defines
  `^[ROTCI]-[0-9]{3,}$` and a kind-aware matcher. Release, Objective, Task,
  Check, Issue, evidence/index, and router validation now use it. Tests cover
  all five valid families, 3+ digits, longer suffixes, wrong families, old
  bare IDs, and malformed IDs.
- Generated identities: `internal/data/write.go` emits `C-###` and `I-###`,
  Check ordering in `internal/data/project.go` reads the hyphenated numeric
  suffix, and `internal/migrate/plan.go` allocates `X-###` for R/O/T/I.
  Generator and allocator tests assert those forms.
- Guidance and fixtures: V2 examples in the canonical Design and Check skills,
  shared Issue-capture reference, and scaffold `AGENTS.md` use hyphens. The
  Design skill, Check skill, and Issue-capture reference each compare
  byte-identically with their V2 scaffold copy (`cmp -s`, exit 0). V2 data,
  board, doctor, init, root integration, and migration golden fixtures were
  updated; migration goldens were regenerated through the documented test
  procedure.
- Converter compatibility: focused testing showed
  `ownerObjectiveID` split `O-001-example` at the new internal hyphen and
  produced `O`. `internal/migrate/convert.go` now extracts the complete
  `O-###` directory identity, with regression coverage in
  `TestOwnerObjectiveID_extractsHyphenatedIDs`.

Verification evidence:

- PASS: `go test ./internal/data`, `go test ./internal/doctor`,
  `go test ./internal/migrate`, `go test ./internal/board/v2`,
  `go test ./internal/init`, and `go test . ./internal/board`.
- PASS: final `make test-full` ran the uncached full Go suite and Linux,
  Darwin, and Windows builds. The full gate initially found stale IDs in root
  resume-matrix and board-dispatch fixtures; those were corrected before the
  passing rerun. The migration package was the slowest portion at 1m54s.
- PASS: `make build`.
- Diagnostic only: `git diff --check` reports CRLF line endings as trailing
  whitespace in already-dirty user files. Those line endings and unrelated
  edits were preserved; this diagnostic is not a configured Task gate.

Extra reads were limited to the V2 Issue/evidence/index/router validators,
Check ordering and writers, migration conversion, targeted data/migration/
board/doctor/init/root tests and fixtures, and V2 scaffold templates/skills.
No Savepoint CLI command was run. No Check record or technical clearance was
written. T-019 remains `in_progress` at `audit`; no Task Check or owner waiver
has been recorded, and only the owner may set this Task to `done`.
