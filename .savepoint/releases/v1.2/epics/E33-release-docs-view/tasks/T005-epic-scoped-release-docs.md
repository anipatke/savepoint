---
id: E33-release-docs-view/T005-epic-scoped-release-docs
status: done
objective: Make the Epic Docs [3] subview load epic-scoped PRD/Design, and add a top-level shortcut for the project-wide PRD/Design
depends_on:
    - E33-release-docs-view/T003-release-docs-renderer
    - E33-release-docs-view/T004-release-docs-verification
complexity_tier: medium
complexity_reason: "Splits one document source into two scopes: a new epic-scoped loader path plus a new top-level overlay, keybinding, and state, reusing the existing renderer."
---

# T005: Epic-Scoped Release Docs + Project Docs Shortcut

## Problem

The Epic detail overlay's Docs [3] subview is reached from a specific epic and
is titled "RELEASE DOCS", so it reads as "this epic's documents". But the loader
(`loadReleaseDocsCmd(m.Root)` → `data.LoadReleaseDocs(root)`) always reads the
project-wide `.savepoint/PRD.md` and `.savepoint/Design.md` regardless of which
epic is selected. Every epic shows the same project docs.

Two corrections are needed:

1. Docs [3] must be **epic-scoped** — it should show the selected epic's own
   PRD/Design, not the project-level documents.
2. The project-level PRD/Design are still valuable, so they need a **separate
   entry point** from the top-level board (a global shortcut), independent of
   any epic.

## Context Files

- `internal/data/release_doc.go`
- `internal/data/release_doc_test.go`
- `internal/board/io.go`
- `internal/board/update.go`
- `internal/board/model.go`
- `internal/board/view.go`
- `internal/board/epic_panel.go`
- `internal/board/help.go`
- `internal/board/epic_panel_test.go`
- `internal/board/update_test.go`
- `internal/board/io_test.go`

## Decisions to Confirm

These shape the implementation; defaults below are the recommendation unless the
reviewer says otherwise.

- **Epic doc file naming** — read per-epic docs from the epic directory using the
  `E##-` short-ID prefix convention that matches its siblings (`E##-Detail.md`,
  `E##-Audit.md`): i.e. `E33-PRD.md` and `E33-Design.md`. (Alternative: bare
  `PRD.md`/`Design.md` inside the epic dir, as some archived v1 epics used.)
  Default: **prefixed** (`E##-PRD.md` / `E##-Design.md`).
- **No project fallback in Docs [3]** — when an epic has no PRD/Design, show the
  existing missing/empty state scoped to the epic; do **not** silently fall back
  to the project docs. Keeping the scopes separate is the whole point. Default:
  **no fallback**.
- **Top-level shortcut key** — recommend `P` (project docs) since `p` is taken
  (router priority) and `d` is defects. Final key chosen to fit the footer/help
  layout. Default: **`P`**.
- **Subview retitle** — rename the epic subview heading from "RELEASE DOCS" to
  "EPIC DOCS" so the scope is unambiguous, and reserve "RELEASE DOCS" / "PROJECT
  DOCS" wording for the new top-level overlay. Default: **retitle**.

## Acceptance Criteria

### Epic-scoped Docs [3]

- [x] Docs [3] loads the selected epic's own PRD/Design from the epic directory,
      not `.savepoint/PRD.md` / `.savepoint/Design.md`.
- [x] Switching the selected epic and reopening the detail overlay shows that
      epic's documents (no cross-epic bleed).
- [x] An epic with no PRD/Design renders the existing read-only missing/empty
      state, scoped to the epic (e.g. names the epic, not the project), and never
      falls back to project-level docs.
- [x] A read error (non-`IsNotExist`) surfaces with path context as a status
      message, consistent with the existing loader contract.

### Project-level docs shortcut

- [x] A global keybinding on the top-level board opens a project-level docs
      overlay reading `.savepoint/PRD.md` / `.savepoint/Design.md`.
- [x] The overlay reuses the existing Release Docs renderer (PRD/Design selector,
      scrollable body, per-doc scroll offset, missing/empty states) without the
      epic Detail/Audit/Docs tab strip.
- [x] The shortcut is discoverable: it appears in the footer hint and the help
      overlay.
- [x] `esc`/`q` closes the overlay and returns to the board.

### General

