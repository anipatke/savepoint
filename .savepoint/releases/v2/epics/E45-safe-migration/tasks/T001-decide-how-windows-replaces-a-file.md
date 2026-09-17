---
id: E45-safe-migration/T001-decide-how-windows-replaces-a-file
title: Decide how Windows replaces a file
status: planned
objective: Settle the platform replacement primitive by observed evidence and implement it behind one documented contract.
depends_on: []
complexity_tier: spike
complexity_reason: Platform behavior must be observed before the primitive and its retry policy can be chosen.
---

# T001: Decide how Windows replaces a file

## Problem

Every later task in this epic replaces files, and none of them can be written honestly until one question is answered: what actually happens on Windows when Savepoint replaces a protected file that another process holds open, or that carries a read-only attribute?

`internal/data/write.go`'s `replaceV2File` writes a same-directory temporary file and calls `os.Rename`, which Go maps to `MoveFileEx` with `MOVEFILE_REPLACE_EXISTING`. That is a reasonable guess about Windows behavior, not evidence of it — there is no Windows test in this repository and CI runs `ubuntu-latest` only. `internal/init/write.go`'s `AtomicWrite` falls back to a truncating copy when rename fails, which is exactly the primitive a migration must never use on a protected replacement: it destroys the original before the replacement is durable.

This task produces the observed evidence, the decision, and the primitive. It is the epic's first task because T008 and T009 build the recovery protocol on top of whatever it decides.

## Context Files

- `internal/data/write.go`
- `internal/data/write_linux_test.go`
- `internal/init/write.go`
- `internal/init/upgrade_failure_test.go`
- `internal/migrate/replace.go`
- `internal/migrate/replace_windows.go`
- `internal/migrate/replace_test.go`
- `internal/migrate/replace_windows_test.go`
- `go.mod`
- `go.sum`
- `.savepoint/Guardrails.md`

## Acceptance Criteria

- [ ] The task records a decision naming the chosen primitive and the observed behavior behind it: replacing a destination held open by another process, replacing a destination carrying the read-only attribute, and what the destination's attributes look like afterwards.
- [ ] The evidence is reproducible: the exact build and run commands are recorded, and the run happened against an NTFS path from the Windows host, not an ext4 path under WSL, because the WSL filesystem does not exercise Windows replacement semantics.
- [ ] Each rejected alternative — `os.Rename`/`MoveFileEx`, `ReplaceFileW` via `golang.org/x/sys/windows`, and retry-with-backoff over either — is recorded with the observed reason it was rejected, not a predicted one.
- [ ] `internal/migrate` exposes one replacement function with a single documented contract; any behavior that differs between Windows and Unix is explicit in that documentation and covered by a platform-tagged test, satisfying CFG-02.
- [ ] The replacement writes a complete same-directory temporary file, flushes and closes it, runs a caller-supplied final check, and only then replaces the destination; a cross-directory temporary file is never used.
- [ ] A failure before the replacement leaves the destination byte-identical and removes the temporary file, or returns an error naming the temporary file it could not remove.
- [ ] A transient sharing violation is retried a bounded, named number of times; a non-transient error is not retried; exhausted retries return an error naming the path and the underlying condition.
- [ ] `AtomicWrite`'s truncating-copy fallback is unreachable from this path, and a test proves the migration replacement never truncates an existing destination before its replacement is durable.
- [ ] If `golang.org/x/sys/windows` is adopted, `go.mod` records it as a direct dependency and the task records the DEP-01 justification; if it is not adopted, the task records why the standard library sufficed.

## Implementation Plan

- [ ] Write a self-contained experiment as Go tests in `internal/migrate/replace_windows_test.go` covering: destination open by another handle, destination read-only, destination with non-default attributes, and interruption between temp write and replacement.
- [ ] Cross-compile the test binary with `GOOS=windows go test -c ./internal/migrate/` and run it from the Windows host against a directory on an NTFS volume; capture the full output verbatim.
- [ ] Repeat the same experiment on Linux to establish the baseline the contract must match.
- [ ] Choose the primitive from the observed results and record the decision, the evidence, and the rejected alternatives in the Context Log.
- [ ] Implement the chosen primitive in `internal/migrate/replace.go` with the Windows-specific path in `replace_windows.go` behind a build tag, following the `write_linux_test.go` platform-test pattern already in the repository.
- [ ] Implement the bounded retry policy only for the transient conditions the experiment actually observed; do not add retries for conditions that did not occur.
- [ ] Test the contract on Linux: success, final-check rejection, temp cleanup on failure, and destination integrity after each failure mode.
- [ ] Run the focused `internal/migrate` suite, then `make build && make test`.

## Context Log

Pending.
