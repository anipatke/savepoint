---
id: T-066
title: Decide which health results block an Objective Check
objective: O-030
status: done
depends_on: []
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: planning-o030-20261001}
complexity_tier: medium
complexity_reason: Pure policy over existing snapshot types plus one strict config field; many distinct states to keep separate.
check_waiver:
    task: T-066
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-01T19:47:15Z"
---

# Decide which health results block an Objective Check

## Outcome

Given an official snapshot and the project's health configuration, Code Health produces one deterministic verdict that says, per configured instance, whether it blocks clearance, is reported only, or is not configured, and why, in plain words that never confuse a failed tool with unhealthy code.

## User Check

None beyond the Full Objective Check.

## Done When

- `CapabilityConfig` accepts an optional `blocking` flag. It is valid only for coverage, complexity, and duplication; on tests or dependency vulnerabilities it is rejected with a named config error. Existing configs without it decode unchanged, and the snapshot ID and series digest of an instance without the flag are unchanged.
- A pure evaluation over a snapshot and config (no IO) returns one verdict per result: a failing test count above zero blocks; any critical or high vulnerability blocks; Needs Attention blocks only when `blocking` is set (Watch never does); unknown-severity vulnerabilities are reported as needing review without blocking; a required instance that is failed, timed out, unavailable, absent, unsupported, cancelled, or stale blocks; the same outcomes on an optional instance are reported only; partial evidence is labelled incomplete and blocks only via a default or opt-in rule; `not_configured` is reported as not configured.
- The verdict refuses a manual snapshot with a named error rather than evaluating it.
- Every verdict line names which of these it is: collection failure, incomplete, stale, unhealthy measurement, or not configured. The whole report says "blocks clearance" or "does not block clearance" and never says "clear" or "healthy" for the Objective.
- A plain-text renderer for the verdict is deterministic and is covered by golden-style tests.
- Table tests cover each rule above, including a failing optional test instance (blocks), a failed optional coverage instance (reported), and a stale required test report (blocks).

## Context Files

`.savepoint/objectives/O-030-objective-check-code-health-integration/Objective.md`; `internal/codehealth/config.go`, `internal/codehealth/config_test.go`, `internal/codehealth/identity.go`, `internal/codehealth/model.go`, `internal/codehealth/snapshot.go`, `internal/codehealth/classification.go`, `internal/codehealth/errors.go`, `internal/codehealth/collect.go`; new `internal/codehealth/gate.go`, `internal/codehealth/gate_test.go`.

## Design References

Design section 1 (Code Health readers); R-007 Confirmed Design Decisions; O-030 Confirmed Design Decisions.

## Guardrails

DATA-03, CFG-01, ARCH-02, ARCH-04, STYLE-07, STYLE-09, TEST-01, TEST-02, TEST-04, TEST-07.

## Implementation Plan

1. Confirm `Snapshot`, `CapabilityResult`, `Summary`, and the vulnerability severity detail keys exist as described; return REPLAN REQUIRED if not.
2. Add `Blocking` to `CapabilityConfig` with validation. Make sure `Digest()` keeps the same value when the flag is absent, so no comparison series is reset; record the decision on whether setting it changes the digest.
3. Add the gate types and a pure `Evaluate`-style function that takes the snapshot and config, matching results to instances by capability, provider, and name.
4. Add the renderer; keep the wording in data tables, not scattered through the logic.
5. Write table tests and renderer tests.

## Boundaries

No command, no Check record changes, no skill or doc wording, no new outcomes or classifications, and no change to how `Collect` classifies or saves.

## Technical Verification

Focused `make test-focused TEST=./internal/codehealth/...` while iterating; `make build && make test-fast` at handoff.

## Technical Evidence

Gate: `make build && make test-fast` passed (2026-10-02, after gofmt/go vet clean).

Per criterion:
- `blocking` flag: `CapabilityConfig.Blocking` added; rejected on tests/vulnerabilities with `ErrBlockingNotAllowed` (`TestBlockingFlagValidation`). Configs without it decode unchanged and `Digest()` ignores it, so series are not reset (`TestBlockingFlagDoesNotChangeDigest`). Decision: setting it does not change the digest (policy, not measurement).
- Pure evaluation: `Evaluate` in `gate.go`; every rule is a case in `TestEvaluateRules` (failing optional tests block, failed optional coverage reported, stale required tests block, unknown severity reported, partial blocks only via default/opt-in, Watch never blocks).
- Manual snapshot refused with `ErrManualSnapshot` (`TestEvaluateRefusesManualSnapshot`).
- Each line names its kind; headlines are only "blocks clearance" / "does not block clearance"; test asserts no "healthy"/"clear" wording.
- Deterministic renderer: `TestRenderGolden`.

Files changed: `internal/codehealth/{config.go,errors.go,gate.go,gate_test.go}`. No extra reads beyond Context Files (read `classification.go`/`reader_osv.go` listed or adjacent; `collect.go` only).

Limitations: `Evaluate` trusts a validated snapshot and does not re-validate it; a configured instance with no result is treated as a collection failure; no Task Check requested or waived.

## Drift Notes

Design and Codebase Map wording is reconciled in the skill-and-docs Task, not here.
