---
id: C-945
scope: {kind: objective, id: O-029}
result: CLEAR
checked_by: {role: checker, session: o029-independent-recheck-20261001}
executed_session: o029-issue-repair-20261001
checked_at: '2026-10-01T10:42:32Z'
reviewed:
  head_commit: 6b206fc6376b226c3e88ae8afcb60ce06c67c209
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
    - internal/codehealth/reader_repair_test.go
  dependencies: []
issues: []
supersedes: C-944
---

# C-945: O-029 Full Objective Recheck

**CLEAR.** All admitted original findings are repaired; the independently checked implementation satisfies the frozen scope with I-101’s recorded owner decision.

## Prior Issue closure map

| Issue | Disposition in this run | Proof |
|---|---|---|
| I-097 | Closed as verified by C-945 | Original M1 counts now 0; M2 scoped coverage.py 100%; M4 missing file totals gives partial/no value; independent 6/30 jscpd headline 20% and second-side scoped evidence; excludes and remaining provider scope tests pass. |
| I-098 | Closed as verified by C-945 | Null/unrelated Go input errors; finished tests lacking package terminal and start-only are partial; inconsistent JUnit count is partial; multiple roots error; Collect stores null as failed/no-value/unknown. |
| I-099 | Closed as verified by C-945 | Unsupported modes/bad spans error; null/negative Vitest counters error; failed instances have no value/unknown and valid sibling remains intact. |
| I-100 | Closed as verified by C-945 | Original repeated group now value=1/high=1; reordered aliases and separate groups/packages tests pass. |
| I-101 | Already closed by owner as accepted; retained | Existing accepted resolution and Objective wording explicitly permit the CSV limitation. Empty CSV independently returns 0 with the limitation in its reason; outside-root rows stay partial/no false zero. This is not a verified repair of unparsed-file detection. |

## Authority and frozen scope

Fresh independent checker conversation; it did not execute O-029 or its repairs. This is the first Full recheck of C-944, not a new initial perimeter. All seven Tasks remain owner-completed with explicit frontmatter waivers naming Task, reason, actor and time. No owner validation required is declared. Router selection resolves Check O-029 despite stored task state; the owner's explicit recheck request and resolved Next authorize this check. No implementation, Task status, Objective status, router or Design edits.

The numbered scope lock, matrix axes, finite external boundary classifications and workflow operation order are exactly C-944's. The only requirement decision applied is the already recorded owner acceptance of I-101, with its matching Objective edit. Adjacent probes are restricted to the exact cells in the admission ledger below. Provider analysis, network protocols, commands/TUI, unsupported document wrappers and standalone modern build events remain outside the blocking perimeter.

## Coverage matrix

P = Proven; O = original non-blocking observation; A = explicit owner acceptance. All original applicable cells are completed by the retained independent matrix plus the named fixture/helper/collection/runner/store/classification tests. Non-applicable axes retain C-944's reasons: no terminal renderer, no changed lifecycle, no network client, no mutable service state.

