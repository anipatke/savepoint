import { describe, expect, it } from "vitest";
import { validateRouterState, ROUTER_STATES } from "../../src/domain/router.js";

const VALID = {
  state: "task-building",
  release: "v1",
  epic: "E02-data-model",
  next_action: "Build T004.",
};

describe("validateRouterState — valid", () => {
  it("accepts a full valid state object", () => {
    const r = validateRouterState(VALID);
    expect(r.ok).toBe(true);
    if (r.ok) {
      expect(r.value.state).toBe("task-building");
      expect(r.value.release).toBe("v1");
      expect(r.value.epic).toBe("E02-data-model");
      expect(r.value.next_action).toBe("Build T004.");
    }
  });

  it("accepts a state without the optional epic field", () => {
    const r = validateRouterState({ ...VALID, epic: undefined });
    expect(r.ok).toBe(true);
    if (r.ok) expect(r.value.epic).toBeUndefined();
  });

  it("accepts every valid state value", () => {
    for (const state of ROUTER_STATES) {
      const r = validateRouterState({ ...VALID, state });
      expect(r.ok).toBe(true);
    }
  });
});

describe("validateRouterState — non-object input", () => {
  it("rejects null", () => {
    expect(validateRouterState(null).ok).toBe(false);
  });
  it("rejects a string", () => {
    expect(validateRouterState("task-building").ok).toBe(false);
  });
  it("rejects an array", () => {
    expect(validateRouterState([VALID]).ok).toBe(false);
  });
});

describe("validateRouterState — state field", () => {
  it("rejects missing state", () => {
    const r = validateRouterState({ release: "v1", next_action: "x" });
    expect(r.ok).toBe(false);
    if (!r.ok) expect(r.reason).toMatch(/state/);
  });
  it("rejects an unknown state value", () => {
    const r = validateRouterState({ ...VALID, state: "unknown-state" });
    expect(r.ok).toBe(false);
    if (!r.ok) expect(r.reason).toMatch(/state/);
  });
  it("rejects a numeric state", () => {
    const r = validateRouterState({ ...VALID, state: 42 });
    expect(r.ok).toBe(false);
  });
});

describe("validateRouterState — release field", () => {
  it("rejects missing release", () => {
    const r = validateRouterState({ state: "task-building", next_action: "x" });
    expect(r.ok).toBe(false);
    if (!r.ok) expect(r.reason).toMatch(/release/);
  });
  it("rejects empty release", () => {
    const r = validateRouterState({ ...VALID, release: "   " });
    expect(r.ok).toBe(false);
    if (!r.ok) expect(r.reason).toMatch(/release/);
  });
});

describe("validateRouterState — next_action field", () => {
  it("rejects missing next_action", () => {
    const r = validateRouterState({ state: "task-building", release: "v1" });
    expect(r.ok).toBe(false);
    if (!r.ok) expect(r.reason).toMatch(/next_action/);
  });
  it("rejects empty next_action", () => {
    const r = validateRouterState({ ...VALID, next_action: "  " });
    expect(r.ok).toBe(false);
    if (!r.ok) expect(r.reason).toMatch(/next_action/);
  });
});
