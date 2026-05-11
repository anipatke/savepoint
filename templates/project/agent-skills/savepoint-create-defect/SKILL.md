---
name: savepoint-create-defect
description: Captures release-level Savepoint defect files when a user reports a concrete bug, regression, or broken expectation that should be tracked separately from planned epic/task work.
---

# Savepoint Skill: Create Defect (`create-defect`)

## Objective
Turn a bug report into a focused release-level defect record without starting repair work by default.

## Context
A defect is not an epic and not normal task rework. It records observed incorrect behavior so the project can preserve evidence, reproduce the issue, and later repair it through `defect-building`.

Use this skill when the user reports a concrete bug, regression, failed expectation, or "this is broken" conversation and wants it captured in Savepoint.

## Trigger
This skill is activated by an explicit defect-capture request, or when the user reports a bug that should be tracked separately from the current planned task flow.

## Input
- `.savepoint/router.md` for the active release when the user does not name one.
- Existing `.savepoint/releases/{release}/defects/*.md` files, if any, to choose the next `D###` id.
- User-provided symptom, expected behavior, reproduction steps, and impact.

## Workflow

1. **Classify the report.** Confirm it is a defect:
   - Use a defect for a concrete bug, regression, broken behavior, failed test, or production issue.
   - Update the active task instead when planned work simply needs rework before completion.
   - Use an epic/task when the request is new scope or an enhancement.
2. **Choose release and ID.** Use the user's named release, otherwise use `router.md` `release`. Create `.savepoint/releases/{release}/defects/` if needed. Pick the next unused `D###-slug.md`.
3. **Write the defect file.** Use this structure:
   ```markdown
   ---
   id: {release}/D###-slug
   release: {release}
   status: open
   severity: medium
   title: "One-line bug title"
   reference: E##-epic/T###-task # optional
   ---

   # D###: Title

   ## Symptom

   ## Expected Behavior

   ## Reproduction

   ## Impact

   ## Fix Plan

   ## Acceptance Criteria

   - [ ] Observable repair outcome

   ## Resolution Notes

   Pending.
   ```
4. **Keep repair separate.** Do not edit product code in this skill. If the user wants to start fixing immediately, ask or explicitly transition the router to `defect-building` for the new defect.
5. **Handoff.** Tell the user the defect path and whether the router was left unchanged or moved to `defect-building`.

## Constraints
- **Do not create an epic or task for the defect itself.** Defects are release-level records.
- **Do not set `status: in_progress` unless repair is starting now.** If status becomes `in_progress`, include `stage: build`.
- **Do not run broad code searches.** Capture the report from the user's evidence and named files only.
- **Do not mark defects `resolved` prematurely.** Resolution happens after repair and verification.
- **Defect lifecycle:** `open` → `in_progress` (with `stage: build|test|audit`) → `resolved`. Never use task-style `planned` or `done`.
