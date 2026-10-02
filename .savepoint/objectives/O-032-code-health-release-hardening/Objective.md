---
id: O-032
title: Harden Code Health for the v2.1 release
status: done
depends_on: [O-026, O-027, O-028, O-029, O-030, O-031, O-035, O-036]
release: R-007
priority: medium
rank: 2
---

# O-032: Harden Code Health for the v2.1 release

## Outcome

Code Health is documented, migration-safe, privacy-conscious, cross-platform, performant, and independently verified as a modular foundation for future providers without weakening existing Savepoint behavior.

## Why

Persisted records, external processes, scaffolding, and a substantial TUI surface require explicit release evidence beyond feature-level tests.

## Success Conditions

- New-project scaffolding and existing-project upgrades preserve user content and maintain confirmed Design/config separation.
- Schema compatibility, strict loading, upgrade/downgrade diagnostics, concurrent access, partial writes, malformed files, and retention maintenance have regression coverage.
- Execution is tested for path safety, cancellation, timeout, process cleanup, bounded output, sanitisation, report limits, supported platforms, and network-independent normal tests.
- Packaging proves Savepoint remains one binary and neither bundles nor dynamically installs providers. Prerequisites and unsupported stacks are documented.
- Documentation explains purpose, limitations, capabilities, provider availability, states, thresholds, trends, refresh, Objective Checks, staleness, history, privacy, local-first behavior, and why Code Health is not a guarantee.
- The Full Objective Check covers all providers, partial/polyglot fixtures, history and staleness, Check integration, TUI behavior, existing behavior, Design reconciliation, and modularity.
- Release notes include a textual TUI walkthrough, provider rationale, packaging implications, supported and unsupported cases, limitations, and separate follow-up work.
- This repository’s failed vulnerability collection and broad duplication measurement have evidence-backed diagnoses and proposed repairs; uncertainty and owner decisions remain explicit. Known production complexity hotspots are refactored without behavior loss, targeting at most 20 CCN and aiming for 10 where justified. Runtime/test measurement scopes and thresholds are recorded, intentional template/history duplication is preserved, and any remaining work is separately proposed before implementation.

## Architectural Considerations

Hardening verifies the module boundary. Provider or platform incompatibility is represented honestly as unavailable rather than solved with bundled runtimes or analysis inside Savepoint.

## Boundaries

**In scope:** compatibility, migrations, packaging, privacy, security, performance, platforms, docs, regression evidence, Design reconciliation, investigation and proposed repairs for this repository's unhealthy or unmeasured health signals, and the mandatory Full Objective Check.

**Out of scope:** more providers, dead-code analysis, generic plugins, cloud services, CI/CD integrations, background monitoring, and unrelated features.

## Confirmed Design Decisions

The owner confirmed modularity, project-owned tools, local-first operation, single-binary distribution, explicit upgrades, bounded persisted output, explicit retention maintenance, and the seven-Objective structure on 2026-09-26.

## Confirmed Release-Hardening Design — 2026-10-02

The original 2026-09-26 decisions remain confirmed. The owner confirmed this scope and the health-analysis proposals on 2026-10-02. It narrows the release evidence and includes the subsequently completed O-035 popover and O-036 agent report. It does not approve implementation or add new providers.

### Proposed Observable Outcomes

