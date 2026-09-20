---
type: audit-findings
audited: 2026-07-03
---

# Audit Findings: E32 Audit Register TUI Review

## Main Findings

**Verification outcome: pass.** All five tasks are `done` and every acceptance criterion was verified against the code and its tests. `make build` passes and a fresh (uncached) `go test ./internal/board ./internal/data` run passes. The epic delivers exactly what E32-Detail.md scoped: the `A` overlay with Prompt/Findings/Runs tabs, grouped finding list with cursor and drill-in detail, empty states, watch/reload refresh of audit data, help/footer hints, and linked-findings backlinks in task and epic detail — all read-only, with an explicit test guaranteeing no key mutates a finding or writes audit files.

Task-level verification notes:

- **T001** — Startup, overlay-open, and watch/reload audit loads all exist (`loadAuditBestEffort`, `loadAuditRegisterCmd`, `reloadMsg.audit`); a missing `.savepoint/audit/` tree degrades to an empty set and never blocks startup. The documented AC deviation (the `A`-key trigger living in T002) resolved as described. The watcher tolerates a missing `audit/` directory and picks up a later-created one via the root watch's create events.
- **T002** — Tab strip, `[`/`]`/arrows/`h`/`l` switching, per-tab scroll preservation, esc/`q` close, and all three empty states are implemented and tested.
- **T003** — Findings group by status → severity → ID via the loader's `SortFindings`; the renderer and cursor follow-scroll share one grouping pass (`auditFindingLayout`) so they cannot disagree. Detail shows every field the AC lists; absent links/locations/proof are omitted.
- **T004** — `A:audits` fits at 80/100/120 columns (footer moved to a `footerHints` constant, separators tightened, `p:priority`→`p:prio`); `d:defects`/`D:docs` unchanged and pinned by a dedicated test; render-policy coverage includes the audit overlay and finding detail.
- **T005** — Reverse lookups (`FindingsForTask`/`FindingsForEpic`) tolerate an absent register, backlink rows render in both detail overlays with a clear empty state, and enter/esc round-trips restore the originating overlay via `FindingDetailOrigin`.

**Audit applied 2026-07-03.** All proposed changes below were applied: Design.md now documents the Audit Register overlay and its `A` keybinding with `last_audited` advanced to this epic, the AGENTS.md Codebase Map covers the audit overlay and the data package's audit models/backlinks, E32 is marked `status: audited`, and the router advanced to `E33-audit-register-workflow-guidance/T001-audit-register-skill`. The router-lag and docs-drift findings from the initial audit are resolved by this apply.

Remaining risks carried forward (none blocking):

1. **Short task refs over-match across epics.** `FindingsForTask` matches a short `T###` link by number alone, so a finding linking `T001` appears in the Linked Findings section of *every* epic's T001 task. This mirrors `ValidateAuditFindings` (which also accepts short refs globally), but unlike `ResolveDependency` there is no same-epic scoping — and every epic has a T001. Low risk while findings use full task IDs; recommend the E33 workflow-guidance epic mandate full `E##-epic/T###-slug` IDs in finding `tasks:` links, or a follow-up that scopes short refs through the finding's `epics:` links.
2. **Arrow keys change meaning in detail overlays with linked findings.** When a task or epic has linked findings, `up`/`down`/`k`/`j` drive the finding cursor instead of scrolling the body; body scrolling falls back to `pgup`/`pgdown` only. The footer hint advertises this, but a long task body plus linked findings has no cursor follow-scroll (the linked section sits at the bottom), so the highlighted row can be off-screen until the user pages down. Acceptable for v1.4; worth revisiting if finding lists grow.
3. **Unrelated E22 launcher work is sitting in this working tree.** `internal/data/launcher_config.go` (+test), the `AgentLauncher` wiring in `internal/data/config.go`, `.savepoint/config.yml`, and the v1.3 E22 T001 task file are uncommitted alongside E32. They are out of this epic's scope and should be committed separately so the E32 diff stays reviewable (AGENTS.md "Small diffs").

## Code Style Review

