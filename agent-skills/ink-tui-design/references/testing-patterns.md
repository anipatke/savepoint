# Testing Patterns & Best Practices

## Table of Contents

1. [Test Environment Setup](#test-environment-setup)
2. [Component Tests](#component-tests)
3. [Mock Patterns](#mock-patterns)
4. [Integration Tests](#integration-tests)
5. [Edge Cases](#edge-cases)
6. [Performance Tests](#performance-tests)
7. [Test Utilities](#test-utilities)

---

## Test Environment Setup

### Vitest + happy-dom

```typescript
/**
 * @vitest-environment happy-dom
 */
import { describe, it, expect, beforeEach, vi } from "vitest";
import { render } from "@testing-library/react";
import { Window } from "happy-dom";

describe("Component", () => {
  beforeEach(() => {
    const window = new Window();
    // @ts-expect-error - happy-dom type mismatch
    globalThis.window = window;
    // @ts-expect-error - happy-dom type mismatch
    globalThis.document = window.document;
  });
});
```

### ink-testing-library

```typescript
import { render } from "ink-testing-library";

describe("Select", () => {
  it("should render items", () => {
    const items = [
      { label: "Item 1", value: "1" },
      { label: "Item 2", value: "2" },
    ];

    const { lastFrame } = render(
      <Select items={items} onSelect={vi.fn()} />,
    );

    expect(lastFrame()).toContain("Item 1");
    expect(lastFrame()).toContain("Item 2");
  });
});
```

---

## Component Tests

### Basic Pattern

```typescript
import { render } from "ink-testing-library";
import { describe, it, expect, vi } from "vitest";

describe("Header", () => {
  it("should render title", () => {
    const { lastFrame } = render(<Header title="Test" />);
    expect(lastFrame()).toContain("Test");
  });

  it("should render version when provided", () => {
    const { lastFrame } = render(
      <Header title="Test" version="1.0.0" />,
    );
    expect(lastFrame()).toContain("1.0.0");
  });
});
```

### Key Input Tests

```typescript
import { render } from "ink-testing-library";

describe("Select", () => {
  it("should move selection down on j key", async () => {
    const items = [
      { label: "Item 1", value: "1" },
      { label: "Item 2", value: "2" },
    ];

    const { stdin, lastFrame } = render(
      <Select items={items} onSelect={vi.fn()} />,
    );

    // Initial state
    expect(lastFrame()).toContain("› Item 1");

    // Press j to move down
    stdin.write("j");

    // Selection moved
    expect(lastFrame()).toContain("› Item 2");
  });

  it("should call onSelect on Enter", () => {
    const onSelect = vi.fn();
    const items = [{ label: "Item 1", value: "1" }];

    const { stdin } = render(
      <Select items={items} onSelect={onSelect} />,
    );

    stdin.write("\r"); // Enter key

    expect(onSelect).toHaveBeenCalledWith(items[0]);
  });
});
```

### Async Tests

```typescript
import { render } from "ink-testing-library";
import { vi } from "vitest";

describe("LoadingIndicator", () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it("should show after delay", async () => {
    const { lastFrame } = render(
      <LoadingIndicator isLoading={true} delay={300} />,
    );

    // Not shown initially
    expect(lastFrame()).toBe("");

    // Shown after 300ms
    vi.advanceTimersByTime(300);
    expect(lastFrame()).toContain("Loading");
  });
});
```

---

## Mock Patterns

### External Dependency Mock

```typescript
vi.mock("../../../fs/project.js", () => ({
  findProjectRoot: vi.fn().mockResolvedValue("/repo"),
  readConfig: vi.fn().mockResolvedValue({ defaultRelease: "v1" }),
  getSavepointPath: vi.fn().mockReturnValue("/repo/.savepoint"),
}));
```

### Screen Component Mock

Mock sibling screens to isolate the component under test:

```typescript
const capturedProps: BoardScreenProps[] = [];

vi.mock("../../components/screens/BoardScreen.js", () => ({
  BoardScreen: (props: BoardScreenProps) => {
    capturedProps.push(props);
    return <div>Mocked BoardScreen</div>;
  },
}));

describe("App", () => {
  beforeEach(() => {
    capturedProps.length = 0; // Reset
  });

  it("should pass tasks to BoardScreen", async () => {
    render(<App />);

    await waitFor(() => {
      expect(capturedProps.length).toBeGreaterThan(0);
    });

    expect(capturedProps[0].tasks).toHaveLength(2);
  });
});
```

### Hook Mock

```typescript
vi.mock("../../hooks/useAsyncData.js", () => ({
  useAsyncData: () => ({
    data: [
      { id: "T001", title: "Setup project" },
      { id: "T002", title: "Add tests" },
    ],
    loading: false,
    error: null,
    refresh: vi.fn(),
  }),
}));
```

### Conditional Mock Responses

```typescript
const mockFetchTasks = vi.fn();

vi.mock("../../../domain/tasks.js", () => ({
  fetchTasks: mockFetchTasks,
}));

describe("useAsyncData", () => {
  it("should handle error", async () => {
    mockFetchTasks.mockRejectedValueOnce(new Error("Read error"));
    // Assert error state is set
  });

  it("should return tasks", async () => {
    mockFetchTasks.mockResolvedValueOnce([{ id: "T001" }]);
    // Assert data is returned
  });
});
```

---

## Integration Tests

### Navigation Flow

```typescript
describe("Navigation", () => {
  it("should navigate from board to detail on Enter", async () => {
    const { stdin, lastFrame } = render(<App />);

    // Wait for board
    await waitFor(() => {
      expect(lastFrame()).toContain("Board");
    });

    // Press Enter
    stdin.write("\r");

    // Detail screen appears
    await waitFor(() => {
      expect(lastFrame()).toContain("Detail");
    });
  });

  it("should go back on Escape", async () => {
    const { stdin, lastFrame } = render(<App />);

    // Navigate to detail
    stdin.write("\r");
    await waitFor(() => {
      expect(lastFrame()).toContain("Detail");
    });

    // Press Escape
    stdin.write("\x1B"); // Escape key

    await waitFor(() => {
      expect(lastFrame()).toContain("Board");
    });
  });
});
```

---

## Edge Cases

```typescript
describe("Edge Cases", () => {
  it("should handle empty list", () => {
    const { lastFrame } = render(
      <Select items={[]} onSelect={vi.fn()} />,
    );

    expect(lastFrame()).not.toContain("undefined");
  });

  it("should handle very long labels", () => {
    const items = [{
      label: "A".repeat(200),
      value: "1",
    }];

    const { lastFrame } = render(
      <Select items={items} onSelect={vi.fn()} />,
    );

    // Assert output is defined (not crashed)
    expect(lastFrame()).toBeDefined();
  });
});
```

---

## Performance Tests

```typescript
describe("Performance", () => {
  it("should handle 1000 items without significant lag", async () => {
    const items = Array.from({ length: 1000 }, (_, i) => ({
      label: `Item ${i}`,
      value: String(i),
    }));

    const start = performance.now();

    const { stdin } = render(
      <Select items={items} onSelect={vi.fn()} limit={20} />,
    );

    // 100 navigation keypresses
    for (let i = 0; i < 100; i++) {
      stdin.write("j");
    }

    const duration = performance.now() - start;

    // Must complete within 1 second
    expect(duration).toBeLessThan(1000);
  });
});
```

---

## Test Utilities

### waitFor Helper

```typescript
async function waitFor(
  condition: () => boolean | void,
  timeout = 1000,
): Promise<void> {
  const start = Date.now();

  while (Date.now() - start < timeout) {
    try {
      if (condition()) return;
    } catch {
      // Condition not yet met
    }
    await new Promise((r) => setTimeout(r, 10));
  }

  throw new Error("waitFor timed out");
}
```

### Key Constants

```typescript
const KEYS = {
  ENTER: "\r",
  ESCAPE: "\x1B",
  UP: "\x1B[A",
  DOWN: "\x1B[B",
  LEFT: "\x1B[D",
  RIGHT: "\x1B[C",
  TAB: "\t",
  BACKSPACE: "\x7F",
  CTRL_C: "\x03",
};
```
