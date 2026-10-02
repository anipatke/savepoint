---
id: T-087
title: Simplify diagnostic repair rules
objective: O-032
status: planned
complexity_tier: high
complexity_reason: Ordered predicates and many diagnostic mappings require exact behavior preservation.
owner_validation: {required: false}
planned_by: {role: planner, session: planning-o032-20261002-owner-confirmed}
---

# Simplify diagnostic repair rules

## Outcome

Doctor retains every existing diagnostic and repair instruction while its largest branching functions become readable rule and lookup data.

## User Check

Compare representative doctor diagnostics and repair text before/after, including legacy and malformed cases.

## Done When

- V2ProblemRepair exact-name mappings and fallback retain their current text; v2DiagnosticName retains sentinel matching and special Goal/reference distinctions.
- SuggestRepair retains errors.Is precedence and ordered substring predicate behavior, including overlapping matches, mixed error chains and unknown input. Do not remove legacy compatibility guidance as incidental cleanup.
- Refactor exact mappings into typed lookup data and ordered cases into explicit named rules; no duplicate policy or new diagnostic vocabulary.
- Record exhaustive mapping/default and ordered-overlap regression evidence plus unchanged-scope Lizard before/after measurements. Target each touched production function <=20 CCN, aiming for <=10 where justified; preserve tests and flag any justified residual for owner disposition.

## Context Files

`internal/doctor/repairs.go`; `internal/doctor/repairs_test.go`; `internal/doctor/checks.go`; `internal/doctor/checks_test.go`; `internal/doctor/v2_runtime_test.go`.

## Design References

O-032 Confirmed Release-Hardening Design and Owner-Requested Health Analysis; Design sections 1, 6, 7, 9, 12 and 13 as applicable.

## Guardrails

FS-01, FS-04, FS-05, FS-06, ARCH-01, ARCH-02, ARCH-03, CFG-01, CFG-02, TEST-01, TEST-02, TEST-03, TEST-04, TEST-05, TEST-08, TEST-09, STYLE-01, STYLE-02, STYLE-07.

## Implementation Plan

1. Capture current mapping, precedence and default behavior as independent expected cases.
2. Separate exact copy lookup from ordered matching without changing first-match semantics.
3. Keep special conditional mappings explicit and reuse existing sentinel ownership.
4. Run focused doctor regressions and compare current-scope complexity before handoff.

## Boundaries

Only doctor diagnostic/repair organization; no new repair actions, planning record writes, health scopes or thresholds.

## Technical Verification

For before/after complexity evidence, the owner-confirmed scope authorizes a direct project-owned Lizard invocation over the scoped source paths with temporary output; run no Savepoint health command. If the executable is unavailable, record the measurement gap for the owner instead of installing it.

`make build && make test-fast` for ordinary handoff; focused `make test-focused TEST=... PKGS=...` only during iteration. Record named per-criterion evidence and limitations. The later independent Check applies `agent-skills/references/check-method.md`; the mandatory Full Objective Check runs its fresh full gate and official Code Health collection. This executor never writes a Check or runs savepoint health setup/report/check.

## Technical Evidence

Pending execution: named per-criterion cases, command/time/toolchain/results, files read/changed, decision deliverable where applicable and limitations.

## Drift Notes

Record implemented responsibility or interface changes for planner reconciliation before the mandatory Full Objective Check. If evidence requires a material unknown repair, return REPLAN REQUIRED rather than inventing another scope.
