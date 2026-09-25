---
id: E01-go-setup/T002-entrypoint
status: done
objective: "Create main.go that launches a blank Bubble Tea program"
depends_on: [E01-go-setup/T001-init-module]
---

# T002: Create Entrypoint

## Acceptance Criteria

- `main.go` exists and compiles.
- Running the binary launches a Bubble Tea program.
- Pressing `q` or `ctrl+c` exits cleanly.
- No TypeScript files remain at project root.

## Implementation Plan

- [x] Write `main.go` with `tea.NewProgram` and a minimal model.
- [x] Implement `Init()`, `Update()`, `View()` on the model.
- [x] Handle `tea.KeyMsg` for `q` and `ctrl+c`.
- [x] Verify `go run main.go` launches and quits.
- [x] Delete `src/`, `test/`, `dist/`, `package.json`, `node_modules/`, `tsconfig.json`, `tsup.config.ts`, `vitest.config.js`, `eslint.config.js`, `.prettierrc.json`.
