---
id: E30-epic-status-self-heal/T002-epic-status-doctor-diagnostic
status: done
objective: Surface a non-canonical epic status through savepoint doctor with an actionable repair hint
depends_on: [E30-epic-status-self-heal/T001-epic-status-normalization]
complexity_tier: low
complexity_reason: "Adds one diagnostic call + repair string following the established defect-lifecycle doctor pattern"
---

# T002: Epic Status Doctor Diagnostic

## Problem

Once T001 heals a non-canonical epic status at load, the underlying file still
carries the wrong value silently. Tasks and defects already surface their
load-time heals through `savepoint doctor` (e.g. `DiagnoseDefectLifecycle` +
repair hints in `internal/doctor/repairs.go`), but epic detail files are only
checked for readable frontmatter (`checkEpicDetail`, `internal/doctor/checks.go`).
A healed-but-wrong epic status is therefore never reported, so it never gets
corrected at the source.

## Context Files

- `internal/data/lifecycle.go`
- `internal/data/lifecycle_test.go`
- `internal/doctor/checks.go`
- `internal/doctor/checks_test.go`
- `internal/doctor/repairs.go`
- `internal/doctor/repairs_test.go`

## Acceptance Criteria

- [x] `DiagnoseEpicStatus(status)` in `internal/data` returns a diagnostic when
      the raw status is non-canonical (alias or unknown) and nothing when it is
      canonical, mirroring `DiagnoseDefectLifecycle`'s shape.
- [x] The diagnostic message names the offending value and the canonical set and
      states what it loads as (e.g. `loads as planned`).
- [x] `checkEpicDetail` parses the epic `status` and emits a `doctor` Problem for
      a non-canonical value, without breaking on a missing/empty status.
- [x] `internal/doctor/repairs.go` maps the diagnostic to an actionable hint
      ("Set the epic status to planned, in_progress, done, or audited").
- [x] Tests cover: canonical status → no problem; unknown status → one problem
      with the repair hint; the repair mapping itself.
- [x] `make build && make test` passes.

## Implementation Plan

- [x] In `internal/data/lifecycle.go`, add an epic-status diagnostic type/code
      (mirror `DefectLifecycleDiagnostic` / `DefectLifecycleDiagnosticCode`) and
      `DiagnoseEpicStatus(status string) []EpicStatusDiagnostic` that reports when
      `ResolveEpicStatusAlias`/`IsCanonicalEpicStatus` (from T001) classify the
      value as non-canonical.
- [x] In `internal/doctor/checks.go`, extend `checkEpicDetail` to read the parsed
      `status` field and append a `Problem{File: detailPath, Message: ...}` for
      each `DiagnoseEpicStatus` result. Keep existing frontmatter validation.
- [x] In `internal/doctor/repairs.go`, add a `strings.Contains` branch for the
      epic-status message returning the canonical-set repair hint, matching the
      existing defect-stage branches.
- [x] Add tests in `lifecycle_test.go`, `checks_test.go`, and `repairs_test.go`.
- [x] Run `make build && make test`.

## Notes

- This task is reporting only; it must not rewrite epic files. The board already
  renders correctly after T001; doctor exists so the source file gets fixed.
- Depends on T001 for the canonical vocabulary and classification helpers.

## Context Log

**Read:** `internal/data/lifecycle.go`, `internal/doctor/checks.go`,
`internal/doctor/repairs.go`, `internal/doctor/checks_test.go`,
`internal/doctor/repairs_test.go`, `internal/testutil/fixture.go`.

**Edited:**

- `internal/data/lifecycle.go` — added `EpicStatusDiagnostic` /
  `EpicStatusDiagnosticCode` (`EpicStatusAliasCode`, `EpicStatusInvalidCode`)
  and `DiagnoseEpicStatus(EpicStatus) []EpicStatusDiagnostic`, mirroring
  `DiagnoseDefectLifecycle`. Canonical and empty statuses return no diagnostics;
  aliases and unknown values each return one diagnostic whose message names the
  offending value, the canonical set, and what it loads as.
- `internal/doctor/checks.go` — `checkEpicDetail` now keeps frontmatter
  validation and additionally calls a new `checkEpicStatus` helper that parses
  the raw `status` and appends a Problem per `DiagnoseEpicStatus` result. A
  missing/empty status produces no problem.
- `internal/doctor/repairs.go` — added a branch (placed before the generic
  `release` branch so the file path in `Problem.Error()` cannot shadow it) that
  maps both epic-status messages to "Set the epic status to planned,
  in_progress, done, or audited".
- `internal/data/lifecycle_test.go` — `DiagnoseEpicStatus` canonical/empty →
  none, alias → one `status_alias`, unknown → one `invalid_status`, with message
  content assertions.
- `internal/doctor/checks_test.go` — canonical epic status → no status problem;
  non-canonical → one problem that maps to the epic-status repair hint.
- `internal/doctor/repairs_test.go` — added both epic-status message → hint
  cases.

**Quality gates:** `make build && make test` — all packages pass. (One
intermediate failure: the repair branch initially sat after the generic
`release` branch, which matched the `releases/` path in the Problem file; moved
it up beside the defect branches.)
