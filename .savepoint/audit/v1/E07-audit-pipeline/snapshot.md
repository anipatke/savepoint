# Audit Snapshot: E07-audit-pipeline

Generated: 2026-04-28T10:46:09.021Z

## Router State

```yaml
state: audit-pending
release: v1
epic: E07-audit-pipeline
task: (none)
next_action: Run savepoint audit for E07-audit-pipeline.
```

## Tasks

| ID                                                 | Status | Objective                                                                                                                               |
| -------------------------------------------------- | ------ | --------------------------------------------------------------------------------------------------------------------------------------- |
| E07-audit-pipeline/T001-audit-cli-contract         | done   | Define the audit command contract, explicit epic selection, and skip-with-reason logging path without running the full audit pipeline.  |
| E07-audit-pipeline/T002-quality-gate-runner        | done   | Run configured audit quality gates with deterministic pass/fail results and no snapshot side effects.                                   |
| E07-audit-pipeline/T003-snapshot-and-prompt        | done   | Generate a bounded audit snapshot and AI handoff prompt for an epic without embedding source contents.                                  |
| E07-audit-pipeline/T004-audit-orchestration-router | done   | Wire `savepoint audit` to run gates, write audit artifacts, and move the router into audit-pending only after blocking gates pass.      |
| E07-audit-pipeline/T005-proposal-validation-apply  | done   | Validate audit proposal files and apply approved patch-shaped documentation updates with explicit divergence failures.                  |
| E07-audit-pipeline/T006-audit-review-state         | done   | Model audit proposal review decisions so approve, reject, edit, and apply flows are testable before the Ink UI.                         |
| E07-audit-pipeline/T007-audit-review-ui            | done   | Add the Ink audit review mode that lets users inspect proposals, approve or reject them, edit replacements, and apply approved changes. |
| E07-audit-pipeline/T008-audit-pipeline-integration | done   | Verify the complete audit lifecycle from command execution through review application with end-to-end coverage and full quality gates.  |

## Changed Files

