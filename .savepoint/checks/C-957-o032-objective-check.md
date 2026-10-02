---
id: C-957
scope: {kind: objective, id: O-032}
result: NEEDS WORK
checked_by: {role: checker, session: check-o032-20261002-independent}
executed_session: o032-multiple-prior-executor-sessions
checked_at: '2026-10-02T05:23:00Z'
health_snapshot: sha256:9bfe2b4130e8f95cee2d757028c002f5027072b63ca00f292043557dc30a4a23
reviewed:
  base_commit: 67171b02f11faeb73e0356a9f895c4be3acad231
  head_commit: 136fe69ed1c265f2d11f5b6c2201ee346557fee5
  files:
    - .savepoint/Design.md
    - .savepoint/health/config.json
    - AGENTS.md
    - CHANGELOG.md
    - README.md
    - internal/board/v2/health_bench_test.go
    - internal/board/v2/update.go
    - internal/buildtool/main.go
    - internal/buildtool/main_test.go
    - internal/codehealth/dashboard.go
    - internal/codehealth/discovery.go
    - internal/codehealth/discovery_catalogue.go
    - internal/codehealth/discovery_test.go
    - internal/codehealth/history.go
    - internal/codehealth/history_bench_test.go
    - internal/codehealth/readers_polyglot_test.go
    - internal/codehealth/runner.go
    - internal/codehealth/runner_safety_test.go
    - internal/codehealth/runner_test.go
    - internal/codehealth/storage_failure_test.go
    - internal/codehealth/window.go
    - internal/codehealth/window_test.go
    - internal/data/write_test.go
    - internal/doctor/checks.go
    - internal/doctor/repairs.go
    - internal/doctor/repairs_test.go
    - internal/init/agent_skills_test.go
    - internal/init/skill_validation_test.go
    - internal/init/upgrade.go
    - internal/init/upgrade_test.go
    - main_health_test.go
  dependencies: [go.mod, go.sum, package.json]
issues: [I-115, I-116, I-117, I-118, I-119, I-120, I-121]
supersedes: null
---

# C-957: O-032 Full Objective Check

## Result and authority

NEEDS WORK. Initial independent Full Check in a fresh session; this conversation built none of the work reviewed. All fourteen T-085–T-098 Tasks are owner-done with explicit board-owner waivers. A waiver does not replace this integration Check. No Task/Objective status, owner acceptance, Design or implementation was edited. Main checkout (git-dir equals common-dir), not a lane. Owner pasted `NEXT: Check O-032`; the router still says task with no selected Task. Following that explicit selection is authorized; the router was left byte-identical.

Seven findings are reported together. Native Windows evidence is unverified, one Windows skip violates CFG-03, two history-window cases violate T-098, release copy overstates two limits, stream output failures violate T-088, and synthetic URL credentials reach permanent evidence despite T-091's safety requirement. Direct Issue repairs and a fresh independent recheck are needed; completed Tasks stay done.

## Frozen initial scope and workflow lock

# O-032 initial independent Full Check scope lock

Session check-o032-20261002; baseline 67171b0; head 136fe69 plus the supplied working tree. No executor conversation inherited. All T-085–T-098 outcomes, Done When bullets, O-032 success conditions and confirmed observable outcomes are in scope. Applicable FS, ARCH, CFG, TEST, TPL and REL rules named by Tasks are required; STYLE rules advisory. Fresh make ci includes make test-full. Native Windows success is required evidence, cross-build is insufficient. Version bump is owner-deferred until merge. No implementation repair, Task closure, publication, provider install or new feature.

1. Public entry points: diagnostic name/repair lookup and sentinel matching; test-stream decode/aggregate/output and report publication; board load reducer and selection restoration; discovery/setup proposals; Store config/snapshot/report read/write, PlanPrune/Prune; ExecRunner/Collect/DefaultReaders; LoadWindow/LoadDashboard, series/trend/spark/chip/report; init/upgrade/migrate preservation; main health check and official snapshot reference validation; archive/checksum/version/npm wrapper; README/CHANGELOG/Design.
2. Files are the production/test diffs since baseline plus untracked window and benchmark files, the existing directly relied-on collector/storage/reader/history/render/gate/runtime orchestration, and Task contexts. No unrelated feature or dependency internals. Source and gate input hashes will be recorded before publication.
3. Supported saved JSON includes serializer output and other valid field orders, duplicate fields as accepted by the existing full decoder; missing/malformed/wrong typed/oversized fields must fail or take the defined fallback. Initial review may amend a factual scope error explicitly before verdict.
4. Mandatory matrix below; each row is crossed with each listed cell. N/A cells require the reason listed here. Each criterion additionally gets independent scenario evidence beyond running unchanged tests.

