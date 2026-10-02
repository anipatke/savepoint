---
id: T-086
title: Measure duplication in maintained code
objective: O-032
status: planned
depends_on: [{task: T-085, requires: clear}]
complexity_tier: spike
complexity_reason: A global duplicate percentage cannot establish the maintained-code baseline or safe refactoring scope.
owner_validation: {required: true}
planned_by: {role: planner, session: planning-o032-20261002-owner-confirmed}
---

# Measure duplication in maintained code

## Outcome

A reproducible maintained-code duplication baseline and concrete clone-repair proposals separate real maintenance costs from preserved history and required template mirrors.

## User Check

Review which files are measured and which real duplicate blocks should be changed.

## Done When

- Record whole-repository baseline evidence and proposed production-Go/test-Go scopes, exact provider targets/exclusions/version and denominator; do not calculate a scoped percentage by filtering the capped evidence list.
- Use project-owned duplication tooling directly with temporary report paths when needed; leave confirmed config, canonical/template copies and archived history unchanged.
- Prove whether this provider version has per-file counts or needs explicitly scoped execution. Identify substantive production clone pairs and reusable test setup separately from required mirrors/fixtures.
- Deliver the named duplication-scope-and-repair decision, including before/after comparison rules, trend restart implications, proposed config change for owner review and bounded follow-up repairs. Keep 3/5 percent thresholds; target <=5 percent after justified repairs, with 3 percent an aim rather than an invented result.

## Context Files

`.savepoint/objectives/O-032-code-health-release-hardening/tasks/T-085-find-why-the-vulnerability-scan-fails.md` (dependency evidence); `.savepoint/health/config.json`; `.savepoint/health/report.md`; `internal/codehealth/config.go`; `internal/codehealth/reader_jscpd.go`; `internal/codehealth/reader_jscpd_test.go`; `internal/codehealth/reader_paths.go`; `internal/codehealth/testdata/readers/duplication/jscpd5.json`; `internal/codehealth/testdata/readers/duplication/versioned.json`.

## Design References

O-032 Confirmed Release-Hardening Design and Owner-Requested Health Analysis; Design sections 1, 6, 7, 9, 12 and 13 as applicable.

## Guardrails

FS-01, FS-04, FS-05, FS-06, ARCH-01, ARCH-02, ARCH-03, CFG-01, CFG-02, TEST-01, TEST-02, TEST-03, TEST-04, TEST-05, TEST-08, TEST-09, STYLE-01, STYLE-02, STYLE-07.

## Implementation Plan

1. Inspect saved top evidence and reader scope/denominator contract.
2. Run bounded maintained-code baseline scenarios with temporary output and document actual numerator/denominator.
3. Rank substantive clone pairs with exact locations and preservation risks.
4. Propose scoped measurement and narrow fixes; return to design when implementation details are not settled.

## Boundaries

Diagnosis and proposals only; no health config edit, automatic exclusion of tests, archive/template modification or threshold changes.

## Technical Verification

`make build && make test-fast` for ordinary handoff; focused `make test-focused TEST=... PKGS=...` only during iteration. Record named per-criterion evidence and limitations. The later independent Check applies `agent-skills/references/check-method.md`; the mandatory Full Objective Check runs its fresh full gate and official Code Health collection. This executor never writes a Check or runs savepoint health setup/report/check. Research completion requires the named decision deliverable; it does not claim any recommended repair has landed.

## Technical Evidence

Pending execution: named per-criterion cases, command/time/toolchain/results, files read/changed, decision deliverable where applicable and limitations.

## Drift Notes

Record implemented responsibility or interface changes for planner reconciliation before the mandatory Full Objective Check. If evidence requires a material unknown repair, return REPLAN REQUIRED rather than inventing another scope.
