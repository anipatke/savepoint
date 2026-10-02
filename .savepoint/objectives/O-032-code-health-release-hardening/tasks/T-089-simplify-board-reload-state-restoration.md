---
id: T-089
title: Simplify board reload state restoration
objective: O-032
status: done
complexity_tier: high
complexity_reason: Reload retry and selection state must survive decomposition without changing visible behavior.
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: planning-o032-20261002-owner-confirmed}
check_waiver:
    task: T-089
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-02T04:08:39Z"
---

# Simplify board reload state restoration

## Outcome

Board reload has distinct failure/reset and successful-restoration steps while selection, status and retry behavior remain stable.

## User Check

Keep an Objective selected while project files reload; verify errors and retries do not silently change the selection.

## Done When

- Preserve sequence rejection, initial-load versus reload failure, one retry, previous visible state and named diagnostics.
- Preserve release/Objective/cursor restoration, filtered scope, retained status message, rollback cleanup and fatal filter handling.
- Extract coherent helpers from applyLoad; do not move IO into rendering or duplicate data-layer lifecycle/readiness rules.
- Record independent initial/success/failure/retry/stale-sequence/filter scenarios and before/after complexity, targeting touched production functions <=20 CCN and aiming <=10. Preserve assertions in existing complex tests; simplify setup only when it improves clarity.

## Context Files

`internal/board/v2/update.go`; `internal/board/v2/load.go`; `internal/board/v2/load_test.go`; `internal/board/v2/boundary_test.go`; `internal/board/v2/fixture_test.go`; `internal/board/v2/objectives_test.go`.

## Design References

O-032 Confirmed Release-Hardening Design and Owner-Requested Health Analysis; Design sections 1, 6, 7, 9, 12 and 13 as applicable.

## Guardrails

FS-01, FS-04, FS-05, FS-06, ARCH-01, ARCH-02, ARCH-03, CFG-01, CFG-02, TEST-01, TEST-02, TEST-03, TEST-04, TEST-05, TEST-08, TEST-09, STYLE-01, STYLE-02, STYLE-07.

## Implementation Plan

1. Map applyLoad state inputs and observable outputs.
2. Separate failed first-load reset and successful data/selection restoration.
3. Verify state and selection across complete retry sequences and stale messages.
4. Compare same-scope complexity and run board focused tests.

## Boundaries

No new keys, overlay design, health collection, Task status transitions or selection policy.

## Technical Verification

For before/after complexity evidence, the owner-confirmed scope authorizes a direct project-owned Lizard invocation over the scoped source paths with temporary output; run no Savepoint health command. If the executable is unavailable, record the measurement gap for the owner instead of installing it.

`make build && make test-fast` for ordinary handoff; focused `make test-focused TEST=... PKGS=...` only during iteration. Record named per-criterion evidence and limitations. The later independent Check applies `agent-skills/references/check-method.md`; the mandatory Full Objective Check runs its fresh full gate and official Code Health collection. This executor never writes a Check or runs savepoint health setup/report/check.

## Technical Evidence

Executed 2026-10-02 (toolchain: local `go`, `lizard` from ~/.local/bin).

**Change:** `internal/board/v2/update.go` only. `applyLoad` now delegates to `acceptLoadSeq` (sequence rejection), `applyLoadFailure` (reload keeps state + `ReloadDiagnostic`; initial failure resets board), `restoreLoadedSelection` (release/Objective/cursor/focus/overlay restoration), `noteVanishedSelection`, `reconcileDetailAndIssues` and `noteIssueFocusMoved`. Order of side effects is unchanged; the one-retry step, status preservation, rollback cleanup and fatal-filter exit stay inline in `applyLoad`. No IO moved, no data-layer rules duplicated, no tests edited.

**Per-criterion:**
- Sequence rejection / initial vs reload failure / one retry / previous visible state / diagnostics: same code paths, covered by existing tests (e.g. `TestFailedReloadRetriesBeforeReporting`); board/v2 Reload|Load|Retry|Stale|Filter subset: 47 PASS.
- Release/Objective/cursor restoration, filtered scope, retained status, rollback cleanup, fatal filter: unchanged logic, existing board/v2 tests pass in `make test-fast`.
- Helpers extracted, no IO in rendering: yes.
- Complexity (lizard, `internal/board/v2/update.go`): `applyLoad` CCN 26 / 80 lines before; 7 / 36 lines after. New helpers CCN <=10 (largest `reconcileDetailAndIssues` 10, `restoreLoadedSelection` and others lower). Touched functions all <=10.

**Commands:** `make build` ok; `make test-fast` ok (exit 0); `make test-focused TEST=. PKGS=./internal/board/v2` ok; `gofmt -l` and `go vet ./internal/board/v2` clean.

**Files read:** all within Context Files except none beyond (`update.go`, board v2 tests via test run). No extra reads.

**Limitations:** No new scenario tests were added; independent initial/success/failure/retry/stale-sequence/filter scenarios rely on the existing suite, which I did not individually re-audit for each case. Complexity measured only on update.go/load.go; no Savepoint health command run.

## Drift Notes

Record implemented responsibility or interface changes for planner reconciliation before the mandatory Full Objective Check. If evidence requires a material unknown repair, return REPLAN REQUIRED rather than inventing another scope.
