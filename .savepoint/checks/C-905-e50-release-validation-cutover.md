---
id: C-905
scope: {kind: objective, id: O-001}
result: NEEDS WORK
checked_by: {role: checker, session: e50-objective-check-20260921}
executed_session: migration
checked_at: '2026-09-21T08:06:07Z'
reviewed:
  head_commit: f23f914f9b4aef300616b76f1e6addcd177184bf
  files:
    - .savepoint/objectives/O-001-release-validation-cutover/Objective.md
    - .savepoint/objectives/O-001-release-validation-cutover/tasks/T-001-prepare-the-maintainer-controlled-repository-cutover.md
    - .savepoint/objectives/O-001-release-validation-cutover/tasks/T-002-verify-and-reconcile-the-live-v2-cutover.md
    - .savepoint/archive/v1/.savepoint/releases/v2/epics/E50-release-validation-cutover/E50-Validation.md
    - .savepoint/archive/v1/.savepoint/releases/v2/epics/E50-release-validation-cutover/E50-Cutover.md
    - .savepoint/migrations/v1-to-v2.yml
    - internal/migrate/plan.go
    - internal/migrate/convert.go
    - internal/migrate/end_to_end_test.go
    - internal/data/objective_gate_v2.go
    - internal/data/release_gate_v2.go
    - internal/data/release_cutover.go
  dependencies:
    - E49 and E51 completion evidence preserved in the V1 archive
    - Go toolchain and local filesystem only; no network or external service used
issues: [I-016, I-017]
supersedes: null
---

# C-905: Full Objective Check — O-001 / E50

## Verdict

`NEEDS WORK`. Both owned Tasks are owner-marked done and their optional Task Checks were explicitly waived, but two mandatory Objective outcomes are not proven: the recorded rollback snapshot is absent, and the Release-cutover criterion is circular with O-001's own mandatory Check and closure sequence.

## Frozen Scope Lock

1. Acceptance and gates: every T-001 and T-002 acceptance criterion; O-001's cutover, V2-only runtime, distribution, trial, recovery, documentation, and fresh-check quality gates; applicable Guardrails; mandatory Full Objective integration and Design reconciliation.
2. Public paths: migration preview/apply/recovery/no-op behavior; V2 load; manifest/archive/reference accountability; `ResolveObjectiveCompletion`, `ResolveReleaseCompletion`, and `ResolveReleaseCutover`; board/plain/resume/doctor interpretation; build and test gates.
3. Runtime dependencies and effects: confined filesystem inventory, staging/backup/replace/activation ordering, exact-byte archive writes, persisted recovery state, schema activation, and local Go/build tooling. Network, publication, deployment, tagging, and changelog actions are out of scope.
4. Matrix axes: normal migration, empty/not-applicable Release set, malformed or ambiguous source, source/destination edit conflict, interruption before and after publish boundaries, retry/no-op, live/archive representations, current/missing/stale/needs-work evidence, owner acceptance, and post-cutover backup availability. TTY rendering details and unrelated current O-900/UI edits are not applicable to E50's migration outcome.
5. Materiality boundary: admit only reproducible violations of E50/T-001/T-002 criteria, Guardrails, or mandatory Objective integration through supported repository or exported package paths.

## Coverage Matrix

| Row | Normal | Boundary / malformed | Failure / retry | Representation / integration | Result |
| --- | --- | --- | --- | --- | --- |
| Preview and planning | deterministic applicable plan | ambiguity and invalid lifecycle fail closed | write-free retry | source inventory to plan | Proven |
| Apply and recovery | schema activates last | every publish boundary covered | conflict refusal and resume converge | journal, staging, backup, installed bytes | Proven |
| Conversion | active epic/tasks become V2 records | closed work archives | dependency ambiguity blocks | live identity plus byte-preserved source | Proven |
| Accountability | every source has a destination | missing reference rejected | second run is a no-op | manifest/live/archive hashes | Proven |
| Runtime cutover | V2 loads and normal paths are V2-only | V1 input gets migration guidance | incomplete operation gets recovery guidance | board/plain/resume/doctor share typed data | Proven |
| Evidence gates | Task waivers remain non-CLEAR | missing/stale/needs-work fail closed | superseded evidence does not carry | Objective then Release scope | Issue: I-017 |
| Rollback retention | named snapshot claimed | exact directory and checksum file probed | absent snapshot has no recovery proof | runbook claim versus filesystem | Issue: I-016 |
| Distribution and docs | recorded six-platform/package gates | platform structure checks | package-cache workaround rerun exactly | README/Design/Guardrails/AGENTS/templates | Proven |

