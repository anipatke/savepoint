---
type: epic-design
status: audited
---

# E06: Agent Audit + Audit Tab

## Purpose

Implement agent-driven audit workflow and TUI audit tab for viewing audit findings. Agent reviews built code against acceptance criteria and AGENTS.md Code Style rules, writes findings to epic folder, and user approves/agent applies.

## What this epic adds

- **Audit Tab** in Epic Detail overlay (Detail / Audit tabs, press 1/2 to switch)
- **Audit file storage** in epic folder: `{epic}/E##-Audit.md`
- **Updated audit skill** writes findings to new location with Code Style checklist
- **Updated router workflow** (`audit-pending`) calls agent for audit
- **Migration** of existing audit files from `.savepoint/audit/` to epic folder
- **Agent closes epic** after applying audit proposals

## What's NOT built (from previous approach)

- ~~CLI command `savepoint audit`~~ - Agent does audit
- ~~Quality gate runner CLI~~ - Agent runs gates at task end
- ~~Snapshot CLI generation~~ - Agent reads files directly
- ~~TUI proposal review~~ - Review in Epic Detail Audit tab
- ~~Skip handling CLI~~ - Manual skip allowed

## Components

| Module | Purpose |
|--------|---------|
| `internal/board/model.go` | Add `EpicDetailTab`, `EpicAuditContent` state |
| `internal/board/update.go` | Handle `1`/`2` tab keys, load audit file |
| `internal/board/view.go` | Render tab indicator and active tab |
| `internal/board/epic_panel.go` | Add `RenderEpicAuditTab()` function |
| `agent-skills/savepoint-audit/SKILL.md` | Updated file paths + checklist |
| `.savepoint/router.md` | Updated audit-pending workflow |

## Implemented as

- `internal/board/model.go` adds `EpicDetailTab` and `EpicAuditContent`.
- `internal/board/update.go` resets tab state when opening an epic detail overlay and loads `{epic}/E##-Audit.md` the first time the user presses `2`.
- `internal/board/view.go` switches between `RenderEpicDetail` and `RenderEpicAuditTab`.
- `internal/board/epic_panel.go` renders the tab indicator and the audit body, including only `Main Findings` and `Code Style Review`.
- `agent-skills/savepoint-audit/SKILL.md` writes one epic-local `E##-Audit.md`, keeps replacement blocks under `Proposed Changes`, and documents apply/close.
- Existing audit files were migrated from `.savepoint/audit/` to epic-local audit files for v1.1 E02-E05.
- `internal/data/write.go` includes helper functions for applying exact proposal text and updating audit metadata.
- The `savepoint audit` CLI stub and `internal/audit` gate-runner package were removed because Savepoint audit is agent-led and skill-driven.

## Boundaries

**In scope:**
- TUI Audit tab in Epic Detail overlay
- Audit findings display in overlay
- Code Style Review checklist
- Migration of existing audit files
- Agent applies findings and closes epic

**Out of scope:**
- CLI audit command (not needed)
- Toggling checkboxes in UI (manual edit)
- Full proposal diff view (audit tab sufficient)