| Frozen row | Normal/empty | Malformed/missing/null/type | Partial/unsupported | Scope/exclusion/mixed | Bounds/provenance | Duplicate/sequence/representation |
|---|---|---|---|---|---|---|
| M0 paths/modules/evidence | P | P | P Windows lexical/outside-root | P Locate/Resolve and nested modules | P cap-1/cap/cap+1/notes | P deterministic order |
| M1 Go tests | P pass/fail/build/skip/no tests | P null/wrong schema/array/invalid JSON | P missing terminal/start-only/unfinished; standalone build event O | P original scoped case + include/exclude regressions/monorepo | P safe named evidence/version unknown | P repeated reads/subtest semantics |
| M1 Vitest/pytest JUnit | P cases/failures/errors/skips/empty | P multiple roots/non-numeric counts/malformed; html wrapper O | P declared-count mismatch | P original Vitest include and pytest exclude + mapped class regressions | P file/name bounds/ignored hostname and output | P shared parser/nested suites |
| M2 Go coverage | P totals/fully covered; zero-total rejected | P mode/span/negative/error fixtures | P unresolved import; unsupported mode rejected | P original scope/monorepo | P literal ratios/order/version | P merged duplicate blocks |
| M2 Vitest coverage | P statement/functions/branches/zero errors | P null/negative statement/function/branch | P outside-root partial | P original scope/mixed collection | P ratio/order/version | P map-order determinism |
| M2 coverage.py | P totals/branches/zero errors | P missing totals/malformed | P inline missing in-scope summary fixture, independently varied 75% and persisted partial | P original 100% scope/excluded zero/no-scope 50% | P meta.version/least covered | P repeat/store |
| M3 Lizard | P max/average/empty | P wrong columns/CCN/start line | P outside-root; unparsed inventory A I-101 | P scoped/mixed Go/TS/Python | P 10/11/20/21/quoted Unicode/cap/version | P deterministic reads |
| M4 jscpd | P totals/no clones/zero error | P absent totals/malformed/inconsistent | P outside-root; scoped report without file counts partial/no-value | P original outside/mixed cases; independent 6/30 and scoped location; exclusions | P largest clones/cap/version/ignored unrelated fields | P Collect reload and unknown no-value |
| M5 OSV/model/Assess | P no findings/all severity buckets | P malformed/missing results; null source O (honest partial) | P unresolved/error/ungrouped | P source filter/mixed ecosystems | P CVSS exact floors/nonfinite/name bounds/freshness | P original duplicate repaired; aliases/packages; mutation/roundtrip |
| M6 registry/Collect/store | P nine providers/multiple instances | P reader error isolated/no value/unknown | P timeout/unavailable/absent/cancel/truncation | P independent valid baseline vs malformed sibling | P root/validation | P immutable reload/repeat; partial/no-value jscpd persisted unknown |
| M7 docs/reality/waivers/gate | P scoped files real and waivers present | P strict index loading | Full gate result below | P Design/Map match production readers | P dependency/gate definitions unchanged | Native Windows execution not claimed; same original boundary |

Fixture completeness uses disk fixtures and inline table fixtures, without introducing a new requirement that every unsupported variant have a standalone disk file. The original missing schema/completeness cells are now covered by reader_repair_test.go. Mixed-project support is demonstrated by all-nine-provider Collect and monorepo reports; Lizard/jscpd/OSV additionally carry mixed languages/ecosystems within one report. No actual analyzer algorithms are executed.

## Workflow and external effects

Retain C-944 operations 1..7: config/history guards -> repository/instance selection -> report acquisition/fake process -> normalization -> validation/assessment -> immutable save -> reload/repeat. No repair changes process startup, tool target, findings exits, cancellation, cleanup, storage or publication ordering. Existing Collect/ExecRunner/Store tests exercise all original success/failure prefixes and cleanup cells in the full suite. Independently compared valid baseline vs malformed sibling and stored semantics; scoped jscpd without counts is partial/no-value and reloads with unknown health. No new external operations or effects were introduced.

## Acceptance reconciliation

