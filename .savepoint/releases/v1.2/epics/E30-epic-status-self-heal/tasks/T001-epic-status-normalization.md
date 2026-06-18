---
id: E30-epic-status-self-heal/T001-epic-status-normalization
status: done
objective: Define a canonical epic-status vocabulary and normalize it on load so the board never renders a blank glyph
depends_on: []
complexity_tier: medium
complexity_reason: "New data-layer vocabulary plus board load wiring and a glyph fallback, mirroring the existing defect-status self-heal pattern"
---

# T001: Epic Status Normalization

## Problem

Epic status is read verbatim into the `epicStatuses` map at
`internal/board/board.go:134` and mapped to a glyph by `statusGlyph`
(`internal/board/status.go:9`). The glyph map only knows
`planned / in_progress / done / audited`; every other value falls through to
`statusGlyphDefault`, a blank space. An agent typo — or the stray `epic-design`
already sitting in `E21-Detail.md` — therefore renders an epic with no visible
status. There is no canonical epic-status vocabulary and no load-time heal,
unlike task `stage` and defect `status`.

## Context Files

- `internal/data/lifecycle.go`
- `internal/data/lifecycle_test.go`
- `internal/board/board.go`
- `internal/board/status.go`
- `internal/board/status_test.go`

## Acceptance Criteria

- [x] A canonical epic-status set `{planned, in_progress, done, audited}` is
      defined once in `internal/data` (constants/helper), with no duplicate list
      elsewhere.
- [x] `NormalizeEpicStatusForLoad` returns a canonical status for any input:
      canonical values pass through; known aliases map to their canonical value;
      anything unknown (including empty and router-state leaks like
      `epic-design`) heals to `planned`.
- [x] The `epicStatuses` map built at load applies `NormalizeEpicStatusForLoad`,
      so the board only ever holds canonical epic statuses.
- [x] `statusGlyph` returns a non-blank, visibly distinct fallback glyph for any
      value that is still not in the canonical set (defense in depth), so no epic
      can render an invisible glyph.
- [x] Tests cover: each canonical pass-through, at least one alias, an unknown
      value healing to `planned`, empty healing to `planned`, and the glyph
      fallback being non-blank.
- [x] `make build && make test` passes.

## Implementation Plan

- [x] In `internal/data/lifecycle.go`, add an `EpicStatus`-style canonical set
      and `CanonicalEpicStatuses()` returning `{planned, in_progress, done,
      audited}`. Reuse existing `data` status string constants where they already
      exist (`ColumnPlanned`, `ColumnInProgress`, `ColumnDone`) rather than
      introducing parallel literals; add an `audited` constant.
- [x] Add `IsCanonicalEpicStatus(value) bool`.
- [x] Add `ResolveEpicStatusAlias(value) (string, bool)` mapping known
      non-canonical leaks onto the canonical set (e.g. task-style
      `complete`/`completed` → `done`; router-ish `epic-design`,
      `epic-task-breakdown`, `task-building` → `planned`). Mirror the shape of
      `ResolveDefectStatusAlias`.
- [x] Add `NormalizeEpicStatusForLoad(value) string`: resolve alias first, else
      pass through canonical, else default `planned`. Mirror
      `NormalizeDefectStatusForLoad`.
- [x] In `internal/board/board.go`, wrap the assignment at line ~134 with
      `data.NormalizeEpicStatusForLoad(status)` so the map holds only canonical
      values. Confirm no other reader bypasses the map.
- [x] In `internal/board/status.go`, change the `default` branch of `statusGlyph`
      to render a visible, dim "unknown" glyph (e.g. `?`) instead of
      `statusGlyphDefault`; keep `statusGlyphDefault` only for the genuine empty
      case if needed.
- [x] Add unit tests in `lifecycle_test.go` (normalization + alias + default) and
      `status_test.go` (fallback glyph non-blank).
- [x] Run `make build && make test`.

## Notes

- Keep router-state vocabulary out of the epic-status canonical set; aliases only
  *heal* leaked values, they do not make router states valid epic statuses.
- This task does not add any TUI control to *change* status (see E31) and does
  not rewrite epic files on disk (healing is load-time; doctor reporting is T002).

## Context Log

**Read:** `internal/data/lifecycle.go`, `internal/data/lifecycle_test.go`,
`internal/data/task.go`, `internal/data/defect.go`, `internal/board/board.go`,
`internal/board/status.go`, `internal/board/epic_panel.go`.

**Edited:**

- `internal/data/lifecycle.go` — added `EpicStatus` type, `EpicStatusAudited`
  constant, `CanonicalEpicStatuses()`, `IsCanonicalEpicStatus()`,
  `ResolveEpicStatusAlias()`, and `NormalizeEpicStatusForLoad()`, mirroring the
  defect-status self-heal functions. `planned/in_progress/done` reuse the
  existing `Column*` constants; `audited` is the only new literal.
- `internal/board/board.go` — the `epicStatuses` map assignment now wraps the
  raw frontmatter value in `data.NormalizeEpicStatusForLoad`, so the map only
  ever holds canonical statuses. `watch.go` reloads through the same
  `loadBoardData` path; `epic_panel.go` and `transitions.go` read the same map,
  so no reader bypasses normalization.
- `internal/board/status.go` — `statusGlyph` default branch now renders a dim
  `?` (`statusGlyphUnknown`) for any non-canonical, non-empty value; the genuine
  empty case still returns `statusGlyphDefault`. The `audited` glyph case now
  uses `data.EpicStatusAudited` instead of a bare literal.
- `internal/data/lifecycle_test.go` — added canonical-set, normalization
  (pass-through, aliases, router leaks, unknown, empty), and alias-resolution
  tests.
- `internal/board/status_test.go` — new file: canonical glyphs non-blank,
  unknown falls back to a visible glyph, empty stays blank.

**Quality gates:** `make build && make test` — all packages pass.
