---
id: O-048
title: Animate Tiles
status: planned
depends_on: []
release: G-003
priority: high
rank: 2
---

# O-048: Animate Tiles

## Outcome

A Task tile whose recorded status is `in_progress` at stage `build` or `test` shows a one-cell spinner and its stage in words, "Building…" or "Testing…", in place of today's stage badge. At `audit` (shown as check) or `done`, the tile shows today's badges. A board with no such Task is byte-identical to v2.3.0. The unused `savepoint run` command, the `internal/run` package, and doctor's lost-runner report that O-045 added are removed, and `agent_launcher` goes back to the parsed-but-unused config block it was before O-045.

## Why

The owner scrapped Agent Orchestration Spike 01 on 2026-10-11 (Goal G-003 Decision). The one part worth keeping is seeing from the main board which work is under way. Shipping the run command without the board that used it would leave a half-feature to maintain.

## Success Conditions

1. A tile at `build` shows a moving spinner with "Building…", and a tile at `test` shows one with "Testing…". The words come from the stage recorded in the Task file. The board does not claim it verified them.
2. A tile at `audit`, `done`, or `planned` shows exactly today's badges.
3. The spinner ticks only while at least one shown tile is at `build` or `test`. An idle board runs no timer. `Update` and rendering do no file work.
4. T-131's committed golden files (`internal/board/v2/testdata/golden/`) pass unchanged. Plain (non-TTY) output, help, the key map, and footer hints are unchanged, and no key is added.
5. `savepoint run`, `internal/run`, doctor's lost-runner report, and O-045's `agent_launcher` additions are gone. AGENTS.md's Codebase Map no longer names them, and existing `agent_launcher` config still loads.
6. `make build && make test-fast` and fresh `make test-full` pass, with the native Windows CI job green for the Full Objective Check.

## Architectural Considerations

- The signal is the Task's recorded `status` and `stage` only, the same values `stageBadge` reads. There is no liveness inference, file-age timeout, or run record. A Task left at `build` keeps animating until someone moves it on (owner decision).
- Animation state is a frame counter on `Model`, advanced by a tick `tea.Cmd` that stops by itself when no shown tile qualifies (ARCH-02).
- The spinner is the owner-approved exception to the visual identity's "running animation skip" rule. `.savepoint/visual-identity.md` records it (TPL-02).
- Removal leaves `docs/reviews/agent-orchestration-*.md` and O-045's records as history.

## Confirmed Design

Owner, 2026-10-11, replanning after the spike was scrapped:

- **Signal.** Any Task at `build` or `test` animates, not only the router's selection, and with no recent-change timeout.
- **Run code.** Remove the `savepoint run` command and its package in this Objective rather than ship it unused.
- **Old records.** O-046 and O-047, with Tasks T-131 to T-136, are deleted. Their IDs are retired, and git history keeps their committed versions.
- **Tasks.** Two independent Tasks: animate the tiles, and remove the run command. They touch different code and may run in either order or side by side.

## Boundaries

**In scope:**
- `internal/board/v2`: tile badge, spinner frame state, tick command, tests
- `.savepoint/visual-identity.md`: the spinner exception
- Removal of `cmd/run.go`, `main.go` run dispatch, `internal/run`, `internal/doctor/lost_runs.go`, O-045's `internal/data/launcher_config.go` additions, related tests, and the AGENTS.md Codebase Map lines

**Out of scope:**
- Any run, agent launch, or board-started work; any new key, option, or overlay; animation on Objective rows or in plain output; changes to gates, `ResolveNext`, or the workflow skills
