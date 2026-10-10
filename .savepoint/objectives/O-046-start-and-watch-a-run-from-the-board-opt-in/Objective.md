---
id: O-046
title: Start and watch a run from the board, opt-in
status: planned
depends_on: [O-045]
release: G-003
priority: high
rank: 2
---

# O-046: Start and watch a run from the board, opt-in

## Outcome

When `agent_launcher.enabled` is true, a key on a Task whose Next is Start, Build, or Test and whose gate allows it starts `savepoint run` detached; a Runs overlay lists runs with state, elapsed time, and last event; a live activity view shows real events; a key cancels. A finished run shows the agent's message under an "agent reports" label and no verdict of its own. A run with denials reads "needs you" and names the refused actions. When the flag is false or absent, none of this exists and the board is byte-identical to today.

## Why

This is the step that removes the terminal switch. It must be an addition beside the paste workflow, not a replacement, and the board must stay a reader of files so a run never depends on the board being open (review §3 and the opt-in rule in `docs/reviews/agent-orchestration-pass-1.md`).

## Success Conditions

1. Golden test: with the flag off or absent, rendered board, key map, help overlay, status hints, and plain (non-TTY) output match v2.3.0 for the same project.
2. The run key is offered only when the feature is on, the focused record is a Task, its Next verb is Start/Build/Test, and the same gate `actionsForRecord` uses allows it; a Blocked Task cannot start a run; a second run cannot start while one is `running`. The key itself is chosen at planning time after confirming it is unused in the V2 board.
3. The board spawns the runner via `os.Executable()` detached and holds no child handle; closing the board does not end the run; reopening shows the run with its newer events.
4. The watcher includes `runs`; the activity view reads `events.jsonl` incrementally from the last offset, keeps a bounded in-memory window, coalesces bursts through the existing debounce, and stays responsive on a 50 MB event file.
5. The overview and activity view follow the Atari-Noir sketches in review §4: real events only, no percentages; orange running, purple agent-reported, red needs-you, dim queued.
6. Cancel writes the cancel file and the row reaches `cancelled` without the board touching the process.
7. Configured gates pass, including `make test-full`.

## Architectural Considerations

- Follow the existing overlay pattern (`Issues`, `Health`, `Options`): a nullable overlay struct on `Model`, mutually exclusive with the others, origin restored on close.
- All reads are `tea.Cmd`s; the reducer only appends; no filesystem access in `Update` (ARCH-02).
- Run state shown on the board is an inference from files plus PID liveness; the board never rewrites `run.yml`.
- Tasks are detailed only after O-045 lands and this Objective's design is confirmed.

## Boundaries

**In scope:**
- `internal/board/v2`: new runs overlay and activity view, watcher root, keys, help, hints; config gate; golden tests

**Out of scope:**
- Re-run with a wider allow-list; interactive takeover; permission editing; Check runs; any provider other than Claude Code
