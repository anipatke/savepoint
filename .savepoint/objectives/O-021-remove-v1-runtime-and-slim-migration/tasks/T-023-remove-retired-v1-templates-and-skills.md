---
id: T-023
title: Remove retired V1 templates and skills
objective: O-021
planned_by: {role: planner, session: o021-design-20260923}
status: in_progress
stage: audit
complexity_tier: medium
complexity_reason: "The files themselves are unused, but many init contract tests pin V1 and V2 template trees together, so the tests must be narrowed to V2 without losing the V2 guarantees they also carry."
depends_on: [{task: T-022, requires: clear}]
owner_validation:
    required: false
    accepted_check: ""
---

# T-023: Remove retired V1 templates and skills

## Outcome

The V1 scaffold trees, this repository's retired V1 skills, and the pre-V2
audit-skill upgrade shim are gone. `init` and `upgrade-assets` behave as
they do today for V2 projects, and `upgrade-assets` still refuses a V1
project with its existing migrate guidance.

## User Check

None beyond the gate. Optionally run `savepoint init` in an empty temp
directory and `savepoint upgrade-assets --dry-run` on this repository and
confirm unchanged output.

## Done When

- `templates/project/` and `templates/release/v1/` are deleted.
- `UpgradeProjectAssets` no longer takes a `v1Templates` argument; `main.go`
  and callers are updated. V1-project refusal wording is unchanged.
- The nine retired V1 skill folders in this repository's `agent-skills/`
  (`savepoint-audit-epic`, `savepoint-audit-register`,
  `savepoint-audit-task`, `savepoint-build-task`, `savepoint-create-defect`,
  `savepoint-create-plan`, `savepoint-create-task`, `savepoint-draft-prd`,
  `savepoint-system-design`) and `agent-skills/references/audit-method.md`
  are deleted. `agent-skills/bubbletea-tui-design/` and the four V2 skills
  and three V2 references remain.
- `internal/init/migrate_audit_skill.go` and its test are deleted, and
  `upgradeProjectAssets` no longer calls it. `retire_v1_skills.go` and its
  test remain and still retire the nine V1 skills on a migrated project.
- Tests that pinned V1 trees (for example V1 skill parity in
  `template_freshness_test.go`, V1 sections of `agent_skills_test.go`,
  `lifecycle_test.go`'s `shippedTemplates`, and V1 entries in
  `integration_test.go` and `manifest_test.go`) are removed or narrowed to
  V2. Every V2 guarantee they carried keeps a test.
- A repository search finds no live reference to the deleted paths outside
  `.savepoint/archive/`, immutable Checks, Issues, and closed Task evidence.
- `git diff --check`, `make build && make test-fast` pass.

## Context Files

`templates/project/AGENTS.md`, `templates/release/v1/PRD.md`,
`main.go`, `main_test.go`, `internal/init/upgrade.go`,
`internal/init/upgrade_test.go`, `internal/init/upgrade_schema_test.go`,
`internal/init/migrate_audit_skill.go`,
`internal/init/migrate_audit_skill_test.go`,
`internal/init/retire_v1_skills.go`,
`internal/init/retire_v1_skills_test.go`,
`internal/init/template_freshness_test.go`,
`internal/init/agent_skills_test.go`, `internal/init/lifecycle_test.go`,
`internal/init/integration_test.go`, `internal/init/manifest_test.go`,
`internal/init/scaffold_test.go`, `internal/init/stale_reference_test.go`,
`agent_skills_test.go`.

## Design References

Design sections 2 and 5 (init and upgrade-assets).

## Guardrails

FS-01, FS-02, TPL-01, TPL-02, TPL-03, TPL-04, TEST-03, TEST-06, TEST-08.

## Implementation Plan

1. Delete the template trees and retired skills; build and list failing
   tests.
2. Remove the `v1Templates` parameter and the audit-skill shim call.
3. For each failing test, delete V1-only assertions and keep V2 ones.
4. Search for stale references to the deleted paths and fix live ones.
5. Run `make build && make test-fast`.

## Boundaries

No change to V2 templates, V2 skills, the upgrade provenance manifest
format, or `.savepoint/archive/v1/`. Do not rewrite historical records
that mention deleted paths.

## Technical Verification

