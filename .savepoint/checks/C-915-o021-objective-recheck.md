---
id: C-915
scope: {kind: objective, id: O-021}
result: CLEAR
checked_by: {role: checker, session: o021-objective-recheck-20260924}
executed_session: o021-remediation-20260924
checked_at: '2026-09-23T23:04:15Z'
reviewed:
  base_commit: a6283d5
  head_commit: 10828f2d02f0723b3ab0fc745970badb7a57dfd0
  files:
    - 'internal/data/project.go#sha256=84dd641e1d9094a992f6cd177fbb176289f01cc64d6d8b80beb808b73527283b'
    - 'internal/data/project_test.go#sha256=7562218ae1246917fc78b76643211bbf291eadf5cea2ee3fc015af451a94af1c'
    - 'internal/doctor/checks.go#sha256=d6a9e8297b455e6c93ae30cf3c981eb2ac6dfc51d11cf50b3cef2996985b6ec6'
    - '.savepoint/Design.md#sha256=b992605302d44b1b9a00e1d521a0478c21f214ad58c520aa886515be358fd7e0'
    - 'AGENTS.md#sha256=bf472b5e6242e6f51db97908fa7aae756a3d8cd6a4277d559d2c4b3232a57192'
    - 'README.md#sha256=89b833b1f0e10deb37c6e7600df57be58dd500d65f62b7e50fa915d998a169ae'
  dependencies: []
issues: [I-042]
supersedes: C-914
---

# C-915: O-021 Full Objective Check re-check

## Closure Map

- **I-040** (unreachable `LoadProject` and V1 dispatch): closed, `verified`.
- **I-041** (Design names `LoadProject` as the board path): closed, `verified`.
- No original matrix cell remains open or unverified.

## Independence And Scope

This session did not execute T-027 (`o021-remediation-20260924`) or any other
O-021 Task. It is the same checker session that wrote C-914, which the method
allows for a re-check.

This re-check reuses C-914's frozen scope lock and 22-cell matrix without
amendment. The review covers HEAD `10828f2` plus the uncommitted T-026 and
T-027 working-tree changes.

The owned Tasks are T-022 to T-027. T-027 is the remediation Task, marked
`done` under owner waivers.

## Admission Ledger And Results

| # | Re-check item | Prior claim | Frozen cell | Result |
|---|---|---|---|---|
| 1 | `deadcode .` in `internal/data` | I-040 | C-914 #2 | Proven. The report lists only three `internal/migrate` entries and none in `internal/data` |
| 2 | V1-only code left in `data` | I-040 | C-914 #17 | Proven. `Project`, `LoadProject`, `loadProjectV1`, and `loadProjectV2` are deleted. `NewDiscover`, which `internal/migrate/plan.go` uses, remains |
| 3 | Design board load path | I-041 | C-914 #19 | Proven. `Design.md:182` names `CheckRuntimeSchema`, `LoadV2Index`, router decoding, and `ResolveNext`. `LoadProject` appears only in two boundary tests that forbid it |
| 4 | Remediated tests keep coverage | T-027 criterion 2 | C-914 #5, #22 | Proven. `TestCheckRuntimeSchema` covers absent config, missing version, v2, malformed, unsupported, and unrelated version fields. V1 fixture tests call `NewDiscover` directly |
| 5 | Cells #1, #3–#16, #18, #20, #21 | C-914 | Same cells | Proven. The rebuilt binary on a `v1-history` copy repeats the V1 refusal for resume, doctor, board, and upgrade-assets. It also repeats the outside-Git, modified-path, and untracked-path refusals; preview writes nothing; apply sets schema 2; resume works after apply; the printed undo returns the tree to clean; and `schema_version: 99` is refused |
| 6 | TEST-08 gate | T-027 criterion 6 | C-914 #22 | Proven. A fresh `make build` exited 0. A fresh `make test-full` (go1.26.2 linux/amd64) exited 0 with no FAIL lines, including the linux/darwin/windows builds. `git diff --check` is clean |

## Acceptance Coverage

All nine O-021 Success Conditions are **Proven**. T-027's Done When items are
**Proven**, except that the T-027 record itself is malformed, as described
below.

## Materiality

No verdict Issues remain, so no materiality actions are required.

## Closure Precondition (I-042)

`T-027-remove-stale-project-loader-references.md` has two top-level
`check_waiver` keys: the owner-chat waiver at 23:01:49Z and the board waiver at
23:01:55Z. As a result, `LoadV2Index` refuses this repository's project:
`savepoint resume`, `savepoint doctor`, and the board all fail with
`mapping key "check_waiver" already defined`.

This record does not fit any frozen matrix cell, and it is not a code defect in
O-021's scope. It therefore does not change the technical verdict.

It does block closure. The closure rules need every owned Task shown `done`,
and the data resolvers cannot read T-027. I-042 records the problem. The fix
is metadata-only: the owner keeps one of the two waivers. Under TEST-08, this
Check's full-gate evidence stays reusable after that fix.

## Observations (non-blocking)

- `internal/migrate` still has three unreachable functions. This is outside
  SC1 and is unchanged from C-914.
- I-037 is still open. It can no longer be reproduced through the CLI.
- The work is uncommitted: T-026, T-027, this record, and the Issue updates.
- The duplicate-key write has no confirmed cause. It may be a concurrent hand
  edit plus a board write. See I-031.

## Owner Validation Still Needed

1. Repair I-042 by keeping one `check_waiver` in T-027, then confirm that
   `savepoint doctor` is clean.
2. Accept the O-021 outcome against this Check.
