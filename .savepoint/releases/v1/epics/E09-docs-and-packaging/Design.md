---
type: epic-design
status: active
---

# Epic E09: docs-and-packaging

## Purpose

Prepare Savepoint as a public npm package with clear user-facing documentation, correct package metadata, and validated publish contents.

## What this epic adds

- Public README suitable for npm and GitHub.
- MIT license finalized.
- Package metadata review for repository, keywords, files, bin, engines, and exports.
- `npm pack` validation.
- Local install smoke tests from the packed tarball.
- Documentation for install, init, board, audit, doctor, and the audit loop.
- Release checklist draft for `0.1.0`.

## Components and files

Expected files introduced or extended by this epic:

| Path                             | Purpose                                            |
| -------------------------------- | -------------------------------------------------- |
| `README.md`                      | Public product and usage documentation.            |
| `LICENSE`                        | Final MIT license text.                            |
| `package.json`                   | Publish metadata and package allowlist.            |
| `.npmignore` or `files` field    | Publish contents control.                          |
| `docs/*.md`                      | Optional focused docs if README becomes too large. |
| `scripts/*.ts` or `scripts/*.js` | Optional pack/install smoke helpers.               |
| `test/packaging/**/*.test.ts`    | Packaging validation tests if scriptable.          |

## Architectural delta

Before this epic, the package can function locally. After this epic, the package is shaped for public consumption and can be validated as a packed npm artifact.

This epic does not publish the package; it makes publishing boring and inspectable.

## Boundaries

In scope:

- Public docs and package metadata.
- Packed artifact validation.
- Local install smoke tests.
- Ensure generated `dist/cli.js` is the published binary.

Out of scope:

- Actual npm publish.
- Manual agent matrix.
- Broad marketing site work.
- Major feature changes discovered during docs writing unless they block package correctness.

## Quality gates

- `npm pack --dry-run` or equivalent should show only intended files.
- Packed install smoke should verify `savepoint --help` and `savepoint --version`.
- README examples should match implemented command behavior.

## Design constraints

- Keep docs plain and practical for vibe coders.
- Do not over-document deferred features.
- Treat package contents as a release contract.
