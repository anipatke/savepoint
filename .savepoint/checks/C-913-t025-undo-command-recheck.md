---
id: C-913
scope: {kind: task, id: T-025}
result: CLEAR
checked_by: {role: checker, session: t025-task-check-recheck-20260924}
executed_session: t025-task-check-20260924
checked_at: '2026-09-23T21:18:45Z'
reviewed:
  base_commit: 89ed584
  head_commit: 4f3b2df02eb76c3f0e1669d09faa17655c19cfc6
  files:
    - 'internal/migrate/command.go#sha256=8b4f4a5e3d914483ef5c370026dc25786304352d677d71794e19cb174bdc64b4'
    - 'internal/migrate/command_test.go#sha256=a861aab2cb450a059142ec3f0610ecc05a0e232be8d023aff39b26eb78a69573'
    - cmd/migrate.go
  dependencies: []
issues: []
supersedes: C-912
---

# C-913: T-025 Task Check recheck — migration undo command

## Closure Map

- I-039: closed as `verified` by this CLEAR recheck.
- No other in-scope Issues remain open or unverified.

## Independence And Scope Lock

This is a fresh checker session, separate from `t025-task-check-20260924`,
which made the repair. This requested Task Check uses Quick evidence. It
supersedes C-912 and retains C-912's Scope Lock without amendment: T-025's
eight Done When criteria, User Check, and stated boundaries; the public
`savepoint migrate [dir] [--apply]` path; Git worktree/status and printed
`git restore`/`git clean` behavior; and C-912's 14 matrix cells. Text-width,
TTY/colour, network, and timeout cells remain not applicable for the same
reasons recorded in C-912. No new axes were admitted.

The pre-probe admission ledger was recorded at
`/tmp/savepoint-c913-admission.md`. The table below records the outcomes for
those exact frozen cells. The workflow/side-effect lock from C-912 remains in
force: plan and preflight precede writes; apply writes files, removes sources,
writes the manifest, and activates schema last; failures report written paths
and an undo command. This recheck exercised both undo branches through the
public CLI, including the printed command after a mid-apply failure.

## Coverage Matrix And Admission Results

| # | Frozen C-912 cell | Recheck evidence | Result |
|---|---|---|---|
| 1 | Preview, clean v1-history repo | Public CLI preview; Git status stayed clean. `TestEndToEnd_previewWritesNothing` also passed in the full gate. | Proven |
| 2 | Apply, clean v1-history repo | Public CLI apply succeeded. Re-run `TestRunCommand_realGitApplyIsReversibleAndRejectsDirtyPaths` verified the changed set equals the plan. | Proven |
| 3 | Printed undo after success, tracked paths present | Ran the exact CLI-emitted restore-and-clean string; Git status returned clean and the tracked content was restored. | Proven |
| 4 | Retry apply after undo | Public CLI applied again after undo; the printed undo returned status to clean. | Proven |
| 5 | Second apply to migrated project | Public CLI reported “project is already migrated”; Git status was unchanged. | Proven |
| 6 | Dirty planned `router.md` | Public CLI exited 1, named `.savepoint/router.md`, and left the pre-existing dirty status unchanged. | Proven |
| 7 | Untracked non-planned `notes.txt` | Public CLI apply succeeded and left `notes.txt` outside the planned changes. | Proven |
| 8 | Outside a Git worktree | Public CLI exited 1, advised `git init`, and created no migration manifest. | Proven |
| 9 | Git missing from `PATH` | Public CLI exited 1, advised installing Git, and left status unchanged. | Proven |
| 10 | Mid-apply source-removal failure | With the planned source directory read-only, CLI exited 1, listed written paths, and printed the undo string. Running that exact string restored clean status. `TestApply_failureReportsTouchedPathsAndGitUndo` and the ordered-apply tests passed. | Proven |
| 11 | Empty tracked group, minimal committed project | CLI emitted only `git --literal-pathspecs clean -fdx -- '.savepoint/config.yml' '.savepoint/migrations/v1-to-v2.yml'`. Running it succeeded; Git status was clean, original `Design.md` bytes were unchanged, and generated files were removed. The named regression test passed too. | Proven |
| 12 | Ignored planned path | `TestRunCommand_applyGitRefusalsHappenBeforeAnyWrite/ignored migration path` passed in the re-run migration suite. | Proven |
| 13 | Refusal precedes writeability probe | `TestRunCommand_applyGitRefusalsHappenBeforeAnyWrite` passed; every refusal case observed zero probe calls and unchanged snapshots. | Proven |
| 14 | Golden v1-basic and v1-history output | `TestEndToEnd_goldenConvertedOutput` and `TestEndToEnd_goldenIsReproducible` passed in the fresh full gate. | Proven |

