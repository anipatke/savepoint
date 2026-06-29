---
id: E30-audit-register-templates/T003-run-history-template-and-scaffold
status: planned
objective: Add audit run history guidance and wire audit-register assets into scaffolding.
depends_on:
  - E30-audit-register-templates/T001-audit-prompt-template
  - E30-audit-register-templates/T002-register-and-finding-templates
complexity_tier: medium
complexity_reason: Scaffold wiring touches generated assets and existing init behavior.
---

# T003: Run History Template and Scaffold

## Problem

Audit run history needs an append-only location and the new audit assets must be created for new projects without affecting existing scaffolding behavior.

## Context Files

- `templates/project/.savepoint/audit/runs/README.md`
- `templates/project/.savepoint/audit/register.md`
- `internal/init/scaffold.go`
- `internal/init/scaffold_test.go`
- `internal/init/integration_test.go`

## Acceptance Criteria

- [ ] Run history guidance defines the `YYYY-MM-DD-label.md` naming convention.
- [ ] Run records require date, auditor/model, prompt version, commit SHA, mode, coverage, source audits, and headline counts.
- [ ] New project scaffolding creates the audit prompt, register, findings guidance, and runs guidance.
- [ ] Existing scaffold behavior for router, PRD, Design, AGENTS, and config remains unchanged.
- [ ] Tests cover the generated audit-register paths.

## Implementation Plan

- [ ] Add run history guidance under the project template.
- [ ] Wire the audit-register files into scaffold creation.
- [ ] Add scaffold tests for the new files and directories.
- [ ] Add integration coverage proving existing generated files are unchanged.

## Context Log

Pending.
