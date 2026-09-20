---
id: E49-objective-board/T002-open-the-board-a-v2-project-already-has
title: Open the board a V2 project already has
status: done
objective: Dispatch the board on schema version and stand up internal/board/v2 far enough that a V2 project opens, an empty project is normal, and an invalid one is named rather than hidden.
depends_on:
    - E49-objective-board/T001-give-the-board-a-router-writer-it-cannot-misuse
complexity_tier: high
complexity_reason: New package, new entry-point branch, and the load and failure contracts every later task builds on.
---

# T002: Open the board a V2 project already has

## Problem

`savepoint init` has written V2 projects since E47. `internal/board` still walks `releases/epics/tasks`, so the very first thing a new user does after `init` fails: `data.Discover.ListReleases` returns `releases directory not found` and the board never starts. This task is the one that makes `savepoint board` work at all on the default project shape, and it settles the three contracts every later task in the epic inherits.

**The dispatch is the same one every other command already makes.** `board.Run` resolves the root, calls `data.LoadProject`, and branches on `Project.SchemaVersion` — exactly as `resume`, `doctor`, and `upgrade-assets` do. Doing it at the entry point means the V2 board never sees a V1 project and needs no fallback anywhere inside it, and it means this epic edits exactly one V1 board file. `internal/board/v2` is a new package rather than more files in `internal/board` because that is the only way Go makes a V1 type unreachable from a V2 render path; the package is deliberately shaped so E50 can delete the V1 files and promote it.

**An empty project is the normal case, not an edge case.** A project fresh from `templates/project-v2` has no Objectives, no Tasks, no Checks, and no Issues. It must open — empty columns, no Objective selected, and the planning next action — with no error, no diagnostic, and nothing that reads as broken. This is the state a user is in the first time they ever run the command.

**A project that will not load is a screen, not a crash, and never a fallback.** `LoadV2Index` fails closed on duplicate identity, unsafe paths, missing ownership, dangling references, and a Task missing its required `title`. That last one is why the epic's title requirement is structural: `DecodeTaskV2` rejects a titleless Task, so no such Task can ever reach a card. The board's half of that guarantee is to render the diagnostic as readable text naming the file and the problem, to draw no partial board, to exit nonzero without a TTY, and under no circumstances to retry the load through V1 discovery. A fallback here would convert a named data defect into a board quietly showing something else.

The load itself belongs in a `tea.Cmd` (ARCH-02): `LoadProject`, `ReadStateV2`, `migrate.PendingOperation`, and `ResolveNext`, returning one typed message. Startup and every later reload use that same command, so there is one load path rather than a startup path and a refresh path that can drift.

`cmd/board.go` gains `--objective` and keeps parsing only (ARCH-01). It does not know the schema, so the rejection of a filter that does not apply — `--epic` against a V2 project, `--objective` against a V1 one — happens after dispatch and names both the flag and the schema it was refused for (CFG-01).

## Context Files

- `internal/board/board.go`
- `internal/board/board_test.go`
- `internal/board/io.go`
- `internal/board/model.go`
- `internal/board/tui.go`
- `internal/board/view.go`
- `cmd/board.go`
- `cmd/board_test.go`
- `main.go`
- `internal/data/project.go`
- `internal/data/router_v2.go`
- `internal/data/next.go`
- `internal/data/errors.go`
- `internal/migrate/operation.go`
- `internal/testutil/fixture.go`
- `internal/testutil/fs.go`
- `templates/project-v2/.savepoint/config.yml`
- `templates/project-v2/.savepoint/router.md`
- `AGENTS.md`
- `agent-skills/bubbletea-tui-design/SKILL.md`

## Acceptance Criteria

