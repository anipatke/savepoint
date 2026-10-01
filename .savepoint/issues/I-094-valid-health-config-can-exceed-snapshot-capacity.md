---
id: I-094
title: Valid health configuration can exceed snapshot capacity
type: defect
status: open
source:
  kind: check
  check: C-941
  actor: {role: checker, session: o028-check-20261001}
  at: '2026-10-01T08:38:18Z'
tasks: [T-055, T-057]
checks: [C-941]
guardrail_ids: [CFG-01, TEST-01, TEST-02]
history:
  - at: '2026-10-01T08:50:00Z'
    actor: {role: executor, session: o028-repair-20261001}
    kind: repair_attempted
    note: "Config.Validate now counts the not_configured placeholders against MaxResults, so an over-capacity configuration is refused before any tool runs. Tests: TestConfigCapacityCountsPlaceholdersForUnconfiguredCapabilities, TestCollectRefusesAnOverfullConfigurationBeforeRunningTools (exact capacity saves). Gates: make build, make test-fast, make test-full (go1.26.2 linux/amd64) pass; awaiting a checker."
    check: C-941
---

# I-094: Valid health configuration can exceed snapshot capacity

## Summary

Collect promises one saved snapshot with each configured instance and a not_configured entry for each missing capability. Config accepts 32 instances, but the snapshot accepts only 32 total results. A supported many-instance configuration can run every tool and then fail to save any evidence.

## Evidence

C-941 frozen M10 independent probe: create 32 named Lizard instances with distinct scopes, validate the config successfully, and Collect with a successful fake runner/reader. All 32 tools run. Save fails: `detail exceeds bounds: results: has 36 entries; at most 32`. Expected: validation and collection agree on capacity, and valid config can produce the promised snapshot. Actual: four unconfigured capabilities overflow the bound after execution.

`internal/codehealth/config.go:59` checks only configured entries against MaxResults; `collect.go:122` runs the configured instances; `collect.go:126` appends missing-capability placeholders; `collect.go:159` saves only afterward. Existing sequencing test covers two configured instances, not the capacity boundary. Executor evidence disclosed the limitation but provided no owner exception.

## Proof Needed

Align validation and collection capacity without silently dropping an instance or missing-capability result. Test exact final-result capacity, immediately below/above, one-capability many-instance configs and all-capability configs. Invalid inputs must fail before executing tools. Run a fresh full gate and frozen M10 recheck.
