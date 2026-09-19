---
id: E49-objective-board/T004-navigate-a-project-by-its-objectives
title: Navigate a project by its Objectives
status: done
objective: Replace the release picker and epic panel with an Objective sidebar that filters the columns by recorded ownership and shows each Objective's own status and integration clearance.
depends_on:
    - E49-objective-board/T003-show-the-work-not-the-lifecycle
complexity_tier: medium
complexity_reason: Bounded navigation surface over an existing ownership map, with focus and clamping behavior to get right.
---

# T004: Navigate a project by its Objectives

## Problem

V1's board navigates by release, then by epic — two nested selectors over directory structure. V2 has one level and it is not a directory: an Objective owns a Task because that Task's own `objective` field says so, which is why `index.ObjectiveTasks` exists and why moving a Task file changes neither its identity nor its ownership. The sidebar must read that map and nothing else. Deriving membership from a path, or from an `O###`-shaped directory name, would reintroduce exactly the inference E42 removed.

The sidebar is also where an Objective stops being a folder and starts being a record with its own state. Each row carries the Objective's `title`, its `status`, and its own integration clearance from `ResolveClearance` — because an Objective reaches done only after its Tasks meet completion rules *and* its own independent Check is current. An Objective whose Tasks are all done but whose integration Check is missing is a state the user needs to see from the sidebar, since it is otherwise invisible: every column looks finished. An Objective waiting on an Objective dependency is the other: `ResolveObjectiveDependency` already answers it, and the row should say so rather than the user discovering it by trying to start work.

Selection filters the columns to the selected Objective's Tasks. With no Objective selected the columns show the whole project, which is the honest default for a project whose router names nothing. The initial selection comes from the router's `objective` when it resolves, and from `--objective` when that was passed — and when the router names an Objective that no longer exists, the sidebar shows no selection and the project's own records still render, rather than the board refusing to open over a stale hint. Selection here is view state only: this task writes nothing to `router.md`. The key that records a selection arrives in T008 with the rest of the write path.

An empty project has an empty sidebar, and that is not an error.

## Context Files

- `internal/board/v2/model.go`
- `internal/board/v2/view.go`
- `internal/board/v2/update.go`
- `internal/board/epic_panel.go`
- `internal/board/layout.go`
- `internal/board/util.go`
- `internal/data/objective_v2.go`
- `internal/data/objective_gate_v2.go`
- `internal/data/gate_v2.go`
- `internal/data/project.go`
- `internal/data/router_v2.go`
- `internal/styles/styles.go`
- `.savepoint/visual-identity.md`

## Acceptance Criteria

- [x] A sidebar lists every Objective in `index.Objectives` in stable `O###` order, each row showing its ID, its `title`, and its `status`.
- [x] Each row shows the Objective's own clearance from `ResolveClearance`, distinguishing missing, needs_work, stale, unknown, and current, reusing the badge mapping from T003 rather than a second one.
- [x] An Objective whose every owned Task is done but whose integration clearance is not current is visibly distinguished from an Objective that is genuinely complete.
- [x] An Objective with an unsatisfied Objective dependency shows that wait, read from `ResolveObjectiveDependency`, naming the Objective it waits on.
- [x] Selecting an Objective filters all three columns to the Tasks in `index.ObjectiveTasks` for that ID; membership is never derived from a file path or directory name.
- [x] With no Objective selected, all three columns show every Task in the project.
- [x] The initial selection comes from `--objective` when passed, otherwise from the router's `objective` when it resolves against the index, otherwise nothing.
- [x] A router naming an Objective absent from the index leaves the sidebar with no selection, renders all Tasks, and surfaces the selection diagnostic; the board still opens.
- [x] Focus moves between sidebar and columns predictably, and focus changes colour and glyphs only — sidebar width and column geometry are identical in both focus states.
- [x] Sidebar navigation clamps at both ends, survives a list shorter than the cursor, and is idempotent on repeated keypresses.
- [x] A sidebar longer than the viewport scrolls without wrapping.
- [x] An empty project renders an empty sidebar with no error and no diagnostic.
- [x] No `router.md` write occurs on any navigation or selection in this task, asserted by a byte-and-mtime snapshot around a navigation sequence.
- [x] No release picker, epic panel, epic detail, or epic status appears anywhere in the V2 package.
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Add sidebar state to the V2 model: the ordered Objective IDs, the cursor, the selection, and the panel focus flag.
- [x] Build the ordered list and per-Objective resolved decisions in the load command, not in the view (ARCH-02).
- [x] Add sidebar rendering reusing the T003 badge mapping for clearance and dependency state.
- [x] Add the all-Tasks-done-but-not-cleared distinction and confirm it reads differently from a done Objective.
- [x] Add ownership filtering to the column grouping, keyed off `index.ObjectiveTasks`.
- [x] Add initial-selection resolution in precedence order, including the unresolvable-router case.
- [x] Add key handling for sidebar focus and movement, with clamping.
- [x] Add tests: ordering, filtering by ownership, the unresolved router hint, empty project, clamping, focus geometry, and the no-write snapshot over a navigation sequence.
- [x] Run `go test ./internal/board/...`, then `make build && make test`.

