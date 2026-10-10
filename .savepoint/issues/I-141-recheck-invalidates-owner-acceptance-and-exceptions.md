---
id: I-141
title: A re-check invalidates recorded owner acceptance and exceptions
type: defect
status: resolved
resolution:
  disposition: escalated
  actor: {role: planner, session: planning-o044-20261010}
  at: '2026-10-10T06:05:00Z'
  reason: Repair needs design decisions on evidence fields and routing; owner chose a dedicated Objective in G-002.
escalated_to: O-044
source:
  kind: report
  actor: {role: executor, session: recheck-lifecycle-investigation-20261010}
  at: '2026-10-10T05:51:22Z'
history:
  - at: '2026-10-10T05:51:22Z'
    actor: {role: executor, session: recheck-lifecycle-investigation-20261010}
    kind: observed
    note: Reproduced against TheShed O-002 (Checks C-003 to C-006). Root cause traced to exact-Check binding in the completion gates and skill Closure Rules. Owner asked to record it as an Issue and include it in the next release with the rest of G-002; implementation deferred.
  - at: '2026-10-10T06:05:00Z'
    actor: {role: planner, session: planning-o044-20261010}
    kind: escalated
    note: Owner chose a dedicated Objective. Repair promoted to O-044 (Keep owner decisions across re-checks) in G-002; G-002 Boundaries widened to admit it.
---

# I-141: A re-check invalidates recorded owner acceptance and exceptions

## Summary

Owner acceptance counts only when it names the latest Check, and an owner
exception applies only to the Check it names. Any superseding Check therefore
cancels both, even when it only verifies a repair that leaves the accepted
behavior and waived requirements unchanged. The result is an administrative
loop: waive, accept, re-check, then waive and accept again. Waived evidence
also has no path to "ready to close by exception" once a later Check exists,
so resume and the board keep routing to another Check.

## Evidence

Reproduction: TheShed (`~/code/TheShed`), Objective O-002, Checks C-003 to
C-006, Issue I-008.

- C-003 NEEDS WORK found two defects and missing evidence. C-004 and C-005
  proved the defect repairs. The owner waived the remaining browser, race,
  failure-timing and exact-platform cells (I-008 history, 05:33:01Z) and
  accepted T-004 and T-005 against C-005 (05:38:48Z).
- Design was then repaired (documentation only) and C-006 independently
  verified it, changing no implementation, requirement or waiver scope. C-006
  still records NEEDS WORK and its Completion Disposition asks the owner to
  record the exception again and renew acceptance against C-006.
- `savepoint resume` in TheShed prints `Check O-002 — Project registry and
  ports` and "Record the Objective O-002 integration Check." — another Check,
  not owner closure.

Root cause:

1. Runtime binds decisions to one exact Check ID.
   `applicableException` in `internal/data/gate_v2.go` returns nil unless
   `exception.check == latest Check`; `ownerAcceptedCheck` requires
   `accepted_check == latest Check`. `ResolveObjectiveCompletion`,
   `ResolveTaskCompletion`, `ResolveObjectiveDependency`,
   `ResolveTaskDependencyV2` (`dependency.go`), `taskDoneByOwnerDecision`
   and `InspectTaskConsistency` all reuse these, so every superseding Check
   drops the decision. There is no record of whether a decision still applies
   at a later Check.
2. Decision scope is not used. `exception.requirements` is decoded but never
   compared with what a Check found unmet, so the runtime cannot tell "same
   waived requirements" from "new unmet requirement". Owner acceptance has no
   scope field at all.
