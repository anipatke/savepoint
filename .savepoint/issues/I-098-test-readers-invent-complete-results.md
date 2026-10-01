---
id: I-098
title: Test readers invent complete results from unusable reports
type: defect
status: open
source:
  kind: check
  check: C-944
  actor: {role: checker, session: o029-full-check-20261001}
  at: '2026-10-01T09:35:25Z'
tasks: [T-060]
checks: [C-944]
guardrail_ids: [TEST-01, TEST-02]
severity: high
history:
  - at: '2026-10-01T09:35:25Z'
    actor: {role: checker, session: o029-full-check-20261001}
    kind: observed
    check: C-944
    note: Initial Full O-029 Check; see C-944 frozen matrix and embedded independent harness.
---

# I-098: Test readers invent complete results from unusable reports

## Summary

O-029 tests and trustworthy normalization success conditions; T-060 Go stream partial rule and malformed/unusable-report error rule. Frozen M1 null/wrong-schema, terminal completeness, and JUnit declared-count/root cells.

Go accepts JSON null or unrelated JSON as zero failures; a stream with a completed test but no final package action is complete rather than partial. JUnit ignores declared suite counts, and accepts multiple XML roots.

## Evidence

C-944 harness `M1 Go null event`, `M1 Go wrong schema`, `M1 Go missing final package`, `M1 Go incomplete start only`, `M1 JUnit mismatched declared total`, `M1 JUnit two roots`, and M6 Collect null report.

Smallest reproduction: report bytes `null`. Expected reader error/no measurement. Actual Read returns value 0, Partial=false, reason no tests ran. Through Collect this is persisted as outcome available and classified watch rather than failed/unknown. reader_tests.go:144-151 counts syntactically valid lines and skips absent Package; :174 returns success without any recognized event requirement.

Go run(TestA), pass(TestA), EOF with no package terminal produces Partial=false because reader_tests.go:82 increments unfinished only for tests lacking their own result, and :112 checks only that counter. T-060 explicitly requires partial when tests started without final package action.

JUnit `<testsuite tests="2" failures="1"><testcase name="ok"/></testsuite>` produces complete 0 failures/1 total despite declaring an omitted failing test. `<testsuite/><testsuite/>` also returns complete zero despite being multiple document roots. reader_junit.go:33-71 merely looks for any testsuite tag; no suite attributes or single-root check is read.

Existing TestGoTestReaderRejectsUnusableReports checks invalid JSON syntax and empty bytes, not valid JSON with wrong/null schema. The truncated fixture has an unfinished test, not finished tests lacking a package terminal. JUnit malformed tests cover broken syntax only.

## Proof Needed

Require actual supported events/document structure, track package completion separately from individual test completion, and diagnose inconsistent/incomplete JUnit declared counts rather than presenting complete success. Cover these exact M1 cells and M6 Collect persistence/classification. Preserve genuine empty suites. The standalone Go build-fail variant is separately documented as an observation, not a demand to double-count a normal stream containing both build and test events.

Repair under this Issue by default and preserve every completed Task status. This record grants no clearance; a later independent Check verifies remediation.
