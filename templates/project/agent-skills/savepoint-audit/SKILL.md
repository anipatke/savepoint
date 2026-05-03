---
name: savepoint-audit
description: Performs Savepoint audit-pending work for a completed epic, reviewing implementation against task acceptance criteria and writing the required E##-Audit.md handoff file.
---

# Savepoint Skill: Audit (`audit`)

## Objective
At the end of an epic, review the implemented code against the Epic objectives and Task Acceptance Criteria, and reconcile any documentation "Drift".

## Context
The Audit Gate is Savepoint's wedge. It prevents projects from degrading into chaos. The `audit` agent acts as the Quality Assurance and Documentation Lead. It reviews the work of the `build-task` agent with "fresh eyes," ensuring that the `Design.md` and the `AGENTS.md` Codebase Map reflect the actual reality of the codebase before the next Epic begins.

## Trigger
This skill is activated when the `.savepoint/router.md` state is `audit-pending`.

## Input
- `.savepoint/releases/{release}/epics/{E##-slug}/E##-Detail.md` (the epic design).
- `.savepoint/Design.md` (project architecture).
- `AGENTS.md` (agent guide and codebase map).
- The source and test files modified during the Epic.

## Workflow

1.  **Fresh Eyes Check:** If you are the exact same agent session that just built the `build-task` code for this Epic, you MUST STOP. Tell the user: "Epic complete. Start a new agent session for the audit."
2.  **Verify ACs:** Review the completed tasks for the Epic. Ensure the Acceptance Criteria were actually met by the committed code. Also confirm each task file has a `## Context Files` section — flag any missing ones under `## Main Findings` as a process gap.
3.  **Process Drift Notes (Reconciliation):** Read every task file in the Epic and look for `## Drift Notes`.
4.  **Draft Audit File:** Based on the code changes and the Drift Notes, write exactly ONE file: `.savepoint/releases/{release}/epics/{E##-slug}/E##-Audit.md`. It must use this format:
    ````md
    ---
    type: audit-findings
    audited: YYYY-MM-DD
    ---
    # Audit Findings: E## {Epic Name}

    ## Main Findings
    [human-readable audit summary only: AC verification, important drift, and notable risks. Do not include per-file replacement blocks here.]

    ## Code Style Review
    - [ ] One job per file
    - [ ] One-sentence functions
    - [ ] Test branches
    - [ ] Types are documentation
    - [ ] Build, don't speculate
    - [ ] Errors at boundaries
    - [ ] One source of truth
    - [ ] Comments explain WHY
    - [ ] Content in data files
    - [ ] Small diffs

    ## Proposed Changes
    ### Target File
    path/to/file.md

    ### Replace
    ```
    exact old text
    ```

    ### With
    ```
    replacement text
    ```
    ````
    The TUI Audit tab renders only `## Main Findings` and `## Code Style Review`. Keep `## Proposed Changes` as admin/apply metadata so the Epic Detail panel does not show stale file-change blocks.
5.  **Review Format:** Use `### Target File`, `### Replace`, and `### With` formatting only inside `## Proposed Changes`. Include proposals for:
    *   **Design.md:** Merge the epic's architectural changes into the project-level `Design.md`.
    *   **AGENTS.md:** Refresh the Codebase Map table with new or changed modules.
    *   **Epic E##-Detail.md:** Add "Implemented as:" notes showing where reality deviated from the original plan.
    *   **Implementation fixes:** Include any must-fix code or test changes found during audit.
6.  **Handoff:** Do not apply the proposals yourself. Do not mark the epic audited. Stop and prompt the user to review `.savepoint/releases/{release}/epics/{E##-slug}/E##-Audit.md` (via the TUI Audit tab). Tell them: "Review the audit tab. When ready, say 'apply audit' to apply proposals and close the epic."
7.  **Apply + Close** (only when the user approves by saying "apply audit" or equivalent):
    1.  Read `E##-Audit.md` — extract every `### Target File` / `### Replace` / `### With` block from `## Proposed Changes`.
    2.  For each pair, apply the replacement to the target file named in the preceding `### Target File` line.
    3.  Update `E##-Audit.md` visible sections: rewrite `## Main Findings` and `## Code Style Review` so the TUI Audit tab reflects the applied outcome, resolved findings, remaining risks if any, and final code-style status. Keep `## Proposed Changes` intact as admin/apply trace unless the user asks to remove it.
    4.  Update `E##-Detail.md` frontmatter: set `status: audited`.
    5.  Update `.savepoint/Design.md` frontmatter: set `last_audited: {release}/{E##-slug}`.
    6.  Read `.savepoint/router.md` current state. Advance:
        - If more epics remain in the release: set `state: epic-design`, `epic: {next-epic-slug}`, `task: ""`, `next_action: "Draft epic design"`.
        - If no more epics: set `state: epic-design`, `epic: ""`, `next_action: "Plan next epic"`.
    7.  Print apply summary: "Applied X proposals. Updated audit findings. Epic {E##} closed as audited. Router → {new state}."

## Constraints
- **Do not write product code.** You are an auditor.
- **Do not apply the changes immediately.** Write the proposals document first.
- **One proposals file.** Do not create multiple proposal files.
- **No CLI audit pipeline.** Savepoint audit is agent-led and skill-driven; do not invoke or design around a `savepoint audit` command.
