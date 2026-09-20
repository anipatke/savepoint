---
id: E02-data-model/T003-task-documents
status: done
objective: "Parse task markdown files into typed task documents with validated IDs, statuses, objectives, and dependency IDs."
depends_on:
  - E02-data-model/T001-domain-ids-status
  - E02-data-model/T002-markdown-frontmatter-boundary
---

# T003: Task Documents

## Scope

Build the task reader and task document model on top of the markdown/frontmatter boundary.

## Acceptance Criteria

- Task frontmatter maps into a typed task document.
- Invalid task IDs, statuses, objectives, or dependency IDs are rejected with path-aware errors.
- Markdown body content remains available to callers.
- Valid and invalid task document cases have focused unit tests.

## Implementation Plan

- [x] Reuse the T001 ID/status validators and the T002 markdown/frontmatter boundary as the only parsing inputs.
- [x] Define the typed task frontmatter and task document model, including validated dependency IDs and markdown body.
- [x] Implement task document parsing that rejects invalid IDs, statuses, missing objectives, and malformed `depends_on` values with path-aware errors.
- [x] Keep task parsing read-only and separate pure frontmatter validation from filesystem access.
- [x] Add focused unit tests for valid task files and each invalid frontmatter branch.
