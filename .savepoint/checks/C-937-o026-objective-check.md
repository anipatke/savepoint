---
id: C-937
scope: {kind: objective, id: O-026}
result: NEEDS WORK
checked_by: {role: checker, session: o026-objective-check-20260929}
executed_session: planning-2026-09-26-o026
checked_at: '2026-09-29T10:49:37Z'
reviewed:
  base_commit: a700075f250088ea7ab5dbe61588fed2b15cdf8f
  head_commit: a700075f250088ea7ab5dbe61588fed2b15cdf8f
  files:
    - .savepoint/objectives/O-026-code-health-provider-feasibility/Objective.md
    - .savepoint/objectives/O-026-code-health-provider-feasibility/tasks/T-049-verify-the-six-health-inputs.md
    - .savepoint/objectives/O-026-code-health-provider-feasibility/tasks/T-050-draw-the-health-module-boundary.md
    - .savepoint/releases/R-007-v2-1-code-health/Release.md
    - .savepoint/Design.md
    - .savepoint/Guardrails.md
    - .savepoint/config.yml
    - go.mod
    - internal/data/check_v2.go
    - internal/data/objective_gate_v2.go
    - internal/board/v2/load.go
    - internal/doctor/v2_runtime.go
  dependencies:
    - C-935
    - C-936
issues: [I-084]
supersedes: null
---

# C-937: O-026 Full Objective Check

## Independence and Evidence Mode

This is a fresh checker session independent from the planning/execution session recorded on O-026's Tasks. Full evidence mode covers both owned Tasks, their integration, the Objective outcome and success conditions, R-007, applicable Guardrails, and reconciliation with Design.

## Frozen Scope Lock

1. Acceptance and release gates: O-026's outcome and six success conditions; T-049's six Done When items; T-050's seven Done When items; R-007's provider-selection, modular-boundary, local-first, Windows, and single-binary conditions; TEST-08's current full-gate requirement.
2. Changed records and public design surfaces: O-026, T-049, T-050, R-007, and their proposed provider catalogue, module requests/results, status model, snapshot boundary, and series-compatibility rules. C-935 and C-936 are prior claims to verify, not substitutes for this Check.
3. Relied-on boundaries: current Check model/gate, board load boundary, doctor runtime boundary, project configuration, Go dependency set, and Design sections 1, 3, 8, 11, 12, and 13.
4. Matrix axes: five capabilities; normal, absent, malformed/failed, partial/capped, unsupported, stale, timed-out, and cancelled states where applicable; configured and unconfigured providers; provider/config/scope/report-definition changes; manual and Full-Check snapshots; board, Check, scaffolding, refresh, and doctor consumers. Production execution and rendering are not applicable because O-026 is explicitly planning-only.
5. Materiality boundary: an Issue must violate an O-026/T-049/T-050 acceptance condition, named Guardrail, or R-007 release gate through the supported planning-record path. Later implementation details and unsupported custom providers are out of scope.

## Mandatory Coverage Matrix

| Surface / invariant | Normal | Missing / failure | Boundary / transition | Result |
| --- | --- | --- | --- | --- |
| Provider catalogue across tests, coverage, complexity, duplication, vulnerabilities | Candidate paths and alternatives are documented | Unsupported stacks and absent artifacts stay explicit | Runtime, platform, licence, offline behavior, and installation ownership are described | **Issue** — paths are expressly provisional and unapproved, so no official selection exists |
| Discovery and execution boundary | Read-only discovery proposes; structured argv collection is explicit | unavailable, failed, timed-out, and cancelled are distinct | no shell, installation, implicit scan, or mandatory Savepoint runtime | Proven as a conditional proposal |
| Normalization | Available values carry scope and provenance | absent, unsupported, failed, partial, and stale never become zero | input/row/time caps yield explicit partial results | Proven as a conditional proposal |
| Snapshot and consumers | collector writes immutable identity; readers receive summaries | missing or failed evidence remains visible | manual and official snapshots remain separate; board/doctor do not scan | Proven as a conditional proposal |
| Series compatibility | stable provider/config/report/scope identity is named | unknown vulnerability DB identity is not silently comparable | provider, version, report definition, exclusions, or relevant config starts a new series | Proven as a conditional proposal |
| Five-measure boundary | all five in-scope capabilities are present | unsupported remains visible | hotspots and dead code are excluded from collection, storage, trend, and score | Proven |
| Owner decision / settled contract | expected: explicit confirmation of five paths and boundary | actual: records say confirmation remains pending | Task-level acceptance of C-935/C-936 does not rewrite the explicit provisional status | **Issue I-084** |

