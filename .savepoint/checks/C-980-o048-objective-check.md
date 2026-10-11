---
id: C-980
scope: {kind: objective, id: O-048}
result: NEEDS WORK
checked_by: {role: checker, session: check-o048-20261011}
executed_session: separate-executor-sessions-T-137-T-138
checked_at: '2026-10-11T05:45:00Z'
health_snapshot: sha256:51e37a46de53139ee91d317bd275c4cea472101c52cd5a855ebae73808342533
reviewed:
  base_commit: 12e63a766ac24215358b9ed3eaae659a69bcd573
  head_commit: 12e63a766ac24215358b9ed3eaae659a69bcd573
  files:
    - .savepoint/visual-identity.md
    - AGENTS.md
    - agent-skills/savepoint-task/SKILL.md
    - templates/project-v2/agent-skills/savepoint-task/SKILL.md
    - main.go
    - cmd/run.go
    - cmd/run_test.go
    - internal/board/v2/badges.go
    - internal/board/v2/card.go
    - internal/board/v2/column.go
    - internal/board/v2/model.go
    - internal/board/v2/update.go
    - internal/board/v2/view.go
    - internal/board/v2/spinner.go
    - internal/board/v2/spinner_test.go
    - internal/board/v2/columns_view_test.go
    - internal/board/v2/objectives_test.go
    - internal/board/v2/releases_test.go
    - internal/data/launcher_config.go
    - internal/data/launcher_config_test.go
    - internal/doctor/lost_runs.go
    - internal/doctor/lost_runs_test.go
    - internal/doctor/v2_runtime.go
    - internal/init/claude_guard_test.go
    - internal/run/
  dependencies: []
issues: [I-155, I-156, I-157]
unmet: [O-048-SC1, O-048-SC3, O-048-SC6, O-048-Boundaries, T-137-DW2, T-137-DW5, T-138-DW6]
supersedes: null
---

# C-980: O-048 Full Objective Check

Reviewed state: uncommitted working tree on top of 12e63a7 (branch 2.4.0-AgentOrchestrationSpike). This session did not build either Task. Both Tasks are `done` with owner Task-check waivers (T-137 05:34:55Z, T-138 05:38:57Z); this Full Check covers both.

**Result: NEEDS WORK.** The removal work (T-138) is clean. The tile animation (T-137) works on load and reload but freezes after an Objective or Goal switch (I-155). An unplanned change to the task skill sits in the Objective's diff (I-156). Native Windows CI has not run on this code (I-157).

## Full Check Progress Checklist

- [x] Establish Scope
- [x] Freeze The Check Scope
- [x] Turn Acceptance Into Invariants
- [x] Build The Mandatory Coverage Matrix
- [x] Finite External-Boundary Matrix (not applicable: no server, subprocess, or provider in the changed board code; the removed subprocess code is checked as removal)
- [x] Workflow And Side-Effect Check Lock (timer lifecycle only; no persistence or external effects)
- [x] Matrix Completion Lock
- [x] Parallel Planning Advice (not a feature here; nothing to evaluate)
- [x] Perform The Adversarial Pass
- [x] Re-check After Remediation (not applicable: first run)
- [x] Verify File Reality
- [x] Verify Evidence And Gates
- [x] Collect Code Health Evidence
- [x] Complete The Issues Pass
- [x] Summarize Materiality
- [x] Review Code Style

## Scope Lock

