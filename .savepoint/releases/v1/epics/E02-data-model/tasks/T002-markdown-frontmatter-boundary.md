---
id: E02-data-model/T002-markdown-frontmatter-boundary
status: done
objective: "Add scoped markdown/frontmatter read helpers that return typed boundary errors with file path and reason."
depends_on:
  - E02-data-model/T001-domain-ids-status
---

# T002: Markdown Frontmatter Boundary

## Scope

Create the filesystem boundary for reading markdown documents with structured frontmatter parsing.

## Acceptance Criteria

- Markdown documents can be read from scoped project paths.
- Frontmatter is parsed with a structured YAML/frontmatter parser.
- Missing, malformed, or invalid frontmatter returns errors that include path and reason.
- Parser success and failure branches have focused unit tests.

## Implementation Plan

- [x] Review the T001 domain result/error style and existing test fixture conventions.
- [x] Add a scoped project path helper for locating `.savepoint` documents without exposing arbitrary filesystem reads.
- [x] Add a markdown/frontmatter reader that uses a structured parser and preserves markdown body content.
- [x] Define path-aware boundary errors for missing files, missing frontmatter, malformed frontmatter, and schema rejection.
- [x] Add focused unit tests for successful reads and each parser/boundary failure branch using small temporary fixtures.
