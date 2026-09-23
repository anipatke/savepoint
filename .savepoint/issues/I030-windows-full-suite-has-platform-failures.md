---
id: I030
title: Windows full suite has platform failures
type: verification
status: open
source:
  kind: check
  check: C908
  actor: {role: checker, session: o016-full-check-20260923}
  at: '2026-09-23T01:42:44Z'
tasks: [T016]
checks: [C908]
guardrail_ids: [CFG-02, TEST-08]
severity: high
history:
  - at: '2026-09-23T01:42:44Z'
    actor: {role: checker, session: o016-full-check-20260923}
    kind: observed
    check: C908
    note: The Windows-native command configured by T016 also reproduces board watcher filtering and init guide-casing failures, so the new CI job cannot pass even after the migration setup conflict is isolated.
  - at: '2026-09-23T03:58:56Z'
    actor: {role: owner, session: user}
    kind: deferred
    note: >-
      Owner narrowed O016 to focused Windows migration setup coverage and cross-builds; the full Windows runtime suite is not required. T018 and T019 were removed. Keep I030 open because these failures remain uncorrected; revisit only if the full suite returns to scope.
---

# I030: Windows full suite has platform failures

## Summary

T016 adds a complete Windows test job, but the current Windows suite has two
additional reproducible platform failures. O016 widened the supported gate to
this path and explicitly promises Windows execution, so a configured job that
is known to fail cannot satisfy the Objective.

## Current Scope

The owner later narrowed O016 to focused Windows migration setup coverage and
cross-builds. T018 and T019 were removed from O016, and the full Windows
runtime suite is deferred. These observed failures remain unfixed; this Issue
stays open as deferred follow-up, without a claim of repair or acceptance.

## Evidence

- Windows Go 1.26.2/amd64, from a local NTFS copy of the current working tree:
  `go test -count=1 ./internal/board/v2` fails
  `TestV2WatcherIgnoresExcludedTrees` because an excluded tree produces a
  `v2FileChangeMsg` (`internal/board/v2/watch_test.go:63`).
- `go test -count=1 ./internal/init -run
  TestUpgradeProjectAssets_conflictSidecarKeepsGuideCasing -v` fails because
  upgrade creates a second agent guide under canonical casing
  (`internal/init/upgrade_test.go:436`).
- The complete Windows command configured at `.github/workflows/ci.yml:29`
  propagates these failures and exits nonzero.
- These are supported Windows paths covered by T016's CI acceptance criterion,
  CFG-02, and TEST-08; they are not remote hypothetical configurations.

## Proof Needed

- No repair is required by the current O016 scope. If the full Windows runtime
  suite is reinstated, repair the watcher and guide-casing failures and rerun
  the focused reproductions plus the full Windows suite.
