---
id: I-099
title: Coverage readers accept invalid measurement data
type: defect
status: resolved
source:
  kind: check
  check: C-944
  actor: {role: checker, session: o029-full-check-20261001}
  at: '2026-10-01T09:35:25Z'
tasks: [T-061]
checks: [C-944, C-945]
guardrail_ids: [TEST-01, TEST-02]
severity: medium
history:
  - at: '2026-10-01T09:35:25Z'
    actor: {role: checker, session: o029-full-check-20261001}
    kind: observed
    check: C-944
    note: Initial Full O-029 Check; see C-944 frozen matrix and embedded independent harness.
  - at: '2026-10-01T10:30:00Z'
    actor: {role: executor, session: o029-issue-repair-20261001}
    kind: repair_attempted
    note: 'Go cover reader accepts only set/count/atomic modes and line.col,line.col spans. Vitest reader rejects missing/null or negative statement, function, and branch counters as a failed report. Exact-cell tests plus a Collect test show invalid reports fail with no value and unknown classification while a valid sibling stays intact. make build && make test-fast pass. Awaiting independent recheck.'
  - at: '2026-10-01T10:42:32Z'
    actor: {role: checker, session: o029-independent-recheck-20261001}
    kind: rechecked
    check: C-945
    note: Original reproduction and admitted adjacent frozen cells pass; see C-945 for independent evidence and successful fresh full gate.
resolution:
  disposition: verified
  check: C-945
  actor: {role: checker, session: o029-independent-recheck-20261001}
  at: '2026-10-01T10:42:32Z'
  reason: Independent Full Objective recheck proved the repair within C-944 frozen scope.
---

# I-099: Coverage readers accept invalid measurement data

## Summary

O-029 trustworthy coverage normalization and malformed/unsupported-variant evidence requirements; T-061 provider parsing/error criteria. Frozen M2 Go mode/coordinates and Vitest missing/null/negative-counter cells.

Go validates only the existence of a mode: prefix and a comma after a colon. Vitest counters decode JSON null into integer zero and treat negative execution counts as uncovered statements.

## Evidence

C-944 harness `M2 Go unsupported mode`, `M2 Go bad coordinates`, `M2 Vitest null counter`, `M2 Vitest negative counter`.

`mode: garbage
example.com/m/a/a.go:1.1,2.1 2 1
` is accepted as complete 100%. `mode: set
example.com/m/a/a.go:bad,bad 2 1
` is also accepted as complete 100%. reader_coverage.go:149-154 never checks mode vocabulary; :159-163 never parses source coordinate numbers.

`{"a/a.ts":{"s":{"0":null}}}` and the same report with counter -1 each return complete 0% (one invented/unusable uncovered statement) with no error. reader_coverage.go:183 uses map[string]int64; :218-222 increments total regardless of null/negative validity. These can become bad-code classifications even though the provider supplied no valid measured count.

TestCoverageReadersRejectUnusableReports covers missing mode and gross row syntax, invalid JSON and empty map, not supported-mode vocabulary, malformed coordinates or null/negative counter values.

## Proof Needed

Validate supported Go modes and block positions; reject missing/null or negative Vitest execution counters as provider/report failure rather than measured bad code. Add exact-cell regressions and prove through Collect that invalid reports retain their own failed/no-value/unknown result while valid sibling values remain intact.

Repair under this Issue by default and preserve every completed Task status. This record grants no clearance; a later independent Check verifies remediation.
