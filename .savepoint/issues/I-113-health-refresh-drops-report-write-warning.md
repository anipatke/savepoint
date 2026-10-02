---
id: I-113
title: Popover refresh silently drops report-write warnings
type: defect
status: resolved
source:
  kind: check
  check: C-955
  actor: {role: checker, session: check-o036-20261002-independent}
  at: '2026-10-02T03:01:24Z'
tasks: [T-083]
checks: [C-955, C-956]
guardrail_ids: [TEST-01, TEST-02]
resolution:
  disposition: verified
  check: C-956
  actor: {role: checker, session: recheck-o036-20261002-independent}
  at: '2026-10-02T03:12:42Z'
  reason: Repair independently verified against the frozen C-955 scope.
history:
  - at: '2026-10-02T03:12:42Z'
    actor: {role: checker, session: recheck-o036-20261002-independent}
    kind: rechecked
    note: 'Actual production command bridge preserves a manual snapshot and reports the write failure separately; warning remains visible after reload at all original sizes.'
    check: C-956
---

# I-113: Popover refresh silently drops report-write warnings

## Summary

O-036 success condition 10 and T-083 Done When 1 require a failed report write to be reported as a warning while the saved snapshot and successful refresh stand. The official check does this; the real popover refresh drops the warning entirely. An owner can consequently hand an agent an older report without knowing the latest refresh did not update it.

## Evidence

`internal/board/v2/io.go:648` discards the Collection returned by Collect. Lines 654–655 explicitly discard Collection.ReportErr. `internal/board/v2/health.go:73` carries only a fatal Err in healthRefreshDoneMsg; lines 273–284 reload the dashboard on nil without any warning. The configured refresh writes through Collect, so this is the supported R path.

Independent reproduction (C-955 frozen matrix Popover refresh / actual report failure): initialize a temporary Git repository, make `.savepoint/health/report.md` a directory, save an empty valid health configuration, then call production `refreshHealth(context.Background(), <root>/.savepoint, nil)`. Actual: nil returned; one valid manual snapshot retained; no warning is available to the completion message. Expected: snapshot retained and refresh succeeds, with a visible report-write warning. This obstruction models a real write error without permissions/root/platform assumptions.

`TestO036IndependentActualRefreshReportFailure` ran through a Go overlay and verified the actual function and saved snapshot. Existing TestCollectReportWriteFailureKeepsTheSnapshot proves ReportErr is available; TestRun_reportWriteFailureIsOnlyAWarning proves the official-command side. Neither drives a failed report write through production popover refresh and observes its warning.

## Proof Needed

Carry the report-write warning through refresh completion separately from fatal collection failure; show it to the owner while keeping the snapshot and reloading results. Test the actual refresh/command bridge and visible completion with both successful and blocked report writes, then rerun the frozen C-955 matrix and full gate. Keep all owner-completed Tasks done.

## Recheck Proof

C-956: Actual production command bridge preserves a manual snapshot and reports the write failure separately; warning remains visible after reload at all original sizes. Fresh full and focused gates pass.
