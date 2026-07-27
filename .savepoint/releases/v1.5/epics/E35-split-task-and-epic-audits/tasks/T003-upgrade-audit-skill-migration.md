---
id: E35-split-task-and-epic-audits/T003-upgrade-audit-skill-migration
status: done
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

- [x] Upgrade-assets installs both split skill folders and the shared audit-method reference.
- [x] If `agent-skills/savepoint-audit/SKILL.md` exists, upgrade preserves its content under a documented non-triggerable `.savepoint/migrations/` archive before removing the active generic folder.
- [x] An upgraded project has no triggerable `agent-skills/savepoint-audit/SKILL.md` alias.
- [x] Migration is idempotent and does not overwrite an existing archived legacy file without an explicit safe conflict policy.
- [x] Dry-run reports the migration and split-asset actions without changing the filesystem.
- [x] Upgrade reporting distinguishes migrated legacy content from updated, unchanged, merged, and skipped assets.
- [x] Projects without the legacy skill upgrade successfully without creating an unnecessary archive.
- [x] Upgrade tests cover stock, user-modified, missing, repeated, dry-run, and archive-conflict cases.

## Implementation Plan

- [x] Extend upgrade asset classification to install the shared reference in addition to `SKILL.md` files.
- [x] Add a preflighted legacy-skill migration into a non-triggerable Savepoint archive.
- [x] Remove the active legacy file and empty folder only after its content is safely preserved.
- [x] Add migration reporting and deterministic dry-run behavior.
- [x] Add table-driven upgrade tests for legacy-content and idempotency branches.

## Context Log

**Read:** `internal/init/upgrade.go`, `internal/init/upgrade_test.go`,
`internal/init/scaffold.go`, `internal/init/integration_test.go`, plus the three
T001 assets the upgrade must install.

**Added — `internal/init/migrate_audit_skill.go`:**

- `migrateLegacyAuditSkill(absTarget, dryRun)` returns `nil` when
  `agent-skills/savepoint-audit/SKILL.md` is absent, so a project that never had
  the legacy skill gains no archive and no report entry.
- Preserve-then-delete order: content is archived under `.savepoint/migrations/`
  before the triggerable file is removed, then the legacy folder is removed only
  when empty (a folder holding other user files is left alone).
- Archive path is `.savepoint/migrations/savepoint-audit-SKILL.md` — deliberately
  not named `SKILL.md` and not under `agent-skills/`, so no skill loader can
  discover it. A `README.md` is written alongside it, once, documenting what the
  archive is and that nothing in the directory is agent-loaded.
- Conflict policy in `resolveArchivePath`: an identical existing archive returns
  an empty path (nothing rewritten, migration is idempotent); a differing archive
  is never overwritten — the new copy goes to the next free numbered sibling
  (`savepoint-audit-SKILL.1.md`, …) up to 100, after which upgrade errors rather
  than guessing.

**Edited — `internal/init/upgrade.go`:**

- New `ActionMigrated` action plus a `Migrated:` line in `Format()`, so migration
  is reported distinctly from updated, merged, unchanged, and skipped.
- `UpgradeProjectAssets` runs the migration before the template walk, so an
  interrupted upgrade never leaves the old alias triggerable beside its
  replacements.
- `isPackageSkillAsset` replaces the inline `/SKILL.md` suffix test: package-owned
  skill assets are now skill entrypoints **and** anything under
  `agent-skills/references/`, which is how the shared audit method reaches
  existing projects. Dry-run, create-missing, and idempotency paths are shared
  with skills, so no branch was duplicated.

**Edited — `internal/init/upgrade_test.go`:** four fixtures used
`agent-skills/savepoint-audit/SKILL.md` as a generic stand-in skill; that name now
means "legacy alias", so they were retargeted to `savepoint-audit-epic`. Behavior
under test is unchanged.

**Added — `internal/init/migrate_audit_skill_test.go`:** split-skill and shared
reference installation, shared reference refresh, stock legacy migration
(archive + removal + README + empty-folder cleanup), user-modified content
preserved verbatim, missing legacy creates no archive, repeated upgrade
idempotent (archive dir stays at README + one copy), archive conflict writes a
numbered sibling without touching the original, dry-run reports `migrated` while
leaving both the legacy file and the archive untouched, and `Format()` counting.

**Verification:**

- `make build && make test` — pass, all packages.
- End-to-end with a real binary against a scaffolded scratch project seeded with
  a locally edited legacy skill:
  - dry run printed `Migrated: 1` and left `agent-skills/savepoint-audit/SKILL.md`
    in place;
  - real run archived the file verbatim to
    `.savepoint/migrations/savepoint-audit-SKILL.md`, wrote the README, removed
    the legacy folder, and left the split skills installed;
  - a third run reported zero migrations and the archive dir still held exactly
    two files;
  - deleting `agent-skills/references/` and re-running reported
    `updated agent-skills/references/audit-method.md` and reinstalled it.

  The first attempt at this check silently used a stale binary (the rebuild had
  failed from the wrong working directory) and appeared to show no migration;
  the result above is from a verified fresh build.

**Health Check:** skipped — this repo has no `.savepoint/Health-Check.md`; the
file exists only as `templates/project/.savepoint/Health-Check.md` for generated
projects.

## Drift Notes

- New module file `internal/init/migrate_audit_skill.go`. The `AGENTS.md` Codebase
  Map row for `internal/init/` reads "Target validation, scaffold writing from
  templates, managed AGENTS.md merge behavior, and safe project asset refresh" —
  it does not mention legacy-asset migration or the `.savepoint/migrations/`
  archive. T004 or a follow-up should extend that row; recorded here so the map
  gap is explicit.
- `.savepoint/migrations/` is a new generated-project directory that no template
  or Design.md directory listing describes. `templates/project/.savepoint/Design.md`
  lists `.savepoint/` contents; the archive only appears after an upgrade, so it
  is deliberately not scaffolded — but the listing may want a note.
- `internal/init/integration_test.go` still uses `agent-skills/savepoint-audit/SKILL.md`
  as an in-memory scaffold fixture. It is a scaffold test, not a migration test, so
  it passes unchanged, but T004 owns updating that fixture and adding upgrade
  integration coverage.
