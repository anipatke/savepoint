---
type: audit-findings
audited: 2026-06-16
---

# Audit Findings: E21 Document Template Optimisation

## Main Findings

E21's implementation is functionally complete against the task acceptance criteria. The revised `AGENTS.md`, `PRD.md`, and `Design.md` templates are present; the new `templates/project/.savepoint/Concept.md` exists with `type: project-concept`; and `internal/init/template_freshness_test.go` now covers the document frontmatter and AGENTS lifecycle terminology guards. The full quality gate `make build && make test` passes.

The documentation drift issue found during audit has been applied. `.savepoint/releases/v1.2/epics/E21-document-template-optimisation/E21-Detail.md` now correctly states that `savepoint init` scaffolds `Concept.md` automatically by walking embedded `templates/project`, while `upgrade-assets` still skips `.savepoint` project-state files for existing projects.

E21 is marked `audited`, `.savepoint/Design.md` records `last_audited: v1.2/E21-document-template-optimisation`, and the router now advances to `E30-epic-status-self-heal/T001-epic-status-normalization`.

## Code Style Review

- [x] One job per file - E21 changes stayed scoped to document templates, live AGENTS guidance, and template freshness tests.
- [x] One job per function - new test helpers were not introduced unnecessarily; assertions follow the existing local test style.
- [x] Test branches - freshness coverage now checks PRD, Design, Concept, and AGENTS lifecycle terminology.
- [x] Types document intent - the new table-driven frontmatter test uses explicit path/want fields.
- [x] Build only what is needed - no speculative product code was added.
- [x] Handle errors at boundaries - no boundary behavior changed; existing template read failures still fail tests directly.
- [x] One source of truth - E21-Detail now matches the actual init scaffold behavior for `Concept.md`.
- [x] Comments explain WHY - no low-value code comments were added.
- [x] Content in data files - template copy remains in markdown templates, not Go logic.
- [x] Small diffs - implementation is limited to the expected template and freshness-test surfaces.

## Proposed Changes

### Target File
.savepoint/releases/v1.2/epics/E21-document-template-optimisation/E21-Detail.md

### Replace
```md
**Out of scope:**
- Changes to agent skills (`agent-skills/` or `templates/project/agent-skills/`).
- Changes to `router.md`, `config.yml`, or `visual-identity.md` templates.
- Live project `PRD.md` and `Design.md` (project-state files, not templates).
- Adding `Concept.md` to the init command scaffold (deferred — tracked below).

## Open decisions

- **Init scaffold inclusion:** `Concept.md` is authored here but not wired into `savepoint init` output. A follow-on task or epic should add it to the init scaffold and the `upgrade-assets` skip-list once the template has been used in at least one real project.
```

### With
```md
**Out of scope:**
- Changes to agent skills (`agent-skills/` or `templates/project/agent-skills/`).
- Changes to `router.md`, `config.yml`, or `visual-identity.md` templates.
- Live project `PRD.md` and `Design.md` (project-state files, not templates).
- Adding `Concept.md` to `upgrade-assets` for existing projects (deferred because `.savepoint/` project-state files are intentionally skipped).

## Open decisions

- **Existing-project upgrade inclusion:** `savepoint init` now scaffolds `Concept.md` automatically because it walks embedded `templates/project`. Existing Savepoint projects do not receive `.savepoint/Concept.md` through `upgrade-assets`; decide after real use whether that command should create the file or continue treating it as project state.
```
