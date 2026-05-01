# Component Design Patterns

## Table of Contents

1. [Directory Structure](#directory-structure)
2. [Component Classification](#component-classification)
3. [Three-Layer Layout](#three-layer-layout)
4. [React.memo Optimization](#reactmemo-optimization)
5. [Controlled / Uncontrolled Modes](#controlled--uncontrolled-modes)
6. [Loop-Free Navigation](#loop-free-navigation)
7. [ErrorBoundary Pattern](#errorboundary-pattern)

---

## Directory Structure

Match the project layout under `src/tui/`:

```
src/tui/
├── App.tsx                  # Root component (screen management)
├── Board.tsx                # Main board screen
├── DetailPane.tsx           # Task detail overlay
├── state/
│   ├── reducer.ts           # Board selection, navigation, refresh
│   ├── app-reducer.ts       # App-level state
│   └── view-state.ts        # Derived view models
├── io/
│   ├── gates.ts             # Transition gate evaluation
│   └── write-status.ts      # File-backed status writes
├── theme/
│   ├── palette.ts           # Color constants
│   ├── capability.ts        # Terminal color detection
│   └── index.ts             # Theme exports
├── render/
│   ├── plain-table.ts       # Non-TTY table rendering
│   └── errors.ts            # Warning / error formatting
└── audit-review/            # Audit-specific screens
    ├── AuditReviewApp.tsx
    ├── ProposalList.tsx
    ├── OperationDetail.tsx
    ├── state.ts
    └── summary.ts
```

## Component Classification

### Screen

A complete interactive view.

**Responsibilities:**

- Handle keyboard input via `useInput`
- Use Header / Content / Footer three-layer layout
- Own screen-specific business logic

**Naming:** `<Concept>Screen.tsx` (or `App.tsx`, `Board.tsx` for top-level)

```typescript
interface BoardScreenProps {
  tasks: TaskItem[];
  stats: TaskStats;
  onSelect: (task: TaskItem) => void;
  onQuit?: () => void;
  onRefresh?: () => void;
  loading?: boolean;
  error?: Error | null;
}

export function BoardScreen({
  tasks,
  stats,
  onSelect,
  onQuit,
  onRefresh,
  loading,
  error,
}: BoardScreenProps) {
  const { rows } = useTerminalSize();

  useInput((input, key) => {
    if (key.escape || input === "q") {
      onQuit?.();
    }
    if (input === "r") {
      onRefresh?.();
    }
  });

  return (
    <Box flexDirection="column" height={rows}>
      <Header title="Savepoint Board" />
      <Box flexDirection="column" flexGrow={1}>
        {loading ? (
          <LoadingIndicator />
        ) : (
          <Select items={tasks} onSelect={onSelect} />
        )}
      </Box>
      <Footer actions={[{ key: "q", label: "Quit" }]} />
    </Box>
  );
}
```

### Part

Reusable, stateless rendering units.

**Responsibilities:**

- Display only
- Optimized with `React.memo`
- All data received via props

**Naming:** `<Feature>.tsx`

```typescript
interface HeaderProps {
  title: string;
  version?: string;
}

export const Header = React.memo(function Header({
  title,
  version,
}: HeaderProps) {
  return (
    <Box borderStyle="single" paddingX={1}>
      <Text bold>{title}</Text>
      {version && <Text dimColor> v{version}</Text>}
    </Box>
  );
});
```

### Common

Basic input primitives.

**Responsibilities:**

- Support both controlled and uncontrolled modes
- Generic interface reusable across the project
- May contain their own `useInput` for local keyboard handling

**Naming:** `<BaseType>.tsx` (e.g., `Select.tsx`, `Input.tsx`, `Confirm.tsx`)

---

## Three-Layer Layout

Standard screen layout pattern.

```typescript
export function StandardScreen({ children }: PropsWithChildren) {
  const { rows } = useTerminalSize();

  return (
    <Box flexDirection="column" height={rows}>
      {/* Layer 1: Header — fixed height */}
      <Header title="App Title" />

      {/* Layer 2: Content — flexGrow fills remaining height */}
      <Box flexDirection="column" flexGrow={1}>
        {children}
      </Box>

      {/* Layer 3: Footer — fixed height */}
      <Footer actions={footerActions} />
    </Box>
  );
}
```

### Dynamic Height Calculation

```typescript
const { rows } = useTerminalSize();

// Fixed section line counts
const headerLines = 2;    // border + title
const filterLines = 1;    // filter input
const statsLines = 1;     // statistics
const footerLines = 1;    // footer

const fixedLines = headerLines + filterLines + statsLines + footerLines;
const contentHeight = rows - fixedLines;
const listLimit = Math.max(5, contentHeight); // minimum 5 lines

return (
  <Select
    items={items}
    limit={listLimit}
    onSelect={handleSelect}
  />
);
```

---

## React.memo Optimization

### Basic Pattern

```typescript
export const Header = React.memo(function Header(props: HeaderProps) {
  return <Box>...</Box>;
});
```

### Custom Comparison Function

> **Important:** When using a custom comparison function, callback props
> (`onSelect`, etc.) must have stable references. Memoize callbacks in the
> parent with `useCallback`.

Compare by array contents rather than reference:

```typescript
function arePropsEqual<T extends SelectItem>(
  prevProps: SelectProps<T>,
  nextProps: SelectProps<T>,
): boolean {
  // Simple value comparison
  if (
    prevProps.limit !== nextProps.limit ||
    prevProps.disabled !== nextProps.disabled ||
    prevProps.onSelect !== nextProps.onSelect
  ) {
    return false;
  }

  // Array length comparison
  if (prevProps.items.length !== nextProps.items.length) {
    return false;
  }

  // Array element comparison
  for (let i = 0; i < prevProps.items.length; i++) {
    const prevItem = prevProps.items[i];
    const nextItem = nextProps.items[i];
    if (
      prevItem.value !== nextItem.value ||
      prevItem.label !== nextItem.label
    ) {
      return false;
    }
  }

  return true;
}

export const Select = React.memo(
  SelectComponent,
  arePropsEqual,
) as typeof SelectComponent;
```

---

## Controlled / Uncontrolled Modes

Allow parent components to control state when needed.

```typescript
interface SelectProps<T> {
  items: T[];
  onSelect: (item: T) => void;
  // Uncontrolled mode
  initialIndex?: number;
  // Controlled mode
  selectedIndex?: number;
  onSelectedIndexChange?: (index: number) => void;
}

function SelectComponent<T>({
  items,
  onSelect,
  initialIndex = 0,
  selectedIndex: externalSelectedIndex,
  onSelectedIndexChange,
}: SelectProps<T>) {
  // Internal state (uncontrolled)
  const [internalSelectedIndex, setInternalSelectedIndex] =
    useState(initialIndex);

  // Determine mode
  const isControlled = externalSelectedIndex !== undefined;
  const selectedIndex = isControlled
    ? externalSelectedIndex
    : internalSelectedIndex;

  // Unified update function
  const updateSelectedIndex = (value: number | ((prev: number) => number)) => {
    const newIndex = typeof value === "function" ? value(selectedIndex) : value;

    if (!isControlled) {
      setInternalSelectedIndex(newIndex);
    }

    onSelectedIndexChange?.(newIndex);
  };

  // ...
}
```

---

## Loop-Free Navigation

Stop at list boundaries rather than wrapping around.

```typescript
useInput((input, key) => {
  if (key.upArrow || input === "k") {
    // Move up (stop at 0)
    updateSelectedIndex((current) => Math.max(0, current - 1));
  } else if (key.downArrow || input === "j") {
    // Move down (stop at last item)
    updateSelectedIndex((current) => Math.min(items.length - 1, current + 1));
  }
});
```

---

## ErrorBoundary Pattern

> **Note:** ErrorBoundary must be a class component due to React constraints.
> This is the only exception to the "function components + hooks" rule.

```typescript
interface Props {
  children: React.ReactNode;
  fallback?: React.ReactNode;
}

interface State {
  hasError: boolean;
  error: Error | null;
}

export class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    console.error("ErrorBoundary caught:", error, errorInfo);
  }

  render() {
    if (this.state.hasError) {
      return this.props.fallback ?? (
        <Box>
          <Text color="red">Error: {this.state.error?.message}</Text>
        </Box>
      );
    }
    return this.props.children;
  }
}
```
