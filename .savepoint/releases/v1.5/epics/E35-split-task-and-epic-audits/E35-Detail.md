---
type: epic-design
status: planned
---

# E35: Split Task and Epic Audits

## Purpose

Give task review and epic closeout distinct, unambiguous skills while keeping their audit reasoning consistent through one shared method.

## What this epic adds

### Structural split
- A read-only `savepoint-audit-task` skill for explicit audit or re-audit requests on one in-progress task.
- An independent-session `savepoint-audit-epic` skill for `audit-pending` and explicit completed-epic audits.
- A shared audit-method reference covering invariant reasoning, scenario matrices, evidence classification, and consolidated findings.
- Generated router and agent guidance that treats task audit as a request-qualified `task-building` override.
- An upgrade migration that removes the old generic skill from the active catalog without discarding legacy user content.
- Automated scaffold, migration, contract, and stale-reference coverage.

### Enriched audit rigor
- **Frozen scope locks.** Initial audit freezes a numbered scope lock (AC, guardrails, gates, changed files, entry points, dependencies, matrix axes, materiality boundary) before adversarial probes. Re-audits reuse it immutably without adding axes or dependency layers.
- **Mandatory coverage matrices.** Both skills build a concrete matrix with named axes (public surfaces, input shape, state, environment, boundaries, sequences, representations, text classes) before verdict. Prose checklist is not matrix evidence.
- **Workflow and side-effect audit locks.** For multi-step or side-effecting work, a per-operation failure-timing table with independent oracle; matrix completion lock before verdict.
- **Admission ledgers and convergence limits.** Re-audit requires an exact frozen matrix cell for every blocking result. Default limit: initial + one re-audit + one targeted remediation, then stop.
- **Credible-blocker exception.** Named types (secret exposure, cross-tenant access, destructive data loss, sensitive-data harm) can override the frozen perimeter.
- **File reality evidence.** Every file named in task logs must exist, be intentionally deleted, or be recorded as discarded scratch work.
- **Materiality tables.** Every finding gets likelihood/impact/materiality/recommendation.
- **Repository handoff results.** Epics produce CLEAR TO COMMIT/PUSH or NOT READY TO COMMIT/PUSH.
- **Guardrails Verification.** Named subsection in epic audits with Full health check.
- **Word count guidelines.** Epic audit 500-900 words; task audit 350-600 words.
- **Final Response Output.** Compact chat summary with verdict, materiality table, gate result, and audit file link.

## Components and files

| Module | Purpose |
|--------|---------|
| `agent-skills/` | Canonical split skills, shared audit method, audit-register, build-task, and guardrails/health-check references |
| `templates/project/agent-skills/` | Generated copies of canonical audit assets |
| `templates/project/.savepoint/router.md` | State routing plus the explicit task-audit override |
| `AGENTS.md` and `templates/project/AGENTS.md` | Live and generated skill catalogs, audit authority rules, and Guardrails/Health-Check discoverability |
| `templates/project/.savepoint/Design.md` | Directory listing for Guardrails.md and Health-Check.md; audit pipeline reference |
| `agent-skills/savepoint-audit-register/SKILL.md` and template copy | Update to reference `savepoint-audit-epic` |
| `internal/init/upgrade.go` | Existing-project migration and shared-reference installation |
| `internal/init/*_test.go` | Scaffold, upgrade, contract, freshness, and integration proof |
| `README.md` | User-facing audit and compatibility guidance |

## Architectural delta

Audit selection changes from one phase-wide skill to two intent-specific skills. Router state remains unchanged: task audit is dispatch metadata on `task-building`, while epic audit remains the `audit-pending` phase workflow. Shared reasoning lives in one non-triggerable reference file consumed by both skills. New scaffolds omit the old generic folder. Upgrade-assets relocates any legacy generic `SKILL.md` to a non-triggerable migration archive before removing its active folder, then installs the split skills and shared method.

## Boundaries

**In scope:**

- Canonical and generated skill content (structural split + enriched rigor)
- Shared audit-method content (scope locks, coverage matrices, adversarial pass, re-audit, materiality)
- Router, build-task, contributor, and user guidance (including cross-reference updates to AGENTS.md, Design.md, audit-register, README)
- Init and upgrade-assets behavior
- Static contract and filesystem-level regression tests (including Guardrails.md and Health-Check.md existence)
- Preservation of historical release records containing the old name

**Out of scope:**

- New workflow states or task lifecycle values
- Automatic task completion, epic approval, or audit finding remediation
- Runtime policing of agent filesystem access beyond the generated contract
- Rewriting prior release artifacts to use the new terminology

## Quality gates

- `make build && make test` passes.
- Fresh scaffold and upgrade tests prove split-skill and shared-method behavior.
- Contract tests prove enriched rigor: scope locks, coverage matrices, side-effect locks, materiality tables, file reality evidence, Guardrails Verification, repository handoff.
- Fresh scaffold includes Guardrails.md and Health-Check.md at expected paths.
- A stale-reference check finds no generic skill in live generators, templates, catalogs, or documentation while excluding historical release records.
- Both split skill folders pass the repository's automated skill validator.

## Build-order dependency

This epic builds after E37 and E38: its skill contracts reference `.savepoint/Guardrails.md` and `.savepoint/Health-Check.md`, and T004 asserts their scaffold existence. See the PRD Build order section.

## Open decisions

None. Existing-project upgrades archive legacy generic skill content outside `agent-skills/`, remove the triggerable folder, and install the split skills.