| Row | Applicable cells | N/A |
|---|---|---|
| D doctor | every exact mapping/default; each sentinel; wrapped/joined; ordered pair overlaps; empty/unknown | numeric/output sink axes: pure strings, no numeric state |
| B build stream | valid/malformed/mixed events; exact/below/above scanner bound; tee/read/writer/child failure; interrupt/early exit; report old/new/temp cleanup; sorted output | network redirects: local subprocess |
| L reload | first/success/fail/retry/repeat; stale/equal/unsequenced; Goal/Objective/cursor/filter/detail/issue/status/rollback; initial vs reloaded | terminal Unicode layout unchanged by this refactor |
| S store | versions -1/0/1/2; absent/empty/malformed/trailing/wrong type; regular/symlink/directory; exact/above byte limits; same/different concurrent saves/read replacement; obstruction before/after temp; retry; manual 9/10/11; official preservation/ties; prune failure | global transaction unsupported; no public pruning command |
| E runner/collector | configured vs actual target/argv/cwd; unavailable/start/error/exit 0/1/127; stdout and stderr cap -1/0/+1 and overflow; control/UTF8/bounded reason/no environment/report leakage; deadline/cancel/orphan/tree/retry; temporary and project report cleanup; cancel before save; report failure after save | redirects not handled by Savepoint; tools own network; detached descendants outside tree guarantee |
| P readers/gate | all 9 providers x measured/partial/malformed/missing/boundary; polyglot independence; scope/version/rename/sibling; manual/official; required/optional/blocking; report stale; full Check reference missing/manual/official | provider install/network outside normal tests |
| W history window | official counts 0–11, >10; leading/trailing/interleaved manuals, manual-only; equal timestamps/identity ordering; valid reordered/duplicate headers; invalid head/body inside/outside, filename mismatch; threshold equivalence full vs window; all rows/trends/baseline/spark/history/chip/report/sign-off; sparse/changing series; no writes | Unicode width not changed; renderer existing text-class suite covers unchanged layout |
| A adoption | fresh/upgrade/repeat/dry-run; authored config/Design/snapshot/report/ignore/skill; conflict sidecar; schema 1/2/unknown; partial failure recovery | downgrade conversion is explicitly unsupported |
| R distribution | six targets; single member/checksum/version; provider-free PATH; wrapper; invalid archive/version/member; native Windows current tree | publication out of scope |
| F performance/docs | 10/100/1000 single/multi/heavy; load/project/report/freshness vs fixed frame normal/history/multi-instance; benchmark repeatability; 5 signals/9 reports/setup/check/report/H/h/R/privacy/limitations/retention/measurements/followups; Design module contracts | universal machine speed budget not promised |

5. Workflow/effect inventory and oracle:

| Order/workflow | Effect | Failure timing / final state / cleanup | Oracle |
|---|---|---|---|
| Collect validate → full history → repository → per-instance report prep/run/read/classify → cancellation guard → snapshot link → report | tools; temporary reports; one immutable snapshot; derived report | pre-save failure/cancel no snapshot; instance failures contained; report failure warning after snapshot survives; temp cleanup independent | byte/ID inventory, fake runner calls, process liveness, direct values |
| Save config/report encode → validate paths → temp write/sync/close/chmod → rename | atomic replacement | old bytes until replacement; failed temp/rename cleans temp; repeat no mtime change; owner ignore preserved | retained bytes/mtime and independently parsed results |
| Save snapshot validate/encode → existing check → temp → exclusive link → remove temp | create-only evidence | concurrent identical one creator; different identity independent; conflict no overwrite | IDs and bytes, directory inventory |
| Prune load/validate all → plan → oldest-manual removes | only older manuals removed | any invalid file prevents removal; removal failure returns removed prefix and retained suffix; no official removed | origin/time-sorted independent inventory |
| Dashboard config → list → head → sort → window decode → projection/report presence | no write or process | bad name/head/window fails; old body damage intentionally unseen; alternate field order fallback full decode | full decoder ordering and expected ID sets; before/after file hashes |
| Build stream pipe → start → read/tee → wait → failures/stderr/timing → validate stream → publish reports | subprocess and atomic reports | read/tee/child failures retain previous complete report, drains after overlong token | independent event/output expectations and report bytes |
| Upgrade schema/manifest → compare managed assets → changed write/sidecar → manifest | user content retained | dry run no write; conflict retained plus incoming sidecar; named partial recovery; repeat sidecar unchanged | byte and mtime capture |
| Distribution compile → archive → checksums → verification → package → native smoke | local build artifacts only | mixed/malformed archive refuses verification | tar member inventory and sha256sum; binary version |

Admission: reproducible supported-path violation of a named criterion or guardrail inside a row above with a credible consequence. Unverified native evidence stays unverified, never translated into a claim the code is defective. Finish every row even after findings.


Initial factual amendment after the reader probe: an empty Lizard CSV is a supported no-functions report under the owner's recorded I-101 accepted limitation, not malformed data. The initial harness incorrectly treated it as invalid; its expected outcome was corrected to zero plus the explicit cannot-detect-unparsed-files reason and the complete affected reader matrix rerun. No new blocking axis was added. Sparse-series trend shortening is explicitly permitted by T-097; it is a copy finding only. Error-URL credential probe belongs to the original E secret-safe external-failure cell and T-091 criterion 2.

## Completed coverage matrix

Cells below classify the entire applicable inventory in the lock; N/A cells keep the reasons in that lock. A passing gate is supporting evidence, not a substitute for these outcomes.