3. Skill prose restates the exact-Check rule.
   `agent-skills/savepoint-check/SKILL.md` Closure Rules ("acceptance naming a
   Check a later run has superseded does not count"; an exception "applies
   only to the Check it names") and its scaffold copy lead checkers to demand
   renewal, as C-006 did.
4. In TheShed the waiver exists only as an I-008 `owner_decision` history
   note, not as a structured `exception` on O-002, so the gate cannot see it
   at all. Nothing in the Check workflow prompts recording it in structured
   form.

## Planned Fix

Recorded so the work can resume without re-deriving it. A partial diff of
the evidence decoder was set aside, not committed.

- **Scope and provenance.** Keep `exception.check` and
  `owner_validation.accepted_check` as the originating Check. Add optional
  `owner_validation.scope` (accepted behavior or criterion IDs); keep
  `exception.requirements` as the exception's scope.
- **Carry-forward record.** Add an append-only `carried_forward` list to
  `exception` and `owner_validation`: `{check, applies, assessed_by,
  assessed_at, reason, material_change}`. A checker entry assesses whether the
  decision still applies at a later Check; `material_change` is required when
  `applies: false`. An owner entry renews the decision at that Check and keeps
  the original provenance. Entries are appended in Check order, at most one
  per Check per role, never rewritten. This sits on the decision (as
  freshness does), so existing immutable Checks such as C-006 can be assessed
  without being edited.
- **Applicability resolver.** One resolver in `internal/data` returns
  applies / unassessed / materially changed for a decision at the latest
  Check, and every gate listed above uses it instead of exact-ID equality.
- **Scope binding.** Add optional `unmet: [requirement IDs]` to NEEDS WORK
  Checks. An applicable exception grants completion only when it covers every
  unmet ID; an uncovered ID is a named blocker. Older Checks without `unmet`
  rely on the checker's applicability assessment.
- **Clearance vs completion.** Clearance stays NEEDS WORK when evidence is
  waived; never CLEAR. Completion eligibility reports "Ready to close by
  exception" when proof plus applicable owner decisions satisfy it.
- **Routing.** Applicable decisions route resume and the board to owner
  closure (`Close`). An unassessed decision routes to an applicability
  assessment (new `Assess` verb mapped to `savepoint-check`; not a new Check).
  A materially changed decision routes to the owner (`Accept`) for that
  decision only, naming the change. The board's accept action appends an
  owner renewal entry instead of overwriting the originating Check.
- **Guidance.** Update the savepoint-check Closure Rules (and scaffold copy),
  check-method re-check guidance, the AGENTS.md Next-verb list and
  Verification Policy, and README, so the rules match the runtime. Checks
  record exceptions in structured form, not only as Issue history.
- **Existing records.** No blanket reapproval: legacy decisions keep their
  originating Check and become "unassessed" at a later Check until a checker
  records an assessment. No owner decision is created or inferred, and no
  Task or Objective is closed automatically.

## Proof Needed

- Regression (TheShed O-002 shape): an Objective with done Tasks, owner
  acceptance and an evidence exception recorded against C-005; C-005 NEEDS
  WORK finds Design drift plus the waived evidence; a superseding C-006 NEEDS
  WORK verifies the Design repair with the same implementation, requirements
  and exception scope, and a checker carries the decisions forward. Resume and
  the board show `Close` / ready to close by exception, with clearance still
  NEEDS WORK, and ask for no new waiver, acceptance or Check.
- Before the assessment is recorded, the same Objective routes to `Assess`,
  not `Check`.
- A material change to accepted behavior requires renewal of that acceptance
  only; another decision in the same index keeps applying.
- A material change to an exception's scope (an unmet requirement outside
  `exception.requirements`, or a checker `applies: false`) requires renewal of
  that exception only, naming the change.
- An owner renewal entry restores applicability and keeps the originating
  Check.
- Decoder rejects: an owner entry with `applies: false`, a missing reason, a
  missing `material_change` when not applying, out-of-order or duplicate
  entries, and an entry naming the originating Check.
- Prior Checks stay byte-identical; `make build && make test-fast` and
  `make test-full` pass; skill and scaffold copies stay byte-identical.

## Scope Note

G-002's Boundaries currently exclude "new lifecycle states, fields, or
commands" and "Go runtime behaviour changes" beyond O-043. This fix needs new
evidence fields and runtime gate changes, so including it in G-002 requires
the owner to widen those Boundaries or plan it as its own Objective.
The owner chose its own Objective: O-044, with G-002 Boundaries widened.
