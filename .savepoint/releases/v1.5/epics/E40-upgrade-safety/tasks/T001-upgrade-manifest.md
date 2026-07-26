---
id: E40-upgrade-safety/T001-upgrade-manifest
status: planned
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

- [ ] `internal/init/manifest.go` defines a manifest model with a schema `version` field and a path → SHA-256 map, loaded from and saved to `.savepoint/.upgrade-manifest.yml`.
- [ ] The manifest covers `agent-skills/*/SKILL.md` paths only. `AGENTS.md` is excluded — its marker pair is already its ownership signal — and no `.savepoint/` path is included.
- [ ] Paths are stored repo-relative with forward slashes so a manifest written on Windows is readable on Linux.
- [ ] Loading a missing manifest returns an empty manifest and no error; loading a malformed one returns an error naming the file.
- [ ] `savepoint init` writes a manifest covering every skill file it scaffolds, so fresh projects start with exact provenance.
- [ ] The manifest is saved via `AtomicWrite`, and is not written at all during a dry-run upgrade.
- [ ] Unit tests cover: round-trip save/load, missing file, malformed file, forward-slash normalization, and that a dry run leaves the manifest untouched.
- [ ] `make build && make test` passes.

## Implementation Plan

- [ ] Define the manifest type, `LoadManifest(dir)`, `Record(path, content)`, and `Save(dir)` in `internal/init/manifest.go`.
- [ ] Hash with `crypto/sha256` over exact file bytes, no line-ending normalization, so the hash matches what was written.
- [ ] Call the manifest write at the end of the `init` scaffold path.
- [ ] Have `UpgradeProjectAssets` load the manifest at start and save it once at the end, skipping the save on dry runs.
- [ ] Add `manifest_test.go` covering the AC list.

## Notes

Recording happens in this task; *acting* on the recorded hashes is T002. After T001 alone, upgrade behaviour is unchanged — the manifest is written and ignored. This keeps the diff reviewable and means a bad manifest cannot destroy files before the conflict logic exists to protect them.
