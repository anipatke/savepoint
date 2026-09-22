---
name: savepoint-design
description: Maintains Savepoint Design and Guardrails and the current Objective, detailing Tasks only for the next ready Objective, when router state is design or an executor returns REPLAN REQUIRED.
---

# Savepoint Skill: Design

## Purpose

Own the project's technical shape at exactly the grain the next unit of work needs: `Design.md` describes implemented reality, `Guardrails.md` holds durable constraints, and one Objective at a time carries detailed Tasks. Everything further out stays a named outcome with Boundaries, not a plan, because detailed plans written ahead of the decisions that shape them are the thing that gets thrown away.

This skill does not write production code, and it does not settle product choices that belong to the owner. Technical readiness does not require the owner to review code.

## Trigger

Use this skill when router `state` is `design`. When an executor returns `REPLAN REQUIRED`, that routes back into this skill's workflow rather than into a separate phase.

## Read

- `.savepoint/router.md`
- `.savepoint/Idea.md`, when it exists, for the original intent
- `.savepoint/Design.md`
- `.savepoint/Guardrails.md`
- The current Objective file
- Targeted implementation evidence needed to settle readiness for the next Objective — read to verify interfaces, ownership, and constraints, not to design speculatively

Read nothing else. Do not detail Tasks for any Objective beyond the next one, and do not open source files outside the evidence a readiness decision actually needs.

## Workflow

1. Read the router, the Idea when present, Design, Guardrails, and the current Objective.
2. Update `Design.md` to describe implemented reality, and `Guardrails.md` to hold durable project constraints; keep Objective deltas as the record of planned change until reconciliation.
3. Keep exactly one Objective active. Objectives beyond it stay named outcomes with Boundaries — no detailed Tasks.
4. Before detailing an Objective's Tasks, check the readiness gate below. Do not detail Tasks for an Objective that is not ready.
5. When the implementation approach for a piece of work is unknown, write a bounded research Task with a named decision deliverable instead of a confident plan the executor will discover is fiction.
6. Split any Task that carries multiple unrelated outcomes or an unresolved architectural decision into separate Tasks.
7. When an executor returns `REPLAN REQUIRED`, treat it as re-entry here: reassess Design, Guardrails, or the Objective as needed, then resume from step 3.
8. Route product choices to the owner instead of inferring them. Technical readiness — settled interfaces, scoped constraints, known dependencies, a verification approach — does not require the owner to review code.
9. When a Task needs a verification approach, name it and reference `agent-skills/references/check-method.md` for how it will later be evaluated; do not restate that method here.
10. When the next Objective's Tasks are detailed and approved, set router `state: task` for the first unblocked planned Task and update `next_action` to execute it with `savepoint-task`.

## Verification Contract

Apply this contract to every implementation, not only to migration work:

- Every Task needs implementation evidence and configured quality-gate results.
- A Task Check is optional. If the owner skips it, record an explicit waiver
  in the Task evidence naming the Task, reason, actor, and time. The waiver is
  not technical `CLEAR` and does not waive acceptance criteria or guardrails.
  It satisfies a Task dependency that requires `clear` — the waiver stands in
  as the owner's own completion decision — but never one that requires
  `accepted`, since there is no Check for the owner to have accepted.
- The Full Objective Check is mandatory before Objective closure. It reviews
  every owned Task, including waived Tasks, cross-Task integration, and
  reconciliation against this Design.
- A Release Check is mandatory whenever a Release exists, and exact owner
  acceptance of its current Check remains required.

## Optional Release Boundary

When the owner chose a Release boundary in Idea, define it as a delivery promise that can be navigated across its member Objectives. A Release is still optional: when no Release is useful, create no Release record and continue through Idea → Design → Task → Check with no missing-record error or extra phase.

For an opted-in Release:

