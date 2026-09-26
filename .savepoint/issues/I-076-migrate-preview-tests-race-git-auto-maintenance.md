---
id: I-076
title: Migrate preview tests race Git auto-maintenance in their fixture repository
type: verification
status: resolved
source:
  kind: report
  actor: {role: owner, session: user}
  at: '2026-09-26T04:54:11Z'
severity: low
history:
  - at: '2026-09-26T04:54:11Z'
    actor: {role: owner, session: user}
    kind: observed
    note: >-
      v2 CI run 36219062865 (ae76b3d), Linux ci job, failed
      TestMainMigratePreviewDefaultWritesNothing and
      TestMainMigrateVerboseListsEveryPlannedRecord; the windows-tests job
      passed. Found while fixing I-061 and I-075 before merging v2 into master.
  - at: '2026-09-26T04:54:11Z'
    actor: {role: executor, session: v2-main-flaky}
    kind: repair_attempted
    note: >-
      Both test Git fixture helpers now set maintenance.auto false and gc.auto
      0 before committing. go test -count=20 -run TestMainMigrate and
      go test -count=5 ./internal/migrate passed; make build and make
      test-full passed.
  - at: '2026-09-26T04:54:11Z'
    actor: {role: owner, session: user}
    kind: owner_decision
    note: >-
      Owner instructed the executor to mark fixed Issues resolved before
      pushing. No independent Check was run and no technical CLEAR is
      implied.
resolution:
  disposition: accepted
  actor: {role: owner, session: user}
  at: '2026-09-26T04:54:11Z'
  reason: Owner accepted the test repair committed to v2.
---

# I-076: Migrate preview tests race Git auto-maintenance in their fixture repository

## Summary

The I-069 repair made the migrate preview run the Git working-tree check, so
the preview tests now build a real Git repository in their fixture
(`initMigrateGitRepo` in `main_test.go`, `initGitRepository` in
`internal/migrate/command_test.go`). The fixture's `git commit` can start
Git's detached automatic maintenance, which creates and removes
`.git/objects/maintenance.lock` in the background. The tests snapshot every
file under the fixture, `.git` included, to prove the preview writes
nothing, so the snapshot sometimes sees that lock file appear or vanish.
Product behavior is unaffected.

## Evidence

- Run 36219062865 (v2, ae76b3d), `ci` job:
  `main_test.go:555: snapshot ...: open .../.git/objects/maintenance.lock:
  no such file or directory` in TestMainMigratePreviewDefaultWritesNothing,
  and `main_test.go:586: file count changed: before 50, after 49` in
  TestMainMigrateVerboseListsEveryPlannedRecord.
- The same tests passed on the lane branch and in the merged local
  `make test-full` runs.

## Proof Needed

- The fixture repositories never start background maintenance, so snapshots
  of the fixture are stable.
- The affected tests pass repeatedly; `make build` and `make test-full`
  pass; the next Linux CI run passes.

## Repair Attempt Evidence

- `main_test.go` `initMigrateGitRepo` and
  `internal/migrate/command_test.go` `initGitRepository` set
  `maintenance.auto false` and `gc.auto 0` before `git add` and
  `git commit`.
- `make build` and `make test-full` passed; `go test -count=20 -run TestMainMigrate .` and
  `go test -count=5 ./internal/migrate` passed.
