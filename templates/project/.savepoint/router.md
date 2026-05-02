# Agent State Machine

This file routes the agent based on the project's current state. Read this whenever you start a session.

## Read order on every session

1. This file (you are here)
2. The current state below to know what to do next
3. The active epic Design
4. The active task file, when a task is selected

Read `.savepoint/PRD.md` only for project vision changes. Read `.savepoint/Design.md` only for architecture changes or audit closeout. Read `.savepoint/releases/{release}/{release}-PRD.md` only for release planning or epic-order questions.

**Conditional read (token discipline):** if your active task touches **Ink/TUI implementation**, also read `agent-skills/ink-tui-design/SKILL.md` after Design.md as the execution guide. If it touches **TUI rendering, theme, or visual design**, also read `.savepoint/visual-identity.md` as the visual guardrails. Otherwise skip the extra files — they are tokens you do not need.

## Current state

```yaml
state: pre-implementation
release: v{{RELEASE_NUMBER}}
epic: none
task: none
next_action: "The project has its PRD and Design locked but no epics defined yet. Help the user define the epics list and confirm priority."
```

## Manual audit override

If the user explicitly asks you to audit an epic, perform the audit for that epic even if the router has not reached `state: audit-pending` yet.

Persist the audit artifacts before replying:

- Ensure `.savepoint/audit/{E##-epic}/snapshot.md` exists. Create a manual snapshot once if needed.
- Write the proposal bundle to `.savepoint/audit/{release}/{E##-epic}/proposals.md`.
- Do not stop at chat-only findings. The filesystem artifact is part of the task output.

## State → next action

<!-- AGENT: Read the state above. Find the matching block below. Follow it. -->

### `state: pre-implementation`

The project has its PRD and Design locked but no epics defined yet.

**Next action:**

1. Read `.savepoint/releases/{release}/{release}-PRD.md` — the release scope (epic list lives there).
2. Help the user define the epics list and confirm priority.
3. For each epic in order, create the directory `.savepoint/releases/{release}/epics/E##-{epic-name}/` with a `Design.md` stub.
4. When epic E01 (scaffolding) is created, transition to `state: epic-design` for that epic.

**Do not** start writing code. We are still in planning.

### `state: epic-design`

An epic exists but its `Design.md` is empty or a stub.

**Next action:** Walk the user through filling out the epic's `Design.md`:

- What is this epic adding to the system?
- What components / files does it touch?
- What's the architectural delta vs the current state?

When complete, transition to `state: epic-task-breakdown` for this epic.

### `state: epic-task-breakdown`

Epic Design exists but tasks are missing or not fully planned.

**Next action:**

1. Re-read the epic Design.
2. Create or update the full epic task list — each task **independently buildable**, **objective-led**, with declared `depends_on`.
3. Each task file lives at `.savepoint/releases/{release}/epics/{E##-epic}/tasks/TNNN-slug.md` with frontmatter:
   ```yaml
   ---
   id: {E##-epic}/TNNN-slug
   status: planned
   objective: "<one sentence>"
   depends_on: []
   ---
   ```
4. In the same pass, write each task's `## Acceptance Criteria` (observable outcomes that define "done") and `## Implementation Plan` (the build checklist that satisfies them) as inline `- [ ]` checkboxes, plus a `## Context Log` section with `Files read`, `Estimated input tokens`, and `Notes` fields.
5. When every task is planned, transition to `state: task-building` for the first unblocked task.

### `state: task-planning`

Reserved for repair or late-added tasks. Normal epic planning happens during `state: epic-task-breakdown`.

**Next action:** Read the task's `objective`. Write `## Acceptance Criteria` (observable outcomes) before `## Implementation Plan` (build checklist). Use inline `- [ ]` checkboxes under the plan, add a `## Context Log` section, set `status: planned`, and stop.

### `state: task-building`

Task is `in_progress`. All `depends_on` are `done`.

**Next action:** Execute the plan. Tick checkboxes as you complete them. The implementation checklist exists to satisfy the task's acceptance criteria; every checked box should map to an observable outcome in `## Acceptance Criteria`. Edit code per the **Code Style** rules in `AGENTS.md`. Before setting `status: done`, update the task's `## Context Log`. When all checkboxes tick, run the full quality-gate suite, set `status: done`, update the router, and **stop for user review**.

**Do not start the next task without user acknowledgment.**

### `state: audit-pending`

The last task in an epic is `done`. Audit must run before the next epic starts.

**Context gate:** If you just built this epic in the current session, you **must not** audit it. Close this session. The user should start a new session for the audit.

**Next action (fresh session only):** Confirm `.savepoint/audit/{E##-epic}/snapshot.md` exists. If it is missing while the audit CLI is still unavailable, create one manual snapshot from the known epic scope once; do not search broadly for replacement inputs. Then read the snapshot, read the epic's `Design.md`, and read only the files listed as changed. Write one patch-shaped proposal bundle to `.savepoint/audit/{E##-epic}/proposals.md`:

- `Design.md` section — merge only the epic delta into project architecture.
- `AGENTS.md` section — refresh Codebase Map entries from changed-module metadata; preserve existing rows.
- `epic-Design.md` section — add "Implemented as:" notes and deltas from the original plan.
- `Quality Review` section — semantic-review findings against the 10 Code Style rules.

Prefer delta-only edits (`Insert After`, `Replace`, `Delete`) anchored to exact text. Do not quote and replace entire large sections unless the whole section genuinely changed.

Proposal format:

```md
## Target File

`path/to/file.md`

## Replace

<exact old heading, marker, or section>

## With

<new content>
```

Quality review section format:

```md
## Must Fix Before Close

## Must Fix Before Next Epic

## Carry Forward

## Already Fixed
```

After proposals are approved, apply approved proposals to live files, mark the epic `Design.md` as `status: audited`, update project `Design.md` `last_audited`, refresh `AGENTS.md` Codebase Map, and advance this router to the next epic state.

Stop. The user reviews proposals in the TUI before commit actions.

## Capability check

If you are not Claude Opus / Gemini 2.5 Pro / GPT-5.5 / equivalent, surface a warning to the user:

> _"Heads up — I'm running on a lighter model. Savepoint's planning steps work best with top-tier models because the embedded prompts are detailed. I'll do my best, but consider switching the model for PRD/Design/Task-breakdown steps."_

Then proceed.
 proceed.
