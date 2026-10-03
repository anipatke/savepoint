---
id: I-129
title: Code Health onboarding gives no guidance on installing tools or producing test reports
type: other
status: open
source:
  kind: report
  actor: {role: owner, session: deep-time-migration-2026-10-03}
  at: '2026-10-03T00:00:00Z'
history:
  - at: '2026-10-03T00:00:00Z'
    actor: {role: executor, session: user-request}
    kind: observed
    note: 'Owner called this an onboarding improvement opportunity after a deep-time health check returned no data for any signal. Not repaired.'
  - at: '2026-10-03T00:00:00Z'
    actor: {role: executor, session: user-request}
    kind: observed
    note: 'Further onboarding evidence from deep-time. (1) pytest-junit and vitest-junit shared junit.xml; with pytest not installed, the pytest line reported 3 failing, the vitest report failures counted twice, and the headline blocked clearance. (2) Coverage was Not suggested because @vitest/coverage-v8 was absent; the message gave no install command or version rule (it must equal the vitest version, 5.0.3 here). (3) The suggested report command, vitest run --coverage --coverage.provider=v8 --coverage.reporter=json, wrote no report while any test failed, because Vitest skips coverage on failure unless --coverage.reportOnFailure=true; the result is the same silent no report found. (4) junit.xml and coverage/ are generated at the project root and are not gitignored. Not repaired.'
  - at: '2026-10-03T00:00:00Z'
    actor: {role: executor, session: user-request}
    kind: repair_attempted
    note: 'Repaired in code, except one decision. (1) health setup: a missing tool says how to install it (lizard: pip install lizard; jscpd: npm install -g jscpd; osv-scanner: its install page), a configured report that is absent names the command that writes it, and a missing @vitest/coverage-v8 says to install it at vitest''s version. (2) health check and report: a collection failure line carries the same remedy, and when no signal produced data a second line says so (No signal produced data, so nothing was judged. Do not read this as a pass.); the headline and verdict contract are unchanged. (3) Setup no longer proposes two test runners with one report path (pytest gets pytest-junit.xml when vitest already uses junit.xml); setup and doctor warn about an existing config that shares a report path and never edit it. (4) The Vitest coverage command now includes --coverage.reportOnFailure=true, and the coverage.py one no longer stops when pytest fails. (5) setup lists generated report files Git does not ignore. (6) README Code Health says who installs tools and produces reports. Decision kept as is: an all-empty health check still saves an official snapshot, because the Check records its ID as evidence that the check ran; the output now says plainly that nothing was judged. Verified against the real deep-time project (doctor and setup show the shared junit.xml warning and the ignore note). go test ./... passes.'
  - at: '2026-10-03T00:40:00Z'
    actor: {role: executor, session: user-request}
    kind: rechecked
    note: 'Rechecked with the published 2.1.4 in a clean scratch project (vitest and pytest declared, lockfile present, none of lizard, jscpd, osv-scanner or pytest on PATH, no reports). health setup: install hint for each missing tool, pytest-junit.xml separate from junit.xml, the vitest coverage package instruction, and the not-ignored report files note. After setup --apply, health check O-001 prints ''No signal produced data, so nothing was judged. Do not read this as a pass.'' and a remedy on every failure line (Install lizard (pip install lizard), jscpd (npm install -g jscpd), the osv-scanner install page; the exact pytest and vitest report commands). Not run: reportOnFailure end to end in a project with a failing test beyond the earlier deep-time run, which wrote coverage-final.json. Executor evidence only; it does not claim verified. The Issue stays open for a checker or an explicit owner decision.'
---
# I-129: Code Health onboarding gives no guidance on installing tools or producing test reports

## Summary

A new user can run `savepoint health check` successfully and still get no evidence, with no indication of what to do next. Savepoint never runs tests, so the two test signals need a JUnit report the user must produce first. The three tool-backed signals need `lizard`, `jscpd` and `osv-scanner` on PATH. Savepoint reports each as "unavailable" or "no report was found" but names no remedy, and the headline still says "Code Health does not block clearance". All signals are optional, so a run with no data at all reads like a pass.

## Evidence

deep-time, 2026-10-02/03. Official snapshots from 22:39 to 22:47 had complexity, duplication and dependency vulnerabilities collected. The 23:51 runs of `savepoint health check O-002` reported all three as "was not found or cannot be run", and both test signals as "no report was found" (`junit.xml` was never produced, in any snapshot). Output: "Code Health does not block clearance." followed by six "reported only; collection failure" lines.

`lizard`, `jscpd` and `osv-scanner` were not installed anywhere on the machine at 23:51. Where they had been available earlier is not established; the first run's `osv-scanner` exit status 127 with a DNS error suggests a different, sandboxed environment. That explanation is inferred, not verified. `python3` also had no `pip` or `ensurepip`, so `lizard` could not be installed by the obvious route.

Both test providers were configured with the same report path, `junit.xml`; if both ran they would overwrite each other. Whether setup proposes that is not checked here.

## Proof Needed

- `health setup` preview lists, for each proposed tool that is not on PATH, how to get it (for example `pip install lizard`, `npm i -g jscpd`, the osv-scanner install page), and for each test provider the command that writes its report at the configured path.
- `health check` and `health report` turn "unavailable" and "no report" into the same actionable hint, and say plainly when no signal produced data, not only "does not block clearance". Decide whether an all-empty check should still save an official snapshot.
- Two test providers cannot be proposed with the same report path, and a provider whose runner is not installed is not counted from another provider's report.
- A "Not suggested" reason for a missing package says how to add it, including the version rule (for Vitest coverage, the same version as vitest).
- The suggested Vitest coverage command includes `--coverage.reportOnFailure=true`, so a first run with a failing test still produces a report; check the other suggested report commands for the same trap.
- Setup suggests ignoring the generated report paths (`junit.xml`, `coverage/`), or says plainly that they are generated files.
- README Code Health section lists the tool prerequisites. Coordinate wording with I-127.
- Each of these is shown by a test on a project with no tools and no reports, and the existing health tests keep passing.
