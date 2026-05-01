---
type: epic-design
status: audited
audited: 2026-04-27
---

# Epic E01: scaffolding

## Purpose

Create the minimum production-ready TypeScript package foundation for `savepoint`, including build, test, lint, formatting, package metadata, and source tree conventions. This epic establishes the first auditable codebase baseline without implementing user-facing Savepoint behavior.

## What this epic adds

- A single-package Node 20.10+ TypeScript project using ESM.
- npm package metadata for package name `savepoint` and binary name `savepoint`.
- A `tsup` build that emits `dist/cli.js` with a Node shebang.
- Strict TypeScript configuration.
- Vitest test setup for future unit and integration tests.
- ESLint and Prettier configuration for mechanical quality gates.
- A minimal source layout that future epics can extend.
- Baseline npm scripts for build, typecheck, lint, format check, and test.
- Root documentation and ignore files needed for ordinary OSS development.

## Components and files

Expected files introduced by this epic:

| Path                 | Purpose                                                                  |
| -------------------- | ------------------------------------------------------------------------ |
| `package.json`       | Package metadata, bin mapping, scripts, dependencies, files allowlist.   |
| `package-lock.json`  | Locked dependency graph for repeatable installs.                         |
| `tsconfig.json`      | Strict TypeScript compiler settings for ESM Node output.                 |
| `tsup.config.ts`     | Build entry and output configuration.                                    |
| `vitest.config.ts`   | Test runner configuration.                                               |
| `eslint.config.js`   | Flat ESLint config for TypeScript.                                       |
| `.prettierrc.json`   | Formatting rules.                                                        |
| `.gitignore`         | Ignore generated artifacts and local environment files.                  |
| `README.md`          | Minimal development README until release docs expand it.                 |
| `LICENSE`            | MIT license placeholder for the package.                                 |
| `src/cli.ts`         | Minimal executable entrypoint that prints help/version placeholder text. |
| `src/version.ts`     | Single source for the package version used by the placeholder CLI.       |
| `test/smoke.test.ts` | Baseline test proving the test runner and exported constants work.       |

## Architectural delta

Before this epic, the repository contains only planning markdown and Savepoint workflow metadata. After this epic, it becomes an installable npm package skeleton with a compilable CLI entrypoint.

This does not add the actual Savepoint command system. The placeholder `src/cli.ts` exists only to validate packaging, shebang output, and the future binary surface. Real command dispatch belongs to `cli-foundation`.

## Boundaries

In scope:

- Configure the repo as an ESM-only TypeScript package.
- Make `npm run build`, `npm run typecheck`, `npm run lint`, and `npm test` meaningful.
- Keep runtime code minimal and dependency-light.
- Use package and binary name `savepoint`.

Out of scope:

- Argument parsing beyond placeholder help/version output.
- Reading or writing `.savepoint` project files.
- Implementing `init`, `board`, `audit`, or `doctor`.
- Ink, React, TUI rendering, theming, or visual identity work.
- Template content or Router Pattern prompt assets.
- Release publishing automation beyond package metadata basics.

## Quality gates

The scaffold should support these commands:

```bash
npm run build
npm run typecheck
npm run lint
npm test
```

The first audit after this epic should be able to update `AGENTS.md` with a codebase map that reflects the new top-level modules.

## Design constraints

- Keep files small and obvious; no framework or architecture abstractions before there is behavior to organize.
- Use strict typing from the first commit.
- Avoid `any`.
- Keep config conventional enough that a vibe coder can recognize and modify it.
- Do not add native dependencies.
- Do not add telemetry, update checks, analytics, or postinstall scripts.

## Open decisions

- The exact dependency versions are implementation details, but they should be current stable versions compatible with Node 20.10+.
- `README.md` should stay intentionally sparse in this epic; fuller public-facing docs belong to `docs-and-packaging`.

## Implemented as

Implemented on 2026-04-27 across tasks T001–T005. All quality gates pass.

**Deviations from original design:**

| Item                         | Deviation                | Reason                                                                                                                                                                                     |
| ---------------------------- | ------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `vitest.config.js`           | Extra file not in design | `npm test` requires `--configLoader runner` pointing at a `.js` config on Windows; `.ts` config cannot be loaded directly in this execution path                                           |
| `scripts/vitest-preload.cjs` | Extra file not in design | Windows runner executes a `net use` child_process call during Vitest startup; this CJS preload stubs that call and also stubs esbuild transforms to avoid a secondary Windows compat issue |
| `.prettierignore`            | Extra file not in design | Standard convention; consistent with `.gitignore` exclusions                                                                                                                               |

**Dependency versions (pinned):**

| Package             | Version |
| ------------------- | ------- |
| `typescript`        | ^5.8.3  |
| `tsup`              | ^8.5.1  |
| `vitest`            | ^3.2.4  |
| `eslint`            | ^9.33.0 |
| `typescript-eslint` | ^8.39.0 |
| `prettier`          | ^3.6.2  |
| `@types/node`       | ^24.1.0 |
