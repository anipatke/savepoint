---
id: E40-upgrade-safety/T001-upgrade-manifest
status: done
objective: Add an upgrade manifest at .savepoint/.upgrade-manifest.yml recording the SHA-256 of each skill file as Savepoint last wrote it, written at both init and upgrade-assets.
depends_on: []
complexity_tier: medium
complexity_reason: New state file plus model, but writes hook into two existing well-tested paths.
---

# T001: Upgrade Manifest

## Problem

Upgrade decides what to do by comparing on-disk content to the incoming template. That comparison cannot distinguish "the user edited this file" from "this is the previous version's copy", so `upgrade.go:200` overwrites every differing skill file with no warning and no backup.

Distinguishing the two requires a third data point: what Savepoint itself last wrote. Nothing records that today — `config.yml` carries no version or provenance state.

## Context Files

- `internal/init/manifest.go` (will create)
- `internal/init/manifest_test.go` (will create)
- `internal/init/scaffold.go` (write manifest at init)
- `internal/init/upgrade.go` (read/record during upgrade)
- `internal/init/write.go` (`AtomicWrite` for the save path)

## Acceptance Criteria

- [x] `internal/init/manifest.go` defines a manifest model with a schema `version` field and a path → SHA-256 map, loaded from and saved to `.savepoint/.upgrade-manifest.yml`.
- [x] The manifest covers `agent-skills/*/SKILL.md` paths only. `AGENTS.md` is excluded — its marker pair is already its ownership signal — and no `.savepoint/` path is included.
- [x] Paths are stored repo-relative with forward slashes so a manifest written on Windows is readable on Linux.
- [x] Loading a missing manifest returns an empty manifest and no error; loading a malformed one returns an error naming the file.
- [x] `savepoint init` writes a manifest covering every skill file it scaffolds, so fresh projects start with exact provenance.
- [x] The manifest is saved via `AtomicWrite`, and is not written at all during a dry-run upgrade.
- [x] Unit tests cover: round-trip save/load, missing file, malformed file, forward-slash normalization, and that a dry run leaves the manifest untouched.
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Define the manifest type, `LoadManifest(dir)`, `Record(path, content)`, and `Save(dir)` in `internal/init/manifest.go`.
- [x] Hash with `crypto/sha256` over exact file bytes, no line-ending normalization, so the hash matches what was written.
- [x] Call the manifest write at the end of the `init` scaffold path.
- [x] Have `UpgradeProjectAssets` load the manifest at start and save it once at the end, skipping the save on dry runs.
- [x] Add `manifest_test.go` covering the AC list.

## Context Log

**Read:** `.savepoint/router.md`, `AGENTS.md`, `.savepoint/Guardrails.md`, `agent-skills/savepoint-build-task/SKILL.md`, `E40-Detail.md`, this task, `internal/init/upgrade.go`, `internal/init/scaffold.go`, `internal/init/write.go`, `internal/data/config.go` (yaml load/save convention).

**Edited:** `internal/init/manifest.go` (new), `internal/init/manifest_test.go` (new), `internal/init/scaffold.go`, `internal/init/upgrade.go`.

**Design notes:**

- `Record` filters to `agent-skills/*/SKILL.md` itself via `isManifestPath`, so callers record every asset they write without repeating the scope rule (one source of truth for scope).
- Scaffold records the *interpolated* bytes it writes, not the raw template, so the hash matches the file on disk.
- Upgrade records on write and on unchanged (on-disk content equals the template there), so a project that upgrades before this task gains provenance on the next run.
- No dependency change: `gopkg.in/yaml.v3` was already used by `internal/data` and `internal/doctor`.

**Acceptance evidence** (`internal/init/manifest_test.go`):

| AC | Test |
|---|---|
| Model, load/save round-trip | `TestManifest_roundTrip` |
| Scope is skill entrypoints only | `TestManifest_recordIgnoresPathsOutsideScope`, `TestScaffold_writesManifestForSkills` |
| Forward-slash paths | `TestManifest_recordNormalizesToForwardSlashes` |
| Missing manifest → empty, no error | `TestLoadManifest_missingFileIsEmpty` |
| Malformed manifest → error naming the file | `TestLoadManifest_malformedFileNamesTheFile` |
| `init` writes provenance | `TestScaffold_writesManifestForSkills` |
| Upgrade records hashes | `TestUpgradeProjectAssets_recordsSkillHashes` |
| Dry run writes nothing | `TestUpgradeProjectAssets_dryRunLeavesManifestUntouched`, `TestUpgradeProjectAssets_dryRunDoesNotCreateManifest` |

**Quality gates:** `make build && make test` — pass, all packages ok.

**Health check:** skipped; this project has no `.savepoint/Health-Check.md`.

## Drift Notes

`internal/init/manifest.go` adds upgrade provenance state to `internal/init`, which the AGENTS.md Codebase Map currently describes as "target validation, scaffold writing from templates, managed AGENTS.md merge behavior, and safe project asset refresh". The map entry is module-level and not edited here; if E40's later tasks keep growing this area, the `internal/init/` row should gain "upgrade provenance manifest" at epic audit.

## Notes

Recording happens in this task; *acting* on the recorded hashes is T002. After T001 alone, upgrade behaviour is unchanged — the manifest is written and ignored. This keeps the diff reviewable and means a bad manifest cannot destroy files before the conflict logic exists to protect them.
