---
id: O-026
title: Select trustworthy Code Health providers
status: done
depends_on: []
release: R-007
priority: critical
rank: 1
---

# O-026: Select trustworthy Code Health providers

## Outcome

Savepoint has a confirmed modular boundary and an evidence-backed catalogue selecting one official project-owned tool or standard report path for each of the five Code Health capabilities in scope for V2.1: tests, coverage, complexity, duplication, and dependency vulnerabilities.

## Why

Provider choices must not smuggle language engines, mandatory runtimes, unsafe downloads, or fragile parsing into Savepoint.

## Success Conditions

- Candidates for tests, coverage, complexity, duplication, and dependency vulnerabilities are compared for supported stacks, deterministic machine output, licence, runtime, platforms, offline behavior, performance, report size, and installation ownership.
- One official provider or report path is selected per capability with limitations and fixture strategy.
- Duplication support does not make Node or another ecosystem a mandatory Savepoint runtime.
- Complexity remains external rather than becoming a Savepoint analysis engine. Change-hotspot collection and hotspot-derived signals are out of scope for V2.1.
- The narrow discovery, execution, normalization, snapshot, and dashboard interfaces owned by the separate Code Health module are settled.
- Compatibility rules identify which provider and configuration changes begin a new comparison series.

## Architectural Considerations

Code Health is compiled into Savepoint but owns its domain. Project-owned tools are invoked through structured executable and argument definitions without a shell. Scaffolding chooses from the supported catalogue and requires confirmation.

## Boundaries

**In scope:** bounded provider research, decisions, output contracts, fixtures, modular interfaces, packaging and licence implications, and unsupported cases.

**Out of scope:** change-hotspot collection or derived hotspot signals in V2.1; installing scanners, production adapters, analysis algorithms, dashboard work, Check integration, custom-provider APIs, and a generic plugin framework.

## Confirmed Design Decisions

The owner confirmed the project-owned-tool model, one official path per capability, a supported catalogue, no automatic installation, no AI assessment, and no compromise to Savepoint's single-binary distribution on 2026-09-26.

The owner confirmed Go, JavaScript, TypeScript, and Python as the initial stacks for provider feasibility on 2026-09-26. A supported path may be unavailable for a stack when the evidence does not justify one; the catalogue must state that limitation explicitly.

On 2026-09-27, the owner de-scoped change-hotspot measurement and hotspot-derived attention signals from V2.1. The V2.1 catalogue and module contract cover the other five capabilities; no hotspot provider, collection path, stored metric, or health-score input is in this release's scope.

On 2026-09-27, the owner confirmed the V2.1 measure set as tests, test coverage, complexity, duplication, and dependency vulnerabilities, and directed work to proceed to T-050. Dead-code measurement is not part of the V2.1 measure set.

On 2026-09-29, the owner approved the five official provider/report paths below and the proposed Code Health module boundary and compatibility rules: Go test JSON plus Vitest/pytest JUnit for tests, Go cover profiles plus Vitest V8/coverage.py JSON for coverage, Lizard for complexity, the jscpd v5 native CLI for duplication, and OSV-Scanner v2 for dependency vulnerabilities. The documented limitations, project-owned installation model, unsupported cases, and provisional resource caps remain part of that approval.

## Selected Provider Catalogue (T049, approved 2026-09-29)

**Owner-confirmed official paths for the five in-scope capabilities.** Tools are project-owned: Savepoint does not download or install them, run test suites, bundle them, or make their runtimes mandatory. The report paths below are the supported catalogue selections; alternatives remain comparison evidence. The owner has excluded change hotspots and dead-code measurement from V2.1.

### Tests

