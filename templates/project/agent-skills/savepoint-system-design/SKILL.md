---
name: savepoint-system-design
description: Produces Savepoint system or epic design documents from the PRD when the router is in epic-design or the user asks for architectural design.
---

# Savepoint Skill: System Design (`system-design`)

## Objective
Translate the Product Requirements Document (PRD) into the initial architectural blueprint (`Design.md`) before the planning phase breaks the work down into Epics.

## Context
Before any tasks can be planned, the project needs a high-level technical direction. The `system-design` skill acts as the Staff Engineer. It reads the vision and constraints from the PRD and makes authoritative technical decisions regarding architecture, directory layout, dependency strategies, and workflow rules.

## Trigger
This skill is activated when the `.savepoint/router.md` state is `epic-design` or when the user explicitly asks to design the system based on the PRD.

## Input
- `.savepoint/PRD.md` (The source of truth for "what" and "why").
- Any existing template or shell in `.savepoint/Design.md`.

## Workflow

1.  **Read the PRD:** Analyze `.savepoint/PRD.md` to understand the goal, constraints (e.g., tech stack, budget), and what is out of scope.
2.  **Architectural Mapping:** Make firm, opinionated technical decisions that solve the PRD's goals. Do not offer options unless absolutely necessary; pick a path and justify it briefly.
3.  **Draft `.savepoint/Design.md`:** Write or update the architectural design document. It must include:
    *   **Architecture Model:** The core pattern (e.g., "File-only state machine," "Client-side React with Firebase," "CLI Tool").
    *   **Directory Layout:** A clear, visual tree mapping out the structure of the project (`src/`, `test/`, etc.).
    *   **Hierarchy Semantics:** Definitions of what constitutes an Epic vs. a Task.
    *   **Status Model & Gates:** How work moves from `planned` to `done`.
    *   **Dependencies:** How modules or tasks rely on each other.
    *   **CLI Surface/API Contract:** If applicable, define the boundary interactions.
4.  **Review Against Constraints:** Ensure your design strictly adheres to the constraints and out-of-scope items listed in the PRD. If the PRD says "no database," do not add a database to the architecture.
5.  **Handoff:** Update `.savepoint/router.md` to `state: planning`. Instruct the user to review the `Design.md` and then prompt the agent to start the epic breakdown.

## Constraints
- **Do not write code.**
- **Do not break down tasks.** That is the job of the `create-plan` skill.
- **Maintain Token Discipline:** Keep the `Design.md` concise. It is meant to be read by AI agents on every turn. Rambling design docs destroy context windows.