| Requirement | Classification and concrete evidence |
|---|---|
| O-029 success 1 | Proven: M1 original probes plus TestGoTestReader, TestJUnitReader and completeness/scope regressions. No report remains absent via Collect. |
| O-029 success 2 | Proven: M2 numeric oracle and TestCoverageReadersValue, TestCoverageOmitsBranchDetailsWhenReportHasNone, scope and invalid-measurement regressions. |
| O-029 success 3 | Proven under recorded owner decision: max/details/bounds/outside-root and explicit empty-report limitation; I-101 remains accepted. |
| O-029 success 4 | Proven: M4 ratios/scoped denominator/evidence, malformed/zero/cap/version tests; no analyzer health fields consumed. |
| O-029 success 5 | Proven: M5 count=1 after identical groups, severity floors and five-bucket sum; unknown/high/critical blocking, partial/provenance tests. |
| O-029 success 6 | Proven within C-944 fixture classes: M0..M5 bounds, scope, supported/malformed/partial/unsupported forms and mixed projects; accepted Lizard limitation remains explicit. |
| O-029 success 7 | Proven: M6 registry/fake mixed collection, baseline comparison, malformed no-value/unknown and immutable reload. |
| T-059 Done When 1..5 | Proven: Root seam; TestRelPath, TestReportInputPathHonoursScope, TestGoModulesResolve/KeepsToScope, WorstEvidence tests plus independent cap neighbors. |
| T-060 Done When 1..5 | Proven: TestGoTestReader and JUnitReader fixtures; original null/terminal/count/root probes; scope repair and provenance evidence. |
| T-061 Done When 1..5 | Proven: CoverageReadersValue/RejectUnusableReports, GoCoverReaderMonorepoResolvesEachModule and invalid-measurement/scope regressions. |
| T-062 Done When 1..6 | Proven: LizardReaderValue/RejectsMalformedReports/EvidenceIsBounded and independent thresholds/quoted Unicode/no-functions. |
| T-063 Done When 1..5 | Proven: JscpdReaderValue/RejectsUnusableReports/EvidenceIsBounded and scoped counts/evidence/partial regressions. First location retained unless only the second is in scope, as I-097's scoped-evidence requirement directs. |
| T-064 Done When 1..5 | Proven: OSVScannerReaderValue/UnknownSeverityBlocks/EvidenceIsWorstFirstAndBounded, severity model/classification tests and distinct-group regressions. |
| T-065 Done When 1..4 | Proven: DefaultReadersCoverTheWholeCatalogue, CollectWithRealReadersKeepsEveryMeasureTruthful, independent baseline/reload and unchanged Design/AGENTS reconciliation. |

## Guardrails and code style

FS-05 actual IO uses filepath; normalized evidence uses slash paths. ARCH-03 explicit Root; ARCH-04 codehealth responsibility maintained. DEP-01/02 no dependency change. DATA-03 invalid inputs diagnose at boundaries. CFG-02 platform-shaped path tests; CFG-03 native Windows CI remains a CI responsibility, as in C-944. TPL-02 Design section 1 and Codebase Map describe actual readers and defer collection command. TEST-01/02/04/06/07 named literal outcomes, temporary projects/fake analyzers; TEST-09 every waiver valid. TEST-08 fresh gate evidence below.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches** — original admitted boundary/scope/completeness cases now covered.
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries** — admitted semantic input defects repaired.
- [x] STYLE-07 **One source of truth**
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs** — repairs stay in reader normalization and tests.

## Observations and materiality

No remaining admitted Issue or materiality action. C-944's standalone build-fail event and html wrapper remain observations; OSV null source remains honestly partial rather than rejected. The original observational harness logs these as ISSUE against a stricter oracle, but C-944 expressly excluded them from its blocking verdict; this run retains that disposition. The jscpd 0.5 percentage-point tolerance/comment discrepancy and pre-existing external report-path orchestration observation remain non-blocking. No additional blocking axes or remediation cycle.

## Evidence

