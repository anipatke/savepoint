---
id: T-011
title: Create Tasks through the allocator
objective: O-019
planned_by: {role: planner, session: i024-allocation-design-20260923}
status: done
complexity_tier: high
complexity_reason: Creation must preserve authored content and leave the project loadable across failed writes and concurrent calls.
depends_on: [{task: T-010, requires: clear}]
owner_validation:
    required: false
    accepted_check: ""
check_waiver:
    task: T-011
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-09-24T11:48:16Z"
---

# T-011: Create Tasks through the allocator

## Outcome

One `savepoint create-task` operation accepts an Objective and ID-free Task draft, reserves an ID, writes matching frontmatter and filename, and validates the full V2 index before success.

## Done When

- Refuse an authored ID and occupied file; preserve other supported draft fields and body content.
- Hold the project lock through reservation, file creation, and strict validation. Failure leaves no new Task file; its number stays retired.
- Refuse missing owners and already-invalid projects without modifying existing content.
- Concurrent creation in different Objectives yields distinct IDs and a valid index on Linux and Windows.
- Command dispatch stays thin; filesystem behavior lives in `internal/`.

## Context Files

`main.go`, `main_test.go`, `cmd/init.go`, `cmd/init_test.go`, `internal/data/project.go`, `internal/data/task_v2.go`, `internal/data/write.go`, `internal/data/project_test.go`, `internal/data/task_v2_test.go`, `internal/data/write_test.go`, `AGENTS.md`.

## Design References

Design sections 1, 2, 5, and 6; O-019 Architectural Considerations.

## Guardrails

FS-01, FS-05, FS-06, DATA-01, DATA-03, CFG-01, CFG-02, ARCH-01, ARCH-03, ARCH-04, TEST-01, TEST-02, TEST-03, TEST-04, TEST-05, TEST-08.

## Implementation Plan

1. Specify the ID-free draft contract and success/failure output.
2. Add thin parsing and dispatch to an internal creation operation using T-010's reservation and lock.
3. Create without overwriting content, then strict-load. On failure remove only the new file and retain the reservation.
4. Test malformed input, occupied paths, write/validation failure, invalid projects, and simultaneous requests.

## Boundaries

No general record editor, Task status transition, or planner-selected ID override.

## Technical Verification

Focused command and internal tests, `git diff --check`, and `make build && make test`.

## Technical Evidence

Started 2026-09-24 from the owner-provided `Start T-011` selection. T-011's `requires: clear` dependency on T-010 is met by the owner's recorded Task-check waiver; this satisfies the dependency gate and is not technical `CLEAR`.

Extra read: `.savepoint/objectives/O-019-allocate-task-identities/tasks/T-010-issue-task-identities-once.md`, because the owner explicitly waived T-010's optional Task Check and T-011 declares a `requires: clear` dependency on T-010. The record shows owner completion and the waiver with actor and time.

Additional reads for this executor run: `.savepoint/Design.md` sections 1, 2, 5, and 6, because T-011 cites those architecture sections; `.savepoint/Guardrails.md` rules FS-01, FS-05, FS-06, DATA-01, DATA-03, CFG-01, CFG-02, ARCH-01, ARCH-03, ARCH-04, TEST-01, TEST-02, TEST-03, TEST-04, TEST-05, and TEST-08, because T-011 names them. These establish the CLI/data boundary, preservation and failure obligations, and cross-platform test requirements.

Extra read: `internal/data/task_ids.go`, because creation must keep T-010's project lock held from reservation through file creation and strict validation.

Extra read: `.github/workflows/ci.yml`, because T-011 requires runtime concurrency coverage on Windows and that cross-platform test must run in CI; the workflow already has a Windows test job for the T-010 allocator tests.

Extra read: `internal/data/runtime.go` (`ResolveTarget`, `CheckRuntimeSchema`), because the command must gate V2 projects before creation and resolve the explicitly selected project directory.

Targeted extra fixture lookup: `internal/data/discover_test.go` definitions for `writeV2ObjectiveFixture` and `writeV2TaskFixture`, because internal creation tests need valid linked Objective fixtures for missing-owner and concurrent-creation cases.

Targeted extra read: `internal/data/discover.go` task-directory discovery, to confirm how a pre-existing destination entry can be represented without making the project invalid before the exclusive-create check.

Targeted extra lookup in `internal/data/discover.go` for the Objective filename constant, so the main-command integration test can create a valid V2 project fixture at the canonical path.

Targeted extra declaration search in `internal/data/parser.go` for `ParseV2Document` and `V2SourceDocument`, to determine whether the existing ID-free draft path could reuse that parser. The command returned only declarations; creation will use the split-frontmatter helper in the listed `internal/data/write.go` context.

