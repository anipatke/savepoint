---
name: ink-tui-design
description: |
  Ink.js TUI design and implementation guide with visual design vocabulary.
  Use when:
  (1) Creating or modifying Ink.js components for the savepoint TUI
  (2) Implementing useInput/useApp hooks
  (3) Solving emoji/icon width issues (string-width)
  (4) Terminal responsive layout and dynamic height calculations
  (5) Keyboard input handling (Enter, Ctrl+C, conflicts)
  (6) Multi-width character (CJK) cursor issues
  (7) CLI UI test design (ink-testing-library)
  (8) Performance optimization (React.memo, custom comparators)
  (9) Visual design decisions for terminal rendering
  (10) Non-TTY fallback rendering
---

# Ink TUI Design

A unified guide for building terminal UIs with Ink.js, combining implementation
mechanics and visual design vocabulary. This skill owns the **how** of TUI
implementation. The **what** of visual direction (palette, spacing, signature
patterns) lives in `.savepoint/visual-identity.md` and always wins when the two
diverge.

## Quick Start

### New Component

1. Determine component type: Screen / Part / Common
2. Read [component-patterns.md](references/component-patterns.md) for the pattern
3. Read [ink-gotchas.md](references/ink-gotchas.md) if this component uses
   emoji, `useInput`, or `Ctrl+C` handling
4. Add types where needed
5. Implement the component
6. Write tests (see [testing-patterns.md](references/testing-patterns.md))

### Common Problem Solving

