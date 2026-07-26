<!-- FROZEN FIXTURE: a bundled skill a project has locally tailored, kept
     byte-for-byte. Upgrade must keep it and offer the incoming version beside
     it, or back it up before replacing under the pre-manifest migration. -->

---
name: savepoint-build-task
description: Executes Savepoint task-building work.
---

# Savepoint Skill: Build Task

## Workflow

1. Read the task acceptance criteria and implementation plan.
2. Implement the checklist in order.
3. Run our own gate: `make check`.
4. Post the summary in #eng-builds before handing back.
