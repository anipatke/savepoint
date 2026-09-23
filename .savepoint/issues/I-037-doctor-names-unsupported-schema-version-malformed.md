---
id: I-037
title: Live doctor names an unsupported schema_version "malformed"
type: defect
status: open
source:
  kind: report
  actor: {role: executor, session: t022-build-20260923}
  at: '2026-09-23T11:00:00Z'
tasks: [T-022]
severity: low
history:
  - at: '2026-09-23T11:00:00Z'
    actor: {role: executor, session: t022-build-20260923}
    kind: observed
    note: >-
      Found while redirecting doctor tests from the deleted CheckProject to
      RunV2Checks during T-022. Left unrepaired because T-022 forbids
      behavior change.
---

# I-037: Live doctor names an unsupported schema_version "malformed"

## Summary

`doctor.RunV2Checks` labels every `data.ReadSchemaVersion` error
`[schema-version-malformed]`. A well-formed but unsupported value such as
`schema_version: 99` (or an explicit `schema_version: 1`) returns
`ErrUnsupportedSchemaVersion`, yet the doctor output still says
`[schema-version-malformed] unsupported schema_version: ...` and gives the
malformed repair. The `[schema-version-unsupported]` branch is reached only
when `schema_version` is absent (implicit V1).

The deleted V1 entry point `CheckProject` routed the same error through
`v2DiagnosticName`, which does name it `schema-version-unsupported`, so the
old `TestCheckProject_SchemaVersionUnsupported` passed without ever testing
the live path.

## Evidence

- `internal/doctor/v2_runtime.go` `RunV2Checks`: the `err != nil` branch after
  `data.ReadSchemaVersion` hard-codes `schema-version-malformed`.
- Observed output for `schema_version: 99`:
  `config.yml: [schema-version-malformed] unsupported schema_version: .../config.yml: schema_version 99`.
- T-022 kept `TestCheckProject_SchemaVersionUnsupported` on the live
  unsupported branch by using a config with no `schema_version`.

## Proof Needed

A doctor test showing `schema_version: 99` reports
`[schema-version-unsupported]` with the unsupported repair, while
`schema_version: nope` still reports `[schema-version-malformed]`. O-021's
T-024 (the shared schema-version check) is the natural place to fix it.
