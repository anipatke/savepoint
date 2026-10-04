---
id: T-109
title: "Offer a parallel planning toggle in Advanced Options"
objective: O-037
status: done
depends_on: [{task: T-108, requires: clear}]
owner_validation:
  required: true
  accepted_check: ""
planned_by: {role: planner, session: codex-o033-advisory-replan-2026-10-03}
complexity_tier: medium
complexity_reason: "A settings overlay must save explicitly and retain reload/focus behavior."
planned_reads:
  - "internal/data/feature_preferences.go"
  - "internal/data/config.go"
  - "internal/board/v2/model.go"
  - "internal/board/v2/update.go"
  - "internal/board/v2/view.go"
  - "internal/board/v2/help.go"
  - "internal/board/v2/footer_test.go"
  - "internal/board/v2/io.go"
  - "internal/board/v2/load.go"
  - "internal/board/v2/load_test.go"
  - "internal/board/v2/watch.go"
  - ".savepoint/visual-identity.md"
planned_writes:
  - "internal/board/v2/options.go"
  - "internal/board/v2/options_test.go"
  - "internal/board/v2/model.go"
  - "internal/board/v2/update.go"
  - "internal/board/v2/view.go"
  - "internal/board/v2/help.go"
  - "internal/board/v2/footer_test.go"
  - "internal/board/v2/io.go"
  - "internal/board/v2/load.go"
  - "internal/board/v2/load_test.go"
  - "internal/board/v2/watch.go"
check_waiver:
  task: T-109
  reason: Owner completed this Task via the board without requesting a Task Check.
  actor:
    role: owner
    session: board-owner
  recorded_at: "2026-10-03T22:00:54Z"
---

# Offer a parallel planning toggle in Advanced Options

## Outcome

Owners can switch optional parallel advice on/off from Advanced Options with saved-state feedback and no effect on core Code Health.

## User Check

Open Advanced Options, toggle Parallel planning, restart the board, then toggle it off; confirm saved plans remain and Code Health still appears as before.

## Done When

- A documented Settings action opens a keyboard-accessible Advanced Options screen containing exactly the real Parallel planning option and a short explanation that suggestions can be ignored.
- Board load/reload reads the saved preference; toggles persist through explicit Bubble Tea commands. Show saved state only after success and actionable messages on stale/conflicting or failed writes.
- Settings survive board restart and refresh after external edits to config.yml. Esc/close returns focus to the original surface; overlays and terminal sizing follow existing board patterns.
- No Code Health toggle or behavior changes, provider execution, destructive cleanup, shell calls or rendering IO. Turning this option off preserves all lane metadata and evidence.
- Record per-criterion evidence and a fresh full gate; optional Task Check waiver is recorded only on explicit owner instruction. Mandatory independent Full Objective Check remains required.

## Context Files

- `internal/data/feature_preferences.go`
- `internal/data/config.go`
- `internal/board/v2/model.go`
- `internal/board/v2/update.go`
- `internal/board/v2/view.go`
- `internal/board/v2/help.go`
- `internal/board/v2/footer_test.go`
- `internal/board/v2/io.go`
- `internal/board/v2/load.go`
- `internal/board/v2/load_test.go`
- `internal/board/v2/watch.go`
- `.savepoint/visual-identity.md`

internal/data/feature_preferences.go is intentionally created by the preference prerequisite. Other existing Context Files must exist; a material interface mismatch returns REPLAN REQUIRED.

## Design References

O-037 Confirmed Design — 2026-10-03: preference storage and Advanced Options delivery before lane implementation. Project Design sections 1, 8 and 9 describe existing ownership, rendering and concurrency boundaries.

## Guardrails

FS-01, FS-04, FS-05, FS-06, DATA-01, DATA-02, DATA-03, ARCH-02, ARCH-03, CFG-01, CFG-02, CFG-03, TEST-01, TEST-02, TEST-03, TEST-04, TEST-08, TEST-09, STYLE-01..10.

## Implementation Plan

1. Confirm the prerequisite preference interface; no lane projection dependency; add new options.go view/model support and options_test.go with existing overlay conventions.
2. Choose one unused discoverable Settings key and document it in help/footer. Dispatch toggle saves as explicit commands followed by reload.
3. Watch preference-file changes through the existing refresh boundary and preserve board selection/focus while reloading.
4. Verify default/off/on, save/restart, external edit, stale or unwritable save, keyboard close/focus and unchanged Code Health surfaces.

## Boundaries

Implement the scoped preference/screen only. Anticipated manifests and lane choice remain advice, not execution allowlists or gates. Retain ordinary extra-read logging and return REPLAN REQUIRED only for an actual material plan/acceptance gap. Code Health, completion authority, existing dependency gates and actual-worktree rules stay unchanged. No agent settings toggling for the owner's real project during execution; use temporary fixtures. No Check or owner waiver is invented.