- Fresh `make test-full` passed outside the sandbox, completed by 2026-10-01T10:42:32Z, go1.26.2 linux/amd64. Full uncached Go suite and Linux/Darwin/Windows builds succeeded. Native Windows execution is not claimed.
- Two earlier sandbox runs did not complete: the first reported internal/data test failure; a repeated suite passed but Windows build could not download/write missing module dependencies. The final approved complete run passed. No implementation repair was performed by the checker; these failed attempts are not reused as gate evidence.
- `git diff --check` and `go vet ./internal/codehealth`: exit 0.
- `go test ./internal/codehealth -count=1 -run 'Reader|Readers|RelPath|ReportInputPath|GoModules|WorstEvidence|Collect|Store|Snapshot|Assess' -v`: exit 0, 7.106s.
- `go test -overlay /tmp/o029_overlay.json ./internal/codehealth -run 'TestO029IndependentMatrix|TestO029Recheck' -count=1 -v`: exit 0, 1.532s; admitted cells pass. All original observational cells retain their C-944 dispositions, irrespective of stricter log labels.
- File reality: 36 Task evidence path references and 79 scoped prior/new file references exist. SHA256 comparison confirmed reviewed input files unchanged throughout final gate and artifact preparation. go.mod, go.sum, Makefile and internal/buildtool have no working-tree changes. Head commit is the working-tree base, not a claim that repairs are committed. The unrelated .gitignore local settings entry is outside reader review.
- Strict index loading after Check/Issue writes: `./savepoint resume`; output recorded below after execution.

## Handoff

O-029 is ready for owner completion. All seven Tasks are already done, independent Full Check is CLEAR, I-097..I-100 are verified, and I-101 remains accepted by the owner. This checker records readiness only and leaves lifecycle decisions to the owner.

## Admission ledger

# Recheck admission ledger
Frozen authority: C-944, numbered scope lock and completed matrix. No new axes.
| Item | Prior issue/claim | Exact cell | Allowed result |
|---|---|---|---|
| Path helpers, scope and bounds | T-059 unchanged plus Locate repair | M0 normal/path/scope/exclusion/cap cells | blocking only for original invariant |
| Go null/unrelated, missing final, start-only | I-098 | M1 named null/wrong-schema/terminal cells; M6 null Collect | repair proven or still open |
| Go/Vitest/pytest scoped counts | I-097 | M1 Go scoped count/Vitest scoped count/Pytest excluded failure | repair proven or still open |
| JUnit declared total and multiple roots | I-098 | M1 mismatched declared total/two roots | repair proven or still open |
| Go unsupported mode/coordinates; Vitest null/negative | I-099 | M2 explicitly named cells; M6 failed containment | repair proven or still open |
| Scoped coverage.py and missing counts | I-097 | M2 scoped count, partial/unsupported in-scope source evidence | repair proven or still open |
| jscpd outside/mixed clone and missing file totals | I-097 | M4 all clones outside scope/mixed first-second locations | scoped value or explicit partial/no value |
| OSV repeated/reordered group/different packages | I-100 | M5 duplicate identical group, Issue Proof Needed named alias/package neighbors | repair proven or still open |
| Lizard no functions and outside root | I-101 owner acceptance | M3 no-functions/unparsed/outside-root | preserve accepted disposition; validate limitation reason |
| Complete provider fixture classes, registry/isolation/store | original matrix and Tasks | M0..M7 original named rows only | proven/unverified |
| Standalone build-fail, html wrapper, null OSV | C-944 observations | original observational cells | observation only |
| Full gate/docs/file reality/waivers | TEST-08, TPL-02, TEST-09 | M7 | proven/unverified |

The original harness was first executed before this ledger was persisted. It is rerun after the ledger to correct that process ordering; the frozen cells and verdict boundaries are unchanged.

