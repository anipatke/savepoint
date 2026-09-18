---
id: E45-safe-migration/T009-publish-the-conversion-once
title: Publish the conversion once
status: done
objective: Apply the plan in a fixed order ending with schema activation, and make rerun resume, conflict, or no-op correctly.
depends_on:
    - E45-safe-migration/T004-carry-active-work-across
    - E45-safe-migration/T005-carry-unfinished-follow-up-across
    - E45-safe-migration/T006-carry-the-project-documents-across
    - E45-safe-migration/T007-ask-instead-of-guessing
    - E45-safe-migration/T008-make-an-interrupted-migration-recoverable
complexity_tier: high
complexity_reason: Orders every write, owns the commit point, and decides resume, conflict, and no-op on rerun.
---

# T009: Publish the conversion once

## Problem

Everything before this task computes; this task commits. The order it writes in is the entire safety argument, because the project has to be a valid V1 project at every moment before the commit point and a valid V2 project at every moment after it.

That is what makes `schema_version: 2` the last write. It is one field in one file, and it is the closest thing to atomicity available: before it, a reader sees V1 and the converted records sit unreferenced; after it, a reader sees V2. Removing originals before their archive copies verify, or activating the schema before the records exist, would each open a window where the project is neither.

Rerun is the other half. A user whose first run was interrupted will type the same command again, and the three things that could be true — the same work is still pending, the user changed something in between, or the project is already migrated — must be distinguished rather than blended into a retry.

## Context Files

- `internal/migrate/apply.go`
- `internal/migrate/apply_test.go`
- `internal/migrate/operation.go`
- `internal/migrate/plan.go`
- `internal/migrate/manifest.go`
- `internal/migrate/replace.go`
- `internal/data/config.go`
- `internal/data/project.go`
- `internal/init/migrate_audit_skill.go`
- `.savepoint/Guardrails.md`

## Acceptance Criteria

- [x] Apply refuses a plan that is not appliable, naming every unresolved blocking ambiguity, and writes nothing.
- [x] Apply revalidates every inventoried source hash before its first write and aborts with a named conflict listing the changed paths when any differs, because the reviewed preview no longer describes the project.
- [x] Writes proceed in a fixed order: verify backups, create additive records under `objectives/` and `issues/`, write `archive/v1/` and `migrations/v1-to-v2.yml`, replace modified in-place files, remove originals only after their archive copies verify, and activate `schema_version: 2` last.
- [x] A test interrupting immediately before activation finds a project that still loads as V1, and a test interrupting immediately after finds one that loads as V2 with a complete record set.
- [x] Every archived file hashes equal to its source bytes, including the fixtures' CRLF content and unknown fields, and no original is removed before its archive copy verifies.
- [x] Pre-existing `.savepoint/migrations/` content written by `internal/init` — the archived legacy skill copies and the `README.md` — is present and byte-identical after apply.
- [x] A rerun after an interruption discovers the single incomplete operation, revalidates both sources and installed files, and resumes from the journalled step when both are unchanged.
- [x] A rerun after the user edited a source or an installed file reports a conflict naming the paths and does not overwrite the edit automatically.
- [x] A successful second full run changes no file content, no modification time, and no allocated ID; a test snapshots all three and asserts equality, satisfying FS-04.
- [x] A run against a project already at `schema_version: 2` reports that and writes nothing, and one against an unknown explicit version fails through the existing named diagnostic.
- [x] Apply never claims atomicity in its output or its documentation: it reports a recoverable operation and names its commit point.
- [x] The manifest is written with the complete identity map, archive references, typed legacy prerequisites, and recorded decisions before activation, so a migrated project can always resolve a legacy reference.

## Implementation Plan

