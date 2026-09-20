---
id: E50-release-validation-cutover/T002-make-live-command-routing-v2-only
status: done
objective: Route every ordinary command through the V2 project model and refuse legacy projects with migration guidance.
depends_on:
    - E50-release-validation-cutover/T001-enforce-one-fail-closed-cutover-preflight
complexity_tier: high
complexity_reason: Removes schema dispatch across command, board, doctor, and resume boundaries.
---

# T002: Make live command routing V2-only

## Problem

Board and doctor still preserve transitional V1 execution paths. After cutover, ordinary commands must share the V2 interpretation and legacy input must stop before rendering, gate execution, or writes.

## Context Files

- `.savepoint/releases/v2/epics/E50-release-validation-cutover/E50-Detail.md`
- `main.go`
- `main_board_test.go`
- `main_board_next_parity_test.go`
- `main_resume_test.go`
- `main_resume_matrix_test.go`
- `cmd/board.go`
- `cmd/board_test.go`
- `cmd/doctor.go`
- `cmd/doctor_test.go`
- `cmd/resume.go`
- `cmd/resume_test.go`
- `internal/board/board.go`
- `internal/board/dispatch_test.go`
- `internal/board/v2/run.go`
- `internal/board/v2/run_test.go`
- `internal/migrate/cutover.go`
- `internal/doctor/checks.go`
- `internal/doctor/checks_test.go`
- `internal/resume/resume.go`

## Acceptance Criteria

- [x] Bare startup, board, doctor, and resume load V2 only and agree on the same indexed state, Release context, and Next projection.
- [x] A V1 project receives a named read-only `migrate --dry-run` action before rendering, quality-gate execution, file watching, or writes.
- [x] An incomplete migration operation receives recovery guidance instead of normal command execution.
- [x] V1-only board filters and routing semantics are removed from public command behavior; valid V2 Objective filtering remains explicit.
- [x] No live command imports or accepts V1 task, epic, defect, or audit-register records.
- [x] TTY, narrow, monochrome, and non-TTY regression coverage still passes for the V2 board.

## Implementation Plan

- [x] Collapse main and command adapters onto the V2 load/run path.
- [x] Replace schema dispatch with typed legacy, pending-operation, and invalid-V2 refusals.
- [x] Remove V1-only filter behavior and update argument validation at the command boundary.
- [x] Make doctor consume only the V2 index and canonical gate results.
- [x] Update command, board, doctor, resume, and parity tests for the V2-only contract.
- [x] Run focused command/consumer tests and the required full gates; record evidence.

## Context Log

- `go test ./cmd ./internal/board ./internal/board/v2 ./internal/doctor .`: PASS.
- `go test ./internal/migrate ./internal/data ./internal/resume ./internal/init`: PASS.
- `make build && make test`: PASS (`go test ./...`; `internal/migrate` 132.206s).
- `git diff --check`: PASS.
- Bare startup, V2 board, V1 refusal, pending-operation recovery, and V2 parity are covered by `main_board_test.go`, `main_board_next_parity_test.go`, `main_resume_test.go`, and `main_resume_matrix_test.go`.
- Public board and doctor adapters now expose only V2 Objective/no-filter surfaces; V1 `--release`/`--epic` parsing and schema dispatch are refused before any live board or doctor work.
- `internal/migrate/cutover.go` exposes the read-only runtime blocker projection used consistently by board, doctor, and resume, while valid V2 Release blockers remain visible to their canonical consumers.
- No `.savepoint/Health-Check.md` is present, so the Quick health check is skipped per `AGENTS.md`.
