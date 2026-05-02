---
type: quality-review
epic: E01-scaffolding
reviewed: 2026-04-27
advisory: true
---

# Quality Review: E01-scaffolding

Semantic review of the scaffolded source against the 10 Code Style rules from `AGENTS.md`. Advisory only — no blocking issues.

---

## Rule-by-rule findings

### 1. One job per file ✓

`src/cli.ts` handles argument dispatch and output. `src/version.ts` is a single export. No file does two things.

### 2. One-sentence rule ✓

`printHelp()` — prints the placeholder help banner. `main()` — dispatches on argv and calls the appropriate output function. Both are one-sentence describable.

### 3. Test what branches ⚠ advisory

`src/cli.ts` has a branch on `command === "--version"`. The smoke test in `test/smoke.test.ts` only covers `src/version.ts`. The CLI dispatch logic is untested.

This is acceptable at scaffold stage — the branches contain only placeholder output and will be replaced in `cli-foundation`. Flag for that epic: when real dispatch lands, cover every branch.

### 4. Types are documentation ✓

No `any` anywhere in source or test. Return types are declared (`void`). ESLint is configured to error on `@typescript-eslint/no-explicit-any`.

### 5. Build, don't speculate ✓

No abstraction layers, registries, or plugin systems. Placeholder CLI is minimal: two functions and a `main()` call. Nothing exists for hypothetical future behavior.

### 6. Errors at boundaries ✓ (n/a at scaffold)

No real system boundaries in placeholder code. The CLI entrypoint reads only `process.argv`, which cannot throw. Acceptable for a placeholder.

### 7. One source of truth ✓

Version string is declared once in `src/version.ts` and imported by `src/cli.ts`. `package.json` also contains `"version": "0.1.0"` — this is expected (npm requires it there), not a violation. When the CLI reads its own version at runtime, it imports from `src/version.ts`, not from `package.json` directly.

> **Recommendation for `cli-foundation`:** consider whether `version.ts` should derive from `package.json` at build time (via `tsup` banner injection or a build-time const) to eliminate the manual sync requirement. Not urgent at `0.1.0`.

### 8. Comments explain WHY ✓

No comments in any source file. Config files have no inline comments. Good.

### 9. Content in data files ✓ (n/a at scaffold)

No user-facing content exists yet. Config is in conventional JSON/TS config files.

### 10. Small diffs ✓

Epic split into 5 tasks with clear separation: package metadata → TypeScript build → test runner → lint/format → verification. Each task touched a minimal set of files with declared dependencies.

---

## Structural observations

### Dual vitest config files

`vitest.config.ts` and `vitest.config.js` are identical. The `.js` variant is the one actually used by `npm test` (required for `--configLoader runner`). The `.ts` variant is unused but exists as the "source of truth" for IDE tooling and type checking.

**Recommendation:** document this in a comment in `vitest.config.ts`, or remove it and keep only `.js`. The dual-file situation will confuse future contributors who edit `.ts` expecting it to take effect.

### scripts/vitest-preload.cjs complexity

The preload script stubs two unrelated things: a `net use` child_process call (Windows network drive detection in Vitest's environment) and `esbuild.transform`/`transformSync` (a secondary Windows compat issue). The stubs are correct and well-structured.

**Recommendation:** add a single-line comment per stub explaining the specific Windows issue being worked around (bug reference or symptom description). The WHY is not obvious — a future contributor would not know whether these stubs are still needed after a Vitest upgrade.

### test/smoke.test.ts include pattern

`vitest.config.js` uses `include: ["test/**/*.test.js"]` but the actual test file is `smoke.test.ts`. This works because Vitest transforms TypeScript to JS before matching, but the pattern is misleading.

**Recommendation:** change include to `["test/**/*.test.ts"]` to match what actually exists in the repo.

---

## Must Fix Before Close

None.

## Carry Forward

- Add CLI branch tests when real dispatch is implemented in `E03-cli-foundation`.
- Decide in `E03-cli-foundation` whether `src/version.ts` should derive from `package.json`.

## Already Fixed

- `vitest.config.ts` now documents the dual-config setup.
- `scripts/vitest-preload.cjs` now explains the Windows `net use` and esbuild workarounds.
- Vitest include patterns now match `test/**/*.test.ts`.

## Summary

No blocking issues remain for `E01-scaffolding`.