| Frozen row | Cell classification and concrete evidence |
|---|---|
| D doctor | Passed every old mapping/default and every typed sentinel/ordered overlap: independent baseline oracle from 67171b0 compares 42 message terms pairwise, 57 mapping keys, 40 sentinels pairwise crossed with five refinement strings. TestO032DoctorOracleMatrix passes; production lookup data retains original outputs and precedence. |
| B stream | Normal decode/aggregation/order passed independent values; malformed/oversized/read/tee/process/interrupt/early-exit and report preservation pass named TestRunGoTestStream*, TestRunGoTestWithReports*, TestGoTestStreamComplete in main_test.go/full gate. Diagnostic/timing output failure and primary child failure under a broken sink are Issue I-120. Pipe draining is present and oversized test terminates. No reports written by focused overlay probes. |
| L reload | Independent baseline oracle crosses initial/reloaded x last sequence 0/4 x incoming 0/3/4/5 x success/failure x retry/nonretry x filter empty/valid/missing (192 cases), compares visible output, selection/cursors/status/diagnostics/fatal and next-command presence: passes. Exact preservation tests include TestFailedReloadKeepsLastGoodBoardAndRecovers, TestFailedReloadRetriesBeforeReporting, TestOlderLoadResultDoesNotOverwriteNewer, TestReleaseReloadPreservesFocusAndDiagnosesRemovedSelection, TestReloadRefreshesAnOpenDetailAndClosesADeletedOne and TestUpdateUnknownObjectiveFilterQuitsWithError. |
| S persistence | Version -1/0/1/2, invalid types/trailing/oversize/path/symlink/conflict, concurrent same/different immutable saves, complete old/new replacements, temp/rename obstruction and repeat preservation pass named storage_failure_test.go and storage_test.go tests. Scoped race run passes. Independent PlanPrune/Prune count matrix manuals 0–12 plus four officials proves dry plan, exact removed counts, official ID preservation, retained counts and repeat. Partial removal behavior inspected: return removed prefix plus untouched suffix; existing TestPrunePartialFailureKeepsRemainingValid proves refusal at zero removals. Deterministic after-first-remove OS failure not separately injected; no claim of a new test for it. Windows replacement cell is Issue I-116 and native result unverified I-115. |
| E external boundary | Configured executable/argv/cwd actual match; no shell. Existing controlled helper TestExecRunnerPassesArgumentsAndDirectoryUnchanged, TestExecRunnerStopsTheWholeTreeOnDeadlineAndRunsAgain, TestExecRunnerBoundsPipesHeldByAnOrphan, TestCollectRealProcessTimeoutEndsDescendantsAndLaterInstanceRuns, TestCollectCancelWinsOverTheInstanceDeadline and TestCollectCancelledSavesNothing pass. Stdout 32 MiB -1/exact/+1 and stderr 64 KiB -1/exact/+1/8x, file-report 32 MiB bounds pass; invalid UTF-8/control sanitization traced to cleanLine and covered by TestExecRunnerResults. Available/failed/unavailable/timeout remain distinct; report failure after save warning tests pass. The official sandbox run reproduces actual OSV DNS/socket refusal as failed/unmeasured; network retry produces measured findings. URL-secret failure output is Issue I-121; native Windows cleanup unverified I-115. HTTP redirects/retries belong to providers, not Savepoint's direct API; no network client/install added. |
| P provider/Check integration | All nine DefaultReaders registered and tested. Existing all-measured and mixed failing Go/JS/Python collection cases pass TestCollectWithRealReadersKeepsEveryMeasureTruthful and TestPolyglotReadersKeepEachStackIndependent; scope/version/rename/sibling/manual/official and required/optional/sign-off pass TestPolyglotHistoryComparabilityAcrossChanges and TestPolyglotDashboardAgreesWithTheVerdict. Independent 9-provider x 5-input (empty, brace, null, array, empty object) matrix passes after recording I-101's accepted empty-CSV limitation. No failed report gains a healthy zero. Partial may carry incomplete values. Strict official-only references pass health_snapshot_refs_test.go and healthcheck_test.go. Jscpd no-per-file counts retains only report global totals; no filtered-evidence denominator fabricated. |
| W window | Independent official counts 0–11 x leading/trailing/interleaved manuals: 35/36 pass, exactly ten plus leading manuals is I-117. Repeated timestamp header violates decoder ordering I-118; repeated-origin tested too. Compact/reordered/512-byte-prefix and malformed/wrong-type/trailing/symlink cells pass. Existing window size/tie/body damage inside/outside/name/head tests pass; fallback full decode and no filesystem writes verified in code. Full/dash newest identities and rows within window covered, >10 continuous-series trend/baseline/spark/history/sign-off pass existing TestLoadDashboardBeyondTheWindowKeepsTrendsAndSaysOlderWasNotRead. Sparse-series shortening matches approved policy; public promise is I-119. Manual-only sign-off remains no-official. |
| A adoption | Full temporary-project init/upgrade/legacy/schema/conflict/partial-recovery suite passes. Exact new preservation test TestMainUpgradeAssetsPreservesAdoptedHealthAndAuthoredContent retains config/snapshots/reports/Design/unmanaged bytes, dry-run and repeat mtimes, sidecar behavior. Fresh TestMainInitEndsWithHealthPreviewAndWritesNoHealthConfig passes. Upgrade schema refusal cases and preview-first migration cases retain bytes. Sidecar repeat fix is a narrow pre-write equality check; no managed ownership change. Four skills/three references byte identity passes TestV2SkillSetIsCompleteWithByteParity. |
| R release | Fresh make ci passes full uncached host suite, six cross-builds, archive verification, npm wrapper and pack inventory. Independent sha256sum verifies all six, smoke-test v2.0.5 passes. Source/package boundary inspected: no bundled providers, shell execution/install added; current Task's provider-free PATH transcript is consistent with default unavailability API tests and one-member archives. Native Windows required evidence unavailable, I-115; cross-compilation is explicitly not native runtime proof. |
| F performance/docs | Filesystem load benchmark single/heavy n=100/1000 repeated 3x: medians 2.50/13.06 ms and 22.77/34.22 ms, single1000 4.87 MB/91.1k allocs, heavy1000 about 23.84 MB/148.9k; matches Task reference-host budgets. Heavy n=100 ~20.21 MB/72.3k; additional file overhead is linear in count, not snapshot bytes. Fixed-frame/render cases run independently, no machine-speed unit gate. All 5 signals/9 report formats, manual/official, H/h/R, report handoff, privacy/provider network/env/no installation, retention/deferred cleanup and follow-ups covered in README/CHANGELOG. Unconditional sparse-history/partial-measurement copy is I-119. |

## Acceptance coverage

Numbers are Done When bullets in their recorded order. Proven means independent scenario evidence plus code/contract review supports the rule; advisory style is evaluated separately.

