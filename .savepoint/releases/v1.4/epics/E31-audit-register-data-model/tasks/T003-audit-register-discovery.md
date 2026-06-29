---
id: E31-audit-register-data-model/T003-audit-register-discovery
status: planned
objective: Discover repo-wide audit register files from `.savepoint/audit/`.
depends_on:
  - E31-audit-register-data-model/T001-audit-finding-model
  - E31-audit-register-data-model/T002-audit-run-and-prompt-model
complexity_tier: medium
complexity_reason: Discovery must join several file types while preserving existing release traversal.
---

# T003: Audit Register Discovery

## Problem

Callers need one data entry point that loads prompt, register, findings, and runs without duplicating filesystem traversal.

## Context Files

- `internal/data/audit.go`
- `internal/data/audit_finding.go`
- `internal/data/audit_run.go`
- `internal/data/discover.go`
- `internal/data/discover_test.go`

## Acceptance Criteria

- [ ] A single data helper loads audit prompt, register summary, sorted findings, and sorted runs for a Savepoint root.
- [ ] Findings sort by status priority, severity priority, and ID.
- [ ] Runs sort newest-first by date and label.
- [ ] Projects without `.savepoint/audit/` return an empty audit register without error.
- [ ] Existing task, epic, release, and defect discovery behavior remains unchanged.

## Implementation Plan

- [ ] Add a repo-wide audit discovery helper.
- [ ] Define deterministic finding and run sort orders.
- [ ] Reuse existing root-relative path handling.
- [ ] Add discovery tests for absent, partial, complete, and malformed audit trees.

## Context Log

Pending.
