---
id: E03-cli-simplify/T001-strip-args
status: planned
objective: "Remove init, audit, doctor from args.ts; keep only board + version + help"
depends_on: [E02-domain-phase-model/T004-simplify-router-domain]
---

# T001: Strip Args

## Acceptance Criteria

- `args.ts` defines only `board` as an allowed command.
- `ParseResult` union has no init, audit, or doctor variants.
- `--version` and `--help` still work.
- Unknown commands return appropriate error.
- All args tests pass.

## Implementation Plan

- [ ] Read existing `args.ts` and `test/cli/args.test.ts`.
- [ ] Remove `InitOptions`, `AuditOptions` interfaces.
- [ ] Remove init/audit parse result variants from `ParseResult`.
- [ ] Remove `parseInitArgs`, `parseAuditArgs`.
- [ ] Simplify `parseArgs` to handle only `board`, `--help`, `--version`.
- [ ] Update `test/cli/args.test.ts`.
- [ ] Run `npm test`.
