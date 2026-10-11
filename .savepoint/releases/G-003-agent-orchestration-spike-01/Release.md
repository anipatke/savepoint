---
id: G-003
title: Animate Tiles
status: planned
---

## Decision (2026-10-11)

The owner stopped Agent Orchestration Spike 01 and scrapped its board code: run start, the Runs overlay and activity view, cancel, and the Agent runs option. The first real run on T-134 wrote code but could run no command, because a project with no configured build or test commands gives a headless run an empty allow-list, and the board gave no warning. The branch is kept for one smaller feature: animated "Building…" and "Testing…" labels on Task tiles. The planner replanned the Goal the same day:

- O-046 (board runs) and O-047 (validation) were deleted with Tasks T-131 to T-136. Their IDs are retired, and git history keeps their committed versions.
- O-045 (the detached run command) stays done as history. Its code is removed by O-048.
- O-048, Animate Tiles, is the Goal's one remaining Objective.

The Goal keeps its stable ID and path. Its original spike outcome is in git history and in `docs/reviews/agent-orchestration-spike-01-plan.md`.

## Outcome

On the board, a Task tile that is building or testing shows a small moving spinner with "Building…" or "Testing…", so the owner can see at a glance which work is under way. The words come from the stage recorded in the Task file. Every other tile, and a board with nothing building or testing, looks exactly as in v2.3.0. The unused agent-run command from the spike is gone. Pasting the Next line into an agent terminal remains the way work is done.

## Why

The spike showed that driving agents from the board needs more groundwork than it was worth now. The one part the owner wants is visible progress on the main board, without new commands or records.

## Success Conditions

- O-048 is done: tiles animate at build and test, idle boards run no timer, and the committed golden board snapshots pass unchanged.
- No `run` command, run package, or lost-runner report ships, and existing `agent_launcher` config still loads.
- `make test-full` passes, with native Windows CI green for the Full Objective Check.

## Boundaries

**In scope:** the tile spinner and labels in `internal/board/v2`; the visual-identity exception for it; removal of O-045's run code.

**Out of scope:** starting, watching, or cancelling agent runs; run records; liveness inference; any change to gates, `ResolveNext`, the workflow skills, or plain output.
