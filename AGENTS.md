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

Use `savepoint-create-defect` when the user reports a concrete bug, regression, or broken expectation that should be captured as a release-level defect before repair starts.

Use the `skill` tool when the listed skill is available. If the agent says the skill is not found, read `agent-skills/{skill}/SKILL.md` directly and follow it as the active skill.

Read `.savepoint/PRD.md` only for vision changes, `.savepoint/Design.md` only for architecture/audit.

## Terminology

- Router `state`: the current phase, such as `epic-design`, `task-building`, or `audit-pending`
- Task `status`: only `planned`, `in_progress`, or `done`
- Task `stage` (build/test/audit): **required** when `status: in_progress` — omitting it is a parse error
- Task lifecycle rules are owned by `internal/data`; legacy `phase` is parse compatibility only and must not be used in new task guidance.
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

1. Read task's `## Context Files` using `Read` tool — one call per file, no explore, no glob
2. Read task's `## Acceptance Criteria` + `## Implementation Plan`
3. When starting implementation, set task frontmatter to `status: in_progress` + `stage: build` (both required together)
4. After setting `in_progress`, press `p` in the TUI to mark the focused task as router priority
5. Follow `savepoint-build-task` for execution, checklist updates, AC verification, quality gates, context log, and handoff.
6. **Stop. Prompt user before continuing.**

## Drift Check

- New files/modules not in Codebase Map?
- Architecture changed from Design.md?

If yes → append `## Drift Notes` to task file.

## Audit Handoff

The agent that builds an epic **must not audit it**. Start a fresh session.

## Audit File Structure

- Audit is agent-led via `savepoint-audit`, not a `savepoint audit` CLI pipeline.
- Write exactly one `.savepoint/releases/{release}/epics/{E##-slug}/E##-Audit.md`.
- The TUI Audit tab renders `## Main Findings` and `## Code Style Review` only.
- Keep file-specific `### Target File` / `### Replace` / `### With` blocks under `## Proposed Changes` so admin apply details do not appear in the Epic Detail panel.
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
| `main.go` | CLI entrypoint, --version, embedded template wiring for init and upgrade-assets |
| `cmd/` | CLI command arg parsing and dispatch for init, board, doctor, and upgrade-assets |
| `internal/init/` | Target validation, scaffold writing from templates, managed AGENTS.md merge behavior, and safe project asset refresh |
| `internal/board/` | TUI board, overlays, epic sidebar, Next Activity line, router priority key, detail checklist rendering, status glyphs, forced color profile, debug logging hooks, async update I/O commands, defect summary/overlay/detail rendering, related-defect card markers, shared board utilities |
| `internal/buildtool/` | Makefile helper, cross-compile including Windows targets, archives, distribution checksums |
| `internal/doctor/` | Read-only project diagnostics, integrity checks, defect validation, timed quality gate execution, report formatting, typed repair suggestions |
| `internal/data/` | Task/router/defect models, frontmatter parsing/splitting, lifecycle validation/defaulting, discovery including root-dir and release defect traversal, unified task status constants, canonical write helpers |
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
