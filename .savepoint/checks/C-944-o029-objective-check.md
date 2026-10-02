---
id: C-944
scope: {kind: objective, id: O-029}
result: NEEDS WORK
checked_by: {role: checker, session: o029-full-check-20261001}
executed_session: o029-executor-sessions-20261001
checked_at: '2026-10-01T09:35:25Z'
reviewed:
  base_commit: 32ae341
  head_commit: e55eb7eb88f7dbc0697ecb5ec7166614b9bedf17
  files:
    - internal/codehealth/collect.go
    - internal/codehealth/reader_paths.go
    - internal/codehealth/reader_paths_test.go
    - internal/codehealth/reader_tests.go
    - internal/codehealth/reader_junit.go
    - internal/codehealth/reader_tests_test.go
    - internal/codehealth/reader_coverage.go
    - internal/codehealth/reader_coverage_test.go
    - internal/codehealth/reader_lizard.go
    - internal/codehealth/reader_lizard_test.go
    - internal/codehealth/reader_jscpd.go
    - internal/codehealth/reader_jscpd_test.go
    - internal/codehealth/reader_osv.go
    - internal/codehealth/reader_osv_test.go
    - internal/codehealth/readers.go
    - internal/codehealth/readers_integration_test.go
    - internal/codehealth/classification.go
    - internal/codehealth/classification_test.go
    - internal/codehealth/snapshot.go
    - internal/codehealth/model_test.go
    - AGENTS.md
    - .savepoint/Design.md
    - .gitignore
    - internal/codehealth/testdata/readers/complexity/all-outside.csv
    - internal/codehealth/testdata/readers/complexity/bad-ccn.csv
    - internal/codehealth/testdata/readers/complexity/bad-line.csv
    - internal/codehealth/testdata/readers/complexity/empty.csv
    - internal/codehealth/testdata/readers/complexity/mixed.csv
    - internal/codehealth/testdata/readers/complexity/outside.csv
    - internal/codehealth/testdata/readers/complexity/wrong-columns.csv
    - internal/codehealth/testdata/readers/coverage/coveragepy-branches.json
    - internal/codehealth/testdata/readers/coverage/coveragepy-full.json
    - internal/codehealth/testdata/readers/coverage/coveragepy-malformed.json
    - internal/codehealth/testdata/readers/coverage/coveragepy-no-totals.json
    - internal/codehealth/testdata/readers/coverage/coveragepy-populated.json
    - internal/codehealth/testdata/readers/coverage/coveragepy-zero.json
    - internal/codehealth/testdata/readers/coverage/go-full.out
    - internal/codehealth/testdata/readers/coverage/go-malformed.out
    - internal/codehealth/testdata/readers/coverage/go-merged.out
    - internal/codehealth/testdata/readers/coverage/go-monorepo.out
    - internal/codehealth/testdata/readers/coverage/go-populated.out
    - internal/codehealth/testdata/readers/coverage/go-unresolvable.out
    - internal/codehealth/testdata/readers/coverage/go-zero.out
    - internal/codehealth/testdata/readers/coverage/vitest-full.json
    - internal/codehealth/testdata/readers/coverage/vitest-malformed.json
    - internal/codehealth/testdata/readers/coverage/vitest-outside.json
    - internal/codehealth/testdata/readers/coverage/vitest-populated.json
    - internal/codehealth/testdata/readers/coverage/vitest-zero.json
    - internal/codehealth/testdata/readers/duplication/absolute.json
    - internal/codehealth/testdata/readers/duplication/inconsistent.json
    - internal/codehealth/testdata/readers/duplication/malformed.json
    - internal/codehealth/testdata/readers/duplication/mixed.json
    - internal/codehealth/testdata/readers/duplication/no-statistics.json
    - internal/codehealth/testdata/readers/duplication/none.json
    - internal/codehealth/testdata/readers/duplication/versioned.json
    - internal/codehealth/testdata/readers/duplication/zero-lines.json
    - internal/codehealth/testdata/readers/tests/go-build-failure.jsonl
    - internal/codehealth/testdata/readers/tests/go-empty.jsonl
    - internal/codehealth/testdata/readers/tests/go-fail.jsonl
    - internal/codehealth/testdata/readers/tests/go-malformed.jsonl
    - internal/codehealth/testdata/readers/tests/go-monorepo.jsonl
    - internal/codehealth/testdata/readers/tests/go-pass.jsonl
    - internal/codehealth/testdata/readers/tests/go-truncated.jsonl
    - internal/codehealth/testdata/readers/tests/junit-empty.xml
    - internal/codehealth/testdata/readers/tests/junit-malformed.xml
    - internal/codehealth/testdata/readers/tests/junit-pass.xml
    - internal/codehealth/testdata/readers/tests/pytest-fail.xml
    - internal/codehealth/testdata/readers/tests/vitest-fail.xml
    - internal/codehealth/testdata/readers/vulnerabilities/absolute.json
    - internal/codehealth/testdata/readers/vulnerabilities/all-severities.json
    - internal/codehealth/testdata/readers/vulnerabilities/malformed.json
    - internal/codehealth/testdata/readers/vulnerabilities/missing-severity.json
    - internal/codehealth/testdata/readers/vulnerabilities/multi-ecosystem.json
    - internal/codehealth/testdata/readers/vulnerabilities/no-results.json
    - internal/codehealth/testdata/readers/vulnerabilities/none.json
    - internal/codehealth/testdata/readers/vulnerabilities/scanner-error.json
    - internal/codehealth/testdata/readers/vulnerabilities/ungrouped.json
    - internal/codehealth/testdata/readers/vulnerabilities/unresolved.json
  dependencies: []
issues: [I-097, I-098, I-099, I-100, I-101]
supersedes: null
---