- [x] `board.Run` resolves the root once, calls `data.LoadProject`, and dispatches on `Project.SchemaVersion`; a V1 project reaches today's board with unchanged behavior and every existing `internal/board` test still passes.
- [x] `internal/board/board.go` is the only V1 board file this task changes.
- [x] A new `internal/board/v2` package exists with its model state, its load command, and its top-level view; no file in it references a V1 board type, `data.Task`, `data.Defect`, `data.ReleaseInfo`, `data.EpicInfo`, or `data.AuditRegisterSet`, asserted by a package-level import and reference test.
- [x] The load command performs `data.LoadProject`, `ReadStateV2`, `migrate.PendingOperation`, and `data.ResolveNext`, and returns one typed message; no filesystem access happens in `Update` (ARCH-02).
- [x] A project created from `templates/project-v2` opens with three empty columns, no Objective selected, and no error or diagnostic anywhere on screen.
- [x] A V2 project with Objectives and Tasks opens and reports the counts it loaded, proving the index reached the model.
- [x] A project `LoadV2Index` refuses — duplicate ID, missing owner, dangling reference, and a Task with no `title` each exercised — renders a diagnostic screen naming the offending file path and the problem, draws no columns, and never calls V1 discovery.
- [x] The same failures exit nonzero without a TTY and print the diagnostic to stderr.
- [x] A malformed or unreadable `router.md`, and an unknown router `state` value, each render as a named diagnostic rather than a healed default (DATA-03).
- [x] A project with a pending migration operation opens and reports it, taking the report from `migrate.PendingOperation`.
- [x] `cmd/board.go` parses `--objective` and rejects unknown flags and positional arguments as it already does; the `cmd` package's thin-imports test still passes (ARCH-01).
- [x] `--epic` or `--release` against a V2 project, and `--objective` against a V1 project, each fail with a message naming the flag and the schema version, and no board starts (CFG-01).
- [x] `--objective` naming an Objective that does not exist fails with a named error listing nothing it guessed at.
- [x] `AGENTS.md` Codebase Map gains a row for `internal/board/v2/` and the `internal/board/` and `cmd/` rows describe the dispatch.
- [x] Every test builds its project in a temporary directory; the live `.savepoint/` project is never read or written by a test (TEST-04).
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Read `board.Run`/`RunWithFilters`/`newProjectModel` and decide the exact seam so the V1 path keeps its current call shape.
- [x] Add the schema branch to `board.Run`, passing the already-loaded `*data.Project` down rather than loading twice.
- [x] Create `internal/board/v2` with `model.go` holding the index, router, migration state, resolved `Next`, and view config.
- [x] Add `load.go` with the single load command and its typed success and failure messages.
- [x] Add `view.go` with the empty-project view and the load-diagnostic view, and a minimal status bar.
- [x] Add the package boundary test asserting no V1 board or V1 data type is referenced.
- [x] Extend `cmd/board.go` with `--objective` and its tests.
- [x] Add the filter/schema mismatch diagnostics after dispatch, where the schema is known.
- [x] Add V2 fixture helpers building valid and each invalid project shape in a temporary directory.
- [x] Update the `AGENTS.md` Codebase Map rows.
- [x] Run `go test ./internal/board/... ./cmd/... ./internal/data/...`, then `make build && make test`.

## Context Log

### Files read

`.savepoint/router.md`, `E49-Detail.md`, this task, `.savepoint/Guardrails.md`, `T001` frontmatter (dependency status), `T003`/`T010` headers (to keep the card and plain-output surfaces they own out of this task), and every file in `## Context Files` except `internal/board/model.go`/`view.go`/`io.go`, which were read for the model, layout, and command patterns the V2 package mirrors. Also read for a targeted verification: `internal/board/interfaces.go`, `internal/board/plain.go`, `internal/board/theme.go`, `internal/styles/styles.go`, `internal/data/discover.go` (V2 discovery path rules), `internal/data/config.go` (`ReadSchemaVersion`, `SchemaVersion`), `internal/data/objective_v2.go`/`task_v2.go` (record shapes), `internal/migrate/command.go`, `cmd/resume.go`/`resume_test.go` (the thin-imports assertion), and `main_resume_test.go`/`main_test.go` (the V2 fixture and subprocess harness patterns).

### Files added

- `internal/board/v2/run.go` — `Options` and `Run`: the Bubble Tea program on a terminal, deterministic text otherwise, both over one load.
- `internal/board/v2/model.go` — model state, `Init`, and the reducer; `--objective` resolution by exact global ID and the Objective-in-view rule.
- `internal/board/v2/load.go` — `ProjectState`, the one `projectLoadedMsg`, `loadCmd`, `loadProject`, and the record counts every surface reads.
- `internal/board/v2/view.go` — the board view (header counts, selection line, Next line, three columns, status bar) and the load-diagnostic screen.
- `internal/board/v2/plain.go` — non-TTY rendering, built with `fmt` alone so it carries no escape sequence.
- `internal/board/v2/{fixture,load,view,run,boundary}_test.go`, `internal/board/dispatch_test.go`, `main_board_test.go`.

### Files edited

`internal/board/board.go` (the dispatch seam — the only V1 board file changed), `cmd/board.go`, `cmd/board_test.go`, `main.go`, `AGENTS.md`.

### Decisions worth recording

