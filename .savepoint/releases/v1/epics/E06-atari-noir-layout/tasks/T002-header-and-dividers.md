---
id: E06-atari-noir-layout/T002-header-and-dividers
status: done
objective: "Build the SavePoint Header and horizontal dividers"
depends_on: ["E06-atari-noir-layout/T001-color-system"]
---

# T002: Build SavePoint Header and Dividers

## Acceptance Criteria

- Header renders exactly `▣  S A V E P O I N T` at the top of the TUI.
- Icon `▣` is styled with Atari Orange.
- Text `S A V E P O I N T` is Primary Text, uppercase, with spacing between letters.
- Horizontal dividers `────────────────────────────────────────────────────────` are added below the header.
- Divider color is strictly Border Subtle.
- The header renders cleanly without wrapping, even in narrow terminals.
- No additional icons or animations are present in the header.

## Implementation Plan

- [x] Update `internal/board/view.go` to render the specific header string.
- [x] Add logic in `internal/board/view.go` to insert the top horizontal divider.
- [x] Adjust `internal/board/layout.go` / `internal/board/view.go` calculations if needed so the board fits below the header+divider.

## Context Log

Files read:
- `.savepoint/router.md`
- `.savepoint/releases/v1/epics/E06-atari-noir-layout/Design.md`
- `.savepoint/releases/v1/epics/E06-atari-noir-layout/tasks/T002-header-and-dividers.md`
- `.savepoint/releases/v1/epics/E06-atari-noir-layout/tasks/T001-color-system.md`
- `.savepoint/visual-identity.md`
- `internal/styles/palette.go`
- `internal/styles/styles.go`
- `internal/board/view.go`
- `internal/board/view_test.go`
- `internal/board/layout.go`
- `internal/board/epic_panel_test.go`
- `internal/board/release_test.go`

Estimated input tokens: ~4500

Notes:
- `layout.go` needed no changes — `lipgloss.JoinVertical` naturally accommodates the extra divider line without explicit height accounting.
- The old `Header` style (orange, bold, padding) was replaced with `HeaderIcon` (orange, bold) and `HeaderText` (primary text) to support two-tone header rendering.
- `Divider` style added with `BorderSubtle` foreground for quiet full-width separators.
- Quality gates: `go build ./...` pass, `go test ./...` pass, `go vet ./...` pass.

## Drift Notes

- Drift: `internal/styles/palette.go` and `internal/styles/styles.go` not yet in Codebase Map — pre-existing omission. These are the Atari-Noir color/style modules active across all board rendering.