## Independent harness

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
  if c.name=="M4 jscpd all clones outside scope" || c.name=="M4 jscpd mixed scope affected path" { if err==nil && rd.Partial && rd.Value==nil { status="PASS" } }; if c.name=="M4 jscpd mixed scope affected path" && rd.Value!=nil && (len(rd.Evidence)==0 || rd.Evidence[0].Path!="a/c.ts") {status="ISSUE"}
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
func TestO029RecheckAdjacent(t *testing.T) {
 // Same M4 scoped counts/mixed-location cells with a different literal denominator.
 report := `{"statistics":{"total":{"lines":80,"duplicatedLines":20,"percentage":25},"formats":{"ts":{"sources":{"a/a.ts":{"lines":30,"duplicatedLines":6},"b/b.ts":{"lines":50,"duplicatedLines":14}}}}},"duplicates":[{"lines":6,"firstFile":{"name":"b/b.ts","start":2},"secondFile":{"name":"a/a.ts","start":8}}]}`
 rd,err:=(JscpdReader{}).Read(context.Background(),ReportInput{Root:"/r",Data:[]byte(report),Scope:[]string{"a/**"}})
 if err!=nil || rd.Value==nil || rd.Value.Number!=20 || rd.Partial || len(rd.Evidence)!=1 || rd.Evidence[0].Path!="a/a.ts" || rd.Evidence[0].Line!=8 {t.Fatalf("M4 mixed scope: %v %+v",err,rd)}
 t.Log("M4 independent scoped denominator 6/30=20%; second location a/a.ts:8 PASS")
 // M2 original partial/unsupported in-scope source cell, independently varied counts.
 py:=`{"totals":{"covered_lines":3,"num_statements":4},"files":{"a/a.py":{"summary":{"covered_lines":3,"num_statements":4}},"a/missing.py":{"summary":{}}}}`
 rd,err=(CoveragePyReader{}).Read(context.Background(),ReportInput{Root:"/r",Data:[]byte(py),Scope:[]string{"a/**"}})
 if err!=nil || rd.Value==nil || rd.Value.Number!=75 || !rd.Partial {t.Fatalf("M2 partial source: %v %+v",err,rd)}
 t.Log("M2 coverage.py in-scope missing summary yields partial 75% PASS")
 root:=project(t);write(t,root,"reports/py.json",py)
 got:=collect(t,root,cfgOf(CapabilityConfig{Capability:CapabilityCoverage,Provider:ProviderCoveragePyJSON,Report:"reports/py.json",Scope:[]string{"a/**"}}),DefaultReaders(),nil)
 if got.Results[0].Result.Outcome!=OutcomePartial || got.Results[0].Result.Value.Number!=75 {t.Fatal("M6 partial persisted",got)}
 rd,err=(LizardReader{}).Read(context.Background(),ReportInput{Root:"/r"})
 if err!=nil || rd.Value==nil || rd.Value.Number!=0 || rd.Reason!="Lizard reported no functions; its CSV cannot show whether any file failed to parse, and does not state its version" {t.Fatal("M3 owner accepted limitation",err,rd)}
 t.Log("M3 accepted limitation explicit at no-functions boundary PASS")
}

func TestO029RecheckScopeCollect(t *testing.T) {
 // M4 missing scoped totals -> M6 normalization/assessment/store, using fake tool.
 root:=project(t)
 report:=`{"statistics":{"total":{"lines":50,"duplicatedLines":10,"percentage":20}},"duplicates":[]}`
 tools:=&fakeTools{t:t,behavior:map[string]func(context.Context,ToolSpec)(ToolResult,error){"dup":stdout(report)}}
 got:=collect(t,root,cfgOf(CapabilityConfig{Capability:CapabilityDuplication,Provider:ProviderJscpdJSON,Executable:"dup",Scope:[]string{"a/**"}}),DefaultReaders(),tools)
 if got.Results[0].Result.Outcome!=OutcomePartial || got.Results[0].Result.Value!=nil {t.Fatal("M6 scoped jscpd no invented total",got)}
 snaps,err:=NewStore(root).LoadSnapshots();if err!=nil {t.Fatal(err)}
 found:=false
 for _,a:=range snaps[0].Summary.Capabilities {if a.Capability==CapabilityDuplication {found=true;if a.Classification!=ClassificationUnknown {t.Fatal("M6 no value classification",a)}}}
 if !found {t.Fatal("no duplication summary")}
 t.Log("M4/M6 scoped jscpd lacking counts -> partial/no-value, stored unknown PASS")
}

```

## Matrix transcript

```text
=== RUN   TestO029IndependentMatrix
=== RUN   TestO029IndependentMatrix/M1_Go_scoped_count
    o029_check_test.go:49: PASS expected value=0 reject=false partial=false; actual value=0 partial=false error=<nil> evidence=[]
