# Agents Guide

Welcome, AI agent. This project (`savepoint`) uses its own conventions to manage its own build. You are about to dogfood the workflow.

## Workflow

Before doing anything, read `.savepoint/router.md`. That file routes you to the next file based on the project's current state.

**Available Custom Skills:**
Savepoint ships with skills that define your role at each stage.
- `savepoint-draft-prd`: Acts as a strict Product Manager to interrogate the user and write the PRD.
- `savepoint-system-design`: Acts as a Staff Engineer mapping the architecture.
- `savepoint-create-plan` / `savepoint-create-task`: Acts as a Technical PM breaking work into Epics and ACs.
- `savepoint-build-task`: Acts as the disciplined execution engine writing code and logging drift.
- `savepoint-audit`: Acts as the QA Lead reconciling code with documentation.

**Skill Activation (CRITICAL):**
When you read `.savepoint/router.md`, you MUST activate the corresponding agent skill for the current state before taking any action. Use the `activate_skill` tool (or equivalent) for the appropriate phase:
- `state: draft-prd` -> Activate the `savepoint-draft-prd` skill.
- `state: design` -> Activate the `savepoint-system-design` skill.
- `state: planning` -> Activate the `savepoint-create-plan` skill.
- `state: task-breakdown` -> Activate the `savepoint-create-task` skill.
- `state: in-progress` -> Activate the `savepoint-build-task` skill.
- `state: audit-pending` -> Activate the `savepoint-audit` skill.

When you are about to write code, you must first read, in order:

1. `.savepoint/router.md` — current state and next action
2. The active epic E##-Detail.md: `.savepoint/releases/{release}/epics/{E##-epic}/E##-Detail.md`
3. The active task file: `.savepoint/releases/{release}/epics/{E##-epic}/tasks/{T###}-*.md`
4. Directly touched source/test files

Read `.savepoint/PRD.md` only for project vision changes, major scope questions, or when the router explicitly asks for it.
Read `.savepoint/Design.md` only when the task changes architecture or audit state. Read `.savepoint/releases/{release}/{release}-PRD.md` only when planning epics, changing release scope, or resolving epic order.

**Conditional read:** if the active task touches TUI implementation, also read `agent-skills/ink-tui-design/SKILL.md` as the execution guide. If it touches TUI rendering, theme, or visual design, also read `.savepoint/visual-identity.md` as the visual guardrails. Otherwise skip the extra files — they are tokens you do not need.

**Do not load files outside the current task scope** unless the task requires it. Token discipline is the wedge of this product; we honor it on ourselves.

Planning and implementation are separate handoffs:

- Epic task breakdown and detailed task planning happen together in one pass by one planning agent.
- Each task file must be independently buildable, objective-led, include explicit `depends_on` IDs, contain `## Acceptance Criteria` (observable outcomes) before `## Implementation Plan` (build checklist), and include a `## Context Log` for files read, estimated input tokens, and notes.
- Implementation happens one task at a time and may be handed to any agent. Clear context between tasks by default; rehydrate only from the router, active epic E##-Detail.md, active task file, and directly touched source/test files.
- During implementation, run focused tests for the touched behavior first; reserve the full quality-gate suite for task closeout.

- After all tasks in an epic are `done`, hand the epic back for audit.
- Any explicit audit request overrides the normal handoff timing for that epic. Persist the audit to `.savepoint/audit/{release}/{E##-epic}/snapshot.md` and `.savepoint/audit/{release}/{E##-epic}/proposals.md` before replying; do not stop at chat-only findings.

## Task Status Canon

Task frontmatter `status` must be exactly one of `planned`, `in_progress`, or `done`.

Active task phase is represented separately with `phase: build`, `phase: test`, or `phase: audit`, and `phase` is only valid when `status: in_progress`.

Never write `todo`, `doing`, `blocked`, `review`, `audit`, or phase names into `status`. If a task is blocked, keep its canonical status and document the blocker in the body.

## Task Completion Protocol

When a task reaches `status: done`, you MUST:

1. Verify every `## Acceptance Criteria` line has a passing test or verified manual outcome. A task is not done until its acceptance criteria are satisfied, not merely its implementation checkboxes ticked.
2. Tick all checkboxes in the `## Implementation Plan`.
3. Fill the `## Context Log` (files read, estimated input tokens, notes).
4. Run the full quality-gate suite (`make build && make test`). Record the result in the Context Log.
5. If any gate fails, fix it or document the blocker in the task file before setting `status: done`.
6. Set the task frontmatter to `status: done`.
7. Update `router.md` with the next action (next unblocked task, or `audit-pending` if all tasks done).
8. **Stop. Prompt the user:**
   > "Task {id} is done. Quality gates: {pass/fail list}. Router updated to {next_action}. Review the changes, then tell me to continue."

