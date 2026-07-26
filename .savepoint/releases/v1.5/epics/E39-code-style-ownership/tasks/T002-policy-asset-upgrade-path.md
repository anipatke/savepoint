---
id: E39-code-style-ownership/T002-policy-asset-upgrade-path
status: planned
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

- [ ] Upgrade installs `.savepoint/Guardrails.md` and `.savepoint/Health-Check.md` when they are missing from the target project.
- [ ] An existing copy of either file is left byte-identical and reported `ActionUnchanged`; user edits are never overwritten.
- [ ] The install-if-missing behavior is driven by an explicit policy-asset allowlist, not a broad `.savepoint/` prefix. Every other `.savepoint/` path stays `ActionSkipped`.
- [ ] The `.savepoint/audit/` assets keep their current behavior; the generalization does not regress `upgradeAuditAsset` coverage.
- [ ] Dry-run reports the policy-asset installs without writing to the filesystem.
- [ ] Upgrade reporting distinguishes an installed policy asset from updated, unchanged, merged, migrated, and skipped assets.
- [ ] Repeated upgrades are idempotent: the second run reports unchanged and writes nothing.
- [ ] Upgrade tests cover missing, present, user-modified, repeated, and dry-run cases for both policy assets.
- [ ] An upgraded project ends with code-style guidance reachable from AGENTS.md through the installed `Guardrails.md`.
- [ ] `make build && make test` passes.

## Implementation Plan

- [ ] Extract the policy-asset allowlist and widen the install-if-missing gate that currently keys on the `.savepoint/audit/` prefix alone.
- [ ] Reuse the `upgradeAuditAsset` install-if-missing semantics rather than adding a second write path.
- [ ] Add the reporting action for installed policy assets and thread it through dry-run.
- [ ] Add table-driven upgrade tests for the missing, present, modified, repeated, and dry-run branches.
- [ ] Add an end-to-end assertion that an upgraded project resolves the AGENTS.md code-style pointer.
- [ ] Run `make build && make test` and record results in the context log.

## Context Log

Pending.
