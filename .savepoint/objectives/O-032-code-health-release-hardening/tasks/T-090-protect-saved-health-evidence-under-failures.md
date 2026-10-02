---
id: T-090
title: Protect saved health evidence under failures
objective: O-032
status: planned
depends_on: [{task: T-086, requires: clear}]
complexity_tier: high
complexity_reason: Concurrent persistence and explicit retention require filesystem preservation and failure evidence.
owner_validation: {required: false}
planned_by: {role: planner, session: planning-o032-20261002-owner-confirmed}
---

# Protect saved health evidence under failures

## Outcome

Versioned health records and the derived report remain safe to read, write and maintain under competing access and interrupted or failed writes.

## User Check

Review temporary failure evidence showing official snapshots and owner-authored config/ignore bytes unchanged.

## Done When

- Inventory existing coverage and add meaningful missing cases for supported record versions, unsupported newer/older version refusal, malformed/trailing JSON, unsafe paths and bounded fields/files.
- Exercise concurrent same/different immutable snapshot saves and reads during config/report replacement; readers see complete old/new data or clear diagnostics, never silent corruption.
- Failure before/after temp creation and replacement preserves user config/ignore content and prior valid report; cleanup and repeat behavior are explicit. Do not expand unsupported direct inputs into new product promises.
- Explicit PlanPrune/Prune remains dry-plan then maintenance: only older manual snapshots removed, ten newest retained, all official snapshots preserved, ties deterministic, invalid-history and partial-removal failure named. Saving never prunes.
- Keep config/snapshot versions and IDs unchanged; repair only demonstrated failures inside these contracts, with local regression evidence and fresh full gate.

## Context Files

`.savepoint/objectives/O-032-code-health-release-hardening/tasks/T-086-measure-duplication-in-maintained-code.md` (dependency evidence); `internal/codehealth/storage.go`; `internal/codehealth/storage_test.go`; `internal/codehealth/config.go`; `internal/codehealth/model.go`; `internal/codehealth/snapshot.go`; `internal/codehealth/model_test.go`; `internal/codehealth/report.go`; `internal/codehealth/report_test.go`.

## Design References

O-032 Confirmed Release-Hardening Design and Owner-Requested Health Analysis; Design sections 1, 6, 7, 9, 12 and 13 as applicable.

## Guardrails

FS-01, FS-04, FS-05, FS-06, ARCH-01, ARCH-02, ARCH-03, CFG-01, CFG-02, TEST-01, TEST-02, TEST-03, TEST-04, TEST-05, TEST-08, TEST-09, STYLE-01, STYLE-02, STYLE-07, CFG-03, REL-01, REL-02, REL-03.

## Implementation Plan

1. Map validation, persistence and maintenance contracts against existing named tests.
2. Create temporary-project competing-access and obstruction scenarios using semantic byte/ID oracles.
3. Add missing regression evidence and narrow demonstrated repairs.
4. Verify preservation, cleanup and retention boundaries; record full gate.

## Boundaries

No pruning CLI, background cleanup, schema migration, retention-policy change or new global transaction/lock design.

## Technical Verification

Fresh `make test-full` for this platform/migration/distribution-sensitive Task; focused `make test-focused TEST=... PKGS=...` only during iteration. Record named per-criterion evidence and limitations. The later independent Check applies `agent-skills/references/check-method.md`; the mandatory Full Objective Check runs its fresh full gate and official Code Health collection. This executor never writes a Check or runs savepoint health setup/report/check.

## Technical Evidence

Pending execution: named per-criterion cases, command/time/toolchain/results, files read/changed, decision deliverable where applicable and limitations.

## Drift Notes

Record implemented responsibility or interface changes for planner reconciliation before the mandatory Full Objective Check. If evidence requires a material unknown repair, return REPLAN REQUIRED rather than inventing another scope.