Focused `internal/init` and root tests during iteration;
`make build && make test-fast` at handoff (TEST-03 evidence that V2 user
content is still preserved by upgrade).

## Technical Evidence

Executor session t023-build-20260923, 2026-09-23. Toolchain go1.26.2
linux/amd64 (auto-switched to go1.26.8 by `go test`/`go vet`).

### Per-criterion outcome

1. **`templates/project/` and `templates/release/v1/` deleted — met.**
   `git rm -r` removed both trees (98 files: the full V1 scaffold under
   `templates/project/`, including its nine skills, `references/audit-method.md`,
   and `.savepoint/{Concept,Design,Guardrails,Health-Check,PRD}.md`,
   `.savepoint/audit/*`, `.savepoint/releases/v1/*`; and
   `templates/release/v1/PRD.md`). `ls templates` now shows only
   `project-v2` and `prompts`.
2. **`UpgradeProjectAssets` single-tree signature — met.** Signature is now
   `UpgradeProjectAssets(v2Templates fs.FS, targetDir string, dryRun, force bool)`
   (`internal/init/upgrade.go`); the `v1Templates` parameter and its `_ =
   v1Templates` placeholder are gone. `main.go`'s `upgradeAssetsRunner` calls
   `savepointinit.UpgradeProjectAssets(subV2, opts.Dir, opts.DryRun,
   opts.Force)`. `noteMigrateRoute` (the V1-project refusal note) is
   unchanged; `TestUpgradeProjectAssets_v1ProjectCarriesMigrateRouteNote` and
   `TestUpgradeProjectAssets_refusesV1WithoutMutation` still pass unmodified
   in behavior, only their fixture plumbing narrowed to one tree.
3. **Nine V1 skills and the shared audit-method reference deleted — met.**
   `git rm -r` removed `agent-skills/{savepoint-audit-epic,
   savepoint-audit-register, savepoint-audit-task, savepoint-build-task,
   savepoint-create-defect, savepoint-create-plan, savepoint-create-task,
   savepoint-draft-prd, savepoint-system-design}` and
   `agent-skills/references/audit-method.md`. `ls agent-skills` now shows
   exactly `bubbletea-tui-design`, `references`, `savepoint-check`,
   `savepoint-design`, `savepoint-idea`, `savepoint-task`, `superpowers`
   (unowned by this Task, left alone). `agent-skills/references` now holds
   exactly `check-method.md`, `commands-and-procedures.md`,
   `issue-capture.md`.
4. **Shim deleted, retirement kept — met.** `internal/init/migrate_audit_skill.go`
   and `migrate_audit_skill_test.go` are deleted. `upgradeProjectAssets` no
   longer calls `migrateLegacyAuditSkill`. Its shared archive helpers
   (`resolveArchivePath`, `writeArchive`, `removeDirIfEmpty`, the
   `migrationsDir`/`migrationsReadmeName`/`maxArchiveConflictTries` consts,
   the `migrationsReadme` and `legacyMigrationsReadme` text, both still
   needed by `retireV1Skills`) moved into `retire_v1_skills.go`, dropping the
   now-always-true `upgradeLegacyReadme` parameter from `writeArchive` since
   its only remaining caller (`retireV1Asset`) always passed `true`; the
   legacy-README-upgrade behavior itself is preserved byte-for-byte (proved
   by `TestRetireV1Skills_upgradesStockLegacyMigrationsReadmeButPreservesEdits`,
   still passing). `retire_v1_skills.go` and `retire_v1_skills_test.go` are
   otherwise unchanged in behavior; all `TestRetireV1Skills_*` tests pass and
   still retire the nine V1 skills on a migrated (schema_version 2, V1
   assets still on disk) project.
5. **V1-tree-pinned tests removed or narrowed, V2 guarantees kept — met.**
   Full list below (criterion 5 detail).