1. Allocate a stable global `R###` identity from the first unused number. Keep that identity stable across title or path edits, never silently reuse it, and fail closed on duplicates.
2. Author the Release sections `Outcome`, `Why`, `Success Conditions`, and `Boundaries`. The outcome describes what the delivery promises, not whether it has been published or deployed.
3. Link each member Objective with one optional `release: R###` field. Derive membership from those Objective records; do not maintain a second membership list.
4. Keep Objectives and Tasks in their normal locations and ownership: a Release does not nest files, own Tasks, or recreate an Objective → Task hierarchy.
5. Treat Release `done` as an integration and owner decision: every member Objective is complete, current CLEAR integration evidence exists, material Issues are resolved or explicitly excepted, and the owner has accepted that exact Check. It does not mean published or deployed.

## Objective Artifact Template

Write the Objective file with this structure:

```markdown
---
id: O###
title: Objective Title
status: planned|in_progress|done
depends_on: [O###]
release: R###
last_check: optional-check-id
freshness:
  state: current|stale|unknown
  check: C###
  assessed_by: role/session
  assessed_at: '2026-09-19T00:00:00Z'
  basis: what was compared to reach this state
---

# O###: Objective Title

## Outcome

The concrete result this Objective delivers.

## Why

Why this Objective matters now.

## Success Conditions

Observable conditions that mean this Objective is done.

## Architectural Considerations

Interfaces, data ownership, and constraints this Objective must respect.

## Boundaries

**In scope:**
- Included work

**Out of scope:**
- Excluded work
```

Omit `release` when an Objective is intentionally unassigned; it is a reference to a first-class Release identity, not free-form release text.

Task membership is derived from which Tasks name this Objective as their owner. Do not also maintain a second, manually kept list of member Tasks in the Objective body — that list drifts from the Tasks themselves and becomes a second source of truth.

## Task Artifact Template

Write each Task file with this structure, filled as a worked example rather than an empty skeleton:

```markdown
---
id: T014
title: Resume unfinished work without changing project files
objective: O008
status: planned
depends_on: [{task: T013, requires: clear}]
owner_validation: {required: true}
planned_by: {role: planner, session: planning-example}
---

# T014: Resume unfinished work without changing project files

## Outcome

Reopening an existing V2 project shows selected work, recorded Check freshness, and one understandable next action without changing project files.

## User Check

If the owner requests a local Task Check, open a project with a Task awaiting
that Check; invoke resume as the human; confirm Task/outcome, owner-wait
distinction, evidence date, and next action match the board. No automatic
command execution or newly written evidence should appear. If the owner waives
the local Check, record the waiver in Task evidence and rely on the mandatory
Full Objective Check for integration.

## Done When

Correct Task/Objective and next action; stale/unknown evidence is explicit;
malformed/missing selection is named; file bytes and mtimes unchanged; an
optional Task Check is either current or explicitly waived; required owner
validation is recorded after the mandatory integration evidence.

## Context Files

Name exact paths only — no globs, no directory-only entries: `cmd/board.go`, `cmd/board_test.go` (command pattern); `internal/data/project.go`, `internal/data/next.go` (shared index/projection); `main.go`, `main_test.go`; `internal/resume/resume.go`, `internal/resume/resume_test.go`, `cmd/resume.go`, `cmd/resume_test.go`.

## Design References

Design sections 4, 5, 10, and 11.

## Guardrails

FS-01, FS-03, DATA-02, DATA-03, ARCH-01, ARCH-03, TEST-01..04, TEST-08.

## Implementation Plan

1. Confirm the preceding project/Next APIs exist; return REPLAN REQUIRED if they do not.
2. Add thin resume argument handling with named errors.
3. Add a renderer for the resolved projection: Objective/Task, outcome, implementation vs technical vs owner state, Issues, evidence basis/date, and Next.
4. Keep presentation deterministic; no subprocess to generate evidence.
5. Wire main dispatch through the injected runner pattern.
6. Verify awaiting Check, awaiting owner, stale/unknown evidence, missing target, malformed router, no Task yet, and writer failure.

## Boundaries

No new Task states, evidence collection, Objective creation, automatic model routing, or duplicate gate logic.

## Technical Verification

Focused cmd/resume and data tests, no filesystem changes on success/failure, `make build && make test`.

## Technical Evidence

Pending execution: named cases, results, reviewed source basis, files read/changed, and limitations, recorded after the work lands.

## Drift Notes

New module or architecture delta beyond the documented Codebase Map, reconciled through the planner before Check.
```

