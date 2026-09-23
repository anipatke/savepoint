---
id: C-914
scope: {kind: objective, id: O-021}
result: NEEDS WORK
checked_by: {role: checker, session: o021-objective-check-20260924}
executed_session: unrecorded-o021-executor-sessions
checked_at: '2026-09-23T22:16:20Z'
reviewed:
  base_commit: a6283d5
  head_commit: 10828f2d02f0723b3ab0fc745970badb7a57dfd0
  files:
    - 'AGENTS.md#sha256=bf472b5e6242e6f51db97908fa7aae756a3d8cd6a4277d559d2c4b3232a57192'
    - 'README.md#sha256=89b833b1f0e10deb37c6e7600df57be58dd500d65f62b7e50fa915d998a169ae'
    - '.savepoint/Design.md#sha256=426153dde2560c00c48a2adf49035c1802c8e2f177c5bfe20d4d830e12d4044b'
    - 'internal/data/project.go#sha256=9a566ce491a8010164b042f8bd62cb3396bdb114506d9b68a2d5b16a99bcc3ed'
    - 'internal/data/runtime.go#sha256=bf4188cba65ddbf0496532c35de6d851a1530fa961254da26d3fa6c1b413d52a'
  dependencies: []
issues: [I-040, I-041]
supersedes: null
---

# C-914: O-021 Full Objective Check

## Independence And Scope

This is a new session that started after `/clear`. It did not build any O-021
Task. This Check uses Full evidence.

The review covers HEAD `10828f2`, plus the uncommitted T-026 edits to
`AGENTS.md`, `README.md`, `.savepoint/Design.md`, and the T-026 record. The
baseline is `a6283d5`.

The scope includes all five owned Tasks:

- T-022, T-023, T-024, and T-026 are done under owner Task-check waivers.
- T-025 is done with owner acceptance of C-913.

## Scope Lock

1. **Criteria and gates:** O-021's nine Success Conditions and each owned
   Task's Done When. The Guardrails are ARCH-04, TPL-01, TPL-02, TEST-06, and
   TEST-08. The gate is a fresh `make test-full`.
2. **Entry points:** `savepoint board`, `resume`, `doctor`, `upgrade-assets`,
   and `migrate [--apply]`. The code paths are `data.CheckRuntimeSchema`,
   `data.LoadV2Index`, and `internal/migrate` apply. The guidance files are
   `AGENTS.md`, `Design.md`, and `README.md`.
3. **External effects:** Git work-tree and status checks, the printed undo
   commands, and filesystem writes during apply.
4. **Matrix:** listed below. TTY, colour, text-width, network, and timeout
   cells are not applicable: no Task in this Objective changed rendering or
   touched a network path.
5. **Materiality:** an item counts as an Issue only if it violates a named
   Success Condition, Done When, or Guardrail, and can be reproduced through
   the built binary, `deadcode`, or the repository files.

## Coverage Matrix

