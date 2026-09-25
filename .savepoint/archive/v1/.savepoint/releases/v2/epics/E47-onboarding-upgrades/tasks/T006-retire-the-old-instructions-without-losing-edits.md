---
id: E47-onboarding-upgrades/T006-retire-the-old-instructions-without-losing-edits
title: Retire the old instructions without losing edits
status: done
objective: Remove the nine V1 skills from a project that has become V2, archiving each one before deletion and reconciling the provenance manifest.
depends_on:
    - E47-onboarding-upgrades/T005-refresh-assets-that-match-the-projects-own-version
complexity_tier: medium
complexity_reason: Deletes user-visible files, so ordering, archival, idempotence, and dry-run behavior all have to hold together.
---

# T006: Retire the old instructions without losing edits

## Problem

A migrated project keeps its nine V1 skills on disk. `migrate` archives them as part of its inventory and deliberately does not touch `agent-skills/` beyond that; T005 installs the V2 skills beside them. The result is a project carrying two routing vocabularies at once, where `savepoint-build-task` and `savepoint-task` are both triggerable and an agent picks whichever it reads first.

Retirement removes the nine, and it is the one operation in this epic that deletes files a user can see. The order is not negotiable and is already established by `migrate_audit_skill.go`, which retired the generic audit skill the same way: archive the content under `.savepoint/migrations/`, confirm it is on disk, then remove the triggerable copy. A crash between those steps leaves a recoverable archive and a stale skill, never a deleted skill and no archive.

Every retired skill is archived, not only the edited ones. The manifest can prove which copies were customized, but choosing to discard the rest means deciding on the user's behalf which of their instructions were worth keeping, and an archive of nine short markdown files costs nothing against that. A differing archive already present gets a numbered sibling rather than an overwrite, exactly as the existing conflict policy does.

The trigger is the project being V2 and asking for a refresh — never the package version alone. A V1 project keeps all nine, because they are the instructions its router still points at.

The manifest is the loose end. It records a hash per shipped skill; a retired skill whose entry survives leaves the manifest claiming provenance for a file that no longer exists, and the next upgrade reasoning from a record of something absent. Entries are dropped as their files are, after the removal succeeds.

## Context Files

- `internal/init/migrate_audit_skill.go`
- `internal/init/migrate_audit_skill_test.go`
- `internal/init/upgrade.go`
- `internal/init/upgrade_test.go`
- `internal/init/manifest.go`
- `internal/init/manifest_test.go`
- `internal/init/lifecycle_test.go`
- `internal/init/write.go`
- `templates/project/agent-skills/savepoint-build-task/SKILL.md`
- `AGENTS.md`
- `.savepoint/Guardrails.md`

## Acceptance Criteria

- [x] Upgrading a V2 project removes all nine V1 skill directories — `savepoint-draft-prd`, `savepoint-create-plan`, `savepoint-system-design`, `savepoint-create-task`, `savepoint-build-task`, `savepoint-audit-task`, `savepoint-audit-epic`, `savepoint-audit-register`, `savepoint-create-defect` — and `references/audit-method.md`.
- [x] Each removed file has an archived copy under `.savepoint/migrations/` with byte-identical content, written and verified before its original is removed.
- [x] A user-edited V1 skill is archived with its edits intact; an unmodified one is archived too.
- [x] An existing archive with identical content is reused rather than duplicated; an existing archive with differing content gets a numbered sibling and is never overwritten.
- [x] The `.savepoint/migrations/README.md` explanation covers the retired skills and states that nothing in that directory is loaded or triggerable.
- [x] Manifest entries for the retired skills are removed after their files are, and the saved manifest records no path that does not exist in the project.
- [x] Upgrading a project that is still at schema version 1 retires nothing; all nine skills and the audit method remain byte-identical.
- [x] A V2 project that has already been retired upgrades again with no archive written, no file removed, no manifest change, and no mtime change.
- [x] `--dry-run` reports every retirement it would perform and writes nothing, verified by snapshotting bytes and mtimes.
- [x] A write failure during archival leaves the original skill in place and reports the failure with the path named; no skill is ever removed without a verified archive.
- [x] An `agent-skills/{skill}/` directory left non-empty by user files is not removed, matching the existing legacy-skill behavior.
- [x] The retirement report entries are distinguishable from installs and updates in the formatted output.
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Read `migrate_audit_skill.go` end to end and extract the reusable parts of its archive-conflict and preserve-then-delete policy rather than restating them.
- [x] Add `internal/init/retire_v1_skills.go` holding the retired-asset list and the retirement pass over one project.
- [x] Call it from the V2 branch of the upgrade, before the V2 assets are installed, so an interrupted run never leaves both vocabularies triggerable.
- [x] Extend the migrations README text to name the retired skills.
- [x] Add manifest entry removal, applied only after a successful file removal.
- [x] Add the retirement report entry and its formatted-output counting.
- [x] Write the tests: full retirement, edited-skill preservation, identical and differing archive conflicts, V1 project untouched, idempotent second run, dry run, injected archive-write failure, and non-empty directory retention.
- [x] Run `go test ./internal/init/...`, then `make build && make test`.

