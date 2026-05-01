# Savepoint Skill: Build Task (`build-task`)

## Objective
Act as a disciplined coding agent that strictly follows Savepoint's implementation loop.

## Context
The `build-task` skill is the execution engine. It reads the detailed task plan, writes the code, and proves that the Acceptance Criteria (ACs) have been met. It is strictly constrained to the scope of the single active task. It does not rewrite architecture, and it does not fix unrelated bugs.

## Trigger
This skill is activated when the `.savepoint/router.md` state is `in-progress` and points to a specific task file.

## Input
- `.savepoint/router.md` (Current state).
- The active epic Design: `.savepoint/releases/v1/epics/{E##-epic}/Design.md`.
- The active task file: `.savepoint/releases/v1/epics/{E##-epic}/tasks/{T###}-*.md`.
- Directly touched source/test files.

## Workflow

1.  **Read the Task Plan:** Review the `## Acceptance Criteria` and the `## Implementation Plan`.
2.  **Test-First Implementation:** If testing is part of the project constraints, write the focused tests for the required behavior first.
3.  **Execute:** Write the code to check off every item in the `## Implementation Plan`.
4.  **Validate ACs:** You must verify that every single item in the `## Acceptance Criteria` has been met. A task is not done just because the code is written; it is done when the ACs are demonstrably true.
5.  **Run Quality Gates:** Run the project's build, lint, and test commands (e.g., `npm run build && npm test` or `go test ./...`).
6.  **Log Context:** Fill out the `## Context Log` at the bottom of the task file:
    *   List files read/edited.
    *   Estimate tokens used.
    *   Record the result of the quality gates.
7.  **Log Drift Notes:** **CRITICAL STEP.** Ask yourself:
    *   Did I add new source files, modules, or exports not listed in the `AGENTS.md` Codebase Map?
    *   Did I change the architecture from what `.savepoint/Design.md` describes?
    *   If YES to either, append a `## Drift Notes` section to the bottom of the task file explaining what changed. DO NOT edit `Design.md` or `AGENTS.md` yourself. The `audit` agent will handle that later.
8.  **Status Update:** Change the task frontmatter to `status: done`.
9.  **Handoff:** Update `router.md` to point to the next unblocked task. If all tasks in the Epic are done, update it to `state: audit-pending`.
10. **Stop:** Prompt the user: "Task {id} is done. Quality gates passed. Review the changes, then tell me to continue." Do not start the next task automatically.

## Constraints
- **Stay in scope:** Do not touch files outside of what is required for the Acceptance Criteria.
- **Do not edit architecture documents.** If you must deviate from the plan, write the code and log a "Drift Note" in the task file.