- **Go:** consume the existing `go test -json` event stream. Map `Action`, `Package`, `Test`, `Elapsed`, and `Output`; distinguish `pass`, `fail`, `skip`, and build errors. Keep only bounded, sanitized failure summaries from `Output`.
- **JavaScript and TypeScript:** consume Vitest JUnit XML when Vitest is already the project-owned runner. Map `testsuites/tests/failures/errors`, `testsuite/name/timestamp/tests/failures/errors/skipped/time`, and `testcase/classname/name/time` plus `failure`, `error`, and `skipped` children. Configure relative file names and suppress captured console output in the report; remove any host identity from existing reports.
- **Python:** consume pytest JUnit XML when pytest is already the project-owned runner. Map the same JUnit counters and testcase identity/outcome fields; retain pytest's runner version and the report timestamp.
- Go's toolchain is BSD-3-Clause; Vitest and pytest are MIT. These paths add no runner dependency when the project already uses them. If a project chooses another runner, its license/runtime and report support must be reviewed before adding a path.
- Do not launch tests to create these inputs. Go's structured output is from the Go toolchain; Vitest and pytest require their existing user-managed runtimes/dependencies. Projects using another runner without a supported structured report are `unsupported`; an absent report is `absent`, never a passing empty result. Compare against `gotestsum` (an extra Go tool), Vitest JSON, and Python's text-only `unittest` output; the chosen paths avoid introducing a test runner or Savepoint-specific runner adapter.

### Coverage

- **Go:** consume a fresh `go test -coverprofile` file already produced by the project gate. Each profile row supplies source path, start/end line-column positions, statement count, and execution count; `mode: set|count|atomic` identifies the aggregation mode. Derive covered and total statements from rows; do not rerun tests.
- **JavaScript and TypeScript:** consume Vitest's JSON coverage report using its V8 provider when the existing runtime exposes V8 and the project already owns `@vitest/coverage-v8`. Vitest 3.2+ remaps V8 coverage to source files. Normalize per-file statement/function/branch counters and coverage totals from the Istanbul-compatible report. Non-V8 runtimes and absent coverage configuration are explicitly unsupported/absent. Istanbul is the alternative provider when a project already owns it; it supports more runtimes but instruments source and has higher runtime and memory cost.
- **Python:** consume `coverage.py coverage json`; map `meta`, each entry under `files` and its `summary` (including covered/total statements and `percent_covered`), and `totals`. Respect the project's configured include/omit filters.
- Vitest and its coverage provider are MIT; coverage.py is Apache-2.0. The coverage path remains opt-in and uses only project-owned packages and runtimes.
- Only existing structured coverage artifacts are inputs. Projects with no report, unsupported runner/provider, or stale report do not receive a synthesized zero. Version-lock each selected producer and fixture its report contract.

### Complexity

- **Proposed provider:** Lizard, run by the project-owned Python environment with its CSV report. The documented report columns include function `NLOC`, cyclomatic complexity (`CCN`), token count, parameter count, function identity/line/file, and per-file `LOC`, average metrics, and function count. It lists Go, JavaScript, TypeScript/TSX, and Python support and also offers cppncss-style XML.
- **Comparison:** `gocyclo` is Go-only; ESLint's `complexity` rule and Ruff's C901 are language-specific paths with different configuration/report semantics. Lizard gives one normalizable measurement across the four stacks. It measures cyclomatic rather than cognitive complexity; its own docs describe a syntactic estimate, so values are comparable only within a pinned Lizard version/configuration. Unrecognized or partially parsed files must be marked unsupported rather than reported as complexity zero.
- Lizard is MIT-licensed and requires Python 3.8 or later. It is not bundled; the project owner supplies the executable/runtime. Pin its version and fixture the CSV columns and ordering.

### Duplication

- **Proposed provider:** jscpd 5.x Rust CLI, using its JSON reporter. Its documented report sample maps `duplicates[].format/lines/tokens/firstFile.name/start/end/secondFile.name/start/end` and `statistics.total.lines/tokens/sources/clones/duplicatedLines/percentage`.
- **Comparison:** `dupl` uses a Go AST and only analyzes Go; its own documentation cautions that it can report false positives. jscpd supports the four target languages, emits structured JSON, and v5 is a self-contained native binary, so Savepoint does not gain a Node or Rust runtime requirement. Its upstream packages publish binaries for macOS arm64/x64, Linux arm64/x64 (glibc/musl), and Windows arm64/x64. jscpd is MIT-licensed.
- Use the project's pinned configuration and explicit exclusions, `--workers 1`, and project-selected `--min-tokens`, `--min-lines`, `--max-size`, and `--max-lines`. These parameters change what counts as duplication and therefore begin a new comparison series when changed. Do not consume jscpd's separate health score, complexity estimate, dashboard, or dead-code output.

### Change hotspots — excluded from V2.1

