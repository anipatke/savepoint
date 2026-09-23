# Agents Guide

## Workflow

1. Read `.savepoint/router.md` — state + next action
2. Activate skill per table below
3. Read: router → epic → task → source files

The phase skill is the canonical workflow source. This guide defines routing, terminology, and repo rules only; do not duplicate phase-by-phase prompt instructions here.

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
- CI runs the full gate with `make ci`. A Full Objective or Goal Check requires current successful `make test-full` evidence; the optional Task Check does not replace it.
- Reuse a successful full result only for metadata-only corrections. Record the original command, time, toolchain, and result, then prove code, tests, fixtures, dependencies, and gate definitions are unchanged since that run. Any change to those inputs requires a fresh full run.
- A Task Check is optional, not an automatic implementation gate. If the
  owner skips the optional independent Task Check, the Task evidence must
  carry an explicit owner waiver naming the Task, reason, actor, and time.
  That waiver is not technical `CLEAR` and does not waive any acceptance
  criterion, guardrail, Objective Check, or Goal Check.
- The Full Objective Check is mandatory before an Objective can close. It is
  the V2 equivalent of the epic-level integration gate and covers every owned
  Task, including Tasks whose optional Task Check was waived, plus cross-Task
  integration and Design reconciliation.
- A Goal Check is mandatory whenever a Goal exists. It covers all member
  Objectives and cross-Objective integration, followed by exact owner
  acceptance of the current Check.

The runtime gate resolvers and their tests in `internal/data` enforce this
contract (`CheckWaiver` in `evidence_v2.go` and `gate_v2.go`).

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
- Agents may set a Task to `status: in_progress` when starting implementation.
- Only the user may set a Task to `status: done` or retreat a Task to an earlier status.
- Only `savepoint-check` may write a Check record or close an Issue as `verified`. The owner may close an Issue as `accepted` through an explicit decision with reason, actor, and time; an agent may record that exact decision but may not infer it. `savepoint-design` may close an Issue as `escalated` when it promotes the repair into a new Objective.

## Issue Capture

Use Issue capture when planning, implementation, or a Check surfaces a defect, drift, a guardrail gap, or other durable follow-up that does not belong inside Design or the current Objective's Tasks. “Defect” stays a word the user says; it maps to `type: defect` on the Issue record and does not reopen a separate defect workflow.

- Issues live at `.savepoint/issues/I-###-slug.md`.
- See `agent-skills/references/issue-capture.md` for the artifact template, search-before-creating rule, resolution dispositions, and role boundaries.
- The executor reports repair evidence without granting clearance. A checker closes a proven repair as `verified`; the owner may direct an `accepted` closure after visual inspection without claiming technical `CLEAR`; the planner closes a promoted repair as `escalated`. See `agent-skills/references/issue-capture.md`.

## Implementation

Follow the active skill for execution. During `task`, the canonical flow is `savepoint-task` — it owns the read budget, `status: in_progress` + `stage: build` setting, per-criterion evidence, and the handoff decision between an optional Task Check and the mandatory Full Objective Check.

**Stop. Prompt the user before continuing.** Only the user may mark a task `status: done` or retreat a task to an earlier status.

## Check

`savepoint-check` is the only role that can write a Check record or close an
Issue as `verified`. The owner may explicitly close an Issue as `accepted`
without a Check; this does not claim technical `CLEAR` or waive a mandatory
Objective or Goal Check. The owner closes Tasks and accepts Objective/Goal
outcomes after the required evidence exists.

- A Task Check is optional and runs at Quick evidence only when requested; an explicit owner waiver may skip it, but the waiver is not technical `CLEAR`.
- A Full Objective Check is mandatory, runs at Full evidence, and covers every owned Task (including waived Tasks), cross-Task integration, and reconciliation against `Design.md`.
- A Goal Check is mandatory whenever a Goal exists and covers cross-Objective integration before exact owner acceptance.
- The Check session must be independent from the executor's own session — the same model is allowed, the same session is not.
- Both evidence modes apply `agent-skills/references/check-method.md` in full: scope locks, coverage matrices, the adversarial pass, materiality, and re-check convergence.
- Apply `.savepoint/Guardrails.md` when the project has it; its absence is not a finding.
- A CLEAR Check signed by a checker is current on its own. A freshness assessment is optional and only marks the latest Check `stale` or `unknown`.
- Check records are immutable, at `.savepoint/checks/C-###-slug.md`. A recheck writes a new record naming the one it supersedes; it never edits a prior run.

## Code Style

