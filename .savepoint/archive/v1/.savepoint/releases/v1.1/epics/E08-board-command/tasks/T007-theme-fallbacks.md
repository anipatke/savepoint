---
id: E08-board-command/T007-theme-fallbacks
status: done
objective: "Implement theme loading and color fallbacks"
depends_on: ["E08-board-command/T003-tui-app-shell"]
---

# T007: Theme Fallbacks

## Acceptance Criteria

- Load theme tokens from config.yml
- Support truecolor (24-bit) terminals
- Fall back to 256-color mode if no truecolor
- Fall back to 16-color if no 256-color
- Support NO_COLOR=1 for monochrome
- Apply theme to all UI elements (columns, cards, detail, overlays)

## Implementation Plan

- [x] Add `internal/board/theme.go` (extend existing or create new)
- [x] Implement theme loading from config.yml
- [x] Implement truecolor profile detection
- [x] Implement 256-color fallback mapping
- [x] Implement 16-color fallback mapping
- [x] Support NO_COLOR=1 monochrome mode
- [x] Apply theme to columns, cards, detail, overlays
- [x] Test on different terminal profiles
- [x] Run `make build && make test`