# E01-scaffolding Handoff

Read this instead of the full audit unless debugging scaffold history.

- Package scaffold exists for `savepoint` version `0.1.0`.
- Commands available through npm: `build`, `typecheck`, `lint`, `format:check`, and `test`.
- Windows test execution uses `scripts/vitest-preload.cjs` plus `vitest.config.js`.
- Source baseline: `src/cli.ts`, `src/version.ts`, and `test/smoke.test.ts`.
- Carried forward to `E03-cli-foundation`: real CLI dispatch tests and the decision on whether `version.ts` should derive from `package.json`.