| Problem                  | Reference                                                           |
| ------------------------ | ------------------------------------------------------------------- |
| Emoji width misalignment | [ink-gotchas.md#icon-emoji-width](references/ink-gotchas.md)        |
| Ctrl+C fired twice       | [ink-gotchas.md#ctrlc-handling](references/ink-gotchas.md)          |
| useInput conflicts       | [ink-gotchas.md#useinput-conflicts](references/ink-gotchas.md)      |
| Enter key propagates     | [ink-gotchas.md#enter-double-fire](references/ink-gotchas.md)       |
| CJK cursor drift         | [ink-gotchas.md#cjk-cursor](references/ink-gotchas.md)              |
| Layout breaks on resize  | [ink-gotchas.md#layout-breaks](references/ink-gotchas.md)           |
| console.log interferes   | [ink-gotchas.md#consolelog-interference](references/ink-gotchas.md) |
| Ink color type errors    | [ink-gotchas.md#ink-color-types](references/ink-gotchas.md)         |

## Directory Conventions

Match the project structure under `src/tui/`:

```
src/tui/
├── App.tsx                  # Root Ink component — board, detail pane, input routing
├── Board.tsx                # Main kanban-style board view
├── DetailPane.tsx           # Expanded task detail overlay
├── state/                   # Pure reducer and view-state logic
│   ├── reducer.ts           # Board navigation, selection, refresh state
│   ├── app-reducer.ts       # App-level state (screen stack, mode)
│   └── view-state.ts        # Derived view models
├── io/                      # Transition gates and write helpers
│   ├── gates.ts             # Can-transition? rules before writing status
│   └── write-status.ts      # File-backed status updates with mtime checks
├── theme/                   # Color capability detection and palette
│   ├── index.ts             # Theme exports
│   ├── palette.ts           # Color constants (ANSI 256, true color fallbacks)
│   └── capability.ts        # Terminal color support detection
├── render/                  # Non-TTY and plain-text rendering
│   ├── plain-table.ts       # Deterministic table layout for CI/pipes
│   └── errors.ts            # Warning and error message formatting
└── audit-review/            # Audit mode screens (optional, per epic)
    ├── AuditReviewApp.tsx
    ├── ProposalList.tsx
    ├── OperationDetail.tsx
    ├── state.ts
    └── summary.ts
```

Rules:

- **One job per file.** A component that both renders and manages data state
  should be split into a Screen wrapper (state) and a Part (rendering).
- **State is explicit and local.** Keep UI state in the smallest surface that
  needs it. Derive view data from project files, never invent hidden state.
- **Test branching, navigation, render output, and non-TTY fallbacks.** Every
  component with `if/else/switch` gets a test. Pure rendering without branches
  may skip.

## Component Classification

### Screen

- Represents a complete interactive view
- Handles keyboard input via `useInput`
- Uses a 3-layer layout: Header / Content / Footer
- Owns screen-specific business logic and state

### Part

- Reusable, stateless rendering unit
- Optimized with `React.memo`
- Receives all data through props

### Common

- Basic input primitives (Select, Confirm, Input)
- Supports both controlled and uncontrolled modes

See [component-patterns.md](references/component-patterns.md) for complete
examples of each classification.

## Visual Design for Terminals

> **Override authority:** `.savepoint/visual-identity.md` owns the project's
> Atari-Noir design system (palette, spacing rhythm, signature patterns). The
> following is a general vocabulary. When they conflict, always defer to
> `visual-identity.md`.

### Box Drawing & Borders

Choose border styles that express intent:

| Style       | Characters | Mood                           |
| ----------- | ---------- | ------------------------------ |
| Single line | `┌─┐│└┘`   | Clean, modern, minimal         |
| Double line | `╔═╗║╚╝`   | Bold, formal, retro-mainframe  |
| Rounded     | `╭─╮│╰╯`   | Soft, friendly                 |
| Heavy       | `┏━┓┃┗┛`   | Strong, industrial             |
| ASCII       | `+-+\|`    | Retro, universal compatibility |

This project's default is **single line** in quiet `#1A1A1A` per
`visual-identity.md`. Use double line only for modal title bars or audit
headers. Avoid defaulting to plain ASCII unless targeting very old terminals.

### Color Strategies

| Strategy            | When to use                                                             |
| ------------------- | ----------------------------------------------------------------------- |
| ANSI 16             | Maximum compatibility; craft combinations beyond default red/green/blue |
| 256-color           | Richer palettes; gradients and subtle background variations             |
| True color (24-bit) | Full spectrum; smooth transitions, background tints                     |
| Monochrome          | Single color with intensity variations (dim, normal, bold, reverse)     |

This project prefers **true color with 256/16 fallbacks**. See `theme/capability.ts`
for detection and `theme/palette.ts` for the Atari-Noir color constants.

Create atmosphere with:

- Background color blocks for panels (not giant fills — intentional accents)
- Color-coded semantic meaning (reinforce with text; never rely on color alone)
- Inverted/reverse video for the active/focused row
- Dim text for secondary info, bold for primary

### Typography & Text Hierarchy

The terminal is all typography. Make it count:

- **Weight:** Bold for primary actions and headers; dim for metadata
- **Case:** UPPERCASE for section headers; sentence case for body
- **Symbols:** Enrich with `→ • ◆ ★ ⚡ λ ∴ ≡` and status glyphs
- **Custom bullets:** Replace `-` with `▸ ◉ ✓ ⬢ ›` where it adds clarity
- **Letter spacing:** Simulate with spaces for impact headers where the medium
  allows (fixed-width cells limit this)

### Layout Principles

- **Panels & Windows:** Distinct regions with quiet borders. Depth from
  contrast/glow, not heavy shadows.
- **Columns:** Side-by-side using careful spacing. Keep layout compact and
  readable on narrow terminals first.
- **Whitespace:** Generous padding inside panels; breathing room between
  sections.
- **Density:** Match to purpose. The board view is information-dense; audit
  review breathes more.
- **Hierarchy:** Clear visual distinction between primary content, secondary
  info, and chrome.

### Data Display

| Element           | Technique                                                        |
| ----------------- | ---------------------------------------------------------------- |
| Status indicators | `●` green/live, `○` empty, `◐` partial, `✓` complete, `✗` failed |
| Sparklines        | `▁▂▃▄▅▆▇█` for inline mini-charts                                |
| Horizontal bars   | Block characters for proportional comparison                     |
| Trees             | `├── └── │` for hierarchies                                      |
| Gauges            | `[████████░░]` with percentage                                   |

### Decorative Elements

- **Dividers:** `─────` `═══` `••••••` `░░░░░░` — match tone to section
- **Section markers:** `▶ SECTION`, `[ SECTION ]`, `─── SECTION ───`, `◆ SECTION`
- **Icons:** Prefer semantic Unicode or Nerd Font icons. Always apply width
  overrides for emoji (see Pattern 1 below).

### Anti-Patterns to Avoid

Never ship generic terminal aesthetics:

- Plain unformatted text output
- Default colors without intentional palette
- `[INFO]`, `[ERROR]` prefixes without styling
- Simple `----` dividers everywhere
- Walls of unstructured text
- Generic progress bars without personality
- Boring help text formatting
- Inconsistent spacing and alignment
- Accidental line wrapping (treat as a bug)
- Relying on color alone for state changes

## Key Patterns

### 1. Icon Width Overrides

`string-width` v8+ miscalculates some emoji widths. Override explicitly:

```typescript
const iconWidthOverrides: Record<string, number> = {
  "⚡": 1,
  "✨": 1,
  "✅": 1,
  "⚠️": 1,
  "🟢": 1,
  "🟠": 1,
};

const getIconWidth = (icon: string): number => {
  const baseWidth = stringWidth(icon);
  const override = iconWidthOverrides[icon];
  return override !== undefined ? Math.max(baseWidth, override) : baseWidth;
};

const padIcon = (icon: string, columnWidth = 2): string => {
  const width = getIconWidth(icon);
  const padding = Math.max(0, columnWidth - width);
  return icon + " ".repeat(padding);
};
```

### 2. useInput Conflict Avoidance

Multiple `useInput` hooks all receive input. Guard each handler:

```typescript
useInput((input, key) => {
  if (disabled) return; // Do nothing when this component is inactive
  if (key.return) {
    onSelect(selectedItem);
  }
});
```

### 3. Ctrl+C Handling

```typescript
render(<App />, { exitOnCtrlC: false });

// Inside the component
useInput((input, key) => {
  if (key.ctrl && input === "c") {
    cleanup();
    exit();
  }
});
```

### 4. Dynamic Height Calculation

```typescript
const { rows } = useTerminalSize();
const headerLines = 2;
const footerLines = 1;
const contentHeight = rows - headerLines - footerLines;
const listLimit = Math.max(5, contentHeight);
```

### 5. React.memo + Custom Comparison

```typescript
function arePropsEqual<T>(prev: Props<T>, next: Props<T>): boolean {
  if (prev.items.length !== next.items.length) return false;
  for (let i = 0; i < prev.items.length; i++) {
    if (prev.items[i].value !== next.items[i].value) return false;
  }
  return true;
}

export const Select = React.memo(SelectComponent, arePropsEqual);
```

## Detailed References

- **Ink.js Gotchas & Common Problems**: [references/ink-gotchas.md](references/ink-gotchas.md)
- **Component Design Patterns**: [references/component-patterns.md](references/component-patterns.md)
- **Custom Hooks Guide**: [references/hooks-guide.md](references/hooks-guide.md)
- **Testing Patterns**: [references/testing-patterns.md](references/testing-patterns.md)

## Visual Identity Authority

All aesthetic decisions — palette, spacing, signature patterns, motion, tone —
are owned by `.savepoint/visual-identity.md`. When implementing a TUI feature:

1. Read `visual-identity.md` first if the task touches rendering, theme, or
   visual design.
2. Return to this skill for Ink mechanics, component patterns, and testing.
3. If `visual-identity.md` says "quiet borders" and this skill says "double-line
   headers," `visual-identity.md` wins.
