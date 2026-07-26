---
id: E35-split-task-and-epic-audits/T004-scaffold-and-contract-regression
status: planned
objective: Prove scaffolded and upgraded projects enforce the split audit contracts, include new Guardrails.md and Health-Check.md templates, and pass enriched rigor contract assertions.
depends_on:
  - E35-split-task-and-epic-audits/T003-upgrade-audit-skill-migration
complexity_tier: medium
complexity_reason: Regression proof spans scaffold output, static workflow contracts, enrichment assertions, validator checks, stale references, and new template existence.
---

# T004: Scaffold and Contract Regression

## Problem

The split and enrichment are not durable unless automated tests check generated filesystem output, workflow authority, enriched rigor contracts, skill validity, migration behavior, and removal of stale live references.

## Context Files

- `internal/init/scaffold_test.go`
- `internal/init/integration_test.go`
- `internal/init/template_freshness_test.go`
- `internal/init/skill_validation_test.go`
- `internal/init/upgrade_test.go`
- `templates/project/.savepoint/router.md`
- `templates/project/AGENTS.md`
- `templates/project/agent-skills/savepoint-audit-task/SKILL.md`
- `templates/project/agent-skills/savepoint-audit-epic/SKILL.md`
- `templates/project/agent-skills/references/audit-method.md`
- `templates/project/agent-skills/savepoint-build-task/SKILL.md`
- `templates/project/.savepoint/Guardrails.md`
- `templates/project/.savepoint/Health-Check.md`

## Acceptance Criteria

### Structural split and scaffold (existing coverage)
- [ ] Fresh-scaffold tests prove both split skills and the shared method exist and the generic skill folder does not.
- [ ] Router contract tests prove `audit-pending` selects epic audit and explicit task audit/re-audit selects task audit without a new state.
- [ ] Task-skill contract tests reject epic audit artifact creation, lifecycle writes, non-Quick health checks, and result values outside `CLEAR`/`NEEDS WORK`.
- [ ] Epic-skill contract tests require Full health checks, every completed task, independent session, one `E##-Audit.md`, proposal approval, and closeout authority.
- [ ] Shared-method tests require all scenario classes, applicable behavior matrices, full remediation re-audit, evidence classification, and consolidated finding fields.
- [ ] Both new skill directories pass the repository's frontmatter and structure validator.
- [ ] Upgrade integration coverage proves the documented legacy migration behavior.
- [ ] An automated stale-reference assertion excludes historical `.savepoint/releases/` records but rejects the generic name from live sources and generated output.

### New template existence
- [ ] Fresh-scaffold tests prove `Guardrails.md` and `Health-Check.md` exist at expected paths in `.savepoint/`.
- [ ] The Design.md template freshness test covers the directory listing change (Guardrails.md and Health-Check.md rows).

### Enriched rigor contracts
- [ ] Shared-method contract tests require frozen scope lock with named fields, mandatory coverage matrix with named axes, Workflow And Side-Effect Audit Lock, default convergence limit, admission ledger requirement, and credible-blocker exception.
- [ ] Epic-skill contract tests require materiality table, Guardrails Verification subsection, repository handoff result (CLEAR TO COMMIT/PUSH or NOT READY TO COMMIT/PUSH), and word count guidelines (500-900).
- [ ] Task-skill contract tests require word count guidelines (350-600) and Final Response Output format.
- [ ] Both skill contract tests require file reality evidence and Final Response Output format (verdict, materiality table, gate result, audit file link).

### Overall
- [ ] `make build && make test` passes.

## Implementation Plan

- [ ] Update scaffold and end-to-end init fixtures for both skills, the shared method, absence of the generic folder, and Guardrails.md/Health-Check.md existence.
- [ ] Update freshness coverage to compare all canonical/generated skill and shared-reference assets plus Design.md template.
- [ ] Add focused skill validation for required frontmatter, folder/name agreement, and non-empty workflow sections.
- [ ] Add static task, epic, router, build-handoff, and shared-method contract assertions covering enriched rigor (scope locks, matrices, side-effect locks, materiality, Guardrails Verification, repository handoff, handoff format).
- [ ] Add scaffold existence assertions for Guardrails.md and Health-Check.md.
- [ ] Add a live-source stale-name regression check with an explicit historical-record exclusion.
- [ ] Run the full validation suite and record results in the context log.

## Context Log

Pending.
