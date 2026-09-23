# Agents Guide

## Workflow

1. Read `.savepoint/router.md` — state + next action
2. Activate skill per table below
3. Read: router → Idea/Design/Objective → active Task → its Context Files

The active skill is the canonical workflow source for its state. This guide defines routing, terminology, and repo rules only; do not duplicate skill-by-skill prompt instructions here.

## Skill Activation

| State | Skill |
|-------|-------|
| idea | savepoint-idea |
| design | savepoint-design |
| task | savepoint-task |
| check | savepoint-check |

`REPLAN REQUIRED`, returned by `savepoint-task` on a materially invalid plan, routes back into `savepoint-design`. It is not a fifth router state — `state` stays `design` while the planner resolves what broke.

Use the `skill` tool when the listed skill is available. If the agent says the skill is not found, read `agent-skills/{skill}/SKILL.md` directly and follow it as the active skill.

Three shared references back these four skills and are never triggered directly: `agent-skills/references/check-method.md` (loaded in full by `savepoint-check`), `agent-skills/references/issue-capture.md` (entered by `savepoint-design`, `savepoint-task`, and `savepoint-check` from their own workflow), and `agent-skills/references/commands-and-procedures.md` (loaded by `savepoint-design` for config reconciliation). Each carries `triggerable: false` frontmatter.

Read `.savepoint/Idea.md` only for original intent, `.savepoint/Design.md` only for architecture readiness.

## Verification Policy

- Every Task records per-criterion evidence and runs its configured gate before handoff.
- Focused `make test-focused TEST=...` runs are for iteration. Ordinary Task handoff uses `make build && make test-fast`; migration or platform-sensitive Task handoff uses a fresh `make test-full`.
- CI runs the full gate with `make ci`. A Full Objective or Release Check requires current successful `make test-full` evidence; the optional Task Check does not replace it.
- Reuse a successful full result only for metadata-only corrections. Record the original command, time, toolchain, and result, then prove code, tests, fixtures, dependencies, and gate definitions are unchanged since that run. Any change to those inputs requires a fresh full run.
- A Task Check is optional, not an automatic implementation gate. If the
  owner skips the optional independent Task Check, the Task evidence must
  carry an explicit owner waiver naming the Task, reason, actor, and time.
  That waiver is not technical `CLEAR` and does not waive any acceptance
  criterion, guardrail, Objective Check, or Release Check.
- The Full Objective Check is mandatory before an Objective can close. It is
  the V2 higher-level integration gate and covers every owned Task, including
  Tasks whose optional Task Check was waived, plus cross-Task integration and
  Design reconciliation.
- A Release Check is mandatory whenever a Release exists. It covers all member
  Objectives and cross-Objective integration, followed by exact owner
  acceptance of the current Check.

## Optional Releases

A Release is an optional delivery boundary, not another public state. Ask for one during Idea only when the owner needs a navigable delivery/package promise across multiple Objectives. Without one, the normal Idea → Design → Task → Check path remains complete and must not report a missing Release record.

When an owner opts in, Design gives the Release a stable `R-###` identity and an outcome, success conditions, and Objective links. Objectives carry the single optional `release: R-###` reference; Releases do not nest files, own Tasks, or maintain a second membership list. Release `done` means current CLEAR integration evidence plus the owner's acceptance of that exact Check, not published or deployed.

## Terminology

- Router `state` is the current state: `idea`, `design`, `task`, or `check`.
- Task `status`: only `planned`, `in_progress`, or `done`.
- Task `stage`: **required** when `status: in_progress` — `build` → `test` → `audit`; reaching `audit` means the Task is ready for a Check, and explicitly does not mean it passed.
- Never write `stage: implementation`; use `stage: build` when starting implementation work.
- Agents may set a Task to `status: in_progress` when starting implementation.
- Only the user may set a Task to `status: done` or retreat a Task to an earlier status.
- Only `savepoint-check` may write a Check record or close an Issue as `verified`. The owner may close an Issue as `accepted` through an explicit decision with reason, actor, and time; an agent may record that exact decision but may not infer it. `savepoint-design` may close an Issue as `escalated` when it promotes the repair into a new Objective.

## Issue Capture

Use Issue capture when planning, implementation, or a Check surfaces a defect, drift, a guardrail gap, or other durable follow-up that does not belong inside Design or the current Objective's Tasks. "Defect" stays a word the user says; it maps to `type: defect` on the Issue record and does not reopen a separate defect workflow.

