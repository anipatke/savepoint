---
id: C-935
scope: {kind: task, id: T-049}
result: CLEAR
checked_by: {role: checker, session: t049-task-check-20260927}
executed_session: planning-2026-09-26-o026
checked_at: '2026-09-26T22:35:00Z'
reviewed:
  head_commit: a700075f250088ea7ab5dbe61588fed2b15cdf8f
  files:
    - .savepoint/objectives/O-026-code-health-provider-feasibility/Objective.md
    - .savepoint/objectives/O-026-code-health-provider-feasibility/tasks/T-049-verify-the-six-health-inputs.md
    - .savepoint/releases/R-007-v2-1-code-health/Release.md
    - .savepoint/Guardrails.md
  dependencies: []
issues: []
supersedes: null
---

# C-935: T-049 Task Check — verify the six health inputs

## Independence and Scope

This is a fresh checker session, independent from the planning and execution session (`planning-2026-09-26-o026`) that researched and prepared T-049. Quick evidence mode is used for this requested Task Check ("check T049 as per router.md").

## Scope Lock

1. Criteria: T-049's six Done When items, User Check, and Boundaries (no scanner installation, production code, analysis algorithms, dashboard work, Check integration, custom-provider API, or generic plugin framework).
2. Changed files: `.savepoint/objectives/O-026-code-health-provider-feasibility/Objective.md` and `.savepoint/objectives/O-026-code-health-provider-feasibility/tasks/T-049-verify-the-six-health-inputs.md` (plus uncommitted router selection for T-049).
3. Relied-on context: R-007 Release decisions, Guardrails policy, `go.mod`.
4. Out of scope: provider execution engine, TUI rendering, Objective integration Check, and generic plugin ecosystems.
5. Materiality boundary: a material Issue must violate a Done When item or named Guardrail in the scoped files.

## Acceptance Classification

| Done When Item | Result | Evidence |
|---|---|---|
| 1. Candidate comparison (≥2 candidates where alternatives exist) and one official path selected per capability with explicit rationale and unsupported cases | Proven | All 6 capabilities compared in O-026: Tests (`go test -json`, Vitest JUnit, pytest JUnit vs `gotestsum`, Vitest JSON, `unittest`); Coverage (Go coverprofile, Vitest V8 JSON, `coverage.py` JSON vs Istanbul); Complexity (Lizard CSV vs `gocyclo`, ESLint `complexity`, Ruff C901); Duplication (jscpd v5 Rust CLI JSON vs `dupl`); Change hotspots (git-hotspots JSON vs Code Maat CSV); Vulnerabilities (OSV-Scanner v2 JSON vs `govulncheck`, `npm audit`, `pip-audit`). Unsupported cases explicitly documented for every capability. |
| 2. Eight-dimension comparison: stack support, deterministic machine output and schema stability, licence, platform and runtime requirements, offline behavior, bounded performance and report size, and project-owned installation | Proven | Stack coverage addresses Go, JS, TS, Python across all capabilities. Licences are permissive/standard (BSD-3-Clause, MIT, Apache-2.0; Code Maat's GPL-3.0 noted). Runtimes/platforms identified (including jscpd native binaries avoiding Node/Rust runtimes, Lizard Python 3.8+ requirement, and git-hotspots Zig 0.16 build / Windows gap). Offline behavior documented. Proposed 120s, 32 MiB, and 10,000-row bounds specified with truncation semantics. Project-owned installation model maintained. |
| 3. Candidate output contracts name fields for normalized values, scopes, failure states, freshness, provenance, and sanitized references; representative report or documented sample mapped | Proven | Common record fields defined: `capability`, `status` (`available`, `partial`, `absent`, `unsupported`, `failed`, `stale`), `value`, `unit`, `scope`, `collected_at`, provider/version, schema/version, target commit and working-tree fingerprint, configuration digest, sanitized references. Vulnerability advisory IDs and database timestamp mapped. Sanitized repo-relative paths required. Primary schema and sample links cited for all six paths. |
| 4. Exclusions for generated, vendored, and third-party code; fixture cases identified (valid, absent, malformed, partial, unsupported, stale, excluded-source) | Proven | O-026 specifies project-owned include/exclude configuration and `.gitignore` integration, excluding generated, vendored, and third-party code via globs/path lists rather than guessing. Fixture matrix explicitly enumerates valid populated reports, absent artifacts, malformed reports, empty test suites, partial/truncated reports, unsupported runners/lockfiles/platforms, stale targets/databases, and excluded files. |
| 5. Architecture boundaries: duplication avoids mandatory Savepoint runtime, complexity/hotspots remain external, no scanner installation or production adapter | Proven | jscpd 5.x Rust CLI ships standalone native binaries (macOS, Linux, Windows), avoiding Node or Rust runtime requirements for Savepoint. Complexity (Lizard) and hotspots (git-hotspots/Code Maat) remain external tools. No scanner was installed and no production code or adapter was added to Savepoint. |
| 6. Configured Task gate result and owner confirmation status recorded as evidence | Proven | `make build && make test-fast` passed cleanly (exit 0) and `git diff --check` is clean. The proposed catalogue is documented under "Proposed Provider Catalogue (T049, 2026-09-27)" with explicit "Owner confirmation pending" status. T-049 declares `owner_validation: {required: true}`, correctly separating technical verification from owner confirmation. |

## Commands and Gate Verification

- `git diff --check` — PASS (clean whitespace and formatting).
- `make build && make test-fast` — PASS (Go 1.26.2, linux/amd64; exit code 0; all unit and integration packages passed).
- `./savepoint doctor` — PASS (ALL CLEAN, exit code 0; no diagnostics or orphan records).
- `./savepoint resume` — PASS (strict index loading succeeds).

## Guardrails Review

- ARCH-03: Repo-relative paths with normalized separators mandated; no working-directory or daemon state dependency.
- ARCH-04: Code Health defined as a distinct internal module; no cross-boundary responsibilities added.
- CFG-02, CFG-03: Platform differences explicitly documented; git-hotspots Windows gap flagged as unsupported on Windows until packaged native builds exist.
- DEP-01: Zero new dependencies added to `go.mod`.
- TEST-01, TEST-02: Outcome evidence documented in O-026; comprehensive fixture cases specified for future testing.
- TEST-08: Gate passed via `make build && make test-fast`.

## Issues and Materiality

No in-scope Issues were found. No materiality actions are required.

## Observations (non-blocking)

1. **Change hotspots provider choice**: git-hotspots is an early alpha (`0.1.0-alpha.5`), requires source compilation with Zig 0.16, and has no verified Windows support. The alternative, Code Maat, is mature but requires Java and is GPL-3.0. The catalogue correctly flags this as an explicit owner decision for O-026.
2. **Resource ceilings**: The proposed 120-second deadline, 32-MiB report size, and 10,000-row caps are policy proposals that will require empirical verification with test fixtures in subsequent tasks (O-027 / O-028 / O-029).

## Owner Validation Still Needed

T-049 declares `owner_validation: {required: true}`. This Check provides technical clearance (**CLEAR**), but does not close T-049 or O-026. The owner must review the catalogue and record acceptance naming this Check (`C-935`) before T-049 can be marked `status: done`. The Full Objective Check remains mandatory before O-026 can close.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches**
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [x] STYLE-07 **One source of truth**
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs**
