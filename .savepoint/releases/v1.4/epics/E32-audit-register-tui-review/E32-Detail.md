---
type: epic-design
status: audited
---

# E32: Audit Register TUI Review

## Purpose

Expose the Audit Register in the board as a read-only section where users can review the audit prompt, current findings, finding detail, and run history.

## What this epic adds

- `A` top-level board shortcut for the Audit Register.
- Prompt, Findings, and Runs tabs.
- Read-only finding list grouped for quick scanning.
- Read-only finding detail drill-in.
- Linked-findings backlinks from epic and task detail into the finding detail.
- Footer, help, empty-state, reload, and watch behavior for audit files.

## Components and files

| Module | Purpose |
|--------|---------|
| `internal/board/model.go` | Store audit overlay selection, cursor, and scroll offsets |
| `internal/board/io.go` | Load audit-register data asynchronously |
| `internal/board/update.go` | Route audit overlay keys |
| `internal/board/view.go` | Render the audit overlay on top of the board |
| `internal/board/audit_overlay.go` | Render prompt, findings, and runs tabs |
| `internal/board/audit_detail.go` | Render finding detail |
| `internal/board/detail.go` | Add linked-findings section to task detail |
| `internal/board/epic_panel.go` | Surface linked findings in epic detail |
| `internal/data` | Provide structured audit-register data and work-item reverse lookups |

## Architectural delta

The board gains one read-only overlay backed by the data package's audit-register model. It follows existing release-docs and defects overlay patterns and does not mutate audit files.

## Boundaries

**In scope:**
- `A` key entry point
- Prompt, Findings, and Runs tabs
- Finding detail drill-in
- Linked-findings backlinks from epic and task detail (read-only, drill into finding detail)
- Empty states for missing audit assets
- Help/footer updates

**Out of scope:**
- Editing findings
- Creating audit runs
- Marking findings verified, waived, deferred, or duplicate
- Automatic reconciliation
- Plain-output audit summaries

## Quality gates

- Board tests cover overlay open/close, tab switching, list navigation, detail rendering, empty states, and help/footer text.
- Existing board overlay tests remain unchanged.
- `go test ./internal/board ./internal/data` passes.

## Open decisions

None.