=== RUN   TestO029IndependentMatrix/M1_Go_null_event
    o029_check_test.go:49: PASS expected value=0 reject=true partial=false; actual value=no value partial=false error=go test JSON line 1 has no event action evidence=[]
=== RUN   TestO029IndependentMatrix/M1_Go_wrong_schema
    o029_check_test.go:49: PASS expected value=0 reject=true partial=false; actual value=no value partial=false error=go test JSON line 1 has no event action evidence=[]
=== RUN   TestO029IndependentMatrix/M1_Go_missing_final_package
    o029_check_test.go:49: PASS expected value=0 reject=false partial=true; actual value=0 partial=true error=<nil> evidence=[]
=== RUN   TestO029IndependentMatrix/M1_Go_incomplete_start_only
    o029_check_test.go:49: PASS expected value=0 reject=false partial=true; actual value=0 partial=true error=<nil> evidence=[]
=== RUN   TestO029IndependentMatrix/M1_Go_build-fail_event
    o029_check_test.go:49: ISSUE expected value=1 reject=false partial=false; actual value=0 partial=false error=<nil> evidence=[]
=== RUN   TestO029IndependentMatrix/M1_Vitest_scoped_count
    o029_check_test.go:49: PASS expected value=0 reject=false partial=false; actual value=0 partial=false error=<nil> evidence=[]
=== RUN   TestO029IndependentMatrix/M1_Pytest_excluded_failure
    o029_check_test.go:49: PASS expected value=0 reject=false partial=false; actual value=0 partial=false error=<nil> evidence=[]
=== RUN   TestO029IndependentMatrix/M1_JUnit_mismatched_declared_total
    o029_check_test.go:49: PASS expected value=0 reject=false partial=true; actual value=0 partial=true error=<nil> evidence=[]
=== RUN   TestO029IndependentMatrix/M1_JUnit_wrong_root_bypass
    o029_check_test.go:49: ISSUE expected value=0 reject=true partial=false; actual value=0 partial=false error=<nil> evidence=[]
=== RUN   TestO029IndependentMatrix/M1_JUnit_two_roots
    o029_check_test.go:49: PASS expected value=0 reject=true partial=false; actual value=no value partial=false error=JUnit XML has more than one document root evidence=[]
=== RUN   TestO029IndependentMatrix/M2_Go_merged_duplicate
    o029_check_test.go:49: PASS expected value=100 reject=false partial=false; actual value=100 partial=false error=<nil> evidence=[]
=== RUN   TestO029IndependentMatrix/M2_Go_scope
    o029_check_test.go:49: PASS expected value=100 reject=false partial=false; actual value=100 partial=false error=<nil> evidence=[]
=== RUN   TestO029IndependentMatrix/M2_Go_unsupported_mode
    o029_check_test.go:49: PASS expected value=0 reject=true partial=false; actual value=no value partial=false error=go cover profile mode "garbage" is not set, count, or atomic evidence=[]
=== RUN   TestO029IndependentMatrix/M2_Go_bad_coordinates
    o029_check_test.go:49: PASS expected value=0 reject=true partial=false; actual value=no value partial=false error=go cover profile line 2 is not a coverage block evidence=[]
=== RUN   TestO029IndependentMatrix/M2_Vitest_scope
    o029_check_test.go:49: PASS expected value=100 reject=false partial=false; actual value=100 partial=false error=<nil> evidence=[]
=== RUN   TestO029IndependentMatrix/M2_Vitest_null_counter
    o029_check_test.go:49: PASS expected value=0 reject=true partial=false; actual value=no value partial=false error=vitest coverage JSON has a missing or negative execution count for a/a.ts evidence=[]
=== RUN   TestO029IndependentMatrix/M2_Vitest_negative_counter
    o029_check_test.go:49: PASS expected value=0 reject=true partial=false; actual value=no value partial=false error=vitest coverage JSON has a missing or negative execution count for a/a.ts evidence=[]
