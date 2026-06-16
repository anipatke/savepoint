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
| audit-pending | savepoint-audit |
| defect-building | savepoint-build-task |

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

Audit is agent-led via the `savepoint-audit` skill — follow it for the file layout, section rules, and apply/close flow. The builder must not audit its own epic; start a fresh session.

- Audit file: `.savepoint/releases/{release}/epics/{E##-slug}/E##-Audit.md`
- During audit apply/close, update the same `E##-Audit.md` visible sections so `## Main Findings` and `## Code Style Review` describe the applied outcome, not stale pre-apply blockers.

## Code Style

- **One job per file** — split files when responsibilities mix.
- **One job per function** — small, named, testable units.
- **Test branches** — cover meaningful conditionals and edge cases.
- **Types document intent** — prefer explicit types over comments.
- **Build only what is needed** — no speculative abstractions.
- **Handle errors at boundaries** — validate inputs, APIs, IO, and external data.
- **One source of truth** — no duplicated rules, constants, state, or config.
- **Comments explain why** — not what the code already says.
- **Content lives in data** — keep copy/config out of logic.
- **Small diffs** — minimal, reviewable, behaviour-preserving changes.

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
