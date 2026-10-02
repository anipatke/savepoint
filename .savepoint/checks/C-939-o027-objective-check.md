---
id: C-939
scope: {kind: objective, id: O-027}
result: NEEDS WORK
checked_by: {role: checker, session: o027-full-check-20261001}
executed_session: o027-task-execution-sessions
checked_at: '2026-09-30T21:34:30Z'
reviewed:
  base_commit: 9fdc11f
  head_commit: 86ab71a67cc697901d4c4b9c3081d37472550458
  files:
    - internal/codehealth/classification.go
    - internal/codehealth/classification_test.go
    - internal/codehealth/config.go
    - internal/codehealth/errors.go
    - internal/codehealth/history.go
    - internal/codehealth/identity.go
    - internal/codehealth/identity_test.go
    - internal/codehealth/model.go
    - internal/codehealth/model_test.go
    - internal/codehealth/primitives.go
    - internal/codehealth/repository.go
    - internal/codehealth/repository_test.go
    - internal/codehealth/snapshot.go
    - internal/codehealth/storage.go
    - internal/codehealth/storage_test.go
    - internal/codehealth/testdata/valid-config-v1.json
    - internal/codehealth/testdata/valid-snapshot-v1.json
    - AGENTS.md
    - .savepoint/objectives/O-027-code-health-model-and-history/Objective.md
    - .savepoint/objectives/O-027-code-health-model-and-history/tasks/T-051-give-health-evidence-one-precise-vocabulary.md
    - .savepoint/objectives/O-027-code-health-model-and-history/tasks/T-052-keep-health-history-safe-and-compact.md
    - .savepoint/objectives/O-027-code-health-model-and-history/tasks/T-053-describe-repository-state-without-guessing.md
    - .savepoint/objectives/O-027-code-health-model-and-history/tasks/T-054-explain-health-trends-without-false-precision.md
  dependencies: [go.mod, go.sum, Makefile]
issues: [I-085, I-086, I-087, I-088, I-089, I-090, I-091]
supersedes: null
---

# C-939: O-027 Full Objective Check

## Result and independence

**NEEDS WORK.** Seven reproducible findings affect trustworthy records, repository identity, and native Windows verification. Every owned Task was reviewed, including its recorded owner waiver. This fresh checker conversation did not implement these Tasks. Executor provenance spans the separate T-051/T-052/T-053/T-054 sessions described in their evidence; the required executed_session field uses a collective label for those prior execution sessions, not an authenticated individual session ID.

Scope includes the uncommitted classification.go, history.go and classification_test.go, plus the owner's existing Objective, Task, router and AGENTS edits. HEAD alone is not the complete reviewed tree; source hashes below identify the actual inputs. No production code, existing tests, Task statuses, router selection, or Design were changed by this Check.

The supplied “Check O-027” line is the owner's selection. Router state remained task with O-027 selected and no Task/Issue; this does not override the explicit Check request. This checkout is not a worktree lane (git-dir and common-dir both .git).

## Findings and materiality

| Issue | Finding | Likelihood | Impact | Materiality | Recommendation |
|---|---|---|---|---|---|
| I-085 | Positive high/critical counts can produce Good when total is zero; invalid severity counts persist | Medium: malformed adapter evidence | High: hard blocker hidden | High | Fix count validation and conservative assessment together |
| I-086 | Missing/null measured numbers decode as zero while original ID and Good summary remain valid | Medium: incomplete records | High: missing evidence becomes passing | High | Fix required-number presence validation |
| I-087 | Trailing closing delimiters/garbage accepted by both JSON decoders | Medium: corrupt files | Medium: malformed history accepted | Medium | Require end of input after one record |
| I-088 | Evidence order with equal path/line changes canonical identity | Low: multiple notes at one location | Medium: duplicate identities or failed identical retry | Medium | Define full evidence ordering |
| I-089 | Invalid scope globs silently suppress files | Medium: configuration typo | Medium: false unchanged repository state | Medium | Diagnose malformed patterns |
| I-090 | POSIX colon filenames vanish from input and dirty detection | Low: unusual supported filename | Medium: changed source reported clean | Medium | Fingerprint or explicitly reject with a diagnostic |
| I-091 | Unchanged unusual-filename test fails on native Windows | High: every Windows run | Medium: required platform gate fails | High | Fix fixture selection and rerun native CI |

Each linked Issue contains exact requirements, source locations, minimal reproduction, actual/expected behavior and missing test evidence. All remaining matrix rows were completed after the first finding. No finding is admitted solely for style.

## Acceptance coverage

Criterion numbers follow each Task's Done When list.

