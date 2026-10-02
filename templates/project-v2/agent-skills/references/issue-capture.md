---
type: issue-capture-reference
triggerable: false
---

# Shared Savepoint Issue Capture

This non-triggerable reference owns Issue capture for Design, Task, and Check; each enters it from its own workflow.

An Issue captures durable follow-up from planning, execution, or Check; it is a record, not a fifth router state. "Defect" stays a word the user says; it maps to `type: defect` and does not resurrect a separate defect record, status, or public phase in V2.

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

A Task Check skipped under an explicit owner waiver is not an Issue; apply AGENTS.md's Verification Policy.

## Search Before Creating

Before allocating an `I-###`, search for the same symptom, the same location, the same violated requirement, or the same linked work. No automatic deduplication is assumed; every capture performs this search.

After creating or renaming an Issue or another identity-bearing record outside `savepoint create-task`, run `savepoint resume` to require strict loading of the complete V2 index before handoff.

## Resolution Dispositions

A resolved Issue records exactly one disposition:

- **verified** — the Issue was repaired, and the repair is proven by a Check that recorded `CLEAR`.
- **accepted** — an explicit owner decision, not a `CLEAR` Check or independent proof; visual inspection may inform it. It never waives the mandatory Objective Check. See Role Boundaries for recording and board actions.
- **duplicate** — the same problem as another, canonical Issue. It names that Issue and proves nothing itself.
- **escalated** — repair promoted into an Objective, named in `escalated_to`; it proves nothing itself. See Escalation Retires The Issue.

Reopening a recurring problem reuses the same `I-###` with new, dated evidence rather than allocating a new ID.

## Escalation Retires The Issue

When repair becomes an Objective of its own, immediately set `status: resolved` with `resolution: {disposition: escalated, escalated_to: O-###, actor: {role: planner, ...}, ...}` and append a `kind: escalated` history entry naming that Objective. Do not wait for repair proof or re-verify the retired Issue: the new Objective carries its mandatory Full Check and owner acceptance. A bounded Task inside an existing Objective is not escalation.

## History Is Append-Only

`history` is a frontmatter list of `{at, actor, kind, note, check}` entries: `observed`, `repair_attempted`, `rechecked`, `deferred`, `reopened`, `owner_decision`, or `escalated`. It records the Issue's story in the order it happened. A write that shortens, reorders, or edits a recorded entry is refused — history is appended to, never rewritten.

Deferral is a dated history entry on an open Issue, not a fourth lifecycle state: an Issue stays `open` or `in_progress` while it waits, with the wait recorded as a `deferred` entry naming the reason, never as a new status.

## Role Boundaries

- The **executor** reports repair evidence without independently closing the Issue; it may record an explicit owner-directed `accepted` closure.
- The **checker** verifies Check proof and closes the Issue as `verified`; it may also resolve a confirmed duplicate.
- The **owner** may resolve the Issue as `accepted` with Space or reopen any resolved Issue with Backspace from the board's Issues list. Space resolves an In Progress Issue. The board records the fixed resolution reason, owner actor, time, and append-only history; this does not claim technical `CLEAR`. An agent records an owner decision only when directly instructed.
- The **planner** (`savepoint-design`) closes an Issue as `escalated` at the moment it promotes the repair into a new Objective — see Escalation Retires The Issue above.

## Out-Of-Scope Repair

Default: fix it directly and record repair evidence in the Issue's own
history (`kind: repair_attempted`), leaving the Issue open for a checker or an explicit owner acceptance decision.
Before recording `repair_attempted`, re-run the Issue's Proof Needed scenario
and the neighbouring cases the repair could change, record each result in the
history entry, and mark any platform or case you could not run as unverified.
This is executor evidence only; it never claims `verified`.
After recording a direct repair, advance an Issue-only router selection: use
the one distinct Objective identified by the Issue's linked Tasks and
Objective-scoped Checks, when there is exactly one. Set `objective` to that Objective and clear
`task` and `issue`. If those links identify no Objective or multiple
Objectives, clear `issue` only; do not guess an Objective. Preserve `release:`
and run `savepoint resume` to show the resulting Next line. The Issue remains
open until a checker verifies it or the owner explicitly accepts it.
Escalate to a new Objective only when repair needs an open Design decision or spans multiple Objectives; apply Escalation Retires The Issue immediately. A repair that instead becomes a new, bounded Task
within an existing Objective is not an escalation: the Issue stays open for
the checker as in the default case.
