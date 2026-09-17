---
id: E45-safe-migration/T010-stand-down-other-writers
title: Stand down other writers
status: planned
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

- [ ] `upgrade-assets` consults `PendingOperation` before its first write and refuses with the recovery guidance, naming the operation; `--dry-run` still reports without writing.
- [ ] The board's status and router write commands in `internal/board/io.go` refuse while an operation is incomplete and surface the guidance as an error, with the check placed at the write boundary rather than in a rendering path, keeping ARCH-02 intact.
- [ ] Doctor reports an incomplete migration operation as a named diagnostic naming the operation, the step it stopped at, and the recovery command.
- [ ] Doctor suggests no repair that would write while an operation is incomplete, and doctor itself continues to write no project file.
- [ ] A project with no operation directory behaves exactly as it does today; a test asserts the existing upgrade, board write, and doctor behavior is unchanged when nothing is pending.
- [ ] A completed operation does not block anything: only an incomplete one refuses.
- [ ] Each refusal names the operation rather than failing generically, so a user can tell a migration guard apart from an unrelated write error.
- [ ] No V1 or V2 presentation changes: the board renders as it does today, and this task adds no column, badge, or overlay.
- [ ] `AGENTS.md` gains the `internal/migrate` Codebase Map row that ARCH-04 requires for the new package.

## Implementation Plan

- [ ] Add the `PendingOperation` consultation to `internal/init`'s upgrade write path, before the first write and after the dry-run branch.
- [ ] Add the same consultation to the write commands in `internal/board/io.go`, at the boundary only.
- [ ] Add the doctor diagnostic and the repair-suppression rule.
- [ ] Add the `internal/migrate` row to the `AGENTS.md` Codebase Map.
- [ ] Test each refusal, each unchanged-when-absent path, the completed-operation case, and that doctor writes nothing.
- [ ] Run the focused `internal/init`, `internal/board`, and `internal/doctor` suites, then `make build && make test`.

## Context Log

Pending.
