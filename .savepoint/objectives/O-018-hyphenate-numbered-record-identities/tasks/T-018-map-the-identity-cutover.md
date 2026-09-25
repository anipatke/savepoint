---
id: T-018
title: Map the identity cutover
objective: O-018
status: done
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: o018-design-20260923}
check_waiver:
    task: T-018
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-09-23T05:26:10Z"
---

# T-018: Map the identity cutover

## Outcome

A complete, reviewable inventory and migration contract names each active ID field, path, generator, consumer, historical exception, and recovery state before implementation changes the graph.

## User Check

No separate owner interaction is needed. The optional Task Check may be requested; if skipped, record the owner's explicit waiver in Task evidence.

## Done When

The inventory accounts for all five record kinds, nested references, router, CLI, renderers, V1-to-V2 allocation, templates, skills, tests, fixtures, and this repository's live files. The contract states source drift, collision, mixed input, interruption, archive, and re-run behavior with named expected outcomes. A checkable decision deliverable is recorded in Technical Evidence.

## Context Files

`internal/data/project.go`, `internal/data/router_v2.go`, `internal/data/release_v2.go`, `internal/data/objective_v2.go`, `internal/data/task_v2.go`, `internal/data/check_v2.go`, `internal/data/issue_v2.go`, `internal/data/evidence_v2.go`, `internal/data/write.go`, `internal/migrate/plan.go`, `.savepoint/objectives/O-018-hyphenate-numbered-record-identities/Objective.md`.

## Design References

Design sections 1, 2, 4, 6, 7, and 11. O-018 Confirmed Design Decisions.

## Guardrails

FS-01, FS-03, FS-04, DATA-01..03, ARCH-03, TEST-01..04, TEST-08.

## Implementation Plan

1. Inventory structured identity locations and path conventions from the named sources.
2. Define the canonical active format and legacy migration-only input boundary.
3. Specify preview, guarded apply, preserved Check copy location, progress record, recovery, and ordinary-load refusal.
4. Name unresolved technical choices explicitly; return REPLAN REQUIRED if the design cannot satisfy preservation and recovery together.

## Boundaries

No production behavior change or active graph rewrite in this Task.

## Technical Verification

Review the inventory against each named decoder and generator; record a traceable coverage table. Reference `agent-skills/references/check-method.md` for later independent evaluation. Run `make build && make test-fast` at Task handoff if files change.

## Technical Evidence

### Decision deliverable — identity inventory and migration contract

Inventory completed on 2026-09-23. This Task changes no production behavior and does not rewrite the live record graph.

#### Canonical identity and record map

Use one active V2 grammar for all five record kinds: `^(R|O|T|C|I)-[0-9]{3,}$`. The migration map is a prefix insertion only (a hyphen goes between the kind letter and the digits, giving `R-006`, etc.); retain every digit, including zero padding, and every authored slug. Do not renumber records or derive an identity from a title/path. The dedicated migration reader may recognize the old form while ordinary post-cutover V2 loading accepts only the canonical form.

| Kind | Record location and identity | Typed references and validation/consumer path |
|---|---|---|
| Release (`R`) | `.savepoint/releases/R###-slug/Release.md`, frontmatter `id`; release discovery and decoder in `internal/data/discover.go` / `release_v2.go`. | `Objective.release`, router `release`, and Check `scope.id` when `scope.kind: release`; Objective membership and dangling-reference checks in `project.go`. Shared Check evidence fields are also resolved for Releases. `Task.release` is excluded: it is opaque transitional packaging metadata, not an active membership edge. |
| Objective (`O`) | `.savepoint/objectives/O###-slug/Objective.md`, frontmatter `id`; indexed by `project.go`. | `Task.objective`, `Objective.depends_on[]`, router `objective`, Check `scope.id` for objective scope, and Issue `escalated_to`; references are validated/resolved by the Objective/Task decoders, Issue decoder, and `project.go`. |
| Task (`T`) | `.savepoint/objectives/O###-slug/tasks/T###-slug.md`, frontmatter `id`; ownership is the Task's `objective` field, not its directory. | `Task.depends_on[].task`, Issue `tasks[]`, Check `scope.id` for task scope, router `task`, and `check_waiver.task`; the serialized `Task.release` field is transitional packaging metadata, not a typed Release reference. |
| Check (`C`) | Flat `.savepoint/checks/` record with frontmatter `id`; decoded by `check_v2.go`, linked and ordered by `project.go`. | `Check.scope.id` (typed by `scope.kind`), `Check.supersedes`, `Check.issues[]`; Issue `source.check`, `resolution.check`, `history[].check`, and `checks[]`; Task/Objective/Release evidence `last_check`, `freshness.check`, `owner_validation.accepted_check`, and `exception.check`. |
| Issue (`I`) | Flat `.savepoint/issues/I###-slug.md`, frontmatter `id`; decoded by `issue_v2.go`. | `Issue.duplicate_of`, `escalated_to`, `tasks[]`, `checks[]`; Check `issues[]`; and the Check references in Issue source, resolution, and history described above. `project.go` validates link targets, duplicate graph, and escalation targets. |

