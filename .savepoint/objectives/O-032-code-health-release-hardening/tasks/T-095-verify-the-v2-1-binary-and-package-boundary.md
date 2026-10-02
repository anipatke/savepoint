---
id: T-095
title: Verify the v2.1 binary and package boundary
objective: O-032
status: done
depends_on: [{task: T-087, requires: clear}, {task: T-088, requires: clear}, {task: T-089, requires: clear}, {task: T-090, requires: clear}, {task: T-091, requires: clear}, {task: T-092, requires: clear}, {task: T-093, requires: clear}]
complexity_tier: medium
complexity_reason: Release evidence must cover six archives, checksums and provider-free runtime distribution.
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: planning-o032-20261002-owner-confirmed}
check_waiver:
    task: T-095
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-02T04:59:18Z"
---

# Verify the v2.1 binary and package boundary

## Outcome

Current release artifacts prove Savepoint remains one executable per platform with correct version/checksums and no bundled or installed analysis providers.

## User Check

Read the distribution inventory and verify each archive contains one platform executable and has a checksum.

## Done When

- Run make ci and retain command/time/toolchain/result evidence, six target inventory, archive member validation, checksum verification and native version smoke outcome.
- Demonstrate packaged help/core behavior works without provider executables on PATH; unavailable providers remain honest unavailable states, and no runtime download/install is attempted.
- Inspect package/runtime dependency boundary and document provider prerequisites separately from the shipped binary.
- Record current native Windows full-test CI result when accessible; if unavailable, record the exact missing evidence for the Full Objective Check, never claim cross-compilation is native runtime validation.
- Add narrowly missing packaging failure regressions only if evidence exposes an uncovered release contract; no publishing, tagging or deployment.

## Context Files

`.savepoint/objectives/O-032-code-health-release-hardening/tasks/T-087-simplify-diagnostic-repair-rules.md` (dependency evidence); `.savepoint/objectives/O-032-code-health-release-hardening/tasks/T-088-separate-test-stream-processing-responsibilities.md` (dependency evidence); `.savepoint/objectives/O-032-code-health-release-hardening/tasks/T-089-simplify-board-reload-state-restoration.md` (dependency evidence); `.savepoint/objectives/O-032-code-health-release-hardening/tasks/T-090-protect-saved-health-evidence-under-failures.md` (dependency evidence); `.savepoint/objectives/O-032-code-health-release-hardening/tasks/T-091-verify-provider-processes-stop-safely.md` (dependency evidence); `.savepoint/objectives/O-032-code-health-release-hardening/tasks/T-092-preserve-projects-when-adopting-code-health.md` (dependency evidence); `.savepoint/objectives/O-032-code-health-release-hardening/tasks/T-093-exercise-all-health-readers-together.md` (dependency evidence); `internal/buildtool/main.go`; `internal/buildtool/main_test.go`; `Makefile`; `package.json`; `.github/workflows/ci.yml`; `main.go`; `main_health_test.go`.

## Design References

O-032 Confirmed Release-Hardening Design and Owner-Requested Health Analysis; Design sections 1, 6, 7, 9, 12 and 13 as applicable.

## Guardrails

FS-01, FS-04, FS-05, FS-06, ARCH-01, ARCH-02, ARCH-03, CFG-01, CFG-02, TEST-01, TEST-02, TEST-03, TEST-04, TEST-05, TEST-08, TEST-09, STYLE-01, STYLE-02, STYLE-07, CFG-03, REL-01, REL-02, REL-03.

## Implementation Plan

1. Inspect current archive/checksum/package tests and CI matrix.
2. Build and validate all distribution outputs using existing buildtool gates.
3. Smoke provider-free packaged runtime in a temporary PATH/project environment.
4. Record artifact and platform evidence for independent Objective review.

## Boundaries

No publication, dependency installation, new CI/CD integration or separate release runner; fix only demonstrated distribution contract failures.

## Technical Verification

