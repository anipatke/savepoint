# Agents Guide

## Workflow

1. Read `.savepoint/router.md` — state + next action
2. Activate skill per table below
3. Read: router → epic → task → source files

The phase skill is the canonical workflow source. This guide defines routing, terminology, and repo rules only; do not duplicate phase-by-phase prompt instructions here.

## Skill Activation

| State | Skill |
|-------|-------|
| pre-implementation | savepoint-draft-prd |
| epic-design | savepoint-system-design |
| epic-task-breakdown | savepoint-create-task |
| task-building | savepoint-build-task |
| audit-pending | savepoint-audit-epic |
| defect-building | savepoint-build-task |

An explicit request to audit or re-audit one in-progress task uses `savepoint-audit-task` while `state` stays `task-building`. It is a request-qualified override of the phase skill, not a router state.

Use `savepoint-create-defect` when the user reports a concrete bug, regression, or broken expectation that should be captured as a release-level defect before repair starts.

Use the `skill` tool when the listed skill is available. If the agent says the skill is not found, read `agent-skills/{skill}/SKILL.md` directly and follow it as the active skill.

Read `.savepoint/PRD.md` only for vision changes, `.savepoint/Design.md` only for architecture/audit.

## Terminology

- Router `state`: the current phase, such as `epic-design`, `task-building`, or `audit-pending`
- Task `status`: only `planned`, `in_progress`, or `done`
- Task `stage` (build/test/audit): **required** when `status: in_progress` — omitting it self-heals to `stage: build` on load and is flagged by `savepoint doctor`
- Task lifecycle rules are owned by `internal/data`; legacy `phase` is parse compatibility only and must not be used in new task guidance.
- Never write `stage: implementation`; use `stage: build` when starting implementation work.
- Never: todo, doing, blocked, review, audit
- Agents may set a task to `status: in_progress` when starting implementation.
- Only the user may set a task to `status: done` or retreat a task to an earlier status.

## Defect Workflow

Use a defect conversation when the user reports a concrete bug, regression, broken behavior, or failed expectation that should be repaired without reshaping the planned epic/task backlog.

- Defects live at `.savepoint/releases/{release}/defects/D###-slug.md`.
- Use `agent-skills/savepoint-create-defect/SKILL.md` to capture a new defect file.
- Router state may be `defect-building` with a `defect` field naming the active defect id.
- Defect lifecycle: `open` → `in_progress` (requires `stage: build|test|audit`) → `resolved`. Never use task-style `planned` or `done` in defect files.
- Use the board `d` overlay to inspect defects; do not turn defects into a fourth task column.

## Implementation

Follow the active skill for execution. During `task-building`, the canonical flow is `savepoint-build-task` — it owns the read order, `status: in_progress` + `stage: build` setting, AC verification, quality gates, and handoff.

**Stop. Prompt the user before continuing.** Only the user may mark a task `status: done` or retreat a task to an earlier status.

## Drift Check

- New files/modules not in Codebase Map?
- Architecture changed from Design.md?

If yes → append `## Drift Notes` to task file.

## Audit

Audit is agent-led and split by intent:

- `savepoint-audit-epic` — the `audit-pending` phase workflow, or an explicit audit of a completed epic. It requires a session independent from the builder, runs the Full health check, and writes the single `E##-Audit.md` handoff file. The builder must not audit its own epic; start a fresh session.
- `savepoint-audit-task` — an explicit request to audit or re-audit one in-progress task. Router `state` stays `task-building`, the review is read-only, it runs the Quick health check, and it returns `CLEAR` or `NEEDS WORK` without writing any file.

Both skills load `agent-skills/references/audit-method.md`, the shared non-triggerable audit method: scope locks, coverage matrices, workflow and side-effect locks, adversarial pass, re-audit convergence, and materiality.

When the project has `.savepoint/Guardrails.md` (policy) and `.savepoint/Health-Check.md` (evidence modes), both audits apply them — Quick at task handoff and task audit, Full at epic audit. Skip the related step when either file is absent; absence is not a finding.

