---
id: E35-split-task-and-epic-audits/T001-split-audit-contracts
status: planned
objective: Replace the generic audit skill with distinct task and epic contracts backed by one shared audit method with enriched rigor.
depends_on: []
complexity_tier: high
complexity_reason: The task replaces a cross-workflow contract while preserving audit authority and closeout semantics, and adds frozen scope locks, coverage matrices, side-effect locks, convergence limits, and materiality tables.
---

# T001: Split Audit Contracts

## Problem

The generic audit skill conflates read-only review of one in-progress task with artifact-producing audit and closeout of a completed epic, and lacks rigorous audit machinery (scope locks, coverage matrices, materiality tables).

## Context Files

- `agent-skills/savepoint-audit/SKILL.md`
- `agent-skills/savepoint-audit-register/SKILL.md`
- `agent-skills/savepoint-build-task/SKILL.md`
- `agent-skills/savepoint-audit-task/SKILL.md`
- `agent-skills/savepoint-audit-epic/SKILL.md`
- `agent-skills/references/audit-method.md`
- `templates/project/agent-skills/savepoint-audit/SKILL.md`
- `templates/project/agent-skills/savepoint-audit-register/SKILL.md`
- `templates/project/agent-skills/savepoint-build-task/SKILL.md`
- `templates/project/agent-skills/savepoint-audit-task/SKILL.md`
- `templates/project/agent-skills/savepoint-audit-epic/SKILL.md`
- `templates/project/agent-skills/references/audit-method.md`

## Acceptance Criteria

### Structural split (existing)
- [ ] `savepoint-audit-task` triggers only for an explicit audit or re-audit of one in-progress task, keeps router state `task-building`, performs no writes, uses the Quick health check, and returns only `CLEAR` or `NEEDS WORK`.
- [ ] Task-audit guidance explicitly forbids `E##-Audit.md` creation and task/router lifecycle changes.
- [ ] `savepoint-audit-epic` triggers for `audit-pending` or an explicit completed-epic audit, requires a session independent from the builder, audits every completed task, and uses the Full health check.
- [ ] Epic-audit guidance reviews file reality, drift, and release guardrails; writes exactly one `E##-Audit.md`; and preserves proposal approval, apply, and closeout rules.
- [ ] Both skills use one shared reference requiring general-invariant AC reasoning; normal, boundary, malformed, failure, and bypass scenarios; public-entry-point and alternate-representation checks; and applicable behavior matrices.
- [ ] The shared method requires full re-audit after remediation, treats passing tests as supporting evidence, classifies every criterion as Proven/Finding/Unverified, and returns `NEEDS WORK` for findings or material unverified criteria.
- [ ] Consolidated findings include the violated rule, smallest reproduction, expected and actual behavior, exact file evidence, and missing test evidence.
- [ ] `savepoint-build-task` explicitly defers audit and re-audit requests to `savepoint-audit-task`.
- [ ] The live and generated generic skill folders are removed, and all changed canonical skills match their generated copies.

### Enriched audit rigor (new)
- [ ] The shared method requires a frozen scope lock (numbered, listing AC, guardrails, gates, changed files, entry points, dependencies, matrix axes, materiality boundary) before adversarial probes; re-audit reuses it immutably without adding axes.
- [ ] The shared method requires a mandatory coverage matrix with named axes: public surfaces, input shape, state, environment/output, boundaries, sequences, representations, and text classes; a finite external-boundary matrix for server/dependency code; prose checklist is not matrix evidence.
- [ ] The shared method includes a Workflow And Side-Effect Audit Lock for multi-step or side-effecting work: per-operation failure-timing table, independent oracle, and matrix completion lock before verdict.
- [ ] The shared method defines a default convergence limit: initial audit, one full re-audit, one targeted remediation, then stop.
- [ ] The shared method requires an admission ledger for re-audit: every blocking result maps to an exact frozen matrix cell; credible-blocker exception for secret exposure, cross-tenant access, destructive data loss, or sensitive-data harm.
- [ ] Both skills require file reality evidence: every file named in task logs exists, was intentionally deleted, or is recorded as discarded scratch work.
- [ ] Epic audit produces a materiality table (Finding × Likelihood × Impact × Materiality × Recommendation) for every finding.
- [ ] Epic audit includes a named Guardrails Verification subsection with Full health check.
- [ ] Epic audit produces a repository handoff result: CLEAR TO COMMIT/PUSH or NOT READY TO COMMIT/PUSH.
- [ ] Epic audit Main Findings defaults to 500-900 words; task audit defaults to 350-600 words.
- [ ] Both skills return a Final Response Output: compact chat summary with verdict, materiality table (if NEEDS WORK), gate result, and audit file link.
- [ ] Both skills reference `.savepoint/Guardrails.md` and `.savepoint/Health-Check.md` gracefully (skip when absent).
- [ ] `savepoint-build-task` references Health-Check.md Quick mode at task handoff (skip when absent).

### Source-material adaptation
- [ ] Source-material paths are adapted: `../shared/savepoint-audit-method.md` becomes the real relative path to `agent-skills/references/audit-method.md` in both skills' Read sections, and root `GUARDRAILS.md` references become `.savepoint/Guardrails.md`.
- [ ] "release guardrails audit plan" references are softened to: "the release's guardrails mapping, if your project maintains one — otherwise the relevant `.savepoint/Guardrails.md` rule IDs directly."

## Implementation Plan

- [ ] Add the shared non-triggerable audit-method reference (scope locks, coverage matrices, side-effect locks, adversarial pass, re-audit convergence, materiality) and mirror it into project templates.
- [ ] Adapt source-material paths and terminology: shared-method read path, `.savepoint/Guardrails.md` paths, softened audit-plan wording.
- [ ] Add the task-audit skill with strict read-only scope, Quick health check, evidence classification, materiality table, and result contract.
- [ ] Add the epic-audit skill with Full health check, full-epic coverage, single artifact, frozen scope lock, coverage matrix, proposal approval, guardrails verification, repository handoff, and closeout workflow.
- [ ] Update audit-register integration to name the correct split audit owner without becoming a generic alias.
- [ ] Update build-task guidance to hand explicit audit and re-audit requests to task audit and reference Health-Check.md Quick mode.
- [ ] Remove both copies of the old generic skill and verify canonical/template parity.

## Context Log

Pending.
