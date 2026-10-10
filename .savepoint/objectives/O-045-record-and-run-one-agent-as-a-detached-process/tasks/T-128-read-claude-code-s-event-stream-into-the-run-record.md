---
id: T-128
title: Read Claude Code's event stream into the run record
objective: O-045
status: done
depends_on: [{task: T-127, requires: clear}]
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: planning-g003-20261011}
check_waiver:
    task: T-128
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-10T20:54:30Z"
---

# Read Claude Code's event stream into the run record

## Outcome

The runner understands Claude Code's `stream-json` output well enough to fill the run record honestly: which files were read and edited, which commands ran and how they exited, which actions were denied, the session ID, the provider's usage figures, and the agent's final message under `agent_reported`. Lines it does not understand stay on disk verbatim and are counted as unparsed, never dropped. The argv for a real run is built from `agent_launcher.builder` and the Next line.

## User Check

Not required. Technical Task: the Full Objective Check verifies the parsed fixture against the documented event format.

## Done When

1. A `claude` adapter maps `system/init` (session ID, model), `assistant` content `tool_use` blocks (tool name, file path or command), `user` content `tool_result` blocks (exit code when present), `system/permission_denied`, and the final `result` (`session_id`, `usage`, `total_cost_usd`, `permission_denials`, `result` text) to display events and `run.yml` fields; `stream_event` partial deltas are ignored (O-045 condition 2).
2. Unknown line types and malformed JSON lines are kept verbatim in `events.jsonl`, counted in `run.yml` as `unparsed`, and never fail the run (DATA-03 spirit).
3. `summary.md` holds the agent's final message verbatim under a heading that says it is agent-reported; nothing in `internal/run` interprets or scores it.
4. The argv builder substitutes `{{prompt}}` with the Next line exactly as `internal/resume` prints it and `{{project_root}}` with the absolute project root, from `agent_launcher.builder`; a builder with no command fails with a named error before any file is written (CFG-01). The default builder used when the owner has not set one is the Claude command in O-045's Confirmed Design.
5. Text copied into `run.yml` or `summary.md` has terminal control sequences removed and URL credentials stripped (ARCH-05).
6. Tests replay a recorded `stream-json` fixture captured from a real Claude Code 2.1.296 session (checked in under `testdata`) through the fake provider from T-127 and assert every field above, plus the unknown-line and malformed-line cases (TEST-01, TEST-02, TEST-04). `make build && make test-fast` pass (TEST-08). Per-criterion evidence follows AGENTS.md's Verification Policy, including an explicit owner waiver if the optional Task Check is skipped.

## Context Files

`internal/run/` (as landed by T-127), `internal/data/launcher_config.go`, `internal/data/config.go`, `internal/resume/resume.go`, `internal/codehealth/runner.go` (sanitiser only), `.savepoint/objectives/O-045-record-and-run-one-agent-as-a-detached-process/Objective.md`, `docs/reviews/agent-orchestration-pass-1.md` (section 2, Claude Code column, and section 5).

## Design References