| Criterion | Classification | Independent outcome and evidence |
|---|---|---|
| T-051.1 fixed vocabulary/duplicate instances | Proven | Fixed five capability list; TestCapabilitiesAreTheFixedFive, TestConfigValidation, malformed mixed/null list shapes |
| T-051.2 orthogonal collection/freshness and measured values | Issue | Existing 9×3 outcome/freshness matrix passes; M1_missing_numbers exposes missing measurement becoming zero (I-086) |
| T-051.3 explicit versioned validation | Issue | Version/origin/repository/provenance bounds covered; I-085/I-086/I-089 expose accepted impossible/missing/config values |
| T-051.4 canonical/series identities | Issue | Compatibility-component and commit-exclusion tests pass; evidence-tie permutation changes identity (I-088) |
| T-051.5 reject unsafe/malformed/incompatible records | Issue | Evidence path/unit/timestamp/origin tests pass; trailing JSON and number presence fail (I-086/I-087) |
| T-051.6 valid and invalid fixtures | Issue | Named model tests pass but do not cover newly reproduced invalid numeric/JSON cases; embedded M1 harness fails |
| T-052.1 confined strict persistence | Issue | Symlink/traversal/name/schema tests pass on Linux; loaders inherit malformed JSON acceptance (I-087) |
| T-052.2 immutable collision-safe retry | Proven | TestSaveSnapshotIsIdempotent, ConflictKeepsExistingFile, InterruptedWriteLeftoverIsIgnoredAndRetrySucceeds; create-only hard-link path and temp cleanup traced; integration saves/load identity |
| T-052.3 deterministic series-preserving load | Proven | TestLoadSnapshotsOrderIsDeterministic and IncompatibleSeriesStayApart; integrated stored official/manual history remains separate |
| T-052.4 retain official/newest ten manual | Proven | Independent 0–12 manual sweep with three official records; exact removed counts and official identities checked; existing tie and 14-record tests pass |
| T-052.5 explicit pruning and failure safety | Proven | No pruning in load/save call paths; PlanPrune read-only, repeat no-op; malformed/unwritable store tests preserve records |
| T-052.6 persistence boundary tests | Issue | Existing Linux suite and independent retention integration pass; malformed decoder permutations are missing from existing tests (I-087); Windows safety skips disclosed below |
| T-053.1 useful commits/dirty scoped identity | Issue | Clean/dirty/unborn ordinary cases pass; colon paths and invalid globs produce false unchanged identities (I-089/I-090) |
| T-053.2 stable/sensitive/private fingerprints | Issue | Order, exclusion, no-content-retention and real-file hash tests pass; relevant colon-file content changes are invisible (I-090) |
| T-053.3 actual ancestry and honest distance | Proven | Existing same/ahead/behind/diverged/shallow/missing-object matrix; independent four-commit sequence and reverse relation produce exactly four; dirty same-commit comparison says changed |
| T-053.4 failure/path/platform handling | Issue | Ordinary failures, worktrees and symlinks pass on Linux; malformed glob, colon omission and native fixture failure recorded (I-089/I-090/I-091) |
| T-053.5 structured arguments/cancellation | Proven | Source uses exec.CommandContext and context deadlines; TestRunnerBoundaryUsesArgumentVectors, cancelled/timeout cases, secret-safe runner failure tests |
| T-053.6 independent temporary-repo evidence | Issue | Git and byte-hash oracles used on Linux; native Windows TestObserveUnusualFilenames fails (I-091) |
| T-054.1 thresholds/non-overridable blockers | Issue | Threshold limits/default/override matrix passes for consistent values; high/critical with zero total bypasses blocker (I-085) |
| T-054.2 one-level decline/improvement | Proven | TestAssessMaterialDeclineAndImprovement covers every direction, exact minimum, median last three, floor and unchanged poor values; independent stored 90→85 gives Watch |
| T-054.3 official compatible baselines | Proven | Existing all-series-component reset/manual/partial/failure tests plus independent save-load-assess sweep; manual snapshots never enter the official baseline |
| T-054.4 range/trend after three | Proven | TestAssessTrendWording/ExactWording and GoodIsCappedAtWatchUntilThreeComparableChecks: 0–6, ties, zero, finite filters, newest-five and median-three boundaries |
| T-054.5 honest missing/partial/stale summaries | Proven for explicit outcome/freshness states | Existing 5×9×3 assessment matrix and direct NaN/Inf/nil cases; decoder loss of missing-number state is separately I-086 |
| T-054.6 pure comprehensive boundary tests | Issue | Determinism/input immutability and expected classification tables pass; hard-blocker contradiction is absent from existing expected tables (I-085) |

Objective success conditions 1–3 and 4 are **Issue** through the links above. Condition 5 (incompatible history kept apart) and condition 6 (permanent official / explicit newest-ten manual retention) are **Proven**. No composite score exists. Current evidence does not authorize Objective closure.