- Issues live at `.savepoint/issues/I-###-slug.md`.
- See `agent-skills/references/issue-capture.md` for the artifact template, the search-before-creating rule, resolution dispositions, and role boundaries.
- The executor reports repair evidence on an Issue without closing it; only `savepoint-check` verifies the proof and closes it.

## Implementation

Follow the active skill for execution. During `task`, the canonical flow is `savepoint-task` — it owns the read budget (a Task's own `## Context Files`), `status: in_progress` + `stage: build` setting, per-criterion evidence, and the handoff decision between an optional Task Check and the mandatory Full Objective Check.

**Stop. Prompt the user before continuing.** Only the user may mark a Task `status: done` or retreat a Task to an earlier status. An explicit owner decision may close an Issue as `accepted` without a Check; that decision does not waive a mandatory Objective or Release Check.

## Check

`savepoint-check` is the only role that can write a Check record or close an
Issue as `verified`. The owner may explicitly close an Issue as `accepted`
without a Check; this does not claim technical `CLEAR` or waive a mandatory
Objective or Release Check. The owner closes Tasks and accepts
Objective/Release outcomes after the required evidence exists.

- A Task Check is optional and runs at Quick evidence only when requested; an explicit owner waiver may skip it, but the waiver is not technical `CLEAR`.
- A Full Objective Check is mandatory, runs at Full evidence, and covers every owned Task (including waived Tasks), cross-Task integration, and reconciliation against `Design.md`.
- A Release Check is mandatory whenever a Release exists and covers cross-Objective integration before exact owner acceptance.
- The Check session must be independent from the executor's own session — the same model is allowed, the same session is not.
- Both evidence modes apply `agent-skills/references/check-method.md` in full: scope locks, coverage matrices, the adversarial pass, materiality, and re-check convergence.
- Apply `.savepoint/Guardrails.md` when the project has it; its absence is not a finding.
- Check records are immutable, at `.savepoint/checks/C-###-slug.md`. A recheck writes a new record naming the one it supersedes; it never edits a prior run.

## Existing Codebase Adoption

`savepoint init` into a directory that already holds a codebase still hands the agent an empty `Design.md` describing a system that is already sitting in the repository. This section is what that agent reads before filling it in.

Design is reconstructed from the code through targeted reads — the same read-only discipline every other role in this guide follows — never through a whole-repository scan, an automatic analysis pass, or a call to a model service. Read the files a current Objective or Task actually needs, not everything the repository contains. Intent comes from the owner: why the system exists and who it is for is not something the code can state, so it is asked of the owner and recorded in `.savepoint/Idea.md` through `savepoint-idea`, never inferred from source.

What exists goes to `.savepoint/Design.md`: concrete structure to Components/Codebase Map, what was actually verified to `Current Technical State`. What the code is for goes to `.savepoint/Idea.md`, through the owner. An area not yet read is recorded as unknown; it is never filled in by inference.

Adoption does not rewrite user-authored files. `savepoint init` may add or refresh the Savepoint-managed block in an existing agent guide, preserving every byte outside that block; all other Savepoint files are added under `.savepoint/`.

This guidance degrades when optional files are absent: a V2 project ships no Concept, no Health-Check, no procedures file, and no Release record, and none of those are required before adoption can proceed. Their absence is normal, not a finding.

## Code Style

Code style is project-owned policy: the `STYLE` rules in `.savepoint/Guardrails.md` are the single source of truth. Read them when writing or reviewing code. They are Guideline severity and advisory — they inform review but never block a Check by themselves. If the project has no `.savepoint/Guardrails.md`, code style is not defined for it.

## Build

```bash
make build && make test-fast   # ordinary Task handoff
make test-full                 # migration/platform-sensitive Task or Full Objective/Release Check
make ci                        # CI full gate plus distribution and package checks
```

`make test-focused TEST=...` is an iteration aid. Reuse a prior full result only for a metadata-only correction after recording the original run and proving code, tests, fixtures, dependencies, and gate definitions unchanged.

## Codebase Map

| Module | Purpose |
|--------|---------|

## Context Budget

- **Read only what you need.** A Task's `## Context Files` is the read budget, not a suggestion.
- **No exploratory reads.** A read beyond a Task's Context Files is an extra read: allowed, but logged in the Task's evidence with what was read and why.
- **Token awareness.** Every file read consumes context window.

## CLI Rules

**Never run `savepoint` commands.** The CLI is for the human. Edit files directly.

## Reporting to the Owner

The Context Log stays technical and precise — it's the record a Check session verifies later. Chat replies to the owner are a different audience: a few plain sentences, no jargon, no file/function dumps unless asked. Say what happened and what's next; leave the mechanism in the Context Log.
