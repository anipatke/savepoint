---
id: T-058
title: Add the health setup command
objective: O-028
status: done
depends_on: [{task: T-056, requires: clear}]
owner_validation:
    required: true
    accepted_check: ""
planned_by: {role: planner, session: planning-o028-20261001}
check_waiver:
    task: T-058
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-01T08:02:26Z"
---

# Add the health setup command

## Outcome

The owner runs `savepoint health setup [dir]` to see suggested health tools and what would change, and `savepoint health setup --apply` to save them; `savepoint init` ends with the same preview; re-running setup on a configured project is the upgrade path and never changes confirmed entries silently.

## User Check

In a scratch copy of a small Go or JS project: run `savepoint init`, read the health preview it prints; run `savepoint health setup` and confirm nothing was written; run `--apply` and read `.savepoint/health/config.json`; run setup again and see "no changes"; delete the tool from PATH or the project and see it reported as missing without being removed. Confirm the OSV-Scanner line says the scanner contacts OSV.dev.

## Done When

- `savepoint health setup [dir]` loads the project, runs `Discover`, loads any existing config, and prints a plain-language preview: each proposal, its reason and gap, default exclusions, and a closing line naming `--apply`. It writes nothing (FS-03).
- `--apply` writes only `.savepoint/health/config.json` through `Store.SaveConfig`; a second `--apply` on unchanged input reports unchanged (FS-04). Missing or non-Savepoint directories fail clearly without partial writes (FS-06).
- **Reconciliation on an existing config:** entries the owner already confirmed are kept byte-for-byte, including edited args, thresholds, `required`, and exclusions. New proposals are added only with `--apply`. Configured instances whose executable or inputs are now missing are reported, never removed automatically. The preview separates "new", "unchanged", and "needs attention".
- `savepoint init` prints the same preview after scaffolding and the `--apply` hint; it writes no health config itself. Discovery failure during init is reported and does not fail init.
- Board, resume, doctor, and upgrade-assets never discover or write health files.
- `agent-skills/savepoint-design/SKILL.md` gains one line telling the planner to describe confirmed health tools in Design.md in plain words; its scaffold copy stays byte-identical (TPL-01). AGENTS.md (live and template managed block) and Design.md CLI surface name the new human-only command.

## Context Files

`.savepoint/objectives/O-028-code-health-provider-execution-and-scaffolding/Objective.md`; `main.go`; `cmd/init.go`; `cmd/init_test.go`; new `cmd/health.go`, `cmd/health_test.go`; `internal/codehealth/discovery.go`; `internal/codehealth/storage.go`; new `internal/codehealth/setup.go`, `internal/codehealth/setup_test.go`; `agent-skills/savepoint-design/SKILL.md`; `templates/project-v2/agent-skills/savepoint-design/SKILL.md`; `AGENTS.md`; `templates/project-v2/AGENTS.md`.

## Design References

Design sections 1, 6, and 11; O-028 Confirmed Design Decisions (2026-10-01).

## Guardrails

FS-01, FS-03, FS-04, FS-06, TPL-01, TPL-02, TPL-04, ARCH-01, ARCH-04, CFG-01, TEST-01, TEST-02, TEST-03, TEST-04, TEST-08.

## Implementation Plan

1. Confirm `Discover` and its `Proposal` type exist; return REPLAN REQUIRED if not.
2. Put reconciliation in `internal/codehealth/setup.go`: `Plan(existing *Config, proposals) SetupPlan` (new, unchanged, needs attention) and `SetupPlan.Apply(store)`; rendering to plain text lives beside it as a pure function.
3. Add thin `cmd/health.go` dispatch for `health setup [dir] [--apply]` and wire it in `main.go` and help text (ARCH-01).
4. Call the preview from `cmd/init.go` after a successful scaffold; report and continue on discovery error.
5. Update the design skill (both copies), AGENTS.md CLI Rules/Codebase Map (both copies as applicable), and Design section 6.
6. Tests: preview writes nothing (bytes and mtimes unchanged), apply then re-apply is unchanged, owner-edited entries survive reconciliation, missing tool reported not removed, init prints preview and still succeeds when discovery fails, non-project directory refused, skill copies identical.

