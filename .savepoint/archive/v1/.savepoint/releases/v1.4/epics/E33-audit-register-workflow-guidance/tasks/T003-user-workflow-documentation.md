---
id: E33-audit-register-workflow-guidance/T003-user-workflow-documentation
status: done
objective: Document how users review audit prompt, findings, and runs in markdown and the TUI.
depends_on:
    - E32-audit-register-tui-review/T004-help-footer-and-render-regression
    - E33-audit-register-workflow-guidance/T002-template-agent-guidance
complexity_tier: low
complexity_reason: The task adds user-facing documentation for already-planned behavior.
---

# T003: User Workflow Documentation

## Problem

Users need to know where audit-register files live, how to ask agents to use them, and how the read-only TUI review section fits the workflow.

## Context Files

- `README.md`
- `templates/project/AGENTS.md`
- `templates/project/.savepoint/audit/register.md`
- `templates/project/.savepoint/audit/runs/README.md`
- `templates/project/.savepoint/audit/findings/README.md`

## Acceptance Criteria

- [x] README explains the Audit Register at a high level.
- [x] Documentation names `A` as the board shortcut for read-only audit review.
- [x] Documentation explains that markdown files remain the editable source of truth in v1.4.
- [x] Documentation describes stable finding IDs, run history, and proof-based closure.
- [x] Documentation keeps dashboards, external trackers, and automated matching out of v1.4.

## Implementation Plan

- [x] Add a concise Audit Register section to README.
- [x] Explain the `.savepoint/audit/` file layout.
- [x] Explain the TUI review path and read-only boundary.
- [x] Cross-reference generated project guidance where appropriate.

## Context Log

Read: README.md, templates/project/AGENTS.md (T002 result), audit templates (register.md,
findings/README.md, runs/README.md), internal/board/help.go (confirmed `A` = "open audit
register" in the help overlay).

Added an `## Audit Register` README section between Defect Workflow and Agent Skills:
high-level purpose (convergence over cold scans), the `.savepoint/audit/` layout as a file
tree, stable `F###` IDs, immutable run history vs mutable register, proof-based closure
with user-owned waivers, the `A` board shortcut as read-only review, markdown as the
editable source of truth in v1.4, and an explicit exclusion of dashboards, external
tracker integrations, and automated matching. Cross-referenced generated guidance by
pointing users at AGENTS.md, `.savepoint/audit/prompt.md`, and the
`savepoint-audit-register` skill. Also added `A` to the Board bullet list and
`savepoint-audit-register` to the Agent Skills list. Quality gates: `make build && make
test` pass (all packages ok, including the internal/init template freshness tests).
