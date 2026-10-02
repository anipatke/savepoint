---
id: I-126
title: Migrate, upgrade-assets and doctor do not warn about a stale local savepoint dependency
type: drift
status: open
source:
  kind: report
  actor: {role: owner, session: deep-time-migration-2026-10-03}
  at: '2026-10-03T00:00:00Z'
history:
  - at: '2026-10-03T00:00:00Z'
    actor: {role: executor, session: user-request}
    kind: observed
    note: 'Observed with I-125: deep-time kept savepoint ^1.3.0 in dependencies after migration and nothing warned. Not repaired; deferred by owner.'
---
# I-126: Migrate, upgrade-assets and doctor do not warn about a stale local savepoint dependency

## Summary

After `savepoint migrate --apply` and `upgrade-assets`, a project can still declare an old `savepoint` (below 2) in `package.json` and have it installed in `node_modules`. `npx savepoint` then runs the old binary, so the migrated project appears broken (I-125). Neither command, nor `doctor`, reads the project's package files, so the user gets no hint at the moment they migrate.

## Evidence

deep-time's `package.json` listed `"savepoint": "^1.3.0"` under `dependencies` (a CLI tool beside React) through and after migration. The conflict that kept it from updating was a separate peer-dependency problem, so even a user who tried `npm install savepoint@latest` was blocked and had no explanation tied to Savepoint.

## Proof Needed

When the project root has a `package.json` that declares `savepoint` below 2, `migrate` (preview and apply), `upgrade-assets` and `doctor` each report it with a repair hint: update the dependency, prefer `devDependencies` or running `npx savepoint@latest`. A newer-than-installed check against `node_modules/savepoint/package.json` is optional. No `package.json`, or a 2.x dependency, produces no finding. The check is read-only and must not edit the user's package files.
