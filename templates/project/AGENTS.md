# Agents Guide

## Startup

1. Read `.savepoint/router.md` for state, release, epic, task, and next action.
2. Activate the skill for the current state:

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

Read in order: router → active epic → active task → task `## Context Files` only. Read `.savepoint/PRD.md` only for vision changes, and `.savepoint/Design.md` only for architecture/audit.

## Task Status

- `status` must be only `planned`, `in_progress`, or `done`.
- `stage` (build/test/audit): **required** when `status: in_progress` — omitting it is a parse error
- Never: todo, doing, blocked, review, audit
- Agents may set a task to `status: in_progress` when starting implementation.
- Only the user may set a task to `status: done` or retreat a task to an earlier status.

## Defect Lane

Use a defect file when a bug is discovered **during or after a build phase** — not when a planned task is reworked. Defects are tracked separately from epics so the audit trail stays clean.

- **Location:** `.savepoint/releases/{release}/defects/D###-slug.md`
- **When to create:** regression found in TUI testing, broken build caused by a merged epic, or a production issue traced to a known release
- **When NOT to create:** a planned task that needs rework (update the task instead), a scope change (that is an epic), or a future enhancement
- **Capture skill:** `agent-skills/savepoint-create-defect/SKILL.md`

Defect frontmatter:

```yaml
---
id: {release}/D###-slug
release: {release}
status: planned          # planned | in_progress | done
severity: high           # critical | high | medium | low
title: "One-line description"
introduced: v1.0.3       # optional: version where bug appeared
reference: E12-slug/T003-slug  # optional: related task ID
---
```

Press `d` on the board to open the defect overlay. The `savepoint doctor` command validates defect files and reports malformed frontmatter, invalid status, missing in-progress stage, and broken references.

If the router is in `defect-building`, treat the session as repair work for the named defect rather than normal epic planning or task-building.

## Implementation

1. Read the task's `## Context Files` one file at a time; do not explore, glob, search broadly, or read files outside the task context unless explicitly required.
2. Read the task's `## Acceptance Criteria` and `## Implementation Plan`.
3. Set task frontmatter to `status: in_progress` and `stage: build`, then press `p` in the TUI to mark router priority.
4. Execute the plan in order, tick checkboxes, and verify every AC with a passing test or outcome.
5. Run `make build && make test`.
6. Update `router.md` to the next task or `audit-pending`.
7. Stop and prompt the user before continuing.

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
