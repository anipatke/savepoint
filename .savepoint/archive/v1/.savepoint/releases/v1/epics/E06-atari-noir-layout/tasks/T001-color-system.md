---
id: E06-atari-noir-layout/T001-color-system
status: done
objective: "Implement Atari-Noir Color System in palette and styles"
depends_on: []
---

# T001: Implement Atari-Noir Color System

## Acceptance Criteria

- `internal/styles/palette.go` updated with exact hex values for Background, Surface, Surface 2, Border, Border Subtle, Primary Text, Atari Orange, NPP Green, Vibe Purple.
- Legacy colors or gradients removed from `palette.go`.
- `internal/styles/styles.go` maps styles according to usage rules (e.g., Background is dark `#121212`, Primary text is `#F0E6DA`, Atari Orange `#FC6323` for active/focus, Vibe Purple `#B1A1DF` for metadata, NPP Green `#A4C639` for success).
- No large bright background fills or mixed accent colors in a single element exist in the style definitions.

## Implementation Plan

- [x] Update `internal/styles/palette.go` with Atari-Noir hex, ANSI256, and ANSI16 values.
- [x] Refactor `internal/styles/styles.go` variables to use the updated palette constants strictly.
- [x] Remove unused or conflicting colors.

## Context Log

Files read:
- `internal/styles/palette.go`
- `internal/styles/styles.go`
- `.savepoint/visual-identity.md`
- `.savepoint/router.md`
- `E06-atari-noir-layout/Design.md`

Estimated input tokens: 2400

Notes:
- palette.go already contained correct Atari-Noir hex/256/16 values matching the visual identity spec — no changes needed.
- styles.go already referenced palette constants strictly — no circular refs or string literals.
- No legacy colors or gradients found in palette.go.
- No large bright background fills or mixed accent colors found in styles.go.
- Quality gates: `go build ./...` pass, `go test ./...` pass (all cached), `go vet` pass.