| Task | Per-criterion classification |
|---|---|
| T-085 | 1–4 Proven: diagnosis/control/failure transcripts give bounded inference rather than a claim about lost historical stderr; actual official sandbox failure and successful network retry confirm environmental distinction. Tail repair is owner-directed and bounded. Credential risk remains explicitly covered by I-121 against T-091's stronger acceptance promise. |
| T-086 | 1–4 Proven: separate whole/production/test/maintained denominators, owner-directed config/default-exclusion changes and trend restart are documented; scoped jscpd official result has its own numerator/denominator. Intentional archive/template duplication is preserved. Extra test helper refactors preserve meaningful assertions; thresholds remain unchanged. |
| T-087 | 1–4 Proven: complete independent old/new diagnostic matrix passes; focused responsibility/CCN evidence coherent with source and official complexity evidence. No new diagnostic vocabulary. |
| T-088 | 1,3,4 Proven; 2 Issue I-120: subprocess/stream/report completion retained, but supported diagnostic writer failure is ignored or masks primary child failure. |
| T-089 | 1–4 Proven: 192-case independent old/new load-state matrix and named existing selection/detail/issue/rollback tests pass, helper responsibilities remain local and no IO moved into rendering. |
| T-090 | 1,3–5 Proven (after-first-removal fault test limitation stated above); 2 Issue I-116 for required Windows coverage. Linux bytes/concurrency/race/preservation and no-ID/version-change evidence pass. |
| T-091 | 1,3,4 Proven on Linux; 2 Issue I-121; 5 Unverified I-115. Output/process bounds and distinct instance failures pass, but credential-bearing error output persists and native Windows evidence is missing. |
| T-092 | 1–5 Proven through temporary-project adoption/migration/full gate and managed-skill identity. No auto-config or authored-content overwrite introduced. |
| T-093 | 1–5 Proven through nine-provider/mixed-stack/malformed matrix, isolation, history/gate/reference and report-warning checks, with the recorded owner-accepted I-101 limitation. |
| T-094 | 1–4 Proven: deterministic benchmarks and named bounded decision with host/sizes/limitations; independent reference measurements and fixed-frame test support the findings. Budgets remain host-specific proposals. |
| T-095 | 1–3,5 Proven; 4 Unverified I-115. Full local release/package evidence current; no native Windows head run. Version 2.0.5 is expected under recorded owner deferral, not a mismatch finding. |
| T-096 | 1 Issue I-119 (history promise); 2 Issue I-119 (partial evidence); 3–5 Proven for other documented content, boundaries, packaging and implemented Design reconciliation. Owner copy/walkthrough review remains the owner's decision, not inferred here. |
| T-097 | 1–6 Proven as a research/decision deliverable: alternatives, per-series bounds, sign-off/labels, damage boundary and follow-up Task specified without changing identity/retention. Its implementation claim of identical small history is independently evaluated under T-098, not assumed. |
| T-098 | 1 Issue I-118 (ordering); 2 Issue I-117 (<=10-official equivalence); 3–7 Proven with continuous comparable signals, named old-body boundary, unchanged full readers, reproducible budget and no cache/write/identity/retention change. |

O-032 success conditions 1 (adoption), 4 (packaging boundary), 7 (release notes content apart from inaccurate limits) and 8 (bounded diagnosis/refactor/scope evidence) are supported by the rows above. Conditions 2–3 have platform/diagnostic findings, 5 has inaccurate history/partial guidance, and 6/confirmed outcome 5 cannot be cleared until the complete integration findings are repaired and native evidence is obtained. Confirmed outcomes 1 adoption and the Linux parts of 2 persistence/process and 3 performance/distribution are proven; 3 native Windows is Unverified and 4 honest public guidance has I-119. No broader analyzer, cleanup CLI, provider plugin or publication requirement invented.

## Gate and Code Health evidence

Fresh `make ci` on 2026-10-02, go1.26.2 linux/amd64, AMD Ryzen 7 7800X3D/WSL2, Node v22.22.2/npm 10.9.7: exit 0. Includes `make test-full` (-count=1 host tests and six target builds), build, dist verification and package-check. Initial sandbox attempt passed full tests/build/distribution but npm failed EROFS writing its cache; explicit escalation rerun passed. No code change between attempts. 3,268 tests pass, zero skipped reported by the saved health Go-test reader; all 14 packages pass. All six archives are v2.0.5, linux/darwin/windows x amd64/arm64; archive verification enforces one executable. `sha256sum -c checksums.txt` from dist succeeds for all six. `make smoke-test`: v2.0.5, exit 0. Native Windows CI gap is I-115, not hidden by this gate.

`git diff --check` returns 2 solely for a pre-existing extra blank EOF line at T-094:147. Advisory metadata whitespace observation, no program/gate failure. Check/Issue writes contain no new whitespace errors. `go test -race ./internal/codehealth -run 'TestConcurrentSnapshotSavesAndReads|TestReadersSeeCompleteConfigAndReportDuringReplacement|TestFailedReplacement|TestFailedWrites|TestPrune' -count=1` passes (1.817s package time).

After the full gate, `./savepoint health check O-032` saved official snapshot sha256:7887dae880e2b2c255061957fd15929a949502af21bef6c71981c0011e7a8345. OSV failed under the sandbox (DNS socket operation not permitted); the instance is failed/unmeasured, not zero vulnerabilities. Escalated network rerun saved the cited official snapshot at 2026-10-02T05:07:22Z and exits 0: Code Health does not block clearance. Snapshot/report writes are the explicit Full Check exception, not agent setup/report commands.

Latest measured results: tests 0 failures/3,268 passing; coverage 86.5687% (11,492/13,275); complexity 42 maximum CCN across runtime plus tests (33 functions >20); maintained duplication 2.9224% (2,889/98,857 lines, 320 clones/sources); vulnerabilities 2 groups, both unknown severity (go.mod: golang.org/x/sys 0.38.0 GO-2026-5024, golang.org/x/text 0.3.8 GO-2026-5970). Zero high/critical severity was reported, but unknown severity is not assurance of safety. Advisory triage remains separate from this Check's findings. No provider/dependency upgrade or installation performed; exact vulnerability reachability/severity was not investigated in this release check. Tail-preserving failure output worked in the actual sandbox refusal.

## Materiality