All applicable cells were classified. Duplicate and mixed-type inputs, mutable input, Unicode/text-width classes, browser/server redirects, transaction cleanup, and serialization round trips are not applicable to this planning-only Objective. Provider subprocess failure, timeout, cancellation, bounded output, and secret-safe sanitized evidence are specified at the contract level; their implementation belongs to later Objectives.

## Workflow and Side-Effect Lock

O-026 adds no production command or side-effecting workflow. The proposed future sequence is discovery (read-only), explicit confirmation/configuration, explicit collection, normalization, snapshot persistence, then bounded reads. Current Check, board, and doctor paths were inspected to confirm that the proposal does not claim existing collection behavior. No provider was installed or executed, no dependency was added, and no production file was changed.

## Acceptance Classification

| O-026 success condition | Classification | Evidence |
| --- | --- | --- |
| Compare candidates across the required dimensions | Proven | O-026 records alternatives, supported stacks, output forms, licence/runtime/platform constraints, offline behavior, proposed bounds, and installation ownership. |
| Select one official provider/report path per capability | **Issue** | O-026 lines 48-52 explicitly state that all five paths remain proposals awaiting owner confirmation. |
| No mandatory ecosystem runtime for duplication | Proven | The proposed jscpd native CLI remains project-owned and does not add a Savepoint runtime or Go dependency. |
| Complexity external; hotspots excluded | Proven | Lizard remains external; hotspot material is labelled historical and excluded from the V2.1 contract. |
| Settle narrow module interfaces | **Issue** | O-026 lines 116-118 call the boundary a proposal that does not select or approve the provider inputs; T-050 evidence says final validation remains pending. |
| Define series compatibility | Proven | Capability/provider/version, report/measurement definition, scope/exclusions, and relevant configuration form the series identity, with vulnerability-database uncertainty preserved. |

## Adversarial Pass

- Bypass through Task Checks: C-935 and C-936 are CLEAR only for their conditional Quick scopes and both explicitly require later owner validation; they cannot prove the Objective's selected-provider outcome.
- Representation switch: recording `accepted_check` on the Tasks proves the owner accepted those Task Checks, but it does not negate the authoritative Objective text that the paths and boundary remain provisional and unapproved.
- Missing/failed-state probe: absent reports, unavailable executables, malformed output, caps, timeout, cancellation, unsupported stacks, and stale evidence remain distinct rather than becoming clean zero.
- Compatibility probe: changing provider/version, schema/measurement definition, scope/exclusions, or relevant configuration starts a new series; changing only the observed commit does not.
- Scope bypass: historical hotspot research is not wired into the five-capability contract. No dead-code or hotspot measure survives in the proposed collection/storage/read boundary.

## File Reality and Gates

Every file named in Task evidence exists. The scoped changes are planning records only; no phantom production or test file was claimed. `git diff --check` passed. Fresh `make test-full` passed on 2026-09-29 with all Go tests and Linux, Darwin, and Windows builds succeeding. The worktree also contains unrelated owner work for O-033/I-083, which is excluded from this Check and was not modified.

## Issues and Materiality

| Issue | Likelihood | Impact | Materiality | Recommendation |
| --- | --- | --- | --- | --- |
| I-084 — O-026 provider paths remain unselected | High — the pending state is explicit in every governing record | High — later schema and adapter work would otherwise proceed without the provider decision this Objective exists to make | High | Fix now by recording the owner's provider and boundary decisions, then re-check the frozen scope |

## Design Reconciliation

The conditional boundary is consistent with Design's file-first architecture, `internal/` ownership, IO-free rendering, Check integration, failure-mode, distribution, and testing constraints. It introduces no production behavior yet. The reconciliation cannot be final while the catalogue and boundary remain explicitly provisional; I-084 blocks closure rather than silently promoting those proposals into Design truth.

## Observations

- The proposed 120-second, 32-MiB, and 10,000-row limits remain explicitly unmeasured. This is not an O-026 blocker because the records preserve that limitation and assign validation before release.
- O-027 still describes six signals. T-050 already identifies reconciliation to the confirmed five-measure scope as O-027's prerequisite, so it is outside this Objective's closure perimeter.

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

## Verdict

**NEEDS WORK.** I-084 blocks O-026 technical clearance. The Tasks remain done; remediation belongs under the Issue and must not reopen them. A later independent Full Objective Check must create a new Check record and supersede C-937.
