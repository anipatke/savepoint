---
name: savepoint-check
description: Runs an independent, fresh-session Check on an explicitly requested Task, a mandatory Objective, or a mandatory Release when router state is check, applying the shared check method to write an immutable Check record and any Issues, and closing work only under the recorded conditional-acceptance rules.
---

# Savepoint Skill: Check

## Purpose

Turn recorded evidence into an independent, immutable verdict. A Task Check is an optional local review; the Full Objective Check is mandatory as the higher-level integration gate, and a Release Check is mandatory whenever a Release exists. This is the only role that can write those verdicts or close Issues, so its authority is bounded on both sides: it can assess a requested Task Check without turning that local result into Objective or Release completion, and it can never manufacture clearance by repairing the implementation, rewriting acceptance criteria to match what was built, or updating Design as a form of remediation. Correction always goes back to the planner or executor; this skill verifies the repair afterward, in a later run, and never edits its own prior record to do so.

## Trigger

Use this skill when router `state` is `check` for an explicitly requested
Task Check or for a mandatory Objective/Release Check.

This session must be fresh: independent from the executor's conversation that built the work under review. The same model is allowed; the same session is not. Model names are optional — nothing here authenticates identity, it only refuses to treat the executor's own session, or the owner's or planner's self-report, as an independent Check.

## Read

- `.savepoint/router.md`
- The requested Task, Objective, or Release under review; for an Objective
  Check, every Task it owns; for a Release Check, every member Objective and
  its owned Tasks
- Applicable policy: the `.savepoint/Guardrails.md` rule IDs the scope names, when the project has that file
- Prior Checks and Issues linked to the scope
- The scoped source and test files the work actually changed
- `agent-skills/references/check-method.md` in full

Load `agent-skills/references/check-method.md` completely and apply it as written. It owns scope locks, coverage matrices, the adversarial pass, materiality, and re-audit convergence; this skill does not restate that method.

## Workflow

1. Confirm the session is fresh. If this session built the work under review, state that limitation; do not proceed as an independent Check unless the user explicitly asks to continue anyway.
2. Confirm the scope: a Task Check evaluates one Task's outcome and evidence (the Task Check itself is optional); an Objective Check does everything a Task Check does, plus integration across the Objective's owned Tasks and reconciliation against Design; a Release Check uses `scope.kind: release` to evaluate integration across all member Objectives.
3. Apply `agent-skills/references/check-method.md` in full at the matching evidence mode — Quick for a requested Task Check, Full for the mandatory Objective Check or Release Check.
4. Decide the result. Write one new, immutable Check record — never edit a prior one. A rerun gets a new `C###` and names the run it replaces in `supersedes`.
5. On `NEEDS WORK`: record the Issues found, and hand remediation back to the executor or planner rather than repairing anything here. The executor resumes at `stage: build` inside the same Task.
6. On `CLEAR`: this alone does not close a Task or Objective. Apply the closure rules below to record whether the owner may complete the Task or accept the Objective/Release outcome.
7. Record advisory observations, including `STYLE` guardrail findings, as non-blocking; do not let them change the result.
8. Stop. Do not repair implementation, rewrite acceptance criteria, or update Design as part of this run.

## Write Boundary

This skill may write: the Check record, Issues, evaluation metadata, and
authorized closure evidence when the closure rules below allow it. It may
record that a Task or Objective is ready for owner closure and whether the
owner has accepted a Release Check, but it never sets Task/Objective status to
`done`; it never records owner acceptance on the owner's behalf.

It must never: repair implementation, edit acceptance criteria to match a result, or update Design as a form of remediation. A correction is always routed back to the planner or executor, and a later Check — a new record, never an edit to this one — verifies that the repair actually landed.

## Check Artifact Template

Write each Check record with this structure:

```yaml
id: C###
scope: {kind: task|objective|release, id: T###, O###, or R###}
result: CLEAR|NEEDS WORK
checked_by: {role: checker, session: review-001}
executed_session: build-001
checked_at: '2026-09-19T00:00:00Z'
reviewed:
  base_commit: optional-base-sha
  head_commit: optional-head-sha
  files: []
  dependencies: []
issues: []
supersedes: null
```

