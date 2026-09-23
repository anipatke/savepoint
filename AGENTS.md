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
- CI runs the full gate with `make ci`. A Full Objective or Release Check requires current successful `make test-full` evidence; the optional Task Check does not replace it.
- Reuse a successful full result only for metadata-only corrections. Record the original command, time, toolchain, and result, then prove code, tests, fixtures, dependencies, and gate definitions are unchanged since that run. Any change to those inputs requires a fresh full run.
- A Task Check is optional, not an automatic implementation gate. If the
  owner skips the optional independent Task Check, the Task evidence must
  carry an explicit owner waiver naming the Task, reason, actor, and time.
  That waiver is not technical `CLEAR` and does not waive any acceptance
  criterion, guardrail, Objective Check, or Release Check.
- The Full Objective Check is mandatory before an Objective can close. It is
  the V2 equivalent of the epic-level integration gate and covers every owned
  Task, including Tasks whose optional Task Check was waived, plus cross-Task
  integration and Design reconciliation.
- A Release Check is mandatory whenever a Release exists. It covers all member
  Objectives and cross-Objective integration, followed by exact owner
  acceptance of the current Check.

The runtime gate resolvers and their tests in `internal/data` enforce this
contract (`CheckWaiver` in `evidence_v2.go` and `gate_v2.go`).

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
Objective or Release Check. The owner closes Tasks and accepts
Objective/Release outcomes after the required evidence exists.

- A Task Check is optional and runs at Quick evidence only when requested; an explicit owner waiver may skip it, but the waiver is not technical `CLEAR`.
- A Full Objective Check is mandatory, runs at Full evidence, and covers every owned Task (including waived Tasks), cross-Task integration, and reconciliation against `Design.md`.
- A Release Check is mandatory whenever a Release exists and covers cross-Objective integration before exact owner acceptance.
- The Check session must be independent from the executor's own session — the same model is allowed, the same session is not.
- Both evidence modes apply `agent-skills/references/check-method.md` in full: scope locks, coverage matrices, the adversarial pass, materiality, and re-check convergence.
- Apply `.savepoint/Guardrails.md` when the project has it; its absence is not a finding.
- Check records are immutable, at `.savepoint/checks/C-###-slug.md`. A recheck writes a new record naming the one it supersedes; it never edits a prior run.

## Code Style