All mandatory cells are classified. External server/HTTP response cells are not applicable because migration and cutover are local filesystem workflows with no network effect.

## Acceptance Coverage

- T-001: the runbook, owner authority boundary, preview/accounting evidence, stop conditions, recovery procedure, and exclusions are proven. The promised retained rollback material is not available: **Issue I-016**.
- T-002 criterion 1: schema 2, clean V2 load, no pending operation, and shared typed consumers are supported by the recorded live probe and current full tests: **Proven**.
- T-002 criterion 2: the manifest/archive mapping is supported, but “the verified backup remains available” is false: **Issue I-016**.
- T-002 criterion 3: every Release passing cutover cannot be proven during O-001's own Check because R-006 completion depends on this Check and later closure/Release evidence: **Issue I-017**.
- T-002 criteria 4–7: four-state routing, V2-only active guidance, reconciled documentation, focused/full/distribution evidence, and owner-controlled lifecycle handoff are **Proven**.
- O-001 integration and Design reconciliation: converter, archive, recovery, V2 runtime, and canonical resolver boundaries agree with Design. Final cutover sequencing does not: **Issue I-017**.

## Independent Probes and Gates

- `go test ./internal/migrate -run 'TestEndToEnd_(bothFixturesMigrateThroughTheCommandPath|migratedProjectLoadsCleanThroughLoadV2Index|releaseRecordsMappingsAndCutoverGateAgree|temporaryRepositoryCopyMigratesWithReleaseAccountability|everyReferenceResolves|archivedFilesMatchTheFixtureManifestHashes|secondFullRunChangesNothing|previewWritesNothing)' -count=1` — PASS (`107.391s`).
- `make build` — PASS.
- `make test` — PASS; all packages, with `internal/migrate` completing in `132.290s`.
- `git diff --check` — PASS.
- `test -d /tmp/savepoint-cutover-backups/0e7c9ed` — FAIL (absent).
- `test -f /tmp/savepoint-cutover-backups/0e7c9ed/SHA256SUMS` — FAIL (absent).

## Adversarial Findings

The generalized migration path is real: `planEpic` allocates an Objective for every non-closed epic, `planTasksImpl` allocates a Task for every non-done task below it, and `ConvertObjective`/`ConvertTask` render those targets. The live result of one Objective and two Tasks follows from lifecycle selection; it is not a hand-written E50-only converter. Recovery, conflict, preview, archive, reference, and retry probes passed. The two failures above remain after the full matrix completed.

## Materiality

| Issue | Likelihood | Impact | Materiality | Recommendation |
| --- | --- | --- | --- | --- |
| I-016 | High — the recorded path is absent now | High — the promised verified rollback artifact cannot be used | High | Restore or recreate durable verified rollback material before cutover acceptance. |
| I-017 | High — it occurs on the required O-001→R-006 closure path | High — no conforming Objective Check can prove all criteria | High | Replan the criterion/order, then run a fresh Full Objective Check. |

## Non-blocking Observation

The implementation is large relative to the 20 live converted records, but E50 explicitly required preview, exact-byte inventory/archive accounting, interruption recovery at every publish boundary, edit-conflict refusal, no-op retry, and cross-platform replacement. Those controls dominate the code size and passed their focused probes. The transformation itself is generalized. The product choice can reasonably be called expensive for a one-time migration, but this Check found no acceptance violation merely from that proportionality.

## Owner Action

O-001 is not ready for owner closure. Route I-016 to bounded repair and I-017 to design replanning. After both have evidence, run a new immutable Full Objective Check; do not use this `NEEDS WORK` record as Release evidence.
