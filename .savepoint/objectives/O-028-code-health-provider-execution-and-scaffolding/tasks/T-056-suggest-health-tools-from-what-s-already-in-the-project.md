---
id: T-056
title: Suggest health tools from what's already in the project
objective: O-028
status: in_progress
stage: audit
depends_on: [{task: T-055, requires: clear}]
owner_validation: {required: false}
planned_by: {role: planner, session: planning-o028-20261001}
---

# Suggest health tools from what's already in the project

## Outcome

A read-only `Discover` service inspects a bounded set of project files and PATH and returns confirmation-ready proposals: one per capability and component, each with a plain-language reason, the exact configuration it would write, default exclusions, and any support gap. It never executes, installs, writes, or uses the network.

## User Check

None beyond the Full Objective Check; the owner sees the proposals through `savepoint health setup` in the setup Task.

## Done When

- `Discover(ctx, root, lookPath)` returns a deterministic `Proposal` list. Each proposal names capability, provider, instance name, scope, exclusions, executable (for executed tools), args, report path, reason, and a gap (`missing_executable`, `unsupported_stack`, `no_report_yet`, or none). Every proposed `CapabilityConfig` passes `Validate()`.
- **Detected inputs:** `go.mod`; `package.json` scripts and dev dependencies (Vitest, `@vitest/coverage-v8`), `vitest.config.*`; `pyproject.toml`, `pytest.ini`, `setup.cfg`, `.coveragerc`; npm/pnpm/Yarn/Bun, Go, and Python lockfiles named in the O-026 catalogue; `.jscpd.json`, `lizard` config; existing reports at the catalogue's conventional paths. Unrecognized stacks produce an explicit `unsupported_stack` gap, never a silent omission.
- **Executed tools** (Lizard, jscpd, OSV-Scanner): proposal args include the report output and the confirmed exclusions translated into that tool's own exclude flags; OSV-Scanner uses its normal online mode, and its proposal reason says the scanner sends package names and versions to OSV.dev. A tool not on PATH is still proposed, with `missing_executable`, and never downloaded.
- **Report-only providers** (tests and coverage): proposals point at the report path the project's gate should produce and include the exact flag the owner would add to their own gate (for example `go test -json ... > <path>`); Savepoint does not edit their scripts.
- **Default exclusions:** one catalogue list (`vendor/**`, `node_modules/**`, `third_party/**`, `dist/**`, `build/**`, `**/*.pb.go`, `**/*_generated.*`, `**/*.min.js`) proposed only when the path exists or the pattern is a generated-file suffix; the list lives in data, not scattered in logic.
- **Monorepos:** manifests are searched to a fixed depth (3) skipping excluded and hidden directories; each component yields its own named, scoped instance. No aggregate instance is proposed.
- **Bounds:** at most a fixed number of files read, each capped in size; a larger or unreadable file becomes a named gap, not a crash. Symlinks outside the project are not followed.
- Tests prove no process is started (a `lookPath` fake is the only outside call) and no file is written (tree bytes and mtimes unchanged).

## Context Files

`.savepoint/objectives/O-028-code-health-provider-execution-and-scaffolding/Objective.md`; `.savepoint/objectives/O-026-code-health-provider-feasibility/Objective.md` (catalogue sections only); `internal/codehealth/model.go`; `internal/codehealth/config.go`; `internal/codehealth/primitives.go`; new `internal/codehealth/discovery.go`, `internal/codehealth/discovery_catalogue.go`, `internal/codehealth/discovery_test.go`, and fixtures under `internal/codehealth/testdata/discovery/`.

## Design References

Design sections 1 and 11; O-026 Selected Provider Catalogue and Module Boundary; O-028 Confirmed Design Decisions (2026-10-01).

## Guardrails

FS-03, FS-05, ARCH-03, ARCH-04, CFG-02, CFG-03, DEP-01, DEP-02, STYLE-09, TEST-01, TEST-02, TEST-04, TEST-08.

## Implementation Plan

1. Confirm the first Task's `Name`/`Required` fields and default-timeout lookup exist; return REPLAN REQUIRED if not.
2. Write the catalogue data (detected files, report conventions, per-tool argument templates and exclusion-flag translation, default exclusions) in `discovery_catalogue.go`.
3. Implement bounded component search and per-stack detection in `discovery.go`; parse only the fields needed (`package.json` via `encoding/json`; TOML/INI files by bounded line scan for known keys; no new dependency).
4. Build proposals, validate each as a `CapabilityConfig`, and sort deterministically (component path, capability order, provider).
5. Fixtures: Go-only, Vitest with V8 coverage, pytest with coverage.py, a polyglot monorepo with two components, a repo with vendored and generated code, an unsupported stack, a missing executable, an oversized manifest, and an escaping symlink.

## Boundaries

No execution, no writes, no network, no install hints beyond naming the missing executable, no CLI, no Design.md edits, no report parsing.

