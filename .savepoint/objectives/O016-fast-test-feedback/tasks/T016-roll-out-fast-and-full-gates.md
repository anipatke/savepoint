---
id: T016
title: Put fast and full gates in place
objective: O016
status: planned
depends_on: [{task: T013, requires: clear}, {task: T014, requires: clear}, {task: T015, requires: clear}]
owner_validation: {required: true}
planned_by: {role: planner, session: o016-design-20260923}
complexity_tier: high
complexity_reason: Build commands, CI, guardrails, and shipped workflow instructions must change as one verification contract.
---

# T016: Put fast and full gates in place

## Outcome

Developers and agents have named focused, fast, and full verification paths, and CI reports their timing while running the complete suite.

## User Check

Review the final gate table and evidence-reuse wording. Confirm an ordinary Task, a migration or platform-sensitive Task, an Objective Check, and a metadata-only correction each select the intended evidence.

## Done When

- Make targets implement T013's exact selection. The full gate runs every host-platform Go test, including migration recovery, interruption, and integration, plus cross-build checks. CI also executes Windows-only tests on Windows.
- Ordinary Task handoff requires build plus fast; migration or platform-sensitive Task handoff requires fresh full; CI and Full Objective/Release Checks require full.
- Focused tests remain an iteration aid. Metadata-only full-result reuse requires a record of the original run and proof that code, tests, fixtures, and gate definitions did not change.
- CI runs full and publishes package/test timing with enough attribution to identify a regression.
- Guardrails, active task/check/design skills, their V2 scaffold copies, AGENTS.md, and Make/CI guidance use the same names and responsibilities; canonical and scaffold skills remain byte-identical.
- Recorded warm-machine measurements assess the 2/15/45-second goals without making elapsed time a flaky pass/fail assertion.
- Selection tests or command transcripts prove both fast omission and full inclusion, including a failing-test path; `make build`, fast, and full pass.

## Context Files

`Makefile`, `.github/workflows/ci.yml`, `AGENTS.md`, `.savepoint/Guardrails.md`, `agent-skills/savepoint-task/SKILL.md`, `agent-skills/savepoint-check/SKILL.md`, `agent-skills/savepoint-design/SKILL.md`, `templates/project-v2/agent-skills/savepoint-task/SKILL.md`, `templates/project-v2/agent-skills/savepoint-check/SKILL.md`, `templates/project-v2/agent-skills/savepoint-design/SKILL.md`, `.savepoint/objectives/O016-fast-test-feedback/tasks/T013-measure-and-design-test-gates.md`.

## Design References

Design sections 1, 12, and 13; O016 Confirmed Verification Policy.

## Guardrails

TPL-01..04, TEST-01..04, TEST-07..09, CFG-01..02, ARCH-01, STYLE-07.

## Implementation Plan

1. Implement named fast and full targets from T013's selection, retaining `make test` as the full suite unless the confirmed contract explicitly revises it.
2. Add selection verification and CI timing output; ensure CI runs full and a Windows test job.
3. Reconcile TEST-08 and all active/scaffolded workflow guidance in the same integrated change.
4. Document evidence freshness and metadata-only reuse with explicit input comparison.
5. Exercise ordinary, sensitive, CI, Check, and reuse scenarios; record timings and run the revised complete gate.

## Boundaries

No test removal, test-cache substitution for changed inputs, hidden CI exclusion, or automatic timeout failure from timing goals.

## Technical Verification

Gate selection and failure-path evidence, active/scaffold parity, full CI-equivalent command including cross-build checks and Windows test evidence, warm timing samples, `make build`.

## Technical Evidence

Pending execution: gate selection inventory, named results, timing reports, parity comparison, and limitations.
