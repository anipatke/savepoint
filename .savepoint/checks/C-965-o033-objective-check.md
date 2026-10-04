---
id: C-965
scope: {kind: objective, id: O-033}
result: NEEDS WORK
checked_by: {role: checker, session: check-o033-20261004}
executed_session: o033-t102-t107-20261004
checked_at: '2026-10-04T03:26:27Z'
health_snapshot: sha256:4913096b4939c4b2c03005e6c02d9fe3b452bde0a63e1ee624bb7cdbc5871cdf
reviewed:
  base_commit: 4520c93a7421dd9cc80bfe6e83cb14445aaac506
  head_commit: 53fbe53405cdfc99636f57e830c0ed08a6f38d9a
  files:
    - .savepoint/Design.md
    - AGENTS.md
    - CHANGELOG.md
    - README.md
    - agent-skills/references/check-method.md
    - agent-skills/savepoint-check/SKILL.md
    - agent-skills/savepoint-design/SKILL.md
    - agent-skills/savepoint-task/SKILL.md
    - internal/board/v2/card.go
    - internal/board/v2/column.go
    - internal/board/v2/detail.go
    - internal/board/v2/detail_view.go
    - internal/board/v2/lanes.go
    - internal/board/v2/lanes_test.go
    - internal/board/v2/load.go
    - internal/board/v2/parallel_advice_test.go
    - internal/board/v2/plain.go
    - internal/board/v2/releases_test.go
    - internal/board/v2/update.go
    - internal/data/concurrency_plan_v2.go
    - internal/data/concurrency_plan_v2_test.go
    - internal/data/concurrency_v2.go
    - internal/data/concurrency_v2_test.go
    - internal/data/next.go
    - internal/data/objective_v2.go
    - internal/data/project.go
    - internal/data/task_v2.go
    - internal/data/write_test.go
    - internal/init/template_freshness_test.go
    - internal/init/upgrade_test.go
    - internal/init/v2_scaffold_test.go
    - internal/resume/concurrency.go
    - internal/resume/concurrency_test.go
    - internal/resume/resume.go
    - internal/resume/resume_test.go
    - internal/styles/styles.go
    - main.go
    - main_resume_matrix_test.go
    - templates/project-v2/AGENTS.md
    - templates/project-v2/agent-skills/references/check-method.md
    - templates/project-v2/agent-skills/savepoint-check/SKILL.md
    - templates/project-v2/agent-skills/savepoint-design/SKILL.md
    - templates/project-v2/agent-skills/savepoint-task/SKILL.md
    - internal/data/concurrency_integration_test.go
    - internal/board/v2/concurrency_integration_test.go
  dependencies:
    - internal/data/gate_v2.go
    - internal/data/dependency.go
    - internal/data/write.go
    - internal/board/v2/watch.go
    - internal/board/v2/width.go
    - internal/init/upgrade.go
    - .github/workflows/ci.yml
issues: [I-135, I-136, I-137, I-138, I-139, I-140]
supersedes: null
---

# C-965: O-033 Full Objective Check

NEEDS WORK: five reproducible feature-contract failures plus missing native
Windows evidence. This conversation did not build the reviewed implementation
or its Tasks. Prior activity checked O-037, not O-033 implementation. This is
an initial Full Check, not a recheck of I-083 (resolved by escalation).

Scope includes HEAD and the current uncommitted integration/guidance/docs,
not HEAD alone. Hashes of 45 scoped files in /tmp/o033-reviewed-hashes.json were
checked unchanged after gates/probes. All six owned Tasks are owner-completed
with explicit board-owner Task-check waivers. The checker changed no code,
Task/Objective status, router or Design. Tests use temporary projects/Go overlays;
no real worktree or agent execution was created.

## Frozen scope lock

