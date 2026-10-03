---
id: I-128
title: Every npm publish warns that the bin script name was invalid and removed
type: other
status: resolved
source:
  kind: report
  actor: {role: owner, session: deep-time-migration-2026-10-03}
  at: '2026-10-03T00:00:00Z'
resolution:
  disposition: accepted
  actor: {role: owner, session: user-request}
  at: '2026-10-03T03:10:00Z'
  reason: 'Owner accepted after the executor recheck on released Savepoint 2.1.4 (see the rechecked history entry). Not a CLEAR Check.'
history:
  - at: '2026-10-03T00:00:00Z'
    actor: {role: executor, session: user-request}
    kind: observed
    note: 'Observed in the 2.1.1 and 2.1.2 Publish Package logs. Registry metadata checked: bin is intact (savepoint -> bin/savepoint.js). Not repaired; deferred by owner.'
  - at: '2026-10-03T00:00:00Z'
    actor: {role: executor, session: user-request}
    kind: repair_attempted
    note: 'package.json bin is now bin/savepoint.js. Reproduced the warning with the old ./bin/savepoint.js spelling in a scratch package (npm publish --dry-run prints auto-corrected / bin[savepoint] was invalid and removed) and confirmed the new spelling prints none. npm test passes. Not run: a real publish; confirm the next Publish Package log has no warning.'
  - at: '2026-10-03T00:40:00Z'
    actor: {role: executor, session: user-request}
    kind: rechecked
    note: 'Rechecked on the real publish: the Publish Package log for run 37082389282 (savepoint 2.1.4, published) contains no auto-corrected or bin[savepoint] invalid-and-removed warning (0 matching lines), where 2.1.1 and 2.1.2 had them. Executor evidence only; it does not claim verified. The Issue stays open for a checker or an explicit owner decision.'
  - at: '2026-10-03T03:10:00Z'
    actor: {role: executor, session: user-request}
    kind: owner_decision
    note: 'Owner directed: accept I-128. Recorded as accepted, not verified; no CLEAR Check backs it.'
---
# I-128: Every npm publish warns that the bin script name was invalid and removed

## Summary

`package.json` declares `"bin": {"savepoint": "./bin/savepoint.js"}`. On each publish npm prints `npm warn publish "bin[savepoint]" script name bin/savepoint.js was invalid and removed` and asks for `npm pkg fix`. The published package is fine, so this is noise, but it reads like a broken release and will hide a real warning.

## Evidence

Publish logs of runs 37074490185 (2.1.1) and 37078265982 (2.1.2) show the warning. `https://registry.npmjs.org/savepoint/2.1.2` reports `bin: {savepoint: bin/savepoint.js}`, and `npx savepoint@2.1.2 --version` prints `v2.1.2`.

## Proof Needed

`npm pkg fix` (or editing the path to `bin/savepoint.js`) leaves `package.json` valid, and `npm publish --dry-run` prints no `publish errors corrected` warning. `bin/savepoint.test.js` and `make ci` still pass, and `npx savepoint` still resolves the wrapper after a real publish. Do this with the next release; it needs no release of its own.
