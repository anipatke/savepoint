import { describe, expect, it } from "vitest";
import {
  findCycle,
  findDuplicateIds,
  findMissingDependencies,
} from "../../src/validation/dependencies.js";
import {
  validateTaskFrontmatter,
  type TaskFrontmatter,
} from "../../src/domain/task.js";

function makeTask(id: string, deps: string[] = []): TaskFrontmatter {
  const r = validateTaskFrontmatter({
    id,
    status: "planned",
    objective: "test task",
    depends_on: deps,
  });
  if (!r.ok) throw new Error(`bad fixture: ${r.reason}`);
  return r.value;
}

const T1 = "E02-data-model/T001-alpha";
const T2 = "E02-data-model/T002-bravo";
const T3 = "E02-data-model/T003-charlie";

describe("findDuplicateIds", () => {
  it("returns empty for a collection with unique IDs", () => {
    const tasks = [makeTask(T1), makeTask(T2), makeTask(T3)];
    expect(findDuplicateIds(tasks)).toHaveLength(0);
  });

  it("returns empty for an empty collection", () => {
    expect(findDuplicateIds([])).toHaveLength(0);
  });

  it("reports a duplicate ID once when an ID appears twice", () => {
    const tasks = [makeTask(T1), makeTask(T2), makeTask(T1)];
    const dupes = findDuplicateIds(tasks);
    expect(dupes).toHaveLength(1);
    expect(dupes[0].raw).toBe(T1);
  });

  it("reports each distinct duplicate ID once when multiple IDs are duplicated", () => {
    const tasks = [makeTask(T1), makeTask(T2), makeTask(T1), makeTask(T2)];
    const dupes = findDuplicateIds(tasks);
    expect(dupes).toHaveLength(2);
    const raws = dupes.map((d) => d.raw).sort();
    expect(raws).toEqual([T1, T2].sort());
  });
});

describe("findMissingDependencies", () => {
  it("returns empty when all dependencies are present", () => {
    const tasks = [makeTask(T1), makeTask(T2, [T1]), makeTask(T3, [T1, T2])];
    expect(findMissingDependencies(tasks)).toHaveLength(0);
  });

  it("returns empty for a collection with no dependencies", () => {
    const tasks = [makeTask(T1), makeTask(T2)];
    expect(findMissingDependencies(tasks)).toHaveLength(0);
  });

  it("reports the dependent and missing ID when a dependency is absent", () => {
    const tasks = [makeTask(T2, [T1])];
    const missing = findMissingDependencies(tasks);
    expect(missing).toHaveLength(1);
    expect(missing[0].dependent.raw).toBe(T2);
    expect(missing[0].missing.raw).toBe(T1);
  });

  it("reports one entry per unresolved dependency", () => {
    const tasks = [makeTask(T3, [T1, T2])];
    const missing = findMissingDependencies(tasks);
    expect(missing).toHaveLength(2);
    const missingRaws = missing.map((m) => m.missing.raw).sort();
    expect(missingRaws).toEqual([T1, T2].sort());
  });
});

describe("findCycle", () => {
  it("returns null for an empty collection", () => {
    expect(findCycle([])).toBeNull();
  });

  it("returns null for a valid acyclic graph", () => {
    const tasks = [makeTask(T1), makeTask(T2, [T1]), makeTask(T3, [T1, T2])];
    expect(findCycle(tasks)).toBeNull();
  });

  it("returns null for isolated tasks with no dependencies", () => {
    const tasks = [makeTask(T1), makeTask(T2), makeTask(T3)];
    expect(findCycle(tasks)).toBeNull();
  });

  it("detects a simple two-node cycle", () => {
    const tasks = [makeTask(T1, [T2]), makeTask(T2, [T1])];
    const cycle = findCycle(tasks);
    expect(cycle).not.toBeNull();
    if (cycle) {
      const raws = cycle.map((t) => t.raw);
      expect(raws).toContain(T1);
      expect(raws).toContain(T2);
    }
  });

  it("detects a longer cycle and returns enough IDs to explain it", () => {
    // T1 → T2 → T3 → T1
    const tasks = [makeTask(T1, [T2]), makeTask(T2, [T3]), makeTask(T3, [T1])];
    const cycle = findCycle(tasks);
    expect(cycle).not.toBeNull();
    if (cycle) {
      const raws = cycle.map((t) => t.raw);
      expect(raws).toContain(T1);
      expect(raws).toContain(T2);
      expect(raws).toContain(T3);
      // first and last IDs must match to show the loop
      expect(raws[0]).toBe(raws[raws.length - 1]);
    }
  });

  it("ignores missing dependency IDs when detecting cycles", () => {
    // T1 depends on a task not in the collection — no cycle
    const tasks = [makeTask(T1, [T3])];
    expect(findCycle(tasks)).toBeNull();
  });
});
