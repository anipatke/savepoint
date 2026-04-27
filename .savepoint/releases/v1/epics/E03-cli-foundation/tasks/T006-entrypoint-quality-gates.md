---
id: E03-cli-foundation/T006-entrypoint-quality-gates
status: done
objective: "Connect the process entrypoint to the CLI runner and complete the E03 quality-gate closeout."
depends_on:
  - E03-cli-foundation/T005-cli-runner-dispatch
---

# T006: Entrypoint Quality Gates

## Scope

Finish the CLI foundation by keeping process globals in `src/cli.ts`, running the required quality gates, and creating the E03 audit snapshot handoff.

## Implementation Plan

- [x] Read `.savepoint/router.md`, this epic `Design.md`, this task file, and the directly touched entrypoint/test files.
- [x] Update `src/cli.ts` so process argv, streams, environment, platform, and exit code handling stay at the entrypoint boundary.
- [x] Add or update entrypoint-level tests only if existing coverage does not prove the process boundary behavior.
- [x] Run the focused CLI test files for the touched behavior.
- [x] Run `npm run build`.
- [x] Run `npm run typecheck`.
- [x] Run `npm run lint`.
- [x] Run `npm run format:check`.
- [x] Run `npm test`.
- [x] Create `.savepoint/audit/E03-cli-foundation/snapshot.md` with the implemented file list and quality-gate results for audit handoff.

## Acceptance Criteria

- Process globals are isolated to `src/cli.ts`.
- The full E03 quality-gate command list has been run and recorded.
- The audit snapshot exists for the next `audit-pending` state.
