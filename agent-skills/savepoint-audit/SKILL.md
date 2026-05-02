# Savepoint Skill: Audit (`audit`)

## Objective
At the end of an epic, review the implemented code against the Epic objectives and Task Acceptance Criteria, and reconcile any documentation "Drift".

## Context
The Audit Gate is Savepoint's wedge. It prevents projects from degrading into chaos. The `audit` agent acts as the Quality Assurance and Documentation Lead. It reviews the work of the `build-task` agent with "fresh eyes," ensuring that the `Design.md` and the `AGENTS.md` Codebase Map reflect the actual reality of the codebase before the next Epic begins.

## Trigger
This skill is activated when the `.savepoint/router.md` state is `audit-pending`.

## Input
- `.savepoint/audit/{release}/{E##-slug}/snapshot.md` (what changed).
- `.savepoint/releases/{release}/epics/{E##-slug}/E##-Detail.md` (the epic design).
- `.savepoint/Design.md` (project architecture).
- `AGENTS.md` (agent guide and codebase map).
- The source and test files modified during the Epic.

## Workflow

1.  **Fresh Eyes Check:** If you are the exact same agent session that just built the `build-task` code for this Epic, you MUST STOP. Tell the user: "Epic complete. Start a new agent session for the audit."
2.  **Verify ACs:** Review the completed tasks for the Epic. Ensure the Acceptance Criteria were actually met by the committed code.
3.  **Process Drift Notes (Reconciliation):** Read every task file in the Epic and look for `## Drift Notes`.
4.  **Draft Proposals:** Based on the code changes and the Drift Notes, write exactly ONE file: `.savepoint/audit/{release}/{E##-slug}/proposals.md`. It must contain:
    *   **Design.md section:** Propose updates to merge the epic's architectural changes into the project-level `Design.md`.
    *   **AGENTS.md section:** Propose updates to refresh the Codebase Map table with new or changed modules.
    *   **Epic-E##-Detail.md section:** Add "Implemented as:" notes showing where reality deviated from the original plan.
    *   **Quality Review section:** List any minor code-style infractions that must be fixed before the next Epic.
5.  **Review Format:** Use `## Target File`, `## Replace`, and `## With` formatting in the proposals document so it is easy for a human (or an agent) to apply later.
6.  **Handoff:** Do not apply the proposals yourself. Do not mark the epic audited. Stop and prompt the user to review `.savepoint/audit/{release}/{E##-slug}/proposals.md` (often via the TUI). Once the user approves, the proposals are applied, and the router moves to the next Epic.

## Constraints
- **Do not write product code.** You are an auditor.
- **Do not apply the changes immediately.** Write the proposals document first.
- **One proposals file.** Do not create multiple proposal files.