- Audit file: `.savepoint/releases/{release}/epics/{E##-slug}/E##-Audit.md`
- During audit apply/close, update the same `E##-Audit.md` visible sections so `## Main Findings` and `## Code Style Review` describe the applied outcome, not stale pre-apply blockers.
- When `.savepoint/audit/` exists, also follow the `savepoint-audit-register` skill: read `.savepoint/audit/prompt.md`, reconcile against `.savepoint/audit/register.md` with stable `F###` IDs, and record the run — do not restart from a cold scan.

## Code Style

Code style is project-owned policy: the `STYLE` rules in `.savepoint/Guardrails.md` are the single source of truth. Read them when writing or reviewing code. They are Guideline severity and advisory — they inform review but never block on their own. If the project has no `.savepoint/Guardrails.md`, code style is not defined for it.

## Build

```bash
make build && make test
```

## Codebase Map

| Module | Purpose |
|--------|---------|
| `main.go` | CLI entrypoint, --version, embedded template wiring for init and upgrade-assets, migrate dispatch, and resume dispatch: resolves the target directory and any pending migration state, loads the V2 project and router, resolves the `Next` projection, and renders it (`runResume`) |
| `cmd/` | CLI command arg parsing and dispatch for init, board, doctor, upgrade-assets, migrate, and resume — argument parsing only, no record parsing, gate reading, or rendering (ARCH-01). It parses `board`'s V1 filters (`--release`, `--epic`) and V2 filter (`--objective`) without judging which applies; that is refused behind the dispatch, where the schema is known |
| `internal/init/` | Target validation, scaffold writing from a caller-selected template tree (`init` defaults to `templates/project-v2`), upgrade-assets schema-version dispatch between `templates/project` and `templates/project-v2` via `data.ReadSchemaVersion`, upgrade provenance manifest, managed AGENTS.md merge/conflict behavior, and safe project asset refresh |
| `internal/board/` | The board's schema dispatch — resolve the project root once, `data.LoadProject`, then run the V1 board or `internal/board/v2`, refusing a filter flag the resolved schema has no meaning for (CFG-01) — plus the whole V1 board: TUI board, overlays, epic sidebar, Next Activity line, router priority key, detail checklist rendering, status glyphs, forced color profile, debug logging hooks, async update I/O commands, defect summary/overlay/detail rendering, related-defect card markers, audit register overlay with finding detail and linked-finding backlinks, shared board utilities |
| `internal/board/v2/` | The V2 board, a sibling package no V1 board type or V1 record type is reachable from: model state, the single load command (`data.LoadProject`, `ReadStateV2`, `migrate.PendingOperation`, `data.ResolveNext`) startup and every reload share, the load-diagnostic screen for a project the V2 index refuses, the Objective sidebar — list, cursor, selection, per-Objective status/clearance/integration/dependency state, and Task filtering by ownership read from `index.ObjectiveTasks` — three columns of Task cards labelled by their title, the one badge mapping from typed `data` values (stage, clearance, gate blockers, exception) to glyph, label, and accent, the Next area — one compact block formatted from the single resolved `data.Next`, naming the rung, the selected records, `internal/resume`'s evidence and selection-diagnostic wording, an Issues count by type, and the action, with nothing on it derived from the index or moved by the sidebar's selection — the Task and Objective detail overlay, split into the resolution that reaches the index (identity, lifecycle, ownership, dependency decisions, clearance, the Check chain with latest and superseded marked, and the Issues the index's link maps hang off the record) and the rendering that reaches nothing, scrolls, and returns focus to the surface it was opened from — and non-TTY output leading with those same Next lines |
| `internal/buildtool/` | Makefile helper, cross-compile including Windows targets, archives, distribution checksums |
| `internal/doctor/` | Read-only project diagnostics, integrity checks, defect validation, timed quality gate execution, report formatting, typed repair suggestions |
| `internal/data/` | Task/router/defect models, frontmatter parsing/splitting, lifecycle validation/defaulting, discovery including root-dir and release defect traversal, unified task status constants, canonical write helpers, audit-register models/loaders and finding backlink lookups, the V2 Next projection (`ResolveNext`): the precedence ladder and selection resolution over the E43/E44 gate resolvers, deriving one next action for a project without owning any gate rule itself, plus the Issues relevant to that selection (`Next.Issues`), resolved from the index's own Task/Check/Issue link maps |
| `internal/resume/` | Deterministic, plain-text rendering of a resolved `data.Next` projection to an `io.Writer` (`resume.Render`): selected Objective/Task identity, implementation state, technical clearance, owner-wait, exception, dependency, and Issue phrasing, and the next action — with the evidence and freshness wording held once in its own file. It owns that wording for every surface reporting recorded V2 evidence, not just for its own narrative: `EvidenceLines`, `ActionPhrase`, and `SelectionPhrase` are exported so the V2 board's Next area states the same facts in the same words under a layout of its own, and `ClearancePhrase`, `DependencyPhrase`, `ObjectiveDependencyPhrase`, `ExceptionPhrase`, `ReplanPhrase`, `IssueLine`, and `ActorLabel` are exported for the board's detail overlay, which reports one record's own evidence rather than a whole projection. No filesystem, subprocess, network, or TTY access; no project root or index consulted |
| `internal/migrate/` | One-time V1-to-V2 project conversion; the platform file replacement primitive that operation's writes go through, a read-only source inventory (exact-byte hashing, confined walk), role classification against the frozen fixture vocabulary, and the migrate command body: target validation, preview, decisions loading, and the guarded apply |
| `internal/testutil/` | Shared Go test fixtures and filesystem helpers for internal package tests |
| `internal/styles/` | Atari-Noir palette, TUI styles |
| `templates/` | Scaffold markdown, YAML, prompts, and defect workflow guidance |
| `agent-skills/` | Phase-specific skill guides, including defect capture guidance |

## Context Budget

- **Read only what you need.** Each phase has a strict read budget. Do not read files outside your current phase's context.
- **No exploratory reads.** Read only the files listed in the task's `## Context Files`. Do not glob or search for new information unless explicitly instructed.
- **Token awareness.** Every file read consumes context window. Before reading a file, ask: "Do I need this to complete my current phase?"

## CLI Rules

**Never run `savepoint` commands.** The CLI is for the human. Edit files directly.

## V2 Routing

The table below is the live routing model for a V2 project — a project whose `config.yml` declares `schema_version: 2`. It is active for V2 projects; it is not active for this repository until `E50` migrates this repository itself onto the V2 lifecycle. Until then, the `## Skill Activation` table above governs how this repository routes work, and nothing below changes that.

| Router `state` | Skill |
|-----------------|-------|
| idea | savepoint-idea |
| design | savepoint-design |
| task | savepoint-task |
| check | savepoint-check |

`REPLAN REQUIRED`, returned by an executor that hits a materially invalid plan, routes back into `savepoint-design`. It is not a fifth router state — the state stays `design` while the planner resolves what broke.

Three shared references back these four skills: `agent-skills/references/check-method.md`, `agent-skills/references/issue-capture.md`, and `agent-skills/references/commands-and-procedures.md`. Each carries `triggerable: false` frontmatter and is non-triggerable on its own — it is loaded in full by the skill that owns it (`savepoint-check` loads `check-method.md`; `savepoint-design`, `savepoint-task`, and `savepoint-check` each enter `issue-capture.md` from their own workflow; `savepoint-design` loads `commands-and-procedures.md` for config reconciliation), not invoked directly.

`E47` ships this table as the scaffold default for new V2 projects and retires the nine V1 skills from a project that migrates. `E50` migrates this repository itself onto the V2 lifecycle and removes the transitional V1 readers. Neither happens in this task.
