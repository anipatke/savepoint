---
id: T-012
title: Route planners through Task creation
objective: O-019
planned_by: {role: planner, session: i024-allocation-design-20260923}
status: done
complexity_tier: medium
complexity_reason: Canonical and scaffold guidance, upgrades, and agent scenarios must agree on one required path.
depends_on: [{task: T-011, requires: clear}]
owner_validation:
    required: true
    accepted_check: ""
check_waiver:
    task: T-012
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-09-24T20:48:51Z"
---

# T-012: Route planners through Task creation

## Outcome

Planners use the creation operation without selecting an ID. Guidance permits that one command while retaining the ban on other agent-run `savepoint` commands.

## User Check

Review the planning instructions and a two-Objective example. Confirm the planner supplies an ID-free draft, never chooses a number, and may invoke only the narrow Task creation command.

## Done When

- `AGENTS.md` permits `savepoint create-task` only for new Tasks and continues to prohibit other agent-run `savepoint` commands.
- Canonical design skill and byte-identical V2 scaffold copy require the operation; no example directs manual Task ID allocation.
- Fresh-init and upgrade-assets guidance deliver the rule. Strict V2 loading is required after creating or renaming other identity-bearing records.
- A cross-Objective workflow scenario creates Tasks concurrently and fails on ID reuse or an invalid index.
- Command documentation and Design describe implemented behavior and durable allocation state.

## Context Files

`AGENTS.md`, `agent-skills/savepoint-design/SKILL.md`, `templates/project-v2/agent-skills/savepoint-design/SKILL.md`, `agent-skills/savepoint-check/SKILL.md`, `templates/project-v2/agent-skills/savepoint-check/SKILL.md`, `agent-skills/references/issue-capture.md`, `templates/project-v2/agent-skills/references/issue-capture.md`, `internal/init/agent_skills_test.go`, `internal/init/template_freshness_test.go`, `internal/init/v2_scaffold_test.go`, `.savepoint/Design.md`, `README.md`.

## Design References

Design sections 1, 2, 5, and 6; O-019 Architectural Considerations.

## Guardrails

TPL-01, TPL-02, TPL-04, ARCH-03, TEST-01, TEST-02, TEST-04, TEST-05, TEST-08.

## Implementation Plan

1. Replace manual Task-ID guidance with ID-free draft and command instructions. Document the narrow agent exception and its upgrade path.
2. Align canonical and scaffold skills; reconcile Check/Issue guidance only where identity-bearing validation needs it.
3. Add cross-Objective agent scenarios and scaffold/upgrade assertions.
4. Reconcile Design and docs, then run focused tests, `git diff --check`, and `make build && make test`.

## Boundaries

No permission for unrelated CLI commands or duplicated allocation policy across skills.

## Technical Verification

Focused `internal/init` and workflow scenarios, skill byte comparison, `git diff --check`, and `make build && make test`.

## Technical Evidence

### Acceptance Criteria

1. **Agent command rule — met.** `AGENTS.md` and the V2 scaffold allow `savepoint create-task` only for new Tasks from ID-free drafts; other agent-run Savepoint commands remain prohibited. The guide also requires `savepoint resume` after creating or renaming other identity-bearing records.
2. **Planner guidance — met.** Canonical and scaffold `savepoint-design` require the allocator, show an explicit concurrent O-014/O-015 example, and use an ID-free Task template with no numbered heading. Canonical/scaffold Check and Issue guidance require strict V2 loading after their identity-bearing records are written or renamed.
3. **Fresh init and upgrades — met.** Fresh-init output is checked for the narrow exception. The real V2 template upgrade test checks that the same rule reaches the managed AGENTS guide.
4. **Cross-Objective scenario — met.** `TestMainCreateTaskConcurrentAcrossObjectivesStrictLoadsIndex` launches concurrent `savepoint create-task` processes for O-001 and O-002, fails on duplicate IDs, and strict-loads and checks both owners in the resulting V2 index.
5. **Command and architecture documentation — met.** README documents the ID-free invocation and validation. Design records the command, `.savepoint/task-ids.yml`, `.savepoint/task-ids.lock`, locking behavior, retired IDs, strict load, and validation after other identity changes.

### Verification

- `go test ./internal/init` — passed after adjusting the typed artifact-contract test to model the allocator injecting an ID into the ID-free Task draft. The first run exposed stale assertions and they were corrected.
- `go test . -run 'TestMainCreateTaskConcurrentAcrossObjectivesStrictLoadsIndex|TestMainInitScaffoldsV2ProjectWithEmptyValidIndex'` — passed.
- `cmp -s` for the canonical/scaffold `savepoint-design`, `savepoint-check`, and `issue-capture` files — passed; each pair is byte-identical.
- `make build && make test-fast` — passed.
- `make test-full` — passed, including the full suite and Linux, macOS, and Windows build targets.
- `git diff --check` — passed.

### Files Read

Required workflow context: `AGENTS.md`, `.savepoint/router.md`, `.savepoint/Guardrails.md`, `.savepoint/objectives/O-019-allocate-task-identities/Objective.md`, this Task, both design and Check skills, both Issue-capture references, the three listed `internal/init` test files, `.savepoint/Design.md`, and `README.md`.

### Files Changed

- `AGENTS.md`, `templates/project-v2/AGENTS.md`
- `agent-skills/savepoint-design/SKILL.md`, `templates/project-v2/agent-skills/savepoint-design/SKILL.md`
- `agent-skills/savepoint-check/SKILL.md`, `templates/project-v2/agent-skills/savepoint-check/SKILL.md`
- `agent-skills/references/issue-capture.md`, `templates/project-v2/agent-skills/references/issue-capture.md`
- `.savepoint/Design.md`, `README.md`
- `internal/init/agent_skills_test.go`, `internal/init/template_freshness_test.go`, `internal/init/v2_scaffold_test.go`, `internal/init/provenance_contract_test.go`
- `main_test.go` — concurrent command scenario and fresh-init rule assertion
- This Task record — lifecycle and evidence

### Limitations

No optional Task Check was requested or waived, and no Check record was written. The Task remains `in_progress` at `audit` for the owner's User Check and next-step decision; the mandatory Full Objective Check remains required before Objective closure.

### Extra Reads

- `templates/project-v2/AGENTS.md` and targeted `internal/init` managed-guide/upgrade searches — verify the scaffolded CLI rule and real upgrade path.
- `main.go`, `cmd/create_task.go`, `cmd/init_test.go`, `internal/data/task_create.go`, `internal/data/task_ids.go`, and `internal/data/project_test.go`; targeted `cmd`/`internal` searches — confirm CLI syntax, allocator injection, lock/high-water behavior, strict loading, and existing allocator concurrency coverage.
- `.agents/` scenario search returned no fixtures; `main_test.go` supplied the command-level test harness for the concurrent two-Objective scenario.
- `internal/init/provenance_contract_test.go` — adjust typed decoding to simulate the command injecting the allocator-owned ID into the draft.
- `git status --short`, `git diff --stat`, and targeted `git diff` — review the final scope. Other modified paths for T-010/T-011 and allocator implementation were not changed by T-012.

## Drift Notes

If the managed `AGENTS.md` merge cannot safely deliver the exception to existing projects, return REPLAN REQUIRED.
