---
id: E04-epic-navigation/T003-epic-status-glyphs
status: done
objective: "Add status glyphs to the epic sidebar showing each epic's lifecycle state"
depends_on:
  - E04-epic-navigation/T001-sidebar-focusable-navigation
---

# T003: Epic Status Glyphs in Side Panel

## Acceptance Criteria

- Each epic in the sidebar shows a status glyph prefix: `○` for `planned`, `▶` for `in_progress`, `◉` for `done`, `✓` for `audited`, space for unknown
- The glyph is styled: `○` in dim, `▶` in orange, `◉` in green, `✓` in green — matching the status colour convention
- The status is read from the `status:` field in each epic's `E##-Detail.md` frontmatter at load time
- Only the **side panel** is affected — the `E` key dropdown overlay is unchanged
- The status map is loaded once during `loadBoardData` and cached in the Model (no per-frame file I/O)
- If an epic's detail file is missing or has no `status` frontmatter, the glyph defaults to a space (no indicator)
- All existing epic sidebar rendering (focus highlighting, cursor, active marker, selected marker) works unchanged
- `make build && make test` pass

## Implementation Plan

- [x] Add `EpicStatus map[string]string` to `Model` in `model.go` — keyed by epic ID, values like `"planned"`, `"in_progress"`, `"done"`, `"audited"`
- [x] In `board.go` — update `loadBoardData` to load epic statuses:
  - After listing epics for a release, iterate each epic's directory
  - Build detail file path: `filepath.Join(epic.Path, shortID(epic.ID)+"-Detail.md")`
  - Read the file with `os.ReadFile`
  - Parse frontmatter with `data.NewParser().ParseFrontmatter(string(content))`
  - Extract `"status"` from the returned map; if missing, skip
  - Store in `epicStatuses[epic.ID]` — collect all epics across releases into the map
- [x] Pass `epicStatuses map[string]string` through the `reloadMsg` pipeline so dynamic reloads preserve glyphs
- [x] In `epic_panel.go` — update `RenderEpicSidebar` signature to accept `status map[string]string`
- [x] Add a helper function `epicSidebarGlyph(status map[string]string, epicID string) string`:
  - `"planned"` → styles.CardMeta.Render("○")
  - `"in_progress"` → styles.GlyphBuild.Render("▶")
  - `"done"` → styles.TagDone.Render("◉")
  - `"audited"` → styles.TagDone.Render("✓")
  - default → " " (single space)
- [x] In `RenderEpicSidebar` — prepend the glyph to each epic label: `"▶ E01-tui-optimisation"`
- [x] In `view.go` — pass `m.EpicStatus` to `RenderEpicSidebar`
- [x] Handle `EpicStatus` in `reloadMsg` handler in `update.go` — store the map when tasks reload
- [x] Update `reloadMsg` struct in `watch.go` to carry `epicStatuses map[string]string`
- [x] Run `make build && make test` to verify no regressions

## Context Log

Files read:
- `internal/board/model.go`
- `internal/board/epic_panel.go`
- `internal/board/view.go`
- `internal/board/board.go`
- `internal/board/watch.go`
- `internal/board/update.go`
- `internal/board/card.go` — glyph pattern reference
- `internal/styles/styles.go` — existing style constants
- `internal/data/parser.go` — `ParseFrontmatter` for YAML extraction
- `.savepoint/releases/v1.1/epics/E01-tui-optimisation/E01-Detail.md` — example frontmatter with status field
- `.savepoint/releases/v1.1/epics/E04-epic-navigation/tasks/T001-sidebar-focusable-navigation.md` — dependency status check
- `internal/board/epic_panel_test.go` — updated RenderEpicSidebar test calls

Estimated input tokens: 1800

Notes:
- Reuses `data.NewParser().ParseFrontmatter()` which already handles YAML frontmatter extraction
- Reuses `shortID()` from `card.go` (same package) to derive the detail file name
- The `E` key dropdown in `epic_panel.go` (`RenderEpicDropdown`) is intentionally NOT updated — side panel only per user preference
- Function named `epicSidebarGlyph` (not `epicStatusGlyph`) — takes `(map[string]string, string)` for nil-safe lookup
- `go build` + `go test ./...` pass

## Drift Notes

No drift — all changes are within existing files and follow established patterns.
