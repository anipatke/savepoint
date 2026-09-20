---
id: E02-data-model/T004-release-epic-router-config-readers
status: done
objective: "Add read-only models and readers for release PRDs, epic designs, router state, and project configuration."
depends_on:
  - E02-data-model/T001-domain-ids-status
  - E02-data-model/T002-markdown-frontmatter-boundary
---

# T004: Release, Epic, Router, and Config Readers

## Scope

Expose typed read-only access for the remaining Savepoint planning documents and config state.

## Acceptance Criteria

- Release PRD, epic Design, router state, and config data have typed models.
- Readers keep filesystem access at the boundary.
- Config defaults are applied from one place.
- Missing or malformed input returns path-aware boundary errors.
- Reader branches and config default behavior have focused unit tests.

## Implementation Plan

- [x] Reuse the T001 domain primitives and T002 markdown/frontmatter boundary for all document reads.
- [x] Define typed models for release PRD metadata, epic Design metadata, router state, and project config defaults.
- [x] Implement read-only readers that keep path resolution and file IO at the boundary.
- [x] Apply config defaults from one exported source of truth before returning typed config data.
- [x] Add focused unit tests for valid reads, malformed or missing input, router-state validation, and config default behavior.