- [x] One job per file — overlay rendering (`audit_overlay.go`), finding detail (`audit_detail.go`), and backlinks (`audit_backlinks.go` in both packages) are cleanly split.
- [x] One job per function — small named helpers throughout (`findingStatusLabel`, `findingAt`, `detailFooterHint`).
- [x] Test branches — populated/missing/malformed loads, cursor clamp, empty states, round-trips, and read-only guarantees are all covered.
- [x] Types document intent — `auditTab`, `findingLink`, `AuditState`, `OverlayType` origins.
- [x] Build only what is needed — no editing/mutation surface was built; out-of-scope items stayed out.
- [x] Handle errors at boundaries — tolerant loaders; explicit overlay load surfaces errors, aggregate reload degrades to empty with a debug log.
- [x] One source of truth — `auditFindingLayout` shared by renderer and follow-scroll; `auditTabLabels` sizes the tab strip; `footerHints` constant; status labels derived from status constants.
- [x] Comments explain WHY — e.g. why the layout pass is shared, why the footer separator changed, why best-effort load degrades.
- [x] Content in data files — hints and labels live in constants/data, not inline logic.
- [x] Small diffs — E32's own changes are proportionate; flagged the unrelated E22 files in the tree (Main Finding 3).

## Proposed Changes

### Target File
.savepoint/Design.md

### Replace
```md
last_audited: v1.4/E31-audit-register-data-model
```

### With
```md
last_audited: v1.4/E32-audit-register-tui-review
```

### Target File
.savepoint/Design.md

### Replace
```md
a header release indicator, and a top-level `D` Release Docs overlay for the selected release PRD plus overall PRD and Design (v1.2 E17-E33).
```

### With
```md
a header release indicator, a top-level `D` Release Docs overlay for the selected release PRD plus overall PRD and Design (v1.2 E17-E33), and a read-only top-level `A` Audit Register overlay with Prompt/Findings/Runs tabs, a grouped finding list, finding detail drill-in, and linked-finding backlinks inside task and epic detail overlays (v1.4 E32).
```

### Target File
.savepoint/Design.md

### Replace
```md
`d` opens defects, `D` opens Release Docs, `?` opens help, and `q` quits or closes overlays.
```

### With
```md
`d` opens defects, `D` opens Release Docs, `A` opens the Audit Register, `?` opens help, and `q` quits or closes overlays.
```

### Target File
AGENTS.md

### Replace
```md
| `internal/board/` | TUI board, overlays, epic sidebar, Next Activity line, router priority key, detail checklist rendering, status glyphs, forced color profile, debug logging hooks, async update I/O commands, defect summary/overlay/detail rendering, related-defect card markers, shared board utilities |
```

### With
```md
| `internal/board/` | TUI board, overlays, epic sidebar, Next Activity line, router priority key, detail checklist rendering, status glyphs, forced color profile, debug logging hooks, async update I/O commands, defect summary/overlay/detail rendering, related-defect card markers, audit register overlay with finding detail and linked-finding backlinks, shared board utilities |
```

### Target File
AGENTS.md

### Replace
```md
| `internal/data/` | Task/router/defect models, frontmatter parsing/splitting, lifecycle validation/defaulting, discovery including root-dir and release defect traversal, unified task status constants, canonical write helpers |
```

### With
```md
| `internal/data/` | Task/router/defect models, frontmatter parsing/splitting, lifecycle validation/defaulting, discovery including root-dir and release defect traversal, unified task status constants, canonical write helpers, audit-register models/loaders and finding backlink lookups |
```

### Target File
.savepoint/releases/v1.4/epics/E32-audit-register-tui-review/E32-Detail.md

### Replace
```md
type: epic-design
status: planned
```

### With
```md
type: epic-design
status: audited
```

### Target File
.savepoint/router.md

### Replace
```md
state: task-building
release: v1.4
epic: E32-audit-register-tui-review
task: E32-audit-register-tui-review/T005-work-item-finding-backlinks
next_action: Build E32-audit-register-tui-review/T005-work-item-finding-backlinks.
```

### With
```md
state: task-building
release: v1.4
epic: E33-audit-register-workflow-guidance
task: E33-audit-register-workflow-guidance/T001-audit-register-skill
next_action: Build E33-audit-register-workflow-guidance/T001-audit-register-skill.
```
