---
id: I-033
title: Active prose still documents the bare identity format
type: drift
status: resolved
source:
  kind: check
  check: C-910
  actor: {role: checker, session: o018-full-check-20260923}
  at: '2026-09-23T08:40:00Z'
tasks: [T-020]
checks: [C-910]
guardrail_ids: [TPL-02]
severity: medium
resolution:
  disposition: accepted
  actor: {role: owner, session: user}
  at: '2026-09-23T09:05:00Z'
  reason: >-
    Owner accepted the repair (hyphen inserted in the AGENTS.md and Design.md placeholders and in stale bare references in code comments) without a fresh-session recheck. The repair was made by the same session that wrote C-910. This is owner acceptance, not a technical verified closure.
history:
  - at: '2026-09-23T08:40:00Z'
    actor: {role: checker, session: o018-full-check-20260923}
    kind: observed
    check: C-910
    note: >-
      Design.md and root AGENTS.md still give the identity format and record
      paths as bare R###/O###/T###/C###/I### placeholders after the O-018
      rename. Code comments in the Makefile, the pre-commit hook, and the V2
      board still cite renamed records by their old bare IDs.
  - at: '2026-09-23T08:55:00Z'
    actor: {role: executor, session: o018-full-check-20260923}
    kind: repair_attempted
    note: >-
      At the owner's direction ("fix them now"), inserted the hyphen in every
      X### placeholder in AGENTS.md (2) and .savepoint/Design.md (12), and in
      the bare historical references in the Makefile (T-013), the pre-commit hook
      (I-027), and V2 board comments (I-012, I-013, O-012). A word diff shows
      only hyphen insertions. The proof search now returns no hits. This repair
      was made by the same session that wrote C-910, so it needs a fresh-session
      recheck before closure.
---

# I-033: Active prose still documents the bare identity format

## Violated Requirement

O-018 Success Conditions: "This repository's active records, directories,
filenames, cross-references, `router.md`, and active prose (`AGENTS.md`,
`Design.md`, `Guardrails.md`) are renamed" and "a targeted search finds no
unhyphenated active IDs". T-020 Done When names `AGENTS.md` and `Design.md`.
TPL-02 (guidance must describe the behavior the code has).

## Scenario

`grep -nE '[ROTCI]###' .savepoint/Design.md AGENTS.md` returns 14 hits, for
example:

- `AGENTS.md:64` — "Issues live at `.savepoint/issues/I###-slug.md`."
- `AGENTS.md:88` — "Check records are immutable, at `.savepoint/checks/C###-slug.md`."
- `.savepoint/Design.md:57-62` — the layout tree `R###-slug/Release.md`,
  `O###-slug/`, `tasks/T###-slug.md`.
- `.savepoint/Design.md:29, 32, 42, 70, 84, 116, 148, 190` — `escalated_to: O###`,
  `release: R###`, `C###` record, `T###`/`O###` dependency references,
  `R###` identities.

Expected: every format placeholder and path pattern uses `X-###`, matching
`internal/data/identity_v2.go`. Actual: an agent that reads AGENTS.md (every
session) or Design.md and mints an Issue/Check/Task by hand from these
patterns produces a record that strict load now refuses with `ErrV2InvalidID`.

Secondary locations with bare historical record references (renamed records
now have hyphenated IDs; T-018's inventory named the first two for update):

- `Makefile:16` (`T013`), `scripts/git-hooks/pre-commit:5` (`I027`)
- `internal/board/v2/badges.go:158,195` (`I012`, `O012`),
  `internal/board/v2/objectives.go:191,316` (`I012`, `I013`),
  `internal/board/v2/card.go:109` (`O012`)

## Why It Was Missed

`scripts/hyphenate_record_ids.py` matches only digit IDs (`[ROTCI][0-9]{3,}`),
not `###` placeholders, and T-020's targeted search used the same pattern. The
executor fixed the equivalent placeholders in README.md (I-032) but not in
the two files O-018 names explicitly.

## Proof Needed

`grep -nE '\b[ROTCI](###|[0-9]{3,})\b' AGENTS.md .savepoint/Design.md
.savepoint/Guardrails.md Makefile scripts/git-hooks/pre-commit` plus the
V2 board comments returns no hits. The wording around each placeholder is
unchanged apart from the hyphen.