# C-944: O-029 Full Objective Check

## Result and authority

**NEEDS WORK.** Independent fresh checker session; this conversation did not implement any scoped work. The review includes committed T-059..T-063 plus the existing uncommitted T-064/T-065 code, fixtures and documentation. Head commit identifies the working-tree base, not a claim that every reviewed file is committed. No implementation, Design, router, Task or Objective status was changed by the checker.

All seven Tasks are done with owner-attributed Task-check waivers. No owner validation is declared on these Tasks or Objective. This Check covers the waivers rather than inventing local Task clearance. No prior O-029 Check or directly linked Issue existed; matching reader symptoms/locations were searched before allocating I-097..I-101. I-086 is a repaired snapshot-decoding defect, not the new provider-reader defects here.

## Frozen scope and workflow lock

# Initial O-029 Full Check scope lock

1. Criteria: O-029 seven success conditions, T-059..T-065 every Done When; FS-05, ARCH-03/04, CFG-02/03, DEP-01/02, DATA-03, TPL-02, TEST-01/02/04/06/07/08/09. Required gate make test-full and git diff --check. Owner-completed Tasks remain done.
2. Changed files: readers/helpers and their tests and fixtures, Collect Root seam, severity model/classification changes, registry, Design reader bullet and AGENTS Codebase Map. Public surfaces: RelPath, ReportInput.Path, LoadGoModules/Resolve, WorstEvidence, all nine reader paths, DefaultReaders, Collect, snapshot validation/storage, Assess. Collection command, Check/TUI wiring, installs, network analysis, unrelated lifecycle are outside scope.
3. Runtime: instance configuration/scope -> report acquisition or fake tool -> reader -> bounded normalized validation -> assessment -> immutable store -> reload. Existing process/storage behavior is relied on and covered by full gate; this Objective changes normalization, not those implementations.
4. Frozen rows: M0 path mapping/bounds/module lookup; M1 Go/Vitest/pytest tests; M2 Go/Vitest/coverage.py coverage; M3 Lizard complexity; M4 jscpd duplication; M5 OSV and severity; M6 registration/Collect containment/store; M7 docs, file reality, waivers and gates. For M1..M5 axes: valid, empty, malformed/wrong schema/type/null/missing, partial/unsupported variant, duplicate, mixed language, scope and exclusions, provenance, affected-item bounds/order. M1 adjacent cells: missing package terminal with unfinished OR finished tests, start-only, build-fail event, JUnit suite/case count disagreement and non-JUnit/multiple roots. M2 adjacent: duplicate blocks, modes/coordinates, negative/null counters, missing totals and scoped whole-project totals. M3 boundaries 10/11/20/21, CSV quotation, Unicode names, outside-root/no functions. M4 adjacent: all clones outside scope and mixed-scope first/second locations. M5 adjacent: duplicate identical group, null source, exact severity floors, nonfinite strings, unresolved/errors/ungrouped, empty findings, byte mutation and roundtrip. M0 boundaries MaxEvidence-1/MaxEvidence/MaxEvidence+1, MaxNoteLen; lexical Windows and outside-root paths. M6 sequence valid + malformed + timeout + unavailable + absent, repeated read/store and wrong-schema Go through Collect. Values use literal counts/ratios, not implementation calculations.
5. Supported-path materiality: approved provider reports through public Readers/Collect with accepted instance config; inaccurate health, masked missing/partial measurement, missing affected-item evidence or broken gates qualify. Unsupported arbitrary network protocols and provider analysis internals do not. Renderer/interactivity/color/dumb terminal and Unicode width axes N/A: no renderer changes. Lifecycle transitions N/A beyond per-instance outcomes and repeat collection: no lifecycle changes. Network redirect/retry/server startup N/A: no network client here. External subprocess cells (configured/actual target, ordering, missing executable, success/findings/non-success exit, timeout/cancel, malformed output, cleanup, secret-safe errors) covered by Collect/runner tests; fake runner avoids real analysis as Objective requires.

Workflow lock:
| Order | Operation | Effect | Failure owner/final state | Cleanup | Oracle |
|---|---|---|---|---|---|
|1|Validate config / root / history|none|collection fatal, no tools/snapshot|none|Collect tests|
|2|Observe repository and select instance|read-only|known unavailable identity or fatal preflight|none|repository tests|
|3|Read report or prepare temp report/run tool|temp file/process|instance failed/absent/unavailable/timed_out/cancelled; other instances continue|temporary report cleanup|fake runner trace, runner/Collect tests|
|4|Reader normalization|read-only root module lookup|instance failed or partial with reason|none|independent matrix literal oracle|
|5|Bound validation / classification|normalized result|failed instance no value; unknown classification|none|matrix plus validation tests|
|6|Encode/save immutable snapshot|new snapshot visible only on successful write|collection fatal on storage failure|store handles temporary artifacts|storage tests and reload semantic comparisons|
|7|Reload / repeated collection|read-only / another immutable result|named load error, no repair|none|roundtrip tests|


## Completed coverage matrix

P = proven by named fixture tests plus independent probes; I = captured Issue; U = unverified. N/A is explained in the scope lock. Input byte mutation and serialization do not represent mutable service state: readers are stateless and snapshots own normalized values. No per-reader terminal renderer exists.

