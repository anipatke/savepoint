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
- Two test providers cannot be proposed with the same report path.
- README Code Health section lists the tool prerequisites. Coordinate wording with I-127.
- Each of these is shown by a test on a project with no tools and no reports, and the existing health tests keep passing.
