---
id: E47-onboarding-upgrades/T005-refresh-assets-that-match-the-projects-own-version
title: Refresh assets that match the project's own version
status: done
objective: Make upgrade-assets select its template tree from the project's declared schema_version, leaving V1 projects byte-identical and never writing the version itself.
depends_on:
    - E47-onboarding-upgrades/T001-ship-the-v2-project-assets-from-their-own-tree
complexity_tier: high
complexity_reason: Adds a version gate to a command with a large existing contract surface, across init, cmd, and main.
---

# T005: Refresh assets that match the project's own version

## Problem

`upgrade-assets` installs one asset stream into every project it is pointed at. After T001 there are two streams, and installing the wrong one is worse than installing nothing: a V1 project given the four V2 skills gets instructions for a router state it does not have and record families it does not contain, and a migrated V2 project given the V1 skills gets a routing table pointing at epics and audits that no longer exist.

The gate is the project's own `schema_version` and nothing else. Not the package version, not `.upgrade-manifest.yml`'s version, not whether V2 files happen to be lying around — `internal/data`'s `ReadSchemaVersion` already owns this question and returns named diagnostics for a malformed or unsupported value. Asset installation reads that declaration and never writes it: `migrate` is the only operation that changes what a project is, and an upgrade that could promote a project to V2 would be a migration wearing a different name.

A V1 project's upgrade must come out the same as it does today, byte for byte and action for action. This is the compatibility promise of the transitional period: a user who has not migrated yet is not punished for it. The one addition is an informational line naming `savepoint migrate` as the route to V2 — reported, not enforced, and not an error.

The failure mode to design against is a version check that runs too late. `ReadSchemaVersion` must be consulted before the first write, alongside the manifest-writability probe and the pending-migration guard the command already runs at that boundary, so an unreadable version refuses cleanly with nothing changed (FS-06). A dry run must reach the same decision and still write nothing (FS-03).

## Context Files

- `internal/init/upgrade.go`
- `internal/init/upgrade_test.go`
- `internal/init/upgrade_failure_test.go`
- `internal/init/lifecycle_test.go`
- `internal/init/manifest.go`
- `internal/init/manifest_test.go`
- `internal/data/config.go`
- `internal/data/errors.go`
- `cmd/upgrade-assets.go`
- `cmd/upgrade-assets_test.go`
- `main.go`
- `main_test.go`
- `AGENTS.md`
- `.savepoint/releases/v2/v2-Design.md`

## Acceptance Criteria

- [x] `UpgradeProjectAssets` resolves the project's schema version through `data.ReadSchemaVersion` and selects the V2 tree for `SchemaVersionV2` and the V1 tree otherwise.
- [x] A project with no `schema_version` key receives exactly the actions and exactly the bytes it receives before this task, proven by comparing a full upgrade against a recorded baseline.
- [x] A V1 project's report carries one informational entry naming `savepoint migrate` as the route to V2, and that entry is not a failure, a conflict, or a nonzero exit.
- [x] A V2 project's upgrade installs the four V2 skills and three shared references and installs no V1 skill.
- [x] No upgrade path writes, adds, or removes `schema_version`, proven by comparing `config.yml` bytes before and after on both a V1 and a V2 project.
- [x] A malformed `schema_version` fails with `ErrMalformedSchemaVersion` and an unsupported one with `ErrUnsupportedSchemaVersion`, each naming the config path, with no file written and no sidecar created.
- [x] The version check runs before the first write, at the same boundary as the manifest-writability probe and the pending-migration guard, and an upgrade that changes nothing still touches nothing.
- [x] `--dry-run` reaches the same tree selection and the same refusals on every path above and writes nothing, verified by snapshotting file bytes and mtimes.
- [x] The existing pending-migration refusal, manifest-writability refusal, and partial-failure reporting behave unchanged on both trees.
- [x] A second upgrade of each kind is a no-op: no file content, mtime, or manifest entry changes.
- [x] `cmd/upgrade-assets.go` gains no domain logic; tree selection lives in `internal/init` (ARCH-01).
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Record a baseline of a full V1 upgrade — actions and resulting bytes — as a fixture, before changing any behavior, so the no-regression criterion is checkable rather than asserted.
- [x] Add tree resolution to `internal/init`: read the schema version, map it to a tree, and return the named diagnostic unchanged for a version that cannot be interpreted.
- [x] Pass both embedded sub-filesystems from `main.go` into the upgrade entry point, keeping `cmd/upgrade-assets.go` to parsing and dispatch.
- [x] Place the version read at the existing first-write guard boundary, next to the manifest probe and migration guard.
- [x] Add the informational migrate-route entry for a legacy project, with its own action or note so the report format stays unambiguous.
- [x] Verify against the baseline fixture that the V1 path is unchanged.
- [x] Add the V2-project upgrade test, the malformed and unsupported version tests, the dry-run parity tests, and the second-run no-op tests.
- [x] Update the `internal/init/` Codebase Map row and the `templates/` row in `AGENTS.md` if the mirrored wording no longer describes the code.
- [x] Run `go test ./internal/init/... ./cmd/... ./internal/data/...`, then `make build && make test`.