- **The dispatch's `*data.Project` is the schema signal, not the model's state.** The plan said to pass the loaded project down rather than load twice, but the acceptance criteria also require the load to happen inside the `tea.Cmd` that startup and every reload share. Those cannot both hold for the V2 board, so the dispatch's `LoadProject` decides which board runs and the V2 board loads its own state through `loadCmd`. The V1 path does receive the resolved root (`newProjectModelAtRoot`) so nothing is discovered twice on it.
- **`board.Run` routes a failed load by schema version.** A V2 project the index refuses never reaches `Project.SchemaVersion`, so the error path consults `data.ReadSchemaVersion` alone to decide whether the V2 board should render the diagnostic. A load failure on anything else — including an unreadable or malformed `schema_version` — is returned as it is.
- **`RunTUI` keeps its current call shape.** Threading the resolved root into it would edit a second V1 board file for no behavior change, which AC 2 forbids; it resolves the same root this dispatch just resolved, and E50 deletes it.
- **A pending migration is answered without loading the records**, exactly as `runResume` does: a project mid-conversion holds records whose meaning is not yet settled. The board therefore reports the migration and `migrate`'s own recovery guidance with zero counts and empty columns rather than interpreting half-converted records, and `ResolveNext` returns `NextPendingMigration` from the injected state. `--objective` is not judged on that path, because there is no index to resolve it against.
- **An unknown `--objective` is a refusal, not a screen.** The plain path returns the error directly; the terminal path records it on the model and quits, and `Run` returns it, so both exit nonzero without a board.
- **Colour-profile and height budgeting are not set up here.** The V1 board's `applyColorProfile` is a V1 file, and forced-colour, monochrome, and narrow-width correctness are T010's; this task lays out against width only.

### Evidence

- `go test ./internal/board/... ./cmd/... ./internal/data/... .` — pass. Named cases:
  - `internal/board/v2/load_test.go`: `TestLoadProjectLoadsIndexRouterAndNext`, `TestLoadProjectEmptyProjectFromTemplateIsNormal`, `TestLoadProjectRefusedIndexReportsNamedDiagnostic` (four subtests: titleless Task, duplicate Objective ID, Task owned by a missing Objective, Task depending on a missing record), `TestLoadProjectRouterDiagnostics` (unreadable, malformed body, unknown state, unknown key), `TestLoadProjectPendingMigrationIsReportedFromMigrate`, `TestLoadCmdReturnsTheSameSingleMessage`.
  - `internal/board/v2/view_test.go`: `TestViewEmptyTemplateProjectOpensWithThreeEmptyColumns`, `TestViewReportsTheCountsItLoaded`, `TestViewDiagnosticScreenDrawsNoColumns`, `TestViewReportsPendingMigration`, `TestViewBeforeFirstLoadDrawsNoBoard`, `TestUpdateWindowSizeAndQuit`.
  - `internal/board/v2/run_test.go`: `TestRunWithoutTTYLeadsWithNextAndReportsCounts`, `TestRunWithoutTTYIsDeterministic`, `TestRunWithoutTTYRefusedIndexFailsAndWritesNoBoard` (the same four refusals), `TestRunObjectiveFilterSelectsAndRejects`, `TestUpdateUnknownObjectiveFilterQuitsWithError`, `TestUpdateObjectiveFilterSelectsThatObjective`, `TestUpdateRouterSelectionThatNoLongerResolvesSelectsNothing`.
  - `internal/board/v2/boundary_test.go`: `TestPackageDoesNotImportTheV1Board`, `TestPackageReferencesNoV1RecordType`, `TestUpdateDoesNotPerformIO`.
  - `internal/board/dispatch_test.go`: `TestRunWithFiltersDispatchesV1ProjectToTodaysBoard`, `TestRunWithFiltersDispatchesV2ProjectToTheV2Board`, `TestRunWithFiltersV2LoadFailureReportsTheDiagnosticWithoutV1Fallback`, `TestRunWithFiltersRejectsFiltersTheSchemaHasNoMeaningFor` (three subtests), `TestRunWithFiltersRejectsUnknownObjectiveOnAV2Project`, `TestRunWithFiltersReportsAMissingProject`, `TestRunWithFiltersPassesV1FiltersThrough`.
  - `main_board_test.go`: `TestMainBoardV2ProjectWithoutTTYReportsNextAndExitsZero`, `TestMainBoardV2LoadFailureExitsNonzeroWithDiagnosticOnStderr`, `TestMainBoardRejectsV1FilterOnAV2Project` — end to end through the built command in a temporary project directory, proving the nonzero exit and the diagnostic on stderr.
  - `cmd/board_test.go`: `TestRunBoardObjective`, `TestRunBoardObjectiveAlongsideV1Filters`, `TestRunBoardObjectiveMissingValue`, plus the unchanged `TestPackage_staysThinNoDomainImports`.
  - Existing `internal/board` tests unchanged and passing, including `TestNewProjectModelLoadsReleasesEpicsAndTasks` and `TestNewProjectModelFreshInit`.
- `make build && make test` — pass (`go build ./...`, `go vet ./...`, and the full suite; no package reports a failure).
- Rendered views inspected with `NO_COLOR=1` for the populated, empty-template, and refused-load cases: three intact columns, no wrapping at 100 columns, and no columns at all on the diagnostic screen.

`.savepoint/Health-Check.md` is absent, so no Quick health-check evidence block applies.
