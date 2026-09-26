---
id: I-066
title: An unknown command opens the board, and a board failure panics
type: defect
status: resolved
source:
  kind: report
  actor: {role: owner, session: user}
  at: '2026-09-25T23:46:46Z'
severity: medium
history:
  - at: '2026-09-25T23:46:46Z'
    actor: {role: owner, session: user}
    kind: observed
    note: >-
      "npx savepoint update" with savepoint 2.0.3 in a V1 project printed a Go
      panic with a goroutine trace: "panic: board: schema_version 1: run
      savepoint migrate --dry-run, then savepoint migrate --apply".
  - at: '2026-09-25T23:47:38Z'
    actor: {role: executor, session: i066-repair-20260926}
    kind: repair_attempted
    note: >-
      main.go: the command switch gains a default case that prints
      'unknown command "<word>"' and the usage on stderr and exits 2; a bare
      board failure prints its error on stderr and exits 1 instead of
      panicking. TestMainUnknownCommandIsRejectedWithoutOpeningTheBoard and
      TestMainBareBoardRefusalOnV1ProjectDoesNotPanic both fail on the
      previous main.go and pass with the fix; the existing help, version, and
      V1 board refusal tests still pass. Bare savepoint in this V2 repo still
      renders the board. git diff --check and make build && make test-fast
      passed. Issue remains open for verification or owner acceptance.
  - at: '2026-09-26T03:34:53Z'
    actor: {role: owner, session: user}
    kind: owner_decision
    note: >-
      Owner instructed the executor to mark this Issue resolved. The fix
      shipped in savepoint 2.0.4 on npm (tag v2.0.4), and in real use "npx savepoint update" in a V1 project printed the unknown-command usage without a panic.
      No technical CLEAR is implied.
resolution:
  disposition: accepted
  actor: {role: owner, session: user}
  at: '2026-09-26T03:34:53Z'
  reason: Owner accepted the fix as deployed in savepoint 2.0.4.
---

# I-066: An unknown command opens the board, and a board failure panics

## Summary

`main.go` dispatches known commands in a `switch` with no `default`, so an
unknown first argument such as `update` falls through to the bare board.
When `board.Run()` returns an error, `main` calls `panic(err)`, so a
routine, well-worded refusal (for example the V1-schema migration route)
reaches the user as a Go panic and goroutine trace. `savepoint board`
already prints the same refusal cleanly; only the bare and fall-through
paths panic.

## Evidence

- `main.go:39-95`: the command `switch` has no `default` case.
- `main.go:96-98`: `if err := board.Run(); err != nil { panic(err) }`.
- Owner report: `npx savepoint update` in a V1 project, savepoint 2.0.3,
  `main.main() /home/runner/work/savepoint/savepoint/main.go:97`.

## Proof Needed

- An unknown command prints `unknown command "<word>"` and the usage text
  on stderr, exits 2, and opens no board.
- Bare `savepoint` in a V1 project prints the migration refusal on stderr
  and exits 1 without a panic or goroutine trace.
- Bare `savepoint` still opens the board in a V2 project; known commands
  are unchanged.
- Tests cover both paths; `make build && make test-fast` pass.
