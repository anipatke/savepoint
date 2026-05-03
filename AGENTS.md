# Agents Guide

## Workflow

1. Read `.savepoint/router.md` — state + next action
2. Activate skill per table below
3. Read: router → epic → task → source files

## Skill Activation

| State | Skill |
|-------|-------|
| pre-implementation | savepoint-draft-prd |
| epic-design | savepoint-system-design |
| epic-task-breakdown | savepoint-create-task |
| task-building | savepoint-build-task |
| audit-pending | savepoint-audit |

Use the `skill` tool when the listed skill is available. If the agent says the skill is not found, read `agent-skills/{skill}/SKILL.md` directly and follow it as the active skill.

Read `.savepoint/PRD.md` only for vision changes, `.savepoint/Design.md` only for architecture/audit.

## Task Status

- `status`: only `planned`, `in_progress`, or `done`
- `stage` (build/test/audit): **required** when `status: in_progress` — omitting it is a parse error
- Never: todo, doing, blocked, review, audit
- Agents may set a task to `status: in_progress` when starting implementation.
- Only the user may set a task to `status: done` or retreat a task to an earlier status.

## Implementation

1. Read task's `## Context Files` using `Read` tool — one call per file, no explore, no glob
2. Read task's `## Acceptance Criteria` + `## Implementation Plan`
3. When starting implementation, set task frontmatter to `status: in_progress` + `stage: build` (both required together)
4. After setting `in_progress`, press `p` in the TUI to mark the focused task as router priority
5. Execute in order, tick checkboxes
6. Verify every AC has passing test/outcome
7. Run quality gates (build + test)
8. Update router.md: next task or `audit-pending`
9. **Stop. Prompt user before continuing.**

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
| `main.go` | CLI entrypoint, --version |
| `cmd/` | CLI command arg parsing and dispatch for init, board, and doctor |
| `internal/init/` | Target validation, scaffold writing from templates |
| `internal/board/` | TUI board, overlays, epic sidebar, Next Activity line, router priority key, detail checklist rendering, status glyphs, forced color profile, async update I/O commands, shared board utilities |
| `internal/buildtool/` | Makefile helper, cross-compile, archives |
| `internal/doctor/` | Read-only project diagnostics, integrity checks, timed quality gate execution, report formatting, typed repair suggestions |
| `internal/data/` | Task/router models, frontmatter parsing/splitting, lifecycle validation/defaulting, discovery, canonical write helpers |
| `internal/styles/` | Atari-Noir palette, TUI styles |
| `templates/` | Scaffold markdown, YAML, prompts |
| `agent-skills/` | Phase-specific skill guides |

## CLI Rules

**Never run `savepoint` commands.** The CLI is for the human. Edit files directly.