1. **Adoption and compatibility:** temporary-project regressions exercise fresh initialization, existing schema-2 asset upgrades, legacy migration boundaries, and unsupported health-record versions. Confirmed health configuration, user-authored Design content, snapshots and edited guidance are preserved. Unknown versions fail with a named actionable diagnostic; no speculative schema conversion or downgrade rewrite is added.
2. **Stored evidence and execution safety:** exercise competing reads/writes, interrupted or failed replacement, malformed records, path/symlink refusal, output/report bounds, cancellation, timeout and process-tree cleanup. Immutable official evidence survives all maintenance and report failures; failures stay distinct from measurements. Cover the derived report and warning bridge as well as snapshots. Normal tests use controlled subprocesses and local fixtures rather than installed providers or network access.
3. **Performance and platform evidence:** establish repeatable benchmarks for dashboard load/history and rendering using explicitly sized fixtures, including hundreds and thousands of snapshots. Record time, allocations, fixture size and toolchain; identify growth attributable to history loading separately from tool runtime. Propose any numerical release budget from those measurements before treating it as a gate. Native Windows CI must pass, Linux full/distribution gates must pass, and all six shipped target archives retain the existing single-executable/checksum contract. A platform-sensitive repair requires a fresh full gate.
4. **Public guidance and release notes:** document all five signals and the existing nine provider report formats, prerequisites and unsupported cases, setup/check/report and popover workflows, thresholds/history/trends/staleness, manual versus official evidence, partial/polyglot fixtures, and Check integration. Include a textual TUI walkthrough and the agent-report handoff. Explain that no telemetry or provider installation is added, while an explicitly run project-owned provider may itself use the network and Savepoint inherits its execution environment. Do not promise offline vulnerability collection or provider sandboxing that the code does not enforce.
5. **Independent release integration:** the mandatory Full Objective Check covers the complete feature and hardening outcomes, provider fixture coverage, existing Savepoint behavior, code/module boundaries and Design reconciliation. Task evidence and optional Task-check waivers remain governed by the shared workflow. Publishing, tagging and deployment remain outside this Objective.

### Readiness Evidence

- O-026 through O-031, O-035 and O-036 are owner-completed; O-036 has current CLEAR C-956. Their implemented feature outcomes are available for the release boundary.
- `internal/codehealth/model.go` defines five capabilities, nine approved provider keys and bounded record fields. `runner.go` caps report stdout at 32 MiB and stderr at 64 KiB, runs argument vectors without a shell, and bounds cancellation pipe waits. Platform cancellation is owned by the Unix/Windows runner files.
- `internal/codehealth/storage.go` owns versioned config/snapshot reads, atomic persistence, and explicit PlanPrune/Prune. Existing retention tests protect official snapshots, deterministic ties, repeatability and malformed-history refusal. There is no documented public prune command in Design's CLI surface.
- `internal/init` owns managed scaffold/upgrade preservation; existing schema and preservation tests provide a baseline for health-specific adoption evidence.
- `internal/buildtool` validates exactly one executable per distribution archive and checksums. Existing CI runs `make ci` on Linux and the full test suite natively on Windows. Hardening should strengthen those existing mechanisms rather than create another release runner.
- README currently has no Code Health walkthrough. Project Design now records the implemented dashboard/report boundary and the owner's confirmed tools; exact tool commands stay in health config.

### Owner Decisions Recorded

The owner answered Yes to the proposed scope: release hardening plus diagnosis of the scanner failure, meaningful duplication measurement and targeted complexity fixes. Defer an owner-facing history cleanup command; verify and document existing internal retention behavior. Preserve thresholds and existing feature behavior. Performance remediation and unknown scanner/duplication repairs must be based on bounded research findings, with further planning if they require a material implementation decision. Design confirmation authorizes detailed planning, not execution or owner completion.

### Work Separation

Shape separate work around storage/compatibility regressions, process/platform safety, performance evidence, and public guidance/distribution verification. The allocated Tasks name exact contexts and dependencies. Separate worktrees are useful only where both changed and Context Files can be disjoint; documentation and final integration should follow verified contracts, and shared codehealth files must stay in a sequential lane. Unknown performance remediation becomes bounded research with a decision deliverable rather than a speculative optimization plan.

Detailed Tasks have been allocated by create-task and strict-loaded successfully. Router remains selected on O-032 in design until the owner reviews the concrete Task plan. Research recommendations requiring an unknown implementation are returned to design for bounded follow-up; the current plan does not pretend those repairs are already specified.

