---
id: E31-audit-register-data-model/T001-audit-finding-model
status: planned
objective: Add typed audit finding parsing and validation primitives.
depends_on:
  - E30-audit-register-templates/T002-register-and-finding-templates
complexity_tier: medium
complexity_reason: The model adds a new lifecycle surface while reusing existing frontmatter parsing.
---

# T001: Audit Finding Model

## Problem

Audit findings need typed fields and validation before the board or doctor can treat them as durable project state.

## Context Files

- `internal/data/audit_finding.go`
- `internal/data/audit_finding_test.go`
- `internal/data/parser.go`
- `internal/data/errors.go`
- `internal/data/lifecycle.go`

## Acceptance Criteria

- [ ] Finding files parse from `F###-slug.md` markdown with frontmatter and body sections.
- [ ] Supported statuses match the v1.4 PRD.
- [ ] Required fields include ID, title, status, severity, confidence, proof needed, first seen, and last seen.
- [ ] Optional fields support releases, epics, tasks, defects, guardrail IDs, source locations, duplicate-of, deferred rationale, waiver rationale, and verified proof.
- [ ] Invalid IDs, statuses, severity, confidence, and missing required fields return actionable errors.

## Implementation Plan

- [ ] Define audit finding types and constants in `internal/data`.
- [ ] Reuse the existing frontmatter split/parser helpers.
- [ ] Add filename and frontmatter ID consistency checks.
- [ ] Add table-driven tests for valid, missing, malformed, and optional-field records.
- [ ] Keep parser behavior independent from board and doctor packages.

## Context Log

Pending.
