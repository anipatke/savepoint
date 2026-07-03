---
type: audit-findings
audited: 2026-07-03
---

# Audit Findings: E33 Audit Register Workflow Guidance

## Main Findings

E33's scoped implementation satisfies the task acceptance criteria for the audit-register workflow guidance. The new `savepoint-audit-register` skill defines the required read order, separates mutable register state from immutable run history, preserves stable `F###` IDs, requires named proof before `verified`, and keeps task closure plus waiver decisions under user authority. The generated project `AGENTS.md` adds concise register routing without duplicating the skill body, and the README documents the `A` read-only board overlay, markdown source-of-truth boundary, stable IDs, run history, proof-based closure, and v1.4 exclusions for dashboards, external trackers, and automated matching.

Verification passed during audit: `make build` and `make test` both completed successfully on 2026-07-03.

Applied 2026-07-03: the one documentation drift item found by this audit is resolved. `.savepoint/Design.md` now lists `savepoint-audit-register` in the bundled skill set and describes the register-first reconcile step in the agent audit workflow. The `## Proposed Changes` blocks below were applied as written and remain as the trace. No unresolved findings remain.

Lifecycle note: at audit apply/close, tasks T001-T003 were marked `status: done`, the epic was marked `status: audited`, Design.md `last_audited` was set to `v1.4/E33-audit-register-workflow-guidance`, and the router was advanced.

Coverage note: this audit examined E33's epic file, task files, named guidance/templates, README copy, the audit finding status model, help shortcut copy, and template freshness coverage. Unrelated dirty working-tree changes for E34, doctor/data internals, local config, and editor settings were intentionally not reviewed as part of E33.

## Code Style Review

- [x] One job per file - E33 keeps register workflow, generated guidance, and user documentation in their existing files.
- [x] One job per function - No runtime function changes are in E33's scoped implementation.
- [x] Test branches - No new runtime branches were introduced; template freshness and audit template tests cover the mirrored skill/template expectations.
- [x] Types document intent - Canonical finding statuses continue to live in `internal/data`, and the guidance uses those lifecycle values.
- [x] Build only what is needed - The changes are limited to audit-register workflow guidance and user documentation.
- [x] Handle errors at boundaries - No new IO/API boundary code is in E33's scoped implementation.
- [x] One source of truth - The live and scaffolded skill files match exactly; the register prompt and templates remain markdown sources of truth.
- [x] Comments explain WHY - Added guidance explains audit workflow authority and proof rules rather than restating mechanics.
- [x] Content in data files - Workflow copy lives in markdown guidance/templates, not Go logic.
- [x] Small diffs - E33's scoped changes are documentation/template focused; unrelated dirty files were excluded from this audit.

## Proposed Changes

### Target File
.savepoint/Design.md

### Replace
```md
- **Bundled Agent Skills:** Savepoint ships with custom skills (`savepoint-draft-prd`, `savepoint-system-design`, `savepoint-create-plan`, `savepoint-create-task`, `savepoint-create-defect`, `savepoint-build-task`, `savepoint-audit`) to enforce the state machine and capture release-level defects.
```

### With
```md
- **Bundled Agent Skills:** Savepoint ships with custom skills (`savepoint-draft-prd`, `savepoint-system-design`, `savepoint-create-plan`, `savepoint-create-task`, `savepoint-create-defect`, `savepoint-build-task`, `savepoint-audit`, and `savepoint-audit-register`) to enforce the state machine, capture release-level defects, and converge register-backed audit findings when `.savepoint/audit/` exists.
```

### Target File
.savepoint/Design.md

### Replace
```md
2. Reconcile      — Fresh audit agent reads router, epic detail, task files, Design.md, AGENTS.md, and scoped source/test files.
```

### With
```md
2. Reconcile      — Fresh audit agent reads router, epic detail, task files, Design.md, AGENTS.md, and scoped source/test files. When `.savepoint/audit/` exists, the agent first follows `savepoint-audit-register`: prompt/register/findings/runs, stable `F###` reconciliation, and proof rules.
```
