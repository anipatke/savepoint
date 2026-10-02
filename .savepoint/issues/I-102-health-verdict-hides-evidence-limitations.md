---
id: I-102
title: Health verdict hides evidence limitations behind blocking findings
type: defect
status: resolved
source:
  kind: check
  check: C-946
  actor: {role: checker, session: check-o030-20261001-independent}
  at: '2026-10-01T20:06:17Z'
tasks: [T-066, T-068]
checks: [C-946, C-947]
guardrail_ids: [TEST-01, TEST-02]
resolution:
  disposition: verified
  check: C-947
  actor: {role: checker, session: recheck-o030-20261001-independent}
  at: '2026-10-01T20:15:49Z'
  reason: Original reproduction and frozen-scope recheck passed; C-947 records CLEAR.
history:
  - at: '2026-10-01T20:06:17Z'
    actor: {role: checker, session: check-o030-20261001-independent}
    kind: observed
    check: C-946
    note: Reproduced during the initial Full Objective Check of O-030.
  - at: '2026-10-01T20:15:49Z'
    actor: {role: checker, session: recheck-o030-20261001-independent}
    kind: rechecked
    check: C-947
    note: Repair independently verified against C-946's frozen scope and a fresh full gate.
---

# I-102: Health verdict hides evidence limitations behind blocking findings

## Summary

O-030 success condition 5 requires incomplete coverage, stale artifacts and unhealthy measurements to remain distinct. T-066 Done When 2 explicitly labels partial evidence incomplete, and its default blockers must still block. The current single Kind/Reason drops other applicable statements.

## Evidence

`internal/codehealth/gate.go:161-177`: default blockers and opt-in Needs Attention return before stale/partial processing; stale also wins over partial. At :178-179, unknown-severity review wording is restricted to KindNoFinding, so incomplete/stale results lose that review message as well.

Independent policy matrix: 156 disposition cells passed. A required tests result with `OutcomePartial`, `FreshnessFresh`, `Value.Number=1`, and `Reason="package example/missing has no terminal event"` rendered only:

```
Code Health blocks clearance.
- tests via go-test-json [required]: blocks clearance; unhealthy measurement: 1 failing.
```

Changing freshness to stale produces the same line. Expected: retain the blocking failure AND state incomplete/missing package, plus stale when applicable. Actual: limitation and reason disappear.

Supported collection reproduction: temporary Git V2 project, module `example.com/m`, test report from `internal/codehealth/testdata/readers/tests/go-fail.jsonl`, followed by `{"Action":"start","Package":"example.com/m/missing"}` without a terminal event. Configure required go-test-json to read this report; set report mtime newer than inputs. Run `healthcheck.Run` with the test helpers' fake Lizard runner. Store round trip shows tests outcome `partial/fresh`, while the command's tests line is exactly the unhealthy-only line above. A separate fresh Go coverage report is reused; only Lizard executes. Snapshot still blocks correctly; no false clearance was reproduced.

Commands: `go test ./internal/codehealth -run TestO030IndependentMatrix -v -count=1` and `go test ./internal/healthcheck -run TestO030WorkflowMatrix -v -count=1`. The independent limitation assertions fail. Scratch harnesses are embedded in C-946 for repeatability; no production or test source was retained.

Existing `TestEvaluateRules` cases for partial failures assert Disposition and a single Kind only, so their passing result does not prove preservation of all evidence statements.

## Proof Needed

Preserve independent evidence-quality and measurement statements together in the verdict and command output. Verify fresh/stale crossed with available/partial for passing/failing tests, known severe and unknown-severity vulnerabilities, and opt-in coverage/complexity/duplication. Continue blocking default/opt-in findings regardless of partial evidence. Prove reasons and review wording survive combinations; do not relabel failed tools as bad code.
