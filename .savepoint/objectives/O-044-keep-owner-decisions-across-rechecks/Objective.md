---
id: O-044
title: Keep owner decisions across re-checks
status: done
depends_on: []
release: G-002
priority: high
rank: 5
---

# O-044: Keep owner decisions across re-checks

## Outcome

An owner's acceptance or evidence exception survives a later Check that leaves the accepted behavior, requirements and waived scope unchanged. The later Check records that the decision still applies and why. The owner is asked again only for a decision that a material change actually affected, and is told what changed. When proof plus the owner's still-applicable decisions satisfy completion, `savepoint resume` and the board route to owner closure ("ready to close by exception") instead of another Check. Waived evidence stays unproven and is never reported as `CLEAR`.

## Why

Acceptance and exceptions are bound to one exact Check today, so every re-check cancels them and the owner loops through waive → accept → re-check → waive and accept again. TheShed O-002 hit this across C-003 to C-006: after a documentation-only Design repair was independently verified, resume still said `Check O-002` and C-006 asked the owner to record the same decisions again. Captured as I-141 on 2026-10-10; the owner chose to deliver it as its own Objective in G-002.

## Success Conditions

1. Regression (TheShed O-002 shape): an Objective with done Tasks, owner acceptance and an evidence exception; its Check finds only Design drift; the repair is verified by a superseding independent Check with implementation, requirements and exception scope unchanged. The existing decisions apply and the Objective reads ready to close by exception, with no new waiver, acceptance or Check requested.
2. A material change to accepted behavior requires renewed acceptance of that decision only; a material change to an exception's scope requires renewal of that exception only. Each names the change; other decisions keep applying.
3. Technical clearance and completion eligibility are reported separately: clearance stays `NEEDS WORK` while evidence is waived.
4. Prior Checks stay unchanged and decision history is append-only. No owner decision is created or inferred, and no Task or Objective closes automatically.
5. Existing projects keep their recorded decisions and provenance. A decision recorded against an older Check is assessed for applicability rather than requiring blanket reapproval.
6. The runtime gates, resume, board, workflow skills, scaffold copies and documentation state the same rule; it is not prompt guidance alone.
7. Configured gates pass, including `make test-full`.

## Architectural Considerations

Settled in Confirmed Design below:

- Decisions keep their originating Check (`exception.check`, `owner_validation.accepted_check`) as provenance and gain an append-only `carried_forward` list: a checker entry assesses applicability at a later Check (with a reason, and the material change when it no longer applies); an owner entry renews the decision.
- One applicability resolver in `internal/data` replaces exact-Check equality in Task and Objective completion, dependency, and consistency decisions (DATA-02, STYLE-04).
- Exception scope is checked against the requirements a NEEDS WORK Check records as unmet, when it records them.
- Resume and board routing gain an applicability-assessment action distinct from a new Check, and owner renewal names the affected decision.
- Evidence fields change, so `upgrade-assets` delivers the skill updates and older records keep loading.

## Confirmed Design

Owner, 2026-10-10, Objective decision interview and design confirmation:

- **Where applicability is recorded.** On the decision itself: `exception` and `owner_validation` gain an append-only `carried_forward` list of `{check, applies, assessed_by, assessed_at, reason, material_change}`. A checker entry assesses whether the decision still applies at a later Check; `material_change` is required when it does not. An owner entry renews the decision and always applies. The originating Check stays as provenance, and prior Checks are never edited, so existing records such as TheShed C-006 can be assessed without another Check. `owner_validation` also gains optional `scope` (accepted behavior or criterion IDs); `exception.requirements` remains the exception's scope.
- **Unassessed decisions.** A decision recorded against an earlier Check that has no entry for the latest Check routes to a new `Assess` step (savepoint-check, applicability only, no new Check), not to `Check` and not to owner re-approval.
- **Unmet requirements.** A NEEDS WORK Check may list `unmet: [requirement IDs]`. An applicable exception grants completion only when its requirements cover every listed ID; an uncovered ID is a named blocker. A Check without the list relies on the checker's assessment.
- **Routing.** Covered → `Close` ("ready to close by exception"); unassessed → `Assess`; materially changed → `Accept`, for that decision only, naming the change. Clearance stays `NEEDS WORK` while evidence is waived.
- **Tasks.** Four sequential Tasks: record the fields; apply them in the gates; show them in resume and the board; align the skills and guidance.

## Boundaries

**In scope:**
- Evidence fields and decoding for decision scope and carry-forward, and any Check field needed to state unmet requirements
- Task and Objective completion, dependency and consistency gates in `internal/data`
- Resume wording, Next verbs, board actions and detail for carried, unassessed and changed decisions
- savepoint-check Closure Rules, check-method re-check guidance, AGENTS.md routing and Verification Policy, README, and their scaffold copies
- Regression and decoder tests, including the TheShed O-002 scenario as a fixture

**Out of scope:**
- Editing TheShed's records or recording decisions on its owner's behalf
- Automatic closure of Tasks or Objectives
- Reporting waived evidence as `CLEAR` or relaxing the mandatory Full Objective Check
- Goal-level acceptance or exceptions
