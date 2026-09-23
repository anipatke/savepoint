---
id: T016
title: Put fast and full gates in place
objective: O016
status: done
depends_on: [{task: T013, requires: clear}, {task: T014, requires: clear}, {task: T015, requires: clear}]
owner_validation:
    required: true
    accepted_check: ""
planned_by: {role: planner, session: o016-design-20260923}
complexity_tier: high
complexity_reason: Build commands, CI, guardrails, and shipped workflow instructions must change as one verification contract.
check_waiver:
    task: T016
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-09-23T00:55:28Z"
---

# T016: Put fast and full gates in place

## Outcome

Developers and agents have named focused, fast, and full verification paths, and CI reports their timing while running the complete suite.

## User Check

Review the final gate table and evidence-reuse wording. Confirm an ordinary Task, a migration or platform-sensitive Task, an Objective Check, and a metadata-only correction each select the intended evidence.

## Done When

- Make targets implement T013's exact selection. The full gate runs every host-platform Go test, including migration recovery, interruption, and integration, plus cross-build checks. CI also executes Windows-only tests on Windows.
- Ordinary Task handoff requires build plus fast; migration or platform-sensitive Task handoff requires fresh full; CI and Full Objective/Release Checks require full.
- Focused tests remain an iteration aid. Metadata-only full-result reuse requires a record of the original run and proof that code, tests, fixtures, and gate definitions did not change.
- CI runs full and publishes package/test timing with enough attribution to identify a regression.
- Guardrails, active task/check/design skills, their V2 scaffold copies, AGENTS.md, and Make/CI guidance use the same names and responsibilities; canonical and scaffold skills remain byte-identical.
- Recorded warm-machine measurements assess the 2/15/45-second goals without making elapsed time a flaky pass/fail assertion.
- Selection tests or command transcripts prove both fast omission and full inclusion, including a failing-test path; `make build`, fast, and full pass.

## Context Files

`Makefile`, `.github/workflows/ci.yml`, `AGENTS.md`, `.savepoint/Guardrails.md`, `agent-skills/savepoint-task/SKILL.md`, `agent-skills/savepoint-check/SKILL.md`, `agent-skills/savepoint-design/SKILL.md`, `templates/project-v2/agent-skills/savepoint-task/SKILL.md`, `templates/project-v2/agent-skills/savepoint-check/SKILL.md`, `templates/project-v2/agent-skills/savepoint-design/SKILL.md`, `.savepoint/objectives/O016-fast-test-feedback/tasks/T013-measure-and-design-test-gates.md`.

## Design References

Design sections 1, 12, and 13; O016 Confirmed Verification Policy.

## Guardrails

TPL-01..04, TEST-01..04, TEST-07..09, CFG-01..02, ARCH-01, STYLE-07.

## Implementation Plan

1. Implement named fast and full targets from T013's selection, retaining `make test` as the full suite unless the confirmed contract explicitly revises it.
2. Add selection verification and CI timing output; ensure CI runs full and a Windows test job.
3. Reconcile TEST-08 and all active/scaffolded workflow guidance in the same integrated change.
4. Document evidence freshness and metadata-only reuse with explicit input comparison.
5. Exercise ordinary, sensitive, CI, Check, and reuse scenarios; record timings and run the revised complete gate.

## Boundaries

No test removal, test-cache substitution for changed inputs, hidden CI exclusion, or automatic timeout failure from timing goals.

## Technical Verification

Gate selection and failure-path evidence, active/scaffold parity, full CI-equivalent command including cross-build checks and Windows test evidence, warm timing samples, `make build`.

## Technical Evidence
### Start and scope

- Started 2026-09-23 at `stage: build` after the owner's instruction to begin T016. T013, T014, and T015 were `done` with explicit owner Task-check waivers, satisfying T016's dependencies requiring `clear`.
- Read the router, O016, T016, and dependency records T013-T015. T014 and T015 were extra dependency reads used to confirm their owner closure and waiver state.
- Read `agent-skills/savepoint-task/SKILL.md` for the required workflow. Read Design sections 1, 12, and 13 and the named Guardrail IDs for the confirmed selection and consistency rules.
- Extra implementation reads: `internal/buildtool/main.go` to reuse the existing Makefile helper, and `internal/buildtool/main_test.go` to add failure-propagation coverage. Read `templates/project-v2/AGENTS.md` to align shipped workflow guidance. The scaffold Guardrails file at `templates/project-v2/.savepoint/Guardrails.md` is a generic worked example without TEST-08, so it remains unchanged; the live `.savepoint/Guardrails.md` TEST-08 carries this repository's gate policy. The initially checked `templates/project-v2/Guardrails.md` path does not exist.
- The router selected T016 but had a stale T015 owner-review action. Updated it to T016's implementation action, then updated it at handoff to the audit-stage owner review.

### Gate contract delivered

| Use | Gate |
|---|---|
| Test iteration | `make test-focused TEST=<pattern> [PKGS=./package]` |
| Ordinary Task handoff | `make build && make test-fast` |
| Migration/platform-sensitive Task handoff | Fresh `make test-full` |
| CI | `make ci`, which depends on `make test-full`, native build, distribution, and package checks |
| Full Objective or Release Check | Current successful `make test-full` evidence |
| Metadata-only correction | Reuse a recorded full result only after proving code, tests, fixtures, dependencies, and gate definitions are unchanged |

