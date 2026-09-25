---
id: E31-audit-register-data-model/T004-audit-register-validation
status: done
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

- [x] `verified` findings require named proof.
- [x] `duplicate` findings require an existing `duplicate_of` finding.
- [x] Deferred, waived, and owner-decision findings require rationale fields.
- [x] Task and defect references resolve against discovered work items.
- [x] Broken release, epic, task, defect, duplicate, and proof references return typed validation results suitable for doctor output.

## Implementation Plan

- [x] Add audit-register validation result types.
- [x] Build lookup maps for tasks, defects, and findings.
- [x] Validate lifecycle-specific required fields.
- [x] Validate work-item and duplicate references.
- [x] Add table-driven tests for every validation branch.

## Context Log

Read: this task file, E31-Detail.md, v1.4-PRD.md (finding lifecycle), T001
(`audit_finding.go`) for the finding fields + `DiagnoseFinding` seam, T003
(`audit.go`) for `AuditRegisterSet`, `dependency.go` for the short-ID resolution
idiom (`taskShortID`/`epicShortID`), `lifecycle.go` for the diagnostic-code +
message style, `discover.go` for the `ReleaseInfo/EpicInfo/TaskInfo/DefectInfo`
shapes, and the finding template README for the documented rationale contract.

Deliverables:
- `internal/data/audit_validate.go` — `AuditValidationCode` (10 codes),
  `AuditFindingValidation` (FindingID + Code + Message), `AuditWorkItems`
  (discovered release/epic/task/defect IDs), and `ValidateAuditFindings(findings,
  items)`. Builds membership sets for findings, releases, defects, tasks, and
  epics (plus short-ID sets for tasks/epics), then per finding checks the
  lifecycle-required field for its status and resolves every reference. Output is
  deterministic: finding order, lifecycle before references, references in
  release → epic → task → defect → duplicate order.
- `internal/data/audit_validate_test.go` — table-driven coverage for every
  branch (each lifecycle status missing/satisfied, each reference kind
  missing/resolved, full + short-ID resolution, duplicate-of unknown + self),
  plus a clean-register case, in-set duplicate resolution, and a multi-problem
  ordering test.

Decisions:
- **Rationale = deferral_reason OR waiver_reason** for deferred, owner_decision,
  and waived. The finding template documents `deferral_reason` as the rationale
  for all three; T001 added a dedicated `waiver_reason`. Accepting either avoids
  falsely flagging a finding that followed either convention, while each status
  keeps a distinct code/message. Resolves the T001 drift note between template
  and model. `owner_decision` has no dedicated field, so it relies on this.
- **Verified proof = `verified_proof` non-empty.** That field records the proof
  that satisfied `proof_needed`; the AC's "named proof" maps to it (the "proof
  references" half of AC5).
- **References = the plural link lists** (Releases/Epics/Tasks/Defects) +
  `duplicate_of`. Tasks/epics resolve by exact ID or short ID via the shared
  dependency helpers; releases/defects match exactly. The singular compat
  `work_item` is left unvalidated — it is a freeform pointer with no declared
  type, and the AC enumerates the typed lists. Empty IDs are dropped from the
  sets so an unset ref never resolves against an unset ID.
- **`duplicate_of` resolves against the finding set itself**, regardless of
  status, so a broken link on any finding is caught (AC5); a missing
  `duplicate_of` on a `duplicate` finding is the separate lifecycle check (AC2).
  A self-reference is treated as broken.
- **Validation, not healing.** These problems are not auto-fixable (unlike the
  T001/T002 load-time diagnostics), so this is a pure reporting pass that never
  mutates input — the doctor seam (E34) consumes the typed results.

Quality gates: `gofmt` clean, `go build ./...`, `go vet ./internal/data`, and
`make test` (all packages) pass. No existing files touched (AC: task/defect/
router/release-doc parsing unchanged).
