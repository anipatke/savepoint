---
id: E01-scaffolding/T003-vitest-smoke
status: done
objective: "Add Vitest configuration and a smoke test that proves the scaffolded package can run automated tests."
depends_on: [E01-scaffolding/T002-typescript-build]
---

# Task T003: Vitest Smoke Test

## Implementation Plan

- [x] Add Vitest configuration and include pattern for the scaffold smoke test.
- [x] Add a smoke test that imports scaffolded source and asserts the exported version.
- [x] Add the Vitest dependency and `test` script to `package.json`.
- [x] Adjust the scaffold for the Windows runner path so `npm test` executes reliably in this environment.
- [x] Verify `npm test` and `npm run typecheck` pass.

## Scope

Introduce the test runner and a small smoke test for the existing scaffolded source.

Expected files:

- `vitest.config.ts`
- `test/smoke.test.ts`

Expected `package.json` updates:

- Add Vitest dependency.
- Add `test` script.

## Acceptance Criteria

- `npm test` runs Vitest.
- The smoke test imports TypeScript source.
- The smoke test verifies stable placeholder behavior such as exported version metadata.
- Tests do not depend on generated `dist/` output.

## Out Of Scope

- CLI integration tests.
- Command behavior tests.
- Filesystem fixture tests.