Code style is project-owned policy: the `STYLE` rules in `.savepoint/Guardrails.md` are the single source of truth. Read them when writing or reviewing code. They are Guideline severity and advisory — they inform review but never block on their own. If the project has no `.savepoint/Guardrails.md`, code style is not defined for it.

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
| `main.go` | CLI entrypoint, --version, embedded template wiring for init and upgrade-assets, migrate dispatch, and resume dispatch: resolves the target directory and any pending migration state, loads the V2 project and router, resolves the `Next` projection, and renders it (`runResume`) |
| `cmd/` | CLI command arg parsing and dispatch for init, board, doctor, upgrade-assets, migrate, and resume — argument parsing only, no record parsing, gate reading, or rendering (ARCH-01). It parses `board`'s V1 filters (`--release`, `--epic`) and V2 filter (`--objective`) without judging which applies; that is refused behind the dispatch, where the schema is known |
| `internal/init/` | Target validation, scaffold writing from a caller-selected template tree (`init` defaults to `templates/project-v2`), upgrade-assets schema-version dispatch between `templates/project` and `templates/project-v2` via `data.ReadSchemaVersion`, upgrade provenance manifest, managed AGENTS.md merge/conflict behavior, and safe project asset refresh |
| `internal/board/` | The board's schema dispatch — resolve the project root once, `data.LoadProject`, then run the V1 board or `internal/board/v2`, refusing a filter flag the resolved schema has no meaning for (CFG-01) — plus the whole V1 board: TUI board, overlays, epic sidebar, Next Activity line, router priority key, detail checklist rendering, status glyphs, forced color profile, debug logging hooks, async update I/O commands, defect summary/overlay/detail rendering, related-defect card markers, audit register overlay with finding detail and linked-finding backlinks, shared board utilities |
| `internal/board/v2/` | The V2 board, a sibling package no V1 board type or V1 record type is reachable from: model state, the single load command (`data.LoadProject`, `ReadStateV2`, `migrate.PendingOperation`, `data.ResolveNext`) startup and every reload share, the load-diagnostic screen for a project the V2 index refuses, the Objective sidebar — list, cursor, selection, per-Objective status/clearance/integration/dependency state, and Task filtering by ownership read from `index.ObjectiveTasks` — three columns of Task cards labelled by their title, the one badge mapping from typed `data` values (stage, clearance, gate blockers, waiver, exception) to glyph, label, and accent, the Next area — one compact block formatted from the single resolved `data.Next`, naming the rung, the selected records, `internal/resume`'s evidence and selection-diagnostic wording, an Issues count by type, and the action, with nothing on it derived from the index or moved by the sidebar's selection — the Task and Objective detail overlay, split into the resolution that reaches the index (identity, lifecycle, ownership, dependency decisions, clearance, the Check chain with latest and superseded marked, and the Issues the index's link maps hang off the record) and the rendering that reaches nothing, scrolls, and returns focus to the surface it was opened from — and non-TTY output leading with those same Next lines |
| `internal/buildtool/` | Makefile helper, named Go-test gates and timing summaries, cross-compile including Windows targets, archives, distribution checksums |
| `internal/doctor/` | Read-only project diagnostics, integrity checks, Release readiness through the canonical Release completion resolver, defect validation, timed quality gate execution, report formatting, typed repair suggestions |
| `internal/data/` | Task/router/defect models, frontmatter parsing/splitting, lifecycle validation/defaulting, discovery including root-dir and release defect traversal, unified task status constants, canonical write helpers, audit-register models/loaders and finding backlink lookups, the V2 Next projection (`ResolveNext`): the precedence ladder and selection resolution over the E43/E44 gate resolvers, deriving one next action for a project without owning any gate rule itself, plus the Issues relevant to that selection (`Next.Issues`), resolved from the index's own Task/Check/Issue link maps |
| `internal/resume/` | Deterministic, plain-text rendering of a resolved `data.Next` projection to an `io.Writer` (`resume.Render`): selected Objective/Task identity, implementation state, technical clearance, owner-wait, exception, dependency, and Issue phrasing, and the next action — with the evidence and freshness wording held once in its own file. It owns that wording for every surface reporting recorded V2 evidence, not just for its own narrative: `EvidenceLines`, `ActionPhrase`, and `SelectionPhrase` are exported so the V2 board's Next area states the same facts in the same words under a layout of its own, and `ClearancePhrase`, `DependencyPhrase`, `ObjectiveDependencyPhrase`, `ExceptionPhrase`, `ReplanPhrase`, `IssueLine`, and `ActorLabel` are exported for the board's detail overlay, which reports one record's own evidence rather than a whole projection. No filesystem, subprocess, network, or TTY access; no project root or index consulted |
| `internal/migrate/` | One-time V1-to-V2 project conversion; first-class Release PRD/source mapping; the platform file replacement primitive that operation's writes go through; a read-only source inventory (exact-byte hashing, confined walk); role classification against the frozen fixture vocabulary; and the migrate command body: target validation, preview, decisions loading, guarded apply, interruption recovery, and no-op retry proof |
| `internal/data/release_cutover.go` | Project-level E50 cutover composition: orders every declared Release's `ResolveReleaseCompletion` blockers without introducing a second readiness rule; no-Release projects remain valid |
| `internal/board/v2/releases.go` | Optional `r` Release selector and read-only Release detail path; membership and readiness are derived from the V2 index and canonical data resolvers |
| `internal/testutil/` | Shared Go test fixtures and filesystem helpers for internal package tests |
| `internal/styles/` | Atari-Noir palette, TUI styles |
| `templates/` | Scaffold markdown, YAML, prompts, and defect workflow guidance |
| `agent-skills/` | Phase-specific skill guides, including defect capture guidance |

The V2 `internal/data/` boundary now also owns first-class Release records,
derived Release membership, the canonical per-Release completion decision,
and the `ResolveReleaseCutover` composition consumed by E50. `internal/doctor`
and `internal/board/v2` report that same decision; they do not maintain a
second Release-readiness policy.

After this repository's schema-2 migration, ordinary startup, board, doctor,
resume, and init behavior is V2-only. The V1 readers and V1 board surfaces in
the map are compatibility code reachable only from explicit migration or
preserved historical fixtures.

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

`E47` ships this table as the scaffold default for new V2 projects; E50 activates it here after migration. V1 skills remain available only for the V1 scaffold/upgrade path and byte-preserved history.

## Legacy V1 compatibility (not active)

The following contract is retained only for archived V1 projects and the V1 scaffold; it is not an active route in this schema-2 repository.

| task-building | savepoint-build-task |
| audit-pending | savepoint-audit-epic |

An explicit request uses `savepoint-audit-task` while `state` stays `task-building`; that is not a router state and not a new state here.

Task `stage` (build/test/audit): **required** when `status: in_progress` — Task lifecycle rules are owned by `internal/data`; legacy `phase` is parse compatibility only and must not be used in new task guidance. Only the user may set a task to `status: done`.