CLI contract: `savepoint create-task --objective O-### --draft <path> [dir]`; the ID-free V2 draft may omit `objective` (the command injects the selected Objective) or include its exact matching ID. Success prints the allocated Task ID and path relative to `.savepoint`.

Implementation completed: added `CreateTaskV2`, which strict-loads before and under the allocator lock, exclusively creates the ID-named Task, strict-loads the full index while still locked, then rolls back only the newly created file on failure while keeping the reservation. Added thin command parsing/dispatch and tests for byte preservation, invalid drafts/projects, missing owners, occupied destinations, rollback, and concurrent cross-Objective creation. The Windows CI job now runs the concurrent creation test.

Extra read: `agent-skills/references/issue-capture.md`, because the first full gate run exposed an unrelated intermittent board watcher failure and the active task workflow requires Issue capture assessment for durable follow-up. It requires a matching-Issue search before allocating a new ID.

Targeted Issue search: `.savepoint/issues/` for the watcher debounce symptom and test, as required before deciding whether to capture a new Issue; it returned I-020 and I-030.

Targeted extra Issue reads: I-020 and I-030, because the search returned watcher-related Issues and their exact symptoms had to be compared. I-020 covers duplicated reload diagnostics; I-030 covers Windows watcher filtering and guide casing. Neither is the Linux debounce test's second reload message. The debounce test passed on the immediate full-gate rerun, so no Issue was created for this single non-reproducing test failure.

Acceptance evidence:

- ID-free creation, identity injection, byte-preserved body and supported Task fields, authored-ID and occupied-path refusal: covered by `TestCreateTaskV2InjectsIdentityAndPreservesDraftBytes`, `TestCreateTaskV2RejectsInvalidDraftAndMissingObjectiveWithoutReservation`, `TestCreateTaskV2RefusesOccupiedDestinationAndRetiresID`, and `TestCreateExclusiveTaskFileDoesNotOverwriteExistingFile`.
- Reservation held through strict validation; post-write failure removes the new Task and retires its number: implemented through `withTaskIDReservation` and covered by `TestCreateTaskV2StrictValidationFailureRemovesTaskButRetiresID`.
- Missing Objective and invalid project refusals leave prior files unchanged: covered by `TestCreateTaskV2RejectsInvalidDraftAndMissingObjectiveWithoutReservation` and `TestCreateTaskV2InvalidProjectDoesNotChangeExistingContent`.
- Concurrent creation across two Objectives produces distinct IDs and a valid index on Linux: `TestCreateTaskV2ConcurrentAcrossObjectives` passed. The Windows CI job now runs that same test; its test binary cross-compiles successfully, but Windows runtime execution has not occurred in this session and remains CI evidence.
- Thin CLI parsing/dispatch and success output: covered by command parser/runner tests and `TestMainCreateTaskAllocatesIDAndPrintsCreatedPath`.

Verification on Go 1.26.2 linux/amd64:

- Focused command/data/main tests passed: `go test ./cmd ./internal/data . -run 'CreateTask|CreateExclusiveTaskFile|TestMainHelpPrintsV2CommandContract' -count=1`.
- `make build` passed.
- The first `make test-full` run failed only at the unrelated timing-sensitive `internal/board/v2.TestV2WatcherDebouncesRapidWrites`; all new Task creation tests passed. An immediate `make test-full` rerun passed all `./...` packages and Linux, macOS, and Windows builds on 2026-09-24 at 11:43 UTC.
- Windows data test cross-compilation passed: `GOOS=windows GOARCH=amd64 go test -c -o /tmp/savepoint-internal-data.test.exe ./internal/data`.
- `git diff --check` passed.

Files changed for T-011: `.github/workflows/ci.yml`, `cmd/create_task.go`, `cmd/init_test.go`, `internal/data/task_create.go`, `internal/data/task_ids.go` (reservation callback/lock lifetime), `internal/data/project_test.go`, `main.go`, and `main_test.go`, plus this Task record.

Handoff: implementation and configured Linux full gate are complete at stage `audit`; this does not mark T-011 done or mean it passed an independent Task Check. The Windows runtime test awaits CI.

Owner Task-check waiver: task `T-011`; reason: `Owner explicitly waived the optional Task Check before breaking for tonight.`; actor: `{role: owner, session: user}`; recorded_at: `2026-09-24T11:47:33Z`. This waiver is not technical `CLEAR` and does not mark the Task done.

## Drift Notes

If the current decoder cannot preserve the ID-free draft's required data, return REPLAN REQUIRED before changing the record schema.