- [x] No new panel/card nesting style inconsistent with the board.
- [x] `make build && make test` pass; `go vet ./internal/board/ ./internal/data/`
      clean.

## Implementation Plan

- [x] **Data layer:** add an epic-scoped loader in
      `internal/data/release_doc.go` (e.g. `LoadEpicReleaseDocs(epicDir, shortID
      string) ([]ReleaseDoc, error)`) reusing `ReleaseDoc`, `releaseDocSpecs`,
      and the missing/empty/read-error behavior. Build the per-doc filename from
      the short ID + spec (`shortID + "-" + spec.fileName`, e.g. `E33-PRD.md`).
      Keep `LoadReleaseDocs(root)` for the project scope.
- [x] **Board IO:** rename/clarify `loadReleaseDocsCmd` into two commands — one
      epic-scoped (takes `epicDir`, `shortID`) used by the Epic detail overlay,
      and one project-scoped (takes `root`) used by the new top-level overlay.
- [x] **Epic wiring:** at `update.go` openEpicDetail (~line 610) and the Docs
      `case "3"` handler (~line 401/404), call the epic-scoped command with the
      already-computed `epicDir` / `m.epicDetailEpic()`; stop passing `m.Root`.
- [x] **Empty-state copy:** update `releaseDocBody` missing/empty strings so they
      read as epic-scoped when rendered from the epic subview (and project-scoped
      from the top-level overlay) — consider passing a scope label rather than
      hard-coding.
- [x] **Subview retitle:** update the heading in `RenderEpicReleaseDocs`
      ("EPIC DOCS") per the decision above.
- [x] **Top-level overlay:** add `OverlayReleaseDocs` to `model.go`, model state
      for the project docs (loaded docs, selected index, per-doc offsets), a
      global key handler (`P`) that loads project docs and sets the overlay, and
      `view.go` dispatch. Generalize the renderer (extract a shared
      `renderReleaseDocsView`/`RenderReleaseDocs` from `RenderEpicReleaseDocs`)
      so the top-level overlay reuses selector/body/empty-state without the epic
      tab strip.
- [x] **Navigation in the overlay:** `[`/`]` doc switch, scroll keys, and
      `esc`/`q` close — mirror the epic subview handlers.
- [x] **Discoverability:** add the shortcut to the footer hint and `help.go`.
- [x] **Tests:** epic-scoped loader (present/missing/mixed/read-error) in
      `release_doc_test.go`; Docs [3] loads epic docs and shows scoped empty
      state (`update_test.go`/`io_test.go`); top-level overlay opens, renders
      project docs, and closes (`update_test.go`); renderer tests for the shared
      project-docs view (`epic_panel_test.go`).

## Dependencies and Risks

- **Reopens E33 scope.** E33's original tasks are done and the epic is at
  audit-pending; this task extends the epic beyond its audited scope. If the
  team prefers, this could instead be a new epic — flagged for the workflow
  owner. Filed here as an additional task per request.
- **No epic docs exist yet.** No current epic has `E##-PRD.md` / `E##-Design.md`,
  so Docs [3] will render the epic-scoped empty state until those files are
  authored. Authoring per-epic docs (or scaffolding/templating them) is a likely
  follow-up and is out of scope here — this task delivers the loading/rendering
  path and the project-level shortcut, not the document content.
- **Shared renderer regression risk.** Generalizing `RenderEpicReleaseDocs` must
  not change the epic subview's existing layout/tests; keep the epic entry point
  behavior identical aside from the heading retitle.

## Context Log

- Data layer (`internal/data/release_doc.go`): extracted a shared
  `loadReleaseDocs(dir, fileName func(spec) string)` core. `LoadReleaseDocs(root)`
  keeps project scope (`PRD.md`/`Design.md`); new `LoadEpicReleaseDocs(epicDir,
  shortID)` reads the `E##-`-prefixed siblings (`E33-PRD.md`/`E33-Design.md`) and
  does not fall back to project docs. Tests in `release_doc_test.go` cover
  present/missing/mixed/read-error and assert epic scope ignores a project
  `PRD.md` in the same dir.
- Board IO (`internal/board/io.go`): replaced `loadReleaseDocsCmd` with
  `loadEpicReleaseDocsCmd(epicDir, shortID)` (returns `releaseDocsMsg`) and
  `loadProjectReleaseDocsCmd(root)` (returns the new `projectDocsMsg`, added in
  `watch.go`).
