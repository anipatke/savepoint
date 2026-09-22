---
id: I025
title: Reduce the always-visible board key hints
type: other
status: open
source:
  kind: report
  actor: {role: owner, session: user}
  at: '2026-09-22T10:01:43Z'
severity: medium
history:
  - at: '2026-09-22T10:01:43Z'
    actor: {role: owner, session: user}
    kind: observed
    note: The Objective-focused footer exposes too many simultaneous key hints and is visually overwhelming; keep the footer concise and move the complete reference to Help.
---

# I025: Reduce the always-visible board key hints

## Summary

The Objective-focused footer currently displays navigation, card transfer,
Release selection, selection, detail, Issues, clear, tab, priority recording,
Help, and quit actions at once. The row is difficult to scan and competes with
the board content it is meant to support.

The footer should show only the few actions most relevant to the current focus
and state. The complete keyboard map already has a dedicated Help overlay and
should remain discoverable there.

## Evidence

- `internal/board/v2/view.go:386` composes the Objective-focused hint row from
  `↑↓`, `→`, `r`, `enter`, `v`, `i`, `esc`, `tab`, the contextual action, `?`,
  and `q`.
- The owner reported that this is too many always-visible options.
- `internal/board/v2/help.go` already owns the complete keyboard reference, so
  reducing footer density does not require removing functionality.

## Proof Needed

- Define a small, consistent footer-hint budget and hierarchy for Objective,
  Task-card, overlay, and diagnostic contexts.
- Keep only the primary context-sensitive actions visible; preserve the full
  key map in Help and ensure `?:help` remains readily discoverable.
- Do not duplicate a detailed reload diagnostic in the footer; coordinate with
  I020.
- Verify all keyboard actions still work even when omitted from the compact
  footer, and that context-sensitive owner actions remain visible when they
  become available.
- Test narrow and wide terminals, every focus surface, overlays, no-colour
  output, and line-width limits; pass focused board tests and project gates.