The router's authored selection fields are exactly `release`, `objective`, and `task` in `.savepoint/router.md`. Derived membership/reverse links such as `ObjectiveTasks`, `ReleaseObjectives`, `ScopeChecks`, and Issue/Check backlinks are rebuilt from the rewritten record fields; they are not separate authored identities to migrate. `Task.release` is retained as opaque transitional packaging metadata (the two live values are `v2`), not a Release identity. `guardrail_ids`, actor roles/sessions, timestamps, hashes, commit IDs, and non-record path slugs are not numbered record identities.

The complete structured surface is in the Task's named data files: `release_v2.go`, `objective_v2.go`, `task_v2.go`, `check_v2.go`, `issue_v2.go`, `evidence_v2.go`, and `router_v2.go` decode shapes; `project.go` resolves cross-record links and Check ordering; `write.go` validates router selections and owns Check/Issue creation and Issue link writes. There are currently five family regex declarations across those decoders (Release, Objective, Task, Check, Issue); validation should converge on one kind-aware identity rule. Check ID numeric ordering in `project.go` and Check/Issue next-ID suffix parsing in `write.go` also need to skip the new `-` separator.

#### Paths, generation, and consumers

| Surface | Current behavior and cutover requirement |
|---|---|
| Path discovery | Design §2 and `internal/data/discover.go` define the active roots. Release discovery examines direct children of `.savepoint/releases` that contain `Release.md`; Objective and Task discovery is under `.savepoint/objectives/**/tasks`; Checks and Issues are flat under `.savepoint/checks` and `.savepoint/issues`. Discovery validates that directory/file basenames begin with the record's declared ID. Insert the ID hyphen in every path prefix and preserve the remainder of each path and slug. |
| New Check/Issue IDs | `internal/data/write.go` currently emits `C%03d` and `I%03d`. `CreateCheckV2` writes `checks/{id}.md`; `CreateIssueV2` writes `issues/{id}-{slug}.md`. Both allocation and path generation must emit hyphenated IDs. |
| V1→V2 allocation | `internal/migrate/plan.go`'s `idAllocator` emits `O/T/I/R` plus at least three digits while planning Release, Objective, Task, and Issue targets. It does not allocate Check IDs. Preserve deterministic allocation order and reservations, but generate canonical IDs for the V2 targets; keep V1 source IDs and archived bytes intact. |
| Hand-authored IDs | Objective and Task IDs have no Go generator; active planner/task skill examples and instructions mint them by hand. Release records are likewise authored rather than allocated by a V2 record-creation API. Update active guidance and examples to the same grammar. |
| CLI and display | `cmd/board.go` parses `--release`, `--epic`, and `--objective` as strings; schema-aware meaning/refusal is behind board dispatch. The V2 board, doctor, and resume consume the data index/projection and display IDs from it; they must retain the canonical value and use data validation rather than local identity regexes. `internal/doctor/repairs.go` currently prints `R###`/`O###`/`T###`/`C###`/`I###` guidance and must be updated. |
| Active text/assets | Targeted search found ID syntax/examples in `README.md`, root `AGENTS.md`, `.savepoint/Design.md`, `.savepoint/router.md`, `Makefile`, `scripts/git-hooks/pre-commit`, active root `agent-skills/`, and V2 scaffold `templates/project-v2/` (including its managed guide and skill copies). Update literal references and format examples; keep the established root/scaffold skill-copy equality. The V1 scaffold `templates/project/` and retired V1 skills remain compatibility assets and are not rewritten. |

