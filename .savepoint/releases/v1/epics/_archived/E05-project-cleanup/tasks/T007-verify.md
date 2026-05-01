---
id: E05-project-cleanup/T007-verify
status: planned
objective: "Run full verification: build, typecheck, lint, test"
depends_on: [E05-project-cleanup/T006-clean-package-json]
---

# T007: Verify

## Acceptance Criteria

- `npm run build` passes.
- `npm run typecheck` passes.
- `npm run lint` passes.
- `npm test` passes.
- No errors, no warnings, no dead code.

## Implementation Plan

- [ ] Run `npm run build`.
- [ ] Run `npm run typecheck`.
- [ ] Run `npm run lint`.
- [ ] Run `npm test`.
- [ ] Fix any failures.
- [ ] Repeat until all gates pass.
