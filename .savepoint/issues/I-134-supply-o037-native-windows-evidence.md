---
id: I-134
title: Supply native Windows evidence for O-037
type: verification
status: resolved
source:
  kind: check
  check: C-963
  actor: {role: checker, session: check-o037-20261004}
  at: '2026-10-04T02:31:20Z'
tasks: [T-108, T-109, T-110]
checks: [C-963, C-964]
guardrail_ids: [CFG-03]
resolution:
  disposition: verified
  check: C-964
  actor: {role: checker, session: recheck-o037-20261004}
  at: '2026-10-04T02:43:55Z'
  reason: 'GitHub verified CI run 37171504754 native windows-tests success on exact reviewed 4520c93a7421dd9cc80bfe6e83cb14445aaac506.'
history:
  - at: '2026-10-04T02:31:20Z'
    actor: {role: checker, session: check-o037-20261004}
    kind: observed
    check: C-963
    note: Initial independent Full Objective Check.
  - at: '2026-10-04T02:43:55Z'
    actor: {role: checker, session: recheck-o037-20261004}
    kind: rechecked
    check: C-964
    note: 'CLEAR. GitHub verified CI run 37171504754 native windows-tests success on exact reviewed 4520c93a7421dd9cc80bfe6e83cb14445aaac506.'
---

# I-134: Supply native Windows evidence for O-037

## Summary

Supply native Windows evidence for O-037 before O-037 clearance.

## Evidence

O-037 Readiness and Verification and all three Tasks require current native windows-tests repository CI evidence before independent Full clearance.

The implementation is an uncommitted working-tree change on cdfa699. Task evidence explicitly states Windows CI is pending. Fresh make test-full passes host tests and cross-compiles Windows, which does not execute native tests. `.github/workflows/ci.yml:32–50` defines the required native job. The owner replied “completed” to the bundled evidence/walkthrough question, confirming completion but supplied no CI run or revision proof. Native success for the reviewed source therefore remains unverified.

No Windows failure is alleged.

## Proof Needed

Provide successful native windows-tests CI evidence tied to the reviewed O-037 implementation (run URL/revision and job result). If source changes, refresh host/full evidence too. Do not substitute cross-compilation. C-963 M6 remains the verification perimeter.