A path discrepancy must be resolved before cutover: Design §7 documents Check files as `C###-slug.md`, and the live Checks have that form, while `CreateCheckV2` writes bare `C###.md`. Discovery allows either because it checks the basename prefix. The migration should preserve each existing suffix exactly after changing the ID prefix; T-020/T-019 must decide and test the future Check-creator naming convention so docs and generated records agree.

#### Repository and verification inventory

The active canonical path inventory contains 6 Release, 10 Objective, 21 Task, 5 Check, and 31 Issue files (73 records) under the roots above. Their identities and all links must be discovered from the strict active index at migration time rather than frozen as a hand-maintained list. A targeted live-file search found two `Task.release: v2` compatibility values; preserve these opaque values byte-for-byte and never convert them to `R-v2` or infer a Release edge from them. The path search also found two old Task files under `.savepoint/releases/v2/epics`; V2 discovery does not traverse that legacy subtree, so those files are outside the live index and remain unchanged unless the planner explicitly reclassifies them. V1 migration goldens/fixtures under `internal/migrate/testdata` and `internal/data/testdata/migration`, V1 board/record tests, the V1 `templates/project/` tree, and historical/archived material are compatibility inputs, not active V2 output to rewrite.

The affected test/fixture families are: `internal/data` identity decoders, `project_test.go`, `write_test.go`, router/evidence/graph and V2 end-to-end tests; `internal/board/v2`, `internal/doctor`, `internal/resume`, `cmd/board_test.go`, and `main*_test.go` output/filter fixtures; `internal/migrate` plan/conversion/operation/recovery tests (old V1 inputs remain old, V2 planned outputs become canonical); and `internal/init` template-freshness/skill-copy tests plus all current active V2 fixtures. Migration coverage must include malformed IDs, mixed old/new identities or references, duplicate/colliding target IDs and paths, source drift, interrupted writes at each phase, Check archive byte equality, preservation of opaque `Task.release: v2`, ordinary-load refusal while pending, and safe recovery/no-op rerun. Expected negative result is a named refusal before mutation when preflight fails; unexpected drift during recovery leaves the marker and source bytes in place for explicit recovery rather than overwriting them.

#### Migration contract for T-019/T-021

1. Preview performs a read-only strict load of the complete old-format V2 graph, validates reference closure, calculates the full typed ID and path rename maps, and hashes every source and intended output. It reports all planned field/path changes. Dry-run writes no marker, stage, archive, or project file.
2. Before apply writes anything, reject malformed or unexplained mixed-format input, dangling references, duplicate identities, unsafe paths, any occupied destination, and any source whose bytes differ from the preview hash. These failures report the offending record/path and leave the project unchanged. Do not guess a map for an identity absent from the loaded graph.
3. Preserve the opaque transitional `Task.release` field (currently `v2` in two live Tasks); it is not a Release reference in the V2 index, so keep it unchanged under DATA-01. Proposed durable operation record: `.savepoint/migrations/v2-identity-cutover.yml`, written atomically before the first project replacement, storing the complete map, source/output hashes, planned path moves, and per-file phase/progress. Proposed preserved-Check location: `.savepoint/archive/v2-pre-identity-cutover/checks/<original-check-basename>`; create-only and verify byte-for-byte. A different pre-existing archive copy is a collision and is never overwritten. These names are a T-018 contract proposal; T-019 must confirm them against the existing file-replace helper before implementation.
4. Stage all transformed outputs, archive original Checks, and verify the staged graph before removing old paths. Multi-file replacement is not filesystem-atomic; each replacement uses the shared platform replace primitive, while the operation record makes the sequence recoverable. Ordinary `LoadProject`/board/doctor/resume refuse whenever an operation is unfinished and show the recovery command. Only the dedicated migration reader may load intact old input or the recorded in-progress state.
5. If interrupted, rerun reads the same operation record and accepts a pending source only when it matches its recorded source hash, or an already-applied destination only when it matches the recorded output hash. It then safely finishes remaining steps. Unexpected bytes, missing required stage data, or a conflicting archive/destination stop recovery without deletion or overwrite and leave the recovery marker intact. Archived original Checks remain byte-identical; active Checks contain the mapped ID and references.
6. On success, strict-load the complete new-format graph, mark the operation `complete`, and retain the completed record as audit evidence. A rerun on that unchanged canonical graph reports already migrated/unchanged and writes nothing. Ordinary loading accepts the canonical graph and rejects any old-form identity; V1 history and legacy fixture inputs stay untouched.

