---
id: E40-upgrade-safety/T004-backward-compat-fixtures
status: done
objective: Pin old file shapes and the template/skill lifecycle matrix with regression tests and explicit one-command upgrade guidance.
depends_on: ["E40-upgrade-safety/T002-skill-conflict-policy", "E40-upgrade-safety/T003-agents-marker-conflict"]
complexity_tier: medium
complexity_reason: Fixture authoring, reader-contract assertions, full install/upgrade ownership coverage, and user-facing upgrade guidance span multiple boundaries.
---

# T004: Backward-Compatibility Fixtures

## Problem

`upgrade.go:118-121` skips the whole `.savepoint/` subtree, so `router.md` is never migrated and a project's router keeps its original shape indefinitely. Compatibility therefore has to hold in the reader, not the template — and nothing tests that.

`ReadState` (`internal/data/router.go:28`) currently tolerates unknown and missing YAML keys, because `yaml.Unmarshal` is called without `KnownFields`. That tolerance is accidental rather than asserted: a future change enabling strict decoding, or adding a required field, would break every existing router with no failing test.

It also hard-fails without the literal `## Current state` heading (`router.go:10`) and a fenced YAML block inside it. Those two anchors are a format contract, and nothing records that.

`template_freshness_test.go` asserts templates match *current* guidance. Nothing asserts that *old* files still work, that every installed asset has an explicit upgrade ownership policy, or that user-facing instructions distinguish the required upgrade command from its optional dry-run preview.

## Context Files

- `internal/init/testdata/legacy/` (will create)
- `internal/data/router_test.go`
- `internal/init/upgrade_test.go`
- `internal/init/integration_test.go`
- `internal/data/router.go` (contract under test)
- `README.md`

## Acceptance Criteria

- [x] `internal/init/testdata/legacy/` holds pinned fixtures: a pre-v1.5 `router.md`, an `AGENTS.md` with no markers, an `AGENTS.md` with markers, and a customized `SKILL.md`.
- [x] Fixtures are byte-frozen and commented as such — they represent files already in the wild and must never be "fixed up" to match current templates. A test that fails against a fixture indicates a compatibility break, not a stale fixture.
- [x] A router test asserts the legacy fixture parses, and that every `RouterState` field it carries reads back correctly.
- [x] A router test asserts unknown YAML keys are ignored rather than erroring, pinning the non-strict decode.
- [x] A router test asserts a missing optional key (`defect`) yields the zero value rather than an error.
- [x] A router test asserts the two structural anchors: content without `## Current state`, and content without a fenced YAML block, each produce a clear error. The shipped `templates/project/.savepoint/router.md` is asserted to contain both.
- [x] An upgrade test runs `UpgradeProjectAssets` over a project built from the legacy fixtures and asserts: the marker-less `AGENTS.md` is reported `conflict` and left byte-identical, the marked one is merged in place, and the customized `SKILL.md` survives per the T002 migration rule with its `.bak` written.
- [x] A test asserts a task file using the legacy `phase` frontmatter field still parses, covering the existing `phase` → `stage` compatibility path (`parser.go:82`).
- [x] Install and upgrade tests enforce the E40 Template and skill lifecycle matrix: all shipped assets exist after `init`; only package-owned skills/references and the marked AGENTS.md block refresh in place; Guardrails.md, Health-Check.md, and audit assets install only when missing; every other `.savepoint/` file remains untouched on upgrade.
- [x] The permanent project-owned boundary is explicit and tested: `--force` does not overwrite router, config, PRD, Design, Concept, visual identity, Guardrails, Health-Check, audit state, release PRDs, epics, tasks, defects, or audit records.
- [x] README documents `savepoint upgrade-assets` as the single required per-project update command, describes `--dry-run` as an optional read-only preview, and makes clear that installing or updating the Savepoint binary does not mutate existing projects automatically.
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Author the fixture files, each with a header comment stating its provenance and that it is frozen.
- [x] Add the router reader-contract tests to `internal/data/router_test.go`.
- [x] Add the legacy-project upgrade test to `internal/init/upgrade_test.go`, assembling a temp project from the fixtures.
- [x] Confirm the legacy `phase` task case is covered; add it if the existing parser tests only cover `stage`.
- [x] Add table-driven install/upgrade ownership coverage for every lifecycle-matrix row, including `--force` assertions at the project-owned boundary.
- [x] Update README with the one-command upgrade flow, optional dry run, and conflict-resolution behavior.

## Notes

This is the task that makes the rest of the epic durable. T001–T003 fix three specific defects; this one is the guard that catches the fourth before it ships.

The fixtures encode a policy as much as a test: Savepoint may change its templates freely, but may not change how it *reads* files it has already written without a deliberate migration.

## Context Log

**Read:** `.savepoint/router.md`, `AGENTS.md`, `agent-skills/savepoint-build-task/SKILL.md`, `.savepoint/Guardrails.md` (STYLE), `E40-Detail.md`, `T002`/`T003` task files, this task, `internal/init/upgrade.go`, `internal/init/agents.go`, `internal/init/manifest.go`, `internal/init/scaffold.go`, `internal/init/upgrade_test.go`, `internal/init/integration_test.go`, `internal/init/template_freshness_test.go`, `internal/data/router.go`, `internal/data/router_test.go`, `internal/data/parser.go`, `internal/data/parser_test.go`, `internal/testutil/fs.go`, `README.md`, `templates/project/.savepoint/router.md`.

