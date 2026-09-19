---
id: E50-release-validation-cutover/T006-record-realistic-trials-and-agent-scenarios
status: planned
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

- [ ] A tiny project and an existing-codebase copy each complete preview, apply, recovery or conflict handling, unchanged retry, doctor, board/plain, resume, and archive/reference inspection.
- [ ] Trials use disposable copies, record source hashes and commands, and prove the original projects' bytes and mtimes remain unchanged.
- [ ] One agent executes a fully planned Task without hidden architecture decisions and records reads, extra reads, evidence, and Check handoff.
- [ ] One agent encounters a materially invalid plan, returns `REPLAN REQUIRED`, preserves partial work, and resumes only after a revised plan.
- [ ] A fresh checker detects a seeded material defect, avoids advisory false blockers, records a durable Issue, and verifies the repair through a later Check.
- [ ] Evidence records session/model when supplied, scope, context use, replans, findings, outcomes, commands, and limitations without universal reliability claims.
- [ ] Validation reuses named automated evidence where valid and clearly separates observation from inference.

## Implementation Plan

- [ ] Define the evidence table and reproducible setup for both disposable trial projects.
- [ ] Run the tiny-project and existing-codebase migration/workflow trials with before/after integrity captures.
- [ ] Run the planned-execution, replan, and seeded-defect scenarios in isolated copies and fresh sessions where required.
- [ ] Record commands, outcomes, limitations, and links to supporting automated tests in `E50-Validation.md`.
- [ ] Re-run any automated gate whose evidence changed during the trials.
- [ ] Review the record for overclaims, hidden live mutation, and missing negative-path evidence.

## Context Log

Pending.