1. Criteria: O-048 SC1-SC6 and Boundaries; T-137 DW1-DW6; T-138 DW1-DW6. Guardrails ARCH-02, ARCH-04, TPL-02, CFG-01, CFG-03, DATA-03, STYLE-07, STYLE-09, TEST-01/02/04/06/08.
2. Changed files: the `reviewed.files` list. Entry points: the TUI `Model.Update`/`View` (tile render, spinner tick, load, Objective selection, Goal selection), plain board output, `savepoint` usage and dispatch, `doctor` V2 checks, `agent_launcher` config loading, shipped guidance.
3. Relied-on orchestration: Bubble Tea `tea.Tick` command delivery; the board's load/reload path (`applyLoad`).
4. Matrix axes: tile state (planned, in_progress×{build,test,audit}, done, unknown stage); render surface (board columns, single-column narrow, plain output, detail); width (normal, narrow 11, below label width); timer transitions (open, reload with/without live, tick with/without live, scope change via Objective select/esc/Goal select); removal surfaces (dispatch, usage, imports, doctor, config, guard tests, guidance, `.savepoint/runs`).
5. Admission boundary: supported owner paths in the TUI and CLI; Issues only against named criteria above.

## Coverage Matrix

| Row | Cell | Result | Evidence |
|---|---|---|---|
| Tile in_progress+build | board | Passed | `⠋ Building…`, BUILD badge absent (`columns_view_test`, `TestLiveBuildTileShowsSpinnerAndBuilding`) |
| Tile in_progress+test | board | Passed | `Testing…`, clear accent (`TestLiveTestTileShowsSpinnerAndTesting`) |
| Tile audit / done / planned | board | Passed | byte-identical live vs static render (`TestOtherTilesAreUnchangedWhenLive`) |
| Unknown stage | board | Passed | `liveStageBadge` returns ok=false |
| Live label + other badges | order | Passed | `TestLiveLabelKeepsOtherBadgesInOrder` |
| Narrow tile (11) | wrap | Passed | truncated with `…`, same height (`TestLiveLabelIsCutNotWrappedOnANarrowTile`) |
| Plain output | build/test Task | Passed | `Live` never set by `newTaskCard`; `TestPlainOutputCarriesNoLiveLabel` |
| Idle board | golden | Passed | `internal/board/v2/testdata/golden/` unmodified (`git status` clean) and passing in `make test-full` |
| Help, key map, footer | — | Passed | no change in diff to keys/help/footer files |
| Timer: open/reload with live | start | Passed | `applyLoad` returns `spinnerCmd()`; one timer only |
| Timer: reload with no live | no start | Passed | `TestStageChangeBetweenReloadsStartsAndStopsTheTimer` |
| Timer: tick, no live left | stop | Passed | `TestSpinnerStopsAfterTheLastLiveTileMovesOn` |
| Timer: scope change (Objective select / esc / Goal select) to a view with live tiles after the timer stopped | restart | **Issue I-155** | independent probe: `spinnerRunning=false`, cmd nil, `hasLiveTile()=true` |
| Update/render file work | — | Passed | spinner code reads only `m.Cards`; `tea.Tick` only |
| Removal: dispatch, usage, import | — | Passed | `main.go` diff; `go build` ok |
| Removal: doctor lost runs | — | Passed | call and files gone |
| Removal: `agent_launcher` | — | Passed | `git diff ced16e7^ -- internal/data/launcher_config*.go` empty; repo and scaffold config ship it only as a comment |
| Removal: guard tests | — | Passed | 3 run cases removed, scaffold-copy test kept |
| Removal: guidance | — | Passed | AGENTS.md rows updated; grep of `templates/`, `agent-skills/`, code finds no run-command or `.savepoint/runs` mention outside history |
| Boundaries: workflow skills | — | **Issue I-156** | `agent-skills/savepoint-task/SKILL.md` + template copy changed, unplanned, unrecorded |
| Gates: native Windows CI | — | **Issue I-157** | work uncommitted; no CI run |

## Workflow Lock (spinner timer)

| Order | Operation | State change | Failure / stop | Oracle |
|---|---|---|---|---|
| 1 | load/reload | `spinnerRunning=true` iff live and not running | none | reload test, probe |
| 2 | tick | frame++ and reschedule iff live; else running=false | stops cleanly | tick tests |
| 3 | scope change | `m.Cards` rebuilt | **never restarts a stopped timer** | probe (I-155) |
| 4 | quit | pending tick discarded by Bubble Tea | n/a | — |

