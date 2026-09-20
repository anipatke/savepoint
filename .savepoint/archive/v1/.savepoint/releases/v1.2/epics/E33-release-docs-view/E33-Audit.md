---
type: audit-findings
audited: 2026-06-20
---

# Audit Findings: E33 Release Docs View

## Main Findings

E33's final implementation matches the post-review redesign in T005. The old
Epic detail Docs [3] subview has been removed; Epic detail is back to Detail [1]
and Audit [2]. Release documentation now lives in a top-level `D` Release Docs
overlay that loads the selected release PRD plus the project-wide PRD and Design
through the board command/message flow.

The previous audit finding about Markdown whitespace is resolved. The renderer
now uses `wrapDocText` through `styledWrap` instead of `WrapText`, so raw
document lines preserve blank lines, leading indentation, and interior spacing.
The render tests cover selector labels, selected body content, missing-doc empty
state, narrow width wrapping, unbreakable long tokens, indentation, and the
absence of the old Epic Docs tab strip.

Acceptance criteria are covered across data, IO, reducer, renderer, and view
tests. `LoadReleaseDocs(root, release)` returns exactly the bounded document set
for Release PRD, Overall PRD, and Overall Design; missing documents become
unavailable entries; unexpected read errors surface with path context; `D` opens
the overlay while lowercase `d` still opens defects; `esc` closes the overlay;
document selection and per-document scroll offsets are preserved.

The release-level PRD update requested with this re-audit has been applied:
`.savepoint/releases/v1.2/v1.2-PRD.md` now reflects the E21/E30-E33 release
work, including the Release Docs overlay. E33 has also been marked `audited`,
and the router has advanced to the first planned v1.3 task.

Verification run during re-audit and close: `make build && make test` passed.

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

The release-doc implementation stays within the established board/data boundary:
file loading is centralized in `internal/data`, UI state is in the board model,
filesystem reads happen through commands, and the renderer is reusable for the
single top-level Release Docs overlay. The re-audit PRD proposal has been
applied; no product-code follow-up remains.

## Proposed Changes

### Target File
.savepoint/releases/v1.2/v1.2-PRD.md

### Replace
```md
## Overview

v1.2 turns release-level defects into first-class workflow items without changing the existing three-column task board. It also trims the Savepoint template and skill surface so the workflow contract is easier for agents to follow, while keeping epics and tasks as the primary planning model. The release now also adds explicit task complexity so planning can reflect implementation difficulty in a consistent way, then closes recent lifecycle drift by making the task-to-epic workflow contract a single implementation surface.

## What ships in v1.2

1. **Defect Workflow TUI (E17)** — *Planned*
   - Release-level defect files under `.savepoint/releases/{release}/defects/`
   - Defect data model, router priority, board summary, overlay, detail view, and doctor validation
   - 6 tasks planned
2. **Template and Skill Optimisation (E18)** — *Planned*
   - Simplify the scaffolded and live Savepoint guidance by making skills canonical and trimming redundant prompts
   - Normalize workflow terminology across AGENTS, router, skills, and template tests
   - 3 tasks planned
3. **Task Complexity Field (E19)** — *Planned*
   - Add a complexity tier and short reason to task planning metadata
   - Surface complexity in the TUI and teach `create-task` to assign it from the shared rubric
   - Backfill the existing v1.2 task set with complexity metadata
   - 4 tasks planned
4. **Clean-up Lifecycle (E20)** — *Planned*
   - Centralize task lifecycle parsing, compatibility, validation, and transition rules
   - Align parser, writer, doctor, board transitions, templates, and skills around one contract
   - Add regression coverage for legacy lifecycle metadata and task-to-epic handoff behavior
   - 5 tasks planned

## Epic breakdown

| # | Epic | Status | Tasks |
|---|------|--------|-------|
| 17 | Defect Workflow TUI | Planned | 0/6 done |
| 18 | Template and Skill Optimisation | Planned | 0/3 done |
| 19 | Task Complexity Field | Planned | 0/4 done |
| 20 | Clean-up Lifecycle | Planned | 0/5 done |

## Success criteria

- Release-level defects can be tracked and surfaced in the TUI
- The board remains three columns wide
- Doctor validates defect files and broken defect references
- The template and skill suite has one canonical workflow source and no redundant phase prompt surface
- Task complexity is recorded, displayed, and validated consistently
- Task lifecycle rules have one source of truth across parser, writer, doctor, board, skills, and templates
- All E17, E18, E19, and E20 tasks are planned and ready to build
```