## Completed frozen matrix

The full pre-probe lock and workflow inventory are preserved below. This table records every row's outcome, including mixed cells.

| Cell | Passed cells | Issue cells | Unverified / not applicable |
|---|---|---|---|
| M1 | fixed vocabularies, all outcome/freshness pairs, measured-only values, versions/origins, bad timestamps/IDs, path safety, units, nonfinite values, unknown/duplicate/wrong/mixed shape; reason lengths 199/200/201; mutation-after-validation | missing/null numbers I-086; trailing delimiters I-087; contradictory severity I-085; malformed pattern I-089 | None needed to establish these defects |
| M2 | snapshot deep-copy/order tests for results/details/summary/scope, observation sensitivity, all series/config component tests | same-path/line evidence-note permutations I-088 | Snapshot schema versions other than 1 are rejected; no cross-version join surface exists |
| M3 | absent roots, unsafe directories/files, unknown schema/name mismatch, preserved conflict bytes, temp leftover/retry, identical mtime, load/assess integration | malformed decoding I-087 | Mid-write/sync/close OS faults and secondary unlink failures not dynamically injected; cleanup/error propagation inspected, no claim of power-loss testing |
| M4 | manual 0–12 independent sweep, 14 and ties existing tests, mixed official origins, repeat, read-only plan, no implicit prune, invalid/unwritable refusal | None | Failure after a nonempty removal prefix not dynamically injected; loop/report preservation traced. No transactional rollback promised |
| M5 | normal/staged/deleted/untracked/excluded/unborn/nonrepo/subdirectory/worktree/shallow, symlink/dangling, whitespace/Unicode names on Linux, known runner failures | invalid glob I-089; colon path I-090; Windows filename setup I-091 | Native Windows symlinks unavailable to this account; Unix chmod tests N/A to Windows. Linux exercises both |
| M6 | same/ahead/behind/diverged/unknown/no commit/unavailable, shallow count refusal, cancellation, malformed rev-list and independent four-commit reversal | None | No network/remote fetch surface; direct fabricated Relation values are not ancestry proof |
| M7 | literal threshold transitions, hard minimum overrides for consistent evidence, median decline, nonfinite/empty/current confidence, Overall precedence, pure Summary conversion | zero-total severity contradiction I-085 | Invalid direct model guidance is expected to pass configuration/model validation before use; no extra product contract imposed |
| M8 | official/manual/incompatible/partial/failed histories, empty/short/3–6 observations, ties, zeros, finite filters, bounded UTF-8, stored history pipeline | None beyond M1/M7 integration findings | Display width/grapheme shaping N/A: no TUI renderer |
| M9 | fresh Linux full gate, all six cross-builds, Linux race test/vet, file reality and dependency/gate unchanged check | native Windows package gate fails I-091 | Native full-suite CI success not established; package failure is already concrete. Design follow-up recorded separately |

All invariants have been classified; not dynamically exercised OS fault timings are explicitly unverified as tests, not represented as executed proof. They do not turn this NEEDS WORK result into clearance. The repair recheck must retain this limitation or supply evidence; it must not silently expand those timings into unrelated storage requirements.

External boundary completion: actual Git executable/default and missing override exercised; ordered discovery/queries traced by the runner tests; success/nonzero/missing objects and malformed count responses exercised; cancellation/timeout cases return typed errors, underlying stderr is discarded; repeat observations are read-only. Output cap implementation inspected; exceeding 32 MiB not dynamically exercised. Server connection/redirect/provider retry cells are N/A to this local-only module. No external service was contacted.

Workflow completion: all storage and repository operations in the pre-probe inventory traced; setup/publication conflicts, retry, ignored interrupted temp, invalid history and permission failure exercised. Successful records reloaded semantically, not merely counted; official identities remained present across pruning. Pure classification was exercised after actual storage round trips. Remaining fault-injection limits are disclosed in M3/M4 above.

## Commands and platform evidence

Run on 2026-09-30 UTC (2026-10-01 Australia/Sydney), Go 1.26.2 linux/amd64:

