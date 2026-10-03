---
id: T-108
title: "Remember the parallel planning preference"
objective: O-037
status: planned
depends_on: []
owner_validation: {required: false}
planned_by: {role: planner, session: codex-o033-advisory-replan-2026-10-03}
complexity_tier: medium
complexity_reason: "Safe persistent preferences must preserve authored YAML and refuse stale writes."
planned_reads:
  - "internal/data/config.go"
  - "internal/data/config_test.go"
  - "internal/data/write_splice.go"
  - "internal/data/write.go"
  - "internal/data/project.go"
planned_writes:
  - "internal/data/config.go"
  - "internal/data/config_test.go"
  - "internal/data/feature_preferences.go"
  - "internal/data/feature_preferences_test.go"
---

# Remember the parallel planning preference

## Outcome

Projects retain one optional parallel-planning preference, default off, without changing any Code Health or lifecycle policy.

## User Check

Toggle using the internal writer fixture, reload the project preference, and verify unrelated configuration/comments are unchanged.

## Done When

- Decode optional features.parallel_planning as a boolean, default false when absent; preserve all current Code Health behavior and existing required configuration diagnostics.
- Expose one shared saved preference to board/resume consumers, and a narrow data-owned writer that changes only this key while preserving unrelated keys/comments/authored bytes.
- Refuse stale-source conflicts and malformed configuration; safe file replacement and failed writes leave a clear diagnostic and do not report an unsaved toggle as saved. Repeated unchanged writes are no-ops.
- No registry, generic settings editor, Code Health option, lane gates or project opt-in/backfill. Settings enable presentation/planning advice only.
- Record per-criterion evidence and a fresh full gate; optional Task Check waiver is recorded only on explicit owner instruction. Mandatory independent Full Objective Check remains required.

## Context Files

- `internal/data/config.go`
- `internal/data/config_test.go`
- `internal/data/write_splice.go`
- `internal/data/write.go`
- `internal/data/project.go`

feature_preferences.go and feature_preferences_test.go are explicitly new implementation targets, not missing Context Files.

## Design References

O-037 Confirmed Design — 2026-10-03: preference storage and Advanced Options delivery before lane implementation. Project Design sections 1, 8 and 9 describe existing ownership, rendering and concurrency boundaries.

## Guardrails

FS-01, FS-04, FS-05, FS-06, DATA-01, DATA-02, DATA-03, ARCH-02, ARCH-03, CFG-01, CFG-02, CFG-03, TEST-01, TEST-02, TEST-03, TEST-04, TEST-08, TEST-09, STYLE-01..10.

## Implementation Plan

1. Define typed preference values and pure decoding in config ownership, with named malformed-value diagnostics.
2. Implement an explicit freshness-checked narrow writer in new feature_preferences.go, reusing existing preservation conventions without rewriting the whole Config struct.
3. Verify absent/false/true values, malformed YAML/value, preservation, repeat no-op, changed-source refusal and unwritable/replacement failure in temporary-project tests.
4. Publish the exact read/write interface and source-conflict behavior for the board and shared Next projection.

## Boundaries

Implement the scoped preference/screen only. Anticipated manifests and lane choice remain advice, not execution allowlists or gates. Retain ordinary extra-read logging and return REPLAN REQUIRED only for an actual material plan/acceptance gap. Code Health, completion authority, existing dependency gates and actual-worktree rules stay unchanged. No agent settings toggling for the owner's real project during execution; use temporary fixtures. No Check or owner waiver is invented.

## Technical Verification

Fresh make test-full for configuration replacement/path and reload-watch behavior. Focused tests are iteration only. Native windows-tests CI evidence is produced by repository CI and supplied by the owner before the independent Full Objective Check under agent-skills/references/check-method.md. Owner validation is required for the settings screen.

## Technical Evidence

Pending execution: record actual reads/writes, logged extra reads, per-criterion cases/results, command/time/toolchain, full gate and limitations. No self-clearance.

## Drift Notes

Record material implementation/acceptance deltas for planner reconciliation. Ignoring a suggested lane or differing from an anticipated manifest alone requires no replan.

