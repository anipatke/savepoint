---
id: E01-scaffolding/T002-typescript-build
status: done
objective: "Add the strict TypeScript source layout and `tsup` build that emits the placeholder `savepoint` CLI binary."
depends_on: [E01-scaffolding/T001-package-baseline]
---

# Task T002: TypeScript Build

## Implementation Plan

- [x] Add `tsconfig.json` with strict NodeNext compilation settings for the scaffold.
- [x] Add `tsup.config.ts` and a minimal `src/cli.ts` entrypoint that builds to `dist/cli.js`.
- [x] Add `src/version.ts` as the single source of version metadata used by the CLI placeholder.
- [x] Add TypeScript and build dependencies to `package.json`.
- [x] Generate `package-lock.json` from the updated manifest.
- [x] Verify `npm run build` and `npm run typecheck` pass.

## Scope

Create the minimal source tree and build configuration needed for a compiled CLI placeholder.

Expected files:

- `tsconfig.json`
- `tsup.config.ts`
- `src/cli.ts`
- `src/version.ts`
- `package-lock.json`

Expected `package.json` updates:

- Add TypeScript/build dependencies.
- Add `build` and `typecheck` scripts.

## Acceptance Criteria

- `npm run build` emits `dist/cli.js`.
- The emitted CLI file has a Node shebang.
- `npm run typecheck` passes under strict TypeScript settings.
- The placeholder CLI can print version/help placeholder output.
- No real command dispatch is implemented.

## Out Of Scope

- Argument parsing beyond minimal placeholder output.
- `init`, `board`, `audit`, or `doctor` behavior.
- Tests beyond whatever is necessary to validate the build manually.
