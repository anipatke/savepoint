---
type: audit-findings
audited: 2026-06-29
status: applied
---

# Audit Findings: E30 Audit Register Templates

## Main Findings

E30 is substantially implemented: the prompt, register, findings guidance, run history guidance, scaffold coverage, upgrade-assets behavior, and template freshness coverage are present. I verified the scoped implementation and tests against the task acceptance criteria, and the quality gates pass: `make build`, `make test`, and `go test ./internal/init`.

Finding 1, medium severity — applied: `templates/project/.savepoint/audit/register.md` previously seeded a real-looking `F001` row despite a zero-count convergence summary, conflicting with the stable, non-reusable `F###` contract. The placeholder row is now replaced with a non-row instruction so the first real audit mints `F001` cleanly, and a template freshness assertion guards against the example row returning.

Finding 2, low severity — applied: E30 workflow metadata was stale. T001 evidence is now reconciled (acceptance criteria and implementation plan checked, `Context Log` filled in). Router advancement is handled by this apply/close pass below.

Coverage examined: E30 detail and all four E30 task files; `.savepoint/Design.md`; root `AGENTS.md`; v1.4 finding lifecycle PRD section; audit templates under `templates/project/.savepoint/audit/`; generated project `AGENTS.md`; both bundled `savepoint-audit` skill copies; `internal/init` scaffold, upgrade, integration, and template freshness code/tests; `main.go` embed directives. Coverage skipped: unrelated dirty launcher/config files outside the E30 task context.

## Code Style Review

- [x] One job per file: audit markdown templates, scaffold tests, upgrade tests, and upgrade helper code stay in their expected modules.
- [x] One job per function: the new `isAuditAsset` and `upgradeAuditAsset` helpers are small and single-purpose.
- [x] Test branches: missing assets, edited existing assets, pristine reruns, dry-run behavior, scaffold creation, and real-template upgrade coverage are tested.
- [x] Types document intent: existing `UpgradeAction` reporting and task frontmatter remain explicit.
- [x] Build only what is needed: no parser, board, doctor, or runtime audit-register behavior was added in this epic.
- [x] Handle errors at boundaries: upgrade asset reads, parent directory creation, writes, and stat errors are wrapped at the filesystem boundary.
- [x] One source of truth: the seeded `F001` example row was removed, so the register's zero-count current state and stable-ID rules no longer conflict.
- [x] Comments explain WHY: new comments explain why audit assets are additive-only user-maintained state.
- [x] Content in data files: audit workflow copy lives in templates and skill markdown, not Go logic.
- [x] Small diffs: source changes are narrow and supported by focused tests.

## Proposed Changes

### Target File
templates/project/.savepoint/audit/register.md

### Replace
```md
| F001 | Example finding (replace) | open | medium | medium | E00-example/T000 | 0000-00-00 | 0000-00-00 | Pending |
```

### With
```md
<!-- Add the first real finding as `F001`; do not keep placeholder rows in the register. -->
```

### Target File
internal/init/template_freshness_test.go

### Replace
```go
	assertContains(t, register, "Convergence summary")
```

### With
```go
	assertContains(t, register, "Convergence summary")
	assertContains(t, register, "first real finding as `F001`")
	assertNotContains(t, register, "Example finding")
```

### Target File
.savepoint/releases/v1.4/epics/E30-audit-register-templates/tasks/T001-audit-prompt-template.md

### Replace
```md
## Acceptance Criteria

- [ ] The prompt template explains that each audit must reconcile against the existing register when present.
- [ ] Required per-finding fields include stable ID handling, severity, confidence, source auditor, location, guardrail IDs, proof needed, and work-item mapping.
- [ ] The prompt requires coverage notes for examined and unexamined surfaces.
- [ ] The prompt includes a short changelog section for refinements over time.
- [ ] Generated project guidance points agents to the prompt before starting a register-backed audit.

## Implementation Plan

- [ ] Add `.savepoint/audit/prompt.md` to the project template.
- [ ] Define the required audit output shape in the prompt body.
- [ ] Add a changelog section with an initial version entry.
- [ ] Update generated `AGENTS.md` guidance to mention the audit prompt when the register workflow is used.
- [ ] Align the existing audit skill wording with the new prompt location.

## Context Log

Pending.
```

### With
```md
## Acceptance Criteria

- [x] The prompt template explains that each audit must reconcile against the existing register when present.
- [x] Required per-finding fields include stable ID handling, severity, confidence, source auditor, location, guardrail IDs, proof needed, and work-item mapping.
- [x] The prompt requires coverage notes for examined and unexamined surfaces.
- [x] The prompt includes a short changelog section for refinements over time.
- [x] Generated project guidance points agents to the prompt before starting a register-backed audit.

## Implementation Plan

- [x] Add `.savepoint/audit/prompt.md` to the project template.
- [x] Define the required audit output shape in the prompt body.
- [x] Add a changelog section with an initial version entry.
- [x] Update generated `AGENTS.md` guidance to mention the audit prompt when the register workflow is used.
- [x] Align the existing audit skill wording with the new prompt location.

## Context Log

Read: router.md, AGENTS.md, E30-Detail.md, this T001 task file,
`templates/project/.savepoint/audit/prompt.md`, `templates/project/AGENTS.md`,
`agent-skills/savepoint-audit/SKILL.md`, and the scaffolded audit skill copy.

Delivered by later E30 remediation: added the prompt template, required reconciliation and
per-finding fields, coverage accounting, changelog, generated-project AGENTS guidance, and
audit-skill wording that points register-backed audits at `.savepoint/audit/prompt.md`.

Audit verification on 2026-06-29 confirmed the T001 acceptance criteria are present in the
scoped files. E30 quality gates pass: `make build`, `make test`, and `go test ./internal/init`.
```
