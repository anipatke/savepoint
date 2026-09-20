---
id: E45-safe-migration/T001-decide-how-windows-replaces-a-file
title: Decide how Windows replaces a file
status: done
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

- [x] The task records a decision naming the chosen primitive and the observed behavior behind it: replacing a destination held open by another process, replacing a destination carrying the read-only attribute, and what the destination's attributes look like afterwards.
- [x] The evidence is reproducible: the exact build and run commands are recorded, and the run happened against an NTFS path from the Windows host, not an ext4 path under WSL, because the WSL filesystem does not exercise Windows replacement semantics.
- [x] Each rejected alternative — `os.Rename`/`MoveFileEx`, `ReplaceFileW` via `golang.org/x/sys/windows`, and retry-with-backoff over either — is recorded with the observed reason it was rejected, not a predicted one.
- [x] `internal/migrate` exposes one replacement function with a single documented contract; any behavior that differs between Windows and Unix is explicit in that documentation and covered by a platform-tagged test, satisfying CFG-02.
- [x] The replacement writes a complete same-directory temporary file, flushes and closes it, runs a caller-supplied final check, and only then replaces the destination; a cross-directory temporary file is never used.
- [x] A failure before the replacement leaves the destination byte-identical and removes the temporary file, or returns an error naming the temporary file it could not remove.
- [x] A transient sharing violation is retried a bounded, named number of times; a non-transient error is not retried; exhausted retries return an error naming the path and the underlying condition.
- [x] `AtomicWrite`'s truncating-copy fallback is unreachable from this path, and a test proves the migration replacement never truncates an existing destination before its replacement is durable.
- [x] `golang.org/x/sys/windows` was not adopted; the Context Log records why the standard library sufficed, and `go.mod` and `go.sum` are unchanged.

## Implementation Plan

- [x] Write a self-contained experiment as Go tests in `internal/migrate/replace_windows_test.go` covering: destination open by another handle, destination read-only, destination with non-default attributes, and interruption between temp write and replacement.
- [x] Cross-compile the test binary with `GOOS=windows go test -c ./internal/migrate/` and run it from the Windows host against a directory on an NTFS volume; capture the full output verbatim.
- [x] Repeat the same experiment on Linux to establish the baseline the contract must match.
- [x] Choose the primitive from the observed results and record the decision, the evidence, and the rejected alternatives in the Context Log.
- [x] Implement the chosen primitive in `internal/migrate/replace.go` with the Windows-specific path in `replace_windows.go` behind a build tag, following the `write_linux_test.go` platform-test pattern already in the repository.
- [x] Implement the bounded retry policy only for the transient conditions the experiment actually observed; do not add retries for conditions that did not occur.
- [x] Test the contract on Linux: success, final-check rejection, temp cleanup on failure, and destination integrity after each failure mode.
- [x] Run the focused `internal/migrate` suite, then `make build && make test`.

## Decision

**Windows replaces a protected file with `ReplaceFileW`, bound directly through `syscall.NewLazyDLL`, with a bounded retry over `ERROR_SHARING_VIOLATION` only. Unix uses `os.Rename` with no retry. Both sit behind one exported function, `migrate.ReplaceFile`.**

The decisive observation is one nobody predicted in the epic design: `MoveFileEx` — which is literally what `os.Rename` compiles to on Windows, via `internal/syscall/windows.Rename` — is refused with `ERROR_ACCESS_DENIED` whenever *any* other process holds the destination open, **including under the fully permissive `FILE_SHARE_READ|WRITE|DELETE` share mode that Go's own `os.Open` uses**. `ReplaceFileW` succeeds in exactly that case. A migration runs against a project directory the user may well have open in an editor or indexer, so this is the common case, not the exotic one.

### Observed behavior, both candidates, same inputs

