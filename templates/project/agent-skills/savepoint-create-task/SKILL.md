---
name: savepoint-create-task
description: Plans Savepoint task files during epic-task-breakdown by writing acceptance criteria, implementation checklists, dependencies, and context-log shells.
---

# Savepoint Skill: Create Task (`create-task`)

## Objective
Take high-level tasks identified during the planning phase and build detailed, actionable task plans with strict Acceptance Criteria.

## Context
A task plan is a contract between the planner and the builder. If the task plan is vague, the resulting code will be buggy. The `create-task` skill acts as a Senior Engineer writing tickets for a Junior Developer (the `build-task` agent). It must define exactly *what* constitutes success and *how* to achieve it, without actually writing the code.

## Trigger
This skill is activated when the `.savepoint/router.md` state is `epic-task-breakdown` and the router points to a specific epic.

## Input
- `.savepoint/releases/v1/epics/{E##-epic}/E##-Detail.md` (The active Epic's design).
- The high-level task markdown file (e.g., `.savepoint/releases/v1/epics/{E##-epic}/tasks/T001-slug.md`).

## Workflow

1.  **Read Context:** Understand the goal of the specific task within the context of its parent Epic.
2.  **Define Acceptance Criteria (ACs):** This is the most critical step. Write explicit, observable outcomes. (e.g., "Running `npm test` passes 5 tests for the auth module" or "The `/login` route returns a 200 OK with a valid JWT"). Do not use subjective language like "looks good."
3.  **Draft Implementation Plan:** Create a step-by-step checklist for the `build-task` agent to follow.
    *   List which files need to be created or modified.
    *   Specify which functions or components need to be written.
    *   Include instructions to write tests *first* if applicable.
4.  **Populate Context Files:** Add a `## Context Files` section listing the exact file paths the build agent must read before coding. Pull these from the epic's Components table — only the files this task actually touches or depends on. Example:
    ```markdown
    ## Context Files
    - `internal/board/board.go`
    - `internal/board/model.go`
    - `internal/data/config.go`
    ```
    No globs. No directories. Exact paths only.
5.  **Add Context Log Shell:** Ensure the bottom of the task file includes a `## Context Log` section with placeholders for `Files read:`, `Estimated input tokens:`, and `Notes:`.
5.  **Define Dependencies:** If this task relies on another task being completed first, explicitly declare it in the YAML frontmatter (e.g., `depends_on: [T001-setup]`).
6.  **Status Update:** Change the task frontmatter to `status: planned`.
7.  **Handoff:** Update `.savepoint/router.md` to `state: in-progress` and ensure it points to the newly planned task. Prompt the user to approve the task plan before building begins.

## Constraints
- **Do not write code.** Your job is to plan the work, not execute it.
- **Keep it isolated:** The task plan should touch as few files as possible. If a task plan requires modifying 15 files, it is too large and should be broken down further.
