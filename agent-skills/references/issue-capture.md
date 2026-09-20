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
id: I###
title: Issue Title
type: defect|drift|guardrail|verification|other
status: open|in_progress|resolved
source:
  kind: check|report|migration
  check: optional-C###
  actor: {role: ..., session: ...}
  at: '2026-09-19T00:00:00Z'
tasks: [T###]
checks: [C###]
guardrail_ids: [RULE-ID]
severity: optional
resolution:
  disposition: verified|accepted|duplicate
  check: optional-C###
  actor: {role: ..., session: ...}
  at: '2026-09-19T00:00:00Z'
  reason: optional
duplicate_of: optional-I###
history:
  - at: '2026-09-19T00:00:00Z'
    actor: {role: ..., session: ...}
    kind: observed
    note: optional
    check: optional-C###
```

```markdown
# I###: Issue Title

## Summary

## Evidence

## Proof Needed
```

`tasks`, `checks`, `guardrail_ids`, `severity`, `resolution`, `duplicate_of`, and `history` are optional; write only the ones the Issue actually has.

`type` is descriptive: it names what kind of durable follow-up this is, and it is never sufficient on its own to block a Task. Material blocking comes from a Check recording `NEEDS WORK` against an acceptance criterion or guardrail, not from an Issue's `type`.

An optional Task Check may be absent under an explicit owner waiver; that
absence is not an Issue and does not waive the mandatory Full Objective Check
or, when a Release exists, the mandatory Release Check.

## Search Before Creating

Before allocating an `I###`, look for an existing Issue matching the same symptom, the same location, the same violated requirement, or the same linked work. No automatic deduplication is assumed: nothing in Savepoint runs a matching pass for you, so this search is a manual step every capture takes before naming a new ID.

## Resolution Dispositions

A resolved Issue records exactly one disposition:

- **verified** — the Issue was repaired, and the repair is proven by a Check that recorded `CLEAR`.
- **accepted** — an explicit owner decision to accept the risk. This is an owner decision, not a repair, and it proves nothing.
- **duplicate** — the same problem as another, canonical Issue. It names that Issue and proves nothing itself.

Reopening a recurring problem reuses the same `I###` with new, dated evidence rather than allocating a new ID.

## History Is Append-Only

`history` is a frontmatter list of `{at, actor, kind, note, check}` entries: `observed`, `repair_attempted`, `rechecked`, `deferred`, `reopened`, or `owner_decision`. It records the Issue's story in the order it happened. A write that shortens, reorders, or edits a recorded entry is refused — history is appended to, never rewritten.

Deferral is a dated history entry on an open Issue, not a fourth lifecycle state: an Issue stays `open` or `in_progress` while it waits, with the wait recorded as a `deferred` entry naming the reason, never as a new status.

## Role Boundaries

- The **executor** reports repair evidence on an Issue without closing it.
- The **checker** verifies the proof and closes the Issue.
- The **owner** decides acceptance.

## Out-Of-Scope Repair

Repair discovered while capturing or investigating an Issue that falls outside the current Task's boundaries becomes a new, bounded Task in an Objective. It is never carried out as work performed inside the Issue itself — an Issue records follow-up; it does not repair itself.
