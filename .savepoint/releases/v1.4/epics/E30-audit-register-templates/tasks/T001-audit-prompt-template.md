---
id: E30-audit-register-templates/T001-audit-prompt-template
status: done
objective: Add the scaffolded reusable audit prompt template.
depends_on: []
complexity_tier: low
complexity_reason: The change is a bounded markdown template with focused scaffold expectations.
---

# T001: Audit Prompt Template

## Problem

Savepoint audits do not have a reusable prompt that requires reconciliation against prior findings, coverage accounting, or closure proof.

## Context Files

- `templates/project/.savepoint/audit/prompt.md`
- `templates/project/AGENTS.md`
- `agent-skills/savepoint-audit/SKILL.md`

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
