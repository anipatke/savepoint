---
id: E30-epic-status-self-heal/T002-epic-status-doctor-diagnostic
status: planned
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

- [ ] `DiagnoseEpicStatus(status)` in `internal/data` returns a diagnostic when
      the raw status is non-canonical (alias or unknown) and nothing when it is
      canonical, mirroring `DiagnoseDefectLifecycle`'s shape.
- [ ] The diagnostic message names the offending value and the canonical set and
      states what it loads as (e.g. `loads as planned`).
- [ ] `checkEpicDetail` parses the epic `status` and emits a `doctor` Problem for
      a non-canonical value, without breaking on a missing/empty status.
- [ ] `internal/doctor/repairs.go` maps the diagnostic to an actionable hint
      ("Set the epic status to planned, in_progress, done, or audited").
- [ ] Tests cover: canonical status → no problem; unknown status → one problem
      with the repair hint; the repair mapping itself.
- [ ] `make build && make test` passes.

## Implementation Plan

- [ ] In `internal/data/lifecycle.go`, add an epic-status diagnostic type/code
      (mirror `DefectLifecycleDiagnostic` / `DefectLifecycleDiagnosticCode`) and
      `DiagnoseEpicStatus(status string) []EpicStatusDiagnostic` that reports when
      `ResolveEpicStatusAlias`/`IsCanonicalEpicStatus` (from T001) classify the
      value as non-canonical.
- [ ] In `internal/doctor/checks.go`, extend `checkEpicDetail` to read the parsed
      `status` field and append a `Problem{File: detailPath, Message: ...}` for
      each `DiagnoseEpicStatus` result. Keep existing frontmatter validation.
- [ ] In `internal/doctor/repairs.go`, add a `strings.Contains` branch for the
      epic-status message returning the canonical-set repair hint, matching the
      existing defect-stage branches.
- [ ] Add tests in `lifecycle_test.go`, `checks_test.go`, and `repairs_test.go`.
- [ ] Run `make build && make test`.

## Notes

- This task is reporting only; it must not rewrite epic files. The board already
  renders correctly after T001; doctor exists so the source file gets fixed.
- Depends on T001 for the canonical vocabulary and classification helpers.