| Issue | Likelihood | Impact | Materiality | Recommendation |
|---|---|---|---|---|
| I-115 | High: no current run exists | Medium: supported-platform evidence missing | High for release readiness | Obtain native current-tree CI evidence before clearance. |
| I-116 | High: unconditional Windows skip | Medium: live reader/replacement failure behavior untested | Medium | Replace skip with preservation/refusal evidence; combine with I-115 validation. |
| I-117 | Medium: manual refresh commonly precedes checks | Low: basis/count wording changes at exact boundary | Low | Narrow window-boundary correction and regression. |
| I-118 | Low: repeated fields require an alternate saved representation | High: wrong newest evidence can hide current results | Medium | Make header semantics and full decoder agree or refuse ambiguous input. |
| I-119 | Medium: removed/readded instances and partial readers supported | Medium: readers are promised evidence the window cannot show | Medium | Narrow copy correction to approved limits. |
| I-120 | Low: broken/closed diagnostic sink | Medium: failure diagnostics lost and primary child error masked | Medium | Preserve primary errors and report output failures; retain atomic report behavior. |
| I-121 | Low: provider must emit URL credentials | High: permanent local evidence can retain a secret | High for privacy | Bounded URL-credential protection before persistence; explicit owner exception required to accept unmet criterion. |

## Design reconciliation and observations

Architecture retains single-binary, project-owned providers, five signals/nine readers, read-only renderers, local snapshots, no telemetry/install, full-history collection/prune and derived report warning after snapshot save. Implemented head-based window and stderr-tail changes are reconciled in the current Design tools bullet and AGENTS module map. Exact project tool arguments stay in health config. Modular boundaries are intact. No dependency or schema-version change introduced by O-032. Most complexity improvement is production refactoring; residual aggregate CCN includes test functions and is not itself an opt-in health blocker.

Nonblocking observations: I-101's accepted no-functions/unparsed Lizard CSV limit remains visible; do not reopen it from the initial erroneous empty-report oracle. T-093's JscpdReader doc comment still describes partial/no-value when global totals are retained; useful future copy cleanup, not a new release defect finding. Candidate A's fallback full-decodes noncanonical heads, and manual-only/unusually dense manual history can still load unbounded bodies; those are explicit decision trade-offs, not an invented universal speed guarantee. Older corrupted bodies outside the window are intentionally unseen; full-history paths continue to refuse them. SnapshotLabels still full-loads history, as expressly deferred. T-094's extra EOF blank line and transient root report files present before review were not repaired. Saved health advisories need separate owner triage; neither optional warning nor lack of a Good label automatically created an Issue.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [ ] STYLE-03 **Test branches** — missed leading-manual, repeated-header, broken-sink and credential-output cases; native Windows replacement skipped (window_test.go, main_test.go, runner_safety_test.go, storage_failure_test.go).
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [ ] STYLE-06 **Handle errors at boundaries** — stream diagnostic errors ignored/mask primary error (buildtool/main.go:270,300,318).
- [x] STYLE-07 **One source of truth**
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs**

Style is advisory. The verdict follows named acceptance/guardrail findings and missing native evidence, not unchecked style boxes.

## Owner handoff

No Task/Objective closure or inferred exception. Repair I-116–I-121 directly under their Issues; obtain I-115 native evidence for the resulting final tree. Recheck this same frozen matrix and original reproductions in a fresh session, superseding C-957. Owner acceptance and Task/Objective completion remain owner decisions. No router selection or Goal bytes changed.

## Reproducibility appendix

Review probes were Go overlay files in /tmp, not repository source/test modifications. Commands: `go test -overlay /tmp/o032-overlay.json ./internal/codehealth -run '^TestO032' -count=1 -v`, analogous doctor/board/buildtool packages; separate final codehealth run selected `^TestO032(Sparse|Prune|Reader|Provider)`. Initial syntax/field-name mistakes in scratch harnesses were corrected before completed results; they are not product findings. Production tests/full gate ran without overlays. Appendices below preserve the corrected harness and oracle generators; use absolute target paths ending in `internal/<package>/o032_probe_test.go` mapped to the scratch files in a Go `Replace` overlay. Doctor/board baseline functions come from immutable 67171b0, not the new rule tables.

Gate-input inventory: 416 files; SHA256 of sorted JSON path→content-hash manifest `ebc1df5d9db876a0f9dd84ab02d02d9201ddd2c9076b55333c1d8ceebd0ea949`. Scratch manifest `/tmp/o032-reviewed-inputs.json`. Code, tests, fixtures, dependencies, gate definitions and shipped assets were not edited by this checker. Reviewed dirty tree is listed above; head SHA alone does not describe it.

### Scratch file o032_make_oracles.py

```python
import json,re,subprocess
from pathlib import Path
root=Path.cwd()
def old(path):return subprocess.check_output(['git','show','67171b0:'+path],text=True)
def function(src,name):
 start=src.index('func '+name+'(')
 end=src.find('\nfunc ',start+5)
 if end<0:end=len(src)
 # Comments preceding the next function are harmless in the generated source.
 return src[start:end].replace('func '+name+'(', 'func oracle'+name+'(',1)
rep=old('internal/doctor/repairs.go'); chk=old('internal/doctor/checks.go')
terms=sorted(set(re.findall(r'strings.Contains\(msg, "([^"]+)"\)',rep)))
names=sorted(set(re.findall(r'case "([^"]+)":',rep)))
sentinels=sorted(set(re.findall(r'data\.(Err\w+)',chk[chk.index('func v2DiagnosticName('):])))
src='package doctor\nimport("errors";"strings";"fmt";"testing";"github.com/opencode/savepoint/internal/data")\n'
src+='\n'.join(function(rep,n) for n in ['SuggestRepair','V2ProblemRepair'])+'\n'+function(chk,'v2DiagnosticName')
src+='\nfunc TestO032DoctorOracleMatrix(t *testing.T){\n'
src+='terms:=[]string{'+','.join(json.dumps(x) for x in ['', 'unknown']+terms)+'}\n'
src+='for _,a:=range terms {for _,b:=range terms{ e:=errors.New(a+" "+b);if g,w:=SuggestRepair(e),oracleSuggestRepair(e);g!=w{t.Errorf("repair %q: %q vs %q",e,g,w)}}}\n'
src+='for _,name:=range []string{'+','.join(json.dumps(x) for x in ['', 'unknown']+names)+'}{if g,w:=V2ProblemRepair(name),oracleV2ProblemRepair(name);g!=w{t.Errorf("map %q: %q vs %q",name,g,w)}}\n'
src+='sentinels:=[]error{'+','.join('data.'+s for s in sentinels)+'}\n'
src+='for _,a:=range sentinels {for _,b:=range sentinels {for _,msg:=range []string{"", "Goal id", "release id", "scope names missing release", "quality gate"}{e:=fmt.Errorf("%s: %w",msg,errors.Join(a,b));if g,w:=v2DiagnosticName(e),oraclev2DiagnosticName(e);g!=w{t.Errorf("name %v: %q vs %q",e,g,w)};if g,w:=SuggestRepair(e),oracleSuggestRepair(e);g!=w{t.Errorf("typed %v: %q vs %q",e,g,w)}}}}\n}\n'
Path('/tmp/o032_doctor_probe_test.go').write_text(src)
overlay=json.loads(Path('/tmp/o032-overlay.json').read_text())
overlay['Replace'][str(root/'internal/doctor/o032_probe_test.go')]='/tmp/o032_doctor_probe_test.go'
Path('/tmp/o032-overlay.json').write_text(json.dumps(overlay))
print('Doctor oracle:',len(terms)+2,'message terms x pairwise;',len(names)+2,'mapping keys;',len(sentinels),'sentinels x pairwise x 5 refinements')

```

