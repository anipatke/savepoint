---
id: T-114
title: Drop migration history and use one term per role
objective: O-039
status: planned
depends_on: [{task: T-113, requires: clear}]
complexity_tier: low
complexity_reason: Wording edits in the same parity-checked assets, with content tests.
owner_validation: {required: false}
planned_by: {role: planner, session: plan-o039-20261006}
---

# Drop migration history and use one term per role

## Outcome

Active skills, references, and the managed block no longer narrate V1/V2 history outside the managed block's Required Goal Context. Prose says "owner" for the human decision-maker and "Goal" for the record, and `release:` appears only as the field name.

## User Check

Search the skills and references for "V1", "resurrect", "user may", and "Release" in prose, and find no history narrative or mixed terms.

## Done When

- Each removed history passage is listed in Technical Evidence and is either covered by Required Goal Context or no longer needed by an agent.
- "user" never names the decision-maker where "owner" is meant, and "Release" never names a Goal in prose. Field names and R-### identities are unchanged.
- Live and scaffold copies are byte-identical (TPL-01). Content tests are updated, with one assertion against the old wording returning.
- O-034's scenario walkthroughs are re-walked, with the result recorded per scenario.
- Token weight before and after is recorded.
- `make build && make test-fast` passes.

## Context Files

`templates/project-v2/AGENTS.md`; `AGENTS.md`; the four workflow skills and three references under `agent-skills/` and their copies under `templates/project-v2/agent-skills/`; `internal/init/agent_skills_test.go`; `internal/init/skill_validation_test.go`; `internal/init/template_freshness_test.go`; `.savepoint/objectives/O-034-skills-optimisation/tasks/T-099-review-the-skills-against-what-the-code-health-goal-taught-us.md`.

## Design References

O-039 Confirmed Design; AGENTS.md Terminology.

## Guardrails

TPL-01, TPL-02, TPL-04, FS-02, POL-01, POL-02, TEST-01, TEST-02, TEST-03, TEST-08, STYLE-07, STYLE-10.

## Implementation Plan

1. Record the before token weight.
2. Remove history passages and settle terms file by file, mirroring each change to the scaffold copy.
3. Update content tests.
4. Re-walk the scenarios, record the after weight, and run the gate.

## Boundaries

No rule changes, no field or identity renames, and no Go runtime changes.

## Technical Verification

Focused tests during iteration; handoff gate from AGENTS.md's Verification Policy.

## Technical Evidence

Pending execution.

## Drift Notes

None yet.
