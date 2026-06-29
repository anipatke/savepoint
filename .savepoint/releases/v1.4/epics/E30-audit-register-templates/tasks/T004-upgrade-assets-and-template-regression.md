---
id: E30-audit-register-templates/T004-upgrade-assets-and-template-regression
status: planned
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

- [ ] Upgrade-assets can add missing audit-register scaffold files to existing projects.
- [ ] Upgrade-assets does not overwrite user-edited prompt, register, finding, or run files.
- [ ] Managed/unmanaged file behavior stays consistent with existing project asset rules.
- [ ] Template freshness tests account for the audit-register files.
- [ ] `go test ./internal/init` passes.

## Implementation Plan

- [ ] Extend upgrade asset discovery for audit-register template files.
- [ ] Preserve existing user files using the established safe-write behavior.
- [ ] Add tests for missing, existing pristine, and existing edited audit files.
- [ ] Update template freshness expectations.
- [ ] Verify no generated source or runtime files are introduced.

## Context Log

Pending.
