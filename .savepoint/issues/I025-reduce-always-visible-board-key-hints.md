---
id: I025
title: Reduce the always-visible board key hints
type: other
status: resolved
source:
  kind: report
  actor: {role: owner, session: user}
  at: '2026-09-22T10:01:43Z'
severity: medium
resolution:
  disposition: accepted
  actor: {role: owner, session: user}
  at: '2026-09-22T11:16:02Z'
  reason: Owner explicitly accepted the remaining risk and requested this Issue be marked resolved.
history:
  - at: '2026-09-22T10:01:43Z'
    actor: {role: owner, session: user}
    kind: observed
    note: The Objective-focused footer exposes too many simultaneous key hints and is visually overwhelming; keep the footer concise and move the complete reference to Help.
  - at: '2026-09-22T10:26:36Z'
    actor: {role: owner, session: user}
    kind: owner_decision
    note: Objective navigation selects immediately on up/down, Tab is removed as a surface-focus shortcut, p remains a separate persisted router-selection action, and duplicate-Issue navigation is conditional rather than always advertised.
  - at: '2026-09-22T11:15:00Z'
    actor: {role: executor, session: i025-remediation-20260922}
    kind: repair_attempted
    note: "Repaired directly per the owner decision above: immediate up/down selection (moveObjectiveCursor now calls selectObjective; enter:select removed from handler and footer), Tab removed everywhere (key handler, toggleSidebarFocus, footer hints, Help row), the canonical hint gated on Issues.Detail.DuplicateTarget != nil, the dead I:task issues hint removed, and r:releases/v:detail/enter:detail/esc:clear hidden via new releaseHint/detailHint/clearObjectiveHint helpers (backed by a new hasDetailTarget bounds check shared with detailUnderCursor) whenever they would be no-ops. Help's closing line now reads esc/q:close rather than q:quit, matching that esc and q both only close Help there. Files changed: internal/board/v2/view.go, update.go, help.go, plus existing tests updated for the new tab-free/immediate-select navigation and a new footer_test.go covering each suppression rule. go build, go vet, go test ./internal/board/v2/..., git diff --check, make build, and make test all pass. The Issue remains open for independent Check verification."
  - at: '2026-09-22T11:16:02Z'
    actor: {role: owner, session: user}
    kind: owner_decision
    note: Owner accepted the remaining risk and requested this Issue be marked resolved; no technical CLEAR is implied.
---

# I025: Reduce the always-visible board key hints

## Summary

The Objective-focused footer currently displays navigation, card transfer,
Release selection, selection, detail, Issues, clear, tab, router selection,
Help, and quit actions at once. The row is difficult to scan and competes with
the board content it is meant to support.

The footer should show only the few actions most relevant to the current focus
and state. The complete keyboard map already has a dedicated Help overlay and
should remain discoverable there. Objective navigation should follow the V1
interaction: moving with up/down immediately changes the filtered Task view;
there is no separate `enter:select` step. Surface focus should be arrow-driven;
Tab is not a required or advertised board action.

## Evidence

- `internal/board/v2/view.go:386` composes the Objective-focused hint row from
  `↑↓`, `→`, `r`, `enter`, `v`, `i`, `esc`, `tab`, the contextual action, `?`,
  and `q`; the current row includes actions that are redundant, conditional,
  or misleading.
- The owner reported that this is too many always-visible options.
- `internal/board/v2/help.go` already owns the complete keyboard reference, so
  reducing footer density does not require removing functionality.
- `p:record selection` persists the focused Objective or Task in `router.md`;
  it is distinct from the transient Objective filter and should remain only as
  a clearly labelled contextual owner action.
- The Issues overlay advertises `I:task issues` without handling `I`, and
  `enter:canonical` has meaning only for an Issue carrying `duplicate_of`.

## Proof Needed

- Define a small, consistent footer-hint budget and hierarchy for Objective,
  Task-card, overlay, and diagnostic contexts.
- Make Objective up/down navigation apply the selection immediately; remove
  `enter:select` from the handler and footer.
- Remove the Tab handler, footer hints, and Help entry; arrow keys remain the
  supported way to move between the Objective sidebar and Task columns.
- Keep `p` as a separate, clearly labelled persisted router-selection action;
  keep owner actions visible only when they are available.
- Remove the dead `I:task issues` hint, and show canonical-Issue navigation
  only when the current Issue actually has a canonical target.
- Keep only the primary context-sensitive actions visible; preserve the full
  remaining key map in Help and ensure `?:help` remains readily discoverable.
- Suppress hints for empty or unavailable actions, including clearing an empty
  Objective selection, opening detail on an empty surface, and Release
  selection when no Release exists.
- Correct overlay wording so a close key is not labelled as quit when it only
  closes the current overlay.
- Do not duplicate a detailed reload diagnostic in the footer; coordinate with
  I020.
- Verify all keyboard actions still work even when omitted from the compact
  footer, and that context-sensitive owner actions remain visible when they
  become available.
- Test narrow and wide terminals, every focus surface, overlays, no-colour
  output, and line-width limits; pass focused board tests and project gates.