### Scratch file o032_more_oracles.py

```python
import json,subprocess
from pathlib import Path
root=Path.cwd()
old=subprocess.check_output(['git','show','67171b0:internal/board/v2/update.go'],text=True)
start=old.index('func (m Model) applyLoad(');end=old.index('\nfunc ',start+5)
oracle=old[start:end].replace(' applyLoad(', ' oracleApplyLoad(',1)
src='package v2\nimport("fmt";"reflect";"testing";tea "github.com/charmbracelet/bubbletea")\n'+oracle+'''
func TestO032ReloadOracleMatrix(t *testing.T){
 root:=writeValidProject(t)
 for _,loaded:=range []bool{false,true}{for _,last:=range []uint64{0,4}{for _,seq:=range []uint64{0,3,4,5}{for _,failure:=range []bool{false,true}{for _,retry:=range []bool{false,true}{for _,filter:=range []string{"","O-001","O-999"}{
  m:=openBoard(t,root,"");m.Loaded=loaded;m.lastLoadSeq=last;m.ObjectiveFilter=filter;m.preserveReloadStatus=true;m.StatusMessage="retained status"
  msg:=loadCmd(root)().(projectLoadedMsg);msg.Seq=seq;msg.Retry=retry;if failure{msg.Diagnostic="temporary record diagnostic"}
  a,ac:=m.applyLoad(msg);b,bc:=m.oracleApplyLoad(msg);g,w:=a.(Model),b.(Model)
  if (ac==nil)!=(bc==nil){t.Errorf("cmd mismatch")}
  gv:=[]any{g.Loaded,g.SelectedRelease,g.SelectedObjective,g.ObjectiveCursor,g.FocusedCard,g.Diagnostic,g.ReloadDiagnostic,g.StatusMessage,g.FatalErr,g.lastLoadSeq,g.DetailOffset}
  wv:=[]any{w.Loaded,w.SelectedRelease,w.SelectedObjective,w.ObjectiveCursor,w.FocusedCard,w.Diagnostic,w.ReloadDiagnostic,w.StatusMessage,w.FatalErr,w.lastLoadSeq,w.DetailOffset}
  if !reflect.DeepEqual(gv,wv)||g.View()!=w.View(){t.Errorf("loaded=%v last=%d seq=%d fail=%v retry=%v filter=%s changed behavior: %v vs %v",loaded,last,seq,failure,retry,filter,gv,wv)}
 }}}}}}
}
'''
Path('/tmp/o032_board_probe_test.go').write_text(src)
overlay=json.loads(Path('/tmp/o032-overlay.json').read_text());overlay['Replace'][str(root/'internal/board/v2/o032_probe_test.go')]='/tmp/o032_board_probe_test.go'
Path('/tmp/o032-overlay.json').write_text(json.dumps(overlay))

```

### Scratch file o032_window_probe_test.go

