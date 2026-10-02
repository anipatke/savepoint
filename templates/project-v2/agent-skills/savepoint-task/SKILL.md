---
name: savepoint-task
description: Executes one Savepoint Task within its planned boundaries when router state is task, recording extra reads, lifecycle progress, and handoff evidence for a fresh savepoint-check session, and returning REPLAN REQUIRED on a material gap instead of redesigning.
---

# Savepoint Skill: Task

## Purpose

Build one Task within its planned boundaries and record truthful evidence. A materially invalid plan returns `REPLAN REQUIRED`; this role cannot grant clearance for its own work.

## Goal Context

Every Savepoint project has at least one live Goal selected by the router, and every live Objective names exactly one Goal through `release:`. Goals come from `savepoint init` (G-001), `savepoint migrate`, and the planner, never from this skill. If Next says `Choose a Goal`, or `savepoint doctor` reports a missing Goal or `release:`, report it to the owner; do not pick or create a Goal yourself.

## Trigger

Use this skill when the `Next` line starts with `Start`, `Build`, or `Test`, or router `state` is `task`. A pasted `Next` line is the selection.

## Next

Start from the `Next` line as AGENTS.md's Workflow describes; if `savepoint` is unavailable, follow AGENTS.md rather than guessing. AGENTS.md's Router Selection section says who changes the router.

## Read

- `.savepoint/router.md`
- The active Task
- The owning Objective's boundaries
- The Task's `## Context Files`
- Applicable policy: the `.savepoint/Guardrails.md` rule IDs the Task names, when the project has that file

These Context Files are the read budget. Necessary targeted extra reads are allowed, but logged with what was read and why. The router, `AGENTS.md`, this skill (read directly when the skill tool cannot find it), and the Guardrail IDs the Task names are not extra reads.

## Workflow

1. Confirm the start is allowed: the Task's own dependencies are satisfied and its owning Objective is ready. A blocked start is reported to the planner, never worked around by starting anyway or substituting a different Task.
2. Set the router selection to the owning Objective and active Task, then set the Task `status: in_progress` and `stage: build`, and set the owning Objective `status: in_progress` if it is still `planned`. Follow AGENTS.md's Router Selection section and preserve `release:`.
3. Implement the plan's checklist in scoped order. Treat the `STYLE` guardrail rules as advisory and reference guardrail rule IDs rather than restating rule prose.
4. Start at `in_progress` with `stage: build`, then advance as work completes: `build` → `test` → `audit`. Reaching `audit` means the Task is ready for a Check when one is requested — an optional Task Check or the mandatory Full Objective Check — and explicitly does not mean it passed.
5. Before reading or editing anything outside the Context Files, record the extra read and its reason in the Task's evidence.
6. If the plan turns out to be materially invalid — a Context File doesn't exist, an assumption the plan depends on is false, the described approach can't work — stop and return `REPLAN REQUIRED` instead of redesigning silently. See below.
7. At handoff, verify every acceptance criterion against a concrete outcome, run the applicable gate in Verification Gates below, and record the required technical evidence whether or not an optional Task Check is requested.
8. If the owner requests the optional Task Check, hand off to a fresh `savepoint-check` session. If the owner supplies an explicit Task-check waiver, record that decision and route the evidence to the mandatory Full Objective Check instead. The executor's own session can never be that Check.

In a worktree lane, follow AGENTS.md's Worktree Lanes section: skip the router writes in step 2 and after a direct Issue repair, create no Tasks, Checks, or Issues, and commit on the lane branch without pushing or merging.

## Lifecycle

- A Task Check's `NEEDS WORK` resumes repair at `stage: build` within the same Task. `CLEAR` never sets `status: done`; completion belongs to the owner.
- A mandatory Objective Check's `NEEDS WORK` never retreats a Task that is already `done`. Default to direct repair under its recorded Issue; use new or newly selected work linked to the Objective only when planning is needed. Follow `agent-skills/references/issue-capture.md`, Out-Of-Scope Repair, for proof and the Issue-only router handoff after `repair_attempted`. Preserve `release:` and show the resulting `savepoint resume` Next line.
- Without a requested Task Check, record only an explicit owner waiver using Evidence And Handoff below. Apply AGENTS.md's Verification Policy for its dependency meaning and the mandatory Full Objective Check.

