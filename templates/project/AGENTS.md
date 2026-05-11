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
- Task `stage` (build/test/audit): **required** when `status: in_progress` — omitting it is a parse error
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

1. Read the task's `## Context Files` one file at a time; do not explore, glob, search broadly, or read files outside the task context unless explicitly required.
2. Read the task's `## Acceptance Criteria` and `## Implementation Plan`.
3. Set task frontmatter to `status: in_progress` and `stage: build`, then press `p` in the TUI to mark router priority.
4. Follow `savepoint-build-task` for execution, checklist updates, AC verification, quality gates, context log, and handoff.
5. Stop and prompt the user before continuing.

## Drift Check

- New files/modules not in Codebase Map?
- Architecture changed from Design.md?

If yes → append `## Drift Notes` to task file.

## Audit

- The builder must not audit its own epic; start a fresh audit session.
- Audit is agent-led via `savepoint-audit`, not a `savepoint audit` CLI pipeline.
- Write exactly one `.savepoint/releases/{release}/epics/{E##-slug}/E##-Audit.md`.
- The TUI Audit tab renders only `## Main Findings` and `## Code Style Review`.
- Put file-specific `### Target File`, `### Replace`, and `### With` blocks under `## Proposed Changes`.
- During audit apply/close, update the same `E##-Audit.md` visible sections so they describe the applied outcome, not stale blockers.

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

Build gate: `make build && make test`

## Codebase Map

| Module | Epic | Purpose |
|--------|------|---------|

## CLI Rules

**Never run `savepoint` commands.** The CLI is for the human. Edit files directly.