- **Owner decision (2026-09-27):** Change hotspots are de-scoped for V2.1. The comparison below is retained as research evidence only; no candidate is selected and no hotspot result is collected, stored, trended, or rendered in this release.

| Candidate | Signal and report fit | Runtime, licence, and platform | Local use and bounds | Main trade-off |
|---|---|---|---|---|
| **git-hotspots** | Git-only, file-level history analysis, so the signal itself is language-independent. Its JSON report and score cover change frequency, additions/deletions, recency, co-change, confidence, evidence, and caveats — the closest match to the proposed T-049 normalized fields. | Apache-2.0. Git plus Zig 0.16 is needed to build the current `0.1.0-alpha.5`. There is no published general-purpose binary/package; Windows and non-Linux packaging remain unverified. | Reads local history without network calls. `--limit` and `--since` bound report rows or selected history. No comparable-repository performance result or hard input-byte/process-time ceiling is published. | Best semantic fit, with the weakest release maturity and platform packaging. The report schema is not declared stable; pin a source revision and fixture its output. |
| **Code Maat** | Language-independent VCS mining, including Git. `entity-churn`, age, and coupling are separate analyses with CSV outputs; the documented CLI does not provide one combined hotspot report or JSON schema. Adopting it would require treating those separate measurements as the contract, rather than silently inventing a Savepoint score. | GPL-3.0. Clojure standalone JAR, Java 8+ at runtime. Git logs are prepared locally first; the project guide advises Git Bash on Windows for its log line-ending expectations. Java makes it broadly portable where a compatible JVM is installed. | Local/offline processing. `--rows` caps output and date filters can limit Git history. T-049 has no comparable performance or input-size measurement for it. | Older, flexible VCS analysis, but weaker fit to the single structured hotspot contract and retains the GPL/Java costs. |
| **CodeLore** | Local Git history analysis with CSV, JSON, and other structured formats. Its `hotspots` report ranks revisions together with complexity/health; the documented score is not the same as git-hotspots' frequency/churn/recency/co-change score. It also has a separate coupling analysis. Files without supported complexity data are omitted from hotspot ranking, so source-language coverage affects the result. | GPL-3.0-only. Native Rust CLI; source install requires Rust 1.96+, with released binaries for macOS and Linux arm64/x64 and Windows x64. Windows arm64 is not in the published target list. It is pre-1.0 and under active development. | Runs locally. CLI `--rows` caps output. The separate MCP hotspot read defaults to 50 rows and caps at 500. Docs describe a cold history ingest around 5–30 seconds, not a benchmark for Savepoint projects. No comparable output-byte or deadline test has been run. | Better packaged distribution than git-hotspots and no Java runtime, but shares Code Maat's GPL concern and defines a materially different, complexity-dependent measure. |

**Historical assessment only:** git-hotspots most directly supplies the proposed hotspot signals. CodeLore and Code Maat are alternatives with different definitions and constraints. None of the three was installed or benchmarked for this comparison. This assessment does not select a V2.1 provider; hotspot work is out of scope for this release.

### Dependency vulnerabilities

- **Proposed provider:** OSV-Scanner v2 JSON output, one external CLI for Go, JavaScript/TypeScript, and Python. Prefer resolved lockfiles: `go.mod`; npm, pnpm, Yarn, or Bun lockfiles; and Python `uv`, Poetry, Pipenv, PDM, `pylock.toml`, or fully pinned requirements. Map `results[].source.path/type`, `packages[].package.name/version/ecosystem`, and `vulnerabilities[].id/aliases` plus affected/fixed range and severity data where present. Normalize absolute source paths to repository-relative references.
- **Comparison:** `govulncheck` offers Go-specific call reachability; `npm audit` and `pip-audit` are ecosystem-specific. OSV-Scanner is Apache-2.0 and publishes prebuilt platform binaries. Its normal mode can contact OSV/deps.dev services with package identity/version data; offline vulnerability matching is available only with a previously cached database. Capture scanner version and database snapshot time when exposed; if freshness cannot be established, label it `unknown`/`stale`, not clean. A manifest without a resolved version is best-effort and must be marked partial; no supported lockfile or ecosystem is unsupported.

### Shared normalized record and bounds

