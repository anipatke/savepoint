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
2. The run key `s` is offered only when the feature is on, the focused record is a Task, its Next verb is Start/Build/Test, and the same gate `actionsForRecord` uses allows it; a Blocked Task cannot start a run; a second run cannot start while one is `running` (and not inferred `runner_lost`).
3. The board spawns the runner via `os.Executable()` detached and holds no child handle; closing the board does not end the run; reopening shows the run with its newer events.
4. The watcher includes `runs`; the activity view reads `events.jsonl` incrementally from the last offset, keeps a bounded in-memory window, coalesces bursts through the existing debounce, and stays responsive on a 50 MB event file.
5. The overview and activity view follow the Atari-Noir sketches in review §4: real events only, no percentages; orange running, purple agent-reported, red needs-you, dim queued.
6. Cancel writes the cancel file and the row reaches `cancelled` without the board touching the process.
7. Configured gates pass, including `make test-full`.

## Architectural Considerations

- Follow the existing overlay pattern (`Issues`, `Health`, `Options`): a nullable overlay struct on `Model`, mutually exclusive with the others, origin restored on close.
- All reads are `tea.Cmd`s; the reducer only appends; no filesystem access in `Update` (ARCH-02).
- Run state shown on the board is an inference from files plus PID liveness; the board never rewrites `run.yml`.
- The board imports `internal/run` for reading records, decoding events (`ParseClaudeLine`), `RunnerLost`, and the detached spawn; it never imports `codehealth` and never parses provider lines itself.

## Confirmed Design

Owner, 2026-10-11, design confirmation in a planning session:

- **Keys.** `s` starts a run on the focused Task; `A` opens the Runs overlay; `enter` on a row opens the activity view; `c` cancels the focused or viewed run. All three were unused in the V2 board (`a` is owner acceptance, so the plan's `a` was rejected). With the flag off none of them is bound, hinted, or listed in help.
- **Detached start.** A new `internal/run` helper starts `os.Executable() run --task T-### <project root>` detached (`Setsid` on Unix; `DETACHED_PROCESS|CREATE_NEW_PROCESS_GROUP` on Windows), stdin from the null device, stdout and stderr to `.savepoint/runs/launch.log` (replaced on each start), then releases the process. The board holds no child handle. A refusal before the run directory exists (for example the Claude version floor) is shown from that log when no new run appears within a few seconds, never silently dropped.
- **Board-side gating.** The board checks the flag, the Task gate and Next verb, and that no run is `running` (and not lost) before spawning; the runner's own checks in `Launch` stay authoritative.
- **Reading runs.** Read-only helpers in `internal/run` list run directories and read `events.jsonl` from a byte offset, returning decoded events and the new offset; the board keeps a bounded window (500 events) per viewed run. Reads are `tea.Cmd`s; `runs` joins `v2WatchDirectories` only when the flag is on.
- **Presentation.** Real events only, no percentages; orange running, purple finished with agent-reported outcome, red needs-you (any denials), dim otherwise. The final message sits under an "agent reports" label; the board adds no verdict.
- **Cancel.** `c` writes `<run>/cancel` and nothing else; the row reaches `cancelled` from the runner's record.
- **Out of this Objective.** A running badge on Task cards; re-run, takeover, permission editing.
- **Tasks.** Four sequential Tasks sharing board files, no lanes: feature switch with flag-off golden parity (fixtures captured before any other change); start a run with `s`; Runs overlay and activity view; cancel and needs-you/agent-reports presentation.

## Boundaries

**In scope:**
- `internal/board/v2`: new runs overlay and activity view, watcher root, keys, help, hints; config gate; golden tests
- `internal/run`: detached spawn (Unix and Windows files) and read-only run listing and offset event reading for the board

**Out of scope:**
- Re-run with a wider allow-list; interactive takeover; permission editing; Check runs; any provider other than Claude Code
