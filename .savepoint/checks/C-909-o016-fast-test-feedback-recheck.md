---
id: C-909
scope: {kind: objective, id: O-016}
result: CLEAR
checked_by: {role: checker, session: o016-full-recheck-20260923}
executed_session: o016-t017-focused-20260923
checked_at: '2026-09-23T04:20:00Z'
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
    - internal/migrate/replace_windows_test.go
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
    - .savepoint/objectives/O-016-fast-test-feedback/tasks/T-017-unify-migration-test-setup.md
    - .savepoint/issues/I-029-windows-migration-tests-have-two-testmain-functions.md
    - .savepoint/issues/I-030-windows-full-suite-has-platform-failures.md
  dependencies: []
issues: []
supersedes: C-908
---

# C-909: O-016 Full Objective Re-check

## Verdict

`CLEAR`, on the owner-narrowed O-016 scope. A fresh `make test-full` passed on
Linux with native Linux Node/npm. The I-029 repair is proven: the migration
package now has one `TestMain` on Windows, and the focused Windows check
passes. I independently reproduced that check with Windows Go and ran a
mutation probe that shows it catches a reintroduced duplicate. I-029 is closed
as `verified`. I-030 stays open and deferred by the owner. It is outside the
narrowed scope and does not block this verdict.

This result does not close a Task or O-016. The owner must still complete the
actions under **Owner Action**, including T-017's lifecycle and evidence, before
the Objective can close.

## Independence And Scope Basis

This session did not build any O-016 work. It is independent of the T-013-T-016
executor session (`o016-task-series-20260923`) and the T-017 executor session
(`o016-t017-focused-20260923`). The router still says `state: task` for T-017.
The owner asked for this Full Objective Check directly.

Owner scope decisions applied as recorded in O-016 `## Design Confirmation`:
- 2026-09-23 03:05 UTC: Windows coverage is the focused migration setup check
  plus the existing cross-builds. The full Windows runtime suite is out of
  scope. The Linux `make test-full` gate is still required.
- 2026-09-23 03:55 UTC: T-018 and T-019 were removed. I confirmed no T-018 or
  T-019 file exists under O-016.

## Re-check Admission Ledger

C-908 was the initial run, so its scope lock is frozen. The owner-narrowed
Windows promise replaces C-908's "complete Windows suite" cell. I did not
reinstate that full-suite requirement here, and I added no new axis.

| Re-check item | Prior Issue / claim | Frozen C-908 cell | Allowed result |
| --- | --- | --- | --- |
| One effective Windows `TestMain`; helpers dispatch before fixture prep | I-029 | Shared fixture reuse (Windows); Windows CI | Blocking |
| Focused Windows CI command compiles and passes | I-029, T-017 | Windows CI (narrowed to focused check) | Blocking |
| Fresh Linux full gate on current inputs | TEST-08; T-017 criterion 5 | Linux full selection | Blocking |
| Fast/focused selection unchanged | C-908 Proven | Fast selection; Focused gate | Blocking on regression |
| Parallel/serial ownership unchanged | C-908 Proven | Parallel groups | Blocking on regression |
| Guidance and scaffold parity | C-908 Proven | Guidance and parity | Blocking on regression |
| Board watcher / guide-casing Windows failures | I-030 | Windows CI (full suite), now owner-deferred | Observation only |

## Closure Map Of Prior Issues

| Issue | Status after C-909 |
| --- | --- |
| I-029 | Closed, `verified` by C-909 |
| I-030 | Still open. Owner-deferred on 2026-09-23 03:58 UTC and outside the narrowed O-016 scope. Not re-examined as a blocker and not accepted on the owner's behalf |

## Coverage Matrix

