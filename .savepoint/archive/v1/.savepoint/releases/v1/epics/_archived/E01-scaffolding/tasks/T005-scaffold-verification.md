---
id: E01-scaffolding/T005-scaffold-verification
status: done
objective: "Verify the complete scaffold through the agreed npm scripts and make any final minimal adjustments needed for the first audit baseline."
depends_on: [E01-scaffolding/T004-lint-format-gates]
---

# Task T005: Scaffold Verification

## Implementation Plan

- [x] Run the full scaffold command set: build, typecheck, lint, format check, and tests.
- [x] Inspect the repo for any missing scaffold artifacts or generated files that should stay ignored.
- [x] Make only minimal corrections if the baseline is still inconsistent.
- [x] Mark the scaffold task complete once the verification set passes.

## Scope

Run the full scaffold command set and make small corrective edits if any setup task left the baseline inconsistent.

Expected verification commands:

- `npm run build`
- `npm run typecheck`
- `npm run lint`
- `npm run format:check`
- `npm test`

## Acceptance Criteria

- All expected verification commands pass.
- `package-lock.json` is present and consistent with `package.json`.
- `dist/` remains generated output and is not committed as source.
- The repo is ready for the first epic audit to populate the `AGENTS.md` codebase map.

## Out Of Scope

- Updating `AGENTS.md` codebase map before audit.
- Implementing command behavior.
- Public release packaging polish.
