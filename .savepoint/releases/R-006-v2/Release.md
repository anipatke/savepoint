---
id: R-006
title: Savepoint V2 — Simple workflow, trustworthy completion
status: done
legacy_fields:
    release: v2
    type: release-prd
---
## Outcome

The V1 release promise is preserved in the Legacy Source section below. Source: `.savepoint/releases/v2/v2-PRD.md`.

## Why

This Release carries the V1 delivery boundary forward with a stable V2 identity.

## Success Conditions

- Converted Objectives sourced from this V1 release reference this Release identity.
- Historical completion, when present, remains typed archive evidence rather than a V2 Check.

## Boundaries

Release membership is derived from Objective records; Tasks remain owned by Objectives.

## Legacy Source (verbatim)

Source path: `.savepoint/releases/v2/v2-PRD.md`

SHA-256: `71737d6915e22ab96ce628550ab390fd1208e2bfdc451bfc344a1bb6b83480a1`

Original frontmatter:

```yaml
type: release-prd
status: planned
release: v2
```

````markdown


# Savepoint V2 — Simple workflow, trustworthy completion

## Purpose

An AI-assisted solo builder can understand what to do next, hand a bounded Task to an executor, obtain an independent technical Check, and accept outcomes that need owner judgment. Existing work survives the transition without forcing permanent V1/V2 compatibility.

## Inputs and authority

- Confirmed product decisions: `.savepoint/Savepoint V2 Refactor Prompt.md`, including its opening decision register.
- Proposed architecture and schemas: [v2-Design.md](v2-Design.md).
- Current implementation and build authority: `.savepoint/Design.md`, `.savepoint/Guardrails.md`, and the current V1 phase skills.

This release is managed as V1 `Release → Epic → Task`. Each delivery epic below implements one V2 product Objective. These are not duplicate V2 work records. Current tasks remain `planned` until implementation starts; only the owner marks V1 tasks `done`.

## First-release scope

- Idea, Design, mandatory Objectives, detailed Tasks, independent Checks, and one Issues collection with Defect as a type.
- Detailed planning one Objective at a time; verifiable technical and research outcomes; targeted supporting reads; explicit replanning preserving partial work.
- Conditional owner acceptance, dependency approval prerequisites, current technical clearance, Objective integration gates, and recorded owner exceptions.
- Optional first-class Releases with stable `R###` identities, derived Objective membership, Release-scoped Checks, historical completion references, and an owner-accepted completion gate. A Release is a delivery boundary, not a publishing or deployment workflow.
- Globally stable Task and Issue IDs; reviewed-scope/revision evidence with agent-assessed freshness.
- Safe conversion of active work and unresolved Issues, intact historical archives, reference mappings, and recoverable migration.
- Four public skills, concise shared checking method, config-based verification commands, optional project-specific procedures.
- Fresh/existing-codebase onboarding, safe asset upgrades, unified Issues view, prominent Next action, read-only resume, practical doctor evidence.
- Preserved Atari-Noir identity, narrow-terminal behavior, non-TTY output, local files, no telemetry, and six-platform binary distribution.

## Deferred

Automatic evidence capture and freshness detection, hooks, push enforcement, mandatory model invocation, broad multi-model benchmarking, permanent dual-mode runtime, team/cloud features, and a database. Technical verification still runs through existing configured checks and agents; deferring capture does not defer verification.

## Success conditions

1. A new user can describe the four-step rhythm and find the next action from the board or resume without learning internal router vocabulary.
2. A planned Task can be executed from its explicit context with only recorded targeted supporting reads; material plan gaps return `REPLAN REQUIRED` and preserve partial work.
3. No Savepoint completion control treats executor completion, missing/unknown/stale clearance, or pending required owner acceptance as clean completion.
4. Task dependencies default to technical clearance; explicit approval prerequisites wait for the owner. Dependent Objectives wait for integration clearance; unrelated Objectives do not.
5. Fresh checker sessions verify meaningful failures, preserve Issue identity across rechecks, and require proof before closure.
6. Migration preview is write-free, interruption at every publish boundary is recoverable, a second unchanged run is a no-op, and every Release PRD, user-authored source, and reference has an accountable live or archive destination.
7. E50 cutover is refused until migration is unambiguous and recoverable and the canonical Release completion decision permits every declared Release; technical uncertainty, material Release Issues, and missing owner acceptance remain explicit blockers.
8. Software regression evidence and the three agent scenarios in the design are recorded before release. No live user project is used as a mutable test fixture.

## Ordered delivery

| V1 epic / V2 Objective | Outcome | Dependencies |
|---|---|---|
| E41 / O-001 — Migration source fixtures | Representative V1 work and history have frozen, interpretable preservation evidence. | None |
| E42 / O-002 — Project schema and identity | V2 projects load with stable IDs, strict diagnostics, and preserved author content. | E41 |
| E43 / O-003 — Task and Check gates | Completion and dependency decisions reflect technical clearance and conditional acceptance. | E42 |
| E44 / O-004 — Issues and integration | Follow-up converges through one Issue model and Objective Checks. | E43 |
| E45 / O-005 — Safe migration | Active V1 work converts predictably with archives, references, and recovery. | E44 |
| E46 / O-006 — Agent workflow | Four public skills deliver bounded planning, execution, and independent verification. | E44 |
| E47 / O-007 — Onboarding and upgrades | Fresh and existing projects receive a coherent V2 workflow without silent overwrites. | E45, E46 |
| E48 / O-008 — Next and resume | Users recover current work, evidence, and next action from one shared interpretation. | E44 |
| E49 / O-009 — Objective board | The TUI presents Objectives, Tasks, Checks, and Issues clearly with enforced actions. | E47, E48 |
| E51 / O-011 — First-class releases | Release remains a navigable delivery boundary through V2 migration, board, resume, and cutover. | E49, E45 |
| E50 / O-010 — Release validation and cutover | Evaluated, packaged V2 replaces transitional V1 runtime support safely. | E49, E51 |

Default execution order is E41 through E51, with E51 completing before the E50 cutover. Dependencies permit independent design work, but the current V1 audit handoff still applies between build epics. Each epic has an independently scoped integration/audit outcome. Only E41 has detailed Task plans; design and task breakdown for later epics happen as earlier implementation settles.

## Current readiness

Product choices are confirmed. E51's first-class Release model, migration mapping, recovery proof, board/resume parity, and project-level cutover composition are implemented and tested on temporary fixture and repository copies. `data.ResolveReleaseCompletion` remains the canonical per-Release decision and `data.ResolveReleaseCutover` only composes those decisions for E50; neither claims publication, deployment, tagging, or changelog behavior. The live repository is still V1 and has not been migrated. E50 remains responsible for the independent audit, packaging evidence, final V1-reader retirement, and explicit maintainer cutover.
````
