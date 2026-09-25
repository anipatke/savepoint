---
id: v1.2/D008-filename-style-dependency-resolution
release: v1.2
status: resolved
severity: high
title: "Filename-style same-epic dependencies block task completion"
reference: E17-defect-workflow-tui/T010-defect-resolve-hotkey
---

# D008: Filename-Style Dependency Resolution

## Symptom

Moving T010 to done from the board fails with `dependency "T004-defects-overlay" not found`, even though the related E17 task file exists and is already done.

## Expected Behavior

The board and doctor should resolve same-epic dependency references consistently, including the filename-style `T###-slug` form used by task planning guidance.

## Reproduction

1. Open the board for release `v1.2`, epic `E17-defect-workflow-tui`.
2. Focus `T010-defect-resolve-hotkey`.
3. Attempt to advance it from audit to done.
4. Observe the dependency-not-found status message for `T004-defects-overlay`.

## Impact

Users cannot complete a valid task when its dependency uses a same-epic filename stem rather than a bare `T###` or full task ID. The failure is confusing because the referenced task is present.

## Fix Plan

- Update the shared dependency resolver to accept filename-style same-epic task references.
- Keep slash-qualified dependency references strict so invalid cross-epic full IDs do not silently resolve to same-epic tasks.
- Add board transition and doctor dependency tests for the filename-style form.
- Update task planning guidance and design documentation to describe supported dependency reference forms.

## Acceptance Criteria

- [x] `T004-defects-overlay` resolves to the same-epic task `E17-defect-workflow-tui/T004-defects-overlay`.
- [x] Board task advancement uses the fixed resolver.
- [x] Doctor dependency checks use the fixed resolver.
- [x] Slash-qualified full dependency references remain exact matches.
- [x] Planner guidance documents full IDs plus same-epic shorthand forms.

## Resolution Notes

Implemented in `internal/data.ResolveDependency` by normalizing same-epic, non-slash task references through their `T###` prefix. Added regression tests in `internal/board/transitions_test.go` and `internal/doctor/checks_test.go`. Updated live and scaffolded create-task guidance plus `.savepoint/Design.md`.