# O-033 Full Objective initial scope lock
1. Requirements: Objective SC1–10; T102 DW1–7, T103 DW1–7, T104 DW1–7, T105 DW1–7, T106 DW1–7 (actual enumerated rules), T107 DW1–7; named FS/DATA/TPL/ARCH/CFG/TEST Guardrails. STYLE advisory. Gate fresh build/full, native Windows CI, official health. Owner walkthrough separate from technical clearance.
2. Public surfaces: Task/Objective YAML decode, LoadV2Index planning validation, scopes/path comparison/review digest, ResolveConcurrencyV2/ResolveNext, board grouping/load/update/columns/headings/details, resume Render/ParallelLines/instruction formatter, non-TTY board, init/upgrade guidance. O-033 code from 4520c93 to current HEAD53fbe53 plus current uncommitted T105–107 guidance/integration files. No unrelated O037 fixes.
3. Required orchestration: existing lifecycle/dependency gates; configured preference reader/writer/reload; YAML reader; explicit status/priority writes; index validation; pure projections; text sanitizers/width; canonical/scaffold asset delivery. No shell/worktree orchestration or arbitrary dependency internals.
4. Coverage matrix frozen before probes:
|Cell|Inputs/states/surfaces|Expected rule/oracle|
|M1 decode|old/mixed/new; omitted/empty/null/scalar/map/list/wrong/mixed type/nonfinite/duplicate fields and lane declarations; unresolved owner/local lane; malformed independence pair/reason/path/review|nonfatal named advice diagnostics, required schema/record failures unchanged; load/start still available|
|M2 paths/scopes|exact new file/duplicates/shared reads/writes/write-read/case aliases; slash/backslash/absolute/drive/UNC/traversal/directory/glob/control/Windows-ambiguous; receiver/path copies and mutation after validation|exact lexical portable rule; unknown withholding; no filesystem probing; no stale exemption|
|M3 recommendation|off/on/default; selected Objective, invalid/stale selection/replan; all 3 statuses; direct/transitive deps including done intermediate; clear/accepted/waived/blocked; same lane/empty/singleton; active known/unknown/unlaned/multiple active; remaining lane known/unknown/conflicting; explained overlap/stale digest/wrong direction|one lane head, no incompatible clique; dependency and scope reasons, deterministic results, no new gate|
|M4 transitions/actions|toggle/restart/external reload/stale/failed save; start->inprogress->done, done->backward by owner, selected task after regrouping; ignored lane and wider manifest; priority writes|stable membership/focus/counts, authored bytes preserved, same lifecycle and Next decisions off/on|
|M5 render/navigation|board columns/details/resume/plain; focused/unfocused/compact; middle of lane, heading heights, empty heading, long title; widths30/40/47/80/120/160 and heights20/24/30/48/60; ASCII/control/combining/VS/wide/modifier/flag/joined emoji|ANSI-aware width oracle; focused ID visible if card itself fits; heading nonselectable; no rendering IO/control injection|
|M6 reasons/instructions|selected task/objective, blocked/unknown/singleton/eligible parallel/active, stale router/replan, after edits/status changes; Goal-wide namespace; complete copied session|same reasons for same selection; exact executable Start line selects intended record; objectives/tasks/lane/scopes/prerequisites/worktree rules; no cross-Objective advice|
|M7 persistence/guidance|canonical/scaffold parity; old/clean/edited assets, dry-run/apply/repeat, config choices preserved, no backfill; docs/Design/API|byte comparison and fixture independent oracle, advisory skills, no health switch or self-clearance|
|M8 gates/evidence|file reality, all waivers, full host/crossbuild/native Windows, official health, owner walkthrough|current attributable result, no presumed completion|
External finite boundary: startup/configured target and runtime root via load; read success/missing/unavailable/malformed tested; writes success/refusal/retry/cleanup through existing writer tests; redirects/server response/network timeout N/A local advice/no server. No new subprocess/browser/provider dependency; official health uses established Full Check boundary. No secret fields in new output; record text sanitized.
5. Admit only supported-path violations introduced/touched/promised by this Objective, named criteria/Guardrails and credible consequence. Advice accuracy is feature behavior; ignoring advice is never a finding. Symlink/normalization alias uncertainty assessed within pure lexical promise, no generalized filesystem resolver requirement. No optimal scheduling, actual merge monitoring or new lane-enforcement gate.

Workflow/side effects:
|Order|Operation|Effect/failure/final state|Oracle|
|1|schema/router/index load; typed decode; plan validation|reads only; required fatal versus optional nonfatal|temporary authored input/load diagnostic|
|2|preference read; Next/concurrency/grouping|pure immutable-result calculation; advice withheld doesn't block work|independent matrix and original lifecycle gates|
|3|render/details/resume/plain/instructions|no IO, sanitize, output error propagated|exact Start/text/control and ANSI width|
|4|explicit toggle/action command|narrow config or managed record write; source freshness; temp replacement/error/cleanup|exact bytes, no premature success, lifecycle gate unchanged|
|5|watch/reload/status transition/regroup|focus ID maintained; membership follows status|file/state snapshots and repeated workflow|
|6|init/upgrade|managed asset delivery; edited files preserved/sidecar, no task/config rewrite/backfill|exact template bytes/dry-run/repeat fixture|


