---
id: T004
title: Keep active guidance aligned with the simplified card
objective: O012
planned_by: {role: planner, session: task-card-simplification-20260921}
status: planned
depends_on: [{task: T005, requires: clear}]
owner_validation: {required: false}
---

# T004: Keep active guidance aligned with the simplified card

## Outcome

Current architecture and repository guidance describe the implemented Task-card
vocabulary and its policy boundary, while historical evidence remains intact and
cannot be mistaken for instructions to restore the retired presentation.

## User Check

No additional product choice is required. Compare the reconciled wording with
O012's accepted vocabulary and confirm it does not redefine waiver or exception
as technical `CLEAR`. If the optional local Task Check is skipped, record an
explicit waiver and rely on O012's mandatory Full Objective Check to verify
documentation/code agreement.

## Done When

- `.savepoint/Design.md` records the canonical Task-card outcomes, removal of
  ordinary completion badges, and the unchanged data/gate boundary.
- `AGENTS.md` describes the board's single review-outcome mapping rather than a
  composition that implies the retired completion treatment.
- `.savepoint/Guardrails.md`, `.savepoint/visual-identity.md`, `README.md`, and
  `templates/project-v2/AGENTS.md` are reviewed against the implemented behavior;
  each necessary edit is made and every no-change result is recorded.
- A targeted search across the named active documents finds no instruction to
  render `✓ DONE`, `⚠ DONE`, `BY WAIVER`, `BY EXCEPTION`, `Check (stale)`,
  `Check (unverified)`, or `? CHECKER` on Task cards.
- Active guidance continues to state that waiver and owner-accepted exception
  are not independent technical clearance and do not replace mandatory
  Objective or Release Checks.
- Historical archives, immutable Checks, and Issues are left untouched;
  references there are treated as evidence, not active policy. O900's fixture
  records are reconciled only by T005 for visual coverage.
- Documentation reconciliation, focused tests, `git diff --check`, `make build`,
  and `make test` pass before handoff.

## Context Files

`.savepoint/Design.md`, `.savepoint/Guardrails.md`, `.savepoint/visual-identity.md`, `AGENTS.md`, `README.md`, `templates/project-v2/AGENTS.md`, `internal/board/v2/card.go`, `internal/board/v2/badges.go`.

## Design References

Design sections 1, 4, 7, 8, and 13.

## Guardrails

TPL-02, ARCH-04, TEST-01, TEST-08, STYLE-07, STYLE-08, STYLE-09, STYLE-10.

## Implementation Plan

1. Treat T003's tested renderer as the implementation authority and compare it
   with every named active document before changing prose.
2. Reconcile Design section 8 and any related architecture/status statements
   with the final one-outcome card vocabulary and unchanged gate semantics.
3. Update the AGENTS Codebase Map where its badge-composition description no
   longer matches the board; change other named documents only where a real
   contradiction exists.
4. Search the named active documents for retired Task-card labels and record
   which files changed, which required no edit, and why historical matches were
   intentionally excluded.
5. Verify that active wording cannot be read as making waiver or owner acceptance
   technical `CLEAR`, satisfying an `accepted` dependency, or weakening mandatory
   Objective/Release Checks.
6. Run focused V2 board tests, `git diff --check`, `make build`, and `make test`;
   record exact documentation review and command evidence.

## Boundaries

No production rendering changes, gate-policy changes, archive rewrites, immutable
Check edits, Issue edits, or broad copy refresh. Documentation changes are limited
to statements made inaccurate by T003.

## Technical Verification

Targeted retired-vocabulary searches over the exact Context Files, focused
`go test ./internal/board/v2`, `git diff --check`, `make build`, and `make test`.
The mandatory Full Objective Check applies `agent-skills/references/check-method.md`
to renderer/documentation agreement before O012 can close.

## Technical Evidence

Pending execution: reviewed/changed/no-change document inventory, search output,
focused/full command results, files changed, and limitations.

## Drift Notes

Any documentation conflict that implies a product-policy change rather than the
accepted card presentation is REPLAN REQUIRED and returns to design.
