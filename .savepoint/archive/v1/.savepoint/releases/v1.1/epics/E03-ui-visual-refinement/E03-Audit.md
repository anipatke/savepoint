---
type: audit-findings
audited: 2026-05-02
---

# Audit Findings: E03 UI Visual Refinement

## Main Findings

E03 verified the visual refinement work around the board header, Next Activity line, task detail checklist rendering, shared status glyphs, and deterministic ANSI256 color profile behavior.

The audit found no blocking product-code issue in the reviewed E03 scope. The main drift was documentation/process cleanup: update architecture notes to match the implemented rendering behavior, record where the color profile was actually applied, and check off the completed T006 implementation plan.

## Code Style Review


- [ ] One job per file
- [ ] One-sentence functions
- [ ] Test branches
- [ ] Types are documentation
- [ ] Build, don't speculate
- [ ] Errors at boundaries
- [ ] One source of truth
- [ ] Comments explain WHY
- [ ] Content in data files
- [ ] Small diffs

## Proposed Changes

### Target File

`.savepoint/Design.md`

### Replace

```yaml
last_audited: v1.1/E02-cross-platform-compatibility (2026-05-02)
```

### With

```yaml
last_audited: v1.1/E03-ui-visual-refinement (2026-05-02)
```

---

### Target File

`.savepoint/Design.md`

### Replace

```markdown
- **Board command** (`savepoint board`) reads project state, renders the Atari-Noir TUI board, supports release/epic filtering, detail overlays, task status transitions with mtime-guarded writes, release/epic-scoped router priority markers, fsnotify-based task auto-refresh (epic E06), header Next Activity display, height-aware column/detail viewport scrolling, stable focused/unfocused column border geometry (v1.1 E01), and a focusable wide-screen epic sidebar with purple epic focus, epic detail overlays, and status glyphs loaded from epic detail frontmatter (v1.1 E04).
```

### With

```markdown
- **Board command** (`savepoint board`) reads project state, renders the Atari-Noir TUI board, supports release/epic filtering, detail overlays, task status transitions with mtime-guarded writes, release/epic-scoped router priority markers, fsnotify-based task auto-refresh (epic E06), header Next Activity display, height-aware column/detail viewport scrolling, stable focused/unfocused column border geometry (v1.1 E01), dedicated phase-colored Next Activity line below the header, sentence-boundary checklist rendering in task details, shared status glyph mapping for task cards and the epic sidebar, a forced ANSI256 Lipgloss color profile for board startup (v1.1 E03), and a focusable wide-screen epic sidebar with purple epic focus, epic detail overlays, and status glyphs loaded from epic detail frontmatter (v1.1 E04).
```

---

### Target File

`.savepoint/Design.md`

### Replace

```markdown
**Layout:** single screen with a 3-column task board (`planned`, `in_progress`, `done`), optional epic sidebar on wide terminals, centered overlays for release/epic/help/task/epic-detail views, static Atari-Noir header/footer, full-width dividers, uniform black TUI backgrounds, and navigation hints. The header can show a compact right-aligned Next Activity value from router state. Columns and detail overlays use height-aware viewport slicing with subtle above/more scroll indicators. Focused and unfocused columns preserve the same rounded-border geometry so focus changes do not shift content. On terminals at least 120 columns wide, the epic sidebar is focusable from the Planned column; it uses the purple epic accent for focused panel borders, focused epic labels, and epic detail overlays while task-column focus remains orange. Non-TTY output remains a plain table fallback.
```

### With

```markdown
**Layout:** single screen with a 3-column task board (`planned`, `in_progress`, `done`), optional epic sidebar on wide terminals, centered overlays for release/epic/help/task/epic-detail views, static Atari-Noir header/footer, full-width dividers, uniform black TUI backgrounds, and navigation hints. Active router `next_action` renders as a dedicated full-width line below the header with phase-colored `PLAN`, `BUILD`, or `AUDIT` prefix styling and truncates to terminal width. Columns and detail overlays use height-aware viewport slicing with subtle above/more scroll indicators. Focused and unfocused columns preserve the same rounded-border geometry so focus changes do not shift content. Task detail implementation-plan checkboxes render once per semantic sentence, not once per hard-wrapped markdown line. On terminals at least 120 columns wide, the epic sidebar is focusable from the Planned column; it uses the purple epic accent for focused panel borders, focused epic labels, and epic detail overlays while task-column focus remains orange. Task card and epic sidebar status glyphs share `internal/board/status.go`; task cards use explicit `Task.Status` when available and retain the legacy column/stage glyph fallback when it is not. Non-TTY output remains a plain table fallback.
```

