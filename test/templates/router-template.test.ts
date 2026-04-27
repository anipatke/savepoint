import { describe, expect, it } from "vitest";
import { loadTemplate } from "../../src/templates/index.js";

describe("router template integrity", () => {
  it("contains all required router states", async () => {
    const r = await loadTemplate(".savepoint/router.md");
    expect(r.ok).toBe(true);
    if (!r.ok) return;

    const requiredStates = [
      "pre-implementation",
      "epic-design",
      "epic-task-breakdown",
      "task-planning",
      "task-building",
      "audit-pending",
    ];
    for (const state of requiredStates) {
      expect(r.value).toContain(state);
    }
  });

  it("contains read-order guidance", async () => {
    const r = await loadTemplate(".savepoint/router.md");
    expect(r.ok).toBe(true);
    if (!r.ok) return;
    expect(r.value).toContain("## Read order on every session");
    expect(r.value).toContain("This file (you are here)");
    expect(r.value).toContain("active epic Design");
  });

  it("contains conditional visual-identity read instruction", async () => {
    const r = await loadTemplate(".savepoint/router.md");
    expect(r.ok).toBe(true);
    if (!r.ok) return;
    expect(r.value).toContain("visual-identity.md");
    expect(r.value).toMatch(
      /Conditional read|conditional read|visual design|visual-identity/,
    );
  });

  it("contains embedded agent instruction markers", async () => {
    const r = await loadTemplate(".savepoint/router.md");
    expect(r.ok).toBe(true);
    if (!r.ok) return;
    expect(r.value).toContain("<!-- AGENT:");
  });
});