## Acceptance Coverage

| T-025 acceptance item | Classification | Evidence |
|---|---|---|
| Refuse missing Git, non-worktree, and dirty planned paths before writes, with named next steps | Proven | Matrix rows 6, 8, 9, 12, and 13; unit and public CLI probes. |
| Injected clean-tree check and real-Git integration coverage | Proven | Re-run package suite passed; Git and `sh` were available for the real-Git CLI probes. |
| Apply ordering, error report, and usable undo | Proven | Matrix rows 3, 10, and 11; apply ordering/error tests and independent exact-command probes. |
| Already-migrated apply is a no-op | Proven | Matrix row 5 and `TestRunCommand_secondApplyReportsAlreadyMigratedWithoutGit`. |
| Recovery/cutover/replace code and recovery-only manifest fields remain deleted | Proven | Unchanged from C-912 at the same HEAD commit; fresh full gate passed. |
| Golden conversion output remains unchanged | Proven | Matrix row 14 and fresh full gate. |
| README migration guidance remains in place | Proven | Unchanged from C-912 at the same HEAD commit; no README worktree change. |
| Required diff and migration/platform-sensitive gate | Proven | `git diff --check` and fresh `make test-full` both passed. |

## Guardrails

The Task's named rules were applied: FS-01 and TEST-03 are supported by the
clean-tree refusal and content-preservation probes; FS-03 by the read-only
preview; FS-05 by unchanged path construction and the Windows build; FS-06 by
the named refusal paths and writeability tests; DATA-01 by the migration test
suite; ARCH-01 by the unchanged thin `cmd/migrate.go` dispatch; TEST-02 and
TEST-04 by happy/failure coverage in temporary repositories; and TEST-08 by
the fresh full gate. No guardrail violation was found.

## Commands And Results

- `go test ./internal/migrate -run '^TestRunCommand_printedUndoWorksWithNoTrackedPlannedPaths$' -count=1` — PASS.
- `go test ./internal/migrate ./cmd -count=1` — PASS.
- `git diff --check` — PASS.
- `make test-full` — PASS on `go1.26.2 linux/amd64`; all Go packages and Linux, Darwin, and Windows builds passed.
- Public CLI probes in isolated `/tmp` Git repositories — PASS for matrix rows 1–11.

The first scratch attempt put its output log inside the test repository and
was discarded; the final runs stored logs outside each repository. One `go
build` emitted a non-fatal module stat-cache write warning because the cache
was read-only, but returned successfully, produced the binary used for the
probes, and the fresh full gate passed.

## Issues And Materiality

No material Issues remain in the frozen scope. No materiality actions are
required.

## Observations

- C-912's known empty-directory leftovers after undo remain outside the
  acceptance boundary; Git status is clean and retry succeeded.
- The module-cache warning above did not affect the built CLI or full gate.

## Owner Validation Still Needed

`T-025.owner_validation.required` is true. The owner must record acceptance
of this exact current Check (`C-913`) before marking T-025 done. This Task
Check does not close T-025 or O-021; the mandatory Full Objective Check, and
any applicable Goal Check, remain required.
