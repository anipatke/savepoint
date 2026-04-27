---
id: E01-scaffolding/T004-lint-format-gates
status: done
objective: "Add ESLint and Prettier configuration so the scaffold has mechanical quality gates from the first epic."
depends_on: [E01-scaffolding/T003-vitest-smoke]
---

# Task T004: Lint And Format Gates

## Implementation Plan

- [x] Add ESLint, TypeScript ESLint, Prettier, and supporting config dependencies to `package.json`.
- [x] Add `eslint.config.js` with flat-config rules that reject `any` in TypeScript.
- [x] Add `.prettierrc.json` and ignore rules for generated output.
- [x] Add `lint` and `format:check` scripts to `package.json`.
- [x] Verify `npm run lint` and `npm run format:check` pass on the current scaffold.
- [x] Verify the existing scaffold still passes `npm run build`, `npm test`, and `npm run typecheck`.

## Scope

Add lint and formatting configuration for the TypeScript project.

Expected files:

- `eslint.config.js`
- `.prettierrc.json`

Expected `package.json` updates:

- Add ESLint, TypeScript ESLint, and Prettier dependencies.
- Add `lint` and `format:check` scripts.

## Acceptance Criteria

- `npm run lint` passes.
- `npm run format:check` passes.
- ESLint is configured with the flat config format.
- ESLint rejects `any` usage in project TypeScript.
- Prettier covers source, test, config, and markdown files.

## Out Of Scope

- Language-specific quality presets beyond TypeScript.
- Dependency boundary tooling.
- Auto-format command requirements beyond a check script.
