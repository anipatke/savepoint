---
id: I-125
title: Savepoint 1.x panics on a V2 project instead of pointing to V2
type: defect
status: open
source:
  kind: report
  actor: {role: owner, session: deep-time-migration-2026-10-03}
  at: '2026-10-03T00:00:00Z'
history:
  - at: '2026-10-03T00:00:00Z'
    actor: {role: executor, session: user-request}
    kind: observed
    note: 'Owner ran npx savepoint health check in the deep-time project after migrating to V2 and got a Go panic from the 1.3.0 binary. Not repaired; deferred by owner.'
---
# I-125: Savepoint 1.x panics on a V2 project instead of pointing to V2

## Summary

A project that migrated to V2 but still has `savepoint` ^1.x in `package.json` runs the old 1.3.0 binary through `npx` (npx prefers the local install). That binary looks for `releases/<id>/epics/`, which V2 Goals do not have, and exits with a Go panic instead of a message. The 2.x code cannot help, because the old binary is what runs. Fixing it needs a 1.x patch release; `^1.3.0` includes `1.3.1`, so every project pinned that way would receive it on its next install.

Related: I-066 is the same `panic(err)` habit in the 2.x dispatcher.

## Evidence

deep-time, after migration: `npx savepoint health check` printed `panic: epics directory not found: stat .../.savepoint/releases/G-002-mobile-viewports/epics: no such file or directory`, from `main.main()` at `main.go:66` of the 1.x build. `node_modules/savepoint` was 1.3.0 and `npx savepoint --version` printed `v0.0.0`. `G-002-mobile-viewports` holds only `Release.md`. The 2.1.2 binary, run in the same project, gives the expected usage message.

## Proof Needed

A `1.3.1` built from the `v1.3.0` tag detects a V2 layout (`.savepoint/router.md` and a `releases/G-###` Goal, or the V2 config) before any V1 loading, prints one line telling the user to run `npx savepoint@latest` or remove the old dependency, and exits non-zero without a panic. Test on a V2 fixture and confirm a real V1 project is unaffected. Publish it as a 1.x release without moving the `latest` dist-tag off 2.x (`npm publish --tag v1`); confirm `^1.3.0` still resolves to it.