**Added:** `internal/init/testdata/legacy/{README.md,router.md,AGENTS.unmarked.md,AGENTS.marked.md,SKILL.customized.md,task-legacy-phase.md}`, `internal/init/lifecycle_test.go`.

**Edited:** `internal/init/upgrade_test.go`, `internal/data/router_test.go`, `internal/data/parser_test.go`, `README.md`.

**Design notes:**

- One frozen copy of each fixture serves both packages: `internal/data/router_test.go` reads `../init/testdata/legacy/` rather than keeping a second copy in sync (STYLE-07). The fixture directory `README.md` carries the freeze rule, and each file repeats it in a header comment so it is visible where it is edited.
- The lifecycle-matrix coverage runs against the real `templates/project` tree via `os.DirFS`, not a hand-built `fstest.MapFS`. A hand-built table would only prove the matrix for assets someone remembered to add; walking the shipped tree means a new template lands in an ownership class the moment it ships.
- `wantOwnership` in the test restates the matrix independently of `isPolicyAsset` / `isAuditAsset` / `isPackageSkillAsset`. Asserting production predicates against themselves proves nothing; the duplication is the point.
- Ownership is asserted on a project where *every* installed file has been locally edited, run once without `--force` and once with. That is the sharpest available form of the boundary: `--force` changes what Savepoint does to files it owns, never which files it owns.
- `AGENTS.md` is excluded from the "local edit is gone after `--force`" check because prose outside the marker pair surviving *is* its contract; its refresh is asserted separately by staling the block content and checking the block, the prose, and the marker count.
- The legacy `phase` case already had frontmatter-level coverage, so the new test parses a whole frozen task file — body, checklists, and all — rather than restating what `TestParseTaskFile_readsLegacyPhaseForInProgress` proves.

**Acceptance evidence:**

| AC | Test |
|---|---|
| Legacy router fixture parses, all fields read back | `internal/data/router_test.go` `TestRouterReader_legacyFixtureParses` |
| Unknown YAML keys ignored (non-strict decode pinned) | `internal/data/router_test.go` `TestRouterReader_ignoresUnknownKeys` |
| Missing optional `defect` yields zero value | `internal/data/router_test.go` `TestRouterReaderNoDefectField` (existing) |
| Both structural anchors required; shipped template has both | `internal/data/router_test.go` `TestRouterReader_requiresStructuralAnchors`, `TestRouterReader_shippedTemplateCarriesBothAnchors` |
| Marker-less legacy `AGENTS.md` conflicts, stays byte-identical; pre-manifest skill backed up | `internal/init/upgrade_test.go` `TestUpgradeProjectAssets_legacyProjectWithUnmarkedGuide` |
| Marked legacy `AGENTS.md` merged in place | `internal/init/upgrade_test.go` `TestUpgradeProjectAssets_legacyProjectWithMarkedGuide` |
| Customized skill conflicts once provenance exists | `internal/init/upgrade_test.go` `TestUpgradeProjectAssets_legacyCustomizedSkillConflictsOnceTracked` |
| Legacy `phase` task file parses | `internal/data/parser_test.go` `TestParseTaskFile_legacyPhaseFixtureParses` |
| All shipped assets installed by `init`; manifest scope | `internal/init/lifecycle_test.go` `TestLifecycle_installWritesEveryShippedAsset` |
| Matrix ownership per path, with and without `--force` | `internal/init/lifecycle_test.go` `TestLifecycle_upgradeHonoursOwnership` |
| Package-owned skills, references, and the managed block refresh | `internal/init/lifecycle_test.go` `TestLifecycle_upgradeRefreshesPackageOwnedAssets` |
| Customized skills conflict rather than being overwritten | `internal/init/lifecycle_test.go` `TestLifecycle_upgradeConflictsOnCustomizedSkills` |
| Guardrails, Health-Check, and audit assets install only when missing | `internal/init/lifecycle_test.go` `TestLifecycle_upgradeReinstallsMissingInstallIfMissingAssets`, `TestUpgradeProjectAssets_policyAssetBranches` (existing) |

**Quality gates:** `make build && make test` — pass (all packages ok, `internal/init` 2.2s).

No `.savepoint/Health-Check.md` in this project, so the Quick health check step does not apply.

## Drift Notes

- The epic's Components and files table lists `internal/init/upgrade_test.go` and `internal/data/router_test.go` for this work. Two additions go beyond it: `internal/init/lifecycle_test.go` (the lifecycle-matrix coverage, kept out of the already 1550-line `upgrade_test.go` per STYLE-01) and `internal/data/parser_test.go` (the legacy `phase` fixture case, which belongs with the other parser tests).
- `internal/init/testdata/legacy/` also holds `task-legacy-phase.md`, one fixture more than the epic listed, so the `phase` → `stage` path is pinned against a whole frozen file rather than an inline snippet.
- No module or architectural change: all additions are tests and fixtures. The AGENTS.md Codebase Map needs no update.
