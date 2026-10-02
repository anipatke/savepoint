---
id: T-095
title: Verify the v2.1 binary and package boundary
objective: O-032
status: planned
depends_on: [{task: T-087, requires: clear}, {task: T-088, requires: clear}, {task: T-089, requires: clear}, {task: T-090, requires: clear}, {task: T-091, requires: clear}, {task: T-092, requires: clear}, {task: T-093, requires: clear}]
complexity_tier: medium
complexity_reason: Release evidence must cover six archives, checksums and provider-free runtime distribution.
owner_validation: {required: false}
planned_by: {role: planner, session: planning-o032-20261002-owner-confirmed}
---

# Verify the v2.1 binary and package boundary

## Outcome

Current release artifacts prove Savepoint remains one executable per platform with correct version/checksums and no bundled or installed analysis providers.

## User Check

Read the distribution inventory and verify each archive contains one platform executable and has a checksum.

## Done When

- Run make ci and retain command/time/toolchain/result evidence, six target inventory, archive member validation, checksum verification and native version smoke outcome.
- Demonstrate packaged help/core behavior works without provider executables on PATH; unavailable providers remain honest unavailable states, and no runtime download/install is attempted.
- Inspect package/runtime dependency boundary and document provider prerequisites separately from the shipped binary.
- Record current native Windows full-test CI result when accessible; if unavailable, record the exact missing evidence for the Full Objective Check, never claim cross-compilation is native runtime validation.
- Add narrowly missing packaging failure regressions only if evidence exposes an uncovered release contract; no publishing, tagging or deployment.

## Context Files

`.savepoint/objectives/O-032-code-health-release-hardening/tasks/T-087-simplify-diagnostic-repair-rules.md` (dependency evidence); `.savepoint/objectives/O-032-code-health-release-hardening/tasks/T-088-separate-test-stream-processing-responsibilities.md` (dependency evidence); `.savepoint/objectives/O-032-code-health-release-hardening/tasks/T-089-simplify-board-reload-state-restoration.md` (dependency evidence); `.savepoint/objectives/O-032-code-health-release-hardening/tasks/T-090-protect-saved-health-evidence-under-failures.md` (dependency evidence); `.savepoint/objectives/O-032-code-health-release-hardening/tasks/T-091-verify-provider-processes-stop-safely.md` (dependency evidence); `.savepoint/objectives/O-032-code-health-release-hardening/tasks/T-092-preserve-projects-when-adopting-code-health.md` (dependency evidence); `.savepoint/objectives/O-032-code-health-release-hardening/tasks/T-093-exercise-all-health-readers-together.md` (dependency evidence); `internal/buildtool/main.go`; `internal/buildtool/main_test.go`; `Makefile`; `package.json`; `.github/workflows/ci.yml`; `main.go`; `main_health_test.go`.

## Design References

O-032 Confirmed Release-Hardening Design and Owner-Requested Health Analysis; Design sections 1, 6, 7, 9, 12 and 13 as applicable.

## Guardrails

FS-01, FS-04, FS-05, FS-06, ARCH-01, ARCH-02, ARCH-03, CFG-01, CFG-02, TEST-01, TEST-02, TEST-03, TEST-04, TEST-05, TEST-08, TEST-09, STYLE-01, STYLE-02, STYLE-07, CFG-03, REL-01, REL-02, REL-03.

## Implementation Plan

1. Inspect current archive/checksum/package tests and CI matrix.
2. Build and validate all distribution outputs using existing buildtool gates.
3. Smoke provider-free packaged runtime in a temporary PATH/project environment.
4. Record artifact and platform evidence for independent Objective review.

## Boundaries

No publication, dependency installation, new CI/CD integration or separate release runner; fix only demonstrated distribution contract failures.

## Technical Verification

Fresh `make test-full` for this platform/migration/distribution-sensitive Task; `make ci` additionally for distribution evidence; focused `make test-focused TEST=... PKGS=...` only during iteration. Record named per-criterion evidence and limitations. The later independent Check applies `agent-skills/references/check-method.md`; the mandatory Full Objective Check runs its fresh full gate and official Code Health collection. This executor never writes a Check or runs savepoint health setup/report/check.

## Technical Evidence

Pending execution: named per-criterion cases, command/time/toolchain/results, files read/changed, decision deliverable where applicable and limitations.

## Drift Notes

Record implemented responsibility or interface changes for planner reconciliation before the mandatory Full Objective Check. If evidence requires a material unknown repair, return REPLAN REQUIRED rather than inventing another scope.
