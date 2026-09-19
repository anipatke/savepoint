---
name: savepoint-task
description: Executes one Savepoint Task within its planned boundaries when router state is task, recording extra reads, lifecycle progress, and handoff evidence for a fresh savepoint-check session, and returning REPLAN REQUIRED on a material gap instead of redesigning.
---

# Savepoint Skill: Task

## Purpose

Build exactly one Task within the boundaries the planner already set, and leave behind a truthful record of what happened. This is the role with the most room to quietly lie: widening scope and calling it necessary, redesigning around an inconvenient plan without saying so, ticking acceptance criteria that were never actually verified, or granting itself clearance. This skill closes those off structurally: it can advance a Task's lifecycle and record evidence, but it can never write the Check that clears that evidence, and a materially invalid plan produces `REPLAN REQUIRED` rather than an improvised rewrite.

## Trigger

Use this skill when router `state` is `task`.

## Read

- `.savepoint/router.md`
- The active Task
- The owning Objective's boundaries
- The Task's `## Context Files`
- Applicable policy: the `.savepoint/Guardrails.md` rule IDs the Task names, when the project has that file

These Context Files are the read budget. Any read beyond them is an extra read: allowed, but logged with what was read and why, so replanning frequency and context growth stay measurable instead of anecdotal.

## Workflow

1. Confirm the start is allowed: the Task's own dependencies are satisfied and its owning Objective is ready. A blocked start is reported to the planner, never worked around by starting anyway or substituting a different Task.
2. Set the Task `status: in_progress` and `stage: build`.
3. Implement the plan's checklist in scoped order, writing code that follows the `STYLE` guardrail rules where the project defines them.
4. Advance the lifecycle as work completes: `build` → `test` → `audit`. Reaching `audit` means the Task is ready for a Check, and explicitly does not mean it passed.
5. Before editing anything outside the Context Files, record the extra read and its reason in the Task's evidence.
6. If the plan turns out to be materially invalid — a Context File doesn't exist, an assumption the plan depends on is false, the described approach can't work — stop and return `REPLAN REQUIRED` instead of redesigning silently. See below.
7. At handoff, verify every acceptance criterion against a concrete outcome, run `make build && make test`, and record the required technical evidence.
8. Hand off to a fresh `savepoint-check` session. The executor's own session can never be that Check.

## Write Boundary

This skill may write: scoped implementation for the active Task, recorded evidence (extra reads, per-criterion outcomes, command results, limitations), lifecycle progress (`status` and `stage`), and a replan handoff when one is needed.

It must never: edit the Task's acceptance criteria to match what was actually built, write a Check record, close an Issue, or claim clearance or owner acceptance for its own work. Those are the checker's and owner's authority, not the executor's.

## Lifecycle

- **Start:** `planned` → `in_progress` with `stage: build`. Requires satisfied Task dependencies and a ready owning Objective.
- **Verify implementation:** `stage: build` → `test` → `audit`, recording acceptance-criterion evidence and required command results along the way. `audit` means ready for Check — it is never recorded or described as passed.
- **Replan:** keep the current `status` and `stage`; set the replan reason with handoff evidence; preserve partial work; stop for the planner.
- **After a Check:** a `NEEDS WORK` Check resumes repair at `stage: build` within the same Task. A `CLEAR` Check does not close the Task by itself — completion and `status: done` are the checker's or owner's action, never something this skill sets for itself.

## Extra Reads

The Task's Context Files are the budget, not a suggestion. When a targeted verification read is genuinely necessary beyond them, take it, but record in the Task's evidence what was read and why it was needed. An unlogged extra read defeats the point: it hides how often plans actually hold up against real context growth.

## REPLAN REQUIRED

When the plan is materially invalid, this skill never redesigns around the problem and never quietly ships a different approach than the one the owner and planner agreed to. Instead it:

- sets the replan reason with the handoff evidence the planner needs to understand what broke;
- keeps the Task's current `status` and `stage` unchanged;
- preserves whatever partial work already exists;
- stops, and hands control to the planner (`savepoint-design`) rather than continuing.

A material gap is not a routine surprise to route around — a Context File that doesn't exist, a dependency whose actual interface contradicts the plan, an acceptance criterion that turns out to be technically impossible as written. Stopping honestly here is the success case, not a failure to route around.

## Issue Capture

Enter Issue capture as an entry from this workflow when implementation hits a
defect, drift, or guardrail gap that is not the work this Task owns. Record it
as an Issue rather than expanding scope or repairing it silently; see
`agent-skills/references/issue-capture.md` for the artifact template and
rules. This skill may add evidence to an Issue; only `savepoint-check` may
close one.

## Evidence And Handoff

At handoff, the Task's recorded evidence must include:

- a per-criterion outcome for every acceptance criterion;
- the named commands actually run, including `make build && make test`;
- the files read and the files changed, including every logged extra read;
- stated limitations — anything not verified, or verified only partially.

This evidence is what a fresh `savepoint-check` session will treat as claims to verify, not as proof by itself; see `agent-skills/references/check-method.md` for what that session does with it. Handoff always goes to that fresh session — this skill's own session, having built the Task, can never be the Check that clears it.

## Rules

- Stay within the active Task's scope; do not widen it and call the extra work necessary without a replan.
- Do not edit acceptance criteria to match what was built.
- Do not write a Check record, close an Issue, or claim clearance or owner acceptance for this Task's own work.
- A blocked start (unsatisfied Task dependency, an owning Objective that is not ready) is reported, never worked around.
- Every read beyond the Task's Context Files is logged with what was read and why.
- A materially invalid plan returns `REPLAN REQUIRED` with preserved partial work and unchanged `status`/`stage`; it is never silently redesigned.
- Treat the `STYLE` guardrail rules as advisory: they shape the code written, but do not block handoff on their own, and this skill references guardrail rule IDs rather than restating rule prose.
- Handoff always names a fresh `savepoint-check` session as the next step; the executor's own session is never that Check.
- Use `state` only for router phase, Task `status` only for Task lifecycle, and `stage` only when the Task is `in_progress`.
