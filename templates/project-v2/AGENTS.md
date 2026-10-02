# Agents Guide

## Workflow

1. If the owner pasted a `Next` line, act on it directly; it is the selection, so do not re-run `savepoint resume` to confirm it. Otherwise run the read-only `savepoint resume` command and act on its `Next` line.
2. The Next line's first word says what to do; use it to choose the skill: `Start`, `Build`, or `Test` → `savepoint-task`; `Check` → `savepoint-check`; `Plan` or `Replan` → `savepoint-design`; `Pick a Task in` → select the Objective's next Task, then `savepoint-task`; `Fix` → repair with `savepoint-task` under `issue-capture.md`; `Choose` → report it to the owner, who decides; `Accept`, `Close`, `Blocked`, `Done`, or `Resolved` → report it to the owner, who decides. When an Objective or Task is selected with an Issue, follow the Objective or Task and treat the Issue as context. If a selected record is reported stale, do not substitute other work.
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

After a direct Issue repair selected alone, follow `agent-skills/references/issue-capture.md`, Out-Of-Scope Repair, for the exact router handoff. Preserve `release:`.

Three shared references back these four skills and are never triggered directly: `agent-skills/references/check-method.md` (loaded in full by `savepoint-check`), `agent-skills/references/issue-capture.md` (entered by `savepoint-design`, `savepoint-task`, and `savepoint-check` from their own workflow), and `agent-skills/references/commands-and-procedures.md` (loaded by `savepoint-design` for config reconciliation). Each carries `triggerable: false` frontmatter.

Read `.savepoint/Idea.md` only for original intent, `.savepoint/Design.md` only for architecture readiness.

## Verification Policy

- Every Task records per-criterion evidence and runs its configured gate before handoff.
- Gate commands are project-owned: `quality_gates` in `.savepoint/config.yml` (build, lint, typecheck, test), plus any fuller gate listed under Build below. When no gate is configured, record that none ran; never invent one.
- Focused test runs are for iteration. Ordinary Task handoff runs the configured build and test gates; migration or platform-sensitive Task handoff runs the project's full gate fresh when it defines one separately.
- A Full Objective Check requires current successful full-gate evidence; the optional Task Check does not replace it.
- Reuse a successful full result only for metadata-only corrections. Record the original command, time, toolchain, and result, then prove code, tests, fixtures, dependencies, and gate definitions are unchanged since that run. Any change to those inputs requires a fresh full run.
- A Task Check is optional, not an automatic implementation gate. If the
  owner skips the optional independent Task Check, the Task evidence must
  carry an explicit owner waiver naming the Task, reason, actor, and time.
  That waiver is not technical `CLEAR` and does not waive any acceptance
  criterion, guardrail, or Objective Check.
- A recorded owner Task-check waiver does satisfy a `requires: clear` Task
  dependency: the owner's completion decision stands in for clear there. It
  never satisfies `requires: accepted`, since there is no Check to accept.
- Do not decide whether a dependency blocks by reading this prose. The
  Savepoint runtime's dependency gate decides, and
  `savepoint resume` and the board's transition gate report its result as
  `Blocked:` lines. If neither reports a block, the dependency is met.
- The Full Objective Check is mandatory before an Objective can close. It is
  the V2 higher-level integration gate and covers every owned Task, including
  Tasks whose optional Task Check was waived, plus cross-Task integration and
  Design reconciliation.

A Check session must be independent from the executor's conversation: the same model is allowed, the same session is not. Only `savepoint-check` writes Check records. Each run is immutable at `.savepoint/checks/C-###-slug.md`; a recheck writes a new record naming the one it supersedes.

During a Check, load `agent-skills/references/check-method.md` in full and apply the selected mode: Quick follows its Quick Check Procedure; Full applies every section, including coverage matrices and the adversarial pass. Apply Guardrails when present; absence is not a finding. A signed `CLEAR` Check is current on its own; optional freshness assessment marks only `stale` or `unknown`.

## Required Goal Context

Every Savepoint project must have a live Goal selected by the router, and
every live Objective must name exactly one Goal in its `release:` field.
`savepoint init` creates and selects G-001, titled after the project.
`savepoint migrate` keeps the V1 router's live Goal. When that selection is
missing or unresolvable, it selects a uniquely identifiable existing live
Goal for the converted active work when possible. If selected work belongs
only to a historical Goal, migration creates a live continuation with a G-###
identity and moves that active Objective into it. Converted V1 Releases keep
their R-### identities. An unresolved release lifecycle decision stays in
the preview and blocks Apply.

If the router Goal is missing, blank, or `none`, Next says `Choose a Goal`; use
`g` to select a live Goal. Unknown or archived router selections are reported
as selection diagnostics; use `g` to select a live Goal. If there are no live
Goals, `savepoint doctor` says to create one. An Objective missing `release:` remains
loadable, but resume and the board flag it and doctor names the Objective and
the exact `release:` field using an `R-###` or `G-###` Goal ID. Unknown or
malformed Objective references remain errors. The board shows only the selected
Goal's Objectives and Tasks; it never falls back to a project-wide view.

Existing Goals keep their stable `R-###` identities, paths, and references.
New Goals use stable `G-###` identities under `.savepoint/releases/`
(`Release.md`); Objective and router `release:` fields remain the persisted
compatibility boundary. Converted V1 Releases keep their `R-###` identities.
Goals group Objectives; a Goal is complete when every member Objective is
complete. Goals do not own Tasks, and do not publish, deploy, tag, or generate
changelogs.

