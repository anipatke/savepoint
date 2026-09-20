---
id: E49-objective-board/T009-stay-current-without-losing-your-place
title: Stay current without losing your place
status: done
objective: Watch the V2 file model, reload through the one load command, and preserve focus and selection across a reload without ever reloading from the archive.
depends_on:
    - E49-objective-board/T008-let-the-gates-decide-which-keys-exist
complexity_tier: medium
complexity_reason: Bounded watch and reload behavior, with focus preservation and an excluded-tree rule to get exactly right.
---

# T009: Stay current without losing your place

## Problem

The board is open while an agent works. Tasks change stage, Checks get written, Issues get opened and resolved, and the router moves — all from another process, while the user is looking at the screen. A board that does not notice is worse than one that never claimed to be live, because the badges keep asserting a state the files no longer have.

The watch set follows the V2 file model rather than the V1 one: `objectives/` recursively, since Tasks live under their Objective's directory, plus `checks/`, `issues/`, `router.md`, and `config.yml`. Three exclusions are deliberate. `releases/` and `defects/` have no V2 meaning. `audit/` has no V2 meaning. And `archive/` must never trigger a reload: it holds byte-preserved historical source from migration, and the live board reading it as a change signal would make old bytes look like new activity. `migrations/` is watched only insofar as pending-migration state matters, which the load command already reads through `migrate.PendingOperation`.

Reload goes through the single load command from T002. This is the property worth protecting: one load path, so a reloaded board and a freshly started board over the same files cannot disagree. Every derived value — the index, the resolved `Next`, per-record decisions, Check chains, Issue links — is rebuilt from that one call rather than patched in place.

Preserving the user's place is the part that makes it usable. A reload must keep the selected Objective, the focused column and card, the scroll offsets, and an open overlay pointed at the same record — when those records still exist. When one does not, the board moves to a defined neighbour and says what happened rather than silently jumping to the top or holding a stale reference. A record that changes column under the user is the common case and needs the focus to follow the record, not the position.

A reload that fails must not destroy a working board. An edit that leaves a project temporarily invalid — a half-written Task file, a dangling reference mid-rename — makes `LoadV2Index` fail closed, and the right response is to keep showing the last good state with a visible diagnostic, then recover on the next successful load. A partially written file is also why the reload is debounced: an editor's save can produce several events in a row, and each one must not trigger its own full index rebuild.

## Context Files

- `internal/board/v2/model.go`
- `internal/board/v2/load.go`
- `internal/board/v2/update.go`
- `internal/board/v2/view.go`
- `internal/board/watch.go`
- `internal/board/io.go`
- `internal/board/debug.go`
- `internal/data/project.go`
- `internal/data/discover.go`
- `internal/migrate/operation.go`
- `internal/testutil/fs.go`
- `agent-skills/bubbletea-tui-design/SKILL.md`

## Acceptance Criteria

- [x] The watcher covers `objectives/` recursively, `checks/`, `issues/`, `router.md`, and `config.yml` under the project's `.savepoint` directory.
- [x] A change under `archive/`, `releases/`, `defects/`, or `audit/` triggers no reload.
- [x] A new Objective directory created while the board is open becomes watched without a restart.
- [x] Every reload goes through the same load command startup uses; no second load path exists, asserted structurally.
- [x] A Task edited on disk updates its card, badges, and detail after reload, and the Next area updates with it.
- [x] A newly written Check updates the affected record's clearance badge and appears in its Check history.
- [x] A newly written or resolved Issue updates the Issues overlay and any linked-record listing.
- [x] A router edit updates the reported selection and the selection diagnostic.
- [x] Rapid successive writes to the same file produce one debounced reload, not one per event.
- [x] After a reload the selected Objective, focused column, focused card, and scroll offsets are preserved when those records still exist.
- [x] A focused Task that moved to a different column keeps focus, following the record rather than the position.
- [x] A focused or open record that no longer exists moves focus to a defined neighbour and reports what happened; no stale record reference is held.
- [x] An open detail or Issues overlay stays open on the same record across a reload when that record survives.
- [x] A reload whose load fails keeps the last good board on screen with a visible diagnostic, and a subsequent successful load clears it.
- [x] A migration operation appearing or completing while the board is open is picked up by the reload and changes the reported Next accordingly.
- [x] Reloading writes nothing, asserted by a byte-and-mtime snapshot across a sequence of external edits and reloads.
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Define the watch set and the excluded trees as data, and add the recursive registration for `objectives/`.
- [x] Add the debounce between a filesystem event and the reload command.
- [x] Route reload through the T002 load command and add the structural assertion that there is only one.
- [x] Add identity-based restoration: record the selected Objective ID, focused Task ID, open overlay record ID, and offsets before reload, and re-resolve them after.
- [x] Add the defined-neighbour fallback and its status message for a record that disappeared.
- [x] Add the failed-reload path that preserves the last good model and surfaces the diagnostic.
- [x] Add tests driving real file edits in a temporary project: each watched path, each excluded path, new Objective directory, debounce, column move, disappearance, overlay survival, failed reload and recovery, migration appearing, and the no-write snapshot.
- [x] Run `go test ./internal/board/...`, then `make build && make test`.

