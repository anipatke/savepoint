---
id: T-026
title: Reconcile guidance with the smaller code base
objective: O-021
planned_by: {role: planner, session: o021-design-20260923}
status: done
complexity_tier: low
complexity_reason: "Documentation-only reconciliation plus a line-count report; no production behavior changes."
depends_on: [{task: T-025, requires: clear}]
owner_validation:
    required: false
    accepted_check: ""
check_waiver:
    task: T-026
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-09-23T22:13:37Z"
---

# T-026: Reconcile guidance with the smaller code base

## Outcome

AGENTS.md, `Design.md`, and the V2 scaffold guidance describe the code base
as it now is: no V1 board, no V1 templates or skills, a `data`-owned schema
check, and a clean-git migration. The Objective records the measured size
reduction.

## User Check

Read the AGENTS.md Codebase Map: every row is one or two plain sentences,
and no row describes a deleted package, file, or behavior.

## Done When

- AGENTS.md's Codebase Map lists only live modules, each in one or two
  sentences. Rows for deleted V1 board behavior, `release_cutover.go`,
  recoverable apply, and the V1 template tree are gone or corrected. The
  closing "V1 readers ... compatibility code" paragraph states what
  remains: the V1 readers `internal/migrate` uses.
- `Design.md` no longer claims recoverable apply, interruption recovery,
  or the migration journal (command table, `.savepoint/` tree, test
  table, and header notes), and states the clean-git rule and the `data`
  schema check.
- `templates/project-v2/AGENTS.md` and any V2 skill text that mentions
  recovery or pending migrations are aligned; template freshness tests
  pass.
- The Objective's Technical Evidence (or this Task's) records production
  and test line counts per package against the baseline in O-021's Why,
  and a final `deadcode` report.
- `git diff --check`, `make build && make test-fast` pass.

## Context Files

`AGENTS.md`, `.savepoint/Design.md`, `templates/project-v2/AGENTS.md`,
`README.md`, `internal/init/template_freshness_test.go`,
`internal/init/agent_skills_test.go`.

## Design References

Design sections 1, 2, 10, and 11.

## Guardrails

ARCH-04, TPL-01, TPL-02, TEST-06, TEST-08.

## Implementation Plan

1. Search guidance for references to deleted code and behavior.
2. Rewrite the Codebase Map rows concisely and correct Design sections.
3. Align V2 scaffold guidance and run template tests.
4. Measure and record line counts and `deadcode` output.

## Boundaries

No production code changes. No edits to archived, Check, or closed Issue
records. O-014's router guidance changes stay with O-014.

## Technical Verification

`make build && make test-fast`; the Full Objective Check then runs a fresh
`make test-full`.

## Technical Evidence

T-025 is recorded done with owner acceptance in the router. T-026 is now
`in_progress`, stage `audit`; scoped documentation edits and required gates
are complete.

Context reads: `.savepoint/router.md`, this Task, the owning `Objective.md`,
`AGENTS.md`, `.savepoint/Design.md`, `templates/project-v2/AGENTS.md`,
`README.md`, `internal/init/template_freshness_test.go`,
`internal/init/agent_skills_test.go`, and ARCH-04, TPL-01, TPL-02, TEST-06,
and TEST-08 in `.savepoint/Guardrails.md`.

Extra reads and reasons: searched `agent-skills/` and
`templates/project-v2/agent-skills/` for recovery, pending-migration,
interruption, journal, preflight, and migration-state wording because the Task
requires active V2 guidance alignment; no V2 skill text matched. Read targeted
schema and command-boundary symbols in `internal/data`, `internal/migrate`,
`internal/board`, `internal/doctor`, and `main.go` to verify the Codebase Map
and apply guidance. Read `internal/data/runtime.go` and the Git-path checks in
`internal/migrate/command.go`, `apply.go`, and their tests to confirm the named
schema diagnostic and Git path check. Searched `go.mod` and `go.sum` plus the
local Go module cache for a deadcode runner after `deadcode .` was not on PATH;
the cached runner is `/home/user/go/bin/deadcode`. Read the pre-O-021 Go
snapshot and current tracked Go sources to produce the required per-package
line counts. `git status --short` showed no pre-existing changes beyond this
Task's lifecycle/evidence update.

Documentation findings applied: removed the dead V1 board/package descriptions
and `release_cutover.go` from the Codebase Map, stated that `internal/migrate`
uses the surviving V1 readers, and corrected migration apply and schema-check
guidance. The template AGENTS guide and canonical/scaffold V2 skills had no
recovery or pending-migration claims to change.