| Case | `MoveFileEx` (`os.Rename`) | `ReplaceFileW` |
|---|---|---|
| Destination held open by another process, share RWD | `ERROR_ACCESS_DENIED(5)` | succeeds |
| Destination held open, share READ only | `ERROR_ACCESS_DENIED(5)` | `ERROR_SHARING_VIOLATION(32)` |
| Destination held open, share none | `ERROR_ACCESS_DENIED(5)` | `ERROR_SHARING_VIOLATION(32)` |
| Destination read-only attribute | `ERROR_ACCESS_DENIED(5)`, destination byte-identical, attribute still `0x1` | `ERROR_ACCESS_DENIED(5)`, destination byte-identical, attribute still `0x1` |
| Destination hidden (`0x2`) before replacement | attribute lost — becomes `0x20` (ARCHIVE) | attribute preserved — `0x22` |
| Destination with an explicit `Everyone:(RX)` ACE | ACE lost; inherited ACEs only | ACE preserved, full ACL identical |
| Destination does not exist | creates it silently | `ERROR_FILE_NOT_FOUND(2)`, creates nothing |
| Process dies after temp is durable, before replacement | destination byte-identical, temp file remains — recoverable, primitive-independent | same |
| Failure while held, retried after the holder releases | succeeds | succeeds |

### What the destination looks like afterwards

After a successful `ReplaceFileW` the destination keeps its own attributes and its own ACL; the replacement file's metadata is discarded. That is why the `mode` argument to `ReplaceFile` decides the replaced file's permission bits on Unix but not on Windows — the one documented platform difference, stated in the `ReplaceFile` doc comment and covered by `TestReplaceFile_modeBecomesTheReplacedFileMode` (Unix) and `TestReplaceFile_preservesDestinationAttributesAndACL` (Windows).

### Rejected alternatives, with the observed reason

1. **`os.Rename` / `MoveFileEx` — rejected.** Observed, not predicted: it fails with `ERROR_ACCESS_DENIED` against a destination any other process holds open even under the most permissive share mode, and it silently discards the destination's attributes and ACL. Either alone disqualifies it for replacing a user's file. It has a third, quieter problem: it reports contention and permission denial under the *same* error code, so no honest retry policy can be built on it.
2. **`ReplaceFileW` via `golang.org/x/sys/windows` — rejected as a dependency, adopted as a call.** Observed by inspecting the module actually in the cache: `golang.org/x/sys@v0.38.0/windows` does not bind `ReplaceFileW` at all, and neither does the standard library's `syscall` package or `internal/syscall/windows`. Taking `x/sys` as a direct dependency would therefore have bought nothing but `NewLazySystemDLL`, which `syscall.NewLazyDLL` already provides. The binding is ~30 lines in `replace_windows.go`. **DEP-01: no new direct dependency was added; `go.mod` and `go.sum` are unchanged.**
3. **Retry-with-backoff as the primitive — rejected; as a bounded policy over `ReplaceFileW` — adopted.** Recording this precisely, because the acceptance criterion lists it flatly as rejected and reality is narrower. Retry over `MoveFileEx` is rejected: it cannot tell a transient lock from a read-only attribute, so it would sit and wait on a condition that never clears. Retry as a *substitute* for choosing a primitive is rejected: no amount of retrying makes `MoveFileEx` preserve an ACL. What is adopted is a bounded retry over `ReplaceFileW` for the single condition the experiment observed to be transient — `ERROR_SHARING_VIOLATION(32)`, which cleared as soon as the holding process released. `ERROR_ACCESS_DENIED(5)` is deliberately *not* retried: it was observed only from the read-only attribute, a standing condition. `ERROR_LOCK_VIOLATION(33)` is not retried either, because the experiment never produced it and the plan forbids retrying conditions that did not occur.

### Reproducing the evidence

Run on 2026-09-17 from this WSL2 checkout against the Windows host:

```
GOOS=windows GOARCH=amd64 go test -c -o migrate-probe.exe ./internal/migrate/
mkdir -p /mnt/c/savepoint-e45/tmp && cp migrate-probe.exe /mnt/c/savepoint-e45/
cd /mnt/c/savepoint-e45 && TMP='C:\savepoint-e45\tmp' TEMP='C:\savepoint-e45\tmp' \
  WSLENV='TMP/w:TEMP/w' ./migrate-probe.exe -test.v
```

The `TMP`/`TEMP` override is what forces `t.TempDir()` onto `C:` rather than onto a `\\wsl.localhost` path, and `TestProbeVolumeIsNTFS` calls `GetVolumeInformationW` and fails the run unless the filesystem under test reports `NTFS` — so the transcript proves where it ran rather than asserting it. Every probe logs a line prefixed `EVIDENCE`. The scratch directory was removed after the run; the commands above recreate it.

The Linux baseline is the same binary built for the host: `go test -v ./internal/migrate/`.

