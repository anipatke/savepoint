# Agent State Machine

## Read order

1. This file (router.md)
2. Current state → next action
3. Active epic E##-Detail.md
4. Active task file

Read `.savepoint/PRD.md` only for vision changes. Read `.savepoint/Design.md` only for architecture/audit.

## Current state

```yaml
state: task-building
release: v1.1
epic: E14-structural-improvements
task: E14-structural-improvements/T005-shell-tokenization
next_action: Build E14-structural-improvements/T005-shell-tokenization.
```

## State → action

### pre-implementation

PRD + Design locked, no epics yet.

**Next:** 1) Read release PRD for epic list, 2) Define + confirm epic order, 3) Create epic stubs. Transition to `epic-design` for E01.

### epic-design

Epic E##-Detail.md is empty/stub.

**Next:** Define what this epic adds, files it touches, and architectural delta. Then transition to `epic-task-breakdown`.

### epic-task-breakdown

Epic exists, tasks missing.

**Next:** 1) Re-read epic, 2) Create task files at `tasks/TNNN-slug.md` with `status: planned`, `objective`, `depends_on`, 3) Add `## Implementation Plan` checkboxes per task. When all planned → first unblocked task.

### task-building

Task `in_progress`, depends satisfied.

**Next:** When starting work, set task `status: in_progress` and press `p` in TUI to mark the focused task as router priority. Execute plan, tick checkboxes, run quality gates, update router to next task or `audit-pending`. Stop.

### audit-pending

Epic complete, needs audit before next epic.

**Next:** Fresh audit agent reads epic E##-Detail.md, task files, drift notes, Design.md, AGENTS.md, and scoped changed files. Write one audit file to `.savepoint/releases/{release}/epics/{E##-epic}/E##-Audit.md`:
- `## Main Findings` user-facing narrative only
- `## Code Style Review` checklist against AGENTS.md rules
- `## Proposed Changes` admin/apply blocks using `### Target File`, `### Replace`, `### With`

After user approves: apply proposals, mark epic `status: audited`, update Design.md `last_audited`, advance router.