Go line counts are newline counts per package (`production / tests`). Baseline
is commit `a6283d5` from 2026-09-23, the O-013 completion snapshot referenced
by O-021; it matches O-021's approximate 33.4k production and 55.1k test-line
baseline.

| Package | Baseline | Current | Change |
| --- | ---: | ---: | ---: |
| `.` | 302 / 2,075 | 293 / 2,004 | -9 / -71 |
| `cmd` | 357 / 806 | 350 / 811 | -7 / +5 |
| `internal/board` | 4,930 / 7,826 | 89 / 300 | -4,841 / -7,526 |
| `internal/board/v2` | 6,356 / 5,869 | 6,198 / 5,732 | -158 / -137 |
| `internal/buildtool` | 730 / 413 | 730 / 413 | 0 / 0 |
| `internal/data` | 9,515 / 16,750 | 7,272 / 12,416 | -2,243 / -4,334 |
| `internal/doctor` | 2,496 / 3,720 | 1,334 / 1,722 | -1,162 / -1,998 |
| `internal/init` | 1,659 / 8,960 | 1,541 / 7,672 | -118 / -1,288 |
| `internal/migrate` | 5,934 / 7,999 | 4,248 / 5,936 | -1,686 / -2,063 |
| `internal/resume` | 692 / 650 | 688 / 638 | -4 / -12 |
| `internal/styles` | 319 / 185 | 319 / 185 | 0 / 0 |
| `internal/testutil` | 142 / 0 | 142 / 0 | 0 / 0 |
| **Total** | **33,432 / 55,253** | **23,204 / 37,829** | **-10,228 / -17,424** |

Verification: `git diff --check` passed. `make build && make test-fast`
completed successfully (exit 0). The gate ran the template and guidance cases
`TestProjectTemplatesRejectStaleWorkflowTerms`,
`TestV2WorkflowAssetsHaveNoActiveV1Routing`,
`TestProjectAgentsGuidesLifecycleTerminologyConsistency`,
`TestProjectGuidanceTemplatesMirrorLiveGuidance`,
`TestV2SkillSetIsCompleteWithByteParity`, and
`TestV2SkillSetHasNoObsoleteVocabulary`.

Command history: `deadcode .` exited 127 because the binary was not on PATH.
`/home/user/go/bin/deadcode .` first exited 1 for a missing Go build-cache
entry, then exited 1 because the sandbox could not read `/home/user/.cache/go-build`.
After approval to use the existing cache, the final
`/home/user/go/bin/deadcode .` report completed (exit 0):

```text
internal/data/project.go:26:6: unreachable func: LoadProject
internal/data/project.go:43:6: unreachable func: loadProjectV1
internal/data/project.go:50:6: unreachable func: loadProjectV2
internal/migrate/command.go:79:6: unreachable func: FindProjectRoot
internal/migrate/manifest.go:181:6: unreachable func: UnmarshalManifest
internal/migrate/manifest.go:195:6: unreachable func: WriteManifestCreateOnly
```

No unreachable functions were reported in `internal/board` or
`internal/doctor`. The three `internal/data` findings mean O-021's stated
success condition of zero unreachable functions in that package still needs
the Full Objective Check to assess. The initial missing-cache and sandbox
read-only errors, then approval to read the Go cache, are part of the report's
command history.

Files changed: `AGENTS.md`, `.savepoint/Design.md`, `README.md`, and this Task
record. No production code, test source, scaffold guidance, canonical skill, or
scaffold skill changed.

Acceptance criteria:

1. **Met.** The Codebase Map lists the live modules in one or two sentences
   each, removes the V1 board and deleted cutover/apply/template descriptions,
   and identifies the V1 readers retained for `internal/migrate`.
2. **Met.** Design no longer describes recoverable writes, interruption
   recovery, or a journal. It documents the data-owned schema check and the
   clean-Git-path apply rule.
3. **Met.** The V2 template and skill searches found no stale recovery or
   pending-migration wording. The named template and byte-parity tests passed.
4. **Met as a report.** The before/after per-package counts and final deadcode
   output are recorded. The report still lists three unreachable `internal/data`
   functions, which the Full Objective Check must assess against O-021's
   success condition.
5. **Met.** `git diff --check` and `make build && make test-fast` passed.

## Drift Notes

The final `deadcode .` report identifies three unreachable functions in
`internal/data`; O-021's Full Objective Check must assess them against its
success condition requiring zero unreachable functions. No Task-scope plan
deviation occurred.
