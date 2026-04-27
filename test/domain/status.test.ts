import { describe, expect, it } from "vitest";
import {
  isTaskStatus,
  isTransitionAllowed,
  parseTaskStatus,
  TASK_STATUSES,
  validateTransition,
} from "../../src/domain/status.js";

describe("TASK_STATUSES", () => {
  it("contains the five canonical statuses", () => {
    expect(TASK_STATUSES).toEqual([
      "backlog",
      "planned",
      "in_progress",
      "review",
      "done",
    ]);
  });
});

describe("isTaskStatus", () => {
  it("returns true for each canonical status", () => {
    for (const s of TASK_STATUSES) {
      expect(isTaskStatus(s)).toBe(true);
    }
  });

  it("returns false for unknown strings", () => {
    expect(isTaskStatus("wip")).toBe(false);
    expect(isTaskStatus("")).toBe(false);
    expect(isTaskStatus("DONE")).toBe(false);
  });
});

describe("parseTaskStatus", () => {
  it("accepts each canonical status", () => {
    for (const s of TASK_STATUSES) {
      const r = parseTaskStatus(s);
      expect(r.ok).toBe(true);
      if (r.ok) expect(r.value).toBe(s);
    }
  });

  it("rejects unknown status and includes the value in the reason", () => {
    const r = parseTaskStatus("wip");
    expect(r.ok).toBe(false);
    if (!r.ok) expect(r.reason).toContain("wip");
  });
});

describe("isTransitionAllowed", () => {
  const allowed = [
    ["backlog", "planned"],
    ["planned", "in_progress"],
    ["in_progress", "review"],
    ["in_progress", "planned"],
    ["review", "done"],
    ["review", "in_progress"],
    ["done", "review"],
  ] as const;

  const disallowed = [
    ["backlog", "in_progress"],
    ["backlog", "review"],
    ["backlog", "done"],
    ["planned", "backlog"],
    ["planned", "review"],
    ["planned", "done"],
    ["in_progress", "backlog"],
    ["in_progress", "done"],
    ["review", "backlog"],
    ["review", "planned"],
    ["done", "backlog"],
    ["done", "planned"],
    ["done", "in_progress"],
  ] as const;

  for (const [from, to] of allowed) {
    it(`allows ${from} → ${to}`, () => {
      expect(isTransitionAllowed(from, to)).toBe(true);
    });
  }

  for (const [from, to] of disallowed) {
    it(`disallows ${from} → ${to}`, () => {
      expect(isTransitionAllowed(from, to)).toBe(false);
    });
  }
});

describe("validateTransition", () => {
  it("returns ok for allowed transition", () => {
    const r = validateTransition("planned", "in_progress");
    expect(r.ok).toBe(true);
    if (r.ok) {
      expect(r.value.from).toBe("planned");
      expect(r.value.to).toBe("in_progress");
    }
  });

  it("returns error for disallowed transition with both statuses in reason", () => {
    const r = validateTransition("planned", "done");
    expect(r.ok).toBe(false);
    if (!r.ok) {
      expect(r.reason).toContain("planned");
      expect(r.reason).toContain("done");
    }
  });

  it("returns ok for done → review so completed work can be reopened for audit-stale rechecks", () => {
    const r = validateTransition("done", "review");
    expect(r.ok).toBe(true);
    if (r.ok) {
      expect(r.value.from).toBe("done");
      expect(r.value.to).toBe("review");
    }
  });

  it("returns error for disallowed done transition with the allowed revert in reason", () => {
    const r = validateTransition("done", "in_progress");
    expect(r.ok).toBe(false);
    if (!r.ok) expect(r.reason).toContain("review");
  });
});
