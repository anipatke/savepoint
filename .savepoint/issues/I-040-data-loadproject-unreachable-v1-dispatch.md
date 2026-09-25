---
id: I-040
title: data.LoadProject and its V1 dispatch are unreachable from the savepoint binary
type: drift
status: resolved
source:
  kind: check
  check: C-914
  actor: {role: checker, session: o021-objective-check-20260924}
  at: '2026-09-23T22:16:20Z'
tasks: [T-022, T-024, T-026, T-027]
checks: [C-914, C-915]
severity: low
resolution:
  disposition: verified
  check: C-915
  actor: {role: checker, session: o021-objective-recheck-20260924}
  at: '2026-09-23T23:04:15Z'
history:
  - at: '2026-09-23T22:16:20Z'
    actor: {role: checker, session: o021-objective-check-20260924}
    kind: observed
    check: C-914
    note: >-
      Reproduced with `deadcode .` at HEAD 10828f2 during the O-021 Full
      Objective Check. T-026's evidence had flagged the same output for the
      Check to assess.
  - at: '2026-09-23T22:59:28Z'
    actor: {role: executor, session: o021-remediation-20260924}
    kind: repair_attempted
    note: >-
      T-027 removed the unused Project wrapper and V1/V2 dispatch helpers,
      retargeted their test callers, and retained migration's direct V1
      readers. The fresh deadcode report has no internal/data findings;
      make build and make test-full passed. Awaiting independent Full
      Objective Check C-914 recheck.
  - at: '2026-09-23T23:04:15Z'
    actor: {role: checker, session: o021-objective-recheck-20260924}
    kind: rechecked
    check: C-915
    note: Independent re-check reproduced the Proof Needed; see C-915.
---

# I-040: data.LoadProject and its V1 dispatch are unreachable from the savepoint binary

## Summary

O-021 Success Condition 1 requires `deadcode .` to report no unreachable
functions in `internal/data` apart from any the owner explicitly keeps and
lists in the Objective. The Objective lists none. `deadcode .` still reports
three:

```text
internal/data/project.go:26:6: unreachable func: LoadProject
internal/data/project.go:43:6: unreachable func: loadProjectV1
internal/data/project.go:50:6: unreachable func: loadProjectV2
```

After T-024, runtime commands use `data.CheckRuntimeSchema` plus
`data.LoadV2Index` (`internal/board/v2/load.go:75`, `internal/board/v2/io.go:35`,
`internal/doctor/v2_runtime.go:40`, `main.go:215`). No production code calls
`LoadProject`; only tests do (`main_test.go`, `main_resume_matrix_test.go`,
`internal/data/project_test.go`, `internal/data/migration_history_test.go`,
`internal/data/migration_source_test.go`, `internal/migrate/apply_test.go`,
`internal/init/v2_scaffold_test.go`).

`loadProjectV1` and the `Project.V1 *Discover` field are also V1-only code that
`internal/migrate` does not use (`internal/migrate/plan.go:491` calls
`data.NewDiscover` directly). That conflicts with Success Condition 8: "Nothing
else in `data` is V1-only."

## Evidence

- `/home/user/go/bin/deadcode .` at HEAD `10828f2` (go1.26.2 linux/amd64). The
  output matches T-026's recorded report.
- `grep -rn 'LoadProject(' --include='*.go' . | grep -v _test.go` matches only
  the definition in `internal/data/project.go:26`.

## Proof Needed

Either option closes this Issue:

1. Delete `LoadProject`, `loadProjectV1`, `loadProjectV2`, and the `Project`
   type, or reduce them to what live code uses. Move the test callers to
   `CheckRuntimeSchema` and `LoadV2Index`. Then show a `deadcode .` report
   with no `internal/data` entries, and a passing fresh `make test-full`.
2. The owner explicitly keeps these functions and lists them in O-021's
   Success Conditions, with a reason. That is a planner/owner change, not a
   Check repair.
