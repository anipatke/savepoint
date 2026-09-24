---
id: T-020
title: Rename this project's records
objective: O-018
status: done
depends_on: [{task: T-019, requires: clear}]
owner_validation:
    required: true
    accepted_check: C-919
    accepted_by: {role: owner, session: owner-chat-20260924}
planned_by: {role: planner, session: o018-replan-20260923}
check_waiver:
    task: T-020
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-09-23T08:21:52Z"
---

# T-020: Rename this project's records

## Outcome

This repository's active records and guidance use hyphenated identities, via
a rename script kept in the repo for other projects.

## User Check

Review the board and resume output after the rename and confirm it still
points at the intended O-018 work. The optional Task Check may be requested
or explicitly waived; the Full Objective Check remains mandatory.

## Done When

The script renames directories/files under `.savepoint/{releases,objectives,checks,issues}`
and rewrites identity references in frontmatter, bodies, `router.md`,
`AGENTS.md`, `Design.md`, and `Guardrails.md`, leaving `.savepoint/archive/`
untouched. Strict load, board, doctor, and resume succeed; a targeted search
finds no unhyphenated active IDs; `git diff --check` and fresh
`make test-full` pass.

## Context Files

`.savepoint/router.md`, `AGENTS.md`, `.savepoint/Design.md`,
`.savepoint/Guardrails.md`, `.savepoint/objectives/`, `.savepoint/issues/`,
`.savepoint/checks/`, `.savepoint/releases/`.

## Guardrails

FS-01..04, DATA-01..04, TEST-01..05.

## Implementation Plan

1. Write the script under `scripts/`; match word-bounded `[ROTCI][0-9]{3,}`
   only, so words and version strings are untouched.
2. Run it, review the diff, and fix any false positives by hand.
3. Verify load, board, doctor, resume, and the full gate.

## Boundaries

No migration command, recovery state, or archived Check copies. Check
content changes only in identity references.

## Technical Verification

Fresh `make test-full` at handoff.

## Technical Evidence

Start-of-work check: T-019's dependency requirement was `clear`. T-019's
frontmatter carries an explicit owner `check_waiver` (recorded via the board,
`actor: {role: owner, session: board-owner}`, `2026-09-23T07:57:57Z`), which
`savepoint-task` satisfies a downstream `requires: clear` dependency (it does
not satisfy `requires: accepted`, but T-020 asks for `clear`). Objective O-018
is `planned` and ready. Confirmed `git status` before touching anything: the
working tree already carried T-019's uncommitted implementation, matching the
same-commit boundary in O-018's Confirmed Design Decisions.

Pre-rename baseline (recorded before any T-020 change): `make build` and
`make test-fast` passed against T-019's already-uncommitted state. `doctor`
and `resume` against the live repository both failed to load
(`[invalid_v2] ... objective id "O001" must match O- plus at least three
digits`), confirming T-019's validation change and this repository's
still-bare records were in exactly the pre-rename state O-018 describes.

Acceptance evidence:

