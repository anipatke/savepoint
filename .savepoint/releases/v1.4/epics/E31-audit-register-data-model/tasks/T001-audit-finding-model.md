---
id: E31-audit-register-data-model/T001-audit-finding-model
status: done
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

- [x] Finding files parse from `F###-slug.md` markdown with frontmatter and body sections.
- [x] Supported statuses match the v1.4 PRD.
- [x] Required fields include ID, title, status, severity, confidence, proof needed, first seen, and last seen.
- [x] Optional fields support releases, epics, tasks, defects, guardrail IDs, source locations, duplicate-of, deferred rationale, waiver rationale, and verified proof.
- [x] Invalid IDs, statuses, severity, confidence, and missing required fields heal at load and are reported as actionable warnings via `DiagnoseFinding`.

## Implementation Plan

- [x] Define audit finding types and constants in `internal/data`.
- [x] Reuse the existing frontmatter split/parser helpers.
- [x] Add filename and frontmatter ID consistency checks.
- [x] Add table-driven tests for valid, missing, malformed, and optional-field records.
- [x] Keep parser behavior independent from board and doctor packages.

## Context Log

Read: this task file, E31-Detail.md, E30-audit-register-templates/T002 task + the
scaffolded `findings/README.md`, v1.4-PRD.md (finding lifecycle + file layout), and the
existing `internal/data` models (defect.go, parser.go, write.go, errors.go, release_doc.go)
for the parsing/validation idioms.

Deliverables:
- `internal/data/audit_finding.go` — `AuditFinding` type, `FindingStatus` (10 PRD lifecycle
  values) and `FindingConfidence` (high/medium/low) constants,
  `(*Parser).ParseFindingFile(path, content)`, `NormalizeFindingForLoad`, and
  `DiagnoseFinding` (with `FindingDiagnostic` / `FindingDiagnosticCode`). Reuses
  `SplitFrontmatterBody` / `normalizeLineEndings`, unmarshals directly into the typed struct
  (yaml tags), and retains the markdown body.
- `internal/data/audit_finding_test.go` — table-driven coverage for valid, all-statuses,
  optional-field, missing-frontmatter, malformed-YAML, enum healing, filename ID recovery
  (mismatch/invalid/slug-only), and `DiagnoseFinding` missing-required + invalid-value cases.

Decisions:
- **Healing with warnings, not errors** (per user direction, matching the defect/epic idiom):
  `ParseFindingFile` only returns an error for a structural failure (missing frontmatter,
  malformed YAML). `NormalizeFindingForLoad` heals in place — non-canonical status→`open`,
  severity→`medium`, confidence→`medium`, and a missing/malformed/mismatched frontmatter ID
  is recovered from the filename's `F###`. `DiagnoseFinding` re-derives every healed/missing
  condition from raw frontmatter as actionable warnings, mirroring `DiagnoseDefectLifecycle` /
  `DiagnoseEpicStatus` so E34's doctor can surface them from raw frontmatter. (Earlier draft
  returned hard errors; replaced.)
- Filename `F###` is treated as the authoritative stable ID (the template defines the ID as
  "matches the filename"); the frontmatter ID heals toward it when they disagree.
- `Severity` reuses the shared `DefectSeverity` vocabulary (critical/high/medium/low) to keep
  one source of truth for the severity set.
- T002/T003/T004 scope kept out: no `.savepoint/audit/` discovery, no run/prompt model, no
  work-item link resolution, and no doctor wiring (E34) here — `DiagnoseFinding` is the seam.

Drift note: the E30 finding template (`findings/README.md`) documents a single `work_item`
field plus `deferral_reason`, and omits separate `releases/epics/tasks/defects`,
`waiver_reason`, and `verified_proof`. The T001 AC enumerates those richer optional fields,
so the model supports both: the plural link lists per the AC and a compatibility `work_item`
(+ `source_auditor`) so scaffolded findings still parse without loss. The template may want a
follow-up to expose the richer fields; flagging for E30/E33 owners.

Quality gates: `go build ./...`, `go vet ./internal/data`, and `make test` (all packages)
pass. `gofmt` applied.
