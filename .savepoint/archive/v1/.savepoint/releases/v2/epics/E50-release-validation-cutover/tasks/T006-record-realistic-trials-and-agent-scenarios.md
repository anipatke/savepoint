---
id: E50-release-validation-cutover/T006-record-realistic-trials-and-agent-scenarios
status: done
objective: Record two isolated project trials and three bounded agent scenarios against the packaged V2 workflow.
depends_on:
    - E50-release-validation-cutover/T005-prove-six-platform-distribution-and-cli-contract
complexity_tier: high
complexity_reason: Requires reproducible migration, workflow, and independent-agent evidence across realistic cases.
---

# T006: Record realistic trials and agent scenarios

## Problem

Unit and integration tests cannot prove that the final workflow is understandable and recoverable in realistic use. The release needs bounded human-readable evidence without treating this repository or any live user project as a mutable fixture.

## Context Files

- `.savepoint/releases/v2/epics/E50-release-validation-cutover/E50-Detail.md`
- `.savepoint/releases/v2/epics/E50-release-validation-cutover/E50-Validation.md`
- `.savepoint/releases/v2/v2-Design.md`
- `.savepoint/releases/v2/v2-PRD.md`
- `internal/migrate/end_to_end_test.go`
- `internal/data/e2e_v2_test.go`
- `internal/board/v2/run_test.go`
- `main_resume_matrix_test.go`
- `README.md`

## Acceptance Criteria

- [x] A tiny project and an existing-codebase copy each complete preview, apply, recovery or conflict handling, unchanged retry, doctor, board/plain, resume, and archive/reference inspection.
- [x] Trials use disposable copies, record source hashes and commands, and prove the original projects' bytes and mtimes remain unchanged.
- [x] One agent executes a fully planned Task without hidden architecture decisions and records reads, extra reads, evidence, and Check handoff.
- [x] One agent encounters a materially invalid plan, returns `REPLAN REQUIRED`, preserves partial work, and resumes only after a revised plan.
- [x] A fresh checker detects a seeded material defect, avoids advisory false blockers, records a durable Issue, and verifies the repair through a later Check.
- [x] Evidence records session/model when supplied, scope, context use, replans, findings, outcomes, commands, and limitations without universal reliability claims.
- [x] Validation reuses named automated evidence where valid and clearly separates observation from inference.

## Implementation Plan

- [x] Define the evidence table and reproducible setup for both disposable trial projects.
- [x] Run the tiny-project and existing-codebase migration/workflow trials with before/after integrity captures.
- [x] Run the planned-execution, replan, and seeded-defect scenarios in isolated copies and fresh sessions where required.
- [x] Record commands, outcomes, limitations, and links to supporting automated tests in `E50-Validation.md`.
- [x] Re-run any automated gate whose evidence changed during the trials.
- [x] Review the record for overclaims, hidden live mutation, and missing negative-path evidence.

## Context Log

- 2026-09-20: Read the router, E50 detail, V2 design/PRD, and the task's bounded context files. No `.savepoint/Health-Check.md` is present, so the Quick health check is skipped.
- 2026-09-20: Set this task to `status: in_progress` with `stage: build` before editing.
- 2026-09-20: Added `E50-Validation.md` with two disposable trials, source tree/manifest hashes, mtime-integrity assertions, command transcripts, three fixture-backed agent scenarios, observation/inference boundaries, and limitations.
- 2026-09-20: Focused migration evidence passed, including preview, apply, every publish-boundary recovery, source/destination conflicts, unchanged retry, archive/reference mapping, and repository-copy Release accountability (`go test ./internal/migrate ...`, `ok`, 109.913s; recovery boundary probe `ok`, 15.078s).
- 2026-09-20: Focused data, doctor, board/plain, resume, replan, Issue, and advisory checks passed; the exact commands and outcomes are recorded in `E50-Validation.md`.
- 2026-09-20: Required full gate passed: `make build && make test` (build succeeded; all Go packages passed, including `internal/migrate` in 125.530s).
- 2026-09-20: Record review found no live-project mutation or universal reliability claim. The scenario session/model fields are explicitly `not supplied`; fixture actor sessions are provenance only.
