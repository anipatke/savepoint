---
id: T-106
title: "Teach the workflow to maintain lane plans"
objective: O-033
status: planned
depends_on: [{task: T-103, requires: clear}]
owner_validation: {required: false}
planned_by: {role: planner, session: codex-o033-planning-2026-10-03}
complexity_tier: medium
complexity_reason: "Planner, executor and checker contracts must stay aligned across canonical/scaffold guidance and safe upgrades."
lane: guidance
planned_reads:
  - "templates/project-v2/.savepoint/config.yml"
  - "agent-skills/savepoint-design/SKILL.md"
  - "agent-skills/savepoint-task/SKILL.md"
  - "agent-skills/savepoint-check/SKILL.md"
  - "agent-skills/references/check-method.md"
  - "templates/project-v2/agent-skills/savepoint-design/SKILL.md"
  - "templates/project-v2/agent-skills/savepoint-task/SKILL.md"
  - "templates/project-v2/agent-skills/savepoint-check/SKILL.md"
  - "templates/project-v2/agent-skills/references/check-method.md"
  - "AGENTS.md"
  - "templates/project-v2/AGENTS.md"
  - "internal/init/upgrade.go"
  - "internal/init/upgrade_test.go"
  - "internal/init/v2_scaffold_test.go"
  - "internal/init/skill_validation_test.go"
  - "internal/init/template_freshness_test.go"
planned_writes:
  - "agent-skills/savepoint-design/SKILL.md"
  - "agent-skills/savepoint-task/SKILL.md"
  - "agent-skills/savepoint-check/SKILL.md"
  - "agent-skills/references/check-method.md"
  - "templates/project-v2/agent-skills/savepoint-design/SKILL.md"
  - "templates/project-v2/agent-skills/savepoint-task/SKILL.md"
  - "templates/project-v2/agent-skills/savepoint-check/SKILL.md"
  - "templates/project-v2/agent-skills/references/check-method.md"
  - "AGENTS.md"
  - "templates/project-v2/AGENTS.md"
  - "internal/init/upgrade_test.go"
  - "internal/init/v2_scaffold_test.go"
  - "internal/init/skill_validation_test.go"
---

# Teach the workflow to maintain lane plans

## Outcome

Fresh and upgraded projects receive coherent lane planning, scope-drift and independent verification instructions while preserving owner-edited assets and old Task records.

## User Check

Review a sample newly planned lane pair and an expanded write scope; inspect upgrade fixtures showing managed guidance refreshed while user-authored guidance and records remain unchanged.

Owner validation is not required for this technical Task; owner Task completion remains required.

## Done When

- Document the Advanced Options setting and exact features.parallel_planning key. No Code Health toggle is introduced. With advice enabled the planner may still plan sequentially or omit optional manifests; disabled mode must not demand lane plans. Malformed or stale advice is nonblocking. Test owner-ignore scenarios in planner/executor/checker guidance.
- Planner guidance reads the project preference and optionally offers exact read/write manifests and stable lane keys, explaining semantic dependencies and bound independence explanations. Ordinary sequential planning and omission of advice remain allowed even when enabled. Replace the old blanket disjoint-Context-Files criterion with the confirmed read/write distinction.
- Executors retain normal scoped execution and extra-read rules. They may note changed advisory scopes but must not require REPLAN REQUIRED solely for changed manifests, ignored groupings or a different worktree choice. Existing actual-worktree preservation/identity rules still apply when a real worktree is used.
- Checker/reference guidance verifies recommendation accuracy as feature behavior. Ignoring a lane or manifest is not a finding or clearance blocker. Do not invent blocking policy, new statuses or executor self-clearance.
- Canonical changed skills/reference are byte-identical to their scaffold counterparts. Root and scaffold routing guidance remain consistent without duplicating safety algorithms.
- Fresh init and existing upgrade fixtures prove new guidance delivery, preservation of edited assets/unknown frontmatter/body and no automatic migration/backfill of lane metadata.
- Per-criterion evidence and the configured handoff gate are recorded. Optional Task Checks require a fresh session; skipping one requires an explicit owner waiver naming this Task, reason, actor and time. All work remains subject to the mandatory Full Objective Check.

## Context Files

