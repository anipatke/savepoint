---
id: E31-audit-register-data-model/T002-audit-run-and-prompt-model
status: planned
objective: Add typed audit run, prompt, and register summary records.
depends_on:
  - E30-audit-register-templates/T003-run-history-template-and-scaffold
complexity_tier: medium
complexity_reason: The task adds multiple related read models with tolerant missing-file behavior.
---

# T002: Audit Run and Prompt Model

## Problem

Audit prompt, register summary, and run history files need structured read models for the TUI and doctor.

## Context Files

- `internal/data/audit_run.go`
- `internal/data/audit_run_test.go`
- `internal/data/audit_register.go`
- `internal/data/audit_register_test.go`
- `internal/data/parser.go`
- `internal/data/release_doc.go`

## Acceptance Criteria

- [ ] Prompt reads include path, availability, body, and prompt version when present.
- [ ] Register reads include path, availability, body, and summary counts when present.
- [ ] Run records parse date, label, auditor, prompt version, commit SHA, mode, coverage summary, source audits, and headline counts.
- [ ] Missing prompt, register, findings, or runs return empty available state rather than fatal errors.
- [ ] Malformed run frontmatter returns actionable validation errors.

## Implementation Plan

- [ ] Define read models for audit prompt, register summary, and run records.
- [ ] Reuse release-doc style availability handling where appropriate.
- [ ] Parse run filenames and validate date/label shape.
- [ ] Add tests for absent, empty, valid, and malformed prompt/register/run states.

## Context Log

Pending.