## Context Log

**Files read:** `migrate_audit_skill.go`, `migrate_audit_skill_test.go`, `upgrade.go`, `upgrade_test.go`, `upgrade_schema_test.go`, `upgrade_failure_test.go`, `manifest.go`, `manifest_test.go`, `lifecycle_test.go`, `write.go`, `AGENTS.md`, `.savepoint/Guardrails.md`, `.savepoint/router.md`, `E47-Detail.md`.

**Files edited:**
- `internal/init/migrate_audit_skill.go` — generalized `resolveArchivePath` and `writeArchive` to take a `stem` and an `assetWriter` instead of the hardcoded legacy-audit constant and `AtomicWrite`, so retirement reuses them exactly; extended `migrationsReadme` to name the nine retired V1 skills and the audit-method reference; `migrateLegacyAuditSkill` now takes `write assetWriter`.
- `internal/init/retire_v1_skills.go` (new) — `retiredV1SkillDirs`, `retiredV1AssetPaths`, `retireV1Skills`/`retireV1Asset` (archive-then-delete, manifest drop after removal), `archiveStem`.
- `internal/init/upgrade.go` — added `ActionRetired` and its `Format()` counting/line; `upgradeProjectAssets` gained a `retireV1 bool` parameter, called right after the legacy audit-skill migration and before the template walk; `UpgradeProjectAssets` now computes `isV2` once and passes it through as `retireV1`, and reuses it for the migrate-route info note.
- `internal/init/manifest.go` — added `Manifest.Forget(path)`.
- `internal/init/upgrade_failure_test.go` — updated the three direct `upgradeProjectAssets` call sites for the new trailing `retireV1` argument (`false`, unrelated to this task's behavior).
- `internal/init/retire_v1_skills_test.go` (new) — full retirement, edited-vs-unmodified skill preservation, identical/differing archive conflict, V1 project untouched, idempotent second run, dry run (bytes + mtimes), injected archive-write failure, non-empty directory retention, migrations README content, and report-format distinction.

**Quality gates:** `go build ./...`, `go vet ./...`, `gofmt -l` (clean after formatting `upgrade.go`), `go test ./internal/init/...`, `go test ./...`, `make build && make test` — all pass. No `.savepoint/Health-Check.md` in this project, so the Quick check step is skipped per the build-task skill.

**Post-audit remediation:** `retireV1Skills` now forgets manifest provenance after a confirmed missing retired file as well as after removal. Retirement entries carry the exact selected archive path, including numbered conflicts, and `writeArchive` upgrades only the byte-identical pre-E47 stock migration README on the V2 path; edited README content remains untouched. The focused regression cases are in `internal/init/retire_v1_skills_test.go` and the frozen V1 compatibility case is in `internal/init/upgrade_schema_test.go`.

No drift: no new files/modules outside `internal/init/retire_v1_skills.go` (already named in the epic's Components table) and no architecture change beyond what E47-Detail.md already describes.