Fresh `make test-full` for this platform/migration/distribution-sensitive Task; `make ci` additionally for distribution evidence; focused `make test-focused TEST=... PKGS=...` only during iteration. Record named per-criterion evidence and limitations. The later independent Check applies `agent-skills/references/check-method.md`; the mandatory Full Objective Check runs its fresh full gate and official Code Health collection. This executor never writes a Check or runs savepoint health setup/report/check.

## Technical Evidence

Executed 2026-10-02 on linux/amd64, go1.26.2, node v22.22.2, npm 10.9.7, branch v2.1 at 136fe69 plus the uncommitted T-094/T-097/T-098 working tree.

**Criterion 1 — make ci and inventory.** `make ci` (= `test-full build dist package-check`) exited 0 at 2026-10-02T04:57:26Z (~22s wall; Go test cache not relied on, `-count=1`). First attempt failed: `unexpected distribution archive "savepoint-v1.3.1-darwin-amd64.tar.gz"` — six stale v1.3.1 archives left in untracked `dist/` from a 2026-09-26 build. The guard is correct (it refuses mixed versions); I moved them to the session scratchpad (not deleted) and re-ran. No code change. Result: six archives `savepoint-v2.0.5-{linux,darwin,windows}-{amd64,arm64}.tar.gz`; `sha256sum -c dist/checksums.txt` OK for all six; `tar tzf` shows exactly one member each (`savepoint`, or `savepoint.exe` for windows); `make verify-dist` → "distribution verified: 6 archives (v2.0.5)"; `make smoke-test` → native `v2.0.5`, "smoke test passed"; `npm test` passed; `npm pack --dry-run` lists LICENSE, README, bin/savepoint.js, six platform binaries, package.json (10 files).

**Criterion 2 — provider-free runtime.** Extracted the linux-amd64 archive to a temp dir and ran with `env -i PATH=/nonexistent`: `--version`, `--help`, `init`, `doctor` (ALL CLEAN), `resume`, `health setup` preview and `--apply` all worked. Missing providers (lizard, jscpd, osv-scanner) reported as `gap: missing executable` with "install it yourself, Savepoint does not"; no download/install attempted. Source grep: no `net/http` or install invocation in non-test code; the only `npm install` is the opt-in `init --install` for the user's own project. `health check` was not run end to end here (needs an Objective in the temp project); provider unavailability under check is covered by the T-093 tests.

**Criterion 3 — dependency boundary.** Go deps are bubbletea/fsnotify and transitive libs only; npm package ships only the wrapper plus the six binaries; no provider is bundled or declared as a dependency. Provider prerequisites are stated by `savepoint health setup` output. README.md does not mention Code Health providers at all; documenting them is T-096's scope, not changed here.

**Criterion 4 — native Windows CI.** Not available. `windows-tests` in `.github/workflows/ci.yml` triggers only on push/PR to master and v2; `gh run list --branch v2.1` returns no runs. The latest Windows run seen is on master/fix branch (2026-09-27, success) and predates this branch's work. The windows archives here are cross-compiled and only structure/checksum verified, not runtime validated. Missing for the Full Objective Check: a green `windows-tests` run on the O-032 head commit (e.g. via PR to v2/master).

**Criterion 5 — regressions.** None added; evidence exposed no uncovered contract (stale-version refusal already enforced by `dist`). No publishing, tagging or deployment.

Files read: Context Files listed (buildtool/main.go partially, ci.yml, Makefile, package.json), plus extra reads: `agent-skills/savepoint-task/SKILL.md` and AGENTS.md (workflow), bin/savepoint.js, go.mod, README.md grep, `internal/init/install.go` grep (verify no runtime install). Files changed: this Task file only (router already selected T-095).

Limitations: linux-only runtime smoke; version is still 2.0.5 in package.json/binaries on branch v2.1 (owner decision 2026-10-02: the version is bumped on merge with master, so 2.0.5 here is expected); working tree was dirty with other Tasks' changes, so evidence is for that tree; no Windows native evidence.

## Drift Notes

Record implemented responsibility or interface changes for planner reconciliation before the mandatory Full Objective Check. If evidence requires a material unknown repair, return REPLAN REQUIRED rather than inventing another scope.