- `make test` remains the complete host Go suite and now runs uncached with JSON events. `make test-fast` applies exactly T013's three exclusions: `TestEndToEnd_temporaryRepositoryCopyMigratesWithReleaseAccountability`, `TestApply_recoversAtEveryPublishBoundaryWithoutOverwritingUserEdits`, and `TestEndToEnd_goldenIsReproducible`. `make test-full` runs the uncached full suite without exclusions and all six Linux, macOS, and Windows cross-build targets.
- The shared build helper prints the ten slowest packages and tests with durations. It retains test output for failed packages, prints diagnostics, and returns the Go test process's nonzero exit. Timings are reported only; no elapsed-time timeout was added.
- CI retains the Ubuntu `make ci` job and adds `windows-tests`, which installs Go 1.26 and Node 24, then runs the complete uncached Go suite on `windows-latest` through the timing helper.
- Ordinary, migration/platform-sensitive, CI, Check, focused-iteration, and metadata-only reuse responsibilities are aligned in the Makefile, CI workflow, live Guardrails, root and scaffold `AGENTS.md`, and active task/check/design skills. Each V2 scaffold skill copy is byte-identical to its active source (SHA-256 pairs match for all three skills).

### Selection and failure-path evidence

- `make -n test-fast` printed the anchored T013 `-skip` expression with exactly the three named tests. Make escaping preserves the closing `$` in the shell command.
- `make -n test-full` printed `go run ./internal/buildtool test -json -count=1 ./...` with no `-skip`, followed by `build-linux`, `build-darwin`, and `build-windows`; `make -n ci` showed `test-full`, `build`, `dist`, and `package-check`.
- `go test -count=1 ./internal/buildtool` passed. `TestFocusedTestArgsRequiresPattern` rejects a missing/blank focus; `TestFocusedTestArgsAddsUncachedTimingAndPackage` checks focused selection; `TestRunGoTestCommandReportsTimingAndPropagatesFailure` injects a child that emits passing/failing JSON and exits 7, then verifies the error retains exit code 7 and the timing/diagnostic output.
- A first fast run used a restricted WSL `PATH`; the repository's temporary npm bridge then could not find Windows `cmd.exe`, causing `TestIntegration_InstallDependencies` to fail. Rerunning with the inherited WSL `PATH` and `/tmp/t014-tools` prefixed passed. The wrapper was also adjusted to suppress passing-test chatter while retaining failing-package output before the successful gate runs below.

### Gate runs and warm measurements

Environment: same warm WSL2 Ubuntu 24.04.4 / Go 1.26.2 linux/amd64 machine recorded in T013 (Ryzen 7 7800X3D, 16 logical CPUs). Go test result caching was disabled with `-count=1`; the build cache was warm.

| Command | Result | Wall time | Timing evidence |
|---|---:|---:|---|
| `make test-focused TEST=TestVersion_override PKGS=./internal/buildtool` | Pass | 0.61 s | Within 2 s goal |
| `make test-fast` | Pass | 9.30 s | Within 15 s goal; `internal/migrate` 2.434 s |
| `make test-full` warm run 1 | Pass | 127.47 s | Above 45 s goal; all tests and six cross-builds completed |
| `make test-full` warm run 2, after final CI setup | Pass | 141.47 s | Above 45 s goal; all tests and six cross-builds completed |
| `make build` | Pass | — | Native build completed |

Across the two full samples, `internal/migrate` took 1m49.145s and 2m4.17s. In the final sample, `TestEndToEnd_temporaryRepositoryCopyMigratesWithReleaseAccountability` took 1m47.69s, `TestApply_recoversAtEveryPublishBoundaryWithoutOverwritingUserEdits` took 15.3s, and `TestEndToEnd_goldenIsReproducible` took 1.56s. The 45-second target is missed in both runs and remains a budget, not a flaky gate assertion.

### Acceptance evidence and limitations

1. **Gate selection:** Pass. Fast uses only T013's exact three-test skip expression; full uses `./...` without a skip and performs all six cross-builds.
2. **Handoff responsibilities:** Pass. Ordinary and migration/platform-sensitive Task commands, CI, and mandatory Objective/Release Check requirements use the documented gates.
3. **Evidence reuse:** Pass. Policies require the original command, time, toolchain, and result plus proof the code, tests, fixtures, dependencies, and gate definitions did not change; otherwise a fresh full run is required.
4. **Timing and CI attribution:** Pass. Package/test timing is in command output for Linux CI and the new Windows full-suite job. The 2 s and 15 s goals passed; the 45 s goal did not.
5. **Guidance and scaffold parity:** Pass. The live TEST-08, root/scaffold guides, Make/CI instructions, and active/scaffold task/check/design skills were reconciled; skill SHA-256 pairs match.
6. **Selection and failure path:** Pass. Dry-run transcripts show fast exclusions and full inclusion; the helper's subprocess test proves nonzero propagation and timing output.
7. **Required gates:** Pass locally: `make build`, `make test-fast`, and `make test-full` all exited 0.

The Windows job is configured with Go 1.26 and Node 24. `GOOS=windows GOARCH=amd64 go test -c -o /tmp/savepoint-buildtool.test.exe ./internal/buildtool` passed, but the hosted Windows test job itself cannot be executed from this WSL session and its first CI result remains outstanding. `make ci` wiring was reviewed with `make -n ci`; its test-full, build, and cross-build components were executed separately, while distribution/package steps were not rerun locally. Subsequent changes were limited to this Task's evidence and the router's next-action metadata, so they did not alter the recorded code, test, fixture, dependency, or gate inputs.