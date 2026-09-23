---
id: T-022
title: Delete code the shipped binary cannot reach
objective: O-021
planned_by: {role: planner, session: o021-design-20260923}
status: planned
complexity_tier: medium
complexity_reason: "Mechanical deletion across four packages, but tests entangle live and dead code, so each test file must be split or removed without losing coverage of anything still reachable."
depends_on: []
owner_validation:
    required: false
    accepted_check: ""
---

# T-022: Delete code the shipped binary cannot reach

## Outcome

Every function that `deadcode` reports as unreachable from the `savepoint`
main package is gone, together with the tests that only exercise it. The
V1 board disappears, leaving `internal/board` as the thin `board` command
dispatcher. User-visible behavior does not change.

## User Check

None needed beyond the gate: `savepoint board`, `resume`, `doctor`, `init`,
`upgrade-assets`, and `migrate --dry-run` behave exactly as before on this
repository and on a V1 fixture copy.

## Done When

- `go run golang.org/x/tools/cmd/deadcode@latest .` reports nothing under
  `internal/board`, `internal/data`, `internal/doctor`, `internal/init`, or
  `cmd`. Anything reported under `internal/migrate` that T-025 will delete
  anyway (`apply.go`, `cutover.go`, `operation.go`, `manifest.go` entries)
  may remain and is listed in evidence.
- `internal/board` contains only `board.go` (dispatch and `Filters`) and
  whatever it genuinely calls. Every V1 board file and V1 board test is
  deleted. `internal/board/v2` reachable behavior is untouched.
- The unreachable V1 functions in `internal/data` (audit register/run
  loaders and validators, V1 lifecycle helpers, V1 write helpers, unused
  parser/discover methods) and the unused V2 writers (`CreateCheckV2`,
  `CreateIssueV2`, `WriteIssueV2`, `WriteReleaseV2`,
  `WriteReleaseLifecycleV2`, `WriteReleaseEvidenceV2`, and their private
  helpers) are deleted. V1 readers that `internal/migrate` still reaches stay.
- Unreachable V1 checks in `internal/doctor` are deleted; `RunV2Checks`
  output on this repository is byte-identical before and after.
- Tests that only covered deleted code are deleted; tests covering both
  keep their live cases. No test is weakened to make deletion pass.
- Evidence records production and test line counts per package before
  and after, and the before/after `deadcode` reports.
- `git diff --check`, `make build && make test-fast` pass.

## Context Files

Deletion set: the `deadcode` report, run fresh at start. Board:
`internal/board/board.go`, `internal/board/debug.go`,
`internal/board/legacy.go`, `internal/board/model.go`,
`internal/board/update.go`, `internal/board/view.go`,
`internal/board/dispatch_test.go`, `internal/board/board_test.go`. Data:
`internal/data/audit.go`, `internal/data/audit_backlinks.go`,
`internal/data/audit_finding.go`, `internal/data/audit_register.go`,
`internal/data/audit_run.go`, `internal/data/audit_validate.go`,
`internal/data/lifecycle.go`, `internal/data/write.go`,
`internal/data/write_test.go`, `internal/data/release_gate_v2.go`,
`internal/data/release_doc.go`, `internal/data/parser.go`,
`internal/data/discover.go`, `internal/data/task.go`. Doctor:
`internal/doctor/checks.go`, `internal/doctor/checks_test.go`,
`internal/doctor/interfaces.go`, `internal/doctor/repairs.go`,
`internal/doctor/report.go`, `internal/doctor/gates.go`. Other:
`internal/init/clipboard.go`, `internal/init/upgrade.go`, `cmd/init.go`,
`internal/board/v2/card.go`, `internal/board/v2/model.go`,
`internal/board/v2/objectives.go`, `internal/board/v2/releases.go`,
`internal/board/v2/update.go`, `internal/board/v2/view.go`. Any other
`internal/board/*.go` file is deleted whole once `board.go` no longer
references it; its `_test.go` sibling goes with it.

## Design References

Design sections 2 (Codebase Map) and 11 (testing).

## Guardrails

ARCH-01, ARCH-04, DATA-01, TEST-01, TEST-06, TEST-08.

## Implementation Plan

1. Record baseline line counts and a fresh `deadcode` report.
2. Delete the V1 board: remove every `internal/board` file except
   `board.go`, then restore only what `board.go` needs. Move any board
   test that exercises the live dispatch into `board_test.go`.
3. Delete unreachable `data`, `doctor`, `init`, `cmd`, and `board/v2`
   functions. After each package, run `go build ./... && go vet ./...` and
   its focused tests; delete test cases that no longer compile because
   their subject is gone.
4. Re-run `deadcode` until clean (or only T-025's migrate entries remain).
5. Run `make build && make test-fast`; record counts and results.

## Boundaries

No behavior change, no refactoring of live code, no renaming of
`internal/board/v2`. Do not touch `internal/migrate` beyond compile fixes,
templates, skills, or documentation (T-023 to T-026).

## Technical Verification

`deadcode` report; focused package tests during iteration;
`make build && make test-fast` at handoff. See
`agent-skills/references/check-method.md`.

## Technical Evidence

Pending execution.

## Drift Notes

None expected. If a reported function turns out to be reachable through
reflection, templates, or `go:linkname`, keep it and record why.
