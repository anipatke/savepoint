# Agents Guide

> Legacy V1 scaffold: this asset is retained for migration and historical
> fixtures. New schema-2 projects use `templates/project-v2/AGENTS.md`; the V1
> routing below is not active in the migrated repository.

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

Use the `skill` tool when the listed skill is available. If the agent says the skill is not found, read `agent-skills/{skill}/SKILL.md` directly and follow it as the active skill.

Use `savepoint-create-defect` when the user reports a concrete bug, regression, or broken expectation that should be captured as a release-level defect before repair starts.

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

## Audit Register

`.savepoint/audit/` is the durable register of repo-wide audit findings: `prompt.md` (canonical audit prompt), `register.md` (current reconciled state), `findings/` (one file per stable `F###` finding), and `runs/` (immutable run history).

- When `.savepoint/audit/` exists and audit work starts, follow the `savepoint-audit-register` skill and read `.savepoint/audit/prompt.md` before recording anything.
- Reconcile against `.savepoint/audit/register.md` instead of restarting from a cold scan; a finding seen again keeps its `F###` ID.
- A finding reaches `verified` only with named proof. Only the user grants `waived` or `owner_decision` dispositions.
- The board `A` overlay is read-only review; the markdown files stay the source of truth.

## Code Style

Code style is project-owned policy: the `STYLE` rules in `.savepoint/Guardrails.md` are the single source of truth. Read them when writing or reviewing code. They are Guideline severity and advisory — they inform review but never block on their own. If the project has no `.savepoint/Guardrails.md`, code style is not defined for it.

## Build

```bash
make build && make test
```

## Codebase Map

| Module | Purpose |
|--------|---------|

## Context Budget

- **Read only what you need.** Each phase has a strict read budget. Do not read files outside your current phase's context.
- **No exploratory reads.** Read only the files listed in the task's `## Context Files`. Do not glob or search for new information unless explicitly instructed.
- **Token awareness.** Every file read consumes context window. Before reading a file, ask: "Do I need this to complete my current phase?"

## CLI Rules

**Never run `savepoint` commands.** The CLI is for the human. Edit files directly.

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