## Technical Verification

Focused runs while iterating; `make build && make test-fast` at handoff. Discovery is filesystem- and platform-sensitive, so also run fresh `make test-full` and note Windows path handling.

## Technical Evidence

Executed 2026-10-01. Handoff is ready for the owner's choice of an optional Task Check or a waiver; this is not a pass or clearance.

**Per-criterion outcomes**

- `Discover(ctx, root, lookPath)` deterministic, every config valid: `discovery.go`; `TestDiscoverIsDeterministic` (5 reruns equal) and the `discover` test helper validating every proposal. Proposals carry capability, provider, name, scope, exclusions, executable, args, report, reason and gap through `Proposal` and `CapabilityConfig`.
- Detected inputs: Go (`TestDiscoverGoOnly`), Vitest + V8 and missing V8 (`TestDiscoverVitestWith*`), pytest + coverage.py (`TestDiscoverPytestWithCoveragePy`), Python without pytest, JS without Vitest, unsupported stack (`TestDiscoverUnsupportedStack`), lockfiles including pnpm, existing reports (`TestDiscoverExistingReportClearsGap`). Limitation: `package.json` scripts are matched by the substring "vitest"; TOML/INI are line scans, not parsers.
- Executed tools: `TestDiscoverExecutedToolArgsTranslateExclusions` checks report output plus exclusion translation for all three tools; the OSV-Scanner reason names OSV.dev (`TestDiscoverGoOnly`); `TestDiscoverMissingExecutableIsProposedNotInstalled`.
- Report-only providers: gate flag and report path asserted (`TestDiscoverGoOnly`, monorepo test with `cd <dir> &&` prefix). Savepoint edits no scripts.
- Default exclusions come from the `defaultExclusions` table; `TestDiscoverDefaultExclusions` shows directory patterns only when the directory exists and suffix patterns always.
- Monorepos: `TestDiscoverMonorepoNamesAndScopes` (distinct names and scopes, no root aggregate, `node_modules` unsearched, combined `Config.Validate()` passes). `TestDiscoverSearchDepthIsBounded` for depth 3.
- Bounds: oversized, malformed, escaping-symlink, read-limit and unrepresentable-name cases each become a named gap (`TestDiscoverOversizedManifestIsAGap`, `...MalformedManifest...`, `...EscapingSymlink`, `...ReadLimitIsAGap`, `...UnrepresentableNameBecomesNote`).
- No process or write: `TestDiscoverIsReadOnly` (tree bytes, sizes and mtimes unchanged; `lookPath` only asked about `lizard`, `jscpd`, `osv-scanner`). The package imports no `os/exec` or `net`.

**Commands run:** `make build && make test-fast` (pass); fresh `make test-full` (exit 0, includes linux/darwin/windows cross-builds); `GOOS=windows go vet ./internal/codehealth/` (pass). Focused `go test ./internal/codehealth/ -run 'Discover|Lizard'` while iterating.

**Files read:** `AGENTS.md`, `agent-skills/savepoint-task/SKILL.md`, the Task, O-028 and O-026 (catalogue sections), `model.go`, `config.go`, `primitives.go`. Extra reads: `.savepoint/router.md` and the T-055 frontmatter (to confirm the selection and the dependency), `errors.go` and a grep of `Guardrails.md` (error helpers and rule text), `go.mod` version line. WebFetch of OSV-Scanner docs to check flags (see limitations).

**Files changed:** new `internal/codehealth/discovery.go`, `discovery_catalogue.go`, `discovery_test.go`; `AGENTS.md` Codebase Map entry (drift note).

**Limitations**

- Windows: only `go vet` for `GOOS=windows` and the cross-build ran here; the tests were not run natively on Windows. Paths use `filepath` for disk access and slash form in records. The symlink test skips when symlinks cannot be created.
- OSV-Scanner `scan source -r --format json --output-file` matches its docs. The exclude flag `--experimental-exclude g:<glob>` could not be confirmed from the docs and is from memory; T-057 or the Check must verify it. Lizard `-x`/`-o` and jscpd `--ignore`/`--output`/`--workers` come from memory of the tools' CLIs and are not verified against installed versions.
- Fixtures are built in temporary directories inside the tests, not committed under `testdata/discovery/`, as Test guardrail TEST-04 prefers; the plan listed `testdata`.
- A component nested inside another (root plus sub-packages) gets its own scope while the root keeps an empty scope, so their scopes can overlap in practice; not proposed away here.
- `.lizardrc` is a guess at a Lizard config name and only affects the reason text.
- Added gaps beyond the four named: `unreadable_file`, `oversized_file`, `read_limit_reached`.

## Drift Notes

Update the AGENTS.md Codebase Map entry for `internal/codehealth`: it currently says the package "performs no collection"; discovery is now in scope. Reconcile through the Full Objective Check.
