---
id: I-030
title: Windows full suite has platform failures
type: verification
status: open
source:
    kind: check
    check: C-908
    actor: {role: checker, session: o016-full-check-20260923}
    at: '2026-09-23T01:42:44Z'
tasks: [T-016]
checks: [C-908]
guardrail_ids: [CFG-02, TEST-08]
severity: high
history:
    - at: '2026-09-23T01:42:44Z'
      actor: {role: checker, session: o016-full-check-20260923}
      kind: observed
      check: C-908
      note: The Windows-native command configured by T-016 also reproduces board watcher filtering and init guide-casing failures, so the new CI job cannot pass even after the migration setup conflict is isolated.
    - at: '2026-09-23T03:58:56Z'
      actor: {role: owner, session: user}
      kind: deferred
      note: >-
        Owner narrowed O-016 to focused Windows migration setup coverage and cross-builds; the full Windows runtime suite is not required. T-018 and T-019 were removed. Keep I-030 open because these failures remain uncorrected; revisit only if the full suite returns to scope.
    - at: "2026-09-25T09:26:40Z"
      actor:
        role: owner
        session: board-owner
      kind: owner_decision
      note: Moved from open to in_progress by the owner from the board.
    - at: "2026-09-25T09:26:42Z"
      actor:
        role: owner
        session: board-owner
      kind: owner_decision
      note: Moved from in_progress to open by the owner from the board.
    - at: "2026-09-25T10:36:09Z"
      actor:
        role: executor
        session: i030-windows-repair-20260925
      kind: repair_attempted
      note: >-
        Owner reinstated full Windows support. A native Windows run (go1.26.2
        windows/amd64, NTFS copy of the working tree) found 17 failing tests,
        not two, in six causes, all now repaired. (1) Directory metadata
        writes, including those caused by registering a watch, reloaded the
        board: watch.go now ignores Write/Chmod events on directories, and the
        unchanged-upgrade snapshot reads directory times with Lstat instead of
        the parent listing. (2) Case-insensitive filesystems: guide-casing
        tests check the exact directory entry, and case-collision fixtures skip
        via testutil.SkipIfCaseInsensitive. (3) Path messages: tests expect
        native separators. (4) Permissions: distribution archives set member
        mode 0755 explicitly; the 0600 assertion is skipped on Windows. (5)
        Lock-release failure injection holds a second handle on Windows. (6)
        The real-Git test repo disables core.autocrlf. The Windows CI job now
        runs the full suite. Evidence: full Windows suite passed twice (go run
        ./internal/buildtool test -json -count=1 ./...); make build and make
        test-full passed on Linux go1.26.2; git diff --check clean. Hosted
        windows-latest CI evidence is pending until pushed.
    - at: "2026-09-25T10:42:38Z"
      actor:
        role: executor
        session: i030-windows-repair-20260925
      kind: repair_attempted
      note: >-
        Hosted windows-latest CI run 36125057404 failed: its temp root is the
        short 8.3 path C:\Users\RUNNER~1, and V2 discovery compared each
        record's resolved path with the unresolved root, so every project
        under a short name or symlinked directory was refused as "resolves
        outside the project root". The local NTFS copy had no short name, so
        it did not show this. Discovery now resolves the root the same way
        (newV2PathConfiner), with a symlinked-root regression test that fails
        without the fix. Evidence: full Windows suite passed with TMP set to
        a short 8.3 path (reproducing 343 failures before the fix) and from
        the normal temp path; make build and make test-full passed on Linux.
---

# I-030: Windows full suite has platform failures

## Summary

T-016 adds a complete Windows test job, but the current Windows suite has two
additional reproducible platform failures. O-016 widened the supported gate to
this path and explicitly promises Windows execution, so a configured job that
is known to fail cannot satisfy the Objective.

## Current Scope

The owner later narrowed O-016 to focused Windows migration setup coverage and
cross-builds. T-018 and T-019 were removed from O-016, and the full Windows
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
- These are supported Windows paths covered by T-016's CI acceptance criterion,
  CFG-02, and TEST-08; they are not remote hypothetical configurations.

## Proof Needed

- No repair is required by the current O-016 scope. If the full Windows runtime
  suite is reinstated, repair the watcher and guide-casing failures and rerun
  the focused reproductions plus the full Windows suite.
