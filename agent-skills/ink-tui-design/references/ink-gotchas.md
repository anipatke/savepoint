# Ink.js Gotchas & Common Problems

## Table of Contents

1. [Icon / Emoji Width Problems](#1-icon--emoji-width-problems)
2. [Ctrl+C Handling](#2-ctrlc-handling)
3. [useInput Conflicts](#3-useinput-conflicts)
4. [Enter Double-Fire](#4-enter-double-fire)
5. [CJK / Multi-Width Cursor](#5-cjk--multi-width-cursor)
6. [Layout Breaks](#6-layout-breaks)
7. [console.log Interference](#7-consolelog-interference)
8. [Ink Color Type Errors](#8-ink-color-type-errors)

---

## 1. Icon / Emoji Width Problems

### Problem

`string-width` v8+ changed how Variation Selectors (VS15/VS16) are handled,
causing emoji width calculations to be off. Terminals may render an emoji as
2 cells wide, but `stringWidth()` sometimes returns 1.

### Solution: WIDTH_OVERRIDES Pattern

```typescript
import stringWidth from "string-width";

const iconWidthOverrides: Record<string, number> = {
  "⚡": 1,
  "✨": 1,
  "🐛": 1,
  "🔥": 1,
  "🚀": 1,
  "📌": 1,
  "🟢": 1,
  "🟠": 1,
  "👉": 1,
  "💾": 1,
  "📤": 1,
  "🔃": 1,
  "✅": 1,
  "⚠️": 1,
  "🔗": 1,
  "💻": 1,
  "☁️": 1,
};

const getIconWidth = (icon: string): number => {
  const baseWidth = stringWidth(icon);
  const override = iconWidthOverrides[icon];
  return override !== undefined ? Math.max(baseWidth, override) : baseWidth;
};

// Fixed-width column padding
const COLUMN_WIDTH = 2;
const padIcon = (icon: string): string => {
  const width = getIconWidth(icon);
  const padding = Math.max(0, COLUMN_WIDTH - width);
  return icon + " ".repeat(padding);
};
```

### Affected Areas

- Status icon columns in board view
- Progress indicators
- Any fixed-width layout with inline symbols

---

## 2. Ctrl+C Handling

### Problem

Ink handles `Ctrl+C` by default to exit the app. This can fire twice or
conflict with custom cleanup logic.

### Solution: `exitOnCtrlC: false` + `useInput` + SIGINT

```typescript
import { render, useApp, useInput } from "ink";
import process from "node:process";

function App({ onExit }: { onExit: () => void }) {
  const { exit } = useApp();

  useInput((input, key) => {
    if (key.ctrl && input === "c") {
      onExit();
      exit();
    }
  });

  // Also handle SIGINT directly
  useEffect(() => {
    const handler = () => {
      onExit();
      exit();
    };
    process.on("SIGINT", handler);
    return () => process.off("SIGINT", handler);
  }, [exit, onExit]);

  return <Box>...</Box>;
}

// Disable Ink's default Ctrl+C handler
render(<App onExit={cleanup} />, { exitOnCtrlC: false });
```

---

## 3. useInput Conflicts

### Problem

When multiple `useInput` hooks exist, **all handlers fire** for every keypress.
A parent and child component may both process the same key.

### Solutions

**Pattern 1: `disabled` prop**

```typescript
useInput((input, key) => {
  if (disabled) return; // Ignore input when inactive
  // Handle key...
});
```

**Pattern 2: Mode flag**

```typescript
const [filterMode, setFilterMode] = useState(false);

useInput((input, key) => {
  if (filterMode) return; // Global shortcuts disabled while filtering
  if (input === "c") onCleanupCommand?.();
});
```

**Pattern 3: `blockKeys` prop on Input component**

```typescript
// Inside an Input component
useInput((input) => {
  if (blockKeys && blockKeys.includes(input)) {
    return; // Consume the key; do not propagate to parent
  }
});

// Usage
<Input blockKeys={["c", "r", "f"]} ... />
```

---

## 4. Enter Double-Fire

### Problem

After selecting an item in a Select component, the Enter key event propagates
to the next screen and triggers an unwanted action there.

### Solution: Ready-state buffering

```typescript
const [ready, setReady] = useState(false);

useEffect(() => {
  // Defer input handling until after the first render cycle
  const timer = setTimeout(() => setReady(true), 50);
  return () => clearTimeout(timer);
}, []);

useInput((input, key) => {
  if (!ready) return; // Drop events during initialization
  if (key.return) {
    onSelect(selectedItem);
  }
});
```

---

## 5. CJK / Multi-Width Cursor

### Problem

CJK characters and some symbols display as 2 cells wide but count as 1
character. Calculating cursor position by character count causes drift during
CJK input.

### Solution: Display-width based position calculation

```typescript
function getCharWidth(char: string): number {
  const code = char.codePointAt(0);
  if (!code) return 1;

  // CJK, Hangul, emoji ranges render as 2 cells wide
  if (
    (code >= 0x1100 && code <= 0x115f) || // Hangul Jamo
    (code >= 0x2e80 && code <= 0x9fff) || // CJK
    (code >= 0xac00 && code <= 0xd7af) || // Hangul Syllables
    (code >= 0xf900 && code <= 0xfaff) || // CJK Compatibility
    (code >= 0xfe10 && code <= 0xfe1f) || // Vertical forms
    (code >= 0x1f300 && code <= 0x1f9ff) // Emojis
  ) {
    return 2;
  }
  return 1;
}

function getDisplayWidth(str: string): number {
  return [...str].reduce((width, char) => width + getCharWidth(char), 0);
}

function toDisplayColumn(text: string, charPosition: number): number {
  return getDisplayWidth(text.slice(0, charPosition));
}
```

---

## 6. Layout Breaks

### Problem

- Long lines wrap past terminal width
- Timestamps or labels get pushed to the next line
- Column positions drift

### Solutions

**Safe margin**

```typescript
const { stdout } = useStdout();
const columns = Math.max(20, (stdout?.columns ?? 80) - 1); // 1-cell margin

// Or use 90% width (Gemini CLI pattern)
const safeColumns = Math.floor(columns * 0.9);
```

**Fixed-width columns**

```typescript
const COLUMN_WIDTH = 2;
const padToWidth = (content: string, width: number): string => {
  const actualWidth = stringWidth(content);
  const padding = Math.max(0, width - actualWidth);
  return content + " ".repeat(padding);
};
```

**Truncate to width**

```typescript
function truncateToWidth(text: string, maxWidth: number): string {
  let width = 0;
  let result = "";
  for (const char of text) {
    const charWidth = getCharWidth(char);
    if (width + charWidth > maxWidth) {
      return result + "…";
    }
    width += charWidth;
    result += char;
  }
  return result;
}
```

---

## 7. console.log Interference

### Problem

`console.log` output competes with Ink's rendering. Ink redraws the screen
periodically, so log lines get overwritten or corrupt the UI.

### Solutions

**Structured logs to stderr**

```typescript
import pino from "pino";

const logger = pino({
  transport: {
    target: "pino-pretty",
    options: { destination: 2 }, // stderr
  },
});
```

**Conditional debug logging**

```typescript
if (process.env.DEBUG) {
  console.error("Debug:", data); // stderr, not stdout
}
```

---

## 8. Ink Color Type Errors

### Problem

Using `<Text color="cyan">` in TypeScript can trigger strict type errors
because Ink's color types are narrow.

### Solutions

```typescript
// Use a const assertion
<Text color={"cyan" as const}>...</Text>

// Or define as a variable
const selectedColor = "cyan" as const;
<Text color={selectedColor}>...</Text>

// Or use chalk inside the Text node
import chalk from "chalk";
<Text>{chalk.cyan("Selected")}</Text>
```
