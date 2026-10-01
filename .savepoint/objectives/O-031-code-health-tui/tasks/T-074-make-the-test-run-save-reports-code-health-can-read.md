---
id: T-074
title: Make the test run save reports Code Health can read
objective: O-031
status: planned
depends_on: []
owner_validation: {required: false}
planned_by: {role: planner, session: planning-o031-20261002}
complexity_tier: low
complexity_reason: Build-tool change that tees the existing go test JSON stream to a file and adds a coverage profile; no Code Health code changes.
---

# Make the test run save reports Code Health can read

## Outcome

Running this repository's tests leaves behind a Go test-results file and a coverage file at the conventional paths, so Code Health can report tests and coverage for Savepoint itself, both on a manual refresh and during the Full Objective Check that reuses the full gate.

## User Check

None beyond the Full Objective Check. Afterwards the owner may run `make test`, then `savepoint health setup`, and expect go test and coverage to be found without a "report missing" note.

## Done When

- `make test` and `make test-fast` write `go-test.json` (the raw `go test -json` event stream) and `coverage.out` (a Go coverage profile) at the repository root, the conventional paths in `internal/codehealth/discovery_catalogue.go`. `make test-focused` writes neither.
- The reports are written even when tests fail, and each is replaced atomically (temp file then rename), so an interrupted run never leaves a half-written report in place of a complete one.
- The terminal timing summary and exit codes are unchanged.
- Both files are listed in `.gitignore`.
- Buildtool tests cover: report written on pass, report written on failure, focused run writes no report, and the summary output unchanged. Tests use temporary directories (TEST-04).
- `make test-full` (which includes `make test`) still passes on Linux and Windows CI; paths use `filepath` joins (FS-05, CFG-03).

## Context Files

`.savepoint/objectives/O-031-code-health-tui/Objective.md`; `Makefile`; `.gitignore`; `internal/buildtool/main.go`, `internal/buildtool/main_test.go`; `internal/codehealth/discovery_catalogue.go` (read only, for report paths).

## Design References

Design sections 12 (Distribution & build) and 13 (Testing); O-030 Confirmed Design Decisions (Artifact reuse).

## Guardrails

FS-05, CFG-02, CFG-03, ARCH-03, TEST-01, TEST-02, TEST-04, TEST-08.

## Implementation Plan

1. Confirm `runGoTestCommand` reads the JSON stream line by line; return REPLAN REQUIRED if the stream is not available to tee.
2. Add an opt-in report mode to the `test` subcommand that tees each raw line to a temp file and adds `-coverprofile` to a temp path, renaming both into place after `go test` exits.
3. Turn it on in the `test` and `test-fast` Makefile targets only.
4. Add `.gitignore` entries and tests.

## Boundaries

No health configuration is written (setup stays human-only), no change to Code Health readers or discovery, no change to `make ci` beyond what `make test` already runs.

## Technical Verification

Gate definitions change, so fresh `make test-full` before handoff. Later evaluated per `agent-skills/references/check-method.md`.

## Technical Evidence

Pending execution.

## Drift Notes

Design section 13 should mention the two report files at the Objective Check.
