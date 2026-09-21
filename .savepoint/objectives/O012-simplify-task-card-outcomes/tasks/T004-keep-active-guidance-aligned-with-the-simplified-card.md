---
id: T004
title: Keep active guidance aligned with the simplified card
objective: O012
planned_by: {role: planner, session: task-card-simplification-20260921}
status: done
depends_on: [{task: T005, requires: clear}]
owner_validation:
    required: false
    accepted_check: ""
check_waiver:
    task: T004
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-09-21T11:44:02Z"
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

**Start gate:** T005's `depends_on: [{task: T005, requires: clear}]` was satisfied
before starting — T005 is `status: done` with a recorded owner `check_waiver`
(`recorded_at: 2026-09-21T11:35:50Z`), and `agent-skills/savepoint-task/SKILL.md`
("Without a Task Check") and `.savepoint/Design.md` Section 5 both state an
owner waiver satisfies a downstream Task dependency that requires `clear`.
Confirmed by reading T005's frontmatter directly (a Context File of this Task).

**Document inventory (all six named active documents reviewed):**

- Changed: `.savepoint/Design.md` — added a **Task-card review outcome (O012)**
  bullet to Section 8 (TUI), immediately after the existing **Layout:** bullet.
  It records the six exact badge strings, the fixed precedence (exception →
  waiver → clearance), that the retired `✓ DONE`/`⚠ DONE`/`BY WAIVER`/
  `BY EXCEPTION` badges no longer appear, the open-card actionable-only rule,
  and that waiver/owner-acceptance remain distinct from independent technical
  clearance and do not replace the mandatory Objective/Release Check —
  `internal/data` still owns every underlying decision (Section 1
  cross-reference); this is presentation only.
- Changed: `AGENTS.md` — the `internal/board/v2/` Codebase Map row's badge-
  mapping parenthetical read `(stage, clearance, gate blockers, exception)`,
  omitting `waiver` even though `taskCheckBadge` (pre-existing) and the
  current `taskReviewOutcomeBadge` (T003) both take a `waived bool` input to
  the same one-mapping surface this row describes. Added `waiver` to the
  tuple: `(stage, clearance, gate blockers, waiver, exception)`. The row
  already read "the one badge mapping" (singular) before this edit — it did
  not describe the retired completion-plus-Check composition, so no other
  wording change was needed there.
- No change (reviewed, no contradiction found): `.savepoint/Guardrails.md` —
  TPL-02, TEST-09, and the Enforcement section describe waiver/exception
  policy in general terms only; none names a specific card badge or implies
  the retired presentation.
- No change (reviewed, no contradiction found): `.savepoint/visual-identity.md`
  — read in full (125 lines); contains no card-badge, review-outcome, or
  completion-badge content to reconcile.
- No change (reviewed, no contradiction found): `README.md` — the "terminal
  board" section describes Task columns at a general level
  ("planned, in-progress, and done work, including the current
  build/test/audit stage") and names no specific badge string.
- No change (reviewed, no contradiction found): `templates/project-v2/AGENTS.md`
  — read in full; its Codebase Map table is the fresh-scaffold empty
  placeholder (header row only), and its Verification Policy / Check sections
  already state waiver and owner acceptance are not technical `CLEAR` and do
  not replace the mandatory Objective/Release Check.

**Retired-vocabulary search:** `grep -n -E` for `✓ DONE|⚠ DONE|BY WAIVER|
BY EXCEPTION|Check \(stale\)|Check \(unverified\)|\? CHECKER|CHECKER` across
all six named documents, run before and after the edits above. Before: zero
matches anywhere. After: one match, in the new `Design.md` bullet itself,
which names the retired strings only to state they "no longer appear" — not
an instruction to render them.

**Per-criterion evidence (Done When):**

- `Design.md` records the canonical outcomes, badge removal, and unchanged
  data/gate boundary: see the new Section 8 bullet above, which states all
  three explicitly (six outcome strings + precedence; retired-badge removal;
  "`internal/data` still owns every underlying ... decision ... This is
  presentation only").
- `AGENTS.md` describes the one review-outcome mapping rather than a
  composition implying the retired treatment: confirmed — the row already
  said "the one badge mapping" (singular) pre-edit, and the only change was
  adding the missing `waiver` input.
- `Guardrails.md`, `visual-identity.md`, `README.md`, and
  `templates/project-v2/AGENTS.md` reviewed against implemented behavior:
  see the document-inventory no-change entries above, each with its reason.
- Targeted search finds no instruction to render a retired badge: see the
  retired-vocabulary search above.
- Active guidance still states waiver/owner-acceptance are not independent
  technical clearance and do not replace mandatory Checks: unchanged in
  `AGENTS.md` Verification Policy, `Guardrails.md` TEST-09/Enforcement,
  `README.md`, and `templates/project-v2/AGENTS.md`; newly stated again in
  the added `Design.md` bullet.
- Historical archives, immutable Checks, and Issues untouched: no file under
  `.savepoint/archive/`, `.savepoint/checks/`, or `.savepoint/issues/` was
  read or written by this Task. O900's fixture records (owned by T005) were
  not touched.
- Documentation reconciliation, focused tests, `git diff --check`,
  `make build`, and `make test` pass: see Commands run below.

**Commands run:**

- `go test ./internal/board/v2/...` — `ok` (cached from T005's unmodified
  package; this Task made no source changes).
- `git diff --check` — exit 0, no whitespace errors.
- `make build` — clean.
- `make test` — full run, exit 0; every package `ok` (`internal/init` 3.7s,
  `internal/migrate` 122.7s, rest cached/fast).

**Files changed:** `.savepoint/Design.md`, `AGENTS.md`. No production code,
no test files, no Guardrails/visual-identity/README/template changes (each
reviewed with a recorded reason above).

**Limitations:**

- No fresh User Check or owner waiver has been requested/recorded for this
  Task specifically; per the router's routing and this Objective's plan, this
  evidence hands off to O012's mandatory Full Objective Check rather than an
  optional local Task Check.
- The AGENTS.md badge-mapping omission of `waiver` predates T003 (confirmed
  via `git log --oneline -- AGENTS.md`, whose most recent entry pre-O012 is
  `e8ac64d`, and `git show fb68054 -- AGENTS.md`, which touched nothing in
  this file) — it was not literally "made inaccurate by T003," but Done-When
  explicitly names updating this row, and the omission is a real, if
  pre-existing, drift against the current `internal/board/v2/badges.go`
  (ARCH-04, TPL-02), so it was corrected here rather than left as a gap this
  Objective's guidance reconciliation would otherwise leave behind.

## Drift Notes

Any documentation conflict that implies a product-policy change rather than the
accepted card presentation is REPLAN REQUIRED and returns to design.
