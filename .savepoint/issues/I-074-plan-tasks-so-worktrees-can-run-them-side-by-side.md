---
id: I-074
title: Plan Tasks so worktrees can run them side by side
type: other
status: open
source:
  kind: report
  actor: {role: owner, session: user}
  at: '2026-09-26T04:33:30Z'
severity: low
history:
  - at: '2026-09-26T04:33:30Z'
    actor: {role: owner, session: user}
    kind: observed
    note: >-
      After splitting I-069..I-072 across two Herdr worktrees by hand, the
      owner asked that Objective breakdowns prefer Tasks that worktrees can
      deliver in parallel where practical. Today the planner orders Tasks but
      never says which can run side by side, and every lane prompt had to
      restate the worktree rules.
---

# I-074: Plan Tasks so worktrees can run them side by side

## Summary

When `savepoint-design` details an Objective's Tasks, it splits by outcome
(step 7), records `depends_on`, and lists exact Context Files, but it never
considers whether Tasks could run at the same time in separate worktrees.
The router then hands out one Task at a time. The owner wants breakdowns to
prefer parallel delivery where practical, and to be told plainly which Tasks
can run together.

The ingredients already exist: Tasks with no `depends_on` path between them
are independent, disjoint Context Files predict clean merges, and all Task IDs
are allocated by `savepoint create-task` during planning on the main branch, so
worktree lanes need not create identities. What is missing is guidance, not
machinery. Keep it small: no new fields, statuses, router keys, or commands.

## Evidence

- `agent-skills/savepoint-design/SKILL.md` Workflow steps 5, 7, and 11: Tasks
  are split by outcome and the router selects the "first unblocked planned
  Task"; nothing addresses parallel lanes or file overlap.
- The skills and AGENTS.md do not mention worktrees or branches. Worktree
  hazards found while splitting I-069..I-072:
  - `.savepoint/router.md` holds one selection that `savepoint-task` writes
    at Task start and after an Issue repair, so parallel lanes conflict on it.
  - `.savepoint/task-ids.yml` is tracked per checkout and its lock covers only
    one folder, so `create-task` in two worktrees would allocate the same
    `T-###`.
  - `C-###` and `I-###` have no allocator, so Checks or Issues created in
    parallel lanes collide. Collisions fail closed as duplicate IDs at merge.
- The lane prompts for `fix/i072-doctor-active-goal` and
  `fix/i069-i071-migrate` had to restate these rules by hand.

## Proof Needed

- `savepoint-design` guidance says: where practical, shape an Objective's
  Tasks so independent ones have no `depends_on` path between them and no
  overlapping Context Files; keep genuinely shared work in one sequential
  lane rather than forcing a split.
- The plan presented to the owner names the lanes in plain words, for example
  "T-101 and T-102 can run side by side in worktrees; T-103 waits for both."
- AGENTS.md and `savepoint-task` state the worktree lane rules once: in a
  worktree lane, do not edit `.savepoint/router.md`, do not create Tasks,
  Checks, or Issues (record them as a note for the owner instead), commit on
  the lane branch, and run the Full Objective Check on the main branch after
  merging.
- Canonical skills and their `templates/project-v2` copies stay byte-identical;
  guidance tests cover the new phrases; `make build && make test-fast` pass.