| Row/public paths | Normal/empty | Malformed/missing/null/wrong type | Partial/unsupported variant | Scope/exclusion/mixed project | Boundaries/evidence/provenance | Duplicate/sequence/representation | Verdict |
|---|---|---|---|---|---|---|---|
| M0 RelPath, Path, module discovery/Resolve, WorstEvidence | P: relative/absolute/nested-module/unresolved | P: traversal/control/sensitive/empty rejected | P: Windows drive and slash inputs; outside-root rejected | P: include/exclude and nested modules | P: cap-1/cap/cap+1, bounded notes and ordering | P: deterministic tie order; pure lexical mapping | Proven T-059 |
| M1 Go test JSON | P: passing/failing/build-error/skip/empty fixtures | I-098: null/unrelated object gives complete zero | I-098: finished tests without package terminal and start-only stream; existing unfinished-test case P | I-097: excluded package failure still counted; monorepo path resolution P | P: safe bounded names/unknown version on mapped packages | Go subtest counting documented; standalone build-fail probe observation below | Issue |
| M1 Vitest/pytest JUnit | P: cases/failures/errors/skips/empty | I-098: multiple roots accepted; wrong-root wrapper observation | I-098: suite declares omitted failure but reading is complete | I-097: excluded test counted; pytest path conversion P | P: file/classname, names bounded, hostname/system-out ignored, version unknown | P: provider-shared parser; declared-count disagreement tested independently | Issue |
| M2 Go cover | P: statement totals, zero-total error | I-099: malformed block coordinates accepted | I-099: unsupported mode accepted; unresolved import P partial | P: scoped and monorepo aggregation | P: covered/total and least-covered order | P: duplicate-block merge and any-run coverage literal oracle | Issue |
| M2 Vitest V8 | P: populated/full/zero-total | I-099: null/negative counters become measurements; wrong type JSON rejected by decoder | P: unmappable file partial | P: scope and mixed collection | P: statement headline and function/branch details | P: input map ordering sorted; shared literal ratios | Issue |
| M2 coverage.py | P: totals/full/zero-total | P: missing/null required totals and inconsistent totals rejected | Partial/unsupported in-scope source evidence U where required by Objective fixture promise | I-097: filtered file evidence but global headline | P: meta.version, statements/branch counts | P: persisted through mixed collection | Issue; fixture evidence U |
| M3 Lizard | P: max/average/threshold counts; empty says no functions | P: column/CCN/start-line errors | I-101 U: unparsed in-scope files; outside-root rows P partial | P: mixed Go/TS/Python and filtered scope | P: 10/11/20/21, quoted comma/Unicode, cap, safe notes/version unknown | P: stateless repeated reads | Unverified unparsed criterion |
| M4 jscpd | P: percentage/count cross-check; zero-lines error | P: no totals/malformed/inconsistent report rejected | P: outside-root clone partial | I-097: global totals/scoped items mismatch; excluded firstFile evidence | P: largest clones/cap/version unknown or stated; ignored analyzer fields | P: fixture collection/storage; tolerance observation below | Issue |
| M5 OSV/classification/model | P: empty/mixed ecosystems/all buckets | P: malformed/no results; null source explicitly partial (observation) | P: unresolved/error/ungrouped/outside-root partial with reason | P: source scope/exclusion helper | P: exact CVSS floors, NaN/infinities, cap and bounded names, five-bucket sum, unknown blocks | I-100: duplicate identical group inflates count; P byte immutability and Reading JSON roundtrip | Issue |
| M6 DefaultReaders/Collect/store/Assess | P: nine-provider completeness + mixed run | P: failed malformed sibling no value/unknown; I-098 null Go bypass persists measured zero | P: timeout/unavailable/absent isolated; existing cancel/truncation tests | P: multiple named instances; baseline-vs-malformed sibling independent exact-value comparison | P: actual root passed, normalized result validation | P: immutable store reload Validate/ID, repeat collection; related failure-prefix/cleanup tests | Integration machinery Proven, reader defects propagate |
| M7 docs/files/waivers/gates | P: every evidence file exists; Design/Map describe readers and no command | P: all Task waivers name Task/actor/time | P: full gate/build outputs | P: changes confined to package responsibility; dependencies unchanged | Style checklist below | P: strict index loading after artifacts; native Windows CI for this working tree not obtained here | Required local gate Proven; native Windows execution unverified |

The fixture requirement in Objective success condition 6 is not claimed complete: real partial/unsupported forms for every provider have not been proven; I-098/I-099 capture missing wrong-schema/unsupported tests and I-101 captures the material unparsed-file gap. A recheck must establish the original per-provider fixture promise without expanding the frozen input axes.

## Per-criterion reconciliation

| Scope | Acceptance classification |
|---|---|
| O-029 success 1 tests | Issue I-098; basic failure/build/empty paths proven |
| O-029 success 2 coverage | Issue I-099; statement headline/branch/function details proven on valid fixtures |
| O-029 success 3 complexity | Unverified I-101 unparsed files; maximum and details proven |
| O-029 success 4 duplication | Valid normalization proven; scoped behavior Issue I-097 |
| O-029 success 5 vulnerabilities | Issue I-100 distinct totals; severity/classification/freshness wording/partial fixtures proven |
| O-029 success 6 reader paths/bounds/scope/fixtures | Bounds and repository-relative paths proven; scope Issue I-097; partial/unsupported fixture completeness Unverified |
| O-029 success 7 registry/containment | Proven; invalid reader-produced measured data still propagates (I-098/I-099) |
| T-059 criteria 1..5 | Proven; Root seam, lexical paths, nested module mapping, deterministic bounds and platform-shaped path cases |
| T-060 criteria 1..5 | 1 Issue I-098 terminal completeness; 2 valid JUnit results proven but unusable representation Issue I-098; 3 Issue I-098 null/wrong schema; 4 unknown provenance proven; 5 named existing fixture files proven, missing variants recorded |
| T-061 criteria 1..5 | 1 Issue I-099 malformed coordinates/mode; 2 Issue I-099 null/negative counts; 3 valid totals/meta proven, scope Issue I-097; 4 valid bounds/zero-total errors proven; 5 existing specified fixtures proven |
| T-062 criteria 1..6 | All Task-specific numeric/format/evidence/empty/version/fixture criteria proven; wider Objective unparsed-files criterion Unverified I-101 |
| T-063 criteria 1..5 | Valid totals/evidence/zero-lines/ignored fields/version/fixtures proven; wider scoped outcome Issue I-097 |
| T-064 criteria 1..5 | 1 Issue I-100 distinct groups; 2 severity validation/blocking proven; 3 bounded affected packages proven; 4 explicit empty/partial/version/database wording proven for supplied fixtures; 5 specified files proven |
| T-065 criteria 1..4 | Proven: registry, fault isolation, snapshot roundtrip, documentation. The fixture integration test alone does not prove semantic input validation; independent matrix supplies that evidence. |

