---
id: E39-code-style-ownership/T002-policy-asset-upgrade-path
status: done
objective: Deliver Guardrails.md and Health-Check.md to existing projects on upgrade by generalizing the install-if-missing gate, without overwriting user edits.
depends_on:
    - E39-code-style-ownership/T001-point-guidance-at-style-rules
complexity_tier: medium
complexity_reason: Widens an upgrade gate that must stay non-destructive, dry-run correct, idempotent, and narrow enough to skip the rest of .savepoint.
---

# T002: Policy-Asset Upgrade Path

## Problem

`upgrade-assets` skips the entire `.savepoint/` subtree (`internal/init/upgrade.go:118-121`), so the E37 and E38 templates reach fresh `init` projects only. An existing project upgrading to v1.5 receives the refreshed AGENTS.md managed block — which after T001 points at `.savepoint/Guardrails.md` — but never receives the file itself. Code-style and health-check guidance would disappear for every upgraded project.

## Context Files

- `internal/init/upgrade.go`
- `internal/init/upgrade_test.go`
- `templates/project/.savepoint/Guardrails.md`
- `templates/project/.savepoint/Health-Check.md`

## Acceptance Criteria

- [x] Upgrade installs `.savepoint/Guardrails.md` and `.savepoint/Health-Check.md` when they are missing from the target project.
- [x] An existing copy of either file is left byte-identical and reported `ActionUnchanged`; user edits are never overwritten.
- [x] The install-if-missing behavior is driven by an explicit policy-asset allowlist, not a broad `.savepoint/` prefix. Every other `.savepoint/` path stays `ActionSkipped`.
- [x] The `.savepoint/audit/` assets keep their current behavior; the generalization does not regress `upgradeAuditAsset` coverage.
- [x] Dry-run reports the policy-asset installs without writing to the filesystem.
- [x] Upgrade reporting distinguishes an installed policy asset from updated, unchanged, merged, migrated, and skipped assets.
- [x] Repeated upgrades are idempotent: the second run reports unchanged and writes nothing.
- [x] Upgrade tests cover missing, present, user-modified, repeated, and dry-run cases for both policy assets.
- [x] An upgraded project ends with code-style guidance reachable from AGENTS.md through the installed `Guardrails.md`.
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Extract the policy-asset allowlist and widen the install-if-missing gate that currently keys on the `.savepoint/audit/` prefix alone.
- [x] Reuse the `upgradeAuditAsset` install-if-missing semantics rather than adding a second write path.
- [x] Add the reporting action for installed policy assets and thread it through dry-run.
- [x] Add table-driven upgrade tests for the missing, present, modified, repeated, and dry-run branches.
- [x] Add an end-to-end assertion that an upgraded project resolves the AGENTS.md code-style pointer.
- [x] Run `make build && make test` and record results in the context log.

## Context Log

### Files read

- `.savepoint/router.md`, `.savepoint/Guardrails.md`, `AGENTS.md`
- `E39-Detail.md`, this task file
- `internal/init/upgrade.go`, `internal/init/upgrade_test.go`, `internal/init/scaffold.go`, `internal/init/template_freshness_test.go`
- `templates/project/.savepoint/Guardrails.md`, `templates/project/.savepoint/Health-Check.md`, `templates/project/AGENTS.md`, `README.md`

### Files edited

- `internal/init/upgrade.go` — added `ActionInstalled` and its `Format()` count; replaced the `isAuditAsset` gate with `installMissingAction(path)`, which resolves an install-if-missing path to the action it reports; added the exact `policyAssets` allowlist; renamed `upgradeAuditAsset` to `installMissingAsset` with an `installAction` parameter so both asset families share the single write path.
- `internal/init/upgrade_test.go` — `policyTemplates` fixture plus `installsMissingPolicyAssets`, table-driven `policyAssetBranches` (missing / pristine / user-modified / dry-run missing / dry-run present, each for both assets), and `policyAssetsIdempotent`; extended `TestUpgradeReport_format` with the installed and migrated counts.
- `internal/init/template_freshness_test.go` — `TestUpgradeDeliversPolicyAssetsFromRealTemplates`: upgrade against the real `templates/project` FS reports both policy assets installed, and the AGENTS.md pointer resolves to `STYLE-01..10` in the installed `Guardrails.md`.
- `README.md` — one paragraph recording the new install-if-missing behavior, since the previous text described upgrade as never touching `.savepoint/` beyond the migrations archive.

### Decisions

- Audit assets keep reporting `ActionUpdated` on install (unchanged behavior, existing tests untouched); only policy assets report the new `ActionInstalled`.
- The shared install path now interpolates `{{PROJECT_NAME}}` exactly as `Scaffold` does, so an upgraded project's `Guardrails.md` reads the same as a freshly scaffolded one. Audit templates contain no placeholders, so their bytes are unaffected — asserted by the pre-existing audit-asset tests.
- No `.savepoint/Health-Check.md` in this repo, so the build skill's Quick check step was skipped; absence is not a finding.

### Quality gates

- `make build && make test` — pass; all packages ok, `internal/init` 1.265s.
- `go vet ./internal/init/` — clean.
- New tests verified green by name: `TestUpgradeDeliversPolicyAssetsFromRealTemplates`, `TestUpgradeProjectAssets_installsMissingPolicyAssets`, `TestUpgradeProjectAssets_policyAssetBranches` (10 subtests), `TestUpgradeProjectAssets_policyAssetsIdempotent`.

## Drift Notes

`.savepoint/Design.md:26` says `upgrade-assets` "skips project-owned state, except" the audit-skill migrations archive. That is now incomplete: upgrade also installs the two allowlisted policy assets when missing. The architecture is the one E39 designed, but the Design.md sentence needs the policy-asset exception added at audit apply time. No new files, modules, or packages; the Codebase Map is unchanged.