- [x] Add `apply.go` implementing the ordered publish as explicit named steps over T008's journal, one function per step.
- [x] Implement source-hash revalidation and the conflict abort before the first write.
- [x] Implement archive writing with post-copy hash verification, and gate original removal on it.
- [x] Implement schema activation as the final single-field write through the T001 primitive.
- [x] Implement rerun resolution: resume, conflict, and already-migrated, each with its own named outcome.
- [x] Add the before-and-after-activation interruption tests asserting the project loads as V1 and as V2 respectively.
- [x] Add the idempotence snapshot test over content, mtimes, and IDs, and the `.savepoint/migrations/` coexistence assertion.
- [x] Run the focused `internal/migrate`, `internal/data`, and `internal/init` suites, then `make build && make test`.

## Context Log

**Files read:** `.savepoint/router.md`, AGENTS.md, `agent-skills/savepoint-build-task/SKILL.md`, E45-Detail.md, this task file, `internal/migrate/{operation,manifest,plan,convert,convert_docs,convert_issues,replace,inventory,classify,decisions}.go` and their tests, `internal/migrate/fixture_test.go`, `internal/data/{config,write,router,parser,project,discover}.go`, `internal/init/migrate_audit_skill.go`, both fixtures' `manifest.yml` and `project/` trees under `internal/data/testdata/migration/`.

**Files edited:**
- `internal/migrate/apply.go` (new) — `Apply`, the ordered publish (`publish`, `publishWrite`, `publishRemoval`, `ensureBackedUp`), source-hash revalidation (`revalidateSources`), the render pass that builds every write up front (`buildApplyBatch`, `stagedWrite`, `applyBatch`), and `activateSchema`, the final single-field `schema_version: 2` write kept deliberately outside the per-path journal (see its doc comment and operation.go's own note on why the "every path verified but not activated" recovery case exists).
- `internal/migrate/plan.go` — added `PlannedTarget.InstallPath()`, since `TargetPath` omits the `.md` extension for `TargetTask` (task allocation only needed the directory-qualified identity) but apply and the manifest both need the real on-disk filename.
- `internal/migrate/manifest.go` — `BuildManifest` now records `t.InstallPath()` instead of the extension-less `t.TargetPath`, so a migrated project's manifest identity map names a path that actually exists.
- `internal/data/config.go` — added `ActivateSchemaV2` (and its `setSchemaVersionNode` helper), a YAML-node patch that sets `schema_version: 2` while preserving every other key, matching the same preserve-unknown-fields pattern `write.go` already uses for V2 records; handles a project with no `config.yml` at all by producing a fresh one.
- `internal/migrate/apply_test.go` (new) — fresh end-to-end apply over both frozen fixtures (V2 load, archive byte-equality, preserved-in-place `Design.md`, config preservation), refusal paths (not appliable, plan conflict, source changed since preview) proven to write nothing via `inventory_test.go`'s existing `snapshotTree`/`assertSnapshotsEqual`, before/after-activation interruption, resume-to-convergence, resume-after-tamper (a verification conflict), second-full-run no-op, and `.savepoint/migrations/` coexistence.

**Design note on resume:** `Apply(root, plan)` does not recompute `Plan` itself — the caller supplies the plan for both a fresh apply and a resume. A resumed entry's `WriteStaged` call is skipped entirely once an entry reaches `StepStaged` or later, and `WriteStaged` itself refuses content that does not hash to the journal's frozen planned hash, so a resume whose caller-supplied plan renders different bytes (e.g. a different injected clock, since `plan.GeneratedAt` is embedded in converted Issue timestamps) fails loudly rather than silently installing the wrong content. `PendingOperationReport` already exposes the original `Journal.CreatedAt`/`OperationID`, which is what a future caller (T011's `--recover` dispatch) would replay through `Plan`'s injected clock/op-ID to reconstruct byte-identical content.

**Quality gates:** `go build ./...`, `go test ./internal/migrate/... ./internal/data/...`, and `make build && make test` all pass. This project has no `.savepoint/Health-Check.md`, so the Quick-check step is skipped per the skill's own instruction.
