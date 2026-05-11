---
id: v1.2/D007-defect-status-language-mismatch
release: v1.2
status: resolved
severity: medium
title: "Defect status language mismatches Defects view"
reference: E17-defect-workflow-tui
---

# D007: Defect Status Language Mismatch

## Symptom

The Defects view presents defect status language as "Open" and "Resolved", but defect frontmatter and existing E17 defect records use task-style statuses such as `planned` and `done`.

## Expected Behavior

Defect files and the Defects view should use the same defect-specific language: `open` for unresolved defects and `resolved` for closed defects.

## Reproduction

1. Open the Defects view for release `v1.2`.
2. Inspect existing E17 defect files under `.savepoint/releases/v1.2/defects/`.
3. Compare the view labels "Open" and "Resolved" with file frontmatter values such as `planned` and `done`.

## Impact

Status mismatch makes defect lifecycle behavior ambiguous and can cause confusion when resolving defects from the TUI or validating defect files.

## Fix Plan

- Update the defect lifecycle model to use `open`, `in_progress`, and `resolved` where applicable.
- Ensure the Defects view labels and file frontmatter map to one consistent source of truth.
- Update existing E17 defect records from `planned`/`done` to `open`/`resolved`.
- Update tests and templates that still expect task-style `planned`/`done` values for defects.

## Acceptance Criteria

- [x] Defect files use `open` for unresolved defects.
- [x] Defect files use `resolved` for resolved defects.
- [x] The Defects view displays statuses from the same defect lifecycle semantics as defect frontmatter.
- [x] Existing E17 defect records no longer use `planned` or `done`.
- [x] Doctor validation accepts the canonical defect statuses and rejects task-only defect statuses where appropriate.
- [x] Defect templates and skill guidance use `open`/`resolved` language.

## Resolution Notes

Introduced `DefectStatus` type (`open`, `in_progress`, `resolved`) in `internal/data/defect.go` and `internal/data/parser.go`. Updated `validateDefectLifecycle` and the board overlay, update, and view handlers. All existing defect `.md` files migrated from `planned`/`done` to `open`/`resolved`. Skill guidance and AGENTS.md templates updated. All tests passing.
