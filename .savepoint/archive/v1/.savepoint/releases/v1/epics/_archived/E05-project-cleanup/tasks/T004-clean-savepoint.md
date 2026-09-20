---
id: E05-project-cleanup/T004-clean-savepoint
status: planned
objective: "Delete .savepoint/audit/, .savepoint/metrics/, .savepoint/Design.md, .savepoint/PRD.md, .savepoint/visual-identity.md; trim config.yml"
depends_on: [E05-project-cleanup/T002-delete-dead-tests]
---

# T004: Clean Savepoint Data

## Acceptance Criteria

- `.savepoint/audit/` directory does not exist.
- `.savepoint/metrics/` directory does not exist.
- `.savepoint/Design.md` does not exist.
- `.savepoint/PRD.md` does not exist.
- `.savepoint/visual-identity.md` does not exist.
- `config.yml` contains only theme configuration.
- `router.md` still exists and is valid.

## Implementation Plan

- [ ] Delete `.savepoint/audit/` and all files within.
- [ ] Delete `.savepoint/metrics/` and all files within.
- [ ] Delete `.savepoint/Design.md`.
- [ ] Delete `.savepoint/PRD.md`.
- [ ] Delete `.savepoint/visual-identity.md`.
- [ ] Read existing `config.yml`.
- [ ] Rewrite `config.yml` with theme-only content.