## Adversarial findings and materiality

| Issue | Likelihood | Impact | Materiality | Recommendation |
|---|---|---|---|---|
| I-097 scoped counts | High: monorepo report reused by scoped instances | High: unrelated failures or coverage misclassify health | High | Fix scope semantics across affected readers together |
| I-098 unusable/incomplete tests | Medium: interrupted/corrupt/wrong report | High: recorded complete failure count hides missing evidence | High | Fix event/document validation and completeness |
| I-099 unusable coverage | Low: malformed counters or unsupported variant | Medium: invented good/bad measurements | Medium | Fix boundary validation narrowly |
| I-100 duplicate OSV groups | Low: duplicate input groups | Medium: inflated counts/severity/trends | Medium | Deduplicate or explicitly refuse duplicates |
| I-101 unparsed complexity evidence | Medium: parsing gaps possible; handling not established | Medium: unknown complexity may be shown as zero | Medium | Establish provider evidence or return requirement decision to planner/owner |

## Evidence and gates

- Fresh `make test-full`: exit 0, 2026-10-01 around 09:26–09:27 UTC, go1.26.2 linux/amd64. Full Go suite and linux/darwin/windows build steps succeeded. This is fresh evidence on the exact reviewed working tree; no implementation edits followed it. Native Windows test execution/CI on this uncommitted tree was not obtained, so cross-compilation is not described as native Windows evidence.
- `git diff --check`: exit 0. `go vet ./internal/codehealth`: exit 0.
- Focused named reader/helper/Collect/runner/store/classification tests: exit 0 (2.907s). Exact command: `go test ./internal/codehealth -run 'Test(RelPath|ReportInputPath|GoModules|WorstEvidence|CollectHandsReaders|GoTestReader|JUnitReader|CoverageReaders|GoCoverage|LizardReader|LizardEvidence|JscpdReader|JscpdEvidence|OSVScanner|DefaultReaders|CollectWithRealReaders|VulnerabilitySeverity|AssessCurrentValue|Collect|ExecRunner|Store)' -count=1`.
- Independent matrix: `go test -overlay=/tmp/o029-overlay.json ./internal/codehealth -run TestO029IndependentMatrix -v -count=1`. Harness exit 0 means all observational cells completed; the explicit ISSUE log lines are failed invariants, not a technical pass. Full output and source embedded below.
- File reality: 30 distinct Task-evidence path references inspected; none missing. All new fixture/source files named in reviewed.files exist. go.mod/go.sum and gate definitions have no O-029 changes. The source fixtures are handwritten; no external analysis or network-dependent test ran.
- Provider shape checks use primary sources: [OSV results schema](https://raw.githubusercontent.com/google/osv-scanner/main/pkg/models/results.go), [Lizard producer](https://raw.githubusercontent.com/terryyin/lizard/master/lizard.py), [Go build JSON documentation](https://pkg.go.dev/cmd/go#hdr-Build_json_encoding). These corroborate boundaries only and do not replace independent normalization probes.

Guardrails: FS-05 lexical slash normalization is paired with filepath for actual filesystem access; ARCH-03 uses explicit root; ARCH-04 reader responsibility stays in codehealth; DEP-01/02 no new dependency; CFG-02 platform-shaped path tests and platform builds pass, CFG-03 native execution is not certified here; TEST-01/02 have the gaps recorded above; TEST-04 tests use temporary projects/fakes; TEST-06/07 named outcome tests, not percentage claims; TEST-08 local full gate passed; TEST-09 owner waivers present. TPL-02/Design's broad malformed-report promise is limited by I-098/I-099, so documentation is not treated as stronger evidence than runtime.

## Observations (non-blocking)

- A standalone modern `build-fail` event with ImportPath is ignored by GoTestReader. The primary Go documentation confirms this event form, but a complete normal Go test build-failure stream also carries package fail events; the standalone probe does not establish undercounting of a complete run. Retain this variant in the frozen matrix as a partial/unsupported representation, without requiring double counting.
- JUnit accepts an html wrapper containing testsuite. This is outside normal Vitest/pytest document roots; multiple roots and contradictory declared counts are the admitted I-098 cases.
- OSV null result entry yields partial zero with an explicit missing-path reason; Assess treats partial data as unknown. The harness initially labeled this an ISSUE against strict rejection, but containment meets honest unavailability; it is not admitted as a separate blocker.
- jscpd's cross-check tolerance is 0.5 percentage points although the code comment says rounding to two decimals. No criterion fixes a numerical tolerance; this discrepancy is advisory.
- Lizard's affected-item list suppresses functions at or below the Good threshold, a documented executor interpretation. This review does not invent a requirement for healthy functions to be listed as affected.
- The prior O-028 observation on hand-configured executed reports outside Savepoint's report directory is unchanged pre-existing orchestration, outside this Objective's introduced normalization boundary; no new Issue is admitted here.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [ ] STYLE-03 **Test branches** — reader boundary/scope/completeness branches miss the cases in I-097..I-101.
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [ ] STYLE-06 **Handle errors at boundaries** — valid JSON/XML syntax is accepted without adequate report semantics; reader_tests.go:144, reader_junit.go:50, reader_coverage.go:149/218.
- [x] STYLE-07 **One source of truth**
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs** — provider-sized files/Tasks form the planned feature increment.

## Handoff

O-029 is not ready to close. Repair I-097..I-100 directly under the Issues by default and establish I-101's missing evidence/requirement decision; retain all done Task statuses. A fresh recheck uses this exact frozen scope/matrix, reproductions and original gate. No owner acceptance or exception is recorded.

## Independent harness

Save this source to a scratch file, create a Go overlay mapping a virtual `internal/codehealth/o029_check_test.go` to it, and run the command above. It relies only on existing temporary-project/fake-runner test helpers. The scratch files are not production edits.

```go
package codehealth

import (
 "context"
 "encoding/json"
 "fmt"
 "testing"
)

// Independent initial Check matrix. Observational assertions preserve all
// actual outcomes so one failed invariant never stops the remaining cells.
func TestO029IndependentMatrix(t *testing.T) {
 root := goRoot(t)
 type cell struct { name string; reader Reader; provider ProviderKey; data string; scope, exclusions []string; want float64; reject, partial bool }
 cells := []cell{
  {name:"M1 Go scoped count", reader:GoTestReader{}, provider:ProviderGoTestJSON, scope:[]string{"a/**"}, data:"{\"Action\":\"pass\",\"Package\":\"example.com/m/a\",\"Test\":\"TestOK\"}\n{\"Action\":\"pass\",\"Package\":\"example.com/m/a\"}\n{\"Action\":\"fail\",\"Package\":\"example.com/m/b\",\"Test\":\"TestBad\"}\n{\"Action\":\"fail\",\"Package\":\"example.com/m/b\"}"},
  {name:"M1 Go null event",reader:GoTestReader{},provider:ProviderGoTestJSON,data:"null",reject:true},
  {name:"M1 Go wrong schema",reader:GoTestReader{},provider:ProviderGoTestJSON,data:"{\"results\":[]}",reject:true},
  {name:"M1 Go missing final package",reader:GoTestReader{},provider:ProviderGoTestJSON,data:"{\"Action\":\"run\",\"Package\":\"example.com/m/a\",\"Test\":\"TestA\"}\n{\"Action\":\"pass\",\"Package\":\"example.com/m/a\",\"Test\":\"TestA\"}",partial:true},
  {name:"M1 Go incomplete start only",reader:GoTestReader{},provider:ProviderGoTestJSON,data:"{\"Action\":\"start\",\"Package\":\"example.com/m/a\"}",partial:true},
  {name:"M1 Go build-fail event",reader:GoTestReader{},provider:ProviderGoTestJSON,data:"{\"Action\":\"build-fail\",\"ImportPath\":\"example.com/m/a\"}",want:1},
  {name:"M1 Vitest scoped count",reader:JUnitReader{},provider:ProviderVitestJUnit,scope:[]string{"a/**"},data:`<testsuite><testcase classname="a/ok.test.ts" name="ok"/><testcase classname="b/bad.test.ts" name="bad"><failure/></testcase></testsuite>`},
  {name:"M1 Pytest excluded failure",reader:JUnitReader{},provider:ProviderPytestJUnit,exclusions:[]string{"b/**"},data:`<testsuite><testcase file="a/ok.py" name="ok"/><testcase file="b/bad.py" name="bad"><error/></testcase></testsuite>`},
  {name:"M1 JUnit mismatched declared total",reader:JUnitReader{},provider:ProviderVitestJUnit,data:`<testsuite tests="2" failures="1"><testcase name="ok"/></testsuite>`,partial:true},
  {name:"M1 JUnit wrong root bypass",reader:JUnitReader{},provider:ProviderVitestJUnit,data:`<html><testsuite/></html>`,reject:true},
  {name:"M1 JUnit two roots",reader:JUnitReader{},provider:ProviderVitestJUnit,data:`<testsuite/><testsuite/>`,reject:true},
  {name:"M2 Go merged duplicate",reader:GoCoverReader{},provider:ProviderGoCoverProfile,data:"mode: count\nexample.com/m/a/a.go:1.1,2.1 2 0\nexample.com/m/a/a.go:1.1,2.1 2 3\n",want:100},
  {name:"M2 Go scope",reader:GoCoverReader{},provider:ProviderGoCoverProfile,scope:[]string{"a/**"},data:"mode: set\nexample.com/m/a/a.go:1.1,2.1 2 1\nexample.com/m/b/b.go:1.1,2.1 2 0\n",want:100},
  {name:"M2 Go unsupported mode",reader:GoCoverReader{},provider:ProviderGoCoverProfile,data:"mode: garbage\nexample.com/m/a/a.go:1.1,2.1 2 1\n",reject:true},
  {name:"M2 Go bad coordinates",reader:GoCoverReader{},provider:ProviderGoCoverProfile,data:"mode: set\nexample.com/m/a/a.go:bad,bad 2 1\n",reject:true},
  {name:"M2 Vitest scope",reader:VitestCoverageReader{},provider:ProviderVitestV8,scope:[]string{"a/**"},data:`{"a/a.ts":{"s":{"0":1}},"b/b.ts":{"s":{"0":0}}}`,want:100},
  {name:"M2 Vitest null counter",reader:VitestCoverageReader{},provider:ProviderVitestV8,data:`{"a/a.ts":{"s":{"0":null}}}`,reject:true},
  {name:"M2 Vitest negative counter",reader:VitestCoverageReader{},provider:ProviderVitestV8,data:`{"a/a.ts":{"s":{"0":-1}}}`,reject:true},
  {name:"M2 coverage.py scoped count",reader:CoveragePyReader{},provider:ProviderCoveragePyJSON,scope:[]string{"a/**"},data:`{"totals":{"covered_lines":1,"num_statements":2},"files":{"a/a.py":{"summary":{"covered_lines":1,"num_statements":1}},"b/b.py":{"summary":{"covered_lines":0,"num_statements":1}}}}`,want:100},
  {name:"M2 coverage.py missing total",reader:CoveragePyReader{},provider:ProviderCoveragePyJSON,data:`{"totals":{"num_statements":1}}`,reject:true},
  {name:"M3 Lizard thresholds",reader:LizardReader{},provider:ProviderLizardCSV,data:"1,10,1,0,1,f,a.go,f,f(),1,1\n1,11,1,0,1,g,b.go,g,g(),1,1\n1,20,1,0,1,h,c.go,h,h(),1,1\n1,21,1,0,1,i,d.go,i,i(),1,1\n",want:21},
  {name:"M3 Lizard quoted Unicode",reader:LizardReader{},provider:ProviderLizardCSV,data:"1,22,1,0,1,loc,src/é.go,\"f,é\",f(),1,1\n",want:22},
  {name:"M4 jscpd all clones outside scope",reader:JscpdReader{},provider:ProviderJscpdJSON,scope:[]string{"a/**"},data:`{"statistics":{"total":{"lines":100,"duplicatedLines":20,"percentage":20}},"duplicates":[{"lines":20,"firstFile":{"name":"b/b.ts","start":1},"secondFile":{"name":"b/c.ts","start":1}}]}`,reject:true},
  {name:"M4 jscpd mixed scope affected path",reader:JscpdReader{},provider:ProviderJscpdJSON,scope:[]string{"a/**"},data:`{"statistics":{"total":{"lines":100,"duplicatedLines":20,"percentage":20}},"duplicates":[{"lines":20,"firstFile":{"name":"b/b.ts","start":1},"secondFile":{"name":"a/c.ts","start":1}}]}`,want:20},
  {name:"M5 OSV duplicate group",reader:OSVScannerReader{},provider:ProviderOSVScannerJSON,data:`{"results":[{"source":{"path":"go.mod"},"packages":[{"package":{"name":"x","version":"1"},"groups":[{"ids":["X"],"max_severity":"7"},{"ids":["X"],"max_severity":"7"}]}]}]}`,want:1},
  {name:"M5 OSV null source",reader:OSVScannerReader{},provider:ProviderOSVScannerJSON,data:`{"results":[null]}`,reject:true},
  {name:"M5 OSV scope",reader:OSVScannerReader{},provider:ProviderOSVScannerJSON,scope:[]string{"a/**"},data:`{"results":[{"source":{"path":"b/go.mod"},"packages":[{"package":{"name":"x","version":"1"},"groups":[{"ids":["X"],"max_severity":"7"}]}]}]}`},
 }
 for _, c := range cells { t.Run(c.name,func(t *testing.T){
  rd,err:=c.reader.Read(context.Background(),ReportInput{Root:root,Provider:c.provider,Data:[]byte(c.data),Scope:c.scope,Exclusions:c.exclusions})
  status:="PASS"; actual:="no value"; if rd.Value!=nil { actual=fmt.Sprintf("%g",rd.Value.Number) }
  if c.reject { if err==nil {status="ISSUE"} } else if err!=nil || rd.Value==nil || rd.Value.Number!=c.want || rd.Partial!=c.partial {status="ISSUE"}
  if c.name=="M4 jscpd mixed scope affected path" && (len(rd.Evidence)==0 || rd.Evidence[0].Path!="a/c.ts") {status="ISSUE"}
  t.Logf("%s expected value=%g reject=%v partial=%v; actual value=%s partial=%v error=%v evidence=%v",status,c.want,c.reject,c.partial,actual,rd.Partial,err,rd.Evidence)
 }) }
 // Alternate public surface: all-invalid Go report through Collect and store.
 p:=project(t); write(t,p,"go.mod","module example.com/m\n"); write(t,p,"report.jsonl","null\n")
 got:=collect(t,p,cfgOf(CapabilityConfig{Capability:CapabilityTests,Provider:ProviderGoTestJSON,Report:"report.jsonl"}),DefaultReaders(),nil)
 snaps,err:=NewStore(p).LoadSnapshots(); if err!=nil {t.Fatal(err)}
 for _,s:=range snaps[0].Summary.Capabilities { if s.Capability==CapabilityTests { t.Logf("M6 Collect null Go report: outcome=%s value=%v persisted classification=%s",got.Results[0].Result.Outcome,got.Results[0].Result.Value,s.Classification) } }
 // Exact CVSS boundaries and non-finite variants against a literal oracle.
 for score,want:=range map[string]string{"NaN":"unknown","+Inf":"unknown","-Inf":"unknown","0":"unknown","0.1":"low","3.9":"low","4":"medium","6.9":"medium","7":"high","8.9":"high","9":"critical","10":"critical","10.1":"unknown"} { got,_:=osvBucket(score); if got!=want {t.Errorf("M5 score %s = %s want %s",score,got,want)} }
 // Mutation and round trip: bytes are not changed by a reader.
 b:=[]byte(`{"results":[]}`); before:=string(b); rd,err:=OSVScannerReader{}.Read(context.Background(),ReportInput{Root:root,Data:b}); if err!=nil || string(b)!=before {t.Fatal("M5 mutable input",err)}; enc,_:=json.Marshal(rd); var back Reading; if json.Unmarshal(enc,&back)!=nil || back.Value.Number!=0 {t.Fatal("M5 roundtrip")}
 for _,n:=range []int{MaxEvidence-1,MaxEvidence,MaxEvidence+1} { var items []RankedEvidence; for i:=0;i<n;i++ {items=append(items,RankedEvidence{Ref:EvidenceRef{Path:fmt.Sprintf("src/f%d.go",i)},Rank:float64(i)})}; got:=WorstEvidence(items); if len(got)!=min(n,MaxEvidence) || got[0].Path!=fmt.Sprintf("src/f%d.go",n-1) {t.Fatal("M0 cap/order",n,got)} }
 t.Log("M0 exact evidence cap and neighbors PASS")
 // Isolation checked against a standalone baseline, with literal expected metrics.
 r:=project(t); write(t,r,"go.mod","module example.com/m\n"); write(t,r,"test.jsonl","{\"Action\":\"pass\",\"Package\":\"example.com/m/a\",\"Test\":\"OK\"}\n{\"Action\":\"pass\",\"Package\":\"example.com/m/a\"}\n"); write(t,r,"bad.xml","<");
 cfg:=cfgOf(CapabilityConfig{Capability:CapabilityTests,Provider:ProviderGoTestJSON,Report:"test.jsonl"})
 base:=collect(t,r,cfg,DefaultReaders(),nil)
 cfg.Capabilities=append(cfg.Capabilities,CapabilityConfig{Capability:CapabilityTests,Provider:ProviderVitestJUnit,Report:"bad.xml"})
 mixed:=collect(t,r,cfg,DefaultReaders(),nil)
 if base.Results[0].Result.Value.Number!=0 || mixed.Results[0].Result.Value.Number!=0 || detailMap(Reading{Details:mixed.Results[0].Result.Details})["total_tests"]!=1 || mixed.Results[1].Result.Outcome!=OutcomeFailed || mixed.Results[1].Result.Value!=nil {t.Fatal("M6 independent isolation",base,mixed)}
 t.Log("M6 baseline-versus-malformed sibling exact 0 failures/1 total preserved; repeat/store PASS")
}
```

## Matrix transcript

```text
=== RUN   TestO029IndependentMatrix
=== RUN   TestO029IndependentMatrix/M1_Go_scoped_count
    o029_check_test.go:49: ISSUE expected value=0 reject=false partial=false; actual value=1 partial=false error=<nil> evidence=[]
=== RUN   TestO029IndependentMatrix/M1_Go_null_event
    o029_check_test.go:49: ISSUE expected value=0 reject=true partial=false; actual value=0 partial=false error=<nil> evidence=[]
=== RUN   TestO029IndependentMatrix/M1_Go_wrong_schema
    o029_check_test.go:49: ISSUE expected value=0 reject=true partial=false; actual value=0 partial=false error=<nil> evidence=[]
=== RUN   TestO029IndependentMatrix/M1_Go_missing_final_package
    o029_check_test.go:49: ISSUE expected value=0 reject=false partial=true; actual value=0 partial=false error=<nil> evidence=[]
=== RUN   TestO029IndependentMatrix/M1_Go_incomplete_start_only
    o029_check_test.go:49: ISSUE expected value=0 reject=false partial=true; actual value=0 partial=false error=<nil> evidence=[]
=== RUN   TestO029IndependentMatrix/M1_Go_build-fail_event
    o029_check_test.go:49: ISSUE expected value=1 reject=false partial=false; actual value=0 partial=false error=<nil> evidence=[]
=== RUN   TestO029IndependentMatrix/M1_Vitest_scoped_count
    o029_check_test.go:49: ISSUE expected value=0 reject=false partial=false; actual value=1 partial=false error=<nil> evidence=[]
=== RUN   TestO029IndependentMatrix/M1_Pytest_excluded_failure
    o029_check_test.go:49: ISSUE expected value=0 reject=false partial=false; actual value=1 partial=false error=<nil> evidence=[]
=== RUN   TestO029IndependentMatrix/M1_JUnit_mismatched_declared_total
    o029_check_test.go:49: ISSUE expected value=0 reject=false partial=true; actual value=0 partial=false error=<nil> evidence=[]
=== RUN   TestO029IndependentMatrix/M1_JUnit_wrong_root_bypass
    o029_check_test.go:49: ISSUE expected value=0 reject=true partial=false; actual value=0 partial=false error=<nil> evidence=[]
=== RUN   TestO029IndependentMatrix/M1_JUnit_two_roots
    o029_check_test.go:49: ISSUE expected value=0 reject=true partial=false; actual value=0 partial=false error=<nil> evidence=[]
=== RUN   TestO029IndependentMatrix/M2_Go_merged_duplicate
    o029_check_test.go:49: PASS expected value=100 reject=false partial=false; actual value=100 partial=false error=<nil> evidence=[]
=== RUN   TestO029IndependentMatrix/M2_Go_scope
    o029_check_test.go:49: PASS expected value=100 reject=false partial=false; actual value=100 partial=false error=<nil> evidence=[]
=== RUN   TestO029IndependentMatrix/M2_Go_unsupported_mode
    o029_check_test.go:49: ISSUE expected value=0 reject=true partial=false; actual value=100 partial=false error=<nil> evidence=[]
=== RUN   TestO029IndependentMatrix/M2_Go_bad_coordinates
    o029_check_test.go:49: ISSUE expected value=0 reject=true partial=false; actual value=100 partial=false error=<nil> evidence=[]
=== RUN   TestO029IndependentMatrix/M2_Vitest_scope
    o029_check_test.go:49: PASS expected value=100 reject=false partial=false; actual value=100 partial=false error=<nil> evidence=[]
=== RUN   TestO029IndependentMatrix/M2_Vitest_null_counter
    o029_check_test.go:49: ISSUE expected value=0 reject=true partial=false; actual value=0 partial=false error=<nil> evidence=[{a/a.ts 0 0.0% of 1 statements covered}]
=== RUN   TestO029IndependentMatrix/M2_Vitest_negative_counter
    o029_check_test.go:49: ISSUE expected value=0 reject=true partial=false; actual value=0 partial=false error=<nil> evidence=[{a/a.ts 0 0.0% of 1 statements covered}]
=== RUN   TestO029IndependentMatrix/M2_coverage.py_scoped_count
    o029_check_test.go:49: ISSUE expected value=100 reject=false partial=false; actual value=50 partial=false error=<nil> evidence=[]
=== RUN   TestO029IndependentMatrix/M2_coverage.py_missing_total
    o029_check_test.go:49: PASS expected value=0 reject=true partial=false; actual value=no value partial=false error=coverage.py JSON has no totals for statements evidence=[]
=== RUN   TestO029IndependentMatrix/M3_Lizard_thresholds
    o029_check_test.go:49: PASS expected value=21 reject=false partial=false; actual value=21 partial=false error=<nil> evidence=[{d.go 1 i has complexity 21} {c.go 1 h has complexity 20} {b.go 1 g has complexity 11}]
=== RUN   TestO029IndependentMatrix/M3_Lizard_quoted_Unicode
    o029_check_test.go:49: PASS expected value=22 reject=false partial=false; actual value=22 partial=false error=<nil> evidence=[{src/é.go 1 f,é has complexity 22}]
=== RUN   TestO029IndependentMatrix/M4_jscpd_all_clones_outside_scope
    o029_check_test.go:49: ISSUE expected value=0 reject=true partial=false; actual value=20 partial=false error=<nil> evidence=[]
=== RUN   TestO029IndependentMatrix/M4_jscpd_mixed_scope_affected_path
    o029_check_test.go:49: ISSUE expected value=20 reject=false partial=false; actual value=20 partial=false error=<nil> evidence=[{b/b.ts 1 20 lines repeated at a/c.ts:1}]
=== RUN   TestO029IndependentMatrix/M5_OSV_duplicate_group
    o029_check_test.go:49: ISSUE expected value=1 reject=false partial=false; actual value=2 partial=false error=<nil> evidence=[{go.mod 0 x 1: X} {go.mod 0 x 1: X}]
=== RUN   TestO029IndependentMatrix/M5_OSV_null_source
    o029_check_test.go:49: ISSUE expected value=0 reject=true partial=false; actual value=0 partial=true error=<nil> evidence=[]
=== RUN   TestO029IndependentMatrix/M5_OSV_scope
    o029_check_test.go:49: PASS expected value=0 reject=false partial=false; actual value=0 partial=false error=<nil> evidence=[]
=== NAME  TestO029IndependentMatrix
    o029_check_test.go:55: M6 Collect null Go report: outcome=available value=&{0 count} persisted classification=watch
    o029_check_test.go:61: M0 exact evidence cap and neighbors PASS
    o029_check_test.go:69: M6 baseline-versus-malformed sibling exact 0 failures/1 total preserved; repeat/store PASS
--- PASS: TestO029IndependentMatrix (0.08s)
    --- PASS: TestO029IndependentMatrix/M1_Go_scoped_count (0.00s)
    --- PASS: TestO029IndependentMatrix/M1_Go_null_event (0.00s)
    --- PASS: TestO029IndependentMatrix/M1_Go_wrong_schema (0.00s)
    --- PASS: TestO029IndependentMatrix/M1_Go_missing_final_package (0.00s)
    --- PASS: TestO029IndependentMatrix/M1_Go_incomplete_start_only (0.00s)
    --- PASS: TestO029IndependentMatrix/M1_Go_build-fail_event (0.00s)
    --- PASS: TestO029IndependentMatrix/M1_Vitest_scoped_count (0.00s)
    --- PASS: TestO029IndependentMatrix/M1_Pytest_excluded_failure (0.00s)
    --- PASS: TestO029IndependentMatrix/M1_JUnit_mismatched_declared_total (0.00s)
    --- PASS: TestO029IndependentMatrix/M1_JUnit_wrong_root_bypass (0.00s)
    --- PASS: TestO029IndependentMatrix/M1_JUnit_two_roots (0.00s)
    --- PASS: TestO029IndependentMatrix/M2_Go_merged_duplicate (0.00s)
    --- PASS: TestO029IndependentMatrix/M2_Go_scope (0.00s)
    --- PASS: TestO029IndependentMatrix/M2_Go_unsupported_mode (0.00s)
    --- PASS: TestO029IndependentMatrix/M2_Go_bad_coordinates (0.00s)
    --- PASS: TestO029IndependentMatrix/M2_Vitest_scope (0.00s)
    --- PASS: TestO029IndependentMatrix/M2_Vitest_null_counter (0.00s)
    --- PASS: TestO029IndependentMatrix/M2_Vitest_negative_counter (0.00s)
    --- PASS: TestO029IndependentMatrix/M2_coverage.py_scoped_count (0.00s)
    --- PASS: TestO029IndependentMatrix/M2_coverage.py_missing_total (0.00s)
    --- PASS: TestO029IndependentMatrix/M3_Lizard_thresholds (0.00s)
    --- PASS: TestO029IndependentMatrix/M3_Lizard_quoted_Unicode (0.00s)
    --- PASS: TestO029IndependentMatrix/M4_jscpd_all_clones_outside_scope (0.00s)
    --- PASS: TestO029IndependentMatrix/M4_jscpd_mixed_scope_affected_path (0.00s)
    --- PASS: TestO029IndependentMatrix/M5_OSV_duplicate_group (0.00s)
    --- PASS: TestO029IndependentMatrix/M5_OSV_null_source (0.00s)
    --- PASS: TestO029IndependentMatrix/M5_OSV_scope (0.00s)
PASS
ok  	github.com/opencode/savepoint/internal/codehealth	0.084s
```