## Context Log

**Files read:** `internal/init/upgrade.go`, `internal/init/upgrade_test.go`, `internal/init/upgrade_failure_test.go`, `internal/init/lifecycle_test.go`, `internal/init/manifest.go`, `internal/init/manifest_test.go`, `internal/data/config.go`, `internal/data/errors.go`, `internal/data/project.go`, `cmd/upgrade-assets.go`, `cmd/upgrade-assets_test.go`, `main.go`, `main_test.go`, `AGENTS.md`, `internal/init/migrate_audit_skill_test.go`, `internal/init/template_freshness_test.go`.

**Files edited:**
- `internal/init/upgrade.go` — added `ActionInfo` and `noteMigrateRoute`; renamed the old single-tree `UpgradeProjectAssets` wrapper to unexported `upgradeAssetsFromTree` (unchanged body); added a new exported `UpgradeProjectAssets(v1Templates, v2Templates fs.FS, targetDir string, dryRun, force bool)` that resolves `data.ReadSchemaVersion` before any write, selects the matching tree, delegates to `upgradeAssetsFromTree`, and appends the migrate-route informational entry for a V1 project; `Format()` now reports an `Info` count.
- `main.go` — `upgradeAssetsRunner` now builds both embedded sub-filesystems (`templates/project`, `templates/project-v2`) and passes both into `UpgradeProjectAssets`.
- `internal/init/upgrade_test.go`, `upgrade_failure_test.go`, `lifecycle_test.go`, `manifest_test.go`, `migrate_audit_skill_test.go`, `template_freshness_test.go` — mechanically renamed every real call site from `UpgradeProjectAssets(...)` to `upgradeAssetsFromTree(...)` (same single-tree signature); these tests exercise upgrade policy independent of version dispatch, so their content is otherwise untouched.
- `internal/init/upgrade_schema_test.go` (new) — version-dispatch coverage: tree selection for missing/explicit `schema_version`, the migrate-route note, malformed/unsupported refusals (clean, config-path-named, no manifest created), dry-run parity including refusal parity, `schema_version` left untouched on both trees, second-run no-op (content + mtime) on both trees, the pending-migration guard firing through dispatch on both trees, and a baseline-regression test comparing the V1 dispatch path byte-for-byte against a direct `upgradeAssetsFromTree` call.
- `main_test.go` — added `TestMainUpgradeAssetsV1ProjectInstallsV1SkillsAndNamesMigrateRoute` and `TestMainUpgradeAssetsV2ProjectInstallsOnlyV2Skills`, exercising the real embedded template trees end to end through `savepoint upgrade-assets`.
- `AGENTS.md` — updated the `internal/init/` Codebase Map row to describe schema-version dispatch instead of "upgrade-assets still serves `templates/project`".

**Quality gates:** `go build ./...`, `go vet ./...`, `go test ./...`, and `make build && make test` all pass.