| # | Cell | Evidence | Result |
|---|---|---|---|
| 1 | SC1: deadcode in `internal/board` and `internal/doctor` | `/home/user/go/bin/deadcode .` reports no entries | Proven |
| 2 | SC1: deadcode in `internal/data` | Reports `LoadProject`, `loadProjectV1`, and `loadProjectV2` as unreachable; none is listed in the Objective | **Issue I-040** |
| 3 | SC2: V1 board removed; dispatch lives in the smallest package | `internal/board` contains only `board.go` and `debug.go` (89 prod lines) plus `v2/` | Proven |
| 4 | SC3: `internal/migrate` import boundary | No non-test imports in board, board/v2, doctor, init, resume, or cmd. `main.go` uses `migrate.` only in the migrate dispatch (lines 169–179). One test-only import is in `internal/data/next_test.go` | Proven |
| 5 | SC3: schema 1 gets one named diagnostic (resume, board, doctor) | The built binary on a V1 copy of `v1-history` prints `schema_version 1: run savepoint migrate --dry-run, then savepoint migrate --apply` and exits 1 for all three | Proven |
| 6 | SC3: upgrade-assets on V1 | Refuses with migrate guidance and writes nothing | Proven |
| 7 | SC3 boundary: `schema_version: 99` | doctor and resume print `unsupported schema_version … 99` and exit 1 | Proven |
| 8 | SC4: golden fixtures | The diff since `a6283d5` removes only the manifest's `generated_at` and `operation_id` lines, and the golden note explains why. `TestEndToEnd_golden*` passed in the fresh full gate | Proven |
| 9 | SC4: preview is the default | `savepoint migrate` with no flag previews and writes nothing | Proven |
| 10 | SC5: apply outside Git | Exit 1 with `git init` guidance | Proven |
| 11 | SC5: modified planned path | Exit 1, names `M .savepoint/router.md` | Proven |
| 12 | SC5: untracked path in a planned directory | Exit 1, names `?? .savepoint/untracked.md` | Proven |
| 13 | SC5: ignored file at a planned destination | Refused before any write (reported as a destination conflict) | Proven |
| 14 | SC5: clean apply, then undo, then re-apply | Apply succeeds and sets `schema_version: 2` last. Running the exact printed undo string (eval) returns `git status` to clean. After commit, a second apply is a no-op. `resume` works after migration | Proven |
| 15 | SC6: journal, recovery, pending-op, `NextPendingMigration`, cutover preflight, and `replace*.go` removed | No production matches except legacy-directory skip comments in `inventory.go`. `replace*.go` is absent | Proven |
| 16 | SC7: V1 templates, the nine retired skills, `audit-method.md`, and the audit-skill shim removed | `templates/` holds `project-v2` and `prompts`. The nine skill folders and `audit-method.md` were deleted after `a6283d5`. `migrate_audit_skill.go` is absent. `retire_v1_skills.go`, `bubbletea-tui-design`, and the non-V1 `superpowers` folder remain | Proven |
| 17 | SC8: only migration-needed V1 readers remain in `data` | `loadProjectV1` and `Project.V1` are V1-only and `internal/migrate` does not use them | **Issue I-040** |
| 18 | SC9: AGENTS.md Codebase Map | Every row is one or two sentences, with no deleted package or behavior. The closing paragraph names the V1 readers that `internal/migrate` uses | Proven |
| 19 | SC9: Design.md describes the smaller shape | Command table, tree, test table, header, and the clean-Git rule are all corrected. Line 182 still names `data.LoadProject` as the board load path | **Issue I-041** |
| 20 | SC9: before/after line counts | T-026 records per-package counts: 33,432/55,253 → 23,204/37,829. Spot-checked `internal/board` at 89 prod lines | Proven |
| 21 | TPL-01 skill byte parity; template mirroring | `TestV2SkillSetIsCompleteWithByteParity` and `TestProjectGuidanceTemplatesMirrorLiveGuidance` pass in the fresh full gate | Proven |
| 22 | TEST-08 gate | Fresh `make test-full` at 2026-09-24T08:14+10:00 (go1.26.2 linux/amd64) exited 0 with no FAIL lines, including the build-linux/darwin/windows steps. `git diff --check` is clean and `make build` exits 0 | Proven |

## Acceptance Coverage

| Success Condition | Classification |
|---|---|
| 1. deadcode clean in board/doctor/data | Issue (I-040) |
| 2. V1 board gone, thin dispatch | Proven |
| 3. No runtime `migrate` import; one data schema gate | Proven |
| 4. Goldens byte-identical or explained diff; preview default | Proven |
| 5. Clean-Git apply with undo guidance | Proven |
| 6. Transaction machinery removed | Proven |
| 7. V1 templates, skills, and shim removed | Proven |
| 8. Only migrate-needed V1 readers remain in data | Issue (I-040) |
| 9. Guidance reconciled; `make test-full`; size report | Issue (I-041); gate and size report Proven |

Per-Task: T-022, T-023, T-024, and T-025 outcomes are met by rows 1–17,
except the leftover in row 2. T-026's Done When items are met except the
Design line in row 19. T-026's own evidence recorded the deadcode output and
deferred it to this Check.

## Materiality

| Issue | Likelihood | Impact | Materiality | Recommendation |
|---|---|---|---|---|
| I-040 unreachable `LoadProject` + V1 dispatch | High: every deadcode run reports it | Low: no user-visible effect, but it is the Objective's headline success condition | Low | Fix now as one narrow change (delete it and repoint the test callers), or the owner lists it as kept in O-021 |
| I-041 Design names `LoadProject` as the board path | Medium: agents read Design for architecture | Low | Low | Fix together with I-040 |

## Observations (non-blocking)

- `internal/migrate` still has unreachable `FindProjectRoot`,
  `UnmarshalManifest`, and `WriteManifestCreateOnly`. SC1 does not cover
  `internal/migrate`, so these are not Issues. Removing them is optional
  cleanup.
- I-037 (doctor labels an unsupported schema as "malformed") remains open. It
  can no longer be reproduced through `savepoint doctor`, because
  `CheckRuntimeSchema` now fails first with `unsupported schema_version`. Its
  Proof Needed asks for a `RunV2Checks` test, which this Check did not verify.
  The owner may accept it or leave it open.
- `migrate` still generates and previews an operation ID although the manifest
  no longer records one. This is harmless.
- T-026's edits are uncommitted. They should be committed with the
  remediation.

## Owner Validation Still Needed

None beyond closing I-040 and I-041. After remediation, a new Check that
supersedes this one is required before O-021 can close.