- `make test-full`: **PASS**, fresh host-wide tests and linux/darwin/windows amd64/arm64 builds; no previous run reused. Codehealth package 1.322s in full gate.
- `go test -race -count=1 ./internal/codehealth`: **PASS**, 2.257s.
- `go vet ./internal/codehealth`: **PASS**.
- `git diff --check`: **PASS**.
- `go test -overlay /tmp/o027-overlay.json -count=1 -run TestO027Independent -v ./internal/codehealth`: **FAIL**, completed all subtests: six finding groups fail; independent stored-history 0–12 sweep, bounds/shape checks and four-commit ancestry sequence pass. Final missing-number probe was strengthened to preserve the original valid zero snapshot ID; its focused rerun also fails.
- `GOOS=windows GOARCH=amd64 go test -c -o /tmp/o027-codehealth.test.exe ./internal/codehealth`: **PASS** (cross-compiled test binary).
- Windows PowerShell native execution of that binary from the package fixture directory with quoted `"-test.v"`: **FAIL**, exit 1, `TestObserveUnusualFilenames` fails opening a filename containing a double quote. All other entered tests pass or skip. Windows temp projects were automatically cleaned by testing.T.
- Native Windows skips: symlink creation lacks account privilege; Unix file/directory permission-bit scenarios are not applicable as written. These are not claimed as passing Windows safety tests.
- The first PowerShell attempt parsed unquoted `-test.v` as `-test` and exited 2 before testing; corrected quoted invocation above is the evidence run.

The checker used temporary overlay tests; no existing test/source file was patched. The embedded harness below is the durable reproduction artifact. /tmp copies and compiled executable are disposable scratch. No new dependency was added; go.mod, go.sum and Makefile have no diff from 9fdc11f.

## Guardrails, file reality and architecture reconciliation

DATA-03 / CFG-01 fail on I-085/I-086/I-087/I-089. CFG-02 / CFG-03 fail on I-090/I-091. Linux TEST-08 full gate passes, but that does not negate the explicit native Windows failure. TEST-01/02 gaps are tied to the individual reproductions. FS-01/03/04/05/06 and TEST-03 preservation paths have no reproduced destructive-write finding. DATA-02 is satisfied: Code Health outcome/freshness vocabulary does not redefine Task lifecycle. ARCH-03 uses explicit project roots in reviewed call paths; ARCH-04 retains the Code Health module boundary; DEP-01 adds no dependency. TEST-04 temp-project isolation and TEST-09 recorded owner waivers are satisfied.

All files named as created in the four Tasks exist, including formerly missing Context Files and both fixtures. Historical notes saying no waiver was recorded are superseded by owner-authored frontmatter waivers. T-052's old unrelated .claude/settings.local.json note is historical, not a missing implementation artifact.

Design.md currently documents the older project/record architecture and contains no Code Health package or health-directory schema section. The implemented module is consistent with O-027's confirmed boundary: versioned JSON records under .savepoint/health, pure classification/history, local Git observation, no provider execution or UI integration. Task evidence says config.json; T-051's earlier reference to planned config.yml is superseded by T-052's explicit persistence choice. Record this architecture delta for planning reconciliation; this Check does not rewrite Design to remediate it. No sixth capability, generic provider API, automatic Issue generation, or composite score was introduced.

## Owner validation and observations

- T-053 and T-054 declare owner validation. Their owner board-completion waivers are recorded; no new acceptance is inferred here. The owner may still review repository wording and the compact classification scenarios during remediation. T-054's body simultaneously says the User Check table was reviewed and later says it was not presented; reconcile that historical wording in the evidence without manufacturing acceptance.
- Stored observations expose full incompatible/manual series via LoadSnapshots; the bounded assessment explanation may truncate trailing history notes. No UI visibility contract is tested here.
- Code Health's static validation is not a recomputation of all historical classifications. This review does not impose history-dependent validation on every stored summary.
- Existing Windows issue I-030 was searched and concerns different repaired modules. No prior O-027 Check or matching Issue existed; these are new findings, not a broadened recheck.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [ ] STYLE-03 **Test branches** — independent M1/M2/M5/M7 probes expose missed branches in config.go, identity.go, repository.go and classification.go.
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [ ] STYLE-06 **Handle errors at boundaries** — config.go:159 and repository.go:163 suppress malformed input information.
- [x] STYLE-07 **One source of truth**
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs** — four scoped module Tasks; no unrelated implementation changes.

## Handoff

Repair I-085 through I-091 directly under their Issue records, retaining all owner-completed Task statuses. A fresh independent Full Check must reuse this frozen scope and verify the reproductions plus the affected gates. This Check does not close O-027, repair implementation, or record owner acceptance.

## Reviewed source fingerprints