Code style is project-owned policy: the `STYLE` rules in `.savepoint/Guardrails.md` are the single source of truth. Read them when writing or reviewing code. They are Guideline severity and advisory — they inform review but never block on their own. If the project has no `.savepoint/Guardrails.md`, code style is not defined for it.

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
| `main.go` | Wires CLI commands, version output, and embedded V2 templates. Its resume path loads the V2 project and router, resolves `data.Next`, and renders the result. |
| `cmd/` | Parses arguments and dispatches init, board, doctor, upgrade-assets, migrate, and resume commands. It leaves project records, gates, and rendering to `internal/` packages. |
| `internal/init/` | Validates targets and scaffolds `templates/project-v2`. Upgrade-assets checks the project schema through `internal/data` and safely refreshes managed guidance and assets. |
| `internal/board/` | Owns schema-aware board dispatch and rejects filter flags that do not apply to the project. V2 board rendering lives in `internal/board/v2`. |
| `internal/board/v2/` | Implements the V2 TUI and non-TTY board, using the shared `data.Next` projection for the next action. It renders Objective, Task, and Goal navigation and details. |
| `internal/buildtool/` | Runs named Go build and test gates and prepares cross-platform binaries, archives, and checksums. |
| `internal/doctor/` | Runs read-only project diagnostics and configured quality gates. It reports Goal readiness through the canonical `internal/data` resolver and formats repair guidance. |
| `internal/data/` | Loads projects and owns the shared schema-version check, V2 records, indexes, lifecycle and gate decisions, Goal completion, and the `Next` projection. The retained V1 readers are used by `internal/migrate` to parse conversion inputs. |
| `internal/resume/` | Renders a resolved `data.Next` projection and shared evidence wording to plain text. It performs no filesystem, subprocess, network, or TTY access. |
| `internal/migrate/` | Previews and converts V1 project files into V2 files and a source manifest. Apply requires a Git work tree and checks planned paths for modified, untracked, or ignored files before writing directly; a partial failure reports written paths and Git undo commands. |
| `internal/testutil/` | Provides shared Go test fixtures and filesystem helpers for internal packages. |
| `internal/styles/` | Defines the TUI palette and styles. |
| `templates/` | Contains the V2 scaffold's Markdown, YAML, prompts, and workflow guidance. |
| `agent-skills/` | Contains the four active V2 skills and their shared references. |

After schema-2 migration, ordinary board, resume, doctor, init, and
upgrade-assets behavior uses the V2 runtime. The V1 readers retained in
`internal/data` are used by `internal/migrate` to parse legacy source projects.

## Context Budget

- **Read only what you need.** Each phase has a strict read budget. Do not read files outside your current phase's context.
- **No exploratory reads.** Read only the files listed in the task's `## Context Files`. Do not glob or search for new information unless explicitly instructed.
- **Token awareness.** Every file read consumes context window. Before reading a file, ask: "Do I need this to complete my current phase?"

## CLI Rules

**Never run `savepoint` commands.** The CLI is for the human. Edit files directly.

## Reporting to the Owner

The Context Log stays technical and precise — it's the record a Check session verifies later. Chat replies to the owner are a different audience: a few plain sentences, no jargon, no file/function dumps unless asked. Say what happened and what's next; leave the mechanism in the Context Log.

## V2 Routing

This section records the V2 routing contract for a project whose `config.yml` declares `schema_version: 2`. It is active in this migrated repository; in a legacy V1 scaffold it is not active until migration, while this repository's active table is the four-state table above.

| Router `state` | Skill |
|-----------------|-------|
| idea | savepoint-idea |
| design | savepoint-design |
| task | savepoint-task |
| check | savepoint-check |

`REPLAN REQUIRED`, returned by an executor that hits a materially invalid plan, routes back into `savepoint-design`. It is not a fifth router state — the state stays `design` while the planner resolves what broke.

Three shared references back these four skills: `agent-skills/references/check-method.md`, `agent-skills/references/issue-capture.md`, and `agent-skills/references/commands-and-procedures.md`. Each carries `triggerable: false` frontmatter and is non-triggerable on its own — it is loaded in full by the skill that owns it (`savepoint-check` loads `check-method.md`; `savepoint-design`, `savepoint-task`, and `savepoint-check` each enter `issue-capture.md` from their own workflow; `savepoint-design` loads `commands-and-procedures.md` for config reconciliation), not invoked directly.

`E47` ships this table as the scaffold default for new V2 projects; E50 activates it here after migration. The V1 scaffold and its skills are gone (O-021); they remain only as byte-preserved history.

## Legacy V1 compatibility (not active)

The following contract is retained only for reading archived V1 projects and this repository's own historical records; it is not an active route in this schema-2 repository, and the V1 scaffold that once shipped it is gone (O-021).

| task-building | savepoint-build-task |
| audit-pending | savepoint-audit-epic |

An explicit request uses `savepoint-audit-task` while `state` stays `task-building`; that is not a router state and not a new state here.

Task `stage` (build/test/audit): **required** when `status: in_progress` — Task lifecycle rules are owned by `internal/data`; legacy `phase` is parse compatibility only and must not be used in new task guidance. Only the user may set a task to `status: done`.
