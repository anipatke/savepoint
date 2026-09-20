---
id: E45-safe-migration/T010-stand-down-other-writers
title: Stand down other writers
status: done
objective: Make upgrade-assets, board writes, and doctor repairs refuse while a migration operation is incomplete.
depends_on:
    - E45-safe-migration/T008-make-an-interrupted-migration-recoverable
complexity_tier: medium
complexity_reason: Three existing write boundaries plus one diagnostic, with no new presentation or policy.
---

# T010: Stand down other writers

## Problem

An incomplete migration leaves the project in a recorded intermediate state that only the migration knows how to finish. Savepoint's other write paths do not know that. `upgrade-assets` would refresh assets whose pre-migration copies are sitting in a backup directory, and the board's status and router writes would edit V1 records that migration is midway through archiving. Either one turns a recoverable operation into a conflict.

Three write boundaries exist today: `internal/init`'s asset refresh, `internal/board/io.go`'s status and router writes, and `internal/doctor`'s repair suggestions. Each needs the same one-line question before it writes, and the same answer when the answer is no: name the operation and say how to finish or recover it.

The temptation here is to do more — to render migration state in the board, or to have doctor offer to repair it. Neither belongs in this task. E49 owns V2 presentation, and an incomplete migration is finished by the migration command, not by a repair suggestion.

## Context Files

- `internal/doctor/checks.go`
- `internal/doctor/repairs.go`
- `internal/doctor/checks_test.go`
- `internal/board/io.go`
- `internal/board/io_test.go`
- `internal/init/upgrade.go`
- `internal/init/upgrade_test.go`
- `internal/migrate/operation.go`
- `AGENTS.md`

## Acceptance Criteria

- [x] `upgrade-assets` consults `PendingOperation` before its first write and refuses with the recovery guidance, naming the operation; `--dry-run` still reports without writing.
- [x] The board's status and router write commands in `internal/board/io.go` refuse while an operation is incomplete and surface the guidance as an error, with the check placed at the write boundary rather than in a rendering path, keeping ARCH-02 intact.
- [x] Doctor reports an incomplete migration operation as a named diagnostic naming the operation, the step it stopped at, and the recovery command.
- [x] Doctor suggests no repair that would write while an operation is incomplete, and doctor itself continues to write no project file.
- [x] A project with no operation directory behaves exactly as it does today; a test asserts the existing upgrade, board write, and doctor behavior is unchanged when nothing is pending.
- [x] A completed operation does not block anything: only an incomplete one refuses.
- [x] Each refusal names the operation rather than failing generically, so a user can tell a migration guard apart from an unrelated write error.
- [x] No V1 or V2 presentation changes: the board renders as it does today, and this task adds no column, badge, or overlay.
- [x] `AGENTS.md` gains the `internal/migrate` Codebase Map row that ARCH-04 requires for the new package.

## Implementation Plan

- [x] Add the `PendingOperation` consultation to `internal/init`'s upgrade write path, before the first write and after the dry-run branch.
- [x] Add the same consultation to the write commands in `internal/board/io.go`, at the boundary only.
- [x] Add the doctor diagnostic and the repair-suppression rule.
- [x] Add the `internal/migrate` row to the `AGENTS.md` Codebase Map.
- [x] Test each refusal, each unchanged-when-absent path, the completed-operation case, and that doctor writes nothing.
- [x] Run the focused `internal/init`, `internal/board`, and `internal/doctor` suites, then `make build && make test`.

## Context Log

**Files read:** `internal/migrate/operation.go` (PendingOperation/RecoveryGuidance contract), `internal/init/upgrade.go`, `internal/board/io.go` + `update.go`, `internal/doctor/checks.go` + `repairs.go` + `report.go`, their `_test.go` files, `.savepoint/Guardrails.md` (ARCH-02, ARCH-04, FS-01, FS-06), `AGENTS.md` (Codebase Map row for `internal/migrate` already present from an earlier task — no change needed here).

**Files edited:**
- `internal/init/upgrade.go`: added `refusePendingMigration`, wired via the existing `beforeFirstWrite` seam used by the manifest-writability guard — a dry run never calls the wrapped `write`, so it keeps previewing normally with a migration pending, matching AC1.
- `internal/board/io.go`: added `pendingMigrationMsg(root)` and consulted it as the first statement in every write `tea.Cmd` closure — `writeRouterTaskCmd`, `writeRouterReleaseEpicCmd`, `writeTaskStatusCmd`, `writeDefectStatusCmd`, `writeEpicStatusCmd`. The latter three gained a `root` parameter; `update.go`'s five call sites now pass `m.Root`. The check runs inside the already-IO-performing closure, never in `Update()`, so ARCH-02 holds.
- `internal/doctor/checks.go`: added `CheckMigration(root)`, reporting `[migrate-operation-incomplete]` or `[migrate-multiple-operations]` Problems with a repair that only names `savepoint migrate --recover` — doctor never runs it.
- `internal/doctor/report.go`: wired `CheckMigration` into `DiagnosticReport` (`Migration` field, `RunAllChecks`, `HasProblems`, `Format`'s "Migration Check" section).

**Root path convention:** `internal/board` and `internal/doctor` pass the `.savepoint` directory as `root` (matching every existing check in both packages), so each guard calls `migrate.PendingOperation(filepath.Dir(root))` to reach the project directory `PendingOperation` expects. `internal/init`'s `absTarget` is already the project directory, so it calls `migrate.PendingOperation(absTarget)` directly.

**Tests added:**
- `internal/init/upgrade_test.go`: refusal on a real write with a pending operation (file untouched, error names the op); dry run still reports normally with a pending operation.
- `internal/board/io_test.go`: refusal (message names the operation, target file untouched) for `writeEpicStatusCmd`, `writeDefectStatusCmd`, `writeTaskStatusCmd`, `writeRouterTaskCmd`, `writeRouterReleaseEpicCmd`. The router-cmd tests write no `router.md` at all, proving the guard fires before any read.
- `internal/doctor/checks_test.go`: `CheckMigration` none/incomplete/multiple-operations cases.
- `internal/doctor/report_test.go`: `RunAllChecks` reports the incomplete operation as a problem naming it and the recovery command, and a before/after file listing of the project directory proves the run wrote nothing; "Migration Check" added to the format-sections assertion.
- The "completed operation blocks nothing" AC is covered transitively: `PendingOperation` already returns `nil` once `apply.go` removes the operation directory after a successful, activated migration (proven by `internal/migrate`'s own `assertNoPendingOperation` tests from T008/T009), and every guard here is a direct call to that same detector — no separate end-to-end re-test was needed.
- All pre-existing tests in the three packages pass unmodified, which is the "unchanged when nothing is pending" proof for every write path not given its own new test.

**Quality gates:** `go build ./...`, `go vet ./...`, `make build && make test` all clean. No `.savepoint/Health-Check.md` in this project, so the Quick check step was skipped per the build-task skill.

**Not done (out of scope / not applicable):** No TUI session was run in this environment, so the router-priority `p` keypress mentioned in the skill's workflow step 4 was not exercised interactively; it is unrelated to this task's guard logic.