## Context Log

**Files read:** `.savepoint/router.md`, `E45-Detail.md`, this task file, `.savepoint/Guardrails.md`, `internal/data/write.go`, `internal/data/write_linux_test.go`, `internal/init/write.go`, `go.mod`, `AGENTS.md`. Targeted verification reads outside `## Context Files`: `golang.org/x/sys@v0.38.0/windows` and `$GOROOT/src/os/file_windows.go` + `internal/syscall/windows/syscall_windows.go`, to establish that `ReplaceFileW` is unbound in both and that `os.Rename` on Windows is `MoveFileEx` with `MOVEFILE_REPLACE_EXISTING`. `.savepoint/Health-Check.md` is absent, so the Quick health check step does not apply.

**Files created:** `internal/migrate/replace.go`, `replace_unix.go`, `replace_windows.go`, `replace_test.go`, `replace_unix_test.go`, `replace_windows_test.go`.

**Files edited:** `AGENTS.md` (one Codebase Map row for `internal/migrate`, ARCH-04).

**Not edited, deliberately:** `internal/data/write.go` and `internal/init/write.go`. The experiment shows `replaceV2File`'s `os.Rename` carries the same Windows weakness the migration path now avoids, and that `AtomicWrite`'s fallback would truncate; both are outside this task's scope. See `## Drift Notes`.

**Quality gates**

- `go test -v ./internal/migrate/` (linux/amd64) — ok, 13 tests pass, 0 skipped.
- `./migrate-probe.exe -test.v` (windows/amd64, NTFS, Windows host) — PASS, 22 tests: 8 probe tests, 4 Windows contract tests, 10 portable contract tests, of which `TestReplaceFile_symlinkDestinationIsRefused` skips (creating a symlink needs privilege this account lacks) and `TestReplaceFile_doesNotReachAtomicWrite` skips (it scans package source, which does not travel with a cross-compiled binary).
- `GOOS=windows GOARCH=amd64 go build ./...` — ok (REL-01).
- `go vet ./internal/migrate/` for both `GOOS=linux` and `GOOS=windows` — clean.
- `make build && make test` — ok, all packages pass.
- `git diff go.mod go.sum` — empty (DEP-01).

**Named failure-path cases (TEST-02):** `TestReplaceFile_rejectedFinalCheckLeavesDestinationIntact`, `TestReplaceFile_neverTruncatesTheDestination`, `TestReplaceFile_missingDestinationIsRefused`, `TestReplaceFile_directoryDestinationIsRefused`, `TestReplaceFile_symlinkDestinationIsRefused`, `TestReplaceFile_refusedReplacementLeavesDestinationIntact` (Unix), `TestReplaceFile_exhaustedRetriesNameThePathAndCondition` (Windows), `TestReplaceFile_readOnlyDestinationIsNotRetried` (Windows).

**Retry bound:** `replaceAttempts = 5`, `replaceRetryDelay = 100ms`, so a destination nobody releases costs 400ms before a truthful error. `TestReplaceFile_exhaustedRetriesNameThePathAndCondition` measures 0.41s, which is the bound behaving as designed rather than a slow test.

## Drift Notes

**New module not previously in the Codebase Map.** `internal/migrate` is created by this task. A Codebase Map row was added to `AGENTS.md` describing only what the package holds today — the replacement primitive. Later E45 tasks will widen that row as inventory, planning, conversion, and the operation journal land; the row is deliberately narrow now rather than describing a package that does not exist yet.

**A finding about existing code, outside this task's scope.** The experiment establishes something about code this task did not touch: `internal/data`'s `replaceV2File` uses `os.Rename`, and on Windows that call is refused with `ERROR_ACCESS_DENIED` whenever another process holds the record open — even under the permissive share mode — and it discards the destination's attributes and ACL when it does succeed. So the board's V2 status and evidence writes are weaker on Windows than they read as being, and `internal/init`'s `AtomicWrite` additionally falls back to a truncating copy when its rename fails, which on Windows is precisely the common case. Changing either is outside T001, which owns the migration path only. Recording it here so the E45 audit can decide whether it belongs to this release or to a defect of its own.

**No architectural delta.** The package boundary, the preview/apply split, and the publish order in `E45-Detail.md` are unchanged; this task settles only the Open decision the design left bounded.