## Terminology

- Router `state` is the current state: `idea`, `design`, `task`, or `check`.
- Task `status`: only `planned`, `in_progress`, or `done`.
- Task `stage`: **required** when `status: in_progress` — `build` → `test` → `audit`; reaching `audit` means the Task is ready for a Check, and explicitly does not mean it passed.
- Never write `stage: implementation`; use `stage: build` when starting implementation work.
- Agents may set a Task to `status: in_progress` when starting implementation, and its owning Objective from `planned` to `in_progress` at the same time. That is the only Objective status change an agent makes.
- Only the user may set a Task to `status: done` or retreat a Task to an earlier status. **Stop. Prompt the user before continuing** when that decision is required.
- Only `savepoint-check` may write a Check record or close an Issue as `verified`. The owner may resolve an Issue as `accepted` from the board's Issues panel and reopen a resolved one. Board resolution records the fixed reason, owner actor, and time; it is not technical `CLEAR`. An agent may record an owner decision only when directly instructed. `savepoint-design` may close an Issue as `escalated` when it promotes the repair into a new Objective.

## Issue Capture

Capture defects, drift, guardrail gaps, and durable follow-up that does not belong inside Design or the current Objective's Tasks as Issues at `.savepoint/issues/I-###-slug.md`. "Defect" maps to `type: defect`, not a separate workflow. Enter `agent-skills/references/issue-capture.md` for search-before-creating, artifact fields, resolution authority, append-only history, and repair routing.

## Worktree Lanes

The owner may run independent Tasks or Issue repairs side by side in git worktrees, starting each lane by pasting its `Next` line. A lane is a worktree the owner names as one, or any checkout where `git rev-parse --git-dir` differs from `git rev-parse --git-common-dir`. In a lane:

- Do not edit `.savepoint/router.md`; skip every router write the active skill would make. The owner sets the router on the main branch after merging.
- Do not create Tasks, Checks, or Issues. Their IDs are allocated per checkout and would collide at merge. Record a needed one as a note in the Task or Issue evidence for the owner.
- Commit on the lane branch; do not push or merge. Optional Task Checks and the Full Objective Check run on the main branch after the lane merges.

`savepoint-design` shapes Tasks for lanes where practical; this section only sets what an agent may write inside one.

## Existing Codebase Adoption

Design is reconstructed from the code through targeted reads needed by the current Objective or Task: no whole-repository scan, automatic analysis pass, or call to a model service. What exists goes to `.savepoint/Design.md`: structure to Components/Codebase Map, verified behaviour to Current Technical State. Unread areas are recorded as unknown. What the code is for goes to `.savepoint/Idea.md`: owner intent is recorded in `.savepoint/Idea.md` through `savepoint-idea`, never inferred from source.

Adoption preserves user-authored files. `savepoint init` may add or refresh the Savepoint-managed block in an existing agent guide, preserving every byte outside that block; other Savepoint files are added under `.savepoint/`. Concept, Health-Check, and a procedures file are optional. Their absence is normal, not a finding. A Goal is required: `savepoint init` supplies G-001; migration keeps or selects an existing live Goal, or creates a continuation for historical work. Follow Required Goal Context for `Choose a Goal` and `savepoint doctor` repair guidance.

## Code Style

Code style is project-owned policy: the `STYLE` rules in `.savepoint/Guardrails.md` are the single source of truth. Read them when writing or reviewing code. They are Guideline severity and advisory — they inform review but never block a Check by themselves. If the project has no `.savepoint/Guardrails.md`, code style is not defined for it.

## Build

List the project's gate commands here when they differ from, or add to, `quality_gates` in `.savepoint/config.yml` — for example a slower full gate for migration, platform-sensitive work, and Full Objective Checks.

See Verification Policy for gate evidence and reuse requirements.

## Codebase Map

| Module | Purpose |
|--------|---------|

## Context Budget

Follow the active skill's Read section and the Task's `## Context Files`. Do not scan or search outside that skill's scope. A necessary targeted extra read during Task execution is allowed only when logged with the file and reason as `savepoint-task` requires. Read Idea only for original intent and Design only for architecture readiness.

## CLI Rules

Agents may run `savepoint resume`, a read-only command that prints `Next` without writing project files. No other `savepoint` command is for agents except the narrow Task creation operation below. `savepoint health setup [dir] [--apply]` is human-only: it suggests health tools and, with `--apply`, saves them; agents never run it.

Exception: agents may run `savepoint health check O-### [dir]` only during a Full Objective Check, after the full gate, to collect one official Code Health snapshot and record its ID in the Check. Task Checks and all other activity never run it.

Exception: agents may run `savepoint create-task --objective O-### --draft <path> [dir]` only to create a new Task from an ID-free draft. The command assigns the project-wide Task ID and strict-loads the V2 index before reporting success. Do not use it to edit or rename a Task, and do not choose or write Task IDs manually. After creating or renaming any other identity-bearing V2 record, run `savepoint resume` to require strict loading of the full V2 index.

## Reporting to the Owner

The Context Log stays technical and precise — it's the record a Check session verifies later. Chat replies to the owner are a different audience: a few plain sentences, no jargon, no file/function dumps unless asked. Say what happened and what's next; leave the mechanism in the Context Log.