## Owner-Requested Health Analysis — 2026-10-02

The owner explicitly requested analysis and proposed fixes for the three concerning indicators while planning O-032. This authorizes investigation and proposals, not silent changes to health configuration, thresholds, production code or installed tools. The owner subsequently confirmed the release scope and proposed health work on 2026-10-02.

### Evidence Basis

Read the saved `.savepoint/health/report.md` and official snapshot `sha256:bcc7043ac9e301037c48d37851b740eaec8e79a9e6e2f8de32ef0f39f0cacbb6` (measured 2 Oct 13:11 local time), confirmed health configuration and targeted source functions. No tool collection or health command was run in design. These are saved measurements, not a fresh assessment of later code.

The snapshot has two Needs Attention labels and one Unknown, rather than three measured red findings. Complexity is 46 CCN; duplication is 9.7586%; dependency vulnerabilities are unmeasured because osv-scanner failed. Tests pass (3,205) and coverage is 86.2613%. All instances are currently advisory. Advisories still warrant the requested analysis; scanner failure provides no assurance either way about dependencies.

### Complexity: Structural Refactoring Proposal

**Observed:** `internal/doctor/repairs.go` SuggestRepair and V2ProblemRepair each score 46; `internal/doctor/checks.go` v2DiagnosticName scores 44; `internal/buildtool/main.go` runGoTestStream scores 31; `internal/board/v2/update.go` applyLoad scores 26. Many listed test functions also exceed the watch threshold, including two at 36. The current Good/watch limits are 10/20.

**Cause supported by source:** doctor functions enumerate large sets of diagnostic-to-message mappings in switches. SuggestRepair also has ordered string predicates, whose precedence matters. Test-stream processing and board reload orchestrate several separate responsibilities. A large switch can be more declarative than intrinsically risky, so preserve behavior instead of treating the score as proof of a bug.

**Proposed fixes:** move exact diagnostic-to-repair copy into typed lookup data with explicit fallbacks; represent ordered sentinel/predicate matching as an ordered rule list without changing precedence. Factor streaming decode, timing aggregation and failure-output reporting into focused helpers. Split board reload failure/reset from successful selection restoration only where existing state invariants remain explicit. For test hotspots, separate scenario setup and assertions or use coherent table cases, preserving independent assertions and failure readability; do not omit tests just to lower the number.

**Verification:** table coverage for every existing diagnostic and default, overlapping error/predicate precedence, mixed error chains, event stream malformed/read/write/process failures, and board first-load/retry/selection restoration. Compare before/after on identical scopes and fixtures with the same thresholds. Proposed target is removal of red-level production hotspots (at most 20 per function), aiming for 10 where a clear decomposition warrants it. A residual justified declarative/test hotspot gets an explicit disposition, not a hidden threshold increase. The owner confirmed these targeted refactoring boundaries and targets on 2026-10-02.

### Duplication: Scope Reconciliation Before Refactoring

**Observed:** 9.8% exceeds the current Good/watch limits of 3%/5%. The saved top evidence names archived V1 planning records, historical Check records, canonical skills mirrored into templates, and golden migration fixtures. TPL-01 explicitly requires skill/template byte identity, and archives are preserved evidence.

**Consequence:** the whole-repository percentage mixes maintained code with intentional/history duplication. The listed examples justify a scope investigation; they do not establish how much of the 9.8% is explained by each category or what runtime-only duplication would be. The original jscpd JSON report is not retained at its configured path, so an exact adjusted percentage cannot be calculated from the capped snapshot evidence alone.

**Proposed fixes:** agree a meaningful maintained-code measurement, excluding archived/planning records and intentionally mirrored generated assets from that metric. Keep production Go and test Go visible as separate scopes/instances where the provider can report reliable scoped totals. Preserve the canonical/template equality test and golden-fixture integrity as explicit release checks. Inventory all actual clone pairs under the agreed scope, then consolidate substantive runtime duplication into existing shared helpers and reusable test setup where that improves clarity. Never edit history, break mandated template identity, or weaken thresholds to manufacture a green score.