Design §1 (file-first), §6 (CLI surface). O-045 Confirmed Design. Claude Code headless documentation as cited in the review (https://code.claude.com/docs/en/headless).

## Guardrails

DATA-03, ARCH-04, ARCH-05, CFG-01, STYLE-07, STYLE-09, TEST-01, TEST-02, TEST-04, TEST-06, TEST-08.

## Implementation Plan

1. Confirm T-127's record and streaming shapes; return REPLAN REQUIRED if `events.jsonl` is not the raw provider line per line.
2. Capture a short real `claude -p --output-format stream-json --verbose` session in a scratch project as the fixture; redact nothing that is not a secret.
3. Add the adapter as a pure function over lines; keep provider-specific knowledge in one file.
4. Add the argv builder over `data.AgentLauncher` with the default Claude command spec in data, not logic (STYLE-09).
5. Wire the adapter into the runner's line loop; update `run.yml` on `result`.
6. Tests for every mapped type, unknown lines, malformed lines, and the argv builder.

## Boundaries

No cancel, timeout, or liveness; no CLI subcommand; no board change; no second provider; no interpretation of agent text.

## Technical Verification

Focused `make test-focused TEST=...` during iteration; `make build && make test-fast` for handoff. The Full Objective Check needs fresh `make test-full` under `agent-skills/references/check-method.md`.

## Technical Evidence

### Per-criterion outcomes

1. **Adapter mapping — met.** `internal/run/claude.go` `ParseClaudeLine` maps `system/init` (session ID, model), `assistant` `tool_use` blocks (Read → file read; Edit/Write/MultiEdit/NotebookEdit → file edited; Bash → command; other tools → generic `tool`), `user` `tool_result` blocks, `system/permission_denied`, and `result` (`session_id`, `usage`, `total_cost_usd`, `permission_denials`, `result` text) onto `Event`s; `events.go` `observer.apply` folds them into `run.yml` fields (`session_id`, `model`, `files_read`, `files_edited`, `commands`, `denials`, `usage`, `total_cost_usd`, `agent_reported`). `stream_event` returns no events. Tests: `TestParseClaudeLine_mapsEachDocumentedLineType`, `TestReplay_realSessionWithEditAndDenialFillsTheRecord`, `TestReplay_failingCommandRecordsItsExitCode`.
2. **Unknown and malformed lines — met.** Every stdout line is written to `events.jsonl` before it is interpreted; an unknown top-level type, malformed JSON, `null`, or a non-object `message` yields one `unparsed` event that increments `run.yml` `unparsed`; the run still finishes. Tests: `TestReplay_unknownAndMalformedLinesStayOnDiskAndAreCounted` (3 unparsed, file byte-identical, state finished) and the real fixtures (1 unparsed each: a `rate_limit_event`).
3. **summary.md — met.** Written atomically at finish: heading `# Agent-reported summary`, a notice that the text is the agent's own and unchecked, then the final message. Without a `result` line it says no final message was reported. Nothing in `internal/run` interprets or scores the text. Tests: `TestReplay_realSessionWithEditAndDenialFillsTheRecord` (summary content), `TestReplay_aRunWithoutAResultLineStillGetsASummaryAndNoReport`.
4. **Argv builder — met.** `internal/run/argv.go` `BuildArgv` substitutes `{{prompt}}` (exactly as given; the caller passes the Next line as `resume` prints it) and `{{project_root}}` (absolute) in `agent_launcher.builder` args in one pass, no shell. An empty command returns `agent_launcher.builder.command is required: ...`; the function touches no file, so nothing is written. A nil builder uses `data.DefaultAgentBuilder(quality_gates)` (the command spec is data in `internal/data/launcher_config.go`): `claude -p {{prompt}} --output-format stream-json --verbose --permission-mode acceptEdits [--allowedTools Bash(<first word> *)...] --permission-prompts none`. Tests: `TestBuildArgv`, `TestDefaultAgentBuilder_allowsTheFirstWordOfEachGateOnce`, `TestDefaultAgentBuilder_noGatesAllowsNothingAndOmitsTheFlag`.
5. **Sanitising — met.** `cleanText`/`oneLine` in `sanitize.go` remove CSI/OSC/C1 sequences and other control characters and strip URL credentials from every provider string copied into `run.yml` or `summary.md` (paths, commands, session ID, model, denial text, usage strings, final message). Tests: `TestReplay_textCopiedIntoTheRecordHasControlSequencesAndCredentialsRemoved`, `TestCleanText_removesSequencesButKeepsLines`.
6. **Fixture and gate — met.** Two fixtures captured from real Claude Code 2.1.296 sessions are checked in under `internal/run/testdata/claude-2.1.296/` (`edit-and-denial.jsonl`: read, edit, commands, a real permission denial; `failing-command.jsonl`: a real non-zero exit) and replayed through the T-127 fake provider (`TestMain` mode `replay`). Tests use temp dirs and no network (TEST-04). `make build && make test-fast`: pass (exit 0). No owner waiver of the optional Task Check is recorded; none was requested.

### Commands run

- `claude -p ... --output-format stream-json --verbose --permission-mode acceptEdits --allowedTools ... --permission-prompts none --model haiku` in a scratch project, twice, to capture the fixtures (scanned for emails, tokens, URL credentials: none).
- `go test ./internal/run ./internal/data` (focused iteration); `go test -race -count=1 ./internal/run` (pass); `GOOS=windows go vet ./internal/run` (pass).
- A temporary live test (deleted) ran the default builder through `Start` against the real `claude` binary: state `finished`, with session ID, model, usage, cost and `agent_reported: ok` filled in.
- `make build && make test-fast`: exit 0.

### Files

- Read (Context Files): `internal/run/*` (T-127), `internal/data/launcher_config.go`, `internal/data/config.go`, `internal/resume/resume.go` (lines 40-80), the O-045 Objective, `docs/reviews/agent-orchestration-pass-1.md` (sections 2 and 5), the Guardrails rows for the IDs this Task names.
- Extra reads (logged): `.savepoint/config.yml` (quality_gates shape, to derive allowed tools); `go.mod` (module path).
- `internal/codehealth/runner.go` was not read: T-127 already copied its sanitiser into `internal/run/sanitize.go`.
- Changed: `internal/run/{argv,claude,events,record,run,sanitize}.go`, `internal/run/{claude_test,run_test}.go`, `internal/run/testdata/claude-2.1.296/*.jsonl` (new), `internal/data/launcher_config.go` and `_test.go`, `AGENTS.md` (Codebase Map row for `internal/run`), this Task.

### Limitations and notes for the planner

- Claude Code's stream has no numeric exit code field. A failed command's result text starts `Exit code N` (observed live), and that is what is parsed. A result with `is_error: false` is recorded as exit 0, which is an inference from the provider's convention, not a stated code. A denied command has no exit code.
- Counted as `unparsed`: any top-level type other than `system`, `assistant`, `user`, `result`, `stream_event` (in real sessions, `rate_limit_event`). Housekeeping `system` subtypes (hooks, thinking tokens, retries) are ignored on purpose and not counted.
- Not verified: `stream_event` partial deltas from a real `--include-partial-messages` session (a synthetic line only); a native Windows run (only `GOOS=windows go vet`); Claude versions other than 2.1.296.
- Lists in `run.yml` are capped at 500 entries with `lists_truncated: true`; the full stream remains in `events.jsonl`.
- Drift for the planner: `data.AgentLauncher.Validate` still rejects `enabled: true` with no `builder`, so the default builder is reachable only once T-130 relaxes that rule or resolves the builder itself. `BuildArgv(nil, ...)` already returns the default. T-130 should set `Spec.Provider = run.ProviderClaude` (this selects the adapter) and pass the Next line as `resume` prints it.
- No change was made to the Task's acceptance criteria.

## Drift Notes

Record architecture deltas and reconcile through the planner before Check.
