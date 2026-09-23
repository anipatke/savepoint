---
id: I029
title: Windows migration tests have two TestMain functions
type: defect
status: resolved
source:
  kind: check
  check: C908
  actor: {role: checker, session: o016-full-check-20260923}
  at: '2026-09-23T01:42:44Z'
tasks: [T014, T016, T017]
checks: [C908, C909]
guardrail_ids: [CFG-02, TEST-08]
severity: high
resolution:
  disposition: verified
  check: C909
  actor: {role: checker, session: o016-full-recheck-20260923}
  at: '2026-09-23T04:20:00Z'
  reason: One effective migration TestMain on Windows with helper dispatch in init; focused Windows check passes, a duplicate-TestMain mutation fails it, and a fresh Linux make test-full passed.
history:
  - at: '2026-09-23T01:42:44Z'
    actor: {role: checker, session: o016-full-check-20260923}
    kind: observed
    check: C908
    note: The Windows-native O016 full-suite probe fails during internal/migrate setup because T014 added an unconditional TestMain while replace_windows_test.go already defines the package TestMain on Windows.
  - at: '2026-09-23T03:58:56Z'
    actor: {role: executor, session: o016-t017-focused-20260923}
    kind: repair_attempted
    note: >-
      T017 removed the Windows TestMain collision by dispatching helper modes in init. The focused Windows command passed in 1.912s and make build-all passed. make test-full did not pass in the available WSL environment, so the Issue remains open for an independent Check.
  - at: '2026-09-23T04:20:00Z'
    actor: {role: checker, session: o016-full-recheck-20260923}
    kind: rechecked
    note: C909 reproduced the focused Windows migration check (Go 1.26.2 windows/amd64, all three tests pass), proved it fails when a second TestMain is reintroduced, and recorded a fresh passing Linux make test-full with native npm; I029 is resolved.
    check: C909
---

# I029: Windows migration tests have two TestMain functions

## Summary

O016's shared migration-fixture setup added `TestMain` in
`internal/migrate/end_to_end_test.go`, but the package already has the
Windows-only `TestMain` in `internal/migrate/replace_windows_test.go`. The
package therefore cannot compile on Windows, and T016's new Windows CI job
cannot run the complete migration suite.

## Current Scope

O016 has since replaced the full Windows CI suite with a focused migration
setup check. This Issue remains about the TestMain conflict; the complete
Windows package suite is no longer a required O016 gate.

## Evidence

- `internal/migrate/end_to_end_test.go:64` defines the new unconditional
  `TestMain` that prepares shared converted fixtures.
- `internal/migrate/replace_windows_test.go:69` defines the existing Windows
  package `TestMain`.
- Windows Go 1.26.2/amd64, from a local NTFS copy of the current working tree:
  `go test -count=1 ./internal/migrate` fails at setup with
  `multiple definitions of TestMain`.
- `.github/workflows/ci.yml:29` makes this package part of the new complete
  Windows suite, so the failure violates T016's Windows CI criterion, CFG-02,
  and TEST-08.

## Proof Needed

- Keep one package `TestMain`; Windows helper modes must dispatch before shared
  fixture preparation.
- Review the focused Windows check for shared fixture results, private copies,
  held-open helpers, and interruption helpers. The full Windows package suite
  is not required under the narrowed O016 scope.
- Require a passing fresh Linux `make test-full` and an independent Check of
  this repair before closing I029.