```text
51eb26f5ac253421aa6dc197cb855217715193e477881e98d7d332adc97a0f44  internal/codehealth/classification.go
70c38a5c3ee63fcb9eb1bcaa897f92bf4b1b410685b8f29a33342cd5d0d96d61  internal/codehealth/classification_test.go
2549f76363da391689ed5a07ef050b73de5bc4db716fc0ebed05c5ea19e67365  internal/codehealth/config.go
b9422b1b3dfbe1aec51c38c4a4c6d933dafc8f221009f9143de96ac89bb88880  internal/codehealth/errors.go
3ce4cbb64e8ea53390aa8768e6278ebd1e65ec28394c7e23172891c74855d431  internal/codehealth/history.go
8be80ed31d434538f1f27fa6c4ac7cca48ade0941203870e960a0d6be5080939  internal/codehealth/identity.go
02154f603637f8dbf4be24b5b1a2148bfcf8329e5e380168d1e4141b6e174396  internal/codehealth/identity_test.go
e53c24da5e762e232e8f6bb3af96dcc53dea2fe763e3976c8cd7b0bdb1f71c90  internal/codehealth/model.go
3d39e8f02f1ec6ec31b3796cc685ffd2a6843f93b8ad9db4b094a62d662e11c9  internal/codehealth/model_test.go
bb59f6d20d4f7e2055f9f8de247797e4a3e8c837e637956e155b19bd689c54cf  internal/codehealth/primitives.go
112b0530d76171682a33354ca8b1b63ea617aef50495e1ff89d6b3c828077373  internal/codehealth/repository.go
0064544a7591dd9e4e1b0cb68330e00c513cd3080b9950d5cb9b88bb5c4718c4  internal/codehealth/repository_test.go
56c1f18f238ec1fbb4dfc4d167cbdde506140b07e2dbc01268ebdd523737b718  internal/codehealth/snapshot.go
64633bcbf20f9ec6f6f47a8e9113da13e7ef1a5c0b195e2fa3c28e4ec0373b06  internal/codehealth/storage.go
63f5ff055b29a99280af18e9775445066739012ca61710a332d0123fafde6508  internal/codehealth/storage_test.go
83f5d16ff2641d34e5a3087fcae6e86773ee3acf0f7d1a761d02e625a2a498b1  internal/codehealth/testdata/valid-config-v1.json
cc0d7df9cacbd5c0a93630785b65a25f56e8205610a660425a35730e9ab3ea06  internal/codehealth/testdata/valid-snapshot-v1.json
```

## Pre-probe scope and workflow lock

# O-027 initial Full Check scope lock

1. T-051 criteria 1–6; T-052 1–6; T-053 1–6; T-054 1–6; O-027 six success conditions. Apply their named Guardrails plus FS-03, TEST-03, CFG-03 and mandatory fresh make test-full. Owner-confirmed thresholds and three-official-observation Good cap govern.
2. Scope: all internal/codehealth source, tests and two fixtures introduced by e974891, 78de829, 86ab71a and uncommitted classification/history files; AGENTS map and O-027 evidence. Public surfaces: vocabulary helpers; Config.Validate/DecodeConfig; Snapshot.Validate/DecodeSnapshot/Canonical/ComputeID; CapabilityConfig.Digest/CapabilityResult.SeriesID; Store constructor/config/save/load/prune; InputScope.Validate, GitRunner.Run, ObserveRepository, Observation.Unborn/Describe, Relate, Relation.Describe, DescribeRepositoryError; Assess, Assessment.Summary, Overall.
3. Dependencies: standard library JSON, hashes, filesystem, exec and local Git. No provider execution, CLI/TUI integration, network, Check snapshot wiring, or concurrent hostile filesystem replacement in this scope. History API explicitly requires earlier observations oldest-first; misuse of that order is outside scope. Existing input records and direct model calls are both in scope.
4. Frozen matrix rows (each includes named existing tests plus independent scenarios):

| Cell | Surfaces and inputs | State/representation/boundary | Oracle |
|---|---|---|---|
| M1 | config/snapshot validation and decode; five capabilities, nine outcomes, three freshness values | missing/empty/unknown/duplicate/wrong/mixed/null; trailing JSON delimiters; measured/value restrictions; detail severity count consistency; bounded fields at limit±1; revalidation after mutation | literal allowed states and JSON validity |
| M2 | canonical snapshot and series/config identity | permutations including evidence same path/line with different notes; nil/empty, mutation isolation, commit excluded from series; each compatibility component changed | byte-equivalence of reordered content; hash change/no-change |
| M3 | Store config/snapshot/load | absent/dirty dirs; version/name/content mismatch; conflicting/existing/symlink files; retry, temp leftover; config replacement; immutable save-load-assess integration | byte preservation, IDs and independently specified values |
| M4 | PlanPrune/Prune | manual 0–12, ties, mixed official series, repeat, malformed and unwritable stores; load/save never prune | chronological IDs; official bytes permanent; report equals removals |
| M5 | ObserveRepository/InputScope/GitRunner | clean/dirty/staged/deleted/untracked/excluded/unborn/nonrepo/subdirectory/worktree/shallow; unusual POSIX colon filename; malformed glob; enumeration order; symlink | real Git and file changes; diagnostics instead of silent omission |
| M6 | Relate/descriptions | same/ahead/behind/diverged/unknown/no commit/unavailable; failures/cancel/timeout; counts only if proven | independent Git commands and literal descriptions |
| M7 | Assess/Overall/Summary | all capabilities and states; threshold limit±epsilon; material decline limit±epsilon; one-level worsening, improvement; high/critical including contradictory total/severity, missing severity; direct invalid models/configuration | literal expected classifications; validation at storage boundary |
| M8 | history integration | 0–6 observations, official/manual/partial/failed/incompatible; ties/zeros/nonfinite; three-point cap, median last three and window five; UTF-8 bounded explanation | literal range/count/classification; storage round trip |
| M9 | gates/platform/architecture/evidence | full host tests and six cross-builds; native Windows evidence; file reality; Design and map reconciliation | command outcomes and scoped source |