## REPLAN REQUIRED

A materially invalid plan is never silently redesigned. This skill:

- records a `replan:` block in the Task frontmatter, with handoff evidence in the body that the planner needs to understand what broke;
- keeps the Task's current `status` and `stage` unchanged;
- preserves whatever partial work already exists;
- stops, and hands control to the planner (`savepoint-design`) rather than continuing.

`savepoint resume` routes to `Replan` only from this frontmatter block; a REPLAN REQUIRED written only in the body leaves the router on the Task:

```yaml
replan:
  reason: Context File internal/store/lock.go does not exist
  recorded_by: {role: executor, session: <session>}
  recorded_at: '2026-09-19T00:00:00Z'
```

Material gaps include a missing Context File, a contradictory dependency interface, or a technically impossible acceptance criterion. Stop for replanning; do not route around them.

## Issue Capture

Enter Issue capture as an entry from this workflow when implementation hits a
defect, drift, or guardrail gap that is not the work this Task owns. Record it
as an Issue rather than expanding scope or repairing it silently; see
`agent-skills/references/issue-capture.md` for the artifact template and
rules. This skill may add evidence to an Issue; it may record an explicit owner
`accepted` closure but cannot infer acceptance or close one as `verified`.

## Verification Gates

Follow AGENTS.md's Verification Policy: ordinary handoff requires the configured build and test gates; migration/platform-sensitive handoff requires the project's fresh full gate when separately defined. Commands come from `quality_gates` in `.savepoint/config.yml` and the project's Build rules. When none is configured, record that none ran; never invent a gate. Focused tests are iteration aids only. Apply the shared metadata-only reuse requirements exactly.

## Acting On A Code Health Report

A Code Health report is advisory. When it lists a signal, use the thresholds in `.savepoint/health/config.json`:

- Bring a signal that Needs Attention back to the watch line. The aim is not a target; do not chase it.
- Prefer production code that is risky or often changed over tests, and leave flat dispatch tables alone.
- Do not lower a score by moving branches elsewhere; a split must make the code easier to read.
- Record in the Task evidence what remains above the watch line; never silently stop or silently continue.
- Narrowing the measured scope, for example excluding test files, is the owner's decision.

When the project's Guardrails define a complexity rule, cite its ID rather than restating it.

## Evidence And Handoff

At handoff, the Task's recorded evidence must include:

- a per-criterion outcome for every acceptance criterion;
- the named commands actually run, including the applicable build plus fast gate or the full gate;
- the files read and the files changed, including every logged extra read;
- stated limitations — anything not verified, or verified only partially.

Record an owner's Task-check waiver only on the owner's explicit instruction, in this shape (`task` is this Task's own ID); prose in the body does not satisfy the gates:

```yaml
check_waiver:
  task: T-###
  reason: Owner closed this Task without requesting a Task Check.
  actor: {role: owner, session: <session>}
  recorded_at: '2026-09-19T00:00:00Z'
```

A fresh `savepoint-check` session treats this evidence as claims to verify, not proof; the executor's own session can never be that Check. Requested local Checks and the mandatory Full Objective Check use `agent-skills/references/check-method.md`. A waiver satisfies `requires: clear`, never `requires: accepted`, and creates no technical `CLEAR`.

## Rules

This skill may write: scoped implementation, recorded evidence (extra reads, per-criterion outcomes, commands, limitations), lifecycle progress, and a replan handoff.

Never widen scope, edit the Task's acceptance criteria to match implementation, write a Check record, close an Issue on your own, invent a Task-check waiver, or claim clearance or owner acceptance. Record an `accepted` Issue resolution only on explicit owner instruction; reopening is also an owner action. A blocked start is reported, never worked around. Only the owner sets `status: done` or retreats status; apply AGENTS.md's lifecycle terminology.
