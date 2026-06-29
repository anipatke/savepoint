---
id: E30-audit-register-templates/T004-upgrade-assets-and-template-regression
status: done
objective: Ensure upgrade-assets can add audit-register documentation safely.
depends_on:
    - E30-audit-register-templates/T003-run-history-template-and-scaffold
complexity_tier: medium
complexity_reason: Upgrade behavior must preserve user files while adding new managed assets.
---

# T004: Upgrade Assets and Template Regression

## Problem

Existing projects need a safe path to receive audit-register documentation assets without overwriting user-maintained audit state.

## Context Files

- `internal/init/upgrade.go`
- `internal/init/upgrade_test.go`
- `internal/init/template_freshness_test.go`
- `internal/init/write.go`
- `templates/project/.savepoint/audit/prompt.md`
- `templates/project/.savepoint/audit/register.md`

## Acceptance Criteria

- [x] Upgrade-assets can add missing audit-register scaffold files to existing projects.
- [x] Upgrade-assets does not overwrite user-edited prompt, register, finding, or run files.
- [x] Managed/unmanaged file behavior stays consistent with existing project asset rules.
- [x] Template freshness tests account for the audit-register files.
- [x] `go test ./internal/init` passes.

## Implementation Plan

- [x] Extend upgrade asset discovery for audit-register template files.
- [x] Preserve existing user files using the established safe-write behavior.
- [x] Add tests for missing, existing pristine, and existing edited audit files.
- [x] Update template freshness expectations.
- [x] Verify no generated source or runtime files are introduced.

## Context Log

**Read:** `internal/init/upgrade.go`, `internal/init/upgrade_test.go`, `internal/init/template_freshness_test.go`, `internal/init/write.go`, `internal/init/scaffold.go`, `internal/init/scaffold_test.go`, `internal/init/integration_test.go`, `main.go`, epic `E30-Detail.md`, and the four `templates/project/.savepoint/audit/` files.

**Edited:**
- `internal/init/upgrade.go` — added an audit-asset branch to the walk before the generic `.savepoint` skip. New `auditAssetPrefix` const, `isAuditAsset`, and `upgradeAuditAsset` helper implement additive-only behavior: a missing file under `.savepoint/audit/` is written (`ActionUpdated`) via `AtomicWrite`; an existing file is left untouched (`ActionUnchanged`). Dry-run reports `ActionUpdated` for missing files without writing.
- `internal/init/upgrade_test.go` — added `auditTemplates`/`actionFor` helpers and four cases: adds missing assets, preserves user-edited register (and still adds missing siblings), reports pristine assets unchanged, and dry-run does not write.
- `internal/init/template_freshness_test.go` — added `TestUpgradeAddsAuditRegisterTemplatesFromRealTemplates`, driving the real `templates/project` tree through `UpgradeProjectAssets` (add then idempotent rerun) plus an `upgradeActionFor` helper.

**Decisions:** Audit-register files are user-maintained state, so they follow create-if-missing (unmanaged) semantics — distinct from `agent-skills/*/SKILL.md`, which stay managed/overwritten. Content is written verbatim (no interpolation); the audit templates contain no `{{…}}` tokens. The generic `.savepoint` skip is preserved for non-audit files, so `TestUpgradeProjectAssets_skipsSavepointDir` still holds.

**Quality gates:** `go test ./internal/init` → ok. `make build && make test` → all packages ok. No new files/modules or architecture delta beyond the documented `internal/init` responsibilities; no Drift Notes required.
