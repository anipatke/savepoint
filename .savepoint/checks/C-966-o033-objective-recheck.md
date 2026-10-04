---
id: C-966
scope: {kind: objective, id: O-033}
result: CLEAR
checked_by: {role: checker, session: check-o033-recheck-20261004}
executed_session: o033-t102-t107-20261004-and-owner-repairs-97eff73
checked_at: '2026-10-04T03:46:04Z'
health_snapshot: sha256:37d89f3767d1063cb381385e65674a75a20b748508d14cf0bf1d76eb7f92d3b3
reviewed:
  base_commit: 4520c93a7421dd9cc80bfe6e83cb14445aaac506
  head_commit: 19a6807c821be1d6372c96f93fd8c55edf91ba91
  files:
    - .github/workflows/ci.yml
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
    - internal/board/v2/concurrency_integration_test.go
    - internal/board/v2/detail.go
    - internal/board/v2/detail_view.go
    - internal/board/v2/lanes.go
    - internal/board/v2/lanes_test.go
    - internal/board/v2/load.go
    - internal/board/v2/o033_repairs_test.go
    - internal/board/v2/parallel_advice_test.go
    - internal/board/v2/plain.go
    - internal/board/v2/releases_test.go
    - internal/board/v2/update.go
    - internal/data/concurrency_integration_test.go
    - internal/data/concurrency_plan_v2.go
    - internal/data/concurrency_plan_v2_test.go
    - internal/data/concurrency_repairs_test.go
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
    - internal/resume/concurrency_repairs_test.go
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
  dependencies:
    - internal/data/gate_v2.go
    - internal/data/dependency.go
    - internal/data/write.go
    - internal/board/v2/watch.go
    - internal/board/v2/width.go
    - internal/init/upgrade.go
issues: []
supersedes: C-965
---

# C-966: O-033 Full Objective Recheck

## Closure map

