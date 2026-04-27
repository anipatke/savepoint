# Agent State Machine

This file routes the agent based on the project's current state. Read this whenever you start a session.

## Read order on every session

1. This file (you are here)
2. The current state below to know what to do next
3. The active epic Design
4. The active task file, when a task is selected

Read `.savepoint/PRD.md` only for project vision changes. Read `.savepoint/Design.md` only for architecture changes or audit closeout. Read `.savepoint/releases/v1/PRD.md` only for release planning or epic-order questions.

**Conditional read (token discipline):** if your active task touches **TUI rendering, theme, or visual design**, also read `.savepoint/visual-identity.md` after Design.md. Otherwise skip it — it's ~1.8K tokens you don't need.

## Current state

```yaml
state: epic-task-breakdown
release: v1
epic: E02-data-model
next_action: "Break E02-data-model into T### task files only. Read only the E02 Design unless an explicit dependency requires more context."
```

## State → next action

<!-- AGENT: Read the state above. Find the matching block below. Follow it. -->

### `state: pre-implementation`

The project has its PRD and Design locked but no epics defined yet.

**Next action:**

1. Read `.savepoint/releases/v1/PRD.md` — the v1 release scope (epic list lives there).
2. Help the user define the epics list and confirm priority.
3. For each epic in order, create the directory `.savepoint/releases/v1/epics/E##-{epic-name}/` with a `Design.md` stub.
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

Epic Design exists but no tasks yet.

**Next action:**

1. Re-read the epic Design.
2. Propose a task list — each task **independently buildable**, **objective-led**, with declared `depends_on`.
3. Create `.savepoint/releases/v1/epics/{E##-epic}/tasks/TNNN-slug.md` for each, with frontmatter:
   ```yaml
   ---
   id: {E##-epic}/TNNN-slug
   status: backlog
   objective: "<one sentence>"
   depends_on: []
   ---
   ```
4. When the user approves the list, transition to `state: task-planning`.

### `state: task-planning`

Tasks exist in `backlog`. User has selected one to plan.

**Next action:** Read the task's `objective`. Write the implementation plan as inline `- [ ]` checkboxes under a `## Implementation Plan` heading. Set `status: planned`. Stop.

### `state: task-building`

Task is `in_progress`. All `depends_on` are `done`.

**Next action:** Execute the plan. Tick checkboxes as you complete them. Edit code per the **Code Style** rules in `AGENTS.md`. When all checkboxes tick, set `status: review` and stop.

### `state: audit-pending`

The last task in an epic is `done`. Audit must run before the next epic starts.

**Next action:** Read `.savepoint/audit/{E##-epic}/snapshot.md`. Read the epic's `Design.md`. Read only the files listed as changed. Write patch-shaped proposed updates to `.savepoint/audit/{E##-epic}/proposals/`:

- `Design.md` — merge epic delta into project architecture
- `AGENTS.md` — refresh the Codebase Map section between markers
- `epic-Design.md` — add "implemented as:" section noting deltas from the original plan
- `quality-review.md` — semantic-review findings against the 10 Code Style rules (advisory only)

Proposal format:

```md
## Target File

`path/to/file.md`

## Replace

<exact old heading, marker, or section>

## With

<new content>
```

Quality review format:

```md
## Must Fix Before Close

## Carry Forward

## Already Fixed
```

After proposals are approved, apply approved proposals to live files, mark the epic `Design.md` as `status: audited`, update project `Design.md` `last_audited`, refresh `AGENTS.md` Codebase Map, and advance this router to the next epic state.

Stop. The user reviews proposals in the TUI before commit actions.

## Capability check

If you are not Claude Opus / Gemini 2.5 Pro / GPT-5.5 / equivalent, surface a warning to the user:

> _"Heads up — I'm running on a lighter model. Savepoint's planning steps work best with top-tier models because the embedded prompts are detailed. I'll do my best, but consider switching the model for PRD/Design/Task-breakdown steps."_

Then proceed.
