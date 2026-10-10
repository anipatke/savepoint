---
id: T-112
title: Keep only repo-specific rules in our AGENTS.md
objective: O-039
status: planned
depends_on: []
complexity_tier: low
complexity_reason: One project-owned file edited above the managed block, plus a parity check.
owner_validation: {required: false}
planned_by: {role: planner, session: plan-o039-20261006}
---

# Keep only repo-specific rules in our AGENTS.md

## Outcome

This repository's AGENTS.md has no section above `<!-- SAVEPOINT:BEGIN -->` that repeats the managed block. What remains there is repository-specific: make gate commands, the Codebase Map, runtime gate code pointers, the build-from-source rule for `savepoint`, and Reporting to the Owner.

## User Check

Open AGENTS.md and confirm that Workflow, Skill Activation, Router Selection, Verification Policy, Required Goal Context, Terminology, Issue Capture, and Worktree Lanes each appear once, inside the managed block.

## Done When

- Every section above the managed block is either repository-specific or removed. Each removed sentence that differs from the managed block is listed in Technical Evidence with where it now lives, or why it was dropped. The `health report` human-only rule is listed for the Task "Give each shared rule one home" to move into the template.
- The managed block's bytes are unchanged (FS-02).
- `savepoint resume` and `savepoint doctor` still run cleanly.
- Token weight of AGENTS.md before and after is recorded.
- `make build && make test-fast` passes.

## Context Files

`AGENTS.md`; `templates/project-v2/AGENTS.md`; `internal/init/agents.go`; `internal/init/agents_test.go`.

## Design References

O-039 Confirmed Design; Design section 1 (V2 agent workflow assets, upgrade-assets).

## Guardrails

TPL-01, TPL-02, TPL-04, FS-02, POL-01, POL-02, TEST-01, TEST-02, TEST-03, TEST-08, STYLE-07, STYLE-10.

## Implementation Plan

1. Record AGENTS.md token weight.
2. Diff each section above the block against its managed counterpart; keep repository-specific text and remove the rest.
3. Note every dropped difference and where it now lives.
4. Run resume, doctor, and the gate; record the after weight.

## Boundaries

No edits inside the managed block or to the template; no rule changes.

## Technical Verification

Focused tests during iteration; handoff gate from AGENTS.md's Verification Policy.

## Technical Evidence

Pending execution.

## Drift Notes

None yet.
