---
id: E46-agent-workflow-assets/T006-check-in-a-fresh-session-and-record-the-proof
title: Check in a fresh session and record the proof
status: done
objective: Add the savepoint-check skill so an independent session writes an immutable Check and Issues, and may close work only under the recorded rules.
depends_on:
    - E46-agent-workflow-assets/T001-write-the-shared-checking-method-once
    - E46-agent-workflow-assets/T005-execute-a-plan-and-say-what-actually-happened
complexity_tier: medium
complexity_reason: One mirrored skill pair, but it holds closure authority, owner acceptance, and the no-repair boundary.
---

# T006: Check in a fresh session and record the proof

## Problem

The checker is the only role that can turn work into a completion, so its authority has to be bounded on both sides. It must be able to close a technical Task when clearance is current and nothing material blocks it — otherwise every project waits on the owner for routine work. It must not be able to manufacture that clearance: no repairing the implementation and then passing it, no rewriting acceptance criteria to match what was built, no updating Design as a form of remediation.

Freshness is the other boundary and it is structural rather than stylistic. A Check written in the executor's conversation is not independent, whatever it says. The same model is fine; the same session is not.

The output is an immutable record, one file per run, with a new `C###` and `supersedes` on a recheck. Corrections go back to the planner or executor and a later Check verifies the repair — the record of the failed run stays.

## Context Files

- `agent-skills/savepoint-check/SKILL.md`
- `templates/project/agent-skills/savepoint-check/SKILL.md`
- `agent-skills/references/check-method.md`
- `agent-skills/savepoint-audit-task/SKILL.md`
- `agent-skills/savepoint-audit-epic/SKILL.md`
- `internal/init/agent_skills_test.go`
- `internal/data/check_v2.go`
- `internal/data/evidence_v2.go`
- `internal/data/objective_gate_v2.go`
- `.savepoint/releases/v2/v2-Design.md`

## Acceptance Criteria

- [x] `agent-skills/savepoint-check/SKILL.md` exists with frontmatter `name: savepoint-check`, passes the existing structure validation, and has a byte-identical shipped copy.
- [x] Its trigger is router `state: check`, and it states the fresh-session requirement: independent from the executor conversation, same model permitted, model names optional.
- [x] It loads `agent-skills/references/check-method.md` in full and does not restate the method.
- [x] Its write boundary names the Check record, Issues, evaluation metadata, and authorized closure only; it forbids repairing implementation, editing acceptance criteria, and updating Design as remediation.
- [x] It carries the Check artifact template matching `internal/data/check_v2.go`: `id`, `scope: {kind: task|objective, id}`, `result`, `checked_by`, `executed_session`, `checked_at`, `reviewed` with files and dependencies, `issues`, `supersedes`.
- [x] It states that each run writes a new immutable `C###` record, that a recheck sets `supersedes`, and that no record is edited after the fact.
- [x] It distinguishes the two scopes: a Task Check evaluates one Task's outcome and evidence; an Objective Check adds integration across owned Tasks and Design reconciliation.
- [x] It states the closure rules: a checker may complete a technical Task with current clearance and no unexcepted material blocker, while a Task with `owner_validation.required` needs recorded owner acceptance naming the accepted Check.
- [x] It states that a record lacking sufficient scope or evidence cannot support completion, and that stale or unknown freshness blocks normal completion rather than being re-assessed by assertion.
- [x] It requires a failed Check to record `NEEDS WORK` with Issues and hand remediation back to the executor or planner, with the executor resuming at `stage: build` inside the same Task.
- [x] It states that advisory observations, including `STYLE` guardrail rules, are non-blocking and are recorded as such.
- [x] `internal/init/agent_skills_test.go` asserts the fresh-session rule, the write boundary and forbidden actions, the Check template fields, the two scopes, and the conditional-acceptance closure rules.
- [x] The V1 audit skills and `agent-skills/references/audit-method.md` are unchanged and still present in both trees.

## Implementation Plan

- [x] Read design sections 4, 5, and 8, plus `internal/data/check_v2.go` and `evidence_v2.go`, to match the recorded field names exactly.
- [x] Write `agent-skills/savepoint-check/SKILL.md` with Purpose, Trigger, Read, Workflow, Rules, and the Check artifact template.
- [x] Write the closure rules, including the conditional owner acceptance path and the stale/unknown freshness block.
- [x] Write the NEEDS WORK path: Issues recorded, remediation handed back, record retained, later Check verifies.
- [x] Point at the shared method for scope locks, coverage, adversarial pass, materiality, and re-audit convergence instead of repeating them.
- [x] Mirror the file to the template tree.
- [x] Extend `internal/init/agent_skills_test.go` with the savepoint-check contract case.
- [x] Run `go test ./internal/init/... ./internal/data/...`, then `make build && make test`.

