---
id: C-908
scope: {kind: objective, id: O-016}
result: NEEDS WORK
checked_by: {role: checker, session: o016-full-check-20260923}
executed_session: o016-task-series-20260923
checked_at: '2026-09-23T01:42:44Z'
reviewed:
  base_commit: e7e72bb3d736028d208bc070f012fd6c66e8c930
  head_commit: e7e72bb3d736028d208bc070f012fd6c66e8c930
  files:
    - .github/workflows/ci.yml
    - Makefile
    - internal/buildtool/main.go
    - internal/buildtool/main_test.go
    - internal/migrate/apply_test.go
    - internal/migrate/command_test.go
    - internal/migrate/cutover_test.go
    - internal/migrate/end_to_end_test.go
    - internal/migrate/operation_test.go
    - AGENTS.md
    - .savepoint/Guardrails.md
    - agent-skills/savepoint-task/SKILL.md
    - agent-skills/savepoint-check/SKILL.md
    - agent-skills/savepoint-design/SKILL.md
    - templates/project-v2/AGENTS.md
    - templates/project-v2/agent-skills/savepoint-task/SKILL.md
    - templates/project-v2/agent-skills/savepoint-check/SKILL.md
    - templates/project-v2/agent-skills/savepoint-design/SKILL.md
    - .savepoint/objectives/O-016-fast-test-feedback/Objective.md
    - .savepoint/objectives/O-016-fast-test-feedback/tasks/T-013-measure-and-design-test-gates.md
    - .savepoint/objectives/O-016-fast-test-feedback/tasks/T-014-share-read-only-migration-results.md
    - .savepoint/objectives/O-016-fast-test-feedback/tasks/T-015-parallelize-safe-test-groups.md
    - .savepoint/objectives/O-016-fast-test-feedback/tasks/T-016-roll-out-fast-and-full-gates.md
  dependencies: []
issues: [I-029, I-030]
supersedes: null
---

# C-908: O-016 Full Objective Check

## Verdict

`NEEDS WORK`. Linux focused, fast, and full gates are correctly selected and
pass with useful timing output, but the complete Windows suite configured by
T-016 is known to fail. I-029 records O-016's duplicate Windows `TestMain`; I-030
records the two additional Windows platform failures exposed by the new job.
O-016 is not ready for owner closure.

## Evidence Mode

Full Objective evidence. This run reviewed all four owned Tasks (T-013-T-016),
including their owner-waived optional Task Checks, cross-Task integration, and
reconciliation against Design and Guardrails.

## Frozen Scope Lock

1. Acceptance and policy: O-016's eight Success Conditions; every Done-When
   criterion in T-013-T-016; FS-01, TPL-01..04, ARCH-01, CFG-01..02,
   TEST-01..05, TEST-07..09, and STYLE-03, STYLE-05, STYLE-07 as advisory.
2. Changed behavior and public surfaces: `make test-focused`, `test-fast`,
   `test-full`, and `ci`; the buildtool test runner and timing formatter; the
   migration package's prepared-result and parallel-test paths; the Linux and
   Windows CI jobs; and the active/scaffold guidance named by T-016.
3. Relied-on orchestration: Make invokes the Go buildtool, which starts
   `go test`, consumes JSON events, reports timing, preserves diagnostics and
   exit status, then full/CI gates run cross-build and packaging dependencies.
   Migration TestMain prepares and protects fixture results before tests,
   while mutating and global-hook scenarios retain isolated/serial ownership.
4. Matrix axes: focused pattern present/missing; fast/full/CI selection;
   included/excluded/failing tests; Linux and Windows hosts; normal and race
   execution; shared read-only versus private mutable fixtures; serial global
   hooks versus isolated parallel tests; command start/output/failure/exit;
   active versus scaffold guidance; and all four Task lifecycle/evidence rows.
5. Materiality boundary: a finding must violate O-016/T-013-T-016 or a named
   Guardrail through a supported gate, test, CI, or shipped-guidance path.
   Unrelated packages are admitted only where T-016's new complete Windows job
   explicitly promises to run them.

## Coverage Matrix

