---
id: I-075
title: The board watcher debounce test fails when CI stalls
type: verification
status: resolved
source:
  kind: report
  actor: {role: owner, session: user}
  at: '2026-09-26T04:44:12Z'
severity: low
history:
  - at: '2026-09-26T04:44:12Z'
    actor: {role: owner, session: user}
    kind: observed
    note: >-
      v2 CI failed on c9a5e50 (run 36217074871) in
      TestV2WatcherDebouncesRapidWrites with "rapid write burst produced a
      second reload message"; the next push (7adeff8), carrying the same
      code, passed. The owner asked to log and investigate it before merging
      v2 into master.
  - at: '2026-09-26T04:49:37Z'
    actor: {role: executor, session: v2-main-flaky}
    kind: repair_attempted
    note: >-
      Added watchV2FilesAfter so the test sets its own quiet interval;
      watchV2Files keeps the 100 ms production debounce. The test uses a
      500 ms quiet interval against 10 ms write spacing and waits 1 s for a
      second reload. It passed 10 of 10 runs, and with the debounce
      temporarily removed it failed with "rapid write burst produced a second
      reload message". make build and make test-full passed.
  - at: '2026-09-26T04:49:37Z'
    actor: {role: owner, session: user}
    kind: owner_decision
    note: >-
      Owner instructed the executor to mark fixed Issues resolved before
      pushing. No independent Check was run and no technical CLEAR is
      implied.
resolution:
  disposition: accepted
  actor: {role: owner, session: user}
  at: '2026-09-26T04:49:37Z'
  reason: Owner accepted the test repair committed to v2.
---

# I-075: The board watcher debounce test fails when CI stalls

## Summary

`TestV2WatcherDebouncesRapidWrites` writes `config.yml` five times with a
10 ms sleep between writes and expects the watcher to report exactly one
reload, because the production debounce is 100 ms. The test assumes every
10 ms sleep finishes well within 100 ms. On a busy CI runner a single stall of
about 90 ms between two writes is enough for the debounce to fire mid-burst;
the remaining writes then produce a second reload and the test fails. The
product behavior is fine: an extra reload is harmless and the board reloads
correctly. The flaw is the test's timing margin.

## Evidence

- Run 36217074871 (v2, c9a5e50), `ci` job (Linux, `make ci`):
  `watch_test.go:111: rapid write burst produced a second reload message:
  v2.v2FileChangeMsg` and `--- FAIL: TestV2WatcherDebouncesRapidWrites
  (0.56s)`. The 0.56 s run time, against about 0.35 s expected, shows the
  test ran slowly.
- Run 36217625363 (v2, 7adeff8, which contains c9a5e50) passed; so did the
  local `make test-fast` and `make test-full` runs of the same code.
- In the last 60 CI runs this is the only failure of this test.
- Not reproduced locally: `go test -count=40 -run
  TestV2WatcherDebouncesRapidWrites ./internal/board/v2` passed 40 of 40
  idle and 40 of 40 with every CPU saturated (WSL2, 2026-09-26). The
  scheduler-stall explanation fits the failure message and slow run time but
  is inferred, not observed.
- `internal/board/v2/watch.go`: `v2WatchDebounce = 100 * time.Millisecond`;
  `debounceV2Events` resets the timer on each relevant event and returns once
  the filesystem has been quiet for the interval.
- `internal/board/v2/watch_test.go:83-114`: five writes 10 ms apart, then a
  second watch call must stay silent for 200 ms.

## Proof Needed

- The test proves "one reload per burst" without depending on the runner
  keeping every sleep under the production debounce. For example, drive
  `debounceV2Events` with a test-only interval much longer than the write
  spacing (such as 500 ms), or pass the interval into the watch command, so
  a normal CI stall cannot split the burst. Production keeps 100 ms.
- The test still fails if the debounce is removed (each write reloading).
- The test passes repeatedly, for example `go test -count=50 -run
  TestV2WatcherDebouncesRapidWrites ./internal/board/v2` under load, and
  `make build && make test-fast` pass.

## Repair Attempt Evidence

- `internal/board/v2/watch.go`: `watchV2Files` delegates to the new
  `watchV2FilesAfter(watcher, root, delay)`; production behavior and the
  100 ms `v2WatchDebounce` are unchanged.
- `internal/board/v2/watch_test.go` `TestV2WatcherDebouncesRapidWrites` uses a
  500 ms quiet interval, 50 times the 10 ms write spacing, and waits 1 s
  for a second reload.
- `go test -count=10 -run TestV2Watcher ./internal/board/v2` passed. With
  `debounceV2Events` temporarily returning at once, the test failed with
  "rapid write burst produced a second reload message" (then restored).
- `go test -count=50 -run TestV2WatcherDebouncesRapidWrites
  ./internal/board/v2` passed 50 of 50 with every CPU saturated.
- `make build` and `make test-full` passed (2026-09-26T04:49:37Z).
