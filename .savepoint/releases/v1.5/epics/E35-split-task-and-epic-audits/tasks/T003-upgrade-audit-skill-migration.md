---
id: E35-split-task-and-epic-audits/T003-upgrade-audit-skill-migration
status: planned
objective: Migrate existing projects from the generic audit skill to the split skills without losing legacy content.
depends_on:
  - E35-split-task-and-epic-audits/T001-split-audit-contracts
  - E35-split-task-and-epic-audits/T002-generated-routing-and-guidance
complexity_tier: high
complexity_reason: Upgrade must relocate user content, remove an active alias, install new assets, and remain dry-run safe.
---

# T003: Upgrade Audit Skill Migration

## Problem

The current upgrade path copies only `SKILL.md` files and would leave an existing generic audit skill active while failing to install the shared method reference.

## Context Files

- `internal/init/upgrade.go`
- `internal/init/upgrade_test.go`
- `templates/project/agent-skills/savepoint-audit-task/SKILL.md`
- `templates/project/agent-skills/savepoint-audit-epic/SKILL.md`
- `templates/project/agent-skills/references/audit-method.md`

## Acceptance Criteria

- [ ] Upgrade-assets installs both split skill folders and the shared audit-method reference.
- [ ] If `agent-skills/savepoint-audit/SKILL.md` exists, upgrade preserves its content under a documented non-triggerable `.savepoint/migrations/` archive before removing the active generic folder.
- [ ] An upgraded project has no triggerable `agent-skills/savepoint-audit/SKILL.md` alias.
- [ ] Migration is idempotent and does not overwrite an existing archived legacy file without an explicit safe conflict policy.
- [ ] Dry-run reports the migration and split-asset actions without changing the filesystem.
- [ ] Upgrade reporting distinguishes migrated legacy content from updated, unchanged, merged, and skipped assets.
- [ ] Projects without the legacy skill upgrade successfully without creating an unnecessary archive.
- [ ] Upgrade tests cover stock, user-modified, missing, repeated, dry-run, and archive-conflict cases.

## Implementation Plan

- [ ] Extend upgrade asset classification to install the shared reference in addition to `SKILL.md` files.
- [ ] Add a preflighted legacy-skill migration into a non-triggerable Savepoint archive.
- [ ] Remove the active legacy file and empty folder only after its content is safely preserved.
- [ ] Add migration reporting and deterministic dry-run behavior.
- [ ] Add table-driven upgrade tests for legacy-content and idempotency branches.

## Context Log

Pending.
