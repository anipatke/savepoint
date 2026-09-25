---
type: epic-design
status: audited
---

# E43: Make Task completion trustworthy

## Purpose

Decide Task progress and completion from recorded evidence instead of a status field, so no control treats executor self-report, missing clearance, stale clearance, or a pending owner decision as clean completion.

This is the V1 delivery epic for V2 product Objective O003. The current V1 lifecycle remains authoritative while it is built: E43 makes the V2 decision canonical in data, it does not switch any live workflow onto it.

## What this epic adds

- Immutable Check records under `.savepoint/checks/`, one file per completed evaluation, with scope, result, actor provenance, reviewed basis, and `supersedes` linkage for reruns.
- A recorded freshness assessment on Task and Objective records: clearance is `current` only when an assessment says so; an absent assessment is `unknown`, never current.
- Explicit actor authority — planner, executor, checker, owner — attached to Checks, acceptance, and exceptions, so a recorded outcome names who produced it.
- Conditional owner acceptance: technical Tasks close on current clearance; Tasks declaring `owner_validation.required` also need the owner's acceptance of a named Check.
- Task dependency gates that decide satisfaction from clearance and acceptance, turning E42's structural `{task, requires: clear|accepted}` graph into an answer.
- A replan flag that pauses execution and preserves partial work without inventing a fourth Task status.
- Recorded owner exceptions that name requirement IDs, reason, owner provenance, time, and the affected Check — and are reported as completion by exception rather than as a CLEAR result.
- One canonical gate API in `internal/data` that every V2 consumer will call, plus named diagnostics for malformed or inconsistent evidence surfaced through doctor.

## Components and files

| Module | Purpose |
|--------|---------|
| `internal/data/check_v2.go` (new) | Define and strictly decode the shared Check record: ID, scope, result, actors, timestamps, reviewed basis, issue references, `supersedes`. |
| `internal/data/evidence_v2.go` (new) | Define and decode the evidence blocks shared by Task and Objective records: actor references, freshness, owner validation, exception, replan. |
| `internal/data/gate_v2.go` (new) | Evaluate clearance state, Task start/advance/completion decisions, completion authority, and typed blocker reasons. |
| `internal/data/task_v2.go` | Extend `TaskV2` with the shared evidence fields; keep strict decoding and V1 isolation unchanged. |
| `internal/data/objective_v2.go` | Decode the same evidence fields on Objectives so E44 evaluates them without a second schema. |
| `internal/data/discover.go` | Add confined `.savepoint/checks/` discovery using the existing V2 path confinement and path/ID rules. |
| `internal/data/project.go` | Extend `V2Index` with Checks keyed by C ID, a scope-target index, and next-ID allocation over active records. |
| `internal/data/dependency.go` | Evaluate Task dependency satisfaction for `requires: clear` and `requires: accepted` over the index. |
| `internal/data/write.go` | Create Check records create-only, and patch nested Task evidence fields through the existing preserving write path. |
| `internal/data/errors.go` | Add the named Check/evidence diagnostics that doctor reports. |
| `internal/doctor/checks.go`, `internal/doctor/repairs.go` | Report the new diagnostics by stable name with manual repair guidance; no new parsing or policy. |
| `internal/data/*_test.go`, `internal/doctor/*_test.go` | Cover Check decoding, supersedes chains, clearance states, authority, dependency gates, exceptions, replan, external-edit detection, and diagnostic integration on temporary projects. |

## Architectural delta

E42 established structure; E43 adds evaluation on top of it. `LoadV2Index` gains a third record family — Checks discovered from `.savepoint/checks/`, keyed by global `C###` ID, with a scope index from each target Task or Objective ID to its Checks in recorded order. Structural problems keep failing closed at load. Evaluation never returns an error for a recoverable state: missing, stale, or unknown clearance is a decision with a named reason, not a load failure.

Check records are immutable. A rerun is a new `C###` with `supersedes` naming the previous Check for the same scope; the index validates that the target exists, shares the scope, and forms no cycle or fork. Writes are create-only and refuse an existing path, so a corrected outcome must be a new record. New IDs are allocated as the next unused ID across active records; E45 extends allocation with migration reservations, and archived identities are never reused.

The Check schema is shared by both scopes via `scope.kind: task|objective`, and `issues:` decodes as `I###` identity references only. E43 evaluates Task scope. Objective integration gating and the Issue records themselves belong to E44, which consumes this schema rather than adding one.

Task and Objective frontmatter carry the same evidence block: `last_check: C###`; `freshness: {state: current|stale|unknown, check: C###, assessed_by: {role, session}, assessed_at, basis}`; `owner_validation: {required: bool, accepted_check: C###}`; an optional `exception: {requirements: [ID], reason, owner, recorded_at, check}`; and an optional `replan: {reason, recorded_by: {role, session}, recorded_at}`. There is no mutable `checked: true` flag. Unknown or malformed evidence values are named diagnostics at load, never healed into a completion-capable state; absent optional evidence is legal and evaluates to `unknown`.

