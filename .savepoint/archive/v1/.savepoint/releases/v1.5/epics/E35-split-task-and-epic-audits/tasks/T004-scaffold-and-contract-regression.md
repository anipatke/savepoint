---
id: E35-split-task-and-epic-audits/T004-scaffold-and-contract-regression
status: done
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
- [x] Fresh-scaffold tests prove both split skills and the shared method exist and the generic skill folder does not.
- [x] Router contract tests prove `audit-pending` selects epic audit and explicit task audit/re-audit selects task audit without a new state.
- [x] Task-skill contract tests reject epic audit artifact creation, lifecycle writes, non-Quick health checks, and result values outside `CLEAR`/`NEEDS WORK`.
- [x] Epic-skill contract tests require Full health checks, every completed task, independent session, one `E##-Audit.md`, proposal approval, and closeout authority.
- [x] Shared-method tests require all scenario classes, applicable behavior matrices, full remediation re-audit, evidence classification, and consolidated finding fields.
- [x] Both new skill directories pass the repository's frontmatter and structure validator.
- [x] Upgrade integration coverage proves the documented legacy migration behavior.
- [x] An automated stale-reference assertion excludes historical `.savepoint/releases/` records but rejects the generic name from live sources and generated output.

### New template existence
- [x] Fresh-scaffold tests prove `Guardrails.md` and `Health-Check.md` exist at expected paths in `.savepoint/`.
- [x] The Design.md template freshness test covers the directory listing change (Guardrails.md and Health-Check.md rows).

### Enriched rigor contracts
- [x] Shared-method contract tests require frozen scope lock with named fields, mandatory coverage matrix with named axes, Workflow And Side-Effect Audit Lock, default convergence limit, admission ledger requirement, and credible-blocker exception.
- [x] Epic-skill contract tests require materiality table, Guardrails Verification subsection, repository handoff result (CLEAR TO COMMIT/PUSH or NOT READY TO COMMIT/PUSH), and word count guidelines (500-900).
- [x] Task-skill contract tests require word count guidelines (350-600) and Final Response Output format.
- [x] Both skill contract tests require file reality evidence and Final Response Output format (verdict, materiality table, gate result, audit file link).

### Overall
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Update scaffold and end-to-end init fixtures for both skills, the shared method, absence of the generic folder, and Guardrails.md/Health-Check.md existence.
- [x] Update freshness coverage to compare all canonical/generated skill and shared-reference assets plus Design.md template.
- [x] Add focused skill validation for required frontmatter, folder/name agreement, and non-empty workflow sections.
- [x] Add static task, epic, router, build-handoff, and shared-method contract assertions covering enriched rigor (scope locks, matrices, side-effect locks, materiality, Guardrails Verification, repository handoff, handoff format).
- [x] Add scaffold existence assertions for Guardrails.md and Health-Check.md.
- [x] Add a live-source stale-name regression check with an explicit historical-record exclusion.
- [x] Run the full validation suite and record results in the context log.

## Context Log

### Files read

`.savepoint/router.md`, `AGENTS.md`, `agent-skills/savepoint-build-task/SKILL.md`,
`agent-skills/savepoint-audit-task/SKILL.md`, `agent-skills/savepoint-audit-epic/SKILL.md`,
`agent-skills/references/audit-method.md`, `templates/project/.savepoint/router.md`,
`templates/project/.savepoint/Guardrails.md`, `templates/project/.savepoint/Health-Check.md`,
`templates/project/.savepoint/Design.md`, `internal/init/upgrade.go`,
`internal/init/migrate_audit_skill.go`, and the existing `internal/init` test files.

### Files created

- `internal/init/audit_contract_test.go` — static routing, task-skill, epic-skill,
  and shared-method contract assertions. Every assertion runs against both the
  live copy and the scaffolded template copy.
- `internal/init/skill_validation_test.go` — frontmatter/folder-name agreement and
  non-empty section structure for all `savepoint-*` skills in both trees, plus the
  non-triggerable frontmatter contract for the shared method reference.
- `internal/init/stale_reference_test.go` — live-source and generated-output
  rejection of the retired `savepoint-audit` name.

### Files edited

- `internal/init/scaffold_test.go` — added `scaffoldFromRealTemplates` plus fresh
  scaffold assertions for the split skills, the shared method, absence of the
  generic skill folder, and interpolated `Guardrails.md`/`Health-Check.md`.
- `internal/init/integration_test.go` — the end-to-end fixture still installed
  `agent-skills/savepoint-audit/SKILL.md`; replaced with the two split skills and
  the shared reference (this was the open item recorded in T003's drift notes).
- `internal/init/template_freshness_test.go` — added Design.md directory-listing
  coverage, Guardrails/Health-Check template contracts, and a real-templates
  upgrade test proving the legacy migration end to end.

### Design decisions

- **Wrap-tolerant prose matching.** The skills and shared method are hard-wrapped
  markdown, so literal `strings.Contains` on a contract sentence fails purely on
  reflow. `assertPhrase` collapses whitespace on both sides, keeping the assertion
  about meaning rather than line width. Three assertions failed on wraps before
  this was added.
- **Stale-reference exclusions are narrow and named.** The scan skips this
  repository's own `.savepoint/releases/` (historical records — rewriting them
  would falsify history) and `README.md` (must name the retired skill to document
  what upgrade does to projects that still have it). `templates/project/.savepoint/releases/`
  is a scaffold stub, not history, so it stays in scope. A dedicated test asserts
  the exclusion holds and that historical records actually exist, so the exclusion
  cannot silently become a no-op.
- **Structure validation is tiered.** `## Trigger` and `## Workflow` are required
  of every skill (all nine satisfy this today); the fuller
  `## Purpose`/`## Read`/`## Rules` set is required of the two audit skills, whose
  read scope and rules are the contract. `savepoint-create-defect` uses
  `## Objective`/`## Constraints`, so a blanket requirement would have been a
  false constraint rather than a real one.

### Quality gates

- `make build && make test` — pass. All 9 packages ok; `internal/init` 1.18s.
- `git diff --check` — clean.
- Mutation check (each mutation applied, run, then reverted; tree confirmed
  restored and green):
  - added `| audit-pending | savepoint-audit |` to the router template →
    `TestLiveSourcesDoNotReferenceGenericAuditSkill` and
    `TestScaffoldedProjectDoesNotReferenceGenericAuditSkill` both failed with the
    offending line and number;
  - deleted `templates/project/.savepoint/Guardrails.md` →
    `TestScaffold_installsGuardrailsAndHealthCheck` and
    `TestProjectGuardrailsAndHealthCheckTemplatesExist` failed;
  - changed the task skill's Quick-only rule to Full →
    `TestTaskAuditSkillContract` failed on the exact missing phrase.

### Notes

`gofmt -w internal/init/*.go` reformatted an unrelated const alignment in
`clipboard.go`; reverted via `git checkout` to keep the diff scoped.