- Script: `scripts/hyphenate_record_ids.py` renames identity-named entries
  (word-bounded `^[ROTCI][0-9]{3,}` at the start of a directory/file name)
  under `.savepoint/{releases,objectives,checks,issues}` to their hyphenated
  form, recursing into matched entries to catch nested identity-named
  children (e.g. a renamed Objective directory's `tasks/T0##-*.md` files),
  and rewrites every word-bounded bare identity occurrence across those
  renamed files plus `.savepoint/router.md`, `AGENTS.md`,
  `.savepoint/Design.md`, and `.savepoint/Guardrails.md`. An entry whose own
  name does not match the identity pattern (e.g. `.savepoint/releases/v1`,
  `v1.1`, ..., `v2`, the un-archived V1 source hierarchy sitting alongside
  the V2 `R0##-*` release directories) is left untouched and not recursed
  into, so V1-scheme `T0##`/`D0##` identities inside that hierarchy are never
  touched. Supports `--dry-run` and `--root`; uses `git mv` (falls back to a
  plain rename outside a git repository).
- Dry run reviewed before executing: renamed 27 directories/files (6
  releases, 10 objectives, 20 nested tasks, 5 checks, 31 issues — some
  objectives have no tasks) and identified 76 files for content rewrite (74
  changed, `AGENTS.md`/`Guardrails.md` had no bare identities to rewrite). No
  false positives: the V1 `releases/v1*`/`v2` trees and `.savepoint/archive/`
  did not appear in the plan.
- Ran for real; `git status` showed every rename as `RM` (git-detected
  rename + the subsequent content rewrite). Diffed a sample
  (`O-018`'s own Objective.md, `router.md`) and confirmed hyphenated
  `id:`/`objective:`/`release:`/`task:` fields and headings.
- Targeted search `grep -rnoE '\b[ROTCI][0-9]{3,}\b'` across
  `.savepoint/releases/R-*`, `.savepoint/objectives`, `.savepoint/checks`,
  `.savepoint/issues`, `.savepoint/router.md`, `.savepoint/Design.md`,
  `.savepoint/Guardrails.md`, and `AGENTS.md`: no matches. `.savepoint/v1*`
  legacy release trees and `.savepoint/archive/` confirmed unchanged via
  `git status` (no entries).
- Strict load, board, doctor, resume: rebuilt (`go build`), then `doctor`
  and `resume` both loaded the project successfully (no `invalid_v2`
  error). Compared `doctor`'s findings against a `git worktree` built from
  the last commit (before T-019/T-020): the same set of pre-existing
  findings (done-without-clearance Tasks, an incomplete Release, open
  Issues) appears in both, only reformatted to the hyphenated IDs — the
  rename changed no project state, only identity formatting. `board` has no
  non-interactive flag; ran it with stdin redirected from `/dev/null`,
  which prints its non-interactive summary (confirms it loads) showing
  hyphenated `T-020`, `O-018`, `R-006`.
- `git diff HEAD --check`: no output (no whitespace/line-ending
  violations).
- Fresh `make test-full` (2026-09-23T08:13:00Z–08:15:15Z, after the README
  repair below): PASS — full Go suite plus Linux/Darwin/Windows builds.

Extra read/write beyond Context Files: a repo-wide search after the rename
found three stale bare-identity examples/placeholders in `README.md`
(`--objective O001`, the `O001-example`/`R001-example` directory sample, and
`R###`/`O###` format placeholders including `release: R###`) — not this
repository's own records and not one of T-020's or O-018's named files, so
out of Context-File scope. Per the Out-Of-Scope Repair default
(`agent-skills/references/issue-capture.md`), fixed directly as a small,
self-contained doc correction and recorded **Issue I-032** (`type: drift`,
`status: open`) with observed/repair-attempted history; left open for a
checker or explicit owner acceptance. No Check record or Issue closure was
written by this Task.

Verification evidence:

- PASS: `make build`.
- PASS: `make test-fast` (pre-rename baseline) and a fresh `make test-full`
  (post-rename, post-README-repair; see above).
- PASS: `doctor` and `resume` load the renamed project; `board` loads
  (non-interactive summary).
- PASS: `git diff HEAD --check`.
- PASS: targeted bare-identity search returns no matches across the
  Task's Context Files.

Limitations: `board`'s interactive TUI was not exercised in a TTY (this
environment is headless); only its load path and non-interactive summary
were verified. Command help text (`--help` output, e.g. `board
[--objective <objective>]`) was not changed and was not audited for stale
identity examples beyond README.md. `.savepoint/Idea.md` and
`.savepoint/visual-identity.md` were checked (no bare identities present)
but are outside T-020's named scope, so were not modified.

Files read beyond Context Files: `README.md` (extra read/write, logged
above). `internal/migrate/classify.go`, `AGENTS.md`'s workflow description,
and a `git worktree` of the last commit were read only to confirm the
V1-hierarchy exclusion and the pre-existing-findings comparison; nothing in
them was changed.

Files changed: `scripts/hyphenate_record_ids.py` (new); every renamed
directory/file under `.savepoint/{releases,objectives,checks,issues}` per
the dry-run plan above; `.savepoint/router.md`; `README.md`;
`.savepoint/issues/I-032-readme-identity-examples-were-still-bare.md` (new);
this Task's own frontmatter/evidence.

### Correction after C-910 (2026-09-23)

The dry-run claim "No false positives" above was inaccurate. The script also
rewrote three bare-versus-hyphenated examples, so each compared an ID with
itself, and it missed the `X###` placeholders in AGENTS.md and Design.md.
C-910 recorded these as I-034 and I-033. Both were repaired at the owner's
direction; see those Issues' `repair_attempted` history. This Task's status
is unchanged.