| Area | Expected invariant | Evidence (this run) | Classification |
| --- | --- | --- | --- |
| T-013 evidence | Baseline, inventory, isolation decision, and gate design recorded; owner acceptance plus waiver | Record unchanged since C-908 (mtime 2026-09-22T22:44Z); `status: done`, waiver with Task/reason/actor/time | Proven |
| T-014 shared reuse | Read-only consumers share protected results; mutating consumers use fresh copies | `TestEndToEnd_sharedMigratedResultsAreReadOnly` passed in Linux full and Windows focused runs; Windows consumers receive private copies | Proven |
| T-015 parallel groups | 53 isolated tests parallel; hook and probe owners serial | 53 `t.Parallel` calls (apply 13, command 1, cutover 5, e2e 19, operation 15). A scan for hook, probe, env, and chdir mutators found only 4 serial tests. The race run passed | Proven |
| Focused gate | Requires a pattern, uncached, reports timing | Real run passed with timing summary; missing `TEST` exits nonzero (make exit 2) | Proven |
| Fast selection | Excludes exactly T-013's three parents | `make -n test-fast` shows the anchored three-name skip; fresh `make test-fast` passed in 8.27 s wall | Proven |
| Linux full selection | Every host test uncached plus six cross-builds | Fresh `make test-full` exit 0. `internal/migrate` 1m48.107s, then build-linux, build-darwin, and build-windows ran | Proven |
| npm integration inside full | Runs against real native npm, not skipped | `npm` resolves to `~/.nvm/versions/node/v24.16.0/bin/npm`; `TestIntegration_InstallDependencies` PASS (not SKIP) | Proven |
| Timing visibility | Slowest packages and tests attributed; no elapsed-time assertions | All gate summaries name the slowest entries | Proven |
| Failure propagation | Child failure stays nonzero with diagnostics | Buildtool tests pass in the full suite. The Windows mutation probe shows `[setup failed]` and nonzero exit through the CI command | Proven |
| One Windows `TestMain` (I-029) | Exactly one effective `TestMain`; helpers exit in `init` before fixture prep | Source: only `end_to_end_test.go:64` defines `TestMain`; `replace_windows_test.go:68` is `init()`. `GOOS=windows go test -c` and `go vet` pass | Proven |
| Windows focused CI command (T-017) | The three named tests compile and pass on Windows | Windows Go 1.26.2/amd64 from an NTFS copy: all three tests and their subtests PASS; `internal/migrate` 1.890 s (1.982 s via the buildtool wrapper) | Proven |
| Regression guard (TEST-05) | CI job fails if the duplicate returns | Added a second `TestMain` to the NTFS copy only. The CI command reported `FAIL internal/migrate [setup failed]`, and Windows `go vet` exited 1 | Proven |
| CI routing | Linux runs `make ci`; Windows runs only the focused check | `.github/workflows/ci.yml`: `ci` job runs `make ci`. The `windows-tests` job runs the focused command matching T-017's Technical Verification | Proven |
| Guidance and parity | Make, CI, Guardrails TEST-08, root and scaffold guides, and skills agree; skill copies byte-identical | `cmp` identical for task, check, and design skills and the three references. TEST-08 and AGENTS text match the Makefile | Proven |
| Design reconciliation | Design §13 build and cross-build statements still hold | `make build-all` targets and `make ci` unchanged in meaning. Design makes no Windows-runtime-suite promise that the narrowing contradicts | Proven |
| T-017 lifecycle and evidence | Task handoff records a fresh full gate | T-017 is `in_progress`/`stage: test`. Its evidence still says there is no passing full gate, and it has no Task-check waiver | Owner/executor action (see below). The technical outcome is proven by this run |

No cell is unverified. There is no network or server boundary. The subprocess
boundary covers the buildtool wrapper, Windows helper re-exec, and npm. For
each, this run checked the configured and actual targets, success,
non-success, diagnostics, and exit propagation. Redirect, retry, and
cancellation are not promised by these gates.

## Commands And Results

Linux: Go 1.26.2 linux/amd64, WSL2 kernel 6.18.33.2-microsoft-standard-WSL2,
native Node v24.16.0 and npm 11.13.0 from nvm (first on `PATH`). No npm bridge
or `TMPDIR` override.

- `make test-full`: **pass**, exit 0. Started 2026-09-23T04:14:15Z and
  finished 04:16:21Z (2:06.39 wall). Slowest tests:
  `TestEndToEnd_temporaryRepositoryCopyMigratesWithReleaseAccountability`
  1m31.74s, `TestApply_recoversAtEveryPublishBoundaryWithoutOverwritingUserEdits`
  15.1s, and `TestEndToEnd_goldenIsReproducible` 1.72s. All three cross-build
  phases completed.
- Input fingerprint at run time and after all probes (unchanged): HEAD
  `e7e72bb3d736028d208bc070f012fd6c66e8c930`. `git diff | sha256sum` =
  `a1e3a570df1fcc8e25ec1ebdc5490a68afdf1a73faf7b17ee4c8c6aa355bb3a1`, and
  `git status --porcelain | sha256sum` =
  `e5b354d7236278c29a70fd457230bb702fd2394a6994d3cdb53d2bf29a53ffe5`. This
  record and the I-029/O-016 metadata edits were made after the run and are
  metadata-only.
