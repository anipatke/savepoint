---
type: issue-capture-reference
triggerable: false
---

# Shared Savepoint Issue Capture

This reference is not a skill and never triggers on its own. `savepoint-design`, `savepoint-task`, and `savepoint-check` each enter Issue capture as an entry from their own workflow and point here instead of restating these rules; the rule prose lives in this file alone.

An Issue is a record, not a phase: durable follow-up is captured from wherever it is noticed — a planner spotting drift, an executor hitting something out of scope, a checker recording a blocker — without a fourth router state or a separate defect phase. "Defect" stays a word the user says; it maps to `type: defect` and does not resurrect a separate defect record, status, or public phase in V2.

## Issue Artifact Template

Write each Issue file with this structure:

```yaml
id: I-###
title: Issue Title
type: defect|drift|guardrail|verification|other
status: open|in_progress|resolved
source:
  kind: check|report|migration
  check: optional-C-###
  actor: {role: ..., session: ...}
  at: '2026-09-19T00:00:00Z'
tasks: [T-###]
checks: [C-###]
guardrail_ids: [RULE-ID]
severity: optional
resolution:
  disposition: verified|accepted|duplicate|escalated
  check: optional-C-###
  actor: {role: ..., session: ...}
  at: '2026-09-19T00:00:00Z'
  reason: optional
duplicate_of: optional-I-###
escalated_to: optional-O-###
history:
  - at: '2026-09-19T00:00:00Z'
    actor: {role: ..., session: ...}
    kind: observed
    note: optional
    check: optional-C-###
```

```markdown
# I-###: Issue Title

## Summary

## Evidence

## Proof Needed
```

`tasks`, `checks`, `guardrail_ids`, `severity`, `resolution`, `duplicate_of`, `escalated_to`, and `history` are optional; write only the ones the Issue actually has.

`type` is descriptive: it names what kind of durable follow-up this is, and it is never sufficient on its own to block a Task. Material blocking comes from a Check recording `NEEDS WORK` against an acceptance criterion or guardrail, not from an Issue's `type`.

An optional Task Check may be absent under an explicit owner waiver; that
absence is not an Issue and does not waive the mandatory Full Objective Check
or, when a Release exists, the mandatory Release Check.

## Search Before Creating

Before allocating an `I-###`, look for an existing Issue matching the same symptom, the same location, the same violated requirement, or the same linked work. No automatic deduplication is assumed: nothing in Savepoint runs a matching pass for you, so this search is a manual step every capture takes before naming a new ID.

## Resolution Dispositions

A resolved Issue records exactly one disposition:

- **verified** — the Issue was repaired, and the repair is proven by a Check that recorded `CLEAR`.
- **accepted** — an explicit owner decision to close the Issue, including after visual inspection. Record the reason, owner actor, and time. This is an owner decision, not a `CLEAR` Check or independent proof, and it does not waive a mandatory Objective or Release Check. An agent may record the owner's exact decision but may not infer acceptance.
- **duplicate** — the same problem as another, canonical Issue. It names that Issue and proves nothing itself.
- **escalated** — the Issue's repair was promoted into a tracked Objective. It names that Objective in `escalated_to` and proves nothing itself; the Objective's own mandatory Check and owner acceptance carry the proof from here, not a later re-verification of this Issue.

Reopening a recurring problem reuses the same `I-###` with new, dated evidence rather than allocating a new ID.

## Escalation Retires The Issue

When an Issue's repair becomes a new Objective — not a Task inside the current Objective, but an Objective of its own — retire the Issue immediately rather than leaving it open until that Objective's work is later verified: set `status: resolved` with `resolution: {disposition: escalated, escalated_to: O-###, actor: {role: planner, ...}, ...}`, and append a `kind: escalated` history entry naming the Objective. Escalation does not require Check proof at closure time: the promoted Objective carries its own mandatory Full Objective Check and owner acceptance. An owner-directed `accepted` resolution also closes without Check proof, while making no technical clearance claim. `savepoint-design` performs this closure at the moment it plans the remediation Objective; owner-directed `accepted` closure is the other path that needs no Check proof.

## History Is Append-Only

`history` is a frontmatter list of `{at, actor, kind, note, check}` entries: `observed`, `repair_attempted`, `rechecked`, `deferred`, `reopened`, or `owner_decision`. It records the Issue's story in the order it happened. A write that shortens, reorders, or edits a recorded entry is refused — history is appended to, never rewritten.

Deferral is a dated history entry on an open Issue, not a fourth lifecycle state: an Issue stays `open` or `in_progress` while it waits, with the wait recorded as a `deferred` entry naming the reason, never as a new status.

## Role Boundaries

- The **executor** reports repair evidence without independently closing the Issue; it may record an explicit owner-directed `accepted` closure.
- The **checker** verifies Check proof and closes the Issue as `verified`; it may also resolve a confirmed duplicate.
- The **owner** may close the Issue as `accepted` through an explicit decision with reason, actor, and time. An agent records that decision only when directly instructed, without claiming technical `CLEAR`.
- The **planner** (`savepoint-design`) closes an Issue as `escalated` at the moment it promotes the repair into a new Objective — see Escalation Retires The Issue above.

## Out-Of-Scope Repair

Default: fix it directly and record repair evidence in the Issue's own
history (`kind: repair_attempted`), leaving the Issue open for a checker or an explicit owner acceptance decision.
Escalate only when the repair itself needs planning — an open Design
decision, or work spanning multiple Objectives — where the repair becomes a
new Objective of its own. When that happens, retire the Issue immediately
with disposition `escalated` rather than leaving it open (see Escalation
Retires The Issue above). A repair that instead becomes a new, bounded Task
within an existing Objective is not an escalation: the Issue stays open for
the checker as in the default case.
