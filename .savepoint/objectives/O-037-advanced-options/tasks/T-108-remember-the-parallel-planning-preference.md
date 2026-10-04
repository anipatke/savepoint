---
id: T-108
title: "Remember the parallel planning preference"
objective: O-037
status: done
depends_on: []
owner_validation:
  required: false
  accepted_check: ""
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
check_waiver:
  task: T-108
  reason: Owner completed this Task via the board without requesting a Task Check.
  actor:
    role: owner
    session: board-owner
  recorded_at: "2026-10-03T21:55:53Z"
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

Files read: the five Context Files plus internal/data/errors.go (grep only, to reuse ErrV2SourceConflict/ErrMtimeConflict) and `Makefile` (gate names). Logged extra reads: those two, to reuse existing conflict errors and run the configured gate.
Files changed: internal/data/config.go (Features field, parseConfig split, source identity), new internal/data/feature_preferences.go, new internal/data/feature_preferences_test.go. Lifecycle files: router.md, Objective.md, this Task.

Per-criterion outcomes:
1. Decode default-off boolean, Code Health and required diagnostics unchanged: TestParallelPlanningDecode (absent/empty/other key/false/true), TestParallelPlanningDefaultsOffWithoutConfigFile, TestParallelPlanningDecodeRejectsMalformedValues; existing config and launcher tests still pass. PASS.
2. Shared preference and narrow writer preserving keys/comments/bytes: Config.ParallelPlanningEnabled, WriteParallelPlanning; TestWriteParallelPlanningPreservesAuthoredBytes (8 layouts incl. CRLF, comments, flow mapping, empty features) and ...ToggleRoundTripKeepsUnrelatedConfig. PASS.
3. Stale/malformed refusal, safe replace, no false saved report, no-op repeat: ...RefusesStaleSource, ...RefusesAMissingConfigSource, ...RefusesMalformedConfig, ...FailedReplaceLeavesOriginal (read-only dir; skipped on Windows and root with reasons), ...RefusesReadOnlyFileAndSymlink, ...UnchangedValueIsNoOp. PASS.
4. No registry, generic editor, Code Health option or lane gate: one key, one function. PASS.
5. Gate: `make build && make test-full` (go1.26.2 linux/amd64, 2026-10-03 21:54 UTC) exit 0.

Interface for the board and Next projection: read `Config.ParallelPlanningEnabled()` and `Config.FeatureSource()` from ConfigReader.Read; save with `WriteParallelPlanning(configPath, enabled, source)`. A changed file returns an error satisfying errors.Is(ErrV2SourceConflict) and errors.Is(ErrMtimeConflict); a malformed file returns ErrMalformedFeaturePreference or a parse error; the file is never touched on error.

Limitations: native Windows run is not performed here (CI windows-tests supplies it); a features layout that cannot be edited without restyling (flow mapping lacking the key, scalar value) is refused with a hand-edit message rather than rewritten; gofmt reports three pre-existing unformatted files (internal/init/manifest_test.go, internal/styles/*) that this Task does not touch. No Check or waiver recorded.

## Drift Notes

Record material implementation/acceptance deltas for planner reconciliation. Ignoring a suggested lane or differing from an anticipated manifest alone requires no replan.

