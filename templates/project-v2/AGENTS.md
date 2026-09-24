# Agents Guide

## Workflow

1. Run the read-only `savepoint resume` command and act on its `Next` line. A `Next` line pasted by the owner is also an explicit selection.
2. The Next line's first word says what to do; use it to choose the skill: `Start`, `Build`, or `Test` → `savepoint-task`; `Check` → `savepoint-check`; `Plan` or `Replan` → `savepoint-design`; `Pick a Task in` → select the Objective's next Task, then `savepoint-task`; `Fix` → repair with `savepoint-task` under `issue-capture.md`; `Accept`, `Close`, `Blocked`, `Done`, or `Resolved` → report it to the owner, who decides. When an Objective or Task is selected with an Issue, follow the Objective or Task and treat the Issue as context. If a selected record is reported stale, do not substitute other work.
3. Activate the skill per the table below and follow its Read section and the active Task's Context Files.

If the Next line does not name an Objective, Task, or Issue, follow the router `state` and the resume guidance; do not guess a record. If the `savepoint` binary is unavailable, read the router selection, report that the tool is missing, and do not guess the next step.

The active skill is the canonical workflow source. This guide defines routing, terminology, and repo rules only; do not duplicate skill-by-skill prompt instructions here.

## Skill Activation

| State | Skill |
|-------|-------|
| idea | savepoint-idea |
| design | savepoint-design |
| task | savepoint-task |
| check | savepoint-check |

`REPLAN REQUIRED`, returned by `savepoint-task` on a materially invalid plan, routes back into `savepoint-design`. It is not a fifth router state — `state` stays `design` while the planner resolves what broke.

Use the `skill` tool when the listed skill is available. If the agent says the skill is not found, read `agent-skills/{skill}/SKILL.md` directly and follow it as the active skill.

## Router Selection

The board advances the router when the owner closes a Task. `savepoint-design` selects the next Objective, and `savepoint-task` selects the Task it starts. Agents and board actions never blank or change `release:`; only an explicit owner Goal choice changes it.

An owner may ask, `set router to O-### [T-###] [I-###]`, or select an Issue alone with `set router to I-###`. Confirm every named record exists and, when a Task is named, require an Objective and confirm the Task belongs to it. If a record is missing or the Task belongs elsewhere, do not edit the router; explain why. For an Issue-only selection, clear `objective` and `task`; for an Objective/Task selection, clear `issue` unless the owner named one. Edit only the `objective`, `task`, and `issue` selection keys, leaving `release:` byte-for-byte unchanged unless the owner names a Release as part of an explicit Goal choice. Finish by running `savepoint resume` and showing its `Next` line.

After a direct Issue repair selected alone, `savepoint-task` records `repair_attempted`, then moves the router to the Issue's Objective when its linked Tasks and Objective-scoped Checks identify exactly one. It clears `task` and `issue`; if there is no single Objective, it clears `issue` only. It leaves `release:` unchanged and runs `savepoint resume` to show the next step.

Three shared references back these four skills and are never triggered directly: `agent-skills/references/check-method.md` (loaded in full by `savepoint-check`), `agent-skills/references/issue-capture.md` (entered by `savepoint-design`, `savepoint-task`, and `savepoint-check` from their own workflow), and `agent-skills/references/commands-and-procedures.md` (loaded by `savepoint-design` for config reconciliation). Each carries `triggerable: false` frontmatter.

Read `.savepoint/Idea.md` only for original intent, `.savepoint/Design.md` only for architecture readiness.

## Verification Policy

