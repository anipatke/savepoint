# Changelog

## v2.2.0 — Advanced Options and parallel planning

### Added

- **Advanced Options.** Press `o` on the board to open a settings screen with
  one option, **Parallel planning**, off by default. `enter` or `space` saves it
  to `features.parallel_planning` in `.savepoint/config.yml`; `esc`, `q` or `o`
  closes the screen and returns focus to where you were.
- New projects start with `features.parallel_planning: false`. Existing
  projects with no `features` key stay off, and upgrading never rewrites your
  `config.yml`, so an explicit choice survives.
- Saving changes only that one line: comments, key order and unrelated keys are
  kept. If `config.yml` changed since the board loaded it, the save is refused
  and explained, and the screen shows the file as it now stands.

- **Optional parallel-planning metadata.** An Objective may declare `lanes`;
  a Task may name a `lane` and exact `planned_reads` / `planned_writes`.
  Nothing is required or backfilled.
- **Lane / Proposed worktree headings.** With the option on, the selected
  Objective's columns group Tasks under stable lane headings; the Goal-wide
  view namespaces them by Objective.
- **Conservative suggestions with reasons.** Which Tasks may start together,
  and why others are not suggested, appear identically in `savepoint resume`,
  Objective and Task details and the plain board, with a copyable instruction
  for a fresh session per ready Task.
- Planning, task and check guidance now describes the optional advice.
- Repairs from the O-033 check: a lane key declared twice is dropped and its
  Tasks get no suggestion; authored lane titles cannot inject terminal
  controls; a stale router selection withholds launch advice on the board and
  details exactly as in `savepoint resume`; each copyable instruction opens
  with a standalone `Start T-### — …` line; and planning diagnostics name the
  record, file and field that could not be used.

### Notes

- Advice is only advice. It never blocks work, changes Code Health, or
  decides what is done; start, advance and completion decisions are identical
  with it on, off or ignored, and malformed advice degrades to a visible
  diagnostic.
- Safety is limited to what the records state: dependencies, exact planned
  paths and recorded active work. Unknown scope withholds a suggestion, and a
  suggestion is not a guarantee of independence.
- You prepare any worktree, prerequisites and branch yourself. Savepoint
  creates and monitors none, and the owner merges and runs Checks on `main`.
- You can edit the file by hand instead:

  ```yaml
  features:
    parallel_planning: true
  ```

## v2.1.0 — Code Health

Released 2026-10-03.

### Added

- **Code Health.** Five signals (tests, coverage, complexity, duplication,
  dependency vulnerabilities) read from nine report formats: `go test -json`,
  Vitest and pytest JUnit XML, Go cover profile, Vitest V8 and coverage.py
  coverage, lizard CSV, jscpd JSON and osv-scanner JSON.
- `savepoint health setup`, `health check O-###` and `health report`. Checks
  are official (permanent, Full Objective Check only) or manual (board refresh,
  never counts for sign-off). A Check records an official snapshot ID.
- Board: header chip, `H` popover with labels (Good, Watch, Needs attention,
  Unknown), trends, staleness and sign-off wording, `h` history and `R` refresh.
- `.savepoint/health/report.md`, a derived agent-readable report that a
  human can hand to an agent to investigate and fix findings.

### Textual walkthrough

1. Install the tools you want; run `savepoint health setup` and review the
   preview, then `--apply`.
2. Work as usual. Open `savepoint board`, press `H`, move through the five
   signals, `h` for history, `R` to refresh manually.
3. In a Full Objective Check the checker runs the full gate and then
   `savepoint health check O-###`; the Check records the snapshot.
4. Give `.savepoint/health/report.md` to an agent to fix what is not Good.

### Why providers are not bundled

Savepoint stays one binary. It reads reports from tools your project already
owns, so language support grows without Savepoint bundling runtimes, and a tool
that cannot run is reported as unavailable instead of faked. The six release
archives each contain exactly one executable; verified in T-095 on linux/amd64
with no provider installed and no network (`make ci` exit 0; checksums OK).

### Hardening in this release

- Stored evidence: strict version and record validation, bounded record sizes,
  concurrent read/write and failed-replacement tests, `upgrade-assets` that
  preserves health config, snapshots, reports and edited guidance. (T-090, T-092)
- Execution: no shell, bounded output (the tail of long error output is kept),
  timeout and cancellation stop the whole process tree, and credentials inside
  URLs in tool errors are dropped before anything is saved. (T-085, T-091)
- All nine readers exercised together on a Go, JavaScript and Python mix. (T-093)
- Performance: history load was linear in snapshot bytes (84 ms at 1,000
  normal snapshots; 1.2 s and 1.2 GB at 1,000 maximum-size ones). The dashboard
  now validates only a bounded window and reads older files just for their time
  and origin: 33 ms and 0.56 s (170 MB) for the same cases. Render and Git
  freshness are constant. One host (Ryzen 7 7800X3D, go1.26.2, Linux/WSL2, warm
  cache); not measured on Windows or macOS. (T-094, T-097, T-098)
- Internal code health: diagnostic repair, test-stream and board reload
  hotspots refactored with identical behaviour (CCN 46/46/44 to 5/2/4, 31 to 8,
  26 to 7). (T-087, T-088, T-089)

### Health investigation outcomes for this repository

- **Vulnerability scan failure:** `osv-scanner` exits 127 for any error. Network
  failure was reproduced and is the most likely cause, but the original error
  text was lost, so this is inferred, not proven. The lost-error defect is
  fixed. A real measurement awaits the Full Objective Check with network. Until
  then dependency vulnerabilities are **not measured**; nothing here says the
  project is free of them.
- **Duplication:** the whole-repository figure of about 10% was mostly archived
  history, project records and required template mirrors. By owner decision the
  measurement now excludes `.savepoint/**`, `templates/**` and `**/testdata/**`:
  4.51% (4,381 of 97,040 lines), almost all in Go tests. This is a new series; it
  is a scope change, not an improvement, and is not comparable with the old
  figure. New projects get the same exclusions by default.

### Limitations

- Code Health is evidence, not a guarantee. Failed, timed-out or unavailable
  evidence has no value and is never clean. Partial evidence may carry a
  measured value, but it is incomplete and never counts as Good.
- The dashboard's bounded history window can show a shorter trend for a signal
  that was added or removed within the older history.
- Savepoint adds no telemetry or installation, but a project-owned provider may
  use the network (osv-scanner does) and inherits Savepoint's environment.
  Providers are not sandboxed and there is no offline vulnerability mode.
- Supported stacks have readers for the nine formats only; anything else is
  unsupported for that signal. No dead-code analysis, more providers, plugins or
  background monitoring.
- No automatic pruning. Official snapshots are permanent. Internal maintenance
  can keep the ten newest manual snapshots; an owner cleanup command is deferred.
- Windows was verified by cross-compilation and archive structure only on this
  branch; a native `windows-tests` run on the release head is still needed.
- Performance budgets are proposals measured on one machine.
- A few small duplications in production code and about 4,000 duplicated lines
  in Go tests remain; maintained-code duplication is under the 5% watch limit
  and above the 3% aim.

### Follow-up work (separate, not in this release)

- Owner-facing history cleanup command.
- Test-setup helpers to reduce duplication in `init` and `data` tests.
- Optional offline vulnerability mode (needs a provider data download).
