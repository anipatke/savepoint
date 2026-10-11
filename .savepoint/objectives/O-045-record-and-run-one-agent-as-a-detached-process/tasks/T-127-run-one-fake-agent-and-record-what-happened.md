---
id: T-127
title: Run one fake agent and record what happened
objective: O-045
status: done
depends_on: []
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: planning-g003-20261011}
check_waiver:
    task: T-127
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-10T20:38:06Z"
---

# Run one fake agent and record what happened

## Outcome

A new `internal/run` package can start one configured command in the project root, stream its stdout line by line into an append-only `events.jsonl`, keep `run.yml` current with a heartbeat while it runs, and record the final state (`finished` or `failed`) with exit code and a sanitised stderr tail. No board, no CLI wiring, no real provider yet: a fake provider script is enough to prove the record.

## User Check

Not required. Technical Task: the Full Objective Check verifies the record shape against O-045's Confirmed Design.

## Done When

1. `run.Start(spec)` creates `.savepoint/runs/<yyyymmddThhmmssZ>-<task>/` and writes `run.yml` with `provider`, `command`, `cwd`, `task`, `next_line`, `started_at`, `updated_at`, `pid`, `state: running` before the process starts (O-045 condition 1).
2. Every stdout line of the child is appended verbatim to `events.jsonl` with `O_APPEND`; nothing is parsed or dropped in this Task (condition 2's "kept verbatim" half).
3. `updated_at` is rewritten atomically at least every 5 s while the child runs; `run.yml` is never written non-atomically (FS-04).
4. On exit, `state` is `finished` with `exit_code` for exit 0, and `failed` with `exit_code` and a `stderr_tail` for non-zero exit; the tail is one printable line with URL credentials removed, reusing the health runner's sanitiser or a copy of it (condition 4, ARCH-05).
5. The child runs with no shell, `CI=1 NO_COLOR=1`, its own process group on Unix, and the Windows tree-kill hook, generalised from `internal/codehealth/runner_unix.go` and `runner_windows.go` without the board or `internal/run` importing `codehealth` (ARCH-04, DEP-02).
6. Tests use a fake provider (a Go test helper binary or script) and temporary directories only (TEST-04): running, finished, failed, and the heartbeat case. `make build && make test-fast` pass (TEST-08). Per-criterion evidence follows AGENTS.md's Verification Policy, including an explicit owner waiver if the optional Task Check is skipped.

## Context Files

`internal/codehealth/runner.go`, `internal/codehealth/runner_unix.go`, `internal/codehealth/runner_windows.go`, `internal/data/write.go`, `internal/data/launcher_config.go`, `.savepoint/objectives/O-045-record-and-run-one-agent-as-a-detached-process/Objective.md`, `docs/reviews/agent-orchestration-pass-1.md` (section 5 only).

## Design References

Design §1 (file-first, token-efficiency), §6 (CLI surface, agent-runnable commands), §11 (failure modes). O-045 Confirmed Design.

## Guardrails

FS-04, FS-05, FS-06, ARCH-03, ARCH-04, ARCH-05, DEP-02, CFG-02, TEST-01, TEST-02, TEST-04, TEST-08.

## Implementation Plan

1. Confirm `codehealth.ExecRunner`'s process-group and tree-kill code still matches O-045's Confirmed Design; return REPLAN REQUIRED if the runner cannot be generalised without changing health behaviour.
2. Add `internal/run`: `Spec`, `Record` (the `run.yml` fields), an atomic record writer, and `Start`.
3. Move the platform process setup into a small shared package or duplicate it under `internal/run` with a note; keep `codehealth` behaviour byte-identical.
4. Stream stdout to `events.jsonl`; keep a bounded stderr tail.
5. Heartbeat goroutine; final state on exit.
6. Tests with the fake provider on Unix and Windows build tags.

## Boundaries

No parsing of provider output, no cancel or timeout (next Tasks), no CLI subcommand, no board change, no Task or router write, no new third-party dependency.

## Technical Verification

Focused `make test-focused TEST=...` during iteration; `make build && make test-fast` for handoff. Process and path work is platform-sensitive: the Full Objective Check needs the native Windows CI job's result, produced by pushing the branch, before it starts.

## Technical Evidence

Executor session: task-t127-20261011. Claims below are for a fresh Check to verify.

**Per-criterion outcomes**
1. `run.Start` creates `.savepoint/runs/<yyyymmddThhmmssZ>-T-###/` and writes `run.yml` (state `running`, all listed fields) before launching; `pid` is added by a second atomic write right after the process starts, since a PID cannot exist earlier. Tests: `TestStart_finishedRunKeepsStdoutVerbatim`, `TestStart_recordsRunningBeforeAndWhileChildRuns`.
2. Each stdout line is appended whole with `O_APPEND`, nothing parsed or dropped, including non-JSON lines and a final line without a newline (newline added). Test: `TestStart_finishedRunKeepsStdoutVerbatim`.
3. Heartbeat rewrites `updated_at` every 5 s by default (`DefaultHeartbeat`); `run.yml` is only written through a temp file plus rename. Test: `TestStart_heartbeatRewritesUpdatedAtWhileRunning` (20 ms interval, injected clock, run.yml parsed on every poll).
4. Exit 0 gives `finished` with `exit_code`; non-zero gives `failed` with `exit_code` and a one-line `stderr_tail` (control characters and URL credentials removed, end kept). Test: `TestStart_failedRunRecordsExitCodeAndSanitisedStderrTail`; `TestSanitizeTail_keepsEndOfOversizedText`. The sanitiser is a copy of the health runner's (`internal/run/sanitize.go`).
5. No shell (`exec.CommandContext` with argv); env adds `CI=1 NO_COLOR=1` (the fake agent exits 9 if absent); process group plus group kill on Unix and `taskkill` tree kill on Windows are copied into `proc_unix.go` / `proc_windows.go`. `internal/run` does not import `codehealth`, and `internal/codehealth` is unchanged.
6. Tests use a fake provider (the test binary re-run via `TestMain`) and temp dirs only. Also covered: refusal outside a Savepoint project, empty command, path-escaping task ID (nothing written), and a missing executable (recorded `failed`). `go vet` passes for linux, windows and darwin.

**Commands run**
- `go test -count=1 -race ./internal/run`: ok
- `GOOS=windows go vet ./internal/run`, `GOOS=darwin go vet ./internal/run`: ok
- `make build && make test-fast`: both exit 0

**Files read:** the Context Files listed above, except the `docs/reviews/...` section 5 read as `sed -n 185,247p`.
**Extra reads (logged):** `.savepoint/router.md` and `internal/data/write.go` lines 250-320 (how the project writes files atomically); `go.mod`, `Makefile`, `.savepoint/Guardrails.md` guardrail rows (to confirm gate commands and rule text); a search for `gopkg.in/yaml.v3` use (already imported by `internal/data`; no new dependency).
**Files changed:** new `internal/run/{run,record,sanitize,proc_unix,proc_windows,run_test}.go`; `AGENTS.md` (one Codebase Map row for `internal/run/`, required by ARCH-04); this Task and O-045 status.

**Limitations**
- Windows and macOS code was only compiled and vetted here, not run. Native Windows CI evidence is still needed before the Full Objective Check.
- `make test-full` was not run (not required for ordinary handoff; the Objective Check needs it).
- No cancel or timeout yet (next Tasks); the process context is `context.Background()`.
- `go.mod` still lists `gopkg.in/yaml.v3` as `// indirect`, as before; left untouched.
- Optional Task Check not requested; no owner waiver recorded.

## Drift Notes

Record architecture deltas and reconcile through the planner before Check.
