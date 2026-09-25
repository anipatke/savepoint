---
id: I-019
title: Objective Check repair must not require reopening completed Tasks
type: defect
status: resolved
source:
  kind: report
  actor: {role: owner, session: user}
  at: '2026-09-22T09:36:11Z'
checks: [C-906]
severity: high
resolution:
  disposition: accepted
  actor: {role: owner, session: user}
  at: '2026-09-22T21:13:35Z'
  reason: Owner directed closure after visual inspection and waived an independent Issue Check; this accepts the repair without claiming technical CLEAR.
history:
  - at: '2026-09-22T09:36:11Z'
    actor: {role: owner, session: user}
    kind: observed
    note: The owner rejected the workflow requirement to retreat a completed Task before repairing an Objective Check finding; completed Tasks must remain closed and remediation must be represented as new work linked to the Objective and Issue.
  - at: '2026-09-22T00:00:00Z'
    actor: {role: executor, session: v2-chat}
    kind: repair_attempted
    note: Split repair routing by Check scope in agent-skills/savepoint-check/SKILL.md (workflow step 5 and Rules), its byte-identical templates/project-v2 copy (TPL-01), agent-skills/savepoint-task/SKILL.md (new "After a mandatory Objective or Release Check" lifecycle line, mirrored to its templates/project-v2 copy), and .savepoint/Design.md Section 7 step 4 — a Task Check's NEEDS WORK still resumes stage:build inside that Task, but an Objective/Release Check's NEEDS WORK now routes to new or newly selected work linked to the Objective without retreating a done Task. Added TestResolveObjectiveCompletion_repairAndRecheckDoesNotRequireRetreatingDoneTasks in internal/data/objective_gate_v2_test.go proving ResolveObjectiveCompletion already blocks on NEEDS WORK and unblocks on a fresh superseding CLEAR Check while owned Tasks remain done throughout — the gate resolvers required no code change, only the skill/Design instruction text did. Also updated internal/init/agent_skills_test.go's doc-consistency assertions (TestSavepointCheckSkillNeedsWorkPath, new TestSavepointTaskSkillObjectiveCheckNeedsWorkDoesNotRetreatDoneTasks) to match the new split-routing wording. `go build ./...`, `go vet ./...`, and `go test ./...` pass except a pre-existing, unrelated failure, since repaired separately as I-028 — TestSharedIssueCaptureRoleBoundariesAndRepairRouting expected agent-skills/references/issue-capture.md to contain the phrase "becomes a new, bounded Task in an Objective", which was absent on this branch before and after this repair (verified via git stash) — out of scope for this Issue and not touched here.
  - at: '2026-09-22T20:52:28Z'
    actor: {role: executor, session: codex-i019}
    kind: repair_attempted
    note: Recorded the implementation plan and extended the Objective completion regression to add a remediation Task while prior Tasks remain done, require a superseding CLEAR Check, and block closure until a material Issue linked to that current Check is resolved. ResolveObjectiveCompletion now enforces that Issue blocker and refuses an unrelated Objective exception; the Check skill, V2 scaffold copy, and Design describe the same closure rule. Renamed the shared Check-to-Issue helper for Objective and Release use. Focused regression and affected package tests passed. make build passed; make test failed only in TestIntegration_InstallDependencies because WSL resolved Windows npm from a UNC path. make build and go test ./... -skip '^TestIntegration_InstallDependencies$' passed. git diff --check passed. Fresh independent Check is still required; I-019 remains open.
  - at: '2026-09-22T21:07:20Z'
    actor: {role: executor, session: codex-i019}
    kind: repair_attempted
    note: Reran the complete configured quality gate with Linux Node v22.22.2 and npm 10.9.7 prepended to the WSL PATH. make test completed successfully across all packages, including internal/init's npm integration test and internal/migrate. The prior failure was caused by the command session resolving Windows npm; no test was skipped in this run. I-019 remains open pending a fresh independent savepoint-check session.
  - at: '2026-09-22T21:13:35Z'
    actor: {role: owner, session: user}
    kind: owner_decision
    note: Owner explicitly directed I-019 closure after visual inspection, waived a separate Issue Check, and accepted the implemented repair. The full make build and make test gates passed; this accepted disposition is not an independent technical CLEAR.
---

# I-019: Objective Check repair must not require reopening completed Tasks

## Summary

The active Check workflow says every `NEEDS WORK` result returns the executor
to `stage: build` inside the same Task. When a mandatory Full Objective Check
finds a cross-Task or integration problem after all owned Tasks are already
done, that rule forces the owner to manually retreat a completed Task before
any repair can begin.

That is unnecessary lifecycle bookkeeping. Completed Tasks are historical
records of work the owner already closed. An Objective-level finding should
block the Objective through its linked Issue and current Check, while repair is
captured as new remediation work under the Objective. The owner should not
have to reopen old Tasks merely to satisfy the agent router.

## Evidence

- `agent-skills/savepoint-check/SKILL.md:38` and line 114 require the executor
  to resume at `stage: build` inside the same Task for every `NEEDS WORK`
  result, without distinguishing Task Checks from Objective/Release Checks.
- `.savepoint/Design.md:144` repeats the same unconditional repair route.
- `AGENTS.md:55` and `.savepoint/Design.md:105` reserve retreating a Done Task
  to the owner. Together, these rules made C-906 tell the owner to reopen T-005
  even though I-018 is an Objective integration finding and T-005's completed
  history remains valid.
- The owner explicitly rejected that required interaction on
  2026-09-22: link the finding to O-012, keep O-012 open until repaired, and do
  not require micromanagement of completed Tasks.

## Proof Needed

- Split repair routing by Check scope: a failed Task Check may return its
  active Task to build when appropriate, while an Objective or Release Check
  records linked Issues and creates or selects new remediation work without
  retreating completed Tasks.
- Reconcile the canonical Check skill, Design, active AGENTS guidance, and the
  byte-identical V2 scaffold skill copy required by TPL-01.
- Add workflow tests proving a failed Objective Check can be remediated and
  rechecked while all previously completed Tasks remain `done`.
- Confirm Objective closure still requires every linked material Issue to be
  resolved or explicitly accepted and a fresh current `CLEAR` Objective Check.

## Implementation Plan

1. Preserve the scope-specific repair instructions in the Check and Task skills, Design, and V2 scaffold copies. Keep completed Tasks closed; route an Objective finding to a new or newly selected Task under that Objective, linked through the Issue's `tasks` and `checks` lists.
2. Extend the Objective completion regression to exercise a new remediation Task after `NEEDS WORK`: the Objective remains blocked while it is unfinished, then requires a fresh current `CLEAR` Check after repair. Assert the original Tasks remain `done` at every step.
3. Enforce resolution of material Issues linked to the current Objective Check in the canonical Objective completion gate, alongside owned Task completion and current clearance. A checker may verify the Issue, or the owner may explicitly accept it; a direct Objective link in the Issue schema is outside I-019.
4. Run build and test gates and append repair evidence here. A fresh independent Check can close the Issue as `verified`; an explicit owner decision can instead close it as `accepted` without claiming technical `CLEAR`.

## Owner Disposition

The owner directed an `accepted` resolution after visual inspection and waived a separate independent Issue Check. This closes I-019 by owner decision and does not record technical `CLEAR` or replace any mandatory Objective or Release Check.
