---
id: E07-init-command/T006-clipboard
status: done
objective: "Implement best-effort clipboard copy"
depends_on: ["E07-init-command/T005-magic-prompt"]
---

# T006: Clipboard

## Acceptance Criteria

- Detects platform (Windows/macOS/Linux)
- Copies magic prompt to clipboard
- Does NOT fail if clipboard unavailable (best-effort)
- Returns success/skipped/failed status without throwing on error
- Works on PowerShell, iTerm2, terminal emulators

## Implementation Plan

- [x] Add `internal/init/clipboard.go`
- [x] Implement `CopyToClipboard(text) ClipboardResult`
- [x] Detect platform: Windows (clip.exe), macOS (pbcopy), Linux (xclip/xsel)
- [x] Fall back gracefully if clipboard tools unavailable
- [x] Log skip or failure without failing init
- [x] Test clipboard status handling on the current platform
- [x] Run quality gates (`go test ./internal/init ./cmd`, `go build ./...`, `go test ./...`)
