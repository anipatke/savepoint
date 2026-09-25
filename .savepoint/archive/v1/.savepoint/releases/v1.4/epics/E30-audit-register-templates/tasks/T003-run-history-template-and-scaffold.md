---
id: E30-audit-register-templates/T003-run-history-template-and-scaffold
status: done
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

- [x] Run history guidance defines the `YYYY-MM-DD-label.md` naming convention.
- [x] Run records require date, auditor/model, prompt version, commit SHA, mode, coverage, source audits, and headline counts.
- [x] New project scaffolding creates the audit prompt, register, findings guidance, and runs guidance.
- [x] Existing scaffold behavior for router, PRD, Design, AGENTS, and config remains unchanged.
- [x] Tests cover the generated audit-register paths.

## Implementation Plan

- [x] Add run history guidance under the project template.
- [x] Wire the audit-register files into scaffold creation.
- [x] Add scaffold tests for the new files and directories.
- [x] Add integration coverage proving existing generated files are unchanged.

## Context Log

Read: T003 task, T001/T002 (done) for register/findings/prompt templates, v1.4-PRD.md
(audit file layout + run record requirements), E30-Detail.md, scaffold.go, scaffold_test.go,
integration_test.go, template_freshness_test.go, main.go (embed directives + `fs.Sub`).

Key finding: scaffolding is fully generic — `Scaffold` walks the embedded `templates/project`
subtree via `fs.WalkDir`, and `main.go` embeds `.savepoint` with `all:` so hidden/nested
files are included. No allowlist exists, so "wiring" the audit assets into scaffold creation
just means adding the template file; no `scaffold.go` change was needed (or made).

Deliverable:
- `templates/project/.savepoint/audit/runs/README.md` — append-only run history guidance:
  `YYYY-MM-DD-label.md` naming, immutability rule, required frontmatter (date, auditor,
  model, prompt_version, commit, mode, coverage, source_audits, headline counts net_new/
  reopened/verified/deferred/coverage_gaps), required body sections (Scope, Coverage,
  Findings, Reconciliation), and the link back to register convergence totals.

Tests added:
- `TestProjectAuditRegisterTemplatesExist` (template_freshness_test.go) — validates the real
  embedded prompt/register/findings/runs templates, including the runs naming convention and
  every required run-record field.
- `TestScaffold_createsAuditRegisterAssets` (scaffold_test.go) — proves the four audit assets
  (incl. nested `runs/` and `findings/` dirs) are created.
- `TestIntegration_AuditAssetsDoNotAlterExistingFiles` (integration_test.go) — proves audit
  assets are scaffolded while config/PRD/Design/router content and the merged AGENTS.md stay
  unchanged. Also extended `runInitPipeline` MapFS + `TestIntegration_EmptyDirectory` entries
  to cover the audit paths.

Verification: `make build && make test` pass. End-to-end `savepoint init` in a temp dir emits
all four `.savepoint/audit/` assets including `runs/README.md` from the real embedded FS.

Note for T004: upgrade-assets wiring + template freshness for the audit assets is owned by
E30 T004; this task covers init scaffolding only.
