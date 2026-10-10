---
id: T-120
title: Record what each owner decision covers and whether it still applies
objective: O-044
status: done
depends_on: []
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: planning-o044-20261010}
check_waiver:
    task: T-120
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-10T06:09:31Z"
---

# Record what each owner decision covers and whether it still applies

## Outcome

Task and Objective records can state what an owner acceptance covers and, for each later Check, whether a checker found the acceptance or exception still applies (with the reason, or the material change) or the owner renewed it. A NEEDS WORK Check can list the requirements it found unmet. These fields load, save and fail clearly when malformed; existing records load unchanged.

## User Check

Open a copy of a project whose Objective carries an `exception` and a Task carrying `owner_validation`, add `carried_forward` entries and a Check `unmet` list as documented in O-044 Confirmed Design, and run `savepoint resume`: it loads. Break one entry (for example drop `reason`) and confirm resume names the file and field.

## Done When

1. `owner_validation` decodes optional `scope` and `carried_forward`; both require `accepted_check`. `exception` decodes optional `carried_forward`. Each entry needs `check` (C-###, not the originating Check), `applies`, `assessed_by` (role checker or owner), `assessed_at` (RFC 3339) and `reason`; `material_change` is required when `applies: false`; an owner entry must apply; entries are in Check ID order with at most one per Check per role. Each violation is a named diagnostic (DATA-03).
2. A Check decodes optional `unmet: [requirement IDs]`, with no blank entries, accepted only on a NEEDS WORK Check.
3. Every `carried_forward.check` must name an existing Check scoped to the same record as the originating Check; a dangling or cross-scope reference is a named load error.
4. The evidence writer round-trips the new fields and preserves unknown frontmatter and body text (DATA-01); a record without them writes byte-identically to today.
5. Existing fixtures and TheShed-shaped records without the new fields load and behave exactly as before; no gate decision changes in this Task.
6. Decoder, reference and writer tests cover each happy path and each rejection (TEST-01, TEST-02); `make build && make test-fast` passes (TEST-08). Per-criterion evidence follows AGENTS.md's Verification Policy, including an explicit owner waiver if the optional Task Check is skipped.

## Context Files

`internal/data/evidence_v2.go`, `internal/data/evidence_v2_test.go`, `internal/data/check_v2.go`, `internal/data/check_v2_test.go`, `internal/data/project.go`, `internal/data/project_test.go`, `internal/data/write.go`, `internal/data/write_test.go`, `internal/data/errors.go`, `.savepoint/objectives/O-044-keep-owner-decisions-across-rechecks/Objective.md`, `.savepoint/issues/I-141-recheck-invalidates-owner-acceptance-and-exceptions.md`.

## Design References

O-044 Confirmed Design (where applicability is recorded; unmet requirements). I-141 Root cause and Planned Fix.

## Guardrails

DATA-01, DATA-03, STYLE-04, STYLE-06, STYLE-07, TEST-01, TEST-02, TEST-04, TEST-08.

## Implementation Plan

1. Add `DecisionCarry` and the new fields to `OwnerValidation`, `Exception` and `CheckV2`, with frontmatter structs.
2. Decode and validate entries in one shared helper used by both decision kinds; reuse the existing actor and timestamp decoders.
3. Extend `checkEvidenceReferences` to validate carried Check references and same-scope ownership.
4. Extend the evidence patches in `write.go` to encode the new fields, omitting empty ones.
5. Add decoder, reference and round-trip tests, including a record with none of the new fields.

## Boundaries

Data shape only: no change to completion, dependency, consistency or Next decisions, resume or board output, or skills. Do not edit any recorded Check or existing project record.

## Technical Verification

Focused `make test-focused TEST=...` during iteration; `make build && make test-fast` for handoff. The Full Objective Check needs fresh `make test-full` under `agent-skills/references/check-method.md`.

## Technical Evidence

Executor session, 2026-10-10. Stage advanced build → test → audit; ready for a Check, not passed.

**Per-criterion outcomes**

1. Met. `owner_validation` decodes optional `scope` and `carried_forward`, both rejected without `accepted_check`; `exception` decodes optional `carried_forward`. One shared decoder (`decodeDecisionCarriesV2`, `internal/data/carry_v2.go`) enforces: `check` C-### and not the originating Check, `applies` present, `assessed_by` role checker/owner, RFC 3339 `assessed_at`, `reason`, `material_change` when `applies: false`, owner entries must apply, Check ID order, one entry per Check per role. Each violation names file and field. Tests: `TestDecodeEvidenceV2_carriedForwardRejections` (12 cases), `_carriedForwardOrderAndDuplicates`, `_scopeAndCarriedForwardRequireAcceptedCheck`, `_ownerValidationScopeBlankEntryRejected`, plus the valid cases.
2. Met. `CheckV2.Unmet` decodes `unmet`, rejects blank entries and any use on a CLEAR Check (`TestDecodeCheckV2_unmet`, `_unmetRejections`).
3. Met. `checkEvidenceReferences` now calls `checkDecisionCarryReferences`: dangling carried Check is `ErrV2EvidenceMissingReference`; a Check scoped elsewhere than the originating Check is `ErrV2EvidenceMalformed` naming both scopes (`TestLoadV2Index_carriedForwardReferences`, 6 cases).
4. Met. `ownerValidationV2Patch` and `exceptionV2Patch` write the new fields with `omitempty`; round trip for Task and Objective preserves unknown frontmatter and body (`TestWriteEvidenceV2_roundTripsScopeAndCarriedForwardForTaskAndObjective`); a record without them writes none of the new keys (`TestWriteTaskEvidenceV2_omitsScopeAndCarriedForwardWhenEmpty`), and the existing no-op byte/mtime tests still pass.
5. Met. No gate, resume, board or skill code changed; the whole existing suite passes unmodified.
6. Met. `make build` exit 0 and `make test-fast` exit 0 (go1.26.2 linux/amd64, 2026-10-10T06:09Z). Optional Task Check not requested yet; no owner waiver recorded.

**Files read:** the Task's Context Files (`evidence_v2.go`, `check_v2.go`, `project.go`, `write.go`, `errors.go`, existing evidence/check/project/write tests as needed, O-044 Objective). `I-141` was not read: the Task and Objective already carry its Planned Fix. Extra reads: `AGENTS.md`, `agent-skills/savepoint-task/SKILL.md`, `.savepoint/router.md` (not extra per the skill); `internal/data/discover_test.go` lines 96-110 and 590-595 to reuse existing fixture helpers (test setup only).

**Files changed:** `internal/data/carry_v2.go` (new), `internal/data/carry_v2_test.go` (new), `internal/data/evidence_v2.go`, `internal/data/check_v2.go`, `internal/data/project.go`, `internal/data/write.go`; this Task and O-044 status.

**Limitations:** `make test-full` not run (not required for ordinary handoff; the Full Objective Check needs it). `unmet` is decoded but there is no Check writer in `internal/data` to encode it, so none was added. Carried entries are validated against Checks of the same scope only; whether a carried Check is later than the originating one is not enforced, as the Task did not ask for it.

## Drift Notes

None yet.
