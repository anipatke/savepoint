---
id: I-092
title: Discovered health reports never reach collection
type: defect
status: open
source:
  kind: check
  check: C-941
  actor: {role: checker, session: o028-check-20261001}
  at: '2026-10-01T08:38:18Z'
tasks: [T-056, T-057, T-058]
checks: [C-941]
guardrail_ids: [TEST-01, TEST-02]
history:
  - at: '2026-10-01T08:50:00Z'
    actor: {role: executor, session: o028-repair-20261001}
    kind: repair_attempted
    note: "Collect now reads the report file a suggested tool writes (Report path) instead of stdout; own reports under .savepoint/health/reports are cleared before and after each run and the directory is created. Tests: TestCollectReadsTheReportASuggestedToolWrites (Lizard, jscpd, OSV: success, silent run, nonzero exit). Gates: make build, make test-fast, make test-full (go1.26.2 linux/amd64) pass; awaiting a checker."
    check: C-941
---

# I-092: Discovered health reports never reach collection

## Summary

O-028 promises supported tools can be discovered, configured, and run. Discovery's three executed providers write reports to fixed project paths, but Collect reads only stdout unless an argument is exactly `{report}`. None of those proposals contains that placeholder. This breaks the supported setup-to-collection path and bypasses the promised temporary report cleanup.

## Evidence

C-941 frozen M5, independent harness `M5 discovered executed reports reach reader`: Discover a Go project with all executables found; use each proposal as Collect's config; a fake tool obeys the proposed file output path, writes `actual-report`, and emits `console-status` to stdout. A reader checks which bytes arrive. Lizard, jscpd, and OSV all receive `console-status`, produce `failed`, and leave the actual report inside the repository. This proves the orchestration mismatch without depending on real report normalization, which belongs to O-029.

Locations: `internal/codehealth/discovery_catalogue.go:100` constructs all three fixed output arguments; `internal/codehealth/discovery.go:402` assembles proposals; `internal/codehealth/collect.go:309` returns stdout when no temporary report was prepared; `collect.go:323` recognizes only the literal placeholder. Expected: the output selected by setup reaches the reader through a compatible bounded, cleaned-up report channel. Actual: configured file reports are ignored.

Existing `TestDiscoverExecutedToolArgsTranslateExclusions` asserts proposal args in isolation; `TestCollectWithRealProcessesAndReportFiles` manually supplies compatible placeholder args. Neither connects them.

## Proof Needed

Exercise Discover → Plan/Apply → Collect for all three executed providers, with successful file output, missing/malformed output, nonzero exit and cancellation cleanup. Keep jscpd's output-directory convention compatible with collection. No production reader is required to prove delivery. Run a fresh full gate and recheck frozen M5 and its named adjacent cells.