## Context Log

Built the Objective sidebar, ownership filtering, and the focus model over the V2 board T003 left.

**Read:** `.savepoint/router.md`, `.savepoint/Guardrails.md` (STYLE rules; no `.savepoint/Health-Check.md` exists, so the Quick check step does not apply), `E49-Detail.md`, this task, `.savepoint/visual-identity.md`, and the Context Files: `internal/board/v2/{model,view,update,load,card,column,badges,run,plain}.go`, `internal/board/epic_panel.go`, `internal/data/{objective_v2,objective_gate_v2,gate_v2,project,router_v2,next}.go`, `internal/styles/styles.go`.

**Edited:**

- `internal/board/v2/objectives.go` (new) — `ObjectiveRow`, `objectiveRows`, `ownedTasksComplete`, `unsatisfiedObjectiveWaits`, `taskIDsInView` (the one place membership is decided), and `renderSidebar`/`renderObjectiveRow`.
- `internal/board/v2/badges.go` — `objectiveIntegrationBadge` and `objectiveWaitBadge`, added to the one mapping rather than a second copy.
- `internal/board/v2/model.go` — `Objectives`, `ObjectiveCursor`, `SidebarFocused`, `sidebarVisible`, `objectiveIndex`, `cardCount`.
- `internal/board/v2/update.go` — per-surface key dispatch, `tab` focus, cursor clamping, `enter`/`esc` selection, `clampObjectiveCursor`, and focus handed back when a resize drops the sidebar.
- `internal/board/v2/card.go` — `groupTaskCardsFor`, with `groupTaskCards` kept as the unfiltered case.
- `internal/board/v2/column.go` — `columnCursor` replaces the two positional focus arguments, so a column's scroll window follows the cursor while only the accent follows surface focus.
- `internal/board/v2/view.go` — the sidebar in place of T003's reserved space, `selectionDiagnosticSummary`, the "n of m tasks" filter statement, and focus-aware hints.
- `internal/board/v2/{plain,run,load}.go` — non-TTY counts now honour the same resolved selection the TUI does; the now-unused `tasksInColumn` removed.
- `internal/styles/styles.go` — `SidebarSelected`, one token in the existing palette.
- `AGENTS.md` — Codebase Map row for `internal/board/v2/` (ARCH-04).

**Evidence** (`internal/board/v2/objectives_test.go` unless noted, over temporary fixture projects only — TEST-04):

- `TestSidebarListsEveryObjectiveInOrder`, `TestSidebarShowsEachObjectivesOwnClearance` (all five states), `TestSidebarSeparatesFinishedTasksFromAFinishedObjective`, `TestSidebarNamesTheObjectiveAWaitIsOn`.
- `TestSelectionFiltersColumnsByRecordedOwnership` — the fixture files T004 under O001's directory while its record names O003, so the assertion fails if membership ever comes from a path.
- `TestNoSelectionShowsEveryTaskInTheProject`, `TestInitialSelectionPrecedence`, `TestRouterNamingAMissingObjectiveOpensTheBoardAnyway`.
- `TestSidebarNavigationClampsAndIsIdempotent`, `TestSidebarSurvivesAReloadThatShortensTheList`, `TestSidebarScrollsRatherThanWrapping`, `TestEmptyProjectRendersAnEmptySidebar`, `TestNarrowTerminalOffersNoSidebarFocus`, `TestFocusChangesAccentsNotGeometry` (line-by-line widths).
- `TestNavigationWritesNothingToTheProject` — content-and-mtime snapshot of every file under `.savepoint/` around a sixteen-key navigate/select/clear sequence.
- `boundary_test.go: TestPackageCarriesNoReleaseOrEpicSurface` — no `Release`, `Epic`, `Defect`, or `AuditRegister` identifier in the package.
- `run_test.go: TestRunWithoutTTYCountsOnlyTheSelectedObjectivesTasks`.

**Quality gates:** `go test ./internal/board/...` ok; `make build && make test` ok (all packages); `go vet ./...` clean; `gofmt -l` reports nothing in the files touched.

**Decisions worth the next session's attention:**

- Rows and their resolved decisions are built in `applyLoad`, the same seam T003 builds cards in, rather than inside `loadCmd`. Both are outside rendering, so ARCH-02 holds and there is one place where resolvers run per load.
- `tab` is the only key that crosses surfaces. Letting the columns' `left` hand focus to the sidebar would have made the leftmost column's clamp non-idempotent, which the ACs require.
- The sidebar draws only at `sidebarBreakpoint` (120 columns) and wider, the geometry T003 set; below it the sidebar is absent and `tab` is not offered. Narrow-width behaviour belongs to T010.
- A reload re-applies the initial-selection precedence, so a selection the user made by hand does not survive one. Preserving place across a reload is T009's subject.
- `--objective` and the router selection now filter the non-TTY counts too, so the flag means the same thing on both surfaces rather than being silently ignored when piped.