`title` and `objective` are separate required fields with different jobs. `title` is a short, plain-English phrase written for the task's owner; it must never be the Outcome text, a truncation of it, or a restatement of the technical objective. `objective` is always an `O###` reference to the owning Objective, never free text, and every Task belongs to exactly one Objective. `owner_validation.required` and `planned_by` are recorded at planning time, not left as placeholders.

Title readability itself — whether a generated title actually reads clearly to the owner — is evaluated by agent scenarios in E50; the checks in this repository assert only that `title` and `objective` are distinct required fields and that the no-reuse rule above is stated, not that any given title reads well.

## Design Template

`Design.md` carries these sections: Architecture, Components/Codebase Map, Interfaces and Data Flow, Boundaries, Decisions, Current Technical State. Design describes implemented reality concisely; it is not a backlog. Planned change lives in the active Objective's deltas until the work lands, at which point Design is reconciled to match.

## Guardrails Rule

`Guardrails.md` holds durable project constraints only, once Design is mature enough to know what they are — not task-specific instructions and not a running commentary on in-flight work. Keep roughly 10–20 substantive rules, using stable category IDs (like `FS-01`, `TEST-03`) so tasks and checks can reference a rule instead of restating its prose. Guardrails, not this skill or any task file, owns severity and exception authority for those rules. Tasks and checks must not duplicate guardrail prose; they cite the rule ID.

## Commands And Procedures Reconciliation

When reconciling a migrated project's commands and procedures, reference
`agent-skills/references/commands-and-procedures.md` for the config contract
and the Health-Check.md mapping; do not restate that mapping here.

## Issue Capture

Enter Issue capture as an entry from this workflow when planning surfaces
drift, a guardrail gap, or follow-up that does not belong inside Design or the
current Objective's Tasks. Capture it as an Issue instead of folding it into
Design or a Task plan; see `agent-skills/references/issue-capture.md` for the
artifact template and rules. This skill may read and reference an Issue; it
does not close one — with one exception: when this skill promotes an Issue's
repair into a new Objective, it retires that Issue immediately with
disposition `escalated` naming the new Objective, per "Escalation Retires The
Issue" in `agent-skills/references/issue-capture.md`. That is the only Issue
closure this skill performs.

## Readiness Gate

Detail an Objective's Tasks only when all of the following are settled:

- Interfaces and data ownership for the work are settled, not still open questions.
- Constraints on the work are scoped, not open-ended.
- The outcomes of any Task or Objective this one depends on are known.
- There is a verification approach for the work, even if the detailed check itself runs later.

When any of these is not yet true, that gap is the thing to resolve next — either by settling it directly or by writing a bounded research Task that produces the missing decision.

## Rules

- Write only `Design.md`, `Guardrails.md`, the current Objective, the next Objective's detailed Tasks, and the routing handoff. Do not write production code, and do not detail backlog beyond the next Objective.
- An unknown implementation approach becomes a bounded research Task with a named decision deliverable, not a speculative plan.
- Split a Task that carries multiple unrelated outcomes or an unresolved architectural decision.
- Do not create a Release merely to fill the template, and do not add a Release-owned Task list or a separate release helper document.
- `title` must never be the Outcome text, a truncation of it, or a restatement of the technical objective; it is a short, plain-English phrase for the task's owner.
- Route product choices to the owner; do not settle them by inference. Technical readiness does not require the owner to review code.
- Reference `agent-skills/references/check-method.md` for verification method; do not restate it here.
- Use `state` only for router phase, task `status` only for task lifecycle, and `stage` only when an item is `in_progress`.
