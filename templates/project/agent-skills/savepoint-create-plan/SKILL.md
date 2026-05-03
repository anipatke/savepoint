---
name: savepoint-create-plan
description: Creates Savepoint release epics from the PRD and Design documents when the router is in pre-implementation planning.
---

# Savepoint Skill: Create Plan (`create-plan`)

## Objective
Turn the Product Requirements Document (PRD) and the architectural `Design.md` into Epics (logical slices of work) and high-level tasks that can be implemented independently.

## Context
Savepoint enforces small scopes to prevent AI agents from "wandering." The `create-plan` skill acts as the Technical Project Manager. It takes the grand vision and the architectural blueprint and slices it into a sequence of achievable, testable milestones (Epics).

## Trigger
This skill is activated when the `.savepoint/router.md` state is `pre-implementation`.

## Input
- `.savepoint/PRD.md` (Vision and constraints)
- `.savepoint/Design.md` (Architecture and layout)
- `.savepoint/releases/v1/v1-PRD.md` (Optional: Release-scoped PRD)

## Workflow

1.  **Read the Context:** Consume the PRD and Design documents to understand the scope and technical constraints.
2.  **Define Epics:** Group the work into high-level features or milestones. Name them clearly (e.g., `E01-scaffolding`, `E02-database`, `E03-auth`). Each Epic must represent a deliverable slice of value.
3.  **Draft Epic Designs:** For each Epic, create a shell `E##-Detail.md` inside `.savepoint/releases/v1/epics/{E##-epic-name}/E##-Detail.md`. This file should describe the *delta* (what this specific epic adds to the overall architecture).
4.  **Breakdown Tasks (High Level):** Inside each Epic folder, list out the high-level tasks required to complete it. Do not write full implementation plans yet—just identify the discrete chunks of work (e.g., `T001-setup-repo.md`, `T002-init-db.md`).
5.  **Order and Dependency:** Ensure the Epics and tasks are ordered logically. You cannot build the frontend auth UI (E03) before the database models (E02) exist.
6.  **Handoff:** Update `.savepoint/router.md` to `state: task-breakdown` and point it at the first Epic. Prompt the user to review the Epic list.

## Constraints
- **Do not write code.**
- **Do not write detailed implementation steps.** That is the job of the `create-task` skill. Keep the task outlines high-level at this stage.