`reviewed.files` and `reviewed.dependencies` are populated with real path or content-hash entries, or explicit absent entries — never left implicit. `issues` lists the `I###` references this run opened. `supersedes` names the prior `C###` this run replaces, or stays empty on a first run. The record body carries outcome coverage, test and command results, negative and boundary probes, applicable Guardrails, owner validation still needed, and nonblocking observations.

Each run writes a new record with a new `C###`. A recheck never edits the superseded record; it sets its own `supersedes` and leaves the prior run intact as history.

## Closure Rules

- A checker may complete a technical Task's optional Check — one with no `owner_validation.required` — once its clearance is current and no unexcepted material blocker remains; only the owner may set the Task's `status: done`.
- A Task with no requested Task Check may be owner-closed only when its implementation evidence is complete and an explicit Task-check waiver names the Task, reason, actor, and time. The waiver skips only the optional local Check: it is not technical `CLEAR`, does not satisfy a dependency that explicitly requires `clear`, and does not waive any acceptance criterion, guardrail, Objective Check, or Release Check.
- A Task declaring `owner_validation.required` additionally needs the owner's recorded acceptance naming this same current Check; acceptance naming a Check a later run has superseded does not count.
- An Objective closes only after every Task it owns is done and the mandatory Objective integration Check is current, with the same conditional owner-acceptance rule applied at the Objective level. The Full Objective Check reviews every owned Task, including waived Task Checks. An unfinished owned Task is never excused by an Objective-level exception — cross-Task repair goes back through Tasks, and no Objective Check ever closes a Task directly.
- A Release Check is mandatory whenever a Release exists. It reviews cross-Objective integration for `R###`, reuses ordinary Issues for material findings, and does not invent a parallel release audit. Release `done` requires at least one member Objective, every member Objective complete, current CLEAR integration evidence, resolved or explicitly excepted material Issues, and the owner's acceptance of that exact current Check. The checker never supplies that acceptance, and Release `done` does not mean published or deployed.
- A record lacking sufficient scope or evidence cannot support completion. Stale or unknown freshness blocks normal completion; it is never waived through by re-asserting "current" without a fresh assessment.
- A recorded owner exception can grant completion despite an unmet requirement, but it is reported as completion by exception, never as a `CLEAR` result or as current clearance, and it applies only to the Check it names.

## Issue Capture

Enter Issue capture as an entry from this workflow when a Check finds
something that blocks the verdict. A `NEEDS WORK` Check records the Issues it
finds; see `agent-skills/references/issue-capture.md` for the artifact
template and rules. A failed Release Check creates or reuses ordinary Issues;
it does not create release-only findings. This skill is the one role that may close an Issue, after
verifying its proof.

## Rules

- This session must be independent from the executor's conversation under review; state the limitation plainly when it is not, and do not call the result independent unless the user explicitly says to continue anyway.
- Load and apply `agent-skills/references/check-method.md` in full; do not restate its scope-lock, coverage-matrix, adversarial-pass, materiality, or convergence mechanics here.
- Write only the Check record, Issues, evaluation metadata, and authorized closure. Never repair implementation, edit acceptance criteria, or update Design as remediation — route corrections back to the planner or executor.
- Every Check run is a new immutable `C###` record; a recheck sets `supersedes` and never edits a prior run.
- Apply Quick evidence only for a requested Task Check. Apply Full evidence for the mandatory Objective Check and Release Check; an Objective Check additionally covers cross-Task integration and Design reconciliation, and a Task-only Check never substitutes for either mandatory integration gate.
- Apply the Release Check scope when `scope.kind: release`: inspect all member Objectives and their cross-Objective integration, then reuse ordinary Issues and hand owner acceptance back to the owner.
- A `NEEDS WORK` result records Issues and hands remediation to the executor or planner; the executor resumes at `stage: build` inside the same Task.
- Apply the closure rules above exactly; do not use stale, unknown, or missing clearance to satisfy an invoked Task Check or a mandatory Objective/Release Check, and do not treat a Task-check waiver or other exception as a `CLEAR` result.
- Treat advisory observations, including `STYLE` guardrail rules, as non-blocking; record them, but never let them change the result on their own.
- Use `state` only for router phase, Task `status` only for Task lifecycle, and `stage` only when the Task is `in_progress`.