- `templates/project-v2/.savepoint/config.yml`
- `agent-skills/savepoint-design/SKILL.md`
- `agent-skills/savepoint-task/SKILL.md`
- `agent-skills/savepoint-check/SKILL.md`
- `agent-skills/references/check-method.md`
- `templates/project-v2/agent-skills/savepoint-design/SKILL.md`
- `templates/project-v2/agent-skills/savepoint-task/SKILL.md`
- `templates/project-v2/agent-skills/savepoint-check/SKILL.md`
- `templates/project-v2/agent-skills/references/check-method.md`
- `AGENTS.md`
- `templates/project-v2/AGENTS.md`
- `internal/init/upgrade.go`
- `internal/init/upgrade_test.go`
- `internal/init/v2_scaffold_test.go`
- `internal/init/skill_validation_test.go`
- `internal/init/template_freshness_test.go`

## Design References

O-033 Confirmed Design — Revised 2026-10-03, including Persisted plan and ownership, Canonical safety projection, Board/resume/handoff, and Guidance and verification. Project Design sections 1, 4, 5, 7, 8 and 9 provide implemented boundaries; Objective deltas govern the planned feature.

## Guardrails

FS-01, FS-02, FS-04, FS-05, DATA-01, DATA-02, DATA-03, DATA-05, TPL-01, TPL-02, TPL-04, ARCH-01, ARCH-02, ARCH-03, CFG-02, CFG-03, TEST-01, TEST-02, TEST-03, TEST-04, TEST-08, TEST-09, STYLE-01..10.

## Implementation Plan

O-037 already delivers the off default and upgrade preservation. Consume and document that behavior without duplicating its configuration implementation or enabling the owner's real project.

1. Use the confirmed Objective contract and predecessor evidence as the documentation source; log an extra read if exact runtime metadata examples need inspection.
2. Update planner Task/Objective examples and readiness guidance, executor non-enforcement rules and checker coverage prompts; retain focused context budgets.
3. Synchronise canonical/scaffold files and managed AGENTS content through existing source conventions.
4. Extend existing fixture-based scaffold/upgrade/skill tests. No production upgrade change unless evidence proves the declared existing upgrade path cannot deliver the assets; return REPLAN REQUIRED for material gaps.

## Boundaries

This Task implements its scoped plan and records execution evidence and authorised lifecycle updates. planned_writes describes anticipated work, not an enforcement allowlist. Necessary targeted reads must be logged before access. A difference from advisory manifests may be noted but is not itself a blocker or reason for REPLAN REQUIRED; existing material acceptance/architecture gaps still use the ordinary replan policy. Do not update another Task, allocate identity records in a lane, automate worktree/branch creation, push, merge, change Goal selection, add lifecycle states, or write a Check. In a lane skip router writes and commit locally. Set the Objective in_progress only when starting its first Task; later lanes must not rewrite the shared Objective simply to restate that status. Design reconciliation belongs to the final integration Task on main, not parallel lanes.

## Technical Verification

make build && make test-fast; explicitly record scaffold parity and preservation tests. No savepoint upgrade-assets agent command; exercise internal test fixtures.

Use focused make test-focused TEST=... only for iteration. Record named happy-path and failure cases for each Done When criterion, actual commands/time/toolchain/results, reads and changes, and limitations. Independent verification follows agent-skills/references/check-method.md; no executor session performs its own Check.

## Technical Evidence

Pending execution. Record per-criterion evidence, extra-read reasons, actual scope versus manifests, gate result, owner validation when required, and any explicit owner Task-check waiver. Planning evidence is not technical clearance.

## Drift Notes

Pending execution. Record material acceptance/architecture deviations under normal policy. Lane choice or anticipated-scope deviation alone is advisory and never a required replan.


## Advisory Feature Contract — Revised 2026-10-03

Parallel planning defaults off via features.parallel_planning. Code Health remains core and is unchanged. Lane metadata and read/write manifests are suggestions: the owner and all skills may ignore them even when enabled. No lane rule feeds existing start, advance, completion or Check gates, and no manifest deviation alone requires replan. Actual dependencies, acceptance criteria, Guardrails and rules for a real worktree remain authoritative. This revised Objective contract supersedes any earlier mandatory lane wording.