Every normalized capability record carries `capability`, `status`, `value`, `unit`, `scope`, `collected_at`, provider name/version, report/schema version where available, target commit and working-tree fingerprint, configuration digest, and sanitized evidence references. Status is one of `available`, `partial`, `absent`, `unsupported`, `failed`, or `stale`; missing or malformed data never becomes numeric zero. Vulnerability records also carry advisory IDs and the database snapshot timestamp/identity where available. Paths are repository-relative with normalized separators. Do not retain absolute host paths, hostnames, usernames, private remote URLs, source snippets, raw test output, or commit/author text.

Use only explicit project-owned include/exclude configuration plus `.gitignore` where the selected provider honors it. Exclude generated, vendored, and third-party sources through configured globs/path lists; do not guess based on directory names alone. Fixtures must cover: valid populated reports for each supported stack; absent artifact; malformed JSON/XML/CSV/Go profile; empty suite versus no tests; partial/truncated report; unsupported runner/runtime/lockfile/platform; stale target or vulnerability database; and generated, vendored, and third-party files that must be excluded.

No provider was installed or benchmarked for T049. Among the in-scope options, jscpd supports source size/line limits and worker count. This does not establish safe Savepoint process/output ceilings. Proposed Savepoint defaults for the next interface task are a 120-second provider deadline, 32 MiB maximum report input, and 10,000 normalized rows per capability. Reaching any cap produces `partial` with an explicit truncation reason, never a clean result. These are policy proposals, not measured performance claims; fixture and representative-repository measurements must validate or revise them before release.

### Primary evidence reviewed

