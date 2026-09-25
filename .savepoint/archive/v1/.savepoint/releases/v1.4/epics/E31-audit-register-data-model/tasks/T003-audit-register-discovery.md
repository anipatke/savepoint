---
id: E31-audit-register-data-model/T003-audit-register-discovery
status: done
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

- [x] A single data helper loads audit prompt, register summary, sorted findings, and sorted runs for a Savepoint root.
- [x] Findings sort by status priority, severity priority, and ID.
- [x] Runs sort newest-first by date and label.
- [x] Projects without `.savepoint/audit/` return an empty audit register without error.
- [x] Existing task, epic, release, and defect discovery behavior remains unchanged.

## Implementation Plan

- [x] Add a repo-wide audit discovery helper.
- [x] Define deterministic finding and run sort orders.
- [x] Reuse existing root-relative path handling.
- [x] Add discovery tests for absent, partial, complete, and malformed audit trees.

## Context Log

Read: this task file, E31-Detail.md, T001 (`audit_finding.go`) + T002 (`audit_run.go`,
`audit_register.go`) done outputs for the per-artifact loaders and diagnostic seams, and
`discover.go`/`discover_test.go` for the existing traversal idioms this must not disturb.

Deliverables:
- `internal/data/audit.go` — `AuditRegisterSet` (root, prompt, register, sorted findings,
  sorted runs) and `LoadAuditRegisterSet(root)`, the single audit entry point. It composes
  T002's tolerant `LoadAuditPrompt`/`LoadAuditRegister`/`LoadAuditFindings`/`LoadAuditRuns`
  rather than re-walking `audit/`, then applies the deterministic orders. Exposes
  `SortFindings` (status → severity → ID) and `SortRuns` (date desc, label asc tiebreak) for
  reuse by future board/doctor work, plus the `findingStatusPriority`/`findingSeverityPriority`
  rank helpers.
- `internal/data/audit_test.go` — `LoadAuditRegisterSet` coverage for absent, partial
  (findings only), complete (README files skipped), and malformed-aborts trees, plus direct
  `SortFindings`/`SortRuns` order tests with scrambled inputs.

Decisions:
- **Compose, don't re-traverse.** AC4 (absent-tolerant) falls out of the per-artifact loaders'
  existing missing-file handling, so the set helper adds only ordering. A malformed finding/run
  still aborts the whole load with the underlying path-qualified error (matches T002's AC4/AC5
  split). No early audit/-dir existence check — the four loaders already short-circuit absence.
- **Status priority = canonical lifecycle order.** `findingStatusPriority` indexes into
  `findingStatuses`, keeping that slice the single source of truth; active states (open…) sort
  above resolved ones (waived, duplicate). Severity is critical→low. Both rank unknowns last,
  though load-time healing means a loaded finding never carries one.
- **Run tiebreak = label ascending.** "Newest-first by date and label" reads as date the
  primary descending key; labels aren't temporal, so an ascending alphabetical label is the
  deterministic secondary. `SliceStable` keeps it reproducible.
- **Naming:** used `LoadAuditRegisterSet`, the name T002's context log reserved for this
  combined entry point. Cross-record validation (proof/`duplicate_of`/work-item links) stays
  out — that is T004.

Quality gates: `gofmt` clean, `go build ./...`, `go vet ./internal/data`, and `make test`
(all packages) pass. Existing `discover.go`/`discover_test.go` untouched (AC5).