Initial-pass oracle correction: the first copied-Start probe compared a raw
control-bearing title against sanitized advice. The affected cell was rerun
with a plain "Core work" title; it still lacks a standalone Start line. The
sanitization difference itself is correct. Same-selection parity was completed
with both missing Task and missing Issue (latter retains the exact same T-002).
These refine the existing M6 cells rather than adding axes.

## Completed coverage matrix

| Cell | Result and evidence |
|---|---|
| M1 legacy/mixed/new, absent/empty, invalid shapes/types, nonfinite, unresolved lanes/owners | Passed: original concurrency_plan_v2_test and integration tables plus independent mixed-node/digest probes; named nonfatal diagnostics for supported advisory value errors; required record/schema failures remain fatal |
| M1 duplicate lane declarations | I-135: diagnosed but core Task still recommended in group; missing-title adjacent case withholds correctly |
| M1 duplicate top-level frontmatter fields | Existing YAML structural rejection retained, classified required-document error rather than a new advisory feature defect; no silent last-wins |
| M2 paths/unknown scopes/review binding | Passed documented lexical path classes and copy accessors; existing alias/path/independence validation tables, independent ScopeAndDigest. Edited either manifest invalidates persisted digest on reload; no filesystem probe |
| M2 direct mutable model / lexical limitations | Observations: ObjectivePlan exported slices are mutable, unlike Task scope Paths copies; no supported runtime caller mutates them during this workflow. Pure index projection trusts prior validation. Unsupported Windows path names admitted by validator noted below; no new runtime alias resolver demanded |
| M3 status × scope-conflict classes | Passed independent AdviceMatrix: 3 states × 6 scope classes, on/off membership and unchanged start gate, no unsafe group/active recommendation for shared writes/read-write/case alias/unknown |
| M3 dependency/waiver/accepted/replan/active/remaining scopes | Passed original named projection tests: prerequisitesUseOrdinaryStartGate, objectiveDependencyBlocksEveryTask, dependenciesNeverTogether, oneTaskPerLane, remainingLaneScopeConflictsWithhold, recordedActiveWork, staleExplanationWithholdsAndSaysSo, replanAndStaleSelectionWithholdLaunchAdvice. Integration loads persistent fixtures after edits |
| M3 determinism/cliques/empty/singletons | Passed independent eight conflict graphs, every emitted group checked by direct write-set intersection; ten repeated projections per graph. Existing incompatibleTriangleNeverGroupsAConflictingPair, singletonAndEmptyClaimNothing, isDeterministic also pass |
| M3/M6 stale selection across consumers | I-137: data/Next withholds correctly; board consumers omit diagnostic and advertise Start under same stale Issue selection |
| M4 on/off/default/restart/watch/refusal/status/focus | Passed full options suite and concurrency integration, plus independent OffReloadReadOnly byte/focus snapshot. All status membership columns/grouping preserved; existing gates unchanged when lane/manifest advice ignored |
| M4 managed writes | Passed WriteTaskV2_preservesPlanningMetadataByteForByte and WriteObjectiveGroupOrderV2_preservesLanesAndIndependenceByteForByte; underlying persistence ordering unchanged. No new planning writer or lane side effect |
| M5 Unicode/width/height | Passed independent 7 text classes × 6 widths × 5 heights = 210 board cases, including long repeated headings; independent 4 widths × 6 heights × 6 focused cards = 144 focused-window boundary cases; exact ANSI-aware widths and IDs checked. Existing compact/mid-lane/heading geometry tests also pass |
| M5 authored control text | I-136: raw clear-screen ESC[2J and BEL survive in interactive heading output. Shared advice formatter strips them correctly; plain board strips controls |
| M5 modes/backends/no IO | Source tracing and existing redirected/non-TTY/TUI backend tests passed; no changed cursor/profile policy. Rendering receives loaded values and no source file access; commands own writes |
| M6 normal/off/blocked/unknown/eligible/Goal-wide parity | Passed existing ParallelAdvice/ParallelLines and new IntegrationSurfaces tests; no cross-Objective opportunities; identical instruction blocks across surfaces for valid selection; feature-off advice absent with metadata preserved |
| M6 copied session selection | I-138: copied T-004 instruction contains prose but no standalone canonical Start selection when router still selects T-002. Deterministic routing-contract probe, no inference of actual model behavior |
| M6 actionable planning diagnostics | I-139: c.Diagnostics/PlanDiagnostics collected but unused by public formatter/surfaces; tells reader to see nonexistent visible diagnostics |
| M7 canonical/scaffold/guidance/upgrade/adoption | Passed independent byte comparison of all seven skill/reference pairs, original/new upgrade preservation/guidance tests, actual scaffold/runtime suite, no backfill/config/Task rewrite. All changed guidance is advisory, Code Health core, existing worktree rules maintained |
| M7 docs/Design | Ownership and advisory boundaries reconciled; README stale-selection and complete-session claims currently fail I-137/I-138. Stored raw planning diagnostics never shown, contradicting visible-diagnostic promise I-139 |
| M8 local gate/health/waivers/files | Passed fresh full suite/crossbuild, official nonblocking health collection and scoped file/hash reality. Six explicit waivers intact |
| M8 native Windows | Unverified I-140: available CI run covers older O-037 commit, no final O-033 source/test revision proof |
| M8 owner walkthrough | Pending owner report for T-104/T-105/T-107. Owner asked why required; no completion/acceptance inferred. This is owner validation, not a technical defect |

Every cell classified before verdict. Issues did not stop later matrix rows.
Finite boundary and workflow inventory 1–6 were source-traced and exercised via
full gate/persistent fixtures and probes. Advice reads and displays only, no
subprocess/network/worktree operations. Explicit config/status/priority writes
retain existing freshness/temp/cleanup behavior; normal execution never reads
lane compatibility as a gate. Missing/malformed config only removes optional
advice; status selection/record validation remains canonical.

## Acceptance coverage

| Requirement | Classification |
|---|---|
| O-033 SC1/2 | Proven: opt-in saved setting, Code Health unchanged |
| O-033 SC3 | Proven stable membership/grouping/counts; control-text rendering gap I-136 |
| O-033 SC4 | Issue I-137: same selected Task/index differs on stale evidence |
| O-033 SC5 | Issue I-135: malformed duplicate lane still recommended; ordinary conflicts/dependencies proven |
| O-033 SC6/7 | Proven: advice never becomes execution/Check policy; dependencies/authority/worktree rules unchanged |
| O-033 SC8 | Issues I-138/I-139: copied Start routing and visible diagnostics incomplete |
| O-033 SC9/10 | Proven: off restores ungrouped view without record loss, old projects/config safe |
| T-102 DW1/2 | I-135/I-139 for affected advice/visibility; typed diagnostic/local lane data otherwise proven |
| T-102 DW3/4/5/6 | Documented path/scope/digest/managed-preservation behavior proven with exported-slice/lexical limitations retained as observations |
| T-103 DW1/2/3/4/5/6 | Core conservative projection and cliques proven; consumer propagation I-135/I-137 |
| T-104 DW1/2/3/4/5/6 | Grouping/focus/counts/window/determinism proven; authored heading controls I-136 |
| T-105 DW1/2/3/4/5/6 | Optional/common formatter/worktree restrictions proven; I-137/I-138/I-139 and heading bypass I-136 |
| T-106 DW1–6 | Proven: advisory guidance, parity, preserved upgrade delivery/no backfill |
| T-107 DW1/2/3/4/5 | Settings/isolation/persisted integration/docs/Design proven in normal cases; inherited I-135–I-139 gaps |
| All Task final evidence/waiver criteria; T-107 DW6 | Fresh host gates/waivers proven; native platform gate Unverified I-140. Owner walkthrough still pending |

## Evidence and commands

- Fresh `make build && make test-full`, 2026-10-04 UTC, Go 1.26.2 linux/amd64:
  exit 0; host packages and linux/darwin/windows cross-builds. No executor gate
  reused. `git diff --check` exit 0. No configured separate lint/typecheck gate.
- Independent `/tmp/o033-overlay.json` adds only scratch Go tests. Data
  TestO033AdviceMatrix/ScopeAndDigest/CliqueMatrix: exit 0. Initial complete
  data pass fails only duplicate-lane suppression as recorded; portable-name
  and duplicate-frontmatter probes are separately classified observations.
- Board complete TestO033 pass: 210 text/sizing cases and off/reload snapshot
  passed; controls, stale-selection, diagnostics and Start-line probes failed
  as recorded. WindowBoundaryMatrix: 144 cases passed. Plain-title Start
  reproduction and exact same-Task stale Issue probe fail independently.
  Cache access refusals resolved through approved escalation, not test failure.
- Existing complete named unit/integration tables run in fresh full gate:
  metadata/independence/preservation, dependency/accepted/waiver/active/future
  lane conflicts, options failure/restart, board focus/columns/heading windows,
  valid instruction parity, old/mixed records, post-replan/status changes,
  Goal-wide namespace and no cross-Objective comparisons.
- Independent Python comparison: seven canonical/scaffold skill/reference
  pairs byte-identical; 45 scoped current source/test/doc paths exist and
  content hashes unchanged after gate and probes.
- GitHub `gh run list --branch v2.20 --limit 6 --json ...`, independently read
  with approved network access: only successful CI run 37171504754,
  head 4520c93a7421dd9cc80bfe6e83cb14445aaac506; does not cover O-033 HEAD
  or uncommitted integration/guidance tests. No Windows failure alleged.
- Official `savepoint health check O-033` after full gate: first snapshot
  sha256:c307f7538454f5bc649bf07867440501a58fd555cbe43893134f53b39a209ecd
  had optional OSV collection failure due sandbox DNS. Approved network rerun
  saved frontmatter snapshot, all five instances collected, Code Health does
  not block clearance. Both immutable snapshots retained; failed collection
  was never classified as unhealthy code.

## Adversarial pass

Nonfatal malformed Task scope works, but malformed Objective duplicate lane
is still used (I-135). Input copies preserve Task manifests; stale persisted
explanations filtered on reload. Shared-read/work overlap/dependency paths and
first-fit clique properties tested independently. Normal owner status changes
and ignored lane metadata leave gate decisions identical. Consumer input
propagation bypasses stale-selection withholding (I-137). Interactive headings
bypass text sanitation (I-136). Public copied instruction must obey the canonical
routing contract, not just mention a Task (I-138). Generic unknown-scope text
cannot substitute for the promised visible named diagnostic (I-139). No auth,
billing, tenant, cloud or new secret boundary is introduced. Full local and
crossbuild evidence cannot establish native execution (I-140).

## Materiality

| Issue | Likelihood | Impact | Materiality | Recommendation |
|---|---|---|---|---|
| I-135 | Medium (edited lane declarations) | Low (ambiguous advice) | Medium | Narrow affected-lane suppression |
| I-136 | Low (control-bearing authored title) | Medium (terminal cleared/altered) | Medium | Sanitize new heading text |
| I-137 | Medium (stale router context) | Low (conflicting launch advice) | Medium | Propagate matching selection evidence |
| I-138 | High (every copied parallel instruction) | Medium (fresh session routing ambiguous) | Medium | Emit explicit standalone Start line |
| I-139 | Medium (malformed/stale plan editing) | Low (cause/file hidden) | Medium | Show collected sanitized diagnostics |
| I-140 | High (no current attributable run) | Medium (platform unverified) | Medium | Native CI after narrow repairs |

## Nonblocking observations

- Duplicate top-level YAML fields remain required document failures; malformed
  value-shape diagnostics are advisory. No arbitrary syntax recovery required.
- Validator accepts Windows-reserved/device/illegal-character strings such as
  nul, CON.go and x|y.go. Those are unusable authored paths on some platforms;
  documented case-alias/path-class tests pass. No extra alias/device policy
  or runtime filesystem probe was silently added to this Check.
- ObjectivePlan exported slices are mutable to Go consumers; Task manifest
  accessors clone. Current runtime readers do not mutate Objective planning
  slices and no supported persisted edit bypasses reload validation. This
  limits the immutability claim without establishing a user-visible defect.
- Pure lexical comparisons cannot discover symlink/Unicode-normalization
  aliases. Advice is explicitly not guaranteed independence, and no filesystem
  monitoring/resolver was promised; existing case-alias handling is tested.
- Original replacement helper's OS short-write/sync/close failures were source
  traced rather than newly fault-injected; existing real refusal/cleanup tests
  and byte-preservation tests passed. No new persistence implementation here.
- Owner walkthrough remains separate from technical findings; requesting it
  does not authorise acceptance or change owner-completed Task statuses.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [ ] STYLE-03 **Test branches** — duplicate-lane suppression, stale consumer parity, heading controls and diagnostic rendering missing from existing tests; I-135–I-139.
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [x] STYLE-07 **One source of truth**
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs**

## Handoff

Direct narrow repairs belong to I-135–I-139; keep all completed Tasks done.
Provide native CI proof for I-140 after the final repaired inputs. Recheck this
frozen M1–M8 matrix and original reproductions independently; do not broaden the
scope or turn advice into a gate. Owner walkthrough/acceptance remains pending.