### With
```md
## Overview

v1.2 turns release-level defects into first-class workflow items without changing
the existing three-column task board. It also tightens Savepoint's agent-facing
workflow contract, lifecycle metadata, template surface, epic status handling,
and board context visibility so humans and agents can understand the active
release without leaving the TUI.

## What ships in v1.2

1. **Defect Workflow TUI (E17)** — *Audited*
   - Release-level defect files under `.savepoint/releases/{release}/defects/`
   - Defect data model, router priority, board summary, overlay, detail view, related-defect card markers, and doctor validation
   - 11 tasks done
2. **Template and Skill Optimisation (E18)** — *Audited*
   - Simplify scaffolded and live Savepoint guidance by making skills canonical and trimming redundant prompt surfaces
   - Normalize workflow terminology across AGENTS, router, skills, templates, and tests
   - 4 tasks done
3. **Task Complexity Field (E19)** — *Audited*
   - Add a complexity tier and short reason to task planning metadata
   - Surface complexity in the TUI and teach `create-task` to assign it from the shared rubric
   - Backfill the v1.2 task set with complexity metadata
   - 4 tasks done
4. **Clean-up Lifecycle (E20)** — *Audited*
   - Centralize task lifecycle parsing, compatibility, validation, and transition rules
   - Align parser, writer, doctor, board transitions, templates, and skills around one contract
   - Add regression coverage for legacy lifecycle metadata and task-to-epic handoff behavior
   - 6 tasks done
5. **Document Template Optimisation (E21)** — *Audited*
   - Refresh AGENTS, PRD, and Design templates and add a scaffolded Concept document
   - Keep live guidance aligned where it applies to active projects
   - 5 tasks done
6. **Epic Status Self-Heal (E30)** — *Audited*
   - Normalize epic status values on load so sidebar glyphs do not disappear
   - Surface non-canonical epic statuses through doctor diagnostics
   - 2 tasks done
7. **Mark Epic Audited Shortcut (E31)** — *Audited*
   - Add an Epic detail shortcut that writes `status: audited` through the existing data helper
   - Refresh in-memory epic status state so the board glyph updates immediately
   - 1 task done
8. **Header Release Indicator (E32)** — *Audited*
   - Render the selected release in the board header, with or without open defects
   - 1 task done
9. **Release Docs View (E33)** — *Audited*
   - Add a top-level `D` Release Docs overlay for the selected release PRD plus Overall PRD and Overall Design
   - Load docs through the data layer and board command/message flow
   - Preserve Markdown line structure while wrapping read-only document bodies
   - 5 tasks done

## Epic breakdown

| # | Epic | Status | Tasks |
|---|------|--------|-------|
| 17 | Defect Workflow TUI | Audited | 11/11 done |
| 18 | Template and Skill Optimisation | Audited | 4/4 done |
| 19 | Task Complexity Field | Audited | 4/4 done |
| 20 | Clean-up Lifecycle | Audited | 6/6 done |
| 21 | Document Template Optimisation | Audited | 5/5 done |
| 30 | Epic Status Self-Heal | Audited | 2/2 done |
| 31 | Mark Epic Audited Shortcut | Audited | 1/1 done |
| 32 | Header Release Indicator | Audited | 1/1 done |
| 33 | Release Docs View | Audited | 5/5 done |

## Success criteria

- Release-level defects can be tracked and surfaced in the TUI
- The board remains three columns wide
- Doctor validates defect files and broken defect references
- The template and skill suite has one canonical workflow source and no redundant phase prompt surface
- Task complexity is recorded, displayed, and validated consistently
- Task lifecycle rules have one source of truth across parser, writer, doctor, board, skills, and templates
- Epic status drift self-heals on board load and is visible to doctor diagnostics
- The board header shows the selected release
- The board can open read-only release documentation for the selected release without changing workflow state
- E17-E33 implementation tasks are complete and audited
```
