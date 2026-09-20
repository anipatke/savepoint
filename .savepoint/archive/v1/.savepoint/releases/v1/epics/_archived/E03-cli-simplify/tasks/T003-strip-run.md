---
id: E03-cli-simplify/T003-strip-run
status: planned
objective: "Remove init, audit, doctor dispatch from run.ts; keep only board"
depends_on: [E03-cli-simplify/T001-strip-args]
---

# T003: Strip Run

## Acceptance Criteria

- `run.ts` imports only `runBoard` from commands.
- `runCli` dispatches only `board`, `--version`, and `--help`.
- No dead imports (init, audit, doctor).
- All runner tests pass.

## Implementation Plan

- [ ] Read existing `run.ts` and `test/cli/run.test.ts`.
- [ ] Remove imports of `runInit`, `runAudit`, `runDoctor`.
- [ ] Remove init/audit/doctor dispatch branches.
- [ ] Update `test/cli/run.test.ts`.
- [ ] Run `npm test`.
