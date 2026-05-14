---
type: epic-design
status: audited
---

# E20: Strategic npm Packaging

## Purpose

Replace the tactical single-package binary bundle with a durable npm packaging model that installs the correct Savepoint executable for each supported platform and catches artifact mistakes before publish.

## What this epic adds

- A small root npm package that keeps the stable `savepoint` CLI entrypoint.
- Platform-specific optional npm packages for Windows, Linux, and macOS on amd64/x64 and arm64.
- Binary format validation for PE, ELF, and Mach-O artifacts before publish.
- Packed-install smoke tests that exercise the npm artifact rather than only the source tree.
- CI/publish workflow updates that build, verify, and publish the package set deliberately.

## Components and files

| Module | Purpose |
|--------|---------|
| `package.json` | Root npm package manifest, bin entrypoint, scripts, files, and optional dependency wiring |
| `bin/savepoint.js` | Runtime platform resolver and CLI launcher |
| `internal/buildtool` | Cross-platform binary build, npm artifact layout, validation, and tests |
| `.github/workflows/publish.yml` | Publish-time build, packed smoke tests, and npm publish ordering |
| `README.md` | User-facing install and upgrade guidance if the package shape changes commands or expectations |

## Architectural delta

Before this epic, the root npm package owns all binary delivery, making artifact correctness depend on one broad package layout. After this epic, the root package remains the user-facing CLI wrapper while platform packages own native binaries. The build tool becomes the source of truth for package generation and binary validation, and CI proves the packed npm artifacts install and execute on supported operating systems.

## Boundaries

**In scope:**
- Root wrapper package, platform package manifests, buildtool generation, binary validation, packed smoke tests, publish workflow updates, and minimal install docs.

**Out of scope:**
- Changing Savepoint CLI commands or runtime behavior.
- Changing Go release archives outside what is needed to share build validation.
- Dropping support for any currently listed target without explicit user approval.

## Quality gates

- Root package can resolve and launch the correct optional platform package.
- Windows artifacts validate as PE, Linux artifacts validate as ELF, and macOS artifacts validate as Mach-O.
- Packed install smoke tests run against npm tarballs, not only the working tree.
- Publish workflow cannot publish the root package before required platform artifacts are built and verified.
- `make build && make test` passes, apart from any explicitly documented pre-existing unrelated failures.

## Open decisions

None.
