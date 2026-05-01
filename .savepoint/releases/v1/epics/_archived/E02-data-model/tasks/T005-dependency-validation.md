---
id: E02-data-model/T005-dependency-validation
status: done
objective: "Validate task dependency graphs for missing dependencies, duplicate IDs, and cycles using pure domain functions."
depends_on:
  - E02-data-model/T003-task-documents
---

# T005: Dependency Validation

## Scope

Implement dependency resolution and graph validation over parsed task documents.

## Acceptance Criteria

- Missing dependency IDs are reported with the dependent task ID.
- Duplicate task IDs are detected.
- Dependency cycles are detected and reported with enough IDs to explain the cycle.
- Validation functions are pure and do not read files.
- Missing dependency, duplicate ID, and cycle branches have focused unit tests.

## Implementation Plan

- [x] Build on the T003 task document model and keep validation functions independent of filesystem access.
- [x] Implement duplicate task ID detection over a task document collection.
- [x] Implement missing dependency detection that reports the dependent task ID and unresolved dependency ID.
- [x] Implement dependency cycle detection that returns enough task IDs to explain the cycle.
- [x] Add focused unit tests for valid graphs, missing dependencies, duplicate IDs, simple cycles, and longer cycles.
