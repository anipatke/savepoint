---
id: T-087
title: Simplify diagnostic repair rules
objective: O-032
status: done
complexity_tier: high
complexity_reason: Ordered predicates and many diagnostic mappings require exact behavior preservation.
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: planning-o032-20261002-owner-confirmed}
check_waiver:
    task: T-087
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-02T03:55:28Z"
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

Executed 2026-10-02 by the executor on go1.26.2 linux/amd64.

**Per-criterion outcomes**

1. V2ProblemRepair / v2DiagnosticName: before refactor a throwaway golden dump (1003 lines, not committed) recorded every output; after refactor it was byte-identical (`cmp`). It covered every exact name plus empty/unknown input, every sentinel wrapped plain and with "release id", "Goal id" and "scope names missing release" wording, and the generic fallback. New permanent tests: `TestV2ProblemRepair_everyNamedMappingIsDistinctFromDefault`, `TestV2DiagnosticName_specialCases`.
2. SuggestRepair: the same golden dump covered all message terms individually, pairwise overlaps, the four typed sentinels wrapped with message wording, and `errors.Join` chains; identical. Legacy V1 guidance is unchanged. New permanent tests: `TestSuggestRepair_typedSentinelBeatsMessageWording`, `TestSuggestRepair_firstMatchingMessageRuleWins`.
3. Exact mappings are now the typed `v2ProblemRepairs` map and the ordered `v2DiagnosticRules` / `typedRepairRules` / `messageRepairRules` slices. No new vocabulary or duplicated policy.
4. Lizard (direct invocation over `internal/doctor/repairs.go` and `checks.go`, output in scratchpad; no Savepoint health command run): SuggestRepair CCN 46 -> 5, V2ProblemRepair 46 -> 2, v2DiagnosticName 44 -> 4. New helpers are all <=3 CCN. `releaseDiagnosticsForIndex` (CCN 14, untouched) is out of scope and unchanged. No residual above 20.

**Commands:** `make build && make test-fast` passed at 2026-10-02T03:54Z; `go test ./internal/doctor -count=1` passed.

**Files read:** the five Context Files that exist were read in part (repairs.go, repairs_test.go, checks.go, checks_test.go); `v2_runtime_test.go` was not needed. Extra read: `internal/doctor/v2_runtime.go` was only grepped to confirm the v2DiagnosticName call site, and `agent-skills/savepoint-task/SKILL.md` and the O-032 Objective were read for workflow.
**Files changed:** `internal/doctor/repairs.go`, `internal/doctor/checks.go`, `internal/doctor/repairs_test.go`, router selection (T-086 -> T-087), this Task.

**Limitations:** The golden dump was a temporary equivalence proof, not kept as a test; permanent tests are narrower. V2ConsistencyRepair and GateSuggestion were left as is (already low CCN). No Task Check requested and no waiver recorded.

## Drift Notes

Record implemented responsibility or interface changes for planner reconciliation before the mandatory Full Objective Check. If evidence requires a material unknown repair, return REPLAN REQUIRED rather than inventing another scope.