- Epic wiring (`internal/board/update.go`): the Docs `case "3"` handler and
  `openEpicDetail` now load via `loadEpicReleaseDocsCmd` using the epic dir /
  short ID, so Docs [3] is epic-scoped. `RenderEpicReleaseDocs` heading retitled
  "RELEASE DOCS" → "EPIC DOCS".
- Top-level overlay: added `OverlayReleaseDocs` and `ProjectDocsState`
  (`ProjectDocs`/`ProjectDocIndex`/`ProjectDocOffsets`) in `model.go`; `P` board
  key opens it and loads project docs; `handleReleaseDocsOverlay` plus
  `selectProjectDoc`/`scrollProjectDoc`/`selectedProjectDocOffset`/
  `clampProjectDocIndex` mirror the epic subview nav on the separate state;
  `view.go` dispatches to the new shared `RenderReleaseDocs` (selector + body, no
  tab strip). The renderer is factored through `renderReleaseDocsOverlay`, which
  reserves header rows from the viewport so both entry points share body/selector/
  empty-state rendering.
- Discoverability: footer hints normalized to a compact consistent style to fit
  80 cols with the new `P:docs` entry; `help.go` documents `P`.
- Tests: `release_docs_overlay_test.go` covers the project renderer (header, no
  tab strip, selector/body, footer, missing-doc state) and wiring (`P` opens +
  loads, project state distinct from epic docs, esc closes, select/scroll with
  per-doc offsets, `View()` renders the overlay); `io_test.go` renamed to the
  project command + added an epic-scoped command test; `update_test.go` key-`3`
  test seeds epic-dir docs.
- Quality gates: `make build && make test` pass; `go vet ./internal/board/
  ./internal/data/` clean.

## Drift Notes

- The "no epic docs exist yet" risk is unchanged: Docs [3] renders the
  epic-scoped empty state (e.g. `(PRD not found at E33-PRD.md)`) until per-epic
  documents are authored. Authoring/scaffolding them remains a follow-up.
- Empty-state copy was left as-is; the document `Path` (`E33-PRD.md` vs `PRD.md`)
  already distinguishes epic from project scope in the "not found" message, so no
  separate scope label was threaded through `releaseDocBody`.

## Redesign (supersedes the above)

After review, the epic-scoped approach was rejected as misleading and removed.
The final design has **two** documentation scopes only, surfaced through a single
top-level overlay:

- **Removed the Epic detail Docs [3] subview entirely** — the epic detail overlay
  is back to Detail [1] / Audit [2]. Deleted `RenderEpicReleaseDocs`,
  `LoadEpicReleaseDocs`, `loadEpicReleaseDocsCmd`, the epic `ReleaseDocs*` state,
  the tab-3 key handling, the `3:Docs` footers, and the open-epic doc preload.
- **Top-level "Release Docs" overlay** (renamed from "Project Docs"), opened by a
  case-sensitive **`D`** (lowercase `d` stays defects). It shows three documents:
  - **Release PRD** — `releases/<release>/<release>-PRD.md` for the selected
    release (e.g. `releases/v1.2/v1.2-PRD.md`); tracks the active release.
  - **Overall PRD** — root `.savepoint/PRD.md` (release-independent).
  - **Overall Design** — root `.savepoint/Design.md` (release-independent).
  There is no release-level Design.
- `data.LoadReleaseDocs(root, release)` now returns those three via a spec list
  whose `rel(release)` resolves each path; new IDs `ReleaseDocReleasePRD`,
  `ReleaseDocOverallPRD`, `ReleaseDocOverallDesign`. `loadReleaseDocsCmd(root,
  release)` feeds the one `releaseDocsMsg`; board state collapsed to a single
  `ReleaseDocsState` (the duplicate `ProjectDocs*` set was deleted).
- Footer hint `D:docs`; help row updated. Tests reworked accordingly
  (`release_doc_test.go`, `io_test.go`, `release_docs_overlay_test.go`; epic doc
  tests removed from `epic_panel_test.go` / `update_test.go`). `make build &&
  make test` pass; `go vet ./...` clean.

The renderer internals retained from the earlier work (`renderReleaseDocBody`/
`wrapDocText`/`styledWrap`/`releaseDocBody`/`renderReleaseDocSelector`) still back
the overlay body; the whitespace-preserving wrapper fix carried over unchanged.
