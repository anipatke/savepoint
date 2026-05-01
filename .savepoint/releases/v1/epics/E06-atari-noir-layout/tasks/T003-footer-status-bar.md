---
id: E06-atari-noir-layout/T003-footer-status-bar
status: done
objective: "Implement the static Footer Status Bar and bottom divider"
depends_on: ["E06-atari-noir-layout/T002-header-and-dividers"]
---

# T003: Implement Footer Status Bar

## Acceptance Criteria

- Footer replaces the dynamic status string with exactly: `PLAN      │      BUILD      │      AUDIT`.
- PLAN is colored Vibe Purple.
- BUILD is colored Atari Orange.
- AUDIT is colored NPP Green.
- The divider `│` is colored Border Subtle.
- The text is bold or slightly emphasized, with even spacing between sections.
- A full-width horizontal divider is rendered above the footer.
- The footer contains no icons, no repetition of the logo, and no extra labels/metrics.
- The main application content sits cleanly between the top and bottom dividers without overlapping.

## Implementation Plan

- [x] Update `renderStatus` in `internal/board/view.go` to render the new static phase footer.
- [x] Style the footer segments according to the color rules.
- [x] Add the bottom horizontal divider above the footer.
- [x] Adjust height calculations so the board area fits exactly between the top and bottom dividers. (No height calculation changes needed — `lipgloss.JoinVertical` naturally accommodates the extra divider line.)

## Context Log

Files read:
- `.savepoint/router.md`
- `.savepoint/releases/v1/epics/E06-atari-noir-layout/Design.md`
- `.savepoint/releases/v1/epics/E06-atari-noir-layout/tasks/T003-footer-status-bar.md`
- `.savepoint/releases/v1/epics/E06-atari-noir-layout/tasks/T002-header-and-dividers.md`
- `.savepoint/visual-identity.md`
- `internal/board/view.go`
- `internal/board/view_test.go`
- `internal/board/layout.go`
- `internal/board/model.go`
- `internal/styles/palette.go`
- `internal/styles/styles.go`

Estimated input tokens: ~5500

Notes:
- Replaced `renderStatus` (dynamic status bar with release/epic/column info) with `renderFooter` (static phase footer: PLAN │ BUILD │ AUDIT).
- Added `FooterPhasePlan` (purple, bold), `FooterPhaseBuild` (orange, bold), `FooterPhaseAudit` (green, bold), and `FooterDivider` (border subtle) styles.
- Added full-width bottom divider above the footer matching the top divider style.
- Removed unused `fmt` import from `view.go` after `renderStatus` deletion.
- `StatusBar` style and `StatusMessage` model field remain unused but harmless; left for future cleanup.
- Quality gates: `go build ./...` pass, `go test ./...` pass, `go vet ./...` pass.