| Area | Expected invariant | Evidence | Classification |
| --- | --- | --- | --- |
| Task evidence | T-013-T-016 are owner-done; every skipped Task Check has an explicit owner waiver | All four records are `status: done` with Task-specific owner/time/reason waivers | Proven |
| Focused gate | Requires a pattern, remains uncached, selects requested packages, reports timing | Real focused run passed; missing `TEST` failed with the documented usage and exit 1 | Proven |
| Fast selection | Excludes exactly T-013's three slow migration parents and keeps all other packages | `make -n test-fast`; fresh `make test-fast` pass with migration at 2.318s and attributed timing | Proven |
| Linux full selection | Runs every host test uncached and all six cross-build targets | Fresh `make test-full` pass; migration 1m54.096s; Linux, Darwin, and Windows amd64/arm64 builds passed | Proven |
| Timing visibility | Slow packages/tests are attributable; budgets do not create flaky timeouts | Focused/fast/full summaries name the slowest entries; no elapsed-time assertion exists | Proven |
| Shared fixture reuse | Read-only consumers share prepared results; mutating consumers retain private state and cleanup | Focused read-only guard passed; source review and complete Linux suite cover private/recovery/no-op paths | Proven on Linux; Windows blocked by I-029 |
| Parallel groups | Isolated tests may overlap; hook/probe owners stay serial and deterministic | Source inventory, fresh full run, and `go test -race -count=1 ./internal/migrate -run TestApply_` pass (43.529s) | Proven on Linux; Windows blocked by I-029 |
| Failure propagation | Child failure and diagnostics remain visible and nonzero | Buildtool unit suite passed; missing focused pattern propagated exit 1; Windows full command propagated package failures | Proven |
| Windows CI | The configured complete suite compiles and passes on Windows | Windows Go 1.26.2/amd64 native run reproduced migration setup, watcher, and guide-casing failures | Issues - I-029, I-030 |
| Guidance and parity | Make, CI, Guardrails, root/scaffold guides, and active/scaffold skills state one contract | Direct reconciliation; SHA-256 pairs match for task/check/design skills | Proven |
| Design integration | Fast ordinary feedback is separated from complete boundary verification without dropping migration or platform promises | Linux gates and guidance agree; Windows promise is present but cannot yet execute successfully | Issues - I-029, I-030 |

Every mandatory cell was classified. There is no network/server boundary. The
finite subprocess boundary covered configured versus actual targets, start,
success, non-success, diagnostic output, exit propagation, and cleanup;
redirect, retry, and cancellation are not behaviors promised by these gates.

## Independent Adversarial Pass

- Ran the public focused target with and without its required pattern.
- Compared Make dry runs for fast, full, and CI selection; full has no skip and
  CI reaches full before build/distribution/package steps.
- Ran the read-only prepared-result guard separately and a race-enabled Apply
  matrix that includes isolated parallel tests and the serial recovery owner.
- Ran the complete Windows command from an NTFS copy because Windows Go cannot
  lock `go.mod` through the WSL UNC share. This proved the CI configuration's
  actual host behavior rather than relying on YAML or cross-compilation.
- Re-ran each Windows failure at package/test granularity. The migration setup
  failure is caused by two package `TestMain` definitions; the watcher and
  casing failures reproduce independently.

## Issues

- I-029 - O-016 added an unconditional migration `TestMain` while a Windows-only
  package `TestMain` already exists, so `internal/migrate` does not compile on
  Windows.
- I-030 - the new complete Windows suite also reproduces the board watcher
  exclusion and init guide-casing failures, so T-016's Windows job cannot pass
  after the setup conflict alone is removed.

## Commands

- `git diff --check` - pass.
- `go test -count=1 ./internal/buildtool` - pass.
- `make test-focused TEST=TestFocusedTestArgsRequiresPattern PKGS=./internal/buildtool` - pass with timing summary.
- `make test-focused` - expected refusal; exit 1 with the required-pattern usage.
- `go test -count=1 ./internal/migrate -run TestEndToEnd_sharedMigratedResultsAreReadOnly` - pass.
- `go test -race -count=1 ./internal/migrate -run TestApply_` - pass in 43.529s.
- `make test-fast` - pass; slowest package 4.488s, migration 2.318s.
- `make test-full` - pass; migration 1m54.096s and all six cross-builds completed.
- Windows `go run ./internal/buildtool test -json -count=1 ./...` - fail, with all failures propagated and timed.
- Windows focused package/test reproductions - fail as recorded in I-029 and I-030.

The Windows probes used Go 1.26.2 windows/amd64 from a local NTFS copy of the
same working tree. Linux probes used Go 1.26.2 linux/amd64 under WSL2.

## Materiality

| Issue | Likelihood | Impact | Materiality | Recommendation |
| --- | --- | --- | --- | --- |
| I-029 | High - every Windows migration package build sees both functions | High - package setup fails before any migration test runs | High | Fix now by composing one platform-correct package setup, then run Windows migration and full gates |
| I-030 | High - both focused tests reproduce on the supported Windows runner path | High - the new required Windows CI job remains red | High | Fix now in new O-016 remediation work, then run the complete Windows command |

## Non-blocking Observations

- The fresh Linux full suite remains well above the 45-second measured target:
  `internal/migrate` alone took 1m54.096s. O-016 defines timing budgets as
  reported goals rather than flaky pass/fail assertions, so this is visible
  performance debt, not a third Issue in this Check.
- The first Windows attempt from the WSL UNC path failed before tests because
  Windows Go could not lock `go.mod`; the NTFS-copy run replaced that harness
  attempt and is the evidence used by this Check.

## Owner Action

Keep T-013-T-016 `done`. Plan new remediation Task(s) under O-016 linked to I-029
and I-030, then request a new Full Objective Check. This checker did not repair
the implementation, retreat any completed Task, or mark O-016 complete.
