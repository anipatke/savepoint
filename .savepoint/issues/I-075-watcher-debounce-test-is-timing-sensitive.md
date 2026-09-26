---
id: I-075
title: The board watcher debounce test fails when CI stalls
type: verification
status: open
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
