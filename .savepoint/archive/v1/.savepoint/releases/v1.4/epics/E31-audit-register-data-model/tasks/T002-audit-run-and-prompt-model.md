---
id: E31-audit-register-data-model/T002-audit-run-and-prompt-model
status: done
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

- [x] Prompt reads include path, availability, body, and prompt version when present.
- [x] Register reads include path, availability, body, and summary counts when present.
- [x] Run records parse date, label, auditor, prompt version, commit SHA, mode, coverage summary, source audits, and headline counts.
- [x] Missing prompt, register, findings, or runs return empty available state rather than fatal errors.
- [x] Malformed run frontmatter returns actionable validation errors.

## Implementation Plan

- [x] Define read models for audit prompt, register summary, and run records.
- [x] Reuse release-doc style availability handling where appropriate.
- [x] Parse run filenames and validate date/label shape.
- [x] Add tests for absent, empty, valid, and malformed prompt/register/run states.

## Context Log

Read: this T002 task file, E31-Detail.md, T001 (done) `audit_finding.go`/`_test.go` for the
parse/normalize/diagnose idioms, the E30 templates (`prompt.md`, `register.md`,
`runs/README.md`) for the exact on-disk shapes, `release_doc.go` for the
availability-handling pattern, `parser.go`/`write.go` (`SplitFrontmatterBody`,
`normalizeLineEndings`), and T003/T004 task files to fix the epic's internal boundaries.

Deliverables:
- `internal/data/audit_register.go` — `AuditPrompt` (path/available/body/version) +
  `LoadAuditPrompt`, `AuditRegister` (path/available/body + `RegisterSummary` and
  `HasSummary`) + `LoadAuditRegister`. Both reuse the `ReleaseDoc` absent-tolerance idiom:
  missing file → unavailable, unexpected read error → path-qualified error. Prompt version is
  parsed from the `## Changelog` section (highest `**vN**`); register summary counts are
  parsed best-effort from the convergence table. Also holds the shared `audit/` path
  constants for the domain.
- `internal/data/audit_run.go` — `AuditRun` read model + `(*Parser).ParseRunFile`,
  `ParseRunFilename` (`YYYY-MM-DD-label.md` shape), `AuditMode` constants, `DiagnoseRun`
  (with `RunDiagnostic`/`RunDiagnosticCode`), and the directory loaders `LoadAuditRuns` /
  `LoadAuditFindings`.
- `internal/data/audit_register_test.go`, `internal/data/audit_run_test.go` — table-driven
  coverage for absent, empty/no-summary, valid, malformed-frontmatter, and diagnose paths.

Decisions:
- **No healing for runs** (unlike findings). A run is immutable history, so `ParseRunFile`
  preserves the recorded `mode`/`date` verbatim; `DiagnoseRun` reports missing required
  fields, an invalid mode/date, a date that disagrees with the filename, and a filename that
  breaks the naming rule. This mirrors `DiagnoseFinding` and is the seam E34's doctor reads.
- **AC5 vs AC4 split:** absence (missing file/dir) is tolerated and returns an empty/
  unavailable state; *malformed content* (no frontmatter / bad YAML) returns an actionable,
  path-qualified error from `ParseRunFile` and aborts `LoadAuditRuns`/`LoadAuditFindings`.
- **`HasSummary`** distinguishes an absent convergence table from a genuine all-zero one, so
  callers don't read "0 net-new" from a register that simply has no summary yet.
- **Boundary with T003/T004:** sorting (findings by status/severity/ID, runs newest-first)
  and the single combined `LoadAuditRegisterSet`-style entry point belong to T003 discovery;
  cross-record validation (proof for `verified`, resolvable `duplicate_of`/work-item links)
  belongs to T004. T002 provides the per-artifact tolerant loaders T003 composes — the
  parallel to T001 leaving `DiagnoseFinding` as the seam. `LoadAuditFindings` lives here
  (not in T001's `audit_finding.go`) as the sibling of `LoadAuditRuns`, satisfying AC4's
  findings-absence clause without editing a done task's files.

Quality gates: `go build ./...`, `go vet ./internal/data`, `make test` (all packages) pass;
`gofmt` clean.
