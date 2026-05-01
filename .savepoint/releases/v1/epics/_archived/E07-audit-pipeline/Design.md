---
type: epic-design
status: audited
---

# Epic E07: audit-pipeline

## Purpose

Implement `savepoint audit`, the hard gate that reconciles built code with Savepoint documentation before the next epic can start. This epic ships the v0.1.0 audit loop without AI semantic review content.

## What this epic adds

- Quality gate runner for configured lint, typecheck, and test commands.
- Gitignore-respecting snapshot generation.
- Changed-files listing without embedding code contents in the snapshot.
- Router transition into audit-pending state.
- Audit prompt output for the user to give to an AI agent.
- Proposal directory creation.
- TUI review mode for approving, rejecting, or editing proposed documentation updates.
- Commit step that applies approved proposals to live files.
- Audit log entries, including skipped audits and reasons.

## Components and files

Expected files introduced or extended by this epic:

| Path                           | Purpose                                         |
| ------------------------------ | ----------------------------------------------- |
| `src/commands/audit.ts`        | Audit command entrypoint.                       |
| `src/audit/quality-gates.ts`   | Configured command runner.                      |
| `src/audit/snapshot.ts`        | File tree and changed-file snapshot generation. |
| `src/audit/router-state.ts`    | Audit-pending state updates.                    |
| `src/audit/prompts.ts`         | Audit handoff prompt generation.                |
| `src/audit/proposals.ts`       | Proposal discovery and validation.              |
| `src/audit/apply-proposals.ts` | Approved proposal application.                  |
| `src/audit/log.ts`             | Audit log writes.                               |
| `src/tui/audit-review/*.tsx`   | Proposal review UI.                             |
| `test/audit/**/*.test.ts`      | Audit pipeline tests.                           |

## Implemented as

- `src/commands/audit.ts` resolves the target epic from `--epic` or router state, supports `--skip --reason`, runs configured quality gates, writes bounded audit artifacts, transitions the router into `audit-pending`, logs audit events, and prints the audit handoff prompt.
- `src/audit/quality-gates.ts`, `snapshot.ts`, `router-state.ts`, `prompts.ts`, `proposals.ts`, `apply-proposals.ts`, and `log.ts` keep audit gates, snapshot generation, router updates, prompt text, proposal parsing, proposal application, and audit logging separately testable.
- `src/commands/board.ts` detects proposal availability and launches `src/tui/audit-review/AuditReviewApp.tsx` after the board exits through the audit-review entry path.
- `src/tui/audit-review/state.ts`, `summary.ts`, `ProposalList.tsx`, `OperationDetail.tsx`, and `AuditReviewApp.tsx` model and render proposal review decisions, replacement edits, apply status, high-divergence acknowledgement, and audit review outcomes.
- `test/audit/`, `test/commands/audit.test.ts`, `test/cli/{args,run}.test.ts`, and `test/tui/audit-review/` cover the audit command contract, quality gates, snapshots, prompt generation, proposal validation/application, review state, Ink interactions, and end-to-end lifecycle behavior.

Implementation deltas from the original plan:

- Proposal review is launched from the board audit-entry path rather than a separate `savepoint audit review` command.
- The audit command creates `.savepoint/audit/{epic}/proposals/`; reconciliation writes one `proposals/proposals.md` bundle, and the loader rejects multiple bundles to keep review state aligned with the workflow.
- AI semantic review remains manual and advisory through the handoff prompt, matching the epic's out-of-scope boundary for API calls and automatic semantic findings.

## Architectural delta

Before this epic, Savepoint can manage workflow state but cannot enforce the documentation audit gate. After this epic, completing an epic can produce a bounded audit context and block progress until documentation proposals are reviewed.

The audit pipeline bridges CLI, filesystem, config, templates, and TUI review mode, but each step should remain separately testable.

## Boundaries

In scope:

- Implement the six-step audit pipeline except AI semantic review generation.
- Support `--skip --reason`.
- Support `--epic` when an explicit epic is needed.
- Generate `.savepoint/audit/{epic}/snapshot.md`.
- Generate `.savepoint/audit/{epic}/proposals/`.

Out of scope:

- Writing semantic-review findings automatically.
- Calling AI APIs.
- Inferring code quality beyond configured mechanical gates.
- Complex merge conflict resolution for proposal application.

## Quality gates

- Tests should cover passing/failing quality gates, snapshot generation, skip logging, proposal validation, and approved proposal application.
- Snapshot generation must not include source code contents.
- Failed blocking quality gates must stop the audit before snapshot generation.

## Design constraints

- Keep audit context bounded to protect token budgets.
- Treat proposal application as a filesystem boundary with explicit failure states.
- Preserve user edits and require confirmation for high-divergence proposal application through the TUI.