## Technical Verification

Fresh make test-full for configuration replacement/path and reload-watch behavior. Focused tests are iteration only. Native windows-tests CI evidence is produced by repository CI and supplied by the owner before the independent Full Objective Check under agent-skills/references/check-method.md. Owner validation is required for the settings screen.

## Technical Evidence

Executor session: claude-sonnet-5-5, 2026-10-03. No self-clearance; no Check or owner waiver recorded.

**Files changed.** New: `internal/board/v2/options.go`, `internal/board/v2/options_test.go`. Edited: `model.go` (Options field), `update.go` (key dispatch, `optionsSavedMsg` case, overlay-origin restore, `quit()` helper extracted from the q/ctrl+c branch and reused by the options screen), `view.go` (overlay render, hints), `help.go` (help row), `load.go` (`ProjectState.Features` read once per load). `io.go`, `watch.go`, `footer_test.go`, `load_test.go` needed no change: `config.yml` was already in `v2WatchFiles`, and the save command lives in `options.go`.

**Logged extra reads** (outside Context Files, all read-only, for conventions or fixtures): T-108 task file head and the O-037 Objective diff (dependency/start check); `internal/board/v2/releases.go`, `actions.go`, `health.go`, `width.go` (overlay renderer, used keys, `fitLine`); `internal/styles/styles.go` (style names); test helpers in `fixture_test.go`, `view_test.go`, `columns_view_test.go`, `objectives_test.go`, `releases_test.go`, `boundary_test.go`, `objective_close_test.go`, `watch_test.go` (helpers and one failing footer assertion); `templates/project-v2/.savepoint` listing.

**Per-criterion results.**
- Settings action opens a keyboard Advanced Options screen with exactly the one real option and an explanation that suggestions can be ignored: `o` opens it (documented in Help and the footer). `TestOptionsOpensWithParallelPlanningOffAndExplainsIt` (one `Parallel planning:` row, text present), `TestOptionsKeyIsDocumentedInFooterAndHelp`. PASS.
- Load/reload reads the saved preference; toggles persist through explicit commands; saved state shown only after success; actionable stale/failed messages: `loadFeatureState` runs only in `loadProject`; `writeParallelPlanningCmd` wraps `data.WriteParallelPlanning`; "Saved." appears only when the read-back matches. `TestOptionsToggleSavesAndSurvivesRestart`, `TestOptionsStaleSaveIsRefusedAndExplained` (refusal, file untouched, retry succeeds and keeps the external edit), `TestOptionsUnwritableConfigReportsAndSavesNothing` (skipped on Windows or root), `TestOptionsUnreadableConfigRefusesToggleWithoutWriting`, `TestOptionsMalformedPreferenceIsReportedNotHealed`, `TestOptionsSecondToggleWhileSavingIsIgnored`. PASS.
- Survives restart, refreshes after external config.yml edits; Esc/close returns focus; overlay and sizing follow board patterns: `TestOptionsRefreshesAfterAnExternalConfigEdit`, `TestOptionsCloseReturnsFocusToTheOriginSurface` (esc, q, o; sidebar and column origins), `TestOptionsFitsTheTerminal` (100x30, 60x20, 40x14), `TestOptionsCtrlCStillQuits`. PASS.
- No Code Health change, provider execution, destructive cleanup, shell calls or rendering IO; turning the option off preserves other records: `TestOptionsToggleLeavesSavedPlansAndCodeHealthAlone` (every non-config file byte-identical, header chip and closed board render identical). The only file written is config.yml, through the T-108 writer. PASS.
- Evidence and fresh full gate: below. Optional Task Check waiver not recorded (owner has not instructed one).

**Commands.** `make build && make test-fast`: pass. `make test-full`: exit 0, started 2026-10-03T21:59:45Z, finished 2026-10-03T22:00:02Z, go1.26.2 linux/amd64. A first fast run failed `TestSpaceOnDoneObjectiveDoesNothingAndKeysRestoreHints` because my `o:options` hint pushed `?:help` out of the truncated 120-column footer; fixed by placing `o:options` last, then both gates were rerun.

**Limitations.** Native Windows CI evidence is not produced here (owner supplies it per Technical Verification). Owner validation of the screen is still required; I ran no interactive TUI session and used only temporary fixtures, never the real project's settings. The footer hint `o:options` is the first thing truncated on narrow terminals; Help always lists it. The unwritable-config test is skipped when run as root or on Windows.

## Drift Notes

Record material implementation/acceptance deltas for planner reconciliation. Ignoring a suggested lane or differing from an anticipated manifest alone requires no replan.


Prerequisite-created setting interface: `internal/data/feature_preferences.go` is intentionally created by T-108; its absence before that dependency completes is expected.
