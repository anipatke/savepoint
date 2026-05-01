# Custom Hooks Design Guide

## Table of Contents

1. [Design Principles](#design-principles)
2. [useScreenState](#usescreenstate)
3. [useTerminalSize](#useterminalsize)
4. [useAsyncData](#useasyncdata)
5. [useInput Patterns](#useinput-patterns)
6. [useSpinner](#usespinner)

---

## Design Principles

### Single Responsibility

Each hook handles exactly one concern.

```typescript
// Good: single responsibility
function useScreenState() { ... }   // screen navigation only
function useTerminalSize() { ... }  // size monitoring only

// Bad: multiple concerns
function useAppState() {
  // screen navigation, data fetching, terminal size...
}
```

### Return Objects

Group related state and functions in a returned object.

```typescript
interface UseAsyncDataResult<T> {
  data: T[];
  loading: boolean;
  error: Error | null;
  refresh: () => Promise<void>;
}

function useAsyncData<T>(): UseAsyncDataResult<T> {
  // ...
  return { data, loading, error, refresh };
}
```

### Cleanup

Always return a cleanup function from `useEffect`.

```typescript
useEffect(() => {
  const handler = () => { ... };
  process.stdout.on("resize", handler);

  return () => {
    process.stdout.removeListener("resize", handler);
  };
}, []);
```

### Exact Dependencies

Specify dependency arrays accurately for `useCallback` / `useEffect` / `useMemo`.

```typescript
const loadData = useCallback(async () => {
  setLoading(true);
  const data = await fetchData(options); // uses options
  setData(data);
  setLoading(false);
}, [options]); // options is a dependency
```

---

## useScreenState

Stack-based screen navigation.

```typescript
import { useState, useCallback } from "react";
import type { ScreenType } from "../types.js";

export interface ScreenStateResult {
  currentScreen: ScreenType;
  navigateTo: (screen: ScreenType) => void;
  goBack: () => void;
  reset: () => void;
}

const INITIAL_SCREEN: ScreenType = "board";

export function useScreenState(): ScreenStateResult {
  const [history, setHistory] = useState<ScreenType[]>([INITIAL_SCREEN]);

  const currentScreen = history[history.length - 1] ?? INITIAL_SCREEN;

  const navigateTo = useCallback((screen: ScreenType) => {
    setHistory((prev) => [...prev, screen]);
  }, []);

  const goBack = useCallback(() => {
    setHistory((prev) => {
      if (prev.length <= 1) {
        return prev; // Cannot go back from initial screen
      }
      return prev.slice(0, -1);
    });
  }, []);

  const reset = useCallback(() => {
    setHistory([INITIAL_SCREEN]);
  }, []);

  return {
    currentScreen,
    navigateTo,
    goBack,
    reset,
  };
}
```

### Usage

```typescript
function App() {
  const { currentScreen, navigateTo, goBack } = useScreenState();

  const renderScreen = () => {
    switch (currentScreen) {
      case "board":
        return (
          <BoardScreen
            onSelect={(task) => {
              setTask(task);
              navigateTo("detail");
            }}
          />
        );
      case "detail":
        return (
          <DetailScreen
            onBack={goBack}
            onAction={handleAction}
          />
        );
      default:
        return null;
    }
  };

  return <Box>{renderScreen()}</Box>;
}
```

---

## useTerminalSize

Terminal size detection and resize monitoring.

> **Note:** In non-TTY environments (CI, tests), `process.stdout.rows` /
> `columns` may be `undefined`. Always provide fallback values. Also check
> `process.stdout.isTTY` and skip resize event listeners in non-TTY mode.

```typescript
import { useState, useEffect } from "react";

export interface TerminalSize {
  rows: number;
  columns: number;
}

export function useTerminalSize(): TerminalSize {
  const [size, setSize] = useState<TerminalSize>(() => ({
    rows: process.stdout.rows || 24,
    columns: process.stdout.columns || 80,
  }));

  useEffect(() => {
    const handleResize = () => {
      setSize({
        rows: process.stdout.rows || 24,
        columns: process.stdout.columns || 80,
      });
    };

    process.stdout.on("resize", handleResize);
    return () => {
      process.stdout.removeListener("resize", handleResize);
    };
  }, []);

  return size;
}
```

### Usage

```typescript
function BoardScreen() {
  const { rows, columns } = useTerminalSize();

  const headerHeight = 2;
  const footerHeight = 1;
  const availableRows = rows - headerHeight - footerHeight;

  return (
    <Box flexDirection="column">
      <Header />
      <Select items={items} limit={Math.max(5, availableRows)} />
      <Footer />
    </Box>
  );
}
```

---

## useAsyncData

Async data fetching with loading, error, and refresh states.

```typescript
import { useState, useCallback, useEffect } from "react";

export interface UseAsyncDataOptions {
  enableAutoRefresh?: boolean;
  refreshInterval?: number;
}

export interface UseAsyncDataResult<T> {
  data: T[];
  loading: boolean;
  error: Error | null;
  lastUpdated: Date | null;
  refresh: () => Promise<void>;
}

export function useAsyncData<T>(
  fetcher: () => Promise<T[]>,
  options?: UseAsyncDataOptions,
): UseAsyncDataResult<T> {
  const { enableAutoRefresh = false, refreshInterval = 30000 } = options ?? {};

  const [data, setData] = useState<T[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);
  const [lastUpdated, setLastUpdated] = useState<Date | null>(null);

  const loadData = useCallback(async () => {
    setLoading(true);
    setError(null);

    try {
      const result = await fetcher();
      setData(result);
      setLastUpdated(new Date());
    } catch (err) {
      setError(err instanceof Error ? err : new Error(String(err)));
    } finally {
      setLoading(false);
    }
  }, [fetcher]);

  // Initial load
  useEffect(() => {
    loadData();
  }, [loadData]);

  // Auto-refresh
  useEffect(() => {
    if (!enableAutoRefresh) return;

    const intervalId = setInterval(loadData, refreshInterval);
    return () => clearInterval(intervalId);
  }, [enableAutoRefresh, refreshInterval, loadData]);

  return {
    data,
    loading,
    error,
    lastUpdated,
    refresh: loadData,
  };
}
```

### Cancel-Safe Pattern

Prevent state updates after unmount:

```typescript
const loadData = useCallback(async () => {
  let cancelled = false;

  setLoading(true);

  try {
    const data = await fetchData();

    if (!cancelled) {
      setData(data);
    }
  } catch (err) {
    if (!cancelled) {
      setError(err);
    }
  } finally {
    if (!cancelled) {
      setLoading(false);
    }
  }

  return () => {
    cancelled = true;
  };
}, []);
```

---

## useInput Patterns

### Basic Pattern

```typescript
import { useInput } from "ink";

function MyComponent() {
  useInput((input, key) => {
    if (key.escape) {
      handleEscape();
    }
    if (key.return) {
      handleEnter();
    }
    if (key.upArrow || input === "k") {
      moveUp();
    }
    if (key.downArrow || input === "j") {
      moveDown();
    }
  });
}
```

### Conditional Handling

```typescript
function MyComponent({ disabled, mode }: Props) {
  useInput((input, key) => {
    if (disabled) return;

    if (mode === "edit") {
      // Edit mode handling
    } else {
      // Normal mode handling
    }
  });
}
```

### Global Shortcuts

```typescript
function App() {
  const [filterMode, setFilterMode] = useState(false);

  useInput((input, key) => {
    if (filterMode) return; // Disabled while filtering

    switch (input) {
      case "q":
        handleQuit();
        break;
      case "r":
        handleRefresh();
        break;
      case "c":
        handleCleanup();
        break;
    }
  });

  return (
    <Box>
      <FilterInput
        onFocus={() => setFilterMode(true)}
        onBlur={() => setFilterMode(false)}
      />
    </Box>
  );
}
```

---

## useSpinner

Braille-pattern spinner for loading states.

```typescript
const SPINNER_FRAMES = ["⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧"];

function useSpinner(isActive: boolean) {
  const [frameIndex, setFrameIndex] = useState(0);
  const frameRef = useRef(0);

  useEffect(() => {
    if (!isActive) return;

    const interval = setInterval(() => {
      frameRef.current = (frameRef.current + 1) % SPINNER_FRAMES.length;
      setFrameIndex(frameRef.current);
    }, 80); // 80ms = 12.5 FPS

    return () => clearInterval(interval);
  }, [isActive]);

  return SPINNER_FRAMES[frameIndex];
}

// Usage
function LoadingIndicator({ loading }: { loading: boolean }) {
  const frame = useSpinner(loading);

  if (!loading) return null;

  return <Text>{frame} Loading...</Text>;
}
```
