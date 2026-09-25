---
id: E02-domain-phase-model/T003-simplify-config
status: planned
objective: "Remove quality_gates, audit, verify_strict from config model; keep theme only"
depends_on: []
---

# T003: Simplify Config

## Acceptance Criteria

- `config.ts` contains only `ThemeConfig` and a minimal `SavepointConfig` wrapper.
- `CONFIG_DEFAULTS` has no `quality_gates`, `audit`, or `verify_strict` fields.
- `applyConfigDefaults` only processes theme keys.
- `readers/config.ts` still reads and returns a valid config (theme defaults if missing).
- All config tests pass.

## Implementation Plan

- [ ] Read existing `config.ts` and `test/domain/config.test.ts`.
- [ ] Remove `QualityGates`, `AuditConfig` interfaces.
- [ ] Remove `quality_gates`, `audit`, `verify_strict` from `SavepointConfig`.
- [ ] Update `CONFIG_DEFAULTS` to theme-only.
- [ ] Simplify `applyConfigDefaults` to only handle theme fields.
- [ ] Update `test/domain/config.test.ts`.
- [ ] Run `npm test`.
