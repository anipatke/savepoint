---
type: audit-findings
audited: 2026-07-26
---

# Audit Findings: E35 Split Task and Epic Audits

## Main Findings

### Verdict

CLEAR after apply. The single finding below was approved and applied: the
canonical architecture record now describes the upgrade migration introduced by
this epic. The split audit behavior, migration, scaffold, and regression
coverage were already working. No owner-run evidence or waiver is needed; the
repository handoff result is **CLEAR TO COMMIT/PUSH**.

### What Needs Attention

Nothing outstanding. The one finding raised by this audit is resolved.

#### Resolved — the architecture record omitted the new upgrade behavior

E35 makes `upgrade-assets` refresh shared files below
`agent-skills/references/` and archive a retired generic audit skill below
`.savepoint/migrations/` before removing its triggerable copy. The live Design
said the command refreshes only `agent-skills/**/SKILL.md` and skips other
project state, leaving the repository's canonical map behind both the
implementation and the user-facing README, with T003's recorded drift
unreconciled.

Applied: the `upgrade-assets` architecture bullet at `.savepoint/Design.md:26`
now names shared references under `agent-skills/references/` and the
preserve-before-remove migration exception. No code or test change was required,
and none was made.

Evidence: `.savepoint/Design.md:26` now agrees with
`internal/init/upgrade.go:104-115` and
`internal/init/migrate_audit_skill.go:10-41`; the drift recorded in
`.savepoint/releases/v1.5/epics/E35-split-task-and-epic-audits/tasks/T003-upgrade-audit-skill-migration.md`
is reconciled.

### Materiality Summary

No materiality actions are required. The one finding was applied at closeout:

| Finding | Likelihood | Impact | Materiality | Recommendation |
|---|---|---|---|---|
| The architecture record omits the new upgrade behavior — **resolved** | High — every architecture reader received the stale command boundary | Low — runtime behavior and user documentation were correct | Medium — architecture reconciliation is an explicit epic audit gate | Fixed now via the Design.md replacement below; applied |

### What Is Proven / Not Proven

**Proven**

- Live and scaffolded projects contain distinct task- and epic-audit skills,
  one non-triggerable shared method, matching canonical/template copies, and no
  active generic audit folder.
- Router, AGENTS, build-task, audit-register, Design template, and README
  guidance select the correct skill without adding a task-audit router state.
- The skill contracts cover scope freezing, concrete coverage tables,
  side-effect review, convergence, file reality, materiality, health-check
  modes, artifact authority, and compact handoff formats.
- Upgrade behavior covers absent, stock, locally modified, repeated, dry-run,
  and archive-conflict states; legacy content is preserved before removal,
  split assets and the shared reference install, and reporting distinguishes
  migration.
- Fresh scaffolds install the split assets plus interpolated Guardrails.md and
  Health-Check.md. Skill validation and stale-reference checks pass over both
  live and generated trees while retaining historical records.

- The live architecture record is current for the upgrade migration, reconciled
  at closeout by the applied Design.md correction.

**Not proven**

- Nothing. Every acceptance criterion is proven; there are no unverified
  criteria and no owner-run gates outstanding.

### Audit Evidence

- **Scope lock:** all acceptance criteria in T001–T004; live/template audit
  assets and guidance; `Scaffold`; `UpgradeProjectAssets`, migration/report
  helpers, atomic writes, and their tests. Historical release wording and
  unrelated init behavior were outside the supported boundary.
- **Coverage and workflow result:** checked public entry points and generated
  representations across fresh, missing, stock, modified, duplicate,
  conflicting, repeated, dry-run, and normal-write states. The migration
  sequence was traced from target validation through archive selection/write,
  legacy removal, template installation, and report output. Network,
  authentication, numeric, terminal, and Unicode-width cells are not
  applicable to this filesystem/Markdown change.
- **File reality and drift:** every created or edited file named by the four
  task logs exists; the two generic skill entrypoints are intentionally
  deleted; canonical/template pairs are byte-identical. T001/T002 drift is
  resolved. T003's migration-map drift remains as the finding above; T004
  introduced no additional architecture drift.
- **Gates:** focused E35 `internal/init` tests pass; `git diff --check` is clean;
  `make build` passes; `make test` passes for all nine tested packages.

### Guardrails Verification

- **Rule IDs checked:** none mapped. This repository has no live
  `.savepoint/Guardrails.md`; the file is an E35 scaffold template.
- **Health check mode:** Full. This repository has no live
  `.savepoint/Health-Check.md`, so its file-specific ceremony is not
  applicable; the epic acceptance, file-reality, drift, focused-test, build,
  and full-test evidence was still completed.
- **File reality evidence:** complete, with intentional generic-skill deletions
  accounted for.
- **Waivers or unresolved findings:** no waivers; the one architecture-drift
  finding was applied at closeout and none remain.

### Non-Blocking Observations

At audit time the router still named E35/T004 under `task-building` and the epic
detail remained `status: planned`, although all four task files were `done`. The
explicit completed-epic audit request validly overrode that routing state.
Lifecycle closeout has since been completed: the epic is `audited`, Design.md
`last_audited` points at this epic, and the router has advanced.

## Code Style Review

- [x] One job per file
- [x] One job per function
- [x] Test branches
- [x] Types document intent
- [x] Build only what is needed
- [x] Handle errors at boundaries
- [x] One source of truth — the live Design.md upgrade boundary was stale and is now corrected
- [x] Comments explain WHY
- [x] Content in data files
- [x] Small diffs

## Proposed Changes

### Target File

.savepoint/Design.md

### Replace

```md
- **Upgrade-assets command** (`savepoint upgrade-assets [dir] [--dry-run] [--force]`) refreshes package-owned `agent-skills/**/SKILL.md` files and the managed block in the root agent guide from embedded templates for existing Savepoint projects, while skipping `.savepoint/PRD.md`, `.savepoint/Design.md`, `.savepoint/releases/**`, and other project state.
```

### With

```md
- **Upgrade-assets command** (`savepoint upgrade-assets [dir] [--dry-run] [--force]`) refreshes package-owned `agent-skills/**/SKILL.md` files, shared references under `agent-skills/references/`, and the managed block in the root agent guide from embedded templates for existing Savepoint projects. It skips project-owned state, except that a retired generic audit skill is preserved under the non-triggerable `.savepoint/migrations/` archive before its active copy is removed.
```
