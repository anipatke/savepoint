---
id: E51-first-class-releases/T002-require-release-integration-evidence-and-owner-acceptance
status: done
objective: Decide Release completion from member Objectives, independent Release evidence, material Issues, and acceptance by the owner.
depends_on:
    - E51-first-class-releases/T001-give-releases-identity-without-rebuilding-the-hierarchy
complexity_tier: high
complexity_reason: Extends immutable Check scope and composes evidence, issue, objective, exception, and owner-authority rules.
---

# T002: Require Release integration evidence and owner acceptance

## Problem

A Release record is only first-class if Savepoint can distinguish work in progress from a delivery promise that has actually been integrated and accepted. Reusing Objective status alone would fabricate release readiness, while a separate release-audit model would duplicate the Check and Issue contracts.

## Context Files

- `.savepoint/releases/v2/epics/E51-first-class-releases/E51-Detail.md`
- `internal/data/release_v2.go`
- `internal/data/release_v2_test.go`
- `internal/data/check_v2.go`
- `internal/data/check_v2_test.go`
- `internal/data/evidence_v2.go`
- `internal/data/evidence_v2_test.go`
- `internal/data/objective_gate_v2.go`
- `internal/data/objective_gate_v2_test.go`
- `internal/data/issue_v2.go`
- `internal/data/issue_v2_test.go`
- `internal/data/write.go`
- `internal/data/write_test.go`
- `internal/data/release_gate_v2.go`
- `internal/data/release_gate_v2_test.go`

## Acceptance Criteria

- [x] `Check.scope.kind` accepts exactly `task|objective|release`; a Release Check target must resolve to an indexed Release.
- [x] Release Checks use the existing immutable ID, reviewed-scope, supersession, Issue, freshness, and executor/checker-independence contracts.
- [x] Release completion is refused when the Release has no member Objectives or any member Objective is not validly complete.
- [x] Release completion is refused for missing, unknown, stale, NEEDS WORK, or superseded Release evidence.
- [x] Release completion is refused while a material linked Issue remains unresolved unless a scoped owner exception explicitly covers it.
- [x] A current CLEAR Release Check still requires owner acceptance naming that exact Check before `status: done` is allowed.
- [x] A material Release/Objective requirement change or a superseding Check makes prior freshness or owner acceptance inapplicable without silently rewriting history.
- [x] A migrated legacy-completion reference can explain historical `done`, but is typed and displayed separately and never resolves as a new CLEAR Check.
- [x] Managed acceptance and lifecycle writes preserve unknown fields and body bytes, reject stale source documents, and are idempotent on an unchanged second call.
- [x] Existing Task and Objective gate outcomes are unchanged by the addition of Release scope.

## Implementation Plan

- [x] Extend the Check scope decoder and index links with Release targets.
- [x] Add Release evidence resolution using the existing freshness and latest/superseded Check vocabulary.
- [x] Add the canonical Release completion decision and blocker kinds in `internal/data`.
- [x] Compose member Objective completion, Release Check, Issue, exception, and owner-acceptance results without reimplementing their underlying rules.
- [x] Define and validate the typed historical-completion reference used only by migration.
- [x] Add source-preserving, mtime-guarded Release acceptance and lifecycle writes behind canonical decisions.
- [x] Test every refusal and allowed path, stale/superseded evidence, repeat writes, conflicts, exceptions, and non-regression of existing gates.

## Context Log

- Read: router, E51 detail, this task, Guardrails, and the existing Check, Evidence, Objective gate, Issue, write, project/index, discovery, parser, and gate implementations/tests.
- Edited: `internal/data/check_v2.go`, `internal/data/errors.go`, `internal/data/gate_v2.go`, `internal/data/project.go`, `internal/data/release_v2.go`, `internal/data/write.go`; added `internal/data/release_gate_v2.go` and `internal/data/release_gate_v2_test.go`.
- Verification: Release scope decoding/target resolution, member Objective composition, all five non-current clearance outcomes, unresolved Issue blocking and scoped exceptions, superseded acceptance/freshness, historical completion separation, source-preserving/no-op/stale writes, and existing gate non-regression are covered by `internal/data/release_gate_v2_test.go` plus the existing package suite.
- Quality gates: `go test ./internal/data`, `go test ./...`, and `make build && make test` passed.
