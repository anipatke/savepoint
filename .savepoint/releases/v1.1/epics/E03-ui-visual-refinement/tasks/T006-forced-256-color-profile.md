---
id: E03-ui-visual-refinement/T006-forced-256-color-profile
status: planned
objective: "Force 256-color terminal profile at init to ensure consistent background rendering across all terminals"
depends_on: []
---

# T006: Forced 256-Color Profile for Terminal Consistency

## Acceptance Criteria

- On startup, `lipgloss` is configured to use `Force256Color` profile instead of relying on terminal detection
- The board's background renders as uniform black (`"#000000"` / ANSI256 `"232"`) on every color-capable terminal including PowerShell, CMD, Windows Terminal, VS Code terminal, and all Unix terminals
- No color artifacts or background mismatches are visible when running on PowerShell
- Existing style definitions using TrueColor hex values still resolve correctly via Lipgloss's 256-color fallback
- `make build && make test` pass with no regressions

## Implementation Plan

- [ ] Read `main.go` — understand startup flow and identify where to inject profile forcing
- [ ] Edit `main.go` — add `lipgloss.SetColorProfile(lipgloss.Force256Color)` call before `board.Run()`
- [ ] Run `make build && make test` to verify no regressions

## Context Log

Files read:
- `main.go`
- `internal/styles/palette.go`
- `internal/styles/styles.go`

Estimated input tokens: 400

Notes:
- `lipgloss.Force256Color` is supported on all Go platforms and all color-capable terminals
- This is a single-line change — no new dependencies or imports beyond the existing `github.com/charmbracelet/lipgloss` package
- TrueColor hex values (`"#000000"`) are mapped to nearest 256-color values (`"232"`) via palette.go's `color()` helper; these are already configured
