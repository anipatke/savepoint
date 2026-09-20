---
id: E33-release-docs-view/T002-epic-doc-subview-state
status: done
objective: Add board state, IO messages, and key handling for the Release Docs Epic subview
depends_on:
  - E33-release-docs-view/T001-release-doc-data-loader
complexity_tier: medium
complexity_reason: "Coordinates model state, async file loading, and local key handling across board modules."
---

# T002: Epic Doc Subview State

## Problem

The Epic detail overlay needs explicit state for switching into Release Docs,
selecting a document, and scrolling it without interfering with existing global
board navigation.

## Context Files

- `internal/board/model.go`
- `internal/board/io.go`
- `internal/board/update.go`
- `internal/board/help.go`
- `internal/board/epic_panel_test.go`

## Acceptance Criteria

- [x] Board model stores the selected Epic detail subview, selected release doc,
      loaded docs, and per-doc scroll offset.
- [x] Release docs are loaded through board command/message flow, not direct
      filesystem reads inside `Update()`.
- [x] Users can switch between Epic detail (`1`), Audit (`2`), and Release Docs
      (`3`) using numbered subview keys — matching the existing overlay pattern —
      without changing router state.
- [x] Users can switch between PRD and Design while in Release Docs using keys
      that do not collide with the `1`/`2`/`3` subview keys (`[`/`]`, also
      left/right).
- [x] Up/down scrolling applies to the selected document and does not corrupt
      scroll state for the other document.
- [x] Existing Epic detail and Audit navigation tests continue to pass.

## Implementation Plan

- [x] Add focused Epic detail subview state fields to `Model` or an existing
      embedded state group in `internal/board/model.go`.
- [x] Add release-doc load message and command helpers in `internal/board/io.go`
      using the data loader from T001.
- [x] Trigger release-doc loading when the Epic detail overlay opens or when the
      board data refresh path already reloads release-scoped metadata.
- [x] Extend `handleEpicDetailOverlay` in `internal/board/update.go` with a `3`
      case for Release Docs (alongside the existing `1`/`2` cases), plus document
      selection (`[`/`]`) and document scrolling keys.
- [x] Update help text only if the existing overlay help pattern exposes local
      Epic detail keys.
- [x] Add reducer-level tests for subview switching and scroll preservation.

## Context Log

- Read: `internal/board/model.go`, `update.go`, `io.go`, `watch.go`,
  `view.go`, `epic_panel.go`, `help.go`, and `update_test.go` for the existing
  Epic detail overlay tab/scroll pattern; `internal/data/release_doc.go` (T001)
  for the loader API; T003 task to keep rendering out of scope.
- `internal/board/model.go`: extended `EpicState` with `ReleaseDocs`,
  `ReleaseDocIndex`, and `ReleaseDocOffsets` (per-doc scroll keyed by
  `data.ReleaseDocID`); documented `EpicDetailTab` value `2` = Release Docs.
- `internal/board/watch.go`: added `releaseDocsMsg`. `internal/board/io.go`:
  added `loadReleaseDocsCmd(root)` wrapping `data.LoadReleaseDocs`, surfacing
  loader errors as `errorMsg` (no filesystem reads in `Update`).
- `internal/board/update.go`: handle `releaseDocsMsg` (store + clamp index);
  added `3` subview case (lazy-loads docs when uncached), `[`/`]` (also
  left/right, `h`/`l`) doc selection gated on tab 2, and tab-2-aware
  up/down/pgup/pgdown scrolling via `scrollReleaseDoc`. Helpers
  `selectReleaseDoc`, `scrollReleaseDoc`, `selectedReleaseDoc`,
  `clampReleaseDocIndex`. `openEpicDetailOverlay` resets release-doc state and
  batches `loadReleaseDocsCmd` with the detail read.
- Tests: added six reducer tests in `update_test.go` covering tab switch +
  load, cached no-reload, bracket selection + clamp, tab gating, per-doc scroll
  preservation + top clamp, msg index clamp, and overlay-open reset.
- No `help.go` change: the global help overlay does not enumerate the local
  `1`/`2` tab keys (those live in the overlay's own footer, rendered in T003),
  so there was no existing pattern to extend here.
- Quality gates: `make build && make test` pass.

## Drift Notes

Rendering (the `3` tab indicator, Release Docs body, and `view.go` dispatch)
is intentionally deferred to T003 per its scope, so pressing `3` switches state
and loads docs but still shows the Detail body until T003 wires the renderer.
Upper-bound scroll clamping is also left to the renderer, matching the existing
Detail/Audit behavior where `down` increments the offset unbounded and the view
slices safely. Doc selection additionally accepts `h`/`l` alongside `[`/`]` and
left/right since they were unused in this overlay and match the board's vim
navigation idiom.
