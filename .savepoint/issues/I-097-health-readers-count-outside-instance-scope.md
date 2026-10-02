---
id: I-097
title: Health readers count outside the instance scope
type: defect
status: resolved
source:
  kind: check
  check: C-944
  actor: {role: checker, session: o029-full-check-20261001}
  at: '2026-10-01T09:35:25Z'
tasks: [T-060, T-061, T-063]
checks: [C-944, C-945]
guardrail_ids: [TEST-01, TEST-02]
severity: high
history:
  - at: '2026-10-01T09:35:25Z'
    actor: {role: checker, session: o029-full-check-20261001}
    kind: observed
    check: C-944
    note: Initial Full O-029 Check; see C-944 frozen matrix and embedded independent harness.
  - at: '2026-10-01T10:30:00Z'
    actor: {role: executor, session: o029-issue-repair-20261001}
    kind: repair_attempted
    note: 'Go and JUnit readers now count only tests inside the instance scope (owned-but-excluded packages and mapped-but-excluded cases are dropped; unplaceable ones still count). coverage.py rebuilds its headline and branch details from in-scope per-file summaries when a scope is set. jscpd rebuilds totals from per-file format sources, withholds the value as partial when a scoped report has none, and keeps clone evidence on an in-scope file. Regression tests in reader_repair_test.go; make build && make test-fast pass. Awaiting independent recheck.'
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

# I-097: Health readers count outside the instance scope

## Summary

O-029 success condition 6 requires each reader to preserve instance scope. T-060, T-061 and T-063 outcomes promise the configured scope. Frozen M1/M2/M4 scope and exclusions cells.

Go and JUnit filter only affected-item evidence, not headline counts. coverage.py uses whole-report totals while filtering only file evidence. jscpd retains whole-report statistics even when every clone is outside scope, and a mixed-scope clone points at an excluded first file.

## Evidence

Independent harness cells `M1 Go scoped count`, `M1 Vitest scoped count`, `M1 Pytest excluded failure`, `M2 coverage.py scoped count`, `M4 jscpd all clones outside scope`, and `M4 jscpd mixed scope affected path` in C-944.

Smallest Go/JUnit reproduction: scope a/**, passing test in a, failing test in b. Expected 0 scoped failures; actual 1, without affected items. Go: reader_tests.go:56-101; JUnit: reader_junit.go:57-66, scope applied only by appendJUnitEvidence at :106. Python coverage: a/a.py is 1/1, b/b.py 0/1, totals 1/2, scope a/**; expected scoped 100%, actual 50% available at reader_coverage.go:276-311. jscpd: scope a/**, totals 20/100, only clone b/b.ts/b/c.ts; expected explicit unusable/partial scope mismatch if totals cannot be reconstructed, actual available 20% at reader_jscpd.go:93-116. With firstFile b/b.ts and secondFile a/c.ts, actual evidence Path=b/b.ts outside scope.

TestJUnitReaderHonoursScope explicitly expects the leaked failure count, and TestJscpdReaderValue's scope case expects unchanged whole-project totals. Those tests prove the implementation, not the promised scope invariant.

## Proof Needed

Filter tests before counting, distinguish unresolvable paths from known excluded paths, and compute scoped coverage from relevant per-file counts. For jscpd, use provider data that can establish the scoped denominator or explicitly reject/mark partial a report whose scope cannot be established; do not label whole-project totals as a complete scoped result. Keep clone evidence within the declared scope or explain a deliberate cross-scope evidence representation. Prove includes/excludes and mixed-project counts across the four affected reader paths, plus unaffected Go/Vitest coverage, Lizard and OSV scope cells.

Repair under this Issue by default and preserve every completed Task status. This record grants no clearance; a later independent Check verifies remediation.