External-boundary matrix: Git executable target/default/override, launch per query and ordered discovery, missing executable/nonrepo/nonzero/success, malformed output, cancellation/deadline, bounded output, secret-safe errors, repeat and no writes. Network connection/refusal/redirect/retry protocols N/A: local Git only. Filesystem failures are M3/M4. TTY/color/width/output sinks N/A: pure returned strings and no rendering. Unicode grapheme widths N/A; UTF-8 byte bounding applies M8. Lifecycle terminal revival/authorization/billing N/A: immutable observations, no lifecycle reducer.
5. Material findings require an actual supported-path violation of these criteria, introduced or promised here, with concrete consequence. No new product requirements. Windows execution unavailable locally must be reported distinctly from cross-builds.

## Workflow inventory

| Order | Operation | Effect | Failure/final state | Cleanup | Oracle |
|---|---|---|---|---|---|
| S1 | validate/encode then confined directory checks | may create health directories after validation | fatal before file publication | no record to remove | existing bytes unchanged |
| S2 | config read or snapshot identity collision read | none | fatal on unsafe/conflicting record; unchanged on identical | none | bytes/mtime |
| S3 | temp create/write/sync/close/chmod | unpublished temp | fatal; old record remains | remove temp (secondary failure may leave ignored temp) | old snapshot valid |
| S4 | config rename or snapshot exclusive hard link | record becomes visible | fatal if publish fails; identical race is unchanged success | unlink snapshot temp | atomic complete decodable record |
| P1 | load/validate/sort entire history and plan | none | fatal on any invalid history | none | no change on invalid load |
| P2 | delete planned manual IDs oldest-first | sequential explicit removals | fatal on first removal failure; report removed prefix and retained rest | none; no rollback promised | official bytes and remaining IDs |
| R1 | validate scope, Git discovery/head/shallow/list | read-only subprocesses | named failure; no observation | context cancel/process reap | ordered runner calls |
| R2 | hash relevant files, query changed tracked paths | read-only files/Git | named read failure or cancellation | close files | hashes and Git status |
| R3 | verify commits/ancestry/count | read-only Git | unknown/unavailable or cancellation; no invented count | process cleanup | rev-list independent counts |
| A1 | validate/load result, select series, classify, summarize | pure values | explicit unknown/caveat/blocker | none | literal classifications |

Failure injection will use existing permission/conflict tests and independent malformed-record/round-trip probes. Power-loss durability, malicious races, external provider processes and native Windows runtime are not silently claimed tested.


The initial note that native Windows was unavailable was corrected when Windows PowerShell and Go/Git were discovered through WSL. This changes evidence availability, not the frozen Windows/platform cell or acceptance perimeter. Native package execution then reproduced I-091.

## Reproduction harness

Save this Go block to /tmp/o027_check_test.go. The harness relies on existing package test fixture helpers and uses fresh temporary projects. Map a virtual internal/codehealth/o027_check_test.go to it with Go's overlay option; do not add it to production source as part of this Check.

```json
{"Replace":{"/home/user/code/savepoint/internal/codehealth/o027_check_test.go":"/tmp/o027_check_test.go"}}
```