| Status | Path                                                                                   |
| ------ | -------------------------------------------------------------------------------------- |
| M      | .claude/settings.local.json                                                            |
| M      | .savepoint/Design.md                                                                   |
| M      | .savepoint/metrics/context-bench.md                                                    |
| M      | .savepoint/releases/v1/v1-PRD.md                                                          |
| M      | .savepoint/releases/v1/epics/E05-init-command/E05-Detail.md                                |
| M      | .savepoint/releases/v1/epics/E05-init-command/tasks/T001-init-cli-contract.md          |
| M      | .savepoint/releases/v1/epics/E05-init-command/tasks/T002-target-validation.md          |
| M      | .savepoint/releases/v1/epics/E05-init-command/tasks/T003-scaffold-writer.md            |
| M      | .savepoint/releases/v1/epics/E05-init-command/tasks/T004-magic-prompt-and-clipboard.md |
| M      | .savepoint/releases/v1/epics/E05-init-command/tasks/T005-dev-deps-install-option.md    |
| M      | .savepoint/releases/v1/epics/E05-init-command/tasks/T006-init-command-integration.md   |
| M      | .savepoint/releases/v1/epics/E06-tui-board/E06-Detail.md                                   |
| D      | .savepoint/releases/v1/epics/E08-doctor-command/E08-Detail.md                              |
| D      | .savepoint/releases/v1/epics/E09-docs-and-packaging/E09-Detail.md                          |
| D      | .savepoint/releases/v1/epics/E10-release-validation/E10-Detail.md                          |
| M      | .savepoint/router.md                                                                   |
| M      | .savepoint/visual-identity.md                                                          |
| M      | AGENTS.md                                                                              |
| M      | CLAUDE.md                                                                              |
| M      | GEMINI.md                                                                              |
| M      | README.md                                                                              |
| M      | eslint.config.js                                                                       |
| M      | package-lock.json                                                                      |
| M      | package.json                                                                           |
| M      | src/cli.ts                                                                             |
| M      | src/cli/args.ts                                                                        |
| M      | src/cli/environment.ts                                                                 |
| M      | src/cli/exit-codes.ts                                                                  |
| M      | src/cli/help.ts                                                                        |
| M      | src/cli/run.ts                                                                         |
| M      | src/commands/audit.ts                                                                  |
| M      | src/commands/board.ts                                                                  |
| M      | src/commands/init.ts                                                                   |
| M      | src/domain/router.ts                                                                   |
| M      | src/fs/markdown.ts                                                                     |
| M      | src/templates/manifest.ts                                                              |
| M      | templates/project/.savepoint/router.md                                                 |
| M      | templates/project/.savepoint/visual-identity.md                                        |
| M      | templates/project/AGENTS.md                                                            |
| M      | templates/prompts/audit-reconciliation.prompt.md                                       |
| M      | templates/prompts/task-breakdown.prompt.md                                             |
| M      | templates/prompts/task-building.prompt.md                                              |
| M      | templates/prompts/task-planning.prompt.md                                              |
| M      | test/cli/args.test.ts                                                                  |
| M      | test/cli/commands.test.ts                                                              |
| M      | test/cli/run.test.ts                                                                   |
| M      | test/domain/router.test.ts                                                             |
| M      | test/templates/project-templates.test.ts                                               |
| M      | tsconfig.json                                                                          |
| M      | vitest.config.js                                                                       |
| M      | vitest.config.ts                                                                       |
| ??     | .savepoint/audit/E05-init-command/                                                     |
| ??     | .savepoint/audit/E06-tui-board/                                                        |
| ??     | .savepoint/audit/E07-audit-pipeline/                                                   |
| ??     | .savepoint/releases/v1/epics/E06-tui-board/tasks/                                      |
| ??     | .savepoint/releases/v1/epics/E07-audit-pipeline/tasks/                                 |
| ??     | .savepoint/releases/v1/epics/E08-board-workflow-cleanup/                               |
| ??     | .savepoint/releases/v1/epics/E09-doctor-command/                                       |
| ??     | .savepoint/releases/v1/epics/E10-docs-and-packaging/                                   |
| ??     | .savepoint/releases/v1/epics/E11-release-validation/                                   |
| ??     | agent-skills/                                                                          |
| ??     | ink-cli-ui-design.zip                                                                  |
| ??     | src/audit/                                                                             |
| ??     | src/init/                                                                              |
| ??     | src/tui/                                                                               |
| ??     | templates/project/.savepoint/metrics/                                                  |
| ??     | templates/prompts/magic-prompt.prompt.md                                               |
| ??     | test/audit/                                                                            |
| ??     | test/commands/                                                                         |
| ??     | test/init/                                                                             |
| ??     | test/tui/                                                                              |

## Tracked Source and Test Files

- src/cli.ts
- src/cli/args.ts
- src/cli/environment.ts
- src/cli/exit-codes.ts
- src/cli/help.ts
- src/cli/run.ts
- src/commands/audit.ts
- src/commands/board.ts
- src/commands/doctor.ts
- src/commands/init.ts
- src/commands/result.ts
- src/domain/config.ts
- src/domain/epic.ts
- src/domain/ids.ts
- src/domain/release.ts
- src/domain/router.ts
- src/domain/status.ts
- src/domain/task.ts
- src/fs/markdown.ts
- src/fs/project.ts
- src/readers/config.ts
- src/readers/epic.ts
- src/readers/release.ts
- src/readers/router.ts
- src/readers/tasks.ts
- src/templates/index.ts
- src/templates/load.ts
- src/templates/manifest.ts
- src/templates/paths.ts
- src/templates/render.ts
- src/validation/dependencies.ts
- src/version.js
- src/version.ts
- test/cli/args.test.ts
- test/cli/commands.test.ts
- test/cli/environment.test.ts
- test/cli/help.test.ts
- test/cli/run.test.ts
- test/domain/config.test.ts
- test/domain/epic.test.ts
- test/domain/ids.test.ts
- test/domain/release.test.ts
- test/domain/router.test.ts
- test/domain/status.test.ts
- test/domain/task.test.ts
- test/fs/markdown.test.ts
- test/fs/project.test.ts
- test/readers/config.test.ts
- test/readers/epic.test.ts
- test/readers/release.test.ts
- test/readers/router.test.ts
- test/readers/tasks.test.ts
- test/smoke.test.js
- test/smoke.test.ts
- test/templates/project-templates.test.ts
- test/templates/prompt-templates.test.ts
- test/templates/render-integrity.test.ts
- test/templates/router-template.test.ts
- test/templates/template-registry.test.ts
- test/validation/dependencies.test.ts