**Do not start the next task. Do not advance past this point without user acknowledgment.**

## Task Closeout Meta-Check

After marking a task `done` and before prompting the user, ask yourself:

- Did this task add new source files, modules, or exports not in the Codebase Map?
- Did this task change the architecture from what `.savepoint/Design.md` describes?

If yes, append a `## Drift Notes` section to the task file:
  - `Drift: {file} added, not yet in Codebase Map.`
  - `Drift: {section} in Design.md may need update.`

Drift notes are lightweight annotations. They do **not** replace the epic audit. They flag what the next audit should reconcile.

## Audit Handoff Rule

The agent session that builds an epic **must not** run its audit. Audit requires fresh eyes.

When all tasks in an epic are `done`:
1. Update `router.md` to `state: audit-pending` for that epic.
2. Stop. Tell the user: "Epic {id} is complete. Start a new agent session for the audit."
3. The user starts a fresh session. The new agent reads `router.md`, sees `audit-pending`, and follows the audit-reconciliation instructions.

**If you are in the same session that built the epic, you must not audit it.**

## Build / Test / Run

```bash
make build        # go run ./internal/buildtool build
make test         # go test ./...
make run          # go run main.go
make clean        # go run ./internal/buildtool clean
make build-all    # cross-compile linux-amd64, linux-arm64, darwin-amd64, darwin-arm64 into dist/
make dist         # build-all + create versioned tar.gz archives in dist/ (VERSION= override supported)
make smoke-test   # build local binary and run --version to validate exit 0
```

## Code Style

1. **One job per file.** If a file does two things, split it.
2. **One-sentence rule.** If you can't describe a function in one sentence, refactor.
3. **Test what branches.** Logic with if/else/switch gets a test. Pure rendering: skip.
4. **Types are documentation.** No `any`. Let the compiler help.
5. **Build, don't speculate.** No code for hypothetical futures.
6. **Errors at boundaries.** Handle failure where data enters or leaves the system.
7. **One source of truth.** State lives in one place. No syncing copies.
8. **Comments explain WHY,** not what. If removing the comment wouldn't confuse a future reader, delete it.
9. **Content in data files.** Markdown/JSON/YAML, not strings in code.
10. **Small diffs.** Each task touches as few files as possible.

## Codebase Map

| Module                               | Purpose                                                                                              |
| ------------------------------------ | ---------------------------------------------------------------------------------------------------- |
| `main.go`                            | CLI entrypoint, root command wiring, and `--version` handling via build-time version injection        |
| `internal/board/`                    | TUI board models, layout, rendering, overlays, focusable epic sidebar navigation, epic detail overlays, epic status glyph loading, task transitions, router priority markers, and fsnotify refresh |
| `internal/buildtool/`                | Go-native Makefile helper for cleanup, local builds, Linux/Darwin cross-compilation, archives, and smoke tests |
| `internal/data/`                     | Task/router/config models, frontmatter parsing, checklist state parsing, mtime-guarded writes, discovery, and generic file readers |
| `internal/styles/`                   | Atari-Noir palette constants, terminal color fallbacks, shared TUI styles, stable column border styles, scroll indicators, purple epic navigation/detail styles, semantic glyph/tag styles, and footer/header styling |
| `cmd/`                               | Additional CLI subcommands (if any)                                                                  |
| `templates/`                         | Default project scaffold markdown, YAML assets, and agent prompt templates                           |
| `agent-skills/`                      | Custom skill guides for different agent phases (`draft-prd`, `audit`, etc.)                          |

## CLI rules for agents

**Never run `savepoint` commands.** The CLI is for the human. Edit files directly.

(For this repo specifically: `savepoint` doesn't exist yet — we're building it. Even once it exists, this rule stands.)

## Recommended planning models

For PRD/Design/Task planning, this workflow assumes a top-tier model: Claude Opus, Gemini 2.5 Pro, GPT-5.5, or equivalent. Lighter models may not follow embedded prompt instructions reliably. If you are not one of those, advise the user before proceeding with planning steps.
 instructions reliably. If you are not one of those, advise the user before proceeding with planning steps.