6. **No live reference to deleted paths — met.** A repository-wide search for
   `templates/project/`, `templates/release/v1`, each of the nine retired
   skill names, `references/audit-method.md`, and `migrate_audit_skill` found
   only: intentional retirement-list constants in `retire_v1_skills.go` /
   `retire_v1_skills_test.go` / `internal/migrate/classify_test.go` (these
   name real, still-existing paths in an already-migrated *user* project,
   which retirement must still recognize); synthetic fixture path strings
   reused as example names in `upgrade_test.go`/`manifest_test.go`/
   `integration_test.go` (never resolved against a real tree); negative
   "must never appear" checks in `main_test.go`/`template_freshness_test.go`/
   `prompt_test.go`; `internal/init/testdata/legacy/*` (frozen pre-existing
   fixtures); `.savepoint/archive/`, `.savepoint/releases/*` (historical
   release records), `.savepoint/objectives/O-001` and `O-018` (closed Task
   evidence), and `project-audit/audit_report_fable_5.md` (a dated,
   frozen third-party audit report) — all excluded categories. Three live,
   non-excluded docs did claim the V1 scaffold still existed and were fixed:
   `.savepoint/Guardrails.md` TPL-01 (dropped the `templates/project/` /
   nine-V1-skill clause, kept the V2 byte-parity clause), `.savepoint/Design.md`
   (the "Template assets" bullet no longer names `templates/project/`),
   and `AGENTS.md` (the "V2 Routing" section's "V1 skills remain available
   only for the V1 scaffold/upgrade path" sentence, and the "Legacy V1
   compatibility" section's "and the V1 scaffold" clause — both corrected to
   say the scaffold is gone; the state-name vocabulary in that section
   (`task-building`, `audit-pending`) was kept as-is because
   `internal/data/lifecycle.go` and `internal/doctor/repairs.go` still parse
   it for archived V1 projects and this repository's own historical router
   records (Objective boundary: V1 `data` readers migrate needs stay), which
   is unrelated to the deleted scaffold/skill files.
7. **Gate — met.** `git diff --check` clean. `make build && make test-fast`
   passed at 2026-09-23. `go build ./...`, `go vet ./...`, and `go test ./...`
   (the full suite, including the slow `internal/migrate` end-to-end tests)
   all passed, run repeatedly through the narrowing iterations and once more
   at the end. `gofmt -l` on every file this Task touched reports nothing.

### Criterion 5 detail: tests removed or narrowed

Deleted outright (subject fully retired, no V2 analog):
`internal/init/audit_contract_test.go` (223 lines: V1 audit-skill and
audit-method contract, router `audit-pending` routing); within
`internal/init/agent_skills_test.go`:
`TestSavepointBuildTaskSkillPresentInBothTrees`,
`TestV1AuditSkillsAndSharedMethodPresentInBothTrees` (+`v1AuditSkillNames`),
`TestSavepointCreateDefectSkillPresentInBothTrees`,
`TestV1SkillSetUnchangedAlongsideV2Routing` (+`v1SkillNames`),
`TestAgentsGuidesCarryV2RoutingSection`,
`TestV2DesignHandoffMatchesRoutingContract` (+`v2RoutingRequiredPhrases`) —
the last two compared a "## V2 Routing" section between live `AGENTS.md` and
`templates/project/AGENTS.md`; `templates/project-v2/AGENTS.md` never had
that section, so there is no surviving tree to compare against, and the
section's own guarantee (shipped guidance matches live) is not otherwise
claimed for it; within `internal/init/skill_validation_test.go`:
`TestSplitAuditSkillsPassStructureValidation`,
`TestSharedAuditMethodIsNonTriggerableReference` (+`skillRoots`,
`v1TemplateSkillRoot`, `splitAuditSkillNames`); within
`internal/init/template_freshness_test.go`: `TestProjectTemplatesUseCurrentWorkflow`
(V1 router/AGENTS.md/audit-epic wording), `TestProjectConceptTemplateExists`
(V2 ships no Concept.md), `TestUpgradeMigratesLegacyAuditSkillFromRealTemplates`,
`TestProjectAuditRegisterTemplatesExist`, `TestUpgradeAddsAuditRegisterTemplatesFromRealTemplates`
(V2 ships no `.savepoint/audit/`); within `internal/init/scaffold_test.go`:
the split-audit-skill half of `TestScaffold_installsSplitAuditSkillsAndSharedMethod`
(deleted whole, no V2 audit skills exist); within `internal/data/router_test.go`:
`TestRouterReader_shippedTemplateCarriesBothAnchors` (V1 `ReadState` against
the V1 template; superseded by the already-existing
`TestReadStateV2_shippedTemplateDecodesCleanly` in `router_v2_test.go`,
unmodified, which proves the same guarantee for the V2 template via
`ReadStateV2`).

Narrowed to the V2 tree only, keeping every V2-relevant assertion:
`internal/init/lifecycle_matrix_test.go` (`realTrees`→`realV2Tree`;
dropped the "migrated V2 project" and "legacy V1 project" cases from
`TestLifecycleMatrix_fullLoopAcrossProjectKinds`,
`_absentOptionalFilesProduceNoFinding`, and
`_dryRunMatchesRealRunAcrossProjectKinds` — coverage for those states
against synthetic trees already exists in `retire_v1_skills_test.go`;
`TestLifecycleMatrix_retirementAndConflictReportedTogether` now builds its
"migrated project carrying V1 skills" fixture with the existing
`v2ProjectWithLegacySkills` helper instead of a real V1 scaffold);
`internal/init/lifecycle_test.go` (`shippedTemplates` now walks
`templates/project-v2`; all five `TestLifecycle_*` tests pass unmodified
otherwise, since the ownership matrix (`wantOwnership`) and test bodies were
already tree-agnostic); `internal/init/scaffold_test.go`
(`scaffoldFromRealTemplates` now scaffolds from `templates/project-v2`;
`TestScaffold_installsGuardrailsAndHealthCheck` renamed
`TestScaffold_installsGuardrails` and narrowed to Guardrails.md only, since
V2 ships no Health-Check.md); `internal/init/stale_reference_test.go`
unmodified but now exercises the V2 tree transitively through
`scaffoldFromRealTemplates`; `internal/init/upgrade_schema_test.go`
(`v1v2Templates`→single V2 `fstest.MapFS`; deleted
`TestUpgradeAssetsFromTree_preservesFrozenPreE47Fixture`, which tested the
now-deleted legacy-audit-skill migration specifically); `internal/init/retire_v1_skills_test.go`
(`retirementTemplates`→single V2 `fstest.MapFS`; dropped the incidental
legacy-audit-skill-file write from the "stock README is upgraded" subtest of
`TestRetireV1Skills_upgradesStockLegacyMigrationsReadmeButPreservesEdits`,
since it no longer does anything now the shim is gone — the README-upgrade
behavior it verifies now comes from `retireV1Skills` alone, still exercised);
`internal/init/template_freshness_test.go` (`TestProjectGuidanceTemplatesMirrorLiveGuidance`
dropped the V1-vs-live canonical-vocabulary comparison — root `AGENTS.md` and
`templates/project-v2/AGENTS.md` serve different audiences and were never
required to share prose, unlike the V1 template which was scaffolded from
the same source — kept and narrowed the V2 skill/reference byte-parity and
completeness assertions; `TestProjectTemplatesRejectStaleWorkflowTerms`
narrowed to live `AGENTS.md` only; `TestProjectGuardrailsAndHealthCheckTemplatesExist`
renamed `TestProjectGuardrailsTemplateExists`, narrowed to the V2
Guardrails.md (V2 ships no Health-Check.md); `TestProjectDocumentTemplatesHaveTypeFrontmatter`
repointed from V1's PRD/Design/Concept to V2's Design/Idea/Guardrails, all of
which do carry `type:` frontmatter; `TestProjectAgentsGuidesLifecycleTerminologyConsistency`
repointed at `templates/project-v2/AGENTS.md`, which does carry the same
canonical status phrase; `TestUpgradeDeliversPolicyAssetsFromRealTemplates`
repointed at the V2 tree and narrowed to Guardrails.md, the only policy asset
it ships); `main_test.go` and root `agent_skills_test.go` (`v1OnlySkills`/
`v1OnlyPaths` negative-check lists needed no change; `TestBundledSavepointSkillsHaveDiscoveryFrontmatter`,
`TestProjectAgentGuideIncludesLocalSkillFallback`, and
`TestScaffoldedSavepointSkillsMatchBundledSkills` dropped their `templates/project`
tree argument, keeping the V2 tree check, which already carries the fallback
instruction and frontmatter this proves).

Unchanged (already single-tree or V1-agnostic, confirmed still passing after
the signature/deletion changes): `internal/init/upgrade_test.go` (its
`v1Templates, v2Templates` two-arg call sites became single-arg via the same
mechanical `UpgradeProjectAssets(v2, ...)` rewrite; no test logic changed),
`internal/init/upgrade_failure_test.go`, `internal/init/manifest_test.go`,
`internal/init/integration_test.go`, `internal/init/prompt_test.go`.

### Files read

Budgeted Context Files as listed, plus extra reads beyond that budget,
logged with reason: `internal/init/lifecycle_matrix_test.go`,
`internal/init/skill_validation_test.go` (both pinned the real V1 tree and
had to be narrowed, but weren't named in the Task's Context Files list);
`internal/init/audit_contract_test.go` (found failing after the deletion,
turned out to be wholly V1-scoped); `internal/data/router_test.go` (root
`make test-full` surfaced its shipped-V1-template dependency);
`.savepoint/Guardrails.md`, `.savepoint/Design.md`, `AGENTS.md` (root)
(criterion 6's stale-reference search surfaced factually-false claims that
V1 skills/templates still ship, in project governance docs outside this
Task's Context Files); `internal/migrate/classify_test.go`,
`internal/init/lifecycle_test.go` (confirming which V1-skill-name string
literals are intentional retirement-list data versus stale references);
`templates/project-v2/.savepoint/{Design,Guardrails,Idea}.md`,
`templates/project-v2/AGENTS.md` (checking what V2-equivalent content
existed before narrowing tests to point at it instead of inventing new
assertions); `examples.md`, `project-audit/audit_report_fable_5.md`
(both matched the stale-reference grep; read to confirm they are,
respectively, a standalone worked-example document unrelated to this
project's shipped assets, and a frozen dated audit report — neither edited).

### Files changed

- Deleted: `templates/project/` (98 files), `templates/release/v1/PRD.md`,
  the nine `agent-skills/savepoint-*/SKILL.md` directories named in
  criterion 3, `agent-skills/references/audit-method.md`,
  `internal/init/migrate_audit_skill.go`,
  `internal/init/migrate_audit_skill_test.go`,
  `internal/init/audit_contract_test.go`.
- Edited (production): `internal/init/upgrade.go` (signature, dropped
  `migrateLegacyAuditSkill` call, one stale comment), `internal/init/retire_v1_skills.go`
  (absorbed the shared archive helpers), `main.go` (call site).
- Edited (tests): `internal/init/{upgrade_schema,retire_v1_skills,
  lifecycle_matrix,lifecycle,scaffold,skill_validation,template_freshness,
  agent_skills,upgrade}_test.go`, root `agent_skills_test.go`,
  `internal/data/router_test.go`.
- Edited (docs): `.savepoint/Guardrails.md`, `.savepoint/Design.md`,
  `AGENTS.md`.
- `.savepoint/router.md`: corrected a stale `next_action` left over from
  T-022's handoff (still named T-022 after `task:` had already been
  advanced to T-023) before starting this Task; not itself part of this
  Task's Done When.
- The Task record itself (lifecycle and this evidence).

### Limitations

- `savepoint init` and `savepoint upgrade-assets --dry-run` (the User Check)
  were not run as commands (AGENTS.md: agents never run `savepoint`); the
  equivalent behavior is covered by the passing `internal/init` and root
  test suites, including tests that scaffold and upgrade from the real
  `templates/project-v2` tree on disk. The owner's User Check covers running
  the actual binary.
- `deadcode` was not re-run; T-022 already covers dead-function detection
  and this Task's deletions are whole files/directories with no remaining
  callers (proved by `go build`/`go vet`/`go test` all passing), not a
  function-level trim.
- `examples.md` and `project-audit/audit_report_fable_5.md` still name
  retired V1 skills/paths; left untouched as a standalone example document
  and a frozen dated audit report, respectively, neither of which is live
  shipped guidance.

### Handoff

`stage: audit` — ready for an optional Task Check or, under an explicit
owner waiver, for the mandatory Full Objective Check. This is not a pass.

## Drift Notes

None expected. `writeArchive` lost its `upgradeLegacyReadme` parameter
(always `true` at its one remaining call site after the shim's removal);
behavior is unchanged and proven so by the still-passing legacy-README-format
upgrade test, so this is recorded as a drift note rather than a plan
deviation.