```go
package codehealth

import (
 "context"
 "encoding/json"
 "fmt"
 "os"
 "path/filepath"
 "reflect"
 "runtime"
 "strings"
 "testing"
)

func TestO027Independent(t *testing.T) {
 t.Run("M1_missing_numbers",func(t *testing.T){
  for _,value:=range []string{`{"unit":"count"}`,`{"number":null,"unit":"count"}`} {
   r:=result(CapabilityTests,0)
   a:=Assess(OriginOfficial,r,nil,official(CapabilityTests,0,0))
   base:=validSnapshot(t);base.Results=[]CapabilityResult{r};base.Summary=Summary{Overall:a.Classification,Capabilities:[]CapabilitySummary{a.Summary()}};base=reseal(base)
   if err:=base.Validate();err!=nil{t.Fatal(err)}
   original,_:=json.Marshal(base)
   var obj map[string]any
   if err:=json.Unmarshal(original,&obj);err!=nil{t.Fatal(err)}
   var v any;if err:=json.Unmarshal([]byte(value),&v);err!=nil{t.Fatal(err)}
   obj["results"].([]any)[0].(map[string]any)["value"]=v
   raw,_:=json.Marshal(obj)
   // Keep the original ID: removing a measured zero must be detectable.
   got,err:=DecodeSnapshot(raw)
   if err==nil {t.Errorf("missing/null number accepted as measured %v, summary %s, original ID unchanged",got.Results[0].Value.Number,got.Summary.Overall)}
  }
 })
 t.Run("M1_trailing_delimiters", func(t *testing.T) {
  for _, suffix := range []string{"}", "]", " } garbage", " ] {}"} {
   t.Run(suffix, func(t *testing.T) {
    raw := append(readFixture(t,"valid-config-v1.json"), []byte(suffix)...)
    if json.Valid(raw) { t.Fatal("oracle expected invalid JSON") }
    if _,err:=DecodeConfig(raw); err==nil { t.Error("invalid JSON config accepted") }
    raw=append(readFixture(t,"valid-snapshot-v1.json"), []byte(suffix)...)
    if _,err:=DecodeSnapshot(raw); err==nil { t.Error("invalid JSON snapshot accepted") }
   })
  }
 })
 t.Run("M2_evidence_tie_permutation",func(t *testing.T){
  a:=validSnapshot(t)
  a.Results[0].Evidence=[]EvidenceRef{{Path:"report.json",Line:1,Note:"first"},{Path:"report.json",Line:1,Note:"second"}}
  a=reseal(a)
  if err:=a.Validate();err!=nil {t.Fatal(err)}
  b:=a.Canonical()
  for i:=range b.Results {if b.Results[i].Capability==CapabilityCoverage {e:=b.Results[i].Evidence;e[0],e[1]=e[1],e[0]}}
  if a.ComputeID()!=b.ComputeID(){t.Errorf("same evidence permutation changes identity: %s versus %s",a.ComputeID(),b.ComputeID())}
 })
 t.Run("M7_severity_validation_integration",func(t *testing.T){
  for _,tc:=range []struct{name string;total,high,critical float64}{
   {"zero_total_high",0,1,0},{"zero_total_critical",0,0,1},{"negative_high",2,-1,0},{"fractional_high",2,0.5,0},
  } {t.Run(tc.name,func(t *testing.T){
   r:=vulnResult(tc.total,Detail{Key:"high",Number:tc.high},Detail{Key:"critical",Number:tc.critical})
   a:=Assess(OriginOfficial,r,nil,official(CapabilityDependencyVulnerability,tc.total,tc.total))
   s:=validSnapshot(t);s.Results=[]CapabilityResult{r};s.Summary=Summary{Overall:a.Classification,Capabilities:[]CapabilitySummary{a.Summary()}};s=reseal(s)
   err:=s.Validate()
   t.Logf("validation=%v classification=%s explanation=%s",err,a.Classification,a.Explanation)
   if err==nil {t.Error("impossible severity counts accepted by snapshot contract")}
   if tc.high>0 || tc.critical>0 {if a.Classification!=ClassificationNeedsAttention{t.Error("positive severity count did not block")}}
  })}
 })
 t.Run("M5_malformed_pattern",func(t *testing.T){
  root:=newRepo(t)
  for _,p:=range []string{"[", "src/[abc", "src/[z-a]"} {
   scope:=InputScope{Include:[]string{p}}
   // Descending ranges are syntactically accepted by path.Match; only
   // unterminated brackets are mandatory rejection cases here.
   if strings.Contains(p,"]"){continue}
   if err:=scope.Validate();err==nil{t.Errorf("malformed glob %q accepted",p)}
   cfg:=validConfig(t);cfg.Capabilities[0].Scope=[]string{p}
   if err:=cfg.Validate();err==nil{t.Errorf("malformed config glob %q accepted",p)}
   before,err:=ObserveRepository(context.Background(),GitRunner{},root,scope)
   write(t,root,"a.go","changed")
   after,e2:=ObserveRepository(context.Background(),GitRunner{},root,scope)
   t.Logf("pattern=%q errors=%v/%v clean=%v unchanged=%v",p,err,e2,!after.Identity.Dirty,before.Identity==after.Identity)
  }
 })
 t.Run("M5_colon_filename",func(t *testing.T){
  if runtime.GOOS=="windows"{t.Skip("colon filename cannot exist on Windows")}
  root:=newRepo(t)
  writeFile(t,filepath.Join(root,"source:part.go"),"before")
  git(t,root,"add",".");git(t,root,"commit","-m","colon")
  before,err:=ObserveRepository(context.Background(),GitRunner{},root,InputScope{});if err!=nil{t.Fatal(err)}
  writeFile(t,filepath.Join(root,"source:part.go"),"after")
  after,err:=ObserveRepository(context.Background(),GitRunner{},root,InputScope{});if err!=nil{t.Fatal(err)}
  t.Logf("before=%+v after=%+v",before,after)
  if !after.Identity.Dirty || before.Identity.InputFingerprint==after.Identity.InputFingerprint {t.Error("tracked content change silently ignored")}
 })
 t.Run("M3_M4_M8_stored_history",func(t *testing.T){
  for n:=0;n<=12;n++ {t.Run(fmt.Sprint(n),func(t *testing.T){
   st,root:=newProject(t)
   var wantedOfficial []string
   for i:=0;i<3;i++ {
    s:=snapAt(t,OriginOfficial,i+1);r:=result(CapabilityCoverage,90)
    r.CollectedAt=s.CreatedAt;s.Results=[]CapabilityResult{r};s.Summary=Summary{Overall:ClassificationWatch,Capabilities:[]CapabilitySummary{{Capability:r.Capability,Provider:r.Provenance.Provider,Classification:ClassificationWatch}}};s=reseal(s)
    mustSave(t,st,s);wantedOfficial=append(wantedOfficial,s.ID)
   }
   for i:=0;i<n;i++ {mustSave(t,st,snapAt(t,OriginManual,i+20))}
   loaded,err:=st.LoadSnapshots();if err!=nil||len(loaded)!=n+3{t.Fatalf("load: %d %v",len(loaded),err)}
   var h []HistoryEntry
   for _,s:=range loaded{for _,r:=range s.Results{h=append(h,HistoryEntry{Origin:s.Origin,Result:r})}}
   a:=Assess(OriginManual,result(CapabilityCoverage,85),nil,h)
   if a.Classification!=ClassificationWatch || a.Trend.Min!=90 || a.Trend.Max!=90 || a.Trend.Observations!=3 {t.Errorf("wrong stored baseline %+v",a)}
   p,err:=st.PlanPrune();if err!=nil{t.Fatal(err)}
   unchanged,_:=st.LoadSnapshots();if len(unchanged)!=n+3{t.Fatal("plan mutated history")}
   got,err:=st.Prune();if err!=nil||!reflect.DeepEqual(got,p){t.Fatalf("prune mismatch %v",err)}
   if len(got.Removed)!=max(0,n-10){t.Fatal("wrong retention")}
   for _,id:=range wantedOfficial{if _,err:=os.Stat(filepath.Join(root,".savepoint","health","snapshots",strings.TrimPrefix(id,"sha256:")+".json"));err!=nil{t.Fatal("lost official",err)}}
   again,err:=st.Prune();if err!=nil||len(again.Removed)!=0{t.Fatal("repeat prune")}
  })}
 })
 t.Run("M1_bounds_and_shapes",func(t *testing.T){
  for _,n:=range []int{199,200,201}{s:=validSnapshot(t);s.Results[1].Reason=strings.Repeat("x",n);s=reseal(s);err:=s.Validate();if (err==nil)!=(n<=200){t.Errorf("reason length %d: %v",n,err)}}
  for _,raw:=range []string{`{"version":1,"capabilities":[null]}`,`{"version":1,"capabilities":[1]}`,`{"version":1,"capabilities":["tests"]}`} {if _,err:=DecodeConfig([]byte(raw));err==nil{t.Errorf("bad shape accepted: %s",raw)}}
 })
 t.Run("M6_independent_ancestry_sequence",func(t *testing.T){
  root:=newRepo(t);first:=observe(t,root,InputScope{})
  for i:=0;i<4;i++ {commitFile(t,root,"step.go",fmt.Sprint(i))}
  last:=observe(t,root,InputScope{})
  r,err:=Relate(context.Background(),GitRunner{},root,first.Identity,last)
  if err!=nil||r.Kind!=RelationBehind||r.CurrentOnly!=4||!r.Counted{t.Fatalf("expected four behind: %+v %v",r,err)}
  reverse,err:=Relate(context.Background(),GitRunner{},root,last.Identity,first)
  if err!=nil||reverse.Kind!=RelationAhead||reverse.RecordedOnly!=4{t.Fatalf("expected four ahead: %+v %v",reverse,err)}
  write(t,root,"step.go","uncommitted")
  dirty:=observe(t,root,InputScope{});same,_:=Relate(context.Background(),GitRunner{},root,last.Identity,dirty)
  if same.InputsMatch||!strings.Contains(same.Describe(),"changed"){t.Fatal("same commit lost changed input wording")}
 })
}

```
