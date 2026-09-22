---
id: I019
title: Objective Check repair must not require reopening completed Tasks
type: defect
status: open
source:
  kind: report
  actor: {role: owner, session: user}
  at: '2026-09-22T09:36:11Z'
checks: [C906]
severity: high
history:
  - at: '2026-09-22T09:36:11Z'
    actor: {role: owner, session: user}
    kind: observed
    note: The owner rejected the workflow requirement to retreat a completed Task before repairing an Objective Check finding; completed Tasks must remain closed and remediation must be represented as new work linked to the Objective and Issue.
---

# I019: Objective Check repair must not require reopening completed Tasks

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
  to the owner. Together, these rules made C906 tell the owner to reopen T005
  even though I018 is an Objective integration finding and T005's completed
  history remains valid.
- The owner explicitly rejected that required interaction on
  2026-09-22: link the finding to O012, keep O012 open until repaired, and do
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
