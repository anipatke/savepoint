---
id: O-023
title: Show a Check's code style review on the board
status: in_progress
depends_on: [O-014]
release: R-006
---

# O-023: Show a Check's code style review on the board

## Outcome

A Task, Objective, or Goal detail overlay shows the `## Code Style Review`
checklist from that record's latest Check, exactly as the checker wrote it, so
the owner can see which `STYLE-*` rules were ticked without opening the Check
file.

## Why

V1's TUI Audit tab showed each epic audit's Code Style Review. V2 kept the
`STYLE-*` Guardrails but dropped both the checklist and the view. Commit
`2323ba8` restored the checklist in Check records through
`agent-skills/references/check-method.md`; the board still only shows a
Check's ID, result, and author, so the owner cannot inspect it there.

On 2026-09-24 the owner confirmed: show the style checklist only (not the whole
Check body), from the latest Check only, as a planned Objective in R-006 that
follows O-014.

## Success Conditions

- When the record's latest Check has a `## Code Style Review` section, the
  detail overlay shows a `CODE STYLE` section after `CHECKS`, naming that
  Check's ID and listing its lines verbatim, ticked and unticked state
  included.
- When the latest Check has no such section (every Check before C-916), the
  section says so and names the Check; it never shows an empty or invented
  checklist.
- When the record has no Check, no `CODE STYLE` section appears; `CHECKS`
  already says `(no Check recorded)`.
- A superseded Check's checklist is never shown.
- The non-TTY plain table and card badges are unchanged.

## Architectural Considerations

- The Check body is already loaded (`CheckV2.Source.Body`), and the latest
  Check is `index.LatestCheck[id]`. No data-layer, decoder, or schema change is
  needed.
- The style review is advisory and presentation only: it never feeds
  `ResolveClearance`, badges, gates, or Next.
- Extract the section while resolving `RecordDetail`, not in rendering
  (ARCH-02). Do not reuse `data.extractChecklistSection`; it drops the tick
  state.

## Boundaries

**In scope:** a `CODE STYLE` section in the V2 detail overlay for Task,
Objective, and Goal records; Design section 8 reconciliation.

**Out of scope:** showing whole Check bodies, selecting or opening individual
Checks, checklists from superseded Checks, per-Task style checklists for waived
Tasks, parsing or scoring STYLE rule IDs, and editing existing Check records.