```go
package codehealth

import (
 "context"
 "encoding/json"
 "fmt"
 "os"
 "reflect"
 "strings"
 "testing"
)

func TestO032SparseHistoryDocumentation(t *testing.T) {
 store,root:=newProject(t);p:=historyProfiles[0];cfg:=historyConfig(p);store.SaveConfig(cfg)
 for n:=0;n<16;n++ {
  s:=historySnapshot(cfg,p,n);s.Origin=OriginOfficial;s.Retention=RetentionPermanent
  if n>=3&&n<15 {s.Results=nil;s.Summary.Capabilities=nil;s.Summary.Overall=ClassificationUnknown}
  s.ID=s.ComputeID();mustSave(t,store,s)
 }
 got,err:=LoadDashboard(root);if err!=nil {t.Fatal(err)};full:=fullDashboard(t,root)
 t.Logf("window trend=%q spark=%d; full trend=%q spark=%d",got.Rows[1].Trend,len([]rune(got.Rows[1].Spark)),full.Rows[1].Trend,len([]rune(full.Rows[1].Spark)))
 if got.Rows[1].Trend==full.Rows[1].Trend {t.Fatal("probe did not demonstrate the documented sparse-history limit")}
}

func TestO032PruneIndependentOriginCountMatrix(t *testing.T){
 for manuals:=0;manuals<=12;manuals++{t.Run(fmt.Sprint(manuals),func(t *testing.T){
  store,_:=newProject(t);p:=historyProfiles[0];cfg:=historyConfig(p);var officialIDs []string
  for n:=0;n<manuals+4;n++{s:=historySnapshot(cfg,p,n);s.Origin=OriginManual;s.Retention=RetentionPrunable;if n<4{s.Origin=OriginOfficial;s.Retention=RetentionPermanent};s.ID=s.ComputeID();mustSave(t,store,s);if n<4{officialIDs=append(officialIDs,s.ID)}}
  all,err:=store.LoadSnapshots();if err!=nil{t.Fatal(err)}
  plan,err:=store.PlanPrune();if err!=nil{t.Fatal(err)}
  afterPlan,_:=store.LoadSnapshots();if !reflect.DeepEqual(snapshotIDs(all),snapshotIDs(afterPlan)){t.Fatal("dry plan wrote")}
  wantRemoved:=max(0,manuals-10);if len(plan.Removed)!=wantRemoved{t.Fatalf("plan removed %d want %d",len(plan.Removed),wantRemoved)}
  result,err:=store.Prune();if err!=nil{t.Fatal(err)};if !reflect.DeepEqual(plan,result){t.Fatal("plan and actual differ")}
  remaining,_:=store.LoadSnapshots();if len(remaining)!=4+min(manuals,10){t.Fatal("wrong remaining")}
  for _,id:=range officialIDs{found:=false;for _,s:=range remaining{found=found||s.ID==id};if !found{t.Fatal("official lost")}}
  again,err:=store.Prune();if err!=nil||len(again.Removed)!=0{t.Fatal("repeat changed")}
 })}
}

func TestO032ReaderEmptyMalformedMatrix(t *testing.T){
 for key,reader:=range DefaultReaders(){for _,raw:=range []string{"", "{", "null", "[]", "{}"}{t.Run(string(key)+"/"+fmt.Sprintf("%q",raw),func(t *testing.T){
  reading,err:=reader.Read(context.Background(),ReportInput{Provider:key,Root:t.TempDir(),Data:[]byte(raw)})
  if key==ProviderLizardCSV&&raw=="" {
   if err!=nil||reading.Value==nil||reading.Value.Number!=0||!strings.Contains(reading.Reason,"cannot show whether any file failed to parse"){t.Error("accepted I-101 limitation lost")}
  } else if err==nil&&!reading.Partial&&reading.Value!=nil {t.Errorf("non-report input claims complete value: %+v",reading)}
 })}}
}

func TestO032ProviderErrorCredentialPersistence(t *testing.T){
 root:=project(t);const credential="o032-secret-for-probe"
 tools:=&fakeTools{t:t,behavior:map[string]func(context.Context,ToolSpec)(ToolResult,error){"provider":func(context.Context,ToolSpec)(ToolResult,error){return ToolResult{ExitCode:127,Stderr:sanitizeTail("proxy connection failed https://alice:"+credential+"@proxy.example.invalid")},nil}}}
 cfg:=cfgOf(lizardInstance("local","provider","a/**"))
 c:=collect(t,root,cfg,Readers{ProviderLizardCSV:okReader(1,UnitCCN)},tools)
 r:=collectedFor(t,c,CapabilityComplexity,"local").Result
 if strings.Contains(r.Reason,credential){t.Error("URL credential persisted in snapshot reason")}
 snaps,err:=NewStore(root).LoadSnapshots();if err!=nil{t.Fatal(err)}
 data,_:=json.Marshal(snaps);if strings.Contains(string(data),credential){t.Error("immutable on-disk snapshot contains credential")}
}

func TestO032WindowIndependentMatrix(t *testing.T) {
 for officials:=0; officials<=11; officials++ {
  for _, shape:=range []string{"leading", "trailing", "interleaved"} {
   t.Run(fmt.Sprintf("official=%d/%s", officials,shape),func(t *testing.T){
    store,root:=newProject(t); p:=historyProfiles[0]; cfg:=historyConfig(p)
    if _,err:=store.SaveConfig(cfg);err!=nil {t.Fatal(err)}
    n:=0
    save:=func(origin Origin){
     s:=historySnapshot(cfg,p,n); n++; s.Origin=origin
     if origin==OriginOfficial {s.Retention=RetentionPermanent} else {s.Retention=RetentionPrunable}
     s.ID=s.ComputeID(); mustSave(t,store,s)
    }
    if shape=="leading" {save(OriginManual);save(OriginManual)}
    for i:=0;i<officials;i++ {save(OriginOfficial);if shape=="interleaved" {save(OriginManual)}}
    if shape=="trailing" {save(OriginManual);save(OriginManual)}
    if n==0 {save(OriginManual)}
    all,err:=store.LoadSnapshots();if err!=nil {t.Fatal(err)}
    win,cut,err:=store.LoadWindow();if err!=nil {t.Fatal(err)}
    if officials<=10 {
     if cut || !reflect.DeepEqual(snapshotIDs(win),snapshotIDs(all)){t.Errorf("inside-window promise: officials=%d all=%d window=%d cut=%v",officials,len(all),len(win),cut)}
     got,err:=LoadDashboard(root);if err!=nil {t.Fatal(err)}
     full:=fullDashboard(t,root)
     if !reflect.DeepEqual(got,full){t.Errorf("dashboard differs with at most 10 officials: history %d vs %d; basis %q vs %q",len(got.History),len(full.History),got.Rows[0].Basis,full.Rows[0].Basis)}
    } else {
     if len(win)<10 {t.Errorf("window lost officials")}
    }
   })
  }
 }
}

func TestO032WindowDuplicateHeaderMatchesFullDecoder(t *testing.T) {
 for _, key:=range []string{"created_at","origin"} {
  t.Run(key,func(t *testing.T){
   store,root,_:=windowHistory(t,40); all,err:=store.LoadSnapshots();if err!=nil {t.Fatal(err)}
   newest:=all[len(all)-1]; path:=snapshotFileOf(t,root,newest)
   data,err:=os.ReadFile(path);if err!=nil {t.Fatal(err)}
   // Valid JSON using last-value-wins semantics of DecodeSnapshot. Original ID remains correct.
   prefix:=`"created_at":"2000-01-01T00:00:00Z",`
   if key=="origin" {prefix=`"origin":"manual",`}
   altered:=strings.Replace(string(data),"{","{"+prefix,1)
   if err:=os.WriteFile(path,[]byte(altered),0644);err!=nil {t.Fatal(err)}
   decoded,err:=DecodeSnapshot([]byte(altered));if err!=nil {t.Fatalf("full decoder rejected supported input: %v",err)}
   if decoded.ID!=newest.ID {t.Fatal("fixture changed identity")}
   got,err:=LoadDashboard(root);if err!=nil {t.Fatalf("window rejects a valid full-decoder snapshot: %v",err)}
   full:=fullDashboard(t,root)
   if got.SnapshotID!=full.SnapshotID {t.Errorf("newest identity differs: window %s full %s",got.SnapshotID,full.SnapshotID)}
   if !reflect.DeepEqual(got.History,full.History) {t.Errorf("history order differs for valid duplicate %s",key)}
  })
 }
}

func TestO032WindowRepresentationMatrix(t *testing.T) {
 for _,shape:=range []string{"compact","reordered","512prefix","wrongtype","trailing","symlink"} {
  t.Run(shape,func(t *testing.T){
   store,root,_:=windowHistory(t,8); all,_:=store.LoadSnapshots(); path:=snapshotFileOf(t,root,all[len(all)-1]); data,_:=os.ReadFile(path)
   valid:=true
   switch shape {
   case "compact": var b strings.Builder;_ = b; var v any;json.Unmarshal(data,&v);data,_=json.Marshal(v)
   case "reordered":var v map[string]json.RawMessage;json.Unmarshal(data,&v);data,_=json.Marshal(v)
   case "512prefix":data=append([]byte(strings.Repeat(" ",512)),data...)
   case "wrongtype":data=[]byte(strings.Replace(string(data),`"origin": "official"`,`"origin": 1`,1));valid=false
   case "trailing":data=append(data,']');valid=false
   case "symlink":target:=path+".outside";if err:=os.WriteFile(target,data,0644);err!=nil {t.Fatal(err)};if err:=os.Remove(path);err!=nil {t.Fatal(err)};if err:=os.Symlink(target,path);err!=nil {t.Fatal(err)};valid=false
   }
   if shape!="symlink" {if err:=os.WriteFile(path,data,0644);err!=nil {t.Fatal(err)}}
   _,_,err:=store.LoadWindow();if (err==nil)!=valid {t.Errorf("shape %s valid=%v err=%v",shape,valid,err)}
  })
 }
}

```