Clearance resolves to exactly one of `current`, `needs_work`, `stale`, `unknown`, or `missing` by deterministic rules, with no materiality inference: clearance is `current` only when the record's latest Check is CLEAR and a freshness assessment names that same Check with `state: current`. A freshness assessment naming any other Check is `stale`; no assessment is `unknown`; a NEEDS WORK latest Check is `needs_work`; no Check at all is `missing`. Owner acceptance binds to the Check it names, so a superseding Check invalidates acceptance until the owner re-accepts.

Gate decisions are the single place completion eligibility is decided. Start requires a ready plan, no replan flag, and satisfied dependencies. Completion requires current clearance, plus owner acceptance when the Task declares it — technical Tasks close under checker authority, owner-validated Tasks under owner acceptance or checker finalization of an acceptance that still applies. Each decision returns whether it is allowed, which actor may act, and typed blockers naming the unmet requirement, so consumers render reasons instead of re-deriving rules. A decision allowed only by a recorded exception is returned as allowed-by-exception with the exception attached, never as clearance.

Dependency satisfaction reads the same evaluation: `requires: clear` needs the dependency Task done with current clearance; `requires: accepted` additionally needs recorded owner acceptance of that Task's current Check. Both re-resolve evidence from files at decision time rather than trusting a cached judgement.

Data detects inconsistent direct edits — a `done` Task with no current clearance, evidence naming a missing Check, an acceptance pointing at a superseded Check — and reports them as named diagnostics. It does not claim to authenticate actors or prevent external writes; that limit is stated, not engineered around.

Reference: `.savepoint/releases/v2/v2-Design.md`, sections 4 and 5.

## Boundaries

**In scope:**

- Check record schema, confined discovery, index, supersedes validation, and create-only writes.
- Shared evidence decoding on Task and Objective records, and preserving managed writes of Task evidence fields.
- Clearance evaluation, Task start/advance/completion gates, completion authority, dependency satisfaction, exceptions, and the replan flag.
- Named diagnostics for malformed evidence and inconsistent external edits, reported through doctor.

**Out of scope:**

- Issue records, issue reconciliation, and Objective integration gating; E44 owns them and consumes this Check schema.
- Automatic freshness detection, repository hash scanning, Git hooks, actor authentication, or checker-run repair. Freshness is an agent-recorded assessment.
- Rewiring live consumers. `internal/board/transitions.go`, router, and Next still decide from V1 `Task` records; E48 and E49 adopt the gate API when their consumers move to V2. E43's obligation is that the decision exists once, in data, so those epics do not re-implement it. This narrows the component list in the previous E43 stub, which named `internal/board/transitions.go` as a target before E42 settled the loader boundary.
- Activating the V2 policy change to DATA-05. Owner-only completion remains the active V1 rule until cutover in E50; E43 may model conditional authority but must not change how agents treat V1 tasks.
- Migration of V1 audit runs, findings, or defects into Checks and Issues; E45 owns it. E43 fabricates no Check for historical work.
- V2 skills and shipped guidance describing the checker workflow; E46 owns them.

## Dependencies

- E42-project-schema-identity supplies the schema dispatch, strict record decoding, identity-keyed index, path confinement, and preserving write path this epic extends.

## Quality gates

- Clearance tests produce an explicit result for each of current, needs_work, stale, unknown, and missing, including a freshness assessment that names a superseded Check.
- Completion tests distinguish a technical Task closing under checker authority from an owner-validated Task blocked until acceptance, and prove acceptance of a superseded Check does not close it.
- Dependency tests cover `requires: clear` satisfied and unsatisfied, `requires: accepted` waiting on the owner, and a dependency that is done but not currently cleared.
- Check record tests cover strict decoding, duplicate and malformed C IDs, scope-target mismatch, missing targets, valid and forked/cyclic supersedes chains, and refusal to overwrite an existing Check file.
- Exception tests prove an exception records requirement IDs, reason, owner, time, and affected Check, is reported as completion by exception, and never reports a CLEAR result or carries to a superseding Check.
- Replan tests prove the flag blocks start/advance, preserves status, stage, and partial work, and clears only through an explicit write.
- External-edit tests prove a hand-edited `done` Task without clearance and evidence referencing a missing Check produce named diagnostics rather than silent acceptance or a rewritten status.
- Write tests prove Task evidence patches preserve unknown frontmatter fields, authored bodies, and line-ending form, and that a no-op write leaves the file untouched.
- Doctor tests prove every new diagnostic reports under a stable name with a manual repair suggestion, and that doctor writes no project files.
- Implementation handoff requires focused `internal/data` and `internal/doctor` tests plus `make build && make test`, with named cases and outcomes recorded in the Tasks.
- Epic closeout requires a fresh independent V1 epic audit; this planning session is not that audit. Current Guardrails apply, with STYLE advisory.

## Open decisions

None. Exact Go identifiers may be refined during task breakdown, but the Check schema, evidence field names, clearance rules, authority model, and consumer boundary above are fixed.
