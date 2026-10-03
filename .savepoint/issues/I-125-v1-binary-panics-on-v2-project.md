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
  - at: '2026-10-03T00:00:00Z'
    actor: {role: executor, session: user-request}
    kind: repair_attempted
    note: 'Repair prepared, not released. On branch release/1.3.x (commit 38fb0c8, from v1.3.1) the 1.x binary checks schema_version in the nearest .savepoint/config.yml before any command except --version and init, prints how to run the current version, and exits 1. Reproduced on the real deep-time project: the 1.x build printed the notice for health check, board, doctor and an unknown command, exit 1, with --version unchanged; V1 fixtures (no schema_version, schema_version 1, indented or commented key) are left alone. 1.x go test ./... passes. The branch also adds a dist-tag step to publish.yml so an older major publishes under v<major> and cannot become latest; the same step is on master. Not run: a real publish. The issue stays open until 1.3.2 is released and a checker verifies it.'
  - at: '2026-10-03T03:20:00Z'
    actor: {role: executor, session: user-request}
    kind: rechecked
    note: 'Released as savepoint@1.3.2 (tag v1.3.2, branch release/1.3.x, commit 2fc470a), published by CI under dist-tag legacy-1; latest stayed 2.1.4. Two defects were caught on the way: the guard lived in its own file but the 1.x build tool compiles main.go alone (undefined: v2ProjectNotice in make ci), so it was moved into main.go; and npm rejects v1 as a dist-tag (a valid semver range), which failed the first publish before anything reached the registry, so the workflow now uses legacy-<major> on master and the 1.x branch. Checked on the published package: npx savepoint@1.3.2 health check in a V2 project prints the V2 notice and exits 1; plain npx savepoint is still 2.1.4; ^1.3.0 now resolves to 1.3.2; --version prints v1.3.2. 1.x go test and make ci pass. Caveat: a project whose lockfile pins 1.3.0 only receives 1.3.2 after npm update. Executor evidence only; the Issue stays open for a checker or an owner decision.'
---
# I-125: Savepoint 1.x panics on a V2 project instead of pointing to V2

## Summary

A project that migrated to V2 but still has `savepoint` ^1.x in `package.json` runs the old 1.3.0 binary through `npx` (npx prefers the local install). That binary looks for `releases/<id>/epics/`, which V2 Goals do not have, and exits with a Go panic instead of a message. The 2.x code cannot help, because the old binary is what runs. Fixing it needs a 1.x patch release (1.3.2; 1.3.1 already exists). `^1.3.0` includes it, so a project pinned that way receives it when it next updates or installs without a lockfile, not before.

Related: I-066 is the same `panic(err)` habit in the 2.x dispatcher.

## Evidence

deep-time, after migration: `npx savepoint health check` printed `panic: epics directory not found: stat .../.savepoint/releases/G-002-mobile-viewports/epics: no such file or directory`, from `main.main()` at `main.go:66` of the 1.x build. `node_modules/savepoint` was 1.3.0 and `npx savepoint --version` printed `v0.0.0`. `G-002-mobile-viewports` holds only `Release.md`. The 2.1.2 binary, run in the same project, gives the expected usage message.

## Proof Needed

A `1.3.2` built from the `v1.3.1` tag detects a V2 layout (`.savepoint/router.md` and a `releases/G-###` Goal, or the V2 config) before any V1 loading, prints one line telling the user to run `npx savepoint@latest` or remove the old dependency, and exits non-zero without a panic. Test on a V2 fixture and confirm a real V1 project is unaffected. Publish it as 1.3.2 without moving the `latest` dist-tag off 2.x (published under dist-tag `legacy-1`, because npm rejects `v1`; the publish workflow now chooses the tag); confirm `^1.3.0` still resolves to it. Note: 1.3.1 already exists on npm, and a project whose lockfile pins 1.3.0 only receives the patch after `npm update savepoint`.