---

### Target File

`AGENTS.md`

### Replace

```markdown
| `internal/board/` | TUI board, overlays, epic sidebar, status glyphs |
```

## With

```markdown
| `internal/board/` | TUI board, overlays, epic sidebar, Next Activity line, detail checklist rendering, status glyphs, forced color profile |
```

---

### Target File

`.savepoint/releases/v1.1/epics/E03-ui-visual-refinement/E03-Detail.md`

### Replace

```yaml
status: planned
```

### With

```yaml
status: audited
```

---

### Target File

`.savepoint/releases/v1.1/epics/E03-ui-visual-refinement/E03-Detail.md`

### Replace

```markdown
| `internal/board/card.go` | Task card glyph determination (updated for shared helper) |
| `internal/board/view_test.go` | Header, formatting, and checkbox tests |
```

### With

```markdown
| `internal/board/card.go` | Task card glyph determination (updated for shared helper) |
| `internal/board/board.go` | Board startup; sets the Lipgloss color profile to ANSI256 before model initialization |
| `internal/board/view_test.go` | Header, formatting, and checkbox tests |
```

---

### Target File

`.savepoint/releases/v1.1/epics/E03-ui-visual-refinement/E03-Detail.md`

### Insert After

```markdown
| `internal/board/card.go` | Task card glyph determination (updated for shared helper) |
| `internal/board/board.go` | Board startup; sets the Lipgloss color profile to ANSI256 before model initialization |
| `internal/board/view_test.go` | Header, formatting, and checkbox tests |
```

### With

```markdown
## Implemented as

- `internal/board/view.go` renders `next_action` as a separate line below the header through `renderNextActivityLine`, using existing footer phase styles for `PLAN`, `BUILD`, and `AUDIT`.
- `internal/board/layout.go` accounts for the optional Next Activity line when calculating board chrome and content height.
- `internal/data/parser.go` joins hard-wrapped checklist continuation lines before rendering so markdown wrap points do not create duplicate checklist items.
- `internal/board/detail.go` splits checklist item text on semantic sentence boundaries and emits one `[ ]` or `[x]` marker per sentence.
- `internal/data/task.go` adds `Task.Status` plus status constants, including `audited`, for shared board glyph rendering.
- `internal/board/status.go` centralizes the planned, in-progress, done, and audited status glyph mapping used by task cards and the epic sidebar.
- `internal/board/card.go` uses explicit `Task.Status` when present, while preserving the legacy column/stage glyph fallback for older task data.
- `internal/board/epic_panel.go` delegates epic status glyph rendering to the shared status helper.
- `internal/board/board.go` sets Lipgloss to the ANSI256 color profile at board startup. This satisfies the deterministic 256-color rendering intent, although the implementation lives at the board boundary instead of `main.go`.
```

---

### Target File

`.savepoint/releases/v1.1/epics/E03-ui-visual-refinement/tasks/T006-forced-256-color-profile.md`

### Replace

```markdown
## Implementation Plan

- [ ] Read `main.go` — understand startup flow and identify where to inject profile forcing
- [ ] Edit `main.go` — add `lipgloss.SetColorProfile(lipgloss.Force256Color)` call before `board.Run()`
- [ ] Run `make build && make test` to verify no regressions
```

### With

```markdown
## Implementation Plan

- [x] Read board startup flow and identify where to inject profile forcing
- [x] Edit `internal/board/board.go` — add `lipgloss.SetColorProfile(termenv.ANSI256)` before model initialization
- [x] Run equivalent build and test gates in this Windows shell
```

---