| Prior Issue | Outcome | Independent proof |
|---|---|---|
| I-135 | Closed, verified | Duplicate-core and missing-title original probes pass; malformed lane Tasks receive no recommendation; ordinary start still allowed; unaffected valid lanes remain candidates |
| I-136 | Closed, verified | Original ESC[2J/BEL lane title is sanitized in loaded heading and interactive output; adjacent Unicode/width/focus matrices pass |
| I-137 | Closed, verified | Same T-002 with missing I-999 and missing T-999 router scenarios no longer advertise instructions; board helper carries the same selection diagnostic as resume; normal/replan/status surfaces pass |
| I-138 | Closed, verified | Copied T-004 instruction has a standalone exact NextLine Start selection while shared router names T-002; prerequisites/scopes/worktree restrictions preserved |
| I-139 | Closed, verified | Original malformed glob emits named file/record/field/message; diagnostic-only body remains visible and sanitized; off hides the section |
| I-140 | Closed, verified | Independently fetched CI run 37174579468 has successful ci and native windows-tests on exact reviewed HEAD 19a6807 |

CLEAR within C-965's frozen M1–M8 scope. This checker conversation did not
implement the Tasks or repairs. Owner remediation reports were treated as
claims; no executor result was used as independent clearance. This is the one
full recheck following the initial Check. No implementation, router, Task
status, acceptance criteria or Design was changed by the checker.

## Frozen scope and admission ledger

The immutable scope below is carried from C-965. Existing axes apply to the
changed helpers; no new public value, dependency layer or acceptance meaning
was admitted as a blocking perimeter. Initial nonblocking observations remain
observations. Owner validation remains distinct from technical evidence.


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


# O-033 recheck admission ledger (before probes)
Frozen source: C-965 scope M1–M8 and original issue Proof Needed; no expansion.
| Item | Prior claim/Issue | Exact frozen cell | Allowed result |
|---|---|---|---|
| Original malformed duplicate + missing-title pair, ordinary start invariant | I-135 | M1 duplicate lane declarations, original TestO033MalformedLaneSuppressesAdvice | resolve or remain I-135 |
| Original ESC[2J and BEL authored heading | I-136 | M5 authored control text TestO033HeadingControls | resolve or remain I-136 |
| Original missing Task T-999 and same Task T-002/missing Issue I-999 across consumers | I-137 | M3/M6 stale selection, TestO033StaleSelectionParity/StaleIssueSameTaskParity | resolve or remain I-137 |
| Original plain-title T-004 copied block while router selects T-002 | I-138 | M6 copied session TestO033FreshSessionStartLine | resolve or remain I-138 |
| Original T-004 bad/*.go and no-other-advice diagnostic fixture | I-139 | M1/M6 actionable diagnostics TestO033MalformedDiagnosticVisible + formatter regression | resolve or remain I-139 |
| Current SHA Windows full test job | I-140 | M8 native Windows | resolve or remain unverified |
| Original 18 status/conflict, eight graphs × ten repeats, scopes/digest | unchanged invariant | M2/M3 AdviceMatrix/CliqueMatrix/ScopeAndDigest | pass or exact existing-cell Issue |
| Original 210 heading text/size +144 focus boundary cases | I-136 adjacent sizing invariant | M5 HeadingTextMatrix/WindowBoundaryMatrix | pass or exact existing-cell Issue |
| Original off/reload bytes+focus, original package tables+integration | repair nonregression | M1–M7 rows named in C-965 | pass or exact existing-cell Issue |
| Canonical parity, scope file reality, waivers, full host/crossbuild, health | required gate | M7/M8 | proven or unverified |
| Owner walkthrough | C-965 pending | M8 owner validation | remains pending absent owner acceptance |


## Completed coverage matrix

| Frozen cell | Recheck classification and evidence |
|---|---|
| M1 old/mixed/new, absent/empty/null/scalar/map/list/wrong/mixed types, nonfinite, unresolved owner/local lane and malformed independence | Passed full concurrency_plan_v2_test and integration tables; original mixed-node independent ScopeAndDigest probe; required structural YAML errors remain errors, advisory value diagnostics remain nonfatal |
| M1 duplicate-core/board and missing-title original reproductions | Passed TestO033MalformedLaneSuppressesAdvice both cells; diagnostic retained, no ambiguous group and ordinary start allowed; new duplicate-lane regression preserves unrelated board/extra candidates |
| M2 exact portable paths, conflicts, copies, reviewed-scope invalidation | Passed existing Validate/Compare path tables and independent ScopeAndDigest; Paths copy cannot mutate authored scope; edited manifests stale the digest; lexical limits unchanged |
| M3 all statuses × scope classes | Passed independent AdviceMatrix: 3 statuses × disjoint/shared-read/shared-write/write-read/case-alias/unknown; lane membership and existing start gate unchanged off/on, no incompatible groups/active suggestions |
| M3 dependency paths, clear/accepted/waived/blocked, same lane, active and remaining lane scope, replan | Passed full ResolveConcurrencyV2 named tables and loaded-project integration; dependencies never treated as explanatory overlap; unknown active/future lane scope withheld |
| M3 empty/singleton/determinism | Passed eight independent shared-write graphs, every group checked by direct write-set intersections, ten repeated projections per graph; original singleton/empty/first-fit deterministic tables passed |
| M3/M6 missing Task T-999 and missing Issue I-999 with valid T-002 | Passed original StaleSelectionParity and StaleIssueSameTaskParity plus new EverySurface test; projectConcurrency used by lane readiness/details/plain; withheld heading readiness blank, resume and board receive matching selection evidence |
| M4 preference default/on/off/restart/external reload/save refusal and lifecycle | Passed full OptionsJourney, OptionsLeaveLifecycleAndBoardIdenticalWhenOn and Integration_adviceNeverChangesExistingDecisions; independent OffReloadReadOnly confirms focus and all non-config bytes unchanged |
| M4 start/in-progress/done/regroup/priority preservation | Passed HeadingPersistsWhileATaskMovesThroughEveryColumn, IntegrationSurfacesAgreeAfterStatusChanges and managed-write preservation tests; advice does not enter any lifecycle gate |
| M5 text/size/navigation | Passed original 210 text-class × width × height cases and 144 focused-window boundary cases; ANSI-aware width and focused ID oracle, existing compact/mid-lane/count/nonselectable tests passed |
| M5 authored ESC[2J/BEL heading | Passed loaded-board isolated lane reproduction and new TestLaneHeadingControlsAreStripped; heading title sanitized at lane layout, Unicode title retained; shared formatter/plain sanitation unchanged |
| M5 backend/modes/rendering IO | Passed existing full redirected/TUI/no-TTY tests and source trace: new helper reads loaded state only, renderers have no filesystem/subprocess operation; generated style escapes retained |
| M6 normal/off/blocked/unknown/singleton/group/active, Goal-wide namespace, status/replan | Passed original ParallelAdvice/ParallelLines tables and IntegrationSurfaces*; valid T-002 instruction blocks match across resume/Objective/Task/plain, no cross-Objective opportunity |
| M6 copied T-004 block under T-002 shared router | Passed TestO033FreshSessionStartLine; standalone Start line exact for intended Task; savepoint-task separate, scope/prerequisites/router/Goal/no allocation/local commit/no push/merge/main Check instructions retained |
| M6 actionable diagnostics including otherwise empty/withheld body | Passed original glob probe and new ParallelLinesShowPlanningDiagnosticsEvenWithNoOtherAdvice; formatter starts with diagnostics before withheld/normal advice and sanitizes complete output; Task focus includes its own diagnostics and Objective focus includes Objective-wide diagnostics |
| M7 skills/scaffolds/init/upgrade/adoption | Seven canonical/scaffold pairs independently byte-compared; full Guidances/UpgradeDelivers/UpgradePreserves/Scaffold fixtures passed; no config or Task backfill, edited assets preserved, advisor remains optional |
| M7 Design/documentation | Source ownership matches shipped Design section 8; README/CHANGELOG advisory and manual workflow claims reconcile; shared Start/diagnostic/selection behavior now meets original failed promises |
| M8 files/waivers/full host/crossbuild | All 49 reviewed paths present, hashes unchanged after verification; six owner waivers retained; fresh make build && make test-full exit 0; git diff --check exit 0 |
| M8 native Windows | Passed attributable current CI run37174579468, HEAD19a6807, ci and windows-tests success; workflow runs Windows native full uncached Go suite |
| M8 official health | Passed official collection after full gate; final snapshot all five instances collected, nonblocking verdict |
| M8 owner walkthrough/acceptance | Pending owner validation; no acceptance naming C-966 supplied or written. Technical CLEAR does not claim owner walkthrough or close the Objective |

Every original matrix cell was reclassified, with remaining cells completed
after probes. Full suite covers original package tables, including unchanged
persistence/refusal/retry/cleanup cases; independent probes supplement rather
than replace them. No scope amendment or new blocking cell.

## Workflow and finite external boundary

The six-operation inventory from C-965 remains applicable. Source trace:
index load/decode/validation reads only; project preference and canonical
Next/concurrency projection use loaded values; output is pure; explicit
config/status/priority commands alone write; watch/reload rebuilds focus by
ID; init/upgrade owns managed asset delivery. The repaired decoder removes
ambiguous lane entries only from in-memory advice; it does not rewrite files.
projectConcurrency adds selection evidence without adding execution policy.
Diagnostic/Start formatting and heading sanitation add no side effects.

Successful/missing/malformed input and real refusal/cleanup/retry paths are
covered by fresh full fixtures; preservation uses byte-level oracles. Original
OS short-write/sync/close failure injection remains source-traced as in C-965,
not silently elevated to a new recheck requirement. Configured/runtime root,
startup/reload order and unchanged action writer ownership reconciled.
Server responses/redirects/network timeouts are not applicable to local
advice. The official health provider boundary is separate: sandbox network
collection failure was retried with approved access and never called bad code.

## Acceptance coverage and Guardrails

| Requirements | Classification |
|---|---|
| O-033 SC1–2 | Proven: saved opt-in/default/off and core Code Health contract unchanged |
| O-033 SC3–5 | Proven: stable grouping/counts/focus, same pure projection and selection evidence, conservative dependencies/path/unknown-scope recommendations; original I-135–137 repaired |
| O-033 SC6–7 | Proven: ignored advice/lane/manifest adds no gate; real dependencies/owner/worktree authority preserved |
| O-033 SC8 | Proven: standalone intended Start and complete optional instructions, named sanitized nonfatal diagnostics; I-138–139 repaired |
| O-033 SC9–10 | Proven: off/legacy/mixed compatibility, record/config preservation and settings conflict refusal |
| T-102 DW1–7 | Proven across parsing/scopes/independence/byte-preservation/regression tables, fresh gate and existing explicit owner waiver; lexical/exported-value observations retained |
| T-103 DW1–7 | Proven across canonical preference/eligibility/dependency/active/future/clique/stale selection tables, independent graph/status matrix and gate/waiver |
| T-104 DW1–7 | Proven across grouping/layout/navigation/headings/status/readiness tests and independent geometry/text matrix, gate/waiver |
| T-105 DW1–7 | Proven across shared formatter/parity/Start/instructions/diagnostic/control/blocked/off/no-cross-Objective tests, gate/waiver |
| T-106 DW1–7 | Proven across advisory canonical guidance, independent parity and fresh init/upgrade preservation fixtures, gate/waiver |
| T-107 DW1–7 | Proven across persisted old/mixed/new/status/replan/settings/decision integration, documentation/Design reconciliation and current host/native CI evidence; owner validation remains separate |

Applicable FS-01/02/03/04/05/06 and DATA-01/02/03/04/05 maintain their existing
preservation, diagnostic, lifecycle and authority contract in full fixtures;
new planning metadata is advisory and not silently healed into authored records.
TPL-01/02/03/04 parity/delivery/described behavior, ARCH-01/02/03/04 pure
projection/rendering/package ownership, CFG-01/02/03 settings and Windows,
TEST-01 through TEST-09 (where applicable) evidence/regressions/temp fixtures/
full gate/waivers satisfied. No new dependency or platform skip; DEP-01/02
unchanged. CI make ci additionally supplies release distribution/build evidence.
Code Style remains advisory. All technical acceptance rules classified Proven;
no material Issue or technical unverified cell remains.

## Evidence and commands

- Fresh `make build && make test-full`, 2026-10-04 UTC, Go1.26.2 linux/amd64,
  exit 0. Full uncached host Go suite and Linux/Darwin/Windows crossbuilds.
  No executor full-gate reuse. `git diff --check` exit 0; project config has
  no additional lint/typecheck command.
- `go test -overlay /tmp/o033-overlay.json ./internal/data ./internal/board/v2
  -run TestO033 -count=1 -v`, exit 0 in final isolated pass. Scratch tests and
  results `/tmp/o033-recheck-probes-final.txt`; original 18 scope/status,
  eight graphs × ten repeats, 210 text/size and 144 focus cases all pass.
  Overlay API adaptation only replaces withParallel(index,true,detail) with
  withParallel(state,detail), matching remediation. No repository test edits.
- Heading oracle clarification: original writeAdviceProject also has an
  existing unrelated Task title with BEL/SGR. Whole-view control assertion
  initially still failed because of that card, while the new heading showed
  safe "Core". Plain-title isolation of T-004 preserved the exact lane
  ESC[2J/BEL input and passed. The Task-card observation is outside I-136's
  changed heading perimeter; not a new remediation demand or scope amendment.
- New repair regressions all ran in full gate: duplicateLaneSuppressesItsRecommendations,
  LaneHeadingControlsAreStripped, StaleSelectionWithholdsLaunchAdviceOnEverySurface,
  InstructionStartsWithAStandaloneNextLine, ParallelLinesShowPlanningDiagnosticsEvenWithNoOtherAdvice.
  The owner's mutation tests for two fixes are supplemental only; independent
  original reproductions plus source/output oracle support all five repairs.
- Independently read `gh run view 37174579468 --json ...`: successful CI,
  head19a6807c821be1d6372c96f93fd8c55edf91ba91, ci success at03:37:59Z,
  windows-tests success at03:39:26Z, 2026-10-04. Run URL:
  https://github.com/anipatke/savepoint/actions/runs/37174579468.
  Inspected workflow: windows-latest, Go1.26, native buildtool test -json
  -count=1 ./...; push trigger v2.20 change matches owner's reported change.
- Official `./savepoint health check O-033` after full gate, final frontmatter
  snapshot sha256:37d89f3767d1063cb381385e65674a75a20b748508d14cf0bf1d76eb7f92d3b3,
  all five instances collected, does not block clearance. Initial immutable
  snapshot sha256:e2c454e7fbca32e40e943b8287d29281266497e8171048fd8041707978b18a3a
  retained with optional OSV sandbox collection failure; approved retry resolved
  it. Neither optional failure nor measurements were converted into a defect.
- File reality: 49 scope paths read/exist; content hashes checked unchanged
  before record writes; original unrelated coverage/Go JSON scratch left alone.
  Seven canonical/scaffold skill/reference pairs independently byte-identical.
  All six Tasks done with dated owner waivers; no Task status changed.

## Adversarial pass

Malformed advice still loads and leaves start eligibility alone; dropping
ambiguous lanes cannot suppress independent valid lanes. Existing constructors,
loaded records and scope copies/edited explanation inputs preserve withholding.
Direct/transitive dependencies, active work and incompatible clique edges cannot
be bypassed by advice on/off or ignored manifest. Shared helper transports
missing Task/Issue diagnostic to every changed board consumer; headings do not
invent ready-to-start labels under withholding. Shared formatter retains
sanitation on diagnostics and standalone Start; copied intended Task differs
from shared router without a router write. Unicode/width/focus matrix and
redirected/no-TTY fixtures establish adjacent output behavior. Side effects
remain in commands, not renderers. No new auth, cloud, billing, secret or
worktree orchestration boundary was introduced.

## Materiality and observations

No materiality actions are required; every prior Issue is verified by this
CLEAR recheck. No new Issue admitted.

Original observations remain: structural duplicate YAML fields are fatal;
portable path validation is lexical and does not discover symlink/Unicode
normalization aliases or reject every Windows device/illegal name; exported
ObjectivePlan slices remain mutable to Go callers, with no supported runtime
mutation; original OS syscall failures source-traced, not newly injected.
The existing ordinary Task-card title in the scratch fixture retains BEL/SGR;
new lane headings and optional instructions are sanitized. This is a separate
nonblocking existing-surface observation under the frozen recheck policy,
not a revived lane-heading finding.

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

## Owner handoff

C-966 supersedes C-965 without editing it. I-135–I-140 resolve as verified and
retain original evidence/history. Technical integration is CLEAR. Owner
walkthrough/acceptance for T-104, T-105 and T-107 remains pending and was not
inferred from CI or the repair report; acceptance must name current C-966
under savepoint-check Closure Rules. Only the owner may complete O-033.
No completion/acceptance, router/Goal update, implementation repair, commit,
push or merge performed in this recheck.
