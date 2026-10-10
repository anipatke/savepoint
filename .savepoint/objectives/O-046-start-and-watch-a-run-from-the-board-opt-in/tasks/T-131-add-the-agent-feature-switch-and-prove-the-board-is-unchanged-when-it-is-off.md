---
id: T-131
title: Add the agent feature switch and prove the board is unchanged when it is off
objective: O-046
status: planned
depends_on: []
owner_validation: {required: false}
planned_by: {role: planner, session: planning-o046-20261011}
---

# Add the agent feature switch and prove the board is unchanged when it is off

## Outcome

The board knows whether `agent_launcher.enabled` is on, read once per load, and a golden test proves that with it off or absent the rendered TUI board, key handling, help overlay, footer hints, and plain (non-TTY) output are byte-identical to the board before O-046.

## User Check

None required beyond the golden test; the owner sees no change.

## Done When

1. Golden fixtures of the current board are captured **before any other O-046 change** from a temporary fixture project (TEST-04): TUI board at a fixed size with Task and Objective focus, help overlay for each, footer hints, and plain output. Their source commit is recorded in Technical Evidence.
2. The board's load reads `agent_launcher.enabled` through `internal/data`'s existing config reader into load state; rendering performs no IO (ARCH-02). An invalid launcher config surfaces the existing named config diagnostic, not a silent default (CFG-01, DATA-03).
3. Golden tests assert byte-identity for flag absent and flag `false`, and a test proves the flag value reaches the model when `true` (no visible change yet).
4. `make build && make test-fast` and fresh `make test-full` pass. Per-criterion evidence follows AGENTS.md's Verification Policy, with an explicit owner waiver if the optional Task Check is skipped.

## Context Files

`internal/board/v2/load.go`, `internal/board/v2/load_test.go`, `internal/board/v2/model.go`, `internal/board/v2/view.go`, `internal/board/v2/help.go`, `internal/board/v2/plain.go`, `internal/board/v2/fixture_test.go`, `internal/board/v2/footer_test.go`, `internal/board/v2/view_test.go`, `internal/data/config.go`, `internal/data/launcher_config.go`, `.savepoint/objectives/O-046-start-and-watch-a-run-from-the-board-opt-in/Objective.md`.

## Design References

Design §8 (TUI, Advanced Options read-once pattern, board persistence and refresh). O-046 Confirmed Design.

## Guardrails

ARCH-02, ARCH-03, CFG-01, DATA-03, TEST-01, TEST-02, TEST-04, TEST-08.

## Implementation Plan

1. Build the fixture project and write the golden files from the unchanged code; commit them first.
2. Add the flag to the load message and Model, following how `features.parallel_planning` is carried.
3. Add golden comparisons for absent and false, and a model test for true.
4. Gates.

## Boundaries

No keys, overlays, watcher roots, or spawning; no change to `resume` or `doctor`.

## Technical Verification

Golden byte comparisons plus focused tests; handoff gates per AGENTS.md's Verification Policy. Evaluated later under `agent-skills/references/check-method.md`.

## Technical Evidence

Pending execution.

## Drift Notes

Record architecture deltas and reconcile through the planner before Check.