## Context Log

**Dependency note:** T006 depends on T005 (`savepoint-task`), which is `status: in_progress`/`stage: audit` — built and handed off, but not yet marked `done` by the owner (only the owner may do that). Flagged to the user before starting; the owner explicitly chose to proceed with T006 now rather than wait, since T005's implementation is already in the working tree and its own gates pass.

**Files read (Context Files budget):** `agent-skills/references/check-method.md`, `agent-skills/savepoint-audit-task/SKILL.md`, `agent-skills/savepoint-audit-epic/SKILL.md`, `internal/init/agent_skills_test.go`, `internal/data/check_v2.go`, `internal/data/evidence_v2.go`, `internal/data/objective_gate_v2.go`, `.savepoint/releases/v2/v2-Design.md`. `agent-skills/savepoint-check/SKILL.md` and `templates/project/agent-skills/savepoint-check/SKILL.md` did not exist yet — this task creates them.

**Extra reads beyond the Context Files budget, logged with reason:** none this task. Router, epic detail, `AGENTS.md`, `agent-skills/savepoint-design/SKILL.md`, `internal/init/skill_validation_test.go`, `internal/data/gate_v2.go`, and `internal/data/task_v2.go` were already read and current from the prior T005 session in this same conversation, so no new read was needed to reconfirm them.

**Files written:**
- `agent-skills/savepoint-check/SKILL.md` (new) and `templates/project/agent-skills/savepoint-check/SKILL.md` (new, byte-identical mirror).
- `internal/init/agent_skills_test.go` — added `TestSavepointCheckSkillPassesStructureValidation`, `TestSavepointCheckSkillFreshSessionRule`, `TestSavepointCheckSkillLoadsMethodWithoutRestating`, `TestSavepointCheckSkillWriteBoundaryAndForbiddenActions`, `TestSavepointCheckSkillArtifactTemplate`, `TestSavepointCheckSkillTwoScopes`, `TestSavepointCheckSkillClosureRules`, `TestSavepointCheckSkillNeedsWorkPath`, `TestSavepointCheckSkillAdvisoryObservationsNonBlocking`, `TestSavepointCheckSkillLiveAndTemplateMatch`, `TestV1AuditSkillsAndSharedMethodPresentInBothTrees`.

**Per-criterion outcome:** every acceptance criterion above is backed by a passing assertion in `internal/init/agent_skills_test.go` (see test names above) except the structure/frontmatter/byte-parity criterion, which is covered by both the new `TestSavepointCheckSkillPassesStructureValidation`/`TestSavepointCheckSkillLiveAndTemplateMatch` and the pre-existing generic `TestSavepointSkillsHaveValidFrontmatter`, `TestSavepointSkillsHaveNonEmptyTriggerAndWorkflow`, and `TestProjectGuidanceTemplatesMirrorLiveGuidance` (auto-discovers `savepoint-*` dirs, so it now also covers `savepoint-check`). V1 audit skills and shared method unchanged: `git diff --stat` for `agent-skills/savepoint-audit-task`, `agent-skills/savepoint-audit-epic`, `agent-skills/references/audit-method.md` and their template mirrors is empty.

**Commands run:**
- `go test ./internal/init/... -run 'SavepointCheck|V1AuditSkills' -v` — all 11 new/targeted tests pass.
- `go test ./internal/init/... ./internal/data/...` — both packages pass (no regressions from the new directory or the dynamic skill-count parity test).
- `make build && make test` — build succeeds; full suite passes across all packages.

**Limitations:** No independent Check has been run against this Task itself; per the skill just written, this session cannot be that Check. `internal/data/check_v2.go`, `evidence_v2.go`, and `objective_gate_v2.go` are read-only context here — this task does not wire the V2 skill into any V2 runtime code (none exists yet on this branch), only ships the skill guidance and its contract tests, matching the epic's stated scope. T005 is still awaiting the owner's `done` mark; that is a pre-existing limitation carried forward, not something this task changes.