## Context Log

Implemented the V2 live-reload surface. The watcher registers the V2 file
model recursively, ignores historical/V1 trees, learns newly created
Objective directories, and reports one trailing-edge debounced change. The
watcher also observes `.migration` so a pending operation changes the board's
Next projection as it appears or completes.

Reload messages dispatch the existing `loadCmd`; no second loader exists.
Successful loads restore Objective and Task identity, follow a Task across
columns, retain detail/Issues overlays and their offsets, and move deleted
records to a defined neighbour with a status message. Failed external loads
keep the last good board visible with a named diagnostic and recover on the
next valid load. Reloads perform no writes.

Named evidence:

- `TestV2WatchSetIncludesLiveFilesAndExcludesHistoricalTrees`
- `TestV2WatcherReloadsLiveFilesAndLearnsNewObjectiveDirectories`
- `TestV2WatcherIgnoresExcludedTrees`
- `TestV2WatcherDebouncesRapidWrites`
- `TestV2WatcherReportsMigrationAppearingAndCompleting`
- `TestReloadRefreshesCheckAndIssueSurfaces`
- `TestReloadFollowsTaskIdentityAcrossColumnsAndPreservesOverlay`
- `TestFailedReloadKeepsLastGoodBoardAndRecovers`
- `TestReloadRouterEditUpdatesSelectionDiagnostic`
- `TestReloadWritesNothing`
- `TestReloadUsesTheStartupLoadCommand`

Files read: `.savepoint/router.md`, `.savepoint/Guardrails.md`, the E49 epic
detail, this task, `agent-skills/savepoint-build-task/SKILL.md`,
`agent-skills/bubbletea-tui-design/SKILL.md`, the task Context Files, the
existing V2 run/load/model/update/view/detail/Issues tests and fixtures, the
V1 watcher/IO patterns used for comparison, and `.savepoint/visual-identity.md`
for the reload diagnostic's existing terminal treatment.

Files edited:

- `internal/board/v2/watch.go` — V2 watch set, exclusions, recursive directory
  registration, migration watching, and trailing-edge debounce.
- `internal/board/v2/model.go`, `run.go`, `update.go`, `view.go` — watcher
  lifecycle, single-command reload dispatch, identity restoration, failed-load
  preservation, overlay restoration, and visible reload diagnostics.
- `internal/board/v2/watch_test.go` — temporary-project watcher, reload,
  preservation, recovery, migration, and no-write evidence.
- `internal/board/v2/boundary_test.go` — structural assertion for the single
  load command.
- `internal/board/v2/columns_view_test.go` — keeps the existing shrink-focus
  test's fixture structurally valid after record removal.
- this task file — acceptance/plan checklists and handoff evidence.

Quality gates:

- `go test ./internal/board/v2` — passed.
- `go test ./internal/board/...` — passed.
- `make build && make test` — passed for all packages.
- `.savepoint/Health-Check.md` is absent, so Quick health-check evidence does
  not apply.

The task remains `status: in_progress` for user review; only the user may mark
it `done`.
