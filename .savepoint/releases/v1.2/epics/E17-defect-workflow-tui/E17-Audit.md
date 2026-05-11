---
type: audit-findings
audited: 2026-05-11
---

# Audit Findings: E17 Defect Workflow TUI

## Main Findings

E17 is closed as audited. The implementation is functionally present: the data layer has a defect model and lifecycle validation, the board loads release defects, renders open-defect pressure, provides the `d` defects overlay, opens defect details, shows related task-card markers, and supports the `space` resolve shortcut. Doctor has defect diagnostics, and live/scaffold guidance includes the defect lane.

The latest quality gates pass:

- `make build` passed.
- `make test` passed.

Audit apply reconciled the lifecycle documentation drift:

- `T010-defect-resolve-hotkey.md` is now `status: done`, clearing the user-owned status blocker noted during audit.
- `.savepoint/Design.md`, `E17-Detail.md`, and E17 task ledger text now describe defect statuses as `open` -> `in_progress` -> `resolved`.
- Template and skill surfaces were checked; `AGENTS.md`, `templates/project/AGENTS.md`, `agent-skills/savepoint-create-defect/SKILL.md`, and `templates/project/agent-skills/savepoint-create-defect/SKILL.md` already use the defect-specific lifecycle.
- E17 was marked audited, `.savepoint/Design.md` `last_audited` was advanced to `v1.2/E17-defect-workflow-tui`, and the router now points to E18.

No must-fix code defect was found in the T010 hotkey path. Pressing `space` on an `open` defect writes `status: resolved`, preserves body/unrelated frontmatter, clears `stage`, keeps the user in the Defects overlay, and no-ops clearly for `resolved` and `in_progress` defects.

Residual process note: some earlier E17 task files still have incomplete checklist/context-log bookkeeping despite being marked `done`. That does not block E17 closeout because the implemented behavior is covered by source tests and the lifecycle drift has been reconciled.

## Code Style Review

- [x] One job per file
- [x] One job per function
- [x] Test branches
- [x] Types document intent
- [x] Build only what is needed
- [x] Handle errors at boundaries
- [x] One source of truth
- [x] Comments explain WHY
- [x] Content in data files
- [x] Small diffs

## Proposed Changes

### Target File
.savepoint/Design.md

### Replace
```md
Defects use the same `planned`, `in_progress`, and `done` status vocabulary. `stage` is required while a defect is `in_progress`. Router state may enter `defect-building` with a `defect` field naming the active repair item, which the board renders as a `DEFECT` Next Activity line.
```

### With
```md
Defects use defect-specific lifecycle statuses: `open`, `in_progress`, and `resolved`. `stage` is required while a defect is `in_progress`, and must be absent once the defect is `open` or `resolved`. Router state may enter `defect-building` with a `defect` field naming the active repair item, which the board renders as a `DEFECT` Next Activity line.
```

### Target File
.savepoint/releases/v1.2/epics/E17-defect-workflow-tui/E17-Detail.md

### Replace
```md
- Defect model, parser, release-scoped discovery, and lifecycle validation live in `internal/data`.
- Board loading reads release defects alongside tasks and renders an open-defect header signal without adding a fourth column.
- The `d` key opens a defect overlay; enter opens a defect detail overlay that renders the evidence sections from the defect markdown body.
- Related task cards can show compact defect markers when a defect `reference` matches a visible task.
- Doctor validates defect files for parse errors, required fields, invalid status/stage, and broken task-like references.
- Scaffolded AGENTS/router templates and README now document the defect lane, but `defect-building` state guidance needs the audit proposal updates before close.
```

### With
```md
- Defect model, parser, release-scoped discovery, and lifecycle validation live in `internal/data` using the defect-specific statuses `open`, `in_progress`, and `resolved`.
- Board loading reads release defects alongside tasks and renders an open-defect header signal without adding a fourth column.
- The `d` key opens a defect overlay; enter opens a defect detail overlay that renders the evidence sections from the defect markdown body.
- Pressing `space` on an `open` defect in the overlay resolves it through the canonical defect write helper while preserving unrelated frontmatter and body content.
- Related task cards can show compact defect markers when a defect `reference` matches a visible task.
- Doctor validates defect files for parse errors, required fields, invalid status/stage, and broken task-like references.
- Scaffolded AGENTS/router templates, defect skill guidance, and README now document the defect lane and defect-specific lifecycle language.
```

### Target File
.savepoint/releases/v1.2/epics/E17-defect-workflow-tui/tasks/T010-defect-resolve-hotkey.md

### Replace
```md
- [x] In the Defects overlay, pressing `space` on a selected defect with `status: planned` updates that defect file to `status: done`.
```

### With
```md
- [x] In the Defects overlay, pressing `space` on a selected defect with `status: open` updates that defect file to `status: resolved`.
```

### Target File
.savepoint/releases/v1.2/epics/E17-defect-workflow-tui/tasks/T010-defect-resolve-hotkey.md

### Replace
```md
- [x] Pressing `space` on a selected defect that is already `done` does not rewrite the defect file or produce an error state.
```

### With
```md
- [x] Pressing `space` on a selected defect that is already `resolved` does not rewrite the defect file or produce an error state.
```

### Target File
.savepoint/releases/v1.2/epics/E17-defect-workflow-tui/tasks/T001-defect-data-model.md

### Replace
```md
- [ ] Defect status is limited to `planned`, `in_progress`, and `done`
```

### With
```md
- [ ] Defect status is limited to `open`, `in_progress`, and `resolved`
```

### Target File
.savepoint/releases/v1.2/epics/E17-defect-workflow-tui/tasks/T010-defect-resolve-hotkey.md

### Replace
```md
Allow a user to resolve a planned defect directly from the Defects overlay by pressing `space` on the focused defect row.
```

### With
```md
Allow a user to resolve an open defect directly from the Defects overlay by pressing `space` on the focused defect row.
```

### Target File
.savepoint/releases/v1.2/epics/E17-defect-workflow-tui/tasks/T010-defect-resolve-hotkey.md

### Replace
```md
- [x] Add or update board tests covering planned-to-done, done no-op, in-progress lifecycle behavior, and focus stability.
```

### With
```md
- [x] Add or update board tests covering open-to-resolved, resolved no-op, in-progress lifecycle behavior, and focus stability.
```
