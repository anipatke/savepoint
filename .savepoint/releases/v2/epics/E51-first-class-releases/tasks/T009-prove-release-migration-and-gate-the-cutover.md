---
id: E51-first-class-releases/T009-prove-release-migration-and-gate-the-cutover
status: planned
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

- [ ] An end-to-end fixture with multiple Releases, active/settled work, duplicate scoped legacy IDs, release PRDs, unresolved follow-up, and historical evidence previews and migrates with every source accounted for.
- [ ] Dry-run creates no files, directories, probes, backups, manifests, or mtime changes.
- [ ] Interruption at each publish boundary is recoverable; retry refuses conflicting user edits and resumes identical work without overwriting them.
- [ ] A second unchanged migration changes no bytes, mtimes, IDs, mappings, evidence, router selection, or archive content.
- [ ] After migration, data load, doctor, board, the `r` Release selector, plain output, and resume agree on membership, readiness, diagnostics, and next action.
- [ ] A temporary copy of this repository migrates and exposes accountable live/archive destinations for every release PRD and release-scoped reference.
- [ ] E50 cutover is explicitly refused while the V2 Release is ambiguous, technically unclear, blocked by material Issues, or awaiting owner acceptance.
- [ ] E50 consumes the canonical Release completion decision and adds no duplicate cutover/readiness rule.
- [ ] Release remains optional in public documentation; publishing, deployment, tagging, and changelog behavior are not claimed.
- [ ] V2 PRD/design, project Design, README, E50 plan, and AGENTS Codebase Map describe the implemented result without leaving optional-string metadata as current architecture.
- [ ] Focused packages, unchanged V1 regression suites, `go test ./...`, `make build`, and `make test` pass with recorded named evidence.
- [ ] The live repository is never migrated by tests or by this task; all migration evidence uses temporary copies.

## Implementation Plan

- [ ] Extend the end-to-end migration fixture and interruption matrix with first-class Release records and mappings.
- [ ] Drive the migrated fixture through data load, doctor, interactive selector reducers, plain board, and resume.
- [ ] Create a temporary copy of this repository and record its preview/apply/retry/no-op/accountability evidence.
- [ ] Wire E50's pre-cutover decision to the canonical Release gate and assert every refusal.
- [ ] Reconcile release architecture and public workflow documentation after the implementation is proven.
- [ ] Update the Codebase Map and E50 dependency/gate wording without claiming the live cutover occurred.
- [ ] Run all focused and full quality gates and record results in the Context Log.
- [ ] Hand the completed epic to a fresh independent V1 audit session.

## Context Log

Pending.