**Implementation constraint:** JscpdReader can recompute scoped totals when per-file counts exist; newer report formats may provide only global totals. Scoped execution must therefore configure the provider itself to scan the agreed targets/exclusions, and verify the resulting denominator. Merely filtering displayed evidence does not change a global percentage. Changing scope restarts comparability for that series; label old measurements honestly.

**Verification:** fixtures spanning runtime/test/archived/template content demonstrate numerator and denominator, scope/exclusion agreement and no missing clone evidence. Record whole-repo versus scoped measurements clearly and exact tool/report version. Proposed repair target is at or below 5% for the agreed maintained-code scope, aiming for 3%; confirm after baseline rather than promise an unknown result.

### Dependency Vulnerabilities: Bounded Diagnostic Proposal

**Observed:** snapshot outcome failed, with reason beginning `osv-scanner exited with status 127: Scanning dir . Starting filesystem walk for root: / Scanned ...go.mod...`. The persisted reason is truncated. No raw osv JSON report is retained. A project-owned executable exists at `/home/user/.local/bin/osv-scanner` and is an ELF binary, so absent executable is not supported as the diagnosis. The evidence does not distinguish a missing child tool, configuration/report incompatibility or a provider/environment/network failure.

**Proposed next step:** bounded research produces a named decision: reproduce the configured invocation in an owner-authorized diagnostic setting, capture exit status and sanitized complete diagnostics plus report existence, establish tool version and prerequisites, determine whether dependency queries require network access, and recommend the smallest setup/runner/reader correction. Do not install/upgrade tools or change the provider without explicit scope authorization. Normal regression tests use captured local fixtures and fake subprocesses; live collection remains an explicit owner refresh or Full Objective Check activity under the workflow.

**Desired outcome:** the existing provider yields valid measured evidence, or the limitation and owner decision are recorded plainly. Only a successful report can establish zero findings or reveal vulnerabilities. If findings appear, prioritize them by severity and reachability; do not invent dependency upgrades from this failed scan.

### Proposed Addition to Release Work

Treat diagnosis/scope agreement as prerequisites to targeted refactoring. The doctor/buildtool/board hotspots are distinct outcomes and should be split into bounded sequential or independent lanes according to their exact source/test contexts. Duplication measurement changes and scanner setup diagnosis must precede any claim that the repository's health improved. Preserve passing tests and current coverage while comparing measurements, and finish with the existing mandatory Full Objective Check.

## Planning Handoff Evidence

2026-10-02: Owner design confirmation recorded. All eight prerequisite Objectives are done; existing runtime/storage/runner/reader/upgrade/build interfaces are known. Unknown scanner cause, maintained-code duplication baseline and performance budget have bounded research decision deliverables. No new global policy or speculative runtime interface is introduced. Exact Task contexts are verified present, except CHANGELOG.md is explicitly a new documentation file. Dependency evidence paths are named in downstream Context Files. IDs were assigned exclusively by savepoint create-task, with strict full-index loading after creation. All Tasks are planned; no production source was edited by planning.

Parallel execution can use independent doctor, test-stream, board-reload and adoption lanes, whose implementation and Context Files are disjoint. Their common owning Objective/policy is loaded by the active skill, not duplicated into each Task's Context Files. The Code Health investigation/storage/process lane follows its explicit dependencies and shared interfaces. Performance and integration follow their prerequisites; distribution and documentation wait for the verified work they describe. Lane sessions follow AGENTS: no router or identity creation, commit locally, and merge before main-branch Checks.

Plan review remains an owner decision under savepoint-design step 5. On approval, select the first unblocked scanner diagnosis Task and set router state task, preserving the Goal bytes. The mandatory Full Objective Check is performed later by an independent checker and is not an executor Task.
