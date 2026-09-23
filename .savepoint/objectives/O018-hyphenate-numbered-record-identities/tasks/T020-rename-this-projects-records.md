---
id: T020
title: Rename this project's records
objective: O018
status: planned
depends_on: [{task: T019, requires: clear}]
owner_validation: {required: true}
planned_by: {role: planner, session: o018-replan-20260923}
---

# T020: Rename this project's records

## Outcome

This repository's active records and guidance use hyphenated identities, via
a rename script kept in the repo for other projects.

## User Check

Review the board and resume output after the rename and confirm it still
points at the intended O-018 work. The optional Task Check may be requested
or explicitly waived; the Full Objective Check remains mandatory.

## Done When

The script renames directories/files under `.savepoint/{releases,objectives,checks,issues}`
and rewrites identity references in frontmatter, bodies, `router.md`,
`AGENTS.md`, `Design.md`, and `Guardrails.md`, leaving `.savepoint/archive/`
untouched. Strict load, board, doctor, and resume succeed; a targeted search
finds no unhyphenated active IDs; `git diff --check` and fresh
`make test-full` pass.

## Context Files

`.savepoint/router.md`, `AGENTS.md`, `.savepoint/Design.md`,
`.savepoint/Guardrails.md`, `.savepoint/objectives/`, `.savepoint/issues/`,
`.savepoint/checks/`, `.savepoint/releases/`.

## Guardrails

FS-01..04, DATA-01..04, TEST-01..05.

## Implementation Plan

1. Write the script under `scripts/`; match word-bounded `[ROTCI][0-9]{3,}`
   only, so words and version strings are untouched.
2. Run it, review the diff, and fix any false positives by hand.
3. Verify load, board, doctor, resume, and the full gate.

## Boundaries

No migration command, recovery state, or archived Check copies. Check
content changes only in identity references.

## Technical Verification

Fresh `make test-full` at handoff.

## Technical Evidence

Pending execution.
