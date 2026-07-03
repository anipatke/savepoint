---
id: E30-audit-register-templates/T002-register-and-finding-templates
status: done
objective: Add scaffolded register and finding-folder guidance for durable audit findings.
depends_on:
    - E30-audit-register-templates/T001-audit-prompt-template
complexity_tier: low
complexity_reason: The task adds file templates and field guidance without parsing behavior.
---

# T002: Register and Finding Templates

## Problem

There is no durable markdown location for current audit findings, stable finding IDs, or convergence summaries.

## Context Files

- `templates/project/.savepoint/audit/register.md`
- `templates/project/.savepoint/audit/findings/README.md`
- `templates/project/.savepoint/audit/prompt.md`

## Acceptance Criteria

- [x] The register template includes columns or sections for ID, title, status, severity, confidence, linked work item, first seen, last seen, and proof.
- [x] The findings folder guidance defines the `F###-slug.md` naming convention.
- [x] Finding guidance lists required frontmatter and body sections.
- [x] Status values match the v1.4 PRD lifecycle.
- [x] The register template includes a convergence summary area for net-new, reopened, verified, deferred, and coverage-gap counts.

## Implementation Plan

- [x] Add the repo-wide register template.
- [x] Add findings folder guidance with a complete finding record shape.
- [x] Document stable ID and duplicate handling rules.
- [x] Document proof requirements for `verified` findings.
- [x] Cross-reference the prompt template's reconciliation requirements.

## Context Log

Read: router.md, AGENTS.md, E30-Detail.md, T001/T002 task files, v1.4-PRD.md (finding
lifecycle + audit file layout), templates/project/AGENTS.md, both savepoint-audit SKILL.md
copies.

Remediation: T001 was marked `status: done` with a `Pending` context log and no
deliverable — `templates/project/.savepoint/audit/prompt.md` did not exist. Per user
direction, created the missing prompt template (reconciliation rules, required per-finding
fields, coverage accounting, changelog), added the register-backed audit pointer to
`templates/project/AGENTS.md`, and aligned both `savepoint-audit` SKILL.md copies (template
+ live, kept in sync for the freshness tests).

T002 deliverables:
- `templates/project/.savepoint/audit/register.md` — convergence summary (net-new,
  reopened, verified, deferred, coverage gaps) + findings table (ID, title, status,
  severity, confidence, work item, first/last seen, proof).
- `templates/project/.savepoint/audit/findings/README.md` — `F###-slug.md` naming, stable
  ID + duplicate rules, required frontmatter + body sections, `verified` proof requirement,
  cross-reference to `../prompt.md` reconciliation.

All status values match the v1.4 PRD lifecycle (`open`, `triaged`, `mapped`, `in_progress`,
`fixed`, `verified`, `deferred`, `owner_decision`, `waived`, `duplicate`).

Quality gates: `make build && make test` pass (initial failure in
`TestScaffoldedSavepointSkillsMatchBundledSkills` / `TestProjectGuidanceTemplatesMirrorLiveGuidance`
fixed by syncing the live `agent-skills/savepoint-audit/SKILL.md`).

Drift note: T003 (init) and T004 (upgrade-assets + freshness) still own scaffolding the new
`.savepoint/audit/` files into generated projects — these templates are not yet wired into
init copy lists.