The dedicated V2→V2 operation must remain separate from the V1→V2 conversion body in `internal/migrate/plan.go`, while following its preview → guarded apply → interruption recovery → collision rejection shape and reusing its atomic platform file-replace primitive as O-018 specifies. The command/subcommand name and the exact persisted progress-record schema were not fixed by the Objective; T-019 should settle them without changing the contract above. This mapping surfaces no preservation/recovery conflict that requires `REPLAN REQUIRED`.

### Acceptance-criterion evidence

- All five record families, typed nested references, paths, generators, consumers, templates/skills, test/fixture families, and current active root counts are covered in the tables above.
- The contract names source drift, collision, mixed input, interruption, archive, ordinary-load blocking, and rerun outcomes, including which bytes are accepted and which operations refuse before mutation.
- Verification commands: `make build` passed as the first half of `make build && make test-fast`. `make test-fast` failed twice in `internal/init.TestIntegration_InstallDependencies` (`internal/init/integration_test.go:298`): WSL resolved Windows npm 11.9.0, which treated the temporary WSL path as an unsupported UNC cwd and failed with EPERM opening `C:\Windows\package-lock.json`. The environment probe found no Linux `node`/`npm` on PATH. This is a failed test gate; it is not recorded as `CLEAR`.
- Owner direction received 2026-09-23: make no further test runs for this file-only Task. No retry was run after that direction; the `make test-fast` gate remains failed/unverified, and the limitation is carried to review.
- Decision deliverable recorded here for independent Task/Objective review; no source behavior or live graph rewrite performed in T-018.

### Reads and changes

- Read the Task's Context Files: `internal/data/project.go`, `router_v2.go`, `release_v2.go`, `objective_v2.go`, `task_v2.go`, `check_v2.go`, `issue_v2.go`, `evidence_v2.go`, `write.go`, `internal/migrate/plan.go`, and the owning `Objective.md`.
- Reference reads named by the Task: `.savepoint/Design.md` §§1, 2, 4, 6, 7, 11; `.savepoint/Guardrails.md` rules FS-01/03/04, DATA-01..03, ARCH-03, TEST-01..04, TEST-08.
- Extra source read: `internal/data/discover.go`, to verify exact roots, discovery depth, and filename-prefix checks needed for path mapping. A targeted `release:` search of live Objective/Task files found two Objective Release links and two opaque Task `release: v2` values; the latter are excluded from the identity map.
- Extra targeted searches (no other source bodies opened): identity patterns/references and consumers in `cmd/**`, `internal/data/**`, `internal/board/**`, `internal/doctor/**`, `internal/resume/**`, `internal/migrate/**`, `internal/init/**`, `main*.go`; WSL Node/npm resolution diagnostic; identity examples in `README.md`, `AGENTS.md`, `.savepoint/**`, `agent-skills/**`, `templates/project-v2/**`, `Makefile`, and `scripts/git-hooks/pre-commit`; active record paths under `.savepoint/releases`, `objectives`, `checks`, and `issues`; and legacy V1 fixture/archive paths for classification.
- Changed only this Task record's lifecycle/evidence. No implementation, test, fixture, template, skill, router, Objective, or live graph file was changed.
- Limitations: no migration code or identity-format behavior is implemented in T-018; `make test-fast` did not pass and was not rerun after the owner's direction; no Task Check was requested or waived. The task is review-ready at `stage: audit` and still needs the owner's completion action.

## Drift Notes

- `CreateCheckV2` writes `C###.md`, while Design §7 and current active Check paths use `C###-slug.md`; preserve suffixes during migration and reconcile the future creator convention in O-018 before cutover. `Task.release` is a separate nonidentity compatibility field (`v2` in two live Tasks) and must remain opaque.
- `.savepoint/releases/v2/epics` contains two old Task files but is outside V2 discovery. Preserve it as non-indexed legacy material unless O-018 explicitly reclassifies it.
- The exact identity-migration CLI name and progress-record schema are not set by O-018; settle them in T-019 while retaining this migration contract. `task_v2.go` explicitly treats `Task.release` as transitional packaging metadata and the live `v2` values confirm it is not a Release ID; preserve it, do not include it in the R-ID map.