## Adversarial Pass

- Bypass of the start rule through another path that rebuilds `m.Cards`: found (`selectObjective`, `applyReleaseSelection`) → I-155.
- Double timers: a reload while a tick is pending does not start a second (`spinnerRunning` guard). Passed.
- Representation: plain output and detail never set `Live`. Passed.
- Width: Braille frames and `…` are single-cell; truncation keeps one line. Passed.
- Test helper changes in `objectives_test.go`/`releases_test.go` stop at a spinner tick; reviewed, they do not mask other commands (the tick is the last link). Passed.

## Gates And Evidence

| Command | Time (UTC) | Toolchain | Result |
|---|---|---|---|
| `git diff --check HEAD` | 2026-10-11T05:40Z | — | exit 0 |
| `make build` | 2026-10-11T05:40Z | go1.26.2 linux/amd64 | exit 0 |
| `make test-full` (fresh) | 2026-10-11T05:40Z | go1.26.2 linux/amd64 | exit 0, includes Windows/darwin cross-builds |
| Probe for I-155 (temporary in-package test, deleted) | 2026-10-11T05:38Z | go1.26.2 | reproduced freeze |
| Native Windows CI | — | — | not run (I-157) |

File reality: every file named in T-137 and T-138 evidence exists or is recorded as deleted (`git status` shows the `D` entries). No phantom files.

## Code Health

`savepoint health check O-048 .` → snapshot `sha256:51e37a46…342533` (created). Verdict: "Code Health does not block clearance." All five instances optional, no blocking finding. This supports the evidence; it does not grant clearance.

## Acceptance Classification

| Criterion | Class |
|---|---|
| O-048 SC1 / T-137 DW1 | Issue (moves on load; freezes after scope change, I-155) |
| O-048 SC2 | Proven |
| O-048 SC3 / T-137 DW2 | Issue (I-155); no-file-work and idle-board parts proven |
| O-048 SC4 / T-137 DW3 | Proven |
| O-048 SC5 / T-138 DW1-DW5 | Proven |
| O-048 SC6 / T-138 DW6 | Unverified: local gates pass, native Windows CI missing (I-157) |
| T-137 DW4 (visual identity) | Proven |
| T-137 DW5 (tests) | Issue: no scope-change timer test (I-155) |
| T-137 DW6 | Proven (gates pass); owner validation still to record |
| O-048 Boundaries | Issue (I-156) |

## Materiality

| Issue | Likelihood | Impact | Materiality | Recommendation |
|---|---|---|---|---|
| I-155 spinner freezes after Objective/Goal switch | High (moving the sidebar cursor triggers it) | Low-Medium (spinner looks stuck; words still right) | Medium | Fix now: start the timer from the scope-change paths; add tests |
| I-156 unplanned skill change | Certain (in the diff) | Low-Medium (changes every executor's workflow without plan) | Low | Owner accepts and planner records it in O-048/Design, or revert |
| I-157 Windows CI missing | Certain | Low (removal only, cross-build ok) | Low | Commit and run CI; name the job |

## Owner Validation Still Needed

T-137 declares `owner_validation.required: true` with `accepted_check: ""`; its User Check (watch the tile animate at build and test, stop at check) is still to be recorded. Best done after I-155 is fixed.

## Observations (non-blocking)

- Design §8 "Task tile run animation" still says "planned, not built" and that O-045's run code "still ship[s]"; the planner should reconcile it once O-048 lands.
- In the narrow single-column layout the timer runs while a live tile is in scope but its column is off screen (recorded in T-137 Limitations). Harmless, small cost.
- T-138 notes the agent command guard blocked running the built binary; usage removal is verified from the diff.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [ ] STYLE-03 **Test branches** — the scope-change branch of the timer is untested (`spinner_test.go`), see I-155.
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [x] STYLE-07 **One source of truth**
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs**
