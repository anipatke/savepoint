---
id: E02-data-model/T006-epic-task-set-reader
status: done
objective: "Read an epic task directory into a validated task set using the task reader and dependency validation."
depends_on:
  - E02-data-model/T003-task-documents
  - E02-data-model/T005-dependency-validation
---

# T006: Epic Task Set Reader

## Scope

Connect the file-backed task reader to dependency validation for an entire epic's task directory.

## Acceptance Criteria

- An epic task directory is read into a typed task set.
- Dependency validation runs after all task files are parsed.
- File boundary errors and graph validation errors remain distinguishable.
- Valid task sets and invalid graph cases have focused unit tests with small fixtures.

## Implementation Plan

- [x] Reuse the T003 task reader and T005 dependency validation without duplicating parsing or graph logic.
- [x] Add an epic task directory reader that discovers task markdown files and parses each into typed task documents.
- [x] Run dependency validation only after all readable task files have been parsed.
- [x] Return distinguishable results for file boundary failures versus dependency graph validation failures.
- [x] Add focused fixture tests for a valid task set, parse failures, missing dependencies, duplicate IDs, and cycles.
