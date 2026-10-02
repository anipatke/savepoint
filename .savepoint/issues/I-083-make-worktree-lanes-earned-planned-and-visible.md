---
id: I-083
title: Make worktree lanes earned, planned, and visible on the board
type: other
status: resolved
source:
  kind: report
  actor: {role: owner, session: user}
  at: '2026-09-27T00:14:23Z'
severity: low
resolution:
  disposition: escalated
  actor: {role: planner, session: worktree-support-design-20260927}
  at: '2026-09-27T04:17:21Z'
  reason: Owner promoted the follow-up into the Worktree Support Goal and Objective O-033; detailed Task design remains intentionally deferred.
escalated_to: O-033
history:
  - at: '2026-09-27T00:14:23Z'
    actor: {role: owner, session: user}
    kind: observed
    note: >-
      Raised from a review of the templates/project-v2 skills at 4f17da5
      against three principles: sequential by default, respect dependencies,
      and highlight parallel opportunities. Follows I-074. Owner wants this
      solved as V2 follow-up work, keeping v2.1 (R-007) focused on Code
      Health. Needs more design thinking first, including how lanes appear
      in the board UI. Not ready for Task detailing.
  - at: '2026-09-27T04:17:21Z'
    actor: {role: planner, session: worktree-support-design-20260927}
    kind: escalated
    note: >-
      Promoted into G-001 Worktree Support and O-033. The owner requested an
      Objective-only plan so Task detailing waits for further design discussion.
---

# I-083: Make worktree lanes earned, planned, and visible on the board

## Summary

I-074 added lane guidance with no new fields, statuses, router keys, or
commands. Running lanes since then shows the guidance leans toward
parallelism ("where practical") instead of making it the exception, tests
independence with the wrong signal, and leaves the lane plan only in chat.
The owner has no generated `Next` line for a second lane and no board view of
what is running where.

This Issue deliberately revisits I-074's "no new fields" boundary. Target:
V2 (R-006) follow-up, not v2.1: R-007 stays focused on Code Health. R-006 is
already `done`, so promoting this to an Objective needs an owner Goal
decision (reopen R-006, or a V2.0.x continuation Goal).

Intended principles:

1. Sequential by default. One main lane; a separate lane is earned by passing
   an explicit test, never assumed.
2. Respect dependencies, including implicit ones through shared writes.
3. Highlight parallel opportunities in the plan, in `resume`, and on the
   board.

## Evidence

- `savepoint-design` step 7 says "Where practical, shape Tasks so independent
  ones can run side by side", which rewards splitting rather than requiring
  independence to be proven.
- Step 7's independence test is "no overlapping Context Files", but Context
  Files are the read budget. Shared reads are harmless; shared writes cause
  conflicts. Nothing records a Task's planned write set.
- Serialization points are unnamed. The Task template's own worked example
  writes `src/cli.ts` (dispatch), a file most command Tasks touch.
- Step 5 names lanes only in the owner review in chat; nothing persists, so a
  replan or a later session loses the lane plan.
- The router selects one Task. `resume` produces one `Next` line, so the
  owner hand-writes the prompt for every other lane.
- A lane sets `status: in_progress` only on its branch. Main still shows the
  Task `planned` and can select it again.
- `savepoint-task` step 1 checks dependencies against the lane's copy of the
  Task files. A dependency merged to main after the lane branched is invisible
  there: a false block, or a base missing the contract it builds on.
- AGENTS.md Worktree Lanes has no merge protocol: no merge order, no gate run
  on main between merges, and no rule for a `REPLAN REQUIRED` in one lane that
  invalidates a sibling lane.
- `check-method.md` has no adversarial question for lane-built Objectives
  (duplicated helpers, divergent conventions, conflicts resolved silently at
  merge).

## Candidate Direction

Recorded to seed planning, not yet confirmed by the owner:

- Optional Task frontmatter `lane:`; absent means the main sequential lane.
  Order within a lane stays `depends_on`. One field, no second Task list.
- Lane test the planner must pass to assign a separate lane: no `depends_on`
  path to other lanes; disjoint planned writes; no serialization point
  written; every shared contract already landed on main.
- Split Context Files into read and write (tagged entries or a planned
  `## Files Changed` section).
- A Serialization Points list in `Design.md`.
- Contract first, fan out, join: shared interfaces in a main-lane foundation
  Task; lane Tasks depend on it with `requires: clear`; the Full Objective
  Check is the join.
- Claim on main before launching a lane; lanes branch from main after their
  dependencies merge.
- Merge in dependency order; run configured gates on main after each merge;
  failures become Issues captured on main.
- A lane `REPLAN REQUIRED` stops that lane; the planner checks sibling lanes
  on main.

## Open Questions

- Is `lane:` a planner-chosen label or derived by the runtime from
  `depends_on` plus write sets? Deriving avoids drift; labelling keeps owner
  intent explicit.
- Board UI: how should lanes appear in a meaningful way? Options to weigh:
  - a lane badge or grouping on Task cards, with main lane visually distinct;
  - a "Ready in parallel" group showing Tasks that are unblocked now;
  - a copyable `Next` line per ready lane;
  - claimed / running-in-worktree / merged indicators, possibly detected
    from `git worktree list` and branch state rather than stored;
  - a warning when two in-flight lanes plan writes to the same file or a
    serialization point.
- Should `resume` list all ready lanes, or keep one `Next` and add a separate
  read-only lanes view to preserve the single-next-action model?
- How does claiming work without breaking "agents may set `in_progress` when
  starting": an owner board action, or a planner write on main?
- Does lane state belong in Task records at all, or only in the board's live
  view of git?

## Proof Needed

To be finalised when this is planned. At minimum:

- Guidance states sequential by default, and a separate lane requires the
  lane test.
- Independence is judged on planned writes and serialization points, not
  reads.
- The lane plan survives beyond the chat review and after a replan.
- The owner can start every ready lane without hand-writing its prompt, and
  can see on the board which lanes are ready, claimed, running, and merged.
- A Task running in a lane cannot be selected again on main.
- Canonical skills and `templates/project-v2` copies stay byte-identical;
  guidance tests cover new phrases; `make build && make test-fast` pass.
