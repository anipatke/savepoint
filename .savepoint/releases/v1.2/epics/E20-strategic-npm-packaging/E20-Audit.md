---
type: audit-findings
audited: 2026-05-14
---

# Audit Findings: E20 Strategic npm Packaging

## Main Findings

E20 is applied and closed. The root `savepoint` npm package remains the user-facing JS launcher, platform packages are generated as `savepoint-{os}-{arch}` with npm `os`/`cpu` metadata, binary validation covers PE, ELF, and Mach-O artifacts, packed install smoke testing exercises the installed package, and the publish workflow now verifies before publishing platform packages ahead of the root package.

Task acceptance criteria are satisfied for T001 through T004. The launcher resolves all six supported platform targets and reports unsupported or missing platform package cases clearly. `internal/buildtool` now owns npm manifest generation, target-aware binary validation, and packed smoke testing. The publish workflow has separate `verify` and `publish` stages, keeps npm auth scoped to publish steps, asserts every platform artifact exists, and publishes the root package last. README install/update guidance and the E20 release verification notes document the new package shape and safe recovery path.

Verification run during audit:

- `make build` passed.
- `go test ./internal/buildtool` passed.
- `npm.cmd test` passed all 13 launcher tests.
- `make pack-smoke` passed after allowing npm to use its normal cache/log directories outside the workspace.
- `make test` still fails at `TestBundledSavepointSkillsHaveDiscoveryFrontmatter` because `agent-skills/savepoint-audit/SKILL.md` is reported as missing YAML frontmatter. This matches the task handoff notes as a pre-existing, unrelated repository failure. It remains a release-quality residual risk until fixed outside E20.

The audit close update has been applied: `.savepoint/Design.md` now records the strategic npm packaging architecture and `last_audited` is set to `v1.2/E20-strategic-npm-packaging`. No E20 product-code changes were required during audit apply.

## Code Style Review

- [x] One job per file
- [x] One job per function
- [x] Test branches
- [x] Types document intent
- [x] Build only what is needed
- [x] Handle errors at boundaries
- [x] One source of truth
- [x] Comments explain WHY
- [x] Content in data files
- [x] Small diffs

The implementation is scoped and cohesive: JS runtime resolution is split into launcher, platform mapping, and binary location; Go buildtool additions isolate npm package generation, binary-format validation, and packed smoke testing; tests cover supported, unsupported, malformed, missing, and platform-specific branches. No speculative abstractions or duplicated package target lists beyond the current Go/JS boundary were introduced.

## Proposed Changes

No product-code replacement blocks were proposed for E20. The normal audit-close documentation update was applied to `.savepoint/Design.md` so the architecture record reflects root-plus-platform npm packaging, binary validation, packed smoke tests, and publish ordering.
