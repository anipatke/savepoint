---
id: E01-scaffolding/T001-package-baseline
status: done
objective: "Establish the `savepoint` npm package identity and baseline repository hygiene files without implementing CLI behavior."
depends_on: []
---

# Task T001: Package Baseline

## Implementation Plan

- [x] Draft `package.json` with package name, binary mapping, Node engine, and module type.
- [x] Add `.gitignore` for generated, dependency, coverage, and local environment artifacts.
- [x] Add minimal `README.md` content that describes the repo as the `savepoint` scaffold.
- [x] Add MIT `LICENSE` text for the package baseline.
- [x] Verify the files align with the package and binary name `savepoint`.

## Scope

Create the root package and repository metadata that every later scaffolding task builds on.

Expected files:

- `package.json`
- `.gitignore`
- `README.md`
- `LICENSE`

## Acceptance Criteria

- `package.json` uses package name `savepoint`.
- `package.json` declares binary name `savepoint` pointing at `dist/cli.js`.
- `package.json` declares Node `>=20.10`.
- `package.json` is ESM-only with `"type": "module"`.
- `README.md` is intentionally minimal and development-focused.
- `LICENSE` contains MIT license text.
- `.gitignore` excludes generated build, dependency, coverage, and local environment artifacts.

## Out Of Scope

- TypeScript configuration.
- Build tooling.
- Test tooling.
- Lint or format tooling.
- Real command behavior.
