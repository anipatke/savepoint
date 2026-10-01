---
id: T-055
title: Let one project configure several copies of a tool
objective: O-028
status: done
depends_on: []
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: planning-o028-20261001}
check_waiver:
    task: T-055
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-01T07:40:24Z"
---

# Let one project configure several copies of a tool

## Outcome

The Code Health configuration can hold several named, scoped instances of one capability and provider, mark each instance required or optional, and resolve a provider-specific default timeout, and every snapshot result says which instance produced it, without changing existing comparison series.

## User Check

None beyond the Full Objective Check; this is a schema extension with no owner-visible surface yet.

## Done When

- `CapabilityConfig` gains an optional `name` (short token) and `required` (bool, default false). The instance key becomes capability + provider + name. Two instances with the same key are still refused with `ErrDuplicateInstance`; two instances sharing capability and provider must each have a non-empty `name` and non-identical `scope`, or validation names the problem.
- `CapabilityResult` and `CapabilitySummary` carry the instance `name` (omitted when empty), so existing snapshot identity vectors and series identities are unchanged for unnamed instances. Renaming an instance does not start a new series; scope, exclusions, and the config digest still do.
- The catalogue gives every provider a default timeout (Lizard and jscpd 60s, OSV-Scanner 120s; report-only providers for tests and coverage have none because nothing is executed). An `EffectiveTimeout()` helper returns the override when set, else the default. `MaxTimeoutSecond` rises to 600; out-of-range values are refused.
- Config and snapshot `version` stay 1. Unknown fields are still rejected.
- Classification and history filtering treat each instance separately: two instances never share a trend or combine values.

## Context Files

`.savepoint/objectives/O-028-code-health-provider-execution-and-scaffolding/Objective.md`; `internal/codehealth/model.go`; `internal/codehealth/config.go`; `internal/codehealth/snapshot.go`; `internal/codehealth/identity.go`; `internal/codehealth/classification.go`; `internal/codehealth/model_test.go`; `internal/codehealth/identity_test.go`; `internal/codehealth/classification_test.go`.

## Design References

Design sections 1 and 12; O-026 "Confirmed Code Health Module Boundary"; O-028 Confirmed Design Decisions (2026-10-01).

## Guardrails

DATA-03, ARCH-04, CFG-01, STYLE-07, TEST-01, TEST-02, TEST-06, TEST-08.

## Implementation Plan

1. Add the provider default-timeout table beside `providerCapability` in `model.go`, with a lookup that reports "not executed" for report-only providers. Raise `MaxTimeoutSecond` to 600.
2. Extend `CapabilityConfig` with `Name` and `Required`; validate `Name` as a token; change `instanceKey` to include it; add the shared-provider rule (non-empty distinct names, non-identical scope).
3. Add `Name` (omitempty) to `CapabilityResult` and `CapabilitySummary`, validated as a token; keep it out of `SeriesID` and the config `Digest`.
4. Confirm the snapshot ID vector test still passes unchanged; add a vector with a named instance.
5. Make sure `Assess` history filtering keys on series identity per instance, and `Overall` never merges values across instances.
6. Tests: duplicate named/unnamed instances, same scope refused, name round-trip, timeout default and override and bounds, series unchanged by rename, two instances kept apart in history.

## Boundaries

No discovery, execution, storage layout, CLI, or report parsing. No schema version bump. No change to classification thresholds.

## Technical Verification

Focused `make test-focused TEST=./internal/codehealth/...` while iterating; `make build && make test-fast` at handoff.

## Technical Evidence

Executed 2026-10-01 by the task executor. No extra reads beyond Context Files except `internal/codehealth/history.go`, `primitives.go` (to see `selectSeries`, `validateToken`) and `AGENTS.md`/`savepoint-task` skill for workflow.

Per-criterion outcome:

- **Name/required, instance key, duplicates:** met. `CapabilityConfig.Name` (token) and `Required` (omitempty) added; `instanceKey` is capability+provider+name; same key gives `ErrDuplicateInstance`; shared capability+provider requires a non-empty name and a non-identical scope (compared as sets) or `ErrInvalidConfig` names the instance. Tests: `TestConfigInstances`, `TestConfigNameAndRequiredRoundTrip`.
- **Result/summary carry name; identities unchanged:** met. `CapabilityResult.Name` and `CapabilitySummary.Name` are omitempty; pinned `fixtureSnapshotID` and `fixtureConfigDigest` tests pass unchanged. Name is absent from `seriesKey` and `digestKey`: `TestRenameKeepsSeriesButScopeStartsOne`, `TestNameStaysOutOfConfigDigest`; scope change still starts a new series. Canonical sort and duplicate/summary matching include the name: `TestSnapshotNamedInstances`, `TestSnapshotIDIgnoresInstanceOrder`.
- **Default timeouts:** met. `DefaultTimeoutSeconds` (Lizard/jscpd 60, OSV 120, report-only none), `CapabilityConfig.EffectiveTimeout()`, `MaxTimeoutSecond` 600; out-of-range refused: `TestEffectiveTimeout`, `TestEveryExecutedProviderHasADefaultTimeout`, `TestTimeoutBounds`.
- **Versions and unknown fields:** met. Versions stay 1; unknown config field still rejected (`TestConfigNameAndRequiredRoundTrip`, existing decode tests).
- **Instances kept apart:** met. `Assess` carries the name; history already filters by series identity, which includes scope, so no change to `selectSeries` was needed (filtering by name would contradict rename-keeps-series). `Overall` never merges values. `TestInstancesAreAssessedSeparately`.

Commands run: `go vet ./internal/codehealth`; `go test ./internal/codehealth`; `make build && make test-fast` (passed).

Limitations: a `not_configured` result carrying a name is refused (my addition, beyond the written criteria). Two instances of one provider with identical scope are refused only inside one config, not across a rename over time. `gofmt -l` lists three unrelated pre-existing files (`internal/init/manifest_test.go`, `internal/styles/*`); untouched. Task Check not requested; stage is `audit`, not passed.

## Drift Notes

None expected; the AGENTS.md Codebase Map entry for `internal/codehealth` already covers configuration records.