- Every Task records per-criterion evidence and runs its configured gate before handoff.
- Focused `make test-focused TEST=...` runs are for iteration. Ordinary Task handoff uses `make build && make test-fast`; migration or platform-sensitive Task handoff uses a fresh `make test-full`.
- CI runs the full gate with `make ci`. A Full Objective or Goal Check requires current successful `make test-full` evidence; the optional Task Check does not replace it.
- Reuse a successful full result only for metadata-only corrections. Record the original command, time, toolchain, and result, then prove code, tests, fixtures, dependencies, and gate definitions are unchanged since that run. Any change to those inputs requires a fresh full run.
- A Task Check is optional, not an automatic implementation gate. If the
  owner skips the optional independent Task Check, the Task evidence must
  carry an explicit owner waiver naming the Task, reason, actor, and time.
  That waiver is not technical `CLEAR` and does not waive any acceptance
  criterion, guardrail, Objective Check, or Goal Check.
- A recorded owner Task-check waiver does satisfy a `requires: clear` Task
  dependency: the owner's completion decision stands in for clear there. It
  never satisfies `requires: accepted`, since there is no Check to accept.
- Do not decide whether a dependency blocks by reading this prose. The
  runtime gate (`ResolveTaskDependencyV2` in `internal/data`) decides, and
  `savepoint resume` and the board's transition gate report its result as
  `Blocked:` lines. If neither reports a block, the dependency is met.
- The Full Objective Check is mandatory before an Objective can close. It is
  the V2 higher-level integration gate and covers every owned Task, including
  Tasks whose optional Task Check was waived, plus cross-Task integration and
  Design reconciliation.
- A Goal Check is mandatory whenever a Goal exists. It covers all member
  Objectives and cross-Objective integration, followed by exact owner
  acceptance of the current Check.

## Optional Goals

A Goal is an optional delivery context for related Objectives, not another
workflow state or a required step. Without one, the ordinary Idea → Design →
Task → Check path remains complete. In the V2 board, `g` opens the Goal
selector; `r` remains an undisplayed compatibility alias.

Existing storage is unchanged: Goals use stable `R-###` records under
`.savepoint/releases/` (`Release.md`), Objective `release:` references, router
`release:` selections, and Check `scope.kind: release`. These names are a
compatibility boundary, not the public board vocabulary. Goals group
Objectives, do not own Tasks, and do not publish, deploy, tag, or generate
changelogs. A Goal Check retains the existing cross-Objective integration and
exact-owner-acceptance requirements.

## Terminology

- Router `state` is the current state: `idea`, `design`, `task`, or `check`.
- Task `status`: only `planned`, `in_progress`, or `done`.
- Task `stage`: **required** when `status: in_progress` — `build` → `test` → `audit`; reaching `audit` means the Task is ready for a Check, and explicitly does not mean it passed.
- Never write `stage: implementation`; use `stage: build` when starting implementation work.
- Agents may set a Task to `status: in_progress` when starting implementation, and its owning Objective from `planned` to `in_progress` at the same time. That is the only Objective status change an agent makes.
- Only the user may set a Task to `status: done` or retreat a Task to an earlier status.
- Only `savepoint-check` may write a Check record or close an Issue as `verified`. The owner may close an Issue as `accepted` through an explicit decision with reason, actor, and time; an agent may record that exact decision but may not infer it. `savepoint-design` may close an Issue as `escalated` when it promotes the repair into a new Objective.

## Issue Capture

Use Issue capture when planning, implementation, or a Check surfaces a defect, drift, a guardrail gap, or other durable follow-up that does not belong inside Design or the current Objective's Tasks. "Defect" stays a word the user says; it maps to `type: defect` on the Issue record and does not reopen a separate defect workflow.

- Issues live at `.savepoint/issues/I-###-slug.md`.
- See `agent-skills/references/issue-capture.md` for the artifact template, the search-before-creating rule, resolution dispositions, and role boundaries.
- The executor reports repair evidence on an Issue without closing it; only `savepoint-check` verifies the proof and closes it.

## Implementation