- `git diff --check`: pass. `make build`: pass.
- `make test-fast`: pass, 8.27 s wall.
- `make test-focused TEST=TestFocusedTestArgsRequiresPattern PKGS=./internal/buildtool`: pass. `make test-focused` with no pattern: refused, nonzero exit.
- `make -n test-fast`, `make -n test-full`, `make -n ci`: selections as recorded in the matrix.
- `go test -race -count=1 ./internal/migrate -run 'TestApply_|TestOperation_|TestPreflightCutover_'`: pass, 43.605 s.
- `go test -count=1 -v -run TestIntegration_InstallDependencies ./internal/init`: PASS.
- `GOOS=windows GOARCH=amd64 go test -c ./internal/migrate` and `go vet`: pass.

Windows: Go 1.26.2 windows/amd64, run from a fresh NTFS copy of the same
working tree (excluding `.git`). The copy was removed afterwards.

- `go run ./internal/buildtool test -json -count=1 ./internal/migrate -run '^(TestEndToEnd_sharedMigratedResultsAreReadOnly|TestProbeDestinationHeldOpen|TestProbeInterruptionBeforeReplace)$'`: pass, package 1.982s.
- `go test -count=1 -v` with the same selection: `TestProbeDestinationHeldOpen`, `TestProbeInterruptionBeforeReplace`, and `TestEndToEnd_sharedMigratedResultsAreReadOnly` PASS; `ok internal/migrate 1.890s`.
- Mutation probe with a duplicate `TestMain` added to the copy only: the CI command failed with `[setup failed]`, and `go vet` exited 1.

## Applicable Guardrails

- CFG-02: the Windows difference (private copies, helper dispatch in `init`)
  is explicit and tested on Windows by the focused job. Satisfied.
- TEST-05: the I-029 regression is guarded by the focused Windows CI job. The
  mutation probe proves it fails on regression. Satisfied.
- TEST-08: current successful `make test-full` evidence exists for this
  Check. Satisfied for the Objective Check. T-017's own handoff record still
  needs to cite it (Owner Action).
- TEST-09: T-013-T-016 carry waivers. T-017 has none yet (Owner Action).
- TPL-01: byte parity holds. FS-01, TEST-01..04, TEST-07, CFG-01, and ARCH-01
  show no regression.
- STYLE rules: advisory only; nothing affects the verdict.

## Materiality

No Issues were admitted in this run, so no materiality actions are required.

## Non-blocking Observations

- The full gate still takes about 2 minutes, against the 45-second target.
  `internal/migrate` is 1m48s, and one repository-copy test is 1m31s. O-016
  treats timing as reported budgets, not assertions. C-908 recorded the same
  performance debt.
- I-030 (board watcher exclusion and init guide casing on Windows) is still
  unrepaired and owner-deferred. T-017 also recorded unrelated Windows failures
  in the migration package's full run: inventory case collision, a read-only
  install, and cutover directory mtimes. No Issue covers them yet. Neither
  blocks O-016 under the narrowed scope. If the full Windows suite returns to
  scope, capture the migration-package failures as a new Issue or extend I-030.
- The shipped V2 scaffold `AGENTS.md` and skills now name repository-specific
  targets (`make test-fast`, `make test-full`, "in this repository ... `make
  ci`"). New user projects may not define these targets. T-016 required
  scaffold parity by design, and C-908 accepted the guidance cell. This is a
  possible TPL-02 follow-up for the planner, not a finding here.
- I-031 (reload diagnostics on Issue records) is unrelated to O-016's
  acceptance. It was not assessed here.

## Owner Action

Only the owner can take these steps. This Check did not take them.

1. T-017: it is `in_progress`/`stage: test`, and its evidence still says there
   is no passing full gate. Before owner closure, its evidence must cite a
   passing fresh `make test-full`. That can be this run: command, times,
   toolchain, and input fingerprint are above; citing it is a metadata-only
   reuse. It also needs either an explicit owner Task-check waiver (Task,
   reason, actor, time) or a requested Task Check. Then the owner may set
   T-017 `done`.
2. T-013-T-016 stay `done`. No Task was retreated.
3. After every owned Task is `done`, the owner may accept and close O-016 on
   this current Check. O-016 does not declare `owner_validation.required`.
   I-029 is resolved. C-909 links no open Issues.
4. I-030 stays open and deferred. It is not accepted on the owner's behalf.
