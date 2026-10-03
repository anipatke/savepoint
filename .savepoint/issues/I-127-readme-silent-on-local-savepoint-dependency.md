---
id: I-127
title: README does not say where the savepoint package belongs or why npx can pick a stale copy
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
    note: 'Observed with I-125 and I-126: a migrating user had no guidance that savepoint is a dev tool or that npx prefers a local install. Not repaired; deferred by owner.'
  - at: '2026-10-03T00:00:00Z'
    actor: {role: executor, session: user-request}
    kind: repair_attempted
    note: 'README Upgrading now says: update or remove an existing savepoint dependency (npm install -D savepoint@latest, or npx savepoint@latest), that npx runs the project copy first so an old 1.x install runs instead, and that doctor, migrate and upgrade-assets warn about one (I-126). Short, same tone. Not run: a read-through by the owner.'
---
# I-127: README does not say where the savepoint package belongs or why npx can pick a stale copy

## Summary

The README tells users to run `npx savepoint ...` and, for V1 projects, `npx savepoint migrate`, but never says that a local install takes precedence over the latest published version, or that `savepoint` is a development tool that belongs in `devDependencies` (or nowhere, when run with `npx savepoint@latest`). A V1 project that already depends on `savepoint` ^1.x runs the old binary even when the user follows the README exactly.

## Evidence

README.md lines 42-45 and 111-112 show bare `npx savepoint` commands. deep-time kept the old dependency in `dependencies` and hit I-125 by following them.

## Proof Needed

The README's migration note states: remove or update an existing `savepoint` dependency before migrating; prefer `devDependencies` or `npx savepoint@latest`; and a local copy wins over `npx`. It stays short and matches the project's tone (see the README tone commit). Keep it in step with the wording the I-126 findings use.
