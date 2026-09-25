---
type: epic-design
status: active
---

# Epic E09: doctor-command

## Purpose

Implement `savepoint doctor`, the diagnostic command for corrupted or inconsistent Savepoint projects. Doctor reports problems and suggested repairs, but does not make destructive changes automatically.

## What this epic adds

- Integrity checks for config, router state, release/epic/task structure, YAML/frontmatter, task acceptance criteria, dependency references, cycles, duplicate task IDs, audit state, and orphaned tasks.
- Ad-hoc quality gate runner.
- Human-readable diagnostic output.
- Machine-testable diagnostic objects.
- Suggested repair messages for common failure modes.
- Layer-2 AI review prompt printing hook for future expansion, without implementing v0.2 semantic review.

## Components and files

Expected files introduced or extended by this epic:

| Path                          | Purpose                       |
| ----------------------------- | ----------------------------- |
| `src/commands/doctor.ts`      | Doctor command entrypoint.    |
| `src/doctor/checks/*.ts`      | Individual integrity checks.  |
| `src/doctor/report.ts`        | Diagnostic report formatting. |
| `src/doctor/repairs.ts`       | Suggested repair text.        |
| `src/doctor/quality-gates.ts` | Ad-hoc quality gate reuse.    |
| `test/doctor/**/*.test.ts`    | Diagnostic behavior tests.    |

## Architectural delta

Before this epic, validation exists inside specific workflows. After this epic, users have a central command for diagnosing project health outside normal board/audit flows.

Doctor should reuse the data model and audit quality-gate runner instead of introducing parallel logic.

## Boundaries

In scope:

- Detect corrupt YAML/frontmatter.
- Detect missing config.
- Detect missing dependencies and dependency cycles.
- Detect duplicate task IDs.
- Detect missing or malformed task acceptance criteria after the board workflow cleanup makes them required.
- Detect tasks in nonexistent epics.
- Detect audit proposals without matching audit-pending state.
- Run configured quality gates on demand.

Out of scope:

- Auto-moving orphaned tasks without user action.
- Repairing files destructively.
- Launching the TUI.
- Calling AI APIs.

## Quality gates

- Each failure mode listed in the project Design should have at least one test.
- Reports should include file paths and actionable messages.
- Exit codes should distinguish clean, diagnosed problems, and internal command failure.

## Design constraints

- Doctor diagnoses and proposes; it does not silently mutate.
- Keep diagnostics deterministic and suitable for CI output.
- Reuse existing validation functions wherever possible.