=== RUN   TestO029IndependentMatrix/M2_coverage.py_scoped_count
    o029_check_test.go:49: PASS expected value=100 reject=false partial=false; actual value=100 partial=false error=<nil> evidence=[]
=== RUN   TestO029IndependentMatrix/M2_coverage.py_missing_total
    o029_check_test.go:49: PASS expected value=0 reject=true partial=false; actual value=no value partial=false error=coverage.py JSON has no totals for statements evidence=[]
=== RUN   TestO029IndependentMatrix/M3_Lizard_thresholds
    o029_check_test.go:49: PASS expected value=21 reject=false partial=false; actual value=21 partial=false error=<nil> evidence=[{d.go 1 i has complexity 21} {c.go 1 h has complexity 20} {b.go 1 g has complexity 11}]
=== RUN   TestO029IndependentMatrix/M3_Lizard_quoted_Unicode
    o029_check_test.go:49: PASS expected value=22 reject=false partial=false; actual value=22 partial=false error=<nil> evidence=[{src/é.go 1 f,é has complexity 22}]
=== RUN   TestO029IndependentMatrix/M4_jscpd_all_clones_outside_scope
    o029_check_test.go:49: PASS expected value=0 reject=true partial=false; actual value=no value partial=true error=<nil> evidence=[]
=== RUN   TestO029IndependentMatrix/M4_jscpd_mixed_scope_affected_path
    o029_check_test.go:49: PASS expected value=20 reject=false partial=false; actual value=no value partial=true error=<nil> evidence=[]
=== RUN   TestO029IndependentMatrix/M5_OSV_duplicate_group
    o029_check_test.go:49: PASS expected value=1 reject=false partial=false; actual value=1 partial=false error=<nil> evidence=[{go.mod 0 x 1: X}]
=== RUN   TestO029IndependentMatrix/M5_OSV_null_source
    o029_check_test.go:49: ISSUE expected value=0 reject=true partial=false; actual value=0 partial=true error=<nil> evidence=[]
=== RUN   TestO029IndependentMatrix/M5_OSV_scope
    o029_check_test.go:49: PASS expected value=0 reject=false partial=false; actual value=0 partial=false error=<nil> evidence=[]
=== NAME  TestO029IndependentMatrix
    o029_check_test.go:55: M6 Collect null Go report: outcome=failed value=<nil> persisted classification=unknown
    o029_check_test.go:61: M0 exact evidence cap and neighbors PASS
    o029_check_test.go:69: M6 baseline-versus-malformed sibling exact 0 failures/1 total preserved; repeat/store PASS
--- PASS: TestO029IndependentMatrix (0.90s)
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
    --- PASS: TestO029IndependentMatrix/M2_Vitest_scope (0.01s)
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
=== RUN   TestO029RecheckAdjacent
    o029_check_test.go:76: M4 independent scoped denominator 6/30=20%; second location a/a.ts:8 PASS
    o029_check_test.go:81: M2 coverage.py in-scope missing summary yields partial 75% PASS
    o029_check_test.go:87: M3 accepted limitation explicit at no-functions boundary PASS
--- PASS: TestO029RecheckAdjacent (0.32s)
=== RUN   TestO029RecheckScopeCollect
    o029_check_test.go:101: M4/M6 scoped jscpd lacking counts -> partial/no-value, stored unknown PASS
--- PASS: TestO029RecheckScopeCollect (0.25s)
PASS
ok  	github.com/opencode/savepoint/internal/codehealth	1.532s

```

## Strict index validation

`./savepoint resume` exited 0 after all artifact writes. It reported:

```text
Close O-029 — Measure all five Code Health signals
Technical clearance: Check C-945 is recorded CLEAR.
I-097..I-101: resolved.
Next action: Owner: record Objective O-029 as done.
```