Follow the active skill for execution. During `task`, the canonical flow is `savepoint-task` — it owns the read budget (a Task's own `## Context Files`), `status: in_progress` + `stage: build` setting, per-criterion evidence, and the handoff decision between an optional Task Check and the mandatory Full Objective Check.

**Stop. Prompt the user before continuing.** Only the user may mark a Task `status: done` or retreat a Task to an earlier status. An explicit owner decision may close an Issue as `accepted` without a Check; that decision does not waive a mandatory Objective or Goal Check.

## Check

`savepoint-check` is the only role that can write a Check record or close an
Issue as `verified`. The owner may explicitly close an Issue as `accepted`
without a Check; this does not claim technical `CLEAR` or waive a mandatory
Objective or Goal Check. The owner closes Tasks and accepts Objective/Goal
outcomes after the required evidence exists.

- A Task Check is optional and runs at Quick evidence only when requested; an explicit owner waiver may skip it, but the waiver is not technical `CLEAR` (it still satisfies a `requires: clear` dependency; see Verification Policy).
- A Full Objective Check is mandatory, runs at Full evidence, and covers every owned Task (including waived Tasks), cross-Task integration, and reconciliation against `Design.md`.
- A Goal Check is mandatory whenever a Goal exists and covers cross-Objective integration before exact owner acceptance.
- The Check session must be independent from the executor's own session — the same model is allowed, the same session is not.
- Both evidence modes apply `agent-skills/references/check-method.md` in full: scope locks, coverage matrices, the adversarial pass, materiality, and re-check convergence.
- Apply `.savepoint/Guardrails.md` when the project has it; its absence is not a finding.
- A CLEAR Check signed by a checker is current on its own. A freshness assessment is optional and only marks the latest Check `stale` or `unknown`.
- Check records are immutable, at `.savepoint/checks/C-###-slug.md`. A recheck writes a new record naming the one it supersedes; it never edits a prior run.

## Existing Codebase Adoption

`savepoint init` into a directory that already holds a codebase still hands the agent an empty `Design.md` describing a system that is already sitting in the repository. This section is what that agent reads before filling it in.

Design is reconstructed from the code through targeted reads — the same read-only discipline every other role in this guide follows — never through a whole-repository scan, an automatic analysis pass, or a call to a model service. Read the files a current Objective or Task actually needs, not everything the repository contains. Intent comes from the owner: why the system exists and who it is for is not something the code can state, so it is asked of the owner and recorded in `.savepoint/Idea.md` through `savepoint-idea`, never inferred from source.

What exists goes to `.savepoint/Design.md`: concrete structure to Components/Codebase Map, what was actually verified to `Current Technical State`. What the code is for goes to `.savepoint/Idea.md`, through the owner. An area not yet read is recorded as unknown; it is never filled in by inference.

Adoption does not rewrite user-authored files. `savepoint init` may add or refresh the Savepoint-managed block in an existing agent guide, preserving every byte outside that block; all other Savepoint files are added under `.savepoint/`.

This guidance degrades when optional files are absent: a V2 project may have no Goal (stored as an `R-###` Release record), no Concept, no Health-Check, or no procedures file, and none is required before adoption can proceed. Their absence is normal, not a finding.

## Code Style

Code style is project-owned policy: the `STYLE` rules in `.savepoint/Guardrails.md` are the single source of truth. Read them when writing or reviewing code. They are Guideline severity and advisory — they inform review but never block a Check by themselves. If the project has no `.savepoint/Guardrails.md`, code style is not defined for it.

## Build

```bash
make build && make test-fast   # ordinary Task handoff
make test-full                 # migration/platform-sensitive Task or Full Objective/Goal Check
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

Agents may run exactly `savepoint resume`, a read-only command that prints `Next` without writing project files. Every other `savepoint` command is for the human.

## Reporting to the Owner

The Context Log stays technical and precise — it's the record a Check session verifies later. Chat replies to the owner are a different audience: a few plain sentences, no jargon, no file/function dumps unless asked. Say what happened and what's next; leave the mechanism in the Context Log.
