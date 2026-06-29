---
id: E31-audit-register-data-model/T004-audit-register-validation
status: planned
objective: Validate audit-register consistency against tasks, defects, and proof rules.
depends_on:
  - E31-audit-register-data-model/T003-audit-register-discovery
complexity_tier: high
complexity_reason: Cross-record validation touches audit, task, defect, and dependency semantics.
---

# T004: Audit Register Validation

## Problem

The register can only be trusted if verified findings have proof, duplicate links resolve, and mapped work items exist.

## Context Files

- `internal/data/audit_validate.go`
- `internal/data/audit_validate_test.go`
- `internal/data/audit.go`
- `internal/data/task.go`
- `internal/data/defect.go`
- `internal/data/dependency.go`

## Acceptance Criteria

- [ ] `verified` findings require named proof.
- [ ] `duplicate` findings require an existing `duplicate_of` finding.
- [ ] Deferred, waived, and owner-decision findings require rationale fields.
- [ ] Task and defect references resolve against discovered work items.
- [ ] Broken release, epic, task, defect, duplicate, and proof references return typed validation results suitable for doctor output.

## Implementation Plan

- [ ] Add audit-register validation result types.
- [ ] Build lookup maps for tasks, defects, and findings.
- [ ] Validate lifecycle-specific required fields.
- [ ] Validate work-item and duplicate references.
- [ ] Add table-driven tests for every validation branch.

## Context Log

Pending.
