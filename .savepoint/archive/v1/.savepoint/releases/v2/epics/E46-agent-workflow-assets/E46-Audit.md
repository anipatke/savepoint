---
type: audit-findings
audited: 2026-09-19
---

# Audit Findings: E46 Guide planning, execution, and independent checking

## Main Findings

### Verdict

CLEAR. Both initial findings are closed, all E46 acceptance coverage is proven at the repository boundary, and no owner waiver or post-push evidence is required. The epic is **CLEAR TO COMMIT/PUSH**.

### Closure Map

- **Required provenance is written but not decoded — closed.** `TaskV2` now strictly decodes required planner `planned_by` provenance. `CheckV2` now strictly decodes required `executed_session` provenance and rejects a Check whose execution and checking session are the same. Valid, missing, malformed, blank, and wrong-authority cases have direct decoder coverage, and doctor fixtures exercise the stricter records through project loading.
- **The design handoff names the wrong router state — closed.** Both copies of `savepoint-design` now set `state: task` and hand execution to `savepoint-task`. A cross-contract test compares this outgoing transition with both V2 routing tables and rejects the stale V1 `task-building` handoff.

### Materiality Summary

No materiality actions are required.

### What Is Proven / Not Proven

Proven: the four V2 skills and three references exist with canonical/shipped byte parity; sensitive role/write ownership remains partitioned; Task and Check templates emit provenance that their typed decoders retain and validate; fresh checking cannot claim the execution session as the checking session; the design handoff stays inside the documented four-state V2 router; the configured build gate runs between typecheck and test; and all focused and repository gates pass.

Not proven by this epic: semantic title quality and the complete owner-facing V2 runtime flow. Those are explicitly assigned to E50 agent evaluation and the later cutover epics, rather than claimed by E46's text and schema contracts.

### Audit Evidence

- Scope lock: unchanged from the initial audit — T001-T009 acceptance criteria; four public V2 skills; three shared references; shipped mirrors; routing guides; typed Task/Check contracts; init contract tests; doctor V2 project loading; and T008's configured build-gate path.
- Coverage and workflow result: canonical/shipped parity, role/write ownership, lifecycle/routing boundaries, artifact-schema fidelity, malformed and bypass cases, project-level loading, and quality-gate ordering/failure behavior all pass. The two initially failed cells now pass.
- File reality and drift: all scoped files exist, mirrored assets are byte-identical, and E46's implemented workflow/configuration delta is reconciled into `.savepoint/Design.md`. No unexplained phantom file remains.
- Gates: `git diff --check` passed; focused `internal/data`, `internal/init`, and `internal/doctor` tests passed; `make build && make test` passed across all packages. `.savepoint/Health-Check.md` remains absent, so its optional procedure step is not applicable.

### Guardrails Verification

- Rule IDs checked: DATA-01, DATA-02, DATA-03, TPL-01, TPL-02, TPL-03, TPL-04, ARCH-03, ARCH-04, CFG-01, TEST-01 through TEST-08, POL-01, POL-02.
- Health check mode: Full.
- Evidence: strict provenance decoder tests, project-level doctor tests, cross-contract routing tests, byte-parity tests, `git diff --check`, and the complete build/test gate.
- File reality evidence: canonical skills/references, shipped mirrors, typed records, tests, config, task records, and Design reconciliation are present.
- Waivers or unresolved findings: none.

### Non-Blocking Observations

The router changed concurrently to active defect `D003-ci-workflow-missing-token-permissions`. Closing E46 intentionally preserves that defect workflow instead of overwriting it with E47 routing; advance to E47 after D003 is resolved.

## Code Style Review

- [x] STYLE-01 **One job per file** — split files when responsibilities mix.
- [x] STYLE-02 **One job per function** — small, named, testable units.
- [x] STYLE-03 **Test branches** — cover meaningful conditionals and edge cases.
- [x] STYLE-04 **Types document intent** — prefer explicit types over comments.
- [x] STYLE-05 **Build only what is needed** — no speculative abstractions.
- [x] STYLE-06 **Handle errors at boundaries** — validate inputs, APIs, IO, and external data.
- [x] STYLE-07 **One source of truth** — no duplicated rules, constants, state, or config.
- [x] STYLE-08 **Comments explain why** — not what the code already says.
- [x] STYLE-09 **Content lives in data** — keep copy/config out of logic.
- [x] STYLE-10 **Small diffs** — minimal, reviewable, behaviour-preserving changes.
