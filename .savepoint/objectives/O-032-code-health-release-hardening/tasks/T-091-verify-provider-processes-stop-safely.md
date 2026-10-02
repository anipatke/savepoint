---
id: T-091
title: Verify provider processes stop safely
objective: O-032
status: planned
depends_on: [{task: T-085, requires: clear}]
complexity_tier: high
complexity_reason: Unix and Windows process-tree cleanup and bounded reports are platform-sensitive.
owner_validation: {required: false}
planned_by: {role: planner, session: planning-o032-20261002-owner-confirmed}
---

# Verify provider processes stop safely

## Outcome

Provider execution has repeatable safety evidence for cancellation, timeout, cleanup, argument/path handling and bounded sanitized failures on supported platforms.

## User Check

Verify a cancelled helper process tree ends and later collection can run normally.

## Done When

- Use controlled local helper processes to prove child/descendant termination, timeout/cancel ordering, pipe release and repeat after failure; retain distinct unavailable/failed/timed_out outcomes.
- Probe stdout report limit and stderr cap at below/exact/above boundaries, malformed/oversized file reports, no-shell argument handling and safe diagnostics without source/secrets leakage.
- Verify independent instances remain independent and cancellation before save writes no snapshot. Report failure after save remains separate from collection failure.
- Use scanner decision as evidence, not permission for an unknown setup/install change; only repair demonstrated execution-contract defects. Provider-specific unknown changes return to design.
- Record Linux full evidence and the existing native Windows CI verification path/results; no Windows skip for a supported real-world case without justified policy. Normal tests run without provider installations or network.

## Context Files

`.savepoint/objectives/O-032-code-health-release-hardening/tasks/T-085-find-why-the-vulnerability-scan-fails.md` (dependency evidence); `internal/codehealth/runner.go`; `internal/codehealth/runner_unix.go`; `internal/codehealth/runner_windows.go`; `internal/codehealth/runner_test.go`; `internal/codehealth/runner_alive_unix_test.go`; `internal/codehealth/runner_alive_windows_test.go`; `internal/codehealth/collect.go`; `internal/codehealth/collect_test.go`; `internal/codehealth/collect_repair_test.go`.

## Design References

O-032 Confirmed Release-Hardening Design and Owner-Requested Health Analysis; Design sections 1, 6, 7, 9, 12 and 13 as applicable.

## Guardrails

FS-01, FS-04, FS-05, FS-06, ARCH-01, ARCH-02, ARCH-03, CFG-01, CFG-02, TEST-01, TEST-02, TEST-03, TEST-04, TEST-05, TEST-08, TEST-09, STYLE-01, STYLE-02, STYLE-07, CFG-03, REL-01, REL-02, REL-03.

## Implementation Plan

1. Reconcile scanner research with execution invariants.
2. Inventory existing runner and collector failure tests.
3. Add boundary and subprocess-tree scenarios; repair only bounded reproduced safety defects.
4. Run relevant local tests and fresh full gate; record native platform limits honestly.

## Boundaries

No provider installation, shell execution, network-dependent normal tests, new capability, thresholds or provider sandbox promise.

## Technical Verification

Fresh `make test-full` for this platform/migration/distribution-sensitive Task; focused `make test-focused TEST=... PKGS=...` only during iteration. Record named per-criterion evidence and limitations. The later independent Check applies `agent-skills/references/check-method.md`; the mandatory Full Objective Check runs its fresh full gate and official Code Health collection. This executor never writes a Check or runs savepoint health setup/report/check.

## Technical Evidence

Pending execution: named per-criterion cases, command/time/toolchain/results, files read/changed, decision deliverable where applicable and limitations.

## Drift Notes

Record implemented responsibility or interface changes for planner reconciliation before the mandatory Full Objective Check. If evidence requires a material unknown repair, return REPLAN REQUIRED rather than inventing another scope.