## Boundaries

No collection or refresh command, no interactive prompts, no Design.md writes by the command, no tool installation, no changes to doctor or the board.

## Technical Verification

Focused runs while iterating. Scaffolding and file writes are covered by TEST-03: fresh `make test-full` at handoff.

## Technical Evidence

Executor evidence for a fresh check to verify; not a Check and not clearance.

**Per criterion**

- Preview command: `savepoint health setup [dir]` resolves the target, runs `Discover`, loads any config, prints New / Unchanged / Needs attention / Not suggested, reasons, gaps, exclusions and a closing `--apply` line; writes nothing. Evidence: `TestPreviewWritesNothing` (bytes and mtimes), `TestMainHealthSetupPreviewWritesNothing`.
- `--apply`: writes only `.savepoint/health/config.json` via `Store.SaveConfig`; second apply reports "No changes". Evidence: `TestApplyThenReapplyIsUnchanged`, `TestMainHealthSetupApplyWritesOnlyConfigAndRepeatsAsUnchanged`. Missing and non-project directories fail with no `.savepoint` created: `TestMainHealthSetupRefusesMissingAndNonSavepointDirectories`, `TestApplyRefusesNonProjectWithoutWriting`. A corrupt existing config fails clearly: `TestUnusableExistingConfigFailsClearly`.
- Reconciliation: edited args, thresholds, `required`, timeout and exclusions survive (`TestReconciliationKeepsOwnerEdits`); a missing executable or scope is reported and kept (`TestMissingToolIsReportedNotRemoved`, `TestMissingScopeIsReported`); a hand-formatted file is not rewritten when nothing is new (`TestApplyWithNothingNewKeepsHandEditedFileBytes`); a proposal that would conflict with the existing config is refused with nothing written (`TestConflictingProposalIsNotSavedAndIsExplained`).
- `savepoint init` prints the same preview after scaffolding and writes no health config (`TestMainInitEndsWithHealthPreviewAndWritesNoHealthConfig`); a discovery failure only warns on stderr (`TestMainInitStillSucceedsWhenHealthDiscoveryFails`).
- Board, resume, doctor, upgrade-assets: not touched; no discovery or health writes added to them.
- Docs: design skill line added to both copies (identical; `TestScaffoldedSavepointSkillsMatchBundledSkills` passes); AGENTS.md CLI Rules in both copies, live Codebase Map rows for `cmd/` and `internal/codehealth/`, and Design section 6 name the human-only command.
- OSV line: preview states the scanner, not Savepoint, contacts OSV.dev (asserted in tests).

**Commands run**: `make build && make test-fast` passed; `make test-full` passed (exit 0, includes cross-platform builds). Manual preview of a scratch Go project read as expected.

**Files changed**: new `cmd/health.go`, `cmd/health_test.go`, `internal/codehealth/setup.go`, `internal/codehealth/setup_test.go`, `main_health_test.go`; edited `main.go`, `AGENTS.md`, `templates/project-v2/AGENTS.md`, both design skill copies, `.savepoint/Design.md`.

**Extra reads** (outside Context Files): `internal/codehealth/config.go` (Config/CapabilityConfig fields and validation), `cmd/resume.go` (command pattern), `internal/codehealth/discovery_catalogue.go` (default exclusions), `internal/codehealth/discovery_test.go` (test helpers), `main_test.go` (subprocess test helpers), `.savepoint/Guardrails.md` (named rule text), `.savepoint/Design.md` section 6.

**Limitations**: `cmd/init_test.go` was read but not changed; init behavior is tested in `main_health_test.go`. The owner User Check (scratch project walkthrough) is not yet done. Setup also adds proposals whose tool is not installed yet (they are marked and reported as needing attention later). The discovery-failure init test triggers failure through a corrupt existing config and cannot see the stderr warning on success.

## Drift Notes

Design section 6 CLI table and AGENTS.md CLI Rules gain `savepoint health setup`; reconcile at the Full Objective Check.