### Scratch file o032_buildtool_probe_test.go

```go
package main
import("bytes";"errors";"os/exec";"strings";"testing")
func TestO032StreamOutputFailureMatrix(t *testing.T){
 t.Run("success child broken timing sink",func(t *testing.T){
  if err:=runGoTestStream(streamHelperCommand("ok"),failingWriter{},nil);err==nil {t.Error("output writer failed but stream returned success")}
 })
 t.Run("child failure plus broken diagnostic sink",func(t *testing.T){
  err:=runGoTestStream(streamHelperCommand("failed-package"),failingWriter{},nil)
  var child *exec.ExitError;if !errors.As(err,&child)||child.ExitCode()!=1{t.Errorf("primary child exit lost: %v",err)}
 })
 t.Run("normal event independent totals",func(t *testing.T){
  summary:=newTestStreamSummary();var out,tee bytes.Buffer
  raw:="invalid-line\n"+`{"Action":"output","Package":"b","Output":"detail"}`+"\n"+`{"Action":"fail","Package":"b","Elapsed":1}`+"\n"+`{"Action":"pass","Package":"a","Test":"T","Elapsed":0.5}`+"\n"
  if err:=consumeTestStream(strings.NewReader(raw),&out,&tee,summary);err!=nil{t.Fatal(err)}
  if tee.String()!=raw||out.String()!="invalid-line\n"||len(summary.failedPackages)!=1||len(summary.packages)!=1||len(summary.tests)!=1 {t.Fatal("independent aggregation oracle mismatch")}
 })
}

```

### Independent command outputs

```text
=== RUN   TestO032DoctorOracleMatrix
--- PASS: TestO032DoctorOracleMatrix (0.03s)
PASS
ok  	github.com/opencode/savepoint/internal/doctor	0.028s
```

```text
=== RUN   TestO032ReloadOracleMatrix
--- PASS: TestO032ReloadOracleMatrix (0.20s)
PASS
ok  	github.com/opencode/savepoint/internal/board/v2	0.207s
```

```text
=== RUN   TestO032StreamOutputFailureMatrix
=== RUN   TestO032StreamOutputFailureMatrix/success_child_broken_timing_sink
    o032_probe_test.go:5: output writer failed but stream returned success
=== RUN   TestO032StreamOutputFailureMatrix/child_failure_plus_broken_diagnostic_sink
    o032_probe_test.go:9: primary child exit lost: write failed package output: disk full
=== RUN   TestO032StreamOutputFailureMatrix/normal_event_independent_totals
--- FAIL: TestO032StreamOutputFailureMatrix (0.00s)
    --- FAIL: TestO032StreamOutputFailureMatrix/success_child_broken_timing_sink (0.00s)
    --- FAIL: TestO032StreamOutputFailureMatrix/child_failure_plus_broken_diagnostic_sink (0.00s)
    --- PASS: TestO032StreamOutputFailureMatrix/normal_event_independent_totals (0.00s)
FAIL
FAIL	github.com/opencode/savepoint/internal/buildtool	0.007s
FAIL
```

Selected window results: exactly-ten-leading manuals changes dashboard basis; repeated created_at hides newest ID. Sparse-history output: window `No trend yet`, zero spark points; full `Steady over 4 official checks, 90% to 92%`, four spark points. Corrected 45-cell reader matrix and thirteen retention counts pass; synthetic credential test reports `URL credential persisted in snapshot reason` and `immutable on-disk snapshot contains credential`. No real credential used. Full raw logs remain in /tmp/o032-*.log.
