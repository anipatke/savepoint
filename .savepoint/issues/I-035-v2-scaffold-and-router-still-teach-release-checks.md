---
id: I-035
title: V2 scaffold and live router still teach Release Checks
type: drift
status: resolved
source:
  kind: check
  check: C-911
  actor: {role: checker, session: o013-full-check-20260923}
  at: '2026-09-23T09:45:00Z'
tasks: [T-008]
checks: [C-911]
guardrail_ids: [TPL-02]
severity: medium
resolution:
  disposition: accepted
  actor: {role: owner, session: user}
  at: '2026-09-23T09:47:00Z'
  reason: >-
    Owner directed closure of the repair without a fresh-session recheck. The repair was made by the same session that wrote C-911, so this is owner acceptance, not a technical verified closure.
history:
  - at: '2026-09-23T09:45:00Z'
    actor: {role: checker, session: o013-full-check-20260923}
    kind: observed
    check: C-911
    note: >-
      The shipped V2 scaffold's .savepoint/Design.md, Guardrails.md, and
      router.md, and this repository's live .savepoint/router.md, still
      describe the optional context as a Release and its integration gate as a
      Release Check. T-008's Context Files omitted these files.
  - at: '2026-09-23T09:47:00Z'
    actor: {role: executor, session: o013-full-check-20260923}
    kind: repair_attempted
    note: >-
      At the owner's direction, reworded templates/project-v2/.savepoint/Design.md (Goal Check paragraph), Guardrails.md (TEST-04 and the check-method sentence), router.md (check state and Goal Check paragraph), and the live .savepoint/router.md (check row, lifecycle Goal Check, and Goals paragraph). Search of templates/project-v2 and .savepoint/router.md now finds Release only as the R-### / Release.md storage boundary. internal/init tests, make build, and make test-full pass (2026-09-23, go1.26.2 linux/amd64).
---

# I-035: V2 scaffold and live router still teach Release Checks

## Violated Requirement

O-013 Success Condition 5: "Active architecture, repository guidance, V2
skills, and shipped V2 templates describe Goals while documenting the
compatibility boundary where internal `Release` names remain." Guardrail
TPL-02: shipped guidance must describe current behavior. The board, resume,
skills, and AGENTS.md now call it a Goal Check, so the scaffold contradicts
them.

## Scenario

`savepoint init` writes `templates/project-v2/.savepoint/` into a new V2
project. Those files, plus this repository's own router, still say:

- `templates/project-v2/.savepoint/Design.md:55` — "A Release Check is
  mandatory whenever a Release exists".
- `templates/project-v2/.savepoint/Guardrails.md:47` (TEST-04) — "the
  mandatory Objective or Release Check".
- `templates/project-v2/.savepoint/router.md:40` — "the mandatory
  Objective/Release scope"; `:57` — "A Release Check is mandatory whenever a
  Release exists".
- `.savepoint/router.md:31` — "the mandatory Objective/Release scope";
  `:45` — "Release Check is mandatory whenever a Release exists"; `:48` —
  "Releases are optional delivery boundaries".

Expected: Goal / Goal Check wording, with Release named only as the storage
compatibility boundary. Actual: Release is the public term in all six
places.

T-008's evidence says its search found no remaining "Release Check" or
"optional Release" instructions. That search covered only its listed Context
Files, which did not include these four files.

## Proof Needed

- Each passage above uses Goal / Goal Check public wording.
- A search of `templates/project-v2/` and `.savepoint/router.md` for
  `Release Check`, `a Release exists`, and `Objective/Release` returns nothing.
- `internal/init` scaffold tests still pass. If any test pins the old
  wording, update it.
