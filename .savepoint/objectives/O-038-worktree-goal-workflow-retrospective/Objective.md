---
id: O-038
title: Review the optional parallel workflow
status: in_progress
depends_on: [O-037, O-033]
release: G-001
priority: medium
rank: 3
---

# O-038: Review the optional parallel workflow

## Outcome

The planner reviews this Goal's workflow lessons and records justified improvements, or an explicit no-change conclusion with reasons, after Advanced Options and parallel suggestions are complete.

## Why

The design skill requires a final workflow retrospective for each Goal; the new advisory boundary should remain simple and preserve owner choice.

## Success Conditions

- Review the Goal's Task/replan/Check/Issue evidence against canonical skills, references, AGENTS routing and scaffold documents.
- Record justified improvements or no change with reasons, preserving optional advice, core Code Health, owner authority and independent verification.
- Keep canonical/scaffold guidance consistent and verify any actual changes through configured gates and the mandatory Full Objective Check.

## Architectural Considerations

This is the final Goal member required by savepoint-design. It adds no feature framework, lifecycle state or Goal-owned Task list. Detailed Tasks wait until preceding work and lessons exist.

## Boundaries

**In scope:** evidence-backed workflow retrospective and bounded guidance improvements.

**Out of scope:** unrelated product changes, mandatory lanes, Code Health deactivation and detailed speculative Tasks before the dependencies complete.

## Retrospective Review — 2026-10-04

Evidence reviewed for G-001: O-033 (T-102–T-107) and O-037 (T-108–T-110) Task records; Checks C-963 and C-965 (NEEDS WORK) and rechecks C-964 and C-966 (CLEAR); Issues I-131–I-140, with I-074, I-083 and the earlier I-115 as lessons carried in; the canonical skills, shared references, AGENTS.md routing and `.github/workflows/ci.yml`. This repository owns the packaged skills, so a change to them would also have to keep the scaffold copies identical (TPL-01). None is needed.

Findings:

- No Task returned REPLAN REQUIRED. The one planning churn was owner-directed during O-033 planning (opt-in behavior, Code Health kept core, Advanced Options split into O-037 and reordered first) and was absorbed by revising planned Tasks in place through `create-task` identities.
- Both Full Objective Checks first returned NEEDS WORK. Most findings were ordinary edge-case defects (I-131, I-132, I-135, I-137, I-139) covered by existing rules such as DATA-01 and DATA-03, and all were repaired and verified on recheck.
- Missing native Windows CI evidence recurred at three consecutive Full Checks (I-115 for O-032, I-134 for O-037, I-140 for O-033). The shared root cause is the push trigger's hard-coded branch list in `.github/workflows/ci.yml`: each new working branch (v2.1, then v2.20) produced no native run until it was added by hand. The existing planner rule to name who produces platform evidence did not prevent this.
- I-136 found authored lane titles passing raw terminal control sequences to the interactive board. No Guardrail covered that class.

## Confirmed Retrospective Decisions — 2026-10-04

Owner confirmed through the Objective Decision Interview:

1. **CI on every pushed branch.** The CI push trigger runs on all branches rather than a hard-coded list, so native `windows-tests` evidence exists for any pushed working branch. The pull-request trigger is unchanged. Tag pushes stay with `publish.yml`. This is delivered by one implementation Task.
2. **Guardrail ARCH-05 (Required).** Text from project files or external tool output must have terminal control sequences removed before any terminal surface displays it. Written into `Guardrails.md` by the planner on 2026-10-04.

No change, with reasons:

- **Copied worktree instructions lacking a routable Start line (I-138).** Repaired with a regression test. It occurred once, and the AGENTS.md Workflow already defines the routing contract.
- **Owner-directed design revision during planning.** The skill already requires explicit design confirmation. The revision was a legitimate owner choice and was handled without lost work or identity churn.
- **Lane guidance in the planner, executor and checker skills.** No executor confusion, REPLAN REQUIRED or Check finding concerned the guidance itself, and the advisory boundary held.
- No new field, state, command or Goal-owned Task list is added. Optional advice, core Code Health, owner authority and independent verification are preserved.

## Readiness and Verification

Interfaces and ownership are settled: the CI workflow file is repository configuration, and Guardrails are planner-owned. The dependencies O-033 and O-037 are complete. Verification uses `make build && make test-fast` for the Task, a post-push CI run on the working branch showing `ci` and native `windows-tests` triggered by `push`, and the mandatory Full Objective Check with fresh `make test-full` under `agent-skills/references/check-method.md`.
