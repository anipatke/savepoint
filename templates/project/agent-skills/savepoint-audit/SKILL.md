---
name: savepoint-audit
description: Performs Savepoint audit-pending work for a completed epic, reviewing implementation against task acceptance criteria and writing the required E##-Audit.md handoff file.
---

# Savepoint Skill: Audit

## Purpose

Review a completed epic with fresh eyes, verify task acceptance criteria, capture drift, and write one audit handoff file.

## Trigger

Use this skill when router `state` is `audit-pending`, or when the user explicitly asks for an epic audit.

## Read

- `.savepoint/router.md`
- Active epic detail file
- Task files for the active epic
- `.savepoint/Design.md`
- `AGENTS.md`
- Scoped source/test files changed by the epic

## Workflow

1. Stop if you are the same agent session that just built the epic.
2. Verify every completed task against its acceptance criteria and context log.
3. Review task `## Drift Notes`.
4. Write exactly one `.savepoint/releases/{release}/epics/{E##-slug}/E##-Audit.md`.
5. Put user-facing findings under `## Main Findings` and code-style checks under `## Code Style Review`.
6. Put admin replacement blocks only under `## Proposed Changes`.
7. Stop and ask the user to review the audit before any proposals are applied.

## Artifact Template

Write `.savepoint/releases/{release}/epics/{E##-slug}/E##-Audit.md` with this structure:

````markdown
---
type: audit-findings
audited: {date}
---

# Audit Findings: E## Epic Title

## Main Findings

User-facing audit summary, including concrete verification outcomes and any unresolved findings.

## Code Style Review

- [ ] One job per file
- [ ] One job per function
- [ ] Test branches
- [ ] Types document intent
- [ ] Build only what is needed
- [ ] Handle errors at boundaries
- [ ] One source of truth
- [ ] Comments explain WHY
- [ ] Content in data files
- [ ] Small diffs

## Proposed Changes

### Target File
path/to/file

### Replace
```md
Existing text to replace.
```

### With
```md
Replacement text.
```
````

## Apply + Close

Only after the user says to apply the audit:

1. Apply approved `## Proposed Changes` blocks.
2. Update `E##-Audit.md` visible sections to describe the applied outcome.
3. Mark the epic audited and advance the router.
4. Report: "Updated audit findings."

## Rules

- Do not write product code during audit.
- Do not apply proposals before approval.
- Do not create more than one audit file for an epic.
- Use `state` only for router phase, task `status` only for task lifecycle, and `stage` only when an item is `in_progress`.
