---
id: E05-project-cleanup/T006-clean-package-json
status: planned
objective: "Remove ink-testing-library and verify no unused dependencies"
depends_on: [E05-project-cleanup/T001-delete-dead-src]
---

# T006: Clean Package JSON

## Acceptance Criteria

- `ink-testing-library` removed from devDependencies.
- All remaining dependencies are used by the simplified codebase.
- `npm install` succeeds.
- Build still works.

## Implementation Plan

- [ ] Read `package.json`.
- [ ] Remove `ink-testing-library` from devDependencies.
- [ ] Verify `@types/react` is still needed (yes, Ink uses React).
- [ ] Run `npm install`.
- [ ] Run `npm run build`.
