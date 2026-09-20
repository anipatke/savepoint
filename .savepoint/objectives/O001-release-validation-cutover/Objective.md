---
id: O001
title: 'E50: Release evaluated V2 and retire transitional runtime'
status: planned
release: R006
---
## Migrated from V1

Relocated verbatim from the V1 epic body at `.savepoint/releases/v2/epics/E50-release-validation-cutover/E50-Detail.md` (release `v2`).

## V1 Body (verbatim)

# E50: Release evaluated V2 and retire transitional runtime

## Purpose

Prove the V2 workflow on realistic projects and supported distribution targets, then make V2 the only live runtime without weakening migration recovery or Release accountability.

V1 delivery epic for V2 product Objective O010. The current V1 lifecycle governs this work until the maintainer-controlled repository cutover succeeds.

## What this epic adds

- One fail-closed cutover preflight that combines migration applicability, recoverability, clean V2 loading, and the canonical Release cutover decision.
- V2-only board, doctor, resume, init, and shipped workflow defaults; V1 parsing remains reachable only from explicit migration and frozen historical fixtures.
- Six-platform package and checksum evidence plus CLI help and public documentation that match the shipped V2 behavior.
- Three recorded agent scenarios and two isolated migration trials, including an existing-codebase copy, with limitations stated instead of universal model claims.
- A maintainer-controlled runbook and post-cutover verification path for migrating this repository without using it as a mutable test fixture.

## Components and files

| Module | Purpose |
|--------|---------|
| `internal/data/` | Load and validate V2 state and expose `ResolveReleaseCutover` as the only Release readiness composition. |
| `internal/migrate/` | Prove an appliable, recoverable V1-to-V2 operation and retain the isolated V1 input boundary after cutover. |
| `main.go`, `cmd/` | Route live commands to V2 only and return named migration guidance for legacy input. |
| `internal/board/`, `internal/board/v2/`, `internal/doctor/`, `internal/resume/` | Remove live V1 interpretation while preserving one shared V2 Next, evidence, and diagnostic model. |
| `internal/init/`, `templates/`, `agent-skills/` | Ship and activate only the V2 workflow for new or migrated projects while preserving upgrade provenance. |
| `internal/buildtool/`, `bin/savepoint.js`, `package.json`, `Makefile`, `.github/workflows/` | Build and verify the six supported archives, launchers, and checksums without publishing. |
| `README.md`, `.savepoint/Design.md`, `.savepoint/Guardrails.md`, `AGENTS.md` | Reconcile public behavior, active architecture, policy, and repository routing at the proven cutover. |
| `E50-Validation.md`, `E50-Cutover.md` | Record reproducible scenario/trial evidence and the maintainer-controlled live cutover handoff. |

## Architectural delta

Normal command execution stops dispatching by V1/V2 schema and accepts only a valid V2 project. A V1 project receives a named, read-only migration action; an incomplete operation receives recovery guidance. The explicit migration command keeps its confined V1 readers, frozen fixtures, exact-byte archive mapping, and recoverable publish protocol.

Cutover eligibility is fail-closed. Migration planning must be unambiguous and recoverable, the candidate project must load cleanly as V2, and every declared Release must pass `data.ResolveReleaseCutover`. E50 may compose operational prerequisites around that result, but it must not reproduce Release completion rules.

The package default and live repository workflow become the four V2 phases only after software, distribution, trial, and agent-scenario evidence are recorded. The live repository migration remains an explicit maintainer action; agents prepare and verify it but do not run the Savepoint CLI themselves.

Reference: `.savepoint/releases/v2/v2-Design.md`.

## Boundaries

**In scope:**

- Cutover refusal and recovery behavior, V2-only live command routing, migration-only V1 compatibility, V2 scaffold/workflow activation, six-platform packaging proof, realistic isolated trials, agent evaluation evidence, and this repository's explicit cutover handoff.

**Out of scope:**

- Automatic publication, deployment, tagging, changelog generation, telemetry, hooks, push enforcement, a benchmark program, universal model-quality claims, or a permanent `--v1` runtime.

## Dependencies

- E49-objective-board
- E51-first-class-releases, including a clear independent re-audit before any E50 task enters `in_progress`

## Quality gates

- Cutover refuses ambiguous migration input, incomplete recovery, invalid V2 state, technically unclear Releases, unexcepted material Release Issues, stale or missing Release evidence, and missing exact owner acceptance.
- Release readiness comes only from `data.ResolveReleaseCutover`, which delegates to `ResolveReleaseCompletion`; tests detect any parallel E50 policy.
- V1 source readers are reachable only from explicit migration or frozen fixtures; board, doctor, resume, init, and ordinary startup are V2-only.
- Three named agent scenarios and two isolated trial projects are recorded with session/model when supplied, scope, reads, replans, findings, limitations, and no live-project mutation.
- Linux, macOS, and Windows archives for amd64 and arm64 are built, launched or structurally verified as appropriate, and covered by distribution checksums.
- Migration preview remains write-free; backup, every publish boundary, recovery, edit-conflict retry, unchanged second apply, archive mapping, and reference integrity retain named evidence.
- Public help, README, active Design, Guardrails, AGENTS routing, canonical skills, and shipped templates describe the same V2-only runtime.
- Every implementation handoff runs focused tests plus `make build && make test`; epic closeout requires a fresh independent V1 audit.

## Open decisions

None. The live repository apply is intentionally reserved for an explicit maintainer action after E51 re-audit and all pre-cutover E50 evidence pass.
