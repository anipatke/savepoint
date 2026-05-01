---
id: E07-audit-pipeline/T005-proposal-validation-apply
status: done
objective: Validate audit proposal files and apply approved patch-shaped documentation updates with explicit divergence failures.
depends_on:
  - E07-audit-pipeline/T003-snapshot-and-prompt
---

# T005: proposal-validation-apply

## Implementation Plan

- [x] Add `src/audit/proposals.ts` to discover proposal files in `.savepoint/audit/{epic}/proposals/` and parse patch-shaped operations for documentation targets.
- [x] Validate proposal target paths, required sections, duplicate operations, empty replacements, and malformed proposal documents with path-aware errors.
- [x] Add `src/audit/apply-proposals.ts` to apply approved replace, insert, and delete operations only when anchors still match the live file.
- [x] Fail with explicit divergence errors instead of overwriting user edits when a proposal no longer matches the target document.
- [x] Preserve compatibility with any existing board handoff signal until the review UI moves fully to the proposal directory.
- [x] Add tests for discovery, validation failures, successful apply, partial approval, divergence, and scoped path protection.
- [x] Run the focused proposal validation and apply tests before handing off the task.

## Context Log

- Files read: `AGENTS.md`; `.savepoint/router.md`; `.savepoint/releases/v1/epics/E07-audit-pipeline/Design.md`; `.savepoint/releases/v1/epics/E07-audit-pipeline/tasks/T005-proposal-validation-apply.md`; `src/audit/snapshot.ts`; `src/audit/prompts.ts`; `src/audit/router-state.ts`; `src/audit/log.ts`; `src/fs/project.ts`; `test/audit/snapshot.test.ts`; `test/audit/quality-gates.test.ts`; `templates/prompts/audit-reconciliation.prompt.md`; `.savepoint/audit/E06-tui-board/proposals.md`; `.savepoint/audit/E05-init-command/proposals.md`
- Estimated input tokens: ~18,500
- Notes: Implemented `proposals.ts` (discovery, parse, validate) and `apply-proposals.ts` (applyOperation, applyApprovedOperations). All 46 focused tests pass. Proposal format supports Replace/Insert After/Delete with fenced or plain-text anchors; scoped path protection rejects non-.md/.yml/.yaml and path traversal targets; divergence detection skips writing any operation to a file if any op on it diverges.
