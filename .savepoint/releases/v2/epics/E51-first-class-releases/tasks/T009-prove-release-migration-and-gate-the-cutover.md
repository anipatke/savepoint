---
id: E51-first-class-releases/T009-prove-release-migration-and-gate-the-cutover
status: done
objective: Prove first-class Releases end to end on migrated copies and make canonical Release readiness a prerequisite for E50 cutover.
depends_on:
    - E51-first-class-releases/T004-migrate-every-release-and-its-source-promise
    - E51-first-class-releases/T005-diagnose-release-structure-and-readiness
    - E51-first-class-releases/T006-teach-the-workflow-when-a-release-is-worth-using
    - E51-first-class-releases/T008-show-one-release-answer-on-board-and-resume
complexity_tier: high
complexity_reason: Closes the epic with cross-command migration, recovery, documentation, compatibility, and live-cutover evidence.
---

# T009: Prove Release migration and gate the cutover

## Problem

The epic is not complete when each package passes in isolation. A migrated project must preserve every V1 Release and source promise, recover from interruption, load through every V2 consumer, retain the `r` switch, and prevent E50 from cutting over on ambiguous or unaccepted Release state.

## Context Files

- `.savepoint/releases/v2/epics/E51-first-class-releases/E51-Detail.md`
- `.savepoint/releases/v2/epics/E50-release-validation-cutover/E50-Detail.md`
- `.savepoint/releases/v2/v2-PRD.md`
- `.savepoint/releases/v2/v2-Design.md`
- `.savepoint/Design.md`
- `README.md`
- `AGENTS.md`
- `internal/data/e2e_v2_test.go`
- `internal/data/release_gate_v2.go`
- `internal/data/release_cutover.go`
- `internal/data/release_gate_v2_test.go`
- `internal/doctor/checks.go`
- `internal/migrate/apply.go`
- `internal/migrate/apply_test.go`
- `internal/migrate/operation.go`
- `internal/migrate/operation_test.go`
- `internal/migrate/manifest.go`
- `internal/migrate/manifest_test.go`
- `internal/migrate/fixture_test.go`
- `internal/migrate/end_to_end_test.go`
- `main_board_next_parity_test.go`
- `main_resume_matrix_test.go`

## Acceptance Criteria

- [x] An end-to-end fixture with multiple Releases, active/settled work, duplicate scoped legacy IDs, release PRDs, unresolved follow-up, and historical evidence previews and migrates with every source accounted for.
- [x] Dry-run creates no files, directories, probes, backups, manifests, or mtime changes.
- [x] Interruption at each publish boundary is recoverable; retry refuses conflicting user edits and resumes identical work without overwriting them.
- [x] A second unchanged migration changes no bytes, mtimes, IDs, mappings, evidence, router selection, or archive content.
- [x] After migration, data load, doctor, board, the `r` Release selector, plain output, and resume agree on membership, readiness, diagnostics, and next action.
- [x] A temporary copy of this repository migrates and exposes accountable live/archive destinations for every release PRD and release-scoped reference.
- [x] E50 cutover is explicitly refused while the V2 Release is ambiguous, technically unclear, blocked by material Issues, or awaiting owner acceptance.
- [x] E50 consumes the canonical Release completion decision and adds no duplicate cutover/readiness rule.
- [x] Release remains optional in public documentation; publishing, deployment, tagging, and changelog behavior are not claimed.
- [x] V2 PRD/design, project Design, README, E50 plan, and AGENTS Codebase Map describe the implemented result without leaving optional-string metadata as current architecture.
- [x] Focused packages, unchanged V1 regression suites, `go test ./...`, `make build`, and `make test` pass with recorded named evidence.
- [x] The live repository is never migrated by tests or by this task; all migration evidence uses temporary copies.

## Implementation Plan

- [x] Extend the end-to-end migration fixture and interruption matrix with first-class Release records and mappings.
- [x] Drive the migrated fixture through data load, doctor, interactive selector reducers, plain board, and resume.
- [x] Create a temporary copy of this repository and record its preview/apply/retry/no-op/accountability evidence.
- [x] Wire E50's pre-cutover decision to the canonical Release gate and assert every refusal.
- [x] Reconcile release architecture and public workflow documentation after the implementation is proven.
- [x] Update the Codebase Map and E50 dependency/gate wording without claiming the live cutover occurred.
- [x] Run all focused and full quality gates and record results in the Context Log.
- [ ] Hand the completed epic to a fresh independent V1 audit session.

## Context Log

- `go test ./internal/data ./internal/migrate . -run 'TestResolveReleaseCutover|TestPlan_ambiguousReleaseDispositionBlocksCutover|TestEndToEnd_releaseRecordsMappingsAndCutoverGateAgree|TestMigratedReleaseFlowsThroughDoctorBoardSelectorPlainAndResume|TestApply_recoversAtEveryPublishBoundaryWithoutOverwritingUserEdits' -count=1`: PASS (`internal/data` 0.014s, `internal/migrate` 16.838s, root 0.205s).
- `go test ./internal/migrate -run TestEndToEnd_temporaryRepositoryCopyMigratesWithReleaseAccountability -count=1`: PASS (105.036s); the full working-tree copy was temporary and the second apply was byte/mtime/identity/manifest stable.
- `make build`: PASS.
- `make test`: PASS; this runs `go test ./...`, including unchanged V1 regression suites and `internal/migrate` (131.676s).
- `go test ./...`: PASS (direct run after the task update; `internal/migrate` 132.404s).
- `git diff --check`: PASS.
- No `.savepoint/Health-Check.md` is present, so the repository Quick health check was skipped per `AGENTS.md`; independent epic audit remains the next handoff.
- Migration tests never target the live repository; all preview/apply/interruption/no-op evidence uses fixture or temporary-copy roots.