- Go test JSON and coverage: [Go command documentation](https://go.dev/cmd/go/), [test2json event format](https://pkg.go.dev/cmd/test2json), [coverage profiles](https://go.dev/doc/build-cover), and [Go license](https://go.dev/LICENSE).
- Vitest: [JUnit reporters](https://vitest.dev/guide/reporters.html), [coverage providers](https://vitest.dev/guide/coverage.html), and [MIT license](https://github.com/vitest-dev/vitest/blob/main/LICENSE). Python: [pytest JUnit XML](https://docs.pytest.org/en/stable/how-to/output.html), [pytest MIT license](https://github.com/pytest-dev/pytest/blob/main/LICENSE), [coverage.py JSON](https://coverage.readthedocs.io/en/latest/commands/cmd_json.html), and [coverage.py Apache-2.0 license](https://github.com/coveragepy/coveragepy/blob/main/LICENSE.txt).
- Complexity: [Lizard supported languages, output and requirements](https://pypi.org/project/lizard/) and [its MIT license](https://github.com/terryyin/lizard/blob/master/LICENSE.txt); alternatives: [gocyclo](https://github.com/fzipp/gocyclo), [ESLint complexity](https://eslint.org/docs/latest/rules/complexity), and [Ruff C901](https://docs.astral.sh/ruff/settings/).
- Duplication: [jscpd v5 engine, platforms, filters and reporter](https://github.com/kucherenko/jscpd/blob/master/docs/rust.md), [JSON sample](https://jscpd.dev/reporters/json), and [dupl](https://github.com/mibk/dupl).
- Hotspots: [git-hotspots README and CLI](https://github.com/arsham/git-hotspots), [git-hotspots user guide](https://github.com/arsham/git-hotspots/blob/main/docs/user-guide.md), [Code Maat README and CLI](https://github.com/adamtornhill/code-maat), and [CodeLore README and license](https://github.com/emrecdr/codelore) plus its [advanced guide](https://github.com/emrecdr/codelore/blob/main/docs/advanced-usage.md).
- Vulnerabilities: [OSV-Scanner supported lockfiles](https://google.github.io/osv-scanner/supported-languages-and-lockfiles/), [JSON output](https://google.github.io/osv-scanner/output/), [offline matching](https://google.github.io/osv-scanner/usage/), and [OSV-Scanner repository/license](https://github.com/google/osv-scanner).

## Confirmed Code Health Module Boundary (T050, approved 2026-09-29)

**Status:** owner-approved boundary for the five confirmed V2.1 measures and selected provider paths above.

### Ownership and storage

- A new `internal/codehealth` package owns the fixed five-capability catalogue, scaffold-time discovery, structured external-tool execution, provider-specific report normalization, health snapshot persistence, and a narrow read model. It does not become a generic plugin or user-authored command framework.
- Store project-owned Code Health settings in `.savepoint/health/config.yml`; store normalized, immutable snapshots under `.savepoint/health/snapshots/`. O-027 defines their versioned schemas and identity rules before implementation.
- Configuration contains a catalogue capability/provider key, executable path, argument array, expected report/artifact path, timeout, measurement options, scope and exclusions, and project thresholds. Execute with an explicit executable and argument vector in the project root; never interpret a shell string or install a tool.
- Discovery reads project structure and existing tool/report configuration to return proposals for owner confirmation. It does not run a scanner, create a measurement, alter configuration, or call the network.

### Service contract and callers

| Caller | Request and result | Boundary |
| --- | --- | --- |
| Scaffolding | `Discover(project facts)` returns candidate capability/provider pairs, detected inputs, support gaps, and a plain-language reason for each proposal. | Read-only; no execution or writes. The owner confirms before a proposal becomes configuration. |
| Explicit refresh or Full Objective Check | `Collect(mode, configured capabilities, repository identity)` returns one collection outcome per capability plus an immutable snapshot identity. | Run configured providers sequentially. Consume fresh test/coverage artifacts rather than rerunning tests; execute only selected project-owned analysis/scanner tools. Capture timeout, cancellation, exit status, bounded output, and collection time. |
| Check workflow | The Full Objective Check receives the snapshot identity and records that reference on its Check; it does not embed report data. | `CheckV2` currently has no Code Health snapshot reference. O-027 must add a typed reference without overloading `ReviewedBasis` or `Issues`. Collection remains an explicit Check/refresh operation, never a side effect of decoding a Check. |
| TUI | A read-only `Summary` and `History` query returns bounded current values, status, explanations, and comparable history points. | `internal/board/v2` receives the read model through its command/load boundary. Rendering performs no filesystem or process IO; opening the board never scans. |
| Doctor | No Code Health collection call. | Existing schema, planning-record, and configured-quality-gate diagnostics remain separate from health providers. |

The collector and reader are separate interfaces. The collector owns execution and writes snapshots; consumers receive immutable summaries and snapshot references only. The proposed resource limits are 120 seconds per provider, 32 MiB of report input, and 10,000 normalized rows per capability. These remain unmeasured caps from T049 and require validation before release.

### Result and error states

The service must preserve the distinction among: `not_configured` (no confirmed provider configuration), `unavailable` (configured executable or input cannot be accessed), `unsupported` (no supported path for the detected stack or report), `absent` (expected fresh artifact is missing), `available` (valid evidence, including a measured zero), `partial` (valid but incomplete or capped), `failed` (execution or parsing error), `timed_out`, `cancelled`, and `stale` (evidence does not match the current repository/configuration or freshness basis). A failure or missing report is never represented as a clean zero. O-027 must settle the exact versioned status model and any additional `unknown` freshness case.

### Series compatibility

- A comparison series is keyed by capability, provider identity and version, report/schema version, measurement-definition version, effective scope/exclusion fingerprint, and relevant provider configuration. A change to any of these starts a new series; values from old series remain readable but are not joined into a trend.
- Each snapshot also records target commit, relevant working-tree fingerprint, collection time, and sanitized provider provenance. These describe the point and determine freshness; the target commit changing between observations is expected and does not itself reset a series.
- Record the vulnerability database identity/time when exposed. If the database identity changes or cannot be established, do not silently compare scanner results as a code-only trend; mark comparability unknown or start a separate series until O-027 defines the exact rule.
- Manual snapshots stay separate from official Full Objective Check baselines. Full Check snapshots are permanent; manual retention and pruning follow the owner-confirmed ten-snapshot policy, with pruning explicit and outside collection.

### O-027 handoff

O-027 currently describes six signals; before it starts, its scope must be reconciled to the owner's five-measure V2.1 decision. O-027 must define snapshot and config schemas, stable snapshot IDs, status serialization, Check references, repository/dirty-tree fingerprints, manual versus official series, retention, stale-data rules, classification/trend behavior, vulnerability-database comparability, and fixtures for valid, missing, failed, partial, unsupported, cancelled, timed-out, stale, and incompatible results. No O-027 implementation or schema is included in T050.
